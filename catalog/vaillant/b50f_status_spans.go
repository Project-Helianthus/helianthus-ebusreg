package vaillant

import (
	b50FStatusSpanBytes "bytes"
	_ "embed"
	b50FStatusSpanJSON "encoding/json"
	b50FStatusSpanFmt "fmt"
	b50FStatusSpanIO "io"
)

//go:embed data/ab12_b50f_status_spans.json
var b50FStatusSpanStatusJSON string

type B50FStatusSpanKey struct {
	B50FStatusSpanFamily            string  `json:"family"`
	B50FStatusSpanVariant           string  `json:"variant"`
	B50FStatusSpanPrimary           byte    `json:"primary"`
	B50FStatusSpanSecondary         *byte   `json:"secondary"`
	B50FStatusSpanOperation         *byte   `json:"operation"`
	B50FStatusSpanDirection         *string `json:"direction"`
	B50FStatusSpanField             string  `json:"field"`
	B50FStatusSpanGroup             *byte   `json:"group"`
	B50FStatusSpanRequestedInstance *byte   `json:"requested_instance"`
	B50FStatusSpanRegister          *uint16 `json:"register"`
	B50FStatusSpanCommandRole       *string `json:"command_role"`
}

type B50FStatusSpanBit struct {
	B50FStatusSpanBit      byte    `json:"bit"`
	B50FStatusSpanLabel    *string `json:"label"`
	B50FStatusSpanCategory *string `json:"category"`
	B50FStatusSpanReserved bool    `json:"reserved"`
}

type B50FStatusSpanEntry struct {
	B50FStatusSpanKey             B50FStatusSpanKey   `json:"key"`
	B50FStatusSpanPosition        *int                `json:"position"`
	B50FStatusSpanWidthBits       *int                `json:"width_bits"`
	B50FStatusSpanGrade           string              `json:"grade"`
	B50FStatusSpanQualification   string              `json:"qualification"`
	B50FStatusSpanCoordinateGrade string              `json:"coordinate_grade"`
	B50FStatusSpanAllocatedMask   *byte               `json:"allocated_mask"`
	B50FStatusSpanReservedMask    *byte               `json:"reserved_mask"`
	B50FStatusSpanBits            []B50FStatusSpanBit `json:"bits"`
	B50FStatusSpanRawFallback     string              `json:"raw_fallback"`
	B50FStatusSpanUnknown         []string            `json:"unknown"`
}

func b50FStatusSpanOptionalEqual[b50FStatusSpanValue comparable](b50FStatusSpanLeft, b50FStatusSpanRight *b50FStatusSpanValue) bool {
	if b50FStatusSpanLeft == nil || b50FStatusSpanRight == nil {
		return b50FStatusSpanLeft == nil && b50FStatusSpanRight == nil
	}
	return *b50FStatusSpanLeft == *b50FStatusSpanRight
}

func b50FStatusSpanKeyEqual(b50FStatusSpanLeft, b50FStatusSpanRight B50FStatusSpanKey) bool {
	return b50FStatusSpanLeft.B50FStatusSpanFamily == b50FStatusSpanRight.B50FStatusSpanFamily &&
		b50FStatusSpanLeft.B50FStatusSpanVariant == b50FStatusSpanRight.B50FStatusSpanVariant &&
		b50FStatusSpanLeft.B50FStatusSpanPrimary == b50FStatusSpanRight.B50FStatusSpanPrimary &&
		b50FStatusSpanOptionalEqual(b50FStatusSpanLeft.B50FStatusSpanSecondary, b50FStatusSpanRight.B50FStatusSpanSecondary) &&
		b50FStatusSpanOptionalEqual(b50FStatusSpanLeft.B50FStatusSpanOperation, b50FStatusSpanRight.B50FStatusSpanOperation) &&
		b50FStatusSpanOptionalEqual(b50FStatusSpanLeft.B50FStatusSpanDirection, b50FStatusSpanRight.B50FStatusSpanDirection) &&
		b50FStatusSpanLeft.B50FStatusSpanField == b50FStatusSpanRight.B50FStatusSpanField &&
		b50FStatusSpanOptionalEqual(b50FStatusSpanLeft.B50FStatusSpanGroup, b50FStatusSpanRight.B50FStatusSpanGroup) &&
		b50FStatusSpanOptionalEqual(b50FStatusSpanLeft.B50FStatusSpanRequestedInstance, b50FStatusSpanRight.B50FStatusSpanRequestedInstance) &&
		b50FStatusSpanOptionalEqual(b50FStatusSpanLeft.B50FStatusSpanRegister, b50FStatusSpanRight.B50FStatusSpanRegister) &&
		b50FStatusSpanOptionalEqual(b50FStatusSpanLeft.B50FStatusSpanCommandRole, b50FStatusSpanRight.B50FStatusSpanCommandRole)
}

