package vaillant

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func ab43Decode(ab43T *testing.T, ab43Data []byte) map[string]any {
	ab43T.Helper()
	var ab43Value map[string]any
	if ab43Error := json.Unmarshal(ab43Data, &ab43Value); ab43Error != nil {
		ab43T.Fatal(ab43Error)
	}
	return ab43Value
}

func ab43Expected(ab43T *testing.T) map[string]any {
	ab43T.Helper()
	ab43Value := ab43Decode(ab43T, []byte(`{
 "id":"ab43_simulator_adc","representation_grade":"Proven","candidate_grade":"Hypothesis",
 "physical_qualification":"Unknown","coordinate_system":"logical_including_prefix_not_physical_offsets",
 "limitations":["partial_catalog","remaining_selectors_unknown","address_space_unknown","support_unknown","zero_not_established_sentinel","signed_carrier_does_not_admit_negative_values","formal_and_operational_domains_are_distinct","no_runtime_registration_or_wire_builder"],
 "entries":[]}`))
	ab43Request := []string{
		`{"name":"relay","logical_position":3,"width":1,"raw_domain":[0,255],"signed":false,"byte_order":null,"conversion":null,"unit":null,"bit_roles":null}`,
		`{"name":"PWM1","logical_position":4,"width":1,"raw_domain":[0,255],"signed":false,"byte_order":null,"conversion":{"divide_by":2},"unit":"percent"}`,
		`{"name":"PWM2","logical_position":5,"width":1,"raw_domain":[0,255],"signed":false,"byte_order":null,"conversion":{"divide_by":2},"unit":"percent"}`,
		`{"name":"stepper","logical_position":6,"width":1,"raw_domain":[0,255],"signed":false,"byte_order":null,"conversion":null,"unit":"steps","mappings":{"255":"calibrate/go-to-zero"},"steps_domain":[0,240],"unmapped_regions":[[241,254]]}`,
		`{"name":"LED","logical_position":7,"width":1,"raw_domain":[0,255],"signed":false,"byte_order":null,"conversion":null,"unit":null,"mappings":{"0":"disabled","1":"enabled"},"unmapped_regions":[[2,255]]}`,
		`{"name":"ED","logical_position":8,"width":1,"raw_domain":[0,255],"signed":false,"byte_order":null,"conversion":{"divide_by":2},"unit":"percent"}`,
		`{"name":"temperature-table","logical_position":9,"width":1,"raw_domain":[0,255],"signed":false,"byte_order":null,"conversion":null,"unit":null,"mappings":{"0":{"label":"standard3K3","candidate_celsius_range":[-13.5,116]},"1":{"label":"extended1K","candidate_celsius_range":[-20,155]},"2":{"label":"low-temperature15K","candidate_celsius_range":[-39,80]},"3":{"label":"standard1K","candidate_celsius_range":[-13.9,130]}},"unmapped_regions":[[4,255]]}`,
	}
	ab43Positive := []string{
		`{"name":"ADC","logical_position":3,"width":2,"raw_domain":[0,1023],"signed":false,"byte_order":"little_endian","conversion":null,"unit":null}`,
		`{"name":"temperature1","logical_position":5,"width":2,"raw_domain":[0,3200],"signed":true,"byte_order":"little_endian","conversion":{"divide_by":16},"unit":"degrees_C","candidate_physical_domain":[0,200],"sentinel":null}`,
		`{"name":"temperature2","logical_position":7,"width":2,"raw_domain":[0,3200],"signed":true,"byte_order":"little_endian","conversion":{"divide_by":16},"unit":"degrees_C","candidate_physical_domain":[0,200],"sentinel":null}`,
		`{"name":"temperature3","logical_position":9,"width":2,"raw_domain":[0,3200],"signed":true,"byte_order":"little_endian","conversion":{"divide_by":16},"unit":"degrees_C","candidate_physical_domain":[0,200],"sentinel":null}`,
		`{"name":"temperature4","logical_position":11,"width":2,"raw_domain":[0,3200],"signed":true,"byte_order":"little_endian","conversion":{"divide_by":16},"unit":"degrees_C","candidate_physical_domain":[0,200],"sentinel":null}`,
		`{"name":"reserved","logical_position":13,"width":1,"raw_domain":[0,255],"signed":false,"byte_order":null,"conversion":null,"unit":null,"reserved_regions":[[13,1]]}`,
		`{"name":"impeller","logical_position":14,"width":2,"raw_domain":[0,1000],"signed":false,"byte_order":"little_endian","conversion":null,"unit":"Hz"}`,
		`{"name":"vortex","logical_position":16,"width":2,"raw_domain":[0,350],"signed":false,"byte_order":"little_endian","conversion":null,"unit":"Hz","candidate_operational_domain":[5,350],"candidate_hysteresis_hz":{"on":5,"off":4},"sentinel":null}`,
	}
	for ab43Index, ab43Alternative := range []string{"request", "positive", "negative"} {
		ab43Row := ab43Decode(ab43T, []byte(`{"family":"Vaillant simulator","variant":null,"service":"B520","operation":0,"physical_extent":null,"grade":"Hypothesis","raw_fallback":"preserve_raw_bytes","reserved_regions":[],"remaining_selectors":null,"address_space":null,"support":null}`))
		ab43Row["id"] = "ab43_b520_op00_" + ab43Alternative
		ab43Row["alternative"] = ab43Alternative
		ab43Row["direction"] = "response"
		ab43Row["logical_extent"] = nil
		ab43Fields := []any{}
		var ab43Declarations []string
		switch ab43Index {
		case 0:
			ab43Row["direction"] = "request"
			ab43Row["logical_extent"] = float64(10)
			ab43Declarations = ab43Request
		case 1:
			ab43Row["logical_extent"] = float64(18)
			ab43Declarations = ab43Positive
			ab43Row["reserved_regions"] = []any{[]any{float64(13), float64(1)}}
		}
		for _, ab43Declaration := range ab43Declarations {
			ab43Field := ab43Decode(ab43T, []byte(ab43Declaration))
			ab43Field["default"] = float64(0)
			ab43Field["grade"] = "Hypothesis"
			ab43Field["raw_fallback"] = "preserve_raw_bytes"
			if _, ab43Present := ab43Field["mappings"]; !ab43Present {
				ab43Field["mappings"] = map[string]any{}
			}
			if _, ab43Present := ab43Field["reserved_regions"]; !ab43Present {
				ab43Field["reserved_regions"] = []any{}
			}
			ab43Fields = append(ab43Fields, ab43Field)
		}
		ab43Row["fields"] = ab43Fields
		ab43Value["entries"] = append(ab43Value["entries"].([]any), ab43Row)
	}
	return ab43Value
}

