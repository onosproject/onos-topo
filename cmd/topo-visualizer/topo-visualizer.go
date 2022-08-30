// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/onosproject/onos-lib-go/pkg/cli"
	visualizer "github.com/onosproject/onos-topo/pkg/tools/topo-visualizer"
	"github.com/spf13/cobra"
	"os"
)

const (
	serviceAddress = "onos-topo:5150"
)

// The main entry point
func main() {
	if err := getRootCommand().Execute(); err != nil {
		println(err)
		os.Exit(1)
	}
}

func getRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "topo-visualizer",
		Short: "Starts HTTP/WS server for visualizing topology entities and relations",
		RunE:  runServer,
	}
	cli.AddEndpointFlags(cmd, serviceAddress)
	return cmd
}

func runServer(cmd *cobra.Command, args []string) error {
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	return visualizer.NewServer(conn).Serve()
}
