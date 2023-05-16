// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/spf13/cobra"

	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-topo/pkg/manager"
)

var log = logging.GetLogger()

// The main entry point
func main() {
	cmd := &cobra.Command{
		Use:  "onos-topo",
		RunE: runRootCommand,
	}
	cli.AddServiceEndpointFlags(cmd, "onos-topo gRPC")
	cli.Run(cmd)
}

func runRootCommand(cmd *cobra.Command, _ []string) error {
	flags, err := cli.ExtractServiceEndpointFlags(cmd)
	if err != nil {
		return err
	}

	log.Infof("Starting onos-topo")
	return cli.RunDaemon(manager.NewManager(manager.Config{ServiceFlags: flags}))
}
