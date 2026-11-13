package vaillant

import (
	b503HistoryFlagBinary "encoding/binary"
	b503HistoryFlagReflect "reflect"
	b503HistoryFlagTesting "testing"
)

func b503HistoryFlagPointer[b503HistoryFlagValue any](b503HistoryFlagInput b503HistoryFlagValue) *b503HistoryFlagValue {
	return &b503HistoryFlagInput
}

func b503HistoryFlagCatalog(b503HistoryFlagT *b503HistoryFlagTesting.T) B503HistoryMetadata {
	b503HistoryFlagT.Helper()
	b503HistoryFlagGot, b503HistoryFlagError := B503HistoryCatalog()
	if b503HistoryFlagError != nil {
		b503HistoryFlagT.Fatal(b503HistoryFlagError)
	}
	return b503HistoryFlagGot
}

func b503HistoryFlagExpectedField(b503HistoryFlagName string, b503HistoryFlagOperation *uint8, b503HistoryFlagWidth *int, b503HistoryFlagSignedness, b503HistoryFlagOrder, b503HistoryFlagConversion string) B503HistoryField {
	return B503HistoryField{
		B503HistoryFlagKey:       B503HistoryFlagKey{B503HistoryFlagFamily: "B503", B503HistoryFlagPrimary: 181, B503HistoryFlagSecondary: b503HistoryFlagPointer(uint8(3)), B503HistoryFlagOperation: b503HistoryFlagOperation, B503HistoryFlagField: b503HistoryFlagName},
		B503HistoryFlagWidthBits: b503HistoryFlagWidth, B503HistoryFlagSignedness: b503HistoryFlagSignedness, B503HistoryFlagByteOrder: b503HistoryFlagOrder, B503HistoryFlagConversion: b503HistoryFlagConversion,
		B503HistoryFlagEnumeration: map[int]string{}, B503HistoryFlagReplacement: []int{}, B503HistoryFlagReserved: []int{}, B503HistoryFlagBits: map[int]string{},
		B503HistoryFlagGrade: "Hypothesis", B503HistoryFlagRawFallback: "preserve original field bytes", B503HistoryFlagBranches: map[int]string{}, B503HistoryFlagNoItemTail: []int{},
	}
}

