// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package store

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/atomix/atomix-go-client/pkg/atomix"
	"github.com/atomix/atomix-go-framework/pkg/atomix/meta"
	"github.com/google/uuid"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/onosproject/onos-lib-go/pkg/logging"

	_map "github.com/atomix/atomix-go-client/pkg/atomix/map"
	"github.com/gogo/protobuf/proto"
	topoapi "github.com/onosproject/onos-api/go/onos/topo"
)

var log = logging.GetLogger()

// NewAtomixStore returns a new persistent Store
func NewAtomixStore(client atomix.Client) (Store, error) {
	objects, err := client.GetMap(context.Background(), "onos-topo-objects")
	if err != nil {
		return nil, err
	}
	store := &atomixStore{
		objects:  objects,
		cache:    make(map[topoapi.ID]topoapi.Object),
		watchers: make(map[uuid.UUID]chan<- topoapi.Event),
		relations: relationMaps{
			targets: make(map[topoapi.ID][]topoapi.ID),
			sources: make(map[topoapi.ID][]topoapi.ID),
			lock:    sync.RWMutex{},
		},
	}

	// watch the atomixStore for changes
	// when relations are deleted, remove the implied relations from the store target asnd source maps
	// when objects are deleted, remove their entries. The corresponding relations should be deleted as well by the delete method, so we do not have to search for them.
	// when objects are added, add their entry to the map
	// when a relation is added, add the implied relation to the store target and source maps
	mapCh := make(chan _map.Event)
	if err := objects.Watch(context.Background(), mapCh, _map.WithReplay()); err != nil {
		return nil, errors.FromAtomix(err)
	}
	go store.watchStoreEvents(mapCh)
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
	Delete(ctx context.Context, id topoapi.ID, revision topoapi.Revision) error

	// List streams objects to the given channel
	List(ctx context.Context, filters *topoapi.Filters) ([]topoapi.Object, error)

	// Watch streams object events to the given channel
	Watch(ctx context.Context, ch chan<- topoapi.Event, filters *topoapi.Filters, opts ...WatchOption) error
}

// WatchOption is a configuration option for Watch calls
type WatchOption interface {
	apply(*watchOptions)
}

// watchReplyOption is an option to replay events on watch
type watchReplayOption struct {
	replay bool
}

func (o watchReplayOption) apply(opts *watchOptions) {
	opts.replay = o.replay
}

// WithReplay returns a WatchOption that replays past changes
func WithReplay() WatchOption {
	return watchReplayOption{true}
}

type watchOptions struct {
	replay bool
}

// atomixStore is the object implementation of the Store
type atomixStore struct {
	objects    _map.Map
	cache      map[topoapi.ID]topoapi.Object
	cacheMu    sync.RWMutex
	relations  relationMaps
	watchers   map[uuid.UUID]chan<- topoapi.Event
	watchersMu sync.RWMutex
}

type relationMaps struct {
	// map of entity IDs to list of relations where that entity is a source of the relation
	sources map[topoapi.ID][]topoapi.ID
	// map of entity IDs to list of relations where that entity is a target of the relation
	targets map[topoapi.ID][]topoapi.ID
	lock    sync.RWMutex
}

func (s *atomixStore) watchStoreEvents(mapCh chan _map.Event) {
	for event := range mapCh {
		obj, err := decodeObject(event.Entry)
		if err != nil {
			continue
		}

		var eventType topoapi.EventType
		switch event.Type {
		case _map.EventReplay:
			eventType = topoapi.EventType_NONE
			s.cacheMu.Lock()
			s.cache[obj.ID] = *obj
			s.cacheMu.Unlock()
			s.registerSrcTgt(obj, true)
		case _map.EventInsert:
			eventType = topoapi.EventType_ADDED
			s.cacheMu.Lock()
			s.cache[obj.ID] = *obj
			s.cacheMu.Unlock()
			s.registerSrcTgt(obj, true)
		case _map.EventUpdate:
			eventType = topoapi.EventType_UPDATED
			s.cacheMu.Lock()
			s.cache[obj.ID] = *obj
			s.cacheMu.Unlock()
		case _map.EventRemove:
			eventType = topoapi.EventType_REMOVED
			s.cacheMu.Lock()
			delete(s.cache, topoapi.ID(event.Entry.Key))
			s.cacheMu.Unlock()
			s.unregisterSrcTgt(obj)
		}

		s.watchersMu.RLock()
		for _, watcher := range s.watchers {
			watcher <- topoapi.Event{
				Type:   eventType,
				Object: *obj,
			}
		}
		s.watchersMu.RUnlock()
	}
}

