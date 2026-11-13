package vaillant

import (
	b51AMirrorFieldBinary "encoding/binary"
	b51AMirrorFieldFmt "fmt"
	b51AMirrorFieldReflect "reflect"
	b51AMirrorFieldTesting "testing"
)

func b51AMirrorFieldPtr[B51AMirrorFieldT any](b51AMirrorFieldValue B51AMirrorFieldT) *B51AMirrorFieldT {
	return &b51AMirrorFieldValue
}

func b51AMirrorFieldExpectedField(b51AMirrorFieldOperation uint8, b51AMirrorFieldDirection, b51AMirrorFieldName string, b51AMirrorFieldPosition, b51AMirrorFieldWidth int, b51AMirrorFieldMaximum int64) B51AMirrorField {
	b51AMirrorFieldField := B51AMirrorField{
		B51AMirrorFieldKey:             B51AMirrorKey{B51AMirrorFieldFamily: "vaillant", B51AMirrorFieldVariant: "mirror-diagnosis", B51AMirrorFieldPrimary: 181, B51AMirrorFieldSecondary: b51AMirrorFieldPtr(uint8(26)), B51AMirrorFieldOperation: b51AMirrorFieldPtr(b51AMirrorFieldOperation), B51AMirrorFieldDirection: b51AMirrorFieldDirection, B51AMirrorFieldField: b51AMirrorFieldName},
		B51AMirrorFieldLogicalPosition: b51AMirrorFieldPtr(b51AMirrorFieldPosition), B51AMirrorFieldWidthBytes: b51AMirrorFieldPtr(b51AMirrorFieldWidth), B51AMirrorFieldSignedness: "unsigned", B51AMirrorFieldByteOrder: "not_applicable", B51AMirrorFieldConversion: "identity",
		B51AMirrorFieldFormalMin: b51AMirrorFieldPtr(int64(0)), B51AMirrorFieldFormalMax: b51AMirrorFieldPtr(b51AMirrorFieldMaximum), B51AMirrorFieldEnumeration: map[int64]string{}, B51AMirrorFieldBits: []B51AMirrorBit{}, B51AMirrorFieldAlternatives: []B51AMirrorAlternative{}, B51AMirrorFieldReserved: []int{},
		B51AMirrorFieldGrade: SchemaHypothesis, B51AMirrorFieldRawFallback: "preserve original bytes including unnamed and reserved bits", B51AMirrorFieldCorrelation: "Unknown", B51AMirrorFieldUnmapped: "Unknown",
	}
	if b51AMirrorFieldWidth == 2 {
		b51AMirrorFieldField.B51AMirrorFieldByteOrder = "little_endian"
	}
	return b51AMirrorFieldField
}

func b51AMirrorFieldExpectedBit(b51AMirrorFieldIndex int, b51AMirrorFieldName string, b51AMirrorFieldRole string) B51AMirrorBit {
	b51AMirrorFieldBit := B51AMirrorBit{B51AMirrorFieldIndex: b51AMirrorFieldIndex, B51AMirrorFieldRole: b51AMirrorFieldRole, B51AMirrorFieldDefault: b51AMirrorFieldPtr(false), B51AMirrorFieldGrade: SchemaHypothesis}
	if b51AMirrorFieldName != "" {
		b51AMirrorFieldBit.B51AMirrorFieldName = b51AMirrorFieldPtr(b51AMirrorFieldName)
	}
	if b51AMirrorFieldRole == "unnamed" {
		b51AMirrorFieldBit.B51AMirrorFieldDefault = nil
		b51AMirrorFieldBit.B51AMirrorFieldGrade = SchemaUnknown
	}
	return b51AMirrorFieldBit
}

