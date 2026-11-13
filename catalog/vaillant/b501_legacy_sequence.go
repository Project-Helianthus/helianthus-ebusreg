package vaillant

import (
	_ "embed"
	b501LegacySequenceJSON "encoding/json"
	b501LegacySequenceFmt "fmt"
)

//go:embed data/ab34_b501_legacy_sequence.json
var b501LegacySequenceEmbedded string

type B501LegacySequenceKey struct {
	B501LegacySequenceFamily        string  `json:"ab34_family"`
	B501LegacySequenceVariant       string  `json:"ab34_variant"`
	B501LegacySequencePrimary       uint8   `json:"ab34_primary"`
	B501LegacySequenceSecondary     *uint8  `json:"ab34_secondary"`
	B501LegacySequenceOperation     *uint8  `json:"ab34_operation"`
	B501LegacySequenceDirection     string  `json:"ab34_direction"`
	B501LegacySequenceField         string  `json:"ab34_field"`
	B501LegacySequenceAlternative   string  `json:"ab34_alternative"`
	B501LegacySequenceGroup         *uint8  `json:"ab34_group"`
	B501LegacySequenceInstanceScope string  `json:"ab34_instance_scope"`
	B501LegacySequenceRegister      *uint16 `json:"ab34_register"`
}

func b501LegacySequenceValidate(b501LegacySequenceData []byte) ([]b501LegacySequenceJSON.RawMessage, error) {
	var b501LegacySequenceRows []b501LegacySequenceJSON.RawMessage
	if b501LegacySequenceError := b501LegacySequenceJSON.Unmarshal(b501LegacySequenceData, &b501LegacySequenceRows); b501LegacySequenceError != nil {
		return nil, b501LegacySequenceError
	}
	if b501LegacySequenceRows == nil {
		return nil, b501LegacySequenceFmt.Errorf("ab34 array required")
	}
	b501LegacySequenceSeen := make(map[string]bool)
	for _, b501LegacySequenceRow := range b501LegacySequenceRows {
		var b501LegacySequenceObject map[string]b501LegacySequenceJSON.RawMessage
		if b501LegacySequenceError := b501LegacySequenceJSON.Unmarshal(b501LegacySequenceRow, &b501LegacySequenceObject); b501LegacySequenceError != nil {
			return nil, b501LegacySequenceError
		}
		var b501LegacySequenceKey B501LegacySequenceKey
		if b501LegacySequenceError := b501LegacySequenceJSON.Unmarshal(b501LegacySequenceObject["ab34_identity"], &b501LegacySequenceKey); b501LegacySequenceError != nil {
			return nil, b501LegacySequenceError
		}
		var b501LegacySequenceAxes map[string]b501LegacySequenceJSON.RawMessage
		if b501LegacySequenceError := b501LegacySequenceJSON.Unmarshal(b501LegacySequenceObject["ab34_identity"], &b501LegacySequenceAxes); b501LegacySequenceError != nil {
			return nil, b501LegacySequenceError
		}
		b501LegacySequenceCanonical, b501LegacySequenceError := b501LegacySequenceJSON.Marshal(b501LegacySequenceKey)
		if b501LegacySequenceError != nil {
			return nil, b501LegacySequenceError
		}
		var b501LegacySequenceRequired map[string]b501LegacySequenceJSON.RawMessage
		if b501LegacySequenceError = b501LegacySequenceJSON.Unmarshal(b501LegacySequenceCanonical, &b501LegacySequenceRequired); b501LegacySequenceError != nil {
			return nil, b501LegacySequenceError
		}
		for b501LegacySequenceAxis := range b501LegacySequenceRequired {
			if _, b501LegacySequencePresent := b501LegacySequenceAxes[b501LegacySequenceAxis]; !b501LegacySequencePresent {
				return nil, b501LegacySequenceFmt.Errorf("ab34 missing identity axis %s", b501LegacySequenceAxis)
			}
		}
		if len(b501LegacySequenceAxes) != len(b501LegacySequenceRequired) || string(b501LegacySequenceAxes["ab34_primary"]) == "null" || b501LegacySequenceKey.B501LegacySequenceFamily == "" || b501LegacySequenceKey.B501LegacySequenceVariant == "" || b501LegacySequenceKey.B501LegacySequenceDirection == "" || b501LegacySequenceKey.B501LegacySequenceField == "" || b501LegacySequenceKey.B501LegacySequenceAlternative == "" || b501LegacySequenceKey.B501LegacySequenceInstanceScope == "" {
			return nil, b501LegacySequenceFmt.Errorf("ab34 incomplete identity")
		}
		if b501LegacySequenceSeen[string(b501LegacySequenceCanonical)] {
			return nil, b501LegacySequenceFmt.Errorf("ab34 duplicate identity")
		}
		b501LegacySequenceSeen[string(b501LegacySequenceCanonical)] = true
	}
	return b501LegacySequenceRows, nil
}

func B501LegacySequenceEntries() (b501LegacySequenceJSON.RawMessage, error) {
	b501LegacySequenceData := []byte(b501LegacySequenceEmbedded)
	if _, b501LegacySequenceError := b501LegacySequenceValidate(b501LegacySequenceData); b501LegacySequenceError != nil {
		return nil, b501LegacySequenceError
	}
	return append(b501LegacySequenceJSON.RawMessage(nil), b501LegacySequenceData...), nil
}

func b501LegacySequenceLookup(b501LegacySequenceData []byte, b501LegacySequenceKey B501LegacySequenceKey) (b501LegacySequenceJSON.RawMessage, bool, error) {
	b501LegacySequenceRows, b501LegacySequenceError := b501LegacySequenceValidate(b501LegacySequenceData)
	if b501LegacySequenceError != nil {
		return nil, false, b501LegacySequenceError
	}
	b501LegacySequenceWant, b501LegacySequenceError := b501LegacySequenceJSON.Marshal(b501LegacySequenceKey)
	if b501LegacySequenceError != nil {
		return nil, false, b501LegacySequenceError
	}
	for _, b501LegacySequenceRow := range b501LegacySequenceRows {
		var b501LegacySequenceObject struct {
			B501LegacySequenceIdentity B501LegacySequenceKey `json:"ab34_identity"`
		}
		if b501LegacySequenceError = b501LegacySequenceJSON.Unmarshal(b501LegacySequenceRow, &b501LegacySequenceObject); b501LegacySequenceError != nil {
			return nil, false, b501LegacySequenceError
		}
		b501LegacySequenceIdentity, b501LegacySequenceError := b501LegacySequenceJSON.Marshal(b501LegacySequenceObject.B501LegacySequenceIdentity)
		if b501LegacySequenceError != nil {
			return nil, false, b501LegacySequenceError
		}
		if string(b501LegacySequenceIdentity) == string(b501LegacySequenceWant) {
			return append(b501LegacySequenceJSON.RawMessage(nil), b501LegacySequenceRow...), true, nil
		}
	}
	return nil, false, nil
}

func B501LegacySequenceLookup(b501LegacySequenceKey B501LegacySequenceKey) (b501LegacySequenceJSON.RawMessage, bool, error) {
	return b501LegacySequenceLookup([]byte(b501LegacySequenceEmbedded), b501LegacySequenceKey)
}
