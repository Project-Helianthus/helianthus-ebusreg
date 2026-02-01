package router

import (
	"context"
	"errors"
	"testing"

	ebuserrors "github.com/d3vi1/helianthus-ebusgo/errors"
	"github.com/d3vi1/helianthus-ebusgo/protocol"
	"github.com/d3vi1/helianthus-ebusreg/registry"
)

type mockPlane struct {
	name           string
	methods        []registry.Method
	subscriptions  []Subscription
	broadcastCalls int
	lastBroadcast  protocol.Frame
	broadcastErr   error
	buildRequest   func(registry.Method, map[string]any) (protocol.Frame, error)
	decodeResponse func(registry.Method, protocol.Frame) (any, error)
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

func (plane *mockPlane) DecodeResponse(method registry.Method, response protocol.Frame) (any, error) {
	return plane.decodeResponse(method, response)
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
		decodeResponse: func(method registry.Method, response protocol.Frame) (any, error) {
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
		decodeResponse: func(method registry.Method, response protocol.Frame) (any, error) {
			return nil, nil
		},
	}

	_, err := router.Invoke(context.Background(), plane, "missing", map[string]any{})
	if !errors.Is(err, ebuserrors.ErrInvalidPayload) {
		t.Fatalf("expected ErrInvalidPayload, got %v", err)
	}
}
