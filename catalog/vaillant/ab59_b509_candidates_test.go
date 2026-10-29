package vaillant

import (
	ab59bytes "bytes"
	ab59binary "encoding/binary"
	ab59fmt "fmt"
	ab59reflect "reflect"
	ab59strings "strings"
	ab59testing "testing"
)

func ab59Pointer[ab59T any](ab59Value ab59T) *ab59T { return &ab59Value }

func ab59Key(ab59Operation byte, ab59Variant, ab59Field string, ab59Direction *string) Ab59Key {
	return Ab59Key{"vaillant", ab59Variant, 0xb5, ab59Pointer(byte(9)), ab59Pointer(ab59Operation), ab59Direction, ab59Field}
}

func ab59Catalog(ab59T *ab59testing.T) Ab59Catalog {
	ab59T.Helper()
	ab59Catalog, ab59Error := Ab59B509Candidates()
	if ab59Error != nil {
		ab59T.Fatal(ab59Error)
	}
	return ab59Catalog
}

func ab59Lookup(ab59T *ab59testing.T, ab59Key Ab59Key) Ab59Entry {
	ab59T.Helper()
	ab59Entry, ab59Found, ab59Error := Ab59Lookup(ab59Key)
	if ab59Error != nil || !ab59Found {
		ab59T.Fatalf("ab59 lookup %#v: found=%v error=%v", ab59Key, ab59Found, ab59Error)
	}
	return ab59Entry
}

