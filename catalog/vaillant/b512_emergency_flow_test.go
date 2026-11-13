package vaillant

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"
)

type b512EmergencyFlowCase struct {
	B512EmergencyFlowID      string             `json:"id"`
	B512EmergencyFlowInput   map[string]any     `json:"input"`
	B512EmergencyFlowNumbers map[string]float64 `json:"numbers"`
}

func b512EmergencyFlowCases(b512EmergencyFlowT *testing.T) []b512EmergencyFlowCase {
	b512EmergencyFlowT.Helper()
	b512EmergencyFlowData, b512EmergencyFlowErr := os.ReadFile("testdata/ab15ap15_b512_emergency_flow.json")
	if b512EmergencyFlowErr != nil {
		b512EmergencyFlowT.Fatal(b512EmergencyFlowErr)
	}
	var b512EmergencyFlowCases []b512EmergencyFlowCase
	if b512EmergencyFlowErr := json.Unmarshal(b512EmergencyFlowData, &b512EmergencyFlowCases); b512EmergencyFlowErr != nil {
		b512EmergencyFlowT.Fatal(b512EmergencyFlowErr)
	}
	if len(b512EmergencyFlowCases) != 4 {
		b512EmergencyFlowT.Fatal("fixture count")
	}
	return b512EmergencyFlowCases
}

func b512EmergencyFlowEvaluate(b512EmergencyFlowT *testing.T, b512EmergencyFlowInput any) map[string]any {
	b512EmergencyFlowT.Helper()
	b512EmergencyFlowData, b512EmergencyFlowErr := json.Marshal(b512EmergencyFlowInput)
	if b512EmergencyFlowErr != nil {
		b512EmergencyFlowT.Fatal(b512EmergencyFlowErr)
	}
	b512EmergencyFlowOutput, b512EmergencyFlowErr := B512EmergencyFlowEvaluate(b512EmergencyFlowData)
	if b512EmergencyFlowErr != nil {
		b512EmergencyFlowT.Fatal(b512EmergencyFlowErr)
	}
	var b512EmergencyFlowResult map[string]any
	if b512EmergencyFlowErr := json.Unmarshal(b512EmergencyFlowOutput, &b512EmergencyFlowResult); b512EmergencyFlowErr != nil {
		b512EmergencyFlowT.Fatal(b512EmergencyFlowErr)
	}
	return b512EmergencyFlowResult
}

func b512EmergencyFlowField(b512EmergencyFlowOutput map[string]any, b512EmergencyFlowName string) map[string]any {
	return b512EmergencyFlowOutput["fields"].(map[string]any)[b512EmergencyFlowName].(map[string]any)
}

