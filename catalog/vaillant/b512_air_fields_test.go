package vaillant

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

func b512AirFieldGet(b512AirFieldT *testing.T, b512AirFieldOperation uint16, b512AirFieldDirection, b512AirFieldName string) B512AirField {
	b512AirFieldT.Helper()
	b512AirFieldField, b512AirFieldOK := B512LookupAirField(B512AirKey{B512AirFieldFamily: "vaillant", B512AirFieldVariant: "B512", B512AirFieldPrimary: 0xb512, B512AirFieldOperation: B512AirFieldAxis{B512AirFieldPresent: true, B512AirFieldValue: b512AirFieldOperation}, B512AirFieldDirection: b512AirFieldDirection, B512AirFieldField: b512AirFieldName})
	if !b512AirFieldOK {
		b512AirFieldT.Fatalf("ab88 missing %d/%s/%s", b512AirFieldOperation, b512AirFieldDirection, b512AirFieldName)
	}
	return b512AirFieldField
}

func TestAb88ProductionFields(b512AirFieldT *testing.T) {
	b512AirFieldExpected := []struct {
		b512AirFieldOperation                   uint16
		b512AirFieldDirection, b512AirFieldName string
		b512AirFieldPosition, b512AirFieldWidth int
		b512AirFieldSigned                      bool
		b512AirFieldDivisor                     float64
		b512AirFieldUnit                        string
		b512AirFieldDefault                     *float64
	}{}
	b512AirFieldZero, b512AirFieldOne := 0.0, 1.0
	b512AirFieldAdd := func(b512AirFieldOperation uint16, b512AirFieldDirection, b512AirFieldName string, b512AirFieldPosition, b512AirFieldWidth int, b512AirFieldSigned bool, b512AirFieldDivisor float64, b512AirFieldUnit string, b512AirFieldDefault *float64) {
		b512AirFieldExpected = append(b512AirFieldExpected, struct {
			b512AirFieldOperation                   uint16
			b512AirFieldDirection, b512AirFieldName string
			b512AirFieldPosition, b512AirFieldWidth int
			b512AirFieldSigned                      bool
			b512AirFieldDivisor                     float64
			b512AirFieldUnit                        string
			b512AirFieldDefault                     *float64
		}{b512AirFieldOperation, b512AirFieldDirection, b512AirFieldName, b512AirFieldPosition, b512AirFieldWidth, b512AirFieldSigned, b512AirFieldDivisor, b512AirFieldUnit, b512AirFieldDefault})
	}
	b512AirFieldAdd(11, "request", "target_fan_1", 3, 2, false, 1, "rpm", &b512AirFieldZero)
	b512AirFieldAdd(11, "request", "target_fan_2", 5, 2, false, 1, "rpm", &b512AirFieldZero)
	b512AirFieldAdd(11, "request", "relay", 7, 1, false, 1, "", &b512AirFieldZero)
	b512AirFieldAdd(11, "positive", "actual_fan_1", 3, 2, false, 1, "rpm", &b512AirFieldZero)
	b512AirFieldAdd(11, "positive", "actual_fan_2", 5, 2, false, 1, "rpm", &b512AirFieldZero)
	b512AirFieldAdd(11, "positive", "status", 7, 1, false, 1, "", &b512AirFieldZero)
	b512AirFieldAdd(11, "positive", "temperature_1", 8, 2, true, 16, "degrees Celsius", &b512AirFieldZero)
	b512AirFieldAdd(11, "positive", "temperature_2", 10, 2, true, 16, "degrees Celsius", &b512AirFieldZero)
	b512AirFieldAdd(11, "positive", "temperature_3", 12, 2, true, 16, "degrees Celsius", &b512AirFieldZero)
	for b512AirFieldIndex, b512AirFieldName := range []string{"pulse", "P", "I", "D"} {
		b512AirFieldAdd(12, "request", b512AirFieldName, 3+b512AirFieldIndex, 1, false, 1, "", nil)
	}
	b512AirFieldAdd(12, "request", "minimum_pwm", 7, 2, true, 256, "percent", &b512AirFieldOne)
	b512AirFieldAdd(12, "request", "maximum_pwm", 9, 2, true, 256, "percent", &b512AirFieldOne)
	b512AirFieldAdd(12, "request", "fan_count", 11, 1, false, 1, "", nil)
	b512AirFieldAdd(12, "positive", "feedback", 3, 1, false, 1, "", nil)
	b512AirFieldAdd(13, "positive", "digital", 3, 1, false, 1, "", &b512AirFieldZero)
	b512AirFieldAdd(13, "positive", "temperature_1", 4, 2, true, 16, "degrees Celsius", nil)
	b512AirFieldAdd(13, "positive", "temperature_2", 6, 2, true, 16, "degrees Celsius", nil)
	b512AirFieldAdd(14, "positive", "temperature", 3, 2, true, 16, "degrees Celsius", &b512AirFieldZero)
	b512AirFieldAdd(14, "positive", "humidity", 5, 2, false, 16, "percent", &b512AirFieldZero)
	b512AirFieldAdd(14, "positive", "air_quality", 7, 2, false, 1, "", &b512AirFieldZero)
	b512AirFieldMetadata := B512AirFields()
	if len(b512AirFieldMetadata.B512AirFieldFields) != len(b512AirFieldExpected) {
		b512AirFieldT.Fatal("ab88 finite field count")
	}
	b512AirFieldSeen := map[B512AirKey]bool{}
	for _, b512AirFieldWant := range b512AirFieldExpected {
		b512AirFieldT.Run(fmt.Sprintf("ab88_%d_%s_%s", b512AirFieldWant.b512AirFieldOperation, b512AirFieldWant.b512AirFieldDirection, b512AirFieldWant.b512AirFieldName), func(b512AirFieldT *testing.T) {
			b512AirFieldGot := b512AirFieldGet(b512AirFieldT, b512AirFieldWant.b512AirFieldOperation, b512AirFieldWant.b512AirFieldDirection, b512AirFieldWant.b512AirFieldName)
			if b512AirFieldSeen[b512AirFieldGot.B512AirFieldKey] {
				b512AirFieldT.Fatal("ab88 duplicate")
			}
			b512AirFieldSeen[b512AirFieldGot.B512AirFieldKey] = true
			b512AirFieldSignedness, b512AirFieldOrder, b512AirFieldConversion := "unsigned", "not_applicable", "identity"
			if b512AirFieldWant.b512AirFieldSigned {
				b512AirFieldSignedness = "signed"
			}
			if b512AirFieldWant.b512AirFieldWidth == 2 {
				b512AirFieldOrder = "little_endian"
			}
			if b512AirFieldWant.b512AirFieldDivisor != 1 {
				b512AirFieldConversion = "divide"
			}
			if b512AirFieldGot.B512AirFieldLogicalPosition == nil || *b512AirFieldGot.B512AirFieldLogicalPosition != b512AirFieldWant.b512AirFieldPosition || b512AirFieldGot.B512AirFieldWidth == nil || *b512AirFieldGot.B512AirFieldWidth != b512AirFieldWant.b512AirFieldWidth || b512AirFieldGot.B512AirFieldSignedness != b512AirFieldSignedness || b512AirFieldGot.B512AirFieldByteOrder != b512AirFieldOrder || b512AirFieldGot.B512AirFieldConversion != b512AirFieldConversion || b512AirFieldGot.B512AirFieldDivisor != b512AirFieldWant.b512AirFieldDivisor || !reflect.DeepEqual(b512AirFieldGot.B512AirFieldDefault, b512AirFieldWant.b512AirFieldDefault) {
				b512AirFieldT.Fatalf("ab88 layout/conversion/default: %+v", b512AirFieldGot)
			}
			if (b512AirFieldGot.B512AirFieldUnit == nil) != (b512AirFieldWant.b512AirFieldUnit == "") || (b512AirFieldGot.B512AirFieldUnit != nil && *b512AirFieldGot.B512AirFieldUnit != b512AirFieldWant.b512AirFieldUnit) {
				b512AirFieldT.Fatal("ab88 unit")
			}
			if b512AirFieldGot.B512AirFieldGrade != "Hypothesis" || b512AirFieldGot.B512AirFieldRawFallback != "preserve original field bytes including unlabeled and reserved values" || b512AirFieldGot.B512AirFieldReplacement != nil {
				b512AirFieldT.Fatal("ab88 evidence or replacement")
			}
			var b512AirFieldDefaultLabel *string
			if b512AirFieldWant.b512AirFieldName == "feedback" {
				b512AirFieldRejected := "rejected"
				b512AirFieldDefaultLabel = &b512AirFieldRejected
			}
			if !reflect.DeepEqual(b512AirFieldGot.B512AirFieldDefaultLabel, b512AirFieldDefaultLabel) {
				b512AirFieldT.Fatal("ab88 default label")
			}
			b512AirFieldBounds := []float64(nil)
			if b512AirFieldWant.b512AirFieldName == "minimum_pwm" || b512AirFieldWant.b512AirFieldName == "maximum_pwm" {
				b512AirFieldBounds = []float64{-128, 127.99609375}
			}
			if !reflect.DeepEqual(b512AirFieldGot.B512AirFieldBounds, b512AirFieldBounds) {
				b512AirFieldT.Fatal("ab88 bounds")
			}
			b512AirFieldEnum := map[string]string{}
			var b512AirFieldOtherwise *string
			switch b512AirFieldWant.b512AirFieldName {
			case "status":
				b512AirFieldEnum = map[string]string{"0": "idle", "1": "running", "255": "fault"}
			case "feedback":
				b512AirFieldEnum = map[string]string{"1": "accepted"}
				b512AirFieldRejected := "rejected"
				b512AirFieldOtherwise = &b512AirFieldRejected
			case "digital":
				b512AirFieldEnum = map[string]string{"0": "false", "1": "true"}
			}
			if !reflect.DeepEqual(b512AirFieldGot.B512AirFieldEnumeration, b512AirFieldEnum) || !reflect.DeepEqual(b512AirFieldGot.B512AirFieldOtherwise, b512AirFieldOtherwise) {
				b512AirFieldT.Fatal("ab88 enumeration")
			}
			b512AirFieldBits := map[string]map[string]string{}
			b512AirFieldReserved := []int{}
			if b512AirFieldWant.b512AirFieldName == "relay" {
				b512AirFieldBits = map[string]map[string]string{"0": {"1": "15V"}, "1": {"1": "tray-heater relay"}}
				b512AirFieldReserved = []int{2, 3, 4, 5, 6, 7}
			}
			if !reflect.DeepEqual(b512AirFieldGot.B512AirFieldBits, b512AirFieldBits) || !reflect.DeepEqual(b512AirFieldGot.B512AirFieldReservedBits, b512AirFieldReserved) {
				b512AirFieldT.Fatal("ab88 relay/reserved roles")
			}
		})
	}
	if b512AirFieldMetadata.B512AirFieldRepresentation != "Proven" || b512AirFieldMetadata.B512AirFieldPhysicalSupport != "Unknown" || b512AirFieldMetadata.B512AirFieldAcknowledgement != "Unknown" || b512AirFieldMetadata.B512AirFieldComplete || b512AirFieldMetadata.B512AirFieldCoordinates != "logical positions; not wire offsets" || len(b512AirFieldMetadata.B512AirFieldLimitations) != 4 {
		b512AirFieldT.Fatal("ab88 qualification boundary")
	}
	if !reflect.DeepEqual(b512AirFieldMetadata.B512AirFieldOperations, []B512AirOperation{{11, true, "Hypothesis"}, {12, true, "Hypothesis"}, {13, false, "Hypothesis"}, {14, false, "Hypothesis"}}) {
		b512AirFieldT.Fatal("ab88 no-argument operations")
	}
}

