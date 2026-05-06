package registry

import "testing"

/*
RED tests for cruise plan address-table-registry-w19-26 M0B — these reference
M1 types/funcs that intentionally don't exist yet.
*/

func TestAddressSlotLookupParity(t *testing.T) {
	registry := NewDeviceRegistry(nil)

	registered08 := registry.Register(DeviceInfo{
		Address:         0x08,
		Manufacturer:    "vaillant",
		DeviceID:        "BAI00",
		SerialNumber:    "NETX3-BOILER",
		SoftwareVersion: "1.0",
		HardwareVersion: "2.0",
	})
	registry.Register(DeviceInfo{
		Address:         0x15,
		Manufacturer:    "vaillant",
		DeviceID:        "BASV2",
		SerialNumber:    "NETX3-BASV2",
		SoftwareVersion: "1.1",
		HardwareVersion: "2.1",
	})

	var slot *AddressSlot
	var ok bool
	slot, ok = registry.Lookup(0x08)
	if !ok {
		t.Fatalf("Lookup(0x08) ok = false; want true")
	}
	if slot == nil {
		t.Fatalf("Lookup(0x08) slot = nil; want AddressSlot")
	}

	unset, ok := registry.Lookup(0xAA)
	if ok {
		t.Fatalf("Lookup(0xAA) ok = true; want false")
	}
	if unset != nil {
		t.Fatalf("Lookup(0xAA) slot = %v; want nil", unset)
	}

	var iterated08 DeviceEntry
	registry.Iterate(func(entry DeviceEntry) bool {
		if entry.Address() == 0x08 {
			iterated08 = entry
			return false
		}
		return true
	})
	if iterated08 == nil {
		t.Fatalf("Iterate did not return registered address 0x08")
	}
	if iterated08 != registered08 {
		t.Fatalf("legacy iteration entry = %p; want Register return %p", iterated08, registered08)
	}
	if slot.Device != registered08 {
		t.Fatalf("slot.Device = %p; want legacy DeviceEntry %p", slot.Device, registered08)
	}
}

func TestAddressSlotAliasing(t *testing.T) {
	registry := NewDeviceRegistry(nil)
	registry.Register(DeviceInfo{
		Address:         0xF1,
		Manufacturer:    "vaillant",
		DeviceID:        "NETX3",
		SerialNumber:    "NETX3",
		SoftwareVersion: "1.0",
		HardwareVersion: "1.0",
	})
	registry.Register(DeviceInfo{
		Address:         0xF6,
		Manufacturer:    "vaillant",
		DeviceID:        "NETX3",
		SerialNumber:    "NETX3",
		SoftwareVersion: "1.0",
		HardwareVersion: "1.0",
	})

	if err := registry.AliasAddresses(0xF1, 0xF6); err != nil {
		t.Fatalf("AliasAddresses(0xF1, 0xF6) error = %v", err)
	}

	var masterSlot *AddressSlot
	var ok bool
	masterSlot, ok = registry.Lookup(0xF1)
	if !ok {
		t.Fatalf("Lookup(0xF1) ok = false; want true")
	}
	var slaveSlot *AddressSlot
	slaveSlot, ok = registry.Lookup(0xF6)
	if !ok {
		t.Fatalf("Lookup(0xF6) ok = false; want true")
	}
	if masterSlot.Device == nil {
		t.Fatalf("Lookup(0xF1).Device = nil; want aliased device")
	}
	if masterSlot.Device != slaveSlot.Device {
		t.Fatalf("NETX3 slots are not aliased: 0xF1=%p 0xF6=%p", masterSlot.Device, slaveSlot.Device)
	}
}

func TestAddressSlotInternalFields(t *testing.T) {
	registry := NewDeviceRegistry(nil)
	registry.Register(DeviceInfo{
		Address:         0x15,
		Manufacturer:    "vaillant",
		DeviceID:        "BASV2",
		SerialNumber:    "BASV2-IDENTITY",
		SoftwareVersion: "3.0",
		HardwareVersion: "4.0",
	})

	wantSource := DiscoverySourceActiveConfirmed
	wantVerification := VerificationStateIdentityConfirmed

	var slot *AddressSlot
	var ok bool
	slot, ok = registry.Lookup(0x15)
	if !ok {
		t.Fatalf("Lookup(0x15) ok = false; want true")
	}
	if slot == nil {
		t.Fatalf("Lookup(0x15) slot = nil; want AddressSlot")
	}
	if slot.Addr != 0x15 {
		t.Fatalf("slot.Addr = 0x%02X; want 0x15", slot.Addr)
	}
	if slot.DiscoverySource != wantSource {
		t.Fatalf("slot.DiscoverySource = %v; want %v", slot.DiscoverySource, wantSource)
	}
	if slot.VerificationState != wantVerification {
		t.Fatalf("slot.VerificationState = %v; want %v", slot.VerificationState, wantVerification)
	}
	if slot.Device == nil {
		t.Fatalf("slot.Device = nil; want confirmed device")
	}
	if len(slot.Device.Faces) == 0 {
		t.Fatalf("slot.Device.Faces is empty; want at least one bus face")
	}
	for _, face := range slot.Device.Faces {
		if face.Addr == 0x15 {
			return
		}
	}
	t.Fatalf("slot.Device.Faces does not include address 0x15: %+v", slot.Device.Faces)
}
