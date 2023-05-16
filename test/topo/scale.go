// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package topo

import (
	"context"
	"fmt"
	"github.com/gogo/protobuf/types"
	topoapi "github.com/onosproject/onos-api/go/onos/topo"
	utils "github.com/onosproject/onos-topo/test/utils"
)

// TestScale tests topo at scale
func (s *TestSuite) TestScale() {
	s.T().Logf("Creating connection")
	conn, err := utils.CreateConnection()
	s.NoError(err)
	s.T().Logf("Creating Topo Client")
	client := topoapi.NewTopoClient(conn)

	s.T().Logf("Creating 100 nodes, with 6 cells each Nd 6 node-cell relations")
	for n := 0; n < 100; n++ {
		err = s.CreateEntity(client, fmt.Sprintf("node%d", n+1), "e2node", []*types.Any{{TypeUrl: "onos.topo.Location", Value: []byte(`{"lat": 123.0, "lng": 321.0}`)}}, nil)
		s.NoError(err)

		for c := 0; c < 6; c++ {
			err = s.CreateEntity(client, fmt.Sprintf("cell%d%d", n+1, c+1), "e2cell", []*types.Any{{TypeUrl: "onos.topo.Location", Value: []byte(`{"lat": 123.0, "lng": 321.0}`)}}, nil)
			s.NoError(err)
			err = CreateRelation(client, fmt.Sprintf("node%d", n+1), fmt.Sprintf("cell%d%d", n+1, c+1), "contains")
			s.NoError(err)
		}
	}

	// Filter e2nodes; there should be 100
	s.T().Logf("Getting all 'e2nodes'")
	res, err := client.List(context.Background(), &topoapi.ListRequest{Filters: &topoapi.Filters{
		KindFilter: &topoapi.Filter{Filter: &topoapi.Filter_Equal_{Equal_: &topoapi.EqualFilter{Value: "e2node"}}},
	}})
	s.NoError(err)
	s.Len(res.Objects, 100)

	// Filter e2cells; there should be 600
	s.T().Logf("Getting all 'e2cells'")
	res, err = client.List(context.Background(), &topoapi.ListRequest{Filters: &topoapi.Filters{
		KindFilter: &topoapi.Filter{Filter: &topoapi.Filter_Equal_{Equal_: &topoapi.EqualFilter{Value: "e2cell"}}},
	}})
	s.NoError(err)
	s.Len(res.Objects, 600)

	// Filter contains relations; there should be 600
	s.T().Logf("Getting all 'contains' relations'")
	res, err = client.List(context.Background(), &topoapi.ListRequest{Filters: &topoapi.Filters{
		KindFilter: &topoapi.Filter{Filter: &topoapi.Filter_Equal_{Equal_: &topoapi.EqualFilter{Value: "contains"}}},
	}})
	s.NoError(err)
	s.Len(res.Objects, 600)

	// Filter e2cells of an e2node; there should be 6
	s.T().Logf("Getting all cells of a node")
	res, err = client.List(context.Background(), &topoapi.ListRequest{Filters: &topoapi.Filters{
		RelationFilter: &topoapi.RelationFilter{
			SrcId:        "node10",
			RelationKind: "contains",
			TargetKind:   "e2cell",
		},
	}})
	s.NoError(err)
	s.Len(res.Objects, 6)

	// Filter e2nodes, e2cells, contains, but only entities; there should be 700
	s.T().Logf("Getting all nodes and cells")
	res, err = client.List(context.Background(), &topoapi.ListRequest{Filters: &topoapi.Filters{
		KindFilter:  &topoapi.Filter{Filter: &topoapi.Filter_In{In: &topoapi.InFilter{Values: []string{"e2node", "e2cell", "contains"}}}},
		ObjectTypes: []topoapi.Object_Type{topoapi.Object_ENTITY},
	}})
	s.NoError(err)
	s.Len(res.Objects, 700)
}
