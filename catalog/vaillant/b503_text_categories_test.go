package vaillant

import (
	b503TextCategoryBytesPackage "bytes"
	b503TextCategoryJSON "encoding/json"
	b503TextCategoryReflect "reflect"
	b503TextCategoryTesting "testing"
)

func b503TextCategoryPointer[b503TextCategoryType any](b503TextCategoryValue b503TextCategoryType) *b503TextCategoryType {
	return &b503TextCategoryValue
}

func b503TextCategoryFind(b503TextCategoryTest *b503TextCategoryTesting.T, b503TextCategoryVariant string, b503TextCategoryOperation *uint8, b503TextCategoryField string) B503TextCategoryEntry {
	b503TextCategoryTest.Helper()
	b503TextCategoryEntry, b503TextCategoryFound, b503TextCategoryError := B503TextCategoryLookup(B503TextCategoryKey{B503TextCategoryFamily: "vaillant", B503TextCategoryVariant: b503TextCategoryVariant, B503TextCategoryPrimary: 181, B503TextCategorySecondary: b503TextCategoryPointer(uint8(3)), B503TextCategoryOperation: b503TextCategoryOperation, B503TextCategoryField: b503TextCategoryField})
	if b503TextCategoryError != nil || !b503TextCategoryFound {
		b503TextCategoryTest.Fatalf("lookup %s/%v/%s: %v %v", b503TextCategoryVariant, b503TextCategoryOperation, b503TextCategoryField, b503TextCategoryFound, b503TextCategoryError)
	}
	return b503TextCategoryEntry
}

