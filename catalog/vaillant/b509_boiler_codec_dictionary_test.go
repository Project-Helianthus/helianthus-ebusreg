package vaillant

import (
	b509BoilerCodecDictionaryBytes "bytes"
	b509BoilerCodecDictionaryJSONPackage "encoding/json"
	b509BoilerCodecDictionaryReflect "reflect"
	b509BoilerCodecDictionaryStrings "strings"
	b509BoilerCodecDictionaryTesting "testing"
)

func b509BoilerCodecDictionaryPointer[B509BoilerCodecDictionaryT any](b509BoilerCodecDictionaryValue B509BoilerCodecDictionaryT) *B509BoilerCodecDictionaryT {
	return &b509BoilerCodecDictionaryValue
}

func b509BoilerCodecDictionaryFixture() B509BoilerCodecDictionaryEntry {
	return B509BoilerCodecDictionaryEntry{
		B509BoilerCodecDictionaryKey:      B509BoilerCodecDictionaryKey{B509BoilerCodecDictionaryFamily: "boiler", B509BoilerCodecDictionaryVariant: "direct-symbolic-candidate", B509BoilerCodecDictionaryPrimary: b509BoilerCodecDictionaryPointer(byte(0xB5)), B509BoilerCodecDictionarySecondary: b509BoilerCodecDictionaryPointer(byte(9)), B509BoilerCodecDictionaryOperation: b509BoilerCodecDictionaryPointer(byte(13)), B509BoilerCodecDictionaryDirection: "read", B509BoilerCodecDictionarySelector: []uint16{0x18, 0}},
		B509BoilerCodecDictionaryQuantity: "flow_temperature", B509BoilerCodecDeclaration: "DATA2c", B509BoilerCodecDictionaryGrade: "Hypothesis", B509BoilerCodecDictionaryConflicts: []string{},
		B509BoilerCodecDictionaryEvidence: map[string]string{"metadata": "Proven", "qualification": "Unknown", "physical_support": "Unknown", "presence": "Unknown", "freshness": "Unknown", "effects": "Unknown", "persistence": "Unknown", "delivery": "Unknown", "execution": "Unknown"},
	}
}

func b509BoilerCodecDictionaryEncode(b509BoilerCodecDictionaryTest *b509BoilerCodecDictionaryTesting.T, b509BoilerCodecDictionaryValue any) []byte {
	b509BoilerCodecDictionaryTest.Helper()
	b509BoilerCodecDictionaryData, b509BoilerCodecDictionaryError := b509BoilerCodecDictionaryJSONPackage.Marshal(b509BoilerCodecDictionaryValue)
	if b509BoilerCodecDictionaryError != nil {
		b509BoilerCodecDictionaryTest.Fatal(b509BoilerCodecDictionaryError)
	}
	return b509BoilerCodecDictionaryData
}