func TestAb15Ap15Synthetic(b512EmergencyFlowT *testing.T) {
	for _, b512EmergencyFlowCase := range b512EmergencyFlowCases(b512EmergencyFlowT) {
		b512EmergencyFlowT.Run(b512EmergencyFlowCase.B512EmergencyFlowID, func(b512EmergencyFlowT *testing.T) {
			b512EmergencyFlowOutput := b512EmergencyFlowEvaluate(b512EmergencyFlowT, b512EmergencyFlowCase.B512EmergencyFlowInput)
			if b512EmergencyFlowOutput["state"] != "complete" || b512EmergencyFlowOutput["physical"] != "Unknown" || b512EmergencyFlowOutput["qualified"] != false {
				b512EmergencyFlowT.Fatal(b512EmergencyFlowOutput)
			}
			if !reflect.DeepEqual(b512EmergencyFlowOutput["context"], b512EmergencyFlowCase.B512EmergencyFlowInput["context"]) || !reflect.DeepEqual(b512EmergencyFlowOutput["raw"], b512EmergencyFlowCase.B512EmergencyFlowInput["fields"]) {
				b512EmergencyFlowT.Fatal("preservation")
			}
			for b512EmergencyFlowName, b512EmergencyFlowNumber := range b512EmergencyFlowCase.B512EmergencyFlowNumbers {
				b512EmergencyFlowField := b512EmergencyFlowField(b512EmergencyFlowOutput, b512EmergencyFlowName)
				if b512EmergencyFlowField["number"] != b512EmergencyFlowNumber || b512EmergencyFlowField["state"] != "Hypothesis" || b512EmergencyFlowField["calculation"] != "Proven" {
					b512EmergencyFlowT.Fatal(b512EmergencyFlowName, b512EmergencyFlowField)
				}
			}
			b512EmergencyFlowContext := b512EmergencyFlowCase.B512EmergencyFlowInput["context"].(map[string]any)
			b512EmergencyFlowFlags := b512EmergencyFlowField(b512EmergencyFlowOutput, "flags")
			if b512EmergencyFlowContext["direction"] == "positive" {
				b512EmergencyFlowCandidate := b512EmergencyFlowOutput["candidate"].(map[string]any)
				if !reflect.DeepEqual(b512EmergencyFlowCandidate["declared_positive_prefix"], []any{float64(181), float64(18), float64(12)}) || b512EmergencyFlowCandidate["literal_reply_conflict"] != true || b512EmergencyFlowCandidate["physical_association"] != "Unknown" || b512EmergencyFlowCandidate["requested_operation"] != b512EmergencyFlowContext["operation"] {
					b512EmergencyFlowT.Fatal(b512EmergencyFlowCandidate)
				}
			}
			if b512EmergencyFlowContext["view"] == "flow-io" {
				if b512EmergencyFlowFlags["flags"].(map[string]any)["compressor-off"] != true {
					b512EmergencyFlowT.Fatal(b512EmergencyFlowFlags)
				}
				if b512EmergencyFlowContext["direction"] == "positive" {
					if b512EmergencyFlowFlags["reserved"] != float64(128) || b512EmergencyFlowFlags["unmodeled"] != nil || b512EmergencyFlowField(b512EmergencyFlowOutput, "setpoint")["candidate_domain"] != false {
						b512EmergencyFlowT.Fatal(b512EmergencyFlowOutput)
					}
				} else if b512EmergencyFlowFlags["unmodeled"] != float64(128) || b512EmergencyFlowFlags["reserved"] != float64(0) {
					b512EmergencyFlowT.Fatal(b512EmergencyFlowFlags)
				}
			} else if b512EmergencyFlowFlags["unmodeled"] != float64(128) || b512EmergencyFlowFlags["reserved"] != float64(0) {
				b512EmergencyFlowT.Fatal(b512EmergencyFlowFlags)
			}
		})
	}
}

