package registry

import (
	"fmt"
	"slices"
	"testing"
	"time"
)

func TestIdentityEnrichmentPreservesBothVR940fAliasPairs(t *testing.T) {
	for _, first := range []byte{0x04, 0xf6} {
		t.Run(fmt.Sprintf("first_%02x", first), func(t *testing.T) {
			reg := NewDeviceRegistry(nil)
			observedAt := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
			for _, address := range []byte{0xff, 0x04, 0xf1, 0xf6} {
				reg.RegisterPassiveObserved(DeviceInfo{Address: address}, SlotRoleUnknown, observedAt)
			}
			for _, pair := range [][2]byte{{0xff, 0x04}, {0xf1, 0xf6}} {
				if err := reg.AliasAddresses(pair[0], pair[1]); err != nil {
					t.Fatal(err)
				}
			}

			second := byte(0xf6)
			if first == second {
				second = 0x04
			}
			for _, address := range []byte{first, second} {
				reg.Register(DeviceInfo{Address: address, Manufacturer: "Vaillant", DeviceID: "NETX3", SerialNumber: "SYNTHETIC-VR940F", SoftwareVersion: "0129", HardwareVersion: "0404"})
			}

			canonical, _ := reg.Lookup(first)
			for _, address := range []byte{0x04, 0xff, 0xf1, 0xf6} {
				entry, ok := reg.Lookup(address)
				if !ok || entry != canonical {
					t.Fatalf("address %02x lost its evidenced alias pair during serial enrichment", address)
				}
				slot, ok := reg.LookupSlot(address)
				if !ok || slot.Device != entry {
					t.Fatalf("address %02x slot does not reference the canonical entry", address)
				}
				if slot.FirstObservedAt != observedAt {
					t.Fatalf("address %02x lost its original observation timestamp", address)
				}
			}
			addresses := canonical.Addresses()
			slices.Sort(addresses)
			if !slices.Equal(addresses, []byte{0x04, 0xf1, 0xf6, 0xff}) {
				t.Fatalf("addresses = %x", addresses)
			}
			count := 0
			reg.Iterate(func(DeviceEntry) bool { count++; return true })
			if count != 1 {
				t.Fatalf("entries = %d; want one physical VR940f", count)
			}
		})
	}
}

func TestIdentityEnrichmentLeavesConflictingCompanionWithItsOriginalDevice(t *testing.T) {
	reg := NewDeviceRegistry(nil)
	reg.Register(DeviceInfo{Address: 0x04, Manufacturer: "Vaillant", DeviceID: "NETX3", SerialNumber: "SYNTHETIC-A"})
	reg.Register(DeviceInfo{Address: 0xf1, Manufacturer: "Vaillant", DeviceID: "NETX3", SerialNumber: "SYNTHETIC-B"})
	if err := reg.AliasAddresses(0xf1, 0xf6); err != nil {
		t.Fatal(err)
	}

	// A new identity at one address is not evidence that the previous
	// physical device's other address moved with it.
	reg.Register(DeviceInfo{Address: 0xf6, Manufacturer: "Vaillant", DeviceID: "NETX3", SerialNumber: "SYNTHETIC-A"})
	original, _ := reg.Lookup(0xf1)
	moved, _ := reg.Lookup(0xf6)
	if original == moved || original.SerialNumber() != "SYNTHETIC-B" || moved.SerialNumber() != "SYNTHETIC-A" {
		t.Fatal("conflicting stable identities were merged")
	}
	if !slices.Equal(original.Addresses(), []byte{0xf1}) {
		t.Fatalf("original addresses = %x", original.Addresses())
	}
}
