package vaillant

import (
	b50ERuleValueRowBytes "bytes"
	b50ERuleValueRowJSON "encoding/json"
	b50ERuleValueRowFmt "fmt"
	b50ERuleValueRowReflect "reflect"
	b50ERuleValueRowTesting "testing"
)

func b50ERuleValueRowPointer[b50ERuleValueRowType any](b50ERuleValueRowInput b50ERuleValueRowType) *b50ERuleValueRowType {
	return &b50ERuleValueRowInput
}

func b50ERuleValueRowKey(b50ERuleValueRowSelector byte) B50ERuleValueRowKey {
	return B50ERuleValueRowKey{b50ERuleValueRowPointer("vaillant"), b50ERuleValueRowPointer("B50E"), b50ERuleValueRowPointer(byte(0xB5)), b50ERuleValueRowPointer(byte(0x0E)), b50ERuleValueRowPointer(byte(2)), b50ERuleValueRowPointer("positive"), b50ERuleValueRowPointer(b50ERuleValueRowSelector), b50ERuleValueRowPointer("rule_value_row")}
}

func b50ERuleValueRowExpectedField(b50ERuleValueRowOffset int, b50ERuleValueRowWidth *int, b50ERuleValueRowEncoding string) B50ERuleValueRowField {
	return B50ERuleValueRowField{B50ERuleValueRowOffset: b50ERuleValueRowPointer(b50ERuleValueRowOffset), B50ERuleValueRowWidth: b50ERuleValueRowWidth, B50ERuleValueRowEncoding: b50ERuleValueRowEncoding, B50ERuleValueRowByteOrder: "Unknown", B50ERuleValueRowEnumeration: map[string]string{}}
}

