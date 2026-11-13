package vaillant

import (
	b51ADiagnosticVariantJSON "encoding/json"
	b51ADiagnosticVariantReflect "reflect"
	b51ADiagnosticVariantStrings "strings"
	b51ADiagnosticVariantTesting "testing"
)

func b51ADiagnosticVariantPointer[b51ADiagnosticVariantType any](b51ADiagnosticVariantValue b51ADiagnosticVariantType) *b51ADiagnosticVariantType {
	return &b51ADiagnosticVariantValue
}

func b51ADiagnosticVariantGet(b51ADiagnosticVariantT *b51ADiagnosticVariantTesting.T, b51ADiagnosticVariantOp uint8, b51ADiagnosticVariantDirection, b51ADiagnosticVariantField string) B51ADiagnosticVariantEntry {
	b51ADiagnosticVariantT.Helper()
	b51ADiagnosticVariantEntry, b51ADiagnosticVariantFound, b51ADiagnosticVariantErr := B51ADiagnosticVariantLookup(B51ADiagnosticVariantKey{B51ADiagnosticVariantFamily: "vaillant", B51ADiagnosticVariant: "B51A", B51ADiagnosticVariantPrimary: 181, B51ADiagnosticVariantSecondary: b51ADiagnosticVariantPointer(uint8(26)), B51ADiagnosticVariantOperation: b51ADiagnosticVariantPointer(b51ADiagnosticVariantOp), B51ADiagnosticVariantDirection: b51ADiagnosticVariantDirection, B51ADiagnosticVariantField: b51ADiagnosticVariantField})
	if b51ADiagnosticVariantErr != nil || !b51ADiagnosticVariantFound {
		b51ADiagnosticVariantT.Fatalf("lookup %d/%s/%s: %v %v", b51ADiagnosticVariantOp, b51ADiagnosticVariantDirection, b51ADiagnosticVariantField, b51ADiagnosticVariantFound, b51ADiagnosticVariantErr)
	}
	return b51ADiagnosticVariantEntry
}

func b51ADiagnosticVariantEqual(b51ADiagnosticVariantT *b51ADiagnosticVariantTesting.T, b51ADiagnosticVariantGot, b51ADiagnosticVariantWant any) {
	b51ADiagnosticVariantT.Helper()
	if !b51ADiagnosticVariantReflect.DeepEqual(b51ADiagnosticVariantGot, b51ADiagnosticVariantWant) {
		b51ADiagnosticVariantT.Fatalf("got %#v, want %#v", b51ADiagnosticVariantGot, b51ADiagnosticVariantWant)
	}
}

func TestAb98CoordinatesAndGrades(b51ADiagnosticVariantT *b51ADiagnosticVariantTesting.T) {
	b51ADiagnosticVariantEntries, b51ADiagnosticVariantErr := B51ADiagnosticVariantEntries()
	if b51ADiagnosticVariantErr != nil {
		b51ADiagnosticVariantT.Fatal(b51ADiagnosticVariantErr)
	}
	b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, len(b51ADiagnosticVariantEntries), 37)
	for _, b51ADiagnosticVariantEntry := range b51ADiagnosticVariantEntries {
		b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantEntry.B51ADiagnosticVariantGrade, "Hypothesis")
		b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantEntry.B51ADiagnosticVariantMetadataPresence, "Proven")
		b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantEntry.B51ADiagnosticVariantPhysicalQualification, "Unknown")
		b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantEntry.B51ADiagnosticVariantInvocation, "Unknown")
		b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantEntry.B51ADiagnosticVariantNativeContract, "https://github.com/Project-Helianthus/helianthus-docs-ebus/blob/main/protocols/ebus-services/ebus-overview.md")
	}
	for _, b51ADiagnosticVariantOp := range []uint8{1, 4, 5, 7, 8} {
		for _, b51ADiagnosticVariantDirection := range []string{"request", "positive"} {
			b51ADiagnosticVariantEntry := b51ADiagnosticVariantGet(b51ADiagnosticVariantT, b51ADiagnosticVariantOp, b51ADiagnosticVariantDirection, "sequence")
			b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantEntry.B51ADiagnosticVariantPosition, b51ADiagnosticVariantPointer(uint16(3)))
			b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantEntry.B51ADiagnosticVariantWidthBits, b51ADiagnosticVariantPointer(uint16(8)))
			b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantEntry.B51ADiagnosticVariantDomain, []uint64{0, 255})
			b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantEntry.B51ADiagnosticVariantDefault, b51ADiagnosticVariantPointer(uint64(0)))
		}
	}
}

