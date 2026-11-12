package ebusstandard

import (
	standardFFBodyStateBytes "bytes"
	standardFFBodyStateJSON "encoding/json"
	standardFFBodyStateFmt "fmt"
	standardFFBodyStateReflect "reflect"
	standardFFBodyStateStrings "strings"
	standardFFBodyStateTesting "testing"
)

func standardFFBodyStateByte(standardFFBodyStateValue uint8) *uint8 { return &standardFFBodyStateValue }
func standardFFBodyStateWord(standardFFBodyStateValue uint16) *uint16 {
	return &standardFFBodyStateValue
}

func standardFFBodyStateTestKey(standardFFBodyStateSecondary uint8) StandardFFBodyStateKey {
	standardFFBodyStateDirection := "none"
	if standardFFBodyStateSecondary >= 3 {
		standardFFBodyStateDirection = "abstract-request"
	}
	return StandardFFBodyStateKey{StandardFFBodyStateFamily: "standard-service", StandardFFBodyStateVariant: standardFFBodyStateFmt.Sprintf("FF%02d", standardFFBodyStateSecondary), StandardFFBodyStatePrimary: 255,
		StandardFFBodyStateSecondary: standardFFBodyStateByte(standardFFBodyStateSecondary), StandardFFBodyStateDirection: standardFFBodyStateDirection, StandardFFBodyStateField: "body-completeness",
		StandardFFBodyStateAlternative: "none", StandardFFBodyStateInstanceScope: "none"}
}

func standardFFBodyStateObject(standardFFBodyStateT *standardFFBodyStateTesting.T, standardFFBodyStateData []byte) map[string]standardFFBodyStateJSON.RawMessage {
	standardFFBodyStateT.Helper()
	var standardFFBodyStateResult map[string]standardFFBodyStateJSON.RawMessage
	if standardFFBodyStateErr := standardFFBodyStateJSON.Unmarshal(standardFFBodyStateData, &standardFFBodyStateResult); standardFFBodyStateErr != nil {
		standardFFBodyStateT.Fatal(standardFFBodyStateErr)
	}
	return standardFFBodyStateResult
}

func standardFFBodyStateExpect(standardFFBodyStateT *standardFFBodyStateTesting.T, standardFFBodyStateRow map[string]standardFFBodyStateJSON.RawMessage, standardFFBodyStateField string, standardFFBodyStateWant string) {
	standardFFBodyStateT.Helper()
	if !standardFFBodyStateBytes.Equal(standardFFBodyStateBytes.TrimSpace(standardFFBodyStateRow[standardFFBodyStateField]), []byte(standardFFBodyStateWant)) {
		standardFFBodyStateT.Fatalf("ab54 %s = %s, want %s", standardFFBodyStateField, standardFFBodyStateRow[standardFFBodyStateField], standardFFBodyStateWant)
	}
}

func standardFFBodyStateGet(standardFFBodyStateT *standardFFBodyStateTesting.T, standardFFBodyStateKey StandardFFBodyStateKey) standardFFBodyStateJSON.RawMessage {
	standardFFBodyStateT.Helper()
	standardFFBodyStateRow, standardFFBodyStateFound, standardFFBodyStateErr := StandardFFBodyStateLookup(standardFFBodyStateKey)
	if standardFFBodyStateErr != nil || !standardFFBodyStateFound {
		standardFFBodyStateT.Fatalf("ab54 missing row: found=%v error=%v", standardFFBodyStateFound, standardFFBodyStateErr)
	}
	return standardFFBodyStateRow
}

