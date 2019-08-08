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
	"github.com/golang/protobuf/ptypes"
	"github.com/onosproject/onos-topo/pkg/northbound/device"
	"github.com/spf13/cobra"
	"io"
	"os"
	"text/tabwriter"
	"time"
)

func getGetDeviceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "device <id>",
		Aliases: []string{"devices"},
		Args:    cobra.MaximumNArgs(1),
		Short:   "Get a device",
		Run:     runGetDeviceCommand,
	}
	cmd.Flags().BoolP("verbose", "v", false, "whether to print the device with verbose output")
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	return cmd
}

func runGetDeviceCommand(cmd *cobra.Command, args []string) {
	verbose, _ := cmd.Flags().GetBool("verbose")
	noHeaders, _ := cmd.Flags().GetBool("no-headers")

	conn := getConnection()
	defer conn.Close()

	client := device.NewDeviceServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if len(args) == 0 {
		stream, err := client.List(ctx, &device.ListRequest{})
		if err != nil {
			ExitWithError(ExitBadConnection, err)
		}

		writer := new(tabwriter.Writer)
		writer.Init(os.Stdout, 0, 0, 3, ' ', tabwriter.FilterHTML)

		if !noHeaders {
			if verbose {
				fmt.Fprintln(writer, "ID\tADDRESS\tVERSION\tUSER\tPASSWORD")
			} else {
				fmt.Fprintln(writer, "ID\tADDRESS\tVERSION")
			}
		}

		for {
			response, err := stream.Recv()
			if err == io.EOF {
				break
			} else if err != nil {
				ExitWithError(ExitError, err)
			}

			dvc := response.Device
			if verbose {
				fmt.Fprintln(writer, fmt.Sprintf("%s\t%s\t%s\t%s\t%s", dvc.Id, dvc.Address, dvc.SoftwareVersion, dvc.Credentials.User, dvc.Credentials.Password))
			} else {
				fmt.Fprintln(writer, fmt.Sprintf("%s\t%s\t%s", dvc.Id, dvc.Address, dvc.SoftwareVersion))
			}
		}
		writer.Flush()
	} else {
		response, err := client.Get(ctx, &device.GetRequest{
			DeviceId: args[0],
		})
		if err != nil {
			ExitWithError(ExitBadConnection, err)
		}

		dvc := response.Device

		writer := new(tabwriter.Writer)
		writer.Init(os.Stdout, 0, 0, 3, ' ', tabwriter.FilterHTML)
		fmt.Fprintln(writer, fmt.Sprintf("ID\t%s", dvc.Id))
		fmt.Fprintln(writer, fmt.Sprintf("ADDRESS\t%s", dvc.Address))
		fmt.Fprintln(writer, fmt.Sprintf("VERSION\t%s", dvc.SoftwareVersion))

		if verbose {
			fmt.Fprintln(writer, fmt.Sprintf("USER\t%s", dvc.Credentials.User))
			fmt.Fprintln(writer, fmt.Sprintf("PASSWORD\t%s", dvc.Credentials.Password))
		}
		writer.Flush()
	}
}

func getAddDeviceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "device <id> [args]",
		Aliases: []string{"devices"},
		Args:    cobra.ExactArgs(1),
		Short:   "Add a device",
		Run:     runAddDeviceCommand,
	}
	cmd.Flags().StringP("address", "a", "", "the address of the device")
	cmd.Flags().StringP("user", "u", "", "the device username")
	cmd.Flags().StringP("password", "p", "", "the device password")
	cmd.Flags().StringP("version", "v", "", "the device software version")
	cmd.Flags().String("key", "", "the TLS key")
	cmd.Flags().String("cert", "", "the TLS certificate")
	cmd.Flags().String("ca-cert", "", "the TLS CA certificate")
	cmd.Flags().DurationP("timeout", "t", 30*time.Second, "the device connection timeout")
	return cmd
}

func runAddDeviceCommand(cmd *cobra.Command, args []string) {
	id := args[0]
	address, _ := cmd.Flags().GetString("address")
	user, _ := cmd.Flags().GetString("user")
	password, _ := cmd.Flags().GetString("password")
	version, _ := cmd.Flags().GetString("version")
	key, _ := cmd.Flags().GetString("key")
	cert, _ := cmd.Flags().GetString("cert")
	caCert, _ := cmd.Flags().GetString("ca-cert")
	timeout, _ := cmd.Flags().GetDuration("timeout")

	conn := getConnection()
	defer conn.Close()

	client := device.NewDeviceServiceClient(conn)

	dvc := &device.Device{
		Id:              id,
		Address:         address,
		SoftwareVersion: version,
		Timeout:         ptypes.DurationProto(timeout),
		Credentials: &device.Credentials{
			User:     user,
			Password: password,
		},
		Tls: &device.TlsConfig{
			Cert:   cert,
			Key:    key,
			CaCert: caCert,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := client.Add(ctx, &device.AddRequest{
		Device: dvc,
	})
	if err != nil {
		ExitWithError(ExitBadConnection, err)
	} else {
		ExitWithOutput("Added device %s", id)
	}
}

func getUpdateDeviceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "device <id> [args]",
		Aliases: []string{"devices"},
		Args:    cobra.ExactArgs(1),
		Short:   "Update a device",
		Run:     runUpdateDeviceCommand,
	}
	cmd.Flags().StringP("address", "a", "", "the address of the device")
	cmd.Flags().StringP("user", "u", "", "the device username")
	cmd.Flags().StringP("password", "p", "", "the device password")
	cmd.Flags().StringP("version", "v", "", "the device software version")
	cmd.Flags().String("key", "", "the TLS key")
	cmd.Flags().String("cert", "", "the TLS certificate")
	cmd.Flags().String("ca-cert", "", "the TLS CA certificate")
	cmd.Flags().DurationP("timeout", "t", 30*time.Second, "the device connection timeout")
	return cmd
}

