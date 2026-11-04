package vaillant

import (
	ab21embed "embed"
	ab21json "encoding/json"
	ab21fmt "fmt"
)

//go:embed data/ab21_b519_reply_policy.json
var ab21Files ab21embed.FS

type Ab21Key struct {
	Ab21Family        string  `json:"family"`
	Ab21Variant       string  `json:"variant"`
	Ab21Primary       uint8   `json:"primary"`
	Ab21Secondary     *uint8  `json:"secondary"`
	Ab21Operation     *uint8  `json:"operation"`
	Ab21Direction     string  `json:"direction"`
	Ab21Field         string  `json:"field"`
	Ab21Alternative   string  `json:"alternative"`
	Ab21Group         *uint8  `json:"group"`
	Ab21InstanceScope string  `json:"instance_scope"`
	Ab21Register      *uint16 `json:"register"`
}

func ab21Validate(ab21Data []byte) ([]ab21json.RawMessage, error) {
	var ab21Rows []ab21json.RawMessage
	if ab21Error := ab21json.Unmarshal(ab21Data, &ab21Rows); ab21Error != nil {
		return nil, ab21Error
	}
	if len(ab21Rows) == 0 {
		return nil, ab21fmt.Errorf("ab21: empty catalog")
	}
	ab21Seen := make(map[string]bool)
	for ab21Index, ab21Row := range ab21Rows {
		var ab21Object map[string]ab21json.RawMessage
		if ab21Error := ab21json.Unmarshal(ab21Row, &ab21Object); ab21Error != nil {
			return nil, ab21Error
		}
		var ab21Axes map[string]ab21json.RawMessage
		if ab21Error := ab21json.Unmarshal(ab21Object["key"], &ab21Axes); ab21Error != nil {
			return nil, ab21Error
		}
		if len(ab21Axes) != 11 {
			return nil, ab21fmt.Errorf("ab21: incomplete identity at %d", ab21Index)
		}
		for _, ab21Axis := range []string{"family", "variant", "primary", "secondary", "operation", "direction", "field", "alternative", "group", "instance_scope", "register"} {
			if _, ab21Present := ab21Axes[ab21Axis]; !ab21Present {
				return nil, ab21fmt.Errorf("ab21: missing axis %s", ab21Axis)
			}
		}
		var ab21Key Ab21Key
		if ab21Error := ab21json.Unmarshal(ab21Object["key"], &ab21Key); ab21Error != nil {
			return nil, ab21Error
		}
		if ab21Key.Ab21Family == "" || ab21Key.Ab21Variant == "" || ab21Key.Ab21Direction == "" || ab21Key.Ab21Field == "" || ab21Key.Ab21Alternative == "" || ab21Key.Ab21InstanceScope == "" || string(ab21Axes["primary"]) == "null" {
			return nil, ab21fmt.Errorf("ab21: invalid identity at %d", ab21Index)
		}
		ab21Identity, ab21Error := ab21json.Marshal(ab21Key)
		if ab21Error != nil {
			return nil, ab21Error
		}
		if ab21Seen[string(ab21Identity)] {
			return nil, ab21fmt.Errorf("ab21: duplicate identity at %d", ab21Index)
		}
		ab21Seen[string(ab21Identity)] = true
	}
	return ab21Rows, nil
}

func Ab21Entries() (ab21json.RawMessage, error) {
	ab21Data, ab21Error := ab21Files.ReadFile("data/ab21_b519_reply_policy.json")
	if ab21Error != nil {
		return nil, ab21Error
	}
	if _, ab21Error := ab21Validate(ab21Data); ab21Error != nil {
		return nil, ab21Error
	}
	return append(ab21json.RawMessage(nil), ab21Data...), nil
}

func ab21Lookup(ab21Data []byte, ab21Key Ab21Key) (ab21json.RawMessage, bool, error) {
	ab21Rows, ab21Error := ab21Validate(ab21Data)
	if ab21Error != nil {
		return nil, false, ab21Error
	}
	ab21Wanted, ab21Error := ab21json.Marshal(ab21Key)
	if ab21Error != nil {
		return nil, false, ab21Error
	}
	for _, ab21Row := range ab21Rows {
		var ab21Entry struct {
			Ab21Key Ab21Key `json:"key"`
		}
		if ab21Error := ab21json.Unmarshal(ab21Row, &ab21Entry); ab21Error != nil {
			return nil, false, ab21Error
		}
		ab21Identity, ab21Error := ab21json.Marshal(ab21Entry.Ab21Key)
		if ab21Error != nil {
			return nil, false, ab21Error
		}
		if string(ab21Wanted) == string(ab21Identity) {
			return append(ab21json.RawMessage(nil), ab21Row...), true, nil
		}
	}
	return nil, false, nil
}

func Ab21Lookup(ab21Key Ab21Key) (ab21json.RawMessage, bool, error) {
	ab21Data, ab21Error := ab21Files.ReadFile("data/ab21_b519_reply_policy.json")
	if ab21Error != nil {
		return nil, false, ab21Error
	}
	return ab21Lookup(ab21Data, ab21Key)
}
