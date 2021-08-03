// Copyright 2021-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package store

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/atomix/atomix-go-client/pkg/atomix"
	"github.com/atomix/atomix-go-framework/pkg/atomix/meta"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/onosproject/onos-lib-go/pkg/logging"

	_map "github.com/atomix/atomix-go-client/pkg/atomix/map"
	"github.com/gogo/protobuf/proto"
	topoapi "github.com/onosproject/onos-api/go/onos/topo"
)

var log = logging.GetLogger("store", "topo")

// NewAtomixStore returns a new persistent Store
func NewAtomixStore(client atomix.Client) (Store, error) {
	objects, err := client.GetMap(context.Background(), "onos-topo-objects")
	if err != nil {
		return nil, err
	}
	store := &atomixStore{
		objects: objects,
		relations: relationMaps{
			targets: map[string][]topoapi.ID{},
			sources: map[string][]topoapi.ID{},
			tgtLock: sync.RWMutex{},
			srcLock: sync.RWMutex{},
		},
	}

	// watch the atomixStore for changes
	// when relations are deleted, remove the implied relations from the store target asnd source maps
	// when objects are deleted, remove their entries. The corresponding relations should be deleted as well by the delete method, so we do not have to search for them.
	// when objects are added, add their entry to the map
	// when a relation is added, add the implied relation to the store target and source maps
	mapCh := make(chan _map.Event)
	if err := objects.Watch(context.Background(), mapCh, make([]_map.WatchOption, 0)...); err != nil {
		// log.Errorf("Failed to start indexer: %s", err)
		return nil, errors.FromAtomix(err)
	}
	go func() {
		for event := range mapCh {
			obj, err := decodeObject(event.Entry)
			if err != nil {
				continue
			}

			switch event.Type {
			case _map.EventReplay:
				store.registerSrcTgt(obj)
			case _map.EventInsert:
				store.registerSrcTgt(obj)
			case _map.EventRemove:
				store.unregisterSrcTgt(obj)
			}

		}
	}()
	return store, nil
}

// Store stores topology information
type Store interface {
	io.Closer

	// Create creates an object in the store
	Create(ctx context.Context, object *topoapi.Object) error

	// Update updates an existing object in the store
	Update(ctx context.Context, object *topoapi.Object) error

	// Get retrieves an object from the store
	Get(ctx context.Context, id topoapi.ID) (*topoapi.Object, error)

	// Delete deletes a object from the store
	Delete(ctx context.Context, id topoapi.ID) error

	// List streams objects to the given channel
	List(ctx context.Context, filters *topoapi.Filters) ([]topoapi.Object, error)

	// Watch streams object events to the given channel
	Watch(ctx context.Context, ch chan<- topoapi.Event, filters *topoapi.Filters, opts ...WatchOption) error
}

// WatchOption is a configuration option for Watch calls
type WatchOption interface {
	apply([]_map.WatchOption) []_map.WatchOption
}

// watchReplyOption is an option to replay events on watch
type watchReplayOption struct {
}

// Temporary container used to help return relations, source entities, and target entities for relation filter
type relationTargetContainer struct {
	relation *topoapi.Object
	entity   *topoapi.Object
}

func (o watchReplayOption) apply(opts []_map.WatchOption) []_map.WatchOption {
	return append(opts, _map.WithReplay())
}

// WithReplay returns a WatchOption that replays past changes
func WithReplay() WatchOption {
	return watchReplayOption{}
}

// atomixStore is the object implementation of the Store
type atomixStore struct {
	objects   _map.Map
	relations relationMaps
}

type relationMaps struct {
	targets map[string][]topoapi.ID
	sources map[string][]topoapi.ID
	tgtLock sync.RWMutex
	srcLock sync.RWMutex
}

