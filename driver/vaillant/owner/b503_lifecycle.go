package owner

import (
	"context"
	"sync"

	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/owner/b503session"
)

type B503Runtime struct {
	lifecycleMu sync.Mutex
	manager     *b503session.Manager
}

func (runtime *B503Runtime) EBusDriverActivated(generation uint64) {
	if runtime == nil || runtime.manager == nil || generation == 0 {
		return
	}
	runtime.lifecycleMu.Lock()
	defer runtime.lifecycleMu.Unlock()
	if generation <= runtime.manager.LifecycleTransportKey().TransportEpoch {
		return
	}
	runtime.manager.OnEpochAdvance(context.Background(), generation)
}

func (runtime *B503Runtime) EBusDriverWithdrawn(generation uint64) {
	if runtime == nil || runtime.manager == nil || generation == 0 {
		return
	}
	runtime.disconnectIfCurrent(b503session.TransportKey{
		AdapterInstanceID: runtime.manager.LifecycleTransportKey().AdapterInstanceID,
		TransportEpoch:    generation,
	})
}

func (runtime *B503Runtime) disconnectIfCurrent(expected b503session.TransportKey) {
	if runtime == nil || runtime.manager == nil {
		return
	}
	runtime.lifecycleMu.Lock()
	defer runtime.lifecycleMu.Unlock()
	if runtime.manager.LifecycleTransportKey() != expected {
		return
	}
	runtime.manager.OnTransportDisconnect()
}
