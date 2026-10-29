package vaillant

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

func ab05Fixture(ab05T *testing.T) []Ab05Record {
	ab05T.Helper()
	ab05Data, ab05Error := os.ReadFile("testdata/ab05_synthetic_metadata.json")
	if ab05Error != nil {
		ab05T.Fatal(ab05Error)
	}
	var ab05Records []Ab05Record
	if ab05Error = json.Unmarshal(ab05Data, &ab05Records); ab05Error != nil {
		ab05T.Fatal(ab05Error)
	}
	return ab05Records
}

func ab05Int(ab05Value int) *int { return &ab05Value }

func TestAb05KeyDimensions(ab05T *testing.T) {
	ab05Cases := []struct {
		ab05Name   string
		ab05Change func(*Ab05Key)
	}{
		{"ab05_family", func(ab05Key *Ab05Key) { ab05Key.Ab05Family = "beta" }},
		{"ab05_variant", func(ab05Key *Ab05Key) { ab05Key.Ab05Variant = "ab05_secondary" }},
		{"ab05_direction", func(ab05Key *Ab05Key) { ab05Key.Ab05Direction = "request" }},
		{"ab05_opcode", func(ab05Key *Ab05Key) { ab05Key.Ab05Opcode = 6 }},
		{"ab05_group", func(ab05Key *Ab05Key) { ab05Key.Ab05Group = 4 }},
		{"ab05_instance", func(ab05Key *Ab05Key) { ab05Key.Ab05Instance = 2 }},
		{"ab05_register", func(ab05Key *Ab05Key) { ab05Key.Ab05Register = 35 }},
		{"ab05_alternative", func(ab05Key *Ab05Key) { ab05Key.Ab05Alternative = "ab05_other" }},
	}
	for _, ab05Case := range ab05Cases {
		ab05T.Run(ab05Case.ab05Name, func(ab05T *testing.T) {
			ab05First := ab05Fixture(ab05T)[0]
			ab05Second := ab05First
			ab05Case.ab05Change(&ab05Second.Ab05Key)
			ab05Second.Ab05RawFallback = []byte{42}
			ab05Records := []Ab05Record{ab05First, ab05Second}
			if ab05Error := Ab05Validate(ab05Records); ab05Error != nil {
				ab05T.Fatal(ab05Error)
			}
			for ab05Index, ab05Want := range ab05Records {
				ab05Got, ab05Found, ab05Error := Ab05Query(ab05Records, ab05Want.Ab05Key)
				if ab05Error != nil || !ab05Found || !reflect.DeepEqual(ab05Got, ab05Want) {
					ab05T.Fatalf("ab05 record %d: got %#v, found %v, error %v", ab05Index, ab05Got, ab05Found, ab05Error)
				}
			}
			ab05Second.Ab05Key = ab05First.Ab05Key
			if Ab05Validate([]Ab05Record{ab05First, ab05Second}) == nil {
				ab05T.Fatal("ab05 duplicate accepted")
			}
		})
	}
}

func TestAb05AlternativesAndOverlap(ab05T *testing.T) {
	ab05Records := ab05Fixture(ab05T)
	for _, ab05Index := range []int{0, 1} {
		ab05Got, ab05Found, ab05Error := Ab05Query(ab05Records, ab05Records[ab05Index].Ab05Key)
		if ab05Error != nil || !ab05Found || !reflect.DeepEqual(ab05Got, ab05Records[ab05Index]) {
			ab05T.Fatalf("ab05 alternative flattened or overlap lost: %#v %v", ab05Got, ab05Error)
		}
	}
	if len(ab05Records[0].Ab05Fields) != 2 || *ab05Records[0].Ab05Fields[0].Ab05Width != 2 || *ab05Records[1].Ab05Fields[0].Ab05Width != 1 {
		ab05T.Fatal("ab05 fixture lost distinct overlapping candidates")
	}
}

