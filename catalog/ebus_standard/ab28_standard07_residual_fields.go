package ebus_standard_catalog

import (
	_ "embed"
	ab28json "encoding/json"
	ab28fmt "fmt"
)

type Ab28Key struct {
	Ab28Family        string  `json:"family"`
	Ab28Variant       string  `json:"variant"`
	Ab28Primary       uint8   `json:"primary"`
	Ab28Secondary     *uint8  `json:"secondary"`
	Ab28Operation     *uint8  `json:"operation"`
	Ab28Direction     string  `json:"direction"`
	Ab28Field         string  `json:"field"`
	Ab28Alternative   string  `json:"alternative"`
	Ab28Group         *uint8  `json:"group"`
	Ab28InstanceScope string  `json:"instance_scope"`
	Ab28Register      *uint16 `json:"register"`
}

//go:embed testdata/ab28_standard07_residual_fields.json
var ab28Embedded string

func ab28Validate(ab28Data []byte) ([]ab28json.RawMessage, error) {
	var ab28Rows []ab28json.RawMessage
	if ab28Err := ab28json.Unmarshal(ab28Data, &ab28Rows); ab28Err != nil {
		return nil, ab28Err
	}
	if ab28Rows == nil {
		return nil, ab28fmt.Errorf("ab28: expected JSON array")
	}
	ab28Seen := make(map[string]bool)
	for _, ab28Row := range ab28Rows {
		var ab28Envelope struct {
			Ab28Key *Ab28Key `json:"key"`
		}
		if ab28Err := ab28json.Unmarshal(ab28Row, &ab28Envelope); ab28Err != nil {
			return nil, ab28Err
		}
		if ab28Envelope.Ab28Key == nil {
			return nil, ab28fmt.Errorf("ab28: missing key")
		}
		var ab28Object map[string]ab28json.RawMessage
		if ab28Err := ab28json.Unmarshal(ab28Row, &ab28Object); ab28Err != nil {
			return nil, ab28Err
		}
		var ab28Axes map[string]ab28json.RawMessage
		if ab28Err := ab28json.Unmarshal(ab28Object["key"], &ab28Axes); ab28Err != nil {
			return nil, ab28Err
		}
		for _, ab28Axis := range []string{"family", "variant", "primary", "secondary", "operation", "direction", "field", "alternative", "group", "instance_scope", "register"} {
			if _, ab28OK := ab28Axes[ab28Axis]; !ab28OK {
				return nil, ab28fmt.Errorf("ab28: missing axis %s", ab28Axis)
			}
		}
		ab28Key := *ab28Envelope.Ab28Key
		if ab28Key.Ab28Family == "" || ab28Key.Ab28Variant == "" || ab28Key.Ab28Direction == "" || ab28Key.Ab28Field == "" || ab28Key.Ab28Alternative == "" || ab28Key.Ab28InstanceScope == "" || string(ab28Axes["primary"]) == "null" {
			return nil, ab28fmt.Errorf("ab28: incomplete identity")
		}
		ab28Identity, ab28Err := ab28json.Marshal(ab28Key)
		if ab28Err != nil {
			return nil, ab28Err
		}
		if ab28Seen[string(ab28Identity)] {
			return nil, ab28fmt.Errorf("ab28: duplicate key")
		}
		ab28Seen[string(ab28Identity)] = true
	}
	return ab28Rows, nil
}

func Ab28Entries() (ab28json.RawMessage, error) {
	if _, ab28Err := ab28Validate([]byte(ab28Embedded)); ab28Err != nil {
		return nil, ab28Err
	}
	return ab28json.RawMessage([]byte(ab28Embedded)), nil
}

func ab28Lookup(ab28Data []byte, ab28Key Ab28Key) (ab28json.RawMessage, bool, error) {
	ab28Rows, ab28Err := ab28Validate(ab28Data)
	if ab28Err != nil {
		return nil, false, ab28Err
	}
	ab28Wanted, ab28Err := ab28json.Marshal(ab28Key)
	if ab28Err != nil {
		return nil, false, ab28Err
	}
	for _, ab28Row := range ab28Rows {
		var ab28Envelope struct {
			Ab28Key Ab28Key `json:"key"`
		}
		if ab28Err := ab28json.Unmarshal(ab28Row, &ab28Envelope); ab28Err != nil {
			return nil, false, ab28Err
		}
		ab28Identity, ab28Err := ab28json.Marshal(ab28Envelope.Ab28Key)
		if ab28Err != nil {
			return nil, false, ab28Err
		}
		if string(ab28Identity) == string(ab28Wanted) {
			return append(ab28json.RawMessage(nil), ab28Row...), true, nil
		}
	}
	return nil, false, nil
}

func Ab28Lookup(ab28Key Ab28Key) (ab28json.RawMessage, bool, error) {
	return ab28Lookup([]byte(ab28Embedded), ab28Key)
}
