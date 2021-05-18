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

package topo

import topoapi "github.com/onosproject/onos-api/go/onos/topo"

func match(object *topoapi.Object, filters *topoapi.Filters) bool {
	return filters == nil || (matchKinds(object, filters.KindFilters) && matchLabels(object, filters.LabelFilters))
}

func matchLabels(object *topoapi.Object, filters []*topoapi.Filter) bool {
	for _, filter := range filters {
		if !matchLabel(object, filter) {
			return false
		}
	}
	return true
}

func matchKinds(object *topoapi.Object, filters []*topoapi.Filter) bool {
	if len(filters) > 0 {
		return matchKind(object, filters[0])
	}
	return true
}

func matchLabel(object *topoapi.Object, filter *topoapi.Filter) bool {
	eqo := filter.GetEqual_()
	if eqo != nil {
		return object.Labels[filter.Key] == eqo.Value
	}
	igo := filter.GetIn()
	if igo != nil {
		actual := object.Labels[filter.Key]
		for _, v := range igo.Values {
			if v == actual {
				return true
			}
		}
		return false
	}
	ngo := filter.GetNot()
	if ngo != nil {
		return !matchLabel(object, ngo.Inner)
	}
	return false
}

func matchKind(object *topoapi.Object, filter *topoapi.Filter) bool {
	if object.Type != topoapi.Object_ENTITY && object.Type != topoapi.Object_RELATION {
		return false
	}

	eqo := filter.GetEqual_()
	if eqo != nil {
		if object.Type == topoapi.Object_ENTITY {
			return string(object.GetEntity().KindID) == eqo.Value
		}
		return string(object.GetRelation().KindID) == eqo.Value
	}
	igo := filter.GetIn()
	if igo != nil {
		var actual string
		if object.Type == topoapi.Object_ENTITY {
			actual = string(object.GetEntity().KindID)
		} else {
			actual = string(object.GetRelation().KindID)
		}
		for _, v := range igo.Values {
			if v == actual {
				return true
			}
		}
		return false
	}
	ngo := filter.GetNot()
	if ngo != nil {
		return !matchKind(object, ngo.Inner)
	}
	return false
}