func TestAb96Facts(b503TextCategoryTest *b503TextCategoryTesting.T) {
	b503TextCategoryEntries, b503TextCategoryError := B503TextCategoryEntries()
	if b503TextCategoryError != nil || len(b503TextCategoryEntries) != 9 {
		b503TextCategoryTest.Fatalf("finite entries: %d %v", len(b503TextCategoryEntries), b503TextCategoryError)
	}
	for _, b503TextCategoryEntry := range b503TextCategoryEntries {
		if b503TextCategoryEntry.B503TextCategoryNativeContract != "" {
			b503TextCategoryTest.Fatal("native contract must remain absent until companion publication")
		}
		if b503TextCategoryEntry.B503TextCategoryGrade != "Hypothesis" || b503TextCategoryEntry.B503TextCategoryMetadataPresence != "Proven" || b503TextCategoryEntry.B503TextCategoryPhysicalQualification != "Unknown" || b503TextCategoryEntry.B503TextCategoryPosition != nil || b503TextCategoryEntry.B503TextCategoryPadding != nil || b503TextCategoryEntry.B503TextCategoryKey.B503TextCategoryDirection != nil || b503TextCategoryEntry.B503TextCategoryRawPreservation != "Retain raw values and partial last-known-good fields independently." {
			b503TextCategoryTest.Fatalf("qualification/unknown/raw boundary: %+v", b503TextCategoryEntry)
		}
		if !b503TextCategoryReflect.DeepEqual(b503TextCategoryEntry.B503TextCategoryUnknown, []string{"absolute coordinates", "termination on the wire", "success", "session exclusivity", "physical effects", "prefix-only positive form does not establish an empty physical response", "direction"}) {
			b503TextCategoryTest.Fatalf("unknown boundaries: %v", b503TextCategoryEntry.B503TextCategoryUnknown)
		}
	}
	for _, b503TextCategoryOperation := range []uint8{3, 6} {
		b503TextCategoryUnlock := b503TextCategoryFind(b503TextCategoryTest, "common", b503TextCategoryPointer(b503TextCategoryOperation), "unlock")
		if b503TextCategoryUnlock.B503TextCategoryWidthCharacters == nil || *b503TextCategoryUnlock.B503TextCategoryWidthCharacters != 8 || b503TextCategoryUnlock.B503TextCategoryEncoding != nil || b503TextCategoryUnlock.B503TextCategoryWidthBits != nil {
			b503TextCategoryTest.Fatal("unlock must declare eight characters without encoding or byte width")
		}
		for _, b503TextCategoryText := range []string{"1234567", "12345678", "123456789"} {
			if (len(b503TextCategoryText) == *b503TextCategoryUnlock.B503TextCategoryWidthCharacters) != (b503TextCategoryText == "12345678") {
				b503TextCategoryTest.Fatal("declared character width")
			}
		}
		b503TextCategoryTrailing := b503TextCategoryFind(b503TextCategoryTest, "common", b503TextCategoryPointer(b503TextCategoryOperation), "unlock_trailing_byte")
		if b503TextCategoryTrailing.B503TextCategoryWidthBits == nil || *b503TextCategoryTrailing.B503TextCategoryWidthBits != 8 || b503TextCategoryTrailing.B503TextCategoryUnsigned == nil || !*b503TextCategoryTrailing.B503TextCategoryUnsigned || b503TextCategoryTrailing.B503TextCategoryDefault == nil || *b503TextCategoryTrailing.B503TextCategoryDefault != 0 || b503TextCategoryTrailing.B503TextCategoryFollows == nil || *b503TextCategoryTrailing.B503TextCategoryFollows != "unlock" {
			b503TextCategoryTest.Fatal("separate trailing unsigned byte/default")
		}
	}
	for _, b503TextCategoryOperation := range []uint8{4, 5} {
		b503TextCategoryText := b503TextCategoryFind(b503TextCategoryTest, "common", b503TextCategoryPointer(b503TextCategoryOperation), "text")
		if b503TextCategoryText.B503TextCategoryEncoding == nil || *b503TextCategoryText.B503TextCategoryEncoding != "ASCII" || b503TextCategoryText.B503TextCategoryExtent == nil || *b503TextCategoryText.B503TextCategoryExtent != "through end of data extent" || b503TextCategoryText.B503TextCategoryWidthBits != nil || b503TextCategoryText.B503TextCategoryWidthCharacters != nil {
			b503TextCategoryTest.Fatal("variable ASCII extent")
		}
		b503TextCategoryBytes := []byte{0x41, 0x42}
		if *b503TextCategoryText.B503TextCategoryEncoding != "ASCII" || string(b503TextCategoryBytes) != "AB" {
			b503TextCategoryTest.Fatal("ASCII AB fixture")
		}
	}
	for _, b503TextCategoryText := range []string{"", "AB", "ABCDE"} {
		b503TextCategoryMode := b503TextCategoryFind(b503TextCategoryTest, "common", b503TextCategoryPointer(uint8(5)), "mode")
		if b503TextCategoryMode.B503TextCategoryPosition != nil || b503TextCategoryMode.B503TextCategoryFollows != nil || b503TextCategoryMode.B503TextCategoryWidthBits == nil || *b503TextCategoryMode.B503TextCategoryWidthBits != 8 || b503TextCategoryMode.B503TextCategoryUnsigned == nil || !*b503TextCategoryMode.B503TextCategoryUnsigned || b503TextCategoryMode.B503TextCategoryDefault != nil || b503TextCategoryMode.B503TextCategoryLabels != nil {
			b503TextCategoryTest.Fatalf("mode byte/unknown coordinate for text length %d", len(b503TextCategoryText))
		}
	}
	b503TextCategoryCommon := b503TextCategoryFind(b503TextCategoryTest, "common", nil, "categories")
	b503TextCategoryReset := b503TextCategoryFind(b503TextCategoryTest, "reset", nil, "categories")
	if !b503TextCategoryReflect.DeepEqual(b503TextCategoryCommon.B503TextCategoryLabels, map[uint8]*string{1: nil, 2: b503TextCategoryPointer("predictive maintenance"), 3: nil, 4: nil, 5: nil}) || !b503TextCategoryReflect.DeepEqual(b503TextCategoryReset.B503TextCategoryLabels, map[uint8]*string{1: b503TextCategoryPointer("error"), 3: b503TextCategoryPointer("status")}) {
		b503TextCategoryTest.Fatal("variant-specific category labels")
	}
	if b503TextCategoryCommon.B503TextCategoryWidthBits != nil || b503TextCategoryReset.B503TextCategoryWidthBits != nil {
		b503TextCategoryTest.Fatal("unprovided category width")
	}
}