func b50ERuleValueRowExpected() []B50ERuleValueRowEntry {
	var b50ERuleValueRowRows []B50ERuleValueRowEntry
	for b50ERuleValueRowIndex, b50ERuleValueRowCase := range []struct {
		b50ERuleValueRowEncoding                                                               string
		b50ERuleValueRowWidth, b50ERuleValueRowCompanionOffset, b50ERuleValueRowCompanionWidth *int
		b50ERuleValueRowCompanionEncoding                                                      string
	}{
		{"opaque", nil, nil, nil, ""},
		{"boolean_u8", b50ERuleValueRowPointer(1), nil, nil, ""},
		{"u8", b50ERuleValueRowPointer(1), nil, nil, ""},
		{"i8", b50ERuleValueRowPointer(1), nil, nil, ""},
		{"u16", b50ERuleValueRowPointer(2), nil, nil, ""},
		{"i16", b50ERuleValueRowPointer(2), nil, nil, ""},
		{"u32", b50ERuleValueRowPointer(4), nil, nil, ""},
		{"i32", b50ERuleValueRowPointer(4), nil, nil, ""},
		{"secure_boolean_u8", b50ERuleValueRowPointer(1), nil, nil, ""},
		{"u8", b50ERuleValueRowPointer(1), b50ERuleValueRowPointer(1), b50ERuleValueRowPointer(1), "raw_inverse"},
		{"i8", b50ERuleValueRowPointer(1), b50ERuleValueRowPointer(1), b50ERuleValueRowPointer(1), "raw_inverse"},
		{"u16", b50ERuleValueRowPointer(2), b50ERuleValueRowPointer(2), b50ERuleValueRowPointer(2), "raw_inverse"},
		{"i16", b50ERuleValueRowPointer(2), b50ERuleValueRowPointer(2), b50ERuleValueRowPointer(2), "raw_inverse"},
		{"u32", b50ERuleValueRowPointer(4), b50ERuleValueRowPointer(4), b50ERuleValueRowPointer(4), "raw_inverse"},
		{"i32", b50ERuleValueRowPointer(4), b50ERuleValueRowPointer(4), b50ERuleValueRowPointer(4), "raw_inverse"},
		{"opaque", nil, nil, nil, ""},
		{"opaque", nil, nil, nil, ""},
		{"i16", b50ERuleValueRowPointer(2), b50ERuleValueRowPointer(2), b50ERuleValueRowPointer(1), "status_u8"},
		{"i16", b50ERuleValueRowPointer(2), b50ERuleValueRowPointer(2), b50ERuleValueRowPointer(2), "raw_inverse"},
	} {
		b50ERuleValueRowRow := B50ERuleValueRowEntry{B50ERuleValueRowKey: b50ERuleValueRowKey(byte(b50ERuleValueRowIndex)), B50ERuleValue: b50ERuleValueRowExpectedField(0, b50ERuleValueRowCase.b50ERuleValueRowWidth, b50ERuleValueRowCase.b50ERuleValueRowEncoding), B50ERuleValueRowBase: "row-relative", B50ERuleValueRowGrade: "Hypothesis", B50ERuleValueRowMetadataPresence: "Proven", B50ERuleValueRowPhysicalSupport: "Unknown", B50ERuleValueRowEffects: "Unknown", B50ERuleValueRowPersistence: "Unknown", B50ERuleValueRowDelivery: "Unknown", B50ERuleValueRowDefaultQualification: "Unknown", B50ERuleValueRowRawFallback: "Preserve isolated row bytes and independently useful fields after partial failure; retain prior good values with uncertain freshness. Inverse bytes are opaque companions, without complement validation. No numeric decoding without byte-order qualification. Opaque end-of-data alternatives have unknown width, never zero width. OP03 buffer bytes are outside this table."}
		if b50ERuleValueRowCase.b50ERuleValueRowCompanionOffset != nil {
			b50ERuleValueRowCompanion := b50ERuleValueRowExpectedField(*b50ERuleValueRowCase.b50ERuleValueRowCompanionOffset, b50ERuleValueRowCase.b50ERuleValueRowCompanionWidth, b50ERuleValueRowCase.b50ERuleValueRowCompanionEncoding)
			b50ERuleValueRowRow.B50ERuleValueRowCompanion = &b50ERuleValueRowCompanion
		}
		if b50ERuleValueRowIndex == 8 {
			b50ERuleValueRowRow.B50ERuleValue.B50ERuleValueRowEnumeration = map[string]string{"0F": "true", "F0": "false"}
		}
		if b50ERuleValueRowIndex >= 17 {
			b50ERuleValueRowRow.B50ERuleValue.B50ERuleValueRowScale = b50ERuleValueRowPointer(1.0 / 16)
			b50ERuleValueRowRow.B50ERuleValue.B50ERuleValueRowUnit = b50ERuleValueRowPointer("Celsius")
			b50ERuleValueRowRow.B50ERuleValue.B50ERuleValueRowCarrierMinimum = b50ERuleValueRowPointer(-32768)
			b50ERuleValueRowRow.B50ERuleValue.B50ERuleValueRowCarrierMaximum = b50ERuleValueRowPointer(32767)
			b50ERuleValueRowRow.B50ERuleValue.B50ERuleValueRowScaledMinimum = b50ERuleValueRowPointer(-2048.0)
			b50ERuleValueRowRow.B50ERuleValue.B50ERuleValueRowScaledMaximum = b50ERuleValueRowPointer(2047.9375)
		}
		if b50ERuleValueRowIndex == 17 {
			b50ERuleValueRowRow.B50ERuleValueRowCompanion.B50ERuleValueRowEnumeration = map[string]string{"00": "okay", "55": "short-circuit", "AA": "interrupted"}
		}
		b50ERuleValueRowRows = append(b50ERuleValueRowRows, b50ERuleValueRowRow)
	}
	return b50ERuleValueRowRows
}

