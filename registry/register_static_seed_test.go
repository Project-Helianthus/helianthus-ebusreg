package registry

import (
	"testing"
	"time"
)

// TestRegisterStaticSeed_StampsCorrectLabels asserts that
// RegisterStaticSeed stamps the address slot with
// DiscoverySourceStaticSeed / VerificationStateCandidate (NOT the
// ActiveConfirmed / IdentityConfirmed labels that Register uses).
// This is the central P3.5 contract.
func TestRegisterStaticSeed_StampsCorrectLabels(t *testing.T) {
	t.Parallel()

	reg := NewDeviceRegistry(nil)
	now := time.Date(2026, 5, 9, 12, 0, 0, 0, time.UTC)
	entry := reg.RegisterStaticSeed(DeviceInfo{
		Address:      0xF1,
		Manufacturer: "Vaillant",
		DeviceID:     "NETX3",
	}, SlotRoleMaster, now)
	if entry == nil {
		t.Fatalf("RegisterStaticSeed returned nil entry")
	}

	slot, ok := reg.LookupSlot(0xF1)
	if !ok || slot == nil {
		t.Fatalf("LookupSlot(0xF1) ok=%v slot=%v", ok, slot)
	}
	if slot.DiscoverySource != DiscoverySourceStaticSeed {
		t.Errorf("slot.DiscoverySource = %v; want DiscoverySourceStaticSeed", slot.DiscoverySource)
	}
	if slot.VerificationState != VerificationStateCandidate {
		t.Errorf("slot.VerificationState = %v; want VerificationStateCandidate", slot.VerificationState)
	}
	if slot.Role != SlotRoleMaster {
		t.Errorf("slot.Role = %v; want SlotRoleMaster", slot.Role)
	}
}

// TestRegisterStaticSeed_AttachesIdentity asserts the underlying
// identity-merge body still runs — the entry has the seeded
// Manufacturer / DeviceID and is queryable via Lookup.
func TestRegisterStaticSeed_AttachesIdentity(t *testing.T) {
	t.Parallel()

	reg := NewDeviceRegistry(nil)
	reg.RegisterStaticSeed(DeviceInfo{
		Address:      0xF1,
		Manufacturer: "Vaillant",
		DeviceID:     "NETX3",
	}, SlotRoleMaster, time.Now())

	entry, ok := reg.Lookup(0xF1)
	if !ok || entry == nil {
		t.Fatalf("Lookup(0xF1) ok=%v entry=%v", ok, entry)
	}
	if got, want := entry.Manufacturer(), "Vaillant"; got != want {
		t.Errorf("entry.Manufacturer() = %q; want %q", got, want)
	}
	if got, want := entry.DeviceID(), "NETX3"; got != want {
		t.Errorf("entry.DeviceID() = %q; want %q", got, want)
	}

	slot, _ := reg.LookupSlot(0xF1)
	if slot.Device == nil {
		t.Errorf("slot.Device is nil; want pointer to the registered entry")
	}
}

// TestRegisterStaticSeed_DoesNotDowngradeActiveConfirmed asserts that
// calling RegisterStaticSeed on a slot already marked ActiveConfirmed
// (e.g. via a prior Register call) does NOT downgrade the discovery
// label. The identity-merge body MAY still update fields per the
// Register monotonic merge rules, but the AddressSlot label stays at
// ActiveConfirmed / IdentityConfirmed.
func TestRegisterStaticSeed_DoesNotDowngradeActiveConfirmed(t *testing.T) {
	t.Parallel()

	reg := NewDeviceRegistry(nil)
	// Land at ActiveConfirmed first.
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

	// Now seed; must not downgrade.
	reg.RegisterStaticSeed(DeviceInfo{
		Address:      0xF1,
		Manufacturer: "Vaillant",
		DeviceID:     "NETX3",
	}, SlotRoleMaster, time.Now())

	post, _ := reg.LookupSlot(0xF1)
	if post.DiscoverySource != DiscoverySourceActiveConfirmed {
		t.Errorf("slot.DiscoverySource after RegisterStaticSeed = %v; want ActiveConfirmed (no downgrade)", post.DiscoverySource)
	}
	if post.VerificationState != VerificationStateIdentityConfirmed {
		t.Errorf("slot.VerificationState after RegisterStaticSeed = %v; want IdentityConfirmed (no downgrade)", post.VerificationState)
	}
}

// TestRegisterStaticSeed_PreservesAliasIdentity verifies that the P0
// alias-identity invariant (PR #136 — secondary keys retained via
// identityKeyAliases on a Register call that promotes a different
// identityKey) holds when seed registration is later followed by an
// AliasAddresses call. The seed's DeviceID and Manufacturer must
// survive the alias.
//
// Uses distinct SerialNumber values per face so each entry has a
// non-empty identityKey — that's the precondition for AliasAddresses
// to actually exercise the identityKeyAliases transfer path
// (per Codex P3.5 review FINDING_1 — without serial/mac the
// canonicalPhysicalIdentity yields no key and the alias path silently
// no-ops).
func TestRegisterStaticSeed_PreservesAliasIdentity(t *testing.T) {
	t.Parallel()

	reg := NewDeviceRegistry(nil)
	// Seed both NETX3 face addresses with distinct stable identity
	// keys so AliasAddresses takes the identityKey-merge path.
	reg.RegisterStaticSeed(DeviceInfo{
		Address:      0xF1,
		Manufacturer: "Vaillant",
		DeviceID:     "NETX3",
		SerialNumber: "SN-F1",
	}, SlotRoleMaster, time.Now())
	reg.RegisterStaticSeed(DeviceInfo{
		Address:      0xF6,
		Manufacturer: "Vaillant",
		DeviceID:     "NETX3",
		SerialNumber: "SN-F6",
	}, SlotRoleSlave, time.Now())

	if err := reg.AliasAddresses(0xF1, 0xF6); err != nil {
		t.Fatalf("AliasAddresses(0xF1, 0xF6) error = %v", err)
	}

	// Lookup via either alias must resolve to the merged entry with
	// the seeded identity preserved.
	for _, addr := range []byte{0xF1, 0xF6} {
		entry, ok := reg.Lookup(addr)
		if !ok || entry == nil {
			t.Fatalf("Lookup(0x%02X) after alias: ok=%v entry=%v", addr, ok, entry)
		}
		if got, want := entry.DeviceID(), "NETX3"; got != want {
			t.Errorf("aliased entry[0x%02X] DeviceID = %q; want %q", addr, got, want)
		}
		if got, want := entry.Manufacturer(), "Vaillant"; got != want {
			t.Errorf("aliased entry[0x%02X] Manufacturer = %q; want %q", addr, got, want)
		}
	}

	// Both slots still labelled static_seed/candidate.
	for _, addr := range []byte{0xF1, 0xF6} {
		slot, _ := reg.LookupSlot(addr)
		if slot.DiscoverySource != DiscoverySourceStaticSeed {
			t.Errorf("slot[0x%02X].DiscoverySource = %v; want StaticSeed", addr, slot.DiscoverySource)
		}
	}
}
