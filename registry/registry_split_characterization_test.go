package registry

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestRegistrySplitCharacterization(t *testing.T) {
	reg := NewDeviceRegistry(nil)
	reg.Register(DeviceInfo{Address: 0x10, Manufacturer: "Vaillant", DeviceID: "BASV2", SerialNumber: "SN-10"})
	reg.Register(DeviceInfo{Address: 0x15, Manufacturer: "Vaillant", DeviceID: "BASV2", SerialNumber: "SN-15"})
	if err := reg.AliasAddresses(0x10, 0x15); err != nil {
		t.Fatalf("AliasAddresses: %v", err)
	}
	primary, ok := reg.Lookup(0x10)
	if !ok {
		t.Fatal("primary alias lookup failed")
	}
	alias, ok := reg.Lookup(0x15)
	if !ok || alias != primary {
		t.Fatalf("alias lookup = (%v, %t); want canonical entry", alias, ok)
	}
	if got, want := primary.Addresses(), []byte{0x10, 0x15}; !reflect.DeepEqual(got, want) {
		t.Fatalf("alias address order = %v; want %v", got, want)
	}

	observedAt := time.Date(2026, 8, 27, 12, 0, 0, 0, time.UTC)
	reg.RegisterPassiveObserved(DeviceInfo{Address: 0x26}, SlotRoleSlave, observedAt)
	before, ok := reg.LookupSlotSnapshot(0x26)
	if !ok {
		t.Fatal("passive slot snapshot missing")
	}
	reg.Register(DeviceInfo{Address: 0x26, Manufacturer: "Vaillant", DeviceID: "VR_71"})
	after, ok := reg.LookupSlotSnapshot(0x26)
	if !ok || before.DiscoverySource != DiscoverySourcePassiveObserved || after.DiscoverySource != DiscoverySourceActiveConfirmed {
		t.Fatalf("slot observation order = (%v, %v, %t); want passive then active", before.DiscoverySource, after.DiscoverySource, ok)
	}
	if !after.DeviceAttached {
		t.Fatal("active snapshot lost attached device")
	}

	first := make([]DeviceEntrySnapshot, 0)
	reg.IterateSnapshots(func(snapshot DeviceEntrySnapshot) bool {
		first = append(first, snapshot)
		return true
	})
	second := make([]DeviceEntrySnapshot, 0)
	reg.IterateSnapshots(func(snapshot DeviceEntrySnapshot) bool {
		second = append(second, snapshot)
		return true
	})
	if !reflect.DeepEqual(first, second) {
		t.Fatalf("snapshot iteration is non-deterministic: first=%#v second=%#v", first, second)
	}
	if len(first) == 0 || len(first[0].Addresses) == 0 {
		t.Fatal("snapshot fixture missing copied addresses")
	}
	original := first[0].Addresses[0]
	first[0].Addresses[0] = 0xff
	entry, ok := reg.Lookup(original)
	if !ok || entry.Addresses()[0] != original {
		t.Fatal("snapshot address mutation leaked into registry")
	}

	if _, err := CanonicalIndexForEntry(nil); !errors.Is(err, ErrProjectionInvalidNode) {
		t.Fatalf("CanonicalIndexForEntry(nil) error = %v; want %v", err, ErrProjectionInvalidNode)
	}
}
