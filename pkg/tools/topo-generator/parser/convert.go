// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package parser

import (
	"strings"
	"regexp"
	"strconv"
	"github.com/google/uuid"
	"github.com/onosproject/onos-topo/pkg/tools/topo-generator/reader"
)

type Underlay struct {
	Networks	[]Network 
}

type Network struct {
	Entity_Id		string      
	Name			string		
    	Display_Name		string  
	Switches		[]Switches 
	Links			[]Link		
}

type Switches struct {
    	Entity_Id		string          
	Name			string
	Model_Id		string	        
	Role			string	        
	P4RT_Address		string
	P4RT_Port		int
	Insecure		bool
	Ports			[]Ports 
}

type Ports struct {
	Entity_Id		string	
	Name			string
	Display_Name		string	
	Speed			string	
	Port_Number		int		
	Channel_Number		int		
}

type Link struct {
	Source				string	
	Source_Name			string
	Destination			string	
	Dest_Name			string
	Link_Type			string
	// for both unidirectional and bidirectional
	URI				string
	URI_Name			string
	UUID_Source			string
	UUID_Dest			string
	UUID_Source_Name		string
	UUID_Dest_Name			string
	// for bi-directional case (default)
	URI1 				string
	URI1_Name 			string
	UUID_Source1			string
	UUID_Dest1			string
	UUID_Source1_Name		string
	UUID_Dest1_Name			string
}

func Convert(result reader.Underlay) Underlay {
	var underlay Underlay
	var networks []Network
	
	// writing the entity-kind-relationship file
	reg, _ := regexp.Compile("[/:]+")

	for _, n := range result.Networks {
		// network
		var network Network
		network.Entity_Id = n.Entity_Id
		network.Display_Name = n.Display_Name
		network.Name = reg.ReplaceAllString(n.Entity_Id, ".")
		var switches []Switches
		var links []Link
		
		// switches
		for _, s := range n.Switches {
			var sw Switches
			sw.Entity_Id = s.Entity_Id
			sw.Name = reg.ReplaceAllString(s.Entity_Id, ".")
			sw.Model_Id = s.Model_Id
			sw.Role = s.Role
			split := strings.Split(s.P4RT_Server_Endpoint, ":")
			sw.P4RT_Address = split[0]
			intVar, _ := strconv.Atoi(split[1])
			sw.P4RT_Port = intVar
			// default is false
			if s.TLS_insecure == 0 {
				sw.Insecure = true
			}
			var ports []Ports

			// ports
			for _, p := range s.Ports {
				var port Ports
				port.Entity_Id = p.Entity_Id
				port.Display_Name = p.Display_Name
				port.Name = reg.ReplaceAllString(p.Entity_Id, ".")
				port.Speed = p.Speed
				port.Port_Number = p.Port_Number
				port.Channel_Number = p.Channel_Number
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
			link.Source_Name = reg.ReplaceAllString(l.Source, ".")
			link.Destination = l.Destination
			link.Dest_Name = reg.ReplaceAllString(l.Destination, ".")
			link.Link_Type = l.Link_Type
			link.URI = l.Source + "-" + l.Destination
			link.URI_Name = reg.ReplaceAllString(link.URI, ".")
			link.UUID_Source = "uuid:" + uuid.New().String()
			link.UUID_Source_Name = reg.ReplaceAllString(link.UUID_Source, ".")
			link.UUID_Dest = "uuid:" + uuid.New().String()
			link.UUID_Dest_Name = reg.ReplaceAllString(link.UUID_Dest, ".")
			
			// handles whether the link is  bidirectional
			if l.Link_Type == "" {
				link.URI1 = l.Destination + "-" + l.Source
				link.URI1_Name = reg.ReplaceAllString(link.URI1, ".")
				link.UUID_Dest1 = "uuid:" + uuid.New().String()
				link.UUID_Dest1_Name = reg.ReplaceAllString(link.UUID_Dest1, ".")
				link.UUID_Source1 = "uuid:" + uuid.New().String()
				link.UUID_Source1_Name = reg.ReplaceAllString(link.UUID_Source1, ".")
			}
			
			links = append(links, link)
		}
		
		network.Links = links
		networks = append(networks, network)
	}

	underlay.Networks = networks
	return underlay
}
