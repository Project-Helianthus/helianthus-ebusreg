package vaillant

import (
	_ "embed"
	"encoding/json"
	"reflect"
)

type B555NumericUnitKey struct {
	B555NumericUnitFamily    string  `json:"family"`
	B555NumericUnitVariant   string  `json:"variant"`
	B555NumericUnitPrimary   uint16  `json:"primary"`
	B555NumericUnitSecondary *uint8  `json:"secondary"`
	B555NumericUnitOperation *uint8  `json:"operation"`
	B555NumericUnitDirection *string `json:"direction"`
	B555NumericUnitField     string  `json:"field"`
	B555NumericUnitOpcode    *uint8  `json:"opcode"`
	B555NumericUnitGroup     *uint8  `json:"group"`
	B555NumericUnitInstance  *uint8  `json:"instance"`
	B555NumericUnitRegister  *uint16 `json:"register"`
}

type B555NumericUnitReplacement struct {
	B555NumericUnitRaw     uint64      `json:"raw"`
	B555NumericUnitRole    string      `json:"role"`
	B555NumericUnitContext string      `json:"context"`
	B555NumericUnitGrade   SchemaGrade `json:"grade"`
}

type B555NumericField struct {
	B555NumericUnitKey             B555NumericUnitKey           `json:"key"`
	B555NumericUnitLogicalPosition *int                         `json:"logical_position"`
	B555NumericUnitWidthBits       *int                         `json:"width_bits"`
	B555NumericUnitSignedness      string                       `json:"signedness"`
	B555NumericUnitByteOrder       string                       `json:"byte_order"`
	B555NumericUnitConversion      string                       `json:"conversion"`
	B555NumericUnitDivisor         *float64                     `json:"divisor"`
	B555NumericUnit                *string                      `json:"unit"`
	B555NumericUnitCode            *uint8                       `json:"unit_code"`
	B555NumericUnitDefault         *float64                     `json:"default"`
	B555NumericUnitRawMin          *uint64                      `json:"raw_min"`
	B555NumericUnitRawMax          *uint64                      `json:"raw_max"`
	B555NumericMin                 *float64                     `json:"numeric_min"`
	B555NumericMax                 *float64                     `json:"numeric_max"`
	B555NumericUnitEnumeration     map[uint64]string            `json:"enumeration"`
	B555NumericUnitReplacement     []B555NumericUnitReplacement `json:"replacement"`
	B555NumericUnitReserved        []uint64                     `json:"reserved"`
	B555NumericUnitGrade           SchemaGrade                  `json:"grade"`
	B555NumericUnitConversionGrade SchemaGrade                  `json:"conversion_grade"`
	B555NumericUnitUnitGrade       SchemaGrade                  `json:"unit_grade"`
	B555NumericUnitRawFallback     string                       `json:"raw_fallback"`
	B555NumericUnitEnforcement     SchemaGrade                  `json:"enforcement"`
	B555NumericUnitPeerEffects     SchemaGrade                  `json:"peer_effects"`
}

type B555NumericMetadata struct {
	B555NumericUnitRepresentationGrade   SchemaGrade        `json:"representation_grade"`
	B555NumericUnitPhysicalQualification SchemaGrade        `json:"physical_qualification"`
	B555NumericUnitNativeSupport         SchemaGrade        `json:"native_support"`
	B555NumericUnitFreshness             SchemaGrade        `json:"freshness"`
	B555NumericUnitPersistence           SchemaGrade        `json:"persistence"`
	B555NumericUnitExecution             SchemaGrade        `json:"execution"`
	B555NumericUnitComplete              bool               `json:"complete"`
	B555NumericUnitSourceContract        string             `json:"source_contract"`
	B555NumericUnitLimitations           []string           `json:"limitations"`
	B555NumericUnitEntries               []B555NumericField `json:"entries"`
}

//go:embed data/ab94_b555_numeric_units.json
var b555NumericUnitNumericJSON string

func B555NumericUnits() (B555NumericMetadata, error) {
	var b555NumericUnitMetadata B555NumericMetadata
	b555NumericUnitError := json.Unmarshal([]byte(b555NumericUnitNumericJSON), &b555NumericUnitMetadata)
	return b555NumericUnitMetadata, b555NumericUnitError
}

func b555NumericUnitLookup(b555NumericUnitEntries []B555NumericField, b555NumericUnitKey B555NumericUnitKey) (B555NumericField, bool) {
	for _, b555NumericUnitEntry := range b555NumericUnitEntries {
		if reflect.DeepEqual(b555NumericUnitEntry.B555NumericUnitKey, b555NumericUnitKey) {
			return b555NumericUnitEntry, true
		}
	}
	return B555NumericField{}, false
}

func B555LookupNumericUnits(b555NumericUnitKey B555NumericUnitKey) (B555NumericField, bool, error) {
	b555NumericUnitMetadata, b555NumericUnitError := B555NumericUnits()
	if b555NumericUnitError != nil {
		return B555NumericField{}, false, b555NumericUnitError
	}
	b555NumericUnitField, b555NumericUnitFound := b555NumericUnitLookup(b555NumericUnitMetadata.B555NumericUnitEntries, b555NumericUnitKey)
	return b555NumericUnitField, b555NumericUnitFound, nil
}