func b512AirFieldEvaluate(b512AirFieldField B512AirField, b512AirFieldBytes []byte) (float64, bool) {
	if b512AirFieldField.B512AirFieldWidth == nil || len(b512AirFieldBytes) != *b512AirFieldField.B512AirFieldWidth {
		return 0, false
	}
	b512AirFieldRaw := int64(b512AirFieldBytes[0])
	if *b512AirFieldField.B512AirFieldWidth == 2 {
		if b512AirFieldField.B512AirFieldByteOrder != "little_endian" {
			return 0, false
		}
		b512AirFieldRaw |= int64(b512AirFieldBytes[1]) << 8
		if b512AirFieldField.B512AirFieldSignedness == "signed" {
			b512AirFieldRaw = int64(int16(b512AirFieldRaw))
		}
	}
	b512AirFieldValue := float64(b512AirFieldRaw)
	switch b512AirFieldField.B512AirFieldConversion {
	case "divide":
		b512AirFieldValue /= b512AirFieldField.B512AirFieldDivisor
	case "identity":
	default:
		return 0, false
	}
	return b512AirFieldValue, true
}

func TestAb88SyntheticValues(b512AirFieldT *testing.T) {
	for _, b512AirFieldCase := range []struct {
		b512AirFieldOp                    uint16
		b512AirFieldDir, b512AirFieldName string
		b512AirFieldBytes                 []byte
		b512AirFieldWant                  float64
	}{
		{11, "positive", "temperature_1", []byte{0x10, 0}, 1}, {11, "positive", "temperature_2", []byte{0xf0, 0xff}, -1},
		{12, "request", "minimum_pwm", []byte{0, 1}, 1}, {12, "request", "minimum_pwm", []byte{0, 0x80}, -128}, {12, "request", "maximum_pwm", []byte{0xff, 0x7f}, 127.99609375},
		{14, "positive", "humidity", []byte{0x40, 6}, 100}, {14, "positive", "humidity", []byte{0x50, 6}, 101}, {14, "positive", "air_quality", []byte{0xff, 0xff}, 65535},
	} {
		b512AirFieldField := b512AirFieldGet(b512AirFieldT, b512AirFieldCase.b512AirFieldOp, b512AirFieldCase.b512AirFieldDir, b512AirFieldCase.b512AirFieldName)
		if b512AirFieldGot, b512AirFieldOK := b512AirFieldEvaluate(b512AirFieldField, b512AirFieldCase.b512AirFieldBytes); !b512AirFieldOK || b512AirFieldGot != b512AirFieldCase.b512AirFieldWant {
			b512AirFieldT.Fatalf("ab88 %s got %v want %v", b512AirFieldCase.b512AirFieldName, b512AirFieldGot, b512AirFieldCase.b512AirFieldWant)
		}
		for _, b512AirFieldInvalid := range [][]byte{b512AirFieldCase.b512AirFieldBytes[:1], append(append([]byte{}, b512AirFieldCase.b512AirFieldBytes...), 0)} {
			if _, b512AirFieldOK := b512AirFieldEvaluate(b512AirFieldField, b512AirFieldInvalid); b512AirFieldOK {
				b512AirFieldT.Fatal("ab88 field extent")
			}
		}
	}
	for _, b512AirFieldDir := range []string{"request", "positive"} {
		b512AirFieldLogical := []byte{0, 0, 0, 2, 1, 4, 3, 0, 0x10, 0, 0xf0, 0xff, 0, 0}
		b512AirFieldPrefix := "target"
		if b512AirFieldDir == "positive" {
			b512AirFieldPrefix = "actual"
		}
		for b512AirFieldIndex, b512AirFieldWant := range []float64{258, 772} {
			b512AirFieldField := b512AirFieldGet(b512AirFieldT, 11, b512AirFieldDir, fmt.Sprintf("%s_fan_%d", b512AirFieldPrefix, b512AirFieldIndex+1))
			b512AirFieldPosition := *b512AirFieldField.B512AirFieldLogicalPosition
			if b512AirFieldGot, b512AirFieldOK := b512AirFieldEvaluate(b512AirFieldField, b512AirFieldLogical[b512AirFieldPosition:b512AirFieldPosition+*b512AirFieldField.B512AirFieldWidth]); !b512AirFieldOK || b512AirFieldGot != b512AirFieldWant {
				b512AirFieldT.Fatal("ab88 fan identity")
			}
		}
	}
	b512AirFieldRelay := b512AirFieldGet(b512AirFieldT, 11, "request", "relay")
	if b512AirFieldRelay.B512AirFieldBits["2"] != nil || b512AirFieldRelay.B512AirFieldReservedBits[0] != 2 {
		b512AirFieldT.Fatal("ab88 bit 2 reserved")
	}
	b512AirFieldFeedback := b512AirFieldGet(b512AirFieldT, 12, "positive", "feedback")
	for b512AirFieldRaw, b512AirFieldWant := range map[int]string{1: "accepted", 2: "rejected"} {
		b512AirFieldLabel := b512AirFieldFeedback.B512AirFieldEnumeration[strconv.Itoa(b512AirFieldRaw)]
		if b512AirFieldLabel == "" && b512AirFieldFeedback.B512AirFieldOtherwise != nil {
			b512AirFieldLabel = *b512AirFieldFeedback.B512AirFieldOtherwise
		}
		if b512AirFieldLabel != b512AirFieldWant {
			b512AirFieldT.Fatal("ab88 feedback")
		}
	}
	if b512AirFieldGet(b512AirFieldT, 13, "positive", "digital").B512AirFieldEnumeration["2"] != "" {
		b512AirFieldT.Fatal("ab88 digital 2 unlabeled")
	}
}

