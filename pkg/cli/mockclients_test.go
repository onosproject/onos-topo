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

// Client Mocks
package cli

import (
	"context"
	"github.com/onosproject/onos-topo/api/device"
	"google.golang.org/grpc"
)

type mockDeviceServiceClient struct {
	test string
}

func (m *mockDeviceServiceClient) Add(ctx context.Context, request *device.AddRequest, opts ...grpc.CallOption) (*device.AddResponse, error) {
	addedDevice := request.Device
	return &device.AddResponse{Device: addedDevice}, nil // Just reflect it
}

func (m *mockDeviceServiceClient) Update(ctx context.Context, request *device.UpdateRequest, opts ...grpc.CallOption) (*device.UpdateResponse, error) {
	updatedDevice := request.Device
	return &device.UpdateResponse{Device: updatedDevice}, nil
}

func (m *mockDeviceServiceClient) Get(ctx context.Context, request *device.GetRequest, opts ...grpc.CallOption) (*device.GetResponse, error) {
	return &device.GetResponse{Device: generateDeviceData(1)[0]}, nil
}

func (m *mockDeviceServiceClient) List(ctx context.Context, in *device.ListRequest, opts ...grpc.CallOption) (device.DeviceService_ListClient, error) {
	return nil, nil
}

func (m *mockDeviceServiceClient) Remove(ctx context.Context, request *device.RemoveRequest, opts ...grpc.CallOption) (*device.RemoveResponse, error) {
	return &device.RemoveResponse{}, nil
}

// setUpMockClients sets up factories to create mocks of top level clients used by the CLI
func setUpMockClients() {
	device.DeviceServiceClientFactory = func(cc *grpc.ClientConn) device.DeviceServiceClient {
		return &mockDeviceServiceClient{test: ""}
	}
}
