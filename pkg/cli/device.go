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
	"strings"
	"text/tabwriter"
	"time"

	"github.com/onosproject/onos-topo/api/topo"
	"github.com/onosproject/onos-topo/pkg/bulk"

	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/onosproject/onos-topo/api/device"
	"github.com/spf13/cobra"
)

func getGetDeviceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "device <id>",
		Aliases: []string{"devices"},
		Args:    cobra.MaximumNArgs(1),
		Short:   "Get a device",
		RunE:    runGetDeviceCommand,
	}
	cmd.Flags().BoolP("verbose", "v", false, "whether to print the device with verbose output")
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	return cmd
}

func runGetDeviceCommand(cmd *cobra.Command, args []string) error {
	verbose, _ := cmd.Flags().GetBool("verbose")
	noHeaders, _ := cmd.Flags().GetBool("no-headers")

	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	outputWriter := cli.GetOutput()
	writer := new(tabwriter.Writer)
	writer.Init(outputWriter, 0, 0, 3, ' ', tabwriter.FilterHTML)

	client := device.CreateDeviceServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if len(args) == 0 {
		stream, err := client.List(ctx, &device.ListRequest{})
		if err != nil {
			cli.Output("list error")
			return err
		}

		if !noHeaders {
			if verbose {
				_, _ = fmt.Fprintln(writer, "ID\tDISPLAYNAME\tADDRESS\tVERSION\tTYPE\tSTATE\tUSER\tPASSWORD\tATTRIBUTES")
			} else {
				_, _ = fmt.Fprintln(writer, "ID\tDISPLAYNAME\tADDRESS\tVERSION\tTYPE\tSTATE")
			}
		}

		for {
			response, err := stream.Recv()
			if err == io.EOF {
				break
			} else if err != nil {
				cli.Output("recv error")
				return err
			}

			dev := response.Device
			state := stateString(dev)
			if verbose {
				attributesBuf := bytes.Buffer{}
				for key, attribute := range dev.Attributes {
					attributesBuf.WriteString(key)
					attributesBuf.WriteString(": ")
					attributesBuf.WriteString(attribute)
					attributesBuf.WriteString(", ")
				}
				_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n", dev.ID, dev.Displayname, dev.Address, dev.Version, dev.Type,
					state, dev.Credentials.User, dev.Credentials.Password, attributesBuf.String())
			} else {
				_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\t%s\n", dev.ID, dev.Displayname, dev.Address, dev.Version, dev.Type, state)
			}
		}
	} else {
		response, err := client.Get(ctx, &device.GetRequest{
			ID: device.ID(args[0]),
		})
		if err != nil {
			cli.Output("get error")
			return err
		}

		dev := response.Device

		state := stateString(dev)
		_, _ = fmt.Fprintf(writer, "ID\t%s\n", dev.ID)
		_, _ = fmt.Fprintf(writer, "DisplayName\t%s\n", dev.Displayname)
		_, _ = fmt.Fprintf(writer, "ADDRESS\t%s\n", dev.Address)
		_, _ = fmt.Fprintf(writer, "VERSION\t%s\n", dev.Version)
		_, _ = fmt.Fprintf(writer, "TYPE\t%s\n", dev.Type)
		_, _ = fmt.Fprintf(writer, "STATE\t%s\n", state)
		if verbose {
			_, _ = fmt.Fprintf(writer, "USER\t%s\n", dev.Credentials.User)
			_, _ = fmt.Fprintf(writer, "PASSWORD\t%s\n", dev.Credentials.Password)
			for key, attribute := range dev.Attributes {
				_, _ = fmt.Fprintf(writer, "%s\t%s\n", strings.ToUpper(key), attribute)
			}
		}
	}
	return writer.Flush()
}

func stateString(dev *device.Device) string {
	stateBuf := bytes.Buffer{}
	for index, protocol := range dev.Protocols {
		stateBuf.WriteString(protocol.Protocol.String())
		stateBuf.WriteString(": {Connectivity: ")
		stateBuf.WriteString(protocol.ConnectivityState.String())
		stateBuf.WriteString(", Channel: ")
		stateBuf.WriteString(protocol.ChannelState.String())
		stateBuf.WriteString(", Service: ")
		stateBuf.WriteString(protocol.ServiceState.String())
		stateBuf.WriteString("}")
		if index != len(dev.Protocols) && len(dev.Protocols) != 1 {
			stateBuf.WriteString("\n")
		}
	}
	return stateBuf.String()
}

func getAddDeviceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "device <id> [args]",
		Aliases: []string{"devices"},
		Args:    cobra.ExactArgs(1),
		Short:   "Add a device",
		RunE:    runAddDeviceCommand,
	}
	cmd.Flags().StringP("type", "t", "", "the type of the device")
	cmd.Flags().StringP("role", "r", "", "the device role")
	cmd.Flags().StringP("target", "g", "", "the device target name")
	cmd.Flags().StringP("address", "a", "", "the address of the device")
	cmd.Flags().StringP("user", "u", "", "the device username")
	cmd.Flags().StringP("password", "p", "", "the device password")
	cmd.Flags().StringP("version", "v", "", "the device software version")
	cmd.Flags().StringP("displayname", "d", "", "A user friendly display name")
	cmd.Flags().String("key", "", "the TLS key")
	cmd.Flags().String("cert", "", "the TLS certificate")
	cmd.Flags().String("ca-cert", "", "the TLS CA certificate")
	cmd.Flags().Bool("plain", false, "whether to connect over a plaintext connection")
	cmd.Flags().Bool("insecure", false, "whether to enable skip verification")
	cmd.Flags().Duration("timeout", 5*time.Second, "the device connection timeout")
	cmd.Flags().StringToString("attributes", map[string]string{}, "an arbitrary mapping of device attributes")

	_ = cmd.MarkFlagRequired("version")
	_ = cmd.MarkFlagRequired("type")
	return cmd
}

func runAddDeviceCommand(cmd *cobra.Command, args []string) error {
	id := args[0]
	deviceType, _ := cmd.Flags().GetString("type")
	deviceRole, _ := cmd.Flags().GetString("role")
	deviceTarget, _ := cmd.Flags().GetString("target")
	address, _ := cmd.Flags().GetString("address")
	user, _ := cmd.Flags().GetString("user")
	password, _ := cmd.Flags().GetString("password")
	version, _ := cmd.Flags().GetString("version")
	displayName, _ := cmd.Flags().GetString("displayname")
	key, _ := cmd.Flags().GetString("key")
	cert, _ := cmd.Flags().GetString("cert")
	caCert, _ := cmd.Flags().GetString("ca-cert")
	plain, _ := cmd.Flags().GetBool("plain")
	insecure, _ := cmd.Flags().GetBool("insecure")
	timeout, _ := cmd.Flags().GetDuration("timeout")
	attributes, _ := cmd.Flags().GetStringToString("attributes")

	// Target defaults to the ID
	if deviceTarget == "" {
		deviceTarget = id
	}

	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := device.CreateDeviceServiceClient(conn)

	dev := &device.Device{
		ID:          device.ID(id),
		Type:        device.Type(deviceType),
		Role:        device.Role(deviceRole),
		Address:     address,
		Target:      deviceTarget,
		Version:     version,
		Displayname: displayName,
		Timeout:     &timeout,
		Credentials: device.Credentials{
			User:     user,
			Password: password,
		},
		TLS: device.TlsConfig{
			Cert:     cert,
			Key:      key,
			CaCert:   caCert,
			Plain:    plain,
			Insecure: insecure,
		},
		Attributes: attributes,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err = client.Add(ctx, &device.AddRequest{
		Device: dev,
	})
	if err != nil {
		return err
	}
	cli.Output("Added device %s \n", id)
	return nil
}

func getUpdateDeviceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "device <id> [args]",
		Aliases: []string{"devices"},
		Args:    cobra.ExactArgs(1),
		Short:   "Update a device",
		RunE:    runUpdateDeviceCommand,
	}
	cmd.Flags().StringP("type", "t", "", "the type of the device")
	cmd.Flags().StringP("role", "r", "", "the device role")
	cmd.Flags().StringP("target", "g", "", "the device target name")
	cmd.Flags().StringP("address", "a", "", "the address of the device")
	cmd.Flags().StringP("user", "u", "", "the device username")
	cmd.Flags().StringP("password", "p", "", "the device password")
	cmd.Flags().StringP("version", "v", "", "the device software version")
	cmd.Flags().StringP("displayname", "d", "", "A user friendly display name")
	cmd.Flags().String("key", "", "the TLS key")
	cmd.Flags().String("cert", "", "the TLS certificate")
	cmd.Flags().String("ca-cert", "", "the TLS CA certificate")
	cmd.Flags().Bool("plain", false, "whether to connect over a plaintext connection")
	cmd.Flags().Bool("insecure", false, "whether to enable skip verification")
	cmd.Flags().Duration("timeout", 30*time.Second, "the device connection timeout")
	cmd.Flags().StringToString("attributes", map[string]string{}, "an arbitrary mapping of device attributes")
	return cmd
}

