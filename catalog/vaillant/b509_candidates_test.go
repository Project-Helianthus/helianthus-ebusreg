package vaillant

import (
	b509CandidateBytes "bytes"
	b509CandidateBinary "encoding/binary"
	b509CandidateFmt "fmt"
	b509CandidateReflect "reflect"
	b509CandidateStrings "strings"
	b509CandidateTesting "testing"
)

func b509CandidatePointer[b509CandidateT any](b509CandidateValue b509CandidateT) *b509CandidateT {
	return &b509CandidateValue
}

func b509CandidateKey(b509CandidateOperation byte, b509CandidateVariant, b509CandidateField string, b509CandidateDirection *string) B509CandidateKey {
	return B509CandidateKey{"vaillant", b509CandidateVariant, 0xb5, b509CandidatePointer(byte(9)), b509CandidatePointer(b509CandidateOperation), b509CandidateDirection, b509CandidateField}
}

func b509CandidateCatalog(b509CandidateT *b509CandidateTesting.T) B509CandidateCatalog {
	b509CandidateT.Helper()
	b509CandidateCatalog, b509CandidateError := B509Candidates()
	if b509CandidateError != nil {
		b509CandidateT.Fatal(b509CandidateError)
	}
	return b509CandidateCatalog
}

func b509CandidateLookup(b509CandidateT *b509CandidateTesting.T, b509CandidateKey B509CandidateKey) B509CandidateEntry {
	b509CandidateT.Helper()
	b509CandidateEntry, b509CandidateFound, b509CandidateError := B509CandidateLookup(b509CandidateKey)
	if b509CandidateError != nil || !b509CandidateFound {
		b509CandidateT.Fatalf("ab59 lookup %#v: found=%v error=%v", b509CandidateKey, b509CandidateFound, b509CandidateError)
	}
	return b509CandidateEntry
}

