package vaillant

import (
	b506RFSequenceBytes "bytes"
	b506RFSequenceJSON "encoding/json"
	b506RFSequenceFmt "fmt"
	b506RFSequenceReflect "reflect"
	b506RFSequenceTesting "testing"
)

func b506RFSequenceDecode(b506RFSequenceT *b506RFSequenceTesting.T, b506RFSequenceRaw []byte) map[string]b506RFSequenceJSON.RawMessage {
	b506RFSequenceT.Helper()
	var b506RFSequenceObject map[string]b506RFSequenceJSON.RawMessage
	if b506RFSequenceError := b506RFSequenceJSON.Unmarshal(b506RFSequenceRaw, &b506RFSequenceObject); b506RFSequenceError != nil {
		b506RFSequenceT.Fatal(b506RFSequenceError)
	}
	return b506RFSequenceObject
}

func b506RFSequenceRows(b506RFSequenceT *b506RFSequenceTesting.T) []b506RFSequenceJSON.RawMessage {
	b506RFSequenceT.Helper()
	b506RFSequenceData, b506RFSequenceError := B506RFSequenceEntries()
	if b506RFSequenceError != nil {
		b506RFSequenceT.Fatal(b506RFSequenceError)
	}
	var b506RFSequenceResult []b506RFSequenceJSON.RawMessage
	if b506RFSequenceError := b506RFSequenceJSON.Unmarshal(b506RFSequenceData, &b506RFSequenceResult); b506RFSequenceError != nil {
		b506RFSequenceT.Fatal(b506RFSequenceError)
	}
	return b506RFSequenceResult
}

func b506RFSequenceKeyOf(b506RFSequenceT *b506RFSequenceTesting.T, b506RFSequenceRaw []byte) B506RFSequenceKey {
	b506RFSequenceT.Helper()
	var b506RFSequenceEntry b506RFSequenceRow
	if b506RFSequenceError := b506RFSequenceJSON.Unmarshal(b506RFSequenceRaw, &b506RFSequenceEntry); b506RFSequenceError != nil {
		b506RFSequenceT.Fatal(b506RFSequenceError)
	}
	return b506RFSequenceEntry.B506RFSequenceKey
}

func b506RFSequenceWant(b506RFSequenceT *b506RFSequenceTesting.T, b506RFSequenceObject map[string]b506RFSequenceJSON.RawMessage, b506RFSequenceField, b506RFSequenceExpected string) {
	b506RFSequenceT.Helper()
	if string(b506RFSequenceObject[b506RFSequenceField]) != b506RFSequenceExpected {
		b506RFSequenceT.Fatalf("%s = %s; want %s", b506RFSequenceField, b506RFSequenceObject[b506RFSequenceField], b506RFSequenceExpected)
	}
}

