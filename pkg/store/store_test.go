// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package store

import (
	"context"
	"github.com/atomix/go-sdk/pkg/test"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"testing"
	"time"

	"github.com/onosproject/onos-api/go/onos/topo"
	"github.com/stretchr/testify/assert"
)

func TestTopoStore(t *testing.T) {
	cluster := test.NewClient()
	defer cluster.Close()

	store1, err := NewAtomixStore(cluster)
	assert.NoError(t, err)

	store2, err := NewAtomixStore(cluster)
	assert.NoError(t, err)

	// List the objects; there should be none
	noobjects, err := store1.List(context.TODO(), nil)
	assert.NoError(t, err)
	assert.Len(t, noobjects, 0)

	ch := make(chan topo.Event)
	err = store2.Watch(context.Background(), ch, nil)
	assert.NoError(t, err)

	k1 := &topo.Object{
		ID:   "foo",
		Type: topo.Object_KIND,
		Obj:  &topo.Object_Kind{Kind: &topo.Kind{Name: "Foo"}},
	}
	err = store1.Create(context.TODO(), k1)
	assert.NoError(t, err)

	k2 := &topo.Object{
		ID:   "bar",
		Type: topo.Object_KIND,
		Obj:  &topo.Object_Kind{Kind: &topo.Kind{Name: "Bar"}},
	}
	err = store1.Create(context.TODO(), k2)
	assert.NoError(t, err)

	obj1 := &topo.Object{
		ID:     "o1",
		Type:   topo.Object_ENTITY,
		Obj:    &topo.Object_Entity{Entity: &topo.Entity{KindID: topo.ID("foo")}},
		Labels: map[string]string{},
	}
	obj1.Labels["env"] = "test"
	obj1.Labels["area"] = "ran"

	obj2 := &topo.Object{
		ID:     "o2",
		Type:   topo.Object_ENTITY,
		Obj:    &topo.Object_Entity{Entity: &topo.Entity{KindID: topo.ID("bar")}},
		Labels: map[string]string{},
	}
	obj2.Labels["env"] = "production"
	obj2.Labels["area"] = "ran"

	// Create a new object
	err = store1.Create(context.TODO(), obj1)
	assert.NoError(t, err)
	assert.Equal(t, topo.ID("o1"), obj1.ID)
	assert.NotEqual(t, topo.Revision(0), obj1.Revision)

	// Get the object
	obj1, err = store2.Get(context.TODO(), "o1")
	assert.NoError(t, err)
	assert.NotNil(t, obj1)
	assert.Equal(t, topo.ID("o1"), obj1.ID)
	assert.NotEqual(t, topo.Revision(0), obj1.Revision)
	assert.Equal(t, "test", obj1.Labels["env"])
	assert.Equal(t, "ran", obj1.Labels["area"])

	// Create another object
	err = store2.Create(context.TODO(), obj2)
	assert.NoError(t, err)
	assert.Equal(t, topo.ID("o2"), obj2.ID)
	assert.NotEqual(t, topo.Revision(0), obj2.Revision)

	// Verify events were received for the kinds
	topoEvent := nextEvent(t, ch)
	assert.Equal(t, topo.ID("foo"), topoEvent.ID)
	topoEvent = nextEvent(t, ch)
	assert.Equal(t, topo.ID("bar"), topoEvent.ID)

	// Verify events were received for the objects
	topoEvent = nextEvent(t, ch)
	assert.Equal(t, topo.ID("o1"), topoEvent.ID)
	topoEvent = nextEvent(t, ch)
	assert.Equal(t, topo.ID("o2"), topoEvent.ID)

	// Update one of the objects
	err = obj2.SetAspect(&topo.Location{Lat: 1, Lng: 2})
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

	err = obj11.SetAspect(&topo.Location{Lat: 2, Lng: 1})
	assert.NoError(t, err)
	err = store1.Update(context.TODO(), obj11)
	assert.NoError(t, err)

	err = obj12.SetAspect(&topo.E2Node{})
	assert.NoError(t, err)
	err = store2.Update(context.TODO(), obj12)
	assert.Error(t, err)

	// Verify events were received again
	topoEvent = nextEvent(t, ch)
	assert.Equal(t, topo.ID("o2"), topoEvent.ID)
	topoEvent = nextEvent(t, ch)
	assert.Equal(t, topo.ID("o2"), topoEvent.ID)
	topoEvent = nextEvent(t, ch)
	assert.Equal(t, topo.ID("o1"), topoEvent.ID)

	// Verify the attribute values
	obj2g, err := store1.Get(context.TODO(), obj2.ID)
	assert.NoError(t, err)
	loc := &topo.Location{}
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
	assert.Nil(t, obj2)
	assert.True(t, errors.IsNotFound(err))

	obj := &topo.Object{
		ID: "o1",
	}

	err = store1.Create(context.TODO(), obj)
	assert.Error(t, err)

	obj = &topo.Object{
		ID:   "o2",
		Type: topo.Object_ENTITY,
	}

	err = store1.Create(context.TODO(), obj)
	assert.NoError(t, err)

	ch = make(chan topo.Event)
	err = store1.Watch(context.TODO(), ch, nil, WithReplay())
	assert.NoError(t, err)

	obj = nextEvent(t, ch)
	assert.NotNil(t, obj)
	obj = nextEvent(t, ch)
	assert.NotNil(t, obj)
}