func TestAb86ProductionMetadata(b503HistoryFlagT *b503HistoryFlagTesting.T) {
	b503HistoryFlagGot := b503HistoryFlagCatalog(b503HistoryFlagT)
	b503HistoryFlagExpected := B503HistoryMetadata{B503HistoryFlagRepresentation: "Proven", B503HistoryFlagQualification: "Unknown", B503HistoryFlagUnknown: []string{
		"century", "timezone", "full calendar validity", "copied-span adaptation", "physical ACK", "reset effect", "persistence", "freshness", "execution permission", "unprovided layouts and flag meanings", "broader family completeness",
	}}
	b503HistoryFlagMarker := b503HistoryFlagExpectedField("marker", b503HistoryFlagPointer(uint8(1)), b503HistoryFlagPointer(8), "unsigned", "not_applicable", "branch-before-calendar")
	b503HistoryFlagMarker.B503HistoryFlagEnumeration = map[int]string{0: "no_item", 1: "code_only", 2: "calendar_components"}
	b503HistoryFlagMarker.B503HistoryFlagBranches = map[int]string{0: "require-seven-FF-tail", 1: "retain-code-without-calendar", 2: "validate-each-component-before-decimal"}
	b503HistoryFlagMarker.B503HistoryFlagNoItemTail = []int{255, 255, 255, 255, 255, 255, 255}
	b503HistoryFlagExpected.B503HistoryFlagFields = append(b503HistoryFlagExpected.B503HistoryFlagFields, b503HistoryFlagMarker, b503HistoryFlagExpectedField("error_code", b503HistoryFlagPointer(uint8(1)), nil, "Unknown", "Unknown", "identity"))
	for b503HistoryFlagIndex, b503HistoryFlagName := range []string{"minute", "hour", "day", "month", "year"} {
		b503HistoryFlagField := b503HistoryFlagExpectedField(b503HistoryFlagName, b503HistoryFlagPointer(uint8(1)), b503HistoryFlagPointer(8), "unsigned", "not_applicable", "decimal-after-validity")
		b503HistoryFlagField.B503HistoryFlagBounds = [][]int{{0, 59}, {0, 23}, {1, 31}, {1, 12}, {0, 99}}[b503HistoryFlagIndex]
		b503HistoryFlagField.B503HistoryFlagReplacement = []int{255}
		b503HistoryFlagExpected.B503HistoryFlagFields = append(b503HistoryFlagExpected.B503HistoryFlagFields, b503HistoryFlagField)
	}
	b503HistoryFlagReset := b503HistoryFlagExpectedField("reset_acceptance", nil, b503HistoryFlagPointer(8), "unsigned", "not_applicable", "identity")
	b503HistoryFlagReset.B503HistoryFlagDefault = b503HistoryFlagPointer(0)
	b503HistoryFlagReset.B503HistoryFlagEnumeration = map[int]string{0: "false", 1: "true"}
	b503HistoryFlagExpected.B503HistoryFlagFields = append(b503HistoryFlagExpected.B503HistoryFlagFields, b503HistoryFlagReset)
	b503HistoryFlagFlags := b503HistoryFlagExpectedField("history_flags", b503HistoryFlagPointer(uint8(7)), b503HistoryFlagPointer(8), "unsigned", "not_applicable", "bitset")
	b503HistoryFlagFlags.B503HistoryFlagBits = map[int]string{0: "history_flag_0", 1: "history_flag_1", 2: "history_flag_2", 3: "history_flag_3", 4: "history_flag_4"}
	b503HistoryFlagFlags.B503HistoryFlagReserved = []int{5, 6, 7}
	b503HistoryFlagExpected.B503HistoryFlagFields = append(b503HistoryFlagExpected.B503HistoryFlagFields, b503HistoryFlagFlags)
	for _, b503HistoryFlagName := range []string{"current_error_1", "current_error_2", "current_error_3", "current_error_4", "current_error_5"} {
		b503HistoryFlagExpected.B503HistoryFlagFields = append(b503HistoryFlagExpected.B503HistoryFlagFields, b503HistoryFlagExpectedField(b503HistoryFlagName, b503HistoryFlagPointer(uint8(7)), b503HistoryFlagPointer(16), "unsigned", "little_endian", "identity"))
	}
	b503HistoryFlagFeedback := b503HistoryFlagExpectedField("feedback", b503HistoryFlagPointer(uint8(7)), b503HistoryFlagPointer(8), "unsigned", "not_applicable", "identity")
	b503HistoryFlagFeedback.B503HistoryFlagEnumeration = map[int]string{1: "accepted"}
	b503HistoryFlagFeedback.B503HistoryFlagOtherwise = b503HistoryFlagPointer("rejected")
	b503HistoryFlagExpected.B503HistoryFlagFields = append(b503HistoryFlagExpected.B503HistoryFlagFields, b503HistoryFlagFeedback)
	if !b503HistoryFlagReflect.DeepEqual(b503HistoryFlagGot, b503HistoryFlagExpected) {
		b503HistoryFlagT.Fatalf("ab86 production declarations differ: got %#v want %#v", b503HistoryFlagGot, b503HistoryFlagExpected)
	}
	for _, b503HistoryFlagField := range b503HistoryFlagGot.B503HistoryFlagFields {
		b503HistoryFlagLookupField, b503HistoryFlagFound, b503HistoryFlagError := B503HistoryLookup(b503HistoryFlagField.B503HistoryFlagKey)
		if b503HistoryFlagError != nil || !b503HistoryFlagFound || !b503HistoryFlagReflect.DeepEqual(b503HistoryFlagLookupField, b503HistoryFlagField) {
			b503HistoryFlagT.Fatal("ab86 lookup control", b503HistoryFlagField.B503HistoryFlagKey, b503HistoryFlagError)
		}
	}
}

