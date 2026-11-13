package vaillant

import (
	b509IdentityShapeBytes "bytes"
	b509IdentityShapeJSON "encoding/json"
	b509IdentityShapeReflect "reflect"
	b509IdentityShapeTesting "testing"
)

func b509IdentityShapeNumber(b509IdentityShapeValue uint16) *uint16 { return &b509IdentityShapeValue }

func b509IdentityShapeTestKey(b509IdentityShapeOperation uint16, b509IdentityShapeDirection, b509IdentityShapeField, b509IdentityShapeAlternative string) B509IdentityShapeKey {
	return B509IdentityShapeKey{B509IdentityShapeFamily: "vaillant", B509IdentityShapeVariant: "B509", B509IdentityShapePrimary: b509IdentityShapeNumber(181), B509IdentityShapeSecondary: b509IdentityShapeNumber(9), B509IdentityShapeOperation: b509IdentityShapeNumber(b509IdentityShapeOperation), B509IdentityShapeDirection: b509IdentityShapeDirection, B509IdentityShapeField: b509IdentityShapeField, B509IdentityShapeAlternative: b509IdentityShapeAlternative, B509IdentityShapeInstanceScope: "none"}
}

func b509IdentityShapeDecode(b509IdentityShapeT *b509IdentityShapeTesting.T, b509IdentityShapeRaw []byte) map[string]b509IdentityShapeJSON.RawMessage {
	b509IdentityShapeT.Helper()
	var b509IdentityShapeRow map[string]b509IdentityShapeJSON.RawMessage
	if b509IdentityShapeError := b509IdentityShapeJSON.Unmarshal(b509IdentityShapeRaw, &b509IdentityShapeRow); b509IdentityShapeError != nil {
		b509IdentityShapeT.Fatal(b509IdentityShapeError)
	}
	return b509IdentityShapeRow
}

func b509IdentityShapeEqualJSON(b509IdentityShapeT *b509IdentityShapeTesting.T, b509IdentityShapeGot []byte, b509IdentityShapeWant string) {
	b509IdentityShapeT.Helper()
	var b509IdentityShapeActual, b509IdentityShapeExpected any
	if b509IdentityShapeError := b509IdentityShapeJSON.Unmarshal(b509IdentityShapeGot, &b509IdentityShapeActual); b509IdentityShapeError != nil {
		b509IdentityShapeT.Fatal(b509IdentityShapeError)
	}
	if b509IdentityShapeError := b509IdentityShapeJSON.Unmarshal([]byte(b509IdentityShapeWant), &b509IdentityShapeExpected); b509IdentityShapeError != nil {
		b509IdentityShapeT.Fatal(b509IdentityShapeError)
	}
	if !b509IdentityShapeReflect.DeepEqual(b509IdentityShapeActual, b509IdentityShapeExpected) {
		b509IdentityShapeT.Fatalf("ab46 metadata: got %s; want %s", b509IdentityShapeGot, b509IdentityShapeWant)
	}
}

func b509IdentityShapeGet(b509IdentityShapeT *b509IdentityShapeTesting.T, b509IdentityShapeKey B509IdentityShapeKey) map[string]b509IdentityShapeJSON.RawMessage {
	b509IdentityShapeT.Helper()
	b509IdentityShapeRaw, b509IdentityShapeFound, b509IdentityShapeError := B509IdentityShapeLookup(b509IdentityShapeKey)
	if b509IdentityShapeError != nil || !b509IdentityShapeFound {
		b509IdentityShapeT.Fatalf("ab46 missing exact row: %v", b509IdentityShapeError)
	}
	return b509IdentityShapeDecode(b509IdentityShapeT, b509IdentityShapeRaw)
}

