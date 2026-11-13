package vaillant

import (
	b509Op59QuantityBytes "bytes"
	_ "embed"
	b509Op59QuantityJSON "encoding/json"
	b509Op59QuantityFmt "fmt"
	b509Op59QuantityIO "io"
	b509Op59QuantityReflect "reflect"
)

//go:embed data/ab14_b509_op59_quantities.json
var b509Op59QuantityData string

type B509Op59QuantityKey struct {
	B509Op59QuantityFamily            string  `json:"family"`
	B509Op59QuantityVariant           string  `json:"variant"`
	B509Op59QuantityPrimary           uint8   `json:"primary"`
	B509Op59QuantitySecondary         *uint8  `json:"secondary"`
	B509Op59QuantityOperation         *uint8  `json:"operation"`
	B509Op59QuantityDirection         *string `json:"direction"`
	B509Op59QuantityField             string  `json:"field"`
	B509Op59QuantityGroup             *uint8  `json:"group"`
	B509Op59QuantityRequestedInstance *uint8  `json:"requested_instance"`
	B509Op59QuantityRegister          *uint16 `json:"register"`
	B509Op59QuantityCommandRole       *string `json:"command_role"`
}

type B509Op59QuantityCodec struct {
	B509Op59QuantityRepresentation string   `json:"representation"`
	B509Op59QuantityWidth          *uint8   `json:"width_bits"`
	B509Op59QuantityDivisor        *float64 `json:"divisor"`
	B509Op59QuantityUnit           *string  `json:"unit"`
}

type B509Op59QuantityEntry struct {
	B509Op59QuantityKey                   B509Op59QuantityKey    `json:"key"`
	B509Op59QuantityGrade                 string                 `json:"grade"`
	B509Op59QuantityMetadataPresence      string                 `json:"metadata_presence"`
	B509Op59QuantityPhysicalQualification string                 `json:"physical_qualification"`
	B509Op59QuantityCodec                 *B509Op59QuantityCodec `json:"codec"`
	B509Op59QuantityPosition              *uint16                `json:"position"`
	B509Op59QuantityLabel                 *string                `json:"label"`
	B509Op59QuantityPreserveRaw           bool                   `json:"preserve_raw"`
	B509Op59QuantityUnknown               []string               `json:"unknown"`
}

func B509Op59QuantityEntries() ([]B509Op59QuantityEntry, error) {
	var b509Op59QuantityEntries []B509Op59QuantityEntry
	b509Op59QuantityDecoder := b509Op59QuantityJSON.NewDecoder(b509Op59QuantityBytes.NewBufferString(b509Op59QuantityData))
	b509Op59QuantityDecoder.DisallowUnknownFields()
	if b509Op59QuantityError := b509Op59QuantityDecoder.Decode(&b509Op59QuantityEntries); b509Op59QuantityError != nil {
		return nil, b509Op59QuantityError
	}
	if b509Op59QuantityError := b509Op59QuantityDecoder.Decode(new(any)); b509Op59QuantityError != b509Op59QuantityIO.EOF {
		return nil, b509Op59QuantityFmt.Errorf("ab14: trailing embedded data")
	}
	if len(b509Op59QuantityEntries) != 5 {
		return nil, b509Op59QuantityFmt.Errorf("ab14: expected five finite entries")
	}
	for b509Op59QuantityIndex, b509Op59QuantityEntry := range b509Op59QuantityEntries {
		if b509Op59QuantityEntry.B509Op59QuantityKey.B509Op59QuantityFamily == "" || b509Op59QuantityEntry.B509Op59QuantityKey.B509Op59QuantityVariant == "" || b509Op59QuantityEntry.B509Op59QuantityKey.B509Op59QuantityPrimary == 0 || b509Op59QuantityEntry.B509Op59QuantityKey.B509Op59QuantityField == "" || b509Op59QuantityEntry.B509Op59QuantityGrade != "Hypothesis" || b509Op59QuantityEntry.B509Op59QuantityMetadataPresence != "Proven" || b509Op59QuantityEntry.B509Op59QuantityPhysicalQualification != "Unknown" || !b509Op59QuantityEntry.B509Op59QuantityPreserveRaw || len(b509Op59QuantityEntry.B509Op59QuantityUnknown) == 0 {
			return nil, b509Op59QuantityFmt.Errorf("ab14: invalid entry %d", b509Op59QuantityIndex)
		}
		if b509Op59QuantityEntry.B509Op59QuantityCodec == nil || (b509Op59QuantityEntry.B509Op59QuantityCodec.B509Op59QuantityRepresentation != "raw" && b509Op59QuantityEntry.B509Op59QuantityCodec.B509Op59QuantityRepresentation != "signedLE16" && b509Op59QuantityEntry.B509Op59QuantityCodec.B509Op59QuantityRepresentation != "unsignedLE16") || (b509Op59QuantityEntry.B509Op59QuantityCodec.B509Op59QuantityDivisor != nil && *b509Op59QuantityEntry.B509Op59QuantityCodec.B509Op59QuantityDivisor <= 0) {
			return nil, b509Op59QuantityFmt.Errorf("ab14: invalid codec %d", b509Op59QuantityIndex)
		}
		b509Op59QuantityCodec := b509Op59QuantityEntry.B509Op59QuantityCodec
		if b509Op59QuantityCodec.B509Op59QuantityRepresentation == "raw" {
			if b509Op59QuantityCodec.B509Op59QuantityWidth != nil || b509Op59QuantityCodec.B509Op59QuantityDivisor != nil || b509Op59QuantityCodec.B509Op59QuantityUnit != nil {
				return nil, b509Op59QuantityFmt.Errorf("ab14: raw limit has unprovided metadata")
			}
		} else if b509Op59QuantityCodec.B509Op59QuantityWidth == nil || *b509Op59QuantityCodec.B509Op59QuantityWidth != 16 || b509Op59QuantityCodec.B509Op59QuantityDivisor == nil || b509Op59QuantityCodec.B509Op59QuantityUnit == nil || *b509Op59QuantityCodec.B509Op59QuantityUnit == "" {
			return nil, b509Op59QuantityFmt.Errorf("ab14: incomplete quantity codec")
		}
		for _, b509Op59QuantityPrior := range b509Op59QuantityEntries[:b509Op59QuantityIndex] {
			if b509Op59QuantityReflect.DeepEqual(b509Op59QuantityPrior.B509Op59QuantityKey, b509Op59QuantityEntry.B509Op59QuantityKey) {
				return nil, b509Op59QuantityFmt.Errorf("ab14: duplicate key")
			}
		}
	}
	return b509Op59QuantityEntries, nil
}

func B509Op59QuantityLookup(b509Op59QuantityKey B509Op59QuantityKey) (B509Op59QuantityEntry, bool, error) {
	b509Op59QuantityEntries, b509Op59QuantityError := B509Op59QuantityEntries()
	if b509Op59QuantityError != nil {
		return B509Op59QuantityEntry{}, false, b509Op59QuantityError
	}
	for _, b509Op59QuantityEntry := range b509Op59QuantityEntries {
		if b509Op59QuantityReflect.DeepEqual(b509Op59QuantityEntry.B509Op59QuantityKey, b509Op59QuantityKey) {
			return b509Op59QuantityEntry, true, nil
		}
	}
	return B509Op59QuantityEntry{}, false, nil
}
