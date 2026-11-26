package b503session

import (
	"context"
	"errors"
)

// ReadOperation validates the target-bound owner, performs at most one epoch
// refresh when required, and dispatches the admitted read exactly once.
func (m *Manager) ReadOperation(ctx context.Context, target byte, dispatch DispatchFunc) ([]byte, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	m.MarkQualifiedTarget(target)
	pending, err := m.beginOwnerOperation(ctx, target, nil, OperationRead)
	if err != nil {
		return nil, err
	}
	operationCtx, cancel := m.bindPendingContext(ctx, pending)
	defer cancel(nil)
	outcome := dispatchOnce(withPendingOperation(operationCtx, pending.OperationID), target, dispatch)
	m.stateMu.Lock()
	if m.pendingOperation != pending {
		m.stateMu.Unlock()
		if outcome.Err != nil {
			return nil, outcome.Err
		}
		return nil, ErrTransportDown
	}
	pending.Emitted = outcome.Emitted
	m.pendingOperation = nil
	if outcome.Err != nil {
		if errors.Is(outcome.Err, ErrTransportDown) {
			epoch := m.transport.TransportEpoch
			m.toDisabledLocked()
			m.createCleanupLocked(target, epoch, pending.OperationID)
			m.lastRefreshTransportDown = true
			m.stateMu.Unlock()
			return nil, outcome.Err
		}
		epoch, cleanup := m.consumeIdleExpiryAfterReadLocked(target, pending.OperationID)
		m.stateMu.Unlock()
		if cleanup {
			m.attemptCleanup(context.Background(), epoch)
		}
		return nil, outcome.Err
	}
	if outcome.Native != NativeACK {
		epoch, cleanup := m.consumeIdleExpiryAfterReadLocked(target, pending.OperationID)
		m.stateMu.Unlock()
		if cleanup {
			m.attemptCleanup(context.Background(), epoch)
		}
		return nil, ErrTransportDown
	}
	if m.state != Active || m.transport != pending.Transport || m.activeTarget != target {
		epoch, cleanup := m.consumeIdleExpiryAfterReadLocked(target, pending.OperationID)
		m.stateMu.Unlock()
		if cleanup {
			m.attemptCleanup(context.Background(), epoch)
		}
		return nil, ErrTransportDown
	}
	m.idleExpiryPending = false
	m.armIdleTimerLocked()
	m.stateMu.Unlock()
	return outcome.Response, nil
}

// DisableOperation owns the explicit native disable. The owner remains held
// while the dispatcher waits for poll quiescence, then is released after the
// emitted or pre-emission terminal outcome. Its exact outcome is recorded, but
// ACK or NAK alone never clears cleanup or makes re-Enable admissible.
func (m *Manager) DisableOperation(ctx context.Context, key SessionKey, target byte, dispatch DispatchFunc) error {
	if ctx == nil {
		ctx = context.Background()
	}
	m.MarkQualifiedTarget(target)
	pending, err := m.beginOwnerOperation(ctx, target, &key, OperationDisable)
	if err != nil {
		return err
	}
	operationCtx, cancel := m.bindPendingContext(ctx, pending)
	defer cancel(nil)
	// Once the current owner has admitted an explicit Disable, idle expiry
	// must not race it to a second terminal path while it waits for poll
	// quiescence. The owner remains held until emission or invalidation.
	m.stateMu.Lock()
	if m.pendingOperation == pending && m.mutexHeld {
		if m.idleTimer != nil {
			m.idleTimer.Stop()
			m.idleTimer = nil
		}
		m.idleTimerGen++
	}
	m.stateMu.Unlock()
	outcome := dispatchOnce(withPendingOperation(operationCtx, pending.OperationID), target, dispatch)
	m.stateMu.Lock()
	if m.pendingOperation != pending {
		m.stateMu.Unlock()
		if outcome.Err != nil {
			return outcome.Err
		}
		return ErrTransportDown
	}
	pending.Emitted = pending.Emitted || outcome.Emitted
	m.pendingOperation = nil
	if !pending.Emitted && m.state == Active && m.mutexHeld && m.activeTarget == target {
		// The dispatcher never reached bus.Send. Keep the authenticated owner so
		// the caller can retry the exact Disable once poll quiescence or its own
		// bounded context permits it. A cleanup attempt would falsely record a
		// wire event that did not occur and would discard the only valid token.
		m.idleExpiryPending = false
		m.armIdleTimerLocked()
		m.stateMu.Unlock()
		if outcome.Err != nil {
			return outcome.Err
		}
		return ErrTransportDown
	}
	epoch := m.transport.TransportEpoch
	if m.mutexHeld {
		m.releaseOwnerLocked()
	}
	m.state = Disabled
	m.createCleanupLocked(target, epoch, pending.OperationID)
	attemptID := m.cleanup.GatewayCleanupAttemptID
	m.cleanup.LastAttemptEpoch = epoch
	m.cleanup.LastAttempted = true
	m.stateMu.Unlock()
	m.recordCleanupOutcome(attemptID, epoch, outcome)
	return outcome.Err
}

