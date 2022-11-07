// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

// Package manager is is the main coordinator for the ONOS topology subsystem.
package manager

import (
	"github.com/atomix/go-sdk/pkg/client"
	"github.com/atomix/go-sdk/pkg/primitive"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-lib-go/pkg/northbound"
	service "github.com/onosproject/onos-topo/pkg/northbound"
	"github.com/onosproject/onos-topo/pkg/store"
)

var log = logging.GetLogger()

// Config is a manager configuration
type Config struct {
	CAPath       string
	KeyPath      string
	CertPath     string
	GRPCPort     int
	AtomixClient primitive.Client
}

// NewManager creates a new manager
func NewManager(config Config) *Manager {
	log.Info("Creating Manager")
	return &Manager{
		Config: config,
	}
}

// Manager single point of entry for the topology system.
type Manager struct {
	Config Config
}

// Run starts a synchronizer based on the devices and the northbound services.
func (m *Manager) Run() {
	log.Info("Starting Manager")
	if err := m.Start(); err != nil {
		log.Fatal("Unable to run Manager", err)
	}
}

// Start starts the manager
func (m *Manager) Start() error {
	err := m.startNorthboundServer()
	if err != nil {
		return err
	}
	return nil
}

// startNorthboundServer starts the northbound gRPC server
func (m *Manager) startNorthboundServer() error {
	s := northbound.NewServer(northbound.NewServerCfg(
		m.Config.CAPath,
		m.Config.KeyPath,
		m.Config.CertPath,
		int16(m.Config.GRPCPort),
		true,
		northbound.SecurityConfig{}))

	if m.Config.AtomixClient == nil {
		m.Config.AtomixClient = client.NewClient()
	}

	topoStore, err := store.NewAtomixStore(m.Config.AtomixClient)
	if err != nil {
		return err
	}

	s.AddService(logging.Service{})
	s.AddService(service.NewService(topoStore))

	doneCh := make(chan error)
	go func() {
		err := s.Serve(func(started string) {
			log.Info("Started NBI on ", started)
			close(doneCh)
		})
		if err != nil {
			doneCh <- err
		}
	}()
	return <-doneCh
}

// Close kills the channels and manager related objects
func (m *Manager) Close() {
	log.Info("Closing Manager")
}
