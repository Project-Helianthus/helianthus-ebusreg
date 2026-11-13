package vaillant

import (
	b510CandidateReflect "reflect"
	b510CandidateTesting "testing"
)

func b510CandidatePointer[b510CandidateT any](b510CandidateValue b510CandidateT) *b510CandidateT {
	return &b510CandidateValue
}

func b510CandidateGet(b510CandidateT *b510CandidateTesting.T, b510CandidateOperation uint8, b510CandidateField string) B510CandidateEntry {
	b510CandidateT.Helper()
	b510CandidateEntry, b510CandidateOK := B510CandidateLookup(B510CandidateKey{"vaillant", "B510", 0xb5, b510CandidatePointer(uint8(0x10)), b510CandidatePointer(b510CandidateOperation), b510CandidatePointer("request"), b510CandidateField})
	if !b510CandidateOK {
		b510CandidateT.Fatalf("missing OP%02x %s", b510CandidateOperation, b510CandidateField)
	}
	return b510CandidateEntry
}

func TestAb65FiniteCatalog(b510CandidateT *b510CandidateTesting.T) {
	b510CandidateCatalog := B510Candidates()
	if b510CandidateCatalog.B510CandidateMetadataPresence != "Proven" || b510CandidateCatalog.B510CandidatePhysicalQualification != "Unknown" || b510CandidateCatalog.B510CandidateRoutableInvocation != "Unknown" || len(b510CandidateCatalog.B510CandidateEntries) != 12 || len(b510CandidateCatalog.B510CandidateLimitations) != 5 {
		b510CandidateT.Fatal("finite scope or grades changed")
	}
	b510CandidateFields := map[uint8][]string{0: {"disable"}, 2: {"remote_status", "CH_setpoint_upper_limit", "DHW_setpoint_upper_limit"}, 3: {"climate", "capacity", "cost", "force"}, 4: {"energy"}, 5: {"DHW_setpoint", "adjustment"}, 6: {"status"}}
	for b510CandidateOperation, b510CandidateNames := range b510CandidateFields {
		for _, b510CandidateName := range b510CandidateNames {
			b510CandidateEntry := b510CandidateGet(b510CandidateT, b510CandidateOperation, b510CandidateName)
			if b510CandidateEntry.B510CandidateGrade != "Hypothesis" || b510CandidateEntry.B510CandidatePosition != nil || b510CandidateEntry.B510CandidateWidthBits == nil || *b510CandidateEntry.B510CandidateWidthBits != 8 || b510CandidateEntry.B510CandidateInterpretation != "Unknown" {
				b510CandidateT.Fatal("invented qualification or coordinates", b510CandidateName)
			}
		}
	}
	for _, b510CandidateEntry := range b510CandidateCatalog.B510CandidateEntries {
		b510CandidateGot, b510CandidateOK := B510CandidateLookup(b510CandidateEntry.B510CandidateKey)
		if !b510CandidateOK || !b510CandidateReflect.DeepEqual(b510CandidateGot, b510CandidateEntry) {
			b510CandidateT.Fatal("accessor/lookup mismatch")
		}
	}
}

func TestAb65ExactKeyAxes(b510CandidateT *b510CandidateTesting.T) {
	b510CandidateControl := b510CandidateGet(b510CandidateT, 0, "disable").B510CandidateKey
	b510CandidateChanges := map[string]func(*B510CandidateKey){
		"ab65_family":        func(b510CandidateKey *B510CandidateKey) { b510CandidateKey.B510CandidateFamily = "other" },
		"ab65_family_absent": func(b510CandidateKey *B510CandidateKey) { b510CandidateKey.B510CandidateFamily = "" },
		"ab65_variant":       func(b510CandidateKey *B510CandidateKey) { b510CandidateKey.B510CandidateVariant = "other" },
		"ab65_primary":       func(b510CandidateKey *B510CandidateKey) { b510CandidateKey.B510CandidatePrimary = 0 },
		"ab65_secondary": func(b510CandidateKey *B510CandidateKey) {
			b510CandidateKey.B510CandidateSecondary = b510CandidatePointer(uint8(0))
		},
		"ab65_secondary_absent": func(b510CandidateKey *B510CandidateKey) { b510CandidateKey.B510CandidateSecondary = nil },
		"ab65_operation": func(b510CandidateKey *B510CandidateKey) {
			b510CandidateKey.B510CandidateOperation = b510CandidatePointer(uint8(1))
		},
		"ab65_operation_absent": func(b510CandidateKey *B510CandidateKey) { b510CandidateKey.B510CandidateOperation = nil },
		"ab65_direction": func(b510CandidateKey *B510CandidateKey) {
			b510CandidateKey.B510CandidateDirection = b510CandidatePointer("response")
		},
		"ab65_direction_absent": func(b510CandidateKey *B510CandidateKey) { b510CandidateKey.B510CandidateDirection = nil },
		"ab65_field":            func(b510CandidateKey *B510CandidateKey) { b510CandidateKey.B510CandidateField = "other" },
	}
	for b510CandidateName, b510CandidateChange := range b510CandidateChanges {
		b510CandidateT.Run(b510CandidateName, func(b510CandidateT *b510CandidateTesting.T) {
			if _, b510CandidateOK := B510CandidateLookup(b510CandidateControl); !b510CandidateOK {
				b510CandidateT.Fatal("invalid positive control")
			}
			b510CandidateKey := b510CandidateControl
			b510CandidateChange(&b510CandidateKey)
			if _, b510CandidateOK := B510CandidateLookup(b510CandidateKey); b510CandidateOK {
				b510CandidateT.Fatal("axis wildcarded")
			}
		})
	}
}

