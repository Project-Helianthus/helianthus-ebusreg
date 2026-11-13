package vaillant

import (
	b523SensorFieldBinary "encoding/binary"
	b523SensorFieldJSON "encoding/json"
	b523SensorFieldFmt "fmt"
	b523SensorFieldReflect "reflect"
	b523SensorFieldTesting "testing"
)

func b523SensorFieldPtr[B523SensorFieldValue any](b523SensorFieldValue B523SensorFieldValue) *B523SensorFieldValue {
	return &b523SensorFieldValue
}

func b523SensorFieldCatalog(b523SensorFieldT *b523SensorFieldTesting.T) B523SensorCatalog {
	b523SensorFieldT.Helper()
	b523SensorFieldCatalog, b523SensorFieldError := B523SensorFields()
	if b523SensorFieldError != nil {
		b523SensorFieldT.Fatal(b523SensorFieldError)
	}
	return b523SensorFieldCatalog
}

func b523SensorFieldExpectedBase(b523SensorFieldVariant string, b523SensorFieldOperation uint16, b523SensorFieldName string, b523SensorFieldPosition, b523SensorFieldWidth int) B523SensorField {
	return B523SensorField{
		B523SensorFieldKey:      B523SensorKey{B523SensorFieldFamily: "Vaillant", B523SensorFieldVariant: b523SensorFieldVariant, B523SensorFieldPrimary: b523SensorFieldPtr(uint16(181)), B523SensorFieldSecondary: b523SensorFieldPtr(uint16(35)), B523SensorFieldOperation: b523SensorFieldPtr(b523SensorFieldOperation), B523SensorFieldField: b523SensorFieldName},
		B523SensorFieldPosition: b523SensorFieldPtr(b523SensorFieldPosition), B523SensorFieldWidth: b523SensorFieldPtr(b523SensorFieldWidth), B523SensorFieldSignedness: "Unknown", B523SensorFieldByteOrder: "Unknown", B523SensorFieldEnumeration: map[string]string{}, B523SensorFieldSelection: "none", B523SensorFieldBranches: []B523SensorBranch{}, B523SensorFieldGrade: "Hypothesis", B523SensorFieldRawFallback: "preserve original bytes including reserved bits",
	}
}

