package vaillant

import (
	ab21aq21json "encoding/json"
	ab21aq21os "os"
	ab21aq21reflect "reflect"
	ab21aq21strings "strings"
	ab21aq21testing "testing"
)

type ab21aq21Case struct {
	Ab21Aq21ID       string         `json:"id"`
	Ab21Aq21Input    map[string]any `json:"input"`
	Ab21Aq21Expected map[string]any `json:"expected"`
}

func ab21aq21Cases(ab21aq21T *ab21aq21testing.T) []ab21aq21Case {
	ab21aq21T.Helper()
	ab21aq21Data, ab21aq21Err := ab21aq21os.ReadFile("testdata/ab21aq21_b51a_diagnostic_fields.json")
	if ab21aq21Err != nil {
		ab21aq21T.Fatal(ab21aq21Err)
	}
	var ab21aq21Cases []ab21aq21Case
	if ab21aq21Err := ab21aq21json.Unmarshal(ab21aq21Data, &ab21aq21Cases); ab21aq21Err != nil {
		ab21aq21T.Fatal(ab21aq21Err)
	}
	if len(ab21aq21Cases) != 9 {
		ab21aq21T.Fatal("fixture count", len(ab21aq21Cases))
	}
	return ab21aq21Cases
}

func ab21aq21Eval(ab21aq21T *ab21aq21testing.T, ab21aq21Input map[string]any) map[string]any {
	ab21aq21T.Helper()
	ab21aq21Data, ab21aq21Err := ab21aq21json.Marshal(ab21aq21Input)
	if ab21aq21Err != nil {
		ab21aq21T.Fatal(ab21aq21Err)
	}
	ab21aq21Output, ab21aq21Err := Ab21Aq21Evaluate(ab21aq21Data)
	if ab21aq21Err != nil {
		ab21aq21T.Fatal(ab21aq21Err)
	}
	var ab21aq21Result map[string]any
	if ab21aq21Err := ab21aq21json.Unmarshal(ab21aq21Output, &ab21aq21Result); ab21aq21Err != nil {
		ab21aq21T.Fatal(ab21aq21Err)
	}
	if ab21aq21Result["qualified"] != false || ab21aq21Result["physical"] != "Unknown" {
		ab21aq21T.Fatal("qualification", ab21aq21Result)
	}
	if !ab21aq21reflect.DeepEqual(ab21aq21Input, ab21aq21Result["input"]) {
		ab21aq21T.Fatal("input preservation")
	}
	if !ab21aq21reflect.DeepEqual(ab21aq21Input["fields"], ab21aq21Result["raw"]) {
		ab21aq21T.Fatal("raw preservation")
	}
	return ab21aq21Result
}

func ab21aq21Subset(ab21aq21T *ab21aq21testing.T, ab21aq21Got, ab21aq21Want any) {
	ab21aq21T.Helper()
	if ab21aq21Map, ab21aq21OK := ab21aq21Want.(map[string]any); ab21aq21OK {
		ab21aq21Actual, ab21aq21OK := ab21aq21Got.(map[string]any)
		if !ab21aq21OK {
			ab21aq21T.Fatalf("want object, got %v", ab21aq21Got)
		}
		for ab21aq21Key, ab21aq21Value := range ab21aq21Map {
			ab21aq21Subset(ab21aq21T, ab21aq21Actual[ab21aq21Key], ab21aq21Value)
		}
		return
	}
	if !ab21aq21reflect.DeepEqual(ab21aq21Got, ab21aq21Want) {
		ab21aq21T.Fatalf("got %v want %v", ab21aq21Got, ab21aq21Want)
	}
}

