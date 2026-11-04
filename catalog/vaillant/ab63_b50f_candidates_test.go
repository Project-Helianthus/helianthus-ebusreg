package vaillant

import (
	ab63fmt "fmt"
	ab63reflect "reflect"
	ab63testing "testing"
)

func ab63Byte(ab63Value byte) *byte       { return &ab63Value }
func ab63String(ab63Value string) *string { return &ab63Value }
func ab63Key(ab63Operation *byte, ab63Field string) Ab63Key {
	return Ab63Key{Ab63Family: "vaillant", Ab63Variant: "B50F", Ab63Primary: 0xB5, Ab63Secondary: ab63Byte(0x0F), Ab63Operation: ab63Operation, Ab63Field: ab63Field}
}
func ab63Get(ab63T *ab63testing.T, ab63Operation *byte, ab63Name string) Ab63Entry {
	ab63T.Helper()
	ab63Entry, ab63OK := Ab63B50FLookup(ab63Key(ab63Operation, ab63Name))
	if !ab63OK {
		ab63T.Fatalf("missing %v/%s", ab63Operation, ab63Name)
	}
	return ab63Entry
}

func TestAb63FiniteCatalog(ab63T *ab63testing.T) {
	ab63Catalog := Ab63B50FCandidates()
	if ab63Catalog.Ab63Presence != Ab05Proven || ab63Catalog.Ab63Qualification != Ab05Unknown || ab63Catalog.Ab63Invocation != Ab05Unknown || len(ab63Catalog.Ab63Entries) != 21 {
		ab63T.Fatal("finite catalog or grade separation")
	}
	ab63Expected := map[string]int{"demand": 1, "stack-voltage": 1, "feedback": 1, "temperature_1": 2, "temperature_1_error": 2, "temperature_2": 2, "temperature_2_error": 2, "temperature_3": 1, "temperature_3_error": 1, "integrity_1": 4, "integrity_2": 4}
	for _, ab63Entry := range ab63Catalog.Ab63Entries {
		ab63Expected[ab63Entry.Ab63Key.Ab63Field]--
		if ab63Entry.Ab63Key.Ab63Family != "vaillant" || ab63Entry.Ab63Key.Ab63Variant != "B50F" || ab63Entry.Ab63Key.Ab63Primary != 0xB5 || ab63Entry.Ab63Key.Ab63Secondary == nil || *ab63Entry.Ab63Key.Ab63Secondary != 0x0F {
			ab63T.Fatal("native key axes")
		}
		if ab63Entry.Ab63Field.Ab05Position == nil && ab63Entry.Ab63Coordinates != "Unknown" {
			ab63T.Fatal("unknown scalar coordinates")
		}
		if ab63Entry.Ab63Field.Ab05Width == nil {
			ab63T.Fatal("known carrier width lost")
		}
		if ab63Entry.Ab63Key.Ab63Field == "demand" || ab63Entry.Ab63Key.Ab63Field == "feedback" {
			if *ab63Entry.Ab63Field.Ab05Width != 1 || ab63Entry.Ab63Field.Ab05Position != nil || ab63Entry.Ab63Field.Ab05Codec != "raw" {
				ab63T.Fatal("raw byte metadata")
			}
		}
		if ab63Entry.Ab63Key.Ab63Field == "stack-voltage" && *ab63Entry.Ab63Field.Ab05Width != 2 {
			ab63T.Fatal("stack raw width")
		}
		if ab63Entry.Ab63Field.Ab05Grade != Ab05Hypothesis || ab63Entry.Ab63Field.Ab05Name != ab63Entry.Ab63Key.Ab63Field || ab63Entry.Ab63Key.Ab63Direction != nil || ab63Entry.Ab63Field.Ab05RawFallback != nil {
			ab63T.Fatal("grade, name, unknown direction or native fallback")
		}
		ab63Found, ab63OK := Ab63B50FLookup(ab63Entry.Ab63Key)
		if !ab63OK || !ab63reflect.DeepEqual(ab63Found, ab63Entry) {
			ab63T.Fatal("accessor/lookup mismatch")
		}
	}
	for ab63Name, ab63Count := range ab63Expected {
		if ab63Count != 0 {
			ab63T.Fatalf("unexpected field count %s: %d", ab63Name, ab63Count)
		}
	}
}