func TestAb54Inventory(standardFFBodyStateT *standardFFBodyStateTesting.T) {
	standardFFBodyStateData, standardFFBodyStateErr := StandardFFBodyStateEntries()
	if standardFFBodyStateErr != nil {
		standardFFBodyStateT.Fatal(standardFFBodyStateErr)
	}
	var standardFFBodyStateRows []standardFFBodyStateJSON.RawMessage
	if standardFFBodyStateErr := standardFFBodyStateJSON.Unmarshal(standardFFBodyStateData, &standardFFBodyStateRows); standardFFBodyStateErr != nil {
		standardFFBodyStateT.Fatal(standardFFBodyStateErr)
	}
	if len(standardFFBodyStateRows) != 7 {
		standardFFBodyStateT.Fatalf("ab54 identities = %d, want 7", len(standardFFBodyStateRows))
	}
	standardFFBodyStateQueries := []string{"none", "none", "none", "net-status-query", "monitored-participants-query", "failed-nodes-query", "required-services-query"}
	standardFFBodyStateNN := []string{"Unknown", "0", "0", "Unknown", "1..10", "Unknown", "1+2n for complete service pairs"}
	for standardFFBodyStateSB := uint8(0); standardFFBodyStateSB < 7; standardFFBodyStateSB++ {
		standardFFBodyStateRaw := standardFFBodyStateGet(standardFFBodyStateT, standardFFBodyStateTestKey(standardFFBodyStateSB))
		if !standardFFBodyStateBytes.Equal(standardFFBodyStateRaw, standardFFBodyStateRows[standardFFBodyStateSB]) {
			standardFFBodyStateT.Fatal("ab54 lookup did not return whole exact row")
		}
		standardFFBodyStateRow := standardFFBodyStateObject(standardFFBodyStateT, standardFFBodyStateRaw)
		standardFFBodyStateExpect(standardFFBodyStateT, standardFFBodyStateRow, "id", standardFFBodyStateFmt.Sprintf("\"ab54-ff%02d\"", standardFFBodyStateSB))
		standardFFBodyStateExpect(standardFFBodyStateT, standardFFBodyStateRow, "grade", `"Proven"`)
		standardFFBodyStateExpect(standardFFBodyStateT, standardFFBodyStateRow, "invocable", "false")
		standardFFBodyStateExpect(standardFFBodyStateT, standardFFBodyStateRow, "executable", "false")
		standardFFBodyStateExpect(standardFFBodyStateT, standardFFBodyStateRow, "semantic_validity", `"Unknown"`)
		standardFFBodyStateExpect(standardFFBodyStateT, standardFFBodyStateRow, "physical_qualification", `"Unknown"`)
		standardFFBodyStateExpect(standardFFBodyStateT, standardFFBodyStateRow, "query_identity", `"`+standardFFBodyStateQueries[standardFFBodyStateSB]+`"`)
		standardFFBodyStateReply := `"Unknown"`
		if standardFFBodyStateSB >= 3 {
			standardFFBodyStateReply = `"same-service"`
		}
		standardFFBodyStateExpect(standardFFBodyStateT, standardFFBodyStateRow, "reply_service", standardFFBodyStateReply)
		standardFFBodyStateForm := standardFFBodyStateObject(standardFFBodyStateT, standardFFBodyStateRow["raw_form"])
		standardFFBodyStateExpect(standardFFBodyStateT, standardFFBodyStateForm, "PB", `"FF"`)
		standardFFBodyStateExpect(standardFFBodyStateT, standardFFBodyStateForm, "SB", standardFFBodyStateFmt.Sprintf("\"%02d\"", standardFFBodyStateSB))
		standardFFBodyStateExpect(standardFFBodyStateT, standardFFBodyStateForm, "NN", `"`+standardFFBodyStateNN[standardFFBodyStateSB]+`"`)
		standardFFBodyStateExpect(standardFFBodyStateT, standardFFBodyStateForm, "body", `"Unknown"`)
		standardFFBodyStateUnknowns := standardFFBodyStateObject(standardFFBodyStateT, standardFFBodyStateRow["unknowns"])
		for _, standardFFBodyStateAxis := range []string{"native_emission", "forwarding_eligibility", "target_identity", "correlation", "lifecycle_effects", "installed_support", "physical_ack", "exact_body_fields", "pair_meaning", "presence", "freshness"} {
			standardFFBodyStateExpect(standardFFBodyStateT, standardFFBodyStateUnknowns, standardFFBodyStateAxis, `"Unknown"`)
		}
		if standardFFBodyStateSB != 4 {
			standardFFBodyStateExpect(standardFFBodyStateT, standardFFBodyStateRow, "minimum_overhead_width", "null")
		}
		if standardFFBodyStateSB != 6 {
			standardFFBodyStateExpect(standardFFBodyStateT, standardFFBodyStateRow, "service_pair_width", "null")
		}
		if standardFFBodyStateSB != 4 && standardFFBodyStateSB != 6 {
			standardFFBodyStateExpect(standardFFBodyStateT, standardFFBodyStateRow, "cases", "[]")
		}
	}
	standardFFBodyStateReset := standardFFBodyStateObject(standardFFBodyStateT, standardFFBodyStateGet(standardFFBodyStateT, standardFFBodyStateTestKey(1)))
	standardFFBodyStateExpect(standardFFBodyStateT, standardFFBodyStateReset, "policy_boundary", `"NN0; reset of dynamic target configuration requires explicit emission policy; no emitter added"`)
	standardFFBodyStateAvailability := standardFFBodyStateObject(standardFFBodyStateT, standardFFBodyStateGet(standardFFBodyStateT, standardFFBodyStateTestKey(2)))
	standardFFBodyStateExpect(standardFFBodyStateT, standardFFBodyStateAvailability, "policy_boundary", `"NN0 independent of responder availability; no responder inference"`)
}

