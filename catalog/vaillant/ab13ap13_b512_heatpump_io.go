package vaillant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

func ab13ap13Read(ab13ap13Decoder *json.Decoder) (any, error) {
	ab13ap13Token, ab13ap13Err := ab13ap13Decoder.Token()
	if ab13ap13Err != nil {
		return nil, ab13ap13Err
	}
	if ab13ap13Delimiter, ab13ap13OK := ab13ap13Token.(json.Delim); ab13ap13OK {
		switch ab13ap13Delimiter {
		case '{':
			ab13ap13Object := map[string]any{}
			for ab13ap13Decoder.More() {
				ab13ap13KeyToken, ab13ap13Err := ab13ap13Decoder.Token()
				if ab13ap13Err != nil {
					return nil, ab13ap13Err
				}
				ab13ap13Key, ab13ap13OK := ab13ap13KeyToken.(string)
				if !ab13ap13OK {
					return nil, fmt.Errorf("object key required")
				}
				if _, ab13ap13Exists := ab13ap13Object[ab13ap13Key]; ab13ap13Exists {
					return nil, fmt.Errorf("duplicate key %q", ab13ap13Key)
				}
				ab13ap13Value, ab13ap13Err := ab13ap13Read(ab13ap13Decoder)
				if ab13ap13Err != nil {
					return nil, ab13ap13Err
				}
				ab13ap13Object[ab13ap13Key] = ab13ap13Value
			}
			_, ab13ap13Err = ab13ap13Decoder.Token()
			return ab13ap13Object, ab13ap13Err
		case '[':
			ab13ap13Array := []any{}
			for ab13ap13Decoder.More() {
				ab13ap13Value, ab13ap13Err := ab13ap13Read(ab13ap13Decoder)
				if ab13ap13Err != nil {
					return nil, ab13ap13Err
				}
				ab13ap13Array = append(ab13ap13Array, ab13ap13Value)
			}
			_, ab13ap13Err = ab13ap13Decoder.Token()
			return ab13ap13Array, ab13ap13Err
		default:
			return nil, fmt.Errorf("unexpected delimiter")
		}
	}
	return ab13ap13Token, nil
}

func ab13ap13Uint(ab13ap13Value any, ab13ap13Bits int) (uint64, error) {
	ab13ap13Number, ab13ap13OK := ab13ap13Value.(json.Number)
	if !ab13ap13OK {
		return 0, fmt.Errorf("unsigned integer required")
	}
	return strconv.ParseUint(string(ab13ap13Number), 10, ab13ap13Bits)
}

func ab13ap13ValidateArguments(ab13ap13Value any) error {
	switch ab13ap13Item := ab13ap13Value.(type) {
	case json.Number:
		if _, ab13ap13Err := ab13ap13Uint(ab13ap13Item, 32); ab13ap13Err != nil {
			return fmt.Errorf("argument number: %w", ab13ap13Err)
		}
	case map[string]any:
		for _, ab13ap13Child := range ab13ap13Item {
			if ab13ap13Err := ab13ap13ValidateArguments(ab13ap13Child); ab13ap13Err != nil {
				return ab13ap13Err
			}
		}
	case []any:
		for _, ab13ap13Child := range ab13ap13Item {
			if ab13ap13Err := ab13ap13ValidateArguments(ab13ap13Child); ab13ap13Err != nil {
				return ab13ap13Err
			}
		}
	}
	return nil
}

