package registry

import (
	"context"
	"errors"
	"testing"

	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
)

// TestScanDirectedRejectsNilTargets confirms that ScanDirected refuses
// to fall through to DefaultScanTargets when given nil input. AD10 in
// startup-admission-discovery-w17-26 requires this: startup admission
// on non-ebusd-tcp direct transports MUST never emit a full-range scan.
func TestScanDirectedRejectsNilTargets(t *testing.T) {
	registry := NewDeviceRegistry(nil)
	bus := &mockScanBus{}
	entries, err := ScanDirected(context.Background(), bus, registry, 0x71, nil)
	if !errors.Is(err, ErrScanDirectedEmptyTargets) {
		t.Fatalf("expected ErrScanDirectedEmptyTargets, got err=%v", err)
	}
	if entries != nil {
		t.Errorf("expected nil entries on error, got %d", len(entries))
	}
	if len(bus.calls) != 0 {
		t.Fatalf("expected zero bus traffic on nil targets; got %d frames (this would defeat the directed-scan contract)", len(bus.calls))
	}
}

// TestScanDirectedRejectsEmptyTargets mirrors the nil case for an
// explicitly empty non-nil slice, which must also be refused.
func TestScanDirectedRejectsEmptyTargets(t *testing.T) {
	registry := NewDeviceRegistry(nil)
	bus := &mockScanBus{}
	entries, err := ScanDirected(context.Background(), bus, registry, 0x71, []byte{})
	if !errors.Is(err, ErrScanDirectedEmptyTargets) {
		t.Fatalf("expected ErrScanDirectedEmptyTargets, got err=%v", err)
	}
	if entries != nil {
		t.Errorf("expected nil entries on error, got %d", len(entries))
	}
	if len(bus.calls) != 0 {
		t.Fatalf("expected zero bus traffic on empty targets; got %d frames", len(bus.calls))
	}
}

// TestScanDirectedFiltersInitiatorCapableAndSYNESCTargets verifies the
// directed API honors the same target-capability rules as the legacy
// Scan: initiator-capable addresses and SYN (0xAA) / ESC (0xA9) are
// never probed even if the caller erroneously includes them. Per AD10
// filtering must be "consistent with scan.go:64-67".
func TestScanDirectedFiltersInitiatorCapableAndSYNESCTargets(t *testing.T) {
	registry := NewDeviceRegistry(nil)
	bus := &mockScanBus{}
	// Include a realistic target (0x08 BAI), SYN (0xAA), ESC (0xA9), and
	// an initiator-capable address (0x00). Directed scan should only
	// probe 0x08.
	_, err := ScanDirected(context.Background(), bus, registry, 0x71, []byte{0x08, 0xAA, 0xA9, 0x00})
	// ErrNoSuchDevice is expected from mockScanBus for unhandled targets;
	// we're asserting against the send-side, not the response-side.
	_ = err

	for _, frame := range bus.calls {
		switch frame.Target {
		case 0xAA:
			t.Errorf("filter leak: directed scan sent frame to SYN (0xAA)")
		case 0xA9:
			t.Errorf("filter leak: directed scan sent frame to ESC (0xA9)")
		case 0x00:
			t.Errorf("filter leak: directed scan sent frame to initiator-capable address 0x00")
		}
	}
}

// TestScanDirectedEquivalentToScanOnNonEmptyTargets pins the invariant
// that ScanDirected and Scan produce identical wire traffic (same
// frame set, same order) when given identical non-empty targets.
// This is the directed API's compatibility contract: it is a
// refusal-to-default wrapper, not a behavioral fork.
func TestScanDirectedEquivalentToScanOnNonEmptyTargets(t *testing.T) {
	targets := []byte{0x08, 0x15, 0x26}

	regA := NewDeviceRegistry(nil)
	busA := &mockScanBus{}
	_, _ = Scan(context.Background(), busA, regA, 0x71, targets)

	regB := NewDeviceRegistry(nil)
	busB := &mockScanBus{}
	_, _ = ScanDirected(context.Background(), busB, regB, 0x71, targets)

	if len(busA.calls) != len(busB.calls) {
		t.Fatalf("frame-count divergence: Scan=%d ScanDirected=%d", len(busA.calls), len(busB.calls))
	}
	for i := range busA.calls {
		a, b := busA.calls[i], busB.calls[i]
		if a.Source != b.Source || a.Target != b.Target || a.Primary != b.Primary || a.Secondary != b.Secondary {
			t.Errorf("frame %d divergence: Scan=%+v ScanDirected=%+v", i, a, b)
		}
	}
}

// TestScanNilTargetsStillFallsThroughToDefault is the regression guard
// for the ebusd-tcp sanctioned bounded-retry path. Introducing
// ScanDirected MUST NOT silently alter Scan's historical fall-through
// semantics — external tooling and the ebusd-tcp adapter rely on it.
func TestScanNilTargetsStillFallsThroughToDefault(t *testing.T) {
	registry := NewDeviceRegistry(nil)
	bus := &mockScanBus{}
	_, _ = Scan(context.Background(), bus, registry, 0x71, nil)
	if len(bus.calls) == 0 {
		t.Fatalf("regression: Scan(nil) emitted zero frames — fall-through to DefaultScanTargets is broken")
	}
	// Spot check: at least one frame should target a known slave-range
	// address that DefaultScanTargets includes (e.g. 0x08 BAI).
	sentToBAI := false
	for _, frame := range bus.calls {
		if frame.Target == 0x08 {
			sentToBAI = true
			break
		}
	}
	if !sentToBAI {
		t.Errorf("regression: Scan(nil) fall-through did not probe 0x08 (BAI); DefaultScanTargets may be broken")
	}
}

// TestScanDirectedDeduplicatesTargets verifies that duplicate target
// addresses in the caller's list are deduplicated (reusing the same
// dedupeScanTargets logic as Scan) so the directed API cannot be
// tricked into probing the same address multiple times per pass.
func TestScanDirectedDeduplicatesTargets(t *testing.T) {
	registry := NewDeviceRegistry(nil)
	bus := &mockScanBus{
		responses: map[byte]*protocol.Frame{},
	}
	_, _ = ScanDirected(context.Background(), bus, registry, 0x71, []byte{0x08, 0x08, 0x08, 0x15, 0x15})
	// With two distinct capable targets, the first pass should emit
	// at most 2 frames (mockScanBus returns ErrNoSuchDevice so no
	// retry pass fires; the collision-retry loop needs collided or
	// retried addresses, and ErrNoSuchDevice is neither).
	if len(bus.calls) > 2 {
		t.Errorf("expected ≤2 frames after dedupe, got %d: %+v", len(bus.calls), bus.calls)
	}
}
