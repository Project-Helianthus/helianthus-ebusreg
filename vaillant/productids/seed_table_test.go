package productids_test

// RED tests for cruise plan address-table-registry-w19-26 M0B/productids — references M5A LoadSeedTable that intentionally doesn't exist yet.

import (
	"testing"

	"github.com/Project-Helianthus/helianthus-ebusreg/vaillant/productids"
)

const staticSeedSource = "vaillant_static_seed_w19_26"

func TestStaticSeedTable_NETX3_OwnsThreeAddresses(t *testing.T) {
	entry := findStaticSeedEntry(t, productids.LoadSeedTable(true), "Vaillant", "NETX3")

	assertSeedAddresses(t, entry.Addresses, []productids.SeedAddressEntry{
		{Addr: 0xF1, Role: "master", Confidence: "candidate"},
		{Addr: 0xF6, Role: "slave", Confidence: "candidate"},
		{Addr: 0x04, Role: "slave", Confidence: "candidate"},
	})
}

func TestStaticSeedTable_BASV2_OwnsTwoAddresses(t *testing.T) {
	entry := findStaticSeedEntry(t, productids.LoadSeedTable(true), "Vaillant", "BASV2")

	assertSeedAddresses(t, entry.Addresses, []productids.SeedAddressEntry{
		{Addr: 0x15, Role: "slave", Confidence: "candidate"},
		{Addr: 0xEC, Role: "slave", Confidence: "candidate"},
	})
}

func TestStaticSeedTable_FeatureFlagDefaultFalse(t *testing.T) {
	if got := productids.LoadSeedTable(false); len(got) != 0 {
		t.Fatalf("expected disabled static seed table to be empty, got %d entries", len(got))
	}
}

func TestStaticSeedEntry_Source(t *testing.T) {
	entries := productids.LoadSeedTable(true)
	if len(entries) == 0 {
		t.Fatal("expected enabled static seed table to return entries")
	}

	for _, entry := range entries {
		if entry.Source != staticSeedSource {
			t.Fatalf("expected Source %q, got %q", staticSeedSource, entry.Source)
		}
	}
}

func findStaticSeedEntry(t *testing.T, entries []productids.StaticSeedEntry, manufacturer, deviceID string) productids.StaticSeedEntry {
	t.Helper()

	for _, entry := range entries {
		if entry.Manufacturer == manufacturer && entry.DeviceID == deviceID {
			return entry
		}
	}

	t.Fatalf("expected static seed entry for Manufacturer=%q DeviceID=%q", manufacturer, deviceID)
	return productids.StaticSeedEntry{}
}

func assertSeedAddresses(t *testing.T, got, want []productids.SeedAddressEntry) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("expected %d seed addresses, got %d: %#v", len(want), len(got), got)
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("address[%d]: expected %#v, got %#v", i, want[i], got[i])
		}
	}
}
