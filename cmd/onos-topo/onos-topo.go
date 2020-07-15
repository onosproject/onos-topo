// Copyright 2019-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/*
Package onos-topo is the main entry point to the ONOS topology subsystem.

Arguments

-caPath <the location of a CA certificate>

-keyPath <the location of a client private key>

-certPath <the location of a client certificate>


See ../../docs/run.md for how to run the application.
*/
package main

import (
	"flag"
	"github.com/onosproject/onos-lib-go/pkg/auth"
	"os"

	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-lib-go/pkg/northbound"
	"github.com/onosproject/onos-topo/pkg/manager"
	"github.com/onosproject/onos-topo/pkg/northbound/admin"
	"github.com/onosproject/onos-topo/pkg/northbound/device"
	"github.com/onosproject/onos-topo/pkg/northbound/diags"
	"github.com/onosproject/onos-topo/pkg/northbound/topo"
)

var log = logging.GetLogger("main")

// The main entry point
func main() {
	caPath := flag.String("caPath", "", "path to CA certificate")
	keyPath := flag.String("keyPath", "", "path to client private key")
	certPath := flag.String("certPath", "", "path to client certificate")
	flag.Parse()

	log.Info("Starting onos-topo")

	mgr, err := manager.NewManager()
	if err != nil {
		log.Fatal("Unable to load onos-topo ", err)
	} else {
		mgr.Run()
		err = startServer(*caPath, *keyPath, *certPath)
		if err != nil {
			log.Fatal("Unable to start onos-topo ", err)
		}
	}
}

// Creates gRPC server and registers various services; then serves.
func startServer(caPath string, keyPath string, certPath string) error {
	oidcServer := os.Getenv(auth.OIDCServerURL)
	securityConfig := northbound.SecurityConfig{}
	if oidcServer != "" {
		log.Infof("Authentication checking enabled. Using oidc server %s", oidcServer)
		securityConfig.AuthenticationEnabled = true
	}

	s := northbound.NewServer(northbound.NewServerCfg(caPath, keyPath, certPath,
		5150, true, securityConfig))
	s.AddService(admin.Service{})
	s.AddService(diags.Service{})
	s.AddService(logging.Service{})

	deviceService, err := device.NewService()
	if err != nil {
		return err
	}
	s.AddService(deviceService)

	topoService, err := topo.NewService()
	if err != nil {
		return err
	}
	s.AddService(topoService)

	return s.Serve(func(started string) {
		log.Info("Started NBI on ", started)
	})
}
