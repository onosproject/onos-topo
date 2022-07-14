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

// Link contains the entity information
type Link struct {
	URI                string
	URIName            string
	OriginatesRelation Originates
	TerminatesRelation Terminates
}

// Originates contains the relation information
type Originates struct {
	URI        string
	URIName    string
	UUID       string
	UUIDName   string
	Source     string
	SourceName string
}

// Terminates contains the relation information
type Terminates struct {
	URI         string
	URIName     string
	UUID        string
	UUIDName    string
	Destination string
	DestName    string
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
			link.URI = l.Source + "-" + l.Destination
			link.URIName = reg.ReplaceAllString(link.URI, ".")

			var originate Originates
			originate.URI = link.URI
			originate.URIName = link.URIName
			originate.Source = l.Source
			originate.SourceName = reg.ReplaceAllString(l.Source, ".")
			originate.UUID = "uuid:" + uuid.New().String()
			originate.UUIDName = reg.ReplaceAllString(originate.UUID, ".")
			link.OriginatesRelation = originate

			var terminate Terminates
			terminate.URI = link.URI
			terminate.URIName = link.URIName
			terminate.Destination = l.Destination
			terminate.DestName = reg.ReplaceAllString(l.Destination, ".")
			terminate.UUID = "uuid:" + uuid.New().String()
			terminate.UUIDName = reg.ReplaceAllString(terminate.UUID, ".")
			link.TerminatesRelation = terminate

			links = append(links, link)

			// handles whether the link is  bidirectional
			if l.LinkType == "" {
				var link Link
				link.URI = l.Destination + "-" + l.Source
				link.URIName = reg.ReplaceAllString(link.URI, ".")

				var originate Originates
				originate.URI = link.URI
				originate.URIName = link.URIName
				originate.Source = l.Destination
				originate.SourceName = reg.ReplaceAllString(l.Destination, ".")
				originate.UUID = "uuid:" + uuid.New().String()
				originate.UUIDName = reg.ReplaceAllString(originate.UUID, ".")
				link.OriginatesRelation = originate

				var terminate Terminates
				terminate.URI = link.URI
				terminate.URIName = link.URIName
				terminate.Destination = l.Source
				terminate.DestName = reg.ReplaceAllString(l.Source, ".")
				terminate.UUID = "uuid:" + uuid.New().String()
				terminate.UUIDName = reg.ReplaceAllString(terminate.UUID, ".")
				link.TerminatesRelation = terminate

				links = append(links, link)
			}
		}

		network.Links = links
		networks = append(networks, network)
	}

	underlay.Networks = networks
	return underlay
}
