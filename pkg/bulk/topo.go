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
	TopoKinds     []TopoKind
	TopoEntities  []TopoEntity
	TopoRelations []TopoRelation
}

// TopoKind - required to get around the "oneof" Obj
type TopoKind struct {
	Ref   *topo.Reference
	Type  topo.Object_Type
	Obj   *topo.Object_Kind
	Attrs *topo.Attributes
}

// TopoKindToTopoObject - convert to Object
func TopoKindToTopoObject(topoKind *TopoKind) *topo.Object {
	return &topo.Object{
		Ref:   topoKind.Ref,
		Type:  topoKind.Type,
		Obj:   topoKind.Obj,
		Attrs: topoKind.Attrs,
	}
}

// TopoEntity - required to get around the "oneof" Obj
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

// TopoRelation - required to get around the "oneof" Obj
type TopoRelation struct {
	Ref   *topo.Reference
	Type  topo.Object_Type
	Obj   *topo.Object_Relation
	Attrs *topo.Attributes
}

// TopoRelationToTopoObject - convert to Object
func TopoRelationToTopoObject(topoRelation *TopoRelation) *topo.Object {
	return &topo.Object{
		Ref:   topoRelation.Ref,
		Type:  topoRelation.Type,
		Obj:   topoRelation.Obj,
		Attrs: topoRelation.Attrs,
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
	if len(config.TopoKinds) == 0 {
		return fmt.Errorf("no kinds found")
	}

	for _, kind := range config.TopoKinds {
		topoKind := kind // pin
		if topoKind.Type != topo.Object_KIND {
			return fmt.Errorf("unexpected type %v for TopoKind", topoKind.Type)
		} else if topoKind.Ref == nil || topoKind.Ref.GetID() == "" {
			return fmt.Errorf("empty ref for TopoKind")
		} else if topoKind.Obj.Kind.GetName() == "" {
			return fmt.Errorf("empty name for TopoKind")
		}
	}

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

	for _, relation := range config.TopoRelations {
		topoRelation := relation // pin
		if topoRelation.Type != topo.Object_RELATION {
			return fmt.Errorf("unexpected type %v for TopoRelation", topoRelation.Type)
		} else if topoRelation.Ref == nil || topoRelation.Ref.GetID() == "" {
			return fmt.Errorf("empty ref for TopoRelation")
		} else if topoRelation.Obj.Relation.SourceRef == nil ||
			topoRelation.Obj.Relation.SourceRef.ID == "" {
			return fmt.Errorf("empty source ref for TopoRelation")
		} else if topoRelation.Obj.Relation.TargetRef == nil ||
			topoRelation.Obj.Relation.TargetRef.ID == "" {
			return fmt.Errorf("empty target ref for TopoRelation")
		}
	}

	return nil
}
