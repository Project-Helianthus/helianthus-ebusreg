package vaillant

import (
	ab53Binary "encoding/binary"
	ab53Reflect "reflect"
	ab53Testing "testing"
)

func ab53Byte(ab53Value uint8) *uint8 { return &ab53Value }

func ab53Key(ab53Variant string, ab53Operation *uint8, ab53Field string) Ab53Key {
	return Ab53Key{Ab53Family: "vaillant", Ab53Variant: ab53Variant, Ab53Primary: 0xb5, Ab53Secondary: ab53Byte(4), Ab53Operation: ab53Operation, Ab53Direction: "response", Ab53Field: ab53Field}
}

func ab53Get(ab53T *ab53Testing.T, ab53KeyValue Ab53Key) Ab53Entry {
	ab53T.Helper()
	ab53Entry, ab53Found, ab53Error := Ab53LookupB504Candidate(ab53KeyValue)
	if ab53Error != nil || !ab53Found {
		ab53T.Fatalf("lookup %+v: found=%v error=%v", ab53KeyValue, ab53Found, ab53Error)
	}
	return ab53Entry
}

func TestAb53ExactKeyAxes(ab53T *ab53Testing.T) {
	ab53Control := ab53Key("outside-temperature-div16", ab53Byte(0x16), "corrected_outside_temperature")
	ab53Get(ab53T, ab53Control)
	ab53Cases := []struct {
		Ab53Name   string
		Ab53Change func(*Ab53Key)
	}{
		{"ab53Family", func(ab53K *Ab53Key) { ab53K.Ab53Family = "other-family" }},
		{"ab53Variant", func(ab53K *Ab53Key) { ab53K.Ab53Variant = "other-variant" }},
		{"ab53Primary", func(ab53K *Ab53Key) { ab53K.Ab53Primary = 0xb4 }},
		{"ab53Secondary", func(ab53K *Ab53Key) { ab53K.Ab53Secondary = ab53Byte(5) }},
		{"ab53SecondaryAbsent", func(ab53K *Ab53Key) { ab53K.Ab53Secondary = nil }},
		{"ab53SecondaryZero", func(ab53K *Ab53Key) { ab53K.Ab53Secondary = ab53Byte(0) }},
		{"ab53Operation", func(ab53K *Ab53Key) { ab53K.Ab53Operation = ab53Byte(0x17) }},
		{"ab53OperationAbsent", func(ab53K *Ab53Key) { ab53K.Ab53Operation = nil }},
		{"ab53OperationZero", func(ab53K *Ab53Key) { ab53K.Ab53Operation = ab53Byte(0) }},
		{"ab53Direction", func(ab53K *Ab53Key) { ab53K.Ab53Direction = "request" }},
		{"ab53Field", func(ab53K *Ab53Key) { ab53K.Ab53Field = "other_field" }},
	}
	for _, ab53Case := range ab53Cases {
		ab53T.Run(ab53Case.Ab53Name, func(ab53T *ab53Testing.T) {
			ab53K := ab53Control
			ab53Case.Ab53Change(&ab53K)
			_, ab53Found, ab53Error := Ab53LookupB504Candidate(ab53K)
			if ab53Error != nil || ab53Found {
				ab53T.Fatalf("altered axis matched: %v %v", ab53Found, ab53Error)
			}
		})
	}
	ab53Unknown := ab53Key("solar-scalar-two", nil, "solar_state")
	ab53Get(ab53T, ab53Unknown)
	ab53Unknown.Ab53Operation = ab53Byte(0)
	if _, ab53Found, ab53Error := Ab53LookupB504Candidate(ab53Unknown); ab53Error != nil || ab53Found {
		ab53T.Fatal("absent operation matched explicit zero")
	}
}

