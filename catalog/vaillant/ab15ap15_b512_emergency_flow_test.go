package vaillant

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"
)

type ab15ap15Case struct {
	Ab15Ap15ID      string             `json:"id"`
	Ab15Ap15Input   map[string]any     `json:"input"`
	Ab15Ap15Numbers map[string]float64 `json:"numbers"`
}

func ab15ap15Cases(ab15ap15T *testing.T) []ab15ap15Case {
	ab15ap15T.Helper()
	ab15ap15Data, ab15ap15Err := os.ReadFile("testdata/ab15ap15_b512_emergency_flow.json")
	if ab15ap15Err != nil {
		ab15ap15T.Fatal(ab15ap15Err)
	}
	var ab15ap15Cases []ab15ap15Case
	if ab15ap15Err := json.Unmarshal(ab15ap15Data, &ab15ap15Cases); ab15ap15Err != nil {
		ab15ap15T.Fatal(ab15ap15Err)
	}
	if len(ab15ap15Cases) != 4 {
		ab15ap15T.Fatal("fixture count")
	}
	return ab15ap15Cases
}

func ab15ap15Evaluate(ab15ap15T *testing.T, ab15ap15Input any) map[string]any {
	ab15ap15T.Helper()
	ab15ap15Data, ab15ap15Err := json.Marshal(ab15ap15Input)
	if ab15ap15Err != nil {
		ab15ap15T.Fatal(ab15ap15Err)
	}
	ab15ap15Output, ab15ap15Err := Ab15Ap15Evaluate(ab15ap15Data)
	if ab15ap15Err != nil {
		ab15ap15T.Fatal(ab15ap15Err)
	}
	var ab15ap15Result map[string]any
	if ab15ap15Err := json.Unmarshal(ab15ap15Output, &ab15ap15Result); ab15ap15Err != nil {
		ab15ap15T.Fatal(ab15ap15Err)
	}
	return ab15ap15Result
}

func ab15ap15Field(ab15ap15Output map[string]any, ab15ap15Name string) map[string]any {
	return ab15ap15Output["fields"].(map[string]any)[ab15ap15Name].(map[string]any)
}

func TestAb15Ap15Synthetic(ab15ap15T *testing.T) {
	for _, ab15ap15Case := range ab15ap15Cases(ab15ap15T) {
		ab15ap15T.Run(ab15ap15Case.Ab15Ap15ID, func(ab15ap15T *testing.T) {
			ab15ap15Output := ab15ap15Evaluate(ab15ap15T, ab15ap15Case.Ab15Ap15Input)
			if ab15ap15Output["state"] != "complete" || ab15ap15Output["physical"] != "Unknown" || ab15ap15Output["qualified"] != false {
				ab15ap15T.Fatal(ab15ap15Output)
			}
			if !reflect.DeepEqual(ab15ap15Output["context"], ab15ap15Case.Ab15Ap15Input["context"]) || !reflect.DeepEqual(ab15ap15Output["raw"], ab15ap15Case.Ab15Ap15Input["fields"]) {
				ab15ap15T.Fatal("preservation")
			}
			for ab15ap15Name, ab15ap15Number := range ab15ap15Case.Ab15Ap15Numbers {
				ab15ap15Field := ab15ap15Field(ab15ap15Output, ab15ap15Name)
				if ab15ap15Field["number"] != ab15ap15Number || ab15ap15Field["state"] != "Hypothesis" || ab15ap15Field["calculation"] != "Proven" {
					ab15ap15T.Fatal(ab15ap15Name, ab15ap15Field)
				}
			}
			ab15ap15Context := ab15ap15Case.Ab15Ap15Input["context"].(map[string]any)
			ab15ap15Flags := ab15ap15Field(ab15ap15Output, "flags")
			if ab15ap15Context["direction"] == "positive" {
				ab15ap15Candidate := ab15ap15Output["candidate"].(map[string]any)
				if !reflect.DeepEqual(ab15ap15Candidate["declared_positive_prefix"], []any{float64(181), float64(18), float64(12)}) || ab15ap15Candidate["literal_reply_conflict"] != true || ab15ap15Candidate["physical_association"] != "Unknown" || ab15ap15Candidate["requested_operation"] != ab15ap15Context["operation"] {
					ab15ap15T.Fatal(ab15ap15Candidate)
				}
			}
			if ab15ap15Context["view"] == "flow-io" {
				if ab15ap15Flags["flags"].(map[string]any)["compressor-off"] != true {
					ab15ap15T.Fatal(ab15ap15Flags)
				}
				if ab15ap15Context["direction"] == "positive" {
					if ab15ap15Flags["reserved"] != float64(128) || ab15ap15Flags["unmodeled"] != nil || ab15ap15Field(ab15ap15Output, "setpoint")["candidate_domain"] != false {
						ab15ap15T.Fatal(ab15ap15Output)
					}
				} else if ab15ap15Flags["unmodeled"] != float64(128) || ab15ap15Flags["reserved"] != float64(0) {
					ab15ap15T.Fatal(ab15ap15Flags)
				}
			} else if ab15ap15Flags["unmodeled"] != float64(128) || ab15ap15Flags["reserved"] != float64(0) {
				ab15ap15T.Fatal(ab15ap15Flags)
			}
		})
	}
}

