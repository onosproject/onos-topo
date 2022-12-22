// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

// Package manager is is the main coordinator for the ONOS topology subsystem.
package manager

import (
	"github.com/atomix/go-sdk/pkg/client"
	"github.com/atomix/go-sdk/pkg/primitive"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-lib-go/pkg/northbound"
	service "github.com/onosproject/onos-topo/pkg/northbound"
	"github.com/onosproject/onos-topo/pkg/store"
)

var log = logging.GetLogger("manager")

// Config is a manager configuration
type Config struct {
	AtomixClient primitive.Client
	ServiceFlags *cli.ServiceEndpointFlags
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
	cli.Daemon
	Config Config
}

// Start starts the manager
func (m *Manager) Start() error {
	log.Info("Starting Manager")
	err := m.startNorthboundServer()
	if err != nil {
		return err
	}
	return nil
}

// Stop stops the channels and manager related objects
func (m *Manager) Stop() {
	log.Info("Stopping Manager")
}

// startNorthboundServer starts the northbound gRPC server
func (m *Manager) startNorthboundServer() error {
	s := northbound.NewServer(cli.ServerConfigFromFlags(m.Config.ServiceFlags, northbound.SecurityConfig{}))
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
