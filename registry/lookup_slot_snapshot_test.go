package registry

import (
	"sync"
	"testing"
	"time"
)

// P8.1 — LookupSlotSnapshot value-typed snapshot API.
//
// Addresses the lock-free read advisory: callers reading a live
// AddressSlot pointer returned by LookupSlot risk a race with
// concurrent writers (Register / MarkSlot* / RegisterPassiveObserved).
// LookupSlotSnapshot copies the observable fields under r.mu.RLock
// so the caller's subsequent reads are immune to mutation.

// TestLookupSlotSnapshot_ReturnsCurrentSlotState asserts the
// snapshot reflects the slot's current label / role / timing fields.
func TestLookupSlotSnapshot_ReturnsCurrentSlotState(t *testing.T) {
	t.Parallel()

	reg := NewDeviceRegistry(nil)
	now := time.Date(2026, 5, 9, 12, 0, 0, 0, time.UTC)
	reg.RegisterPassiveObserved(DeviceInfo{Address: 0xF1}, SlotRoleMaster, now)

	snap, ok := reg.LookupSlotSnapshot(0xF1)
	if !ok {
		t.Fatalf("LookupSlotSnapshot(0xF1) ok=false")
	}
	if snap.Addr != 0xF1 {
		t.Errorf("snap.Addr = 0x%02X; want 0xF1", snap.Addr)
	}
	if snap.Role != SlotRoleMaster {
		t.Errorf("snap.Role = %v; want SlotRoleMaster", snap.Role)
	}
	if snap.DiscoverySource != DiscoverySourcePassiveObserved {
		t.Errorf("snap.DiscoverySource = %v; want PassiveObserved", snap.DiscoverySource)
	}
	if snap.VerificationState != VerificationStateCorroborated {
		t.Errorf("snap.VerificationState = %v; want Corroborated", snap.VerificationState)
	}
	if !snap.FirstObservedAt.Equal(now) {
		t.Errorf("snap.FirstObservedAt = %v; want %v", snap.FirstObservedAt, now)
	}
	if !snap.LastObservedAt.Equal(now) {
		t.Errorf("snap.LastObservedAt = %v; want %v", snap.LastObservedAt, now)
	}
	if !snap.DeviceAttached {
		t.Errorf("snap.DeviceAttached = false; want true (RegisterPassiveObserved attaches slot.Device)")
	}
}

// TestLookupSlotSnapshot_AbsentSlotReturnsZero verifies the false
// return for an address with no slot.
func TestLookupSlotSnapshot_AbsentSlotReturnsZero(t *testing.T) {
	t.Parallel()

	reg := NewDeviceRegistry(nil)
	snap, ok := reg.LookupSlotSnapshot(0x42)
	if ok {
		t.Errorf("LookupSlotSnapshot(0x42) ok=true; want false (no slot)")
	}
	if snap.Addr != 0 || snap.DiscoverySource != 0 {
		t.Errorf("snap = %+v; want zero value when ok=false", snap)
	}
}

// TestLookupSlotSnapshot_ReflectsLaterUpgrade asserts that calling
// LookupSlotSnapshot AFTER a registry upgrade returns the upgraded
// labels — proves the snapshot is taken at call time, not at
// insertion time.
func TestLookupSlotSnapshot_ReflectsLaterUpgrade(t *testing.T) {
	t.Parallel()

	reg := NewDeviceRegistry(nil)
	reg.RegisterPassiveObserved(DeviceInfo{Address: 0x15}, SlotRoleSlave, time.Now())

	pre, _ := reg.LookupSlotSnapshot(0x15)
	if pre.DiscoverySource != DiscoverySourcePassiveObserved {
		t.Fatalf("pre.DiscoverySource = %v; want PassiveObserved", pre.DiscoverySource)
	}

	reg.Register(DeviceInfo{
		Address:      0x15,
		Manufacturer: "Vaillant",
		DeviceID:     "BASV2",
		SerialNumber: "SN-X",
	})

	post, _ := reg.LookupSlotSnapshot(0x15)
	if post.DiscoverySource != DiscoverySourceActiveConfirmed {
		t.Errorf("post.DiscoverySource = %v; want ActiveConfirmed (snapshot taken AFTER upgrade)", post.DiscoverySource)
	}
	if post.VerificationState != VerificationStateIdentityConfirmed {
		t.Errorf("post.VerificationState = %v; want IdentityConfirmed", post.VerificationState)
	}
}

// TestLookupSlotSnapshot_RaceFreeUnderConcurrentWrites is the central
// P8.1 contract: the snapshot's observable fields are stable for the
// caller's lifetime regardless of concurrent registry mutations.
//
// Methodology: kick off a writer goroutine that hammers the slot via
// RegisterPassiveObserved (which mutates DiscoverySource /
// VerificationState / FirstObservedAt / LastObservedAt under
// r.mu.Lock). The reader takes a snapshot, then re-reads each field
// from its local copy several times — none of those reads should
// race with the writer because the snapshot is value-typed and
// disconnected from registry storage.
//
// The race detector (-race) is the authoritative gate: a torn read
// or unsync access via this code path would fail.
func TestLookupSlotSnapshot_RaceFreeUnderConcurrentWrites(t *testing.T) {
	t.Parallel()

	reg := NewDeviceRegistry(nil)
	addr := byte(0x10)
	reg.RegisterPassiveObserved(DeviceInfo{Address: addr}, SlotRoleMaster, time.Now())

	var wg sync.WaitGroup
	wg.Add(1)
	stop := make(chan struct{})
	go func() {
		defer wg.Done()
		i := 0
		for {
			select {
			case <-stop:
				return
			default:
			}
			reg.RegisterPassiveObserved(DeviceInfo{Address: addr}, SlotRoleMaster, time.Now())
			if i%10 == 0 {
				// Periodically promote to ActiveConfirmed and back —
				// exercises the monotonic-ladder branches under load.
				reg.Register(DeviceInfo{Address: addr, SerialNumber: "SN-RACE"})
			}
			i++
		}
	}()

	for i := 0; i < 200; i++ {
		snap, ok := reg.LookupSlotSnapshot(addr)
		if !ok {
			t.Fatalf("LookupSlotSnapshot(0x10) ok=false at iter %d", i)
		}
		// Each field read should be a plain memory read of a value
		// already copied under RLock — no race with the writer.
		_ = snap.Addr
		_ = snap.Role
		_ = snap.DiscoverySource
		_ = snap.VerificationState
		_ = snap.FirstObservedAt
		_ = snap.LastObservedAt
		_ = snap.DeviceAttached
	}

	close(stop)
	wg.Wait()
}