func TestAb98KeyAxes(b51ADiagnosticVariantT *b51ADiagnosticVariantTesting.T) {
	b51ADiagnosticVariantBase := b51ADiagnosticVariantGet(b51ADiagnosticVariantT, 1, "request", "CRC").B51ADiagnosticVariantKey
	b51ADiagnosticVariantChanges := []struct {
		b51ADiagnosticVariantName   string
		b51ADiagnosticVariantChange func(*B51ADiagnosticVariantKey)
	}{
		{"family", func(b51ADiagnosticVariantKey *B51ADiagnosticVariantKey) {
			b51ADiagnosticVariantKey.B51ADiagnosticVariantFamily = "other-family"
		}},
		{"variant", func(b51ADiagnosticVariantKey *B51ADiagnosticVariantKey) {
			b51ADiagnosticVariantKey.B51ADiagnosticVariant = "B524"
		}},
		{"primary", func(b51ADiagnosticVariantKey *B51ADiagnosticVariantKey) {
			b51ADiagnosticVariantKey.B51ADiagnosticVariantPrimary = 0
		}},
		{"secondary", func(b51ADiagnosticVariantKey *B51ADiagnosticVariantKey) {
			b51ADiagnosticVariantKey.B51ADiagnosticVariantSecondary = b51ADiagnosticVariantPointer(uint8(0))
		}},
		{"operation", func(b51ADiagnosticVariantKey *B51ADiagnosticVariantKey) {
			b51ADiagnosticVariantKey.B51ADiagnosticVariantOperation = b51ADiagnosticVariantPointer(uint8(0))
		}},
		{"direction", func(b51ADiagnosticVariantKey *B51ADiagnosticVariantKey) {
			b51ADiagnosticVariantKey.B51ADiagnosticVariantDirection = "positive"
		}},
		{"field", func(b51ADiagnosticVariantKey *B51ADiagnosticVariantKey) {
			b51ADiagnosticVariantKey.B51ADiagnosticVariantField = "unlisted"
		}},
	}
	for _, b51ADiagnosticVariantCase := range b51ADiagnosticVariantChanges {
		b51ADiagnosticVariantT.Run(b51ADiagnosticVariantCase.b51ADiagnosticVariantName, func(b51ADiagnosticVariantT *b51ADiagnosticVariantTesting.T) {
			b51ADiagnosticVariantKey := b51ADiagnosticVariantBase
			b51ADiagnosticVariantCase.b51ADiagnosticVariantChange(&b51ADiagnosticVariantKey)
			_, b51ADiagnosticVariantFound, b51ADiagnosticVariantErr := B51ADiagnosticVariantLookup(b51ADiagnosticVariantKey)
			if b51ADiagnosticVariantFound || b51ADiagnosticVariantErr != nil {
				b51ADiagnosticVariantT.Fatalf("nonmatching key: %v %v", b51ADiagnosticVariantFound, b51ADiagnosticVariantErr)
			}
		})
	}
	b51ADiagnosticVariantOriginal := b51ADiagnosticVariantData
	defer func() { b51ADiagnosticVariantData = b51ADiagnosticVariantOriginal }()
	b51ADiagnosticVariantSynthetic := b51ADiagnosticVariantGet(b51ADiagnosticVariantT, 1, "request", "flags")
	b51ADiagnosticVariantSynthetic.B51ADiagnosticVariantKey = B51ADiagnosticVariantKey{B51ADiagnosticVariantFamily: "synthetic-family", B51ADiagnosticVariant: "B524", B51ADiagnosticVariantPrimary: 181, B51ADiagnosticVariantSecondary: b51ADiagnosticVariantPointer(uint8(36)), B51ADiagnosticVariantOperation: b51ADiagnosticVariantPointer(uint8(3)), B51ADiagnosticVariantDirection: "request", B51ADiagnosticVariantField: "synthetic-field", B51ADiagnosticVariantGroup: b51ADiagnosticVariantPointer(uint8(7)), B51ADiagnosticVariantRequestedInstance: b51ADiagnosticVariantPointer(uint8(9)), B51ADiagnosticVariantRegister: b51ADiagnosticVariantPointer(uint16(513)), B51ADiagnosticVariantCommandRole: b51ADiagnosticVariantPointer("read")}
	b51ADiagnosticVariantChanges[5].b51ADiagnosticVariantChange = func(b51ADiagnosticVariantKey *B51ADiagnosticVariantKey) {
		b51ADiagnosticVariantKey.B51ADiagnosticVariantDirection = "positive"
	}
	b51ADiagnosticVariantChanges[1].b51ADiagnosticVariantChange = func(b51ADiagnosticVariantKey *B51ADiagnosticVariantKey) {
		b51ADiagnosticVariantKey.B51ADiagnosticVariant = "B51A"
	}
	b51ADiagnosticVariantChanges = append(b51ADiagnosticVariantChanges,
		struct {
			b51ADiagnosticVariantName   string
			b51ADiagnosticVariantChange func(*B51ADiagnosticVariantKey)
		}{"group", func(b51ADiagnosticVariantKey *B51ADiagnosticVariantKey) {
			b51ADiagnosticVariantKey.B51ADiagnosticVariantGroup = b51ADiagnosticVariantPointer(uint8(0))
		}},
		struct {
			b51ADiagnosticVariantName   string
			b51ADiagnosticVariantChange func(*B51ADiagnosticVariantKey)
		}{"requested-instance", func(b51ADiagnosticVariantKey *B51ADiagnosticVariantKey) {
			b51ADiagnosticVariantKey.B51ADiagnosticVariantRequestedInstance = b51ADiagnosticVariantPointer(uint8(0))
		}},
		struct {
			b51ADiagnosticVariantName   string
			b51ADiagnosticVariantChange func(*B51ADiagnosticVariantKey)
		}{"register", func(b51ADiagnosticVariantKey *B51ADiagnosticVariantKey) {
			b51ADiagnosticVariantKey.B51ADiagnosticVariantRegister = b51ADiagnosticVariantPointer(uint16(257))
		}},
		struct {
			b51ADiagnosticVariantName   string
			b51ADiagnosticVariantChange func(*B51ADiagnosticVariantKey)
		}{"command-role", func(b51ADiagnosticVariantKey *B51ADiagnosticVariantKey) {
			b51ADiagnosticVariantKey.B51ADiagnosticVariantCommandRole = b51ADiagnosticVariantPointer("write")
		}},
	)
	b51ADiagnosticVariantInput, b51ADiagnosticVariantErr := b51ADiagnosticVariantJSON.Marshal([]B51ADiagnosticVariantEntry{b51ADiagnosticVariantSynthetic})
	if b51ADiagnosticVariantErr != nil {
		b51ADiagnosticVariantT.Fatal(b51ADiagnosticVariantErr)
	}
	b51ADiagnosticVariantSetData(b51ADiagnosticVariantInput)
	b51ADiagnosticVariantHit, b51ADiagnosticVariantFound, b51ADiagnosticVariantErr := B51ADiagnosticVariantLookup(b51ADiagnosticVariantSynthetic.B51ADiagnosticVariantKey)
	if b51ADiagnosticVariantErr != nil || !b51ADiagnosticVariantFound {
		b51ADiagnosticVariantT.Fatalf("synthetic baseline: %v %v", b51ADiagnosticVariantFound, b51ADiagnosticVariantErr)
	}
	b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantHit, b51ADiagnosticVariantSynthetic)
	for _, b51ADiagnosticVariantCase := range b51ADiagnosticVariantChanges {
		b51ADiagnosticVariantT.Run("B524/"+b51ADiagnosticVariantCase.b51ADiagnosticVariantName, func(b51ADiagnosticVariantT *b51ADiagnosticVariantTesting.T) {
			b51ADiagnosticVariantKey := b51ADiagnosticVariantSynthetic.B51ADiagnosticVariantKey
			b51ADiagnosticVariantCase.b51ADiagnosticVariantChange(&b51ADiagnosticVariantKey)
			_, b51ADiagnosticVariantFound, b51ADiagnosticVariantErr := B51ADiagnosticVariantLookup(b51ADiagnosticVariantKey)
			if b51ADiagnosticVariantFound || b51ADiagnosticVariantErr != nil {
				b51ADiagnosticVariantT.Fatalf("axis ignored: %v %v", b51ADiagnosticVariantFound, b51ADiagnosticVariantErr)
			}
		})
	}
	for _, b51ADiagnosticVariantAxis := range []string{"secondary", "operation", "group", "requested-instance", "register", "command-role"} {
		b51ADiagnosticVariantT.Run("absence/"+b51ADiagnosticVariantAxis, func(b51ADiagnosticVariantT *b51ADiagnosticVariantTesting.T) {
			b51ADiagnosticVariantNil := b51ADiagnosticVariantSynthetic
			b51ADiagnosticVariantZero := b51ADiagnosticVariantSynthetic
			switch b51ADiagnosticVariantAxis {
			case "secondary":
				b51ADiagnosticVariantNil.B51ADiagnosticVariantKey.B51ADiagnosticVariantSecondary = nil
				b51ADiagnosticVariantZero.B51ADiagnosticVariantKey.B51ADiagnosticVariantSecondary = b51ADiagnosticVariantPointer(uint8(0))
			case "operation":
				b51ADiagnosticVariantNil.B51ADiagnosticVariantKey.B51ADiagnosticVariantOperation = nil
				b51ADiagnosticVariantZero.B51ADiagnosticVariantKey.B51ADiagnosticVariantOperation = b51ADiagnosticVariantPointer(uint8(0))
			case "group":
				b51ADiagnosticVariantNil.B51ADiagnosticVariantKey.B51ADiagnosticVariantGroup = nil
				b51ADiagnosticVariantZero.B51ADiagnosticVariantKey.B51ADiagnosticVariantGroup = b51ADiagnosticVariantPointer(uint8(0))
			case "requested-instance":
				b51ADiagnosticVariantNil.B51ADiagnosticVariantKey.B51ADiagnosticVariantRequestedInstance = nil
				b51ADiagnosticVariantZero.B51ADiagnosticVariantKey.B51ADiagnosticVariantRequestedInstance = b51ADiagnosticVariantPointer(uint8(0))
			case "register":
				b51ADiagnosticVariantNil.B51ADiagnosticVariantKey.B51ADiagnosticVariantRegister = nil
				b51ADiagnosticVariantZero.B51ADiagnosticVariantKey.B51ADiagnosticVariantRegister = b51ADiagnosticVariantPointer(uint16(0))
			case "command-role":
				b51ADiagnosticVariantNil.B51ADiagnosticVariantKey.B51ADiagnosticVariantCommandRole = nil
				b51ADiagnosticVariantZero.B51ADiagnosticVariantKey.B51ADiagnosticVariantCommandRole = b51ADiagnosticVariantPointer("")
			}
			b51ADiagnosticVariantNil.B51ADiagnosticVariantInterpretations = []string{"absent"}
			b51ADiagnosticVariantZero.B51ADiagnosticVariantInterpretations = []string{"zero"}
			b51ADiagnosticVariantInput, b51ADiagnosticVariantErr := b51ADiagnosticVariantJSON.Marshal([]B51ADiagnosticVariantEntry{b51ADiagnosticVariantNil, b51ADiagnosticVariantZero})
			if b51ADiagnosticVariantErr != nil {
				b51ADiagnosticVariantT.Fatal(b51ADiagnosticVariantErr)
			}
			b51ADiagnosticVariantSetData(b51ADiagnosticVariantInput)
			for _, b51ADiagnosticVariantWant := range []B51ADiagnosticVariantEntry{b51ADiagnosticVariantNil, b51ADiagnosticVariantZero} {
				b51ADiagnosticVariantGot, b51ADiagnosticVariantFound, b51ADiagnosticVariantErr := B51ADiagnosticVariantLookup(b51ADiagnosticVariantWant.B51ADiagnosticVariantKey)
				if b51ADiagnosticVariantErr != nil || !b51ADiagnosticVariantFound {
					b51ADiagnosticVariantT.Fatalf("absence: %v %v", b51ADiagnosticVariantFound, b51ADiagnosticVariantErr)
				}
				b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantGot, b51ADiagnosticVariantWant)
			}
		})
	}
}