func TestAb17FiniteMetadata(b506RFSequenceT *b506RFSequenceTesting.T) {
	b506RFSequenceEntries := b506RFSequenceRows(b506RFSequenceT)
	if len(b506RFSequenceEntries) != 6 {
		b506RFSequenceT.Fatalf("rows: %d", len(b506RFSequenceEntries))
	}
	for _, b506RFSequenceRaw := range b506RFSequenceEntries {
		b506RFSequenceKey := b506RFSequenceKeyOf(b506RFSequenceT, b506RFSequenceRaw)
		b506RFSequenceLookup, b506RFSequenceFound, b506RFSequenceError := B506RFSequenceLookup(b506RFSequenceKey)
		if b506RFSequenceError != nil || !b506RFSequenceFound || !b506RFSequenceBytes.Equal(b506RFSequenceLookup, b506RFSequenceRaw) {
			b506RFSequenceT.Fatalf("whole row lookup: %v %v", b506RFSequenceFound, b506RFSequenceError)
		}
		b506RFSequenceObject := b506RFSequenceDecode(b506RFSequenceT, b506RFSequenceLookup)
		b506RFSequenceWant(b506RFSequenceT, b506RFSequenceObject, "grade", `"Hypothesis"`)
		b506RFSequenceWant(b506RFSequenceT, b506RFSequenceObject, "qualification", `"Unknown"`)
		b506RFSequenceWant(b506RFSequenceT, b506RFSequenceObject, "ab17_raw_form", `"retain-original-bytes"`)
		for _, b506RFSequenceField := range []string{"ab17_extent_bytes", "ab17_minimum_bytes", "ab17_maximum_bytes", "ab17_field_position"} {
			b506RFSequenceWant(b506RFSequenceT, b506RFSequenceObject, b506RFSequenceField, "null")
		}
		b506RFSequenceUnknown := b506RFSequenceDecode(b506RFSequenceT, b506RFSequenceObject["ab17_unknown"])
		for _, b506RFSequenceField := range []string{"native_support", "physical_effects", "persistence", "delivery"} {
			b506RFSequenceWant(b506RFSequenceT, b506RFSequenceUnknown, b506RFSequenceField, `"Unknown"`)
		}
		if b506RFSequenceKey.B506RFSequenceFamily != "vaillant" || b506RFSequenceKey.B506RFSequenceVariant != "B506-RF" || b506RFSequenceKey.B506RFSequencePrimary != 181 || b506RFSequenceKey.B506RFSequenceSecondary == nil || *b506RFSequenceKey.B506RFSequenceSecondary != 6 || b506RFSequenceKey.B506RFSequenceDirection != "none" || b506RFSequenceKey.B506RFSequenceGroup != nil || b506RFSequenceKey.B506RFSequenceRegister != nil || b506RFSequenceKey.B506RFSequenceInstanceScope != "none" {
			b506RFSequenceT.Fatal("unexpected finite identity")
		}
		if b506RFSequenceKey.B506RFSequenceOperation == nil {
			b506RFSequenceT.Fatal("missing operation")
		}
		if b506RFSequenceKey.B506RFSequenceField != "stop" && *b506RFSequenceKey.B506RFSequenceOperation != 13 {
			b506RFSequenceT.Fatal("RF metadata operation")
		}
	}
	b506RFSequenceCount := b506RFSequenceDecode(b506RFSequenceT, b506RFSequenceDecode(b506RFSequenceT, b506RFSequenceEntries[0])["ab17_count"])
	b506RFSequenceWant(b506RFSequenceT, b506RFSequenceCount, "ab17_zero_meaning", `"unlimited"`)
	b506RFSequenceWant(b506RFSequenceT, b506RFSequenceCount, "ab17_finite_minimum", "1")
	b506RFSequenceWant(b506RFSequenceT, b506RFSequenceCount, "ab17_finite_maximum", "255")
	b506RFSequenceWant(b506RFSequenceT, b506RFSequenceCount, "ab17_finite_meaning", `"count-value"`)
	for _, b506RFSequenceCase := range []struct {
		b506RFSequenceCount    int
		b506RFSequenceExpected string
	}{{0, "unlimited"}, {1, "finite-1"}, {2, "finite-2"}, {255, "finite-255"}} {
		b506RFSequenceT.Run(b506RFSequenceFmt.Sprintf("ab17-count-%d", b506RFSequenceCase.b506RFSequenceCount), func(b506RFSequenceT *b506RFSequenceTesting.T) {
			var b506RFSequenceMinimum, b506RFSequenceMaximum int
			var b506RFSequenceZero, b506RFSequenceMeaning string
			for _, b506RFSequenceDecodeCase := range []struct {
				b506RFSequenceName   string
				b506RFSequenceTarget any
			}{{"ab17_finite_minimum", &b506RFSequenceMinimum}, {"ab17_finite_maximum", &b506RFSequenceMaximum}, {"ab17_zero_meaning", &b506RFSequenceZero}, {"ab17_finite_meaning", &b506RFSequenceMeaning}} {
				if b506RFSequenceError := b506RFSequenceJSON.Unmarshal(b506RFSequenceCount[b506RFSequenceDecodeCase.b506RFSequenceName], b506RFSequenceDecodeCase.b506RFSequenceTarget); b506RFSequenceError != nil {
					b506RFSequenceT.Fatal(b506RFSequenceError)
				}
			}
			b506RFSequenceActual := b506RFSequenceZero
			if b506RFSequenceCase.b506RFSequenceCount != 0 {
				if b506RFSequenceCase.b506RFSequenceCount < b506RFSequenceMinimum || b506RFSequenceCase.b506RFSequenceCount > b506RFSequenceMaximum || b506RFSequenceMeaning != "count-value" {
					b506RFSequenceT.Fatal("finite count unmapped")
				}
				b506RFSequenceActual = b506RFSequenceFmt.Sprintf("finite-%d", b506RFSequenceCase.b506RFSequenceCount)
			}
			if b506RFSequenceActual != b506RFSequenceCase.b506RFSequenceExpected {
				b506RFSequenceT.Fatalf("count meaning %q", b506RFSequenceActual)
			}
		})
	}
	for b506RFSequenceSequence, b506RFSequenceExpected := range []string{`"ab17-class-a"`, `"ab17-class-a"`, `"ab17-class-a"`, `"ab17-class-b"`} {
		b506RFSequenceKey := b506RFSequenceKeyOf(b506RFSequenceT, b506RFSequenceEntries[1])
		b506RFSequenceKey.B506RFSequenceAlternative = b506RFSequenceFmt.Sprint(b506RFSequenceSequence)
		b506RFSequenceRaw, b506RFSequenceFound, b506RFSequenceError := B506RFSequenceLookup(b506RFSequenceKey)
		if b506RFSequenceError != nil || !b506RFSequenceFound {
			b506RFSequenceT.Fatalf("sequence %d: %v", b506RFSequenceSequence, b506RFSequenceError)
		}
		b506RFSequenceObject := b506RFSequenceDecode(b506RFSequenceT, b506RFSequenceRaw)
		b506RFSequenceWant(b506RFSequenceT, b506RFSequenceObject, "ab17_extent_class", b506RFSequenceExpected)
		for _, b506RFSequencePayload := range [][]byte{{}, {0x23, 0x45}} {
			b506RFSequenceT.Run(b506RFSequenceFmt.Sprintf("ab17-sequence-%d-payload-length-%d", b506RFSequenceSequence, len(b506RFSequencePayload)), func(b506RFSequenceT *b506RFSequenceTesting.T) {
				b506RFSequenceWant(b506RFSequenceT, b506RFSequenceObject, "qualification", `"Unknown"`)
				b506RFSequenceWant(b506RFSequenceT, b506RFSequenceObject, "ab17_extent_bytes", "null")
			})
		}
	}
	for b506RFSequenceSequence := 4; b506RFSequenceSequence <= 255; b506RFSequenceSequence++ {
		b506RFSequenceKey := b506RFSequenceKeyOf(b506RFSequenceT, b506RFSequenceEntries[1])
		b506RFSequenceKey.B506RFSequenceAlternative = b506RFSequenceFmt.Sprint(b506RFSequenceSequence)
		if b506RFSequenceRaw, b506RFSequenceFound, b506RFSequenceError := B506RFSequenceLookup(b506RFSequenceKey); b506RFSequenceError != nil || b506RFSequenceFound || b506RFSequenceRaw != nil {
			b506RFSequenceT.Fatalf("unmapped sequence %d", b506RFSequenceSequence)
		}
	}
	b506RFSequenceStop := b506RFSequenceDecode(b506RFSequenceT, b506RFSequenceEntries[5])
	b506RFSequenceStopKey := b506RFSequenceKeyOf(b506RFSequenceT, b506RFSequenceEntries[5])
	if b506RFSequenceStopKey.B506RFSequenceField != "stop" || *b506RFSequenceStopKey.B506RFSequenceOperation != 14 {
		b506RFSequenceT.Fatal("stop identity")
	}
	var b506RFSequenceRelated []int
	if b506RFSequenceError := b506RFSequenceJSON.Unmarshal(b506RFSequenceStop["ab17_related_operations"], &b506RFSequenceRelated); b506RFSequenceError != nil {
		b506RFSequenceT.Fatal(b506RFSequenceError)
	}
	if !b506RFSequenceReflect.DeepEqual(b506RFSequenceRelated, []int{12, 13}) {
		b506RFSequenceT.Fatal("stop associations")
	}
	for _, b506RFSequenceField := range []string{"ab17_activation_order", "ab17_delivery_guarantee", "ab17_completed_stop"} {
		b506RFSequenceWant(b506RFSequenceT, b506RFSequenceStop, b506RFSequenceField, `"Unknown"`)
	}
	for _, b506RFSequenceField := range []string{"ab17_retry_count", "ab17_timeout"} {
		b506RFSequenceWant(b506RFSequenceT, b506RFSequenceStop, b506RFSequenceField, "null")
	}
}

