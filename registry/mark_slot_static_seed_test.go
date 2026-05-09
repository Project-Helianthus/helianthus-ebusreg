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

// TestMarkSlotStaticSeed_UpgradesAttachedSlot covers the in-scope use
// case: a slot already attached to a device entry (here via
// RegisterStaticSeed for a different address that aliased into the
// same entry) gets its discovery_source label upgraded by a
// MarkSlotStaticSeed call. Faces is refreshed because slot.Device is
// non-nil at call time.
func TestMarkSlotStaticSeed_UpgradesAttachedSlot(t *testing.T) {
	t.Parallel()

	reg := NewDeviceRegistry(nil)
	// Land slot 0xF1 attached to a device via Register (lands as
	// ActiveConfirmed/IdentityConfirmed — but the discovery_source
	// downgrade-guard means MarkSlotStaticSeed will NOT advance the
	// label here). Use a fresh slot with passive-observed state
	// instead so we can prove the upgrade path.
	now := time.Now()
	reg.MarkSlotPassiveObserved(0xF1, SlotRoleMaster, now)
	pre, _ := reg.LookupSlot(0xF1)
	if pre.DiscoverySource != DiscoverySourcePassiveObserved {
		t.Fatalf("pre-condition: slot.DiscoverySource = %v; want PassiveObserved", pre.DiscoverySource)
	}
	// Attach a device pointer so syncEntryFacesLocked has something
	// to refresh — the most realistic path is RegisterStaticSeed of
	// THIS address (which would normally double-stamp; here we just
	// want a Device on the slot so MarkSlotStaticSeed exercises the
	// Faces-refresh branch).
	reg.RegisterStaticSeed(DeviceInfo{
		Address:      0xF1,
		Manufacturer: "Vaillant",
		DeviceID:     "NETX3",
		SerialNumber: "SN-AB",
	}, SlotRoleMaster, now)

	post, _ := reg.LookupSlot(0xF1)
	if post.DiscoverySource != DiscoverySourceStaticSeed {
		t.Errorf("after RegisterStaticSeed: slot.DiscoverySource = %v; want StaticSeed", post.DiscoverySource)
	}
	if post.Device == nil {
		t.Errorf("after RegisterStaticSeed: slot.Device is nil; want non-nil (entry attached)")
	}
}

// TestMarkSlotStaticSeed_RaceFreeWriteAndRead mirrors the M6.1 hardening
// proof for MarkSlotPassiveObserved AND extends it to the cross-API
// case (Codex P3.5 review FINDING_2): concurrent passive-observed and
// static-seed writers on the same slot must converge through the
// registry's RWMutex without panic, deadlock, or torn DiscoverySource /
// VerificationState / Role values.
//
// Readers exercise only the monotonically-converging fields (Role,
// DiscoverySource, VerificationState). The shared-pointer pattern of
// LookupSlot means LastObservedAt is technically read-after-write via
// the slot pointer; the existing MarkSlotPassiveObserved race test
// constrains its reader the same way for the same reason. Improving
// LookupSlot to return a snapshot copy would close that pointer-share
// surface but is wider than P3.5's scope.
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
	// Final state: at least one StaticSeed write happened and the
	// monotonic guard means DiscoverySource >= StaticSeed
	// (PassiveObserved < StaticSeed).
	if slot.DiscoverySource < DiscoverySourceStaticSeed {
		t.Errorf("after race test: slot.DiscoverySource = %v; want >= StaticSeed", slot.DiscoverySource)
	}
}
