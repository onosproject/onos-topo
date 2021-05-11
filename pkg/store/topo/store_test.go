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
	"github.com/gogo/protobuf/types"
	topoapi "github.com/onosproject/onos-api/go/onos/topo"
	"github.com/onosproject/onos-lib-go/pkg/atomix"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTopoStore(t *testing.T) {
	_, address := atomix.StartLocalNode()

	store1, err := newLocalStore(address)
	assert.NoError(t, err)
	defer store1.Close()

	store2, err := newLocalStore(address)
	assert.NoError(t, err)
	defer store2.Close()

	ch := make(chan topoapi.Event)
	err = store2.Watch(context.Background(), ch)
	assert.NoError(t, err)

	obj1 := &topoapi.Object{
		ID: "o1",
	}
	obj2 := &topoapi.Object{
		ID: "o2",
	}

	// Create a new object
	err = store1.Create(context.TODO(), obj1)
	assert.NoError(t, err)
	assert.Equal(t, topoapi.ID("o1"), obj1.ID)
	assert.NotEqual(t, topoapi.Revision(0), obj1.Revision)

	// Get the object
	obj1, err = store2.Get(context.TODO(), "o1")
	assert.NoError(t, err)
	assert.NotNil(t, obj1)
	assert.Equal(t, topoapi.ID("o1"), obj1.ID)
	assert.NotEqual(t, topoapi.Revision(0), obj1.Revision)

	// Create another object
	err = store2.Create(context.TODO(), obj2)
	assert.NoError(t, err)
	assert.Equal(t, topoapi.ID("o2"), obj2.ID)
	assert.NotEqual(t, topoapi.Revision(0), obj2.Revision)

	// Verify events were received for the objects
	topoEvent := nextEvent(t, ch)
	assert.Equal(t, topoapi.ID("o1"), topoEvent.ID)
	topoEvent = nextEvent(t, ch)
	assert.Equal(t, topoapi.ID("o2"), topoEvent.ID)

	// Update one of the objects
	obj2.Attributes = make(map[string]*types.Any)
	foo, err := types.MarshalAny(&topoapi.Location{Lat: 1, Lng: 2})
	assert.NoError(t, err)
	obj2.Attributes["foo"] = foo
	revision := obj2.Revision
	err = store1.Update(context.TODO(), obj2)
	assert.NoError(t, err)
	assert.NotEqual(t, revision, obj2.Revision)

	// Read and then update the object
	obj2, err = store2.Get(context.TODO(), "o2")
	assert.NoError(t, err)
	assert.NotNil(t, obj2)
	err = store1.Update(context.TODO(), obj2)
	assert.NoError(t, err)
	assert.NotEqual(t, revision, obj2.Revision)

	// Verify that concurrent updates fail
	obj11, err := store1.Get(context.TODO(), "o1")
	assert.NoError(t, err)
	obj12, err := store2.Get(context.TODO(), "o1")
	assert.NoError(t, err)

	obj11.Attributes = make(map[string]*types.Any)
	bar, err := types.MarshalAny(&topoapi.Location{Lat: 2, Lng: 1})
	assert.NoError(t, err)
	obj11.Attributes["foo"] = bar
	err = store1.Update(context.TODO(), obj11)
	assert.NoError(t, err)

	obj12.Attributes = make(map[string]*types.Any)
	foobar, err := types.MarshalAny(&topoapi.E2Node{})
	assert.NoError(t, err)
	obj12.Attributes["foo"] = foobar
	err = store2.Update(context.TODO(), obj12)
	assert.Error(t, err)

	// Verify events were received again
	topoEvent = nextEvent(t, ch)
	assert.Equal(t, topoapi.ID("o2"), topoEvent.ID)
	topoEvent = nextEvent(t, ch)
	assert.Equal(t, topoapi.ID("o2"), topoEvent.ID)
	topoEvent = nextEvent(t, ch)
	assert.Equal(t, topoapi.ID("o1"), topoEvent.ID)

	// Verify the attribute values
	obj2g, err := store1.Get(context.TODO(), obj2.ID)
	assert.NoError(t, err)
	var loc topoapi.Location
	err = types.UnmarshalAny(obj2g.Attributes["foo"], &loc)
	assert.NoError(t, err)
	assert.Equal(t, 1.0, loc.Lat)
	assert.Equal(t, 2.0, loc.Lng)

	// List the objects
	objects, err := store1.List(context.TODO())
	assert.NoError(t, err)
	assert.Len(t, objects, 2)

	// Delete an object
	err = store1.Delete(context.TODO(), obj2.ID)
	assert.NoError(t, err)
	obj2, err = store2.Get(context.TODO(), "o2")
	assert.Error(t, err)
	assert.True(t, errors.IsNotFound(err))
	assert.Nil(t, obj2)

	obj := &topoapi.Object{
		ID: "o1",
	}

	err = store1.Create(context.TODO(), obj)
	assert.Error(t, err)

	obj = &topoapi.Object{
		ID: "o2",
	}

	err = store1.Create(context.TODO(), obj)
	assert.NoError(t, err)

	ch = make(chan topoapi.Event)
	err = store1.Watch(context.TODO(), ch, WithReplay())
	assert.NoError(t, err)

	obj = nextEvent(t, ch)
	assert.NotNil(t, obj)
	obj = nextEvent(t, ch)
	assert.NotNil(t, obj)
}

func nextEvent(t *testing.T, ch chan topoapi.Event) *topoapi.Object {
	select {
	case c := <-ch:
		return &c.Object
	case <-time.After(5 * time.Second):
		t.FailNow()
	}
	return nil
}