func b503HistoryFlagNamed(b503HistoryFlagT *b503HistoryFlagTesting.T, b503HistoryFlagName string) B503HistoryField {
	b503HistoryFlagT.Helper()
	for _, b503HistoryFlagField := range b503HistoryFlagCatalog(b503HistoryFlagT).B503HistoryFlagFields {
		if b503HistoryFlagField.B503HistoryFlagKey.B503HistoryFlagField == b503HistoryFlagName {
			return b503HistoryFlagField
		}
	}
	b503HistoryFlagT.Fatal("ab86 missing field", b503HistoryFlagName)
	return B503HistoryField{}
}

func b503HistoryFlagCalendar(b503HistoryFlagField B503HistoryField, b503HistoryFlagRaw int) (int, bool) {
	for _, b503HistoryFlagReplacement := range b503HistoryFlagField.B503HistoryFlagReplacement {
		if b503HistoryFlagRaw == b503HistoryFlagReplacement {
			return 0, false
		}
	}
	if len(b503HistoryFlagField.B503HistoryFlagBounds) != 2 || b503HistoryFlagRaw < b503HistoryFlagField.B503HistoryFlagBounds[0] || b503HistoryFlagRaw > b503HistoryFlagField.B503HistoryFlagBounds[1] {
		return 0, false
	}
	if b503HistoryFlagField.B503HistoryFlagConversion != "decimal-after-validity" {
		return 0, false
	}
	return b503HistoryFlagRaw, true
}

func b503HistoryFlagHistoryCase(b503HistoryFlagT *b503HistoryFlagTesting.T, b503HistoryFlagMarker int, b503HistoryFlagTail []int, b503HistoryFlagCode int, b503HistoryFlagComponents []int) (string, int, []int) {
	b503HistoryFlagT.Helper()
	b503HistoryFlagField := b503HistoryFlagNamed(b503HistoryFlagT, "marker")
	switch b503HistoryFlagField.B503HistoryFlagBranches[b503HistoryFlagMarker] {
	case "require-seven-FF-tail":
		if b503HistoryFlagReflect.DeepEqual(b503HistoryFlagTail, b503HistoryFlagField.B503HistoryFlagNoItemTail) {
			return b503HistoryFlagField.B503HistoryFlagEnumeration[b503HistoryFlagMarker], 0, nil
		}
	case "retain-code-without-calendar":
		if b503HistoryFlagNamed(b503HistoryFlagT, "error_code").B503HistoryFlagConversion == "identity" {
			return b503HistoryFlagField.B503HistoryFlagEnumeration[b503HistoryFlagMarker], b503HistoryFlagCode, nil
		}
	case "validate-each-component-before-decimal":
		if len(b503HistoryFlagComponents) != 5 {
			return "Unknown", 0, nil
		}
		b503HistoryFlagResult := []int{}
		for b503HistoryFlagIndex, b503HistoryFlagName := range []string{"minute", "hour", "day", "month", "year"} {
			b503HistoryFlagValue, b503HistoryFlagValid := b503HistoryFlagCalendar(b503HistoryFlagNamed(b503HistoryFlagT, b503HistoryFlagName), b503HistoryFlagComponents[b503HistoryFlagIndex])
			if !b503HistoryFlagValid {
				return "Unknown", b503HistoryFlagCode, nil
			}
			b503HistoryFlagResult = append(b503HistoryFlagResult, b503HistoryFlagValue)
		}
		return b503HistoryFlagField.B503HistoryFlagEnumeration[b503HistoryFlagMarker], b503HistoryFlagCode, b503HistoryFlagResult
	}
	return "Unknown", 0, nil
}

