package vaillant

import (
	"bytes"
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"
)

func ab19aq19EvaluateTest(ab19aq19T *testing.T, ab19aq19Input any) map[string]any {
	ab19aq19T.Helper()
	ab19aq19JSON, ab19aq19Err := json.Marshal(ab19aq19Input)
	if ab19aq19Err != nil {
		ab19aq19T.Fatal(ab19aq19Err)
	}
	ab19aq19Output, ab19aq19Err := Ab19Aq19Evaluate(ab19aq19JSON)
	if ab19aq19Err != nil {
		ab19aq19T.Fatal(ab19aq19Err)
	}
	ab19aq19Result := map[string]any{}
	if ab19aq19Err := json.Unmarshal(ab19aq19Output, &ab19aq19Result); ab19aq19Err != nil {
		ab19aq19T.Fatal(ab19aq19Err)
	}
	return ab19aq19Result
}

func ab19aq19Cases(ab19aq19T *testing.T) []map[string]any {
	ab19aq19T.Helper()
	ab19aq19Data, ab19aq19Err := os.ReadFile("testdata/ab19aq19_b51a_handshake_fields.json")
	if ab19aq19Err != nil {
		ab19aq19T.Fatal(ab19aq19Err)
	}
	var ab19aq19Cases []map[string]any
	if ab19aq19Err := json.Unmarshal(ab19aq19Data, &ab19aq19Cases); ab19aq19Err != nil {
		ab19aq19T.Fatal(ab19aq19Err)
	}
	return ab19aq19Cases
}

func ab19aq19Path(ab19aq19Value any, ab19aq19Path string) any {
	for _, ab19aq19Key := range strings.Split(ab19aq19Path, ".") {
		ab19aq19Map, ab19aq19OK := ab19aq19Value.(map[string]any)
		if !ab19aq19OK {
			return nil
		}
		ab19aq19Value = ab19aq19Map[ab19aq19Key]
	}
	return ab19aq19Value
}

func TestAb19Aq19Fixtures(ab19aq19T *testing.T) {
	for _, ab19aq19Case := range ab19aq19Cases(ab19aq19T) {
		ab19aq19T.Run(ab19aq19Case["id"].(string), func(ab19aq19T *testing.T) {
			ab19aq19Input := ab19aq19Case["input"].(map[string]any)
			ab19aq19Output := ab19aq19EvaluateTest(ab19aq19T, ab19aq19Input)
			for ab19aq19PathKey, ab19aq19Expected := range ab19aq19Case["expected"].(map[string]any) {
				if ab19aq19Got := ab19aq19Path(ab19aq19Output, ab19aq19PathKey); !reflect.DeepEqual(ab19aq19Got, ab19aq19Expected) {
					ab19aq19T.Errorf("%s = %v want %v", ab19aq19PathKey, ab19aq19Got, ab19aq19Expected)
				}
			}
			for _, ab19aq19Key := range []string{"context", "fields", "arguments", "extra"} {
				if !reflect.DeepEqual(ab19aq19Output[ab19aq19Key], ab19aq19Input[ab19aq19Key]) {
					ab19aq19T.Fatal("raw retention", ab19aq19Key)
				}
			}
			for ab19aq19Key, ab19aq19Raw := range ab19aq19Input["fields"].(map[string]any) {
				if !reflect.DeepEqual(ab19aq19Path(ab19aq19Output, "results."+ab19aq19Key+".raw"), ab19aq19Raw) {
					ab19aq19T.Fatal("result raw", ab19aq19Key)
				}
			}
			if ab19aq19Output["physical"] != "Unknown" || ab19aq19Output["qualified"] != false {
				ab19aq19T.Fatal("qualification")
			}
			ab19aq19JSON, _ := json.Marshal(ab19aq19Output)
			for _, ab19aq19Forbidden := range []string{"success", "crc-valid", "validated"} {
				if bytes.Contains(ab19aq19JSON, []byte(ab19aq19Forbidden)) {
					ab19aq19T.Fatal("invented claim", ab19aq19Forbidden)
				}
			}
			if ab19aq19Path(ab19aq19Output, "results.crc.grade") != nil && ab19aq19Path(ab19aq19Output, "results.crc.grade") != "Hypothesis" {
				ab19aq19T.Fatal("calculation promoted meaning")
			}
		})
	}
}

