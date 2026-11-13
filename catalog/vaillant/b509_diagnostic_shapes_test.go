package vaillant

import (
	b509DiagnosticShapeJSON "encoding/json"
	b509DiagnosticShapeReflect "reflect"
	b509DiagnosticShapeTesting "testing"
)

func b509DiagnosticShapeNumber(b509DiagnosticShapeValue uint16) *uint16 {
	return &b509DiagnosticShapeValue
}
func b509DiagnosticShapeInt(b509DiagnosticShapeValue int) *int { return &b509DiagnosticShapeValue }
func b509DiagnosticShapeBound(b509DiagnosticShapeValue int64) *int64 {
	return &b509DiagnosticShapeValue
}

func b509DiagnosticShapeKey(b509DiagnosticShapeOperation uint16, b509DiagnosticShapeDirection, b509DiagnosticShapeField string) B509DiagnosticKey {
	return B509DiagnosticKey{B509DiagnosticShapeFamily: "vaillant", B509DiagnosticShapeVariant: "B509", B509DiagnosticShapePrimary: b509DiagnosticShapeNumber(0xB5), B509DiagnosticShapeSecondary: b509DiagnosticShapeNumber(9), B509DiagnosticShapeOperation: b509DiagnosticShapeNumber(b509DiagnosticShapeOperation), B509DiagnosticShapeDirection: b509DiagnosticShapeDirection, B509DiagnosticShapeField: b509DiagnosticShapeField}
}

func b509DiagnosticShapeGet(b509DiagnosticShapeT *b509DiagnosticShapeTesting.T, b509DiagnosticShapeOperation uint16, b509DiagnosticShapeDirection, b509DiagnosticShapeName string) B509DiagnosticField {
	b509DiagnosticShapeT.Helper()
	b509DiagnosticShapeField, b509DiagnosticShapeOK, b509DiagnosticShapeError := B509LookupDiagnosticField(b509DiagnosticShapeKey(b509DiagnosticShapeOperation, b509DiagnosticShapeDirection, b509DiagnosticShapeName))
	if b509DiagnosticShapeError != nil || !b509DiagnosticShapeOK {
		b509DiagnosticShapeT.Fatalf("ab87 missing %s: %v", b509DiagnosticShapeName, b509DiagnosticShapeError)
	}
	return b509DiagnosticShapeField
}

func b509DiagnosticShapeExpected(b509DiagnosticShapeOperation uint16, b509DiagnosticShapeDirection, b509DiagnosticShapeName string, b509DiagnosticShapeWidth *int, b509DiagnosticShapeSign, b509DiagnosticShapeOrder, b509DiagnosticShapeConversion string, b509DiagnosticShapeMin, b509DiagnosticShapeMax *int64, b509DiagnosticShapeMapping string) B509DiagnosticField {
	return B509DiagnosticField{B509DiagnosticShapeKey: b509DiagnosticShapeKey(b509DiagnosticShapeOperation, b509DiagnosticShapeDirection, b509DiagnosticShapeName), B509DiagnosticShapeWidthBits: b509DiagnosticShapeWidth, B509DiagnosticShapeSignedness: b509DiagnosticShapeSign, B509DiagnosticShapeByteOrder: b509DiagnosticShapeOrder, B509DiagnosticShapeConversion: b509DiagnosticShapeConversion, B509DiagnosticShapeMinimum: b509DiagnosticShapeMin, B509DiagnosticShapeMaximum: b509DiagnosticShapeMax, B509DiagnosticShapeGrade: "Hypothesis", B509DiagnosticShapeMappingGrade: b509DiagnosticShapeMapping, B509DiagnosticShapeRawFallback: "preserve original bytes"}
}

