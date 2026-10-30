package vaillant

import (
	ab86binary "encoding/binary"
	ab86reflect "reflect"
	ab86testing "testing"
)

func ab86Pointer[ab86Value any](ab86Input ab86Value) *ab86Value { return &ab86Input }

func ab86Catalog(ab86T *ab86testing.T) Ab86HistoryMetadata {
	ab86T.Helper()
	ab86Got, ab86Error := Ab86HistoryCatalog()
	if ab86Error != nil {
		ab86T.Fatal(ab86Error)
	}
	return ab86Got
}

func ab86ExpectedField(ab86Name string, ab86Operation *uint8, ab86Width *int, ab86Signedness, ab86Order, ab86Conversion string) Ab86HistoryField {
	return Ab86HistoryField{
		Ab86Key:       Ab86Key{Ab86Family: "B503", Ab86Primary: 181, Ab86Secondary: ab86Pointer(uint8(3)), Ab86Operation: ab86Operation, Ab86Field: ab86Name},
		Ab86WidthBits: ab86Width, Ab86Signedness: ab86Signedness, Ab86ByteOrder: ab86Order, Ab86Conversion: ab86Conversion,
		Ab86Enumeration: map[int]string{}, Ab86Replacement: []int{}, Ab86Reserved: []int{}, Ab86Bits: map[int]string{},
		Ab86Grade: "Hypothesis", Ab86RawFallback: "preserve original field bytes", Ab86Branches: map[int]string{}, Ab86NoItemTail: []int{},
	}
}

func TestAb86ProductionMetadata(ab86T *ab86testing.T) {
	ab86Got := ab86Catalog(ab86T)
	ab86Expected := Ab86HistoryMetadata{Ab86Representation: "Proven", Ab86Qualification: "Unknown", Ab86Unknown: []string{
		"century", "timezone", "full calendar validity", "copied-span adaptation", "physical ACK", "reset effect", "persistence", "freshness", "execution permission", "unprovided layouts and flag meanings", "broader family completeness",
	}}
	ab86Marker := ab86ExpectedField("marker", ab86Pointer(uint8(1)), ab86Pointer(8), "unsigned", "not_applicable", "branch-before-calendar")
	ab86Marker.Ab86Enumeration = map[int]string{0: "no_item", 1: "code_only", 2: "calendar_components"}
	ab86Marker.Ab86Branches = map[int]string{0: "require-seven-FF-tail", 1: "retain-code-without-calendar", 2: "validate-each-component-before-decimal"}
	ab86Marker.Ab86NoItemTail = []int{255, 255, 255, 255, 255, 255, 255}
	ab86Expected.Ab86Fields = append(ab86Expected.Ab86Fields, ab86Marker, ab86ExpectedField("error_code", ab86Pointer(uint8(1)), nil, "Unknown", "Unknown", "identity"))
	for ab86Index, ab86Name := range []string{"minute", "hour", "day", "month", "year"} {
		ab86Field := ab86ExpectedField(ab86Name, ab86Pointer(uint8(1)), ab86Pointer(8), "unsigned", "not_applicable", "decimal-after-validity")
		ab86Field.Ab86Bounds = [][]int{{0, 59}, {0, 23}, {1, 31}, {1, 12}, {0, 99}}[ab86Index]
		ab86Field.Ab86Replacement = []int{255}
		ab86Expected.Ab86Fields = append(ab86Expected.Ab86Fields, ab86Field)
	}
	ab86Reset := ab86ExpectedField("reset_acceptance", nil, ab86Pointer(8), "unsigned", "not_applicable", "identity")
	ab86Reset.Ab86Default = ab86Pointer(0)
	ab86Reset.Ab86Enumeration = map[int]string{0: "false", 1: "true"}
	ab86Expected.Ab86Fields = append(ab86Expected.Ab86Fields, ab86Reset)
	ab86Flags := ab86ExpectedField("history_flags", ab86Pointer(uint8(7)), ab86Pointer(8), "unsigned", "not_applicable", "bitset")
	ab86Flags.Ab86Bits = map[int]string{0: "history_flag_0", 1: "history_flag_1", 2: "history_flag_2", 3: "history_flag_3", 4: "history_flag_4"}
	ab86Flags.Ab86Reserved = []int{5, 6, 7}
	ab86Expected.Ab86Fields = append(ab86Expected.Ab86Fields, ab86Flags)
	for _, ab86Name := range []string{"current_error_1", "current_error_2", "current_error_3", "current_error_4", "current_error_5"} {
		ab86Expected.Ab86Fields = append(ab86Expected.Ab86Fields, ab86ExpectedField(ab86Name, ab86Pointer(uint8(7)), ab86Pointer(16), "unsigned", "little_endian", "identity"))
	}
	ab86Feedback := ab86ExpectedField("feedback", ab86Pointer(uint8(7)), ab86Pointer(8), "unsigned", "not_applicable", "identity")
	ab86Feedback.Ab86Enumeration = map[int]string{1: "accepted"}
	ab86Feedback.Ab86Otherwise = ab86Pointer("rejected")
	ab86Expected.Ab86Fields = append(ab86Expected.Ab86Fields, ab86Feedback)
	if !ab86reflect.DeepEqual(ab86Got, ab86Expected) {
		ab86T.Fatalf("ab86 production declarations differ: got %#v want %#v", ab86Got, ab86Expected)
	}
	for _, ab86Field := range ab86Got.Ab86Fields {
		ab86LookupField, ab86Found, ab86Error := Ab86HistoryLookup(ab86Field.Ab86Key)
		if ab86Error != nil || !ab86Found || !ab86reflect.DeepEqual(ab86LookupField, ab86Field) {
			ab86T.Fatal("ab86 lookup control", ab86Field.Ab86Key, ab86Error)
		}
	}
}

