package vaillant

import (
	ab05ap05json "encoding/json"
	ab05ap05os "os"
	ab05ap05reflect "reflect"
	ab05ap05strings "strings"
	ab05ap05testing "testing"
)

func ab05ap05Input() map[string]any {
	ab05ap05Slots := []any{}
	for ab05ap05Index := 0; ab05ap05Index < 9; ab05ap05Index++ {
		ab05ap05Slots = append(ab05ap05Slots, map[string]any{"registered": true, "capable": true, "prior": 10, "prior_previous": 7, "prior_changed": false, "outcome": "success", "value": 0})
	}
	return map[string]any{"context": map[string]any{"primary": 181, "secondary": 22, "operation": 14, "direction": "local", "view": "cycle", "family": "synthetic-family", "variant": "candidate"}, "fields": map[string]any{"opaque": []any{0, 255}}, "arguments": map[string]any{"slots": ab05ap05Slots, "failure_counter": 19}}
}
func ab05ap05Run(ab05ap05T *ab05ap05testing.T, ab05ap05InputValue map[string]any) map[string]any {
	ab05ap05T.Helper()
	ab05ap05Encoded, ab05ap05Err := ab05ap05json.Marshal(ab05ap05InputValue)
	if ab05ap05Err != nil {
		ab05ap05T.Fatal(ab05ap05Err)
	}
	ab05ap05Output, ab05ap05Err := Ab05Ap05Evaluate(ab05ap05Encoded)
	if ab05ap05Err != nil {
		ab05ap05T.Fatal(ab05ap05Err)
	}
	ab05ap05Decoded := map[string]any{}
	if ab05ap05Err = ab05ap05json.Unmarshal(ab05ap05Output, &ab05ap05Decoded); ab05ap05Err != nil {
		ab05ap05T.Fatal(ab05ap05Err)
	}
	return ab05ap05Decoded
}
func ab05ap05Args(ab05ap05InputValue map[string]any) map[string]any {
	return ab05ap05InputValue["arguments"].(map[string]any)
}
func ab05ap05Slot(ab05ap05InputValue map[string]any, ab05ap05Index int) map[string]any {
	return ab05ap05Args(ab05ap05InputValue)["slots"].([]any)[ab05ap05Index].(map[string]any)
}
func ab05ap05Fields(ab05ap05Output map[string]any) map[string]any {
	return ab05ap05Output["fields"].(map[string]any)
}
func ab05ap05Outcome(ab05ap05Output map[string]any, ab05ap05Index int) map[string]any {
	return ab05ap05Fields(ab05ap05Output)["slots"].(map[string]any)["outcomes"].([]any)[ab05ap05Index].(map[string]any)
}
func ab05ap05Equal(ab05ap05T *ab05ap05testing.T, ab05ap05Actual, ab05ap05Expected any) {
	ab05ap05T.Helper()
	if !ab05ap05reflect.DeepEqual(ab05ap05Actual, ab05ap05Expected) {
		ab05ap05T.Fatalf("got %#v want %#v", ab05ap05Actual, ab05ap05Expected)
	}
}
func ab05ap05Number(ab05ap05Value any) any { return ab05ap05Value.(map[string]any)["number"] }
func ab05ap05Reject(ab05ap05T *ab05ap05testing.T, ab05ap05InputValue map[string]any) {
	ab05ap05T.Helper()
	ab05ap05Encoded, _ := ab05ap05json.Marshal(ab05ap05InputValue)
	ab05ap05Result, ab05ap05Err := Ab05Ap05Evaluate(ab05ap05Encoded)
	if ab05ap05Err == nil || ab05ap05Result != nil {
		ab05ap05T.Fatalf("accepted invalid input: %s", ab05ap05Encoded)
	}
}