func TestAb46EntriesAndWholeRows(b509IdentityShapeT *b509IdentityShapeTesting.T) {
	b509IdentityShapeRaw, b509IdentityShapeError := B509IdentityShapeEntries()
	if b509IdentityShapeError != nil {
		b509IdentityShapeT.Fatal(b509IdentityShapeError)
	}
	var b509IdentityShapeRows []b509IdentityShapeJSON.RawMessage
	if b509IdentityShapeError = b509IdentityShapeJSON.Unmarshal(b509IdentityShapeRaw, &b509IdentityShapeRows); b509IdentityShapeError != nil || len(b509IdentityShapeRows) != 6 {
		b509IdentityShapeT.Fatalf("ab46 finite array: count %d error %v", len(b509IdentityShapeRows), b509IdentityShapeError)
	}
	b509IdentityShapeKeys := []B509IdentityShapeKey{
		b509IdentityShapeTestKey(83, "read", "hardware_identifier", "packed_decimal_candidate"),
		b509IdentityShapeTestKey(83, "write", "hardware_identifier", "packed_decimal_candidate"),
		b509IdentityShapeTestKey(84, "Unknown", "identity_content", "link_address_six"),
		b509IdentityShapeTestKey(84, "Unknown", "identity_content", "conflicting_seven"),
		b509IdentityShapeTestKey(84, "Unknown", "identity_content", "segmented_characters"),
		b509IdentityShapeTestKey(84, "Unknown", "identity_state", "none"),
	}
	for b509IdentityShapeIndex, b509IdentityShapeKey := range b509IdentityShapeKeys {
		b509IdentityShapeGot, b509IdentityShapeFound, b509IdentityShapeError := B509IdentityShapeLookup(b509IdentityShapeKey)
		if b509IdentityShapeError != nil || !b509IdentityShapeFound || !b509IdentityShapeBytes.Equal(b509IdentityShapeGot, b509IdentityShapeRows[b509IdentityShapeIndex]) {
			b509IdentityShapeT.Fatal("ab46 lookup must return the entire original row")
		}
		b509IdentityShapeRow := b509IdentityShapeDecode(b509IdentityShapeT, b509IdentityShapeGot)
		var b509IdentityShapeActualKey B509IdentityShapeKey
		if b509IdentityShapeError = b509IdentityShapeJSON.Unmarshal(b509IdentityShapeRow["key"], &b509IdentityShapeActualKey); b509IdentityShapeError != nil || !b509IdentityShapeReflect.DeepEqual(b509IdentityShapeActualKey, b509IdentityShapeKey) {
			b509IdentityShapeT.Fatal("ab46 independent key contract")
		}
		b509IdentityShapeEqualJSON(b509IdentityShapeT, b509IdentityShapeRow["grade"], `"Hypothesis"`)
		for _, b509IdentityShapeAxis := range []string{"qualification", "presence", "freshness", "execution"} {
			b509IdentityShapeEqualJSON(b509IdentityShapeT, b509IdentityShapeRow[b509IdentityShapeAxis], `"Unknown"`)
		}
		b509IdentityShapeEqualJSON(b509IdentityShapeT, b509IdentityShapeRow["unresolved"], `{"field_positions":"Unknown","segment_selectors":"Unknown","padding":"Unknown","termination":"Unknown","transport_correlation":"Unknown","completion_predicate":"Unknown"}`)
		b509IdentityShapeEqualJSON(b509IdentityShapeT, b509IdentityShapeRow["native_request_admission"], `false`)
		b509IdentityShapeEqualJSON(b509IdentityShapeT, b509IdentityShapeRow["complete_identity_inference"], `false`)
		b509IdentityShapeEqualJSON(b509IdentityShapeT, b509IdentityShapeRow["raw_fallback"], `"preserve all original bytes and unaffected fields"`)
	}
}

