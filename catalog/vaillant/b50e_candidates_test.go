package vaillant

import (
	b50ECandidateReflect "reflect"
	b50ECandidateTesting "testing"
)

func b50ECandidatePointer[b50ECandidateValue any](b50ECandidateInput b50ECandidateValue) *b50ECandidateValue {
	return &b50ECandidateInput
}

func b50ECandidateComparisonKey() B50ECandidateKey {
	return B50ECandidateKey{B50ECandidateFamily: "vaillant", B50ECandidateVariant: "B50E", B50ECandidatePrimary: 0xB5, B50ECandidateSecondary: b50ECandidatePointer(byte(0x0E)), B50ECandidateField: "rule_comparison"}
}

func b50ECandidateClearKey() B50ECandidateKey {
	return B50ECandidateKey{B50ECandidateFamily: "vaillant", B50ECandidateVariant: "B50E", B50ECandidatePrimary: 0xB5, B50ECandidateSecondary: b50ECandidatePointer(byte(0x0E)), B50ECandidateOperation: b50ECandidatePointer(byte(1)), B50ECandidateDirection: b50ECandidatePointer("response"), B50ECandidateField: "clear_rule_positive_index"}
}

func b50ECandidateRequireEntry(b50ECandidateT *b50ECandidateTesting.T, b50ECandidateKey B50ECandidateKey) B50ECandidateEntry {
	b50ECandidateT.Helper()
	b50ECandidateEntry, b50ECandidateFound, b50ECandidateError := B50ELookupCandidate(b50ECandidateKey)
	if b50ECandidateError != nil || !b50ECandidateFound {
		b50ECandidateT.Fatalf("lookup: found=%v error=%v", b50ECandidateFound, b50ECandidateError)
	}
	return b50ECandidateEntry
}

func TestAb61ProductionMetadata(b50ECandidateT *b50ECandidateTesting.T) {
	b50ECandidateCatalog, b50ECandidateError := B50ECandidates()
	if b50ECandidateError != nil {
		b50ECandidateT.Fatal(b50ECandidateError)
	}
	if b50ECandidateCatalog.B50ECandidateMetadataPresence != "Proven" || b50ECandidateCatalog.B50ECandidatePhysicalQualification != "Unknown" || len(b50ECandidateCatalog.B50ECandidateEntries) != 2 || len(b50ECandidateCatalog.B50ECandidateLimitations) != 6 {
		b50ECandidateT.Fatalf("unexpected bounded catalog: %+v", b50ECandidateCatalog)
	}
	for b50ECandidateIndex, b50ECandidateKey := range []B50ECandidateKey{b50ECandidateComparisonKey(), b50ECandidateClearKey()} {
		b50ECandidateEntry := b50ECandidateRequireEntry(b50ECandidateT, b50ECandidateKey)
		if !b50ECandidateReflect.DeepEqual(b50ECandidateEntry, b50ECandidateCatalog.B50ECandidateEntries[b50ECandidateIndex]) {
			b50ECandidateT.Fatal("lookup differs from accessor")
		}
		if b50ECandidateEntry.B50ECandidateGrade != "Hypothesis" || b50ECandidateEntry.B50ECandidatePosition != nil || b50ECandidateEntry.B50ECandidateWidthBits == nil || b50ECandidateEntry.B50ECandidateRawFallback == "" {
			b50ECandidateT.Fatalf("lost grade, unknown position, carrier width or raw fallback: %+v", b50ECandidateEntry)
		}
	}
	b50ECandidateComparison := b50ECandidateCatalog.B50ECandidateEntries[0]
	b50ECandidateClear := b50ECandidateCatalog.B50ECandidateEntries[1]
	if *b50ECandidateComparison.B50ECandidateWidthBits != 8 || b50ECandidateComparison.B50ECandidateClearIndex != nil || b50ECandidateComparison.B50ECandidateComparison == nil || *b50ECandidateClear.B50ECandidateWidthBits != 16 || b50ECandidateClear.B50ECandidateComparison != nil || b50ECandidateClear.B50ECandidateClearIndex == nil {
		b50ECandidateT.Fatal("candidate alternatives merged or widths changed")
	}
	if b50ECandidateClear.B50ECandidateClearIndex.B50ECandidateByteOrder != "little_endian" || !b50ECandidateClear.B50ECandidateClearIndex.B50ECandidateUnsigned || b50ECandidateClear.B50ECandidateClearIndex.B50EAllRulesCandidate != 65535 {
		b50ECandidateT.Fatal("clear-index metadata changed")
	}
}

