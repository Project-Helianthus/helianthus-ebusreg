package vaillant

import (
	b510OperatingDefaultBytes "bytes"
	b510OperatingDefaultJSON "encoding/json"
	b510OperatingDefaultReflect "reflect"
	b510OperatingDefaultTesting "testing"
)

func b510OperatingDefaultDecode(b510OperatingDefaultT *b510OperatingDefaultTesting.T, b510OperatingDefaultRaw []byte) map[string]any {
	b510OperatingDefaultT.Helper()
	var b510OperatingDefaultResult map[string]any
	if b510OperatingDefaultErr := b510OperatingDefaultJSON.Unmarshal(b510OperatingDefaultRaw, &b510OperatingDefaultResult); b510OperatingDefaultErr != nil {
		b510OperatingDefaultT.Fatal(b510OperatingDefaultErr)
	}
	return b510OperatingDefaultResult
}

func b510OperatingDefaultMarshal(b510OperatingDefaultT *b510OperatingDefaultTesting.T, b510OperatingDefaultValue any) []byte {
	b510OperatingDefaultT.Helper()
	b510OperatingDefaultRaw, b510OperatingDefaultErr := b510OperatingDefaultJSON.Marshal(b510OperatingDefaultValue)
	if b510OperatingDefaultErr != nil {
		b510OperatingDefaultT.Fatal(b510OperatingDefaultErr)
	}
	return b510OperatingDefaultRaw
}

func b510OperatingDefaultKeyFor(b510OperatingDefaultField string) B510OperatingDefaultKey {
	b510OperatingDefaultSecondary, b510OperatingDefaultOperation := uint8(16), uint8(0)
	return B510OperatingDefaultKey{B510OperatingDefaultFamily: "vaillant", B510OperatingDefaultVariant: "B510", B510OperatingDefaultPrimary: 181, B510OperatingDefaultSecondary: &b510OperatingDefaultSecondary, B510OperatingDefaultOperation: &b510OperatingDefaultOperation, B510OperatingDefaultDirection: "request", B510OperatingDefaultField: b510OperatingDefaultField, B510OperatingDefaultAlternative: "none", B510OperatingDefaultInstanceScope: "none"}
}

