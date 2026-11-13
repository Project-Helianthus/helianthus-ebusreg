package vaillant

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func simulatorADCDecode(simulatorADCT *testing.T, simulatorADCData []byte) map[string]any {
	simulatorADCT.Helper()
	var simulatorADCValue map[string]any
	if simulatorADCError := json.Unmarshal(simulatorADCData, &simulatorADCValue); simulatorADCError != nil {
		simulatorADCT.Fatal(simulatorADCError)
	}
	return simulatorADCValue
}

func simulatorADCExpected(simulatorADCT *testing.T) map[string]any {
	simulatorADCT.Helper()
	simulatorADCValue := simulatorADCDecode(simulatorADCT, []byte(`{
 "id":"ab43_simulator_adc","representation_grade":"Proven","candidate_grade":"Hypothesis",
 "physical_qualification":"Unknown","coordinate_system":"logical_including_prefix_not_physical_offsets",
 "limitations":["partial_catalog","remaining_selectors_unknown","address_space_unknown","support_unknown","zero_not_established_sentinel","signed_carrier_does_not_admit_negative_values","formal_and_operational_domains_are_distinct","no_runtime_registration_or_wire_builder"],
 "entries":[]}`))
	simulatorADCRequest := []string{
		`{"name":"relay","logical_position":3,"width":1,"raw_domain":[0,255],"signed":false,"byte_order":null,"conversion":null,"unit":null,"bit_roles":null}`,
		`{"name":"PWM1","logical_position":4,"width":1,"raw_domain":[0,255],"signed":false,"byte_order":null,"conversion":{"divide_by":2},"unit":"percent"}`,
		`{"name":"PWM2","logical_position":5,"width":1,"raw_domain":[0,255],"signed":false,"byte_order":null,"conversion":{"divide_by":2},"unit":"percent"}`,
		`{"name":"stepper","logical_position":6,"width":1,"raw_domain":[0,255],"signed":false,"byte_order":null,"conversion":null,"unit":"steps","mappings":{"255":"calibrate/go-to-zero"},"steps_domain":[0,240],"unmapped_regions":[[241,254]]}`,
		`{"name":"LED","logical_position":7,"width":1,"raw_domain":[0,255],"signed":false,"byte_order":null,"conversion":null,"unit":null,"mappings":{"0":"disabled","1":"enabled"},"unmapped_regions":[[2,255]]}`,
		`{"name":"ED","logical_position":8,"width":1,"raw_domain":[0,255],"signed":false,"byte_order":null,"conversion":{"divide_by":2},"unit":"percent"}`,
		`{"name":"temperature-table","logical_position":9,"width":1,"raw_domain":[0,255],"signed":false,"byte_order":null,"conversion":null,"unit":null,"mappings":{"0":{"label":"standard3K3","candidate_celsius_range":[-13.5,116]},"1":{"label":"extended1K","candidate_celsius_range":[-20,155]},"2":{"label":"low-temperature15K","candidate_celsius_range":[-39,80]},"3":{"label":"standard1K","candidate_celsius_range":[-13.9,130]}},"unmapped_regions":[[4,255]]}`,
	}
	simulatorADCPositive := []string{
		`{"name":"ADC","logical_position":3,"width":2,"raw_domain":[0,1023],"signed":false,"byte_order":"little_endian","conversion":null,"unit":null}`,
		`{"name":"temperature1","logical_position":5,"width":2,"raw_domain":[0,3200],"signed":true,"byte_order":"little_endian","conversion":{"divide_by":16},"unit":"degrees_C","candidate_physical_domain":[0,200],"sentinel":null}`,
		`{"name":"temperature2","logical_position":7,"width":2,"raw_domain":[0,3200],"signed":true,"byte_order":"little_endian","conversion":{"divide_by":16},"unit":"degrees_C","candidate_physical_domain":[0,200],"sentinel":null}`,
		`{"name":"temperature3","logical_position":9,"width":2,"raw_domain":[0,3200],"signed":true,"byte_order":"little_endian","conversion":{"divide_by":16},"unit":"degrees_C","candidate_physical_domain":[0,200],"sentinel":null}`,
		`{"name":"temperature4","logical_position":11,"width":2,"raw_domain":[0,3200],"signed":true,"byte_order":"little_endian","conversion":{"divide_by":16},"unit":"degrees_C","candidate_physical_domain":[0,200],"sentinel":null}`,
		`{"name":"reserved","logical_position":13,"width":1,"raw_domain":[0,255],"signed":false,"byte_order":null,"conversion":null,"unit":null,"reserved_regions":[[13,1]]}`,
		`{"name":"impeller","logical_position":14,"width":2,"raw_domain":[0,1000],"signed":false,"byte_order":"little_endian","conversion":null,"unit":"Hz"}`,
		`{"name":"vortex","logical_position":16,"width":2,"raw_domain":[0,350],"signed":false,"byte_order":"little_endian","conversion":null,"unit":"Hz","candidate_operational_domain":[5,350],"candidate_hysteresis_hz":{"on":5,"off":4},"sentinel":null}`,
	}
	for simulatorADCIndex, simulatorADCAlternative := range []string{"request", "positive", "negative"} {
		simulatorADCRow := simulatorADCDecode(simulatorADCT, []byte(`{"family":"Vaillant simulator","variant":null,"service":"B520","operation":0,"physical_extent":null,"grade":"Hypothesis","raw_fallback":"preserve_raw_bytes","reserved_regions":[],"remaining_selectors":null,"address_space":null,"support":null}`))
		simulatorADCRow["id"] = "ab43_b520_op00_" + simulatorADCAlternative
		simulatorADCRow["alternative"] = simulatorADCAlternative
		simulatorADCRow["direction"] = "response"
		simulatorADCRow["logical_extent"] = nil
		simulatorADCFields := []any{}
		var simulatorADCDeclarations []string
		switch simulatorADCIndex {
		case 0:
			simulatorADCRow["direction"] = "request"
			simulatorADCRow["logical_extent"] = float64(10)
			simulatorADCDeclarations = simulatorADCRequest
		case 1:
			simulatorADCRow["logical_extent"] = float64(18)
			simulatorADCDeclarations = simulatorADCPositive
			simulatorADCRow["reserved_regions"] = []any{[]any{float64(13), float64(1)}}
		}
		for _, simulatorADCDeclaration := range simulatorADCDeclarations {
			simulatorADCField := simulatorADCDecode(simulatorADCT, []byte(simulatorADCDeclaration))
			simulatorADCField["default"] = float64(0)
			simulatorADCField["grade"] = "Hypothesis"
			simulatorADCField["raw_fallback"] = "preserve_raw_bytes"
			if _, simulatorADCPresent := simulatorADCField["mappings"]; !simulatorADCPresent {
				simulatorADCField["mappings"] = map[string]any{}
			}
			if _, simulatorADCPresent := simulatorADCField["reserved_regions"]; !simulatorADCPresent {
				simulatorADCField["reserved_regions"] = []any{}
			}
			simulatorADCFields = append(simulatorADCFields, simulatorADCField)
		}
		simulatorADCRow["fields"] = simulatorADCFields
		simulatorADCValue["entries"] = append(simulatorADCValue["entries"].([]any), simulatorADCRow)
	}
	return simulatorADCValue
}