func TestAb61ExactKeyAxes(b50ECandidateT *b50ECandidateTesting.T) {
	for _, b50ECandidateControl := range []B50ECandidateKey{b50ECandidateComparisonKey(), b50ECandidateClearKey()} {
		b50ECandidateRequireEntry(b50ECandidateT, b50ECandidateControl)
		b50ECandidateChanges := []struct {
			b50ECandidateName   string
			b50ECandidateChange func(*B50ECandidateKey)
		}{
			{"ab61Family", func(b50ECandidateKey *B50ECandidateKey) { b50ECandidateKey.B50ECandidateFamily = "other_family" }},
			{"ab61FamilyAbsent", func(b50ECandidateKey *B50ECandidateKey) { b50ECandidateKey.B50ECandidateFamily = "" }},
			{"ab61Variant", func(b50ECandidateKey *B50ECandidateKey) { b50ECandidateKey.B50ECandidateVariant = "other_variant" }},
			{"ab61VariantAbsent", func(b50ECandidateKey *B50ECandidateKey) { b50ECandidateKey.B50ECandidateVariant = "" }},
			{"ab61Primary", func(b50ECandidateKey *B50ECandidateKey) { b50ECandidateKey.B50ECandidatePrimary = 0 }},
			{"ab61Secondary", func(b50ECandidateKey *B50ECandidateKey) {
				b50ECandidateKey.B50ECandidateSecondary = b50ECandidatePointer(byte(2))
			}},
			{"ab61SecondaryAbsent", func(b50ECandidateKey *B50ECandidateKey) { b50ECandidateKey.B50ECandidateSecondary = nil }},
			{"ab61SecondaryZero", func(b50ECandidateKey *B50ECandidateKey) {
				b50ECandidateKey.B50ECandidateSecondary = b50ECandidatePointer(byte(0))
			}},
			{"ab61Operation", func(b50ECandidateKey *B50ECandidateKey) {
				b50ECandidateKey.B50ECandidateOperation = b50ECandidatePointer(byte(2))
			}},
			{"ab61OperationZero", func(b50ECandidateKey *B50ECandidateKey) {
				b50ECandidateKey.B50ECandidateOperation = b50ECandidatePointer(byte(0))
			}},
			{"ab61Direction", func(b50ECandidateKey *B50ECandidateKey) {
				b50ECandidateKey.B50ECandidateDirection = b50ECandidatePointer("request")
			}},
			{"ab61DirectionEmpty", func(b50ECandidateKey *B50ECandidateKey) {
				b50ECandidateKey.B50ECandidateDirection = b50ECandidatePointer("")
			}},
			{"ab61Field", func(b50ECandidateKey *B50ECandidateKey) { b50ECandidateKey.B50ECandidateField = "other_field" }},
			{"ab61FieldAbsent", func(b50ECandidateKey *B50ECandidateKey) { b50ECandidateKey.B50ECandidateField = "" }},
		}
		if b50ECandidateControl.B50ECandidateOperation != nil {
			b50ECandidateChanges = append(b50ECandidateChanges, struct {
				b50ECandidateName   string
				b50ECandidateChange func(*B50ECandidateKey)
			}{"ab61OperationAbsent", func(b50ECandidateKey *B50ECandidateKey) { b50ECandidateKey.B50ECandidateOperation = nil }})
			b50ECandidateChanges = append(b50ECandidateChanges, struct {
				b50ECandidateName   string
				b50ECandidateChange func(*B50ECandidateKey)
			}{"ab61DirectionAbsent", func(b50ECandidateKey *B50ECandidateKey) { b50ECandidateKey.B50ECandidateDirection = nil }})
		}
		for _, b50ECandidateCase := range b50ECandidateChanges {
			b50ECandidateT.Run("ab61"+b50ECandidateControl.B50ECandidateField+"/"+b50ECandidateCase.b50ECandidateName, func(b50ECandidateT *b50ECandidateTesting.T) {
				b50ECandidateChanged := b50ECandidateControl
				b50ECandidateCase.b50ECandidateChange(&b50ECandidateChanged)
				b50ECandidateEntry, b50ECandidateFound, b50ECandidateError := B50ELookupCandidate(b50ECandidateChanged)
				if b50ECandidateError != nil || b50ECandidateFound || !b50ECandidateReflect.DeepEqual(b50ECandidateEntry, B50ECandidateEntry{}) {
					b50ECandidateT.Fatalf("non-exact key matched: %+v found=%v error=%v", b50ECandidateChanged, b50ECandidateFound, b50ECandidateError)
				}
			})
		}
	}
}

