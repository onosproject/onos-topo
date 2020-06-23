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

// Package topo :
package topo

import (
	"context"

	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-lib-go/pkg/northbound"
	"github.com/onosproject/onos-topo/api/topo"
	topoapi "github.com/onosproject/onos-topo/api/topo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var log = logging.GetLogger("northbound", "topo")

// NewService returns a new topo Service
func NewService() (northbound.Service, error) {
	objectStore, err := NewAtomixStore()
	if err != nil {
		return nil, err
	}
	return &Service{
		store: objectStore,
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
		objectStore: s.store,
	}
	topoapi.RegisterTopoServer(r, server)
}

// Server implements the gRPC service for administrative facilities.
type Server struct {
	objectStore Store
}

// TopoClientFactory : Default TopoClient creation.
var TopoClientFactory = func(cc *grpc.ClientConn) topoapi.TopoClient {
	return topoapi.NewTopoClient(cc)
}

// CreateTopoClient creates and returns a new topo entity client
func CreateTopoClient(cc *grpc.ClientConn) topoapi.TopoClient {
	return TopoClientFactory(cc)
}

// ValidateEntity validates the given entity
func ValidateEntity(entity *topoapi.Entity) error {
	return nil
}

// Write :
func (s *Server) Write(ctx context.Context, request *topoapi.WriteRequest) (*topoapi.WriteResponse, error) {
	for _, update := range request.Updates {
		object := update.Object
		switch update.Type {
		case topo.Update_INSERT:
			log.Infof("Insert object %v", object)
			if err := s.objectStore.Store(object); err != nil {
				return nil, err
			}
		default:
			log.Infof("Invalid type %v", object)
		}
	}
	return &topoapi.WriteResponse{}, nil
}

// Read :
func (s *Server) Read(ctx context.Context, request *topoapi.ReadRequest) (*topoapi.ReadResponse, error) {
	var objects []*topo.Object

	for _, ref := range request.Refs {
		id := ref.ID
		object, err := s.objectStore.Load(id)
		if err != nil {
			return nil, err
		} else if object == nil {
			log.Infof("Not found object %s", string(id))
			return nil, status.Error(codes.NotFound, string(id))
		}
		log.Infof("Read object %v", object)
		objects = append(objects, object)
	}

	return &topoapi.ReadResponse{
		Objects: objects,
	}, nil
}

// StreamChannel :
func (s *Server) StreamChannel(topo.Topo_StreamChannelServer) error {
	return nil
}
