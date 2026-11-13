package vaillant

import (
	pumpRoleBytes "bytes"
	pumpRoleJSON "encoding/json"
	pumpRoleReflect "reflect"
	pumpRoleTesting "testing"
)

func pumpRoleDecode(pumpRoleT *pumpRoleTesting.T, pumpRoleData []byte) map[string]any {
	pumpRoleT.Helper()
	var pumpRoleValue map[string]any
	if pumpRoleError := pumpRoleJSON.Unmarshal(pumpRoleData, &pumpRoleValue); pumpRoleError != nil {
		pumpRoleT.Fatal(pumpRoleError)
	}
	return pumpRoleValue
}

func pumpRoleExpect(pumpRoleT *pumpRoleTesting.T, pumpRoleGot any, pumpRoleExpected string) {
	pumpRoleT.Helper()
	var pumpRoleWant any
	if pumpRoleError := pumpRoleJSON.Unmarshal([]byte(pumpRoleExpected), &pumpRoleWant); pumpRoleError != nil {
		pumpRoleT.Fatal(pumpRoleError)
	}
	if !pumpRoleReflect.DeepEqual(pumpRoleGot, pumpRoleWant) {
		pumpRoleT.Fatalf("ab31 got %#v; want %#v", pumpRoleGot, pumpRoleWant)
	}
}

func TestAb31Catalog(pumpRoleT *pumpRoleTesting.T) {
	pumpRoleCatalog := pumpRoleDecode(pumpRoleT, PumpRoles())
	pumpRoleExpect(pumpRoleT, pumpRoleCatalog["catalog"], `"ab31_pump_roles"`)
	pumpRoleExpect(pumpRoleT, pumpRoleCatalog["representation_grade"], `"Proven"`)
	pumpRoleExpect(pumpRoleT, pumpRoleCatalog["coordinate_system"], `"logical positions include prefix positions; not physical offsets"`)
	pumpRoleExpect(pumpRoleT, pumpRoleCatalog["limitations"], `[
"Candidate interpretations remain Hypothesis; physical support remains Unknown.",
"Boiler-pump correspondence is role-qualified; universal boiler support and appliance exclusivity are Unknown.",
"Candidate broadcast, timer and later register change establish no causal chain, physical support or exclusive use.",
"No runtime registration, wire builder, write admission, automatic selection or native codec promotion.",
"Request, positive response and negative response forms are Unknown; direction is null where unknown."]`)
	pumpRoleEntries := pumpRoleCatalog["entries"].([]any)
	if len(pumpRoleCatalog) != 5 || len(pumpRoleEntries) != 3 {
		pumpRoleT.Fatal("ab31 catalog or alternative count changed")
	}
	for pumpRoleIndex, pumpRoleExpected := range []string{
		`{"id":"ab31_b508_op05_pump","family":"B508","variant":"pump role","direction":null,"service":"B508","operation":"OP05","alternative":"pump state with conflicting defaults"}`,
		`{"id":"ab31_b508_op0b_solar","family":"B508","variant":"solar cascading","direction":null,"service":"B508","operation":"OP0B","alternative":"independent tank flags"}`,
		`{"id":"ab31_noise_reduction_broadcast","family":"Unknown","variant":"noise-reduction","direction":"broadcast","service":null,"operation":null,"alternative":"distinct noise-reduction candidate"}`,
	} {
		pumpRoleEntry := pumpRoleEntries[pumpRoleIndex].(map[string]any)
		pumpRoleIdentity := make(map[string]any)
		for _, pumpRoleAxis := range []string{"id", "family", "variant", "direction", "service", "operation", "alternative"} {
			pumpRoleValue, pumpRolePresent := pumpRoleEntry[pumpRoleAxis]
			if !pumpRolePresent {
				pumpRoleT.Fatalf("ab31 missing axis %s", pumpRoleAxis)
			}
			pumpRoleIdentity[pumpRoleAxis] = pumpRoleValue
		}
		pumpRoleExpect(pumpRoleT, pumpRoleIdentity, pumpRoleExpected)
		pumpRoleExpect(pumpRoleT, pumpRoleEntry["grade"], `"Hypothesis"`)
		pumpRoleExpect(pumpRoleT, pumpRoleEntry["physical_qualification"], `"Unknown"`)
		pumpRoleExpect(pumpRoleT, pumpRoleEntry["logical_extent"], `null`)
		pumpRoleExpect(pumpRoleT, pumpRoleEntry["raw_fallback"], `"preserve native bytes, including unmapped values and reserved bits"`)
		if len(pumpRoleEntry) != 12 {
			pumpRoleT.Fatal("ab31 row fields changed")
		}
	}
	pumpRoleExpect(pumpRoleT, pumpRoleEntries[0].(map[string]any)["fields"], `[{
"name":"pump state","logical_position":null,"width_bytes":null,"raw_domain":[0,1],"signedness":"Unknown","byte_order":null,"conversion":null,"unit":null,"default":null,"mappings":{"0":"blocked","1":"released"},"reserved_regions":null,"grade":"Hypothesis","raw_fallback":"preserve native bytes, including unmapped values and reserved bits","conflicting_default_descriptions":["blocked","released"]}]`)
	pumpRoleExpect(pumpRoleT, pumpRoleEntries[1].(map[string]any)["fields"], `[
{"name":"cycle-pause","logical_position":3,"width_bytes":null,"raw_domain":null,"signedness":"Unknown","byte_order":null,"conversion":null,"unit":null,"default":null,"mappings":{},"reserved_regions":null,"grade":"Hypothesis","raw_fallback":"preserve native bytes, including unmapped values and reserved bits"},
{"name":"tank-one flag","logical_position":4,"width_bytes":1,"raw_domain":{"min":0,"max":255},"signedness":"unsigned","byte_order":null,"conversion":{"mask":1,"shift":0},"unit":null,"default":null,"mappings":{"0":"flag clear","1":"flag set"},"reserved_regions":[{"bits":[1,7]}],"grade":"Hypothesis","raw_fallback":"preserve native bytes, including unmapped values and reserved bits"},
{"name":"tank-two flag","logical_position":5,"width_bytes":1,"raw_domain":{"min":0,"max":255},"signedness":"unsigned","byte_order":null,"conversion":{"mask":1,"shift":0},"unit":null,"default":null,"mappings":{"0":"flag clear","1":"flag set"},"reserved_regions":[{"bits":[1,7]}],"grade":"Hypothesis","raw_fallback":"preserve native bytes, including unmapped values and reserved bits"}]`)
	pumpRoleExpect(pumpRoleT, pumpRoleEntries[2].(map[string]any)["fields"], `[]`)
}

