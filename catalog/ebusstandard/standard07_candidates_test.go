package ebusstandard

import (
	standard07CandidateReflect "reflect"
	standard07CandidateTesting "testing"
)

func standard07CandidateByte(standard07CandidateValue byte) *byte { return &standard07CandidateValue }

func standard07CandidateIdentifierKey() Standard07CandidateKey {
	return Standard07CandidateKey{Standard07CandidateFamily: "standard", Standard07CandidateVariant: "identification", Standard07CandidatePrimary: 7, Standard07CandidateSecondary: standard07CandidateByte(4), Standard07CandidateDirection: "response", Standard07CandidateField: "identifier"}
}

func TestAb77ExactKey(standard07CandidateT *standard07CandidateTesting.T) {
	standard07CandidateControl := standard07CandidateIdentifierKey()
	standard07CandidateEntry, standard07CandidateOK := Standard07CandidateLookup(standard07CandidateControl)
	if !standard07CandidateOK || standard07CandidateEntry.Standard07CandidateLocalParser == nil {
		standard07CandidateT.Fatal("valid independent key missing")
	}
	for _, standard07CandidateAxis := range []string{"family", "variant", "primary", "secondary", "secondary_absent", "secondary_zero", "operation_zero", "direction", "field"} {
		standard07CandidateT.Run(standard07CandidateAxis, func(standard07CandidateT *standard07CandidateTesting.T) {
			standard07CandidateKey := standard07CandidateIdentifierKey()
			switch standard07CandidateAxis {
			case "family":
				standard07CandidateKey.Standard07CandidateFamily = ""
			case "variant":
				standard07CandidateKey.Standard07CandidateVariant = "existence_inquiry"
			case "primary":
				standard07CandidateKey.Standard07CandidatePrimary = 0
			case "secondary":
				standard07CandidateKey.Standard07CandidateSecondary = standard07CandidateByte(254)
			case "secondary_absent":
				standard07CandidateKey.Standard07CandidateSecondary = nil
			case "secondary_zero":
				standard07CandidateKey.Standard07CandidateSecondary = standard07CandidateByte(0)
			case "operation_zero":
				standard07CandidateKey.Standard07CandidateOperation = standard07CandidateByte(0)
			case "direction":
				standard07CandidateKey.Standard07CandidateDirection = "request"
			case "field":
				standard07CandidateKey.Standard07CandidateField = "missing"
			}
			if _, standard07CandidateOK := Standard07CandidateLookup(standard07CandidateKey); standard07CandidateOK {
				standard07CandidateT.Fatal("changed axis matched")
			}
			if _, standard07CandidateOK := Standard07CandidateLookup(standard07CandidateControl); !standard07CandidateOK {
				standard07CandidateT.Fatal("control lost")
			}
		})
	}
	if standard07CandidateEqualByte(nil, standard07CandidateByte(0)) || standard07CandidateEqualByte(standard07CandidateByte(0), nil) || !standard07CandidateEqualByte(standard07CandidateByte(0), standard07CandidateByte(0)) {
		standard07CandidateT.Fatal("absent/zero equality")
	}
}

