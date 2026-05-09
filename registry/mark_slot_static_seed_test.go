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

// TestMarkSlotStaticSeed_DoesNotAttachOrphanSlot documents the scope
// limit: MarkSlotStaticSeed on an address with no prior Device pointer
// only stamps the AddressSlot's labels; it does NOT attach the address
// to any device entry, does NOT add it to r.entries, and does NOT
// update any device's Faces. To plant identity for a new seeded
// address, callers must use RegisterStaticSeed (which composes
// registerLocked) — or AliasAddresses to attach an existing entry
// after-the-fact. Codex P3.5 review pass 2 caught the misleading
// previous doc claim that this API could mark "additional faces".
func TestMarkSlotStaticSeed_DoesNotAttachOrphanSlot(t *testing.T) {
	t.Parallel()

	reg := NewDeviceRegistry(nil)
	reg.MarkSlotStaticSeed(0x04, SlotRoleSlave, time.Now())

	if entry, ok := reg.Lookup(0x04); ok || entry != nil {
		t.Errorf("Lookup(0x04) ok=%v entry=%v; want (false, nil) — orphan slot must not produce a device entry", ok, entry)
	}
	slot, ok := reg.LookupSlot(0x04)
	if !ok || slot == nil {
		t.Fatalf("LookupSlot(0x04) ok=%v slot=%v; want non-nil slot (label-only stamping)", ok, slot)
	}
	if slot.Device != nil {
		t.Errorf("orphan slot.Device = %v; want nil", slot.Device)
	}
	if slot.DiscoverySource != DiscoverySourceStaticSeed {
		t.Errorf("orphan slot.DiscoverySource = %v; want DiscoverySourceStaticSeed", slot.DiscoverySource)
	}
}

// TestMarkSlotStaticSeed_RefreshesFacesOnAttachedSlot covers the
// in-scope use case: a slot already attached to a device entry by a
// prior RegisterStaticSeed call has its Role and LastObservedAt
// updated by a subsequent MarkSlotStaticSeed call, AND the
// Faces-refresh branch (slot.Device != nil → syncEntryFacesLocked)
// actually runs — i.e. the device's BusFace list reflects the role.
//
// Test design: uses a target-class address (0x15 BASV2 target). The
// slot is initially seeded with SlotRoleUnknown, then a later
// MarkSlotStaticSeed(0x15, SlotRoleSlave) call sets the explicit
// target role on the slot. After the call, entry.Faces[0x15].Role
// MUST be SlotRoleSlave — that can ONLY hold if syncEntryFacesLocked
// actually ran on the attached entry; if the refresh branch were
// skipped, the cached Face would still carry SlotRoleUnknown.
func TestMarkSlotStaticSeed_RefreshesFacesOnAttachedSlot(t *testing.T) {
	t.Parallel()

	reg := NewDeviceRegistry(nil)
	earlier := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)

	// Step 1: attach slot 0x15 to a device entry via RegisterStaticSeed
	// with NO role hint (SlotRoleUnknown). The slot's Role stays
	// Unknown; entry.Faces[0].Role is Unknown.
	entry := reg.RegisterStaticSeed(DeviceInfo{
		Address:      0x15,
		Manufacturer: "Vaillant",
		DeviceID:     "BASV2",
		SerialNumber: "SN-15",
	}, SlotRoleUnknown, earlier)
	if entry == nil {
		t.Fatalf("RegisterStaticSeed returned nil entry")
	}
	preFace, ok := faceForAddress(entry, 0x15)
	if !ok {
		t.Fatalf("pre-condition: entry has no Face for 0x15")
	}
	if preFace.Role != SlotRoleUnknown {
		t.Fatalf("pre-condition: Face[0x15].Role = %v; want SlotRoleUnknown", preFace.Role)
	}

	// Step 2: MarkSlotStaticSeed on the now-attached slot with an
	// explicit role. Faces refresh must propagate the role into
	// entry.Faces.
	later := earlier.Add(24 * time.Hour)
	reg.MarkSlotStaticSeed(0x15, SlotRoleSlave, later)

	slot, ok := reg.LookupSlot(0x15)
	if !ok || slot == nil {
		t.Fatalf("LookupSlot(0x15) ok=%v slot=%v", ok, slot)
	}
	if slot.Role != SlotRoleSlave {
		t.Errorf("after MarkSlotStaticSeed: slot.Role = %v; want SlotRoleSlave", slot.Role)
	}
	if !slot.LastObservedAt.Equal(later) {
		t.Errorf("after MarkSlotStaticSeed: slot.LastObservedAt = %v; want %v", slot.LastObservedAt, later)
	}

	// Critical assertion: entry.Faces was refreshed by the
	// slot.Device != nil branch — Face[0x15].Role is now SlotRoleSlave
	// (the target role), NOT SlotRoleUnknown. Without Faces-refresh,
	// the Face would still hold SlotRoleUnknown.
	postFace, ok := faceForAddress(entry, 0x15)
	if !ok {
		t.Fatalf("after MarkSlotStaticSeed: entry has no Face for 0x15")
	}
	if postFace.Role != SlotRoleSlave {
		t.Errorf("after MarkSlotStaticSeed: Face[0x15].Role = %v; want SlotRoleSlave (Faces must be refreshed)", postFace.Role)
	}
}

