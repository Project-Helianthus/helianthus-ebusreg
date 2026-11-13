package vaillant

import (
	b516ResidualFieldBytes "bytes"
	b516ResidualFieldBinary "encoding/binary"
	b516ResidualFieldJSON "encoding/json"
	b516ResidualFieldReflect "reflect"
	b516ResidualFieldStrings "strings"
	b516ResidualFieldTesting "testing"
)

func b516ResidualFieldPtr(b516ResidualFieldValue uint16) *uint16 { return &b516ResidualFieldValue }

func b516ResidualFieldKey(b516ResidualFieldOperation uint16, b516ResidualFieldField string) B516ResidualFieldKey {
	return B516ResidualFieldKey{B516ResidualFieldFamily: "Vaillant", B516ResidualFieldVariant: "Unknown", B516ResidualFieldPrimary: 181, B516ResidualFieldSecondary: b516ResidualFieldPtr(22), B516ResidualFieldOperation: b516ResidualFieldPtr(b516ResidualFieldOperation), B516ResidualFieldDirection: "Unknown", B516ResidualField: b516ResidualFieldField, B516ResidualFieldAlternative: "none", B516ResidualFieldInstanceScope: "none"}
}

func b516ResidualFieldObject(b516ResidualFieldT *b516ResidualFieldTesting.T, b516ResidualFieldRaw []byte) map[string]any {
	b516ResidualFieldT.Helper()
	var b516ResidualFieldValue map[string]any
	if b516ResidualFieldError := b516ResidualFieldJSON.Unmarshal(b516ResidualFieldRaw, &b516ResidualFieldValue); b516ResidualFieldError != nil {
		b516ResidualFieldT.Fatal(b516ResidualFieldError)
	}
	return b516ResidualFieldValue
}

func b516ResidualFieldGet(b516ResidualFieldT *b516ResidualFieldTesting.T, b516ResidualFieldOperation uint16, b516ResidualFieldField string) map[string]any {
	b516ResidualFieldT.Helper()
	b516ResidualFieldRaw, b516ResidualFieldFound, b516ResidualFieldError := B516ResidualFieldLookup(b516ResidualFieldKey(b516ResidualFieldOperation, b516ResidualFieldField))
	if b516ResidualFieldError != nil || !b516ResidualFieldFound {
		b516ResidualFieldT.Fatalf("lookup %s: %t %v", b516ResidualFieldField, b516ResidualFieldFound, b516ResidualFieldError)
	}
	return b516ResidualFieldObject(b516ResidualFieldT, b516ResidualFieldRaw)
}

func b516ResidualFieldEqual(b516ResidualFieldT *b516ResidualFieldTesting.T, b516ResidualFieldActual, b516ResidualFieldExpected any) {
	b516ResidualFieldT.Helper()
	b516ResidualFieldWant, b516ResidualFieldError := b516ResidualFieldJSON.Marshal(b516ResidualFieldExpected)
	if b516ResidualFieldError != nil {
		b516ResidualFieldT.Fatal(b516ResidualFieldError)
	}
	b516ResidualFieldGot, b516ResidualFieldError := b516ResidualFieldJSON.Marshal(b516ResidualFieldActual)
	if b516ResidualFieldError != nil || !b516ResidualFieldBytes.Equal(b516ResidualFieldGot, b516ResidualFieldWant) {
		b516ResidualFieldT.Fatalf("got %s want %s (%v)", b516ResidualFieldGot, b516ResidualFieldWant, b516ResidualFieldError)
	}
}