func TestAb48FiniteMetadata(b510OperatingDefaultT *b510OperatingDefaultTesting.T) {
	b510OperatingDefaultExpected := map[string]string{
		"mode":            `{"default":"forced-heating","numeric_mode_code":null,"symbolic_contexts":["forced-heating","forced-DHW","forced-off"],"mode_encodings":"Unknown"}`,
		"modulation":      `{"declared_default":67,"width_bits":null,"domain":null}`,
		"options":         `{"bit_defaults":[1,0,0,0,0,0,0,0],"bit2_context":{"field":"CH_setpoint","clear_role":"heating","set_role":"cooling","raw_value_changes":false},"width_bits":null}`,
		"loading_pump":    `{"width_bits":8,"signedness":"unsigned","conversion":"raw / 2","unit":"percent","formal_domain":[0,127.5],"declared_default":255.0,"default_conflict":true,"raw_example":{"bytes":"FF","value":127.5},"forced_mode_dependency":true,"contexts":[{"mode":"forced-heating","numeric_mode_code":null,"meaning":"Unknown","activation":"Unknown"},{"mode":"forced-DHW","numeric_mode_code":null,"meaning":"Unknown","activation":"Unknown"},{"mode":"forced-off","numeric_mode_code":null,"meaning":"Unknown","activation":"Unknown"}]}`,
		"CH_setpoint":     `{"width_bits":null,"zero_degrees_role":["stop","pump-stop"],"transition":{"from":0,"to":"positive","candidate_minimum_seconds":300,"execution_guarantee":false,"enforcement":"Unknown"},"other_transitions":"Unknown"}`,
		"pump_activation": `{"width_bits":null,"forced_mode_dependency":true,"contexts":[{"mode":"forced-heating","numeric_mode_code":null,"activation":"Unknown"},{"mode":"forced-DHW","numeric_mode_code":null,"activation":"Unknown"},{"mode":"forced-off","numeric_mode_code":null,"activation":"Unknown"}]}`,
	}
	b510OperatingDefaultEntries, b510OperatingDefaultErr := B510OperatingDefaultEntries()
	if b510OperatingDefaultErr != nil {
		b510OperatingDefaultT.Fatal(b510OperatingDefaultErr)
	}
	var b510OperatingDefaultRows []b510OperatingDefaultJSON.RawMessage
	if b510OperatingDefaultErr = b510OperatingDefaultJSON.Unmarshal(b510OperatingDefaultEntries, &b510OperatingDefaultRows); b510OperatingDefaultErr != nil || len(b510OperatingDefaultRows) != 6 {
		b510OperatingDefaultT.Fatalf("ab48 finite array: %d %v", len(b510OperatingDefaultRows), b510OperatingDefaultErr)
	}
	for b510OperatingDefaultField, b510OperatingDefaultDetails := range b510OperatingDefaultExpected {
		b510OperatingDefaultT.Run("ab48_"+b510OperatingDefaultField, func(b510OperatingDefaultT *b510OperatingDefaultTesting.T) {
			b510OperatingDefaultKey := b510OperatingDefaultKeyFor(b510OperatingDefaultField)
			b510OperatingDefaultRaw, b510OperatingDefaultFound, b510OperatingDefaultErr := B510OperatingDefaultLookup(b510OperatingDefaultKey)
			if b510OperatingDefaultErr != nil || !b510OperatingDefaultFound {
				b510OperatingDefaultT.Fatalf("ab48 lookup: %v %v", b510OperatingDefaultFound, b510OperatingDefaultErr)
			}
			b510OperatingDefaultExpectedRow := map[string]any{"id": "ab48_" + b510OperatingDefaultField, "key": b510OperatingDefaultKey, "grade": "Hypothesis", "metadata_presence": "Proven", "physical_qualification": "Unknown", "presence": "Unknown", "freshness": "Unknown", "execution": "Unknown", "native_contract": "https://github.com/Project-Helianthus/helianthus-docs-ebus/blob/main/protocols/vaillant/ebus-vaillant-b510.md", "position": nil, "raw_form": "preserve original bytes", "unknown": []string{"other positions", "peer acceptance", "effects"}, "details": b510OperatingDefaultDecode(b510OperatingDefaultT, []byte(b510OperatingDefaultDetails))}
			if !b510OperatingDefaultReflect.DeepEqual(b510OperatingDefaultDecode(b510OperatingDefaultT, b510OperatingDefaultRaw), b510OperatingDefaultDecode(b510OperatingDefaultT, b510OperatingDefaultMarshal(b510OperatingDefaultT, b510OperatingDefaultExpectedRow))) {
				b510OperatingDefaultT.Fatalf("ab48 whole row mismatch: %s", b510OperatingDefaultRaw)
			}
			b510OperatingDefaultMatches := 0
			for _, b510OperatingDefaultRow := range b510OperatingDefaultRows {
				if b510OperatingDefaultBytes.Equal(b510OperatingDefaultRow, b510OperatingDefaultRaw) {
					b510OperatingDefaultMatches++
				}
			}
			if b510OperatingDefaultMatches != 1 {
				b510OperatingDefaultT.Fatal("ab48 lookup must retain whole exact raw row")
			}
		})
	}
}

