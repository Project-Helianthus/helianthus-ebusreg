package registry

import (
	"testing"
	"time"
)

// P8 — RegisterPassiveObserved correct-label contract.
//
// The gateway's address_table_inserter previously called Register
// (which stamps DiscoverySourceActiveConfirmed /
// VerificationStateIdentityConfirmed) followed by
// MarkSlotPassiveObserved. That ordering recorded the wrong original source,
// because discovery provenance is immutable once the face is admitted.
//
// RegisterPassiveObserved performs identity-merge AND the passive-
// observation slot stamping atomically under one lock acquisition.

// TestRegisterPassiveObserved_StampsCorrectLabels asserts that
// RegisterPassiveObserved stamps PassiveObserved / Corroborated
// — NOT the ActiveConfirmed / IdentityConfirmed labels Register uses.
// This is the central P8 contract.
func TestRegisterPassiveObserved_StampsCorrectLabels(t *testing.T) {
	t.Parallel()

	reg := NewDeviceRegistry(nil)
	now := time.Date(2026, 5, 9, 12, 0, 0, 0, time.UTC)
	entry := reg.RegisterPassiveObserved(DeviceInfo{
		Address: 0xF1,
	}, SlotRoleMaster, now)
	if entry == nil {
		t.Fatalf("RegisterPassiveObserved returned nil entry")
	}

	slot, ok := reg.LookupSlot(0xF1)
	if !ok || slot == nil {
		t.Fatalf("LookupSlot(0xF1) ok=%v slot=%v", ok, slot)
	}
	if slot.DiscoverySource != DiscoverySourcePassiveObserved {
		t.Errorf("slot.DiscoverySource = %v; want DiscoverySourcePassiveObserved", slot.DiscoverySource)
	}
	if slot.VerificationState != VerificationStateCorroborated {
		t.Errorf("slot.VerificationState = %v; want VerificationStateCorroborated", slot.VerificationState)
	}
	if slot.Role != SlotRoleMaster {
		t.Errorf("slot.Role = %v; want SlotRoleMaster", slot.Role)
	}
	if !slot.FirstObservedAt.Equal(now) {
		t.Errorf("slot.FirstObservedAt = %v; want %v", slot.FirstObservedAt, now)
	}
	if !slot.LastObservedAt.Equal(now) {
		t.Errorf("slot.LastObservedAt = %v; want %v", slot.LastObservedAt, now)
	}
}

// TestRegisterPassiveObserved_AttachesEntryToSlot asserts that the
// slot's Device pointer is set after RegisterPassiveObserved — the
// passive insert must be queryable via Lookup, just like Register.
// Without slot.Device the address-table snapshot would emit nil
// device fields and break operator visibility.
func TestRegisterPassiveObserved_AttachesEntryToSlot(t *testing.T) {
	t.Parallel()

	reg := NewDeviceRegistry(nil)
	reg.RegisterPassiveObserved(DeviceInfo{
		Address: 0xF1,
	}, SlotRoleMaster, time.Now())

	slot, _ := reg.LookupSlot(0xF1)
	if slot.Device == nil {
		t.Errorf("slot.Device is nil; want pointer to the registered entry")
	}

	entry, ok := reg.Lookup(0xF1)
	if !ok || entry == nil {
		t.Fatalf("Lookup(0xF1) ok=%v entry=%v", ok, entry)
	}
}

// TestRegisterPassiveObserved_RetainsExistingActiveSource asserts that a
// passive registration preserves an earlier active admission source.
func TestRegisterPassiveObserved_RetainsExistingActiveSource(t *testing.T) {
	t.Parallel()

	reg := NewDeviceRegistry(nil)
	reg.Register(DeviceInfo{
		Address:      0xF1,
		Manufacturer: "Vaillant",
		DeviceID:     "NETX3",
		SerialNumber: "SN-Z",
	})
	pre, _ := reg.LookupSlot(0xF1)
	if pre.DiscoverySource != DiscoverySourceActiveConfirmed {
		t.Fatalf("pre-condition: slot.DiscoverySource = %v; want ActiveConfirmed", pre.DiscoverySource)
	}

	reg.RegisterPassiveObserved(DeviceInfo{
		Address: 0xF1,
	}, SlotRoleMaster, time.Now())

	post, _ := reg.LookupSlot(0xF1)
	if post.DiscoverySource != DiscoverySourceActiveConfirmed {
		t.Errorf("slot.DiscoverySource after RegisterPassiveObserved = %v; want ActiveConfirmed (no downgrade)", post.DiscoverySource)
	}
	if post.VerificationState != VerificationStateIdentityConfirmed {
		t.Errorf("slot.VerificationState after RegisterPassiveObserved = %v; want IdentityConfirmed (no downgrade)", post.VerificationState)
	}
}

