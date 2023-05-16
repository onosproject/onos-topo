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
	AddEndpointFlags(cmd, serviceAddress)
	return cmd
}

func runServer(cmd *cobra.Command, _ []string) error {
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	return visualizer.NewServer(conn).Serve()
}

// FIXME: Remove this after clearing up the onos-lib-go fiasco.

const (
	// ServiceAddress command option
	ServiceAddress = "service-address"
	// TLSCertPathFlag command option
	TLSCertPathFlag = "tls-cert-path"
	// TLSKeyPathFlag command option
	TLSKeyPathFlag = "tls-key-path"
	// NoTLSFlag command option
	NoTLSFlag = "no-tls"
)

// AddEndpointFlags adds service address, TLS cert path and TLS key path option to the command.
func AddEndpointFlags(cmd *cobra.Command, defaultAddress string) {
	cmd.Flags().String(ServiceAddress, defaultAddress, "service address; defaults to "+defaultAddress)
	cmd.Flags().String(TLSKeyPathFlag, "", "path to client private key")
	cmd.Flags().String(TLSCertPathFlag, "", "path to client certificate")
	cmd.Flags().Bool(NoTLSFlag, false, "if present, do not use TLS")
}
