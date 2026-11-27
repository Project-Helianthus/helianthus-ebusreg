package owner

// Outcome, emission fence and lifecycle behaviour of the history dispatcher.

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"

	ebuserrors "github.com/Project-Helianthus/helianthus-ebusgo/errors"
	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/owner/b503session"
)

func TestHistoryDispatcherOutcomePreCanceledIsNotEmitted(t *testing.T) {
	disp, bus, _ := newTestDispatcher(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	outcome := disp.InvokeB503Outcome(ctx, 0x15, []byte{0x00, 0x03})
	if outcome.Emitted {
		t.Fatal("pre-canceled outcome marked emitted")
	}
	if outcome.Err == nil || !errors.Is(outcome.Err, errRawFrameUpstreamTimeout) {
		t.Fatalf("pre-canceled outcome error = %v", outcome.Err)
	}
	if calls := bus.callCount(); calls != 0 {
		t.Fatalf("pre-canceled bus.Send count = %d, want 0", calls)
	}
}

func TestHistoryDispatcherOutcomeNAKIsEmittedAndNativeNAK(t *testing.T) {
	disp, bus, _ := newTestDispatcher(t)
	bus.setErr([2]byte{0x00, 0x03}, fmt.Errorf("target rejected request: %w", ebuserrors.ErrNACK))
	outcome := disp.InvokeB503Outcome(context.Background(), 0x15, []byte{0x00, 0x03})
	if !outcome.Emitted || outcome.Native != b503session.NativeNAK {
		t.Fatalf("NAK outcome = %+v, want emitted NativeNAK", outcome)
	}
	if outcome.Err == nil || !errors.Is(outcome.Err, errRawFrameUpstreamRPCFailed) {
		t.Fatalf("NAK outcome error = %v", outcome.Err)
	}
	if calls := bus.callCount(); calls != 1 {
		t.Fatalf("NAK bus.Send count = %d, want 1", calls)
	}
}

func TestHistoryDispatcherOutcomeNAKTextWithoutSentinelRemainsAmbiguous(t *testing.T) {
	d, bus, _ := newTestDispatcher(t)
	bus.setErr([2]byte{0x00, 0x03}, errors.New("proxy diagnostic contains NAK text only"))
	outcome := d.InvokeB503Outcome(context.Background(), 0x15, []byte{0x00, 0x03})
	if !outcome.Emitted || outcome.Native != b503session.NativeAmbiguous {
		t.Fatalf("text-only NAK outcome = %+v, want emitted NativeAmbiguous", outcome)
	}
}

func TestHistoryDispatcherOutcomeACKIsNativeEvidenceOnly(t *testing.T) {
	disp, bus, _ := newTestDispatcher(t)
	bus.setResp([2]byte{0x00, 0x03}, []byte{0x01, 0x00})
	outcome := disp.InvokeB503Outcome(context.Background(), 0x15, []byte{0x00, 0x03})
	if outcome.Err != nil || !outcome.Emitted || outcome.Native != b503session.NativeACK {
		t.Fatalf("ACK outcome = %+v, want exact emitted NativeACK", outcome)
	}
	if calls := bus.callCount(); calls != 1 {
		t.Fatalf("ACK bus.Send count = %d, want 1", calls)
	}
}

func TestHistoryDispatcherEnableDisconnectBeforeEmissionReturnsIdleWithoutCleanup(t *testing.T) {
	bus := newB503DispatcherMockBus()
	bus.setResp([2]byte{0x00, 0x03}, []byte{0x01, 0x00})
	mgr := b503session.New(
		b503session.TransportKey{AdapterInstanceID: "test", TransportEpoch: 1},
		30*time.Second,
		nil,
	)
	readMu := newB503BlockingReadLocker()
	dispatcher := NewRawFrameDispatcher(bus, adapterSource, readMu, mgr, time.Second)

	done := make(chan error, 1)
	go func() {
		_, err := mgr.EnableOperation(context.Background(), 0x15, func(ctx context.Context, target byte) b503session.DispatchOutcome {
			return InvokeB503Operation(ctx, dispatcher, target, []byte{0x00, 0x03})
		})
		done <- err
	}()
	<-readMu.entered
	if pending, ok := mgr.PendingOperation(); !ok || pending.Kind != b503session.OperationEnable || pending.Emitted {
		t.Fatalf("pre-emission pending enable = %+v, present=%v", pending, ok)
	}
	if calls := bus.callCount(); calls != 0 {
		t.Fatalf("pre-emission bus.Send count = %d, want 0", calls)
	}

	mgr.OnTransportDisconnect()
	if got := mgr.StatusSnapshot(); got.State != b503session.Idle || got.Owned {
		t.Fatalf("pre-emission disconnect snapshot = %+v, want ownerless Idle", got)
	}
	if _, ok := mgr.PendingOperation(); ok {
		t.Fatal("pre-emission disconnect retained pending enable")
	}
	if cleanup, ok := mgr.CleanupObligation(); ok {
		t.Fatalf("pre-emission disconnect cleanup = %+v, want none", cleanup)
	}
	close(readMu.release)
	if err := <-done; !errors.Is(err, b503session.ErrTransportDown) {
		t.Fatalf("pre-emission disconnect error = %v, want exact TRANSPORT_DOWN", err)
	}
	if calls := bus.callCount(); calls != 0 {
		t.Fatalf("disconnected queued enable bus.Send count = %d, want 0", calls)
	}

	mgr.OnEpochAdvance(context.Background(), 2)
	if cleanup, ok := mgr.CleanupObligation(); ok {
		t.Fatalf("reconnect manufactured cleanup = %+v", cleanup)
	}
	if _, err := mgr.EnableOperation(context.Background(), 0x15, func(context.Context, byte) b503session.DispatchOutcome {
		return b503session.DispatchOutcome{Emitted: true, Native: b503session.NativeACK, Transport: mgr.TransportKey()}
	}); err != nil {
		t.Fatalf("new explicit enable after reconnect = %v", err)
	}
}

func TestHistoryDispatcherEpochAdvanceAfterAdmissionDoesNotWriteReplacementGeneration(t *testing.T) {
	bus := newB503DispatcherMockBus()
	bus.setResp([2]byte{0x00, 0x03}, []byte{0x01, 0x00})
	mgr := b503session.New(
		b503session.TransportKey{AdapterInstanceID: "test", TransportEpoch: 1},
		30*time.Second,
		nil,
	)
	dispatcher := NewRawFrameDispatcher(bus, adapterSource, &sync.Mutex{}, mgr, time.Second)
	epochAttempted := make(chan struct{})
	epochAdvanced := make(chan struct{})
	dispatcher.afterPendingEmission = func() {
		go func() {
			close(epochAttempted)
			mgr.OnEpochAdvance(context.Background(), 2)
			close(epochAdvanced)
		}()
		<-epochAttempted
	}
	bus.onSend = func(frame protocol.Frame, _ *b503DispatcherMockBus) {
		if len(frame.Data) < 2 || frame.Data[0] != 0x00 || frame.Data[1] != 0x03 {
			t.Fatalf("wire frame = %x, want SERVICE_WRITE selector 00 03", frame.Data)
		}
		select {
		case <-epochAdvanced:
			t.Fatal("epoch advance completed before admitted SERVICE_WRITE entered bus.Send")
		default:
		}
		if got := mgr.DispatchEpoch(); got != 1 {
			t.Fatalf("SERVICE_WRITE entered replacement epoch %d, want issuing epoch 1", got)
		}
	}

	_, err := mgr.EnableOperation(context.Background(), 0x15, func(ctx context.Context, target byte) b503session.DispatchOutcome {
		return InvokeB503Operation(ctx, dispatcher, target, []byte{0x00, 0x03})
	})
	if err != nil && !errors.Is(err, errRawFrameStaleEpoch) && !errors.Is(err, b503session.ErrTransportDown) {
		t.Fatalf("admitted enable outcome = %v, want success or conservative replacement-generation result", err)
	}
	select {
	case <-epochAdvanced:
	case <-time.After(time.Second):
		t.Fatal("epoch advance remained blocked after bus.Send returned")
	}
	if got := mgr.DispatchEpoch(); got != 2 {
		t.Fatalf("epoch after send = %d, want replacement epoch 2", got)
	}
	if calls := bus.callCount(); calls != 1 {
		t.Fatalf("SERVICE_WRITE bus.Send count = %d, want exactly 1 issuing-generation send", calls)
	}
}

func TestHistoryDispatcherDisableReadMuDeadlineBeforeEmissionPreservesOwner(t *testing.T) {
	bus := newB503DispatcherMockBus()
	bus.setResp([2]byte{0x00, 0x03}, []byte{0x01, 0x00})
	mgr := b503session.New(
		b503session.TransportKey{AdapterInstanceID: "test", TransportEpoch: 1},
		0, // Idle expiry is outside this pre-emission cancellation test.
		nil,
	)
	key, err := mgr.EnableOperation(context.Background(), 0x15, func(context.Context, byte) b503session.DispatchOutcome {
		return b503session.DispatchOutcome{Emitted: true, Native: b503session.NativeACK, Transport: mgr.TransportKey()}
	})
	if err != nil {
		t.Fatalf("setup enable: %v", err)
	}
	before, ok := mgr.CleanupObligation()
	if !ok || before.LastAttempted {
		t.Fatalf("setup cleanup evidence = %+v, present=%v", before, ok)
	}
	readMu := newB503BlockingReadLocker()
	dispatcher := NewRawFrameDispatcher(bus, adapterSource, readMu, mgr, 0)
	// Arm expiry only after the dispatcher requests readMu. A wall-clock
	// deadline started before goroutine admission can expire before Lock and
	// leave the test waiting forever for readMu.entered under runner load.
	ctx := newB503LaggingParentContext()
	var expireOnce, releaseOnce sync.Once
	expire := func() {
		expireOnce.Do(func() {
			ctx.expire()
			// Deliver the held deadline callback synchronously so the Manager's
			// derived operation context is expired before readMu is released.
			ctx.mu.Lock()
			f := ctx.afterFunc
			ctx.afterFunc = nil
			ctx.mu.Unlock()
			if f != nil {
				f()
			}
		})
	}
	release := func() { releaseOnce.Do(func() { close(readMu.release) }) }
	operationCtx := make(chan context.Context, 1)
	done := make(chan error, 1)
	joined := make(chan struct{})
	t.Cleanup(func() {
		expire()
		release()
		select {
		case <-joined:
		case <-hangGuard(t):
			t.Error("Disable worker did not join after cancellation and readMu release")
		}
	})
	go func() {
		defer close(joined)
		done <- mgr.DisableOperation(ctx, key, 0x15, func(dispatchCtx context.Context, target byte) b503session.DispatchOutcome {
			operationCtx <- dispatchCtx
			return InvokeB503Operation(dispatchCtx, dispatcher, target, []byte{0x00, 0x03})
		})
	}()
	select {
	case <-readMu.entered:
	case <-joined:
		t.Fatalf("Disable returned before requesting readMu: %v", <-done)
	case <-hangGuard(t):
		t.Fatal("Disable did not request readMu")
	}
	if pending, ok := mgr.PendingOperation(); !ok || pending.Kind != b503session.OperationDisable || pending.Emitted {
		t.Fatalf("blocked disable = %+v, present=%v; want pending without emission", pending, ok)
	}
	if calls := bus.callCount(); calls != 0 {
		t.Fatalf("blocked disable bus.Send count = %d, want 0 before expiry", calls)
	}
	if got := mgr.StatusSnapshot(); got.State != b503session.Active || !got.Owned {
		t.Fatalf("blocked disable state = %+v, want owned Active", got)
	}
	dispatchCtx := <-operationCtx // Published before the readMu entry barrier.
	expire()
	select {
	case <-dispatchCtx.Done():
	default:
		t.Fatal("derived operation cancellation was not delivered before readMu release")
	}
	if err := dispatchCtx.Err(); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("derived operation error = %v, want deadline exceeded", err)
	}
	release()
	select {
	case <-joined:
	case <-hangGuard(t):
		t.Fatal("expired Disable did not complete after readMu release")
	}
	if err := <-done; !errors.Is(err, errRawFrameUpstreamTimeout) {
		t.Errorf("pre-emission readMu deadline error = %v, want upstream timeout", err)
	}
	if calls := bus.callCount(); calls != 0 {
		t.Errorf("readMu-deadline bus.Send count = %d, want 0", calls)
	}
	if got := mgr.StatusSnapshot(); got.State != b503session.Active || !got.Owned {
		t.Errorf("pre-emission readMu deadline state = %+v, want owned Active", got)
	}
	after, ok := mgr.CleanupObligation()
	if !ok || !reflect.DeepEqual(after, before) {
		t.Errorf("pre-emission readMu deadline recorded a nonexistent cleanup attempt: before=%+v after=%+v present=%v", before, after, ok)
	}
	retryCtx, cancelRetry := context.WithCancel(context.Background())
	t.Cleanup(cancelRetry)
	done = make(chan error, 1)
	joined = make(chan struct{})
	go func() {
		defer close(joined)
		done <- mgr.DisableOperation(retryCtx, key, 0x15, func(dispatchCtx context.Context, target byte) b503session.DispatchOutcome {
			return InvokeB503Operation(dispatchCtx, dispatcher, target, []byte{0x00, 0x03})
		})
	}()
	select {
	case <-joined:
	case <-hangGuard(t):
		t.Fatal("nonexpired retry Disable did not join")
	}
	if err := <-done; err != nil {
		t.Fatalf("nonexpired retry disable: %v", err)
	}
	if calls := bus.callCount(); calls != 1 {
		t.Fatalf("nonexpired retry bus.Send count = %d, want 1", calls)
	}
	retried, ok := mgr.CleanupObligation()
	if !ok || !retried.LastAttempted || !retried.LastEmitted || retried.LastNative != b503session.NativeACK || retried.LastErr != nil {
		t.Fatalf("nonexpired retry cleanup = %+v, present=%v; want emitted ACK", retried, ok)
	}
	if got := mgr.StatusSnapshot(); got.Owned || got.State != b503session.Idle {
		t.Fatalf("nonexpired retry state = %+v, want ownerless Idle", got)
	}
}

