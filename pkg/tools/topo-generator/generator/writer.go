// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package generator

import (
	"html/template"
	"os"
	"github.com/onosproject/onos-topo/pkg/tools/topo-generator/parser"
	"github.com/onosproject/onos-lib-go/pkg/logging"
)

var log = logging.GetLogger()

func WriteFile(underlay parser.Underlay, filename string) {
	// writing the entity-kind-relationship file
	file, _ := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	for _, n := range underlay.Networks {
		// network
		t, err := template.ParseFiles("github.com/onosproject/onos-topo/pkg/tools/topo-generator/generator/templates/network.yaml")
		if err != nil {
			log.Error(err.Error())
			return
		}
		err = t.Execute(file, n)
		
		// switches
		for _, s := range n.Switches {
			t, err := template.ParseFiles("github.com/onosproject/onos-topo/pkg/tools/topo-generator/generator/templates/switch.yaml")
			if err != nil {
				log.Error(err.Error())
				return
			}
			err = t.Execute(file, s)
			
			// ports
			for _, p := range s.Ports {
				t, err := template.ParseFiles("github.com/onosproject/onos-topo/pkg/tools/topo-generator/generator/templates/port.yaml")
				if err != nil {
					log.Error(err.Error())
					return
				}
				err = t.Execute(file, p)
			}
		}

		// links
		for _, l := range n.Links {
			// handles whether the link is unidirectional or bidirectional
			if l.Link_Type == "unidirectional" {
				t, err := template.ParseFiles("github.com/onosproject/onos-topo/pkg/tools/topo-generator/generator/templates/unidirectional.yaml")
				if err != nil {
					log.Error(err.Error())
					return
				}
				err = t.Execute(file, l)
			} else {
				t, err := template.ParseFiles("github.com/onosproject/onos-topo/pkg/tools/topo-generator/generator/templates/bidirectional.yaml")
				if err != nil {
					log.Error(err.Error())
					return
				}
				err = t.Execute(file, l)
			}
		}
	}
	defer file.Close()
}
