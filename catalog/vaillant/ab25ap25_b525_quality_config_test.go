package vaillant

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

func ab25ap25Input(ab25ap25Operation int, ab25ap25Direction, ab25ap25View string, ab25ap25Fields map[string]any) map[string]any {
	return map[string]any{"context": map[string]any{"primary": 181, "secondary": 37, "operation": ab25ap25Operation, "direction": ab25ap25Direction, "view": ab25ap25View, "family": "synthetic-family", "variant": "synthetic-variant", "opaque": map[string]any{"flag": false, "zero": 0, "null": nil}}, "fields": ab25ap25Fields}
}

func ab25ap25Evaluate(ab25ap25T *testing.T, ab25ap25InputObject map[string]any) map[string]any {
	ab25ap25T.Helper()
	ab25ap25Bytes, ab25ap25Error := json.Marshal(ab25ap25InputObject)
	if ab25ap25Error != nil {
		ab25ap25T.Fatal(ab25ap25Error)
	}
	ab25ap25Output, ab25ap25Error := Ab25Ap25Evaluate(ab25ap25Bytes)
	if ab25ap25Error != nil {
		ab25ap25T.Fatal(ab25ap25Error)
	}
	var ab25ap25Result map[string]any
	if ab25ap25Error := json.Unmarshal(ab25ap25Output, &ab25ap25Result); ab25ap25Error != nil {
		ab25ap25T.Fatal(ab25ap25Error)
	}
	var ab25ap25RoundTrip map[string]any
	_ = json.Unmarshal(ab25ap25Bytes, &ab25ap25RoundTrip)
	if !reflect.DeepEqual(ab25ap25RoundTrip["context"], ab25ap25Result["context"]) || !reflect.DeepEqual(ab25ap25RoundTrip["fields"], ab25ap25Result["raw"]) {
		ab25ap25T.Fatalf("raw/context changed: %s", ab25ap25Output)
	}
	if ab25ap25Arguments, ab25ap25Exists := ab25ap25RoundTrip["arguments"]; ab25ap25Exists && !reflect.DeepEqual(ab25ap25Arguments, ab25ap25Result["arguments"]) {
		ab25ap25T.Fatal("arguments changed")
	}
	if ab25ap25Result["physical"] != "Unknown" || ab25ap25Result["qualified"] != false {
		ab25ap25T.Fatal("qualification gained")
	}
	return ab25ap25Result
}

func ab25ap25Field(ab25ap25Output map[string]any, ab25ap25Name string) map[string]any {
	return ab25ap25Output["fields"].(map[string]any)[ab25ap25Name].(map[string]any)
}

func TestAb25Ap25RSSI(ab25ap25T *testing.T) {
	ab25ap25Bytes, ab25ap25Error := os.ReadFile("testdata/ab25ap25_b525_quality_config.json")
	if ab25ap25Error != nil {
		ab25ap25T.Fatal(ab25ap25Error)
	}
	var ab25ap25Cases []map[string]any
	if ab25ap25Error := json.Unmarshal(ab25ap25Bytes, &ab25ap25Cases); ab25ap25Error != nil {
		ab25ap25T.Fatal(ab25ap25Error)
	}
	if len(ab25ap25Cases) != 4 {
		ab25ap25T.Fatal("missing cases")
	}
	for _, ab25ap25Case := range ab25ap25Cases {
		ab25ap25T.Run(ab25ap25Case["id"].(string), func(ab25ap25T *testing.T) {
			ab25ap25Result := ab25ap25Evaluate(ab25ap25T, ab25ap25Input(3, "positive", "quality", map[string]any{"rssi": []any{ab25ap25Case["byte"]}, "quality": []int{0}}))
			ab25ap25RSSI := ab25ap25Field(ab25ap25Result, "rssi")
			if ab25ap25Result["state"] != "complete" || ab25ap25RSSI["number"] != ab25ap25Case["number"] || ab25ap25RSSI["unit"] != "Unknown" || ab25ap25RSSI["state"] != "Hypothesis" || ab25ap25RSSI["calculation"] != "Proven" {
				ab25ap25T.Fatalf("RSSI: %v", ab25ap25RSSI)
			}
		})
	}
}