func TestAb15Ap15Fields(ab15ap15T *testing.T) {
	for ab15ap15Index := 0; ab15ap15Index < 4; ab15ap15Index++ {
		for ab15ap15Name := range ab15ap15Cases(ab15ap15T)[ab15ap15Index].Ab15Ap15Input["fields"].(map[string]any) {
			for _, ab15ap15Mode := range []string{"absent", "empty", "short", "long", "zero", "changed"} {
				ab15ap15T.Run("ab15ap15_"+string(rune('A'+ab15ap15Index))+ab15ap15Name+ab15ap15Mode, func(ab15ap15T *testing.T) {
					ab15ap15Input := ab15ap15Cases(ab15ap15T)[ab15ap15Index].Ab15Ap15Input
					ab15ap15Raw := ab15ap15Input["fields"].(map[string]any)
					ab15ap15Original := ab15ap15Raw[ab15ap15Name].([]any)
					ab15ap15Control := ab15ap15Evaluate(ab15ap15T, ab15ap15Input)
					switch ab15ap15Mode {
					case "absent":
						delete(ab15ap15Raw, ab15ap15Name)
					case "empty":
						ab15ap15Raw[ab15ap15Name] = []any{}
					case "short":
						ab15ap15Raw[ab15ap15Name] = ab15ap15Original[:len(ab15ap15Original)-1]
					case "long":
						ab15ap15Raw[ab15ap15Name] = append(ab15ap15Original, float64(0))
					case "zero", "changed":
						ab15ap15New := make([]any, len(ab15ap15Original))
						for ab15ap15I := range ab15ap15New {
							if ab15ap15Mode == "zero" {
								ab15ap15New[ab15ap15I] = float64(0)
							} else {
								ab15ap15New[ab15ap15I] = float64(2)
							}
						}
						ab15ap15Raw[ab15ap15Name] = ab15ap15New
					}
					ab15ap15Output := ab15ap15Evaluate(ab15ap15T, ab15ap15Input)
					ab15ap15FieldResult := ab15ap15Field(ab15ap15Output, ab15ap15Name)
					if ab15ap15Mode == "zero" || ab15ap15Mode == "changed" {
						if ab15ap15Output["state"] != "complete" || reflect.DeepEqual(ab15ap15FieldResult, ab15ap15Field(ab15ap15Control, ab15ap15Name)) {
							ab15ap15T.Fatal(ab15ap15Output)
						}
						if ab15ap15Name != "flags" && ab15ap15Mode == "changed" {
							ab15ap15Expected := map[string]float64{"temperature1": 32.125, "temperature2": 32.125, "temperature3": 32.125, "pressure": 0.2, "flow": 20, "pump_status": 2, "switch": 2, "co2_1": 514, "co2_2": 514, "setpoint": 2}
							if ab15ap15FieldResult["number"] != ab15ap15Expected[ab15ap15Name] {
								ab15ap15T.Fatal(ab15ap15Name, ab15ap15FieldResult)
							}
						}
						if ab15ap15Name != "flags" && ab15ap15Mode == "zero" && ab15ap15FieldResult["number"] != float64(0) {
							ab15ap15T.Fatal(ab15ap15FieldResult)
						}
					} else if ab15ap15Output["state"] != "partial" || ab15ap15FieldResult["state"] != "incomplete" || ab15ap15FieldResult["number"] != nil || ab15ap15FieldResult["label"] != nil {
						ab15ap15T.Fatal(ab15ap15Output)
					}
					ab15ap15FieldRaw, ab15ap15HasRaw := ab15ap15FieldResult["raw"]
					if ab15ap15Mode == "absent" {
						if ab15ap15HasRaw {
							ab15ap15T.Fatal("invented raw")
						}
					} else if !ab15ap15HasRaw || !reflect.DeepEqual(ab15ap15FieldRaw, ab15ap15Output["raw"].(map[string]any)[ab15ap15Name]) {
						ab15ap15T.Fatal("per-field raw")
					}
					for ab15ap15Sibling := range ab15ap15Raw {
						if ab15ap15Sibling != ab15ap15Name && !reflect.DeepEqual(ab15ap15Field(ab15ap15Output, ab15ap15Sibling), ab15ap15Field(ab15ap15Control, ab15ap15Sibling)) {
							ab15ap15T.Fatal("sibling lost", ab15ap15Sibling)
						}
					}
					if !reflect.DeepEqual(ab15ap15Output["raw"], ab15ap15Raw) {
						ab15ap15T.Fatal("raw changed")
					}
				})
			}
		}
	}
	for _, ab15ap15Vector := range []struct {
		ab15ap15Byte int
		ab15ap15Want float64
	}{{128, -1280}, {127, 1270}, {255, -10}, {0, 0}} {
		ab15ap15Input := ab15ap15Cases(ab15ap15T)[0].Ab15Ap15Input
		ab15ap15Input["fields"].(map[string]any)["flow"] = []int{ab15ap15Vector.ab15ap15Byte}
		if ab15ap15Field(ab15ap15Evaluate(ab15ap15T, ab15ap15Input), "flow")["number"] != ab15ap15Vector.ab15ap15Want {
			ab15ap15T.Fatal(ab15ap15Vector)
		}
	}
	ab15ap15Input := ab15ap15Cases(ab15ap15T)[2].Ab15Ap15Input
	ab15ap15Input["fields"].(map[string]any)["co2_1"] = []int{1, 0}
	ab15ap15Input["fields"].(map[string]any)["co2_2"] = []int{255, 255}
	ab15ap15Output := ab15ap15Evaluate(ab15ap15T, ab15ap15Input)
	if ab15ap15Field(ab15ap15Output, "co2_1")["number"] != float64(1) || ab15ap15Field(ab15ap15Output, "co2_2")["number"] != float64(65535) {
		ab15ap15T.Fatal("swapped control")
	}
}