func TestAb77Metadata(standard07CandidateT *standard07CandidateTesting.T) {
	standard07CandidateData := Standard07CandidateMetadata()
	if standard07CandidateData.Standard07CandidateUnknownGrade != "Unknown" || standard07CandidateData.Standard07CandidatePresenceGrade != "Proven" || len(standard07CandidateData.Standard07CandidateEntries) != 7 || len(standard07CandidateData.Standard07CandidateUnknown) != 5 || standard07CandidateData.Standard07CandidateScope == "" {
		standard07CandidateT.Fatal("catalog boundary")
	}
	for standard07CandidateIndex, standard07CandidateField := range []string{"identifier", "major", "minor", "product_number", "version_control", "revision_control", "existence_inquiry"} {
		standard07CandidateKey := standard07CandidateIdentifierKey()
		standard07CandidateKey.Standard07CandidateField = standard07CandidateField
		if standard07CandidateIndex == 6 {
			standard07CandidateKey.Standard07CandidateVariant = "existence_inquiry"
			standard07CandidateKey.Standard07CandidateSecondary = standard07CandidateByte(254)
			standard07CandidateKey.Standard07CandidateDirection = "request"
		}
		standard07CandidateEntry, standard07CandidateOK := Standard07CandidateLookup(standard07CandidateKey)
		if !standard07CandidateOK || !standard07CandidateReflect.DeepEqual(standard07CandidateEntry, standard07CandidateData.Standard07CandidateEntries[standard07CandidateIndex]) || standard07CandidateEntry.Standard07CandidateGrade != "Hypothesis" || standard07CandidateEntry.Standard07CandidateWidth != nil {
			standard07CandidateT.Fatalf("field %s", standard07CandidateField)
		}
		if standard07CandidateIndex == 4 || standard07CandidateIndex == 5 {
			if standard07CandidateEntry.Standard07CandidateLogicalPosition == nil || *standard07CandidateEntry.Standard07CandidateLogicalPosition != standard07CandidateIndex+4 || standard07CandidateEntry.Standard07CandidateControl == nil {
				standard07CandidateT.Fatal("logical overlap missing")
			}
			standard07CandidateControl := standard07CandidateEntry.Standard07CandidateControl
			if standard07CandidateControl.Standard07CandidateNotOkayMinimum != 0 || standard07CandidateControl.Standard07CandidateNotOkayMaximum != 98 || standard07CandidateControl.Standard07CandidateDefault != "not_okay" || standard07CandidateControl.Standard07CandidateOkayValues != nil {
				standard07CandidateT.Fatal("unsupported okay claim")
			}
		} else if standard07CandidateEntry.Standard07CandidateLogicalPosition != nil || standard07CandidateEntry.Standard07CandidateControl != nil {
			standard07CandidateT.Fatal("invented native coordinates")
		}
		if standard07CandidateIndex == 6 {
			if standard07CandidateEntry.Standard07CandidateDirectReplyRequired == nil || *standard07CandidateEntry.Standard07CandidateDirectReplyRequired {
				standard07CandidateT.Fatal("existence direct reply required")
			}
		} else if standard07CandidateEntry.Standard07CandidateDirectReplyRequired != nil {
			standard07CandidateT.Fatal("invented answer policy")
		}
		if (standard07CandidateEntry.Standard07CandidateLocalParser != nil) != (standard07CandidateIndex == 0) {
			standard07CandidateT.Fatal("local parser attached to wrong field")
		}
	}
}

func TestAb77FreshCopies(standard07CandidateT *standard07CandidateTesting.T) {
	standard07CandidateBefore := Standard07CandidateMetadata()
	standard07CandidateChanged := Standard07CandidateMetadata()
	standard07CandidateChanged.Standard07CandidateUnknown[0] = "changed"
	*standard07CandidateChanged.Standard07CandidateEntries[0].Standard07CandidateKey.Standard07CandidateSecondary = 0
	standard07CandidateChanged.Standard07CandidateEntries[0].Standard07CandidateLocalParser.Standard07CandidateNumericOffsets[0] = 0
	standard07CandidateChanged.Standard07CandidateEntries[0].Standard07CandidateLocalParser.Standard07CandidateProductOffsets[0] = 0
	*standard07CandidateChanged.Standard07CandidateEntries[4].Standard07CandidateLogicalPosition = 0
	standard07CandidateChanged.Standard07CandidateEntries[4].Standard07CandidateControl.Standard07CandidateDefault = "changed"
	*standard07CandidateChanged.Standard07CandidateEntries[6].Standard07CandidateDirectReplyRequired = true
	standard07CandidateChanged.Standard07CandidateEntries[0].Standard07CandidateKey.Standard07CandidateFamily = "changed"
	if !standard07CandidateReflect.DeepEqual(standard07CandidateBefore, Standard07CandidateMetadata()) {
		standard07CandidateT.Fatal("shared accessor storage")
	}
	standard07CandidateLookup, standard07CandidateOK := Standard07CandidateLookup(standard07CandidateIdentifierKey())
	if !standard07CandidateOK {
		standard07CandidateT.Fatal("missing lookup")
	}
	standard07CandidateLookup.Standard07CandidateLocalParser.Standard07CandidateNumericOffsets[0] = 0
	*standard07CandidateLookup.Standard07CandidateKey.Standard07CandidateSecondary = 0
	standard07CandidateAgain, standard07CandidateOK := Standard07CandidateLookup(standard07CandidateIdentifierKey())
	if !standard07CandidateOK || !standard07CandidateReflect.DeepEqual(standard07CandidateAgain, standard07CandidateBefore.Standard07CandidateEntries[0]) {
		standard07CandidateT.Fatal("shared lookup storage")
	}
}

