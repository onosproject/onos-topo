// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package topo

import (
	"context"
	topoapi "github.com/onosproject/onos-api/go/onos/topo"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	utils "github.com/onosproject/onos-topo/test/utils"
	"testing"

	"gotest.tools/assert"
)

const (
	numRequestedE2Nodes = 50
)

//var (
//	initialEnbID  = 155000
//	serviceModels = []string{"kpm2", "rcpre2"}
//	controllers   = []string{"e2t-1"}
//)

var log = logging.GetLogger("topo")

func (s *TestSuite) TestAddRemoveDevice(t *testing.T) {

	//ToDo - how to define it properly?? What is the correct way?
	//st, err := utils.NewClient()
	//assert.NilError(t, err, "Error creating client")

	//topo := topoapi.Object{}
	//err = st.Create(context.Background(), &topo)
	//assert.NilError(t, err, "Error creating topo")

	st := utils.Rnib{}
	//st, err := utils.NewStore()
	//assert.NilError(t, err, "Updating device failed"

	e2NodeID := topoapi.ID("13b4f7")

	// ToDo - should I take it from somewhere else?
	serviceModels := make(map[string]*topoapi.ServiceModelInfo)
	serviceModels["kpmv2"] = &topoapi.ServiceModelInfo{
		OID:  "1.3.6.1.4.1.53148.1.2.2.2",
		Name: "KPMv2",
	}

	// create or update E2 node entities
	err := st.CreateOrUpdateE2Node(context.Background(), e2NodeID, serviceModels)
	assert.NilError(t, err, "Updating device failed")

	// checking whether it is actually stored
	out, err := st.GetE2Relation(context.Background(), e2NodeID)
	assert.NilError(t, err, "Updating device failed")
	assert.Equal(t, out, e2NodeID)

	// m stands for channel manager
	// removing device from the store
	err = st.DeleteE2Relation(context.Background(), out)
	assert.NilError(t, err, "Removing device failed")

	// checking whether it was deleted
	res, _ := st.GetE2Relation(context.Background(), e2NodeID)
	//assert.NilError(t, err, "Updating device failed")
	assert.Equal(t, res, topoapi.ID(""))
}

//func (s *TestSuite) TestMultiE2Nodes(t *testing.T) {
//
//	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
//	defer cancel()
//	topoSdkClient, err := utils.NewTopoClient()
//	assert.NilError(t, err)
//	topoEventChan := make(chan topoapi.Event)
//	err = topoSdkClient.WatchE2Connections(ctx, topoEventChan)
//	assert.NilError(t, err)
//
//	nodeClient := utils.GetRansimNodeClient(t, sim)
//	assert.Assert(t, nodeClient != nil)
//
//	defaultNumNodes := utils.GetNumNodes(t, nodeClient)
//
//	for i := 0; i < numRequestedE2Nodes; i++ {
//		enbID := i + initialEnbID
//		createNodeRequest := &modelapi.CreateNodeRequest{
//			Node: &ransimtypes.Node{
//				GnbID:         ransimtypes.GnbID(enbID),
//				ServiceModels: serviceModels,
//				Controllers:   controllers,
//				CellNCGIs:     []ransimtypes.NCGI{},
//			},
//		}
//		e2node, err := nodeClient.CreateNode(ctx, createNodeRequest)
//		assert.NilError(t, err)
//		assert.Assert(t, e2node != nil)
//	}
//	numNodes := utils.GetNumNodes(t, nodeClient)
//	assert.Equal(t, numRequestedE2Nodes+defaultNumNodes, numNodes)
//
//	utils.CountTopoAddedOrNoneEvent(topoEventChan, numNodes)
//
//	e2nodes := utils.GetNodes(t, nodeClient)
//	for _, e2node := range e2nodes {
//		_, err = nodeClient.DeleteNode(ctx, &modelapi.DeleteNodeRequest{
//			GnbID: e2node.GnbID,
//		})
//		assert.NoError(t, err)
//	}
//
//	utils.CountTopoRemovedEvent(topoEventChan, numNodes)
//
//	numNodes = utils.GetNumNodes(t, nodeClient)
//	assert.Equal(t, 0, numNodes)
//	err = sim.Uninstall()
//	assert.NilError(t, err)
//}
