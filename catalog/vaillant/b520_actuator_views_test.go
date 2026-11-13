package vaillant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"
)

func b520ActuatorViewInput(b520ActuatorViewDirection string, b520ActuatorViewBytes []int) map[string]any {
	return map[string]any{
		"context": map[string]any{"primary": 181, "secondary": 32, "operation": 0, "direction": b520ActuatorViewDirection, "view": "logical-object", "family": map[string]any{"opaque": []any{"synthetic", 0}}, "variant": []any{nil, "candidate"}},
		"fields":  map[string]any{"logical_object": b520ActuatorViewBytes, "unmodelled": []int{0, 255}},
	}
}

func b520ActuatorViewCall(b520ActuatorViewT *testing.T, b520ActuatorViewInput any) map[string]any {
	b520ActuatorViewT.Helper()
	b520ActuatorViewJSON, b520ActuatorViewErr := json.Marshal(b520ActuatorViewInput)
	if b520ActuatorViewErr != nil {
		b520ActuatorViewT.Fatal(b520ActuatorViewErr)
	}
	b520ActuatorViewResult, b520ActuatorViewErr := B520ActuatorViewInterpret(b520ActuatorViewJSON)
	if b520ActuatorViewErr != nil {
		b520ActuatorViewT.Fatal(b520ActuatorViewErr)
	}
	b520ActuatorViewOutput := map[string]any{}
	if b520ActuatorViewErr := json.Unmarshal(b520ActuatorViewResult, &b520ActuatorViewOutput); b520ActuatorViewErr != nil {
		b520ActuatorViewT.Fatal(b520ActuatorViewErr)
	}
	if b520ActuatorViewOutput["physical"] != "Unknown" || b520ActuatorViewOutput["qualified"] != false {
		b520ActuatorViewT.Fatalf("native qualification inferred: %s", b520ActuatorViewResult)
	}
	for b520ActuatorViewName, b520ActuatorViewValue := range b520ActuatorViewOutput["fields"].(map[string]any) {
		b520ActuatorViewResultField := b520ActuatorViewValue.(map[string]any)
		b520ActuatorViewEqual(b520ActuatorViewT, b520ActuatorViewResultField["grade"], "Hypothesis")
		if b520ActuatorViewName == "relay" || b520ActuatorViewName == "reserved" || b520ActuatorViewResultField["state"] == "incomplete" || b520ActuatorViewResultField["state"] == "unmapped" || b520ActuatorViewResultField["label"] != nil {
			if _, b520ActuatorViewExists := b520ActuatorViewResultField["number"]; b520ActuatorViewExists {
				b520ActuatorViewT.Fatalf("raw-only, labelled or incomplete field %s gained a number", b520ActuatorViewName)
			}
		}
	}
	b520ActuatorViewRoundTripJSON, b520ActuatorViewErr := json.Marshal(b520ActuatorViewOutput)
	if b520ActuatorViewErr != nil {
		b520ActuatorViewT.Fatal(b520ActuatorViewErr)
	}
	var b520ActuatorViewRoundTrip map[string]any
	if b520ActuatorViewErr := json.Unmarshal(b520ActuatorViewRoundTripJSON, &b520ActuatorViewRoundTrip); b520ActuatorViewErr != nil {
		b520ActuatorViewT.Fatal(b520ActuatorViewErr)
	}
	if !reflect.DeepEqual(b520ActuatorViewOutput, b520ActuatorViewRoundTrip) {
		b520ActuatorViewT.Fatal("JSON round trip changed state or presence")
	}
	return b520ActuatorViewOutput
}

func b520ActuatorViewField(b520ActuatorViewOutput map[string]any, b520ActuatorViewName string) map[string]any {
	return b520ActuatorViewOutput["fields"].(map[string]any)[b520ActuatorViewName].(map[string]any)
}

func b520ActuatorViewEqual(b520ActuatorViewT *testing.T, b520ActuatorViewGot, b520ActuatorViewWant any) {
	b520ActuatorViewT.Helper()
	b520ActuatorViewGotJSON, _ := json.Marshal(b520ActuatorViewGot)
	b520ActuatorViewWantJSON, _ := json.Marshal(b520ActuatorViewWant)
	if !bytes.Equal(b520ActuatorViewGotJSON, b520ActuatorViewWantJSON) {
		b520ActuatorViewT.Fatalf("got %s; want %s", b520ActuatorViewGotJSON, b520ActuatorViewWantJSON)
	}
}

