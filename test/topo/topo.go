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
	"fmt"
	"github.com/gogo/protobuf/types"
	topoapi "github.com/onosproject/onos-api/go/onos/topo"
	utils "github.com/onosproject/onos-topo/test/utils"
	"testing"

	"gotest.tools/assert"
)

func createEntity(client topoapi.TopoClient, id string, kindID string, aspectList []*types.Any) error {
	aspects := map[string]*types.Any{}
	for _, aspect := range aspectList {
		aspects[aspect.TypeUrl] = aspect
	}
	_, err := client.Create(context.Background(), &topoapi.CreateRequest{
		Object: &topoapi.Object{
			ID:      topoapi.ID(id),
			Type:    topoapi.Object_ENTITY,
			Aspects: aspects,
			Obj:     &topoapi.Object_Entity{Entity: &topoapi.Entity{KindID: topoapi.ID(kindID)}},
		},
	})
	return err
}

func createRelation(client topoapi.TopoClient, src string, tgt string, kindID string) error {
	_, err := client.Create(context.Background(), &topoapi.CreateRequest{
		Object: &topoapi.Object{
			ID:   topoapi.ID(src + tgt + kindID),
			Type: topoapi.Object_RELATION,
			Obj: &topoapi.Object_Relation{
				Relation: &topoapi.Relation{
					SrcEntityID: topoapi.ID(src),
					TgtEntityID: topoapi.ID(tgt),
					KindID:      topoapi.ID(kindID),
				},
			},
		},
	})
	return err
}

// TestAddRemoveDevice adds devices to the storage, lists and checks that they are in database and removes devices from the storage
func (s *TestSuite) TestAddRemoveDevice(t *testing.T) {
	t.Logf("Creating connection")
	conn, err := utils.CreateConnection()
	assert.NilError(t, err)
	t.Logf("Creating Topo Client")
	client := topoapi.NewTopoClient(conn)

	t.Logf("Adding first device to the topo store")
	err = createEntity(client, "1", "testKind", []*types.Any{
		{TypeUrl: "onos.topo.Location", Value: []byte(`{"lat": 123, "lng": 321}`)},
		{TypeUrl: "foo", Value: []byte(`{"s": "barfoo", "n": 314, "b": true}`)},
	})
	assert.NilError(t, err)

	t.Logf("Adding second device to the topo store")
	err = createEntity(client, "2", "testKind", []*types.Any{
		{TypeUrl: "onos.topo.Location", Value: []byte(`{"lat": 111, "lng": 222}`)},
		{TypeUrl: "foo", Value: []byte(`{"s": "foobar", "n": 628, "b": true}`)},
	})
	assert.NilError(t, err)

	t.Logf("Checking whether added device exists")
	gres, err := client.Get(context.Background(), &topoapi.GetRequest{
		ID: "1",
	})
	assert.NilError(t, err)
	assert.Equal(t, topoapi.ID("1"), gres.Object.ID)

	t.Logf("Listing all devices")
	res, err := client.List(context.Background(), &topoapi.ListRequest{})
	assert.NilError(t, err)
	t.Logf("Verifying that there are two devices stored")
	assert.Equal(t, len(res.Objects) == 2 &&
		(res.Objects[0].ID == "1" || res.Objects[1].ID == "1"), true)

	t.Logf("Updating first device")
	obj := gres.Object
	obj.Aspects["onos.topo.Location"] = &types.Any{TypeUrl: "onos.topo.Location", Value: []byte(`{"lat": 123, "lng": 321}`)}

	ures, err := client.Update(context.Background(), &topoapi.UpdateRequest{
		Object: obj,
	})
	assert.NilError(t, err)
	assert.Assert(t, ures != nil)

	t.Logf("Deleting first device")
	_, err = client.Delete(context.Background(), &topoapi.DeleteRequest{
		ID: "1",
	})
	assert.NilError(t, err)

	t.Logf("Listing all devices and verifying that there is only second device left")
	res, err = client.List(context.Background(), &topoapi.ListRequest{})
	assert.NilError(t, err)
	assert.Equal(t, len(res.Objects) == 1 && res.Objects[0].ID == "2", true)
}