func TestAb87ProductionContract(b509DiagnosticShapeT *b509DiagnosticShapeTesting.T) {
	b509DiagnosticShapeCatalog, b509DiagnosticShapeError := B509DiagnosticShapes()
	if b509DiagnosticShapeError != nil {
		b509DiagnosticShapeT.Fatal(b509DiagnosticShapeError)
	}
	b509DiagnosticShapeExpectedFields := []B509DiagnosticField{
		b509DiagnosticShapeExpected(0x13, "request", "arguments", b509DiagnosticShapeInt(0), "Unknown", "Unknown", "no_arguments", nil, nil, "Unknown"),
		b509DiagnosticShapeExpected(0x13, "response", "result", nil, "Unknown", "Unknown", "Boolean", nil, nil, "Unknown"),
		b509DiagnosticShapeExpected(0x34, "request", "address", b509DiagnosticShapeInt(32), "unsigned", "little_endian", "identity", b509DiagnosticShapeBound(0), b509DiagnosticShapeBound(4294967295), "Hypothesis"),
		b509DiagnosticShapeExpected(0x34, "request", "requested_count", b509DiagnosticShapeInt(8), "unsigned", "not_applicable", "identity", b509DiagnosticShapeBound(0), b509DiagnosticShapeBound(255), "Hypothesis"),
		b509DiagnosticShapeExpected(0x34, "response", "address", b509DiagnosticShapeInt(32), "unsigned", "little_endian", "identity", b509DiagnosticShapeBound(0), b509DiagnosticShapeBound(4294967295), "Hypothesis"),
		b509DiagnosticShapeExpected(0x34, "response", "data", nil, "Unknown", "Unknown", "raw_bytes", nil, nil, "Unknown"),
		b509DiagnosticShapeExpected(0x35, "request", "address", b509DiagnosticShapeInt(32), "unsigned", "little_endian", "identity", b509DiagnosticShapeBound(0), b509DiagnosticShapeBound(4294967295), "Hypothesis"),
		b509DiagnosticShapeExpected(0x35, "request", "data", nil, "Unknown", "Unknown", "raw_bytes", nil, nil, "Unknown"),
		b509DiagnosticShapeExpected(0x35, "response", "write_status", nil, "Unknown", "Unknown", "status_domain", b509DiagnosticShapeBound(0), b509DiagnosticShapeBound(6), "Unknown"),
		b509DiagnosticShapeExpected(0x51, "Unknown", "component", nil, "Unknown", "Unknown", "identity", b509DiagnosticShapeBound(0), b509DiagnosticShapeBound(7), "Hypothesis"),
		b509DiagnosticShapeExpected(0x51, "Unknown", "initiator", b509DiagnosticShapeInt(8), "unsigned", "not_applicable", "identity", b509DiagnosticShapeBound(0), b509DiagnosticShapeBound(255), "Hypothesis"),
		b509DiagnosticShapeExpected(0x51, "Unknown", "identity", b509DiagnosticShapeInt(40), "not_applicable", "not_applicable", "ASCII_exact", nil, nil, "Hypothesis"),
		b509DiagnosticShapeExpected(0x51, "Unknown", "software_component_1", nil, "Unknown", "Unknown", "decimal_scalar", nil, nil, "Unknown"),
		b509DiagnosticShapeExpected(0x51, "Unknown", "software_component_2", nil, "Unknown", "Unknown", "decimal_scalar", nil, nil, "Unknown"),
		b509DiagnosticShapeExpected(0x51, "Unknown", "software_component_3", nil, "Unknown", "Unknown", "decimal_scalar", nil, nil, "Unknown"),
		b509DiagnosticShapeExpected(0x51, "Unknown", "software_component_4", nil, "Unknown", "Unknown", "decimal_scalar", nil, nil, "Unknown"),
	}
	b509DiagnosticShapeExpectedFields[5].B509DiagnosticShapeRecommendedBytes = b509DiagnosticShapeInt(10)
	b509DiagnosticShapeExpectedFields[8].B509DiagnosticShapeEnumeration = map[string]*string{"0": nil, "1": nil, "2": nil, "3": nil, "4": nil, "5": nil, "6": nil}
	if !b509DiagnosticShapeReflect.DeepEqual(b509DiagnosticShapeCatalog.B509DiagnosticShapeFields, b509DiagnosticShapeExpectedFields) {
		b509DiagnosticShapeT.Fatal("ab87 finite metadata differs from independent contract")
	}
	if b509DiagnosticShapeCatalog.B509DiagnosticShapeRepresentation != "Proven" || b509DiagnosticShapeCatalog.B509DiagnosticShapeQualification != "Unknown" || b509DiagnosticShapeCatalog.B509DiagnosticShapeExecution != "Unknown" || b509DiagnosticShapeCatalog.B509DiagnosticShapePhysicalEffects != "Unknown" || b509DiagnosticShapeCatalog.B509DiagnosticShapeResponseCorrelation != "Unknown" || len(b509DiagnosticShapeCatalog.B509DiagnosticShapeLimitations) != 6 {
		b509DiagnosticShapeT.Fatal("ab87 evidence boundaries changed")
	}
	for _, b509DiagnosticShapeField := range b509DiagnosticShapeExpectedFields {
		b509DiagnosticShapeGot, b509DiagnosticShapeOK, b509DiagnosticShapeError := B509LookupDiagnosticField(b509DiagnosticShapeField.B509DiagnosticShapeKey)
		if b509DiagnosticShapeError != nil || !b509DiagnosticShapeOK || !b509DiagnosticShapeReflect.DeepEqual(b509DiagnosticShapeGot, b509DiagnosticShapeField) {
			b509DiagnosticShapeT.Fatal("ab87 exact lookup changed")
		}
	}
	if _, b509DiagnosticShapeOK, b509DiagnosticShapeError := B509LookupDiagnosticField(b509DiagnosticShapeKey(0x35, "request", "requested_count")); b509DiagnosticShapeError != nil || b509DiagnosticShapeOK {
		b509DiagnosticShapeT.Fatal("ab87 OP35 must not gain OP34 count")
	}
}

