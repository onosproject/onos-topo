// Copyright 2020-present Open Networking Foundation.
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

package bulk

import (
	"gotest.tools/assert"
	"testing"
)

func Test_LoadConfig1(t *testing.T) {
	deviceConfig = nil
	config, err := GetDeviceConfig("topo-load-example.yaml")
	assert.NilError(t, err, "Unexpected error loading towerConfig")
	assert.Equal(t, 2, len(config.Devices), "Unexpected number of towers")

	tower1 := config.Devices[0]
	assert.Equal(t, "315010-0001420", string(tower1.ID))
	assert.Equal(t, 6, len(tower1.Attributes))
}
