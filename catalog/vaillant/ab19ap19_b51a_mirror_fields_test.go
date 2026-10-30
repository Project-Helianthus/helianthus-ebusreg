package vaillant

import (
	ab19ap19json "encoding/json"
	ab19ap19os "os"
	ab19ap19reflect "reflect"
	ab19ap19testing "testing"
)

func ab19ap19Cases(ab19ap19T *ab19ap19testing.T) []map[string]any {
	ab19ap19T.Helper()
	ab19ap19Bytes, ab19ap19Err := ab19ap19os.ReadFile("testdata/ab19ap19_b51a_mirror_fields.json")
	if ab19ap19Err != nil {
		ab19ap19T.Fatal(ab19ap19Err)
	}
	var ab19ap19Cases []map[string]any
	if ab19ap19Err := ab19ap19json.Unmarshal(ab19ap19Bytes, &ab19ap19Cases); ab19ap19Err != nil {
		ab19ap19T.Fatal(ab19ap19Err)
	}
	return ab19ap19Cases
}
func ab19ap19Run(ab19ap19T *ab19ap19testing.T, ab19ap19Input any) map[string]any {
	ab19ap19T.Helper()
	ab19ap19Bytes, _ := ab19ap19json.Marshal(ab19ap19Input)
	ab19ap19Output, ab19ap19Err := Ab19Ap19Evaluate(ab19ap19Bytes)
	if ab19ap19Err != nil {
		ab19ap19T.Fatal(ab19ap19Err)
	}
	var ab19ap19Result map[string]any
	if ab19ap19Err := ab19ap19json.Unmarshal(ab19ap19Output, &ab19ap19Result); ab19ap19Err != nil {
		ab19ap19T.Fatal(ab19ap19Err)
	}
	return ab19ap19Result
}
func ab19ap19Field(ab19ap19Output map[string]any, ab19ap19Name string) map[string]any {
	return ab19ap19Output["fields"].(map[string]any)[ab19ap19Name].(map[string]any)
}
func TestAb19Ap19Vectors(ab19ap19T *ab19ap19testing.T) {
	for _, ab19ap19Case := range ab19ap19Cases(ab19ap19T) {
		ab19ap19T.Run(ab19ap19Case["id"].(string), func(ab19ap19T *ab19ap19testing.T) {
			ab19ap19Input := ab19ap19Case["input"].(map[string]any)
			ab19ap19Output := ab19ap19Run(ab19ap19T, ab19ap19Input)
			if ab19ap19Output["state"] != "complete" || ab19ap19Output["physical"] != "Unknown" || ab19ap19Output["qualified"] != false || !ab19ap19reflect.DeepEqual(ab19ap19Output["context"], ab19ap19Input["context"]) || !ab19ap19reflect.DeepEqual(ab19ap19Output["raw"], ab19ap19Input["fields"]) {
				ab19ap19T.Fatal(ab19ap19Output)
			}
			for ab19ap19Name, ab19ap19Want := range ab19ap19Case["numbers"].(map[string]any) {
				ab19ap19Field := ab19ap19Field(ab19ap19Output, ab19ap19Name)
				if ab19ap19Field["number"] != ab19ap19Want || ab19ap19Field["state"] != "Hypothesis" || ab19ap19Field["calculation"] != "Proven" {
					ab19ap19T.Fatal(ab19ap19Field)
				}
			}
			if ab19ap19Address, ab19ap19OK := ab19ap19Output["fields"].(map[string]any)["address"]; ab19ap19OK && ab19ap19Address.(map[string]any)["logical_validity"] != "Unknown" {
				ab19ap19T.Fatal(ab19ap19Address)
			}
		})
	}
}
func TestAb19Ap19IndependentFields(ab19ap19T *ab19ap19testing.T) {
	for ab19ap19Index := 0; ab19ap19Index < 4; ab19ap19Index++ {
		for ab19ap19Name := range ab19ap19Cases(ab19ap19T)[ab19ap19Index]["input"].(map[string]any)["fields"].(map[string]any) {
			for _, ab19ap19Replacement := range []any{nil, []any{}, []any{float64(0)}, []any{float64(255)}, []any{float64(0), float64(0)}, []any{float64(255), float64(255)}, []any{float64(1), float64(2), float64(3)}} {
				ab19ap19Input := ab19ap19Cases(ab19ap19T)[ab19ap19Index]["input"].(map[string]any)
				ab19ap19Raw := ab19ap19Input["fields"].(map[string]any)
				ab19ap19Original := ab19ap19Run(ab19ap19T, ab19ap19Input)
				if ab19ap19Replacement == nil {
					delete(ab19ap19Raw, ab19ap19Name)
				} else {
					ab19ap19Raw[ab19ap19Name] = ab19ap19Replacement
				}
				ab19ap19Output := ab19ap19Run(ab19ap19T, ab19ap19Input)
				for ab19ap19Sibling, ab19ap19Expected := range ab19ap19Original["fields"].(map[string]any) {
					if ab19ap19Sibling != ab19ap19Name && !ab19ap19reflect.DeepEqual(ab19ap19Expected, ab19ap19Output["fields"].(map[string]any)[ab19ap19Sibling]) {
						ab19ap19T.Fatal("sibling changed", ab19ap19Sibling)
					}
				}
				ab19ap19Result := ab19ap19Field(ab19ap19Output, ab19ap19Name)
				if ab19ap19Replacement == nil {
					if _, ab19ap19Present := ab19ap19Result["raw"]; ab19ap19Present {
						ab19ap19T.Fatal("invented missing raw")
					}
				} else if !ab19ap19reflect.DeepEqual(ab19ap19Result["raw"], ab19ap19Replacement) {
					ab19ap19T.Fatal("changed field raw lost", ab19ap19Name)
				}
				ab19ap19Width := 1
				if ab19ap19Name == "address" || ab19ap19Name == "mirror_address" || ab19ap19Name == "symbol_address" {
					ab19ap19Width = 2
				}
				ab19ap19Incomplete := ab19ap19Replacement == nil
				if ab19ap19Replacement != nil && ab19ap19Name != "data" {
					ab19ap19Incomplete = len(ab19ap19Replacement.([]any)) != ab19ap19Width
				}
				if ab19ap19Incomplete {
					if ab19ap19Result["state"] != "incomplete" || ab19ap19Result["number"] != nil || ab19ap19Output["state"] != "partial" {
						ab19ap19T.Fatal(ab19ap19Output)
					}
				} else {
					if ab19ap19Result["state"] != "Hypothesis" {
						ab19ap19T.Fatal(ab19ap19Result)
					}
					if ab19ap19Name == "data" {
						if ab19ap19Result["label"] != "opaque" || ab19ap19Result["number"] != nil {
							ab19ap19T.Fatal(ab19ap19Result)
						}
					} else if ab19ap19Width == 1 && ab19ap19Result["number"] != ab19ap19Replacement.([]any)[0] {
						ab19ap19T.Fatal("scalar value changed", ab19ap19Name, ab19ap19Result)
					}
				}
			}
		}
	}
}
func TestAb19Ap19StatusAndFlags(ab19ap19T *ab19ap19testing.T) {
	for ab19ap19Status := 0; ab19ap19Status < 256; ab19ap19Status++ {
		ab19ap19Input := ab19ap19Cases(ab19ap19T)[1]["input"].(map[string]any)
		ab19ap19Input["fields"].(map[string]any)["status"] = []int{ab19ap19Status}
		ab19ap19Output := ab19ap19Run(ab19ap19T, ab19ap19Input)
		ab19ap19Alternatives := ab19ap19Field(ab19ap19Output, "status")["alternatives"].(map[string]any)
		if len(ab19ap19Alternatives) != 2 {
			ab19ap19T.Fatal("both status interpretations required", ab19ap19Alternatives)
		}
		ab19ap19Boolean := ab19ap19Alternatives["whole_boolean"].(map[string]any)
		ab19ap19Bit := ab19ap19Alternatives["bit0_out_of_area"].(map[string]any)
		if len(ab19ap19Alternatives) != 2 || ab19ap19Bit["value"] != (ab19ap19Status%2 == 1) {
			ab19ap19T.Fatal(ab19ap19Alternatives)
		}
		if ab19ap19Status < 2 {
			if ab19ap19Boolean["value"] != (ab19ap19Status == 1) {
				ab19ap19T.Fatal(ab19ap19Boolean)
			}
		} else if ab19ap19Boolean["state"] != "Unknown" || ab19ap19Boolean["value"] != nil {
			ab19ap19T.Fatal(ab19ap19Boolean)
		}
		if len(ab19ap19Field(ab19ap19Output, "data")["raw"].([]any)) != 2 {
			ab19ap19T.Fatal("returned data resized")
		}
	}
	for _, ab19ap19Index := range []int{0, 2, 3} {
		ab19ap19Names := map[int][]string{0: {"triple-storage"}, 2: {"triple-storage", "direct-update"}, 3: {"not-written", "invalid-symbol", "invalid-parameter", "storage-error"}}[ab19ap19Index]
		for ab19ap19Bit := -1; ab19ap19Bit < 8; ab19ap19Bit++ {
			ab19ap19Number := 0
			if ab19ap19Bit >= 0 {
				ab19ap19Number = 1 << ab19ap19Bit
			}
			ab19ap19Input := ab19ap19Cases(ab19ap19T)[ab19ap19Index]["input"].(map[string]any)
			ab19ap19Input["fields"].(map[string]any)["flags"] = []int{ab19ap19Number}
			ab19ap19Result := ab19ap19Field(ab19ap19Run(ab19ap19T, ab19ap19Input), "flags")
			for ab19ap19Index, ab19ap19Name := range ab19ap19Names {
				if ab19ap19Result["flags"].(map[string]any)[ab19ap19Name] != (ab19ap19Index == ab19ap19Bit) {
					ab19ap19T.Fatal(ab19ap19Result)
				}
			}
			ab19ap19Reserved := 0
			if ab19ap19Bit >= len(ab19ap19Names) {
				ab19ap19Reserved = ab19ap19Number
			}
			if ab19ap19Result["reserved"] != float64(ab19ap19Reserved) {
				ab19ap19T.Fatal(ab19ap19Result)
			}
		}
	}
}
func TestAb19Ap19ContextAndOwnership(ab19ap19T *ab19ap19testing.T) {
	for _, ab19ap19Key := range []string{"primary", "secondary", "operation", "direction", "view", "family", "variant", "opaque"} {
		for _, ab19ap19Value := range []any{nil, false, float64(0), "other", map[string]any{"nested": []any{false, nil, float64(0)}}} {
			ab19ap19Input := ab19ap19Cases(ab19ap19T)[0]["input"].(map[string]any)
			ab19ap19Context := ab19ap19Input["context"].(map[string]any)
			if ab19ap19Key == "primary" || ab19ap19Key == "secondary" || ab19ap19Key == "operation" {
				if ab19ap19Value != nil && ab19ap19Value != float64(0) {
					continue
				}
			}
			if ab19ap19Key == "direction" || ab19ap19Key == "view" {
				if ab19ap19Value != nil && ab19ap19Value != "other" {
					continue
				}
			}
			ab19ap19Context[ab19ap19Key] = ab19ap19Value
			ab19ap19Input["arguments"] = map[string]any{"opaque": []any{false, nil, float64(0), "", []any{}}}
			ab19ap19Output := ab19ap19Run(ab19ap19T, ab19ap19Input)
			if !ab19ap19reflect.DeepEqual(ab19ap19Output["context"], ab19ap19Context) || !ab19ap19reflect.DeepEqual(ab19ap19Output["arguments"], ab19ap19Input["arguments"]) || ab19ap19Output["qualified"] != false {
				ab19ap19T.Fatal(ab19ap19Output)
			}
			ab19ap19Unsupported := ab19ap19Key != "family" && ab19ap19Key != "variant" && ab19ap19Key != "opaque"
			if ab19ap19Unsupported {
				if ab19ap19Output["state"] != "Unknown" || len(ab19ap19Output["fields"].(map[string]any)) != 0 {
					ab19ap19T.Fatal(ab19ap19Output)
				}
			} else if ab19ap19Output["state"] != "complete" {
				ab19ap19T.Fatal(ab19ap19Output)
			}
			delete(ab19ap19Context, ab19ap19Key)
			ab19ap19Output = ab19ap19Run(ab19ap19T, ab19ap19Input)
			if _, ab19ap19Present := ab19ap19Output["context"].(map[string]any)[ab19ap19Key]; ab19ap19Present {
				ab19ap19T.Fatal("missing context invented")
			}
		}
	}
	ab19ap19Input := ab19ap19Cases(ab19ap19T)[0]["input"].(map[string]any)
	ab19ap19Input["context"].(map[string]any)["register"] = float64(6)
	ab19ap19Input["arguments"] = map[string]any{"register": float64(6), "sequence": float64(255), "length": float64(3)}
	ab19ap19Output := ab19ap19Run(ab19ap19T, ab19ap19Input)
	if ab19ap19Output["state"] != "complete" || len(ab19ap19Output["fields"].(map[string]any)) != 4 || ab19ap19Field(ab19ap19Output, "sequence")["number"] != float64(0) {
		ab19ap19T.Fatal("selector alias or sequence correlation", ab19ap19Output)
	}
	ab19ap19Bytes, _ := ab19ap19json.Marshal(ab19ap19Input)
	ab19ap19Result, ab19ap19Err := Ab19Ap19Evaluate(ab19ap19Bytes)
	if ab19ap19Err != nil {
		ab19ap19T.Fatal(ab19ap19Err)
	}
	ab19ap19Saved := append([]byte(nil), ab19ap19Result...)
	for ab19ap19Index := range ab19ap19Bytes {
		ab19ap19Bytes[ab19ap19Index] = ' '
	}
	if string(ab19ap19Saved) != string(ab19ap19Result) {
		ab19ap19T.Fatal("input aliases result")
	}
	for ab19ap19Index := range ab19ap19Result {
		ab19ap19Result[ab19ap19Index] = ' '
	}
	ab19ap19Again, _ := ab19ap19json.Marshal(ab19ap19Input)
	ab19ap19Result, ab19ap19Err = Ab19Ap19Evaluate(ab19ap19Again)
	if ab19ap19Err != nil || string(ab19ap19Result) != string(ab19ap19Saved) {
		ab19ap19T.Fatal("result aliases later call")
	}
	ab19ap19Output["raw"].(map[string]any)["sequence"].([]any)[0] = float64(99)
	if ab19ap19Field(ab19ap19Output, "sequence")["raw"].([]any)[0] != float64(0) || ab19ap19Run(ab19ap19T, ab19ap19Input)["raw"].(map[string]any)["sequence"].([]any)[0] != float64(0) {
		ab19ap19T.Fatal("nested storage aliases")
	}
}
func TestAb19Ap19Reject(ab19ap19T *ab19ap19testing.T) {
	for _, ab19ap19Input := range []string{``, `{`, `[]`, `{} {}`, `{"context":{},"fields":{},"fields":{}}`, `{"context":{"x":{"a":1,"a":2}},"fields":{}}`, `{"context":{},"fields":{},"arguments":{"x":[{"a":0,"a":1}]}}`, `{"context":null,"fields":{}}`, `{"context":{},"fields":null}`, `{"context":{},"fields":{},"arguments":false}`} {
		ab19ap19Output, ab19ap19Err := Ab19Ap19Evaluate([]byte(ab19ap19Input))
		if ab19ap19Err == nil || ab19ap19Output != nil {
			ab19ap19T.Fatal("accepted invalid JSON", ab19ap19Input)
		}
	}
	for _, ab19ap19Invalid := range []any{nil, false, "0", float64(-1), float64(256), float64(0.5), ab19ap19json.Number("0.00000000000000000001"), ab19ap19json.Number("255.00000000000000000001"), map[string]any{}} {
		for _, ab19ap19Key := range []string{"sequence", "flags", "length", "address", "mirror_address", "symbol_address", "status", "data"} {
			ab19ap19Input := ab19ap19Cases(ab19ap19T)[0]["input"].(map[string]any)
			ab19ap19Input["fields"].(map[string]any)[ab19ap19Key] = []any{ab19ap19Invalid}
			ab19ap19Bytes, _ := ab19ap19json.Marshal(ab19ap19Input)
			ab19ap19Output, ab19ap19Err := Ab19Ap19Evaluate(ab19ap19Bytes)
			if ab19ap19Err == nil || ab19ap19Output != nil {
				ab19ap19T.Fatal("accepted invalid byte", ab19ap19Key, ab19ap19Invalid)
			}
		}
	}
	for _, ab19ap19Key := range []string{"primary", "secondary", "operation", "direction", "view"} {
		ab19ap19Input := ab19ap19Cases(ab19ap19T)[0]["input"].(map[string]any)
		ab19ap19Input["context"].(map[string]any)[ab19ap19Key] = []int{}
		ab19ap19Bytes, _ := ab19ap19json.Marshal(ab19ap19Input)
		if _, ab19ap19Err := Ab19Ap19Evaluate(ab19ap19Bytes); ab19ap19Err == nil {
			ab19ap19T.Fatal("accepted wrong context kind")
		}
	}
}

