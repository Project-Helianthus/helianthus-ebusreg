package b503session

import (
	"context"
	"sync"
	"time"
)

// TransportKey identifies a specific adapter incarnation on the bus.
// TransportEpoch advances on every successful transport reconnect so the
// session FSM can detect that its owning incarnation is gone.
type TransportKey struct {
	AdapterInstanceID string
	TransportEpoch    uint64
}

// SessionKey is the (transport_key, issuer_token) tuple returned by
// Enable. A caller holding the full SessionKey may Disable; only the
// transport_key is required for Read (spec §6.2).
type SessionKey struct {
	Transport   TransportKey
	IssuerToken string
}

// RefreshFunc is invoked during epoch-advance handling. It MUST return
// one of:
//   - (newTransportKey, nil)           refresh succeeded; session re-homes.
//   - (TransportKey{}, ErrTransportDown) transport went down during refresh.
//   - (TransportKey{}, other-error)    refresh failed for another reason;
//     session becomes permanently busy
//     until the next Enable.
type RefreshFunc func(ctx context.Context) (TransportKey, error)

// NativeResult is the exact device-side wire outcome relevant to the B503
// lifecycle. ACK and NAK are evidence about the native exchange only; neither
// proves that device-side session state has settled.
type NativeResult uint8

const (
	// NativeAmbiguous is the native ambiguous value of the Vaillant owner.
	NativeAmbiguous NativeResult = iota
	// NativeACK is the native ack value of the Vaillant owner.
	NativeACK
	// NativeNAK is the native nak value of the Vaillant owner.
	NativeNAK
)

// DispatchOutcome records the conservative wire boundary for one native B503
// operation. Emitted becomes true as soon as bus.Send is entered; it does not
// claim that the target received or accepted the frame.
type DispatchOutcome struct {
	Response  []byte
	Emitted   bool
	Native    NativeResult
	Transport TransportKey
	Err       error
}

// DispatchFunc emits exactly one target-bound native operation.
type DispatchFunc func(context.Context, byte) DispatchOutcome

type pendingOperationContextKey struct{}

// PendingOperationIDFromContext returns the process-local lifecycle operation
// identity installed by Manager immediately before it invokes a dispatcher.
// Dispatchers use it only to reject an operation invalidated while waiting for
// transport quiescence; it grants no caller authority.
func PendingOperationIDFromContext(ctx context.Context) (string, bool) {
	if ctx == nil {
		return "", false
	}
	operationID, ok := ctx.Value(pendingOperationContextKey{}).(string)
	return operationID, ok && operationID != ""
}

func withPendingOperation(ctx context.Context, operationID string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, pendingOperationContextKey{}, operationID)
}

// CleanupObligationSnapshot is read-only evidence for deterministic tests and
// diagnostics. It carries no caller authority.
type CleanupObligationSnapshot struct {
	Target                  byte
	GatewayCleanupAttemptID string
	TransportEpoch          uint64
	OriginEmitted           bool
	OriginNative            NativeResult
	OriginErr               error
	LastAttempted           bool
	LastAttemptEpoch        uint64
	LastEmitted             bool
	LastNative              NativeResult
	LastErr                 error
}

type cleanupObligation struct {
	CleanupObligationSnapshot
	sourceOperationID string
}

// OperationKind identifies the native lifecycle operation currently admitted
// by Manager.
type OperationKind uint8

const (
	// OperationEnable is the operation enable value of the Vaillant owner.
	OperationEnable OperationKind = iota + 1
	// OperationRead is the operation read value of the Vaillant owner.
	OperationRead
	// OperationDisable is the operation disable value of the Vaillant owner.
	OperationDisable
)

// PendingOperationSnapshot is non-authoritative diagnostic evidence. The
// operation ID is process-local and grants no caller authority.
type PendingOperationSnapshot struct {
	OperationID string
	Kind        OperationKind
	Target      byte
	Transport   TransportKey
	Emitted     bool
}

type pendingOperation struct {
	PendingOperationSnapshot
	token       string
	cancel      context.CancelCauseFunc
	invalidated bool
	cleanup     bool
}

