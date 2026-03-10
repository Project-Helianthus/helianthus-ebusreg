package router

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	ebuserrors "github.com/Project-Helianthus/helianthus-ebusgo/errors"
	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
	"github.com/Project-Helianthus/helianthus-ebusgo/types"
	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
	"github.com/Project-Helianthus/helianthus-ebusreg/schema"
)

type mockPlane struct {
	name           string
	methods        []registry.Method
	subscriptions  []Subscription
	broadcastCalls int
	lastBroadcast  protocol.Frame
	broadcastErr   error
	buildRequest   func(registry.Method, map[string]any) (protocol.Frame, error)
	decodeResponse func(registry.Method, protocol.Frame, map[string]any) (any, error)
}

type mockTemplate struct {
	primary   byte
	secondary byte
}

func (template mockTemplate) Primary() byte {
	return template.primary
}

func (template mockTemplate) Secondary() byte {
	return template.secondary
}

type mockMethod struct {
	name     string
	readOnly bool
	template registry.FrameTemplate
	response schema.SchemaSelector
}

func (method mockMethod) Name() string {
	return method.name
}

func (method mockMethod) ReadOnly() bool {
	return method.readOnly
}

func (method mockMethod) Template() registry.FrameTemplate {
	return method.template
}

func (method mockMethod) ResponseSchema() schema.SchemaSelector {
	return method.response
}

func (plane *mockPlane) Name() string {
	return plane.name
}

func (plane *mockPlane) Methods() []registry.Method {
	return plane.methods
}

func (plane *mockPlane) Subscriptions() []Subscription {
	return plane.subscriptions
}

func (plane *mockPlane) OnBroadcast(frame protocol.Frame) error {
	plane.broadcastCalls++
	plane.lastBroadcast = frame
	return plane.broadcastErr
}

func (plane *mockPlane) BuildRequest(method registry.Method, params map[string]any) (protocol.Frame, error) {
	return plane.buildRequest(method, params)
}

func (plane *mockPlane) DecodeResponse(method registry.Method, response protocol.Frame, params map[string]any) (any, error) {
	return plane.decodeResponse(method, response, params)
}

type mockBus struct {
	lastRequest protocol.Frame
	response    *protocol.Frame
	err         error
	callCount   int
}

func (bus *mockBus) Send(ctx context.Context, frame protocol.Frame) (*protocol.Frame, error) {
	bus.callCount++
	bus.lastRequest = frame
	return bus.response, bus.err
}

type coalescingBus struct {
	mu        sync.Mutex
	response  protocol.Frame
	delay     time.Duration
	callCount int
}

func (bus *coalescingBus) Send(ctx context.Context, frame protocol.Frame) (*protocol.Frame, error) {
	if bus.delay > 0 {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(bus.delay):
		}
	}

	bus.mu.Lock()
	bus.callCount++
	response := cloneFrame(bus.response)
	bus.mu.Unlock()
	return &response, nil
}

func (bus *coalescingBus) Calls() int {
	bus.mu.Lock()
	defer bus.mu.Unlock()
	return bus.callCount
}

func TestBusEventRouter_BroadcastFanOut(t *testing.T) {
	t.Parallel()

	router := NewBusEventRouter(&mockBus{})

	planeA := &mockPlane{
		name: "A",
		subscriptions: []Subscription{
			{Primary: 0xB5, Secondary: 0x16},
		},
	}
	planeB := &mockPlane{
		name: "B",
		subscriptions: []Subscription{
			{Primary: 0xB5, Secondary: 0x04},
		},
	}
	planeC := &mockPlane{
		name: "C",
		subscriptions: []Subscription{
			{Primary: 0xB5, Secondary: 0x16},
		},
	}

	router.SetPlanes([]Plane{planeA, planeB, planeC})

	errors := router.HandleBroadcast(protocol.Frame{Primary: 0xB5, Secondary: 0x16})
	if len(errors) != 0 {
		t.Fatalf("unexpected errors: %v", errors)
	}
	if planeA.broadcastCalls != 1 || planeC.broadcastCalls != 1 {
		t.Fatalf("expected plane A and C to receive broadcast")
	}
	if planeB.broadcastCalls != 0 {
		t.Fatalf("expected plane B to not receive broadcast")
	}
}