func TestAb88RelayBitRoles(b512AirFieldT *testing.T) {
	b512AirFieldRelay := b512AirFieldGet(b512AirFieldT, 11, "request", "relay")
	for _, b512AirFieldCase := range []struct {
		b512AirFieldRaw  byte
		b512AirFieldWant []string
	}{
		{0x00, []string{}},
		{0x01, []string{"15V"}},
		{0x02, []string{"tray-heater relay"}},
		{0x03, []string{"15V", "tray-heater relay"}},
		{0x04, []string{}},
		{0x05, []string{"15V"}},
		{0x06, []string{"tray-heater relay"}},
		{0x07, []string{"15V", "tray-heater relay"}},
	} {
		b512AirFieldT.Run(fmt.Sprintf("ab88_%02x", b512AirFieldCase.b512AirFieldRaw), func(b512AirFieldT *testing.T) {
			b512AirFieldLabels := []string{}
			for b512AirFieldBit := 0; b512AirFieldBit < 8; b512AirFieldBit++ {
				b512AirFieldState := (b512AirFieldCase.b512AirFieldRaw >> b512AirFieldBit) & 1
				b512AirFieldLabel := b512AirFieldRelay.B512AirFieldBits[strconv.Itoa(b512AirFieldBit)][strconv.Itoa(int(b512AirFieldState))]
				if b512AirFieldLabel != "" {
					b512AirFieldLabels = append(b512AirFieldLabels, b512AirFieldLabel)
				}
			}
			if !reflect.DeepEqual(b512AirFieldLabels, b512AirFieldCase.b512AirFieldWant) {
				b512AirFieldT.Fatalf("ab88 relay roles got %v want %v", b512AirFieldLabels, b512AirFieldCase.b512AirFieldWant)
			}
		})
	}
}