func TestAb43ExactCatalog(simulatorADCT *testing.T) {
	simulatorADCGot := simulatorADCDecode(simulatorADCT, SimulatorADC())
	simulatorADCWant := simulatorADCExpected(simulatorADCT)
	if !reflect.DeepEqual(simulatorADCGot, simulatorADCWant) {
		simulatorADCT.Error("ab43 catalog differs from finite declaration contract")
	}
	for simulatorADCIndex, simulatorADCEntry := range simulatorADCWant["entries"].([]any) {
		simulatorADCRow := simulatorADCEntry.(map[string]any)
		simulatorADCT.Run("ab43_"+simulatorADCRow["alternative"].(string), func(simulatorADCT *testing.T) {
			simulatorADCActual := simulatorADCGot["entries"].([]any)[simulatorADCIndex].(map[string]any)
			for simulatorADCKey, simulatorADCValue := range simulatorADCRow {
				if simulatorADCKey == "fields" {
					continue
				}
				simulatorADCT.Run("ab43_"+simulatorADCKey, func(simulatorADCT *testing.T) {
					if simulatorADCGotValue, simulatorADCPresent := simulatorADCActual[simulatorADCKey]; !simulatorADCPresent || !reflect.DeepEqual(simulatorADCGotValue, simulatorADCValue) {
						simulatorADCT.Fatalf("ab43 %s: got %#v want %#v", simulatorADCKey, simulatorADCActual[simulatorADCKey], simulatorADCValue)
					}
				})
			}
			for simulatorADCFieldIndex, simulatorADCExpectedField := range simulatorADCRow["fields"].([]any) {
				simulatorADCField := simulatorADCExpectedField.(map[string]any)
				simulatorADCT.Run("ab43_"+simulatorADCField["name"].(string), func(simulatorADCT *testing.T) {
					simulatorADCActualField := simulatorADCActual["fields"].([]any)[simulatorADCFieldIndex].(map[string]any)
					for simulatorADCKey, simulatorADCValue := range simulatorADCField {
						simulatorADCT.Run("ab43_"+simulatorADCKey, func(simulatorADCT *testing.T) {
							if simulatorADCGotValue, simulatorADCPresent := simulatorADCActualField[simulatorADCKey]; !simulatorADCPresent || !reflect.DeepEqual(simulatorADCGotValue, simulatorADCValue) {
								simulatorADCT.Fatalf("ab43 %s: got %#v want %#v", simulatorADCKey, simulatorADCActualField[simulatorADCKey], simulatorADCValue)
							}
						})
					}
				})
			}
		})
	}
}

