package vaillant

import (
	b506CandidateBytes "bytes"
	b506CandidateReflect "reflect"
	b506CandidateTesting "testing"
)

func b506CandidateByte(b506CandidateValue byte) *byte       { return &b506CandidateValue }
func b506CandidateString(b506CandidateValue string) *string { return &b506CandidateValue }

func b506CandidateKey(b506CandidateVariant string, b506CandidateOperation *byte, b506CandidateField string) B506CandidateKey {
	return B506CandidateKey{B506CandidateFamily: "vaillant", B506CandidateVariant: b506CandidateVariant, B506CandidatePrimary: 0xB5, B506CandidateSecondary: b506CandidateByte(6), B506CandidateOperation: b506CandidateOperation, B506CandidateField: b506CandidateField}
}

func b506CandidateEntry(b506CandidateT *b506CandidateTesting.T, b506CandidateKeyValue B506CandidateKey) B506CandidateEntry {
	b506CandidateT.Helper()
	b506CandidateEntryValue, b506CandidateOK := B506CandidateLookup(b506CandidateKeyValue)
	if !b506CandidateOK {
		b506CandidateT.Fatalf("missing declaration: %+v", b506CandidateKeyValue)
	}
	return b506CandidateEntryValue
}

func TestAb57ExactKeys(b506CandidateT *b506CandidateTesting.T) {
	b506CandidateKeys := []B506CandidateKey{
		b506CandidateKey("production-test", b506CandidateByte(0x11), "parameter"),
		b506CandidateKey("backlight-only", b506CandidateByte(2), "outputs"),
		b506CandidateKey("three-output", b506CandidateByte(2), "outputs"),
		b506CandidateKey("eeprom-checksum", nil, "checksum_type"),
	}
	b506CandidateCatalog := B506Candidates()
	if len(b506CandidateCatalog.B506CandidateEntries) != len(b506CandidateKeys) {
		b506CandidateT.Fatal("finite entry count")
	}
	for _, b506CandidateControl := range b506CandidateKeys {
		b506CandidateFound := b506CandidateEntry(b506CandidateT, b506CandidateControl)
		if !b506CandidateReflect.DeepEqual(b506CandidateFound.B506CandidateKey, b506CandidateControl) {
			b506CandidateT.Fatal("key changed")
		}
		b506CandidateMutations := map[string]func(*B506CandidateKey){
			"ab57_family_empty":     func(b506CandidateK *B506CandidateKey) { b506CandidateK.B506CandidateFamily = "" },
			"ab57_variant_empty":    func(b506CandidateK *B506CandidateKey) { b506CandidateK.B506CandidateVariant = "" },
			"ab57_field_empty":      func(b506CandidateK *B506CandidateKey) { b506CandidateK.B506CandidateField = "" },
			"ab57_family":           func(b506CandidateK *B506CandidateKey) { b506CandidateK.B506CandidateFamily = "ab57_other" },
			"ab57_variant":          func(b506CandidateK *B506CandidateKey) { b506CandidateK.B506CandidateVariant = "ab57_other" },
			"ab57_primary":          func(b506CandidateK *B506CandidateKey) { b506CandidateK.B506CandidatePrimary = 0 },
			"ab57_secondary_absent": func(b506CandidateK *B506CandidateKey) { b506CandidateK.B506CandidateSecondary = nil },
			"ab57_secondary_zero":   func(b506CandidateK *B506CandidateKey) { b506CandidateK.B506CandidateSecondary = b506CandidateByte(0) },
			"ab57_operation_zero":   func(b506CandidateK *B506CandidateKey) { b506CandidateK.B506CandidateOperation = b506CandidateByte(0) },
			"ab57_direction": func(b506CandidateK *B506CandidateKey) {
				b506CandidateK.B506CandidateDirection = b506CandidateString("request")
			},
			"ab57_direction_empty": func(b506CandidateK *B506CandidateKey) {
				b506CandidateK.B506CandidateDirection = b506CandidateString("")
			},
			"ab57_field": func(b506CandidateK *B506CandidateKey) { b506CandidateK.B506CandidateField = "ab57_other" },
		}
		if b506CandidateControl.B506CandidateOperation != nil {
			b506CandidateMutations["ab57_operation_absent"] = func(b506CandidateK *B506CandidateKey) { b506CandidateK.B506CandidateOperation = nil }
		}
		for b506CandidateName, b506CandidateMutate := range b506CandidateMutations {
			b506CandidateT.Run("ab57_"+b506CandidateControl.B506CandidateVariant+"/"+b506CandidateName, func(b506CandidateT *b506CandidateTesting.T) {
				b506CandidateChanged := b506CandidateControl
				b506CandidateMutate(&b506CandidateChanged)
				if _, b506CandidateOK := B506CandidateLookup(b506CandidateChanged); b506CandidateOK {
					b506CandidateT.Fatal("changed axis matched")
				}
				b506CandidateEntry(b506CandidateT, b506CandidateControl)
			})
		}
	}
}