func TestBusEventRouter_BroadcastCollectsErrors(t *testing.T) {
	t.Parallel()

	router := NewBusEventRouter(&mockBus{})
	plane := &mockPlane{
		name: "A",
		subscriptions: []Subscription{
			{Primary: 0xB5, Secondary: 0x16},
		},
		broadcastErr: errors.New("boom"),
	}

	router.SetPlanes([]Plane{plane})
	errors := router.HandleBroadcast(protocol.Frame{Primary: 0xB5, Secondary: 0x16})
	if len(errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errors))
	}
}

func TestBusEventRouter_InvokeRoutesResponse(t *testing.T) {
	t.Parallel()

	method := mockMethod{
		name:     "get_status",
		readOnly: true,
		template: mockTemplate{primary: 0xB5, secondary: 0x04},
	}

	plane := &mockPlane{
		name:    "heating",
		methods: []registry.Method{method},
		buildRequest: func(method registry.Method, params map[string]any) (protocol.Frame, error) {
			return protocol.Frame{Source: 0x10, Target: 0x08, Primary: 0xB5, Secondary: 0x04}, nil
		},
		decodeResponse: func(method registry.Method, response protocol.Frame, params map[string]any) (any, error) {
			return "ok", nil
		},
	}

	response := &protocol.Frame{Source: 0x08, Target: 0x10, Primary: 0xB5, Secondary: 0x04}
	bus := &mockBus{response: response}
	router := NewBusEventRouter(bus)

	result, err := router.Invoke(context.Background(), plane, "get_status", map[string]any{})
	if err != nil {
		t.Fatalf("Invoke error = %v", err)
	}
	if result != "ok" {
		t.Fatalf("Invoke result = %v; want ok", result)
	}
	if bus.callCount != 1 {
		t.Fatalf("expected Send to be called once")
	}
}

func TestBusEventRouter_InvokeMissingMethod(t *testing.T) {
	t.Parallel()

	router := NewBusEventRouter(&mockBus{})
	plane := &mockPlane{
		name:    "heating",
		methods: []registry.Method{},
		buildRequest: func(method registry.Method, params map[string]any) (protocol.Frame, error) {
			return protocol.Frame{}, nil
		},
		decodeResponse: func(method registry.Method, response protocol.Frame, params map[string]any) (any, error) {
			return nil, nil
		},
	}

	_, err := router.Invoke(context.Background(), plane, "missing", map[string]any{})
	if !errors.Is(err, ebuserrors.ErrInvalidPayload) {
		t.Fatalf("expected ErrInvalidPayload, got %v", err)
	}
}

type decodingPlane struct {
	*mockPlane
	decoded map[string]types.Value
	handled bool
	err     error
}

func (plane *decodingPlane) DecodeBroadcast(frame protocol.Frame) (map[string]types.Value, bool, error) {
	return plane.decoded, plane.handled, plane.err
}

