package vaillant

import (
	ab19ao19bytes "bytes"
	ab19ao19json "encoding/json"
	ab19ao19fmt "fmt"
	ab19ao19io "io"
	ab19ao19strconv "strconv"
	ab19ao19utf8 "unicode/utf8"
)

func ab19ao19Read(ab19ao19Decoder *ab19ao19json.Decoder) (any, error) {
	ab19ao19Token, ab19ao19Err := ab19ao19Decoder.Token()
	if ab19ao19Err != nil {
		return nil, ab19ao19Err
	}
	if ab19ao19Delimiter, ab19ao19OK := ab19ao19Token.(ab19ao19json.Delim); ab19ao19OK {
		switch ab19ao19Delimiter {
		case '{':
			ab19ao19Object := map[string]any{}
			for ab19ao19Decoder.More() {
				ab19ao19KeyToken, ab19ao19Err := ab19ao19Decoder.Token()
				if ab19ao19Err != nil {
					return nil, ab19ao19Err
				}
				ab19ao19Key, ab19ao19OK := ab19ao19KeyToken.(string)
				if !ab19ao19OK {
					return nil, ab19ao19fmt.Errorf("invalid object key")
				}
				if _, ab19ao19Exists := ab19ao19Object[ab19ao19Key]; ab19ao19Exists {
					return nil, ab19ao19fmt.Errorf("duplicate key %q", ab19ao19Key)
				}
				ab19ao19Value, ab19ao19Err := ab19ao19Read(ab19ao19Decoder)
				if ab19ao19Err != nil {
					return nil, ab19ao19Err
				}
				ab19ao19Object[ab19ao19Key] = ab19ao19Value
			}
			if _, ab19ao19Err := ab19ao19Decoder.Token(); ab19ao19Err != nil {
				return nil, ab19ao19Err
			}
			return ab19ao19Object, nil
		case '[':
			ab19ao19Array := []any{}
			for ab19ao19Decoder.More() {
				ab19ao19Value, ab19ao19Err := ab19ao19Read(ab19ao19Decoder)
				if ab19ao19Err != nil {
					return nil, ab19ao19Err
				}
				ab19ao19Array = append(ab19ao19Array, ab19ao19Value)
			}
			if _, ab19ao19Err := ab19ao19Decoder.Token(); ab19ao19Err != nil {
				return nil, ab19ao19Err
			}
			return ab19ao19Array, nil
		default:
			return nil, ab19ao19fmt.Errorf("unexpected delimiter")
		}
	}
	return ab19ao19Token, nil
}

func ab19ao19Integer(ab19ao19Value any) (int, bool) {
	ab19ao19Number, ab19ao19OK := ab19ao19Value.(ab19ao19json.Number)
	if !ab19ao19OK {
		return 0, false
	}
	ab19ao19Int, ab19ao19Err := ab19ao19strconv.Atoi(string(ab19ao19Number))
	return ab19ao19Int, ab19ao19Err == nil
}

