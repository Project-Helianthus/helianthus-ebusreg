package b503session

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"sync/atomic"
	"time"
)

// armIdleTimerLocked (re)starts the idle-timeout timer. Safe to call
// repeatedly; any previously-running timer is stopped first. Each call
// increments idleTimerGen so an in-flight stale callback (stopped timer
// whose goroutine was already running at Stop() time) can detect its
// generation is no longer current and exit without touching the FSM.
func (m *Manager) armIdleTimerLocked() {
	if m.idleTimer != nil {
		m.idleTimer.Stop()
	}
	if m.idleTimeout <= 0 {
		return
	}
	m.idleTimerGen++
	gen := m.idleTimerGen
	m.idleTimer = time.AfterFunc(m.idleTimeout, func() { m.idleTimerFired(gen) })
}

// idleTimerFired is the timer callback. It acquires stateMu itself; it
// MUST NOT be called with stateMu held. If `gen` does not match the
// current idleTimerGen, this callback belongs to a previously-stopped
// timer whose goroutine was already running when Stop() was called — it
// MUST return without side effects to avoid prematurely disabling a
// rearmed active session.
func (m *Manager) idleTimerFired(gen uint64) {
	m.stateMu.Lock()
	if gen != m.idleTimerGen {
		m.stateMu.Unlock()
		return // stale callback from a stopped-then-rearmed timer
	}
	if m.state != Active {
		m.stateMu.Unlock()
		return
	}
	if pending := m.pendingOperation; pending != nil {
		// An admitted read owns the current Active session until it reports its
		// native outcome. Record this elapsed timer and let ReadOperation either
		// re-arm on ACK or perform one conservative cleanup on failure.
		if pending.Kind == OperationRead {
			m.idleExpiryPending = true
		}
		m.stateMu.Unlock()
		return
	}
	target := m.activeTarget
	if m.pendingEpoch > m.transport.TransportEpoch {
		attemptedEpoch := m.pendingEpoch
		pending := &pendingOperation{PendingOperationSnapshot: PendingOperationSnapshot{
			OperationID: newCleanupAttemptID(), Kind: OperationDisable, Target: target, Transport: m.transport,
		}}
		m.pendingOperation = pending
		m.state = Refreshing
		refresh := m.refresh
		m.stateMu.Unlock()
		m.idleRefreshAndDisable(pending, attemptedEpoch, refresh)
		return
	}
	epoch := m.transport.TransportEpoch
	attemptID := newCleanupAttemptID()
	m.toDisabledLocked()
	m.createCleanupLocked(target, epoch, attemptID)
	m.stateMu.Unlock()
	m.attemptCleanup(context.Background(), epoch)
}

// idleRefreshAndDisable gives an expired owner the same one-shot rebind fence
// as an admitted Disable before emitting the timer's defensive disable.
func (m *Manager) idleRefreshAndDisable(pending *pendingOperation, attemptedEpoch uint64, refresh RefreshFunc) {
	var (
		transport TransportKey
		err       error
	)
	if refresh == nil {
		err = ErrTransportDown
	} else {
		transport, err = refresh(context.Background())
	}
	m.stateMu.Lock()
	if m.pendingOperation != pending || !m.mutexHeld || m.state != Refreshing {
		m.stateMu.Unlock()
		return
	}
	if pending.invalidated || err != nil || transport.TransportEpoch < attemptedEpoch || transport.TransportEpoch < m.observedEpoch {
		m.pendingOperation = nil
		latest := m.observedEpoch
		m.pendingEpoch = 0
		m.transport.TransportEpoch = latest
		m.toDisabledLocked()
		m.createCleanupLocked(pending.Target, latest, pending.OperationID)
		m.stateMu.Unlock()
		return
	}
	m.transport = transport
	m.observedEpoch = transport.TransportEpoch
	m.pendingEpoch = 0
	m.state = Active
	pending.Transport = transport
	dispatch := m.cleanupDispatch
	m.stateMu.Unlock()
	if dispatch == nil {
		m.stateMu.Lock()
		if m.pendingOperation == pending {
			m.pendingOperation = nil
			m.toDisabledLocked()
			m.createCleanupLocked(pending.Target, transport.TransportEpoch, pending.OperationID)
		}
		m.stateMu.Unlock()
		return
	}
	ctx, cancel := m.bindPendingContext(context.Background(), pending)
	defer cancel(nil)
	outcome := dispatchOnce(withPendingOperation(ctx, pending.OperationID), pending.Target, dispatch)
	m.stateMu.Lock()
	if m.pendingOperation != pending {
		m.stateMu.Unlock()
		return
	}
	if pending.invalidated || m.observedEpoch > transport.TransportEpoch {
		m.pendingOperation = nil
		m.pendingEpoch = 0
		m.transport.TransportEpoch = m.observedEpoch
		m.toDisabledLocked()
		m.createCleanupLocked(pending.Target, m.observedEpoch, pending.OperationID)
		if m.cleanup != nil {
			m.cleanup.TransportEpoch = m.observedEpoch
		}
		m.stateMu.Unlock()
		return
	}
	m.pendingOperation = nil
	m.releaseOwnerLocked()
	m.state = Disabled
	m.createCleanupLocked(pending.Target, transport.TransportEpoch, pending.OperationID)
	m.cleanup.LastAttempted = true
	m.cleanup.LastAttemptEpoch = transport.TransportEpoch
	m.cleanup.LastEmitted = outcome.Emitted
	m.cleanup.LastNative = outcome.Native
	m.cleanup.LastErr = outcome.Err
	m.stateMu.Unlock()
}

