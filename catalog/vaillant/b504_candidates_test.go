package vaillant

import (
	b504CandidateBinary "encoding/binary"
	b504CandidateReflect "reflect"
	b504CandidateTesting "testing"
)

func b504CandidateByte(b504CandidateValue uint8) *uint8 { return &b504CandidateValue }

func b504CandidateKey(b504CandidateVariant string, b504CandidateOperation *uint8, b504CandidateField string) B504CandidateKey {
	return B504CandidateKey{B504CandidateFamily: "vaillant", B504CandidateVariant: b504CandidateVariant, B504CandidatePrimary: 0xb5, B504CandidateSecondary: b504CandidateByte(4), B504CandidateOperation: b504CandidateOperation, B504CandidateDirection: "response", B504CandidateField: b504CandidateField}
}

func b504CandidateGet(b504CandidateT *b504CandidateTesting.T, b504CandidateKeyValue B504CandidateKey) B504CandidateEntry {
	b504CandidateT.Helper()
	b504CandidateEntry, b504CandidateFound, b504CandidateError := B504LookupCandidate(b504CandidateKeyValue)
	if b504CandidateError != nil || !b504CandidateFound {
		b504CandidateT.Fatalf("lookup %+v: found=%v error=%v", b504CandidateKeyValue, b504CandidateFound, b504CandidateError)
	}
	return b504CandidateEntry
}

func TestAb53ExactKeyAxes(b504CandidateT *b504CandidateTesting.T) {
	b504CandidateControl := b504CandidateKey("outside-temperature-div16", b504CandidateByte(0x16), "corrected_outside_temperature")
	b504CandidateGet(b504CandidateT, b504CandidateControl)
	b504CandidateCases := []struct {
		B504CandidateName   string
		B504CandidateChange func(*B504CandidateKey)
	}{
		{"ab53Family", func(b504CandidateK *B504CandidateKey) { b504CandidateK.B504CandidateFamily = "other-family" }},
		{"ab53Variant", func(b504CandidateK *B504CandidateKey) { b504CandidateK.B504CandidateVariant = "other-variant" }},
		{"ab53Primary", func(b504CandidateK *B504CandidateKey) { b504CandidateK.B504CandidatePrimary = 0xb4 }},
		{"ab53Secondary", func(b504CandidateK *B504CandidateKey) { b504CandidateK.B504CandidateSecondary = b504CandidateByte(5) }},
		{"ab53SecondaryAbsent", func(b504CandidateK *B504CandidateKey) { b504CandidateK.B504CandidateSecondary = nil }},
		{"ab53SecondaryZero", func(b504CandidateK *B504CandidateKey) { b504CandidateK.B504CandidateSecondary = b504CandidateByte(0) }},
		{"ab53Operation", func(b504CandidateK *B504CandidateKey) {
			b504CandidateK.B504CandidateOperation = b504CandidateByte(0x17)
		}},
		{"ab53OperationAbsent", func(b504CandidateK *B504CandidateKey) { b504CandidateK.B504CandidateOperation = nil }},
		{"ab53OperationZero", func(b504CandidateK *B504CandidateKey) { b504CandidateK.B504CandidateOperation = b504CandidateByte(0) }},
		{"ab53Direction", func(b504CandidateK *B504CandidateKey) { b504CandidateK.B504CandidateDirection = "request" }},
		{"ab53Field", func(b504CandidateK *B504CandidateKey) { b504CandidateK.B504CandidateField = "other_field" }},
	}
	for _, b504CandidateCase := range b504CandidateCases {
		b504CandidateT.Run(b504CandidateCase.B504CandidateName, func(b504CandidateT *b504CandidateTesting.T) {
			b504CandidateK := b504CandidateControl
			b504CandidateCase.B504CandidateChange(&b504CandidateK)
			_, b504CandidateFound, b504CandidateError := B504LookupCandidate(b504CandidateK)
			if b504CandidateError != nil || b504CandidateFound {
				b504CandidateT.Fatalf("altered axis matched: %v %v", b504CandidateFound, b504CandidateError)
			}
		})
	}
	b504CandidateUnknown := b504CandidateKey("solar-scalar-two", nil, "solar_state")
	b504CandidateGet(b504CandidateT, b504CandidateUnknown)
	b504CandidateUnknown.B504CandidateOperation = b504CandidateByte(0)
	if _, b504CandidateFound, b504CandidateError := B504LookupCandidate(b504CandidateUnknown); b504CandidateError != nil || b504CandidateFound {
		b504CandidateT.Fatal("absent operation matched explicit zero")
	}
}