func TestAb51FiniteMetadata(b516ResidualFieldT *b516ResidualFieldTesting.T) {
	b516ResidualFieldSource, b516ResidualFieldError := B516ResidualFieldEntries()
	if b516ResidualFieldError != nil {
		b516ResidualFieldT.Fatal(b516ResidualFieldError)
	}
	var b516ResidualFieldRows []b516ResidualFieldJSON.RawMessage
	if b516ResidualFieldError = b516ResidualFieldJSON.Unmarshal(b516ResidualFieldSource, &b516ResidualFieldRows); b516ResidualFieldError != nil || len(b516ResidualFieldRows) != 20 {
		b516ResidualFieldT.Fatalf("finite array: %d %v", len(b516ResidualFieldRows), b516ResidualFieldError)
	}
	b516ResidualFieldSequence := []string{"setpoint", "mode", "special_function", "message_1", "message_2", "message_3", "message_4", "message_5", "minute", "hour", "day", "month", "year"}
	for _, b516ResidualFieldRow := range b516ResidualFieldRows {
		b516ResidualFieldValue := b516ResidualFieldObject(b516ResidualFieldT, b516ResidualFieldRow)
		var b516ResidualFieldIdentity struct {
			B516ResidualFieldKey B516ResidualFieldKey `json:"key"`
		}
		if b516ResidualFieldError = b516ResidualFieldJSON.Unmarshal(b516ResidualFieldRow, &b516ResidualFieldIdentity); b516ResidualFieldError != nil {
			b516ResidualFieldT.Fatal(b516ResidualFieldError)
		}
		b516ResidualFieldLookupRow, b516ResidualFieldFound, b516ResidualFieldLookupError := B516ResidualFieldLookup(b516ResidualFieldIdentity.B516ResidualFieldKey)
		if !b516ResidualFieldFound || b516ResidualFieldLookupError != nil || !b516ResidualFieldBytes.Equal(b516ResidualFieldLookupRow, b516ResidualFieldRow) {
			b516ResidualFieldT.Fatal("whole exact row lost", b516ResidualFieldLookupError)
		}
		for _, b516ResidualFieldName := range []string{"physical_qualification", "presence", "freshness", "execution", "position_grade", "variant_association_grade"} {
			b516ResidualFieldEqual(b516ResidualFieldT, b516ResidualFieldValue[b516ResidualFieldName], "Unknown")
		}
		b516ResidualFieldEqual(b516ResidualFieldT, b516ResidualFieldValue["grade"], "Hypothesis")
		b516ResidualFieldEqual(b516ResidualFieldT, b516ResidualFieldValue["metadata_presence"], "Proven")
		b516ResidualFieldEqual(b516ResidualFieldT, b516ResidualFieldValue["position"], nil)
		if len(b516ResidualFieldValue["unresolved"].([]any)) == 0 || len(b516ResidualFieldValue["excluded"].([]any)) != 3 {
			b516ResidualFieldT.Fatal("unresolved/exclusion boundary lost")
		}
		b516ResidualFieldEqual(b516ResidualFieldT, b516ResidualFieldValue["excluded"], []string{"Passive energy, OP14 aggregation, time-positive alternatives, energy-type flags, meter/clear enums and physical identity are separate.", "OP03 decimal software-version and identity-bearing structures have no supplied fields and are not instantiated.", "Voltage conversion is variant-specific and Unknown here."})
		if !b516ResidualFieldStrings.Contains(b516ResidualFieldValue["native_contract"].(string), "no native promotion") {
			b516ResidualFieldT.Fatal("native boundary lost")
		}
		b516ResidualFieldWidth, b516ResidualFieldWidthGrade := any(nil), "Unknown"
		if b516ResidualFieldIdentity.B516ResidualFieldKey.B516ResidualField == "setpoint" || b516ResidualFieldStrings.HasPrefix(b516ResidualFieldIdentity.B516ResidualFieldKey.B516ResidualField, "message_") {
			b516ResidualFieldWidth, b516ResidualFieldWidthGrade = float64(8), "Hypothesis"
		}
		if *b516ResidualFieldIdentity.B516ResidualFieldKey.B516ResidualFieldOperation == 13 {
			b516ResidualFieldWidth, b516ResidualFieldWidthGrade = float64(16), "Hypothesis"
		}
		b516ResidualFieldEqual(b516ResidualFieldT, b516ResidualFieldValue["width_bits"], b516ResidualFieldWidth)
		b516ResidualFieldEqual(b516ResidualFieldT, b516ResidualFieldValue["width_grade"], b516ResidualFieldWidthGrade)
		b516ResidualFieldUnit := any(nil)
		if b516ResidualFieldIdentity.B516ResidualFieldKey.B516ResidualField == "setpoint" {
			b516ResidualFieldUnit = "degrees Celsius"
		}
		b516ResidualFieldEqual(b516ResidualFieldT, b516ResidualFieldValue["unit"], b516ResidualFieldUnit)
	}
	for _, b516ResidualFieldField := range b516ResidualFieldSequence {
		b516ResidualFieldValue := b516ResidualFieldGet(b516ResidualFieldT, 8, b516ResidualFieldField)
		b516ResidualFieldEqual(b516ResidualFieldT, b516ResidualFieldValue["sequence"], b516ResidualFieldSequence)
		if b516ResidualFieldStrings.HasPrefix(b516ResidualFieldField, "message_") {
			b516ResidualFieldEqual(b516ResidualFieldT, b516ResidualFieldValue["independent"], true)
		}
	}
	b516ResidualFieldEqual(b516ResidualFieldT, b516ResidualFieldGet(b516ResidualFieldT, 8, "setpoint")["raw_form"], map[string]any{"encoding": "unsigned", "divide_by": 2, "preserve": "original bytes"})
	for _, b516ResidualFieldField := range []string{"mode", "special_function", "message_1", "message_2", "message_3", "message_4", "message_5"} {
		b516ResidualFieldEqual(b516ResidualFieldT, b516ResidualFieldGet(b516ResidualFieldT, 8, b516ResidualFieldField)["raw_form"], map[string]any{"encoding": "opaque", "preserve": "original bytes"})
	}
	for _, b516ResidualFieldField := range []string{"page", "data"} {
		b516ResidualFieldValue := b516ResidualFieldGet(b516ResidualFieldT, 14, b516ResidualFieldField)
		b516ResidualFieldEqual(b516ResidualFieldT, b516ResidualFieldValue["maximum_length"], nil)
		b516ResidualFieldEqual(b516ResidualFieldT, b516ResidualFieldValue["raw_form"], map[string]any{"encoding": "opaque", "preserve": "original bytes"})
	}
	b516ResidualFieldValue := b516ResidualFieldGet(b516ResidualFieldT, 14, "data")
	b516ResidualFieldEqual(b516ResidualFieldT, b516ResidualFieldValue["variable"], true)
	b516ResidualFieldEqual(b516ResidualFieldT, b516ResidualFieldValue["empty_metadata"], true)
	b516ResidualFieldEqual(b516ResidualFieldT, b516ResidualFieldValue["peer_acceptance"], "Unknown")
	b516ResidualFieldValue = b516ResidualFieldGet(b516ResidualFieldT, 2, "unresolved_time_associations")
	b516ResidualFieldEqual(b516ResidualFieldT, b516ResidualFieldValue["association_grade"], "Unknown")
	b516ResidualFieldEqual(b516ResidualFieldT, b516ResidualFieldValue["conflicts"], map[string]any{
		"minute": "minute metadata on a second field conflicts with unprovided field associations",
		"day":    map[string]any{"default": 0, "calendar_domain": []int{1, 31}, "status": "out-of-domain calendar conflict", "association_grade": "Unknown"},
		"month":  map[string]any{"default": 0, "calendar_domain": []int{1, 12}, "status": "out-of-domain calendar conflict", "association_grade": "Unknown"},
	})
}

