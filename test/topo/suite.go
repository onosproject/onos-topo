// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package topo

import (
	"github.com/onosproject/helmit/pkg/helm"
	"github.com/onosproject/helmit/pkg/input"
	"github.com/onosproject/helmit/pkg/test"
	"github.com/onosproject/onos-test/pkg/onostest"
)

type testSuite struct {
	test.Suite
}

// TestSuite is the onos-topo test suite
type TestSuite struct {
	testSuite
}

const onosTopoComponentName = "onos-topo"

// SetupTestSuite sets up the onos-topo test suite
func (s *TestSuite) SetupTestSuite(c *input.Context) error {

	registry := c.GetArg("registry").String("")
	err := helm.Chart(onosTopoComponentName, onostest.OnosChartRepo).
		Release(onosTopoComponentName).
		Set("image.tag", "latest").
		Set("global.image.registry", registry).
		Install(true)
	if err != nil {
		return err
	}
	return nil
}
