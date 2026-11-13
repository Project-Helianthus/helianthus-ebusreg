package vaillant

import (
	b511CandidateReflect "reflect"
	b511CandidateTesting "testing"
)

func b511CandidateByte(b511CandidateValue uint8) *uint8 { return &b511CandidateValue }

func b511CandidateEntry(b511CandidateT *b511CandidateTesting.T, b511CandidateVariant string, b511CandidateOperation *uint8, b511CandidateDirection, b511CandidateField string) B511CandidateEntry {
	b511CandidateT.Helper()
	b511CandidateResult, b511CandidateFound, b511CandidateError := B511CandidateLookup(B511CandidateB511Key{
		B511CandidateFamily: "vaillant", B511CandidateVariant: b511CandidateVariant, B511CandidateOperation: b511CandidateOperation,
		B511CandidateDirection: b511CandidateDirection, B511CandidateField: b511CandidateField,
	})
	if b511CandidateError != nil || !b511CandidateFound {
		b511CandidateT.Fatalf("lookup %s/%s: found=%v error=%v", b511CandidateVariant, b511CandidateField, b511CandidateFound, b511CandidateError)
	}
	return b511CandidateResult
}

func TestAb67ExactKeyAxes(b511CandidateT *b511CandidateTesting.T) {
	b511CandidateControl := B511CandidateB511Key{B511CandidateFamily: "vaillant", B511CandidateVariant: "OP0C", B511CandidateOperation: b511CandidateByte(12), B511CandidateDirection: "request", B511CandidateField: "mode"}
	b511CandidateChanges := map[string]func(*B511CandidateB511Key){
		"family":  func(b511CandidateKey *B511CandidateB511Key) { b511CandidateKey.B511CandidateFamily = "other-family" },
		"variant": func(b511CandidateKey *B511CandidateB511Key) { b511CandidateKey.B511CandidateVariant = "OP04" },
		"primary-absent-zero": func(b511CandidateKey *B511CandidateB511Key) {
			b511CandidateKey.B511CandidatePrimary = b511CandidateByte(0)
		},
		"secondary-absent-zero": func(b511CandidateKey *B511CandidateB511Key) {
			b511CandidateKey.B511CandidateSecondary = b511CandidateByte(0)
		},
		"operation": func(b511CandidateKey *B511CandidateB511Key) {
			b511CandidateKey.B511CandidateOperation = b511CandidateByte(4)
		},
		"operation-absent": func(b511CandidateKey *B511CandidateB511Key) { b511CandidateKey.B511CandidateOperation = nil },
		"direction":        func(b511CandidateKey *B511CandidateB511Key) { b511CandidateKey.B511CandidateDirection = "response" },
		"field":            func(b511CandidateKey *B511CandidateB511Key) { b511CandidateKey.B511CandidateField = "flags" },
	}
	for b511CandidateName, b511CandidateChange := range b511CandidateChanges {
		b511CandidateT.Run(b511CandidateName, func(b511CandidateT *b511CandidateTesting.T) {
			if b511CandidateValue, b511CandidateFound, b511CandidateError := B511CandidateLookup(b511CandidateControl); b511CandidateError != nil || !b511CandidateFound || b511CandidateValue.B511CandidateValues[2] != "request-enter" {
				b511CandidateT.Fatalf("valid independent control failed: %v %v", b511CandidateFound, b511CandidateError)
			}
			b511CandidateKey := b511CandidateControl
			b511CandidateChange(&b511CandidateKey)
			if _, b511CandidateFound, b511CandidateError := B511CandidateLookup(b511CandidateKey); b511CandidateFound || b511CandidateError != nil {
				b511CandidateT.Fatalf("changed axis matched: %v %v", b511CandidateFound, b511CandidateError)
			}
		})
	}
	for b511CandidateAxis := 0; b511CandidateAxis < 4; b511CandidateAxis++ {
		b511CandidateUnknown := b511CandidateControl
		switch b511CandidateAxis {
		case 0:
			b511CandidateUnknown.B511CandidateFamily = ""
		case 1:
			b511CandidateUnknown.B511CandidateVariant = ""
		case 2:
			b511CandidateUnknown.B511CandidateDirection = ""
		case 3:
			b511CandidateUnknown.B511CandidateField = ""
		}
		if _, b511CandidateFound, b511CandidateError := B511CandidateLookup(b511CandidateUnknown); b511CandidateFound || b511CandidateError != nil {
			b511CandidateT.Fatal("unknown string axis wildcard-matched")
		}
	}

	b511CandidateZero := b511CandidateEntry(b511CandidateT, "OP00-reserved", b511CandidateByte(0), "response", "status_byte_one").B511CandidateKey
	b511CandidateZero.B511CandidateOperation = nil
	if _, b511CandidateFound, b511CandidateError := B511CandidateLookup(b511CandidateZero); b511CandidateFound || b511CandidateError != nil {
		b511CandidateT.Fatal("absent operation matched zero")
	}
	b511CandidatePartial := B511CandidateB511Key{B511CandidateFamily: "vaillant", B511CandidateVariant: "10", B511CandidatePrimary: b511CandidateByte(0xB5), B511CandidateDirection: "request", B511CandidateField: "prefix"}
	if _, b511CandidateFound, b511CandidateError := B511CandidateLookup(b511CandidatePartial); !b511CandidateFound || b511CandidateError != nil {
		b511CandidateT.Fatal("known primary control failed")
	}
	for _, b511CandidatePrimary := range []*uint8{nil, b511CandidateByte(0), b511CandidateByte(0xF5)} {
		b511CandidateChanged := b511CandidatePartial
		b511CandidateChanged.B511CandidatePrimary = b511CandidatePrimary
		if _, b511CandidateFound, b511CandidateError := B511CandidateLookup(b511CandidateChanged); b511CandidateFound || b511CandidateError != nil {
			b511CandidateT.Fatal("changed known primary matched")
		}
	}
}

