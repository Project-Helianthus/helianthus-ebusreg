package router

import (
	"sync"

	"github.com/d3vi1/helianthus-ebusgo/protocol"
	"github.com/d3vi1/helianthus-ebusreg/registry"
)

type Subscription struct {
	Primary   byte
	Secondary byte
}

type Plane interface {
	Name() string
	Methods() []registry.Method
	Subscriptions() []Subscription
	OnBroadcast(frame protocol.Frame) error
}

type BusEventRouter struct {
	mu            sync.RWMutex
	subscriptions map[subscriptionKey][]Plane
}

type subscriptionKey struct {
	primary   byte
	secondary byte
}

func NewBusEventRouter() *BusEventRouter {
	return &BusEventRouter{
		subscriptions: make(map[subscriptionKey][]Plane),
	}
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
	}
	return errors
}
