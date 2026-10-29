package vaillant

import (
	ab69reflect "reflect"
	ab69strconv "strconv"
	ab69testing "testing"
)

func ab69Byte(ab69Value byte) *byte { return &ab69Value }

func ab69Get(ab69T *ab69testing.T, ab69Operation byte, ab69Variant, ab69Direction, ab69Field string) Ab69Candidate {
	ab69T.Helper()
	ab69Entry, ab69OK := Ab69Lookup(Ab69Key{"B512", ab69Variant, 181, ab69Byte(18), ab69Byte(ab69Operation), ab69Direction, ab69Field})
	if !ab69OK {
		ab69T.Fatalf("missing OP%02d %s", ab69Operation, ab69Field)
	}
	return ab69Entry
}

func ab69Meaning(ab69Entry Ab69Candidate, ab69Raw int) string {
	if ab69Value, ab69OK := ab69Entry.Ab69Mappings[ab69strconv.Itoa(ab69Raw)]; ab69OK {
		return ab69Value
	}
	for _, ab69Range := range ab69Entry.Ab69Ranges {
		if ab69Raw >= ab69Range.Ab69Min && ab69Raw <= ab69Range.Ab69Max {
			return ab69Range.Ab69Meaning
		}
	}
	return "Unknown"
}

func TestAb69Metadata(ab69T *ab69testing.T) {
	ab69Catalog := Ab69B512Candidates()
	if ab69Catalog.Ab69RepresentationGrade != "Proven" || ab69Catalog.Ab69PhysicalQualification != "Unknown" || ab69Catalog.Ab69InstalledSupport != "Unknown" || ab69Catalog.Ab69ResponseAdaptation != "Unknown" || ab69Catalog.Ab69ActualSelection != "Unknown" || ab69Catalog.Ab69Scope == "" || len(ab69Catalog.Ab69Entries) != 15 {
		ab69T.Fatal("catalog boundary changed")
	}
	for _, ab69Entry := range ab69Catalog.Ab69Entries {
		ab69Found, ab69OK := Ab69Lookup(ab69Entry.Ab69Key)
		if !ab69OK || !ab69reflect.DeepEqual(ab69Found, ab69Entry) || ab69Entry.Ab69Grade != "Hypothesis" {
			ab69T.Fatal("entry/lookup/grade mismatch")
		}
		if *ab69Entry.Ab69Key.Ab69Operation != 4 && (ab69Entry.Ab69TimeoutSeconds != nil || ab69Entry.Ab69TimeoutRestart != nil) {
			ab69T.Fatal("timeout leaked")
		}
		if ab69Entry.Ab69Key.Ab69Field != "virtual_error" && ab69Meaning(ab69Entry, 65535) == "no-error" {
			ab69T.Fatal("virtual error sentinel leaked across fields")
		}
		if *ab69Entry.Ab69Key.Ab69Operation != 5 && ab69Entry.Ab69ExternalMode != nil {
			ab69T.Fatal("external mode leaked")
		}
	}
}

func TestAb69ExactKey(ab69T *ab69testing.T) {
	ab69Control := Ab69Key{"B512", "circulation", 181, ab69Byte(18), ab69Byte(0), "request", "setpoint"}
	ab69Cases := []struct {
		ab69Name   string
		ab69Mutate func(*Ab69Key)
	}{
		{"family", func(ab69Key *Ab69Key) { ab69Key.Ab69Family = "Unknown" }},
		{"variant", func(ab69Key *Ab69Key) { ab69Key.Ab69Variant = "Unknown" }},
		{"primary", func(ab69Key *Ab69Key) { ab69Key.Ab69Primary = 0 }},
		{"secondary-absent", func(ab69Key *Ab69Key) { ab69Key.Ab69Secondary = nil }},
		{"secondary-zero", func(ab69Key *Ab69Key) { ab69Key.Ab69Secondary = ab69Byte(0) }},
		{"operation-absent", func(ab69Key *Ab69Key) { ab69Key.Ab69Operation = nil }},
		{"operation-other", func(ab69Key *Ab69Key) { ab69Key.Ab69Operation = ab69Byte(2) }},
		{"direction", func(ab69Key *Ab69Key) { ab69Key.Ab69Direction = "positive" }},
		{"field", func(ab69Key *Ab69Key) { ab69Key.Ab69Field = "Unknown" }},
	}
	for _, ab69Case := range ab69Cases {
		ab69T.Run("ab69-"+ab69Case.ab69Name, func(ab69T *ab69testing.T) {
			if _, ab69OK := Ab69Lookup(ab69Control); !ab69OK {
				ab69T.Fatal("valid control failed")
			}
			ab69Changed := ab69Control
			ab69Case.ab69Mutate(&ab69Changed)
			if _, ab69OK := Ab69Lookup(ab69Changed); ab69OK {
				ab69T.Fatal("axis ignored or wildcarded")
			}
		})
	}
	if ab69OptionalEqual(nil, ab69Byte(0)) || ab69OptionalEqual(ab69Byte(0), nil) || !ab69OptionalEqual(nil, nil) {
		ab69T.Fatal("absent conflated with zero")
	}
}

