package registry

import (
	"sync"
	"testing"
	"time"
)

func directQualifiedIdentity(address byte, serial string) DeviceInfo {
	return DeviceInfo{
		Address:      address,
		Manufacturer: "Acme",
		DeviceID:     "VR_71",
		SerialNumber: serial,
	}
}

func currentQualifiedIdentityWitness(registry *DeviceRegistry, address byte) (QualifiedIdentityWitness, bool) {
	var witness QualifiedIdentityWitness
	ok := registry.WithCurrentQualifiedIdentityWitness(address, func(current QualifiedIdentityWitness) {
		witness = current
	})
	return witness, ok
}

func TestQualifiedIdentityWitness_DirectExactObservationOnly(t *testing.T) {
	registry := NewDeviceRegistry(nil)
	if _, ok := currentQualifiedIdentityWitness(registry, 0x10); ok {
		t.Fatal("unobserved address produced a witness")
	}

	registry.Register(DeviceInfo{Address: 0x10, Manufacturer: " acme ", DeviceID: "vr_71", SerialNumber: " sn-1 "})
	witness, ok := currentQualifiedIdentityWitness(registry, 0x10)
	if !ok {
		t.Fatal("direct complete observation did not produce a witness")
	}
	if witness.Address != 0x10 || witness.IdentityAuthority != (QualifiedIdentityAuthority{Manufacturer: "ACME", DeviceID: "VR_71", SerialNumber: "SN-1"}) {
		t.Fatalf("witness authority = %#v; want exact normalized authority", witness)
	}
	if witness.ObservationProvenance != QualifiedIdentityObservationProvenanceDirect || !witness.Current || !witness.Immutable || witness.RegistryObservationGeneration == 0 || witness.RegistryProofGeneration == 0 {
		t.Fatalf("witness contract fields = %#v; want direct/current/immutable positive generations", witness)
	}

	registry.RegisterStaticSeed(directQualifiedIdentity(0x11, "SN-SEED"), SlotRoleSlave, time.Unix(1, 0))
	registry.RegisterPassiveObserved(directQualifiedIdentity(0x12, "SN-PASSIVE"), SlotRoleSlave, time.Unix(1, 0))
	registry.Register(DeviceInfo{Address: 0x13, Manufacturer: "Acme", DeviceID: "VR_71"})
	registry.confirmScanIdentity(0x13)
	registry.Register(directQualifiedIdentity(0x14, "0xFFFFFFFF"))
	for _, address := range []byte{0x11, 0x12, 0x13, 0x14} {
		if _, ok := currentQualifiedIdentityWitness(registry, address); ok {
			t.Fatalf("non-direct or incomplete address 0x%02X produced a witness", address)
		}
	}
}

func TestQualifiedIdentityWitness_SparseAndChangedAuthorityLifecycle(t *testing.T) {
	for _, change := range []struct {
		name string
		info DeviceInfo
	}{
		{name: "manufacturer", info: DeviceInfo{Address: 0x10, Manufacturer: "Other"}},
		{name: "device ID", info: DeviceInfo{Address: 0x10, DeviceID: "VR_72"}},
		{name: "serial", info: DeviceInfo{Address: 0x10, SerialNumber: "SN-2"}},
	} {
		t.Run(change.name, func(t *testing.T) {
			registry := NewDeviceRegistry(nil)
			registry.Register(directQualifiedIdentity(0x10, "SN-1"))
			first, ok := currentQualifiedIdentityWitness(registry, 0x10)
			if !ok {
				t.Fatal("setup did not produce a witness")
			}

			registry.Register(DeviceInfo{Address: 0x10, SoftwareVersion: "0204"})
			sparse, ok := currentQualifiedIdentityWitness(registry, 0x10)
			if !ok || sparse.RegistryObservationGeneration != first.RegistryObservationGeneration || sparse.RegistryProofGeneration != first.RegistryProofGeneration {
				t.Fatalf("sparse LKG refresh = %#v, present=%v; want unchanged direct witness %#v", sparse, ok, first)
			}

			registry.Register(change.info)
			if _, ok := currentQualifiedIdentityWitness(registry, 0x10); ok {
				t.Fatal("changed partial triple member retained a witness")
			}

			registry.Register(directQualifiedIdentity(0x10, "SN-3"))
			fresh, ok := currentQualifiedIdentityWitness(registry, 0x10)
			if !ok || fresh.RegistryProofGeneration <= first.RegistryProofGeneration || fresh.RegistryObservationGeneration <= first.RegistryObservationGeneration {
				t.Fatalf("fresh direct observation = %#v, present=%v; want newer witness than %#v", fresh, ok, first)
			}
		})
	}
}