func b523SensorFieldExpectedFields() []B523SensorField {
	var b523SensorFieldFields []B523SensorField
	for b523SensorFieldVariant, b523SensorFieldZero := range map[string]string{"zero-valid-value": "valid-value", "zero-valid-value-or-low-level": "valid-value-or-low-level"} {
		for _, b523SensorFieldOperation := range []struct {
			B523SensorFieldOp                                                  uint16
			B523Sensors, B523SensorFieldPositions                              []int
			B523SensorFieldFirst, B523SensorFieldLast, B523SensorFieldReserved int
		}{
			{3, []int{1, 2, 3, 4, 5, 6}, []int{3, 5, 7, 9, 11, 13}, 16, 17, 4},
			{6, []int{1, 2, 3, 4, 5, 6, 7}, []int{3, 5, 7, 9, 11, 13, 15}, 17, 18, 6},
			{7, []int{8, 9, 10, 11, 12, 13}, []int{3, 5, 7, 9, 12, 14}, 16, 17, 4},
		} {
			for b523SensorFieldIndex, b523SensorFieldSensor := range b523SensorFieldOperation.B523Sensors {
				b523SensorFieldField := b523SensorFieldExpectedBase(b523SensorFieldVariant, b523SensorFieldOperation.B523SensorFieldOp, b523SensorFieldFmt.Sprintf("S%d", b523SensorFieldSensor), b523SensorFieldOperation.B523SensorFieldPositions[b523SensorFieldIndex], 16)
				b523SensorFieldField.B523SensorFieldReplacementRaw = b523SensorFieldPtr(uint16(32768))
				b523SensorFieldField.B523SensorFieldSelection = "unsigned raw equality before signed conversion"
				b523SensorFieldField.B523SensorFieldBranches = []B523SensorBranch{
					{B523SensorFieldName: "replacement", B523SensorFieldSignedness: "unsigned", B523SensorFieldByteOrder: "little-endian", B523SensorFieldDivideBy: b523SensorFieldPtr(1.0), B523SensorFieldDefault: b523SensorFieldPtr(32768.0), B523SensorFieldBounds: []float64{0, 65535}, B523SensorFieldGrade: "Hypothesis"},
					{B523SensorFieldName: "normal", B523SensorFieldSignedness: "signed", B523SensorFieldByteOrder: "little-endian", B523SensorFieldDivideBy: b523SensorFieldPtr(16.0), B523SensorFieldUnit: b523SensorFieldPtr("Celsius"), B523SensorFieldDefault: b523SensorFieldPtr(0.0), B523SensorFieldBounds: []float64{-32768, 32767}, B523SensorFieldGrade: "Hypothesis"},
				}
				b523SensorFieldFields = append(b523SensorFieldFields, b523SensorFieldField)
				b523SensorFieldPosition := b523SensorFieldOperation.B523SensorFieldFirst
				if b523SensorFieldIndex >= 4 {
					b523SensorFieldPosition = b523SensorFieldOperation.B523SensorFieldLast
				}
				b523SensorFieldField = b523SensorFieldExpectedBase(b523SensorFieldVariant, b523SensorFieldOperation.B523SensorFieldOp, b523SensorFieldFmt.Sprintf("S%d.status", b523SensorFieldSensor), b523SensorFieldPosition, 2)
				b523SensorFieldField.B523SensorFieldBitOffset = b523SensorFieldPtr(b523SensorFieldIndex % 4 * 2)
				b523SensorFieldField.B523SensorFieldSignedness = "unsigned"
				b523SensorFieldField.B523SensorFieldByteOrder = "not_applicable"
				b523SensorFieldField.B523SensorFieldDivideBy = b523SensorFieldPtr(1.0)
				b523SensorFieldField.B523SensorFieldDefault = b523SensorFieldPtr(0.0)
				b523SensorFieldField.B523SensorFieldBounds = []float64{0, 3}
				b523SensorFieldField.B523SensorFieldEnumeration = map[string]string{"0": b523SensorFieldZero, "1": "open-circuit", "2": "closed-circuit", "3": "high-level"}
				b523SensorFieldFields = append(b523SensorFieldFields, b523SensorFieldField)
			}
			b523SensorFieldField := b523SensorFieldExpectedBase(b523SensorFieldVariant, b523SensorFieldOperation.B523SensorFieldOp, "status.reserved", b523SensorFieldOperation.B523SensorFieldLast, 8-b523SensorFieldOperation.B523SensorFieldReserved)
			b523SensorFieldField.B523SensorFieldBitOffset = b523SensorFieldPtr(b523SensorFieldOperation.B523SensorFieldReserved)
			b523SensorFieldField.B523SensorFieldReserved = true
			b523SensorFieldField.B523SensorFieldSignedness = "unsigned"
			b523SensorFieldField.B523SensorFieldByteOrder = "not_applicable"
			b523SensorFieldFields = append(b523SensorFieldFields, b523SensorFieldField)
			if b523SensorFieldOperation.B523SensorFieldOp == 3 {
				b523SensorFieldField = b523SensorFieldExpectedBase(b523SensorFieldVariant, 3, "PWM", 15, 8)
				b523SensorFieldField.B523SensorFieldSignedness = "unsigned"
				b523SensorFieldField.B523SensorFieldByteOrder = "not_applicable"
				b523SensorFieldField.B523SensorFieldDivideBy = b523SensorFieldPtr(2.0)
				b523SensorFieldField.B523SensorFieldUnit = b523SensorFieldPtr("percent")
				b523SensorFieldField.B523SensorFieldBounds = []float64{0, 255}
				b523SensorFieldFields = append(b523SensorFieldFields, b523SensorFieldField)
			}
			if b523SensorFieldOperation.B523SensorFieldOp == 7 {
				b523SensorFieldField = b523SensorFieldExpectedBase(b523SensorFieldVariant, 7, "opaque", 11, 8)
				b523SensorFieldField.B523SensorFieldSignedness = "unsigned"
				b523SensorFieldField.B523SensorFieldByteOrder = "not_applicable"
				b523SensorFieldField.B523SensorFieldBounds = []float64{0, 255}
				b523SensorFieldFields = append(b523SensorFieldFields, b523SensorFieldField)
			}
		}
	}
	return b523SensorFieldFields
}