func TestAb67Requests(b511CandidateT *b511CandidateTesting.T) {
	b511CandidateMode := b511CandidateEntry(b511CandidateT, "OP0C", b511CandidateByte(12), "request", "mode")
	if *b511CandidateMode.B511CandidatePosition != 3 || !b511CandidateReflect.DeepEqual(b511CandidateMode.B511CandidateValues, map[uint8]string{0: "not-in-standby", 1: "in-standby", 2: "request-enter", 3: "request-exit"}) || b511CandidateMode.B511CandidateDefault != nil {
		b511CandidateT.Fatal("mode or absent argument changed")
	}
	for _, b511CandidateOperation := range []uint8{7, 9} {
		b511CandidateVariant := map[uint8]string{7: "OP07", 9: "OP09"}[b511CandidateOperation]
		b511CandidateFlags := b511CandidateEntry(b511CandidateT, b511CandidateVariant, b511CandidateByte(b511CandidateOperation), "request", "flags")
		if *b511CandidateFlags.B511CandidatePosition != 3 || !b511CandidateReflect.DeepEqual(b511CandidateFlags.B511CandidateBits, map[uint8]string{0: "assistant", 1: "factory-initialization"}) || !b511CandidateReflect.DeepEqual(b511CandidateFlags.B511CandidateReservedBits, []uint8{2, 3, 4, 5, 6, 7}) {
			b511CandidateT.Fatal("flags declarations changed")
		}
		b511CandidateRaw := uint8(3)
		b511CandidateNamed := make(map[string]bool)
		for b511CandidateBit, b511CandidateLabel := range b511CandidateFlags.B511CandidateBits {
			b511CandidateNamed[b511CandidateLabel] = b511CandidateRaw&(1<<b511CandidateBit) != 0
		}
		if !b511CandidateReflect.DeepEqual(b511CandidateNamed, map[string]bool{"assistant": true, "factory-initialization": true}) {
			b511CandidateT.Fatal("raw 03 lost a named flag")
		}
	}
	b511CandidateAssistant := b511CandidateEntry(b511CandidateT, "OP0B", b511CandidateByte(11), "request", "flags")
	if *b511CandidateAssistant.B511CandidatePosition != 3 || !b511CandidateReflect.DeepEqual(b511CandidateAssistant.B511CandidateBits, map[uint8]string{0: "assistant"}) {
		b511CandidateT.Fatal("OP0B inherited factory initialization")
	}
	b511CandidateBoolean := b511CandidateEntry(b511CandidateT, "OP04", b511CandidateByte(4), "request", "boolean")
	if *b511CandidateBoolean.B511CandidatePosition != 3 || !b511CandidateReflect.DeepEqual(b511CandidateBoolean.B511CandidateValues, map[uint8]string{0: "false", 1: "true"}) {
		b511CandidateT.Fatal("boolean changed")
	}
	for b511CandidateIndex, b511CandidateField := range []string{"temperature", "actuator", "setpoint", "configuration"} {
		b511CandidateOperation := []uint8{10, 10, 13, 13}[b511CandidateIndex]
		b511CandidateVariant := []string{"OP0A", "OP0A", "OP0D", "OP0D"}[b511CandidateIndex]
		if *b511CandidateEntry(b511CandidateT, b511CandidateVariant, b511CandidateByte(b511CandidateOperation), "request", b511CandidateField).B511CandidatePosition != []int{3, 5, 3, 4}[b511CandidateIndex] {
			b511CandidateT.Fatal("request coordinate changed")
		}
	}
}

