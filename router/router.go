package router

import (
	"context"
	"encoding/json"
	"errors"
	"expvar"
	"fmt"
	"strings"
	"sync"
	"time"

	ebuserrors "github.com/Project-Helianthus/helianthus-ebusgo/errors"
	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
	"github.com/Project-Helianthus/helianthus-ebusgo/types"
	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
)

// Subscription selects broadcast frames by their primary and secondary
// command bytes.
type Subscription struct {
	Primary   byte
	Secondary byte
}

// ErrBroadcastEventOverflow is reported by HandleBroadcast when the events
// channel is full. The decoded event is dropped, the plane's own handling
// has already run, and the overflow counter is incremented.
var ErrBroadcastEventOverflow = errors.New("router broadcast event queue full")

var routerBroadcastEventOverflowTotal = expvar.NewInt("router_broadcast_event_overflow_total")

// BroadcastEvent is one decoded broadcast delivered on Events: the plane
// that decoded it, the frame and the decoded values by name.
type BroadcastEvent struct {
	Plane  string
	Frame  protocol.Frame
	Values map[string]types.Value
}

// BroadcastDecoder is an optional interface of a Plane. DecodeBroadcast
// returns handled false for frames it does not turn into an event; an error
// is reported by HandleBroadcast and no event is emitted.
type BroadcastDecoder interface {
	DecodeBroadcast(frame protocol.Frame) (map[string]types.Value, bool, error)
}

// Plane is a registry plane that the router can drive: it names the
// broadcasts it subscribes to, receives them through OnBroadcast, and builds
// requests and decodes responses for its methods. DecodeResponse receives
// a private copy of the response frame on every call.
type Plane interface {
	Name() string
	Methods() []registry.Method
	Subscriptions() []Subscription
	OnBroadcast(frame protocol.Frame) error
	BuildRequest(method registry.Method, params map[string]any) (protocol.Frame, error)
	DecodeResponse(method registry.Method, response protocol.Frame, params map[string]any) (any, error)
}

// Bus sends one request frame and returns its response; a nil response
// without an error is treated as an invalid payload.
type Bus interface {
	Send(ctx context.Context, frame protocol.Frame) (*protocol.Frame, error)
}

// BusEventRouter fans broadcasts out to subscribed planes and invokes plane
// methods over a Bus. Concurrent read-only invocations with the same plane,
// method and parameters share one bus request, and a successful response is
// reused for the refresh interval of its class. Methods are safe for
// concurrent use; create a router with NewBusEventRouter or
// NewBusEventRouterWithOptions.
type BusEventRouter struct {
	bus           Bus
	mu            sync.RWMutex
	subscriptions map[subscriptionKey][]Plane
	events        chan BroadcastEvent

	invokeMu    sync.Mutex
	invokeCache map[string]*invokeEntry
	options     RouterOptions
	now         func() time.Time
}

type subscriptionKey struct {
	primary   byte
	secondary byte
}

// RefreshClass selects how long a read-only invocation result may be reused.
type RefreshClass string

// RefreshClass values: state data uses StateRefreshInterval and
// configuration data uses ConfigRefreshInterval.
const (
	RefreshClassState  RefreshClass = "state"
	RefreshClassConfig RefreshClass = "config"
)

const (
	defaultStateRefreshInterval  = 1 * time.Minute
	defaultConfigRefreshInterval = 5 * time.Minute
)

// ClassifyInvocationFunc overrides the refresh class that the router infers
// from a method name. Returning "" keeps the inferred class; a cache_class
// or refresh_class parameter of "state" or "config" still takes precedence.
type ClassifyInvocationFunc func(
	planeName string,
	methodName string,
	params map[string]any,
) RefreshClass

// RouterOptions configures invocation coalescing. Non-positive intervals
// select the defaults of one minute for state and five minutes for
// configuration. DisableCoalescing sends every invocation to the bus.
type RouterOptions struct {
	StateRefreshInterval  time.Duration
	ConfigRefreshInterval time.Duration
	ClassifyInvocation    ClassifyInvocationFunc
	DisableCoalescing     bool
}

type invokeEntry struct {
	running bool
	done    chan struct{}

	lastOKAt     time.Time
	lastFrame    protocol.Frame
	hasLastFrame bool
}

// NewBusEventRouter returns a router over bus with default options.
func NewBusEventRouter(bus Bus) *BusEventRouter {
	return NewBusEventRouterWithOptions(bus, RouterOptions{})
}