func TestAb17IndependentAxes(b506RFSequenceT *b506RFSequenceTesting.T) {
	b506RFSequenceBase := b506RFSequenceRows(b506RFSequenceT)[0]
	b506RFSequenceBaseKey := b506RFSequenceKeyOf(b506RFSequenceT, b506RFSequenceBase)
	b506RFSequenceZero := uint8(0)
	b506RFSequenceOther := uint8(7)
	b506RFSequenceRegister := uint16(0)
	for _, b506RFSequenceCase := range []struct {
		b506RFSequenceName   string
		b506RFSequenceChange func(*B506RFSequenceKey)
	}{
		{"ab17-family", func(b506RFSequenceKey *B506RFSequenceKey) { b506RFSequenceKey.B506RFSequenceFamily = "ab17-other" }},
		{"ab17-variant", func(b506RFSequenceKey *B506RFSequenceKey) { b506RFSequenceKey.B506RFSequenceVariant = "ab17-other" }},
		{"ab17-primary", func(b506RFSequenceKey *B506RFSequenceKey) { b506RFSequenceKey.B506RFSequencePrimary = 182 }},
		{"ab17-secondary", func(b506RFSequenceKey *B506RFSequenceKey) {
			b506RFSequenceKey.B506RFSequenceSecondary = &b506RFSequenceOther
		}},
		{"ab17-operation", func(b506RFSequenceKey *B506RFSequenceKey) {
			b506RFSequenceKey.B506RFSequenceOperation = &b506RFSequenceOther
		}},
		{"ab17-direction", func(b506RFSequenceKey *B506RFSequenceKey) { b506RFSequenceKey.B506RFSequenceDirection = "ab17-other" }},
		{"ab17-field", func(b506RFSequenceKey *B506RFSequenceKey) { b506RFSequenceKey.B506RFSequenceField = "ab17-other" }},
		{"ab17-alternative", func(b506RFSequenceKey *B506RFSequenceKey) { b506RFSequenceKey.B506RFSequenceAlternative = "ab17-other" }},
		{"ab17-group", func(b506RFSequenceKey *B506RFSequenceKey) {
			b506RFSequenceKey.B506RFSequenceGroup = &b506RFSequenceOther
		}},
		{"ab17-instance-scope", func(b506RFSequenceKey *B506RFSequenceKey) { b506RFSequenceKey.B506RFSequenceInstanceScope = "dynamic" }},
		{"ab17-register", func(b506RFSequenceKey *B506RFSequenceKey) {
			b506RFSequenceOtherRegister := uint16(1)
			b506RFSequenceKey.B506RFSequenceRegister = &b506RFSequenceOtherRegister
		}},
		{"ab17-secondary-absent-zero", func(b506RFSequenceKey *B506RFSequenceKey) { b506RFSequenceKey.B506RFSequenceSecondary = nil }},
		{"ab17-operation-absent-zero", func(b506RFSequenceKey *B506RFSequenceKey) { b506RFSequenceKey.B506RFSequenceOperation = nil }},
		{"ab17-group-absent-zero", func(b506RFSequenceKey *B506RFSequenceKey) {
			b506RFSequenceKey.B506RFSequenceGroup = &b506RFSequenceZero
		}},
		{"ab17-register-absent-zero", func(b506RFSequenceKey *B506RFSequenceKey) {
			b506RFSequenceKey.B506RFSequenceRegister = &b506RFSequenceRegister
		}},
	} {
		b506RFSequenceT.Run(b506RFSequenceCase.b506RFSequenceName, func(b506RFSequenceT *b506RFSequenceTesting.T) {
			b506RFSequenceControl := b506RFSequenceBaseKey
			if b506RFSequenceCase.b506RFSequenceName == "ab17-secondary-absent-zero" {
				b506RFSequenceControl.B506RFSequenceSecondary = &b506RFSequenceZero
			}
			if b506RFSequenceCase.b506RFSequenceName == "ab17-operation-absent-zero" {
				b506RFSequenceControl.B506RFSequenceOperation = &b506RFSequenceZero
			}
			b506RFSequenceChanged := b506RFSequenceControl
			b506RFSequenceCase.b506RFSequenceChange(&b506RFSequenceChanged)
			if b506RFSequenceRaw, b506RFSequenceFound, b506RFSequenceError := B506RFSequenceLookup(b506RFSequenceChanged); b506RFSequenceError != nil || b506RFSequenceFound || b506RFSequenceRaw != nil {
				b506RFSequenceT.Fatal("public lookup fell back")
			}
			b506RFSequenceObjects := []map[string]b506RFSequenceJSON.RawMessage{b506RFSequenceDecode(b506RFSequenceT, b506RFSequenceBase), b506RFSequenceDecode(b506RFSequenceT, b506RFSequenceBase)}
			for b506RFSequenceIndex, b506RFSequenceKey := range []B506RFSequenceKey{b506RFSequenceControl, b506RFSequenceChanged} {
				b506RFSequenceKeyJSON, b506RFSequenceError := b506RFSequenceJSON.Marshal(b506RFSequenceKey)
				if b506RFSequenceError != nil {
					b506RFSequenceT.Fatal(b506RFSequenceError)
				}
				b506RFSequenceObjects[b506RFSequenceIndex]["key"] = b506RFSequenceKeyJSON
				b506RFSequenceObjects[b506RFSequenceIndex]["ab17_id"] = b506RFSequenceJSON.RawMessage(b506RFSequenceFmt.Sprintf(`"ab17-control-%d"`, b506RFSequenceIndex))
			}
			b506RFSequenceData, b506RFSequenceError := b506RFSequenceJSON.Marshal(b506RFSequenceObjects)
			if b506RFSequenceError != nil {
				b506RFSequenceT.Fatal(b506RFSequenceError)
			}
			if _, b506RFSequenceError := b506RFSequenceEntriesFrom(string(b506RFSequenceData)); b506RFSequenceError != nil {
				b506RFSequenceT.Fatalf("distinct identity rejected: %v", b506RFSequenceError)
			}
			for b506RFSequenceIndex, b506RFSequenceKey := range []B506RFSequenceKey{b506RFSequenceControl, b506RFSequenceChanged} {
				b506RFSequenceRaw, b506RFSequenceFound, b506RFSequenceError := b506RFSequenceLookupFrom(string(b506RFSequenceData), b506RFSequenceKey)
				if b506RFSequenceError != nil || !b506RFSequenceFound {
					b506RFSequenceT.Fatalf("valid control: %v", b506RFSequenceError)
				}
				b506RFSequenceWant(b506RFSequenceT, b506RFSequenceDecode(b506RFSequenceT, b506RFSequenceRaw), "ab17_id", b506RFSequenceFmt.Sprintf(`"ab17-control-%d"`, b506RFSequenceIndex))
			}
		})
	}
}