func TestAb68ExactDictionary(b509BoilerCodecDictionaryTest *b509BoilerCodecDictionaryTesting.T) {
	b509BoilerCodecDictionaryExpected := map[[2]uint16][3]string{
		{0x18, 0}: {"flow_temperature", "DATA2c", ""},
		{0x44, 0}: {"central_heating_pump_active", "on/off", ""},
		{0x02, 0}: {"water_pressure", "DATA2b", ""},
		{0x05, 0}: {"flame_active", "UCH (0x0F=on, 0xF0=off)", ""},
		{0xBB, 0}: {"gas_valve_active", "UCH (0x0F=on, 0xF0=off)", ""},
		{0x83, 0}: {"fan_speed", "UIN", ""},
		{0x2E, 0}: {"modulation_setpoint", "SIN / 10", "setpoint versus achieved modulation unresolved"},
		{0xAB, 0}: {"boiler_state_code", "UCH", ""},
		{0x54, 0}: {"diverter_valve_position", "UCH", "position enum versus percent unresolved"},
		{0x58, 0}: {"dhw_demand_active", "on/off", ""},
		{0x55, 0}: {"dhw_water_flow_rate", "UIN / 100", ""},
		{0x0E, 4}: {"maximum_heating_flow_temperature_setpoint", "DATA2c", ""},
		{0xA5, 0}: {"maximum_heating_flow_temperature_setpoint", "DATA2c", ""},
		{0x0F, 4}: {"maximum_dhw_flow_temperature_setpoint", "DATA2c", ""},
		{0x07, 4}: {"heating_partial_load_power", "UCH", "unit Unknown"},
		{0x08, 4}: {"dhw_partial_load_power", "UCH", "unit Unknown"},
		{0x9E, 0}: {"storage_loading_pump_level", "percent0", ""},
		{0x39, 0}: {"flow_temperature_setpoint", "DATA2c", ""},
		{0xEA, 3}: {"dhw_temperature_setpoint", "DATA2c", ""},
		{0x7B, 0}: {"circulation_pump_active", "on/off", ""},
		{0x3F, 0}: {"external_pump_active", "on/off", ""},
		{0xF2, 3}: {"heating_switch_active", "on/off", ""},
		{0x24, 0}: {"fan_speed_setpoint", "UIN", ""},
		{0xA4, 0}: {"flame_ionization_current", "SIN / 10", "current versus voltage unresolved"},
		{0xFB, 0}: {"primary_circuit_flow_rate", "UIN / 100", "litres/minute versus litres/hour unresolved"},
	}
	b509BoilerCodecDictionaryEntries, b509BoilerCodecDictionaryError := B509BoilerCodecDictionaryEntries()
	if b509BoilerCodecDictionaryError != nil {
		b509BoilerCodecDictionaryTest.Fatal(b509BoilerCodecDictionaryError)
	}
	if len(b509BoilerCodecDictionaryEntries) != 25 {
		b509BoilerCodecDictionaryTest.Fatalf("ab68 exact set: got %d rows", len(b509BoilerCodecDictionaryEntries))
	}
	b509BoilerCodecDictionarySeen := make(map[[2]uint16]bool)
	for _, b509BoilerCodecDictionaryEntry := range b509BoilerCodecDictionaryEntries {
		b509BoilerCodecDictionarySelector := [2]uint16{b509BoilerCodecDictionaryEntry.B509BoilerCodecDictionaryKey.B509BoilerCodecDictionarySelector[0], b509BoilerCodecDictionaryEntry.B509BoilerCodecDictionaryKey.B509BoilerCodecDictionarySelector[1]}
		b509BoilerCodecDictionaryWant, b509BoilerCodecDictionaryFound := b509BoilerCodecDictionaryExpected[b509BoilerCodecDictionarySelector]
		if !b509BoilerCodecDictionaryFound || b509BoilerCodecDictionarySeen[b509BoilerCodecDictionarySelector] {
			b509BoilerCodecDictionaryTest.Fatalf("ab68 unexpected/duplicate selector: %04X", b509BoilerCodecDictionarySelector)
		}
		b509BoilerCodecDictionarySeen[b509BoilerCodecDictionarySelector] = true
		b509BoilerCodecDictionaryControl := b509BoilerCodecDictionaryFixture()
		b509BoilerCodecDictionaryControl.B509BoilerCodecDictionaryKey.B509BoilerCodecDictionarySelector = []uint16{b509BoilerCodecDictionarySelector[0], b509BoilerCodecDictionarySelector[1]}
		b509BoilerCodecDictionaryControl.B509BoilerCodecDictionaryQuantity = b509BoilerCodecDictionaryWant[0]
		b509BoilerCodecDictionaryControl.B509BoilerCodecDeclaration = b509BoilerCodecDictionaryWant[1]
		if b509BoilerCodecDictionaryWant[2] != "" {
			b509BoilerCodecDictionaryControl.B509BoilerCodecDictionaryConflicts = []string{b509BoilerCodecDictionaryWant[2]}
		}
		if !b509BoilerCodecDictionaryReflect.DeepEqual(b509BoilerCodecDictionaryEntry, b509BoilerCodecDictionaryControl) {
			b509BoilerCodecDictionaryTest.Errorf("ab68 metadata for %04X: got %+v, want %+v", b509BoilerCodecDictionarySelector, b509BoilerCodecDictionaryEntry, b509BoilerCodecDictionaryControl)
		}
		b509BoilerCodecDictionaryResult, b509BoilerCodecDictionaryOK, b509BoilerCodecDictionaryLookupError := B509BoilerCodecDictionaryLookup(b509BoilerCodecDictionaryControl.B509BoilerCodecDictionaryKey)
		if b509BoilerCodecDictionaryLookupError != nil || !b509BoilerCodecDictionaryOK || !b509BoilerCodecDictionaryReflect.DeepEqual(b509BoilerCodecDictionaryResult, b509BoilerCodecDictionaryControl) {
			b509BoilerCodecDictionaryTest.Errorf("ab68 lookup %04X: %+v %v %v", b509BoilerCodecDictionarySelector, b509BoilerCodecDictionaryResult, b509BoilerCodecDictionaryOK, b509BoilerCodecDictionaryLookupError)
		}
	}
	if len(b509BoilerCodecDictionarySeen) != len(b509BoilerCodecDictionaryExpected) {
		b509BoilerCodecDictionaryTest.Fatal("ab68 missing finite tuple")
	}
}