func b51AMirrorFieldExpectedFields() []B51AMirrorField {
	b51AMirrorFieldFields := []B51AMirrorField{}
	for _, b51AMirrorFieldOperation := range []uint8{0, 1, 2, 3, 6} {
		for _, b51AMirrorFieldDirection := range []string{"request", "positive"} {
			b51AMirrorFieldField := b51AMirrorFieldExpectedField(b51AMirrorFieldOperation, b51AMirrorFieldDirection, "sequence", 3, 1, 255)
			b51AMirrorFieldField.B51AMirrorFieldDefault = b51AMirrorFieldPtr(int64(0))
			if b51AMirrorFieldOperation == 0 || b51AMirrorFieldOperation == 1 {
				b51AMirrorFieldField.B51AMirrorFieldCorrelation = "candidate_equal_sequence"
			}
			b51AMirrorFieldFields = append(b51AMirrorFieldFields, b51AMirrorFieldField)
		}
	}
	b51AMirrorFieldField := b51AMirrorFieldExpectedField(0, "positive", "acknowledgement", 4, 1, 255)
	b51AMirrorFieldField.B51AMirrorFieldDefault = b51AMirrorFieldPtr(int64(0))
	b51AMirrorFieldField.B51AMirrorFieldBits = append(b51AMirrorFieldField.B51AMirrorFieldBits, b51AMirrorFieldExpectedBit(0, "denied", "named"))
	for b51AMirrorFieldIndex := 1; b51AMirrorFieldIndex < 8; b51AMirrorFieldIndex++ {
		b51AMirrorFieldField.B51AMirrorFieldBits = append(b51AMirrorFieldField.B51AMirrorFieldBits, b51AMirrorFieldExpectedBit(b51AMirrorFieldIndex, "", "unnamed"))
	}
	b51AMirrorFieldFields = append(b51AMirrorFieldFields, b51AMirrorFieldField)
	b51AMirrorFieldField = b51AMirrorFieldExpectedField(2, "request", "flags", 4, 1, 255)
	b51AMirrorFieldBit := b51AMirrorFieldExpectedBit(0, "triple-storage", "named")
	b51AMirrorFieldBit.B51AMirrorFieldDefault = nil
	b51AMirrorFieldField.B51AMirrorFieldBits = append(b51AMirrorFieldField.B51AMirrorFieldBits, b51AMirrorFieldBit)
	for b51AMirrorFieldIndex := 1; b51AMirrorFieldIndex < 8; b51AMirrorFieldIndex++ {
		b51AMirrorFieldField.B51AMirrorFieldBits = append(b51AMirrorFieldField.B51AMirrorFieldBits, b51AMirrorFieldExpectedBit(b51AMirrorFieldIndex, "", "unnamed"))
	}
	b51AMirrorFieldFields = append(b51AMirrorFieldFields, b51AMirrorFieldField, b51AMirrorFieldExpectedField(2, "request", "length", 5, 1, 255))
	b51AMirrorFieldField = b51AMirrorFieldExpectedField(2, "request", "mirror_address", 6, 2, 65535)
	b51AMirrorFieldField.B51AMirrorFieldLogicalMin = b51AMirrorFieldPtr(int64(0))
	b51AMirrorFieldField.B51AMirrorFieldLogicalLimit = b51AMirrorFieldPtr("mirror-length - 1")
	b51AMirrorFieldFields = append(b51AMirrorFieldFields, b51AMirrorFieldField)
	b51AMirrorFieldField = b51AMirrorFieldExpectedField(2, "positive", "positive", 4, 1, 255)
	b51AMirrorFieldField.B51AMirrorFieldEnumeration = map[int64]string{0: "false", 1: "true"}
	b51AMirrorFieldField.B51AMirrorFieldUnmapped = "2..255 unresolved"
	b51AMirrorFieldField.B51AMirrorFieldAlternatives = []B51AMirrorAlternative{{B51AMirrorFieldName: "competing out-of-area bit", B51AMirrorFieldWidthBytes: b51AMirrorFieldPtr(1), B51AMirrorFieldCondition: "unresolved alternative to whole-u8 boolean", B51AMirrorFieldEnumeration: map[int64]string{}, B51AMirrorFieldBits: []B51AMirrorBit{b51AMirrorFieldExpectedBit(0, "out-of-area", "named")}, B51AMirrorFieldGrade: SchemaHypothesis}}
	b51AMirrorFieldField.B51AMirrorFieldAlternatives[0].B51AMirrorFieldBits[0].B51AMirrorFieldDefault = nil
	b51AMirrorFieldFields = append(b51AMirrorFieldFields, b51AMirrorFieldField)
	for _, b51AMirrorFieldOperation := range []uint8{2, 3} {
		b51AMirrorFieldDirection, b51AMirrorFieldPosition, b51AMirrorFieldLimit := "positive", 5, "end-delimited; maximum unset; returned length independent of requested length"
		if b51AMirrorFieldOperation == 3 {
			b51AMirrorFieldDirection = "request"
			b51AMirrorFieldPosition = 9
			b51AMirrorFieldLimit = "end-delimited; maximum unset; empty data does not prove acceptance"
		}
		b51AMirrorFieldField = b51AMirrorFieldExpectedField(b51AMirrorFieldOperation, b51AMirrorFieldDirection, "data", b51AMirrorFieldPosition, 1, 255)
		b51AMirrorFieldField.B51AMirrorFieldWidthBytes = nil
		b51AMirrorFieldField.B51AMirrorFieldFormalMin = nil
		b51AMirrorFieldField.B51AMirrorFieldFormalMax = nil
		b51AMirrorFieldField.B51AMirrorFieldSignedness = "Unknown"
		b51AMirrorFieldField.B51AMirrorFieldConversion = "opaque"
		b51AMirrorFieldField.B51AMirrorFieldLogicalMin = b51AMirrorFieldPtr(int64(0))
		b51AMirrorFieldField.B51AMirrorFieldLogicalLimit = b51AMirrorFieldPtr(b51AMirrorFieldLimit)
		b51AMirrorFieldFields = append(b51AMirrorFieldFields, b51AMirrorFieldField)
	}
	b51AMirrorFieldField = b51AMirrorFieldExpectedField(3, "request", "flags", 4, 1, 255)
	b51AMirrorFieldField.B51AMirrorFieldDefault = b51AMirrorFieldPtr(int64(0))
	b51AMirrorFieldField.B51AMirrorFieldBits = []B51AMirrorBit{b51AMirrorFieldExpectedBit(0, "triple-storage", "named"), b51AMirrorFieldExpectedBit(1, "direct-update", "named")}
	b51AMirrorFieldField.B51AMirrorFieldReserved = []int{2, 3, 4, 5, 6, 7}
	for b51AMirrorFieldIndex := 2; b51AMirrorFieldIndex < 8; b51AMirrorFieldIndex++ {
		b51AMirrorFieldField.B51AMirrorFieldBits = append(b51AMirrorFieldField.B51AMirrorFieldBits, b51AMirrorFieldExpectedBit(b51AMirrorFieldIndex, "", "reserved"))
	}
	b51AMirrorFieldFields = append(b51AMirrorFieldFields, b51AMirrorFieldField)
	for b51AMirrorFieldIndex, b51AMirrorFieldName := range []string{"mirror_address", "symbol_address"} {
		b51AMirrorFieldField = b51AMirrorFieldExpectedField(3, "request", b51AMirrorFieldName, 5+2*b51AMirrorFieldIndex, 2, 65535)
		b51AMirrorFieldField.B51AMirrorFieldDefault = b51AMirrorFieldPtr(int64(0))
		b51AMirrorFieldField.B51AMirrorFieldLogicalLimit = b51AMirrorFieldPtr("independent " + b51AMirrorFieldName + " limit Unknown")
		b51AMirrorFieldFields = append(b51AMirrorFieldFields, b51AMirrorFieldField)
	}
	b51AMirrorFieldField = b51AMirrorFieldExpectedField(3, "positive", "flags", 4, 1, 255)
	b51AMirrorFieldField.B51AMirrorFieldDefault = b51AMirrorFieldPtr(int64(0))
	b51AMirrorFieldField.B51AMirrorFieldReserved = []int{4, 5, 6, 7}
	for b51AMirrorFieldIndex, b51AMirrorFieldName := range []string{"not-written", "invalid-symbol", "invalid-parameter", "storage-error"} {
		b51AMirrorFieldField.B51AMirrorFieldBits = append(b51AMirrorFieldField.B51AMirrorFieldBits, b51AMirrorFieldExpectedBit(b51AMirrorFieldIndex, b51AMirrorFieldName, "named"))
	}
	for b51AMirrorFieldIndex := 4; b51AMirrorFieldIndex < 8; b51AMirrorFieldIndex++ {
		b51AMirrorFieldField.B51AMirrorFieldBits = append(b51AMirrorFieldField.B51AMirrorFieldBits, b51AMirrorFieldExpectedBit(b51AMirrorFieldIndex, "", "reserved"))
	}
	b51AMirrorFieldFields = append(b51AMirrorFieldFields, b51AMirrorFieldField, b51AMirrorFieldExpectedField(6, "request", "level", 4, 1, 255))
	b51AMirrorFieldField = b51AMirrorFieldExpectedField(6, "request", "point", 5, 1, 255)
	b51AMirrorFieldField.B51AMirrorFieldLogicalMin = b51AMirrorFieldPtr(int64(0))
	b51AMirrorFieldField.B51AMirrorFieldLogicalMax = b51AMirrorFieldPtr(int64(99))
	b51AMirrorFieldFields = append(b51AMirrorFieldFields, b51AMirrorFieldField, b51AMirrorFieldExpectedField(6, "request", "value", 6, 2, 65535))
	b51AMirrorFieldField = b51AMirrorFieldExpectedField(6, "positive", "save_status", 4, 1, 255)
	b51AMirrorFieldField.B51AMirrorFieldDefault = b51AMirrorFieldPtr(int64(0))
	b51AMirrorFieldField.B51AMirrorFieldEnumeration = map[int64]string{0: "value-not-saved", 1: "point-not-found", 2: "write-protected", 3: "range-violation", 4: "misconfiguration"}
	b51AMirrorFieldField.B51AMirrorFieldUnmapped = "5..255 unlabeled; no implicit success"
	b51AMirrorFieldFields = append(b51AMirrorFieldFields, b51AMirrorFieldField)
	for b51AMirrorFieldIndex, b51AMirrorFieldName := range []string{"mirror_address", "symbol_address"} {
		b51AMirrorFieldField = b51AMirrorFieldExpectedField(6, "positive", b51AMirrorFieldName, 5+2*b51AMirrorFieldIndex, 2, 65535)
		b51AMirrorFieldField.B51AMirrorFieldLogicalLimit = b51AMirrorFieldPtr("independent " + b51AMirrorFieldName + " limit Unknown")
		b51AMirrorFieldFields = append(b51AMirrorFieldFields, b51AMirrorFieldField)
	}
	b51AMirrorFieldField = b51AMirrorFieldExpectedField(6, "positive", "value", 9, 2, 65535)
	b51AMirrorFieldField.B51AMirrorFieldWidthBytes = nil
	for b51AMirrorFieldIndex, b51AMirrorFieldName := range []string{"formal-u16", "display-one-byte", "display-two-byte"} {
		b51AMirrorFieldField.B51AMirrorFieldAlternatives = append(b51AMirrorFieldField.B51AMirrorFieldAlternatives, B51AMirrorAlternative{B51AMirrorFieldName: b51AMirrorFieldName, B51AMirrorFieldWidthBytes: b51AMirrorFieldPtr([]int{2, 1, 2}[b51AMirrorFieldIndex]), B51AMirrorFieldCondition: []string{"formal description", "display type selects one byte", "display type selects two bytes"}[b51AMirrorFieldIndex], B51AMirrorFieldEnumeration: map[int64]string{}, B51AMirrorFieldBits: []B51AMirrorBit{}, B51AMirrorFieldGrade: SchemaHypothesis})
	}
	return append(b51AMirrorFieldFields, b51AMirrorFieldField)
}

