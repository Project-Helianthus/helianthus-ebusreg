package vaillant

import (
	"encoding/binary"
	"encoding/json"
	"reflect"
	"testing"
)

func b555NumericUnitPointer[b555NumericUnitType any](b555NumericUnitValue b555NumericUnitType) *b555NumericUnitType {
	return &b555NumericUnitValue
}

func b555NumericUnitExpectedField(b555NumericUnitName string, b555NumericUnitWidth int, b555NumericUnitOrder, b555NumericUnitConversion string, b555NumericUnitRawMax uint64, b555NumericUnitNumericMax float64) B555NumericField {
	return B555NumericField{
		B555NumericUnitKey:       B555NumericUnitKey{B555NumericUnitFamily: "vaillant", B555NumericUnitVariant: "generic", B555NumericUnitPrimary: 0xB555, B555NumericUnitField: b555NumericUnitName},
		B555NumericUnitWidthBits: b555NumericUnitPointer(b555NumericUnitWidth), B555NumericUnitSignedness: "unsigned", B555NumericUnitByteOrder: b555NumericUnitOrder,
		B555NumericUnitConversion: b555NumericUnitConversion, B555NumericUnitRawMin: b555NumericUnitPointer(uint64(0)), B555NumericUnitRawMax: b555NumericUnitPointer(b555NumericUnitRawMax),
		B555NumericMin: b555NumericUnitPointer(float64(0)), B555NumericMax: b555NumericUnitPointer(b555NumericUnitNumericMax),
		B555NumericUnitEnumeration: map[uint64]string{}, B555NumericUnitReplacement: []B555NumericUnitReplacement{}, B555NumericUnitReserved: []uint64{},
		B555NumericUnitGrade: SchemaHypothesis, B555NumericUnitConversionGrade: SchemaHypothesis, B555NumericUnitUnitGrade: SchemaUnknown,
		B555NumericUnitRawFallback: "preserve original bytes; no implicit sentinel selection", B555NumericUnitEnforcement: SchemaUnknown, B555NumericUnitPeerEffects: SchemaUnknown,
	}
}

func b555NumericUnitExpected() []B555NumericField {
	b555NumericUnitNumeric := b555NumericUnitExpectedField("numeric", 16, "little_endian", "divide", 65535, 6553.5)
	b555NumericUnitNumeric.B555NumericUnitDivisor = b555NumericUnitPointer(float64(10))
	b555NumericUnitContext := b555NumericUnitNumeric
	b555NumericUnitContext.B555NumericUnitKey.B555NumericUnitVariant = "synthetic-contextual-replacement"
	b555NumericUnitContext.B555NumericUnitReplacement = []B555NumericUnitReplacement{{B555NumericUnitRaw: 65535, B555NumericUnitRole: "contextual replacement", B555NumericUnitContext: "synthetic-only qualified condition", B555NumericUnitGrade: SchemaHypothesis}}
	b555NumericUnitMin := b555NumericUnitExpectedField("minimum_setpoint", 8, "not_applicable", "identity", 255, 255)
	b555NumericUnitMin.B555NumericUnitKey.B555NumericUnitSecondary = b555NumericUnitPointer(uint8(0xA3))
	b555NumericUnitMax := b555NumericUnitMin
	b555NumericUnitMax.B555NumericUnitKey.B555NumericUnitField = "maximum_setpoint"
	b555NumericUnitCelsius := b555NumericUnitMin
	b555NumericUnitCelsius.B555NumericUnitKey.B555NumericUnitVariant = "synthetic-qualified-celsius"
	b555NumericUnitCelsius.B555NumericUnit = b555NumericUnitPointer("degrees C")
	b555NumericUnitCelsius.B555NumericUnitUnitGrade = SchemaHypothesis
	b555NumericUnitCelsius.B555NumericUnitReplacement = []B555NumericUnitReplacement{{B555NumericUnitRaw: 255, B555NumericUnitRole: "unavailable", B555NumericUnitContext: "synthetic-only qualified condition", B555NumericUnitGrade: SchemaHypothesis}}
	b555NumericUnitUnit := b555NumericUnitExpectedField("raw_numeric", 0, "Unknown", "Unknown", 0, 0)
	b555NumericUnitUnit.B555NumericUnitKey.B555NumericUnitVariant = "unqualified-unit-2"
	b555NumericUnitUnit.B555NumericUnitWidthBits = nil
	b555NumericUnitUnit.B555NumericUnitSignedness = "Unknown"
	b555NumericUnitUnit.B555NumericUnitCode = b555NumericUnitPointer(uint8(2))
	b555NumericUnitUnit.B555NumericUnitRawMin = nil
	b555NumericUnitUnit.B555NumericUnitRawMax = nil
	b555NumericUnitUnit.B555NumericMin = nil
	b555NumericUnitUnit.B555NumericMax = nil
	b555NumericUnitUnit.B555NumericUnitGrade = SchemaUnknown
	b555NumericUnitUnit.B555NumericUnitConversionGrade = SchemaUnknown
	return []B555NumericField{b555NumericUnitNumeric, b555NumericUnitContext, b555NumericUnitMin, b555NumericUnitMax, b555NumericUnitCelsius, b555NumericUnitUnit}
}

