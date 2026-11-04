package vaillant

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"
)

func ab05aq05Decode(ab05aq05T *testing.T, ab05aq05Input []byte) map[string]any {
	ab05aq05T.Helper()
	ab05aq05Output, ab05aq05Error := Ab05Aq05Evaluate(ab05aq05Input)
	if ab05aq05Error != nil {
		ab05aq05T.Fatal(ab05aq05Error)
	}
	var ab05aq05Object map[string]any
	if ab05aq05Error = json.Unmarshal(ab05aq05Output, &ab05aq05Object); ab05aq05Error != nil {
		ab05aq05T.Fatal(ab05aq05Error)
	}
	return ab05aq05Object
}

func ab05aq05Input(ab05aq05View string, ab05aq05Operation int, ab05aq05Fields, ab05aq05Arguments map[string]any) map[string]any {
	return map[string]any{"context": map[string]any{"primary": 181, "secondary": 5, "operation": ab05aq05Operation, "direction": "request", "view": ab05aq05View, "family": "synthetic-family", "variant": "synthetic-variant"}, "fields": ab05aq05Fields, "arguments": ab05aq05Arguments}
}

func ab05aq05Run(ab05aq05T *testing.T, ab05aq05Object map[string]any) map[string]any {
	ab05aq05T.Helper()
	ab05aq05Bytes, ab05aq05Error := json.Marshal(ab05aq05Object)
	if ab05aq05Error != nil {
		ab05aq05T.Fatal(ab05aq05Error)
	}
	return ab05aq05Decode(ab05aq05T, ab05aq05Bytes)
}

func TestAb05Aq05Fixture(ab05aq05T *testing.T) {
	ab05aq05Data, ab05aq05Error := os.ReadFile("testdata/ab05aq05_b505_conditional_roles.json")
	if ab05aq05Error != nil {
		ab05aq05T.Fatal(ab05aq05Error)
	}
	var ab05aq05Cases []struct {
		Ab05Aq05ID       string          `json:"id"`
		Ab05Aq05Input    json.RawMessage `json:"input"`
		Ab05Aq05Key      string          `json:"key"`
		Ab05Aq05Expected map[string]any  `json:"expected"`
	}
	if ab05aq05Error = json.Unmarshal(ab05aq05Data, &ab05aq05Cases); ab05aq05Error != nil {
		ab05aq05T.Fatal(ab05aq05Error)
	}
	if len(ab05aq05Cases) != 15 {
		ab05aq05T.Fatalf("fixture case count %d", len(ab05aq05Cases))
	}
	for _, ab05aq05Case := range ab05aq05Cases {
		ab05aq05T.Run(ab05aq05Case.Ab05Aq05ID, func(ab05aq05T *testing.T) {
			ab05aq05Output := ab05aq05Decode(ab05aq05T, ab05aq05Case.Ab05Aq05Input)
			if ab05aq05Output["state"] != "complete" || ab05aq05Output["physical"] != "Unknown" || ab05aq05Output["qualified"] != false {
				ab05aq05T.Fatal(ab05aq05Output)
			}
			ab05aq05Result := ab05aq05Output["results"].(map[string]any)[ab05aq05Case.Ab05Aq05Key].(map[string]any)
			if ab05aq05Result["grade"] != "Hypothesis" || ab05aq05Result["state"] != "complete" {
				ab05aq05T.Fatal(ab05aq05Result)
			}
			for ab05aq05Key, ab05aq05Expected := range ab05aq05Case.Ab05Aq05Expected {
				if !reflect.DeepEqual(ab05aq05Result[ab05aq05Key], ab05aq05Expected) {
					ab05aq05T.Fatalf("%s: got %v want %v", ab05aq05Key, ab05aq05Result[ab05aq05Key], ab05aq05Expected)
				}
			}
		})
	}
}