func runUpdateDeviceCommand(cmd *cobra.Command, args []string) {
	id := args[0]

	conn := getConnection()
	defer conn.Close()

	client := device.NewDeviceServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	response, err := client.Get(ctx, &device.GetRequest{
		DeviceId: id,
	})
	if err != nil {
		ExitWithError(ExitBadConnection, err)
	}

	cancel()
	dvc := response.Device

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
		dvc.SoftwareVersion = version
	}
	if cmd.Flags().Changed("key") {
		key, _ := cmd.Flags().GetString("key")
		dvc.Tls.Key = key
	}
	if cmd.Flags().Changed("cert") {
		cert, _ := cmd.Flags().GetString("cert")
		dvc.Tls.Cert = cert
	}
	if cmd.Flags().Changed("ca-cert") {
		caCert, _ := cmd.Flags().GetString("ca-cert")
		dvc.Tls.CaCert = caCert
	}
	if cmd.Flags().Changed("timeout") {
		timeout, _ := cmd.Flags().GetDuration("timeout")
		dvc.Timeout = ptypes.DurationProto(timeout)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err = client.Update(ctx, &device.UpdateRequest{
		Device: dvc,
	})
	if err != nil {
		ExitWithError(ExitBadConnection, err)
	} else {
		ExitWithOutput("Updated device %s", id)
	}
}

func getRemoveDeviceCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "device <id> [args]",
		Aliases: []string{"devices"},
		Args:    cobra.ExactArgs(1),
		Short:   "Remove a device",
		Run:     runRemoveDeviceCommand,
	}
}

func runRemoveDeviceCommand(cmd *cobra.Command, args []string) {
	id := args[0]

	conn := getConnection()
	defer conn.Close()

	client := device.NewDeviceServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := client.Remove(ctx, &device.RemoveRequest{
		Device: &device.Device{
			Id: id,
		},
	})
	if err != nil {
		ExitWithError(ExitBadConnection, err)
	} else {
		ExitWithOutput("Removed device %s", id)
	}
}

func getWatchDeviceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "device <id> [args]",
		Aliases: []string{"devices"},
		Args:    cobra.MaximumNArgs(1),
		Short:   "Watch for device changes",
		Run:     runWatchDeviceCommand,
	}
	cmd.Flags().BoolP("verbose", "v", false, "whether to print the device with verbose output")
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	return cmd
}

func runWatchDeviceCommand(cmd *cobra.Command, args []string) {
	var id string
	if len(args) > 0 {
		id = args[0]
	}

	verbose, _ := cmd.Flags().GetBool("verbose")
	noHeaders, _ := cmd.Flags().GetBool("no-headers")

	conn := getConnection()
	defer conn.Close()

	client := device.NewDeviceServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	stream, err := client.List(ctx, &device.ListRequest{
		Subscribe: true,
	})
	if err != nil {
		ExitWithError(ExitBadConnection, err)
	}

	writer := new(tabwriter.Writer)
	writer.Init(os.Stdout, 0, 0, 3, ' ', tabwriter.FilterHTML)

	if !noHeaders {
		if verbose {
			fmt.Fprintln(writer, "EVENT\tID\tADDRESS\tVERSION\tUSER\tPASSWORD")
		} else {
			fmt.Fprintln(writer, "EVENT\tID\tADDRESS\tVERSION")
		}
		writer.Flush()
	}

	for {
		response, err := stream.Recv()
		if err == io.EOF {
			ExitWithSuccess()
		} else if err != nil {
			ExitWithError(ExitError, err)
		}

		device := response.Device
		if id != "" && device.Id != id {
			continue
		}

		if verbose {
			fmt.Fprintln(writer, fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s", response.Type, device.Id, device.Address, device.SoftwareVersion, device.Credentials.User, device.Credentials.Password))
		} else {
			fmt.Fprintln(writer, fmt.Sprintf("%s\t%s\t%s\t%s", response.Type, device.Id, device.Address, device.SoftwareVersion))
		}
		writer.Flush()
	}
}