func TestAb94ProductionMetadata(b555NumericUnitT *testing.T) {
	b555NumericUnitGot, b555NumericUnitError := B555NumericUnits()
	if b555NumericUnitError != nil {
		b555NumericUnitT.Fatal(b555NumericUnitError)
	}
	b555NumericUnitWant := B555NumericMetadata{
		B555NumericUnitRepresentationGrade: SchemaProven, B555NumericUnitPhysicalQualification: SchemaUnknown, B555NumericUnitNativeSupport: SchemaUnknown,
		B555NumericUnitFreshness: SchemaUnknown, B555NumericUnitPersistence: SchemaUnknown, B555NumericUnitExecution: SchemaUnknown,
		B555NumericUnitSourceContract: "AI23 supplied functional facts; supplemental declarations only",
		B555NumericUnitLimitations:    []string{"Logical positions are unknown; no physical layout is asserted.", "Synthetic contextual variants are examples, not identified device variants.", "No invocation, transport, group/day mapping, persistence or admission is enabled.", "Unit-code 2 has raw numeric material but no qualified semantic conversion."},
		B555NumericUnitEntries:        b555NumericUnitExpected(),
	}
	if !reflect.DeepEqual(b555NumericUnitGot, b555NumericUnitWant) {
		b555NumericUnitT.Fatalf("ab94 production metadata differs: got %+v", b555NumericUnitGot)
	}
	for _, b555NumericUnitExpected := range b555NumericUnitWant.B555NumericUnitEntries {
		b555NumericUnitGotField, b555NumericUnitFound, b555NumericUnitError := B555LookupNumericUnits(b555NumericUnitExpected.B555NumericUnitKey)
		if b555NumericUnitError != nil || !b555NumericUnitFound || !reflect.DeepEqual(b555NumericUnitGotField, b555NumericUnitExpected) {
			b555NumericUnitT.Fatalf("ab94 exact lookup: %v %v %+v", b555NumericUnitError, b555NumericUnitFound, b555NumericUnitGotField)
		}
	}
}

func b555NumericUnitEvaluate(b555NumericUnitField B555NumericField, b555NumericUnitRaw []byte) (float64, bool) {
	if b555NumericUnitField.B555NumericUnitWidthBits == nil || b555NumericUnitField.B555NumericUnitRawMin == nil || b555NumericUnitField.B555NumericUnitRawMax == nil {
		return 0, false
	}
	var b555NumericUnitNumber uint64
	switch {
	case *b555NumericUnitField.B555NumericUnitWidthBits == 16 && b555NumericUnitField.B555NumericUnitByteOrder == "little_endian" && len(b555NumericUnitRaw) == 2:
		b555NumericUnitNumber = uint64(binary.LittleEndian.Uint16(b555NumericUnitRaw))
	case *b555NumericUnitField.B555NumericUnitWidthBits == 8 && b555NumericUnitField.B555NumericUnitByteOrder == "not_applicable" && len(b555NumericUnitRaw) == 1:
		b555NumericUnitNumber = uint64(b555NumericUnitRaw[0])
	default:
		return 0, false
	}
	if b555NumericUnitField.B555NumericUnitSignedness != "unsigned" || b555NumericUnitNumber < *b555NumericUnitField.B555NumericUnitRawMin || b555NumericUnitNumber > *b555NumericUnitField.B555NumericUnitRawMax {
		return 0, false
	}
	b555NumericUnitValue := float64(b555NumericUnitNumber)
	switch b555NumericUnitField.B555NumericUnitConversion {
	case "divide":
		if b555NumericUnitField.B555NumericUnitDivisor == nil || *b555NumericUnitField.B555NumericUnitDivisor <= 0 {
			return 0, false
		}
		b555NumericUnitValue /= *b555NumericUnitField.B555NumericUnitDivisor
	case "identity":
	default:
		return 0, false
	}
	if b555NumericUnitField.B555NumericMin == nil || b555NumericUnitField.B555NumericMax == nil || b555NumericUnitValue < *b555NumericUnitField.B555NumericMin || b555NumericUnitValue > *b555NumericUnitField.B555NumericMax {
		return 0, false
	}
	return b555NumericUnitValue, true
}