func TestAb23Ao23Fixtures(b520ActuatorViewT *testing.T) {
	b520ActuatorViewJSON, b520ActuatorViewErr := os.ReadFile("testdata/ab23ao23_b520_actuator_views.json")
	if b520ActuatorViewErr != nil {
		b520ActuatorViewT.Fatal(b520ActuatorViewErr)
	}
	var b520ActuatorViewFixtures []map[string]json.RawMessage
	if b520ActuatorViewErr := json.Unmarshal(b520ActuatorViewJSON, &b520ActuatorViewFixtures); b520ActuatorViewErr != nil {
		b520ActuatorViewT.Fatal(b520ActuatorViewErr)
	}
	if len(b520ActuatorViewFixtures) != 2 {
		b520ActuatorViewT.Fatal("expected two independent vectors")
	}
	for _, b520ActuatorViewFixture := range b520ActuatorViewFixtures {
		var b520ActuatorViewID, b520ActuatorViewDirection string
		var b520ActuatorViewBytes []int
		var b520ActuatorViewNumbers map[string]float64
		for b520ActuatorViewKey, b520ActuatorViewTarget := range map[string]any{"id": &b520ActuatorViewID, "direction": &b520ActuatorViewDirection, "bytes": &b520ActuatorViewBytes, "numbers": &b520ActuatorViewNumbers} {
			if b520ActuatorViewErr := json.Unmarshal(b520ActuatorViewFixture[b520ActuatorViewKey], b520ActuatorViewTarget); b520ActuatorViewErr != nil {
				b520ActuatorViewT.Fatal(b520ActuatorViewErr)
			}
		}
		b520ActuatorViewT.Run(b520ActuatorViewID, func(b520ActuatorViewT *testing.T) {
			b520ActuatorViewInput := b520ActuatorViewInput(b520ActuatorViewDirection, b520ActuatorViewBytes)
			b520ActuatorViewOutput := b520ActuatorViewCall(b520ActuatorViewT, b520ActuatorViewInput)
			b520ActuatorViewEqual(b520ActuatorViewT, b520ActuatorViewOutput["state"], "Hypothesis")
			b520ActuatorViewEqual(b520ActuatorViewT, b520ActuatorViewOutput["raw"], b520ActuatorViewInput["fields"])
			b520ActuatorViewEqual(b520ActuatorViewT, b520ActuatorViewOutput["context"], b520ActuatorViewInput["context"])
			for b520ActuatorViewName, b520ActuatorViewNumber := range b520ActuatorViewNumbers {
				b520ActuatorViewField := b520ActuatorViewField(b520ActuatorViewOutput, b520ActuatorViewName)
				b520ActuatorViewEqual(b520ActuatorViewT, b520ActuatorViewField["number"], b520ActuatorViewNumber)
				b520ActuatorViewEqual(b520ActuatorViewT, b520ActuatorViewField["grade"], "Hypothesis")
				b520ActuatorViewEqual(b520ActuatorViewT, b520ActuatorViewField["state"], "Hypothesis")
			}
		})
	}
}

func TestAb23Ao23RequestMappings(b520ActuatorViewT *testing.T) {
	for _, b520ActuatorViewStepper := range []int{0, 240, 241, 254, 255} {
		b520ActuatorViewT.Run(fmt.Sprintf("ab23ao23_stepper_%d", b520ActuatorViewStepper), func(b520ActuatorViewT *testing.T) {
			b520ActuatorViewOutput := b520ActuatorViewCall(b520ActuatorViewT, b520ActuatorViewInput("request", []int{83, 12, 251, 37, 255, 255, b520ActuatorViewStepper, 0, 255, 0}))
			b520ActuatorViewField := b520ActuatorViewField(b520ActuatorViewOutput, "stepper")
			if b520ActuatorViewStepper <= 240 {
				b520ActuatorViewEqual(b520ActuatorViewT, b520ActuatorViewField["number"], b520ActuatorViewStepper)
			} else {
				if _, b520ActuatorViewExists := b520ActuatorViewField["number"]; b520ActuatorViewExists {
					b520ActuatorViewT.Fatal("unmapped/calibration carrier gained a number")
				}
				if b520ActuatorViewStepper == 255 {
					b520ActuatorViewEqual(b520ActuatorViewT, b520ActuatorViewField["label"], "calibrate/go-to-zero")
				} else {
					b520ActuatorViewEqual(b520ActuatorViewT, b520ActuatorViewField["state"], "unmapped")
					b520ActuatorViewEqual(b520ActuatorViewT, b520ActuatorViewOutput["state"], "partial")
				}
			}
			for _, b520ActuatorViewName := range []string{"PWM1", "PWM2", "ED"} {
				b520ActuatorViewEqual(b520ActuatorViewT, b520ActuatorViewFieldFrom(b520ActuatorViewOutput, b520ActuatorViewName, "number"), 127.5)
				b520ActuatorViewEqual(b520ActuatorViewT, b520ActuatorViewFieldFrom(b520ActuatorViewOutput, b520ActuatorViewName, "calculation"), "Proven")
			}
		})
	}
	for b520ActuatorViewIndex, b520ActuatorViewLabels := range map[int][]string{7: {"disabled", "enabled"}, 9: {"standard3K3", "extended1K", "low-temperature15K", "standard1K"}} {
		for b520ActuatorViewValue := 0; b520ActuatorViewValue <= len(b520ActuatorViewLabels); b520ActuatorViewValue++ {
			b520ActuatorViewT.Run(fmt.Sprintf("ab23ao23_map_%d_%d", b520ActuatorViewIndex, b520ActuatorViewValue), func(b520ActuatorViewT *testing.T) {
				b520ActuatorViewBytes := []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
				b520ActuatorViewBytes[b520ActuatorViewIndex] = b520ActuatorViewValue
				b520ActuatorViewName := "LED"
				if b520ActuatorViewIndex == 9 {
					b520ActuatorViewName = "temperature-table"
				}
				b520ActuatorViewOutput := b520ActuatorViewCall(b520ActuatorViewT, b520ActuatorViewInput("request", b520ActuatorViewBytes))
				b520ActuatorViewField := b520ActuatorViewField(b520ActuatorViewOutput, b520ActuatorViewName)
				if b520ActuatorViewValue == len(b520ActuatorViewLabels) {
					b520ActuatorViewEqual(b520ActuatorViewT, b520ActuatorViewField["state"], "unmapped")
					if _, b520ActuatorViewExists := b520ActuatorViewField["number"]; b520ActuatorViewExists {
						b520ActuatorViewT.Fatal("unmapped number")
					}
				} else {
					b520ActuatorViewEqual(b520ActuatorViewT, b520ActuatorViewField["label"], b520ActuatorViewLabels[b520ActuatorViewValue])
					if b520ActuatorViewIndex == 9 {
						b520ActuatorViewEqual(b520ActuatorViewT, b520ActuatorViewField["flags"].(map[string]any)["candidate_celsius_range"], [][2]float64{{-13.5, 116}, {-20, 155}, {-39, 80}, {-13.9, 130}}[b520ActuatorViewValue])
					}
				}
			})
		}
	}
}