func b509DiagnosticShapeUnsigned(b509DiagnosticShapeField B509DiagnosticField, b509DiagnosticShapeRaw []byte) (uint64, bool) {
	if b509DiagnosticShapeField.B509DiagnosticShapeWidthBits == nil || *b509DiagnosticShapeField.B509DiagnosticShapeWidthBits != 8*len(b509DiagnosticShapeRaw) || b509DiagnosticShapeField.B509DiagnosticShapeSignedness != "unsigned" || b509DiagnosticShapeField.B509DiagnosticShapeConversion != "identity" {
		return 0, false
	}
	var b509DiagnosticShapeValue uint64
	switch b509DiagnosticShapeField.B509DiagnosticShapeByteOrder {
	case "little_endian":
		for b509DiagnosticShapeIndex, b509DiagnosticShapeByte := range b509DiagnosticShapeRaw {
			b509DiagnosticShapeValue |= uint64(b509DiagnosticShapeByte) << uint(8*b509DiagnosticShapeIndex)
		}
	case "not_applicable":
		if len(b509DiagnosticShapeRaw) != 1 {
			return 0, false
		}
		b509DiagnosticShapeValue = uint64(b509DiagnosticShapeRaw[0])
	default:
		return 0, false
	}
	return b509DiagnosticShapeValue, true
}

func b509DiagnosticShapeInBounds(b509DiagnosticShapeField B509DiagnosticField, b509DiagnosticShapeValue int64) bool {
	return (b509DiagnosticShapeField.B509DiagnosticShapeMinimum == nil || b509DiagnosticShapeValue >= *b509DiagnosticShapeField.B509DiagnosticShapeMinimum) && (b509DiagnosticShapeField.B509DiagnosticShapeMaximum == nil || b509DiagnosticShapeValue <= *b509DiagnosticShapeField.B509DiagnosticShapeMaximum)
}

func b509DiagnosticShapeFitsASCII(b509DiagnosticShapeField B509DiagnosticField, b509DiagnosticShapeValue string) bool {
	if b509DiagnosticShapeField.B509DiagnosticShapeConversion != "ASCII_exact" || b509DiagnosticShapeField.B509DiagnosticShapeWidthBits == nil || len(b509DiagnosticShapeValue)*8 != *b509DiagnosticShapeField.B509DiagnosticShapeWidthBits {
		return false
	}
	for _, b509DiagnosticShapeByte := range []byte(b509DiagnosticShapeValue) {
		if b509DiagnosticShapeByte > 127 {
			return false
		}
	}
	return true
}

