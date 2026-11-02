package vaillant

import (
	ab11aq11bytes "bytes"
	ab11aq11json "encoding/json"
	ab11aq11os "os"
	ab11aq11testing "testing"
)

type ab11aq11Case struct {
	Ab11Aq11Incomplete []string                `json:"incomplete"`
	Ab11Aq11ID         string                  `json:"id"`
	Ab11Aq11Input      ab11aq11json.RawMessage `json:"input"`
	Ab11Aq11State      string                  `json:"state"`
	Ab11Aq11Eligible   bool                    `json:"eligible"`
	Ab11Aq11False      []string                `json:"false"`
}

func ab11aq11Cases(ab11aq11T *ab11aq11testing.T) []ab11aq11Case {
	ab11aq11T.Helper()
	ab11aq11Data, ab11aq11Err := ab11aq11os.ReadFile("testdata/ab11aq11_b517_local_guard.json")
	if ab11aq11Err != nil {
		ab11aq11T.Fatal(ab11aq11Err)
	}
	var ab11aq11Cases []ab11aq11Case
	if ab11aq11Err := ab11aq11json.Unmarshal(ab11aq11Data, &ab11aq11Cases); ab11aq11Err != nil {
		ab11aq11T.Fatal(ab11aq11Err)
	}
	if len(ab11aq11Cases) != 22 {
		ab11aq11T.Fatalf("fixture count: %d", len(ab11aq11Cases))
	}
	return ab11aq11Cases
}

func ab11aq11Object(ab11aq11T *ab11aq11testing.T, ab11aq11Data []byte) map[string]any {
	ab11aq11T.Helper()
	var ab11aq11Object map[string]any
	if ab11aq11Err := ab11aq11json.Unmarshal(ab11aq11Data, &ab11aq11Object); ab11aq11Err != nil {
		ab11aq11T.Fatal(ab11aq11Err)
	}
	return ab11aq11Object
}

func ab11aq11Evaluate(ab11aq11T *ab11aq11testing.T, ab11aq11Input any) map[string]any {
	ab11aq11T.Helper()
	ab11aq11Data, ab11aq11Err := ab11aq11json.Marshal(ab11aq11Input)
	if ab11aq11Err != nil {
		ab11aq11T.Fatal(ab11aq11Err)
	}
	ab11aq11Output, ab11aq11Err := Ab11Aq11Evaluate(ab11aq11Data)
	if ab11aq11Err != nil {
		ab11aq11T.Fatal(ab11aq11Err)
	}
	return ab11aq11Object(ab11aq11T, ab11aq11Output)
}

func TestAb11Aq11Fixture(ab11aq11T *ab11aq11testing.T) {
	for _, ab11aq11Case := range ab11aq11Cases(ab11aq11T) {
		ab11aq11T.Run(ab11aq11Case.Ab11Aq11ID, func(ab11aq11T *ab11aq11testing.T) {
			ab11aq11Input := ab11aq11Object(ab11aq11T, ab11aq11Case.Ab11Aq11Input)
			ab11aq11Output := ab11aq11Evaluate(ab11aq11T, ab11aq11Input)
			if ab11aq11Output["state"] != ab11aq11Case.Ab11Aq11State || ab11aq11Output["physical"] != "Unknown" || ab11aq11Output["qualified"] != false {
				ab11aq11T.Fatalf("envelope: %v", ab11aq11Output)
			}
			for ab11aq11Key, ab11aq11Value := range ab11aq11Input {
				ab11aq11Before, _ := ab11aq11json.Marshal(ab11aq11Value)
				ab11aq11After, _ := ab11aq11json.Marshal(ab11aq11Output[ab11aq11Key])
				if !ab11aq11bytes.Equal(ab11aq11Before, ab11aq11After) {
					ab11aq11T.Fatalf("raw retention %s", ab11aq11Key)
				}
			}
			ab11aq11Results := ab11aq11Output["results"].(map[string]any)
			if len(ab11aq11Results) != 8 {
				ab11aq11T.Fatalf("results: %v", ab11aq11Results)
			}
			for ab11aq11Name, ab11aq11Value := range ab11aq11Results {
				ab11aq11Result := ab11aq11Value.(map[string]any)
				ab11aq11ExpectedState := "complete"
				for _, ab11aq11Incomplete := range ab11aq11Case.Ab11Aq11Incomplete {
					if ab11aq11Name == ab11aq11Incomplete {
						ab11aq11ExpectedState = "incomplete"
					}
				}
				if ab11aq11Result["state"] != ab11aq11ExpectedState {
					ab11aq11T.Errorf("state %s: %v", ab11aq11Name, ab11aq11Result)
				}
				ab11aq11Expected := true
				if ab11aq11Name == "eligible" {
					ab11aq11Expected = ab11aq11Case.Ab11Aq11Eligible
				}
				for _, ab11aq11False := range ab11aq11Case.Ab11Aq11False {
					if ab11aq11Name == ab11aq11False {
						ab11aq11Expected = false
					}
				}
				if ab11aq11Result["flags"].(map[string]any)["value"] != ab11aq11Expected {
					ab11aq11T.Errorf("predicate %s: %v", ab11aq11Name, ab11aq11Result)
				}
				if _, ab11aq11OK := ab11aq11Result["raw"]; !ab11aq11OK {
					ab11aq11T.Error("missing raw")
				}
				if ab11aq11Result["grade"] != "Hypothesis" || (ab11aq11Result["state"] == "complete" && ab11aq11Result["calculation"] != "Proven") || (ab11aq11Result["state"] == "incomplete" && ab11aq11Result["calculation"] != nil) {
					ab11aq11T.Error("grades")
				}
			}
			if ab11aq11Case.Ab11Aq11State == "partial" {
				if ab11aq11Results["eligible"].(map[string]any)["state"] != "incomplete" || ab11aq11Results["ready_state"].(map[string]any)["state"] != "complete" {
					ab11aq11T.Error("partial sibling state")
				}
			}
		})
	}
}