// Disable transitions Active -> Disabled -> Idle. Requires a full
// SessionKey match. ErrWrongToken for issuer_token mismatch; ErrNotActive
// if the FSM is not currently Active.
//
// This lower-level transition remains available to tests that isolate the
// owner gate. Production paths use DisableOperation so a native disable is
// emitted and evidenced.
func (m *Manager) Disable(key SessionKey) error {
	m.stateMu.Lock()
	defer m.stateMu.Unlock()

	if m.state == Refreshing {
		return ErrSessionBusy
	}
	if m.state != Active {
		return ErrNotActive
	}
	if key.Transport != m.transport {
		return ErrWrongToken
	}
	if key.IssuerToken != m.activeToken {
		return ErrWrongToken
	}
	m.toDisabledLocked()
	m.toIdleLocked()
	return nil
}

// Read checks that the FSM is Active for the provided transport and
// resets the idle timer. issuer_token is intentionally not part of the
// Read contract (spec §6.2).
//
// This lower-level transition remains available to tests that isolate the
// owner gate. Production paths use ReadOperation so refresh and the triggering
// dispatch remain atomic.
func (m *Manager) Read(transport TransportKey) error {
	m.stateMu.Lock()
	defer m.stateMu.Unlock()

	if m.refreshFailed || m.state == Refreshing {
		return ErrSessionBusy
	}
	if m.state != Active {
		return ErrNotActive
	}
	if transport != m.transport {
		return ErrNotActive
	}
	m.armIdleTimerLocked()
	return nil
}

// State returns the public FSM state. Refreshing is visible while an epoch
// refresh holds the gate, so a caller never observes Disabled with ownership.
func (m *Manager) State() State {
	m.stateMu.Lock()
	defer m.stateMu.Unlock()
	if m.state == Disabled && (m.cleanup != nil || m.restarted) {
		return Idle
	}
	return m.state
}

// IsOwned reports whether the session gate is currently held by any
// owner — true when FSM is Enabling, Active, or Refreshing. Resolver-layer code (e.g.
// capability probe in mcp/vaillant_b503.go) uses this to classify the
// surface as SESSION_BUSY during the full held-owner window, not just
// when State() reports Active/Enabling. Without this, a slow refresh
// could publish AVAILABLE while a concurrent Enable would still return
// ErrSessionBusy.
func (m *Manager) IsOwned() bool {
	m.stateMu.Lock()
	defer m.stateMu.Unlock()
	return m.mutexHeld
}

// StatusSnapshot returns the public state and ownership gate from one stateMu
// critical section. This is the stable observation used by MCP and
// GraphQL session-status surfaces.
func (m *Manager) StatusSnapshot() StatusSnapshot {
	m.stateMu.Lock()
	defer m.stateMu.Unlock()
	state := m.state
	if state == Disabled && (m.cleanup != nil || m.restarted) {
		state = Idle
	}
	return StatusSnapshot{
		State: state,
		Owned: m.mutexHeld,
	}
}

// TransportKey returns the manager's current transport key. Exposed so
// resolver-layer callers (mcp/vaillant_b503.go) can pass the correct key to
// Read/Disable after OnEpochAdvance has re-homed the session.
func (m *Manager) TransportKey() TransportKey {
	m.stateMu.Lock()
	defer m.stateMu.Unlock()
	return m.transport
}

// LastRefreshTransportDown reports whether the most recent OnEpochAdvance
// resolved with ErrTransportDown. Callers use this to distinguish
// transport-down teardown from ordinary Idle state when surfacing a public
// error code (spec §7.1 / plan AD14). Cleared by Enable and ResetForRestart.
func (m *Manager) LastRefreshTransportDown() bool {
	m.stateMu.Lock()
	defer m.stateMu.Unlock()
	return m.lastRefreshTransportDown
}

// OnTransportDisconnect performs owner-conditional cleanup (spec §7.4).
// A pre-emission Enable returns directly to Idle without cleanup. Every other
// held-owner state releases into internal Disabled with retained cleanup.
// An ownerless state is unchanged.
func (m *Manager) OnTransportDisconnect() {
	m.stateMu.Lock()
	defer m.stateMu.Unlock()
	if !m.mutexHeld {
		return
	}
	if pending := m.pendingOperation; pending != nil {
		pending.invalidated = true
		if pending.cancel != nil {
			pending.cancel(ErrTransportDown)
		}
	}
	if pending := m.pendingOperation; pending != nil && pending.Kind == OperationEnable && !pending.Emitted {
		m.pendingOperation = nil
		m.pendingEpoch = 0
		m.releaseOwnerLocked()
		m.clearCleanupForOperationLocked(pending.OperationID)
		m.state = Idle
		return
	}
	target := m.activeTarget
	sourceID := ""
	if m.pendingOperation != nil && m.pendingOperation.Kind == OperationEnable {
		target = m.pendingOperation.Target
		sourceID = m.pendingOperation.OperationID
	}
	epoch := m.transport.TransportEpoch
	m.toDisabledLocked()
	m.createCleanupLocked(target, epoch, sourceID)
}

