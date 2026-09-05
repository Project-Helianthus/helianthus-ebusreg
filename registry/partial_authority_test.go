package registry

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
)

func TestPartialObservationCannotCreateCrossAddressAuthority(t *testing.T) {
	old := DeviceInfo{Address: 0x10, Manufacturer: "Vaillant", DeviceID: "OLD30", SerialNumber: "SN-OLD", MacAddress: "02:00:00:00:00:01"}

	for _, caller := range []struct {
		name     string
		register func(*DeviceRegistry, DeviceInfo)
	}{
		{name: "active", register: func(r *DeviceRegistry, info DeviceInfo) { r.Register(info) }},
		{name: "passive", register: func(r *DeviceRegistry, info DeviceInfo) {
			r.RegisterPassiveObserved(info, SlotRoleSlave, time.Unix(1, 0))
		}},
		{name: "static", register: func(r *DeviceRegistry, info DeviceInfo) { r.RegisterStaticSeed(info, SlotRoleSlave, time.Unix(1, 0)) }},
	} {
		t.Run(caller.name, func(t *testing.T) {
			r := NewDeviceRegistry(nil)
			r.Register(old)
			caller.register(r, DeviceInfo{Address: old.Address, Manufacturer: "Vaillant", DeviceID: "NEW30"})

			local, _ := r.Lookup(old.Address)
			if local.DeviceID() != "NEW30" || local.SerialNumber() != old.SerialNumber {
				t.Fatalf("same-address LKG = (%q, %q); want NEW30/%q", local.DeviceID(), local.SerialNumber(), old.SerialNumber)
			}
			independent := old
			independent.Address, independent.DeviceID = 0x11, "NEW30"
			if joined := r.Register(independent); joined == local {
				t.Fatal("partial model update synthesized cross-address authority")
			}
			if indexed, ok := r.lookupByIdentity(independent); !ok || indexed == local {
				t.Fatal("synthetic retained-field triple became selectable")
			}
		})
	}
}

type partialAuthorityScanBus struct{}

func (partialAuthorityScanBus) Send(_ context.Context, frame protocol.Frame) (*protocol.Frame, error) {
	if frame.Primary == scanPrimary && frame.Secondary == scanSecondary {
		return &protocol.Frame{Source: frame.Target, Target: frame.Source, Primary: frame.Primary, Secondary: frame.Secondary, Data: []byte{0xB5, 'N', 'E', 'W', '3', '0', 0x01, 0x02, 0x03, 0x04}}, nil
	}
	return nil, errors.New("serial unavailable")
}

func TestDirected0704PartialObservationCannotCreateCrossAddressAuthority(t *testing.T) {
	r := NewDeviceRegistry(nil)
	old := DeviceInfo{Address: 0x20, Manufacturer: "Vaillant", DeviceID: "OLD30", SerialNumber: "SN-OLD", MacAddress: "02:00:00:00:00:01"}
	r.Register(old)
	if _, err := Scan(context.Background(), partialAuthorityScanBus{}, r, 0x10, []byte{old.Address}); err != nil {
		t.Fatal(err)
	}

	local, _ := r.Lookup(old.Address)
	independent := old
	independent.Address, independent.DeviceID = 0x21, "NEW30"
	if joined := r.Register(independent); joined == local {
		t.Fatal("directed 07/04 without B5/09 serial created cross-address authority")
	}
}

