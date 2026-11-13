package vaillant

import (
	b50FCandidateFmt "fmt"
	b50FCandidateReflect "reflect"
	b50FCandidateTesting "testing"
)

func b50FCandidateByte(b50FCandidateValue byte) *byte       { return &b50FCandidateValue }
func b50FCandidateString(b50FCandidateValue string) *string { return &b50FCandidateValue }
func b50FCandidateKey(b50FCandidateOperation *byte, b50FCandidateField string) B50FCandidateKey {
	return B50FCandidateKey{B50FCandidateFamily: "vaillant", B50FCandidateVariant: "B50F", B50FCandidatePrimary: 0xB5, B50FCandidateSecondary: b50FCandidateByte(0x0F), B50FCandidateOperation: b50FCandidateOperation, B50FCandidateField: b50FCandidateField}
}
func b50FCandidateGet(b50FCandidateT *b50FCandidateTesting.T, b50FCandidateOperation *byte, b50FCandidateName string) B50FCandidateEntry {
	b50FCandidateT.Helper()
	b50FCandidateEntry, b50FCandidateOK := B50FCandidateLookup(b50FCandidateKey(b50FCandidateOperation, b50FCandidateName))
	if !b50FCandidateOK {
		b50FCandidateT.Fatalf("missing %v/%s", b50FCandidateOperation, b50FCandidateName)
	}
	return b50FCandidateEntry
}

func TestAb63FiniteCatalog(b50FCandidateT *b50FCandidateTesting.T) {
	b50FCandidateCatalog := B50FCandidates()
	if b50FCandidateCatalog.B50FCandidatePresence != SchemaProven || b50FCandidateCatalog.B50FCandidateQualification != SchemaUnknown || b50FCandidateCatalog.B50FCandidateInvocation != SchemaUnknown || len(b50FCandidateCatalog.B50FCandidateEntries) != 21 {
		b50FCandidateT.Fatal("finite catalog or grade separation")
	}
	b50FCandidateExpected := map[string]int{"demand": 1, "stack-voltage": 1, "feedback": 1, "temperature_1": 2, "temperature_1_error": 2, "temperature_2": 2, "temperature_2_error": 2, "temperature_3": 1, "temperature_3_error": 1, "integrity_1": 4, "integrity_2": 4}
	for _, b50FCandidateEntry := range b50FCandidateCatalog.B50FCandidateEntries {
		b50FCandidateExpected[b50FCandidateEntry.B50FCandidateKey.B50FCandidateField]--
		if b50FCandidateEntry.B50FCandidateKey.B50FCandidateFamily != "vaillant" || b50FCandidateEntry.B50FCandidateKey.B50FCandidateVariant != "B50F" || b50FCandidateEntry.B50FCandidateKey.B50FCandidatePrimary != 0xB5 || b50FCandidateEntry.B50FCandidateKey.B50FCandidateSecondary == nil || *b50FCandidateEntry.B50FCandidateKey.B50FCandidateSecondary != 0x0F {
			b50FCandidateT.Fatal("native key axes")
		}
		if b50FCandidateEntry.B50FCandidateField.SchemaPosition == nil && b50FCandidateEntry.B50FCandidateCoordinates != "Unknown" {
			b50FCandidateT.Fatal("unknown scalar coordinates")
		}
		if b50FCandidateEntry.B50FCandidateField.SchemaWidth == nil {
			b50FCandidateT.Fatal("known carrier width lost")
		}
		if b50FCandidateEntry.B50FCandidateKey.B50FCandidateField == "demand" || b50FCandidateEntry.B50FCandidateKey.B50FCandidateField == "feedback" {
			if *b50FCandidateEntry.B50FCandidateField.SchemaWidth != 1 || b50FCandidateEntry.B50FCandidateField.SchemaPosition != nil || b50FCandidateEntry.B50FCandidateField.SchemaCodec != "raw" {
				b50FCandidateT.Fatal("raw byte metadata")
			}
		}
		if b50FCandidateEntry.B50FCandidateKey.B50FCandidateField == "stack-voltage" && *b50FCandidateEntry.B50FCandidateField.SchemaWidth != 2 {
			b50FCandidateT.Fatal("stack raw width")
		}
		if b50FCandidateEntry.B50FCandidateField.SchemaGrade != SchemaHypothesis || b50FCandidateEntry.B50FCandidateField.SchemaName != b50FCandidateEntry.B50FCandidateKey.B50FCandidateField || b50FCandidateEntry.B50FCandidateKey.B50FCandidateDirection != nil || b50FCandidateEntry.B50FCandidateField.SchemaRawFallback != nil {
			b50FCandidateT.Fatal("grade, name, unknown direction or native fallback")
		}
		b50FCandidateFound, b50FCandidateOK := B50FCandidateLookup(b50FCandidateEntry.B50FCandidateKey)
		if !b50FCandidateOK || !b50FCandidateReflect.DeepEqual(b50FCandidateFound, b50FCandidateEntry) {
			b50FCandidateT.Fatal("accessor/lookup mismatch")
		}
	}
	for b50FCandidateName, b50FCandidateCount := range b50FCandidateExpected {
		if b50FCandidateCount != 0 {
			b50FCandidateT.Fatalf("unexpected field count %s: %d", b50FCandidateName, b50FCandidateCount)
		}
	}
}

