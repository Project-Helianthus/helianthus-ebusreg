package registry

import (
	"reflect"
	"testing"
	"time"
)

func admissionQualifiedIdentity(address byte, serial string) DeviceInfo {
	return DeviceInfo{
		Address:      address,
		Manufacturer: "ACME",
		DeviceID:     "VR_71",
		SerialNumber: serial,
	}
}

type admissionRegistryState struct {
	observationGeneration uint64
	proofGeneration       uint64
	entries               map[byte]*deviceEntry
	order                 []*deviceEntry
	slots                 [256]AddressSlot
	slotPresent           [256]bool
}

func snapshotAdmissionRegistryState(r *DeviceRegistry) admissionRegistryState {
	r.mu.RLock()
	defer r.mu.RUnlock()

	state := admissionRegistryState{
		observationGeneration: r.observationGeneration,
		proofGeneration:       r.proofGeneration,
		entries:               make(map[byte]*deviceEntry, len(r.entries)),
		order:                 append([]*deviceEntry(nil), r.order...),
	}
	for address, entry := range r.entries {
		state.entries[address] = entry
	}
	for address, slot := range r.addressTable {
		if slot != nil {
			state.slotPresent[address] = true
			state.slots[address] = *slot
		}
	}
	return state
}

func requireAdmissionRegistryUnchanged(t *testing.T, before, after admissionRegistryState) {
	t.Helper()
	if !reflect.DeepEqual(before, after) {
		t.Fatalf("rejected admission mutated registry\nbefore: %#v\nafter:  %#v", before, after)
	}
}

func TestAdmitPassiveCompanionWithCurrentQualifiedIdentityWitness_AdmitsAndRefreshesPassiveSlaveSlot(t *testing.T) {
	now := time.Date(2026, 9, 6, 12, 0, 0, 0, time.UTC)

	t.Run("new companion slot", func(t *testing.T) {
		registry := NewDeviceRegistry(nil)
		registry.Register(admissionQualifiedIdentity(0x10, "SN-NEW"))
		before := snapshotAdmissionRegistryState(registry)

		companion, admitted := registry.AdmitPassiveCompanionWithCurrentQualifiedIdentityWitness(0x10, now)
		if !admitted || companion != 0x15 {
			t.Fatalf("admission = (0x%02X, %v); want (0x15, true)", companion, admitted)
		}
		slot, ok := registry.LookupSlotSnapshot(companion)
		if !ok || slot.Role != SlotRoleSlave || slot.DiscoverySource != DiscoverySourcePassiveObserved ||
			slot.VerificationState != VerificationStateCorroborated || !slot.FirstObservedAt.Equal(now) || !slot.LastObservedAt.Equal(now) || slot.DeviceAttached {
			t.Fatalf("companion slot = %#v, present=%v; want detached passive/corroborated target at %v", slot, ok, now)
		}
		after := snapshotAdmissionRegistryState(registry)
		if after.observationGeneration != before.observationGeneration+1 || after.proofGeneration != before.proofGeneration {
			t.Fatalf("generations after admission = observation:%d proof:%d; want observation:%d proof:%d", after.observationGeneration, after.proofGeneration, before.observationGeneration+1, before.proofGeneration)
		}
	})

	t.Run("refreshes an existing passive slot", func(t *testing.T) {
		registry := NewDeviceRegistry(nil)
		earlier := now.Add(-time.Minute)
		registry.MarkSlotPassiveObserved(0x15, SlotRoleSlave, earlier)
		registry.Register(admissionQualifiedIdentity(0x10, "SN-REFRESH"))

		companion, admitted := registry.AdmitPassiveCompanionWithCurrentQualifiedIdentityWitness(0x10, now)
		if !admitted || companion != 0x15 {
			t.Fatalf("refresh admission = (0x%02X, %v); want (0x15, true)", companion, admitted)
		}
		slot, ok := registry.LookupSlotSnapshot(0x15)
		if !ok || slot.Role != SlotRoleSlave || !slot.FirstObservedAt.Equal(earlier) || !slot.LastObservedAt.Equal(now) {
			t.Fatalf("refreshed slot = %#v, present=%v; want first=%v last=%v passive target", slot, ok, earlier, now)
		}
	})

	t.Run("synchronizes an attached entry face", func(t *testing.T) {
		registry := NewDeviceRegistry(nil)
		companionEntry := registry.Register(DeviceInfo{Address: 0x15})
		registry.Register(admissionQualifiedIdentity(0x10, "SN-ATTACHED"))

		companion, admitted := registry.AdmitPassiveCompanionWithCurrentQualifiedIdentityWitness(0x10, now)
		if !admitted || companion != 0x15 {
			t.Fatalf("attached admission = (0x%02X, %v); want (0x15, true)", companion, admitted)
		}
		if address, ok := companionEntry.AddressByRole(SlotRoleSlave); !ok || address != 0x15 {
			t.Fatalf("attached entry target face = (0x%02X, %v); want (0x15, true)", address, ok)
		}
		slot, ok := registry.LookupSlotSnapshot(0x15)
		if !ok || !slot.DeviceAttached || slot.Role != SlotRoleSlave {
			t.Fatalf("attached companion slot = %#v, present=%v; want attached target role", slot, ok)
		}
	})
}