func TestAb86SyntheticHistory(b503HistoryFlagT *b503HistoryFlagTesting.T) {
	b503HistoryFlagLabel, b503HistoryFlagCode, b503HistoryFlagDate := b503HistoryFlagHistoryCase(b503HistoryFlagT, 0, []int{255, 255, 255, 255, 255, 255, 255}, 0, nil)
	if b503HistoryFlagLabel != "no_item" || b503HistoryFlagCode != 0 || b503HistoryFlagDate != nil {
		b503HistoryFlagT.Fatal("ab86 no item")
	}
	b503HistoryFlagLabel, _, _ = b503HistoryFlagHistoryCase(b503HistoryFlagT, 0, []int{255, 255, 255, 255, 255, 255, 0}, 0, nil)
	if b503HistoryFlagLabel != "Unknown" {
		b503HistoryFlagT.Fatal("ab86 no-item tail must match")
	}
	b503HistoryFlagLabel, b503HistoryFlagCode, b503HistoryFlagDate = b503HistoryFlagHistoryCase(b503HistoryFlagT, 1, nil, 0x1234, []int{255, 255, 255, 255, 255})
	if b503HistoryFlagLabel != "code_only" || b503HistoryFlagCode != 0x1234 || b503HistoryFlagDate != nil {
		b503HistoryFlagT.Fatal("ab86 marker 01 retains code without calendar")
	}
	b503HistoryFlagLabel, b503HistoryFlagCode, b503HistoryFlagDate = b503HistoryFlagHistoryCase(b503HistoryFlagT, 2, nil, 0x1234, []int{59, 23, 31, 12, 99})
	if b503HistoryFlagLabel != "calendar_components" || b503HistoryFlagCode != 0x1234 || !b503HistoryFlagReflect.DeepEqual(b503HistoryFlagDate, []int{59, 23, 31, 12, 99}) {
		b503HistoryFlagT.Fatal("ab86 component boundaries")
	}
	for b503HistoryFlagIndex, b503HistoryFlagName := range []string{"minute", "hour", "day", "month", "year"} {
		b503HistoryFlagField := b503HistoryFlagNamed(b503HistoryFlagT, b503HistoryFlagName)
		for _, b503HistoryFlagRaw := range []int{255, []int{60, 24, 32, 13, 100}[b503HistoryFlagIndex], []int{-1, -1, 0, 0, -1}[b503HistoryFlagIndex]} {
			if _, b503HistoryFlagValid := b503HistoryFlagCalendar(b503HistoryFlagField, b503HistoryFlagRaw); b503HistoryFlagValid {
				b503HistoryFlagT.Fatal("ab86 invalid/replacement converted", b503HistoryFlagName, b503HistoryFlagRaw)
			}
		}
		b503HistoryFlagLow := []int{0, 0, 1, 1, 0}[b503HistoryFlagIndex]
		if b503HistoryFlagValue, b503HistoryFlagValid := b503HistoryFlagCalendar(b503HistoryFlagField, b503HistoryFlagLow); !b503HistoryFlagValid || b503HistoryFlagValue != b503HistoryFlagLow {
			b503HistoryFlagT.Fatal("ab86 lower bound", b503HistoryFlagName)
		}
	}
	b503HistoryFlagLabel, _, b503HistoryFlagDate = b503HistoryFlagHistoryCase(b503HistoryFlagT, 2, nil, 1, []int{255, 255, 255, 255, 255})
	if b503HistoryFlagLabel != "Unknown" || b503HistoryFlagDate != nil {
		b503HistoryFlagT.Fatal("ab86 replacement has no calendar claim")
	}
}