func nextEvent(t *testing.T, ch chan topo.Event) *topo.Object {
	select {
	case c := <-ch:
		return &c.Object
	case <-time.After(5 * time.Second):
		t.FailNow()
	}
	return nil
}

func TestList(t *testing.T) {
	cluster := test.NewClient()
	defer cluster.Close()

	// Store def
	store, _ := NewAtomixStore(cluster)

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
	objects, err = store.List(context.TODO(), &topo.Filters{ObjectTypes: []topo.Object_Type{topo.Object_ENTITY}})
	assert.NoError(t, err)
	assert.Len(t, objects, 8) // 2 nodes + 6 cells

	// List the objects with label filter
	objects, err = store.List(context.TODO(), &topo.Filters{LabelFilters: []*topo.Filter{
		{
			Filter: &topo.Filter_Equal_{
				Equal_: &topo.EqualFilter{Value: "production"},
			},
			Key: "env",
		},
	}})
	assert.NoError(t, err)
	assert.Len(t, objects, 3) // node 1234, node 2001, and cell 87893172902461441 have the "env": "production" label

	// List the objects with kind filter
	objects, err = store.List(context.TODO(), &topo.Filters{KindFilter: &topo.Filter{
		Filter: &topo.Filter_Not{
			Not: &topo.NotFilter{
				Inner: &topo.Filter{
					Filter: &topo.Filter_Equal_{
						Equal_: &topo.EqualFilter{Value: "e2-cell"},
					},
				},
			},
		},
	},
	})
	assert.NoError(t, err)
	assert.Len(t, objects, 16) // 2 nodes + 8 cell-neighbors + 6 node-cells

	// List the objects with relation filter. this has an implicit scope of topo.RelationFilterScope_TARGET_ONLY
	objects, err = store.List(context.TODO(), &topo.Filters{
		RelationFilter: &topo.RelationFilter{SrcId: "1234", RelationKind: "e2-node-cell", TargetKind: ""},
	})
	assert.NoError(t, err)
	assert.Len(t, objects, 2) // the 1234 node has two cells

	// List the objects with relation filter and scope All
	objects, err = store.List(context.TODO(), &topo.Filters{
		RelationFilter: &topo.RelationFilter{SrcId: "1234", RelationKind: "e2-node-cell", TargetKind: "", Scope: topo.RelationFilterScope_ALL},
	})
	assert.NoError(t, err)
	assert.Len(t, objects, 5) // 1 node, 2 relations, 2 cells

	// List the objects with relation filter and scope Source and Target
	objects, err = store.List(context.TODO(), &topo.Filters{
		RelationFilter: &topo.RelationFilter{SrcId: "1234", RelationKind: "e2-node-cell", TargetKind: "", Scope: topo.RelationFilterScope_SOURCE_AND_TARGETS},
	})
	assert.NoError(t, err)
	assert.Len(t, objects, 3) // 1 node, 2 cells

	objects, err = store.List(context.TODO(), &topo.Filters{
		RelationFilter: &topo.RelationFilter{SrcId: "87893172902445058", RelationKind: "e2-cell-neighbor", TargetKind: "", Scope: topo.RelationFilterScope_SOURCE_AND_TARGETS},
	})
	assert.NoError(t, err)
	assert.Len(t, objects, 3) // 58 and neighbors 57 and 59

	objects, err = store.List(context.TODO(), &topo.Filters{
		RelationFilter: &topo.RelationFilter{SrcId: "87893172902445058", RelationKind: "e2-cell-neighbor", TargetKind: "", Scope: topo.RelationFilterScope_TARGETS_ONLY},
	})
	assert.NoError(t, err)
	assert.Len(t, objects, 2) // neighbors 57 and 59
	assert.Equal(t, topo.Object_ENTITY, objects[0].Type)

	objects, err = store.List(context.TODO(), &topo.Filters{
		RelationFilter: &topo.RelationFilter{SrcId: "87893172902445058", RelationKind: "e2-cell-neighbor", TargetKind: "", Scope: topo.RelationFilterScope_RELATIONS_ONLY},
	})
	assert.NoError(t, err)
	assert.Len(t, objects, 2) // neighbors 57 and 59 relations
	assert.Equal(t, topo.Object_RELATION, objects[0].Type)

	objects, err = store.List(context.TODO(), &topo.Filters{
		RelationFilter: &topo.RelationFilter{SrcId: "87893172902445058", RelationKind: "e2-cell-neighbor", TargetKind: "", Scope: topo.RelationFilterScope_RELATIONS_AND_TARGETS},
	})
	assert.NoError(t, err)
	assert.Len(t, objects, 4) // neighbors 57 and 59 entities and corresponding relations

	objects, err = store.List(context.TODO(), &topo.Filters{
		RelationFilter: &topo.RelationFilter{TargetId: "87893172902445058", RelationKind: "e2-cell-neighbor", TargetKind: "", Scope: topo.RelationFilterScope_TARGETS_ONLY},
	})
	assert.NoError(t, err)
	assert.Len(t, objects, 2) // neighbors 57 and 59
	assert.Equal(t, topo.Object_ENTITY, objects[0].Type)

	objects, err = store.List(context.TODO(), &topo.Filters{
		RelationFilter: &topo.RelationFilter{TargetId: "87893172902445058", RelationKind: "e2-cell-neighbor", TargetKind: "e2-cell", Scope: topo.RelationFilterScope_TARGETS_ONLY},
	})
	assert.NoError(t, err)
	assert.Len(t, objects, 2) // neighbors 57 and 59
	assert.Equal(t, topo.Object_ENTITY, objects[0].Type)

	objects, err = store.List(context.TODO(), &topo.Filters{
		RelationFilter: &topo.RelationFilter{TargetId: "87893172902445058", RelationKind: "e2-cell-neighbor", TargetKind: "e2-cell", Scope: topo.RelationFilterScope_RELATIONS_ONLY},
	})
	assert.NoError(t, err)
	assert.Len(t, objects, 2) // neighbors 57 and 59 relations
	assert.Equal(t, topo.Object_RELATION, objects[0].Type)

	objects, err = store.List(context.TODO(), &topo.Filters{
		RelationFilter: &topo.RelationFilter{TargetId: "87893172902445058", RelationKind: "e2-cell-neighbor", TargetKind: "e2-cell", Scope: topo.RelationFilterScope_RELATIONS_AND_TARGETS},
	})
	assert.NoError(t, err)
	assert.Len(t, objects, 4) // neighbors 57 and 59 entities and corresponding relations

	// List the objects with object type filter
	objects, err = store.List(context.TODO(), &topo.Filters{
		ObjectTypes: []topo.Object_Type{topo.Object_ENTITY},
	})
	assert.NoError(t, err)
	assert.Len(t, objects, 8) // nodes + cells

	// List the objects with object type filter and aspect Location
	objects, err = store.List(context.TODO(), &topo.Filters{
		WithAspects: []string{"onos.topo.Location"},
	})
	assert.NoError(t, err)
	assert.Len(t, objects, 2) // nodes

	// No test for relation filter with target kind: cell-neighbor, node-cell do not have different target kinds
}

