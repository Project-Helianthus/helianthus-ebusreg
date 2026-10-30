package vaillant

import (
	ab99json "encoding/json"
	ab99reflect "reflect"
	ab99testing "testing"
)

func ab99Pointer[ab99T any](ab99Value ab99T) *ab99T { return &ab99Value }
func ab99Key(ab99Variant string, ab99Operation uint8, ab99Direction, ab99Field string) Ab99Key {
	return Ab99Key{Ab99Family: "vaillant", Ab99Variant: ab99Variant, Ab99Primary: 0xb5, Ab99Secondary: ab99Pointer(uint8(0x23)), Ab99Operation: ab99Pointer(ab99Operation), Ab99Direction: ab99Direction, Ab99Field: ab99Field}
}
func ab99Get(ab99T *ab99testing.T, ab99Key Ab99Key) Ab99Entry {
	ab99T.Helper()
	ab99Entry, ab99OK, ab99Error := Ab99Lookup(ab99Key)
	if ab99Error != nil || !ab99OK {
		ab99T.Fatalf("lookup: %v, %v", ab99OK, ab99Error)
	}
	return ab99Entry
}
func TestAb99Finite(ab99T *ab99testing.T) {
	ab99Entries, ab99Error := Ab99Entries()
	if ab99Error != nil || len(ab99Entries) != 34 {
		ab99T.Fatalf("finite catalog: %d %v", len(ab99Entries), ab99Error)
	}
	for _, ab99Entry := range ab99Entries {
		if ab99Entry.Ab99Grade != "Hypothesis" || ab99Entry.Ab99Support != "Unknown" || ab99Entry.Ab99PhysicalAdaptation != "Unknown" || ab99Entry.Ab99ExecutedMovement != "Unknown" || ab99Entry.Ab99SelectedInterpretation != "Unknown" || ab99Entry.Ab99UnknownMetadata != "Unknown" || !ab99Entry.Ab99PreserveRaw {
			ab99T.Fatal("grade/admission")
		}
		if !ab99reflect.DeepEqual(ab99Get(ab99T, ab99Entry.Ab99Key), ab99Entry) {
			ab99T.Fatal("lookup/accessor disagreement")
		}
		if ab99Entry.Ab99Key.Ab99Group != nil || ab99Entry.Ab99Key.Ab99RequestedInstance != nil || ab99Entry.Ab99Key.Ab99Register != nil || ab99Entry.Ab99Key.Ab99CommandRole != nil {
			ab99T.Fatal("invented B524 axes")
		}
	}
}
func TestAb99ExactAxes(ab99T *ab99testing.T) {
	ab99Control := ab99Key("configuration", 0, "request", "S1")
	ab99Changes := map[string]func(*Ab99Key){
		"ab99_family":           func(ab99Key *Ab99Key) { ab99Key.Ab99Family = "other" },
		"ab99_variant":          func(ab99Key *Ab99Key) { ab99Key.Ab99Variant = "other" },
		"ab99_primary":          func(ab99Key *Ab99Key) { ab99Key.Ab99Primary = 0 },
		"ab99_secondary":        func(ab99Key *Ab99Key) { ab99Key.Ab99Secondary = ab99Pointer(uint8(0)) },
		"ab99_secondary_absent": func(ab99Key *Ab99Key) { ab99Key.Ab99Secondary = nil },
		"ab99_operation":        func(ab99Key *Ab99Key) { ab99Key.Ab99Operation = ab99Pointer(uint8(1)) },
		"ab99_operation_absent": func(ab99Key *Ab99Key) { ab99Key.Ab99Operation = nil },
		"ab99_direction":        func(ab99Key *Ab99Key) { ab99Key.Ab99Direction = "positive" },
		"ab99_field":            func(ab99Key *Ab99Key) { ab99Key.Ab99Field = "other" },
		"ab99_group_zero":       func(ab99Key *Ab99Key) { ab99Key.Ab99Group = ab99Pointer(uint8(0)) },
		"ab99_instance_zero":    func(ab99Key *Ab99Key) { ab99Key.Ab99RequestedInstance = ab99Pointer(uint8(0)) },
		"ab99_register_zero":    func(ab99Key *Ab99Key) { ab99Key.Ab99Register = ab99Pointer(uint16(0)) },
		"ab99_role":             func(ab99Key *Ab99Key) { ab99Key.Ab99CommandRole = ab99Pointer("read") },
	}
	for ab99Name, ab99Change := range ab99Changes {
		ab99T.Run(ab99Name, func(ab99T *ab99testing.T) {
			ab99Get(ab99T, ab99Control)
			ab99Changed := ab99Control
			ab99Change(&ab99Changed)
			ab99Entry, ab99OK, ab99Error := Ab99Lookup(ab99Changed)
			if ab99Error != nil || ab99OK || !ab99reflect.DeepEqual(ab99Entry, Ab99Entry{}) {
				ab99T.Fatal("key axis ignored")
			}
		})
	}
}
func TestAb99OptionalAxes(ab99T *ab99testing.T) {
	ab99Saved := ab99Data
	ab99T.Cleanup(func() { ab99Data = ab99Saved })
	ab99Entries, ab99Error := Ab99Entries()
	if ab99Error != nil {
		ab99T.Fatal(ab99Error)
	}
	ab99Control := ab99Entries[0]
	ab99Control.Ab99Key.Ab99Group = ab99Pointer(uint8(1))
	ab99Control.Ab99Key.Ab99RequestedInstance = ab99Pointer(uint8(2))
	ab99Control.Ab99Key.Ab99Register = ab99Pointer(uint16(3))
	ab99Control.Ab99Key.Ab99CommandRole = ab99Pointer("read")
	ab99Data, ab99Error = ab99json.Marshal([]Ab99Entry{ab99Control})
	if ab99Error != nil {
		ab99T.Fatal(ab99Error)
	}
	ab99Changes := map[string]func(*Ab99Key){
		"ab99_group":           func(ab99Key *Ab99Key) { ab99Key.Ab99Group = ab99Pointer(uint8(0)) },
		"ab99_group_absent":    func(ab99Key *Ab99Key) { ab99Key.Ab99Group = nil },
		"ab99_instance":        func(ab99Key *Ab99Key) { ab99Key.Ab99RequestedInstance = ab99Pointer(uint8(0)) },
		"ab99_instance_absent": func(ab99Key *Ab99Key) { ab99Key.Ab99RequestedInstance = nil },
		"ab99_register":        func(ab99Key *Ab99Key) { ab99Key.Ab99Register = ab99Pointer(uint16(0)) },
		"ab99_register_absent": func(ab99Key *Ab99Key) { ab99Key.Ab99Register = nil },
		"ab99_role":            func(ab99Key *Ab99Key) { ab99Key.Ab99CommandRole = ab99Pointer("write") },
		"ab99_role_absent":     func(ab99Key *Ab99Key) { ab99Key.Ab99CommandRole = nil },
	}
	for ab99Name, ab99Change := range ab99Changes {
		ab99T.Run(ab99Name, func(ab99T *ab99testing.T) {
			ab99Get(ab99T, ab99Control.Ab99Key)
			ab99Key := ab99Control.Ab99Key
			ab99Change(&ab99Key)
			if _, ab99OK, ab99Error := Ab99Lookup(ab99Key); ab99OK || ab99Error != nil {
				ab99T.Fatal("B524 axis ignored")
			}
		})
	}
}
func TestAb99OP00(ab99T *ab99testing.T) {
	for ab99Index, ab99Field := range []string{"relay_groups_134", "relay_groups_256", "S1", "S2", "S3", "S4", "S5", "S6"} {
		ab99Entry := ab99Get(ab99T, ab99Key("configuration", 0, "request", ab99Field))
		ab99Want := map[uint8]string{0: "kindA", 1: "kindB", 2: "digital"}
		if ab99Index < 2 {
			ab99Want = map[uint8]string{0: "output", 1: "mixer"}
		}
		if *ab99Entry.Ab99Position != uint8(ab99Index+3) || *ab99Entry.Ab99Default != 0 || *ab99Entry.Ab99RequestExtent != 11 || *ab99Entry.Ab99PositiveExtent != 4 || ab99Entry.Ab99WidthBits != nil || !ab99reflect.DeepEqual(ab99Entry.Ab99Values, ab99Want) {
			ab99T.Fatal("OP00 mapping", ab99Field)
		}
		if ab99Index >= 2 && ab99Entry.Ab99OtherValues != "undefined-or-ignored" {
			ab99T.Fatal("undefined enum")
		}
	}
}
func TestAb99OP04(ab99T *ab99testing.T) {
	for ab99Index, ab99Field := range []string{"R1/R7/R8", "R2/R9/R10", "R3/R11/R12"} {
		ab99Entry := ab99Get(ab99T, ab99Key("configuration", 4, "request", ab99Field))
		if *ab99Entry.Ab99Position != uint8(ab99Index+3) || !ab99reflect.DeepEqual(ab99Entry.Ab99Values, map[uint8]string{0: "output", 1: "mixer"}) || *ab99Entry.Ab99Default != 0 {
			ab99T.Fatal("OP04 relays")
		}
	}
	for ab99Index, ab99Field := range []string{"S1", "S2", "S3", "S4", "S5", "S6", "S7", "S8", "S9", "S10", "S11", "S12", "S13"} {
		ab99Entry := ab99Get(ab99T, ab99Key("configuration", 4, "request", ab99Field))
		if *ab99Entry.Ab99Position != []uint8{6, 6, 7, 7, 8, 8, 9, 9, 10, 10, 11, 12, 12}[ab99Index] || *ab99Entry.Ab99RequestExtent != 13 || *ab99Entry.Ab99PositiveExtent != 4 || ab99Entry.Ab99WidthBits != nil || ab99Entry.Ab99OneHot == nil || *ab99Entry.Ab99OneHot || !ab99reflect.DeepEqual(ab99Entry.Ab99Alternatives, []string{"independent sensor flags", "competing nibble-code interpretation unresolved", "repeated sensor associations unresolved"}) {
			ab99T.Fatal("OP04 alternatives", ab99Field)
		}
		ab99Offset := []uint8{0, 4, 0, 4, 0, 4, 0, 4, 0, 4, 0, 0, 4}[ab99Index]
		ab99Count := 4
		if ab99Index == 10 {
			ab99Count = 8
		}
		if len(ab99Entry.Ab99Bits) != ab99Count {
			ab99T.Fatal("bit count")
		}
		for ab99BitIndex, ab99Bit := range ab99Entry.Ab99Bits {
			ab99Label := ""
			if ab99BitIndex < 3 {
				ab99Label = []string{"kindA", "kindB", "digital"}[ab99BitIndex]
			} else if (ab99Index == 8 || ab99Index == 9) && ab99BitIndex == 3 {
				ab99Label = "PWM"
			}
			if ab99Bit.Ab99Index != ab99Offset+uint8(ab99BitIndex) || ab99Bit.Ab99Reserved != (ab99Label == "") || (ab99Label == "" && ab99Bit.Ab99Label != nil) || (ab99Label != "" && (ab99Bit.Ab99Label == nil || *ab99Bit.Ab99Label != ab99Label)) {
				ab99T.Fatal("sensor bit mapping", ab99Field, ab99BitIndex)
			}
			ab99Default := uint8(0)
			if ab99BitIndex == 0 {
				ab99Default = 1
			}
			if ab99Bit.Ab99Default != ab99Default {
				ab99T.Fatal("sensor defaults")
			}
			for ab99Raw := 0; ab99Raw < 256; ab99Raw++ {
				if ((uint8(ab99Raw) >> ab99Bit.Ab99Index) & 1) != ((uint8(ab99Raw) >> uint8(int(ab99Offset)+ab99BitIndex)) & 1) {
					ab99T.Fatal("independent bit extraction")
				}
			}
		}
		if uint8(3<<ab99Offset)&(1<<ab99Entry.Ab99Bits[0].Ab99Index) == 0 || uint8(3<<ab99Offset)&(1<<ab99Entry.Ab99Bits[1].Ab99Index) == 0 {
			ab99T.Fatal("two flags preserved")
		}
	}
}
func TestAb99Mixer(ab99T *ab99testing.T) {
	for _, ab99Variant := range []string{"mixer-smaller", "mixer-larger"} {
		ab99Circuit := ab99Get(ab99T, ab99Key(ab99Variant, 2, "request", "circuit"))
		ab99Want := map[uint8]string{0: "first group", 1: "second group"}
		if ab99Variant == "mixer-larger" {
			ab99Want[2] = "third group"
		}
		if *ab99Circuit.Ab99Position != 3 || !ab99reflect.DeepEqual(ab99Circuit.Ab99Values, ab99Want) || ab99Circuit.Ab99OtherValues != "ignored" {
			ab99T.Fatal("circuit eligibility")
		}
		ab99Mode := ab99Get(ab99T, ab99Key(ab99Variant, 2, "request", "mode"))
		if *ab99Mode.Ab99Position != 4 || !ab99reflect.DeepEqual(ab99Mode.Ab99Values, map[uint8]string{0: "stop/pump-off/close", 1: "heating-regulation", 2: "cooling-regulation", 3: "pump-off/open", 4: "pump-on/close", 5: "pump-on/open"}) || ab99Mode.Ab99OtherValues != "ignored" {
			ab99T.Fatal("modes")
		}
		ab99Setpoint := ab99Get(ab99T, ab99Key(ab99Variant, 2, "request", "setpoint"))
		if *ab99Setpoint.Ab99Position != 5 || *ab99Setpoint.Ab99WidthBits != 8 || *ab99Setpoint.Ab99Default != 0 || *ab99Setpoint.Ab99Unit != "Celsius" || !ab99reflect.DeepEqual(ab99Setpoint.Ab99RelevantModes, []uint8{1, 2}) || !ab99reflect.DeepEqual(ab99Setpoint.Ab99IrrelevantModes, []uint8{0, 3, 4, 5}) || len(ab99Setpoint.Ab99Candidates) != 2 {
			ab99T.Fatal("setpoint relevance")
		}
		ab99Signed := ab99Setpoint.Ab99Candidates[0]
		ab99Other := ab99Setpoint.Ab99Candidates[1]
		ab99Raw := uint8(0xbe)
		if ab99Signed.Ab99Signed == nil || !*ab99Signed.Ab99Signed || *ab99Signed.Ab99Scale != 0.5 || *ab99Signed.Ab99Minimum != -64 || *ab99Signed.Ab99Maximum != 63.5 || float64(int8(ab99Raw))**ab99Signed.Ab99Scale != -33 || *ab99Other.Ab99Raw != 190 || *ab99Other.Ab99Value != 95 || *ab99Other.Ab99Clamp != 95 || ab99Other.Ab99Scale != nil || ab99Other.Ab99Signed != nil {
			ab99T.Fatal("competing BE interpretations")
		}
		ab99Feedback := ab99Get(ab99T, ab99Key(ab99Variant, 2, "positive", "feedback"))
		if *ab99Feedback.Ab99Position != 3 || ab99Feedback.Ab99WidthBits != nil || ab99Feedback.Ab99Values != nil || ab99Feedback.Ab99Candidates != nil {
			ab99T.Fatal("raw feedback")
		}
		ab99Movement := ab99Get(ab99T, ab99Key(ab99Variant, 2, "positive", "movement"))
		if ab99Movement.Ab99Signed == nil || !*ab99Movement.Ab99Signed || *ab99Movement.Ab99Position != 4 || *ab99Movement.Ab99WidthBits != 8 || *ab99Movement.Ab99Unit != "percent" || ab99Movement.Ab99OtherValues != "unresolved" || !ab99reflect.DeepEqual(ab99Movement.Ab99Intervals, []Ab99Interval{{-100, -1, "closing"}, {0, 0, "stopped"}, {1, 100, "opening"}}) {
			ab99T.Fatal("movement ranges")
		}
		for ab99Index, ab99Raw := range []uint8{0x9c, 0x64, 0x9b, 0x65, 0} {
			ab99Value := int(int8(ab99Raw))
			ab99Label := "unresolved"
			for _, ab99Interval := range ab99Movement.Ab99Intervals {
				if ab99Value >= ab99Interval.Ab99Minimum && ab99Value <= ab99Interval.Ab99Maximum {
					ab99Label = ab99Interval.Ab99Label
				}
			}
			if ab99Label != []string{"closing", "opening", "unresolved", "unresolved", "stopped"}[ab99Index] {
				ab99T.Fatal("movement polarity/bounds")
			}
		}
	}
}
func ab99Mutate(ab99Value ab99reflect.Value) {
	switch ab99Value.Kind() {
	case ab99reflect.Pointer:
		if !ab99Value.IsNil() {
			ab99Mutate(ab99Value.Elem())
		}
	case ab99reflect.Struct:
		for ab99Index := 0; ab99Index < ab99Value.NumField(); ab99Index++ {
			ab99Mutate(ab99Value.Field(ab99Index))
		}
	case ab99reflect.Slice:
		for ab99Index := 0; ab99Index < ab99Value.Len(); ab99Index++ {
			ab99Mutate(ab99Value.Index(ab99Index))
		}
	case ab99reflect.Map:
		for _, ab99Key := range ab99Value.MapKeys() {
			ab99Value.SetMapIndex(ab99Key, ab99reflect.ValueOf("ab99_changed"))
		}
	case ab99reflect.String:
		ab99Value.SetString("ab99_changed")
	case ab99reflect.Bool:
		ab99Value.SetBool(!ab99Value.Bool())
	case ab99reflect.Uint8, ab99reflect.Uint16:
		ab99Value.SetUint(99)
	case ab99reflect.Int:
		ab99Value.SetInt(99)
	case ab99reflect.Float64:
		ab99Value.SetFloat(99)
	}
}
func TestAb99Copies(ab99T *ab99testing.T) {
	ab99Before, ab99Error := Ab99Entries()
	if ab99Error != nil {
		ab99T.Fatal(ab99Error)
	}
	ab99Returned, ab99Error := Ab99Entries()
	if ab99Error != nil {
		ab99T.Fatal(ab99Error)
	}
	ab99Mutate(ab99reflect.ValueOf(&ab99Returned).Elem())
	ab99After, ab99Error := Ab99Entries()
	if ab99Error != nil || !ab99reflect.DeepEqual(ab99Before, ab99After) {
		ab99T.Fatal("accessor storage shared")
	}
	for _, ab99Entry := range ab99Before {
		ab99Got := ab99Get(ab99T, ab99Entry.Ab99Key)
		ab99Mutate(ab99reflect.ValueOf(&ab99Got).Elem())
		if !ab99reflect.DeepEqual(ab99Get(ab99T, ab99Entry.Ab99Key), ab99Entry) {
			ab99T.Fatal("lookup storage shared")
		}
	}
}
func TestAb99Malformed(ab99T *ab99testing.T) {
	ab99Saved := ab99Data
	ab99T.Cleanup(func() { ab99Data = ab99Saved })
	ab99Entries, ab99Error := Ab99Entries()
	if ab99Error != nil {
		ab99T.Fatal(ab99Error)
	}
	ab99Duplicate, ab99Error := ab99json.Marshal([]Ab99Entry{ab99Entries[0], ab99Entries[0]})
	if ab99Error != nil {
		ab99T.Fatal(ab99Error)
	}
	for _, ab99Invalid := range [][]byte{[]byte("{"), []byte("null"), []byte("[]"), []byte("[{}]"), append(append([]byte{}, ab99Saved...), []byte(" true")...), ab99Duplicate, []byte(`[{"ab99_unknown":1}]`)} {
		ab99Data = ab99Invalid
		if ab99Got, ab99Error := Ab99Entries(); ab99Error == nil || ab99Got != nil {
			ab99T.Fatal("invalid catalog exposed")
		}
		if ab99Got, ab99OK, ab99Error := Ab99Lookup(ab99Key("missing", 0, "request", "missing")); ab99Error == nil || ab99OK || !ab99reflect.DeepEqual(ab99Got, Ab99Entry{}) {
			ab99T.Fatal("invalid catalog hidden as miss")
		}
	}
}

