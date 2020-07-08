// Copyright 2020-present Open Networking Foundation.
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

package bulk

import (
	"fmt"
	configlib "github.com/onosproject/onos-lib-go/pkg/config"
	"github.com/onosproject/onos-topo/api/topo"
)

var topoConfig *TopoConfig

// TopoConfig - the top level object
type TopoConfig struct {
	TopoEntities      []TopoEntity
	TopoRelationships []TopoRelationship
}

// TopoEntity - required to get around the "onoof" Obj
type TopoEntity struct {
	Ref   *topo.Reference
	Type  topo.Object_Type
	Obj   *topo.Object_Entity
	Attrs *topo.Attributes
}

// TopoEntityToTopoObject - convert to Object
func TopoEntityToTopoObject(topoEntity *TopoEntity) *topo.Object {
	return &topo.Object{
		Ref:   topoEntity.Ref,
		Type:  topoEntity.Type,
		Obj:   topoEntity.Obj,
		Attrs: topoEntity.Attrs,
	}
}

// TopoRelationship - required to get around the "onoof" Obj
type TopoRelationship struct {
	Ref   *topo.Reference
	Type  topo.Object_Type
	Obj   *topo.Object_Relationship
	Attrs *topo.Attributes
}

// TopoRelationshipToTopoObject - convert to Object
func TopoRelationshipToTopoObject(topoRelationship *TopoRelationship) *topo.Object {
	return &topo.Object{
		Ref:   topoRelationship.Ref,
		Type:  topoRelationship.Type,
		Obj:   topoRelationship.Obj,
		Attrs: topoRelationship.Attrs,
	}
}

// ClearTopo - reset the config - needed for tests
func ClearTopo() {
	topoConfig = nil
}

// GetTopoConfig gets the onos-topo configuration
func GetTopoConfig(location string) (TopoConfig, error) {
	if topoConfig == nil {
		topoConfig = &TopoConfig{}
		if err := configlib.LoadNamedConfig(location, topoConfig); err != nil {
			return TopoConfig{}, err
		}
		if err := TopoChecker(topoConfig); err != nil {
			return TopoConfig{}, err
		}
	}
	return *topoConfig, nil
}

// TopoChecker - check everything is within bounds
func TopoChecker(config *TopoConfig) error {
	if len(config.TopoEntities) == 0 {
		return fmt.Errorf("no entities found")
	}

	for _, entity := range config.TopoEntities {
		topoEntity := entity // pin
		if topoEntity.Type != topo.Object_ENTITY {
			return fmt.Errorf("unexpected type %v for TopoEntity", topoEntity.Type)
		} else if topoEntity.Ref == nil || topoEntity.Ref.GetID() == "" {
			return fmt.Errorf("empty ref for TopoEntity")
		}
	}

	for _, relationship := range config.TopoRelationships {
		topoRelationship := relationship // pin
		if topoRelationship.Type != topo.Object_RELATIONSHIP {
			return fmt.Errorf("unexpected type %v for TopoRelationship", topoRelationship.Type)
		} else if topoRelationship.Ref == nil || topoRelationship.Ref.GetID() == "" {
			return fmt.Errorf("empty ref for TopoRelationship")
		} else if topoRelationship.Obj.Relationship.SourceRef == nil ||
			topoRelationship.Obj.Relationship.SourceRef.ID == "" {
			return fmt.Errorf("empty source ref for TopoRelationship")
		} else if topoRelationship.Obj.Relationship.TargetRef == nil ||
			topoRelationship.Obj.Relationship.TargetRef.ID == "" {
			return fmt.Errorf("empty target ref for TopoRelationship")
		}
	}

	return nil
}