func b516ResidualFieldTime(b516ResidualFieldT *b516ResidualFieldTesting.T, b516ResidualFieldMetadata map[string]any, b516ResidualFieldRaw byte) (int, string) {
	b516ResidualFieldT.Helper()
	b516ResidualFieldForm := b516ResidualFieldMetadata["raw_form"].(map[string]any)
	b516ResidualFieldEqual(b516ResidualFieldT, b516ResidualFieldForm, map[string]any{"encoding": "packed-decimal", "selection_order": []string{"raw FF replacement", "decimal validation", "decimal conversion"}, "replacement_raw": 255, "invalid": "retain raw", "preserve": "original bytes"})
	if int(b516ResidualFieldRaw) == int(b516ResidualFieldForm["replacement_raw"].(float64)) {
		return 0, "replacement"
	}
	if b516ResidualFieldRaw>>4 > 9 || b516ResidualFieldRaw&15 > 9 {
		return int(b516ResidualFieldRaw), "raw"
	}
	return int(b516ResidualFieldRaw>>4)*10 + int(b516ResidualFieldRaw&15), "value"
}

func TestAb51SyntheticTime(b516ResidualFieldT *b516ResidualFieldTesting.T) {
	b516ResidualFieldOrder := []string{"minute", "hour", "day", "month", "year"}
	b516ResidualFieldVector := []byte{0x12, 0x23, 0x31, 0x12, 0x99}
	b516ResidualFieldWant := []int{12, 23, 31, 12, 99}
	for b516ResidualFieldIndex, b516ResidualFieldField := range b516ResidualFieldOrder {
		b516ResidualFieldMetadata := b516ResidualFieldGet(b516ResidualFieldT, 8, b516ResidualFieldField)
		b516ResidualFieldEqual(b516ResidualFieldT, b516ResidualFieldMetadata["time_order"], b516ResidualFieldOrder)
		b516ResidualFieldEqual(b516ResidualFieldT, b516ResidualFieldMetadata["century"], nil)
		b516ResidualFieldEqual(b516ResidualFieldT, b516ResidualFieldMetadata["timezone"], nil)
		for _, b516ResidualFieldCase := range []struct {
			B516ResidualFieldRaw   byte
			B516ResidualFieldWant  int
			B516ResidualFieldState string
		}{{0x12, 12, "value"}, {0xff, 0, "replacement"}, {0xfa, 250, "raw"}, {b516ResidualFieldVector[b516ResidualFieldIndex], b516ResidualFieldWant[b516ResidualFieldIndex], "value"}} {
			b516ResidualFieldGot, b516ResidualFieldState := b516ResidualFieldTime(b516ResidualFieldT, b516ResidualFieldMetadata, b516ResidualFieldCase.B516ResidualFieldRaw)
			if b516ResidualFieldGot != b516ResidualFieldCase.B516ResidualFieldWant || b516ResidualFieldState != b516ResidualFieldCase.B516ResidualFieldState {
				b516ResidualFieldT.Fatal("time rule changed", b516ResidualFieldField, b516ResidualFieldGot, b516ResidualFieldState)
			}
		}
	}
	b516ResidualFieldEqual(b516ResidualFieldT, b516ResidualFieldVector, []byte{0x12, 0x23, 0x31, 0x12, 0x99})
}

