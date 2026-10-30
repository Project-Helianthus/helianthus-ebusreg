package vaillant

import (
	ab36bytes "bytes"
	_ "embed"
	ab36json "encoding/json"
	ab36fmt "fmt"
)

//go:embed data/ab36_b524_name_annotations.json
var ab36Data string

type Ab36Key struct {
	Ab36Family        string  `json:"family"`
	Ab36Variant       string  `json:"variant"`
	Ab36Primary       uint8   `json:"primary"`
	Ab36Secondary     *uint8  `json:"secondary"`
	Ab36Operation     *uint8  `json:"operation"`
	Ab36Direction     string  `json:"direction"`
	Ab36Field         string  `json:"field"`
	Ab36Alternative   string  `json:"alternative"`
	Ab36Group         *uint8  `json:"group"`
	Ab36InstanceScope string  `json:"instance_scope"`
	Ab36Register      *uint16 `json:"register"`
}

func ab36Validate(ab36Source []byte) ([]ab36json.RawMessage, map[string]ab36json.RawMessage, error) {
	var ab36Rows []ab36json.RawMessage
	if ab36Error := ab36json.Unmarshal(ab36Source, &ab36Rows); ab36Error != nil {
		return nil, nil, ab36Error
	}
	if ab36Rows == nil {
		return nil, nil, ab36fmt.Errorf("ab36: expected an array")
	}
	ab36Index := make(map[string]ab36json.RawMessage, len(ab36Rows))
	for _, ab36Row := range ab36Rows {
		var ab36Object map[string]ab36json.RawMessage
		if ab36Error := ab36json.Unmarshal(ab36Row, &ab36Object); ab36Error != nil {
			return nil, nil, ab36Error
		}
		ab36RawKey, ab36Present := ab36Object["key"]
		if !ab36Present || ab36bytes.Equal(ab36bytes.TrimSpace(ab36RawKey), []byte("null")) {
			return nil, nil, ab36fmt.Errorf("ab36: missing key")
		}
		var ab36Key Ab36Key
		ab36Decoder := ab36json.NewDecoder(ab36bytes.NewReader(ab36RawKey))
		ab36Decoder.DisallowUnknownFields()
		if ab36Error := ab36Decoder.Decode(&ab36Key); ab36Error != nil {
			return nil, nil, ab36Error
		}
		var ab36Fields map[string]ab36json.RawMessage
		if ab36Error := ab36json.Unmarshal(ab36RawKey, &ab36Fields); ab36Error != nil {
			return nil, nil, ab36Error
		}
		ab36RawPrimary, ab36Present := ab36Fields["primary"]
		if !ab36Present || ab36Key.Ab36Family == "" || ab36Key.Ab36Variant == "" || ab36Key.Ab36Direction == "" || ab36Key.Ab36Field == "" || ab36Key.Ab36Alternative == "" || ab36Key.Ab36InstanceScope == "" || ab36bytes.Equal(ab36bytes.TrimSpace(ab36RawPrimary), []byte("null")) {
			return nil, nil, ab36fmt.Errorf("ab36: incomplete key")
		}
		ab36Encoded, ab36Error := ab36json.Marshal(ab36Key)
		if ab36Error != nil {
			return nil, nil, ab36Error
		}
		if _, ab36Present := ab36Index[string(ab36Encoded)]; ab36Present {
			return nil, nil, ab36fmt.Errorf("ab36: duplicate key")
		}
		ab36Index[string(ab36Encoded)] = ab36Row
	}
	return ab36Rows, ab36Index, nil
}

func Ab36Entries() (ab36json.RawMessage, error) {
	if _, _, ab36Error := ab36Validate([]byte(ab36Data)); ab36Error != nil {
		return nil, ab36Error
	}
	return ab36json.RawMessage(ab36Data), nil
}

func Ab36Lookup(ab36Key Ab36Key) (ab36json.RawMessage, bool, error) {
	_, ab36Index, ab36Error := ab36Validate([]byte(ab36Data))
	if ab36Error != nil {
		return nil, false, ab36Error
	}
	ab36Encoded, ab36Error := ab36json.Marshal(ab36Key)
	if ab36Error != nil {
		return nil, false, ab36Error
	}
	ab36Row, ab36Found := ab36Index[string(ab36Encoded)]
	if !ab36Found {
		return nil, false, nil
	}
	return append(ab36json.RawMessage(nil), ab36Row...), true, nil
}
