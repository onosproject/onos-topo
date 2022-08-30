// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package visualizer

import (
	"encoding/json"
	"github.com/onosproject/onos-api/go/onos/topo"
)

// TopoEvent is used to serialize topology events to JSON
type TopoEvent struct {
	Event    string            `json:"event"`
	Type     string            `json:"type"`
	UUID     topo.UUID         `json:"uuid"`
	ID       topo.ID           `json:"id"`
	Revision topo.Revision     `json:"revision"`
	Entity   *TopoEntity       `json:"entity"`
	Relation *TopoRelation     `json:"relation"`
	Kind     *TopoKind         `json:"kind"`
	Aspects  map[string]string `json:"aspects"`
	// TODO: add aspects and labels
}

// TopoEntity represents an entity
type TopoEntity struct {
	Kind topo.ID `json:"kind"`
}

// TopoRelation represents a relation
type TopoRelation struct {
	Kind topo.ID `json:"kind"`
	Src  topo.ID `json:"src"`
	Tgt  topo.ID `json:"tgt"`
}

// TopoKind represents a kind
type TopoKind struct {
	Name string `json:"name"`
}

// EncodeTopoEvent transforms the watch response event into JSON
func EncodeTopoEvent(msg *topo.WatchResponse) ([]byte, error) {
	o := msg.Event.Object
	te := &TopoEvent{
		Event:    eventType(msg.Event.Type),
		UUID:     o.UUID,
		ID:       o.ID,
		Revision: o.Revision,
		Aspects:  make(map[string]string),
	}

	switch o.Type {
	case topo.Object_ENTITY:
		ent := o.GetEntity()
		te.Type = "entity"
		te.Entity = &TopoEntity{
			Kind: ent.KindID,
		}
	case topo.Object_RELATION:
		rel := o.GetRelation()
		te.Type = "relation"
		te.Relation = &TopoRelation{
			Kind: rel.KindID,
			Src:  rel.SrcEntityID,
			Tgt:  rel.TgtEntityID,
		}
	case topo.Object_KIND:
		kind := o.GetKind()
		te.Type = "kind"
		te.Kind = &TopoKind{
			Name: kind.Name,
		}
	}

	for at, av := range o.Aspects {
		te.Aspects[at] = string(av.Value)
	}

	return json.Marshal(te)
}

func eventType(eventType topo.EventType) string {
	switch eventType {
	case topo.EventType_ADDED:
		return "added"
	case topo.EventType_UPDATED:
		return "updated"
	case topo.EventType_REMOVED:
		return "removed"
	}
	return "replay"
}