func TestAb63ExactKeyAxes(b50FCandidateT *b50FCandidateTesting.T) {
	for b50FCandidateIndex, b50FCandidateEntry := range B50FCandidates().B50FCandidateEntries {
		for _, b50FCandidateAxis := range []string{"family", "variant", "primary", "secondary", "secondary_absent", "secondary_zero", "operation", "operation_zero_or_absent", "direction", "direction_empty", "field"} {
			b50FCandidateT.Run(b50FCandidateFmt.Sprintf("ab63_%d_%s", b50FCandidateIndex, b50FCandidateAxis), func(b50FCandidateT *b50FCandidateTesting.T) {
				b50FCandidateKeyValue := b50FCandidateEntry.B50FCandidateKey
				switch b50FCandidateAxis {
				case "family":
					b50FCandidateKeyValue.B50FCandidateFamily = ""
				case "variant":
					b50FCandidateKeyValue.B50FCandidateVariant = ""
				case "primary":
					b50FCandidateKeyValue.B50FCandidatePrimary = 0
				case "secondary":
					b50FCandidateKeyValue.B50FCandidateSecondary = b50FCandidateByte(16)
				case "secondary_absent":
					b50FCandidateKeyValue.B50FCandidateSecondary = nil
				case "secondary_zero":
					b50FCandidateKeyValue.B50FCandidateSecondary = b50FCandidateByte(0)
				case "operation":
					b50FCandidateKeyValue.B50FCandidateOperation = b50FCandidateByte(255)
				case "operation_zero_or_absent":
					if b50FCandidateKeyValue.B50FCandidateOperation == nil {
						b50FCandidateKeyValue.B50FCandidateOperation = b50FCandidateByte(0)
					} else {
						b50FCandidateKeyValue.B50FCandidateOperation = nil
					}
				case "direction":
					b50FCandidateKeyValue.B50FCandidateDirection = b50FCandidateString("response")
				case "direction_empty":
					b50FCandidateKeyValue.B50FCandidateDirection = b50FCandidateString("")
				case "field":
					b50FCandidateKeyValue.B50FCandidateField = "ab63_unknown"
				}
				if _, b50FCandidateOK := B50FCandidateLookup(b50FCandidateKeyValue); b50FCandidateOK {
					b50FCandidateT.Fatal("axis wildcarded or conflated")
				}
				if _, b50FCandidateOK := B50FCandidateLookup(b50FCandidateEntry.B50FCandidateKey); !b50FCandidateOK {
					b50FCandidateT.Fatal("valid control rejected")
				}
			})
		}
	}
}

