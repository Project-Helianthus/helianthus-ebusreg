package registry

import (
	"testing"
	"time"
)

type observationGenerationProvider struct{}

func (observationGenerationProvider) Name() string                    { return "generation-test" }
func (observationGenerationProvider) Match(DeviceInfo) bool           { return false }
func (observationGenerationProvider) CreatePlanes(DeviceInfo) []Plane { return nil }

func readObservationGeneration(t *testing.T, reg *DeviceRegistry) uint64 {
	t.Helper()
	var generation uint64
	if !reg.WithObservationGeneration(func(current uint64) { generation = current }) {
		t.Fatal("WithObservationGeneration rejected a valid callback")
	}
	return generation
}

func TestWithObservationGenerationRejectsNilInputs(t *testing.T) {
	var nilRegistry *DeviceRegistry
	if nilRegistry.WithObservationGeneration(func(uint64) {}) {
		t.Fatal("nil registry accepted an observation callback")
	}
	if NewDeviceRegistry(nil).WithObservationGeneration(nil) {
		t.Fatal("nil observation callback accepted")
	}
}

func TestObservationGenerationAdvancesForPublicMutations(t *testing.T) {
	reg := NewDeviceRegistry(nil)
	want := uint64(0)
	assertGeneration := func(name string) {
		t.Helper()
		want++
		if got := readObservationGeneration(t, reg); got != want {
			t.Fatalf("%s generation = %d; want %d", name, got, want)
		}
	}

	reg.RegisterProvider(observationGenerationProvider{})
	assertGeneration("RegisterProvider")
	reg.Register(DeviceInfo{Address: 0x10, Manufacturer: "Vaillant", DeviceID: "BASV2"})
	assertGeneration("Register")
	reg.RegisterStaticSeed(DeviceInfo{Address: 0x15, Manufacturer: "Vaillant", DeviceID: "BASV2"}, SlotRoleUnknown, time.Unix(1, 0))
	assertGeneration("RegisterStaticSeed")
	reg.MarkSlotStaticSeed(0x15, SlotRoleUnknown, time.Unix(2, 0))
	assertGeneration("MarkSlotStaticSeed")
	reg.MarkSlotPassiveObserved(0x15, SlotRoleUnknown, time.Unix(3, 0))
	assertGeneration("MarkSlotPassiveObserved")
	reg.RegisterPassiveObserved(DeviceInfo{Address: 0x26, Manufacturer: "Vaillant", DeviceID: "VR_71"}, SlotRoleUnknown, time.Unix(4, 0))
	assertGeneration("RegisterPassiveObserved")
	if err := reg.AliasAddresses(0x10, 0x15); err != nil {
		t.Fatalf("AliasAddresses: %v", err)
	}
	assertGeneration("AliasAddresses")
}

func TestWithObservationGenerationSerializesConcurrentWriter(t *testing.T) {
	reg := NewDeviceRegistry(nil)
	reg.Register(DeviceInfo{Address: 0x10, Manufacturer: "Vaillant", DeviceID: "BASV2"})

	var insideGeneration uint64
	if !reg.WithObservationGeneration(func(current uint64) {
		insideGeneration = current
		if reg.mu.TryLock() {
			reg.mu.Unlock()
			t.Error("registry write lock acquired inside observation read critical section")
		}
	}) {
		t.Fatal("WithObservationGeneration rejected a valid callback")
	}
	reg.Register(DeviceInfo{Address: 0x26, Manufacturer: "Vaillant", DeviceID: "VR_71"})

	if after := readObservationGeneration(t, reg); after != insideGeneration+1 {
		t.Fatalf("generation after serialized writer = %d; want %d", after, insideGeneration+1)
	}
}