func runUpdateDeviceCommand(cmd *cobra.Command, args []string) error {
	id := args[0]

	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return nil
	}
	defer conn.Close()

	client := device.CreateDeviceServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	response, err := client.Get(ctx, &device.GetRequest{
		ID: device.ID(id),
	})
	if err != nil {
		return err
	}

	cancel()
	dvc := response.Device

	if cmd.Flags().Changed("type") {
		deviceType, _ := cmd.Flags().GetString("type")
		dvc.Type = device.Type(deviceType)
	}
	if cmd.Flags().Changed("target") {
		deviceTarget, _ := cmd.Flags().GetString("target")
		dvc.Target = deviceTarget
	}
	if cmd.Flags().Changed("role") {
		deviceRole, _ := cmd.Flags().GetString("role")
		dvc.Role = device.Role(deviceRole)
	}
	if cmd.Flags().Changed("address") {
		address, _ := cmd.Flags().GetString("address")
		dvc.Address = address
	}
	if cmd.Flags().Changed("user") {
		user, _ := cmd.Flags().GetString("user")
		dvc.Credentials.User = user
	}
	if cmd.Flags().Changed("password") {
		password, _ := cmd.Flags().GetString("password")
		dvc.Credentials.Password = password
	}
	if cmd.Flags().Changed("version") {
		version, _ := cmd.Flags().GetString("version")
		dvc.Version = version
	}
	if cmd.Flags().Changed("displayname") {
		displayName, _ := cmd.Flags().GetString("displayname")
		dvc.Displayname = displayName
	}
	if cmd.Flags().Changed("key") {
		key, _ := cmd.Flags().GetString("key")
		dvc.TLS.Key = key
	}
	if cmd.Flags().Changed("cert") {
		cert, _ := cmd.Flags().GetString("cert")
		dvc.TLS.Cert = cert
	}
	if cmd.Flags().Changed("ca-cert") {
		caCert, _ := cmd.Flags().GetString("ca-cert")
		dvc.TLS.CaCert = caCert
	}
	if cmd.Flags().Changed("plain") {
		plain, _ := cmd.Flags().GetBool("plain")
		dvc.TLS.Plain = plain
	}
	if cmd.Flags().Changed("insecure") {
		insecure, _ := cmd.Flags().GetBool("insecure")
		dvc.TLS.Insecure = insecure
	}
	if cmd.Flags().Changed("timeout") {
		timeout, _ := cmd.Flags().GetDuration("timeout")
		dvc.Timeout = &timeout
	}
	if cmd.Flags().Changed("attributes") {
		attributes, _ := cmd.Flags().GetStringToString("attributes")
		dvc.Attributes = attributes
	}

	ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err = client.Update(ctx, &device.UpdateRequest{
		Device: dvc,
	})
	if err != nil {
		return err
	}
	cli.Output("Updated device %s", id)
	return nil
}

func getRemoveDeviceCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "device <id> [args]",
		Aliases: []string{"devices"},
		Args:    cobra.ExactArgs(1),
		Short:   "Remove a device",
		RunE:    runRemoveDeviceCommand,
	}
}

func runRemoveDeviceCommand(cmd *cobra.Command, args []string) error {
	id := args[0]

	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := device.CreateDeviceServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err = client.Remove(ctx, &device.RemoveRequest{
		Device: &device.Device{
			ID: device.ID(id),
		},
	})
	if err != nil {
		return err
	}
	cli.Output("Removed device %s", id)
	return nil
}

func getWatchDeviceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "device <id> [args]",
		Aliases: []string{"devices"},
		Args:    cobra.MaximumNArgs(1),
		Short:   "Watch for device changes",
		RunE:    runWatchDeviceCommand,
	}
	cmd.Flags().BoolP("verbose", "v", false, "whether to print the device with verbose output")
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	return cmd
}

func runWatchDeviceCommand(cmd *cobra.Command, args []string) error {
	var id string
	if len(args) > 0 {
		id = args[0]
	}

	verbose, _ := cmd.Flags().GetBool("verbose")
	noHeaders, _ := cmd.Flags().GetBool("no-headers")

	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := device.CreateDeviceServiceClient(conn)

	stream, err := client.List(context.Background(), &device.ListRequest{
		Subscribe: true,
	})
	if err != nil {
		return err
	}

	writer := new(tabwriter.Writer)
	writer.Init(cli.GetOutput(), 0, 0, 3, ' ', tabwriter.FilterHTML)

	if !noHeaders {
		if verbose {
			_, _ = fmt.Fprintln(writer, "EVENT\tID\tADDRESS\tVERSION\tUSER\tPASSWORD")
		} else {
			_, _ = fmt.Fprintln(writer, "EVENT\tID\tADDRESS\tVERSION")
		}
		_ = writer.Flush()
	}

	for {
		response, err := stream.Recv()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}

		dev := response.Device
		if id != "" && dev.ID != device.ID(id) {
			continue
		}

		if verbose {
			_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\t%s\n", response.Type, dev.ID, dev.Address, dev.Version, dev.Credentials.User, dev.Credentials.Password)
		} else {
			_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\t%s\n", response.Type, dev.ID, dev.Address, dev.Version)
		}
		_ = writer.Flush()
	}
}