func TestAb15Ap15Fields(b512EmergencyFlowT *testing.T) {
	for b512EmergencyFlowIndex := 0; b512EmergencyFlowIndex < 4; b512EmergencyFlowIndex++ {
		for b512EmergencyFlowName := range b512EmergencyFlowCases(b512EmergencyFlowT)[b512EmergencyFlowIndex].B512EmergencyFlowInput["fields"].(map[string]any) {
			for _, b512EmergencyFlowMode := range []string{"absent", "empty", "short", "long", "zero", "changed"} {
				b512EmergencyFlowT.Run("ab15ap15_"+string(rune('A'+b512EmergencyFlowIndex))+b512EmergencyFlowName+b512EmergencyFlowMode, func(b512EmergencyFlowT *testing.T) {
					b512EmergencyFlowInput := b512EmergencyFlowCases(b512EmergencyFlowT)[b512EmergencyFlowIndex].B512EmergencyFlowInput
					b512EmergencyFlowRaw := b512EmergencyFlowInput["fields"].(map[string]any)
					b512EmergencyFlowOriginal := b512EmergencyFlowRaw[b512EmergencyFlowName].([]any)
					b512EmergencyFlowControl := b512EmergencyFlowEvaluate(b512EmergencyFlowT, b512EmergencyFlowInput)
					switch b512EmergencyFlowMode {
					case "absent":
						delete(b512EmergencyFlowRaw, b512EmergencyFlowName)
					case "empty":
						b512EmergencyFlowRaw[b512EmergencyFlowName] = []any{}
					case "short":
						b512EmergencyFlowRaw[b512EmergencyFlowName] = b512EmergencyFlowOriginal[:len(b512EmergencyFlowOriginal)-1]
					case "long":
						b512EmergencyFlowRaw[b512EmergencyFlowName] = append(b512EmergencyFlowOriginal, float64(0))
					case "zero", "changed":
						b512EmergencyFlowNew := make([]any, len(b512EmergencyFlowOriginal))
						for b512EmergencyFlowI := range b512EmergencyFlowNew {
							if b512EmergencyFlowMode == "zero" {
								b512EmergencyFlowNew[b512EmergencyFlowI] = float64(0)
							} else {
								b512EmergencyFlowNew[b512EmergencyFlowI] = float64(2)
							}
						}
						b512EmergencyFlowRaw[b512EmergencyFlowName] = b512EmergencyFlowNew
					}
					b512EmergencyFlowOutput := b512EmergencyFlowEvaluate(b512EmergencyFlowT, b512EmergencyFlowInput)
					b512EmergencyFlowFieldResult := b512EmergencyFlowField(b512EmergencyFlowOutput, b512EmergencyFlowName)
					if b512EmergencyFlowMode == "zero" || b512EmergencyFlowMode == "changed" {
						if b512EmergencyFlowOutput["state"] != "complete" || reflect.DeepEqual(b512EmergencyFlowFieldResult, b512EmergencyFlowField(b512EmergencyFlowControl, b512EmergencyFlowName)) {
							b512EmergencyFlowT.Fatal(b512EmergencyFlowOutput)
						}
						if b512EmergencyFlowName != "flags" && b512EmergencyFlowMode == "changed" {
							b512EmergencyFlowExpected := map[string]float64{"temperature1": 32.125, "temperature2": 32.125, "temperature3": 32.125, "pressure": 0.2, "flow": 20, "pump_status": 2, "switch": 2, "co2_1": 514, "co2_2": 514, "setpoint": 2}
							if b512EmergencyFlowFieldResult["number"] != b512EmergencyFlowExpected[b512EmergencyFlowName] {
								b512EmergencyFlowT.Fatal(b512EmergencyFlowName, b512EmergencyFlowFieldResult)
							}
						}
						if b512EmergencyFlowName != "flags" && b512EmergencyFlowMode == "zero" && b512EmergencyFlowFieldResult["number"] != float64(0) {
							b512EmergencyFlowT.Fatal(b512EmergencyFlowFieldResult)
						}
					} else if b512EmergencyFlowOutput["state"] != "partial" || b512EmergencyFlowFieldResult["state"] != "incomplete" || b512EmergencyFlowFieldResult["number"] != nil || b512EmergencyFlowFieldResult["label"] != nil {
						b512EmergencyFlowT.Fatal(b512EmergencyFlowOutput)
					}
					b512EmergencyFlowFieldRaw, b512EmergencyFlowHasRaw := b512EmergencyFlowFieldResult["raw"]
					if b512EmergencyFlowMode == "absent" {
						if b512EmergencyFlowHasRaw {
							b512EmergencyFlowT.Fatal("invented raw")
						}
					} else if !b512EmergencyFlowHasRaw || !reflect.DeepEqual(b512EmergencyFlowFieldRaw, b512EmergencyFlowOutput["raw"].(map[string]any)[b512EmergencyFlowName]) {
						b512EmergencyFlowT.Fatal("per-field raw")
					}
					for b512EmergencyFlowSibling := range b512EmergencyFlowRaw {
						if b512EmergencyFlowSibling != b512EmergencyFlowName && !reflect.DeepEqual(b512EmergencyFlowField(b512EmergencyFlowOutput, b512EmergencyFlowSibling), b512EmergencyFlowField(b512EmergencyFlowControl, b512EmergencyFlowSibling)) {
							b512EmergencyFlowT.Fatal("sibling lost", b512EmergencyFlowSibling)
						}
					}
					if !reflect.DeepEqual(b512EmergencyFlowOutput["raw"], b512EmergencyFlowRaw) {
						b512EmergencyFlowT.Fatal("raw changed")
					}
				})
			}
		}
	}
	for _, b512EmergencyFlowVector := range []struct {
		b512EmergencyFlowByte int
		b512EmergencyFlowWant float64
	}{{128, -1280}, {127, 1270}, {255, -10}, {0, 0}} {
		b512EmergencyFlowInput := b512EmergencyFlowCases(b512EmergencyFlowT)[0].B512EmergencyFlowInput
		b512EmergencyFlowInput["fields"].(map[string]any)["flow"] = []int{b512EmergencyFlowVector.b512EmergencyFlowByte}
		if b512EmergencyFlowField(b512EmergencyFlowEvaluate(b512EmergencyFlowT, b512EmergencyFlowInput), "flow")["number"] != b512EmergencyFlowVector.b512EmergencyFlowWant {
			b512EmergencyFlowT.Fatal(b512EmergencyFlowVector)
		}
	}
	b512EmergencyFlowInput := b512EmergencyFlowCases(b512EmergencyFlowT)[2].B512EmergencyFlowInput
	b512EmergencyFlowInput["fields"].(map[string]any)["co2_1"] = []int{1, 0}
	b512EmergencyFlowInput["fields"].(map[string]any)["co2_2"] = []int{255, 255}
	b512EmergencyFlowOutput := b512EmergencyFlowEvaluate(b512EmergencyFlowT, b512EmergencyFlowInput)
	if b512EmergencyFlowField(b512EmergencyFlowOutput, "co2_1")["number"] != float64(1) || b512EmergencyFlowField(b512EmergencyFlowOutput, "co2_2")["number"] != float64(65535) {
		b512EmergencyFlowT.Fatal("swapped control")
	}
}