func TestAb74ExactMetadata(b50ERuleValueRowT *b50ERuleValueRowTesting.T) {
	b50ERuleValueRowRows, b50ERuleValueRowError := B50ERuleValueRowEntries()
	if b50ERuleValueRowError != nil {
		b50ERuleValueRowT.Fatal(b50ERuleValueRowError)
	}
	if !b50ERuleValueRowReflect.DeepEqual(b50ERuleValueRowRows, b50ERuleValueRowExpected()) {
		b50ERuleValueRowT.Fatal("ab74 independent nineteen-row contract differs")
	}
	for _, b50ERuleValueRowRow := range b50ERuleValueRowExpected() {
		b50ERuleValueRowFound, b50ERuleValueRowOK, b50ERuleValueRowError := B50ERuleValueRowLookup(b50ERuleValueRowRow.B50ERuleValueRowKey)
		if b50ERuleValueRowError != nil || !b50ERuleValueRowOK || !b50ERuleValueRowReflect.DeepEqual(b50ERuleValueRowFound, b50ERuleValueRowRow) {
			b50ERuleValueRowT.Fatal("ab74 exact control", b50ERuleValueRowError)
		}
	}
}

func TestAb74IdentityAxes(b50ERuleValueRowT *b50ERuleValueRowTesting.T) {
	for b50ERuleValueRowSelector := byte(0); b50ERuleValueRowSelector <= 18; b50ERuleValueRowSelector++ {
		for b50ERuleValueRowAxis := 0; b50ERuleValueRowAxis < b50ERuleValueRowReflect.TypeOf(B50ERuleValueRowKey{}).NumField(); b50ERuleValueRowAxis++ {
			for b50ERuleValueRowMutation := 0; b50ERuleValueRowMutation < 3; b50ERuleValueRowMutation++ {
				b50ERuleValueRowInput := b50ERuleValueRowKey(b50ERuleValueRowSelector)
				b50ERuleValueRowSlot := b50ERuleValueRowReflect.ValueOf(&b50ERuleValueRowInput).Elem().Field(b50ERuleValueRowAxis)
				if b50ERuleValueRowMutation == 0 {
					b50ERuleValueRowSlot.SetZero()
				} else {
					b50ERuleValueRowValue := b50ERuleValueRowReflect.New(b50ERuleValueRowSlot.Type().Elem())
					if b50ERuleValueRowValue.Elem().Kind() == b50ERuleValueRowReflect.String {
						if b50ERuleValueRowMutation == 2 {
							b50ERuleValueRowValue.Elem().SetString("ab74_other")
						}
					} else if b50ERuleValueRowMutation == 2 {
						b50ERuleValueRowValue.Elem().SetUint(19)
						if b50ERuleValueRowAxis == 4 {
							b50ERuleValueRowValue.Elem().SetUint(3)
						}
					}
					b50ERuleValueRowSlot.Set(b50ERuleValueRowValue)
				}
				b50ERuleValueRowFound, b50ERuleValueRowOK, b50ERuleValueRowError := B50ERuleValueRowLookup(b50ERuleValueRowInput)
				b50ERuleValueRowWant := b50ERuleValueRowAxis == 6 && b50ERuleValueRowMutation == 1
				if b50ERuleValueRowError != nil || b50ERuleValueRowOK != b50ERuleValueRowWant || (b50ERuleValueRowWant && !b50ERuleValueRowReflect.DeepEqual(b50ERuleValueRowFound, b50ERuleValueRowExpected()[0])) {
					b50ERuleValueRowT.Fatalf("ab74 selector %d axis %d mutation %d: %v %v", b50ERuleValueRowSelector, b50ERuleValueRowAxis, b50ERuleValueRowMutation, b50ERuleValueRowOK, b50ERuleValueRowError)
				}
			}
		}
	}
}

func b50ERuleValueRowMarshal(b50ERuleValueRowT *b50ERuleValueRowTesting.T, b50ERuleValueRowInput any) []byte {
	b50ERuleValueRowT.Helper()
	b50ERuleValueRowRaw, b50ERuleValueRowError := b50ERuleValueRowJSON.Marshal(b50ERuleValueRowInput)
	if b50ERuleValueRowError != nil {
		b50ERuleValueRowT.Fatal(b50ERuleValueRowError)
	}
	return b50ERuleValueRowRaw
}