func TestAb17Validation(b506RFSequenceT *b506RFSequenceTesting.T) {
	b506RFSequenceBase := b506RFSequenceRows(b506RFSequenceT)[0]
	b506RFSequenceKey := b506RFSequenceKeyOf(b506RFSequenceT, b506RFSequenceBase)
	for b506RFSequenceIndex, b506RFSequenceData := range []string{"{", "null", "{}", "[null]", "[{}]", string(b506RFSequenceBase), "[] false", "[" + string(b506RFSequenceBase) + "," + string(b506RFSequenceBase) + "]"} {
		b506RFSequenceT.Run(b506RFSequenceFmt.Sprintf("ab17-invalid-%d", b506RFSequenceIndex), func(b506RFSequenceT *b506RFSequenceTesting.T) {
			if _, b506RFSequenceError := b506RFSequenceEntriesFrom(b506RFSequenceData); b506RFSequenceError == nil {
				b506RFSequenceT.Fatal("invalid catalog admitted")
			}
			if b506RFSequenceRaw, b506RFSequenceFound, b506RFSequenceError := b506RFSequenceLookupFrom(b506RFSequenceData, b506RFSequenceKey); b506RFSequenceError == nil || b506RFSequenceFound || b506RFSequenceRaw != nil {
				b506RFSequenceT.Fatal("lookup did not fail closed")
			}
		})
	}
	for _, b506RFSequenceAxis := range []string{"family", "variant", "primary", "secondary", "operation", "direction", "field", "alternative", "group", "instance_scope", "register"} {
		b506RFSequenceT.Run("ab17-missing-"+b506RFSequenceAxis, func(b506RFSequenceT *b506RFSequenceTesting.T) {
			b506RFSequenceObject := b506RFSequenceDecode(b506RFSequenceT, b506RFSequenceBase)
			b506RFSequenceAxes := b506RFSequenceDecode(b506RFSequenceT, b506RFSequenceObject["key"])
			delete(b506RFSequenceAxes, b506RFSequenceAxis)
			b506RFSequenceKeyJSON, b506RFSequenceError := b506RFSequenceJSON.Marshal(b506RFSequenceAxes)
			if b506RFSequenceError != nil {
				b506RFSequenceT.Fatal(b506RFSequenceError)
			}
			b506RFSequenceObject["key"] = b506RFSequenceKeyJSON
			b506RFSequenceData, b506RFSequenceError := b506RFSequenceJSON.Marshal([]any{b506RFSequenceObject})
			if b506RFSequenceError != nil {
				b506RFSequenceT.Fatal(b506RFSequenceError)
			}
			if _, b506RFSequenceError := b506RFSequenceEntriesFrom(string(b506RFSequenceData)); b506RFSequenceError == nil {
				b506RFSequenceT.Fatal("missing identity admitted")
			}
		})
	}
}

