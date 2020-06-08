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
// WITHOUT WARRANTIES OR CONDITIONS OF ANY etype, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package topo

import "errors"

// AttrKind ...
type AttrKind string

const (
	// Role is the role of an entity
	Role AttrKind = "AK_ROLE"
	// Revision is the revision of an entity
	Revision AttrKind = "AK_REVISION"
)

// AttrVal ...
type AttrVal string

// IsAttrKindValid validates Entity Type
func (ak AttrKind) IsAttrKindValid() error {
	switch ak {
	case Role, Revision:
		return nil
	}
	return errors.New("Inalid entity type")
}