func TestAb65ScalesAndEnums(b510CandidateT *b510CandidateTesting.T) {
	b510CandidateCapacity := b510CandidateGet(b510CandidateT, 3, "capacity")
	if !b510CandidateReflect.DeepEqual(b510CandidateCapacity.B510CandidateRawDomain, []uint8{0, 250}) || b510CandidateCapacity.B510CandidateDivisor == nil || *b510CandidateCapacity.B510CandidateDivisor != 10 || b510CandidateCapacity.B510CandidateUnit == nil || *b510CandidateCapacity.B510CandidateUnit != "kW" {
		b510CandidateT.Fatal("capacity domain/scale")
	}
	for b510CandidateIndex, b510CandidateRaw := range []uint8{0, 250, 251} {
		b510CandidateInside := b510CandidateRaw >= b510CandidateCapacity.B510CandidateRawDomain[0] && b510CandidateRaw <= b510CandidateCapacity.B510CandidateRawDomain[1]
		if b510CandidateInside != []bool{true, true, false}[b510CandidateIndex] {
			b510CandidateT.Fatal("capacity domain")
		}
		if b510CandidateInside && float64(b510CandidateRaw) / *b510CandidateCapacity.B510CandidateDivisor != []float64{0, 25}[b510CandidateIndex] {
			b510CandidateT.Fatal("capacity example")
		}
	}
	for _, b510CandidateEntry := range []B510CandidateEntry{b510CandidateGet(b510CandidateT, 2, "CH_setpoint_upper_limit"), b510CandidateGet(b510CandidateT, 2, "DHW_setpoint_upper_limit"), b510CandidateGet(b510CandidateT, 5, "DHW_setpoint")} {
		if b510CandidateEntry.B510CandidateDivisor == nil || float64(1) / *b510CandidateEntry.B510CandidateDivisor != 0.5 || b510CandidateEntry.B510CandidateUnit == nil || *b510CandidateEntry.B510CandidateUnit != "degrees Celsius" {
			b510CandidateT.Fatal("temperature 01")
		}
	}
	b510CandidateExpected := map[string]map[uint8]string{
		"climate": {0: "off", 1: "heating", 2: "cooling"},
		"force":   {0: "normal", 1: "maximum", 2: "minimum", 3: "off", 4: "cooling-leak-test", 5: "cooling-charge-test", 6: "capacity-test", 7: "I/O-test"},
		"energy":  {0: "unknown", 1: "gas", 2: "electricity-one", 3: "electricity-two", 4: "fuel", 5: "pellet"},
	}
	for b510CandidateField, b510CandidateValues := range b510CandidateExpected {
		b510CandidateOperation := uint8(3)
		if b510CandidateField == "energy" {
			b510CandidateOperation = 4
		}
		b510CandidateEntry := b510CandidateGet(b510CandidateT, b510CandidateOperation, b510CandidateField)
		if !b510CandidateReflect.DeepEqual(b510CandidateEntry.B510CandidateValues, b510CandidateValues) {
			b510CandidateT.Fatal("enum labels", b510CandidateField)
		}
		if _, b510CandidateMapped := b510CandidateEntry.B510CandidateValues[8]; b510CandidateMapped {
			b510CandidateT.Fatal("unmapped value 8")
		}
	}
	b510CandidateCost := b510CandidateGet(b510CandidateT, 3, "cost")
	if !b510CandidateReflect.DeepEqual(b510CandidateCost.B510CandidateRawDomain, []uint8{0, 250}) || b510CandidateCost.B510CandidateUnit == nil || *b510CandidateCost.B510CandidateUnit != "reciprocal-MWh" || b510CandidateCost.B510CandidateDivisor != nil {
		b510CandidateT.Fatal("cost metadata")
	}
	b510CandidateRemote := b510CandidateGet(b510CandidateT, 2, "remote_status")
	if b510CandidateRemote.B510CandidateDivisor != nil || b510CandidateRemote.B510CandidateValues != nil || !b510CandidateReflect.DeepEqual(b510CandidateRemote.B510CandidateRawDomain, []uint8{0, 255}) {
		b510CandidateT.Fatal("raw remote status")
	}
}