func b51AMirrorFieldCatalog(b51AMirrorFieldT *b51AMirrorFieldTesting.T) B51AMirrorMetadata {
	b51AMirrorFieldT.Helper()
	b51AMirrorFieldMetadata, b51AMirrorFieldError := B51AMirrorFields()
	if b51AMirrorFieldError != nil {
		b51AMirrorFieldT.Fatal(b51AMirrorFieldError)
	}
	return b51AMirrorFieldMetadata
}

func b51AMirrorFieldGet(b51AMirrorFieldT *b51AMirrorFieldTesting.T, b51AMirrorFieldOperation uint8, b51AMirrorFieldDirection, b51AMirrorFieldName string) B51AMirrorField {
	b51AMirrorFieldT.Helper()
	b51AMirrorFieldKey := B51AMirrorKey{B51AMirrorFieldFamily: "vaillant", B51AMirrorFieldVariant: "mirror-diagnosis", B51AMirrorFieldPrimary: 181, B51AMirrorFieldSecondary: b51AMirrorFieldPtr(uint8(26)), B51AMirrorFieldOperation: b51AMirrorFieldPtr(b51AMirrorFieldOperation), B51AMirrorFieldDirection: b51AMirrorFieldDirection, B51AMirrorFieldField: b51AMirrorFieldName}
	b51AMirrorFieldField, b51AMirrorFieldFound, b51AMirrorFieldError := B51ALookupMirrorField(b51AMirrorFieldKey)
	if b51AMirrorFieldError != nil || !b51AMirrorFieldFound {
		b51AMirrorFieldT.Fatalf("ab90 missing %v: %v", b51AMirrorFieldKey, b51AMirrorFieldError)
	}
	return b51AMirrorFieldField
}

