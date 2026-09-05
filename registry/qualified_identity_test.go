package registry

import (
	"fmt"
	"testing"
	"time"
)

func qualifiedIdentity(address byte, serial string) DeviceInfo {
	return DeviceInfo{
		Address:      address,
		Manufacturer: "Vaillant",
		DeviceID:     "BASV2",
		SerialNumber: serial,
	}
}

func requireSeparateEntries(t *testing.T, registry *DeviceRegistry, first, second byte) {
	t.Helper()
	firstEntry, firstOK := registry.Lookup(first)
	secondEntry, secondOK := registry.Lookup(second)
	if !firstOK || !secondOK {
		t.Fatalf("lookups present = (%v, %v); want both", firstOK, secondOK)
	}
	if firstEntry == secondEntry {
		t.Fatalf("0x%02X and 0x%02X unexpectedly share an entry", first, second)
	}
}

func TestQualifiedIdentity_RequiresExactCompleteTriple(t *testing.T) {
	t.Parallel()

	t.Run("different device ID remains distinct", func(t *testing.T) {
		registry := NewDeviceRegistry(nil)
		registry.Register(qualifiedIdentity(0x15, "SN-159"))
		other := qualifiedIdentity(0xEC, "SN-159")
		other.DeviceID = "SOL00"
		registry.Register(other)
		requireSeparateEntries(t, registry, 0x15, 0xEC)
	})

	t.Run("normalized exact triple merges 0x15 and 0xEC", func(t *testing.T) {
		registry := NewDeviceRegistry(nil)
		registry.Register(qualifiedIdentity(0x15, " sn-159 "))
		other := qualifiedIdentity(0xEC, "SN-159")
		other.Manufacturer = " vaillant "
		other.DeviceID = "basv2"
		registry.Register(other)
		first, _ := registry.Lookup(0x15)
		second, _ := registry.Lookup(0xEC)
		if first != second {
			t.Fatal("complete normalized triples must merge")
		}
	})

	for _, partial := range []struct {
		name   string
		mutate func(*DeviceInfo)
	}{
		{name: "missing manufacturer", mutate: func(info *DeviceInfo) { info.Manufacturer = "" }},
		{name: "missing device ID", mutate: func(info *DeviceInfo) { info.DeviceID = "" }},
		{name: "missing serial", mutate: func(info *DeviceInfo) { info.SerialNumber = "" }},
	} {
		t.Run(partial.name, func(t *testing.T) {
			registry := NewDeviceRegistry(nil)
			registry.Register(qualifiedIdentity(0x15, "SN-159"))
			other := qualifiedIdentity(0xEC, "SN-159")
			partial.mutate(&other)
			registry.Register(other)
			requireSeparateEntries(t, registry, 0x15, 0xEC)
		})
	}
}

func TestQualifiedIdentity_PreservesPunctuationAtMemberBoundaries(t *testing.T) {
	t.Parallel()

	t.Run("ambiguous delimiter layouts remain distinct", func(t *testing.T) {
		registry := NewDeviceRegistry(nil)
		first := registry.Register(DeviceInfo{Address: 0x10, Manufacturer: "A|B", DeviceID: "C", SerialNumber: "D"})
		second := registry.Register(DeviceInfo{Address: 0x11, Manufacturer: "A", DeviceID: "B", SerialNumber: "C|D"})
		if first == second {
			t.Fatal("distinct triples with preserved punctuation merged")
		}
	})

	t.Run("punctuation in every member still supports normalized equality", func(t *testing.T) {
		registry := NewDeviceRegistry(nil)
		first := registry.Register(DeviceInfo{Address: 0x12, Manufacturer: " a|b, ", DeviceID: " c:d|e ", SerialNumber: " serial|1,2:3 "})
		second := registry.Register(DeviceInfo{Address: 0x13, Manufacturer: "A|B,", DeviceID: "C:D|E", SerialNumber: "SERIAL|1,2:3"})
		if first != second {
			t.Fatal("normalized complete triple with punctuation did not merge")
		}
	})
}