func TestAb19Ap19IndependentCoordinates(ab19ap19T *ab19ap19testing.T) {
	for ab19ap19Index := 0; ab19ap19Index < 4; ab19ap19Index++ {
		for _, ab19ap19Sequence := range []float64{0, 255} {
			ab19ap19Input := ab19ap19Cases(ab19ap19T)[ab19ap19Index]["input"].(map[string]any)
			ab19ap19Input["fields"].(map[string]any)["sequence"] = []float64{ab19ap19Sequence}
			ab19ap19Input["arguments"] = map[string]any{"sequence": 255 - ab19ap19Sequence, "length": 3, "register": 6}
			ab19ap19Output := ab19ap19Run(ab19ap19T, ab19ap19Input)
			if ab19ap19Output["state"] != "complete" || ab19ap19Field(ab19ap19Output, "sequence")["number"] != ab19ap19Sequence || !ab19ap19reflect.DeepEqual(ab19ap19Output["context"], ab19ap19Input["context"]) {
				ab19ap19T.Fatal("sequence correlation", ab19ap19Output)
			}
			if ab19ap19Index == 1 && len(ab19ap19Field(ab19ap19Output, "data")["raw"].([]any)) != 2 {
				ab19ap19T.Fatal("requested length resized data")
			}
		}
	}
	for _, ab19ap19Name := range []string{"address", "mirror_address", "symbol_address"} {
		for _, ab19ap19Pair := range []struct {
			ab19ap19Bytes []int
			ab19ap19Want  float64
		}{{[]int{0, 0}, 0}, {[]int{255, 255}, 65535}, {[]int{52, 18}, 4660}, {[]int{120, 86}, 22136}} {
			ab19ap19Index := 2
			if ab19ap19Name == "address" {
				ab19ap19Index = 0
			}
			ab19ap19Input := ab19ap19Cases(ab19ap19T)[ab19ap19Index]["input"].(map[string]any)
			ab19ap19Input["fields"].(map[string]any)[ab19ap19Name] = ab19ap19Pair.ab19ap19Bytes
			ab19ap19Result := ab19ap19Field(ab19ap19Run(ab19ap19T, ab19ap19Input), ab19ap19Name)
			if ab19ap19Result["number"] != ab19ap19Pair.ab19ap19Want || ab19ap19Result["logical_validity"] != "Unknown" || ab19ap19Result["label"] != nil {
				ab19ap19T.Fatal(ab19ap19Result)
			}
		}
	}
	for _, ab19ap19Argument := range []any{map[string]any{}, map[string]any{"opaque": false}, map[string]any{"opaque": 0}, map[string]any{"opaque": []any{}}, map[string]any{"opaque": nil}} {
		ab19ap19Input := ab19ap19Cases(ab19ap19T)[0]["input"].(map[string]any)
		ab19ap19Input["arguments"] = ab19ap19Argument
		ab19ap19Output := ab19ap19Run(ab19ap19T, ab19ap19Input)
		ab19ap19Bytes, _ := ab19ap19json.Marshal(ab19ap19Argument)
		var ab19ap19Want any
		_ = ab19ap19json.Unmarshal(ab19ap19Bytes, &ab19ap19Want)
		if !ab19ap19reflect.DeepEqual(ab19ap19Want, ab19ap19Output["arguments"]) || ab19ap19Output["state"] != "complete" {
			ab19ap19T.Fatal(ab19ap19Output)
		}
	}
}