func TestAb25Ap25Quality(ab25ap25T *testing.T) {
	for ab25ap25Bit := -1; ab25ap25Bit < 8; ab25ap25Bit++ {
		ab25ap25ByteValue := 0
		if ab25ap25Bit >= 0 {
			ab25ap25ByteValue = 1 << ab25ap25Bit
		}
		ab25ap25T.Run(fmt.Sprintf("ab25ap25_bit_%d", ab25ap25Bit), func(ab25ap25T *testing.T) {
			ab25ap25Result := ab25ap25Evaluate(ab25ap25T, ab25ap25Input(3, "positive", "quality", map[string]any{"rssi": []int{255}, "quality": []int{ab25ap25ByteValue}}))
			ab25ap25Quality := ab25ap25Field(ab25ap25Result, "quality")
			ab25ap25Flags, ab25ap25OK := ab25ap25Quality["flags"].(map[string]any)
			if !ab25ap25OK {
				ab25ap25T.Fatal("missing independent quality flags")
			}
			if len(ab25ap25Flags) != 8 || ab25ap25Quality["checksum_verified"] != false || ab25ap25Quality["state"] != "Hypothesis" {
				ab25ap25T.Fatal(ab25ap25Quality)
			}
			for ab25ap25Index := 0; ab25ap25Index < 8; ab25ap25Index++ {
				ab25ap25Name := fmt.Sprintf("bit%d", ab25ap25Index)
				if ab25ap25Index == 7 {
					ab25ap25Name = "crc_okay"
				}
				if ab25ap25Flags[ab25ap25Name] != (ab25ap25Index == ab25ap25Bit) {
					ab25ap25T.Fatal(ab25ap25Flags)
				}
			}
			for _, ab25ap25Name := range []string{"number", "score", "label"} {
				if _, ab25ap25Exists := ab25ap25Quality[ab25ap25Name]; ab25ap25Exists {
					ab25ap25T.Fatal("invented quality quantity")
				}
			}
		})
	}
}

func TestAb25Ap25Configuration(ab25ap25T *testing.T) {
	for _, ab25ap25Operation := range []int{24, 25, 26, 27} {
		for _, ab25ap25Direction := range []string{"request", "positive"} {
			for _, ab25ap25Config := range []int{0, 129, 254, 1} {
				for ab25ap25DataCase := 0; ab25ap25DataCase < 3; ab25ap25DataCase++ {
					ab25ap25T.Run(fmt.Sprintf("ab25ap25_%d_%s_%d_%d", ab25ap25Operation, ab25ap25Direction, ab25ap25Config, ab25ap25DataCase), func(ab25ap25T *testing.T) {
						ab25ap25Fields := map[string]any{"configuration": []int{ab25ap25Config}}
						if ab25ap25DataCase == 0 {
							ab25ap25Fields["data"] = []int{0, 170, 255}
						}
						if ab25ap25DataCase == 1 {
							ab25ap25Fields["data"] = []int{}
						}
						ab25ap25Result := ab25ap25Evaluate(ab25ap25T, ab25ap25Input(ab25ap25Operation, ab25ap25Direction, "configuration", ab25ap25Fields))
						ab25ap25ConfigResult := ab25ap25Field(ab25ap25Result, "configuration")
						if ab25ap25ConfigResult["state"] != "Hypothesis" || ab25ap25ConfigResult["unmodeled"] != float64(ab25ap25Config&254) || ab25ap25ConfigResult["flags"].(map[string]any)["bit0"] != (ab25ap25Config&1 != 0) {
							ab25ap25T.Fatal(ab25ap25ConfigResult)
						}
						if !reflect.DeepEqual(ab25ap25ConfigResult["candidate_default"], map[string]any{"number": float64(0), "label": "read", "state": "Hypothesis"}) {
							ab25ap25T.Fatal("default metadata")
						}
						ab25ap25Data := ab25ap25Field(ab25ap25Result, "data")
						if ab25ap25DataCase == 2 {
							if ab25ap25Result["state"] != "partial" || ab25ap25Data["state"] != "Unknown" {
								ab25ap25T.Fatal(ab25ap25Result)
							}
							if _, ab25ap25Exists := ab25ap25Data["raw"]; ab25ap25Exists {
								ab25ap25T.Fatal("missing became empty")
							}
						} else if ab25ap25Result["state"] != "complete" || ab25ap25Data["state"] != "Hypothesis" {
							ab25ap25T.Fatal(ab25ap25Result)
						}
					})
				}
			}
		}
	}
}