func (s *TestSuite) TestScale(t *testing.T) {
	t.Logf("Creating connection")
	conn, err := utils.CreateConnection()
	assert.NilError(t, err)
	t.Logf("Creating Topo Client")
	client := topoapi.NewTopoClient(conn)

	t.Logf("Creating 100 nodes, with 6 cells each Nd 6 node-cell relations")
	for n := 0; n < 100; n++ {
		err = createEntity(client, fmt.Sprintf("node%d", n+1), "e2node", []*types.Any{{TypeUrl: "onos.topo.Location", Value: []byte(`{"lat": 123.0, "lng": 321.0}`)}})
		assert.NilError(t, err)

		for c := 0; c < 6; c++ {
			err = createEntity(client, fmt.Sprintf("cell%d%d", n+1, c+1), "e2cell", []*types.Any{{TypeUrl: "onos.topo.Location", Value: []byte(`{"lat": 123.0, "lng": 321.0}`)}})
			assert.NilError(t, err)
			err = createRelation(client, fmt.Sprintf("node%d", n+1), fmt.Sprintf("cell%d%d", n+1, c+1), "contains")
			assert.NilError(t, err)
		}
	}

	// Filter e2nodes; there should be 100
	t.Logf("Getting all 'e2nodes'")
	res, err := client.List(context.Background(), &topoapi.ListRequest{Filters: &topoapi.Filters{
		KindFilter: &topoapi.Filter{Filter: &topoapi.Filter_Equal_{Equal_: &topoapi.EqualFilter{Value: "e2node"}}},
	}})
	assert.NilError(t, err)
	assert.Equal(t, len(res.Objects), 100)

	// Filter e2cells; there should be 600
	t.Logf("Getting all 'e2cells'")
	res, err = client.List(context.Background(), &topoapi.ListRequest{Filters: &topoapi.Filters{
		KindFilter: &topoapi.Filter{Filter: &topoapi.Filter_Equal_{Equal_: &topoapi.EqualFilter{Value: "e2cell"}}},
	}})
	assert.NilError(t, err)
	assert.Equal(t, len(res.Objects), 600)

	// Filter contains relations; there should be 600
	t.Logf("Getting all 'contains' relations'")
	res, err = client.List(context.Background(), &topoapi.ListRequest{Filters: &topoapi.Filters{
		KindFilter: &topoapi.Filter{Filter: &topoapi.Filter_Equal_{Equal_: &topoapi.EqualFilter{Value: "contains"}}},
	}})
	assert.NilError(t, err)
	assert.Equal(t, len(res.Objects), 600)

	// Filter e2cells of an e2node; there should be 6
	t.Logf("Getting all cells of a node")
	res, err = client.List(context.Background(), &topoapi.ListRequest{Filters: &topoapi.Filters{
		RelationFilter: &topoapi.RelationFilter{
			SrcId:        "node10",
			RelationKind: "contains",
			TargetKind:   "e2cell",
		},
	}})
	assert.NilError(t, err)
	assert.Equal(t, len(res.Objects), 6)

	// Filter e2nodes, e2cells, contains, but only entities; there should be 700
	t.Logf("Getting all nodes and cells")
	res, err = client.List(context.Background(), &topoapi.ListRequest{Filters: &topoapi.Filters{
		KindFilter: &topoapi.Filter{Filter: &topoapi.Filter_In{In: &topoapi.InFilter{Values: []string{"e2node", "e2cell", "contains"}}}},
		ObjectTypes: []topoapi.Object_Type{topoapi.Object_ENTITY},
	}})
	assert.NilError(t, err)
	assert.Equal(t, len(res.Objects), 700)
}
