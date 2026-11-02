package ebus_standard_catalog

import (
	ab54bytes "bytes"
	_ "embed"
	ab54json "encoding/json"
	ab54fmt "fmt"
)

type Ab54Key struct {
	Ab54Family        string  `json:"family"`
	Ab54Variant       string  `json:"variant"`
	Ab54Primary       uint8   `json:"primary"`
	Ab54Secondary     *uint8  `json:"secondary"`
	Ab54Operation     *uint8  `json:"operation"`
	Ab54Direction     string  `json:"direction"`
	Ab54Field         string  `json:"field"`
	Ab54Alternative   string  `json:"alternative"`
	Ab54Group         *uint16 `json:"group"`
	Ab54InstanceScope string  `json:"instance_scope"`
	Ab54Register      *uint16 `json:"register"`
}

//go:embed testdata/ab54_standard_ff_body_states.json
var ab54BodyStatesJSON []byte

func ab54UniqueFields(ab54Decoder *ab54json.Decoder) error {
	ab54Token, ab54Err := ab54Decoder.Token()
	if ab54Err != nil {
		return ab54Err
	}
	ab54Delimiter, ab54Composite := ab54Token.(ab54json.Delim)
	if !ab54Composite {
		return nil
	}
	ab54Seen := make(map[string]bool)
	for ab54Decoder.More() {
		if ab54Delimiter == '{' {
			ab54NameToken, ab54Err := ab54Decoder.Token()
			if ab54Err != nil {
				return ab54Err
			}
			ab54Name, ab54OK := ab54NameToken.(string)
			if !ab54OK || ab54Seen[ab54Name] {
				return ab54fmt.Errorf("ab54: duplicate JSON member")
			}
			ab54Seen[ab54Name] = true
		}
		if ab54Err := ab54UniqueFields(ab54Decoder); ab54Err != nil {
			return ab54Err
		}
	}
	_, ab54Err = ab54Decoder.Token()
	return ab54Err
}

func ab54Validate(ab54Data []byte) (ab54json.RawMessage, map[string]ab54json.RawMessage, error) {
	ab54Data = ab54bytes.TrimSpace(ab54Data)
	if len(ab54Data) == 0 || ab54Data[0] != '[' {
		return nil, nil, ab54fmt.Errorf("ab54: expected JSON array")
	}
	var ab54Rows []ab54json.RawMessage
	if ab54Err := ab54json.Unmarshal(ab54Data, &ab54Rows); ab54Err != nil {
		return nil, nil, ab54Err
	}
	if ab54Err := ab54UniqueFields(ab54json.NewDecoder(ab54bytes.NewReader(ab54Data))); ab54Err != nil {
		return nil, nil, ab54Err
	}
	ab54Index := make(map[string]ab54json.RawMessage, len(ab54Rows))
	for _, ab54Row := range ab54Rows {
		var ab54Envelope map[string]ab54json.RawMessage
		if ab54Err := ab54json.Unmarshal(ab54Row, &ab54Envelope); ab54Err != nil {
			return nil, nil, ab54Err
		}
		var ab54Key Ab54Key
		ab54Decoder := ab54json.NewDecoder(ab54bytes.NewReader(ab54Envelope["key"]))
		ab54Decoder.DisallowUnknownFields()
		if ab54Err := ab54Decoder.Decode(&ab54Key); ab54Err != nil {
			return nil, nil, ab54Err
		}
		var ab54Fields map[string]ab54json.RawMessage
		if ab54Err := ab54json.Unmarshal(ab54Envelope["key"], &ab54Fields); ab54Err != nil {
			return nil, nil, ab54Err
		}
		for ab54Name := range ab54Fields {
			switch ab54Name {
			case "family", "variant", "primary", "secondary", "operation", "direction", "field", "alternative", "group", "instance_scope", "register":
			default:
				return nil, nil, ab54fmt.Errorf("ab54: unknown key axis")
			}
		}
		if len(ab54Fields["primary"]) == 0 || ab54bytes.Equal(ab54Fields["primary"], []byte("null")) ||
			ab54Key.Ab54Family == "" || ab54Key.Ab54Variant == "" || ab54Key.Ab54Direction == "" ||
			ab54Key.Ab54Field == "" || ab54Key.Ab54Alternative == "" || ab54Key.Ab54InstanceScope == "" {
			return nil, nil, ab54fmt.Errorf("ab54: incomplete key; use none for non-applicable string axes")
		}
		ab54Identity, ab54Err := ab54json.Marshal(ab54Key)
		if ab54Err != nil {
			return nil, nil, ab54Err
		}
		if _, ab54Exists := ab54Index[string(ab54Identity)]; ab54Exists {
			return nil, nil, ab54fmt.Errorf("ab54: duplicate key")
		}
		ab54Index[string(ab54Identity)] = append(ab54json.RawMessage(nil), ab54Row...)
	}
	return append(ab54json.RawMessage(nil), ab54Data...), ab54Index, nil
}

func Ab54Entries() (ab54json.RawMessage, error) {
	ab54Entries, _, ab54Err := ab54Validate(ab54BodyStatesJSON)
	return ab54Entries, ab54Err
}

func ab54Lookup(ab54Data []byte, ab54Key Ab54Key) (ab54json.RawMessage, bool, error) {
	_, ab54Index, ab54Err := ab54Validate(ab54Data)
	if ab54Err != nil {
		return nil, false, ab54Err
	}
	ab54Identity, ab54Err := ab54json.Marshal(ab54Key)
	if ab54Err != nil {
		return nil, false, ab54Err
	}
	ab54Row, ab54Found := ab54Index[string(ab54Identity)]
	return append(ab54json.RawMessage(nil), ab54Row...), ab54Found, nil
}

func Ab54Lookup(ab54Key Ab54Key) (ab54json.RawMessage, bool, error) {
	return ab54Lookup(ab54BodyStatesJSON, ab54Key)
}