func TestAb43IdentityAxes(simulatorADCT *testing.T) {
	simulatorADCWant := simulatorADCExpected(simulatorADCT)
	for simulatorADCIndex := range simulatorADCWant["entries"].([]any) {
		for _, simulatorADCAxis := range []string{"family", "variant", "direction", "service", "operation", "alternative"} {
			for _, simulatorADCMode := range []string{"vary", "drop"} {
				simulatorADCT.Run(fmt.Sprintf("ab43_%d_%s_%s", simulatorADCIndex, simulatorADCAxis, simulatorADCMode), func(simulatorADCT *testing.T) {
					simulatorADCControl := simulatorADCDecode(simulatorADCT, SimulatorADC())
					if !reflect.DeepEqual(simulatorADCControl, simulatorADCWant) {
						simulatorADCT.Fatal("ab43 invalid control")
					}
					simulatorADCRow := simulatorADCControl["entries"].([]any)[simulatorADCIndex].(map[string]any)
					if simulatorADCMode == "drop" {
						delete(simulatorADCRow, simulatorADCAxis)
					} else if simulatorADCAxis == "operation" {
						simulatorADCRow[simulatorADCAxis] = float64(1)
					} else {
						simulatorADCRow[simulatorADCAxis] = "ab43_other"
					}
					simulatorADCBytes, simulatorADCError := json.Marshal(simulatorADCControl)
					if simulatorADCError != nil {
						simulatorADCT.Fatal(simulatorADCError)
					}
					if reflect.DeepEqual(simulatorADCDecode(simulatorADCT, simulatorADCBytes), simulatorADCWant) {
						simulatorADCT.Fatal("ab43 changed identity axis survived exact comparison")
					}
				})
			}
		}
	}
	simulatorADCCollapsed := simulatorADCDecode(simulatorADCT, SimulatorADC())
	simulatorADCCollapsed["entries"].([]any)[2] = simulatorADCCollapsed["entries"].([]any)[1]
	if reflect.DeepEqual(simulatorADCCollapsed, simulatorADCWant) {
		simulatorADCT.Fatal("ab43 collapsed alternatives survived comparison")
	}
}