func TestAb59FiniteContract(b509CandidateT *b509CandidateTesting.T) {
	b509CandidateCatalog := b509CandidateCatalog(b509CandidateT)
	if b509CandidateCatalog.B509CandidateMetadataPresence != SchemaProven || b509CandidateCatalog.B509CandidateWireAdaptation != SchemaUnknown || b509CandidateCatalog.B509CandidateInstalledSupport != SchemaUnknown || b509CandidateCatalog.B509CandidateInvocation != SchemaUnknown {
		b509CandidateT.Fatal("ab59 evidence boundaries changed")
	}
	b509CandidateLimitations := []string{
		"Finite candidate declarations only; the broader family catalog remains incomplete.",
		"Candidate semantics are Hypothesis; metadata presence does not establish routable invocation or physical support.",
		"Unprovided coordinates, directions and widths remain null; null keys never wildcard-match.",
		"Raw carriers must be preserved; these declarations do not adapt wire data or select alternatives.",
		"Generic OP0D/0E access and identity scanning are outside these entries.",
	}
	if !b509CandidateReflect.DeepEqual(b509CandidateCatalog.B509CandidateLimitations, b509CandidateLimitations) {
		b509CandidateT.Fatal("ab59 limitations changed")
	}
	b509CandidateExpected := []B509CandidateEntry{}
	b509CandidateAdd := func(b509CandidateOperation byte, b509CandidateVariant, b509CandidateField string, b509CandidateDirection *string, b509CandidateWidth *int, b509CandidateSignedness, b509CandidateOrder, b509CandidateConversion string, b509CandidateDivisor *int, b509CandidateUnit *string, b509CandidateDefault, b509CandidateLower, b509CandidateUpper *int64, b509CandidateMappings map[string]string, b509CandidateOther *string, b509CandidateState string) {
		b509CandidateExpected = append(b509CandidateExpected, B509CandidateEntry{
			B509CandidateKey:       b509CandidateKey(b509CandidateOperation, b509CandidateVariant, b509CandidateField, b509CandidateDirection),
			B509CandidateWidthBits: b509CandidateWidth, B509CandidateSignedness: b509CandidateSignedness, B509CandidateByteOrder: b509CandidateOrder,
			B509CandidateConversion: b509CandidateConversion, B509CandidateDivisor: b509CandidateDivisor, B509CandidateUnit: b509CandidateUnit,
			B509CandidateDefault: b509CandidateDefault, B509CandidateLower: b509CandidateLower, B509CandidateUpper: b509CandidateUpper,
			B509CandidateMappings: b509CandidateMappings, B509CandidateOtherFeedback: b509CandidateOther, B509CandidateState: b509CandidateState, B509CandidateGrade: SchemaHypothesis,
		})
	}
	b509CandidateAdd(0x07, "selector", "selector", b509CandidatePointer("request"), b509CandidatePointer(8), "unsigned", "not_applicable", "raw", nil, nil, nil, nil, nil, map[string]string{}, nil, "raw")
	b509CandidateAdd(0x07, "positive-u16", "value", b509CandidatePointer("response"), b509CandidatePointer(16), "unsigned", "little_endian", "raw", nil, nil, b509CandidatePointer(int64(0)), nil, nil, map[string]string{}, nil, "positive-candidate")
	b509CandidateAdd(0x07, "positive-u32", "value", b509CandidatePointer("response"), b509CandidatePointer(32), "unsigned", "little_endian", "raw", nil, nil, b509CandidatePointer(int64(0)), nil, nil, map[string]string{}, nil, "positive-candidate")
	b509CandidateAdd(0x32, "status-query", "address", b509CandidatePointer("request"), nil, "Unknown", "Unknown", "Unknown", nil, nil, nil, nil, nil, map[string]string{}, nil, "raw")
	b509CandidateAdd(0x33, "raw-quantities", "address", nil, nil, "Unknown", "Unknown", "Unknown", nil, nil, nil, nil, nil, map[string]string{}, nil, "raw")
	for _, b509CandidateField := range []string{"minimum", "maximum", "default", "step", "factor"} {
		b509CandidateAdd(0x33, "raw-quantities", b509CandidateField, nil, b509CandidatePointer(16), "signed", "little_endian", "raw", nil, nil, nil, nil, nil, map[string]string{}, nil, "raw")
	}
	b509CandidateAdd(0x39, "raw-u16", "value", nil, b509CandidatePointer(16), "unsigned", "Unknown", "raw", nil, nil, nil, nil, nil, map[string]string{}, nil, "raw")
	b509CandidateAdd(0x41, "raw-0f-f0", "value", nil, b509CandidatePointer(8), "unsigned", "not_applicable", "raw", nil, nil, nil, nil, nil, map[string]string{"15": "raw 0F", "240": "raw F0"}, nil, "raw")
	b509CandidateAdd(0x44, "default-domain-conflict", "value", nil, b509CandidatePointer(8), "unsigned", "not_applicable", "Unknown", nil, b509CandidatePointer("minutes"), b509CandidatePointer(int64(0)), b509CandidatePointer(int64(1)), b509CandidatePointer(int64(255)), map[string]string{}, nil, "raw")
	b509CandidateAdd(0x59, "quantities", "temperature", nil, b509CandidatePointer(16), "signed", "little_endian", "divide", b509CandidatePointer(16), b509CandidatePointer("degrees Celsius"), nil, nil, nil, map[string]string{}, nil, "raw")
	b509CandidateAdd(0x59, "quantities", "flow", nil, b509CandidatePointer(16), "unsigned", "little_endian", "divide", b509CandidatePointer(100), b509CandidatePointer("litres per minute"), nil, nil, nil, map[string]string{}, nil, "raw")
	b509CandidateAdd(0x59, "quantities", "volume", nil, b509CandidatePointer(16), "unsigned", "little_endian", "divide", b509CandidatePointer(10), b509CandidatePointer("litres"), nil, nil, nil, map[string]string{}, nil, "raw")
	for _, b509CandidateField := range []string{"minimum_volume", "maximum_volume"} {
		b509CandidateAdd(0x59, "quantities", b509CandidateField, nil, nil, "Unknown", "Unknown", "raw", nil, nil, nil, nil, nil, map[string]string{}, nil, "raw")
	}
	b509CandidateAdd(0x06, "prefix-only-positive", "prefix", b509CandidatePointer("response"), nil, "Unknown", "Unknown", "Unknown", nil, nil, nil, nil, nil, map[string]string{}, nil, "prefix-only-positive")
	b509CandidateAdd(0x06, "positive-with-feedback", "feedback", b509CandidatePointer("response"), b509CandidatePointer(8), "unsigned", "not_applicable", "raw", nil, nil, nil, nil, nil, map[string]string{"1": "accepted"}, b509CandidatePointer("rejected"), "explicit-feedback")
	b509CandidateAdd(0x06, "negative-without-extra-discriminator", "prefix", b509CandidatePointer("response"), nil, "Unknown", "Unknown", "Unknown", nil, nil, nil, nil, nil, map[string]string{}, nil, "ambiguous-negative")
	if len(b509CandidateCatalog.B509CandidateEntries) != 21 {
		b509CandidateT.Fatalf("ab59 expected 21 finite entries, got %d", len(b509CandidateCatalog.B509CandidateEntries))
	}
	for b509CandidateIndex, b509CandidateWant := range b509CandidateExpected {
		b509CandidateT.Run(b509CandidateFmt.Sprintf("ab59_row_%d", b509CandidateIndex), func(b509CandidateT *b509CandidateTesting.T) {
			b509CandidateGot := b509CandidateCatalog.B509CandidateEntries[b509CandidateIndex]
			if !b509CandidateReflect.DeepEqual(b509CandidateGot, b509CandidateLookup(b509CandidateT, b509CandidateWant.B509CandidateKey)) {
				b509CandidateT.Fatal("ab59 accessor and lookup differ")
			}
			b509CandidateGot.B509CandidateNotes = nil
			if !b509CandidateReflect.DeepEqual(b509CandidateGot, b509CandidateWant) {
				b509CandidateT.Fatalf("ab59 metadata differs: got %#v want %#v", b509CandidateGot, b509CandidateWant)
			}
		})
	}
	b509CandidateNotes := map[int][]string{
		0: {"cannot choose the reply width"}, 1: {"incomplete/raw"}, 2: {"incomplete/raw"},
		3: {"does not toggle", "Address width is Unknown"}, 4: {"Address width is Unknown"},
		5: {"enforcement are Unknown", "not a converter"}, 6: {"enforcement are Unknown", "not a converter"},
		7: {"enforcement are Unknown", "not a converter"}, 8: {"enforcement are Unknown", "not a converter"}, 9: {"enforcement are Unknown", "not a converter"},
		10: {"Complement algorithm is Unknown"}, 11: {"Other values have no established labels"},
		12: {"conflicts with domain 1..255", "Hypothesis"}, 16: {"without division"}, 17: {"without division"},
		18: {"does not select an alternative", "successful reset"}, 19: {"only in this explicitly selected alternative", "does not establish successful reset"}, 20: {"No extra discriminator", "does not establish physical NACK"},
	}
	for b509CandidateIndex, b509CandidateFragments := range b509CandidateNotes {
		for _, b509CandidateFragment := range b509CandidateFragments {
			if !b509CandidateStrings.Contains(b509CandidateStrings.Join(b509CandidateCatalog.B509CandidateEntries[b509CandidateIndex].B509CandidateNotes, " "), b509CandidateFragment) {
				b509CandidateT.Fatalf("ab59 row %d missing caveat %q", b509CandidateIndex, b509CandidateFragment)
			}
		}
	}
}

