// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package generator

import (
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-topo/pkg/tools/topo-generator/parser"
	"html/template"
	"os"
	"strings"
)

var log = logging.GetLogger()

// WriteFile will create the output file we desire (e-k-r file)
func WriteFile(result parser.NetworkLayer, filename string) error {
	// writing the entity-kind-relationship file
	pwd, _ := os.Getwd()
	s := strings.Split(pwd, "/onos-topo/")
	networkPath := s[0] + "/onos-topo/pkg/tools/topo-generator/generator/templates/network.yaml"
	switchPath := s[0] + "/onos-topo/pkg/tools/topo-generator/generator/templates/switch.yaml"
	portPath := s[0] + "/onos-topo/pkg/tools/topo-generator/generator/templates/port.yaml"
	linkPath := s[0] + "/onos-topo/pkg/tools/topo-generator/generator/templates/link.yaml"
	originatesPath := s[0] + "/onos-topo/pkg/tools/topo-generator/generator/templates/originates.yaml"
	terminatesPath := s[0] + "/onos-topo/pkg/tools/topo-generator/generator/templates/terminates.yaml"

	file, _ := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	for _, n := range result.Networks {
		// network
		t, err := template.ParseFiles(networkPath)
		if err != nil {
			log.Error(err.Error())
			return err
		}
		err = t.Execute(file, n)
		if err != nil {
			log.Error(err.Error())
			return err
		}

		// switches
		for _, s := range n.Switches {
			t, err := template.ParseFiles(switchPath)
			if err != nil {
				log.Error(err.Error())
				return err
			}
			err = t.Execute(file, s)
			if err != nil {
				log.Error(err.Error())
				return err
			}

			// ports
			for _, p := range s.Ports {
				t, err := template.ParseFiles(portPath)
				if err != nil {
					log.Error(err.Error())
					return err
				}
				err = t.Execute(file, p)
				if err != nil {
					log.Error(err.Error())
					return err
				}
			}
		}

		// links
		for _, l := range n.Links {
			// handles whether the link is unidirectional or bidirectional
			t, err := template.ParseFiles(linkPath)
			if err != nil {
				log.Error(err.Error())
				return err
			}
			err = t.Execute(file, l)
			if err != nil {
				log.Error(err.Error())
				return err
			}

			t1, err1 := template.ParseFiles(originatesPath)
			if err1 != nil {
				log.Error(err1.Error())
				return err1
			}
			err1 = t1.Execute(file, l.OriginatesRelation)
			if err1 != nil {
				log.Error(err1.Error())
				return err1
			}

			t2, err2 := template.ParseFiles(terminatesPath)
			if err2 != nil {
				log.Error(err2.Error())
				return err2
			}
			err2 = t2.Execute(file, l.TerminatesRelation)
			if err2 != nil {
				log.Error(err2.Error())
				return err2
			}
		}
	}
	defer file.Close()
	return nil
}
