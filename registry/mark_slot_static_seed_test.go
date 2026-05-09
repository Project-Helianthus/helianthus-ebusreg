package registry

import (
	"sync"
	"testing"
	"time"
)

// TestMarkSlotStaticSeed_BasicWrite asserts the API stamps Role,
// DiscoverySource=StaticSeed, VerificationState=Candidate, and
// timestamps on a fresh slot.
func TestMarkSlotStaticSeed_BasicWrite(t *testing.T) {
	t.Parallel()

	reg := NewDeviceRegistry(nil)
	now := time.Date(2026, 5, 9, 12, 0, 0, 0, time.UTC)
	reg.MarkSlotStaticSeed(0xF1, SlotRoleMaster, now)

	slot, ok := reg.LookupSlot(0xF1)
	if !ok || slot == nil {
		t.Fatalf("LookupSlot(0xF1) ok=%v slot=%v; want non-nil slot", ok, slot)
	}
	if slot.Role != SlotRoleMaster {
		t.Errorf("slot.Role = %v; want SlotRoleMaster", slot.Role)
	}
	if slot.DiscoverySource != DiscoverySourceStaticSeed {
		t.Errorf("slot.DiscoverySource = %v; want DiscoverySourceStaticSeed", slot.DiscoverySource)
	}
	if slot.VerificationState != VerificationStateCandidate {
		t.Errorf("slot.VerificationState = %v; want VerificationStateCandidate", slot.VerificationState)
	}
	if !slot.FirstObservedAt.Equal(now) {
		t.Errorf("slot.FirstObservedAt = %v; want %v", slot.FirstObservedAt, now)
	}
	if !slot.LastObservedAt.Equal(now) {
		t.Errorf("slot.LastObservedAt = %v; want %v", slot.LastObservedAt, now)
	}
}

// TestMarkSlotStaticSeed_MonotonicNoDowngrade asserts that re-marking
// a slot already at higher DiscoverySource (ActiveConfirmed) does NOT
// downgrade it to StaticSeed. Same shape as
// TestMarkSlotPassiveObserved_MonotonicMetadata.
func TestMarkSlotStaticSeed_MonotonicNoDowngrade(t *testing.T) {
	t.Parallel()

	reg := NewDeviceRegistry(nil)
	// Bring the slot to ActiveConfirmed via Register.
	reg.Register(DeviceInfo{Address: 0xF1, Manufacturer: "Vaillant", DeviceID: "NETX3", SerialNumber: "SN-A"})
	pre, _ := reg.LookupSlot(0xF1)
	if pre.DiscoverySource != DiscoverySourceActiveConfirmed {
		t.Fatalf("pre-condition: slot.DiscoverySource = %v; want ActiveConfirmed", pre.DiscoverySource)
	}

	// Static-seed call must not downgrade.
	reg.MarkSlotStaticSeed(0xF1, SlotRoleMaster, time.Now())
	post, _ := reg.LookupSlot(0xF1)
	if post.DiscoverySource != DiscoverySourceActiveConfirmed {
		t.Errorf("slot.DiscoverySource after MarkSlotStaticSeed = %v; want ActiveConfirmed (no downgrade)", post.DiscoverySource)
	}
	if post.VerificationState != VerificationStateIdentityConfirmed {
		t.Errorf("slot.VerificationState after MarkSlotStaticSeed = %v; want IdentityConfirmed (no downgrade)", post.VerificationState)
	}
}

// TestMarkSlotStaticSeed_MonotonicUpgradeFromPassive asserts that the
// passive→static-seed transition correctly upgrades DiscoverySource
// from PassiveObserved to StaticSeed (per the enum order
// Unknown < PassiveObserved < StaticSeed < ActiveConfirmed). Note:
// this is the asymmetry codified in Codex Pass 4 of the P6+P3.5 plan
// — static-seed pre-known taxonomy "outranks" passive-only inference,
// so a passively-observed slot subsequently learned to be a static
// seed advances; the reverse direction (ActiveConfirmed -> StaticSeed)
// is rejected by MonotonicNoDowngrade above.
func TestMarkSlotStaticSeed_MonotonicUpgradeFromPassive(t *testing.T) {
	t.Parallel()

	reg := NewDeviceRegistry(nil)
	earlier := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
	reg.MarkSlotPassiveObserved(0xF1, SlotRoleMaster, earlier)
	pre, _ := reg.LookupSlot(0xF1)
	if pre.DiscoverySource != DiscoverySourcePassiveObserved {
		t.Fatalf("pre-condition: slot.DiscoverySource = %v; want PassiveObserved", pre.DiscoverySource)
	}

	later := earlier.Add(24 * time.Hour)
	reg.MarkSlotStaticSeed(0xF1, SlotRoleMaster, later)

	post, _ := reg.LookupSlot(0xF1)
	if post.DiscoverySource != DiscoverySourceStaticSeed {
		t.Errorf("slot.DiscoverySource after passive->static-seed = %v; want StaticSeed (upgrade)", post.DiscoverySource)
	}
	// VerificationState was Corroborated (from passive); StaticSeed
	// brings Candidate which is < Corroborated, so the slot must
	// retain Corroborated.
	if post.VerificationState != VerificationStateCorroborated {
		t.Errorf("slot.VerificationState after passive->static-seed = %v; want Corroborated (retained)", post.VerificationState)
	}
}

// TestMarkSlotStaticSeed_RaceFreeWriteAndRead mirrors the M6.1 hardening
// proof for MarkSlotPassiveObserved. The API must serialize via the
// registry's RWMutex; concurrent writers and readers must not race.
func TestMarkSlotStaticSeed_RaceFreeWriteAndRead(t *testing.T) {
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
				reg.MarkSlotStaticSeed(0xF1, SlotRoleMaster, time.Now())
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
	if slot.DiscoverySource != DiscoverySourceStaticSeed {
		t.Errorf("after race test: slot.DiscoverySource = %v; want StaticSeed", slot.DiscoverySource)
	}
}
