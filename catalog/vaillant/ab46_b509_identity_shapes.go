package vaillant

import (
	ab46Bytes "bytes"
	_ "embed"
	ab46JSON "encoding/json"
	ab46Fmt "fmt"
	ab46IO "io"
	ab46Reflect "reflect"
)

type Ab46Key struct {
	Ab46Family        string  `json:"family"`
	Ab46Variant       string  `json:"variant"`
	Ab46Primary       *uint16 `json:"primary"`
	Ab46Secondary     *uint16 `json:"secondary"`
	Ab46Operation     *uint16 `json:"operation"`
	Ab46Direction     string  `json:"direction"`
	Ab46Field         string  `json:"field"`
	Ab46Alternative   string  `json:"alternative"`
	Ab46Group         *uint16 `json:"group"`
	Ab46InstanceScope string  `json:"instance_scope"`
	Ab46Register      *uint16 `json:"register"`
}

//go:embed data/ab46_b509_identity_shapes.json
var ab46Data string

func ab46Validate(ab46Source []byte) ([]ab46JSON.RawMessage, []Ab46Key, error) {
	var ab46Rows []ab46JSON.RawMessage
	if ab46Error := ab46JSON.Unmarshal(ab46Source, &ab46Rows); ab46Error != nil {
		return nil, nil, ab46Error
	}
	if ab46Rows == nil {
		return nil, nil, ab46Fmt.Errorf("ab46: expected an array")
	}
	ab46Keys := make([]Ab46Key, 0, len(ab46Rows))
	for ab46Index, ab46Raw := range ab46Rows {
		var ab46Object map[string]ab46JSON.RawMessage
		if ab46Error := ab46JSON.Unmarshal(ab46Raw, &ab46Object); ab46Error != nil || ab46Object == nil {
			return nil, nil, ab46Fmt.Errorf("ab46: row %d must be an object", ab46Index)
		}
		var ab46Key Ab46Key
		ab46Decoder := ab46JSON.NewDecoder(ab46Bytes.NewReader(ab46Object["key"]))
		ab46Decoder.DisallowUnknownFields()
		if ab46Error := ab46Decoder.Decode(&ab46Key); ab46Error != nil {
			return nil, nil, ab46Error
		}
		if ab46Error := ab46Decoder.Decode(new(any)); ab46Error != ab46IO.EOF {
			return nil, nil, ab46Fmt.Errorf("ab46: invalid key")
		}
		if ab46Key.Ab46Family == "" || ab46Key.Ab46Variant == "" || ab46Key.Ab46Primary == nil || ab46Key.Ab46Direction == "" || ab46Key.Ab46Field == "" || ab46Key.Ab46Alternative == "" || ab46Key.Ab46InstanceScope == "" {
			return nil, nil, ab46Fmt.Errorf("ab46: incomplete key at row %d", ab46Index)
		}
		for _, ab46Previous := range ab46Keys {
			if ab46Reflect.DeepEqual(ab46Previous, ab46Key) {
				return nil, nil, ab46Fmt.Errorf("ab46: duplicate key at row %d", ab46Index)
			}
		}
		ab46Keys = append(ab46Keys, ab46Key)
	}
	return ab46Rows, ab46Keys, nil
}

func Ab46Entries() (ab46JSON.RawMessage, error) {
	if _, _, ab46Error := ab46Validate([]byte(ab46Data)); ab46Error != nil {
		return nil, ab46Error
	}
	return ab46JSON.RawMessage([]byte(ab46Data)), nil
}

func ab46Lookup(ab46Source []byte, ab46Key Ab46Key) (ab46JSON.RawMessage, bool, error) {
	ab46Rows, ab46Keys, ab46Error := ab46Validate(ab46Source)
	if ab46Error != nil {
		return nil, false, ab46Error
	}
	for ab46Index, ab46Candidate := range ab46Keys {
		if ab46Reflect.DeepEqual(ab46Candidate, ab46Key) {
			return append(ab46JSON.RawMessage(nil), ab46Rows[ab46Index]...), true, nil
		}
	}
	return nil, false, nil
}

func Ab46Lookup(ab46Key Ab46Key) (ab46JSON.RawMessage, bool, error) {
	return ab46Lookup([]byte(ab46Data), ab46Key)
}