func TestAb15Ap15ContextAndOwnership(b512EmergencyFlowT *testing.T) {
	for _, b512EmergencyFlowKey := range []string{"primary", "secondary", "operation", "direction", "view", "family", "variant", "opaque"} {
		for _, b512EmergencyFlowMode := range []string{"absent", "null", "zero", "changed"} {
			b512EmergencyFlowT.Run("ab15ap15_"+b512EmergencyFlowKey+b512EmergencyFlowMode, func(b512EmergencyFlowT *testing.T) {
				b512EmergencyFlowInput := b512EmergencyFlowCases(b512EmergencyFlowT)[0].B512EmergencyFlowInput
				b512EmergencyFlowContext := b512EmergencyFlowInput["context"].(map[string]any)
				switch b512EmergencyFlowMode {
				case "absent":
					delete(b512EmergencyFlowContext, b512EmergencyFlowKey)
				case "null":
					b512EmergencyFlowContext[b512EmergencyFlowKey] = nil
				default:
					switch b512EmergencyFlowKey {
					case "primary", "secondary", "operation":
						if b512EmergencyFlowMode == "zero" {
							b512EmergencyFlowContext[b512EmergencyFlowKey] = float64(0)
						} else {
							b512EmergencyFlowContext[b512EmergencyFlowKey] = float64(12)
						}
					case "direction", "view":
						b512EmergencyFlowContext[b512EmergencyFlowKey] = ""
						if b512EmergencyFlowMode == "changed" {
							b512EmergencyFlowContext[b512EmergencyFlowKey] = "other"
						}
					default:
						b512EmergencyFlowContext[b512EmergencyFlowKey] = map[string]any{"nested": []any{false, float64(0), nil}}
					}
				}
				b512EmergencyFlowOutput := b512EmergencyFlowEvaluate(b512EmergencyFlowT, b512EmergencyFlowInput)
				b512EmergencyFlowWant := "Unknown"
				if b512EmergencyFlowKey == "family" || b512EmergencyFlowKey == "variant" || b512EmergencyFlowKey == "opaque" {
					b512EmergencyFlowWant = "complete"
				}
				if b512EmergencyFlowOutput["state"] != b512EmergencyFlowWant || b512EmergencyFlowOutput["qualified"] != false || !reflect.DeepEqual(b512EmergencyFlowOutput["context"], b512EmergencyFlowContext) || !reflect.DeepEqual(b512EmergencyFlowOutput["raw"], b512EmergencyFlowInput["fields"]) {
					b512EmergencyFlowT.Fatal(b512EmergencyFlowOutput)
				}
				if b512EmergencyFlowWant == "Unknown" && len(b512EmergencyFlowOutput["fields"].(map[string]any)) != 0 {
					b512EmergencyFlowT.Fatal("unsupported formula")
				}
			})
		}
	}
	for _, b512EmergencyFlowKey := range []string{"family", "variant", "opaque"} {
		for _, b512EmergencyFlowValue := range []any{false, float64(0), "", []any{}, map[string]any{}} {
			b512EmergencyFlowInput := b512EmergencyFlowCases(b512EmergencyFlowT)[0].B512EmergencyFlowInput
			b512EmergencyFlowInput["context"].(map[string]any)[b512EmergencyFlowKey] = b512EmergencyFlowValue
			b512EmergencyFlowOutput := b512EmergencyFlowEvaluate(b512EmergencyFlowT, b512EmergencyFlowInput)
			if b512EmergencyFlowOutput["state"] != "complete" || b512EmergencyFlowOutput["qualified"] != false || !reflect.DeepEqual(b512EmergencyFlowOutput["context"], b512EmergencyFlowInput["context"]) {
				b512EmergencyFlowT.Fatal(b512EmergencyFlowKey, b512EmergencyFlowOutput)
			}
		}
	}
	for _, b512EmergencyFlowArgument := range []any{nil, false, float64(0), "", []any{}, []any{false, float64(0)}, map[string]any{"nested": false}} {
		b512EmergencyFlowInput := b512EmergencyFlowCases(b512EmergencyFlowT)[0].B512EmergencyFlowInput
		b512EmergencyFlowInput["arguments"] = map[string]any{"opaque": b512EmergencyFlowArgument}
		b512EmergencyFlowOutput := b512EmergencyFlowEvaluate(b512EmergencyFlowT, b512EmergencyFlowInput)
		if !reflect.DeepEqual(b512EmergencyFlowOutput["arguments"], b512EmergencyFlowInput["arguments"]) || b512EmergencyFlowOutput["state"] != "complete" {
			b512EmergencyFlowT.Fatal(b512EmergencyFlowOutput)
		}
	}
	b512EmergencyFlowAbsentInput := b512EmergencyFlowCases(b512EmergencyFlowT)[0].B512EmergencyFlowInput
	if _, b512EmergencyFlowExists := b512EmergencyFlowEvaluate(b512EmergencyFlowT, b512EmergencyFlowAbsentInput)["arguments"]; b512EmergencyFlowExists {
		b512EmergencyFlowT.Fatal("invented arguments")
	}
	b512EmergencyFlowAbsentInput["arguments"] = map[string]any{}
	if b512EmergencyFlowArguments, b512EmergencyFlowExists := b512EmergencyFlowEvaluate(b512EmergencyFlowT, b512EmergencyFlowAbsentInput)["arguments"]; !b512EmergencyFlowExists || !reflect.DeepEqual(b512EmergencyFlowArguments, map[string]any{}) {
		b512EmergencyFlowT.Fatal("empty arguments")
	}
	b512EmergencyFlowInput := b512EmergencyFlowCases(b512EmergencyFlowT)[0].B512EmergencyFlowInput
	b512EmergencyFlowInput["arguments"] = map[string]any{"opaque": []any{false, float64(0)}}
	b512EmergencyFlowInput["fields"].(map[string]any)["unmapped"] = []any{float64(255)}
	b512EmergencyFlowOutput := b512EmergencyFlowEvaluate(b512EmergencyFlowT, b512EmergencyFlowInput)
	if b512EmergencyFlowField(b512EmergencyFlowOutput, "unmapped")["state"] != "Unknown" || b512EmergencyFlowField(b512EmergencyFlowOutput, "unmapped")["number"] != nil {
		b512EmergencyFlowT.Fatal(b512EmergencyFlowOutput)
	}
	b512EmergencyFlowBytes, _ := json.Marshal(b512EmergencyFlowInput)
	b512EmergencyFlowReturned, b512EmergencyFlowErr := B512EmergencyFlowEvaluate(b512EmergencyFlowBytes)
	if b512EmergencyFlowErr != nil {
		b512EmergencyFlowT.Fatal(b512EmergencyFlowErr)
	}
	b512EmergencyFlowSaved := append([]byte(nil), b512EmergencyFlowReturned...)
	for b512EmergencyFlowI := range b512EmergencyFlowBytes {
		b512EmergencyFlowBytes[b512EmergencyFlowI] = ' '
	}
	if string(b512EmergencyFlowReturned) != string(b512EmergencyFlowSaved) {
		b512EmergencyFlowT.Fatal("input alias")
	}
	b512EmergencyFlowOutput["raw"].(map[string]any)["flow"].([]any)[0] = float64(0)
	b512EmergencyFlowOutput["context"].(map[string]any)["family"] = false
	b512EmergencyFlowOutput["arguments"].(map[string]any)["opaque"].([]any)[0] = true
	b512EmergencyFlowField(b512EmergencyFlowOutput, "flow")["raw"].([]any)[0] = float64(1)
	if !reflect.DeepEqual(b512EmergencyFlowInput, b512EmergencyFlowCasesWithExtras(b512EmergencyFlowT)) {
		b512EmergencyFlowT.Fatal("decoded input alias")
	}
	for b512EmergencyFlowI := range b512EmergencyFlowReturned {
		b512EmergencyFlowReturned[b512EmergencyFlowI] = ' '
	}
	b512EmergencyFlowNewBytes, _ := json.Marshal(b512EmergencyFlowInput)
	b512EmergencyFlowNewReturned, b512EmergencyFlowErr := B512EmergencyFlowEvaluate(b512EmergencyFlowNewBytes)
	if b512EmergencyFlowErr != nil || string(b512EmergencyFlowNewReturned) != string(b512EmergencyFlowSaved) {
		b512EmergencyFlowT.Fatal("later call changed", b512EmergencyFlowErr)
	}
}