func TestAb54BodyCompleteness(standardFFBodyStateT *standardFFBodyStateTesting.T) {
	for _, standardFFBodyStateSB := range []uint8{4, 6} {
		standardFFBodyStateRow := standardFFBodyStateObject(standardFFBodyStateT, standardFFBodyStateGet(standardFFBodyStateT, standardFFBodyStateTestKey(standardFFBodyStateSB)))
		var standardFFBodyStateCases []map[string]standardFFBodyStateJSON.RawMessage
		if standardFFBodyStateErr := standardFFBodyStateJSON.Unmarshal(standardFFBodyStateRow["cases"], &standardFFBodyStateCases); standardFFBodyStateErr != nil {
			standardFFBodyStateT.Fatal(standardFFBodyStateErr)
		}
		standardFFBodyStateCount := 2
		if standardFFBodyStateSB == 6 {
			standardFFBodyStateCount = 5
		}
		if len(standardFFBodyStateCases) != standardFFBodyStateCount {
			standardFFBodyStateT.Fatal("ab54 case inventory changed")
		}
		for standardFFBodyStateIndex, standardFFBodyStateCase := range standardFFBodyStateCases {
			standardFFBodyStateNN := standardFFBodyStateIndex + 1
			standardFFBodyStateExpect(standardFFBodyStateT, standardFFBodyStateCase, "id", standardFFBodyStateFmt.Sprintf("\"ab54-ff%02d-length%d\"", standardFFBodyStateSB, standardFFBodyStateNN))
			standardFFBodyStateExpect(standardFFBodyStateT, standardFFBodyStateCase, "nn", standardFFBodyStateFmt.Sprint(standardFFBodyStateNN))
			standardFFBodyStateExpect(standardFFBodyStateT, standardFFBodyStateCase, "grade", `"Proven"`)
			standardFFBodyStateExpect(standardFFBodyStateT, standardFFBodyStateCase, "semantic_validity", `"Unknown"`)
			var standardFFBodyStateRaw []int
			if standardFFBodyStateErr := standardFFBodyStateJSON.Unmarshal(standardFFBodyStateCase["raw"], &standardFFBodyStateRaw); standardFFBodyStateErr != nil {
				standardFFBodyStateT.Fatal(standardFFBodyStateErr)
			}
			if !standardFFBodyStateReflect.DeepEqual(standardFFBodyStateRaw, make([]int, standardFFBodyStateNN)) {
				standardFFBodyStateT.Fatal("ab54 synthetic raw body changed")
			}
			if standardFFBodyStateSB == 4 {
				standardFFBodyStateExpect(standardFFBodyStateT, standardFFBodyStateRow, "minimum_overhead_width", "2")
				standardFFBodyStateExpect(standardFFBodyStateT, standardFFBodyStateCase, "minimum_overhead_met", []string{"false", "true"}[standardFFBodyStateIndex])
				standardFFBodyStateExpect(standardFFBodyStateT, standardFFBodyStateCase, "body_complete", []string{"false", "null"}[standardFFBodyStateIndex])
				standardFFBodyStateExpect(standardFFBodyStateT, standardFFBodyStateCase, "state", []string{`"truncated/raw"`, `"minimum-overhead-only"`}[standardFFBodyStateIndex])
			} else {
				standardFFBodyStateExpect(standardFFBodyStateT, standardFFBodyStateRow, "service_pair_width", "2")
				standardFFBodyStateExpect(standardFFBodyStateT, standardFFBodyStateCase, "complete_pairs", []string{"0", "0", "1", "1", "2"}[standardFFBodyStateIndex])
				standardFFBodyStateExpect(standardFFBodyStateT, standardFFBodyStateCase, "pair_shape_complete", []string{"true", "false", "true", "false", "true"}[standardFFBodyStateIndex])
				standardFFBodyStateExpect(standardFFBodyStateT, standardFFBodyStateCase, "body_complete", []string{"null", "false", "null", "false", "null"}[standardFFBodyStateIndex])
				standardFFBodyStateState, standardFFBodyStateRemainder := `"arithmetic-only"`, []int{}
				if standardFFBodyStateNN == 2 || standardFFBodyStateNN == 4 {
					standardFFBodyStateState, standardFFBodyStateRemainder = `"incomplete/raw-remainder"`, []int{0}
				}
				standardFFBodyStateExpect(standardFFBodyStateT, standardFFBodyStateCase, "state", standardFFBodyStateState)
				var standardFFBodyStateGotRemainder []int
				if standardFFBodyStateErr := standardFFBodyStateJSON.Unmarshal(standardFFBodyStateCase["raw_remainder"], &standardFFBodyStateGotRemainder); standardFFBodyStateErr != nil {
					standardFFBodyStateT.Fatal(standardFFBodyStateErr)
				}
				if !standardFFBodyStateReflect.DeepEqual(standardFFBodyStateGotRemainder, standardFFBodyStateRemainder) {
					standardFFBodyStateT.Fatal("ab54 lost trailing raw byte")
				}
			}
		}
	}
}

