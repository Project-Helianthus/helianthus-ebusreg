package registry

import "testing"

// TestIdentityMerge_FourAddressesSameSerial asserts the operator-pinned
// scenario: NETX3 publishes itself on 4 bus addresses (0xF1 source +
// 0xF6 companion via canonical pair, plus 0xFF source + 0x04 wrap-pair
// companion). Once enrichment populates each address's SerialNumber +
// Manufacturer (e.g. via B509 ScanID probe), all four MUST resolve to
// a single DeviceEntry exposing all 4 addresses — operator should see
// ONE NETX3 device, not four "unknown" entries.
func TestIdentityMerge_FourAddressesSameSerial(t *testing.T) {
	t.Parallel()

	reg := NewDeviceRegistry(nil)

	// Each address is registered independently with the same enrichment
	// identity. This emulates the post-enrichment path where each slot
	// has been probed and the result includes manufacturer + serial.
	netx3 := DeviceInfo{
		Manufacturer: "Vaillant",
		DeviceID:     "NETX3",
		SerialNumber: "SN-NETX3-PHYS",
	}
	for _, addr := range []byte{0xF1, 0xF6, 0xFF, 0x04} {
		info := netx3
		info.Address = addr
		reg.Register(info)
	}

	// All four addresses resolve to the same DeviceEntry.
	primary := byte(0)
	for _, addr := range []byte{0xF1, 0xF6, 0xFF, 0x04} {
		entry, ok := reg.Lookup(addr)
		if !ok {
			t.Fatalf("Lookup(0x%02X) ok=false", addr)
		}
		if primary == 0 {
			primary = entry.Address()
		} else if entry.Address() != primary {
			t.Fatalf("0x%02X primary=0x%02X; want 0x%02X (all same SerialNumber should alias)", addr, entry.Address(), primary)
		}
	}

	// And exposes all 4 addresses via Addresses().
	entry, _ := reg.Lookup(0xF1)
	addrs := entry.Addresses()
	want := map[byte]bool{0xF1: false, 0xF6: false, 0xFF: false, 0x04: false}
	for _, a := range addrs {
		if _, expected := want[a]; expected {
			want[a] = true
		}
	}
	for a, found := range want {
		if !found {
			t.Errorf("Addresses() missing 0x%02X; got %v", a, addrs)
		}
	}
}

// TestIdentityMerge_DistinctSerialsStaySeparate asserts that addresses
// with DIFFERENT identity (e.g. BAI ecoTEC at 0x08 vs BASV2 at 0x15)
// remain separate DeviceEntries even when registered to the same
// registry instance.
func TestIdentityMerge_DistinctSerialsStaySeparate(t *testing.T) {
	t.Parallel()

	reg := NewDeviceRegistry(nil)
	reg.Register(DeviceInfo{Address: 0x08, Manufacturer: "Vaillant", DeviceID: "BAI00", SerialNumber: "SN-BAI"})
	reg.Register(DeviceInfo{Address: 0x15, Manufacturer: "Vaillant", DeviceID: "BASV2", SerialNumber: "SN-BASV2"})

	bai, _ := reg.Lookup(0x08)
	basv2, _ := reg.Lookup(0x15)
	if bai.Address() == basv2.Address() {
		t.Fatalf("distinct identities unexpectedly aliased: BAI primary=0x%02X, BASV2 primary=0x%02X", bai.Address(), basv2.Address())
	}
}

// TestIdentityMerge_LateEnrichmentAliasesPriorEntry asserts the most
// important M6 scenario: an address is FIRST observed passively (with
// no identity), THEN enriched with SerialNumber matching an existing
// device. The two MUST be aliased — proving identity-merge works for
// the inserter→enrichment timeline operators see in production.
func TestIdentityMerge_LateEnrichmentAliasesPriorEntry(t *testing.T) {
	t.Parallel()

	reg := NewDeviceRegistry(nil)

	// Step 1: 0xF6 known via active scan with full identity.
	reg.Register(DeviceInfo{Address: 0xF6, Manufacturer: "Vaillant", DeviceID: "NETX3", SerialNumber: "SN-NETX3"})
	pre, _ := reg.Lookup(0xF6)
	if pre.Address() != 0xF6 {
		t.Fatalf("setup: 0xF6 primary = 0x%02X; want 0xF6", pre.Address())
	}

	// Step 2: 0xF1 inserted passively via inserter (no identity yet).
	reg.Register(DeviceInfo{Address: 0xF1})
	mid, _ := reg.Lookup(0xF1)
	if mid.Address() == 0xF6 {
		t.Fatalf("0xF1 unexpectedly aliased to 0xF6 with no identity match: primary=0x%02X", mid.Address())
	}

	// Step 3: enrichment learns 0xF1's identity (e.g. via probe) — same
	// SerialNumber as 0xF6. Re-Register with full identity.
	reg.Register(DeviceInfo{Address: 0xF1, Manufacturer: "Vaillant", DeviceID: "NETX3", SerialNumber: "SN-NETX3"})

	post0xF1, _ := reg.Lookup(0xF1)
	post0xF6, _ := reg.Lookup(0xF6)
	if post0xF1.Address() != post0xF6.Address() {
		t.Fatalf("after late enrichment, 0xF1 primary=0x%02X 0xF6 primary=0x%02X; want same primary",
			post0xF1.Address(), post0xF6.Address())
	}

	addrs := post0xF1.Addresses()
	hasF1, hasF6 := false, false
	for _, a := range addrs {
		if a == 0xF1 {
			hasF1 = true
		}
		if a == 0xF6 {
			hasF6 = true
		}
	}
	if !hasF1 || !hasF6 {
		t.Errorf("post-enrichment Addresses() = %v; want both 0xF1 and 0xF6", addrs)
	}
}