func TestAb19Aq19Context(ab19aq19T *testing.T) {
	for _, ab19aq19Case := range ab19aq19Cases(ab19aq19T) {
		for _, ab19aq19Key := range []string{"primary", "secondary", "operation", "direction", "view", "family", "variant"} {
			ab19aq19T.Run(ab19aq19Case["id"].(string)+"/"+ab19aq19Key, func(ab19aq19T *testing.T) {
				ab19aq19JSON, _ := json.Marshal(ab19aq19Case["input"])
				var ab19aq19Input map[string]any
				if ab19aq19Err := json.Unmarshal(ab19aq19JSON, &ab19aq19Input); ab19aq19Err != nil {
					ab19aq19T.Fatal(ab19aq19Err)
				}
				ab19aq19Context := ab19aq19Input["context"].(map[string]any)
				switch ab19aq19Key {
				case "primary", "secondary":
					ab19aq19Context[ab19aq19Key] = 0
				case "operation":
					ab19aq19Context[ab19aq19Key] = 0
				default:
					ab19aq19Context[ab19aq19Key] = "other"
				}
				ab19aq19Output := ab19aq19EvaluateTest(ab19aq19T, ab19aq19Input)
				if ab19aq19Key == "family" || ab19aq19Key == "variant" {
					if ab19aq19Output["qualified"] != false || ab19aq19Output["state"] != "complete" {
						ab19aq19T.Fatal("opaque identity")
					}
					return
				}
				if ab19aq19Output["state"] != "Unknown" {
					ab19aq19T.Fatal("unsupported context")
				}
				for _, ab19aq19Result := range ab19aq19Output["results"].(map[string]any) {
					if len(ab19aq19Result.(map[string]any)) != 2 {
						ab19aq19T.Fatal("formula leaked")
					}
				}
			})
		}
	}
}

func TestAb19Aq19Invalid(ab19aq19T *testing.T) {
	for _, ab19aq19Input := range []string{
		`{`, `[]`, `{} {}`, `{"context":{},"context":{},"fields":{}}`,
		`{"context":{},"fields":{},"extra":[{"x":0,"x":1}]}`,
		`{"context":{"operation":null},"fields":{}}`,
		`{"context":{"primary":true},"fields":{}}`,
		`{"context":{"secondary":1.5},"fields":{}}`,
		`{"context":{"operation":256},"fields":{}}`,
		`{"context":{"direction":null},"fields":{}}`,
		`{"context":{"view":false},"fields":{}}`,
		`{"context":{},"fields":{"sequence":null}}`,
		`{"context":{},"fields":{"sequence":["1"]}}`,
		`{"context":{},"fields":{"sequence":[false]}}`,
		`{"context":{},"fields":{"sequence":[null]}}`,
		`{"context":{},"fields":{"sequence":[-1]}}`,
		`{"context":{},"fields":{"sequence":[256]}}`,
		`{"context":{},"fields":{"sequence":[1.5]}}`,
		`{"context":{},"fields":{},"arguments":null}`,
		`{"context":{},"fields":{},"arguments":{"value":1}}`,
		`{"context":{},"fields":{},"extra":1e999}`,
		`{"context":{},"fields":{}} true`,
	} {
		ab19aq19Output, ab19aq19Err := Ab19Aq19Evaluate(json.RawMessage(ab19aq19Input))
		if ab19aq19Err == nil || ab19aq19Output != nil {
			ab19aq19T.Errorf("accepted %s", ab19aq19Input)
		}
	}
}