// TestRegisterPassiveObserved_FollowedByActiveScanConfirmsIdentity asserts
// that active evidence advances verification while retaining passive origin.
func TestRegisterPassiveObserved_FollowedByActiveScanConfirmsIdentity(t *testing.T) {
	t.Parallel()

	reg := NewDeviceRegistry(nil)
	reg.RegisterPassiveObserved(DeviceInfo{
		Address: 0x15,
	}, SlotRoleSlave, time.Now())
	pre, _ := reg.LookupSlot(0x15)
	if pre.DiscoverySource != DiscoverySourcePassiveObserved {
		t.Fatalf("pre-condition: slot.DiscoverySource = %v; want PassiveObserved", pre.DiscoverySource)
	}

	reg.Register(DeviceInfo{
		Address:      0x15,
		Manufacturer: "Vaillant",
		DeviceID:     "BASV2",
		SerialNumber: "SN-X",
	})

	post, _ := reg.LookupSlot(0x15)
	if post.DiscoverySource != DiscoverySourcePassiveObserved {
		t.Errorf("slot.DiscoverySource after directed scan = %v; want retained PassiveObserved", post.DiscoverySource)
	}
	if post.VerificationState != VerificationStateIdentityConfirmed {
		t.Errorf("slot.VerificationState after directed scan = %v; want IdentityConfirmed", post.VerificationState)
	}
}

// TestRegisterPassiveObserved_StaticSeedAfterPassiveRetainsSource asserts
// that later taxonomy knowledge does not rewrite passive origin.
func TestRegisterPassiveObserved_StaticSeedAfterPassiveRetainsSource(t *testing.T) {
	t.Parallel()

	reg := NewDeviceRegistry(nil)
	reg.RegisterPassiveObserved(DeviceInfo{
		Address: 0xF1,
	}, SlotRoleMaster, time.Now())

	reg.MarkSlotStaticSeed(0xF1, SlotRoleMaster, time.Now())

	slot, _ := reg.LookupSlot(0xF1)
	if slot.DiscoverySource != DiscoverySourcePassiveObserved {
		t.Errorf("slot.DiscoverySource = %v; want retained PassiveObserved", slot.DiscoverySource)
	}
}

// TestRegisterPassiveObserved_RaceFreeWriteAndRead exercises the lock
// discipline: concurrent RegisterPassiveObserved calls and Lookup /
// LookupSlot reads must not produce torn state. Race detector is the
// authoritative gate.
func TestRegisterPassiveObserved_RaceFreeWriteAndRead(t *testing.T) {
	t.Parallel()

	reg := NewDeviceRegistry(nil)
	done := make(chan struct{})
	go func() {
		for i := 0; i < 100; i++ {
			reg.RegisterPassiveObserved(DeviceInfo{
				Address: byte(0x10 + (i % 16)),
			}, SlotRoleMaster, time.Now())
		}
		close(done)
	}()
	for i := 0; i < 200; i++ {
		_, _ = reg.LookupSlot(byte(0x10 + (i % 16)))
		_, _ = reg.Lookup(byte(0x10 + (i % 16)))
	}
	<-done
}

