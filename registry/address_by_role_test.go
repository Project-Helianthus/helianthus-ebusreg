package registry

import (
	"testing"
	"time"
)

// Phase C M-C6a: AddressByRole + PrimaryDisplayAddress contract tests.

// TestAddressByRole_BAIPairResolvesSlave asserts the operator-pinned
// scenario: BAI ecoTEC at 0x03↔0x08 has Address()=0x03 (initiator,
// possibly the wrong byte for M2S routing) but
// AddressByRole(SlotRoleSlave)=0x08 returns the correct slave byte.
func TestAddressByRole_BAIPairResolvesSlave(t *testing.T) {
	t.Parallel()

	reg := NewDeviceRegistry(nil)
	// Register slave 0x08 with full identity first.
	reg.Register(DeviceInfo{Address: 0x08, Manufacturer: "Vaillant", DeviceID: "BAI00", SerialNumber: "SN-BAI"})
	// Then alias 0x03 to 0x08 (mimics the inserter flow post-A.7b).
	if err := reg.AliasAddresses(0x03, 0x08); err != nil {
		t.Fatalf("AliasAddresses: %v", err)
	}

	entry, ok := reg.Lookup(0x08)
	if !ok || entry == nil {
		t.Fatalf("Lookup(0x08) ok=%v entry=%v", ok, entry)
	}

	// Mark slot roles so AddressByRole can find them. Active scan
	// path normally does this; the test does it explicitly.
	reg.MarkSlotPassiveObserved(0x08, SlotRoleSlave, time.Time{})
	reg.MarkSlotPassiveObserved(0x03, SlotRoleMaster, time.Time{})

	// Re-lookup so Faces are refreshed.
	entry, _ = reg.Lookup(0x08)

	slave, slaveOK := entry.AddressByRole(SlotRoleSlave)
	if !slaveOK || slave != 0x08 {
		t.Errorf("AddressByRole(SlotRoleSlave) = (0x%02X, %v); want (0x08, true)", slave, slaveOK)
	}
	master, masterOK := entry.AddressByRole(SlotRoleMaster)
	if !masterOK || master != 0x03 {
		t.Errorf("AddressByRole(SlotRoleMaster) = (0x%02X, %v); want (0x03, true)", master, masterOK)
	}
}

// TestAddressByRole_NoMatchReturnsFalse asserts a slave-only device
// (e.g. VR_71 at 0x26 — no canonical master companion in v1 table)
// returns (0, false) when queried for SlotRoleMaster.
func TestAddressByRole_NoMatchReturnsFalse(t *testing.T) {
	t.Parallel()

	reg := NewDeviceRegistry(nil)
	reg.Register(DeviceInfo{Address: 0x26, Manufacturer: "Vaillant", DeviceID: "VR_71", SerialNumber: "SN-VR71"})
	reg.MarkSlotPassiveObserved(0x26, SlotRoleSlave, time.Time{})

	entry, _ := reg.Lookup(0x26)
	if _, ok := entry.AddressByRole(SlotRoleMaster); ok {
		t.Errorf("AddressByRole(SlotRoleMaster) on slave-only device returned ok=true; want false")
	}
	slave, ok := entry.AddressByRole(SlotRoleSlave)
	if !ok || slave != 0x26 {
		t.Errorf("AddressByRole(SlotRoleSlave) = (0x%02X, %v); want (0x26, true)", slave, ok)
	}
}

// TestAddressByRole_FallsBackToAddressClassWhenRoleUnknown asserts the
// active-scan path: registry.Register populates entry.Faces but does
// NOT call MarkSlotPassiveObserved, so face.Role stays SlotRoleUnknown.
// AddressByRole MUST fall back to protocol.AddressClassOf-derived role
// so M2S routing works after plain active scan (Codex P2 from PR #134).
func TestAddressByRole_FallsBackToAddressClassWhenRoleUnknown(t *testing.T) {
	t.Parallel()

	reg := NewDeviceRegistry(nil)
	reg.Register(DeviceInfo{Address: 0x08, Manufacturer: "Vaillant", DeviceID: "BAI00", SerialNumber: "SN-BAI"})
	entry, _ := reg.Lookup(0x08)

	slave, ok := entry.AddressByRole(SlotRoleSlave)
	if !ok || slave != 0x08 {
		t.Errorf("AddressByRole(SlotRoleSlave) on active-scanned slave = (0x%02X, %v); want (0x08, true) via AddressClass fallback", slave, ok)
	}
	if _, ok := entry.AddressByRole(SlotRoleMaster); ok {
		t.Errorf("AddressByRole(SlotRoleMaster) on slave-only entry returned ok=true; want false")
	}
}

// TestPrimaryDisplayAddress_MatchesAddress asserts the new method is
// equivalent to Address() for now (M-C6a is additive; M-C6c removes
// Address). Same primary returned.
func TestPrimaryDisplayAddress_MatchesAddress(t *testing.T) {
	t.Parallel()

	reg := NewDeviceRegistry(nil)
	reg.Register(DeviceInfo{Address: 0x15, Manufacturer: "Vaillant", DeviceID: "BASV2"})
	entry, _ := reg.Lookup(0x15)

	if got, want := entry.PrimaryDisplayAddress(), entry.Address(); got != want {
		t.Errorf("PrimaryDisplayAddress() = 0x%02X; want %02X (matches Address)", got, want)
	}
}