func TestAb05Aq05Identity(ab05aq05T *testing.T) {
	for _, ab05aq05Axis := range []string{"primary", "secondary", "operation-absent", "operation-zero", "operation-other", "direction", "view", "family", "variant"} {
		ab05aq05T.Run(ab05aq05Axis, func(ab05aq05T *testing.T) {
			ab05aq05Object := ab05aq05Input("dhw-enable", 49, map[string]any{"enabled": []int{1}}, map[string]any{})
			ab05aq05Context := ab05aq05Object["context"].(map[string]any)
			switch ab05aq05Axis {
			case "primary":
				ab05aq05Context["primary"] = 180
			case "secondary":
				ab05aq05Context["secondary"] = 6
			case "operation-absent":
				delete(ab05aq05Context, "operation")
			case "operation-zero":
				ab05aq05Context["operation"] = 0
			case "operation-other":
				ab05aq05Context["operation"] = 9
			default:
				ab05aq05Context[ab05aq05Axis] = "other"
			}
			ab05aq05Object["extra"] = map[string]any{"null": nil, "zero": 0, "false": false, "empty": []any{}}
			ab05aq05Output := ab05aq05Run(ab05aq05T, ab05aq05Object)
			ab05aq05State := "Unknown"
			if ab05aq05Axis == "family" || ab05aq05Axis == "variant" {
				ab05aq05State = "complete"
			}
			if ab05aq05Output["state"] != ab05aq05State || ab05aq05Output["qualified"] != false {
				ab05aq05T.Fatal(ab05aq05Output)
			}
			ab05aq05Result := ab05aq05Output["results"].(map[string]any)["enabled"].(map[string]any)
			if ab05aq05State == "Unknown" {
				if len(ab05aq05Result) != 2 || ab05aq05Result["state"] != "Unknown" {
					ab05aq05T.Fatal(ab05aq05Result)
				}
			}
			ab05aq05Bytes, _ := json.Marshal(ab05aq05Object)
			var ab05aq05Original map[string]any
			_ = json.Unmarshal(ab05aq05Bytes, &ab05aq05Original)
			for _, ab05aq05Key := range []string{"context", "fields", "arguments", "extra"} {
				if !reflect.DeepEqual(ab05aq05Output[ab05aq05Key], ab05aq05Original[ab05aq05Key]) {
					ab05aq05T.Fatal("lost raw", ab05aq05Key)
				}
			}
		})
	}
	for _, ab05aq05Operation := range []int{61, 90} {
		ab05aq05Object := ab05aq05Input("unchanged", ab05aq05Operation, map[string]any{}, map[string]any{"replacement": true, "previous": 17, "requested": 23})
		if ab05aq05Run(ab05aq05T, ab05aq05Object)["state"] != "complete" {
			ab05aq05T.Fatal("operation alias")
		}
	}
}

func TestAb05Aq05PartialAndOpaque(ab05aq05T *testing.T) {
	for _, ab05aq05Field := range []string{"enabled", "disable", "start"} {
		for _, ab05aq05Width := range []int{-1, 0, 1, 2} {
			ab05aq05View := "dhw-enable"
			ab05aq05Operation := 49
			ab05aq05Args := map[string]any{}
			if ab05aq05Field == "disable" {
				ab05aq05View = "parallel-charge"
				ab05aq05Operation = 24
				ab05aq05Args = map[string]any{"role": "mixer", "parallel": true}
			}
			if ab05aq05Field == "start" {
				ab05aq05View = "programmed-start"
				ab05aq05Operation = 9
				ab05aq05Args = map[string]any{"program_mode": 0}
			}
			ab05aq05Fields := map[string]any{"independent": []int{17}, "empty": []int{}}
			if ab05aq05Width >= 0 {
				ab05aq05Fields[ab05aq05Field] = make([]int, ab05aq05Width)
			}
			ab05aq05Object := ab05aq05Input(ab05aq05View, ab05aq05Operation, ab05aq05Fields, ab05aq05Args)
			ab05aq05Output := ab05aq05Run(ab05aq05T, ab05aq05Object)
			ab05aq05Complete := ab05aq05Width == 1 || (ab05aq05Field == "start" && ab05aq05Width >= 0)
			ab05aq05State := "partial"
			if ab05aq05Complete {
				ab05aq05State = "complete"
			}
			if ab05aq05Output["state"] != ab05aq05State {
				ab05aq05T.Fatal(ab05aq05Field, ab05aq05Width, ab05aq05Output)
			}
			ab05aq05Results := ab05aq05Output["results"].(map[string]any)
			if !reflect.DeepEqual(ab05aq05Results["independent"].(map[string]any)["raw"], []any{float64(17)}) || !reflect.DeepEqual(ab05aq05Results["empty"].(map[string]any)["raw"], []any{}) {
				ab05aq05T.Fatal("siblings lost")
			}
			ab05aq05Result := ab05aq05Results[ab05aq05Field].(map[string]any)
			_, ab05aq05HasRaw := ab05aq05Result["raw"]
			if ab05aq05HasRaw != (ab05aq05Width >= 0) {
				ab05aq05T.Fatal("missing/empty collapsed")
			}
			if !ab05aq05Complete {
				if _, ab05aq05Exists := ab05aq05Result["number"]; ab05aq05Exists {
					ab05aq05T.Fatal("invented number")
				}
			}
		}
	}
}

