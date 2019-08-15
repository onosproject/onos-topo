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

package device

import (
	"context"
	"errors"
	"github.com/atomix/atomix-go-client/pkg/client/map_"
	"github.com/atomix/atomix-go-client/pkg/client/session"
	"github.com/gogo/protobuf/proto"
	"github.com/onosproject/onos-topo/pkg/util"
	"sync"
	"time"
)

// NewAtomixStore returns a new persistent Store
func NewAtomixStore() (Store, error) {
	client, err := util.GetAtomixClient()
	if err != nil {
		return nil, err
	}

	group, err := client.GetGroup(context.Background(), util.GetAtomixRaftGroup())
	if err != nil {
		return nil, err
	}

	devices, err := group.GetMap(context.Background(), "devices", session.WithTimeout(30*time.Second))
	if err != nil {
		return nil, err
	}

	return &atomixStore{
		devices: devices,
	}, nil
}

// NewLocalStore returns a new local device store
func NewLocalStore() Store {
	return &localStore{
		devices:  make(map[ID]Device),
		watchers: make([]chan<- *Event, 0),
	}
}

// Store stores topology information
type Store interface {
	// Load loads a device from the store
	Load(deviceID ID) (*Device, error)

	// Store stores a device in the store
	Store(*Device) error

	// Delete deletes a device from the store
	Delete(*Device) error

	// List streams devices to the given channel
	List(chan<- *Device) error

	// Watch streams device events to the given channel
	Watch(chan<- *Event) error
}

// localStore is a local implementation of the device Store
type localStore struct {
	devices  map[ID]Device
	mu       sync.RWMutex
	revision uint64
	watchers []chan<- *Event
}

func (s *localStore) Load(deviceID ID) (*Device, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	device, ok := s.devices[deviceID]
	if !ok {
		return nil, nil
	}
	return &device, nil
}

func (s *localStore) Store(device *Device) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if device.Revision == 0 {
		_, ok := s.devices[device.ID]
		if !ok {
			s.revision++
			device.Revision = Revision(s.revision)
			s.devices[device.ID] = *device
			s.broadcastEvent(&Event{
				Type:   EventInserted,
				Device: device,
			})
		} else {
			return errors.New("device already exists")
		}
	} else {
		storedDevice, ok := s.devices[device.ID]
		if ok && device.Revision == storedDevice.Revision {
			s.revision++
			device.Revision = Revision(s.revision)
			s.devices[device.ID] = *device
			s.broadcastEvent(&Event{
				Type:   EventUpdated,
				Device: device,
			})
		} else {
			return errors.New("unknown device")
		}
	}
	return nil
}

func (s *localStore) Delete(device *Device) error {
	if device.Revision == 0 {
		return errors.New("no device revision provided")
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	storedDevice, ok := s.devices[device.ID]
	if ok && storedDevice.Revision == device.Revision {
		delete(s.devices, device.ID)
		s.broadcastEvent(&Event{
			Type:   EventRemoved,
			Device: &storedDevice,
		})
		return nil
	}
	return errors.New("device out of date")
}

func (s *localStore) List(ch chan<- *Device) error {
	go func() {
		s.mu.RLock()
		devices := make([]*Device, 0, len(s.devices))
		for _, device := range s.devices {
			devices = append(devices, &device)
		}
		s.mu.RUnlock()

		defer close(ch)
		for _, device := range devices {
			ch <- device
		}
	}()
	return nil
}

func (s *localStore) Watch(ch chan<- *Event) error {
	go func() {
		s.mu.Lock()
		s.watchers = append(s.watchers, ch)
		devices := make([]*Device, 0, len(s.devices))
		for _, device := range s.devices {
			devices = append(devices, &device)
		}
		defer s.mu.Unlock()

		for _, device := range devices {
			ch <- &Event{
				Type:   EventNone,
				Device: device,
			}
		}
	}()
	return nil
}

func (s *localStore) broadcastEvent(event *Event) {
	for _, watcher := range s.watchers {
		watcher <- event
	}
}

// atomixStore is the device implementation of the Store
type atomixStore struct {
	devices map_.Map
}

func (s *atomixStore) Load(deviceID ID) (*Device, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	kv, err := s.devices.Get(ctx, string(deviceID))
	if err != nil {
		return nil, err
	}
	return decodeDevice(kv.Key, kv.Value, kv.Version)
}

func (s *atomixStore) Store(device *Device) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	bytes, err := proto.Marshal(device)
	if err != nil {
		return err
	}

	// Put the device in the map using an optimistic lock if this is an update
	var kv *map_.KeyValue
	if device.Revision == 0 {
		kv, err = s.devices.Put(ctx, string(device.ID), bytes)
	} else {
		kv, err = s.devices.Put(ctx, string(device.ID), bytes, map_.WithVersion(int64(device.Revision)))
	}

	if err != nil {
		return err
	}

	// Update the device metadata
	device.Revision = Revision(kv.Version)
	return err
}

func (s *atomixStore) Delete(device *Device) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if device.Revision > 0 {
		_, err := s.devices.Remove(ctx, string(device.ID), map_.WithVersion(int64(device.Revision)))
		return err
	}
	_, err := s.devices.Remove(ctx, string(device.ID))
	return err
}

func (s *atomixStore) List(ch chan<- *Device) error {
	mapCh := make(chan *map_.KeyValue)
	if err := s.devices.Entries(context.Background(), mapCh); err != nil {
		return err
	}

	go func() {
		defer close(ch)
		for kv := range mapCh {
			if device, err := decodeDevice(kv.Key, kv.Value, kv.Version); err == nil {
				ch <- device
			}
		}
	}()
	return nil
}

func (s *atomixStore) Watch(ch chan<- *Event) error {
	mapCh := make(chan *map_.MapEvent)
	if err := s.devices.Watch(context.Background(), mapCh, map_.WithReplay()); err != nil {
		return err
	}

	go func() {
		defer close(ch)
		for event := range mapCh {
			if device, err := decodeDevice(event.Key, event.Value, event.Version); err == nil {
				ch <- &Event{
					Type:   EventType(event.Type),
					Device: device,
				}
			}
		}
	}()
	return nil
}

func decodeDevice(key string, value []byte, version int64) (*Device, error) {
	device := &Device{}
	if err := proto.Unmarshal(value, device); err != nil {
		return nil, err
	}
	device.ID = ID(key)
	device.Revision = Revision(version)
	return device, nil
}

// EventType provides the type for a device event
type EventType string

const (
	EventNone     EventType = ""
	EventInserted EventType = "inserted"
	EventUpdated  EventType = "updated"
	EventRemoved  EventType = "removed"
)

// Event is a store event for a device
type Event struct {
	Type   EventType
	Device *Device
}