func TestAb05AbsentExtents(ab05T *testing.T) {
	ab05Records := ab05Fixture(ab05T)
	ab05Header, ab05Found, ab05Error := Ab05Query(ab05Records, ab05Records[2].Ab05Key)
	if ab05Error != nil || !ab05Found || ab05Header.Ab05BodyWidth != nil || len(ab05Header.Ab05Fields) != 0 || !reflect.DeepEqual(ab05Header.Ab05RawFallback, []byte{171, 205}) {
		ab05T.Fatalf("ab05 header ceased to be partial: %#v %v", ab05Header, ab05Error)
	}
	ab05Opaque, ab05Found, ab05Error := Ab05Query(ab05Records, ab05Records[3].Ab05Key)
	if ab05Error != nil || !ab05Found || ab05Opaque.Ab05Fields[0].Ab05Position != nil || ab05Opaque.Ab05Fields[0].Ab05Width != nil {
		ab05T.Fatalf("ab05 unknown extent became explicit: %#v %v", ab05Opaque, ab05Error)
	}
	ab05Records[2].Ab05BodyWidth = ab05Int(0)
	ab05Zero, _, ab05Error := Ab05Query(ab05Records, ab05Records[2].Ab05Key)
	if ab05Error != nil || ab05Zero.Ab05BodyWidth == nil || *ab05Zero.Ab05BodyWidth != 0 {
		ab05T.Fatal("ab05 explicit empty body lost")
	}
}

func TestAb05UnknownGradeAndCodec(ab05T *testing.T) {
	ab05Records := ab05Fixture(ab05T)
	ab05Header, _, ab05Error := Ab05Query(ab05Records, ab05Records[2].Ab05Key)
	if ab05Error != nil || ab05Header.Ab05Grade != Ab05Unknown {
		ab05T.Fatal("ab05 missing grade promoted")
	}
	ab05Opaque, _, ab05Error := Ab05Query(ab05Records, ab05Records[3].Ab05Key)
	if ab05Error != nil || ab05Opaque.Ab05Grade != Ab05Unknown || ab05Opaque.Ab05Fields[0].Ab05Grade != Ab05Unknown || ab05Opaque.Ab05Fields[0].Ab05Codec != "ab05_unrecognized" || !reflect.DeepEqual(ab05Opaque.Ab05Fields[0].Ab05RawFallback, []byte{171, 205}) || !reflect.DeepEqual(ab05Opaque.Ab05RawFallback, []byte{171, 205}) {
		ab05T.Fatalf("ab05 opaque evidence changed: %#v %v", ab05Opaque, ab05Error)
	}
	if ab05Records[2].Ab05Grade != "" || ab05Records[3].Ab05Fields[0].Ab05Grade != "" {
		ab05T.Fatal("ab05 supplied evidence mutated")
	}
	for _, ab05Grade := range []Ab05Grade{Ab05Unknown, Ab05Hypothesis, Ab05Proven} {
		ab05Records[3].Ab05Grade = ab05Grade
		ab05Records[3].Ab05Fields[0].Ab05Grade = ab05Grade
		ab05Got, _, ab05Error := Ab05Query(ab05Records, ab05Records[3].Ab05Key)
		if ab05Error != nil || ab05Got.Ab05Grade != ab05Grade || ab05Got.Ab05Fields[0].Ab05Grade != ab05Grade || ab05Got.Ab05Fields[0].Ab05Codec != "ab05_unrecognized" {
			ab05T.Fatal("ab05 grade conflated with codec")
		}
	}
}

