// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/gogo/protobuf/types"
	topoapi "github.com/onosproject/onos-api/go/onos/topo"
	"github.com/onosproject/onos-lib-go/pkg/certs"
	"github.com/onosproject/onos-lib-go/pkg/grpc/retry"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"sync"
)

// Utility to populate onos-topo with a large scale topology for testing purposes.
func main() {
	createPodDeployment()
}

var log = logging.GetLogger("topo-scale")

const (
	superSpineCount      = 4
	portsPerSuperSpine   = 64
	spineSuperSpineTrunk = 2

	podCount       = 16
	spinePerPod    = 2
	rackPerPod     = 6
	serversPerPod  = 12
	leafSpineTrunk = 4

	maxVMsPerServer = 20
	portsPerSpine   = 32
	portsPerLeaf    = 32
	portsPerIPU     = 2
)

func assertNoError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Builder hods state to assist generating various fabric topologies
type Builder struct {
	lock     sync.RWMutex
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
	b.lock.Lock()
	defer b.lock.Unlock()
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
	labels := Labels{"tier": "superspines"}
	for ssi := 1; ssi <= superSpineCount; ssi++ {
		createDevice(client, "switch", fmt.Sprintf("superspine-%d", ssi), portsPerSuperSpine, labels)
	}

	var wg sync.WaitGroup
	wg.Add(podCount)
	for pi := 1; pi <= podCount; pi++ {
		createPod(client, pi, fmt.Sprintf("pod-%02d", pi), builder, &wg)
	}
	wg.Wait()

	for ssi := 1; ssi <= superSpineCount; ssi++ {
		podNumber := (ssi-1)*podCount/superSpineCount + 1
		podName := fmt.Sprintf("pod-%02d", podNumber)
		rackName := fmt.Sprintf("rack-%02d-1", podNumber)
		superspineName := fmt.Sprintf("superspine-%d", ssi)
		assertNoError(addLabels(client, superspineName, Labels{"pod": podName, "rack": rackName}))
		assertNoError(createRelation(client, rackName, superspineName, "contains"))
	}
}

// Labels is a map of labeled values
type Labels map[string]string

func createPod(client topoapi.TopoClient, podID int, podName string, builder *Builder, wg *sync.WaitGroup) {
	log.Infof("Creating pod %s...", podName)
	podLabels := Labels{"pod": podName}
	err := createEntity(client, podName, "pod", nil, podLabels)
	assertNoError(err)

	for si := 1; si <= spinePerPod; si++ {
		spineName := fmt.Sprintf("spine-%02d-%d", podID, si)
		createDevice(client, "switch", spineName, portsPerSpine, podLabels)
		assertNoError(createRelation(client, podName, spineName, "contains"))

		for ssi := 1; ssi <= superSpineCount; ssi++ {
			superSpineName := fmt.Sprintf("superspine-%d", ssi)
			createLinkTrunk(client, superSpineName, spineName, spineSuperSpineTrunk, "link", podLabels, builder)
		}
	}

	go func() {
		for ri := 1; ri <= rackPerPod; ri++ {
			rackName := fmt.Sprintf("rack-%02d-%d", podID, ri)
			rackLabels := Labels{"pod": podName, "rack": rackName}
			createRack(client, podID, ri, rackName, rackLabels)
			assertNoError(createRelation(client, podName, rackName, "contains"))

			leafName := fmt.Sprintf("leaf-%02d-%d", podID, ri)
			for si := 1; si <= spinePerPod; si++ {
				spineName := fmt.Sprintf("spine-%02d-%d", podID, si)
				createLinkTrunk(client, spineName, leafName, leafSpineTrunk, "link", rackLabels, builder)
			}
		}

		for si := 1; si <= spinePerPod; si++ {
			rackID := (si-1)*rackPerPod/spinePerPod + 1
			rackName := fmt.Sprintf("rack-%02d-%d", podID, rackID)
			spineName := fmt.Sprintf("spine-%02d-%d", podID, si)
			assertNoError(addLabels(client, spineName, Labels{"rack": rackName}))
			assertNoError(createRelation(client, rackName, spineName, "contains"))
		}

		wg.Done()
	}()
}

func addLabels(client topoapi.TopoClient, id string, labels Labels) error {
	ctx := context.Background()
	resp, err := client.Get(ctx, &topoapi.GetRequest{ID: topoapi.ID(id)})
	if err != nil {
		return err
	}
	if len(resp.Object.Labels) == 0 {
		resp.Object.Labels = labels
	} else {
		for l, v := range labels {
			resp.Object.Labels[l] = v
		}
	}
	_, err = client.Update(ctx, &topoapi.UpdateRequest{Object: resp.Object})
	return err
}