func TestAb11Aq11Context(ab11aq11T *ab11aq11testing.T) {
	for ab11aq11Name, ab11aq11Value := range map[string]any{"primary": 180, "secondary": 24, "operation": 0, "direction": "remote", "view": "B518"} {
		ab11aq11T.Run(ab11aq11Name, func(ab11aq11T *ab11aq11testing.T) {
			ab11aq11Input := ab11aq11Object(ab11aq11T, ab11aq11Cases(ab11aq11T)[0].Ab11Aq11Input)
			ab11aq11Input["context"].(map[string]any)[ab11aq11Name] = ab11aq11Value
			ab11aq11Output := ab11aq11Evaluate(ab11aq11T, ab11aq11Input)
			if ab11aq11Output["state"] != "Unknown" || len(ab11aq11Output["results"].(map[string]any)) != 0 || ab11aq11Output["qualified"] != false {
				ab11aq11T.Fatalf("unsupported context: %v", ab11aq11Output)
			}
			ab11aq11Before, _ := ab11aq11json.Marshal(ab11aq11Input["context"])
			ab11aq11After, _ := ab11aq11json.Marshal(ab11aq11Output["context"])
			if !ab11aq11bytes.Equal(ab11aq11Before, ab11aq11After) {
				ab11aq11T.Error("context retention")
			}
		})
	}
	for _, ab11aq11Name := range []string{"family", "variant"} {
		ab11aq11Input := ab11aq11Object(ab11aq11T, ab11aq11Cases(ab11aq11T)[0].Ab11Aq11Input)
		delete(ab11aq11Input["context"].(map[string]any), ab11aq11Name)
		ab11aq11Output := ab11aq11Evaluate(ab11aq11T, ab11aq11Input)
		if _, ab11aq11Present := ab11aq11Output["context"].(map[string]any)[ab11aq11Name]; ab11aq11Present {
			ab11aq11T.Error("opaque absence lost")
		}
		if ab11aq11Output["state"] != "complete" || ab11aq11Output["qualified"] != false {
			ab11aq11T.Error("opaque absence qualified")
		}
		for _, ab11aq11Value := range []any{nil, false, 0, "B518", []any{}} {
			ab11aq11Input := ab11aq11Object(ab11aq11T, ab11aq11Cases(ab11aq11T)[0].Ab11Aq11Input)
			ab11aq11Input["context"].(map[string]any)[ab11aq11Name] = ab11aq11Value
			ab11aq11Output := ab11aq11Evaluate(ab11aq11T, ab11aq11Input)
			if ab11aq11Output["qualified"] != false || ab11aq11Output["physical"] != "Unknown" || ab11aq11Output["state"] != "complete" {
				ab11aq11T.Error("opaque identity")
			}
		}
	}
}

