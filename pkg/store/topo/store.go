// Copyright 2019-present Open Networking Foundation.
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

package topo

import (
	"context"
	"github.com/atomix/go-client/pkg/client/util/net"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"io"
	"time"

	_map "github.com/atomix/go-client/pkg/client/map"
	"github.com/atomix/go-client/pkg/client/primitive"
	"github.com/gogo/protobuf/proto"
	topoapi "github.com/onosproject/onos-api/go/onos/topo"
	"github.com/onosproject/onos-lib-go/pkg/atomix"
	"github.com/onosproject/onos-topo/pkg/config"
)

var log = logging.GetLogger("store", "topo")

// NewAtomixStore returns a new persistent Store
func NewAtomixStore() (Store, error) {
	ricConfig, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	database, err := atomix.GetDatabase(ricConfig.Atomix, ricConfig.Atomix.GetDatabase(atomix.DatabaseTypeConsensus))
	if err != nil {
		return nil, err
	}

	objects, err := database.GetMap(context.Background(), "objects")
	if err != nil {
		return nil, err
	}

	return &atomixStore{
		objects: objects,
	}, nil
}

// NewLocalStore returns a new local object store
func NewLocalStore() (Store, error) {
	_, address := atomix.StartLocalNode()
	return newLocalStore(address)
}

func newLocalStore(address net.Address) (Store, error) {
	name := primitive.Name{
		Namespace: "local",
		Name:      "objects",
	}

	session, err := primitive.NewSession(context.TODO(), primitive.Partition{ID: 1, Address: address})
	if err != nil {
		return nil, err
	}

	objects, err := _map.New(context.Background(), name, []*primitive.Session{session})
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
	List(ctx context.Context) ([]topoapi.Object, error)

	// Watch streams object events to the given channel
	Watch(ctx context.Context, ch chan<- topoapi.Event, opts ...WatchOption) error
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

	log.Infof("Creating object %+v", object)
	bytes, err := proto.Marshal(object)
	if err != nil {
		log.Errorf("Failed to create object %+v: %s", object, err)
		return err
	}

	// Put the object in the map using an optimistic lock if this is an update
	entry, err := s.objects.Put(ctx, string(object.ID), bytes, _map.IfNotSet())
	if err != nil {
		log.Errorf("Failed to create object %+v: %s", object, err)
		return err
	}

	object.Revision = topoapi.Revision(entry.Version)
	return err
}

func (s *atomixStore) Update(ctx context.Context, object *topoapi.Object) error {
	if object.ID == "" {
		return errors.NewInvalid("ID cannot be empty")
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
	entry, err := s.objects.Put(ctx, string(object.ID), bytes, _map.IfVersion(_map.Version(object.Revision)))
	if err != nil {
		log.Errorf("Failed to update object %+v: %s", object, err)
		return errors.FromAtomix(err)
	}
	object.Revision = topoapi.Revision(entry.Version)
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
	return decodeObject(entry)
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

func (s *atomixStore) List(ctx context.Context) ([]topoapi.Object, error) {
	mapCh := make(chan *_map.Entry)
	if err := s.objects.Entries(ctx, mapCh); err != nil {
		return nil, err
	}

	eps := make([]topoapi.Object, 0)

	for entry := range mapCh {
		if ep, err := decodeObject(entry); err == nil {
			eps = append(eps, *ep)
		}
	}
	return eps, nil
}

func (s *atomixStore) Watch(ctx context.Context, ch chan<- topoapi.Event, opts ...WatchOption) error {
	watchOpts := make([]_map.WatchOption, 0)
	for _, opt := range opts {
		watchOpts = opt.apply(watchOpts)
	}

	mapCh := make(chan *_map.Event)
	if err := s.objects.Watch(context.Background(), mapCh, watchOpts...); err != nil {
		return err
	}

	go func() {
		defer close(ch)
		for event := range mapCh {
			if object, err := decodeObject(event.Entry); err == nil {
				var eventType topoapi.EventType
				switch event.Type {
				case _map.EventNone:
					eventType = topoapi.EventType_NONE
				case _map.EventInserted:
					eventType = topoapi.EventType_ADDED
				case _map.EventUpdated:
					eventType = topoapi.EventType_UPDATED
				case _map.EventRemoved:
					eventType = topoapi.EventType_REMOVED
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

func decodeObject(entry *_map.Entry) (*topoapi.Object, error) {
	object := &topoapi.Object{}
	if err := proto.Unmarshal(entry.Value, object); err != nil {
		return nil, err
	}
	object.ID = topoapi.ID(entry.Key)
	object.Revision = topoapi.Revision(entry.Version)
	return object, nil
}
