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

// Package admin implements the northbound administrative gRPC service for the topology subsystem.
package device

import (
	"context"
	"github.com/onosproject/onos-topo/pkg/manager"
	"github.com/onosproject/onos-topo/pkg/northbound"
	"github.com/onosproject/onos-topo/pkg/northbound/proto"
	"github.com/onosproject/onos-topo/pkg/store"
	"google.golang.org/grpc"
)

// NewService returns a new device Service
func NewService(mgr *manager.Manager) northbound.Service {
	return &Service{
		mgr: mgr,
	}
}

// Service is a Service implementation for administration.
type Service struct {
	northbound.Service
	mgr *manager.Manager
}

// Register registers the Service with the gRPC server.
func (s Service) Register(r *grpc.Server) {
	server := &Server{
		deviceStore: s.mgr.DeviceStore,
	}
	proto.RegisterDeviceServiceServer(r, server)
}

// Server implements the gRPC service for administrative facilities.
type Server struct {
	deviceStore store.DeviceStore
}

func (s *Server) Add(ctx context.Context, request *proto.AddDeviceRequest) (*proto.AddDeviceResponse, error) {
	device := request.Device
	if err := s.deviceStore.Store(device); err != nil {
		return nil, err
	}
	return &proto.AddDeviceResponse{
		Metadata: device.Metadata,
	}, nil
}

func (s *Server) Update(ctx context.Context, request *proto.UpdateDeviceRequest) (*proto.UpdateDeviceResponse, error) {
	device := request.Device
	if err := s.deviceStore.Store(device); err != nil {
		return nil, err
	}
	return &proto.UpdateDeviceResponse{
		Metadata: device.Metadata,
	}, nil
}

func (s *Server) Get(ctx context.Context, request *proto.GetDeviceRequest) (*proto.GetDeviceResponse, error) {
	device, err := s.deviceStore.Load(request.DeviceId)
	if err != nil {
		return nil, err
	}
	return &proto.GetDeviceResponse{
		Device: device,
	}, nil
}

func (s *Server) List(request *proto.ListRequest, server proto.DeviceService_ListServer) error {
	if request.Subscribe {
		ch := make(chan *store.DeviceEvent)
		if err := s.deviceStore.Watch(ch); err != nil {
			return err
		}

		for event := range ch {
			var t proto.ListResponse_Type
			switch event.Type {
			case store.EventNone:
				t = proto.ListResponse_NONE
			case store.DeviceInserted:
				t = proto.ListResponse_ADDED
			case store.DeviceUpdated:
				t = proto.ListResponse_UPDATED
			case store.DeviceRemoved:
				t = proto.ListResponse_REMOVED
			}
			server.Send(&proto.ListResponse{
				Type:   t,
				Device: event.Device,
			})
		}
	} else {
		ch := make(chan *proto.Device)
		if err := s.deviceStore.List(ch); err != nil {
			return err
		}

		for device := range ch {
			server.Send(&proto.ListResponse{
				Type:   proto.ListResponse_NONE,
				Device: device,
			})
		}
	}
	return nil
}

func (s *Server) Remove(ctx context.Context, request *proto.RemoveDeviceRequest) (*proto.RemoveDeviceResponse, error) {
	err := s.deviceStore.Delete(request.Device)
	if err != nil {
		return nil, err
	}
	return &proto.RemoveDeviceResponse{}, nil
}