// Manager is the single-owner live-monitor session FSM.
//
// Two mutexes cooperate:
//
//   - mu (liveMonitorMu): the ownership gate. Acquired on the
//     Idle->Enabling transition; released exactly once on entry to
//     Disabled from a held-owner state (Enabling, Active, or Refreshing).
//     This is distinct from the B524 readMu used by the poller. The
//     lock lifecycle is tracked by mutexHeld so cleanup paths never
//     double-release.
//
//   - stateMu: field-level protection for state, transport, activeToken,
//     idleTimer, and mutexHeld. All public methods acquire stateMu for
//     the duration of their FSM inspection/mutation; they acquire or
//     release mu as the FSM dictates.
//
// This separation avoids the deadlock trap of having the idle-timer
// callback try to re-enter a mutex that is semantically held by a
// different logical agent (the session owner).
type Manager struct {
	mu sync.Mutex // liveMonitorMu — ownership gate, distinct from B524 readMu.
	// emissionMu serializes a lifecycle advance with the narrow interval from
	// pending-operation admission through entry into the generation-bound bus
	// substrate. It is distinct from both the owner gate and B524 readMu: the
	// dispatcher takes it only after readMu is acquired, while lifecycle events
	// take it before stateMu and never acquire readMu.
	emissionMu sync.Mutex

	stateMu      sync.Mutex
	state        State
	transport    TransportKey
	activeToken  string
	idleTimeout  time.Duration
	idleTimer    *time.Timer
	idleTimerGen uint64 // generation counter — stale callback guard
	// idleExpiryPending records that the timer elapsed while an admitted Read
	// was awaiting its native outcome. The callback must not tear down that
	// owner; ReadOperation resolves the expiry after the outcome.
	idleExpiryPending bool
	refresh           RefreshFunc
	mutexHeld         bool
	// refreshFailed sticks to true when a refresh returned a non-nil,
	// non-ErrTransportDown error. Subsequent Reads then surface
	// ErrSessionBusy until ResetForRestart or a new Enable. It is
	// cleared by Enable and ResetForRestart.
	refreshFailed bool
	// lastRefreshTransportDown is set when OnEpochAdvance's refresh returned
	// ErrTransportDown; cleared on Enable / ResetForRestart. See spec §7.1
	// and plan AD14 for the resolver-layer contract this enables.
	lastRefreshTransportDown bool

	activeTarget     byte
	pendingOperation *pendingOperation
	cleanup          *cleanupObligation
	cleanupDispatch  DispatchFunc
	observedEpoch    uint64
	pendingEpoch     uint64
	qualified        map[byte]struct{}
	restarted        bool
	restartFences    map[byte]struct{}
}

// StatusSnapshot is one coherent, public observation of the live-monitor
// session. State and Owned are captured from that same observation.
//
// Callers that render or serialize both values must use StatusSnapshot rather
// than combining State and IsOwned, because a lifecycle transition may occur
// between separate calls to those accessors.
type StatusSnapshot struct {
	State State
	Owned bool
}

// New constructs a Manager bound to the initial transport incarnation.
// idleTimeout is the Active->Disabled auto-disable window (spec §7.6;
// production value 30s). refresh is invoked during epoch-advance handling.
func New(transport TransportKey, idleTimeout time.Duration, refresh RefreshFunc) *Manager {
	return &Manager{
		state:         Idle,
		transport:     transport,
		idleTimeout:   idleTimeout,
		refresh:       refresh,
		observedEpoch: transport.TransportEpoch,
		qualified:     make(map[byte]struct{}),
		restartFences: make(map[byte]struct{}),
	}
}

// SetCleanupDispatcher installs the process-local native disable hook used by
// idle expiry and ambiguous enable cleanup. It does not itself emit I/O.
func (m *Manager) SetCleanupDispatcher(dispatch DispatchFunc) {
	m.stateMu.Lock()
	m.cleanupDispatch = dispatch
	m.stateMu.Unlock()
}

// MarkQualifiedTarget records a finite registry-qualified B503 target. After
// ResetForRestart, newly enumerated qualified targets inherit the restart
// UNKNOWN fence until separately authorized recovery outside public v1.
func (m *Manager) MarkQualifiedTarget(target byte) {
	m.stateMu.Lock()
	m.qualified[target] = struct{}{}
	if m.restarted {
		m.restartFences[target] = struct{}{}
	}
	m.stateMu.Unlock()
}

// TargetBlocked reports whether cleanup or restart recovery makes target
// unavailable. This is deliberately target-specific.
func (m *Manager) TargetBlocked(target byte) bool {
	m.stateMu.Lock()
	defer m.stateMu.Unlock()
	if _, fenced := m.restartFences[target]; fenced {
		return true
	}
	// The cleanup write remains target-bound, but the physical live-monitor
	// slot is single-owner per transport. Any retained cleanup therefore makes
	// every target on that transport unavailable for a new claim.
	return m.cleanup != nil
}