func TestAb46HardwareCandidate(b509IdentityShapeT *b509IdentityShapeTesting.T) {
	for _, b509IdentityShapeDirection := range []string{"read", "write"} {
		b509IdentityShapeRow := b509IdentityShapeGet(b509IdentityShapeT, b509IdentityShapeTestKey(83, b509IdentityShapeDirection, "hardware_identifier", "packed_decimal_candidate"))
		b509IdentityShapeEqualJSON(b509IdentityShapeT, b509IdentityShapeRow["width_bytes"], `6`)
		b509IdentityShapeEqualJSON(b509IdentityShapeT, b509IdentityShapeRow["raw_form"], `"six opaque bytes"`)
		b509IdentityShapeEqualJSON(b509IdentityShapeT, b509IdentityShapeRow["flags"], `{"grade":"Unknown","bits":"Unknown"}`)
		b509IdentityShapeEqualJSON(b509IdentityShapeT, b509IdentityShapeRow["candidate"], `{"grade":"Hypothesis","encoding":"packed_decimal_pairs","nibble_min":0,"nibble_max":9,"invalid_nibble":"candidate absent; preserve all raw bytes","security_identity":false}`)
		b509IdentityShapeEqualJSON(b509IdentityShapeT, b509IdentityShapeRow["synthetic_cases"], `[
			{"id":"ab46_decimal_pairs","raw":[18,52,86,120,144,18],"candidate_decimal_pairs":[12,34,56,78,90,12],"candidate_valid":true},
			{"id":"ab46_invalid_nibble","raw":[250,52,86,120,144,18],"candidate_decimal_pairs":null,"candidate_valid":false}
		]`)
	}
}

func TestAb46SegmentedAlternatives(b509IdentityShapeT *b509IdentityShapeTesting.T) {
	for _, b509IdentityShapeCase := range []struct {
		b509IdentityShapeAlternative string
		b509IdentityShapeWidth       string
		b509IdentityShapeRaw         string
		b509IdentityShapeForm        string
	}{
		{"link_address_six", `6`, `[1,2,3,4,5,6]`, `"six opaque link-address bytes"`},
		{"conflicting_seven", `7`, `[1,2,3,4,5,6,7]`, `"seven opaque bytes; extra byte uninterpreted"`},
		{"segmented_characters", `null`, `[1,2,3]`, `"raw character segments; boundaries unresolved"`},
	} {
		b509IdentityShapeRow := b509IdentityShapeGet(b509IdentityShapeT, b509IdentityShapeTestKey(84, "Unknown", "identity_content", b509IdentityShapeCase.b509IdentityShapeAlternative))
		b509IdentityShapeEqualJSON(b509IdentityShapeT, b509IdentityShapeRow["width_bytes"], b509IdentityShapeCase.b509IdentityShapeWidth)
		b509IdentityShapeEqualJSON(b509IdentityShapeT, b509IdentityShapeRow["synthetic_raw"], b509IdentityShapeCase.b509IdentityShapeRaw)
		b509IdentityShapeEqualJSON(b509IdentityShapeT, b509IdentityShapeRow["raw_form"], b509IdentityShapeCase.b509IdentityShapeForm)
		b509IdentityShapeEqualJSON(b509IdentityShapeT, b509IdentityShapeRow["selection_by_length"], `false`)
		b509IdentityShapeEqualJSON(b509IdentityShapeT, b509IdentityShapeRow["interpretation"], `"Unknown"`)
	}
}

func TestAb46FeedbackAndPartialRetention(b509IdentityShapeT *b509IdentityShapeTesting.T) {
	b509IdentityShapeRow := b509IdentityShapeGet(b509IdentityShapeT, b509IdentityShapeTestKey(84, "Unknown", "identity_state", "none"))
	b509IdentityShapeEqualJSON(b509IdentityShapeT, b509IdentityShapeRow["properties"], `["feedback_success","status_success","content_extent","segment_completeness","complete_identity"]`)
	b509IdentityShapeEqualJSON(b509IdentityShapeT, b509IdentityShapeRow["synthetic_cases"], `[
		{"id":"ab46_empty_success","feedback_success":true,"status_success":true,"content_extent":"empty","segments":[],"segment_completeness":false,"complete_identity":false},
		{"id":"ab46_partial_success","feedback_success":true,"status_success":true,"content_extent":"partial","segments":[{"raw":[65,66],"valid":true},{"raw":[250],"valid":false}],"segment_completeness":false,"complete_identity":false},
		{"id":"ab46_status_failure","feedback_success":true,"status_success":false,"content_extent":"synthetic_full","segments":[{"raw":[65,66],"valid":true},{"raw":[67,68],"valid":true}],"segment_completeness":true,"complete_identity":false},
		{"id":"ab46_segment_failure","feedback_success":true,"status_success":true,"content_extent":"partial","previous_segments":[[65,66],[67,68]],"segments":[{"raw":[69,70],"valid":true},{"raw":[250],"valid":false}],"retained_segments":[[69,70],[67,68]],"segment_completeness":false,"complete_identity":false}
	]`)
}

