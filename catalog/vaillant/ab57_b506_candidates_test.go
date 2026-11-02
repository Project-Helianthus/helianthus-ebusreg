package vaillant

import (
	ab57Bytes "bytes"
	ab57Reflect "reflect"
	ab57Testing "testing"
)

func ab57Byte(ab57Value byte) *byte       { return &ab57Value }
func ab57String(ab57Value string) *string { return &ab57Value }

func ab57Key(ab57Variant string, ab57Operation *byte, ab57Field string) Ab57Key {
	return Ab57Key{Ab57Family: "vaillant", Ab57Variant: ab57Variant, Ab57Primary: 0xB5, Ab57Secondary: ab57Byte(6), Ab57Operation: ab57Operation, Ab57Field: ab57Field}
}

func ab57Entry(ab57T *ab57Testing.T, ab57KeyValue Ab57Key) Ab57Entry {
	ab57T.Helper()
	ab57EntryValue, ab57OK := Ab57B506Lookup(ab57KeyValue)
	if !ab57OK {
		ab57T.Fatalf("missing declaration: %+v", ab57KeyValue)
	}
	return ab57EntryValue
}

func TestAb57ExactKeys(ab57T *ab57Testing.T) {
	ab57Keys := []Ab57Key{
		ab57Key("production-test", ab57Byte(0x11), "parameter"),
		ab57Key("backlight-only", ab57Byte(2), "outputs"),
		ab57Key("three-output", ab57Byte(2), "outputs"),
		ab57Key("eeprom-checksum", nil, "checksum_type"),
	}
	ab57Catalog := Ab57B506Candidates()
	if len(ab57Catalog.Ab57Entries) != len(ab57Keys) {
		ab57T.Fatal("finite entry count")
	}
	for _, ab57Control := range ab57Keys {
		ab57Found := ab57Entry(ab57T, ab57Control)
		if !ab57Reflect.DeepEqual(ab57Found.Ab57Key, ab57Control) {
			ab57T.Fatal("key changed")
		}
		ab57Mutations := map[string]func(*Ab57Key){
			"ab57_family_empty":     func(ab57K *Ab57Key) { ab57K.Ab57Family = "" },
			"ab57_variant_empty":    func(ab57K *Ab57Key) { ab57K.Ab57Variant = "" },
			"ab57_field_empty":      func(ab57K *Ab57Key) { ab57K.Ab57Field = "" },
			"ab57_family":           func(ab57K *Ab57Key) { ab57K.Ab57Family = "ab57_other" },
			"ab57_variant":          func(ab57K *Ab57Key) { ab57K.Ab57Variant = "ab57_other" },
			"ab57_primary":          func(ab57K *Ab57Key) { ab57K.Ab57Primary = 0 },
			"ab57_secondary_absent": func(ab57K *Ab57Key) { ab57K.Ab57Secondary = nil },
			"ab57_secondary_zero":   func(ab57K *Ab57Key) { ab57K.Ab57Secondary = ab57Byte(0) },
			"ab57_operation_zero":   func(ab57K *Ab57Key) { ab57K.Ab57Operation = ab57Byte(0) },
			"ab57_direction":        func(ab57K *Ab57Key) { ab57K.Ab57Direction = ab57String("request") },
			"ab57_direction_empty":  func(ab57K *Ab57Key) { ab57K.Ab57Direction = ab57String("") },
			"ab57_field":            func(ab57K *Ab57Key) { ab57K.Ab57Field = "ab57_other" },
		}
		if ab57Control.Ab57Operation != nil {
			ab57Mutations["ab57_operation_absent"] = func(ab57K *Ab57Key) { ab57K.Ab57Operation = nil }
		}
		for ab57Name, ab57Mutate := range ab57Mutations {
			ab57T.Run("ab57_"+ab57Control.Ab57Variant+"/"+ab57Name, func(ab57T *ab57Testing.T) {
				ab57Changed := ab57Control
				ab57Mutate(&ab57Changed)
				if _, ab57OK := Ab57B506Lookup(ab57Changed); ab57OK {
					ab57T.Fatal("changed axis matched")
				}
				ab57Entry(ab57T, ab57Control)
			})
		}
	}
}