func TestAb25Ap25Incomplete(ab25ap25T *testing.T) {
	for _, ab25ap25View := range []string{"quality", "configuration"} {
		ab25ap25Names := []string{"rssi", "quality"}
		ab25ap25Operation := 3
		if ab25ap25View == "configuration" {
			ab25ap25Names = []string{"configuration", "data"}
			ab25ap25Operation = 24
		}
		for _, ab25ap25Name := range ab25ap25Names {
			for ab25ap25Case := 0; ab25ap25Case < 3; ab25ap25Case++ {
				if ab25ap25Name == "data" && ab25ap25Case > 0 {
					continue
				}
				ab25ap25Fields := map[string]any{ab25ap25Names[0]: []int{0}, ab25ap25Names[1]: []int{0}}
				delete(ab25ap25Fields, ab25ap25Name)
				if ab25ap25Case == 1 {
					ab25ap25Fields[ab25ap25Name] = []int{}
				}
				if ab25ap25Case == 2 {
					ab25ap25Fields[ab25ap25Name] = []int{0, 255}
				}
				ab25ap25Result := ab25ap25Evaluate(ab25ap25T, ab25ap25Input(ab25ap25Operation, "positive", ab25ap25View, ab25ap25Fields))
				if ab25ap25Result["state"] != "partial" || ab25ap25Field(ab25ap25Result, ab25ap25Name)["state"] != "Unknown" {
					ab25ap25T.Fatal(ab25ap25Result)
				}
				for _, ab25ap25Sibling := range ab25ap25Names {
					if ab25ap25Sibling != ab25ap25Name && ab25ap25Field(ab25ap25Result, ab25ap25Sibling)["state"] != "Hypothesis" {
						ab25ap25T.Fatal("lost sibling")
					}
				}
				for _, ab25ap25Key := range []string{"number", "label", "flags"} {
					if _, ab25ap25Exists := ab25ap25Field(ab25ap25Result, ab25ap25Name)[ab25ap25Key]; ab25ap25Exists {
						ab25ap25T.Fatal("invented incomplete value")
					}
				}
			}
		}
	}
}

func TestAb25Ap25ContextArguments(ab25ap25T *testing.T) {
	for _, ab25ap25Key := range []string{"primary", "secondary", "operation", "direction", "view", "family", "variant", "opaque"} {
		for _, ab25ap25Value := range []any{nil, 0, false, "other", map[string]any{"nested": []any{0, false, nil}}} {
			if (ab25ap25Key == "primary" || ab25ap25Key == "secondary" || ab25ap25Key == "operation") && ab25ap25Value != nil && ab25ap25Value != 0 {
				continue
			}
			if (ab25ap25Key == "direction" || ab25ap25Key == "view") && ab25ap25Value != nil && ab25ap25Value != "other" {
				continue
			}
			ab25ap25InputObject := ab25ap25Input(3, "positive", "quality", map[string]any{"rssi": []int{0}, "quality": []int{0}})
			ab25ap25InputObject["context"].(map[string]any)[ab25ap25Key] = ab25ap25Value
			ab25ap25Result := ab25ap25Evaluate(ab25ap25T, ab25ap25InputObject)
			if (ab25ap25Key == "family" || ab25ap25Key == "variant" || ab25ap25Key == "opaque") && ab25ap25Result["state"] != "complete" {
				ab25ap25T.Fatal("opaque key affected admission")
			}
			if ab25ap25Key != "family" && ab25ap25Key != "variant" && ab25ap25Key != "opaque" && (ab25ap25Result["state"] != "Unknown" || len(ab25ap25Result["fields"].(map[string]any)) != 0) {
				ab25ap25T.Fatal("selector admitted")
			}
		}
		ab25ap25InputObject := ab25ap25Input(3, "positive", "quality", map[string]any{"rssi": []int{0}, "quality": []int{0}})
		delete(ab25ap25InputObject["context"].(map[string]any), ab25ap25Key)
		ab25ap25Result := ab25ap25Evaluate(ab25ap25T, ab25ap25InputObject)
		if ab25ap25Key != "family" && ab25ap25Key != "variant" && ab25ap25Key != "opaque" && ab25ap25Result["state"] != "Unknown" {
			ab25ap25T.Fatal("absent selector admitted")
		}
	}
	for _, ab25ap25Value := range []any{nil, false, 0, []any{}, map[string]any{}, "opaque"} {
		ab25ap25InputObject := ab25ap25Input(3, "positive", "quality", map[string]any{"rssi": []int{0}, "quality": []int{0}, "extra": []int{255}})
		ab25ap25InputObject["arguments"] = map[string]any{"extra": ab25ap25Value}
		ab25ap25Evaluate(ab25ap25T, ab25ap25InputObject)
	}
}