func TestAb90ExactProductionMetadata(b51AMirrorFieldT *b51AMirrorFieldTesting.T) {
	b51AMirrorFieldGot := b51AMirrorFieldCatalog(b51AMirrorFieldT)
	b51AMirrorFieldWant := B51AMirrorMetadata{B51AMirrorFieldRepresentationGrade: SchemaProven, B51AMirrorFieldSemanticsGrade: SchemaHypothesis, B51AMirrorFieldPhysicalAdaptation: SchemaUnknown, B51AMirrorFieldStorageDurability: SchemaUnknown, B51AMirrorFieldSuccessSemantics: SchemaUnknown, B51AMirrorFieldComplete: false, B51AMirrorFieldNativeContract: "protocols/vaillant/ebus-vaillant-B51A.md#candidate-mirror-and-diagnosis-layouts", B51AMirrorFieldLimitations: []string{
		"Finite supplemental mirror/diagnosis declarations; configured-register tables remain separate.",
		"Logical positions include service prefix; physical response adaptation is Unknown.",
		"Sequence equality for reset/01 is candidate correlation only; no retry or freshness policy.",
		"No invocation, registration, transport, persistence or execution authority.",
		"Unprovided properties and broader family coverage remain Unknown.",
	}, B51AMirrorFieldFields: b51AMirrorFieldExpectedFields()}
	if !b51AMirrorFieldReflect.DeepEqual(b51AMirrorFieldGot, b51AMirrorFieldWant) {
		b51AMirrorFieldT.Error("ab90 production metadata differs from independent contract")
	}
	if len(b51AMirrorFieldGot.B51AMirrorFieldFields) != 28 {
		b51AMirrorFieldT.Fatalf("ab90 field count %d", len(b51AMirrorFieldGot.B51AMirrorFieldFields))
	}
	for b51AMirrorFieldIndex, b51AMirrorFieldField := range b51AMirrorFieldGot.B51AMirrorFieldFields {
		b51AMirrorFieldT.Run(b51AMirrorFieldFmt.Sprintf("ab90_field_%d", b51AMirrorFieldIndex), func(b51AMirrorFieldT *b51AMirrorFieldTesting.T) {
			if !b51AMirrorFieldReflect.DeepEqual(b51AMirrorFieldField, b51AMirrorFieldWant.B51AMirrorFieldFields[b51AMirrorFieldIndex]) {
				b51AMirrorFieldT.Fatalf("ab90 field %s/%s op %d mismatch", b51AMirrorFieldField.B51AMirrorFieldKey.B51AMirrorFieldDirection, b51AMirrorFieldField.B51AMirrorFieldKey.B51AMirrorFieldField, *b51AMirrorFieldField.B51AMirrorFieldKey.B51AMirrorFieldOperation)
			}
			b51AMirrorFieldFound, b51AMirrorFieldOK, b51AMirrorFieldError := B51ALookupMirrorField(b51AMirrorFieldField.B51AMirrorFieldKey)
			if b51AMirrorFieldError != nil || !b51AMirrorFieldOK || !b51AMirrorFieldReflect.DeepEqual(b51AMirrorFieldField, b51AMirrorFieldFound) {
				b51AMirrorFieldT.Fatal("ab90 lookup mismatch")
			}
		})
	}
}

func b51AMirrorFieldRead(b51AMirrorFieldT *b51AMirrorFieldTesting.T, b51AMirrorFieldField B51AMirrorField, b51AMirrorFieldRaw []byte, b51AMirrorFieldAlternative string) uint64 {
	b51AMirrorFieldT.Helper()
	b51AMirrorFieldWidth := b51AMirrorFieldField.B51AMirrorFieldWidthBytes
	if b51AMirrorFieldAlternative != "" {
		for _, b51AMirrorFieldChoice := range b51AMirrorFieldField.B51AMirrorFieldAlternatives {
			if b51AMirrorFieldChoice.B51AMirrorFieldName == b51AMirrorFieldAlternative {
				b51AMirrorFieldWidth = b51AMirrorFieldChoice.B51AMirrorFieldWidthBytes
			}
		}
	}
	if b51AMirrorFieldWidth == nil || b51AMirrorFieldField.B51AMirrorFieldLogicalPosition == nil {
		b51AMirrorFieldT.Fatal("ab90 no unconditional scalar layout")
	}
	b51AMirrorFieldPosition := *b51AMirrorFieldField.B51AMirrorFieldLogicalPosition
	if b51AMirrorFieldPosition+*b51AMirrorFieldWidth > len(b51AMirrorFieldRaw) {
		b51AMirrorFieldT.Fatal("ab90 short synthetic carrier")
	}
	if *b51AMirrorFieldWidth == 1 {
		return uint64(b51AMirrorFieldRaw[b51AMirrorFieldPosition])
	}
	if *b51AMirrorFieldWidth == 2 && b51AMirrorFieldField.B51AMirrorFieldByteOrder == "little_endian" {
		return uint64(b51AMirrorFieldBinary.LittleEndian.Uint16(b51AMirrorFieldRaw[b51AMirrorFieldPosition : b51AMirrorFieldPosition+2]))
	}
	b51AMirrorFieldT.Fatal("ab90 unsupported synthetic layout")
	return 0
}

