package owner

// Frame routing and error mapping of the history dispatcher.

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/owner/b503session"
)

// --- Read selector dispatch tests (5 tools × byte-stream verification) ---

func TestHistoryDispatcherErrorsCurrentRoutesViaBusSend(t *testing.T) {
	disp, bus, _ := newTestDispatcher(t)
	// LOCAL_CAPTURE-shaped 5-slot response with first slot = 0x0119
	// (the worked example from B503.md §5.3).
	bus.setResp([2]byte{0x00, 0x01}, []byte{0x19, 0x01, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF})

	resp, err := disp.Invoke(context.Background(), 0x08, []byte{0x00, 0x01})
	if err != nil {
		t.Fatalf("Invoke errors.current = %v; want nil", err)
	}
	if len(resp) != 10 {
		t.Fatalf("len(resp) = %d; want 10", len(resp))
	}
	if resp[0] != 0x19 || resp[1] != 0x01 {
		t.Fatalf("resp[0:2] = %02x %02x; want 19 01", resp[0], resp[1])
	}
	if bus.callCount() != 1 {
		t.Fatalf("bus call count = %d; want 1", bus.callCount())
	}
	frame, _ := bus.lastCall()
	if frame.Primary != 0xB5 || frame.Secondary != 0x03 {
		t.Fatalf("frame PB/SB = %02x/%02x; want b5/03", frame.Primary, frame.Secondary)
	}
	if frame.Source != adapterSource {
		t.Fatalf("frame.Source = %02x; want %02x", frame.Source, adapterSource)
	}
	if frame.Target != 0x08 {
		t.Fatalf("frame.Target = %02x; want 08", frame.Target)
	}
	if len(frame.Data) != 2 || frame.Data[0] != 0x00 || frame.Data[1] != 0x01 {
		t.Fatalf("frame.Data = %x; want 00 01", frame.Data)
	}
}

func TestHistoryDispatcherErrorsHistoryRoutesViaBusSend(t *testing.T) {
	disp, bus, _ := newTestDispatcher(t)
	// Errorhistory request: family=01, selector=01, plus an index byte.
	bus.setResp([2]byte{0x01, 0x01}, []byte{0x03, 0x77, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF})

	payload := []byte{0x01, 0x01, 0x03}
	resp, err := disp.Invoke(context.Background(), 0x08, payload)
	if err != nil {
		t.Fatalf("Invoke errors.history = %v; want nil", err)
	}
	if len(resp) != 11 {
		t.Fatalf("len(resp) = %d; want 11", len(resp))
	}
	frame, _ := bus.lastCall()
	if len(frame.Data) != 3 {
		t.Fatalf("frame.Data length = %d; want 3 (family+selector+index)", len(frame.Data))
	}
}

func TestHistoryDispatcherServiceCurrentRoutesViaBusSend(t *testing.T) {
	disp, bus, _ := newTestDispatcher(t)
	bus.setResp([2]byte{0x00, 0x02}, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF})

	resp, err := disp.Invoke(context.Background(), 0x08, []byte{0x00, 0x02})
	if err != nil {
		t.Fatalf("Invoke service.current = %v; want nil", err)
	}
	if len(resp) != 10 {
		t.Fatalf("len(resp) = %d; want 10", len(resp))
	}
	frame, _ := bus.lastCall()
	if frame.Data[0] != 0x00 || frame.Data[1] != 0x02 {
		t.Fatalf("frame.Data prefix = %x; want 00 02", frame.Data[:2])
	}
}

func TestHistoryDispatcherServiceHistoryRoutesViaBusSend(t *testing.T) {
	disp, bus, _ := newTestDispatcher(t)
	bus.setResp([2]byte{0x01, 0x02}, []byte{0x05, 0x42, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF})

	payload := []byte{0x01, 0x02, 0x05}
	resp, err := disp.Invoke(context.Background(), 0x08, payload)
	if err != nil {
		t.Fatalf("Invoke service.history = %v; want nil", err)
	}
	if len(resp) != 11 {
		t.Fatalf("len(resp) = %d; want 11", len(resp))
	}
}