func TestAb51SyntheticWords(b516ResidualFieldT *b516ResidualFieldTesting.T) {
	b516ResidualFieldNames := []string{"word_1", "word_2", "word_3", "word_4"}
	b516ResidualFieldRaw := []byte{0, 0, 1, 0, 0, 1, 255, 255}
	b516ResidualFieldWant := []uint16{0, 1, 256, 65535}
	for b516ResidualFieldIndex, b516ResidualFieldName := range b516ResidualFieldNames {
		b516ResidualFieldMetadata := b516ResidualFieldGet(b516ResidualFieldT, 13, b516ResidualFieldName)
		b516ResidualFieldEqual(b516ResidualFieldT, b516ResidualFieldMetadata["raw_form"], map[string]any{"encoding": "unsigned", "byte_order": "little-endian", "preserve": "original bytes"})
		b516ResidualFieldEqual(b516ResidualFieldT, b516ResidualFieldMetadata["independent"], true)
		b516ResidualFieldEqual(b516ResidualFieldT, b516ResidualFieldMetadata["label"], nil)
		b516ResidualFieldEqual(b516ResidualFieldT, b516ResidualFieldMetadata["sentinel"], nil)
		if b516ResidualFieldBinary.LittleEndian.Uint16(b516ResidualFieldRaw[2*b516ResidualFieldIndex:]) != b516ResidualFieldWant[b516ResidualFieldIndex] {
			b516ResidualFieldT.Fatal("independent opaque words")
		}
	}
	b516ResidualFieldEqual(b516ResidualFieldT, b516ResidualFieldRaw, []byte{0, 0, 1, 0, 0, 1, 255, 255})
}

func b516ResidualFieldEncode(b516ResidualFieldT *b516ResidualFieldTesting.T, b516ResidualFieldValue any) []byte {
	b516ResidualFieldT.Helper()
	b516ResidualFieldRaw, b516ResidualFieldError := b516ResidualFieldJSON.Marshal(b516ResidualFieldValue)
	if b516ResidualFieldError != nil {
		b516ResidualFieldT.Fatal(b516ResidualFieldError)
	}
	return b516ResidualFieldRaw
}

