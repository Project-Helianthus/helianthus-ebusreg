package vaillant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/big"
	"strconv"
	"strings"
)

func ab05aq05Read(ab05aq05Decoder *json.Decoder) (any, error) {
	ab05aq05Token, ab05aq05Error := ab05aq05Decoder.Token()
	if ab05aq05Error != nil {
		return nil, ab05aq05Error
	}
	switch ab05aq05Value := ab05aq05Token.(type) {
	case json.Delim:
		switch ab05aq05Value {
		case '{':
			ab05aq05Object := map[string]any{}
			for ab05aq05Decoder.More() {
				ab05aq05Key, ab05aq05Err := ab05aq05Decoder.Token()
				if ab05aq05Err != nil {
					return nil, ab05aq05Err
				}
				ab05aq05Name := ab05aq05Key.(string)
				if _, ab05aq05Exists := ab05aq05Object[ab05aq05Name]; ab05aq05Exists {
					return nil, fmt.Errorf("duplicate key %q", ab05aq05Name)
				}
				ab05aq05Child, ab05aq05Err := ab05aq05Read(ab05aq05Decoder)
				if ab05aq05Err != nil {
					return nil, ab05aq05Err
				}
				ab05aq05Object[ab05aq05Name] = ab05aq05Child
			}
			_, ab05aq05Error = ab05aq05Decoder.Token()
			return ab05aq05Object, ab05aq05Error
		case '[':
			ab05aq05Array := []any{}
			for ab05aq05Decoder.More() {
				ab05aq05Child, ab05aq05Err := ab05aq05Read(ab05aq05Decoder)
				if ab05aq05Err != nil {
					return nil, ab05aq05Err
				}
				ab05aq05Array = append(ab05aq05Array, ab05aq05Child)
			}
			_, ab05aq05Error = ab05aq05Decoder.Token()
			return ab05aq05Array, ab05aq05Error
		}
	case json.Number:
		if ab05aq05Err := ab05aq05Bound(string(ab05aq05Value)); ab05aq05Err != nil {
			return nil, ab05aq05Err
		}
		ab05aq05Number, ab05aq05Err := ab05aq05Value.Float64()
		if ab05aq05Err != nil || math.IsInf(ab05aq05Number, 0) || math.IsNaN(ab05aq05Number) {
			return nil, fmt.Errorf("nonfinite number")
		}
	}
	return ab05aq05Token, nil
}

// ab05aq05Bound rejects numeric tokens whose length or decimal exponent would
// make exact rational arithmetic and rendering unbounded.
func ab05aq05Bound(ab05aq05Token string) error {
	if len(ab05aq05Token) > 64 {
		return fmt.Errorf("number token too long")
	}
	if ab05aq05Index := strings.IndexAny(ab05aq05Token, "eE"); ab05aq05Index >= 0 {
		ab05aq05Exponent, ab05aq05Err := strconv.Atoi(ab05aq05Token[ab05aq05Index+1:])
		if ab05aq05Err != nil || ab05aq05Exponent > 400 || ab05aq05Exponent < -400 {
			return fmt.Errorf("number exponent out of range")
		}
	}
	return nil
}

func ab05aq05Byte(ab05aq05Value any) (int, bool) {
	ab05aq05Number, ab05aq05OK := ab05aq05Value.(json.Number)
	if !ab05aq05OK {
		return 0, false
	}
	ab05aq05Rational, ab05aq05Valid := new(big.Rat).SetString(string(ab05aq05Number))
	if !ab05aq05Valid || !ab05aq05Rational.IsInt() || !ab05aq05Rational.Num().IsInt64() {
		return 0, false
	}
	ab05aq05Integer := ab05aq05Rational.Num().Int64()
	return int(ab05aq05Integer), ab05aq05Integer >= 0 && ab05aq05Integer <= 255
}