func TestAb59ExactKeyAxes(b509CandidateT *b509CandidateTesting.T) {
	for b509CandidateIndex, b509CandidateEntry := range b509CandidateCatalog(b509CandidateT).B509CandidateEntries {
		b509CandidateControl := b509CandidateEntry.B509CandidateKey
		b509CandidateChanges := map[string]func(*B509CandidateKey){
			"ab59_family":           func(b509CandidateKey *B509CandidateKey) { b509CandidateKey.B509CandidateFamily = "ab59_other" },
			"ab59_variant":          func(b509CandidateKey *B509CandidateKey) { b509CandidateKey.B509CandidateVariant = "ab59_other" },
			"ab59_primary":          func(b509CandidateKey *B509CandidateKey) { b509CandidateKey.B509CandidatePrimary = 0 },
			"ab59_secondary_absent": func(b509CandidateKey *B509CandidateKey) { b509CandidateKey.B509CandidateSecondary = nil },
			"ab59_secondary_zero": func(b509CandidateKey *B509CandidateKey) {
				b509CandidateKey.B509CandidateSecondary = b509CandidatePointer(byte(0))
			},
			"ab59_operation_absent": func(b509CandidateKey *B509CandidateKey) { b509CandidateKey.B509CandidateOperation = nil },
			"ab59_operation_zero": func(b509CandidateKey *B509CandidateKey) {
				b509CandidateKey.B509CandidateOperation = b509CandidatePointer(byte(0))
			},
			"ab59_direction": func(b509CandidateKey *B509CandidateKey) {
				b509CandidateKey.B509CandidateDirection = b509CandidatePointer("ab59_other")
			},
			"ab59_direction_empty": func(b509CandidateKey *B509CandidateKey) {
				b509CandidateKey.B509CandidateDirection = b509CandidatePointer("")
			},
			"ab59_field": func(b509CandidateKey *B509CandidateKey) { b509CandidateKey.B509CandidateField = "ab59_other" },
		}
		if b509CandidateControl.B509CandidateDirection != nil {
			b509CandidateChanges["ab59_direction_absent"] = func(b509CandidateKey *B509CandidateKey) { b509CandidateKey.B509CandidateDirection = nil }
		} else {
			b509CandidateChanges["ab59_direction_present"] = func(b509CandidateKey *B509CandidateKey) {
				b509CandidateKey.B509CandidateDirection = b509CandidatePointer("response")
			}
		}
		for b509CandidateName, b509CandidateChange := range b509CandidateChanges {
			b509CandidateT.Run(b509CandidateFmt.Sprintf("ab59_row_%d/%s", b509CandidateIndex, b509CandidateName), func(b509CandidateT *b509CandidateTesting.T) {
				if !b509CandidateReflect.DeepEqual(b509CandidateLookup(b509CandidateT, b509CandidateControl), b509CandidateEntry) {
					b509CandidateT.Fatal("ab59 invalid positive control")
				}
				b509CandidateChanged := b509CandidateControl
				b509CandidateChange(&b509CandidateChanged)
				b509CandidateGot, b509CandidateFound, b509CandidateError := B509CandidateLookup(b509CandidateChanged)
				if b509CandidateError != nil || b509CandidateFound || !b509CandidateReflect.DeepEqual(b509CandidateGot, B509CandidateEntry{}) {
					b509CandidateT.Fatalf("ab59 changed axis matched: %v %v", b509CandidateFound, b509CandidateError)
				}
				if !b509CandidateReflect.DeepEqual(b509CandidateLookup(b509CandidateT, b509CandidateControl), b509CandidateEntry) {
					b509CandidateT.Fatal("ab59 control damaged")
				}
			})
		}
	}
}