func TestAb90SyntheticSequenceAndReset(b51AMirrorFieldT *b51AMirrorFieldTesting.T) {
	for _, b51AMirrorFieldOperation := range []uint8{0, 1, 2, 3, 6} {
		for _, b51AMirrorFieldDirection := range []string{"request", "positive"} {
			b51AMirrorFieldSequence := b51AMirrorFieldGet(b51AMirrorFieldT, b51AMirrorFieldOperation, b51AMirrorFieldDirection, "sequence")
			for _, b51AMirrorFieldValue := range []byte{0, 255} {
				if b51AMirrorFieldRead(b51AMirrorFieldT, b51AMirrorFieldSequence, []byte{0, 0, 0, b51AMirrorFieldValue}, "") != uint64(b51AMirrorFieldValue) {
					b51AMirrorFieldT.Fatal("ab90 sequence lost")
				}
			}
			b51AMirrorFieldExpected := "Unknown"
			if b51AMirrorFieldOperation == 0 || b51AMirrorFieldOperation == 1 {
				b51AMirrorFieldExpected = "candidate_equal_sequence"
			}
			if b51AMirrorFieldSequence.B51AMirrorFieldCorrelation != b51AMirrorFieldExpected {
				b51AMirrorFieldT.Fatal("ab90 correlation leaked across operations")
			}
		}
		b51AMirrorFieldRequest := b51AMirrorFieldGet(b51AMirrorFieldT, b51AMirrorFieldOperation, "request", "sequence")
		b51AMirrorFieldPositive := b51AMirrorFieldGet(b51AMirrorFieldT, b51AMirrorFieldOperation, "positive", "sequence")
		b51AMirrorFieldMatched := b51AMirrorFieldRequest.B51AMirrorFieldCorrelation == "candidate_equal_sequence" && b51AMirrorFieldRead(b51AMirrorFieldT, b51AMirrorFieldRequest, []byte{0, 0, 0, 0}, "") == b51AMirrorFieldRead(b51AMirrorFieldT, b51AMirrorFieldPositive, []byte{0, 0, 0, 255}, "")
		if b51AMirrorFieldMatched {
			b51AMirrorFieldT.Fatal("ab90 mismatched sequence correlated")
		}
	}
	b51AMirrorFieldAck := b51AMirrorFieldGet(b51AMirrorFieldT, 0, "positive", "acknowledgement")
	b51AMirrorFieldRaw := b51AMirrorFieldRead(b51AMirrorFieldT, b51AMirrorFieldAck, []byte{0, 0, 0, 0, 2}, "")
	if b51AMirrorFieldRaw&(1<<b51AMirrorFieldAck.B51AMirrorFieldBits[0].B51AMirrorFieldIndex) != 0 || b51AMirrorFieldRaw&(1<<b51AMirrorFieldAck.B51AMirrorFieldBits[1].B51AMirrorFieldIndex) == 0 || b51AMirrorFieldAck.B51AMirrorFieldBits[1].B51AMirrorFieldGrade != SchemaUnknown || b51AMirrorFieldAck.B51AMirrorFieldBits[1].B51AMirrorFieldName != nil {
		b51AMirrorFieldT.Fatal("ab90 reset02 is unnamed, not denied")
	}
}

func TestAb90SyntheticReadAndWrite(b51AMirrorFieldT *b51AMirrorFieldTesting.T) {
	b51AMirrorFieldLength := b51AMirrorFieldGet(b51AMirrorFieldT, 2, "request", "length")
	b51AMirrorFieldData := b51AMirrorFieldGet(b51AMirrorFieldT, 2, "positive", "data")
	b51AMirrorFieldRequested := b51AMirrorFieldRead(b51AMirrorFieldT, b51AMirrorFieldLength, []byte{0, 0, 0, 0, 0, 3}, "")
	b51AMirrorFieldReturned := []byte{0, 0, 0, 0, 1, 0x12, 0x34}
	b51AMirrorFieldActual := len(b51AMirrorFieldReturned[*b51AMirrorFieldData.B51AMirrorFieldLogicalPosition:])
	if b51AMirrorFieldRequested != 3 || b51AMirrorFieldActual != 2 || b51AMirrorFieldData.B51AMirrorFieldWidthBytes != nil || b51AMirrorFieldData.B51AMirrorFieldFormalMax != nil || *b51AMirrorFieldData.B51AMirrorFieldLogicalMin != 0 {
		b51AMirrorFieldT.Fatal("ab90 requested/actual lengths conflated")
	}
	b51AMirrorFieldPositive := b51AMirrorFieldGet(b51AMirrorFieldT, 2, "positive", "positive")
	for b51AMirrorFieldRaw, b51AMirrorFieldWant := range map[byte]string{0: "false", 1: "true", 2: ""} {
		b51AMirrorFieldValue := b51AMirrorFieldRead(b51AMirrorFieldT, b51AMirrorFieldPositive, []byte{0, 0, 0, 0, b51AMirrorFieldRaw}, "")
		if b51AMirrorFieldPositive.B51AMirrorFieldEnumeration[int64(b51AMirrorFieldValue)] != b51AMirrorFieldWant {
			b51AMirrorFieldT.Fatal("ab90 whole-u8 alternatives collapsed")
		}
	}
	b51AMirrorFieldMirror := b51AMirrorFieldGet(b51AMirrorFieldT, 3, "request", "mirror_address")
	b51AMirrorFieldSymbol := b51AMirrorFieldGet(b51AMirrorFieldT, 3, "request", "symbol_address")
	b51AMirrorFieldWrite := []byte{0, 0, 0, 0, 0, 0x34, 0x12, 0x78, 0x56}
	if b51AMirrorFieldRead(b51AMirrorFieldT, b51AMirrorFieldMirror, b51AMirrorFieldWrite, "") != 0x1234 || b51AMirrorFieldRead(b51AMirrorFieldT, b51AMirrorFieldSymbol, b51AMirrorFieldWrite, "") != 0x5678 {
		b51AMirrorFieldT.Fatal("ab90 mirror/symbol swapped")
	}
	b51AMirrorFieldOpaque := b51AMirrorFieldGet(b51AMirrorFieldT, 3, "request", "data")
	if len(b51AMirrorFieldWrite[*b51AMirrorFieldOpaque.B51AMirrorFieldLogicalPosition:]) != 0 || b51AMirrorFieldOpaque.B51AMirrorFieldConversion != "opaque" || b51AMirrorFieldOpaque.B51AMirrorFieldWidthBytes != nil || *b51AMirrorFieldOpaque.B51AMirrorFieldLogicalMin != 0 || b51AMirrorFieldCatalog(b51AMirrorFieldT).B51AMirrorFieldSuccessSemantics != SchemaUnknown {
		b51AMirrorFieldT.Fatal("ab90 empty write inferred acceptance")
	}
}