func TestAb96ExactKeyAxes(b503TextCategoryTest *b503TextCategoryTesting.T) {
	b503TextCategorySaved := b503TextCategoryData
	b503TextCategoryTest.Cleanup(func() { b503TextCategoryData = b503TextCategorySaved })
	b503TextCategoryKey := B503TextCategoryKey{B503TextCategoryFamily: "vaillant", B503TextCategoryVariant: "synthetic", B503TextCategoryPrimary: 181, B503TextCategorySecondary: b503TextCategoryPointer(uint8(36)), B503TextCategoryOperation: b503TextCategoryPointer(uint8(3)), B503TextCategoryDirection: b503TextCategoryPointer("request"), B503TextCategoryField: "synthetic", B503TextCategoryGroup: b503TextCategoryPointer(uint8(1)), B503TextCategoryRequestedInstance: b503TextCategoryPointer(uint8(2)), B503TextCategoryRegister: b503TextCategoryPointer(uint16(3)), B503TextCategoryCommandRole: b503TextCategoryPointer("read")}
	b503TextCategoryEntry := b503TextCategoryFind(b503TextCategoryTest, "common", b503TextCategoryPointer(uint8(3)), "unlock")
	b503TextCategoryEntry.B503TextCategoryKey = b503TextCategoryKey
	var b503TextCategoryError error
	b503TextCategoryData, b503TextCategoryError = b503TextCategoryJSON.Marshal([]B503TextCategoryEntry{b503TextCategoryEntry})
	if b503TextCategoryError != nil {
		b503TextCategoryTest.Fatal(b503TextCategoryError)
	}
	b503TextCategoryGot, b503TextCategoryFound, b503TextCategoryError := B503TextCategoryLookup(b503TextCategoryKey)
	if b503TextCategoryError != nil || !b503TextCategoryFound || !b503TextCategoryReflect.DeepEqual(b503TextCategoryGot.B503TextCategoryKey, b503TextCategoryKey) {
		b503TextCategoryTest.Fatal("synthetic B524 exact key", b503TextCategoryError)
	}
	for _, b503TextCategoryCase := range []struct {
		b503TextCategoryName   string
		b503TextCategoryChange func(*B503TextCategoryKey)
	}{
		{"family", func(b503TextCategoryKey *B503TextCategoryKey) { b503TextCategoryKey.B503TextCategoryFamily = "other" }},
		{"variant", func(b503TextCategoryKey *B503TextCategoryKey) { b503TextCategoryKey.B503TextCategoryVariant = "other" }},
		{"primary", func(b503TextCategoryKey *B503TextCategoryKey) { b503TextCategoryKey.B503TextCategoryPrimary = 180 }},
		{"secondary", func(b503TextCategoryKey *B503TextCategoryKey) {
			b503TextCategoryKey.B503TextCategorySecondary = b503TextCategoryPointer(uint8(35))
		}},
		{"operation", func(b503TextCategoryKey *B503TextCategoryKey) {
			b503TextCategoryKey.B503TextCategoryOperation = b503TextCategoryPointer(uint8(4))
		}},
		{"direction", func(b503TextCategoryKey *B503TextCategoryKey) {
			b503TextCategoryKey.B503TextCategoryDirection = b503TextCategoryPointer("response")
		}},
		{"field", func(b503TextCategoryKey *B503TextCategoryKey) { b503TextCategoryKey.B503TextCategoryField = "other" }},
		{"group", func(b503TextCategoryKey *B503TextCategoryKey) {
			b503TextCategoryKey.B503TextCategoryGroup = b503TextCategoryPointer(uint8(2))
		}},
		{"requested_instance", func(b503TextCategoryKey *B503TextCategoryKey) {
			b503TextCategoryKey.B503TextCategoryRequestedInstance = b503TextCategoryPointer(uint8(3))
		}},
		{"register", func(b503TextCategoryKey *B503TextCategoryKey) {
			b503TextCategoryKey.B503TextCategoryRegister = b503TextCategoryPointer(uint16(4))
		}},
		{"command_role", func(b503TextCategoryKey *B503TextCategoryKey) {
			b503TextCategoryKey.B503TextCategoryCommandRole = b503TextCategoryPointer("write")
		}},
	} {
		b503TextCategoryTest.Run(b503TextCategoryCase.b503TextCategoryName, func(b503TextCategoryTest *b503TextCategoryTesting.T) {
			b503TextCategoryChanged := b503TextCategoryKey
			b503TextCategoryCase.b503TextCategoryChange(&b503TextCategoryChanged)
			b503TextCategoryGot, b503TextCategoryFound, b503TextCategoryError := B503TextCategoryLookup(b503TextCategoryChanged)
			if b503TextCategoryError != nil || b503TextCategoryFound || !b503TextCategoryReflect.DeepEqual(b503TextCategoryGot, B503TextCategoryEntry{}) {
				b503TextCategoryTest.Fatal("changed axis matched", b503TextCategoryError)
			}
		})
	}
}

