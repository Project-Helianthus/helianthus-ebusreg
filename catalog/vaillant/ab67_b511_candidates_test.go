package vaillant

import (
	ab67reflect "reflect"
	ab67testing "testing"
)

func ab67Byte(ab67Value uint8) *uint8 { return &ab67Value }

func ab67Entry(ab67T *ab67testing.T, ab67Variant string, ab67Operation *uint8, ab67Direction, ab67Field string) Ab67B511Entry {
	ab67T.Helper()
	ab67Result, ab67Found, ab67Error := Ab67B511Lookup(Ab67B511Key{
		Ab67Family: "vaillant", Ab67Variant: ab67Variant, Ab67Operation: ab67Operation,
		Ab67Direction: ab67Direction, Ab67Field: ab67Field,
	})
	if ab67Error != nil || !ab67Found {
		ab67T.Fatalf("lookup %s/%s: found=%v error=%v", ab67Variant, ab67Field, ab67Found, ab67Error)
	}
	return ab67Result
}

func TestAb67ExactKeyAxes(ab67T *ab67testing.T) {
	ab67Control := Ab67B511Key{Ab67Family: "vaillant", Ab67Variant: "OP0C", Ab67Operation: ab67Byte(12), Ab67Direction: "request", Ab67Field: "mode"}
	ab67Changes := map[string]func(*Ab67B511Key){
		"family":                func(ab67Key *Ab67B511Key) { ab67Key.Ab67Family = "other-family" },
		"variant":               func(ab67Key *Ab67B511Key) { ab67Key.Ab67Variant = "OP04" },
		"primary-absent-zero":   func(ab67Key *Ab67B511Key) { ab67Key.Ab67Primary = ab67Byte(0) },
		"secondary-absent-zero": func(ab67Key *Ab67B511Key) { ab67Key.Ab67Secondary = ab67Byte(0) },
		"operation":             func(ab67Key *Ab67B511Key) { ab67Key.Ab67Operation = ab67Byte(4) },
		"operation-absent":      func(ab67Key *Ab67B511Key) { ab67Key.Ab67Operation = nil },
		"direction":             func(ab67Key *Ab67B511Key) { ab67Key.Ab67Direction = "response" },
		"field":                 func(ab67Key *Ab67B511Key) { ab67Key.Ab67Field = "flags" },
	}
	for ab67Name, ab67Change := range ab67Changes {
		ab67T.Run(ab67Name, func(ab67T *ab67testing.T) {
			if ab67Value, ab67Found, ab67Error := Ab67B511Lookup(ab67Control); ab67Error != nil || !ab67Found || ab67Value.Ab67Values[2] != "request-enter" {
				ab67T.Fatalf("valid independent control failed: %v %v", ab67Found, ab67Error)
			}
			ab67Key := ab67Control
			ab67Change(&ab67Key)
			if _, ab67Found, ab67Error := Ab67B511Lookup(ab67Key); ab67Found || ab67Error != nil {
				ab67T.Fatalf("changed axis matched: %v %v", ab67Found, ab67Error)
			}
		})
	}
	for ab67Axis := 0; ab67Axis < 4; ab67Axis++ {
		ab67Unknown := ab67Control
		switch ab67Axis {
		case 0:
			ab67Unknown.Ab67Family = ""
		case 1:
			ab67Unknown.Ab67Variant = ""
		case 2:
			ab67Unknown.Ab67Direction = ""
		case 3:
			ab67Unknown.Ab67Field = ""
		}
		if _, ab67Found, ab67Error := Ab67B511Lookup(ab67Unknown); ab67Found || ab67Error != nil {
			ab67T.Fatal("unknown string axis wildcard-matched")
		}
	}

	ab67Zero := ab67Entry(ab67T, "OP00-reserved", ab67Byte(0), "response", "status_byte_one").Ab67Key
	ab67Zero.Ab67Operation = nil
	if _, ab67Found, ab67Error := Ab67B511Lookup(ab67Zero); ab67Found || ab67Error != nil {
		ab67T.Fatal("absent operation matched zero")
	}
	ab67Partial := Ab67B511Key{Ab67Family: "vaillant", Ab67Variant: "10", Ab67Primary: ab67Byte(0xB5), Ab67Direction: "request", Ab67Field: "prefix"}
	if _, ab67Found, ab67Error := Ab67B511Lookup(ab67Partial); !ab67Found || ab67Error != nil {
		ab67T.Fatal("known primary control failed")
	}
	for _, ab67Primary := range []*uint8{nil, ab67Byte(0), ab67Byte(0xF5)} {
		ab67Changed := ab67Partial
		ab67Changed.Ab67Primary = ab67Primary
		if _, ab67Found, ab67Error := Ab67B511Lookup(ab67Changed); ab67Found || ab67Error != nil {
			ab67T.Fatal("changed known primary matched")
		}
	}
}