func TestAb93ProductionMetadata(b523SensorFieldT *b523SensorFieldTesting.T) {
	b523SensorFieldCatalog := b523SensorFieldCatalog(b523SensorFieldT)
	if b523SensorFieldCatalog.B523SensorFieldRepresentationGrade != "Proven" || b523SensorFieldCatalog.B523SensorFieldQualification != "Unknown" || b523SensorFieldCatalog.B523SensorFieldComplete || len(b523SensorFieldCatalog.B523SensorFieldLimitations) != 6 {
		b523SensorFieldT.Fatal("qualification/completeness boundary")
	}
	b523SensorFieldExpected := b523SensorFieldExpectedFields()
	if len(b523SensorFieldCatalog.B523SensorFieldFields) != 86 || len(b523SensorFieldExpected) != 86 {
		b523SensorFieldT.Fatal("finite field count")
	}
	b523SensorFieldSeen := map[string]bool{}
	for _, b523SensorFieldWant := range b523SensorFieldExpected {
		b523SensorFieldGot, b523SensorFieldFound, b523SensorFieldError := B523LookupSensorField(b523SensorFieldWant.B523SensorFieldKey)
		if b523SensorFieldError != nil || !b523SensorFieldFound || !b523SensorFieldReflect.DeepEqual(b523SensorFieldGot, b523SensorFieldWant) {
			b523SensorFieldT.Fatalf("metadata mismatch %v: got %+v want %+v, %v", b523SensorFieldWant.B523SensorFieldKey, b523SensorFieldGot, b523SensorFieldWant, b523SensorFieldError)
		}
		b523SensorFieldBytes, _ := b523SensorFieldJSON.Marshal(b523SensorFieldWant.B523SensorFieldKey)
		if b523SensorFieldSeen[string(b523SensorFieldBytes)] {
			b523SensorFieldT.Fatal("duplicate key")
		}
		b523SensorFieldSeen[string(b523SensorFieldBytes)] = true
	}
	for _, b523SensorFieldField := range b523SensorFieldCatalog.B523SensorFieldFields {
		b523SensorFieldBytes, _ := b523SensorFieldJSON.Marshal(b523SensorFieldField.B523SensorFieldKey)
		if !b523SensorFieldSeen[string(b523SensorFieldBytes)] {
			b523SensorFieldT.Fatal("unexpected key")
		}
	}
}

func b523SensorFieldEvaluate(b523SensorFieldT *b523SensorFieldTesting.T, b523SensorFieldField B523SensorField, b523SensorFieldBytes []byte) (string, float64, string, []byte) {
	b523SensorFieldT.Helper()
	b523SensorFieldRaw := append([]byte(nil), b523SensorFieldBytes...)
	if b523SensorFieldField.B523SensorFieldRawFallback != "preserve original bytes including reserved bits" {
		b523SensorFieldT.Fatal("raw preservation")
	}
	if b523SensorFieldField.B523SensorFieldReplacementRaw != nil {
		b523SensorFieldWord := b523SensorFieldBinary.LittleEndian.Uint16(b523SensorFieldBytes[*b523SensorFieldField.B523SensorFieldPosition:])
		b523SensorFieldBranch := b523SensorFieldField.B523SensorFieldBranches[1]
		if b523SensorFieldWord == *b523SensorFieldField.B523SensorFieldReplacementRaw {
			b523SensorFieldBranch = b523SensorFieldField.B523SensorFieldBranches[0]
		}
		b523SensorFieldValue := float64(b523SensorFieldWord)
		if b523SensorFieldBranch.B523SensorFieldSignedness == "signed" {
			b523SensorFieldValue = float64(int16(b523SensorFieldWord))
		}
		return b523SensorFieldBranch.B523SensorFieldName, b523SensorFieldValue / *b523SensorFieldBranch.B523SensorFieldDivideBy, "", b523SensorFieldRaw
	}
	b523SensorFieldValue := uint16(b523SensorFieldBytes[*b523SensorFieldField.B523SensorFieldPosition])
	if b523SensorFieldField.B523SensorFieldBitOffset != nil {
		b523SensorFieldValue = (b523SensorFieldValue >> *b523SensorFieldField.B523SensorFieldBitOffset) & ((1 << *b523SensorFieldField.B523SensorFieldWidth) - 1)
	}
	if b523SensorFieldField.B523SensorFieldReserved {
		return "reserved", float64(b523SensorFieldValue), "", b523SensorFieldRaw
	}
	b523SensorFieldNumber := float64(b523SensorFieldValue)
	if b523SensorFieldField.B523SensorFieldDivideBy != nil {
		b523SensorFieldNumber /= *b523SensorFieldField.B523SensorFieldDivideBy
	}
	return "value", b523SensorFieldNumber, b523SensorFieldField.B523SensorFieldEnumeration[b523SensorFieldFmt.Sprint(b523SensorFieldValue)], b523SensorFieldRaw
}