func TestAdmitPassiveCompanionWithCurrentQualifiedIdentityWitness_RejectsWithoutMutation(t *testing.T) {
	now := time.Date(2026, 9, 6, 12, 0, 0, 0, time.UTC)

	t.Run("nil registry", func(t *testing.T) {
		var registry *DeviceRegistry
		if companion, admitted := registry.AdmitPassiveCompanionWithCurrentQualifiedIdentityWitness(0x10, now); admitted || companion != 0 {
			t.Fatalf("nil admission = (0x%02X, %v); want (0, false)", companion, admitted)
		}
	})

	for _, tc := range []struct {
		name   string
		source byte
		at     time.Time
		setup  func(*DeviceRegistry)
	}{
		{name: "zero observation time", source: 0x10, setup: func(r *DeviceRegistry) { r.Register(admissionQualifiedIdentity(0x10, "SN-ZERO")) }},
		{name: "noncanonical source", source: 0x26, at: now, setup: func(r *DeviceRegistry) { r.Register(admissionQualifiedIdentity(0x26, "SN-NONCANONICAL")) }},
		{name: "missing witness", source: 0x10, at: now, setup: func(*DeviceRegistry) {}},
		{name: "partial identity", source: 0x10, at: now, setup: func(r *DeviceRegistry) {
			r.Register(DeviceInfo{Address: 0x10, Manufacturer: "ACME", DeviceID: "VR_71"})
		}},
		{name: "last-known fields alone", source: 0x10, at: now, setup: func(r *DeviceRegistry) {
			r.Register(DeviceInfo{Address: 0x10, SoftwareVersion: "0204"})
		}},
		{name: "static evidence", source: 0x10, at: now, setup: func(r *DeviceRegistry) {
			r.RegisterStaticSeed(admissionQualifiedIdentity(0x10, "SN-STATIC"), SlotRoleMaster, now)
		}},
		{name: "passive evidence", source: 0x10, at: now, setup: func(r *DeviceRegistry) {
			r.RegisterPassiveObserved(admissionQualifiedIdentity(0x10, "SN-PASSIVE"), SlotRoleMaster, now)
		}},
		{name: "alias only", source: 0x10, at: now, setup: func(r *DeviceRegistry) {
			r.Register(admissionQualifiedIdentity(0x00, "SN-ALIAS"))
			if err := r.AliasAddresses(0x00, 0x10); err != nil {
				t.Fatalf("AliasAddresses: %v", err)
			}
		}},
		{name: "replaced non-direct authority", source: 0x10, at: now, setup: func(r *DeviceRegistry) {
			r.Register(admissionQualifiedIdentity(0x10, "SN-OLD"))
			r.RegisterStaticSeed(admissionQualifiedIdentity(0x10, "SN-REPLACED"), SlotRoleMaster, now)
		}},
		{name: "retired direct proof after sparse identity change", source: 0x10, at: now, setup: func(r *DeviceRegistry) {
			r.Register(admissionQualifiedIdentity(0x10, "SN-RETIRED"))
			r.Register(DeviceInfo{Address: 0x10, Manufacturer: "OTHER"})
		}},
		{name: "conflicted authority", source: 0x10, at: now, setup: func(r *DeviceRegistry) {
			incumbent := admissionQualifiedIdentity(0x11, "SN-CONFLICT")
			incumbent.MacAddress = "02:00:00:00:00:01"
			r.Register(incumbent)
			conflict := incumbent
			conflict.Address = 0x10
			conflict.MacAddress = "02:00:00:00:00:02"
			r.Register(conflict)
		}},
		{name: "companion already observed in source-side role", source: 0x10, at: now, setup: func(r *DeviceRegistry) {
			r.Register(admissionQualifiedIdentity(0x10, "SN-ROLE"))
			r.MarkSlotPassiveObserved(0x15, SlotRoleMaster, now)
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			registry := NewDeviceRegistry(nil)
			tc.setup(registry)
			before := snapshotAdmissionRegistryState(registry)
			companion, admitted := registry.AdmitPassiveCompanionWithCurrentQualifiedIdentityWitness(tc.source, tc.at)
			if admitted || companion != 0 {
				t.Fatalf("rejected admission = (0x%02X, %v); want (0, false)", companion, admitted)
			}
			requireAdmissionRegistryUnchanged(t, before, snapshotAdmissionRegistryState(registry))
		})
	}
}

