package vaillant

import (
	ab20bytes "bytes"
	ab20json "encoding/json"
	ab20reflect "reflect"
	ab20testing "testing"
)

func ab20KeyFor(ab20Variant, ab20Direction, ab20Field string) Ab20Key {
	ab20Secondary := uint8(22)
	return Ab20Key{Ab20Family: "vaillant", Ab20Variant: ab20Variant, Ab20Primary: 181, Ab20Secondary: &ab20Secondary, Ab20Direction: ab20Direction, Ab20Field: ab20Field, Ab20Alternative: "none", Ab20InstanceScope: "none"}
}

func ab20RowFor(ab20T *ab20testing.T, ab20Key Ab20Key) map[string]ab20json.RawMessage {
	ab20T.Helper()
	ab20Raw, ab20Found, ab20Err := Ab20Lookup(ab20Key)
	if ab20Err != nil || !ab20Found {
		ab20T.Fatalf("ab20 lookup: %v %v", ab20Found, ab20Err)
	}
	var ab20Row map[string]ab20json.RawMessage
	if ab20Err := ab20json.Unmarshal(ab20Raw, &ab20Row); ab20Err != nil {
		ab20T.Fatal(ab20Err)
	}
	return ab20Row
}

func ab20Want(ab20T *ab20testing.T, ab20Row map[string]ab20json.RawMessage, ab20Field, ab20Expected string) {
	ab20T.Helper()
	var ab20Got, ab20Want any
	if ab20Err := ab20json.Unmarshal(ab20Row[ab20Field], &ab20Got); ab20Err != nil {
		ab20T.Fatalf("ab20 %s: %v", ab20Field, ab20Err)
	}
	if ab20Err := ab20json.Unmarshal([]byte(ab20Expected), &ab20Want); ab20Err != nil {
		ab20T.Fatal(ab20Err)
	}
	if !ab20reflect.DeepEqual(ab20Got, ab20Want) {
		ab20T.Fatalf("ab20 %s = %v, want %v", ab20Field, ab20Got, ab20Want)
	}
}

func TestAb20DirectionalRepresentations(ab20T *ab20testing.T) {
	ab20Request := ab20RowFor(ab20T, ab20KeyFor("energy-type", "request", "energy-type"))
	ab20Read := ab20RowFor(ab20T, ab20KeyFor("energy-type", "read-positive", "energy-type"))
	ab20Write := ab20RowFor(ab20T, ab20KeyFor("energy-type", "write-positive", "energy-type"))
	ab20Want(ab20T, ab20Request, "representation", `"raw-u8-enum"`)
	ab20Want(ab20T, ab20Request, "enum_labels", `"Unknown"`)
	ab20Want(ab20T, ab20Read, "representation", `"low-six-and-independent-flags"`)
	ab20Want(ab20T, ab20Read, "flags", `[{"bit":6,"mask":64},{"bit":7,"mask":128}]`)
	var ab20Flags []struct {
		Ab20Bit  uint8 `json:"bit"`
		Ab20Mask uint8 `json:"mask"`
	}
	if ab20Err := ab20json.Unmarshal(ab20Read["flags"], &ab20Flags); ab20Err != nil {
		ab20T.Fatal(ab20Err)
	}
	for _, ab20Flag := range ab20Flags {
		if uint8(0xC1)&ab20Flag.Ab20Mask != ab20Flag.Ab20Mask {
			ab20T.Fatalf("ab20 missing high flag %d", ab20Flag.Ab20Bit)
		}
	}
	ab20Want(ab20T, ab20Write, "representation", `"scalar-six-bit"`)
	for _, ab20Row := range []map[string]ab20json.RawMessage{ab20Request, ab20Write} {
		if _, ab20Exists := ab20Row["flags"]; ab20Exists {
			ab20T.Fatal("ab20 flags assigned outside read-positive")
		}
	}
	if _, ab20Exists := ab20Request["value_mask"]; ab20Exists {
		ab20T.Fatal("ab20 request flattened into scalar")
	}
	for _, ab20Case := range []struct {
		ab20Name     string
		ab20Row      map[string]ab20json.RawMessage
		ab20Masks    map[string]uint8
		ab20Expected map[string]uint8
	}{
		{"ab20-request", ab20Request, map[string]uint8{"raw_mask": 255}, map[string]uint8{"raw_mask": 193}},
		{"ab20-read", ab20Read, map[string]uint8{"raw_mask": 255, "value_mask": 63}, map[string]uint8{"raw_mask": 193, "value_mask": 1}},
		{"ab20-write", ab20Write, map[string]uint8{"raw_mask": 255, "value_mask": 63, "unmodeled_mask": 192}, map[string]uint8{"raw_mask": 193, "value_mask": 1, "unmodeled_mask": 192}},
	} {
		ab20T.Run(ab20Case.ab20Name, func(ab20T *ab20testing.T) {
			for ab20Field, ab20ExpectedMask := range ab20Case.ab20Masks {
				var ab20Mask uint8
				if ab20Err := ab20json.Unmarshal(ab20Case.ab20Row[ab20Field], &ab20Mask); ab20Err != nil {
					ab20T.Fatal(ab20Err)
				}
				if ab20Mask != ab20ExpectedMask || uint8(0xC1)&ab20Mask != ab20Case.ab20Expected[ab20Field] {
					ab20T.Fatalf("ab20 %s: mask %d", ab20Field, ab20Mask)
				}
			}
		})
	}
}