func TestAb74RejectInvalidCatalog(b50ERuleValueRowT *b50ERuleValueRowTesting.T) {
	for _, b50ERuleValueRowRaw := range [][]byte{nil, []byte("{"), []byte("null"), []byte("[]"), append(b50ERuleValueRowBytes.Clone(b50ERuleValueRowRowsJSON), []byte(" {}")...)} {
		if b50ERuleValueRowRows, b50ERuleValueRowError := b50ERuleValueRowDecode(b50ERuleValueRowRaw); b50ERuleValueRowError == nil || b50ERuleValueRowRows != nil {
			b50ERuleValueRowT.Fatal("ab74 malformed catalog admitted")
		}
	}
	for _, b50ERuleValueRowMutation := range []string{"ab74_partial", "ab74_extra", "ab74_duplicate", "ab74_unknown", "ab74_zero_width", "ab74_order", "ab74_default", "ab74_base"} {
		b50ERuleValueRowRows := b50ERuleValueRowExpected()
		switch b50ERuleValueRowMutation {
		case "ab74_partial":
			b50ERuleValueRowRows = b50ERuleValueRowRows[:18]
		case "ab74_extra":
			b50ERuleValueRowRows = append(b50ERuleValueRowRows, b50ERuleValueRowRows[0])
		case "ab74_duplicate":
			b50ERuleValueRowRows[1].B50ERuleValueRowKey = b50ERuleValueRowRows[0].B50ERuleValueRowKey
		case "ab74_unknown":
			b50ERuleValueRowRows[0].B50ERuleValueRowKey.B50ERuleValueRowSelector = b50ERuleValueRowPointer(byte(19))
		case "ab74_zero_width":
			b50ERuleValueRowRows[0].B50ERuleValue.B50ERuleValueRowWidth = b50ERuleValueRowPointer(0)
		case "ab74_order":
			b50ERuleValueRowRows[17].B50ERuleValue.B50ERuleValueRowByteOrder = "little-endian"
		case "ab74_default":
			b50ERuleValueRowRows[1].B50ERuleValue.B50ERuleValueRowDefault = b50ERuleValueRowPointer(0.0)
		case "ab74_base":
			b50ERuleValueRowRows[1].B50ERuleValue.B50ERuleValueRowOffset = b50ERuleValueRowPointer(1)
		}
		if b50ERuleValueRowGot, b50ERuleValueRowError := b50ERuleValueRowDecode(b50ERuleValueRowMarshal(b50ERuleValueRowT, b50ERuleValueRowRows)); b50ERuleValueRowError == nil || b50ERuleValueRowGot != nil {
			b50ERuleValueRowT.Fatal(b50ERuleValueRowMutation, "admitted")
		}
	}
	var b50ERuleValueRowRows []map[string]any
	if b50ERuleValueRowError := b50ERuleValueRowJSON.Unmarshal(b50ERuleValueRowRowsJSON, &b50ERuleValueRowRows); b50ERuleValueRowError != nil {
		b50ERuleValueRowT.Fatal(b50ERuleValueRowError)
	}
	for _, b50ERuleValueRowObject := range []map[string]any{b50ERuleValueRowRows[0], b50ERuleValueRowRows[0]["key"].(map[string]any), b50ERuleValueRowRows[0]["value"].(map[string]any)} {
		for b50ERuleValueRowName, b50ERuleValueRowValue := range b50ERuleValueRowObject {
			delete(b50ERuleValueRowObject, b50ERuleValueRowName)
			if _, b50ERuleValueRowError := b50ERuleValueRowDecode(b50ERuleValueRowMarshal(b50ERuleValueRowT, b50ERuleValueRowRows)); b50ERuleValueRowError == nil {
				b50ERuleValueRowT.Fatalf("ab74 missing %s admitted", b50ERuleValueRowName)
			}
			b50ERuleValueRowObject[b50ERuleValueRowName] = b50ERuleValueRowValue
		}
		b50ERuleValueRowObject["ab74_unknown"] = true
		if _, b50ERuleValueRowError := b50ERuleValueRowDecode(b50ERuleValueRowMarshal(b50ERuleValueRowT, b50ERuleValueRowRows)); b50ERuleValueRowError == nil {
			b50ERuleValueRowT.Fatal("ab74 unknown property admitted")
		}
		delete(b50ERuleValueRowObject, "ab74_unknown")
	}
	if _, b50ERuleValueRowError := b50ERuleValueRowDecode(b50ERuleValueRowRowsJSON); b50ERuleValueRowError != nil {
		b50ERuleValueRowT.Fatal("ab74 valid control", b50ERuleValueRowError)
	}
}