func b512EmergencyFlowCasesWithExtras(b512EmergencyFlowT *testing.T) map[string]any {
	b512EmergencyFlowInput := b512EmergencyFlowCases(b512EmergencyFlowT)[0].B512EmergencyFlowInput
	b512EmergencyFlowInput["arguments"] = map[string]any{"opaque": []any{false, float64(0)}}
	b512EmergencyFlowInput["fields"].(map[string]any)["unmapped"] = []any{float64(255)}
	return b512EmergencyFlowInput
}

func TestAb15Ap15Reject(b512EmergencyFlowT *testing.T) {
	for b512EmergencyFlowIndex, b512EmergencyFlowBad := range []string{
		``, `{`, `[]`, `null`, `{} {}`, `{"context":{},"context":{},"fields":{}}`,
		`{"context":{"x":{"a":1,"a":2}},"fields":{}}`,
		`{"context":{},"fields":{"x":[{"a":1,"a":2}]}}`,
		`{"context":{},"fields":{},"arguments":{"x":[{"a":1,"a":2}]}}`,
		`{"context":null,"fields":{}}`, `{"context":{},"fields":null}`,
		`{"context":{},"fields":{},"arguments":false}`,
	} {
		b512EmergencyFlowOutput, b512EmergencyFlowErr := B512EmergencyFlowEvaluate([]byte(b512EmergencyFlowBad))
		if b512EmergencyFlowErr == nil || b512EmergencyFlowOutput != nil {
			b512EmergencyFlowT.Fatal(b512EmergencyFlowIndex, b512EmergencyFlowErr, string(b512EmergencyFlowOutput))
		}
	}
	for _, b512EmergencyFlowBad := range []string{`null`, `false`, `true`, `"0"`, `{}`, `[false]`, `[null]`, `["1"]`, `[1.5]`, `[-1]`, `[256]`, `[1e100]`, `[[]]`} {
		b512EmergencyFlowInput := `{"context":{},"fields":{"flags":` + b512EmergencyFlowBad + `}}`
		if b512EmergencyFlowOutput, b512EmergencyFlowErr := B512EmergencyFlowEvaluate([]byte(b512EmergencyFlowInput)); b512EmergencyFlowErr == nil || b512EmergencyFlowOutput != nil {
			b512EmergencyFlowT.Fatal(b512EmergencyFlowBad)
		}
	}
	for _, b512EmergencyFlowKey := range []string{"primary", "secondary", "operation", "direction", "view"} {
		for _, b512EmergencyFlowBad := range []any{false, []any{}, map[string]any{}} {
			b512EmergencyFlowInput := b512EmergencyFlowCases(b512EmergencyFlowT)[0].B512EmergencyFlowInput
			b512EmergencyFlowInput["context"].(map[string]any)[b512EmergencyFlowKey] = b512EmergencyFlowBad
			b512EmergencyFlowBytes, _ := json.Marshal(b512EmergencyFlowInput)
			if b512EmergencyFlowOutput, b512EmergencyFlowErr := B512EmergencyFlowEvaluate(b512EmergencyFlowBytes); b512EmergencyFlowErr == nil || b512EmergencyFlowOutput != nil {
				b512EmergencyFlowT.Fatal(b512EmergencyFlowKey, b512EmergencyFlowBad)
			}
		}
	}
	for _, b512EmergencyFlowKey := range []string{"primary", "secondary", "operation"} {
		for _, b512EmergencyFlowBad := range []any{1.5, -1, 256, "181"} {
			b512EmergencyFlowInput := b512EmergencyFlowCases(b512EmergencyFlowT)[0].B512EmergencyFlowInput
			b512EmergencyFlowInput["context"].(map[string]any)[b512EmergencyFlowKey] = b512EmergencyFlowBad
			b512EmergencyFlowBytes, _ := json.Marshal(b512EmergencyFlowInput)
			if _, b512EmergencyFlowErr := B512EmergencyFlowEvaluate(b512EmergencyFlowBytes); b512EmergencyFlowErr == nil {
				b512EmergencyFlowT.Fatal(b512EmergencyFlowKey, b512EmergencyFlowBad)
			}
		}
	}
}