func b51ADiagnosticVariantSetData(b51ADiagnosticVariantInput []byte) {
	b51ADiagnosticVariantData = b51ADiagnosticVariantInput
}

func TestAb98FlagsAndConditionalFields(b51ADiagnosticVariantT *b51ADiagnosticVariantTesting.T) {
	for _, b51ADiagnosticVariantCase := range []struct {
		b51ADiagnosticVariantOp        uint8
		b51ADiagnosticVariantDirection string
		b51ADiagnosticVariantLabels    map[uint8][]string
	}{
		{1, "request", map[uint8][]string{0: {"power-on/start"}, 1: {"change"}, 2: {"CRC-support"}}},
		{1, "positive", map[uint8][]string{0: {"valid"}, 1: {"invalid"}, 2: {"blocked"}, 3: {"change"}, 4: {"denied"}, 5: {"invalid-range"}, 6: {"not-user-adjustable"}}},
		{5, "positive", map[uint8][]string{0: {"disabled; following data omitted"}, 1: {"writable"}, 2: {"enum-vs-range"}, 3: {"additional-four-byte/no-range", "read-only over writable"}, 6: {"text alternative 6"}, 7: {"text alternative 7"}}},
		{7, "request", map[uint8][]string{0: {"normal-replication"}, 1: {"inverted-replication"}, 2: {"refresh"}}},
	} {
		b51ADiagnosticVariantEntry := b51ADiagnosticVariantGet(b51ADiagnosticVariantT, b51ADiagnosticVariantCase.b51ADiagnosticVariantOp, b51ADiagnosticVariantCase.b51ADiagnosticVariantDirection, "flags")
		b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, len(b51ADiagnosticVariantEntry.B51ADiagnosticVariantBits), 8)
		b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantEntry.B51ADiagnosticVariantExclusive, b51ADiagnosticVariantPointer(false))
		for b51ADiagnosticVariantIndex := uint8(0); b51ADiagnosticVariantIndex < 8; b51ADiagnosticVariantIndex++ {
			b51ADiagnosticVariantBit := b51ADiagnosticVariantEntry.B51ADiagnosticVariantBits[b51ADiagnosticVariantIndex]
			b51ADiagnosticVariantLabels, b51ADiagnosticVariantNamed := b51ADiagnosticVariantCase.b51ADiagnosticVariantLabels[b51ADiagnosticVariantIndex]
			b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantBit.B51ADiagnosticVariantBit, b51ADiagnosticVariantIndex)
			b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantBit.B51ADiagnosticVariantReserved, !b51ADiagnosticVariantNamed)
			b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantBit.B51ADiagnosticVariantDefault, false)
			if b51ADiagnosticVariantNamed {
				b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantBit.B51ADiagnosticVariantLabels, b51ADiagnosticVariantLabels)
			} else {
				b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, len(b51ADiagnosticVariantBit.B51ADiagnosticVariantLabels), 0)
			}
			b51ADiagnosticVariantRaw := uint8(1) << b51ADiagnosticVariantIndex
			b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantRaw&(1<<b51ADiagnosticVariantBit.B51ADiagnosticVariantBit), b51ADiagnosticVariantRaw)
		}
	}
	for b51ADiagnosticVariantIndex, b51ADiagnosticVariantField := range []string{"DSN", "start", "offset"} {
		for _, b51ADiagnosticVariantDirection := range []string{"request", "positive"} {
			b51ADiagnosticVariantEntry := b51ADiagnosticVariantGet(b51ADiagnosticVariantT, 1, b51ADiagnosticVariantDirection, b51ADiagnosticVariantField)
			b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantEntry.B51ADiagnosticVariantPosition, b51ADiagnosticVariantPointer(uint16(5+2*b51ADiagnosticVariantIndex)))
			b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantEntry.B51ADiagnosticVariantWidthBits, (*uint16)(nil))
			b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantEntry.B51ADiagnosticVariantByteOrder, (*string)(nil))
			if b51ADiagnosticVariantDirection == "positive" {
				b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantEntry.B51ADiagnosticVariantConditions, []B51ADiagnosticVariantCondition{{B51ADiagnosticVariantField: "flags", B51ADiagnosticVariantBit: 3, B51ADiagnosticVariantSet: true}})
				for _, b51ADiagnosticVariantRaw := range []uint8{0, 8, 255} {
					b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantApplies(b51ADiagnosticVariantEntry.B51ADiagnosticVariantConditions, b51ADiagnosticVariantRaw), b51ADiagnosticVariantRaw&8 != 0)
				}
			} else {
				b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, len(b51ADiagnosticVariantEntry.B51ADiagnosticVariantConditions), 0)
			}
		}
	}
	b51ADiagnosticVariantCRC := b51ADiagnosticVariantGet(b51ADiagnosticVariantT, 1, "request", "CRC")
	b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantCRC.B51ADiagnosticVariantDefault, (*uint64)(nil))
	b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantGet(b51ADiagnosticVariantT, 1, "positive", "family-code").B51ADiagnosticVariantDefault, (*uint64)(nil))
	b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantCRC.B51ADiagnosticVariantPosition, b51ADiagnosticVariantPointer(uint16(10)))
	b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantCRC.B51ADiagnosticVariantWidthBits, b51ADiagnosticVariantPointer(uint16(16)))
	b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantCRC.B51ADiagnosticVariantByteOrder, b51ADiagnosticVariantPointer("little-endian"))
	b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantCRC.B51ADiagnosticVariantInterpretations, []string{"Mirror candidate. Unsupported CRC-support flag is 0."})
	b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantCRC.B51ADiagnosticVariantUnknown, []string{"replacement number", "physical length", "qualification", "invocation"})
	b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantGet(b51ADiagnosticVariantT, 1, "positive", "family-code").B51ADiagnosticVariantPosition, b51ADiagnosticVariantPointer(uint16(10)))
}