func b520ActuatorViewFieldFrom(b520ActuatorViewOutput map[string]any, b520ActuatorViewName, b520ActuatorViewKey string) any {
	return b520ActuatorViewField(b520ActuatorViewOutput, b520ActuatorViewName)[b520ActuatorViewKey]
}

func TestAb23Ao23Domains(b520ActuatorViewT *testing.T) {
	for b520ActuatorViewIndex, b520ActuatorViewName := range map[int]string{3: "ADC", 5: "temperature1", 7: "temperature2", 9: "temperature3", 11: "temperature4", 14: "impeller", 16: "vortex"} {
		b520ActuatorViewValues := []int{0, 1023, 1024}
		b520ActuatorViewLimit := 1023
		if b520ActuatorViewIndex >= 5 && b520ActuatorViewIndex <= 11 {
			b520ActuatorViewValues = []int{0, 16, 3072, 3200, 3201, 65535}
			b520ActuatorViewLimit = 3200
		} else if b520ActuatorViewIndex == 14 {
			b520ActuatorViewValues = []int{0, 4, 1000, 1001}
			b520ActuatorViewLimit = 1000
		} else if b520ActuatorViewIndex == 16 {
			b520ActuatorViewValues = []int{0, 4, 5, 350, 351}
			b520ActuatorViewLimit = 350
		}
		for _, b520ActuatorViewValue := range b520ActuatorViewValues {
			b520ActuatorViewT.Run(fmt.Sprintf("ab23ao23_%s_%d", b520ActuatorViewName, b520ActuatorViewValue), func(b520ActuatorViewT *testing.T) {
				b520ActuatorViewBytes := make([]int, 18)
				b520ActuatorViewBytes[b520ActuatorViewIndex], b520ActuatorViewBytes[b520ActuatorViewIndex+1] = b520ActuatorViewValue%256, b520ActuatorViewValue/256
				b520ActuatorViewOutput := b520ActuatorViewCall(b520ActuatorViewT, b520ActuatorViewInput("positive", b520ActuatorViewBytes))
				b520ActuatorViewField := b520ActuatorViewField(b520ActuatorViewOutput, b520ActuatorViewName)
				b520ActuatorViewWant := float64(b520ActuatorViewValue)
				if b520ActuatorViewIndex >= 5 && b520ActuatorViewIndex <= 11 {
					b520ActuatorViewWant = map[int]float64{0: 0, 16: 1, 3072: 192, 3200: 200, 3201: 200.0625, 65535: -0.0625}[b520ActuatorViewValue]
					b520ActuatorViewEqual(b520ActuatorViewT, b520ActuatorViewField["calculation"], "Proven")
				}
				b520ActuatorViewEqual(b520ActuatorViewT, b520ActuatorViewField["number"], b520ActuatorViewWant)
				b520ActuatorViewEqual(b520ActuatorViewT, b520ActuatorViewField["candidate_domain"], b520ActuatorViewValue <= b520ActuatorViewLimit)
				b520ActuatorViewEqual(b520ActuatorViewT, b520ActuatorViewOutput["state"], "Hypothesis")
				if b520ActuatorViewIndex == 16 {
					b520ActuatorViewEqual(b520ActuatorViewT, b520ActuatorViewField["flags"], map[string]any{"candidate_operational_domain": b520ActuatorViewValue >= 5 && b520ActuatorViewValue <= 350, "candidate_hysteresis_hz": map[string]int{"on": 5, "off": 4}})
				}
			})
		}
	}
}

