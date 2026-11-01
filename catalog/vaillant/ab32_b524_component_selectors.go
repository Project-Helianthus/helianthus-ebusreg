package vaillant

import (
	ab32embed "embed"
	ab32json "encoding/json"
	ab32fmt "fmt"
	ab32reflect "reflect"
	ab32strings "strings"
)

//go:embed data/ab32_b524_component_selectors.json
var ab32Files ab32embed.FS

type Ab32Key struct {
	Ab32Family        string  `json:"family"`
	Ab32Variant       string  `json:"variant"`
	Ab32Primary       uint8   `json:"primary"`
	Ab32Secondary     *uint8  `json:"secondary"`
	Ab32Operation     *uint8  `json:"operation"`
	Ab32Direction     string  `json:"direction"`
	Ab32Field         string  `json:"field"`
	Ab32Alternative   string  `json:"alternative"`
	Ab32Group         *uint8  `json:"group"`
	Ab32InstanceScope string  `json:"instance_scope"`
	Ab32Register      *uint16 `json:"register"`
}

func ab32Validate(ab32Data []byte) ([]ab32json.RawMessage, []Ab32Key, error) {
	var ab32Rows []ab32json.RawMessage
	if ab32Error := ab32json.Unmarshal(ab32Data, &ab32Rows); ab32Error != nil {
		return nil, nil, ab32Error
	}
	if ab32Rows == nil {
		return nil, nil, ab32fmt.Errorf("ab32: expected array")
	}
	ab32Keys := make([]Ab32Key, 0, len(ab32Rows))
	ab32IDs := make(map[string]bool)
	for _, ab32Row := range ab32Rows {
		var ab32Envelope struct {
			Ab32ID  string              `json:"id"`
			Ab32Key ab32json.RawMessage `json:"key"`
		}
		if ab32Error := ab32json.Unmarshal(ab32Row, &ab32Envelope); ab32Error != nil {
			return nil, nil, ab32Error
		}
		if !ab32strings.HasPrefix(ab32Envelope.Ab32ID, "ab32-") || ab32IDs[ab32Envelope.Ab32ID] {
			return nil, nil, ab32fmt.Errorf("ab32: invalid or duplicate id")
		}
		ab32IDs[ab32Envelope.Ab32ID] = true
		var ab32Fields map[string]ab32json.RawMessage
		if ab32Error := ab32json.Unmarshal(ab32Envelope.Ab32Key, &ab32Fields); ab32Error != nil {
			return nil, nil, ab32Error
		}
		ab32Axes := []string{"family", "variant", "primary", "secondary", "operation", "direction", "field", "alternative", "group", "instance_scope", "register"}
		if len(ab32Fields) != len(ab32Axes) {
			return nil, nil, ab32fmt.Errorf("ab32: incomplete identity")
		}
		for _, ab32Axis := range ab32Axes {
			if _, ab32Present := ab32Fields[ab32Axis]; !ab32Present {
				return nil, nil, ab32fmt.Errorf("ab32: missing %s", ab32Axis)
			}
		}
		if string(ab32Fields["primary"]) == "null" {
			return nil, nil, ab32fmt.Errorf("ab32: null primary")
		}
		var ab32Key Ab32Key
		if ab32Error := ab32json.Unmarshal(ab32Envelope.Ab32Key, &ab32Key); ab32Error != nil {
			return nil, nil, ab32Error
		}
		for _, ab32Text := range []string{ab32Key.Ab32Family, ab32Key.Ab32Variant, ab32Key.Ab32Direction, ab32Key.Ab32Field, ab32Key.Ab32Alternative, ab32Key.Ab32InstanceScope} {
			if ab32Text == "" {
				return nil, nil, ab32fmt.Errorf("ab32: empty identity string")
			}
		}
		for _, ab32Previous := range ab32Keys {
			if ab32reflect.DeepEqual(ab32Previous, ab32Key) {
				return nil, nil, ab32fmt.Errorf("ab32: duplicate key")
			}
		}
		ab32Keys = append(ab32Keys, ab32Key)
	}
	return ab32Rows, ab32Keys, nil
}

func ab32Read() ([]byte, []ab32json.RawMessage, []Ab32Key, error) {
	ab32Data, ab32Error := ab32Files.ReadFile("data/ab32_b524_component_selectors.json")
	if ab32Error != nil {
		return nil, nil, nil, ab32Error
	}
	ab32Rows, ab32Keys, ab32Error := ab32Validate(ab32Data)
	return ab32Data, ab32Rows, ab32Keys, ab32Error
}

func Ab32Entries() (ab32json.RawMessage, error) {
	ab32Data, _, _, ab32Error := ab32Read()
	if ab32Error != nil {
		return nil, ab32Error
	}
	return append(ab32json.RawMessage(nil), ab32Data...), nil
}

func Ab32Lookup(ab32Key Ab32Key) (ab32json.RawMessage, bool, error) {
	ab32Data, _, _, ab32Error := ab32Read()
	if ab32Error != nil {
		return nil, false, ab32Error
	}
	return ab32LookupData(ab32Data, ab32Key)
}

func ab32LookupData(ab32Data []byte, ab32Key Ab32Key) (ab32json.RawMessage, bool, error) {
	ab32Rows, ab32Keys, ab32Error := ab32Validate(ab32Data)
	if ab32Error != nil {
		return nil, false, ab32Error
	}
	for ab32Index, ab32Candidate := range ab32Keys {
		if ab32reflect.DeepEqual(ab32Candidate, ab32Key) {
			return append(ab32json.RawMessage(nil), ab32Rows[ab32Index]...), true, nil
		}
	}
	return nil, false, nil
}