func TestAb11Aq11Reject(ab11aq11T *ab11aq11testing.T) {
	for _, ab11aq11Data := range []string{`{`, `[]`, `null`, `{} {}`, `{"x":1,"x":2}`, `{"x":{"a":1,"\u0061":2}}`, `{"x":[{"a":1,"a":2}]}`, `{"x":1e999}`, `{"x":NaN}`} {
		ab11aq11Output, ab11aq11Err := Ab11Aq11Evaluate([]byte(ab11aq11Data))
		if ab11aq11Err == nil || ab11aq11Output != nil {
			ab11aq11T.Errorf("accepted %s", ab11aq11Data)
		}
	}
	for _, ab11aq11Group := range []string{"context", "fields", "arguments"} {
		ab11aq11Names := map[string][]string{"context": {"primary", "secondary", "operation"}, "fields": {"request", "input"}, "arguments": {"ready_state", "final_option", "received_length"}}[ab11aq11Group]
		for _, ab11aq11Name := range ab11aq11Names {
			for _, ab11aq11Bad := range []any{nil, true, "3", -1, 1.5, 65536, 1.0000000001} {
				ab11aq11Input := ab11aq11Object(ab11aq11T, ab11aq11Cases(ab11aq11T)[0].Ab11Aq11Input)
				if ab11aq11Group == "fields" {
					ab11aq11Input[ab11aq11Group].(map[string]any)[ab11aq11Name] = []any{ab11aq11Bad}
				} else {
					ab11aq11Input[ab11aq11Group].(map[string]any)[ab11aq11Name] = ab11aq11Bad
				}
				ab11aq11Data, _ := ab11aq11json.Marshal(ab11aq11Input)
				ab11aq11Output, ab11aq11Err := Ab11Aq11Evaluate(ab11aq11Data)
				if ab11aq11Err == nil || ab11aq11Output != nil {
					ab11aq11T.Errorf("accepted %s.%s=%v", ab11aq11Group, ab11aq11Name, ab11aq11Bad)
				}
			}
		}
	}
	for _, ab11aq11Name := range []string{"ready_state", "final_option", "received_length"} {
		ab11aq11Input := ab11aq11Object(ab11aq11T, ab11aq11Cases(ab11aq11T)[0].Ab11Aq11Input)
		delete(ab11aq11Input["arguments"].(map[string]any), ab11aq11Name)
		ab11aq11Data, _ := ab11aq11json.Marshal(ab11aq11Input)
		if ab11aq11Output, ab11aq11Err := Ab11Aq11Evaluate(ab11aq11Data); ab11aq11Err == nil || ab11aq11Output != nil {
			ab11aq11T.Errorf("missing argument %s", ab11aq11Name)
		}
	}
}

func TestAb11Aq11Memory(ab11aq11T *ab11aq11testing.T) {
	ab11aq11Input := append([]byte(nil), ab11aq11Cases(ab11aq11T)[0].Ab11Aq11Input...)
	ab11aq11Original := append([]byte(nil), ab11aq11Input...)
	ab11aq11First, ab11aq11Err := Ab11Aq11Evaluate(ab11aq11Input)
	if ab11aq11Err != nil {
		ab11aq11T.Fatal(ab11aq11Err)
	}
	ab11aq11Saved := append([]byte(nil), ab11aq11First...)
	for ab11aq11Index := range ab11aq11Input {
		ab11aq11Input[ab11aq11Index] = 'x'
	}
	if !ab11aq11bytes.Equal(ab11aq11First, ab11aq11Saved) {
		ab11aq11T.Fatal("input aliases output")
	}
	for ab11aq11Index := range ab11aq11First {
		ab11aq11First[ab11aq11Index] = 'x'
	}
	ab11aq11Second, ab11aq11Err := Ab11Aq11Evaluate(ab11aq11Original)
	if ab11aq11Err != nil || !ab11aq11bytes.Equal(ab11aq11Second, ab11aq11Saved) {
		ab11aq11T.Fatal("cross-call storage")
	}
}