func TestAb68KeyAxes(b509BoilerCodecDictionaryTest *b509BoilerCodecDictionaryTesting.T) {
	b509BoilerCodecDictionaryCases := map[string]func(*B509BoilerCodecDictionaryKey){
		"ab68Family": func(b509BoilerCodecDictionaryKey *B509BoilerCodecDictionaryKey) {
			b509BoilerCodecDictionaryKey.B509BoilerCodecDictionaryFamily = "interface-unit"
		},
		"ab68UnknownFamily": func(b509BoilerCodecDictionaryKey *B509BoilerCodecDictionaryKey) {
			b509BoilerCodecDictionaryKey.B509BoilerCodecDictionaryFamily = "ab68-unknown"
		},
		"ab68Variant": func(b509BoilerCodecDictionaryKey *B509BoilerCodecDictionaryKey) {
			b509BoilerCodecDictionaryKey.B509BoilerCodecDictionaryVariant = "ab68-other"
		},
		"ab68Primary": func(b509BoilerCodecDictionaryKey *B509BoilerCodecDictionaryKey) {
			b509BoilerCodecDictionaryKey.B509BoilerCodecDictionaryPrimary = b509BoilerCodecDictionaryPointer(byte(0))
		},
		"ab68Secondary": func(b509BoilerCodecDictionaryKey *B509BoilerCodecDictionaryKey) {
			b509BoilerCodecDictionaryKey.B509BoilerCodecDictionarySecondary = b509BoilerCodecDictionaryPointer(byte(0))
		},
		"ab68OperationWrite": func(b509BoilerCodecDictionaryKey *B509BoilerCodecDictionaryKey) {
			b509BoilerCodecDictionaryKey.B509BoilerCodecDictionaryOperation = b509BoilerCodecDictionaryPointer(byte(0x0E))
		},
		"ab68OperationAlternate": func(b509BoilerCodecDictionaryKey *B509BoilerCodecDictionaryKey) {
			b509BoilerCodecDictionaryKey.B509BoilerCodecDictionaryOperation = b509BoilerCodecDictionaryPointer(byte(0x29))
		},
		"ab68OperationZero": func(b509BoilerCodecDictionaryKey *B509BoilerCodecDictionaryKey) {
			b509BoilerCodecDictionaryKey.B509BoilerCodecDictionaryOperation = b509BoilerCodecDictionaryPointer(byte(0))
		},
		"ab68DirectionRequest": func(b509BoilerCodecDictionaryKey *B509BoilerCodecDictionaryKey) {
			b509BoilerCodecDictionaryKey.B509BoilerCodecDictionaryDirection = "request"
		},
		"ab68DirectionPositive": func(b509BoilerCodecDictionaryKey *B509BoilerCodecDictionaryKey) {
			b509BoilerCodecDictionaryKey.B509BoilerCodecDictionaryDirection = "positive"
		},
		"ab68SelectorFirst": func(b509BoilerCodecDictionaryKey *B509BoilerCodecDictionaryKey) {
			b509BoilerCodecDictionaryKey.B509BoilerCodecDictionarySelector[0]++
		},
		"ab68SelectorSecond": func(b509BoilerCodecDictionaryKey *B509BoilerCodecDictionaryKey) {
			b509BoilerCodecDictionaryKey.B509BoilerCodecDictionarySelector[1]++
		},
		"ab68SelectorZero": func(b509BoilerCodecDictionaryKey *B509BoilerCodecDictionaryKey) {
			b509BoilerCodecDictionaryKey.B509BoilerCodecDictionarySelector = []uint16{0, 0}
		},
		"ab68SelectorReversed": func(b509BoilerCodecDictionaryKey *B509BoilerCodecDictionaryKey) {
			b509BoilerCodecDictionaryKey.B509BoilerCodecDictionarySelector = []uint16{0, 0x18}
		},
	}
	for b509BoilerCodecDictionaryName, b509BoilerCodecDictionaryMutate := range b509BoilerCodecDictionaryCases {
		b509BoilerCodecDictionaryTest.Run(b509BoilerCodecDictionaryName, func(b509BoilerCodecDictionaryTest *b509BoilerCodecDictionaryTesting.T) {
			b509BoilerCodecDictionaryFirst, b509BoilerCodecDictionarySecond := b509BoilerCodecDictionaryFixture(), b509BoilerCodecDictionaryFixture()
			b509BoilerCodecDictionaryMutate(&b509BoilerCodecDictionarySecond.B509BoilerCodecDictionaryKey)
			b509BoilerCodecDictionaryResult, b509BoilerCodecDictionaryOK, b509BoilerCodecDictionaryError := B509BoilerCodecDictionaryLookup(b509BoilerCodecDictionarySecond.B509BoilerCodecDictionaryKey)
			if b509BoilerCodecDictionaryError != nil || b509BoilerCodecDictionaryOK || !b509BoilerCodecDictionaryReflect.DeepEqual(b509BoilerCodecDictionaryResult, B509BoilerCodecDictionaryEntry{}) {
				b509BoilerCodecDictionaryTest.Fatalf("ab68 unexpected admission: %+v %v %v", b509BoilerCodecDictionaryResult, b509BoilerCodecDictionaryOK, b509BoilerCodecDictionaryError)
			}
			b509BoilerCodecDictionarySecond.B509BoilerCodecDictionaryQuantity = "ab68-synthetic-control"
			b509BoilerCodecDictionaryData := b509BoilerCodecDictionaryEncode(b509BoilerCodecDictionaryTest, []B509BoilerCodecDictionaryEntry{b509BoilerCodecDictionaryFirst, b509BoilerCodecDictionarySecond})
			for _, b509BoilerCodecDictionaryControl := range []B509BoilerCodecDictionaryEntry{b509BoilerCodecDictionaryFirst, b509BoilerCodecDictionarySecond} {
				b509BoilerCodecDictionaryResult, b509BoilerCodecDictionaryOK, b509BoilerCodecDictionaryError = b509BoilerCodecDictionaryLookup(b509BoilerCodecDictionaryData, b509BoilerCodecDictionaryControl.B509BoilerCodecDictionaryKey)
				if b509BoilerCodecDictionaryError != nil || !b509BoilerCodecDictionaryOK || !b509BoilerCodecDictionaryReflect.DeepEqual(b509BoilerCodecDictionaryResult, b509BoilerCodecDictionaryControl) {
					b509BoilerCodecDictionaryTest.Fatalf("ab68 independent key control: %+v %v %v", b509BoilerCodecDictionaryResult, b509BoilerCodecDictionaryOK, b509BoilerCodecDictionaryError)
				}
			}
		})
	}
	b509BoilerCodecDictionaryAbsent := map[string]func(*B509BoilerCodecDictionaryKey){
		"ab68Family": func(b509BoilerCodecDictionaryKey *B509BoilerCodecDictionaryKey) {
			b509BoilerCodecDictionaryKey.B509BoilerCodecDictionaryFamily = ""
		},
		"ab68Variant": func(b509BoilerCodecDictionaryKey *B509BoilerCodecDictionaryKey) {
			b509BoilerCodecDictionaryKey.B509BoilerCodecDictionaryVariant = ""
		},
		"ab68Primary": func(b509BoilerCodecDictionaryKey *B509BoilerCodecDictionaryKey) {
			b509BoilerCodecDictionaryKey.B509BoilerCodecDictionaryPrimary = nil
		},
		"ab68Secondary": func(b509BoilerCodecDictionaryKey *B509BoilerCodecDictionaryKey) {
			b509BoilerCodecDictionaryKey.B509BoilerCodecDictionarySecondary = nil
		},
		"ab68Operation": func(b509BoilerCodecDictionaryKey *B509BoilerCodecDictionaryKey) {
			b509BoilerCodecDictionaryKey.B509BoilerCodecDictionaryOperation = nil
		},
		"ab68Direction": func(b509BoilerCodecDictionaryKey *B509BoilerCodecDictionaryKey) {
			b509BoilerCodecDictionaryKey.B509BoilerCodecDictionaryDirection = ""
		},
		"ab68Selector": func(b509BoilerCodecDictionaryKey *B509BoilerCodecDictionaryKey) {
			b509BoilerCodecDictionaryKey.B509BoilerCodecDictionarySelector = nil
		},
	}
	for b509BoilerCodecDictionaryName, b509BoilerCodecDictionaryMutate := range b509BoilerCodecDictionaryAbsent {
		b509BoilerCodecDictionaryTest.Run(b509BoilerCodecDictionaryName+"Absent", func(b509BoilerCodecDictionaryTest *b509BoilerCodecDictionaryTesting.T) {
			b509BoilerCodecDictionaryEntry := b509BoilerCodecDictionaryFixture()
			b509BoilerCodecDictionaryMutate(&b509BoilerCodecDictionaryEntry.B509BoilerCodecDictionaryKey)
			if _, b509BoilerCodecDictionaryOK, b509BoilerCodecDictionaryError := B509BoilerCodecDictionaryLookup(b509BoilerCodecDictionaryEntry.B509BoilerCodecDictionaryKey); b509BoilerCodecDictionaryOK || b509BoilerCodecDictionaryError != nil {
				b509BoilerCodecDictionaryTest.Fatal("ab68 absent key matched")
			}
			if b509BoilerCodecDictionaryEntries, b509BoilerCodecDictionaryError := b509BoilerCodecDictionaryParse(b509BoilerCodecDictionaryEncode(b509BoilerCodecDictionaryTest, []B509BoilerCodecDictionaryEntry{b509BoilerCodecDictionaryEntry})); b509BoilerCodecDictionaryError == nil || b509BoilerCodecDictionaryEntries != nil {
				b509BoilerCodecDictionaryTest.Fatal("ab68 incomplete source admitted")
			}
		})
	}
}

