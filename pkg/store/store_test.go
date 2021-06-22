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
	"testing"
	"time"

	"github.com/atomix/atomix-go-client/pkg/atomix/test"
	"github.com/atomix/atomix-go-client/pkg/atomix/test/rsm"
	topoapi "github.com/onosproject/onos-api/go/onos/topo"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestTopoStore(t *testing.T) {
	test := test.NewTest(
		rsm.NewProtocol(),
		test.WithReplicas(1),
		test.WithPartitions(1))
	assert.NoError(t, test.Start())
	defer test.Stop()

	client1, err := test.NewClient("node-1")
	assert.NoError(t, err)

	client2, err := test.NewClient("node-2")
	assert.NoError(t, err)

	store1, err := NewAtomixStore(client1)
	assert.NoError(t, err)

	store2, err := NewAtomixStore(client2)
	assert.NoError(t, err)

	ch := make(chan topoapi.Event)
	err = store2.Watch(context.Background(), ch, nil)
	assert.NoError(t, err)

	k1 := &topoapi.Object{
		ID:   "foo",
		Type: topoapi.Object_KIND,
		Obj:  &topoapi.Object_Kind{Kind: &topoapi.Kind{Name: "Foo"}},
	}
	err = store1.Create(context.TODO(), k1)
	assert.NoError(t, err)

	k2 := &topoapi.Object{
		ID:   "bar",
		Type: topoapi.Object_KIND,
		Obj:  &topoapi.Object_Kind{Kind: &topoapi.Kind{Name: "Bar"}},
	}
	err = store1.Create(context.TODO(), k2)
	assert.NoError(t, err)

	obj1 := &topoapi.Object{
		ID:     "o1",
		Type:   topoapi.Object_ENTITY,
		Obj:    &topoapi.Object_Entity{Entity: &topoapi.Entity{KindID: topoapi.ID("foo")}},
		Labels: map[string]string{},
	}
	obj1.Labels["env"] = "test"
	obj1.Labels["area"] = "ran"

	obj2 := &topoapi.Object{
		ID:     "o2",
		Type:   topoapi.Object_ENTITY,
		Obj:    &topoapi.Object_Entity{Entity: &topoapi.Entity{KindID: topoapi.ID("bar")}},
		Labels: map[string]string{},
	}
	obj2.Labels["env"] = "production"
	obj2.Labels["area"] = "ran"

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
	assert.Equal(t, "test", obj1.Labels["env"])
	assert.Equal(t, "ran", obj1.Labels["area"])

	// Create another object
	err = store2.Create(context.TODO(), obj2)
	assert.NoError(t, err)
	assert.Equal(t, topoapi.ID("o2"), obj2.ID)
	assert.NotEqual(t, topoapi.Revision(0), obj2.Revision)

	// Verify events were received for the kinds
	topoEvent := nextEvent(t, ch)
	assert.Equal(t, topoapi.ID("foo"), topoEvent.ID)
	topoEvent = nextEvent(t, ch)
	assert.Equal(t, topoapi.ID("bar"), topoEvent.ID)

	// Verify events were received for the objects
	topoEvent = nextEvent(t, ch)
	assert.Equal(t, topoapi.ID("o1"), topoEvent.ID)
	topoEvent = nextEvent(t, ch)
	assert.Equal(t, topoapi.ID("o2"), topoEvent.ID)

	// Update one of the objects
	err = obj2.SetAspect(&topoapi.Location{Lat: 1, Lng: 2})
	assert.NoError(t, err)
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
	assert.Equal(t, "production", obj2.Labels["env"])
	assert.Equal(t, "ran", obj2.Labels["area"])

	// Verify that concurrent updates fail
	obj11, err := store1.Get(context.TODO(), "o1")
	assert.NoError(t, err)
	obj12, err := store2.Get(context.TODO(), "o1")
	assert.NoError(t, err)

	err = obj11.SetAspect(&topoapi.Location{Lat: 2, Lng: 1})
	assert.NoError(t, err)
	err = store1.Update(context.TODO(), obj11)
	assert.NoError(t, err)

	err = obj12.SetAspect(&topoapi.E2Node{})
	assert.NoError(t, err)
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
	loc := obj2g.GetAspect(&topoapi.Location{}).(*topoapi.Location)
	assert.NotNil(t, loc)
	assert.Equal(t, 1.0, loc.Lat)
	assert.Equal(t, 2.0, loc.Lng)

	// List the objects
	objects, err := store1.List(context.TODO(), nil)
	assert.NoError(t, err)
	assert.Len(t, objects, 4)

	// List the objects with label filter
	objects, err = store1.List(context.TODO(), &topoapi.Filters{LabelFilters: []*topoapi.Filter{
		{
			Filter: &topoapi.Filter_Equal_{
				Equal_: &topoapi.EqualFilter{Value: "production"},
			},
			Key: "env",
		},
	}})
	assert.NoError(t, err)
	assert.Len(t, objects, 1)
	assert.Equal(t, "o2", string(objects[0].ID))

	// List the objects with kind filter
	objects, err = store1.List(context.TODO(), &topoapi.Filters{KindFilters: []*topoapi.Filter{
		{
			Filter: &topoapi.Filter_Not{
				Not: &topoapi.NotFilter{
					Inner: &topoapi.Filter{
						Filter: &topoapi.Filter_Equal_{
							Equal_: &topoapi.EqualFilter{Value: "bar"},
						},
					},
				},
			},
		},
	}})
	assert.NoError(t, err)
	assert.Len(t, objects, 1)
	assert.Equal(t, "o1", string(objects[0].ID))

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
		ID:   "o2",
		Type: topoapi.Object_ENTITY,
	}

	err = store1.Create(context.TODO(), obj)
	assert.NoError(t, err)

	ch = make(chan topoapi.Event)
	err = store1.Watch(context.TODO(), ch, nil, WithReplay())
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