func ab86Named(ab86T *ab86testing.T, ab86Name string) Ab86HistoryField {
	ab86T.Helper()
	for _, ab86Field := range ab86Catalog(ab86T).Ab86Fields {
		if ab86Field.Ab86Key.Ab86Field == ab86Name {
			return ab86Field
		}
	}
	ab86T.Fatal("ab86 missing field", ab86Name)
	return Ab86HistoryField{}
}

func ab86Calendar(ab86Field Ab86HistoryField, ab86Raw int) (int, bool) {
	for _, ab86Replacement := range ab86Field.Ab86Replacement {
		if ab86Raw == ab86Replacement {
			return 0, false
		}
	}
	if len(ab86Field.Ab86Bounds) != 2 || ab86Raw < ab86Field.Ab86Bounds[0] || ab86Raw > ab86Field.Ab86Bounds[1] {
		return 0, false
	}
	if ab86Field.Ab86Conversion != "decimal-after-validity" {
		return 0, false
	}
	return ab86Raw, true
}

func ab86HistoryCase(ab86T *ab86testing.T, ab86Marker int, ab86Tail []int, ab86Code int, ab86Components []int) (string, int, []int) {
	ab86T.Helper()
	ab86Field := ab86Named(ab86T, "marker")
	switch ab86Field.Ab86Branches[ab86Marker] {
	case "require-seven-FF-tail":
		if ab86reflect.DeepEqual(ab86Tail, ab86Field.Ab86NoItemTail) {
			return ab86Field.Ab86Enumeration[ab86Marker], 0, nil
		}
	case "retain-code-without-calendar":
		if ab86Named(ab86T, "error_code").Ab86Conversion == "identity" {
			return ab86Field.Ab86Enumeration[ab86Marker], ab86Code, nil
		}
	case "validate-each-component-before-decimal":
		if len(ab86Components) != 5 {
			return "Unknown", 0, nil
		}
		ab86Result := []int{}
		for ab86Index, ab86Name := range []string{"minute", "hour", "day", "month", "year"} {
			ab86Value, ab86Valid := ab86Calendar(ab86Named(ab86T, ab86Name), ab86Components[ab86Index])
			if !ab86Valid {
				return "Unknown", ab86Code, nil
			}
			ab86Result = append(ab86Result, ab86Value)
		}
		return ab86Field.Ab86Enumeration[ab86Marker], ab86Code, ab86Result
	}
	return "Unknown", 0, nil
}