func TestPartialObservationAuthorityLifecycle(t *testing.T) {
	old := DeviceInfo{Address: 0x30, Manufacturer: "Vaillant", DeviceID: "OLD30", SerialNumber: "SN-OLD", MacAddress: "02:00:00:00:00:01"}

	t.Run("wholly sparse refresh retains legitimate current authority", func(t *testing.T) {
		r := NewDeviceRegistry(nil)
		current := r.Register(old)
		r.Register(DeviceInfo{Address: old.Address, SoftwareVersion: "0204"})
		if indexed, ok := r.lookupByIdentity(old); !ok || indexed != current {
			t.Fatal("wholly sparse refresh retired unchanged current authority")
		}
		other := old
		other.Address = 0x31
		if joined := r.Register(other); joined != current {
			t.Fatal("complete compatible observation did not use retained current authority")
		}
	})

	t.Run("complete correction establishes current authority and retires history", func(t *testing.T) {
		r := NewDeviceRegistry(nil)
		r.Register(old)
		r.Register(DeviceInfo{Address: old.Address, Manufacturer: "Vaillant", DeviceID: "NEW30"})
		if _, ok := r.lookupByIdentity(old); ok {
			t.Fatal("partial model update retained old cross-address authority")
		}

		corrected := old
		corrected.DeviceID = "NEW30"
		current := r.Register(corrected)
		if _, ok := r.lookupByIdentity(old); ok {
			t.Fatal("historical triple remained selectable after complete correction")
		}
		if indexed, ok := r.lookupByIdentity(corrected); !ok || indexed != current {
			t.Fatal("complete current observation did not establish corrected authority")
		}
		other := corrected
		other.Address = 0x31
		if joined := r.Register(other); joined != current {
			t.Fatal("complete compatible correction did not merge")
		}
	})

	t.Run("serial and MAC conflicts retain the incumbent", func(t *testing.T) {
		r := NewDeviceRegistry(nil)
		incumbent := r.Register(old)
		macConflict := old
		macConflict.Address, macConflict.MacAddress = 0x31, "02:00:00:00:00:02"
		if disputed := r.Register(macConflict); disputed == incumbent {
			t.Fatal("conflicting MAC joined incumbent")
		}
		serialConflict := old
		serialConflict.Address, serialConflict.SerialNumber = 0x32, "SN-OTHER"
		if distinct := r.Register(serialConflict); distinct == incumbent {
			t.Fatal("conflicting serial joined incumbent")
		}
		if indexed, ok := r.lookupByIdentity(old); !ok || indexed != incumbent {
			t.Fatal("conflicting evidence displaced incumbent authority")
		}
	})
}

func TestPartialAuthorityAliasOrdersKeepIncumbent(t *testing.T) {
	for _, order := range [][2]byte{{0x41, 0x42}, {0x42, 0x41}} {
		t.Run("alias order", func(t *testing.T) {
			r := NewDeviceRegistry(nil)
			incumbentInfo := DeviceInfo{Address: 0x40, Manufacturer: "Vaillant", DeviceID: "BASV2", SerialNumber: "SN-ALIAS", MacAddress: "02:00:00:00:00:01"}
			incumbent := r.Register(incumbentInfo)
			disputed := incumbentInfo
			disputed.Address, disputed.MacAddress = 0x41, "02:00:00:00:00:02"
			r.Register(disputed)
			r.Register(DeviceInfo{Address: 0x42})
			if err := r.AliasAddresses(order[0], order[1]); err != nil {
				t.Fatal(err)
			}
			if indexed, ok := r.lookupByIdentity(incumbentInfo); !ok || indexed != incumbent {
				t.Fatal("alias rebuild displaced incumbent authority")
			}
		})
	}
}

func TestPartialAuthorityDoesNotReviveStaleTopology(t *testing.T) {
	const source, companion = byte(0x50), byte(0x51)
	seededAt := time.Unix(1, 0)
	r := NewDeviceRegistry(nil)
	old := DeviceInfo{Manufacturer: "Vaillant", DeviceID: "BASV2", SerialNumber: "SN-TOPOLOGY"}
	old.Address = source
	r.RegisterStaticSeed(old, SlotRoleSlave, seededAt)
	old.Address = companion
	r.RegisterStaticSeed(old, SlotRoleSlave, seededAt)
	if err := r.AliasAddresses(source, companion); err != nil {
		t.Fatal(err)
	}

	corrected := old
	corrected.SerialNumber = "SN-TOPOLOGY-CORRECTED"
	r.RegisterStaticSeed(corrected, SlotRoleSlave, seededAt)
	r.RegisterStaticSeed(old, SlotRoleSlave, seededAt)
	r.Register(DeviceInfo{Address: source, Manufacturer: old.Manufacturer, DeviceID: old.DeviceID, SerialNumber: old.SerialNumber})

	slot, ok := r.LookupSlotSnapshot(companion)
	if !ok || slot.DiscoverySource != DiscoverySourceStaticSeed || slot.VerificationState != VerificationStateCandidate || !slot.FirstObservedAt.Equal(seededAt) {
		t.Fatalf("rejoined companion = %#v, present=%v; want retained static candidate without stale topology confirmation", slot, ok)
	}
}