func TestAb69PumpValveCases(ab69T *ab69testing.T) {
	ab69Cases := []struct {
		ab69Op                                byte
		ab69Variant, ab69Direction, ab69Field string
		ab69Raw                               int
		ab69Want                              string
	}{
		{0, "circulation", "request", "setpoint", 0, "off"},
		{0, "circulation", "request", "setpoint", 100, "on/speed"},
		{0, "circulation", "request", "setpoint", 101, "unresolved"},
		{0, "circulation", "request", "setpoint", 255, "unresolved"},
		{2, "valve", "request", "valve", 0, "CH"},
		{2, "valve", "request", "valve", 50, "midpoint"},
		{2, "valve", "request", "valve", 100, "DHW"},
		{2, "valve", "request", "valve", 254, "standby"},
		{2, "valve", "request", "valve", 255, "unchanged"},
		{2, "valve", "positive", "contact", 0, "open"},
		{2, "valve", "positive", "contact", 1, "closed"},
		{2, "valve", "positive", "virtual_error", 65535, "no-error"},
		{2, "valve", "positive", "tank_temperature", 65535, "Unknown"},
		{4, "heating-pump", "request", "setpoint", 0, "off regardless of internal demand"},
		{4, "heating-pump", "request", "setpoint", 100, "on/speed"},
		{4, "heating-pump", "request", "setpoint", 255, "replacement/no-remote-control"},
	}
	for _, ab69Case := range ab69Cases {
		ab69Entry := ab69Get(ab69T, ab69Case.ab69Op, ab69Case.ab69Variant, ab69Case.ab69Direction, ab69Case.ab69Field)
		if ab69Meaning(ab69Entry, ab69Case.ab69Raw) != ab69Case.ab69Want {
			ab69T.Fatalf("%+v", ab69Case)
		}
	}
	ab69Circulation := ab69Get(ab69T, 0, "circulation", "request", "setpoint")
	if !ab69reflect.DeepEqual(ab69Circulation.Ab69RawDomain, []int{0, 255}) || *ab69Circulation.Ab69Default != 0 || *ab69Circulation.Ab69Unit != "percent" {
		ab69T.Fatal("circulation domain/default/unit")
	}
	ab69Ack := ab69Get(ab69T, 0, "circulation", "positive", "acknowledgement")
	if ab69Ack.Ab69Position != nil || ab69Ack.Ab69Width != nil || !ab69reflect.DeepEqual(ab69Ack.Ab69Notes, []string{"prefix-only", "physical acknowledgement: Unknown"}) {
		ab69T.Fatal("ack invented")
	}
	ab69Feedback := ab69Get(ab69T, 2, "valve", "positive", "valve_feedback")
	if !ab69reflect.DeepEqual(ab69Feedback.Ab69Alternatives, []string{"actual", "setpoint fallback"}) || !ab69reflect.DeepEqual(ab69Feedback.Ab69Notes, []string{"actual-versus-fallback selection: Unknown"}) {
		ab69T.Fatal("fallback promoted")
	}
	ab69Tank := ab69Get(ab69T, 2, "valve", "positive", "tank_temperature")
	ab69Bytes := []byte{0xf0, 0xff}
	ab69Raw := int(ab69Bytes[0]) | int(ab69Bytes[1])<<8
	if *ab69Tank.Ab69Signed && ab69Raw >= 32768 {
		ab69Raw -= 65536
	}
	if float64(ab69Raw) / *ab69Tank.Ab69Divisor != -1 || *ab69Tank.Ab69ByteOrder != "little-endian" || !ab69reflect.DeepEqual(ab69Tank.Ab69RawDomain, []int{-32768, 32767}) {
		ab69T.Fatal("signed temperature conversion")
	}
	ab69Error := ab69Get(ab69T, 2, "valve", "positive", "virtual_error")
	if *ab69Error.Ab69Signed || *ab69Error.Ab69ByteOrder != "little-endian" || !ab69reflect.DeepEqual(ab69Error.Ab69Notes, []string{"supplements fault-service information"}) {
		ab69T.Fatal("virtual error metadata")
	}
	ab69Heating := ab69Get(ab69T, 4, "heating-pump", "request", "setpoint")
	if *ab69Heating.Ab69TimeoutSeconds != 900 || *ab69Heating.Ab69TimeoutRestart != "every received service" {
		ab69T.Fatal("timeout")
	}
	ab69Accepted := ab69Get(ab69T, 4, "heating-pump", "positive", "feedback")
	if ab69Accepted.Ab69Position != nil {
		ab69T.Fatal("unknown feedback position invented")
	}
	for ab69Raw := 0; ab69Raw <= 255; ab69Raw++ {
		ab69Want := "rejected"
		if ab69Raw == 1 {
			ab69Want = "accepted"
		}
		if ab69Meaning(ab69Accepted, ab69Raw) != ab69Want {
			ab69T.Fatal("acceptance byte")
		}
	}
	for ab69Channel := 1; ab69Channel <= 3; ab69Channel++ {
		ab69Relay := ab69Get(ab69T, 5, "relay", "request", "relay_"+ab69strconv.Itoa(ab69Channel))
		if *ab69Relay.Ab69Position != ab69Channel+2 || *ab69Relay.Ab69Width != 1 || *ab69Relay.Ab69ExternalMode != 8 || ab69Relay.Ab69TimeoutSeconds != nil || ab69Relay.Ab69TimeoutRestart != nil || !ab69reflect.DeepEqual(ab69Relay.Ab69RawDomain, []int{0, 255}) {
			ab69T.Fatal("relay independent metadata")
		}
		for _, ab69Case := range []struct {
			ab69Raw  int
			ab69Want string
		}{{0, "off"}, {1, "on-or-speed"}, {100, "on-or-speed"}, {255, "unchanged"}} {
			if ab69Meaning(ab69Relay, ab69Case.ab69Raw) != ab69Case.ab69Want {
				ab69T.Fatal("relay meaning")
			}
		}
	}
}