func TestAb88ExactIdentity(b512AirFieldT *testing.T) {
	b512AirFieldControl := b512AirFieldGet(b512AirFieldT, 12, "request", "minimum_pwm")
	for _, b512AirFieldB524 := range []bool{false, true} {
		b512AirFieldField := b512AirFieldControl
		if b512AirFieldB524 {
			b512AirFieldField.B512AirFieldKey.B512AirFieldVariant = "B524"
			b512AirFieldField.B512AirFieldKey.B512AirFieldPrimary = 0xb524
			b512AirFieldField.B512AirFieldKey.B512AirFieldOpcode = B512AirFieldAxis{true, 1}
			b512AirFieldField.B512AirFieldKey.B512AirFieldGroup = B512AirFieldAxis{true, 2}
			b512AirFieldField.B512AirFieldKey.B512AirFieldInstance = B512AirFieldAxis{true, 3}
			b512AirFieldField.B512AirFieldKey.B512AirFieldRegister = B512AirFieldAxis{true, 4}
		}
		b512AirFieldKey := b512AirFieldField.B512AirFieldKey
		for b512AirFieldIndex := 0; b512AirFieldIndex < reflect.TypeOf(b512AirFieldKey).NumField(); b512AirFieldIndex++ {
			b512AirFieldT.Run(fmt.Sprintf("ab88_%t_%s", b512AirFieldB524, reflect.TypeOf(b512AirFieldKey).Field(b512AirFieldIndex).Name), func(b512AirFieldT *testing.T) {
				b512AirFieldChanged := b512AirFieldKey
				b512AirFieldAxis := reflect.ValueOf(&b512AirFieldChanged).Elem().Field(b512AirFieldIndex)
				switch b512AirFieldAxis.Kind() {
				case reflect.String:
					b512AirFieldAxis.SetString(b512AirFieldAxis.String() + "_ab88_other")
				case reflect.Uint16:
					b512AirFieldAxis.SetUint(b512AirFieldAxis.Uint() + 1)
				case reflect.Struct:
					b512AirFieldAxis.FieldByName("B512AirFieldValue").SetUint(b512AirFieldAxis.FieldByName("B512AirFieldValue").Uint() + 1)
				}
				if _, b512AirFieldOK := b512AirFieldFindAirField([]B512AirField{b512AirFieldField}, b512AirFieldChanged); b512AirFieldOK {
					b512AirFieldT.Fatal("ab88 axis ignored")
				}
				if !b512AirFieldB524 {
					if _, b512AirFieldOK := B512LookupAirField(b512AirFieldChanged); b512AirFieldOK {
						b512AirFieldT.Fatal("ab88 public wildcard")
					}
				}
				if b512AirFieldAxis.Kind() == reflect.Struct {
					b512AirFieldChanged = b512AirFieldKey
					b512AirFieldAxis = reflect.ValueOf(&b512AirFieldChanged).Elem().Field(b512AirFieldIndex)
					b512AirFieldAxis.FieldByName("B512AirFieldPresent").SetBool(!b512AirFieldAxis.FieldByName("B512AirFieldPresent").Bool())
					if _, b512AirFieldOK := b512AirFieldFindAirField([]B512AirField{b512AirFieldField}, b512AirFieldChanged); b512AirFieldOK {
						b512AirFieldT.Fatal("ab88 absent equals present")
					}
					b512AirFieldChanged = b512AirFieldKey
					b512AirFieldAxis = reflect.ValueOf(&b512AirFieldChanged).Elem().Field(b512AirFieldIndex)
					b512AirFieldAxis.FieldByName("B512AirFieldValue").SetUint(0)
					b512AirFieldAxis.FieldByName("B512AirFieldPresent").SetBool(true)
					b512AirFieldZero := b512AirFieldField
					b512AirFieldZero.B512AirFieldKey = b512AirFieldChanged
					if _, b512AirFieldOK := b512AirFieldFindAirField([]B512AirField{b512AirFieldZero}, b512AirFieldChanged); !b512AirFieldOK {
						b512AirFieldT.Fatal("ab88 explicit zero control")
					}
					b512AirFieldAxis.FieldByName("B512AirFieldPresent").SetBool(false)
					if _, b512AirFieldOK := b512AirFieldFindAirField([]B512AirField{b512AirFieldZero}, b512AirFieldChanged); b512AirFieldOK {
						b512AirFieldT.Fatal("ab88 zero presence ignored")
					}
				}
				if _, b512AirFieldOK := b512AirFieldFindAirField([]B512AirField{b512AirFieldField}, b512AirFieldKey); !b512AirFieldOK {
					b512AirFieldT.Fatal("ab88 valid control")
				}
			})
		}
	}
}