func TestAb19Ap19RawBoundaries(ab19ap19T *ab19ap19testing.T) {
	for _, ab19ap19Name := range []string{"sequence", "flags", "length", "address", "mirror_address", "symbol_address", "status", "data"} {
		for _, ab19ap19Value := range []any{nil, false, 0, "", map[string]any{}} {
			ab19ap19Input := ab19ap19Cases(ab19ap19T)[0]["input"].(map[string]any)
			ab19ap19Input["fields"].(map[string]any)[ab19ap19Name] = ab19ap19Value
			ab19ap19Bytes, _ := ab19ap19json.Marshal(ab19ap19Input)
			if ab19ap19Output, ab19ap19Err := Ab19Ap19Evaluate(ab19ap19Bytes); ab19ap19Err == nil || ab19ap19Output != nil {
				ab19ap19T.Fatal("field must be array", ab19ap19Name)
			}
		}
	}
	for _, ab19ap19Value := range []any{[]any{}, []int{0}, []int{255}} {
		ab19ap19Input := ab19ap19Cases(ab19ap19T)[0]["input"].(map[string]any)
		ab19ap19Input["fields"].(map[string]any)["configured_register"] = ab19ap19Value
		ab19ap19Output := ab19ap19Run(ab19ap19T, ab19ap19Input)
		if len(ab19ap19Output["fields"].(map[string]any)) != 4 || ab19ap19Output["state"] != "complete" {
			ab19ap19T.Fatal("unknown field interpreted", ab19ap19Output)
		}
		ab19ap19Bytes, _ := ab19ap19json.Marshal(ab19ap19Value)
		var ab19ap19Want any
		_ = ab19ap19json.Unmarshal(ab19ap19Bytes, &ab19ap19Want)
		if !ab19ap19reflect.DeepEqual(ab19ap19Output["raw"].(map[string]any)["configured_register"], ab19ap19Want) {
			ab19ap19T.Fatal("unknown raw lost")
		}
	}
}