// NewBusEventRouterWithOptions returns a router over bus with options,
// replacing non-positive refresh intervals by their defaults. The events
// channel buffers 64 events.
func NewBusEventRouterWithOptions(bus Bus, options RouterOptions) *BusEventRouter {
	if options.StateRefreshInterval <= 0 {
		options.StateRefreshInterval = defaultStateRefreshInterval
	}
	if options.ConfigRefreshInterval <= 0 {
		options.ConfigRefreshInterval = defaultConfigRefreshInterval
	}

	return &BusEventRouter{
		bus:           bus,
		subscriptions: make(map[subscriptionKey][]Plane),
		events:        make(chan BroadcastEvent, 64),
		invokeCache:   make(map[string]*invokeEntry),
		options:       options,
		now:           time.Now,
	}
}

// Events returns the channel of decoded broadcast events, or nil for a nil
// router. The channel is never closed; a consumer that falls behind loses
// events with ErrBroadcastEventOverflow.
func (router *BusEventRouter) Events() <-chan BroadcastEvent {
	if router == nil {
		return nil
	}
	return router.events
}

// SetPlanes replaces all broadcast subscriptions with those of planes. A
// plane subscribed to the same command bytes more than once receives such
// a frame once per subscription.
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

// HandleBroadcast passes frame to every plane subscribed to its command
// bytes, in SetPlanes order, and queues a BroadcastEvent for each plane that
// decodes it. Queueing never blocks; every plane error, decode error and
// queue overflow is returned, and an empty result means none occurred.
func (router *BusEventRouter) HandleBroadcast(frame protocol.Frame) []error {
	key := subscriptionKey{primary: frame.Primary, secondary: frame.Secondary}

	router.mu.RLock()
	planes := append([]Plane(nil), router.subscriptions[key]...)
	router.mu.RUnlock()

	if len(planes) == 0 {
		return nil
	}

	errs := make([]error, 0)
	for _, plane := range planes {
		if err := plane.OnBroadcast(frame); err != nil {
			errs = append(errs, err)
		}

		decoder, ok := plane.(BroadcastDecoder)
		if !ok {
			continue
		}
		decoded, handled, err := decoder.DecodeBroadcast(frame)
		if err != nil {
			errs = append(errs, err)
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
			routerBroadcastEventOverflowTotal.Add(1)
			errs = append(errs, fmt.Errorf("router.HandleBroadcast plane=%s: %w", plane.Name(), ErrBroadcastEventOverflow))
		}
	}
	return errs
}

// Invoke builds the request of the named method of plane, sends it on the
// bus and returns the decoded response. Mutating methods are always sent.
// Read-only methods are coalesced: callers with the same plane name, method
// name and JSON-encoded params wait for one in-flight request, and a
// successful response is decoded again for later callers until the refresh
// interval of its class has passed. Failures are not cached. A nil plane or
// unknown method wraps the ebusgo ErrInvalidPayload; a cancelled ctx stops
// waiting for another caller's request.
func (router *BusEventRouter) Invoke(ctx context.Context, plane Plane, methodName string, params map[string]any) (any, error) {
	if plane == nil {
		return nil, fmt.Errorf("router.Invoke missing plane: %w", ebuserrors.ErrInvalidPayload)
	}
	if ctx == nil {
		ctx = context.Background()
	}

	method, ok := findMethod(plane.Methods(), methodName)
	if !ok {
		return nil, fmt.Errorf("router.Invoke missing method %q: %w", methodName, ebuserrors.ErrInvalidPayload)
	}

	request, err := plane.BuildRequest(method, params)
	if err != nil {
		return nil, fmt.Errorf("router.Invoke build plane=%s method=%s: %w", plane.Name(), method.Name(), err)
	}

	maxAge := router.invokeMaxAge(plane.Name(), method, params)
	if maxAge <= 0 || request.Target == protocol.AddressBroadcast {
		// A broadcast has no response frame to cache or replay.
		decoded, _, err := router.sendAndDecode(ctx, plane, method, params, request)
		return decoded, err
	}

	cacheKey, err := buildInvokeCacheKey(plane.Name(), method.Name(), params)
	if err != nil {
		return nil, fmt.Errorf("router.Invoke key plane=%s method=%s: %w", plane.Name(), method.Name(), err)
	}

	for {
		router.invokeMu.Lock()
		entry := router.invokeCache[cacheKey]
		if entry == nil {
			entry = &invokeEntry{}
			router.invokeCache[cacheKey] = entry
		}

		if !entry.running && entry.hasLastFrame {
			age := router.now().Sub(entry.lastOKAt)
			if age >= 0 && age <= maxAge {
				frame := cloneFrame(entry.lastFrame)
				router.invokeMu.Unlock()
				decoded, decodeErr := plane.DecodeResponse(method, frame, params)
				if decodeErr != nil {
					return nil, fmt.Errorf(
						"router.Invoke decode plane=%s method=%s: %w",
						plane.Name(),
						method.Name(),
						decodeErr,
					)
				}
				return decoded, nil
			}
		}

		if entry.running {
			done := entry.done
			router.invokeMu.Unlock()
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-done:
				continue
			}
		}

		entry.running = true
		entry.done = make(chan struct{})
		done := entry.done
		router.invokeMu.Unlock()

		decoded, frame, invokeErr := router.sendAndDecode(ctx, plane, method, params, request)

		router.invokeMu.Lock()
		entry.running = false
		if invokeErr == nil {
			entry.lastOKAt = router.now()
			entry.lastFrame = cloneFrame(frame)
			entry.hasLastFrame = true
		}
		close(done)
		router.invokeMu.Unlock()

		if invokeErr != nil {
			return nil, invokeErr
		}
		return decoded, nil
	}
}

