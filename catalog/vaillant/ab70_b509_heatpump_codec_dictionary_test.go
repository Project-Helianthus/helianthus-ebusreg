package vaillant

import (
	ab70bytes "bytes"
	ab70json "encoding/json"
	ab70reflect "reflect"
	ab70strings "strings"
	ab70testing "testing"
)

func ab70Pointer[ab70T any](ab70Value ab70T) *ab70T { return &ab70Value }

func ab70Key(ab70First, ab70Second uint16) Ab70Key {
	return Ab70Key{"heat-pump", "direct-symbolic-candidate", ab70Pointer(byte(181)), ab70Pointer(byte(9)), ab70Pointer(byte(13)), "read", []uint16{ab70First, ab70Second}}
}

func TestAb70ExactDictionary(ab70T *ab70testing.T) {
	ab70Expected := []struct {
		ab70First    uint16
		ab70Second   uint16
		ab70Quantity string
		ab70Codec    string
		ab70Unit     string
		ab70Conflict string
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
	ab70Entries, ab70Error := Ab70Entries()
	if ab70Error != nil || len(ab70Entries) != len(ab70Expected) {
		ab70T.Fatalf("ab70 exact count: %d %v", len(ab70Entries), ab70Error)
	}
	ab70Seen := make(map[[2]uint16]bool)
	for _, ab70Want := range ab70Expected {
		ab70Key := ab70Key(ab70Want.ab70First, ab70Want.ab70Second)
		ab70Got, ab70Found, ab70Error := Ab70Lookup(ab70Key)
		if ab70Error != nil || !ab70Found {
			ab70T.Fatalf("ab70 missing %v: %v", ab70Key, ab70Error)
		}
		ab70Entry := Ab70Entry{
			Ab70Key: ab70Key, Ab70Quantity: ab70Want.ab70Quantity, Ab70CodecDeclaration: ab70Want.ab70Codec,
			Ab70Grade: "Hypothesis", Ab70SoftwareContract: "Proven", Ab70Qualification: "Unknown", Ab70Presence: "Unknown", Ab70Freshness: "Unknown", Ab70Execution: "Unknown", Ab70PhysicalEffects: "Unknown", Ab70Persistence: "Unknown", Ab70Delivery: "Unknown",
			Ab70Notes: []string{"Literal codec declaration; no executable converter."},
		}
		if ab70Want.ab70Unit != "" {
			ab70Entry.Ab70Unit = ab70Pointer(ab70Want.ab70Unit)
		}
		if ab70Want.ab70First == 0xE2 {
			ab70Entry.Ab70WidthBits = ab70Pointer(8)
			ab70Entry.Ab70Signedness = ab70Pointer("unsigned")
			ab70Entry.Ab70DeclaredScale = ab70Pointer(100)
		}
		if ab70Want.ab70First == 0xE4 && ab70Want.ab70Second == 1 {
			ab70Entry.Ab70Signedness = ab70Pointer("signed")
		}
		if ab70Want.ab70Conflict != "" {
			ab70Entry.Ab70Notes = append(ab70Entry.Ab70Notes, ab70Want.ab70Conflict)
		}
		if !ab70reflect.DeepEqual(ab70Entry, ab70Got) {
			ab70T.Errorf("ab70 selector %v got %#v want %#v", ab70Key.Ab70Selector, ab70Got, ab70Entry)
		}
		ab70Seen[[2]uint16{ab70Want.ab70First, ab70Want.ab70Second}] = true
	}
	for _, ab70Entry := range ab70Entries {
		if len(ab70Entry.Ab70Key.Ab70Selector) != 2 || !ab70Seen[[2]uint16{ab70Entry.Ab70Key.Ab70Selector[0], ab70Entry.Ab70Key.Ab70Selector[1]}] {
			ab70T.Fatal("ab70 extra tuple")
		}
		delete(ab70Seen, [2]uint16{ab70Entry.Ab70Key.Ab70Selector[0], ab70Entry.Ab70Key.Ab70Selector[1]})
	}
}

func TestAb70KeyAxes(ab70T *ab70testing.T) {
	ab70Changes := []struct {
		ab70Name   string
		ab70Change func(*Ab70Key)
	}{
		{"ab70Family", func(ab70K *Ab70Key) { ab70K.Ab70Family = "heat-pump-interface-unit" }},
		{"ab70Variant", func(ab70K *Ab70Key) { ab70K.Ab70Variant = "unknown" }},
		{"ab70Primary", func(ab70K *Ab70Key) { ab70K.Ab70Primary = ab70Pointer(byte(0)) }},
		{"ab70Secondary", func(ab70K *Ab70Key) { ab70K.Ab70Secondary = ab70Pointer(byte(0)) }},
		{"ab70Operation", func(ab70K *Ab70Key) { ab70K.Ab70Operation = ab70Pointer(byte(14)) }},
		{"ab70Alternate", func(ab70K *Ab70Key) { ab70K.Ab70Operation = ab70Pointer(byte(41)) }},
		{"ab70Request", func(ab70K *Ab70Key) { ab70K.Ab70Direction = "request" }},
		{"ab70Positive", func(ab70K *Ab70Key) { ab70K.Ab70Direction = "positive" }},
		{"ab70First", func(ab70K *Ab70Key) { ab70K.Ab70Selector[0] = 0 }},
		{"ab70Second", func(ab70K *Ab70Key) { ab70K.Ab70Selector[1] = 1 }},
		{"ab70Reverse", func(ab70K *Ab70Key) { ab70K.Ab70Selector = []uint16{0, 0xD0} }},
		{"ab70AbsentPrimary", func(ab70K *Ab70Key) { ab70K.Ab70Primary = nil }},
		{"ab70AbsentSecondary", func(ab70K *Ab70Key) { ab70K.Ab70Secondary = nil }},
		{"ab70AbsentOperation", func(ab70K *Ab70Key) { ab70K.Ab70Operation = nil }},
		{"ab70AbsentSelector", func(ab70K *Ab70Key) { ab70K.Ab70Selector = nil }},
		{"ab70PartialSelector", func(ab70K *Ab70Key) { ab70K.Ab70Selector = []uint16{0xD0} }},
		{"ab70ExtraSelector", func(ab70K *Ab70Key) { ab70K.Ab70Selector = []uint16{0xD0, 0, 0} }},
	}
	ab70Original, _, _ := Ab70Lookup(ab70Key(0xD0, 0))
	for _, ab70Case := range ab70Changes {
		ab70T.Run(ab70Case.ab70Name, func(ab70T *ab70testing.T) {
			ab70Changed := ab70Key(0xD0, 0)
			ab70Case.ab70Change(&ab70Changed)
			if _, ab70Found, ab70Error := Ab70Lookup(ab70Changed); ab70Found || ab70Error != nil {
				ab70T.Fatalf("ab70 admitted altered key: %v", ab70Error)
			}
			ab70Control := ab70Original
			ab70Control.Ab70Key = ab70Changed
			if _, ab70Found, ab70Error := ab70Find([]Ab70Entry{ab70Control}, ab70Changed); !ab70Found || ab70Error != nil {
				ab70T.Fatal("ab70 exact changed control failed")
			}
			if _, ab70Found, _ := ab70Find([]Ab70Entry{ab70Control}, ab70Key(0xD0, 0)); ab70Found {
				ab70T.Fatal("ab70 omitted key axis")
			}
		})
	}
}

func TestAb70ParserAndOwnership(ab70T *ab70testing.T) {
	ab70Entries, ab70Error := Ab70Entries()
	if ab70Error != nil {
		ab70T.Fatal(ab70Error)
	}
	ab70Entry := ab70Entries[0]
	ab70Entry.Ab70CodecDeclaration = "ab70-unknown-label"
	ab70Entry.Ab70Mapping = map[string]string{"ab70-zero": "ab70-native"}
	ab70Entry.Ab70WidthBits = ab70Pointer(0)
	ab70Entry.Ab70Sentinel = ab70Pointer(0)
	ab70Entry.Ab70Unit = ab70Pointer("")
	ab70Data, ab70Error := ab70json.Marshal([]Ab70Entry{ab70Entry})
	if ab70Error != nil {
		ab70T.Fatal(ab70Error)
	}
	ab70Original := append([]byte(nil), ab70Data...)
	ab70Parsed, ab70Error := ab70Parse(ab70Data)
	if ab70Error != nil || !ab70reflect.DeepEqual(ab70Parsed, []Ab70Entry{ab70Entry}) || !ab70bytes.Equal(ab70Data, ab70Original) {
		ab70T.Fatal("ab70 raw or metadata-only control failed", ab70Error)
	}
	ab70Copy, ab70Found, ab70Error := ab70Find(ab70Parsed, ab70Entry.Ab70Key)
	if !ab70Found || ab70Error != nil {
		ab70T.Fatal("ab70 copy failed")
	}
	ab70Copy.Ab70Key.Ab70Selector[0] = 0
	*ab70Copy.Ab70Key.Ab70Primary = 0
	*ab70Copy.Ab70Key.Ab70Secondary = 0
	*ab70Copy.Ab70Key.Ab70Operation = 0
	*ab70Copy.Ab70WidthBits = 9
	*ab70Copy.Ab70Sentinel = 9
	*ab70Copy.Ab70Unit = "ab70-changed"
	ab70Copy.Ab70Mapping["ab70-zero"] = "ab70-changed"
	ab70Copy.Ab70Notes[0] = "ab70-changed"
	if !ab70reflect.DeepEqual(ab70Parsed, []Ab70Entry{ab70Entry}) {
		ab70T.Fatal("ab70 nested data shared")
	}
	for ab70Index := range ab70Data {
		ab70Data[ab70Index] = ' '
	}
	if !ab70reflect.DeepEqual(ab70Parsed, []Ab70Entry{ab70Entry}) {
		ab70T.Fatal("ab70 input memory shared")
	}
	ab70Entries[0].Ab70Notes[0] = "ab70-changed"
	ab70Entries[0].Ab70Key.Ab70Selector[0] = 0
	ab70Again, _, _ := Ab70Lookup(ab70Key(0xD0, 0))
	if ab70Again.Ab70Notes[0] != "Literal codec declaration; no executable converter." {
		ab70T.Fatal("ab70 entries shared")
	}
	ab70Again.Ab70Notes[0] = "ab70-changed"
	*ab70Again.Ab70Key.Ab70Primary = 0
	ab70Again, ab70Found, ab70Error = Ab70Lookup(ab70Key(0xD0, 0))
	if !ab70Found || ab70Error != nil || ab70Again.Ab70Notes[0] != "Literal codec declaration; no executable converter." {
		ab70T.Fatal("ab70 repeated lookup shared")
	}
	ab70Duplicate, _ := ab70json.Marshal([]Ab70Entry{ab70Entry, ab70Entry})
	for _, ab70Bad := range []string{"{", "null", "[]", "[{}]", string(ab70Duplicate), string(ab70Original) + " []", ab70strings.Replace(string(ab70Original), `"quantity":`, `"ab70_extra":0,"quantity":`, 1), ab70strings.Replace(string(ab70Original), `"selector":[208,0]`, `"selector":[208]`, 1), ab70strings.Replace(string(ab70Original), `"selector":[208,0]`, `"selector":[208,0,0]`, 1), ab70strings.Replace(string(ab70Original), `"selector":[208,0]`, `"selector":[256,0]`, 1)} {
		if ab70Result, ab70Error := ab70Parse([]byte(ab70Bad)); ab70Error == nil || ab70Result != nil {
			ab70T.Fatal("ab70 malformed/duplicate admitted", ab70Bad)
		}
	}
	ab70Value := ab70reflect.ValueOf(&ab70Entry).Elem()
	for ab70Index := 0; ab70Index < ab70Value.NumField(); ab70Index++ {
		ab70Field := ab70Value.Type().Field(ab70Index)
		if ab70Field.Name == "Ab70Quantity" || ab70Field.Name == "Ab70CodecDeclaration" || ab70Field.Name == "Ab70Grade" || ab70Field.Name == "Ab70SoftwareContract" || ab70strings.Contains("Ab70Qualification Ab70Presence Ab70Freshness Ab70Execution Ab70PhysicalEffects Ab70Persistence Ab70Delivery", ab70Field.Name) {
			ab70Saved := ab70Value.Field(ab70Index).String()
			ab70Value.Field(ab70Index).SetString("")
			ab70Bad, _ := ab70json.Marshal([]Ab70Entry{ab70Entry})
			if _, ab70Error := ab70Parse(ab70Bad); ab70Error == nil {
				ab70T.Errorf("ab70 missing %s admitted", ab70Field.Name)
			}
			ab70Value.Field(ab70Index).SetString(ab70Saved)
		}
	}
	ab70Absent := ab70Entry
	ab70Absent.Ab70WidthBits = nil
	ab70Absent.Ab70Sentinel = nil
	ab70Absent.Ab70Unit = nil
	if ab70reflect.DeepEqual(ab70Absent, ab70Entry) {
		ab70T.Fatal("ab70 absent equals explicit zero")
	}
}

func TestAb70OptionalAndPartialAxes(ab70T *ab70testing.T) {
	ab70Entries, ab70Error := Ab70Entries()
	if ab70Error != nil {
		ab70T.Fatal(ab70Error)
	}
	ab70Optional := []struct {
		ab70Name   string
		ab70Change func(*Ab70Entry)
	}{
		{"ab70Unit", func(ab70E *Ab70Entry) { ab70E.Ab70Unit = ab70Pointer("") }},
		{"ab70Width", func(ab70E *Ab70Entry) { ab70E.Ab70WidthBits = ab70Pointer(0) }},
		{"ab70Scale", func(ab70E *Ab70Entry) { ab70E.Ab70DeclaredScale = ab70Pointer(0) }},
		{"ab70Sentinel", func(ab70E *Ab70Entry) { ab70E.Ab70Sentinel = ab70Pointer(0) }},
		{"ab70Signedness", func(ab70E *Ab70Entry) { ab70E.Ab70Signedness = ab70Pointer("") }},
		{"ab70ByteOrder", func(ab70E *Ab70Entry) { ab70E.Ab70ByteOrder = ab70Pointer("") }},
		{"ab70Converter", func(ab70E *Ab70Entry) { ab70E.Ab70Converter = ab70Pointer("ab70-opaque") }},
		{"ab70Mapping", func(ab70E *Ab70Entry) { ab70E.Ab70Mapping = map[string]string{"ab70-zero": "ab70-label"} }},
	}
	for _, ab70Case := range ab70Optional {
		ab70T.Run(ab70Case.ab70Name, func(ab70T *ab70testing.T) {
			ab70Entry := ab70Entries[0]
			ab70Case.ab70Change(&ab70Entry)
			ab70Data, _ := ab70json.Marshal([]Ab70Entry{ab70Entry})
			ab70Parsed, ab70Error := ab70Parse(ab70Data)
			if ab70Error != nil || !ab70reflect.DeepEqual(ab70Parsed, []Ab70Entry{ab70Entry}) || ab70reflect.DeepEqual(ab70Parsed, ab70Entries[:1]) {
				ab70T.Fatal("ab70 optional field lost", ab70Error)
			}
			ab70Copy, _, _ := ab70Find(ab70Parsed, ab70Entry.Ab70Key)
			ab70Value := ab70reflect.ValueOf(&ab70Copy).Elem()
			for ab70Index := 0; ab70Index < ab70Value.NumField(); ab70Index++ {
				ab70Field := ab70Value.Field(ab70Index)
				if ab70Field.Kind() == ab70reflect.Pointer && !ab70Field.IsNil() {
					if ab70Field.Elem().Kind() == ab70reflect.String {
						ab70Field.Elem().SetString("ab70-mutated")
					} else {
						ab70Field.Elem().SetInt(7)
					}
				}
			}
			if ab70Copy.Ab70Mapping != nil {
				ab70Copy.Ab70Mapping["ab70-zero"] = "ab70-mutated"
			}
			if !ab70reflect.DeepEqual(ab70Parsed, []Ab70Entry{ab70Entry}) {
				ab70T.Fatal("ab70 optional pointer shared")
			}
		})
	}
	ab70KeyChanges := []func(*Ab70Key){
		func(ab70K *Ab70Key) { ab70K.Ab70Family = "" },
		func(ab70K *Ab70Key) { ab70K.Ab70Variant = "" },
		func(ab70K *Ab70Key) { ab70K.Ab70Direction = "" },
		func(ab70K *Ab70Key) { ab70K.Ab70Primary = nil },
		func(ab70K *Ab70Key) { ab70K.Ab70Secondary = nil },
		func(ab70K *Ab70Key) { ab70K.Ab70Operation = nil },
	}
	for ab70Index, ab70Change := range ab70KeyChanges {
		ab70Entry := ab70Entries[0]
		ab70Change(&ab70Entry.Ab70Key)
		ab70Data, _ := ab70json.Marshal([]Ab70Entry{ab70Entry})
		if ab70Parsed, ab70Error := ab70Parse(ab70Data); ab70Error == nil || ab70Parsed != nil {
			ab70T.Fatalf("ab70 partial axis %d admitted", ab70Index)
		}
	}
	for _, ab70Change := range []func(*Ab70Key){
		func(ab70K *Ab70Key) { ab70K.Ab70Primary = ab70Pointer(byte(0)) },
		func(ab70K *Ab70Key) { ab70K.Ab70Secondary = ab70Pointer(byte(0)) },
		func(ab70K *Ab70Key) { ab70K.Ab70Operation = ab70Pointer(byte(0)) },
		func(ab70K *Ab70Key) { ab70K.Ab70Selector = []uint16{0, 0} },
	} {
		ab70Entry := ab70Entries[0]
		ab70Change(&ab70Entry.Ab70Key)
		ab70Data, _ := ab70json.Marshal([]Ab70Entry{ab70Entry})
		ab70Parsed, ab70Error := ab70Parse(ab70Data)
		if ab70Error != nil || !ab70reflect.DeepEqual(ab70Parsed, []Ab70Entry{ab70Entry}) {
			ab70T.Fatal("ab70 numeric zero rejected", ab70Error)
		}
	}
}

func TestAb70FailedReadPreservesOwnedData(ab70T *ab70testing.T) {
	ab70Saved := ab70JSON
	defer func() { ab70JSON = ab70Saved }()
	ab70Good, ab70Found, ab70Error := Ab70Lookup(ab70Key(0xBC, 0))
	if !ab70Found || ab70Error != nil {
		ab70T.Fatal("ab70 valid control failed")
	}
	ab70Snapshot, ab70Error := ab70json.Marshal(ab70Good)
	if ab70Error != nil {
		ab70T.Fatal(ab70Error)
	}
	ab70JSON = []byte(`[{"key":{}}]`)
	if ab70Rows, ab70Error := Ab70Entries(); ab70Error == nil || ab70Rows != nil {
		ab70T.Fatal("ab70 partial data admitted")
	}
	if ab70Row, ab70Found, ab70Error := Ab70Lookup(ab70Key(0xBC, 0)); ab70Error == nil || ab70Found || !ab70reflect.DeepEqual(ab70Row, Ab70Entry{}) {
		ab70T.Fatal("ab70 lookup error lost")
	}
	ab70After, _ := ab70json.Marshal(ab70Good)
	if !ab70bytes.Equal(ab70Snapshot, ab70After) {
		ab70T.Fatal("ab70 last owned evidence changed after failure")
	}
	ab70JSON = ab70Saved
	ab70Again, ab70Found, ab70Error := Ab70Lookup(ab70Key(0xBC, 0))
	if !ab70Found || ab70Error != nil || !ab70reflect.DeepEqual(ab70Good, ab70Again) {
		ab70T.Fatal("ab70 recovery control failed")
	}
}