// CleanupObligation returns a copy of the current process-local obligation.
func (m *Manager) CleanupObligation() (CleanupObligationSnapshot, bool) {
	m.stateMu.Lock()
	defer m.stateMu.Unlock()
	if m.cleanup == nil {
		return CleanupObligationSnapshot{}, false
	}
	return m.cleanup.CleanupObligationSnapshot, true
}

// PendingOperation returns the current typed target-bound operation, if any.
func (m *Manager) PendingOperation() (PendingOperationSnapshot, bool) {
	m.stateMu.Lock()
	defer m.stateMu.Unlock()
	if m.pendingOperation == nil {
		return PendingOperationSnapshot{}, false
	}
	return m.pendingOperation.PendingOperationSnapshot, true
}

// AdmitPendingEmission marks the exact pending lifecycle operation as emitted
// at the dispatcher's post-quiesce boundary. A transport lifecycle event that
// already invalidated the operation makes admission fail without wire I/O.
func (m *Manager) AdmitPendingEmission(operationID string, target byte) bool {
	m.stateMu.Lock()
	defer m.stateMu.Unlock()
	return m.admitPendingEmissionLocked(operationID, target)
}

// BeginPendingEmission admits one typed lifecycle operation and holds the
// lifecycle-emission fence until the dispatcher has entered its generation-bound
// send substrate. Call the returned release exactly once after bus.Send returns.
// A lifecycle advance which acquired the fence first invalidates the pending
// operation, so this method returns false without a wire attempt.
func (m *Manager) BeginPendingEmission(operationID string, target byte) (release func(), ok bool) {
	m.emissionMu.Lock()
	m.stateMu.Lock()
	ok = m.admitPendingEmissionLocked(operationID, target)
	m.stateMu.Unlock()
	if !ok {
		m.emissionMu.Unlock()
		return nil, false
	}
	return func() { m.emissionMu.Unlock() }, true
}

func (m *Manager) admitPendingEmissionLocked(operationID string, target byte) bool {
	pending := m.pendingOperation
	if pending == nil || pending.invalidated || pending.OperationID != operationID || pending.Target != target {
		return false
	}
	switch pending.Kind {
	case OperationEnable:
		if m.state != Enabling {
			return false
		}
	case OperationRead:
		if m.state != Active {
			return false
		}
	case OperationDisable:
		// A normal Disable retains the current owner's Active gate. The only
		// ownerless disable admission is the Manager-created defensive cleanup
		// after an emitted lifecycle operation left settlement unproven.
		if pending.cleanup {
			if m.state != Disabled || m.mutexHeld {
				return false
			}
			break
		}
		if m.state != Active || !m.mutexHeld {
			return false
		}
	default:
		return false
	}
	pending.Emitted = true
	return true
}

// DispatchEpoch is the latest transport epoch observed by lifecycle routing.
// Dispatchers compare this value with request-side issue metadata so a stale
// completion cannot be rescued by a still-bound old owner key.
func (m *Manager) DispatchEpoch() uint64 {
	m.stateMu.Lock()
	defer m.stateMu.Unlock()
	return m.observedEpoch
}

// LifecycleTransportKey identifies the latest admitted transport generation
// without rebinding a surviving owner. Driver lifecycle withdrawal correlates
// against this key; owner operations still use TransportKey and perform the
// documented bounded refresh before dispatch.
func (m *Manager) LifecycleTransportKey() TransportKey {
	m.stateMu.Lock()
	defer m.stateMu.Unlock()
	return TransportKey{
		AdapterInstanceID: m.transport.AdapterInstanceID,
		TransportEpoch:    m.observedEpoch,
	}
}

// Enable transitions Idle -> Enabling -> Active and returns a SessionKey
// with a freshly-minted 16-byte hex issuer_token. Returns ErrSessionBusy
// if the FSM is not Idle.
//
// This lower-level transition remains available to tests that isolate the
// owner gate. Production B503 paths use EnableOperation so native ACK/NAK and
// emission state own the transition.
func (m *Manager) Enable(ctx context.Context) (SessionKey, error) {
	m.stateMu.Lock()
	defer m.stateMu.Unlock()

	if m.state != Idle {
		return SessionKey{}, ErrSessionBusy
	}

	// Idle -> Enabling -> Active, atomically from the caller's perspective.
	m.state = Enabling
	// Acquire ownership gate. The mutex is held logically by the session,
	// not by this goroutine; we release it on entry to Disabled.
	m.mu.Lock()
	m.mutexHeld = true

	token, err := newIssuerToken()
	if err != nil {
		// Roll back: release gate, revert state.
		m.mu.Unlock()
		m.mutexHeld = false
		m.state = Idle
		return SessionKey{}, err
	}
	m.activeToken = token
	m.activeTarget = 0
	m.state = Active
	m.refreshFailed = false
	m.lastRefreshTransportDown = false
	m.armIdleTimerLocked()

	return SessionKey{Transport: m.transport, IssuerToken: token}, nil
}