func TestAb15Ap15ContextAndOwnership(ab15ap15T *testing.T) {
	for _, ab15ap15Key := range []string{"primary", "secondary", "operation", "direction", "view", "family", "variant", "opaque"} {
		for _, ab15ap15Mode := range []string{"absent", "null", "zero", "changed"} {
			ab15ap15T.Run("ab15ap15_"+ab15ap15Key+ab15ap15Mode, func(ab15ap15T *testing.T) {
				ab15ap15Input := ab15ap15Cases(ab15ap15T)[0].Ab15Ap15Input
				ab15ap15Context := ab15ap15Input["context"].(map[string]any)
				switch ab15ap15Mode {
				case "absent":
					delete(ab15ap15Context, ab15ap15Key)
				case "null":
					ab15ap15Context[ab15ap15Key] = nil
				default:
					switch ab15ap15Key {
					case "primary", "secondary", "operation":
						if ab15ap15Mode == "zero" {
							ab15ap15Context[ab15ap15Key] = float64(0)
						} else {
							ab15ap15Context[ab15ap15Key] = float64(12)
						}
					case "direction", "view":
						ab15ap15Context[ab15ap15Key] = ""
						if ab15ap15Mode == "changed" {
							ab15ap15Context[ab15ap15Key] = "other"
						}
					default:
						ab15ap15Context[ab15ap15Key] = map[string]any{"nested": []any{false, float64(0), nil}}
					}
				}
				ab15ap15Output := ab15ap15Evaluate(ab15ap15T, ab15ap15Input)
				ab15ap15Want := "Unknown"
				if ab15ap15Key == "family" || ab15ap15Key == "variant" || ab15ap15Key == "opaque" {
					ab15ap15Want = "complete"
				}
				if ab15ap15Output["state"] != ab15ap15Want || ab15ap15Output["qualified"] != false || !reflect.DeepEqual(ab15ap15Output["context"], ab15ap15Context) || !reflect.DeepEqual(ab15ap15Output["raw"], ab15ap15Input["fields"]) {
					ab15ap15T.Fatal(ab15ap15Output)
				}
				if ab15ap15Want == "Unknown" && len(ab15ap15Output["fields"].(map[string]any)) != 0 {
					ab15ap15T.Fatal("unsupported formula")
				}
			})
		}
	}
	for _, ab15ap15Key := range []string{"family", "variant", "opaque"} {
		for _, ab15ap15Value := range []any{false, float64(0), "", []any{}, map[string]any{}} {
			ab15ap15Input := ab15ap15Cases(ab15ap15T)[0].Ab15Ap15Input
			ab15ap15Input["context"].(map[string]any)[ab15ap15Key] = ab15ap15Value
			ab15ap15Output := ab15ap15Evaluate(ab15ap15T, ab15ap15Input)
			if ab15ap15Output["state"] != "complete" || ab15ap15Output["qualified"] != false || !reflect.DeepEqual(ab15ap15Output["context"], ab15ap15Input["context"]) {
				ab15ap15T.Fatal(ab15ap15Key, ab15ap15Output)
			}
		}
	}
	for _, ab15ap15Argument := range []any{nil, false, float64(0), "", []any{}, []any{false, float64(0)}, map[string]any{"nested": false}} {
		ab15ap15Input := ab15ap15Cases(ab15ap15T)[0].Ab15Ap15Input
		ab15ap15Input["arguments"] = map[string]any{"opaque": ab15ap15Argument}
		ab15ap15Output := ab15ap15Evaluate(ab15ap15T, ab15ap15Input)
		if !reflect.DeepEqual(ab15ap15Output["arguments"], ab15ap15Input["arguments"]) || ab15ap15Output["state"] != "complete" {
			ab15ap15T.Fatal(ab15ap15Output)
		}
	}
	ab15ap15AbsentInput := ab15ap15Cases(ab15ap15T)[0].Ab15Ap15Input
	if _, ab15ap15Exists := ab15ap15Evaluate(ab15ap15T, ab15ap15AbsentInput)["arguments"]; ab15ap15Exists {
		ab15ap15T.Fatal("invented arguments")
	}
	ab15ap15AbsentInput["arguments"] = map[string]any{}
	if ab15ap15Arguments, ab15ap15Exists := ab15ap15Evaluate(ab15ap15T, ab15ap15AbsentInput)["arguments"]; !ab15ap15Exists || !reflect.DeepEqual(ab15ap15Arguments, map[string]any{}) {
		ab15ap15T.Fatal("empty arguments")
	}
	ab15ap15Input := ab15ap15Cases(ab15ap15T)[0].Ab15Ap15Input
	ab15ap15Input["arguments"] = map[string]any{"opaque": []any{false, float64(0)}}
	ab15ap15Input["fields"].(map[string]any)["unmapped"] = []any{float64(255)}
	ab15ap15Output := ab15ap15Evaluate(ab15ap15T, ab15ap15Input)
	if ab15ap15Field(ab15ap15Output, "unmapped")["state"] != "Unknown" || ab15ap15Field(ab15ap15Output, "unmapped")["number"] != nil {
		ab15ap15T.Fatal(ab15ap15Output)
	}
	ab15ap15Bytes, _ := json.Marshal(ab15ap15Input)
	ab15ap15Returned, ab15ap15Err := Ab15Ap15Evaluate(ab15ap15Bytes)
	if ab15ap15Err != nil {
		ab15ap15T.Fatal(ab15ap15Err)
	}
	ab15ap15Saved := append([]byte(nil), ab15ap15Returned...)
	for ab15ap15I := range ab15ap15Bytes {
		ab15ap15Bytes[ab15ap15I] = ' '
	}
	if string(ab15ap15Returned) != string(ab15ap15Saved) {
		ab15ap15T.Fatal("input alias")
	}
	ab15ap15Output["raw"].(map[string]any)["flow"].([]any)[0] = float64(0)
	ab15ap15Output["context"].(map[string]any)["family"] = false
	ab15ap15Output["arguments"].(map[string]any)["opaque"].([]any)[0] = true
	ab15ap15Field(ab15ap15Output, "flow")["raw"].([]any)[0] = float64(1)
	if !reflect.DeepEqual(ab15ap15Input, ab15ap15CasesWithExtras(ab15ap15T)) {
		ab15ap15T.Fatal("decoded input alias")
	}
	for ab15ap15I := range ab15ap15Returned {
		ab15ap15Returned[ab15ap15I] = ' '
	}
	ab15ap15NewBytes, _ := json.Marshal(ab15ap15Input)
	ab15ap15NewReturned, ab15ap15Err := Ab15Ap15Evaluate(ab15ap15NewBytes)
	if ab15ap15Err != nil || string(ab15ap15NewReturned) != string(ab15ap15Saved) {
		ab15ap15T.Fatal("later call changed", ab15ap15Err)
	}
}