func TestAb54ExactAxes(standardFFBodyStateT *standardFFBodyStateTesting.T) {
	standardFFBodyStateControl := standardFFBodyStateTestKey(4)
	standardFFBodyStateRaw := standardFFBodyStateGet(standardFFBodyStateT, standardFFBodyStateControl)
	standardFFBodyStateMutations := []struct {
		StandardFFBodyStateName   string
		StandardFFBodyStateChange func(*StandardFFBodyStateKey)
	}{
		{"ab54-family", func(standardFFBodyStateK *StandardFFBodyStateKey) {
			standardFFBodyStateK.StandardFFBodyStateFamily = "other-family"
		}},
		{"ab54-variant", func(standardFFBodyStateK *StandardFFBodyStateKey) {
			standardFFBodyStateK.StandardFFBodyStateVariant = "FF03"
		}},
		{"ab54-primary", func(standardFFBodyStateK *StandardFFBodyStateKey) {
			standardFFBodyStateK.StandardFFBodyStatePrimary = 254
		}},
		{"ab54-primary-zero", func(standardFFBodyStateK *StandardFFBodyStateKey) {
			standardFFBodyStateK.StandardFFBodyStatePrimary = 0
		}},
		{"ab54-secondary", func(standardFFBodyStateK *StandardFFBodyStateKey) {
			standardFFBodyStateK.StandardFFBodyStateSecondary = standardFFBodyStateByte(3)
		}},
		{"ab54-secondary-absent", func(standardFFBodyStateK *StandardFFBodyStateKey) {
			standardFFBodyStateK.StandardFFBodyStateSecondary = nil
		}},
		{"ab54-secondary-zero", func(standardFFBodyStateK *StandardFFBodyStateKey) {
			standardFFBodyStateK.StandardFFBodyStateSecondary = standardFFBodyStateByte(0)
		}},
		{"ab54-operation-zero", func(standardFFBodyStateK *StandardFFBodyStateKey) {
			standardFFBodyStateK.StandardFFBodyStateOperation = standardFFBodyStateByte(0)
		}},
		{"ab54-operation-value", func(standardFFBodyStateK *StandardFFBodyStateKey) {
			standardFFBodyStateK.StandardFFBodyStateOperation = standardFFBodyStateByte(1)
		}},
		{"ab54-direction", func(standardFFBodyStateK *StandardFFBodyStateKey) {
			standardFFBodyStateK.StandardFFBodyStateDirection = "abstract-positive"
		}},
		{"ab54-field", func(standardFFBodyStateK *StandardFFBodyStateKey) {
			standardFFBodyStateK.StandardFFBodyStateField = "other-field"
		}},
		{"ab54-alternative", func(standardFFBodyStateK *StandardFFBodyStateKey) {
			standardFFBodyStateK.StandardFFBodyStateAlternative = "other-alternative"
		}},
		{"ab54-group-zero", func(standardFFBodyStateK *StandardFFBodyStateKey) {
			standardFFBodyStateK.StandardFFBodyStateGroup = standardFFBodyStateWord(0)
		}},
		{"ab54-group-value", func(standardFFBodyStateK *StandardFFBodyStateKey) {
			standardFFBodyStateK.StandardFFBodyStateGroup = standardFFBodyStateWord(1)
		}},
		{"ab54-instance-scope", func(standardFFBodyStateK *StandardFFBodyStateKey) {
			standardFFBodyStateK.StandardFFBodyStateInstanceScope = "other-scope"
		}},
		{"ab54-register-zero", func(standardFFBodyStateK *StandardFFBodyStateKey) {
			standardFFBodyStateK.StandardFFBodyStateRegister = standardFFBodyStateWord(0)
		}},
		{"ab54-register-value", func(standardFFBodyStateK *StandardFFBodyStateKey) {
			standardFFBodyStateK.StandardFFBodyStateRegister = standardFFBodyStateWord(1)
		}},
	}
	for _, standardFFBodyStateMutation := range standardFFBodyStateMutations {
		standardFFBodyStateT.Run(standardFFBodyStateMutation.StandardFFBodyStateName, func(standardFFBodyStateT *standardFFBodyStateTesting.T) {
			standardFFBodyStateKey := standardFFBodyStateControl
			standardFFBodyStateMutation.StandardFFBodyStateChange(&standardFFBodyStateKey)
			if standardFFBodyStateGot, standardFFBodyStateFound, standardFFBodyStateErr := StandardFFBodyStateLookup(standardFFBodyStateKey); standardFFBodyStateErr != nil || standardFFBodyStateFound || standardFFBodyStateGot != nil {
				standardFFBodyStateT.Fatal("ab54 changed axis selected production row")
			}
			standardFFBodyStateOther := standardFFBodyStateObject(standardFFBodyStateT, standardFFBodyStateRaw)
			standardFFBodyStateOther["key"], _ = standardFFBodyStateJSON.Marshal(standardFFBodyStateKey)
			standardFFBodyStateOther["id"] = standardFFBodyStateJSON.RawMessage(`"ab54-synthetic-axis-control"`)
			standardFFBodyStateOtherRaw, _ := standardFFBodyStateJSON.Marshal(standardFFBodyStateOther)
			standardFFBodyStateSource := append(append(append([]byte{'['}, standardFFBodyStateRaw...), ','), standardFFBodyStateOtherRaw...)
			standardFFBodyStateSource = append(standardFFBodyStateSource, ']')
			for _, standardFFBodyStateExpected := range []struct {
				StandardFFBodyStateKey StandardFFBodyStateKey
				StandardFFBodyStateRaw []byte
			}{{standardFFBodyStateControl, standardFFBodyStateRaw}, {standardFFBodyStateKey, standardFFBodyStateOtherRaw}} {
				standardFFBodyStateGot, standardFFBodyStateFound, standardFFBodyStateErr := standardFFBodyStateLookup(standardFFBodyStateSource, standardFFBodyStateExpected.StandardFFBodyStateKey)
				if standardFFBodyStateErr != nil || !standardFFBodyStateFound || !standardFFBodyStateBytes.Equal(standardFFBodyStateGot, standardFFBodyStateExpected.StandardFFBodyStateRaw) {
					standardFFBodyStateT.Fatalf("ab54 valid independent control: %v", standardFFBodyStateErr)
				}
			}
		})
	}
	standardFFBodyStateZero := standardFFBodyStateTestKey(0)
	standardFFBodyStateGet(standardFFBodyStateT, standardFFBodyStateZero)
	standardFFBodyStateZero.StandardFFBodyStateSecondary = nil
	if _, standardFFBodyStateFound, standardFFBodyStateErr := StandardFFBodyStateLookup(standardFFBodyStateZero); standardFFBodyStateErr != nil || standardFFBodyStateFound {
		standardFFBodyStateT.Fatal("ab54 absent SB matched zero")
	}
}

