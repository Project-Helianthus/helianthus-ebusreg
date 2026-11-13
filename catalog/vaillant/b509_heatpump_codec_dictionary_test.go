package vaillant

import (
	b509HeatpumpCodecDictionaryBytes "bytes"
	b509HeatpumpCodecDictionaryJSONPackage "encoding/json"
	b509HeatpumpCodecDictionaryReflect "reflect"
	b509HeatpumpCodecDictionaryStrings "strings"
	b509HeatpumpCodecDictionaryTesting "testing"
)

func b509HeatpumpCodecDictionaryPointer[b509HeatpumpCodecDictionaryT any](b509HeatpumpCodecDictionaryValue b509HeatpumpCodecDictionaryT) *b509HeatpumpCodecDictionaryT {
	return &b509HeatpumpCodecDictionaryValue
}

func b509HeatpumpCodecDictionaryKey(b509HeatpumpCodecDictionaryFirst, b509HeatpumpCodecDictionarySecond uint16) B509HeatpumpCodecDictionaryKey {
	return B509HeatpumpCodecDictionaryKey{"heat-pump", "direct-symbolic-candidate", b509HeatpumpCodecDictionaryPointer(byte(181)), b509HeatpumpCodecDictionaryPointer(byte(9)), b509HeatpumpCodecDictionaryPointer(byte(13)), "read", []uint16{b509HeatpumpCodecDictionaryFirst, b509HeatpumpCodecDictionarySecond}}
}