func Ab05Aq05Evaluate(ab05aq05Input json.RawMessage) (json.RawMessage, error) {
	ab05aq05Decoder := json.NewDecoder(bytes.NewReader(ab05aq05Input))
	ab05aq05Decoder.UseNumber()
	ab05aq05Parsed, ab05aq05Error := ab05aq05Read(ab05aq05Decoder)
	if ab05aq05Error != nil {
		return nil, ab05aq05Error
	}
	if _, ab05aq05Error = ab05aq05Decoder.Token(); ab05aq05Error != io.EOF {
		return nil, fmt.Errorf("trailing JSON")
	}
	ab05aq05Output, ab05aq05OK := ab05aq05Parsed.(map[string]any)
	if !ab05aq05OK {
		return nil, fmt.Errorf("input must be object")
	}
	ab05aq05Context, ab05aq05OK := ab05aq05Output["context"].(map[string]any)
	if !ab05aq05OK {
		return nil, fmt.Errorf("context must be object")
	}
	ab05aq05Fields, ab05aq05OK := ab05aq05Output["fields"].(map[string]any)
	if !ab05aq05OK {
		return nil, fmt.Errorf("fields must be object")
	}
	ab05aq05Arguments := map[string]any{}
	if ab05aq05Value, ab05aq05Present := ab05aq05Output["arguments"]; ab05aq05Present {
		ab05aq05Arguments, ab05aq05OK = ab05aq05Value.(map[string]any)
		if !ab05aq05OK {
			return nil, fmt.Errorf("arguments must be object")
		}
	}
	for _, ab05aq05Key := range []string{"primary", "secondary", "operation"} {
		if ab05aq05Value, ab05aq05Present := ab05aq05Context[ab05aq05Key]; ab05aq05Present {
			if _, ab05aq05OK = ab05aq05Byte(ab05aq05Value); !ab05aq05OK {
				return nil, fmt.Errorf("invalid %s", ab05aq05Key)
			}
		}
	}
	for _, ab05aq05Key := range []string{"direction", "view", "family", "variant"} {
		if ab05aq05Value, ab05aq05Present := ab05aq05Context[ab05aq05Key]; ab05aq05Present {
			if _, ab05aq05OK = ab05aq05Value.(string); !ab05aq05OK {
				return nil, fmt.Errorf("invalid %s", ab05aq05Key)
			}
		}
	}
	ab05aq05Results := map[string]any{}
	for ab05aq05Name, ab05aq05Raw := range ab05aq05Fields {
		ab05aq05Array, ab05aq05Valid := ab05aq05Raw.([]any)
		if !ab05aq05Valid {
			return nil, fmt.Errorf("field %s must be byte array", ab05aq05Name)
		}
		for _, ab05aq05Value := range ab05aq05Array {
			if _, ab05aq05Valid = ab05aq05Byte(ab05aq05Value); !ab05aq05Valid {
				return nil, fmt.Errorf("invalid field byte")
			}
		}
		ab05aq05Results[ab05aq05Name] = map[string]any{"raw": ab05aq05Raw, "state": "Unknown"}
	}
	ab05aq05Output["state"] = "Unknown"
	ab05aq05Output["results"] = ab05aq05Results
	ab05aq05Output["physical"] = "Unknown"
	ab05aq05Output["qualified"] = false
	ab05aq05Primary, _ := ab05aq05Byte(ab05aq05Context["primary"])
	ab05aq05Secondary, _ := ab05aq05Byte(ab05aq05Context["secondary"])
	ab05aq05Operation, ab05aq05HasOperation := ab05aq05Byte(ab05aq05Context["operation"])
	ab05aq05View, _ := ab05aq05Context["view"].(string)
	ab05aq05Supported := ab05aq05Primary == 181 && ab05aq05Secondary == 5 && ab05aq05Context["direction"] == "request" && ab05aq05HasOperation
	ab05aq05Supported = ab05aq05Supported && ((ab05aq05View == "unchanged" && (ab05aq05Operation == 90 || ab05aq05Operation == 61)) || (ab05aq05View == "dhw-enable" && ab05aq05Operation == 49) || (ab05aq05View == "parallel-charge" && (ab05aq05Operation == 24 || ab05aq05Operation == 39)) || (ab05aq05View == "flow-excess" && ab05aq05Operation == 23) || (ab05aq05View == "programmed-start" && ab05aq05Operation == 9))
	if !ab05aq05Supported {
		return json.Marshal(ab05aq05Output)
	}
	ab05aq05Output["state"] = "complete"
	ab05aq05Result := map[string]any{"state": "complete", "grade": "Hypothesis"}
	ab05aq05Name := ab05aq05View
	switch ab05aq05View {
	case "unchanged":
		ab05aq05Replacement, ab05aq05Valid := ab05aq05Arguments["replacement"].(bool)
		if !ab05aq05Valid {
			return nil, fmt.Errorf("replacement must be boolean")
		}
		for _, ab05aq05Key := range []string{"previous", "requested"} {
			ab05aq05Value, ab05aq05Present := ab05aq05Arguments[ab05aq05Key]
			if !ab05aq05Present {
				return nil, fmt.Errorf("missing %s", ab05aq05Key)
			}
			switch ab05aq05Value.(type) {
			case map[string]any, []any:
				return nil, fmt.Errorf("%s must be scalar", ab05aq05Key)
			}
		}
		ab05aq05Result["raw"] = ab05aq05Arguments["requested"]
		ab05aq05Result["value"] = ab05aq05Arguments["requested"]
		ab05aq05Result["label"] = "numeric-candidate"
		ab05aq05Result["replacement_code"] = "Unknown"
		if ab05aq05Replacement {
			ab05aq05Result["value"] = ab05aq05Arguments["previous"]
			ab05aq05Result["label"] = "unchanged"
		}
		if ab05aq05Number, ab05aq05Numeric := ab05aq05Result["value"].(json.Number); ab05aq05Numeric {
			ab05aq05Result["number"] = ab05aq05Number
		}
	case "flow-excess":
		ab05aq05Demand, ab05aq05Valid := ab05aq05Arguments["demand"].(json.Number)
		if !ab05aq05Valid {
			return nil, fmt.Errorf("demand must be numeric")
		}
		ab05aq05Excess, ab05aq05Valid := ab05aq05Arguments["excess"].(json.Number)
		if !ab05aq05Valid {
			return nil, fmt.Errorf("excess must be numeric")
		}
		ab05aq05Left, ab05aq05Valid := new(big.Rat).SetString(string(ab05aq05Demand))
		if !ab05aq05Valid {
			return nil, fmt.Errorf("invalid demand")
		}
		ab05aq05Right, ab05aq05Valid := new(big.Rat).SetString(string(ab05aq05Excess))
		if !ab05aq05Valid {
			return nil, fmt.Errorf("invalid excess")
		}
		ab05aq05Sum := new(big.Rat).Add(ab05aq05Left, ab05aq05Right)
		ab05aq05Denominator := new(big.Int).Set(ab05aq05Sum.Denom())
		ab05aq05Scale := 0
		for ab05aq05Denominator.Cmp(big.NewInt(1)) != 0 {
			ab05aq05Divisor := int64(2)
			if ab05aq05Denominator.Bit(0) != 0 {
				ab05aq05Divisor = 5
			}
			ab05aq05Denominator.Quo(ab05aq05Denominator, big.NewInt(ab05aq05Divisor))
			ab05aq05Scale++
		}
		ab05aq05Result["raw"] = map[string]any{"demand": ab05aq05Demand, "excess": ab05aq05Excess}
		ab05aq05Result["number"] = json.Number(ab05aq05Sum.FloatString(ab05aq05Scale))
		ab05aq05Result["calculation"] = "Proven"
	default:
		ab05aq05Name = "enabled"
		if ab05aq05View == "parallel-charge" {
			ab05aq05Name = "disable"
		}
		if ab05aq05View == "programmed-start" {
			ab05aq05Name = "start"
		}
		ab05aq05Role := ""
		ab05aq05Parallel := false
		ab05aq05Mode := 0
		if ab05aq05View == "parallel-charge" {
			ab05aq05Role, ab05aq05OK = ab05aq05Arguments["role"].(string)
			if !ab05aq05OK || (ab05aq05Role != "mixer" && ab05aq05Role != "heating" && ab05aq05Role != "dhw") {
				return nil, fmt.Errorf("invalid role")
			}
			ab05aq05Parallel, ab05aq05OK = ab05aq05Arguments["parallel"].(bool)
			if !ab05aq05OK {
				return nil, fmt.Errorf("parallel must be boolean")
			}
		}
		if ab05aq05View == "programmed-start" {
			ab05aq05Mode, ab05aq05OK = ab05aq05Byte(ab05aq05Arguments["program_mode"])
			if !ab05aq05OK {
				return nil, fmt.Errorf("invalid program_mode")
			}
		}
		ab05aq05Raw, ab05aq05Present := ab05aq05Fields[ab05aq05Name]
		if ab05aq05Present {
			ab05aq05Result["raw"] = ab05aq05Raw
		}
		ab05aq05Array, _ := ab05aq05Raw.([]any)
		if !ab05aq05Present || (ab05aq05View != "programmed-start" && len(ab05aq05Array) != 1) {
			ab05aq05Result["state"] = "incomplete"
			ab05aq05Output["state"] = "partial"
			break
		}
		if ab05aq05View == "programmed-start" {
			ab05aq05Result["label"] = "ignored"
			if ab05aq05Mode == 0 {
				ab05aq05Result["label"] = "applicable"
			}
			break
		}
		ab05aq05ByteValue, _ := ab05aq05Byte(ab05aq05Array[0])
		ab05aq05Result["label"] = "unmapped"
		if ab05aq05View == "dhw-enable" {
			if ab05aq05ByteValue == 0 {
				ab05aq05Result["label"] = "frost-only-eligible"
			}
			if ab05aq05ByteValue == 1 {
				ab05aq05Result["label"] = "ordinary-DHW-eligible"
			}
			break
		}
		switch {
		case ab05aq05Role == "dhw":
			ab05aq05Result["label"] = "not-applicable"
		case ab05aq05ByteValue == 0:
			ab05aq05Result["label"] = "suppression-not-requested"
		case ab05aq05ByteValue == 1:
			ab05aq05Result["label"] = "heating-suppressed"
			if ab05aq05Parallel {
				ab05aq05Result["label"] = "exemption-unknown"
				if ab05aq05Role == "mixer" {
					ab05aq05Result["label"] = "mixer-parallel-exception"
				}
			}
		}
	}
	ab05aq05Results[ab05aq05Name] = ab05aq05Result
	return json.Marshal(ab05aq05Output)
}
