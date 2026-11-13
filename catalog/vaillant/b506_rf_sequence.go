package vaillant

import (
	b506RFSequenceBytes "bytes"
	_ "embed"
	b506RFSequenceJSON "encoding/json"
	b506RFSequenceFmt "fmt"
	b506RFSequenceStrings "strings"
)

//go:embed data/ab17_b506_rf_sequence.json
var b506RFSequenceCatalog string

type B506RFSequenceKey struct {
	B506RFSequenceFamily        string  `json:"family"`
	B506RFSequenceVariant       string  `json:"variant"`
	B506RFSequencePrimary       uint8   `json:"primary"`
	B506RFSequenceSecondary     *uint8  `json:"secondary"`
	B506RFSequenceOperation     *uint8  `json:"operation"`
	B506RFSequenceDirection     string  `json:"direction"`
	B506RFSequenceField         string  `json:"field"`
	B506RFSequenceAlternative   string  `json:"alternative"`
	B506RFSequenceGroup         *uint8  `json:"group"`
	B506RFSequenceInstanceScope string  `json:"instance_scope"`
	B506RFSequenceRegister      *uint16 `json:"register"`
}

type b506RFSequenceRow struct {
	B506RFSequenceID            string            `json:"ab17_id"`
	B506RFSequenceKey           B506RFSequenceKey `json:"key"`
	B506RFSequenceGrade         string            `json:"grade"`
	B506RFSequenceQualification string            `json:"qualification"`
	B506RFSequenceRawForm       string            `json:"ab17_raw_form"`
	B506RFSequenceUnknown       map[string]string `json:"ab17_unknown"`
}

func b506RFSequenceEntriesFrom(b506RFSequenceData string) ([]b506RFSequenceJSON.RawMessage, error) {
	var b506RFSequenceRows []b506RFSequenceJSON.RawMessage
	if b506RFSequenceError := b506RFSequenceJSON.Unmarshal([]byte(b506RFSequenceData), &b506RFSequenceRows); b506RFSequenceError != nil {
		return nil, b506RFSequenceError
	}
	if b506RFSequenceRows == nil {
		return nil, b506RFSequenceFmt.Errorf("ab17: expected array")
	}
	b506RFSequenceSeen := make(map[string]bool)
	for _, b506RFSequenceRaw := range b506RFSequenceRows {
		var b506RFSequenceEntry b506RFSequenceRow
		if b506RFSequenceError := b506RFSequenceJSON.Unmarshal(b506RFSequenceRaw, &b506RFSequenceEntry); b506RFSequenceError != nil {
			return nil, b506RFSequenceError
		}
		var b506RFSequenceObject map[string]b506RFSequenceJSON.RawMessage
		if b506RFSequenceError := b506RFSequenceJSON.Unmarshal(b506RFSequenceRaw, &b506RFSequenceObject); b506RFSequenceError != nil {
			return nil, b506RFSequenceError
		}
		var b506RFSequenceAxes map[string]b506RFSequenceJSON.RawMessage
		if b506RFSequenceError := b506RFSequenceJSON.Unmarshal(b506RFSequenceObject["key"], &b506RFSequenceAxes); b506RFSequenceError != nil {
			return nil, b506RFSequenceError
		}
		for _, b506RFSequenceAxis := range []string{"family", "variant", "primary", "secondary", "operation", "direction", "field", "alternative", "group", "instance_scope", "register"} {
			if _, b506RFSequencePresent := b506RFSequenceAxes[b506RFSequenceAxis]; !b506RFSequencePresent {
				return nil, b506RFSequenceFmt.Errorf("ab17: missing axis %s", b506RFSequenceAxis)
			}
			if b506RFSequenceAxis != "secondary" && b506RFSequenceAxis != "operation" && b506RFSequenceAxis != "group" && b506RFSequenceAxis != "register" && b506RFSequenceBytes.Equal(b506RFSequenceAxes[b506RFSequenceAxis], []byte("null")) {
				return nil, b506RFSequenceFmt.Errorf("ab17: null axis %s", b506RFSequenceAxis)
			}
		}
		b506RFSequenceDecoder := b506RFSequenceJSON.NewDecoder(b506RFSequenceBytes.NewReader(b506RFSequenceObject["key"]))
		b506RFSequenceDecoder.DisallowUnknownFields()
		if b506RFSequenceError := b506RFSequenceDecoder.Decode(&b506RFSequenceEntry.B506RFSequenceKey); b506RFSequenceError != nil {
			return nil, b506RFSequenceError
		}
		if !b506RFSequenceStrings.HasPrefix(b506RFSequenceEntry.B506RFSequenceID, "ab17-") || b506RFSequenceEntry.B506RFSequenceGrade != "Hypothesis" || b506RFSequenceEntry.B506RFSequenceQualification != "Unknown" || b506RFSequenceEntry.B506RFSequenceRawForm == "" || len(b506RFSequenceEntry.B506RFSequenceUnknown) == 0 || b506RFSequenceEntry.B506RFSequenceKey.B506RFSequenceFamily == "" || b506RFSequenceEntry.B506RFSequenceKey.B506RFSequenceVariant == "" || b506RFSequenceEntry.B506RFSequenceKey.B506RFSequenceDirection == "" || b506RFSequenceEntry.B506RFSequenceKey.B506RFSequenceField == "" || b506RFSequenceEntry.B506RFSequenceKey.B506RFSequenceAlternative == "" || b506RFSequenceEntry.B506RFSequenceKey.B506RFSequenceInstanceScope == "" {
			return nil, b506RFSequenceFmt.Errorf("ab17: incomplete metadata")
		}
		b506RFSequenceIdentity, b506RFSequenceError := b506RFSequenceJSON.Marshal(b506RFSequenceEntry.B506RFSequenceKey)
		if b506RFSequenceError != nil {
			return nil, b506RFSequenceError
		}
		if b506RFSequenceSeen[string(b506RFSequenceIdentity)] {
			return nil, b506RFSequenceFmt.Errorf("ab17: duplicate key")
		}
		b506RFSequenceSeen[string(b506RFSequenceIdentity)] = true
	}
	return b506RFSequenceRows, nil
}