func b512AirFieldMutate(b512AirFieldField *B512AirField) {
	*b512AirFieldField.B512AirFieldLogicalPosition = 99
	*b512AirFieldField.B512AirFieldWidth = 99
	if b512AirFieldField.B512AirFieldUnit != nil {
		*b512AirFieldField.B512AirFieldUnit = "ab88_changed"
	}
	if b512AirFieldField.B512AirFieldDefault != nil {
		*b512AirFieldField.B512AirFieldDefault = 99
	}
	if b512AirFieldField.B512AirFieldDefaultLabel != nil {
		*b512AirFieldField.B512AirFieldDefaultLabel = "ab88_changed"
	}
	if b512AirFieldField.B512AirFieldOtherwise != nil {
		*b512AirFieldField.B512AirFieldOtherwise = "ab88_changed"
	}
	for b512AirFieldIndex := range b512AirFieldField.B512AirFieldBounds {
		b512AirFieldField.B512AirFieldBounds[b512AirFieldIndex] = 99
	}
	for b512AirFieldKey := range b512AirFieldField.B512AirFieldEnumeration {
		b512AirFieldField.B512AirFieldEnumeration[b512AirFieldKey] = "ab88_changed"
	}
	for _, b512AirFieldBits := range b512AirFieldField.B512AirFieldBits {
		for b512AirFieldKey := range b512AirFieldBits {
			b512AirFieldBits[b512AirFieldKey] = "ab88_changed"
		}
	}
	for b512AirFieldIndex := range b512AirFieldField.B512AirFieldReservedBits {
		b512AirFieldField.B512AirFieldReservedBits[b512AirFieldIndex] = 99
	}
}

