package vaillant

import (
	ab13json "encoding/json"
	ab13reflect "reflect"
	ab13testing "testing"
)

func ab13Pointer[ab13T any](ab13Value ab13T) *ab13T { return &ab13Value }

func ab13Key(ab13Operation uint8, ab13Field string) Ab13Key {
	return Ab13Key{Ab13Family: "vaillant", Ab13Variant: "B507", Ab13Primary: 0xb5, Ab13Secondary: ab13Pointer(uint8(7)), Ab13Operation: ab13Pointer(ab13Operation), Ab13Field: ab13Field}
}

func ab13Get(ab13T *ab13testing.T, ab13Operation uint8, ab13Field string) Ab13Entry {
	ab13T.Helper()
	ab13Entry, ab13Found, ab13Error := Ab13Lookup(ab13Key(ab13Operation, ab13Field))
	if ab13Error != nil || !ab13Found {
		ab13T.Fatalf("lookup: found=%v error=%v", ab13Found, ab13Error)
	}
	return ab13Entry
}

func ab13Install(ab13T *ab13testing.T, ab13Entries []Ab13Entry) {
	ab13T.Helper()
	ab13Saved := ab13Data
	ab13T.Cleanup(func() { ab13Data = ab13Saved })
	ab13Encoded, ab13Error := ab13json.Marshal(ab13Entries)
	if ab13Error != nil {
		ab13T.Fatal(ab13Error)
	}
	ab13Data = ab13Encoded
}

func TestAb13FiniteMetadata(ab13T *ab13testing.T) {
	ab13Entries, ab13Error := Ab13Entries()
	if ab13Error != nil || len(ab13Entries) != 6 {
		ab13T.Fatalf("entries=%d error=%v", len(ab13Entries), ab13Error)
	}
	for _, ab13Expected := range []struct {
		ab13Operation        uint8
		ab13Field, ab13Grade string
	}{
		{6, "capacity", "Hypothesis"}, {6, "cost", "Hypothesis"}, {6, "temperature", "Hypothesis"},
		{8, "preference_word", "Hypothesis"}, {9, "preference_word", "Unknown"}, {11, "preference_word", "Unknown"},
	} {
		ab13Entry := ab13Get(ab13T, ab13Expected.ab13Operation, ab13Expected.ab13Field)
		if ab13Entry.Ab13Grade != ab13Expected.ab13Grade || ab13Entry.Ab13MetadataPresence != "Proven" || ab13Entry.Ab13PhysicalQualification != "Unknown" || ab13Entry.Ab13Invocation != "Unknown" || !ab13Entry.Ab13PreserveRaw {
			ab13T.Fatal("grade/qualification/raw boundary")
		}
		if ab13Entry.Ab13Position != nil || ab13Entry.Ab13Key.Ab13Direction != nil || ab13Entry.Ab13Key.Ab13Group != nil || ab13Entry.Ab13Key.Ab13RequestedInstance != nil || ab13Entry.Ab13Key.Ab13Register != nil || ab13Entry.Ab13Key.Ab13CommandRole != nil || ab13Entry.Ab13Reference != "Unknown" {
			ab13T.Fatal("invented coordinates or reference")
		}
		if ab13Entry.Ab13NativeContract != "https://github.com/Project-Helianthus/helianthus-docs-ebus/blob/main/protocols/vaillant/ebus-vaillant-B507.md" {
			ab13T.Fatal("native contract")
		}
	}
}

func TestAb13Capacity(ab13T *ab13testing.T) {
	ab13Entry := ab13Get(ab13T, 6, "capacity")
	if ab13Entry.Ab13WidthBits == nil || *ab13Entry.Ab13WidthBits != 8 || ab13Entry.Ab13Encoding == nil || *ab13Entry.Ab13Encoding != "unsigned" || ab13Entry.Ab13Divisor == nil || *ab13Entry.Ab13Divisor != 10 || ab13Entry.Ab13Unit == nil || *ab13Entry.Ab13Unit != "kW" || !ab13reflect.DeepEqual(ab13Entry.Ab13RawDomain, []uint16{0, 250}) {
		ab13T.Fatal("capacity candidate")
	}
	for _, ab13Case := range []struct {
		ab13Raw    uint16
		ab13Inside bool
		ab13KW     float64
	}{{0, true, 0}, {1, true, 0.1}, {250, true, 25}, {251, false, 0}} {
		ab13Inside := ab13Case.ab13Raw >= ab13Entry.Ab13RawDomain[0] && ab13Case.ab13Raw <= ab13Entry.Ab13RawDomain[1]
		if ab13Inside != ab13Case.ab13Inside {
			ab13T.Fatal("out-of-domain must not clip")
		}
		if ab13Inside && float64(ab13Case.ab13Raw) / *ab13Entry.Ab13Divisor != ab13Case.ab13KW {
			ab13T.Fatal("unsigned division candidate")
		}
	}
	if !ab13reflect.DeepEqual(ab13Entry.Ab13Alternatives, []string{"Negative mode-related prose remains unresolved; no signed reinterpretation selected."}) {
		ab13T.Fatal("unresolved alternative lost")
	}
}