func TestAb59FiniteContract(ab59T *ab59testing.T) {
	ab59Catalog := ab59Catalog(ab59T)
	if ab59Catalog.Ab59MetadataPresence != Ab05Proven || ab59Catalog.Ab59WireAdaptation != Ab05Unknown || ab59Catalog.Ab59InstalledSupport != Ab05Unknown || ab59Catalog.Ab59Invocation != Ab05Unknown {
		ab59T.Fatal("ab59 evidence boundaries changed")
	}
	ab59Limitations := []string{
		"Finite candidate declarations only; the broader family catalog remains incomplete.",
		"Candidate semantics are Hypothesis; metadata presence does not establish routable invocation or physical support.",
		"Unprovided coordinates, directions and widths remain null; null keys never wildcard-match.",
		"Raw carriers must be preserved; these declarations do not adapt wire data or select alternatives.",
		"Generic OP0D/0E access and identity scanning are outside these entries.",
	}
	if !ab59reflect.DeepEqual(ab59Catalog.Ab59Limitations, ab59Limitations) {
		ab59T.Fatal("ab59 limitations changed")
	}
	ab59Expected := []Ab59Entry{}
	ab59Add := func(ab59Operation byte, ab59Variant, ab59Field string, ab59Direction *string, ab59Width *int, ab59Signedness, ab59Order, ab59Conversion string, ab59Divisor *int, ab59Unit *string, ab59Default, ab59Lower, ab59Upper *int64, ab59Mappings map[string]string, ab59Other *string, ab59State string) {
		ab59Expected = append(ab59Expected, Ab59Entry{
			Ab59Key:       ab59Key(ab59Operation, ab59Variant, ab59Field, ab59Direction),
			Ab59WidthBits: ab59Width, Ab59Signedness: ab59Signedness, Ab59ByteOrder: ab59Order,
			Ab59Conversion: ab59Conversion, Ab59Divisor: ab59Divisor, Ab59Unit: ab59Unit,
			Ab59Default: ab59Default, Ab59Lower: ab59Lower, Ab59Upper: ab59Upper,
			Ab59Mappings: ab59Mappings, Ab59OtherFeedback: ab59Other, Ab59State: ab59State, Ab59Grade: Ab05Hypothesis,
		})
	}
	ab59Add(7, "selector", "selector", ab59Pointer("request"), ab59Pointer(8), "unsigned", "not_applicable", "raw", nil, nil, nil, nil, nil, map[string]string{}, nil, "raw")
	ab59Add(7, "positive-u16", "value", ab59Pointer("response"), ab59Pointer(16), "unsigned", "little_endian", "raw", nil, nil, ab59Pointer(int64(0)), nil, nil, map[string]string{}, nil, "positive-candidate")
	ab59Add(7, "positive-u32", "value", ab59Pointer("response"), ab59Pointer(32), "unsigned", "little_endian", "raw", nil, nil, ab59Pointer(int64(0)), nil, nil, map[string]string{}, nil, "positive-candidate")
	ab59Add(32, "status-query", "address", ab59Pointer("request"), nil, "Unknown", "Unknown", "Unknown", nil, nil, nil, nil, nil, map[string]string{}, nil, "raw")
	ab59Add(33, "raw-quantities", "address", nil, nil, "Unknown", "Unknown", "Unknown", nil, nil, nil, nil, nil, map[string]string{}, nil, "raw")
	for _, ab59Field := range []string{"minimum", "maximum", "default", "step", "factor"} {
		ab59Add(33, "raw-quantities", ab59Field, nil, ab59Pointer(16), "signed", "little_endian", "raw", nil, nil, nil, nil, nil, map[string]string{}, nil, "raw")
	}
	ab59Add(39, "raw-u16", "value", nil, ab59Pointer(16), "unsigned", "Unknown", "raw", nil, nil, nil, nil, nil, map[string]string{}, nil, "raw")
	ab59Add(41, "raw-0f-f0", "value", nil, ab59Pointer(8), "unsigned", "not_applicable", "raw", nil, nil, nil, nil, nil, map[string]string{"15": "raw 0F", "240": "raw F0"}, nil, "raw")
	ab59Add(44, "default-domain-conflict", "value", nil, ab59Pointer(8), "unsigned", "not_applicable", "Unknown", nil, ab59Pointer("minutes"), ab59Pointer(int64(0)), ab59Pointer(int64(1)), ab59Pointer(int64(255)), map[string]string{}, nil, "raw")
	ab59Add(59, "quantities", "temperature", nil, ab59Pointer(16), "signed", "little_endian", "divide", ab59Pointer(16), ab59Pointer("degrees Celsius"), nil, nil, nil, map[string]string{}, nil, "raw")
	ab59Add(59, "quantities", "flow", nil, ab59Pointer(16), "unsigned", "little_endian", "divide", ab59Pointer(100), ab59Pointer("litres per minute"), nil, nil, nil, map[string]string{}, nil, "raw")
	ab59Add(59, "quantities", "volume", nil, ab59Pointer(16), "unsigned", "little_endian", "divide", ab59Pointer(10), ab59Pointer("litres"), nil, nil, nil, map[string]string{}, nil, "raw")
	for _, ab59Field := range []string{"minimum_volume", "maximum_volume"} {
		ab59Add(59, "quantities", ab59Field, nil, nil, "Unknown", "Unknown", "raw", nil, nil, nil, nil, nil, map[string]string{}, nil, "raw")
	}
	ab59Add(6, "prefix-only-positive", "prefix", ab59Pointer("response"), nil, "Unknown", "Unknown", "Unknown", nil, nil, nil, nil, nil, map[string]string{}, nil, "prefix-only-positive")
	ab59Add(6, "positive-with-feedback", "feedback", ab59Pointer("response"), ab59Pointer(8), "unsigned", "not_applicable", "raw", nil, nil, nil, nil, nil, map[string]string{"1": "accepted"}, ab59Pointer("rejected"), "explicit-feedback")
	ab59Add(6, "negative-without-extra-discriminator", "prefix", ab59Pointer("response"), nil, "Unknown", "Unknown", "Unknown", nil, nil, nil, nil, nil, map[string]string{}, nil, "ambiguous-negative")
	if len(ab59Catalog.Ab59Entries) != 21 {
		ab59T.Fatalf("ab59 expected 21 finite entries, got %d", len(ab59Catalog.Ab59Entries))
	}
	for ab59Index, ab59Want := range ab59Expected {
		ab59T.Run(ab59fmt.Sprintf("ab59_row_%d", ab59Index), func(ab59T *ab59testing.T) {
			ab59Got := ab59Catalog.Ab59Entries[ab59Index]
			if !ab59reflect.DeepEqual(ab59Got, ab59Lookup(ab59T, ab59Want.Ab59Key)) {
				ab59T.Fatal("ab59 accessor and lookup differ")
			}
			ab59Got.Ab59Notes = nil
			if !ab59reflect.DeepEqual(ab59Got, ab59Want) {
				ab59T.Fatalf("ab59 metadata differs: got %#v want %#v", ab59Got, ab59Want)
			}
		})
	}
	ab59Notes := map[int][]string{
		0: {"cannot choose the reply width"}, 1: {"incomplete/raw"}, 2: {"incomplete/raw"},
		3: {"does not toggle", "Address width is Unknown"}, 4: {"Address width is Unknown"},
		5: {"enforcement are Unknown", "not a converter"}, 6: {"enforcement are Unknown", "not a converter"},
		7: {"enforcement are Unknown", "not a converter"}, 8: {"enforcement are Unknown", "not a converter"}, 9: {"enforcement are Unknown", "not a converter"},
		10: {"Complement algorithm is Unknown"}, 11: {"Other values have no established labels"},
		12: {"conflicts with domain 1..255", "Hypothesis"}, 16: {"without division"}, 17: {"without division"},
		18: {"does not select an alternative", "successful reset"}, 19: {"only in this explicitly selected alternative", "does not establish successful reset"}, 20: {"No extra discriminator", "does not establish physical NACK"},
	}
	for ab59Index, ab59Fragments := range ab59Notes {
		for _, ab59Fragment := range ab59Fragments {
			if !ab59strings.Contains(ab59strings.Join(ab59Catalog.Ab59Entries[ab59Index].Ab59Notes, " "), ab59Fragment) {
				ab59T.Fatalf("ab59 row %d missing caveat %q", ab59Index, ab59Fragment)
			}
		}
	}
}