func TestAb05Aq05Reject(ab05aq05T *testing.T) {
	for _, ab05aq05Bad := range []string{``, `{`, `[]`, `{} {}`, `{"context":{},"fields":{},"x":{"a":1,"a":2}}`, `{"context":{},"fields":{},"x":[{"a":1,"a":2}]}`, `{"context":{},"context":{},"fields":{}}`, `{"context":{},"fields":{},"x":NaN}`, `{"context":{},"fields":{},"x":1e999}`, `{"context":{},"fields":null}`, `{"context":null,"fields":{}}`, `{"context":{},"fields":{},"arguments":null}`} {
		if ab05aq05Output, ab05aq05Error := Ab05Aq05Evaluate([]byte(ab05aq05Bad)); ab05aq05Error == nil || ab05aq05Output != nil {
			ab05aq05T.Fatalf("accepted %s", ab05aq05Bad)
		}
	}
	for _, ab05aq05Key := range []string{"primary", "secondary", "operation", "enabled", "program_mode", "replacement", "previous", "requested", "role", "parallel", "demand", "excess"} {
		ab05aq05BadValues := []any{nil, "1", true, 1.5, -1, 256}
		if ab05aq05Key == "replacement" || ab05aq05Key == "parallel" {
			ab05aq05BadValues = []any{nil, "true", 0, []any{}}
		}
		if ab05aq05Key == "previous" || ab05aq05Key == "requested" {
			ab05aq05BadValues = []any{[]any{}, map[string]any{}}
		}
		if ab05aq05Key == "role" {
			ab05aq05BadValues = []any{nil, 0, false, "other"}
		}
		if ab05aq05Key == "demand" || ab05aq05Key == "excess" {
			ab05aq05BadValues = []any{nil, "1", true, []any{}}
		}
		for _, ab05aq05Bad := range ab05aq05BadValues {
			ab05aq05Object := ab05aq05Input("dhw-enable", 49, map[string]any{"enabled": []int{1}}, map[string]any{})
			ab05aq05Context := ab05aq05Object["context"].(map[string]any)
			ab05aq05Args := ab05aq05Object["arguments"].(map[string]any)
			switch ab05aq05Key {
			case "primary", "secondary", "operation":
				ab05aq05Context[ab05aq05Key] = ab05aq05Bad
			case "enabled":
				ab05aq05Object["fields"] = map[string]any{"enabled": []any{ab05aq05Bad}}
			case "program_mode":
				ab05aq05Context["view"] = "programmed-start"
				ab05aq05Context["operation"] = 9
				ab05aq05Args[ab05aq05Key] = ab05aq05Bad
			case "replacement", "previous", "requested":
				ab05aq05Context["view"] = "unchanged"
				ab05aq05Context["operation"] = 90
				ab05aq05Args["replacement"] = true
				ab05aq05Args["previous"] = 17
				ab05aq05Args["requested"] = 23
				ab05aq05Args[ab05aq05Key] = ab05aq05Bad
			case "role", "parallel":
				ab05aq05Context["view"] = "parallel-charge"
				ab05aq05Context["operation"] = 24
				ab05aq05Args["role"] = "mixer"
				ab05aq05Args["parallel"] = true
				ab05aq05Args[ab05aq05Key] = ab05aq05Bad
			case "demand", "excess":
				ab05aq05Context["view"] = "flow-excess"
				ab05aq05Context["operation"] = 23
				ab05aq05Args["demand"] = 40
				ab05aq05Args["excess"] = 5
				ab05aq05Args[ab05aq05Key] = ab05aq05Bad
			}
			ab05aq05Bytes, _ := json.Marshal(ab05aq05Object)
			if ab05aq05Output, ab05aq05Error := Ab05Aq05Evaluate(ab05aq05Bytes); ab05aq05Error == nil || ab05aq05Output != nil {
				ab05aq05T.Fatalf("accepted %s=%v", ab05aq05Key, ab05aq05Bad)
			}
		}
	}
	for _, ab05aq05Case := range []struct {
		Ab05Aq05View      string
		Ab05Aq05Operation int
		Ab05Aq05Args      map[string]any
	}{{"unchanged", 90, map[string]any{"replacement": true, "previous": 17, "requested": 23}}, {"parallel-charge", 24, map[string]any{"role": "mixer", "parallel": true}}, {"flow-excess", 23, map[string]any{"demand": 40, "excess": 5}}, {"programmed-start", 9, map[string]any{"program_mode": 0}}} {
		for ab05aq05Key := range ab05aq05Case.Ab05Aq05Args {
			ab05aq05Args := map[string]any{}
			for ab05aq05Name, ab05aq05Value := range ab05aq05Case.Ab05Aq05Args {
				if ab05aq05Name != ab05aq05Key {
					ab05aq05Args[ab05aq05Name] = ab05aq05Value
				}
			}
			ab05aq05Bytes, _ := json.Marshal(ab05aq05Input(ab05aq05Case.Ab05Aq05View, ab05aq05Case.Ab05Aq05Operation, map[string]any{}, ab05aq05Args))
			if ab05aq05Output, ab05aq05Error := Ab05Aq05Evaluate(ab05aq05Bytes); ab05aq05Error == nil || ab05aq05Output != nil {
				ab05aq05T.Fatal("accepted missing", ab05aq05Key)
			}
		}
	}
}