func b51ADiagnosticVariantApplies(b51ADiagnosticVariantConditions []B51ADiagnosticVariantCondition, b51ADiagnosticVariantRaw uint8) bool {
	for _, b51ADiagnosticVariantCondition := range b51ADiagnosticVariantConditions {
		if (b51ADiagnosticVariantRaw&(1<<b51ADiagnosticVariantCondition.B51ADiagnosticVariantBit) != 0) != b51ADiagnosticVariantCondition.B51ADiagnosticVariantSet {
			return false
		}
	}
	return true
}

func TestAb98DiagnosticAlternatives(b51ADiagnosticVariantT *b51ADiagnosticVariantTesting.T) {
	b51ADiagnosticVariantEnable := b51ADiagnosticVariantGet(b51ADiagnosticVariantT, 4, "positive", "enable-bits")
	b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantEnable.B51ADiagnosticVariantPosition, b51ADiagnosticVariantPointer(uint16(4)))
	b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantEnable.B51ADiagnosticVariantWidthBits, b51ADiagnosticVariantPointer(uint16(8)))
	b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantEnable.B51ADiagnosticVariantCount, b51ADiagnosticVariantPointer(uint16(13)))
	b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantEnable.B51ADiagnosticVariantLogicalFootprint, b51ADiagnosticVariantPointer(uint16(17)))
	b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantEnable.B51ADiagnosticVariantTail, []uint16{0, 13, 0})
	b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, int(*b51ADiagnosticVariantEnable.B51ADiagnosticVariantCount)*int(*b51ADiagnosticVariantEnable.B51ADiagnosticVariantWidthBits), 104)
	for b51ADiagnosticVariantIndex := 0; b51ADiagnosticVariantIndex < 104; b51ADiagnosticVariantIndex++ {
		b51ADiagnosticVariantRaw := make([]byte, 17)
		b51ADiagnosticVariantRaw[4+b51ADiagnosticVariantIndex/8] = 1 << (b51ADiagnosticVariantIndex % 8)
		b51ADiagnosticVariantSeen := 0
		for b51ADiagnosticVariantByte := 0; b51ADiagnosticVariantByte < int(*b51ADiagnosticVariantEnable.B51ADiagnosticVariantCount); b51ADiagnosticVariantByte++ {
			for b51ADiagnosticVariantBit := 0; b51ADiagnosticVariantBit < int(*b51ADiagnosticVariantEnable.B51ADiagnosticVariantWidthBits); b51ADiagnosticVariantBit++ {
				if b51ADiagnosticVariantRaw[int(*b51ADiagnosticVariantEnable.B51ADiagnosticVariantPosition)+b51ADiagnosticVariantByte]&(1<<b51ADiagnosticVariantBit) != 0 {
					b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantByte*8+b51ADiagnosticVariantBit, b51ADiagnosticVariantIndex)
					b51ADiagnosticVariantSeen++
				}
			}
		}
		b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantSeen, 1)
	}
	b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantEnable.B51ADiagnosticVariantValues, map[uint64]string{0: "disabled", 1: "enabled"})
	b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantEnable.B51ADiagnosticVariantDefault, b51ADiagnosticVariantPointer(uint64(0)))
	for _, b51ADiagnosticVariantOp := range []uint8{4, 5} {
		b51ADiagnosticVariantLevel := b51ADiagnosticVariantGet(b51ADiagnosticVariantT, b51ADiagnosticVariantOp, "request", "level")
		b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantLevel.B51ADiagnosticVariantPosition, b51ADiagnosticVariantPointer(uint16(4)))
		if b51ADiagnosticVariantOp == 4 {
			b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantLevel.B51ADiagnosticVariantWidthBits, b51ADiagnosticVariantPointer(uint16(8)))
		} else {
			b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantLevel.B51ADiagnosticVariantWidthBits, (*uint16)(nil))
			b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantLevel.B51ADiagnosticVariantDomain, []uint64(nil))
		}
		b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantLevel.B51ADiagnosticVariantDefault, b51ADiagnosticVariantPointer(uint64(0)))
	}
	b51ADiagnosticVariantPoint := b51ADiagnosticVariantGet(b51ADiagnosticVariantT, 5, "request", "point")
	b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantPoint.B51ADiagnosticVariantWidthBits, (*uint16)(nil))
	b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantPoint.B51ADiagnosticVariantPosition, b51ADiagnosticVariantPointer(uint16(5)))
	b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantPoint.B51ADiagnosticVariantDomain, []uint64{0, 99})
	b51ADiagnosticVariantFollowing := b51ADiagnosticVariantGet(b51ADiagnosticVariantT, 5, "positive", "following-data")
	b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantFollowing.B51ADiagnosticVariantConditions, []B51ADiagnosticVariantCondition{{B51ADiagnosticVariantField: "flags", B51ADiagnosticVariantBit: 0, B51ADiagnosticVariantSet: false}})
	b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, len(b51ADiagnosticVariantFollowing.B51ADiagnosticVariantAlternatives), 3)
	for b51ADiagnosticVariantRaw := 0; b51ADiagnosticVariantRaw <= 255; b51ADiagnosticVariantRaw++ {
		b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantApplies(b51ADiagnosticVariantFollowing.B51ADiagnosticVariantConditions, uint8(b51ADiagnosticVariantRaw)), b51ADiagnosticVariantRaw&1 == 0)
	}
	b51ADiagnosticVariantRange := b51ADiagnosticVariantFollowing.B51ADiagnosticVariantAlternatives[0]
	b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantRange.B51ADiagnosticVariantLabel, "range")
	b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, len(b51ADiagnosticVariantRange.B51ADiagnosticVariantFields), 5)
	b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantRange.B51ADiagnosticVariantInterpretations, []string{"Candidate min+n*step with n*step<=max; enforcement Unknown."})
	for b51ADiagnosticVariantIndex, b51ADiagnosticVariantWant := range []struct {
		b51ADiagnosticVariantName     string
		b51ADiagnosticVariantPosition uint16
		b51ADiagnosticVariantWidth    *uint16
		b51ADiagnosticVariantOrder    *string
	}{
		{"display-type", 5, nil, nil}, {"value", 6, b51ADiagnosticVariantPointer(uint16(16)), b51ADiagnosticVariantPointer("little-endian")}, {"minimum", 8, nil, nil}, {"maximum", 10, nil, nil}, {"step", 12, b51ADiagnosticVariantPointer(uint16(8)), nil},
	} {
		b51ADiagnosticVariantGot := b51ADiagnosticVariantRange.B51ADiagnosticVariantFields[b51ADiagnosticVariantIndex]
		b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantGot.B51ADiagnosticVariantDefault, b51ADiagnosticVariantPointer(uint64(0)))
		b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantGot.B51ADiagnosticVariantKey.B51ADiagnosticVariantField, b51ADiagnosticVariantWant.b51ADiagnosticVariantName)
		b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantGot.B51ADiagnosticVariantPosition, b51ADiagnosticVariantPointer(b51ADiagnosticVariantWant.b51ADiagnosticVariantPosition))
		b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantGot.B51ADiagnosticVariantWidthBits, b51ADiagnosticVariantWant.b51ADiagnosticVariantWidth)
		b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantGot.B51ADiagnosticVariantByteOrder, b51ADiagnosticVariantWant.b51ADiagnosticVariantOrder)
	}
	b51ADiagnosticVariantEnumRaw := []byte{0, 0, 0, 0, 4, 0, 11, 23, 37, 41, 59, 67}
	b51ADiagnosticVariantEnum := b51ADiagnosticVariantFollowing.B51ADiagnosticVariantAlternatives[1]
	b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantEnum.B51ADiagnosticVariantLabel, "enum")
	b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, len(b51ADiagnosticVariantEnum.B51ADiagnosticVariantFields), 6)
	b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantEnum.B51ADiagnosticVariantInterpretations, []string{"Six independent alternatives; repeat last when fewer. No separate current-value field."})
	for b51ADiagnosticVariantIndex, b51ADiagnosticVariantName := range []string{"enum-alternative-1", "enum-alternative-2", "enum-alternative-3", "enum-alternative-4", "enum-alternative-5", "enum-alternative-6"} {
		b51ADiagnosticVariantField := b51ADiagnosticVariantEnum.B51ADiagnosticVariantFields[b51ADiagnosticVariantIndex]
		b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantEnumRaw[*b51ADiagnosticVariantField.B51ADiagnosticVariantPosition], []byte{11, 23, 37, 41, 59, 67}[b51ADiagnosticVariantIndex])
		b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantField.B51ADiagnosticVariantKey.B51ADiagnosticVariantField, b51ADiagnosticVariantName)
		b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantField.B51ADiagnosticVariantPosition, b51ADiagnosticVariantPointer(uint16(6+b51ADiagnosticVariantIndex)))
		b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantField.B51ADiagnosticVariantWidthBits, b51ADiagnosticVariantPointer(uint16(8)))
	}
	b51ADiagnosticVariantFour := b51ADiagnosticVariantFollowing.B51ADiagnosticVariantAlternatives[2]
	b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantFour.B51ADiagnosticVariantLabel, "four-byte")
	b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, len(b51ADiagnosticVariantFour.B51ADiagnosticVariantFields), 4)
	b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantFour.B51ADiagnosticVariantFields[0].B51ADiagnosticVariantKey.B51ADiagnosticVariantField, "additional-value")
	b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantFour.B51ADiagnosticVariantFields[0].B51ADiagnosticVariantPosition, b51ADiagnosticVariantPointer(uint16(6)))
	b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantFour.B51ADiagnosticVariantFields[0].B51ADiagnosticVariantWidthBits, b51ADiagnosticVariantPointer(uint16(32)))
	b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantFour.B51ADiagnosticVariantFields[0].B51ADiagnosticVariantByteOrder, b51ADiagnosticVariantPointer("little-endian"))
	for b51ADiagnosticVariantIndex, b51ADiagnosticVariantName := range []string{"reserved-10", "reserved-11", "reserved-12"} {
		b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantFour.B51ADiagnosticVariantFields[b51ADiagnosticVariantIndex+1].B51ADiagnosticVariantKey.B51ADiagnosticVariantField, b51ADiagnosticVariantName)
		b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantFour.B51ADiagnosticVariantFields[b51ADiagnosticVariantIndex+1].B51ADiagnosticVariantPosition, b51ADiagnosticVariantPointer(uint16(10+b51ADiagnosticVariantIndex)))
		b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantFour.B51ADiagnosticVariantFields[b51ADiagnosticVariantIndex+1].B51ADiagnosticVariantWidthBits, b51ADiagnosticVariantPointer(uint16(8)))
	}
	b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantRange.B51ADiagnosticVariantConditions, b51ADiagnosticVariantFollowing.B51ADiagnosticVariantConditions)
	b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantEnum.B51ADiagnosticVariantConditions, b51ADiagnosticVariantFollowing.B51ADiagnosticVariantConditions)
	b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantFour.B51ADiagnosticVariantConditions, []B51ADiagnosticVariantCondition{{B51ADiagnosticVariantField: "flags", B51ADiagnosticVariantBit: 0, B51ADiagnosticVariantSet: false}, {B51ADiagnosticVariantField: "flags", B51ADiagnosticVariantBit: 3, B51ADiagnosticVariantSet: true}})
	b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantApplies(b51ADiagnosticVariantRange.B51ADiagnosticVariantConditions, 8), true)
	b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantApplies(b51ADiagnosticVariantFour.B51ADiagnosticVariantConditions, 8), true)
	for _, b51ADiagnosticVariantAlternative := range b51ADiagnosticVariantFollowing.B51ADiagnosticVariantAlternatives {
		b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantApplies(b51ADiagnosticVariantAlternative.B51ADiagnosticVariantConditions, 255), false)
	}
	b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantFollowing.B51ADiagnosticVariantUnknown, []string{"alternative selection", "step enforcement", "physical length", "qualification", "invocation"})
}

