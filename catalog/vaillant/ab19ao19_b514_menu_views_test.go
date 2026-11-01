package vaillant

import (
	ab19ao19json "encoding/json"
	ab19ao19fmt "fmt"
	ab19ao19os "os"
	ab19ao19reflect "reflect"
	ab19ao19strings "strings"
	ab19ao19testing "testing"
)

func ab19ao19Input(ab19ao19Operation int, ab19ao19View string, ab19ao19Fields map[string]any) map[string]any {
	return map[string]any{"context": map[string]any{"primary": 181, "secondary": 20, "operation": ab19ao19Operation, "direction": "positive", "view": ab19ao19View, "family": map[string]any{"opaque": []any{"synthetic-family", 0}}, "variant": []any{"synthetic-variant", false}}, "fields": ab19ao19Fields}
}
func ab19ao19Call(ab19ao19T *ab19ao19testing.T, ab19ao19InputObject map[string]any) map[string]any {
	ab19ao19T.Helper()
	ab19ao19InputJSON, ab19ao19Err := ab19ao19json.Marshal(ab19ao19InputObject)
	if ab19ao19Err != nil {
		ab19ao19T.Fatal(ab19ao19Err)
	}
	ab19ao19OutputJSON, ab19ao19Err := Ab19Ao19Interpret(ab19ao19InputJSON)
	if ab19ao19Err != nil {
		ab19ao19T.Fatal(ab19ao19Err)
	}
	ab19ao19Output := map[string]any{}
	if ab19ao19Err := ab19ao19json.Unmarshal(ab19ao19OutputJSON, &ab19ao19Output); ab19ao19Err != nil {
		ab19ao19T.Fatal(ab19ao19Err)
	}
	if ab19ao19Output["physical"] != "Unknown" || ab19ao19Output["qualified"] != false {
		ab19ao19T.Fatal("qualification promoted", ab19ao19Output)
	}
	ab19ao19Expected := map[string]any{}
	if ab19ao19Err := ab19ao19json.Unmarshal(ab19ao19InputJSON, &ab19ao19Expected); ab19ao19Err != nil {
		ab19ao19T.Fatal(ab19ao19Err)
	}
	for ab19ao19Key, ab19ao19Value := range map[string]any{"context": ab19ao19Expected["context"], "raw": ab19ao19Expected["fields"]} {
		if !ab19ao19reflect.DeepEqual(ab19ao19Output[ab19ao19Key], ab19ao19Value) {
			ab19ao19T.Fatalf("lost %s: %#v", ab19ao19Key, ab19ao19Output)
		}
	}
	ab19ao19ActualArguments, ab19ao19ActualPresent := ab19ao19Output["arguments"]
	ab19ao19ExpectedArguments, ab19ao19ExpectedPresent := ab19ao19Expected["arguments"]
	if ab19ao19ActualPresent != ab19ao19ExpectedPresent || !ab19ao19reflect.DeepEqual(ab19ao19ActualArguments, ab19ao19ExpectedArguments) {
		ab19ao19T.Fatal("argument presence changed")
	}
	for _, ab19ao19Field := range ab19ao19Output["fields"].(map[string]any) {
		if ab19ao19Field.(map[string]any)["grade"] != "Hypothesis" {
			ab19ao19T.Fatal("family grade changed")
		}
	}
	return ab19ao19Output
}
func ab19ao19Field(ab19ao19Output map[string]any, ab19ao19Name string) map[string]any {
	return ab19ao19Output["fields"].(map[string]any)[ab19ao19Name].(map[string]any)
}
func ab19ao19NoNumber(ab19ao19T *ab19ao19testing.T, ab19ao19Field map[string]any) {
	ab19ao19T.Helper()
	if _, ab19ao19Present := ab19ao19Field["number"]; ab19ao19Present {
		ab19ao19T.Fatal("unexpected number", ab19ao19Field)
	}
}
func TestAb19Ao19Feedback(ab19ao19T *ab19ao19testing.T) {
	for _, ab19ao19View := range []string{"generic-result", "specific-a-result", "specific-b-result"} {
		for ab19ao19Byte := 0; ab19ao19Byte <= 255; ab19ao19Byte++ {
			ab19ao19T.Run(ab19ao19fmt.Sprintf("%s/%d", ab19ao19View, ab19ao19Byte), func(ab19ao19T *ab19ao19testing.T) {
				ab19ao19Output := ab19ao19Call(ab19ao19T, ab19ao19Input(6, ab19ao19View, map[string]any{"feedback": []int{ab19ao19Byte, 99}, "unmodelled": []int{4, 5}}))
				ab19ao19Expected := ""
				if ab19ao19Byte == 1 {
					ab19ao19Expected = "not-successful"
					if ab19ao19View == "generic-result" {
						ab19ao19Expected = "in-progress"
					}
				}
				if ab19ao19Byte == 255 && ab19ao19View == "generic-result" {
					ab19ao19Expected = "unsupported"
				}
				if ab19ao19Byte >= 3 && ab19ao19View != "generic-result" {
					ab19ao19Expected = "invalid"
				}
				ab19ao19FieldValue := ab19ao19Field(ab19ao19Output, "feedback")
				ab19ao19Label, _ := ab19ao19FieldValue["label"].(string)
				if ab19ao19Label != ab19ao19Expected || ab19ao19Output["state"] != "Hypothesis" {
					ab19ao19T.Fatal(ab19ao19Output)
				}
				ab19ao19State := "Hypothesis"
				if ab19ao19Expected == "" {
					ab19ao19State = "unmapped"
				}
				if ab19ao19FieldValue["state"] != ab19ao19State || !ab19ao19reflect.DeepEqual(ab19ao19FieldValue["raw"], []any{float64(ab19ao19Byte)}) {
					ab19ao19T.Fatal(ab19ao19FieldValue)
				}
				ab19ao19NoNumber(ab19ao19T, ab19ao19FieldValue)
			})
		}
	}
	for ab19ao19Byte := 0; ab19ao19Byte <= 255; ab19ao19Byte++ {
		ab19ao19T.Run(ab19ao19fmt.Sprintf("progress/%d", ab19ao19Byte), func(ab19ao19T *ab19ao19testing.T) {
			ab19ao19Output := ab19ao19Call(ab19ao19T, ab19ao19Input(7, "progress", map[string]any{"feedback": []int{ab19ao19Byte}}))
			ab19ao19Expected := ""
			if ab19ao19Byte < 7 {
				ab19ao19Expected = []string{"idle", "progress", "passed", "warning", "combustion-error", "gas-type-error", "aborted"}[ab19ao19Byte]
			}
			if ab19ao19Byte == 255 {
				ab19ao19Expected = "unsupported"
			}
			ab19ao19FieldValue := ab19ao19Field(ab19ao19Output, "feedback")
			ab19ao19Label, _ := ab19ao19FieldValue["label"].(string)
			ab19ao19State := "Hypothesis"
			if ab19ao19Expected == "" {
				ab19ao19State = "unmapped"
			}
			if ab19ao19Label != ab19ao19Expected || ab19ao19FieldValue["state"] != ab19ao19State {
				ab19ao19T.Fatal(ab19ao19Output)
			}
			ab19ao19NoNumber(ab19ao19T, ab19ao19FieldValue)
		})
	}
}
func ab19ao19Pressure() []int { return []int{19, 23, 29, 1, 1, 0, 0, 1, 52, 18, 0, 100, 0} }
func ab19ao19PressureCheck(ab19ao19T *ab19ao19testing.T, ab19ao19Data []int, ab19ao19Expected map[string]float64) map[string]any {
	ab19ao19T.Helper()
	ab19ao19Output := ab19ao19Call(ab19ao19T, ab19ao19Input(8, "pressure", map[string]any{"logical_object": ab19ao19Data, "other": []int{7}}))
	for ab19ao19Name, ab19ao19Number := range ab19ao19Expected {
		ab19ao19Value := ab19ao19Field(ab19ao19Output, ab19ao19Name)
		if ab19ao19Value["number"] != ab19ao19Number || ab19ao19Value["calculation"] != "Proven" || ab19ao19Value["state"] != "Hypothesis" {
			ab19ao19T.Fatalf("%s: %#v", ab19ao19Name, ab19ao19Value)
		}
	}
	ab19ao19Step := ab19ao19Field(ab19ao19Output, "step")
	if len(ab19ao19Data) > 10 {
		if ab19ao19Step["state"] != "unit-conflict" || !ab19ao19reflect.DeepEqual(ab19ao19Step["flags"], []any{"constant-1-mbar-versus-default-0"}) {
			ab19ao19T.Fatal(ab19ao19Step)
		}
	}
	ab19ao19NoNumber(ab19ao19T, ab19ao19Step)
	ab19ao19NoNumber(ab19ao19T, ab19ao19Field(ab19ao19Output, "feedback"))
	return ab19ao19Output
}
func TestAb19Ao19Pressure(ab19ao19T *ab19ao19testing.T) {
	for _, ab19ao19Case := range []struct {
		ab19ao19Index, ab19ao19Byte int
		ab19ao19Name                string
		ab19ao19Number              float64
	}{
		{4, 2, "pressure", 2}, {5, 1, "pressure", 257}, {6, 3, "minimum", 259}, {7, 2, "minimum", 512}, {8, 53, "maximum", 4661}, {9, 19, "maximum", 4916}, {11, 200, "flow", 2}, {12, 1, "flow", 3.56},
	} {
		ab19ao19T.Run(ab19ao19fmt.Sprintf("word/%d", ab19ao19Case.ab19ao19Index), func(ab19ao19T *ab19ao19testing.T) {
			ab19ao19Data := ab19ao19Pressure()
			ab19ao19Data[ab19ao19Case.ab19ao19Index] = ab19ao19Case.ab19ao19Byte
			ab19ao19Expected := map[string]float64{"pressure": 1, "minimum": 256, "maximum": 4660, "flow": 1}
			ab19ao19Expected[ab19ao19Case.ab19ao19Name] = ab19ao19Case.ab19ao19Number
			ab19ao19PressureCheck(ab19ao19T, ab19ao19Data, ab19ao19Expected)
		})
	}
	for _, ab19ao19Index := range []int{0, 1, 2, 3, 10, 13} {
		ab19ao19T.Run(ab19ao19fmt.Sprintf("opaque-feedback-step-tail/%d", ab19ao19Index), func(ab19ao19T *ab19ao19testing.T) {
			ab19ao19Data := append(ab19ao19Pressure(), 81)
			ab19ao19Data[ab19ao19Index] = 255
			ab19ao19Output := ab19ao19PressureCheck(ab19ao19T, ab19ao19Data, map[string]float64{"pressure": 1, "minimum": 256, "maximum": 4660, "flow": 1})
			if ab19ao19Output["state"] != "Hypothesis" || ab19ao19Field(ab19ao19Output, "feedback")["state"] != "unmapped" {
				ab19ao19T.Fatal(ab19ao19Output)
			}
		})
	}
	for _, ab19ao19Step := range []int{0, 1, 2, 255} {
		ab19ao19Data := ab19ao19Pressure()
		ab19ao19Data[10] = ab19ao19Step
		ab19ao19PressureCheck(ab19ao19T, ab19ao19Data, map[string]float64{"pressure": 1, "minimum": 256, "maximum": 4660, "flow": 1})
	}
	for ab19ao19Length := 0; ab19ao19Length <= 13; ab19ao19Length++ {
		ab19ao19T.Run(ab19ao19fmt.Sprintf("length/%d", ab19ao19Length), func(ab19ao19T *ab19ao19testing.T) {
			ab19ao19Data := ab19ao19Pressure()[:ab19ao19Length]
			ab19ao19Expected := map[string]float64{}
			for ab19ao19Name, ab19ao19End := range map[string]int{"pressure": 6, "minimum": 8, "maximum": 10, "flow": 13} {
				if ab19ao19Length >= ab19ao19End {
					ab19ao19Expected[ab19ao19Name] = map[string]float64{"pressure": 1, "minimum": 256, "maximum": 4660, "flow": 1}[ab19ao19Name]
				}
			}
			ab19ao19Output := ab19ao19PressureCheck(ab19ao19T, ab19ao19Data, ab19ao19Expected)
			ab19ao19State := "partial"
			if ab19ao19Length == 13 {
				ab19ao19State = "Hypothesis"
			}
			if ab19ao19Output["state"] != ab19ao19State {
				ab19ao19T.Fatal(ab19ao19Output)
			}
			for ab19ao19Name, ab19ao19End := range map[string]int{"feedback": 4, "pressure": 6, "minimum": 8, "maximum": 10, "step": 11, "flow": 13} {
				if ab19ao19Length < ab19ao19End {
					ab19ao19Value := ab19ao19Field(ab19ao19Output, ab19ao19Name)
					if ab19ao19Value["state"] != "incomplete" {
						ab19ao19T.Fatal(ab19ao19Value)
					}
					ab19ao19NoNumber(ab19ao19T, ab19ao19Value)
				}
			}
		})
	}
	ab19ao19Data := make([]int, 13)
	ab19ao19PressureCheck(ab19ao19T, ab19ao19Data, map[string]float64{"pressure": 0, "minimum": 0, "maximum": 0, "flow": 0})
}
func TestAb19Ao19Selection(ab19ao19T *ab19ao19testing.T) {
	for _, ab19ao19Selection := range []struct {
		ab19ao19Operation         int
		ab19ao19View, ab19ao19Raw string
	}{{6, "generic-result", "feedback"}, {6, "specific-a-result", "feedback"}, {6, "specific-b-result", "feedback"}, {7, "progress", "feedback"}, {8, "pressure", "logical_object"}} {
		for _, ab19ao19Key := range []string{"primary", "secondary", "operation", "direction", "view", "family", "variant"} {
			for _, ab19ao19Missing := range []bool{false, true} {
				ab19ao19T.Run(ab19ao19fmt.Sprintf("%d/%s/%s/%t", ab19ao19Selection.ab19ao19Operation, ab19ao19Selection.ab19ao19View, ab19ao19Key, ab19ao19Missing), func(ab19ao19T *ab19ao19testing.T) {
					ab19ao19Data := []int{1}
					if ab19ao19Selection.ab19ao19Operation == 8 {
						ab19ao19Data = ab19ao19Pressure()
					}
					ab19ao19InputObject := ab19ao19Input(ab19ao19Selection.ab19ao19Operation, ab19ao19Selection.ab19ao19View, map[string]any{ab19ao19Selection.ab19ao19Raw: ab19ao19Data})
					ab19ao19Control := ab19ao19Call(ab19ao19T, ab19ao19InputObject)
					ab19ao19Context := ab19ao19InputObject["context"].(map[string]any)
					if ab19ao19Missing {
						delete(ab19ao19Context, ab19ao19Key)
					} else {
						ab19ao19Context[ab19ao19Key] = map[string]any{"primary": 182, "secondary": 21, "operation": 2, "direction": "request", "view": "other", "family": map[string]any{"opaque": []any{0, true}}, "variant": []any{17, "other"}}[ab19ao19Key]
					}
					ab19ao19Output := ab19ao19Call(ab19ao19T, ab19ao19InputObject)
					if ab19ao19Key == "family" || ab19ao19Key == "variant" {
						if !ab19ao19reflect.DeepEqual(ab19ao19Output["fields"], ab19ao19Control["fields"]) {
							ab19ao19T.Fatal("opaque label selected view")
						}
					} else if ab19ao19Output["state"] != "Unknown" || len(ab19ao19Output["fields"].(map[string]any)) != 0 {
						ab19ao19T.Fatal(ab19ao19Output)
					}
				})
			}
		}
	}
	for _, ab19ao19Operation := range []int{0, 2, 5, 9} {
		ab19ao19Output := ab19ao19Call(ab19ao19T, ab19ao19Input(ab19ao19Operation, "pressure", map[string]any{"logical_object": ab19ao19Pressure()}))
		if ab19ao19Output["state"] != "Unknown" {
			ab19ao19T.Fatal(ab19ao19Output)
		}
	}
}
func TestAb19Ao19Presence(ab19ao19T *ab19ao19testing.T) {
	for _, ab19ao19Operation := range []int{6, 7, 8} {
		for _, ab19ao19Mode := range []string{"missing", "empty", "zero"} {
			ab19ao19T.Run(ab19ao19fmt.Sprintf("%d/%s", ab19ao19Operation, ab19ao19Mode), func(ab19ao19T *ab19ao19testing.T) {
				ab19ao19View := map[int]string{6: "generic-result", 7: "progress", 8: "pressure"}[ab19ao19Operation]
				ab19ao19Key := "feedback"
				if ab19ao19Operation == 8 {
					ab19ao19Key = "logical_object"
				}
				ab19ao19Raw := map[string]any{"other": []int{0}}
				if ab19ao19Mode == "empty" {
					ab19ao19Raw[ab19ao19Key] = []int{}
				}
				if ab19ao19Mode == "zero" {
					ab19ao19Raw[ab19ao19Key] = []int{0}
				}
				ab19ao19InputObject := ab19ao19Input(ab19ao19Operation, ab19ao19View, ab19ao19Raw)
				ab19ao19Output := ab19ao19Call(ab19ao19T, ab19ao19InputObject)
				ab19ao19State := "partial"
				if ab19ao19Mode == "zero" && ab19ao19Operation != 8 {
					ab19ao19State = "Hypothesis"
				}
				if ab19ao19Output["state"] != ab19ao19State {
					ab19ao19T.Fatal(ab19ao19Output)
				}
				ab19ao19InputObject["arguments"] = map[string]any{}
				ab19ao19Call(ab19ao19T, ab19ao19InputObject)
			})
		}
	}
}
func TestAb19Ao19Reject(ab19ao19T *ab19ao19testing.T) {
	for ab19ao19Index, ab19ao19InputText := range []string{
		``, `{`, `null`, `[]`, `{} {}`, `{"context":{},"fields":{},"fields":{}}`,
		`{"context":{"family":{"x":0,"x":1}},"fields":{}}`, `{"context":{"variant":[{"x":0,"x":1}]},"fields":{}}`,
		`{"context":{},"fields":{"feedback":[],"feedback":[0]}}`, `{"context":{},"fields":{},"arguments":{"x":1,"x":2}}`,
		`{"context":null,"fields":{}}`, `{"context":{},"fields":null}`, `{"context":[],"fields":{}}`, `{"context":{},"fields":[]}`, `{"context":{},"fields":{},"extra":0}`,
		`{"context":{},"fields":{},"arguments":null}`, `{"context":{},"fields":{},"arguments":[]}`, `{"context":{},"fields":{},"arguments":{"flag":true}}`,
	} {
		ab19ao19T.Run(ab19ao19fmt.Sprintf("structure/%d", ab19ao19Index), func(ab19ao19T *ab19ao19testing.T) {
			ab19ao19Output, ab19ao19Err := Ab19Ao19Interpret([]byte(ab19ao19InputText))
			if ab19ao19Err == nil || ab19ao19Output != nil {
				ab19ao19T.Fatal("accepted invalid structure", ab19ao19InputText)
			}
		})
	}
	for _, ab19ao19Key := range []string{"feedback", "logical_object", "unmodelled"} {
		for _, ab19ao19Bad := range []string{`null`, `0`, `"AA=="`, `{}`, `[null]`, `[true]`, `["1"]`, `[-1]`, `[256]`, `[1.5]`, `[1.0]`, `[1e0]`, `[999999999999999999999999999999]`, `[[0]]`, `[{}]`} {
			ab19ao19InputText := `{"context":{},"fields":{"` + ab19ao19Key + `":` + ab19ao19Bad + `}}`
			ab19ao19Output, ab19ao19Err := Ab19Ao19Interpret([]byte(ab19ao19InputText))
			if ab19ao19Err == nil || ab19ao19Output != nil {
				ab19ao19T.Fatal("accepted invalid bytes", ab19ao19InputText)
			}
		}
	}
	for _, ab19ao19Key := range []string{"primary", "secondary", "operation", "direction", "view"} {
		for _, ab19ao19Bad := range []any{nil, true, []int{}, map[string]any{}} {
			ab19ao19InputObject := ab19ao19Input(8, "pressure", map[string]any{})
			ab19ao19InputObject["context"].(map[string]any)[ab19ao19Key] = ab19ao19Bad
			ab19ao19InputJSON, _ := ab19ao19json.Marshal(ab19ao19InputObject)
			ab19ao19Output, ab19ao19Err := Ab19Ao19Interpret(ab19ao19InputJSON)
			if ab19ao19Err == nil || ab19ao19Output != nil {
				ab19ao19T.Fatal("accepted invalid context", string(ab19ao19InputJSON))
			}
		}
	}
	for _, ab19ao19InputText := range []string{`{"context":{"primary":"181"},"fields":{}}`, `{"context":{"secondary":20.5},"fields":{}}`, `{"context":{"operation":6.0},"fields":{}}`, `{"context":{"direction":0},"fields":{}}`, `{"context":{"view":1},"fields":{}}`} {
		ab19ao19Output, ab19ao19Err := Ab19Ao19Interpret([]byte(ab19ao19InputText))
		if ab19ao19Err == nil || ab19ao19Output != nil {
			ab19ao19T.Fatal("accepted wrong kind")
		}
	}
}
func TestAb19Ao19Isolation(ab19ao19T *ab19ao19testing.T) {
	ab19ao19InputJSON, _ := ab19ao19json.Marshal(ab19ao19Input(8, "pressure", map[string]any{"logical_object": ab19ao19Pressure()}))
	ab19ao19Original := append([]byte(nil), ab19ao19InputJSON...)
	ab19ao19First, ab19ao19Err := Ab19Ao19Interpret(ab19ao19InputJSON)
	if ab19ao19Err != nil {
		ab19ao19T.Fatal(ab19ao19Err)
	}
	ab19ao19Snapshot := append([]byte(nil), ab19ao19First...)
	for ab19ao19Index := range ab19ao19InputJSON {
		ab19ao19InputJSON[ab19ao19Index] = ' '
	}
	if string(ab19ao19First) != string(ab19ao19Snapshot) {
		ab19ao19T.Fatal("input aliases output")
	}
	ab19ao19Second, ab19ao19Err := Ab19Ao19Interpret(ab19ao19Original)
	if ab19ao19Err != nil {
		ab19ao19T.Fatal(ab19ao19Err)
	}
	ab19ao19Decoded := map[string]any{}
	if ab19ao19Err := ab19ao19json.Unmarshal(ab19ao19First, &ab19ao19Decoded); ab19ao19Err != nil {
		ab19ao19T.Fatal(ab19ao19Err)
	}
	ab19ao19Decoded["context"].(map[string]any)["family"].(map[string]any)["opaque"].([]any)[0] = "changed"
	ab19ao19Decoded["raw"].(map[string]any)["logical_object"].([]any)[4] = float64(200)
	if ab19ao19Field(ab19ao19Decoded, "pressure")["raw"].([]any)[0] != float64(1) {
		ab19ao19T.Fatal("raw evidence aliases semantic raw")
	}
	ab19ao19Field(ab19ao19Decoded, "pressure")["raw"].([]any)[0] = float64(201)
	if ab19ao19Decoded["raw"].(map[string]any)["logical_object"].([]any)[4] != float64(200) {
		ab19ao19T.Fatal("semantic raw aliases raw evidence")
	}
	if ab19ao19Field(ab19ao19Decoded, "pressure")["number"] != float64(1) {
		ab19ao19T.Fatal("raw mutation changed calculation")
	}
	for ab19ao19Index := range ab19ao19First {
		ab19ao19First[ab19ao19Index] = ' '
	}
	ab19ao19Third, ab19ao19Err := Ab19Ao19Interpret(ab19ao19Original)
	if ab19ao19Err != nil || string(ab19ao19Second) != string(ab19ao19Snapshot) || string(ab19ao19Third) != string(ab19ao19Snapshot) {
		ab19ao19T.Fatal("output aliases another call")
	}
}
func TestAb19Ao19Fixture(ab19ao19T *ab19ao19testing.T) {
	ab19ao19Data, ab19ao19Err := ab19ao19os.ReadFile("testdata/ab19ao19_b514_menu_views.json")
	if ab19ao19Err != nil {
		ab19ao19T.Fatal(ab19ao19Err)
	}
	ab19ao19Cases := []struct {
		Ab19Ao19ID      string             `json:"id"`
		Ab19Ao19Input   map[string]any     `json:"input"`
		Ab19Ao19State   string             `json:"state"`
		Ab19Ao19Numbers map[string]float64 `json:"numbers"`
	}{}
	if ab19ao19Err := ab19ao19json.Unmarshal(ab19ao19Data, &ab19ao19Cases); ab19ao19Err != nil {
		ab19ao19T.Fatal(ab19ao19Err)
	}
	if len(ab19ao19Cases) == 0 {
		ab19ao19T.Fatal("empty fixture")
	}
	for _, ab19ao19Case := range ab19ao19Cases {
		if !ab19ao19strings.HasPrefix(ab19ao19Case.Ab19Ao19ID, "ab19ao19") {
			ab19ao19T.Fatal("fixture ID")
		}
		ab19ao19Output := ab19ao19Call(ab19ao19T, ab19ao19Case.Ab19Ao19Input)
		if ab19ao19Output["state"] != ab19ao19Case.Ab19Ao19State {
			ab19ao19T.Fatal(ab19ao19Output)
		}
		for ab19ao19Name, ab19ao19Number := range ab19ao19Case.Ab19Ao19Numbers {
			if ab19ao19Field(ab19ao19Output, ab19ao19Name)["number"] != ab19ao19Number {
				ab19ao19T.Fatal(ab19ao19Output)
			}
		}
	}
}

