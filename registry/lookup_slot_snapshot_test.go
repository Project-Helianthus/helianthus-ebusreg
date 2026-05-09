package registry

import (
	"sync"
	"sync/atomic"
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
	var writes uint64
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
			atomic.AddUint64(&writes, 1)
			if i%10 == 0 {
				// Periodically promote to ActiveConfirmed and back —
				// exercises the monotonic-ladder branches under load.
				reg.Register(DeviceInfo{Address: addr, SerialNumber: "SN-RACE"})
				atomic.AddUint64(&writes, 1)
			}
			i++
		}
	}()

	// Codex P8.1 review MINOR FINDING_1 — start barrier: ensure the
	// writer has actually mutated the slot at least once before the
	// reader phase begins, so the race detector reliably exercises
	// concurrent read/write overlap. Without this barrier the reader
	// could finish before the writer goroutine got CPU.
	deadline := time.Now().Add(2 * time.Second)
	for atomic.LoadUint64(&writes) == 0 {
		if time.Now().After(deadline) {
			close(stop)
			wg.Wait()
			t.Fatal("writer goroutine produced no writes within 2s — start barrier exceeded")
		}
		runtimeYield()
	}

	// Reader phase: drive enough iterations that the writer overlaps
	// with reads. Track the writes counter to confirm overlap actually
	// happened (vs. the writer completing before reads start).
	startWrites := atomic.LoadUint64(&writes)
	const readIters = 500
	for i := 0; i < readIters; i++ {
		snap, ok := reg.LookupSlotSnapshot(addr)
		if !ok {
			t.Fatalf("LookupSlotSnapshot(0x10) ok=false at iter %d", i)
		}
		// Each field read is a plain memory read of a value already
		// copied under RLock — no race with the writer.
		_ = snap.Addr
		_ = snap.Role
		_ = snap.DiscoverySource
		_ = snap.VerificationState
		_ = snap.FirstObservedAt
		_ = snap.LastObservedAt
		_ = snap.DeviceAttached
	}
	endWrites := atomic.LoadUint64(&writes)
	if endWrites <= startWrites {
		t.Errorf("writer produced %d additional writes during reader phase; want > 0 (race window not exercised)", endWrites-startWrites)
	}

	close(stop)
	wg.Wait()
}

// runtimeYield is a tiny helper to relinquish a scheduling slot to
// the writer goroutine without burning CPU in a tight spin.
func runtimeYield() {
	// time.Sleep(0) on Go is a runtime.Gosched() equivalent in many
	// schedulers, but explicit Gosched is clearer.
	time.Sleep(time.Microsecond)
}