func TestAb53TemperatureAlternatives(b504CandidateT *b504CandidateTesting.T) {
	b504CandidateCases := []struct {
		B504CandidateVariant  string
		B504CandidateDivisor  int
		B504CandidatePositive float64
		B504CandidateNegative float64
	}{
		{"outside-temperature-div16", 16, 1, -1},
		{"outside-temperature-div256", 256, 0.0625, -0.0625},
	}
	for _, b504CandidateCase := range b504CandidateCases {
		b504CandidateT.Run("ab53"+b504CandidateCase.B504CandidateVariant, func(b504CandidateT *b504CandidateTesting.T) {
			b504CandidateEntry := b504CandidateGet(b504CandidateT, b504CandidateKey(b504CandidateCase.B504CandidateVariant, b504CandidateByte(0x16), "corrected_outside_temperature"))
			if b504CandidateEntry.B504CandidateDivisor == nil || *b504CandidateEntry.B504CandidateDivisor != b504CandidateCase.B504CandidateDivisor || b504CandidateEntry.B504CandidateSignedness != "signed" || b504CandidateEntry.B504CandidateByteOrder != "little-endian" || b504CandidateEntry.B504CandidateWidth == nil || *b504CandidateEntry.B504CandidateWidth != 2 || b504CandidateEntry.B504CandidateUnit == nil || *b504CandidateEntry.B504CandidateUnit != "degrees Celsius" {
				b504CandidateT.Fatal("temperature candidate changed", b504CandidateEntry)
			}
			b504CandidatePositive := float64(int16(b504CandidateBinary.LittleEndian.Uint16([]byte{0x10, 0}))) / float64(*b504CandidateEntry.B504CandidateDivisor)
			b504CandidateNegative := float64(int16(b504CandidateBinary.LittleEndian.Uint16([]byte{0xf0, 0xff}))) / float64(*b504CandidateEntry.B504CandidateDivisor)
			if b504CandidatePositive != b504CandidateCase.B504CandidatePositive || b504CandidateNegative != b504CandidateCase.B504CandidateNegative {
				b504CandidateT.Fatalf("synthetic fields: %v %v", b504CandidatePositive, b504CandidateNegative)
			}
		})
	}
	b504CandidateTail := b504CandidateGet(b504CandidateT, b504CandidateKey("unused-tail", b504CandidateByte(0x16), "unused_tail"))
	if b504CandidateTail.B504CandidatePosition != nil || b504CandidateTail.B504CandidateWidth != nil || b504CandidateTail.B504CandidateGrade != SchemaUnknown {
		b504CandidateT.Fatal("tail acquired an extent or semantics")
	}
}

func b504CandidateLabel(b504CandidateEntry B504CandidateEntry, b504CandidateRaw int) string {
	if b504CandidateValue, b504CandidateFound := b504CandidateEntry.B504CandidateLabels[b504CandidateRaw]; b504CandidateFound {
		return b504CandidateValue
	}
	if b504CandidateRaw > 0 && b504CandidateEntry.B504CandidatePositiveLabel != nil {
		return *b504CandidateEntry.B504CandidatePositiveLabel
	}
	return ""
}