func TestAb99AbsentZero(ab99T *ab99testing.T) {
	ab99Saved := ab99Data
	ab99T.Cleanup(func() { ab99Data = ab99Saved })
	ab99Entries, ab99Error := Ab99Entries()
	if ab99Error != nil {
		ab99T.Fatal(ab99Error)
	}
	ab99Control := ab99Entries[0]
	ab99Control.Ab99Key.Ab99Group = ab99Pointer(uint8(0))
	ab99Control.Ab99Key.Ab99RequestedInstance = ab99Pointer(uint8(0))
	ab99Control.Ab99Key.Ab99Register = ab99Pointer(uint16(0))
	ab99Control.Ab99Key.Ab99CommandRole = ab99Pointer("")
	ab99Data, ab99Error = ab99json.Marshal([]Ab99Entry{ab99Control})
	if ab99Error != nil {
		ab99T.Fatal(ab99Error)
	}
	ab99Changes := map[string]func(*Ab99Key){
		"ab99_operation": func(ab99Key *Ab99Key) { ab99Key.Ab99Operation = nil },
		"ab99_group":     func(ab99Key *Ab99Key) { ab99Key.Ab99Group = nil },
		"ab99_instance":  func(ab99Key *Ab99Key) { ab99Key.Ab99RequestedInstance = nil },
		"ab99_register":  func(ab99Key *Ab99Key) { ab99Key.Ab99Register = nil },
		"ab99_role":      func(ab99Key *Ab99Key) { ab99Key.Ab99CommandRole = nil },
	}
	for ab99Name, ab99Change := range ab99Changes {
		ab99T.Run(ab99Name, func(ab99T *ab99testing.T) {
			ab99Get(ab99T, ab99Control.Ab99Key)
			ab99Key := ab99Control.Ab99Key
			ab99Change(&ab99Key)
			if _, ab99OK, ab99Error := Ab99Lookup(ab99Key); ab99OK || ab99Error != nil {
				ab99T.Fatal("absent axis conflated with explicit zero")
			}
		})
	}
}

