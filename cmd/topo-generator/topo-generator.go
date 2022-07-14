// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"flag"
	"github.com/onosproject/onos-topo/pkg/tools/topo-generator/generator"
	"github.com/onosproject/onos-topo/pkg/tools/topo-generator/parser"
	"github.com/onosproject/onos-topo/pkg/tools/topo-generator/reader"
)

// main entry point
func main() {
	input := flag.String("input_file_path", "hr-1.yaml", "input_file")
	output := flag.String("output_file_path", "ekr-1.yaml", "output_file")
	flag.Parse()
	generator.WriteFile(parser.Convert(reader.ReadFile(*input)), *output)
}