func (router *BusEventRouter) sendAndDecode(
	ctx context.Context,
	plane Plane,
	method registry.Method,
	params map[string]any,
	request protocol.Frame,
) (any, protocol.Frame, error) {
	response, err := router.bus.Send(ctx, request)
	if err != nil {
		return nil, protocol.Frame{}, fmt.Errorf(
			"router.Invoke send plane=%s method=%s: %w",
			plane.Name(),
			method.Name(),
			err,
		)
	}
	if response == nil && request.Target == protocol.AddressBroadcast {
		// A broadcast has no response: the transport reports a sent
		// broadcast as (nil, nil), so there is nothing to decode.
		return nil, protocol.Frame{}, nil
	}
	if response == nil {
		return nil, protocol.Frame{}, fmt.Errorf(
			"router.Invoke empty response plane=%s method=%s: %w",
			plane.Name(),
			method.Name(),
			ebuserrors.ErrInvalidPayload,
		)
	}
	frame := cloneFrame(*response)

	decoded, err := plane.DecodeResponse(method, frame, params)
	if err != nil {
		return nil, protocol.Frame{}, fmt.Errorf(
			"router.Invoke decode plane=%s method=%s: %w",
			plane.Name(),
			method.Name(),
			err,
		)
	}
	return decoded, frame, nil
}

func (router *BusEventRouter) invokeMaxAge(
	planeName string,
	method registry.Method,
	params map[string]any,
) time.Duration {
	if router == nil || router.options.DisableCoalescing || method == nil || !method.ReadOnly() {
		return 0
	}

	class := inferRefreshClass(method.Name())
	if router.options.ClassifyInvocation != nil {
		override := router.options.ClassifyInvocation(planeName, method.Name(), params)
		if override != "" {
			class = override
		}
	}
	if paramClass, ok := refreshClassFromParams(params); ok {
		class = paramClass
	}

	if class == RefreshClassConfig {
		return router.options.ConfigRefreshInterval
	}
	return router.options.StateRefreshInterval
}

func inferRefreshClass(methodName string) RefreshClass {
	normalized := strings.ToLower(strings.TrimSpace(methodName))
	if normalized == "" {
		return RefreshClassState
	}

	configTokens := []string{
		"config",
		"setpoint",
		"desired",
		"mode",
		"curve",
		"holiday",
		"name",
		"limit",
		"offset",
		"schedule",
		"timer",
		"mapping",
	}
	for _, token := range configTokens {
		if strings.Contains(normalized, token) {
			return RefreshClassConfig
		}
	}
	return RefreshClassState
}

func refreshClassFromParams(params map[string]any) (RefreshClass, bool) {
	if len(params) == 0 {
		return "", false
	}
	for _, key := range []string{"cache_class", "refresh_class"} {
		value, ok := params[key]
		if !ok {
			continue
		}
		text, ok := value.(string)
		if !ok {
			continue
		}
		normalized := strings.ToLower(strings.TrimSpace(text))
		switch normalized {
		case string(RefreshClassState):
			return RefreshClassState, true
		case string(RefreshClassConfig):
			return RefreshClassConfig, true
		}
	}
	return "", false
}

func buildInvokeCacheKey(
	planeName string,
	methodName string,
	params map[string]any,
) (string, error) {
	normalizedParams := params
	if normalizedParams == nil {
		normalizedParams = map[string]any{}
	}
	serialized, err := json.Marshal(normalizedParams)
	if err != nil {
		return "", err
	}
	return planeName + "|" + methodName + "|" + string(serialized), nil
}

func cloneFrame(frame protocol.Frame) protocol.Frame {
	return protocol.Frame{
		Source:    frame.Source,
		Target:    frame.Target,
		Primary:   frame.Primary,
		Secondary: frame.Secondary,
		Data:      append([]byte(nil), frame.Data...),
	}
}

func findMethod(methods []registry.Method, name string) (registry.Method, bool) {
	for _, method := range methods {
		if method.Name() == name {
			return method, true
		}
	}
	return nil, false
}