func ab15ap15CasesWithExtras(ab15ap15T *testing.T) map[string]any {
	ab15ap15Input := ab15ap15Cases(ab15ap15T)[0].Ab15Ap15Input
	ab15ap15Input["arguments"] = map[string]any{"opaque": []any{false, float64(0)}}
	ab15ap15Input["fields"].(map[string]any)["unmapped"] = []any{float64(255)}
	return ab15ap15Input
}

func TestAb15Ap15Reject(ab15ap15T *testing.T) {
	for ab15ap15Index, ab15ap15Bad := range []string{
		``, `{`, `[]`, `null`, `{} {}`, `{"context":{},"context":{},"fields":{}}`,
		`{"context":{"x":{"a":1,"a":2}},"fields":{}}`,
		`{"context":{},"fields":{"x":[{"a":1,"a":2}]}}`,
		`{"context":{},"fields":{},"arguments":{"x":[{"a":1,"a":2}]}}`,
		`{"context":null,"fields":{}}`, `{"context":{},"fields":null}`,
		`{"context":{},"fields":{},"arguments":false}`,
	} {
		ab15ap15Output, ab15ap15Err := Ab15Ap15Evaluate([]byte(ab15ap15Bad))
		if ab15ap15Err == nil || ab15ap15Output != nil {
			ab15ap15T.Fatal(ab15ap15Index, ab15ap15Err, string(ab15ap15Output))
		}
	}
	for _, ab15ap15Bad := range []string{`null`, `false`, `true`, `"0"`, `{}`, `[false]`, `[null]`, `["1"]`, `[1.5]`, `[-1]`, `[256]`, `[1e100]`, `[[]]`} {
		ab15ap15Input := `{"context":{},"fields":{"flags":` + ab15ap15Bad + `}}`
		if ab15ap15Output, ab15ap15Err := Ab15Ap15Evaluate([]byte(ab15ap15Input)); ab15ap15Err == nil || ab15ap15Output != nil {
			ab15ap15T.Fatal(ab15ap15Bad)
		}
	}
	for _, ab15ap15Key := range []string{"primary", "secondary", "operation", "direction", "view"} {
		for _, ab15ap15Bad := range []any{false, []any{}, map[string]any{}} {
			ab15ap15Input := ab15ap15Cases(ab15ap15T)[0].Ab15Ap15Input
			ab15ap15Input["context"].(map[string]any)[ab15ap15Key] = ab15ap15Bad
			ab15ap15Bytes, _ := json.Marshal(ab15ap15Input)
			if ab15ap15Output, ab15ap15Err := Ab15Ap15Evaluate(ab15ap15Bytes); ab15ap15Err == nil || ab15ap15Output != nil {
				ab15ap15T.Fatal(ab15ap15Key, ab15ap15Bad)
			}
		}
	}
	for _, ab15ap15Key := range []string{"primary", "secondary", "operation"} {
		for _, ab15ap15Bad := range []any{1.5, -1, 256, "181"} {
			ab15ap15Input := ab15ap15Cases(ab15ap15T)[0].Ab15Ap15Input
			ab15ap15Input["context"].(map[string]any)[ab15ap15Key] = ab15ap15Bad
			ab15ap15Bytes, _ := json.Marshal(ab15ap15Input)
			if _, ab15ap15Err := Ab15Ap15Evaluate(ab15ap15Bytes); ab15ap15Err == nil {
				ab15ap15T.Fatal(ab15ap15Key, ab15ap15Bad)
			}
		}
	}
}