// faceForAddress is a small test helper that returns the BusFace
// matching the given address from a device entry. The DeviceEntry
// public API does not expose Faces directly; we use the *deviceEntry
// concrete type's exported Faces field via the test package's
// same-package access.
func faceForAddress(entry DeviceEntry, addr byte) (BusFace, bool) {
	d, ok := entry.(*deviceEntry)
	if !ok {
		return BusFace{}, false
	}
	for _, face := range d.Faces {
		if face.Addr == addr {
			return face, true
		}
	}
	return BusFace{}, false
}

// TestMarkSlotStaticSeed_RaceFreeWriteAndRead mirrors the M6.1 hardening
// proof for MarkSlotPassiveObserved AND extends it to the cross-API
// case (Codex P3.5 review FINDING_2): concurrent static-seed writers
// and passive-observed writers on the same slot must serialize via
// the registry's RWMutex without panic, deadlock, or torn slot state
// after wg.Wait completes.
//
// Reader scope is intentionally narrow: only LookupSlot's RLock
// (which IS held during the slot map read) and slot != nil. Reading
// slot.Role / DiscoverySource / VerificationState through the
// shared pointer that LookupSlot returns races with writers that
// hold r.mu — that's a pre-existing LookupSlot API surface
// (returns a *AddressSlot rather than a copy) which the existing
// MarkSlotPassiveObserved race test already touches by reading the
// same fields, but it is timing-sensitive. Hardening the API to
// return a snapshot is wider than P3.5's scope. After wg.Wait,
// final-state assertions hold the mutex (via the same
// LookupSlot path) and use the post-wait timing to read steady-state
// fields safely.
func TestMarkSlotStaticSeed_RaceFreeWriteAndRead(t *testing.T) {
	reg := NewDeviceRegistry(nil)
	const writers = 50
	const passiveWriters = 50
	const readers = 100
	const iterations = 100

	var wg sync.WaitGroup
	wg.Add(writers + passiveWriters + readers)

	for i := 0; i < writers; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				reg.MarkSlotStaticSeed(0xF1, SlotRoleMaster, time.Now())
			}
		}()
	}
	// Concurrent passive-observed writers exercise the cross-API
	// serialization path. Both APIs hit r.mu, so neither set of
	// mutations may interleave at the byte level.
	for i := 0; i < passiveWriters; i++ {
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
				_ = slot
				_ = ok
			}
		}()
	}
	wg.Wait()

	slot, ok := reg.LookupSlot(0xF1)
	if !ok || slot == nil {
		t.Fatalf("after race test: LookupSlot(0xF1) ok=%v slot=%v", ok, slot)
	}
	// Final state: at least one StaticSeed write happened and the
	// monotonic guard means DiscoverySource >= StaticSeed
	// (PassiveObserved < StaticSeed).
	if slot.DiscoverySource < DiscoverySourceStaticSeed {
		t.Errorf("after race test: slot.DiscoverySource = %v; want >= StaticSeed", slot.DiscoverySource)
	}
}
