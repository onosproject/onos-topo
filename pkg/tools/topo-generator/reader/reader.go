// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package reader

import (
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

// NetworkLayer stores all the networks
type NetworkLayer struct {
	Networks []Network `yaml:"underlay_networks,flow"`
}

// Network stores information pertaining each network
type Network struct {
	EntityID    string   `yaml:"entity_id"`
	DisplayName string   `yaml:"display_name"`
	Switches    []Switch `yaml:"switches,flow"`
	Links       []Link   `yaml:"links,flow"`
}

// Switch stores info for each switch
type Switch struct {
	EntityID           string `yaml:"entity_id"`
	DisplayName        string `yaml:"display_name"`
	ModelID            string `yaml:"model_id"`
	Role               string `yaml:"role"`
	ManagementEndpoint string `yaml:"management_endpoint"`
	P4RTServerEndpoint string `yaml:"p4rt_server_endpoint"`
	TLSInsecure        int    `yaml:"tls_insecure"`
	Ports              []Port `yaml:"ports,flow"`
}

// Port contains the port info
type Port struct {
	EntityID      string `yaml:"entity_id"`
	DisplayName   string `yaml:"display_name"`
	Speed         string `yaml:"speed"`
	PortNumber    int    `yaml:"port_number"`
	ChannelNumber int    `yaml:"channel_number"`
}

// Link handles links between ports
type Link struct {
	Source      string `yaml:"src"`
	Destination string `yaml:"dst"`
	LinkType    string `yaml:"type"`
}

var log = logging.GetLogger()

// ReadFile converts the human-readable file into the struct system above
func ReadFile(filename string) NetworkLayer {
	// reading in the human-readable schema
	var result NetworkLayer
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Error(err.Error())
		return result
	}
	err = yaml.Unmarshal(content, &result)
	if err != nil {
		log.Error("Failed to parse file ", err)
	}
	return result
}
