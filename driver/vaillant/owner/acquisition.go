package owner

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Project-Helianthus/helianthus-ebusgo/protocol/vaillant/b503"
	"github.com/Project-Helianthus/helianthus-ebusgo/transport"

	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/owner/b503session"
)

const (
	// ConfiguredMonitorTarget is the configured live-monitor target fenced at startup.
	ConfiguredMonitorTarget byte = 0x08

	acquisitionAdapterID = "ebus"
	historySessionIdle   = 30 * time.Second
)

// Acquisition is one Vaillant acquisition owner tied to its own history session.
type Acquisition struct {
	poller       *vaillantSemanticPoller
	substrate    Substrate
	manager      *b503session.Manager
	dispatcher   *RawFrameDispatcher
	gate         *dispatcherGate
	refreshCalls atomic.Int32
}

type dispatcherGate struct {
	mu      *sync.Mutex
	waiting atomic.Int32
}

func (g *dispatcherGate) Lock() {
	if g == nil || g.mu == nil {
		return
	}
	g.waiting.Add(1)
	g.mu.Lock()
	g.waiting.Add(-1)
}

func (g *dispatcherGate) Unlock() {
	if g == nil || g.mu == nil {
		return
	}
	g.mu.Unlock()
}

// AcquisitionView is the source and published zone and hot-water view of one owner.
type AcquisitionView struct {
	Source   byte
	Zones    int
	HotWater bool
}

// NewAcquisition returns an owner for cfg and substrate, tied to its own history session.
func NewAcquisition(cfg Config, substrate Substrate) *Acquisition {
	acquired := &Acquisition{
		poller:    newVaillantSemanticPoller(cfg, substrate, newRetainedSemanticState(), newPublicationHub(), nil),
		substrate: substrate,
	}
	acquired.manager = newHistorySession(acquired)
	acquired.gate = &dispatcherGate{mu: &acquired.poller.readMu}
	acquired.dispatcher = newHistoryDispatcher(cfg, substrate, acquired.poller, acquired.manager, acquired.gate)
	acquired.manager.SetCleanupDispatcher(func(ctx context.Context, target byte) b503session.DispatchOutcome {
		return InvokeB503Operation(ctx, acquired.dispatcher, target, b503.EncodeLiveMonitorMain())
	})
	return acquired
}

func newHistorySession(acquired *Acquisition) *b503session.Manager {
	manager := b503session.New(
		b503session.TransportKey{AdapterInstanceID: acquisitionAdapterID},
		historySessionIdle,
		func(context.Context) (b503session.TransportKey, error) {
			acquired.refreshCalls.Add(1)
			return b503session.TransportKey{}, b503session.ErrTransportDown
		},
	)
	manager.ResetForRestart(ConfiguredMonitorTarget)
	return manager
}

func newHistoryDispatcher(cfg Config, substrate Substrate, poller *vaillantSemanticPoller, manager *b503session.Manager, gate sync.Locker) *RawFrameDispatcher {
	var bus b503Bus
	if substrate.Bus != nil {
		bus = substrate.Bus
	}
	return NewRawFrameDispatcherWithSourceProvider(bus, func() (byte, bool) {
		source, _ := poller.sourceAndCompanion()
		return source, source != 0
	}, gate, manager, cfg.SemanticRequestTimeout)
}

// Runtime returns the lifecycle runtime shared by the owner and its dispatcher.
func (a *Acquisition) Runtime() *B503Runtime {
	if a == nil || a.dispatcher == nil {
		return nil
	}
	return a.dispatcher.runtime
}

// History returns the session manager tied to this owner.
func (a *Acquisition) History() *b503session.Manager {
	if a == nil {
		return nil
	}
	return a.manager
}

// Schedule returns a copy of the owner's recurring task intervals.
func (a *Acquisition) Schedule() map[string]time.Duration {
	if a == nil || a.poller == nil {
		return nil
	}
	schedule := a.poller.pollSchedule()
	if schedule == nil {
		return map[string]time.Duration{}
	}
	out := make(map[string]time.Duration, len(schedule.recurring))
	for _, task := range schedule.recurring {
		out[string(task.key)] = task.interval
	}
	return out
}

// Views returns the owner's source and published zone and hot-water view.
func (a *Acquisition) Views() AcquisitionView {
	if a == nil || a.poller == nil {
		return AcquisitionView{}
	}
	source, _ := a.poller.sourceAndCompanion()
	view := AcquisitionView{Source: source}
	published := publicationOwner(a.poller)
	if published == nil || published.state == nil {
		return view
	}
	view.Zones = len(published.state.Zones())
	view.HotWater = published.state.DHW() != nil
	return view
}

// SetSource rebinds the initiator the owner admits.
func (a *Acquisition) SetSource(source byte) {
	if a == nil || a.poller == nil {
		return
	}
	a.poller.SetSourceAndCompanion(source, 0)
}

// Source returns the initiator the owner admits.
func (a *Acquisition) Source() byte {
	if a == nil || a.poller == nil {
		return 0
	}
	source, _ := a.poller.sourceAndCompanion()
	return source
}

// Transport returns the substrate transport the owner was built with.
func (a *Acquisition) Transport() transport.RawTransport {
	if a == nil {
		return nil
	}
	return a.substrate.Transport
}

// LockDispatcher holds the poller exchange lock the dispatcher uses.
func (a *Acquisition) LockDispatcher() {
	if a != nil && a.gate != nil {
		a.gate.Lock()
	}
}

// UnlockDispatcher releases the poller exchange lock.
func (a *Acquisition) UnlockDispatcher() {
	if a != nil && a.gate != nil {
		a.gate.Unlock()
	}
}

// DispatcherWaiting reports how many callers are blocked on the dispatcher lock.
func (a *Acquisition) DispatcherWaiting() int {
	if a == nil || a.gate == nil {
		return 0
	}
	return int(a.gate.waiting.Load())
}

// Invoke sends one payload through the owner's dispatcher.
func (a *Acquisition) Invoke(ctx context.Context, target byte, payload []byte) ([]byte, error) {
	if a == nil || a.dispatcher == nil {
		return nil, errRawFrameMisconfigured
	}
	return a.dispatcher.Invoke(ctx, target, payload)
}

// HoldStartup keeps the owner's polling loops waiting. Call it before Begin.
func (a *Acquisition) HoldStartup() {
	if a == nil || a.poller == nil || a.poller.startupBarrier != nil {
		return
	}
	a.poller.startupBarrier = make(chan struct{})
}

// Begin starts the owner's schedule on ctx.
func (a *Acquisition) Begin(ctx context.Context) {
	if a == nil || a.poller == nil {
		return
	}
	a.poller.Start(ctx)
}

// Suspend stops acquisition at the next task boundary.
func (a *Acquisition) Suspend() {
	if a == nil || a.poller == nil {
		return
	}
	a.poller.SuspendAcquisition()
}

// RefreshCalls reports how often the history refresh stub has run.
func (a *Acquisition) RefreshCalls() int {
	if a == nil {
		return 0
	}
	return int(a.refreshCalls.Load())
}

// SourceUnavailable reports whether err is a source-admission failure.
func SourceUnavailable(err error) bool {
	return errors.Is(err, errRawFrameSourceNotAdmitted)
}