func TestAb59ExactKeyAxes(ab59T *ab59testing.T) {
	for ab59Index, ab59Entry := range ab59Catalog(ab59T).Ab59Entries {
		ab59Control := ab59Entry.Ab59Key
		ab59Changes := map[string]func(*Ab59Key){
			"ab59_family":           func(ab59Key *Ab59Key) { ab59Key.Ab59Family = "ab59_other" },
			"ab59_variant":          func(ab59Key *Ab59Key) { ab59Key.Ab59Variant = "ab59_other" },
			"ab59_primary":          func(ab59Key *Ab59Key) { ab59Key.Ab59Primary = 0 },
			"ab59_secondary_absent": func(ab59Key *Ab59Key) { ab59Key.Ab59Secondary = nil },
			"ab59_secondary_zero":   func(ab59Key *Ab59Key) { ab59Key.Ab59Secondary = ab59Pointer(byte(0)) },
			"ab59_operation_absent": func(ab59Key *Ab59Key) { ab59Key.Ab59Operation = nil },
			"ab59_operation_zero":   func(ab59Key *Ab59Key) { ab59Key.Ab59Operation = ab59Pointer(byte(0)) },
			"ab59_direction":        func(ab59Key *Ab59Key) { ab59Key.Ab59Direction = ab59Pointer("ab59_other") },
			"ab59_direction_empty":  func(ab59Key *Ab59Key) { ab59Key.Ab59Direction = ab59Pointer("") },
			"ab59_field":            func(ab59Key *Ab59Key) { ab59Key.Ab59Field = "ab59_other" },
		}
		if ab59Control.Ab59Direction != nil {
			ab59Changes["ab59_direction_absent"] = func(ab59Key *Ab59Key) { ab59Key.Ab59Direction = nil }
		} else {
			ab59Changes["ab59_direction_present"] = func(ab59Key *Ab59Key) { ab59Key.Ab59Direction = ab59Pointer("response") }
		}
		for ab59Name, ab59Change := range ab59Changes {
			ab59T.Run(ab59fmt.Sprintf("ab59_row_%d/%s", ab59Index, ab59Name), func(ab59T *ab59testing.T) {
				if !ab59reflect.DeepEqual(ab59Lookup(ab59T, ab59Control), ab59Entry) {
					ab59T.Fatal("ab59 invalid positive control")
				}
				ab59Changed := ab59Control
				ab59Change(&ab59Changed)
				ab59Got, ab59Found, ab59Error := Ab59Lookup(ab59Changed)
				if ab59Error != nil || ab59Found || !ab59reflect.DeepEqual(ab59Got, Ab59Entry{}) {
					ab59T.Fatalf("ab59 changed axis matched: %v %v", ab59Found, ab59Error)
				}
				if !ab59reflect.DeepEqual(ab59Lookup(ab59T, ab59Control), ab59Entry) {
					ab59T.Fatal("ab59 control damaged")
				}
			})
		}
	}
}

