package vaillant

import (
	b512DHWStepperBytesPackage "bytes"
	_ "embed"
	b512DHWStepperJSON "encoding/json"
	b512DHWStepperFmt "fmt"
	b512DHWStepperIO "io"
	b512DHWStepperMaps "maps"
	b512DHWStepperReflect "reflect"
	b512DHWStepperSlices "slices"
	b512DHWStepperSync "sync"
)

//go:embed data/ab97_b512_dhw_stepper.json
var b512DHWStepperData []byte

type B512DHWStepperKey struct {
	B512DHWStepperFamily            string  `json:"family"`
	B512DHWStepperVariant           string  `json:"variant"`
	B512DHWStepperPrimary           uint8   `json:"primary"`
	B512DHWStepperSecondary         *uint8  `json:"secondary"`
	B512DHWStepperOperation         *uint8  `json:"operation"`
	B512DHWStepperDirection         string  `json:"direction"`
	B512DHWStepperField             string  `json:"field"`
	B512DHWStepperGroup             *uint8  `json:"group"`
	B512DHWStepperRequestedInstance *uint8  `json:"requested_instance"`
	B512DHWStepperRegister          *uint16 `json:"register"`
	B512DHWStepperCommandRole       *string `json:"command_role"`
}

type B512DHWStepperEntry struct {
	B512DHWStepperCoordinateBasis *string           `json:"coordinate_basis"`
	B512DHWStepperKey             B512DHWStepperKey `json:"key"`
	B512DHWStepperPosition        *uint8            `json:"position"`
	B512DHWStepperWidthBits       *uint8            `json:"width_bits"`
	B512DHWStepperBit             *uint8            `json:"bit"`
	B512DHWStepperDivisor         *float64          `json:"divisor"`
	B512DHWStepperUnit            *string           `json:"unit"`
	B512DHWStepperFormalUnit      *string           `json:"formal_unit"`
	B512DHWStepperBooleanDefault  *bool             `json:"boolean_default"`
	B512DHWStepperDefault         *int              `json:"default"`
	B512DHWStepperValues          map[string]string `json:"values"`
	B512DHWStepperLabel           *string           `json:"label"`
	B512DHWStepperLabelGrade      string            `json:"label_grade"`
	B512DHWStepperSigned          *bool             `json:"signed"`
	B512DHWStepperByteOrder       *string           `json:"byte_order"`
	B512DHWStepperReserved        bool              `json:"reserved"`
	B512DHWStepperGrade           string            `json:"grade"`
	B512DHWStepperQualification   string            `json:"qualification"`
	B512DHWStepperRawPolicy       string            `json:"raw_policy"`
	B512DHWStepperUnknown         []string          `json:"unknown"`
	B512DHWStepperNativeContract  string            `json:"native_contract"`
}

// The embedded catalog is validated once and never exposed to callers.
var b512DHWStepperCatalog = b512DHWStepperSync.OnceValues(b512DHWStepperLoadEntries)

func B512DHWStepperEntries() ([]B512DHWStepperEntry, error) {
	b512DHWStepperEntries, b512DHWStepperErr := b512DHWStepperCatalog()
	if b512DHWStepperErr != nil {
		return nil, b512DHWStepperErr
	}
	b512DHWStepperResult := make([]B512DHWStepperEntry, len(b512DHWStepperEntries))
	for b512DHWStepperIndex, b512DHWStepperEntry := range b512DHWStepperEntries {
		b512DHWStepperResult[b512DHWStepperIndex] = b512DHWStepperCloneEntry(b512DHWStepperEntry)
	}
	return b512DHWStepperResult, nil
}

