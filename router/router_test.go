package router

import (
	"errors"
	"testing"

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

func TestBusEventRouter_BroadcastFanOut(t *testing.T) {
	t.Parallel()

	router := NewBusEventRouter()

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

	router := NewBusEventRouter()
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