func TestAb93SyntheticSensorWords(b523SensorFieldT *b523SensorFieldTesting.T) {
	b523SensorFieldFields := b523SensorFieldCatalog(b523SensorFieldT).B523SensorFieldFields
	b523SensorFieldCount := 0
	for _, b523SensorFieldTarget := range b523SensorFieldFields {
		if b523SensorFieldTarget.B523SensorFieldReplacementRaw == nil {
			continue
		}
		b523SensorFieldCount++
		for _, b523SensorFieldCase := range []struct {
			B523SensorFieldBytes  []byte
			B523SensorFieldBranch string
			B523SensorFieldValue  float64
		}{{[]byte{0, 128}, "replacement", 32768}, {[]byte{16, 0}, "normal", 1}, {[]byte{240, 255}, "normal", -1}} {
			b523SensorFieldT.Run(b523SensorFieldFmt.Sprintf("ab93_%s_%d_%s_%s", b523SensorFieldTarget.B523SensorFieldKey.B523SensorFieldVariant, *b523SensorFieldTarget.B523SensorFieldKey.B523SensorFieldOperation, b523SensorFieldTarget.B523SensorFieldKey.B523SensorFieldField, b523SensorFieldCase.B523SensorFieldBranch+b523SensorFieldFmt.Sprint(b523SensorFieldCase.B523SensorFieldValue)), func(b523SensorFieldT *b523SensorFieldTesting.T) {
				b523SensorFieldBytes := make([]byte, 20)
				copy(b523SensorFieldBytes[*b523SensorFieldTarget.B523SensorFieldPosition:], b523SensorFieldCase.B523SensorFieldBytes)
				for _, b523SensorFieldField := range b523SensorFieldFields {
					if b523SensorFieldField.B523SensorFieldReplacementRaw == nil || b523SensorFieldField.B523SensorFieldKey.B523SensorFieldVariant != b523SensorFieldTarget.B523SensorFieldKey.B523SensorFieldVariant || *b523SensorFieldField.B523SensorFieldKey.B523SensorFieldOperation != *b523SensorFieldTarget.B523SensorFieldKey.B523SensorFieldOperation {
						continue
					}
					b523SensorFieldBranch, b523SensorFieldValue, _, b523SensorFieldRaw := b523SensorFieldEvaluate(b523SensorFieldT, b523SensorFieldField, b523SensorFieldBytes)
					b523SensorFieldWantBranch, b523SensorFieldWantValue := "normal", 0.0
					if b523SensorFieldField.B523SensorFieldKey.B523SensorFieldField == b523SensorFieldTarget.B523SensorFieldKey.B523SensorFieldField {
						b523SensorFieldWantBranch, b523SensorFieldWantValue = b523SensorFieldCase.B523SensorFieldBranch, b523SensorFieldCase.B523SensorFieldValue
					}
					if b523SensorFieldBranch != b523SensorFieldWantBranch || b523SensorFieldValue != b523SensorFieldWantValue || !b523SensorFieldReflect.DeepEqual(b523SensorFieldRaw, b523SensorFieldBytes) {
						b523SensorFieldT.Fatalf("word isolation: %s %s %v", b523SensorFieldField.B523SensorFieldKey.B523SensorFieldField, b523SensorFieldBranch, b523SensorFieldValue)
					}
				}
			})
		}
	}
	if b523SensorFieldCount != 38 {
		b523SensorFieldT.Fatalf("executed %d sensor occurrences, want 19 per variant", b523SensorFieldCount)
	}
}