func TestAb94SyntheticNumbersAndContexts(b555NumericUnitT *testing.T) {
	b555NumericUnitMetadata, b555NumericUnitError := B555NumericUnits()
	if b555NumericUnitError != nil {
		b555NumericUnitT.Fatal(b555NumericUnitError)
	}
	b555NumericUnitRows := b555NumericUnitMetadata.B555NumericUnitEntries
	for b555NumericUnitIndex, b555NumericUnitRaw := range [][]byte{{0, 0}, {1, 0}, {255, 255}} {
		b555NumericUnitValue, b555NumericUnitOK := b555NumericUnitEvaluate(b555NumericUnitRows[0], b555NumericUnitRaw)
		if !b555NumericUnitOK || b555NumericUnitValue != []float64{0, 0.1, 6553.5}[b555NumericUnitIndex] {
			b555NumericUnitT.Fatalf("ab94 numeric %v: %v %v", b555NumericUnitRaw, b555NumericUnitValue, b555NumericUnitOK)
		}
	}
	for _, b555NumericUnitIndex := range []int{0, 1} {
		b555NumericUnitValue, b555NumericUnitOK := b555NumericUnitEvaluate(b555NumericUnitRows[b555NumericUnitIndex], []byte{255, 255})
		if !b555NumericUnitOK || b555NumericUnitValue != 6553.5 {
			b555NumericUnitT.Fatal("ab94 context erased formal numeric alternative")
		}
	}
	if len(b555NumericUnitRows[0].B555NumericUnitReplacement) != 0 || len(b555NumericUnitRows[1].B555NumericUnitReplacement) != 1 || b555NumericUnitRows[1].B555NumericUnitReplacement[0].B555NumericUnitRaw != 65535 {
		b555NumericUnitT.Fatal("ab94 replacement branch conflated")
	}
	for _, b555NumericUnitIndex := range []int{2, 3} {
		b555NumericUnitValue, b555NumericUnitOK := b555NumericUnitEvaluate(b555NumericUnitRows[b555NumericUnitIndex], []byte{255})
		if !b555NumericUnitOK || b555NumericUnitValue != 255 || b555NumericUnitRows[b555NumericUnitIndex].B555NumericUnit != nil || len(b555NumericUnitRows[b555NumericUnitIndex].B555NumericUnitReplacement) != 0 {
			b555NumericUnitT.Fatal("ab94 generic A3 acquired qualified semantics")
		}
	}
	if b555NumericUnitRows[4].B555NumericUnit == nil || *b555NumericUnitRows[4].B555NumericUnit != "degrees C" || b555NumericUnitRows[4].B555NumericUnitReplacement[0].B555NumericUnitRaw != 255 {
		b555NumericUnitT.Fatal("ab94 synthetic qualification missing")
	}
	b555NumericUnitRaw := []byte{1, 0}
	b555NumericUnitBefore := append([]byte(nil), b555NumericUnitRaw...)
	if _, b555NumericUnitOK := b555NumericUnitEvaluate(b555NumericUnitRows[5], b555NumericUnitRaw); b555NumericUnitOK {
		b555NumericUnitT.Fatal("ab94 unit 2 inherited conversion")
	}
	if !reflect.DeepEqual(b555NumericUnitRaw, b555NumericUnitBefore) || b555NumericUnitRows[5].B555NumericUnit != nil || b555NumericUnitRows[5].B555NumericUnitConversionGrade != SchemaUnknown || *b555NumericUnitRows[5].B555NumericUnitCode != 2 {
		b555NumericUnitT.Fatal("ab94 unit 2 raw/unknown lost")
	}
}

