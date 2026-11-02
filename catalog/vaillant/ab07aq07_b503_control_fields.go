package vaillant

import (
	ab07aq07bytes "bytes"
	ab07aq07json "encoding/json"
	ab07aq07fmt "fmt"
	ab07aq07io "io"
	ab07aq07strconv "strconv"
)

func ab07aq07Decode(ab07aq07Decoder *ab07aq07json.Decoder) (any, error) {
	ab07aq07Token, ab07aq07Error := ab07aq07Decoder.Token()
	if ab07aq07Error != nil {
		return nil, ab07aq07Error
	}
	if ab07aq07Delimiter, ab07aq07OK := ab07aq07Token.(ab07aq07json.Delim); ab07aq07OK {
		switch ab07aq07Delimiter {
		case '{':
			ab07aq07Object := map[string]any{}
			for ab07aq07Decoder.More() {
				ab07aq07KeyToken, ab07aq07Error := ab07aq07Decoder.Token()
				if ab07aq07Error != nil {
					return nil, ab07aq07Error
				}
				ab07aq07Key, ab07aq07OK := ab07aq07KeyToken.(string)
				if !ab07aq07OK {
					return nil, ab07aq07fmt.Errorf("invalid object key")
				}
				if _, ab07aq07Exists := ab07aq07Object[ab07aq07Key]; ab07aq07Exists {
					return nil, ab07aq07fmt.Errorf("duplicate key %q", ab07aq07Key)
				}
				ab07aq07Value, ab07aq07Error := ab07aq07Decode(ab07aq07Decoder)
				if ab07aq07Error != nil {
					return nil, ab07aq07Error
				}
				ab07aq07Object[ab07aq07Key] = ab07aq07Value
			}
			_, ab07aq07Error := ab07aq07Decoder.Token()
			return ab07aq07Object, ab07aq07Error
		case '[':
			ab07aq07Array := []any{}
			for ab07aq07Decoder.More() {
				ab07aq07Value, ab07aq07Error := ab07aq07Decode(ab07aq07Decoder)
				if ab07aq07Error != nil {
					return nil, ab07aq07Error
				}
				ab07aq07Array = append(ab07aq07Array, ab07aq07Value)
			}
			_, ab07aq07Error := ab07aq07Decoder.Token()
			return ab07aq07Array, ab07aq07Error
		}
		return nil, ab07aq07fmt.Errorf("unexpected delimiter")
	}
	return ab07aq07Token, nil
}

func ab07aq07Uint(ab07aq07Value any, ab07aq07Bits int) (uint64, error) {
	ab07aq07Number, ab07aq07OK := ab07aq07Value.(ab07aq07json.Number)
	if !ab07aq07OK {
		return 0, ab07aq07fmt.Errorf("unsigned integer required")
	}
	return ab07aq07strconv.ParseUint(string(ab07aq07Number), 10, ab07aq07Bits)
}