func b512DHWStepperLoadEntries() ([]B512DHWStepperEntry, error) {
	b512DHWStepperDecoder := b512DHWStepperJSON.NewDecoder(b512DHWStepperBytesPackage.NewReader(b512DHWStepperData))
	b512DHWStepperDecoder.DisallowUnknownFields()
	var b512DHWStepperEntries []B512DHWStepperEntry
	if b512DHWStepperErr := b512DHWStepperDecoder.Decode(&b512DHWStepperEntries); b512DHWStepperErr != nil {
		return nil, b512DHWStepperErr
	}
	var b512DHWStepperTail any
	if b512DHWStepperErr := b512DHWStepperDecoder.Decode(&b512DHWStepperTail); b512DHWStepperErr != b512DHWStepperIO.EOF {
		return nil, b512DHWStepperFmt.Errorf("ab97 trailing data")
	}
	if len(b512DHWStepperEntries) == 0 {
		return nil, b512DHWStepperFmt.Errorf("ab97 empty catalog")
	}
	for b512DHWStepperIndex, b512DHWStepperEntry := range b512DHWStepperEntries {
		if b512DHWStepperEntry.B512DHWStepperKey.B512DHWStepperFamily == "" || b512DHWStepperEntry.B512DHWStepperKey.B512DHWStepperVariant == "" || b512DHWStepperEntry.B512DHWStepperKey.B512DHWStepperField == "" || (b512DHWStepperEntry.B512DHWStepperKey.B512DHWStepperDirection != "request" && b512DHWStepperEntry.B512DHWStepperKey.B512DHWStepperDirection != "positive") || b512DHWStepperEntry.B512DHWStepperGrade != "Hypothesis" || b512DHWStepperEntry.B512DHWStepperQualification != "Unknown" || b512DHWStepperEntry.B512DHWStepperRawPolicy == "" || b512DHWStepperEntry.B512DHWStepperNativeContract == "" || len(b512DHWStepperEntry.B512DHWStepperUnknown) == 0 {
			return nil, b512DHWStepperFmt.Errorf("ab97 invalid entry %d", b512DHWStepperIndex)
		}
		if b512DHWStepperEntry.B512DHWStepperDivisor != nil && *b512DHWStepperEntry.B512DHWStepperDivisor <= 0 {
			return nil, b512DHWStepperFmt.Errorf("ab97 invalid divisor")
		}
		if b512DHWStepperEntry.B512DHWStepperBit != nil && (*b512DHWStepperEntry.B512DHWStepperBit > 7 || b512DHWStepperEntry.B512DHWStepperWidthBits == nil || *b512DHWStepperEntry.B512DHWStepperWidthBits != 1) {
			return nil, b512DHWStepperFmt.Errorf("ab97 invalid bit")
		}
		for b512DHWStepperPrevious := 0; b512DHWStepperPrevious < b512DHWStepperIndex; b512DHWStepperPrevious++ {
			if b512DHWStepperReflect.DeepEqual(b512DHWStepperEntry.B512DHWStepperKey, b512DHWStepperEntries[b512DHWStepperPrevious].B512DHWStepperKey) {
				return nil, b512DHWStepperFmt.Errorf("ab97 duplicate key")
			}
		}
	}
	return b512DHWStepperEntries, nil
}

func B512DHWStepperLookup(b512DHWStepperKey B512DHWStepperKey) (B512DHWStepperEntry, bool, error) {
	b512DHWStepperEntries, b512DHWStepperErr := b512DHWStepperCatalog()
	if b512DHWStepperErr != nil {
		return B512DHWStepperEntry{}, false, b512DHWStepperErr
	}
	for _, b512DHWStepperEntry := range b512DHWStepperEntries {
		if b512DHWStepperReflect.DeepEqual(b512DHWStepperKey, b512DHWStepperEntry.B512DHWStepperKey) {
			return b512DHWStepperCloneEntry(b512DHWStepperEntry), true, nil
		}
	}
	return B512DHWStepperEntry{}, false, nil
}

func b512DHWStepperClonePointer[T any](b512DHWStepperValue *T) *T {
	if b512DHWStepperValue == nil {
		return nil
	}
	b512DHWStepperCopy := *b512DHWStepperValue
	return &b512DHWStepperCopy
}