func createRack(client topoapi.TopoClient, podID int, rackID int, rackName string, labels Labels) {
	log.Infof("Creating rack %s...", rackName)
	err := createEntity(client, rackName, "rack", nil, labels)
	assertNoError(err)

	leafName := fmt.Sprintf("leaf-%02d-%d", podID, rackID)
	createDevice(client, "leaf", leafName, portsPerLeaf, labels)
	err = createRelation(client, rackName, leafName, "contains")
	assertNoError(err)

	for si := 1; si <= serversPerPod; si++ {
		serverName := fmt.Sprintf("server-%02d-%d-%02d", podID, rackID, si)
		createServer(client, podID, rackID, si, serverName, labels)
		err = createRelation(client, rackName, serverName, "contains")
		assertNoError(err)
	}
}

func createServer(client topoapi.TopoClient, podID int, rackID int, serverID int, serverName string, labels Labels) {
	err := createEntity(client, serverName, "server", nil, labels)
	assertNoError(err)

	ipuName := fmt.Sprintf("ipu-%02d-%d-%02d", podID, rackID, serverID)
	createDevice(client, "ipu", ipuName, portsPerIPU+maxVMsPerServer, labels)
	err = createRelation(client, serverName, ipuName, "contains")
	createBidirectionalLink(client, ipuName, fmt.Sprintf("leaf-%02d-%d", podID, rackID), "link", labels)
	assertNoError(err)

	for i := 1; i <= maxVMsPerServer; i++ {
		vmName := fmt.Sprintf("vm-%02d-%d-%02d-%02d", podID, rackID, serverID, i)
		err := createEntity(client, vmName, "vm", nil, labels)
		assertNoError(err)
		err = createRelation(client, serverName, vmName, "contains")
		assertNoError(err)

		createBidirectionalLink(client, fmt.Sprintf("%s/%d", ipuName, i), vmName, "edge-link", labels)
		assertNoError(err)
	}
}

func createDevice(client topoapi.TopoClient, kindID string, deviceName string, portCount int, labels Labels) {
	log.Infof("Creating %s %s with %d ports...", kindID, deviceName, portCount)
	err := createEntity(client, deviceName, kindID, nil, labels)
	assertNoError(err)

	for pn := 1; pn <= portCount; pn++ {
		portName := fmt.Sprintf("%s/%d", deviceName, pn)
		err = createEntity(client, portName, "port", nil, labels)
		assertNoError(err)
		err = createRelation(client, deviceName, portName, "has")
		assertNoError(err)
	}
}

func createBidirectionalLink(client topoapi.TopoClient, id1 string, id2 string, kindID string, labels Labels) {
	createLink(client, id1, id2, kindID, labels)
	createLink(client, id2, id1, kindID, labels)
}

func createLink(client topoapi.TopoClient, id1 string, id2 string, kindID string, labels Labels) {
	linkName := fmt.Sprintf("%s-%s", id1, id2)
	err := createEntity(client, linkName, kindID, nil, labels)
	assertNoError(err)
	err = createRelation(client, id1, linkName, "originates")
	assertNoError(err)
	err = createRelation(client, id2, linkName, "terminates")
	assertNoError(err)
}

// Create a trunk of specified number of links between two Devices
func createLinkTrunk(client topoapi.TopoClient, src string, tgt string, count int, kindID string, labels Labels, builder *Builder) {
	log.Infof("Creating trunk with %d links between %s and %s...", count, src, tgt)
	for i := 0; i < count; i++ {
		createBidirectionalLink(client, builder.NextDevicePortID(src), builder.NextDevicePortID(tgt), kindID, labels)
	}
}

// createEntity creates an entity object
func createEntity(client topoapi.TopoClient, id string, kindID string, aspectList []*types.Any, labels map[string]string) error {
	aspects := map[string]*types.Any{}
	for _, aspect := range aspectList {
		aspects[aspect.TypeUrl] = aspect
	}
	_, err := client.Create(context.Background(), &topoapi.CreateRequest{
		Object: &topoapi.Object{
			ID:      topoapi.ID(id),
			Type:    topoapi.Object_ENTITY,
			Aspects: aspects,
			Obj:     &topoapi.Object_Entity{Entity: &topoapi.Entity{KindID: topoapi.ID(kindID)}},
			Labels:  labels,
		},
	})
	return err
}

// CreateRelation creates a relation object
func createRelation(client topoapi.TopoClient, src string, tgt string, kindID string) error {
	_, err := client.Create(context.Background(), &topoapi.CreateRequest{
		Object: &topoapi.Object{
			ID:   topoapi.ID(src + tgt + kindID),
			Type: topoapi.Object_RELATION,
			Obj: &topoapi.Object_Relation{
				Relation: &topoapi.Relation{
					SrcEntityID: topoapi.ID(src),
					TgtEntityID: topoapi.ID(tgt),
					KindID:      topoapi.ID(kindID),
				},
			},
		},
	})
	return err
}