func TestAb15Ap15FlagPredicates(ab15ap15T *testing.T) {
	for ab15ap15Index := 0; ab15ap15Index < 4; ab15ap15Index++ {
		for _, ab15ap15Byte := range []int{0, 1, 2, 4, 64, 128, 255} {
			ab15ap15Input := ab15ap15Cases(ab15ap15T)[ab15ap15Index].Ab15Ap15Input
			ab15ap15Input["fields"].(map[string]any)["flags"] = []int{ab15ap15Byte}
			ab15ap15Flags := ab15ap15Field(ab15ap15Evaluate(ab15ap15T, ab15ap15Input), "flags")
			ab15ap15Label := []string{"compressor-off", "compressor-off", "emergency", "LED"}[ab15ap15Index]
			if ab15ap15Flags["flags"].(map[string]any)[ab15ap15Label] != (ab15ap15Byte%2 == 1) {
				ab15ap15T.Fatal(ab15ap15Flags)
			}
			if ab15ap15Index == 1 {
				if ab15ap15Flags["flags"].(map[string]any)["pump-relay"] != (ab15ap15Byte/2%2 == 1) || ab15ap15Flags["reserved"] != float64(ab15ap15Byte/4*4) || ab15ap15Flags["unmodeled"] != nil {
					ab15ap15T.Fatal(ab15ap15Flags)
				}
			} else if ab15ap15Flags["reserved"] != float64(ab15ap15Byte%128/2*2) || ab15ap15Flags["unmodeled"] != float64(ab15ap15Byte/128*128) {
				ab15ap15T.Fatal(ab15ap15Flags)
			}
		}
	}
	for _, ab15ap15Byte := range []int{0, 100, 101, 255} {
		ab15ap15Input := ab15ap15Cases(ab15ap15T)[1].Ab15Ap15Input
		ab15ap15Input["fields"].(map[string]any)["setpoint"] = []int{ab15ap15Byte}
		ab15ap15Result := ab15ap15Field(ab15ap15Evaluate(ab15ap15T, ab15ap15Input), "setpoint")
		if ab15ap15Result["number"] != float64(ab15ap15Byte) || ab15ap15Result["candidate_domain"] != (ab15ap15Byte < 101) {
			ab15ap15T.Fatal(ab15ap15Result)
		}
	}
	for _, ab15ap15Index := range []int{1, 3} {
		ab15ap15Input := ab15ap15Cases(ab15ap15T)[ab15ap15Index].Ab15Ap15Input
		ab15ap15Input["context"].(map[string]any)["operation"] = float64(12)
		ab15ap15Output := ab15ap15Evaluate(ab15ap15T, ab15ap15Input)
		if ab15ap15Output["state"] != "Unknown" || strings.Contains(string(ab15ap15MustJSON(ab15ap15Output)), "declared_positive_prefix") {
			ab15ap15T.Fatal(ab15ap15Output)
		}
	}
}