func b503HistoryFlagLabel(b503HistoryFlagField B503HistoryField, b503HistoryFlagRaw int) string {
	if b503HistoryFlagLabel, b503HistoryFlagFound := b503HistoryFlagField.B503HistoryFlagEnumeration[b503HistoryFlagRaw]; b503HistoryFlagFound {
		return b503HistoryFlagLabel
	}
	if b503HistoryFlagField.B503HistoryFlagOtherwise != nil {
		return *b503HistoryFlagField.B503HistoryFlagOtherwise
	}
	return "unlabeled"
}

func TestAb86SyntheticFlagsWordsFeedback(b503HistoryFlagT *b503HistoryFlagTesting.T) {
	b503HistoryFlagFlags := b503HistoryFlagNamed(b503HistoryFlagT, "history_flags")
	for _, b503HistoryFlagRaw := range []int{0x1F, 0xE0} {
		b503HistoryFlagNames := []string{}
		for b503HistoryFlagBit := 0; b503HistoryFlagBit < 8; b503HistoryFlagBit++ {
			if b503HistoryFlagName, b503HistoryFlagFound := b503HistoryFlagFlags.B503HistoryFlagBits[b503HistoryFlagBit]; b503HistoryFlagFound && b503HistoryFlagRaw&(1<<b503HistoryFlagBit) != 0 {
				b503HistoryFlagNames = append(b503HistoryFlagNames, b503HistoryFlagName)
			}
		}
		b503HistoryFlagExpected := []string{}
		if b503HistoryFlagRaw == 0x1F {
			b503HistoryFlagExpected = []string{"history_flag_0", "history_flag_1", "history_flag_2", "history_flag_3", "history_flag_4"}
		}
		if !b503HistoryFlagReflect.DeepEqual(b503HistoryFlagNames, b503HistoryFlagExpected) {
			b503HistoryFlagT.Fatal("ab86 reserved bits included", b503HistoryFlagNames)
		}
	}
	for b503HistoryFlagIndex, b503HistoryFlagName := range []string{"current_error_1", "current_error_2", "current_error_3", "current_error_4", "current_error_5"} {
		b503HistoryFlagField := b503HistoryFlagNamed(b503HistoryFlagT, b503HistoryFlagName)
		b503HistoryFlagBytes := [][]byte{{1, 0}, {0, 1}, {0x34, 0x12}, {0, 0}, {255, 255}}[b503HistoryFlagIndex]
		var b503HistoryFlagValue uint16
		if b503HistoryFlagField.B503HistoryFlagWidthBits != nil && *b503HistoryFlagField.B503HistoryFlagWidthBits == 16 && b503HistoryFlagField.B503HistoryFlagConversion == "identity" && b503HistoryFlagField.B503HistoryFlagSignedness == "unsigned" {
			switch b503HistoryFlagField.B503HistoryFlagByteOrder {
			case "little_endian":
				b503HistoryFlagValue = b503HistoryFlagBinary.LittleEndian.Uint16(b503HistoryFlagBytes)
			case "big_endian":
				b503HistoryFlagValue = b503HistoryFlagBinary.BigEndian.Uint16(b503HistoryFlagBytes)
			default:
				b503HistoryFlagT.Fatal("ab86 unknown byte order")
			}
		} else {
			b503HistoryFlagT.Fatal("ab86 word metadata")
		}
		if b503HistoryFlagValue != []uint16{1, 256, 0x1234, 0, 65535}[b503HistoryFlagIndex] {
			b503HistoryFlagT.Fatal("ab86 independent word", b503HistoryFlagName, b503HistoryFlagValue)
		}
	}
	b503HistoryFlagReset := b503HistoryFlagNamed(b503HistoryFlagT, "reset_acceptance")
	b503HistoryFlagFeedback := b503HistoryFlagNamed(b503HistoryFlagT, "feedback")
	for b503HistoryFlagRaw := 0; b503HistoryFlagRaw <= 255; b503HistoryFlagRaw++ {
		b503HistoryFlagResetWant := "unlabeled"
		b503HistoryFlagFeedbackWant := "rejected"
		if b503HistoryFlagRaw == 0 {
			b503HistoryFlagResetWant = "false"
		}
		if b503HistoryFlagRaw == 1 {
			b503HistoryFlagResetWant = "true"
			b503HistoryFlagFeedbackWant = "accepted"
		}
		if b503HistoryFlagLabel(b503HistoryFlagReset, b503HistoryFlagRaw) != b503HistoryFlagResetWant || b503HistoryFlagLabel(b503HistoryFlagFeedback, b503HistoryFlagRaw) != b503HistoryFlagFeedbackWant {
			b503HistoryFlagT.Fatal("ab86 reset/feedback alternatives", b503HistoryFlagRaw)
		}
	}
}