func TestAb67Requests(ab67T *ab67testing.T) {
	ab67Mode := ab67Entry(ab67T, "OP0C", ab67Byte(12), "request", "mode")
	if *ab67Mode.Ab67Position != 3 || !ab67reflect.DeepEqual(ab67Mode.Ab67Values, map[uint8]string{0: "not-in-standby", 1: "in-standby", 2: "request-enter", 3: "request-exit"}) || ab67Mode.Ab67Default != nil {
		ab67T.Fatal("mode or absent argument changed")
	}
	for _, ab67Operation := range []uint8{7, 9} {
		ab67Variant := map[uint8]string{7: "OP07", 9: "OP09"}[ab67Operation]
		ab67Flags := ab67Entry(ab67T, ab67Variant, ab67Byte(ab67Operation), "request", "flags")
		if *ab67Flags.Ab67Position != 3 || !ab67reflect.DeepEqual(ab67Flags.Ab67Bits, map[uint8]string{0: "assistant", 1: "factory-initialization"}) || !ab67reflect.DeepEqual(ab67Flags.Ab67ReservedBits, []uint8{2, 3, 4, 5, 6, 7}) {
			ab67T.Fatal("flags declarations changed")
		}
		ab67Raw := uint8(3)
		ab67Named := make(map[string]bool)
		for ab67Bit, ab67Label := range ab67Flags.Ab67Bits {
			ab67Named[ab67Label] = ab67Raw&(1<<ab67Bit) != 0
		}
		if !ab67reflect.DeepEqual(ab67Named, map[string]bool{"assistant": true, "factory-initialization": true}) {
			ab67T.Fatal("raw 03 lost a named flag")
		}
	}
	ab67Assistant := ab67Entry(ab67T, "OP0B", ab67Byte(11), "request", "flags")
	if *ab67Assistant.Ab67Position != 3 || !ab67reflect.DeepEqual(ab67Assistant.Ab67Bits, map[uint8]string{0: "assistant"}) {
		ab67T.Fatal("OP0B inherited factory initialization")
	}
	ab67Boolean := ab67Entry(ab67T, "OP04", ab67Byte(4), "request", "boolean")
	if *ab67Boolean.Ab67Position != 3 || !ab67reflect.DeepEqual(ab67Boolean.Ab67Values, map[uint8]string{0: "false", 1: "true"}) {
		ab67T.Fatal("boolean changed")
	}
	for ab67Index, ab67Field := range []string{"temperature", "actuator", "setpoint", "configuration"} {
		ab67Operation := []uint8{10, 10, 13, 13}[ab67Index]
		ab67Variant := []string{"OP0A", "OP0A", "OP0D", "OP0D"}[ab67Index]
		if *ab67Entry(ab67T, ab67Variant, ab67Byte(ab67Operation), "request", ab67Field).Ab67Position != []int{3, 5, 3, 4}[ab67Index] {
			ab67T.Fatal("request coordinate changed")
		}
	}
}

func TestAb67SpecialFunctions(ab67T *ab67testing.T) {
	for ab67Index, ab67Variant := range []string{"OP0F-first", "OP0F-second", "OP0F-third"} {
		ab67Expected := [][]uint8{{0, 1, 3, 4, 5, 6, 7, 8, 10}, {0, 1, 4, 7, 10, 11}, {0, 1, 4, 7, 10}}[ab67Index]
		ab67Special := ab67Entry(ab67T, ab67Variant, ab67Byte(15), "request", "special_function")
		if !ab67reflect.DeepEqual(ab67Special.Ab67Membership, ab67Expected) || ab67Special.Ab67Values != nil {
			ab67T.Fatal("membership merged or labels invented")
		}
		for ab67Index, ab67Field := range []string{"binding", "special_function", "heating", "error_maintenance_cooling_allowed", "wish_temperature", "quick_veto_hour", "quick_veto_minute", "event_end_hour", "event_end_minute"} {
			if *ab67Entry(ab67T, ab67Variant, ab67Byte(15), "request", ab67Field).Ab67Position != ab67Index+3 {
				ab67T.Fatal("OP0F coordinate changed")
			}
		}
		ab67Wish := ab67Entry(ab67T, ab67Variant, ab67Byte(15), "request", "wish_temperature")
		if *ab67Wish.Ab67WidthBits != 8 || float64(1)**ab67Wish.Ab67Scale != 0.5 || *ab67Wish.Ab67Unit != "degrees" {
			ab67T.Fatal("wish raw 01 conversion changed")
		}
		for _, ab67Field := range []string{"quick_veto_hour", "quick_veto_minute", "event_end_hour", "event_end_minute"} {
			ab67Clock := ab67Entry(ab67T, ab67Variant, ab67Byte(15), "request", ab67Field)
			if !ab67reflect.DeepEqual(ab67Clock.Ab67RawDomain, []int{0, 255}) || ab67Clock.Ab67ClockValidation != Ab05Unknown || ab67Clock.Ab67Scale != nil || ab67Clock.Ab67Unit != nil {
				ab67T.Fatal("raw 255 promoted to clock")
			}
		}
	}
}