func TestAb63Demand(b50FCandidateT *b50FCandidateTesting.T) {
	b50FCandidateEntry := b50FCandidateGet(b50FCandidateT, b50FCandidateByte(3), "demand")
	if !b50FCandidateReflect.DeepEqual(b50FCandidateEntry.B50FCandidateMasks, map[string]byte{"burner": 1, "reformer-gas-valve": 2, "gas-flow": 4, "fan-power": 8, "valve": 16}) || b50FCandidateEntry.B50FCandidateReserved == nil || *b50FCandidateEntry.B50FCandidateReserved != 224 || b50FCandidateEntry.B50FCandidateInterpretation != "demand bits; not applied actuation" || !b50FCandidateReflect.DeepEqual(b50FCandidateEntry.B50FCandidateExclusions, []string{"applied actuation"}) {
		b50FCandidateT.Fatal("demand metadata")
	}
	for b50FCandidateIndex, b50FCandidateRaw := range []byte{0, 1, 31, 224, 255} {
		b50FCandidateWant := [][]bool{{false, false, false, false, false}, {true, false, false, false, false}, {true, true, true, true, true}, {false, false, false, false, false}, {true, true, true, true, true}}[b50FCandidateIndex]
		for b50FCandidateBit, b50FCandidateName := range []string{"burner", "reformer-gas-valve", "gas-flow", "fan-power", "valve"} {
			if (b50FCandidateRaw&b50FCandidateEntry.B50FCandidateMasks[b50FCandidateName] != 0) != b50FCandidateWant[b50FCandidateBit] {
				b50FCandidateT.Fatal("demand example")
			}
		}
		if b50FCandidateRaw&*b50FCandidateEntry.B50FCandidateReserved != []byte{0, 0, 0, 224, 224}[b50FCandidateIndex] {
			b50FCandidateT.Fatal("reserved bits")
		}
	}
}

func TestAb63TemperatureAndConflict(b50FCandidateT *b50FCandidateTesting.T) {
	for b50FCandidateOp, b50FCandidateCount := range map[byte]int{0: 3, 1: 2} {
		for b50FCandidateNumber := 1; b50FCandidateNumber <= b50FCandidateCount; b50FCandidateNumber++ {
			b50FCandidateName := b50FCandidateFmt.Sprintf("temperature_%d", b50FCandidateNumber)
			b50FCandidateEntry := b50FCandidateGet(b50FCandidateT, b50FCandidateByte(b50FCandidateOp), b50FCandidateName)
			b50FCandidateError := b50FCandidateGet(b50FCandidateT, b50FCandidateByte(b50FCandidateOp), b50FCandidateName+"_error")
			if b50FCandidateEntry.B50FCandidateField.SchemaPosition != nil || b50FCandidateEntry.B50FCandidateField.SchemaWidth == nil || *b50FCandidateEntry.B50FCandidateField.SchemaWidth != 2 || b50FCandidateEntry.B50FCandidateField.SchemaCodec != "signed little-endian int16" || b50FCandidateEntry.B50FCandidateField.SchemaUnit != "degrees Celsius" || b50FCandidateEntry.B50FCandidateDivisor == nil || *b50FCandidateEntry.B50FCandidateDivisor != 16 {
				b50FCandidateT.Fatal("temperature metadata")
			}
			if b50FCandidateError.B50FCandidateField.SchemaPosition != nil || b50FCandidateError.B50FCandidateField.SchemaWidth == nil || *b50FCandidateError.B50FCandidateField.SchemaWidth != 1 || b50FCandidateError.B50FCandidateField.SchemaCodec != "raw" || b50FCandidateError.B50FCandidateDivisor != nil || b50FCandidateError.B50FCandidateInterpretation != "raw error byte immediately following "+b50FCandidateName || !b50FCandidateReflect.DeepEqual(b50FCandidateError.B50FCandidateExclusions, []string{"temperature conversion"}) {
				b50FCandidateT.Fatal("error byte independence")
			}
			for b50FCandidateIndex, b50FCandidatePair := range [][2]byte{{240, 255}, {0, 0}, {16, 0}} {
				for _, b50FCandidateErrorRaw := range []byte{0, 1, 255} {
					b50FCandidateRaw := [3]byte{b50FCandidatePair[0], b50FCandidatePair[1], b50FCandidateErrorRaw}
					b50FCandidateValue := float64(int16(uint16(b50FCandidateRaw[0])|uint16(b50FCandidateRaw[1])<<8)) / float64(*b50FCandidateEntry.B50FCandidateDivisor)
					if b50FCandidateValue != []float64{-1, 0, 1}[b50FCandidateIndex] || b50FCandidateRaw[2] != b50FCandidateErrorRaw {
						b50FCandidateT.Fatal("signed sample or adjacent error")
					}
				}
			}
		}
	}
	b50FCandidateStack := b50FCandidateGet(b50FCandidateT, b50FCandidateByte(1), "stack-voltage")
	b50FCandidateRaw := [2]byte{16, 0}
	if b50FCandidateStack.B50FCandidateField.SchemaCodec != "raw" || b50FCandidateStack.B50FCandidateField.SchemaUnit != "" || b50FCandidateStack.B50FCandidateDivisor != nil || b50FCandidateStack.B50FCandidateField.SchemaPosition != nil || !b50FCandidateReflect.DeepEqual(b50FCandidateStack.B50FCandidateConflict, []string{"stack-voltage label", "temperature conversion metadata"}) || !b50FCandidateReflect.DeepEqual(b50FCandidateStack.B50FCandidateExclusions, []string{"voltage scale", "selected temperature conversion"}) || b50FCandidateRaw != [2]byte{16, 0} {
		b50FCandidateT.Fatal("stack conversion conflict resolved or raw changed")
	}
}

