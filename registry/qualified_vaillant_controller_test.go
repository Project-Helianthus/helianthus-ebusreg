package registry

import (
	"sync"
	"testing"
	"time"
)

func directQualifiedBASV2(address byte, serial string) DeviceInfo {
	return DeviceInfo{
		Address: address, Manufacturer: " Vaillant ", DeviceID: " basv2 ", SerialNumber: serial,
		SoftwareVersion: " 0507 ", HardwareVersion: "1704",
	}
}

func TestQualifiedVaillantController_ExactDirectTupleOnly(t *testing.T) {
	registry := NewDeviceRegistry(nil)
	registry.Register(directQualifiedBASV2(0x15, "SN-1"))
	controller, ok := registry.CurrentQualifiedVaillantController(0x15)
	if !ok {
		t.Fatal("exact direct BASV2 observation did not qualify")
	}
	if controller.IdentityAuthority != (QualifiedIdentityAuthority{Manufacturer: "VAILLANT", DeviceID: "BASV2", SerialNumber: "SN-1"}) ||
		controller.Product != "BASV2" || controller.SoftwareVersion != "0507" || controller.HardwareVersion != "1704" ||
		controller.Address != 0x15 || controller.ObservedAt.IsZero() {
		t.Fatalf("controller = %#v; want exact normalized current tuple", controller)
	}
	if controller.RuleRevision != qualifiedVaillantControllerRuleRevision || controller.NativeEvidence != (NativeEvidenceReference{Kind: "protocol.observation", Reference: qualifiedVaillantControllerEvidenceRef, SHA256: qualifiedVaillantControllerEvidenceSHA}) ||
		!controller.Current || !controller.Immutable || controller.RegistryObservationGeneration == 0 || controller.RegistryProofGeneration == 0 {
		t.Fatalf("controller contract fields = %#v", controller)
	}

	for _, candidate := range []DeviceInfo{
		{Address: 0x15, Manufacturer: "Vaillant", DeviceID: "BASV2", SerialNumber: "SN-2", SoftwareVersion: "0508", HardwareVersion: "1704"},
		{Address: 0x15, Manufacturer: "Vaillant", DeviceID: "BASV2", SerialNumber: "SN-3", SoftwareVersion: "0507", HardwareVersion: "1705"},
		{Address: 0x15, Manufacturer: "Vaillant", DeviceID: "CTLV2", SerialNumber: "SN-4", SoftwareVersion: "0507", HardwareVersion: "1704"},
		{Address: 0x15, Manufacturer: "Vaillant", DeviceID: "BASV2", SerialNumber: "0xFFFFFFFF", SoftwareVersion: "0507", HardwareVersion: "1704"},
		{Address: 0x15, Manufacturer: "Vaillant", DeviceID: "BASV2", SerialNumber: "SN-5", SoftwareVersion: "", HardwareVersion: "1704"},
	} {
		registry = NewDeviceRegistry(nil)
		registry.Register(candidate)
		if _, ok := registry.CurrentQualifiedVaillantController(0x15); ok {
			t.Fatalf("unsupported or incomplete candidate qualified: %#v", candidate)
		}
	}
}

func TestQualifiedVaillantController_SparseAndLifecycleRetirement(t *testing.T) {
	registry := NewDeviceRegistry(nil)
	registry.Register(directQualifiedBASV2(0x15, "SN-1"))
	first, ok := registry.CurrentQualifiedVaillantController(0x15)
	if !ok {
		t.Fatal("setup did not qualify")
	}

	registry.Register(DeviceInfo{Address: 0x15, SoftwareVersion: "0507"})
	sparse, ok := registry.CurrentQualifiedVaillantController(0x15)
	if !ok || sparse.RegistryObservationGeneration != first.RegistryObservationGeneration || sparse.RegistryProofGeneration != first.RegistryProofGeneration {
		t.Fatalf("sparse matching refresh = %#v, present=%v; want retained direct proof %#v", sparse, ok, first)
	}

	registry.Register(DeviceInfo{Address: 0x15, SoftwareVersion: "0508"})
	if _, ok := registry.CurrentQualifiedVaillantController(0x15); ok {
		t.Fatal("changed software version retained qualification")
	}
	registry.Register(directQualifiedBASV2(0x15, "SN-2"))
	second, ok := registry.CurrentQualifiedVaillantController(0x15)
	if !ok || second.RegistryProofGeneration <= first.RegistryProofGeneration {
		t.Fatalf("fresh direct observation = %#v, present=%v; want newer proof", second, ok)
	}

	registry.RegisterStaticSeed(directQualifiedBASV2(0x15, "SN-SEED"), SlotRoleSlave, time.Unix(1, 0))
	if _, ok := registry.CurrentQualifiedVaillantController(0x15); ok {
		t.Fatal("non-direct replacement retained exact-address qualification")
	}
}

func TestQualifiedVaillantController_ConflictRetiresWithoutReplacement(t *testing.T) {
	registry := NewDeviceRegistry(nil)
	incumbent := directQualifiedBASV2(0x15, "SN-1")
	incumbent.MacAddress = "02:00:00:00:00:01"
	registry.Register(incumbent)
	if _, ok := registry.CurrentQualifiedVaillantController(0x15); !ok {
		t.Fatal("setup did not qualify")
	}
	conflict := incumbent
	conflict.MacAddress = "02:00:00:00:00:02"
	registry.Register(conflict)
	if _, ok := registry.CurrentQualifiedVaillantController(0x15); ok {
		t.Fatal("conflicting observation retained current qualification")
	}
}

func TestQualifiedVaillantController_ValueIsolationAndConcurrentReadRetire(t *testing.T) {
	registry := NewDeviceRegistry(nil)
	registry.Register(directQualifiedBASV2(0x15, "SN-1"))
	controller, ok := registry.CurrentQualifiedVaillantController(0x15)
	if !ok {
		t.Fatal("setup did not qualify")
	}
	controller.IdentityAuthority.Manufacturer = "MUTATED"
	controller.NativeEvidence.SHA256 = "MUTATED"
	controller.Current = false
	isolate, ok := registry.CurrentQualifiedVaillantController(0x15)
	if !ok || isolate.IdentityAuthority.Manufacturer != "VAILLANT" || isolate.NativeEvidence.SHA256 != qualifiedVaillantControllerEvidenceSHA || !isolate.Current {
		t.Fatalf("caller mutation leaked: %#v, present=%v", isolate, ok)
	}

	var workers sync.WaitGroup
	start := make(chan struct{})
	workers.Add(2)
	go func() {
		defer workers.Done()
		<-start
		for i := 0; i < 200; i++ {
			registry.Register(directQualifiedBASV2(0x15, "SN-WRITE"))
		}
	}()
	go func() {
		defer workers.Done()
		<-start
		for i := 0; i < 200; i++ {
			if current, ok := registry.CurrentQualifiedVaillantController(0x15); ok {
				if current.Address != 0x15 || current.IdentityAuthority.Manufacturer != "VAILLANT" || current.RegistryProofGeneration == 0 {
					t.Errorf("concurrent controller = %#v", current)
				}
			}
		}
	}()
	close(start)
	workers.Wait()
}