func TestAb21Aq21Fixtures(ab21aq21T *ab21aq21testing.T) {
	for _, ab21aq21Case := range ab21aq21Cases(ab21aq21T) {
		ab21aq21T.Run(ab21aq21Case.Ab21Aq21ID, func(ab21aq21T *ab21aq21testing.T) {
			ab21aq21Output := ab21aq21Eval(ab21aq21T, ab21aq21Case.Ab21Aq21Input)
			if ab21aq21Output["state"] != "complete" {
				ab21aq21T.Fatal(ab21aq21Output)
			}
			ab21aq21Subset(ab21aq21T, ab21aq21Output["results"], ab21aq21Case.Ab21Aq21Expected)
			for ab21aq21Key, ab21aq21Value := range ab21aq21Output["results"].(map[string]any) {
				ab21aq21Result := ab21aq21Value.(map[string]any)
				if !ab21aq21reflect.DeepEqual(ab21aq21Result["raw"], ab21aq21Case.Ab21Aq21Input["fields"].(map[string]any)[ab21aq21Key]) {
					ab21aq21T.Fatal("result raw")
				}
				if ab21aq21Key == "alternatives" {
					if _, ab21aq21Present := ab21aq21Result["number"]; ab21aq21Present {
						ab21aq21T.Fatal("enum current number")
					}
				}
				if ab21aq21Result["state"] == "Unknown" {
					if _, ab21aq21Present := ab21aq21Result["number"]; ab21aq21Present {
						ab21aq21T.Fatal("disabled interpreted", ab21aq21Key)
					}
					continue
				}
				if ab21aq21Result["state"] != "Hypothesis" {
					ab21aq21T.Fatal("candidate grade", ab21aq21Result)
				}
				if _, ab21aq21Present := ab21aq21Result["number"]; ab21aq21Present && ab21aq21Result["calculation"] != "Proven" {
					ab21aq21T.Fatal("calculation grade")
				}
				for _, ab21aq21Forbidden := range []string{"current", "current_value", "success", "effective_authorization"} {
					if _, ab21aq21Present := ab21aq21Result[ab21aq21Forbidden]; ab21aq21Present {
						ab21aq21T.Fatal("invented result", ab21aq21Forbidden)
					}
				}
			}
		})
	}
}

func TestAb21Aq21EnablePositions(ab21aq21T *ab21aq21testing.T) {
	for ab21aq21Position := 0; ab21aq21Position < 104; ab21aq21Position++ {
		ab21aq21Input := ab21aq21Cases(ab21aq21T)[0].Ab21Aq21Input
		ab21aq21Input["fields"].(map[string]any)["enable"].([]any)[ab21aq21Position/8] = float64(uint64(1) << uint(ab21aq21Position%8))
		ab21aq21Bits := ab21aq21Eval(ab21aq21T, ab21aq21Input)["results"].(map[string]any)["enable"].(map[string]any)["flags"].([]any)
		if len(ab21aq21Bits) != 104 {
			ab21aq21T.Fatal("bit count")
		}
		for ab21aq21Index, ab21aq21Bit := range ab21aq21Bits {
			if ab21aq21Bit != (ab21aq21Index == ab21aq21Position) {
				ab21aq21T.Fatal("bit position", ab21aq21Position, ab21aq21Index)
			}
		}
	}
}

func TestAb21Aq21ContextAxes(ab21aq21T *ab21aq21testing.T) {
	for _, ab21aq21Key := range []string{"primary", "secondary", "operation", "direction", "view"} {
		for _, ab21aq21Absent := range []bool{false, true} {
			ab21aq21Input := ab21aq21Cases(ab21aq21T)[1].Ab21Aq21Input
			ab21aq21Context := ab21aq21Input["context"].(map[string]any)
			if ab21aq21Key == "direction" || ab21aq21Key == "view" {
				ab21aq21Context[ab21aq21Key] = "unsupported"
			} else {
				ab21aq21Context[ab21aq21Key] = float64(0)
			}
			if ab21aq21Absent {
				delete(ab21aq21Context, ab21aq21Key)
			}
			ab21aq21Output := ab21aq21Eval(ab21aq21T, ab21aq21Input)
			if ab21aq21Output["state"] != "Unknown" || len(ab21aq21Output["results"].(map[string]any)) != 0 {
				ab21aq21T.Fatal("context admitted", ab21aq21Key, ab21aq21Absent)
			}
		}
	}
	for _, ab21aq21Key := range []string{"family", "variant"} {
		for _, ab21aq21Value := range []any{nil, false, float64(0), "", "another-family", []any{}} {
			ab21aq21Input := ab21aq21Cases(ab21aq21T)[1].Ab21Aq21Input
			ab21aq21Input["context"].(map[string]any)[ab21aq21Key] = ab21aq21Value
			if ab21aq21Eval(ab21aq21T, ab21aq21Input)["state"] != "complete" {
				ab21aq21T.Fatal("opaque identity selects layout")
			}
		}
	}
}