func (s *atomixStore) Create(ctx context.Context, object *topoapi.Object) error {
	if object.Type == topoapi.Object_UNSPECIFIED {
		return errors.NewInvalid("Type cannot be unspecified")
	}

	// set a uuid
	uuid, err := uuid.NewRandom()
	if err != nil {
		return errors.FromAtomix(err)
	}
	object.UUID = topoapi.UUID(uuid.String())
	// If an object is a relation and its ID is empty, build one.
	if object.Type == topoapi.Object_RELATION {
		if object.ID == "" {
			object.ID = topoapi.ID("uuid:" + string(object.UUID))
		}
		if _, err := s.objects.Get(ctx, string(object.GetRelation().SrcEntityID)); err != nil {
			err = errors.FromAtomix(err)
			if !errors.IsNotFound(err) {
				log.Errorf("Failed to create Object %+v: %v", object, err)
				return err
			}
			log.Warnf("Source Entity does not exist")
			return errors.NewInvalid("Source Entity does not exist")
		}
		if _, err := s.objects.Get(ctx, string(object.GetRelation().TgtEntityID)); err != nil {
			err = errors.FromAtomix(err)
			if !errors.IsNotFound(err) {
				log.Errorf("Failed to create Object %+v: %v", object, err)
				return err
			}
			log.Warnf("Target Entity does not exist")
			return errors.NewInvalid("Target Entity does not exist")
		}
	} else if object.ID == "" {
		return errors.NewInvalid("ID cannot be empty")
	}

	log.Infof("Creating Object %+v", object)
	bytes, err := proto.Marshal(object)
	if err != nil {
		log.Errorf("Failed to create Object %+v: %v", object, err)
		return errors.NewInvalid(err.Error())
	}

	// Put the object in the map using an optimistic lock if this is an update
	entry, err := s.objects.Put(ctx, string(object.ID), bytes, _map.IfNotSet())
	if err != nil {
		err = errors.FromAtomix(err)
		if !errors.IsAlreadyExists(err) {
			log.Errorf("Failed to create Object %+v: %v", object, err)
		} else {
			log.Warnf("Failed to create Object %+v: %v", object, err)
		}
		return err
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

	log.Infof("Updating Object %+v", object)
	bytes, err := proto.Marshal(object)
	if err != nil {
		log.Errorf("Failed to update Object %+v", object, err)
		return errors.NewInvalid(err.Error())
	}

	// Update the object in the map
	entry, err := s.objects.Put(ctx, string(object.ID), bytes, _map.IfMatch(meta.NewRevision(meta.Revision(object.Revision))))
	if err != nil {
		err = errors.FromAtomix(err)
		if !errors.IsNotFound(err) && !errors.IsConflict(err) {
			log.Errorf("Failed to update Object %+v: %v", object, err)
		} else {
			log.Warnf("Failed to update Object %+v: %v", object, err)
		}
		return err
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
		err = errors.FromAtomix(err)
		if !errors.IsNotFound(err) {
			log.Errorf("Failed to get Object '%s': %v", id, err)
		} else {
			log.Warnf("Failed to get Object '%s': %v", id, err)
		}
		return nil, err
	}
	obj, err := decodeObject(*entry)
	if err != nil {
		return nil, err
	}
	s.addSrcTgts(obj)
	return obj, nil
}

func (s *atomixStore) Delete(ctx context.Context, id topoapi.ID, revision topoapi.Revision) error {
	if id == "" {
		return errors.NewInvalid("ID cannot be empty")
	}

	err := s.deleteRelatedRelations(ctx, id)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	log.Infof("Deleting Object '%s'", id)

	if revision == 0 {
		_, err = s.objects.Remove(ctx, string(id))
	} else {
		_, err = s.objects.Remove(ctx, string(id), _map.IfMatch(meta.NewRevision(meta.Revision(revision))))
	}
	if err != nil {
		err = errors.FromAtomix(err)
		if !errors.IsNotFound(err) && !errors.IsConflict(err) {
			log.Errorf("Failed to delete Object '%s': %v", id, err)
		} else {
			log.Warnf("Failed to delete Object '%s': %v", id, err)
		}
		return err
	}
	return nil
}

func (s *atomixStore) deleteRelatedRelations(ctx context.Context, id topoapi.ID) error {
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
						err = errors.FromAtomix(err)
						if !errors.IsNotFound(err) {
							return err
						}
					}
				}
			}
		}
	}
	return nil
}

