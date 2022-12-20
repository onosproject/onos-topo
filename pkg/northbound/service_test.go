// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package northbound

import (
	"context"
	"github.com/atomix/go-sdk/pkg/primitive"
	"github.com/atomix/go-sdk/pkg/test"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"net"
	"sync"
	"testing"

	topoapi "github.com/onosproject/onos-api/go/onos/topo"
	"github.com/onosproject/onos-lib-go/pkg/northbound"
	"github.com/onosproject/onos-topo/pkg/store"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

var lis *bufconn.Listener

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func newTestService(client primitive.Client) (northbound.Service, error) {
	store, err := store.NewAtomixStore(client)
	if err != nil {
		return nil, err
	}
	return &Service{
		store: store,
	}, nil
}

func createServerConnection(t *testing.T, client primitive.Client) *grpc.ClientConn {
	lis = bufconn.Listen(1024 * 1024)
	s, err := newTestService(client)
	assert.NoError(t, err)
	assert.NotNil(t, s)
	server := grpc.NewServer()
	s.Register(server)

	go func() {
		if err := server.Serve(lis); err != nil {
			assert.NoError(t, err, "Server exited with error: %v", err)
		}
	}()

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	return conn
}

func TestServiceBasics(t *testing.T) {
	cluster := test.NewClient()
	defer cluster.Close()

	conn := createServerConnection(t, cluster)
	client := topoapi.NewTopoClient(conn)

	_, err := client.Create(context.Background(), &topoapi.CreateRequest{
		Object: &topoapi.Object{
			ID:   "1",
			Type: topoapi.Object_ENTITY,
		},
	})
	assert.NoError(t, err)

	_, err = client.Create(context.Background(), &topoapi.CreateRequest{
		Object: &topoapi.Object{
			ID:   "2",
			Type: topoapi.Object_ENTITY,
		},
	})
	assert.NoError(t, err)

	gres, err := client.Get(context.Background(), &topoapi.GetRequest{
		ID: "1",
	})
	assert.NoError(t, err)
	assert.Equal(t, topoapi.ID("1"), gres.Object.ID)

	res, err := client.List(context.Background(), &topoapi.ListRequest{})
	assert.NoError(t, err)
	assert.Condition(t, func() bool {
		return len(res.Objects) == 2 &&
			(res.Objects[0].ID == "1" || res.Objects[1].ID == "1")
	})

	obj := gres.Object
	err = obj.SetAspect(&topoapi.Location{Lat: 3.14, Lng: 6.28})
	assert.NoError(t, err)
	ures, err := client.Update(context.Background(), &topoapi.UpdateRequest{
		Object: obj,
	})
	assert.NoError(t, err)
	assert.NotNil(t, ures)

	obj = ures.Object
	loc := &topoapi.Location{}
	err = obj.GetAspect(loc)
	assert.NoError(t, err)
	assert.Equal(t, 6.28, loc.Lng)

	qst, err := client.Query(context.Background(), &topoapi.QueryRequest{})
	assert.NoError(t, err)
	_, err = qst.Recv()
	assert.NoError(t, err)
	_, err = qst.Recv()
	assert.NoError(t, err)
	_, err = qst.Recv()
	assert.True(t, err == io.EOF)

	_, err = client.Delete(context.Background(), &topoapi.DeleteRequest{
		ID:       obj.ID,
		Revision: obj.Revision,
	})
	assert.NoError(t, err)

	qst, err = client.Query(context.Background(), &topoapi.QueryRequest{})
	assert.NoError(t, err)
	qres, err := qst.Recv()
	assert.NoError(t, err)
	assert.Condition(t, func() bool {
		return qres.Object.ID == "2"
	})
	qres, err = qst.Recv()
	assert.True(t, err == io.EOF)

	res, err = client.List(context.Background(), &topoapi.ListRequest{})
	assert.NoError(t, err)
	assert.Condition(t, func() bool {
		return len(res.Objects) == 1 && res.Objects[0].ID == "2"
	})
}

func TestWatchBasics(t *testing.T) {
	cluster := test.NewClient()
	defer cluster.Close()

	conn := createServerConnection(t, cluster)
	client := topoapi.NewTopoClient(conn)

	cres, err := client.Create(context.Background(), &topoapi.CreateRequest{
		Object: &topoapi.Object{
			ID:   "1",
			Type: topoapi.Object_ENTITY,
		},
	})
	assert.NoError(t, err)
	obj := cres.Object

	res, err := client.Watch(context.Background(), &topoapi.WatchRequest{})
	assert.NoError(t, err)

	var wg sync.WaitGroup
	wg.Add(1)

	var pause sync.WaitGroup
	pause.Add(1)
	go func() {
		e, err := res.Recv()
		assert.NoError(t, err)
		assert.Equal(t, topoapi.EventType_NONE, e.Event.Type)
		assert.Equal(t, topoapi.ID("1"), e.Event.Object.ID)
		pause.Done()

		e, err = res.Recv()
		assert.NoError(t, err)
		assert.Equal(t, topoapi.EventType_ADDED, e.Event.Type)
		assert.Equal(t, topoapi.ID("2"), e.Event.Object.ID)

		e, err = res.Recv()
		assert.NoError(t, err)
		assert.Equal(t, topoapi.EventType_REMOVED, e.Event.Type)
		assert.Equal(t, topoapi.ID("1"), e.Event.Object.ID)

		wg.Done()
	}()

	// Pause before adding a new item to validate that existing items are processed first
	pause.Wait()
	_, err = client.Create(context.Background(), &topoapi.CreateRequest{
		Object: &topoapi.Object{
			ID:   "2",
			Type: topoapi.Object_ENTITY,
		},
	})
	assert.NoError(t, err)

	_, err = client.Delete(context.Background(), &topoapi.DeleteRequest{
		ID:       obj.ID,
		Revision: obj.Revision,
	})
	assert.NoError(t, err)

	wg.Wait()
}

func TestBadIDAdd(t *testing.T) {
	cluster := test.NewClient()
	defer cluster.Close()

	conn := createServerConnection(t, cluster)
	client := topoapi.NewTopoClient(conn)

	_, err := client.Create(context.Background(), &topoapi.CreateRequest{
		Object: &topoapi.Object{},
	})
	assert.Error(t, err)
}

func TestBadTypeAdd(t *testing.T) {
	cluster := test.NewClient()
	defer cluster.Close()

	conn := createServerConnection(t, cluster)
	client := topoapi.NewTopoClient(conn)

	_, err := client.Create(context.Background(), &topoapi.CreateRequest{
		Object: &topoapi.Object{ID: "foo"},
	})
	assert.Error(t, err)
}

func TestBadRemove(t *testing.T) {
	cluster := test.NewClient()
	defer cluster.Close()

	conn := createServerConnection(t, cluster)
	client := topoapi.NewTopoClient(conn)

	_, err := client.Delete(context.Background(), &topoapi.DeleteRequest{})
	assert.Error(t, err)
}

func TestSort(t *testing.T) {
	cluster := test.NewClient()
	defer cluster.Close()

	conn := createServerConnection(t, cluster)
	client := topoapi.NewTopoClient(conn)

	_, err := client.Create(context.Background(), &topoapi.CreateRequest{
		Object: &topoapi.Object{
			ID:   "a",
			Type: topoapi.Object_ENTITY,
		},
	})
	assert.NoError(t, err)

	_, err = client.Create(context.Background(), &topoapi.CreateRequest{
		Object: &topoapi.Object{
			ID:   "b",
			Type: topoapi.Object_ENTITY,
		},
	})
	assert.NoError(t, err)

	_, err = client.Create(context.Background(), &topoapi.CreateRequest{
		Object: &topoapi.Object{
			ID:   "c",
			Type: topoapi.Object_ENTITY,
		},
	})
	assert.NoError(t, err)

	res, err := client.List(context.Background(), &topoapi.ListRequest{
		SortOrder: topoapi.SortOrder_ASCENDING,
	})
	assert.NoError(t, err)
	assert.Condition(t, func() bool {
		return len(res.Objects) == 3 &&
			(res.Objects[0].ID == "a" && res.Objects[1].ID == "b" && res.Objects[2].ID == "c")
	})
	res, err = client.List(context.Background(), &topoapi.ListRequest{
		SortOrder: topoapi.SortOrder_DESCENDING,
	})
	assert.NoError(t, err)
	assert.Condition(t, func() bool {
		return len(res.Objects) == 3 &&
			(res.Objects[0].ID == "c" && res.Objects[1].ID == "b" && res.Objects[2].ID == "a")
	})
}