func TestAb63ExactKeyAxes(ab63T *ab63testing.T) {
	for ab63Index, ab63Entry := range Ab63B50FCandidates().Ab63Entries {
		for _, ab63Axis := range []string{"family", "variant", "primary", "secondary", "secondary_absent", "secondary_zero", "operation", "operation_zero_or_absent", "direction", "direction_empty", "field"} {
			ab63T.Run(ab63fmt.Sprintf("ab63_%d_%s", ab63Index, ab63Axis), func(ab63T *ab63testing.T) {
				ab63KeyValue := ab63Entry.Ab63Key
				switch ab63Axis {
				case "family":
					ab63KeyValue.Ab63Family = ""
				case "variant":
					ab63KeyValue.Ab63Variant = ""
				case "primary":
					ab63KeyValue.Ab63Primary = 0
				case "secondary":
					ab63KeyValue.Ab63Secondary = ab63Byte(16)
				case "secondary_absent":
					ab63KeyValue.Ab63Secondary = nil
				case "secondary_zero":
					ab63KeyValue.Ab63Secondary = ab63Byte(0)
				case "operation":
					ab63KeyValue.Ab63Operation = ab63Byte(255)
				case "operation_zero_or_absent":
					if ab63KeyValue.Ab63Operation == nil {
						ab63KeyValue.Ab63Operation = ab63Byte(0)
					} else {
						ab63KeyValue.Ab63Operation = nil
					}
				case "direction":
					ab63KeyValue.Ab63Direction = ab63String("response")
				case "direction_empty":
					ab63KeyValue.Ab63Direction = ab63String("")
				case "field":
					ab63KeyValue.Ab63Field = "ab63_unknown"
				}
				if _, ab63OK := Ab63B50FLookup(ab63KeyValue); ab63OK {
					ab63T.Fatal("axis wildcarded or conflated")
				}
				if _, ab63OK := Ab63B50FLookup(ab63Entry.Ab63Key); !ab63OK {
					ab63T.Fatal("valid control rejected")
				}
			})
		}
	}
}

func TestAb63Demand(ab63T *ab63testing.T) {
	ab63Entry := ab63Get(ab63T, ab63Byte(3), "demand")
	if !ab63reflect.DeepEqual(ab63Entry.Ab63Masks, map[string]byte{"burner": 1, "reformer-gas-valve": 2, "gas-flow": 4, "fan-power": 8, "valve": 16}) || ab63Entry.Ab63Reserved == nil || *ab63Entry.Ab63Reserved != 224 || ab63Entry.Ab63Interpretation != "demand bits; not applied actuation" || !ab63reflect.DeepEqual(ab63Entry.Ab63Exclusions, []string{"applied actuation"}) {
		ab63T.Fatal("demand metadata")
	}
	for ab63Index, ab63Raw := range []byte{0, 1, 31, 224, 255} {
		ab63Want := [][]bool{{false, false, false, false, false}, {true, false, false, false, false}, {true, true, true, true, true}, {false, false, false, false, false}, {true, true, true, true, true}}[ab63Index]
		for ab63Bit, ab63Name := range []string{"burner", "reformer-gas-valve", "gas-flow", "fan-power", "valve"} {
			if (ab63Raw&ab63Entry.Ab63Masks[ab63Name] != 0) != ab63Want[ab63Bit] {
				ab63T.Fatal("demand example")
			}
		}
		if ab63Raw&*ab63Entry.Ab63Reserved != []byte{0, 0, 0, 224, 224}[ab63Index] {
			ab63T.Fatal("reserved bits")
		}
	}
}