func TestAb46EveryKeyAxis(b509IdentityShapeT *b509IdentityShapeTesting.T) {
	b509IdentityShapeControl := b509IdentityShapeTestKey(83, "read", "hardware_identifier", "packed_decimal_candidate")
	b509IdentityShapeType := b509IdentityShapeReflect.TypeOf(b509IdentityShapeControl)
	for b509IdentityShapeIndex := 0; b509IdentityShapeIndex < b509IdentityShapeType.NumField(); b509IdentityShapeIndex++ {
		b509IdentityShapeT.Run("ab46_"+b509IdentityShapeType.Field(b509IdentityShapeIndex).Name, func(b509IdentityShapeT *b509IdentityShapeTesting.T) {
			if _, b509IdentityShapeFound, b509IdentityShapeError := B509IdentityShapeLookup(b509IdentityShapeControl); b509IdentityShapeError != nil || !b509IdentityShapeFound {
				b509IdentityShapeT.Fatal("ab46 missing positive control")
			}
			b509IdentityShapeChanged := b509IdentityShapeControl
			b509IdentityShapeAxis := b509IdentityShapeReflect.ValueOf(&b509IdentityShapeChanged).Elem().Field(b509IdentityShapeIndex)
			if b509IdentityShapeAxis.Kind() == b509IdentityShapeReflect.String {
				b509IdentityShapeAxis.SetString("ab46_distinct")
			} else {
				b509IdentityShapeAxis.Set(b509IdentityShapeReflect.ValueOf(b509IdentityShapeNumber(0)))
			}
			if _, b509IdentityShapeFound, b509IdentityShapeError := B509IdentityShapeLookup(b509IdentityShapeChanged); b509IdentityShapeError != nil || b509IdentityShapeFound {
				b509IdentityShapeT.Fatal("ab46 changed axis matched")
			}
			if b509IdentityShapeAxis.Kind() == b509IdentityShapeReflect.String {
				for _, b509IdentityShapeToken := range []string{"", "*", "none"} {
					b509IdentityShapeAxis.SetString(b509IdentityShapeToken)
					if b509IdentityShapeReflect.DeepEqual(b509IdentityShapeChanged, b509IdentityShapeControl) {
						continue
					}
					if _, b509IdentityShapeFound, b509IdentityShapeError := B509IdentityShapeLookup(b509IdentityShapeChanged); b509IdentityShapeError != nil || b509IdentityShapeFound {
						b509IdentityShapeT.Fatal("ab46 string token became a wildcard")
					}
				}
			}
			if b509IdentityShapeAxis.Kind() != b509IdentityShapeReflect.Ptr {
				return
			}
			b509IdentityShapeAxis.Set(b509IdentityShapeReflect.ValueOf(b509IdentityShapeNumber(65535)))
			if _, b509IdentityShapeFound, b509IdentityShapeError := B509IdentityShapeLookup(b509IdentityShapeChanged); b509IdentityShapeError != nil || b509IdentityShapeFound {
				b509IdentityShapeT.Fatal("ab46 numeric axis ignored")
			}
			if b509IdentityShapeType.Field(b509IdentityShapeIndex).Name == "B509IdentityShapePrimary" {
				return
			}
			b509IdentityShapeAxis.Set(b509IdentityShapeReflect.Zero(b509IdentityShapeAxis.Type()))
			b509IdentityShapeAbsent := b509IdentityShapeChanged
			b509IdentityShapeAxis.Set(b509IdentityShapeReflect.ValueOf(b509IdentityShapeNumber(0)))
			b509IdentityShapeSource, b509IdentityShapeError := b509IdentityShapeJSON.Marshal([]any{map[string]any{"key": b509IdentityShapeChanged, "grade": "Unknown", "ab46_extension": []int{1, 2}}})
			if b509IdentityShapeError != nil {
				b509IdentityShapeT.Fatal(b509IdentityShapeError)
			}
			if _, b509IdentityShapeFound, b509IdentityShapeError := b509IdentityShapeLookup(b509IdentityShapeSource, b509IdentityShapeChanged); b509IdentityShapeError != nil || !b509IdentityShapeFound {
				b509IdentityShapeT.Fatal("ab46 zero-valued key control failed")
			}
			if _, b509IdentityShapeFound, b509IdentityShapeError := b509IdentityShapeLookup(b509IdentityShapeSource, b509IdentityShapeAbsent); b509IdentityShapeError != nil || b509IdentityShapeFound {
				b509IdentityShapeT.Fatal("ab46 absence matched explicit zero")
			}
		})
	}
}

