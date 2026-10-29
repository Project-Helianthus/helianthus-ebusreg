package vaillant

import (
	ab31bytes "bytes"
	ab31json "encoding/json"
	ab31reflect "reflect"
	ab31testing "testing"
)

func ab31Decode(ab31T *ab31testing.T, ab31Data []byte) map[string]any {
	ab31T.Helper()
	var ab31Value map[string]any
	if ab31Error := ab31json.Unmarshal(ab31Data, &ab31Value); ab31Error != nil {
		ab31T.Fatal(ab31Error)
	}
	return ab31Value
}

func ab31Expect(ab31T *ab31testing.T, ab31Got any, ab31Expected string) {
	ab31T.Helper()
	var ab31Want any
	if ab31Error := ab31json.Unmarshal([]byte(ab31Expected), &ab31Want); ab31Error != nil {
		ab31T.Fatal(ab31Error)
	}
	if !ab31reflect.DeepEqual(ab31Got, ab31Want) {
		ab31T.Fatalf("ab31 got %#v; want %#v", ab31Got, ab31Want)
	}
}

func TestAb31Catalog(ab31T *ab31testing.T) {
	ab31Catalog := ab31Decode(ab31T, Ab31PumpRoles())
	ab31Expect(ab31T, ab31Catalog["catalog"], `"ab31_pump_roles"`)
	ab31Expect(ab31T, ab31Catalog["representation_grade"], `"Proven"`)
	ab31Expect(ab31T, ab31Catalog["coordinate_system"], `"logical positions include prefix positions; not physical offsets"`)
	ab31Expect(ab31T, ab31Catalog["limitations"], `[
"Candidate interpretations remain Hypothesis; physical support remains Unknown.",
"Boiler-pump correspondence is role-qualified; universal boiler support and appliance exclusivity are Unknown.",
"Candidate broadcast, timer and later register change establish no causal chain, physical support or exclusive use.",
"No runtime registration, wire builder, write admission, automatic selection or native codec promotion.",
"Request, positive response and negative response forms are Unknown; direction is null where unknown."]`)
	ab31Entries := ab31Catalog["entries"].([]any)
	if len(ab31Catalog) != 5 || len(ab31Entries) != 3 {
		ab31T.Fatal("ab31 catalog or alternative count changed")
	}
	for ab31Index, ab31Expected := range []string{
		`{"id":"ab31_b508_op05_pump","family":"B508","variant":"pump role","direction":null,"service":"B508","operation":"OP05","alternative":"pump state with conflicting defaults"}`,
		`{"id":"ab31_b508_op0b_solar","family":"B508","variant":"solar cascading","direction":null,"service":"B508","operation":"OP0B","alternative":"independent tank flags"}`,
		`{"id":"ab31_noise_reduction_broadcast","family":"Unknown","variant":"noise-reduction","direction":"broadcast","service":null,"operation":null,"alternative":"distinct noise-reduction candidate"}`,
	} {
		ab31Entry := ab31Entries[ab31Index].(map[string]any)
		ab31Identity := make(map[string]any)
		for _, ab31Axis := range []string{"id", "family", "variant", "direction", "service", "operation", "alternative"} {
			ab31Value, ab31Present := ab31Entry[ab31Axis]
			if !ab31Present {
				ab31T.Fatalf("ab31 missing axis %s", ab31Axis)
			}
			ab31Identity[ab31Axis] = ab31Value
		}
		ab31Expect(ab31T, ab31Identity, ab31Expected)
		ab31Expect(ab31T, ab31Entry["grade"], `"Hypothesis"`)
		ab31Expect(ab31T, ab31Entry["physical_qualification"], `"Unknown"`)
		ab31Expect(ab31T, ab31Entry["logical_extent"], `null`)
		ab31Expect(ab31T, ab31Entry["raw_fallback"], `"preserve native bytes, including unmapped values and reserved bits"`)
		if len(ab31Entry) != 12 {
			ab31T.Fatal("ab31 row fields changed")
		}
	}
	ab31Expect(ab31T, ab31Entries[0].(map[string]any)["fields"], `[{
"name":"pump state","logical_position":null,"width_bytes":null,"raw_domain":[0,1],"signedness":"Unknown","byte_order":null,"conversion":null,"unit":null,"default":null,"mappings":{"0":"blocked","1":"released"},"reserved_regions":null,"grade":"Hypothesis","raw_fallback":"preserve native bytes, including unmapped values and reserved bits","conflicting_default_descriptions":["blocked","released"]}]`)
	ab31Expect(ab31T, ab31Entries[1].(map[string]any)["fields"], `[
{"name":"cycle-pause","logical_position":3,"width_bytes":null,"raw_domain":null,"signedness":"Unknown","byte_order":null,"conversion":null,"unit":null,"default":null,"mappings":{},"reserved_regions":null,"grade":"Hypothesis","raw_fallback":"preserve native bytes, including unmapped values and reserved bits"},
{"name":"tank-one flag","logical_position":4,"width_bytes":1,"raw_domain":{"min":0,"max":255},"signedness":"unsigned","byte_order":null,"conversion":{"mask":1,"shift":0},"unit":null,"default":null,"mappings":{"0":"flag clear","1":"flag set"},"reserved_regions":[{"bits":[1,7]}],"grade":"Hypothesis","raw_fallback":"preserve native bytes, including unmapped values and reserved bits"},
{"name":"tank-two flag","logical_position":5,"width_bytes":1,"raw_domain":{"min":0,"max":255},"signedness":"unsigned","byte_order":null,"conversion":{"mask":1,"shift":0},"unit":null,"default":null,"mappings":{"0":"flag clear","1":"flag set"},"reserved_regions":[{"bits":[1,7]}],"grade":"Hypothesis","raw_fallback":"preserve native bytes, including unmapped values and reserved bits"}]`)
	ab31Expect(ab31T, ab31Entries[2].(map[string]any)["fields"], `[]`)
}