func TestAb70ExactDictionary(b509HeatpumpCodecDictionaryT *b509HeatpumpCodecDictionaryTesting.T) {
	b509HeatpumpCodecDictionaryExpected := []struct {
		b509HeatpumpCodecDictionaryFirst    uint16
		b509HeatpumpCodecDictionarySecond   uint16
		b509HeatpumpCodecDictionaryQuantity string
		b509HeatpumpCodecDictionaryCodec    string
		b509HeatpumpCodecDictionaryUnit     string
		b509HeatpumpCodecDictionaryConflict string
	}{
		{0xD0, 0, "heat_pump_status", "hpstatus enum", "", ""},
		{0xBC, 0, "environmental_energy_yield_total", "energy4", "kWh", "Environmental versus delivered yield remains unresolved."},
		{0xBA, 0, "environmental_power", "power", "kW", ""},
		{0xBB, 0, "environmental_power_ratio", "percent0", "%", ""},
		{0x8D, 0, "refrigerant_superheat", "temp", "K delta", ""},
		{0x8E, 0, "refrigerant_subcooling", "temp", "K delta", ""},
		{0x0F, 0, "source_inlet_temperature", "tempsensor", "C", ""},
		{0x14, 0, "refrigerant_high_pressure", "sensor", "bar", ""},
		{0x15, 0, "refrigerant_low_pressure", "sensor", "bar", ""},
		{0xE4, 0, "compressor_state", "UCH enum", "", ""},
		{0xE9, 0, "compressor_control_state", "UCH enum", "", ""},
		{0xE6, 0, "compressor_minimum_on_duration", "seconds", "s", ""},
		{0xE7, 0, "compressor_minimum_off_duration", "seconds", "s", ""},
		{0xE8, 0, "compressor_minimum_start_interval", "seconds", "s", ""},
		{0xA8, 0, "compressor_dhw_start_count", "counter", "count", ""},
		{0xB4, 0, "compressor_high_pressure_cutout_count", "counter", "count", ""},
		{0xCE, 1, "defrost_active", "yesno", "bool", ""},
		{0x57, 1, "phase_order_status", "phaseok enum", "", ""},
		{0x88, 0, "phase_status_flags", "phase flags", "bitmask", ""},
		{0x3F, 0, "inrush_current_limiter_output_active", "on/off", "", ""},
		{0xA9, 1, "noise_reduction_active", "yesno", "bool", ""},
		{0x24, 1, "noise_reduction_factor", "percent0", "", "Ratio unit remains Unknown."},
		{0xE2, 1, "environmental_power", "UCH", "W", ""},
		{0xE3, 1, "bivalent_mode", "UCH enum", "", ""},
		{0xE4, 1, "source_temperature_limit_at_maximum_flow", "SCH", "C", ""},
		{0xE5, 1, "source_temperature_limit_slope", "UCH", "slope", ""},
		{0xE9, 1, "refrigerant_low_pressure_maximum", "ULG", "", "Pressure unit remains Unknown."},
		{0x80, 0, "compressor_energy_integral", "value", "", "Integral codec and unit remain Unknown."},
	}
	b509HeatpumpCodecDictionaryEntries, b509HeatpumpCodecDictionaryError := B509HeatpumpCodecDictionaryEntries()
	if b509HeatpumpCodecDictionaryError != nil || len(b509HeatpumpCodecDictionaryEntries) != len(b509HeatpumpCodecDictionaryExpected) {
		b509HeatpumpCodecDictionaryT.Fatalf("ab70 exact count: %d %v", len(b509HeatpumpCodecDictionaryEntries), b509HeatpumpCodecDictionaryError)
	}
	b509HeatpumpCodecDictionarySeen := make(map[[2]uint16]bool)
	for _, b509HeatpumpCodecDictionaryWant := range b509HeatpumpCodecDictionaryExpected {
		b509HeatpumpCodecDictionaryKey := b509HeatpumpCodecDictionaryKey(b509HeatpumpCodecDictionaryWant.b509HeatpumpCodecDictionaryFirst, b509HeatpumpCodecDictionaryWant.b509HeatpumpCodecDictionarySecond)
		b509HeatpumpCodecDictionaryGot, b509HeatpumpCodecDictionaryFound, b509HeatpumpCodecDictionaryError := B509HeatpumpCodecDictionaryLookup(b509HeatpumpCodecDictionaryKey)
		if b509HeatpumpCodecDictionaryError != nil || !b509HeatpumpCodecDictionaryFound {
			b509HeatpumpCodecDictionaryT.Fatalf("ab70 missing %v: %v", b509HeatpumpCodecDictionaryKey, b509HeatpumpCodecDictionaryError)
		}
		b509HeatpumpCodecDictionaryEntry := B509HeatpumpCodecDictionaryEntry{
			B509HeatpumpCodecDictionaryKey: b509HeatpumpCodecDictionaryKey, B509HeatpumpCodecDictionaryQuantity: b509HeatpumpCodecDictionaryWant.b509HeatpumpCodecDictionaryQuantity, B509HeatpumpCodecDeclaration: b509HeatpumpCodecDictionaryWant.b509HeatpumpCodecDictionaryCodec,
			B509HeatpumpCodecDictionaryGrade: "Hypothesis", B509HeatpumpCodecDictionarySoftwareContract: "Proven", B509HeatpumpCodecDictionaryQualification: "Unknown", B509HeatpumpCodecDictionaryPresence: "Unknown", B509HeatpumpCodecDictionaryFreshness: "Unknown", B509HeatpumpCodecDictionaryExecution: "Unknown", B509HeatpumpCodecDictionaryPhysicalEffects: "Unknown", B509HeatpumpCodecDictionaryPersistence: "Unknown", B509HeatpumpCodecDictionaryDelivery: "Unknown",
			B509HeatpumpCodecDictionaryNotes: []string{"Literal codec declaration; no executable converter."},
		}
		if b509HeatpumpCodecDictionaryWant.b509HeatpumpCodecDictionaryUnit != "" {
			b509HeatpumpCodecDictionaryEntry.B509HeatpumpCodecDictionaryUnit = b509HeatpumpCodecDictionaryPointer(b509HeatpumpCodecDictionaryWant.b509HeatpumpCodecDictionaryUnit)
		}
		if b509HeatpumpCodecDictionaryWant.b509HeatpumpCodecDictionaryFirst == 0xE2 {
			b509HeatpumpCodecDictionaryEntry.B509HeatpumpCodecDictionaryWidthBits = b509HeatpumpCodecDictionaryPointer(8)
			b509HeatpumpCodecDictionaryEntry.B509HeatpumpCodecDictionarySignedness = b509HeatpumpCodecDictionaryPointer("unsigned")
			b509HeatpumpCodecDictionaryEntry.B509HeatpumpCodecDictionaryDeclaredScale = b509HeatpumpCodecDictionaryPointer(100)
		}
		if b509HeatpumpCodecDictionaryWant.b509HeatpumpCodecDictionaryFirst == 0xE4 && b509HeatpumpCodecDictionaryWant.b509HeatpumpCodecDictionarySecond == 1 {
			b509HeatpumpCodecDictionaryEntry.B509HeatpumpCodecDictionarySignedness = b509HeatpumpCodecDictionaryPointer("signed")
		}
		if b509HeatpumpCodecDictionaryWant.b509HeatpumpCodecDictionaryConflict != "" {
			b509HeatpumpCodecDictionaryEntry.B509HeatpumpCodecDictionaryNotes = append(b509HeatpumpCodecDictionaryEntry.B509HeatpumpCodecDictionaryNotes, b509HeatpumpCodecDictionaryWant.b509HeatpumpCodecDictionaryConflict)
		}
		if !b509HeatpumpCodecDictionaryReflect.DeepEqual(b509HeatpumpCodecDictionaryEntry, b509HeatpumpCodecDictionaryGot) {
			b509HeatpumpCodecDictionaryT.Errorf("ab70 selector %v got %#v want %#v", b509HeatpumpCodecDictionaryKey.B509HeatpumpCodecDictionarySelector, b509HeatpumpCodecDictionaryGot, b509HeatpumpCodecDictionaryEntry)
		}
		b509HeatpumpCodecDictionarySeen[[2]uint16{b509HeatpumpCodecDictionaryWant.b509HeatpumpCodecDictionaryFirst, b509HeatpumpCodecDictionaryWant.b509HeatpumpCodecDictionarySecond}] = true
	}
	for _, b509HeatpumpCodecDictionaryEntry := range b509HeatpumpCodecDictionaryEntries {
		if len(b509HeatpumpCodecDictionaryEntry.B509HeatpumpCodecDictionaryKey.B509HeatpumpCodecDictionarySelector) != 2 || !b509HeatpumpCodecDictionarySeen[[2]uint16{b509HeatpumpCodecDictionaryEntry.B509HeatpumpCodecDictionaryKey.B509HeatpumpCodecDictionarySelector[0], b509HeatpumpCodecDictionaryEntry.B509HeatpumpCodecDictionaryKey.B509HeatpumpCodecDictionarySelector[1]}] {
			b509HeatpumpCodecDictionaryT.Fatal("ab70 extra tuple")
		}
		delete(b509HeatpumpCodecDictionarySeen, [2]uint16{b509HeatpumpCodecDictionaryEntry.B509HeatpumpCodecDictionaryKey.B509HeatpumpCodecDictionarySelector[0], b509HeatpumpCodecDictionaryEntry.B509HeatpumpCodecDictionaryKey.B509HeatpumpCodecDictionarySelector[1]})
	}
}