func TestAb87SyntheticCases(b509DiagnosticShapeT *b509DiagnosticShapeTesting.T) {
	for _, b509DiagnosticShapeCase := range []struct {
		B509DiagnosticShapeOperation uint16
		B509DiagnosticShapeDirection string
	}{{0x34, "request"}, {0x34, "response"}, {0x35, "request"}} {
		b509DiagnosticShapeAddress := b509DiagnosticShapeGet(b509DiagnosticShapeT, b509DiagnosticShapeCase.B509DiagnosticShapeOperation, b509DiagnosticShapeCase.B509DiagnosticShapeDirection, "address")
		b509DiagnosticShapeValue, b509DiagnosticShapeOK := b509DiagnosticShapeUnsigned(b509DiagnosticShapeAddress, []byte{0x78, 0x56, 0x34, 0x12})
		if !b509DiagnosticShapeOK || b509DiagnosticShapeValue != 0x12345678 {
			b509DiagnosticShapeT.Fatal("ab87 u32 little endian address")
		}
	}
	b509DiagnosticShapeCount := b509DiagnosticShapeGet(b509DiagnosticShapeT, 0x34, "request", "requested_count")
	for _, b509DiagnosticShapeRaw := range []byte{0, 255} {
		b509DiagnosticShapeValue, b509DiagnosticShapeOK := b509DiagnosticShapeUnsigned(b509DiagnosticShapeCount, []byte{b509DiagnosticShapeRaw})
		if !b509DiagnosticShapeOK || b509DiagnosticShapeValue != uint64(b509DiagnosticShapeRaw) || !b509DiagnosticShapeInBounds(b509DiagnosticShapeCount, int64(b509DiagnosticShapeValue)) {
			b509DiagnosticShapeT.Fatal("ab87 raw count changed")
		}
	}
	b509DiagnosticShapeData := b509DiagnosticShapeGet(b509DiagnosticShapeT, 0x34, "response", "data")
	b509DiagnosticShapeSynthetic := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	if b509DiagnosticShapeData.B509DiagnosticShapeConversion != "raw_bytes" || b509DiagnosticShapeData.B509DiagnosticShapeRecommendedBytes == nil || *b509DiagnosticShapeData.B509DiagnosticShapeRecommendedBytes != 10 || b509DiagnosticShapeData.B509DiagnosticShapeMaximumBytes != nil || (b509DiagnosticShapeData.B509DiagnosticShapeWidthBits != nil && *b509DiagnosticShapeData.B509DiagnosticShapeWidthBits != 8*len(b509DiagnosticShapeSynthetic)) {
		b509DiagnosticShapeT.Fatal("ab87 recommendation became capacity")
	}
	b509DiagnosticShapeStatus := b509DiagnosticShapeGet(b509DiagnosticShapeT, 0x35, "response", "write_status")
	if !b509DiagnosticShapeInBounds(b509DiagnosticShapeStatus, 6) || b509DiagnosticShapeInBounds(b509DiagnosticShapeStatus, 7) {
		b509DiagnosticShapeT.Fatal("ab87 status domain")
	}
	if _, b509DiagnosticShapeOK := b509DiagnosticShapeStatus.B509DiagnosticShapeEnumeration["7"]; b509DiagnosticShapeOK {
		b509DiagnosticShapeT.Fatal("ab87 status 7 labeled")
	}
	b509DiagnosticShapeIdentity := b509DiagnosticShapeGet(b509DiagnosticShapeT, 0x51, "Unknown", "identity")
	if !b509DiagnosticShapeFitsASCII(b509DiagnosticShapeIdentity, "ABCDE") || b509DiagnosticShapeFitsASCII(b509DiagnosticShapeIdentity, "AB") || b509DiagnosticShapeFitsASCII(b509DiagnosticShapeIdentity, "ABCDEF") || b509DiagnosticShapeFitsASCII(b509DiagnosticShapeIdentity, "ABCDé") || b509DiagnosticShapeIdentity.B509DiagnosticShapePadding != nil {
		b509DiagnosticShapeT.Fatal("ab87 exact ASCII identity")
	}
	b509DiagnosticShapeComponent := b509DiagnosticShapeGet(b509DiagnosticShapeT, 0x51, "Unknown", "component")
	if !b509DiagnosticShapeInBounds(b509DiagnosticShapeComponent, 7) || b509DiagnosticShapeInBounds(b509DiagnosticShapeComponent, 8) {
		b509DiagnosticShapeT.Fatal("ab87 component range")
	}
	b509DiagnosticShapeScalars := []int64{1, 20, 3, 40}
	for b509DiagnosticShapeIndex, b509DiagnosticShapeName := range []string{"software_component_1", "software_component_2", "software_component_3", "software_component_4"} {
		b509DiagnosticShapeField := b509DiagnosticShapeGet(b509DiagnosticShapeT, 0x51, "Unknown", b509DiagnosticShapeName)
		if b509DiagnosticShapeField.B509DiagnosticShapeConversion != "decimal_scalar" || !b509DiagnosticShapeInBounds(b509DiagnosticShapeField, b509DiagnosticShapeScalars[b509DiagnosticShapeIndex]) || b509DiagnosticShapeField.B509DiagnosticShapeWidthBits != nil {
			b509DiagnosticShapeT.Fatal("ab87 four independent decimal scalars")
		}
	}
}