func TestAb54InvalidAndDuplicate(standardFFBodyStateT *standardFFBodyStateTesting.T) {
	standardFFBodyStateRaw := standardFFBodyStateGet(standardFFBodyStateT, standardFFBodyStateTestKey(4))
	standardFFBodyStateDuplicate := "[" + string(standardFFBodyStateRaw) + "," + string(standardFFBodyStateRaw) + "]"
	for standardFFBodyStateIndex, standardFFBodyStateSource := range []string{"", "null", "{}", "[", "[null]", "[{}]", "[1]", "[] false", standardFFBodyStateDuplicate,
		"[" + standardFFBodyStateStrings.Replace(string(standardFFBodyStateRaw), `"primary": 255`, `"primary": 255, "primary": 254`, 1) + "]",
		"[" + standardFFBodyStateStrings.Replace(string(standardFFBodyStateRaw), `"primary": 255`, `"primary": 255, "Primary": 254`, 1) + "]",
		"[" + standardFFBodyStateStrings.Replace(string(standardFFBodyStateRaw), `"key": {`, `"key": {}, "key": {`, 1) + "]",
		"[" + standardFFBodyStateStrings.Replace(string(standardFFBodyStateRaw), `"primary": 255,`, ``, 1) + "]",
		"[" + standardFFBodyStateStrings.Replace(string(standardFFBodyStateRaw), `"primary": 255`, `"primary": null`, 1) + "]",
		"[" + standardFFBodyStateStrings.Replace(string(standardFFBodyStateRaw), `"alternative": "none"`, `"alternative": ""`, 1) + "]",
		"[" + standardFFBodyStateStrings.Replace(string(standardFFBodyStateRaw), `"primary": 255`, `"primary": 256`, 1) + "]",
	} {
		standardFFBodyStateT.Run(standardFFBodyStateFmt.Sprintf("ab54-invalid-%d", standardFFBodyStateIndex), func(standardFFBodyStateT *standardFFBodyStateTesting.T) {
			if standardFFBodyStateEntries, standardFFBodyStateIndex, standardFFBodyStateErr := standardFFBodyStateValidate([]byte(standardFFBodyStateSource)); standardFFBodyStateErr == nil || standardFFBodyStateEntries != nil || standardFFBodyStateIndex != nil {
				standardFFBodyStateT.Fatal("ab54 invalid source succeeded")
			}
			if standardFFBodyStateGot, standardFFBodyStateFound, standardFFBodyStateErr := standardFFBodyStateLookup([]byte(standardFFBodyStateSource), standardFFBodyStateTestKey(4)); standardFFBodyStateErr == nil || standardFFBodyStateFound || standardFFBodyStateGot != nil {
				standardFFBodyStateT.Fatal("ab54 invalid lookup succeeded")
			}
		})
	}
	standardFFBodyStateAbsentKey := standardFFBodyStateTestKey(4)
	standardFFBodyStateAbsentKey.StandardFFBodyStateFamily = "absent"
	if _, standardFFBodyStateFound, standardFFBodyStateErr := standardFFBodyStateLookup([]byte(standardFFBodyStateDuplicate), standardFFBodyStateAbsentKey); standardFFBodyStateErr == nil || standardFFBodyStateFound {
		standardFFBodyStateT.Fatal("ab54 invalid catalog was hidden by a miss")
	}
}