func TestAb57UnionPermutation(b506CandidateT *b506CandidateTesting.T) {
	b506CandidateProduction := b506CandidateEntry(b506CandidateT, b506CandidateKey("production-test", b506CandidateByte(0x11), "parameter"))
	b506CandidateExpected := map[byte]string{0: "opaque-through-end", 1: "opaque-through-end", 2: "opaque-through-end", 4: "opaque-through-end", 5: "opaque-through-end", 3: "raw-byte", 6: "raw-byte", 7: "raw-byte", 8: "raw-byte", 9: "raw-byte", 10: "boolean", 32: "raw-byte", 33: "raw-byte", 34: "raw-byte", 35: "raw-byte", 36: "raw-byte", 37: "raw-byte"}
	if len(b506CandidateProduction.B506CandidateParameters) != len(b506CandidateExpected) {
		b506CandidateT.Fatal("union size")
	}
	for b506CandidateRotation := 0; b506CandidateRotation < len(b506CandidateExpected); b506CandidateRotation++ {
		for b506CandidateSelector := 0; b506CandidateSelector < 256; b506CandidateSelector++ {
			b506CandidateParameter, b506CandidateOK := B506CandidateParameterLookup(b506CandidateProduction, byte(b506CandidateSelector))
			b506CandidateWant, b506CandidateMapped := b506CandidateExpected[byte(b506CandidateSelector)]
			if b506CandidateOK != b506CandidateMapped {
				b506CandidateT.Fatalf("selector %d mapping", b506CandidateSelector)
			}
			if !b506CandidateOK {
				continue
			}
			if b506CandidateParameter.B506CandidateSelector != byte(b506CandidateSelector) || b506CandidateParameter.B506CandidateRepresentation != b506CandidateWant || b506CandidateParameter.B506CandidateMetadata.SchemaCodec != b506CandidateWant || b506CandidateParameter.B506CandidateMetadata.SchemaGrade != SchemaHypothesis || b506CandidateParameter.B506CandidateMetadata.SchemaPosition != nil || b506CandidateParameter.B506CandidateMinimumBytes != nil || b506CandidateParameter.B506CandidateMaximumBytes != nil {
				b506CandidateT.Fatalf("selector %d declaration", b506CandidateSelector)
			}
			if b506CandidateWant == "opaque-through-end" {
				if b506CandidateParameter.B506CandidateMetadata.SchemaWidth != nil {
					b506CandidateT.Fatal("opaque width invented")
				}
			} else if b506CandidateParameter.B506CandidateMetadata.SchemaWidth == nil || *b506CandidateParameter.B506CandidateMetadata.SchemaWidth != 1 {
				b506CandidateT.Fatal("one byte width missing")
			}
		}
		b506CandidateProduction.B506CandidateParameters = append(b506CandidateProduction.B506CandidateParameters[1:], b506CandidateProduction.B506CandidateParameters[0])
	}
	b506CandidateBoolean, _ := B506CandidateParameterLookup(b506CandidateProduction, 10)
	if !b506CandidateReflect.DeepEqual(b506CandidateBoolean.B506CandidateMetadata.SchemaEnum, map[int64]string{0: "no", 1: "yes"}) {
		b506CandidateT.Fatal("Boolean meaning must cover exactly 0 and 1")
	}
	for _, b506CandidateSelector := range []byte{0, 1, 2, 4, 5} {
		b506CandidateOpaque, _ := B506CandidateParameterLookup(b506CandidateProduction, b506CandidateSelector)
		for _, b506CandidateRaw := range [][]byte{{0, 0xAA, 0xFF}, {}} {
			b506CandidateRetained := append([]byte{}, b506CandidateRaw...)
			if b506CandidateOpaque.B506CandidateRepresentation != "opaque-through-end" || !b506CandidateBytes.Equal(b506CandidateRetained, b506CandidateRaw) {
				b506CandidateT.Fatal("opaque representation")
			}
		}
	}
}

func TestAb57OutputAlternatives(b506CandidateT *b506CandidateTesting.T) {
	for _, b506CandidateCase := range []struct {
		B506CandidateVariant     string
		B506CandidateMask        byte
		B506CandidateReserved    byte
		B506CandidateActive      byte
		B506CandidateRawReserved byte
	}{
		{"backlight-only", 1, 254, 1, 6}, {"three-output", 7, 248, 7, 0},
	} {
		b506CandidateValue := b506CandidateEntry(b506CandidateT, b506CandidateKey(b506CandidateCase.B506CandidateVariant, b506CandidateByte(2), "outputs"))
		if b506CandidateValue.B506CandidateOutputMask == nil || b506CandidateValue.B506CandidateReservedMask == nil || *b506CandidateValue.B506CandidateOutputMask != b506CandidateCase.B506CandidateMask || *b506CandidateValue.B506CandidateReservedMask != b506CandidateCase.B506CandidateReserved {
			b506CandidateT.Fatal("variant masks")
		}
		b506CandidateRaw := byte(7)
		if b506CandidateRaw&*b506CandidateValue.B506CandidateOutputMask != b506CandidateCase.B506CandidateActive || b506CandidateRaw&*b506CandidateValue.B506CandidateReservedMask != b506CandidateCase.B506CandidateRawReserved || (b506CandidateRaw&*b506CandidateValue.B506CandidateOutputMask)|(b506CandidateRaw&*b506CandidateValue.B506CandidateReservedMask) != b506CandidateRaw {
			b506CandidateT.Fatal("raw07 alternative or reserved bits lost")
		}
	}
}