func TestAb51EveryIdentityAxis(b516ResidualFieldT *b516ResidualFieldTesting.T) {
	b516ResidualFieldControl := b516ResidualFieldKey(8, "minute")
	for b516ResidualFieldIndex := 0; b516ResidualFieldIndex < b516ResidualFieldReflect.TypeOf(b516ResidualFieldControl).NumField(); b516ResidualFieldIndex++ {
		b516ResidualFieldAxis := b516ResidualFieldReflect.TypeOf(b516ResidualFieldControl).Field(b516ResidualFieldIndex)
		b516ResidualFieldT.Run("ab51_"+b516ResidualFieldAxis.Name, func(b516ResidualFieldT *b516ResidualFieldTesting.T) {
			b516ResidualFieldCandidates := []B516ResidualFieldKey{}
			if b516ResidualFieldAxis.Type.Kind() == b516ResidualFieldReflect.Pointer {
				for _, b516ResidualFieldOption := range []*uint16{nil, b516ResidualFieldPtr(0), b516ResidualFieldPtr(123)} {
					b516ResidualFieldCandidate := b516ResidualFieldControl
					b516ResidualFieldReflect.ValueOf(&b516ResidualFieldCandidate).Elem().Field(b516ResidualFieldIndex).Set(b516ResidualFieldReflect.ValueOf(b516ResidualFieldOption))
					b516ResidualFieldCandidates = append(b516ResidualFieldCandidates, b516ResidualFieldCandidate)
				}
			} else {
				b516ResidualFieldCandidate := b516ResidualFieldControl
				b516ResidualFieldValue := b516ResidualFieldReflect.ValueOf(&b516ResidualFieldCandidate).Elem().Field(b516ResidualFieldIndex)
				if b516ResidualFieldValue.Kind() == b516ResidualFieldReflect.String {
					b516ResidualFieldValue.SetString("ab51-other")
				} else {
					b516ResidualFieldValue.SetUint(0)
				}
				b516ResidualFieldCandidates = append(b516ResidualFieldCandidates, b516ResidualFieldCandidate)
			}
			for _, b516ResidualFieldCandidate := range b516ResidualFieldCandidates {
				_, b516ResidualFieldFound, b516ResidualFieldError := B516ResidualFieldLookup(b516ResidualFieldCandidate)
				if b516ResidualFieldError != nil || b516ResidualFieldFound != b516ResidualFieldReflect.DeepEqual(b516ResidualFieldControl, b516ResidualFieldCandidate) {
					b516ResidualFieldT.Fatal("production axis ignored", b516ResidualFieldAxis.Name, b516ResidualFieldError)
				}
				b516ResidualFieldSource := b516ResidualFieldEncode(b516ResidualFieldT, []any{map[string]any{"key": b516ResidualFieldCandidate, "ab51_fixture": "ab51-axis-control"}})
				for _, b516ResidualFieldQuery := range append([]B516ResidualFieldKey{b516ResidualFieldControl}, b516ResidualFieldCandidates...) {
					b516ResidualFieldRaw, b516ResidualFieldFound, b516ResidualFieldError := b516ResidualFieldLookup(b516ResidualFieldSource, b516ResidualFieldQuery)
					if b516ResidualFieldError != nil || b516ResidualFieldFound != b516ResidualFieldReflect.DeepEqual(b516ResidualFieldCandidate, b516ResidualFieldQuery) {
						b516ResidualFieldT.Fatal("seam axis ignored", b516ResidualFieldAxis.Name, b516ResidualFieldError)
					}
					if b516ResidualFieldFound {
						b516ResidualFieldEqual(b516ResidualFieldT, b516ResidualFieldObject(b516ResidualFieldT, b516ResidualFieldRaw)["ab51_fixture"], "ab51-axis-control")
					}
				}
			}
		})
	}
}