func TestAb05ContradictoryExtents(ab05T *testing.T) {
	ab05Cases := []struct {
		ab05Name   string
		ab05Change func(*Ab05Record)
	}{
		{"ab05_negative_body", func(ab05Record *Ab05Record) { ab05Record.Ab05BodyWidth = ab05Int(-1) }},
		{"ab05_negative_position", func(ab05Record *Ab05Record) { ab05Record.Ab05Fields[0].Ab05Position = ab05Int(-1) }},
		{"ab05_zero_field_width", func(ab05Record *Ab05Record) { ab05Record.Ab05Fields[0].Ab05Width = ab05Int(0) }},
		{"ab05_negative_field_width", func(ab05Record *Ab05Record) { ab05Record.Ab05Fields[0].Ab05Width = ab05Int(-1) }},
		{"ab05_end_exceeds_body", func(ab05Record *Ab05Record) { ab05Record.Ab05Fields[0].Ab05Position = ab05Int(1) }},
		{"ab05_position_exceeds_body", func(ab05Record *Ab05Record) {
			ab05Record.Ab05Fields[0].Ab05Position = ab05Int(3)
			ab05Record.Ab05Fields[0].Ab05Width = nil
		}},
		{"ab05_position_at_body_end", func(ab05Record *Ab05Record) {
			ab05Record.Ab05Fields[0].Ab05Position = ab05Int(2)
			ab05Record.Ab05Fields[0].Ab05Width = nil
		}},
		{"ab05_width_without_position", func(ab05Record *Ab05Record) {
			ab05Record.Ab05Fields[0].Ab05Position = nil
			ab05Record.Ab05Fields[0].Ab05Width = ab05Int(3)
		}},
		{"ab05_overflow_extent", func(ab05Record *Ab05Record) {
			ab05Record.Ab05Fields[0].Ab05Position = ab05Int(1)
			ab05Record.Ab05Fields[0].Ab05Width = ab05Int(int(^uint(0) >> 1))
		}},
		{"ab05_invalid_record_grade", func(ab05Record *Ab05Record) { ab05Record.Ab05Grade = "ab05_invalid" }},
		{"ab05_invalid_field_grade", func(ab05Record *Ab05Record) { ab05Record.Ab05Fields[0].Ab05Grade = "ab05_invalid" }},
	}
	for _, ab05Case := range ab05Cases {
		ab05T.Run(ab05Case.ab05Name, func(ab05T *testing.T) {
			ab05Record := ab05Fixture(ab05T)[0]
			if Ab05Validate([]Ab05Record{ab05Record}) != nil {
				ab05T.Fatal("ab05 invalid control")
			}
			ab05Case.ab05Change(&ab05Record)
			ab05Got, ab05Found, ab05Error := Ab05Query([]Ab05Record{ab05Record}, ab05Record.Ab05Key)
			if ab05Error == nil || ab05Found || !reflect.DeepEqual(ab05Got, Ab05Record{}) {
				ab05T.Fatal("ab05 invalid evidence admitted")
			}
		})
	}
}

func TestAb05QueryIsolation(ab05T *testing.T) {
	ab05Records := ab05Fixture(ab05T)
	ab05Before := ab05Fixture(ab05T)
	ab05Got, _, ab05Error := Ab05Query(ab05Records, ab05Records[0].Ab05Key)
	if ab05Error != nil {
		ab05T.Fatal(ab05Error)
	}
	*ab05Got.Ab05BodyWidth = 99
	*ab05Got.Ab05Fields[0].Ab05Position = 99
	*ab05Got.Ab05Fields[0].Ab05Width = 99
	ab05Got.Ab05RawFallback[0] = 0
	ab05Got.Ab05Fields[0].Ab05RawFallback[0] = 0
	ab05Got.Ab05Fields[0].Ab05Sentinel[0] = 0
	ab05Got.Ab05Fields[0].Ab05Enum[0] = "ab05_changed"
	ab05Got.Ab05Fields[1].Ab05Name = "ab05_changed"
	if !reflect.DeepEqual(ab05Records, ab05Before) {
		ab05T.Fatal("ab05 query changed prior field evidence")
	}
	ab05Missing := ab05Records[0].Ab05Key
	ab05Missing.Ab05Family = "ab05_absent"
	ab05Got, ab05Found, ab05Error := Ab05Query(ab05Records, ab05Missing)
	if ab05Error != nil || ab05Found || !reflect.DeepEqual(ab05Got, Ab05Record{}) {
		ab05T.Fatal("ab05 missing key admitted")
	}
}

