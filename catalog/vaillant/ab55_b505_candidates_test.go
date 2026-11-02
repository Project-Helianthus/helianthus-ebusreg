package vaillant

import (
	ab55Reflect "reflect"
	ab55Testing "testing"
)

func ab55Byte(ab55Value byte) *byte       { return &ab55Value }
func ab55String(ab55Value string) *string { return &ab55Value }

func ab55Key(ab55Field string, ab55Operation *byte) Ab55Key {
	return Ab55Key{Ab55Family: "vaillant", Ab55Variant: "B505", Ab55Primary: 0xB5, Ab55Secondary: ab55Byte(5), Ab55Operation: ab55Operation, Ab55Field: ab55Field}
}

func ab55Lookup(ab55T *ab55Testing.T, ab55Field string, ab55Operation *byte) Ab55Entry {
	ab55T.Helper()
	ab55Entry, ab55OK := Ab55Lookup(ab55Key(ab55Field, ab55Operation))
	if !ab55OK {
		ab55T.Fatalf("missing %s", ab55Field)
	}
	return ab55Entry
}

func TestAb55CatalogContract(ab55T *ab55Testing.T) {
	ab55Catalog := Ab55B505Candidates()
	if ab55Catalog.Ab55MetadataPresence != "Proven" || ab55Catalog.Ab55Qualification != "Unknown" || len(ab55Catalog.Ab55Entries) != 8 || len(ab55Catalog.Ab55Limitations) != 4 {
		ab55T.Fatal("finite catalog or evidence boundary changed")
	}
	ab55Fields := []string{"temperature_1", "temperature_2", "dhw_setpoint", "collector_temperature", "collector_error", "pump_power_display", "remote_control", "priority"}
	ab55Operations := []*byte{ab55Byte(0x32), ab55Byte(0x32), ab55Byte(0x61), ab55Byte(0x4D), ab55Byte(0x4D), nil, nil, ab55Byte(0x4E)}
	ab55Widths := []*int{new(int), new(int), nil, nil, new(int), nil, new(int), new(int)}
	*ab55Widths[0], *ab55Widths[1], *ab55Widths[4], *ab55Widths[6], *ab55Widths[7] = 16, 16, 8, 8, 8
	for ab55Index, ab55Field := range ab55Fields {
		ab55Entry := ab55Lookup(ab55T, ab55Field, ab55Operations[ab55Index])
		if !ab55Reflect.DeepEqual(ab55Entry, ab55Catalog.Ab55Entries[ab55Index]) || !ab55Reflect.DeepEqual(ab55Entry.Ab55Key, ab55Key(ab55Field, ab55Operations[ab55Index])) {
			ab55T.Fatalf("accessor/lookup identity: %s", ab55Field)
		}
		if ab55Entry.Ab55Position != nil || !ab55Reflect.DeepEqual(ab55Entry.Ab55WidthBits, ab55Widths[ab55Index]) || ab55Entry.Ab55Conversion != "Unknown" || ab55Entry.Ab55Bounds != "Unknown" || ab55Entry.Ab55Grade != "Hypothesis" {
			ab55T.Fatalf("unsupported coordinate, width, conversion or grade: %s", ab55Field)
		}
		if ab55Entry.Ab55Unit != nil || ab55Entry.Ab55Freshness != "Unknown" || ab55Entry.Ab55Execution != "Unknown" || ab55Entry.Ab55PhysicalSupport != "Unknown" || ab55Entry.Ab55Acknowledgement != "Unknown" || ab55Entry.Ab55Acceptance != "Unknown" || ab55Entry.Ab55Persistence != "Unknown" || ab55Entry.Ab55DurationEnforcement != "Unknown" {
			ab55T.Fatalf("qualification promoted: %s", ab55Field)
		}
		ab55HasDuration := ab55Field == "dhw_setpoint" || ab55Field == "collector_temperature" || ab55Field == "collector_error"
		if ab55HasDuration {
			if ab55Entry.Ab55DurationSeconds == nil || *ab55Entry.Ab55DurationSeconds != 900 {
				ab55T.Fatalf("duration: %s", ab55Field)
			}
		} else if ab55Entry.Ab55DurationSeconds != nil {
			ab55T.Fatalf("invented duration: %s", ab55Field)
		}
	}
}

