// Copyright 2019-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cli

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configFile = ""
)

const (
	addressKey     = "address"
	defaultAddress = "onos-topo:5150"

	tlsCertPathKey = "tls.certPath"
	tlsKeyPathKey  = "tls.keyPath"
)

func init() {
	cobra.OnInitialize(initConfig)
}

var configOptions = []string{
	addressKey,
	tlsCertPathKey,
	tlsKeyPathKey,
}

func addConfigFlags(cmd *cobra.Command) {
	viper.SetDefault(addressKey, defaultAddress)
	cmd.PersistentFlags().StringP("address", "a", viper.GetString(addressKey), "the onos-topo service address")
	cmd.PersistentFlags().String("tls-key-path", viper.GetString(tlsKeyPathKey), "the path to the TLS key")
	cmd.PersistentFlags().String("tls-cert-path", viper.GetString(tlsCertPathKey), "the path to the TLS certificate")
}

func getConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config {set,get,delete} [args]",
		Short: "Read and update CLI configuration options",
	}
	cmd.AddCommand(newConfigGetCommand())
	cmd.AddCommand(newConfigSetCommand())
	cmd.AddCommand(newConfigDeleteCommand())
	return cmd
}

func newConfigGetCommand() *cobra.Command {
	return &cobra.Command{
		Use:       "get <key>",
		Args:      cobra.ExactArgs(1),
		ValidArgs: configOptions,
		RunE:      runConfigGetCommand,
	}
}

func runConfigGetCommand(_ *cobra.Command, args []string) error {
	value := viper.Get(args[0])
	fmt.Fprintln(GetOutput(), value)
	return nil
}

func newConfigSetCommand() *cobra.Command {
	return &cobra.Command{
		Use:       "set <key> <value>",
		Args:      cobra.ExactArgs(2),
		ValidArgs: configOptions,
		RunE:      runConfigSetCommand,
	}
}

func runConfigSetCommand(_ *cobra.Command, args []string) error {
	viper.Set(args[0], args[1])
	if err := viper.WriteConfig(); err != nil {
		return err
	}

	value := viper.Get(args[0])
	fmt.Fprintln(GetOutput(), value)
	return nil
}

func newConfigDeleteCommand() *cobra.Command {
	return &cobra.Command{
		Use:       "delete <key>",
		Args:      cobra.ExactArgs(1),
		ValidArgs: configOptions,
		RunE:      runConfigDeleteCommand,
	}
}

func runConfigDeleteCommand(_ *cobra.Command, args []string) error {
	viper.Set(args[0], nil)
	if err := viper.WriteConfig(); err != nil {
		return err
	}

	value := viper.Get(args[0])
	fmt.Fprintln(GetOutput(), value)
	return nil
}

func getAddress(cmd *cobra.Command) string {
	address, _ := cmd.PersistentFlags().GetString("address")
	return address
}

func getCertPath(cmd *cobra.Command) string {
	certPath, _ := cmd.PersistentFlags().GetString("tls-cert-path")
	return certPath
}

func getKeyPath(cmd *cobra.Command) string {
	keyPath, _ := cmd.PersistentFlags().GetString("tls-key-path")
	return keyPath
}

func initConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			panic(err)
		}

		viper.SetConfigName("topo")
		viper.AddConfigPath(home + "/.onos")
		viper.AddConfigPath("/etc/onos")
		viper.AddConfigPath(".")
	}

	viper.ReadInConfig()
}
