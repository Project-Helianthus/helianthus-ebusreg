package system

import "testing"

type profileKey struct {
	Group  byte
	Opcode byte
}

func TestB524DiscoveryProfiles(t *testing.T) {
	t.Parallel()

	profiles := map[profileKey]B524GroupProfile{}
	for _, profile := range B524DiscoveryProfiles {
		key := profileKey{Group: profile.Group, Opcode: profile.Opcode}
		if _, exists := profiles[key]; exists {
			t.Fatalf("duplicate profile for group 0x%02X opcode 0x%02X", profile.Group, profile.Opcode)
		}
		profiles[key] = profile
	}

	// Singleton groups
	assertProfile(t, profiles, 0x00, B524OpcodeLocal, 0x00, 0x00FF)
	assertProfile(t, profiles, 0x01, B524OpcodeLocal, 0x00, 0x0013)

	// Multi-instance local groups
	assertProfile(t, profiles, 0x02, B524OpcodeLocal, 0x0A, 0x0025)
	assertProfile(t, profiles, 0x03, B524OpcodeLocal, 0x0A, 0x002F)

	// Dual-namespace groups
	assertProfile(t, profiles, 0x08, B524OpcodeLocal, 0x00, 0x0007)
	assertProfile(t, profiles, 0x08, B524OpcodeRemote, 0x0A, 0x0004)
	assertProfile(t, profiles, 0x09, B524OpcodeLocal, 0x0A, 0x000F)
	assertProfile(t, profiles, 0x09, B524OpcodeRemote, 0x0A, 0x002F)
	assertProfile(t, profiles, 0x0A, B524OpcodeLocal, 0x0A, 0x004D)
	assertProfile(t, profiles, 0x0A, B524OpcodeRemote, 0x0A, 0x003F)

	// Remote-only groups
	assertProfile(t, profiles, 0x0C, B524OpcodeRemote, 0x0A, 0x003F)
}

func TestB524Profile_RegistersDualNamespaceEntries(t *testing.T) {
	t.Parallel()

	type dualEntry struct {
		group             byte
		localInstanceMax  byte
		localRegisterMax  uint16
		remoteInstanceMax byte
		remoteRegisterMax uint16
	}

	dualGroups := []dualEntry{
		{group: 0x08, localInstanceMax: 0x00, localRegisterMax: 0x0007, remoteInstanceMax: 0x0A, remoteRegisterMax: 0x0004},
		{group: 0x09, localInstanceMax: 0x0A, localRegisterMax: 0x000F, remoteInstanceMax: 0x0A, remoteRegisterMax: 0x002F},
		{group: 0x0A, localInstanceMax: 0x0A, localRegisterMax: 0x004D, remoteInstanceMax: 0x0A, remoteRegisterMax: 0x003F},
	}

	profiles := map[profileKey]B524GroupProfile{}
	for _, p := range B524DiscoveryProfiles {
		profiles[profileKey{Group: p.Group, Opcode: p.Opcode}] = p
	}

	for _, de := range dualGroups {
		localKey := profileKey{Group: de.group, Opcode: B524OpcodeLocal}
		local, ok := profiles[localKey]
		if !ok {
			t.Fatalf("missing local profile for dual-namespace group 0x%02X", de.group)
		}
		if local.InstanceMax != de.localInstanceMax {
			t.Fatalf("group 0x%02X local instance max = 0x%02X, want 0x%02X", de.group, local.InstanceMax, de.localInstanceMax)
		}
		if local.RegisterMax != de.localRegisterMax {
			t.Fatalf("group 0x%02X local register max = 0x%04X, want 0x%04X", de.group, local.RegisterMax, de.localRegisterMax)
		}

		remoteKey := profileKey{Group: de.group, Opcode: B524OpcodeRemote}
		remote, ok := profiles[remoteKey]
		if !ok {
			t.Fatalf("missing remote profile for dual-namespace group 0x%02X", de.group)
		}
		if remote.InstanceMax != de.remoteInstanceMax {
			t.Fatalf("group 0x%02X remote instance max = 0x%02X, want 0x%02X", de.group, remote.InstanceMax, de.remoteInstanceMax)
		}
		if remote.RegisterMax != de.remoteRegisterMax {
			t.Fatalf("group 0x%02X remote register max = 0x%04X, want 0x%04X", de.group, remote.RegisterMax, de.remoteRegisterMax)
		}
	}
}

func assertProfile(t *testing.T, profiles map[profileKey]B524GroupProfile, group, opcode, instanceMax byte, registerMax uint16) {
	t.Helper()

	key := profileKey{Group: group, Opcode: opcode}
	profile, ok := profiles[key]
	if !ok {
		t.Fatalf("missing profile for group 0x%02X opcode 0x%02X", group, opcode)
	}
	if profile.InstanceMax != instanceMax {
		t.Fatalf("group 0x%02X opcode 0x%02X instance max = 0x%02X, want 0x%02X", group, opcode, profile.InstanceMax, instanceMax)
	}
	if profile.RegisterMax != registerMax {
		t.Fatalf("group 0x%02X opcode 0x%02X register max = 0x%04X, want 0x%04X", group, opcode, profile.RegisterMax, registerMax)
	}
}