func TestAb99RequiredCoordinates(ab99T *ab99testing.T) {
	ab99Saved := ab99Data
	ab99T.Cleanup(func() { ab99Data = ab99Saved })
	for _, ab99Axis := range []string{"primary", "secondary", "operation"} {
		for _, ab99Form := range []string{"missing", "null", "wrong"} {
			if ab99Axis == "operation" && ab99Form == "wrong" {
				continue
			}
			ab99T.Run(ab99Axis+"/"+ab99Form, func(ab99T *ab99testing.T) {
				var ab99Entries []map[string]any
				if ab99Error := ab99json.Unmarshal(ab99Saved, &ab99Entries); ab99Error != nil {
					ab99T.Fatal(ab99Error)
				}
				ab99Key := ab99Entries[0]["key"].(map[string]any)
				switch ab99Form {
				case "missing":
					delete(ab99Key, ab99Axis)
				case "null":
					ab99Key[ab99Axis] = nil
				case "wrong":
					ab99Key[ab99Axis] = 0
				}
				var ab99Error error
				ab99Data, ab99Error = ab99json.Marshal(ab99Entries)
				if ab99Error != nil {
					ab99T.Fatal(ab99Error)
				}
				if ab99Got, ab99Error := Ab99Entries(); ab99Error == nil || ab99Got != nil {
					ab99T.Fatal("invalid coordinates exposed")
				}
				if ab99Got, ab99OK, ab99Error := Ab99Lookup(Ab99Key{}); ab99Error == nil || ab99OK || !ab99reflect.DeepEqual(ab99Got, Ab99Entry{}) {
					ab99T.Fatal("invalid coordinates hidden as miss")
				}
			})
		}
	}
}
