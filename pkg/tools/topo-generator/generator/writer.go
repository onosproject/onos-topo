// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package generator

import (
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-topo/pkg/tools/topo-generator/parser"
	"html/template"
	"os"
)

var log = logging.GetLogger()

// WriteFile will create the output file we desire (e-k-r file)
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
		if err != nil {
			log.Error(err.Error())
		}

		// switches
		for _, s := range n.Switches {
			t, err := template.ParseFiles("github.com/onosproject/onos-topo/pkg/tools/topo-generator/generator/templates/switch.yaml")
			if err != nil {
				log.Error(err.Error())
				return
			}
			err = t.Execute(file, s)
			if err != nil {
				log.Error(err.Error())
				return
			}

			// ports
			for _, p := range s.Ports {
				t, err := template.ParseFiles("github.com/onosproject/onos-topo/pkg/tools/topo-generator/generator/templates/port.yaml")
				if err != nil {
					log.Error(err.Error())
					return
				}
				err = t.Execute(file, p)
				if err != nil {
					log.Error(err.Error())
					return
				}
			}
		}

		// links
		for _, l := range n.Links {
			// handles whether the link is unidirectional or bidirectional
			if l.LinkType == "unidirectional" {
				t, err := template.ParseFiles("github.com/onosproject/onos-topo/pkg/tools/topo-generator/generator/templates/unidirectional.yaml")
				if err != nil {
					log.Error(err.Error())
					return
				}
				err = t.Execute(file, l)
				if err != nil {
					log.Error(err.Error())
					return
				}
			} else {
				t, err := template.ParseFiles("github.com/onosproject/onos-topo/pkg/tools/topo-generator/generator/templates/bidirectional.yaml")
				if err != nil {
					log.Error(err.Error())
					return
				}
				err = t.Execute(file, l)
				if err != nil {
					log.Error(err.Error())
					return
				}
			}
		}
	}
	defer file.Close()
}