func TestAb61ComparisonSynthetic(b50ECandidateT *b50ECandidateTesting.T) {
	b50ECandidateEntry := b50ECandidateRequireEntry(b50ECandidateT, b50ECandidateComparisonKey())
	if b50ECandidateEntry.B50ECandidateComparison == nil {
		b50ECandidateT.Fatal("missing comparison metadata")
	}
	b50ECandidateMetadata := b50ECandidateEntry.B50ECandidateComparison
	for _, b50ECandidateCase := range []struct {
		b50ECandidateRaw        byte
		b50ECandidateMode       byte
		b50ECandidateStop       bool
		b50ECandidateUnresolved byte
	}{
		{0x00, 0, false, 0}, {0x07, 7, false, 0}, {0x80, 0, true, 0}, {0x87, 7, true, 0}, {0x88, 0, true, 8},
	} {
		b50ECandidateMode := b50ECandidateCase.b50ECandidateRaw & b50ECandidateMetadata.B50ECandidateModeMask
		b50ECandidateStop := b50ECandidateCase.b50ECandidateRaw&b50ECandidateMetadata.B50ECandidateStopOverrideFIFOMask != 0
		b50ECandidateUnresolved := b50ECandidateCase.b50ECandidateRaw & b50ECandidateMetadata.B50ECandidateUnresolvedMask
		if b50ECandidateMode != b50ECandidateCase.b50ECandidateMode || b50ECandidateStop != b50ECandidateCase.b50ECandidateStop || b50ECandidateUnresolved != b50ECandidateCase.b50ECandidateUnresolved {
			b50ECandidateT.Fatalf("raw %02x: mode=%d stop=%v unresolved=%02x", b50ECandidateCase.b50ECandidateRaw, b50ECandidateMode, b50ECandidateStop, b50ECandidateUnresolved)
		}
		if b50ECandidateMode|b50ECandidateUnresolved|(b50ECandidateCase.b50ECandidateRaw&b50ECandidateMetadata.B50ECandidateStopOverrideFIFOMask) != b50ECandidateCase.b50ECandidateRaw {
			b50ECandidateT.Fatal("raw bits lost")
		}
	}
	for b50ECandidateRaw := 0; b50ECandidateRaw < 256; b50ECandidateRaw++ {
		if byte(b50ECandidateRaw)&b50ECandidateMetadata.B50ECandidateModeMask != byte(b50ECandidateRaw%8) || (byte(b50ECandidateRaw)&b50ECandidateMetadata.B50ECandidateStopOverrideFIFOMask != 0) != (b50ECandidateRaw >= 128) || byte(b50ECandidateRaw)&b50ECandidateMetadata.B50ECandidateUnresolvedMask != byte(b50ECandidateRaw)&0x78 {
			b50ECandidateT.Fatalf("independent scalar/flag/reserved bits changed for %02x", b50ECandidateRaw)
		}
	}
}

