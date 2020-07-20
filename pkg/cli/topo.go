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
	"bytes"
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
		Short:   "Get Entity",
		RunE:    runGetEntityCommand,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	return cmd
}

func getGetRelationCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "relation <id>",
		Aliases: []string{"relations"},
		Args:    cobra.MaximumNArgs(1),
		Short:   "Get Relation",
		RunE:    runGetRelationCommand,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	return cmd
}

func getGetKindCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "kind <id>",
		Aliases: []string{"kinds"},
		Args:    cobra.MaximumNArgs(1),
		Short:   "Get Kind",
		RunE:    runGetKindCommand,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	return cmd
}

func getAddEntityCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "entity <id> [args]",
		Args:  cobra.MinimumNArgs(1),
		Short: "Add Entity",
		RunE:  runAddEntityCommand,
	}
	cmd.Flags().StringP("kind", "k", "", "Kind ID")
	//_ = cmd.MarkFlagRequired("kind")
	cmd.Flags().StringToStringP("attributes", "a", map[string]string{}, "an user defined mapping of entity attributes")
	return cmd
}

func getAddRelationCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "relation <id> <src-entity-id> <tgt-entity-id> [args]",
		Args:  cobra.MinimumNArgs(3),
		Short: "Add Relation",
		RunE:  runAddRelationCommand,
	}
	cmd.Flags().StringP("kind", "k", "", "Kind ID")
	//_ = cmd.MarkFlagRequired("kind")
	return cmd
}

func getAddKindCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kind <id> <name> [args]",
		Args:  cobra.MinimumNArgs(2),
		Short: "Add Kind",
		RunE:  runAddKindCommand,
	}
	return cmd
}

func getRemoveObjectCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "object <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Remove an object",
		RunE:  runRemoveObjectCommand,
	}
}

func getWatchEntityCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "entity [id] [args]",
		Short: "Watch Entities",
		Args:  cobra.MaximumNArgs(2),
		RunE:  runWatchEntityCommand,
	}
	cmd.Flags().BoolP("noreplay", "r", false, "do not replay past topo updates")
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	return cmd
}

func getWatchRelationCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "relation <id> [args]",
		Short: "Watch Relations",
		Args:  cobra.MaximumNArgs(2),
		RunE:  runWatchRelationCommand,
	}
	cmd.Flags().BoolP("noreplay", "r", false, "do not replay past topo updates")
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	return cmd
}

func getWatchKindCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kind [id] [args]",
		Short: "Watch Kinds",
		Args:  cobra.MaximumNArgs(2),
		RunE:  runWatchKindCommand,
	}
	cmd.Flags().BoolP("noreplay", "r", false, "do not replay past topo updates")
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	return cmd
}

func getWatchAllCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all [args]",
		Short: "Watch Entities and Relations",
		RunE:  runWatchAllCommand,
	}
	cmd.Flags().BoolP("noreplay", "r", false, "do not replay past topo updates")
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	return cmd
}

func runGetEntityCommand(cmd *cobra.Command, args []string) error {
	return runGetCommand(cmd, args, topo.Object_ENTITY)
}

func runGetRelationCommand(cmd *cobra.Command, args []string) error {
	return runGetCommand(cmd, args, topo.Object_RELATION)
}

func runGetKindCommand(cmd *cobra.Command, args []string) error {
	return runGetCommand(cmd, args, topo.Object_KIND)
}

func runAddEntityCommand(cmd *cobra.Command, args []string) error {
	return writeObject(cmd, args, topo.Object_ENTITY)
}

func runAddRelationCommand(cmd *cobra.Command, args []string) error {
	return writeObject(cmd, args, topo.Object_RELATION)
}

func runAddKindCommand(cmd *cobra.Command, args []string) error {
	return writeObject(cmd, args, topo.Object_KIND)
}

func runRemoveObjectCommand(cmd *cobra.Command, args []string) error {
	id := args[0]

	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := topo.CreateTopoClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err = client.Delete(ctx, &topo.DeleteRequest{ID: topo.ID(id)})
	if err != nil {
		return err
	}
	cli.Output("Removed object %s", id)
	return nil
}

