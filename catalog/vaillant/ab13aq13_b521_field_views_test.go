package vaillant

import (
	ab13aq13bytes "bytes"
	ab13aq13json "encoding/json"
	ab13aq13os "os"
	ab13aq13reflect "reflect"
	ab13aq13strconv "strconv"
	ab13aq13strings "strings"
	ab13aq13testing "testing"
)

func ab13aq13Input(ab13aq13T *ab13aq13testing.T, ab13aq13Direction string) map[string]any {
	ab13aq13T.Helper()
	ab13aq13Fields := map[string]any{"status0": []any{0}, "status1": []any{0}, "brine_percent": []any{0}, "fan_percent": []any{0}, "fan_speed": []any{0}, "temperature0": []any{0, 0}, "temperature1": []any{0, 0}}
	if ab13aq13Direction == "request" {
		ab13aq13Fields = map[string]any{"flags": []any{0}, "appliance_bits": []any{0}, "noise_percent": []any{0}, "frost_temperature": []any{0}}
	}
	return map[string]any{"context": map[string]any{"primary": 181, "secondary": 33, "view": "module-fields", "direction": ab13aq13Direction}, "fields": ab13aq13Fields}
}

func ab13aq13Evaluate(ab13aq13T *ab13aq13testing.T, ab13aq13Input any) map[string]any {
	ab13aq13T.Helper()
	ab13aq13Encoded, ab13aq13Error := ab13aq13json.Marshal(ab13aq13Input)
	if ab13aq13Error != nil {
		ab13aq13T.Fatal(ab13aq13Error)
	}
	ab13aq13Raw, ab13aq13Error := Ab13Aq13Evaluate(ab13aq13Encoded)
	if ab13aq13Error != nil {
		ab13aq13T.Fatal(ab13aq13Error)
	}
	ab13aq13Output := map[string]any{}
	if ab13aq13Error := ab13aq13json.Unmarshal(ab13aq13Raw, &ab13aq13Output); ab13aq13Error != nil {
		ab13aq13T.Fatal(ab13aq13Error)
	}
	if ab13aq13Output["physical"] != "Unknown" || ab13aq13Output["qualified"] != false {
		ab13aq13T.Fatal("qualification changed", ab13aq13Output)
	}
	return ab13aq13Output
}

func ab13aq13Results(ab13aq13Output map[string]any) map[string]any {
	return ab13aq13Output["results"].(map[string]any)
}