// OnEpochAdvance records the latest transport epoch. A surviving owner is
// rebound only when its next admitted ReadOperation or current-owner
// DisableOperation invokes refreshOwner. With no owner, retained cleanup is
// fenced without emitting any recovery write; reconnect and restart cannot
// establish settlement or reconstruct caller ownership.
func (m *Manager) OnEpochAdvance(_ context.Context, newEpoch uint64) {
	m.emissionMu.Lock()
	defer m.emissionMu.Unlock()
	m.stateMu.Lock()
	if newEpoch <= m.observedEpoch {
		m.stateMu.Unlock()
		return
	}
	m.observedEpoch = newEpoch
	m.lastRefreshTransportDown = false
	if pending := m.pendingOperation; pending != nil {
		// The dispatcher revalidates this context immediately before bus.Send.
		// A rollover therefore cancels work still waiting for quiescence; a
		// Send already entered remains conservatively recorded as emitted.
		pending.invalidated = true
		if pending.cancel != nil {
			pending.cancel(ErrTransportDown)
		}
	}
	if pending := m.pendingOperation; m.state == Enabling && pending != nil && pending.Kind == OperationEnable && !pending.Emitted {
		m.transport.TransportEpoch = newEpoch
		m.pendingEpoch = 0
		m.pendingOperation = nil
		m.releaseOwnerLocked()
		m.clearCleanupForOperationLocked(pending.OperationID)
		m.state = Idle
		m.stateMu.Unlock()
		return
	}
	if m.mutexHeld {
		// A surviving owner is rebound only by its next admitted READ or
		// current-owner DISABLE. ENABLE is never a refresh trigger.
		m.pendingEpoch = newEpoch
		m.stateMu.Unlock()
		return
	}
	m.transport.TransportEpoch = newEpoch
	m.pendingEpoch = 0
	m.stateMu.Unlock()
}

// ResetForRestart destroys caller handles and process-local cleanup, releases
// a held gate, and fences known or supplied qualified targets. It emits no
// native operation. The stable public session observation is Idle/owned=false
// while the internal state remains Disabled until out-of-band recovery.
func (m *Manager) ResetForRestart(qualifiedTargets ...byte) {
	m.stateMu.Lock()
	defer m.stateMu.Unlock()
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
	if pending := m.pendingOperation; pending != nil {
		pending.invalidated = true
		if pending.cancel != nil {
			pending.cancel(ErrTransportDown)
		}
	}
	m.pendingOperation = nil
	m.cleanup = nil
	m.refreshFailed = false
	m.lastRefreshTransportDown = false
	m.pendingEpoch = 0
	m.idleExpiryPending = false
	m.state = Disabled
	m.restarted = true
	for target := range m.qualified {
		m.restartFences[target] = struct{}{}
	}
	for _, target := range qualifiedTargets {
		m.qualified[target] = struct{}{}
		m.restartFences[target] = struct{}{}
	}
}

// --- internal helpers (all require stateMu held) ---

// toDisabledLocked transitions to Disabled from any held-owner state and
// releases the ownership gate exactly once. Bumps idleTimerGen so any
// in-flight stale callback is invalidated.
func (m *Manager) toDisabledLocked() {
	if m.idleTimer != nil {
		m.idleTimer.Stop()
		m.idleTimer = nil
	}
	m.idleTimerGen++
	m.idleExpiryPending = false
	m.state = Disabled
	if m.mutexHeld {
		m.mu.Unlock()
		m.mutexHeld = false
	}
	m.activeToken = ""
	m.activeTarget = 0
}

// toIdleLocked finalises Disabled -> Idle.
func (m *Manager) toIdleLocked() {
	if m.cleanup == nil {
		m.state = Idle
	}
}

// consumeIdleExpiryAfterReadLocked turns a timer that elapsed under a pending
// read into one cleanup only when the read did not successfully extend the
// still-current owner. Caller must hold stateMu after removing that read from
// pendingOperation.
func (m *Manager) consumeIdleExpiryAfterReadLocked(target byte, sourceOperationID string) (uint64, bool) {
	if !m.idleExpiryPending {
		return 0, false
	}
	m.idleExpiryPending = false
	if m.state != Active || !m.mutexHeld || m.activeTarget != target {
		return 0, false
	}
	epoch := m.transport.TransportEpoch
	m.toDisabledLocked()
	m.createCleanupLocked(target, epoch, sourceOperationID)
	return epoch, true
}