func TestAb93SyntheticStatusAndReserved(b523SensorFieldT *b523SensorFieldTesting.T) {
	b523SensorFieldFields := b523SensorFieldCatalog(b523SensorFieldT).B523SensorFieldFields
	b523SensorFieldCount, b523SensorFieldReservedCount := 0, 0
	for _, b523SensorFieldTarget := range b523SensorFieldFields {
		if len(b523SensorFieldTarget.B523SensorFieldEnumeration) == 0 && !b523SensorFieldTarget.B523SensorFieldReserved {
			continue
		}
		if b523SensorFieldTarget.B523SensorFieldReserved {
			for b523SensorFieldBit := *b523SensorFieldTarget.B523SensorFieldBitOffset; b523SensorFieldBit < 8; b523SensorFieldBit++ {
				b523SensorFieldReservedCount++
				b523SensorFieldBytes := make([]byte, 20)
				b523SensorFieldBytes[*b523SensorFieldTarget.B523SensorFieldPosition] = 1 << b523SensorFieldBit
				b523SensorFieldBranch, b523SensorFieldValue, b523SensorFieldLabel, b523SensorFieldRaw := b523SensorFieldEvaluate(b523SensorFieldT, b523SensorFieldTarget, b523SensorFieldBytes)
				if b523SensorFieldBranch != "reserved" || b523SensorFieldValue != float64(uint16(1)<<(b523SensorFieldBit-*b523SensorFieldTarget.B523SensorFieldBitOffset)) || b523SensorFieldLabel != "" || !b523SensorFieldReflect.DeepEqual(b523SensorFieldRaw, b523SensorFieldBytes) {
					b523SensorFieldT.Fatal("reserved bit not preserved")
				}
				for _, b523SensorFieldField := range b523SensorFieldFields {
					if len(b523SensorFieldField.B523SensorFieldEnumeration) == 0 || b523SensorFieldField.B523SensorFieldKey.B523SensorFieldVariant != b523SensorFieldTarget.B523SensorFieldKey.B523SensorFieldVariant || *b523SensorFieldField.B523SensorFieldKey.B523SensorFieldOperation != *b523SensorFieldTarget.B523SensorFieldKey.B523SensorFieldOperation {
						continue
					}
					_, b523SensorFieldValue, _, b523SensorFieldRaw := b523SensorFieldEvaluate(b523SensorFieldT, b523SensorFieldField, b523SensorFieldBytes)
					if b523SensorFieldValue != 0 || !b523SensorFieldReflect.DeepEqual(b523SensorFieldRaw, b523SensorFieldBytes) {
						b523SensorFieldT.Fatal("reserved bit became sensor state")
					}
				}
			}
			continue
		}
		b523SensorFieldCount++
		for b523SensorFieldState := 0; b523SensorFieldState < 4; b523SensorFieldState++ {
			b523SensorFieldBytes := make([]byte, 20)
			b523SensorFieldBytes[*b523SensorFieldTarget.B523SensorFieldPosition] = byte(b523SensorFieldState << *b523SensorFieldTarget.B523SensorFieldBitOffset)
			for _, b523SensorFieldField := range b523SensorFieldFields {
				if len(b523SensorFieldField.B523SensorFieldEnumeration) == 0 || b523SensorFieldField.B523SensorFieldKey.B523SensorFieldVariant != b523SensorFieldTarget.B523SensorFieldKey.B523SensorFieldVariant || *b523SensorFieldField.B523SensorFieldKey.B523SensorFieldOperation != *b523SensorFieldTarget.B523SensorFieldKey.B523SensorFieldOperation {
					continue
				}
				b523SensorFieldWant := 0
				if b523SensorFieldField.B523SensorFieldKey.B523SensorFieldField == b523SensorFieldTarget.B523SensorFieldKey.B523SensorFieldField {
					b523SensorFieldWant = b523SensorFieldState
				}
				b523SensorFieldLabels := []string{"valid-value", "open-circuit", "closed-circuit", "high-level"}
				if b523SensorFieldField.B523SensorFieldKey.B523SensorFieldVariant == "zero-valid-value-or-low-level" {
					b523SensorFieldLabels[0] = "valid-value-or-low-level"
				}
				_, b523SensorFieldValue, b523SensorFieldLabel, b523SensorFieldRaw := b523SensorFieldEvaluate(b523SensorFieldT, b523SensorFieldField, b523SensorFieldBytes)
				if b523SensorFieldValue != float64(b523SensorFieldWant) || b523SensorFieldLabel != b523SensorFieldLabels[b523SensorFieldWant] || !b523SensorFieldReflect.DeepEqual(b523SensorFieldRaw, b523SensorFieldBytes) {
					b523SensorFieldT.Fatal("status isolation/alternative label")
				}
			}
		}
	}
	if b523SensorFieldCount != 38 || b523SensorFieldReservedCount != 20 {
		b523SensorFieldT.Fatalf("status/reserved coverage %d/%d", b523SensorFieldCount, b523SensorFieldReservedCount)
	}
}

