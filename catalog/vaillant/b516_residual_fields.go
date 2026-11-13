package vaillant

import (
	b516ResidualFieldBytes "bytes"
	_ "embed"
	b516ResidualFieldJSON "encoding/json"
	b516ResidualFieldFmt "fmt"
	b516ResidualFieldReflect "reflect"
)

//go:embed data/ab51_b516_residual_fields.json
var b516ResidualFieldData string

type B516ResidualFieldKey struct {
	B516ResidualFieldFamily        string  `json:"family"`
	B516ResidualFieldVariant       string  `json:"variant"`
	B516ResidualFieldPrimary       uint16  `json:"primary"`
	B516ResidualFieldSecondary     *uint16 `json:"secondary"`
	B516ResidualFieldOperation     *uint16 `json:"operation"`
	B516ResidualFieldDirection     string  `json:"direction"`
	B516ResidualField              string  `json:"field"`
	B516ResidualFieldAlternative   string  `json:"alternative"`
	B516ResidualFieldGroup         *uint16 `json:"group"`
	B516ResidualFieldInstanceScope string  `json:"instance_scope"`
	B516ResidualFieldRegister      *uint16 `json:"register"`
}

func b516ResidualFieldValidate(b516ResidualFieldSource []byte) ([]b516ResidualFieldJSON.RawMessage, []B516ResidualFieldKey, error) {
	var b516ResidualFieldRows []b516ResidualFieldJSON.RawMessage
	if b516ResidualFieldError := b516ResidualFieldJSON.Unmarshal(b516ResidualFieldSource, &b516ResidualFieldRows); b516ResidualFieldError != nil {
		return nil, nil, b516ResidualFieldError
	}
	if b516ResidualFieldRows == nil {
		return nil, nil, b516ResidualFieldFmt.Errorf("ab51: expected array")
	}
	b516ResidualFieldKeys := make([]B516ResidualFieldKey, 0, len(b516ResidualFieldRows))
	for _, b516ResidualFieldRow := range b516ResidualFieldRows {
		var b516ResidualFieldObject map[string]b516ResidualFieldJSON.RawMessage
		if b516ResidualFieldError := b516ResidualFieldJSON.Unmarshal(b516ResidualFieldRow, &b516ResidualFieldObject); b516ResidualFieldError != nil {
			return nil, nil, b516ResidualFieldError
		}
		var b516ResidualFieldMembers map[string]b516ResidualFieldJSON.RawMessage
		if b516ResidualFieldError := b516ResidualFieldJSON.Unmarshal(b516ResidualFieldObject["key"], &b516ResidualFieldMembers); b516ResidualFieldError != nil {
			return nil, nil, b516ResidualFieldError
		}
		for _, b516ResidualFieldName := range []string{"family", "variant", "primary", "direction", "field", "alternative", "instance_scope"} {
			if len(b516ResidualFieldMembers[b516ResidualFieldName]) == 0 || b516ResidualFieldBytes.Equal(b516ResidualFieldBytes.TrimSpace(b516ResidualFieldMembers[b516ResidualFieldName]), []byte("null")) {
				return nil, nil, b516ResidualFieldFmt.Errorf("ab51: missing key axis %s", b516ResidualFieldName)
			}
		}
		var b516ResidualFieldKey B516ResidualFieldKey
		b516ResidualFieldDecoder := b516ResidualFieldJSON.NewDecoder(b516ResidualFieldBytes.NewReader(b516ResidualFieldObject["key"]))
		b516ResidualFieldDecoder.DisallowUnknownFields()
		if b516ResidualFieldError := b516ResidualFieldDecoder.Decode(&b516ResidualFieldKey); b516ResidualFieldError != nil {
			return nil, nil, b516ResidualFieldError
		}
		if b516ResidualFieldKey.B516ResidualFieldFamily == "" || b516ResidualFieldKey.B516ResidualFieldVariant == "" || b516ResidualFieldKey.B516ResidualFieldDirection == "" || b516ResidualFieldKey.B516ResidualField == "" || b516ResidualFieldKey.B516ResidualFieldAlternative == "" || b516ResidualFieldKey.B516ResidualFieldInstanceScope == "" {
			return nil, nil, b516ResidualFieldFmt.Errorf("ab51: empty string axis; use none when not applicable")
		}
		for _, b516ResidualFieldPrevious := range b516ResidualFieldKeys {
			if b516ResidualFieldReflect.DeepEqual(b516ResidualFieldPrevious, b516ResidualFieldKey) {
				return nil, nil, b516ResidualFieldFmt.Errorf("ab51: duplicate key")
			}
		}
		b516ResidualFieldKeys = append(b516ResidualFieldKeys, b516ResidualFieldKey)
	}
	return b516ResidualFieldRows, b516ResidualFieldKeys, nil
}

func B516ResidualFieldEntries() (b516ResidualFieldJSON.RawMessage, error) {
	b516ResidualFieldSource := []byte(b516ResidualFieldData)
	if _, _, b516ResidualFieldError := b516ResidualFieldValidate(b516ResidualFieldSource); b516ResidualFieldError != nil {
		return nil, b516ResidualFieldError
	}
	return b516ResidualFieldJSON.RawMessage(b516ResidualFieldSource), nil
}

func b516ResidualFieldLookup(b516ResidualFieldSource []byte, b516ResidualFieldKey B516ResidualFieldKey) (b516ResidualFieldJSON.RawMessage, bool, error) {
	b516ResidualFieldRows, b516ResidualFieldKeys, b516ResidualFieldError := b516ResidualFieldValidate(b516ResidualFieldSource)
	if b516ResidualFieldError != nil {
		return nil, false, b516ResidualFieldError
	}
	for b516ResidualFieldIndex, b516ResidualFieldCandidate := range b516ResidualFieldKeys {
		if b516ResidualFieldReflect.DeepEqual(b516ResidualFieldCandidate, b516ResidualFieldKey) {
			return append(b516ResidualFieldJSON.RawMessage(nil), b516ResidualFieldRows[b516ResidualFieldIndex]...), true, nil
		}
	}
	return nil, false, nil
}

func B516ResidualFieldLookup(b516ResidualFieldKey B516ResidualFieldKey) (b516ResidualFieldJSON.RawMessage, bool, error) {
	return b516ResidualFieldLookup([]byte(b516ResidualFieldData), b516ResidualFieldKey)
}