func b503HistoryFlagCheckAxes(b503HistoryFlagT *b503HistoryFlagTesting.T, b503HistoryFlagField B503HistoryField, b503HistoryFlagProduction bool) {
	b503HistoryFlagT.Helper()
	b503HistoryFlagControl, b503HistoryFlagFound := b503HistoryFlagLookup([]B503HistoryField{b503HistoryFlagField}, b503HistoryFlagField.B503HistoryFlagKey)
	if !b503HistoryFlagFound || !b503HistoryFlagReflect.DeepEqual(b503HistoryFlagControl, b503HistoryFlagField) {
		b503HistoryFlagT.Fatal("ab86 invalid axis control")
	}
	b503HistoryFlagKeyType := b503HistoryFlagReflect.TypeOf(b503HistoryFlagField.B503HistoryFlagKey)
	for b503HistoryFlagIndex := 0; b503HistoryFlagIndex < b503HistoryFlagKeyType.NumField(); b503HistoryFlagIndex++ {
		b503HistoryFlagT.Run(b503HistoryFlagKeyType.Field(b503HistoryFlagIndex).Name, func(b503HistoryFlagT *b503HistoryFlagTesting.T) {
			for b503HistoryFlagCase := 0; b503HistoryFlagCase < 3; b503HistoryFlagCase++ {
				b503HistoryFlagKey := b503HistoryFlagField.B503HistoryFlagKey
				b503HistoryFlagAxis := b503HistoryFlagReflect.ValueOf(&b503HistoryFlagKey).Elem().Field(b503HistoryFlagIndex)
				if b503HistoryFlagAxis.Kind() == b503HistoryFlagReflect.Pointer {
					if b503HistoryFlagCase == 0 {
						b503HistoryFlagAxis.SetZero()
					} else {
						b503HistoryFlagNew := b503HistoryFlagReflect.New(b503HistoryFlagAxis.Type().Elem())
						if b503HistoryFlagCase == 2 {
							if b503HistoryFlagNew.Elem().Kind() == b503HistoryFlagReflect.String {
								b503HistoryFlagNew.Elem().SetString("ab86_other")
							} else {
								b503HistoryFlagNew.Elem().SetUint(99)
							}
						}
						b503HistoryFlagAxis.Set(b503HistoryFlagNew)
					}
				} else if b503HistoryFlagAxis.Kind() == b503HistoryFlagReflect.String {
					if b503HistoryFlagCase == 0 {
						b503HistoryFlagAxis.SetString("")
					} else {
						b503HistoryFlagAxis.SetString("ab86_other")
					}
				} else {
					b503HistoryFlagAxis.SetUint(uint64(b503HistoryFlagCase))
				}
				if b503HistoryFlagReflect.DeepEqual(b503HistoryFlagKey, b503HistoryFlagField.B503HistoryFlagKey) {
					continue
				}
				if _, b503HistoryFlagFound := b503HistoryFlagLookup([]B503HistoryField{b503HistoryFlagField}, b503HistoryFlagKey); b503HistoryFlagFound {
					b503HistoryFlagT.Fatal("ab86 axis ignored", b503HistoryFlagCase)
				}
				if b503HistoryFlagProduction {
					if _, b503HistoryFlagFound, b503HistoryFlagError := B503HistoryLookup(b503HistoryFlagKey); b503HistoryFlagError != nil || b503HistoryFlagFound {
						b503HistoryFlagT.Fatal("ab86 production axis ignored", b503HistoryFlagCase, b503HistoryFlagError)
					}
				}
			}
		})
	}
}

