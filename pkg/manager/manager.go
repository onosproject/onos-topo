// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

// Package manager is the main coordinator for the ONOS topology subsystem.
package manager

import (
	"github.com/atomix/go-sdk/pkg/client"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-lib-go/pkg/northbound"
	service "github.com/onosproject/onos-topo/pkg/northbound"
	"github.com/onosproject/onos-topo/pkg/store"
)

var log = logging.GetLogger("manager")

// Config is a manager configuration
type Config struct {
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
	Config    Config
	topoStore store.Store
}

// Start starts the manager
func (m *Manager) Start() error {
	log.Info("Starting Manager")

	var err error
	if m.topoStore, err = store.NewAtomixStore(client.NewClient()); err != nil {
		return err
	}

	s := northbound.NewServer(cli.ServerConfigFromFlags(m.Config.ServiceFlags, northbound.SecurityConfig{}))
	s.AddService(logging.Service{})
	s.AddService(service.NewService(m.topoStore))
	return s.StartInBackground()
}

// Stop stops the channels and manager related objects
func (m *Manager) Stop() {
	log.Info("Stopping Manager")
	_ = m.topoStore.Close()
}
