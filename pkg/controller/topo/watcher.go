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
	topoapi "github.com/onosproject/onos-topo/api/topo"
	topostore "github.com/onosproject/onos-topo/pkg/store/topo"
	"sync"
)

const queueSize = 100

// Watcher is a subscription watcher
type Watcher struct {
	topo   topostore.Store
	cancel context.CancelFunc
	mu     sync.Mutex
}

// Start starts the subscription watcher
func (w *Watcher) Start(ch chan<- controller.ID) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.cancel != nil {
		return nil
	}

	subCh := make(chan topoapi.Event, queueSize)
	ctx, cancel := context.WithCancel(context.Background())
	err := w.topo.Watch(ctx, subCh)
	if err != nil {
		cancel()
		return err
	}
	w.cancel = cancel

	go func() {
		for request := range subCh {
			ch <- controller.NewID(request.Object.ID)
		}
		close(ch)
	}()
	return nil
}

// Stop stops the subscription watcher
func (w *Watcher) Stop() {
	w.mu.Lock()
	if w.cancel != nil {
		w.cancel()
		w.cancel = nil
	}
	w.mu.Unlock()
}

var _ controller.Watcher = &Watcher{}