func TestAb68Validation(b509BoilerCodecDictionaryTest *b509BoilerCodecDictionaryTesting.T) {
	b509BoilerCodecDictionaryValid := b509BoilerCodecDictionaryEncode(b509BoilerCodecDictionaryTest, []B509BoilerCodecDictionaryEntry{b509BoilerCodecDictionaryFixture()})
	b509BoilerCodecDictionaryCases := map[string][]byte{
		"ab68Malformed": []byte("["), "ab68Null": []byte("null"), "ab68Empty": []byte("[]"), "ab68Partial": []byte("[{}]"),
		"ab68Trailing":          append(append([]byte{}, b509BoilerCodecDictionaryValid...), []byte(" {}")...),
		"ab68TrailingMalformed": append(append([]byte{}, b509BoilerCodecDictionaryValid...), byte('x')),
		"ab68UnknownField":      []byte(b509BoilerCodecDictionaryStrings.Replace(string(b509BoilerCodecDictionaryValid), `"quantity":`, `"ab68_extra":0,"quantity":`, 1)),
		"ab68Duplicate":         b509BoilerCodecDictionaryEncode(b509BoilerCodecDictionaryTest, []B509BoilerCodecDictionaryEntry{b509BoilerCodecDictionaryFixture(), b509BoilerCodecDictionaryFixture()}),
	}
	for b509BoilerCodecDictionaryName, b509BoilerCodecDictionaryData := range b509BoilerCodecDictionaryCases {
		b509BoilerCodecDictionaryTest.Run(b509BoilerCodecDictionaryName, func(b509BoilerCodecDictionaryTest *b509BoilerCodecDictionaryTesting.T) {
			b509BoilerCodecDictionaryBefore := append([]byte{}, b509BoilerCodecDictionaryData...)
			if b509BoilerCodecDictionaryEntries, b509BoilerCodecDictionaryError := b509BoilerCodecDictionaryParse(b509BoilerCodecDictionaryData); b509BoilerCodecDictionaryError == nil || b509BoilerCodecDictionaryEntries != nil {
				b509BoilerCodecDictionaryTest.Fatal("ab68 invalid data admitted")
			}
			if b509BoilerCodecDictionaryResult, b509BoilerCodecDictionaryOK, b509BoilerCodecDictionaryError := b509BoilerCodecDictionaryLookup(b509BoilerCodecDictionaryData, b509BoilerCodecDictionaryFixture().B509BoilerCodecDictionaryKey); b509BoilerCodecDictionaryError == nil || b509BoilerCodecDictionaryOK || !b509BoilerCodecDictionaryReflect.DeepEqual(b509BoilerCodecDictionaryResult, B509BoilerCodecDictionaryEntry{}) {
				b509BoilerCodecDictionaryTest.Fatal("ab68 parse error not propagated")
			}
			if !b509BoilerCodecDictionaryBytes.Equal(b509BoilerCodecDictionaryData, b509BoilerCodecDictionaryBefore) {
				b509BoilerCodecDictionaryTest.Fatal("ab68 raw source changed on failure")
			}
		})
	}
	b509BoilerCodecDictionaryInvalid := map[string]func(*B509BoilerCodecDictionaryEntry){
		"ab68ShortSelector": func(b509BoilerCodecDictionaryEntry *B509BoilerCodecDictionaryEntry) {
			b509BoilerCodecDictionaryEntry.B509BoilerCodecDictionaryKey.B509BoilerCodecDictionarySelector = []uint16{24}
		},
		"ab68LongSelector": func(b509BoilerCodecDictionaryEntry *B509BoilerCodecDictionaryEntry) {
			b509BoilerCodecDictionaryEntry.B509BoilerCodecDictionaryKey.B509BoilerCodecDictionarySelector = []uint16{24, 0, 0}
		},
		"ab68WideSelector": func(b509BoilerCodecDictionaryEntry *B509BoilerCodecDictionaryEntry) {
			b509BoilerCodecDictionaryEntry.B509BoilerCodecDictionaryKey.B509BoilerCodecDictionarySelector[0] = 256
		},
		"ab68Quantity": func(b509BoilerCodecDictionaryEntry *B509BoilerCodecDictionaryEntry) {
			b509BoilerCodecDictionaryEntry.B509BoilerCodecDictionaryQuantity = ""
		},
		"ab68Codec": func(b509BoilerCodecDictionaryEntry *B509BoilerCodecDictionaryEntry) {
			b509BoilerCodecDictionaryEntry.B509BoilerCodecDeclaration = ""
		},
		"ab68Grade": func(b509BoilerCodecDictionaryEntry *B509BoilerCodecDictionaryEntry) {
			b509BoilerCodecDictionaryEntry.B509BoilerCodecDictionaryGrade = "Proven"
		},
		"ab68Evidence": func(b509BoilerCodecDictionaryEntry *B509BoilerCodecDictionaryEntry) {
			b509BoilerCodecDictionaryEntry.B509BoilerCodecDictionaryEvidence = nil
		},
	}
	for _, b509BoilerCodecDictionaryAxis := range []string{"metadata", "qualification", "physical_support", "presence", "freshness", "effects", "persistence", "delivery", "execution"} {
		b509BoilerCodecDictionaryAxis := b509BoilerCodecDictionaryAxis
		b509BoilerCodecDictionaryInvalid["ab68Evidence-"+b509BoilerCodecDictionaryAxis] = func(b509BoilerCodecDictionaryEntry *B509BoilerCodecDictionaryEntry) {
			b509BoilerCodecDictionaryEntry.B509BoilerCodecDictionaryEvidence[b509BoilerCodecDictionaryAxis] = "Hypothesis"
		}
	}
	for b509BoilerCodecDictionaryName, b509BoilerCodecDictionaryMutate := range b509BoilerCodecDictionaryInvalid {
		b509BoilerCodecDictionaryTest.Run(b509BoilerCodecDictionaryName, func(b509BoilerCodecDictionaryTest *b509BoilerCodecDictionaryTesting.T) {
			b509BoilerCodecDictionaryEntry := b509BoilerCodecDictionaryFixture()
			b509BoilerCodecDictionaryMutate(&b509BoilerCodecDictionaryEntry)
			if b509BoilerCodecDictionaryEntries, b509BoilerCodecDictionaryError := b509BoilerCodecDictionaryParse(b509BoilerCodecDictionaryEncode(b509BoilerCodecDictionaryTest, []B509BoilerCodecDictionaryEntry{b509BoilerCodecDictionaryEntry})); b509BoilerCodecDictionaryError == nil || b509BoilerCodecDictionaryEntries != nil {
				b509BoilerCodecDictionaryTest.Fatal("ab68 invalid entry admitted")
			}
		})
	}
	if b509BoilerCodecDictionaryEntries, b509BoilerCodecDictionaryError := b509BoilerCodecDictionaryParse(b509BoilerCodecDictionaryValid); b509BoilerCodecDictionaryError != nil || len(b509BoilerCodecDictionaryEntries) != 1 {
		b509BoilerCodecDictionaryTest.Fatalf("ab68 valid control: %v", b509BoilerCodecDictionaryError)
	}
}