func b512DHWStepperCloneEntry(b512DHWStepperEntry B512DHWStepperEntry) B512DHWStepperEntry {
	b512DHWStepperEntry.B512DHWStepperKey.B512DHWStepperSecondary = b512DHWStepperClonePointer(b512DHWStepperEntry.B512DHWStepperKey.B512DHWStepperSecondary)
	b512DHWStepperEntry.B512DHWStepperKey.B512DHWStepperOperation = b512DHWStepperClonePointer(b512DHWStepperEntry.B512DHWStepperKey.B512DHWStepperOperation)
	b512DHWStepperEntry.B512DHWStepperKey.B512DHWStepperGroup = b512DHWStepperClonePointer(b512DHWStepperEntry.B512DHWStepperKey.B512DHWStepperGroup)
	b512DHWStepperEntry.B512DHWStepperKey.B512DHWStepperRequestedInstance = b512DHWStepperClonePointer(b512DHWStepperEntry.B512DHWStepperKey.B512DHWStepperRequestedInstance)
	b512DHWStepperEntry.B512DHWStepperKey.B512DHWStepperRegister = b512DHWStepperClonePointer(b512DHWStepperEntry.B512DHWStepperKey.B512DHWStepperRegister)
	b512DHWStepperEntry.B512DHWStepperKey.B512DHWStepperCommandRole = b512DHWStepperClonePointer(b512DHWStepperEntry.B512DHWStepperKey.B512DHWStepperCommandRole)
	b512DHWStepperEntry.B512DHWStepperCoordinateBasis = b512DHWStepperClonePointer(b512DHWStepperEntry.B512DHWStepperCoordinateBasis)
	b512DHWStepperEntry.B512DHWStepperPosition = b512DHWStepperClonePointer(b512DHWStepperEntry.B512DHWStepperPosition)
	b512DHWStepperEntry.B512DHWStepperWidthBits = b512DHWStepperClonePointer(b512DHWStepperEntry.B512DHWStepperWidthBits)
	b512DHWStepperEntry.B512DHWStepperBit = b512DHWStepperClonePointer(b512DHWStepperEntry.B512DHWStepperBit)
	b512DHWStepperEntry.B512DHWStepperDivisor = b512DHWStepperClonePointer(b512DHWStepperEntry.B512DHWStepperDivisor)
	b512DHWStepperEntry.B512DHWStepperUnit = b512DHWStepperClonePointer(b512DHWStepperEntry.B512DHWStepperUnit)
	b512DHWStepperEntry.B512DHWStepperFormalUnit = b512DHWStepperClonePointer(b512DHWStepperEntry.B512DHWStepperFormalUnit)
	b512DHWStepperEntry.B512DHWStepperBooleanDefault = b512DHWStepperClonePointer(b512DHWStepperEntry.B512DHWStepperBooleanDefault)
	b512DHWStepperEntry.B512DHWStepperDefault = b512DHWStepperClonePointer(b512DHWStepperEntry.B512DHWStepperDefault)
	b512DHWStepperEntry.B512DHWStepperLabel = b512DHWStepperClonePointer(b512DHWStepperEntry.B512DHWStepperLabel)
	b512DHWStepperEntry.B512DHWStepperSigned = b512DHWStepperClonePointer(b512DHWStepperEntry.B512DHWStepperSigned)
	b512DHWStepperEntry.B512DHWStepperByteOrder = b512DHWStepperClonePointer(b512DHWStepperEntry.B512DHWStepperByteOrder)
	b512DHWStepperEntry.B512DHWStepperValues = b512DHWStepperMaps.Clone(b512DHWStepperEntry.B512DHWStepperValues)
	b512DHWStepperEntry.B512DHWStepperUnknown = b512DHWStepperSlices.Clone(b512DHWStepperEntry.B512DHWStepperUnknown)
	return b512DHWStepperEntry
}