func TestList(t *testing.T) {
	test := test.NewTest(
		rsm.NewProtocol(),
		test.WithReplicas(1),
		test.WithPartitions(1))
	assert.NoError(t, test.Start())
	defer test.Stop()
	// Define client, store, and objects

	// Client def
	client, _ := test.NewClient("client")

	// Store def
	store, _ := NewAtomixStore(client)

	// Objects def:
	// - node 1234
	// - node 2001
	// - cell 87893172902461441
	// - cell 87893172902461443
	// - cell 87893172902445057
	// - cell 87893172902445058
	// - cell 87893172902445059
	// - cell 87893172902445060
	// - cell-neighbor: 87893172902461441, 87893172902461443 + vice versa
	// - cell-neighbor: 87893172902445057, 87893172902445058 + vice versa
	// - cell-neighbor: 87893172902445058, 87893172902445059 + vice versa
	// - cell-neighbor: 87893172902445059, 87893172902445060 + vice versa
	// - node-cell: 1234, 87893172902461441
	// - node-cell: 1234, 87893172902461443
	// - node-cell: 2001, 87893172902445057
	// - node-cell: 2001, 87893172902445058
	// - node-cell: 2001, 87893172902445059
	// - node-cell: 2001, 87893172902445060
	createObjectsListTest(store)

	// List the objects
	objects, err := store.List(context.TODO(), nil)
	assert.NoError(t, err)
	assert.Len(t, objects, 22)

	// List the objects with label filter
	objects, err = store.List(context.TODO(), &topoapi.Filters{LabelFilters: []*topoapi.Filter{
		{
			Filter: &topoapi.Filter_Equal_{
				Equal_: &topoapi.EqualFilter{Value: "production"},
			},
			Key: "env",
		},
	}})
	assert.NoError(t, err)
	assert.Len(t, objects, 3)
	// assert.Equal(t, "1234", string(objects[0].ID))

	// List the objects with kind filter
	objects, err = store.List(context.TODO(), &topoapi.Filters{KindFilters: []*topoapi.Filter{
		{
			Filter: &topoapi.Filter_Not{
				Not: &topoapi.NotFilter{
					Inner: &topoapi.Filter{
						Filter: &topoapi.Filter_Equal_{
							Equal_: &topoapi.EqualFilter{Value: "e2-cell"},
						},
					},
				},
			},
		},
	}})
	assert.NoError(t, err)
	assert.Len(t, objects, 16)
	// assert.Equal(t, "1234", string(objects[0].ID))

	// List the objects with relation filter
	objects, err = store.List(context.TODO(), &topoapi.Filters{
		RelationFilter: &topoapi.RelationFilter{SrcId: "1234", RelationKind: "e2-node-cell", TargetKind: ""},
	})
	assert.NoError(t, err)
	assert.Len(t, objects, 2)
	// assert.Equal(t, "1234-87893172902461441", string(objects[0].ID))

	objects, err = store.List(context.TODO(), &topoapi.Filters{
		RelationFilter: &topoapi.RelationFilter{SrcId: "87893172902445058", RelationKind: "e2-cell-neighbor", TargetKind: ""},
	})
	assert.NoError(t, err)
	assert.Len(t, objects, 2)
	// assert.Equal(t, "87893172902445058-87893172902445057", string(objects[0].ID))
	// assert.Equal(t, "87893172902445058-87893172902445059", string(objects[1].ID))

	// No test for relation filter with target kind: cell-neighbor, node-cell do not have different target kinds
}

