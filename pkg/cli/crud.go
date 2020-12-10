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

import "github.com/spf13/cobra"

func getGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get {entity|relation|kind} [args]",
		Short: "Get topology resources",
	}
	cmd.AddCommand(getGetEntityCommand())
	cmd.AddCommand(getGetRelationCommand())
	cmd.AddCommand(getGetKindCommand())
	return cmd
}

func getAddCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add {entity|relation|kind} [args]",
		Short: "Add a topology resource",
	}
	cmd.AddCommand(getAddEntityCommand())
	cmd.AddCommand(getAddRelationCommand())
	cmd.AddCommand(getAddKindCommand())
	return cmd
}

// FIXME: implement update for entity|relation|kind
func getUpdateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update {...} [args]",
		Short: "Update a topology resource",
	}
	return cmd
}

// FIXME: implement remove for entity|relation|kind
// FIXME: fix inconsistency... what is "object"? Is it "entity"?!
func getRemoveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove {object} [args]",
		Short: "Remove a topology resource",
	}
	cmd.AddCommand(getRemoveObjectCommand())
	return cmd
}

func getWatchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "watch {entity|relation|kind|all} [args]",
		Short: "Watch for changes to a topology resource type",
	}
	cmd.AddCommand(getWatchEntityCommand())
	cmd.AddCommand(getWatchRelationCommand())
	cmd.AddCommand(getWatchKindCommand())
	cmd.AddCommand(getWatchAllCommand())
	return cmd
}

func getLoadCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "load {topofile}",
		Short: "Bulk load topo data from a file",
	}

	cmd.AddCommand(getLoadYamlEntitiesCommand())

	return cmd
}