func TestAb93PWMAndOpaque(b523SensorFieldT *b523SensorFieldTesting.T) {
	for _, b523SensorFieldField := range b523SensorFieldCatalog(b523SensorFieldT).B523SensorFieldFields {
		if b523SensorFieldField.B523SensorFieldKey.B523SensorFieldField == "PWM" {
			b523SensorFieldBytes := make([]byte, 20)
			b523SensorFieldBytes[15] = 255
			_, b523SensorFieldValue, _, _ := b523SensorFieldEvaluate(b523SensorFieldT, b523SensorFieldField, b523SensorFieldBytes)
			if b523SensorFieldValue != 127.5 {
				b523SensorFieldT.Fatal("PWM clipped")
			}
		}
		if *b523SensorFieldField.B523SensorFieldKey.B523SensorFieldOperation == 7 && (b523SensorFieldField.B523SensorFieldKey.B523SensorFieldField == "opaque" || b523SensorFieldField.B523SensorFieldKey.B523SensorFieldField == "S12" || b523SensorFieldField.B523SensorFieldKey.B523SensorFieldField == "S13") {
			b523SensorFieldBytes := make([]byte, 20)
			b523SensorFieldBytes[11] = 255
			b523SensorFieldBytes[12] = 16
			b523SensorFieldBytes[14] = 240
			b523SensorFieldBytes[15] = 255
			_, b523SensorFieldValue, _, _ := b523SensorFieldEvaluate(b523SensorFieldT, b523SensorFieldField, b523SensorFieldBytes)
			b523SensorFieldWant := map[string]float64{"opaque": 255, "S12": 1, "S13": -1}[b523SensorFieldField.B523SensorFieldKey.B523SensorFieldField]
			if b523SensorFieldValue != b523SensorFieldWant {
				b523SensorFieldT.Fatal("OP07 intermediate byte conflated with sensor")
			}
		}
	}
}

func b523SensorFieldCopyKey(b523SensorFieldKey B523SensorKey) B523SensorKey {
	b523SensorFieldData, _ := b523SensorFieldJSON.Marshal(b523SensorFieldKey)
	var b523SensorFieldCopy B523SensorKey
	if b523SensorFieldError := b523SensorFieldJSON.Unmarshal(b523SensorFieldData, &b523SensorFieldCopy); b523SensorFieldError != nil {
		panic(b523SensorFieldError)
	}
	return b523SensorFieldCopy
}