func TestAb94ExactIdentityAxes(b555NumericUnitT *testing.T) {
	b555NumericUnitMetadata, b555NumericUnitError := B555NumericUnits()
	if b555NumericUnitError != nil {
		b555NumericUnitT.Fatal(b555NumericUnitError)
	}
	for _, b555NumericUnitRow := range b555NumericUnitMetadata.B555NumericUnitEntries {
		for b555NumericUnitIndex := 0; b555NumericUnitIndex < reflect.TypeOf(b555NumericUnitRow.B555NumericUnitKey).NumField(); b555NumericUnitIndex++ {
			b555NumericUnitName := reflect.TypeOf(b555NumericUnitRow.B555NumericUnitKey).Field(b555NumericUnitIndex).Name
			b555NumericUnitT.Run(b555NumericUnitRow.B555NumericUnitKey.B555NumericUnitVariant+"/"+b555NumericUnitRow.B555NumericUnitKey.B555NumericUnitField+"/"+b555NumericUnitName, func(b555NumericUnitT *testing.T) {
				b555NumericUnitKey := b555NumericUnitRow.B555NumericUnitKey
				b555NumericUnitAxis := reflect.ValueOf(&b555NumericUnitKey).Elem().Field(b555NumericUnitIndex)
				switch b555NumericUnitAxis.Kind() {
				case reflect.String:
					b555NumericUnitAxis.SetString("ab94_distinct")
				case reflect.Uint16:
					b555NumericUnitAxis.SetUint(0)
				case reflect.Ptr:
					b555NumericUnitAxis.Set(reflect.New(b555NumericUnitAxis.Type().Elem()))
				}
				if reflect.DeepEqual(b555NumericUnitKey, b555NumericUnitRow.B555NumericUnitKey) {
					b555NumericUnitT.Fatal("ab94 axis not varied")
				}
				if _, b555NumericUnitFound, b555NumericUnitError := B555LookupNumericUnits(b555NumericUnitKey); b555NumericUnitError != nil || b555NumericUnitFound {
					b555NumericUnitT.Fatalf("ab94 axis ignored: %v %v", b555NumericUnitFound, b555NumericUnitError)
				}
				if _, b555NumericUnitFound, b555NumericUnitError := B555LookupNumericUnits(b555NumericUnitRow.B555NumericUnitKey); b555NumericUnitError != nil || !b555NumericUnitFound {
					b555NumericUnitT.Fatal("ab94 valid control failed")
				}
			})
		}
	}
	// Synthetic B524 coordinates exercise the same matcher without declaring a B524 protocol row.
	b555NumericUnitKey := B555NumericUnitKey{B555NumericUnitFamily: "ab94_family", B555NumericUnitVariant: "ab94_variant", B555NumericUnitPrimary: 0xB524, B555NumericUnitSecondary: b555NumericUnitPointer(uint8(0)), B555NumericUnitOperation: b555NumericUnitPointer(uint8(0)), B555NumericUnitDirection: b555NumericUnitPointer("response"), B555NumericUnitField: "ab94_field", B555NumericUnitOpcode: b555NumericUnitPointer(uint8(0)), B555NumericUnitGroup: b555NumericUnitPointer(uint8(0)), B555NumericUnitInstance: b555NumericUnitPointer(uint8(0)), B555NumericUnitRegister: b555NumericUnitPointer(uint16(0))}
	for b555NumericUnitIndex := 3; b555NumericUnitIndex < reflect.TypeOf(b555NumericUnitKey).NumField(); b555NumericUnitIndex++ {
		if reflect.ValueOf(b555NumericUnitKey).Field(b555NumericUnitIndex).Kind() != reflect.Ptr {
			continue
		}
		for _, b555NumericUnitAbsent := range []bool{false, true} {
			b555NumericUnitChanged := b555NumericUnitKey
			b555NumericUnitAxis := reflect.ValueOf(&b555NumericUnitChanged).Elem().Field(b555NumericUnitIndex)
			if b555NumericUnitAbsent {
				b555NumericUnitAxis.SetZero()
			} else {
				b555NumericUnitAxis.Set(reflect.New(b555NumericUnitAxis.Type().Elem()))
				if b555NumericUnitAxis.Elem().Kind() == reflect.String {
					b555NumericUnitAxis.Elem().SetString("ab94_other")
				} else {
					b555NumericUnitAxis.Elem().SetUint(1)
				}
			}
			b555NumericUnitRows := []B555NumericField{{B555NumericUnitKey: b555NumericUnitKey}}
			if _, b555NumericUnitFound := b555NumericUnitLookup(b555NumericUnitRows, b555NumericUnitChanged); b555NumericUnitFound {
				b555NumericUnitT.Fatalf("ab94 B524 axis %s ignored (absent=%v)", reflect.TypeOf(b555NumericUnitKey).Field(b555NumericUnitIndex).Name, b555NumericUnitAbsent)
			}
			if _, b555NumericUnitFound := b555NumericUnitLookup(b555NumericUnitRows, b555NumericUnitKey); !b555NumericUnitFound {
				b555NumericUnitT.Fatal("ab94 B524 control failed")
			}
		}
	}
}