func TestAb20TimeAndEvidenceBoundaries(ab20T *ab20testing.T) {
	for _, ab20Case := range []struct{ ab20Variant, ab20State string }{{"generic", "header-only"}, {"other", "absent"}} {
		ab20Row := ab20RowFor(ab20T, ab20KeyFor(ab20Case.ab20Variant, "positive", "system-time"))
		ab20Want(ab20T, ab20Row, "positive_state", `"`+ab20Case.ab20State+`"`)
	}
	ab20Data, ab20Err := Ab20Entries()
	if ab20Err != nil {
		ab20T.Fatal(ab20Err)
	}
	var ab20Rows []map[string]ab20json.RawMessage
	if ab20Err := ab20json.Unmarshal(ab20Data, &ab20Rows); ab20Err != nil {
		ab20T.Fatal(ab20Err)
	}
	if len(ab20Rows) != 5 {
		ab20T.Fatalf("ab20 finite area: %d", len(ab20Rows))
	}
	for _, ab20Row := range ab20Rows {
		var ab20Identity Ab20Key
		if ab20Err := ab20json.Unmarshal(ab20Row["key"], &ab20Identity); ab20Err != nil {
			ab20T.Fatal(ab20Err)
		}
		ab20LookedUp := ab20RowFor(ab20T, ab20Identity)
		if !ab20reflect.DeepEqual(ab20LookedUp, ab20Row) {
			ab20T.Fatal("ab20 lookup did not return whole entry")
		}
		ab20Want(ab20T, ab20Row, "grade", `"Hypothesis"`)
		ab20Want(ab20T, ab20Row, "metadata_presence", `"Proven"`)
		ab20Want(ab20T, ab20Row, "raw_form", `"retained without normalization"`)
		for _, ab20Field := range []string{"physical_length", "native_support", "physical_effects", "persistence", "delivery", "timeout", "broadcast"} {
			ab20Want(ab20T, ab20Row, ab20Field, `"Unknown"`)
		}
		ab20Want(ab20T, ab20Row, "unknown", `["encoded operation","exact enum tables","missing positions","native qualification"]`)
		var ab20Key map[string]ab20json.RawMessage
		if ab20Err := ab20json.Unmarshal(ab20Row["key"], &ab20Key); ab20Err != nil {
			ab20T.Fatal(ab20Err)
		}
		for _, ab20Field := range []string{"operation", "position", "enum_table"} {
			if _, ab20Exists := ab20Key[ab20Field]; ab20Exists {
				ab20T.Fatalf("ab20 invented key %s", ab20Field)
			}
			if _, ab20Exists := ab20Row[ab20Field]; ab20Exists {
				ab20T.Fatalf("ab20 invented metadata %s", ab20Field)
			}
		}
	}
}

func ab20AxisChanges() map[string]func(*Ab20Key) {
	return map[string]func(*Ab20Key){
		"ab20-family":           func(ab20Key *Ab20Key) { ab20Key.Ab20Family = "synthetic-family" },
		"ab20-variant":          func(ab20Key *Ab20Key) { ab20Key.Ab20Variant = "synthetic-variant" },
		"ab20-primary":          func(ab20Key *Ab20Key) { ab20Key.Ab20Primary = 0 },
		"ab20-secondary":        func(ab20Key *Ab20Key) { ab20Value := uint8(0); ab20Key.Ab20Secondary = &ab20Value },
		"ab20-operation":        func(ab20Key *Ab20Key) { ab20Value := uint8(0); ab20Key.Ab20Operation = &ab20Value },
		"ab20-direction":        func(ab20Key *Ab20Key) { ab20Key.Ab20Direction = "synthetic-direction" },
		"ab20-field":            func(ab20Key *Ab20Key) { ab20Key.Ab20Field = "synthetic-field" },
		"ab20-alternative":      func(ab20Key *Ab20Key) { ab20Key.Ab20Alternative = "synthetic-alternative" },
		"ab20-group":            func(ab20Key *Ab20Key) { ab20Value := uint8(0); ab20Key.Ab20Group = &ab20Value },
		"ab20-instance-scope":   func(ab20Key *Ab20Key) { ab20Key.Ab20InstanceScope = "dynamic" },
		"ab20-register":         func(ab20Key *Ab20Key) { ab20Value := uint16(0); ab20Key.Ab20Register = &ab20Value },
		"ab20-secondary-absent": func(ab20Key *Ab20Key) { ab20Key.Ab20Secondary = nil },
	}
}