func (s *atomixStore) Create(ctx context.Context, object *topoapi.Object) error {
	// If an object is a relation and its ID is empty, build one.
	if object.ID == "" && object.Type == topoapi.Object_RELATION {
		relation := object.GetRelation()
		object.ID = topoapi.RelationID(relation.SrcEntityID, relation.KindID, relation.TgtEntityID)
	}
	if object.ID == "" {
		return errors.NewInvalid("ID cannot be empty")
	}
	if object.Type == topoapi.Object_UNSPECIFIED {
		return errors.NewInvalid("Type cannot be unspecified")
	}

	log.Infof("Creating object %+v", object)
	bytes, err := proto.Marshal(object)
	if err != nil {
		log.Errorf("Failed to create object %+v: %s", object, err)
		return errors.NewInvalid(err.Error())
	}

	// Put the object in the map using an optimistic lock if this is an update
	entry, err := s.objects.Put(ctx, string(object.ID), bytes, _map.IfNotSet())
	if err != nil {
		log.Errorf("Failed to create object %+v: %s", object, err)
		return errors.FromAtomix(err)
	}

	object.Revision = topoapi.Revision(entry.Revision)
	return nil
}

func (s *atomixStore) Update(ctx context.Context, object *topoapi.Object) error {
	if object.ID == "" {
		return errors.NewInvalid("ID cannot be empty")
	}
	if object.Type == topoapi.Object_UNSPECIFIED {
		return errors.NewInvalid("Type cannot be unspecified")
	}
	if object.Revision == 0 {
		return errors.NewInvalid("object must contain a revision on update")
	}

	log.Infof("Updating object %+v", object)
	bytes, err := proto.Marshal(object)
	if err != nil {
		log.Errorf("Failed to update object %+v: %s", object, err)
		return errors.NewInvalid(err.Error())
	}

	// Update the object in the map
	entry, err := s.objects.Put(ctx, string(object.ID), bytes, _map.IfMatch(meta.NewRevision(meta.Revision(object.Revision))))
	if err != nil {
		log.Errorf("Failed to update object %+v: %s", object, err)
		return errors.FromAtomix(err)
	}
	object.Revision = topoapi.Revision(entry.Revision)
	return nil
}

func (s *atomixStore) Get(ctx context.Context, id topoapi.ID) (*topoapi.Object, error) {
	if id == "" {
		return nil, errors.NewInvalid("ID cannot be empty")
	}

	entry, err := s.objects.Get(ctx, string(id))
	if err != nil {
		return nil, errors.FromAtomix(err)
	}
	obj, err := decodeObject(*entry)
	if err != nil {
		return nil, err
	}
	s.addSrcTgts(obj)
	return obj, nil
}

func (s *atomixStore) Delete(ctx context.Context, id topoapi.ID) error {
	if id == "" {
		return errors.NewInvalid("ID cannot be empty")
	}

	s.relations.srcLock.Lock()
	defer s.relations.srcLock.Unlock()
	s.relations.tgtLock.Lock()
	defer s.relations.tgtLock.Unlock()
	// access the object to determine its properties
	mapObj, err := s.objects.Get(ctx, string(id))
	if err != nil {
		return errors.FromAtomix(err)
	}
	obj, err := decodeObject(*mapObj)
	if err != nil {
		return err
	}

	if obj.GetEntity() != nil {
		// delete the relations
		mapCh := make(chan _map.Entry)
		if err := s.objects.Entries(ctx, mapCh); err != nil {
			return errors.FromAtomix(err)
		}
		for entry := range mapCh {
			if ep, err := decodeObject(entry); err == nil {
				// if object is a relation and its kind and src id matches the filter, create blank entry for its target id
				if ep.Type == topoapi.Object_RELATION && (ep.GetRelation().GetSrcEntityID() == obj.ID || ep.GetRelation().GetTgtEntityID() == obj.ID) {
					// the deletion of the relation should trigger the watch to update the store maps
					_, err = s.objects.Remove(ctx, string(ep.ID))
					if err != nil {
						return errors.FromAtomix(err)
					}
				}
			}
		}
	}

	log.Infof("Deleting object %s", id)
	_, err = s.objects.Remove(ctx, string(id))
	if err != nil {
		log.Errorf("Failed to delete object %s: %s", id, err)
		return errors.FromAtomix(err)
	}
	return nil
}

