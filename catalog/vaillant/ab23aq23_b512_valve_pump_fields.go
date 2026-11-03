package vaillant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

const ab23aq23MaxDepth = 64

func ab23aq23Read(ab23aq23Decoder *json.Decoder) (any, error) {
	return ab23aq23ReadDepth(ab23aq23Decoder, 0)
}

func ab23aq23ReadDepth(ab23aq23Decoder *json.Decoder, ab23aq23Depth int) (any, error) {
	if ab23aq23Depth > ab23aq23MaxDepth {
		return nil, fmt.Errorf("nesting too deep")
	}
	ab23aq23Token, ab23aq23Err := ab23aq23Decoder.Token()
	if ab23aq23Err != nil {
		return nil, ab23aq23Err
	}
	if ab23aq23Delimiter, ab23aq23OK := ab23aq23Token.(json.Delim); ab23aq23OK {
		switch ab23aq23Delimiter {
		case '{':
			ab23aq23Object := map[string]any{}
			for ab23aq23Decoder.More() {
				ab23aq23Key, ab23aq23Err := ab23aq23Decoder.Token()
				if ab23aq23Err != nil {
					return nil, ab23aq23Err
				}
				if _, ab23aq23Exists := ab23aq23Object[ab23aq23Key.(string)]; ab23aq23Exists {
					return nil, fmt.Errorf("duplicate key")
				}
				ab23aq23Value, ab23aq23Err := ab23aq23ReadDepth(ab23aq23Decoder, ab23aq23Depth+1)
				if ab23aq23Err != nil {
					return nil, ab23aq23Err
				}
				ab23aq23Object[ab23aq23Key.(string)] = ab23aq23Value
			}
			_, ab23aq23Err = ab23aq23Decoder.Token()
			return ab23aq23Object, ab23aq23Err
		case '[':
			ab23aq23Array := []any{}
			for ab23aq23Decoder.More() {
				ab23aq23Value, ab23aq23Err := ab23aq23ReadDepth(ab23aq23Decoder, ab23aq23Depth+1)
				if ab23aq23Err != nil {
					return nil, ab23aq23Err
				}
				ab23aq23Array = append(ab23aq23Array, ab23aq23Value)
			}
			_, ab23aq23Err = ab23aq23Decoder.Token()
			return ab23aq23Array, ab23aq23Err
		}
		return nil, fmt.Errorf("unexpected delimiter")
	}
	if ab23aq23Number, ab23aq23OK := ab23aq23Token.(json.Number); ab23aq23OK {
		if _, ab23aq23Err := ab23aq23Number.Float64(); ab23aq23Err != nil {
			return nil, fmt.Errorf("nonfinite number")
		}
	}
	return ab23aq23Token, nil
}

func ab23aq23Byte(ab23aq23Value any) (int, error) {
	ab23aq23Number, ab23aq23OK := ab23aq23Value.(json.Number)
	if !ab23aq23OK {
		return 0, fmt.Errorf("expected byte integer")
	}
	ab23aq23Integer, ab23aq23Err := strconv.Atoi(string(ab23aq23Number))
	if ab23aq23Err != nil || ab23aq23Integer < 0 || ab23aq23Integer > 255 {
		return 0, fmt.Errorf("invalid byte integer")
	}
	return ab23aq23Integer, nil
}

func ab23aq23Bytes(ab23aq23Value any) ([]int, error) {
	ab23aq23Array, ab23aq23OK := ab23aq23Value.([]any)
	if !ab23aq23OK {
		return nil, fmt.Errorf("expected byte array")
	}
	ab23aq23Bytes := make([]int, len(ab23aq23Array))
	for ab23aq23Index, ab23aq23Item := range ab23aq23Array {
		ab23aq23ByteValue, ab23aq23Err := ab23aq23Byte(ab23aq23Item)
		if ab23aq23Err != nil {
			return nil, ab23aq23Err
		}
		ab23aq23Bytes[ab23aq23Index] = ab23aq23ByteValue
	}
	return ab23aq23Bytes, nil
}