func runWatchEntityCommand(cmd *cobra.Command, args []string) error {
	return watch(cmd, args, topo.Object_ENTITY)
}

func runWatchRelationCommand(cmd *cobra.Command, args []string) error {
	return watch(cmd, args, topo.Object_RELATION)
}

func runWatchKindCommand(cmd *cobra.Command, args []string) error {
	return watch(cmd, args, topo.Object_KIND)
}

func runWatchAllCommand(cmd *cobra.Command, args []string) error {
	return watch(cmd, args, topo.Object_UNSPECIFIED)
}

func runGetCommand(cmd *cobra.Command, args []string, objectType topo.Object_Type) error {
	noHeaders, _ := cmd.Flags().GetBool("no-headers")

	if !noHeaders {
		printHeader(false)
	}

	if len(args) == 0 {
		for object := range listObjects(cmd, args, objectType) {
			if object != nil {
				if objectType == topo.Object_UNSPECIFIED || objectType == object.Type {
					printRow(object, false, noHeaders)
				}
			}
		}
	} else {
		id := args[0]
		object, err := getObject(cmd, topo.ID(id))
		if err != nil {
			return err
		}
		if object != nil {
			if objectType == topo.Object_UNSPECIFIED || objectType == object.Type {
				printRow(object, false, noHeaders)
			}
		}
	}

	return nil
}

func writeObject(cmd *cobra.Command, args []string, objectType topo.Object_Type) error {
	var object *topo.Object
	id := args[0]
	attributes, _ := cmd.Flags().GetStringToString("attributes")

	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := topo.CreateTopoClient(conn)

	if objectType == topo.Object_ENTITY {
		kindID, _ := cmd.Flags().GetString("kind")
		entity := &topo.Object_Entity{
			Entity: &topo.Entity{
				KindID: topo.ID(kindID),
			},
		}

		object = &topo.Object{
			ID:         topo.ID(id),
			Type:       objectType,
			Obj:        entity,
			Attributes: attributes,
		}
	} else if objectType == topo.Object_RELATION {
		kindID, _ := cmd.Flags().GetString("kind")
		relation := &topo.Object_Relation{
			Relation: &topo.Relation{
				KindID:      topo.ID(kindID),
				SrcEntityID: topo.ID(args[1]),
				TgtEntityID: topo.ID(args[2]),
			},
		}

		object = &topo.Object{
			ID:   topo.ID(id),
			Type: objectType,
			Obj:  relation,
		}
	} else if objectType == topo.Object_KIND {
		kind := &topo.Object_Kind{
			Kind: &topo.Kind{
				Name: args[1],
			},
		}

		object = &topo.Object{
			ID:   topo.ID(id),
			Type: objectType,
			Obj:  kind,
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err = client.Set(ctx, &topo.SetRequest{Objects: []*topo.Object{object}})
	if err != nil {
		return err
	}
	return nil
}

func listObjects(cmd *cobra.Command, args []string, objectType topo.Object_Type) <-chan *topo.Object {
	out := make(chan *topo.Object)

	go func() {
		defer close(out)
		conn, err := cli.GetConnection(cmd)
		if err != nil {
			return
		}
		defer conn.Close()

		client := topo.CreateTopoClient(conn)

		stream, err := client.List(context.Background(), &topo.ListRequest{})
		if err != nil {
			return
		}
		for {
			response, err := stream.Recv()
			if err == io.EOF {
				// read done.
				break
			} else if err != nil {
				cli.Output("recv error : %v", err)
				return
			}
			out <- response.Object
		}
	}()

	return out
}

func getObject(cmd *cobra.Command, id topo.ID) (*topo.Object, error) {
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := topo.CreateTopoClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	response, err := client.Get(ctx, &topo.GetRequest{ID: id})
	if err != nil {
		cli.Output("get error")
		return nil, err
	}
	return response.Object, nil
}

func watch(cmd *cobra.Command, args []string, objectType topo.Object_Type) error {
	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	noreplay, _ := cmd.Flags().GetBool("noreplay")

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
		ID:       id,
		Noreplay: noreplay,
	})
	if err != nil {
		return err
	}

	if !noHeaders {
		printHeader(true)
	}
	for update := range watchStream(stream) {
		if objectType == topo.Object_UNSPECIFIED || objectType == update.Object.Type {
			printUpdateType(update.Type)
			printRow(update.Object, true, noHeaders)
		}
	}

	return nil
}

