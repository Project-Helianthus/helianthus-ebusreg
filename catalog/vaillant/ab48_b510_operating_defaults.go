package vaillant

import (
	ab48bytes "bytes"
	_ "embed"
	ab48json "encoding/json"
	ab48fmt "fmt"
	ab48reflect "reflect"
)

//go:embed data/ab48_b510_operating_defaults.json
var ab48Data string

type Ab48Key struct {
	Ab48Family        string  `json:"family"`
	Ab48Variant       string  `json:"variant"`
	Ab48Primary       uint8   `json:"primary"`
	Ab48Secondary     *uint8  `json:"secondary"`
	Ab48Operation     *uint8  `json:"operation"`
	Ab48Direction     string  `json:"direction"`
	Ab48Field         string  `json:"field"`
	Ab48Alternative   string  `json:"alternative"`
	Ab48Group         *uint8  `json:"group"`
	Ab48InstanceScope string  `json:"instance_scope"`
	Ab48Register      *uint16 `json:"register"`
}

func ab48Validate(ab48Source []byte) ([]ab48json.RawMessage, []Ab48Key, error) {
	var ab48Rows []ab48json.RawMessage
	if ab48Err := ab48json.Unmarshal(ab48Source, &ab48Rows); ab48Err != nil {
		return nil, nil, ab48Err
	}
	if ab48Rows == nil {
		return nil, nil, ab48fmt.Errorf("ab48 catalog: expected array")
	}
	ab48Keys := make([]Ab48Key, len(ab48Rows))
	for ab48Index, ab48Raw := range ab48Rows {
		var ab48Object map[string]ab48json.RawMessage
		if ab48Err := ab48json.Unmarshal(ab48Raw, &ab48Object); ab48Err != nil {
			return nil, nil, ab48Err
		}
		var ab48Members map[string]ab48json.RawMessage
		if ab48Err := ab48json.Unmarshal(ab48Object["key"], &ab48Members); ab48Err != nil {
			return nil, nil, ab48Err
		}
		for _, ab48Name := range []string{"family", "variant", "primary", "direction", "field", "alternative", "instance_scope"} {
			if ab48Value, ab48OK := ab48Members[ab48Name]; !ab48OK || ab48bytes.Equal(ab48bytes.TrimSpace(ab48Value), []byte("null")) {
				return nil, nil, ab48fmt.Errorf("ab48 catalog: missing key axis %s", ab48Name)
			}
		}
		ab48Decoder := ab48json.NewDecoder(ab48bytes.NewReader(ab48Object["key"]))
		ab48Decoder.DisallowUnknownFields()
		if ab48Err := ab48Decoder.Decode(&ab48Keys[ab48Index]); ab48Err != nil {
			return nil, nil, ab48Err
		}
		ab48Key := ab48Keys[ab48Index]
		if ab48Key.Ab48Family == "" || ab48Key.Ab48Variant == "" || ab48Key.Ab48Direction == "" || ab48Key.Ab48Field == "" || ab48Key.Ab48Alternative == "" || ab48Key.Ab48InstanceScope == "" {
			return nil, nil, ab48fmt.Errorf("ab48 catalog: empty key axis; use none for non-applicable strings")
		}
		for ab48Previous := 0; ab48Previous < ab48Index; ab48Previous++ {
			if ab48reflect.DeepEqual(ab48Keys[ab48Previous], ab48Key) {
				return nil, nil, ab48fmt.Errorf("ab48 catalog: duplicate key")
			}
		}
	}
	return ab48Rows, ab48Keys, nil
}

func Ab48Entries() (ab48json.RawMessage, error) {
	ab48Source := []byte(ab48Data)
	if _, _, ab48Err := ab48Validate(ab48Source); ab48Err != nil {
		return nil, ab48Err
	}
	return ab48Source, nil
}

func ab48LookupSource(ab48Source []byte, ab48Key Ab48Key) (ab48json.RawMessage, bool, error) {
	ab48Rows, ab48Keys, ab48Err := ab48Validate(ab48Source)
	if ab48Err != nil {
		return nil, false, ab48Err
	}
	for ab48Index, ab48Candidate := range ab48Keys {
		if ab48reflect.DeepEqual(ab48Candidate, ab48Key) {
			return ab48bytes.Clone(ab48Rows[ab48Index]), true, nil
		}
	}
	return nil, false, nil
}

func Ab48Lookup(ab48Key Ab48Key) (ab48json.RawMessage, bool, error) {
	return ab48LookupSource([]byte(ab48Data), ab48Key)
}