func TestAb43ExactCatalog(ab43T *testing.T) {
	ab43Got := ab43Decode(ab43T, Ab43SimulatorADC())
	ab43Want := ab43Expected(ab43T)
	if !reflect.DeepEqual(ab43Got, ab43Want) {
		ab43T.Error("ab43 catalog differs from finite declaration contract")
	}
	for ab43Index, ab43Entry := range ab43Want["entries"].([]any) {
		ab43Row := ab43Entry.(map[string]any)
		ab43T.Run("ab43_"+ab43Row["alternative"].(string), func(ab43T *testing.T) {
			ab43Actual := ab43Got["entries"].([]any)[ab43Index].(map[string]any)
			for ab43Key, ab43Value := range ab43Row {
				if ab43Key == "fields" {
					continue
				}
				ab43T.Run("ab43_"+ab43Key, func(ab43T *testing.T) {
					if ab43GotValue, ab43Present := ab43Actual[ab43Key]; !ab43Present || !reflect.DeepEqual(ab43GotValue, ab43Value) {
						ab43T.Fatalf("ab43 %s: got %#v want %#v", ab43Key, ab43Actual[ab43Key], ab43Value)
					}
				})
			}
			for ab43FieldIndex, ab43ExpectedField := range ab43Row["fields"].([]any) {
				ab43Field := ab43ExpectedField.(map[string]any)
				ab43T.Run("ab43_"+ab43Field["name"].(string), func(ab43T *testing.T) {
					ab43ActualField := ab43Actual["fields"].([]any)[ab43FieldIndex].(map[string]any)
					for ab43Key, ab43Value := range ab43Field {
						ab43T.Run("ab43_"+ab43Key, func(ab43T *testing.T) {
							if ab43GotValue, ab43Present := ab43ActualField[ab43Key]; !ab43Present || !reflect.DeepEqual(ab43GotValue, ab43Value) {
								ab43T.Fatalf("ab43 %s: got %#v want %#v", ab43Key, ab43ActualField[ab43Key], ab43Value)
							}
						})
					}
				})
			}
		})
	}
}

