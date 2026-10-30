package vaillant

import (
	ab05ao05bytes "bytes"
	ab05ao05json "encoding/json"
	ab05ao05fmt "fmt"
	ab05ao05io "io"
	ab05ao05strconv "strconv"
	ab05ao05strings "strings"
)

func ab05ao05Value(ab05ao05Decoder *ab05ao05json.Decoder) (any, error) {
	ab05ao05Token, ab05ao05Error := ab05ao05Decoder.Token()
	if ab05ao05Error != nil {
		return nil, ab05ao05Error
	}
	switch ab05ao05Token {
	case ab05ao05json.Delim('{'):
		ab05ao05Object := map[string]any{}
		for ab05ao05Decoder.More() {
			ab05ao05KeyToken, ab05ao05Error := ab05ao05Decoder.Token()
			if ab05ao05Error != nil {
				return nil, ab05ao05Error
			}
			ab05ao05Key := ab05ao05KeyToken.(string)
			if _, ab05ao05Exists := ab05ao05Object[ab05ao05Key]; ab05ao05Exists {
				return nil, ab05ao05fmt.Errorf("duplicate key %q", ab05ao05Key)
			}
			ab05ao05Child, ab05ao05Error := ab05ao05Value(ab05ao05Decoder)
			if ab05ao05Error != nil {
				return nil, ab05ao05Error
			}
			ab05ao05Object[ab05ao05Key] = ab05ao05Child
		}
		_, ab05ao05Error = ab05ao05Decoder.Token()
		return ab05ao05Object, ab05ao05Error
	case ab05ao05json.Delim('['):
		ab05ao05Array := []any{}
		for ab05ao05Decoder.More() {
			ab05ao05Child, ab05ao05Error := ab05ao05Value(ab05ao05Decoder)
			if ab05ao05Error != nil {
				return nil, ab05ao05Error
			}
			ab05ao05Array = append(ab05ao05Array, ab05ao05Child)
		}
		_, ab05ao05Error = ab05ao05Decoder.Token()
		return ab05ao05Array, ab05ao05Error
	default:
		return ab05ao05Token, nil
	}
}