func TestBusEventRouter_EmitsDecodedBroadcastEvents(t *testing.T) {
	t.Parallel()

	router := NewBusEventRouter(&mockBus{})
	plane := &decodingPlane{
		mockPlane: &mockPlane{
			name: "A",
			subscriptions: []Subscription{
				{Primary: 0xB5, Secondary: 0x16},
			},
		},
		decoded: map[string]types.Value{
			"foo": {Value: uint8(1), Valid: true},
		},
		handled: true,
	}
	router.SetPlanes([]Plane{plane})

	errors := router.HandleBroadcast(protocol.Frame{Primary: 0xB5, Secondary: 0x16})
	if len(errors) != 0 {
		t.Fatalf("unexpected errors: %v", errors)
	}

	select {
	case event := <-router.Events():
		if event.Plane != "A" {
			t.Fatalf("event.Plane = %q; want A", event.Plane)
		}
		value, ok := event.Values["foo"]
		if !ok || !value.Valid || value.Value != uint8(1) {
			t.Fatalf("event.Values[foo] = %+v; want valid 1", value)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("timeout waiting for broadcast event")
	}
}

func TestBusEventRouter_SurfacesBroadcastEventOverflow(t *testing.T) {
	t.Parallel()

	before := routerBroadcastEventOverflowTotal.Value()
	router := NewBusEventRouter(&mockBus{})
	plane := &decodingPlane{
		mockPlane: &mockPlane{
			name: "A",
			subscriptions: []Subscription{
				{Primary: 0xB5, Secondary: 0x16},
			},
		},
		decoded: map[string]types.Value{
			"foo": {Value: uint8(1), Valid: true},
		},
		handled: true,
	}
	router.SetPlanes([]Plane{plane})

	for i := 0; i < cap(router.events); i++ {
		router.events <- BroadcastEvent{Plane: "seed"}
	}

	errs := router.HandleBroadcast(protocol.Frame{Primary: 0xB5, Secondary: 0x16})
	if len(errs) != 1 {
		t.Fatalf("errors = %d; want 1", len(errs))
	}
	if !errors.Is(errs[0], ErrBroadcastEventOverflow) {
		t.Fatalf("error = %v; want ErrBroadcastEventOverflow", errs[0])
	}
	if plane.broadcastCalls != 1 {
		t.Fatalf("broadcast calls = %d; want 1", plane.broadcastCalls)
	}
	if got := routerBroadcastEventOverflowTotal.Value(); got != before+1 {
		t.Fatalf("routerBroadcastEventOverflowTotal = %d; want %d", got, before+1)
	}
}

func TestBusEventRouter_SurfacesBroadcastEventOverflowPerPlane(t *testing.T) {
	t.Parallel()

	before := routerBroadcastEventOverflowTotal.Value()
	router := NewBusEventRouter(&mockBus{})
	planeA := &decodingPlane{
		mockPlane: &mockPlane{
			name: "A",
			subscriptions: []Subscription{
				{Primary: 0xB5, Secondary: 0x16},
			},
		},
		decoded: map[string]types.Value{
			"foo": {Value: uint8(1), Valid: true},
		},
		handled: true,
	}
	planeB := &decodingPlane{
		mockPlane: &mockPlane{
			name: "B",
			subscriptions: []Subscription{
				{Primary: 0xB5, Secondary: 0x16},
			},
		},
		decoded: map[string]types.Value{
			"bar": {Value: uint8(2), Valid: true},
		},
		handled: true,
	}
	router.SetPlanes([]Plane{planeA, planeB})

	for i := 0; i < cap(router.events); i++ {
		router.events <- BroadcastEvent{Plane: "seed"}
	}

	errs := router.HandleBroadcast(protocol.Frame{Primary: 0xB5, Secondary: 0x16})
	if len(errs) != 2 {
		t.Fatalf("errors = %d; want 2", len(errs))
	}
	for _, err := range errs {
		if !errors.Is(err, ErrBroadcastEventOverflow) {
			t.Fatalf("error = %v; want ErrBroadcastEventOverflow", err)
		}
	}
	if planeA.broadcastCalls != 1 || planeB.broadcastCalls != 1 {
		t.Fatalf("broadcast calls = (%d, %d); want (1, 1)", planeA.broadcastCalls, planeB.broadcastCalls)
	}
	if got := routerBroadcastEventOverflowTotal.Value(); got != before+2 {
		t.Fatalf("routerBroadcastEventOverflowTotal = %d; want %d", got, before+2)
	}
}

func TestBusEventRouter_InvokeCoalescesConcurrentReadOnlyCalls(t *testing.T) {
	t.Parallel()

	method := mockMethod{
		name:     "read_state",
		readOnly: true,
		template: mockTemplate{primary: 0xB5, secondary: 0x24},
	}

	plane := &mockPlane{
		name:    "system",
		methods: []registry.Method{method},
		buildRequest: func(method registry.Method, params map[string]any) (protocol.Frame, error) {
			return protocol.Frame{Source: 0x31, Target: 0x15, Primary: 0xB5, Secondary: 0x24}, nil
		},
		decodeResponse: func(method registry.Method, response protocol.Frame, params map[string]any) (any, error) {
			return map[string]any{"ok": true}, nil
		},
	}

	bus := &coalescingBus{
		response: protocol.Frame{Source: 0x15, Target: 0x31, Primary: 0xB5, Secondary: 0x24, Data: []byte{0x01}},
		delay:    40 * time.Millisecond,
	}
	eventRouter := NewBusEventRouter(bus)

	const workers = 6
	var wg sync.WaitGroup
	errCh := make(chan error, workers)
	for worker := 0; worker < workers; worker++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := eventRouter.Invoke(context.Background(), plane, "read_state", map[string]any{"addr": 0x10})
			errCh <- err
		}()
	}
	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			t.Fatalf("Invoke error = %v", err)
		}
	}
	if got := bus.Calls(); got != 1 {
		t.Fatalf("bus calls = %d; want 1", got)
	}
}