const depth = 512

func TestQuery(t *testing.T) {
	cluster := test.NewClient()
	defer cluster.Close()

	// Store def
	store, _ := NewAtomixStore(cluster)

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
	ch := make(chan *topo.Object, depth)
	err = store.Query(context.TODO(), ch, nil)
	assert.NoError(t, err)
	assert.Equal(t, consume(ch), 22) // 2 nodes + 6 cells + 8 cell-neighbors + 6 node-cells

	// List the objects
	ch = make(chan *topo.Object, depth)
	err = store.Query(context.TODO(), ch, &topo.Filters{ObjectTypes: []topo.Object_Type{topo.Object_ENTITY}})
	assert.NoError(t, err)
	assert.Equal(t, consume(ch), 8) // 2 nodes + 6 cells

	// List the objects with label filter
	ch = make(chan *topo.Object, depth)
	err = store.Query(context.TODO(), ch, &topo.Filters{LabelFilters: []*topo.Filter{
		{
			Filter: &topo.Filter_Equal_{
				Equal_: &topo.EqualFilter{Value: "production"},
			},
			Key: "env",
		},
	}})
	assert.NoError(t, err)
	assert.Equal(t, 3, consume(ch)) // node 1234, node 2001, and cell 87893172902461441 have the "env": "production" label

	// List the objects with kind filter
	ch = make(chan *topo.Object, depth)
	err = store.Query(context.TODO(), ch, &topo.Filters{KindFilter: &topo.Filter{
		Filter: &topo.Filter_Not{
			Not: &topo.NotFilter{
				Inner: &topo.Filter{
					Filter: &topo.Filter_Equal_{
						Equal_: &topo.EqualFilter{Value: "e2-cell"},
					},
				},
			},
		},
	},
	})
	assert.NoError(t, err)
	assert.Equal(t, 16, consume(ch)) // 2 nodes + 8 cell-neighbors + 6 node-cells

	// List the objects with relation filter. this has an implicit scope of topo.RelationFilterScope_TARGET_ONLY
	ch = make(chan *topo.Object, depth)
	err = store.Query(context.TODO(), ch, &topo.Filters{
		RelationFilter: &topo.RelationFilter{SrcId: "1234", RelationKind: "e2-node-cell", TargetKind: ""},
	})
	assert.NoError(t, err)
	assert.Equal(t, 2, consume(ch)) // the 1234 node has two cells

	// List the objects with relation filter and scope All
	ch = make(chan *topo.Object, depth)
	err = store.Query(context.TODO(), ch, &topo.Filters{
		RelationFilter: &topo.RelationFilter{SrcId: "1234", RelationKind: "e2-node-cell", TargetKind: "", Scope: topo.RelationFilterScope_ALL},
	})
	assert.NoError(t, err)
	assert.Equal(t, 5, consume(ch)) // 1 node, 2 relations, 2 cells

	// List the objects with relation filter and scope Source and Target
	ch = make(chan *topo.Object, depth)
	err = store.Query(context.TODO(), ch, &topo.Filters{
		RelationFilter: &topo.RelationFilter{SrcId: "1234", RelationKind: "e2-node-cell", TargetKind: "", Scope: topo.RelationFilterScope_SOURCE_AND_TARGETS},
	})
	assert.NoError(t, err)
	assert.Equal(t, 3, consume(ch)) // 1 node, 2 cells

	ch = make(chan *topo.Object, depth)
	err = store.Query(context.TODO(), ch, &topo.Filters{
		RelationFilter: &topo.RelationFilter{SrcId: "87893172902445058", RelationKind: "e2-cell-neighbor", TargetKind: "", Scope: topo.RelationFilterScope_SOURCE_AND_TARGETS},
	})
	assert.NoError(t, err)
	assert.Equal(t, 3, consume(ch)) // 58 and neighbors 57 and 59

	ch = make(chan *topo.Object, depth)
	err = store.Query(context.TODO(), ch, &topo.Filters{
		RelationFilter: &topo.RelationFilter{SrcId: "87893172902445058", RelationKind: "e2-cell-neighbor", TargetKind: "", Scope: topo.RelationFilterScope_TARGETS_ONLY},
	})
	assert.NoError(t, err)
	assert.Equal(t, 2, consume(ch)) // neighbors 57 and 59
	//assert.Equal(t, topo.Object_ENTITY, objects[0].Type)

	ch = make(chan *topo.Object, depth)
	err = store.Query(context.TODO(), ch, &topo.Filters{
		RelationFilter: &topo.RelationFilter{SrcId: "87893172902445058", RelationKind: "e2-cell-neighbor", TargetKind: "", Scope: topo.RelationFilterScope_RELATIONS_ONLY},
	})
	assert.NoError(t, err)
	assert.Equal(t, 2, consume(ch)) // neighbors 57 and 59 relations
	//assert.Equal(t, topo.Object_RELATION, objects[0].Type)

	ch = make(chan *topo.Object, depth)
	err = store.Query(context.TODO(), ch, &topo.Filters{
		RelationFilter: &topo.RelationFilter{SrcId: "87893172902445058", RelationKind: "e2-cell-neighbor", TargetKind: "", Scope: topo.RelationFilterScope_RELATIONS_AND_TARGETS},
	})
	assert.NoError(t, err)
	assert.Equal(t, 4, consume(ch)) // neighbors 57 and 59 entities and corresponding relations

	ch = make(chan *topo.Object, depth)
	err = store.Query(context.TODO(), ch, &topo.Filters{
		RelationFilter: &topo.RelationFilter{TargetId: "87893172902445058", RelationKind: "e2-cell-neighbor", TargetKind: "", Scope: topo.RelationFilterScope_TARGETS_ONLY},
	})
	assert.NoError(t, err)
	assert.Equal(t, 2, consume(ch)) // neighbors 57 and 59
	//assert.Equal(t, topo.Object_ENTITY, objects[0].Type)

	ch = make(chan *topo.Object, depth)
	err = store.Query(context.TODO(), ch, &topo.Filters{
		RelationFilter: &topo.RelationFilter{TargetId: "87893172902445058", RelationKind: "e2-cell-neighbor", TargetKind: "e2-cell", Scope: topo.RelationFilterScope_TARGETS_ONLY},
	})
	assert.NoError(t, err)
	assert.Equal(t, 2, consume(ch)) // neighbors 57 and 59
	//assert.Equal(t, topo.Object_ENTITY, objects[0].Type)

	ch = make(chan *topo.Object, depth)
	err = store.Query(context.TODO(), ch, &topo.Filters{
		RelationFilter: &topo.RelationFilter{TargetId: "87893172902445058", RelationKind: "e2-cell-neighbor", TargetKind: "e2-cell", Scope: topo.RelationFilterScope_RELATIONS_ONLY},
	})
	assert.NoError(t, err)
	assert.Equal(t, 2, consume(ch)) // neighbors 57 and 59 relations
	//assert.Equal(t, topo.Object_RELATION, objects[0].Type)

	ch = make(chan *topo.Object, depth)
	err = store.Query(context.TODO(), ch, &topo.Filters{
		RelationFilter: &topo.RelationFilter{TargetId: "87893172902445058", RelationKind: "e2-cell-neighbor", TargetKind: "e2-cell", Scope: topo.RelationFilterScope_RELATIONS_AND_TARGETS},
	})
	assert.NoError(t, err)
	assert.Equal(t, 4, consume(ch)) // neighbors 57 and 59 entities and corresponding relations

	// List the objects with object type filter
	ch = make(chan *topo.Object, depth)
	err = store.Query(context.TODO(), ch, &topo.Filters{
		ObjectTypes: []topo.Object_Type{topo.Object_ENTITY},
	})
	assert.NoError(t, err)
	assert.Equal(t, 8, consume(ch)) // nodes + cells

	// List the objects with object type filter and aspect Location
	ch = make(chan *topo.Object, depth)
	err = store.Query(context.TODO(), ch, &topo.Filters{
		WithAspects: []string{"onos.topo.Location"},
	})
	assert.NoError(t, err)
	assert.Equal(t, 2, consume(ch)) // nodes

	// No test for relation filter with target kind: cell-neighbor, node-cell do not have different target kinds
}