func TestAb11Aq11KindsAndBounds(ab11aq11T *ab11aq11testing.T) {
	for _, ab11aq11Group := range []string{"context", "fields", "arguments"} {
		for _, ab11aq11Bad := range []any{nil, false, 0, "", []any{}} {
			ab11aq11Input := ab11aq11Object(ab11aq11T, ab11aq11Cases(ab11aq11T)[0].Ab11Aq11Input)
			ab11aq11Input[ab11aq11Group] = ab11aq11Bad
			ab11aq11Data, _ := ab11aq11json.Marshal(ab11aq11Input)
			if ab11aq11Output, ab11aq11Err := Ab11Aq11Evaluate(ab11aq11Data); ab11aq11Err == nil || ab11aq11Output != nil {
				ab11aq11T.Errorf("accepted group %s: %v", ab11aq11Group, ab11aq11Bad)
			}
		}
	}
	for _, ab11aq11Name := range []string{"direction", "view"} {
		ab11aq11Input := ab11aq11Object(ab11aq11T, ab11aq11Cases(ab11aq11T)[0].Ab11Aq11Input)
		ab11aq11Input["context"].(map[string]any)[ab11aq11Name] = nil
		ab11aq11Data, _ := ab11aq11json.Marshal(ab11aq11Input)
		if ab11aq11Output, ab11aq11Err := Ab11Aq11Evaluate(ab11aq11Data); ab11aq11Err == nil || ab11aq11Output != nil {
			ab11aq11T.Errorf("accepted context kind %s", ab11aq11Name)
		}
	}
	for _, ab11aq11Name := range []string{"primary", "secondary", "operation"} {
		ab11aq11Input := ab11aq11Object(ab11aq11T, ab11aq11Cases(ab11aq11T)[0].Ab11Aq11Input)
		ab11aq11Input["context"].(map[string]any)[ab11aq11Name] = 256
		ab11aq11Data, _ := ab11aq11json.Marshal(ab11aq11Input)
		if ab11aq11Output, ab11aq11Err := Ab11Aq11Evaluate(ab11aq11Data); ab11aq11Err == nil || ab11aq11Output != nil {
			ab11aq11T.Errorf("accepted context overflow %s", ab11aq11Name)
		}
	}
	for _, ab11aq11Name := range []string{"ready_state", "final_option"} {
		ab11aq11Input := ab11aq11Object(ab11aq11T, ab11aq11Cases(ab11aq11T)[0].Ab11Aq11Input)
		ab11aq11Input["arguments"].(map[string]any)[ab11aq11Name] = 256
		ab11aq11Data, _ := ab11aq11json.Marshal(ab11aq11Input)
		if ab11aq11Output, ab11aq11Err := Ab11Aq11Evaluate(ab11aq11Data); ab11aq11Err == nil || ab11aq11Output != nil {
			ab11aq11T.Errorf("accepted u8 overflow %s", ab11aq11Name)
		}
	}
	for _, ab11aq11Name := range []string{"request", "input"} {
		for _, ab11aq11Bad := range []any{nil, false, "", 0, []any{256}} {
			ab11aq11Input := ab11aq11Object(ab11aq11T, ab11aq11Cases(ab11aq11T)[0].Ab11Aq11Input)
			ab11aq11Input["fields"].(map[string]any)[ab11aq11Name] = ab11aq11Bad
			ab11aq11Data, _ := ab11aq11json.Marshal(ab11aq11Input)
			if ab11aq11Output, ab11aq11Err := Ab11Aq11Evaluate(ab11aq11Data); ab11aq11Err == nil || ab11aq11Output != nil {
				ab11aq11T.Errorf("accepted field kind/width %s", ab11aq11Name)
			}
		}
	}
	for _, ab11aq11Value := range []int{0, 65535} {
		ab11aq11Input := ab11aq11Object(ab11aq11T, ab11aq11Cases(ab11aq11T)[0].Ab11Aq11Input)
		ab11aq11Input["arguments"].(map[string]any)["received_length"] = ab11aq11Value
		ab11aq11Output := ab11aq11Evaluate(ab11aq11T, ab11aq11Input)
		if ab11aq11Output["results"].(map[string]any)["eligible"].(map[string]any)["flags"].(map[string]any)["value"] != false {
			ab11aq11T.Error("bounded length accepted as eligible")
		}
	}
	ab11aq11Input := ab11aq11Object(ab11aq11T, ab11aq11Cases(ab11aq11T)[0].Ab11Aq11Input)
	ab11aq11Input["arguments"].(map[string]any)["opaque"] = map[string]any{"null": nil, "zero": 0, "false": false, "empty": []any{}}
	ab11aq11Output := ab11aq11Evaluate(ab11aq11T, ab11aq11Input)
	ab11aq11Before, _ := ab11aq11json.Marshal(ab11aq11Input["arguments"])
	ab11aq11After, _ := ab11aq11json.Marshal(ab11aq11Output["arguments"])
	if !ab11aq11bytes.Equal(ab11aq11Before, ab11aq11After) {
		ab11aq11T.Error("argument extras retention")
	}
	for _, ab11aq11Key := range []string{"state", "results", "physical", "qualified", "raw"} {
		ab11aq11Input := ab11aq11Object(ab11aq11T, ab11aq11Cases(ab11aq11T)[0].Ab11Aq11Input)
		ab11aq11Input[ab11aq11Key] = map[string]any{"opaque": nil}
		ab11aq11Output := ab11aq11Evaluate(ab11aq11T, ab11aq11Input)
		ab11aq11Before, _ := ab11aq11json.Marshal(ab11aq11Input)
		ab11aq11After, _ := ab11aq11json.Marshal(ab11aq11Output["raw"])
		if !ab11aq11bytes.Equal(ab11aq11Before, ab11aq11After) {
			ab11aq11T.Errorf("reserved output key loses raw input: %s", ab11aq11Key)
		}
	}
	for _, ab11aq11Suffix := range []string{" {}", " true", " null", " []"} {
		ab11aq11Data := append(append([]byte(nil), ab11aq11Cases(ab11aq11T)[0].Ab11Aq11Input...), []byte(ab11aq11Suffix)...)
		if ab11aq11Output, ab11aq11Err := Ab11Aq11Evaluate(ab11aq11Data); ab11aq11Err == nil || ab11aq11Output != nil {
			ab11aq11T.Error("accepted trailing JSON")
		}
	}
}