func TestAb46InvalidAndDuplicateData(b509IdentityShapeT *b509IdentityShapeTesting.T) {
	b509IdentityShapeOriginal := b509IdentityShapeData
	defer func() { b509IdentityShapeData = b509IdentityShapeOriginal }()
	b509IdentityShapeKey := b509IdentityShapeTestKey(83, "read", "hardware_identifier", "packed_decimal_candidate")
	b509IdentityShapeRow, _, _ := B509IdentityShapeLookup(b509IdentityShapeKey)
	b509IdentityShapeDuplicate := append([]byte{'['}, b509IdentityShapeRow...)
	b509IdentityShapeDuplicate = append(b509IdentityShapeDuplicate, ',')
	b509IdentityShapeDuplicate = append(b509IdentityShapeDuplicate, b509IdentityShapeRow...)
	b509IdentityShapeDuplicate = append(b509IdentityShapeDuplicate, ']')
	b509IdentityShapeChangedRow := b509IdentityShapeDecode(b509IdentityShapeT, b509IdentityShapeRow)
	b509IdentityShapeChangedRow["grade"] = b509IdentityShapeJSON.RawMessage(`"Unknown"`)
	b509IdentityShapeChangedRow["width_bytes"] = b509IdentityShapeJSON.RawMessage(`7`)
	b509IdentityShapeChangedRaw, b509IdentityShapeError := b509IdentityShapeJSON.Marshal(b509IdentityShapeChangedRow)
	if b509IdentityShapeError != nil {
		b509IdentityShapeT.Fatal(b509IdentityShapeError)
	}
	for _, b509IdentityShapeBad := range []string{``, `{}`, `null`, `[null]`, `[1]`, `[{}]`, `[{"key":null}]`, `[{"key":{"family":"vaillant"}}]`, `[`, string(b509IdentityShapeDuplicate), "[" + string(b509IdentityShapeRow) + "," + string(b509IdentityShapeChangedRaw) + "]", "[" + string(b509IdentityShapeRow) + ",null]"} {
		b509IdentityShapeData = b509IdentityShapeBad
		if b509IdentityShapeGot, b509IdentityShapeError := B509IdentityShapeEntries(); b509IdentityShapeError == nil || b509IdentityShapeGot != nil {
			b509IdentityShapeT.Fatal("ab46 invalid entries succeeded")
		}
		if b509IdentityShapeGot, b509IdentityShapeFound, b509IdentityShapeError := B509IdentityShapeLookup(b509IdentityShapeKey); b509IdentityShapeError == nil || b509IdentityShapeFound || b509IdentityShapeGot != nil {
			b509IdentityShapeT.Fatal("ab46 invalid lookup succeeded, including an early matching row")
		}
	}
	b509IdentityShapeData = `[]`
	if b509IdentityShapeGot, b509IdentityShapeError := B509IdentityShapeEntries(); b509IdentityShapeError != nil || string(b509IdentityShapeGot) != `[]` {
		b509IdentityShapeT.Fatal("ab46 empty array control")
	}
	if b509IdentityShapeGot, b509IdentityShapeFound, b509IdentityShapeError := B509IdentityShapeLookup(b509IdentityShapeKey); b509IdentityShapeError != nil || b509IdentityShapeFound || b509IdentityShapeGot != nil {
		b509IdentityShapeT.Fatal("ab46 empty array lookup")
	}
}