func TestAb48EveryKeyAxis(b510OperatingDefaultT *b510OperatingDefaultTesting.T) {
	b510OperatingDefaultZero, b510OperatingDefaultOne := uint8(0), uint8(1)
	b510OperatingDefaultRegisterZero, b510OperatingDefaultRegisterOne := uint16(0), uint16(1)
	b510OperatingDefaultChanges := map[string]func(*B510OperatingDefaultKey){
		"ab48_family": func(b510OperatingDefaultKey *B510OperatingDefaultKey) {
			b510OperatingDefaultKey.B510OperatingDefaultFamily = "ab48_other"
		},
		"ab48_variant": func(b510OperatingDefaultKey *B510OperatingDefaultKey) {
			b510OperatingDefaultKey.B510OperatingDefaultVariant = "ab48_other"
		},
		"ab48_primary": func(b510OperatingDefaultKey *B510OperatingDefaultKey) {
			b510OperatingDefaultKey.B510OperatingDefaultPrimary = 0
		},
		"ab48_secondary_absent": func(b510OperatingDefaultKey *B510OperatingDefaultKey) {
			b510OperatingDefaultKey.B510OperatingDefaultSecondary = nil
		},
		"ab48_secondary_zero": func(b510OperatingDefaultKey *B510OperatingDefaultKey) {
			b510OperatingDefaultKey.B510OperatingDefaultSecondary = &b510OperatingDefaultZero
		},
		"ab48_operation_absent": func(b510OperatingDefaultKey *B510OperatingDefaultKey) {
			b510OperatingDefaultKey.B510OperatingDefaultOperation = nil
		},
		"ab48_operation_one": func(b510OperatingDefaultKey *B510OperatingDefaultKey) {
			b510OperatingDefaultKey.B510OperatingDefaultOperation = &b510OperatingDefaultOne
		},
		"ab48_direction": func(b510OperatingDefaultKey *B510OperatingDefaultKey) {
			b510OperatingDefaultKey.B510OperatingDefaultDirection = "response"
		},
		"ab48_field": func(b510OperatingDefaultKey *B510OperatingDefaultKey) {
			b510OperatingDefaultKey.B510OperatingDefaultField = "ab48_other"
		},
		"ab48_alternative": func(b510OperatingDefaultKey *B510OperatingDefaultKey) {
			b510OperatingDefaultKey.B510OperatingDefaultAlternative = "ab48_other"
		},
		"ab48_group_zero": func(b510OperatingDefaultKey *B510OperatingDefaultKey) {
			b510OperatingDefaultKey.B510OperatingDefaultGroup = &b510OperatingDefaultZero
		},
		"ab48_group_one": func(b510OperatingDefaultKey *B510OperatingDefaultKey) {
			b510OperatingDefaultKey.B510OperatingDefaultGroup = &b510OperatingDefaultOne
		},
		"ab48_instance_scope": func(b510OperatingDefaultKey *B510OperatingDefaultKey) {
			b510OperatingDefaultKey.B510OperatingDefaultInstanceScope = "ab48_other"
		},
		"ab48_register_zero": func(b510OperatingDefaultKey *B510OperatingDefaultKey) {
			b510OperatingDefaultKey.B510OperatingDefaultRegister = &b510OperatingDefaultRegisterZero
		},
		"ab48_register_one": func(b510OperatingDefaultKey *B510OperatingDefaultKey) {
			b510OperatingDefaultKey.B510OperatingDefaultRegister = &b510OperatingDefaultRegisterOne
		},
	}
	for b510OperatingDefaultName, b510OperatingDefaultChange := range b510OperatingDefaultChanges {
		b510OperatingDefaultT.Run(b510OperatingDefaultName, func(b510OperatingDefaultT *b510OperatingDefaultTesting.T) {
			b510OperatingDefaultBase := b510OperatingDefaultKeyFor("loading_pump")
			b510OperatingDefaultOriginal, b510OperatingDefaultFound, b510OperatingDefaultErr := B510OperatingDefaultLookup(b510OperatingDefaultBase)
			if b510OperatingDefaultErr != nil || !b510OperatingDefaultFound {
				b510OperatingDefaultT.Fatal("ab48 invalid control")
			}
			b510OperatingDefaultChanged := b510OperatingDefaultKeyFor("loading_pump")
			b510OperatingDefaultChange(&b510OperatingDefaultChanged)
			if b510OperatingDefaultRaw, b510OperatingDefaultFound, b510OperatingDefaultErr := B510OperatingDefaultLookup(b510OperatingDefaultChanged); b510OperatingDefaultErr != nil || b510OperatingDefaultFound || b510OperatingDefaultRaw != nil {
				b510OperatingDefaultT.Fatal("ab48 namespace axis ignored")
			}
			b510OperatingDefaultSecond := b510OperatingDefaultDecode(b510OperatingDefaultT, b510OperatingDefaultOriginal)
			b510OperatingDefaultSecond["key"] = b510OperatingDefaultChanged
			b510OperatingDefaultSource := b510OperatingDefaultMarshal(b510OperatingDefaultT, []any{b510OperatingDefaultJSON.RawMessage(b510OperatingDefaultOriginal), b510OperatingDefaultSecond})
			b510OperatingDefaultWant := b510OperatingDefaultDecode(b510OperatingDefaultT, b510OperatingDefaultMarshal(b510OperatingDefaultT, b510OperatingDefaultSecond))
			b510OperatingDefaultRaw, b510OperatingDefaultFound, b510OperatingDefaultErr := b510OperatingDefaultLookupSource(b510OperatingDefaultSource, b510OperatingDefaultChanged)
			if b510OperatingDefaultErr != nil || !b510OperatingDefaultFound || !b510OperatingDefaultReflect.DeepEqual(b510OperatingDefaultDecode(b510OperatingDefaultT, b510OperatingDefaultRaw), b510OperatingDefaultWant) {
				b510OperatingDefaultT.Fatalf("ab48 distinct valid control: %v", b510OperatingDefaultErr)
			}
		})
	}
}