func TestAb67SpecialFunctions(b511CandidateT *b511CandidateTesting.T) {
	for b511CandidateIndex, b511CandidateVariant := range []string{"OP0F-first", "OP0F-second", "OP0F-third"} {
		b511CandidateExpected := [][]uint8{{0, 1, 3, 4, 5, 6, 7, 8, 10}, {0, 1, 4, 7, 10, 11}, {0, 1, 4, 7, 10}}[b511CandidateIndex]
		b511CandidateSpecial := b511CandidateEntry(b511CandidateT, b511CandidateVariant, b511CandidateByte(15), "request", "special_function")
		if !b511CandidateReflect.DeepEqual(b511CandidateSpecial.B511CandidateMembership, b511CandidateExpected) || b511CandidateSpecial.B511CandidateValues != nil {
			b511CandidateT.Fatal("membership merged or labels invented")
		}
		for b511CandidateIndex, b511CandidateField := range []string{"binding", "special_function", "heating", "error_maintenance_cooling_allowed", "wish_temperature", "quick_veto_hour", "quick_veto_minute", "event_end_hour", "event_end_minute"} {
			if *b511CandidateEntry(b511CandidateT, b511CandidateVariant, b511CandidateByte(15), "request", b511CandidateField).B511CandidatePosition != b511CandidateIndex+3 {
				b511CandidateT.Fatal("OP0F coordinate changed")
			}
		}
		b511CandidateWish := b511CandidateEntry(b511CandidateT, b511CandidateVariant, b511CandidateByte(15), "request", "wish_temperature")
		if *b511CandidateWish.B511CandidateWidthBits != 8 || float64(1)**b511CandidateWish.B511CandidateScale != 0.5 || *b511CandidateWish.B511CandidateUnit != "degrees" {
			b511CandidateT.Fatal("wish raw 01 conversion changed")
		}
		for _, b511CandidateField := range []string{"quick_veto_hour", "quick_veto_minute", "event_end_hour", "event_end_minute"} {
			b511CandidateClock := b511CandidateEntry(b511CandidateT, b511CandidateVariant, b511CandidateByte(15), "request", b511CandidateField)
			if !b511CandidateReflect.DeepEqual(b511CandidateClock.B511CandidateRawDomain, []int{0, 255}) || b511CandidateClock.B511CandidateClockValidation != SchemaUnknown || b511CandidateClock.B511CandidateScale != nil || b511CandidateClock.B511CandidateUnit != nil {
				b511CandidateT.Fatal("raw 255 promoted to clock")
			}
		}
	}
}

func TestAb67ResponseAlternatives(b511CandidateT *b511CandidateTesting.T) {
	b511CandidateReserved := b511CandidateEntry(b511CandidateT, "OP00-reserved", b511CandidateByte(0), "response", "status_byte_one")
	b511CandidateFan := b511CandidateEntry(b511CandidateT, "OP00-fan", b511CandidateByte(0), "response", "status_byte_one")
	if len(b511CandidateReserved.B511CandidateBits) != 0 || !b511CandidateReflect.DeepEqual(b511CandidateReserved.B511CandidateReservedBits, []uint8{7}) || !b511CandidateReflect.DeepEqual(b511CandidateFan.B511CandidateBits, map[uint8]string{7: "fan-running"}) || len(b511CandidateFan.B511CandidateReservedBits) != 0 {
		b511CandidateT.Fatal("raw 80 alternatives collapsed")
	}
	for _, b511CandidateStatus := range []B511CandidateEntry{b511CandidateReserved, b511CandidateFan} {
		if *b511CandidateStatus.B511CandidatePosition != 8 || *b511CandidateStatus.B511CandidateExtent != 12 {
			b511CandidateT.Fatal("OP00 abstract positions changed")
		}
	}
	for b511CandidateIndex, b511CandidateVariant := range []string{"OP07-base", "OP07-extended"} {
		for _, b511CandidateField := range []string{"room_setpoint", "overlapping_structure_unknown"} {
			b511CandidateOverlap := b511CandidateEntry(b511CandidateT, b511CandidateVariant, b511CandidateByte(7), "response", b511CandidateField)
			if *b511CandidateOverlap.B511CandidatePosition != 7 || *b511CandidateOverlap.B511CandidateExtent != 8+b511CandidateIndex {
				b511CandidateT.Fatal("overlap concatenated")
			}
		}
		b511CandidateRoom := b511CandidateEntry(b511CandidateT, b511CandidateVariant, b511CandidateByte(7), "response", "room_setpoint")
		if *b511CandidateRoom.B511CandidateWidthBits != 8 || *b511CandidateRoom.B511CandidateScale != 0.5 || *b511CandidateRoom.B511CandidateUnit != "degrees" {
			b511CandidateT.Fatal("room setpoint codec changed")
		}
	}
	b511CandidatePressure := b511CandidateEntry(b511CandidateT, "OP07-extended", b511CandidateByte(7), "response", "building_circuit_pressure")
	if *b511CandidatePressure.B511CandidatePosition != 8 || *b511CandidatePressure.B511CandidateExtent != 9 || *b511CandidatePressure.B511CandidateUnit != "percent" || !b511CandidateReflect.DeepEqual(b511CandidatePressure.B511CandidateRawDomain, []int{0, 100}) {
		b511CandidateT.Fatal("pressure promoted to bar or extent summed")
	}
}