func consume(ch chan *topo.Object) int {
	count := 0
	for range ch {
		count++
	}
	return count
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
	object := &topo.Object{
		ID:     topo.ID(a.id),
		Type:   topo.Object_ENTITY,
		Obj:    &topo.Object_Entity{Entity: &topo.Entity{KindID: topo.ID("e2-node")}},
		Labels: a.labels,
	}
	err := object.SetAspect(&topo.Location{Lat: 3.14, Lng: 6.28})
	assert.NoError(t, err)
	err = s.Create(context.TODO(), object)
	assert.NoError(t, err)
}

func createCell(t *testing.T, s Store, a auxCell) {
	err := s.Create(context.TODO(), &topo.Object{
		ID:     topo.ID(a.id),
		Type:   topo.Object_ENTITY,
		Obj:    &topo.Object_Entity{Entity: &topo.Entity{KindID: topo.ID("e2-cell")}},
		Labels: a.labels,
	})
	assert.NoError(t, err)
}
func createNodeToCell(t *testing.T, s Store, a auxNodeToCell) {
	err := s.Create(context.TODO(), &topo.Object{
		// ID: intentionally left empty for the auto-generation to take place
		Type:   topo.Object_RELATION,
		Obj:    &topo.Object_Relation{Relation: &topo.Relation{KindID: topo.ID("e2-node-cell"), SrcEntityID: topo.ID(a.srcID), TgtEntityID: topo.ID(a.tgtID)}},
		Labels: a.labels,
	})
	assert.NoError(t, err)
}

// creates both ways
func createCellNeighbors(t *testing.T, s Store, a auxCellNeighbor) {
	err := s.Create(context.TODO(), &topo.Object{
		// ID: intentionally left empty for the auto-generation to take place
		Type:   topo.Object_RELATION,
		Obj:    &topo.Object_Relation{Relation: &topo.Relation{KindID: topo.ID("e2-cell-neighbor"), SrcEntityID: topo.ID(a.srcID), TgtEntityID: topo.ID(a.tgtID)}},
		Labels: a.labels,
	})
	assert.NoError(t, err)
	err = s.Create(context.TODO(), &topo.Object{
		// ID: intentionally left empty for the auto-generation to take place
		Type:   topo.Object_RELATION,
		Obj:    &topo.Object_Relation{Relation: &topo.Relation{KindID: topo.ID("e2-cell-neighbor"), SrcEntityID: topo.ID(a.tgtID), TgtEntityID: topo.ID(a.srcID)}},
		Labels: a.labels,
	})
	assert.NoError(t, err)
}