func TestAb13Aq13Fixture(ab13aq13T *ab13aq13testing.T) {
	ab13aq13Raw, ab13aq13Error := ab13aq13os.ReadFile("testdata/ab13aq13_b521_field_views.json")
	if ab13aq13Error != nil {
		ab13aq13T.Fatal(ab13aq13Error)
	}
	var ab13aq13Cases []struct {
		Ab13Aq13ID       string                    `json:"id"`
		Ab13Aq13Input    map[string]any            `json:"input"`
		Ab13Aq13State    string                    `json:"state"`
		Ab13Aq13Expected map[string]map[string]any `json:"expected"`
	}
	if ab13aq13Error := ab13aq13json.Unmarshal(ab13aq13Raw, &ab13aq13Cases); ab13aq13Error != nil {
		ab13aq13T.Fatal(ab13aq13Error)
	}
	if len(ab13aq13Cases) != 41 {
		ab13aq13T.Fatalf("fixture count %d", len(ab13aq13Cases))
	}
	for _, ab13aq13Case := range ab13aq13Cases {
		ab13aq13T.Run(ab13aq13Case.Ab13Aq13ID, func(ab13aq13T *ab13aq13testing.T) {
			ab13aq13Output := ab13aq13Evaluate(ab13aq13T, ab13aq13Case.Ab13Aq13Input)
			if ab13aq13Output["state"] != ab13aq13Case.Ab13Aq13State {
				ab13aq13T.Fatal(ab13aq13Output)
			}
			if !ab13aq13reflect.DeepEqual(ab13aq13Output["fields"], ab13aq13Case.Ab13Aq13Input["fields"]) {
				ab13aq13T.Fatal("raw fields changed")
			}
			for ab13aq13Name, ab13aq13Expected := range ab13aq13Case.Ab13Aq13Expected {
				ab13aq13Result := ab13aq13Results(ab13aq13Output)[ab13aq13Name].(map[string]any)
				for ab13aq13Key, ab13aq13Value := range ab13aq13Expected {
					if !ab13aq13reflect.DeepEqual(ab13aq13Result[ab13aq13Key], ab13aq13Value) {
						ab13aq13T.Fatalf("%s.%s got %v want %v", ab13aq13Name, ab13aq13Key, ab13aq13Result[ab13aq13Key], ab13aq13Value)
					}
				}
			}
			for _, ab13aq13Value := range ab13aq13Results(ab13aq13Output) {
				ab13aq13Result := ab13aq13Value.(map[string]any)
				if ab13aq13Result["grade"] != "Hypothesis" {
					ab13aq13T.Fatal("candidate grade lost")
				}
				if ab13aq13Result["state"] == "complete" && ab13aq13Result["calculation"] != "Proven" {
					ab13aq13T.Fatal("calculation grade lost")
				}
				for ab13aq13Key := range ab13aq13Result {
					switch ab13aq13Key {
					case "raw", "present", "state", "grade", "calculation", "flags", "reserved", "number", "unit":
					default:
						ab13aq13T.Fatalf("invented result property %s", ab13aq13Key)
					}
				}
			}
		})
	}
}

func TestAb13Aq13Context(ab13aq13T *ab13aq13testing.T) {
	for ab13aq13Key, ab13aq13Values := range map[string][]any{"primary": {0, 180, 182, 255, nil}, "secondary": {0, 32, 34, 255, nil}, "operation": {0, 1, 255, nil}, "direction": {"negative", "local", "", nil}, "view": {"other", "", nil}} {
		for ab13aq13Index, ab13aq13Value := range ab13aq13Values {
			ab13aq13T.Run(ab13aq13Key+ab13aq13strconv.Itoa(ab13aq13Index), func(ab13aq13T *ab13aq13testing.T) {
				ab13aq13Input := ab13aq13Input(ab13aq13T, "positive")
				ab13aq13Input["context"].(map[string]any)[ab13aq13Key] = ab13aq13Value
				ab13aq13Output := ab13aq13Evaluate(ab13aq13T, ab13aq13Input)
				ab13aq13ContextBefore, _ := ab13aq13json.Marshal(ab13aq13Input["context"])
				ab13aq13ContextAfter, _ := ab13aq13json.Marshal(ab13aq13Output["context"])
				if !ab13aq13bytes.Equal(ab13aq13ContextBefore, ab13aq13ContextAfter) {
					ab13aq13T.Fatal("unsupported context lost")
				}
				if ab13aq13Output["state"] != "Unknown" {
					ab13aq13T.Fatal("unsupported context evaluated")
				}
				for _, ab13aq13Value := range ab13aq13Results(ab13aq13Output) {
					ab13aq13Result := ab13aq13Value.(map[string]any)
					if len(ab13aq13Result) != 3 || ab13aq13Result["state"] != "Unknown" {
						ab13aq13T.Fatal("unsupported result", ab13aq13Result)
					}
				}
			})
		}
	}
	for _, ab13aq13Key := range []string{"primary", "secondary", "direction", "view"} {
		ab13aq13Input := ab13aq13Input(ab13aq13T, "positive")
		delete(ab13aq13Input["context"].(map[string]any), ab13aq13Key)
		if ab13aq13Evaluate(ab13aq13T, ab13aq13Input)["state"] != "Unknown" {
			ab13aq13T.Fatal("missing context evaluated", ab13aq13Key)
		}
	}
	for _, ab13aq13Key := range []string{"family", "variant"} {
		for _, ab13aq13Value := range []any{nil, false, 0, "", "synthetic-family", map[string]any{"opaque": true}} {
			ab13aq13Input := ab13aq13Input(ab13aq13T, "positive")
			ab13aq13Input["context"].(map[string]any)[ab13aq13Key] = ab13aq13Value
			if ab13aq13Evaluate(ab13aq13T, ab13aq13Input)["state"] != "complete" {
				ab13aq13T.Fatal("opaque identity selected view")
			}
		}
	}
}