func TestAb68OptionalFieldsAndUnknownCodec(b509BoilerCodecDictionaryTest *b509BoilerCodecDictionaryTesting.T) {
	b509BoilerCodecDictionaryCases := map[string]func(*B509BoilerCodecDictionaryEntry){
		"ab68Width": func(b509BoilerCodecDictionaryEntry *B509BoilerCodecDictionaryEntry) {
			b509BoilerCodecDictionaryEntry.B509BoilerCodecDictionaryWidthBits = b509BoilerCodecDictionaryPointer(0)
		},
		"ab68ByteOrder": func(b509BoilerCodecDictionaryEntry *B509BoilerCodecDictionaryEntry) {
			b509BoilerCodecDictionaryEntry.B509BoilerCodecDictionaryByteOrder = b509BoilerCodecDictionaryPointer("")
		},
		"ab68Unit": func(b509BoilerCodecDictionaryEntry *B509BoilerCodecDictionaryEntry) {
			b509BoilerCodecDictionaryEntry.B509BoilerCodecDictionaryUnit = b509BoilerCodecDictionaryPointer("")
		},
		"ab68Converter": func(b509BoilerCodecDictionaryEntry *B509BoilerCodecDictionaryEntry) {
			b509BoilerCodecDictionaryEntry.B509BoilerCodecDictionaryConverter = b509BoilerCodecDictionaryPointer("")
		},
		"ab68Mapping": func(b509BoilerCodecDictionaryEntry *B509BoilerCodecDictionaryEntry) {
			b509BoilerCodecDictionaryEntry.B509BoilerCodecDictionaryMapping = map[string]string{"0": "ab68-zero"}
		},
		"ab68Sentinel": func(b509BoilerCodecDictionaryEntry *B509BoilerCodecDictionaryEntry) {
			b509BoilerCodecDictionaryEntry.B509BoilerCodecDictionarySentinel = b509BoilerCodecDictionaryPointer(int64(0))
		},
		"ab68UnknownCodec": func(b509BoilerCodecDictionaryEntry *B509BoilerCodecDictionaryEntry) {
			b509BoilerCodecDictionaryEntry.B509BoilerCodecDeclaration = "ab68-unknown-compound / 7"
		},
	}
	for b509BoilerCodecDictionaryName, b509BoilerCodecDictionaryMutate := range b509BoilerCodecDictionaryCases {
		b509BoilerCodecDictionaryTest.Run(b509BoilerCodecDictionaryName, func(b509BoilerCodecDictionaryTest *b509BoilerCodecDictionaryTesting.T) {
			b509BoilerCodecDictionaryEntry := b509BoilerCodecDictionaryFixture()
			b509BoilerCodecDictionaryMutate(&b509BoilerCodecDictionaryEntry)
			b509BoilerCodecDictionaryData := b509BoilerCodecDictionaryEncode(b509BoilerCodecDictionaryTest, []B509BoilerCodecDictionaryEntry{b509BoilerCodecDictionaryEntry})
			b509BoilerCodecDictionaryResult, b509BoilerCodecDictionaryOK, b509BoilerCodecDictionaryError := b509BoilerCodecDictionaryLookup(b509BoilerCodecDictionaryData, b509BoilerCodecDictionaryEntry.B509BoilerCodecDictionaryKey)
			if b509BoilerCodecDictionaryError != nil || !b509BoilerCodecDictionaryOK || !b509BoilerCodecDictionaryReflect.DeepEqual(b509BoilerCodecDictionaryResult, b509BoilerCodecDictionaryEntry) || b509BoilerCodecDictionaryReflect.DeepEqual(b509BoilerCodecDictionaryResult, b509BoilerCodecDictionaryFixture()) {
				b509BoilerCodecDictionaryTest.Fatalf("ab68 explicit zero/unknown metadata lost: %+v %v", b509BoilerCodecDictionaryResult, b509BoilerCodecDictionaryError)
			}
		})
	}
}

