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
	"fmt"

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

// Set :
func (s *Server) Set(ctx context.Context, request *topoapi.SetRequest) (*topoapi.SetResponse, error) {
	for _, object := range request.Objects {
		if err := s.ValidateObject(object); err != nil {
			return nil, err
		}
		if err := s.objectStore.Store(object); err != nil {
			return nil, err
		}
	}
	return &topoapi.SetResponse{}, nil
}

// Get :
func (s *Server) Get(ctx context.Context, request *topoapi.GetRequest) (*topoapi.GetResponse, error) {
	id := request.ID
	object, err := s.objectStore.Load(id)
	if err != nil {
		return nil, err
	} else if object == nil {
		log.Infof("Not found object %s", string(id))
		return nil, status.Error(codes.NotFound, string(id))
	}
	return &topoapi.GetResponse{
		Object: object,
	}, nil
}

// Delete ...
func (s *Server) Delete(ctx context.Context, request *topoapi.DeleteRequest) (*topoapi.DeleteResponse, error) {
	id := request.ID
	err := s.objectStore.Delete(id)
	if err != nil {
		return nil, err
	}
	return &topoapi.DeleteResponse{}, nil
}

// List ..
func (s *Server) List(request *topoapi.ListRequest, server topoapi.Topo_ListServer) error {
	ch := make(chan *topoapi.Object)
	if err := s.objectStore.List(ch); err != nil {
		return err
	}

	for object := range ch {
		err := server.Send(&topoapi.ListResponse{
			Object: object,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// Subscribe ...
func (s *Server) Subscribe(request *topoapi.SubscribeRequest, server topoapi.Topo_SubscribeServer) error {
	var watchOpts []WatchOption
	ch := make(chan *Event)

	if !request.Noreplay {
		watchOpts = append(watchOpts, WithReplay())
	}
	if err := s.objectStore.Watch(ch, watchOpts...); err != nil {
		return err
	}

	return s.Stream(server, ch)
}

// Stream ...
func (s *Server) Stream(server topoapi.Topo_SubscribeServer, ch chan *Event) error {
	for event := range ch {
		var t topoapi.Update_Type
		switch event.Type {
		case EventNone:
			t = topoapi.Update_UNSPECIFIED
		case EventInserted:
			t = topoapi.Update_INSERT
		case EventUpdated:
			t = topoapi.Update_MODIFY
		case EventRemoved:
			t = topoapi.Update_DELETE
		}

		subscribeResponse := &topo.SubscribeResponse{
			Update: &topo.Update{
				Type:   t,
				Object: event.Object,
			},
		}

		if err := server.Send(subscribeResponse); err != nil {
			return err
		}
	}
	return nil
}

// ValidateObject validates the given object
func (s *Server) ValidateObject(object *topoapi.Object) error {
	var kind *topo.Object
	var err error
	switch object.Type {
	case topo.Object_KIND:
	case topo.Object_ENTITY:
		if object.GetEntity().KindID != topo.NullID {
			kind, err = s.Load(object.GetEntity().KindID)
			if err != nil {
				return err
			}
		}
	case topo.Object_RELATION:
		kind, err = s.Load(object.GetRelation().KindID)
		if err != nil {
			return err
		}
		_, err := s.Load(object.GetRelation().SrcEntityID)
		if err != nil {
			return err
		}
		_, err = s.Load(object.GetRelation().TgtEntityID)
		if err != nil {
			return err
		}
	default:
		log.Infof("Invalid type %v", object)
	}

	if kind != nil && object.Type != topo.Object_KIND {
		if kind.Attributes != nil {
			for attrName := range object.Attributes {
				if _, ok := kind.Attributes[attrName]; !ok {
					return fmt.Errorf("Invalid attribute %s", attrName)
				}
			}
			for attrName, val := range kind.Attributes {
				if _, ok := object.Attributes[attrName]; !ok {
					object.Attributes[attrName] = val
				}
			}
		}
	}
	return nil
}

// Load ...
func (s *Server) Load(id topo.ID) (*topo.Object, error) {
	object, err := s.objectStore.Load(id)
	if err != nil {
		return nil, err
	} else if object == nil {
		log.Infof("Not found object %s", string(id))
		return nil, status.Error(codes.NotFound, string(id))
	}
	return object, nil
}
