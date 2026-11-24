package registry

import "testing"

// fakeTargetRule declares a boiler target by manufacturer Acme and device
// BOIL00. It is a caller-supplied rule, as in production.
func fakeTargetRule(info DeviceInfo) TargetDeclaration {
	if info.Manufacturer != "Acme" || info.DeviceID != "BOIL00" {
		return TargetDeclaration{Reason: "not_boiler_class"}
	}
	return TargetDeclaration{
		Declared:          true,
		Class:             "boiler_target",
		Level:             "declared",
		RuleDigest:        "test-rule",
		FlowOpcode:        0x0D,
		FlowSelector:      0x1800,
		FlowElement:       "x",
		FlowEvidenceGrade: "Hypothesis",
	}
}

var boundBoiler = DeviceInfo{Manufacturer: "Acme", DeviceID: "BOIL00"}

func currentTarget(r *DeviceRegistry, address byte, declared DeviceInfo) (QualifiedTarget, bool) {
	var got QualifiedTarget
	ok := r.WithCurrentQualifiedTarget(address, declared, fakeTargetRule, func(q QualifiedTarget) { got = q })
	return got, ok
}

// A declared target needs no serial: it is current with no complete identity
// witness.
func TestQualifiedTarget_CurrentWithoutCompleteWitness(t *testing.T) {
	r := NewDeviceRegistry(nil)
	r.Register(DeviceInfo{Address: 0x20, Manufacturer: "Acme", DeviceID: "BOIL00"})
	if _, ok := currentQualifiedIdentityWitness(r, 0x20); ok {
		t.Fatal("serial-less identity produced a witness")
	}
	got, ok := currentTarget(r, 0x20, boundBoiler)
	if !ok {
		t.Fatal("declared target without a witness produced no result")
	}
	if got.Address != 0x20 || got.DeviceID != "BOIL00" || got.Declaration.FlowOpcode != 0x0D || got.Declaration.FlowSelector != 0x1800 {
		t.Fatalf("target = %#v", got)
	}
}

// The result retires when the address no longer resolves to the bound device,
// with or without a complete identity witness.
func TestQualifiedTarget_RetiresWhenAddressStopsResolving(t *testing.T) {
	t.Run("no witness", func(t *testing.T) {
		r := NewDeviceRegistry(nil)
		r.Register(DeviceInfo{Address: 0x20, Manufacturer: "Acme", DeviceID: "BOIL00"})
		r.Register(DeviceInfo{Address: 0x20, Manufacturer: "Acme", DeviceID: "BOIL01"})
		if _, ok := currentTarget(r, 0x20, boundBoiler); ok {
			t.Fatal("target kept after the address changed to another device")
		}
	})
	t.Run("complete witness", func(t *testing.T) {
		r := NewDeviceRegistry(nil)
		r.Register(DeviceInfo{Address: 0x20, Manufacturer: "Acme", DeviceID: "BOIL00", SerialNumber: "SN-B"})
		if _, ok := currentTarget(r, 0x20, boundBoiler); !ok {
			t.Fatal("setup produced no target")
		}
		r.Register(DeviceInfo{Address: 0x20, Manufacturer: "Acme", DeviceID: "BOIL01", SerialNumber: "SN-C"})
		if _, ok := currentTarget(r, 0x20, boundBoiler); ok {
			t.Fatal("target kept after its witness ended")
		}
	})
}

func TestQualifiedTarget_UnboundOrUndeclaredHasNoResult(t *testing.T) {
	r := NewDeviceRegistry(nil)
	r.Register(DeviceInfo{Address: 0x20, Manufacturer: "Acme", DeviceID: "BOIL00"})
	if _, ok := currentTarget(r, 0x20, DeviceInfo{Manufacturer: "Acme", DeviceID: "BOIL01"}); ok {
		t.Fatal("result returned for a device the caller did not bind")
	}
	if _, ok := currentTarget(r, 0x21, boundBoiler); ok {
		t.Fatal("address without an entry produced a target")
	}
	r.Register(DeviceInfo{Address: 0x22, Manufacturer: "Acme", DeviceID: "HEAT00"})
	if _, ok := currentTarget(r, 0x22, DeviceInfo{Manufacturer: "Acme", DeviceID: "HEAT00"}); ok {
		t.Fatal("undeclared identity produced a target")
	}
	var nilRegistry *DeviceRegistry
	if nilRegistry.WithCurrentQualifiedTarget(0x20, boundBoiler, fakeTargetRule, func(QualifiedTarget) {}) {
		t.Fatal("nil registry returned a target")
	}
	if r.WithCurrentQualifiedTarget(0x20, boundBoiler, nil, func(QualifiedTarget) {}) {
		t.Fatal("nil rule returned a target")
	}
}