func (s *atomixStore) List(ctx context.Context, filters *topoapi.Filters) ([]topoapi.Object, error) {
	mapCh := make(chan _map.Entry)
	if err := s.objects.Entries(ctx, mapCh); err != nil {
		return nil, errors.FromAtomix(err)
	}

	eps := make([]topoapi.Object, 0)

	// first make sure there are filters. if there aren't, return everything with the correct type
	if filters == nil {
		for entry := range mapCh {
			if ep, err := decodeObject(entry); err == nil {
				s.addSrcTgts(ep)
				eps = append(eps, *ep)
			}
		}
		return eps, nil
	}

	if filters.RelationFilter != nil {
		return s.listRelationFilter(ctx, mapCh, filters, eps)
	}

	for entry := range mapCh {
		if ep, err := decodeObject(entry); err == nil {
			if match(ep, filters) {
				if matchType(ep, filters.ObjectTypes) {
					s.addSrcTgts(ep)
					eps = append(eps, *ep)
				}
			}
		}
	}

	return eps, nil
}

func (s *atomixStore) listRelationFilter(ctx context.Context, mapCh chan _map.Entry, filters *topoapi.Filters, eps []topoapi.Object) ([]topoapi.Object, error) {
	filter := filters.RelationFilter

	// contains _all_ relations that have the same kind as the filter and same SrcId as the filter
	entitiesToGet := make(map[topoapi.ID]relationTargetContainer)
	for entry := range mapCh {
		if ep, err := decodeObject(entry); err == nil {
			// if object is a relation and its kind and src id matches the filter, create blank entry for its target id
			if ep.Type == topoapi.Object_RELATION && string(ep.GetRelation().KindID) == filter.GetRelationKind() && string(ep.GetRelation().GetSrcEntityID()) == filter.SrcId {
				entitiesToGet[ep.GetRelation().TgtEntityID] = relationTargetContainer{relation: ep, entity: nil}
			} else
			// if object is an entity, see if satisfies the filter and set its value in entitiesToGet
			if ep.Type == topoapi.Object_ENTITY {
				if value, found := entitiesToGet[ep.ID]; found {
					if filter.TargetKind == "" || string(ep.GetEntity().KindID) == filter.TargetKind {
						temp := value
						temp.entity = ep
						entitiesToGet[ep.ID] = temp
					}
				}
			}
		}
	}
	// to prevent adding a node twice. each relation filter must specify a source id, so we will only ever want to add one node (source)
	foundSource := false
	// iterate over entitiesToGet to obtain missed entities and push onto eps
	for id, relationEntity := range entitiesToGet {

		if relationEntity.entity == nil {
			storeEntity, _ := s.Get(ctx, id)
			if filter.TargetKind == "" || string(storeEntity.GetEntity().KindID) == filter.TargetKind {
				if matchType(storeEntity, filters.ObjectTypes) {
					s.addSrcTgts(storeEntity)
					eps = append(eps, *storeEntity)
					if filter.Scope == topoapi.RelationFilterScope_ALL {
						eps = append(eps, *relationEntity.relation)
					}
					if !foundSource && (filter.Scope == topoapi.RelationFilterScope_ALL || filter.Scope == topoapi.RelationFilterScope_SOURCE_AND_TARGET) {
						src, err := s.Get(ctx, relationEntity.relation.GetRelation().SrcEntityID)
						if err != nil {
							return nil, err
						}
						s.addSrcTgts(src)
						eps = append(eps, *src)

						foundSource = true
					}
				}
			}
		} else {
			if matchType(relationEntity.entity, filters.ObjectTypes) {
				s.addSrcTgts(relationEntity.entity)
				eps = append(eps, *relationEntity.entity)
				if filter.Scope == topoapi.RelationFilterScope_ALL {
					eps = append(eps, *relationEntity.relation)
				}
				if !foundSource && (filter.Scope == topoapi.RelationFilterScope_ALL || filter.Scope == topoapi.RelationFilterScope_SOURCE_AND_TARGET) {
					src, err := s.Get(ctx, relationEntity.relation.GetRelation().SrcEntityID)
					if err != nil {
						return nil, err
					}
					s.addSrcTgts(src)
					eps = append(eps, *src)
					foundSource = true
				}
			}
		}
	}
	return eps, nil
}

