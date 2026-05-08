package registry

import "testing"

// TestAliasAddresses_PreservesSecondaryIdentity asserts that when
// AliasAddresses(a, b) is called and slotA.Device is an empty
// passive-observed entry while slotB.Device is identity-bearing
// (manufacturer + deviceID + serialNumber set), the merged entry
// retains the identity from the secondary side.
//
// Regression scenario from live observation 2026-05-08: BASV2
// actively scanned at 0x15 with full identity, then bus traffic
// from BASV2 master 0x10 caused the inserter to call
// Register({Address: 0x10}) creating an empty entry, followed by
// AliasAddresses(0x10, 0x15). Pre-fix, this destroyed BASV2's
// identity row (delete(r.identity, key) + removeEntry(r.order,
// secondary)) and the canonical entry retained empty info.
//
// Post-fix (P0): absorbIdentityLocked promotes secondary's
// non-empty fields onto canonical, then secondary is removed
// without losing the identity.
func TestAliasAddresses_PreservesSecondaryIdentity(t *testing.T) {
	t.Parallel()

	reg := NewDeviceRegistry(nil)

	// Step 1: register the identity-bearing entry at 0x15 (active scan path).
	reg.Register(DeviceInfo{
		Address:         0x15,
		Manufacturer:    "Vaillant",
		DeviceID:        "BASV2",
		SerialNumber:    "SN-BASV2-001",
		SoftwareVersion: "0204",
		HardwareVersion: "0102",
	})

	// Step 2: register a bare passive observation at 0x10 (gateway
	// inserter path: Register({Address: 0x10}) with no identity).
	reg.Register(DeviceInfo{Address: 0x10})

	// Step 3: alias the canonical pair (gateway calls
	// AliasAddresses(initiator, target)).
	if err := reg.AliasAddresses(0x10, 0x15); err != nil {
		t.Fatalf("AliasAddresses(0x10, 0x15) error = %v", err)
	}

	// Both lookups must return the same entry with full identity.
	for _, addr := range []byte{0x10, 0x15} {
		entry, ok := reg.Lookup(addr)
		if !ok {
			t.Fatalf("Lookup(0x%02X) = false; want entry", addr)
		}
		if got := entry.Manufacturer(); got != "Vaillant" {
			t.Errorf("Lookup(0x%02X).Manufacturer() = %q; want \"Vaillant\"", addr, got)
		}
		if got := entry.DeviceID(); got != "BASV2" {
			t.Errorf("Lookup(0x%02X).DeviceID() = %q; want \"BASV2\"", addr, got)
		}
		if got := entry.SerialNumber(); got != "SN-BASV2-001" {
			t.Errorf("Lookup(0x%02X).SerialNumber() = %q; want \"SN-BASV2-001\"", addr, got)
		}
		if got := entry.SoftwareVersion(); got != "0204" {
			t.Errorf("Lookup(0x%02X).SoftwareVersion() = %q; want \"0204\"", addr, got)
		}
		if got := entry.HardwareVersion(); got != "0102" {
			t.Errorf("Lookup(0x%02X).HardwareVersion() = %q; want \"0102\"", addr, got)
		}
	}

	// Both addresses must appear in the merged entry's address set.
	entry, _ := reg.Lookup(0x10)
	addresses := entry.Addresses()
	have10, have15 := false, false
	for _, a := range addresses {
		if a == 0x10 {
			have10 = true
		}
		if a == 0x15 {
			have15 = true
		}
	}
	if !have10 || !have15 {
		t.Errorf("merged Addresses() = %v; want both 0x10 and 0x15", addresses)
	}

	// Identity-by-SerialNumber lookup must still resolve to the
	// merged entry (not orphaned by the previous identityKey
	// deletion).
	by, ok := reg.lookupByIdentity(DeviceInfo{Manufacturer: "Vaillant", DeviceID: "BASV2", SerialNumber: "SN-BASV2-001"})
	if !ok {
		t.Fatalf("lookupByIdentity by SerialNumber = false; want entry")
	}
	if by.Manufacturer() != "Vaillant" || by.SerialNumber() != "SN-BASV2-001" {
		t.Errorf("lookupByIdentity returned unexpected entry: mfr=%q serial=%q", by.Manufacturer(), by.SerialNumber())
	}

	// Iterate must show exactly one entry (the merged one).
	count := 0
	reg.Iterate(func(e DeviceEntry) bool {
		count++
		return true
	})
	if count != 1 {
		t.Errorf("registry entry count = %d; want 1 (merged BASV2)", count)
	}
}