func TestAb15Ap15FlagPredicates(b512EmergencyFlowT *testing.T) {
	for b512EmergencyFlowIndex := 0; b512EmergencyFlowIndex < 4; b512EmergencyFlowIndex++ {
		for _, b512EmergencyFlowByte := range []int{0, 1, 2, 4, 64, 128, 255} {
			b512EmergencyFlowInput := b512EmergencyFlowCases(b512EmergencyFlowT)[b512EmergencyFlowIndex].B512EmergencyFlowInput
			b512EmergencyFlowInput["fields"].(map[string]any)["flags"] = []int{b512EmergencyFlowByte}
			b512EmergencyFlowFlags := b512EmergencyFlowField(b512EmergencyFlowEvaluate(b512EmergencyFlowT, b512EmergencyFlowInput), "flags")
			b512EmergencyFlowLabel := []string{"compressor-off", "compressor-off", "emergency", "LED"}[b512EmergencyFlowIndex]
			if b512EmergencyFlowFlags["flags"].(map[string]any)[b512EmergencyFlowLabel] != (b512EmergencyFlowByte%2 == 1) {
				b512EmergencyFlowT.Fatal(b512EmergencyFlowFlags)
			}
			if b512EmergencyFlowIndex == 1 {
				if b512EmergencyFlowFlags["flags"].(map[string]any)["pump-relay"] != (b512EmergencyFlowByte/2%2 == 1) || b512EmergencyFlowFlags["reserved"] != float64(b512EmergencyFlowByte/4*4) || b512EmergencyFlowFlags["unmodeled"] != nil {
					b512EmergencyFlowT.Fatal(b512EmergencyFlowFlags)
				}
			} else if b512EmergencyFlowFlags["reserved"] != float64(b512EmergencyFlowByte%128/2*2) || b512EmergencyFlowFlags["unmodeled"] != float64(b512EmergencyFlowByte/128*128) {
				b512EmergencyFlowT.Fatal(b512EmergencyFlowFlags)
			}
		}
	}
	for _, b512EmergencyFlowByte := range []int{0, 100, 101, 255} {
		b512EmergencyFlowInput := b512EmergencyFlowCases(b512EmergencyFlowT)[1].B512EmergencyFlowInput
		b512EmergencyFlowInput["fields"].(map[string]any)["setpoint"] = []int{b512EmergencyFlowByte}
		b512EmergencyFlowResult := b512EmergencyFlowField(b512EmergencyFlowEvaluate(b512EmergencyFlowT, b512EmergencyFlowInput), "setpoint")
		if b512EmergencyFlowResult["number"] != float64(b512EmergencyFlowByte) || b512EmergencyFlowResult["candidate_domain"] != (b512EmergencyFlowByte < 101) {
			b512EmergencyFlowT.Fatal(b512EmergencyFlowResult)
		}
	}
	for _, b512EmergencyFlowIndex := range []int{1, 3} {
		b512EmergencyFlowInput := b512EmergencyFlowCases(b512EmergencyFlowT)[b512EmergencyFlowIndex].B512EmergencyFlowInput
		b512EmergencyFlowInput["context"].(map[string]any)["operation"] = float64(12)
		b512EmergencyFlowOutput := b512EmergencyFlowEvaluate(b512EmergencyFlowT, b512EmergencyFlowInput)
		if b512EmergencyFlowOutput["state"] != "Unknown" || strings.Contains(string(b512EmergencyFlowMustJSON(b512EmergencyFlowOutput)), "declared_positive_prefix") {
			b512EmergencyFlowT.Fatal(b512EmergencyFlowOutput)
		}
	}
}

