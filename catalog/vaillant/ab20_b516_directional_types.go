package vaillant

import (
	_ "embed"
	ab20json "encoding/json"
	ab20fmt "fmt"
	ab20reflect "reflect"
	ab20strings "strings"
)

//go:embed data/ab20_b516_directional_types.json
var ab20Data string

type Ab20Key struct {
	Ab20Family        string  `json:"family"`
	Ab20Variant       string  `json:"variant"`
	Ab20Primary       uint8   `json:"primary"`
	Ab20Secondary     *uint8  `json:"secondary,omitempty"`
	Ab20Operation     *uint8  `json:"operation,omitempty"`
	Ab20Direction     string  `json:"direction"`
	Ab20Field         string  `json:"field"`
	Ab20Alternative   string  `json:"alternative"`
	Ab20Group         *uint8  `json:"group,omitempty"`
	Ab20InstanceScope string  `json:"instance_scope"`
	Ab20Register      *uint16 `json:"register,omitempty"`
}

func ab20Validate(ab20Data []byte) ([]ab20json.RawMessage, error) {
	var ab20Rows []ab20json.RawMessage
	if ab20Err := ab20json.Unmarshal(ab20Data, &ab20Rows); ab20Err != nil {
		return nil, ab20fmt.Errorf("ab20 catalog: %w", ab20Err)
	}
	if len(ab20Rows) == 0 {
		return nil, ab20fmt.Errorf("ab20 catalog: empty array")
	}
	ab20Keys := make([]Ab20Key, 0, len(ab20Rows))
	ab20IDs := make(map[string]bool)
	for _, ab20Row := range ab20Rows {
		var ab20Entry struct {
			Ab20ID  string  `json:"id"`
			Ab20Key Ab20Key `json:"key"`
		}
		if ab20Err := ab20json.Unmarshal(ab20Row, &ab20Entry); ab20Err != nil {
			return nil, ab20fmt.Errorf("ab20 catalog: %w", ab20Err)
		}
		var ab20Object map[string]ab20json.RawMessage
		if ab20Err := ab20json.Unmarshal(ab20Row, &ab20Object); ab20Err != nil {
			return nil, ab20Err
		}
		var ab20Identity map[string]ab20json.RawMessage
		if ab20Err := ab20json.Unmarshal(ab20Object["key"], &ab20Identity); ab20Err != nil {
			return nil, ab20Err
		}
		if ab20Primary, ab20Present := ab20Identity["primary"]; !ab20Present || string(ab20Primary) == "null" {
			return nil, ab20fmt.Errorf("ab20 catalog: missing primary")
		}
		ab20Key := ab20Entry.Ab20Key
		if !ab20strings.HasPrefix(ab20Entry.Ab20ID, "ab20-") || ab20IDs[ab20Entry.Ab20ID] || ab20Key.Ab20Family == "" || ab20Key.Ab20Variant == "" || ab20Key.Ab20Direction == "" || ab20Key.Ab20Field == "" || ab20Key.Ab20Alternative == "" || ab20Key.Ab20InstanceScope == "" {
			return nil, ab20fmt.Errorf("ab20 catalog: incomplete identity or duplicate id")
		}
		for _, ab20Previous := range ab20Keys {
			if ab20reflect.DeepEqual(ab20Previous, ab20Key) {
				return nil, ab20fmt.Errorf("ab20 catalog: duplicate key")
			}
		}
		ab20IDs[ab20Entry.Ab20ID] = true
		ab20Keys = append(ab20Keys, ab20Key)
	}
	return ab20Rows, nil
}

func Ab20Entries() (ab20json.RawMessage, error) {
	if _, ab20Err := ab20Validate([]byte(ab20Data)); ab20Err != nil {
		return nil, ab20Err
	}
	return ab20json.RawMessage([]byte(ab20Data)), nil
}

func Ab20Lookup(ab20Key Ab20Key) (ab20json.RawMessage, bool, error) {
	return ab20Lookup([]byte(ab20Data), ab20Key)
}

func ab20Lookup(ab20Data []byte, ab20Key Ab20Key) (ab20json.RawMessage, bool, error) {
	ab20Rows, ab20Err := ab20Validate(ab20Data)
	if ab20Err != nil {
		return nil, false, ab20Err
	}
	for _, ab20Row := range ab20Rows {
		var ab20Entry struct {
			Ab20Key Ab20Key `json:"key"`
		}
		if ab20Err := ab20json.Unmarshal(ab20Row, &ab20Entry); ab20Err != nil {
			return nil, false, ab20Err
		}
		if ab20reflect.DeepEqual(ab20Key, ab20Entry.Ab20Key) {
			return append(ab20json.RawMessage(nil), ab20Row...), true, nil
		}
	}
	return nil, false, nil
}