func TestAb13Aq13WidthsAndSiblings(ab13aq13T *ab13aq13testing.T) {
	for _, ab13aq13Direction := range []string{"request", "positive"} {
		for ab13aq13Name, ab13aq13Original := range ab13aq13Input(ab13aq13T, ab13aq13Direction)["fields"].(map[string]any) {
			for _, ab13aq13Length := range []int{-1, 0, 1, 2, 3} {
				ab13aq13T.Run(ab13aq13Direction+ab13aq13Name+ab13aq13strconv.Itoa(ab13aq13Length), func(ab13aq13T *ab13aq13testing.T) {
					ab13aq13Input := ab13aq13Input(ab13aq13T, ab13aq13Direction)
					ab13aq13Fields := ab13aq13Input["fields"].(map[string]any)
					delete(ab13aq13Fields, ab13aq13Name)
					if ab13aq13Length >= 0 {
						ab13aq13Array := []any{}
						for ab13aq13Index := 0; ab13aq13Index < ab13aq13Length; ab13aq13Index++ {
							ab13aq13Array = append(ab13aq13Array, 0)
						}
						ab13aq13Fields[ab13aq13Name] = ab13aq13Array
					}
					ab13aq13Fields["opaque"] = []any{}
					ab13aq13Baseline := ab13aq13Results(ab13aq13Evaluate(ab13aq13T, ab13aq13InputWithoutDirectionChange(ab13aq13T, ab13aq13Direction)))
					ab13aq13Output := ab13aq13Evaluate(ab13aq13T, ab13aq13Input)
					ab13aq13Complete := ab13aq13Length == len(ab13aq13Original.([]any))
					ab13aq13Expected := "partial"
					if ab13aq13Complete {
						ab13aq13Expected = "complete"
					}
					if ab13aq13Output["state"] != ab13aq13Expected {
						ab13aq13T.Fatal(ab13aq13Output)
					}
					for ab13aq13Other, ab13aq13Value := range ab13aq13Results(ab13aq13Output) {
						ab13aq13Result := ab13aq13Value.(map[string]any)
						if ab13aq13Other == "opaque" {
							if ab13aq13Result["state"] != "Unknown" || !ab13aq13reflect.DeepEqual(ab13aq13Result["raw"], []any{}) {
								ab13aq13T.Fatal("empty opaque lost")
							}
							continue
						}
						if ab13aq13Other == ab13aq13Name && !ab13aq13Complete {
							if ab13aq13Result["state"] != "incomplete" || ab13aq13Result["present"] != (ab13aq13Length >= 0) {
								ab13aq13T.Fatal("missing/width confused")
							}
							if _, ab13aq13Exists := ab13aq13Result["number"]; ab13aq13Exists {
								ab13aq13T.Fatal("invented number")
							}
						} else if !ab13aq13reflect.DeepEqual(ab13aq13Result, ab13aq13Baseline[ab13aq13Other]) {
							ab13aq13T.Fatal("sibling changed", ab13aq13Other)
						}
					}
				})
			}
		}
	}
}