func TestAb05Ap05Transitions(ab05ap05T *ab05ap05testing.T) {
	ab05ap05Data, ab05ap05Err := ab05ap05os.ReadFile("testdata/ab05ap05_b516_cycle_reducer.json")
	if ab05ap05Err != nil {
		ab05ap05T.Fatal(ab05ap05Err)
	}
	var ab05ap05Cases []map[string]any
	if ab05ap05Err = ab05ap05json.Unmarshal(ab05ap05Data, &ab05ap05Cases); ab05ap05Err != nil {
		ab05ap05T.Fatal(ab05ap05Err)
	}
	for _, ab05ap05Case := range ab05ap05Cases {
		for ab05ap05Index := 0; ab05ap05Index < 9; ab05ap05Index++ {
			ab05ap05InputValue := ab05ap05Input()
			ab05ap05SlotValue := ab05ap05Slot(ab05ap05InputValue, ab05ap05Index)
			ab05ap05SlotValue["prior"] = ab05ap05Case["prior"]
			ab05ap05SlotValue["value"] = ab05ap05Case["value"]
			ab05ap05Output := ab05ap05Run(ab05ap05T, ab05ap05InputValue)
			ab05ap05Result := ab05ap05Outcome(ab05ap05Output, ab05ap05Index)
			ab05ap05Equal(ab05ap05T, ab05ap05Result["changed"].(map[string]any)["value"], ab05ap05Case["changed"])
			ab05ap05Equal(ab05ap05T, ab05ap05Number(ab05ap05Result["current"]), ab05ap05Case["value"])
			ab05ap05Equal(ab05ap05T, ab05ap05Number(ab05ap05Result["previous"]), ab05ap05Case["prior"])
			ab05ap05Equal(ab05ap05T, ab05ap05Number(ab05ap05Fields(ab05ap05Output)["sum"]), ab05ap05Case["value"])
			ab05ap05Equal(ab05ap05T, ab05ap05Output["qualified"], false)
			ab05ap05Equal(ab05ap05T, ab05ap05Result["state"], "Hypothesis")
			ab05ap05Equal(ab05ap05T, ab05ap05Result["current"].(map[string]any)["state"], "Hypothesis")
			ab05ap05Equal(ab05ap05T, ab05ap05Result["current"].(map[string]any)["calculation"], "Proven")
			ab05ap05Equal(ab05ap05T, ab05ap05Output["physical"], "Unknown")
		}
	}
}
func TestAb05Ap05EligibilityOutcomes(ab05ap05T *ab05ap05testing.T) {
	for ab05ap05Index := 0; ab05ap05Index < 9; ab05ap05Index++ {
		for _, ab05ap05Registered := range []bool{false, true} {
			for _, ab05ap05Capable := range []bool{false, true} {
				for _, ab05ap05OutcomeValue := range []string{"success", "read-failure", "parse-failure"} {
					for _, ab05ap05Changed := range []bool{false, true} {
						ab05ap05InputValue := ab05ap05Input()
						ab05ap05SlotValue := ab05ap05Slot(ab05ap05InputValue, ab05ap05Index)
						ab05ap05SlotValue["registered"] = ab05ap05Registered
						ab05ap05SlotValue["capable"] = ab05ap05Capable
						ab05ap05SlotValue["prior_changed"] = ab05ap05Changed
						ab05ap05SlotValue["prior_previous"] = 123
						ab05ap05SlotValue["outcome"] = ab05ap05OutcomeValue
						ab05ap05SlotValue["value"] = 3
						if ab05ap05OutcomeValue != "success" {
							delete(ab05ap05SlotValue, "value")
						}
						ab05ap05Output := ab05ap05Run(ab05ap05T, ab05ap05InputValue)
						ab05ap05Result := ab05ap05Outcome(ab05ap05Output, ab05ap05Index)
						ab05ap05Counts := ab05ap05Fields(ab05ap05Output)["counts"].(map[string]any)
						ab05ap05Admitted := ab05ap05Registered && ab05ap05Capable
						ab05ap05ExpectedOutcome := "skipped"
						if ab05ap05Admitted {
							ab05ap05ExpectedOutcome = ab05ap05OutcomeValue
						}
						ab05ap05Equal(ab05ap05T, ab05ap05Result["outcome"], ab05ap05ExpectedOutcome)
						ab05ap05ExpectedSuccessful := float64(8)
						if ab05ap05Admitted && ab05ap05OutcomeValue == "success" {
							ab05ap05ExpectedSuccessful = 9
						}
						ab05ap05Equal(ab05ap05T, ab05ap05Counts["successful"], ab05ap05ExpectedSuccessful)
						ab05ap05Equal(ab05ap05T, ab05ap05Result["attempted"], ab05ap05Admitted)
						ab05ap05Equal(ab05ap05T, ab05ap05Result["prior_changed"], ab05ap05Changed)
						if !ab05ap05Admitted {
							ab05ap05Equal(ab05ap05T, ab05ap05Number(ab05ap05Result["previous"]), float64(123))
							ab05ap05Equal(ab05ap05T, ab05ap05Counts["skipped"], float64(1))
							ab05ap05Equal(ab05ap05T, ab05ap05Counts["failed"], float64(0))
							ab05ap05Equal(ab05ap05T, ab05ap05Counts["attempted"], float64(8))
						} else {
							ab05ap05Equal(ab05ap05T, ab05ap05Number(ab05ap05Result["previous"]), float64(10))
							ab05ap05Equal(ab05ap05T, ab05ap05Counts["attempted"], float64(9))
						}
						if !ab05ap05Admitted || ab05ap05OutcomeValue != "success" {
							ab05ap05Equal(ab05ap05T, ab05ap05Number(ab05ap05Result["current"]), float64(4294967295))
							ab05ap05Equal(ab05ap05T, ab05ap05Result["changed"].(map[string]any)["state"], "Unknown")
							if _, ab05ap05Present := ab05ap05Result["changed"].(map[string]any)["value"]; ab05ap05Present {
								ab05ap05T.Fatal("invented changed value")
							}
							ab05ap05Equal(ab05ap05T, ab05ap05Output["partial_cycle"], true)
						} else {
							ab05ap05Equal(ab05ap05T, ab05ap05Number(ab05ap05Result["current"]), float64(3))
							ab05ap05Equal(ab05ap05T, ab05ap05Output["partial_cycle"], false)
						}
						if ab05ap05Admitted && ab05ap05OutcomeValue != "success" {
							ab05ap05Equal(ab05ap05T, ab05ap05Counts["failed"], float64(1))
							ab05ap05Equal(ab05ap05T, ab05ap05Output["state"], "partial")
							ab05ap05Equal(ab05ap05T, ab05ap05Fields(ab05ap05Output)["failure_counter"].(map[string]any)["prior"], float64(19))
						} else {
							ab05ap05Equal(ab05ap05T, ab05ap05Number(ab05ap05Fields(ab05ap05Output)["failure_counter"]), float64(0))
						}
						ab05ap05ExpectedSum := float64(0)
						if ab05ap05Admitted && ab05ap05OutcomeValue == "success" {
							ab05ap05ExpectedSum = 3
						}
						ab05ap05Equal(ab05ap05T, ab05ap05Number(ab05ap05Fields(ab05ap05Output)["sum"]), ab05ap05ExpectedSum)
					}
				}
			}
		}
	}
}
func TestAb05Ap05WrapAndPartial(ab05ap05T *ab05ap05testing.T) {
	ab05ap05InputValue := ab05ap05Input()
	ab05ap05Slot(ab05ap05InputValue, 0)["value"] = uint32(4294967294)
	ab05ap05Slot(ab05ap05InputValue, 1)["value"] = 3
	ab05ap05Slot(ab05ap05InputValue, 2)["outcome"] = "read-failure"
	delete(ab05ap05Slot(ab05ap05InputValue, 2), "value")
	ab05ap05Args(ab05ap05InputValue)["failure_counter"] = uint32(4294967295)
	ab05ap05Args(ab05ap05InputValue)["counter_increment"] = 1
	ab05ap05Output := ab05ap05Run(ab05ap05T, ab05ap05InputValue)
	ab05ap05Equal(ab05ap05T, ab05ap05Number(ab05ap05Fields(ab05ap05Output)["sum"]), float64(1))
	ab05ap05Equal(ab05ap05T, ab05ap05Number(ab05ap05Fields(ab05ap05Output)["failure_counter"]), float64(0))
	delete(ab05ap05Args(ab05ap05InputValue), "counter_increment")
	ab05ap05Output = ab05ap05Run(ab05ap05T, ab05ap05InputValue)
	ab05ap05Equal(ab05ap05T, ab05ap05Fields(ab05ap05Output)["failure_counter"].(map[string]any)["prior"], float64(4294967295))
	ab05ap05Equal(ab05ap05T, ab05ap05Output["state"], "partial")
	ab05ap05Slot(ab05ap05InputValue, 2)["capable"] = false
	ab05ap05Args(ab05ap05InputValue)["counter_increment"] = 7
	ab05ap05Output = ab05ap05Run(ab05ap05T, ab05ap05InputValue)
	ab05ap05Equal(ab05ap05T, ab05ap05Number(ab05ap05Fields(ab05ap05Output)["failure_counter"]), float64(0))
}