func TestAb67ResponseAlternatives(ab67T *ab67testing.T) {
	ab67Reserved := ab67Entry(ab67T, "OP00-reserved", ab67Byte(0), "response", "status_byte_one")
	ab67Fan := ab67Entry(ab67T, "OP00-fan", ab67Byte(0), "response", "status_byte_one")
	if len(ab67Reserved.Ab67Bits) != 0 || !ab67reflect.DeepEqual(ab67Reserved.Ab67ReservedBits, []uint8{7}) || !ab67reflect.DeepEqual(ab67Fan.Ab67Bits, map[uint8]string{7: "fan-running"}) || len(ab67Fan.Ab67ReservedBits) != 0 {
		ab67T.Fatal("raw 80 alternatives collapsed")
	}
	for _, ab67Status := range []Ab67B511Entry{ab67Reserved, ab67Fan} {
		if *ab67Status.Ab67Position != 8 || *ab67Status.Ab67Extent != 12 {
			ab67T.Fatal("OP00 abstract positions changed")
		}
	}
	for ab67Index, ab67Variant := range []string{"OP07-base", "OP07-extended"} {
		for _, ab67Field := range []string{"room_setpoint", "overlapping_structure_unknown"} {
			ab67Overlap := ab67Entry(ab67T, ab67Variant, ab67Byte(7), "response", ab67Field)
			if *ab67Overlap.Ab67Position != 7 || *ab67Overlap.Ab67Extent != 8+ab67Index {
				ab67T.Fatal("overlap concatenated")
			}
		}
		ab67Room := ab67Entry(ab67T, ab67Variant, ab67Byte(7), "response", "room_setpoint")
		if *ab67Room.Ab67WidthBits != 8 || *ab67Room.Ab67Scale != 0.5 || *ab67Room.Ab67Unit != "degrees" {
			ab67T.Fatal("room setpoint codec changed")
		}
	}
	ab67Pressure := ab67Entry(ab67T, "OP07-extended", ab67Byte(7), "response", "building_circuit_pressure")
	if *ab67Pressure.Ab67Position != 8 || *ab67Pressure.Ab67Extent != 9 || *ab67Pressure.Ab67Unit != "percent" || !ab67reflect.DeepEqual(ab67Pressure.Ab67RawDomain, []int{0, 100}) {
		ab67T.Fatal("pressure promoted to bar or extent summed")
	}
}