func TestAb05Aq05ScalarsAndMemory(ab05aq05T *testing.T) {
	for _, ab05aq05Value := range []any{nil, false, 0, "", 17.5} {
		for _, ab05aq05Replacement := range []bool{false, true} {
			ab05aq05Object := ab05aq05Input("unchanged", 90, map[string]any{}, map[string]any{"replacement": ab05aq05Replacement, "previous": ab05aq05Value, "requested": ab05aq05Value})
			ab05aq05Output := ab05aq05Run(ab05aq05T, ab05aq05Object)
			ab05aq05Bytes, _ := json.Marshal(ab05aq05Value)
			var ab05aq05Expected any
			_ = json.Unmarshal(ab05aq05Bytes, &ab05aq05Expected)
			if !reflect.DeepEqual(ab05aq05Output["results"].(map[string]any)["unchanged"].(map[string]any)["value"], ab05aq05Expected) {
				ab05aq05T.Fatal("scalar collapsed")
			}
		}
	}
	ab05aq05Object := ab05aq05Input("flow-excess", 23, map[string]any{}, map[string]any{"demand": 0.1, "excess": 0.2})
	if ab05aq05Run(ab05aq05T, ab05aq05Object)["results"].(map[string]any)["flow-excess"].(map[string]any)["number"] != 0.3 {
		ab05aq05T.Fatal("decimal sum")
	}
	ab05aq05InputBytes, _ := json.Marshal(ab05aq05Object)
	ab05aq05Original := append([]byte(nil), ab05aq05InputBytes...)
	ab05aq05First, ab05aq05Error := Ab05Aq05Evaluate(ab05aq05InputBytes)
	if ab05aq05Error != nil {
		ab05aq05T.Fatal(ab05aq05Error)
	}
	ab05aq05Saved := append([]byte(nil), ab05aq05First...)
	for ab05aq05Index := range ab05aq05InputBytes {
		ab05aq05InputBytes[ab05aq05Index] = ' '
	}
	if string(ab05aq05First) != string(ab05aq05Saved) {
		ab05aq05T.Fatal("input aliases output")
	}
	for ab05aq05Index := range ab05aq05First {
		ab05aq05First[ab05aq05Index] = ' '
	}
	ab05aq05Second, ab05aq05Error := Ab05Aq05Evaluate(ab05aq05Original)
	if ab05aq05Error != nil || string(ab05aq05Second) != string(ab05aq05Saved) {
		ab05aq05T.Fatal("cross-call storage")
	}
	if strings.Contains(string(ab05aq05Second), "replacement_code") {
		ab05aq05T.Fatal("cross-view contamination")
	}
}

