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
	"github.com/spf13/cobra"
	viperapi "github.com/spf13/viper"
)

var viper = viperapi.New()

// init initializes the command line
func init() {
	initConfig()
}

// Init is a hook called after cobra initialization
func Init() {
	// noop for now
}

// GetCommand returns the root command for the topo service
func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "topo {get,add,update,remove,watch} [args]",
		Short: "ONOS topology subsystem commands",
	}

	addConfigFlags(cmd)

	cmd.AddCommand(getConfigCommand())
	cmd.AddCommand(getGetCommand())
	cmd.AddCommand(getAddCommand())
	cmd.AddCommand(getUpdateCommand())
	cmd.AddCommand(getRemoveCommand())
	cmd.AddCommand(getWatchCommand())
	return cmd
}
