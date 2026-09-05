package registry

import (
	"reflect"
	"testing"
)

// TestAliasAddresses_PreservesSecondaryIdentity asserts that when
// AliasAddresses(a, b) is called and slotA.Device is an empty
// passive-observed entry while slotB.Device is identity-bearing
// (manufacturer + deviceID + serialNumber set), the merged entry
// retains the identity from the secondary side.
//
// Regression scenario from live observation 2026-05-08: BASV2
// actively scanned at 0x15 with full identity, then bus traffic
// from BASV2 initiator 0x10 caused the inserter to call
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

// TestAliasAddresses_PreservesDistinctIdentityKeys asserts that explicit
// aliasing retains both qualified identity keys as lookup aliases.
func TestAliasAddresses_PreservesDistinctIdentityKeys(t *testing.T) {
	t.Parallel()

	reg := NewDeviceRegistry(nil)

	// Use two complete but distinct triples. Explicit AliasAddresses owns this
	// topology decision; Register itself never joins these addresses.
	reg.Register(DeviceInfo{
		Address:      0x10,
		Manufacturer: "Vaillant",
		DeviceID:     "BASV2",
		SerialNumber: "SN-CANONICAL-001",
		MacAddress:   "AA:BB:CC:DD:EE:01",
	})
	reg.Register(DeviceInfo{
		Address:      0x15,
		Manufacturer: "Vaillant",
		DeviceID:     "BASV2",
		SerialNumber: "SN-DISTINCT-001",
	})

	if err := reg.AliasAddresses(0x10, 0x15); err != nil {
		t.Fatalf("AliasAddresses error = %v", err)
	}

	// Both qualified lookup paths must resolve.
	byCanonical, ok := reg.lookupByIdentity(DeviceInfo{Manufacturer: "Vaillant", DeviceID: "BASV2", SerialNumber: "SN-CANONICAL-001"})
	if !ok {
		t.Fatalf("lookupByIdentity by canonical triple = false; want resolvable")
	}
	bySerial, ok := reg.lookupByIdentity(DeviceInfo{Manufacturer: "Vaillant", DeviceID: "BASV2", SerialNumber: "SN-DISTINCT-001"})
	if !ok {
		t.Fatalf("lookupByIdentity by Serial = false; want resolvable")
	}
	// Both must point at the same entry.
	if byCanonical.MacAddress() != "AA:BB:CC:DD:EE:01" {
		t.Errorf("canonical resolution: mac=%q; want AA:BB:CC:DD:EE:01", byCanonical.MacAddress())
	}
	if bySerial != byCanonical {
		t.Errorf("qualified alias resolution differs: canonical=%p serial=%p", byCanonical, bySerial)
	}
	// Same set of addresses (canonical + secondary's address).
	if !reflect.DeepEqual(byCanonical.Addresses(), bySerial.Addresses()) {
		t.Errorf("merge mismatch: canonical.Addresses=%v bySerial.Addresses=%v", byCanonical.Addresses(), bySerial.Addresses())
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

// TestAliasAddresses_PreservesSurvivingSecondaryIdentity asserts the
// case where the secondary entry has MORE THAN the aliased address
// (e.g. 0x15 + 0x16 share serial; aliasing empty 0x10 to 0x15
// preserves the secondary at 0x16 with intact identity row in
// r.identity). This is the Codex P2 follow-up on PR #136 (2026-05-08):
// pre-fix, absorbIdentityLocked moved identityKey to canonical and
// cleared secondary's, leaving the surviving-secondary entry without
// an identity row → lookupByIdentity could not resolve to it.
func TestAliasAddresses_PreservesSurvivingSecondaryIdentity(t *testing.T) {
	t.Parallel()

	reg := NewDeviceRegistry(nil)

	// Register a multi-face entry: 0x15 + 0x16 share Serial.
	reg.Register(DeviceInfo{
		Address:      0x15,
		Manufacturer: "Vaillant",
		DeviceID:     "BASV2",
		SerialNumber: "SN-MULTI-001",
	})
	reg.Register(DeviceInfo{
		Address:      0x16,
		Manufacturer: "Vaillant",
		DeviceID:     "BASV2",
		SerialNumber: "SN-MULTI-001",
	})
	// At this point: identity-merge collapsed 0x15 + 0x16 into one
	// entry with addresses=[0x15, 0x16].

	// Plant an empty entry at 0x10.
	reg.Register(DeviceInfo{Address: 0x10})

	// Alias 0x10 ↔ 0x15. Pre-fix this incorrectly stripped the
	// (0x15, 0x16)-merged entry's identityKey because absorb fired
	// before the addresses[]==nil check.
	if err := reg.AliasAddresses(0x10, 0x15); err != nil {
		t.Fatalf("AliasAddresses error = %v", err)
	}

	// Lookup by serial must still resolve. The entry it resolves to
	// is implementation-specific (could be the merged-onto-canonical
	// or the surviving-secondary), but identity must be intact.
	by, ok := reg.lookupByIdentity(DeviceInfo{Manufacturer: "Vaillant", DeviceID: "BASV2", SerialNumber: "SN-MULTI-001"})
	if !ok {
		t.Fatalf("lookupByIdentity by SN-MULTI-001 = false; want resolvable")
	}
	if by.Manufacturer() != "Vaillant" || by.SerialNumber() != "SN-MULTI-001" {
		t.Errorf("identity lost: mfr=%q serial=%q", by.Manufacturer(), by.SerialNumber())
	}
	// 0x16 must remain reachable via direct lookup with same identity.
	entry16, ok := reg.Lookup(0x16)
	if !ok {
		t.Fatalf("Lookup(0x16) = false; want entry")
	}
	if entry16.SerialNumber() != "SN-MULTI-001" {
		t.Errorf("0x16 lost SerialNumber: got %q; want SN-MULTI-001", entry16.SerialNumber())
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

// TestRegister_PartialIdentityDoesNotMergeAcrossAddresses keeps model-only
// observations out of the qualified identity index, even after one address
// is later enriched with a complete triple.
func TestRegister_PartialIdentityDoesNotMergeAcrossAddresses(t *testing.T) {
	t.Parallel()

	reg := NewDeviceRegistry(nil)

	// Step 1: register entry A with model evidence only.
	reg.Register(DeviceInfo{
		Address:         0x10,
		Manufacturer:    "Vaillant",
		DeviceID:        "BAI00",
		SoftwareVersion: "0204",
		HardwareVersion: "0102",
	})

	// Step 2: refresh entry A with a complete triple.
	reg.Register(DeviceInfo{
		Address:         0x10,
		Manufacturer:    "Vaillant",
		DeviceID:        "BAI00",
		SerialNumber:    "SN-A-001",
		SoftwareVersion: "0204",
		HardwareVersion: "0102",
	})

	// Step 3: register entry B with the same model but a different serial.
	reg.Register(DeviceInfo{
		Address:         0x11,
		Manufacturer:    "Vaillant",
		DeviceID:        "BAI00",
		SerialNumber:    "SN-B-001",
		SoftwareVersion: "0204",
		HardwareVersion: "0102",
	})

	// Step 4: a bare model-only observation at a new address cannot join
	// either completed triple.
	reg.Register(DeviceInfo{
		Address:         0x12,
		Manufacturer:    "Vaillant",
		DeviceID:        "BAI00",
		SoftwareVersion: "0204",
		HardwareVersion: "0102",
	})

	// A's and B's address groups must remain separate from 0x12.
	entryA, _ := reg.Lookup(0x10)
	entryB, _ := reg.Lookup(0x11)

	if entryA == nil {
		t.Fatalf("Lookup(0x10) returned nil; want entry A")
	}
	if entryB == nil {
		t.Fatalf("Lookup(0x11) returned nil; want entry B")
	}

	for _, a := range entryA.Addresses() {
		if a == 0x12 {
			t.Errorf("entry A (0x10, sn=SN-A-001) absorbed 0x12; want 0x12 NOT in entry A's address set")
		}
	}
	for _, a := range entryB.Addresses() {
		if a == 0x12 {
			t.Errorf("entry B (0x11, sn=SN-B-001) absorbed 0x12; want 0x12 NOT in entry B's address set")
		}
	}

	// A and B must remain distinct entries with their own serials.
	if entryA.SerialNumber() != "SN-A-001" {
		t.Errorf("entryA.SerialNumber() = %q; want \"SN-A-001\"", entryA.SerialNumber())
	}
	if entryB.SerialNumber() != "SN-B-001" {
		t.Errorf("entryB.SerialNumber() = %q; want \"SN-B-001\"", entryB.SerialNumber())
	}

	// The bare model-only observation at 0x12 resolves to its own entry.
	entry12, ok := reg.Lookup(0x12)
	if !ok {
		t.Fatalf("Lookup(0x12) = false; want a fresh entry from bare sig-only Register")
	}
	if entry12.SerialNumber() == "SN-A-001" || entry12.SerialNumber() == "SN-B-001" {
		t.Errorf("entry at 0x12 inherited a serial from A or B (= %q); want empty (fresh entry)", entry12.SerialNumber())
	}
}

// TestPartialModelObservationDoesNotMergeAcrossAddresses keeps the former
// signature-only regression boundary: a partial model observation is not a
// cross-address identity key after explicit alias enrichment.
func TestPartialModelObservationDoesNotMergeAcrossAddresses(t *testing.T) {
	reg := NewDeviceRegistry(nil)

	// Step 1: register canonical at 0x10 with MAC only.
	reg.Register(DeviceInfo{
		Address:      0x10,
		Manufacturer: "Vaillant",
		MacAddress:   "AA:BB:CC:11:22:33",
	})

	// Step 2: register secondary at 0x11 with sig only (no MAC).
	reg.Register(DeviceInfo{
		Address:         0x11,
		Manufacturer:    "Vaillant",
		DeviceID:        "BAI00",
		SoftwareVersion: "1201",
		HardwareVersion: "7603",
	})

	// Step 3: alias 0x11 into canonical at 0x10. AliasAddresses copies
	// the partial model fields into the topology group.
	if err := reg.AliasAddresses(0x10, 0x11); err != nil {
		t.Fatalf("AliasAddresses(0x10, 0x11) err=%v", err)
	}

	// Step 4: a bare model-only observation at a new address remains distinct.
	reg.Register(DeviceInfo{
		Address:         0x12,
		Manufacturer:    "Vaillant",
		DeviceID:        "BAI00",
		SoftwareVersion: "1201",
		HardwareVersion: "7603",
	})

	entry10, ok := reg.Lookup(0x10)
	if !ok {
		t.Fatalf("Lookup(0x10) = false; want merged canonical")
	}
	entry12, ok := reg.Lookup(0x12)
	if !ok {
		t.Fatalf("Lookup(0x12) = false; want a distinct partial observation")
	}
	if entry10 == entry12 {
		t.Error("bare model-only observation unexpectedly merged across addresses")
	}
	if entry10.MacAddress() != "AA:BB:CC:11:22:33" {
		t.Errorf("entry.MacAddress = %q; want AA:BB:CC:11:22:33", entry10.MacAddress())
	}
	if entry10.DeviceID() != "BAI00" {
		t.Errorf("entry.DeviceID = %q; want BAI00", entry10.DeviceID())
	}
}
