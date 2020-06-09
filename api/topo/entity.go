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
// WITHOUT WARRANTIES OR CONDITIONS OF ANY etype, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package topo

import (
	"errors"

	grpc "google.golang.org/grpc"
)

// EntityKind represents an entity's "kind" or "type"
type EntityKind string

const (
	// EKE2Interface represent an 'E2 Interface' Entity etype
	EKE2Interface EntityKind = "ET_E2_INTERFACE"
)

// EntityID is a unique ID used as primary key for entities
type EntityID string

// Entity represent "things"
type Entity struct {
	// Entities have "kinds" or "types"
	kind EntityKind

	// id is a opaque universally unique identifiers (UUID) used as the entity's primary key
	id EntityID

	// attr maps the attributes for this entity. Each entity has a set of attributes.
	attr map[AttrKind][]AttrVal

	// rkContains maps the CONTAINS relationship for this entity
	// An entity can "contain" other entities, e.g. a switch contains ports.
	rkContains map[EntityKind][]EntityID
}

// IsEntityKindValid validates EntityKind
func (ek EntityKind) IsEntityKindValid() error {
	switch ek {
	case EKE2Interface:
		return nil
	}
	return errors.New("Inalid entity type")
}

// EntityServiceClientFactory : Default EntityServiceClient creation.
var EntityServiceClientFactory = func(cc *grpc.ClientConn) EntityServiceClient {
	return NewEntityServiceClient(cc)
}

// CreateEntityServiceClient creates and returns a new topo device client
func CreateEntityServiceClient(cc *grpc.ClientConn) EntityServiceClient {
	return EntityServiceClientFactory(cc)
}