func TestAb21Aq21FieldWidthsAndSiblings(ab21aq21T *ab21aq21testing.T) {
	for _, ab21aq21Case := range ab21aq21Cases(ab21aq21T) {
		for ab21aq21Key, ab21aq21Value := range ab21aq21Case.Ab21Aq21Input["fields"].(map[string]any) {
			for _, ab21aq21Mode := range []string{"absent", "empty", "short", "long"} {
				ab21aq21Input := ab21aq21Cases(ab21aq21T)[0].Ab21Aq21Input
				for _, ab21aq21Candidate := range ab21aq21Cases(ab21aq21T) {
					if ab21aq21Candidate.Ab21Aq21ID == ab21aq21Case.Ab21Aq21ID {
						ab21aq21Input = ab21aq21Candidate.Ab21Aq21Input
					}
				}
				ab21aq21Fields := ab21aq21Input["fields"].(map[string]any)
				ab21aq21Array := ab21aq21Value.([]any)
				switch ab21aq21Mode {
				case "absent":
					delete(ab21aq21Fields, ab21aq21Key)
				case "empty":
					ab21aq21Fields[ab21aq21Key] = []any{}
				case "short":
					ab21aq21Fields[ab21aq21Key] = ab21aq21Array[:len(ab21aq21Array)-1]
				case "long":
					ab21aq21Fields[ab21aq21Key] = append(append([]any{}, ab21aq21Array...), float64(0))
				}
				ab21aq21Output := ab21aq21Eval(ab21aq21T, ab21aq21Input)
				ab21aq21Result := ab21aq21Output["results"].(map[string]any)[ab21aq21Key].(map[string]any)
				if ab21aq21Output["state"] != "partial" || ab21aq21Result["state"] != "incomplete" {
					ab21aq21T.Fatal("width", ab21aq21Key, ab21aq21Mode)
				}
				if _, ab21aq21Present := ab21aq21Result["number"]; ab21aq21Present {
					ab21aq21T.Fatal("incomplete number")
				}
				if _, ab21aq21Present := ab21aq21Result["raw"]; ab21aq21Present != (ab21aq21Mode != "absent") {
					ab21aq21T.Fatal("absence conflated")
				}
				if ab21aq21Case.Ab21Aq21ID == "ab21aq21_range" && ab21aq21Key != "flags" && ab21aq21Key != "value" {
					ab21aq21Subset(ab21aq21T, ab21aq21Output["results"].(map[string]any)["value"], map[string]any{"number": float64(4660), "state": "Hypothesis"})
				}
			}
		}
	}
}

func TestAb21Aq21RejectJSONAndKinds(ab21aq21T *ab21aq21testing.T) {
	for _, ab21aq21Input := range []string{``, `{`, `[]`, `{} {}`, `{"context":{},"fields":{},"x":{"y":1,"y":2}}`, `{"context":{"primary":181,"primary":181},"fields":{}}`, `{"context":{},"fields":{"x":[{"a":0,"a":1}]}}`, `{"context":{},"fields":{},"x":1e999}`, `{"context":{},"fields":{},"x":NaN}`} {
		ab21aq21Output, ab21aq21Err := Ab21Aq21Evaluate([]byte(ab21aq21Input))
		if ab21aq21Err == nil || ab21aq21Output != nil {
			ab21aq21T.Fatal("accepted", ab21aq21Input)
		}
	}
	for _, ab21aq21Key := range []string{"primary", "secondary", "operation", "byte", "value_width", "context", "fields", "arguments", "direction", "view"} {
		for _, ab21aq21Bad := range []any{nil, true, "1", float64(-1), float64(256), 1.5, []any{}} {
			if (ab21aq21Key == "direction" || ab21aq21Key == "view") && ab21aq21Bad == "1" {
				continue
			}
			ab21aq21Input := ab21aq21Cases(ab21aq21T)[5].Ab21Aq21Input
			switch ab21aq21Key {
			case "byte":
				ab21aq21Input["fields"].(map[string]any)["value"] = []any{ab21aq21Bad, 0}
			case "value_width":
				ab21aq21Input["arguments"].(map[string]any)[ab21aq21Key] = ab21aq21Bad
			case "context", "fields", "arguments":
				ab21aq21Input[ab21aq21Key] = ab21aq21Bad
			default:
				ab21aq21Input["context"].(map[string]any)[ab21aq21Key] = ab21aq21Bad
			}
			ab21aq21Data, _ := ab21aq21json.Marshal(ab21aq21Input)
			ab21aq21Output, ab21aq21Err := Ab21Aq21Evaluate(ab21aq21Data)
			if ab21aq21Err == nil || ab21aq21Output != nil {
				ab21aq21T.Fatal("accepted kind", ab21aq21Key, ab21aq21Bad)
			}
		}
	}
	for _, ab21aq21Width := range []any{nil, float64(0), float64(3)} {
		ab21aq21Input := ab21aq21Cases(ab21aq21T)[5].Ab21Aq21Input
		if ab21aq21Width == nil {
			delete(ab21aq21Input, "arguments")
		} else {
			ab21aq21Input["arguments"].(map[string]any)["value_width"] = ab21aq21Width
		}
		ab21aq21Data, _ := ab21aq21json.Marshal(ab21aq21Input)
		if ab21aq21Output, ab21aq21Err := Ab21Aq21Evaluate(ab21aq21Data); ab21aq21Err == nil || ab21aq21Output != nil {
			ab21aq21T.Fatal("inferred width")
		}
	}
}

