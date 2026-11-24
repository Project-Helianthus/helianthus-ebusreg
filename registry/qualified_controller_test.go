package registry

import "testing"

// fakeControllerRule declares one identity: manufacturer Acme, device VR_71 and
// the raw product field 2103. It is a caller-supplied rule, as in production.
func fakeControllerRule(info DeviceInfo) ControllerDeclaration {
	if info.Manufacturer != "Acme" || info.DeviceID != "VR_71" || info.HardwareVersion != "2103" {
		return ControllerDeclaration{Reason: "unlisted_pair"}
	}
	return ControllerDeclaration{
		Declared:      true,
		Family:        "VR_71",
		ProductNumber: "0141",
		Level:         "declared",
		EvidenceGrade: "Hypothesis",
		RuleDigest:    "test-rule",
		Rows:          []DeclaredRow{{Group: 0, Register: 0x0073, Element: "x", EvidenceGrade: "Hypothesis", Level: "declared"}},
	}
}

func controllerRegistry(t *testing.T, address byte) *DeviceRegistry {
	t.Helper()
	r := NewDeviceRegistry(nil)
	r.Register(DeviceInfo{Address: address, Manufacturer: "Acme", DeviceID: "VR_71", SerialNumber: "SN-1", SoftwareVersion: "0123", HardwareVersion: "2103"})
	return r
}

func currentController(r *DeviceRegistry, address byte) (QualifiedController, bool) {
	var got QualifiedController
	ok := r.WithCurrentQualifiedController(address, fakeControllerRule, func(q QualifiedController) { got = q })
	return got, ok
}

func TestQualifiedController_ExactAddressAndWitnessGenerations(t *testing.T) {
	r := controllerRegistry(t, 0x10)
	got, ok := currentController(r, 0x10)
	if !ok {
		t.Fatal("exact address produced no controller result")
	}
	witness, _ := currentQualifiedIdentityWitness(r, 0x10)
	if got.Address != 0x10 || got.DeviceID != "VR_71" || got.SoftwareVersion != "0123" || got.ProductNumberRaw != "2103" {
		t.Fatalf("result identity = %#v", got)
	}
	if got.Declaration.Family != "VR_71" || got.Declaration.RuleDigest != "test-rule" || len(got.Declaration.Rows) != 1 {
		t.Fatalf("declaration = %#v", got.Declaration)
	}
	if got.RegistryObservationGeneration != witness.RegistryObservationGeneration || got.RegistryProofGeneration != witness.RegistryProofGeneration {
		t.Fatalf("result generations = %d/%d; want witness %d/%d", got.RegistryObservationGeneration, got.RegistryProofGeneration, witness.RegistryObservationGeneration, witness.RegistryProofGeneration)
	}
	if _, ok := currentController(r, 0x11); ok {
		t.Fatal("address without an entry produced a controller result")
	}
}

func TestQualifiedController_UndeclaredIdentityHasNoResult(t *testing.T) {
	r := controllerRegistry(t, 0x10)
	r.Register(DeviceInfo{Address: 0x12, Manufacturer: "Acme", DeviceID: "VR_72", SerialNumber: "SN-2", HardwareVersion: "2103"})
	if _, ok := currentController(r, 0x12); ok {
		t.Fatal("undeclared identity produced a controller result")
	}
}

func TestQualifiedController_InvalidInputLeavesRegistryUnchanged(t *testing.T) {
	r := controllerRegistry(t, 0x10)
	before, _ := currentQualifiedIdentityWitness(r, 0x10)
	var obsBefore uint64
	r.WithObservationGeneration(func(g uint64) { obsBefore = g })

	if r.WithCurrentQualifiedController(0x10, nil, func(QualifiedController) {}) {
		t.Fatal("nil rule returned a result")
	}
	if r.WithCurrentQualifiedController(0x10, fakeControllerRule, nil) {
		t.Fatal("nil callback returned a result")
	}
	var nilRegistry *DeviceRegistry
	if nilRegistry.WithCurrentQualifiedController(0x10, fakeControllerRule, func(QualifiedController) {}) {
		t.Fatal("nil registry returned a result")
	}
	rejecting := func(DeviceInfo) ControllerDeclaration { return ControllerDeclaration{Reason: "incomplete"} }
	if r.WithCurrentQualifiedController(0x10, rejecting, func(QualifiedController) { t.Fatal("rejected rule invoked callback") }) {
		t.Fatal("rejecting rule returned a result")
	}

	after, _ := currentQualifiedIdentityWitness(r, 0x10)
	var obsAfter uint64
	r.WithObservationGeneration(func(g uint64) { obsAfter = g })
	if after != before || obsAfter != obsBefore {
		t.Fatalf("registry state changed: witness %#v -> %#v, observation %d -> %d", before, after, obsBefore, obsAfter)
	}
}

