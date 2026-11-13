package vaillant

import (
	b511ClassCandidateBytes "bytes"
	_ "embed"
	b511ClassCandidateJSON "encoding/json"
	b511ClassCandidateFmt "fmt"
)

//go:embed data/ab50_b511_class_candidates.json
var b511ClassCandidateData string

type B511ClassCandidateKey struct {
	B511ClassCandidateFamily        string  `json:"family"`
	B511ClassCandidateVariant       string  `json:"variant"`
	B511ClassCandidatePrimary       uint16  `json:"primary"`
	B511ClassCandidateSecondary     *uint16 `json:"secondary"`
	B511ClassCandidateOperation     *uint16 `json:"operation"`
	B511ClassCandidateDirection     string  `json:"direction"`
	B511ClassCandidateField         string  `json:"field"`
	B511ClassCandidateAlternative   string  `json:"alternative"`
	B511ClassCandidateGroup         *uint16 `json:"group"`
	B511ClassCandidateInstanceScope string  `json:"instance_scope"`
	B511ClassCandidateRegister      *uint16 `json:"register"`
}

func b511ClassCandidateValidate(b511ClassCandidateSource []byte) ([]b511ClassCandidateJSON.RawMessage, []B511ClassCandidateKey, error) {
	var b511ClassCandidateRows []b511ClassCandidateJSON.RawMessage
	if b511ClassCandidateError := b511ClassCandidateJSON.Unmarshal(b511ClassCandidateSource, &b511ClassCandidateRows); b511ClassCandidateError != nil {
		return nil, nil, b511ClassCandidateError
	}
	if b511ClassCandidateRows == nil {
		return nil, nil, b511ClassCandidateFmt.Errorf("ab50: expected array")
	}
	b511ClassCandidateKeys := make([]B511ClassCandidateKey, len(b511ClassCandidateRows))
	b511ClassCandidateSeen := make(map[string]bool)
	for b511ClassCandidateIndex, b511ClassCandidateRow := range b511ClassCandidateRows {
		var b511ClassCandidateObject map[string]b511ClassCandidateJSON.RawMessage
		if b511ClassCandidateError := b511ClassCandidateJSON.Unmarshal(b511ClassCandidateRow, &b511ClassCandidateObject); b511ClassCandidateError != nil {
			return nil, nil, b511ClassCandidateError
		}
		b511ClassCandidateKeyJSON, b511ClassCandidateOK := b511ClassCandidateObject["key"]
		if !b511ClassCandidateOK {
			return nil, nil, b511ClassCandidateFmt.Errorf("ab50: missing key")
		}
		var b511ClassCandidateAxes map[string]b511ClassCandidateJSON.RawMessage
		if b511ClassCandidateError := b511ClassCandidateJSON.Unmarshal(b511ClassCandidateKeyJSON, &b511ClassCandidateAxes); b511ClassCandidateError != nil {
			return nil, nil, b511ClassCandidateError
		}
		if b511ClassCandidatePrimary, b511ClassCandidateOK := b511ClassCandidateAxes["primary"]; !b511ClassCandidateOK || b511ClassCandidateBytes.Equal(b511ClassCandidateBytes.TrimSpace(b511ClassCandidatePrimary), []byte("null")) {
			return nil, nil, b511ClassCandidateFmt.Errorf("ab50: missing primary")
		}
		b511ClassCandidateDecoder := b511ClassCandidateJSON.NewDecoder(b511ClassCandidateBytes.NewReader(b511ClassCandidateKeyJSON))
		b511ClassCandidateDecoder.DisallowUnknownFields()
		if b511ClassCandidateError := b511ClassCandidateDecoder.Decode(&b511ClassCandidateKeys[b511ClassCandidateIndex]); b511ClassCandidateError != nil {
			return nil, nil, b511ClassCandidateError
		}
		b511ClassCandidateKey := b511ClassCandidateKeys[b511ClassCandidateIndex]
		if b511ClassCandidateKey.B511ClassCandidateFamily == "" || b511ClassCandidateKey.B511ClassCandidateVariant == "" || b511ClassCandidateKey.B511ClassCandidateDirection == "" || b511ClassCandidateKey.B511ClassCandidateField == "" || b511ClassCandidateKey.B511ClassCandidateAlternative == "" || b511ClassCandidateKey.B511ClassCandidateInstanceScope == "" {
			return nil, nil, b511ClassCandidateFmt.Errorf("ab50: incomplete key")
		}
		b511ClassCandidateCanonical, b511ClassCandidateError := b511ClassCandidateJSON.Marshal(b511ClassCandidateKey)
		if b511ClassCandidateError != nil {
			return nil, nil, b511ClassCandidateError
		}
		if b511ClassCandidateSeen[string(b511ClassCandidateCanonical)] {
			return nil, nil, b511ClassCandidateFmt.Errorf("ab50: duplicate key")
		}
		b511ClassCandidateSeen[string(b511ClassCandidateCanonical)] = true
	}
	return b511ClassCandidateRows, b511ClassCandidateKeys, nil
}

func b511ClassCandidateEntries(b511ClassCandidateSource []byte) (b511ClassCandidateJSON.RawMessage, error) {
	if _, _, b511ClassCandidateError := b511ClassCandidateValidate(b511ClassCandidateSource); b511ClassCandidateError != nil {
		return nil, b511ClassCandidateError
	}
	return append(b511ClassCandidateJSON.RawMessage(nil), b511ClassCandidateSource...), nil
}

func B511ClassCandidateEntries() (b511ClassCandidateJSON.RawMessage, error) {
	return b511ClassCandidateEntries([]byte(b511ClassCandidateData))
}

func b511ClassCandidateLookup(b511ClassCandidateSource []byte, b511ClassCandidateKey B511ClassCandidateKey) (b511ClassCandidateJSON.RawMessage, bool, error) {
	b511ClassCandidateRows, b511ClassCandidateKeys, b511ClassCandidateError := b511ClassCandidateValidate(b511ClassCandidateSource)
	if b511ClassCandidateError != nil {
		return nil, false, b511ClassCandidateError
	}
	b511ClassCandidateWanted, b511ClassCandidateError := b511ClassCandidateJSON.Marshal(b511ClassCandidateKey)
	if b511ClassCandidateError != nil {
		return nil, false, b511ClassCandidateError
	}
	for b511ClassCandidateIndex, b511ClassCandidateCandidate := range b511ClassCandidateKeys {
		b511ClassCandidateCanonical, b511ClassCandidateError := b511ClassCandidateJSON.Marshal(b511ClassCandidateCandidate)
		if b511ClassCandidateError != nil {
			return nil, false, b511ClassCandidateError
		}
		if b511ClassCandidateBytes.Equal(b511ClassCandidateWanted, b511ClassCandidateCanonical) {
			return append(b511ClassCandidateJSON.RawMessage(nil), b511ClassCandidateRows[b511ClassCandidateIndex]...), true, nil
		}
	}
	return nil, false, nil
}

func B511ClassCandidateLookup(b511ClassCandidateKey B511ClassCandidateKey) (b511ClassCandidateJSON.RawMessage, bool, error) {
	return b511ClassCandidateLookup([]byte(b511ClassCandidateData), b511ClassCandidateKey)
}
