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
	response *protocol.Frame
	calls    []protocol.Frame
	collided bool
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
				Data:      []byte{0x01, 'D', 'E', 'V', '0', '8', 0x05, 0x06, 0x07, 0x08},
			},
			0x09: {
				Source:    0x09,
				Target:    0x10,
				Primary:   scanPrimary,
				Secondary: scanSecondary,
				Data:      []byte{0x0A, 'D', 'E', 'V', '0', '9', 0x0E, 0x0F, 0x10, 0x11},
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
	if entry.Manufacturer() != "0x01" || entry.DeviceID() != "DEV08" || entry.SoftwareVersion() != "0506" || entry.HardwareVersion() != "0708" {
		t.Fatalf("unexpected device 0x08 info: %+v", entry)
	}

	entry, ok = registry.Lookup(0x09)
	if !ok {
		t.Fatalf("expected device 0x09 to be registered")
	}
	if entry.Manufacturer() != "0x0A" || entry.DeviceID() != "DEV09" || entry.SoftwareVersion() != "0E0F" || entry.HardwareVersion() != "1011" {
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
			Data:      []byte{0x01, 'D', 'E', 'V', '0', '8', 0x05, 0x06, 0x07, 0x08},
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
				Data:      []byte{0x01, 'D', 'E', 'V', '0', '8', 0x05, 0x06, 0x07, 0x08},
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
	if len(bus.calls) != 6 {
		t.Fatalf("expected 6 scan calls, got %d", len(bus.calls))
	}
}

type vaillantScanIDBus struct {
	calls []protocol.Frame
}

func (bus *vaillantScanIDBus) Send(ctx context.Context, frame protocol.Frame) (*protocol.Frame, error) {
	bus.calls = append(bus.calls, frame)

	if frame.Primary == scanPrimary && frame.Secondary == scanSecondary {
		if frame.Target != 0x20 {
			return nil, ebuserrors.ErrNoSuchDevice
		}
		return &protocol.Frame{
			Source:    0x30,
			Target:    frame.Source,
			Primary:   scanPrimary,
			Secondary: scanSecondary,
			Data:      []byte{0xB5, 'D', 'E', 'V', '3', '0', 0x05, 0x14, 0x12, 0x04},
		}, nil
	}

	if frame.Primary == vaillantPrimary && frame.Secondary == vaillantScanIDSecondary {
		if frame.Target != 0x30 || len(frame.Data) != 1 {
			return nil, ebuserrors.ErrNoSuchDevice
		}

		var chunk []byte
		switch frame.Data[0] {
		case 0x24:
			chunk = []byte("21220900")
		case 0x25:
			chunk = []byte("20184848")
		case 0x26:
			chunk = []byte("00820054")
		case 0x27:
			chunk = []byte{'0', '9', 'N', '4', 0x00, 0x00, 0x00, 0x00}
		default:
			return nil, ebuserrors.ErrNoSuchDevice
		}

		return &protocol.Frame{
			Source:    frame.Target,
			Target:    frame.Source,
			Primary:   vaillantPrimary,
			Secondary: vaillantScanIDSecondary,
			Data:      append([]byte{0x00}, chunk...),
		}, nil
	}

	return nil, ebuserrors.ErrNoSuchDevice
}

type vaillantScanIDCanceledBus struct {
	calls []protocol.Frame
}

func (bus *vaillantScanIDCanceledBus) Send(ctx context.Context, frame protocol.Frame) (*protocol.Frame, error) {
	bus.calls = append(bus.calls, frame)

	if frame.Primary == scanPrimary && frame.Secondary == scanSecondary {
		return &protocol.Frame{
			Source:    0x30,
			Target:    frame.Source,
			Primary:   scanPrimary,
			Secondary: scanSecondary,
			Data:      []byte{0xB5, 'D', 'E', 'V', '3', '0', 0x05, 0x14, 0x12, 0x04},
		}, nil
	}
	if frame.Primary == vaillantPrimary && frame.Secondary == vaillantScanIDSecondary {
		return nil, context.Canceled
	}
	return nil, ebuserrors.ErrNoSuchDevice
}

type vaillantScanIDTimeoutBus struct{}

func (bus *vaillantScanIDTimeoutBus) Send(ctx context.Context, frame protocol.Frame) (*protocol.Frame, error) {
	if frame.Primary == scanPrimary && frame.Secondary == scanSecondary {
		return &protocol.Frame{
			Source:    0x30,
			Target:    frame.Source,
			Primary:   scanPrimary,
			Secondary: scanSecondary,
			Data:      []byte{0xB5, 'D', 'E', 'V', '3', '0', 0x05, 0x14, 0x12, 0x04},
		}, nil
	}
	if frame.Primary == vaillantPrimary && frame.Secondary == vaillantScanIDSecondary {
		return nil, context.DeadlineExceeded
	}
	return nil, ebuserrors.ErrNoSuchDevice
}

func TestScanUsesDiscoveredAddressForVaillantScanID(t *testing.T) {
	t.Parallel()

	registry := NewDeviceRegistry(nil)
	bus := &vaillantScanIDBus{}

	entries, err := Scan(context.Background(), bus, registry, 0x10, []byte{0x20})
	if err != nil {
		t.Fatalf("Scan error = %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	entry, ok := registry.Lookup(0x30)
	if !ok {
		t.Fatalf("expected device 0x30 to be registered")
	}
	if entry.SerialNumber() != "21-22-09-0020184848-0082-005409-N4" {
		t.Fatalf("unexpected serial number: %q", entry.SerialNumber())
	}

	if len(bus.calls) != 5 {
		t.Fatalf("expected 5 calls (scan + 4 scan.id), got %d", len(bus.calls))
	}
}

func TestScanPropagatesContextCancellationFromVaillantScanID(t *testing.T) {
	t.Parallel()

	registry := NewDeviceRegistry(nil)
	bus := &vaillantScanIDCanceledBus{}

	entries, err := Scan(context.Background(), bus, registry, 0x10, []byte{0x20})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Scan error = %v; want context.Canceled", err)
	}
	if entries != nil {
		t.Fatalf("expected nil entries on cancellation, got %d", len(entries))
	}
}

func TestScanIgnoresScanIDTimeouts(t *testing.T) {
	t.Parallel()

	registry := NewDeviceRegistry(nil)
	bus := &vaillantScanIDTimeoutBus{}

	entries, err := Scan(context.Background(), bus, registry, 0x10, []byte{0x20})
	if err != nil {
		t.Fatalf("Scan error = %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	entry, ok := registry.Lookup(0x30)
	if !ok {
		t.Fatalf("expected device 0x30 to be registered")
	}
	if entry.SerialNumber() != "" {
		t.Fatalf("serial number = %q; want empty on timeout", entry.SerialNumber())
	}
}

func TestScanPreservesKnownSerialWhenScanIDFails(t *testing.T) {
	t.Parallel()

	registry := NewDeviceRegistry(nil)
	registry.Register(DeviceInfo{
		Address:         0x30,
		Manufacturer:    "Vaillant",
		DeviceID:        "DEV30",
		SerialNumber:    "21-22-09-0020184848-0082-005409-N4",
		SoftwareVersion: "0514",
		HardwareVersion: "1204",
	})
	bus := &vaillantScanIDTimeoutBus{}

	entries, err := Scan(context.Background(), bus, registry, 0x10, []byte{0x20})
	if err != nil {
		t.Fatalf("Scan error = %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	entry, ok := registry.Lookup(0x30)
	if !ok {
		t.Fatalf("expected device 0x30 to be registered")
	}
	if entry.SerialNumber() != "21-22-09-0020184848-0082-005409-N4" {
		t.Fatalf("serial number = %q; want preserved value", entry.SerialNumber())
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
				Data:      []byte{0x01, 'D', 'E', 'V', '0', '8', 0x05, 0x06, 0x07, 0x08},
			},
			0x09: {
				Source:    0x09,
				Target:    0x10,
				Primary:   scanPrimary,
				Secondary: scanSecondary,
				Data:      []byte{0x0A, 'D', 'E', 'V', '0', '9', 0x0E, 0x0F, 0x10, 0x11},
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
	if len(bus.calls) != 6 {
		t.Fatalf("expected 6 scan calls, got %d", len(bus.calls))
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

func TestScanSkipsInitiatorCapableAndUnknownTargets(t *testing.T) {
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
				Data:      []byte{0x0A, 'D', 'E', 'V', '0', '9', 0x0E, 0x0F, 0x10, 0x11},
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
	if len(bus.calls) != 5 {
		t.Fatalf("expected 5 scan calls, got %d", len(bus.calls))
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
				Data:      []byte{0x01, 'D', 'E', 'V', '0', '8', 0x05, 0x06, 0x07, 0x08},
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

func TestFormatVaillantSerial(t *testing.T) {
	t.Parallel()

	formatted := formatVaillantSerial("21220900201848480082005409N4")
	if formatted != "21-22-09-0020184848-0082-005409-N4" {
		t.Fatalf("formatted = %q", formatted)
	}

	rawWithSuffix := "2122ABCD0020184848XYZZ005409N4TAIL"
	got := formatVaillantSerial(rawWithSuffix)
	if got != rawWithSuffix {
		t.Fatalf("fallback serial = %q; want full raw value", got)
	}
}