func TestAb86SyntheticHistory(ab86T *ab86testing.T) {
	ab86Label, ab86Code, ab86Date := ab86HistoryCase(ab86T, 0, []int{255, 255, 255, 255, 255, 255, 255}, 0, nil)
	if ab86Label != "no_item" || ab86Code != 0 || ab86Date != nil {
		ab86T.Fatal("ab86 no item")
	}
	ab86Label, _, _ = ab86HistoryCase(ab86T, 0, []int{255, 255, 255, 255, 255, 255, 0}, 0, nil)
	if ab86Label != "Unknown" {
		ab86T.Fatal("ab86 no-item tail must match")
	}
	ab86Label, ab86Code, ab86Date = ab86HistoryCase(ab86T, 1, nil, 0x1234, []int{255, 255, 255, 255, 255})
	if ab86Label != "code_only" || ab86Code != 0x1234 || ab86Date != nil {
		ab86T.Fatal("ab86 marker 01 retains code without calendar")
	}
	ab86Label, ab86Code, ab86Date = ab86HistoryCase(ab86T, 2, nil, 0x1234, []int{59, 23, 31, 12, 99})
	if ab86Label != "calendar_components" || ab86Code != 0x1234 || !ab86reflect.DeepEqual(ab86Date, []int{59, 23, 31, 12, 99}) {
		ab86T.Fatal("ab86 component boundaries")
	}
	for ab86Index, ab86Name := range []string{"minute", "hour", "day", "month", "year"} {
		ab86Field := ab86Named(ab86T, ab86Name)
		for _, ab86Raw := range []int{255, []int{60, 24, 32, 13, 100}[ab86Index], []int{-1, -1, 0, 0, -1}[ab86Index]} {
			if _, ab86Valid := ab86Calendar(ab86Field, ab86Raw); ab86Valid {
				ab86T.Fatal("ab86 invalid/replacement converted", ab86Name, ab86Raw)
			}
		}
		ab86Low := []int{0, 0, 1, 1, 0}[ab86Index]
		if ab86Value, ab86Valid := ab86Calendar(ab86Field, ab86Low); !ab86Valid || ab86Value != ab86Low {
			ab86T.Fatal("ab86 lower bound", ab86Name)
		}
	}
	ab86Label, _, ab86Date = ab86HistoryCase(ab86T, 2, nil, 1, []int{255, 255, 255, 255, 255})
	if ab86Label != "Unknown" || ab86Date != nil {
		ab86T.Fatal("ab86 replacement has no calendar claim")
	}
}

func ab86Label(ab86Field Ab86HistoryField, ab86Raw int) string {
	if ab86Label, ab86Found := ab86Field.Ab86Enumeration[ab86Raw]; ab86Found {
		return ab86Label
	}
	if ab86Field.Ab86Otherwise != nil {
		return *ab86Field.Ab86Otherwise
	}
	return "unlabeled"
}

func TestAb86SyntheticFlagsWordsFeedback(ab86T *ab86testing.T) {
	ab86Flags := ab86Named(ab86T, "history_flags")
	for _, ab86Raw := range []int{0x1F, 0xE0} {
		ab86Names := []string{}
		for ab86Bit := 0; ab86Bit < 8; ab86Bit++ {
			if ab86Name, ab86Found := ab86Flags.Ab86Bits[ab86Bit]; ab86Found && ab86Raw&(1<<ab86Bit) != 0 {
				ab86Names = append(ab86Names, ab86Name)
			}
		}
		ab86Expected := []string{}
		if ab86Raw == 0x1F {
			ab86Expected = []string{"history_flag_0", "history_flag_1", "history_flag_2", "history_flag_3", "history_flag_4"}
		}
		if !ab86reflect.DeepEqual(ab86Names, ab86Expected) {
			ab86T.Fatal("ab86 reserved bits included", ab86Names)
		}
	}
	for ab86Index, ab86Name := range []string{"current_error_1", "current_error_2", "current_error_3", "current_error_4", "current_error_5"} {
		ab86Field := ab86Named(ab86T, ab86Name)
		ab86Bytes := [][]byte{{1, 0}, {0, 1}, {0x34, 0x12}, {0, 0}, {255, 255}}[ab86Index]
		var ab86Value uint16
		if ab86Field.Ab86WidthBits != nil && *ab86Field.Ab86WidthBits == 16 && ab86Field.Ab86Conversion == "identity" && ab86Field.Ab86Signedness == "unsigned" {
			switch ab86Field.Ab86ByteOrder {
			case "little_endian":
				ab86Value = ab86binary.LittleEndian.Uint16(ab86Bytes)
			case "big_endian":
				ab86Value = ab86binary.BigEndian.Uint16(ab86Bytes)
			default:
				ab86T.Fatal("ab86 unknown byte order")
			}
		} else {
			ab86T.Fatal("ab86 word metadata")
		}
		if ab86Value != []uint16{1, 256, 0x1234, 0, 65535}[ab86Index] {
			ab86T.Fatal("ab86 independent word", ab86Name, ab86Value)
		}
	}
	ab86Reset := ab86Named(ab86T, "reset_acceptance")
	ab86Feedback := ab86Named(ab86T, "feedback")
	for ab86Raw := 0; ab86Raw <= 255; ab86Raw++ {
		ab86ResetWant := "unlabeled"
		ab86FeedbackWant := "rejected"
		if ab86Raw == 0 {
			ab86ResetWant = "false"
		}
		if ab86Raw == 1 {
			ab86ResetWant = "true"
			ab86FeedbackWant = "accepted"
		}
		if ab86Label(ab86Reset, ab86Raw) != ab86ResetWant || ab86Label(ab86Feedback, ab86Raw) != ab86FeedbackWant {
			ab86T.Fatal("ab86 reset/feedback alternatives", ab86Raw)
		}
	}
}