func TestAb55ExactKeyAxes(ab55T *ab55Testing.T) {
	ab55Control := ab55Key("temperature_1", ab55Byte(0x32))
	ab55Mutations := map[string]func(*Ab55Key){
		"ab55_family":           func(ab55Key *Ab55Key) { ab55Key.Ab55Family = "other-family" },
		"ab55_variant":          func(ab55Key *Ab55Key) { ab55Key.Ab55Variant = "other-variant" },
		"ab55_primary":          func(ab55Key *Ab55Key) { ab55Key.Ab55Primary = 0 },
		"ab55_secondary":        func(ab55Key *Ab55Key) { ab55Key.Ab55Secondary = ab55Byte(6) },
		"ab55_secondary_absent": func(ab55Key *Ab55Key) { ab55Key.Ab55Secondary = nil },
		"ab55_secondary_zero":   func(ab55Key *Ab55Key) { ab55Key.Ab55Secondary = ab55Byte(0) },
		"ab55_operation":        func(ab55Key *Ab55Key) { ab55Key.Ab55Operation = ab55Byte(0x33) },
		"ab55_operation_absent": func(ab55Key *Ab55Key) { ab55Key.Ab55Operation = nil },
		"ab55_operation_zero":   func(ab55Key *Ab55Key) { ab55Key.Ab55Operation = ab55Byte(0) },
		"ab55_direction":        func(ab55Key *Ab55Key) { ab55Key.Ab55Direction = ab55String("request") },
		"ab55_direction_empty":  func(ab55Key *Ab55Key) { ab55Key.Ab55Direction = ab55String("") },
		"ab55_field":            func(ab55Key *Ab55Key) { ab55Key.Ab55Field = "other-field" },
	}
	for ab55Name, ab55Mutate := range ab55Mutations {
		ab55T.Run(ab55Name, func(ab55T *ab55Testing.T) {
			ab55Key := ab55Control
			ab55Mutate(&ab55Key)
			if _, ab55OK := Ab55Lookup(ab55Key); ab55OK {
				ab55T.Fatal("changed axis matched")
			}
			if _, ab55OK := Ab55Lookup(ab55Control); !ab55OK {
				ab55T.Fatal("valid control missed")
			}
		})
	}
	for _, ab55Axis := range []string{"secondary", "operation", "direction"} {
		ab55Left, ab55Right := ab55Control, ab55Control
		switch ab55Axis {
		case "secondary":
			ab55Left.Ab55Secondary, ab55Right.Ab55Secondary = nil, ab55Byte(0)
		case "operation":
			ab55Left.Ab55Operation, ab55Right.Ab55Operation = nil, ab55Byte(0)
		case "direction":
			ab55Left.Ab55Direction, ab55Right.Ab55Direction = nil, ab55String("")
		}
		if ab55KeysEqual(ab55Left, ab55Right) || ab55KeysEqual(ab55Right, ab55Left) || !ab55KeysEqual(ab55Left, ab55Left) || !ab55KeysEqual(ab55Right, ab55Right) {
			ab55T.Fatalf("absent and zero collapsed: %s", ab55Axis)
		}
	}
	ab55Absent := ab55Key("remote_control", nil)
	if _, ab55OK := Ab55Lookup(ab55Absent); !ab55OK {
		ab55T.Fatal("absent operation control missed")
	}
	ab55Absent.Ab55Operation = ab55Byte(0)
	if _, ab55OK := Ab55Lookup(ab55Absent); ab55OK {
		ab55T.Fatal("zero operation matched absent")
	}
	ab55Second := ab55Lookup(ab55T, "temperature_2", ab55Byte(0x32))
	if ab55Second.Ab55Key.Ab55Field != "temperature_2" {
		ab55T.Fatal("field identities collapsed")
	}
}

func TestAb55IndependentSentinels(ab55T *ab55Testing.T) {
	ab55First := ab55Lookup(ab55T, "temperature_1", ab55Byte(0x32))
	ab55Second := ab55Lookup(ab55T, "temperature_2", ab55Byte(0x32))
	for _, ab55Entry := range []Ab55Entry{ab55First, ab55Second} {
		if !ab55Reflect.DeepEqual(ab55Entry.Ab55Branches, map[string]string{"8000": "preserve-existing", "0000": "circuit-default"}) || ab55Entry.Ab55Fallback != "numeric-candidate" {
			ab55T.Fatal("sentinel branches changed")
		}
	}
	for _, ab55Case := range []struct {
		Ab55Raw  [2]string
		Ab55Want [2]string
	}{
		{[2]string{"8000", "0000"}, [2]string{"preserve-existing", "circuit-default"}},
		{[2]string{"0000", "8000"}, [2]string{"circuit-default", "preserve-existing"}},
		{[2]string{"8000", "8000"}, [2]string{"preserve-existing", "preserve-existing"}},
		{[2]string{"0001", "0000"}, [2]string{"numeric-candidate", "circuit-default"}},
	} {
		ab55Before := ab55Case.Ab55Raw
		ab55Got := [2]string{ab55First.Ab55SelectRaw(ab55Case.Ab55Raw[0]), ab55Second.Ab55SelectRaw(ab55Case.Ab55Raw[1])}
		if ab55Got != ab55Case.Ab55Want || ab55Case.Ab55Raw != ab55Before {
			ab55T.Fatalf("raw %v: got %v", ab55Case.Ab55Raw, ab55Got)
		}
	}
	ab55Other := ab55Lookup(ab55T, "dhw_setpoint", ab55Byte(0x61))
	if ab55Other.Ab55SelectRaw("8000") != "Unknown" || len(ab55Other.Ab55Branches) != 0 {
		ab55T.Fatal("OP32 sentinel leaked into OP61")
	}
}

