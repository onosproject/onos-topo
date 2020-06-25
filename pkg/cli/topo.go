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
		Use:   "entity <id>",
		Args:  cobra.MinimumNArgs(1),
		Short: "Get a topo entity",
		RunE:  runGetEntityCommand,
	}
	/*
		cmd.Flags().StringP("id", "i", "", "the id of the entity")
		cmd.Flags().BoolP("verbose", "v", false, "whether to print the entity with verbose output")

		_ = cmd.MarkFlagRequired("id")
	*/

	return cmd
}

func runGetEntityCommand(cmd *cobra.Command, args []string) error {

	var objects []*topo.Object

	outputWriter := cli.GetOutput()
	writer := new(tabwriter.Writer)
	writer.Init(outputWriter, 0, 0, 3, ' ', tabwriter.FilterHTML)

	objects, err := readObjects(cmd, args)
	if err != nil {
		return err
	}

	if len(objects) != 0 {
		switch obj := objects[0].Obj.(type) {
		case *topo.Object_Entity:
			_, _ = fmt.Fprintf(writer, "ID\t%s\n", objects[0].Ref.GetID())
			_, _ = fmt.Fprintf(writer, "Type\t%s\n", obj.Entity.GetType())
		case nil:
			cli.Output("No object is set")
			// No object is set
		default:
			cli.Output("get error")
			// return ERROR
		}
	}
	return writer.Flush()
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
		Use:   "relation <id>",
		Args:  cobra.MinimumNArgs(1),
		Short: "Get a topo relationship",
		RunE:  runGetRelationCommand,
	}

	return cmd
}

func runGetRelationCommand(cmd *cobra.Command, args []string) error {

	var objects []*topo.Object

	outputWriter := cli.GetOutput()
	writer := new(tabwriter.Writer)
	writer.Init(outputWriter, 0, 0, 3, ' ', tabwriter.FilterHTML)

	objects, err := readObjects(cmd, args)
	if err != nil {
		return err
	}

	if len(objects) != 0 {
		switch obj := objects[0].Obj.(type) {
		case *topo.Object_Relationship:
			_, _ = fmt.Fprintf(writer, "ID\t%s\n", objects[0].Ref.GetID())
			_, _ = fmt.Fprintf(writer, "type\t%s\n", obj.Relationship.GetType())
			for _, ref := range obj.Relationship.GetSourceRefs() {
				_, _ = fmt.Fprintf(writer, "src-entity-id\t%s\n", string(ref.ID))
			}
			for _, ref := range obj.Relationship.GetTargetRefs() {
				_, _ = fmt.Fprintf(writer, "tgt-entity-id\t%s\n", string(ref.ID))
			}
		case nil:
			cli.Output("Error: nil object\n")
		default:
			cli.Output("Error: get error\n")
		}
	}
	return writer.Flush()
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
				SourceRefs: []*topo.Reference{{ID: topo.ID(args[1])}},
				TargetRefs: []*topo.Reference{{ID: topo.ID(args[2])}},
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

func readObjects(cmd *cobra.Command, args []string) ([]*topo.Object, error) {

	var objects []*topo.Object
	id := args[0]

	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := topo.CreateTopoClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if len(args) == 0 {
		// TODO - implement List function
	} else {
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
		Use:   "entity <id> [args]",
		Short: "Watch for topo changes",
		RunE:  runWatchEntityCommand,
	}
	return cmd
}

func runWatchEntityCommand(cmd *cobra.Command, args []string) error {
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

	stream, err := client.Subscribe(context.Background())
	if err != nil {
		return err
	}

	waitc := make(chan struct{})

	writer := new(tabwriter.Writer)
	writer.Init(cli.GetOutput(), 0, 0, 3, ' ', tabwriter.FilterHTML)

	go func() {
		for {
			response, err := stream.Recv()
			if err == io.EOF {
				// read done.
				close(waitc)
				return
			}
			if err != nil {
				cli.Output("Error receiving notification : %v", err)
				close(waitc)
				return
			}

			update := response.Update

			_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\t", update.Type, update.Object.Type, update.Object.Ref.ID)
			switch obj := update.Object.Obj.(type) {
			case *topo.Object_Entity:
				_, _ = fmt.Fprintf(writer, "%s\n", obj.Entity.Type)
			case *topo.Object_Relationship:
				_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\n", obj.Relationship.Type, obj.Relationship.SourceRefs[0], obj.Relationship.TargetRefs[0])
			default:
				_, _ = fmt.Fprintf(writer, "\n")
			}

			_ = writer.Flush()
		}
	}()

	subscribeRequest := &topo.SubscribeRequest{Refs: []*topo.Reference{{ID: id}}}
	if err := stream.Send(subscribeRequest); err != nil {
		close(waitc)
	}

	_ = stream.CloseSend()
	<-waitc
	return nil
}
