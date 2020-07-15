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
	var object *topo.Object
	noHeaders, _ := cmd.Flags().GetBool("no-headers")

	updates := make(chan *topo.Update)
	done := make(chan bool)
	defer close(updates)
	defer close(done)

	go printIt(updates, objectType, done, false, noHeaders)

	object, err := readObjects(cmd, args, objectType)
	if err != nil {
		return err
	}

	updates <- &topo.Update{
		Type:   topo.Update_UNSPECIFIED,
		Object: object,
	}

	updates <- nil
	<-done

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

func readObjects(cmd *cobra.Command, args []string, objectType topo.Object_Type) (*topo.Object, error) {
	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	var object *topo.Object

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
			ID:       topo.NullID,
			Snapshot: true,
		})
		if err != nil {
			return nil, err
		}
		updates := make(chan *topo.Update)
		go watchStream(stream, updates)
		printIt(updates, objectType, nil, false, noHeaders)
	} else {
		id := args[0]
		response, err := client.Get(ctx, &topo.GetRequest{ID: topo.ID(id)})
		if err != nil {
			cli.Output("get error")
			return nil, err
		}
		object = response.Object
	}
	return object, nil
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

	updates := make(chan *topo.Update)

	go watchStream(stream, updates)

	printIt(updates, objectType, nil, true, noHeaders)

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

		updates <- response.Update
	}
}

func printIt(updates chan *topo.Update, objectType topo.Object_Type, done chan bool, watch bool, noHeaders bool) {
	var width = 16
	var prec = width - 1
	writer := new(tabwriter.Writer)
	writer.Init(cli.GetOutput(), 0, 0, 3, ' ', tabwriter.FilterHTML)

	if !noHeaders {
		if watch {
			_, _ = fmt.Fprintf(writer, "%-*.*s", width, prec, "Update Type")
		}
		_, _ = fmt.Fprintf(writer, "%-*.*s%-*.*s%-*.*s%-*.*s\n", width, prec, "Object Type", width, prec, "Reference ID", width, prec, "Object Kind", width, prec, "Attributes")
	}

	for update := range updates {
		u := update
		printUpdateType := func() {
			if watch {
				if u.Type == topo.Update_UNSPECIFIED {
					_, _ = fmt.Fprintf(writer, "%-*.*s", width, prec, "REPLAY")
				} else {
					_, _ = fmt.Fprintf(writer, "%-*.*s", width, prec, u.Type)
				}
			}
		}
		if u == nil {
			break
		}
		switch u.Object.Type {
		case topo.Object_ENTITY:
			e := u.Object.GetEntity()
			printUpdateType()
			if objectType == topo.Object_UNSPECIFIED || objectType == topo.Object_ENTITY {
				_, _ = fmt.Fprintf(writer, "%-*.*s%-*.*s%-*.*s%s\n", width, prec, u.Object.Type, width, prec, u.Object.ID, width, prec, e.KindID, attrsToString(u.Object.Attributes))
			}
		case topo.Object_RELATION:
			r := u.Object.GetRelation()
			printUpdateType()
			if objectType == topo.Object_UNSPECIFIED || objectType == topo.Object_RELATION {
				_, _ = fmt.Fprintf(writer, "%-*.*s%-*.*s%-*.*s", width, prec, u.Object.Type, width, prec, u.Object.ID, width, prec, r.KindID)
				_, _ = fmt.Fprintf(writer, "src=%s, tgt=%s, %s\n", r.SrcEntityID, r.TgtEntityID, attrsToString(u.Object.Attributes))
			}
		case topo.Object_KIND:
			k := u.Object.GetKind()
			printUpdateType()
			if objectType == topo.Object_UNSPECIFIED || objectType == topo.Object_KIND {
				_, _ = fmt.Fprintf(writer, "%-*.*s%-*.*s%-*.*s\n", width, prec, u.Object.Type, width, prec, u.Object.ID, width, prec, k.GetName())
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
