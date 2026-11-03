package vaillant

import (
	_ "embed"
	ab34json "encoding/json"
	ab34fmt "fmt"
)

//go:embed data/ab34_b501_legacy_sequence.json
var ab34Embedded string

type Ab34Key struct {
	Ab34Family        string  `json:"ab34_family"`
	Ab34Variant       string  `json:"ab34_variant"`
	Ab34Primary       uint8   `json:"ab34_primary"`
	Ab34Secondary     *uint8  `json:"ab34_secondary"`
	Ab34Operation     *uint8  `json:"ab34_operation"`
	Ab34Direction     string  `json:"ab34_direction"`
	Ab34Field         string  `json:"ab34_field"`
	Ab34Alternative   string  `json:"ab34_alternative"`
	Ab34Group         *uint8  `json:"ab34_group"`
	Ab34InstanceScope string  `json:"ab34_instance_scope"`
	Ab34Register      *uint16 `json:"ab34_register"`
}

func ab34Validate(ab34Data []byte) ([]ab34json.RawMessage, error) {
	var ab34Rows []ab34json.RawMessage
	if ab34Error := ab34json.Unmarshal(ab34Data, &ab34Rows); ab34Error != nil {
		return nil, ab34Error
	}
	if ab34Rows == nil {
		return nil, ab34fmt.Errorf("ab34 array required")
	}
	ab34Seen := make(map[string]bool)
	for _, ab34Row := range ab34Rows {
		var ab34Object map[string]ab34json.RawMessage
		if ab34Error := ab34json.Unmarshal(ab34Row, &ab34Object); ab34Error != nil {
			return nil, ab34Error
		}
		var ab34Key Ab34Key
		if ab34Error := ab34json.Unmarshal(ab34Object["ab34_identity"], &ab34Key); ab34Error != nil {
			return nil, ab34Error
		}
		var ab34Axes map[string]ab34json.RawMessage
		if ab34Error := ab34json.Unmarshal(ab34Object["ab34_identity"], &ab34Axes); ab34Error != nil {
			return nil, ab34Error
		}
		ab34Canonical, ab34Error := ab34json.Marshal(ab34Key)
		if ab34Error != nil {
			return nil, ab34Error
		}
		var ab34Required map[string]ab34json.RawMessage
		if ab34Error = ab34json.Unmarshal(ab34Canonical, &ab34Required); ab34Error != nil {
			return nil, ab34Error
		}
		for ab34Axis := range ab34Required {
			if _, ab34Present := ab34Axes[ab34Axis]; !ab34Present {
				return nil, ab34fmt.Errorf("ab34 missing identity axis %s", ab34Axis)
			}
		}
		if len(ab34Axes) != len(ab34Required) || string(ab34Axes["ab34_primary"]) == "null" || ab34Key.Ab34Family == "" || ab34Key.Ab34Variant == "" || ab34Key.Ab34Direction == "" || ab34Key.Ab34Field == "" || ab34Key.Ab34Alternative == "" || ab34Key.Ab34InstanceScope == "" {
			return nil, ab34fmt.Errorf("ab34 incomplete identity")
		}
		if ab34Seen[string(ab34Canonical)] {
			return nil, ab34fmt.Errorf("ab34 duplicate identity")
		}
		ab34Seen[string(ab34Canonical)] = true
	}
	return ab34Rows, nil
}

func Ab34Entries() (ab34json.RawMessage, error) {
	ab34Data := []byte(ab34Embedded)
	if _, ab34Error := ab34Validate(ab34Data); ab34Error != nil {
		return nil, ab34Error
	}
	return append(ab34json.RawMessage(nil), ab34Data...), nil
}

func ab34Lookup(ab34Data []byte, ab34Key Ab34Key) (ab34json.RawMessage, bool, error) {
	ab34Rows, ab34Error := ab34Validate(ab34Data)
	if ab34Error != nil {
		return nil, false, ab34Error
	}
	ab34Want, ab34Error := ab34json.Marshal(ab34Key)
	if ab34Error != nil {
		return nil, false, ab34Error
	}
	for _, ab34Row := range ab34Rows {
		var ab34Object struct {
			Ab34Identity Ab34Key `json:"ab34_identity"`
		}
		if ab34Error = ab34json.Unmarshal(ab34Row, &ab34Object); ab34Error != nil {
			return nil, false, ab34Error
		}
		ab34Identity, ab34Error := ab34json.Marshal(ab34Object.Ab34Identity)
		if ab34Error != nil {
			return nil, false, ab34Error
		}
		if string(ab34Identity) == string(ab34Want) {
			return append(ab34json.RawMessage(nil), ab34Row...), true, nil
		}
	}
	return nil, false, nil
}

func Ab34Lookup(ab34Key Ab34Key) (ab34json.RawMessage, bool, error) {
	return ab34Lookup([]byte(ab34Embedded), ab34Key)
}