func TestAb43IndependentBytes(simulatorADCT *testing.T) {
	simulatorADCFirst := SimulatorADC()
	simulatorADCBefore := bytes.Clone(simulatorADCFirst)
	for simulatorADCIndex := range simulatorADCFirst {
		simulatorADCFirst[simulatorADCIndex] ^= 255
	}
	if !bytes.Equal(SimulatorADC(), simulatorADCBefore) {
		simulatorADCT.Fatal("ab43 accessor aliases production bytes")
	}
	if !reflect.DeepEqual(simulatorADCDecode(simulatorADCT, SimulatorADC()), simulatorADCExpected(simulatorADCT)) {
		simulatorADCT.Fatal("ab43 accessor lost production catalog")
	}
}

func simulatorADCFieldAt(simulatorADCT *testing.T, simulatorADCAlternative int, simulatorADCName string) map[string]any {
	simulatorADCT.Helper()
	for _, simulatorADCValue := range simulatorADCDecode(simulatorADCT, SimulatorADC())["entries"].([]any)[simulatorADCAlternative].(map[string]any)["fields"].([]any) {
		simulatorADCField := simulatorADCValue.(map[string]any)
		if simulatorADCField["name"] == simulatorADCName {
			return simulatorADCField
		}
	}
	simulatorADCT.Fatal("ab43 field missing: " + simulatorADCName)
	return nil
}

func simulatorADCInside(simulatorADCDomain any, simulatorADCRaw float64) bool {
	simulatorADCBounds := simulatorADCDomain.([]any)
	return simulatorADCRaw >= simulatorADCBounds[0].(float64) && simulatorADCRaw <= simulatorADCBounds[1].(float64)
}

func simulatorADCCandidate(simulatorADCField map[string]any, simulatorADCRaw []byte) float64 {
	simulatorADCValue := float64(simulatorADCRaw[0])
	if simulatorADCField["width"] == float64(2) {
		var simulatorADCWord uint16
		if simulatorADCField["byte_order"] == "little_endian" {
			simulatorADCWord = binary.LittleEndian.Uint16(simulatorADCRaw)
		} else {
			simulatorADCWord = binary.BigEndian.Uint16(simulatorADCRaw)
		}
		simulatorADCValue = float64(simulatorADCWord)
		if simulatorADCField["signed"] == true {
			simulatorADCValue = float64(int16(simulatorADCWord))
		}
	}
	if simulatorADCConversion, simulatorADCOK := simulatorADCField["conversion"].(map[string]any); simulatorADCOK {
		simulatorADCValue /= simulatorADCConversion["divide_by"].(float64)
	}
	return simulatorADCValue
}

