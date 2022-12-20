// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"crypto/tls"
	"fmt"
	topoapi "github.com/onosproject/onos-api/go/onos/topo"
	"github.com/onosproject/onos-lib-go/pkg/certs"
	"github.com/onosproject/onos-lib-go/pkg/grpc/retry"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-topo/test/topo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Utility to populate onos-topo with a large scale topology for testing purposes.
func main() {
	createPodDeployment()
}

var log = logging.GetLogger("topo-scale")

const (
	superSpineCount    = 4
	portsPerSuperSpine = 64

	podCount      = 16
	spinePerPod   = 2
	rackPerPod    = 6
	serversPerPod = 12

	maxVMsPerServer = 20
	portsPerSpine   = 32
	portsPerLeaf    = 32
	portsPerIPU     = 2
)

func assertNoError(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}

// Builder hods state to assist generating various fabric topologies
type Builder struct {
	nextPort map[string]int
}

// NewBuilder creates a new topology builder context
func NewBuilder() *Builder {
	return &Builder{
		nextPort: make(map[string]int),
	}
}

// NextDevicePortID reserves the next available port ID and returns it
func (b *Builder) NextDevicePortID(deviceID string) string {
	portNumber, ok := b.nextPort[deviceID]
	if !ok {
		portNumber = 1
	}
	portID := fmt.Sprintf("%s/%d", deviceID, portNumber)
	b.nextPort[deviceID] = portNumber + 1
	return portID
}

// GetClientCredentials returns client credentials
func GetClientCredentials() (*tls.Config, error) {
	cert, err := tls.X509KeyPair([]byte(certs.DefaultClientCrt), []byte(certs.DefaultClientKey))
	if err != nil {
		return nil, err
	}
	return &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}, nil
}

// CreateConnection creates gRPC connection to the topo store
func CreateConnection() (*grpc.ClientConn, error) {
	tlsConfig, err := GetClientCredentials()
	if err != nil {
		return nil, err
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
		grpc.WithUnaryInterceptor(retry.RetryingUnaryClientInterceptor()),
	}

	conn, err := grpc.Dial("localhost:5150", opts...)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func createPodDeployment() {
	conn, err := CreateConnection()
	assertNoError(err)
	client := topoapi.NewTopoClient(conn)
	builder := NewBuilder()

	// Create super-spines
	for ssi := 1; ssi <= superSpineCount; ssi++ {
		createDevice(client, "switch", fmt.Sprintf("superspine-%d", ssi), portsPerSuperSpine)
	}

	for pi := 1; pi <= podCount; pi++ {
		podName := fmt.Sprintf("pod-%02d", pi)
		createPod(client, pi, podName, builder)
	}

	// TODO create rack contains for super-spines

}

func createPod(client topoapi.TopoClient, podID int, podName string, builder *Builder) {
	log.Infof("Creating pod %s...", podName)
	err := topo.CreateEntity(client, podName, "pod", nil)
	assertNoError(err)

	for si := 1; si <= spinePerPod; si++ {
		spineName := fmt.Sprintf("spine-%02d-%d", podID, si)
		createDevice(client, "switch", spineName, portsPerSpine)
		err = topo.CreateRelation(client, podName, spineName, "contains")
		assertNoError(err)

		for ssi := 1; ssi <= superSpineCount; ssi++ {
			superSpineName := fmt.Sprintf("superspine-%d", ssi)
			createLinkTrunk(client, superSpineName, spineName, 2, "link", builder)
		}
	}

	for ri := 1; ri <= rackPerPod; ri++ {
		rackName := fmt.Sprintf("rack-%02d-%d", podID, ri)
		createRack(client, podID, ri, rackName)
		err = topo.CreateRelation(client, podName, rackName, "contains")
		assertNoError(err)

		leafName := fmt.Sprintf("leaf-%02d-%d", podID, ri)
		for si := 1; si <= spinePerPod; si++ {
			spineName := fmt.Sprintf("spine-%02d-%d", podID, si)
			createLinkTrunk(client, spineName, leafName, 2, "link", builder)
		}
	}

	// TODO create rack contains for spines
}

func createRack(client topoapi.TopoClient, podID int, rackID int, rackName string) {
	log.Infof("Creating rack %s...", rackName)
	err := topo.CreateEntity(client, rackName, "rack", nil)
	assertNoError(err)

	leafName := fmt.Sprintf("leaf-%02d-%d", podID, rackID)
	createDevice(client, "leaf", leafName, portsPerLeaf)
	err = topo.CreateRelation(client, rackName, leafName, "contains")
	assertNoError(err)

	for si := 1; si <= serversPerPod; si++ {
		serverName := fmt.Sprintf("server-%02d-%d-%02d", podID, rackID, si)
		createServer(client, podID, rackID, si, serverName)
		err = topo.CreateRelation(client, rackName, serverName, "contains")
		assertNoError(err)
	}
}

func createServer(client topoapi.TopoClient, podID int, rackID int, serverID int, serverName string) {
	err := topo.CreateEntity(client, serverName, "server", nil)
	assertNoError(err)

	ipuName := fmt.Sprintf("ipu-%02d-%d-%02d", podID, rackID, serverID)
	createDevice(client, "ipu", ipuName, portsPerIPU+maxVMsPerServer)
	err = topo.CreateRelation(client, serverName, ipuName, "contains")
	assertNoError(err)

	for i := 1; i <= maxVMsPerServer; i++ {
		vmName := fmt.Sprintf("vm-%02d-%d-%02d-%02d", podID, rackID, serverID, i)
		err := topo.CreateEntity(client, vmName, "vm", nil)
		assertNoError(err)
		err = topo.CreateRelation(client, serverName, vmName, "contains")
		assertNoError(err)

		createBidirectionalLink(client, fmt.Sprintf("%s/%d", ipuName, i), vmName, "edge-link")
		assertNoError(err)
	}
}

func createDevice(client topoapi.TopoClient, kindID string, deviceName string, portCount int) {
	log.Infof("Creating %s %s with %d ports...", kindID, deviceName, portCount)
	err := topo.CreateEntity(client, deviceName, kindID, nil)
	assertNoError(err)

	for pn := 1; pn <= portCount; pn++ {
		portName := fmt.Sprintf("%s/%d", deviceName, pn)
		err = topo.CreateEntity(client, portName, "port", nil)
		assertNoError(err)
		err = topo.CreateRelation(client, deviceName, portName, "has")
		assertNoError(err)
	}
}

func createBidirectionalLink(client topoapi.TopoClient, id1 string, id2 string, kindID string) {
	createLink(client, id1, id2, kindID)
	createLink(client, id2, id1, kindID)
}

func createLink(client topoapi.TopoClient, id1 string, id2 string, kindID string) {
	linkName := fmt.Sprintf("%s-%s", id1, id2)
	err := topo.CreateEntity(client, linkName, kindID, nil)
	assertNoError(err)
	err = topo.CreateRelation(client, id1, linkName, "originates")
	assertNoError(err)
	err = topo.CreateRelation(client, id2, linkName, "terminates")
	assertNoError(err)
}

// Create a trunk of specified number of links between two Devices
func createLinkTrunk(client topoapi.TopoClient, src string, tgt string, count int, kindID string, builder *Builder) {
	log.Infof("Creating trunk with %d links between %s and %s...", count, src, tgt)
	for i := 0; i < count; i++ {
		createBidirectionalLink(client, builder.NextDevicePortID(src), builder.NextDevicePortID(tgt), kindID)
	}
}