func TestAb88CopyIsolation(b512AirFieldT *testing.T) {
	b512AirFieldBefore := B512AirFields()
	b512AirFieldCopy := B512AirFields()
	b512AirFieldCopy.B512AirFieldLimitations[0] = "ab88_changed"
	b512AirFieldCopy.B512AirFieldOperations[0].B512AirFieldGrade = "ab88_changed"
	for b512AirFieldIndex := range b512AirFieldCopy.B512AirFieldFields {
		b512AirFieldMutate(&b512AirFieldCopy.B512AirFieldFields[b512AirFieldIndex])
		b512AirFieldLookup, b512AirFieldOK := B512LookupAirField(b512AirFieldBefore.B512AirFieldFields[b512AirFieldIndex].B512AirFieldKey)
		if !b512AirFieldOK {
			b512AirFieldT.Fatal("ab88 missing copy")
		}
		b512AirFieldMutate(&b512AirFieldLookup)
		b512AirFieldAgain, b512AirFieldOK := B512LookupAirField(b512AirFieldBefore.B512AirFieldFields[b512AirFieldIndex].B512AirFieldKey)
		if !b512AirFieldOK || !reflect.DeepEqual(b512AirFieldAgain, b512AirFieldBefore.B512AirFieldFields[b512AirFieldIndex]) {
			b512AirFieldT.Fatal("ab88 lookup shares storage")
		}
	}
	if !reflect.DeepEqual(B512AirFields(), b512AirFieldBefore) {
		b512AirFieldT.Fatal("ab88 accessor shares storage")
	}
}
