package owner

import (
	"context"
	"sync"

	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/owner/b503session"
)

// B503Runtime is the b503 runtime type of the Vaillant owner.
type B503Runtime struct {
	lifecycleMu sync.Mutex
	manager     *b503session.Manager
}

// EBusDriverActivated implements the ebus driver activated step of the Vaillant owner.
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

// EBusDriverWithdrawn implements the ebus driver withdrawn step of the Vaillant owner.
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