func b509CandidatePoison(b509CandidateEntry *B509CandidateEntry) {
	b509CandidateEntry.B509CandidateKey.B509CandidateFamily = "ab59_mutated"
	if b509CandidateEntry.B509CandidateKey.B509CandidateSecondary != nil {
		*b509CandidateEntry.B509CandidateKey.B509CandidateSecondary = 0
	}
	if b509CandidateEntry.B509CandidateKey.B509CandidateOperation != nil {
		*b509CandidateEntry.B509CandidateKey.B509CandidateOperation = 0
	}
	if b509CandidateEntry.B509CandidateKey.B509CandidateDirection != nil {
		*b509CandidateEntry.B509CandidateKey.B509CandidateDirection = "ab59_mutated"
	}
	for _, b509CandidatePointer := range []*int{b509CandidateEntry.B509CandidatePosition, b509CandidateEntry.B509CandidateWidthBits, b509CandidateEntry.B509CandidateDivisor} {
		if b509CandidatePointer != nil {
			*b509CandidatePointer = -999
		}
	}
	for _, b509CandidatePointer := range []*int64{b509CandidateEntry.B509CandidateDefault, b509CandidateEntry.B509CandidateLower, b509CandidateEntry.B509CandidateUpper} {
		if b509CandidatePointer != nil {
			*b509CandidatePointer = -999
		}
	}
	for _, b509CandidatePointer := range []*string{b509CandidateEntry.B509CandidateUnit, b509CandidateEntry.B509CandidateOtherFeedback} {
		if b509CandidatePointer != nil {
			*b509CandidatePointer = "ab59_mutated"
		}
	}
	b509CandidateEntry.B509CandidateMappings["ab59_mutated"] = "ab59_mutated"
	for b509CandidateIndex := range b509CandidateEntry.B509CandidateNotes {
		b509CandidateEntry.B509CandidateNotes[b509CandidateIndex] = "ab59_mutated"
	}
}