func TestAb69CoordinatesPower(ab69T *ab69testing.T) {
	for _, ab69Case := range []struct {
		ab69Op                                byte
		ab69Variant, ab69Direction, ab69Field string
		ab69Position, ab69Width               int
	}{
		{0, "circulation", "request", "setpoint", 3, 1}, {2, "valve", "request", "valve", 3, 1}, {2, "valve", "positive", "valve_feedback", 3, 1}, {2, "valve", "positive", "contact", 4, 1}, {2, "valve", "positive", "tank_temperature", 5, 2}, {2, "valve", "positive", "virtual_error", 7, 2}, {4, "heating-pump", "request", "setpoint", 3, 1}, {7, "power", "positive", "maximum_power", 3, 1}, {7, "power", "positive", "minimum_power", 4, 1},
	} {
		ab69Entry := ab69Get(ab69T, ab69Case.ab69Op, ab69Case.ab69Variant, ab69Case.ab69Direction, ab69Case.ab69Field)
		if *ab69Entry.Ab69Position != ab69Case.ab69Position || *ab69Entry.Ab69Width != ab69Case.ab69Width {
			ab69T.Fatal("coordinate/width")
		}
	}
	ab69Args := ab69Get(ab69T, 7, "power", "request", "arguments")
	if ab69Args.Ab69Position != nil || ab69Args.Ab69Width != nil || !ab69reflect.DeepEqual(ab69Args.Ab69Notes, []string{"no arguments"}) {
		ab69T.Fatal("arguments invented")
	}
	ab69Maximum := ab69Get(ab69T, 7, "power", "positive", "maximum_power")
	if *ab69Maximum.Ab69Default != 0 || *ab69Maximum.Ab69Unit != "kW" || *ab69Maximum.Ab69Signed || !ab69reflect.DeepEqual(ab69Maximum.Ab69RawDomain, []int{0, 255}) {
		ab69T.Fatal("maximum power")
	}
	ab69Minimum := ab69Get(ab69T, 7, "power", "positive", "minimum_power")
	if *ab69Minimum.Ab69Default != 0 || *ab69Minimum.Ab69Signed || *ab69Minimum.Ab69Unit != "percent" || !ab69reflect.DeepEqual(ab69Minimum.Ab69RawDomain, []int{0, 250}) || 250 / *ab69Minimum.Ab69Divisor != 125 || ab69Meaning(ab69Minimum, 251) != "outside declared domain; retain raw" {
		ab69T.Fatal("minimum power clipped or domain expanded")
	}
}

func TestAb69DeepCopy(ab69T *ab69testing.T) {
	ab69Before := Ab69B512Candidates()
	ab69Changed := Ab69B512Candidates()
	ab69Changed.Ab69Entries[0].Ab69Mappings["0"] = "changed"
	ab69Changed.Ab69Entries[0].Ab69Ranges[0].Ab69Min = -1
	ab69Changed.Ab69Entries[0].Ab69RawDomain[1] = 0
	*ab69Changed.Ab69Entries[0].Ab69Position = 99
	*ab69Changed.Ab69Entries[0].Ab69Key.Ab69Operation = 99
	ab69Changed.Ab69Entries[1].Ab69Notes[0] = "changed"
	ab69Changed.Ab69Entries[3].Ab69Alternatives[0] = "changed"
	if !ab69reflect.DeepEqual(ab69Before, Ab69B512Candidates()) {
		ab69T.Fatal("shared accessor storage")
	}
	ab69Lookup := ab69Get(ab69T, 0, "circulation", "request", "setpoint")
	ab69Lookup.Ab69Mappings["0"] = "changed"
	*ab69Lookup.Ab69Width = 99
	if !ab69reflect.DeepEqual(ab69Before.Ab69Entries[0], ab69Get(ab69T, 0, "circulation", "request", "setpoint")) {
		ab69T.Fatal("shared lookup storage")
	}
}