func TestAb13OperationSentinels(ab13T *ab13testing.T) {
	for _, ab13Operation := range []uint8{8, 9, 11} {
		ab13Entry := ab13Get(ab13T, ab13Operation, "preference_word")
		for ab13Index, ab13Raw := range []uint16{0, 65535, 1} {
			ab13Want := ""
			if ab13Operation == 8 {
				ab13Want = []string{"appliance-unavailable candidate", "maximum-efficiency-preference", ""}[ab13Index]
			}
			ab13Got, ab13Mapped := ab13Entry.Ab13Values[ab13Raw]
			if ab13Got != ab13Want || ab13Mapped != (ab13Want != "") {
				ab13T.Fatal("operation-local label", ab13Operation, ab13Raw)
			}
		}
		if ab13Operation == 8 {
			if ab13Entry.Ab13WidthBits == nil || *ab13Entry.Ab13WidthBits != 16 || len(ab13Entry.Ab13Values) != 2 {
				ab13T.Fatal("OP08 word")
			}
		} else if ab13Entry.Ab13WidthBits != nil || ab13Entry.Ab13Values != nil {
			ab13T.Fatal("transferred OP08 metadata")
		}
		if ab13Entry.Ab13Divisor != nil || ab13Entry.Ab13Encoding != nil || ab13Entry.Ab13Unit != nil {
			ab13T.Fatal("invented conversion")
		}
	}
}

func TestAb13CostContexts(ab13T *ab13testing.T) {
	for _, ab13Mode := range []string{"synthetic-mode-a", "synthetic-mode-b"} {
		for _, ab13Reference := range []string{"monetary", "emissions"} {
			ab13T.Run("ab13_"+ab13Mode+"_"+ab13Reference, func(ab13T *ab13testing.T) {
				ab13Raw := uint16(173)
				ab13Entry := ab13Get(ab13T, 6, "cost")
				if !ab13reflect.DeepEqual(ab13Entry.Ab13Dependencies, []string{"operating_mode", "monetary_or_emissions_reference"}) || ab13Entry.Ab13Reference != "Unknown" || ab13Entry.Ab13QuantityOrigin != "calculated_or_estimated_possible" || ab13Entry.Ab13Divisor != nil || ab13Entry.Ab13WidthBits != nil || ab13Entry.Ab13Encoding != nil || ab13Entry.Ab13Unit != nil || ab13Entry.Ab13Values != nil || ab13Entry.Ab13RawDomain != nil || !ab13Entry.Ab13PreserveRaw || ab13Raw != 173 {
					ab13T.Fatal("cost must remain raw with unresolved reference and converter")
				}
			})
		}
	}
	ab13Temperature := ab13Get(ab13T, 6, "temperature")
	if ab13Temperature.Ab13QuantityOrigin != "calculated_or_estimated_possible" || ab13Temperature.Ab13WidthBits != nil || ab13Temperature.Ab13Divisor != nil || ab13Temperature.Ab13Encoding != nil || ab13Temperature.Ab13Unit != nil || ab13Temperature.Ab13Values != nil || ab13Temperature.Ab13RawDomain != nil {
		ab13T.Fatal("temperature is not qualified as direct measurement")
	}
}