func ab86CheckAxes(ab86T *ab86testing.T, ab86Field Ab86HistoryField, ab86Production bool) {
	ab86T.Helper()
	ab86Control, ab86Found := ab86Lookup([]Ab86HistoryField{ab86Field}, ab86Field.Ab86Key)
	if !ab86Found || !ab86reflect.DeepEqual(ab86Control, ab86Field) {
		ab86T.Fatal("ab86 invalid axis control")
	}
	ab86KeyType := ab86reflect.TypeOf(ab86Field.Ab86Key)
	for ab86Index := 0; ab86Index < ab86KeyType.NumField(); ab86Index++ {
		ab86T.Run(ab86KeyType.Field(ab86Index).Name, func(ab86T *ab86testing.T) {
			for ab86Case := 0; ab86Case < 3; ab86Case++ {
				ab86Key := ab86Field.Ab86Key
				ab86Axis := ab86reflect.ValueOf(&ab86Key).Elem().Field(ab86Index)
				if ab86Axis.Kind() == ab86reflect.Pointer {
					if ab86Case == 0 {
						ab86Axis.SetZero()
					} else {
						ab86New := ab86reflect.New(ab86Axis.Type().Elem())
						if ab86Case == 2 {
							if ab86New.Elem().Kind() == ab86reflect.String {
								ab86New.Elem().SetString("ab86_other")
							} else {
								ab86New.Elem().SetUint(99)
							}
						}
						ab86Axis.Set(ab86New)
					}
				} else if ab86Axis.Kind() == ab86reflect.String {
					if ab86Case == 0 {
						ab86Axis.SetString("")
					} else {
						ab86Axis.SetString("ab86_other")
					}
				} else {
					ab86Axis.SetUint(uint64(ab86Case))
				}
				if ab86reflect.DeepEqual(ab86Key, ab86Field.Ab86Key) {
					continue
				}
				if _, ab86Found := ab86Lookup([]Ab86HistoryField{ab86Field}, ab86Key); ab86Found {
					ab86T.Fatal("ab86 axis ignored", ab86Case)
				}
				if ab86Production {
					if _, ab86Found, ab86Error := Ab86HistoryLookup(ab86Key); ab86Error != nil || ab86Found {
						ab86T.Fatal("ab86 production axis ignored", ab86Case, ab86Error)
					}
				}
			}
		})
	}
}