func TestAb54RequiredNullableAxes(standardFFBodyStateT *standardFFBodyStateTesting.T) {
	standardFFBodyStateRaw := standardFFBodyStateGet(standardFFBodyStateT, standardFFBodyStateTestKey(4))
	for _, standardFFBodyStateAxis := range []string{"secondary", "operation", "group", "register"} {
		standardFFBodyStateT.Run("ab54-"+standardFFBodyStateAxis, func(standardFFBodyStateT *standardFFBodyStateTesting.T) {
			for _, standardFFBodyStateValue := range []string{"", "null", "0", "1"} {
				standardFFBodyStateT.Run("ab54-value-"+standardFFBodyStateValue, func(standardFFBodyStateT *standardFFBodyStateTesting.T) {
					standardFFBodyStateRow := standardFFBodyStateObject(standardFFBodyStateT, standardFFBodyStateRaw)
					standardFFBodyStateFields := standardFFBodyStateObject(standardFFBodyStateT, standardFFBodyStateRow["key"])
					delete(standardFFBodyStateFields, standardFFBodyStateAxis)
					if standardFFBodyStateValue != "" {
						standardFFBodyStateFields[standardFFBodyStateAxis] = standardFFBodyStateJSON.RawMessage(standardFFBodyStateValue)
					}
					standardFFBodyStateKeyData, standardFFBodyStateErr := standardFFBodyStateJSON.Marshal(standardFFBodyStateFields)
					if standardFFBodyStateErr != nil {
						standardFFBodyStateT.Fatal(standardFFBodyStateErr)
					}
					standardFFBodyStateRow["key"] = standardFFBodyStateKeyData
					standardFFBodyStateSource, standardFFBodyStateErr := standardFFBodyStateJSON.Marshal([]map[string]standardFFBodyStateJSON.RawMessage{standardFFBodyStateRow})
					if standardFFBodyStateErr != nil {
						standardFFBodyStateT.Fatal(standardFFBodyStateErr)
					}
					var standardFFBodyStateKey StandardFFBodyStateKey
					if standardFFBodyStateErr := standardFFBodyStateJSON.Unmarshal(standardFFBodyStateKeyData, &standardFFBodyStateKey); standardFFBodyStateErr != nil {
						standardFFBodyStateT.Fatal(standardFFBodyStateErr)
					}
					standardFFBodyStateEntries, standardFFBodyStateIndex, standardFFBodyStateValidateErr := standardFFBodyStateValidate(standardFFBodyStateSource)
					standardFFBodyStateGot, standardFFBodyStateFound, standardFFBodyStateLookupErr := standardFFBodyStateLookup(standardFFBodyStateSource, standardFFBodyStateKey)
					if standardFFBodyStateValue == "" {
						if standardFFBodyStateValidateErr == nil || standardFFBodyStateEntries != nil || standardFFBodyStateIndex != nil {
							standardFFBodyStateT.Error("ab54 missing nullable axis passed validation")
						}
						if standardFFBodyStateLookupErr == nil || standardFFBodyStateFound || standardFFBodyStateGot != nil {
							standardFFBodyStateT.Error("ab54 missing nullable axis passed lookup")
						}
						return
					}
					if standardFFBodyStateValidateErr != nil || !standardFFBodyStateBytes.Equal(standardFFBodyStateEntries, standardFFBodyStateSource) || len(standardFFBodyStateIndex) != 1 {
						standardFFBodyStateT.Fatalf("ab54 explicit nullable axis failed validation: %v", standardFFBodyStateValidateErr)
					}
					if standardFFBodyStateLookupErr != nil || !standardFFBodyStateFound || !standardFFBodyStateBytes.Equal(standardFFBodyStateGot, standardFFBodyStateSource[1:len(standardFFBodyStateSource)-1]) {
						standardFFBodyStateT.Fatalf("ab54 explicit nullable axis failed exact lookup: %v", standardFFBodyStateLookupErr)
					}
				})
			}
		})
	}
}