// TestAliasAddresses_PreservesCanonicalIdentity asserts the
// symmetric case: canonical (slotA) is identity-bearing, secondary
// (slotB) is empty. The original behavior already worked here
// because the empty secondary has nothing to lose, but the test
// pins the invariant against future regressions.
func TestAliasAddresses_PreservesCanonicalIdentity(t *testing.T) {
	t.Parallel()

	reg := NewDeviceRegistry(nil)

	reg.Register(DeviceInfo{
		Address:      0x08,
		Manufacturer: "Vaillant",
		DeviceID:     "BAI00",
		SerialNumber: "SN-BAI-001",
	})
	reg.Register(DeviceInfo{Address: 0x03})

	if err := reg.AliasAddresses(0x08, 0x03); err != nil {
		t.Fatalf("AliasAddresses(0x08, 0x03) error = %v", err)
	}

	entry, ok := reg.Lookup(0x03)
	if !ok {
		t.Fatalf("Lookup(0x03) = false; want entry")
	}
	if entry.Manufacturer() != "Vaillant" || entry.DeviceID() != "BAI00" || entry.SerialNumber() != "SN-BAI-001" {
		t.Errorf("Lookup(0x03) lost identity: mfr=%q devID=%q serial=%q",
			entry.Manufacturer(), entry.DeviceID(), entry.SerialNumber())
	}
}

// TestAliasAddresses_BothEmpty exercises the no-identity case
// (e.g. NETX3 0xF1↔0xF6 both passively observed before any active
// scan succeeds for either face). The merge must still group the
// addresses; identity stays empty until a future Register or
// enrichment populates it.
func TestAliasAddresses_BothEmpty(t *testing.T) {
	t.Parallel()

	reg := NewDeviceRegistry(nil)
	reg.Register(DeviceInfo{Address: 0xF1})
	reg.Register(DeviceInfo{Address: 0xF6})

	if err := reg.AliasAddresses(0xF1, 0xF6); err != nil {
		t.Fatalf("AliasAddresses(0xF1, 0xF6) error = %v", err)
	}

	for _, addr := range []byte{0xF1, 0xF6} {
		entry, ok := reg.Lookup(addr)
		if !ok {
			t.Fatalf("Lookup(0x%02X) = false; want entry", addr)
		}
		if entry.Manufacturer() != "" {
			t.Errorf("Lookup(0x%02X).Manufacturer() = %q; want empty (no identity yet)", addr, entry.Manufacturer())
		}
	}

	// Subsequent active scan registers identity for one face;
	// identity-merge through the canonical-pair alias should
	// propagate it onto the merged entry.
	reg.Register(DeviceInfo{
		Address:      0xF6,
		Manufacturer: "Vaillant",
		DeviceID:     "NETX3",
		SerialNumber: "SN-NETX3-001",
	})
	entry, _ := reg.Lookup(0xF1)
	if entry.Manufacturer() != "Vaillant" {
		t.Errorf("after Register(0xF6, NETX3), Lookup(0xF1).Manufacturer() = %q; want \"Vaillant\"", entry.Manufacturer())
	}
	if entry.SerialNumber() != "SN-NETX3-001" {
		t.Errorf("after Register(0xF6, NETX3), Lookup(0xF1).SerialNumber() = %q; want \"SN-NETX3-001\"", entry.SerialNumber())
	}
}