func TestQualifiedIdentity_RejectsSentinelSerialTextualVariants(t *testing.T) {
	t.Parallel()

	// Normalization trims surrounding space, upper-cases, accepts one optional
	// 0x prefix, and ignores leading zeroes only for all-hex sentinel spellings.
	// It deliberately leaves ordinary serial formats (including punctuation) alone.
	for _, serial := range []string{
		"0", "00000000", " 0x00000000 ", "0X000000000",
		"FFFFFFFF", "0xffffffff", "00000000fFfFfFfF",
		"7FFFFFFF", "0x7fffffff", "000000007FFFFFFF",
	} {
		t.Run(fmt.Sprintf("%q", serial), func(t *testing.T) {
			registry := NewDeviceRegistry(nil)
			registry.Register(qualifiedIdentity(0x15, serial))
			registry.Register(qualifiedIdentity(0xEC, serial))
			requireSeparateEntries(t, registry, 0x15, 0xEC)
		})
	}
}

func TestQualifiedIdentity_PreservesSameAddressEnrichment(t *testing.T) {
	t.Parallel()

	registry := NewDeviceRegistry(nil)
	registry.Register(DeviceInfo{Address: 0x15, Manufacturer: "Vaillant", DeviceID: "BASV2"})
	registry.Register(qualifiedIdentity(0x15, "SN-159"))
	entry, ok := registry.Lookup(0x15)
	if !ok || entry.SerialNumber() != "SN-159" || entry.DeviceID() != "BASV2" {
		t.Fatalf("same-address enrichment = %#v, present=%v; want qualified BASV2 identity", entry, ok)
	}
}

func TestQualifiedIdentity_RetiresCorrectedTripleFromIndependentSelection(t *testing.T) {
	t.Parallel()

	for _, correction := range []struct {
		name   string
		mutate func(*DeviceInfo)
	}{
		{name: "manufacturer", mutate: func(info *DeviceInfo) { info.Manufacturer = "Saunier Duval" }},
		{name: "device ID", mutate: func(info *DeviceInfo) { info.DeviceID = "DEV30" }},
		{name: "serial", mutate: func(info *DeviceInfo) { info.SerialNumber = "SN-CORRECTED" }},
	} {
		t.Run(correction.name, func(t *testing.T) {
			registry := NewDeviceRegistry(nil)
			old := DeviceInfo{Address: 0x10, Manufacturer: "Vaillant", DeviceID: "OLD30", SerialNumber: "SN-OLD"}
			corrected := old
			correction.mutate(&corrected)

			registry.Register(old)
			registry.Register(corrected)

			// A sparse same-address observation retains the corrected LKG
			// triple; it must not revive the retired one.
			registry.Register(DeviceInfo{Address: corrected.Address, SoftwareVersion: "0204"})
			current, ok := registry.Lookup(corrected.Address)
			if !ok || canonicalPhysicalIdentity(DeviceInfo{
				Manufacturer: current.Manufacturer(), DeviceID: current.DeviceID(), SerialNumber: current.SerialNumber(),
			}).key() != canonicalPhysicalIdentity(corrected).key() {
				t.Fatalf("sparse correction retention = (%q, %q, %q), present=%v; want current corrected triple", current.Manufacturer(), current.DeviceID(), current.SerialNumber(), ok)
			}

			if _, ok := registry.lookupByIdentity(old); ok {
				t.Fatal("retired triple remained selectable before any independent observation")
			}
			if byCurrent, ok := registry.lookupByIdentity(corrected); !ok || byCurrent != current {
				t.Fatal("current corrected triple was not selectable")
			}

			independentOld := old
			independentOld.Address = 0x11
			registry.Register(independentOld)
			requireSeparateEntries(t, registry, corrected.Address, independentOld.Address)

			independentCurrent := corrected
			independentCurrent.Address = 0x12
			registry.Register(independentCurrent)
			merged, _ := registry.Lookup(independentCurrent.Address)
			if merged != current {
				t.Fatal("independent current triple did not merge with corrected entry")
			}
		})
	}
}