func TestAb05Aq05IndependentInputs(ab05aq05T *testing.T) {
	for _, ab05aq05Case := range []struct {
		Ab05Aq05Demand float64
		Ab05Aq05Excess float64
		Ab05Aq05Want   float64
	}{{41, 5, 46}, {40, 6, 46}, {0, 5, 5}, {40, 0, 40}, {-40, 5, -35}} {
		ab05aq05Result := ab05aq05Run(ab05aq05T, ab05aq05Input("flow-excess", 23, map[string]any{}, map[string]any{"demand": ab05aq05Case.Ab05Aq05Demand, "excess": ab05aq05Case.Ab05Aq05Excess}))["results"].(map[string]any)["flow-excess"].(map[string]any)
		if ab05aq05Result["number"] != ab05aq05Case.Ab05Aq05Want {
			ab05aq05T.Fatal(ab05aq05Result)
		}
	}
	for _, ab05aq05Key := range []string{"previous", "requested"} {
		ab05aq05Args := map[string]any{"replacement": ab05aq05Key == "previous", "previous": 17, "requested": 23}
		ab05aq05Args[ab05aq05Key] = 31
		ab05aq05Result := ab05aq05Run(ab05aq05T, ab05aq05Input("unchanged", 61, map[string]any{"replacement": []int{0}}, ab05aq05Args))["results"].(map[string]any)["unchanged"].(map[string]any)
		if ab05aq05Result["number"] != float64(31) || ab05aq05Result["replacement_code"] != "Unknown" {
			ab05aq05T.Fatal(ab05aq05Result)
		}
	}
	for _, ab05aq05Mode := range []int{0, 1, 255} {
		for _, ab05aq05Start := range [][]int{{}, {0}, {255, 1, 2, 3}} {
			ab05aq05Result := ab05aq05Run(ab05aq05T, ab05aq05Input("programmed-start", 9, map[string]any{"start": ab05aq05Start}, map[string]any{"program_mode": ab05aq05Mode}))["results"].(map[string]any)["start"].(map[string]any)
			ab05aq05Label := "ignored"
			if ab05aq05Mode == 0 {
				ab05aq05Label = "applicable"
			}
			if ab05aq05Result["label"] != ab05aq05Label {
				ab05aq05T.Fatal(ab05aq05Result)
			}
			ab05aq05Bytes, _ := json.Marshal(ab05aq05Start)
			var ab05aq05Expected any
			_ = json.Unmarshal(ab05aq05Bytes, &ab05aq05Expected)
			if !reflect.DeepEqual(ab05aq05Result["raw"], ab05aq05Expected) {
				ab05aq05T.Fatal("start changed")
			}
		}
	}
	for _, ab05aq05Field := range []string{"enabled", "disable", "start"} {
		for _, ab05aq05Bad := range []any{nil, 1, "1", false, []any{nil}, []any{true}, []any{"1"}, []any{256}, []any{1.5}} {
			ab05aq05Bytes, _ := json.Marshal(ab05aq05Input("unknown", 0, map[string]any{ab05aq05Field: ab05aq05Bad}, map[string]any{}))
			if ab05aq05Output, ab05aq05Error := Ab05Aq05Evaluate(ab05aq05Bytes); ab05aq05Error == nil || ab05aq05Output != nil {
				ab05aq05T.Fatal("invalid raw field accepted", ab05aq05Field, ab05aq05Bad)
			}
		}
	}
	for _, ab05aq05Bad := range []string{`{"context":{"operation":1.0000000000000000001},"fields":{}}`, `{"context":{},"fields":{"enabled":[0.99999999999999999999]}}`, `{"context":{},"fields":{}} false`} {
		if ab05aq05Output, ab05aq05Error := Ab05Aq05Evaluate([]byte(ab05aq05Bad)); ab05aq05Error == nil || ab05aq05Output != nil {
			ab05aq05T.Fatal("precision/trailing accepted")
		}
	}
}