func TestAb48WidthAndGradeRetention(b510OperatingDefaultT *b510OperatingDefaultTesting.T) {
	b510OperatingDefaultRaw, _, b510OperatingDefaultErr := B510OperatingDefaultLookup(b510OperatingDefaultKeyFor("loading_pump"))
	if b510OperatingDefaultErr != nil {
		b510OperatingDefaultT.Fatal(b510OperatingDefaultErr)
	}
	for b510OperatingDefaultName, b510OperatingDefaultValue := range map[string]any{"ab48_width": float64(16), "ab48_grade": "Unknown"} {
		b510OperatingDefaultT.Run(b510OperatingDefaultName, func(b510OperatingDefaultT *b510OperatingDefaultTesting.T) {
			b510OperatingDefaultRow := b510OperatingDefaultDecode(b510OperatingDefaultT, b510OperatingDefaultRaw)
			if b510OperatingDefaultName == "ab48_width" {
				b510OperatingDefaultRow["details"].(map[string]any)["width_bits"] = b510OperatingDefaultValue
			} else {
				b510OperatingDefaultRow["grade"] = b510OperatingDefaultValue
			}
			b510OperatingDefaultGot, b510OperatingDefaultFound, b510OperatingDefaultErr := b510OperatingDefaultLookupSource(b510OperatingDefaultMarshal(b510OperatingDefaultT, []any{b510OperatingDefaultRow}), b510OperatingDefaultKeyFor("loading_pump"))
			if b510OperatingDefaultErr != nil || !b510OperatingDefaultFound || !b510OperatingDefaultReflect.DeepEqual(b510OperatingDefaultDecode(b510OperatingDefaultT, b510OperatingDefaultGot), b510OperatingDefaultRow) {
				b510OperatingDefaultT.Fatal("ab48 row metadata normalized or lost")
			}
			if _, _, b510OperatingDefaultErr = b510OperatingDefaultValidate(b510OperatingDefaultMarshal(b510OperatingDefaultT, []any{b510OperatingDefaultJSON.RawMessage(b510OperatingDefaultRaw), b510OperatingDefaultRow})); b510OperatingDefaultErr == nil {
				b510OperatingDefaultT.Fatal("ab48 metadata must not disambiguate duplicate identity")
			}
		})
	}
}