func TestAb19Ao19UnknownSiblings(ab19ao19T *ab19ao19testing.T) {
	for _, ab19ao19Name := range []string{"other", "feedback", "logical_object"} {
		for _, ab19ao19Data := range [][]int{{}, {0}, {255, 1, 2, 3}} {
			ab19ao19Raw := map[string]any{"logical_object": ab19ao19Pressure(), "feedback": []int{1}, "other": []int{0}}
			ab19ao19Raw[ab19ao19Name] = ab19ao19Data
			ab19ao19Output := ab19ao19Call(ab19ao19T, ab19ao19Input(6, "generic-result", ab19ao19Raw))
			if ab19ao19Name != "feedback" && ab19ao19Field(ab19ao19Output, "feedback")["label"] != "in-progress" {
				ab19ao19T.Fatal("unknown sibling erased feedback")
			}
			ab19ao19Output = ab19ao19Call(ab19ao19T, ab19ao19Input(8, "pressure", ab19ao19Raw))
			if ab19ao19Name != "logical_object" && ab19ao19Field(ab19ao19Output, "pressure")["number"] != float64(1) {
				ab19ao19T.Fatal("unknown sibling erased pressure")
			}
		}
	}
	ab19ao19Bad := append([]byte(`{"context":{"opaque":"`), 255)
	ab19ao19Bad = append(ab19ao19Bad, []byte(`"},"fields":{}}`)...)
	ab19ao19Output, ab19ao19Err := Ab19Ao19Interpret(ab19ao19Bad)
	if ab19ao19Err == nil || ab19ao19Output != nil {
		ab19ao19T.Fatal("accepted malformed UTF-8")
	}
}