func TestAb65AdjustmentAndReserved(b510CandidateT *b510CandidateTesting.T) {
	b510CandidateAdjustment := b510CandidateGet(b510CandidateT, 5, "adjustment")
	if len(b510CandidateAdjustment.B510CandidateBits) != 8 {
		b510CandidateT.Fatal("missing adjustment bits")
	}
	for b510CandidateIndex, b510CandidateRaw := range []uint8{0, 3, 0xfc} {
		b510CandidateReserved := uint8(0)
		for _, b510CandidateBit := range b510CandidateAdjustment.B510CandidateBits {
			if b510CandidateBit.B510CandidateReserved {
				b510CandidateReserved |= b510CandidateRaw & (1 << b510CandidateBit.B510CandidateBit)
				continue
			}
			b510CandidateWant := []string{"released", "blocked", "released"}[b510CandidateIndex]
			if b510CandidateBit.B510CandidateValues[(b510CandidateRaw>>b510CandidateBit.B510CandidateBit)&1] != b510CandidateWant {
				b510CandidateT.Fatal("blocked polarity")
			}
		}
		if b510CandidateReserved != []uint8{0, 0, 0xfc}[b510CandidateIndex] {
			b510CandidateT.Fatal("reserved bits erased")
		}
	}
	for b510CandidateIndex, b510CandidateBit := range b510CandidateAdjustment.B510CandidateBits {
		if b510CandidateBit.B510CandidateBit != uint8(b510CandidateIndex) || b510CandidateBit.B510CandidateReserved != (b510CandidateIndex >= 2) || b510CandidateBit.B510CandidateExpected != nil {
			b510CandidateT.Fatal("adjustment bit identity")
		}
	}
	if *b510CandidateAdjustment.B510CandidateBits[0].B510CandidateMeaning != "CH" || *b510CandidateAdjustment.B510CandidateBits[1].B510CandidateMeaning != "DHW" {
		b510CandidateT.Fatal("adjustment field identity")
	}
	b510CandidateStatus := b510CandidateGet(b510CandidateT, 6, "status")
	if len(b510CandidateStatus.B510CandidateBits) != 8 || !b510CandidateReflect.DeepEqual(b510CandidateStatus.B510CandidateBits[0].B510CandidateValues, map[uint8]string{0: "off", 1: "on"}) {
		b510CandidateT.Fatal("OP06 bit0")
	}
	for b510CandidateIndex, b510CandidateBit := range b510CandidateStatus.B510CandidateBits[1:] {
		if b510CandidateBit.B510CandidateBit != uint8(b510CandidateIndex+1) || !b510CandidateBit.B510CandidateReserved || b510CandidateBit.B510CandidateExpected == nil || *b510CandidateBit.B510CandidateExpected != 0 {
			b510CandidateT.Fatal("OP06 reserved zero")
		}
	}
}