func TestAb20ExactAxes(ab20T *ab20testing.T) {
	ab20Base := ab20KeyFor("energy-type", "read-positive", "energy-type")
	ab20Original, ab20Found, ab20Err := Ab20Lookup(ab20Base)
	if ab20Err != nil || !ab20Found {
		ab20T.Fatalf("ab20 base: %v %v", ab20Found, ab20Err)
	}
	for ab20Name, ab20Change := range ab20AxisChanges() {
		ab20T.Run(ab20Name, func(ab20T *ab20testing.T) {
			ab20Changed := ab20Base
			ab20Change(&ab20Changed)
			ab20Raw, ab20Found, ab20Err := Ab20Lookup(ab20Changed)
			if ab20Err != nil || ab20Found || ab20Raw != nil {
				ab20T.Fatalf("ab20 fallback: %v %v", ab20Found, ab20Err)
			}
			ab20Fixture, ab20Err := ab20json.Marshal([]any{map[string]any{"id": "ab20-control-original", "key": ab20Base}, map[string]any{"id": "ab20-control-changed", "key": ab20Changed}})
			if ab20Err != nil {
				ab20T.Fatal(ab20Err)
			}
			for _, ab20Key := range []Ab20Key{ab20Base, ab20Changed} {
				ab20Raw, ab20Found, ab20Err := ab20Lookup(ab20Fixture, ab20Key)
				if ab20Err != nil || !ab20Found {
					ab20T.Fatalf("ab20 independent control: %v %v", ab20Found, ab20Err)
				}
				var ab20Entry struct {
					Ab20Key Ab20Key `json:"key"`
				}
				if ab20Err := ab20json.Unmarshal(ab20Raw, &ab20Entry); ab20Err != nil {
					ab20T.Fatal(ab20Err)
				}
				if !ab20reflect.DeepEqual(ab20Entry.Ab20Key, ab20Key) {
					ab20T.Fatal("ab20 wrong whole row")
				}
			}
		})
	}
	ab20Again, _, _ := Ab20Lookup(ab20Base)
	if !ab20bytes.Equal(ab20Original, ab20Again) {
		ab20T.Fatal("ab20 key mutation changed catalog")
	}
}