func TestAb23Ao23IndependentBytes(b520ActuatorViewT *testing.T) {
	for _, b520ActuatorViewDirection := range []string{"request", "positive"} {
		b520ActuatorViewNames := []string{"relay", "PWM1", "PWM2", "stepper", "LED", "ED", "temperature-table"}
		b520ActuatorViewOffsets, b520ActuatorViewWidths, b520ActuatorViewLength := []int{3, 4, 5, 6, 7, 8, 9}, []int{1, 1, 1, 1, 1, 1, 1}, 10
		if b520ActuatorViewDirection == "positive" {
			b520ActuatorViewNames = []string{"ADC", "temperature1", "temperature2", "temperature3", "temperature4", "reserved", "impeller", "vortex"}
			b520ActuatorViewOffsets, b520ActuatorViewWidths, b520ActuatorViewLength = []int{3, 5, 7, 9, 11, 13, 14, 16}, []int{2, 2, 2, 2, 2, 1, 2, 2}, 18
		}
		b520ActuatorViewControl := b520ActuatorViewCall(b520ActuatorViewT, b520ActuatorViewInput(b520ActuatorViewDirection, make([]int, b520ActuatorViewLength)))
		for b520ActuatorViewPosition := 0; b520ActuatorViewPosition < b520ActuatorViewLength+2; b520ActuatorViewPosition++ {
			b520ActuatorViewT.Run(fmt.Sprintf("ab23ao23_%s_byte_%d", b520ActuatorViewDirection, b520ActuatorViewPosition), func(b520ActuatorViewT *testing.T) {
				b520ActuatorViewBytes := make([]int, b520ActuatorViewLength+2)
				b520ActuatorViewBytes[b520ActuatorViewPosition] = 1
				b520ActuatorViewInput := b520ActuatorViewInput(b520ActuatorViewDirection, b520ActuatorViewBytes)
				b520ActuatorViewOutput := b520ActuatorViewCall(b520ActuatorViewT, b520ActuatorViewInput)
				b520ActuatorViewEqual(b520ActuatorViewT, b520ActuatorViewOutput["raw"], b520ActuatorViewInput["fields"])
				for b520ActuatorViewIndex, b520ActuatorViewName := range b520ActuatorViewNames {
					b520ActuatorViewField := b520ActuatorViewField(b520ActuatorViewOutput, b520ActuatorViewName)
					if b520ActuatorViewPosition < b520ActuatorViewOffsets[b520ActuatorViewIndex] || b520ActuatorViewPosition >= b520ActuatorViewOffsets[b520ActuatorViewIndex]+b520ActuatorViewWidths[b520ActuatorViewIndex] {
						b520ActuatorViewEqual(b520ActuatorViewT, b520ActuatorViewField, b520ActuatorViewFieldFromControl(b520ActuatorViewControl, b520ActuatorViewName))
					} else {
						b520ActuatorViewEqual(b520ActuatorViewT, b520ActuatorViewField["raw"], b520ActuatorViewBytes[b520ActuatorViewOffsets[b520ActuatorViewIndex]:b520ActuatorViewOffsets[b520ActuatorViewIndex]+b520ActuatorViewWidths[b520ActuatorViewIndex]])
						if b520ActuatorViewName == "reserved" {
							b520ActuatorViewEqual(b520ActuatorViewT, b520ActuatorViewField["reserved"], true)
						} else if b520ActuatorViewDirection == "positive" {
							b520ActuatorViewNumber := float64(1)
							if b520ActuatorViewPosition != b520ActuatorViewOffsets[b520ActuatorViewIndex] {
								b520ActuatorViewNumber = 256
							}
							if b520ActuatorViewIndex >= 1 && b520ActuatorViewIndex <= 4 {
								b520ActuatorViewNumber /= 16
							}
							b520ActuatorViewEqual(b520ActuatorViewT, b520ActuatorViewField["number"], b520ActuatorViewNumber)
						} else if b520ActuatorViewName == "PWM1" || b520ActuatorViewName == "PWM2" || b520ActuatorViewName == "ED" {
							b520ActuatorViewEqual(b520ActuatorViewT, b520ActuatorViewField["number"], 0.5)
						}
					}
				}
			})
		}
	}
}