func TestAb96AbsentZero(b503TextCategoryTest *b503TextCategoryTesting.T) {
	b503TextCategorySaved := b503TextCategoryData
	b503TextCategoryTest.Cleanup(func() { b503TextCategoryData = b503TextCategorySaved })
	b503TextCategoryEntry := b503TextCategoryFind(b503TextCategoryTest, "common", nil, "categories")
	b503TextCategoryAxes := []string{"B503TextCategorySecondary", "B503TextCategoryOperation", "B503TextCategoryGroup", "B503TextCategoryRequestedInstance", "B503TextCategoryRegister", "B503TextCategoryDirection", "B503TextCategoryCommandRole"}
	for _, b503TextCategoryAxis := range b503TextCategoryAxes {
		b503TextCategoryTest.Run(b503TextCategoryAxis, func(b503TextCategoryTest *b503TextCategoryTesting.T) {
			b503TextCategoryAbsent := b503TextCategoryEntry
			b503TextCategoryReflect.ValueOf(&b503TextCategoryAbsent.B503TextCategoryKey).Elem().FieldByName(b503TextCategoryAxis).SetZero()
			b503TextCategoryZero := b503TextCategoryAbsent
			b503TextCategoryField := b503TextCategoryReflect.ValueOf(&b503TextCategoryZero.B503TextCategoryKey).Elem().FieldByName(b503TextCategoryAxis)
			b503TextCategoryField.Set(b503TextCategoryReflect.New(b503TextCategoryField.Type().Elem()))
			for _, b503TextCategoryStored := range []B503TextCategoryEntry{b503TextCategoryAbsent, b503TextCategoryZero} {
				var b503TextCategoryError error
				b503TextCategoryData, b503TextCategoryError = b503TextCategoryJSON.Marshal([]B503TextCategoryEntry{b503TextCategoryStored})
				if b503TextCategoryError != nil {
					b503TextCategoryTest.Fatal(b503TextCategoryError)
				}
				for _, b503TextCategoryQuery := range []B503TextCategoryKey{b503TextCategoryAbsent.B503TextCategoryKey, b503TextCategoryZero.B503TextCategoryKey} {
					_, b503TextCategoryFound, b503TextCategoryError := B503TextCategoryLookup(b503TextCategoryQuery)
					if b503TextCategoryError != nil || b503TextCategoryFound != b503TextCategoryReflect.DeepEqual(b503TextCategoryQuery, b503TextCategoryStored.B503TextCategoryKey) {
						b503TextCategoryTest.Fatal("absent/zero collapsed", b503TextCategoryError)
					}
				}
			}
		})
	}
}