func TestBusEventRouter_InvokeUsesRefreshPolicyWindows(t *testing.T) {
	t.Parallel()

	methodState := mockMethod{
		name:     "read_state",
		readOnly: true,
		template: mockTemplate{primary: 0xB5, secondary: 0x24},
	}
	methodConfig := mockMethod{
		name:     "read_config",
		readOnly: true,
		template: mockTemplate{primary: 0xB5, secondary: 0x24},
	}
	plane := &mockPlane{
		name:    "system",
		methods: []registry.Method{methodState, methodConfig},
		buildRequest: func(method registry.Method, params map[string]any) (protocol.Frame, error) {
			return protocol.Frame{Source: 0x31, Target: 0x15, Primary: 0xB5, Secondary: 0x24}, nil
		},
		decodeResponse: func(method registry.Method, response protocol.Frame, params map[string]any) (any, error) {
			return map[string]any{"name": method.Name()}, nil
		},
	}

	bus := &coalescingBus{
		response: protocol.Frame{Source: 0x15, Target: 0x31, Primary: 0xB5, Secondary: 0x24, Data: []byte{0x01}},
	}

	eventRouter := NewBusEventRouterWithOptions(bus, RouterOptions{
		StateRefreshInterval:  30 * time.Millisecond,
		ConfigRefreshInterval: 90 * time.Millisecond,
	})

	_, err := eventRouter.Invoke(context.Background(), plane, "read_state", map[string]any{"addr": 0x10})
	if err != nil {
		t.Fatalf("Invoke state #1 error = %v", err)
	}
	_, err = eventRouter.Invoke(context.Background(), plane, "read_state", map[string]any{"addr": 0x10})
	if err != nil {
		t.Fatalf("Invoke state #2 error = %v", err)
	}
	if got := bus.Calls(); got != 1 {
		t.Fatalf("state warm cache calls = %d; want 1", got)
	}

	time.Sleep(35 * time.Millisecond)
	_, err = eventRouter.Invoke(context.Background(), plane, "read_state", map[string]any{"addr": 0x10})
	if err != nil {
		t.Fatalf("Invoke state #3 error = %v", err)
	}
	if got := bus.Calls(); got != 2 {
		t.Fatalf("state after ttl calls = %d; want 2", got)
	}

	_, err = eventRouter.Invoke(context.Background(), plane, "read_config", map[string]any{"addr": 0x11})
	if err != nil {
		t.Fatalf("Invoke config #1 error = %v", err)
	}
	time.Sleep(35 * time.Millisecond)
	_, err = eventRouter.Invoke(context.Background(), plane, "read_config", map[string]any{"addr": 0x11})
	if err != nil {
		t.Fatalf("Invoke config #2 error = %v", err)
	}
	if got := bus.Calls(); got != 3 {
		t.Fatalf("config still cached calls = %d; want 3", got)
	}
	time.Sleep(65 * time.Millisecond)
	_, err = eventRouter.Invoke(context.Background(), plane, "read_config", map[string]any{"addr": 0x11})
	if err != nil {
		t.Fatalf("Invoke config #3 error = %v", err)
	}
	if got := bus.Calls(); got != 4 {
		t.Fatalf("config after ttl calls = %d; want 4", got)
	}
}