func TestAb67PartialAndAdmission(b511CandidateT *b511CandidateTesting.T) {
	b511CandidateCatalog, b511CandidateError := B511CandidateCatalog()
	if b511CandidateError != nil {
		b511CandidateT.Fatal(b511CandidateError)
	}
	if len(b511CandidateCatalog.B511CandidateEntries) != 47 || b511CandidateCatalog.B511CandidateMetadataPresence != SchemaProven || b511CandidateCatalog.B511CandidatePhysicalQualification != SchemaUnknown {
		b511CandidateT.Fatalf("finite catalog or grade changed: %d", len(b511CandidateCatalog.B511CandidateEntries))
	}
	for b511CandidateIndex, b511CandidateRow := range b511CandidateCatalog.B511CandidateEntries {
		if b511CandidateRow.B511CandidateGrade != SchemaHypothesis || b511CandidateRow.B511CandidateInstalledSupport != SchemaUnknown || b511CandidateRow.B511CandidatePhysicalPrefix != SchemaUnknown || b511CandidateRow.B511CandidateDefault != nil || b511CandidateRow.B511CandidateRawPolicy != "preserve original carrier" || b511CandidateRow.B511CandidateKey.B511CandidateSecondary != nil {
			b511CandidateT.Fatal("claim/default/raw boundary changed")
		}
		if b511CandidateRow.B511CandidateKey.B511CandidateVariant != "10" && b511CandidateRow.B511CandidateKey.B511CandidateVariant != "11" && b511CandidateRow.B511CandidateKey.B511CandidatePrimary != nil {
			b511CandidateT.Fatal("invented unprovided primary prefix")
		}
		b511CandidateKnownWidth := b511CandidateRow.B511CandidateKey.B511CandidateField == "wish_temperature" || b511CandidateRow.B511CandidateKey.B511CandidateField == "room_setpoint"
		if !b511CandidateKnownWidth && b511CandidateRow.B511CandidateWidthBits != nil {
			b511CandidateT.Fatal("invented unprovided width")
		}

		for _, b511CandidateOther := range b511CandidateCatalog.B511CandidateEntries[b511CandidateIndex+1:] {
			if b511CandidateReflect.DeepEqual(b511CandidateRow.B511CandidateKey, b511CandidateOther.B511CandidateKey) {
				b511CandidateT.Fatal("duplicate exact key")
			}
		}
		if b511CandidateRow.B511CandidateKey.B511CandidateVariant == "10" || b511CandidateRow.B511CandidateKey.B511CandidateVariant == "11" {
			b511CandidateExpected := uint8(0xB5)
			if b511CandidateRow.B511CandidateKey.B511CandidateDirection == "response" {
				b511CandidateExpected = 0xF5
			}
			if b511CandidateRow.B511CandidateKey.B511CandidateOperation != nil || *b511CandidateRow.B511CandidateKey.B511CandidatePrimary != b511CandidateExpected || *b511CandidateRow.B511CandidatePosition != 0 || b511CandidateRow.B511CandidateExtent != nil || b511CandidateRow.B511CandidateWidthBits != nil {
				b511CandidateT.Fatal("partial label materialized selector or body")
			}
		}
	}
	for _, b511CandidateVariant := range []string{"10", "11"} {
		for b511CandidateIndex, b511CandidateDirection := range []string{"request", "response"} {
			b511CandidateKey := B511CandidateB511Key{B511CandidateFamily: "vaillant", B511CandidateVariant: b511CandidateVariant, B511CandidatePrimary: b511CandidateByte([]uint8{0xB5, 0xF5}[b511CandidateIndex]), B511CandidateDirection: b511CandidateDirection, B511CandidateField: "prefix"}
			if _, b511CandidateFound, b511CandidateError := B511CandidateLookup(b511CandidateKey); !b511CandidateFound || b511CandidateError != nil {
				b511CandidateT.Fatal("partial control failed")
			}
			for _, b511CandidateOperation := range []uint8{0, 10, 11, 16, 17} {
				b511CandidateKey.B511CandidateOperation = b511CandidateByte(b511CandidateOperation)
				if _, b511CandidateFound, b511CandidateError := B511CandidateLookup(b511CandidateKey); b511CandidateFound || b511CandidateError != nil {
					b511CandidateT.Fatal("digits or absent operation became selector")
				}
			}
		}
	}
	if _, b511CandidateFound, b511CandidateError := B511CandidateLookup(B511CandidateB511Key{}); b511CandidateFound || b511CandidateError != nil {
		b511CandidateT.Fatal("constructed universal empty poll")
	}
}

