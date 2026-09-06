package registry

import (
	"reflect"
	"sync"
	"testing"
	"time"
)

func sourceRetentionIdentity(address byte, serial string) DeviceInfo {
	return DeviceInfo{
		Address:      address,
		Manufacturer: "Vaillant",
		DeviceID:     "BASV2",
		SerialNumber: serial,
	}
}

func TestDiscoverySourceRetention_ActiveVerificationKeepsOriginalSource(t *testing.T) {
	t.Parallel()

	for _, scenario := range []struct {
		name       string
		seed       func(*DeviceRegistry, DeviceInfo)
		wantSource DiscoverySource
		wantState  VerificationState
	}{
		{
			name: "static seed", wantSource: DiscoverySourceStaticSeed,
			wantState: VerificationStateIdentityConfirmed,
			seed: func(registry *DeviceRegistry, info DeviceInfo) {
				registry.RegisterStaticSeed(info, SlotRoleSlave, time.Unix(1, 0))
			},
		},
		{
			name: "passive observation", wantSource: DiscoverySourcePassiveObserved,
			wantState: VerificationStateIdentityConfirmed,
			seed: func(registry *DeviceRegistry, info DeviceInfo) {
				registry.RegisterPassiveObserved(info, SlotRoleSlave, time.Unix(1, 0))
			},
		},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			registry := NewDeviceRegistry(nil)
			info := sourceRetentionIdentity(0x15, "SN-SOURCE-RETENTION-"+scenario.name)
			scenario.seed(registry, info)

			registry.Register(info)
			slot, ok := registry.LookupSlotSnapshot(info.Address)
			if !ok || slot.DiscoverySource != scenario.wantSource || slot.VerificationState != scenario.wantState {
				t.Fatalf("confirmed slot = %#v, present=%v; want source=%v state=%v", slot, ok, scenario.wantSource, scenario.wantState)
			}
		})
	}
}

func TestDiscoverySourceRetention_AliasTopologyAndSnapshotsStayPerFace(t *testing.T) {
	t.Parallel()

	registry := NewDeviceRegistry(nil)
	staticInfo := sourceRetentionIdentity(0x15, "SN-ALIAS-SOURCE-RETENTION")
	passiveInfo := staticInfo
	passiveInfo.Address = 0x16
	registry.RegisterStaticSeed(staticInfo, SlotRoleSlave, time.Unix(1, 0))
	registry.RegisterPassiveObserved(passiveInfo, SlotRoleSlave, time.Unix(2, 0))
	if err := registry.AliasAddresses(staticInfo.Address, passiveInfo.Address); err != nil {
		t.Fatal(err)
	}

	before, ok := registry.LookupSlotSnapshot(passiveInfo.Address)
	if !ok {
		t.Fatal("passive snapshot missing before confirmation")
	}
	registry.Register(staticInfo)
	afterStatic, staticOK := registry.LookupSlotSnapshot(staticInfo.Address)
	afterPassive, passiveOK := registry.LookupSlotSnapshot(passiveInfo.Address)
	if !staticOK || afterStatic.DiscoverySource != DiscoverySourceStaticSeed || afterStatic.VerificationState != VerificationStateIdentityConfirmed {
		t.Fatalf("static face after confirmation = %#v, present=%v", afterStatic, staticOK)
	}
	if !passiveOK || afterPassive.DiscoverySource != DiscoverySourcePassiveObserved || afterPassive.VerificationState != VerificationStateIdentityConfirmed {
		t.Fatalf("passive topology face after confirmation = %#v, present=%v", afterPassive, passiveOK)
	}
	if before.DiscoverySource != DiscoverySourcePassiveObserved || before.VerificationState != VerificationStateCorroborated {
		t.Fatalf("detached pre-confirmation snapshot changed = %#v", before)
	}
}

func TestDiscoverySourceRetention_RejectedConflictDoesNotMutateSlotOrGeneration(t *testing.T) {
	t.Parallel()

	registry := NewDeviceRegistry(nil)
	incumbent := sourceRetentionIdentity(0x15, "SN-CONFLICT-SOURCE-RETENTION")
	incumbent.MacAddress = "02:00:00:00:00:01"
	registry.RegisterPassiveObserved(incumbent, SlotRoleSlave, time.Unix(1, 0))
	before, ok := registry.LookupSlotSnapshot(incumbent.Address)
	if !ok {
		t.Fatal("incumbent slot missing")
	}
	beforeGeneration := readObservationGeneration(t, registry)

	conflicting := incumbent
	conflicting.MacAddress = "02:00:00:00:00:99"
	registry.Register(conflicting)

	after, ok := registry.LookupSlotSnapshot(incumbent.Address)
	if !ok || !reflect.DeepEqual(after, before) {
		t.Fatalf("conflicting active observation mutated slot: before=%#v after=%#v present=%v", before, after, ok)
	}
	if got := readObservationGeneration(t, registry); got != beforeGeneration {
		t.Fatalf("conflicting active observation generation = %d; want unchanged %d", got, beforeGeneration)
	}
}

func TestDiscoverySourceRetention_SnapshotRaceControl(t *testing.T) {
	registry := NewDeviceRegistry(nil)
	info := sourceRetentionIdentity(0x15, "SN-RACE-SOURCE-RETENTION")
	registry.RegisterStaticSeed(info, SlotRoleSlave, time.Unix(1, 0))

	var writers sync.WaitGroup
	writers.Add(1)
	go func() {
		defer writers.Done()
		for i := 0; i < 200; i++ {
			registry.Register(info)
		}
	}()
	for i := 0; i < 400; i++ {
		slot, ok := registry.LookupSlotSnapshot(info.Address)
		if !ok {
			t.Fatalf("snapshot missing at iteration %d", i)
		}
		if slot.DiscoverySource != DiscoverySourceStaticSeed {
			t.Fatalf("snapshot source = %v; want retained StaticSeed", slot.DiscoverySource)
		}
	}
	writers.Wait()
}