func Ab07Aq07Evaluate(ab07aq07Input ab07aq07json.RawMessage) (ab07aq07json.RawMessage, error) {
	ab07aq07Decoder := ab07aq07json.NewDecoder(ab07aq07bytes.NewReader(ab07aq07Input))
	ab07aq07Decoder.UseNumber()
	ab07aq07Value, ab07aq07Error := ab07aq07Decode(ab07aq07Decoder)
	if ab07aq07Error != nil {
		return nil, ab07aq07Error
	}
	if _, ab07aq07Error = ab07aq07Decoder.Token(); ab07aq07Error != ab07aq07io.EOF {
		return nil, ab07aq07fmt.Errorf("trailing JSON")
	}
	ab07aq07Output, ab07aq07OK := ab07aq07Value.(map[string]any)
	if !ab07aq07OK {
		return nil, ab07aq07fmt.Errorf("input object required")
	}
	ab07aq07Context, ab07aq07OK := ab07aq07Output["context"].(map[string]any)
	if !ab07aq07OK {
		return nil, ab07aq07fmt.Errorf("context object required")
	}
	ab07aq07Fields, ab07aq07OK := ab07aq07Output["fields"].(map[string]any)
	if !ab07aq07OK {
		return nil, ab07aq07fmt.Errorf("fields object required")
	}
	for ab07aq07Key, ab07aq07Value := range ab07aq07Fields {
		ab07aq07Array, ab07aq07OK := ab07aq07Value.([]any)
		if !ab07aq07OK {
			return nil, ab07aq07fmt.Errorf("field %s requires byte array", ab07aq07Key)
		}
		for _, ab07aq07Byte := range ab07aq07Array {
			if _, ab07aq07Error := ab07aq07Uint(ab07aq07Byte, 8); ab07aq07Error != nil {
				return nil, ab07aq07Error
			}
		}
	}
	if ab07aq07Value, ab07aq07Present := ab07aq07Output["arguments"]; ab07aq07Present {
		ab07aq07Arguments, ab07aq07OK := ab07aq07Value.(map[string]any)
		if !ab07aq07OK || len(ab07aq07Arguments) != 0 {
			return nil, ab07aq07fmt.Errorf("no arguments defined for these views")
		}
	}
	ab07aq07Supported := true
	ab07aq07Numbers := map[string]uint64{}
	for _, ab07aq07Key := range []string{"primary", "secondary", "operation"} {
		ab07aq07Value, ab07aq07Present := ab07aq07Context[ab07aq07Key]
		if !ab07aq07Present {
			ab07aq07Supported = false
			continue
		}
		ab07aq07Number, ab07aq07Error := ab07aq07Uint(ab07aq07Value, 8)
		if ab07aq07Error != nil {
			return nil, ab07aq07Error
		}
		ab07aq07Numbers[ab07aq07Key] = ab07aq07Number
	}
	for _, ab07aq07Key := range []string{"direction", "view"} {
		if ab07aq07Value, ab07aq07Present := ab07aq07Context[ab07aq07Key]; ab07aq07Present {
			if _, ab07aq07OK := ab07aq07Value.(string); !ab07aq07OK {
				return nil, ab07aq07fmt.Errorf("context %s requires string", ab07aq07Key)
			}
		} else {
			ab07aq07Supported = false
		}
	}
	ab07aq07View, _ := ab07aq07Context["view"].(string)
	ab07aq07Direction, _ := ab07aq07Context["direction"].(string)
	ab07aq07Operation := ab07aq07Numbers["operation"]
	ab07aq07Match := false
	switch ab07aq07View {
	case "reset-acceptance":
		ab07aq07Match = (ab07aq07Operation == 3 || ab07aq07Operation == 6) && ab07aq07Direction == "positive"
	case "current-errors":
		ab07aq07Match = ab07aq07Operation == 7 && ab07aq07Direction == "positive"
	case "category-common":
		ab07aq07Match = ab07aq07Operation == 4 && ab07aq07Direction == "request"
	case "category-reset":
		ab07aq07Match = (ab07aq07Operation == 3 || ab07aq07Operation == 6) && ab07aq07Direction == "request"
	}
	ab07aq07Results := map[string]any{}
	ab07aq07Output["results"] = ab07aq07Results
	ab07aq07Output["state"] = "Unknown"
	ab07aq07Output["physical"] = "Unknown"
	ab07aq07Output["qualified"] = false
	if !ab07aq07Supported || !ab07aq07Match || ab07aq07Numbers["primary"] != 181 || ab07aq07Numbers["secondary"] != 3 {
		return ab07aq07json.Marshal(ab07aq07Output)
	}
	ab07aq07Output["state"] = "complete"
	ab07aq07Widths := map[string]int{}
	switch ab07aq07View {
	case "reset-acceptance":
		ab07aq07Widths["acceptance"] = 1
	case "current-errors":
		ab07aq07Widths["flags"] = 1
		ab07aq07Widths["feedback"] = 1
		for ab07aq07Index := 0; ab07aq07Index < 5; ab07aq07Index++ {
			ab07aq07Widths[ab07aq07fmt.Sprintf("error%d", ab07aq07Index)] = 2
		}
	default:
		ab07aq07Widths["category"] = 1
	}
	for ab07aq07Key, ab07aq07Width := range ab07aq07Widths {
		ab07aq07Raw, ab07aq07Present := ab07aq07Fields[ab07aq07Key]
		ab07aq07Result := map[string]any{"raw": ab07aq07Raw, "present": ab07aq07Present, "state": "incomplete", "grade": "Hypothesis"}
		ab07aq07Results[ab07aq07Key] = ab07aq07Result
		ab07aq07Array, _ := ab07aq07Raw.([]any)
		if !ab07aq07Present || len(ab07aq07Array) != ab07aq07Width {
			ab07aq07Output["state"] = "partial"
			continue
		}
		ab07aq07Number, _ := ab07aq07Uint(ab07aq07Array[0], 8)
		if ab07aq07Width == 2 {
			ab07aq07High, _ := ab07aq07Uint(ab07aq07Array[1], 8)
			ab07aq07Number += ab07aq07High << 8
		}
		ab07aq07Result["state"] = "Hypothesis"
		ab07aq07Result["number"] = ab07aq07Number
		ab07aq07Result["calculation"] = "Proven"
		switch ab07aq07Key {
		case "acceptance":
			ab07aq07Result["default"] = 0
			if ab07aq07Number <= 1 {
				ab07aq07Result["value"] = ab07aq07Number == 1
			} else {
				ab07aq07Result["state"] = "Unknown"
				ab07aq07Result["unmodeled"] = ab07aq07Number
			}
		case "feedback":
			ab07aq07Result["label"] = "rejected"
			if ab07aq07Number == 1 {
				ab07aq07Result["label"] = "accepted"
			}
		case "flags":
			ab07aq07Flags := map[string]bool{}
			for ab07aq07Index := 0; ab07aq07Index < 5; ab07aq07Index++ {
				ab07aq07Flags[ab07aq07fmt.Sprintf("bit%d", ab07aq07Index)] = ab07aq07Number&(1<<ab07aq07Index) != 0
			}
			ab07aq07Result["flags"] = ab07aq07Flags
			ab07aq07Result["reserved"] = ab07aq07Number & 0xe0
		case "category":
			ab07aq07Result["state"] = "Unknown"
			if ab07aq07View == "category-common" {
				ab07aq07Result["candidate_domain"] = []int{1, 2, 3, 4, 5}
				if ab07aq07Number == 2 {
					ab07aq07Result["label"] = "predictive-maintenance"
					ab07aq07Result["state"] = "Hypothesis"
				}
			} else {
				if ab07aq07Number == 1 {
					ab07aq07Result["label"] = "error"
					ab07aq07Result["state"] = "Hypothesis"
				}
				if ab07aq07Number == 3 {
					ab07aq07Result["label"] = "status"
					ab07aq07Result["state"] = "Hypothesis"
				}
			}
		}
	}
	return ab07aq07json.Marshal(ab07aq07Output)
}