func TestAb63IntegrityAndFeedback(b50FCandidateT *b50FCandidateTesting.T) {
	b50FCandidateFeedback := b50FCandidateGet(b50FCandidateT, nil, "feedback")
	if b50FCandidateFeedback.B50FCandidateInterpretation != "acceptance feedback only" || !b50FCandidateReflect.DeepEqual(b50FCandidateFeedback.B50FCandidateExclusions, []string{"integrity validation", "physical effect"}) || len(b50FCandidateFeedback.B50FCandidateField.SchemaEnum) != 256 {
		b50FCandidateT.Fatal("feedback boundaries")
	}
	for b50FCandidateValue := 0; b50FCandidateValue < 256; b50FCandidateValue++ {
		b50FCandidateWant := "rejected"
		if b50FCandidateValue == 1 {
			b50FCandidateWant = "accepted"
		}
		if b50FCandidateFeedback.B50FCandidateField.SchemaEnum[int64(b50FCandidateValue)] != b50FCandidateWant {
			b50FCandidateT.Fatal("feedback domain")
		}
	}
	for b50FCandidateOp, b50FCandidatePositions := range map[byte][2]int{0: {12, 13}, 1: {12, 13}, 2: {11, 12}, 3: {4, 5}} {
		b50FCandidateFirst := b50FCandidateGet(b50FCandidateT, b50FCandidateByte(b50FCandidateOp), "integrity_1")
		b50FCandidateSecond := b50FCandidateGet(b50FCandidateT, b50FCandidateByte(b50FCandidateOp), "integrity_2")
		for b50FCandidateIndex, b50FCandidateEntry := range []B50FCandidateEntry{b50FCandidateFirst, b50FCandidateSecond} {
			if b50FCandidateEntry.B50FCandidateField.SchemaPosition == nil || *b50FCandidateEntry.B50FCandidateField.SchemaPosition != b50FCandidatePositions[b50FCandidateIndex] || b50FCandidateEntry.B50FCandidateField.SchemaWidth == nil || *b50FCandidateEntry.B50FCandidateField.SchemaWidth != 1 || b50FCandidateEntry.B50FCandidateField.SchemaCodec != "raw" || b50FCandidateEntry.B50FCandidateInterpretation != "independent raw integrity byte" || b50FCandidateEntry.B50FCandidateCoordinates != "abstract logical byte position; not established wire offset" || !b50FCandidateReflect.DeepEqual(b50FCandidateEntry.B50FCandidateExclusions, []string{"word byte order: Unknown", "polynomial: Unknown", "initialization: Unknown", "reflection: Unknown", "final transform: Unknown"}) {
				b50FCandidateT.Fatal("integrity interpretation")
			}
		}
		for _, b50FCandidateLeft := range []byte{0, 18, 255} {
			for _, b50FCandidateRight := range []byte{0, 52, 255} {
				var b50FCandidateRaw [14]byte
				b50FCandidateRaw[b50FCandidatePositions[0]] = b50FCandidateLeft
				b50FCandidateRaw[b50FCandidatePositions[1]] = b50FCandidateRight
				if [2]byte{b50FCandidateRaw[*b50FCandidateFirst.B50FCandidateField.SchemaPosition], b50FCandidateRaw[*b50FCandidateSecond.B50FCandidateField.SchemaPosition]} != [2]byte{b50FCandidateLeft, b50FCandidateRight} {
					b50FCandidateT.Fatal("ordered independent bytes")
				}
				for b50FCandidateIndex, b50FCandidateValue := range []int64{0, 1, 2, 255} {
					if b50FCandidateFeedback.B50FCandidateField.SchemaEnum[b50FCandidateValue] != []string{"rejected", "accepted", "rejected", "rejected"}[b50FCandidateIndex] {
						b50FCandidateT.Fatal("feedback coupled to integrity")
					}
				}
			}
		}
	}
}

