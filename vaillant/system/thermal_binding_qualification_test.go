package system

import (
	"testing"

	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
)

func currentController(t *testing.T, deviceID, productRaw string) registry.QualifiedController {
	t.Helper()
	reg := registry.NewDeviceRegistry(nil)
	reg.Register(registry.DeviceInfo{
		Address:         0x15,
		Manufacturer:    "Vaillant",
		DeviceID:        deviceID,
		SerialNumber:    "SN-BIND-" + deviceID,
		SoftwareVersion: "0123",
		HardwareVersion: productRaw,
	})
	var got registry.QualifiedController
	if !reg.WithCurrentQualifiedController(0x15, DeclareController, func(q registry.QualifiedController) { got = q }) {
		t.Fatalf("no controller result for %s", deviceID)
	}
	return got
}

// TestBindingControllerRowsOfDeclaredFamily checks each row the record lists
// for the family 70000 (product 0141) and the family 72000 (product 0179).
func TestBindingControllerRowsOfDeclaredFamily(t *testing.T) {
	controller70 := currentController(t, "70000", "2103")
	b := BindControllerRow(controller70, ControllerRowRef{Group: 0x00, Register: 0x0073})
	if !b.Declared || b.Element != "thermal.measurement.outdoor_air_temperature" || b.EvidenceGrade != "Hypothesis" || b.Level != "declared" || b.Family != "70000" || b.RuleDigest != ControllerQualificationRuleDigest() {
		t.Fatalf("binding 70000 0x0073 = %#v", b)
	}
	if b := BindControllerRow(controller70, ControllerRowRef{Group: 0x01, Register: 0x0004}); !b.Declared || b.EvidenceGrade != "Proven" {
		t.Fatalf("binding 70000 DHW setpoint = %#v", b)
	}
	if b := BindControllerRow(controller70, ControllerRowRef{Group: 0x00, Register: 0x0096}); b.Declared || b.Reason != BindingReasonRowNotListed {
		t.Fatalf("unlisted row 0x0096 = %#v; want row_not_listed", b)
	}
	if b := BindControllerRow(controller70, ControllerRowRef{Group: 0x00, Instance: 0x01, Register: 0x0073}); b.Declared || b.Reason != BindingReasonInstanceNotListed {
		t.Fatalf("instance 01 = %#v; want instance_not_listed", b)
	}

	controller72 := currentController(t, "72000", "7703")
	if b := BindControllerRow(controller72, ControllerRowRef{Group: 0x00, Register: 0x004B}); !b.Declared || b.Family != "72000" {
		t.Fatalf("binding 72000 0x004B = %#v", b)
	}
	if b := BindControllerRow(controller72, ControllerRowRef{Group: 0x03, Register: 0x0022}); !b.Declared || b.Element != "thermal.setpoint.temperature_setpoint" {
		t.Fatalf("72000 zone setpoint 0x0022 = %#v", b)
	}
	// 0x0022 lists 72000 but not 70000: the same row is not declared for 70000.
	if b := BindControllerRow(controller70, ControllerRowRef{Group: 0x03, Register: 0x0022}); b.Declared || b.Reason != BindingReasonRowNotListed {
		t.Fatalf("70000 zone setpoint 0x0022 = %#v; want row_not_listed", b)
	}
}

// TestBindingNotDeclaredWithoutCurrentController gives the named missing
// declaration when the controller result is absent.
func TestBindingNotDeclaredWithoutCurrentController(t *testing.T) {
	b := BindControllerRow(registry.QualifiedController{}, ControllerRowRef{Group: 0x00, Register: 0x0073})
	if b.Declared || b.Reason != BindingReasonNotDeclared {
		t.Fatalf("binding without controller = %#v; want not_declared", b)
	}
}

// TestBindingBoilerFlowRow checks the boiler flow row through a current target.
func TestBindingBoilerFlowRow(t *testing.T) {
	reg := registry.NewDeviceRegistry(nil)
	reg.Register(registry.DeviceInfo{Address: 0x08, Manufacturer: "Vaillant", DeviceID: "BAI00"})
	declared := registry.DeviceInfo{Manufacturer: "Vaillant", DeviceID: "BAI00"}
	var target registry.QualifiedTarget
	if !reg.WithCurrentQualifiedTarget(0x08, declared, DeclareBoilerTarget, func(q registry.QualifiedTarget) { target = q }) {
		t.Fatal("no boiler target result")
	}
	b := BindBoilerFlow(target, BoilerFlowRowOpcode, BoilerFlowRowSelector)
	if !b.Declared || b.Element != "thermal.measurement.heat_generator_flow_temperature" || b.Level != "declared" || b.RuleDigest != ControllerQualificationRuleDigest() {
		t.Fatalf("boiler flow binding = %#v", b)
	}
	if b := BindBoilerFlow(target, 0x0D, 0x1801); b.Declared || b.Reason != BindingReasonRowNotListed {
		t.Fatalf("other selector = %#v; want row_not_listed", b)
	}
	if b := BindBoilerFlow(registry.QualifiedTarget{}, BoilerFlowRowOpcode, BoilerFlowRowSelector); b.Declared || b.Reason != BindingReasonNotDeclared {
		t.Fatalf("binding without target = %#v", b)
	}
}
