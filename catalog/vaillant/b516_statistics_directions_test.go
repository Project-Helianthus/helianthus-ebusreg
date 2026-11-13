package vaillant

import (
	b516StatisticsDirectionBinary "encoding/binary"
	b516StatisticsDirectionJSON "encoding/json"
	b516StatisticsDirectionReflect "reflect"
	b516StatisticsDirectionStrings "strings"
	b516StatisticsDirectionTesting "testing"
)

func b516StatisticsDirectionPointer[b516StatisticsDirectionType any](b516StatisticsDirectionValue b516StatisticsDirectionType) *b516StatisticsDirectionType {
	return &b516StatisticsDirectionValue
}
func b516StatisticsDirectionEqual(b516StatisticsDirectionT *b516StatisticsDirectionTesting.T, b516StatisticsDirectionGot, b516StatisticsDirectionWant any) {
	b516StatisticsDirectionT.Helper()
	if !b516StatisticsDirectionReflect.DeepEqual(b516StatisticsDirectionGot, b516StatisticsDirectionWant) {
		b516StatisticsDirectionT.Fatalf("got %#v; want %#v", b516StatisticsDirectionGot, b516StatisticsDirectionWant)
	}
}
func b516StatisticsDirectionGet(b516StatisticsDirectionT *b516StatisticsDirectionTesting.T, b516StatisticsDirectionVariant string, b516StatisticsDirectionOperation uint8, b516StatisticsDirectionSide, b516StatisticsDirectionField string, b516StatisticsDirectionSelector *uint8) B516StatisticsDirectionEntry {
	b516StatisticsDirectionT.Helper()
	b516StatisticsDirectionEntry, b516StatisticsDirectionFound, b516StatisticsDirectionErr := B516StatisticsDirectionLookup(B516StatisticsDirectionKey{B516StatisticsDirectionFamily: "vaillant", B516StatisticsDirectionVariant: b516StatisticsDirectionVariant, B516StatisticsDirectionOperation: b516StatisticsDirectionPointer(b516StatisticsDirectionOperation), B516StatisticsDirection: b516StatisticsDirectionSide, B516StatisticsDirectionField: b516StatisticsDirectionField, B516StatisticsDirectionSelector: b516StatisticsDirectionSelector})
	if b516StatisticsDirectionErr != nil || !b516StatisticsDirectionFound {
		b516StatisticsDirectionT.Fatalf("lookup %s/%d/%s/%s: %v %v", b516StatisticsDirectionVariant, b516StatisticsDirectionOperation, b516StatisticsDirectionSide, b516StatisticsDirectionField, b516StatisticsDirectionFound, b516StatisticsDirectionErr)
	}
	return b516StatisticsDirectionEntry
}
func TestAb72StatisticsSides(b516StatisticsDirectionT *b516StatisticsDirectionTesting.T) {
	for _, b516StatisticsDirectionCase := range []struct {
		b516StatisticsDirectionVariant, b516StatisticsDirectionSide string
		b516StatisticsDirectionLength                               uint8
		b516StatisticsDirectionValue                                bool
		b516StatisticsDirectionAccess                               *uint8
	}{
		{"ab72-statistics-read", "request", 7, false, b516StatisticsDirectionPointer(uint8(2))},
		{"ab72-statistics-read", "positive", 11, true, nil},
		{"ab72-statistics-write", "request", 11, true, b516StatisticsDirectionPointer(uint8(2))},
		{"ab72-statistics-write", "positive", 7, false, b516StatisticsDirectionPointer(uint8(3))},
	} {
		b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionGet(b516StatisticsDirectionT, b516StatisticsDirectionCase.b516StatisticsDirectionVariant, 16, b516StatisticsDirectionCase.b516StatisticsDirectionSide, "layout", nil).B516StatisticsDirectionLength, b516StatisticsDirectionPointer(b516StatisticsDirectionCase.b516StatisticsDirectionLength))
		b516StatisticsDirectionFlags := b516StatisticsDirectionGet(b516StatisticsDirectionT, b516StatisticsDirectionCase.b516StatisticsDirectionVariant, 16, b516StatisticsDirectionCase.b516StatisticsDirectionSide, "flags", nil)
		b516StatisticsDirectionExpected := map[string]uint8{"time-base-low": 0, "time-base-high": 1}
		if b516StatisticsDirectionCase.b516StatisticsDirectionAccess != nil {
			b516StatisticsDirectionExpected["access"] = *b516StatisticsDirectionCase.b516StatisticsDirectionAccess
		}
		b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionFlags.B516StatisticsDirectionBits, b516StatisticsDirectionExpected)
		b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionFlags.B516StatisticsDirectionValues["time-base:0"], "total")
		b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionFlags.B516StatisticsDirectionValues["time-base:1"], "daily")
		b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionFlags.B516StatisticsDirectionValues["time-base:2"], "monthly")
		b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionFlags.B516StatisticsDirectionValues["time-base:3"], "yearly")
		b516StatisticsDirectionKey := b516StatisticsDirectionFlags.B516StatisticsDirectionKey
		b516StatisticsDirectionKey.B516StatisticsDirectionField = "energy-value"
		b516StatisticsDirectionEntry, b516StatisticsDirectionFound, b516StatisticsDirectionErr := B516StatisticsDirectionLookup(b516StatisticsDirectionKey)
		b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionErr, nil)
		b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionFound, b516StatisticsDirectionCase.b516StatisticsDirectionValue)
		if b516StatisticsDirectionFound {
			b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionEntry.B516StatisticsDirectionEncoding, b516StatisticsDirectionPointer("binary32 LE"))
			b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionEntry.B516StatisticsDirectionUnit, b516StatisticsDirectionPointer("Wh"))
		}
		for _, b516StatisticsDirectionIndex := range []uint8{0, 255} {
			b516StatisticsDirectionPeriod := b516StatisticsDirectionGet(b516StatisticsDirectionT, b516StatisticsDirectionCase.b516StatisticsDirectionVariant, 16, b516StatisticsDirectionCase.b516StatisticsDirectionSide, "period-index", b516StatisticsDirectionPointer(b516StatisticsDirectionIndex))
			b516StatisticsDirectionEnergy := b516StatisticsDirectionGet(b516StatisticsDirectionT, b516StatisticsDirectionCase.b516StatisticsDirectionVariant, 16, b516StatisticsDirectionCase.b516StatisticsDirectionSide, "energy-index", b516StatisticsDirectionPointer(b516StatisticsDirectionIndex))
			b516StatisticsDirectionPeriodRule, b516StatisticsDirectionEnergyRule := "Select table row; ignore date.", "Ignore energy type and operation mode."
			if b516StatisticsDirectionIndex == 255 {
				b516StatisticsDirectionPeriodRule = "Select by date."
				b516StatisticsDirectionEnergyRule = "Select energy type and operation mode."
			}
			b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionPeriod.B516StatisticsDirectionRules, []string{"Total ignores period index and date.", b516StatisticsDirectionPeriodRule})
			b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionEnergy.B516StatisticsDirectionRules, []string{b516StatisticsDirectionEnergyRule})
			b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionPeriod.B516StatisticsDirectionUnknown, []string{"date packing", "offsets", "defaults"})
			b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionPeriod.B516StatisticsDirectionEncoding, (*string)(nil))
		}
		if b516StatisticsDirectionCase.b516StatisticsDirectionSide == "positive" {
			b516StatisticsDirectionReturn := b516StatisticsDirectionGet(b516StatisticsDirectionT, b516StatisticsDirectionCase.b516StatisticsDirectionVariant, 16, "positive", "return-nibble", nil)
			b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionReturn.B516StatisticsDirectionValues, map[string]string{"0": "value-contained", "1": "NOK"})
			b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionReturn.B516StatisticsDirectionUnknown, []string{"meanings 2..15", "offsets"})
		}
	}
}
func TestAb72Index(b516StatisticsDirectionT *b516StatisticsDirectionTesting.T) {
	b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionGet(b516StatisticsDirectionT, "ab72-statistics-index", 17, "request", "layout", nil).B516StatisticsDirectionLength, b516StatisticsDirectionPointer(uint8(0)))
	b516StatisticsDirectionFields := []string{"Version", "Build", "TotalEnergies", "Days", "DailyEnergies", "Months", "MonthlyEnergies", "Years", "YearlyEnergies"}
	b516StatisticsDirectionOrder := b516StatisticsDirectionGet(b516StatisticsDirectionT, "ab72-statistics-index", 17, "positive", "ordered-fields", nil)
	b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionOrder.B516StatisticsDirectionOrder, b516StatisticsDirectionFields)
	b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionOrder.B516StatisticsDirectionRules, []string{"Counts are returned metadata, not installed counts or enforced bounds."})
	for b516StatisticsDirectionIndex, b516StatisticsDirectionField := range b516StatisticsDirectionFields {
		b516StatisticsDirectionEntry := b516StatisticsDirectionGet(b516StatisticsDirectionT, "ab72-statistics-index", 17, "positive", b516StatisticsDirectionField, nil)
		b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionEntry.B516StatisticsDirectionPosition, b516StatisticsDirectionPointer(uint8(b516StatisticsDirectionIndex)))
		b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionEntry.B516StatisticsDirectionWidth, b516StatisticsDirectionPointer(uint8(8)))
		b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionEntry.B516StatisticsDirectionEncoding, b516StatisticsDirectionPointer("u8"))
	}
}
func TestAb72MeterClearAlternatives(b516StatisticsDirectionT *b516StatisticsDirectionTesting.T) {
	for _, b516StatisticsDirectionVariant := range []string{"ab72-meter-scalar-u8", "ab72-meter-low-six-bits", "ab72-clear-scalar-u8", "ab72-clear-low-six-bits"} {
		b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionGet(b516StatisticsDirectionT, b516StatisticsDirectionVariant, 18, "request", "flags", nil).B516StatisticsDirectionBits, map[string]uint8{"total": 0, "daily": 1, "monthly": 2, "yearly": 3})
		b516StatisticsDirectionSelector := b516StatisticsDirectionGet(b516StatisticsDirectionT, b516StatisticsDirectionVariant, 18, "request", "energy-selector", nil)
		b516StatisticsDirectionExpected := map[string]uint8{}
		b516StatisticsDirectionRules := []string{"Scalar u8 selector."}
		if b516StatisticsDirectionStrings.Contains(b516StatisticsDirectionVariant, "low-six") {
			b516StatisticsDirectionExpected = map[string]uint8{"sum-6": 6, "sum-7": 7}
			b516StatisticsDirectionRules = []string{"Low six bits selector; independent sum bits 6 and 7."}
		}
		b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionSelector.B516StatisticsDirectionBits, b516StatisticsDirectionExpected)
		b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionSelector.B516StatisticsDirectionRules, b516StatisticsDirectionRules)
		b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionGet(b516StatisticsDirectionT, b516StatisticsDirectionVariant, 18, "request", "mode", nil).B516StatisticsDirectionEncoding, b516StatisticsDirectionPointer("u8"))
		b516StatisticsDirectionPositive := b516StatisticsDirectionGet(b516StatisticsDirectionT, b516StatisticsDirectionVariant, 18, "positive", "flags", nil)
		b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionPositive.B516StatisticsDirectionBits, map[string]uint8{"not-supported": 0, "energy-unavailable": 1, "mode-unavailable": 2})
		b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionPositive.B516StatisticsDirectionConflicts, []string{"Nested and outer duplicate flags; packing precedence Unknown."})
		b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionPositive.B516StatisticsDirectionUnknown, []string{"packing precedence", "offsets", "defaults", "full variant-specific meter/clear tables"})
	}
}
func TestAb72YieldAndConflicts(b516StatisticsDirectionT *b516StatisticsDirectionTesting.T) {
	for _, b516StatisticsDirectionCase := range []struct{ b516StatisticsDirectionVariant, b516StatisticsDirectionSide string }{{"ab72-yield-read", "positive"}, {"ab72-yield-write", "request"}} {
		b516StatisticsDirectionEntry := b516StatisticsDirectionGet(b516StatisticsDirectionT, b516StatisticsDirectionCase.b516StatisticsDirectionVariant, 24, b516StatisticsDirectionCase.b516StatisticsDirectionSide, "yield-value", nil)
		b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionEntry.B516StatisticsDirectionEncoding, b516StatisticsDirectionPointer("LEu32"))
		b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionEntry.B516StatisticsDirectionUnit, b516StatisticsDirectionPointer("kWh"))
		b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionEntry.B516StatisticsDirectionWidth, b516StatisticsDirectionPointer(uint8(32)))
		b516StatisticsDirectionRaw := []byte{255, 255, 255, 255}
		b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionBinary.LittleEndian.Uint32(b516StatisticsDirectionRaw), uint32(4294967295))
		b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionEntry.B516StatisticsDirectionRules, []string{"All u32 values retain their raw numeric meaning; no sentinel inferred."})
	}
	for _, b516StatisticsDirectionVariant := range []string{"ab72-yield-read", "ab72-yield-write"} {
		for _, b516StatisticsDirectionSide := range []string{"request", "positive"} {
			b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionGet(b516StatisticsDirectionT, b516StatisticsDirectionVariant, 24, b516StatisticsDirectionSide, "yield-type", nil).B516StatisticsDirectionValues, map[string]string{"0": "solar", "1": "environment", "2": "energy", "3": "gas-or-oil"})
		}
	}
	b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionGet(b516StatisticsDirectionT, "ab72-yield-reset", 25, "request", "yield-type", nil).B516StatisticsDirectionValues, map[string]string{"0": "solar", "1": "environment", "2": "energy", "3": "gas-or-oil"})
	b516StatisticsDirectionFeedback := b516StatisticsDirectionGet(b516StatisticsDirectionT, "ab72-yield-reset", 25, "positive", "feedback", nil)
	b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionFeedback.B516StatisticsDirectionValues, map[string]string{"1": "accepted"})
	b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionFeedback.B516StatisticsDirectionUnknown, []string{"other feedback meanings", "effects", "persistence"})
	for _, b516StatisticsDirectionField := range []string{"yield-type", "month", "year-index", "mode", "sum"} {
		b516StatisticsDirectionEntry := b516StatisticsDirectionGet(b516StatisticsDirectionT, "ab72-yield-month", 26, "positive", b516StatisticsDirectionField, nil)
		b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionEntry.B516StatisticsDirectionPosition, (*uint8)(nil))
		b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionEntry.B516StatisticsDirectionUnit, (*string)(nil))
		b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionEntry.B516StatisticsDirectionConflicts, []string{"Unit-type overlaps second sum byte.", "Month 10/11 conflict."})
		b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionEntry.B516StatisticsDirectionUnknown, []string{"offsets", "defaults", "packing precedence", "unit"})
	}
	b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionGet(b516StatisticsDirectionT, "ab72-yield-month", 26, "positive", "sum", nil).B516StatisticsDirectionEncoding, b516StatisticsDirectionPointer("LEu32"))
	b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionGet(b516StatisticsDirectionT, "ab72-yield-month", 26, "positive", "year-index", nil).B516StatisticsDirectionValues, map[string]string{"0": "this-year", "1": "last-year"})
	b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionGet(b516StatisticsDirectionT, "ab72-yield-month", 26, "positive", "mode", nil).B516StatisticsDirectionValues, map[string]string{"0": "all", "1": "heating", "2": "DHW"})
}
func TestAb72ExactKeys(b516StatisticsDirectionT *b516StatisticsDirectionTesting.T) {
	b516StatisticsDirectionBase := b516StatisticsDirectionGet(b516StatisticsDirectionT, "ab72-statistics-read", 16, "request", "period-index", b516StatisticsDirectionPointer(uint8(0))).B516StatisticsDirectionKey
	for _, b516StatisticsDirectionChange := range []func(*B516StatisticsDirectionKey){
		func(b516StatisticsDirectionKey *B516StatisticsDirectionKey) {
			b516StatisticsDirectionKey.B516StatisticsDirectionFamily = "ab72-other"
		},
		func(b516StatisticsDirectionKey *B516StatisticsDirectionKey) {
			b516StatisticsDirectionKey.B516StatisticsDirectionVariant = "ab72-other"
		},
		func(b516StatisticsDirectionKey *B516StatisticsDirectionKey) {
			b516StatisticsDirectionKey.B516StatisticsDirectionOperation = b516StatisticsDirectionPointer(uint8(0))
		},
		func(b516StatisticsDirectionKey *B516StatisticsDirectionKey) {
			b516StatisticsDirectionKey.B516StatisticsDirectionOperation = nil
		},
		func(b516StatisticsDirectionKey *B516StatisticsDirectionKey) {
			b516StatisticsDirectionKey.B516StatisticsDirection = "ab72-other"
		},
		func(b516StatisticsDirectionKey *B516StatisticsDirectionKey) {
			b516StatisticsDirectionKey.B516StatisticsDirectionSelector = nil
		},
		func(b516StatisticsDirectionKey *B516StatisticsDirectionKey) {
			b516StatisticsDirectionKey.B516StatisticsDirectionSelector = b516StatisticsDirectionPointer(uint8(1))
		},
		func(b516StatisticsDirectionKey *B516StatisticsDirectionKey) {
			b516StatisticsDirectionKey.B516StatisticsDirectionField = "ab72-other"
		},
	} {
		b516StatisticsDirectionKey := b516StatisticsDirectionBase
		b516StatisticsDirectionChange(&b516StatisticsDirectionKey)
		_, b516StatisticsDirectionFound, b516StatisticsDirectionErr := B516StatisticsDirectionLookup(b516StatisticsDirectionKey)
		b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionErr, nil)
		b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionFound, false)
	}
	b516StatisticsDirectionEntries, b516StatisticsDirectionErr := B516StatisticsDirectionEntries()
	b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionErr, nil)
	b516StatisticsDirectionEqual(b516StatisticsDirectionT, len(b516StatisticsDirectionEntries), 68)
	for _, b516StatisticsDirectionEntry := range b516StatisticsDirectionEntries {
		b516StatisticsDirectionGot, b516StatisticsDirectionFound, b516StatisticsDirectionErr := B516StatisticsDirectionLookup(b516StatisticsDirectionEntry.B516StatisticsDirectionKey)
		b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionErr, nil)
		b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionFound, true)
		b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionGot, b516StatisticsDirectionEntry)
	}
}
func TestAb72ParserAndOwnership(b516StatisticsDirectionT *b516StatisticsDirectionTesting.T) {
	b516StatisticsDirectionEntry := b516StatisticsDirectionGet(b516StatisticsDirectionT, "ab72-statistics-write", 16, "request", "flags", nil)
	for _, b516StatisticsDirectionRaw := range []string{"{", "null", "[]", "[{}]", string(b516StatisticsDirectionData) + " []", b516StatisticsDirectionStrings.Replace(string(b516StatisticsDirectionData), "\"grade\":", "\"extra\": 0, \"grade\":", 1), b516StatisticsDirectionStrings.Replace(string(b516StatisticsDirectionData), "Hypothesis", "Proven", 1)} {
		_, b516StatisticsDirectionErr := b516StatisticsDirectionParse([]byte(b516StatisticsDirectionRaw))
		if b516StatisticsDirectionErr == nil {
			b516StatisticsDirectionT.Fatal("invalid input admitted")
		}
	}
	b516StatisticsDirectionRaw, b516StatisticsDirectionErr := b516StatisticsDirectionJSON.Marshal([]B516StatisticsDirectionEntry{b516StatisticsDirectionEntry, b516StatisticsDirectionEntry})
	b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionErr, nil)
	_, b516StatisticsDirectionErr = b516StatisticsDirectionParse(b516StatisticsDirectionRaw)
	if b516StatisticsDirectionErr == nil {
		b516StatisticsDirectionT.Fatal("duplicate admitted")
	}
	b516StatisticsDirectionRaw, b516StatisticsDirectionErr = b516StatisticsDirectionJSON.Marshal([]B516StatisticsDirectionEntry{b516StatisticsDirectionEntry})
	b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionErr, nil)
	b516StatisticsDirectionParsed, b516StatisticsDirectionErr := b516StatisticsDirectionParse(b516StatisticsDirectionRaw)
	b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionErr, nil)
	for b516StatisticsDirectionIndex := range b516StatisticsDirectionRaw {
		b516StatisticsDirectionRaw[b516StatisticsDirectionIndex] = 0
	}
	b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionParsed[0], b516StatisticsDirectionEntry)
	b516StatisticsDirectionEntries, b516StatisticsDirectionErr := B516StatisticsDirectionEntries()
	b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionErr, nil)
	for b516StatisticsDirectionIndex := range b516StatisticsDirectionEntries {
		b516StatisticsDirectionOwned := &b516StatisticsDirectionEntries[b516StatisticsDirectionIndex]
		b516StatisticsDirectionOwned.B516StatisticsDirectionKey.B516StatisticsDirectionFamily = "ab72-mutated"
		*b516StatisticsDirectionOwned.B516StatisticsDirectionKey.B516StatisticsDirectionOperation = 0
		if b516StatisticsDirectionOwned.B516StatisticsDirectionKey.B516StatisticsDirectionSelector != nil {
			*b516StatisticsDirectionOwned.B516StatisticsDirectionKey.B516StatisticsDirectionSelector = 42
		}
		for b516StatisticsDirectionKey := range b516StatisticsDirectionOwned.B516StatisticsDirectionBits {
			b516StatisticsDirectionOwned.B516StatisticsDirectionBits[b516StatisticsDirectionKey] = 7
		}
		for b516StatisticsDirectionKey := range b516StatisticsDirectionOwned.B516StatisticsDirectionValues {
			b516StatisticsDirectionOwned.B516StatisticsDirectionValues[b516StatisticsDirectionKey] = "ab72-mutated"
		}
		for _, b516StatisticsDirectionSlice := range [][]string{b516StatisticsDirectionOwned.B516StatisticsDirectionOrder, b516StatisticsDirectionOwned.B516StatisticsDirectionRules, b516StatisticsDirectionOwned.B516StatisticsDirectionConflicts, b516StatisticsDirectionOwned.B516StatisticsDirectionUnknown} {
			for b516StatisticsDirectionIndex := range b516StatisticsDirectionSlice {
				b516StatisticsDirectionSlice[b516StatisticsDirectionIndex] = "ab72-mutated"
			}
		}
		if b516StatisticsDirectionOwned.B516StatisticsDirectionLength != nil {
			*b516StatisticsDirectionOwned.B516StatisticsDirectionLength = 255
		}
		if b516StatisticsDirectionOwned.B516StatisticsDirectionPosition != nil {
			*b516StatisticsDirectionOwned.B516StatisticsDirectionPosition = 255
		}
		if b516StatisticsDirectionOwned.B516StatisticsDirectionWidth != nil {
			*b516StatisticsDirectionOwned.B516StatisticsDirectionWidth = 1
		}
		if b516StatisticsDirectionOwned.B516StatisticsDirectionEncoding != nil {
			*b516StatisticsDirectionOwned.B516StatisticsDirectionEncoding = "ab72-mutated"
		}
		if b516StatisticsDirectionOwned.B516StatisticsDirectionUnit != nil {
			*b516StatisticsDirectionOwned.B516StatisticsDirectionUnit = "ab72-mutated"
		}
	}
	b516StatisticsDirectionFresh, b516StatisticsDirectionErr := B516StatisticsDirectionEntries()
	b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionErr, nil)
	b516StatisticsDirectionNative, b516StatisticsDirectionErr := b516StatisticsDirectionParse(b516StatisticsDirectionData)
	b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionErr, nil)
	b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionFresh, b516StatisticsDirectionNative)
	b516StatisticsDirectionLookup := b516StatisticsDirectionGet(b516StatisticsDirectionT, "ab72-statistics-write", 16, "request", "flags", nil)
	b516StatisticsDirectionLookup.B516StatisticsDirectionBits["access"] = 7
	*b516StatisticsDirectionLookup.B516StatisticsDirectionKey.B516StatisticsDirectionOperation = 0
	b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionGet(b516StatisticsDirectionT, "ab72-statistics-write", 16, "request", "flags", nil), b516StatisticsDirectionEntry)
	for _, b516StatisticsDirectionOwned := range b516StatisticsDirectionFresh {
		b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionOwned.B516StatisticsDirectionGrade, "Hypothesis")
		b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionOwned.B516StatisticsDirectionSoftwareContract, "Proven")
		for _, b516StatisticsDirectionGrade := range []string{b516StatisticsDirectionOwned.B516StatisticsDirectionPhysicalSupport, b516StatisticsDirectionOwned.B516StatisticsDirectionPresence, b516StatisticsDirectionOwned.B516StatisticsDirectionFreshness, b516StatisticsDirectionOwned.B516StatisticsDirectionExecution, b516StatisticsDirectionOwned.B516StatisticsDirectionEffects, b516StatisticsDirectionOwned.B516StatisticsDirectionPersistence, b516StatisticsDirectionOwned.B516StatisticsDirectionDelivery} {
			b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionGrade, "Unknown")
		}
	}
}

