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

package device

import (
	"context"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"io"
	"net"
	"testing"
)

type testListServer struct {
	grpc.ServerStream
	ch chan *ListResponse
}

func (s *testListServer) Send(m *ListResponse) error {
	s.ch <- m
	return nil
}

func TestLocalServer(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	s := grpc.NewServer()

	RegisterDeviceServiceServer(s, &Server{
		deviceStore: NewLocalStore(),
	})

	go func() {
		if err := s.Serve(lis); err != nil {
			t.Fatalf("Server exited with error %v", err)
		}
	}()

	dialer := func(ctx context.Context, address string) (net.Conn, error) {
		return lis.Dial()
	}

	conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(dialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet %v", err)
	}

	client := NewDeviceServiceClient(conn)

	addResponse, err := client.Add(context.Background(), &AddRequest{
		Device: &Device{
			ID:      ID("foo"),
			Address: "foo:1234",
		},
	})
	assert.NoError(t, err)
	assert.NotEqual(t, Revision(0), addResponse.Device.Revision)

	getResponse, err := client.Get(context.Background(), &GetRequest{
		ID: ID("foo"),
	})
	assert.NoError(t, err)
	assert.Equal(t, ID("foo"), getResponse.Device.ID)
	assert.Equal(t, addResponse.Device.Revision, getResponse.Device.Revision)
	assert.Equal(t, "foo:1234", getResponse.Device.Address)

	list, err := client.List(context.Background(), &ListRequest{})
	for {
		listResponse, err := list.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("list failed with error %v", err)
		}
		assert.Equal(t, ID("foo"), listResponse.Device.ID)
		assert.Equal(t, addResponse.Device.Revision, listResponse.Device.Revision)
		assert.Equal(t, "foo:1234", listResponse.Device.Address)
	}

	subscribe, err := client.List(context.Background(), &ListRequest{
		Subscribe: true,
	})

	eventCh := make(chan *ListResponse)
	go func() {
		for {
			subscribeResponse, err := subscribe.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Fatalf("subscribe failed with error %v", err)
			}
			eventCh <- subscribeResponse
		}
	}()

	listResponse := <-eventCh
	assert.Equal(t, ListResponse_NONE, listResponse.Type)
	assert.Equal(t, ID("foo"), listResponse.Device.ID)
	assert.Equal(t, "foo:1234", listResponse.Device.Address)

	addResponse, err = client.Add(context.Background(), &AddRequest{
		Device: &Device{
			ID:      ID("bar"),
			Address: "bar:1234",
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, ID("bar"), addResponse.Device.ID)
	assert.Equal(t, "bar:1234", addResponse.Device.Address)
	assert.NotEqual(t, Revision(0), addResponse.Device.Revision)

	listResponse = <-eventCh
	assert.Equal(t, ListResponse_ADDED, listResponse.Type)
	assert.Equal(t, ID("bar"), listResponse.Device.ID)
	assert.Equal(t, "bar:1234", listResponse.Device.Address)

	_, err = client.Remove(context.Background(), &RemoveRequest{
		Device: getResponse.Device,
	})
	assert.NoError(t, err)

	listResponse = <- eventCh
	assert.Equal(t, ListResponse_REMOVED, listResponse.Type)
	assert.Equal(t, ID("foo"), listResponse.Device.ID)
	assert.Equal(t, "foo:1234", listResponse.Device.Address)
}