func TestAb05Ap05ContextAndOwnership(ab05ap05T *ab05ap05testing.T) {
	for ab05ap05Key, ab05ap05Value := range map[string]any{"primary": 180, "secondary": 21, "operation": 13, "direction": "other", "view": "other"} {
		for _, ab05ap05Mode := range []string{"different", "absent", "null"} {
			ab05ap05InputValue := ab05ap05Input()
			ab05ap05Context := ab05ap05InputValue["context"].(map[string]any)
			ab05ap05Context[ab05ap05Key] = ab05ap05Value
			if ab05ap05Mode == "absent" {
				delete(ab05ap05Context, ab05ap05Key)
			}
			if ab05ap05Mode == "null" {
				ab05ap05Context[ab05ap05Key] = nil
			}
			ab05ap05Output := ab05ap05Run(ab05ap05T, ab05ap05InputValue)
			ab05ap05Equal(ab05ap05T, ab05ap05Output["state"], "Unknown")
			if _, ab05ap05Present := ab05ap05Fields(ab05ap05Output)["sum"]; ab05ap05Present {
				ab05ap05T.Fatal("formula on unsupported context")
			}
			ab05ap05Encoded, _ := ab05ap05json.Marshal(ab05ap05Context)
			var ab05ap05Expected any
			_ = ab05ap05json.Unmarshal(ab05ap05Encoded, &ab05ap05Expected)
			ab05ap05Equal(ab05ap05T, ab05ap05Output["context"], ab05ap05Expected)
		}
	}
	for _, ab05ap05Key := range []string{"family", "variant", "extra"} {
		for _, ab05ap05Value := range []any{nil, "", 0, false, []any{map[string]any{"nested": 1}}} {
			ab05ap05InputValue := ab05ap05Input()
			ab05ap05InputValue["context"].(map[string]any)[ab05ap05Key] = ab05ap05Value
			ab05ap05Output := ab05ap05Run(ab05ap05T, ab05ap05InputValue)
			ab05ap05Equal(ab05ap05T, ab05ap05Output["state"], "complete")
			ab05ap05Equal(ab05ap05T, ab05ap05Output["qualified"], false)
		}
	}
	ab05ap05InputValue := ab05ap05Input()
	ab05ap05Args(ab05ap05InputValue)["opaque"] = map[string]any{"nested": []any{false, nil, 0}}
	ab05ap05Encoded, _ := ab05ap05json.Marshal(ab05ap05InputValue)
	ab05ap05Output, ab05ap05Err := Ab05Ap05Evaluate(ab05ap05Encoded)
	if ab05ap05Err != nil {
		ab05ap05T.Fatal(ab05ap05Err)
	}
	ab05ap05Saved := append([]byte(nil), ab05ap05Output...)
	for ab05ap05Index := range ab05ap05Encoded {
		ab05ap05Encoded[ab05ap05Index] = ' '
	}
	ab05ap05Equal(ab05ap05T, []byte(ab05ap05Output), ab05ap05Saved)
	for ab05ap05Index := range ab05ap05Output {
		ab05ap05Output[ab05ap05Index] = ' '
	}
	ab05ap05Again := ab05ap05Run(ab05ap05T, ab05ap05InputValue)
	ab05ap05Again["raw"].(map[string]any)["opaque"].([]any)[0] = 99
	ab05ap05Equal(ab05ap05T, ab05ap05Fields(ab05ap05Again)["opaque"].(map[string]any)["raw"].([]any)[0], float64(0))
	ab05ap05Equal(ab05ap05T, ab05ap05Run(ab05ap05T, ab05ap05InputValue)["raw"].(map[string]any)["opaque"].([]any)[0], float64(0))
}