// Deprecated: to be replaced by getLoadYamlEntitiesCommand
func getLoadYamlCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "yaml {file}",
		Args:  cobra.ExactArgs(1),
		Short: "Load topo data from a YAML file",
		RunE:  runLoadYamlCommand,
	}
	cmd.Flags().StringArray("attr", []string{""}, "Extra attributes to add to each device in k=v format")
	return cmd
}

func getLoadYamlEntitiesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "yamlentities {file}",
		Args:  cobra.ExactArgs(1),
		Short: "Load topo data from a YAML file",
		RunE:  runLoadYamlEntitiesCommand,
	}
	cmd.Flags().StringArray("attr", []string{""}, "Extra attributes to add to each device in k=v format")
	return cmd
}

// Deprecated: only used for getLoadYamlCommand() above
func runLoadYamlCommand(cmd *cobra.Command, args []string) error {
	var filename string
	if len(args) > 0 {
		filename = args[0]
	}

	extraAttrs, err := cmd.Flags().GetStringArray("attr")
	if err != nil {
		return err
	}
	for _, x := range extraAttrs {
		split := strings.Split(x, "=")
		if len(split) != 2 {
			return fmt.Errorf("expect extra args to be in the format a=b. Rejected: %s", x)
		}
	}

	deviceConfig, err := bulk.GetDeviceConfig(filename)
	if err != nil {
		return err
	}

	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	client := device.CreateDeviceServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	for _, dev := range deviceConfig.TopoDevices {
		if dev.Attributes == nil {
			dev.Attributes = make(map[string]string)
		}
		for _, x := range extraAttrs {
			split := strings.Split(x, "=")
			dev.Attributes[split[0]] = split[1]
		}

		dev := dev // pin
		resp, err := client.Add(ctx, &device.AddRequest{
			Device: &dev,
		})
		if err != nil {
			return err
		}
		if resp.Device.ID != dev.ID {
			return fmt.Errorf("error loading %s in to topo", dev.ID)
		}
	}

	fmt.Printf("Loaded %d topo devices from %s\n", len(deviceConfig.TopoDevices), filename)

	return nil
}

func runLoadYamlEntitiesCommand(cmd *cobra.Command, args []string) error {
	var filename string
	if len(args) > 0 {
		filename = args[0]
	}

	extraAttrs, err := cmd.Flags().GetStringArray("attr")
	if err != nil {
		return err
	}
	for _, x := range extraAttrs {
		cli.Output("runLoadYamlEntitiesCommand %v", x)
		split := strings.Split(x, "=")
		if len(split) != 2 {
			return fmt.Errorf("expect extra args to be in the format a=b. Rejected: %s", x)
		}
	}

	topoConfig, err := bulk.GetTopoConfig(filename)
	if err != nil {
		return err
	}

	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	client := topo.CreateTopoClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	request := topo.SetRequest{
		Objects: make([]*topo.Object, 0),
	}

	for _, kind := range topoConfig.TopoKinds {
		if kind.Attributes == nil {
			a := make(map[string]string)
			kind.Attributes = &a
		}
		for _, x := range extraAttrs {
			split := strings.Split(x, "=")
			(*kind.Attributes)[split[0]] = split[1]
		}

		kind := kind // pin
		request.Objects = append(request.Objects, bulk.TopoKindToTopoObject(&kind))
	}

	for _, entity := range topoConfig.TopoEntities {
		if entity.Attributes == nil {
			a := make(map[string]string)
			entity.Attributes = &a
		}
		for _, x := range extraAttrs {
			split := strings.Split(x, "=")
			(*entity.Attributes)[split[0]] = split[1]
		}

		entity := entity // pin
		request.Objects = append(request.Objects, bulk.TopoEntityToTopoObject(&entity))
	}

	for _, relation := range topoConfig.TopoRelations {
		if relation.Attributes == nil {
			a := make(map[string]string)
			relation.Attributes = &a
		}
		for _, x := range extraAttrs {
			split := strings.Split(x, "=")
			(*relation.Attributes)[split[0]] = split[1]
		}

		relation := relation // pin
		request.Objects = append(request.Objects, bulk.TopoRelationToTopoObject(&relation))
	}

	for _, relation := range topoConfig.TopoRelations {
		if relation.Attributes == nil {
			a := make(map[string]string)
			relation.Attributes = &a
		}
		for _, x := range extraAttrs {
			split := strings.Split(x, "=")
			(*relation.Attributes)[split[0]] = split[1]
		}

		relation := relation // pin
		request.Objects = append(request.Objects, bulk.TopoRelationToTopoObject(&relation))
	}
	_, err = client.Set(ctx, &request)
	if err != nil {
		return err
	}

	fmt.Printf("Loaded %d topo devices from %s\n", len(topoConfig.TopoEntities), filename)

	return nil
}