func TestAb54IndependentStorage(standardFFBodyStateT *standardFFBodyStateTesting.T) {
	standardFFBodyStateBefore := standardFFBodyStateGet(standardFFBodyStateT, standardFFBodyStateTestKey(4))
	standardFFBodyStateEntries, standardFFBodyStateErr := StandardFFBodyStateEntries()
	if standardFFBodyStateErr != nil {
		standardFFBodyStateT.Fatal(standardFFBodyStateErr)
	}
	standardFFBodyStateExpected := append([]byte(nil), standardFFBodyStateEntries...)
	standardFFBodyStateEntries[0] = '!'
	standardFFBodyStateReturned := standardFFBodyStateGet(standardFFBodyStateT, standardFFBodyStateTestKey(4))
	standardFFBodyStateReturned[0] = '!'
	standardFFBodyStateAgain, standardFFBodyStateErr := StandardFFBodyStateEntries()
	if standardFFBodyStateErr != nil || !standardFFBodyStateBytes.Equal(standardFFBodyStateAgain, standardFFBodyStateExpected) || !standardFFBodyStateBytes.Equal(standardFFBodyStateBefore, standardFFBodyStateGet(standardFFBodyStateT, standardFFBodyStateTestKey(4))) {
		standardFFBodyStateT.Fatal("ab54 returned mutation changed later API results")
	}
	standardFFBodyStateSource := append([]byte(nil), standardFFBodyStateExpected...)
	standardFFBodyStateDetached, standardFFBodyStateIndex, standardFFBodyStateErr := standardFFBodyStateValidate(standardFFBodyStateSource)
	if standardFFBodyStateErr != nil {
		standardFFBodyStateT.Fatal(standardFFBodyStateErr)
	}
	for standardFFBodyStateOffset := range standardFFBodyStateSource {
		standardFFBodyStateSource[standardFFBodyStateOffset] = '!'
	}
	if !standardFFBodyStateBytes.Equal(standardFFBodyStateDetached, standardFFBodyStateExpected) {
		standardFFBodyStateT.Fatal("ab54 source aliases entries")
	}
	standardFFBodyStateIdentity, _ := standardFFBodyStateJSON.Marshal(standardFFBodyStateTestKey(4))
	if !standardFFBodyStateBytes.Equal(standardFFBodyStateIndex[string(standardFFBodyStateIdentity)], standardFFBodyStateBefore) {
		standardFFBodyStateT.Fatal("ab54 source aliases row")
	}
	standardFFBodyStateIndex[string(standardFFBodyStateIdentity)][0] = '!'
	if !standardFFBodyStateBytes.Equal(standardFFBodyStateDetached, standardFFBodyStateExpected) {
		standardFFBodyStateT.Fatal("ab54 entries alias indexed row")
	}
}

func TestAb54MetadataPreservation(standardFFBodyStateT *standardFFBodyStateTesting.T) {
	standardFFBodyStateRaw := standardFFBodyStateGet(standardFFBodyStateT, standardFFBodyStateTestKey(4))
	for _, standardFFBodyStateGrade := range []string{"Proven", "Hypothesis", "Unknown"} {
		for _, standardFFBodyStateWidth := range []string{"null", "0", "1", "3"} {
			standardFFBodyStateRow := standardFFBodyStateObject(standardFFBodyStateT, standardFFBodyStateRaw)
			standardFFBodyStateRow["grade"] = standardFFBodyStateJSON.RawMessage(`"` + standardFFBodyStateGrade + `"`)
			standardFFBodyStateRow["minimum_overhead_width"] = standardFFBodyStateJSON.RawMessage(standardFFBodyStateWidth)
			standardFFBodyStateRow["ab54_unresolved"] = standardFFBodyStateJSON.RawMessage(`{"ab54_raw":[0,255],"ab54_grade":"Unknown"}`)
			standardFFBodyStateChanged, _ := standardFFBodyStateJSON.Marshal(standardFFBodyStateRow)
			standardFFBodyStateSource := append(append([]byte{'['}, standardFFBodyStateChanged...), ']')
			standardFFBodyStateGot, standardFFBodyStateFound, standardFFBodyStateErr := standardFFBodyStateLookup(standardFFBodyStateSource, standardFFBodyStateTestKey(4))
			if standardFFBodyStateErr != nil || !standardFFBodyStateFound || !standardFFBodyStateBytes.Equal(standardFFBodyStateGot, standardFFBodyStateChanged) {
				standardFFBodyStateT.Fatal("ab54 width, grade or unresolved raw metadata lost")
			}
		}
	}
}

