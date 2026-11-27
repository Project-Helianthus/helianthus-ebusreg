package owner

import (
	"context"
	"errors"
	"io"
	"sync"
	"testing"

	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/owner/b503session"
	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
)

func TestAcquisitionTiesHistoryToTheDispatcherRuntime(t *testing.T) {
	raw := &acquisitionTransport{done: make(chan struct{})}
	acquired := NewAcquisition(Config{ScanSource: 0x11}, Substrate{
		Registry:  registry.NewDeviceRegistry(nil),
		Bus:       protocol.NewBus(raw, protocol.BusConfig{}, 1),
		Transport: raw,
	})
	t.Cleanup(func() {
		acquired.History().ResetForRestart()
		_ = raw.Close()
	})
	if acquired.Runtime() == nil || acquired.Runtime().manager != acquired.History() {
		t.Fatal("runtime is not the owner's history session")
	}
	if acquired.dispatcher.runtime != acquired.Runtime() || acquired.dispatcher.mgr != acquired.History() {
		t.Fatal("dispatcher does not keep the shared runtime")
	}
	gate, ok := acquired.dispatcher.readMu.(*dispatcherGate)
	if !ok || gate.mu != &acquired.poller.readMu {
		t.Fatal("dispatcher lock is not the poller exchange lock")
	}
	if !acquired.History().TargetBlocked(ConfiguredMonitorTarget) {
		t.Fatal("configured target is not fenced")
	}
	acquired.History().MarkQualifiedTarget(0x21)
	if !acquired.History().TargetBlocked(0x21) {
		t.Fatal("later qualified target is not fenced")
	}
	calls := 0
	_, err := acquired.History().EnableOperation(context.Background(), ConfiguredMonitorTarget, func(context.Context, byte) b503session.DispatchOutcome {
		calls++
		return b503session.DispatchOutcome{}
	})
	if !errors.Is(err, b503session.ErrCleanupPending) || calls != 0 {
		t.Fatalf("enable = %v calls=%d", err, calls)
	}
}

type acquisitionTransport struct {
	mu   sync.Mutex
	done chan struct{}
	shut bool
}

func (t *acquisitionTransport) ReadByte() (byte, error) {
	<-t.done
	return 0, io.EOF
}

func (t *acquisitionTransport) Write(payload []byte) (int, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.shut {
		return 0, io.ErrClosedPipe
	}
	return len(payload), nil
}

func (t *acquisitionTransport) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.shut {
		return nil
	}
	t.shut = true
	close(t.done)
	return nil
}
