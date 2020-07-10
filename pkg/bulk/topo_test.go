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
	"testing"

	"github.com/ghodss/yaml"
	"github.com/onosproject/onos-topo/api/topo"
	"gotest.tools/assert"
)

func Test_LoadConfig2(t *testing.T) {
	topoConfig = nil
	config, err := GetTopoConfig("topo-load-entities-example.yaml")
	assert.NilError(t, err, "Unexpected error loading topo entities")
	//assert.Equal(t, 2, len(config.TopoKinds), "Unexpected number of topo kinds")
	assert.Equal(t, 2, len(config.TopoEntities), "Unexpected number of topo entities")
	assert.Equal(t, 1, len(config.TopoRelations), "Unexpected number of topo relations")

	// TODO - check kinds
	tower1 := config.TopoEntities[0]
	assert.Equal(t, topo.Object_ENTITY, tower1.Type)
	assert.Equal(t, "E2Node", string(tower1.Obj.Entity.GetKind().ID))
	assert.Equal(t, topo.ID("315010-0001420"), tower1.Ref.GetID())
	address, ok := tower1.Attrs.GetAttrs()["address"]
	assert.Assert(t, ok, "error extracting address")
	assert.Equal(t, "ran-simulator:5152", address)

	rel1 := config.TopoRelations[0]
	assert.Equal(t, topo.Object_RELATION, rel1.Type)
	assert.Equal(t, "XnInterface", string(rel1.Obj.Relation.GetKind().ID))
	assert.Equal(t, topo.ID("315010-0001420"), rel1.Obj.Relation.GetSourceRef().GetID())
	assert.Equal(t, topo.ID("315010-0001421"), rel1.Obj.Relation.GetTargetRef().GetID())
	assert.Equal(t, topo.ID("rel1"), rel1.Ref.GetID())
	displayname, ok := rel1.Attrs.GetAttrs()["displayname"]
	assert.Assert(t, ok, "error extracting displayname")
	assert.Equal(t, "Tower 1 - Tower 2", displayname)

}

func Test_LoadConfig3(t *testing.T) {

	topoEntity1 := topo.Object{
		Ref:  &topo.Reference{ID: "entity1"},
		Type: topo.Object_ENTITY,
		Obj: &topo.Object_Entity{
			Entity: &topo.Entity{
				Kind: &topo.Reference{ID: "E2Node"},
			},
		},
		Attrs: &topo.Attributes{Attrs: make(map[string]string)},
	}
	topoEntity1.GetAttrs().GetAttrs()["test1"] = "testvalue1"
	topoEntity1.GetAttrs().GetAttrs()["test2"] = "testvalue2"

	topoEntity2 := topo.Object{
		Ref:  &topo.Reference{ID: "entity2"},
		Type: topo.Object_ENTITY,
		Obj: &topo.Object_Entity{
			Entity: &topo.Entity{
				Kind: &topo.Reference{ID: "E2Node"},
			},
		},
		Attrs: &topo.Attributes{Attrs: make(map[string]string)},
	}
	topoEntity2.GetAttrs().GetAttrs()["test3"] = "testvalue3"
	topoEntity2.GetAttrs().GetAttrs()["test4"] = "testvalue4"

	topoRelation1 := topo.Object{
		Ref:  &topo.Reference{ID: "relation1"},
		Type: topo.Object_RELATION,
		Obj: &topo.Object_Relation{
			Relation: &topo.Relation{
				Kind:      &topo.Reference{ID: "XnInterface"},
				SourceRef: topoEntity1.Ref,
				TargetRef: topoEntity2.Ref,
			}},
		Attrs: &topo.Attributes{Attrs: make(map[string]string)},
	}
	topoRelation1.GetAttrs().GetAttrs()["test3"] = "testvalue3"
	topoRelation1.GetAttrs().GetAttrs()["test4"] = "testvalue4"

	out, err := yaml.Marshal(topoEntity1)
	assert.NilError(t, err, "Unexpected error marshalling entity to YAML")
	t.Log("topoEntity1\n", string(out))
}
