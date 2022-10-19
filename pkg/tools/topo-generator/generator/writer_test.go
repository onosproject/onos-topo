// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0
package generator

import (
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-topo/pkg/tools/topo-generator/parser"
	"github.com/onosproject/onos-topo/pkg/tools/topo-generator/reader"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// TestWriter will verify that the writer is working properly
func TestWriter(t *testing.T) {
	log = logging.GetLogger()
	assert := assert.New(t)
	// TEST ONE: A simple topology containing 2 switches with 1 port each and a link between them
	err := WriteFile(parser.Convert(reader.ReadFile("../test_input_files/hr-1.yaml")), "ekr-1.yaml")
	assert.Equal(err, nil, "correct output")
	e := os.Remove("ekr-1.yaml")
	if e != nil {
		log.Error(e)
	}
	// TEST TWO: A simple topology with 2 networks, the first network having 1 switch with 0 ports
	err2 := WriteFile(parser.Convert(reader.ReadFile("../test_input_files/hr-2.yaml")), "ekr-2.yaml")
	assert.Equal(err2, nil, "correct output")
	e2 := os.Remove("ekr-2.yaml")
	if e2 != nil {
		log.Error(e2)
	}
	// TEST THREE: 2x2 topology
	err3 := WriteFile(parser.Convert(reader.ReadFile("../test_input_files/hr-3.yaml")), "ekr-3.yaml")
	assert.Equal(err3, nil, "correct output")
	e3 := os.Remove("ekr-3.yaml")
	if e3 != nil {
		log.Error(e3)
	}
	// TEST FOUR: The same topology as test one, but the link is unidirectional
	err4 := WriteFile(parser.Convert(reader.ReadFile("../test_input_files/hr-4.yaml")), "ekr-4.yaml")
	assert.Equal(err4, nil, "correct output")
	e4 := os.Remove("ekr-4.yaml")
	if e4 != nil {
		log.Error(e4)
	}
}