func TestAb59DeepCopy(b509CandidateT *b509CandidateTesting.T) {
	b509CandidateBefore := b509CandidateCatalog(b509CandidateT)
	b509CandidateReturned := b509CandidateCatalog(b509CandidateT)
	b509CandidateReturned.B509CandidateLimitations[0] = "ab59_mutated"
	for b509CandidateIndex := range b509CandidateReturned.B509CandidateEntries {
		b509CandidatePoison(&b509CandidateReturned.B509CandidateEntries[b509CandidateIndex])
	}
	if !b509CandidateReflect.DeepEqual(b509CandidateCatalog(b509CandidateT), b509CandidateBefore) {
		b509CandidateT.Fatal("ab59 nested accessor storage leaked")
	}
	for _, b509CandidateEntry := range b509CandidateBefore.B509CandidateEntries {
		b509CandidateReturnedEntry := b509CandidateLookup(b509CandidateT, b509CandidateEntry.B509CandidateKey)
		b509CandidatePoison(&b509CandidateReturnedEntry)
		if !b509CandidateReflect.DeepEqual(b509CandidateLookup(b509CandidateT, b509CandidateEntry.B509CandidateKey), b509CandidateEntry) {
			b509CandidateT.Fatal("ab59 nested lookup storage leaked")
		}
	}
	if !b509CandidateReflect.DeepEqual(b509CandidateCatalog(b509CandidateT), b509CandidateBefore) {
		b509CandidateT.Fatal("ab59 lookup contaminated accessor")
	}
}

func b509CandidateInterpret(b509CandidateEntry B509CandidateEntry, b509CandidateRaw []byte) (float64, bool) {
	if b509CandidateEntry.B509CandidateWidthBits == nil || *b509CandidateEntry.B509CandidateWidthBits != len(b509CandidateRaw)*8 || b509CandidateEntry.B509CandidateByteOrder != "little_endian" {
		return 0, false
	}
	var b509CandidateValue float64
	switch *b509CandidateEntry.B509CandidateWidthBits {
	case 16:
		b509CandidateWord := b509CandidateBinary.LittleEndian.Uint16(b509CandidateRaw)
		switch b509CandidateEntry.B509CandidateSignedness {
		case "signed":
			b509CandidateValue = float64(int16(b509CandidateWord))
		case "unsigned":
			b509CandidateValue = float64(b509CandidateWord)
		default:
			return 0, false
		}
	case 32:
		if b509CandidateEntry.B509CandidateSignedness != "unsigned" {
			return 0, false
		}
		b509CandidateValue = float64(b509CandidateBinary.LittleEndian.Uint32(b509CandidateRaw))
	default:
		return 0, false
	}
	switch b509CandidateEntry.B509CandidateConversion {
	case "raw":
		return b509CandidateValue, true
	case "divide":
		if b509CandidateEntry.B509CandidateDivisor == nil || *b509CandidateEntry.B509CandidateDivisor <= 0 {
			return 0, false
		}
		return b509CandidateValue / float64(*b509CandidateEntry.B509CandidateDivisor), true
	default:
		return 0, false
	}
}

