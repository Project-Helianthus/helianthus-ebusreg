package vaillant

import (
	ab46Bytes "bytes"
	ab46JSON "encoding/json"
	ab46Reflect "reflect"
	ab46Testing "testing"
)

func ab46Number(ab46Value uint16) *uint16 { return &ab46Value }

func ab46TestKey(ab46Operation uint16, ab46Direction, ab46Field, ab46Alternative string) Ab46Key {
	return Ab46Key{Ab46Family: "vaillant", Ab46Variant: "B509", Ab46Primary: ab46Number(181), Ab46Secondary: ab46Number(9), Ab46Operation: ab46Number(ab46Operation), Ab46Direction: ab46Direction, Ab46Field: ab46Field, Ab46Alternative: ab46Alternative, Ab46InstanceScope: "none"}
}

func ab46Decode(ab46T *ab46Testing.T, ab46Raw []byte) map[string]ab46JSON.RawMessage {
	ab46T.Helper()
	var ab46Row map[string]ab46JSON.RawMessage
	if ab46Error := ab46JSON.Unmarshal(ab46Raw, &ab46Row); ab46Error != nil {
		ab46T.Fatal(ab46Error)
	}
	return ab46Row
}

func ab46EqualJSON(ab46T *ab46Testing.T, ab46Got []byte, ab46Want string) {
	ab46T.Helper()
	var ab46Actual, ab46Expected any
	if ab46Error := ab46JSON.Unmarshal(ab46Got, &ab46Actual); ab46Error != nil {
		ab46T.Fatal(ab46Error)
	}
	if ab46Error := ab46JSON.Unmarshal([]byte(ab46Want), &ab46Expected); ab46Error != nil {
		ab46T.Fatal(ab46Error)
	}
	if !ab46Reflect.DeepEqual(ab46Actual, ab46Expected) {
		ab46T.Fatalf("ab46 metadata: got %s; want %s", ab46Got, ab46Want)
	}
}

func ab46Get(ab46T *ab46Testing.T, ab46Key Ab46Key) map[string]ab46JSON.RawMessage {
	ab46T.Helper()
	ab46Raw, ab46Found, ab46Error := Ab46Lookup(ab46Key)
	if ab46Error != nil || !ab46Found {
		ab46T.Fatalf("ab46 missing exact row: %v", ab46Error)
	}
	return ab46Decode(ab46T, ab46Raw)
}

func TestAb46EntriesAndWholeRows(ab46T *ab46Testing.T) {
	ab46Raw, ab46Error := Ab46Entries()
	if ab46Error != nil {
		ab46T.Fatal(ab46Error)
	}
	var ab46Rows []ab46JSON.RawMessage
	if ab46Error = ab46JSON.Unmarshal(ab46Raw, &ab46Rows); ab46Error != nil || len(ab46Rows) != 6 {
		ab46T.Fatalf("ab46 finite array: count %d error %v", len(ab46Rows), ab46Error)
	}
	ab46Keys := []Ab46Key{
		ab46TestKey(83, "read", "hardware_identifier", "packed_decimal_candidate"),
		ab46TestKey(83, "write", "hardware_identifier", "packed_decimal_candidate"),
		ab46TestKey(84, "Unknown", "identity_content", "link_address_six"),
		ab46TestKey(84, "Unknown", "identity_content", "conflicting_seven"),
		ab46TestKey(84, "Unknown", "identity_content", "segmented_characters"),
		ab46TestKey(84, "Unknown", "identity_state", "none"),
	}
	for ab46Index, ab46Key := range ab46Keys {
		ab46Got, ab46Found, ab46Error := Ab46Lookup(ab46Key)
		if ab46Error != nil || !ab46Found || !ab46Bytes.Equal(ab46Got, ab46Rows[ab46Index]) {
			ab46T.Fatal("ab46 lookup must return the entire original row")
		}
		ab46Row := ab46Decode(ab46T, ab46Got)
		var ab46ActualKey Ab46Key
		if ab46Error = ab46JSON.Unmarshal(ab46Row["key"], &ab46ActualKey); ab46Error != nil || !ab46Reflect.DeepEqual(ab46ActualKey, ab46Key) {
			ab46T.Fatal("ab46 independent key contract")
		}
		ab46EqualJSON(ab46T, ab46Row["grade"], `"Hypothesis"`)
		for _, ab46Axis := range []string{"qualification", "presence", "freshness", "execution"} {
			ab46EqualJSON(ab46T, ab46Row[ab46Axis], `"Unknown"`)
		}
		ab46EqualJSON(ab46T, ab46Row["unresolved"], `{"field_positions":"Unknown","segment_selectors":"Unknown","padding":"Unknown","termination":"Unknown","transport_correlation":"Unknown","completion_predicate":"Unknown"}`)
		ab46EqualJSON(ab46T, ab46Row["native_request_admission"], `false`)
		ab46EqualJSON(ab46T, ab46Row["complete_identity_inference"], `false`)
		ab46EqualJSON(ab46T, ab46Row["raw_fallback"], `"preserve all original bytes and unaffected fields"`)
	}
}