func TestAb65DisableModesAndUnresolved(b510CandidateT *b510CandidateTesting.T) {
	b510CandidateDisable := b510CandidateGet(b510CandidateT, 0, "disable")
	if len(b510CandidateDisable.B510CandidateBits) != 8 || b510CandidateDisable.B510CandidateValues != nil {
		b510CandidateT.Fatal("invented mode codes")
	}
	for b510CandidateIndex, b510CandidateBit := range b510CandidateDisable.B510CandidateBits {
		if b510CandidateBit.B510CandidateBit != uint8(b510CandidateIndex) {
			b510CandidateT.Fatal("bit identity")
		}
		if b510CandidateIndex < 3 && (!b510CandidateReflect.DeepEqual(b510CandidateBit.B510CandidateActiveModes, []string{"automatic"}) || !b510CandidateReflect.DeepEqual(b510CandidateBit.B510CandidateIgnoredModes, []string{"forced-heating", "forced-DHW", "forced-off"})) {
			b510CandidateT.Fatal("disable mode applicability")
		}
	}
	for _, b510CandidateMode := range []string{"automatic", "forced-heating", "forced-DHW", "forced-off"} {
		b510CandidateActive := uint8(0)
		for _, b510CandidateBit := range b510CandidateDisable.B510CandidateBits[:3] {
			for _, b510CandidateAllowed := range b510CandidateBit.B510CandidateActiveModes {
				if b510CandidateAllowed == b510CandidateMode {
					b510CandidateActive |= uint8(7) & (1 << b510CandidateBit.B510CandidateBit)
				}
			}
		}
		b510CandidateWant := uint8(0)
		if b510CandidateMode == "automatic" {
			b510CandidateWant = 7
		}
		if b510CandidateActive != b510CandidateWant {
			b510CandidateT.Fatal("07 mode example", b510CandidateMode)
		}
	}
	for _, b510CandidateBit := range b510CandidateDisable.B510CandidateBits[3:5] {
		if b510CandidateBit.B510CandidateGrade != "Unknown" || b510CandidateBit.B510CandidateMeaning == nil || *b510CandidateBit.B510CandidateMeaning != "timer assignment unresolved" || b510CandidateBit.B510CandidateValues != nil {
			b510CandidateT.Fatal("settled competing timer assignments")
		}
	}
	if b510CandidateDisable.B510CandidateBits[5].B510CandidateGrade != "Unknown" || b510CandidateDisable.B510CandidateBits[5].B510CandidateMeaning != nil {
		b510CandidateT.Fatal("omitted meaning invented")
	}
	if *b510CandidateDisable.B510CandidateBits[6].B510CandidateMeaning != "relates to subsequent flame status" {
		b510CandidateT.Fatal("flame relationship")
	}
	b510CandidateLegionella := b510CandidateDisable.B510CandidateBits[7]
	if *b510CandidateLegionella.B510CandidateMeaning != "starts legionella protection" || b510CandidateLegionella.B510CandidatePersistsUntilStatus == nil || !*b510CandidateLegionella.B510CandidatePersistsUntilStatus || b510CandidateLegionella.B510CandidateCommandsTermination == nil || *b510CandidateLegionella.B510CandidateCommandsTermination || b510CandidateLegionella.B510CandidateActiveModes != nil || b510CandidateLegionella.B510CandidateIgnoredModes != nil {
		b510CandidateT.Fatal("bit7 independent persistence")
	}
	for _, b510CandidateRaw := range []uint8{7, 0x80, 0x87, 0x18} {
		if (b510CandidateRaw&0x80 != 0) != (b510CandidateRaw == 0x80 || b510CandidateRaw == 0x87) {
			b510CandidateT.Fatal("bit7 conflated")
		}
	}
}

func b510CandidateMutate(b510CandidateValue b510CandidateReflect.Value) {
	switch b510CandidateValue.Kind() {
	case b510CandidateReflect.Pointer:
		if !b510CandidateValue.IsNil() {
			b510CandidateMutate(b510CandidateValue.Elem())
		}
	case b510CandidateReflect.Struct:
		for b510CandidateIndex := 0; b510CandidateIndex < b510CandidateValue.NumField(); b510CandidateIndex++ {
			b510CandidateMutate(b510CandidateValue.Field(b510CandidateIndex))
		}
	case b510CandidateReflect.Slice:
		for b510CandidateIndex := 0; b510CandidateIndex < b510CandidateValue.Len(); b510CandidateIndex++ {
			b510CandidateMutate(b510CandidateValue.Index(b510CandidateIndex))
		}
	case b510CandidateReflect.Map:
		for _, b510CandidateKey := range b510CandidateValue.MapKeys() {
			b510CandidateValue.SetMapIndex(b510CandidateKey, b510CandidateReflect.ValueOf("ab65_changed"))
		}
	case b510CandidateReflect.String:
		b510CandidateValue.SetString("ab65_changed")
	case b510CandidateReflect.Bool:
		b510CandidateValue.SetBool(!b510CandidateValue.Bool())
	case b510CandidateReflect.Uint8:
		b510CandidateValue.SetUint(99)
	case b510CandidateReflect.Int:
		b510CandidateValue.SetInt(99)
	case b510CandidateReflect.Float64:
		b510CandidateValue.SetFloat(99)
	}
}

func TestAb65DeepCopy(b510CandidateT *b510CandidateTesting.T) {
	b510CandidateBefore := B510Candidates()
	b510CandidateReturned := B510Candidates()
	b510CandidateMutate(b510CandidateReflect.ValueOf(&b510CandidateReturned).Elem())
	if !b510CandidateReflect.DeepEqual(B510Candidates(), b510CandidateBefore) {
		b510CandidateT.Fatal("shared nested accessor backing")
	}
	for _, b510CandidateEntry := range b510CandidateBefore.B510CandidateEntries {
		b510CandidateGot, b510CandidateOK := B510CandidateLookup(b510CandidateEntry.B510CandidateKey)
		if !b510CandidateOK {
			b510CandidateT.Fatal("lookup")
		}
		b510CandidateMutate(b510CandidateReflect.ValueOf(&b510CandidateGot).Elem())
		b510CandidateFresh, b510CandidateOK := B510CandidateLookup(b510CandidateEntry.B510CandidateKey)
		if !b510CandidateOK || !b510CandidateReflect.DeepEqual(b510CandidateFresh, b510CandidateEntry) {
			b510CandidateT.Fatal("shared nested lookup backing")
		}
	}
}