func TestAb59SyntheticQuantities(b509CandidateT *b509CandidateTesting.T) {
	b509CandidateCases := []struct {
		B509CandidateKey      B509CandidateKey
		B509CandidateRaw      []byte
		B509CandidateWant     float64
		B509CandidateComplete bool
	}{
		{b509CandidateKey(0x07, "positive-u16", "value", b509CandidatePointer("response")), []byte{0x34, 0x12}, 4660, true},
		{b509CandidateKey(0x07, "positive-u32", "value", b509CandidatePointer("response")), []byte{0x78, 0x56, 0x34, 0x12}, 305419896, true},
		{b509CandidateKey(0x07, "positive-u16", "value", b509CandidatePointer("response")), []byte{0x78, 0x56, 0x34, 0x12}, 0, false},
		{b509CandidateKey(0x07, "positive-u32", "value", b509CandidatePointer("response")), []byte{0x34, 0x12}, 0, false},
		{b509CandidateKey(0x59, "quantities", "temperature", nil), []byte{0xf0, 0xff}, -1, true},
		{b509CandidateKey(0x59, "quantities", "flow", nil), []byte{0x64, 0}, 1, true},
		{b509CandidateKey(0x59, "quantities", "flow", nil), []byte{0xff, 0xff}, 655.35, true},
		{b509CandidateKey(0x59, "quantities", "volume", nil), []byte{0x0a, 0}, 1, true},
	}
	for b509CandidateIndex, b509CandidateCase := range b509CandidateCases {
		b509CandidateT.Run(b509CandidateFmt.Sprintf("ab59_sample_%d", b509CandidateIndex), func(b509CandidateT *b509CandidateTesting.T) {
			b509CandidateBefore := b509CandidateBytes.Clone(b509CandidateCase.B509CandidateRaw)
			b509CandidateValue, b509CandidateComplete := b509CandidateInterpret(b509CandidateLookup(b509CandidateT, b509CandidateCase.B509CandidateKey), b509CandidateCase.B509CandidateRaw)
			if b509CandidateValue != b509CandidateCase.B509CandidateWant || b509CandidateComplete != b509CandidateCase.B509CandidateComplete || !b509CandidateBytes.Equal(b509CandidateCase.B509CandidateRaw, b509CandidateBefore) {
				b509CandidateT.Fatalf("ab59 sample got %v %v", b509CandidateValue, b509CandidateComplete)
			}
		})
	}
	for b509CandidateCaseIndex, b509CandidateWords := range [][]int16{{-1, 2, 0, 1, -2}, {2, -1, 0, 1, -2}} {
		for b509CandidateIndex, b509CandidateField := range []string{"minimum", "maximum", "default", "step", "factor"} {
			b509CandidateRaw := make([]byte, 2)
			b509CandidateBinary.LittleEndian.PutUint16(b509CandidateRaw, uint16(b509CandidateWords[b509CandidateIndex]))
			b509CandidateEntry := b509CandidateLookup(b509CandidateT, b509CandidateKey(0x33, "raw-quantities", b509CandidateField, nil))
			b509CandidateValue, b509CandidateComplete := b509CandidateInterpret(b509CandidateEntry, b509CandidateRaw)
			if !b509CandidateComplete || b509CandidateValue != float64(b509CandidateWords[b509CandidateIndex]) || b509CandidateEntry.B509CandidateDivisor != nil {
				b509CandidateT.Fatalf("ab59 OP33 case %d field %s lost independent raw quantity", b509CandidateCaseIndex, b509CandidateField)
			}
		}
	}
	b509CandidateRaw39 := []byte{0xff, 0xff}
	b509CandidateEntry39 := b509CandidateLookup(b509CandidateT, b509CandidateKey(0x39, "raw-u16", "value", nil))
	if b509CandidateBinary.LittleEndian.Uint16(b509CandidateRaw39) != 65535 || b509CandidateEntry39.B509CandidateConversion != "raw" || b509CandidateEntry39.B509CandidateDivisor != nil {
		b509CandidateT.Fatal("ab59 OP39 complemented")
	}
	for _, b509CandidateField := range []string{"minimum_volume", "maximum_volume"} {
		b509CandidateEntry := b509CandidateLookup(b509CandidateT, b509CandidateKey(0x59, "quantities", b509CandidateField, nil))
		b509CandidateRaw := []byte{0x0a, 0}
		b509CandidateBefore := b509CandidateBytes.Clone(b509CandidateRaw)
		_, b509CandidateComplete := b509CandidateInterpret(b509CandidateEntry, b509CandidateRaw)
		if b509CandidateComplete || b509CandidateEntry.B509CandidateConversion != "raw" || b509CandidateEntry.B509CandidateDivisor != nil || !b509CandidateBytes.Equal(b509CandidateRaw, b509CandidateBefore) {
			b509CandidateT.Fatal("ab59 raw volume acquired unsupported conversion")
		}
	}
}

