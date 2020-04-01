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

// Unit tests for device CLI
package cli

import (
	"bytes"
	"fmt"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/onosproject/onos-topo/api/device"
	"gotest.tools/assert"
	"strings"
	"testing"
	"time"
)

const version = "1.0.0"
const deviceType = "TestDevice"

func generateDeviceData(count int) []*device.Device {
	devices := make([]*device.Device, count)
	second := time.Second

	for devIdx := range devices {
		devices[devIdx] = &device.Device{
			ID:          device.ID(fmt.Sprintf("test-device-%d", devIdx)),
			Revision:    device.Revision(devIdx),
			Address:     fmt.Sprintf("192.168.0.%d", devIdx),
			Target:      "",
			Version:     version,
			Timeout:     &second,
			Credentials: device.Credentials{},
			TLS:         device.TlsConfig{},
			Type:        deviceType,
			Role:        "leaf",
			Attributes:  nil,
			Protocols:   nil,
		}
	}
	return devices
}

func Test_GetDevice(t *testing.T) {
	outputBuffer := bytes.NewBufferString("")
	cli.CaptureOutput(outputBuffer)

	setUpMockClients()
	getDevices := getGetDeviceCommand()
	args := make([]string, 1)
	args[0] = "test-device-1"
	getDevices.SetArgs(args)
	err := getDevices.Execute()
	assert.NilError(t, err)
	output := outputBuffer.String()
	assert.Assert(t, strings.Contains(output, "test-device"))
}

func Test_AddDevice(t *testing.T) {
	outputBuffer := bytes.NewBufferString("")
	cli.CaptureOutput(outputBuffer)

	setUpMockClients()
	addDevice := getAddDeviceCommand()
	args := make([]string, 7)
	args[0] = "test-device-1" // Name
	args[1] = fmt.Sprintf("--type=%s", deviceType)
	args[2] = fmt.Sprintf("--version=%s", version)
	args[3] = "--address=192.168.0.1"
	args[4] = "--timeout=1s"
	args[5] = "--user=test"
	args[6] = "--role=leaf"
	addDevice.SetArgs(args)
	err := addDevice.Execute()
	assert.NilError(t, err)
	output := outputBuffer.String()
	assert.Assert(t, strings.Contains(output, "Added device test-device-1"))
}

func Test_UpdateDevice(t *testing.T) {
	outputBuffer := bytes.NewBufferString("")
	cli.CaptureOutput(outputBuffer)

	setUpMockClients()
	updateDevice := getUpdateDeviceCommand()
	args := make([]string, 7)
	args[0] = "test-device-1" // Name
	args[1] = fmt.Sprintf("--type=%s", deviceType)
	args[2] = fmt.Sprintf("--version=%s", version)
	args[3] = "--address=192.168.0.1"
	args[4] = "--timeout=1s"
	args[5] = "--user=test"
	args[6] = "--role=leaf"
	updateDevice.SetArgs(args)
	err := updateDevice.Execute()
	assert.NilError(t, err)
	output := outputBuffer.String()
	assert.Assert(t, strings.Contains(output, "Updated device test-device-1"))
}

func Test_RemoveDevice(t *testing.T) {
	outputBuffer := bytes.NewBufferString("")
	cli.CaptureOutput(outputBuffer)

	setUpMockClients()
	removeDevice := getRemoveDeviceCommand()
	args := make([]string, 1)
	args[0] = "test-device-1" // Name
	removeDevice.SetArgs(args)
	err := removeDevice.Execute()
	assert.NilError(t, err)
	output := outputBuffer.String()
	assert.Assert(t, strings.Contains(output, "Removed device test-device-1"))
}