func TestAb53ZeroAndPositiveBranches(b504CandidateT *b504CandidateTesting.T) {
	for _, b504CandidateCase := range []struct {
		B504CandidateVariant   string
		B504CandidateOperation uint8
		B504CandidateField     string
		B504CandidateZero      string
		B504CandidatePositive  string
	}{
		{"setpoint", 0x20, "setpoint", "no-demand", "nonzero candidate"},
		{"protection", 0x21, "protection_temperature", "off", "protection-active"},
	} {
		b504CandidateEntry := b504CandidateGet(b504CandidateT, b504CandidateKey(b504CandidateCase.B504CandidateVariant, b504CandidateByte(b504CandidateCase.B504CandidateOperation), b504CandidateCase.B504CandidateField))
		if b504CandidateEntry.B504CandidateLabels[0] != b504CandidateCase.B504CandidateZero || b504CandidateEntry.B504CandidatePositiveLabel == nil || *b504CandidateEntry.B504CandidatePositiveLabel != b504CandidateCase.B504CandidatePositive {
			b504CandidateT.Fatal("zero/positive distinction lost")
		}
		if b504CandidateLabel(b504CandidateEntry, 0) != b504CandidateCase.B504CandidateZero || b504CandidateLabel(b504CandidateEntry, 1) != b504CandidateCase.B504CandidatePositive {
			b504CandidateT.Fatal("synthetic raw zero/one branch changed")
		}
		if _, b504CandidatePresent := b504CandidateEntry.B504CandidateLabels[1]; b504CandidatePresent || b504CandidateEntry.B504CandidateDivisor != nil {
			b504CandidateT.Fatal("invented exact nonzero arithmetic")
		}
	}
}

func TestAb53ContextualUnavailable(b504CandidateT *b504CandidateTesting.T) {
	b504CandidateRoomKey := b504CandidateKey("room-temperature", b504CandidateByte(0x26), "room_temperature")
	b504CandidateRoom := b504CandidateGet(b504CandidateT, b504CandidateRoomKey)
	if b504CandidateRoom.B504CandidateLabels[255] != "unavailable" {
		b504CandidateT.Fatal("missing contextual FF")
	}
	b504CandidateRoomKey.B504CandidateOperation = b504CandidateByte(0x20)
	if _, b504CandidateFound, b504CandidateError := B504LookupCandidate(b504CandidateRoomKey); b504CandidateError != nil || b504CandidateFound {
		b504CandidateT.Fatal("OP26 label inherited by unrelated operation")
	}
	b504CandidateOther := b504CandidateGet(b504CandidateT, b504CandidateKey("setpoint", b504CandidateByte(0x20), "setpoint"))
	if _, b504CandidatePresent := b504CandidateOther.B504CandidateLabels[255]; b504CandidatePresent {
		b504CandidateT.Fatal("unrelated selector inherited FF label")
	}
	b504CandidateFlag := b504CandidateGet(b504CandidateT, b504CandidateKey("adjustment", b504CandidateByte(0x26), "adjustment_flag"))
	if b504CandidateFlag.B504CandidateLabels[1] != "setpoint change permitted" || len(b504CandidateFlag.B504CandidateLabels) != 1 {
		b504CandidateT.Fatal("adjustment candidate changed")
	}
}

func b504CandidateContains(b504CandidateDomain *B504CandidateDomain, b504CandidateValue int) bool {
	if b504CandidateDomain == nil {
		return false
	}
	if b504CandidateDomain.B504CandidateValues != nil {
		for _, b504CandidateAllowed := range b504CandidateDomain.B504CandidateValues {
			if b504CandidateAllowed == b504CandidateValue {
				return true
			}
		}
		return false
	}
	if b504CandidateDomain.B504CandidateMinimum == nil || b504CandidateDomain.B504CandidateMaximum == nil || b504CandidateValue < *b504CandidateDomain.B504CandidateMinimum || b504CandidateValue > *b504CandidateDomain.B504CandidateMaximum {
		return false
	}
	return b504CandidateDomain.B504CandidateStep == nil || (b504CandidateValue-*b504CandidateDomain.B504CandidateMinimum)%*b504CandidateDomain.B504CandidateStep == 0
}