func TestAb59FeedbackAlternatives(b509CandidateT *b509CandidateTesting.T) {
	b509CandidateFeedback := b509CandidateLookup(b509CandidateT, b509CandidateKey(0x06, "positive-with-feedback", "feedback", b509CandidatePointer("response")))
	for b509CandidateIndex, b509CandidateRaw := range []byte{0, 1, 2, 255} {
		b509CandidateGot, b509CandidateFound := b509CandidateFeedback.B509CandidateMappings[b509CandidateFmt.Sprint(b509CandidateRaw)]
		if !b509CandidateFound && b509CandidateFeedback.B509CandidateOtherFeedback != nil {
			b509CandidateGot = *b509CandidateFeedback.B509CandidateOtherFeedback
		}
		if b509CandidateGot != []string{"rejected", "accepted", "rejected", "rejected"}[b509CandidateIndex] {
			b509CandidateT.Fatalf("ab59 feedback %d got %s", b509CandidateRaw, b509CandidateGot)
		}
	}
	for b509CandidateIndex, b509CandidateVariant := range []string{"prefix-only-positive", "negative-without-extra-discriminator"} {
		b509CandidateEntry := b509CandidateLookup(b509CandidateT, b509CandidateKey(0x06, b509CandidateVariant, "prefix", b509CandidatePointer("response")))
		if b509CandidateEntry.B509CandidateState != []string{"prefix-only-positive", "ambiguous-negative"}[b509CandidateIndex] || len(b509CandidateEntry.B509CandidateMappings) != 0 || b509CandidateEntry.B509CandidateOtherFeedback != nil || b509CandidateEntry.B509CandidateWidthBits != nil {
			b509CandidateT.Fatal("ab59 apparent prefix classified as feedback")
		}
	}
}

func TestAb59AbsentAndZeroKeys(b509CandidateT *b509CandidateTesting.T) {
	b509CandidateBase := b509CandidateCatalog(b509CandidateT).B509CandidateEntries[0]
	for _, b509CandidateAxis := range []string{"ab59_secondary", "ab59_operation"} {
		b509CandidateT.Run(b509CandidateAxis, func(b509CandidateT *b509CandidateTesting.T) {
			b509CandidateAbsent := b509CandidateBase
			b509CandidateZero := b509CandidateBase
			if b509CandidateAxis == "ab59_secondary" {
				b509CandidateAbsent.B509CandidateKey.B509CandidateSecondary = nil
				b509CandidateZero.B509CandidateKey.B509CandidateSecondary = b509CandidatePointer(byte(0))
			} else {
				b509CandidateAbsent.B509CandidateKey.B509CandidateOperation = nil
				b509CandidateZero.B509CandidateKey.B509CandidateOperation = b509CandidatePointer(byte(0))
			}
			for _, b509CandidateControl := range []B509CandidateEntry{b509CandidateAbsent, b509CandidateZero} {
				b509CandidateGot, b509CandidateFound := b509CandidateFind([]B509CandidateEntry{b509CandidateControl}, b509CandidateControl.B509CandidateKey)
				if !b509CandidateFound || !b509CandidateReflect.DeepEqual(b509CandidateGot, b509CandidateControl) {
					b509CandidateT.Fatal("ab59 optional-axis positive control failed")
				}
				b509CandidateChanged := b509CandidateZero.B509CandidateKey
				if b509CandidateReflect.DeepEqual(b509CandidateChanged, b509CandidateControl.B509CandidateKey) {
					b509CandidateChanged = b509CandidateAbsent.B509CandidateKey
				}
				if _, b509CandidateFound = b509CandidateFind([]B509CandidateEntry{b509CandidateControl}, b509CandidateChanged); b509CandidateFound {
					b509CandidateT.Fatal("ab59 absent and zero matched")
				}
			}
		})
	}
}