func TestAb93ExactIdentityAxes(b523SensorFieldT *b523SensorFieldTesting.T) {
	b523SensorFieldOriginal := b523SensorFieldCatalog(b523SensorFieldT).B523SensorFieldFields[0]
	for _, b523SensorFieldB524 := range []bool{false, true} {
		b523SensorFieldControl := b523SensorFieldOriginal
		b523SensorFieldControl.B523SensorFieldKey = b523SensorFieldCopyKey(b523SensorFieldOriginal.B523SensorFieldKey)
		if b523SensorFieldB524 {
			b523SensorFieldControl.B523SensorFieldKey.B523SensorFieldSecondary = b523SensorFieldPtr(uint16(36))
			b523SensorFieldControl.B523SensorFieldKey.B523SensorFieldOperation = nil
			b523SensorFieldControl.B523SensorFieldKey.B523SensorFieldOpcode = b523SensorFieldPtr(uint16(0))
			b523SensorFieldControl.B523SensorFieldKey.B523SensorFieldGroup = b523SensorFieldPtr(uint16(0))
			b523SensorFieldControl.B523SensorFieldKey.B523SensorFieldInstance = b523SensorFieldPtr(uint16(0))
			b523SensorFieldControl.B523SensorFieldKey.B523SensorFieldRegister = b523SensorFieldPtr(uint16(0))
		}
		b523SensorFieldFields := []B523SensorField{b523SensorFieldControl}
		b523SensorFieldEqual := b523SensorFieldCopyKey(b523SensorFieldControl.B523SensorFieldKey)
		if _, b523SensorFieldFound, b523SensorFieldError := b523SensorFieldFindSensorField(b523SensorFieldFields, b523SensorFieldEqual); !b523SensorFieldFound || b523SensorFieldError != nil {
			b523SensorFieldT.Fatal("valid control failed", b523SensorFieldError)
		}
		for _, b523SensorFieldAxis := range []string{"B523SensorFieldFamily", "B523SensorFieldVariant", "B523SensorFieldPrimary", "B523SensorFieldSecondary", "B523SensorFieldOperation", "B523SensorFieldDirection", "B523SensorFieldField", "B523SensorFieldOpcode", "B523SensorFieldGroup", "B523SensorFieldInstance", "B523SensorFieldRegister"} {
			b523SensorFieldT.Run(b523SensorFieldFmt.Sprintf("ab93_%t_%s", b523SensorFieldB524, b523SensorFieldAxis), func(b523SensorFieldT *b523SensorFieldTesting.T) {
				b523SensorFieldKey := b523SensorFieldCopyKey(b523SensorFieldControl.B523SensorFieldKey)
				b523SensorFieldValue := b523SensorFieldReflect.ValueOf(&b523SensorFieldKey).Elem().FieldByName(b523SensorFieldAxis)
				if b523SensorFieldValue.Kind() == b523SensorFieldReflect.String {
					b523SensorFieldValue.SetString("ab93-other")
					if _, b523SensorFieldFound, b523SensorFieldError := b523SensorFieldFindSensorField(b523SensorFieldFields, b523SensorFieldKey); b523SensorFieldFound || b523SensorFieldError != nil {
						b523SensorFieldT.Fatal("ignored string axis", b523SensorFieldError)
					}
					return
				}
				// Each nullable coordinate is tested absent, explicitly zero and nonzero.
				// The current control must match; every other value must miss.
				b523SensorFieldCandidates := []b523SensorFieldReflect.Value{b523SensorFieldReflect.Zero(b523SensorFieldValue.Type())}
				if b523SensorFieldValue.Type().Elem().Kind() == b523SensorFieldReflect.String {
					b523SensorFieldCandidates = append(b523SensorFieldCandidates, b523SensorFieldReflect.ValueOf(b523SensorFieldPtr("")), b523SensorFieldReflect.ValueOf(b523SensorFieldPtr("response")))
				} else {
					b523SensorFieldCandidates = append(b523SensorFieldCandidates, b523SensorFieldReflect.ValueOf(b523SensorFieldPtr(uint16(0))), b523SensorFieldReflect.ValueOf(b523SensorFieldPtr(uint16(123))))
				}
				for _, b523SensorFieldCandidate := range b523SensorFieldCandidates {
					b523SensorFieldValue.Set(b523SensorFieldCandidate)
					b523SensorFieldWant := b523SensorFieldReflect.DeepEqual(b523SensorFieldKey, b523SensorFieldControl.B523SensorFieldKey)
					_, b523SensorFieldFound, b523SensorFieldError := b523SensorFieldFindSensorField(b523SensorFieldFields, b523SensorFieldKey)
					if b523SensorFieldFound != b523SensorFieldWant || b523SensorFieldError != nil {
						b523SensorFieldT.Fatal("nullable exact identity failed", b523SensorFieldError)
					}
					if !b523SensorFieldB524 {
						_, b523SensorFieldFound, b523SensorFieldError = B523LookupSensorField(b523SensorFieldKey)
						if b523SensorFieldFound != b523SensorFieldWant || b523SensorFieldError != nil {
							b523SensorFieldT.Fatal("production nullable identity failed", b523SensorFieldError)
						}
					}
				}
			})
		}
	}
}