func TestAb51ValidationAndCopies(b516ResidualFieldT *b516ResidualFieldTesting.T) {
	b516ResidualFieldOriginal, b516ResidualFieldError := B516ResidualFieldEntries()
	if b516ResidualFieldError != nil {
		b516ResidualFieldT.Fatal(b516ResidualFieldError)
	}
	b516ResidualFieldKeyValue := b516ResidualFieldKey(8, "minute")
	b516ResidualFieldOriginalRow, b516ResidualFieldFound, b516ResidualFieldError := B516ResidualFieldLookup(b516ResidualFieldKeyValue)
	if !b516ResidualFieldFound || b516ResidualFieldError != nil {
		b516ResidualFieldT.Fatal(b516ResidualFieldError)
	}
	b516ResidualFieldCopy, b516ResidualFieldError := B516ResidualFieldEntries()
	if b516ResidualFieldError != nil {
		b516ResidualFieldT.Fatal(b516ResidualFieldError)
	}
	for b516ResidualFieldIndex := range b516ResidualFieldCopy {
		b516ResidualFieldCopy[b516ResidualFieldIndex] = 'x'
	}
	b516ResidualFieldFresh, b516ResidualFieldError := B516ResidualFieldEntries()
	if b516ResidualFieldError != nil || !b516ResidualFieldBytes.Equal(b516ResidualFieldOriginal, b516ResidualFieldFresh) {
		b516ResidualFieldT.Fatal("entries alias", b516ResidualFieldError)
	}
	b516ResidualFieldSource := append([]byte(nil), b516ResidualFieldOriginal...)
	b516ResidualFieldReturned, b516ResidualFieldFound, b516ResidualFieldError := b516ResidualFieldLookup(b516ResidualFieldSource, b516ResidualFieldKeyValue)
	if !b516ResidualFieldFound || b516ResidualFieldError != nil {
		b516ResidualFieldT.Fatal(b516ResidualFieldError)
	}
	for b516ResidualFieldIndex := range b516ResidualFieldSource {
		b516ResidualFieldSource[b516ResidualFieldIndex] = 'x'
	}
	if !b516ResidualFieldBytes.Equal(b516ResidualFieldReturned, b516ResidualFieldOriginalRow) {
		b516ResidualFieldT.Fatal("source alias")
	}
	for b516ResidualFieldIndex := range b516ResidualFieldReturned {
		b516ResidualFieldReturned[b516ResidualFieldIndex] = 'x'
	}
	b516ResidualFieldReturned, b516ResidualFieldFound, b516ResidualFieldError = B516ResidualFieldLookup(b516ResidualFieldKeyValue)
	if !b516ResidualFieldFound || b516ResidualFieldError != nil || !b516ResidualFieldBytes.Equal(b516ResidualFieldReturned, b516ResidualFieldOriginalRow) {
		b516ResidualFieldT.Fatal("lookup alias", b516ResidualFieldError)
	}
	b516ResidualFieldReturned[0] = 'x'
	b516ResidualFieldReturned, b516ResidualFieldFound, b516ResidualFieldError = B516ResidualFieldLookup(b516ResidualFieldKeyValue)
	if !b516ResidualFieldFound || b516ResidualFieldError != nil || !b516ResidualFieldBytes.Equal(b516ResidualFieldReturned, b516ResidualFieldOriginalRow) {
		b516ResidualFieldT.Fatal("public lookup alias", b516ResidualFieldError)
	}
	b516ResidualFieldDuplicate := b516ResidualFieldEncode(b516ResidualFieldT, []b516ResidualFieldJSON.RawMessage{b516ResidualFieldOriginalRow, b516ResidualFieldOriginalRow})
	for _, b516ResidualFieldBad := range [][]byte{nil, []byte("{"), []byte("null"), []byte("{}"), []byte("[null]"), []byte("[{}]"), []byte("[1]"), []byte("[] []"), b516ResidualFieldDuplicate} {
		if _, _, b516ResidualFieldError = b516ResidualFieldValidate(b516ResidualFieldBad); b516ResidualFieldError == nil {
			b516ResidualFieldT.Fatal("invalid source accepted", string(b516ResidualFieldBad))
		}
		if b516ResidualFieldRaw, b516ResidualFieldFound, b516ResidualFieldError := b516ResidualFieldLookup(b516ResidualFieldBad, B516ResidualFieldKey{}); b516ResidualFieldError == nil || b516ResidualFieldFound || b516ResidualFieldRaw != nil {
			b516ResidualFieldT.Fatal("invalid source became absence/success")
		}
	}
	for _, b516ResidualFieldName := range []string{"family", "variant", "primary", "secondary", "operation", "direction", "field", "alternative", "group", "instance_scope", "register"} {
		b516ResidualFieldObjectValue := b516ResidualFieldObject(b516ResidualFieldT, b516ResidualFieldOriginalRow)
		b516ResidualFieldMembers := b516ResidualFieldObjectValue["key"].(map[string]any)
		b516ResidualFieldMembers[b516ResidualFieldName] = []any{"ab51-invalid-type"}
		if _, _, b516ResidualFieldError = b516ResidualFieldValidate(b516ResidualFieldEncode(b516ResidualFieldT, []any{b516ResidualFieldObjectValue})); b516ResidualFieldError == nil {
			b516ResidualFieldT.Fatal("invalid axis accepted", b516ResidualFieldName)
		}
	}
	b516ResidualFieldObjectValue := b516ResidualFieldObject(b516ResidualFieldT, b516ResidualFieldOriginalRow)
	b516ResidualFieldObjectValue["key"].(map[string]any)["ab51_unknown_axis"] = 0
	if _, _, b516ResidualFieldError = b516ResidualFieldValidate(b516ResidualFieldEncode(b516ResidualFieldT, []any{b516ResidualFieldObjectValue})); b516ResidualFieldError == nil {
		b516ResidualFieldT.Fatal("unknown namespace axis accepted")
	}
	b516ResidualFieldMissing := b516ResidualFieldKeyValue
	b516ResidualFieldMissing.B516ResidualField = "ab51-absent"
	if b516ResidualFieldRaw, b516ResidualFieldFound, b516ResidualFieldError := B516ResidualFieldLookup(b516ResidualFieldMissing); b516ResidualFieldError != nil || b516ResidualFieldFound || b516ResidualFieldRaw != nil {
		b516ResidualFieldT.Fatal("absence contract")
	}
}