func TestHistoryDispatcherDispatcherRejectsLaggingParentDeadlineBeforeEmission(t *testing.T) {
	bus := newB503DispatcherMockBus()
	bus.setResp([2]byte{0x00, 0x03}, []byte{0x01, 0x00})
	mgr := b503session.New(
		b503session.TransportKey{AdapterInstanceID: "test", TransportEpoch: 1},
		30*time.Second,
		nil,
	)
	readMu := newB503BlockingReadLocker()
	dispatcher := NewRawFrameDispatcher(bus, adapterSource, readMu, mgr, time.Second)
	ctx := newB503LaggingParentContext()
	done := make(chan B503DispatchOutcome, 1)
	go func() {
		done <- dispatcher.InvokeB503Outcome(ctx, 0x15, []byte{0x00, 0x03})
	}()
	<-readMu.entered
	ctx.expire()
	close(readMu.release)
	if outcome := <-done; outcome.Err == nil || outcome.Emitted {
		t.Fatalf("lagging parent deadline outcome = %+v, want a pre-emission deadline failure", outcome)
	}
	if calls := bus.callCount(); calls != 0 {
		t.Fatalf("lagging-parent-deadline bus.Send count = %d, want 0", calls)
	}
}

func TestHistoryDispatcherDisableDisconnectBeforeEmissionRetainsFenceWithoutWrite(t *testing.T) {
	for _, tc := range []struct {
		name         string
		refreshEpoch uint64
	}{
		{name: "active"},
		{name: "refreshed", refreshEpoch: 2},
	} {
		t.Run(tc.name, func(t *testing.T) {
			bus := newB503DispatcherMockBus()
			var refresh b503session.RefreshFunc
			if tc.refreshEpoch != 0 {
				refresh = func(context.Context) (b503session.TransportKey, error) {
					return b503session.TransportKey{AdapterInstanceID: "test", TransportEpoch: tc.refreshEpoch}, nil
				}
			}
			mgr := b503session.New(
				b503session.TransportKey{AdapterInstanceID: "test", TransportEpoch: 1},
				100*time.Millisecond,
				refresh,
			)
			key, err := mgr.EnableOperation(context.Background(), 0x15, func(context.Context, byte) b503session.DispatchOutcome {
				return b503session.DispatchOutcome{Emitted: true, Native: b503session.NativeACK, Transport: mgr.TransportKey()}
			})
			if err != nil {
				t.Fatalf("setup enable = %v", err)
			}
			if tc.refreshEpoch != 0 {
				mgr.OnEpochAdvance(context.Background(), tc.refreshEpoch)
			}
			readMu := newB503BlockingReadLocker()
			dispatcher := NewRawFrameDispatcher(bus, adapterSource, readMu, mgr, time.Second)
			done := make(chan error, 1)
			go func() {
				done <- mgr.DisableOperation(context.Background(), key, 0x15, func(ctx context.Context, target byte) b503session.DispatchOutcome {
					return InvokeB503Operation(ctx, dispatcher, target, []byte{0x00, 0x03})
				})
			}()
			<-readMu.entered
			// The admitted explicit Disable owns the terminal path while it
			// waits for quiescence; its former idle deadline must not race it.
			time.Sleep(150 * time.Millisecond)
			pending, ok := mgr.PendingOperation()
			if !ok || pending.Kind != b503session.OperationDisable || pending.Emitted {
				t.Fatalf("pre-emission pending disable = %+v, present=%v", pending, ok)
			}
			if got := mgr.StatusSnapshot(); got.State != b503session.Active || !got.Owned {
				t.Fatalf("disable waiting past idle deadline = %+v, want owned Active", got)
			}
			before, ok := mgr.CleanupObligation()
			if !ok || before.Target != 0x15 || before.GatewayCleanupAttemptID == "" {
				t.Fatalf("pre-disconnect cleanup = %+v, present=%v", before, ok)
			}
			if tc.refreshEpoch != 0 && before.TransportEpoch != tc.refreshEpoch {
				t.Fatalf("refreshed cleanup epoch = %d, want %d", before.TransportEpoch, tc.refreshEpoch)
			}

			mgr.OnTransportDisconnect()
			close(readMu.release)
			if err := <-done; !errors.Is(err, b503session.ErrTransportDown) {
				t.Fatalf("pre-emission disable error = %v, want exact TRANSPORT_DOWN", err)
			}
			if calls := bus.callCount(); calls != 0 {
				t.Fatalf("disconnected queued disable bus.Send count = %d, want 0", calls)
			}
			if got := mgr.StatusSnapshot(); got.State != b503session.Idle || got.Owned {
				t.Fatalf("pre-emission disable snapshot = %+v, want ownerless Idle", got)
			}
			after, ok := mgr.CleanupObligation()
			if !ok || after.GatewayCleanupAttemptID != before.GatewayCleanupAttemptID || after.Target != before.Target || after.TransportEpoch != before.TransportEpoch || after.LastEmitted || after.LastNative != b503session.NativeAmbiguous || !errors.Is(after.LastErr, b503session.ErrTransportDown) {
				t.Fatalf("retained cleanup mismatch: before=%+v after=%+v present=%v", before, after, ok)
			}
			mgr.OnEpochAdvance(context.Background(), mgr.TransportKey().TransportEpoch+1)
			mgr.OnEpochAdvance(context.Background(), mgr.TransportKey().TransportEpoch+1)
			if calls := bus.callCount(); calls != 0 {
				t.Fatalf("later epochs emitted recovery writes = %d, want 0", calls)
			}
			var reenableWrites int
			if _, err := mgr.EnableOperation(context.Background(), 0x15, func(context.Context, byte) b503session.DispatchOutcome {
				reenableWrites++
				return b503session.DispatchOutcome{Emitted: true, Native: b503session.NativeACK}
			}); !errors.Is(err, b503session.ErrCleanupPending) || reenableWrites != 0 {
				t.Fatalf("cleanup-fenced re-enable = %v, writes=%d", err, reenableWrites)
			}
		})
	}
}

