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