func TestAb98ReplicationAndVersionDirections(b51ADiagnosticVariantT *b51ADiagnosticVariantTesting.T) {
	b51ADiagnosticVariantStatus := b51ADiagnosticVariantGet(b51ADiagnosticVariantT, 7, "positive", "replication-status")
	b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantStatus.B51ADiagnosticVariantPosition, b51ADiagnosticVariantPointer(uint16(4)))
	b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantStatus.B51ADiagnosticVariantWidthBits, b51ADiagnosticVariantPointer(uint16(8)))
	b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantStatus.B51ADiagnosticVariantDefault, b51ADiagnosticVariantPointer(uint64(0)))
	b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantStatus.B51ADiagnosticVariantValues, map[uint64]string{0: "false", 1: "true"})
	b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantStatus.B51ADiagnosticVariantInterpretations, []string{"Whole byte 0 false / 1 true.", "Nonzero-truth competes with unlabeled 2..255; raw 02 remains unresolved."})
	b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantStatus.B51ADiagnosticVariantUnknown, []string{"truth for 2..255", "replication algorithms", "physical length", "qualification", "invocation"})
	b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantGet(b51ADiagnosticVariantT, 7, "request", "flags").B51ADiagnosticVariantUnknown, []string{"replication algorithms", "physical length", "qualification", "invocation"})
	b51ADiagnosticVariantConfiguration := b51ADiagnosticVariantGet(b51ADiagnosticVariantT, 8, "request", "configuration")
	b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantConfiguration.B51ADiagnosticVariantDefault, b51ADiagnosticVariantPointer(uint64(0)))
	b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantConfiguration.B51ADiagnosticVariantValues, map[uint64]string{1: "initiator-spare", 2: "responder-spare"})
	for _, b51ADiagnosticVariantCase := range []struct {
		b51ADiagnosticVariantDirection, b51ADiagnosticVariantField string
		b51ADiagnosticVariantPosition, b51ADiagnosticVariantWidth  uint16
	}{
		{"request", "responder-written", 4, 16}, {"request", "initiator-to-write", 6, 16}, {"request", "configuration", 8, 8}, {"request", "responder-version", 9, 8}, {"request", "initiator-version", 10, 8},
		{"positive", "initiator-written", 4, 0}, {"positive", "responder-to-write", 6, 0}, {"positive", "responder-version", 8, 8}, {"positive", "initiator-version", 9, 8},
	} {
		b51ADiagnosticVariantEntry := b51ADiagnosticVariantGet(b51ADiagnosticVariantT, 8, b51ADiagnosticVariantCase.b51ADiagnosticVariantDirection, b51ADiagnosticVariantCase.b51ADiagnosticVariantField)
		b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantEntry.B51ADiagnosticVariantPosition, b51ADiagnosticVariantPointer(b51ADiagnosticVariantCase.b51ADiagnosticVariantPosition))
		if b51ADiagnosticVariantCase.b51ADiagnosticVariantWidth == 0 {
			b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantEntry.B51ADiagnosticVariantWidthBits, (*uint16)(nil))
			b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantEntry.B51ADiagnosticVariantByteOrder, (*string)(nil))
			b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantEntry.B51ADiagnosticVariantDomain, []uint64(nil))
		} else {
			b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantEntry.B51ADiagnosticVariantWidthBits, b51ADiagnosticVariantPointer(b51ADiagnosticVariantCase.b51ADiagnosticVariantWidth))
			b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantEntry.B51ADiagnosticVariantDomain, []uint64{0, (1 << b51ADiagnosticVariantCase.b51ADiagnosticVariantWidth) - 1})
		}
		b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantEntry.B51ADiagnosticVariantDefault, b51ADiagnosticVariantPointer(uint64(0)))
		if b51ADiagnosticVariantCase.b51ADiagnosticVariantWidth == 16 {
			b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantEntry.B51ADiagnosticVariantByteOrder, b51ADiagnosticVariantPointer("little-endian"))
		}
	}
	b51ADiagnosticVariantRequest := []byte{0, 0, 0, 0, 0, 0, 0, 0, 2, 41, 53}
	b51ADiagnosticVariantPositive := []byte{0, 0, 0, 0, 0, 0, 0, 0, 67, 79}
	b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantRequest[*b51ADiagnosticVariantGet(b51ADiagnosticVariantT, 8, "request", "responder-version").B51ADiagnosticVariantPosition], byte(41))
	b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantRequest[*b51ADiagnosticVariantGet(b51ADiagnosticVariantT, 8, "request", "initiator-version").B51ADiagnosticVariantPosition], byte(53))
	b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantPositive[*b51ADiagnosticVariantGet(b51ADiagnosticVariantT, 8, "positive", "responder-version").B51ADiagnosticVariantPosition], byte(67))
	b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantPositive[*b51ADiagnosticVariantGet(b51ADiagnosticVariantT, 8, "positive", "initiator-version").B51ADiagnosticVariantPosition], byte(79))
}

