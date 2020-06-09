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
	topoapi "github.com/onosproject/onos-topo/api/topo"
	"google.golang.org/grpc"

	"time"
)

var log = logging.GetLogger("northbound", "topo")

const (
	defaultTimeout       = 5 * time.Second
	entityNamePattern    = `^[a-zA-Z0-9\-:_]{4,40}$`
	entityAddressPattern = `^[a-zA-Z0-9\-_\.]+:[0-9]+$`
	entityVersionPattern = `^(\d+(\.\d+){2,3})$`
	entityAttrKeyPattern = `^[a-zA-Z0-9\-_\.]{1,40}$`
	displayNameMaxLength = 80
)

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
	topoapi.RegisterEntityServiceServer(r, server)
}

// Server implements the gRPC service for administrative facilities.
type Server struct {
}

// EntityServiceClientFactory : Default EntityServiceClient creation.
var EntityServiceClientFactory = func(cc *grpc.ClientConn) topoapi.EntityServiceClient {
	return topoapi.NewEntityServiceClient(cc)
}

// CreateEntityServiceClient creates and returns a new topo entity client
func CreateEntityServiceClient(cc *grpc.ClientConn) topoapi.EntityServiceClient {
	return EntityServiceClientFactory(cc)
}

// ValidateEntity validates the given entity
func ValidateEntity(entity *topoapi.Entity) error {
	return nil
}

// Add :
func (s *Server) Add(ctx context.Context, request *topoapi.AddRequest) (*topoapi.AddResponse, error) {
	return &topoapi.AddResponse{}, nil
}

// Update :
func (s *Server) Update(ctx context.Context, request *topoapi.UpdateRequest) (*topoapi.UpdateResponse, error) {
	return &topoapi.UpdateResponse{}, nil
}

// Get :
func (s *Server) Get(ctx context.Context, request *topoapi.GetRequest) (*topoapi.GetResponse, error) {
	return &topoapi.GetResponse{}, nil
}

// List :
func (s *Server) List(request *topoapi.ListRequest, server topoapi.EntityService_ListServer) error {
	return nil
}

// Remove :
func (s *Server) Remove(ctx context.Context, request *topoapi.RemoveRequest) (*topoapi.RemoveResponse, error) {
	return &topoapi.RemoveResponse{}, nil
}
