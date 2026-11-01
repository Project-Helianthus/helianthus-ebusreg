package vaillant

import (
	ab09ao09JSON "encoding/json"
	ab09ao09OS "os"
	ab09ao09Reflect "reflect"
	ab09ao09Testing "testing"
)

func ab09ao09Input() map[string]any {
	return map[string]any{"context": map[string]any{"primary": 181, "secondary": 15, "operation": 0, "direction": "positive", "view": "thermal-fields", "family": map[string]any{"opaque": []any{"synthetic"}}, "variant": []any{0, "candidate"}}, "fields": map[string]any{"temperature_one": []int{240, 255}, "temperature_two": []int{0, 0}, "temperature_three": []int{16, 0}, "error_one": []int{0}, "error_two": []int{1}, "error_three": []int{2}}}
}
func ab09ao09Call(ab09ao09T *ab09ao09Testing.T, ab09ao09InputValue any) map[string]any {
	ab09ao09T.Helper()
	ab09ao09Data, ab09ao09Err := ab09ao09JSON.Marshal(ab09ao09InputValue)
	if ab09ao09Err != nil {
		ab09ao09T.Fatal(ab09ao09Err)
	}
	ab09ao09Data, ab09ao09Err = Ab09Ao09Interpret(ab09ao09Data)
	if ab09ao09Err != nil {
		ab09ao09T.Fatal(ab09ao09Err)
	}
	var ab09ao09Output map[string]any
	if ab09ao09Err = ab09ao09JSON.Unmarshal(ab09ao09Data, &ab09ao09Output); ab09ao09Err != nil {
		ab09ao09T.Fatal(ab09ao09Err)
	}
	if ab09ao09Output["physical"] != "Unknown" || ab09ao09Output["qualified"] != false {
		ab09ao09T.Fatal("qualification promoted")
	}
	return ab09ao09Output
}
func ab09ao09Result(ab09ao09Output map[string]any, ab09ao09Key string) map[string]any {
	return ab09ao09Output["fields"].(map[string]any)[ab09ao09Key].(map[string]any)
}
func TestAb09Ao09Thermal(ab09ao09T *ab09ao09Testing.T) {
	for ab09ao09Index, ab09ao09Key := range []string{"temperature_one", "temperature_two", "temperature_three", "error_one", "error_two", "error_three"} {
		for _, ab09ao09Bytes := range [][]int{nil, {}, {0}, {16, 0}, {255, 255, 99}} {
			ab09ao09InputValue := ab09ao09Input()
			ab09ao09Raw := ab09ao09InputValue["fields"].(map[string]any)
			if ab09ao09Bytes == nil {
				delete(ab09ao09Raw, ab09ao09Key)
			} else {
				ab09ao09Raw[ab09ao09Key] = ab09ao09Bytes
			}
			ab09ao09Output := ab09ao09Call(ab09ao09T, ab09ao09InputValue)
			ab09ao09Field := ab09ao09Result(ab09ao09Output, ab09ao09Key)
			ab09ao09Width := 1
			if ab09ao09Index < 3 {
				ab09ao09Width = 2
			}
			if len(ab09ao09Bytes) != ab09ao09Width {
				if _, ab09ao09Found := ab09ao09Field["number"]; ab09ao09Found {
					ab09ao09T.Fatal("numeric incomplete field")
				}
				if ab09ao09Output["state"] != "partial" {
					ab09ao09T.Fatal("lost partial")
				}
			}
			if ab09ao09Bytes == nil {
				if _, ab09ao09Found := ab09ao09Field["raw"]; ab09ao09Found || ab09ao09Field["state"] != "missing" {
					ab09ao09T.Fatal("lost absence")
				}
			} else {
				ab09ao09Expected, _ := ab09ao09JSON.Marshal(ab09ao09Bytes)
				ab09ao09Actual, _ := ab09ao09JSON.Marshal(ab09ao09Field["raw"])
				if string(ab09ao09Expected) != string(ab09ao09Actual) {
					ab09ao09T.Fatal("raw alias", ab09ao09Key)
				}
			}
			for ab09ao09Sibling, ab09ao09Expected := range map[string]float64{"temperature_one": -1, "temperature_two": 0, "temperature_three": 1} {
				if ab09ao09Sibling != ab09ao09Key {
					ab09ao09Field = ab09ao09Result(ab09ao09Output, ab09ao09Sibling)
					if ab09ao09Field["number"] != ab09ao09Expected || ab09ao09Field["grade"] != "Hypothesis" || ab09ao09Field["calculation"] != "Proven" {
						ab09ao09T.Fatal("independent arithmetic", ab09ao09Sibling, ab09ao09Field)
					}
				}
			}
		}
	}
	ab09ao09InputValue := ab09ao09Input()
	ab09ao09InputValue["context"].(map[string]any)["operation"] = 1
	ab09ao09Raw := ab09ao09InputValue["fields"].(map[string]any)
	delete(ab09ao09Raw, "temperature_three")
	for _, ab09ao09Bytes := range [][]int{nil, {16, 0}, {0, 0}, {}, {1}, {1, 2, 3}} {
		if ab09ao09Bytes == nil {
			delete(ab09ao09Raw, "stack_voltage")
		} else {
			ab09ao09Raw["stack_voltage"] = ab09ao09Bytes
		}
		ab09ao09Output := ab09ao09Call(ab09ao09T, ab09ao09InputValue)
		ab09ao09Field := ab09ao09Result(ab09ao09Output, "stack_voltage")
		if _, ab09ao09Found := ab09ao09Field["number"]; ab09ao09Found {
			ab09ao09T.Fatal("invented voltage")
		}
		if len(ab09ao09Bytes) == 2 && ab09ao09Field["state"] != "unit-conflict" {
			ab09ao09T.Fatal("lost conflict")
		}
		if ab09ao09Result(ab09ao09Output, "temperature_one")["number"] != float64(-1) {
			ab09ao09T.Fatal("lost sibling")
		}
	}
}
func TestAb09Ao09DemandIntegrity(ab09ao09T *ab09ao09Testing.T) {
	for _, ab09ao09Demand := range []int{0, 1, 2, 4, 8, 16, 31, 32, 64, 128, 224, 255} {
		for ab09ao09Feedback := 0; ab09ao09Feedback < 256; ab09ao09Feedback++ {
			for _, ab09ao09Integrity := range [][]int{nil, {}, {18}, {52}, {18, 52}} {
				ab09ao09InputValue := ab09ao09Input()
				ab09ao09Context := ab09ao09InputValue["context"].(map[string]any)
				ab09ao09Context["operation"] = 3
				ab09ao09Context["view"] = "demand"
				ab09ao09Raw := map[string]any{"demand": []int{ab09ao09Demand}, "feedback": []int{ab09ao09Feedback}, "integrity_two": []int{52}}
				if ab09ao09Integrity != nil {
					ab09ao09Raw["integrity_one"] = ab09ao09Integrity
				}
				ab09ao09InputValue["fields"] = ab09ao09Raw
				ab09ao09Output := ab09ao09Call(ab09ao09T, ab09ao09InputValue)
				ab09ao09Field := ab09ao09Result(ab09ao09Output, "demand")
				if ab09ao09Field["label"] != "demand" || ab09ao09Field["reserved"] != float64(ab09ao09Demand&224) || ab09ao09Field["grade"] != "Hypothesis" {
					ab09ao09T.Fatal("demand promoted or reserved lost")
				}
				for ab09ao09Index, ab09ao09Name := range []string{"burner", "reformer-gas-valve", "gas-flow", "fan-power", "valve"} {
					if ab09ao09Field["flags"].(map[string]any)[ab09ao09Name] != (ab09ao09Demand&(1<<ab09ao09Index) != 0) {
						ab09ao09T.Fatal("demand bit", ab09ao09Name)
					}
				}
				ab09ao09Label := "rejected"
				if ab09ao09Feedback == 1 {
					ab09ao09Label = "accepted"
				}
				if ab09ao09Result(ab09ao09Output, "feedback")["label"] != ab09ao09Label {
					ab09ao09T.Fatal("feedback depends on integrity")
				}
				ab09ao09Field = ab09ao09Result(ab09ao09Output, "integrity_two")
				if ab09ao09Field["crc"] != "Unknown" || ab09ao09Field["word_order"] != "Unknown" || ab09ao09Field["logical_position"] != float64(5) {
					ab09ao09T.Fatal("invented CRC")
				}
				if _, ab09ao09Found := ab09ao09Field["number"]; ab09ao09Found {
					ab09ao09T.Fatal("numeric integrity")
				}
			}
		}
	}
	for _, ab09ao09Operation := range []int{0, 1, 2, 3} {
		ab09ao09InputValue := ab09ao09Input()
		ab09ao09Context := ab09ao09InputValue["context"].(map[string]any)
		ab09ao09Context["operation"] = ab09ao09Operation
		if ab09ao09Operation == 2 {
			ab09ao09Context["view"] = "integrity-only"
		}
		if ab09ao09Operation == 3 {
			ab09ao09Context["view"] = "demand"
		}
		ab09ao09Raw := ab09ao09InputValue["fields"].(map[string]any)
		ab09ao09Raw["integrity_one"] = []int{18}
		ab09ao09Raw["integrity_two"] = []int{52}
		ab09ao09Output := ab09ao09Call(ab09ao09T, ab09ao09InputValue)
		ab09ao09Position := 12
		if ab09ao09Operation == 2 {
			ab09ao09Position = 11
		}
		if ab09ao09Operation == 3 {
			ab09ao09Position = 4
		}
		for ab09ao09Index, ab09ao09Key := range []string{"integrity_one", "integrity_two"} {
			ab09ao09Field := ab09ao09Result(ab09ao09Output, ab09ao09Key)
			if ab09ao09Field["logical_position"] != float64(ab09ao09Position+ab09ao09Index) || ab09ao09Field["state"] != "unknown-order" {
				ab09ao09T.Fatal("integrity position")
			}
		}
		if ab09ao09Operation == 2 && ab09ao09Result(ab09ao09Output, "temperature_one")["state"] != "Unknown" {
			ab09ao09T.Fatal("borrowed thermal")
		}
	}
}
func TestAb09Ao09SelectorsCopies(ab09ao09T *ab09ao09Testing.T) {
	for ab09ao09Key, ab09ao09Value := range map[string]any{"primary": 180, "secondary": 14, "operation": 9, "direction": "negative", "view": "unknown", "family": map[string]any{"other": []int{1, 2}}, "variant": false} {
		ab09ao09InputValue := ab09ao09Input()
		ab09ao09InputValue["context"].(map[string]any)[ab09ao09Key] = ab09ao09Value
		ab09ao09Output := ab09ao09Call(ab09ao09T, ab09ao09InputValue)
		ab09ao09Expected, _ := ab09ao09JSON.Marshal(ab09ao09InputValue["context"])
		ab09ao09Actual, _ := ab09ao09JSON.Marshal(ab09ao09Output["context"])
		if string(ab09ao09Expected) != string(ab09ao09Actual) {
			ab09ao09T.Fatal("context changed")
		}
		if ab09ao09Key != "family" && ab09ao09Key != "variant" {
			if ab09ao09Output["state"] != "Unknown" || ab09ao09Result(ab09ao09Output, "temperature_one")["state"] != "Unknown" {
				ab09ao09T.Fatal("guessed selection", ab09ao09Key)
			}
		} else if ab09ao09Result(ab09ao09Output, "temperature_one")["number"] != float64(-1) {
			ab09ao09T.Fatal("family selected view")
		}
	}
	for _, ab09ao09Key := range []string{"primary", "secondary", "operation", "direction", "view"} {
		ab09ao09InputValue := ab09ao09Input()
		delete(ab09ao09InputValue["context"].(map[string]any), ab09ao09Key)
		if ab09ao09Call(ab09ao09T, ab09ao09InputValue)["state"] != "Unknown" {
			ab09ao09T.Fatal("missing selector", ab09ao09Key)
		}
	}
	ab09ao09InputValue := ab09ao09Input()
	ab09ao09InputValue["arguments"] = map[string]any{}
	ab09ao09InputValue["fields"].(map[string]any)["tail"] = []int{1, 2, 3}
	ab09ao09Data, _ := ab09ao09JSON.Marshal(ab09ao09InputValue)
	ab09ao09Saved := append([]byte(nil), ab09ao09Data...)
	ab09ao09First, ab09ao09Err := Ab09Ao09Interpret(ab09ao09Data)
	if ab09ao09Err != nil {
		ab09ao09T.Fatal(ab09ao09Err)
	}
	for ab09ao09Index := range ab09ao09Data {
		ab09ao09Data[ab09ao09Index] = ' '
	}
	ab09ao09Second, ab09ao09Err := Ab09Ao09Interpret(ab09ao09Saved)
	if ab09ao09Err != nil || string(ab09ao09First) != string(ab09ao09Second) {
		ab09ao09T.Fatal("input alias")
	}
	ab09ao09First[0] = '!'
	ab09ao09Third, _ := Ab09Ao09Interpret(ab09ao09Saved)
	if string(ab09ao09Third) != string(ab09ao09Second) {
		ab09ao09T.Fatal("output alias")
	}
	ab09ao09Output := ab09ao09Call(ab09ao09T, ab09ao09InputValue)
	ab09ao09Result(ab09ao09Output, "tail")["raw"].([]any)[0] = float64(99)
	if ab09ao09Output["raw"].(map[string]any)["tail"].([]any)[0] != float64(1) {
		ab09ao09T.Fatal("nested alias")
	}
	if _, ab09ao09Found := ab09ao09Output["arguments"]; !ab09ao09Found {
		ab09ao09T.Fatal("empty arguments lost")
	}
	delete(ab09ao09InputValue, "arguments")
	if _, ab09ao09Found := ab09ao09Call(ab09ao09T, ab09ao09InputValue)["arguments"]; ab09ao09Found {
		ab09ao09T.Fatal("arguments invented")
	}
}
func TestAb09Ao09Reject(ab09ao09T *ab09ao09Testing.T) {
	for _, ab09ao09Text := range []string{`null`, `[]`, `{}`, `{"context":{},"fields":null}`, `{"context":{},"fields":{"x":null}}`, `{"context":{},"fields":{"x":[256]}}`, `{"context":{},"fields":{"x":[-1]}}`, `{"context":{},"fields":{"x":[1.0]}}`, `{"context":{},"fields":{"x":[true]}}`, `{"context":{},"fields":{"x":["1"]}}`, `{"context":{},"fields":{"x":[{}]}}`, `{"context":{},"fields":{"x":{}}}`, `{"context":{},"fields":{},"arguments":null}`, `{"context":{},"fields":{},"arguments":{"flag":true}}`, `{"context":{"primary":"181"},"fields":{}}`, `{"context":{"operation":0.1},"fields":{}}`, `{"context":{"direction":0},"fields":{}}`, `{"context":{"view":null},"fields":{}}`, `{"context":{},"fields":{},"extra":0}`, `{"context":{},"fields":{},"fields":{}}`, `{"context":{"family":{"x":1,"x":2}},"fields":{}}`, `{"context":{},"fields":{"x":[],"x":[]}}`, `{"context":{},"fields":{}} {}`, `{"context":{},"fields":{}`, `{"context":{},"fields":{"x":[1e999]}}`} {
		ab09ao09Output, ab09ao09Err := Ab09Ao09Interpret([]byte(ab09ao09Text))
		if ab09ao09Err == nil || ab09ao09Output != nil {
			ab09ao09T.Fatal("accepted invalid input", ab09ao09Text)
		}
	}
}
func TestAb09Ao09Fixtures(ab09ao09T *ab09ao09Testing.T) {
	ab09ao09Data, ab09ao09Err := ab09ao09OS.ReadFile("testdata/ab09ao09_b50f_field_views.json")
	if ab09ao09Err != nil {
		ab09ao09T.Fatal(ab09ao09Err)
	}
	var ab09ao09Cases []struct {
		Ab09Ao09ID       string                  `json:"ab09ao09_id"`
		Ab09Ao09Input    ab09ao09JSON.RawMessage `json:"input"`
		Ab09Ao09Expected map[string]any          `json:"expected"`
	}
	if ab09ao09Err = ab09ao09JSON.Unmarshal(ab09ao09Data, &ab09ao09Cases); ab09ao09Err != nil || len(ab09ao09Cases) == 0 {
		ab09ao09T.Fatal("invalid fixtures", ab09ao09Err)
	}
	for _, ab09ao09Case := range ab09ao09Cases {
		ab09ao09T.Run(ab09ao09Case.Ab09Ao09ID, func(ab09ao09T *ab09ao09Testing.T) {
			ab09ao09Data, ab09ao09Err := Ab09Ao09Interpret(ab09ao09Case.Ab09Ao09Input)
			if ab09ao09Err != nil {
				ab09ao09T.Fatal(ab09ao09Err)
			}
			var ab09ao09Output map[string]any
			if ab09ao09Err = ab09ao09JSON.Unmarshal(ab09ao09Data, &ab09ao09Output); ab09ao09Err != nil {
				ab09ao09T.Fatal(ab09ao09Err)
			}
			for ab09ao09Key, ab09ao09Expected := range ab09ao09Case.Ab09Ao09Expected {
				if !ab09ao09Reflect.DeepEqual(ab09ao09Output[ab09ao09Key], ab09ao09Expected) {
					ab09ao09T.Fatal(ab09ao09Key, ab09ao09Output)
				}
			}
		})
	}
}