func TestAb57UnionPermutation(ab57T *ab57Testing.T) {
	ab57Production := ab57Entry(ab57T, ab57Key("production-test", ab57Byte(0x11), "parameter"))
	ab57Expected := map[byte]string{0: "opaque-through-end", 1: "opaque-through-end", 2: "opaque-through-end", 4: "opaque-through-end", 5: "opaque-through-end", 3: "raw-byte", 6: "raw-byte", 7: "raw-byte", 8: "raw-byte", 9: "raw-byte", 10: "boolean", 32: "raw-byte", 33: "raw-byte", 34: "raw-byte", 35: "raw-byte", 36: "raw-byte", 37: "raw-byte"}
	if len(ab57Production.Ab57Parameters) != len(ab57Expected) {
		ab57T.Fatal("union size")
	}
	for ab57Rotation := 0; ab57Rotation < len(ab57Expected); ab57Rotation++ {
		for ab57Selector := 0; ab57Selector < 256; ab57Selector++ {
			ab57Parameter, ab57OK := Ab57ParameterLookup(ab57Production, byte(ab57Selector))
			ab57Want, ab57Mapped := ab57Expected[byte(ab57Selector)]
			if ab57OK != ab57Mapped {
				ab57T.Fatalf("selector %d mapping", ab57Selector)
			}
			if !ab57OK {
				continue
			}
			if ab57Parameter.Ab57Selector != byte(ab57Selector) || ab57Parameter.Ab57Representation != ab57Want || ab57Parameter.Ab57Metadata.Ab05Codec != ab57Want || ab57Parameter.Ab57Metadata.Ab05Grade != Ab05Hypothesis || ab57Parameter.Ab57Metadata.Ab05Position != nil || ab57Parameter.Ab57MinimumBytes != nil || ab57Parameter.Ab57MaximumBytes != nil {
				ab57T.Fatalf("selector %d declaration", ab57Selector)
			}
			if ab57Want == "opaque-through-end" {
				if ab57Parameter.Ab57Metadata.Ab05Width != nil {
					ab57T.Fatal("opaque width invented")
				}
			} else if ab57Parameter.Ab57Metadata.Ab05Width == nil || *ab57Parameter.Ab57Metadata.Ab05Width != 1 {
				ab57T.Fatal("one byte width missing")
			}
		}
		ab57Production.Ab57Parameters = append(ab57Production.Ab57Parameters[1:], ab57Production.Ab57Parameters[0])
	}
	ab57Boolean, _ := Ab57ParameterLookup(ab57Production, 10)
	if !ab57Reflect.DeepEqual(ab57Boolean.Ab57Metadata.Ab05Enum, map[int64]string{0: "no", 1: "yes"}) {
		ab57T.Fatal("Boolean meaning must cover exactly 0 and 1")
	}
	for _, ab57Selector := range []byte{0, 1, 2, 4, 5} {
		ab57Opaque, _ := Ab57ParameterLookup(ab57Production, ab57Selector)
		for _, ab57Raw := range [][]byte{{0, 0xAA, 0xFF}, {}} {
			ab57Retained := append([]byte{}, ab57Raw...)
			if ab57Opaque.Ab57Representation != "opaque-through-end" || !ab57Bytes.Equal(ab57Retained, ab57Raw) {
				ab57T.Fatal("opaque representation")
			}
		}
	}
}

func TestAb57OutputAlternatives(ab57T *ab57Testing.T) {
	for _, ab57Case := range []struct {
		Ab57Variant     string
		Ab57Mask        byte
		Ab57Reserved    byte
		Ab57Active      byte
		Ab57RawReserved byte
	}{
		{"backlight-only", 1, 254, 1, 6}, {"three-output", 7, 248, 7, 0},
	} {
		ab57Value := ab57Entry(ab57T, ab57Key(ab57Case.Ab57Variant, ab57Byte(2), "outputs"))
		if ab57Value.Ab57OutputMask == nil || ab57Value.Ab57ReservedMask == nil || *ab57Value.Ab57OutputMask != ab57Case.Ab57Mask || *ab57Value.Ab57ReservedMask != ab57Case.Ab57Reserved {
			ab57T.Fatal("variant masks")
		}
		ab57Raw := byte(7)
		if ab57Raw&*ab57Value.Ab57OutputMask != ab57Case.Ab57Active || ab57Raw&*ab57Value.Ab57ReservedMask != ab57Case.Ab57RawReserved || (ab57Raw&*ab57Value.Ab57OutputMask)|(ab57Raw&*ab57Value.Ab57ReservedMask) != ab57Raw {
			ab57T.Fatal("raw07 alternative or reserved bits lost")
		}
	}
}