func Ab13Ap13Evaluate(ab13ap13Input json.RawMessage) (json.RawMessage, error) {
	ab13ap13Decoder := json.NewDecoder(bytes.NewReader(ab13ap13Input))
	ab13ap13Decoder.UseNumber()
	ab13ap13Value, ab13ap13Err := ab13ap13Read(ab13ap13Decoder)
	if ab13ap13Err != nil {
		return nil, ab13ap13Err
	}
	if _, ab13ap13Err = ab13ap13Decoder.Token(); ab13ap13Err != io.EOF {
		return nil, fmt.Errorf("trailing JSON")
	}
	ab13ap13Root, ab13ap13OK := ab13ap13Value.(map[string]any)
	if !ab13ap13OK {
		return nil, fmt.Errorf("input object required")
	}
	ab13ap13Context, ab13ap13OK := ab13ap13Root["context"].(map[string]any)
	if !ab13ap13OK {
		return nil, fmt.Errorf("context object required")
	}
	ab13ap13Raw, ab13ap13OK := ab13ap13Root["fields"].(map[string]any)
	if !ab13ap13OK {
		return nil, fmt.Errorf("fields object required")
	}
	for _, ab13ap13Key := range []string{"primary", "secondary", "operation"} {
		if ab13ap13Item, ab13ap13Exists := ab13ap13Context[ab13ap13Key]; ab13ap13Exists && ab13ap13Item != nil {
			if _, ab13ap13Err := ab13ap13Uint(ab13ap13Item, 8); ab13ap13Err != nil {
				return nil, fmt.Errorf("context %s: %w", ab13ap13Key, ab13ap13Err)
			}
		}
	}
	for _, ab13ap13Key := range []string{"view", "direction"} {
		if ab13ap13Item, ab13ap13Exists := ab13ap13Context[ab13ap13Key]; ab13ap13Exists && ab13ap13Item != nil {
			if _, ab13ap13OK := ab13ap13Item.(string); !ab13ap13OK {
				return nil, fmt.Errorf("context %s: string required", ab13ap13Key)
			}
		}
	}
	for ab13ap13Key, ab13ap13Item := range ab13ap13Raw {
		ab13ap13Array, ab13ap13OK := ab13ap13Item.([]any)
		if !ab13ap13OK {
			return nil, fmt.Errorf("field %s: byte array required", ab13ap13Key)
		}
		for _, ab13ap13Byte := range ab13ap13Array {
			if _, ab13ap13Err := ab13ap13Uint(ab13ap13Byte, 8); ab13ap13Err != nil {
				return nil, fmt.Errorf("field %s: %w", ab13ap13Key, ab13ap13Err)
			}
		}
	}
	ab13ap13Fields := map[string]any{}
	ab13ap13Output := map[string]any{"context": ab13ap13Context, "raw": ab13ap13Raw, "state": "Unknown", "fields": ab13ap13Fields, "physical": "Unknown", "qualified": false}
	if ab13ap13Arguments, ab13ap13Exists := ab13ap13Root["arguments"]; ab13ap13Exists {
		if _, ab13ap13OK := ab13ap13Arguments.(map[string]any); !ab13ap13OK {
			return nil, fmt.Errorf("arguments object required")
		}
		if ab13ap13Err := ab13ap13ValidateArguments(ab13ap13Arguments); ab13ap13Err != nil {
			return nil, ab13ap13Err
		}
		ab13ap13Output["arguments"] = ab13ap13Arguments
	}
	ab13ap13Primary, _ := ab13ap13Uint(ab13ap13Context["primary"], 8)
	ab13ap13Secondary, _ := ab13ap13Uint(ab13ap13Context["secondary"], 8)
	ab13ap13Operation, _ := ab13ap13Uint(ab13ap13Context["operation"], 8)
	ab13ap13Direction, _ := ab13ap13Context["direction"].(string)
	ab13ap13View, _ := ab13ap13Context["view"].(string)
	ab13ap13Supported := ab13ap13Primary == 181 && ab13ap13Secondary == 18 && (ab13ap13Direction == "request" || ab13ap13Direction == "positive") && ((ab13ap13Operation == 15 && ab13ap13View == "valve-io") || (ab13ap13Operation == 16 && ab13ap13View == "board-io"))
	if !ab13ap13Supported {
		return json.Marshal(ab13ap13Output)
	}
	ab13ap13Output["state"] = "complete"
	ab13ap13Names := []string{}
	if ab13ap13Operation == 15 && ab13ap13Direction == "request" {
		ab13ap13Names = []string{"output", "valve"}
	}
	if ab13ap13Operation == 15 && ab13ap13Direction == "positive" {
		ab13ap13Names = []string{"temperature1", "temperature2", "input", "pressure", "backup"}
	}
	if ab13ap13Operation == 16 && ab13ap13Direction == "request" {
		ab13ap13Names = []string{"actuator1", "actuator2", "inverse"}
	}
	if ab13ap13Operation == 16 && ab13ap13Direction == "positive" {
		ab13ap13Names = []string{"status", "phase1", "phase2", "adc0", "adc1", "adc2", "adc3", "adc4", "adc5", "board"}
	}
	for _, ab13ap13Name := range ab13ap13Names {
		ab13ap13Result := map[string]any{"state": "incomplete"}
		ab13ap13Fields[ab13ap13Name] = ab13ap13Result
		ab13ap13Item, ab13ap13Exists := ab13ap13Raw[ab13ap13Name]
		if ab13ap13Exists {
			ab13ap13Result["raw"] = ab13ap13Item
		}
		ab13ap13Array, _ := ab13ap13Item.([]any)
		ab13ap13Width := 1
		if ab13ap13Name == "temperature1" || ab13ap13Name == "temperature2" || ab13ap13Name == "inverse" || (len(ab13ap13Name) == 4 && ab13ap13Name[:3] == "adc") {
			ab13ap13Width = 2
		}
		if !ab13ap13Exists || len(ab13ap13Array) != ab13ap13Width {
			ab13ap13Output["state"] = "partial"
			continue
		}
		ab13ap13First, _ := ab13ap13Uint(ab13ap13Array[0], 8)
		ab13ap13Word := ab13ap13First
		if ab13ap13Width == 2 {
			ab13ap13Second, _ := ab13ap13Uint(ab13ap13Array[1], 8)
			ab13ap13Word |= ab13ap13Second << 8
		}
		ab13ap13Result["state"] = "Hypothesis"
		switch ab13ap13Name {
		case "temperature1", "temperature2":
			if ab13ap13Word == 0x8000 {
				ab13ap13Result["label"] = "replacement"
			} else {
				ab13ap13Result["number"] = float64(int16(ab13ap13Word)) / 16
				ab13ap13Result["calculation"] = "Proven"
				ab13ap13Result["unit"] = "Celsius"
			}
		case "pressure":
			ab13ap13Result["number"] = float64(ab13ap13First) / 10
			ab13ap13Result["calculation"] = "Proven"
			ab13ap13Result["unit"] = "bar"
			ab13ap13Result["candidate_domain"] = true
			ab13ap13Result["alternatives"] = []any{map[string]any{"state": "Hypothesis", "label": "candidate-default", "number": 255.0, "candidate_domain": false, "conflict": true}}
		case "phase1", "phase2":
			ab13ap13Result["number"] = ab13ap13First
			ab13ap13Result["unit"] = "percent"
			ab13ap13Result["alternatives"] = []any{map[string]any{"state": "Hypothesis", "label": "single-phase", "approximate": []int{50}}, map[string]any{"state": "Hypothesis", "label": "three-phase", "approximate": []int{0, 33, 66, 100}}}
		case "adc0", "adc1", "adc2", "adc3", "adc4", "adc5":
			ab13ap13Result["number"] = ab13ap13Word
			ab13ap13Result["calculation"] = "Proven"
		case "inverse":
			ab13ap13Result["state"] = "Unknown"
			ab13ap13Result["unmodeled"] = map[string]any{"order": "Unknown", "algorithm": "Unknown"}
		case "valve", "input", "backup", "board":
			ab13ap13Label := ""
			switch ab13ap13Name {
			case "valve":
				if ab13ap13First < 3 {
					ab13ap13Label = []string{"off", "CH", "DHW"}[ab13ap13First]
				}
			case "input":
				if ab13ap13First < 2 {
					ab13ap13Label = []string{"multipurpose", "utility-contact"}[ab13ap13First]
				}
			case "backup":
				if ab13ap13First == 0 {
					ab13ap13Label = "off"
				}
				if ab13ap13First >= 1 && ab13ap13First <= 3 {
					ab13ap13Label = "step"
					ab13ap13Result["number"] = ab13ap13First
				}
				if ab13ap13First == 255 {
					ab13ap13Label = "replacement"
				}
			case "board":
				if ab13ap13First <= 3 || ab13ap13First == 11 || ab13ap13First == 12 {
					ab13ap13Label = "opaque-code-" + strconv.FormatUint(ab13ap13First, 10)
				}
			}
			if ab13ap13Label == "" {
				ab13ap13Result["state"] = "Unknown"
				ab13ap13Result["unmodeled"] = true
			} else {
				ab13ap13Result["label"] = ab13ap13Label
			}
		default:
			ab13ap13Labels := []string{}
			switch ab13ap13Name {
			case "output":
				ab13ap13Labels = []string{"output1", "output2"}
			case "actuator1":
				ab13ap13Labels = []string{"cooling-valve", "passive-mixer1", "passive-mixer2", "cooling-active", "well-pump", "backup-safety", "backup1", "backup2"}
			case "actuator2":
				ab13ap13Labels = []string{"backup3", "external-valve", "external-backup"}
			case "status":
				ab13ap13Labels = []string{"anode1", "anode2", "smart-grid", "well", "safety-limiter", "safety"}
			}
			ab13ap13Flags := map[string]bool{}
			for ab13ap13Bit, ab13ap13Label := range ab13ap13Labels {
				ab13ap13Flags[ab13ap13Label] = ab13ap13First&(1<<ab13ap13Bit) != 0
			}
			ab13ap13Result["flags"] = ab13ap13Flags
			ab13ap13Result["reserved"] = ab13ap13First & (255 << len(ab13ap13Labels))
		}
	}
	return json.Marshal(ab13ap13Output)
}