func b511CandidateMutate(b511CandidateValue b511CandidateReflect.Value) {
	switch b511CandidateValue.Kind() {
	case b511CandidateReflect.Ptr:
		if !b511CandidateValue.IsNil() {
			b511CandidateMutate(b511CandidateValue.Elem())
		}
	case b511CandidateReflect.Struct:
		for b511CandidateIndex := 0; b511CandidateIndex < b511CandidateValue.NumField(); b511CandidateIndex++ {
			b511CandidateMutate(b511CandidateValue.Field(b511CandidateIndex))
		}
	case b511CandidateReflect.Slice:
		for b511CandidateIndex := 0; b511CandidateIndex < b511CandidateValue.Len(); b511CandidateIndex++ {
			b511CandidateMutate(b511CandidateValue.Index(b511CandidateIndex))
		}
	case b511CandidateReflect.Map:
		for _, b511CandidateKey := range b511CandidateValue.MapKeys() {
			b511CandidateValue.SetMapIndex(b511CandidateKey, b511CandidateReflect.ValueOf("mutated"))
		}
	case b511CandidateReflect.String:
		b511CandidateValue.SetString("mutated")
	case b511CandidateReflect.Int:
		b511CandidateValue.SetInt(99)
	case b511CandidateReflect.Uint8:
		b511CandidateValue.SetUint(99)
	case b511CandidateReflect.Float64:
		b511CandidateValue.SetFloat(99)
	}
}

func TestAb67FreshNestedStorage(b511CandidateT *b511CandidateTesting.T) {
	b511CandidateExpected, b511CandidateError := B511CandidateCatalog()
	if b511CandidateError != nil {
		b511CandidateT.Fatal(b511CandidateError)
	}
	b511CandidateChanged, b511CandidateError := B511CandidateCatalog()
	if b511CandidateError != nil {
		b511CandidateT.Fatal(b511CandidateError)
	}
	b511CandidateMutate(b511CandidateReflect.ValueOf(&b511CandidateChanged).Elem())
	b511CandidateFresh, b511CandidateError := B511CandidateCatalog()
	if b511CandidateError != nil || !b511CandidateReflect.DeepEqual(b511CandidateExpected, b511CandidateFresh) {
		b511CandidateT.Fatal("accessor shares nested storage")
	}
	for _, b511CandidateExpectedRow := range b511CandidateExpected.B511CandidateEntries {
		b511CandidateRow, b511CandidateFound, b511CandidateError := B511CandidateLookup(b511CandidateExpectedRow.B511CandidateKey)
		if !b511CandidateFound || b511CandidateError != nil || !b511CandidateReflect.DeepEqual(b511CandidateRow, b511CandidateExpectedRow) {
			b511CandidateT.Fatal("lookup differs from production catalog")
		}
		b511CandidateMutate(b511CandidateReflect.ValueOf(&b511CandidateRow).Elem())
		b511CandidateAgain, b511CandidateFound, b511CandidateError := B511CandidateLookup(b511CandidateExpectedRow.B511CandidateKey)
		if !b511CandidateFound || b511CandidateError != nil || !b511CandidateReflect.DeepEqual(b511CandidateAgain, b511CandidateExpectedRow) {
			b511CandidateT.Fatal("lookup shares nested storage")
		}
	}
}