func TestAb21Aq21PreservationAndMemory(ab21aq21T *ab21aq21testing.T) {
	ab21aq21Input := ab21aq21Cases(ab21aq21T)[1].Ab21Aq21Input
	for _, ab21aq21Key := range []string{"state", "results", "physical", "qualified", "input", "raw", "extra"} {
		ab21aq21Input[ab21aq21Key] = map[string]any{"null": nil, "zero": float64(0), "false": false, "empty": []any{}}
	}
	ab21aq21Input["fields"].(map[string]any)["opaque"] = []any{}
	ab21aq21Output := ab21aq21Eval(ab21aq21T, ab21aq21Input)
	if !ab21aq21reflect.DeepEqual(ab21aq21Output["extra"], ab21aq21Input["extra"]) {
		ab21aq21T.Fatal("extra lost")
	}
	ab21aq21Data, _ := ab21aq21json.Marshal(ab21aq21Input)
	ab21aq21First, ab21aq21Err := Ab21Aq21Evaluate(ab21aq21Data)
	if ab21aq21Err != nil {
		ab21aq21T.Fatal(ab21aq21Err)
	}
	ab21aq21Saved := string(ab21aq21First)
	for ab21aq21Index := range ab21aq21Data {
		ab21aq21Data[ab21aq21Index] = ' '
	}
	if string(ab21aq21First) != ab21aq21Saved {
		ab21aq21T.Fatal("input aliases output")
	}
	for ab21aq21Index := range ab21aq21First {
		ab21aq21First[ab21aq21Index] = ' '
	}
	ab21aq21Next, _ := ab21aq21json.Marshal(ab21aq21Eval(ab21aq21T, ab21aq21Input))
	var ab21aq21Expected any
	if ab21aq21Err := ab21aq21json.Unmarshal([]byte(ab21aq21Saved), &ab21aq21Expected); ab21aq21Err != nil {
		ab21aq21T.Fatal(ab21aq21Err)
	}
	var ab21aq21Actual any
	if ab21aq21Err := ab21aq21json.Unmarshal(ab21aq21Next, &ab21aq21Actual); ab21aq21Err != nil {
		ab21aq21T.Fatal(ab21aq21Err)
	}
	if !ab21aq21reflect.DeepEqual(ab21aq21Actual, ab21aq21Expected) {
		ab21aq21T.Fatal("cross-call storage")
	}
	if ab21aq21strings.Contains(string(ab21aq21Next), `"success"`) {
		ab21aq21T.Fatal("success inferred")
	}
}