func TestAb90SyntheticSave(b51AMirrorFieldT *b51AMirrorFieldTesting.T) {
	b51AMirrorFieldPoint := b51AMirrorFieldGet(b51AMirrorFieldT, 6, "request", "point")
	for b51AMirrorFieldRaw, b51AMirrorFieldWant := range map[byte]bool{99: true, 100: false} {
		b51AMirrorFieldValue := int64(b51AMirrorFieldRead(b51AMirrorFieldT, b51AMirrorFieldPoint, []byte{0, 0, 0, 7, 9, b51AMirrorFieldRaw}, ""))
		b51AMirrorFieldFits := b51AMirrorFieldValue >= *b51AMirrorFieldPoint.B51AMirrorFieldLogicalMin && b51AMirrorFieldValue <= *b51AMirrorFieldPoint.B51AMirrorFieldLogicalMax
		if b51AMirrorFieldFits != b51AMirrorFieldWant {
			b51AMirrorFieldT.Fatal("ab90 point domain changed")
		}
	}
	b51AMirrorFieldStatus := b51AMirrorFieldGet(b51AMirrorFieldT, 6, "positive", "save_status")
	for b51AMirrorFieldRaw, b51AMirrorFieldWant := range map[byte]string{0: "value-not-saved", 1: "point-not-found", 5: ""} {
		if b51AMirrorFieldStatus.B51AMirrorFieldEnumeration[int64(b51AMirrorFieldRead(b51AMirrorFieldT, b51AMirrorFieldStatus, []byte{0, 0, 0, 0, b51AMirrorFieldRaw}, ""))] != b51AMirrorFieldWant {
			b51AMirrorFieldT.Fatal("ab90 save status invented success")
		}
	}
	b51AMirrorFieldValue := b51AMirrorFieldGet(b51AMirrorFieldT, 6, "positive", "value")
	b51AMirrorFieldCarrier := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0x34, 0x12}
	if b51AMirrorFieldValue.B51AMirrorFieldWidthBytes != nil || b51AMirrorFieldRead(b51AMirrorFieldT, b51AMirrorFieldValue, b51AMirrorFieldCarrier, "display-one-byte") != 0x34 || b51AMirrorFieldRead(b51AMirrorFieldT, b51AMirrorFieldValue, b51AMirrorFieldCarrier, "display-two-byte") != 0x1234 || b51AMirrorFieldRead(b51AMirrorFieldT, b51AMirrorFieldValue, b51AMirrorFieldCarrier, "formal-u16") != 0x1234 {
		b51AMirrorFieldT.Fatal("ab90 conditional widths collapsed")
	}
	b51AMirrorFieldRequest := b51AMirrorFieldGet(b51AMirrorFieldT, 6, "request", "value")
	if b51AMirrorFieldRead(b51AMirrorFieldT, b51AMirrorFieldRequest, []byte{0, 0, 0, 0, 0, 0, 0x34, 0x12}, "") != 0x1234 {
		b51AMirrorFieldT.Fatal("ab90 save request narrowed")
	}
}