func (s *atomixStore) Watch(ctx context.Context, ch chan<- topoapi.Event, filters *topoapi.Filters, opts ...WatchOption) error {
	watchOpts := make([]_map.WatchOption, 0)
	for _, opt := range opts {
		watchOpts = opt.apply(watchOpts)
	}

	mapCh := make(chan _map.Event)
	if err := s.objects.Watch(ctx, mapCh, watchOpts...); err != nil {
		return errors.FromAtomix(err)
	}

	go func() {
		defer close(ch)
		for event := range mapCh {
			if object, err := decodeObject(event.Entry); err == nil {
				if !match(object, filters) {
					continue
				}
				var eventType topoapi.EventType
				switch event.Type {
				case _map.EventReplay:
					eventType = topoapi.EventType_NONE
				case _map.EventInsert:
					eventType = topoapi.EventType_ADDED
				case _map.EventRemove:
					eventType = topoapi.EventType_REMOVED
				case _map.EventUpdate:
					eventType = topoapi.EventType_UPDATED
				default:
					eventType = topoapi.EventType_UPDATED
				}
				ch <- topoapi.Event{
					Type:   eventType,
					Object: *object,
				}
			}
		}
	}()
	return nil
}

func (s *atomixStore) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	_ = s.objects.Close(ctx)
	defer cancel()
	return s.objects.Close(ctx)
}

func decodeObject(entry _map.Entry) (*topoapi.Object, error) {
	object := &topoapi.Object{}
	if err := proto.Unmarshal(entry.Value, object); err != nil {
		return nil, errors.NewInvalid(err.Error())
	}
	object.ID = topoapi.ID(entry.Key)
	object.Revision = topoapi.Revision(entry.Revision)
	return object, nil
}

func (s *atomixStore) addSrcTgts(obj *topoapi.Object) {
	if obj.GetEntity() != nil {
		s.relations.srcLock.RLock()
		defer s.relations.srcLock.RUnlock()
		s.relations.tgtLock.RLock()
		defer s.relations.tgtLock.RUnlock()
		if s.relations.sources[string(obj.ID)] != nil {
			(*obj.GetEntity()).SrcRelationIDs = s.relations.sources[string(obj.ID)]
		}
		if s.relations.targets[string(obj.ID)] != nil {
			obj.GetEntity().TgtRelationIDs = s.relations.targets[string(obj.ID)]
		}
	}
}

// when deleting either a relation or entity, create the correspending entries in the store
func (s *atomixStore) registerSrcTgt(obj *topoapi.Object) {
	if entity := obj.GetEntity(); entity != nil {
		s.relations.srcLock.Lock()
		defer s.relations.srcLock.Unlock()
		s.relations.tgtLock.Lock()
		defer s.relations.tgtLock.Unlock()
		s.relations.sources[string(obj.ID)] = make([]topoapi.ID, 0)
		s.relations.targets[string(obj.ID)] = make([]topoapi.ID, 0)
	} else if relation := obj.GetRelation(); relation != nil {
		s.relations.srcLock.Lock()
		defer s.relations.srcLock.Unlock()
		s.relations.tgtLock.Lock()
		defer s.relations.tgtLock.Unlock()
		if list, found := s.relations.sources[string(relation.TgtEntityID)]; found {
			s.relations.sources[string(relation.TgtEntityID)] = append(list, relation.SrcEntityID)
		}
		if list, found := s.relations.targets[string(relation.SrcEntityID)]; found {
			s.relations.targets[string(relation.SrcEntityID)] = append(list, relation.TgtEntityID)
		}
	}
}

// when deleting either a relation or entity, remove the correspending entries in the store
func (s *atomixStore) unregisterSrcTgt(obj *topoapi.Object) {
	if entity := obj.GetEntity(); entity != nil {
		s.relations.srcLock.Lock()
		defer s.relations.srcLock.Unlock()
		s.relations.tgtLock.Lock()
		defer s.relations.tgtLock.Unlock()
		delete(s.relations.sources, string(obj.ID))
		delete(s.relations.targets, string(obj.ID))
	} else if relation := obj.GetRelation(); relation != nil {
		s.relations.srcLock.Lock()
		defer s.relations.srcLock.Unlock()
		s.relations.tgtLock.Lock()
		defer s.relations.tgtLock.Unlock()
		if list, found := s.relations.sources[string(relation.TgtEntityID)]; found {
			index := 0
			for _, id := range list {
				if id != relation.SrcEntityID {
					list[index] = id
					index++
				}
			}
		}
		if list, found := s.relations.targets[string(relation.SrcEntityID)]; found {
			index := 0
			for _, id := range list {
				if id != relation.TgtEntityID {
					list[index] = id
					index++
				}
			}
		}
	}
}