func TestAb67PartialAndAdmission(ab67T *ab67testing.T) {
	ab67Catalog, ab67Error := Ab67B511Catalog()
	if ab67Error != nil {
		ab67T.Fatal(ab67Error)
	}
	if len(ab67Catalog.Ab67Entries) != 47 || ab67Catalog.Ab67MetadataPresence != Ab05Proven || ab67Catalog.Ab67PhysicalQualification != Ab05Unknown {
		ab67T.Fatalf("finite catalog or grade changed: %d", len(ab67Catalog.Ab67Entries))
	}
	for ab67Index, ab67Row := range ab67Catalog.Ab67Entries {
		if ab67Row.Ab67Grade != Ab05Hypothesis || ab67Row.Ab67InstalledSupport != Ab05Unknown || ab67Row.Ab67PhysicalPrefix != Ab05Unknown || ab67Row.Ab67Default != nil || ab67Row.Ab67RawPolicy != "preserve original carrier" || ab67Row.Ab67Key.Ab67Secondary != nil {
			ab67T.Fatal("claim/default/raw boundary changed")
		}
		if ab67Row.Ab67Key.Ab67Variant != "10" && ab67Row.Ab67Key.Ab67Variant != "11" && ab67Row.Ab67Key.Ab67Primary != nil {
			ab67T.Fatal("invented unprovided primary prefix")
		}
		ab67KnownWidth := ab67Row.Ab67Key.Ab67Field == "wish_temperature" || ab67Row.Ab67Key.Ab67Field == "room_setpoint"
		if !ab67KnownWidth && ab67Row.Ab67WidthBits != nil {
			ab67T.Fatal("invented unprovided width")
		}

		for _, ab67Other := range ab67Catalog.Ab67Entries[ab67Index+1:] {
			if ab67reflect.DeepEqual(ab67Row.Ab67Key, ab67Other.Ab67Key) {
				ab67T.Fatal("duplicate exact key")
			}
		}
		if ab67Row.Ab67Key.Ab67Variant == "10" || ab67Row.Ab67Key.Ab67Variant == "11" {
			ab67Expected := uint8(0xB5)
			if ab67Row.Ab67Key.Ab67Direction == "response" {
				ab67Expected = 0xF5
			}
			if ab67Row.Ab67Key.Ab67Operation != nil || *ab67Row.Ab67Key.Ab67Primary != ab67Expected || *ab67Row.Ab67Position != 0 || ab67Row.Ab67Extent != nil || ab67Row.Ab67WidthBits != nil {
				ab67T.Fatal("partial label materialized selector or body")
			}
		}
	}
	for _, ab67Variant := range []string{"10", "11"} {
		for ab67Index, ab67Direction := range []string{"request", "response"} {
			ab67Key := Ab67B511Key{Ab67Family: "vaillant", Ab67Variant: ab67Variant, Ab67Primary: ab67Byte([]uint8{0xB5, 0xF5}[ab67Index]), Ab67Direction: ab67Direction, Ab67Field: "prefix"}
			if _, ab67Found, ab67Error := Ab67B511Lookup(ab67Key); !ab67Found || ab67Error != nil {
				ab67T.Fatal("partial control failed")
			}
			for _, ab67Operation := range []uint8{0, 10, 11, 16, 17} {
				ab67Key.Ab67Operation = ab67Byte(ab67Operation)
				if _, ab67Found, ab67Error := Ab67B511Lookup(ab67Key); ab67Found || ab67Error != nil {
					ab67T.Fatal("digits or absent operation became selector")
				}
			}
		}
	}
	if _, ab67Found, ab67Error := Ab67B511Lookup(Ab67B511Key{}); ab67Found || ab67Error != nil {
		ab67T.Fatal("constructed universal empty poll")
	}
}

func ab67Mutate(ab67Value ab67reflect.Value) {
	switch ab67Value.Kind() {
	case ab67reflect.Ptr:
		if !ab67Value.IsNil() {
			ab67Mutate(ab67Value.Elem())
		}
	case ab67reflect.Struct:
		for ab67Index := 0; ab67Index < ab67Value.NumField(); ab67Index++ {
			ab67Mutate(ab67Value.Field(ab67Index))
		}
	case ab67reflect.Slice:
		for ab67Index := 0; ab67Index < ab67Value.Len(); ab67Index++ {
			ab67Mutate(ab67Value.Index(ab67Index))
		}
	case ab67reflect.Map:
		for _, ab67Key := range ab67Value.MapKeys() {
			ab67Value.SetMapIndex(ab67Key, ab67reflect.ValueOf("mutated"))
		}
	case ab67reflect.String:
		ab67Value.SetString("mutated")
	case ab67reflect.Int:
		ab67Value.SetInt(99)
	case ab67reflect.Uint8:
		ab67Value.SetUint(99)
	case ab67reflect.Float64:
		ab67Value.SetFloat(99)
	}
}

func TestAb67FreshNestedStorage(ab67T *ab67testing.T) {
	ab67Expected, ab67Error := Ab67B511Catalog()
	if ab67Error != nil {
		ab67T.Fatal(ab67Error)
	}
	ab67Changed, ab67Error := Ab67B511Catalog()
	if ab67Error != nil {
		ab67T.Fatal(ab67Error)
	}
	ab67Mutate(ab67reflect.ValueOf(&ab67Changed).Elem())
	ab67Fresh, ab67Error := Ab67B511Catalog()
	if ab67Error != nil || !ab67reflect.DeepEqual(ab67Expected, ab67Fresh) {
		ab67T.Fatal("accessor shares nested storage")
	}
	for _, ab67ExpectedRow := range ab67Expected.Ab67Entries {
		ab67Row, ab67Found, ab67Error := Ab67B511Lookup(ab67ExpectedRow.Ab67Key)
		if !ab67Found || ab67Error != nil || !ab67reflect.DeepEqual(ab67Row, ab67ExpectedRow) {
			ab67T.Fatal("lookup differs from production catalog")
		}
		ab67Mutate(ab67reflect.ValueOf(&ab67Row).Elem())
		ab67Again, ab67Found, ab67Error := Ab67B511Lookup(ab67ExpectedRow.Ab67Key)
		if !ab67Found || ab67Error != nil || !ab67reflect.DeepEqual(ab67Again, ab67ExpectedRow) {
			ab67T.Fatal("lookup shares nested storage")
		}
	}
}