func TestAb13Aq13Reject(ab13aq13T *ab13aq13testing.T) {
	for ab13aq13Index, ab13aq13Raw := range []string{``, `{`, `[]`, `null`, `{} {}`, `{"context":{},"context":{},"fields":{}}`, `{"context":{"primary":181,"primary":181},"fields":{}}`, `{"context":{},"fields":{"x":[0],"x":[1]}}`, `{"context":{},"fields":{},"extra":[{"x":0,"x":1}]}`, `{"context":{},"fields":{},"extra":1e999}`, `{"context":{},"fields":{},"arguments":{"x":1,"x":2}}`} {
		ab13aq13T.Run("syntax"+ab13aq13strconv.Itoa(ab13aq13Index), func(ab13aq13T *ab13aq13testing.T) {
			ab13aq13Output, ab13aq13Error := Ab13Aq13Evaluate([]byte(ab13aq13Raw))
			if ab13aq13Error == nil || ab13aq13Output != nil {
				ab13aq13T.Fatal("invalid JSON accepted")
			}
		})
	}
	for _, ab13aq13Key := range []string{"primary", "secondary", "operation", "direction", "view", "fields", "context", "arguments", "field"} {
		for ab13aq13Index, ab13aq13Bad := range []any{nil, true, "1", -1, 256, 1.5, []any{}, map[string]any{}} {
			if ab13aq13Bad == nil && (ab13aq13Key == "primary" || ab13aq13Key == "secondary" || ab13aq13Key == "operation" || ab13aq13Key == "direction" || ab13aq13Key == "view") {
				continue
			}
			if (ab13aq13Key == "direction" || ab13aq13Key == "view") && ab13aq13Index == 2 {
				continue
			}
			if (ab13aq13Key == "fields" || ab13aq13Key == "context" || ab13aq13Key == "arguments") && ab13aq13Index == 7 {
				continue
			}
			ab13aq13T.Run(ab13aq13Key+ab13aq13strconv.Itoa(ab13aq13Index), func(ab13aq13T *ab13aq13testing.T) {
				ab13aq13Input := ab13aq13Input(ab13aq13T, "positive")
				switch ab13aq13Key {
				case "fields", "context", "arguments":
					ab13aq13Input[ab13aq13Key] = ab13aq13Bad
				case "field":
					ab13aq13Input["fields"].(map[string]any)["status0"] = []any{ab13aq13Bad}
				default:
					ab13aq13Input["context"].(map[string]any)[ab13aq13Key] = ab13aq13Bad
				}
				ab13aq13Raw, _ := ab13aq13json.Marshal(ab13aq13Input)
				ab13aq13Output, ab13aq13Error := Ab13Aq13Evaluate(ab13aq13Raw)
				if ab13aq13Error == nil || ab13aq13Output != nil {
					ab13aq13T.Fatal("wrong kind accepted", string(ab13aq13Raw))
				}
			})
		}
	}
	for _, ab13aq13Bad := range []any{nil, true, "0", 0, 1.5, []any{}} {
		ab13aq13Input := ab13aq13Input(ab13aq13T, "positive")
		ab13aq13Input["fields"].(map[string]any)["status0"] = ab13aq13Bad
		if ab13aq13Array, ab13aq13OK := ab13aq13Bad.([]any); ab13aq13OK && len(ab13aq13Array) == 0 {
			continue
		}
		ab13aq13Raw, _ := ab13aq13json.Marshal(ab13aq13Input)
		if ab13aq13Output, ab13aq13Error := Ab13Aq13Evaluate(ab13aq13Raw); ab13aq13Error == nil || ab13aq13Output != nil {
			ab13aq13T.Fatal("nonarray accepted")
		}
	}
	for _, ab13aq13Value := range []any{nil, false, 0, "", 1.5} {
		ab13aq13Input := ab13aq13Input(ab13aq13T, "positive")
		ab13aq13Input["arguments"] = map[string]any{"undeclared": ab13aq13Value}
		ab13aq13Raw, _ := ab13aq13json.Marshal(ab13aq13Input)
		if ab13aq13Output, ab13aq13Error := Ab13Aq13Evaluate(ab13aq13Raw); ab13aq13Error == nil || ab13aq13Output != nil {
			ab13aq13T.Fatal("undeclared argument accepted")
		}
	}
}