func b523SensorFieldMutateField(b523SensorFieldField *B523SensorField) {
	b523SensorFieldField.B523SensorFieldKey.B523SensorFieldFamily = "ab93-mutated"
	*b523SensorFieldField.B523SensorFieldKey.B523SensorFieldPrimary = 0
	*b523SensorFieldField.B523SensorFieldKey.B523SensorFieldSecondary = 0
	*b523SensorFieldField.B523SensorFieldKey.B523SensorFieldOperation = 0
	*b523SensorFieldField.B523SensorFieldPosition = 0
	*b523SensorFieldField.B523SensorFieldWidth = 0
	if b523SensorFieldField.B523SensorFieldBitOffset != nil {
		*b523SensorFieldField.B523SensorFieldBitOffset = 99
	}
	if b523SensorFieldField.B523SensorFieldDefault != nil {
		*b523SensorFieldField.B523SensorFieldDefault = 99
	}
	if b523SensorFieldField.B523SensorFieldDivideBy != nil {
		*b523SensorFieldField.B523SensorFieldDivideBy = 99
	}
	if b523SensorFieldField.B523SensorFieldUnit != nil {
		*b523SensorFieldField.B523SensorFieldUnit = "ab93-mutated"
	}
	if b523SensorFieldField.B523SensorFieldReplacementRaw != nil {
		*b523SensorFieldField.B523SensorFieldReplacementRaw = 0
	}
	for b523SensorFieldIndex := range b523SensorFieldField.B523SensorFieldBounds {
		b523SensorFieldField.B523SensorFieldBounds[b523SensorFieldIndex] = 99
	}
	b523SensorFieldField.B523SensorFieldEnumeration["0"] = "ab93-mutated"
	for b523SensorFieldIndex := range b523SensorFieldField.B523SensorFieldBranches {
		b523SensorFieldBranch := &b523SensorFieldField.B523SensorFieldBranches[b523SensorFieldIndex]
		b523SensorFieldBranch.B523SensorFieldName = "ab93-mutated"
		*b523SensorFieldBranch.B523SensorFieldDefault = 99
		*b523SensorFieldBranch.B523SensorFieldDivideBy = 99
		if b523SensorFieldBranch.B523SensorFieldUnit != nil {
			*b523SensorFieldBranch.B523SensorFieldUnit = "ab93-mutated"
		}
		b523SensorFieldBranch.B523SensorFieldBounds[0] = 99
	}
}

func TestAb93DetachedAccessors(b523SensorFieldT *b523SensorFieldTesting.T) {
	b523SensorFieldOriginal := b523SensorFieldCatalog(b523SensorFieldT)
	b523SensorFieldChanged := b523SensorFieldCatalog(b523SensorFieldT)
	b523SensorFieldChanged.B523SensorFieldLimitations[0] = "ab93-mutated"
	for b523SensorFieldIndex := range b523SensorFieldChanged.B523SensorFieldFields {
		b523SensorFieldMutateField(&b523SensorFieldChanged.B523SensorFieldFields[b523SensorFieldIndex])
	}
	if !b523SensorFieldReflect.DeepEqual(b523SensorFieldOriginal, b523SensorFieldCatalog(b523SensorFieldT)) {
		b523SensorFieldT.Fatal("catalog shares backing storage")
	}
	for _, b523SensorFieldOriginalField := range b523SensorFieldOriginal.B523SensorFieldFields {
		b523SensorFieldChangedField, b523SensorFieldFound, b523SensorFieldError := B523LookupSensorField(b523SensorFieldOriginalField.B523SensorFieldKey)
		if !b523SensorFieldFound || b523SensorFieldError != nil {
			b523SensorFieldT.Fatal(b523SensorFieldError)
		}
		b523SensorFieldMutateField(&b523SensorFieldChangedField)
		b523SensorFieldFresh, b523SensorFieldFound, b523SensorFieldError := B523LookupSensorField(b523SensorFieldOriginalField.B523SensorFieldKey)
		if !b523SensorFieldFound || b523SensorFieldError != nil || !b523SensorFieldReflect.DeepEqual(b523SensorFieldFresh, b523SensorFieldOriginalField) {
			b523SensorFieldT.Fatal("lookup shares backing storage", b523SensorFieldError)
		}
		b523SensorFieldChangedField, b523SensorFieldFound, b523SensorFieldError = b523SensorFieldFindSensorField(b523SensorFieldOriginal.B523SensorFieldFields, b523SensorFieldOriginalField.B523SensorFieldKey)
		if !b523SensorFieldFound || b523SensorFieldError != nil {
			b523SensorFieldT.Fatal(b523SensorFieldError)
		}
		b523SensorFieldMutateField(&b523SensorFieldChangedField)
		b523SensorFieldFresh, b523SensorFieldFound, b523SensorFieldError = b523SensorFieldFindSensorField(b523SensorFieldOriginal.B523SensorFieldFields, b523SensorFieldOriginalField.B523SensorFieldKey)
		if !b523SensorFieldFound || b523SensorFieldError != nil || !b523SensorFieldReflect.DeepEqual(b523SensorFieldFresh, b523SensorFieldOriginalField) {
			b523SensorFieldT.Fatal("find aliases input", b523SensorFieldError)
		}
	}
	if !b523SensorFieldReflect.DeepEqual(b523SensorFieldOriginal, b523SensorFieldCatalog(b523SensorFieldT)) {
		b523SensorFieldT.Fatal("nested lookup mutation reached catalog")
	}
}
