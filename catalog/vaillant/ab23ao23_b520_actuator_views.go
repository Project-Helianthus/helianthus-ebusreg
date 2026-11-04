package vaillant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode/utf8"
)

func ab23ao23Value(ab23ao23Decoder *json.Decoder) (any, error) {
	ab23ao23Token, ab23ao23Err := ab23ao23Decoder.Token()
	if ab23ao23Err != nil {
		return nil, ab23ao23Err
	}
	if ab23ao23Delimiter, ab23ao23OK := ab23ao23Token.(json.Delim); ab23ao23OK {
		switch ab23ao23Delimiter {
		case '{':
			ab23ao23Object := map[string]any{}
			for ab23ao23Decoder.More() {
				ab23ao23KeyToken, ab23ao23Err := ab23ao23Decoder.Token()
				if ab23ao23Err != nil {
					return nil, ab23ao23Err
				}
				ab23ao23Key, ab23ao23OK := ab23ao23KeyToken.(string)
				if !ab23ao23OK {
					return nil, fmt.Errorf("object key must be a string")
				}
				if _, ab23ao23Exists := ab23ao23Object[ab23ao23Key]; ab23ao23Exists {
					return nil, fmt.Errorf("duplicate object key %q", ab23ao23Key)
				}
				ab23ao23Child, ab23ao23Err := ab23ao23Value(ab23ao23Decoder)
				if ab23ao23Err != nil {
					return nil, ab23ao23Err
				}
				ab23ao23Object[ab23ao23Key] = ab23ao23Child
			}
			ab23ao23End, ab23ao23Err := ab23ao23Decoder.Token()
			if ab23ao23Err != nil || ab23ao23End != json.Delim('}') {
				return nil, fmt.Errorf("unterminated object")
			}
			return ab23ao23Object, nil
		case '[':
			ab23ao23Array := []any{}
			for ab23ao23Decoder.More() {
				ab23ao23Child, ab23ao23Err := ab23ao23Value(ab23ao23Decoder)
				if ab23ao23Err != nil {
					return nil, ab23ao23Err
				}
				ab23ao23Array = append(ab23ao23Array, ab23ao23Child)
			}
			ab23ao23End, ab23ao23Err := ab23ao23Decoder.Token()
			if ab23ao23Err != nil || ab23ao23End != json.Delim(']') {
				return nil, fmt.Errorf("unterminated array")
			}
			return ab23ao23Array, nil
		default:
			return nil, fmt.Errorf("unexpected delimiter")
		}
	}
	return ab23ao23Token, nil
}

func ab23ao23Integer(ab23ao23Value any) (int64, bool) {
	ab23ao23Number, ab23ao23OK := ab23ao23Value.(json.Number)
	if !ab23ao23OK {
		return 0, false
	}
	ab23ao23Int, ab23ao23Err := strconv.ParseInt(string(ab23ao23Number), 10, 64)
	return ab23ao23Int, ab23ao23Err == nil
}

// Check escapes before decoding, which otherwise replaces unpaired surrogates.
// The JSON decoder remains responsible for the rest of the syntax validation.
func ab23ao23ValidSurrogates(input []byte) bool {
	for i := 0; i < len(input); i++ {
		if input[i] != '\\' {
			continue
		}
		i++
		if i >= len(input) || input[i] != 'u' {
			continue // Skip escaped backslashes and other non-Unicode escapes.
		}
		if i+4 >= len(input) {
			return false
		}
		code, err := strconv.ParseUint(string(input[i+1:i+5]), 16, 16)
		if err != nil {
			return false
		}
		i += 4
		if code >= 0xdc00 && code <= 0xdfff {
			return false
		}
		if code >= 0xd800 && code <= 0xdbff {
			if i+6 >= len(input) || input[i+1] != '\\' || input[i+2] != 'u' {
				return false
			}
			low, err := strconv.ParseUint(string(input[i+3:i+7]), 16, 16)
			if err != nil || low < 0xdc00 || low > 0xdfff {
				return false
			}
			i += 6
		}
	}
	return true
}

