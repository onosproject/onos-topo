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

package topo

import (
	"github.com/onosproject/helmit/pkg/helm"
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
const testName = "topo-test"

// SetupTestSuite sets up the onos-topo test suite
func (s *TestSuite) SetupTestSuite() error {
	err := helm.Chart(onostest.ControllerChartName, onostest.AtomixChartRepo).
		Release(onostest.AtomixName(testName, onosTopoComponentName)).
		Set("scope", "Namespace").
		Install(true)
	if err != nil {
		return err
	}

	err = helm.Chart(onostest.RaftStorageControllerChartName, onostest.AtomixChartRepo).
		Release(onostest.RaftReleaseName(onosTopoComponentName)).
		Set("scope", "Namespace").
		Install(true)
	if err != nil {
		return err
	}

	err = helm.Chart(onosTopoComponentName, onostest.OnosChartRepo).
		Release(onosTopoComponentName).
		Set("image.tag", "latest").
		Set("storage.controller", onostest.AtomixController(testName, onosTopoComponentName)).
		Install(true)
	if err != nil {
		return err
	}
	return nil
}