func (s *atomixStore) List(ctx context.Context, filters *topoapi.Filters) ([]topoapi.Object, error) {
	if filters != nil && filters.RelationFilter != nil {
		return s.listRelationFilter(ctx, filters)
	}

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

	for entry := range mapCh {
		if ep, err := decodeObject(entry); err == nil {
			if match(ep, filters) {
				if matchType(ep, filters.ObjectTypes) && matchAspects(ep, filters.WithAspects) {
					s.addSrcTgts(ep)
					eps = append(eps, *ep)
				}
			}
		}
	}

	return eps, nil
}

func (s *atomixStore) listRelationFilter(ctx context.Context, filters *topoapi.Filters) ([]topoapi.Object, error) {
	filter := filters.RelationFilter

	if len(filter.GetSrcId()) > 0 {
		return s.filterRelationEntities(ctx, topoapi.ID(filter.GetSrcId()), filters, false)
	} else if len(filter.GetTargetId()) > 0 {
		return s.filterRelationEntities(ctx, topoapi.ID(filter.GetTargetId()), filters, true)
	}
	return nil, errors.NewInvalid("filter must contain either srcID or targetID")
}

func (s *atomixStore) filterRelationEntities(ctx context.Context, id topoapi.ID, filters *topoapi.Filters, useSrc bool) ([]topoapi.Object, error) {
	results := make([]topoapi.Object, 0)
	obj, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	rfilter := filters.RelationFilter
	if rfilter.Scope == topoapi.RelationFilterScope_ALL || rfilter.Scope == topoapi.RelationFilterScope_SOURCE_AND_TARGETS {
		results = append(results, *obj)
	}

	relations := obj.GetEntity().SrcRelationIDs
	if useSrc {
		relations = obj.GetEntity().TgtRelationIDs
	}

	for _, rid := range relations {
		robj, err := s.Get(ctx, rid)
		if err == nil && robj.Type == topoapi.Object_RELATION {
			rel := robj.GetRelation()
			if len(rfilter.RelationKind) == 0 || string(rel.KindID) == rfilter.RelationKind {
				oid := rel.GetSrcEntityID()
				if !useSrc {
					oid = rel.GetTgtEntityID()
				}
				ent, err := s.Get(ctx, oid)
				if err == nil && (len(rfilter.TargetKind) == 0 || string(ent.GetEntity().KindID) == rfilter.TargetKind) && matchAspects(ent, filters.WithAspects) {
					if rfilter.Scope == topoapi.RelationFilterScope_ALL ||
						rfilter.Scope == topoapi.RelationFilterScope_RELATIONS_ONLY ||
						rfilter.Scope == topoapi.RelationFilterScope_RELATIONS_AND_TARGETS {
						results = append(results, *robj)
					}

					if rfilter.Scope != topoapi.RelationFilterScope_RELATIONS_ONLY {
						results = append(results, *ent)
					}
				}
			}
		}
	}
	return results, nil
}