func TestAb98IndependentStorage(b51ADiagnosticVariantT *b51ADiagnosticVariantTesting.T) {
	b51ADiagnosticVariantBefore, b51ADiagnosticVariantErr := B51ADiagnosticVariantEntries()
	if b51ADiagnosticVariantErr != nil {
		b51ADiagnosticVariantT.Fatal(b51ADiagnosticVariantErr)
	}
	b51ADiagnosticVariantMutated, b51ADiagnosticVariantErr := B51ADiagnosticVariantEntries()
	if b51ADiagnosticVariantErr != nil {
		b51ADiagnosticVariantT.Fatal(b51ADiagnosticVariantErr)
	}
	for b51ADiagnosticVariantIndex := range b51ADiagnosticVariantMutated {
		b51ADiagnosticVariantEntry := &b51ADiagnosticVariantMutated[b51ADiagnosticVariantIndex]
		*b51ADiagnosticVariantEntry.B51ADiagnosticVariantKey.B51ADiagnosticVariantSecondary = 0
		*b51ADiagnosticVariantEntry.B51ADiagnosticVariantKey.B51ADiagnosticVariantOperation = 0
		if b51ADiagnosticVariantEntry.B51ADiagnosticVariantPosition != nil {
			*b51ADiagnosticVariantEntry.B51ADiagnosticVariantPosition = 0
		}
		if len(b51ADiagnosticVariantEntry.B51ADiagnosticVariantDomain) > 0 {
			b51ADiagnosticVariantEntry.B51ADiagnosticVariantDomain[0] = 99
		}
		if len(b51ADiagnosticVariantEntry.B51ADiagnosticVariantBits) > 0 {
			b51ADiagnosticVariantEntry.B51ADiagnosticVariantBits[0].B51ADiagnosticVariantLabels[0] = "mutated"
		}
		if len(b51ADiagnosticVariantEntry.B51ADiagnosticVariantConditions) > 0 {
			b51ADiagnosticVariantEntry.B51ADiagnosticVariantConditions[0].B51ADiagnosticVariantBit = 7
		}
		if len(b51ADiagnosticVariantEntry.B51ADiagnosticVariantAlternatives) > 0 {
			b51ADiagnosticVariantEntry.B51ADiagnosticVariantAlternatives[0].B51ADiagnosticVariantFields[0].B51ADiagnosticVariantKey.B51ADiagnosticVariantField = "mutated"
			*b51ADiagnosticVariantEntry.B51ADiagnosticVariantAlternatives[0].B51ADiagnosticVariantFields[0].B51ADiagnosticVariantPosition = 99
			b51ADiagnosticVariantEntry.B51ADiagnosticVariantAlternatives[0].B51ADiagnosticVariantInterpretations[0] = "mutated"
		}
		if len(b51ADiagnosticVariantEntry.B51ADiagnosticVariantValues) > 0 {
			b51ADiagnosticVariantEntry.B51ADiagnosticVariantValues[0] = "mutated"
		}
		b51ADiagnosticVariantEntry.B51ADiagnosticVariantUnknown[0] = "mutated"
	}
	b51ADiagnosticVariantAfter, b51ADiagnosticVariantErr := B51ADiagnosticVariantEntries()
	if b51ADiagnosticVariantErr != nil {
		b51ADiagnosticVariantT.Fatal(b51ADiagnosticVariantErr)
	}
	b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantAfter, b51ADiagnosticVariantBefore)
	b51ADiagnosticVariantOne := b51ADiagnosticVariantGet(b51ADiagnosticVariantT, 5, "positive", "following-data")
	b51ADiagnosticVariantTwo := b51ADiagnosticVariantGet(b51ADiagnosticVariantT, 5, "positive", "following-data")
	*b51ADiagnosticVariantOne.B51ADiagnosticVariantAlternatives[0].B51ADiagnosticVariantFields[0].B51ADiagnosticVariantPosition = 99
	b51ADiagnosticVariantOne.B51ADiagnosticVariantAlternatives[0].B51ADiagnosticVariantConditions[0].B51ADiagnosticVariantBit = 7
	b51ADiagnosticVariantThree := b51ADiagnosticVariantGet(b51ADiagnosticVariantT, 5, "positive", "following-data")
	b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantTwo, b51ADiagnosticVariantThree)
	b51ADiagnosticVariantFlags := b51ADiagnosticVariantGet(b51ADiagnosticVariantT, 1, "request", "flags")
	b51ADiagnosticVariantFlags.B51ADiagnosticVariantBits[0].B51ADiagnosticVariantLabels[0] = "mutated"
	b51ADiagnosticVariantEqual(b51ADiagnosticVariantT, b51ADiagnosticVariantGet(b51ADiagnosticVariantT, 1, "request", "flags").B51ADiagnosticVariantBits[0].B51ADiagnosticVariantLabels[0], "power-on/start")
}