func ab59Poison(ab59Entry *Ab59Entry) {
	ab59Entry.Ab59Key.Ab59Family = "ab59_mutated"
	if ab59Entry.Ab59Key.Ab59Secondary != nil {
		*ab59Entry.Ab59Key.Ab59Secondary = 0
	}
	if ab59Entry.Ab59Key.Ab59Operation != nil {
		*ab59Entry.Ab59Key.Ab59Operation = 0
	}
	if ab59Entry.Ab59Key.Ab59Direction != nil {
		*ab59Entry.Ab59Key.Ab59Direction = "ab59_mutated"
	}
	for _, ab59Pointer := range []*int{ab59Entry.Ab59Position, ab59Entry.Ab59WidthBits, ab59Entry.Ab59Divisor} {
		if ab59Pointer != nil {
			*ab59Pointer = -999
		}
	}
	for _, ab59Pointer := range []*int64{ab59Entry.Ab59Default, ab59Entry.Ab59Lower, ab59Entry.Ab59Upper} {
		if ab59Pointer != nil {
			*ab59Pointer = -999
		}
	}
	for _, ab59Pointer := range []*string{ab59Entry.Ab59Unit, ab59Entry.Ab59OtherFeedback} {
		if ab59Pointer != nil {
			*ab59Pointer = "ab59_mutated"
		}
	}
	ab59Entry.Ab59Mappings["ab59_mutated"] = "ab59_mutated"
	for ab59Index := range ab59Entry.Ab59Notes {
		ab59Entry.Ab59Notes[ab59Index] = "ab59_mutated"
	}
}

func TestAb59DeepCopy(ab59T *ab59testing.T) {
	ab59Before := ab59Catalog(ab59T)
	ab59Returned := ab59Catalog(ab59T)
	ab59Returned.Ab59Limitations[0] = "ab59_mutated"
	for ab59Index := range ab59Returned.Ab59Entries {
		ab59Poison(&ab59Returned.Ab59Entries[ab59Index])
	}
	if !ab59reflect.DeepEqual(ab59Catalog(ab59T), ab59Before) {
		ab59T.Fatal("ab59 nested accessor storage leaked")
	}
	for _, ab59Entry := range ab59Before.Ab59Entries {
		ab59ReturnedEntry := ab59Lookup(ab59T, ab59Entry.Ab59Key)
		ab59Poison(&ab59ReturnedEntry)
		if !ab59reflect.DeepEqual(ab59Lookup(ab59T, ab59Entry.Ab59Key), ab59Entry) {
			ab59T.Fatal("ab59 nested lookup storage leaked")
		}
	}
	if !ab59reflect.DeepEqual(ab59Catalog(ab59T), ab59Before) {
		ab59T.Fatal("ab59 lookup contaminated accessor")
	}
}