func TestAb13ExactKeyAxes(ab13T *ab13testing.T) {
	ab13Entry := ab13Get(ab13T, 6, "capacity")
	ab13Entry.Ab13Key = Ab13Key{Ab13Family: "synthetic-family", Ab13Variant: "B524", Ab13Primary: 181, Ab13Secondary: ab13Pointer(uint8(0)), Ab13Operation: ab13Pointer(uint8(0)), Ab13Direction: ab13Pointer("request"), Ab13Field: "synthetic-field", Ab13Group: ab13Pointer(uint8(0)), Ab13RequestedInstance: ab13Pointer(uint8(0)), Ab13Register: ab13Pointer(uint16(0)), Ab13CommandRole: ab13Pointer("read")}
	ab13Install(ab13T, []Ab13Entry{ab13Entry})
	for ab13Name, ab13Change := range map[string]func(*Ab13Key){
		"ab13_family":                    func(ab13Key *Ab13Key) { ab13Key.Ab13Family = "other-family" },
		"ab13_variant":                   func(ab13Key *Ab13Key) { ab13Key.Ab13Variant = "B507" },
		"ab13_primary":                   func(ab13Key *Ab13Key) { ab13Key.Ab13Primary = 182 },
		"ab13_secondary":                 func(ab13Key *Ab13Key) { ab13Key.Ab13Secondary = ab13Pointer(uint8(1)) },
		"ab13_secondary_absent":          func(ab13Key *Ab13Key) { ab13Key.Ab13Secondary = nil },
		"ab13_operation":                 func(ab13Key *Ab13Key) { ab13Key.Ab13Operation = ab13Pointer(uint8(1)) },
		"ab13_operation_absent":          func(ab13Key *Ab13Key) { ab13Key.Ab13Operation = nil },
		"ab13_direction":                 func(ab13Key *Ab13Key) { ab13Key.Ab13Direction = ab13Pointer("response") },
		"ab13_direction_absent":          func(ab13Key *Ab13Key) { ab13Key.Ab13Direction = nil },
		"ab13_field":                     func(ab13Key *Ab13Key) { ab13Key.Ab13Field = "other-field" },
		"ab13_group":                     func(ab13Key *Ab13Key) { ab13Key.Ab13Group = ab13Pointer(uint8(1)) },
		"ab13_group_absent":              func(ab13Key *Ab13Key) { ab13Key.Ab13Group = nil },
		"ab13_requested_instance":        func(ab13Key *Ab13Key) { ab13Key.Ab13RequestedInstance = ab13Pointer(uint8(1)) },
		"ab13_requested_instance_absent": func(ab13Key *Ab13Key) { ab13Key.Ab13RequestedInstance = nil },
		"ab13_register":                  func(ab13Key *Ab13Key) { ab13Key.Ab13Register = ab13Pointer(uint16(1)) },
		"ab13_register_high_byte":        func(ab13Key *Ab13Key) { ab13Key.Ab13Register = ab13Pointer(uint16(256)) },
		"ab13_register_absent":           func(ab13Key *Ab13Key) { ab13Key.Ab13Register = nil },
		"ab13_command_role":              func(ab13Key *Ab13Key) { ab13Key.Ab13CommandRole = ab13Pointer("write") },
		"ab13_command_role_absent":       func(ab13Key *Ab13Key) { ab13Key.Ab13CommandRole = nil },
	} {
		ab13T.Run(ab13Name, func(ab13T *ab13testing.T) {
			ab13Control, ab13Found, ab13Error := Ab13Lookup(ab13Entry.Ab13Key)
			if ab13Error != nil || !ab13Found || !ab13reflect.DeepEqual(ab13Control, ab13Entry) {
				ab13T.Fatal("positive control", ab13Error)
			}
			ab13Changed := ab13Entry.Ab13Key
			ab13Change(&ab13Changed)
			ab13Got, ab13Found, ab13Error := Ab13Lookup(ab13Changed)
			if ab13Error != nil || ab13Found || !ab13reflect.DeepEqual(ab13Got, Ab13Entry{}) {
				ab13T.Fatal("axis omitted or wildcarded", ab13Error)
			}
		})
	}
}