func b520ActuatorViewFieldFromControl(b520ActuatorViewOutput map[string]any, b520ActuatorViewName string) map[string]any {
	return b520ActuatorViewField(b520ActuatorViewOutput, b520ActuatorViewName)
}

func TestAb23Ao23Truncation(b520ActuatorViewT *testing.T) {
	for _, b520ActuatorViewDirection := range []string{"request", "positive"} {
		b520ActuatorViewNames, b520ActuatorViewEnds := []string{"relay", "PWM1", "PWM2", "stepper", "LED", "ED", "temperature-table"}, []int{4, 5, 6, 7, 8, 9, 10}
		if b520ActuatorViewDirection == "positive" {
			b520ActuatorViewNames, b520ActuatorViewEnds = []string{"ADC", "temperature1", "temperature2", "temperature3", "temperature4", "reserved", "impeller", "vortex"}, []int{5, 7, 9, 11, 13, 14, 16, 18}
		}
		b520ActuatorViewFull := make([]int, b520ActuatorViewEnds[len(b520ActuatorViewEnds)-1])
		b520ActuatorViewControl := b520ActuatorViewCall(b520ActuatorViewT, b520ActuatorViewInput(b520ActuatorViewDirection, b520ActuatorViewFull))
		for b520ActuatorViewLength := 0; b520ActuatorViewLength <= len(b520ActuatorViewFull); b520ActuatorViewLength++ {
			b520ActuatorViewT.Run(fmt.Sprintf("ab23ao23_%s_length_%d", b520ActuatorViewDirection, b520ActuatorViewLength), func(b520ActuatorViewT *testing.T) {
				b520ActuatorViewOutput := b520ActuatorViewCall(b520ActuatorViewT, b520ActuatorViewInput(b520ActuatorViewDirection, b520ActuatorViewFull[:b520ActuatorViewLength]))
				b520ActuatorViewWantState := "partial"
				if b520ActuatorViewLength == len(b520ActuatorViewFull) {
					b520ActuatorViewWantState = "Hypothesis"
				}
				b520ActuatorViewEqual(b520ActuatorViewT, b520ActuatorViewOutput["state"], b520ActuatorViewWantState)
				for b520ActuatorViewIndex, b520ActuatorViewName := range b520ActuatorViewNames {
					b520ActuatorViewField := b520ActuatorViewField(b520ActuatorViewOutput, b520ActuatorViewName)
					if b520ActuatorViewLength >= b520ActuatorViewEnds[b520ActuatorViewIndex] {
						b520ActuatorViewEqual(b520ActuatorViewT, b520ActuatorViewField, b520ActuatorViewFieldFromControl(b520ActuatorViewControl, b520ActuatorViewName))
					} else {
						b520ActuatorViewEqual(b520ActuatorViewT, b520ActuatorViewField["state"], "incomplete")
						if _, b520ActuatorViewExists := b520ActuatorViewField["number"]; b520ActuatorViewExists {
							b520ActuatorViewT.Fatal("default filled incomplete field")
						}
					}
				}
			})
		}
	}
}

func TestAb23Ao23Selectors(b520ActuatorViewT *testing.T) {
	for _, b520ActuatorViewDirection := range []string{"request", "positive"} {
		for b520ActuatorViewKey, b520ActuatorViewValues := range map[string][]any{"primary": {180, nil}, "secondary": {33, nil}, "operation": {1, nil}, "direction": {"response", "negative", "unknown", nil}, "view": {"unknown", "", nil}} {
			for b520ActuatorViewIndex, b520ActuatorViewValue := range b520ActuatorViewValues {
				b520ActuatorViewT.Run(fmt.Sprintf("ab23ao23_%s_%s_%d", b520ActuatorViewDirection, b520ActuatorViewKey, b520ActuatorViewIndex), func(b520ActuatorViewT *testing.T) {
					b520ActuatorViewInput := b520ActuatorViewInput(b520ActuatorViewDirection, make([]int, 20))
					b520ActuatorViewControl := b520ActuatorViewCall(b520ActuatorViewT, b520ActuatorViewInput)
					b520ActuatorViewEqual(b520ActuatorViewT, b520ActuatorViewControl["state"], "Hypothesis")
					b520ActuatorViewContext := b520ActuatorViewInput["context"].(map[string]any)
					if b520ActuatorViewValue == nil {
						delete(b520ActuatorViewContext, b520ActuatorViewKey)
					} else {
						b520ActuatorViewContext[b520ActuatorViewKey] = b520ActuatorViewValue
					}
					b520ActuatorViewOutput := b520ActuatorViewCall(b520ActuatorViewT, b520ActuatorViewInput)
					b520ActuatorViewEqual(b520ActuatorViewT, b520ActuatorViewOutput["state"], "Unknown")
					b520ActuatorViewEqual(b520ActuatorViewT, b520ActuatorViewOutput["fields"], map[string]any{})
					b520ActuatorViewEqual(b520ActuatorViewT, b520ActuatorViewOutput["raw"], b520ActuatorViewInput["fields"])
					b520ActuatorViewEqual(b520ActuatorViewT, b520ActuatorViewOutput["context"], b520ActuatorViewContext)
				})
			}
		}
	}
	for _, b520ActuatorViewKey := range []string{"family", "variant", "extra"} {
		for _, b520ActuatorViewValue := range []any{nil, 0, "unknown", map[string]any{"opaque": []any{true, 3}}} {
			b520ActuatorViewInput := b520ActuatorViewInput("positive", make([]int, 18))
			b520ActuatorViewInput["context"].(map[string]any)[b520ActuatorViewKey] = b520ActuatorViewValue
			b520ActuatorViewOutput := b520ActuatorViewCall(b520ActuatorViewT, b520ActuatorViewInput)
			b520ActuatorViewEqual(b520ActuatorViewT, b520ActuatorViewOutput["context"], b520ActuatorViewInput["context"])
			b520ActuatorViewEqual(b520ActuatorViewT, b520ActuatorViewOutput["state"], "Hypothesis")
		}
	}
}