func TestAb54OptionalValues(standardFFBodyStateT *standardFFBodyStateTesting.T) {
	for _, standardFFBodyStateAxis := range []string{"ab54-secondary", "ab54-operation", "ab54-group", "ab54-register"} {
		standardFFBodyStateT.Run(standardFFBodyStateAxis, func(standardFFBodyStateT *standardFFBodyStateTesting.T) {
			standardFFBodyStateKeys := []StandardFFBodyStateKey{standardFFBodyStateTestKey(4), standardFFBodyStateTestKey(4)}
			standardFFBodyStateRows := make([]standardFFBodyStateJSON.RawMessage, 2)
			for standardFFBodyStateIndex := range standardFFBodyStateKeys {
				switch standardFFBodyStateAxis {
				case "ab54-secondary":
					standardFFBodyStateKeys[standardFFBodyStateIndex].StandardFFBodyStateSecondary = standardFFBodyStateByte(uint8(standardFFBodyStateIndex))
				case "ab54-operation":
					standardFFBodyStateKeys[standardFFBodyStateIndex].StandardFFBodyStateOperation = standardFFBodyStateByte(uint8(standardFFBodyStateIndex))
				case "ab54-group":
					standardFFBodyStateKeys[standardFFBodyStateIndex].StandardFFBodyStateGroup = standardFFBodyStateWord(uint16(standardFFBodyStateIndex))
				case "ab54-register":
					standardFFBodyStateKeys[standardFFBodyStateIndex].StandardFFBodyStateRegister = standardFFBodyStateWord(uint16(standardFFBodyStateIndex))
				}
				standardFFBodyStateRow := standardFFBodyStateObject(standardFFBodyStateT, standardFFBodyStateGet(standardFFBodyStateT, standardFFBodyStateTestKey(4)))
				standardFFBodyStateRow["key"], _ = standardFFBodyStateJSON.Marshal(standardFFBodyStateKeys[standardFFBodyStateIndex])
				standardFFBodyStateRow["id"] = standardFFBodyStateJSON.RawMessage(standardFFBodyStateFmt.Sprintf("\"ab54-synthetic-optional-%d\"", standardFFBodyStateIndex))
				standardFFBodyStateRows[standardFFBodyStateIndex], _ = standardFFBodyStateJSON.Marshal(standardFFBodyStateRow)
			}
			standardFFBodyStateSource, _ := standardFFBodyStateJSON.Marshal(standardFFBodyStateRows)
			for standardFFBodyStateIndex, standardFFBodyStateKey := range standardFFBodyStateKeys {
				standardFFBodyStateGot, standardFFBodyStateFound, standardFFBodyStateErr := standardFFBodyStateLookup(standardFFBodyStateSource, standardFFBodyStateKey)
				var standardFFBodyStateGotValue, standardFFBodyStateWantValue any
				if standardFFBodyStateErr != nil || !standardFFBodyStateFound {
					standardFFBodyStateT.Fatalf("ab54 explicit zero and nonzero collided: %v", standardFFBodyStateErr)
				}
				if standardFFBodyStateErr := standardFFBodyStateJSON.Unmarshal(standardFFBodyStateGot, &standardFFBodyStateGotValue); standardFFBodyStateErr != nil {
					standardFFBodyStateT.Fatal(standardFFBodyStateErr)
				}
				if standardFFBodyStateErr := standardFFBodyStateJSON.Unmarshal(standardFFBodyStateRows[standardFFBodyStateIndex], &standardFFBodyStateWantValue); standardFFBodyStateErr != nil {
					standardFFBodyStateT.Fatal(standardFFBodyStateErr)
				}
				if !standardFFBodyStateReflect.DeepEqual(standardFFBodyStateGotValue, standardFFBodyStateWantValue) {
					standardFFBodyStateT.Fatal("ab54 wrong optional value selected")
				}
			}
		})
	}
}