func TestAb90IdentityAxes(b51AMirrorFieldT *b51AMirrorFieldTesting.T) {
	b51AMirrorFieldBase := b51AMirrorFieldGet(b51AMirrorFieldT, 2, "request", "length").B51AMirrorFieldKey
	for b51AMirrorFieldIndex := 0; b51AMirrorFieldIndex < b51AMirrorFieldReflect.TypeOf(b51AMirrorFieldBase).NumField(); b51AMirrorFieldIndex++ {
		b51AMirrorFieldAxis := b51AMirrorFieldReflect.TypeOf(b51AMirrorFieldBase).Field(b51AMirrorFieldIndex).Name
		b51AMirrorFieldT.Run("ab90_"+b51AMirrorFieldAxis, func(b51AMirrorFieldT *b51AMirrorFieldTesting.T) {
			b51AMirrorFieldChanged := b51AMirrorFieldBase
			b51AMirrorFieldValue := b51AMirrorFieldReflect.ValueOf(&b51AMirrorFieldChanged).Elem().Field(b51AMirrorFieldIndex)
			switch b51AMirrorFieldValue.Kind() {
			case b51AMirrorFieldReflect.String:
				b51AMirrorFieldValue.SetString("ab90_distinct")
			case b51AMirrorFieldReflect.Uint8:
				b51AMirrorFieldValue.SetUint(0)
			case b51AMirrorFieldReflect.Pointer:
				b51AMirrorFieldNew := b51AMirrorFieldReflect.New(b51AMirrorFieldValue.Type().Elem())
				b51AMirrorFieldNew.Elem().SetUint(99)
				b51AMirrorFieldValue.Set(b51AMirrorFieldNew)
			default:
				b51AMirrorFieldT.Fatal("ab90 unchecked axis")
			}
			b51AMirrorFieldField, b51AMirrorFieldFound, b51AMirrorFieldError := B51ALookupMirrorField(b51AMirrorFieldChanged)
			if b51AMirrorFieldError != nil || b51AMirrorFieldFound || !b51AMirrorFieldReflect.DeepEqual(b51AMirrorFieldField, B51AMirrorField{}) {
				b51AMirrorFieldT.Fatal("ab90 changed identity matched")
			}
			if b51AMirrorFieldMirrorKeysEqual(b51AMirrorFieldBase, b51AMirrorFieldChanged) {
				b51AMirrorFieldT.Fatal("ab90 ignored key dimension")
			}
			b51AMirrorFieldValue.Set(b51AMirrorFieldReflect.ValueOf(b51AMirrorFieldBase).Field(b51AMirrorFieldIndex))
			if _, b51AMirrorFieldFound, b51AMirrorFieldError = B51ALookupMirrorField(b51AMirrorFieldChanged); b51AMirrorFieldError != nil || !b51AMirrorFieldFound {
				b51AMirrorFieldT.Fatal("ab90 invalid restored control")
			}
		})
	}
	for _, b51AMirrorFieldField := range b51AMirrorFieldCatalog(b51AMirrorFieldT).B51AMirrorFieldFields {
		for _, b51AMirrorFieldAxis := range []string{"B51AMirrorFieldSecondary", "B51AMirrorFieldOperation", "B51AMirrorFieldOpcode", "B51AMirrorFieldGroup", "B51AMirrorFieldInstance", "B51AMirrorFieldRegister"} {
			b51AMirrorFieldChanged := b51AMirrorFieldField.B51AMirrorFieldKey
			b51AMirrorFieldValue := b51AMirrorFieldReflect.ValueOf(&b51AMirrorFieldChanged).Elem().FieldByName(b51AMirrorFieldAxis)
			if b51AMirrorFieldValue.IsNil() {
				b51AMirrorFieldValue.Set(b51AMirrorFieldReflect.New(b51AMirrorFieldValue.Type().Elem()))
			} else {
				b51AMirrorFieldValue.SetZero()
			}
			if _, b51AMirrorFieldFound, b51AMirrorFieldError := B51ALookupMirrorField(b51AMirrorFieldChanged); b51AMirrorFieldError != nil || b51AMirrorFieldFound {
				b51AMirrorFieldT.Fatalf("ab90 absent versus supplied %s wildcarded", b51AMirrorFieldAxis)
			}
		}
	}
}

func TestAb90B524IdentityAxes(b51AMirrorFieldT *b51AMirrorFieldTesting.T) {
	b51AMirrorFieldBase := B51AMirrorKey{B51AMirrorFieldFamily: "ab90_synthetic_family", B51AMirrorFieldVariant: "ab90_synthetic_variant", B51AMirrorFieldPrimary: 181, B51AMirrorFieldSecondary: b51AMirrorFieldPtr(uint8(36)), B51AMirrorFieldOperation: b51AMirrorFieldPtr(uint8(0)), B51AMirrorFieldDirection: "request", B51AMirrorFieldField: "ab90_synthetic_field", B51AMirrorFieldOpcode: b51AMirrorFieldPtr(uint8(0)), B51AMirrorFieldGroup: b51AMirrorFieldPtr(uint8(0)), B51AMirrorFieldInstance: b51AMirrorFieldPtr(uint8(0)), B51AMirrorFieldRegister: b51AMirrorFieldPtr(uint16(0))}
	for b51AMirrorFieldIndex := 0; b51AMirrorFieldIndex < b51AMirrorFieldReflect.TypeOf(b51AMirrorFieldBase).NumField(); b51AMirrorFieldIndex++ {
		b51AMirrorFieldAxis := b51AMirrorFieldReflect.TypeOf(b51AMirrorFieldBase).Field(b51AMirrorFieldIndex).Name
		b51AMirrorFieldT.Run("ab90_"+b51AMirrorFieldAxis, func(b51AMirrorFieldT *b51AMirrorFieldTesting.T) {
			b51AMirrorFieldChanged := b51AMirrorFieldBase
			b51AMirrorFieldValue := b51AMirrorFieldReflect.ValueOf(&b51AMirrorFieldChanged).Elem().Field(b51AMirrorFieldIndex)
			switch b51AMirrorFieldValue.Kind() {
			case b51AMirrorFieldReflect.String:
				b51AMirrorFieldValue.SetString("ab90_other")
			case b51AMirrorFieldReflect.Uint8:
				b51AMirrorFieldValue.SetUint(0)
			case b51AMirrorFieldReflect.Pointer:
				b51AMirrorFieldNew := b51AMirrorFieldReflect.New(b51AMirrorFieldValue.Type().Elem())
				b51AMirrorFieldNew.Elem().SetUint(b51AMirrorFieldValue.Elem().Uint() + 1)
				b51AMirrorFieldValue.Set(b51AMirrorFieldNew)
			}
			if b51AMirrorFieldMirrorKeysEqual(b51AMirrorFieldBase, b51AMirrorFieldChanged) {
				b51AMirrorFieldT.Fatal("ab90 B524 axis ignored")
			}
			if b51AMirrorFieldValue.Kind() == b51AMirrorFieldReflect.Pointer {
				b51AMirrorFieldValue.SetZero()
				if b51AMirrorFieldMirrorKeysEqual(b51AMirrorFieldBase, b51AMirrorFieldChanged) {
					b51AMirrorFieldT.Fatal("ab90 B524 absence wildcarded")
				}
				b51AMirrorFieldValue.Set(b51AMirrorFieldReflect.New(b51AMirrorFieldValue.Type().Elem()))
				b51AMirrorFieldZero := b51AMirrorFieldChanged
				b51AMirrorFieldValue.SetZero()
				if b51AMirrorFieldMirrorKeysEqual(b51AMirrorFieldZero, b51AMirrorFieldChanged) {
					b51AMirrorFieldT.Fatal("ab90 B524 absent equals zero")
				}
			}
			b51AMirrorFieldValue.Set(b51AMirrorFieldReflect.ValueOf(b51AMirrorFieldBase).Field(b51AMirrorFieldIndex))
			if !b51AMirrorFieldMirrorKeysEqual(b51AMirrorFieldBase, b51AMirrorFieldChanged) {
				b51AMirrorFieldT.Fatal("ab90 B524 restored control differs")
			}
		})
	}
	if _, b51AMirrorFieldFound, b51AMirrorFieldError := B51ALookupMirrorField(b51AMirrorFieldBase); b51AMirrorFieldError != nil || b51AMirrorFieldFound {
		b51AMirrorFieldT.Fatal("ab90 synthetic B524 acquired production declaration")
	}
}

