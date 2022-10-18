// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0
package parser

import (
	"github.com/onosproject/onos-topo/pkg/tools/topo-generator/reader"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestParser will verify that the converter is working properly
func TestParser(t *testing.T) {
	assert := assert.New(t)

	// TEST ONE: A simple topology containing 2 switches with 1 port each and a link between them
	result := Convert(reader.ReadFile("../test_input_files/hr-1.yaml"))
	assert.Equal(len(result.Networks), 1, "The number of networks should be the same")

	// for the singular network, check switches and links
	assert.Equal(len(result.Networks[0].Links), 2, "The number of links should be the same")
	assert.Equal(len(result.Networks[0].Switches), 2, "The number of switches should be the same")
	// verify the entity_id and display_name for the network
	assert.Equal(result.Networks[0].EntityID, "network-layer:0/underlay-1", "The network entity id should be the same")
	assert.Equal(result.Networks[0].DisplayName, "underlay-1", "The network display name should be the same")

	// for the 2 switches, check ports
	assert.Equal(len(result.Networks[0].Switches[0].Ports), 1, "The number of networks should be the same")
	// verify the entity_id and p4rt_server_endpoint for the switch
	assert.Equal(result.Networks[0].Switches[0].EntityID, "p4rt:1", "The switch entity id should be the same")
	assert.Equal(result.Networks[0].Switches[0].P4RTPort, 50002, "The switch p4rt server endpoint should be the same")

	assert.Equal(len(result.Networks[0].Switches[1].Ports), 1, "The number of networks should be the same")
	// verify the entity_id and p4rt_server_endpoint for the switch
	assert.Equal(result.Networks[0].Switches[1].EntityID, "p4rt:2", "The switch entity id should be the same")
	assert.Equal(result.Networks[0].Switches[1].P4RTPort, 50001, "The switch p4rt server endpoint should be the same")

	// TEST TWO: A simple topology with 2 networks, the first network having 1 switch with 0 ports
	result2 := Convert(reader.ReadFile("../test_input_files/hr-2.yaml"))
	assert.Equal(len(result2.Networks), 2, "The number of networks should be the same")

	// for the first network, check switches
	assert.Equal(len(result2.Networks[0].Switches), 1, "The number of switches should be the same")
	// verify the entity_id and display_name for the network
	assert.Equal(result2.Networks[0].EntityID, "network-layer:0/underlay-1", "The network entity id should be the same")
	assert.Equal(result2.Networks[0].DisplayName, "underlay-1", "The network display name should be the same")

	// verify the entity_id and p4rt_server_endpoint for the switch
	assert.Equal(result2.Networks[0].Switches[0].EntityID, "p4rt:1", "The switch entity id should be the same")
	assert.Equal(result2.Networks[0].Switches[0].P4RTPort, 50002, "The switch p4rt server endpoint should be the same")

	// verify the entity_id and display_name for the second network
	assert.Equal(result2.Networks[1].EntityID, "network-layer:0/underlay-2", "The network entity id should be the same")
	assert.Equal(result2.Networks[1].DisplayName, "underlay-2", "The network display name should be the same")

	// TEST THREE: 2x2 topology
	result3 := Convert(reader.ReadFile("../test_input_files/hr-3.yaml"))
	assert.Equal(len(result3.Networks), 1, "The number of networks should be the same")

	// for the singular network, check switches and links
	assert.Equal(len(result3.Networks[0].Links), 8, "The number of links should be the same")
	assert.Equal(len(result3.Networks[0].Switches), 4, "The number of switches should be the same")
	// verify the entity_id and display_name for the network
	assert.Equal(result3.Networks[0].EntityID, "network-layer:0/underlay-1", "The network entity id should be the same")
	assert.Equal(result3.Networks[0].DisplayName, "underlay-1", "The network display name should be the same")

	// now we want to make sure that each switch has 2 ports and verify the entity_id and p4rt_server_endpoint for the switch
	assert.Equal(len(result3.Networks[0].Switches[0].Ports), 2, "The number of networks should be the same")
	// verify the entity_id and p4rt_server_endpoint for the switch
	assert.Equal(result3.Networks[0].Switches[0].EntityID, "p4rt:1", "The switch entity id should be the same")
	assert.Equal(result3.Networks[0].Switches[0].P4RTPort, 50001, "The switch p4rt server endpoint should be the same")

	assert.Equal(len(result3.Networks[0].Switches[1].Ports), 2, "The number of networks should be the same")
	// verify the entity_id and p4rt_server_endpoint for the switch
	assert.Equal(result3.Networks[0].Switches[1].EntityID, "p4rt:2", "The switch entity id should be the same")
	assert.Equal(result3.Networks[0].Switches[1].P4RTPort, 50002, "The switch p4rt server endpoint should be the same")

	assert.Equal(len(result3.Networks[0].Switches[2].Ports), 2, "The number of networks should be the same")
	// verify the entity_id and p4rt_server_endpoint for the switch
	assert.Equal(result3.Networks[0].Switches[2].EntityID, "p4rt:3", "The switch entity id should be the same")
	assert.Equal(result3.Networks[0].Switches[2].P4RTPort, 50003, "The switch p4rt server endpoint should be the same")

	assert.Equal(len(result3.Networks[0].Switches[3].Ports), 2, "The number of networks should be the same")
	// verify the entity_id and p4rt_server_endpoint for the switch
	assert.Equal(result3.Networks[0].Switches[3].EntityID, "p4rt:4", "The switch entity id should be the same")
	assert.Equal(result3.Networks[0].Switches[3].P4RTPort, 50004, "The switch p4rt server endpoint should be the same")

	// TEST FOUR: The same topology as test one, but the link is unidirectional
	result4 := Convert(reader.ReadFile("../test_input_files/hr-4.yaml"))
	assert.Equal(len(result4.Networks), 1, "The number of networks should be the same")

	// for the singular network, check switches and links
	assert.Equal(len(result4.Networks[0].Links), 1, "The number of links should be the same")
	assert.Equal(len(result4.Networks[0].Switches), 2, "The number of switches should be the same")
	// verify the entity_id and display_name for the network
	assert.Equal(result4.Networks[0].EntityID, "network-layer:0/underlay-1", "The network entity id should be the same")
	assert.Equal(result4.Networks[0].DisplayName, "underlay-1", "The network display name should be the same")

	// for the 2 switches, check ports
	assert.Equal(len(result4.Networks[0].Switches[0].Ports), 1, "The number of networks should be the same")
	// verify the entity_id and p4rt_server_endpoint for the switch
	assert.Equal(result4.Networks[0].Switches[0].EntityID, "p4rt:1", "The switch entity id should be the same")
	assert.Equal(result4.Networks[0].Switches[0].P4RTPort, 50002, "The switch p4rt server endpoint should be the same")

	assert.Equal(len(result4.Networks[0].Switches[1].Ports), 1, "The number of networks should be the same")
	// verify the entity_id and p4rt_server_endpoint for the switch
	assert.Equal(result4.Networks[0].Switches[1].EntityID, "p4rt:2", "The switch entity id should be the same")
	assert.Equal(result4.Networks[0].Switches[1].P4RTPort, 50001, "The switch p4rt server endpoint should be the same")
}
