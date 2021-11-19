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
	"github.com/onosproject/onos-api/go/onos/topo"
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

	// List the objects; there should be none
	noobjects, err := store1.List(context.TODO(), nil)
	assert.NoError(t, err)
	assert.Len(t, noobjects, 0)

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
	loc := &topoapi.Location{}
	err = obj2g.GetAspect(loc)
	assert.NoError(t, err)
	assert.NotNil(t, loc)
	assert.Equal(t, 1.0, loc.Lat)
	assert.Equal(t, 2.0, loc.Lng)

	// Delete an object; with wrong rev then with right rev
	err = store1.Delete(context.TODO(), obj2.ID, obj2.Revision+10)
	assert.Error(t, err)
	err = store1.Delete(context.TODO(), obj2.ID, obj2.Revision)
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
	createObjectsListTest(t, store)

	object, err := store.Get(context.TODO(), "1234")
	assert.NoError(t, err)
	assert.Len(t, object.GetEntity().SrcRelationIDs, 2)
	assert.Len(t, object.GetEntity().TgtRelationIDs, 0)

	// List the objects
	objects, err := store.List(context.TODO(), nil)
	assert.NoError(t, err)
	assert.Len(t, objects, 22) // 2 nodes + 6 cells + 8 cell-neighbors + 6 node-cells

	// List the objects
	objects, err = store.List(context.TODO(), &topoapi.Filters{ObjectTypes: []topoapi.Object_Type{topoapi.Object_ENTITY}})
	assert.NoError(t, err)
	assert.Len(t, objects, 8) // 2 nodes + 6 cells

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
	assert.Len(t, objects, 3) // node 1234, node 2001, and cell 87893172902461441 have the "env": "production" label

	// List the objects with kind filter
	objects, err = store.List(context.TODO(), &topoapi.Filters{KindFilter: &topoapi.Filter{
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
	})
	assert.NoError(t, err)
	assert.Len(t, objects, 16) // 2 nodes + 8 cell-neighbors + 6 node-cells

	// List the objects with relation filter. this has an implicit scope of topoapi.RelationFilterScope_TARGET_ONLY
	objects, err = store.List(context.TODO(), &topoapi.Filters{
		RelationFilter: &topoapi.RelationFilter{SrcId: "1234", RelationKind: "e2-node-cell", TargetKind: ""},
	})
	assert.NoError(t, err)
	assert.Len(t, objects, 2) // the 1234 node has two cells

	// List the objects with relation filter and scope All
	objects, err = store.List(context.TODO(), &topoapi.Filters{
		RelationFilter: &topoapi.RelationFilter{SrcId: "1234", RelationKind: "e2-node-cell", TargetKind: "", Scope: topoapi.RelationFilterScope_ALL},
	})
	assert.NoError(t, err)
	assert.Len(t, objects, 5) // 1 node, 2 relations, 2 cells

	// List the objects with relation filter and scope Source and Target
	objects, err = store.List(context.TODO(), &topoapi.Filters{
		RelationFilter: &topoapi.RelationFilter{SrcId: "1234", RelationKind: "e2-node-cell", TargetKind: "", Scope: topoapi.RelationFilterScope_SOURCE_AND_TARGET},
	})
	assert.NoError(t, err)
	assert.Len(t, objects, 3) // 1 node, 2 cells

	objects, err = store.List(context.TODO(), &topoapi.Filters{
		RelationFilter: &topoapi.RelationFilter{SrcId: "87893172902445058", RelationKind: "e2-cell-neighbor", TargetKind: "", Scope: topoapi.RelationFilterScope_SOURCE_AND_TARGET},
	})
	assert.NoError(t, err)
	assert.Len(t, objects, 3) // 58 and neighbors 57 and 59

	objects, err = store.List(context.TODO(), &topoapi.Filters{
		RelationFilter: &topoapi.RelationFilter{SrcId: "87893172902445058", RelationKind: "e2-cell-neighbor", TargetKind: "", Scope: topoapi.RelationFilterScope_TARGET_ONLY},
	})
	assert.NoError(t, err)
	assert.Len(t, objects, 2) // neighbors 57 and 59

	objects, err = store.List(context.TODO(), &topoapi.Filters{
		RelationFilter: &topoapi.RelationFilter{TargetId: "87893172902445058", RelationKind: "e2-cell-neighbor", TargetKind: "", Scope: topoapi.RelationFilterScope_TARGET_ONLY},
	})
	assert.NoError(t, err)
	assert.Len(t, objects, 2) // neighbors 57 and 59

	objects, err = store.List(context.TODO(), &topoapi.Filters{
		RelationFilter: &topoapi.RelationFilter{TargetId: "87893172902445058", RelationKind: "e2-cell-neighbor", TargetKind: "e2-cell", Scope: topoapi.RelationFilterScope_TARGET_ONLY},
	})
	assert.NoError(t, err)
	assert.Len(t, objects, 2) // neighbors 57 and 59

	// List the objects with object type filter
	objects, err = store.List(context.TODO(), &topoapi.Filters{
		ObjectTypes: []topoapi.Object_Type{topo.Object_ENTITY},
	})
	assert.NoError(t, err)
	assert.Len(t, objects, 8) // nodes + cells

	// List the objects with object type filter and aspect Location
	objects, err = store.List(context.TODO(), &topoapi.Filters{
		WithAspects: []string{"onos.topo.Location"},
	})
	assert.NoError(t, err)
	assert.Len(t, objects, 2) // nodes

	// No test for relation filter with target kind: cell-neighbor, node-cell do not have different target kinds
}