func TestAb25Ap25Invalid(ab25ap25T *testing.T) {
	ab25ap25Valid := `{"context":{"primary":181,"secondary":37,"operation":3,"direction":"positive","view":"quality"},"fields":{"rssi":[0],"quality":[0]}}`
	ab25ap25Cases := []string{``, `{`, `[]`, ab25ap25Valid + ` {}`, strings.Replace(ab25ap25Valid, `"primary":181`, `"primary":181,"primary":181`, 1), strings.Replace(ab25ap25Valid, `"rssi":[0]`, `"rssi":[0],"rssi":[0]`, 1), strings.Replace(ab25ap25Valid, `"rssi":[0]`, `"rssi":null`, 1), strings.Replace(ab25ap25Valid, `"context":{`, `"arguments":null,"context":{`, 1), strings.Replace(ab25ap25Valid, `"context":{`, `"arguments":{"x":{"a":0,"a":1}},"context":{`, 1)}
	for _, ab25ap25Value := range []string{`true`, `false`, `null`, `"0"`, `-1`, `256`, `0.5`, `1.00000000000000000001`, `[]`, `{}`} {
		ab25ap25Cases = append(ab25ap25Cases, strings.Replace(ab25ap25Valid, `"rssi":[0]`, `"rssi":[`+ab25ap25Value+`]`, 1))
	}
	for _, ab25ap25InputBytes := range ab25ap25Cases {
		ab25ap25Output, ab25ap25Error := Ab25Ap25Evaluate([]byte(ab25ap25InputBytes))
		if ab25ap25Error == nil || ab25ap25Output != nil {
			ab25ap25T.Fatalf("invalid admitted %s", ab25ap25InputBytes)
		}
	}
}

func TestAb25Ap25Ownership(ab25ap25T *testing.T) {
	ab25ap25InputBytes := json.RawMessage(`{"context":{"primary":181,"secondary":37,"operation":24,"direction":"request","view":"configuration","opaque":{"a":[0]}},"fields":{"configuration":[0],"data":[0,170,255]},"arguments":{"opaque":{"a":[false]}}}`)
	ab25ap25Saved := append(json.RawMessage(nil), ab25ap25InputBytes...)
	ab25ap25First, ab25ap25Error := Ab25Ap25Evaluate(ab25ap25InputBytes)
	if ab25ap25Error != nil {
		ab25ap25T.Fatal(ab25ap25Error)
	}
	ab25ap25Expected := append(json.RawMessage(nil), ab25ap25First...)
	for ab25ap25Index := range ab25ap25InputBytes {
		ab25ap25InputBytes[ab25ap25Index] = ' '
	}
	if string(ab25ap25First) != string(ab25ap25Expected) {
		ab25ap25T.Fatal("input aliases output")
	}
	var ab25ap25Decoded map[string]any
	_ = json.Unmarshal(ab25ap25First, &ab25ap25Decoded)
	ab25ap25Decoded["raw"].(map[string]any)["data"].([]any)[0] = 99
	ab25ap25Decoded["context"].(map[string]any)["opaque"].(map[string]any)["a"].([]any)[0] = 99
	ab25ap25Decoded["arguments"].(map[string]any)["opaque"].(map[string]any)["a"].([]any)[0] = true
	if ab25ap25Field(ab25ap25Decoded, "data")["raw"].([]any)[0] != float64(0) {
		ab25ap25T.Fatal("raw aliases semantic data")
	}
	for ab25ap25Index := range ab25ap25First {
		ab25ap25First[ab25ap25Index] = ' '
	}
	ab25ap25Second, ab25ap25Error := Ab25Ap25Evaluate(ab25ap25Saved)
	if ab25ap25Error != nil || string(ab25ap25Second) != string(ab25ap25Expected) {
		ab25ap25T.Fatal("cross-call alias")
	}
}