func TestAb21Aq21IndependentFieldValues(ab21aq21T *ab21aq21testing.T) {
	for _, ab21aq21Case := range ab21aq21Cases(ab21aq21T) {
		if ab21aq21Case.Ab21Aq21ID == "ab21aq21_disabled" {
			continue
		}
		for ab21aq21Key, ab21aq21Value := range ab21aq21Case.Ab21Aq21Input["fields"].(map[string]any) {
			if ab21aq21Key == "enable" || ab21aq21Key == "flags" {
				continue
			}
			for ab21aq21Position := range ab21aq21Value.([]any) {
				var ab21aq21Input map[string]any
				for _, ab21aq21Candidate := range ab21aq21Cases(ab21aq21T) {
					if ab21aq21Candidate.Ab21Aq21ID == ab21aq21Case.Ab21Aq21ID {
						ab21aq21Input = ab21aq21Candidate.Ab21Aq21Input
					}
				}
				ab21aq21Array := make([]any, len(ab21aq21Value.([]any)))
				for ab21aq21Index := range ab21aq21Array {
					ab21aq21Array[ab21aq21Index] = float64(0)
				}
				ab21aq21Array[ab21aq21Position] = float64(1)
				ab21aq21Input["fields"].(map[string]any)[ab21aq21Key] = ab21aq21Array
				ab21aq21Result := ab21aq21Eval(ab21aq21T, ab21aq21Input)["results"].(map[string]any)[ab21aq21Key].(map[string]any)
				switch ab21aq21Key {
				case "alternatives":
					ab21aq21Subset(ab21aq21T, ab21aq21Result["alternatives"], ab21aq21Array)
				case "reserved":
					ab21aq21Subset(ab21aq21T, ab21aq21Result["reserved"], ab21aq21Array)
				default:
					ab21aq21Subset(ab21aq21T, ab21aq21Result["number"], []float64{1, 256, 65536, 16777216}[ab21aq21Position])
				}
			}
		}
	}
	for _, ab21aq21Index := range []int{1, 3, 4} {
		ab21aq21Input := ab21aq21Cases(ab21aq21T)[ab21aq21Index].Ab21Aq21Input
		ab21aq21Input["fields"].(map[string]any)["flags"] = []any{float64(1)}
		for ab21aq21Key, ab21aq21Value := range ab21aq21Eval(ab21aq21T, ab21aq21Input)["results"].(map[string]any) {
			if ab21aq21Key == "flags" {
				continue
			}
			ab21aq21Result := ab21aq21Value.(map[string]any)
			if ab21aq21Result["state"] != "Unknown" || len(ab21aq21Result) != 2 {
				ab21aq21T.Fatal("disabled interpretation", ab21aq21Key)
			}
		}
	}
	for ab21aq21Status, ab21aq21Label := range map[int]string{3: "range-violation", 4: "misconfiguration", 255: "unmapped"} {
		ab21aq21Input := ab21aq21Cases(ab21aq21T)[5].Ab21Aq21Input
		ab21aq21Input["fields"].(map[string]any)["status"] = []any{float64(ab21aq21Status)}
		ab21aq21Subset(ab21aq21T, ab21aq21Eval(ab21aq21T, ab21aq21Input)["results"].(map[string]any)["status"], map[string]any{"label": ab21aq21Label})
	}
	for ab21aq21Bit, ab21aq21Name := range []string{"disabled", "writable", "enum-versus-range", "additional-four-byte/no-range", "reserved", "reserved", "text", "text"} {
		ab21aq21Input := ab21aq21Cases(ab21aq21T)[1].Ab21Aq21Input
		ab21aq21Input["fields"].(map[string]any)["flags"] = []any{float64(uint64(1) << ab21aq21Bit)}
		ab21aq21Result := ab21aq21Eval(ab21aq21T, ab21aq21Input)["results"].(map[string]any)["flags"].(map[string]any)
		if ab21aq21Bit < 4 {
			for ab21aq21Key, ab21aq21Value := range ab21aq21Result["flags"].(map[string]any) {
				if ab21aq21Value != (ab21aq21Key == ab21aq21Name) {
					ab21aq21T.Fatal("flag position", ab21aq21Bit)
				}
			}
		} else if ab21aq21Bit < 6 {
			ab21aq21Subset(ab21aq21T, ab21aq21Result["reserved"], []float64{1, 2}[ab21aq21Bit-4])
		} else {
			ab21aq21Subset(ab21aq21T, ab21aq21Result["unmodeled"], []float64{1, 2}[ab21aq21Bit-6])
		}
	}
	ab21aq21Input := ab21aq21Cases(ab21aq21T)[5].Ab21Aq21Input
	ab21aq21Input["arguments"].(map[string]any)["value_width"] = float64(1)
	ab21aq21Result := ab21aq21Eval(ab21aq21T, ab21aq21Input)["results"].(map[string]any)["value"].(map[string]any)
	if ab21aq21Result["state"] != "incomplete" {
		ab21aq21T.Fatal("width inferred from bytes")
	}
}