func (m *Manager) beginOwnerOperation(ctx context.Context, target byte, key *SessionKey, kind OperationKind) (*pendingOperation, error) {
	m.stateMu.Lock()
	if m.state == Refreshing {
		m.stateMu.Unlock()
		return nil, ErrSessionBusy
	}
	if !m.mutexHeld || m.state != Active {
		m.stateMu.Unlock()
		return nil, ErrNotActive
	}
	if m.pendingOperation != nil {
		m.stateMu.Unlock()
		return nil, ErrSessionBusy
	}
	if m.activeTarget != target {
		m.stateMu.Unlock()
		return nil, ErrTargetMismatch
	}
	if key != nil && (key.Transport != m.transport || key.IssuerToken != m.activeToken) {
		m.stateMu.Unlock()
		return nil, ErrWrongToken
	}
	pending := &pendingOperation{PendingOperationSnapshot: PendingOperationSnapshot{
		OperationID: newCleanupAttemptID(), Kind: kind, Target: target, Transport: m.transport,
	}}
	m.pendingOperation = pending
	if m.pendingEpoch <= m.transport.TransportEpoch {
		m.stateMu.Unlock()
		return pending, nil
	}
	attemptedEpoch := m.pendingEpoch
	m.state = Refreshing
	refresh := m.refresh
	m.stateMu.Unlock()

	var (
		newTransport TransportKey
		err          error
	)
	if refresh == nil {
		err = ErrTransportDown
	} else {
		newTransport, err = refresh(ctx)
	}

	m.stateMu.Lock()
	if m.pendingOperation != pending || !m.mutexHeld || m.state != Refreshing {
		m.stateMu.Unlock()
		if err != nil {
			return nil, err
		}
		return nil, ErrTransportDown
	}
	if pending.invalidated || m.observedEpoch > attemptedEpoch || (err == nil && newTransport.TransportEpoch < m.observedEpoch) {
		m.pendingOperation = nil
		m.pendingEpoch = 0
		m.transport.TransportEpoch = m.observedEpoch
		m.toDisabledLocked()
		m.createCleanupLocked(target, m.observedEpoch, pending.OperationID)
		m.stateMu.Unlock()
		return nil, ErrTransportDown
	}
	if err == nil && newTransport.TransportEpoch >= attemptedEpoch {
		m.transport = newTransport
		m.observedEpoch = newTransport.TransportEpoch
		m.pendingEpoch = 0
		m.state = Active
		pending.Transport = newTransport
		if m.cleanup != nil {
			m.cleanup.TransportEpoch = newTransport.TransportEpoch
		}
		m.armIdleTimerLocked()
		m.stateMu.Unlock()
		return pending, nil
	}
	if err == nil {
		err = ErrTransportDown
	}
	if errors.Is(err, ErrTransportDown) {
		m.lastRefreshTransportDown = true
	} else {
		m.refreshFailed = true
	}
	m.transport.TransportEpoch = attemptedEpoch
	m.pendingEpoch = 0
	m.toDisabledLocked()
	m.createCleanupLocked(target, attemptedEpoch, "")
	m.pendingOperation = nil
	m.stateMu.Unlock()
	return nil, err
}