func TestAb53TemperatureAlternatives(ab53T *ab53Testing.T) {
	ab53Cases := []struct {
		Ab53Variant  string
		Ab53Divisor  int
		Ab53Positive float64
		Ab53Negative float64
	}{
		{"outside-temperature-div16", 16, 1, -1},
		{"outside-temperature-div256", 256, 0.0625, -0.0625},
	}
	for _, ab53Case := range ab53Cases {
		ab53T.Run("ab53"+ab53Case.Ab53Variant, func(ab53T *ab53Testing.T) {
			ab53Entry := ab53Get(ab53T, ab53Key(ab53Case.Ab53Variant, ab53Byte(0x16), "corrected_outside_temperature"))
			if ab53Entry.Ab53Divisor == nil || *ab53Entry.Ab53Divisor != ab53Case.Ab53Divisor || ab53Entry.Ab53Signedness != "signed" || ab53Entry.Ab53ByteOrder != "little-endian" || ab53Entry.Ab53Width == nil || *ab53Entry.Ab53Width != 2 || ab53Entry.Ab53Unit == nil || *ab53Entry.Ab53Unit != "degrees Celsius" {
				ab53T.Fatal("temperature candidate changed", ab53Entry)
			}
			ab53Positive := float64(int16(ab53Binary.LittleEndian.Uint16([]byte{0x10, 0}))) / float64(*ab53Entry.Ab53Divisor)
			ab53Negative := float64(int16(ab53Binary.LittleEndian.Uint16([]byte{0xf0, 0xff}))) / float64(*ab53Entry.Ab53Divisor)
			if ab53Positive != ab53Case.Ab53Positive || ab53Negative != ab53Case.Ab53Negative {
				ab53T.Fatalf("synthetic fields: %v %v", ab53Positive, ab53Negative)
			}
		})
	}
	ab53Tail := ab53Get(ab53T, ab53Key("unused-tail", ab53Byte(0x16), "unused_tail"))
	if ab53Tail.Ab53Position != nil || ab53Tail.Ab53Width != nil || ab53Tail.Ab53Grade != Ab05Unknown {
		ab53T.Fatal("tail acquired an extent or semantics")
	}
}

func ab53Label(ab53Entry Ab53Entry, ab53Raw int) string {
	if ab53Value, ab53Found := ab53Entry.Ab53Labels[ab53Raw]; ab53Found {
		return ab53Value
	}
	if ab53Raw > 0 && ab53Entry.Ab53PositiveLabel != nil {
		return *ab53Entry.Ab53PositiveLabel
	}
	return ""
}

func TestAb53ZeroAndPositiveBranches(ab53T *ab53Testing.T) {
	for _, ab53Case := range []struct {
		Ab53Variant   string
		Ab53Operation uint8
		Ab53Field     string
		Ab53Zero      string
		Ab53Positive  string
	}{
		{"setpoint", 0x20, "setpoint", "no-demand", "nonzero candidate"},
		{"protection", 0x21, "protection_temperature", "off", "protection-active"},
	} {
		ab53Entry := ab53Get(ab53T, ab53Key(ab53Case.Ab53Variant, ab53Byte(ab53Case.Ab53Operation), ab53Case.Ab53Field))
		if ab53Entry.Ab53Labels[0] != ab53Case.Ab53Zero || ab53Entry.Ab53PositiveLabel == nil || *ab53Entry.Ab53PositiveLabel != ab53Case.Ab53Positive {
			ab53T.Fatal("zero/positive distinction lost")
		}
		if ab53Label(ab53Entry, 0) != ab53Case.Ab53Zero || ab53Label(ab53Entry, 1) != ab53Case.Ab53Positive {
			ab53T.Fatal("synthetic raw zero/one branch changed")
		}
		if _, ab53Present := ab53Entry.Ab53Labels[1]; ab53Present || ab53Entry.Ab53Divisor != nil {
			ab53T.Fatal("invented exact nonzero arithmetic")
		}
	}
}

