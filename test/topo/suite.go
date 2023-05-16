// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package topo

import (
	"github.com/onosproject/helmit/pkg/test"
	"github.com/onosproject/onos-test/pkg/onostest"
)

// TestSuite is the onos-topo test suite
type TestSuite struct {
	test.Suite
}

const onosTopoComponentName = "onos-topo"

// SetupSuite sets up the onos-topo test suite
func (s *TestSuite) SetupSuite() {
	registry := s.Arg("registry").String()

	install := s.Helm().
		Install(onosTopoComponentName, "onos-umbrella").
		RepoURL(onostest.OnosChartRepo).
		Set("onos-topo.global.image.tag", "latest")

	if registry != "" {
		install.Set("onos-topo.global.image.registry", registry).
			Set("onos-config.global.image.registry", registry).
			Set("onos-umbrella.global.image.registry", registry).
			Set("topo-discovery.global.image.registry", registry).
			Set("device-provisioner.global.image.registry", registry).
			Set("onos-cli.global.image.registry", registry)
	}

	_, err := install.Wait().Get(s.Context())
	s.NoError(err)
}

var _ test.SetupSuite = (*TestSuite)(nil)