func TestAb43IdentityAxes(ab43T *testing.T) {
	ab43Want := ab43Expected(ab43T)
	for ab43Index := range ab43Want["entries"].([]any) {
		for _, ab43Axis := range []string{"family", "variant", "direction", "service", "operation", "alternative"} {
			for _, ab43Mode := range []string{"vary", "drop"} {
				ab43T.Run(fmt.Sprintf("ab43_%d_%s_%s", ab43Index, ab43Axis, ab43Mode), func(ab43T *testing.T) {
					ab43Control := ab43Decode(ab43T, Ab43SimulatorADC())
					if !reflect.DeepEqual(ab43Control, ab43Want) {
						ab43T.Fatal("ab43 invalid control")
					}
					ab43Row := ab43Control["entries"].([]any)[ab43Index].(map[string]any)
					if ab43Mode == "drop" {
						delete(ab43Row, ab43Axis)
					} else if ab43Axis == "operation" {
						ab43Row[ab43Axis] = float64(1)
					} else {
						ab43Row[ab43Axis] = "ab43_other"
					}
					ab43Bytes, ab43Error := json.Marshal(ab43Control)
					if ab43Error != nil {
						ab43T.Fatal(ab43Error)
					}
					if reflect.DeepEqual(ab43Decode(ab43T, ab43Bytes), ab43Want) {
						ab43T.Fatal("ab43 changed identity axis survived exact comparison")
					}
				})
			}
		}
	}
	ab43Collapsed := ab43Decode(ab43T, Ab43SimulatorADC())
	ab43Collapsed["entries"].([]any)[2] = ab43Collapsed["entries"].([]any)[1]
	if reflect.DeepEqual(ab43Collapsed, ab43Want) {
		ab43T.Fatal("ab43 collapsed alternatives survived comparison")
	}
}

func TestAb43IndependentBytes(ab43T *testing.T) {
	ab43First := Ab43SimulatorADC()
	ab43Before := bytes.Clone(ab43First)
	for ab43Index := range ab43First {
		ab43First[ab43Index] ^= 255
	}
	if !bytes.Equal(Ab43SimulatorADC(), ab43Before) {
		ab43T.Fatal("ab43 accessor aliases production bytes")
	}
	if !reflect.DeepEqual(ab43Decode(ab43T, Ab43SimulatorADC()), ab43Expected(ab43T)) {
		ab43T.Fatal("ab43 accessor lost production catalog")
	}
}

func ab43FieldAt(ab43T *testing.T, ab43Alternative int, ab43Name string) map[string]any {
	ab43T.Helper()
	for _, ab43Value := range ab43Decode(ab43T, Ab43SimulatorADC())["entries"].([]any)[ab43Alternative].(map[string]any)["fields"].([]any) {
		ab43Field := ab43Value.(map[string]any)
		if ab43Field["name"] == ab43Name {
			return ab43Field
		}
	}
	ab43T.Fatal("ab43 field missing: " + ab43Name)
	return nil
}

func ab43Inside(ab43Domain any, ab43Raw float64) bool {
	ab43Bounds := ab43Domain.([]any)
	return ab43Raw >= ab43Bounds[0].(float64) && ab43Raw <= ab43Bounds[1].(float64)
}

func ab43Candidate(ab43Field map[string]any, ab43Raw []byte) float64 {
	ab43Value := float64(ab43Raw[0])
	if ab43Field["width"] == float64(2) {
		var ab43Word uint16
		if ab43Field["byte_order"] == "little_endian" {
			ab43Word = binary.LittleEndian.Uint16(ab43Raw)
		} else {
			ab43Word = binary.BigEndian.Uint16(ab43Raw)
		}
		ab43Value = float64(ab43Word)
		if ab43Field["signed"] == true {
			ab43Value = float64(int16(ab43Word))
		}
	}
	if ab43Conversion, ab43OK := ab43Field["conversion"].(map[string]any); ab43OK {
		ab43Value /= ab43Conversion["divide_by"].(float64)
	}
	return ab43Value
}