func TestAb63DeepCopy(b50FCandidateT *b50FCandidateTesting.T) {
	b50FCandidateBefore := B50FCandidates()
	b50FCandidateMutated := B50FCandidates()
	b50FCandidateMutated.B50FCandidateLimitations[0] = "ab63_mutated"
	for b50FCandidateIndex := range b50FCandidateMutated.B50FCandidateEntries {
		b50FCandidateEntry := &b50FCandidateMutated.B50FCandidateEntries[b50FCandidateIndex]
		if b50FCandidateEntry.B50FCandidateKey.B50FCandidateSecondary != nil {
			*b50FCandidateEntry.B50FCandidateKey.B50FCandidateSecondary = 0
		}
		if b50FCandidateEntry.B50FCandidateKey.B50FCandidateOperation != nil {
			*b50FCandidateEntry.B50FCandidateKey.B50FCandidateOperation = 255
		}
		if b50FCandidateEntry.B50FCandidateField.SchemaPosition != nil {
			*b50FCandidateEntry.B50FCandidateField.SchemaPosition = 99
		}
		if b50FCandidateEntry.B50FCandidateField.SchemaWidth != nil {
			*b50FCandidateEntry.B50FCandidateField.SchemaWidth = 99
		}
		if b50FCandidateEntry.B50FCandidateDivisor != nil {
			*b50FCandidateEntry.B50FCandidateDivisor = 99
		}
		if b50FCandidateEntry.B50FCandidateReserved != nil {
			*b50FCandidateEntry.B50FCandidateReserved = 0
		}
		for b50FCandidateName := range b50FCandidateEntry.B50FCandidateMasks {
			b50FCandidateEntry.B50FCandidateMasks[b50FCandidateName] = 0
		}
		for b50FCandidateValue := range b50FCandidateEntry.B50FCandidateField.SchemaEnum {
			b50FCandidateEntry.B50FCandidateField.SchemaEnum[b50FCandidateValue] = "ab63_mutated"
		}
		if len(b50FCandidateEntry.B50FCandidateConflict) > 0 {
			b50FCandidateEntry.B50FCandidateConflict[0] = "ab63_mutated"
		}
		if len(b50FCandidateEntry.B50FCandidateExclusions) > 0 {
			b50FCandidateEntry.B50FCandidateExclusions[0] = "ab63_mutated"
		}
	}
	if !b50FCandidateReflect.DeepEqual(b50FCandidateBefore, B50FCandidates()) {
		b50FCandidateT.Fatal("accessor aliases nested state")
	}
	for _, b50FCandidateEntry := range b50FCandidateBefore.B50FCandidateEntries {
		b50FCandidateReturned := b50FCandidateGet(b50FCandidateT, b50FCandidateEntry.B50FCandidateKey.B50FCandidateOperation, b50FCandidateEntry.B50FCandidateKey.B50FCandidateField)
		*b50FCandidateReturned.B50FCandidateField.SchemaWidth = 99
		if b50FCandidateAgain, b50FCandidateOK := B50FCandidateLookup(b50FCandidateEntry.B50FCandidateKey); !b50FCandidateOK || !b50FCandidateReflect.DeepEqual(b50FCandidateAgain, b50FCandidateEntry) {
			b50FCandidateT.Fatal("lookup aliases state")
		}
	}
}