func TestAb70KeyAxes(b509HeatpumpCodecDictionaryT *b509HeatpumpCodecDictionaryTesting.T) {
	b509HeatpumpCodecDictionaryChanges := []struct {
		b509HeatpumpCodecDictionaryName   string
		b509HeatpumpCodecDictionaryChange func(*B509HeatpumpCodecDictionaryKey)
	}{
		{"ab70Family", func(b509HeatpumpCodecDictionaryK *B509HeatpumpCodecDictionaryKey) {
			b509HeatpumpCodecDictionaryK.B509HeatpumpCodecDictionaryFamily = "heat-pump-interface-unit"
		}},
		{"ab70Variant", func(b509HeatpumpCodecDictionaryK *B509HeatpumpCodecDictionaryKey) {
			b509HeatpumpCodecDictionaryK.B509HeatpumpCodecDictionaryVariant = "unknown"
		}},
		{"ab70Primary", func(b509HeatpumpCodecDictionaryK *B509HeatpumpCodecDictionaryKey) {
			b509HeatpumpCodecDictionaryK.B509HeatpumpCodecDictionaryPrimary = b509HeatpumpCodecDictionaryPointer(byte(0))
		}},
		{"ab70Secondary", func(b509HeatpumpCodecDictionaryK *B509HeatpumpCodecDictionaryKey) {
			b509HeatpumpCodecDictionaryK.B509HeatpumpCodecDictionarySecondary = b509HeatpumpCodecDictionaryPointer(byte(0))
		}},
		{"ab70Operation", func(b509HeatpumpCodecDictionaryK *B509HeatpumpCodecDictionaryKey) {
			b509HeatpumpCodecDictionaryK.B509HeatpumpCodecDictionaryOperation = b509HeatpumpCodecDictionaryPointer(byte(14))
		}},
		{"ab70Alternate", func(b509HeatpumpCodecDictionaryK *B509HeatpumpCodecDictionaryKey) {
			b509HeatpumpCodecDictionaryK.B509HeatpumpCodecDictionaryOperation = b509HeatpumpCodecDictionaryPointer(byte(41))
		}},
		{"ab70Request", func(b509HeatpumpCodecDictionaryK *B509HeatpumpCodecDictionaryKey) {
			b509HeatpumpCodecDictionaryK.B509HeatpumpCodecDictionaryDirection = "request"
		}},
		{"ab70Positive", func(b509HeatpumpCodecDictionaryK *B509HeatpumpCodecDictionaryKey) {
			b509HeatpumpCodecDictionaryK.B509HeatpumpCodecDictionaryDirection = "positive"
		}},
		{"ab70First", func(b509HeatpumpCodecDictionaryK *B509HeatpumpCodecDictionaryKey) {
			b509HeatpumpCodecDictionaryK.B509HeatpumpCodecDictionarySelector[0] = 0
		}},
		{"ab70Second", func(b509HeatpumpCodecDictionaryK *B509HeatpumpCodecDictionaryKey) {
			b509HeatpumpCodecDictionaryK.B509HeatpumpCodecDictionarySelector[1] = 1
		}},
		{"ab70Reverse", func(b509HeatpumpCodecDictionaryK *B509HeatpumpCodecDictionaryKey) {
			b509HeatpumpCodecDictionaryK.B509HeatpumpCodecDictionarySelector = []uint16{0, 0xD0}
		}},
		{"ab70AbsentPrimary", func(b509HeatpumpCodecDictionaryK *B509HeatpumpCodecDictionaryKey) {
			b509HeatpumpCodecDictionaryK.B509HeatpumpCodecDictionaryPrimary = nil
		}},
		{"ab70AbsentSecondary", func(b509HeatpumpCodecDictionaryK *B509HeatpumpCodecDictionaryKey) {
			b509HeatpumpCodecDictionaryK.B509HeatpumpCodecDictionarySecondary = nil
		}},
		{"ab70AbsentOperation", func(b509HeatpumpCodecDictionaryK *B509HeatpumpCodecDictionaryKey) {
			b509HeatpumpCodecDictionaryK.B509HeatpumpCodecDictionaryOperation = nil
		}},
		{"ab70AbsentSelector", func(b509HeatpumpCodecDictionaryK *B509HeatpumpCodecDictionaryKey) {
			b509HeatpumpCodecDictionaryK.B509HeatpumpCodecDictionarySelector = nil
		}},
		{"ab70PartialSelector", func(b509HeatpumpCodecDictionaryK *B509HeatpumpCodecDictionaryKey) {
			b509HeatpumpCodecDictionaryK.B509HeatpumpCodecDictionarySelector = []uint16{0xD0}
		}},
		{"ab70ExtraSelector", func(b509HeatpumpCodecDictionaryK *B509HeatpumpCodecDictionaryKey) {
			b509HeatpumpCodecDictionaryK.B509HeatpumpCodecDictionarySelector = []uint16{0xD0, 0, 0}
		}},
	}
	b509HeatpumpCodecDictionaryOriginal, _, _ := B509HeatpumpCodecDictionaryLookup(b509HeatpumpCodecDictionaryKey(0xD0, 0))
	for _, b509HeatpumpCodecDictionaryCase := range b509HeatpumpCodecDictionaryChanges {
		b509HeatpumpCodecDictionaryT.Run(b509HeatpumpCodecDictionaryCase.b509HeatpumpCodecDictionaryName, func(b509HeatpumpCodecDictionaryT *b509HeatpumpCodecDictionaryTesting.T) {
			b509HeatpumpCodecDictionaryChanged := b509HeatpumpCodecDictionaryKey(0xD0, 0)
			b509HeatpumpCodecDictionaryCase.b509HeatpumpCodecDictionaryChange(&b509HeatpumpCodecDictionaryChanged)
			if _, b509HeatpumpCodecDictionaryFound, b509HeatpumpCodecDictionaryError := B509HeatpumpCodecDictionaryLookup(b509HeatpumpCodecDictionaryChanged); b509HeatpumpCodecDictionaryFound || b509HeatpumpCodecDictionaryError != nil {
				b509HeatpumpCodecDictionaryT.Fatalf("ab70 admitted altered key: %v", b509HeatpumpCodecDictionaryError)
			}
			b509HeatpumpCodecDictionaryControl := b509HeatpumpCodecDictionaryOriginal
			b509HeatpumpCodecDictionaryControl.B509HeatpumpCodecDictionaryKey = b509HeatpumpCodecDictionaryChanged
			if _, b509HeatpumpCodecDictionaryFound, b509HeatpumpCodecDictionaryError := b509HeatpumpCodecDictionaryFind([]B509HeatpumpCodecDictionaryEntry{b509HeatpumpCodecDictionaryControl}, b509HeatpumpCodecDictionaryChanged); !b509HeatpumpCodecDictionaryFound || b509HeatpumpCodecDictionaryError != nil {
				b509HeatpumpCodecDictionaryT.Fatal("ab70 exact changed control failed")
			}
			if _, b509HeatpumpCodecDictionaryFound, _ := b509HeatpumpCodecDictionaryFind([]B509HeatpumpCodecDictionaryEntry{b509HeatpumpCodecDictionaryControl}, b509HeatpumpCodecDictionaryKey(0xD0, 0)); b509HeatpumpCodecDictionaryFound {
				b509HeatpumpCodecDictionaryT.Fatal("ab70 omitted key axis")
			}
		})
	}
}

