package vaillant

import (
	b50ERuleValueRowBytes "bytes"
	_ "embed"
	b50ERuleValueRowJSON "encoding/json"
	b50ERuleValueRowFmt "fmt"
	b50ERuleValueRowIO "io"
	b50ERuleValueRowReflect "reflect"
)

//go:embed data/ab74_b50e_rule_value_rows.json
var b50ERuleValueRowRowsJSON []byte

type B50ERuleValueRowKey struct {
	B50ERuleValueRowFamily    *string `json:"family"`
	B50ERuleValueRowVariant   *string `json:"variant"`
	B50ERuleValueRowPrimary   *byte   `json:"primary"`
	B50ERuleValueRowSecondary *byte   `json:"secondary"`
	B50ERuleValueRowOperation *byte   `json:"operation"`
	B50ERuleValueRowDirection *string `json:"direction"`
	B50ERuleValueRowSelector  *byte   `json:"selector"`
	B50ERuleValueRowField     *string `json:"field"`
}

type B50ERuleValueRowField struct {
	B50ERuleValueRowOffset         *int              `json:"offset"`
	B50ERuleValueRowWidth          *int              `json:"width"`
	B50ERuleValueRowEncoding       string            `json:"encoding"`
	B50ERuleValueRowByteOrder      string            `json:"byte_order"`
	B50ERuleValueRowDefault        *float64          `json:"default"`
	B50ERuleValueRowEnumeration    map[string]string `json:"enumeration"`
	B50ERuleValueRowScale          *float64          `json:"scale"`
	B50ERuleValueRowUnit           *string           `json:"unit"`
	B50ERuleValueRowCarrierMinimum *int              `json:"carrier_minimum"`
	B50ERuleValueRowCarrierMaximum *int              `json:"carrier_maximum"`
	B50ERuleValueRowScaledMinimum  *float64          `json:"scaled_minimum"`
	B50ERuleValueRowScaledMaximum  *float64          `json:"scaled_maximum"`
}

type B50ERuleValueRowEntry struct {
	B50ERuleValueRowKey                  B50ERuleValueRowKey    `json:"key"`
	B50ERuleValue                        B50ERuleValueRowField  `json:"value"`
	B50ERuleValueRowCompanion            *B50ERuleValueRowField `json:"companion"`
	B50ERuleValueRowBase                 string                 `json:"base"`
	B50ERuleValueRowGrade                string                 `json:"grade"`
	B50ERuleValueRowMetadataPresence     string                 `json:"metadata_presence"`
	B50ERuleValueRowPhysicalSupport      string                 `json:"physical_support"`
	B50ERuleValueRowEffects              string                 `json:"effects"`
	B50ERuleValueRowPersistence          string                 `json:"persistence"`
	B50ERuleValueRowDelivery             string                 `json:"delivery"`
	B50ERuleValueRowDefaultQualification string                 `json:"default_qualification"`
	B50ERuleValueRowRawFallback          string                 `json:"raw_fallback"`
}