func TestAb86ExactIdentityAxes(b503HistoryFlagT *b503HistoryFlagTesting.T) {
	for _, b503HistoryFlagField := range b503HistoryFlagCatalog(b503HistoryFlagT).B503HistoryFlagFields {
		b503HistoryFlagT.Run("ab86_"+b503HistoryFlagField.B503HistoryFlagKey.B503HistoryFlagField, func(b503HistoryFlagT *b503HistoryFlagTesting.T) {
			b503HistoryFlagCheckAxes(b503HistoryFlagT, b503HistoryFlagField, true)
		})
	}
	b503HistoryFlagQualified := B503HistoryField{B503HistoryFlagKey: B503HistoryFlagKey{
		B503HistoryFlagFamily: "B524", B503HistoryFlagVariant: b503HistoryFlagPointer("ab86_synthetic"), B503HistoryFlagPrimary: 181,
		B503HistoryFlagSecondary: b503HistoryFlagPointer(uint8(36)), B503HistoryFlagOperation: b503HistoryFlagPointer(uint8(1)), B503HistoryFlagDirection: b503HistoryFlagPointer("response"), B503HistoryFlagField: "ab86_synthetic",
		B503HistoryFlagOpcode: b503HistoryFlagPointer(uint8(2)), B503HistoryFlagGroup: b503HistoryFlagPointer(uint8(3)), B503HistoryFlagInstance: b503HistoryFlagPointer(uint8(4)), B503HistoryFlagRegister: b503HistoryFlagPointer(uint16(5)),
	}}
	b503HistoryFlagT.Run("ab86_qualified_B524", func(b503HistoryFlagT *b503HistoryFlagTesting.T) {
		b503HistoryFlagCheckAxes(b503HistoryFlagT, b503HistoryFlagQualified, false)
	})
	b503HistoryFlagZero := b503HistoryFlagQualified
	b503HistoryFlagZero.B503HistoryFlagKey.B503HistoryFlagOpcode = b503HistoryFlagPointer(uint8(0))
	b503HistoryFlagZero.B503HistoryFlagKey.B503HistoryFlagGroup = b503HistoryFlagPointer(uint8(0))
	b503HistoryFlagZero.B503HistoryFlagKey.B503HistoryFlagInstance = b503HistoryFlagPointer(uint8(0))
	b503HistoryFlagZero.B503HistoryFlagKey.B503HistoryFlagRegister = b503HistoryFlagPointer(uint16(0))
	b503HistoryFlagZero.B503HistoryFlagKey.B503HistoryFlagOperation = b503HistoryFlagPointer(uint8(0))
	b503HistoryFlagZero.B503HistoryFlagKey.B503HistoryFlagSecondary = b503HistoryFlagPointer(uint8(0))
	b503HistoryFlagZero.B503HistoryFlagKey.B503HistoryFlagVariant = b503HistoryFlagPointer("")
	b503HistoryFlagZero.B503HistoryFlagKey.B503HistoryFlagDirection = b503HistoryFlagPointer("")
	b503HistoryFlagT.Run("ab86_explicit_zero", func(b503HistoryFlagT *b503HistoryFlagTesting.T) {
		b503HistoryFlagCheckAxes(b503HistoryFlagT, b503HistoryFlagZero, false)
	})
	if _, b503HistoryFlagFound, b503HistoryFlagError := B503HistoryLookup(b503HistoryFlagQualified.B503HistoryFlagKey); b503HistoryFlagFound || b503HistoryFlagError != nil {
		b503HistoryFlagT.Fatal("ab86 synthetic B524 must not become production metadata")
	}
}

