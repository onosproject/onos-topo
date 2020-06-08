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

package types

import "errors"

// ID is an Entity ID
type ID string

// EntityType is an Entity type
type EntityType string

const (
	// EtE2Interface represent an E2 Interface
	EtE2Interface EntityType = "ET_E2_INTERFACE"
)

// IsEntityTypeValid validates Entity Type
func (et EntityType) IsEntityTypeValid() error {
	switch et {
	case EtE2Interface:
		return nil
	}
	return errors.New("Inalid entity type")
}

// Role is an Entity role
type Role string

// Revision is the Entity revision number
type Revision uint64