func TestAb13Aq13RetentionAndMemory(ab13aq13T *ab13aq13testing.T) {
	ab13aq13Input := ab13aq13Input(ab13aq13T, "positive")
	ab13aq13Input["extra"] = map[string]any{"absent-check": nil, "null": nil, "zero": 0, "false": false, "empty": []any{}, "nested": map[string]any{"value": "synthetic"}}
	ab13aq13Input["arguments"] = map[string]any{}
	ab13aq13Raw, _ := ab13aq13json.Marshal(ab13aq13Input)
	ab13aq13Before := append([]byte(nil), ab13aq13Raw...)
	ab13aq13First, ab13aq13Error := Ab13Aq13Evaluate(ab13aq13Raw)
	if ab13aq13Error != nil {
		ab13aq13T.Fatal(ab13aq13Error)
	}
	if !ab13aq13bytes.Equal(ab13aq13Raw, ab13aq13Before) {
		ab13aq13T.Fatal("input mutated")
	}
	ab13aq13Output := ab13aq13Evaluate(ab13aq13T, ab13aq13Input)
	ab13aq13Normalized := map[string]any{}
	_ = ab13aq13json.Unmarshal(ab13aq13Raw, &ab13aq13Normalized)
	for _, ab13aq13Key := range []string{"context", "fields", "arguments", "extra"} {
		if !ab13aq13reflect.DeepEqual(ab13aq13Output[ab13aq13Key], ab13aq13Normalized[ab13aq13Key]) {
			ab13aq13T.Fatal("retention", ab13aq13Key)
		}
	}
	ab13aq13Second, _ := Ab13Aq13Evaluate(ab13aq13Before)
	ab13aq13First[0] = '!'
	ab13aq13Raw[0] = '!'
	ab13aq13Third, ab13aq13Error := Ab13Aq13Evaluate(ab13aq13Before)
	if ab13aq13Error != nil || !ab13aq13bytes.Equal(ab13aq13Second, ab13aq13Third) {
		ab13aq13T.Fatal("shared storage")
	}
	ab13aq13NoArgs := ab13aq13Evaluate(ab13aq13T, ab13aq13InputWithoutArguments(ab13aq13T))
	if _, ab13aq13Present := ab13aq13NoArgs["arguments"]; ab13aq13Present {
		ab13aq13T.Fatal("invented empty arguments")
	}
	if ab13aq13strings.Contains(string(ab13aq13Third), "offset") || ab13aq13strings.Contains(string(ab13aq13Third), "position") {
		ab13aq13T.Fatal("invented coordinate")
	}
}

func ab13aq13InputWithoutArguments(ab13aq13T *ab13aq13testing.T) map[string]any {
	return ab13aq13Input(ab13aq13T, "positive")
}

func TestAb13Aq13SignedDomain(ab13aq13T *ab13aq13testing.T) {
	for ab13aq13Carrier := 0; ab13aq13Carrier < 65536; ab13aq13Carrier++ {
		ab13aq13Input := ab13aq13Input(ab13aq13T, "positive")
		ab13aq13Input["fields"].(map[string]any)["temperature0"] = []any{ab13aq13Carrier % 256, ab13aq13Carrier / 256}
		ab13aq13Input["fields"].(map[string]any)["temperature1"] = []any{255 - ab13aq13Carrier%256, 255 - ab13aq13Carrier/256}
		ab13aq13Expected := float64(ab13aq13Carrier) / 16
		if ab13aq13Carrier >= 32768 {
			ab13aq13Expected -= 4096
		}
		ab13aq13AllResults := ab13aq13Results(ab13aq13Evaluate(ab13aq13T, ab13aq13Input))
		ab13aq13Result := ab13aq13AllResults["temperature0"].(map[string]any)
		ab13aq13ComplementExpected := float64(65535-ab13aq13Carrier) / 16
		if ab13aq13Carrier < 32768 {
			ab13aq13ComplementExpected -= 4096
		}
		if ab13aq13AllResults["temperature1"].(map[string]any)["number"] != ab13aq13ComplementExpected {
			ab13aq13T.Fatal("temperature1 carrier", ab13aq13Carrier)
		}
		if ab13aq13Result["number"] != ab13aq13Expected {
			ab13aq13T.Fatalf("carrier %d: %v want %v", ab13aq13Carrier, ab13aq13Result["number"], ab13aq13Expected)
		}
	}
}

