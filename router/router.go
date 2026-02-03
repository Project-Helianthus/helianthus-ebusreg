package router

import (
	"context"
	"fmt"
	"sync"

	ebuserrors "github.com/d3vi1/helianthus-ebusgo/errors"
	"github.com/d3vi1/helianthus-ebusgo/protocol"
	"github.com/d3vi1/helianthus-ebusgo/types"
	"github.com/d3vi1/helianthus-ebusreg/registry"
)

type Subscription struct {
	Primary   byte
	Secondary byte
}

type BroadcastEvent struct {
	Plane  string
	Frame  protocol.Frame
	Values map[string]types.Value
}

type BroadcastDecoder interface {
	DecodeBroadcast(frame protocol.Frame) (map[string]types.Value, bool, error)
}

type Plane interface {
	Name() string
	Methods() []registry.Method
	Subscriptions() []Subscription
	OnBroadcast(frame protocol.Frame) error
	BuildRequest(method registry.Method, params map[string]any) (protocol.Frame, error)
	DecodeResponse(method registry.Method, response protocol.Frame, params map[string]any) (any, error)
}

type Bus interface {
	Send(ctx context.Context, frame protocol.Frame) (*protocol.Frame, error)
}

type BusEventRouter struct {
	bus           Bus
	mu            sync.RWMutex
	subscriptions map[subscriptionKey][]Plane
	events        chan BroadcastEvent
}

type subscriptionKey struct {
	primary   byte
	secondary byte
}

func NewBusEventRouter(bus Bus) *BusEventRouter {
	return &BusEventRouter{
		bus:           bus,
		subscriptions: make(map[subscriptionKey][]Plane),
		events:        make(chan BroadcastEvent, 64),
	}
}

func (router *BusEventRouter) Events() <-chan BroadcastEvent {
	if router == nil {
		return nil
	}
	return router.events
}

func (router *BusEventRouter) SetPlanes(planes []Plane) {
	subscriptions := make(map[subscriptionKey][]Plane)
	for _, plane := range planes {
		for _, sub := range plane.Subscriptions() {
			key := subscriptionKey{primary: sub.Primary, secondary: sub.Secondary}
			subscriptions[key] = append(subscriptions[key], plane)
		}
	}

	router.mu.Lock()
	router.subscriptions = subscriptions
	router.mu.Unlock()
}

func (router *BusEventRouter) HandleBroadcast(frame protocol.Frame) []error {
	key := subscriptionKey{primary: frame.Primary, secondary: frame.Secondary}

	router.mu.RLock()
	planes := append([]Plane(nil), router.subscriptions[key]...)
	router.mu.RUnlock()

	if len(planes) == 0 {
		return nil
	}

	errors := make([]error, 0)
	for _, plane := range planes {
		if err := plane.OnBroadcast(frame); err != nil {
			errors = append(errors, err)
		}

		decoder, ok := plane.(BroadcastDecoder)
		if !ok {
			continue
		}
		decoded, handled, err := decoder.DecodeBroadcast(frame)
		if err != nil {
			errors = append(errors, err)
			continue
		}
		if !handled {
			continue
		}

		event := BroadcastEvent{
			Plane:  plane.Name(),
			Frame:  frame,
			Values: decoded,
		}
		select {
		case router.events <- event:
		default:
		}
	}
	return errors
}

func (router *BusEventRouter) Invoke(ctx context.Context, plane Plane, methodName string, params map[string]any) (any, error) {
	if plane == nil {
		return nil, fmt.Errorf("router.Invoke missing plane: %w", ebuserrors.ErrInvalidPayload)
	}

	method, ok := findMethod(plane.Methods(), methodName)
	if !ok {
		return nil, fmt.Errorf("router.Invoke missing method %q: %w", methodName, ebuserrors.ErrInvalidPayload)
	}

	request, err := plane.BuildRequest(method, params)
	if err != nil {
		return nil, fmt.Errorf("router.Invoke build plane=%s method=%s: %w", plane.Name(), method.Name(), err)
	}

	response, err := router.bus.Send(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("router.Invoke send plane=%s method=%s: %w", plane.Name(), method.Name(), err)
	}
	if response == nil {
		return nil, fmt.Errorf("router.Invoke empty response plane=%s method=%s: %w", plane.Name(), method.Name(), ebuserrors.ErrInvalidPayload)
	}

	decoded, err := plane.DecodeResponse(method, *response, params)
	if err != nil {
		return nil, fmt.Errorf("router.Invoke decode plane=%s method=%s: %w", plane.Name(), method.Name(), err)
	}
	return decoded, nil
}

func findMethod(methods []registry.Method, name string) (registry.Method, bool) {
	for _, method := range methods {
		if method.Name() == name {
			return method, true
		}
	}
	return nil, false
}