func b50ERuleValueRowObserve(b50ERuleValueRowRow B50ERuleValueRowEntry, b50ERuleValueRowRaw []byte) ([]byte, []byte, []byte) {
	b50ERuleValueRowEvidence := b50ERuleValueRowBytes.Clone(b50ERuleValueRowRaw)
	b50ERuleValueRowExtract := func(b50ERuleValueRowField *B50ERuleValueRowField) []byte {
		if b50ERuleValueRowField == nil || b50ERuleValueRowField.B50ERuleValueRowOffset == nil || b50ERuleValueRowField.B50ERuleValueRowWidth == nil || b50ERuleValueRowRow.B50ERuleValueRowBase != "row-relative" {
			return nil
		}
		b50ERuleValueRowStart, b50ERuleValueRowLength := *b50ERuleValueRowField.B50ERuleValueRowOffset, *b50ERuleValueRowField.B50ERuleValueRowWidth
		if b50ERuleValueRowStart < 0 || b50ERuleValueRowLength <= 0 || b50ERuleValueRowStart > len(b50ERuleValueRowRaw) || b50ERuleValueRowLength > len(b50ERuleValueRowRaw)-b50ERuleValueRowStart {
			return nil
		}
		return b50ERuleValueRowBytes.Clone(b50ERuleValueRowRaw[b50ERuleValueRowStart : b50ERuleValueRowStart+b50ERuleValueRowLength])
	}
	return b50ERuleValueRowExtract(&b50ERuleValueRowRow.B50ERuleValue), b50ERuleValueRowExtract(b50ERuleValueRowRow.B50ERuleValueRowCompanion), b50ERuleValueRowEvidence
}