func dispatchOnce(ctx context.Context, target byte, dispatch DispatchFunc) DispatchOutcome {
	if ctx == nil {
		ctx = context.Background()
	}
	if dispatch == nil {
		return DispatchOutcome{Err: ErrTransportDown}
	}
	return dispatch(ctx, target)
}

// bindPendingContext installs the cancellation handle while stateMu is held.
// A lifecycle event that wins the race invalidates and cancels this context
// before the dispatcher reaches its final pre-emission admission.
func (m *Manager) bindPendingContext(ctx context.Context, pending *pendingOperation) (context.Context, context.CancelCauseFunc) {
	operationCtx, cancel := context.WithCancelCause(ctx)
	m.stateMu.Lock()
	if m.pendingOperation != pending || pending.invalidated {
		cancel(ErrTransportDown)
	} else {
		pending.cancel = cancel
	}
	m.stateMu.Unlock()
	return operationCtx, cancel
}

func (m *Manager) releaseOwnerLocked() {
	if m.idleTimer != nil {
		m.idleTimer.Stop()
		m.idleTimer = nil
	}
	m.idleTimerGen++
	if m.mutexHeld {
		m.mu.Unlock()
		m.mutexHeld = false
	}
	m.activeToken = ""
	m.activeTarget = 0
}

func (m *Manager) createCleanupLocked(target byte, epoch uint64, sourceOperationID string) {
	if m.cleanup != nil {
		return
	}
	m.cleanup = &cleanupObligation{
		CleanupObligationSnapshot: CleanupObligationSnapshot{
			Target: target, GatewayCleanupAttemptID: newCleanupAttemptID(), TransportEpoch: epoch,
		},
		sourceOperationID: sourceOperationID,
	}
	m.state = Disabled
}

func (m *Manager) clearCleanupForOperationLocked(operationID string) {
	if m.cleanup != nil && m.cleanup.sourceOperationID == operationID {
		m.cleanup = nil
	}
}

func (m *Manager) attemptCleanup(_ context.Context, epoch uint64) {
	m.stateMu.Lock()
	if m.cleanup == nil || m.restarted || (m.cleanup.LastAttempted && m.cleanup.LastAttemptEpoch == epoch) || m.cleanupDispatch == nil {
		m.stateMu.Unlock()
		return
	}
	attemptID := m.cleanup.GatewayCleanupAttemptID
	target := m.cleanup.Target
	m.cleanup.LastAttemptEpoch = epoch
	m.cleanup.LastAttempted = true
	m.cleanup.TransportEpoch = epoch
	dispatch := m.cleanupDispatch
	pending := &pendingOperation{PendingOperationSnapshot: PendingOperationSnapshot{OperationID: newCleanupAttemptID(), Kind: OperationDisable, Target: target, Transport: m.transport}, cleanup: true}
	m.pendingOperation = pending
	m.stateMu.Unlock()
	cleanupCtx, timeout := context.WithTimeout(context.Background(), 5*time.Second)
	defer timeout()
	operationCtx, cancel := m.bindPendingContext(cleanupCtx, pending)
	defer cancel(nil)
	outcome := dispatchOnce(withPendingOperation(operationCtx, pending.OperationID), target, dispatch)
	m.stateMu.Lock()
	if m.pendingOperation == pending {
		m.pendingOperation = nil
	}
	m.stateMu.Unlock()
	m.recordCleanupOutcome(attemptID, epoch, outcome)
}

// recordCleanupOutcome retains exact native evidence without treating ACK or
// NAK as proof that the device-side session settled.
func (m *Manager) recordCleanupOutcome(attemptID string, epoch uint64, outcome DispatchOutcome) {
	m.stateMu.Lock()
	defer m.stateMu.Unlock()
	if m.cleanup == nil || m.cleanup.GatewayCleanupAttemptID != attemptID || !m.cleanup.LastAttempted || m.cleanup.LastAttemptEpoch != epoch {
		return
	}
	m.cleanup.LastEmitted = outcome.Emitted
	m.cleanup.LastNative = outcome.Native
	m.cleanup.LastErr = outcome.Err
}

var fallbackAttemptID atomic.Uint64

func newCleanupAttemptID() string {
	if id, err := newIssuerToken(); err == nil {
		return id
	}
	return fmt.Sprintf("fallback-%d", fallbackAttemptID.Add(1))
}

// newIssuerToken returns 16 random bytes hex-encoded (32 chars).
func newIssuerToken() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(b[:]), nil
}