func TestAb43SyntheticBoundaries(simulatorADCT *testing.T) {
	for _, simulatorADCName := range []string{"PWM1", "PWM2", "ED"} {
		simulatorADCT.Run("ab43_"+simulatorADCName, func(simulatorADCT *testing.T) {
			simulatorADCField := simulatorADCFieldAt(simulatorADCT, 0, simulatorADCName)
			if simulatorADCCandidate(simulatorADCField, []byte{255}) != 127.5 || !simulatorADCInside(simulatorADCField["raw_domain"], 255) {
				simulatorADCT.Fatal("ab43 percent clipped or unavailable")
			}
		})
	}
	simulatorADCStepper := simulatorADCFieldAt(simulatorADCT, 0, "stepper")
	for _, simulatorADCCase := range []struct {
		simulatorADCRaw      float64
		simulatorADCSteps    bool
		simulatorADCUnmapped bool
		simulatorADCLabel    any
	}{
		{240, true, false, nil}, {241, false, true, nil}, {254, false, true, nil}, {255, false, false, "calibrate/go-to-zero"},
	} {
		simulatorADCT.Run(fmt.Sprintf("ab43_stepper_%g", simulatorADCCase.simulatorADCRaw), func(simulatorADCT *testing.T) {
			if simulatorADCInside(simulatorADCStepper["steps_domain"], simulatorADCCase.simulatorADCRaw) != simulatorADCCase.simulatorADCSteps || simulatorADCInside(simulatorADCStepper["unmapped_regions"].([]any)[0], simulatorADCCase.simulatorADCRaw) != simulatorADCCase.simulatorADCUnmapped || simulatorADCStepper["mappings"].(map[string]any)[fmt.Sprint(simulatorADCCase.simulatorADCRaw)] != simulatorADCCase.simulatorADCLabel {
				simulatorADCT.Fatal("ab43 stepper alternatives collapsed")
			}
		})
	}
	for _, simulatorADCName := range []string{"temperature1", "temperature2", "temperature3", "temperature4"} {
		simulatorADCT.Run("ab43_"+simulatorADCName, func(simulatorADCT *testing.T) {
			simulatorADCField := simulatorADCFieldAt(simulatorADCT, 1, simulatorADCName)
			if simulatorADCCandidate(simulatorADCField, []byte{0, 12}) != 192 || !simulatorADCInside(simulatorADCField["raw_domain"], 3072) || !simulatorADCInside(simulatorADCField["candidate_physical_domain"], 192) {
				simulatorADCT.Fatal("ab43 valid temperature changed")
			}
			if simulatorADCCandidate(simulatorADCField, []byte{255, 255}) != -0.0625 || simulatorADCInside(simulatorADCField["raw_domain"], -1) || simulatorADCInside(simulatorADCField["candidate_physical_domain"], -0.0625) || simulatorADCField["sentinel"] != nil {
				simulatorADCT.Fatal("ab43 negative carrier admitted or sentinel invented")
			}
		})
	}
	simulatorADCADC := simulatorADCFieldAt(simulatorADCT, 1, "ADC")
	for _, simulatorADCCase := range []struct {
		simulatorADCRaw   []byte
		simulatorADCValue float64
		simulatorADCIn    bool
	}{{[]byte{255, 3}, 1023, true}, {[]byte{0, 4}, 1024, false}} {
		if simulatorADCCandidate(simulatorADCADC, simulatorADCCase.simulatorADCRaw) != simulatorADCCase.simulatorADCValue || simulatorADCInside(simulatorADCADC["raw_domain"], simulatorADCCase.simulatorADCValue) != simulatorADCCase.simulatorADCIn {
			simulatorADCT.Fatal("ab43 ADC boundary changed")
		}
	}
	simulatorADCVortex := simulatorADCFieldAt(simulatorADCT, 1, "vortex")
	for _, simulatorADCCase := range []struct {
		simulatorADCRaw         float64
		simulatorADCFormal      bool
		simulatorADCOperational bool
	}{{0, true, false}, {4, true, false}, {5, true, true}, {350, true, true}, {351, false, false}} {
		simulatorADCT.Run(fmt.Sprintf("ab43_vortex_%g", simulatorADCCase.simulatorADCRaw), func(simulatorADCT *testing.T) {
			if simulatorADCInside(simulatorADCVortex["raw_domain"], simulatorADCCase.simulatorADCRaw) != simulatorADCCase.simulatorADCFormal || simulatorADCInside(simulatorADCVortex["candidate_operational_domain"], simulatorADCCase.simulatorADCRaw) != simulatorADCCase.simulatorADCOperational || simulatorADCVortex["sentinel"] != nil {
				simulatorADCT.Fatal("ab43 vortex domains conflated")
			}
		})
	}
}