func TestQualifiedIdentity_StaticSeedPromotionRequiresQualifiedObservation(t *testing.T) {
	t.Parallel()

	registry := NewDeviceRegistry(nil)
	registry.RegisterStaticSeed(DeviceInfo{Address: 0x15, Manufacturer: "Vaillant", DeviceID: "BASV2"}, SlotRoleSlave, time.Unix(1, 0))
	registry.Register(DeviceInfo{Address: 0x15, Manufacturer: "Vaillant", DeviceID: "BASV2", SerialNumber: "0x00000000"})
	partial, _ := registry.LookupSlot(0x15)
	if partial.VerificationState != VerificationStateCandidate {
		t.Fatalf("sentinel observation state = %v; want Candidate", partial.VerificationState)
	}

	registry.Register(qualifiedIdentity(0x15, "SN-159"))
	confirmed, _ := registry.LookupSlot(0x15)
	if confirmed.VerificationState != VerificationStateIdentityConfirmed {
		t.Fatalf("qualified observation state = %v; want IdentityConfirmed", confirmed.VerificationState)
	}
}

func TestQualifiedIdentity_CurrentSessionConfirmationIsPerFace(t *testing.T) {
	t.Parallel()

	for _, scenario := range []struct {
		name       string
		register   func(*DeviceRegistry, DeviceInfo, time.Time)
		address    byte
		wantSource DiscoverySource
		wantState  VerificationState
	}{
		{
			name: "static face in exact triple group remains candidate",
			register: func(registry *DeviceRegistry, info DeviceInfo, observedAt time.Time) {
				registry.RegisterStaticSeed(info, SlotRoleSlave, observedAt)
			},
			address: 0x15, wantSource: DiscoverySourceStaticSeed, wantState: VerificationStateCandidate,
		},
		{
			name: "passive face in exact triple group remains corroborated",
			register: func(registry *DeviceRegistry, info DeviceInfo, observedAt time.Time) {
				registry.RegisterPassiveObserved(info, SlotRoleMaster, observedAt)
			},
			address: 0x31, wantSource: DiscoverySourcePassiveObserved, wantState: VerificationStateCorroborated,
		},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			registry := NewDeviceRegistry(nil)
			observedAt := time.Unix(1, 0)
			scenario.register(registry, qualifiedIdentity(scenario.address, "SN-PER-FACE"), observedAt)

			// This is the only current-session active observation. An exact
			// triple merge joins the entries but does not evidence the older face.
			registry.Register(qualifiedIdentity(0xEC, "SN-PER-FACE"))

			older, ok := registry.LookupSlotSnapshot(scenario.address)
			if !ok || older.DiscoverySource != scenario.wantSource || older.VerificationState != scenario.wantState || !older.FirstObservedAt.Equal(observedAt) {
				t.Fatalf("older face = %#v, present=%v; want source=%v state=%v original=%v", older, ok, scenario.wantSource, scenario.wantState, observedAt)
			}
			direct, ok := registry.LookupSlotSnapshot(0xEC)
			if !ok || direct.DiscoverySource != DiscoverySourceActiveConfirmed || direct.VerificationState != VerificationStateIdentityConfirmed {
				t.Fatalf("directly observed face = %#v, present=%v; want active-confirmed/identity-confirmed", direct, ok)
			}
		})
	}
}