func TestAb05Ap05Validation(ab05ap05T *ab05ap05testing.T) {
	for _, ab05ap05Text := range []string{``, `{`, `[]`, `{} {}`, `{"context":{},"context":{},"fields":{}}`, `{"context":{"x":[{"a":1,"a":2}]},"fields":{}}`} {
		ab05ap05Result, ab05ap05Err := Ab05Ap05Evaluate([]byte(ab05ap05Text))
		if ab05ap05Err == nil || ab05ap05Result != nil {
			ab05ap05T.Fatalf("accepted %s", ab05ap05Text)
		}
	}
	for _, ab05ap05Key := range []string{"registered", "capable", "prior_changed", "prior", "prior_previous", "outcome", "value"} {
		for _, ab05ap05Bad := range []any{nil, []any{}, map[string]any{}, "bad", -1, 1.5, uint64(4294967296)} {
			ab05ap05InputValue := ab05ap05Input()
			ab05ap05Slot(ab05ap05InputValue, 0)[ab05ap05Key] = ab05ap05Bad
			ab05ap05Reject(ab05ap05T, ab05ap05InputValue)
		}
		ab05ap05InputValue := ab05ap05Input()
		delete(ab05ap05Slot(ab05ap05InputValue, 0), ab05ap05Key)
		ab05ap05Reject(ab05ap05T, ab05ap05InputValue)
	}
	for _, ab05ap05Key := range []string{"failure_counter", "counter_increment"} {
		for _, ab05ap05Bad := range []any{nil, false, "0", []any{}, -1, 0.5, uint64(4294967296)} {
			ab05ap05InputValue := ab05ap05Input()
			ab05ap05Args(ab05ap05InputValue)[ab05ap05Key] = ab05ap05Bad
			ab05ap05Reject(ab05ap05T, ab05ap05InputValue)
		}
	}
	for _, ab05ap05Count := range []int{0, 8, 10} {
		ab05ap05InputValue := ab05ap05Input()
		ab05ap05Slots := []any{}
		for ab05ap05Index := 0; ab05ap05Index < ab05ap05Count; ab05ap05Index++ {
			ab05ap05Slots = append(ab05ap05Slots, ab05ap05Slot(ab05ap05InputValue, 0))
		}
		ab05ap05Args(ab05ap05InputValue)["slots"] = ab05ap05Slots
		ab05ap05Reject(ab05ap05T, ab05ap05InputValue)
	}
	for _, ab05ap05Key := range []string{"slots", "failure_counter"} {
		ab05ap05InputValue := ab05ap05Input()
		delete(ab05ap05Args(ab05ap05InputValue), ab05ap05Key)
		ab05ap05Reject(ab05ap05T, ab05ap05InputValue)
	}
	for _, ab05ap05Bad := range []any{nil, false, 0, "bytes", []any{256}, []any{-1}, []any{true}, []any{0.5}, []any{"0"}} {
		ab05ap05InputValue := ab05ap05Input()
		ab05ap05InputValue["fields"].(map[string]any)["opaque"] = ab05ap05Bad
		ab05ap05Reject(ab05ap05T, ab05ap05InputValue)
	}
	for _, ab05ap05Value := range [][]any{{}, {0}, {255, 0, 1}} {
		ab05ap05InputValue := ab05ap05Input()
		ab05ap05InputValue["fields"].(map[string]any)["opaque"] = ab05ap05Value
		ab05ap05Output := ab05ap05Run(ab05ap05T, ab05ap05InputValue)
		ab05ap05Equal(ab05ap05T, ab05ap05Number(ab05ap05Fields(ab05ap05Output)["sum"]), float64(0))
		ab05ap05Equal(ab05ap05T, ab05ap05Fields(ab05ap05Output)["opaque"].(map[string]any)["state"], "Unknown")
	}
	ab05ap05InputValue := ab05ap05Input()
	ab05ap05Encoded, _ := ab05ap05json.Marshal(ab05ap05InputValue)
	for _, ab05ap05Replacement := range []string{`"prior":true`, `"prior":10.0`, `"prior":1e1`} {
		ab05ap05Text := ab05ap05strings.Replace(string(ab05ap05Encoded), `"prior":10`, ab05ap05Replacement, 1)
		if _, ab05ap05Err := Ab05Ap05Evaluate([]byte(ab05ap05Text)); ab05ap05Err == nil {
			ab05ap05T.Fatal("accepted wrong numeric kind")
		}
	}
}