func b51AMirrorFieldMutate(b51AMirrorFieldValue b51AMirrorFieldReflect.Value) {
	switch b51AMirrorFieldValue.Kind() {
	case b51AMirrorFieldReflect.Struct:
		for b51AMirrorFieldIndex := 0; b51AMirrorFieldIndex < b51AMirrorFieldValue.NumField(); b51AMirrorFieldIndex++ {
			b51AMirrorFieldMutate(b51AMirrorFieldValue.Field(b51AMirrorFieldIndex))
		}
	case b51AMirrorFieldReflect.Pointer:
		if !b51AMirrorFieldValue.IsNil() {
			b51AMirrorFieldMutate(b51AMirrorFieldValue.Elem())
		}
	case b51AMirrorFieldReflect.Slice:
		for b51AMirrorFieldIndex := 0; b51AMirrorFieldIndex < b51AMirrorFieldValue.Len(); b51AMirrorFieldIndex++ {
			b51AMirrorFieldMutate(b51AMirrorFieldValue.Index(b51AMirrorFieldIndex))
		}
	case b51AMirrorFieldReflect.Map:
		for _, b51AMirrorFieldKey := range b51AMirrorFieldValue.MapKeys() {
			b51AMirrorFieldValue.SetMapIndex(b51AMirrorFieldKey, b51AMirrorFieldReflect.ValueOf("ab90_mutated"))
		}
	case b51AMirrorFieldReflect.String:
		b51AMirrorFieldValue.SetString("ab90_mutated")
	case b51AMirrorFieldReflect.Bool:
		b51AMirrorFieldValue.SetBool(!b51AMirrorFieldValue.Bool())
	case b51AMirrorFieldReflect.Int, b51AMirrorFieldReflect.Int64:
		b51AMirrorFieldValue.SetInt(b51AMirrorFieldValue.Int() + 1)
	case b51AMirrorFieldReflect.Uint8, b51AMirrorFieldReflect.Uint16:
		b51AMirrorFieldValue.SetUint(b51AMirrorFieldValue.Uint() + 1)
	}
}

func TestAb90DeepCopyIsolation(b51AMirrorFieldT *b51AMirrorFieldTesting.T) {
	b51AMirrorFieldBefore := b51AMirrorFieldCatalog(b51AMirrorFieldT)
	b51AMirrorFieldCopy := b51AMirrorFieldCatalog(b51AMirrorFieldT)
	b51AMirrorFieldMutate(b51AMirrorFieldReflect.ValueOf(&b51AMirrorFieldCopy).Elem())
	if b51AMirrorFieldReflect.DeepEqual(b51AMirrorFieldBefore, b51AMirrorFieldCopy) {
		b51AMirrorFieldT.Fatal("ab90 mutation control did not change")
	}
	if !b51AMirrorFieldReflect.DeepEqual(b51AMirrorFieldBefore, b51AMirrorFieldCatalog(b51AMirrorFieldT)) {
		b51AMirrorFieldT.Fatal("ab90 accessor shares nested storage")
	}
	for _, b51AMirrorFieldExpected := range b51AMirrorFieldBefore.B51AMirrorFieldFields {
		b51AMirrorFieldResult, b51AMirrorFieldFound, b51AMirrorFieldError := B51ALookupMirrorField(b51AMirrorFieldExpected.B51AMirrorFieldKey)
		if b51AMirrorFieldError != nil || !b51AMirrorFieldFound {
			b51AMirrorFieldT.Fatal("ab90 lookup control missing")
		}
		b51AMirrorFieldMutate(b51AMirrorFieldReflect.ValueOf(&b51AMirrorFieldResult).Elem())
		b51AMirrorFieldFresh, b51AMirrorFieldFound, b51AMirrorFieldError := B51ALookupMirrorField(b51AMirrorFieldExpected.B51AMirrorFieldKey)
		if b51AMirrorFieldError != nil || !b51AMirrorFieldFound || !b51AMirrorFieldReflect.DeepEqual(b51AMirrorFieldFresh, b51AMirrorFieldExpected) {
			b51AMirrorFieldT.Fatal("ab90 lookup shares nested storage")
		}
	}
}
