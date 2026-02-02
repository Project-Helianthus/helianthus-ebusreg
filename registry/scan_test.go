package registry

import (
	"context"
	"errors"
	"testing"
	"time"

	ebuserrors "github.com/d3vi1/helianthus-ebusgo/errors"
	"github.com/d3vi1/helianthus-ebusgo/protocol"
)

type mockScanBus struct {
	responses map[byte]*protocol.Frame
	errors    map[byte]error
	calls     []protocol.Frame
}

func (bus *mockScanBus) Send(ctx context.Context, frame protocol.Frame) (*protocol.Frame, error) {
	bus.calls = append(bus.calls, frame)
	if err, ok := bus.errors[frame.Target]; ok {
		return nil, err
	}
	if response, ok := bus.responses[frame.Target]; ok {
		return response, nil
	}
	return nil, ebuserrors.ErrNoSuchDevice
}

type collisionThenResponseBus struct {
	response  *protocol.Frame
	calls     []protocol.Frame
	collided  bool
}

func (bus *collisionThenResponseBus) Send(ctx context.Context, frame protocol.Frame) (*protocol.Frame, error) {
	bus.calls = append(bus.calls, frame)
	if !bus.collided {
		bus.collided = true
		return nil, ebuserrors.ErrBusCollision
	}
	return bus.response, nil
}