func ab59Interpret(ab59Entry Ab59Entry, ab59Raw []byte) (float64, bool) {
	if ab59Entry.Ab59WidthBits == nil || *ab59Entry.Ab59WidthBits != len(ab59Raw)*8 || ab59Entry.Ab59ByteOrder != "little_endian" {
		return 0, false
	}
	var ab59Value float64
	switch *ab59Entry.Ab59WidthBits {
	case 16:
		ab59Word := ab59binary.LittleEndian.Uint16(ab59Raw)
		switch ab59Entry.Ab59Signedness {
		case "signed":
			ab59Value = float64(int16(ab59Word))
		case "unsigned":
			ab59Value = float64(ab59Word)
		default:
			return 0, false
		}
	case 32:
		if ab59Entry.Ab59Signedness != "unsigned" {
			return 0, false
		}
		ab59Value = float64(ab59binary.LittleEndian.Uint32(ab59Raw))
	default:
		return 0, false
	}
	switch ab59Entry.Ab59Conversion {
	case "raw":
		return ab59Value, true
	case "divide":
		if ab59Entry.Ab59Divisor == nil || *ab59Entry.Ab59Divisor <= 0 {
			return 0, false
		}
		return ab59Value / float64(*ab59Entry.Ab59Divisor), true
	default:
		return 0, false
	}
}

func TestAb59SyntheticQuantities(ab59T *ab59testing.T) {
	ab59Cases := []struct {
		Ab59Key      Ab59Key
		Ab59Raw      []byte
		Ab59Want     float64
		Ab59Complete bool
	}{
		{ab59Key(7, "positive-u16", "value", ab59Pointer("response")), []byte{0x34, 0x12}, 4660, true},
		{ab59Key(7, "positive-u32", "value", ab59Pointer("response")), []byte{0x78, 0x56, 0x34, 0x12}, 305419896, true},
		{ab59Key(7, "positive-u16", "value", ab59Pointer("response")), []byte{0x78, 0x56, 0x34, 0x12}, 0, false},
		{ab59Key(7, "positive-u32", "value", ab59Pointer("response")), []byte{0x34, 0x12}, 0, false},
		{ab59Key(59, "quantities", "temperature", nil), []byte{0xf0, 0xff}, -1, true},
		{ab59Key(59, "quantities", "flow", nil), []byte{0x64, 0}, 1, true},
		{ab59Key(59, "quantities", "flow", nil), []byte{0xff, 0xff}, 655.35, true},
		{ab59Key(59, "quantities", "volume", nil), []byte{0x0a, 0}, 1, true},
	}
	for ab59Index, ab59Case := range ab59Cases {
		ab59T.Run(ab59fmt.Sprintf("ab59_sample_%d", ab59Index), func(ab59T *ab59testing.T) {
			ab59Before := ab59bytes.Clone(ab59Case.Ab59Raw)
			ab59Value, ab59Complete := ab59Interpret(ab59Lookup(ab59T, ab59Case.Ab59Key), ab59Case.Ab59Raw)
			if ab59Value != ab59Case.Ab59Want || ab59Complete != ab59Case.Ab59Complete || !ab59bytes.Equal(ab59Case.Ab59Raw, ab59Before) {
				ab59T.Fatalf("ab59 sample got %v %v", ab59Value, ab59Complete)
			}
		})
	}
	for ab59CaseIndex, ab59Words := range [][]int16{{-1, 2, 0, 1, -2}, {2, -1, 0, 1, -2}} {
		for ab59Index, ab59Field := range []string{"minimum", "maximum", "default", "step", "factor"} {
			ab59Raw := make([]byte, 2)
			ab59binary.LittleEndian.PutUint16(ab59Raw, uint16(ab59Words[ab59Index]))
			ab59Entry := ab59Lookup(ab59T, ab59Key(33, "raw-quantities", ab59Field, nil))
			ab59Value, ab59Complete := ab59Interpret(ab59Entry, ab59Raw)
			if !ab59Complete || ab59Value != float64(ab59Words[ab59Index]) || ab59Entry.Ab59Divisor != nil {
				ab59T.Fatalf("ab59 OP33 case %d field %s lost independent raw quantity", ab59CaseIndex, ab59Field)
			}
		}
	}
	ab59Raw39 := []byte{0xff, 0xff}
	ab59Entry39 := ab59Lookup(ab59T, ab59Key(39, "raw-u16", "value", nil))
	if ab59binary.LittleEndian.Uint16(ab59Raw39) != 65535 || ab59Entry39.Ab59Conversion != "raw" || ab59Entry39.Ab59Divisor != nil {
		ab59T.Fatal("ab59 OP39 complemented")
	}
	for _, ab59Field := range []string{"minimum_volume", "maximum_volume"} {
		ab59Entry := ab59Lookup(ab59T, ab59Key(59, "quantities", ab59Field, nil))
		ab59Raw := []byte{0x0a, 0}
		ab59Before := ab59bytes.Clone(ab59Raw)
		_, ab59Complete := ab59Interpret(ab59Entry, ab59Raw)
		if ab59Complete || ab59Entry.Ab59Conversion != "raw" || ab59Entry.Ab59Divisor != nil || !ab59bytes.Equal(ab59Raw, ab59Before) {
			ab59T.Fatal("ab59 raw volume acquired unsupported conversion")
		}
	}
}