func TestAb23Ao23Presence(b520ActuatorViewT *testing.T) {
	for _, b520ActuatorViewDirection := range []string{"request", "positive"} {
		for _, b520ActuatorViewKind := range []string{"missing", "empty", "zero"} {
			b520ActuatorViewT.Run("ab23ao23_"+b520ActuatorViewDirection+"_"+b520ActuatorViewKind, func(b520ActuatorViewT *testing.T) {
				b520ActuatorViewInput := b520ActuatorViewInput(b520ActuatorViewDirection, []int{})
				b520ActuatorViewFields := b520ActuatorViewInput["fields"].(map[string]any)
				switch b520ActuatorViewKind {
				case "missing":
					delete(b520ActuatorViewFields, "logical_object")
				case "zero":
					b520ActuatorViewFields["logical_object"] = []int{0, 0, 0, 0, 0}
				}
				b520ActuatorViewInput["arguments"] = map[string]any{}
				b520ActuatorViewOutput := b520ActuatorViewCall(b520ActuatorViewT, b520ActuatorViewInput)
				b520ActuatorViewEqual(b520ActuatorViewT, b520ActuatorViewOutput["raw"], b520ActuatorViewFields)
				b520ActuatorViewEqual(b520ActuatorViewT, b520ActuatorViewOutput["arguments"], map[string]any{})
				b520ActuatorViewName := "relay"
				if b520ActuatorViewDirection == "positive" {
					b520ActuatorViewName = "ADC"
				}
				b520ActuatorViewField := b520ActuatorViewField(b520ActuatorViewOutput, b520ActuatorViewName)
				if b520ActuatorViewKind == "missing" {
					b520ActuatorViewEqual(b520ActuatorViewT, b520ActuatorViewField["raw"], nil)
				} else if b520ActuatorViewKind == "empty" {
					b520ActuatorViewEqual(b520ActuatorViewT, b520ActuatorViewField["raw"], []int{})
				} else if b520ActuatorViewDirection == "positive" {
					b520ActuatorViewEqual(b520ActuatorViewT, b520ActuatorViewField["number"], 0)
				}
				delete(b520ActuatorViewInput, "arguments")
				if _, b520ActuatorViewExists := b520ActuatorViewCall(b520ActuatorViewT, b520ActuatorViewInput)["arguments"]; b520ActuatorViewExists {
					b520ActuatorViewT.Fatal("absent arguments gained presence")
				}
			})
		}
	}
}