func TestAb74RawCompanions(b50ERuleValueRowT *b50ERuleValueRowTesting.T) {
	for _, b50ERuleValueRowCase := range []struct {
		b50ERuleValueRowSelector                                   byte
		b50ERuleValueRowValueWidth, b50ERuleValueRowCompanionWidth int
	}{
		{9, 1, 1}, {10, 1, 1}, {11, 2, 2}, {12, 2, 2}, {13, 4, 4}, {14, 4, 4}, {18, 2, 2},
	} {
		b50ERuleValueRowRow, b50ERuleValueRowOK, b50ERuleValueRowError := B50ERuleValueRowLookup(b50ERuleValueRowKey(b50ERuleValueRowCase.b50ERuleValueRowSelector))
		if b50ERuleValueRowError != nil || !b50ERuleValueRowOK {
			b50ERuleValueRowT.Fatal(b50ERuleValueRowError)
		}
		for b50ERuleValueRowIndex := 0; b50ERuleValueRowIndex < b50ERuleValueRowCase.b50ERuleValueRowCompanionWidth; b50ERuleValueRowIndex++ {
			b50ERuleValueRowInput := b50ERuleValueRowBytes.Repeat([]byte{0x23}, b50ERuleValueRowCase.b50ERuleValueRowValueWidth+b50ERuleValueRowCase.b50ERuleValueRowCompanionWidth)
			b50ERuleValueRowInput[b50ERuleValueRowCase.b50ERuleValueRowValueWidth+b50ERuleValueRowIndex] = byte(0x40 + b50ERuleValueRowIndex)
			b50ERuleValueRowValue, b50ERuleValueRowCompanion, b50ERuleValueRowEvidence := b50ERuleValueRowObserve(b50ERuleValueRowRow, b50ERuleValueRowInput)
			if !b50ERuleValueRowBytes.Equal(b50ERuleValueRowValue, b50ERuleValueRowInput[:b50ERuleValueRowCase.b50ERuleValueRowValueWidth]) || !b50ERuleValueRowBytes.Equal(b50ERuleValueRowCompanion, b50ERuleValueRowInput[b50ERuleValueRowCase.b50ERuleValueRowValueWidth:]) || !b50ERuleValueRowBytes.Equal(b50ERuleValueRowEvidence, b50ERuleValueRowInput) {
				b50ERuleValueRowT.Fatalf("ab74 row %d companion byte %d lost", b50ERuleValueRowCase.b50ERuleValueRowSelector, b50ERuleValueRowIndex)
			}
			b50ERuleValueRowValue[0], b50ERuleValueRowCompanion[0], b50ERuleValueRowEvidence[0] = 0, 0, 0
			if b50ERuleValueRowInput[0] != 0x23 || b50ERuleValueRowInput[b50ERuleValueRowCase.b50ERuleValueRowValueWidth] == 0 {
				b50ERuleValueRowT.Fatal("ab74 returned bytes alias input")
			}
		}
		for b50ERuleValueRowLength := 0; b50ERuleValueRowLength < b50ERuleValueRowCase.b50ERuleValueRowValueWidth+b50ERuleValueRowCase.b50ERuleValueRowCompanionWidth; b50ERuleValueRowLength++ {
			b50ERuleValueRowInput := b50ERuleValueRowBytes.Repeat([]byte{0x23}, b50ERuleValueRowLength)
			b50ERuleValueRowValue, b50ERuleValueRowCompanion, b50ERuleValueRowEvidence := b50ERuleValueRowObserve(b50ERuleValueRowRow, b50ERuleValueRowInput)
			if (b50ERuleValueRowValue != nil) != (b50ERuleValueRowLength >= b50ERuleValueRowCase.b50ERuleValueRowValueWidth) || b50ERuleValueRowCompanion != nil || !b50ERuleValueRowBytes.Equal(b50ERuleValueRowEvidence, b50ERuleValueRowInput) {
				b50ERuleValueRowT.Fatal("ab74 partial row lost independent field or raw evidence")
			}
		}
		b50ERuleValueRowInput := b50ERuleValueRowBytes.Repeat([]byte{0x23}, b50ERuleValueRowCase.b50ERuleValueRowValueWidth+b50ERuleValueRowCase.b50ERuleValueRowCompanionWidth+1)
		_, _, b50ERuleValueRowEvidence := b50ERuleValueRowObserve(b50ERuleValueRowRow, b50ERuleValueRowInput)
		if !b50ERuleValueRowBytes.Equal(b50ERuleValueRowEvidence, b50ERuleValueRowInput) {
			b50ERuleValueRowT.Fatal("ab74 extra raw byte discarded")
		}
	}
	for _, b50ERuleValueRowSelector := range []byte{0, 15, 16} {
		b50ERuleValueRowRow, _, b50ERuleValueRowError := B50ERuleValueRowLookup(b50ERuleValueRowKey(b50ERuleValueRowSelector))
		if b50ERuleValueRowError != nil {
			b50ERuleValueRowT.Fatal(b50ERuleValueRowError)
		}
		b50ERuleValueRowValue, b50ERuleValueRowCompanion, b50ERuleValueRowEvidence := b50ERuleValueRowObserve(b50ERuleValueRowRow, []byte{0, 0x55, 0xAB})
		if b50ERuleValueRowValue != nil || b50ERuleValueRowCompanion != nil || !b50ERuleValueRowBytes.Equal(b50ERuleValueRowEvidence, []byte{0, 0x55, 0xAB}) {
			b50ERuleValueRowT.Fatal("ab74 opaque row acquired width")
		}
	}
}