func TestAb57ChecksumScalar(b506CandidateT *b506CandidateTesting.T) {
	b506CandidateValue := b506CandidateEntry(b506CandidateT, b506CandidateKey("eeprom-checksum", nil, "checksum_type"))
	if b506CandidateValue.B506CandidateChecksumInterpretation != "scalar, not flags" || len(b506CandidateValue.B506CandidateChecksumValues) != 128 || b506CandidateValue.B506CandidateChecksumWidth != nil || b506CandidateValue.B506CandidateAlgorithm != SchemaUnknown || b506CandidateValue.B506CandidateCompletion != SchemaUnknown || b506CandidateValue.B506CandidateChecksumData != SchemaUnknown {
		b506CandidateT.Fatal("checksum properties conflated")
	}
	for b506CandidateRaw := 0; b506CandidateRaw < 256; b506CandidateRaw++ {
		b506CandidateLabel, b506CandidateOK := b506CandidateValue.B506CandidateChecksumValues[byte(b506CandidateRaw)]
		b506CandidateWant := "scalar-type"
		if b506CandidateRaw == 0 {
			b506CandidateWant = "query-supported-types"
		}
		if b506CandidateOK != (b506CandidateRaw < 128) || (b506CandidateOK && b506CandidateLabel != b506CandidateWant) {
			b506CandidateT.Fatalf("checksum %02X must remain one scalar declaration or unmapped", b506CandidateRaw)
		}
	}
}

func TestAb57FreshCopies(b506CandidateT *b506CandidateTesting.T) {
	b506CandidateBefore := B506Candidates()
	b506CandidateChanged := B506Candidates()
	b506CandidateChanged.B506CandidateLimitations[0] = "ab57_changed"
	*b506CandidateChanged.B506CandidateEntries[0].B506CandidateKey.B506CandidateSecondary = 0
	*b506CandidateChanged.B506CandidateEntries[0].B506CandidateKey.B506CandidateOperation = 0
	*b506CandidateChanged.B506CandidateEntries[0].B506CandidateParameters[5].B506CandidateMetadata.SchemaWidth = 99
	b506CandidateChanged.B506CandidateEntries[0].B506CandidateParameters[16].B506CandidateMetadata.SchemaEnum[0] = "ab57_changed"
	*b506CandidateChanged.B506CandidateEntries[1].B506CandidateOutputMask = 255
	b506CandidateChanged.B506CandidateEntries[3].B506CandidateChecksumValues[3] = "flags"
	if !b506CandidateReflect.DeepEqual(B506Candidates(), b506CandidateBefore) {
		b506CandidateT.Fatal("accessor shared nested storage")
	}
	b506CandidateLookup := b506CandidateEntry(b506CandidateT, b506CandidateKey("production-test", b506CandidateByte(0x11), "parameter"))
	*b506CandidateLookup.B506CandidateKey.B506CandidateSecondary = 0
	b506CandidateLookup.B506CandidateParameters[16].B506CandidateMetadata.SchemaEnum[1] = "ab57_changed"
	if !b506CandidateReflect.DeepEqual(b506CandidateEntry(b506CandidateT, b506CandidateKey("production-test", b506CandidateByte(0x11), "parameter")), b506CandidateBefore.B506CandidateEntries[0]) {
		b506CandidateT.Fatal("lookup shared storage")
	}
}

func TestAb57ClaimBoundaries(b506CandidateT *b506CandidateTesting.T) {
	b506CandidateCatalog := B506Candidates()
	if b506CandidateCatalog.B506CandidatePresence != SchemaProven || b506CandidateCatalog.B506CandidateQualification != SchemaUnknown || len(b506CandidateCatalog.B506CandidateLimitations) != 6 {
		b506CandidateT.Fatal("claim boundaries")
	}
	for _, b506CandidateValue := range b506CandidateCatalog.B506CandidateEntries {
		if b506CandidateValue.B506CandidateGrade != SchemaHypothesis || b506CandidateValue.B506CandidateKey.B506CandidateDirection != nil || b506CandidateValue.B506CandidateChecksumWidth != nil || b506CandidateValue.B506CandidateAlgorithm != SchemaUnknown || b506CandidateValue.B506CandidateCompletion != SchemaUnknown || b506CandidateValue.B506CandidateChecksumData != SchemaUnknown {
			b506CandidateT.Fatal("unknown properties promoted")
		}
	}
}