func Ab23Aq23Evaluate(ab23aq23Input json.RawMessage) (json.RawMessage, error) {
	ab23aq23Decoder := json.NewDecoder(bytes.NewReader(ab23aq23Input))
	ab23aq23Decoder.UseNumber()
	ab23aq23Value, ab23aq23Err := ab23aq23Read(ab23aq23Decoder)
	if ab23aq23Err != nil {
		return nil, ab23aq23Err
	}
	if _, ab23aq23Err = ab23aq23Decoder.Token(); ab23aq23Err != io.EOF {
		return nil, fmt.Errorf("trailing JSON")
	}
	ab23aq23Root, ab23aq23OK := ab23aq23Value.(map[string]any)
	if !ab23aq23OK {
		return nil, fmt.Errorf("expected object")
	}
	ab23aq23Context, ab23aq23OK := ab23aq23Root["context"].(map[string]any)
	if !ab23aq23OK {
		return nil, fmt.Errorf("expected context object")
	}
	ab23aq23Fields, ab23aq23OK := ab23aq23Root["fields"].(map[string]any)
	if !ab23aq23OK {
		return nil, fmt.Errorf("expected fields object")
	}
	ab23aq23Numbers := map[string]int{}
	for _, ab23aq23Key := range []string{"primary", "secondary", "operation"} {
		if ab23aq23Item, ab23aq23Exists := ab23aq23Context[ab23aq23Key]; ab23aq23Exists {
			ab23aq23Number, ab23aq23Err := ab23aq23Byte(ab23aq23Item)
			if ab23aq23Err != nil {
				return nil, ab23aq23Err
			}
			ab23aq23Numbers[ab23aq23Key] = ab23aq23Number
		}
	}
	for _, ab23aq23Key := range []string{"direction", "view"} {
		if ab23aq23Item, ab23aq23Exists := ab23aq23Context[ab23aq23Key]; ab23aq23Exists {
			if _, ab23aq23OK := ab23aq23Item.(string); !ab23aq23OK {
				return nil, fmt.Errorf("expected context string")
			}
		}
	}
	ab23aq23Raw := map[string][]int{}
	for ab23aq23Key, ab23aq23Item := range ab23aq23Fields {
		ab23aq23BytesValue, ab23aq23Err := ab23aq23Bytes(ab23aq23Item)
		if ab23aq23Err != nil {
			return nil, ab23aq23Err
		}
		ab23aq23Raw[ab23aq23Key] = ab23aq23BytesValue
	}
	ab23aq23Arguments := map[string]any{}
	if ab23aq23Item, ab23aq23Exists := ab23aq23Root["arguments"]; ab23aq23Exists {
		ab23aq23Arguments, ab23aq23OK = ab23aq23Item.(map[string]any)
		if !ab23aq23OK {
			return nil, fmt.Errorf("expected arguments object")
		}
	}
	ab23aq23Operation, ab23aq23HasOperation := ab23aq23Numbers["operation"]
	ab23aq23Direction := ab23aq23Context["direction"]
	ab23aq23Supported := ab23aq23HasOperation && ab23aq23Numbers["primary"] == 181 && ab23aq23Numbers["secondary"] == 18 && ab23aq23Context["view"] == "actuator-fields" && (ab23aq23Direction == "request" || ab23aq23Direction == "positive") && (ab23aq23Operation == 0 || ab23aq23Operation == 2 || ab23aq23Operation == 3 || ab23aq23Operation == 4 || ab23aq23Operation == 5)
	if ab23aq23Direction == "positive" && (ab23aq23Operation == 0 || ab23aq23Operation == 5) {
		ab23aq23Supported = false
	}
	ab23aq23Modes := []int{}
	if ab23aq23Supported {
		if ab23aq23Operation == 5 && ab23aq23Direction == "request" {
			ab23aq23Modes, ab23aq23Err = ab23aq23Bytes(ab23aq23Arguments["external_modes"])
			if ab23aq23Err != nil || len(ab23aq23Modes) != 3 || len(ab23aq23Arguments) != 1 {
				return nil, fmt.Errorf("expected exactly three external modes")
			}
		} else if len(ab23aq23Arguments) != 0 {
			return nil, fmt.Errorf("unexpected arguments")
		}
	}
	ab23aq23Original := ab23aq23Root
	ab23aq23Root = map[string]any{}
	for ab23aq23Key, ab23aq23Item := range ab23aq23Original {
		ab23aq23Root[ab23aq23Key] = ab23aq23Item
	}
	ab23aq23Root["input"] = ab23aq23Original
	ab23aq23Results := map[string]any{}
	for ab23aq23Key, ab23aq23Item := range ab23aq23Fields {
		ab23aq23Results[ab23aq23Key] = map[string]any{"raw": ab23aq23Item, "state": "Unknown"}
	}
	ab23aq23Root["physical"] = "Unknown"
	ab23aq23Root["qualified"] = false
	ab23aq23Root["state"] = "Unknown"
	ab23aq23Root["results"] = ab23aq23Results
	if !ab23aq23Supported {
		return json.Marshal(ab23aq23Root)
	}
	ab23aq23Widths := map[string]int{}
	if ab23aq23Direction == "request" {
		switch ab23aq23Operation {
		case 0, 4:
			ab23aq23Widths["setpoint"] = 1
		case 2:
			ab23aq23Widths["valve"] = 1
		case 3:
			ab23aq23Widths["selector"] = 1
		case 5:
			ab23aq23Widths = map[string]int{"channel0": 1, "channel1": 1, "channel2": 1}
		}
	} else {
		switch ab23aq23Operation {
		case 2:
			ab23aq23Widths = map[string]int{"valve": 1, "contact": 1, "tank": 2, "virtual_error": 2}
		case 3:
			ab23aq23Widths = map[string]int{"status": 1, "value": 2}
		case 4:
			ab23aq23Widths["feedback"] = 1
		}
	}
	ab23aq23Root["state"] = "complete"
	for ab23aq23Key, ab23aq23Width := range ab23aq23Widths {
		ab23aq23Result := map[string]any{"state": "incomplete", "meaning": "Hypothesis"}
		if ab23aq23Item, ab23aq23Exists := ab23aq23Fields[ab23aq23Key]; ab23aq23Exists {
			ab23aq23Result["raw"] = ab23aq23Item
		}
		ab23aq23Results[ab23aq23Key] = ab23aq23Result
		ab23aq23Data := ab23aq23Raw[ab23aq23Key]
		if len(ab23aq23Data) != ab23aq23Width {
			ab23aq23Root["state"] = "partial"
			continue
		}
		ab23aq23Number := ab23aq23Data[0]
		if ab23aq23Width == 2 {
			ab23aq23Number += 256 * ab23aq23Data[1]
		}
		ab23aq23Result["state"] = "complete"
		ab23aq23Result["number"] = ab23aq23Number
		ab23aq23Result["calculation"] = "Proven"
		ab23aq23Label := "unmapped"
		switch ab23aq23Key {
		case "setpoint", "channel0", "channel1", "channel2":
			if ab23aq23Number == 0 {
				ab23aq23Label = "off"
			} else if ab23aq23Number <= 100 {
				ab23aq23Label = "on-or-speed"
			}
			if ab23aq23Operation == 4 && ab23aq23Number == 255 {
				ab23aq23Label = "no-remote-control"
				ab23aq23Result["alternatives"] = []string{"replacement", "no-remote-control"}
			}
			if ab23aq23Operation == 5 {
				if ab23aq23Number == 255 {
					ab23aq23Label = "unchanged"
				}
				ab23aq23Result["prerequisite"] = ab23aq23Modes[int(ab23aq23Key[len(ab23aq23Key)-1]-'0')] == 8
			}
			if ab23aq23Operation == 4 {
				ab23aq23Result["candidate_domain"] = map[string]any{"timeout_seconds": 900, "reset": "every-service"}
			}
		case "valve":
			if ab23aq23Direction == "positive" {
				delete(ab23aq23Result, "number")
				delete(ab23aq23Result, "calculation")
				ab23aq23Result["alternatives"] = []string{"actual", "setpoint-fallback"}
				ab23aq23Label = "unresolved"
			} else {
				switch ab23aq23Number {
				case 0:
					ab23aq23Label = "CH"
				case 50:
					ab23aq23Label = "midpoint"
				case 100:
					ab23aq23Label = "DHW"
				case 254:
					ab23aq23Label = "standby"
				case 255:
					ab23aq23Label = "unchanged"
				}
			}
		case "contact":
			if ab23aq23Number == 0 {
				ab23aq23Label = "open"
			}
			if ab23aq23Number == 1 {
				ab23aq23Label = "closed"
			}
		case "tank":
			ab23aq23Result["number"] = float64(int16(ab23aq23Number)) / 16
			ab23aq23Label = "Celsius"
		case "virtual_error":
			if ab23aq23Number == 65535 {
				ab23aq23Label = "no-error"
			}
		case "selector":
			if ab23aq23Number == 0 {
				ab23aq23Label = "valve-temperature-input"
			}
		case "status":
			switch ab23aq23Number {
			case 0:
				ab23aq23Label = "okay"
			case 1:
				ab23aq23Label = "disconnected"
			case 2:
				ab23aq23Label = "short-circuit"
			}
		case "value":
			ab23aq23Label = "raw"
		case "feedback":
			ab23aq23Label = "rejected"
			if ab23aq23Number == 1 {
				ab23aq23Label = "accepted"
			}
		}
		ab23aq23Result["label"] = ab23aq23Label
	}
	return json.Marshal(ab23aq23Root)
}