func TestQualifiedIdentityWitness_AliasAndReplacementDoNotSubstituteAuthority(t *testing.T) {
	registry := NewDeviceRegistry(nil)
	registry.Register(directQualifiedIdentity(0x10, "SN-1"))
	first, ok := currentQualifiedIdentityWitness(registry, 0x10)
	if !ok {
		t.Fatal("setup did not produce a witness")
	}
	if err := registry.AliasAddresses(0x10, 0x11); err != nil {
		t.Fatal(err)
	}
	if _, ok := currentQualifiedIdentityWitness(registry, 0x11); ok {
		t.Fatal("topology alias propagated another address's witness")
	}

	replacement := directQualifiedIdentity(0x10, "SN-2")
	registry.Register(replacement)
	current, ok := currentQualifiedIdentityWitness(registry, 0x10)
	if !ok || current.IdentityAuthority.SerialNumber != "SN-2" || current.RegistryProofGeneration <= first.RegistryProofGeneration {
		t.Fatalf("replacement witness = %#v, present=%v; want fresh exact-address SN-2 witness", current, ok)
	}
	registry.mu.RLock()
	if registry.qualifiedIdentityWitnessCurrentLocked(0x10, first) {
		t.Error("cached pre-replacement witness remained current")
	}
	registry.mu.RUnlock()

	registry.RegisterStaticSeed(directQualifiedIdentity(0x10, "SN-3"), SlotRoleSlave, time.Unix(2, 0))
	if _, ok := currentQualifiedIdentityWitness(registry, 0x10); ok {
		t.Fatal("non-direct replacement left a current witness")
	}
}

func TestQualifiedIdentityWitness_ConflictedAuthorityCannotProduceWitness(t *testing.T) {
	registry := NewDeviceRegistry(nil)
	incumbent := directQualifiedIdentity(0x10, "SN-1")
	incumbent.MacAddress = "02:00:00:00:00:01"
	registry.Register(incumbent)
	if _, ok := currentQualifiedIdentityWitness(registry, incumbent.Address); !ok {
		t.Fatal("setup direct authority did not produce a witness")
	}

	conflict := incumbent
	conflict.Address = 0x11
	conflict.MacAddress = "02:00:00:00:00:02"
	registry.Register(conflict)
	if _, ok := currentQualifiedIdentityWitness(registry, conflict.Address); ok {
		t.Fatal("conflicting qualified authority produced a witness")
	}
}

func TestQualifiedIdentityWitness_CurrentnessChecksExactAddressAndGenerations(t *testing.T) {
	registry := NewDeviceRegistry(nil)
	registry.Register(directQualifiedIdentity(0x10, "SN-1"))
	witness, ok := currentQualifiedIdentityWitness(registry, 0x10)
	if !ok {
		t.Fatal("setup did not produce a witness")
	}

	registry.mu.RLock()
	if !registry.qualifiedIdentityWitnessCurrentLocked(0x10, witness) {
		t.Error("issued witness was not current")
	}
	wrongAddress := witness
	wrongAddress.Address = 0x11
	if registry.qualifiedIdentityWitnessCurrentLocked(0x11, wrongAddress) {
		t.Error("wrong address witness was current")
	}
	wrongObservationGeneration := witness
	wrongObservationGeneration.RegistryObservationGeneration++
	if registry.qualifiedIdentityWitnessCurrentLocked(0x10, wrongObservationGeneration) {
		t.Error("wrong observation generation was current")
	}
	wrongProofGeneration := witness
	wrongProofGeneration.RegistryProofGeneration++
	if registry.qualifiedIdentityWitnessCurrentLocked(0x10, wrongProofGeneration) {
		t.Error("wrong proof generation was current")
	}
	registry.mu.RUnlock()
}

func TestQualifiedIdentityWitness_ValueIsolationAndConcurrentAccess(t *testing.T) {
	registry := NewDeviceRegistry(nil)
	registry.Register(directQualifiedIdentity(0x10, "SN-0"))
	witness, ok := currentQualifiedIdentityWitness(registry, 0x10)
	if !ok {
		t.Fatal("setup did not produce a witness")
	}
	witness.IdentityAuthority.Manufacturer = "MUTATED"
	witness.Current = false
	witness.Immutable = false
	isolated, ok := currentQualifiedIdentityWitness(registry, 0x10)
	if !ok || isolated.IdentityAuthority.Manufacturer != "ACME" || !isolated.Current || !isolated.Immutable {
		t.Fatalf("caller mutation leaked into registry: %#v, present=%v", isolated, ok)
	}

	var writersAndReaders sync.WaitGroup
	start := make(chan struct{})
	writersAndReaders.Add(2)
	go func() {
		defer writersAndReaders.Done()
		<-start
		for i := 0; i < 200; i++ {
			registry.Register(directQualifiedIdentity(0x10, "SN-WRITE"))
		}
	}()
	go func() {
		defer writersAndReaders.Done()
		<-start
		for i := 0; i < 200; i++ {
			registry.WithCurrentQualifiedIdentityWitness(0x10, func(current QualifiedIdentityWitness) {
				if current.Address != 0x10 || current.IdentityAuthority.Manufacturer != "ACME" || current.RegistryObservationGeneration == 0 || current.RegistryProofGeneration == 0 {
					t.Errorf("concurrent witness = %#v; want complete exact value", current)
				}
				current.IdentityAuthority.Manufacturer = "CALLER-MUTATION"
			})
		}
	}()
	close(start)
	writersAndReaders.Wait()
}