func TestHistoryDispatcherEnableDisconnectAfterEmissionRetainsUnknownFence(t *testing.T) {
	bus := newB503DispatcherMockBus()
	bus.setErr([2]byte{0x00, 0x03}, errors.New("ebus: transport closed"))
	enteredSend := make(chan struct{})
	releaseSend := make(chan struct{})
	bus.onSend = func(protocol.Frame, *b503DispatcherMockBus) { close(enteredSend) }
	bus.blockUntil = releaseSend
	mgr := b503session.New(
		b503session.TransportKey{AdapterInstanceID: "test", TransportEpoch: 1},
		30*time.Second,
		nil,
	)
	dispatcher := NewRawFrameDispatcher(bus, adapterSource, &sync.Mutex{}, mgr, time.Second)
	done := make(chan error, 1)
	go func() {
		_, err := mgr.EnableOperation(context.Background(), 0x15, func(ctx context.Context, target byte) b503session.DispatchOutcome {
			return InvokeB503Operation(ctx, dispatcher, target, []byte{0x00, 0x03})
		})
		done <- err
	}()
	<-enteredSend
	if pending, ok := mgr.PendingOperation(); !ok || !pending.Emitted {
		t.Fatalf("in-flight emitted enable = %+v, present=%v", pending, ok)
	}
	mgr.OnTransportDisconnect()
	close(releaseSend)
	if err := <-done; !errors.Is(err, b503session.ErrTransportDown) {
		t.Fatalf("post-emission disconnect error = %v, want TRANSPORT_DOWN", err)
	}
	if calls := bus.callCount(); calls != 1 {
		t.Fatalf("post-emission enable bus.Send count = %d, want 1", calls)
	}
	cleanup, ok := mgr.CleanupObligation()
	if !ok || cleanup.Target != 0x15 || !cleanup.OriginEmitted || cleanup.OriginNative != b503session.NativeAmbiguous || cleanup.GatewayCleanupAttemptID == "" {
		t.Fatalf("post-emission cleanup = %+v, present=%v", cleanup, ok)
	}
	if got := mgr.StatusSnapshot(); got.State != b503session.Idle || got.Owned {
		t.Fatalf("post-emission snapshot = %+v, want ownerless public Idle", got)
	}
	mgr.OnEpochAdvance(context.Background(), 2)
	mgr.OnEpochAdvance(context.Background(), 3)
	if calls := bus.callCount(); calls != 1 {
		t.Fatalf("automatic reconnect writes raised count to %d, want 1", calls)
	}
}
