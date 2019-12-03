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
	"os"
)

const (
	addressKey     = "address"
	defaultAddress = "onos-topo:5150"

	noTLSKey       = "no-tls"
	tlsCertPathKey = "tls.certPath"
	tlsKeyPathKey  = "tls.keyPath"
)

func init() {
	cobra.OnInitialize(initConfig)
}

var (
	configOptions = []string{
		addressKey,
		noTLSKey,
		tlsCertPathKey,
		tlsKeyPathKey,
	}

	configFile = ""
)

func addConfigFlags(cmd *cobra.Command) {
	viper.SetDefault(addressKey, defaultAddress)
	cmd.PersistentFlags().StringP("address", "a", viper.GetString(addressKey), "the onos-topo service address")
	cmd.PersistentFlags().String("tls-cert-path", viper.GetString(tlsCertPathKey), "the path to the TLS certificate")
	cmd.PersistentFlags().String("tls-key-path", viper.GetString(tlsKeyPathKey), "the path to the TLS key")
	cmd.PersistentFlags().Bool("no-tls", viper.GetBool("no-tls"), "if present, do not use TLS")
	cmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default: $HOME/.onos/topo.yaml)")

	_ = viper.BindPFlag(addressKey, cmd.PersistentFlags().Lookup("address"))
	_ = viper.BindPFlag(tlsCertPathKey, cmd.PersistentFlags().Lookup("tls-cert-path"))
	_ = viper.BindPFlag(tlsKeyPathKey, cmd.PersistentFlags().Lookup("tls-key-path"))
	_ = viper.BindPFlag(noTLSKey, cmd.PersistentFlags().Lookup("no-tls"))
}

func getConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config {init,set,get,delete} [args]",
		Short: "Manage the CLI configuration",
	}
	cmd.AddCommand(newConfigInitCommand())
	cmd.AddCommand(newConfigGetCommand())
	cmd.AddCommand(newConfigSetCommand())
	cmd.AddCommand(newConfigDeleteCommand())
	return cmd
}

func newConfigInitCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize the onos topo CLI configuration",
		Run: func(cmd *cobra.Command, args []string) {
			if err := viper.ReadInConfig(); err == nil {
				ExitWithSuccess()
			}
			home, err := homedir.Dir()
			if err != nil {
				ExitWithError(ExitError, err)
			}
			err = os.MkdirAll(home+"/.onos", 0777)
			if err != nil {
				ExitWithError(ExitError, err)
			}
			f, err := os.Create(home + "/.onos/topo.yaml")
			if err != nil {
				ExitWithError(ExitError, err)
			} else {
				f.Close()
			}
			err = viper.WriteConfig()
			if err != nil {
				ExitWithError(ExitError, err)
			} else {
				ExitWithOutput("Created ~/.onos/topo.yaml\n")
			}
		},
	}
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
	address, _ := cmd.Flags().GetString("address")
	return address
}

func getCertPath(cmd *cobra.Command) string {
	certPath, _ := cmd.Flags().GetString("tls-cert-path")
	return certPath
}

func getKeyPath(cmd *cobra.Command) string {
	keyPath, _ := cmd.Flags().GetString("tls-key-path")
	return keyPath
}

func noTLS(cmd *cobra.Command) bool {
	tls, _ := cmd.Flags().GetBool("no-tls")
	return tls
}

func initConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			ExitWithError(ExitError, err)
		}

		viper.SetConfigName("topo")
		viper.AddConfigPath(home + "/.onos")

		_ = viper.ReadInConfig()
	}
}