func b555NumericUnitMutate(b555NumericUnitValue reflect.Value) {
	switch b555NumericUnitValue.Kind() {
	case reflect.Struct:
		for b555NumericUnitIndex := 0; b555NumericUnitIndex < b555NumericUnitValue.NumField(); b555NumericUnitIndex++ {
			b555NumericUnitMutate(b555NumericUnitValue.Field(b555NumericUnitIndex))
		}
	case reflect.Ptr:
		if !b555NumericUnitValue.IsNil() {
			b555NumericUnitMutate(b555NumericUnitValue.Elem())
		}
	case reflect.Slice:
		for b555NumericUnitIndex := 0; b555NumericUnitIndex < b555NumericUnitValue.Len(); b555NumericUnitIndex++ {
			b555NumericUnitMutate(b555NumericUnitValue.Index(b555NumericUnitIndex))
		}
		if b555NumericUnitValue.Type().Elem().Kind() == reflect.Uint64 {
			b555NumericUnitValue.Set(reflect.Append(b555NumericUnitValue, reflect.ValueOf(uint64(7))))
		}
	case reflect.Map:
		b555NumericUnitValue.SetMapIndex(reflect.ValueOf(uint64(7)), reflect.ValueOf("ab94_changed"))
	case reflect.String:
		b555NumericUnitValue.SetString("ab94_changed")
	case reflect.Uint8, reflect.Uint16, reflect.Uint64:
		b555NumericUnitValue.SetUint(7)
	case reflect.Int:
		b555NumericUnitValue.SetInt(7)
	case reflect.Float64:
		b555NumericUnitValue.SetFloat(7)
	case reflect.Bool:
		b555NumericUnitValue.SetBool(true)
	}
}

func TestAb94DeepCopyIsolation(b555NumericUnitT *testing.T) {
	b555NumericUnitBefore, b555NumericUnitError := B555NumericUnits()
	if b555NumericUnitError != nil {
		b555NumericUnitT.Fatal(b555NumericUnitError)
	}
	b555NumericUnitSnapshot, b555NumericUnitError := json.Marshal(b555NumericUnitBefore)
	if b555NumericUnitError != nil {
		b555NumericUnitT.Fatal(b555NumericUnitError)
	}
	b555NumericUnitReturned, b555NumericUnitError := B555NumericUnits()
	if b555NumericUnitError != nil {
		b555NumericUnitT.Fatal(b555NumericUnitError)
	}
	b555NumericUnitMutate(reflect.ValueOf(&b555NumericUnitReturned).Elem())
	for _, b555NumericUnitRow := range b555NumericUnitBefore.B555NumericUnitEntries {
		b555NumericUnitLookupResult, b555NumericUnitFound, b555NumericUnitError := B555LookupNumericUnits(b555NumericUnitRow.B555NumericUnitKey)
		if b555NumericUnitError != nil || !b555NumericUnitFound {
			b555NumericUnitT.Fatal("ab94 lookup failed")
		}
		b555NumericUnitMutate(reflect.ValueOf(&b555NumericUnitLookupResult).Elem())
		b555NumericUnitAgain, b555NumericUnitFound, b555NumericUnitError := B555LookupNumericUnits(b555NumericUnitRow.B555NumericUnitKey)
		if b555NumericUnitError != nil || !b555NumericUnitFound || !reflect.DeepEqual(b555NumericUnitAgain, b555NumericUnitRow) {
			b555NumericUnitT.Fatal("ab94 lookup shared nested storage")
		}
	}
	b555NumericUnitAfter, b555NumericUnitError := B555NumericUnits()
	if b555NumericUnitError != nil {
		b555NumericUnitT.Fatal(b555NumericUnitError)
	}
	b555NumericUnitAfterBytes, b555NumericUnitError := json.Marshal(b555NumericUnitAfter)
	if b555NumericUnitError != nil {
		b555NumericUnitT.Fatal(b555NumericUnitError)
	}
	if string(b555NumericUnitSnapshot) != string(b555NumericUnitAfterBytes) {
		b555NumericUnitT.Fatal("ab94 accessor shared nested storage")
	}
}