// EnableOperation owns the target-bound enable handshake. A pre-emission
// cancellation returns to Idle with no cleanup. Every emitted outcome creates
// a process-local cleanup obligation because native ACK/NAK does not prove the
// device-side session state settled. An exact ACK may activate the caller
// handle, while every emitted failure releases caller ownership and retains
// cleanup under UNKNOWN.
func (m *Manager) EnableOperation(ctx context.Context, target byte, dispatch DispatchFunc) (SessionKey, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	m.MarkQualifiedTarget(target)
	m.stateMu.Lock()
	if _, fenced := m.restartFences[target]; fenced {
		m.stateMu.Unlock()
		return SessionKey{}, ErrCleanupPending
	}
	if m.cleanup != nil {
		m.stateMu.Unlock()
		return SessionKey{}, ErrCleanupPending
	}
	if m.state != Idle {
		m.stateMu.Unlock()
		return SessionKey{}, ErrSessionBusy
	}
	m.state = Enabling
	m.mu.Lock()
	m.mutexHeld = true
	token, err := newIssuerToken()
	if err != nil {
		m.releaseOwnerLocked()
		m.state = Idle
		m.stateMu.Unlock()
		return SessionKey{}, err
	}
	opID := newCleanupAttemptID()
	pending := &pendingOperation{
		PendingOperationSnapshot: PendingOperationSnapshot{
			OperationID: opID, Kind: OperationEnable, Target: target, Transport: m.transport,
		},
		token: token,
	}
	operationCtx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)
	pending.cancel = cancel
	m.pendingOperation = pending
	m.activeToken = token
	m.activeTarget = target
	m.refreshFailed = false
	m.lastRefreshTransportDown = false
	m.stateMu.Unlock()

	outcome := dispatchOnce(withPendingOperation(operationCtx, pending.OperationID), target, dispatch)

	m.stateMu.Lock()
	if m.pendingOperation != pending {
		m.stateMu.Unlock()
		if outcome.Err != nil {
			return SessionKey{}, outcome.Err
		}
		return SessionKey{}, ErrTransportDown
	}
	pending.Emitted = outcome.Emitted
	m.pendingOperation = nil
	advanced := m.pendingEpoch > pending.Transport.TransportEpoch
	if advanced {
		m.transport.TransportEpoch = m.pendingEpoch
	}
	if m.state != Enabling {
		if !outcome.Emitted {
			m.clearCleanupForOperationLocked(opID)
			if m.cleanup == nil {
				m.state = Idle
			}
		} else if m.cleanup != nil && m.cleanup.sourceOperationID == opID {
			m.cleanup.OriginEmitted = true
			m.cleanup.OriginNative = outcome.Native
			m.cleanup.OriginErr = outcome.Err
		}
		m.stateMu.Unlock()
		if outcome.Err != nil {
			return SessionKey{}, outcome.Err
		}
		return SessionKey{}, ErrTransportDown
	}

	if !outcome.Emitted {
		m.releaseOwnerLocked()
		m.state = Idle
		m.stateMu.Unlock()
		return SessionKey{}, outcome.Err
	}

	// Once the enable enters bus.Send, retain a fresh target/epoch-bound
	// cleanup obligation regardless of ACK, NAK, or ambiguous completion.
	m.createCleanupLocked(target, m.transport.TransportEpoch, opID)
	m.cleanup.OriginEmitted = outcome.Emitted
	m.cleanup.OriginNative = outcome.Native
	m.cleanup.OriginErr = outcome.Err
	if outcome.Native == NativeACK && outcome.Err == nil && !advanced {
		m.state = Active
		m.armIdleTimerLocked()
		key := SessionKey{Transport: m.transport, IssuerToken: token}
		m.stateMu.Unlock()
		return key, nil
	}
	m.releaseOwnerLocked()
	m.state = Disabled
	epoch := m.transport.TransportEpoch
	m.stateMu.Unlock()
	m.attemptCleanup(ctx, epoch)
	if outcome.Err != nil {
		return SessionKey{}, outcome.Err
	}
	return SessionKey{}, ErrTransportDown
}
