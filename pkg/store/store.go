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
	return &atomixStore{
		objects: objects,
	}, nil
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

func (o watchReplayOption) apply(opts []_map.WatchOption) []_map.WatchOption {
	return append(opts, _map.WithReplay())
}

// WithReplay returns a WatchOption that replays past changes
func WithReplay() WatchOption {
	return watchReplayOption{}
}

// atomixStore is the object implementation of the Store
type atomixStore struct {
	objects _map.Map
}

func (s *atomixStore) Create(ctx context.Context, object *topoapi.Object) error {
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
	return decodeObject(*entry)
}

func (s *atomixStore) Delete(ctx context.Context, id topoapi.ID) error {
	if id == "" {
		return errors.NewInvalid("ID cannot be empty")
	}

	log.Infof("Deleting object %s", id)
	_, err := s.objects.Remove(ctx, string(id))
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

	// first make sure there are filters. if there aren't, just return everything
	if filters == nil {
		for entry := range mapCh {
			if ep, err := decodeObject(entry); err == nil {
				eps = append(eps, *ep)
			}
		}
		return eps, nil
	}

	if filters.RelationFilter != nil {
		filter := filters.RelationFilter

		// contains _all_ relations that have the same kind as the filter, same SrcId as the filter, and has a target ID that (at time of their reading in mapCh) had not been seen
		entitiesToGet := make(map[topoapi.ID]*topoapi.Object)

		for entry := range mapCh {
			if ep, err := decodeObject(entry); err == nil {
				// if object is a relation and its kind and src id matches the filter, push the destination id
				if ep.Type == topoapi.Object_RELATION && string(ep.GetRelation().KindID) == filter.GetRelationKind() && string(ep.GetRelation().GetSrcEntityID()) == filter.SrcId {
					entitiesToGet[ep.GetRelation().TgtEntityID] = nil
				} else
				// if object is an entity, see if satisfies some relation. else, put in unresolved_entities
				if ep.Type == topoapi.Object_ENTITY {
					if map_entity, found := entitiesToGet[ep.ID]; found {
						if filter.TargetKind == "" || ep.GetKind().Name == filter.TargetKind {
							*map_entity = *ep
						}
					}
				}
			}
		}
		// iterate over entities to make sure we did not miss any valid ones (due to the corresponding relationship being seen first)
		for id := range entitiesToGet {
			if entitiesToGet[id] == nil {
				entity, _ := s.Get(ctx, id)
				if filter.TargetKind == "" || entity.GetKind().Name == filter.TargetKind {
					eps = append(eps, *entity)
				}
			} else {
				eps = append(eps, *entitiesToGet[id])
			}
		}
	} else {
		for entry := range mapCh {
			if ep, err := decodeObject(entry); err == nil {
				if match(ep, filters) {
					eps = append(eps, *ep)
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