func TestAb48OptionalAbsentVersusZero(b510OperatingDefaultT *b510OperatingDefaultTesting.T) {
	b510OperatingDefaultZero, b510OperatingDefaultRegisterZero := uint8(0), uint16(0)
	for _, b510OperatingDefaultAxis := range []string{"secondary", "operation", "group", "register"} {
		b510OperatingDefaultT.Run("ab48_"+b510OperatingDefaultAxis, func(b510OperatingDefaultT *b510OperatingDefaultTesting.T) {
			b510OperatingDefaultAbsent, b510OperatingDefaultExplicit := b510OperatingDefaultKeyFor("mode"), b510OperatingDefaultKeyFor("mode")
			switch b510OperatingDefaultAxis {
			case "secondary":
				b510OperatingDefaultAbsent.B510OperatingDefaultSecondary, b510OperatingDefaultExplicit.B510OperatingDefaultSecondary = nil, &b510OperatingDefaultZero
			case "operation":
				b510OperatingDefaultAbsent.B510OperatingDefaultOperation, b510OperatingDefaultExplicit.B510OperatingDefaultOperation = nil, &b510OperatingDefaultZero
			case "group":
				b510OperatingDefaultAbsent.B510OperatingDefaultGroup, b510OperatingDefaultExplicit.B510OperatingDefaultGroup = nil, &b510OperatingDefaultZero
			case "register":
				b510OperatingDefaultAbsent.B510OperatingDefaultRegister, b510OperatingDefaultExplicit.B510OperatingDefaultRegister = nil, &b510OperatingDefaultRegisterZero
			}
			b510OperatingDefaultSource := b510OperatingDefaultMarshal(b510OperatingDefaultT, []any{map[string]any{"key": b510OperatingDefaultAbsent, "id": "ab48_absent"}, map[string]any{"key": b510OperatingDefaultExplicit, "id": "ab48_zero"}})
			for b510OperatingDefaultIndex, b510OperatingDefaultKey := range []B510OperatingDefaultKey{b510OperatingDefaultAbsent, b510OperatingDefaultExplicit} {
				b510OperatingDefaultRaw, b510OperatingDefaultFound, b510OperatingDefaultErr := b510OperatingDefaultLookupSource(b510OperatingDefaultSource, b510OperatingDefaultKey)
				if b510OperatingDefaultErr != nil || !b510OperatingDefaultFound || b510OperatingDefaultDecode(b510OperatingDefaultT, b510OperatingDefaultRaw)["id"] != []string{"ab48_absent", "ab48_zero"}[b510OperatingDefaultIndex] {
					b510OperatingDefaultT.Fatalf("ab48 absent and zero conflated: %v", b510OperatingDefaultErr)
				}
			}
		})
	}
}

func TestAb48InvalidAndDuplicate(b510OperatingDefaultT *b510OperatingDefaultTesting.T) {
	b510OperatingDefaultRow, _, _ := B510OperatingDefaultLookup(b510OperatingDefaultKeyFor("mode"))
	b510OperatingDefaultCases := map[string][]byte{"ab48_syntax": []byte("["), "ab48_null": []byte("null"), "ab48_object": []byte("{}"), "ab48_null_row": []byte("[null]"), "ab48_missing_key": []byte("[{}]"), "ab48_trailing": []byte("[] []"), "ab48_duplicate": b510OperatingDefaultMarshal(b510OperatingDefaultT, []any{b510OperatingDefaultJSON.RawMessage(b510OperatingDefaultRow), b510OperatingDefaultJSON.RawMessage(b510OperatingDefaultRow)}), "ab48_invalid_later": b510OperatingDefaultMarshal(b510OperatingDefaultT, []any{b510OperatingDefaultJSON.RawMessage(b510OperatingDefaultRow), nil})}
	for _, b510OperatingDefaultAxis := range []string{"family", "variant", "primary", "direction", "field", "alternative", "instance_scope"} {
		b510OperatingDefaultChanged := b510OperatingDefaultDecode(b510OperatingDefaultT, b510OperatingDefaultRow)
		delete(b510OperatingDefaultChanged["key"].(map[string]any), b510OperatingDefaultAxis)
		b510OperatingDefaultCases["ab48_missing_"+b510OperatingDefaultAxis] = b510OperatingDefaultMarshal(b510OperatingDefaultT, []any{b510OperatingDefaultChanged})
	}
	b510OperatingDefaultSaved := b510OperatingDefaultData
	defer func() { b510OperatingDefaultData = b510OperatingDefaultSaved }()
	for b510OperatingDefaultName, b510OperatingDefaultSource := range b510OperatingDefaultCases {
		b510OperatingDefaultT.Run(b510OperatingDefaultName, func(b510OperatingDefaultT *b510OperatingDefaultTesting.T) {
			b510OperatingDefaultData = string(b510OperatingDefaultSource)
			if b510OperatingDefaultRaw, b510OperatingDefaultErr := B510OperatingDefaultEntries(); b510OperatingDefaultErr == nil || b510OperatingDefaultRaw != nil {
				b510OperatingDefaultT.Fatal("ab48 invalid entries accepted")
			}
			if b510OperatingDefaultRaw, b510OperatingDefaultFound, b510OperatingDefaultErr := B510OperatingDefaultLookup(b510OperatingDefaultKeyFor("mode")); b510OperatingDefaultErr == nil || b510OperatingDefaultFound || b510OperatingDefaultRaw != nil {
				b510OperatingDefaultT.Fatal("ab48 invalid lookup accepted")
			}
		})
	}
}

