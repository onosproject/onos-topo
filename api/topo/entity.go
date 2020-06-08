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

import "errors"

// EntityKind ...
type EntityKind string

// EntityID is a unique ID used as primary key for entities
type EntityID string

const (
	// EKE2Interface represent an 'E2 Interface' Entity etype
	EKE2Interface EntityKind = "ET_E2_INTERFACE"
)

// Entity ...
type Entity struct {
	// kind is the ...
	kind EntityKind

	// id is the entity's UUID
	id EntityID

	attr map[AttrKind][]AttrVal

	// rkContains stores the RKCONTAINS relationship kind
	rkContains map[EntityKind][]EntityID
}

// IsEntityKindValid validates Entity Type
func (ek EntityKind) IsEntityKindValid() error {
	switch ek {
	case EKE2Interface:
		return nil
	}
	return errors.New("Inalid entity type")
}
