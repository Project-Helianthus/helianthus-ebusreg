package registry

import (
	"context"
	"testing"
	"time"

	ebuserrors "github.com/Project-Helianthus/helianthus-ebusgo/errors"
	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
)

const scanIdentitySerial = "21-22-09-0020184848-0082-005409-N4"

func requireIdentityConfirmed(t *testing.T, registry *DeviceRegistry, address byte) AddressSlotSnapshot {
	t.Helper()
	slot, ok := registry.LookupSlotSnapshot(address)
	if !ok || slot.DiscoverySource != DiscoverySourceActiveConfirmed || slot.VerificationState != VerificationStateIdentityConfirmed {
		t.Fatalf("slot %02x = %#v, present=%v; want active-confirmed/identity-confirmed", address, slot, ok)
	}
	return slot
}

func TestScanDirected0704ConfirmsRespondingFaceWithoutSerial(t *testing.T) {
	registry := NewDeviceRegistry(nil)

	if _, err := Scan(context.Background(), &vaillantScanIDTimeoutBus{}, registry, 0x10, []byte{0x20}); err != nil {
		t.Fatal(err)
	}
	requireIdentityConfirmed(t, registry, 0x30)
	// 0x30 responded to the probe of 0x20. The target has no serial proof of
	// its own, so its confirmation is possible only through that fresh edge.
	requireIdentityConfirmed(t, registry, 0x20)
}

func TestScanDirected0704ConfirmsStaticRespondingFaceWithoutReplacingTimestamp(t *testing.T) {
	registry := NewDeviceRegistry(nil)
	seededAt := time.Unix(1, 0)
	registry.RegisterStaticSeed(DeviceInfo{Address: 0x30, Manufacturer: "Vaillant", DeviceID: "DEV30"}, SlotRoleSlave, seededAt)

	if _, err := Scan(context.Background(), &vaillantScanIDTimeoutBus{}, registry, 0x10, []byte{0x20}); err != nil {
		t.Fatal(err)
	}
	if slot := requireIdentityConfirmed(t, registry, 0x30); !slot.FirstObservedAt.Equal(seededAt) {
		t.Fatalf("first observation = %v; want retained seed timestamp %v", slot.FirstObservedAt, seededAt)
	}
}

func TestScanDirected0704RecordsFreshTopologyAfterIdentitySplit(t *testing.T) {
	const source, target = byte(0x30), byte(0x20)
	seededAt := time.Unix(1, 0)
	registry := NewDeviceRegistry(nil)
	info := func(address byte, serial string) DeviceInfo {
		return DeviceInfo{Address: address, Manufacturer: "Vaillant", DeviceID: "DEV30", SerialNumber: serial}
	}

	registry.RegisterStaticSeed(info(source, "SN-SCAN-TOPOLOGY-A"), SlotRoleSlave, seededAt)
	registry.RegisterStaticSeed(info(target, "SN-SCAN-TOPOLOGY-A"), SlotRoleSlave, seededAt)
	if err := registry.AliasAddresses(source, target); err != nil {
		t.Fatal(err)
	}
	registry.RegisterStaticSeed(info(target, "SN-SCAN-TOPOLOGY-B"), SlotRoleSlave, seededAt)
	registry.RegisterStaticSeed(info(target, "SN-SCAN-TOPOLOGY-A"), SlotRoleSlave, seededAt)

	// The valid directed response supplies a new source-target edge after the
	// split, so it may confirm its target despite the retired old edge.
	if _, err := Scan(context.Background(), &vaillantScanIDTimeoutBus{}, registry, 0x10, []byte{target}); err != nil {
		t.Fatal(err)
	}
	if slot := requireIdentityConfirmed(t, registry, source); !slot.FirstObservedAt.Equal(seededAt) {
		t.Fatalf("source first observation = %v; want retained seed timestamp %v", slot.FirstObservedAt, seededAt)
	}
	if slot := requireIdentityConfirmed(t, registry, target); !slot.FirstObservedAt.Equal(seededAt) {
		t.Fatalf("target first observation = %v; want retained seed timestamp %v", slot.FirstObservedAt, seededAt)
	}
}

