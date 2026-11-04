package vaillant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

func ab03aq03Validate(ab03aq03Decoder *json.Decoder) error {
	ab03aq03Token, ab03aq03Err := ab03aq03Decoder.Token()
	if ab03aq03Err != nil {
		return ab03aq03Err
	}
	if ab03aq03Delimiter, ab03aq03OK := ab03aq03Token.(json.Delim); ab03aq03OK {
		switch ab03aq03Delimiter {
		case '{':
			ab03aq03Seen := map[string]bool{}
			for ab03aq03Decoder.More() {
				ab03aq03Key, ab03aq03Err := ab03aq03Decoder.Token()
				if ab03aq03Err != nil {
					return ab03aq03Err
				}
				ab03aq03Name, ab03aq03OK := ab03aq03Key.(string)
				if !ab03aq03OK || ab03aq03Seen[ab03aq03Name] {
					return fmt.Errorf("duplicate or invalid object key")
				}
				ab03aq03Seen[ab03aq03Name] = true
				if ab03aq03Err := ab03aq03Validate(ab03aq03Decoder); ab03aq03Err != nil {
					return ab03aq03Err
				}
			}
		case '[':
			for ab03aq03Decoder.More() {
				if ab03aq03Err := ab03aq03Validate(ab03aq03Decoder); ab03aq03Err != nil {
					return ab03aq03Err
				}
			}
		default:
			return fmt.Errorf("unexpected delimiter")
		}
		_, ab03aq03Err = ab03aq03Decoder.Token()
		return ab03aq03Err
	}
	if ab03aq03Number, ab03aq03OK := ab03aq03Token.(json.Number); ab03aq03OK {
		if _, ab03aq03Err := strconv.ParseFloat(string(ab03aq03Number), 64); ab03aq03Err != nil {
			return fmt.Errorf("nonfinite number")
		}
	}
	return nil
}

func ab03aq03Object(ab03aq03Raw json.RawMessage) (map[string]json.RawMessage, error) {
	if len(bytes.TrimSpace(ab03aq03Raw)) == 0 || bytes.TrimSpace(ab03aq03Raw)[0] != '{' {
		return nil, fmt.Errorf("expected object")
	}
	var ab03aq03Object map[string]json.RawMessage
	ab03aq03Err := json.Unmarshal(ab03aq03Raw, &ab03aq03Object)
	return ab03aq03Object, ab03aq03Err
}

func ab03aq03Byte(ab03aq03Raw json.RawMessage) (int, error) {
	ab03aq03Number, ab03aq03Err := strconv.Atoi(string(ab03aq03Raw))
	if ab03aq03Err != nil || ab03aq03Number < 0 || ab03aq03Number > 255 {
		return 0, fmt.Errorf("expected integer in 0..255")
	}
	return ab03aq03Number, nil
}