func createObjectsListTest(s Store) {
	_ = s.Create(context.TODO(), &topoapi.Object{
		ID:     "1234",
		Type:   topoapi.Object_ENTITY,
		Obj:    &topoapi.Object_Entity{Entity: &topoapi.Entity{KindID: topoapi.ID("e2-node")}},
		Labels: map[string]string{"env": "production"},
	})
	_ = s.Create(context.TODO(), &topoapi.Object{
		ID:     "2001",
		Type:   topoapi.Object_ENTITY,
		Obj:    &topoapi.Object_Entity{Entity: &topoapi.Entity{KindID: topoapi.ID("e2-node")}},
		Labels: map[string]string{"env": "production"},
	})
	_ = s.Create(context.TODO(), &topoapi.Object{
		ID:     "87893172902461441",
		Type:   topoapi.Object_ENTITY,
		Obj:    &topoapi.Object_Entity{Entity: &topoapi.Entity{KindID: topoapi.ID("e2-cell")}},
		Labels: map[string]string{"env": "production"},
	})
	_ = s.Create(context.TODO(), &topoapi.Object{
		ID:     "87893172902461443",
		Type:   topoapi.Object_ENTITY,
		Obj:    &topoapi.Object_Entity{Entity: &topoapi.Entity{KindID: topoapi.ID("e2-cell")}},
		Labels: map[string]string{"env": "dev"},
	})
	_ = s.Create(context.TODO(), &topoapi.Object{
		ID:     "87893172902445057",
		Type:   topoapi.Object_ENTITY,
		Obj:    &topoapi.Object_Entity{Entity: &topoapi.Entity{KindID: topoapi.ID("e2-cell")}},
		Labels: map[string]string{"env": "dev"},
	})
	_ = s.Create(context.TODO(), &topoapi.Object{
		ID:     "87893172902445058",
		Type:   topoapi.Object_ENTITY,
		Obj:    &topoapi.Object_Entity{Entity: &topoapi.Entity{KindID: topoapi.ID("e2-cell")}},
		Labels: map[string]string{},
	})
	_ = s.Create(context.TODO(), &topoapi.Object{
		ID:     "87893172902445059",
		Type:   topoapi.Object_ENTITY,
		Obj:    &topoapi.Object_Entity{Entity: &topoapi.Entity{KindID: topoapi.ID("e2-cell")}},
		Labels: map[string]string{},
	})
	_ = s.Create(context.TODO(), &topoapi.Object{
		ID:     "87893172902445060",
		Type:   topoapi.Object_ENTITY,
		Obj:    &topoapi.Object_Entity{Entity: &topoapi.Entity{KindID: topoapi.ID("e2-cell")}},
		Labels: map[string]string{},
	})
	_ = s.Create(context.TODO(), &topoapi.Object{
		ID:     "87893172902461441-87893172902461443",
		Type:   topoapi.Object_RELATION,
		Obj:    &topoapi.Object_Relation{Relation: &topoapi.Relation{KindID: "e2-cell-neighbor", SrcEntityID: "87893172902461441", TgtEntityID: "87893172902461443"}},
		Labels: map[string]string{},
	})
	_ = s.Create(context.TODO(), &topoapi.Object{
		ID:     "87893172902445057-87893172902445058",
		Type:   topoapi.Object_RELATION,
		Obj:    &topoapi.Object_Relation{Relation: &topoapi.Relation{KindID: "e2-cell-neighbor", SrcEntityID: "87893172902445057", TgtEntityID: "87893172902445058"}},
		Labels: map[string]string{},
	})
	_ = s.Create(context.TODO(), &topoapi.Object{
		ID:     "87893172902445058-87893172902445059",
		Type:   topoapi.Object_RELATION,
		Obj:    &topoapi.Object_Relation{Relation: &topoapi.Relation{KindID: "e2-cell-neighbor", SrcEntityID: "87893172902445058", TgtEntityID: "87893172902445059"}},
		Labels: map[string]string{},
	})
	_ = s.Create(context.TODO(), &topoapi.Object{
		ID:     "87893172902445059-87893172902445060",
		Type:   topoapi.Object_RELATION,
		Obj:    &topoapi.Object_Relation{Relation: &topoapi.Relation{KindID: "e2-cell-neighbor", SrcEntityID: "87893172902445059", TgtEntityID: "87893172902445060"}},
		Labels: map[string]string{},
	})
	_ = s.Create(context.TODO(), &topoapi.Object{
		ID:     "87893172902461443-87893172902461441",
		Type:   topoapi.Object_RELATION,
		Obj:    &topoapi.Object_Relation{Relation: &topoapi.Relation{KindID: "e2-cell-neighbor", SrcEntityID: "87893172902461443", TgtEntityID: "87893172902461441"}},
		Labels: map[string]string{},
	})
	_ = s.Create(context.TODO(), &topoapi.Object{
		ID:     "87893172902445058-87893172902445057",
		Type:   topoapi.Object_RELATION,
		Obj:    &topoapi.Object_Relation{Relation: &topoapi.Relation{KindID: "e2-cell-neighbor", SrcEntityID: "87893172902445058", TgtEntityID: "87893172902445057"}},
		Labels: map[string]string{},
	})
	_ = s.Create(context.TODO(), &topoapi.Object{
		ID:     "87893172902445059-87893172902445058",
		Type:   topoapi.Object_RELATION,
		Obj:    &topoapi.Object_Relation{Relation: &topoapi.Relation{KindID: "e2-cell-neighbor", SrcEntityID: "87893172902445059", TgtEntityID: "87893172902445058"}},
		Labels: map[string]string{},
	})
	_ = s.Create(context.TODO(), &topoapi.Object{
		ID:     "87893172902445060-87893172902445059",
		Type:   topoapi.Object_RELATION,
		Obj:    &topoapi.Object_Relation{Relation: &topoapi.Relation{KindID: "e2-cell-neighbor", SrcEntityID: "87893172902445060", TgtEntityID: "87893172902445059"}},
		Labels: map[string]string{},
	})
	_ = s.Create(context.TODO(), &topoapi.Object{
		ID:     "1234-87893172902461441",
		Type:   topoapi.Object_RELATION,
		Obj:    &topoapi.Object_Relation{Relation: &topoapi.Relation{KindID: "e2-node-cell", SrcEntityID: "1234", TgtEntityID: "87893172902461441"}},
		Labels: map[string]string{},
	})
	_ = s.Create(context.TODO(), &topoapi.Object{
		ID:     "1234-87893172902461443",
		Type:   topoapi.Object_RELATION,
		Obj:    &topoapi.Object_Relation{Relation: &topoapi.Relation{KindID: "e2-node-cell", SrcEntityID: "1234", TgtEntityID: "87893172902461443"}},
		Labels: map[string]string{},
	})
	_ = s.Create(context.TODO(), &topoapi.Object{
		ID:     "2001-87893172902445057",
		Type:   topoapi.Object_RELATION,
		Obj:    &topoapi.Object_Relation{Relation: &topoapi.Relation{KindID: "e2-node-cell", SrcEntityID: "2001", TgtEntityID: "87893172902445057"}},
		Labels: map[string]string{},
	})
	_ = s.Create(context.TODO(), &topoapi.Object{
		ID:     "2001-87893172902445058",
		Type:   topoapi.Object_RELATION,
		Obj:    &topoapi.Object_Relation{Relation: &topoapi.Relation{KindID: "e2-node-cell", SrcEntityID: "2001", TgtEntityID: "87893172902445058"}},
		Labels: map[string]string{},
	})
	_ = s.Create(context.TODO(), &topoapi.Object{
		ID:     "2001-87893172902445059",
		Type:   topoapi.Object_RELATION,
		Obj:    &topoapi.Object_Relation{Relation: &topoapi.Relation{KindID: "e2-node-cell", SrcEntityID: "2001", TgtEntityID: "87893172902445059"}},
		Labels: map[string]string{},
	})
	_ = s.Create(context.TODO(), &topoapi.Object{
		ID:     "2001-87893172902445060",
		Type:   topoapi.Object_RELATION,
		Obj:    &topoapi.Object_Relation{Relation: &topoapi.Relation{KindID: "e2-node-cell", SrcEntityID: "2001", TgtEntityID: "87893172902445060"}},
		Labels: map[string]string{},
	})
}
