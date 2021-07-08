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

// TestScale tests topo at scale
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
		KindFilter:  &topoapi.Filter{Filter: &topoapi.Filter_In{In: &topoapi.InFilter{Values: []string{"e2node", "e2cell", "contains"}}}},
		ObjectTypes: []topoapi.Object_Type{topoapi.Object_ENTITY},
	}})
	assert.NilError(t, err)
	assert.Equal(t, len(res.Objects), 700)
}