func TestAb53ContextualUnavailable(ab53T *ab53Testing.T) {
	ab53RoomKey := ab53Key("room-temperature", ab53Byte(0x26), "room_temperature")
	ab53Room := ab53Get(ab53T, ab53RoomKey)
	if ab53Room.Ab53Labels[255] != "unavailable" {
		ab53T.Fatal("missing contextual FF")
	}
	ab53RoomKey.Ab53Operation = ab53Byte(0x20)
	if _, ab53Found, ab53Error := Ab53LookupB504Candidate(ab53RoomKey); ab53Error != nil || ab53Found {
		ab53T.Fatal("OP26 label inherited by unrelated operation")
	}
	ab53Other := ab53Get(ab53T, ab53Key("setpoint", ab53Byte(0x20), "setpoint"))
	if _, ab53Present := ab53Other.Ab53Labels[255]; ab53Present {
		ab53T.Fatal("unrelated selector inherited FF label")
	}
	ab53Flag := ab53Get(ab53T, ab53Key("adjustment", ab53Byte(0x26), "adjustment_flag"))
	if ab53Flag.Ab53Labels[1] != "setpoint change permitted" || len(ab53Flag.Ab53Labels) != 1 {
		ab53T.Fatal("adjustment candidate changed")
	}
}

func ab53Contains(ab53Domain *Ab53Domain, ab53Value int) bool {
	if ab53Domain == nil {
		return false
	}
	if ab53Domain.Ab53Values != nil {
		for _, ab53Allowed := range ab53Domain.Ab53Values {
			if ab53Allowed == ab53Value {
				return true
			}
		}
		return false
	}
	if ab53Domain.Ab53Minimum == nil || ab53Domain.Ab53Maximum == nil || ab53Value < *ab53Domain.Ab53Minimum || ab53Value > *ab53Domain.Ab53Maximum {
		return false
	}
	return ab53Domain.Ab53Step == nil || (ab53Value-*ab53Domain.Ab53Minimum)%*ab53Domain.Ab53Step == 0
}

func TestAb53SolarDomains(ab53T *ab53Testing.T) {
	ab53Three := ab53Get(ab53T, ab53Key("solar-scalar-three", nil, "solar_state"))
	ab53Two := ab53Get(ab53T, ab53Key("solar-scalar-two", nil, "solar_state"))
	if !ab53Reflect.DeepEqual(ab53Three.Ab53Domain.Ab53Values, []int{0, 1, 2}) || !ab53Reflect.DeepEqual(ab53Two.Ab53Domain.Ab53Values, []int{0, 1}) {
		ab53T.Fatal("scalar domains merged or changed")
	}
	if !ab53Contains(ab53Three.Ab53Domain, 2) || ab53Contains(ab53Two.Ab53Domain, 2) {
		ab53T.Fatal("code 2 admitted into two-value table")
	}
	ab53Bits := ab53Get(ab53T, ab53Key("solar-bits", nil, "solar_state"))
	if ab53Bits.Ab53Domain != nil || ab53Bits.Ab53Grade != Ab05Unknown {
		ab53T.Fatal("bit descriptions promoted to scalar domain")
	}
}

func TestAb53DurationAlternatives(ab53T *ab53Testing.T) {
	ab53First := ab53Get(ab53T, ab53Key("preheating-five-minute-60", nil, "preheating_duration"))
	ab53Second := ab53Get(ab53T, ab53Key("preheating-300", nil, "preheating_duration"))
	if ab53First.Ab53Domain == nil || ab53First.Ab53Domain.Ab53Step == nil || *ab53First.Ab53Domain.Ab53Step != 5 || ab53Second.Ab53Domain == nil || ab53Second.Ab53Domain.Ab53Step != nil {
		ab53T.Fatal("duration step metadata changed")
	}
	for _, ab53Case := range []struct {
		Ab53Value             int
		Ab53First, Ab53Second bool
	}{
		{0, true, true}, {1, false, true}, {60, true, true}, {65, false, true}, {300, false, true}, {301, false, false}, {-1, false, false},
	} {
		if ab53Contains(ab53First.Ab53Domain, ab53Case.Ab53Value) != ab53Case.Ab53First || ab53Contains(ab53Second.Ab53Domain, ab53Case.Ab53Value) != ab53Case.Ab53Second {
			ab53T.Fatalf("duration %d domains changed", ab53Case.Ab53Value)
		}
	}
}