// TestRegisterPassiveObserved_IdentityMergePreservesPassiveLabels
// covers Codex P8 review NIT FINDING_2: when two passively-observed
// addresses share a stable identity (SerialNumber / MacAddress), the
// underlying registerLocked merges them into a single device entry
// with both addresses, AND both slots stay labelled at
// PassiveObserved/Corroborated (NOT promoted to ActiveConfirmed).
// Mirrors TestRegisterStaticSeed_PreservesAliasIdentity for the
// passive-observed path.
func TestRegisterPassiveObserved_IdentityMergePreservesPassiveLabels(t *testing.T) {
	t.Parallel()

	reg := NewDeviceRegistry(nil)
	now := time.Now()
	reg.RegisterPassiveObserved(DeviceInfo{
		Address:      0xF1,
		Manufacturer: "Vaillant",
		DeviceID:     "NETX3",
		SerialNumber: "SN-MERGE",
	}, SlotRoleMaster, now)
	reg.RegisterPassiveObserved(DeviceInfo{
		Address:      0xF6,
		Manufacturer: "Vaillant",
		DeviceID:     "NETX3",
		SerialNumber: "SN-MERGE",
	}, SlotRoleSlave, now)

	// Identity merge: both addresses must resolve to the same entry.
	entry1, ok1 := reg.Lookup(0xF1)
	entry2, ok2 := reg.Lookup(0xF6)
	if !ok1 || !ok2 {
		t.Fatalf("Lookup(0xF1)=%v Lookup(0xF6)=%v; both must resolve", ok1, ok2)
	}
	if entry1.SerialNumber() != "SN-MERGE" || entry2.SerialNumber() != "SN-MERGE" {
		t.Errorf("merged entries SerialNumber should be SN-MERGE; got %q / %q", entry1.SerialNumber(), entry2.SerialNumber())
	}
	// Codex P8 pass 2 NIT — assert SAME entry (alias), not just two
	// entries with the same serial. The merged entry's Addresses()
	// must contain both 0xF1 and 0xF6.
	addrs1 := entry1.Addresses()
	addrs2 := entry2.Addresses()
	gotF1, gotF6 := false, false
	for _, a := range addrs1 {
		if a == 0xF1 {
			gotF1 = true
		}
		if a == 0xF6 {
			gotF6 = true
		}
	}
	if !gotF1 || !gotF6 {
		t.Errorf("entry1.Addresses() = %v; want both 0xF1 and 0xF6 (single merged entry)", addrs1)
	}
	if len(addrs1) != len(addrs2) {
		t.Errorf("entry1.Addresses() len=%d != entry2.Addresses() len=%d (different entries — merge failed)", len(addrs1), len(addrs2))
	}

	// Both slots stay at passive_observed/corroborated — the identity
	// merge must NOT silently promote either slot to active_confirmed.
	for _, addr := range []byte{0xF1, 0xF6} {
		slot, ok := reg.LookupSlot(addr)
		if !ok {
			t.Fatalf("LookupSlot(0x%02X) ok=%v", addr, ok)
		}
		if slot.DiscoverySource != DiscoverySourcePassiveObserved {
			t.Errorf("slot[0x%02X].DiscoverySource = %v; want PassiveObserved", addr, slot.DiscoverySource)
		}
		if slot.VerificationState != VerificationStateCorroborated {
			t.Errorf("slot[0x%02X].VerificationState = %v; want Corroborated", addr, slot.VerificationState)
		}
	}
}

// TestMarkSlotPassiveObserved_StillStampsLabel proves the refactor
// preserved MarkSlotPassiveObserved's existing behaviour — calls
// against a pre-attached slot still mark PassiveObserved /
// Corroborated. (Direct callers exist; the refactor must not break
// them.)
func TestMarkSlotPassiveObserved_StillStampsLabel(t *testing.T) {
	t.Parallel()

	reg := NewDeviceRegistry(nil)
	reg.MarkSlotPassiveObserved(0x10, SlotRoleMaster, time.Now())

	slot, ok := reg.LookupSlot(0x10)
	if !ok {
		t.Fatalf("LookupSlot(0x10) ok=%v", ok)
	}
	if slot.DiscoverySource != DiscoverySourcePassiveObserved {
		t.Errorf("slot.DiscoverySource = %v; want PassiveObserved", slot.DiscoverySource)
	}
	if slot.VerificationState != VerificationStateCorroborated {
		t.Errorf("slot.VerificationState = %v; want Corroborated", slot.VerificationState)
	}
}