func TestBusEventRouter_InvokeRefreshClassParamOverridesMethodHeuristic(t *testing.T) {
	t.Parallel()

	method := mockMethod{
		name:     "read_state",
		readOnly: true,
		template: mockTemplate{primary: 0xB5, secondary: 0x24},
	}
	plane := &mockPlane{
		name:    "system",
		methods: []registry.Method{method},
		buildRequest: func(method registry.Method, params map[string]any) (protocol.Frame, error) {
			return protocol.Frame{Source: 0x31, Target: 0x15, Primary: 0xB5, Secondary: 0x24}, nil
		},
		decodeResponse: func(method registry.Method, response protocol.Frame, params map[string]any) (any, error) {
			return map[string]any{"name": method.Name()}, nil
		},
	}
	bus := &coalescingBus{
		response: protocol.Frame{Source: 0x15, Target: 0x31, Primary: 0xB5, Secondary: 0x24, Data: []byte{0x01}},
	}

	eventRouter := NewBusEventRouterWithOptions(bus, RouterOptions{
		StateRefreshInterval:  25 * time.Millisecond,
		ConfigRefreshInterval: 80 * time.Millisecond,
	})

	params := map[string]any{"addr": 0x12, "cache_class": "config"}
	_, err := eventRouter.Invoke(context.Background(), plane, "read_state", params)
	if err != nil {
		t.Fatalf("Invoke #1 error = %v", err)
	}
	time.Sleep(30 * time.Millisecond)
	_, err = eventRouter.Invoke(context.Background(), plane, "read_state", params)
	if err != nil {
		t.Fatalf("Invoke #2 error = %v", err)
	}
	if got := bus.Calls(); got != 1 {
		t.Fatalf("calls = %d; want 1 with config override", got)
	}
}

func TestBusEventRouter_InvokeWriteMethodsDoNotUseCache(t *testing.T) {
	t.Parallel()

	method := mockMethod{
		name:     "set_register",
		readOnly: false,
		template: mockTemplate{primary: 0xB5, secondary: 0x24},
	}
	plane := &mockPlane{
		name:    "system",
		methods: []registry.Method{method},
		buildRequest: func(method registry.Method, params map[string]any) (protocol.Frame, error) {
			return protocol.Frame{Source: 0x31, Target: 0x15, Primary: 0xB5, Secondary: 0x24}, nil
		},
		decodeResponse: func(method registry.Method, response protocol.Frame, params map[string]any) (any, error) {
			return map[string]any{"ok": true}, nil
		},
	}
	bus := &coalescingBus{
		response: protocol.Frame{Source: 0x15, Target: 0x31, Primary: 0xB5, Secondary: 0x24, Data: []byte{0x01}},
	}
	eventRouter := NewBusEventRouterWithOptions(bus, RouterOptions{
		StateRefreshInterval:  1 * time.Minute,
		ConfigRefreshInterval: 5 * time.Minute,
	})

	_, err := eventRouter.Invoke(context.Background(), plane, "set_register", map[string]any{"addr": 0x20})
	if err != nil {
		t.Fatalf("Invoke #1 error = %v", err)
	}
	_, err = eventRouter.Invoke(context.Background(), plane, "set_register", map[string]any{"addr": 0x20})
	if err != nil {
		t.Fatalf("Invoke #2 error = %v", err)
	}
	if got := bus.Calls(); got != 2 {
		t.Fatalf("calls = %d; want 2 for write method", got)
	}
}