func TestAb87ExactIdentityAxes(b509DiagnosticShapeT *b509DiagnosticShapeTesting.T) {
	b509DiagnosticShapeControl := b509DiagnosticShapeGet(b509DiagnosticShapeT, 0x34, "request", "address")
	for _, b509DiagnosticShapeB524 := range []bool{false, true} {
		b509DiagnosticShapeRow := b509DiagnosticShapeControl
		if b509DiagnosticShapeB524 {
			b509DiagnosticShapeRow.B509DiagnosticShapeKey.B509DiagnosticShapeVariant = "B524"
			b509DiagnosticShapeRow.B509DiagnosticShapeKey.B509DiagnosticShapeSecondary = b509DiagnosticShapeNumber(0x24)
			b509DiagnosticShapeRow.B509DiagnosticShapeKey.B509DiagnosticShapeOpcode = b509DiagnosticShapeNumber(0)
			b509DiagnosticShapeRow.B509DiagnosticShapeKey.B509DiagnosticShapeGroup = b509DiagnosticShapeNumber(0)
			b509DiagnosticShapeRow.B509DiagnosticShapeKey.B509DiagnosticShapeInstance = b509DiagnosticShapeNumber(0)
			b509DiagnosticShapeRow.B509DiagnosticShapeKey.B509DiagnosticShapeRegister = b509DiagnosticShapeNumber(0)
		}
		b509DiagnosticShapeRows := []B509DiagnosticField{b509DiagnosticShapeRow}
		b509DiagnosticShapeType := b509DiagnosticShapeReflect.TypeOf(b509DiagnosticShapeRow.B509DiagnosticShapeKey)
		for b509DiagnosticShapeIndex := 0; b509DiagnosticShapeIndex < b509DiagnosticShapeType.NumField(); b509DiagnosticShapeIndex++ {
			b509DiagnosticShapeT.Run("ab87_"+b509DiagnosticShapeRow.B509DiagnosticShapeKey.B509DiagnosticShapeVariant+"_"+b509DiagnosticShapeType.Field(b509DiagnosticShapeIndex).Name, func(b509DiagnosticShapeT *b509DiagnosticShapeTesting.T) {
				if _, b509DiagnosticShapeOK, b509DiagnosticShapeError := b509DiagnosticShapeLookupField(b509DiagnosticShapeRows, b509DiagnosticShapeRow.B509DiagnosticShapeKey); b509DiagnosticShapeError != nil || !b509DiagnosticShapeOK {
					b509DiagnosticShapeT.Fatal("ab87 invalid control")
				}
				b509DiagnosticShapeChanged := b509DiagnosticShapeRow.B509DiagnosticShapeKey
				b509DiagnosticShapeAxis := b509DiagnosticShapeReflect.ValueOf(&b509DiagnosticShapeChanged).Elem().Field(b509DiagnosticShapeIndex)
				if b509DiagnosticShapeAxis.Kind() == b509DiagnosticShapeReflect.String {
					b509DiagnosticShapeAxis.SetString("ab87_distinct")
				} else if b509DiagnosticShapeAxis.IsNil() {
					b509DiagnosticShapeAxis.Set(b509DiagnosticShapeReflect.ValueOf(b509DiagnosticShapeNumber(0)))
				} else {
					b509DiagnosticShapeAxis.Set(b509DiagnosticShapeReflect.Zero(b509DiagnosticShapeAxis.Type()))
				}
				if _, b509DiagnosticShapeOK, b509DiagnosticShapeError := b509DiagnosticShapeLookupField(b509DiagnosticShapeRows, b509DiagnosticShapeChanged); b509DiagnosticShapeError != nil || b509DiagnosticShapeOK {
					b509DiagnosticShapeT.Fatal("ab87 axis aliased absent/zero")
				}
				if b509DiagnosticShapeAxis.Kind() == b509DiagnosticShapeReflect.Ptr {
					b509DiagnosticShapeAxis.Set(b509DiagnosticShapeReflect.ValueOf(b509DiagnosticShapeNumber(0)))
					if !b509DiagnosticShapeReflect.DeepEqual(b509DiagnosticShapeChanged, b509DiagnosticShapeRow.B509DiagnosticShapeKey) {
						if _, b509DiagnosticShapeOK, _ := b509DiagnosticShapeLookupField(b509DiagnosticShapeRows, b509DiagnosticShapeChanged); b509DiagnosticShapeOK {
							b509DiagnosticShapeT.Fatal("ab87 zero wildcard")
						}
					}
					b509DiagnosticShapeAxis.Set(b509DiagnosticShapeReflect.ValueOf(b509DiagnosticShapeNumber(65535)))
					if _, b509DiagnosticShapeOK, _ := b509DiagnosticShapeLookupField(b509DiagnosticShapeRows, b509DiagnosticShapeChanged); b509DiagnosticShapeOK {
						b509DiagnosticShapeT.Fatal("ab87 numeric axis ignored")
					}
				}
			})
		}
	}
}

