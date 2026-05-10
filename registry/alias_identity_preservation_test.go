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

// TestAliasAddresses_PreservesDistinctIdentityKeys asserts that
// when canonical and secondary entries have DIFFERENT identityKeys
// (e.g. canonical has a MAC-derived key, secondary has a
// serial-derived key), the merge re-points secondary's key at
// canonical in r.identity rather than deleting it. After the alias,
// both keys must resolve to the merged entry. (Codex P2 round-3
// finding 2026-05-08 on PR #136.)
func TestAliasAddresses_PreservesDistinctIdentityKeys(t *testing.T) {
	t.Parallel()

	reg := NewDeviceRegistry(nil)

	// Canonical at 0x10: identity by MAC only (no serial).
	reg.Register(DeviceInfo{
		Address:      0x10,
		Manufacturer: "Vaillant",
		DeviceID:     "BASV2",
		MacAddress:   "AA:BB:CC:DD:EE:01",
	})
	// Secondary at 0x15: identity by Serial only (no mac).
	reg.Register(DeviceInfo{
		Address:      0x15,
		Manufacturer: "Vaillant",
		DeviceID:     "BASV2",
		SerialNumber: "SN-DISTINCT-001",
	})

	if err := reg.AliasAddresses(0x10, 0x15); err != nil {
		t.Fatalf("AliasAddresses error = %v", err)
	}

	// Both lookup paths must resolve.
	byMac, ok := reg.lookupByIdentity(DeviceInfo{Manufacturer: "Vaillant", DeviceID: "BASV2", MacAddress: "AA:BB:CC:DD:EE:01"})
	if !ok {
		t.Fatalf("lookupByIdentity by MAC = false; want resolvable")
	}
	bySerial, ok := reg.lookupByIdentity(DeviceInfo{Manufacturer: "Vaillant", DeviceID: "BASV2", SerialNumber: "SN-DISTINCT-001"})
	if !ok {
		t.Fatalf("lookupByIdentity by Serial = false; want resolvable")
	}
	// Both must point at the same entry.
	if byMac.MacAddress() != "AA:BB:CC:DD:EE:01" {
		t.Errorf("by-MAC resolution: mac=%q; want AA:BB:CC:DD:EE:01", byMac.MacAddress())
	}
	if bySerial.SerialNumber() != "SN-DISTINCT-001" {
		t.Errorf("by-Serial resolution: serial=%q; want SN-DISTINCT-001", bySerial.SerialNumber())
	}
	// Same set of addresses (canonical + secondary's address).
	if !reflect.DeepEqual(byMac.Addresses(), bySerial.Addresses()) {
		t.Errorf("merge mismatch: byMac.Addresses=%v bySerial.Addresses=%v", byMac.Addresses(), bySerial.Addresses())
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

// TestRegister_FallbackSignatureNotPreservedAsAlias asserts that when
// an entry is first registered with only a model signature (no serial /
// MAC observed yet) and is later refreshed with a stable serial-derived
// key, the OLD `sig|...` fallback key is NOT preserved as an identity
// alias. Otherwise a second device with the same fingerprint that
// becomes ambiguous-by-signature would silently merge into the first
// entry on a subsequent bare sig-only observation, bypassing
// `lookupCompatibleBySignatureLocked`'s ambiguity-refusal scan.
//
// (Codex P2 round-7 finding 2026-05-08 on PR #136 thread
// PRRT_kwDORGIkfM6ArzFY: "Don't keep fallback signatures as identity
// aliases".)
func TestRegister_FallbackSignatureNotPreservedAsAlias(t *testing.T) {
	t.Parallel()

	reg := NewDeviceRegistry(nil)

	// Step 1: register entry A at 0x10 with signature-only identity
	// (no serial, no MAC). canonicalPhysicalIdentity.key() falls
	// back to "sig|VAILLANT|BAI00|0204|0102".
	reg.Register(DeviceInfo{
		Address:         0x10,
		Manufacturer:    "Vaillant",
		DeviceID:        "BAI00",
		SoftwareVersion: "0204",
		HardwareVersion: "0102",
	})

	// Step 2: refresh entry A with a stable serial — promotes
	// identityKey from "sig|..." to "sn|VAILLANT|SN-A-001".
	reg.Register(DeviceInfo{
		Address:         0x10,
		Manufacturer:    "Vaillant",
		DeviceID:        "BAI00",
		SerialNumber:    "SN-A-001",
		SoftwareVersion: "0204",
		HardwareVersion: "0102",
	})

	// Step 3: register entry B at 0x11 with the SAME signature but a
	// DIFFERENT serial. Both entries now share the same fallback
	// model signature, so `lookupCompatibleBySignatureLocked` would
	// correctly refuse the ambiguous match.
	reg.Register(DeviceInfo{
		Address:         0x11,
		Manufacturer:    "Vaillant",
		DeviceID:        "BAI00",
		SerialNumber:    "SN-B-001",
		SoftwareVersion: "0204",
		HardwareVersion: "0102",
	})

	// Step 4: a bare sig-only observation arrives at a NEW address
	// 0x12. canonicalPhysicalIdentity.key() returns
	// "sig|VAILLANT|BAI00|0204|0102". If the old sig key was
	// preserved as an alias to A, registerLocked's
	// `r.identity[identityKey]` lookup hits A directly, bypassing
	// the ambiguity scan, and 0x12 is incorrectly merged into A.
	reg.Register(DeviceInfo{
		Address:         0x12,
		Manufacturer:    "Vaillant",
		DeviceID:        "BAI00",
		SoftwareVersion: "0204",
		HardwareVersion: "0102",
	})

	// Expected: 0x12 is NOT merged into A. It must either be its
	// own new entry (because the ambiguity scan refused both A and
	// B) or remain unbound from any identity row. Critically, A's
	// addresses must NOT include 0x12, and B's must not either.
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
			t.Errorf("entry A (0x10, sn=SN-A-001) absorbed 0x12 via stale sig|... alias; want 0x12 NOT in entry A's address set")
		}
	}
	for _, a := range entryB.Addresses() {
		if a == 0x12 {
			t.Errorf("entry B (0x11, sn=SN-B-001) absorbed 0x12 via stale sig|... alias; want 0x12 NOT in entry B's address set")
		}
	}

	// A and B must remain distinct entries with their own serials.
	if entryA.SerialNumber() != "SN-A-001" {
		t.Errorf("entryA.SerialNumber() = %q; want \"SN-A-001\"", entryA.SerialNumber())
	}
	if entryB.SerialNumber() != "SN-B-001" {
		t.Errorf("entryB.SerialNumber() = %q; want \"SN-B-001\"", entryB.SerialNumber())
	}

	// The bare sig-only observation at 0x12 should resolve to its
	// own entry (registerLocked falls through to creating a fresh
	// entry when neither identity-by-key nor the ambiguity-checked
	// signature lookup matches).
	entry12, ok := reg.Lookup(0x12)
	if !ok {
		t.Fatalf("Lookup(0x12) = false; want a fresh entry from bare sig-only Register")
	}
	if entry12.SerialNumber() == "SN-A-001" || entry12.SerialNumber() == "SN-B-001" {
		t.Errorf("entry at 0x12 inherited a serial from A or B (= %q); want empty (fresh entry)", entry12.SerialNumber())
	}
}
