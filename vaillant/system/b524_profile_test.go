package system

import "testing"

func TestB524DiscoveryProfiles(t *testing.T) {
	t.Parallel()

	profiles := map[byte]B524GroupProfile{}
	for _, profile := range B524DiscoveryProfiles {
		profiles[profile.Group] = profile
	}

	assertProfile(t, profiles, 0x02, B524OpcodeLocal, 0x0A, 0x0021)
	assertProfile(t, profiles, 0x03, B524OpcodeLocal, 0x0A, 0x002F)
	assertProfile(t, profiles, 0x09, B524OpcodeRemote, 0x0A, 0x002F)
	assertProfile(t, profiles, 0x0A, B524OpcodeRemote, 0x0A, 0x003F)
	assertProfile(t, profiles, 0x0C, B524OpcodeRemote, 0x0A, 0x003F)
}

func assertProfile(t *testing.T, profiles map[byte]B524GroupProfile, group, opcode, instanceMax byte, registerMax uint16) {
	t.Helper()

	profile, ok := profiles[group]
	if !ok {
		t.Fatalf("missing profile for group 0x%02X", group)
	}
	if profile.Opcode != opcode {
		t.Fatalf("group 0x%02X opcode = 0x%02X, want 0x%02X", group, profile.Opcode, opcode)
	}
	if profile.InstanceMax != instanceMax {
		t.Fatalf("group 0x%02X instance max = 0x%02X, want 0x%02X", group, profile.InstanceMax, instanceMax)
	}
	if profile.RegisterMax != registerMax {
		t.Fatalf("group 0x%02X register max = 0x%04X, want 0x%04X", group, profile.RegisterMax, registerMax)
	}
}