func TestAb98MalformedData(b51ADiagnosticVariantT *b51ADiagnosticVariantTesting.T) {
	b51ADiagnosticVariantOriginal := b51ADiagnosticVariantData
	defer func() { b51ADiagnosticVariantData = b51ADiagnosticVariantOriginal }()
	b51ADiagnosticVariantEntry := b51ADiagnosticVariantGet(b51ADiagnosticVariantT, 1, "request", "sequence")
	b51ADiagnosticVariantDuplicate, b51ADiagnosticVariantErr := b51ADiagnosticVariantJSON.Marshal([]B51ADiagnosticVariantEntry{b51ADiagnosticVariantEntry, b51ADiagnosticVariantEntry})
	if b51ADiagnosticVariantErr != nil {
		b51ADiagnosticVariantT.Fatal(b51ADiagnosticVariantErr)
	}
	for _, b51ADiagnosticVariantInput := range []string{"{", "null", "[]", "[{}]", string(b51ADiagnosticVariantOriginal) + " []", string(b51ADiagnosticVariantDuplicate), b51ADiagnosticVariantStrings.Replace(string(b51ADiagnosticVariantOriginal), `"Hypothesis"`, `"qualified"`, 1), b51ADiagnosticVariantStrings.Replace(string(b51ADiagnosticVariantOriginal), `"position": 3`, `"position": -1`, 1), b51ADiagnosticVariantStrings.Replace(string(b51ADiagnosticVariantOriginal), `"grade":`, `"unexpected":`, 1),
		b51ADiagnosticVariantStrings.Replace(string(b51ADiagnosticVariantOriginal), `"bit": 0`, `"bit": 8`, 1),
		b51ADiagnosticVariantStrings.Replace(string(b51ADiagnosticVariantOriginal), `"width_bits": 8`, `"width_bits": 0`, 1),
		b51ADiagnosticVariantStrings.Replace(string(b51ADiagnosticVariantOriginal), `"little-endian"`, `"guessed-endian"`, 1)} {
		b51ADiagnosticVariantData = []byte(b51ADiagnosticVariantInput)
		b51ADiagnosticVariantEntries, b51ADiagnosticVariantErr := B51ADiagnosticVariantEntries()
		if b51ADiagnosticVariantErr == nil || b51ADiagnosticVariantEntries != nil {
			b51ADiagnosticVariantT.Fatalf("malformed catalog accepted")
		}
		b51ADiagnosticVariantGot, b51ADiagnosticVariantFound, b51ADiagnosticVariantErr := B51ADiagnosticVariantLookup(b51ADiagnosticVariantEntry.B51ADiagnosticVariantKey)
		if b51ADiagnosticVariantErr == nil || b51ADiagnosticVariantFound || !b51ADiagnosticVariantReflect.DeepEqual(b51ADiagnosticVariantGot, B51ADiagnosticVariantEntry{}) {
			b51ADiagnosticVariantT.Fatal("malformed lookup succeeded")
		}
	}
}
