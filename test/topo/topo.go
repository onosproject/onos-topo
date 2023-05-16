// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package topo

import (
	"context"
	"github.com/gogo/protobuf/types"
	topoapi "github.com/onosproject/onos-api/go/onos/topo"
	utils "github.com/onosproject/onos-topo/test/utils"
)

// CreateEntity creates an entity object
func (s *TestSuite) CreateEntity(client topoapi.TopoClient, id string, kindID string, aspectList []*types.Any, labels map[string]string) error {
	aspects := map[string]*types.Any{}
	for _, aspect := range aspectList {
		aspects[aspect.TypeUrl] = aspect
	}
	_, err := client.Create(s.Context(), &topoapi.CreateRequest{
		Object: &topoapi.Object{
			ID:      topoapi.ID(id),
			Type:    topoapi.Object_ENTITY,
			Aspects: aspects,
			Obj:     &topoapi.Object_Entity{Entity: &topoapi.Entity{KindID: topoapi.ID(kindID)}},
			Labels:  labels,
		},
	})
	return err
}

// CreateRelation creates a relation object
func CreateRelation(client topoapi.TopoClient, src string, tgt string, kindID string) error {
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
func (s *TestSuite) TestAddRemoveDevice() {
	s.T().Logf("Creating connection")
	conn, err := utils.CreateConnection()
	s.NoError(err)
	s.T().Logf("Creating Topo Client")
	client := topoapi.NewTopoClient(conn)

	s.T().Logf("Adding first device to the topo store")
	err = s.CreateEntity(client, "1", "testKind", []*types.Any{
		{TypeUrl: "onos.topo.Location", Value: []byte(`{"lat": 123, "lng": 321}`)},
		{TypeUrl: "foo", Value: []byte(`{"s": "barfoo", "n": 314, "b": true}`)},
	}, nil)
	s.NoError(err)

	s.T().Logf("Adding second device to the topo store")
	err = s.CreateEntity(client, "2", "testKind", []*types.Any{
		{TypeUrl: "onos.topo.Location", Value: []byte(`{"lat": 111, "lng": 222}`)},
		{TypeUrl: "foo", Value: []byte(`{"s": "foobar", "n": 628, "b": true}`)},
	}, nil)
	s.NoError(err)

	s.T().Logf("Checking whether added device exists")
	gres, err := client.Get(context.Background(), &topoapi.GetRequest{
		ID: "1",
	})
	s.NoError(err)
	s.Equal(topoapi.ID("1"), gres.Object.ID)

	equalFilter := &topoapi.Filter_Equal_{Equal_: &topoapi.EqualFilter{Value: "testKind"}}
	filters := &topoapi.Filters{
		KindFilter: &topoapi.Filter{
			Filter: equalFilter,
		},
	}

	s.T().Logf("Listing all devices")
	res, err := client.List(context.Background(), &topoapi.ListRequest{
		Filters: filters,
	})
	s.NoError(err)
	s.T().Logf("Verifying that there are two devices stored")
	s.Len(res.Objects, 2)
	s.True(string(res.Objects[0].ID) == "1" || string(res.Objects[1].ID) == "1")

	s.T().Logf("Updating first device")
	obj := gres.Object
	obj.Aspects["onos.topo.Location"] = &types.Any{TypeUrl: "onos.topo.Location", Value: []byte(`{"lat": 123, "lng": 321}`)}

	ures, err := client.Update(context.Background(), &topoapi.UpdateRequest{
		Object: obj,
	})
	s.NoError(err)
	s.NotNil(ures)

	s.T().Logf("Deleting first device")
	_, err = client.Delete(context.Background(), &topoapi.DeleteRequest{
		ID: "1",
	})
	s.NoError(err)

	s.T().Logf("Listing all devices and verifying that there is only second device left")
	res, err = client.List(context.Background(), &topoapi.ListRequest{
		Filters: filters,
	})
	s.NoError(err)
	s.Len(res.Objects, 1)
	s.Equal("2", string(res.Objects[0].ID))
}