func TestAb46HardwareCandidate(ab46T *ab46Testing.T) {
	for _, ab46Direction := range []string{"read", "write"} {
		ab46Row := ab46Get(ab46T, ab46TestKey(83, ab46Direction, "hardware_identifier", "packed_decimal_candidate"))
		ab46EqualJSON(ab46T, ab46Row["width_bytes"], `6`)
		ab46EqualJSON(ab46T, ab46Row["raw_form"], `"six opaque bytes"`)
		ab46EqualJSON(ab46T, ab46Row["flags"], `{"grade":"Unknown","bits":"Unknown"}`)
		ab46EqualJSON(ab46T, ab46Row["candidate"], `{"grade":"Hypothesis","encoding":"packed_decimal_pairs","nibble_min":0,"nibble_max":9,"invalid_nibble":"candidate absent; preserve all raw bytes","security_identity":false}`)
		ab46EqualJSON(ab46T, ab46Row["synthetic_cases"], `[
			{"id":"ab46_decimal_pairs","raw":[18,52,86,120,144,18],"candidate_decimal_pairs":[12,34,56,78,90,12],"candidate_valid":true},
			{"id":"ab46_invalid_nibble","raw":[250,52,86,120,144,18],"candidate_decimal_pairs":null,"candidate_valid":false}
		]`)
	}
}

func TestAb46SegmentedAlternatives(ab46T *ab46Testing.T) {
	for _, ab46Case := range []struct {
		ab46Alternative string
		ab46Width       string
		ab46Raw         string
		ab46Form        string
	}{
		{"link_address_six", `6`, `[1,2,3,4,5,6]`, `"six opaque link-address bytes"`},
		{"conflicting_seven", `7`, `[1,2,3,4,5,6,7]`, `"seven opaque bytes; extra byte uninterpreted"`},
		{"segmented_characters", `null`, `[1,2,3]`, `"raw character segments; boundaries unresolved"`},
	} {
		ab46Row := ab46Get(ab46T, ab46TestKey(84, "Unknown", "identity_content", ab46Case.ab46Alternative))
		ab46EqualJSON(ab46T, ab46Row["width_bytes"], ab46Case.ab46Width)
		ab46EqualJSON(ab46T, ab46Row["synthetic_raw"], ab46Case.ab46Raw)
		ab46EqualJSON(ab46T, ab46Row["raw_form"], ab46Case.ab46Form)
		ab46EqualJSON(ab46T, ab46Row["selection_by_length"], `false`)
		ab46EqualJSON(ab46T, ab46Row["interpretation"], `"Unknown"`)
	}
}

func TestAb46FeedbackAndPartialRetention(ab46T *ab46Testing.T) {
	ab46Row := ab46Get(ab46T, ab46TestKey(84, "Unknown", "identity_state", "none"))
	ab46EqualJSON(ab46T, ab46Row["properties"], `["feedback_success","status_success","content_extent","segment_completeness","complete_identity"]`)
	ab46EqualJSON(ab46T, ab46Row["synthetic_cases"], `[
		{"id":"ab46_empty_success","feedback_success":true,"status_success":true,"content_extent":"empty","segments":[],"segment_completeness":false,"complete_identity":false},
		{"id":"ab46_partial_success","feedback_success":true,"status_success":true,"content_extent":"partial","segments":[{"raw":[65,66],"valid":true},{"raw":[250],"valid":false}],"segment_completeness":false,"complete_identity":false},
		{"id":"ab46_status_failure","feedback_success":true,"status_success":false,"content_extent":"synthetic_full","segments":[{"raw":[65,66],"valid":true},{"raw":[67,68],"valid":true}],"segment_completeness":true,"complete_identity":false},
		{"id":"ab46_segment_failure","feedback_success":true,"status_success":true,"content_extent":"partial","previous_segments":[[65,66],[67,68]],"segments":[{"raw":[69,70],"valid":true},{"raw":[250],"valid":false}],"retained_segments":[[69,70],[67,68]],"segment_completeness":false,"complete_identity":false}
	]`)
}

