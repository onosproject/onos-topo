// Copyright 2021-present Open Networking Foundation.
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

package northbound

import (
	"context"
	"fmt"
	"sort"

	"github.com/onosproject/onos-lib-go/pkg/errors"

	topoapi "github.com/onosproject/onos-api/go/onos/topo"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-lib-go/pkg/northbound"
	"github.com/onosproject/onos-topo/pkg/store"
	"google.golang.org/grpc"
)

var log = logging.GetLogger("northbound", "topo")

// NewService returns a new topo Service
func NewService(store store.Store) northbound.Service {
	return &Service{
		store: store,
	}
}

// Service is a Service implementation for administration.
type Service struct {
	store store.Store
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
	objectStore store.Store
}

// Create creates a new topology object
func (s *Server) Create(ctx context.Context, req *topoapi.CreateRequest) (*topoapi.CreateResponse, error) {
	log.Infof("Received CreateRequest %+v", req)
	object := req.Object
	err := s.objectStore.Create(ctx, object)
	if err != nil {
		log.Warnf("CreateRequest %+v failed: %v", req, err)
		return nil, errors.Status(err).Err()
	}
	res := &topoapi.CreateResponse{
		Object: object,
	}
	log.Infof("Sending CreateResponse %+v", res)
	return res, nil
}

// Get retrieves the specified topology object
func (s *Server) Get(ctx context.Context, req *topoapi.GetRequest) (*topoapi.GetResponse, error) {
	log.Infof("Received GetRequest %+v", req)
	object, err := s.objectStore.Get(ctx, req.ID)
	if err != nil {
		log.Warnf("GetRequest %+v failed: %v", req, err)
		return nil, errors.Status(err).Err()
	}
	res := &topoapi.GetResponse{
		Object: object,
	}
	log.Infof("Sending GetResponse %+v", res)
	return res, nil
}

// Update creates an existing topology object
func (s *Server) Update(ctx context.Context, req *topoapi.UpdateRequest) (*topoapi.UpdateResponse, error) {
	log.Infof("Received UpdateRequest %+v", req)
	err := s.objectStore.Update(ctx, req.Object)
	if err != nil {
		log.Warnf("UpdateRequest %+v failed: %v", req, err)
		return nil, errors.Status(err).Err()
	}
	res := &topoapi.UpdateResponse{
		Object: req.Object,
	}
	log.Infof("Sending UpdateResponse %+v", res)
	return res, nil
}

// Delete removes the specified topology object
func (s *Server) Delete(ctx context.Context, req *topoapi.DeleteRequest) (*topoapi.DeleteResponse, error) {
	log.Infof("Received DeleteRequest %+v", req)
	err := s.objectStore.Delete(ctx, req.ID)
	if err != nil {
		log.Warnf("DeleteRequest %+v failed: %v", req, err)
		return nil, errors.Status(err).Err()
	}
	res := &topoapi.DeleteResponse{}
	log.Infof("Sending DeleteResponse %+v", res)
	return res, nil
}

// List returns list of all objects
func (s *Server) List(ctx context.Context, req *topoapi.ListRequest) (*topoapi.ListResponse, error) {
	log.Infof("Received ListRequest %+v", req)
	objects, err := s.objectStore.List(ctx, req.Filters)
	if err != nil {
		log.Warnf("ListRequest %+v failed: %v", req, err)
		return nil, errors.Status(err).Err()
	}

	if req.SortOrder != topoapi.SortOrder_UNORDERED {
		sort.Slice(objects, func(i, j int) bool {
			if req.SortOrder == topoapi.SortOrder_DESCENDING {
				return objects[i].ID > objects[j].ID
			}
			return objects[i].ID < objects[j].ID
		})
	}

	res := &topoapi.ListResponse{
		Objects: objects,
	}
	log.Infof("Sending ListResponse %+v", res)
	return res, nil
}

// Watch streams topology changes
func (s *Server) Watch(req *topoapi.WatchRequest, server topoapi.Topo_WatchServer) error {
	log.Infof("Received WatchRequest %+v", req)
	var watchOpts []store.WatchOption
	if !req.Noreplay {
		watchOpts = append(watchOpts, store.WithReplay())
	}

	ch := make(chan topoapi.Event)
	if err := s.objectStore.Watch(server.Context(), ch, req.Filters, watchOpts...); err != nil {
		log.Warnf("WatchTerminationsRequest %+v failed: %v", req, err)
		return errors.Status(err).Err()
	}

	return s.Stream(server, ch)
}

// Stream is the ongoing stream for WatchTerminations request
func (s *Server) Stream(server topoapi.Topo_WatchServer, ch chan topoapi.Event) error {
	for event := range ch {
		res := &topoapi.WatchResponse{
			Event: event,
		}

		log.Infof("Sending WatchResponse %+v", res)
		if err := server.Send(res); err != nil {
			log.Warnf("WatchResponse %+v failed: %v", res, err)
			return err
		}
	}
	return nil
}

// ValidateObject validates the given object
func (s *Server) ValidateObject(ctx context.Context, object *topoapi.Object) error {
	var kind *topoapi.Object
	var err error
	switch object.Type {
	case topoapi.Object_KIND:
	case topoapi.Object_ENTITY:
		if object.GetEntity().KindID != topoapi.NullID {
			kind, err = s.objectStore.Get(ctx, object.GetEntity().KindID)
			if err != nil {
				return err
			}
		}
	case topoapi.Object_RELATION:
		kind, err = s.objectStore.Get(ctx, object.GetRelation().KindID)
		if err != nil {
			return err
		}
		_, err := s.objectStore.Get(ctx, object.GetRelation().SrcEntityID)
		if err != nil {
			return err
		}
		_, err = s.objectStore.Get(ctx, object.GetRelation().TgtEntityID)
		if err != nil {
			return err
		}
	default:
		log.Infof("Invalid type %v", object)
	}

	if kind != nil && object.Type != topoapi.Object_KIND {
		if kind.Aspects != nil {
			for attrName := range object.Aspects {
				if _, ok := kind.Aspects[attrName]; !ok {
					return fmt.Errorf("Invalid attribute %s", attrName)
				}
			}
			for attrName, val := range kind.Aspects {
				if _, ok := object.Aspects[attrName]; !ok {
					object.Aspects[attrName] = val
				}
			}
		}
	}
	return nil
}