func TestAb68IndependentOwnership(b509BoilerCodecDictionaryTest *b509BoilerCodecDictionaryTesting.T) {
	b509BoilerCodecDictionaryKey := b509BoilerCodecDictionaryFixture().B509BoilerCodecDictionaryKey
	b509BoilerCodecDictionaryKey.B509BoilerCodecDictionarySelector = []uint16{0x2E, 0}
	b509BoilerCodecDictionaryOriginal, b509BoilerCodecDictionaryOK, b509BoilerCodecDictionaryError := B509BoilerCodecDictionaryLookup(b509BoilerCodecDictionaryKey)
	if b509BoilerCodecDictionaryError != nil || !b509BoilerCodecDictionaryOK {
		b509BoilerCodecDictionaryTest.Fatal(b509BoilerCodecDictionaryError)
	}
	b509BoilerCodecDictionaryEntry, _, _ := B509BoilerCodecDictionaryLookup(b509BoilerCodecDictionaryKey)
	*b509BoilerCodecDictionaryEntry.B509BoilerCodecDictionaryKey.B509BoilerCodecDictionaryPrimary = 0
	*b509BoilerCodecDictionaryEntry.B509BoilerCodecDictionaryKey.B509BoilerCodecDictionarySecondary = 0
	*b509BoilerCodecDictionaryEntry.B509BoilerCodecDictionaryKey.B509BoilerCodecDictionaryOperation = 0
	b509BoilerCodecDictionaryEntry.B509BoilerCodecDictionaryKey.B509BoilerCodecDictionarySelector[0] = 0
	b509BoilerCodecDictionaryEntry.B509BoilerCodecDictionaryConflicts[0] = "ab68-mutated"
	b509BoilerCodecDictionaryEntry.B509BoilerCodecDictionaryEvidence["execution"] = "Proven"
	b509BoilerCodecDictionaryEntries, b509BoilerCodecDictionaryError := B509BoilerCodecDictionaryEntries()
	if b509BoilerCodecDictionaryError != nil {
		b509BoilerCodecDictionaryTest.Fatal(b509BoilerCodecDictionaryError)
	}
	b509BoilerCodecDictionaryBaseline := b509BoilerCodecDictionaryEncode(b509BoilerCodecDictionaryTest, b509BoilerCodecDictionaryEntries)
	b509BoilerCodecDictionaryEntries[0].B509BoilerCodecDictionaryEvidence["metadata"] = "ab68-mutated"
	if b509BoilerCodecDictionaryEntries[1].B509BoilerCodecDictionaryEvidence["metadata"] != "Proven" {
		b509BoilerCodecDictionaryTest.Fatal("ab68 rows share nested state")
	}
	for b509BoilerCodecDictionaryIndex := range b509BoilerCodecDictionaryEntries {
		b509BoilerCodecDictionaryEntries[b509BoilerCodecDictionaryIndex].B509BoilerCodecDictionaryKey.B509BoilerCodecDictionarySelector[0] = 0
		b509BoilerCodecDictionaryEntries[b509BoilerCodecDictionaryIndex].B509BoilerCodecDictionaryEvidence["metadata"] = "Unknown"
	}
	b509BoilerCodecDictionaryFreshEntries, b509BoilerCodecDictionaryError := B509BoilerCodecDictionaryEntries()
	if b509BoilerCodecDictionaryError != nil || !b509BoilerCodecDictionaryBytes.Equal(b509BoilerCodecDictionaryEncode(b509BoilerCodecDictionaryTest, b509BoilerCodecDictionaryFreshEntries), b509BoilerCodecDictionaryBaseline) {
		b509BoilerCodecDictionaryTest.Fatal("ab68 entry enumeration shares nested state")
	}
	b509BoilerCodecDictionaryAgain, b509BoilerCodecDictionaryOK, b509BoilerCodecDictionaryError := B509BoilerCodecDictionaryLookup(b509BoilerCodecDictionaryKey)
	if b509BoilerCodecDictionaryError != nil || !b509BoilerCodecDictionaryOK || !b509BoilerCodecDictionaryReflect.DeepEqual(b509BoilerCodecDictionaryOriginal, b509BoilerCodecDictionaryAgain) {
		b509BoilerCodecDictionaryTest.Fatal("ab68 public result shares nested state")
	}
	b509BoilerCodecDictionaryFixtureEntry := b509BoilerCodecDictionaryFixture()
	b509BoilerCodecDictionaryFixtureEntry.B509BoilerCodecDictionaryMapping = map[string]string{"0": "ab68-zero"}
	b509BoilerCodecDictionaryFixtureEntry.B509BoilerCodecDictionaryWidthBits = b509BoilerCodecDictionaryPointer(0)
	b509BoilerCodecDictionaryFixtureEntry.B509BoilerCodecDictionaryByteOrder = b509BoilerCodecDictionaryPointer("ab68-order")
	b509BoilerCodecDictionaryFixtureEntry.B509BoilerCodecDictionaryUnit = b509BoilerCodecDictionaryPointer("ab68-unit")
	b509BoilerCodecDictionaryFixtureEntry.B509BoilerCodecDictionaryConverter = b509BoilerCodecDictionaryPointer("ab68-declaration")
	b509BoilerCodecDictionaryFixtureEntry.B509BoilerCodecDictionarySentinel = b509BoilerCodecDictionaryPointer(int64(0))
	b509BoilerCodecDictionaryData := b509BoilerCodecDictionaryEncode(b509BoilerCodecDictionaryTest, []B509BoilerCodecDictionaryEntry{b509BoilerCodecDictionaryFixtureEntry})
	b509BoilerCodecDictionaryBefore := append([]byte{}, b509BoilerCodecDictionaryData...)
	b509BoilerCodecDictionaryResult, b509BoilerCodecDictionaryOK, b509BoilerCodecDictionaryError := b509BoilerCodecDictionaryLookup(b509BoilerCodecDictionaryData, b509BoilerCodecDictionaryFixtureEntry.B509BoilerCodecDictionaryKey)
	if b509BoilerCodecDictionaryError != nil || !b509BoilerCodecDictionaryOK {
		b509BoilerCodecDictionaryTest.Fatal(b509BoilerCodecDictionaryError)
	}
	b509BoilerCodecDictionaryResult.B509BoilerCodecDictionaryMapping["0"] = "ab68-changed"
	*b509BoilerCodecDictionaryResult.B509BoilerCodecDictionaryWidthBits = 8
	*b509BoilerCodecDictionaryResult.B509BoilerCodecDictionaryByteOrder = "ab68-changed"
	*b509BoilerCodecDictionaryResult.B509BoilerCodecDictionaryUnit = "ab68-changed"
	*b509BoilerCodecDictionaryResult.B509BoilerCodecDictionaryConverter = "ab68-changed"
	*b509BoilerCodecDictionaryResult.B509BoilerCodecDictionarySentinel = 1
	b509BoilerCodecDictionaryAgain, b509BoilerCodecDictionaryOK, b509BoilerCodecDictionaryError = b509BoilerCodecDictionaryLookup(b509BoilerCodecDictionaryData, b509BoilerCodecDictionaryFixtureEntry.B509BoilerCodecDictionaryKey)
	if b509BoilerCodecDictionaryError != nil || !b509BoilerCodecDictionaryOK || !b509BoilerCodecDictionaryReflect.DeepEqual(b509BoilerCodecDictionaryAgain, b509BoilerCodecDictionaryFixtureEntry) || !b509BoilerCodecDictionaryBytes.Equal(b509BoilerCodecDictionaryBefore, b509BoilerCodecDictionaryData) {
		b509BoilerCodecDictionaryTest.Fatal("ab68 private seam shares nested state or alters raw bytes")
	}
	for b509BoilerCodecDictionaryIndex := range b509BoilerCodecDictionaryData {
		b509BoilerCodecDictionaryData[b509BoilerCodecDictionaryIndex] = ' '
	}
	if !b509BoilerCodecDictionaryReflect.DeepEqual(b509BoilerCodecDictionaryAgain, b509BoilerCodecDictionaryFixtureEntry) {
		b509BoilerCodecDictionaryTest.Fatal("ab68 parsed values alias raw source memory")
	}
}