func b50ERuleValueRowDecode(b50ERuleValueRowInput []byte) ([]B50ERuleValueRowEntry, error) {
	var b50ERuleValueRowEntries []B50ERuleValueRowEntry
	b50ERuleValueRowDecoder := b50ERuleValueRowJSON.NewDecoder(b50ERuleValueRowBytes.NewReader(b50ERuleValueRowInput))
	b50ERuleValueRowDecoder.DisallowUnknownFields()
	if b50ERuleValueRowError := b50ERuleValueRowDecoder.Decode(&b50ERuleValueRowEntries); b50ERuleValueRowError != nil {
		return nil, b50ERuleValueRowError
	}
	if b50ERuleValueRowError := b50ERuleValueRowDecoder.Decode(new(any)); b50ERuleValueRowError != b50ERuleValueRowIO.EOF {
		return nil, b50ERuleValueRowFmt.Errorf("ab74 trailing JSON: %v", b50ERuleValueRowError)
	}
	b50ERuleValueRowCanonical, b50ERuleValueRowError := b50ERuleValueRowJSON.Marshal(b50ERuleValueRowEntries)
	if b50ERuleValueRowError != nil {
		return nil, b50ERuleValueRowError
	}
	var b50ERuleValueRowRawShape, b50ERuleValueRowTypedShape any
	if b50ERuleValueRowError = b50ERuleValueRowJSON.Unmarshal(b50ERuleValueRowInput, &b50ERuleValueRowRawShape); b50ERuleValueRowError != nil {
		return nil, b50ERuleValueRowError
	}
	if b50ERuleValueRowError = b50ERuleValueRowJSON.Unmarshal(b50ERuleValueRowCanonical, &b50ERuleValueRowTypedShape); b50ERuleValueRowError != nil {
		return nil, b50ERuleValueRowError
	}
	if !b50ERuleValueRowCompleteShape(b50ERuleValueRowRawShape, b50ERuleValueRowTypedShape) {
		return nil, b50ERuleValueRowFmt.Errorf("ab74 missing explicit metadata field")
	}
	if len(b50ERuleValueRowEntries) != 19 {
		return nil, b50ERuleValueRowFmt.Errorf("ab74 expected nineteen rows")
	}
	b50ERuleValueRowSeen := make(map[byte]bool)
	for _, b50ERuleValueRowEntry := range b50ERuleValueRowEntries {
		b50ERuleValueRowKey := b50ERuleValueRowEntry.B50ERuleValueRowKey
		if b50ERuleValueRowKey.B50ERuleValueRowFamily == nil || *b50ERuleValueRowKey.B50ERuleValueRowFamily != "vaillant" || b50ERuleValueRowKey.B50ERuleValueRowVariant == nil || *b50ERuleValueRowKey.B50ERuleValueRowVariant != "B50E" || b50ERuleValueRowKey.B50ERuleValueRowPrimary == nil || *b50ERuleValueRowKey.B50ERuleValueRowPrimary != 0xB5 || b50ERuleValueRowKey.B50ERuleValueRowSecondary == nil || *b50ERuleValueRowKey.B50ERuleValueRowSecondary != 0x0E || b50ERuleValueRowKey.B50ERuleValueRowOperation == nil || *b50ERuleValueRowKey.B50ERuleValueRowOperation != 2 || b50ERuleValueRowKey.B50ERuleValueRowDirection == nil || *b50ERuleValueRowKey.B50ERuleValueRowDirection != "positive" || b50ERuleValueRowKey.B50ERuleValueRowField == nil || *b50ERuleValueRowKey.B50ERuleValueRowField != "rule_value_row" || b50ERuleValueRowKey.B50ERuleValueRowSelector == nil || *b50ERuleValueRowKey.B50ERuleValueRowSelector > 18 {
			return nil, b50ERuleValueRowFmt.Errorf("ab74 invalid row identity")
		}
		b50ERuleValueRowSelector := *b50ERuleValueRowKey.B50ERuleValueRowSelector
		if b50ERuleValueRowSeen[b50ERuleValueRowSelector] {
			return nil, b50ERuleValueRowFmt.Errorf("ab74 duplicate row key")
		}
		b50ERuleValueRowSeen[b50ERuleValueRowSelector] = true
		if b50ERuleValueRowEntry.B50ERuleValueRowBase != "row-relative" || b50ERuleValueRowEntry.B50ERuleValueRowGrade != "Hypothesis" || b50ERuleValueRowEntry.B50ERuleValueRowMetadataPresence != "Proven" || b50ERuleValueRowEntry.B50ERuleValueRowPhysicalSupport != "Unknown" || b50ERuleValueRowEntry.B50ERuleValueRowEffects != "Unknown" || b50ERuleValueRowEntry.B50ERuleValueRowPersistence != "Unknown" || b50ERuleValueRowEntry.B50ERuleValueRowDelivery != "Unknown" || b50ERuleValueRowEntry.B50ERuleValueRowDefaultQualification != "Unknown" || b50ERuleValueRowEntry.B50ERuleValueRowRawFallback == "" {
			return nil, b50ERuleValueRowFmt.Errorf("ab74 incomplete qualification")
		}
		if b50ERuleValueRowEntry.B50ERuleValue.B50ERuleValueRowOffset == nil || *b50ERuleValueRowEntry.B50ERuleValue.B50ERuleValueRowOffset != 0 {
			return nil, b50ERuleValueRowFmt.Errorf("ab74 invalid value base")
		}
		for _, b50ERuleValueRowField := range []*B50ERuleValueRowField{&b50ERuleValueRowEntry.B50ERuleValue, b50ERuleValueRowEntry.B50ERuleValueRowCompanion} {
			if b50ERuleValueRowField == nil {
				continue
			}
			if b50ERuleValueRowField.B50ERuleValueRowOffset == nil || *b50ERuleValueRowField.B50ERuleValueRowOffset < 0 || b50ERuleValueRowField.B50ERuleValueRowEncoding == "" || b50ERuleValueRowField.B50ERuleValueRowByteOrder != "Unknown" || b50ERuleValueRowField.B50ERuleValueRowDefault != nil || (b50ERuleValueRowField.B50ERuleValueRowWidth != nil && *b50ERuleValueRowField.B50ERuleValueRowWidth <= 0) {
				return nil, b50ERuleValueRowFmt.Errorf("ab74 invalid field declaration")
			}
		}
		if b50ERuleValueRowSelector == 0 || b50ERuleValueRowSelector == 15 || b50ERuleValueRowSelector == 16 {
			if b50ERuleValueRowEntry.B50ERuleValue.B50ERuleValueRowWidth != nil || b50ERuleValueRowEntry.B50ERuleValue.B50ERuleValueRowEncoding != "opaque" {
				return nil, b50ERuleValueRowFmt.Errorf("ab74 opaque width must remain absent")
			}
		} else if b50ERuleValueRowEntry.B50ERuleValue.B50ERuleValueRowWidth == nil {
			return nil, b50ERuleValueRowFmt.Errorf("ab74 missing known width")
		}
		if (b50ERuleValueRowSelector >= 9 && b50ERuleValueRowSelector <= 14 || b50ERuleValueRowSelector >= 17) != (b50ERuleValueRowEntry.B50ERuleValueRowCompanion != nil) {
			return nil, b50ERuleValueRowFmt.Errorf("ab74 missing or extra companion")
		}
	}
	return b50ERuleValueRowEntries, nil
}