func TestAb05JSONNumericKeyPresence(ab05T *testing.T) {
	ab05Base := map[string]any{
		"ab05_family": "ab05_synthetic", "ab05_variant": "ab05_primary",
		"ab05_direction": "request", "ab05_alternative": "ab05_candidate",
		"ab05_opcode": 0, "ab05_group": 0, "ab05_instance": 0, "ab05_register": 0,
	}
	ab05Decode := func(ab05T *testing.T, ab05Key map[string]any) (Ab05Record, error) {
		ab05T.Helper()
		ab05Data, ab05Error := json.Marshal(map[string]any{"ab05_key": ab05Key})
		if ab05Error != nil {
			ab05T.Fatal(ab05Error)
		}
		var ab05Record Ab05Record
		ab05Error = json.Unmarshal(ab05Data, &ab05Record)
		return ab05Record, ab05Error
	}
	ab05Zero, ab05Error := ab05Decode(ab05T, ab05Base)
	if ab05Error != nil {
		ab05T.Fatal(ab05Error)
	}
	ab05Want := Ab05Key{Ab05Family: "ab05_synthetic", Ab05Variant: "ab05_primary", Ab05Direction: "request", Ab05Alternative: "ab05_candidate"}
	ab05Got, ab05Found, ab05Error := Ab05Query([]Ab05Record{ab05Zero}, ab05Want)
	if ab05Error != nil || !ab05Found || ab05Got.Ab05Key != ab05Want {
		ab05T.Fatalf("ab05 explicit zero selectors rejected: %#v %v %v", ab05Got, ab05Found, ab05Error)
	}
	for _, ab05Name := range []string{"ab05_opcode", "ab05_group", "ab05_instance", "ab05_register"} {
		for _, ab05Mode := range []string{"absent", "null", "misspelled"} {
			ab05T.Run(ab05Name+"/"+ab05Mode, func(ab05T *testing.T) {
				ab05Key := make(map[string]any, len(ab05Base))
				for ab05Name, ab05Value := range ab05Base {
					ab05Key[ab05Name] = ab05Value
				}
				delete(ab05Key, ab05Name)
				switch ab05Mode {
				case "null":
					ab05Key[ab05Name] = nil
				case "misspelled":
					ab05Key[ab05Name+"_typo"] = 0
				}
				if _, ab05Error := ab05Decode(ab05T, ab05Key); ab05Error == nil {
					ab05T.Fatal("ab05 missing numeric selector accepted")
				}
			})
		}
	}
	ab05Base["ab05_unknown"] = 0
	if _, ab05Error := ab05Decode(ab05T, ab05Base); ab05Error == nil {
		ab05T.Fatal("ab05 unknown key property accepted")
	}
}

func TestAb05InvalidStringKeys(ab05T *testing.T) {
	ab05Cases := []struct {
		ab05Name   string
		ab05Change func(*Ab05Key)
	}{
		{"ab05_family", func(ab05Key *Ab05Key) { ab05Key.Ab05Family = "" }},
		{"ab05_variant", func(ab05Key *Ab05Key) { ab05Key.Ab05Variant = "" }},
		{"ab05_direction", func(ab05Key *Ab05Key) { ab05Key.Ab05Direction = "" }},
		{"ab05_alternative", func(ab05Key *Ab05Key) { ab05Key.Ab05Alternative = "" }},
		{"ab05_invalid_direction", func(ab05Key *Ab05Key) { ab05Key.Ab05Direction = "ab05_invalid" }},
	}
	for _, ab05Case := range ab05Cases {
		ab05T.Run(ab05Case.ab05Name, func(ab05T *testing.T) {
			ab05Record := ab05Fixture(ab05T)[0]
			ab05Case.ab05Change(&ab05Record.Ab05Key)
			if _, ab05Found, ab05Error := Ab05Query([]Ab05Record{ab05Record}, ab05Record.Ab05Key); ab05Error == nil || ab05Found {
				ab05T.Fatal("ab05 invalid string dimension accepted")
			}
		})
	}
}

func TestAb05QueryEmptyBytes(ab05T *testing.T) {
	for _, ab05Bytes := range [][]byte{nil, {}} {
		ab05Record := ab05Fixture(ab05T)[0]
		ab05Record.Ab05RawFallback = ab05Bytes
		ab05Record.Ab05Fields[0].Ab05RawFallback = ab05Bytes
		ab05Record.Ab05Fields[0].Ab05Sentinel = ab05Bytes
		ab05Got, ab05Found, ab05Error := Ab05Query([]Ab05Record{ab05Record}, ab05Record.Ab05Key)
		if ab05Error != nil || !ab05Found || !reflect.DeepEqual(ab05Got, ab05Record) {
			ab05T.Fatalf("ab05 empty byte presence lost: %#v %v", ab05Got, ab05Error)
		}
	}
}