func TestAb70ParserAndOwnership(b509HeatpumpCodecDictionaryT *b509HeatpumpCodecDictionaryTesting.T) {
	b509HeatpumpCodecDictionaryEntries, b509HeatpumpCodecDictionaryError := B509HeatpumpCodecDictionaryEntries()
	if b509HeatpumpCodecDictionaryError != nil {
		b509HeatpumpCodecDictionaryT.Fatal(b509HeatpumpCodecDictionaryError)
	}
	b509HeatpumpCodecDictionaryEntry := b509HeatpumpCodecDictionaryEntries[0]
	b509HeatpumpCodecDictionaryEntry.B509HeatpumpCodecDeclaration = "ab70-unknown-label"
	b509HeatpumpCodecDictionaryEntry.B509HeatpumpCodecDictionaryMapping = map[string]string{"ab70-zero": "ab70-native"}
	b509HeatpumpCodecDictionaryEntry.B509HeatpumpCodecDictionaryWidthBits = b509HeatpumpCodecDictionaryPointer(0)
	b509HeatpumpCodecDictionaryEntry.B509HeatpumpCodecDictionarySentinel = b509HeatpumpCodecDictionaryPointer(0)
	b509HeatpumpCodecDictionaryEntry.B509HeatpumpCodecDictionaryUnit = b509HeatpumpCodecDictionaryPointer("")
	b509HeatpumpCodecDictionaryData, b509HeatpumpCodecDictionaryError := b509HeatpumpCodecDictionaryJSONPackage.Marshal([]B509HeatpumpCodecDictionaryEntry{b509HeatpumpCodecDictionaryEntry})
	if b509HeatpumpCodecDictionaryError != nil {
		b509HeatpumpCodecDictionaryT.Fatal(b509HeatpumpCodecDictionaryError)
	}
	b509HeatpumpCodecDictionaryOriginal := append([]byte(nil), b509HeatpumpCodecDictionaryData...)
	b509HeatpumpCodecDictionaryParsed, b509HeatpumpCodecDictionaryError := b509HeatpumpCodecDictionaryParse(b509HeatpumpCodecDictionaryData)
	if b509HeatpumpCodecDictionaryError != nil || !b509HeatpumpCodecDictionaryReflect.DeepEqual(b509HeatpumpCodecDictionaryParsed, []B509HeatpumpCodecDictionaryEntry{b509HeatpumpCodecDictionaryEntry}) || !b509HeatpumpCodecDictionaryBytes.Equal(b509HeatpumpCodecDictionaryData, b509HeatpumpCodecDictionaryOriginal) {
		b509HeatpumpCodecDictionaryT.Fatal("ab70 raw or metadata-only control failed", b509HeatpumpCodecDictionaryError)
	}
	b509HeatpumpCodecDictionaryCopy, b509HeatpumpCodecDictionaryFound, b509HeatpumpCodecDictionaryError := b509HeatpumpCodecDictionaryFind(b509HeatpumpCodecDictionaryParsed, b509HeatpumpCodecDictionaryEntry.B509HeatpumpCodecDictionaryKey)
	if !b509HeatpumpCodecDictionaryFound || b509HeatpumpCodecDictionaryError != nil {
		b509HeatpumpCodecDictionaryT.Fatal("ab70 copy failed")
	}
	b509HeatpumpCodecDictionaryCopy.B509HeatpumpCodecDictionaryKey.B509HeatpumpCodecDictionarySelector[0] = 0
	*b509HeatpumpCodecDictionaryCopy.B509HeatpumpCodecDictionaryKey.B509HeatpumpCodecDictionaryPrimary = 0
	*b509HeatpumpCodecDictionaryCopy.B509HeatpumpCodecDictionaryKey.B509HeatpumpCodecDictionarySecondary = 0
	*b509HeatpumpCodecDictionaryCopy.B509HeatpumpCodecDictionaryKey.B509HeatpumpCodecDictionaryOperation = 0
	*b509HeatpumpCodecDictionaryCopy.B509HeatpumpCodecDictionaryWidthBits = 9
	*b509HeatpumpCodecDictionaryCopy.B509HeatpumpCodecDictionarySentinel = 9
	*b509HeatpumpCodecDictionaryCopy.B509HeatpumpCodecDictionaryUnit = "ab70-changed"
	b509HeatpumpCodecDictionaryCopy.B509HeatpumpCodecDictionaryMapping["ab70-zero"] = "ab70-changed"
	b509HeatpumpCodecDictionaryCopy.B509HeatpumpCodecDictionaryNotes[0] = "ab70-changed"
	if !b509HeatpumpCodecDictionaryReflect.DeepEqual(b509HeatpumpCodecDictionaryParsed, []B509HeatpumpCodecDictionaryEntry{b509HeatpumpCodecDictionaryEntry}) {
		b509HeatpumpCodecDictionaryT.Fatal("ab70 nested data shared")
	}
	for b509HeatpumpCodecDictionaryIndex := range b509HeatpumpCodecDictionaryData {
		b509HeatpumpCodecDictionaryData[b509HeatpumpCodecDictionaryIndex] = ' '
	}
	if !b509HeatpumpCodecDictionaryReflect.DeepEqual(b509HeatpumpCodecDictionaryParsed, []B509HeatpumpCodecDictionaryEntry{b509HeatpumpCodecDictionaryEntry}) {
		b509HeatpumpCodecDictionaryT.Fatal("ab70 input memory shared")
	}
	b509HeatpumpCodecDictionaryEntries[0].B509HeatpumpCodecDictionaryNotes[0] = "ab70-changed"
	b509HeatpumpCodecDictionaryEntries[0].B509HeatpumpCodecDictionaryKey.B509HeatpumpCodecDictionarySelector[0] = 0
	b509HeatpumpCodecDictionaryAgain, _, _ := B509HeatpumpCodecDictionaryLookup(b509HeatpumpCodecDictionaryKey(0xD0, 0))
	if b509HeatpumpCodecDictionaryAgain.B509HeatpumpCodecDictionaryNotes[0] != "Literal codec declaration; no executable converter." {
		b509HeatpumpCodecDictionaryT.Fatal("ab70 entries shared")
	}
	b509HeatpumpCodecDictionaryAgain.B509HeatpumpCodecDictionaryNotes[0] = "ab70-changed"
	*b509HeatpumpCodecDictionaryAgain.B509HeatpumpCodecDictionaryKey.B509HeatpumpCodecDictionaryPrimary = 0
	b509HeatpumpCodecDictionaryAgain, b509HeatpumpCodecDictionaryFound, b509HeatpumpCodecDictionaryError = B509HeatpumpCodecDictionaryLookup(b509HeatpumpCodecDictionaryKey(0xD0, 0))
	if !b509HeatpumpCodecDictionaryFound || b509HeatpumpCodecDictionaryError != nil || b509HeatpumpCodecDictionaryAgain.B509HeatpumpCodecDictionaryNotes[0] != "Literal codec declaration; no executable converter." {
		b509HeatpumpCodecDictionaryT.Fatal("ab70 repeated lookup shared")
	}
	b509HeatpumpCodecDictionaryDuplicate, _ := b509HeatpumpCodecDictionaryJSONPackage.Marshal([]B509HeatpumpCodecDictionaryEntry{b509HeatpumpCodecDictionaryEntry, b509HeatpumpCodecDictionaryEntry})
	for _, b509HeatpumpCodecDictionaryBad := range []string{"{", "null", "[]", "[{}]", string(b509HeatpumpCodecDictionaryDuplicate), string(b509HeatpumpCodecDictionaryOriginal) + " []", b509HeatpumpCodecDictionaryStrings.Replace(string(b509HeatpumpCodecDictionaryOriginal), `"quantity":`, `"ab70_extra":0,"quantity":`, 1), b509HeatpumpCodecDictionaryStrings.Replace(string(b509HeatpumpCodecDictionaryOriginal), `"selector":[208,0]`, `"selector":[208]`, 1), b509HeatpumpCodecDictionaryStrings.Replace(string(b509HeatpumpCodecDictionaryOriginal), `"selector":[208,0]`, `"selector":[208,0,0]`, 1), b509HeatpumpCodecDictionaryStrings.Replace(string(b509HeatpumpCodecDictionaryOriginal), `"selector":[208,0]`, `"selector":[256,0]`, 1)} {
		if b509HeatpumpCodecDictionaryResult, b509HeatpumpCodecDictionaryError := b509HeatpumpCodecDictionaryParse([]byte(b509HeatpumpCodecDictionaryBad)); b509HeatpumpCodecDictionaryError == nil || b509HeatpumpCodecDictionaryResult != nil {
			b509HeatpumpCodecDictionaryT.Fatal("ab70 malformed/duplicate admitted", b509HeatpumpCodecDictionaryBad)
		}
	}
	b509HeatpumpCodecDictionaryValue := b509HeatpumpCodecDictionaryReflect.ValueOf(&b509HeatpumpCodecDictionaryEntry).Elem()
	for b509HeatpumpCodecDictionaryIndex := 0; b509HeatpumpCodecDictionaryIndex < b509HeatpumpCodecDictionaryValue.NumField(); b509HeatpumpCodecDictionaryIndex++ {
		b509HeatpumpCodecDictionaryField := b509HeatpumpCodecDictionaryValue.Type().Field(b509HeatpumpCodecDictionaryIndex)
		if b509HeatpumpCodecDictionaryField.Name == "B509HeatpumpCodecDictionaryQuantity" || b509HeatpumpCodecDictionaryField.Name == "B509HeatpumpCodecDeclaration" || b509HeatpumpCodecDictionaryField.Name == "B509HeatpumpCodecDictionaryGrade" || b509HeatpumpCodecDictionaryField.Name == "B509HeatpumpCodecDictionarySoftwareContract" || b509HeatpumpCodecDictionaryStrings.Contains("B509HeatpumpCodecDictionaryQualification B509HeatpumpCodecDictionaryPresence B509HeatpumpCodecDictionaryFreshness B509HeatpumpCodecDictionaryExecution B509HeatpumpCodecDictionaryPhysicalEffects B509HeatpumpCodecDictionaryPersistence B509HeatpumpCodecDictionaryDelivery", b509HeatpumpCodecDictionaryField.Name) {
			b509HeatpumpCodecDictionarySaved := b509HeatpumpCodecDictionaryValue.Field(b509HeatpumpCodecDictionaryIndex).String()
			b509HeatpumpCodecDictionaryValue.Field(b509HeatpumpCodecDictionaryIndex).SetString("")
			b509HeatpumpCodecDictionaryBad, _ := b509HeatpumpCodecDictionaryJSONPackage.Marshal([]B509HeatpumpCodecDictionaryEntry{b509HeatpumpCodecDictionaryEntry})
			if _, b509HeatpumpCodecDictionaryError := b509HeatpumpCodecDictionaryParse(b509HeatpumpCodecDictionaryBad); b509HeatpumpCodecDictionaryError == nil {
				b509HeatpumpCodecDictionaryT.Errorf("ab70 missing %s admitted", b509HeatpumpCodecDictionaryField.Name)
			}
			b509HeatpumpCodecDictionaryValue.Field(b509HeatpumpCodecDictionaryIndex).SetString(b509HeatpumpCodecDictionarySaved)
		}
	}
	b509HeatpumpCodecDictionaryAbsent := b509HeatpumpCodecDictionaryEntry
	b509HeatpumpCodecDictionaryAbsent.B509HeatpumpCodecDictionaryWidthBits = nil
	b509HeatpumpCodecDictionaryAbsent.B509HeatpumpCodecDictionarySentinel = nil
	b509HeatpumpCodecDictionaryAbsent.B509HeatpumpCodecDictionaryUnit = nil
	if b509HeatpumpCodecDictionaryReflect.DeepEqual(b509HeatpumpCodecDictionaryAbsent, b509HeatpumpCodecDictionaryEntry) {
		b509HeatpumpCodecDictionaryT.Fatal("ab70 absent equals explicit zero")
	}
}

