package owner

// Mock bus, lock and context seams shared by the history dispatcher tests.

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/owner/b503session"
)

// adapterSource is the admitted eBUS source used by this test dispatcher.
// Production wiring receives the admitted source from startup admission.
const adapterSource byte = 0x7F

// b503DispatcherMockBus implements b503Bus with a programmable response
// table keyed on the first 2 bytes of the request payload (the §2 family/
// selector prefix). Mirrors the LOCAL_CAPTURE wire shape so the dispatcher
// path is exercised end-to-end with realistic byte streams.
type b503DispatcherMockBus struct {
	mu sync.Mutex

	respByPrefix map[string]*protocol.Frame
	errByPrefix  map[string]error
	defaultErr   error

	calls []protocol.Frame
	// onSend, if set, is invoked synchronously inside Send before
	// returning. Tests use it to inject context cancellation, transport
	// disconnects, or epoch rollovers mid-call.
	onSend func(frame protocol.Frame, mb *b503DispatcherMockBus)
	// blockUntil, if non-nil, makes Send block until the channel closes
	// or ctx fires. Used to model long-running bus turnaround.
	blockUntil <-chan struct{}
}

// b503BlockingReadLocker stops a dispatcher exactly where it has requested
// poll quiescence but has not yet reached the B503 emission boundary.
type b503BlockingReadLocker struct {
	entered chan struct{}
	release chan struct{}
	once    sync.Once
}

func newB503BlockingReadLocker() *b503BlockingReadLocker {
	return &b503BlockingReadLocker{entered: make(chan struct{}), release: make(chan struct{})}
}

func (l *b503BlockingReadLocker) Lock() {
	l.once.Do(func() { close(l.entered) })
	<-l.release
}

func (*b503BlockingReadLocker) Unlock() {}

// b503LaggingParentContext models a caller deadline whose propagation to a
// derived context has not yet run. context.WithTimeout uses AfterFunc when a
// parent provides it; this test seam deliberately holds that callback so the
// final pre-emission check must consult the caller context directly.
type b503LaggingParentContext struct {
	done chan struct{}

	mu        sync.Mutex
	err       error
	afterFunc func()
}

func newB503LaggingParentContext() *b503LaggingParentContext {
	return &b503LaggingParentContext{done: make(chan struct{})}
}

func (c *b503LaggingParentContext) Deadline() (time.Time, bool) { return time.Time{}, false }

func (c *b503LaggingParentContext) Done() <-chan struct{} { return c.done }

func (c *b503LaggingParentContext) Err() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.err
}

func (*b503LaggingParentContext) Value(any) any { return nil }

func (c *b503LaggingParentContext) AfterFunc(f func()) func() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.afterFunc = f
	return func() bool {
		c.mu.Lock()
		defer c.mu.Unlock()
		if c.afterFunc == nil {
			return false
		}
		c.afterFunc = nil
		return true
	}
}

func (c *b503LaggingParentContext) expire() {
	c.mu.Lock()
	c.err = context.DeadlineExceeded
	close(c.done)
	c.mu.Unlock()
}

// classifyBarrierError blocks the first Error call. bus.Send returns
// the value without formatting it, so reaching this barrier proves Invoke has
// already completed its post-Send epoch check and entered error classification.
type classifyBarrierError struct {
	once             sync.Once
	classification   chan struct{}
	continueClassify <-chan struct{}
}

func (err *classifyBarrierError) Error() string {
	err.once.Do(func() {
		close(err.classification)
		<-err.continueClassify
	})
	return "ebus: transport closed"
}

func newB503DispatcherMockBus() *b503DispatcherMockBus {
	return &b503DispatcherMockBus{
		respByPrefix: make(map[string]*protocol.Frame),
		errByPrefix:  make(map[string]error),
	}
}

func (mb *b503DispatcherMockBus) Send(ctx context.Context, frame protocol.Frame) (*protocol.Frame, error) {
	mb.mu.Lock()
	mb.calls = append(mb.calls, cloneFrame(frame))
	hook := mb.onSend
	block := mb.blockUntil
	mb.mu.Unlock()

	if hook != nil {
		hook(frame, mb)
	}

	if block != nil {
		select {
		case <-block:
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	mb.mu.Lock()
	defer mb.mu.Unlock()
	if len(frame.Data) >= 2 {
		key := string([]byte{frame.Data[0], frame.Data[1]})
		if e, ok := mb.errByPrefix[key]; ok {
			return nil, e
		}
		if r, ok := mb.respByPrefix[key]; ok {
			out := cloneFrame(*r)
			return &out, nil
		}
	}
	if mb.defaultErr != nil {
		return nil, mb.defaultErr
	}
	return nil, errors.New("b503DispatcherMockBus: no canned response for prefix")
}

func cloneFrame(f protocol.Frame) protocol.Frame {
	return protocol.Frame{
		Source:    f.Source,
		Target:    f.Target,
		Primary:   f.Primary,
		Secondary: f.Secondary,
		Data:      append([]byte(nil), f.Data...),
	}
}

func (mb *b503DispatcherMockBus) setResp(prefix [2]byte, data []byte) {
	mb.mu.Lock()
	defer mb.mu.Unlock()
	mb.respByPrefix[string(prefix[:])] = &protocol.Frame{Data: append([]byte(nil), data...)}
}

func (mb *b503DispatcherMockBus) setErr(prefix [2]byte, err error) {
	mb.mu.Lock()
	defer mb.mu.Unlock()
	mb.errByPrefix[string(prefix[:])] = err
}

func (mb *b503DispatcherMockBus) callCount() int {
	mb.mu.Lock()
	defer mb.mu.Unlock()
	return len(mb.calls)
}

func (mb *b503DispatcherMockBus) lastCall() (protocol.Frame, bool) {
	mb.mu.Lock()
	defer mb.mu.Unlock()
	if len(mb.calls) == 0 {
		return protocol.Frame{}, false
	}
	return mb.calls[len(mb.calls)-1], true
}

// newTestDispatcher builds a RawFrameDispatcher backed by a fresh mock
// bus and a fresh Manager.
func newTestDispatcher(t *testing.T) (*RawFrameDispatcher, *b503DispatcherMockBus, *b503session.Manager) {
	t.Helper()
	bus := newB503DispatcherMockBus()
	mgr := b503session.New(
		b503session.TransportKey{AdapterInstanceID: "test", TransportEpoch: 1},
		30*time.Second,
		func(ctx context.Context) (b503session.TransportKey, error) {
			return b503session.TransportKey{}, b503session.ErrTransportDown
		},
	)
	var readMu sync.Mutex
	disp := NewRawFrameDispatcher(bus, adapterSource, &readMu, mgr, 2*time.Second)
	return disp, bus, mgr
}

// hangGuard bounds a broken barrier without making elapsed time an ordering
// oracle. Half the test budget stays reserved for cancellation and joining.
func hangGuard(t *testing.T) <-chan time.Time {
	t.Helper()
	budget := time.Minute
	if deadline, ok := t.Deadline(); ok {
		budget = time.Until(deadline) / 2
	}
	timer := time.NewTimer(budget)
	t.Cleanup(func() { timer.Stop() })
	return timer.C
}