func TestAb63TemperatureAndConflict(ab63T *ab63testing.T) {
	for ab63Op, ab63Count := range map[byte]int{0: 3, 1: 2} {
		for ab63Number := 1; ab63Number <= ab63Count; ab63Number++ {
			ab63Name := ab63fmt.Sprintf("temperature_%d", ab63Number)
			ab63Entry := ab63Get(ab63T, ab63Byte(ab63Op), ab63Name)
			ab63Error := ab63Get(ab63T, ab63Byte(ab63Op), ab63Name+"_error")
			if ab63Entry.Ab63Field.Ab05Position != nil || ab63Entry.Ab63Field.Ab05Width == nil || *ab63Entry.Ab63Field.Ab05Width != 2 || ab63Entry.Ab63Field.Ab05Codec != "signed little-endian int16" || ab63Entry.Ab63Field.Ab05Unit != "degrees Celsius" || ab63Entry.Ab63Divisor == nil || *ab63Entry.Ab63Divisor != 16 {
				ab63T.Fatal("temperature metadata")
			}
			if ab63Error.Ab63Field.Ab05Position != nil || ab63Error.Ab63Field.Ab05Width == nil || *ab63Error.Ab63Field.Ab05Width != 1 || ab63Error.Ab63Field.Ab05Codec != "raw" || ab63Error.Ab63Divisor != nil || ab63Error.Ab63Interpretation != "raw error byte immediately following "+ab63Name || !ab63reflect.DeepEqual(ab63Error.Ab63Exclusions, []string{"temperature conversion"}) {
				ab63T.Fatal("error byte independence")
			}
			for ab63Index, ab63Pair := range [][2]byte{{240, 255}, {0, 0}, {16, 0}} {
				for _, ab63ErrorRaw := range []byte{0, 1, 255} {
					ab63Raw := [3]byte{ab63Pair[0], ab63Pair[1], ab63ErrorRaw}
					ab63Value := float64(int16(uint16(ab63Raw[0])|uint16(ab63Raw[1])<<8)) / float64(*ab63Entry.Ab63Divisor)
					if ab63Value != []float64{-1, 0, 1}[ab63Index] || ab63Raw[2] != ab63ErrorRaw {
						ab63T.Fatal("signed sample or adjacent error")
					}
				}
			}
		}
	}
	ab63Stack := ab63Get(ab63T, ab63Byte(1), "stack-voltage")
	ab63Raw := [2]byte{16, 0}
	if ab63Stack.Ab63Field.Ab05Codec != "raw" || ab63Stack.Ab63Field.Ab05Unit != "" || ab63Stack.Ab63Divisor != nil || ab63Stack.Ab63Field.Ab05Position != nil || !ab63reflect.DeepEqual(ab63Stack.Ab63Conflict, []string{"stack-voltage label", "temperature conversion metadata"}) || !ab63reflect.DeepEqual(ab63Stack.Ab63Exclusions, []string{"voltage scale", "selected temperature conversion"}) || ab63Raw != [2]byte{16, 0} {
		ab63T.Fatal("stack conversion conflict resolved or raw changed")
	}
}