func TestHistoryDispatcherLiveMonitorRoutesViaBusSend(t *testing.T) {
	disp, bus, _ := newTestDispatcher(t)
	// LiveMonitorMain request: family=00, selector=03. Response carries
	// status + function in the first two bytes.
	bus.setResp([2]byte{0x00, 0x03}, []byte{0x01, 0x42, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	resp, err := disp.Invoke(context.Background(), 0x08, []byte{0x00, 0x03})
	if err != nil {
		t.Fatalf("Invoke live_monitor = %v; want nil", err)
	}
	if len(resp) != 8 {
		t.Fatalf("len(resp) = %d; want 8", len(resp))
	}
	if resp[0] != 0x01 || resp[1] != 0x42 {
		t.Fatalf("resp[0:2] = %02x %02x; want 01 42", resp[0], resp[1])
	}
}

// --- Reject malformed payload (§12.2 explicit reject) ---

func TestHistoryDispatcherRejectsPayloadStartingWithNamespaceBytes(t *testing.T) {
	disp, bus, _ := newTestDispatcher(t)
	_, err := disp.Invoke(context.Background(), 0x08, []byte{0xB5, 0x03, 0x00, 0x01})
	if err == nil {
		t.Fatalf("Invoke with payload starting b5 03 = nil; want errRawFrameMalformedPayload")
	}
	if !errors.Is(err, errRawFrameMalformedPayload) {
		t.Fatalf("err = %v; want errors.Is(_, errRawFrameMalformedPayload)", err)
	}
	if bus.callCount() != 0 {
		t.Fatalf("bus.Send was called %d times; want 0 (reject must happen before wire emission)", bus.callCount())
	}
}

// --- Error mapping table (§12.4) ---

func TestHistoryDispatcherTransportDownFiresOnTransportDisconnect(t *testing.T) {
	disp, bus, mgr := newTestDispatcher(t)
	if _, err := mgr.Enable(context.Background()); err != nil {
		t.Fatalf("mgr.Enable = %v", err)
	}
	bus.setErr([2]byte{0x00, 0x01}, errors.New("ebus: transport closed"))

	_, err := disp.Invoke(context.Background(), 0x08, []byte{0x00, 0x01})
	if err == nil {
		t.Fatalf("Invoke = nil; want TRANSPORT_DOWN-class error")
	}
	if !errors.Is(err, errRawFrameTransportDown) {
		t.Fatalf("err = %v; want errors.Is(_, errRawFrameTransportDown)", err)
	}
	if mgr.IsOwned() {
		t.Fatalf("Manager.IsOwned() = true after transport-down; want false (OnTransportDisconnect should fire)")
	}
}

func TestHistoryDispatcherStaleSendErrorCannotClearUnsettledCleanup(t *testing.T) {
	disp, bus, mgr := newTestDispatcher(t)
	runtime := &B503Runtime{manager: mgr}
	var cleanupWrites atomic.Int32
	mgr.SetCleanupDispatcher(func(context.Context, byte) b503session.DispatchOutcome {
		cleanupWrites.Add(1)
		return b503session.DispatchOutcome{Emitted: true, Native: b503session.NativeACK}
	})
	disp.disconnectIfCurrent = runtime.disconnectIfCurrent
	if _, err := mgr.Enable(context.Background()); err != nil {
		t.Fatalf("old generation Enable() error = %v", err)
	}

	classification := make(chan struct{})
	continueClassify := make(chan struct{})
	defer func() {
		select {
		case <-continueClassify:
		default:
			close(continueClassify)
		}
	}()
	bus.setErr([2]byte{0x00, 0x01}, &classifyBarrierError{
		classification:   classification,
		continueClassify: continueClassify,
	})
	invokeDone := make(chan error, 1)
	go func() {
		_, err := disp.Invoke(context.Background(), 0x08, []byte{0x00, 0x01})
		invokeDone <- err
	}()

	select {
	case <-classification:
		// The request from generation 1 passed its end-epoch check. Retire it
		// and install a new issuer before the stale send error is classified.
	case <-time.After(time.Second):
		t.Fatal("Invoke did not reach the post-epoch-check classification barrier")
	}
	runtime.EBusDriverWithdrawn(1)
	before, ok := mgr.CleanupObligation()
	if !ok {
		t.Fatal("generation-1 withdrawal did not retain cleanup")
	}
	runtime.EBusDriverActivated(2)
	afterRecovery, ok := mgr.CleanupObligation()
	if !ok || afterRecovery != before {
		t.Fatalf("reconnect mutated cleanup: before=%+v after=%+v present=%v", before, afterRecovery, ok)
	}
	if got := cleanupWrites.Load(); got != 0 {
		t.Fatalf("reconnect cleanup writes = %d, want zero", got)
	}
	var reenableCalls int
	_, err := mgr.EnableOperation(context.Background(), 0, func(context.Context, byte) b503session.DispatchOutcome {
		reenableCalls++
		return b503session.DispatchOutcome{Emitted: true, Native: b503session.NativeACK}
	})
	if !errors.Is(err, b503session.ErrCleanupPending) || reenableCalls != 0 {
		t.Fatalf("re-enable = %v, dispatches=%d; want cleanup pending without write", err, reenableCalls)
	}
	close(continueClassify)

	select {
	case err := <-invokeDone:
		if err == nil || !errors.Is(err, errRawFrameTransportDown) {
			t.Fatalf("stale Invoke error = %v, want errRawFrameTransportDown", err)
		}
	case <-time.After(time.Second):
		t.Fatal("stale Invoke did not finish classification")
	}
	if mgr.IsOwned() {
		t.Fatal("generation-1 send error resurrected caller ownership")
	}
	if got := mgr.TransportKey().TransportEpoch; got != 2 {
		t.Fatalf("recovered transport epoch = %d, want 2", got)
	}
	afterStale, ok := mgr.CleanupObligation()
	if !ok || afterStale.GatewayCleanupAttemptID != before.GatewayCleanupAttemptID {
		t.Fatalf("stale error mutated cleanup = %+v, present=%v", afterStale, ok)
	}
}

func TestHistoryDispatcherCtxCanceledMapsToUpstreamTimeout(t *testing.T) {
	disp, bus, _ := newTestDispatcher(t)
	gate := make(chan struct{})
	bus.blockUntil = gate
	defer close(gate)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(20 * time.Millisecond)
		cancel()
	}()
	_, err := disp.Invoke(ctx, 0x08, []byte{0x00, 0x01})
	if err == nil {
		t.Fatalf("Invoke under canceled ctx = nil; want UPSTREAM_TIMEOUT-class error")
	}
	if !errors.Is(err, errRawFrameUpstreamTimeout) {
		t.Fatalf("err = %v; want errors.Is(_, errRawFrameUpstreamTimeout)", err)
	}
	if errors.Is(err, b503session.ErrTransportDown) {
		t.Fatalf("ordinary caller cancellation err = %v; must not map TRANSPORT_DOWN", err)
	}
}

func TestHistoryDispatcherLifecycleCanceledContextMapsTransportDown(t *testing.T) {
	disp, bus, _ := newTestDispatcher(t)
	ctx, cancel := context.WithCancelCause(context.Background())
	cancel(b503session.ErrTransportDown)
	outcome := disp.InvokeB503Outcome(ctx, 0x15, []byte{0x00, 0x03})
	if outcome.Err == nil || !errors.Is(outcome.Err, b503session.ErrTransportDown) {
		t.Fatalf("lifecycle-canceled outcome error = %v, want TRANSPORT_DOWN", outcome.Err)
	}
	if errors.Is(outcome.Err, errRawFrameUpstreamTimeout) {
		t.Fatalf("lifecycle-canceled outcome error = %v, must not map upstream timeout", outcome.Err)
	}
	if outcome.Emitted || bus.callCount() != 0 {
		t.Fatalf("lifecycle-canceled outcome = %+v, sends=%d; want no emission", outcome, bus.callCount())
	}
}