func TestAb19Aq19WidthsAndFields(ab19aq19T *testing.T) {
	for _, ab19aq19Case := range ab19aq19Cases(ab19aq19T) {
		for ab19aq19Key := range ab19aq19Case["input"].(map[string]any)["fields"].(map[string]any) {
			for _, ab19aq19Width := range []int{-1, 0, 1, 2, 4} {
				ab19aq19JSON, _ := json.Marshal(ab19aq19Case["input"])
				var ab19aq19Input map[string]any
				if ab19aq19Err := json.Unmarshal(ab19aq19JSON, &ab19aq19Input); ab19aq19Err != nil {
					ab19aq19T.Fatal(ab19aq19Err)
				}
				ab19aq19Fields := ab19aq19Input["fields"].(map[string]any)
				ab19aq19Data := []any{}
				for ab19aq19Index := 0; ab19aq19Index < ab19aq19Width; ab19aq19Index++ {
					ab19aq19Data = append(ab19aq19Data, float64(7+ab19aq19Index))
				}
				if ab19aq19Width < 0 {
					delete(ab19aq19Fields, ab19aq19Key)
				} else {
					ab19aq19Fields[ab19aq19Key] = ab19aq19Data
				}
				ab19aq19Output := ab19aq19EvaluateTest(ab19aq19T, ab19aq19Input)
				ab19aq19Result := ab19aq19Path(ab19aq19Output, "results."+ab19aq19Key).(map[string]any)
				ab19aq19Opaque := ab19aq19Key == "dsn" || ab19aq19Key == "start" || ab19aq19Key == "offset"
				if ab19aq19Opaque {
					if ab19aq19Width >= 0 && !reflect.DeepEqual(ab19aq19Result["raw"], ab19aq19Data) {
						ab19aq19T.Fatal("opaque lost")
					}
					if ab19aq19Result["state"] == "incomplete" || ab19aq19Output["state"] != "complete" {
						ab19aq19T.Fatal("opaque width imposed")
					}
				} else {
					ab19aq19Required := 1
					if ab19aq19Key == "crc" {
						ab19aq19Required = 2
					}
					if ab19aq19Width != ab19aq19Required {
						if ab19aq19Output["state"] != "partial" || ab19aq19Result["state"] != "incomplete" || ab19aq19Result["number"] != nil || ab19aq19Result["flags"] != nil {
							ab19aq19T.Fatal("incomplete conversion", ab19aq19Key, ab19aq19Width)
						}
					} else if ab19aq19Result["state"] != "complete" {
						ab19aq19T.Fatal("valid width")
					}
					if ab19aq19Key != "sequence" && ab19aq19Fields["sequence"] != nil && ab19aq19Path(ab19aq19Output, "results.sequence.state") != "complete" {
						ab19aq19T.Fatal("sibling lost")
					}
				}
				if ab19aq19Width < 0 {
					if _, ab19aq19Present := ab19aq19Result["raw"]; ab19aq19Present {
						ab19aq19T.Fatal("missing became null")
					}
				}
			}
		}
	}
}

func TestAb19Aq19Memory(ab19aq19T *testing.T) {
	ab19aq19JSON, _ := json.Marshal(ab19aq19Cases(ab19aq19T)[3]["input"])
	ab19aq19Original := append([]byte(nil), ab19aq19JSON...)
	ab19aq19First, ab19aq19Err := Ab19Aq19Evaluate(ab19aq19JSON)
	if ab19aq19Err != nil {
		ab19aq19T.Fatal(ab19aq19Err)
	}
	ab19aq19Saved := append([]byte(nil), ab19aq19First...)
	for ab19aq19Index := range ab19aq19JSON {
		ab19aq19JSON[ab19aq19Index] = ' '
	}
	if !bytes.Equal(ab19aq19First, ab19aq19Saved) {
		ab19aq19T.Fatal("input alias")
	}
	for ab19aq19Index := range ab19aq19First {
		ab19aq19First[ab19aq19Index] = ' '
	}
	ab19aq19Second, ab19aq19Err := Ab19Aq19Evaluate(ab19aq19Original)
	if ab19aq19Err != nil || !bytes.Equal(ab19aq19Second, ab19aq19Saved) {
		ab19aq19T.Fatal("cross call alias")
	}
}