func b503TextCategoryMutate(b503TextCategoryValue b503TextCategoryReflect.Value) {
	switch b503TextCategoryValue.Kind() {
	case b503TextCategoryReflect.Pointer:
		if !b503TextCategoryValue.IsNil() {
			b503TextCategoryMutate(b503TextCategoryValue.Elem())
		}
	case b503TextCategoryReflect.Struct:
		for b503TextCategoryIndex := 0; b503TextCategoryIndex < b503TextCategoryValue.NumField(); b503TextCategoryIndex++ {
			b503TextCategoryMutate(b503TextCategoryValue.Field(b503TextCategoryIndex))
		}
	case b503TextCategoryReflect.Slice:
		for b503TextCategoryIndex := 0; b503TextCategoryIndex < b503TextCategoryValue.Len(); b503TextCategoryIndex++ {
			b503TextCategoryMutate(b503TextCategoryValue.Index(b503TextCategoryIndex))
		}
	case b503TextCategoryReflect.Map:
		for _, b503TextCategoryKey := range b503TextCategoryValue.MapKeys() {
			b503TextCategoryMutate(b503TextCategoryValue.MapIndex(b503TextCategoryKey))
		}
		b503TextCategoryValue.Clear()
	case b503TextCategoryReflect.String:
		if b503TextCategoryValue.CanSet() {
			b503TextCategoryValue.SetString("changed")
		}
	case b503TextCategoryReflect.Uint8, b503TextCategoryReflect.Uint16:
		if b503TextCategoryValue.CanSet() {
			b503TextCategoryValue.SetUint(0)
		}
	case b503TextCategoryReflect.Int:
		if b503TextCategoryValue.CanSet() {
			b503TextCategoryValue.SetInt(99)
		}
	case b503TextCategoryReflect.Bool:
		if b503TextCategoryValue.CanSet() {
			b503TextCategoryValue.SetBool(false)
		}
	}
}

func TestAb96IndependentStorage(b503TextCategoryTest *b503TextCategoryTesting.T) {
	b503TextCategoryExpected, b503TextCategoryError := B503TextCategoryEntries()
	if b503TextCategoryError != nil {
		b503TextCategoryTest.Fatal(b503TextCategoryError)
	}
	b503TextCategoryEntries, b503TextCategoryError := B503TextCategoryEntries()
	if b503TextCategoryError != nil {
		b503TextCategoryTest.Fatal(b503TextCategoryError)
	}
	for b503TextCategoryIndex := range b503TextCategoryEntries {
		b503TextCategoryLookup, b503TextCategoryFound, b503TextCategoryError := B503TextCategoryLookup(b503TextCategoryExpected[b503TextCategoryIndex].B503TextCategoryKey)
		if b503TextCategoryError != nil || !b503TextCategoryFound {
			b503TextCategoryTest.Fatal(b503TextCategoryError)
		}
		b503TextCategoryMutate(b503TextCategoryReflect.ValueOf(&b503TextCategoryLookup).Elem())
		b503TextCategoryAgain, b503TextCategoryFound, b503TextCategoryError := B503TextCategoryLookup(b503TextCategoryExpected[b503TextCategoryIndex].B503TextCategoryKey)
		if b503TextCategoryError != nil || !b503TextCategoryFound || !b503TextCategoryReflect.DeepEqual(b503TextCategoryAgain, b503TextCategoryExpected[b503TextCategoryIndex]) {
			b503TextCategoryTest.Fatal("lookup aliases nested storage")
		}
	}
	b503TextCategoryMutate(b503TextCategoryReflect.ValueOf(&b503TextCategoryEntries).Elem())
	b503TextCategoryAgain, b503TextCategoryError := B503TextCategoryEntries()
	if b503TextCategoryError != nil || !b503TextCategoryReflect.DeepEqual(b503TextCategoryAgain, b503TextCategoryExpected) {
		b503TextCategoryTest.Fatal("entries alias nested storage")
	}
}