func TestAb31IdentitySerialization(pumpRoleT *pumpRoleTesting.T) {
	for _, pumpRoleAxis := range []string{"family", "variant", "direction", "service", "operation", "alternative"} {
		pumpRoleT.Run("ab31_"+pumpRoleAxis, func(pumpRoleT *pumpRoleTesting.T) {
			for _, pumpRoleOriginal := range pumpRoleDecode(pumpRoleT, PumpRoles())["entries"].([]any) {
				pumpRoleFirst := pumpRoleOriginal.(map[string]any)
				pumpRoleSecond := make(map[string]any)
				for pumpRoleKey, pumpRoleValue := range pumpRoleFirst {
					pumpRoleSecond[pumpRoleKey] = pumpRoleValue
				}
				pumpRoleSecond[pumpRoleAxis] = "ab31_independent_alternative"
				pumpRoleData, pumpRoleError := pumpRoleJSON.Marshal(map[string]any{"entries": []any{pumpRoleFirst, pumpRoleSecond}})
				if pumpRoleError != nil {
					pumpRoleT.Fatal(pumpRoleError)
				}
				pumpRoleRoundTrip := pumpRoleDecode(pumpRoleT, pumpRoleData)["entries"].([]any)
				if len(pumpRoleRoundTrip) != 2 || !pumpRoleReflect.DeepEqual(pumpRoleRoundTrip[0], pumpRoleFirst) || !pumpRoleReflect.DeepEqual(pumpRoleRoundTrip[1], pumpRoleSecond) || pumpRoleReflect.DeepEqual(pumpRoleRoundTrip[0], pumpRoleRoundTrip[1]) {
					pumpRoleT.Fatal("ab31 independent identity axis lost")
				}
			}
		})
	}
}

func TestAb31IndependentTankFlags(pumpRoleT *pumpRoleTesting.T) {
	pumpRoleFields := pumpRoleDecode(pumpRoleT, PumpRoles())["entries"].([]any)[1].(map[string]any)["fields"].([]any)
	for _, pumpRoleCase := range []struct {
		pumpRoleName     string
		pumpRoleRaw      [6]byte
		pumpRoleFlags    [2]byte
		pumpRoleReserved [2]byte
	}{
		{"ab31_one", [6]byte{0, 0, 0, 0, 1, 0}, [2]byte{1, 0}, [2]byte{}},
		{"ab31_two", [6]byte{0, 0, 0, 0, 0, 1}, [2]byte{0, 1}, [2]byte{}},
		{"ab31_reserved", [6]byte{0, 0, 0, 0, 128, 0}, [2]byte{}, [2]byte{128, 0}},
	} {
		pumpRoleT.Run(pumpRoleCase.pumpRoleName, func(pumpRoleT *pumpRoleTesting.T) {
			pumpRoleBefore := pumpRoleCase.pumpRoleRaw
			for pumpRoleIndex, pumpRoleValue := range pumpRoleFields[1:] {
				pumpRoleField := pumpRoleValue.(map[string]any)
				pumpRolePosition := int(pumpRoleField["logical_position"].(float64))
				pumpRoleConversion := pumpRoleField["conversion"].(map[string]any)
				pumpRoleMask := byte(pumpRoleConversion["mask"].(float64))
				pumpRoleShift := uint(pumpRoleConversion["shift"].(float64))
				pumpRoleRaw := pumpRoleCase.pumpRoleRaw[pumpRolePosition]
				if (pumpRoleRaw&pumpRoleMask)>>pumpRoleShift != pumpRoleCase.pumpRoleFlags[pumpRoleIndex] || pumpRoleRaw & ^pumpRoleMask != pumpRoleCase.pumpRoleReserved[pumpRoleIndex] {
					pumpRoleT.Fatal("ab31 flags coupled or reserved evidence changed")
				}
			}
			if pumpRoleCase.pumpRoleRaw != pumpRoleBefore {
				pumpRoleT.Fatal("ab31 raw evidence changed")
			}
		})
	}
}

func TestAb31Copy(pumpRoleT *pumpRoleTesting.T) {
	pumpRoleFirst := PumpRoles()
	pumpRoleBefore := append([]byte(nil), pumpRoleFirst...)
	for pumpRoleIndex := range pumpRoleFirst {
		pumpRoleFirst[pumpRoleIndex] = '!'
	}
	if !pumpRoleBytes.Equal(PumpRoles(), pumpRoleBefore) {
		pumpRoleT.Fatal("ab31 returned bytes alias embedded catalog")
	}
}