func TestQualifiedIdentity_ConfirmedTopologyMayPromoteCompanionFace(t *testing.T) {
	t.Parallel()

	registry := NewDeviceRegistry(nil)
	seededAt := time.Unix(1, 0)
	registry.RegisterStaticSeed(DeviceInfo{Address: 0x15, Manufacturer: "Vaillant", DeviceID: "BASV2"}, SlotRoleSlave, seededAt)
	registry.Register(DeviceInfo{Address: 0xEC, Manufacturer: "Vaillant", DeviceID: "BASV2"})
	if err := registry.AliasAddresses(0xEC, 0x15); err != nil {
		t.Fatal(err)
	}

	// AliasAddresses is the explicit source-target/canonical-companion
	// topology relation. A later qualified active observation at 0xEC may
	// therefore confirm that companion without using entry membership alone.
	registry.Register(qualifiedIdentity(0xEC, "SN-TOPOLOGY"))

	companion, ok := registry.LookupSlotSnapshot(0x15)
	if !ok || companion.DiscoverySource != DiscoverySourceStaticSeed || companion.VerificationState != VerificationStateIdentityConfirmed || !companion.FirstObservedAt.Equal(seededAt) {
		t.Fatalf("topology companion = %#v, present=%v; want preserved static provenance and identity confirmation", companion, ok)
	}
}

func TestQualifiedIdentity_SplitRejoinInvalidatesPriorTopology(t *testing.T) {
	const source, companion = byte(0x20), byte(0x21)
	seededAt := time.Unix(1, 0)
	registry := NewDeviceRegistry(nil)

	registry.RegisterStaticSeed(qualifiedIdentity(source, "SN-TOPOLOGY-A"), SlotRoleSlave, seededAt)
	registry.RegisterStaticSeed(qualifiedIdentity(companion, "SN-TOPOLOGY-A"), SlotRoleSlave, seededAt)
	if err := registry.AliasAddresses(source, companion); err != nil {
		t.Fatal(err)
	}

	// The conflicting same-address identity detaches companion and must retire
	// its old topology evidence. Rejoining by the exact triple alone does not
	// recreate that evidence.
	registry.RegisterStaticSeed(qualifiedIdentity(companion, "SN-TOPOLOGY-B"), SlotRoleSlave, seededAt)
	registry.RegisterStaticSeed(qualifiedIdentity(companion, "SN-TOPOLOGY-A"), SlotRoleSlave, seededAt)
	registry.Register(qualifiedIdentity(source, "SN-TOPOLOGY-A"))

	slot, ok := registry.LookupSlotSnapshot(companion)
	if !ok || slot.DiscoverySource != DiscoverySourceStaticSeed || slot.VerificationState != VerificationStateCandidate || !slot.FirstObservedAt.Equal(seededAt) {
		t.Fatalf("rejoined companion = %#v, present=%v; want retained static candidate", slot, ok)
	}
}

func TestQualifiedIdentity_FreshTopologyAfterSplitRejoinMayPromoteCompanion(t *testing.T) {
	const source, companion = byte(0x20), byte(0x21)
	seededAt := time.Unix(1, 0)
	registry := NewDeviceRegistry(nil)

	registry.RegisterStaticSeed(qualifiedIdentity(source, "SN-TOPOLOGY-A"), SlotRoleSlave, seededAt)
	registry.RegisterStaticSeed(qualifiedIdentity(companion, "SN-TOPOLOGY-A"), SlotRoleSlave, seededAt)
	if err := registry.AliasAddresses(source, companion); err != nil {
		t.Fatal(err)
	}
	registry.RegisterStaticSeed(qualifiedIdentity(companion, "SN-TOPOLOGY-B"), SlotRoleSlave, seededAt)
	registry.RegisterStaticSeed(qualifiedIdentity(companion, "SN-TOPOLOGY-A"), SlotRoleSlave, seededAt)

	// This is fresh explicit topology evidence, recorded after the split.
	if err := registry.AliasAddresses(source, companion); err != nil {
		t.Fatal(err)
	}
	registry.Register(qualifiedIdentity(source, "SN-TOPOLOGY-A"))

	slot, ok := registry.LookupSlotSnapshot(companion)
	if !ok || slot.DiscoverySource != DiscoverySourceStaticSeed || slot.VerificationState != VerificationStateIdentityConfirmed || !slot.FirstObservedAt.Equal(seededAt) {
		t.Fatalf("fresh-topology companion = %#v, present=%v; want static identity-confirmed", slot, ok)
	}
}
