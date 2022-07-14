// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package reader

import (
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

// Underlay stores all the networks
type Underlay struct {
	Networks []Network `yaml:"underlay_networks,flow"`
}

// Network stores information pertaining each network
type Network struct {
	EntityID    string     `yaml:"entity_id"`
	DisplayName string     `yaml:"display_name"`
	Switches    []Switches `yaml:"switches,flow"`
	Links       []Link     `yaml:"links,flow"`
}

// Switches stores info for each switch
type Switches struct {
	EntityID           string  `yaml:"entity_id"`
	DisplayName        string  `yaml:"display_name"`
	ModelID            string  `yaml:"model_id"`
	Role               string  `yaml:"role"`
	ManagementEndpoint string  `yaml:"management_endpoint"`
	P4RTServerEndpoint string  `yaml:"p4rt_server_endpoint"`
	TLSInsecure        int     `yaml:"tls_insecure"`
	Ports              []Ports `yaml:"ports,flow"`
}

// Ports contains the port info
type Ports struct {
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
	LinkType    string `yaml:"link_type"`
}

var log = logging.GetLogger()

// ReadFile converts the human-readable file into the struct system above
func ReadFile(filename string) Underlay {
	// reading in the human-readable schema
	var result Underlay
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
