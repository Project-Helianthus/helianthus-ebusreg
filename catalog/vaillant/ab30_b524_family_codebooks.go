package vaillant

import (
	_ "embed"
	ab30json "encoding/json"
	ab30fmt "fmt"
	ab30reflect "reflect"
)

//go:embed data/ab30_b524_family_codebooks.json
var ab30Embedded string

type Ab30Key struct {
	Ab30Family        string  `json:"family"`
	Ab30Variant       string  `json:"variant"`
	Ab30Primary       uint8   `json:"primary"`
	Ab30Secondary     *uint8  `json:"secondary"`
	Ab30Operation     *uint8  `json:"operation"`
	Ab30Direction     string  `json:"direction"`
	Ab30Field         string  `json:"field"`
	Ab30Alternative   string  `json:"alternative"`
	Ab30Group         *uint8  `json:"group"`
	Ab30InstanceScope string  `json:"instance_scope"`
	Ab30Register      *uint16 `json:"register"`
}

func ab30Validate(ab30Data []byte) ([]ab30json.RawMessage, []Ab30Key, error) {
	var ab30Rows []ab30json.RawMessage
	if ab30Error := ab30json.Unmarshal(ab30Data, &ab30Rows); ab30Error != nil {
		return nil, nil, ab30Error
	}
	if ab30Rows == nil {
		return nil, nil, ab30fmt.Errorf("ab30 entries must be an array")
	}
	ab30Keys := make([]Ab30Key, 0, len(ab30Rows))
	for _, ab30Row := range ab30Rows {
		var ab30Entry struct {
			Ab30Key Ab30Key `json:"key"`
		}
		if ab30Error := ab30json.Unmarshal(ab30Row, &ab30Entry); ab30Error != nil {
			return nil, nil, ab30Error
		}
		var ab30Object map[string]ab30json.RawMessage
		if ab30Error := ab30json.Unmarshal(ab30Row, &ab30Object); ab30Error != nil {
			return nil, nil, ab30Error
		}
		var ab30Identity map[string]ab30json.RawMessage
		if ab30Error := ab30json.Unmarshal(ab30Object["key"], &ab30Identity); ab30Error != nil {
			return nil, nil, ab30Error
		}
		for _, ab30Axis := range []string{"family", "variant", "primary", "secondary", "operation", "direction", "field", "alternative", "group", "instance_scope", "register"} {
			ab30Value, ab30Present := ab30Identity[ab30Axis]
			if !ab30Present || (string(ab30Value) == "null" && ab30Axis != "secondary" && ab30Axis != "operation" && ab30Axis != "group" && ab30Axis != "register") {
				return nil, nil, ab30fmt.Errorf("ab30 missing identity axis %s", ab30Axis)
			}
		}
		if ab30Entry.Ab30Key.Ab30Family == "" || ab30Entry.Ab30Key.Ab30Variant == "" || ab30Entry.Ab30Key.Ab30Direction == "" || ab30Entry.Ab30Key.Ab30Field == "" || ab30Entry.Ab30Key.Ab30Alternative == "" || ab30Entry.Ab30Key.Ab30InstanceScope == "" {
			return nil, nil, ab30fmt.Errorf("ab30 string axes require explicit tokens")
		}
		for _, ab30Prior := range ab30Keys {
			if ab30reflect.DeepEqual(ab30Prior, ab30Entry.Ab30Key) {
				return nil, nil, ab30fmt.Errorf("ab30 duplicate identity")
			}
		}
		ab30Keys = append(ab30Keys, ab30Entry.Ab30Key)
	}
	return ab30Rows, ab30Keys, nil
}

func Ab30Entries() (ab30json.RawMessage, error) {
	if _, _, ab30Error := ab30Validate([]byte(ab30Embedded)); ab30Error != nil {
		return nil, ab30Error
	}
	return ab30json.RawMessage([]byte(ab30Embedded)), nil
}

func ab30Lookup(ab30Data []byte, ab30Key Ab30Key) (ab30json.RawMessage, bool, error) {
	ab30Rows, ab30Keys, ab30Error := ab30Validate(ab30Data)
	if ab30Error != nil {
		return nil, false, ab30Error
	}
	for ab30Index, ab30Candidate := range ab30Keys {
		if ab30reflect.DeepEqual(ab30Candidate, ab30Key) {
			return append(ab30json.RawMessage(nil), ab30Rows[ab30Index]...), true, nil
		}
	}
	return nil, false, nil
}

func Ab30Lookup(ab30Key Ab30Key) (ab30json.RawMessage, bool, error) {
	return ab30Lookup([]byte(ab30Embedded), ab30Key)
}