func TestAb46EveryKeyAxis(ab46T *ab46Testing.T) {
	ab46Control := ab46TestKey(83, "read", "hardware_identifier", "packed_decimal_candidate")
	ab46Type := ab46Reflect.TypeOf(ab46Control)
	for ab46Index := 0; ab46Index < ab46Type.NumField(); ab46Index++ {
		ab46T.Run("ab46_"+ab46Type.Field(ab46Index).Name, func(ab46T *ab46Testing.T) {
			if _, ab46Found, ab46Error := Ab46Lookup(ab46Control); ab46Error != nil || !ab46Found {
				ab46T.Fatal("ab46 missing positive control")
			}
			ab46Changed := ab46Control
			ab46Axis := ab46Reflect.ValueOf(&ab46Changed).Elem().Field(ab46Index)
			if ab46Axis.Kind() == ab46Reflect.String {
				ab46Axis.SetString("ab46_distinct")
			} else {
				ab46Axis.Set(ab46Reflect.ValueOf(ab46Number(0)))
			}
			if _, ab46Found, ab46Error := Ab46Lookup(ab46Changed); ab46Error != nil || ab46Found {
				ab46T.Fatal("ab46 changed axis matched")
			}
			if ab46Axis.Kind() == ab46Reflect.String {
				for _, ab46Token := range []string{"", "*", "none"} {
					ab46Axis.SetString(ab46Token)
					if ab46Reflect.DeepEqual(ab46Changed, ab46Control) {
						continue
					}
					if _, ab46Found, ab46Error := Ab46Lookup(ab46Changed); ab46Error != nil || ab46Found {
						ab46T.Fatal("ab46 string token became a wildcard")
					}
				}
			}
			if ab46Axis.Kind() != ab46Reflect.Ptr {
				return
			}
			ab46Axis.Set(ab46Reflect.ValueOf(ab46Number(65535)))
			if _, ab46Found, ab46Error := Ab46Lookup(ab46Changed); ab46Error != nil || ab46Found {
				ab46T.Fatal("ab46 numeric axis ignored")
			}
			if ab46Type.Field(ab46Index).Name == "Ab46Primary" {
				return
			}
			ab46Axis.Set(ab46Reflect.Zero(ab46Axis.Type()))
			ab46Absent := ab46Changed
			ab46Axis.Set(ab46Reflect.ValueOf(ab46Number(0)))
			ab46Source, ab46Error := ab46JSON.Marshal([]any{map[string]any{"key": ab46Changed, "grade": "Unknown", "ab46_extension": []int{1, 2}}})
			if ab46Error != nil {
				ab46T.Fatal(ab46Error)
			}
			if _, ab46Found, ab46Error := ab46Lookup(ab46Source, ab46Changed); ab46Error != nil || !ab46Found {
				ab46T.Fatal("ab46 zero-valued key control failed")
			}
			if _, ab46Found, ab46Error := ab46Lookup(ab46Source, ab46Absent); ab46Error != nil || ab46Found {
				ab46T.Fatal("ab46 absence matched explicit zero")
			}
		})
	}
}