func TestAb63IntegrityAndFeedback(ab63T *ab63testing.T) {
	ab63Feedback := ab63Get(ab63T, nil, "feedback")
	if ab63Feedback.Ab63Interpretation != "acceptance feedback only" || !ab63reflect.DeepEqual(ab63Feedback.Ab63Exclusions, []string{"integrity validation", "physical effect"}) || len(ab63Feedback.Ab63Field.Ab05Enum) != 256 {
		ab63T.Fatal("feedback boundaries")
	}
	for ab63Value := 0; ab63Value < 256; ab63Value++ {
		ab63Want := "rejected"
		if ab63Value == 1 {
			ab63Want = "accepted"
		}
		if ab63Feedback.Ab63Field.Ab05Enum[int64(ab63Value)] != ab63Want {
			ab63T.Fatal("feedback domain")
		}
	}
	for ab63Op, ab63Positions := range map[byte][2]int{0: {12, 13}, 1: {12, 13}, 2: {11, 12}, 3: {4, 5}} {
		ab63First := ab63Get(ab63T, ab63Byte(ab63Op), "integrity_1")
		ab63Second := ab63Get(ab63T, ab63Byte(ab63Op), "integrity_2")
		for ab63Index, ab63Entry := range []Ab63Entry{ab63First, ab63Second} {
			if ab63Entry.Ab63Field.Ab05Position == nil || *ab63Entry.Ab63Field.Ab05Position != ab63Positions[ab63Index] || ab63Entry.Ab63Field.Ab05Width == nil || *ab63Entry.Ab63Field.Ab05Width != 1 || ab63Entry.Ab63Field.Ab05Codec != "raw" || ab63Entry.Ab63Interpretation != "independent raw integrity byte" || ab63Entry.Ab63Coordinates != "abstract logical byte position; not established wire offset" || !ab63reflect.DeepEqual(ab63Entry.Ab63Exclusions, []string{"word byte order: Unknown", "polynomial: Unknown", "initialization: Unknown", "reflection: Unknown", "final transform: Unknown"}) {
				ab63T.Fatal("integrity interpretation")
			}
		}
		for _, ab63Left := range []byte{0, 18, 255} {
			for _, ab63Right := range []byte{0, 52, 255} {
				var ab63Raw [14]byte
				ab63Raw[ab63Positions[0]] = ab63Left
				ab63Raw[ab63Positions[1]] = ab63Right
				if [2]byte{ab63Raw[*ab63First.Ab63Field.Ab05Position], ab63Raw[*ab63Second.Ab63Field.Ab05Position]} != [2]byte{ab63Left, ab63Right} {
					ab63T.Fatal("ordered independent bytes")
				}
				for ab63Index, ab63Value := range []int64{0, 1, 2, 255} {
					if ab63Feedback.Ab63Field.Ab05Enum[ab63Value] != []string{"rejected", "accepted", "rejected", "rejected"}[ab63Index] {
						ab63T.Fatal("feedback coupled to integrity")
					}
				}
			}
		}
	}
}

func TestAb63DeepCopy(ab63T *ab63testing.T) {
	ab63Before := Ab63B50FCandidates()
	ab63Mutated := Ab63B50FCandidates()
	ab63Mutated.Ab63Limitations[0] = "ab63_mutated"
	for ab63Index := range ab63Mutated.Ab63Entries {
		ab63Entry := &ab63Mutated.Ab63Entries[ab63Index]
		if ab63Entry.Ab63Key.Ab63Secondary != nil {
			*ab63Entry.Ab63Key.Ab63Secondary = 0
		}
		if ab63Entry.Ab63Key.Ab63Operation != nil {
			*ab63Entry.Ab63Key.Ab63Operation = 255
		}
		if ab63Entry.Ab63Field.Ab05Position != nil {
			*ab63Entry.Ab63Field.Ab05Position = 99
		}
		if ab63Entry.Ab63Field.Ab05Width != nil {
			*ab63Entry.Ab63Field.Ab05Width = 99
		}
		if ab63Entry.Ab63Divisor != nil {
			*ab63Entry.Ab63Divisor = 99
		}
		if ab63Entry.Ab63Reserved != nil {
			*ab63Entry.Ab63Reserved = 0
		}
		for ab63Name := range ab63Entry.Ab63Masks {
			ab63Entry.Ab63Masks[ab63Name] = 0
		}
		for ab63Value := range ab63Entry.Ab63Field.Ab05Enum {
			ab63Entry.Ab63Field.Ab05Enum[ab63Value] = "ab63_mutated"
		}
		if len(ab63Entry.Ab63Conflict) > 0 {
			ab63Entry.Ab63Conflict[0] = "ab63_mutated"
		}
		if len(ab63Entry.Ab63Exclusions) > 0 {
			ab63Entry.Ab63Exclusions[0] = "ab63_mutated"
		}
	}
	if !ab63reflect.DeepEqual(ab63Before, Ab63B50FCandidates()) {
		ab63T.Fatal("accessor aliases nested state")
	}
	for _, ab63Entry := range ab63Before.Ab63Entries {
		ab63Returned := ab63Get(ab63T, ab63Entry.Ab63Key.Ab63Operation, ab63Entry.Ab63Key.Ab63Field)
		*ab63Returned.Ab63Field.Ab05Width = 99
		if ab63Again, ab63OK := Ab63B50FLookup(ab63Entry.Ab63Key); !ab63OK || !ab63reflect.DeepEqual(ab63Again, ab63Entry) {
			ab63T.Fatal("lookup aliases state")
		}
	}
}
