package vaillant

import (
	"embed"
	"encoding/json"
	"fmt"
	"reflect"
)

type Ab26Key struct {
	Ab26Family        string  `json:"family"`
	Ab26Variant       string  `json:"variant"`
	Ab26Primary       uint8   `json:"primary"`
	Ab26Secondary     *uint8  `json:"secondary"`
	Ab26Operation     *uint8  `json:"operation"`
	Ab26Direction     string  `json:"direction"`
	Ab26Field         string  `json:"field"`
	Ab26Alternative   string  `json:"alternative"`
	Ab26Group         *uint8  `json:"group"`
	Ab26InstanceScope string  `json:"instance_scope"`
	Ab26Register      *uint16 `json:"register"`
}

//go:embed data/ab26_b555_additional_forms.json
var ab26Files embed.FS

func ab26Validate(ab26Data []byte) ([]json.RawMessage, []Ab26Key, error) {
	var ab26Rows []json.RawMessage
	if ab26Error := json.Unmarshal(ab26Data, &ab26Rows); ab26Error != nil {
		return nil, nil, ab26Error
	}
	if ab26Rows == nil {
		return nil, nil, fmt.Errorf("ab26: expected JSON array")
	}
	ab26Keys := make([]Ab26Key, len(ab26Rows))
	for ab26Index, ab26Row := range ab26Rows {
		var ab26Envelope struct {
			Ab26Key Ab26Key `json:"key"`
		}
		if ab26Error := json.Unmarshal(ab26Row, &ab26Envelope); ab26Error != nil {
			return nil, nil, ab26Error
		}
		var ab26Object map[string]json.RawMessage
		if ab26Error := json.Unmarshal(ab26Row, &ab26Object); ab26Error != nil {
			return nil, nil, ab26Error
		}
		var ab26Identity map[string]json.RawMessage
		if ab26Error := json.Unmarshal(ab26Object["key"], &ab26Identity); ab26Error != nil {
			return nil, nil, ab26Error
		}
		for _, ab26Axis := range []string{"family", "variant", "primary", "secondary", "operation", "direction", "field", "alternative", "group", "instance_scope", "register"} {
			ab26Value, ab26Present := ab26Identity[ab26Axis]
			if !ab26Present || (string(ab26Value) == "null" && ab26Axis != "secondary" && ab26Axis != "operation" && ab26Axis != "group" && ab26Axis != "register") {
				return nil, nil, fmt.Errorf("ab26: missing identity axis %s", ab26Axis)
			}
		}
		ab26Key := ab26Envelope.Ab26Key
		if ab26Key.Ab26Family == "" || ab26Key.Ab26Variant == "" || ab26Key.Ab26Direction == "" || ab26Key.Ab26Field == "" || ab26Key.Ab26Alternative == "" || ab26Key.Ab26InstanceScope == "" {
			return nil, nil, fmt.Errorf("ab26: incomplete identity at row %d", ab26Index)
		}
		for _, ab26Previous := range ab26Keys[:ab26Index] {
			if reflect.DeepEqual(ab26Previous, ab26Key) {
				return nil, nil, fmt.Errorf("ab26: duplicate identity at row %d", ab26Index)
			}
		}
		ab26Keys[ab26Index] = ab26Key
	}
	return ab26Rows, ab26Keys, nil
}

func Ab26Entries() (json.RawMessage, error) {
	ab26Data, ab26Error := ab26Files.ReadFile("data/ab26_b555_additional_forms.json")
	if ab26Error != nil {
		return nil, ab26Error
	}
	if _, _, ab26Error = ab26Validate(ab26Data); ab26Error != nil {
		return nil, ab26Error
	}
	return append(json.RawMessage(nil), ab26Data...), nil
}

func Ab26Lookup(ab26Key Ab26Key) (json.RawMessage, bool, error) {
	ab26Data, ab26Error := Ab26Entries()
	if ab26Error != nil {
		return nil, false, ab26Error
	}
	ab26Rows, ab26Keys, ab26Error := ab26Validate(ab26Data)
	if ab26Error != nil {
		return nil, false, ab26Error
	}
	for ab26Index, ab26Candidate := range ab26Keys {
		if reflect.DeepEqual(ab26Candidate, ab26Key) {
			return append(json.RawMessage(nil), ab26Rows[ab26Index]...), true, nil
		}
	}
	return nil, false, nil
}