func B50ERuleValueRowEntries() ([]B50ERuleValueRowEntry, error) {
	return b50ERuleValueRowDecode(b50ERuleValueRowRowsJSON)
}

func B50ERuleValueRowLookup(b50ERuleValueRowKey B50ERuleValueRowKey) (B50ERuleValueRowEntry, bool, error) {
	b50ERuleValueRowEntries, b50ERuleValueRowError := B50ERuleValueRowEntries()
	if b50ERuleValueRowError != nil {
		return B50ERuleValueRowEntry{}, false, b50ERuleValueRowError
	}
	for _, b50ERuleValueRowEntry := range b50ERuleValueRowEntries {
		if b50ERuleValueRowReflect.DeepEqual(b50ERuleValueRowEntry.B50ERuleValueRowKey, b50ERuleValueRowKey) {
			return b50ERuleValueRowEntry, true, nil
		}
	}
	return B50ERuleValueRowEntry{}, false, nil
}

func b50ERuleValueRowCompleteShape(b50ERuleValueRowRaw, b50ERuleValueRowTyped any) bool {
	switch b50ERuleValueRowExpected := b50ERuleValueRowTyped.(type) {
	case map[string]any:
		b50ERuleValueRowActual, b50ERuleValueRowOK := b50ERuleValueRowRaw.(map[string]any)
		if !b50ERuleValueRowOK || len(b50ERuleValueRowActual) != len(b50ERuleValueRowExpected) {
			return false
		}
		for b50ERuleValueRowName, b50ERuleValueRowValue := range b50ERuleValueRowExpected {
			b50ERuleValueRowInput, b50ERuleValueRowPresent := b50ERuleValueRowActual[b50ERuleValueRowName]
			if !b50ERuleValueRowPresent || !b50ERuleValueRowCompleteShape(b50ERuleValueRowInput, b50ERuleValueRowValue) {
				return false
			}
		}
	case []any:
		b50ERuleValueRowActual, b50ERuleValueRowOK := b50ERuleValueRowRaw.([]any)
		if !b50ERuleValueRowOK || len(b50ERuleValueRowActual) != len(b50ERuleValueRowExpected) {
			return false
		}
		for b50ERuleValueRowIndex := range b50ERuleValueRowExpected {
			if !b50ERuleValueRowCompleteShape(b50ERuleValueRowActual[b50ERuleValueRowIndex], b50ERuleValueRowExpected[b50ERuleValueRowIndex]) {
				return false
			}
		}
	}
	return true
}