func TestAb74FiniteEnumerations(b50ERuleValueRowT *b50ERuleValueRowTesting.T) {
	b50ERuleValueRowBoolean, _, b50ERuleValueRowError := B50ERuleValueRowLookup(b50ERuleValueRowKey(8))
	if b50ERuleValueRowError != nil {
		b50ERuleValueRowT.Fatal(b50ERuleValueRowError)
	}
	b50ERuleValueRowStatus, _, b50ERuleValueRowError := B50ERuleValueRowLookup(b50ERuleValueRowKey(17))
	if b50ERuleValueRowError != nil {
		b50ERuleValueRowT.Fatal(b50ERuleValueRowError)
	}
	for b50ERuleValueRowByte := 0; b50ERuleValueRowByte < 256; b50ERuleValueRowByte++ {
		b50ERuleValueRowLabel := b50ERuleValueRowFmt.Sprintf("%02X", b50ERuleValueRowByte)
		b50ERuleValueRowWantBoolean, b50ERuleValueRowWantStatus := "", ""
		switch b50ERuleValueRowByte {
		case 0x0F:
			b50ERuleValueRowWantBoolean = "true"
		case 0xF0:
			b50ERuleValueRowWantBoolean = "false"
		}
		switch b50ERuleValueRowByte {
		case 0x00:
			b50ERuleValueRowWantStatus = "okay"
		case 0x55:
			b50ERuleValueRowWantStatus = "short-circuit"
		case 0xAA:
			b50ERuleValueRowWantStatus = "interrupted"
		}
		b50ERuleValueRowGotBoolean, b50ERuleValueRowBooleanOK := b50ERuleValueRowBoolean.B50ERuleValue.B50ERuleValueRowEnumeration[b50ERuleValueRowLabel]
		b50ERuleValueRowGotStatus, b50ERuleValueRowStatusOK := b50ERuleValueRowStatus.B50ERuleValueRowCompanion.B50ERuleValueRowEnumeration[b50ERuleValueRowLabel]
		if b50ERuleValueRowGotBoolean != b50ERuleValueRowWantBoolean || b50ERuleValueRowBooleanOK != (b50ERuleValueRowWantBoolean != "") || b50ERuleValueRowGotStatus != b50ERuleValueRowWantStatus || b50ERuleValueRowStatusOK != (b50ERuleValueRowWantStatus != "") {
			b50ERuleValueRowT.Fatalf("ab74 byte %s gained or lost meaning", b50ERuleValueRowLabel)
		}
	}
	b50ERuleValueRowOrdinary, _, b50ERuleValueRowError := B50ERuleValueRowLookup(b50ERuleValueRowKey(1))
	if b50ERuleValueRowError != nil || len(b50ERuleValueRowOrdinary.B50ERuleValue.B50ERuleValueRowEnumeration) != 0 {
		b50ERuleValueRowT.Fatal("ab74 incomplete Boolean gained universal meaning")
	}
	b50ERuleValueRowValue, b50ERuleValueRowCompanion, b50ERuleValueRowEvidence := b50ERuleValueRowObserve(b50ERuleValueRowStatus, []byte{0x80, 0x01, 0xFF})
	if !b50ERuleValueRowBytes.Equal(b50ERuleValueRowValue, []byte{0x80, 0x01}) || !b50ERuleValueRowBytes.Equal(b50ERuleValueRowCompanion, []byte{0xFF}) || !b50ERuleValueRowBytes.Equal(b50ERuleValueRowEvidence, []byte{0x80, 0x01, 0xFF}) || b50ERuleValueRowStatus.B50ERuleValue.B50ERuleValueRowByteOrder != "Unknown" {
		b50ERuleValueRowT.Fatal("ab74 unknown status or byte order lost raw candidate")
	}
}