func standard07CandidateEvaluate(standard07CandidateParser Standard07CandidateLocalParser, standard07CandidateBytes []byte, standard07CandidateCount int) (string, int, int, int, bool) {
	if standard07CandidateCount < standard07CandidateParser.Standard07CandidateMinimumDeclaredCount || len(standard07CandidateBytes) < standard07CandidateParser.Standard07CandidateMinimumAvailableCount || (standard07CandidateParser.Standard07CandidateMaximumDeclaredCount != nil && standard07CandidateCount > *standard07CandidateParser.Standard07CandidateMaximumDeclaredCount) {
		return "", 0, 0, 0, false
	}
	for standard07CandidateIndex := 0; standard07CandidateIndex < standard07CandidateParser.Standard07CandidateIdentifierWidth; standard07CandidateIndex++ {
		standard07CandidateByte := standard07CandidateBytes[standard07CandidateParser.Standard07CandidateIdentifierOffset+standard07CandidateIndex]
		if standard07CandidateIndex == 0 && standard07CandidateParser.Standard07CandidateFirstIdentifierNonzero && standard07CandidateByte == 0 {
			return "", 0, 0, 0, false
		}
		if (standard07CandidateByte != 0 || !standard07CandidateParser.Standard07CandidateCharacterZeroAllowed) && (standard07CandidateByte < standard07CandidateParser.Standard07CandidateCharacterMinimum || standard07CandidateByte > standard07CandidateParser.Standard07CandidateCharacterMaximum) {
			return "", 0, 0, 0, false
		}
	}
	for _, standard07CandidateOffset := range standard07CandidateParser.Standard07CandidateNumericOffsets {
		standard07CandidateByte := standard07CandidateBytes[standard07CandidateOffset]
		if standard07CandidateByte&15 > standard07CandidateParser.Standard07CandidateLowNibbleMaximum || (standard07CandidateParser.Standard07CandidateHighNibbleMaximum != nil && standard07CandidateByte>>4 > *standard07CandidateParser.Standard07CandidateHighNibbleMaximum) {
			return "", 0, 0, 0, false
		}
	}
	standard07CandidateDecimal := func(standard07CandidateByte byte) int {
		return standard07CandidateParser.Standard07CandidateDecimalMultiplier*int(standard07CandidateByte>>4) + int(standard07CandidateByte&15)
	}
	return string(standard07CandidateBytes[standard07CandidateParser.Standard07CandidateIdentifierOffset : standard07CandidateParser.Standard07CandidateIdentifierOffset+standard07CandidateParser.Standard07CandidateIdentifierWidth]), standard07CandidateDecimal(standard07CandidateBytes[standard07CandidateParser.Standard07CandidateMajorOffset]), standard07CandidateDecimal(standard07CandidateBytes[standard07CandidateParser.Standard07CandidateMinorOffset]), standard07CandidateDecimal(standard07CandidateBytes[standard07CandidateParser.Standard07CandidateProductOffsets[0]]) + standard07CandidateParser.Standard07CandidateProductMultiplier*standard07CandidateDecimal(standard07CandidateBytes[standard07CandidateParser.Standard07CandidateProductOffsets[1]]), true
}