func TestAb96MalformedData(b503TextCategoryTest *b503TextCategoryTesting.T) {
	b503TextCategorySaved := b503TextCategoryData
	b503TextCategoryTest.Cleanup(func() { b503TextCategoryData = b503TextCategorySaved })
	b503TextCategoryKey := b503TextCategoryFind(b503TextCategoryTest, "common", nil, "categories").B503TextCategoryKey
	for _, b503TextCategoryInvalid := range []string{"{", "null", "[]", "[{}]", string(b503TextCategorySaved) + " []", "[{\"unexpected\":1}]"} {
		b503TextCategoryData = []byte(b503TextCategoryInvalid)
		b503TextCategoryEntries, b503TextCategoryError := B503TextCategoryEntries()
		if b503TextCategoryError == nil || b503TextCategoryEntries != nil {
			b503TextCategoryTest.Fatal("invalid entries succeeded")
		}
		b503TextCategoryEntry, b503TextCategoryFound, b503TextCategoryError := B503TextCategoryLookup(b503TextCategoryKey)
		if b503TextCategoryError == nil || b503TextCategoryFound || !b503TextCategoryReflect.DeepEqual(b503TextCategoryEntry, B503TextCategoryEntry{}) {
			b503TextCategoryTest.Fatal("invalid lookup succeeded")
		}
	}
	b503TextCategoryData = b503TextCategorySaved
	b503TextCategoryEntries, b503TextCategoryError := B503TextCategoryEntries()
	if b503TextCategoryError != nil {
		b503TextCategoryTest.Fatal(b503TextCategoryError)
	}
	b503TextCategoryData, b503TextCategoryError = b503TextCategoryJSON.Marshal([]B503TextCategoryEntry{b503TextCategoryEntries[0], b503TextCategoryEntries[0]})
	if b503TextCategoryError != nil {
		b503TextCategoryTest.Fatal(b503TextCategoryError)
	}
	if _, b503TextCategoryError = B503TextCategoryEntries(); b503TextCategoryError == nil {
		b503TextCategoryTest.Fatal("duplicate accepted")
	}
}

func TestAb96PrimaryPresence(b503TextCategoryTest *b503TextCategoryTesting.T) {
	b503TextCategorySaved := b503TextCategoryData
	b503TextCategoryTest.Cleanup(func() { b503TextCategoryData = b503TextCategorySaved })
	b503TextCategoryExpected := b503TextCategoryFind(b503TextCategoryTest, "common", nil, "categories")
	b503TextCategoryExpected.B503TextCategoryKey.B503TextCategoryPrimary = 0
	b503TextCategoryEncoded, b503TextCategoryError := b503TextCategoryJSON.Marshal([]B503TextCategoryEntry{b503TextCategoryExpected})
	if b503TextCategoryError != nil {
		b503TextCategoryTest.Fatal(b503TextCategoryError)
	}
	for _, b503TextCategoryCase := range []struct {
		b503TextCategoryName    string
		b503TextCategoryPrimary string
		b503TextCategoryValid   bool
	}{
		{"omitted", "", false},
		{"null", "\"primary\":null,", false},
		{"zero", "\"primary\":0,", true},
	} {
		b503TextCategoryTest.Run(b503TextCategoryCase.b503TextCategoryName, func(b503TextCategoryTest *b503TextCategoryTesting.T) {
			b503TextCategoryData = b503TextCategoryBytesPackage.Replace(b503TextCategoryEncoded, []byte("\"primary\":0,"), []byte(b503TextCategoryCase.b503TextCategoryPrimary), 1)
			b503TextCategoryEntries, b503TextCategoryError := B503TextCategoryEntries()
			if b503TextCategoryCase.b503TextCategoryValid {
				if b503TextCategoryError != nil || !b503TextCategoryReflect.DeepEqual(b503TextCategoryEntries, []B503TextCategoryEntry{b503TextCategoryExpected}) {
					b503TextCategoryTest.Fatalf("explicit primary zero: %v %v", b503TextCategoryEntries, b503TextCategoryError)
				}
			} else if b503TextCategoryError == nil || b503TextCategoryEntries != nil {
				b503TextCategoryTest.Fatal("missing primary accepted by entries")
			}
			b503TextCategoryEntry, b503TextCategoryFound, b503TextCategoryError := B503TextCategoryLookup(b503TextCategoryExpected.B503TextCategoryKey)
			if b503TextCategoryCase.b503TextCategoryValid {
				if b503TextCategoryError != nil || !b503TextCategoryFound || !b503TextCategoryReflect.DeepEqual(b503TextCategoryEntry, b503TextCategoryExpected) {
					b503TextCategoryTest.Fatal("explicit primary zero lookup failed", b503TextCategoryError)
				}
			} else if b503TextCategoryError == nil || b503TextCategoryFound || !b503TextCategoryReflect.DeepEqual(b503TextCategoryEntry, B503TextCategoryEntry{}) {
				b503TextCategoryTest.Fatal("missing primary accepted by lookup")
			}
		})
	}
}