type auxNode struct {
	id     string
	labels map[string]string
}
type auxCell struct {
	id     string
	labels map[string]string
}
type auxNodeToCell struct {
	srcID  string
	tgtID  string
	labels map[string]string
}
type auxCellNeighbor struct {
	srcID  string
	tgtID  string
	labels map[string]string
}

func createObjectsListTest(t *testing.T, s Store) {
	createNode(t, s, auxNode{id: "1234", labels: map[string]string{"env": "production"}})
	createNode(t, s, auxNode{id: "2001", labels: map[string]string{"env": "production"}})
	createCell(t, s, auxCell{id: "87893172902461441", labels: map[string]string{"env": "production"}})
	createCell(t, s, auxCell{id: "87893172902461443", labels: map[string]string{"env": "dev"}})
	createCell(t, s, auxCell{id: "87893172902445057", labels: map[string]string{"env": "dev"}})
	createCell(t, s, auxCell{id: "87893172902445058", labels: map[string]string{}})
	createCell(t, s, auxCell{id: "87893172902445059", labels: map[string]string{}})
	createCell(t, s, auxCell{id: "87893172902445060", labels: map[string]string{}})
	createCellNeighbors(t, s, auxCellNeighbor{srcID: "87893172902461441", tgtID: "87893172902461443", labels: map[string]string{}})
	createCellNeighbors(t, s, auxCellNeighbor{srcID: "87893172902445057", tgtID: "87893172902445058", labels: map[string]string{}})
	createCellNeighbors(t, s, auxCellNeighbor{srcID: "87893172902445058", tgtID: "87893172902445059", labels: map[string]string{}})
	createCellNeighbors(t, s, auxCellNeighbor{srcID: "87893172902445059", tgtID: "87893172902445060", labels: map[string]string{}})
	createNodeToCell(t, s, auxNodeToCell{srcID: "1234", tgtID: "87893172902461441", labels: map[string]string{}})
	createNodeToCell(t, s, auxNodeToCell{srcID: "1234", tgtID: "87893172902461443", labels: map[string]string{}})
	createNodeToCell(t, s, auxNodeToCell{srcID: "2001", tgtID: "87893172902445057", labels: map[string]string{}})
	createNodeToCell(t, s, auxNodeToCell{srcID: "2001", tgtID: "87893172902445058", labels: map[string]string{}})
	createNodeToCell(t, s, auxNodeToCell{srcID: "2001", tgtID: "87893172902445059", labels: map[string]string{}})
	createNodeToCell(t, s, auxNodeToCell{srcID: "2001", tgtID: "87893172902445060", labels: map[string]string{}})
}

func createNode(t *testing.T, s Store, a auxNode) {
	object := &topoapi.Object{
		ID:     topo.ID(a.id),
		Type:   topoapi.Object_ENTITY,
		Obj:    &topoapi.Object_Entity{Entity: &topoapi.Entity{KindID: topoapi.ID("e2-node")}},
		Labels: a.labels,
	}
	err := object.SetAspect(&topoapi.Location{Lat: 3.14, Lng: 6.28})
	assert.NoError(t, err)
	err = s.Create(context.TODO(), object)
	assert.NoError(t, err)
}

func createCell(t *testing.T, s Store, a auxCell) {
	err := s.Create(context.TODO(), &topoapi.Object{
		ID:     topo.ID(a.id),
		Type:   topoapi.Object_ENTITY,
		Obj:    &topoapi.Object_Entity{Entity: &topoapi.Entity{KindID: topoapi.ID("e2-cell")}},
		Labels: a.labels,
	})
	assert.NoError(t, err)
}
func createNodeToCell(t *testing.T, s Store, a auxNodeToCell) {
	err := s.Create(context.TODO(), &topoapi.Object{
		// ID: intentionally left empty for the auto-generation to take place
		Type:   topoapi.Object_RELATION,
		Obj:    &topoapi.Object_Relation{Relation: &topoapi.Relation{KindID: topoapi.ID("e2-node-cell"), SrcEntityID: topoapi.ID(a.srcID), TgtEntityID: topoapi.ID(a.tgtID)}},
		Labels: a.labels,
	})
	assert.NoError(t, err)
}

// creates both ways
func createCellNeighbors(t *testing.T, s Store, a auxCellNeighbor) {
	err := s.Create(context.TODO(), &topoapi.Object{
		// ID: intentionally left empty for the auto-generation to take place
		Type:   topoapi.Object_RELATION,
		Obj:    &topoapi.Object_Relation{Relation: &topoapi.Relation{KindID: topoapi.ID("e2-cell-neighbor"), SrcEntityID: topoapi.ID(a.srcID), TgtEntityID: topoapi.ID(a.tgtID)}},
		Labels: a.labels,
	})
	assert.NoError(t, err)
	err = s.Create(context.TODO(), &topoapi.Object{
		// ID: intentionally left empty for the auto-generation to take place
		Type:   topoapi.Object_RELATION,
		Obj:    &topoapi.Object_Relation{Relation: &topoapi.Relation{KindID: topoapi.ID("e2-cell-neighbor"), SrcEntityID: topoapi.ID(a.tgtID), TgtEntityID: topoapi.ID(a.srcID)}},
		Labels: a.labels,
	})
	assert.NoError(t, err)
}