func TestAb57ChecksumScalar(ab57T *ab57Testing.T) {
	ab57Value := ab57Entry(ab57T, ab57Key("eeprom-checksum", nil, "checksum_type"))
	if ab57Value.Ab57ChecksumInterpretation != "scalar, not flags" || len(ab57Value.Ab57ChecksumValues) != 128 || ab57Value.Ab57ChecksumWidth != nil || ab57Value.Ab57Algorithm != Ab05Unknown || ab57Value.Ab57Completion != Ab05Unknown || ab57Value.Ab57ChecksumData != Ab05Unknown {
		ab57T.Fatal("checksum properties conflated")
	}
	for ab57Raw := 0; ab57Raw < 256; ab57Raw++ {
		ab57Label, ab57OK := ab57Value.Ab57ChecksumValues[byte(ab57Raw)]
		ab57Want := "scalar-type"
		if ab57Raw == 0 {
			ab57Want = "query-supported-types"
		}
		if ab57OK != (ab57Raw < 128) || (ab57OK && ab57Label != ab57Want) {
			ab57T.Fatalf("checksum %02X must remain one scalar declaration or unmapped", ab57Raw)
		}
	}
}

func TestAb57FreshCopies(ab57T *ab57Testing.T) {
	ab57Before := Ab57B506Candidates()
	ab57Changed := Ab57B506Candidates()
	ab57Changed.Ab57Limitations[0] = "ab57_changed"
	*ab57Changed.Ab57Entries[0].Ab57Key.Ab57Secondary = 0
	*ab57Changed.Ab57Entries[0].Ab57Key.Ab57Operation = 0
	*ab57Changed.Ab57Entries[0].Ab57Parameters[5].Ab57Metadata.Ab05Width = 99
	ab57Changed.Ab57Entries[0].Ab57Parameters[16].Ab57Metadata.Ab05Enum[0] = "ab57_changed"
	*ab57Changed.Ab57Entries[1].Ab57OutputMask = 255
	ab57Changed.Ab57Entries[3].Ab57ChecksumValues[3] = "flags"
	if !ab57Reflect.DeepEqual(Ab57B506Candidates(), ab57Before) {
		ab57T.Fatal("accessor shared nested storage")
	}
	ab57Lookup := ab57Entry(ab57T, ab57Key("production-test", ab57Byte(0x11), "parameter"))
	*ab57Lookup.Ab57Key.Ab57Secondary = 0
	ab57Lookup.Ab57Parameters[16].Ab57Metadata.Ab05Enum[1] = "ab57_changed"
	if !ab57Reflect.DeepEqual(ab57Entry(ab57T, ab57Key("production-test", ab57Byte(0x11), "parameter")), ab57Before.Ab57Entries[0]) {
		ab57T.Fatal("lookup shared storage")
	}
}

func TestAb57ClaimBoundaries(ab57T *ab57Testing.T) {
	ab57Catalog := Ab57B506Candidates()
	if ab57Catalog.Ab57Presence != Ab05Proven || ab57Catalog.Ab57Qualification != Ab05Unknown || len(ab57Catalog.Ab57Limitations) != 6 {
		ab57T.Fatal("claim boundaries")
	}
	for _, ab57Value := range ab57Catalog.Ab57Entries {
		if ab57Value.Ab57Grade != Ab05Hypothesis || ab57Value.Ab57Key.Ab57Direction != nil || ab57Value.Ab57ChecksumWidth != nil || ab57Value.Ab57Algorithm != Ab05Unknown || ab57Value.Ab57Completion != Ab05Unknown || ab57Value.Ab57ChecksumData != Ab05Unknown {
			ab57T.Fatal("unknown properties promoted")
		}
	}
}