func B50FStatusSpanEntries() ([]B50FStatusSpanEntry, error) {
	var b50FStatusSpanEntries []B50FStatusSpanEntry
	b50FStatusSpanDecoder := b50FStatusSpanJSON.NewDecoder(b50FStatusSpanBytes.NewBufferString(b50FStatusSpanStatusJSON))
	b50FStatusSpanDecoder.DisallowUnknownFields()
	if b50FStatusSpanError := b50FStatusSpanDecoder.Decode(&b50FStatusSpanEntries); b50FStatusSpanError != nil {
		return nil, b50FStatusSpanError
	}
	if b50FStatusSpanError := b50FStatusSpanDecoder.Decode(new(any)); b50FStatusSpanError != b50FStatusSpanIO.EOF {
		return nil, b50FStatusSpanFmt.Errorf("ab12: trailing data")
	}
	if len(b50FStatusSpanEntries) != 8 {
		return nil, b50FStatusSpanFmt.Errorf("ab12: incomplete finite catalog")
	}
	for b50FStatusSpanIndex, b50FStatusSpanEntry := range b50FStatusSpanEntries {
		if b50FStatusSpanEntry.B50FStatusSpanKey.B50FStatusSpanPrimary == 0 || b50FStatusSpanEntry.B50FStatusSpanKey.B50FStatusSpanFamily == "" || b50FStatusSpanEntry.B50FStatusSpanKey.B50FStatusSpanVariant == "" || b50FStatusSpanEntry.B50FStatusSpanKey.B50FStatusSpanField == "" || b50FStatusSpanEntry.B50FStatusSpanWidthBits == nil || *b50FStatusSpanEntry.B50FStatusSpanWidthBits != 8 || b50FStatusSpanEntry.B50FStatusSpanGrade != "Hypothesis" || b50FStatusSpanEntry.B50FStatusSpanQualification != "Unknown" || b50FStatusSpanEntry.B50FStatusSpanCoordinateGrade != "Unknown" || b50FStatusSpanEntry.B50FStatusSpanPosition != nil || b50FStatusSpanEntry.B50FStatusSpanRawFallback == "" || len(b50FStatusSpanEntry.B50FStatusSpanUnknown) == 0 {
			return nil, b50FStatusSpanFmt.Errorf("ab12: invalid entry %d", b50FStatusSpanIndex)
		}
		for b50FStatusSpanPrevious := 0; b50FStatusSpanPrevious < b50FStatusSpanIndex; b50FStatusSpanPrevious++ {
			if b50FStatusSpanKeyEqual(b50FStatusSpanEntries[b50FStatusSpanPrevious].B50FStatusSpanKey, b50FStatusSpanEntry.B50FStatusSpanKey) {
				return nil, b50FStatusSpanFmt.Errorf("ab12: duplicate key")
			}
		}
		if (b50FStatusSpanIndex < 6) != (b50FStatusSpanEntry.B50FStatusSpanAllocatedMask != nil && b50FStatusSpanEntry.B50FStatusSpanReservedMask != nil) {
			return nil, b50FStatusSpanFmt.Errorf("ab12: invalid container/raw partition")
		}
		if b50FStatusSpanEntry.B50FStatusSpanAllocatedMask == nil || b50FStatusSpanEntry.B50FStatusSpanReservedMask == nil {
			if b50FStatusSpanEntry.B50FStatusSpanAllocatedMask != nil || b50FStatusSpanEntry.B50FStatusSpanReservedMask != nil || len(b50FStatusSpanEntry.B50FStatusSpanBits) != 0 {
				return nil, b50FStatusSpanFmt.Errorf("ab12: inconsistent raw byte")
			}
			continue
		}
		if len(b50FStatusSpanEntry.B50FStatusSpanBits) != 8 || *b50FStatusSpanEntry.B50FStatusSpanAllocatedMask&*b50FStatusSpanEntry.B50FStatusSpanReservedMask != 0 || *b50FStatusSpanEntry.B50FStatusSpanAllocatedMask|*b50FStatusSpanEntry.B50FStatusSpanReservedMask != 255 {
			return nil, b50FStatusSpanFmt.Errorf("ab12: incomplete bit partition")
		}
		for b50FStatusSpanIndex, b50FStatusSpanBit := range b50FStatusSpanEntry.B50FStatusSpanBits {
			if int(b50FStatusSpanBit.B50FStatusSpanBit) != b50FStatusSpanIndex || b50FStatusSpanBit.B50FStatusSpanReserved != (*b50FStatusSpanEntry.B50FStatusSpanReservedMask&(1<<b50FStatusSpanIndex) != 0) || (b50FStatusSpanBit.B50FStatusSpanReserved && (b50FStatusSpanBit.B50FStatusSpanLabel != nil || b50FStatusSpanBit.B50FStatusSpanCategory != nil)) {
				return nil, b50FStatusSpanFmt.Errorf("ab12: invalid bit")
			}
		}
	}
	return b50FStatusSpanEntries, nil
}

func B50FStatusSpanLookup(b50FStatusSpanKey B50FStatusSpanKey) (B50FStatusSpanEntry, bool, error) {
	b50FStatusSpanEntries, b50FStatusSpanError := B50FStatusSpanEntries()
	if b50FStatusSpanError != nil {
		return B50FStatusSpanEntry{}, false, b50FStatusSpanError
	}
	for _, b50FStatusSpanEntry := range b50FStatusSpanEntries {
		if b50FStatusSpanKeyEqual(b50FStatusSpanEntry.B50FStatusSpanKey, b50FStatusSpanKey) {
			return b50FStatusSpanEntry, true, nil
		}
	}
	return B50FStatusSpanEntry{}, false, nil
}
