// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/onosproject/onos-topo/pkg/tools/topo-generator/reader"
	"github.com/onosproject/onos-topo/pkg/tools/topo-generator/parser"
	"github.com/onosproject/onos-topo/pkg/tools/topo-generator/generator"
	"flag"
)

// main entry point
func main() {
	input_file_path := flag.String("input_file_path", "hr-1.yaml", "input_file")
	output_file_path := flag.String("output_file_path", "ekr-1.yaml", "output_file")
	flag.Parse()
	generator.WriteFile(parser.Convert(reader.ReadFile(*input_file_path)), *output_file_path)
}