func Ab03Aq03Evaluate(ab03aq03InputRaw json.RawMessage) (json.RawMessage, error) {
	ab03aq03Decoder := json.NewDecoder(bytes.NewReader(ab03aq03InputRaw))
	ab03aq03Decoder.UseNumber()
	if ab03aq03Err := ab03aq03Validate(ab03aq03Decoder); ab03aq03Err != nil {
		return nil, ab03aq03Err
	}
	if _, ab03aq03Err := ab03aq03Decoder.Token(); ab03aq03Err != io.EOF {
		return nil, fmt.Errorf("trailing JSON")
	}
	ab03aq03Input, ab03aq03Err := ab03aq03Object(ab03aq03InputRaw)
	if ab03aq03Err != nil {
		return nil, ab03aq03Err
	}
	ab03aq03Context, ab03aq03Err := ab03aq03Object(ab03aq03Input["context"])
	if ab03aq03Err != nil {
		return nil, ab03aq03Err
	}
	ab03aq03Fields, ab03aq03Err := ab03aq03Object(ab03aq03Input["fields"])
	if ab03aq03Err != nil {
		return nil, ab03aq03Err
	}
	ab03aq03Numbers := map[string]int{}
	for _, ab03aq03Key := range []string{"primary", "secondary", "operation"} {
		if ab03aq03Raw, ab03aq03OK := ab03aq03Context[ab03aq03Key]; ab03aq03OK {
			ab03aq03Number, ab03aq03Err := ab03aq03Byte(bytes.TrimSpace(ab03aq03Raw))
			if ab03aq03Err != nil {
				return nil, ab03aq03Err
			}
			ab03aq03Numbers[ab03aq03Key] = ab03aq03Number
		}
	}
	ab03aq03Strings := map[string]string{}
	for _, ab03aq03Key := range []string{"direction", "view", "family", "variant"} {
		if ab03aq03Raw, ab03aq03OK := ab03aq03Context[ab03aq03Key]; ab03aq03OK {
			if bytes.Equal(bytes.TrimSpace(ab03aq03Raw), []byte("null")) {
				return nil, fmt.Errorf("expected identifier string")
			}
			var ab03aq03Value string
			if ab03aq03Err := json.Unmarshal(ab03aq03Raw, &ab03aq03Value); ab03aq03Err != nil {
				return nil, ab03aq03Err
			}
			ab03aq03Strings[ab03aq03Key] = ab03aq03Value
		}
	}
	ab03aq03Bytes := map[string][]int{}
	for ab03aq03Key, ab03aq03Raw := range ab03aq03Fields {
		if len(bytes.TrimSpace(ab03aq03Raw)) == 0 || bytes.TrimSpace(ab03aq03Raw)[0] != '[' {
			return nil, fmt.Errorf("expected byte array")
		}
		var ab03aq03Array []json.RawMessage
		if ab03aq03Err := json.Unmarshal(ab03aq03Raw, &ab03aq03Array); ab03aq03Err != nil {
			return nil, ab03aq03Err
		}
		for _, ab03aq03Element := range ab03aq03Array {
			ab03aq03Number, ab03aq03Err := ab03aq03Byte(ab03aq03Element)
			if ab03aq03Err != nil {
				return nil, ab03aq03Err
			}
			ab03aq03Bytes[ab03aq03Key] = append(ab03aq03Bytes[ab03aq03Key], ab03aq03Number)
		}
	}
	ab03aq03Arguments := map[string]json.RawMessage{}
	if ab03aq03Raw, ab03aq03OK := ab03aq03Input["arguments"]; ab03aq03OK {
		ab03aq03Arguments, ab03aq03Err = ab03aq03Object(ab03aq03Raw)
		if ab03aq03Err != nil {
			return nil, ab03aq03Err
		}
	}
	ab03aq03Display := 0
	for ab03aq03Key, ab03aq03Raw := range ab03aq03Arguments {
		if ab03aq03Key != "display_percent" {
			return nil, fmt.Errorf("unknown argument")
		}
		ab03aq03Display, ab03aq03Err = ab03aq03Byte(bytes.TrimSpace(ab03aq03Raw))
		if ab03aq03Err != nil {
			return nil, ab03aq03Err
		}
	}
	ab03aq03Operation := ab03aq03Numbers["operation"]
	ab03aq03Override := ab03aq03Strings["view"] == "override" && (ab03aq03Operation == 77 || ab03aq03Operation == 97)
	ab03aq03Priority := ab03aq03Strings["view"] == "priority" && ab03aq03Operation == 78
	ab03aq03Supported := ab03aq03Numbers["primary"] == 181 && ab03aq03Numbers["secondary"] == 5 && ab03aq03Strings["direction"] == "request" && (ab03aq03Override || ab03aq03Priority)
	ab03aq03Output := map[string]any{}
	for ab03aq03Key, ab03aq03Raw := range ab03aq03Input {
		ab03aq03Output[ab03aq03Key] = ab03aq03Raw
	}
	ab03aq03Output["physical"] = "Unknown"
	ab03aq03Output["qualified"] = false
	ab03aq03Output["state"] = "Unknown"
	ab03aq03Results := map[string]any{}
	for ab03aq03Key, ab03aq03Raw := range ab03aq03Fields {
		ab03aq03Results[ab03aq03Key] = map[string]any{"raw": ab03aq03Raw, "state": "Unknown"}
	}
	ab03aq03Output["results"] = ab03aq03Results
	if !ab03aq03Supported {
		return json.Marshal(ab03aq03Output)
	}
	ab03aq03Output["state"] = "complete"
	ab03aq03Output["grade"] = "Hypothesis"
	if ab03aq03Override {
		ab03aq03Output["candidate_validity_seconds"] = 900
	}
	ab03aq03Required := []string{}
	if ab03aq03Priority {
		ab03aq03Required = []string{"priority"}
	}
	if ab03aq03Operation == 77 {
		if _, ab03aq03OK := ab03aq03Arguments["display_percent"]; !ab03aq03OK {
			return nil, fmt.Errorf("missing display_percent")
		}
		ab03aq03Results["display_percent"] = map[string]any{"raw": ab03aq03Arguments["display_percent"], "state": "complete", "grade": "Hypothesis", "number": ab03aq03Display, "candidate_domain": ab03aq03Display == 0 || ab03aq03Display == 25 || ab03aq03Display == 50 || ab03aq03Display == 75 || ab03aq03Display == 100}
		ab03aq03Required = []string{"sensor_status", "remote_control"}
	}
	if ab03aq03Operation != 77 && len(ab03aq03Arguments) != 0 {
		return nil, fmt.Errorf("arguments not supported for view")
	}
	for _, ab03aq03Key := range ab03aq03Required {
		ab03aq03Result := map[string]any{"state": "incomplete", "grade": "Hypothesis"}
		if ab03aq03Raw, ab03aq03OK := ab03aq03Fields[ab03aq03Key]; ab03aq03OK {
			ab03aq03Result["raw"] = ab03aq03Raw
		}
		ab03aq03Results[ab03aq03Key] = ab03aq03Result
		if len(ab03aq03Bytes[ab03aq03Key]) != 1 {
			ab03aq03Output["state"] = "partial"
			continue
		}
		ab03aq03Number := ab03aq03Bytes[ab03aq03Key][0]
		ab03aq03Result["state"] = "complete"
		if ab03aq03Key == "priority" {
			ab03aq03Result["number"] = ab03aq03Number
			ab03aq03Result["calculation"] = "Proven"
			continue
		}
		ab03aq03Result["label"] = "unmapped"
		if ab03aq03Key == "sensor_status" {
			if ab03aq03Number == 85 {
				ab03aq03Result["label"] = "short-circuit"
			}
			if ab03aq03Number == 170 {
				ab03aq03Result["label"] = "open"
			}
		}
		if ab03aq03Key == "remote_control" && ab03aq03Number == 1 {
			ab03aq03Result["label"] = "enabled"
		}
	}
	return json.Marshal(ab03aq03Output)
}