func TestAb51WidthAndGradeRemainRowMetadata(b516ResidualFieldT *b516ResidualFieldTesting.T) {
	b516ResidualFieldRaw, b516ResidualFieldFound, b516ResidualFieldError := B516ResidualFieldLookup(b516ResidualFieldKey(13, "word_1"))
	if !b516ResidualFieldFound || b516ResidualFieldError != nil {
		b516ResidualFieldT.Fatal(b516ResidualFieldError)
	}
	for _, b516ResidualFieldAxis := range []string{"width_bits", "grade"} {
		b516ResidualFieldValue := b516ResidualFieldObject(b516ResidualFieldT, b516ResidualFieldRaw)
		if b516ResidualFieldAxis == "width_bits" {
			b516ResidualFieldValue[b516ResidualFieldAxis] = 8
		} else {
			b516ResidualFieldValue[b516ResidualFieldAxis] = "Unknown"
		}
		b516ResidualFieldSource := b516ResidualFieldEncode(b516ResidualFieldT, []any{b516ResidualFieldValue})
		b516ResidualFieldGot, b516ResidualFieldFound, b516ResidualFieldError := b516ResidualFieldLookup(b516ResidualFieldSource, b516ResidualFieldKey(13, "word_1"))
		if !b516ResidualFieldFound || b516ResidualFieldError != nil {
			b516ResidualFieldT.Fatal(b516ResidualFieldError)
		}
		b516ResidualFieldEqual(b516ResidualFieldT, b516ResidualFieldObject(b516ResidualFieldT, b516ResidualFieldGot), b516ResidualFieldValue)
		if _, _, b516ResidualFieldError = b516ResidualFieldValidate(b516ResidualFieldEncode(b516ResidualFieldT, []any{b516ResidualFieldObject(b516ResidualFieldT, b516ResidualFieldRaw), b516ResidualFieldValue})); b516ResidualFieldError == nil {
			b516ResidualFieldT.Fatal("metadata difference hid duplicate identity", b516ResidualFieldAxis)
		}
	}
}