func TestAb53CatalogBoundariesAndCopies(ab53T *ab53Testing.T) {
	ab53Before, ab53Error := Ab53B504Candidates()
	if ab53Error != nil {
		ab53T.Fatal(ab53Error)
	}
	if len(ab53Before.Ab53Entries) != 12 || ab53Before.Ab53MetadataPresence != Ab05Proven || ab53Before.Ab53PhysicalAdaptation != Ab05Unknown || ab53Before.Ab53InstalledSupport != Ab05Unknown || ab53Before.Ab53PollingPermission != Ab05Unknown || ab53Before.Ab53RoutableInvocation != Ab05Unknown {
		ab53T.Fatal("catalog boundary changed")
	}
	for ab53Index, ab53Entry := range ab53Before.Ab53Entries {
		if ab53Entry.Ab53Position != nil {
			ab53T.Fatal("unprovided coordinate invented")
		}
		ab53ExpectedGrade := Ab05Hypothesis
		if ab53Entry.Ab53Key.Ab53Variant == "unused-tail" || ab53Entry.Ab53Key.Ab53Variant == "solar-bits" {
			ab53ExpectedGrade = Ab05Unknown
		}
		if ab53Entry.Ab53Grade != ab53ExpectedGrade {
			ab53T.Fatal("semantic grade promoted")
		}
		if !ab53Reflect.DeepEqual(ab53Get(ab53T, ab53Entry.Ab53Key), ab53Entry) {
			ab53T.Fatal("lookup/accessor differ")
		}
		for _, ab53Other := range ab53Before.Ab53Entries[ab53Index+1:] {
			if ab53Reflect.DeepEqual(ab53Entry.Ab53Key, ab53Other.Ab53Key) {
				ab53T.Fatal("duplicate key")
			}
		}
	}
	ab53Copy, ab53Error := Ab53B504Candidates()
	if ab53Error != nil {
		ab53T.Fatal(ab53Error)
	}
	ab53Copy.Ab53Limitations[0] = "changed"
	for ab53Index := range ab53Copy.Ab53Entries {
		ab53Entry := &ab53Copy.Ab53Entries[ab53Index]
		ab53Entry.Ab53Key.Ab53Family = "changed"
		if ab53Entry.Ab53Key.Ab53Operation != nil {
			*ab53Entry.Ab53Key.Ab53Operation = 0
		}
		*ab53Entry.Ab53Key.Ab53Secondary = 0
		if ab53Entry.Ab53Width != nil {
			*ab53Entry.Ab53Width = 99
		}
		if ab53Entry.Ab53Divisor != nil {
			*ab53Entry.Ab53Divisor = 99
		}
		if ab53Entry.Ab53Unit != nil {
			*ab53Entry.Ab53Unit = "changed"
		}
		ab53Entry.Ab53Labels[255] = "changed"
		if ab53Entry.Ab53PositiveLabel != nil {
			*ab53Entry.Ab53PositiveLabel = "changed"
		}
		if len(ab53Entry.Ab53Notes) > 0 {
			ab53Entry.Ab53Notes[0] = "changed"
		}
		if ab53Entry.Ab53Domain != nil {
			if len(ab53Entry.Ab53Domain.Ab53Values) > 0 {
				ab53Entry.Ab53Domain.Ab53Values[0] = 99
			}
			if ab53Entry.Ab53Domain.Ab53Minimum != nil {
				*ab53Entry.Ab53Domain.Ab53Minimum = 99
			}
			if ab53Entry.Ab53Domain.Ab53Maximum != nil {
				*ab53Entry.Ab53Domain.Ab53Maximum = 99
			}
			if ab53Entry.Ab53Domain.Ab53Step != nil {
				*ab53Entry.Ab53Domain.Ab53Step = 99
			}
		}
	}
	ab53After, ab53Error := Ab53B504Candidates()
	if ab53Error != nil || !ab53Reflect.DeepEqual(ab53Before, ab53After) {
		ab53T.Fatal("accessor leaked backing storage")
	}
	ab53KeyValue := ab53Key("room-temperature", ab53Byte(0x26), "room_temperature")
	ab53Lookup := ab53Get(ab53T, ab53KeyValue)
	ab53Lookup.Ab53Labels[255] = "changed"
	*ab53Lookup.Ab53Key.Ab53Operation = 0
	if ab53Get(ab53T, ab53KeyValue).Ab53Labels[255] != "unavailable" {
		ab53T.Fatal("lookup leaked backing storage")
	}
}