func b512EmergencyFlowMustJSON(b512EmergencyFlowValue any) []byte {
	b512EmergencyFlowBytes, _ := json.Marshal(b512EmergencyFlowValue)
	return b512EmergencyFlowBytes
}

func TestAb15Ap15EveryFieldRejectsWrongKinds(b512EmergencyFlowT *testing.T) {
	for b512EmergencyFlowIndex := 0; b512EmergencyFlowIndex < 4; b512EmergencyFlowIndex++ {
		for b512EmergencyFlowName := range b512EmergencyFlowCases(b512EmergencyFlowT)[b512EmergencyFlowIndex].B512EmergencyFlowInput["fields"].(map[string]any) {
			for _, b512EmergencyFlowBad := range []any{false, nil, "", float64(0), []any{false}, []any{float64(-1)}, []any{float64(256)}, []any{1.5}} {
				b512EmergencyFlowInput := b512EmergencyFlowCases(b512EmergencyFlowT)[b512EmergencyFlowIndex].B512EmergencyFlowInput
				b512EmergencyFlowInput["fields"].(map[string]any)[b512EmergencyFlowName] = b512EmergencyFlowBad
				b512EmergencyFlowBytes, _ := json.Marshal(b512EmergencyFlowInput)
				if b512EmergencyFlowOutput, b512EmergencyFlowErr := B512EmergencyFlowEvaluate(b512EmergencyFlowBytes); b512EmergencyFlowErr == nil || b512EmergencyFlowOutput != nil {
					b512EmergencyFlowT.Fatal(b512EmergencyFlowIndex, b512EmergencyFlowName, b512EmergencyFlowBad)
				}
			}
		}
	}
}