func watchStream(stream topo.Topo_SubscribeClient) <-chan *topo.Update {
	out := make(chan *topo.Update)

	go func() {
		defer close(out)
		for {
			response, err := stream.Recv()
			if err == io.EOF {
				// read done.
				return
			}
			if err != nil {
				cli.Output("Error receiving notification : %v", err)
				return
			}

			out <- response.Update
		}
	}()
	return out
}

func printHeader(printUpdateType bool) {
	var width = 16
	var prec = width - 1
	writer := new(tabwriter.Writer)
	writer.Init(cli.GetOutput(), 0, 0, 3, ' ', tabwriter.FilterHTML)

	if printUpdateType {
		_, _ = fmt.Fprintf(writer, "%-*.*s", width, prec, "Update Type")
	}
	_, _ = fmt.Fprintf(writer, "%-*.*s%-*.*s%-*.*s%-*.*s\n", width, prec, "Object Type", width, prec, "Object ID", width, prec, "Kind ID", width, prec, "Attributes")
}

func printUpdateType(updateType topo.Update_Type) {
	var width = 16
	var prec = width - 1
	writer := new(tabwriter.Writer)
	writer.Init(cli.GetOutput(), 0, 0, 3, ' ', tabwriter.FilterHTML)
	if updateType == topo.Update_UNSPECIFIED {
		_, _ = fmt.Fprintf(writer, "%-*.*s", width, prec, "REPLAY")
	} else {
		_, _ = fmt.Fprintf(writer, "%-*.*s", width, prec, updateType)
	}
	_ = writer.Flush()
}

func printRow(object *topo.Object, watch bool, noHeaders bool) {
	var width = 16
	var prec = width - 1
	writer := new(tabwriter.Writer)
	writer.Init(cli.GetOutput(), 0, 0, 3, ' ', tabwriter.FilterHTML)

	switch object.Type {
	case topo.Object_ENTITY:
		e := object.GetEntity()
		// printUpdateType()
		_, _ = fmt.Fprintf(writer, "%-*.*s%-*.*s%-*.*s%s\n", width, prec, object.Type, width, prec, object.ID, width, prec, e.KindID, attrsToString(object.Attributes))
	case topo.Object_RELATION:
		r := object.GetRelation()
		// printUpdateType()
		_, _ = fmt.Fprintf(writer, "%-*.*s%-*.*s%-*.*s", width, prec, object.Type, width, prec, object.ID, width, prec, r.KindID)
		_, _ = fmt.Fprintf(writer, "src=%s, tgt=%s, %s\n", r.SrcEntityID, r.TgtEntityID, attrsToString(object.Attributes))
	case topo.Object_KIND:
		k := object.GetKind()
		// printUpdateType()
		_, _ = fmt.Fprintf(writer, "%-*.*s%-*.*s%-*.*s\n", width, prec, object.Type, width, prec, object.ID, width, prec, k.GetName())
	default:
		_, _ = fmt.Fprintf(writer, "\n")
	}
	_ = writer.Flush()
}

func attrsToString(attrs map[string]string) string {
	attributesBuf := bytes.Buffer{}
	first := true
	for key, attribute := range attrs {
		if !first {
			attributesBuf.WriteString(", ")
		} else {
			first = false
		}
		attributesBuf.WriteString(key)
		attributesBuf.WriteString(":")
		attributesBuf.WriteString(attribute)
	}
	return attributesBuf.String()
}
