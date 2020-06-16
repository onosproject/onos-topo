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

	"github.com/onosproject/onos-lib-go/pkg/northbound"
	"github.com/onosproject/onos-topo/api/topo"
	topoapi "github.com/onosproject/onos-topo/api/topo"
	"google.golang.org/grpc"
)

//var log = logging.GetLogger("northbound", "topo")

// NewService returns a new topo Service
func NewService() (northbound.Service, error) {
	return &Service{}, nil
}

// Service is a Service implementation for administration.
type Service struct {
	northbound.Service
}

// Register registers the Service with the gRPC server.
func (s Service) Register(r *grpc.Server) {
	server := &Server{}
	topoapi.RegisterTopoServer(r, server)
}

// Server implements the gRPC service for administrative facilities.
type Server struct {
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
	return &topoapi.WriteResponse{}, nil
}

// Read :
func (s *Server) Read(ctx context.Context, request *topoapi.ReadRequest) (*topoapi.ReadResponse, error) {
	return &topoapi.ReadResponse{}, nil
}

// StreamChannel :
func (s *Server) StreamChannel(topo.Topo_StreamChannelServer) error {
	return nil
}