func B506RFSequenceEntries() (b506RFSequenceJSON.RawMessage, error) {
	if _, b506RFSequenceError := b506RFSequenceEntriesFrom(b506RFSequenceCatalog); b506RFSequenceError != nil {
		return nil, b506RFSequenceError
	}
	return b506RFSequenceJSON.RawMessage([]byte(b506RFSequenceCatalog)), nil
}

func b506RFSequenceLookupFrom(b506RFSequenceData string, b506RFSequenceKey B506RFSequenceKey) (b506RFSequenceJSON.RawMessage, bool, error) {
	b506RFSequenceRows, b506RFSequenceError := b506RFSequenceEntriesFrom(b506RFSequenceData)
	if b506RFSequenceError != nil {
		return nil, false, b506RFSequenceError
	}
	b506RFSequenceWanted, b506RFSequenceError := b506RFSequenceJSON.Marshal(b506RFSequenceKey)
	if b506RFSequenceError != nil {
		return nil, false, b506RFSequenceError
	}
	for _, b506RFSequenceRaw := range b506RFSequenceRows {
		var b506RFSequenceEntry b506RFSequenceRow
		if b506RFSequenceError := b506RFSequenceJSON.Unmarshal(b506RFSequenceRaw, &b506RFSequenceEntry); b506RFSequenceError != nil {
			return nil, false, b506RFSequenceError
		}
		b506RFSequenceIdentity, b506RFSequenceError := b506RFSequenceJSON.Marshal(b506RFSequenceEntry.B506RFSequenceKey)
		if b506RFSequenceError != nil {
			return nil, false, b506RFSequenceError
		}
		if b506RFSequenceBytes.Equal(b506RFSequenceWanted, b506RFSequenceIdentity) {
			return append(b506RFSequenceJSON.RawMessage(nil), b506RFSequenceRaw...), true, nil
		}
	}
	return nil, false, nil
}

func B506RFSequenceLookup(b506RFSequenceKey B506RFSequenceKey) (b506RFSequenceJSON.RawMessage, bool, error) {
	return b506RFSequenceLookupFrom(b506RFSequenceCatalog, b506RFSequenceKey)
}
