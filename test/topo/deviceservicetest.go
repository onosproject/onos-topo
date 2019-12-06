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
	"github.com/onosproject/onos-test/pkg/onit/env"
	"github.com/onosproject/onos-topo/api/device"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

// TestDeviceService : test
func (s *TestSuite) TestDeviceService(t *testing.T) {
	conn, err := env.Topo().Connect()
	assert.NoError(t, err)
	defer conn.Close()

	client := device.NewDeviceServiceClient(conn)

	list, err := client.List(context.Background(), &device.ListRequest{})
	assert.NoError(t, err)

	count := 0
	for {
		_, err := list.Recv()
		if err == io.EOF {
			break
		}
		assert.NoError(t, err)
		count++
	}

	assert.Equal(t, 0, count)

	events := make(chan *device.ListResponse)
	go func() {
		list, err := client.List(context.Background(), &device.ListRequest{
			Subscribe: true,
		})
		assert.NoError(t, err)

		for {
			response, err := list.Recv()
			if err != nil {
				break
			}
			events <- response
		}
	}()

	addResponse, err := client.Add(context.Background(), &device.AddRequest{
		Device: &device.Device{
			ID:      "test1",
			Type:    "Stratum",
			Address: "device-test1:5000",
			Target:  "device-test1",
			Version: "1.0.0",
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, device.ID("test1"), addResponse.Device.ID)
	assert.NotEqual(t, device.Revision(0), addResponse.Device.Revision)

	getResponse, err := client.Get(context.Background(), &device.GetRequest{
		ID: "test1",
	})
	assert.NoError(t, err)

	assert.Equal(t, device.ID("test1"), getResponse.Device.ID)
	assert.Equal(t, addResponse.Device.Revision, getResponse.Device.Revision)

	eventResponse := <-events
	assert.Equal(t, device.ListResponse_ADDED, eventResponse.Type)
	assert.Equal(t, device.ID("test1"), eventResponse.Device.ID)
	assert.Equal(t, addResponse.Device.Revision, eventResponse.Device.Revision)

	list, err = client.List(context.Background(), &device.ListRequest{})
	assert.NoError(t, err)
	for {
		response, err := list.Recv()
		if err == io.EOF {
			break
		}
		assert.NoError(t, err)
		assert.Equal(t, device.ListResponse_NONE, response.Type)
		assert.Equal(t, device.ID("test1"), response.Device.ID)
		assert.Equal(t, addResponse.Device.Revision, response.Device.Revision)
		count++
	}
	assert.Equal(t, 1, count)

	removeResponse, err := client.Remove(context.Background(), &device.RemoveRequest{
		Device: getResponse.Device,
	})
	assert.NoError(t, err)
	assert.NotNil(t, removeResponse)

	eventResponse = <-events
	assert.Equal(t, device.ListResponse_REMOVED, eventResponse.Type)
	assert.Equal(t, device.ID("test1"), eventResponse.Device.ID)
	assert.Equal(t, addResponse.Device.Revision, eventResponse.Device.Revision)
}