func TestAb55OverridesAndPriority(ab55T *ab55Testing.T) {
	ab55Sensor := ab55Lookup(ab55T, "collector_error", ab55Byte(0x4D))
	if !ab55Reflect.DeepEqual(ab55Sensor.Ab55Branches, map[string]string{"55": "short-circuit", "AA": "open"}) || ab55Sensor.Ab55SelectRaw("55") != "short-circuit" || ab55Sensor.Ab55SelectRaw("AA") != "open" || ab55Sensor.Ab55SelectRaw("54") != "Unknown" {
		ab55T.Fatal("sensor mappings changed")
	}
	ab55Pump := ab55Lookup(ab55T, "pump_power_display", nil)
	if !ab55Reflect.DeepEqual(ab55Pump.Ab55DisplaySteps, []int{0, 25, 50, 75, 100}) {
		ab55T.Fatal("display steps changed")
	}
	ab55Remote := ab55Lookup(ab55T, "remote_control", nil)
	if !ab55Reflect.DeepEqual(ab55Remote.Ab55Branches, map[string]string{"01": "enable-associated-output"}) {
		ab55T.Fatal("remote mappings changed")
	}
	for ab55Raw := 0; ab55Raw < 256; ab55Raw++ {
		if ab55Raw != 1 && ab55Remote.Ab55SelectRaw(ab55Hex(byte(ab55Raw))) != "Unknown" {
			ab55T.Fatal("unprovided remote interpretation")
		}
	}
	ab55Priority := ab55Lookup(ab55T, "priority", ab55Byte(0x4E))
	if ab55Priority.Ab55Signedness != "unsigned" || !ab55Reflect.DeepEqual(ab55Priority.Ab55Priority, map[string]string{"0": "lowest", "255": "highest"}) || len(ab55Priority.Ab55Branches) != 0 || ab55Priority.Ab55Fallback != "Unknown" {
		ab55T.Fatal("priority endpoints or temperature independence changed")
	}
}

func ab55Hex(ab55Raw byte) string {
	const ab55Digits = "0123456789ABCDEF"
	return string([]byte{ab55Digits[ab55Raw>>4], ab55Digits[ab55Raw&15]})
}

func TestAb55DeepCopy(ab55T *ab55Testing.T) {
	ab55Before := Ab55B505Candidates()
	ab55Changed := Ab55B505Candidates()
	ab55Changed.Ab55Limitations[0] = "changed"
	for ab55Index := range ab55Changed.Ab55Entries {
		ab55Entry := &ab55Changed.Ab55Entries[ab55Index]
		*ab55Entry.Ab55Key.Ab55Secondary = 0
		if ab55Entry.Ab55Key.Ab55Operation != nil {
			*ab55Entry.Ab55Key.Ab55Operation = 0
		}
		if ab55Entry.Ab55WidthBits != nil {
			*ab55Entry.Ab55WidthBits = 0
		}
		if ab55Entry.Ab55DurationSeconds != nil {
			*ab55Entry.Ab55DurationSeconds = 0
		}
		ab55Entry.Ab55Branches["8000"] = "changed"
		if ab55Entry.Ab55Priority != nil {
			ab55Entry.Ab55Priority["0"] = "changed"
		}
		if ab55Entry.Ab55DisplaySteps != nil {
			ab55Entry.Ab55DisplaySteps[0] = 99
		}
		ab55Entry.Ab55Key.Ab55Field = "changed"
	}
	ab55LookupCopy := ab55Lookup(ab55T, "temperature_1", ab55Byte(0x32))
	ab55LookupCopy.Ab55Branches["8000"] = "changed"
	*ab55LookupCopy.Ab55Key.Ab55Operation = 0
	*ab55LookupCopy.Ab55WidthBits = 0
	if !ab55Reflect.DeepEqual(Ab55B505Candidates(), ab55Before) || !ab55Reflect.DeepEqual(ab55Lookup(ab55T, "temperature_1", ab55Byte(0x32)), ab55Before.Ab55Entries[0]) {
		ab55T.Fatal("shared catalog backing storage")
	}
}
