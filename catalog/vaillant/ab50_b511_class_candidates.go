package vaillant

import (
	ab50Bytes "bytes"
	_ "embed"
	ab50JSON "encoding/json"
	ab50Fmt "fmt"
)

//go:embed data/ab50_b511_class_candidates.json
var ab50Data string

type Ab50Key struct {
	Ab50Family        string  `json:"family"`
	Ab50Variant       string  `json:"variant"`
	Ab50Primary       uint16  `json:"primary"`
	Ab50Secondary     *uint16 `json:"secondary"`
	Ab50Operation     *uint16 `json:"operation"`
	Ab50Direction     string  `json:"direction"`
	Ab50Field         string  `json:"field"`
	Ab50Alternative   string  `json:"alternative"`
	Ab50Group         *uint16 `json:"group"`
	Ab50InstanceScope string  `json:"instance_scope"`
	Ab50Register      *uint16 `json:"register"`
}

func ab50Validate(ab50Source []byte) ([]ab50JSON.RawMessage, []Ab50Key, error) {
	var ab50Rows []ab50JSON.RawMessage
	if ab50Error := ab50JSON.Unmarshal(ab50Source, &ab50Rows); ab50Error != nil {
		return nil, nil, ab50Error
	}
	if ab50Rows == nil {
		return nil, nil, ab50Fmt.Errorf("ab50: expected array")
	}
	ab50Keys := make([]Ab50Key, len(ab50Rows))
	ab50Seen := make(map[string]bool)
	for ab50Index, ab50Row := range ab50Rows {
		var ab50Object map[string]ab50JSON.RawMessage
		if ab50Error := ab50JSON.Unmarshal(ab50Row, &ab50Object); ab50Error != nil {
			return nil, nil, ab50Error
		}
		ab50KeyJSON, ab50OK := ab50Object["key"]
		if !ab50OK {
			return nil, nil, ab50Fmt.Errorf("ab50: missing key")
		}
		var ab50Axes map[string]ab50JSON.RawMessage
		if ab50Error := ab50JSON.Unmarshal(ab50KeyJSON, &ab50Axes); ab50Error != nil {
			return nil, nil, ab50Error
		}
		if ab50Primary, ab50OK := ab50Axes["primary"]; !ab50OK || ab50Bytes.Equal(ab50Bytes.TrimSpace(ab50Primary), []byte("null")) {
			return nil, nil, ab50Fmt.Errorf("ab50: missing primary")
		}
		ab50Decoder := ab50JSON.NewDecoder(ab50Bytes.NewReader(ab50KeyJSON))
		ab50Decoder.DisallowUnknownFields()
		if ab50Error := ab50Decoder.Decode(&ab50Keys[ab50Index]); ab50Error != nil {
			return nil, nil, ab50Error
		}
		ab50Key := ab50Keys[ab50Index]
		if ab50Key.Ab50Family == "" || ab50Key.Ab50Variant == "" || ab50Key.Ab50Direction == "" || ab50Key.Ab50Field == "" || ab50Key.Ab50Alternative == "" || ab50Key.Ab50InstanceScope == "" {
			return nil, nil, ab50Fmt.Errorf("ab50: incomplete key")
		}
		ab50Canonical, ab50Error := ab50JSON.Marshal(ab50Key)
		if ab50Error != nil {
			return nil, nil, ab50Error
		}
		if ab50Seen[string(ab50Canonical)] {
			return nil, nil, ab50Fmt.Errorf("ab50: duplicate key")
		}
		ab50Seen[string(ab50Canonical)] = true
	}
	return ab50Rows, ab50Keys, nil
}

func ab50Entries(ab50Source []byte) (ab50JSON.RawMessage, error) {
	if _, _, ab50Error := ab50Validate(ab50Source); ab50Error != nil {
		return nil, ab50Error
	}
	return append(ab50JSON.RawMessage(nil), ab50Source...), nil
}

func Ab50Entries() (ab50JSON.RawMessage, error) { return ab50Entries([]byte(ab50Data)) }

func ab50Lookup(ab50Source []byte, ab50Key Ab50Key) (ab50JSON.RawMessage, bool, error) {
	ab50Rows, ab50Keys, ab50Error := ab50Validate(ab50Source)
	if ab50Error != nil {
		return nil, false, ab50Error
	}
	ab50Wanted, ab50Error := ab50JSON.Marshal(ab50Key)
	if ab50Error != nil {
		return nil, false, ab50Error
	}
	for ab50Index, ab50Candidate := range ab50Keys {
		ab50Canonical, ab50Error := ab50JSON.Marshal(ab50Candidate)
		if ab50Error != nil {
			return nil, false, ab50Error
		}
		if ab50Bytes.Equal(ab50Wanted, ab50Canonical) {
			return append(ab50JSON.RawMessage(nil), ab50Rows[ab50Index]...), true, nil
		}
	}
	return nil, false, nil
}

func Ab50Lookup(ab50Key Ab50Key) (ab50JSON.RawMessage, bool, error) {
	return ab50Lookup([]byte(ab50Data), ab50Key)
}
