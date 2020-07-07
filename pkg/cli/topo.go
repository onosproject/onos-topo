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
	"context"
	"fmt"
	"io"
	"text/tabwriter"
	"time"

	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/onosproject/onos-topo/api/topo"
	"github.com/spf13/cobra"
)

func getGetEntityCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "entity <id>",
		Aliases: []string{"entities"},
		Args:    cobra.MaximumNArgs(1),
		Short:   "Get a topo entity",
		RunE:    runGetEntityCommand,
	}
	return cmd
}

func runGetEntityCommand(cmd *cobra.Command, args []string) error {
	return runGetCommand(cmd, args, topo.Object_ENTITY)
}

func getAddEntityCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "entity <id> [args]",
		Args:  cobra.MinimumNArgs(1),
		Short: "Add an entity",
		RunE:  runAddEntityCommand,
	}
	cmd.Flags().StringP("type", "t", "", "the type of the entity")
	//_ = cmd.MarkFlagRequired("type")
	return cmd
}

func runAddEntityCommand(cmd *cobra.Command, args []string) error {
	return writeObject(cmd, args, topo.Object_ENTITY, topo.Update_INSERT)
}

func getGetRelationCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "relation <id>",
		Aliases: []string{"relations"},
		Args:    cobra.MaximumNArgs(1),
		Short:   "Get topo relationships",
		RunE:    runGetRelationCommand,
	}

	return cmd
}

func runGetRelationCommand(cmd *cobra.Command, args []string) error {
	return runGetCommand(cmd, args, topo.Object_RELATIONSHIP)
}

func runGetCommand(cmd *cobra.Command, args []string, objectType topo.Object_Type) error {
	var objects []*topo.Object

	updates := make(chan *topo.Update)
	done := make(chan bool)
	defer close(updates)
	defer close(done)

	go printIt(updates, objectType, done)

	objects, err := readObjects(cmd, args, objectType)
	if err != nil {
		return err
	}

	for _, obj := range objects {
		updates <- &topo.Update{
			Type:   topo.Update_UNSPECIFIED,
			Object: obj,
		}
	}

	updates <- &topo.Update{}
	<-done

	return nil
}

func getAddRelationCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "relation <id> <src-entity-id> <tgt-entity-id> [args]",
		Args:  cobra.MinimumNArgs(3),
		Short: "Add a topo relationship",
		RunE:  runAddRelationCommand,
	}
	cmd.Flags().StringP("type", "t", "", "the type of the entity")
	return cmd
}

func runAddRelationCommand(cmd *cobra.Command, args []string) error {
	return writeObject(cmd, args, topo.Object_RELATIONSHIP, topo.Update_INSERT)
}

func writeObject(cmd *cobra.Command, args []string, objectType topo.Object_Type, updateType topo.Update_Type) error {
	id := args[0]

	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := topo.CreateTopoClient(conn)

	updates := make([]*topo.Update, 1)

	if objectType == topo.Object_ENTITY {
		entityType, _ := cmd.Flags().GetString("type")
		object := &topo.Object_Entity{
			Entity: &topo.Entity{
				Type: entityType,
			},
		}

		updates[0] = &topo.Update{
			Type: updateType,
			Object: &topo.Object{
				Ref: &topo.Reference{
					ID: topo.ID(id)},
				Type: objectType,
				Obj:  object,
			},
		}
	} else if objectType == topo.Object_RELATIONSHIP {
		object := &topo.Object_Relationship{
			Relationship: &topo.Relationship{
				SourceRef: &topo.Reference{ID: topo.ID(args[1])},
				TargetRef: &topo.Reference{ID: topo.ID(args[2])},
			},
		}

		updates[0] = &topo.Update{
			Type: updateType,
			Object: &topo.Object{
				Ref: &topo.Reference{
					ID: topo.ID(id)},
				Type: objectType,
				Obj:  object,
			},
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err = client.Write(ctx, &topo.WriteRequest{Updates: updates})
	if err != nil {
		return err
	}
	return nil
}

func readObjects(cmd *cobra.Command, args []string, objectType topo.Object_Type) ([]*topo.Object, error) {

	var objects []*topo.Object

	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := topo.CreateTopoClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if len(args) == 0 {
		stream, err := client.Subscribe(context.Background(), &topo.SubscribeRequest{
			Ref:      &topo.Reference{},
			SnapShot: true,
		})
		if err != nil {
			return nil, err
		}
		updates := make(chan *topo.Update)
		go watchStream(stream, updates)
		printIt(updates, objectType, nil)
	} else {
		id := args[0]
		reference := &topo.Reference{
			ID: topo.ID(id),
		}
		refs := []*topo.Reference{reference}
		response, err := client.Read(ctx, &topo.ReadRequest{Refs: refs})
		if err != nil {
			cli.Output("get error")
			return nil, err
		}
		objects = response.Objects
	}
	return objects, nil
}

func getWatchEntityCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "entity [id] [args]",
		Short: "Watch for entity changes",
		Args:  cobra.MaximumNArgs(2),
		RunE:  runWatchEntityCommand,
	}
	cmd.Flags().BoolP("replay", "r", false, "whether to replay past topo updates")
	return cmd
}

func getWatchRelationCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "relation <id> [args]",
		Short: "Watch for relationship changes",
		Args:  cobra.MaximumNArgs(2),
		RunE:  runWatchRelationCommand,
	}
	cmd.Flags().BoolP("replay", "r", false, "whether to replay past topo updates")
	return cmd
}

func getWatchAllCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all [args]",
		Short: "Watch for entity and relationship changes",
		RunE:  runWatchAllCommand,
	}
	cmd.Flags().BoolP("replay", "r", false, "whether to replay past topo updates")
	return cmd
}

func runWatchEntityCommand(cmd *cobra.Command, args []string) error {
	return watch(cmd, args, topo.Object_ENTITY)
}

func runWatchRelationCommand(cmd *cobra.Command, args []string) error {
	return watch(cmd, args, topo.Object_RELATIONSHIP)
}

func runWatchAllCommand(cmd *cobra.Command, args []string) error {
	return watch(cmd, args, topo.Object_UNSPECIFIED)
}

func watch(cmd *cobra.Command, args []string, objectType topo.Object_Type) error {
	replay, _ := cmd.Flags().GetBool("replay")
	var id topo.ID
	if len(args) > 0 {
		id = topo.ID(args[0])
	} else {
		id = topo.ID(topo.NullID)
	}

	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := topo.CreateTopoClient(conn)

	stream, err := client.Subscribe(context.Background(), &topo.SubscribeRequest{
		Ref: &topo.Reference{
			ID: id,
		},
		WithoutReplay: !replay,
	})
	if err != nil {
		return err
	}

	updates := make(chan *topo.Update)

	go watchStream(stream, updates)

	printIt(updates, objectType, nil)

	return nil
}
func watchStream(stream topo.Topo_SubscribeClient, updates chan *topo.Update) {
	for {
		response, err := stream.Recv()
		if err == io.EOF {
			// read done.
			close(updates)
			return
		}
		if err != nil {
			cli.Output("Error receiving notification : %v", err)
			close(updates)
			return
		}

		for _, update := range response.Updates {
			updates <- update
		}
	}
}

func printIt(updates chan *topo.Update, objectType topo.Object_Type, done chan bool) {
	writer := new(tabwriter.Writer)
	writer.Init(cli.GetOutput(), 0, 0, 3, ' ', tabwriter.FilterHTML)

	for update := range updates {
		if update.Object == nil {
			break
		}
		switch update.Object.Type {
		case topo.Object_ENTITY:
			e := update.Object.GetEntity()
			if objectType == topo.Object_UNSPECIFIED || objectType == topo.Object_ENTITY {
				_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\t%s\n", update.Type, update.Object.Type, update.Object.Ref.ID, e.Type)
			}
		case topo.Object_RELATIONSHIP:
			r := update.Object.GetRelationship()
			if objectType == topo.Object_UNSPECIFIED || objectType == topo.Object_RELATIONSHIP {
				_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t", update.Type, update.Object.Type, update.Object.Ref.ID, r.Type)
				_, _ = fmt.Fprintf(writer, "%s\t%s\n", r.SourceRef.ID, r.TargetRef.ID)
			}
		default:
			_, _ = fmt.Fprintf(writer, "\n")
		}
		_ = writer.Flush()
	}
	if done != nil {
		done <- true
	}
}