func TestAdmitPassiveCompanionWithCurrentQualifiedIdentityWitness_RequiresExactCurrentWitness(t *testing.T) {
	now := time.Date(2026, 9, 6, 12, 0, 0, 0, time.UTC)
	registry := NewDeviceRegistry(nil)
	registry.Register(admissionQualifiedIdentity(0x00, "SN-OTHER-ADDRESS"))
	registry.mu.Lock()
	registry.qualifiedWitnesses[0x10] = registry.qualifiedWitnesses[0x00]
	registry.mu.Unlock()

	before := snapshotAdmissionRegistryState(registry)
	if companion, admitted := registry.AdmitPassiveCompanionWithCurrentQualifiedIdentityWitness(0x10, now); admitted || companion != 0 {
		t.Fatalf("wrong-address witness admission = (0x%02X, %v); want (0, false)", companion, admitted)
	}
	requireAdmissionRegistryUnchanged(t, before, snapshotAdmissionRegistryState(registry))
}

func TestAdmitPassiveCompanionWithCurrentQualifiedIdentityWitness_LinearizesValidationAndPassiveSlot(t *testing.T) {
	now := time.Date(2026, 9, 6, 12, 0, 0, 0, time.UTC)
	registry := NewDeviceRegistry(nil)
	registry.Register(admissionQualifiedIdentity(0x10, "SN-INITIAL"))

	registry.mu.Lock()
	writerStarted := make(chan struct{})
	writerDone := make(chan struct{})
	go func() {
		close(writerStarted)
		registry.RegisterStaticSeed(admissionQualifiedIdentity(0x10, "SN-REPLACEMENT"), SlotRoleMaster, now)
		close(writerDone)
	}()
	<-writerStarted

	if !registry.admitPassiveCompanionWithCurrentQualifiedIdentityWitnessLocked(0x10, 0x15, now) {
		registry.mu.Unlock()
		t.Fatal("locked admission rejected a current exact-source witness")
	}
	select {
	case <-writerDone:
		registry.mu.Unlock()
		t.Fatal("writer interleaved before the admitted passive slot was committed")
	default:
	}
	slot := registry.addressTable[0x15]
	if slot == nil || slot.Role != SlotRoleSlave || slot.DiscoverySource != DiscoverySourcePassiveObserved || slot.VerificationState != VerificationStateCorroborated || !slot.LastObservedAt.Equal(now) {
		registry.mu.Unlock()
		t.Fatalf("slot at linearization point = %#v; want committed passive target", slot)
	}
	registry.mu.Unlock()
	<-writerDone

	if _, ok := currentQualifiedIdentityWitness(registry, 0x10); ok {
		t.Fatal("replacement writer did not retire the old direct witness")
	}
}
