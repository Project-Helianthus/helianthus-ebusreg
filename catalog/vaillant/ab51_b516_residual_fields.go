package vaillant

import (
	ab51Bytes "bytes"
	_ "embed"
	ab51JSON "encoding/json"
	ab51Fmt "fmt"
	ab51Reflect "reflect"
)

//go:embed data/ab51_b516_residual_fields.json
var ab51Data string

type Ab51Key struct {
	Ab51Family        string  `json:"family"`
	Ab51Variant       string  `json:"variant"`
	Ab51Primary       uint16  `json:"primary"`
	Ab51Secondary     *uint16 `json:"secondary"`
	Ab51Operation     *uint16 `json:"operation"`
	Ab51Direction     string  `json:"direction"`
	Ab51Field         string  `json:"field"`
	Ab51Alternative   string  `json:"alternative"`
	Ab51Group         *uint16 `json:"group"`
	Ab51InstanceScope string  `json:"instance_scope"`
	Ab51Register      *uint16 `json:"register"`
}

func ab51Validate(ab51Source []byte) ([]ab51JSON.RawMessage, []Ab51Key, error) {
	var ab51Rows []ab51JSON.RawMessage
	if ab51Error := ab51JSON.Unmarshal(ab51Source, &ab51Rows); ab51Error != nil {
		return nil, nil, ab51Error
	}
	if ab51Rows == nil {
		return nil, nil, ab51Fmt.Errorf("ab51: expected array")
	}
	ab51Keys := make([]Ab51Key, 0, len(ab51Rows))
	for _, ab51Row := range ab51Rows {
		var ab51Object map[string]ab51JSON.RawMessage
		if ab51Error := ab51JSON.Unmarshal(ab51Row, &ab51Object); ab51Error != nil {
			return nil, nil, ab51Error
		}
		var ab51Members map[string]ab51JSON.RawMessage
		if ab51Error := ab51JSON.Unmarshal(ab51Object["key"], &ab51Members); ab51Error != nil {
			return nil, nil, ab51Error
		}
		for _, ab51Name := range []string{"family", "variant", "primary", "direction", "field", "alternative", "instance_scope"} {
			if len(ab51Members[ab51Name]) == 0 || ab51Bytes.Equal(ab51Bytes.TrimSpace(ab51Members[ab51Name]), []byte("null")) {
				return nil, nil, ab51Fmt.Errorf("ab51: missing key axis %s", ab51Name)
			}
		}
		var ab51Key Ab51Key
		ab51Decoder := ab51JSON.NewDecoder(ab51Bytes.NewReader(ab51Object["key"]))
		ab51Decoder.DisallowUnknownFields()
		if ab51Error := ab51Decoder.Decode(&ab51Key); ab51Error != nil {
			return nil, nil, ab51Error
		}
		if ab51Key.Ab51Family == "" || ab51Key.Ab51Variant == "" || ab51Key.Ab51Direction == "" || ab51Key.Ab51Field == "" || ab51Key.Ab51Alternative == "" || ab51Key.Ab51InstanceScope == "" {
			return nil, nil, ab51Fmt.Errorf("ab51: empty string axis; use none when not applicable")
		}
		for _, ab51Previous := range ab51Keys {
			if ab51Reflect.DeepEqual(ab51Previous, ab51Key) {
				return nil, nil, ab51Fmt.Errorf("ab51: duplicate key")
			}
		}
		ab51Keys = append(ab51Keys, ab51Key)
	}
	return ab51Rows, ab51Keys, nil
}

func Ab51Entries() (ab51JSON.RawMessage, error) {
	ab51Source := []byte(ab51Data)
	if _, _, ab51Error := ab51Validate(ab51Source); ab51Error != nil {
		return nil, ab51Error
	}
	return ab51JSON.RawMessage(ab51Source), nil
}

func ab51Lookup(ab51Source []byte, ab51Key Ab51Key) (ab51JSON.RawMessage, bool, error) {
	ab51Rows, ab51Keys, ab51Error := ab51Validate(ab51Source)
	if ab51Error != nil {
		return nil, false, ab51Error
	}
	for ab51Index, ab51Candidate := range ab51Keys {
		if ab51Reflect.DeepEqual(ab51Candidate, ab51Key) {
			return append(ab51JSON.RawMessage(nil), ab51Rows[ab51Index]...), true, nil
		}
	}
	return nil, false, nil
}

func Ab51Lookup(ab51Key Ab51Key) (ab51JSON.RawMessage, bool, error) {
	return ab51Lookup([]byte(ab51Data), ab51Key)
}