func TestAb72FailClosed(b516StatisticsDirectionT *b516StatisticsDirectionTesting.T) {
	b516StatisticsDirectionOriginal := b516StatisticsDirectionData
	defer func() { b516StatisticsDirectionData = b516StatisticsDirectionOriginal }()
	b516StatisticsDirectionValid := b516StatisticsDirectionGet(b516StatisticsDirectionT, "ab72-statistics-read", 16, "request", "flags", nil)
	b516StatisticsDirectionInvalid := b516StatisticsDirectionValid
	b516StatisticsDirectionInvalid.B516StatisticsDirectionKey.B516StatisticsDirectionField = ""
	b516StatisticsDirectionRaw, b516StatisticsDirectionErr := b516StatisticsDirectionJSON.Marshal([]B516StatisticsDirectionEntry{b516StatisticsDirectionValid, b516StatisticsDirectionInvalid})
	b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionErr, nil)
	b516StatisticsDirectionData = b516StatisticsDirectionRaw
	b516StatisticsDirectionEntries, b516StatisticsDirectionErr := B516StatisticsDirectionEntries()
	if b516StatisticsDirectionErr == nil || b516StatisticsDirectionEntries != nil {
		b516StatisticsDirectionT.Fatal("partial catalog admitted")
	}
	b516StatisticsDirectionEntry, b516StatisticsDirectionFound, b516StatisticsDirectionErr := B516StatisticsDirectionLookup(b516StatisticsDirectionValid.B516StatisticsDirectionKey)
	if b516StatisticsDirectionErr == nil || b516StatisticsDirectionFound {
		b516StatisticsDirectionT.Fatal("lookup ignored catalog error")
	}
	b516StatisticsDirectionEqual(b516StatisticsDirectionT, b516StatisticsDirectionEntry, B516StatisticsDirectionEntry{})
}