func Ab23Ao23Interpret(ab23ao23Input json.RawMessage) (json.RawMessage, error) {
	if !utf8.Valid(ab23ao23Input) {
		return nil, fmt.Errorf("input must be UTF-8 JSON")
	}
	if !ab23ao23ValidSurrogates(ab23ao23Input) {
		return nil, fmt.Errorf("input must contain valid Unicode escapes with paired surrogates")
	}
	ab23ao23Decoder := json.NewDecoder(bytes.NewReader(ab23ao23Input))
	ab23ao23Decoder.UseNumber()
	ab23ao23Value, ab23ao23Err := ab23ao23Value(ab23ao23Decoder)
	if ab23ao23Err != nil {
		return nil, ab23ao23Err
	}
	if _, ab23ao23Err = ab23ao23Decoder.Token(); ab23ao23Err != io.EOF {
		return nil, fmt.Errorf("trailing JSON data")
	}
	ab23ao23InputObject, ab23ao23OK := ab23ao23Value.(map[string]any)
	if !ab23ao23OK {
		return nil, fmt.Errorf("input must be an object")
	}
	for ab23ao23Key := range ab23ao23InputObject {
		if ab23ao23Key != "context" && ab23ao23Key != "fields" && ab23ao23Key != "arguments" {
			return nil, fmt.Errorf("unknown input key %q", ab23ao23Key)
		}
	}
	ab23ao23Context, ab23ao23OK := ab23ao23InputObject["context"].(map[string]any)
	if !ab23ao23OK {
		return nil, fmt.Errorf("context must be an object")
	}
	ab23ao23Raw, ab23ao23OK := ab23ao23InputObject["fields"].(map[string]any)
	if !ab23ao23OK {
		return nil, fmt.Errorf("fields must be an object")
	}
	for _, ab23ao23Key := range []string{"primary", "secondary", "operation"} {
		if ab23ao23Selector, ab23ao23Exists := ab23ao23Context[ab23ao23Key]; ab23ao23Exists {
			ab23ao23Number, ab23ao23OK := ab23ao23Selector.(json.Number)
			if !ab23ao23OK || strings.ContainsAny(string(ab23ao23Number), ".eE") {
				return nil, fmt.Errorf("%s must be an integer", ab23ao23Key)
			}
		}
	}
	for _, ab23ao23Key := range []string{"direction", "view"} {
		if ab23ao23Selector, ab23ao23Exists := ab23ao23Context[ab23ao23Key]; ab23ao23Exists {
			if _, ab23ao23OK := ab23ao23Selector.(string); !ab23ao23OK {
				return nil, fmt.Errorf("%s must be a string", ab23ao23Key)
			}
		}
	}
	for ab23ao23Key, ab23ao23Field := range ab23ao23Raw {
		ab23ao23Bytes, ab23ao23OK := ab23ao23Field.([]any)
		if !ab23ao23OK {
			return nil, fmt.Errorf("field %q must be a byte array", ab23ao23Key)
		}
		for _, ab23ao23Byte := range ab23ao23Bytes {
			ab23ao23Int, ab23ao23OK := ab23ao23Integer(ab23ao23Byte)
			if !ab23ao23OK || ab23ao23Int < 0 || ab23ao23Int > 255 {
				return nil, fmt.Errorf("field %q has an invalid byte", ab23ao23Key)
			}
		}
	}
	ab23ao23Fields := map[string]any{}
	ab23ao23Output := map[string]any{"context": ab23ao23Context, "raw": ab23ao23Raw, "state": "Unknown", "fields": ab23ao23Fields, "physical": "Unknown", "qualified": false}
	if ab23ao23Arguments, ab23ao23Exists := ab23ao23InputObject["arguments"]; ab23ao23Exists {
		ab23ao23Object, ab23ao23OK := ab23ao23Arguments.(map[string]any)
		if !ab23ao23OK || len(ab23ao23Object) != 0 {
			return nil, fmt.Errorf("arguments must be an empty object; no scalar arguments are defined")
		}
		ab23ao23Output["arguments"] = ab23ao23Object
	}
	ab23ao23Primary, _ := ab23ao23Integer(ab23ao23Context["primary"])
	ab23ao23Secondary, _ := ab23ao23Integer(ab23ao23Context["secondary"])
	ab23ao23Operation, ab23ao23HasOperation := ab23ao23Integer(ab23ao23Context["operation"])
	ab23ao23Direction, _ := ab23ao23Context["direction"].(string)
	if ab23ao23Primary != 181 || ab23ao23Secondary != 32 || !ab23ao23HasOperation || ab23ao23Operation != 0 || ab23ao23Context["view"] != "logical-object" || (ab23ao23Direction != "request" && ab23ao23Direction != "positive") {
		return json.Marshal(ab23ao23Output)
	}
	ab23ao23Output["state"] = "Hypothesis"
	ab23ao23Object, ab23ao23Present := ab23ao23Raw["logical_object"].([]any)
	ab23ao23Names := []string{"relay", "PWM1", "PWM2", "stepper", "LED", "ED", "temperature-table"}
	ab23ao23Offsets := []int{3, 4, 5, 6, 7, 8, 9}
	ab23ao23Widths := []int{1, 1, 1, 1, 1, 1, 1}
	if ab23ao23Direction == "positive" {
		ab23ao23Names = []string{"ADC", "temperature1", "temperature2", "temperature3", "temperature4", "reserved", "impeller", "vortex"}
		ab23ao23Offsets = []int{3, 5, 7, 9, 11, 13, 14, 16}
		ab23ao23Widths = []int{2, 2, 2, 2, 2, 1, 2, 2}
	}
	for ab23ao23Index, ab23ao23Name := range ab23ao23Names {
		ab23ao23Offset, ab23ao23Width := ab23ao23Offsets[ab23ao23Index], ab23ao23Widths[ab23ao23Index]
		ab23ao23Field := map[string]any{"state": "incomplete", "grade": "Hypothesis"}
		ab23ao23Fields[ab23ao23Name] = ab23ao23Field
		if !ab23ao23Present {
			ab23ao23Field["raw"] = nil
		} else {
			ab23ao23Start := min(ab23ao23Offset, len(ab23ao23Object))
			ab23ao23End := min(ab23ao23Offset+ab23ao23Width, len(ab23ao23Object))
			ab23ao23Field["raw"] = ab23ao23Object[ab23ao23Start:ab23ao23End]
		}
		if len(ab23ao23Object) < ab23ao23Offset+ab23ao23Width {
			ab23ao23Output["state"] = "partial"
			continue
		}
		ab23ao23Field["state"] = "Hypothesis"
		ab23ao23Low, _ := ab23ao23Integer(ab23ao23Object[ab23ao23Offset])
		ab23ao23Carrier := ab23ao23Low
		if ab23ao23Width == 2 {
			ab23ao23High, _ := ab23ao23Integer(ab23ao23Object[ab23ao23Offset+1])
			ab23ao23Carrier += ab23ao23High * 256
		}
		switch ab23ao23Name {
		case "PWM1", "PWM2", "ED":
			ab23ao23Field["number"] = float64(ab23ao23Carrier) / 2
			ab23ao23Field["calculation"] = "Proven"
		case "stepper":
			if ab23ao23Carrier <= 240 {
				ab23ao23Field["number"] = ab23ao23Carrier
			} else if ab23ao23Carrier == 255 {
				ab23ao23Field["label"] = "calibrate/go-to-zero"
			} else {
				ab23ao23Field["state"] = "unmapped"
			}
		case "LED":
			if ab23ao23Carrier <= 1 {
				ab23ao23Field["label"] = []string{"disabled", "enabled"}[ab23ao23Carrier]
			} else {
				ab23ao23Field["state"] = "unmapped"
			}
		case "temperature-table":
			if ab23ao23Carrier <= 3 {
				ab23ao23Field["label"] = []string{"standard3K3", "extended1K", "low-temperature15K", "standard1K"}[ab23ao23Carrier]
				ab23ao23Field["flags"] = map[string]any{"candidate_celsius_range": [][2]float64{{-13.5, 116}, {-20, 155}, {-39, 80}, {-13.9, 130}}[ab23ao23Carrier]}
			} else {
				ab23ao23Field["state"] = "unmapped"
			}
		case "temperature1", "temperature2", "temperature3", "temperature4":
			ab23ao23Field["number"] = float64(int16(ab23ao23Carrier)) / 16
			ab23ao23Field["calculation"] = "Proven"
			ab23ao23Field["candidate_domain"] = ab23ao23Carrier <= 3200
		case "ADC":
			ab23ao23Field["number"] = ab23ao23Carrier
			ab23ao23Field["candidate_domain"] = ab23ao23Carrier <= 1023
		case "impeller":
			ab23ao23Field["number"] = ab23ao23Carrier
			ab23ao23Field["candidate_domain"] = ab23ao23Carrier <= 1000
		case "vortex":
			ab23ao23Field["number"] = ab23ao23Carrier
			ab23ao23Field["candidate_domain"] = ab23ao23Carrier <= 350
			ab23ao23Field["flags"] = map[string]any{"candidate_operational_domain": ab23ao23Carrier >= 5 && ab23ao23Carrier <= 350, "candidate_hysteresis_hz": map[string]int{"on": 5, "off": 4}}
		case "reserved":
			ab23ao23Field["reserved"] = true
		}
		if ab23ao23Field["state"] == "unmapped" {
			ab23ao23Output["state"] = "partial"
		}
	}
	return json.Marshal(ab23ao23Output)
}
