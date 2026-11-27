package owner

// Send failure mapping, source admission and lifecycle withdrawal of the history dispatcher.

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	ebuserrors "github.com/Project-Helianthus/helianthus-ebusgo/errors"
	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/owner/b503session"
)

func TestHistoryDispatcherNAKErrorsMapToUpstreamRPCFailed(t *testing.T) {
	disp, bus, _ := newTestDispatcher(t)
	bus.setErr([2]byte{0x00, 0x01}, fmt.Errorf("device rejected request: %w", ebuserrors.ErrNACK))

	_, err := disp.Invoke(context.Background(), 0x08, []byte{0x00, 0x01})
	if err == nil {
		t.Fatalf("Invoke under NAK = nil; want UPSTREAM_RPC_FAILED-class error")
	}
	if !errors.Is(err, errRawFrameUpstreamRPCFailed) {
		t.Fatalf("err = %v; want errors.Is(_, errRawFrameUpstreamRPCFailed)", err)
	}
	if errors.Is(err, errRawFrameTransportDown) {
		t.Fatalf("err is misclassified as TRANSPORT_DOWN; NAK is a protocol failure not a transport failure")
	}
}

func TestHistoryDispatcherCRCErrorMapsToUpstreamRPCFailed(t *testing.T) {
	disp, bus, _ := newTestDispatcher(t)
	bus.setErr([2]byte{0x00, 0x01}, errors.New("ebus: crc mismatch expected=ab got=cd"))

	_, err := disp.Invoke(context.Background(), 0x08, []byte{0x00, 0x01})
	if err == nil {
		t.Fatalf("Invoke under CRC = nil; want UPSTREAM_RPC_FAILED-class error")
	}
	if !errors.Is(err, errRawFrameUpstreamRPCFailed) {
		t.Fatalf("err = %v; want errors.Is(_, errRawFrameUpstreamRPCFailed)", err)
	}
}

// --- Misconfiguration / safety guards ---

func TestHistoryDispatcherNilBusReturnsMisconfigured(t *testing.T) {
	mgr := b503session.New(b503session.TransportKey{}, 30*time.Second, nil)
	disp := NewRawFrameDispatcher(nil, adapterSource, nil, mgr, 0)
	_, err := disp.Invoke(context.Background(), 0x08, []byte{0x00, 0x01})
	if err == nil || !errors.Is(err, errRawFrameMisconfigured) {
		t.Fatalf("Invoke with nil bus = %v; want errRawFrameMisconfigured", err)
	}
}

func TestHistoryDispatcherNilManagerReturnsMisconfigured(t *testing.T) {
	bus := newB503DispatcherMockBus()
	disp := NewRawFrameDispatcher(bus, adapterSource, nil, nil, 0)
	_, err := disp.Invoke(context.Background(), 0x08, []byte{0x00, 0x01})
	if err == nil || !errors.Is(err, errRawFrameMisconfigured) {
		t.Fatalf("Invoke with nil mgr = %v; want errRawFrameMisconfigured", err)
	}
}

func TestHistoryDispatcherSourceNotAdmittedDisconnectsOwnerWithoutAdvancingEpoch(t *testing.T) {
	bus := newB503DispatcherMockBus()
	mgr := b503session.New(
		b503session.TransportKey{AdapterInstanceID: "gateway", TransportEpoch: 7},
		30*time.Second,
		nil,
	)
	if _, err := mgr.Enable(context.Background()); err != nil {
		t.Fatalf("Enable() error = %v", err)
	}
	disp := NewRawFrameDispatcherWithSourceProvider(bus, func() (byte, bool) {
		return 0, false
	}, nil, mgr, 0)
	disp.disconnectIfCurrent = (&B503Runtime{manager: mgr}).disconnectIfCurrent
	_, err := disp.Invoke(context.Background(), 0x08, []byte{0x00, 0x01})
	if err == nil || !errors.Is(err, errRawFrameSourceNotAdmitted) {
		t.Fatalf("Invoke before admission = %v; want errRawFrameSourceNotAdmitted", err)
	}
	if bus.callCount() != 0 {
		t.Fatalf("bus.Send calls = %d; want 0 before source admission", bus.callCount())
	}
	if mgr.IsOwned() {
		t.Fatal("source-unavailable read retained the old B503 issuer")
	}
	if got := mgr.TransportKey().TransportEpoch; got != 7 {
		t.Fatalf("transport epoch after source-unavailable read = %d, want 7", got)
	}
	// A later lifecycle withdrawal for the same generation is intentionally
	// idempotent: disconnect releases ownership but never manufactures an epoch.
	mgr.OnTransportDisconnect()
	if got := mgr.TransportKey().TransportEpoch; got != 7 {
		t.Fatalf("transport epoch after duplicate disconnect = %d, want 7", got)
	}
}

