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
	"github.com/onosproject/onos-topo/pkg/northbound"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NewService returns a new device Service
func NewService() (northbound.Service, error) {
	deviceStore, err := NewAtomixStore()
	if err != nil {
		return nil, err
	}
	return &Service{
		store: deviceStore,
	}, nil
}

// Service is a Service implementation for administration.
type Service struct {
	northbound.Service
	store Store
}

// Register registers the Service with the gRPC server.
func (s Service) Register(r *grpc.Server) {
	server := &Server{
		deviceStore: s.store,
	}
	RegisterDeviceServiceServer(r, server)
}

// Server implements the gRPC service for administrative facilities.
type Server struct {
	deviceStore Store
}

func (s *Server) Add(ctx context.Context, request *AddRequest) (*AddResponse, error) {
	device := request.Device
	if device == nil {
		return nil, status.Error(codes.InvalidArgument, "no device specified")
	} else if device.Metadata != nil && device.Metadata.Version != 0 {
		return nil, status.Error(codes.InvalidArgument, "device version is already set")
	}
	if err := s.deviceStore.Store(device); err != nil {
		return nil, err
	}
	return &AddResponse{
		Metadata: device.Metadata,
	}, nil
}

func (s *Server) Update(ctx context.Context, request *UpdateRequest) (*UpdateResponse, error) {
	device := request.Device
	if device == nil {
		return nil, status.Error(codes.InvalidArgument, "no device specified")
	} else if device.Metadata == nil || device.Metadata.Version == 0 {
		return nil, status.Error(codes.InvalidArgument, "device version not set")
	}
	if err := s.deviceStore.Store(device); err != nil {
		return nil, err
	}
	return &UpdateResponse{
		Metadata: device.Metadata,
	}, nil
}

func (s *Server) Get(ctx context.Context, request *GetRequest) (*GetResponse, error) {
	device, err := s.deviceStore.Load(request.DeviceId)
	if err != nil {
		return nil, err
	} else if device == nil {
		return nil, status.Error(codes.NotFound, "device not found")
	}
	return &GetResponse{
		Device: device,
	}, nil
}

func (s *Server) List(request *ListRequest, server DeviceService_ListServer) error {
	if request.Subscribe {
		ch := make(chan *Event)
		if err := s.deviceStore.Watch(ch); err != nil {
			return err
		}

		for event := range ch {
			var t ListResponse_Type
			switch event.Type {
			case EventNone:
				t = ListResponse_NONE
			case DeviceInserted:
				t = ListResponse_ADDED
			case DeviceUpdated:
				t = ListResponse_UPDATED
			case DeviceRemoved:
				t = ListResponse_REMOVED
			}
			err := server.Send(&ListResponse{
				Type:   t,
				Device: event.Device,
			})
			if err != nil {
				return err
			}
		}
	} else {
		ch := make(chan *Device)
		if err := s.deviceStore.List(ch); err != nil {
			return err
		}

		for device := range ch {
			err := server.Send(&ListResponse{
				Type:   ListResponse_NONE,
				Device: device,
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Server) Remove(ctx context.Context, request *RemoveRequest) (*RemoveResponse, error) {
	device := request.Device
	err := s.deviceStore.Delete(device)
	if err != nil {
		return nil, err
	}
	return &RemoveResponse{}, nil
}