func TestAb23Ao23OpaqueFields(b520ActuatorViewT *testing.T) {
	for _, b520ActuatorViewDirection := range []string{"request", "positive"} {
		b520ActuatorViewControl := b520ActuatorViewCall(b520ActuatorViewT, b520ActuatorViewInput(b520ActuatorViewDirection, make([]int, 18)))
		for b520ActuatorViewIndex, b520ActuatorViewValue := range []any{nil, []int{}, []int{0}, []int{255}, []int{1, 2, 3}} {
			b520ActuatorViewT.Run(fmt.Sprintf("ab23ao23_%s_opaque_%d", b520ActuatorViewDirection, b520ActuatorViewIndex), func(b520ActuatorViewT *testing.T) {
				b520ActuatorViewInput := b520ActuatorViewInput(b520ActuatorViewDirection, make([]int, 18))
				b520ActuatorViewFields := b520ActuatorViewInput["fields"].(map[string]any)
				if b520ActuatorViewValue == nil {
					delete(b520ActuatorViewFields, "unmodelled")
				} else {
					b520ActuatorViewFields["unmodelled"] = b520ActuatorViewValue
				}
				b520ActuatorViewOutput := b520ActuatorViewCall(b520ActuatorViewT, b520ActuatorViewInput)
				b520ActuatorViewEqual(b520ActuatorViewT, b520ActuatorViewOutput["raw"], b520ActuatorViewFields)
				b520ActuatorViewEqual(b520ActuatorViewT, b520ActuatorViewOutput["fields"], b520ActuatorViewControl["fields"])
			})
		}
	}
	for _, b520ActuatorViewKey := range []string{"family", "variant"} {
		b520ActuatorViewInput := b520ActuatorViewInput("request", make([]int, 10))
		delete(b520ActuatorViewInput["context"].(map[string]any), b520ActuatorViewKey)
		b520ActuatorViewOutput := b520ActuatorViewCall(b520ActuatorViewT, b520ActuatorViewInput)
		b520ActuatorViewEqual(b520ActuatorViewT, b520ActuatorViewOutput["context"], b520ActuatorViewInput["context"])
		b520ActuatorViewEqual(b520ActuatorViewT, b520ActuatorViewOutput["state"], "Hypothesis")
	}
	b520ActuatorViewOutput := b520ActuatorViewCall(b520ActuatorViewT, map[string]any{"context": map[string]any{"primary": 181, "secondary": 32, "operation": json.Number("18446744073709551616"), "direction": "request", "view": "logical-object"}, "fields": map[string]any{}})
	b520ActuatorViewEqual(b520ActuatorViewT, b520ActuatorViewOutput["state"], "Unknown")
}

func TestAb23Ao23InvalidJSON(b520ActuatorViewT *testing.T) {
	for b520ActuatorViewIndex, b520ActuatorViewJSON := range []string{
		``, `{`, `null`, `[]`, `{} {}`, `{"context":{},"fields":{},"context":{}}`,
		`{"context":{"family":{"x":1,"x":2}},"fields":{}}`,
		`{"context":{},"fields":{"x":[],"x":[]}}`,
		`{"context":{},"fields":{},"arguments":{"x":0,"x":1}}`,
		`{"context":null,"fields":{}}`, `{"context":{},"fields":null}`,
		`{"context":{},"fields":{"x":null}}`, `{"context":{},"fields":{"x":"AA=="}}`,
		`{"context":{},"fields":{"x":[-1]}}`, `{"context":{},"fields":{"x":[256]}}`,
		`{"context":{},"fields":{"x":[1.5]}}`, `{"context":{},"fields":{"x":[1.0]}}`,
		`{"context":{},"fields":{"x":[true]}}`, `{"context":{},"fields":{"x":[null]}}`,
		`{"context":{},"fields":{"x":["1"]}}`, `{"context":{},"fields":{"x":[{}]}}`,
		`{"context":{},"fields":{"x":[[]]}}`, `{"context":{},"fields":{"x":[1e0]}}`,
		`{"context":{},"fields":{},"arguments":null}`, `{"context":{},"fields":{},"arguments":[]}`,
		`{"context":{},"fields":{},"arguments":{"positive":true}}`,
		`{"context":{"operation":"0"},"fields":{}}`, `{"context":{"primary":181.1},"fields":{}}`,
		`{"context":{"secondary":null},"fields":{}}`, `{"context":{"direction":0},"fields":{}}`,
		`{"context":{"view":false},"fields":{}}`, `{"context":{},"fields":{},"extra":0}`,
		"{\"context\":{\"family\":\"" + string([]byte{255}) + "\"},\"fields\":{}}",
	} {
		b520ActuatorViewT.Run(fmt.Sprintf("ab23ao23_invalid_%d", b520ActuatorViewIndex), func(b520ActuatorViewT *testing.T) {
			b520ActuatorViewResult, b520ActuatorViewErr := B520ActuatorViewInterpret(json.RawMessage(b520ActuatorViewJSON))
			if b520ActuatorViewErr == nil || b520ActuatorViewResult != nil {
				b520ActuatorViewT.Fatalf("invalid input succeeded: %s => %s, %v", b520ActuatorViewJSON, b520ActuatorViewResult, b520ActuatorViewErr)
			}
		})
	}
}

