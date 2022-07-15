// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0
//
package generator

import (
	"github.com/onosproject/onos-topo/pkg/tools/topo-generator/parser"
	"github.com/onosproject/onos-topo/pkg/tools/topo-generator/reader"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestWriter will verify that the writer is working properly
func TestWriter(t *testing.T) {
	assert := assert.New(t)
	// TEST ONE: A simple topology containing 2 switches with 1 port each and a link between them
	err := WriteFile(parser.Convert(reader.ReadFile("../test_input_files/hr-1.yaml")), "test_output_files/ekr-1.yaml")
	assert.Equal(err, nil, "correct output")
	// TEST TWO: A simple topology with 2 networks, the first network having 1 switch with 0 ports
	err2 := WriteFile(parser.Convert(reader.ReadFile("../test_input_files/hr-2.yaml")), "test_output_files/ekr-2.yaml")
	assert.Equal(err2, nil, "correct output")
	// TEST THREE: 2x2 topology
	err3 := WriteFile(parser.Convert(reader.ReadFile("../test_input_files/hr-3.yaml")), "test_output_files/ekr-3.yaml")
	assert.Equal(err3, nil, "correct output")
	// TEST FOUR: The same topology as test one, but the link is unidirectional
	err4 := WriteFile(parser.Convert(reader.ReadFile("../test_input_files/hr-4.yaml")), "test_output_files/ekr-4.yaml")
	assert.Equal(err4, nil, "correct output")
}