func b50ERuleValueRowMutateMemory(b50ERuleValueRowInput b50ERuleValueRowReflect.Value) {
	switch b50ERuleValueRowInput.Kind() {
	case b50ERuleValueRowReflect.Pointer:
		if !b50ERuleValueRowInput.IsNil() {
			b50ERuleValueRowMutateMemory(b50ERuleValueRowInput.Elem())
		}
	case b50ERuleValueRowReflect.Struct:
		for b50ERuleValueRowIndex := 0; b50ERuleValueRowIndex < b50ERuleValueRowInput.NumField(); b50ERuleValueRowIndex++ {
			b50ERuleValueRowMutateMemory(b50ERuleValueRowInput.Field(b50ERuleValueRowIndex))
		}
	case b50ERuleValueRowReflect.Map:
		if !b50ERuleValueRowInput.IsNil() {
			b50ERuleValueRowInput.SetMapIndex(b50ERuleValueRowReflect.ValueOf("ab74_mutated"), b50ERuleValueRowReflect.ValueOf("ab74_mutated"))
		}
	case b50ERuleValueRowReflect.String:
		b50ERuleValueRowInput.SetString("ab74_mutated")
	case b50ERuleValueRowReflect.Int:
		b50ERuleValueRowInput.SetInt(99)
	case b50ERuleValueRowReflect.Uint8:
		b50ERuleValueRowInput.SetUint(99)
	case b50ERuleValueRowReflect.Float64:
		b50ERuleValueRowInput.SetFloat(99)
	}
}

func TestAb74DetachedMemory(b50ERuleValueRowT *b50ERuleValueRowTesting.T) {
	b50ERuleValueRowBefore, b50ERuleValueRowError := B50ERuleValueRowEntries()
	if b50ERuleValueRowError != nil {
		b50ERuleValueRowT.Fatal(b50ERuleValueRowError)
	}
	b50ERuleValueRowCopy, b50ERuleValueRowError := B50ERuleValueRowEntries()
	if b50ERuleValueRowError != nil {
		b50ERuleValueRowT.Fatal(b50ERuleValueRowError)
	}
	for b50ERuleValueRowIndex := range b50ERuleValueRowCopy {
		b50ERuleValueRowMutateMemory(b50ERuleValueRowReflect.ValueOf(&b50ERuleValueRowCopy[b50ERuleValueRowIndex]).Elem())
	}
	b50ERuleValueRowAfter, b50ERuleValueRowError := B50ERuleValueRowEntries()
	if b50ERuleValueRowError != nil || !b50ERuleValueRowReflect.DeepEqual(b50ERuleValueRowBefore, b50ERuleValueRowAfter) {
		b50ERuleValueRowT.Fatal("ab74 entries share nested memory", b50ERuleValueRowError)
	}
	for b50ERuleValueRowIndex, b50ERuleValueRowExpected := range b50ERuleValueRowBefore {
		b50ERuleValueRowKeyCopy := b50ERuleValueRowKey(byte(b50ERuleValueRowIndex))
		b50ERuleValueRowFirst, b50ERuleValueRowOK, b50ERuleValueRowError := B50ERuleValueRowLookup(b50ERuleValueRowKeyCopy)
		if b50ERuleValueRowError != nil || !b50ERuleValueRowOK {
			b50ERuleValueRowT.Fatal(b50ERuleValueRowError)
		}
		b50ERuleValueRowMutateMemory(b50ERuleValueRowReflect.ValueOf(&b50ERuleValueRowFirst).Elem())
		b50ERuleValueRowAgain, b50ERuleValueRowOK, b50ERuleValueRowError := B50ERuleValueRowLookup(b50ERuleValueRowKeyCopy)
		if b50ERuleValueRowError != nil || !b50ERuleValueRowOK || !b50ERuleValueRowReflect.DeepEqual(b50ERuleValueRowExpected, b50ERuleValueRowAgain) {
			b50ERuleValueRowT.Fatal("ab74 lookup or caller identity shares nested memory", b50ERuleValueRowError)
		}
	}
	b50ERuleValueRowInput := b50ERuleValueRowBytes.Clone(b50ERuleValueRowRowsJSON)
	b50ERuleValueRowDecoded, b50ERuleValueRowError := b50ERuleValueRowDecode(b50ERuleValueRowInput)
	if b50ERuleValueRowError != nil {
		b50ERuleValueRowT.Fatal(b50ERuleValueRowError)
	}
	clear(b50ERuleValueRowInput)
	if !b50ERuleValueRowReflect.DeepEqual(b50ERuleValueRowDecoded, b50ERuleValueRowBefore) {
		b50ERuleValueRowT.Fatal("ab74 decoded catalog aliases raw JSON")
	}
}