func (s *atomixStore) Watch(ctx context.Context, ch chan<- topoapi.Event, filters *topoapi.Filters, opts ...WatchOption) error {
	var watchOpts watchOptions
	for _, opt := range opts {
		opt.apply(&watchOpts)
	}

	// Create separate channels for replay and watch events
	replayCh := make(chan topoapi.Object)
	eventCh := make(chan topoapi.Event)

	// Create a goroutine to first replay existing state to the watcher and then send events
	go func() {
		defer close(ch)

	replayLoop:
		// Process the replay channel first
		for {
			select {
			case object, ok := <-replayCh:
				// If the replay channel is closed, break out of the replay loop
				if !ok {
					break replayLoop
				}
				// If an object is received on the replay channel, write it to
				// the watch channel if it matches the watch filter
				if match(&object, filters) {
					ch <- topoapi.Event{
						Type:   topoapi.EventType_NONE,
						Object: object,
					}
				}
			case <-ctx.Done():
				// If the watch context is closed, drain the replay channel and break out of the replay loop
				go func() {
					for range replayCh {
					}
				}()
				break replayLoop
			}
		}

	eventLoop:
		// Once the replay channel is processed, process the event channel
		for {
			select {
			case event, ok := <-eventCh:
				// If the event channel is closed, break out of the event loop
				if !ok {
					break eventLoop
				}
				// If an event is received on the replay channel, write it to
				// the watch channel if it matches the watch filter
				if match(&event.Object, filters) {
					ch <- event
				}
			case <-ctx.Done():
				// If the watch context is closed, drain the event channel and break out of the event loop
				go func() {
					for range eventCh {
					}
				}()
				break eventLoop
			}
		}
	}()

	// Add the watcher's event channel
	watcherID := uuid.New()
	s.watchersMu.Lock()
	s.watchers[watcherID] = eventCh
	s.watchersMu.Unlock()

	// Get the objects to replay
	var objects []topoapi.Object
	if watchOpts.replay {
		s.cacheMu.RLock()
		objects = make([]topoapi.Object, 0, len(s.cache))
		for _, object := range s.cache {
			objects = append(objects, object)
		}
		s.cacheMu.RUnlock()
	}

	// Replay existing objects in the cache and then close the replay channel
	go func() {
		defer close(replayCh)
		for _, object := range objects {
			replayCh <- object
		}
	}()

	// Remove the watcher and close the event channel once the watch context is done
	go func() {
		<-ctx.Done()
		s.watchersMu.Lock()
		delete(s.watchers, watcherID)
		s.watchersMu.Unlock()
		close(eventCh)
	}()
	return nil
}

func (s *atomixStore) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := s.objects.Close(ctx)
	if err != nil {
		return errors.FromAtomix(err)
	}
	return nil
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
		s.relations.lock.RLock()
		defer s.relations.lock.RUnlock()
		obj.GetEntity().SrcRelationIDs = s.relations.sources[obj.ID]
		obj.GetEntity().TgtRelationIDs = s.relations.targets[obj.ID]
	}
}

// when creating a relation, create the corresponding entries in the store
func (s *atomixStore) registerSrcTgt(obj *topoapi.Object, strict bool) {
	if relation := obj.GetRelation(); relation != nil {
		if strict {
			// check that the connection is valid (src and tgt are in the store). otherwise remove the dangling relation
			if _, err := s.objects.Get(context.Background(), string(relation.SrcEntityID)); err != nil {
				err = errors.FromAtomix(err)
				if errors.IsNotFound(err) {
					_, _ = s.objects.Remove(context.Background(), string(obj.ID))
				} else {
					log.Error(err)
				}
				return
			}
			if _, err := s.objects.Get(context.Background(), string(relation.TgtEntityID)); err != nil {
				err = errors.FromAtomix(err)
				if errors.IsNotFound(err) {
					_, _ = s.objects.Remove(context.Background(), string(obj.ID))
				} else {
					log.Error(err)
				}
				return
			}
		}

		s.relations.lock.Lock()
		defer s.relations.lock.Unlock()
		s.relations.sources[relation.SrcEntityID] = add(s.relations.sources[relation.SrcEntityID], obj.ID)
		s.relations.targets[relation.TgtEntityID] = add(s.relations.targets[relation.TgtEntityID], obj.ID)
	}
}

// when deleting either a relation or entity, remove the corresponding entries in the store
func (s *atomixStore) unregisterSrcTgt(obj *topoapi.Object) {
	if entity := obj.GetEntity(); entity != nil {
		s.relations.lock.Lock()
		defer s.relations.lock.Unlock()
		delete(s.relations.sources, obj.ID)
		delete(s.relations.targets, obj.ID)

	} else if relation := obj.GetRelation(); relation != nil {
		s.relations.lock.Lock()
		defer s.relations.lock.Unlock()
		s.relations.sources[relation.SrcEntityID] = remove(s.relations.sources[relation.SrcEntityID], obj.ID)
		s.relations.targets[relation.TgtEntityID] = remove(s.relations.targets[relation.TgtEntityID], obj.ID)
	}
}

func add(ids []topoapi.ID, id topoapi.ID) []topoapi.ID {
	for _, eid := range ids {
		if eid == id {
			return ids
		}
	}
	return append(ids, id)
}

func remove(ids []topoapi.ID, id topoapi.ID) []topoapi.ID {
	for i, eid := range ids {
		if eid == id {
			ids[i] = ids[len(ids)-1]
			return ids[:len(ids)-1]
		}
	}
	return ids
}