func TestAb13FreshOwnership(ab13T *ab13testing.T) {
	ab13Baseline, ab13Error := Ab13Entries()
	if ab13Error != nil {
		ab13T.Fatal(ab13Error)
	}
	ab13Entries, ab13Error := Ab13Entries()
	if ab13Error != nil {
		ab13T.Fatal(ab13Error)
	}
	*ab13Entries[0].Ab13Key.Ab13Secondary = 0
	*ab13Entries[0].Ab13Key.Ab13Operation = 0
	*ab13Entries[0].Ab13WidthBits = 0
	*ab13Entries[0].Ab13Encoding = "signed8"
	ab13Entries[0].Ab13RawDomain[1] = 255
	*ab13Entries[0].Ab13Divisor = 0.1
	*ab13Entries[0].Ab13Unit = "other"
	ab13Entries[0].Ab13Unknowns[0] = "other"
	ab13Entries[0].Ab13Alternatives[0] = "other"
	ab13Entries[1].Ab13Dependencies[0] = "other"
	ab13Entries[3].Ab13Values[0] = "other"
	ab13Lookup := ab13Get(ab13T, 8, "preference_word")
	ab13Lookup.Ab13Values[0] = "other"
	*ab13Lookup.Ab13Key.Ab13Operation = 0
	ab13Lookup.Ab13Unknowns[0] = "other"
	ab13Again, ab13Error := Ab13Entries()
	if ab13Error != nil || !ab13reflect.DeepEqual(ab13Again, ab13Baseline) || !ab13reflect.DeepEqual(ab13Get(ab13T, 8, "preference_word"), ab13Baseline[3]) {
		ab13T.Fatal("shared storage", ab13Error)
	}
}

func TestAb13OptionalPointerOwnership(ab13T *ab13testing.T) {
	ab13Entry := ab13Get(ab13T, 6, "capacity")
	ab13Entry.Ab13Key.Ab13Direction = ab13Pointer("synthetic-direction")
	ab13Entry.Ab13Key.Ab13Group = ab13Pointer(uint8(0))
	ab13Entry.Ab13Key.Ab13RequestedInstance = ab13Pointer(uint8(0))
	ab13Entry.Ab13Key.Ab13Register = ab13Pointer(uint16(0))
	ab13Entry.Ab13Key.Ab13CommandRole = ab13Pointer("synthetic-role")
	ab13Entry.Ab13Position = ab13Pointer(uint16(0))
	ab13Install(ab13T, []Ab13Entry{ab13Entry})
	ab13Entries, ab13Error := Ab13Entries()
	if ab13Error != nil {
		ab13T.Fatal(ab13Error)
	}
	ab13Lookup, ab13Found, ab13Error := Ab13Lookup(ab13Entry.Ab13Key)
	if ab13Error != nil || !ab13Found {
		ab13T.Fatal("positive control", ab13Error)
	}
	for _, ab13Changed := range []*Ab13Entry{&ab13Entries[0], &ab13Lookup} {
		*ab13Changed.Ab13Key.Ab13Direction = "changed"
		*ab13Changed.Ab13Key.Ab13Group = 1
		*ab13Changed.Ab13Key.Ab13RequestedInstance = 1
		*ab13Changed.Ab13Key.Ab13Register = 1
		*ab13Changed.Ab13Key.Ab13CommandRole = "changed"
		*ab13Changed.Ab13Position = 1
	}
	ab13Again, ab13Found, ab13Error := Ab13Lookup(ab13Entry.Ab13Key)
	if ab13Error != nil || !ab13Found || !ab13reflect.DeepEqual(ab13Again, ab13Entry) {
		ab13T.Fatal("optional pointers share storage", ab13Error)
	}
}

func TestAb13MalformedData(ab13T *ab13testing.T) {
	ab13Saved := ab13Data
	ab13T.Cleanup(func() { ab13Data = ab13Saved })
	for _, ab13Invalid := range []string{"{", "null", "[]", "[{}]", string(ab13Saved) + " true", `[{"unexpected":1}]`} {
		ab13Data = []byte(ab13Invalid)
		ab13Entries, ab13Error := Ab13Entries()
		if ab13Error == nil || ab13Entries != nil {
			ab13T.Fatal("malformed entries succeeded")
		}
		ab13Entry, ab13Found, ab13Error := Ab13Lookup(ab13Key(6, "capacity"))
		if ab13Error == nil || ab13Found || !ab13reflect.DeepEqual(ab13Entry, Ab13Entry{}) {
			ab13T.Fatal("malformed lookup succeeded")
		}
	}
	ab13Data = ab13Saved
	ab13Entries, ab13Error := Ab13Entries()
	if ab13Error != nil {
		ab13T.Fatal(ab13Error)
	}
	ab13Entries = append(ab13Entries, ab13Entries[0])
	ab13Install(ab13T, ab13Entries)
	if _, ab13Error := Ab13Entries(); ab13Error == nil {
		ab13T.Fatal("duplicate key accepted")
	}
}