func TestScanDirected0704ConfirmsOnlyCurrentTopologyAndNotExactTripleMembers(t *testing.T) {
	registry := NewDeviceRegistry(nil)
	seededAt := time.Unix(1, 0)
	passiveAt := time.Unix(2, 0)
	for _, address := range []byte{0x40, 0x41} {
		info := DeviceInfo{Address: address, Manufacturer: "Vaillant", DeviceID: "DEV30", SerialNumber: scanIdentitySerial}
		if address == 0x40 {
			registry.RegisterStaticSeed(info, SlotRoleSlave, seededAt)
		} else {
			registry.RegisterPassiveObserved(info, SlotRoleMaster, passiveAt)
		}
	}

	if _, err := Scan(context.Background(), &vaillantScanIDBus{}, registry, 0x10, []byte{0x20}); err != nil {
		t.Fatal(err)
	}
	entry, ok := registry.Lookup(0x30)
	if !ok || entry.SerialNumber() != scanIdentitySerial {
		t.Fatalf("serial-success control = %#v, present=%v; want %q", entry, ok, scanIdentitySerial)
	}
	// The source and target form explicit current scan topology. The source is
	// directly observed; the target is the only face that may be propagated.
	requireIdentityConfirmed(t, registry, 0x30)
	requireIdentityConfirmed(t, registry, 0x20)

	static, _ := registry.LookupSlotSnapshot(0x40)
	if static.DiscoverySource != DiscoverySourceStaticSeed || static.VerificationState != VerificationStateCandidate || !static.FirstObservedAt.Equal(seededAt) {
		t.Fatalf("unrelated exact-triple static face = %#v", static)
	}
	passive, _ := registry.LookupSlotSnapshot(0x41)
	if passive.DiscoverySource != DiscoverySourcePassiveObserved || passive.VerificationState != VerificationStateCorroborated || !passive.FirstObservedAt.Equal(passiveAt) {
		t.Fatalf("unrelated exact-triple passive face = %#v", passive)
	}
}

func TestScanDirected0704RejectsInvalidOrUnavailableResponsesWithoutConfirmation(t *testing.T) {
	valid := &protocol.Frame{
		Source: 0x20, Target: 0x10, Primary: scanPrimary, Secondary: scanSecondary,
		Data: []byte{0xB5, 'N', 'E', 'T', 'X', '3', 0x01, 0x29, 0x04, 0x04},
	}
	for _, scenario := range []struct {
		name     string
		response *protocol.Frame
		err      error
	}{
		{name: "timeout", err: ebuserrors.ErrTimeout},
		{name: "malformed", response: &protocol.Frame{Source: 0x20, Target: 0x10, Primary: scanPrimary, Secondary: scanSecondary, Data: []byte{0xB5}}},
		{name: "wrong direction", response: &protocol.Frame{Source: valid.Source, Target: 0x71, Primary: valid.Primary, Secondary: valid.Secondary, Data: valid.Data}},
		{name: "unspecified source", response: &protocol.Frame{Target: valid.Target, Primary: valid.Primary, Secondary: valid.Secondary, Data: valid.Data}},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			registry := NewDeviceRegistry(nil)
			seededAt := time.Unix(1, 0)
			registry.RegisterStaticSeed(DeviceInfo{Address: 0x20, Manufacturer: "Vaillant", DeviceID: "NETX3"}, SlotRoleSlave, seededAt)
			bus := &mockScanBus{responses: map[byte]*protocol.Frame{0x20: scenario.response}, errors: map[byte]error{0x20: scenario.err}}
			if _, err := Scan(context.Background(), bus, registry, 0x10, []byte{0x20}); err != nil {
				t.Fatal(err)
			}
			slot, ok := registry.LookupSlotSnapshot(0x20)
			if !ok || slot.DiscoverySource != DiscoverySourceStaticSeed || slot.VerificationState != VerificationStateCandidate || !slot.FirstObservedAt.Equal(seededAt) {
				t.Fatalf("invalid response changed seeded face: %#v, present=%v", slot, ok)
			}
		})
	}
}
