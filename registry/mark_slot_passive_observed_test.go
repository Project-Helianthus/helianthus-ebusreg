package registry

import (
	"sync"
	"testing"
	"time"
)

// TestMarkSlotPassiveObserved_BasicWrite asserts the API stamps Role,
// DiscoverySource=PassiveObserved, VerificationState=Corroborated, and
// timestamps on a fresh slot.
func TestMarkSlotPassiveObserved_BasicWrite(t *testing.T) {
	t.Parallel()

	reg := NewDeviceRegistry(nil)
	now := time.Date(2026, 5, 7, 12, 0, 0, 0, time.UTC)
	reg.MarkSlotPassiveObserved(0xF1, SlotRoleMaster, now)

	slot, ok := reg.LookupSlot(0xF1)
	if !ok || slot == nil {
		t.Fatalf("LookupSlot(0xF1) ok=%v slot=%v; want non-nil slot", ok, slot)
	}
	if slot.Role != SlotRoleMaster {
		t.Errorf("slot.Role = %v; want SlotRoleMaster", slot.Role)
	}
	if slot.DiscoverySource != DiscoverySourcePassiveObserved {
		t.Errorf("slot.DiscoverySource = %v; want DiscoverySourcePassiveObserved", slot.DiscoverySource)
	}
	if slot.VerificationState != VerificationStateCorroborated {
		t.Errorf("slot.VerificationState = %v; want VerificationStateCorroborated", slot.VerificationState)
	}
	if !slot.FirstObservedAt.Equal(now) {
		t.Errorf("slot.FirstObservedAt = %v; want %v", slot.FirstObservedAt, now)
	}
	if !slot.LastObservedAt.Equal(now) {
		t.Errorf("slot.LastObservedAt = %v; want %v", slot.LastObservedAt, now)
	}
}

// TestMarkSlotPassiveObserved_MonotonicMetadata asserts that re-marking
// a slot already at higher confidence (e.g. ActiveConfirmed) does NOT
// downgrade it to PassiveObserved.
func TestMarkSlotPassiveObserved_MonotonicMetadata(t *testing.T) {
	t.Parallel()

	reg := NewDeviceRegistry(nil)
	// Bring the slot to ActiveConfirmed via Register (which calls
	// observeAddressSlotLocked with ActiveConfirmed/IdentityConfirmed).
	reg.Register(DeviceInfo{Address: 0xF1, Manufacturer: "Vaillant", DeviceID: "NETX3", SerialNumber: "SN-A"})
	pre, _ := reg.LookupSlot(0xF1)
	if pre.DiscoverySource != DiscoverySourceActiveConfirmed {
		t.Fatalf("pre-condition: slot.DiscoverySource = %v; want ActiveConfirmed", pre.DiscoverySource)
	}

	// Now passively re-observe — must not downgrade.
	reg.MarkSlotPassiveObserved(0xF1, SlotRoleMaster, time.Now())
	post, _ := reg.LookupSlot(0xF1)
	if post.DiscoverySource != DiscoverySourceActiveConfirmed {
		t.Errorf("slot.DiscoverySource after MarkSlotPassiveObserved = %v; want ActiveConfirmed (no downgrade)", post.DiscoverySource)
	}
	if post.VerificationState != VerificationStateIdentityConfirmed {
		t.Errorf("slot.VerificationState after MarkSlotPassiveObserved = %v; want IdentityConfirmed (no downgrade)", post.VerificationState)
	}
}

// TestMarkSlotPassiveObserved_RaceFreeWriteAndRead is the M6.1 hardening
// proof: 100 goroutines concurrently call MarkSlotPassiveObserved while
// 100 goroutines call LookupSlot. With direct field mutation (pre-A.7
// behavior) the race detector flags a data race. After M6.1 the API
// must serialize via the registry's RWMutex.
func TestMarkSlotPassiveObserved_RaceFreeWriteAndRead(t *testing.T) {
	reg := NewDeviceRegistry(nil)
	const writers = 100
	const readers = 100
	const iterations = 100

	var wg sync.WaitGroup
	wg.Add(writers + readers)

	for i := 0; i < writers; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				reg.MarkSlotPassiveObserved(0xF1, SlotRoleMaster, time.Now())
			}
		}()
	}
	for i := 0; i < readers; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				slot, ok := reg.LookupSlot(0xF1)
				if ok && slot != nil {
					_ = slot.Role
					_ = slot.DiscoverySource
					_ = slot.VerificationState
				}
			}
		}()
	}
	wg.Wait()

	slot, ok := reg.LookupSlot(0xF1)
	if !ok || slot == nil {
		t.Fatalf("after race test: LookupSlot(0xF1) ok=%v slot=%v", ok, slot)
	}
	if slot.DiscoverySource != DiscoverySourcePassiveObserved {
		t.Errorf("after race test: slot.DiscoverySource = %v; want PassiveObserved", slot.DiscoverySource)
	}
}