func TestAb53SolarDomains(b504CandidateT *b504CandidateTesting.T) {
	b504CandidateThree := b504CandidateGet(b504CandidateT, b504CandidateKey("solar-scalar-three", nil, "solar_state"))
	b504CandidateTwo := b504CandidateGet(b504CandidateT, b504CandidateKey("solar-scalar-two", nil, "solar_state"))
	if !b504CandidateReflect.DeepEqual(b504CandidateThree.B504CandidateDomain.B504CandidateValues, []int{0, 1, 2}) || !b504CandidateReflect.DeepEqual(b504CandidateTwo.B504CandidateDomain.B504CandidateValues, []int{0, 1}) {
		b504CandidateT.Fatal("scalar domains merged or changed")
	}
	if !b504CandidateContains(b504CandidateThree.B504CandidateDomain, 2) || b504CandidateContains(b504CandidateTwo.B504CandidateDomain, 2) {
		b504CandidateT.Fatal("code 2 admitted into two-value table")
	}
	b504CandidateBits := b504CandidateGet(b504CandidateT, b504CandidateKey("solar-bits", nil, "solar_state"))
	if b504CandidateBits.B504CandidateDomain != nil || b504CandidateBits.B504CandidateGrade != SchemaUnknown {
		b504CandidateT.Fatal("bit descriptions promoted to scalar domain")
	}
}

func TestAb53DurationAlternatives(b504CandidateT *b504CandidateTesting.T) {
	b504CandidateFirst := b504CandidateGet(b504CandidateT, b504CandidateKey("preheating-five-minute-60", nil, "preheating_duration"))
	b504CandidateSecond := b504CandidateGet(b504CandidateT, b504CandidateKey("preheating-300", nil, "preheating_duration"))
	if b504CandidateFirst.B504CandidateDomain == nil || b504CandidateFirst.B504CandidateDomain.B504CandidateStep == nil || *b504CandidateFirst.B504CandidateDomain.B504CandidateStep != 5 || b504CandidateSecond.B504CandidateDomain == nil || b504CandidateSecond.B504CandidateDomain.B504CandidateStep != nil {
		b504CandidateT.Fatal("duration step metadata changed")
	}
	for _, b504CandidateCase := range []struct {
		B504CandidateValue                      int
		B504CandidateFirst, B504CandidateSecond bool
	}{
		{0, true, true}, {1, false, true}, {60, true, true}, {65, false, true}, {300, false, true}, {301, false, false}, {-1, false, false},
	} {
		if b504CandidateContains(b504CandidateFirst.B504CandidateDomain, b504CandidateCase.B504CandidateValue) != b504CandidateCase.B504CandidateFirst || b504CandidateContains(b504CandidateSecond.B504CandidateDomain, b504CandidateCase.B504CandidateValue) != b504CandidateCase.B504CandidateSecond {
			b504CandidateT.Fatalf("duration %d domains changed", b504CandidateCase.B504CandidateValue)
		}
	}
}

