// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package store

import topoapi "github.com/onosproject/onos-api/go/onos/topo"

func match(object *topoapi.Object, filters *topoapi.Filters) bool {
	return filters == nil ||
		(matchKind(object, filters.KindFilter) && matchLabels(object, filters.LabelFilters) && matchAspects(object, filters.WithAspects))
}

func matchLabels(object *topoapi.Object, filters []*topoapi.Filter) bool {
	for _, filter := range filters {
		if !matchLabel(object, filter) {
			return false
		}
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
	if filter == nil {
		return true
	}
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

// Returns true if object type is any of the given types; false otherwise
func matchType(object *topoapi.Object, types []topoapi.Object_Type) bool {
	if len(types) != 0 {
		for i := range types {
			if object.Type == types[i] {
				return true
			}
		}
		return false
	}
	return true
}

// Returns true if object has all requested aspects; false otherwise
func matchAspects(object *topoapi.Object, aspects []string) bool {
	for i := range aspects {
		if _, ok := object.Aspects[aspects[i]]; !ok {
			return false
		}
	}
	return true
}