func TestAb70OptionalAndPartialAxes(b509HeatpumpCodecDictionaryT *b509HeatpumpCodecDictionaryTesting.T) {
	b509HeatpumpCodecDictionaryEntries, b509HeatpumpCodecDictionaryError := B509HeatpumpCodecDictionaryEntries()
	if b509HeatpumpCodecDictionaryError != nil {
		b509HeatpumpCodecDictionaryT.Fatal(b509HeatpumpCodecDictionaryError)
	}
	b509HeatpumpCodecDictionaryOptional := []struct {
		b509HeatpumpCodecDictionaryName   string
		b509HeatpumpCodecDictionaryChange func(*B509HeatpumpCodecDictionaryEntry)
	}{
		{"ab70Unit", func(b509HeatpumpCodecDictionaryE *B509HeatpumpCodecDictionaryEntry) {
			b509HeatpumpCodecDictionaryE.B509HeatpumpCodecDictionaryUnit = b509HeatpumpCodecDictionaryPointer("")
		}},
		{"ab70Width", func(b509HeatpumpCodecDictionaryE *B509HeatpumpCodecDictionaryEntry) {
			b509HeatpumpCodecDictionaryE.B509HeatpumpCodecDictionaryWidthBits = b509HeatpumpCodecDictionaryPointer(0)
		}},
		{"ab70Scale", func(b509HeatpumpCodecDictionaryE *B509HeatpumpCodecDictionaryEntry) {
			b509HeatpumpCodecDictionaryE.B509HeatpumpCodecDictionaryDeclaredScale = b509HeatpumpCodecDictionaryPointer(0)
		}},
		{"ab70Sentinel", func(b509HeatpumpCodecDictionaryE *B509HeatpumpCodecDictionaryEntry) {
			b509HeatpumpCodecDictionaryE.B509HeatpumpCodecDictionarySentinel = b509HeatpumpCodecDictionaryPointer(0)
		}},
		{"ab70Signedness", func(b509HeatpumpCodecDictionaryE *B509HeatpumpCodecDictionaryEntry) {
			b509HeatpumpCodecDictionaryE.B509HeatpumpCodecDictionarySignedness = b509HeatpumpCodecDictionaryPointer("")
		}},
		{"ab70ByteOrder", func(b509HeatpumpCodecDictionaryE *B509HeatpumpCodecDictionaryEntry) {
			b509HeatpumpCodecDictionaryE.B509HeatpumpCodecDictionaryByteOrder = b509HeatpumpCodecDictionaryPointer("")
		}},
		{"ab70Converter", func(b509HeatpumpCodecDictionaryE *B509HeatpumpCodecDictionaryEntry) {
			b509HeatpumpCodecDictionaryE.B509HeatpumpCodecDictionaryConverter = b509HeatpumpCodecDictionaryPointer("ab70-opaque")
		}},
		{"ab70Mapping", func(b509HeatpumpCodecDictionaryE *B509HeatpumpCodecDictionaryEntry) {
			b509HeatpumpCodecDictionaryE.B509HeatpumpCodecDictionaryMapping = map[string]string{"ab70-zero": "ab70-label"}
		}},
	}
	for _, b509HeatpumpCodecDictionaryCase := range b509HeatpumpCodecDictionaryOptional {
		b509HeatpumpCodecDictionaryT.Run(b509HeatpumpCodecDictionaryCase.b509HeatpumpCodecDictionaryName, func(b509HeatpumpCodecDictionaryT *b509HeatpumpCodecDictionaryTesting.T) {
			b509HeatpumpCodecDictionaryEntry := b509HeatpumpCodecDictionaryEntries[0]
			b509HeatpumpCodecDictionaryCase.b509HeatpumpCodecDictionaryChange(&b509HeatpumpCodecDictionaryEntry)
			b509HeatpumpCodecDictionaryData, _ := b509HeatpumpCodecDictionaryJSONPackage.Marshal([]B509HeatpumpCodecDictionaryEntry{b509HeatpumpCodecDictionaryEntry})
			b509HeatpumpCodecDictionaryParsed, b509HeatpumpCodecDictionaryError := b509HeatpumpCodecDictionaryParse(b509HeatpumpCodecDictionaryData)
			if b509HeatpumpCodecDictionaryError != nil || !b509HeatpumpCodecDictionaryReflect.DeepEqual(b509HeatpumpCodecDictionaryParsed, []B509HeatpumpCodecDictionaryEntry{b509HeatpumpCodecDictionaryEntry}) || b509HeatpumpCodecDictionaryReflect.DeepEqual(b509HeatpumpCodecDictionaryParsed, b509HeatpumpCodecDictionaryEntries[:1]) {
				b509HeatpumpCodecDictionaryT.Fatal("ab70 optional field lost", b509HeatpumpCodecDictionaryError)
			}
			b509HeatpumpCodecDictionaryCopy, _, _ := b509HeatpumpCodecDictionaryFind(b509HeatpumpCodecDictionaryParsed, b509HeatpumpCodecDictionaryEntry.B509HeatpumpCodecDictionaryKey)
			b509HeatpumpCodecDictionaryValue := b509HeatpumpCodecDictionaryReflect.ValueOf(&b509HeatpumpCodecDictionaryCopy).Elem()
			for b509HeatpumpCodecDictionaryIndex := 0; b509HeatpumpCodecDictionaryIndex < b509HeatpumpCodecDictionaryValue.NumField(); b509HeatpumpCodecDictionaryIndex++ {
				b509HeatpumpCodecDictionaryField := b509HeatpumpCodecDictionaryValue.Field(b509HeatpumpCodecDictionaryIndex)
				if b509HeatpumpCodecDictionaryField.Kind() == b509HeatpumpCodecDictionaryReflect.Pointer && !b509HeatpumpCodecDictionaryField.IsNil() {
					if b509HeatpumpCodecDictionaryField.Elem().Kind() == b509HeatpumpCodecDictionaryReflect.String {
						b509HeatpumpCodecDictionaryField.Elem().SetString("ab70-mutated")
					} else {
						b509HeatpumpCodecDictionaryField.Elem().SetInt(7)
					}
				}
			}
			if b509HeatpumpCodecDictionaryCopy.B509HeatpumpCodecDictionaryMapping != nil {
				b509HeatpumpCodecDictionaryCopy.B509HeatpumpCodecDictionaryMapping["ab70-zero"] = "ab70-mutated"
			}
			if !b509HeatpumpCodecDictionaryReflect.DeepEqual(b509HeatpumpCodecDictionaryParsed, []B509HeatpumpCodecDictionaryEntry{b509HeatpumpCodecDictionaryEntry}) {
				b509HeatpumpCodecDictionaryT.Fatal("ab70 optional pointer shared")
			}
		})
	}
	b509HeatpumpCodecDictionaryKeyChanges := []func(*B509HeatpumpCodecDictionaryKey){
		func(b509HeatpumpCodecDictionaryK *B509HeatpumpCodecDictionaryKey) {
			b509HeatpumpCodecDictionaryK.B509HeatpumpCodecDictionaryFamily = ""
		},
		func(b509HeatpumpCodecDictionaryK *B509HeatpumpCodecDictionaryKey) {
			b509HeatpumpCodecDictionaryK.B509HeatpumpCodecDictionaryVariant = ""
		},
		func(b509HeatpumpCodecDictionaryK *B509HeatpumpCodecDictionaryKey) {
			b509HeatpumpCodecDictionaryK.B509HeatpumpCodecDictionaryDirection = ""
		},
		func(b509HeatpumpCodecDictionaryK *B509HeatpumpCodecDictionaryKey) {
			b509HeatpumpCodecDictionaryK.B509HeatpumpCodecDictionaryPrimary = nil
		},
		func(b509HeatpumpCodecDictionaryK *B509HeatpumpCodecDictionaryKey) {
			b509HeatpumpCodecDictionaryK.B509HeatpumpCodecDictionarySecondary = nil
		},
		func(b509HeatpumpCodecDictionaryK *B509HeatpumpCodecDictionaryKey) {
			b509HeatpumpCodecDictionaryK.B509HeatpumpCodecDictionaryOperation = nil
		},
	}
	for b509HeatpumpCodecDictionaryIndex, b509HeatpumpCodecDictionaryChange := range b509HeatpumpCodecDictionaryKeyChanges {
		b509HeatpumpCodecDictionaryEntry := b509HeatpumpCodecDictionaryEntries[0]
		b509HeatpumpCodecDictionaryChange(&b509HeatpumpCodecDictionaryEntry.B509HeatpumpCodecDictionaryKey)
		b509HeatpumpCodecDictionaryData, _ := b509HeatpumpCodecDictionaryJSONPackage.Marshal([]B509HeatpumpCodecDictionaryEntry{b509HeatpumpCodecDictionaryEntry})
		if b509HeatpumpCodecDictionaryParsed, b509HeatpumpCodecDictionaryError := b509HeatpumpCodecDictionaryParse(b509HeatpumpCodecDictionaryData); b509HeatpumpCodecDictionaryError == nil || b509HeatpumpCodecDictionaryParsed != nil {
			b509HeatpumpCodecDictionaryT.Fatalf("ab70 partial axis %d admitted", b509HeatpumpCodecDictionaryIndex)
		}
	}
	for _, b509HeatpumpCodecDictionaryChange := range []func(*B509HeatpumpCodecDictionaryKey){
		func(b509HeatpumpCodecDictionaryK *B509HeatpumpCodecDictionaryKey) {
			b509HeatpumpCodecDictionaryK.B509HeatpumpCodecDictionaryPrimary = b509HeatpumpCodecDictionaryPointer(byte(0))
		},
		func(b509HeatpumpCodecDictionaryK *B509HeatpumpCodecDictionaryKey) {
			b509HeatpumpCodecDictionaryK.B509HeatpumpCodecDictionarySecondary = b509HeatpumpCodecDictionaryPointer(byte(0))
		},
		func(b509HeatpumpCodecDictionaryK *B509HeatpumpCodecDictionaryKey) {
			b509HeatpumpCodecDictionaryK.B509HeatpumpCodecDictionaryOperation = b509HeatpumpCodecDictionaryPointer(byte(0))
		},
		func(b509HeatpumpCodecDictionaryK *B509HeatpumpCodecDictionaryKey) {
			b509HeatpumpCodecDictionaryK.B509HeatpumpCodecDictionarySelector = []uint16{0, 0}
		},
	} {
		b509HeatpumpCodecDictionaryEntry := b509HeatpumpCodecDictionaryEntries[0]
		b509HeatpumpCodecDictionaryChange(&b509HeatpumpCodecDictionaryEntry.B509HeatpumpCodecDictionaryKey)
		b509HeatpumpCodecDictionaryData, _ := b509HeatpumpCodecDictionaryJSONPackage.Marshal([]B509HeatpumpCodecDictionaryEntry{b509HeatpumpCodecDictionaryEntry})
		b509HeatpumpCodecDictionaryParsed, b509HeatpumpCodecDictionaryError := b509HeatpumpCodecDictionaryParse(b509HeatpumpCodecDictionaryData)
		if b509HeatpumpCodecDictionaryError != nil || !b509HeatpumpCodecDictionaryReflect.DeepEqual(b509HeatpumpCodecDictionaryParsed, []B509HeatpumpCodecDictionaryEntry{b509HeatpumpCodecDictionaryEntry}) {
			b509HeatpumpCodecDictionaryT.Fatal("ab70 numeric zero rejected", b509HeatpumpCodecDictionaryError)
		}
	}
}