func TestAb23Ao23UnicodeEscapes(t *testing.T) {
	for _, input := range []string{
		`{"context":{"family":"\ud800"},"fields":{}}`,
		`{"context":{"family":"\udfff"},"fields":{}}`,
		`{"context":{"family":"\ud800\u0041"},"fields":{}}`,
		`{"context":{"family":"\ud800\ud800"},"fields":{}}`,
		`{"context":{"family":"\udc00\ud800"},"fields":{}}`,
		`{"context":{"family":"\ud800x\udc00"},"fields":{}}`,
		`{"context":{"family":"\ud800\\udc00"},"fields":{}}`,
		`{"context":{"family":{"\ud800":[]}},"fields":{}}`,
		`{"context":{"family":["\udc00"]},"fields":{}}`,
		`{"context":{},"fields":{"\ud800":[]}}`,
	} {
		t.Run(input, func(t *testing.T) {
			output, err := B520ActuatorViewInterpret(json.RawMessage(input))
			if err == nil || output != nil {
				t.Fatalf("unpaired surrogate accepted: %s, %v", output, err)
			}
		})
	}
	for _, tc := range []struct{ encoded, want string }{
		{`"\ud800\udc00"`, "\U00010000"},
		{`"\uDBFF\uDFFF"`, "\U0010ffff"},
		{`"\ud83d\ude00"`, "\U0001f600"},
		{`"\\ud800"`, `\ud800`},
		{`"\\\ud800\udc00"`, "\\\U00010000"},
		{`"\ufffd"`, "\ufffd"},
		{`"�"`, "\ufffd"},
		{`"\"\u0041\\"`, "\"A\\"},
	} {
		t.Run(tc.encoded, func(t *testing.T) {
			input := json.RawMessage(`{"context":{"family":` + tc.encoded + `},"fields":{}}`)
			output, err := B520ActuatorViewInterpret(input)
			if err != nil {
				t.Fatal(err)
			}
			var decoded struct{ Context struct{ Family string } }
			if err := json.Unmarshal(output, &decoded); err != nil {
				t.Fatal(err)
			}
			if decoded.Context.Family != tc.want {
				t.Fatalf("got %q; want %q", decoded.Context.Family, tc.want)
			}
		})
	}
}

func TestAb23Ao23Isolation(b520ActuatorViewT *testing.T) {
	b520ActuatorViewInput := b520ActuatorViewInput("positive", []int{0, 0, 0, 255, 3, 0, 12, 255, 255, 0, 0, 16, 0, 171, 232, 3, 94, 1})
	b520ActuatorViewInput["arguments"] = map[string]any{}
	b520ActuatorViewJSON, _ := json.Marshal(b520ActuatorViewInput)
	b520ActuatorViewOriginal := append(json.RawMessage(nil), b520ActuatorViewJSON...)
	b520ActuatorViewFirst, b520ActuatorViewErr := B520ActuatorViewInterpret(b520ActuatorViewJSON)
	if b520ActuatorViewErr != nil {
		b520ActuatorViewT.Fatal(b520ActuatorViewErr)
	}
	b520ActuatorViewSecond, b520ActuatorViewErr := B520ActuatorViewInterpret(b520ActuatorViewJSON)
	if b520ActuatorViewErr != nil {
		b520ActuatorViewT.Fatal(b520ActuatorViewErr)
	}
	b520ActuatorViewSnapshot := append(json.RawMessage(nil), b520ActuatorViewFirst...)
	for b520ActuatorViewIndex := range b520ActuatorViewJSON {
		b520ActuatorViewJSON[b520ActuatorViewIndex] = 'x'
	}
	if !bytes.Equal(b520ActuatorViewFirst, b520ActuatorViewSnapshot) {
		b520ActuatorViewT.Fatal("input mutation affected output")
	}
	var b520ActuatorViewNested map[string]any
	if b520ActuatorViewErr := json.Unmarshal(b520ActuatorViewFirst, &b520ActuatorViewNested); b520ActuatorViewErr != nil {
		b520ActuatorViewT.Fatal(b520ActuatorViewErr)
	}
	b520ActuatorViewNested["raw"].(map[string]any)["logical_object"].([]any)[5] = 99
	b520ActuatorViewEqual(b520ActuatorViewT, b520ActuatorViewFieldFrom(b520ActuatorViewNested, "temperature1", "raw"), []int{0, 12})
	b520ActuatorViewField(b520ActuatorViewNested, "temperature1")["raw"].([]any)[0] = 88
	b520ActuatorViewEqual(b520ActuatorViewT, b520ActuatorViewFieldFrom(b520ActuatorViewNested, "temperature2", "raw"), []int{255, 255})
	b520ActuatorViewNested["context"].(map[string]any)["family"].(map[string]any)["opaque"].([]any)[0] = "changed"
	b520ActuatorViewNested["arguments"].(map[string]any)["changed"] = true
	for b520ActuatorViewIndex := range b520ActuatorViewFirst {
		b520ActuatorViewFirst[b520ActuatorViewIndex] = 'y'
	}
	if !bytes.Equal(b520ActuatorViewSecond, b520ActuatorViewSnapshot) {
		b520ActuatorViewT.Fatal("output mutation affected another result")
	}
	b520ActuatorViewThird, b520ActuatorViewErr := B520ActuatorViewInterpret(b520ActuatorViewOriginal)
	if b520ActuatorViewErr != nil || !reflect.DeepEqual(b520ActuatorViewThird, b520ActuatorViewSecond) {
		b520ActuatorViewT.Fatal("mutation affected later call", b520ActuatorViewErr)
	}
}