func TestAb05Ap05Boundaries(ab05ap05T *ab05ap05testing.T) {
	for _, ab05ap05Key := range []string{"primary", "secondary", "operation"} {
		for _, ab05ap05Bad := range []any{false, "181", -1, 1.5, 256, []any{}} {
			ab05ap05InputValue := ab05ap05Input()
			ab05ap05InputValue["context"].(map[string]any)[ab05ap05Key] = ab05ap05Bad
			ab05ap05Reject(ab05ap05T, ab05ap05InputValue)
		}
	}
	for _, ab05ap05Key := range []string{"direction", "view"} {
		for _, ab05ap05Bad := range []any{false, 0, []any{}, map[string]any{}} {
			ab05ap05InputValue := ab05ap05Input()
			ab05ap05InputValue["context"].(map[string]any)[ab05ap05Key] = ab05ap05Bad
			ab05ap05Reject(ab05ap05T, ab05ap05InputValue)
		}
	}
	for _, ab05ap05Key := range []string{"prior", "prior_previous", "value"} {
		for _, ab05ap05Bad := range []any{false, true} {
			ab05ap05InputValue := ab05ap05Input()
			ab05ap05Slot(ab05ap05InputValue, 0)[ab05ap05Key] = ab05ap05Bad
			ab05ap05Reject(ab05ap05T, ab05ap05InputValue)
		}
	}
	for _, ab05ap05Key := range []string{"arguments", "context", "fields"} {
		ab05ap05InputValue := ab05ap05Input()
		delete(ab05ap05InputValue, ab05ap05Key)
		ab05ap05Reject(ab05ap05T, ab05ap05InputValue)
		for _, ab05ap05Bad := range []any{nil, false, []any{}, 0, ""} {
			ab05ap05InputValue := ab05ap05Input()
			ab05ap05InputValue[ab05ap05Key] = ab05ap05Bad
			ab05ap05Reject(ab05ap05T, ab05ap05InputValue)
		}
	}
	for _, ab05ap05Bad := range []any{nil, false, map[string]any{}, ""} {
		ab05ap05InputValue := ab05ap05Input()
		ab05ap05Args(ab05ap05InputValue)["slots"] = ab05ap05Bad
		ab05ap05Reject(ab05ap05T, ab05ap05InputValue)
	}
	ab05ap05InputValue := ab05ap05Input()
	ab05ap05Slot(ab05ap05InputValue, 0)["outcome"] = "parse-failure"
	ab05ap05Reject(ab05ap05T, ab05ap05InputValue)
	for _, ab05ap05Counter := range []uint32{0, 1, 4294967295} {
		for _, ab05ap05Increment := range []uint32{0, 1, 4294967295} {
			ab05ap05InputValue := ab05ap05Input()
			ab05ap05Args(ab05ap05InputValue)["failure_counter"] = ab05ap05Counter
			ab05ap05Args(ab05ap05InputValue)["counter_increment"] = ab05ap05Increment
			ab05ap05Output := ab05ap05Run(ab05ap05T, ab05ap05InputValue)
			ab05ap05Equal(ab05ap05T, ab05ap05Number(ab05ap05Fields(ab05ap05Output)["failure_counter"]), float64(0))
		}
	}
	for _, ab05ap05Key := range []string{"family", "variant", "opaque"} {
		for _, ab05ap05Value := range []any{nil, false, "", "other", 0, []any{map[string]any{"nested": false}}} {
			ab05ap05InputValue := ab05ap05Input()
			ab05ap05InputValue["context"].(map[string]any)[ab05ap05Key] = ab05ap05Value
			ab05ap05Args(ab05ap05InputValue)["extra"] = ab05ap05Value
			ab05ap05Output := ab05ap05Run(ab05ap05T, ab05ap05InputValue)
			ab05ap05Encoded, _ := ab05ap05json.Marshal(ab05ap05InputValue)
			var ab05ap05Expected map[string]any
			_ = ab05ap05json.Unmarshal(ab05ap05Encoded, &ab05ap05Expected)
			ab05ap05Equal(ab05ap05T, ab05ap05Output["context"], ab05ap05Expected["context"])
			ab05ap05Equal(ab05ap05T, ab05ap05Output["arguments"], ab05ap05Expected["arguments"])
		}
	}
	ab05ap05InputValue = ab05ap05Input()
	delete(ab05ap05InputValue["fields"].(map[string]any), "opaque")
	ab05ap05Output := ab05ap05Run(ab05ap05T, ab05ap05InputValue)
	ab05ap05Equal(ab05ap05T, ab05ap05Output["raw"], map[string]any{})
	ab05ap05InputValue["fields"].(map[string]any)["other"] = []any{0}
	ab05ap05Output = ab05ap05Run(ab05ap05T, ab05ap05InputValue)
	ab05ap05Equal(ab05ap05T, ab05ap05Output["raw"].(map[string]any)["other"], []any{float64(0)})
}

