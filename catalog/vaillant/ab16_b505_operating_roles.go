package vaillant

import (
	_ "embed"
	ab16json "encoding/json"
	ab16fmt "fmt"
)

//go:embed data/ab16_b505_operating_roles.json
var ab16Data string

type Ab16Key struct {
	Ab16Family        string
	Ab16Variant       string
	Ab16Primary       uint8
	Ab16Secondary     *uint8
	Ab16Operation     *uint8
	Ab16Direction     string
	Ab16Field         string
	Ab16Alternative   string
	Ab16Group         *uint8
	Ab16InstanceScope string
	Ab16Register      *uint16
}

func ab16Validate(ab16Data []byte) (map[string]ab16json.RawMessage, error) {
	var ab16Rows []ab16json.RawMessage
	if ab16Err := ab16json.Unmarshal(ab16Data, &ab16Rows); ab16Err != nil {
		return nil, ab16Err
	}
	if ab16Rows == nil {
		return nil, ab16fmt.Errorf("ab16: expected array")
	}
	ab16Index := make(map[string]ab16json.RawMessage, len(ab16Rows))
	for _, ab16Row := range ab16Rows {
		var ab16Envelope struct {
			Ab16Identity *Ab16Key `json:"identity"`
		}
		if ab16Err := ab16json.Unmarshal(ab16Row, &ab16Envelope); ab16Err != nil {
			return nil, ab16Err
		}
		if ab16Envelope.Ab16Identity == nil {
			return nil, ab16fmt.Errorf("ab16: missing identity")
		}
		ab16Key := ab16Envelope.Ab16Identity
		if ab16Key.Ab16Family == "" || ab16Key.Ab16Variant == "" || ab16Key.Ab16Direction == "" || ab16Key.Ab16Field == "" || ab16Key.Ab16Alternative == "" || ab16Key.Ab16InstanceScope == "" {
			return nil, ab16fmt.Errorf("ab16: incomplete string axes")
		}
		ab16Encoded, ab16Err := ab16json.Marshal(ab16Key)
		if ab16Err != nil {
			return nil, ab16Err
		}
		if _, ab16Found := ab16Index[string(ab16Encoded)]; ab16Found {
			return nil, ab16fmt.Errorf("ab16: duplicate identity")
		}
		ab16Index[string(ab16Encoded)] = append(ab16json.RawMessage(nil), ab16Row...)
	}
	return ab16Index, nil
}

func Ab16Entries() (ab16json.RawMessage, error) {
	if _, ab16Err := ab16Validate([]byte(ab16Data)); ab16Err != nil {
		return nil, ab16Err
	}
	return ab16json.RawMessage(ab16Data), nil
}

func ab16Lookup(ab16Data []byte, ab16Key Ab16Key) (ab16json.RawMessage, bool, error) {
	ab16Index, ab16Err := ab16Validate(ab16Data)
	if ab16Err != nil {
		return nil, false, ab16Err
	}
	ab16Encoded, ab16Err := ab16json.Marshal(ab16Key)
	if ab16Err != nil {
		return nil, false, ab16Err
	}
	ab16Row, ab16Found := ab16Index[string(ab16Encoded)]
	return append(ab16json.RawMessage(nil), ab16Row...), ab16Found, nil
}

func Ab16Lookup(ab16Key Ab16Key) (ab16json.RawMessage, bool, error) {
	return ab16Lookup([]byte(ab16Data), ab16Key)
}