func Ab05Ao05Interpret(ab05ao05Input ab05ao05json.RawMessage) (ab05ao05json.RawMessage, error) {
	ab05ao05Decoder := ab05ao05json.NewDecoder(ab05ao05bytes.NewReader(ab05ao05Input))
	ab05ao05Decoder.UseNumber()
	ab05ao05ValueRoot, ab05ao05Error := ab05ao05Value(ab05ao05Decoder)
	if ab05ao05Error != nil {
		return nil, ab05ao05Error
	}
	if _, ab05ao05Error = ab05ao05Decoder.Token(); ab05ao05Error != ab05ao05io.EOF {
		return nil, ab05ao05fmt.Errorf("trailing JSON")
	}
	ab05ao05Root, ab05ao05OK := ab05ao05ValueRoot.(map[string]any)
	if !ab05ao05OK {
		return nil, ab05ao05fmt.Errorf("input must be an object")
	}
	for ab05ao05Key := range ab05ao05Root {
		if ab05ao05Key != "context" && ab05ao05Key != "fields" && ab05ao05Key != "arguments" {
			return nil, ab05ao05fmt.Errorf("unknown input key %q", ab05ao05Key)
		}
	}
	ab05ao05Context, ab05ao05OK := ab05ao05Root["context"].(map[string]any)
	if !ab05ao05OK {
		return nil, ab05ao05fmt.Errorf("context must be an object")
	}
	ab05ao05Raw, ab05ao05OK := ab05ao05Root["fields"].(map[string]any)
	if !ab05ao05OK {
		return nil, ab05ao05fmt.Errorf("fields must be an object")
	}
	for _, ab05ao05Key := range []string{"primary", "secondary", "operation"} {
		if ab05ao05Item, ab05ao05Present := ab05ao05Context[ab05ao05Key]; ab05ao05Present {
			ab05ao05Number, ab05ao05OK := ab05ao05Item.(ab05ao05json.Number)
			if !ab05ao05OK {
				return nil, ab05ao05fmt.Errorf("%s must be an integer", ab05ao05Key)
			}
			if ab05ao05strings.ContainsAny(string(ab05ao05Number), ".eE") {
				return nil, ab05ao05fmt.Errorf("%s must be an integer", ab05ao05Key)
			}
		}
	}
	for _, ab05ao05Key := range []string{"direction", "view"} {
		if ab05ao05Item, ab05ao05Present := ab05ao05Context[ab05ao05Key]; ab05ao05Present {
			if _, ab05ao05OK := ab05ao05Item.(string); !ab05ao05OK {
				return nil, ab05ao05fmt.Errorf("%s must be a string", ab05ao05Key)
			}
		}
	}
	ab05ao05ResultFields := map[string]any{}
	for ab05ao05Key, ab05ao05Item := range ab05ao05Raw {
		ab05ao05Array, ab05ao05OK := ab05ao05Item.([]any)
		if !ab05ao05OK {
			return nil, ab05ao05fmt.Errorf("%s must be a byte array", ab05ao05Key)
		}
		for _, ab05ao05Item := range ab05ao05Array {
			ab05ao05Number, ab05ao05OK := ab05ao05Item.(ab05ao05json.Number)
			if !ab05ao05OK {
				return nil, ab05ao05fmt.Errorf("invalid byte in %s", ab05ao05Key)
			}
			if _, ab05ao05Error := ab05ao05strconv.ParseUint(string(ab05ao05Number), 10, 8); ab05ao05Error != nil {
				return nil, ab05ao05fmt.Errorf("invalid byte in %s", ab05ao05Key)
			}
		}
		ab05ao05ResultFields[ab05ao05Key] = map[string]any{"raw": ab05ao05Array, "state": "Unknown"}
	}
	ab05ao05Output := map[string]any{"context": ab05ao05Context, "raw": ab05ao05Raw, "state": "Unknown", "fields": ab05ao05ResultFields, "physical": "Unknown", "qualified": false}
	if ab05ao05Arguments, ab05ao05Present := ab05ao05Root["arguments"]; ab05ao05Present {
		ab05ao05Object, ab05ao05OK := ab05ao05Arguments.(map[string]any)
		if !ab05ao05OK || len(ab05ao05Object) != 0 {
			return nil, ab05ao05fmt.Errorf("arguments must be an empty object: no argument flags are defined")
		}
		ab05ao05Output["arguments"] = ab05ao05Object
	}
	ab05ao05Selected := ""
	if ab05ao05Context["primary"] == ab05ao05json.Number("181") && ab05ao05Context["secondary"] == ab05ao05json.Number("7") && ab05ao05Context["direction"] == "positive" {
		if ab05ao05Context["operation"] == ab05ao05json.Number("6") && ab05ao05Context["view"] == "capacity" {
			ab05ao05Selected = "capacity"
		}
		if ab05ao05Context["operation"] == ab05ao05json.Number("8") && ab05ao05Context["view"] == "preference" {
			ab05ao05Selected = "preference"
		}
	}
	if ab05ao05Selected != "" {
		ab05ao05Field := map[string]any{"state": "incomplete", "grade": "Hypothesis"}
		ab05ao05ResultFields[ab05ao05Selected] = ab05ao05Field
		ab05ao05Output["state"] = "partial"
		ab05ao05Array, ab05ao05Present := ab05ao05Raw[ab05ao05Selected].([]any)
		if ab05ao05Present {
			ab05ao05Field["raw"] = ab05ao05Array
		}
		ab05ao05Width := 1
		if ab05ao05Selected == "preference" {
			ab05ao05Width = 2
		}
		if ab05ao05Present && len(ab05ao05Array) == ab05ao05Width {
			ab05ao05Field["state"] = "Hypothesis"
			ab05ao05Output["state"] = "Hypothesis"
			ab05ao05First, _ := ab05ao05strconv.Atoi(string(ab05ao05Array[0].(ab05ao05json.Number)))
			if ab05ao05Selected == "capacity" {
				ab05ao05Field["number"] = float64(ab05ao05First) / 10
				ab05ao05Field["calculation"] = "Proven"
				ab05ao05Field["candidate_domain"] = ab05ao05First <= 250
				if ab05ao05First > 250 {
					ab05ao05Field["state"] = "outside-candidate-domain"
				}
			} else {
				ab05ao05Second, _ := ab05ao05strconv.Atoi(string(ab05ao05Array[1].(ab05ao05json.Number)))
				switch ab05ao05First + 256*ab05ao05Second {
				case 0:
					ab05ao05Field["label"] = "appliance-unavailable"
				case 65535:
					ab05ao05Field["label"] = "maximum-efficiency-preference"
				default:
					ab05ao05Field["state"] = "unmapped"
					ab05ao05Output["state"] = "partial"
				}
			}
		}
	}
	return ab05ao05json.Marshal(ab05ao05Output)
}