func TestAb86ExactIdentityAxes(ab86T *ab86testing.T) {
	for _, ab86Field := range ab86Catalog(ab86T).Ab86Fields {
		ab86T.Run("ab86_"+ab86Field.Ab86Key.Ab86Field, func(ab86T *ab86testing.T) { ab86CheckAxes(ab86T, ab86Field, true) })
	}
	ab86Qualified := Ab86HistoryField{Ab86Key: Ab86Key{
		Ab86Family: "B524", Ab86Variant: ab86Pointer("ab86_synthetic"), Ab86Primary: 181,
		Ab86Secondary: ab86Pointer(uint8(36)), Ab86Operation: ab86Pointer(uint8(1)), Ab86Direction: ab86Pointer("response"), Ab86Field: "ab86_synthetic",
		Ab86Opcode: ab86Pointer(uint8(2)), Ab86Group: ab86Pointer(uint8(3)), Ab86Instance: ab86Pointer(uint8(4)), Ab86Register: ab86Pointer(uint16(5)),
	}}
	ab86T.Run("ab86_qualified_B524", func(ab86T *ab86testing.T) { ab86CheckAxes(ab86T, ab86Qualified, false) })
	ab86Zero := ab86Qualified
	ab86Zero.Ab86Key.Ab86Opcode = ab86Pointer(uint8(0))
	ab86Zero.Ab86Key.Ab86Group = ab86Pointer(uint8(0))
	ab86Zero.Ab86Key.Ab86Instance = ab86Pointer(uint8(0))
	ab86Zero.Ab86Key.Ab86Register = ab86Pointer(uint16(0))
	ab86Zero.Ab86Key.Ab86Operation = ab86Pointer(uint8(0))
	ab86Zero.Ab86Key.Ab86Secondary = ab86Pointer(uint8(0))
	ab86Zero.Ab86Key.Ab86Variant = ab86Pointer("")
	ab86Zero.Ab86Key.Ab86Direction = ab86Pointer("")
	ab86T.Run("ab86_explicit_zero", func(ab86T *ab86testing.T) { ab86CheckAxes(ab86T, ab86Zero, false) })
	if _, ab86Found, ab86Error := Ab86HistoryLookup(ab86Qualified.Ab86Key); ab86Found || ab86Error != nil {
		ab86T.Fatal("ab86 synthetic B524 must not become production metadata")
	}
}

func ab86Corrupt(ab86Value ab86reflect.Value) {
	switch ab86Value.Kind() {
	case ab86reflect.Pointer:
		if !ab86Value.IsNil() {
			ab86Corrupt(ab86Value.Elem())
		}
	case ab86reflect.Struct:
		for ab86Index := 0; ab86Index < ab86Value.NumField(); ab86Index++ {
			ab86Corrupt(ab86Value.Field(ab86Index))
		}
	case ab86reflect.Slice:
		for ab86Index := 0; ab86Index < ab86Value.Len(); ab86Index++ {
			ab86Corrupt(ab86Value.Index(ab86Index))
		}
	case ab86reflect.Map:
		for _, ab86Key := range ab86Value.MapKeys() {
			ab86Value.SetMapIndex(ab86Key, ab86reflect.Zero(ab86Value.Type().Elem()))
		}
	case ab86reflect.String:
		if ab86Value.CanSet() {
			ab86Value.SetString("ab86_mutated")
		}
	case ab86reflect.Int:
		if ab86Value.CanSet() {
			ab86Value.SetInt(42)
		}
	case ab86reflect.Uint8, ab86reflect.Uint16:
		if ab86Value.CanSet() {
			ab86Value.SetUint(42)
		}
	}
}

func TestAb86NestedCopyIsolation(ab86T *ab86testing.T) {
	ab86Before := ab86Catalog(ab86T)
	ab86Changed := ab86Catalog(ab86T)
	ab86Corrupt(ab86reflect.ValueOf(&ab86Changed).Elem())
	if !ab86reflect.DeepEqual(ab86Catalog(ab86T), ab86Before) {
		ab86T.Fatal("ab86 accessor aliases embedded metadata")
	}
	for _, ab86Field := range ab86Before.Ab86Fields {
		ab86Key := ab86Field.Ab86Key
		ab86Changed, ab86Found, ab86Error := Ab86HistoryLookup(ab86Key)
		if !ab86Found || ab86Error != nil {
			ab86T.Fatal("ab86 lookup control", ab86Error)
		}
		ab86Corrupt(ab86reflect.ValueOf(&ab86Changed).Elem())
		ab86Again, ab86Found, ab86Error := Ab86HistoryLookup(ab86Field.Ab86Key)
		if !ab86Found || ab86Error != nil || !ab86reflect.DeepEqual(ab86Again, ab86Field) {
			ab86T.Fatal("ab86 lookup aliases nested data", ab86Error)
		}
	}
}