func TestAb70FailedReadPreservesOwnedData(b509HeatpumpCodecDictionaryT *b509HeatpumpCodecDictionaryTesting.T) {
	b509HeatpumpCodecDictionarySaved := b509HeatpumpCodecDictionaryJSON
	defer func() { b509HeatpumpCodecDictionaryJSON = b509HeatpumpCodecDictionarySaved }()
	b509HeatpumpCodecDictionaryGood, b509HeatpumpCodecDictionaryFound, b509HeatpumpCodecDictionaryError := B509HeatpumpCodecDictionaryLookup(b509HeatpumpCodecDictionaryKey(0xBC, 0))
	if !b509HeatpumpCodecDictionaryFound || b509HeatpumpCodecDictionaryError != nil {
		b509HeatpumpCodecDictionaryT.Fatal("ab70 valid control failed")
	}
	b509HeatpumpCodecDictionarySnapshot, b509HeatpumpCodecDictionaryError := b509HeatpumpCodecDictionaryJSONPackage.Marshal(b509HeatpumpCodecDictionaryGood)
	if b509HeatpumpCodecDictionaryError != nil {
		b509HeatpumpCodecDictionaryT.Fatal(b509HeatpumpCodecDictionaryError)
	}
	b509HeatpumpCodecDictionaryJSON = []byte(`[{"key":{}}]`)
	if b509HeatpumpCodecDictionaryRows, b509HeatpumpCodecDictionaryError := B509HeatpumpCodecDictionaryEntries(); b509HeatpumpCodecDictionaryError == nil || b509HeatpumpCodecDictionaryRows != nil {
		b509HeatpumpCodecDictionaryT.Fatal("ab70 partial data admitted")
	}
	if b509HeatpumpCodecDictionaryRow, b509HeatpumpCodecDictionaryFound, b509HeatpumpCodecDictionaryError := B509HeatpumpCodecDictionaryLookup(b509HeatpumpCodecDictionaryKey(0xBC, 0)); b509HeatpumpCodecDictionaryError == nil || b509HeatpumpCodecDictionaryFound || !b509HeatpumpCodecDictionaryReflect.DeepEqual(b509HeatpumpCodecDictionaryRow, B509HeatpumpCodecDictionaryEntry{}) {
		b509HeatpumpCodecDictionaryT.Fatal("ab70 lookup error lost")
	}
	b509HeatpumpCodecDictionaryAfter, _ := b509HeatpumpCodecDictionaryJSONPackage.Marshal(b509HeatpumpCodecDictionaryGood)
	if !b509HeatpumpCodecDictionaryBytes.Equal(b509HeatpumpCodecDictionarySnapshot, b509HeatpumpCodecDictionaryAfter) {
		b509HeatpumpCodecDictionaryT.Fatal("ab70 last owned evidence changed after failure")
	}
	b509HeatpumpCodecDictionaryJSON = b509HeatpumpCodecDictionarySaved
	b509HeatpumpCodecDictionaryAgain, b509HeatpumpCodecDictionaryFound, b509HeatpumpCodecDictionaryError := B509HeatpumpCodecDictionaryLookup(b509HeatpumpCodecDictionaryKey(0xBC, 0))
	if !b509HeatpumpCodecDictionaryFound || b509HeatpumpCodecDictionaryError != nil || !b509HeatpumpCodecDictionaryReflect.DeepEqual(b509HeatpumpCodecDictionaryGood, b509HeatpumpCodecDictionaryAgain) {
		b509HeatpumpCodecDictionaryT.Fatal("ab70 recovery control failed")
	}
}