func TestHistoryDispatcherCurrentLifecycleWithdrawalReleasesUnrefreshedOwner(t *testing.T) {
	mgr := b503session.New(
		b503session.TransportKey{AdapterInstanceID: "gateway", TransportEpoch: 1},
		30*time.Second,
		func(context.Context) (b503session.TransportKey, error) {
			t.Fatal("withdrawal before an owner operation must not invoke refresh")
			return b503session.TransportKey{}, nil
		},
	)
	runtime := &B503Runtime{manager: mgr}
	if _, err := mgr.EnableOperation(context.Background(), 0x15, func(context.Context, byte) b503session.DispatchOutcome {
		return b503session.DispatchOutcome{
			Emitted: true, Native: b503session.NativeACK,
			Transport: b503session.TransportKey{AdapterInstanceID: "gateway", TransportEpoch: 1},
		}
	}); err != nil {
		t.Fatalf("EnableOperation() error = %v", err)
	}
	before, ok := mgr.CleanupObligation()
	if !ok || before.GatewayCleanupAttemptID == "" || before.TransportEpoch != 1 {
		t.Fatalf("active cleanup = %+v, present=%v", before, ok)
	}

	runtime.EBusDriverActivated(2)
	if got := mgr.TransportKey().TransportEpoch; got != 1 {
		t.Fatalf("activation falsely rebound owner transport epoch = %d, want 1", got)
	}
	if got := mgr.LifecycleTransportKey().TransportEpoch; got != 2 {
		t.Fatalf("lifecycle epoch = %d, want admitted generation 2", got)
	}
	runtime.EBusDriverWithdrawn(1)
	if !mgr.IsOwned() {
		t.Fatal("stale generation-1 withdrawal released generation-2 lifecycle owner")
	}

	runtime.EBusDriverWithdrawn(2)
	if got := mgr.StatusSnapshot(); got.State != b503session.Idle || got.Owned {
		t.Fatalf("generation-2 withdrawal snapshot = %+v, want ownerless public Idle", got)
	}
	after, ok := mgr.CleanupObligation()
	if !ok || after.GatewayCleanupAttemptID != before.GatewayCleanupAttemptID || after.Target != before.Target || after.TransportEpoch != before.TransportEpoch {
		t.Fatalf("withdrawal cleanup drift: before=%+v after=%+v present=%v", before, after, ok)
	}
}

func TestHistoryDispatcherStaleSourceFallbackCannotClearUnsettledCleanup(t *testing.T) {
	bus := newB503DispatcherMockBus()
	mgr := b503session.New(
		b503session.TransportKey{AdapterInstanceID: "gateway", TransportEpoch: 7},
		30*time.Second,
		nil,
	)
	if _, err := mgr.Enable(context.Background()); err != nil {
		t.Fatalf("old generation Enable() error = %v", err)
	}
	runtime := &B503Runtime{manager: mgr}
	var cleanupWrites atomic.Int32
	mgr.SetCleanupDispatcher(func(context.Context, byte) b503session.DispatchOutcome {
		cleanupWrites.Add(1)
		return b503session.DispatchOutcome{Emitted: true, Native: b503session.NativeACK}
	})
	var reenableErr error
	disp := NewRawFrameDispatcherWithSourceProvider(bus, func() (byte, bool) {
		// Invoke already captured epoch 7. Complete the exact lifecycle
		// withdrawal/recovery before returning the stale unavailable result.
		runtime.EBusDriverWithdrawn(7)
		runtime.EBusDriverActivated(8)
		_, reenableErr = mgr.EnableOperation(context.Background(), 0, func(context.Context, byte) b503session.DispatchOutcome {
			t.Fatal("cleanup-fenced re-enable reached dispatcher")
			return b503session.DispatchOutcome{}
		})
		return 0, false
	}, nil, mgr, 0)
	disp.disconnectIfCurrent = runtime.disconnectIfCurrent

	_, err := disp.Invoke(context.Background(), 0x08, []byte{0x00, 0x01})
	if err == nil || !errors.Is(err, errRawFrameSourceNotAdmitted) {
		t.Fatalf("Invoke with stale unavailable result = %v, want errRawFrameSourceNotAdmitted", err)
	}
	if !errors.Is(reenableErr, b503session.ErrCleanupPending) {
		t.Fatalf("epoch-8 re-enable = %v, want ErrCleanupPending", reenableErr)
	}
	if mgr.IsOwned() {
		t.Fatal("epoch-7 source fallback resurrected caller ownership")
	}
	if got := mgr.TransportKey().TransportEpoch; got != 8 {
		t.Fatalf("recovered transport epoch = %d, want 8", got)
	}
	if obligation, ok := mgr.CleanupObligation(); !ok || obligation.TransportEpoch != 7 || obligation.LastNative != b503session.NativeAmbiguous || obligation.LastAttemptEpoch != 0 {
		t.Fatalf("epoch-8 cleanup = %+v, present=%v", obligation, ok)
	}
	if got := cleanupWrites.Load(); got != 0 {
		t.Fatalf("epoch-8 recovery cleanup writes = %d, want zero", got)
	}
	if bus.callCount() != 0 {
		t.Fatalf("bus.Send calls = %d, want 0 for unavailable source", bus.callCount())
	}
}

// --- Source-level invariants: production must not contain the stub literal ---

// TestHistoryDispatcherNoStubLiteralInProductionWiring asserts that the
// stub-error literal "production raw-frame dispatch not yet wired" no
// longer appears anywhere in production wiring code (M6 acceptance §10).
func TestHistoryDispatcherNoStubLiteralInProductionWiring(t *testing.T) {
	// `go test ./cmd/gateway/...` runs with cwd=cmd/gateway, so plain
	// glob "*.go" finds the right files.
	files, err := filepath.Glob("*.go")
	if err != nil || len(files) == 0 {
		t.Fatalf("filepath.Glob(*.go) = %v / %d files; cannot scan production sources", err, len(files))
	}
	bad := "production raw-frame dispatch not yet wired"
	var hits []string
	for _, f := range files {
		if strings.HasSuffix(f, "_test.go") {
			continue
		}
		raw, readErr := os.ReadFile(f)
		if readErr != nil {
			continue
		}
		if strings.Contains(string(raw), bad) {
			hits = append(hits, f)
		}
	}
	if len(hits) > 0 {
		t.Fatalf("M6 acceptance §10: stub-error literal %q still present in production: %v", bad, hits)
	}
}