func TestAb87DeepCopy(b509DiagnosticShapeT *b509DiagnosticShapeTesting.T) {
	b509DiagnosticShapeBefore, b509DiagnosticShapeError := B509DiagnosticShapes()
	if b509DiagnosticShapeError != nil {
		b509DiagnosticShapeT.Fatal(b509DiagnosticShapeError)
	}
	b509DiagnosticShapeCopy, _ := B509DiagnosticShapes()
	b509DiagnosticShapeCopy.B509DiagnosticShapeLimitations[0] = "ab87_changed"
	*b509DiagnosticShapeCopy.B509DiagnosticShapeFields[2].B509DiagnosticShapeKey.B509DiagnosticShapeOperation = 0
	*b509DiagnosticShapeCopy.B509DiagnosticShapeFields[2].B509DiagnosticShapeWidthBits = 16
	*b509DiagnosticShapeCopy.B509DiagnosticShapeFields[2].B509DiagnosticShapeMaximum = 1
	*b509DiagnosticShapeCopy.B509DiagnosticShapeFields[5].B509DiagnosticShapeRecommendedBytes = 1
	b509DiagnosticShapeCopy.B509DiagnosticShapeFields[8].B509DiagnosticShapeEnumeration["7"] = nil
	b509DiagnosticShapeCopy.B509DiagnosticShapeFields[0] = B509DiagnosticField{}
	b509DiagnosticShapeAfter, _ := B509DiagnosticShapes()
	if !b509DiagnosticShapeReflect.DeepEqual(b509DiagnosticShapeAfter, b509DiagnosticShapeBefore) {
		b509DiagnosticShapeT.Fatal("ab87 accessor shared backing")
	}
	b509DiagnosticShapeLookup := b509DiagnosticShapeGet(b509DiagnosticShapeT, 0x35, "response", "write_status")
	*b509DiagnosticShapeLookup.B509DiagnosticShapeKey.B509DiagnosticShapePrimary = 0
	*b509DiagnosticShapeLookup.B509DiagnosticShapeMaximum = 7
	b509DiagnosticShapeLookup.B509DiagnosticShapeEnumeration["7"] = nil
	if b509DiagnosticShapeReflect.DeepEqual(b509DiagnosticShapeGet(b509DiagnosticShapeT, 0x35, "response", "write_status"), b509DiagnosticShapeLookup) {
		b509DiagnosticShapeT.Fatal("ab87 lookup shared backing")
	}
	b509DiagnosticShapeSource := b509DiagnosticShapeBefore.B509DiagnosticShapeFields[8]
	b509DiagnosticShapeLabel := "ab87_label"
	b509DiagnosticShapeUnit := "ab87_unit"
	b509DiagnosticShapePadding := "ab87_padding"
	b509DiagnosticShapeSource.B509DiagnosticShapeEnumeration = map[string]*string{"6": &b509DiagnosticShapeLabel}
	b509DiagnosticShapeSource.B509DiagnosticShapeReplacement = []int64{1}
	b509DiagnosticShapeSource.B509DiagnosticShapeReserved = []int64{2}
	b509DiagnosticShapeSource.B509DiagnosticShapeLogicalPosition = b509DiagnosticShapeInt(0)
	b509DiagnosticShapeSource.B509DiagnosticShapeUnit = &b509DiagnosticShapeUnit
	b509DiagnosticShapeSource.B509DiagnosticShapeDefault = b509DiagnosticShapeBound(0)
	b509DiagnosticShapeSource.B509DiagnosticShapeMaximumBytes = b509DiagnosticShapeInt(11)
	b509DiagnosticShapeSource.B509DiagnosticShapePadding = &b509DiagnosticShapePadding
	b509DiagnosticShapeSnapshot, _ := b509DiagnosticShapeJSON.Marshal(b509DiagnosticShapeSource)
	b509DiagnosticShapeDetached, b509DiagnosticShapeOK, b509DiagnosticShapeError := b509DiagnosticShapeLookupField([]B509DiagnosticField{b509DiagnosticShapeSource}, b509DiagnosticShapeSource.B509DiagnosticShapeKey)
	if b509DiagnosticShapeError != nil || !b509DiagnosticShapeOK {
		b509DiagnosticShapeT.Fatal("ab87 synthetic copy control")
	}
	*b509DiagnosticShapeDetached.B509DiagnosticShapeEnumeration["6"] = "ab87_changed"
	b509DiagnosticShapeDetached.B509DiagnosticShapeReplacement[0] = 9
	b509DiagnosticShapeDetached.B509DiagnosticShapeReserved[0] = 9
	*b509DiagnosticShapeDetached.B509DiagnosticShapeLogicalPosition = 9
	*b509DiagnosticShapeDetached.B509DiagnosticShapeUnit = "ab87_changed"
	*b509DiagnosticShapeDetached.B509DiagnosticShapeDefault = 9
	*b509DiagnosticShapeDetached.B509DiagnosticShapeMaximumBytes = 9
	*b509DiagnosticShapeDetached.B509DiagnosticShapePadding = "ab87_changed"
	b509DiagnosticShapeAgain, _, _ := b509DiagnosticShapeLookupField([]B509DiagnosticField{b509DiagnosticShapeSource}, b509DiagnosticShapeSource.B509DiagnosticShapeKey)
	b509DiagnosticShapeActual, _ := b509DiagnosticShapeJSON.Marshal(b509DiagnosticShapeAgain)
	if string(b509DiagnosticShapeActual) != string(b509DiagnosticShapeSnapshot) {
		b509DiagnosticShapeT.Fatal("ab87 nested lookup storage shared")
	}
}
