package vaillant

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:embed data/ab40_b509_symbolic_annotations.json
var ab40Data []byte

type Ab40Key struct {
	Ab40Family        string  `json:"family"`
	Ab40Variant       string  `json:"variant"`
	Ab40Primary       uint8   `json:"primary"`
	Ab40Secondary     *uint8  `json:"secondary,omitempty"`
	Ab40Operation     *uint8  `json:"operation,omitempty"`
	Ab40Direction     string  `json:"direction"`
	Ab40Field         string  `json:"field"`
	Ab40Alternative   string  `json:"alternative"`
	Ab40Group         *uint8  `json:"group,omitempty"`
	Ab40InstanceScope string  `json:"instance_scope"`
	Ab40Register      *uint16 `json:"register,omitempty"`
	Ab40Class         string  `json:"class"`
	Ab40Selector      *string `json:"selector,omitempty"`
}

func ab40Identity(ab40Key Ab40Key) string {
	ab40Encoded, _ := json.Marshal(ab40Key)
	return string(ab40Encoded)
}

func ab40Validate(ab40Source []byte) ([]json.RawMessage, map[string]json.RawMessage, error) {
	var ab40Rows []json.RawMessage
	if ab40Error := json.Unmarshal(ab40Source, &ab40Rows); ab40Error != nil {
		return nil, nil, ab40Error
	}
	if ab40Rows == nil {
		return nil, nil, fmt.Errorf("ab40: expected array")
	}
	ab40Index := make(map[string]json.RawMessage, len(ab40Rows))
	for _, ab40Raw := range ab40Rows {
		var ab40Key Ab40Key
		if ab40Error := json.Unmarshal(ab40Raw, &ab40Key); ab40Error != nil {
			return nil, nil, ab40Error
		}
		var ab40Object map[string]json.RawMessage
		if ab40Error := json.Unmarshal(ab40Raw, &ab40Object); ab40Error != nil {
			return nil, nil, ab40Error
		}
		for _, ab40Axis := range []string{"family", "variant", "primary", "direction", "field", "alternative", "instance_scope", "class"} {
			if ab40Value, ab40Present := ab40Object[ab40Axis]; !ab40Present || string(ab40Value) == "null" {
				return nil, nil, fmt.Errorf("ab40: missing axis %s", ab40Axis)
			}
		}
		for _, ab40Axis := range []string{"secondary", "operation", "group", "register", "selector"} {
			if ab40Value, ab40Present := ab40Object[ab40Axis]; ab40Present && string(ab40Value) == "null" {
				return nil, nil, fmt.Errorf("ab40: null optional axis %s", ab40Axis)
			}
		}
		if ab40Key.Ab40Family == "" || ab40Key.Ab40Variant == "" || ab40Key.Ab40Direction == "" || ab40Key.Ab40Field == "" || ab40Key.Ab40Alternative == "" || ab40Key.Ab40InstanceScope == "" || ab40Key.Ab40Class == "" {
			return nil, nil, fmt.Errorf("ab40: empty string axis; use none when inapplicable")
		}
		ab40ID := ab40Identity(ab40Key)
		if _, ab40Present := ab40Index[ab40ID]; ab40Present {
			return nil, nil, fmt.Errorf("ab40: duplicate key")
		}
		ab40Index[ab40ID] = append(json.RawMessage(nil), ab40Raw...)
	}
	return ab40Rows, ab40Index, nil
}

func Ab40Entries() (json.RawMessage, error) {
	if _, _, ab40Error := ab40Validate(ab40Data); ab40Error != nil {
		return nil, ab40Error
	}
	return append(json.RawMessage(nil), ab40Data...), nil
}

func Ab40Lookup(ab40Key Ab40Key) (json.RawMessage, bool, error) {
	_, ab40Index, ab40Error := ab40Validate(ab40Data)
	if ab40Error != nil {
		return nil, false, ab40Error
	}
	ab40Raw, ab40Present := ab40Index[ab40Identity(ab40Key)]
	return append(json.RawMessage(nil), ab40Raw...), ab40Present, nil
}
