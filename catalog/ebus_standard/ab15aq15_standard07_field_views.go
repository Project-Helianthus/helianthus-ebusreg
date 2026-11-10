package ebus_standard_catalog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/big"
	"strconv"
)

func ab15aq15Read(ab15aq15Decoder *json.Decoder, ab15aq15Depth int) (any, error) {
	if ab15aq15Depth > 64 {
		return nil, fmt.Errorf("nesting too deep")
	}
	ab15aq15Token, ab15aq15Err := ab15aq15Decoder.Token()
	if ab15aq15Err != nil {
		return nil, ab15aq15Err
	}
	if ab15aq15Number, ab15aq15OK := ab15aq15Token.(json.Number); ab15aq15OK {
		ab15aq15Value, ab15aq15Err := strconv.ParseFloat(string(ab15aq15Number), 64)
		if ab15aq15Err != nil || math.IsInf(ab15aq15Value, 0) || math.IsNaN(ab15aq15Value) {
			return nil, fmt.Errorf("invalid number")
		}
	}
	if ab15aq15Delimiter, ab15aq15OK := ab15aq15Token.(json.Delim); ab15aq15OK {
		if ab15aq15Delimiter == '{' {
			ab15aq15Object := map[string]any{}
			for ab15aq15Decoder.More() {
				ab15aq15KeyToken, ab15aq15Err := ab15aq15Decoder.Token()
				if ab15aq15Err != nil {
					return nil, ab15aq15Err
				}
				ab15aq15Key, ab15aq15OK := ab15aq15KeyToken.(string)
				if !ab15aq15OK {
					return nil, fmt.Errorf("invalid key")
				}
				if _, ab15aq15Found := ab15aq15Object[ab15aq15Key]; ab15aq15Found {
					return nil, fmt.Errorf("duplicate key %q", ab15aq15Key)
				}
				ab15aq15Object[ab15aq15Key], ab15aq15Err = ab15aq15Read(ab15aq15Decoder, ab15aq15Depth+1)
				if ab15aq15Err != nil {
					return nil, ab15aq15Err
				}
			}
			_, ab15aq15Err := ab15aq15Decoder.Token()
			return ab15aq15Object, ab15aq15Err
		}
		if ab15aq15Delimiter == '[' {
			ab15aq15Array := []any{}
			for ab15aq15Decoder.More() {
				ab15aq15Value, ab15aq15Err := ab15aq15Read(ab15aq15Decoder, ab15aq15Depth+1)
				if ab15aq15Err != nil {
					return nil, ab15aq15Err
				}
				ab15aq15Array = append(ab15aq15Array, ab15aq15Value)
			}
			_, ab15aq15Err := ab15aq15Decoder.Token()
			return ab15aq15Array, ab15aq15Err
		}
		return nil, fmt.Errorf("unexpected delimiter")
	}
	return ab15aq15Token, nil
}

func ab15aq15Byte(ab15aq15Value any) (int, bool) {
	ab15aq15Number, ab15aq15OK := ab15aq15Value.(json.Number)
	if !ab15aq15OK {
		return 0, false
	}
	ab15aq15Rational, ab15aq15OK := new(big.Rat).SetString(string(ab15aq15Number))
	if !ab15aq15OK || !ab15aq15Rational.IsInt() || !ab15aq15Rational.Num().IsInt64() {
		return 0, false
	}
	ab15aq15Integer := ab15aq15Rational.Num().Int64()
	return int(ab15aq15Integer), ab15aq15Integer >= 0 && ab15aq15Integer <= 255
}