func ab13aq13InputWithoutDirectionChange(ab13aq13T *ab13aq13testing.T, ab13aq13Direction string) map[string]any {
	return ab13aq13Input(ab13aq13T, ab13aq13Direction)
}

func TestAb13Aq13OpaqueExtent(ab13aq13T *ab13aq13testing.T) {
	for ab13aq13Direction, ab13aq13Length := range map[string]int{"request": 7, "positive": 12} {
		ab13aq13Input := ab13aq13Input(ab13aq13T, ab13aq13Direction)
		ab13aq13Payload := []any{}
		for ab13aq13Index := 0; ab13aq13Index < ab13aq13Length; ab13aq13Index++ {
			ab13aq13Payload = append(ab13aq13Payload, 255)
		}
		ab13aq13Input["fields"] = map[string]any{"opaque": ab13aq13Payload}
		ab13aq13Output := ab13aq13Evaluate(ab13aq13T, ab13aq13Input)
		if ab13aq13Output["state"] != "partial" {
			ab13aq13T.Fatal("abstract extent interpreted")
		}
		for ab13aq13Name, ab13aq13Value := range ab13aq13Results(ab13aq13Output) {
			ab13aq13Result := ab13aq13Value.(map[string]any)
			if ab13aq13Name == "opaque" {
				if ab13aq13Result["state"] != "Unknown" {
					ab13aq13T.Fatal("opaque field interpreted")
				}
				continue
			}
			if ab13aq13Result["state"] != "incomplete" || ab13aq13Result["present"] != false {
				ab13aq13T.Fatal("unlocated field invented")
			}
			if _, ab13aq13Exists := ab13aq13Result["number"]; ab13aq13Exists {
				ab13aq13T.Fatal("unlocated scalar invented")
			}
		}
	}
}

func TestAb13Aq13UnsignedDomain(ab13aq13T *ab13aq13testing.T) {
	for _, ab13aq13Direction := range []string{"request", "positive"} {
		for ab13aq13Carrier := 0; ab13aq13Carrier < 256; ab13aq13Carrier++ {
			ab13aq13Input := ab13aq13Input(ab13aq13T, ab13aq13Direction)
			ab13aq13Expected := map[string]float64{"brine_percent": float64(ab13aq13Carrier), "fan_percent": float64(ab13aq13Carrier), "fan_speed": float64(ab13aq13Carrier * 10)}
			if ab13aq13Direction == "request" {
				ab13aq13Expected = map[string]float64{"noise_percent": float64(ab13aq13Carrier), "frost_temperature": float64(ab13aq13Carrier)}
			}
			for ab13aq13Name := range ab13aq13Expected {
				ab13aq13Input["fields"].(map[string]any)[ab13aq13Name] = []any{ab13aq13Carrier}
			}
			ab13aq13Results := ab13aq13Results(ab13aq13Evaluate(ab13aq13T, ab13aq13Input))
			for ab13aq13Name, ab13aq13Number := range ab13aq13Expected {
				if ab13aq13Results[ab13aq13Name].(map[string]any)["number"] != ab13aq13Number {
					ab13aq13T.Fatal("unsigned carrier", ab13aq13Name, ab13aq13Carrier)
				}
			}
		}
	}
}