func TestAb51OptionalOmissionAndRequiredAxes(b516ResidualFieldT *b516ResidualFieldTesting.T) {
	b516ResidualFieldRaw, b516ResidualFieldFound, b516ResidualFieldError := B516ResidualFieldLookup(b516ResidualFieldKey(8, "minute"))
	if !b516ResidualFieldFound || b516ResidualFieldError != nil {
		b516ResidualFieldT.Fatal(b516ResidualFieldError)
	}
	for _, b516ResidualFieldAxis := range []string{"secondary", "operation", "group", "register"} {
		b516ResidualFieldValue := b516ResidualFieldObject(b516ResidualFieldT, b516ResidualFieldRaw)
		delete(b516ResidualFieldValue["key"].(map[string]any), b516ResidualFieldAxis)
		var b516ResidualFieldIdentity struct {
			B516ResidualFieldKey B516ResidualFieldKey `json:"key"`
		}
		if b516ResidualFieldError = b516ResidualFieldJSON.Unmarshal(b516ResidualFieldEncode(b516ResidualFieldT, b516ResidualFieldValue), &b516ResidualFieldIdentity); b516ResidualFieldError != nil {
			b516ResidualFieldT.Fatal(b516ResidualFieldError)
		}
		b516ResidualFieldSource := b516ResidualFieldEncode(b516ResidualFieldT, []any{b516ResidualFieldValue})
		if _, b516ResidualFieldFound, b516ResidualFieldError = b516ResidualFieldLookup(b516ResidualFieldSource, b516ResidualFieldIdentity.B516ResidualFieldKey); !b516ResidualFieldFound || b516ResidualFieldError != nil {
			b516ResidualFieldT.Fatal("omitted optional axis lost", b516ResidualFieldAxis, b516ResidualFieldError)
		}
		b516ResidualFieldNull := b516ResidualFieldObject(b516ResidualFieldT, b516ResidualFieldEncode(b516ResidualFieldT, b516ResidualFieldValue))
		b516ResidualFieldNull["key"].(map[string]any)[b516ResidualFieldAxis] = nil
		if _, _, b516ResidualFieldError = b516ResidualFieldValidate(b516ResidualFieldEncode(b516ResidualFieldT, []any{b516ResidualFieldValue, b516ResidualFieldNull})); b516ResidualFieldError == nil {
			b516ResidualFieldT.Fatal("null/omitted duplicate hidden", b516ResidualFieldAxis)
		}
		b516ResidualFieldZero := b516ResidualFieldObject(b516ResidualFieldT, b516ResidualFieldEncode(b516ResidualFieldT, b516ResidualFieldValue))
		b516ResidualFieldZero["key"].(map[string]any)[b516ResidualFieldAxis] = 0
		if _, _, b516ResidualFieldError = b516ResidualFieldValidate(b516ResidualFieldEncode(b516ResidualFieldT, []any{b516ResidualFieldValue, b516ResidualFieldZero})); b516ResidualFieldError != nil {
			b516ResidualFieldT.Fatal("zero collapsed into absence", b516ResidualFieldAxis, b516ResidualFieldError)
		}
	}
	for _, b516ResidualFieldAxis := range []string{"family", "variant", "primary", "direction", "field", "alternative", "instance_scope"} {
		for _, b516ResidualFieldMode := range []string{"ab51-omitted", "ab51-null", "ab51-empty"} {
			if b516ResidualFieldAxis == "primary" && b516ResidualFieldMode == "ab51-empty" {
				continue
			}
			b516ResidualFieldValue := b516ResidualFieldObject(b516ResidualFieldT, b516ResidualFieldRaw)
			b516ResidualFieldMembers := b516ResidualFieldValue["key"].(map[string]any)
			switch b516ResidualFieldMode {
			case "ab51-omitted":
				delete(b516ResidualFieldMembers, b516ResidualFieldAxis)
			case "ab51-null":
				b516ResidualFieldMembers[b516ResidualFieldAxis] = nil
			case "ab51-empty":
				b516ResidualFieldMembers[b516ResidualFieldAxis] = ""
			}
			if _, _, b516ResidualFieldError = b516ResidualFieldValidate(b516ResidualFieldEncode(b516ResidualFieldT, []any{b516ResidualFieldValue})); b516ResidualFieldError == nil {
				b516ResidualFieldT.Fatal("invalid required axis accepted", b516ResidualFieldAxis, b516ResidualFieldMode)
			}
		}
	}
}