func TestAb19Aq19IndependentBitsAndNumbers(ab19aq19T *testing.T) {
	for _, ab19aq19Direction := range []string{"request", "positive"} {
		ab19aq19Labels := []string{"power-on/start", "change", "CRC-support"}
		if ab19aq19Direction == "positive" {
			ab19aq19Labels = []string{"valid", "invalid", "blocked", "change", "denied", "invalid-range", "not-user-adjustable"}
		}
		for ab19aq19Bit := 0; ab19aq19Bit < 8; ab19aq19Bit++ {
			ab19aq19Input := map[string]any{"context": map[string]any{"primary": 181, "secondary": 26, "operation": 1, "view": "handshake", "direction": ab19aq19Direction}, "fields": map[string]any{"sequence": []int{13}, "flags": []int{1 << ab19aq19Bit}, "crc": []int{1, 2}, "family_code": []int{27}}}
			ab19aq19Output := ab19aq19EvaluateTest(ab19aq19T, ab19aq19Input)
			for ab19aq19Index, ab19aq19Label := range ab19aq19Labels {
				if ab19aq19Path(ab19aq19Output, "results.flags.flags."+ab19aq19Label) != (ab19aq19Index == ab19aq19Bit) {
					ab19aq19T.Fatal("wrong bit", ab19aq19Direction, ab19aq19Label, ab19aq19Bit)
				}
			}
			ab19aq19Reserved := 0
			if ab19aq19Bit >= len(ab19aq19Labels) {
				ab19aq19Reserved = 1 << ab19aq19Bit
			}
			if ab19aq19Path(ab19aq19Output, "results.flags.reserved") != float64(ab19aq19Reserved) {
				ab19aq19T.Fatal("reserved")
			}
			if ab19aq19Path(ab19aq19Output, "results.sequence.number") != float64(13) {
				ab19aq19T.Fatal("sequence")
			}
			if ab19aq19Direction == "request" && ab19aq19Path(ab19aq19Output, "results.crc.number") != float64(513) {
				ab19aq19T.Fatal("CRC word")
			}
			if ab19aq19Direction == "positive" && ab19aq19Path(ab19aq19Output, "results.family_code.number") != float64(27) {
				ab19aq19T.Fatal("family code")
			}
		}
	}
	for ab19aq19Bit := 0; ab19aq19Bit < 8; ab19aq19Bit++ {
		ab19aq19Input := map[string]any{"context": map[string]any{"primary": 181, "secondary": 26, "view": "reset", "direction": "positive"}, "fields": map[string]any{"sequence": []int{0}, "ack": []int{1 << ab19aq19Bit}}}
		ab19aq19Output := ab19aq19EvaluateTest(ab19aq19T, ab19aq19Input)
		if ab19aq19Path(ab19aq19Output, "results.ack.flags.denied") != (ab19aq19Bit == 0) {
			ab19aq19T.Fatal("ack bit")
		}
		ab19aq19Expected := 0
		if ab19aq19Bit > 0 {
			ab19aq19Expected = 1 << ab19aq19Bit
		}
		if ab19aq19Path(ab19aq19Output, "results.ack.unmodeled") != float64(ab19aq19Expected) {
			ab19aq19T.Fatal("ack unmodeled")
		}
	}
}

func TestAb19Aq19AbsentIdentityAndOpaqueValues(ab19aq19T *testing.T) {
	for _, ab19aq19Key := range []string{"primary", "secondary", "operation", "view", "direction"} {
		ab19aq19Input := ab19aq19Cases(ab19aq19T)[3]["input"].(map[string]any)
		delete(ab19aq19Input["context"].(map[string]any), ab19aq19Key)
		if ab19aq19EvaluateTest(ab19aq19T, ab19aq19Input)["state"] != "Unknown" {
			ab19aq19T.Fatal("missing context admitted", ab19aq19Key)
		}
	}
	for _, ab19aq19Key := range []string{"family", "variant"} {
		for _, ab19aq19Value := range []any{nil, false, float64(0), "", []any{}} {
			ab19aq19Input := ab19aq19Cases(ab19aq19T)[0]["input"].(map[string]any)
			ab19aq19Input["context"].(map[string]any)[ab19aq19Key] = ab19aq19Value
			ab19aq19Output := ab19aq19EvaluateTest(ab19aq19T, ab19aq19Input)
			if !reflect.DeepEqual(ab19aq19Output["context"], ab19aq19Input["context"]) || ab19aq19Output["qualified"] != false {
				ab19aq19T.Fatal("opaque value")
			}
		}
	}
}