func TestAb53CatalogBoundariesAndCopies(b504CandidateT *b504CandidateTesting.T) {
	b504CandidateBefore, b504CandidateError := B504Candidates()
	if b504CandidateError != nil {
		b504CandidateT.Fatal(b504CandidateError)
	}
	if len(b504CandidateBefore.B504CandidateEntries) != 12 || b504CandidateBefore.B504CandidateMetadataPresence != SchemaProven || b504CandidateBefore.B504CandidatePhysicalAdaptation != SchemaUnknown || b504CandidateBefore.B504CandidateInstalledSupport != SchemaUnknown || b504CandidateBefore.B504CandidatePollingPermission != SchemaUnknown || b504CandidateBefore.B504CandidateRoutableInvocation != SchemaUnknown {
		b504CandidateT.Fatal("catalog boundary changed")
	}
	for b504CandidateIndex, b504CandidateEntry := range b504CandidateBefore.B504CandidateEntries {
		if b504CandidateEntry.B504CandidatePosition != nil {
			b504CandidateT.Fatal("unprovided coordinate invented")
		}
		b504CandidateExpectedGrade := SchemaHypothesis
		if b504CandidateEntry.B504CandidateKey.B504CandidateVariant == "unused-tail" || b504CandidateEntry.B504CandidateKey.B504CandidateVariant == "solar-bits" {
			b504CandidateExpectedGrade = SchemaUnknown
		}
		if b504CandidateEntry.B504CandidateGrade != b504CandidateExpectedGrade {
			b504CandidateT.Fatal("semantic grade promoted")
		}
		if !b504CandidateReflect.DeepEqual(b504CandidateGet(b504CandidateT, b504CandidateEntry.B504CandidateKey), b504CandidateEntry) {
			b504CandidateT.Fatal("lookup/accessor differ")
		}
		for _, b504CandidateOther := range b504CandidateBefore.B504CandidateEntries[b504CandidateIndex+1:] {
			if b504CandidateReflect.DeepEqual(b504CandidateEntry.B504CandidateKey, b504CandidateOther.B504CandidateKey) {
				b504CandidateT.Fatal("duplicate key")
			}
		}
	}
	b504CandidateCopy, b504CandidateError := B504Candidates()
	if b504CandidateError != nil {
		b504CandidateT.Fatal(b504CandidateError)
	}
	b504CandidateCopy.B504CandidateLimitations[0] = "changed"
	for b504CandidateIndex := range b504CandidateCopy.B504CandidateEntries {
		b504CandidateEntry := &b504CandidateCopy.B504CandidateEntries[b504CandidateIndex]
		b504CandidateEntry.B504CandidateKey.B504CandidateFamily = "changed"
		if b504CandidateEntry.B504CandidateKey.B504CandidateOperation != nil {
			*b504CandidateEntry.B504CandidateKey.B504CandidateOperation = 0
		}
		*b504CandidateEntry.B504CandidateKey.B504CandidateSecondary = 0
		if b504CandidateEntry.B504CandidateWidth != nil {
			*b504CandidateEntry.B504CandidateWidth = 99
		}
		if b504CandidateEntry.B504CandidateDivisor != nil {
			*b504CandidateEntry.B504CandidateDivisor = 99
		}
		if b504CandidateEntry.B504CandidateUnit != nil {
			*b504CandidateEntry.B504CandidateUnit = "changed"
		}
		b504CandidateEntry.B504CandidateLabels[255] = "changed"
		if b504CandidateEntry.B504CandidatePositiveLabel != nil {
			*b504CandidateEntry.B504CandidatePositiveLabel = "changed"
		}
		if len(b504CandidateEntry.B504CandidateNotes) > 0 {
			b504CandidateEntry.B504CandidateNotes[0] = "changed"
		}
		if b504CandidateEntry.B504CandidateDomain != nil {
			if len(b504CandidateEntry.B504CandidateDomain.B504CandidateValues) > 0 {
				b504CandidateEntry.B504CandidateDomain.B504CandidateValues[0] = 99
			}
			if b504CandidateEntry.B504CandidateDomain.B504CandidateMinimum != nil {
				*b504CandidateEntry.B504CandidateDomain.B504CandidateMinimum = 99
			}
			if b504CandidateEntry.B504CandidateDomain.B504CandidateMaximum != nil {
				*b504CandidateEntry.B504CandidateDomain.B504CandidateMaximum = 99
			}
			if b504CandidateEntry.B504CandidateDomain.B504CandidateStep != nil {
				*b504CandidateEntry.B504CandidateDomain.B504CandidateStep = 99
			}
		}
	}
	b504CandidateAfter, b504CandidateError := B504Candidates()
	if b504CandidateError != nil || !b504CandidateReflect.DeepEqual(b504CandidateBefore, b504CandidateAfter) {
		b504CandidateT.Fatal("accessor leaked backing storage")
	}
	b504CandidateKeyValue := b504CandidateKey("room-temperature", b504CandidateByte(0x26), "room_temperature")
	b504CandidateLookup := b504CandidateGet(b504CandidateT, b504CandidateKeyValue)
	b504CandidateLookup.B504CandidateLabels[255] = "changed"
	*b504CandidateLookup.B504CandidateKey.B504CandidateOperation = 0
	if b504CandidateGet(b504CandidateT, b504CandidateKeyValue).B504CandidateLabels[255] != "unavailable" {
		b504CandidateT.Fatal("lookup leaked backing storage")
	}
}