func TestAb09Ao09EveryFieldAxis(ab09ao09T *ab09ao09Testing.T) {
	for _, ab09ao09Operation := range []int{0, 1, 2, 3} {
		ab09ao09Control := ab09ao09Input()
		ab09ao09Context := ab09ao09Control["context"].(map[string]any)
		ab09ao09Context["operation"] = ab09ao09Operation
		ab09ao09Raw := ab09ao09Control["fields"].(map[string]any)
		if ab09ao09Operation == 1 {
			delete(ab09ao09Raw, "temperature_three")
			ab09ao09Raw["stack_voltage"] = []int{16, 0}
		}
		if ab09ao09Operation == 2 {
			ab09ao09Context["view"] = "integrity-only"
			ab09ao09Raw = map[string]any{}
		}
		if ab09ao09Operation == 3 {
			ab09ao09Context["view"] = "demand"
			ab09ao09Raw = map[string]any{"demand": []int{31}, "feedback": []int{1}}
		}
		ab09ao09Raw["integrity_one"] = []int{18}
		ab09ao09Raw["integrity_two"] = []int{52}
		ab09ao09Control["fields"] = ab09ao09Raw
		ab09ao09Baseline := ab09ao09Call(ab09ao09T, ab09ao09Control)
		for ab09ao09Key, ab09ao09Original := range ab09ao09Raw {
			for _, ab09ao09Bytes := range [][]int{nil, {}, {0}, {255}, {0, 0}, {240, 255}, {1, 2, 3}} {
				if ab09ao09Bytes == nil {
					delete(ab09ao09Raw, ab09ao09Key)
				} else {
					ab09ao09Raw[ab09ao09Key] = ab09ao09Bytes
				}
				ab09ao09Output := ab09ao09Call(ab09ao09T, ab09ao09Control)
				for ab09ao09Sibling := range ab09ao09Raw {
					if ab09ao09Sibling != ab09ao09Key && !ab09ao09Reflect.DeepEqual(ab09ao09Result(ab09ao09Output, ab09ao09Sibling), ab09ao09Result(ab09ao09Baseline, ab09ao09Sibling)) {
						ab09ao09T.Fatal("sibling coupling", ab09ao09Operation, ab09ao09Key, ab09ao09Sibling)
					}
				}
				ab09ao09Width := len(ab09ao09Original.([]int))
				ab09ao09FieldValue, ab09ao09Present := ab09ao09Output["fields"].(map[string]any)[ab09ao09Key]
				if !ab09ao09Present {
					if ab09ao09Bytes != nil || (ab09ao09Key != "integrity_one" && ab09ao09Key != "integrity_two") {
						ab09ao09T.Fatal("lost field")
					}
					continue
				}
				ab09ao09Field := ab09ao09FieldValue.(map[string]any)
				if ab09ao09Bytes == nil {
					if _, ab09ao09Found := ab09ao09Field["raw"]; ab09ao09Found || ab09ao09Field["state"] != "missing" {
						ab09ao09T.Fatal("absence collapsed")
					}
				} else {
					ab09ao09Expected, _ := ab09ao09JSON.Marshal(ab09ao09Bytes)
					ab09ao09Actual, _ := ab09ao09JSON.Marshal(ab09ao09Field["raw"])
					if string(ab09ao09Actual) != string(ab09ao09Expected) {
						ab09ao09T.Fatal("raw substituted", ab09ao09Key)
					}
					if len(ab09ao09Bytes) != ab09ao09Width {
						ab09ao09State := "incomplete"
						if len(ab09ao09Bytes) > ab09ao09Width {
							ab09ao09State = "unmapped"
						}
						if ab09ao09Field["state"] != ab09ao09State {
							ab09ao09T.Fatal("width guessed", ab09ao09Key)
						}
						for _, ab09ao09Semantic := range []string{"number", "label", "flags", "reserved"} {
							if _, ab09ao09Found := ab09ao09Field[ab09ao09Semantic]; ab09ao09Found {
								ab09ao09T.Fatal("incomplete semantics", ab09ao09Key)
							}
						}
					}
				}
			}
			ab09ao09Raw[ab09ao09Key] = ab09ao09Original
		}
	}
	for _, ab09ao09Key := range []string{"primary", "secondary", "operation", "direction", "view"} {
		for _, ab09ao09Value := range []any{nil, true, []int{}, map[string]any{}} {
			ab09ao09InputValue := ab09ao09Input()
			ab09ao09InputValue["context"].(map[string]any)[ab09ao09Key] = ab09ao09Value
			ab09ao09Data, _ := ab09ao09JSON.Marshal(ab09ao09InputValue)
			ab09ao09Result, ab09ao09Err := Ab09Ao09Interpret(ab09ao09Data)
			if ab09ao09Err == nil || ab09ao09Result != nil {
				ab09ao09T.Fatal("wrong selector kind", ab09ao09Key)
			}
		}
	}
}