func TestAb17IndependentBuffers(b506RFSequenceT *b506RFSequenceTesting.T) {
	b506RFSequenceBaseline := b506RFSequenceRows(b506RFSequenceT)[0]
	b506RFSequenceKey := b506RFSequenceKeyOf(b506RFSequenceT, b506RFSequenceBaseline)
	b506RFSequenceFirst, b506RFSequenceError := B506RFSequenceEntries()
	if b506RFSequenceError != nil {
		b506RFSequenceT.Fatal(b506RFSequenceError)
	}
	b506RFSequenceSecond, b506RFSequenceError := B506RFSequenceEntries()
	if b506RFSequenceError != nil {
		b506RFSequenceT.Fatal(b506RFSequenceError)
	}
	b506RFSequenceSnapshot := append([]byte(nil), b506RFSequenceSecond...)
	for b506RFSequenceIndex := range b506RFSequenceFirst {
		b506RFSequenceFirst[b506RFSequenceIndex] = 'x'
	}
	if !b506RFSequenceBytes.Equal(b506RFSequenceSecond, b506RFSequenceSnapshot) {
		b506RFSequenceT.Fatal("entries buffers share memory")
	}
	b506RFSequenceThird, b506RFSequenceError := B506RFSequenceEntries()
	if b506RFSequenceError != nil || !b506RFSequenceBytes.Equal(b506RFSequenceThird, b506RFSequenceSnapshot) {
		b506RFSequenceT.Fatal("entries mutation persisted")
	}
	b506RFSequenceFirstRow, b506RFSequenceFound, b506RFSequenceError := B506RFSequenceLookup(b506RFSequenceKey)
	if b506RFSequenceError != nil || !b506RFSequenceFound {
		b506RFSequenceT.Fatal("lookup")
	}
	b506RFSequenceSecondRow, b506RFSequenceFound, b506RFSequenceError := B506RFSequenceLookup(b506RFSequenceKey)
	if b506RFSequenceError != nil || !b506RFSequenceFound {
		b506RFSequenceT.Fatal("lookup")
	}
	for b506RFSequenceIndex := range b506RFSequenceFirstRow {
		b506RFSequenceFirstRow[b506RFSequenceIndex] = 'x'
	}
	if !b506RFSequenceBytes.Equal(b506RFSequenceSecondRow, b506RFSequenceBaseline) {
		b506RFSequenceT.Fatal("row buffers share memory")
	}
	b506RFSequenceThirdRow, b506RFSequenceFound, b506RFSequenceError := B506RFSequenceLookup(b506RFSequenceKey)
	if b506RFSequenceError != nil || !b506RFSequenceFound || !b506RFSequenceBytes.Equal(b506RFSequenceThirdRow, b506RFSequenceBaseline) {
		b506RFSequenceT.Fatal("row mutation persisted")
	}
	b506RFSequenceFourth, b506RFSequenceError := B506RFSequenceEntries()
	if b506RFSequenceError != nil || !b506RFSequenceBytes.Equal(b506RFSequenceFourth, b506RFSequenceSnapshot) {
		b506RFSequenceT.Fatal("row mutation affected entries")
	}
}