func TestAb25Ap25Admission(ab25ap25T *testing.T) {
	for _, ab25ap25View := range []string{"quality", "configuration"} {
		ab25ap25Operation := 3
		ab25ap25Names := []string{"rssi", "quality"}
		if ab25ap25View == "configuration" {
			ab25ap25Operation = 24
			ab25ap25Names = []string{"configuration", "data"}
		}
		for _, ab25ap25Name := range ab25ap25Names {
			for _, ab25ap25Bad := range []any{nil, false, 0, "", []any{nil}, []any{false}, []any{"0"}, []any{-1}, []any{256}, []any{0.5}} {
				ab25ap25InputObject := ab25ap25Input(ab25ap25Operation, "positive", ab25ap25View, map[string]any{ab25ap25Names[0]: []int{0}, ab25ap25Names[1]: []int{0}})
				ab25ap25InputObject["fields"].(map[string]any)[ab25ap25Name] = ab25ap25Bad
				ab25ap25Bytes, _ := json.Marshal(ab25ap25InputObject)
				ab25ap25Output, ab25ap25Error := Ab25Ap25Evaluate(ab25ap25Bytes)
				if ab25ap25Error == nil || ab25ap25Output != nil {
					ab25ap25T.Fatalf("invalid %s admitted: %s", ab25ap25Name, ab25ap25Bytes)
				}
			}
		}
		for _, ab25ap25Key := range []string{"primary", "secondary", "operation", "direction", "view"} {
			ab25ap25InputObject := ab25ap25Input(ab25ap25Operation, "positive", ab25ap25View, map[string]any{ab25ap25Names[0]: []int{0}, ab25ap25Names[1]: []int{0}})
			ab25ap25InputObject["context"].(map[string]any)[ab25ap25Key] = nil
			ab25ap25Result := ab25ap25Evaluate(ab25ap25T, ab25ap25InputObject)
			if ab25ap25Result["state"] != "Unknown" || len(ab25ap25Result["fields"].(map[string]any)) != 0 {
				ab25ap25T.Fatal("null selector admitted")
			}
		}
	}
	for _, ab25ap25Direction := range []string{"request", "negative", ""} {
		ab25ap25Result := ab25ap25Evaluate(ab25ap25T, ab25ap25Input(3, ab25ap25Direction, "quality", map[string]any{"rssi": []int{0}, "quality": []int{0}}))
		if ab25ap25Result["state"] != "Unknown" {
			ab25ap25T.Fatal("quality direction admitted")
		}
	}
	for _, ab25ap25Operation := range []int{3, 23, 28} {
		ab25ap25Result := ab25ap25Evaluate(ab25ap25T, ab25ap25Input(ab25ap25Operation, "positive", "configuration", map[string]any{"configuration": []int{0}, "data": []int{}}))
		if ab25ap25Result["state"] != "Unknown" {
			ab25ap25T.Fatal("configuration operation admitted")
		}
	}
	ab25ap25Valid := json.RawMessage(`{"context":{"primary":181,"secondary":37,"operation":3,"direction":"positive","view":"quality"},"fields":{"rssi":[0e0],"quality":[1.0]}}`)
	if _, ab25ap25Error := Ab25Ap25Evaluate(ab25ap25Valid); ab25ap25Error != nil {
		ab25ap25T.Fatal("integral JSON numbers rejected", ab25ap25Error)
	}
}
