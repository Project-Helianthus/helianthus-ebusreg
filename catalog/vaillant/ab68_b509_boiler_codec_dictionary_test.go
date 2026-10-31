package vaillant

import (
	ab68bytes "bytes"
	ab68json "encoding/json"
	ab68reflect "reflect"
	ab68strings "strings"
	ab68testing "testing"
)

func ab68Pointer[Ab68T any](ab68Value Ab68T) *Ab68T { return &ab68Value }

func ab68Fixture() Ab68Entry {
	return Ab68Entry{
		Ab68Key:      Ab68Key{Ab68Family: "boiler", Ab68Variant: "direct-symbolic-candidate", Ab68Primary: ab68Pointer(byte(0xB5)), Ab68Secondary: ab68Pointer(byte(9)), Ab68Operation: ab68Pointer(byte(13)), Ab68Direction: "read", Ab68Selector: []uint16{0x18, 0}},
		Ab68Quantity: "flow_temperature", Ab68CodecDeclaration: "DATA2c", Ab68Grade: "Hypothesis", Ab68Conflicts: []string{},
		Ab68Evidence: map[string]string{"metadata": "Proven", "qualification": "Unknown", "physical_support": "Unknown", "presence": "Unknown", "freshness": "Unknown", "effects": "Unknown", "persistence": "Unknown", "delivery": "Unknown", "execution": "Unknown"},
	}
}

func ab68Encode(ab68Test *ab68testing.T, ab68Value any) []byte {
	ab68Test.Helper()
	ab68Data, ab68Error := ab68json.Marshal(ab68Value)
	if ab68Error != nil {
		ab68Test.Fatal(ab68Error)
	}
	return ab68Data
}