func TestAb48StorageIsolation(b510OperatingDefaultT *b510OperatingDefaultTesting.T) {
	b510OperatingDefaultBefore, b510OperatingDefaultErr := B510OperatingDefaultEntries()
	if b510OperatingDefaultErr != nil {
		b510OperatingDefaultT.Fatal(b510OperatingDefaultErr)
	}
	b510OperatingDefaultSource := b510OperatingDefaultBytes.Clone(b510OperatingDefaultBefore)
	b510OperatingDefaultRows, _, b510OperatingDefaultErr := b510OperatingDefaultValidate(b510OperatingDefaultSource)
	if b510OperatingDefaultErr != nil {
		b510OperatingDefaultT.Fatal(b510OperatingDefaultErr)
	}
	b510OperatingDefaultRowBefore := b510OperatingDefaultBytes.Clone(b510OperatingDefaultRows[0])
	b510OperatingDefaultLookupBefore, _, _ := B510OperatingDefaultLookup(b510OperatingDefaultKeyFor("mode"))
	b510OperatingDefaultLookupCopy, _, _ := b510OperatingDefaultLookupSource(b510OperatingDefaultSource, b510OperatingDefaultKeyFor("mode"))
	for b510OperatingDefaultIndex := range b510OperatingDefaultSource {
		b510OperatingDefaultSource[b510OperatingDefaultIndex] = 'x'
	}
	if !b510OperatingDefaultBytes.Equal(b510OperatingDefaultRows[0], b510OperatingDefaultRowBefore) || !b510OperatingDefaultBytes.Equal(b510OperatingDefaultLookupCopy, b510OperatingDefaultLookupBefore) {
		b510OperatingDefaultT.Fatal("ab48 source mutation aliases results")
	}
	b510OperatingDefaultReturned, _ := B510OperatingDefaultEntries()
	for b510OperatingDefaultIndex := range b510OperatingDefaultReturned {
		b510OperatingDefaultReturned[b510OperatingDefaultIndex] = 'x'
	}
	for b510OperatingDefaultIndex := range b510OperatingDefaultLookupCopy {
		b510OperatingDefaultLookupCopy[b510OperatingDefaultIndex] = 'x'
	}
	b510OperatingDefaultPublicCopy, _, _ := B510OperatingDefaultLookup(b510OperatingDefaultKeyFor("mode"))
	for b510OperatingDefaultIndex := range b510OperatingDefaultPublicCopy {
		b510OperatingDefaultPublicCopy[b510OperatingDefaultIndex] = 'x'
	}
	b510OperatingDefaultAfter, _ := B510OperatingDefaultEntries()
	b510OperatingDefaultLookupAfter, _, _ := B510OperatingDefaultLookup(b510OperatingDefaultKeyFor("mode"))
	if !b510OperatingDefaultBytes.Equal(b510OperatingDefaultBefore, b510OperatingDefaultAfter) || !b510OperatingDefaultBytes.Equal(b510OperatingDefaultLookupBefore, b510OperatingDefaultLookupAfter) {
		b510OperatingDefaultT.Fatal("ab48 returned memory aliases catalog")
	}
}