func TestScanRegistersDevices(t *testing.T) {
	t.Parallel()

	registry := NewDeviceRegistry(nil)
	bus := &mockScanBus{
		responses: map[byte]*protocol.Frame{
			0x08: {
				Source:    0x08,
				Target:    0x10,
				Primary:   scanPrimary,
				Secondary: scanSecondary,
				Data:      []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
			},
			0x09: {
				Source:    0x09,
				Target:    0x10,
				Primary:   scanPrimary,
				Secondary: scanSecondary,
				Data:      []byte{0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11},
			},
		},
		errors: map[byte]error{
			0x21: ebuserrors.ErrNoSuchDevice,
		},
	}

	entries, err := Scan(context.Background(), bus, registry, 0x10, []byte{0x08, 0x09, 0x21})
	if err != nil {
		t.Fatalf("Scan error = %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}

	entry, ok := registry.Lookup(0x08)
	if !ok {
		t.Fatalf("expected device 0x08 to be registered")
	}
	if entry.Manufacturer() != "0102" || entry.DeviceID() != "0304" || entry.SoftwareVersion() != "0506" || entry.HardwareVersion() != "0708" {
		t.Fatalf("unexpected device 0x08 info: %+v", entry)
	}

	entry, ok = registry.Lookup(0x09)
	if !ok {
		t.Fatalf("expected device 0x09 to be registered")
	}
	if entry.Manufacturer() != "0A0B" || entry.DeviceID() != "0C0D" || entry.SoftwareVersion() != "0E0F" || entry.HardwareVersion() != "1011" {
		t.Fatalf("unexpected device 0x09 info: %+v", entry)
	}

	if len(bus.calls) != 3 {
		t.Fatalf("expected 3 scan calls, got %d", len(bus.calls))
	}
}

func TestScanRetriesBusCollision(t *testing.T) {
	t.Parallel()

	registry := NewDeviceRegistry(nil)
	bus := &collisionThenResponseBus{
		response: &protocol.Frame{
			Source:    0x08,
			Target:    0x10,
			Primary:   scanPrimary,
			Secondary: scanSecondary,
			Data:      []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
		},
	}

	entries, err := Scan(context.Background(), bus, registry, 0x10, []byte{0x08})
	if err != nil {
		t.Fatalf("Scan error = %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if _, ok := registry.Lookup(0x08); !ok {
		t.Fatalf("expected device 0x08 to be registered")
	}
	if len(bus.calls) != 2 {
		t.Fatalf("expected 2 scan calls, got %d", len(bus.calls))
	}
}

func TestScanSkipsTimeoutAndNACK(t *testing.T) {
	t.Parallel()

	registry := NewDeviceRegistry(nil)
	bus := &mockScanBus{
		responses: map[byte]*protocol.Frame{
			0x08: {
				Source:    0x08,
				Target:    0x10,
				Primary:   scanPrimary,
				Secondary: scanSecondary,
				Data:      []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
			},
		},
		errors: map[byte]error{
			0x21: ebuserrors.ErrTimeout,
			0x22: ebuserrors.ErrNACK,
		},
	}

	entries, err := Scan(context.Background(), bus, registry, 0x10, []byte{0x08, 0x21, 0x22})
	if err != nil {
		t.Fatalf("Scan error = %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if _, ok := registry.Lookup(0x08); !ok {
		t.Fatalf("expected device 0x08 to be registered")
	}
	if len(bus.calls) != 3 {
		t.Fatalf("expected 3 scan calls, got %d", len(bus.calls))
	}
}

func TestScanSkipsContextDeadlineExceeded(t *testing.T) {
	t.Parallel()

	registry := NewDeviceRegistry(nil)
	bus := &mockScanBus{
		responses: map[byte]*protocol.Frame{
			0x08: {
				Source:    0x08,
				Target:    0x10,
				Primary:   scanPrimary,
				Secondary: scanSecondary,
				Data:      []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
			},
			0x09: {
				Source:    0x09,
				Target:    0x10,
				Primary:   scanPrimary,
				Secondary: scanSecondary,
				Data:      []byte{0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11},
			},
		},
		errors: map[byte]error{
			0x21: context.DeadlineExceeded,
		},
	}

	entries, err := Scan(context.Background(), bus, registry, 0x10, []byte{0x08, 0x21, 0x09})
	if err != nil {
		t.Fatalf("Scan error = %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if _, ok := registry.Lookup(0x08); !ok {
		t.Fatalf("expected device 0x08 to be registered")
	}
	if _, ok := registry.Lookup(0x09); !ok {
		t.Fatalf("expected device 0x09 to be registered")
	}
	if len(bus.calls) != 3 {
		t.Fatalf("expected 3 scan calls, got %d", len(bus.calls))
	}
}

func TestScanFailsWhenContextDone(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
	defer cancel()

	registry := NewDeviceRegistry(nil)
	bus := &mockScanBus{}

	entries, err := Scan(ctx, bus, registry, 0x10, []byte{0x08})
	if err == nil {
		t.Fatalf("expected Scan error")
	}
	if entries != nil {
		t.Fatalf("expected no entries, got %d", len(entries))
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected DeadlineExceeded, got %v", err)
	}
	if len(bus.calls) != 0 {
		t.Fatalf("expected 0 scan calls, got %d", len(bus.calls))
	}
}

func TestScanSkipsMasterAndUnknownTargets(t *testing.T) {
	t.Parallel()

	registry := NewDeviceRegistry(nil)
	bus := &mockScanBus{}

	entries, err := Scan(context.Background(), bus, registry, 0x10, []byte{0x10, protocol.SymbolEscape})
	if err != nil {
		t.Fatalf("Scan error = %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
	if len(bus.calls) != 0 {
		t.Fatalf("expected 0 scan calls, got %d", len(bus.calls))
	}
}

func TestScanSkipsInvalidPayloadResponses(t *testing.T) {
	t.Parallel()

	registry := NewDeviceRegistry(nil)
	bus := &mockScanBus{
		responses: map[byte]*protocol.Frame{
			0x08: {
				Source:    0x08,
				Target:    0x10,
				Primary:   scanPrimary,
				Secondary: scanSecondary,
				Data:      []byte{0x01},
			},
			0x09: {
				Source:    0x09,
				Target:    0x10,
				Primary:   scanPrimary,
				Secondary: scanSecondary,
				Data:      []byte{0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11},
			},
		},
	}

	entries, err := Scan(context.Background(), bus, registry, 0x10, []byte{0x08, 0x09})
	if err != nil {
		t.Fatalf("Scan error = %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if _, ok := registry.Lookup(0x09); !ok {
		t.Fatalf("expected device 0x09 to be registered")
	}
	if len(bus.calls) != 2 {
		t.Fatalf("expected 2 scan calls, got %d", len(bus.calls))
	}
}

func TestScanSkipsCRCMismatchAndFailsOnInvalidPayloadError(t *testing.T) {
	t.Parallel()

	registry := NewDeviceRegistry(nil)
	bus := &mockScanBus{
		responses: map[byte]*protocol.Frame{
			0x08: {
				Source:    0x08,
				Target:    0x10,
				Primary:   scanPrimary,
				Secondary: scanSecondary,
				Data:      []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
			},
		},
		errors: map[byte]error{
			0x21: ebuserrors.ErrCRCMismatch,
			0x22: ebuserrors.ErrInvalidPayload,
		},
	}

	entries, err := Scan(context.Background(), bus, registry, 0x10, []byte{0x08, 0x21, 0x22})
	if err == nil {
		t.Fatalf("expected Scan error")
	}
	if entries != nil {
		t.Fatalf("expected no entries, got %d", len(entries))
	}
	if !errors.Is(err, ebuserrors.ErrInvalidPayload) {
		t.Fatalf("expected ErrInvalidPayload, got %v", err)
	}
	if len(bus.calls) != 3 {
		t.Fatalf("expected 3 scan calls, got %d", len(bus.calls))
	}
}