func TestAb68ExactDictionary(ab68Test *ab68testing.T) {
	ab68Expected := map[[2]uint16][3]string{
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
	ab68Entries, ab68Error := Ab68Entries()
	if ab68Error != nil {
		ab68Test.Fatal(ab68Error)
	}
	if len(ab68Entries) != 25 {
		ab68Test.Fatalf("ab68 exact set: got %d rows", len(ab68Entries))
	}
	ab68Seen := make(map[[2]uint16]bool)
	for _, ab68Entry := range ab68Entries {
		ab68Selector := [2]uint16{ab68Entry.Ab68Key.Ab68Selector[0], ab68Entry.Ab68Key.Ab68Selector[1]}
		ab68Want, ab68Found := ab68Expected[ab68Selector]
		if !ab68Found || ab68Seen[ab68Selector] {
			ab68Test.Fatalf("ab68 unexpected/duplicate selector: %04X", ab68Selector)
		}
		ab68Seen[ab68Selector] = true
		ab68Control := ab68Fixture()
		ab68Control.Ab68Key.Ab68Selector = []uint16{ab68Selector[0], ab68Selector[1]}
		ab68Control.Ab68Quantity = ab68Want[0]
		ab68Control.Ab68CodecDeclaration = ab68Want[1]
		if ab68Want[2] != "" {
			ab68Control.Ab68Conflicts = []string{ab68Want[2]}
		}
		if !ab68reflect.DeepEqual(ab68Entry, ab68Control) {
			ab68Test.Errorf("ab68 metadata for %04X: got %+v, want %+v", ab68Selector, ab68Entry, ab68Control)
		}
		ab68Result, ab68OK, ab68LookupError := Ab68Lookup(ab68Control.Ab68Key)
		if ab68LookupError != nil || !ab68OK || !ab68reflect.DeepEqual(ab68Result, ab68Control) {
			ab68Test.Errorf("ab68 lookup %04X: %+v %v %v", ab68Selector, ab68Result, ab68OK, ab68LookupError)
		}
	}
	if len(ab68Seen) != len(ab68Expected) {
		ab68Test.Fatal("ab68 missing finite tuple")
	}
}

func TestAb68KeyAxes(ab68Test *ab68testing.T) {
	ab68Cases := map[string]func(*Ab68Key){
		"ab68Family":             func(ab68Key *Ab68Key) { ab68Key.Ab68Family = "interface-unit" },
		"ab68UnknownFamily":      func(ab68Key *Ab68Key) { ab68Key.Ab68Family = "ab68-unknown" },
		"ab68Variant":            func(ab68Key *Ab68Key) { ab68Key.Ab68Variant = "ab68-other" },
		"ab68Primary":            func(ab68Key *Ab68Key) { ab68Key.Ab68Primary = ab68Pointer(byte(0)) },
		"ab68Secondary":          func(ab68Key *Ab68Key) { ab68Key.Ab68Secondary = ab68Pointer(byte(0)) },
		"ab68OperationWrite":     func(ab68Key *Ab68Key) { ab68Key.Ab68Operation = ab68Pointer(byte(0x0E)) },
		"ab68OperationAlternate": func(ab68Key *Ab68Key) { ab68Key.Ab68Operation = ab68Pointer(byte(0x29)) },
		"ab68OperationZero":      func(ab68Key *Ab68Key) { ab68Key.Ab68Operation = ab68Pointer(byte(0)) },
		"ab68DirectionRequest":   func(ab68Key *Ab68Key) { ab68Key.Ab68Direction = "request" },
		"ab68DirectionPositive":  func(ab68Key *Ab68Key) { ab68Key.Ab68Direction = "positive" },
		"ab68SelectorFirst":      func(ab68Key *Ab68Key) { ab68Key.Ab68Selector[0]++ },
		"ab68SelectorSecond":     func(ab68Key *Ab68Key) { ab68Key.Ab68Selector[1]++ },
		"ab68SelectorZero":       func(ab68Key *Ab68Key) { ab68Key.Ab68Selector = []uint16{0, 0} },
		"ab68SelectorReversed":   func(ab68Key *Ab68Key) { ab68Key.Ab68Selector = []uint16{0, 0x18} },
	}
	for ab68Name, ab68Mutate := range ab68Cases {
		ab68Test.Run(ab68Name, func(ab68Test *ab68testing.T) {
			ab68First, ab68Second := ab68Fixture(), ab68Fixture()
			ab68Mutate(&ab68Second.Ab68Key)
			ab68Result, ab68OK, ab68Error := Ab68Lookup(ab68Second.Ab68Key)
			if ab68Error != nil || ab68OK || !ab68reflect.DeepEqual(ab68Result, Ab68Entry{}) {
				ab68Test.Fatalf("ab68 unexpected admission: %+v %v %v", ab68Result, ab68OK, ab68Error)
			}
			ab68Second.Ab68Quantity = "ab68-synthetic-control"
			ab68Data := ab68Encode(ab68Test, []Ab68Entry{ab68First, ab68Second})
			for _, ab68Control := range []Ab68Entry{ab68First, ab68Second} {
				ab68Result, ab68OK, ab68Error = ab68Lookup(ab68Data, ab68Control.Ab68Key)
				if ab68Error != nil || !ab68OK || !ab68reflect.DeepEqual(ab68Result, ab68Control) {
					ab68Test.Fatalf("ab68 independent key control: %+v %v %v", ab68Result, ab68OK, ab68Error)
				}
			}
		})
	}
	ab68Absent := map[string]func(*Ab68Key){
		"ab68Family":    func(ab68Key *Ab68Key) { ab68Key.Ab68Family = "" },
		"ab68Variant":   func(ab68Key *Ab68Key) { ab68Key.Ab68Variant = "" },
		"ab68Primary":   func(ab68Key *Ab68Key) { ab68Key.Ab68Primary = nil },
		"ab68Secondary": func(ab68Key *Ab68Key) { ab68Key.Ab68Secondary = nil },
		"ab68Operation": func(ab68Key *Ab68Key) { ab68Key.Ab68Operation = nil },
		"ab68Direction": func(ab68Key *Ab68Key) { ab68Key.Ab68Direction = "" },
		"ab68Selector":  func(ab68Key *Ab68Key) { ab68Key.Ab68Selector = nil },
	}
	for ab68Name, ab68Mutate := range ab68Absent {
		ab68Test.Run(ab68Name+"Absent", func(ab68Test *ab68testing.T) {
			ab68Entry := ab68Fixture()
			ab68Mutate(&ab68Entry.Ab68Key)
			if _, ab68OK, ab68Error := Ab68Lookup(ab68Entry.Ab68Key); ab68OK || ab68Error != nil {
				ab68Test.Fatal("ab68 absent key matched")
			}
			if ab68Entries, ab68Error := ab68Parse(ab68Encode(ab68Test, []Ab68Entry{ab68Entry})); ab68Error == nil || ab68Entries != nil {
				ab68Test.Fatal("ab68 incomplete source admitted")
			}
		})
	}
}

func TestAb68Validation(ab68Test *ab68testing.T) {
	ab68Valid := ab68Encode(ab68Test, []Ab68Entry{ab68Fixture()})
	ab68Cases := map[string][]byte{
		"ab68Malformed": []byte("["), "ab68Null": []byte("null"), "ab68Empty": []byte("[]"), "ab68Partial": []byte("[{}]"),
		"ab68Trailing":          append(append([]byte{}, ab68Valid...), []byte(" {}")...),
		"ab68TrailingMalformed": append(append([]byte{}, ab68Valid...), byte('x')),
		"ab68UnknownField":      []byte(ab68strings.Replace(string(ab68Valid), `"quantity":`, `"ab68_extra":0,"quantity":`, 1)),
		"ab68Duplicate":         ab68Encode(ab68Test, []Ab68Entry{ab68Fixture(), ab68Fixture()}),
	}
	for ab68Name, ab68Data := range ab68Cases {
		ab68Test.Run(ab68Name, func(ab68Test *ab68testing.T) {
			ab68Before := append([]byte{}, ab68Data...)
			if ab68Entries, ab68Error := ab68Parse(ab68Data); ab68Error == nil || ab68Entries != nil {
				ab68Test.Fatal("ab68 invalid data admitted")
			}
			if ab68Result, ab68OK, ab68Error := ab68Lookup(ab68Data, ab68Fixture().Ab68Key); ab68Error == nil || ab68OK || !ab68reflect.DeepEqual(ab68Result, Ab68Entry{}) {
				ab68Test.Fatal("ab68 parse error not propagated")
			}
			if !ab68bytes.Equal(ab68Data, ab68Before) {
				ab68Test.Fatal("ab68 raw source changed on failure")
			}
		})
	}
	ab68Invalid := map[string]func(*Ab68Entry){
		"ab68ShortSelector": func(ab68Entry *Ab68Entry) { ab68Entry.Ab68Key.Ab68Selector = []uint16{24} },
		"ab68LongSelector":  func(ab68Entry *Ab68Entry) { ab68Entry.Ab68Key.Ab68Selector = []uint16{24, 0, 0} },
		"ab68WideSelector":  func(ab68Entry *Ab68Entry) { ab68Entry.Ab68Key.Ab68Selector[0] = 256 },
		"ab68Quantity":      func(ab68Entry *Ab68Entry) { ab68Entry.Ab68Quantity = "" },
		"ab68Codec":         func(ab68Entry *Ab68Entry) { ab68Entry.Ab68CodecDeclaration = "" },
		"ab68Grade":         func(ab68Entry *Ab68Entry) { ab68Entry.Ab68Grade = "Proven" },
		"ab68Evidence":      func(ab68Entry *Ab68Entry) { ab68Entry.Ab68Evidence = nil },
	}
	for _, ab68Axis := range []string{"metadata", "qualification", "physical_support", "presence", "freshness", "effects", "persistence", "delivery", "execution"} {
		ab68Axis := ab68Axis
		ab68Invalid["ab68Evidence-"+ab68Axis] = func(ab68Entry *Ab68Entry) { ab68Entry.Ab68Evidence[ab68Axis] = "Hypothesis" }
	}
	for ab68Name, ab68Mutate := range ab68Invalid {
		ab68Test.Run(ab68Name, func(ab68Test *ab68testing.T) {
			ab68Entry := ab68Fixture()
			ab68Mutate(&ab68Entry)
			if ab68Entries, ab68Error := ab68Parse(ab68Encode(ab68Test, []Ab68Entry{ab68Entry})); ab68Error == nil || ab68Entries != nil {
				ab68Test.Fatal("ab68 invalid entry admitted")
			}
		})
	}
	if ab68Entries, ab68Error := ab68Parse(ab68Valid); ab68Error != nil || len(ab68Entries) != 1 {
		ab68Test.Fatalf("ab68 valid control: %v", ab68Error)
	}
}

func TestAb68OptionalFieldsAndUnknownCodec(ab68Test *ab68testing.T) {
	ab68Cases := map[string]func(*Ab68Entry){
		"ab68Width":        func(ab68Entry *Ab68Entry) { ab68Entry.Ab68WidthBits = ab68Pointer(0) },
		"ab68ByteOrder":    func(ab68Entry *Ab68Entry) { ab68Entry.Ab68ByteOrder = ab68Pointer("") },
		"ab68Unit":         func(ab68Entry *Ab68Entry) { ab68Entry.Ab68Unit = ab68Pointer("") },
		"ab68Converter":    func(ab68Entry *Ab68Entry) { ab68Entry.Ab68Converter = ab68Pointer("") },
		"ab68Mapping":      func(ab68Entry *Ab68Entry) { ab68Entry.Ab68Mapping = map[string]string{"0": "ab68-zero"} },
		"ab68Sentinel":     func(ab68Entry *Ab68Entry) { ab68Entry.Ab68Sentinel = ab68Pointer(int64(0)) },
		"ab68UnknownCodec": func(ab68Entry *Ab68Entry) { ab68Entry.Ab68CodecDeclaration = "ab68-unknown-compound / 7" },
	}
	for ab68Name, ab68Mutate := range ab68Cases {
		ab68Test.Run(ab68Name, func(ab68Test *ab68testing.T) {
			ab68Entry := ab68Fixture()
			ab68Mutate(&ab68Entry)
			ab68Data := ab68Encode(ab68Test, []Ab68Entry{ab68Entry})
			ab68Result, ab68OK, ab68Error := ab68Lookup(ab68Data, ab68Entry.Ab68Key)
			if ab68Error != nil || !ab68OK || !ab68reflect.DeepEqual(ab68Result, ab68Entry) || ab68reflect.DeepEqual(ab68Result, ab68Fixture()) {
				ab68Test.Fatalf("ab68 explicit zero/unknown metadata lost: %+v %v", ab68Result, ab68Error)
			}
		})
	}
}

func TestAb68IndependentOwnership(ab68Test *ab68testing.T) {
	ab68Key := ab68Fixture().Ab68Key
	ab68Key.Ab68Selector = []uint16{0x2E, 0}
	ab68Original, ab68OK, ab68Error := Ab68Lookup(ab68Key)
	if ab68Error != nil || !ab68OK {
		ab68Test.Fatal(ab68Error)
	}
	ab68Entry, _, _ := Ab68Lookup(ab68Key)
	*ab68Entry.Ab68Key.Ab68Primary = 0
	*ab68Entry.Ab68Key.Ab68Secondary = 0
	*ab68Entry.Ab68Key.Ab68Operation = 0
	ab68Entry.Ab68Key.Ab68Selector[0] = 0
	ab68Entry.Ab68Conflicts[0] = "ab68-mutated"
	ab68Entry.Ab68Evidence["execution"] = "Proven"
	ab68Entries, ab68Error := Ab68Entries()
	if ab68Error != nil {
		ab68Test.Fatal(ab68Error)
	}
	ab68Baseline := ab68Encode(ab68Test, ab68Entries)
	ab68Entries[0].Ab68Evidence["metadata"] = "ab68-mutated"
	if ab68Entries[1].Ab68Evidence["metadata"] != "Proven" {
		ab68Test.Fatal("ab68 rows share nested state")
	}
	for ab68Index := range ab68Entries {
		ab68Entries[ab68Index].Ab68Key.Ab68Selector[0] = 0
		ab68Entries[ab68Index].Ab68Evidence["metadata"] = "Unknown"
	}
	ab68FreshEntries, ab68Error := Ab68Entries()
	if ab68Error != nil || !ab68bytes.Equal(ab68Encode(ab68Test, ab68FreshEntries), ab68Baseline) {
		ab68Test.Fatal("ab68 entry enumeration shares nested state")
	}
	ab68Again, ab68OK, ab68Error := Ab68Lookup(ab68Key)
	if ab68Error != nil || !ab68OK || !ab68reflect.DeepEqual(ab68Original, ab68Again) {
		ab68Test.Fatal("ab68 public result shares nested state")
	}
	ab68FixtureEntry := ab68Fixture()
	ab68FixtureEntry.Ab68Mapping = map[string]string{"0": "ab68-zero"}
	ab68FixtureEntry.Ab68WidthBits = ab68Pointer(0)
	ab68FixtureEntry.Ab68ByteOrder = ab68Pointer("ab68-order")
	ab68FixtureEntry.Ab68Unit = ab68Pointer("ab68-unit")
	ab68FixtureEntry.Ab68Converter = ab68Pointer("ab68-declaration")
	ab68FixtureEntry.Ab68Sentinel = ab68Pointer(int64(0))
	ab68Data := ab68Encode(ab68Test, []Ab68Entry{ab68FixtureEntry})
	ab68Before := append([]byte{}, ab68Data...)
	ab68Result, ab68OK, ab68Error := ab68Lookup(ab68Data, ab68FixtureEntry.Ab68Key)
	if ab68Error != nil || !ab68OK {
		ab68Test.Fatal(ab68Error)
	}
	ab68Result.Ab68Mapping["0"] = "ab68-changed"
	*ab68Result.Ab68WidthBits = 8
	*ab68Result.Ab68ByteOrder = "ab68-changed"
	*ab68Result.Ab68Unit = "ab68-changed"
	*ab68Result.Ab68Converter = "ab68-changed"
	*ab68Result.Ab68Sentinel = 1
	ab68Again, ab68OK, ab68Error = ab68Lookup(ab68Data, ab68FixtureEntry.Ab68Key)
	if ab68Error != nil || !ab68OK || !ab68reflect.DeepEqual(ab68Again, ab68FixtureEntry) || !ab68bytes.Equal(ab68Before, ab68Data) {
		ab68Test.Fatal("ab68 private seam shares nested state or alters raw bytes")
	}
	for ab68Index := range ab68Data {
		ab68Data[ab68Index] = ' '
	}
	if !ab68reflect.DeepEqual(ab68Again, ab68FixtureEntry) {
		ab68Test.Fatal("ab68 parsed values alias raw source memory")
	}
}
