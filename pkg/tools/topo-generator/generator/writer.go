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

const (
	templatePath = "/onos-topo/pkg/tools/topo-generator/generator/templates/template.yaml"
)

var log = logging.GetLogger()

// WriteFile will create the output file we desire (e-k-r file)
func WriteFile(result parser.NetworkLayer, filename string) error {
	// writing the entity-kind-relationship file
	pwd, _ := os.Getwd()
	s := strings.Split(pwd, "/onos-topo/")
	path := s[0] + templatePath

	file, _ := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	t, err := template.ParseFiles(path)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	err = t.Execute(file, result)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}

	defer file.Close()
	return nil
}