func TestAb59FeedbackAlternatives(ab59T *ab59testing.T) {
	ab59Feedback := ab59Lookup(ab59T, ab59Key(6, "positive-with-feedback", "feedback", ab59Pointer("response")))
	for ab59Index, ab59Raw := range []byte{0, 1, 2, 255} {
		ab59Got, ab59Found := ab59Feedback.Ab59Mappings[ab59fmt.Sprint(ab59Raw)]
		if !ab59Found && ab59Feedback.Ab59OtherFeedback != nil {
			ab59Got = *ab59Feedback.Ab59OtherFeedback
		}
		if ab59Got != []string{"rejected", "accepted", "rejected", "rejected"}[ab59Index] {
			ab59T.Fatalf("ab59 feedback %d got %s", ab59Raw, ab59Got)
		}
	}
	for ab59Index, ab59Variant := range []string{"prefix-only-positive", "negative-without-extra-discriminator"} {
		ab59Entry := ab59Lookup(ab59T, ab59Key(6, ab59Variant, "prefix", ab59Pointer("response")))
		if ab59Entry.Ab59State != []string{"prefix-only-positive", "ambiguous-negative"}[ab59Index] || len(ab59Entry.Ab59Mappings) != 0 || ab59Entry.Ab59OtherFeedback != nil || ab59Entry.Ab59WidthBits != nil {
			ab59T.Fatal("ab59 apparent prefix classified as feedback")
		}
	}
}

func TestAb59AbsentAndZeroKeys(ab59T *ab59testing.T) {
	ab59Base := ab59Catalog(ab59T).Ab59Entries[0]
	for _, ab59Axis := range []string{"ab59_secondary", "ab59_operation"} {
		ab59T.Run(ab59Axis, func(ab59T *ab59testing.T) {
			ab59Absent := ab59Base
			ab59Zero := ab59Base
			if ab59Axis == "ab59_secondary" {
				ab59Absent.Ab59Key.Ab59Secondary = nil
				ab59Zero.Ab59Key.Ab59Secondary = ab59Pointer(byte(0))
			} else {
				ab59Absent.Ab59Key.Ab59Operation = nil
				ab59Zero.Ab59Key.Ab59Operation = ab59Pointer(byte(0))
			}
			for _, ab59Control := range []Ab59Entry{ab59Absent, ab59Zero} {
				ab59Got, ab59Found := ab59Find([]Ab59Entry{ab59Control}, ab59Control.Ab59Key)
				if !ab59Found || !ab59reflect.DeepEqual(ab59Got, ab59Control) {
					ab59T.Fatal("ab59 optional-axis positive control failed")
				}
				ab59Changed := ab59Zero.Ab59Key
				if ab59reflect.DeepEqual(ab59Changed, ab59Control.Ab59Key) {
					ab59Changed = ab59Absent.Ab59Key
				}
				if _, ab59Found = ab59Find([]Ab59Entry{ab59Control}, ab59Changed); ab59Found {
					ab59T.Fatal("ab59 absent and zero matched")
				}
			}
		})
	}
}
