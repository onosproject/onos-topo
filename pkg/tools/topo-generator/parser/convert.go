// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package parser

import (
	"github.com/google/uuid"
	"github.com/onosproject/onos-topo/pkg/tools/topo-generator/reader"
	"regexp"
	"strconv"
	"strings"
)

// Underlay keeps track of the networks
type Underlay struct {
	Networks []Network
}

// Network contains all the info required for the e-k-r file
type Network struct {
	EntityID    string
	Name        string
	DisplayName string
	Switches    []Switches
	Links       []Link
}

// Switches contains all required switch info
type Switches struct {
	EntityID    string
	Name        string
	ModelID     string
	Role        string
	P4RTAddress string
	P4RTPort    int
	Insecure    bool
	Ports       []Ports
}

// Ports contains relevant port information
type Ports struct {
	EntityID      string
	Name          string
	DisplayName   string
	Speed         string
	PortNumber    int
	ChannelNumber int
}

// Link contains information for both unidirectional and bidirectional links
type Link struct {
	Source      string
	SourceName  string
	Destination string
	DestName    string
	LinkType    string
	// for both unidirectional and bidirectional
	URI       string
	URIName   string
	UUID1     string
	UUID2     string
	UUID1Name string
	UUID2Name string
	// for bi-directional case (default)
	FlippedURI     string
	FlippedURIName string
	UUID3          string
	UUID4          string
	UUID3Name      string
	UUID4Name      string
}

// Convert takes the struct system from reader and converts it to these structs
func Convert(result reader.Underlay) Underlay {
	var underlay Underlay
	var networks []Network

	// writing the entity-kind-relationship file
	reg, _ := regexp.Compile("[/:]+")

	for _, n := range result.Networks {
		// network
		var network Network
		network.EntityID = n.EntityID
		network.DisplayName = n.DisplayName
		network.Name = reg.ReplaceAllString(n.EntityID, ".")
		var switches []Switches
		var links []Link

		// switches
		for _, s := range n.Switches {
			var sw Switches
			sw.EntityID = s.EntityID
			sw.Name = reg.ReplaceAllString(s.EntityID, ".")
			sw.ModelID = s.ModelID
			sw.Role = s.Role
			split := strings.Split(s.P4RTServerEndpoint, ":")
			sw.P4RTAddress = split[0]
			intVar, _ := strconv.Atoi(split[1])
			sw.P4RTPort = intVar
			// default is false
			if s.TLSInsecure == 0 {
				sw.Insecure = true
			}
			var ports []Ports

			// ports
			for _, p := range s.Ports {
				var port Ports
				port.EntityID = p.EntityID
				port.DisplayName = p.DisplayName
				port.Name = reg.ReplaceAllString(p.EntityID, ".")
				port.Speed = p.Speed
				port.PortNumber = p.PortNumber
				port.ChannelNumber = p.ChannelNumber
				ports = append(ports, port)
			}

			sw.Ports = ports
			switches = append(switches, sw)
		}

		network.Switches = switches

		// links
		for _, l := range n.Links {
			var link Link
			link.Source = l.Source
			link.SourceName = reg.ReplaceAllString(l.Source, ".")
			link.Destination = l.Destination
			link.DestName = reg.ReplaceAllString(l.Destination, ".")
			link.LinkType = l.LinkType
			link.URI = l.Source + "-" + l.Destination
			link.URIName = reg.ReplaceAllString(link.URI, ".")
			link.UUID1 = "uuid:" + uuid.New().String()
			link.UUID1Name = reg.ReplaceAllString(link.UUID1, ".")
			link.UUID2 = "uuid:" + uuid.New().String()
			link.UUID2Name = reg.ReplaceAllString(link.UUID2, ".")

			// handles whether the link is  bidirectional
			if l.LinkType == "" {
				link.FlippedURI = l.Destination + "-" + l.Source
				link.FlippedURIName = reg.ReplaceAllString(link.FlippedURI, ".")
				link.UUID3 = "uuid:" + uuid.New().String()
				link.UUID3Name = reg.ReplaceAllString(link.UUID3, ".")
				link.UUID4 = "uuid:" + uuid.New().String()
				link.UUID4Name = reg.ReplaceAllString(link.UUID4, ".")
			}

			links = append(links, link)
		}

		network.Links = links
		networks = append(networks, network)
	}

	underlay.Networks = networks
	return underlay
}