func TestAb77LocalBoundaries(standard07CandidateT *standard07CandidateTesting.T) {
	standard07CandidateEntry, standard07CandidateOK := Standard07CandidateLookup(standard07CandidateIdentifierKey())
	if !standard07CandidateOK || standard07CandidateEntry.Standard07CandidateLocalParser == nil {
		standard07CandidateT.Fatal("missing parser metadata")
	}
	standard07CandidateParser := *standard07CandidateEntry.Standard07CandidateLocalParser
	standard07CandidateExpected := Standard07CandidateLocalParser{Standard07CandidateMinimumDeclaredCount: 10, Standard07CandidateMinimumAvailableCount: 10, Standard07CandidateIdentifierOffset: 1, Standard07CandidateIdentifierWidth: 5, Standard07CandidateFirstIdentifierNonzero: true, Standard07CandidateCharacterZeroAllowed: true, Standard07CandidateCharacterMinimum: 32, Standard07CandidateCharacterMaximum: 126, Standard07CandidateNumericOffsets: []int{6, 7, 8, 9}, Standard07CandidateLowNibbleMaximum: 9, Standard07CandidateDecimalMultiplier: 10, Standard07CandidateMajorOffset: 6, Standard07CandidateMinorOffset: 7, Standard07CandidateProductOffsets: []int{8, 9}, Standard07CandidateProductMultiplier: 100, Standard07CandidateGrade: "Hypothesis", Standard07CandidateScope: "partial local parser only; not physical BCD qualification"}
	if !standard07CandidateReflect.DeepEqual(standard07CandidateParser, standard07CandidateExpected) {
		standard07CandidateT.Fatal("local predicate differs from supplied boundary")
	}
	for _, standard07CandidateCase := range []string{"valid", "count9", "count11", "short", "first_zero", "internal_zero", "low_character", "high_character", "FA", "byte_zero"} {
		standard07CandidateT.Run(standard07CandidateCase, func(standard07CandidateT *standard07CandidateTesting.T) {
			standard07CandidateBytes := []byte{0xA1, 0x41, 0x42, 0x43, 0x44, 0x45, 0xF9, 0x99, 0x17, 0x04}
			standard07CandidateCount := 10
			standard07CandidateWantOK := true
			standard07CandidateWantID := "ABCDE"
			switch standard07CandidateCase {
			case "count9":
				standard07CandidateCount = 9
				standard07CandidateWantOK = false
			case "count11":
				standard07CandidateCount = 11
			case "short":
				standard07CandidateBytes = standard07CandidateBytes[:9]
				standard07CandidateCount = 11
				standard07CandidateWantOK = false
			case "first_zero":
				standard07CandidateBytes[1] = 0
				standard07CandidateWantOK = false
			case "internal_zero":
				standard07CandidateBytes[2] = 0
				standard07CandidateWantID = "A\x00CDE"
			case "low_character":
				standard07CandidateBytes[3] = 0x1F
				standard07CandidateWantOK = false
			case "high_character":
				standard07CandidateBytes[3] = 0x7F
				standard07CandidateWantOK = false
			case "FA":
				standard07CandidateBytes[6] = 0xFA
				standard07CandidateWantOK = false
			case "byte_zero":
				standard07CandidateBytes[0] = 0
			}
			standard07CandidateID, standard07CandidateMajor, standard07CandidateMinor, standard07CandidateProduct, standard07CandidateOK := standard07CandidateEvaluate(standard07CandidateParser, standard07CandidateBytes, standard07CandidateCount)
			if standard07CandidateOK != standard07CandidateWantOK {
				standard07CandidateT.Fatal("admission")
			}
			if standard07CandidateOK && (standard07CandidateID != standard07CandidateWantID || standard07CandidateMajor != 159 || standard07CandidateMinor != 99 || standard07CandidateProduct != 417) {
				standard07CandidateT.Fatalf("result %q %d %d %d", standard07CandidateID, standard07CandidateMajor, standard07CandidateMinor, standard07CandidateProduct)
			}
		})
	}
	for _, standard07CandidateOffset := range []int{6, 7, 8, 9} {
		standard07CandidateBytes := []byte{0xA1, 0x41, 0x42, 0x43, 0x44, 0x45, 0xF9, 0x99, 0x17, 0x04}
		standard07CandidateBytes[standard07CandidateOffset] = 0xFA
		if _, _, _, _, standard07CandidateOK := standard07CandidateEvaluate(standard07CandidateParser, standard07CandidateBytes, 10); standard07CandidateOK {
			standard07CandidateT.Fatalf("low nibble accepted at %d", standard07CandidateOffset)
		}
	}
	for standard07CandidateLength := 0; standard07CandidateLength < 10; standard07CandidateLength++ {
		if _, _, _, _, standard07CandidateOK := standard07CandidateEvaluate(standard07CandidateParser, make([]byte, standard07CandidateLength), 10); standard07CandidateOK {
			standard07CandidateT.Fatal("short buffer admitted")
		}
	}
}