func TestAb20ValidationAndCopies(ab20T *ab20testing.T) {
	ab20Key := ab20KeyFor("generic", "positive", "system-time")
	for _, ab20Invalid := range []string{`{`, `null`, `[]`, `{}`, `[null]`, `[{}]`, `[] []`, `[{"id":"ab20-bad","key":{"primary":256}}]`} {
		if _, ab20Err := ab20Validate([]byte(ab20Invalid)); ab20Err == nil {
			ab20T.Fatalf("ab20 accepted %s", ab20Invalid)
		}
		if _, ab20Found, ab20Err := ab20Lookup([]byte(ab20Invalid), ab20Key); ab20Err == nil || ab20Found {
			ab20T.Fatal("ab20 lookup did not fail closed")
		}
	}
	ab20Duplicate, ab20Err := ab20json.Marshal([]any{map[string]any{"id": "ab20-first", "key": ab20Key}, map[string]any{"id": "ab20-second", "key": ab20Key}})
	if ab20Err != nil {
		ab20T.Fatal(ab20Err)
	}
	if _, ab20Err := ab20Validate(ab20Duplicate); ab20Err == nil {
		ab20T.Fatal("ab20 duplicate key admitted")
	}
	if _, ab20Found, ab20Err := ab20Lookup(ab20Duplicate, ab20Key); ab20Err == nil || ab20Found {
		ab20T.Fatal("ab20 duplicate lookup admitted")
	}
	ab20Entries, ab20Err := Ab20Entries()
	if ab20Err != nil {
		ab20T.Fatal(ab20Err)
	}
	ab20Row, ab20Found, ab20Err := Ab20Lookup(ab20Key)
	if ab20Err != nil || !ab20Found {
		ab20T.Fatalf("ab20 lookup: %v %v", ab20Found, ab20Err)
	}
	ab20ExpectedEntries := append([]byte(nil), ab20Entries...)
	ab20ExpectedRow := append([]byte(nil), ab20Row...)
	ab20SecondEntries, _ := Ab20Entries()
	ab20SecondRow, _, _ := Ab20Lookup(ab20Key)
	for ab20Index := range ab20Entries {
		ab20Entries[ab20Index] = 'x'
	}
	for ab20Index := range ab20Row {
		ab20Row[ab20Index] = 'x'
	}
	ab20LaterEntries, ab20Err := Ab20Entries()
	if ab20Err != nil {
		ab20T.Fatal(ab20Err)
	}
	ab20LaterRow, ab20Found, ab20Err := Ab20Lookup(ab20Key)
	if ab20Err != nil || !ab20Found {
		ab20T.Fatalf("ab20 later: %v %v", ab20Found, ab20Err)
	}
	if !ab20bytes.Equal(ab20ExpectedEntries, ab20SecondEntries) || !ab20bytes.Equal(ab20ExpectedEntries, ab20LaterEntries) || !ab20bytes.Equal(ab20ExpectedRow, ab20SecondRow) || !ab20bytes.Equal(ab20ExpectedRow, ab20LaterRow) {
		ab20T.Fatal("ab20 shared return memory")
	}
}

func TestAb20OptionalPresence(ab20T *ab20testing.T) {
	for _, ab20Name := range []string{"ab20-secondary", "ab20-operation", "ab20-group", "ab20-register"} {
		ab20T.Run(ab20Name, func(ab20T *ab20testing.T) {
			ab20Absent := ab20KeyFor("generic", "positive", "system-time")
			ab20Absent.Ab20Secondary = nil
			ab20Zero := ab20Absent
			ab20AxisChanges()[ab20Name](&ab20Zero)
			ab20Fixture, ab20Err := ab20json.Marshal([]any{map[string]any{"id": "ab20-absent", "key": ab20Absent}, map[string]any{"id": "ab20-zero", "key": ab20Zero}})
			if ab20Err != nil {
				ab20T.Fatal(ab20Err)
			}
			for ab20Index, ab20Key := range []Ab20Key{ab20Absent, ab20Zero} {
				ab20Raw, ab20Found, ab20Err := ab20Lookup(ab20Fixture, ab20Key)
				if ab20Err != nil || !ab20Found {
					ab20T.Fatalf("ab20 optional presence: %v %v", ab20Found, ab20Err)
				}
				var ab20Entry struct {
					Ab20ID string `json:"id"`
				}
				if ab20Err := ab20json.Unmarshal(ab20Raw, &ab20Entry); ab20Err != nil {
					ab20T.Fatal(ab20Err)
				}
				if ab20Entry.Ab20ID != []string{"ab20-absent", "ab20-zero"}[ab20Index] {
					ab20T.Fatal("ab20 absent collapsed with zero")
				}
			}
		})
	}
}

func TestAb20IdentityValidation(ab20T *ab20testing.T) {
	ab20Encoded, ab20Err := ab20json.Marshal(ab20KeyFor("generic", "positive", "system-time"))
	if ab20Err != nil {
		ab20T.Fatal(ab20Err)
	}
	for _, ab20Field := range []string{"family", "variant", "primary", "direction", "field", "alternative", "instance_scope"} {
		ab20T.Run("ab20-"+ab20Field, func(ab20T *ab20testing.T) {
			for _, ab20Null := range []bool{false, true} {
				var ab20Identity map[string]ab20json.RawMessage
				if ab20Err := ab20json.Unmarshal(ab20Encoded, &ab20Identity); ab20Err != nil {
					ab20T.Fatal(ab20Err)
				}
				delete(ab20Identity, ab20Field)
				if ab20Null {
					ab20Identity[ab20Field] = ab20json.RawMessage(`null`)
				}
				ab20Fixture, ab20Err := ab20json.Marshal([]any{map[string]any{"id": "ab20-incomplete", "key": ab20Identity}})
				if ab20Err != nil {
					ab20T.Fatal(ab20Err)
				}
				if _, ab20Err := ab20Validate(ab20Fixture); ab20Err == nil {
					ab20T.Fatal("ab20 incomplete identity admitted")
				}
			}
		})
	}
}