func TestAb31IdentitySerialization(ab31T *ab31testing.T) {
	for _, ab31Axis := range []string{"family", "variant", "direction", "service", "operation", "alternative"} {
		ab31T.Run("ab31_"+ab31Axis, func(ab31T *ab31testing.T) {
			for _, ab31Original := range ab31Decode(ab31T, Ab31PumpRoles())["entries"].([]any) {
				ab31First := ab31Original.(map[string]any)
				ab31Second := make(map[string]any)
				for ab31Key, ab31Value := range ab31First {
					ab31Second[ab31Key] = ab31Value
				}
				ab31Second[ab31Axis] = "ab31_independent_alternative"
				ab31Data, ab31Error := ab31json.Marshal(map[string]any{"entries": []any{ab31First, ab31Second}})
				if ab31Error != nil {
					ab31T.Fatal(ab31Error)
				}
				ab31RoundTrip := ab31Decode(ab31T, ab31Data)["entries"].([]any)
				if len(ab31RoundTrip) != 2 || !ab31reflect.DeepEqual(ab31RoundTrip[0], ab31First) || !ab31reflect.DeepEqual(ab31RoundTrip[1], ab31Second) || ab31reflect.DeepEqual(ab31RoundTrip[0], ab31RoundTrip[1]) {
					ab31T.Fatal("ab31 independent identity axis lost")
				}
			}
		})
	}
}

func TestAb31IndependentTankFlags(ab31T *ab31testing.T) {
	ab31Fields := ab31Decode(ab31T, Ab31PumpRoles())["entries"].([]any)[1].(map[string]any)["fields"].([]any)
	for _, ab31Case := range []struct {
		ab31Name     string
		ab31Raw      [6]byte
		ab31Flags    [2]byte
		ab31Reserved [2]byte
	}{
		{"ab31_one", [6]byte{0, 0, 0, 0, 1, 0}, [2]byte{1, 0}, [2]byte{}},
		{"ab31_two", [6]byte{0, 0, 0, 0, 0, 1}, [2]byte{0, 1}, [2]byte{}},
		{"ab31_reserved", [6]byte{0, 0, 0, 0, 128, 0}, [2]byte{}, [2]byte{128, 0}},
	} {
		ab31T.Run(ab31Case.ab31Name, func(ab31T *ab31testing.T) {
			ab31Before := ab31Case.ab31Raw
			for ab31Index, ab31Value := range ab31Fields[1:] {
				ab31Field := ab31Value.(map[string]any)
				ab31Position := int(ab31Field["logical_position"].(float64))
				ab31Conversion := ab31Field["conversion"].(map[string]any)
				ab31Mask := byte(ab31Conversion["mask"].(float64))
				ab31Shift := uint(ab31Conversion["shift"].(float64))
				ab31Raw := ab31Case.ab31Raw[ab31Position]
				if (ab31Raw&ab31Mask)>>ab31Shift != ab31Case.ab31Flags[ab31Index] || ab31Raw & ^ab31Mask != ab31Case.ab31Reserved[ab31Index] {
					ab31T.Fatal("ab31 flags coupled or reserved evidence changed")
				}
			}
			if ab31Case.ab31Raw != ab31Before {
				ab31T.Fatal("ab31 raw evidence changed")
			}
		})
	}
}

func TestAb31Copy(ab31T *ab31testing.T) {
	ab31First := Ab31PumpRoles()
	ab31Before := append([]byte(nil), ab31First...)
	for ab31Index := range ab31First {
		ab31First[ab31Index] = '!'
	}
	if !ab31bytes.Equal(Ab31PumpRoles(), ab31Before) {
		ab31T.Fatal("ab31 returned bytes alias embedded catalog")
	}
}
