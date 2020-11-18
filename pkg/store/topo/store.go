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
	"io"
	"time"

	_map "github.com/atomix/go-client/pkg/client/map"
	"github.com/atomix/go-client/pkg/client/primitive"
	"github.com/gogo/protobuf/proto"
	"github.com/onosproject/onos-lib-go/pkg/atomix"
	topoapi "github.com/onosproject/onos-topo/api/topo"
	"github.com/onosproject/onos-topo/pkg/config"
)

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
	node, address := atomix.StartLocalNode()
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
		closer:  node.Stop,
	}, nil
}

// Store stores topology information
type Store interface {
	io.Closer

	// Load loads a object from the store
	Load(objectID topoapi.ID) (*topoapi.Object, error)

	// Store stores a object in the store
	Store(*topoapi.Object) error

	// Delete deletes a object from the store
	Delete(topoapi.ID) error

	// List streams objects to the given channel
	List(chan<- *topoapi.Object) error

	// Watch streams object events to the given channel
	Watch(chan<- *Event, ...WatchOption) error
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
	closer  func() error
}

func (s *atomixStore) Load(objectID topoapi.ID) (*topoapi.Object, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	entry, err := s.objects.Get(ctx, string(objectID))
	if err != nil {
		return nil, err
	} else if entry == nil {
		return nil, nil
	}
	return decodeObject(entry)
}

func (s *atomixStore) Store(object *topoapi.Object) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	bytes, err := proto.Marshal(object)
	if err != nil {
		return err
	}

	// Put the object in the map using an optimistic lock if this is an update
	_, err = s.objects.Put(ctx, string(object.ID), bytes)

	if err != nil {
		return err
	}

	return err
}

func (s *atomixStore) Delete(id topoapi.ID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := s.objects.Remove(ctx, string(id))
	return err
}

func (s *atomixStore) List(ch chan<- *topoapi.Object) error {
	mapCh := make(chan *_map.Entry)
	if err := s.objects.Entries(context.Background(), mapCh); err != nil {
		return err
	}

	go func() {
		defer close(ch)
		for entry := range mapCh {
			if object, err := decodeObject(entry); err == nil {
				ch <- object
			}
		}
	}()
	return nil
}

func (s *atomixStore) Watch(ch chan<- *Event, opts ...WatchOption) error {
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
				ch <- &Event{
					Type:   EventType(event.Type),
					Object: object,
				}
			}
		}
	}()
	return nil
}

func (s *atomixStore) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	_ = s.objects.Close(ctx)
	cancel()
	if s.closer != nil {
		return s.closer()
	}
	return nil
}

func decodeObject(entry *_map.Entry) (*topoapi.Object, error) {
	object := &topoapi.Object{}
	if err := proto.Unmarshal(entry.Value, object); err != nil {
		return nil, err
	}
	object.ID = topoapi.ID(entry.Key)
	return object, nil
}

// EventType provides the type for a object event
type EventType string

const (
	// EventNone is no event
	EventNone EventType = ""
	// EventInserted is inserted
	EventInserted EventType = "inserted"
	// EventUpdated is updated
	EventUpdated EventType = "updated"
	// EventRemoved is removed
	EventRemoved EventType = "removed"
)

// Event is a store event for a object
type Event struct {
	Type   EventType
	Object *topoapi.Object
}
