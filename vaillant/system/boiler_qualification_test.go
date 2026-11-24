package system

import (
	"encoding/json"
	"testing"

	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
)

// TestBoilerQualificationVectors replays every boiler vector of the record.
func TestBoilerQualificationVectors(t *testing.T) {
	fixture := loadQualificationFixture(t)
	checked := 0
	for _, v := range fixture.Vectors {
		if v.Class != "boiler" {
			continue
		}
		checked++
		t.Run(v.ID, func(t *testing.T) {
			got := QualifyBoilerIdentification(0x15, 0x15, mustHexBytes(t, v.PayloadHex))
			switch v.Expect.Outcome {
			case "publish":
				if got.Outcome != BoilerOutcomePublish || got.DeviceID != v.Expect.DeviceID || got.Level != "declared" {
					t.Fatalf("got %#v; want publish for %s", got, v.Expect.DeviceID)
				}
			case "reject_heat_pump_collision":
				if got.Outcome != BoilerOutcomeHeatPumpCollision || got.DeviceID != v.Expect.DeviceID {
					t.Fatalf("got %#v; want heat-pump collision for %s", got, v.Expect.DeviceID)
				}
			case "unknown":
				if got.Outcome != BoilerOutcomeUnknown || got.Reason != v.Expect.Reason {
					t.Fatalf("got %#v; want unknown %s", got, v.Expect.Reason)
				}
			case "reject_without_mutation":
				if got.Outcome != BoilerOutcomeRejectWithoutChange || got.Reason != v.Expect.Reason {
					t.Fatalf("got %#v; want reject %s", got, v.Expect.Reason)
				}
			default:
				t.Fatalf("unhandled outcome %q", v.Expect.Outcome)
			}
		})
	}
	if checked != 5 {
		t.Fatalf("boiler vectors checked = %d; want the 5 of the record", checked)
	}
}

// TestBoilerTargetDeclarationNeedsNoSerial checks the registry declaration of a
// boiler target: the class decides it, with no serial and no software field,
// and a heat-pump identity is refused.
func TestBoilerTargetDeclarationNeedsNoSerial(t *testing.T) {
	boiler := DeclareBoilerTarget(registry.DeviceInfo{Address: 0x08, Manufacturer: "Vaillant", DeviceID: "BAI00"})
	if !boiler.Declared || boiler.FlowOpcode != 0x0D || boiler.FlowSelector != 0x1800 || boiler.FlowQualifiedBinding {
		t.Fatalf("boiler declaration = %#v", boiler)
	}
	if boiler.FlowElement != "thermal.measurement.heat_generator_flow_temperature" || boiler.FlowEvidenceGrade != "Hypothesis" {
		t.Fatalf("flow row = %#v", boiler)
	}
	for _, id := range []string{"EHP00", "HMU00"} {
		if DeclareBoilerTarget(registry.DeviceInfo{Address: 0x08, Manufacturer: "Vaillant", DeviceID: id}).Declared {
			t.Fatalf("heat-pump identity %s declared as a boiler target", id)
		}
	}
	if DeclareBoilerTarget(registry.DeviceInfo{Address: 0x08, Manufacturer: "0x08", DeviceID: "BAI00"}).Declared {
		t.Fatal("other manufacturer declared as a boiler target")
	}
	if DeclareBoilerTarget(registry.DeviceInfo{Address: 0x08, Manufacturer: "Vaillant", DeviceID: "BAI"}).Declared {
		t.Fatal("malformed device ID declared as a boiler target")
	}
}

// TestBoilerRecordFlowRowMatchesRecord checks the constants against the record.
func TestBoilerQualificationRecordFlowRow(t *testing.T) {
	var record struct {
		BoilerTarget struct {
			FlowRow struct {
				Opcode   string `json:"opcode"`
				Selector string `json:"selector"`
			} `json:"flow_row"`
		} `json:"boiler_target"`
	}
	if err := json.Unmarshal(controllerQualificationRuleJSON, &record); err != nil {
		t.Fatal(err)
	}
	if record.BoilerTarget.FlowRow.Opcode != "0x0D" || record.BoilerTarget.FlowRow.Selector != "18 00" {
		t.Fatalf("record flow row = %#v", record.BoilerTarget.FlowRow)
	}
}