func TestAb46InvalidAndDuplicateData(ab46T *ab46Testing.T) {
	ab46Original := ab46Data
	defer func() { ab46Data = ab46Original }()
	ab46Key := ab46TestKey(83, "read", "hardware_identifier", "packed_decimal_candidate")
	ab46Row, _, _ := Ab46Lookup(ab46Key)
	ab46Duplicate := append([]byte{'['}, ab46Row...)
	ab46Duplicate = append(ab46Duplicate, ',')
	ab46Duplicate = append(ab46Duplicate, ab46Row...)
	ab46Duplicate = append(ab46Duplicate, ']')
	ab46ChangedRow := ab46Decode(ab46T, ab46Row)
	ab46ChangedRow["grade"] = ab46JSON.RawMessage(`"Unknown"`)
	ab46ChangedRow["width_bytes"] = ab46JSON.RawMessage(`7`)
	ab46ChangedRaw, ab46Error := ab46JSON.Marshal(ab46ChangedRow)
	if ab46Error != nil {
		ab46T.Fatal(ab46Error)
	}
	for _, ab46Bad := range []string{``, `{}`, `null`, `[null]`, `[1]`, `[{}]`, `[{"key":null}]`, `[{"key":{"family":"vaillant"}}]`, `[`, string(ab46Duplicate), "[" + string(ab46Row) + "," + string(ab46ChangedRaw) + "]", "[" + string(ab46Row) + ",null]"} {
		ab46Data = ab46Bad
		if ab46Got, ab46Error := Ab46Entries(); ab46Error == nil || ab46Got != nil {
			ab46T.Fatal("ab46 invalid entries succeeded")
		}
		if ab46Got, ab46Found, ab46Error := Ab46Lookup(ab46Key); ab46Error == nil || ab46Found || ab46Got != nil {
			ab46T.Fatal("ab46 invalid lookup succeeded, including an early matching row")
		}
	}
	ab46Data = `[]`
	if ab46Got, ab46Error := Ab46Entries(); ab46Error != nil || string(ab46Got) != `[]` {
		ab46T.Fatal("ab46 empty array control")
	}
	if ab46Got, ab46Found, ab46Error := Ab46Lookup(ab46Key); ab46Error != nil || ab46Found || ab46Got != nil {
		ab46T.Fatal("ab46 empty array lookup")
	}
}

func TestAb46IndependentStorageAndOpaqueMetadata(ab46T *ab46Testing.T) {
	ab46Key := ab46TestKey(83, "read", "hardware_identifier", "packed_decimal_candidate")
	ab46Before, ab46Error := Ab46Entries()
	if ab46Error != nil {
		ab46T.Fatal(ab46Error)
	}
	ab46Entries, _ := Ab46Entries()
	ab46Entries[0] = '!'
	ab46LookupRaw, ab46Found, ab46Error := Ab46Lookup(ab46Key)
	if ab46Error != nil || !ab46Found {
		ab46T.Fatal("ab46 lookup control")
	}
	ab46LookupRaw[0] = '!'
	ab46After, _ := Ab46Entries()
	if !ab46Bytes.Equal(ab46Before, ab46After) {
		ab46T.Fatal("ab46 returned memory changed catalog")
	}
	ab46Again, _, _ := Ab46Lookup(ab46Key)
	if len(ab46Again) == 0 || ab46Again[0] != '{' {
		ab46T.Fatal("ab46 lookup shared storage")
	}
	for _, ab46Grade := range []string{"Proven", "Hypothesis", "Unknown"} {
		ab46Source, ab46Error := ab46JSON.Marshal([]any{map[string]any{"key": ab46Key, "grade": ab46Grade, "width_bytes": 0, "ab46_extension": map[string]any{"raw": []int{0, 250}, "unresolved": "Unknown"}}})
		if ab46Error != nil {
			ab46T.Fatal(ab46Error)
		}
		ab46Rows, ab46Keys, ab46Error := ab46Validate(ab46Source)
		if ab46Error != nil || len(ab46Rows) != 1 || !ab46Reflect.DeepEqual(ab46Keys[0], ab46Key) {
			ab46T.Fatal("ab46 private seam control")
		}
		ab46Source[0] = '!'
		ab46EqualJSON(ab46T, ab46Decode(ab46T, ab46Rows[0])["grade"], `"`+ab46Grade+`"`)
		ab46Copy := append([]byte{'['}, ab46Rows[0]...)
		ab46Copy = append(ab46Copy, ']')
		ab46Got, ab46Found, ab46Error := ab46Lookup(ab46Copy, ab46Key)
		if ab46Error != nil || !ab46Found || !ab46Bytes.Equal(ab46Got, ab46Rows[0]) {
			ab46T.Fatal("ab46 whole opaque metadata row changed")
		}
		ab46Copy[1] = '!'
		ab46EqualJSON(ab46T, ab46Decode(ab46T, ab46Got)["width_bytes"], `0`)
		ab46EqualJSON(ab46T, ab46Decode(ab46T, ab46Got)["ab46_extension"], `{"raw":[0,250],"unresolved":"Unknown"}`)
		ab46Got[0] = '!'
		if ab46Rows[0][0] != '{' {
			ab46T.Fatal("ab46 source row mutated by lookup result")
		}
	}
}