func TestAb48SyntheticScenarios(b510OperatingDefaultT *b510OperatingDefaultTesting.T) {
	b510OperatingDefaultPumpRaw, _, _ := B510OperatingDefaultLookup(b510OperatingDefaultKeyFor("loading_pump"))
	b510OperatingDefaultPump := b510OperatingDefaultDecode(b510OperatingDefaultT, b510OperatingDefaultPumpRaw)["details"].(map[string]any)
	if b510OperatingDefaultPump["conversion"] != "raw / 2" || float64(byte(0xff))/2 != 127.5 || b510OperatingDefaultPump["raw_example"].(map[string]any)["value"] != 127.5 || b510OperatingDefaultPump["declared_default"] != 255.0 {
		b510OperatingDefaultT.Fatal("ab48 raw FF and conflicting literal lost")
	}
	b510OperatingDefaultOptionsRaw, _, _ := B510OperatingDefaultLookup(b510OperatingDefaultKeyFor("options"))
	b510OperatingDefaultContext := b510OperatingDefaultDecode(b510OperatingDefaultT, b510OperatingDefaultOptionsRaw)["details"].(map[string]any)["bit2_context"].(map[string]any)
	b510OperatingDefaultRaw := byte(42)
	for b510OperatingDefaultIndex, b510OperatingDefaultOption := range []byte{0, 4} {
		b510OperatingDefaultRole := b510OperatingDefaultContext["clear_role"]
		if b510OperatingDefaultOption&4 != 0 {
			b510OperatingDefaultRole = b510OperatingDefaultContext["set_role"]
		}
		if b510OperatingDefaultRole != []string{"heating", "cooling"}[b510OperatingDefaultIndex] || b510OperatingDefaultRaw != 42 || b510OperatingDefaultContext["raw_value_changes"] != false {
			b510OperatingDefaultT.Fatal("ab48 option context changed raw")
		}
	}
	b510OperatingDefaultSetpointRaw, _, _ := B510OperatingDefaultLookup(b510OperatingDefaultKeyFor("CH_setpoint"))
	b510OperatingDefaultDetails := b510OperatingDefaultDecode(b510OperatingDefaultT, b510OperatingDefaultSetpointRaw)["details"].(map[string]any)
	b510OperatingDefaultTransition := b510OperatingDefaultDetails["transition"].(map[string]any)
	for b510OperatingDefaultIndex, b510OperatingDefaultPair := range [][2]int{{0, 1}, {1, 2}} {
		b510OperatingDefaultCandidate := b510OperatingDefaultPair[0] == int(b510OperatingDefaultTransition["from"].(float64)) && b510OperatingDefaultPair[1] > 0
		if b510OperatingDefaultCandidate != (b510OperatingDefaultIndex == 0) {
			b510OperatingDefaultT.Fatal("ab48 distinct transition collapsed")
		}
		if b510OperatingDefaultCandidate && b510OperatingDefaultTransition["candidate_minimum_seconds"] != float64(300) {
			b510OperatingDefaultT.Fatal("ab48 candidate interval")
		}
		if !b510OperatingDefaultCandidate && b510OperatingDefaultDetails["other_transitions"] != "Unknown" {
			b510OperatingDefaultT.Fatal("ab48 invented interval")
		}
	}
	if b510OperatingDefaultTransition["execution_guarantee"] != false || b510OperatingDefaultTransition["enforcement"] != "Unknown" {
		b510OperatingDefaultT.Fatal("ab48 candidate promoted to guarantee")
	}
}
