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

package subscription

import (
	"context"
	"github.com/onosproject/onos-lib-go/pkg/controller"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	topoapi "github.com/onosproject/onos-topo/api/topo"
	topostore "github.com/onosproject/onos-topo/pkg/store/topo"
	"time"

	"github.com/onosproject/onos-lib-go/pkg/logging"
)

var log = logging.GetLogger("controller", "subscription")

const defaultTimeout = 30 * time.Second

// NewController returns a new network controller
func NewController(topoStore topostore.Store) *controller.Controller {
	c := controller.NewController("Subscription")
	c.Watch(&Watcher{
		topo: topoStore,
	})
	c.Reconcile(&Reconciler{
		topo: topoStore,
	})
	return c
}

// Reconciler is a topology change reconciler
type Reconciler struct {
	topo topostore.Store
}

// Reconcile reconciles the state of a topology object change
func (r *Reconciler) Reconcile(id controller.ID) (controller.Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	sub, err := r.topo.Get(ctx, id.Value.(topoapi.ID))
	if err != nil {
		if errors.IsNotFound(err) {
			return controller.Result{}, nil
		}
		return controller.Result{}, err
	}

	log.Infof("Reconciling object %+v", sub)
	// TODO: implement this appropriately
	return controller.Result{}, nil
}