func ab15ap15MustJSON(ab15ap15Value any) []byte {
	ab15ap15Bytes, _ := json.Marshal(ab15ap15Value)
	return ab15ap15Bytes
}

func TestAb15Ap15EveryFieldRejectsWrongKinds(ab15ap15T *testing.T) {
	for ab15ap15Index := 0; ab15ap15Index < 4; ab15ap15Index++ {
		for ab15ap15Name := range ab15ap15Cases(ab15ap15T)[ab15ap15Index].Ab15Ap15Input["fields"].(map[string]any) {
			for _, ab15ap15Bad := range []any{false, nil, "", float64(0), []any{false}, []any{float64(-1)}, []any{float64(256)}, []any{1.5}} {
				ab15ap15Input := ab15ap15Cases(ab15ap15T)[ab15ap15Index].Ab15Ap15Input
				ab15ap15Input["fields"].(map[string]any)[ab15ap15Name] = ab15ap15Bad
				ab15ap15Bytes, _ := json.Marshal(ab15ap15Input)
				if ab15ap15Output, ab15ap15Err := Ab15Ap15Evaluate(ab15ap15Bytes); ab15ap15Err == nil || ab15ap15Output != nil {
					ab15ap15T.Fatal(ab15ap15Index, ab15ap15Name, ab15ap15Bad)
				}
			}
		}
	}
}