func TestAb05Aq05BoundedNumbers(ab05aq05T *testing.T) {
	ab05aq05Prefix := `{"context":{"primary":181,"secondary":5,"operation":23,"direction":"request","view":"flow-excess"},"fields":{},"arguments":{`
	for _, ab05aq05Bad := range []string{`"demand":1e-100000,"excess":5`, `"demand":1,"excess":1e1000000`, `"demand":1e-401,"excess":5`, `"demand":0.` + strings.Repeat("0", 70) + `1,"excess":5`} {
		if ab05aq05Output, ab05aq05Error := Ab05Aq05Evaluate([]byte(ab05aq05Prefix + ab05aq05Bad + `}}`)); ab05aq05Error == nil || ab05aq05Output != nil {
			ab05aq05T.Fatalf("accepted %s", ab05aq05Bad)
		}
	}
	ab05aq05Output, ab05aq05Error := Ab05Aq05Evaluate([]byte(ab05aq05Prefix + `"demand":40,"excess":5e0}}`))
	if ab05aq05Error != nil || !strings.Contains(string(ab05aq05Output), `"number":45`) {
		ab05aq05T.Fatalf("bounded sum failed: %s %v", ab05aq05Output, ab05aq05Error)
	}
}

func TestAb05Aq05ReservedKeysAndCollisions(ab05aq05T *testing.T) {
	for _, ab05aq05Key := range []string{"state", "results", "physical", "qualified"} {
		ab05aq05Object := ab05aq05Input("dhw-enable", 49, map[string]any{"enabled": []int{1}}, map[string]any{})
		ab05aq05Object[ab05aq05Key] = "Proven"
		ab05aq05Bytes, _ := json.Marshal(ab05aq05Object)
		if ab05aq05Output, ab05aq05Error := Ab05Aq05Evaluate(ab05aq05Bytes); ab05aq05Error == nil || ab05aq05Output != nil {
			ab05aq05T.Fatalf("accepted reserved key %s", ab05aq05Key)
		}
	}
	ab05aq05Object := ab05aq05Input("flow-excess", 23, map[string]any{"flow-excess": []int{7, 8}}, map[string]any{"demand": json.Number("40"), "excess": json.Number("5")})
	ab05aq05Results := ab05aq05Run(ab05aq05T, ab05aq05Object)["results"].(map[string]any)
	ab05aq05View := ab05aq05Results["flow-excess"].(map[string]any)
	if !reflect.DeepEqual(ab05aq05View["field_raw"], []any{float64(7), float64(8)}) || ab05aq05View["number"] != float64(45) {
		ab05aq05T.Fatalf("collision lost field raw: %v", ab05aq05View)
	}
}

func TestAb05Aq05CanonicalSum(ab05aq05T *testing.T) {
	ab05aq05Object := ab05aq05Input("flow-excess", 23, map[string]any{}, map[string]any{"demand": json.Number("0.25"), "excess": json.Number("0.1")})
	ab05aq05Bytes, _ := json.Marshal(ab05aq05Object)
	ab05aq05Output, ab05aq05Error := Ab05Aq05Evaluate(ab05aq05Bytes)
	if ab05aq05Error != nil {
		ab05aq05T.Fatal(ab05aq05Error)
	}
	if !strings.Contains(string(ab05aq05Output), `"number":0.35,`) && !strings.Contains(string(ab05aq05Output), `"number":0.35}`) {
		ab05aq05T.Fatalf("non-canonical sum: %s", ab05aq05Output)
	}
}

func TestAb05Aq05RejectDeepNesting(ab05aq05T *testing.T) {
	ab05aq05Deep := `{"context":{},"fields":{},"x":` + strings.Repeat("[", 100000) + strings.Repeat("]", 100000) + `}`
	if ab05aq05Output, ab05aq05Error := Ab05Aq05Evaluate([]byte(ab05aq05Deep)); ab05aq05Error == nil || ab05aq05Output != nil {
		ab05aq05T.Fatalf("accepted over-deep input")
	}
	ab05aq05Shallow := `{"context":{},"fields":{},"x":` + strings.Repeat("[", 8) + strings.Repeat("]", 8) + `}`
	if _, ab05aq05Error := Ab05Aq05Evaluate([]byte(ab05aq05Shallow)); ab05aq05Error != nil {
		ab05aq05T.Fatalf("rejected shallow input: %v", ab05aq05Error)
	}
}