func TestAb46IndependentStorageAndOpaqueMetadata(b509IdentityShapeT *b509IdentityShapeTesting.T) {
	b509IdentityShapeKey := b509IdentityShapeTestKey(83, "read", "hardware_identifier", "packed_decimal_candidate")
	b509IdentityShapeBefore, b509IdentityShapeError := B509IdentityShapeEntries()
	if b509IdentityShapeError != nil {
		b509IdentityShapeT.Fatal(b509IdentityShapeError)
	}
	b509IdentityShapeEntries, _ := B509IdentityShapeEntries()
	b509IdentityShapeEntries[0] = '!'
	b509IdentityShapeLookupRaw, b509IdentityShapeFound, b509IdentityShapeError := B509IdentityShapeLookup(b509IdentityShapeKey)
	if b509IdentityShapeError != nil || !b509IdentityShapeFound {
		b509IdentityShapeT.Fatal("ab46 lookup control")
	}
	b509IdentityShapeLookupRaw[0] = '!'
	b509IdentityShapeAfter, _ := B509IdentityShapeEntries()
	if !b509IdentityShapeBytes.Equal(b509IdentityShapeBefore, b509IdentityShapeAfter) {
		b509IdentityShapeT.Fatal("ab46 returned memory changed catalog")
	}
	b509IdentityShapeAgain, _, _ := B509IdentityShapeLookup(b509IdentityShapeKey)
	if len(b509IdentityShapeAgain) == 0 || b509IdentityShapeAgain[0] != '{' {
		b509IdentityShapeT.Fatal("ab46 lookup shared storage")
	}
	for _, b509IdentityShapeGrade := range []string{"Proven", "Hypothesis", "Unknown"} {
		b509IdentityShapeSource, b509IdentityShapeError := b509IdentityShapeJSON.Marshal([]any{map[string]any{"key": b509IdentityShapeKey, "grade": b509IdentityShapeGrade, "width_bytes": 0, "ab46_extension": map[string]any{"raw": []int{0, 250}, "unresolved": "Unknown"}}})
		if b509IdentityShapeError != nil {
			b509IdentityShapeT.Fatal(b509IdentityShapeError)
		}
		b509IdentityShapeRows, b509IdentityShapeKeys, b509IdentityShapeError := b509IdentityShapeValidate(b509IdentityShapeSource)
		if b509IdentityShapeError != nil || len(b509IdentityShapeRows) != 1 || !b509IdentityShapeReflect.DeepEqual(b509IdentityShapeKeys[0], b509IdentityShapeKey) {
			b509IdentityShapeT.Fatal("ab46 private seam control")
		}
		b509IdentityShapeSource[0] = '!'
		b509IdentityShapeEqualJSON(b509IdentityShapeT, b509IdentityShapeDecode(b509IdentityShapeT, b509IdentityShapeRows[0])["grade"], `"`+b509IdentityShapeGrade+`"`)
		b509IdentityShapeCopy := append([]byte{'['}, b509IdentityShapeRows[0]...)
		b509IdentityShapeCopy = append(b509IdentityShapeCopy, ']')
		b509IdentityShapeGot, b509IdentityShapeFound, b509IdentityShapeError := b509IdentityShapeLookup(b509IdentityShapeCopy, b509IdentityShapeKey)
		if b509IdentityShapeError != nil || !b509IdentityShapeFound || !b509IdentityShapeBytes.Equal(b509IdentityShapeGot, b509IdentityShapeRows[0]) {
			b509IdentityShapeT.Fatal("ab46 whole opaque metadata row changed")
		}
		b509IdentityShapeCopy[1] = '!'
		b509IdentityShapeEqualJSON(b509IdentityShapeT, b509IdentityShapeDecode(b509IdentityShapeT, b509IdentityShapeGot)["width_bytes"], `0`)
		b509IdentityShapeEqualJSON(b509IdentityShapeT, b509IdentityShapeDecode(b509IdentityShapeT, b509IdentityShapeGot)["ab46_extension"], `{"raw":[0,250],"unresolved":"Unknown"}`)
		b509IdentityShapeGot[0] = '!'
		if b509IdentityShapeRows[0][0] != '{' {
			b509IdentityShapeT.Fatal("ab46 source row mutated by lookup result")
		}
	}
}