func Ab15Aq15Evaluate(ab15aq15Input json.RawMessage) (json.RawMessage, error) {
	ab15aq15Decoder := json.NewDecoder(bytes.NewReader(ab15aq15Input))
	ab15aq15Decoder.UseNumber()
	ab15aq15Value, ab15aq15Err := ab15aq15Read(ab15aq15Decoder, 0)
	if ab15aq15Err != nil {
		return nil, ab15aq15Err
	}
	if _, ab15aq15Err = ab15aq15Decoder.Token(); ab15aq15Err != io.EOF {
		return nil, fmt.Errorf("trailing JSON")
	}
	ab15aq15Output, ab15aq15OK := ab15aq15Value.(map[string]any)
	if !ab15aq15OK {
		return nil, fmt.Errorf("input must be object")
	}
	ab15aq15Context, ab15aq15OK := ab15aq15Output["context"].(map[string]any)
	if !ab15aq15OK {
		return nil, fmt.Errorf("context must be object")
	}
	ab15aq15Fields, ab15aq15OK := ab15aq15Output["fields"].(map[string]any)
	if !ab15aq15OK {
		return nil, fmt.Errorf("fields must be object")
	}
	for _, ab15aq15Key := range []string{"primary", "secondary", "operation"} {
		if ab15aq15Value, ab15aq15Found := ab15aq15Context[ab15aq15Key]; ab15aq15Found {
			if _, ab15aq15OK := ab15aq15Byte(ab15aq15Value); !ab15aq15OK {
				return nil, fmt.Errorf("invalid %s", ab15aq15Key)
			}
		}
	}
	for _, ab15aq15Key := range []string{"direction", "view"} {
		if ab15aq15Value, ab15aq15Found := ab15aq15Context[ab15aq15Key]; ab15aq15Found {
			if _, ab15aq15OK := ab15aq15Value.(string); !ab15aq15OK {
				return nil, fmt.Errorf("invalid %s", ab15aq15Key)
			}
		}
	}
	for _, ab15aq15Key := range []string{"results", "state", "physical", "qualified"} {
		if _, ab15aq15Found := ab15aq15Output[ab15aq15Key]; ab15aq15Found {
			return nil, fmt.Errorf("reserved key %q in input", ab15aq15Key)
		}
	}
	ab15aq15HasArguments := false
	if ab15aq15Value, ab15aq15Found := ab15aq15Output["arguments"]; ab15aq15Found {
		ab15aq15Arguments, ab15aq15OK := ab15aq15Value.(map[string]any)
		if !ab15aq15OK {
			return nil, fmt.Errorf("arguments must be object")
		}
		ab15aq15HasArguments = len(ab15aq15Arguments) != 0
	}
	ab15aq15Results := map[string]any{}
	for ab15aq15Key, ab15aq15Value := range ab15aq15Fields {
		ab15aq15Array, ab15aq15OK := ab15aq15Value.([]any)
		if !ab15aq15OK {
			return nil, fmt.Errorf("field %s must be byte array", ab15aq15Key)
		}
		for _, ab15aq15Value := range ab15aq15Array {
			if _, ab15aq15OK := ab15aq15Byte(ab15aq15Value); !ab15aq15OK {
				return nil, fmt.Errorf("invalid field byte")
			}
		}
		ab15aq15Results[ab15aq15Key] = map[string]any{"raw": ab15aq15Value, "state": "Unknown"}
	}
	ab15aq15Output["results"] = ab15aq15Results
	ab15aq15Output["state"] = "Unknown"
	ab15aq15Output["physical"] = "Unknown"
	ab15aq15Output["qualified"] = false
	ab15aq15Primary, _ := ab15aq15Byte(ab15aq15Context["primary"])
	ab15aq15Secondary, ab15aq15HasSecondary := ab15aq15Byte(ab15aq15Context["secondary"])
	_, ab15aq15HasOperation := ab15aq15Context["operation"]
	ab15aq15Direction, _ := ab15aq15Context["direction"].(string)
	if ab15aq15Primary != 7 || !ab15aq15HasSecondary || (ab15aq15Secondary != 0 && ab15aq15Secondary != 2 && ab15aq15Secondary != 3 && ab15aq15Secondary != 5) || ab15aq15HasOperation || ab15aq15Context["view"] != "field-candidates" || (ab15aq15Direction != "request" && ab15aq15Direction != "positive") {
		return json.Marshal(ab15aq15Output)
	}
	if ab15aq15HasArguments {
		return nil, fmt.Errorf("no arguments defined for this view")
	}
	ab15aq15Widths := map[string]int{}
	if ab15aq15Direction == "positive" {
		ab15aq15Widths = map[string]int{"version": 1, "revision": 1, "bitmap": 8}
	} else {
		switch ab15aq15Secondary {
		case 0:
			ab15aq15Widths = map[string]int{"outdoor_temperature": 2, "weekday": 1, "second": 1, "minute": 1, "hour": 1, "day": 1, "month": 1, "year": 1}
		case 2:
			ab15aq15Widths = map[string]int{"outdoor_temperature": 2, "duration": 1}
		case 3:
			ab15aq15Widths = map[string]int{"block": 1}
		case 5:
			ab15aq15Widths = map[string]int{"secondary_block": 1, "primary_block": 1}
		}
	}
	ab15aq15Output["state"] = "complete"
	for ab15aq15Key, ab15aq15Width := range ab15aq15Widths {
		ab15aq15Raw, ab15aq15Present := ab15aq15Fields[ab15aq15Key]
		ab15aq15Result := map[string]any{"state": "incomplete", "grade": "Hypothesis"}
		if ab15aq15Present {
			ab15aq15Result["raw"] = ab15aq15Raw
		} else {
			ab15aq15Result["absent"] = true
		}
		ab15aq15Results[ab15aq15Key] = ab15aq15Result
		ab15aq15Array, _ := ab15aq15Raw.([]any)
		if !ab15aq15Present || len(ab15aq15Array) != ab15aq15Width {
			ab15aq15Output["state"] = "partial"
			continue
		}
		ab15aq15Result["state"] = "complete"
		ab15aq15First, _ := ab15aq15Byte(ab15aq15Array[0])
		if ab15aq15Direction == "positive" {
			continue
		}
		switch ab15aq15Key {
		case "outdoor_temperature":
			ab15aq15Second, _ := ab15aq15Byte(ab15aq15Array[1])
			ab15aq15Signed := int(int16(uint16(ab15aq15First | ab15aq15Second<<8)))
			ab15aq15Result["number"] = float64(ab15aq15Signed) / 256
			ab15aq15Result["unit"] = "Celsius"
			ab15aq15Result["calculation"] = "Proven"
			if ab15aq15Secondary == 2 {
				ab15aq15Result["candidate_domain"] = ab15aq15Signed >= -12800 && ab15aq15Signed <= 12800
			}
		case "weekday":
			ab15aq15Labels := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
			if ab15aq15First < 7 {
				ab15aq15Result["label"] = ab15aq15Labels[ab15aq15First]
			} else {
				ab15aq15Result["label"] = "unmapped"
			}
		case "duration":
			ab15aq15Result["number"] = ab15aq15First
			ab15aq15Result["unit"] = "minutes"
			ab15aq15Result["calculation"] = "Proven"
		case "second", "minute", "hour", "day", "month", "year":
			ab15aq15Domains := map[string][2]int{"second": {0, 59}, "minute": {0, 59}, "hour": {0, 23}, "day": {1, 31}, "month": {1, 12}, "year": {0, 99}}
			ab15aq15Result["candidate_domain"] = ab15aq15Domains[ab15aq15Key]
			ab15aq15Result["conversion"] = "Unknown"
		}
	}
	return json.Marshal(ab15aq15Output)
}