// A conflicting direct observation replaces the witness with one for the new
// authority, and the old declared result is not returned for it.
func TestQualifiedController_ConflictingObservationRetiresWithWitness(t *testing.T) {
	r := controllerRegistry(t, 0x10)
	first, ok := currentController(r, 0x10)
	if !ok {
		t.Fatal("setup produced no result")
	}
	before, _ := currentQualifiedIdentityWitness(r, 0x10)
	r.Register(DeviceInfo{Address: 0x10, Manufacturer: "Acme", DeviceID: "VR_72", SerialNumber: "SN-1", HardwareVersion: "2103"})
	after, ok := currentQualifiedIdentityWitness(r, 0x10)
	if !ok || after.IdentityAuthority.DeviceID != "VR_72" {
		t.Fatalf("witness after conflict = %#v, present=%v; want the new authority", after, ok)
	}
	if after.RegistryProofGeneration == before.RegistryProofGeneration {
		t.Fatal("conflicting direct observation did not advance the proof generation")
	}
	if second, ok := currentController(r, 0x10); ok {
		t.Fatalf("old declared result returned for the new authority: %#v (previous %#v)", second, first)
	}
}

func TestQualifiedController_SparseRefreshKeepsResult(t *testing.T) {
	r := controllerRegistry(t, 0x10)
	first, _ := currentController(r, 0x10)
	r.Register(DeviceInfo{Address: 0x10, SoftwareVersion: "0204"})
	sparse, ok := currentController(r, 0x10)
	if !ok {
		t.Fatal("matching sparse refresh retired the controller result")
	}
	if sparse.RegistryObservationGeneration != first.RegistryObservationGeneration || sparse.RegistryProofGeneration != first.RegistryProofGeneration {
		t.Fatalf("sparse refresh changed generations: %#v -> %#v", first, sparse)
	}
}

// Retiring the result does not end the witness: a rule that stops declaring
// the identity hides the result only.
func TestQualifiedController_RetiringResultDoesNotEndWitness(t *testing.T) {
	r := controllerRegistry(t, 0x10)
	declineAll := func(DeviceInfo) ControllerDeclaration { return ControllerDeclaration{Reason: "unlisted_pair"} }
	if r.WithCurrentQualifiedController(0x10, declineAll, func(QualifiedController) {}) {
		t.Fatal("declining rule returned a result")
	}
	if _, ok := currentQualifiedIdentityWitness(r, 0x10); !ok {
		t.Fatal("retiring the result ended the witness")
	}
}

// Without a complete identity there is no witness; the family rule needs no
// serial, so the result is declared with zero witness generations.
func TestQualifiedController_NoSerialDeclaredWithoutWitness(t *testing.T) {
	r := NewDeviceRegistry(nil)
	r.Register(DeviceInfo{Address: 0x13, Manufacturer: "Acme", DeviceID: "VR_71", SoftwareVersion: "0123", HardwareVersion: "2103"})
	if _, ok := currentQualifiedIdentityWitness(r, 0x13); ok {
		t.Fatal("serial-less identity produced a witness")
	}
	got, ok := currentController(r, 0x13)
	if !ok {
		t.Fatal("serial-less declared identity produced no result")
	}
	if got.RegistryObservationGeneration != 0 || got.RegistryProofGeneration != 0 {
		t.Fatalf("serial-less result generations = %d/%d; want zero", got.RegistryObservationGeneration, got.RegistryProofGeneration)
	}
}

// The result's rows are a copy: changing them cannot change a later result.
func TestQualifiedController_RowsAreCopied(t *testing.T) {
	r := controllerRegistry(t, 0x10)
	got, _ := currentController(r, 0x10)
	got.Declaration.Rows[0].Register = 0xFFFF
	again, _ := currentController(r, 0x10)
	if again.Declaration.Rows[0].Register != 0x0073 {
		t.Fatal("mutating a result changed a later result")
	}
}
