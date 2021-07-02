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

//var (
//	initialEnbID  = 155000
//	serviceModels = []string{"kpm2", "rcpre2"}
//	controllers   = []string{"e2t-1"}
//)

var log = logging.GetLogger("topo")

func (s *TestSuite) TestAddRemoveDevice(t *testing.T) {

	conn, err := utils.CreateConnection()
	client := topoapi.NewTopoClient(conn)

	_, err = client.Create(context.Background(), &topoapi.CreateRequest{
		Object: &topoapi.Object{
			ID:   "1",
			Type: topoapi.Object_ENTITY,
		},
	})
	assert.NilError(t, err)

	_, err = client.Create(context.Background(), &topoapi.CreateRequest{
		Object: &topoapi.Object{
			ID:   "2",
			Type: topoapi.Object_ENTITY,
		},
	})
	assert.NilError(t, err)

	gres, err := client.Get(context.Background(), &topoapi.GetRequest{
		ID: "1",
	})
	assert.NilError(t, err)
	assert.Equal(t, topoapi.ID("1"), gres.Object.ID)

	res, err := client.List(context.Background(), &topoapi.ListRequest{})
	assert.NilError(t, err)
	assert.Equal(t, len(res.Objects) == 2 &&
		(res.Objects[0].ID == "1" || res.Objects[1].ID == "1"), true)
	//assert.Condition(t, func() bool {
	//	return len(res.Objects) == 2 &&
	//		(res.Objects[0].ID == "1" || res.Objects[1].ID == "1")
	//})

	obj := gres.Object
	//obj.Attributes = make(map[string]string)
	//obj.Attributes["foo"] = "bar"
	ures, err := client.Update(context.Background(), &topoapi.UpdateRequest{
		Object: obj,
	})
	assert.NilError(t, err)
	assert.Assert(t, ures != nil)
	//assert.Equal(t, ures.Object.Attributes["foo"], "bar")

	_, err = client.Delete(context.Background(), &topoapi.DeleteRequest{
		ID: "1",
	})
	assert.NilError(t, err)

	res, err = client.List(context.Background(), &topoapi.ListRequest{})
	assert.NilError(t, err)
	assert.Equal(t, len(res.Objects) == 1 && res.Objects[0].ID == "2", true)
	//assert.Condition(t, func() bool {
	//	return len(res.Objects) == 1 && res.Objects[0].ID == "2"
	//})
}