func TestAb43SyntheticBoundaries(ab43T *testing.T) {
	for _, ab43Name := range []string{"PWM1", "PWM2", "ED"} {
		ab43T.Run("ab43_"+ab43Name, func(ab43T *testing.T) {
			ab43Field := ab43FieldAt(ab43T, 0, ab43Name)
			if ab43Candidate(ab43Field, []byte{255}) != 127.5 || !ab43Inside(ab43Field["raw_domain"], 255) {
				ab43T.Fatal("ab43 percent clipped or unavailable")
			}
		})
	}
	ab43Stepper := ab43FieldAt(ab43T, 0, "stepper")
	for _, ab43Case := range []struct {
		ab43Raw      float64
		ab43Steps    bool
		ab43Unmapped bool
		ab43Label    any
	}{
		{240, true, false, nil}, {241, false, true, nil}, {254, false, true, nil}, {255, false, false, "calibrate/go-to-zero"},
	} {
		ab43T.Run(fmt.Sprintf("ab43_stepper_%g", ab43Case.ab43Raw), func(ab43T *testing.T) {
			if ab43Inside(ab43Stepper["steps_domain"], ab43Case.ab43Raw) != ab43Case.ab43Steps || ab43Inside(ab43Stepper["unmapped_regions"].([]any)[0], ab43Case.ab43Raw) != ab43Case.ab43Unmapped || ab43Stepper["mappings"].(map[string]any)[fmt.Sprint(ab43Case.ab43Raw)] != ab43Case.ab43Label {
				ab43T.Fatal("ab43 stepper alternatives collapsed")
			}
		})
	}
	for _, ab43Name := range []string{"temperature1", "temperature2", "temperature3", "temperature4"} {
		ab43T.Run("ab43_"+ab43Name, func(ab43T *testing.T) {
			ab43Field := ab43FieldAt(ab43T, 1, ab43Name)
			if ab43Candidate(ab43Field, []byte{0, 12}) != 192 || !ab43Inside(ab43Field["raw_domain"], 3072) || !ab43Inside(ab43Field["candidate_physical_domain"], 192) {
				ab43T.Fatal("ab43 valid temperature changed")
			}
			if ab43Candidate(ab43Field, []byte{255, 255}) != -0.0625 || ab43Inside(ab43Field["raw_domain"], -1) || ab43Inside(ab43Field["candidate_physical_domain"], -0.0625) || ab43Field["sentinel"] != nil {
				ab43T.Fatal("ab43 negative carrier admitted or sentinel invented")
			}
		})
	}
	ab43ADC := ab43FieldAt(ab43T, 1, "ADC")
	for _, ab43Case := range []struct {
		ab43Raw   []byte
		ab43Value float64
		ab43In    bool
	}{{[]byte{255, 3}, 1023, true}, {[]byte{0, 4}, 1024, false}} {
		if ab43Candidate(ab43ADC, ab43Case.ab43Raw) != ab43Case.ab43Value || ab43Inside(ab43ADC["raw_domain"], ab43Case.ab43Value) != ab43Case.ab43In {
			ab43T.Fatal("ab43 ADC boundary changed")
		}
	}
	ab43Vortex := ab43FieldAt(ab43T, 1, "vortex")
	for _, ab43Case := range []struct {
		ab43Raw         float64
		ab43Formal      bool
		ab43Operational bool
	}{{0, true, false}, {4, true, false}, {5, true, true}, {350, true, true}, {351, false, false}} {
		ab43T.Run(fmt.Sprintf("ab43_vortex_%g", ab43Case.ab43Raw), func(ab43T *testing.T) {
			if ab43Inside(ab43Vortex["raw_domain"], ab43Case.ab43Raw) != ab43Case.ab43Formal || ab43Inside(ab43Vortex["candidate_operational_domain"], ab43Case.ab43Raw) != ab43Case.ab43Operational || ab43Vortex["sentinel"] != nil {
				ab43T.Fatal("ab43 vortex domains conflated")
			}
		})
	}
}