func b503HistoryFlagCorrupt(b503HistoryFlagValue b503HistoryFlagReflect.Value) {
	switch b503HistoryFlagValue.Kind() {
	case b503HistoryFlagReflect.Pointer:
		if !b503HistoryFlagValue.IsNil() {
			b503HistoryFlagCorrupt(b503HistoryFlagValue.Elem())
		}
	case b503HistoryFlagReflect.Struct:
		for b503HistoryFlagIndex := 0; b503HistoryFlagIndex < b503HistoryFlagValue.NumField(); b503HistoryFlagIndex++ {
			b503HistoryFlagCorrupt(b503HistoryFlagValue.Field(b503HistoryFlagIndex))
		}
	case b503HistoryFlagReflect.Slice:
		for b503HistoryFlagIndex := 0; b503HistoryFlagIndex < b503HistoryFlagValue.Len(); b503HistoryFlagIndex++ {
			b503HistoryFlagCorrupt(b503HistoryFlagValue.Index(b503HistoryFlagIndex))
		}
	case b503HistoryFlagReflect.Map:
		for _, b503HistoryFlagKey := range b503HistoryFlagValue.MapKeys() {
			b503HistoryFlagValue.SetMapIndex(b503HistoryFlagKey, b503HistoryFlagReflect.Zero(b503HistoryFlagValue.Type().Elem()))
		}
	case b503HistoryFlagReflect.String:
		if b503HistoryFlagValue.CanSet() {
			b503HistoryFlagValue.SetString("ab86_mutated")
		}
	case b503HistoryFlagReflect.Int:
		if b503HistoryFlagValue.CanSet() {
			b503HistoryFlagValue.SetInt(42)
		}
	case b503HistoryFlagReflect.Uint8, b503HistoryFlagReflect.Uint16:
		if b503HistoryFlagValue.CanSet() {
			b503HistoryFlagValue.SetUint(42)
		}
	}
}

func TestAb86NestedCopyIsolation(b503HistoryFlagT *b503HistoryFlagTesting.T) {
	b503HistoryFlagBefore := b503HistoryFlagCatalog(b503HistoryFlagT)
	b503HistoryFlagChanged := b503HistoryFlagCatalog(b503HistoryFlagT)
	b503HistoryFlagCorrupt(b503HistoryFlagReflect.ValueOf(&b503HistoryFlagChanged).Elem())
	if !b503HistoryFlagReflect.DeepEqual(b503HistoryFlagCatalog(b503HistoryFlagT), b503HistoryFlagBefore) {
		b503HistoryFlagT.Fatal("ab86 accessor aliases embedded metadata")
	}
	for _, b503HistoryFlagField := range b503HistoryFlagBefore.B503HistoryFlagFields {
		b503HistoryFlagKey := b503HistoryFlagField.B503HistoryFlagKey
		b503HistoryFlagChanged, b503HistoryFlagFound, b503HistoryFlagError := B503HistoryLookup(b503HistoryFlagKey)
		if !b503HistoryFlagFound || b503HistoryFlagError != nil {
			b503HistoryFlagT.Fatal("ab86 lookup control", b503HistoryFlagError)
		}
		b503HistoryFlagCorrupt(b503HistoryFlagReflect.ValueOf(&b503HistoryFlagChanged).Elem())
		b503HistoryFlagAgain, b503HistoryFlagFound, b503HistoryFlagError := B503HistoryLookup(b503HistoryFlagField.B503HistoryFlagKey)
		if !b503HistoryFlagFound || b503HistoryFlagError != nil || !b503HistoryFlagReflect.DeepEqual(b503HistoryFlagAgain, b503HistoryFlagField) {
			b503HistoryFlagT.Fatal("ab86 lookup aliases nested data", b503HistoryFlagError)
		}
	}
}