func TestAb05Ap05CounterAndPreviousBoundaries(ab05ap05T *ab05ap05testing.T) {
	for _, ab05ap05Case := range []struct {
		ab05ap05Prior     uint32
		ab05ap05Increment uint32
		ab05ap05Expected  float64
	}{
		{0, 0, 0}, {0, 4294967295, 4294967295}, {4294967295, 0, 4294967295}, {4294967295, 1, 0}, {4294967295, 4294967295, 4294967294},
	} {
		ab05ap05InputValue := ab05ap05Input()
		ab05ap05Args(ab05ap05InputValue)["failure_counter"] = ab05ap05Case.ab05ap05Prior
		ab05ap05Args(ab05ap05InputValue)["counter_increment"] = ab05ap05Case.ab05ap05Increment
		ab05ap05Slot(ab05ap05InputValue, 0)["outcome"] = "parse-failure"
		delete(ab05ap05Slot(ab05ap05InputValue, 0), "value")
		ab05ap05Output := ab05ap05Run(ab05ap05T, ab05ap05InputValue)
		ab05ap05Equal(ab05ap05T, ab05ap05Number(ab05ap05Fields(ab05ap05Output)["failure_counter"]), ab05ap05Case.ab05ap05Expected)
		ab05ap05Equal(ab05ap05T, ab05ap05Fields(ab05ap05Output)["failure_counter"].(map[string]any)["native_increments"], "Unknown")
	}
	for _, ab05ap05Previous := range []uint32{0, 4294967295} {
		for _, ab05ap05Prior := range []uint32{0, 4294967295} {
			ab05ap05InputValue := ab05ap05Input()
			ab05ap05Slot(ab05ap05InputValue, 0)["prior_previous"] = ab05ap05Previous
			ab05ap05Slot(ab05ap05InputValue, 0)["prior"] = ab05ap05Prior
			ab05ap05Slot(ab05ap05InputValue, 0)["registered"] = false
			ab05ap05Output := ab05ap05Run(ab05ap05T, ab05ap05InputValue)
			ab05ap05Equal(ab05ap05T, ab05ap05Number(ab05ap05Outcome(ab05ap05Output, 0)["previous"]), float64(ab05ap05Previous))
			ab05ap05Slot(ab05ap05InputValue, 0)["registered"] = true
			ab05ap05Slot(ab05ap05InputValue, 0)["outcome"] = "read-failure"
			delete(ab05ap05Slot(ab05ap05InputValue, 0), "value")
			ab05ap05Output = ab05ap05Run(ab05ap05T, ab05ap05InputValue)
			ab05ap05Equal(ab05ap05T, ab05ap05Number(ab05ap05Outcome(ab05ap05Output, 0)["previous"]), float64(ab05ap05Prior))
		}
	}
	for _, ab05ap05Key := range []string{"primary", "secondary", "operation"} {
		ab05ap05InputValue := ab05ap05Input()
		ab05ap05InputValue["context"].(map[string]any)[ab05ap05Key] = 0
		ab05ap05Output := ab05ap05Run(ab05ap05T, ab05ap05InputValue)
		ab05ap05Equal(ab05ap05T, ab05ap05Output["state"], "Unknown")
		ab05ap05Equal(ab05ap05T, ab05ap05Output["context"].(map[string]any)[ab05ap05Key], float64(0))
	}
}
