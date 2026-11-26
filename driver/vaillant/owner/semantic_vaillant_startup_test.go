package owner

import (
	"testing"
)

func TestVaillantOwnerStartupRadioDeviceIncludeSkipsDisconnectedRegulatorSlots(t *testing.T) {
	t.Parallel()

	classAddress := uint8(0x26)

	if include, _ := startupRadioDeviceInclude(remoteRegulators.group, false, &classAddress); include {
		t.Fatal("disconnected regulator slot included; want skipped")
	}
	if include, mode := startupRadioDeviceInclude(remoteFunctionalModules.group, false, &classAddress); !include || mode != "inventory" {
		t.Fatalf("functional module identity evidence include=%v mode=%q; want inventory include", include, mode)
	}
	unknownClassAddress := uint8(0x00)
	if include, _ := startupRadioDeviceInclude(remoteFunctionalModules.group, false, &unknownClassAddress); include {
		t.Fatal("unknown disconnected functional module slot included; want skipped")
	}
	if include, _ := startupRadioDeviceInclude(remoteFunctionalModules.group, false, nil); include {
		t.Fatal("empty functional module slot included; want skipped")
	}
}

func TestVaillantOwnerStartupRadioFullScanGroupsAreSelectedPerGroup(t *testing.T) {
	t.Parallel()

	allUnseeded := startupRadioFullScanGroups(nil)
	for _, grp := range remoteDeviceGroups {
		if !allUnseeded[grp.group] {
			t.Fatalf("empty radio discovery should require full scan for group 0x%02x", grp.group)
		}
	}

	discovered := map[radioDeviceKey]*vaillantRadioDeviceSnapshot{
		{Group: remoteRegulators.group, Instance: 0}: nil,
	}
	seededRegulator := startupRadioFullScanGroups(discovered)
	if seededRegulator[remoteRegulators.group] {
		t.Fatal("low-slot seeded regulator group should stay on fast scan")
	}
	if !seededRegulator[remoteThermostats.group] {
		t.Fatal("unseeded thermostat group should still require full scan")
	}

	discovered[radioDeviceKey{Group: remoteThermostats.group, Instance: semanticStartupSlotFastMaxInstance + 1}] = nil
	highSeededThermostat := startupRadioFullScanGroups(discovered)
	if !highSeededThermostat[remoteThermostats.group] {
		t.Fatal("seeded high thermostat slot should require full scan")
	}
}