func TestAb61ClearEchoSynthetic(b50ECandidateT *b50ECandidateTesting.T) {
	b50ECandidateEntry := b50ECandidateRequireEntry(b50ECandidateT, b50ECandidateClearKey())
	if b50ECandidateEntry.B50ECandidateClearIndex == nil || b50ECandidateEntry.B50ECandidateWidthBits == nil {
		b50ECandidateT.Fatal("missing clear metadata")
	}
	for _, b50ECandidateCase := range []struct {
		b50ECandidateRaw      []byte
		b50ECandidateIndex    uint16
		b50ECandidateComplete bool
		b50ECandidateAll      bool
	}{
		{[]byte{0x34, 0x12}, 4660, true, false}, {[]byte{0xFF, 0xFF}, 65535, true, true}, {[]byte{0xFF}, 0, false, false},
	} {
		b50ECandidateSaved := append([]byte(nil), b50ECandidateCase.b50ECandidateRaw...)
		b50ECandidateComplete := len(b50ECandidateCase.b50ECandidateRaw)*8 == *b50ECandidateEntry.B50ECandidateWidthBits
		var b50ECandidateIndex uint16
		if b50ECandidateComplete {
			if b50ECandidateEntry.B50ECandidateClearIndex.B50ECandidateByteOrder != "little_endian" {
				b50ECandidateT.Fatal("not little endian")
			}
			b50ECandidateIndex = uint16(b50ECandidateCase.b50ECandidateRaw[0]) | uint16(b50ECandidateCase.b50ECandidateRaw[1])<<8
		}
		b50ECandidateAll := b50ECandidateComplete && b50ECandidateIndex == b50ECandidateEntry.B50ECandidateClearIndex.B50EAllRulesCandidate
		if b50ECandidateComplete != b50ECandidateCase.b50ECandidateComplete || b50ECandidateIndex != b50ECandidateCase.b50ECandidateIndex || b50ECandidateAll != b50ECandidateCase.b50ECandidateAll || !b50ECandidateReflect.DeepEqual(b50ECandidateSaved, b50ECandidateCase.b50ECandidateRaw) {
			b50ECandidateT.Fatal("echo value, completeness, sentinel or raw preservation changed")
		}
	}
	b50ECandidateOtherOperation := b50ECandidateClearKey()
	b50ECandidateOtherOperation.B50ECandidateOperation = b50ECandidatePointer(byte(2))
	if _, b50ECandidateFound, b50ECandidateError := B50ELookupCandidate(b50ECandidateOtherOperation); b50ECandidateError != nil || b50ECandidateFound {
		b50ECandidateT.Fatal("FFFF clear candidate escaped OP01")
	}
}

func TestAb61FreshDeepCopies(b50ECandidateT *b50ECandidateTesting.T) {
	b50ECandidateBefore, b50ECandidateError := B50ECandidates()
	if b50ECandidateError != nil {
		b50ECandidateT.Fatal(b50ECandidateError)
	}
	b50ECandidateChanged, b50ECandidateError := B50ECandidates()
	if b50ECandidateError != nil {
		b50ECandidateT.Fatal(b50ECandidateError)
	}
	b50ECandidateChanged.B50ECandidateLimitations[0] = "changed"
	*b50ECandidateChanged.B50ECandidateEntries[0].B50ECandidateKey.B50ECandidateSecondary = 0
	*b50ECandidateChanged.B50ECandidateEntries[0].B50ECandidateWidthBits = 0
	b50ECandidateChanged.B50ECandidateEntries[0].B50ECandidateComparison.B50ECandidateModeMask = 0
	*b50ECandidateChanged.B50ECandidateEntries[1].B50ECandidateKey.B50ECandidateOperation = 0
	*b50ECandidateChanged.B50ECandidateEntries[1].B50ECandidateKey.B50ECandidateDirection = "changed"
	b50ECandidateChanged.B50ECandidateEntries[1].B50ECandidateClearIndex.B50EAllRulesCandidate = 0
	b50ECandidateChanged.B50ECandidateEntries[1].B50ECandidateRawFallback = "changed"
	b50ECandidateLookup := b50ECandidateRequireEntry(b50ECandidateT, b50ECandidateClearKey())
	*b50ECandidateLookup.B50ECandidateWidthBits = 0
	*b50ECandidateLookup.B50ECandidateKey.B50ECandidateSecondary = 0
	*b50ECandidateLookup.B50ECandidateKey.B50ECandidateOperation = 0
	*b50ECandidateLookup.B50ECandidateKey.B50ECandidateDirection = "changed"
	b50ECandidateLookup.B50ECandidateClearIndex.B50ECandidateByteOrder = "changed"
	b50ECandidateComparisonLookup := b50ECandidateRequireEntry(b50ECandidateT, b50ECandidateComparisonKey())
	b50ECandidateComparisonLookup.B50ECandidateComparison.B50ECandidateUnresolvedMask = 0
	b50ECandidateAfter, b50ECandidateError := B50ECandidates()
	if b50ECandidateError != nil || !b50ECandidateReflect.DeepEqual(b50ECandidateBefore, b50ECandidateAfter) {
		b50ECandidateT.Fatal("accessor shares backing storage")
	}
	if !b50ECandidateReflect.DeepEqual(b50ECandidateRequireEntry(b50ECandidateT, b50ECandidateClearKey()), b50ECandidateBefore.B50ECandidateEntries[1]) || !b50ECandidateReflect.DeepEqual(b50ECandidateRequireEntry(b50ECandidateT, b50ECandidateComparisonKey()), b50ECandidateBefore.B50ECandidateEntries[0]) {
		b50ECandidateT.Fatal("lookup shares backing storage")
	}
}