func Ab19Ao19Interpret(ab19ao19Input ab19ao19json.RawMessage) (ab19ao19json.RawMessage, error) {
	if !ab19ao19utf8.Valid(ab19ao19Input) {
		return nil, ab19ao19fmt.Errorf("invalid UTF-8 JSON")
	}
	ab19ao19Decoder := ab19ao19json.NewDecoder(ab19ao19bytes.NewReader(ab19ao19Input))
	ab19ao19Decoder.UseNumber()
	ab19ao19Value, ab19ao19Err := ab19ao19Read(ab19ao19Decoder)
	if ab19ao19Err != nil {
		return nil, ab19ao19Err
	}
	if _, ab19ao19Err = ab19ao19Decoder.Token(); ab19ao19Err != ab19ao19io.EOF {
		return nil, ab19ao19fmt.Errorf("trailing JSON")
	}
	ab19ao19InputObject, ab19ao19OK := ab19ao19Value.(map[string]any)
	if !ab19ao19OK {
		return nil, ab19ao19fmt.Errorf("input must be an object")
	}
	for ab19ao19Key := range ab19ao19InputObject {
		if ab19ao19Key != "context" && ab19ao19Key != "fields" && ab19ao19Key != "arguments" {
			return nil, ab19ao19fmt.Errorf("unknown input key %q", ab19ao19Key)
		}
	}
	ab19ao19Context, ab19ao19OK := ab19ao19InputObject["context"].(map[string]any)
	if !ab19ao19OK {
		return nil, ab19ao19fmt.Errorf("context must be an object")
	}
	ab19ao19Raw, ab19ao19OK := ab19ao19InputObject["fields"].(map[string]any)
	if !ab19ao19OK {
		return nil, ab19ao19fmt.Errorf("fields must be an object")
	}
	for ab19ao19Key, ab19ao19Value := range ab19ao19Context {
		switch ab19ao19Key {
		case "primary", "secondary", "operation":
			if _, ab19ao19OK := ab19ao19Integer(ab19ao19Value); !ab19ao19OK {
				return nil, ab19ao19fmt.Errorf("%s must be an integer", ab19ao19Key)
			}
		case "direction", "view":
			if _, ab19ao19OK := ab19ao19Value.(string); !ab19ao19OK {
				return nil, ab19ao19fmt.Errorf("%s must be a string", ab19ao19Key)
			}
		}
	}
	for ab19ao19Key, ab19ao19Value := range ab19ao19Raw {
		ab19ao19Array, ab19ao19OK := ab19ao19Value.([]any)
		if !ab19ao19OK {
			return nil, ab19ao19fmt.Errorf("%s must be a byte array", ab19ao19Key)
		}
		for _, ab19ao19Element := range ab19ao19Array {
			ab19ao19Byte, ab19ao19OK := ab19ao19Integer(ab19ao19Element)
			if !ab19ao19OK || ab19ao19Byte < 0 || ab19ao19Byte > 255 {
				return nil, ab19ao19fmt.Errorf("invalid byte in %s", ab19ao19Key)
			}
		}
	}
	ab19ao19Fields := map[string]any{}
	ab19ao19Output := map[string]any{"context": ab19ao19Context, "raw": ab19ao19Raw, "fields": ab19ao19Fields, "state": "Unknown", "physical": "Unknown", "qualified": false}
	if ab19ao19Arguments, ab19ao19Present := ab19ao19InputObject["arguments"]; ab19ao19Present {
		ab19ao19Object, ab19ao19OK := ab19ao19Arguments.(map[string]any)
		if !ab19ao19OK || len(ab19ao19Object) != 0 {
			return nil, ab19ao19fmt.Errorf("no arguments are defined for these views")
		}
		ab19ao19Output["arguments"] = ab19ao19Object
	}
	ab19ao19Primary, _ := ab19ao19Integer(ab19ao19Context["primary"])
	ab19ao19Secondary, _ := ab19ao19Integer(ab19ao19Context["secondary"])
	ab19ao19Operation, _ := ab19ao19Integer(ab19ao19Context["operation"])
	ab19ao19View, _ := ab19ao19Context["view"].(string)
	ab19ao19Selected := ab19ao19Primary == 181 && ab19ao19Secondary == 20 && ab19ao19Context["direction"] == "positive"
	if !ab19ao19Selected {
		return ab19ao19json.Marshal(ab19ao19Output)
	}
	ab19ao19Add := func(ab19ao19Name string, ab19ao19Data []any, ab19ao19Start, ab19ao19Width int) map[string]any {
		ab19ao19End := ab19ao19Start + ab19ao19Width
		if ab19ao19Start > len(ab19ao19Data) {
			ab19ao19Start = len(ab19ao19Data)
		}
		if ab19ao19End > len(ab19ao19Data) {
			ab19ao19End = len(ab19ao19Data)
		}
		ab19ao19Field := map[string]any{"raw": ab19ao19Data[ab19ao19Start:ab19ao19End], "state": "incomplete", "grade": "Hypothesis"}
		ab19ao19Fields[ab19ao19Name] = ab19ao19Field
		return ab19ao19Field
	}
	if (ab19ao19Operation == 6 && (ab19ao19View == "generic-result" || ab19ao19View == "specific-a-result" || ab19ao19View == "specific-b-result")) || (ab19ao19Operation == 7 && ab19ao19View == "progress") {
		ab19ao19Data, ab19ao19Present := ab19ao19Raw["feedback"].([]any)
		if !ab19ao19Present {
			ab19ao19Data = []any{}
		}
		ab19ao19Field := ab19ao19Add("feedback", ab19ao19Data, 0, 1)
		ab19ao19Output["state"] = "partial"
		if len(ab19ao19Data) > 0 {
			ab19ao19Output["state"] = "Hypothesis"
			ab19ao19Field["state"] = "unmapped"
			ab19ao19Byte, _ := ab19ao19Integer(ab19ao19Data[0])
			ab19ao19Label := ""
			if ab19ao19Operation == 6 {
				if ab19ao19Byte == 1 {
					if ab19ao19View == "generic-result" {
						ab19ao19Label = "in-progress"
					} else {
						ab19ao19Label = "not-successful"
					}
				}
				if ab19ao19Byte == 255 && ab19ao19View == "generic-result" {
					ab19ao19Label = "unsupported"
				}
				if ab19ao19Byte >= 3 && ab19ao19View != "generic-result" {
					ab19ao19Label = "invalid"
				}
			} else {
				ab19ao19Label = map[int]string{0: "idle", 1: "progress", 2: "passed", 3: "warning", 4: "combustion-error", 5: "gas-type-error", 6: "aborted", 255: "unsupported"}[ab19ao19Byte]
			}
			if ab19ao19Label != "" {
				ab19ao19Field["label"] = ab19ao19Label
				ab19ao19Field["state"] = "Hypothesis"
			}
		}
	} else if ab19ao19Operation == 8 && ab19ao19View == "pressure" {
		ab19ao19Data, _ := ab19ao19Raw["logical_object"].([]any)
		if ab19ao19Data == nil {
			ab19ao19Data = []any{}
		}
		ab19ao19Output["state"] = "Hypothesis"
		if len(ab19ao19Data) < 13 {
			ab19ao19Output["state"] = "partial"
		}
		for ab19ao19Index, ab19ao19Name := range []string{"feedback", "pressure", "minimum", "maximum", "step", "flow"} {
			ab19ao19Start := []int{3, 4, 6, 8, 10, 11}[ab19ao19Index]
			ab19ao19Width := []int{1, 2, 2, 2, 1, 2}[ab19ao19Index]
			ab19ao19Field := ab19ao19Add(ab19ao19Name, ab19ao19Data, ab19ao19Start, ab19ao19Width)
			if len(ab19ao19Data) < ab19ao19Start+ab19ao19Width {
				continue
			}
			ab19ao19Number, _ := ab19ao19Integer(ab19ao19Data[ab19ao19Start])
			if ab19ao19Width == 2 {
				ab19ao19High, _ := ab19ao19Integer(ab19ao19Data[ab19ao19Start+1])
				ab19ao19Number += ab19ao19High << 8
			}
			ab19ao19Field["state"] = "Hypothesis"
			switch ab19ao19Name {
			case "feedback":
				ab19ao19Field["state"] = "unmapped"
			case "step":
				ab19ao19Field["state"] = "unit-conflict"
				ab19ao19Field["flags"] = []string{"constant-1-mbar-versus-default-0"}
			default:
				ab19ao19Field["calculation"] = "Proven"
				if ab19ao19Name == "flow" {
					ab19ao19Field["number"] = float64(ab19ao19Number) / 100
				} else {
					ab19ao19Field["number"] = ab19ao19Number
				}
			}
		}
	}
	return ab19ao19json.Marshal(ab19ao19Output)
}
