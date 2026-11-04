package vaillant

import (
	ab34bytes "bytes"
	ab34json "encoding/json"
	ab34reflect "reflect"
	ab34strings "strings"
	ab34testing "testing"
)

func ab34BaseKey() Ab34Key {
	ab34Secondary := uint8(1)
	return Ab34Key{Ab34Family: "legacy-interface", Ab34Variant: "legacy-computer-interface", Ab34Primary: 181, Ab34Secondary: &ab34Secondary, Ab34Direction: "initiator-to-initiator", Ab34Field: "sequencing", Ab34Alternative: "none", Ab34InstanceScope: "none"}
}

func ab34Decode(ab34T *ab34testing.T, ab34Data []byte) map[string]any {
	ab34T.Helper()
	var ab34Object map[string]any
	if ab34Error := ab34json.Unmarshal(ab34Data, &ab34Object); ab34Error != nil {
		ab34T.Fatal(ab34Error)
	}
	return ab34Object
}

func ab34Marshal(ab34T *ab34testing.T, ab34Value any) []byte {
	ab34T.Helper()
	ab34Data, ab34Error := ab34json.Marshal(ab34Value)
	if ab34Error != nil {
		ab34T.Fatal(ab34Error)
	}
	return ab34Data
}

func ab34Get(ab34T *ab34testing.T, ab34Key Ab34Key) []byte {
	ab34T.Helper()
	ab34Row, ab34Found, ab34Error := Ab34Lookup(ab34Key)
	if ab34Error != nil || !ab34Found || len(ab34Row) == 0 {
		ab34T.Fatalf("ab34 lookup: found=%v error=%v", ab34Found, ab34Error)
	}
	return ab34Row
}

func TestAb34Contract(ab34T *ab34testing.T) {
	ab34Data, ab34Error := Ab34Entries()
	if ab34Error != nil {
		ab34T.Fatal(ab34Error)
	}
	var ab34Rows []ab34json.RawMessage
	if ab34Error = ab34json.Unmarshal(ab34Data, &ab34Rows); ab34Error != nil || len(ab34Rows) != 3 {
		ab34T.Fatalf("ab34 finite array: %d %v", len(ab34Rows), ab34Error)
	}
	for ab34Index, ab34Role := range []string{"request", "positive-status", "return-data"} {
		ab34Key := ab34BaseKey()
		ab34Service := "B501"
		if ab34Index == 1 {
			ab34Key.Ab34Field, ab34Key.Ab34Alternative = "status", "positive"
		}
		if ab34Index == 2 {
			ab34Secondary := uint8(2)
			ab34Key.Ab34Secondary = &ab34Secondary
			ab34Key.Ab34Field = "data"
			ab34Service = "B502"
		}
		ab34Raw := ab34Get(ab34T, ab34Key)
		if !ab34bytes.Equal(ab34Raw, ab34Rows[ab34Index]) {
			ab34T.Fatal("ab34 lookup must return the whole original row")
		}
		ab34Row := ab34Decode(ab34T, ab34Raw)
		ab34Want := map[string]any{
			"ab34_id":                   "ab34_legacy_" + ab34Role,
			"ab34_identity":             ab34Decode(ab34T, ab34Marshal(ab34T, ab34Key)),
			"ab34_representation_grade": "Proven", "ab34_candidate_grade": "Hypothesis",
			"ab34_service": ab34Service, "ab34_service_role": ab34Role,
			"ab34_raw_form": map[string]any{"ab34_service": ab34Service, "ab34_payload": nil, "ab34_policy": "preserve complete native bytes and unaffected last-known-good fields"},
			"ab34_unknown":  map[string]any{"ab34_synchronization_bytes": "Unknown", "ab34_state_order_beyond_wait": "Unknown", "ab34_timeout_duration": "Unknown", "ab34_segmentation": "Unknown", "ab34_retries": "Unknown", "ab34_return_codec": "Unknown", "ab34_transport_adaptation": "Unknown", "ab34_native_support": "Unknown", "ab34_physical_effects": "Unknown", "ab34_persistence": "Unknown", "ab34_delivery": "Unknown"},
		}
		if ab34Index == 0 {
			ab34Want["ab34_sequence"] = ab34Decode(ab34T, []byte(`{"ab34_synchronization_required":true,"ab34_synchronization_token":null,"ab34_timeout_cooldown_min_ms":3500,"ab34_cooldown_event":"timeout","ab34_success_cooldown_required":false,"ab34_request_service":"B501","ab34_positive_status_service":"B501","ab34_return_data_service":"B502","ab34_status_is_return_data":false}`))
		} else {
			ab34Want["ab34_context_only"] = true
		}
		if !ab34reflect.DeepEqual(ab34Row, ab34Want) {
			ab34T.Errorf("ab34 %s differs: got %s want %s", ab34Role, ab34Raw, ab34Marshal(ab34T, ab34Want))
		}
	}
}

func TestAb34TimeoutBoundary(ab34T *ab34testing.T) {
	ab34Sequence := ab34Decode(ab34T, ab34Get(ab34T, ab34BaseKey()))["ab34_sequence"].(map[string]any)
	for _, ab34Case := range []struct {
		ab34Name    string
		ab34Event   string
		ab34Elapsed int
		ab34Want    bool
	}{
		{"ab34_timeout_3499", "timeout", 3499, false},
		{"ab34_timeout_3500", "timeout", 3500, true},
		{"ab34_timeout_3501", "timeout", 3501, true},
		{"ab34_success_0", "successful-return", 0, true},
		{"ab34_success_3499", "successful-return", 3499, true},
	} {
		ab34T.Run(ab34Case.ab34Name, func(ab34T *ab34testing.T) {
			ab34Allowed := ab34Case.ab34Event != ab34Sequence["ab34_cooldown_event"] || ab34Case.ab34Elapsed >= int(ab34Sequence["ab34_timeout_cooldown_min_ms"].(float64))
			if ab34Allowed != ab34Case.ab34Want {
				ab34T.Fatal("ab34 timeout-only cooldown boundary differs")
			}
		})
	}
	for _, ab34Name := range []string{"ab34_timeout_duration_ms", "ab34_timeout_duration", "ab34_synchronization_bytes"} {
		if _, ab34Present := ab34Sequence[ab34Name]; ab34Present {
			ab34T.Fatalf("ab34 invented %s", ab34Name)
		}
	}
}

func TestAb34EveryIdentityAxis(ab34T *ab34testing.T) {
	ab34Base := ab34BaseKey()
	ab34Original := ab34Get(ab34T, ab34Base)
	for ab34Index := 0; ab34Index < ab34reflect.TypeOf(ab34Base).NumField(); ab34Index++ {
		ab34Name := ab34reflect.TypeOf(ab34Base).Field(ab34Index).Name
		ab34T.Run("ab34_"+ab34Name, func(ab34T *ab34testing.T) {
			ab34Changed := ab34BaseKey()
			ab34Field := ab34reflect.ValueOf(&ab34Changed).Elem().Field(ab34Index)
			switch ab34Field.Kind() {
			case ab34reflect.String:
				ab34Field.SetString("ab34_synthetic_other")
				if ab34Name == "Ab34Family" {
					ab34Field.SetString("modern-controller")
				}
			case ab34reflect.Uint8:
				ab34Field.SetUint(180)
			case ab34reflect.Pointer:
				ab34Pointer := ab34reflect.New(ab34Field.Type().Elem())
				ab34Pointer.Elem().SetUint(0)
				ab34Field.Set(ab34Pointer)
			}
			ab34Row, ab34Found, ab34Error := Ab34Lookup(ab34Changed)
			if ab34Error != nil || ab34Found || ab34Row != nil {
				ab34T.Fatalf("ab34 changed axis matched: %v %v", ab34Found, ab34Error)
			}
			ab34Control := ab34Decode(ab34T, ab34Original)
			ab34Control["ab34_identity"] = ab34Changed
			ab34Data := ab34Marshal(ab34T, []any{ab34Decode(ab34T, ab34Original), ab34Control})
			if _, ab34Error = ab34Validate(ab34Data); ab34Error != nil {
				ab34T.Fatalf("ab34 distinct-axis valid control rejected: %v", ab34Error)
			}
			ab34Row, ab34Found, ab34Error = ab34Lookup(ab34Data, ab34Changed)
			if ab34Error != nil || !ab34Found || !ab34reflect.DeepEqual(ab34Decode(ab34T, ab34Row), ab34Decode(ab34T, ab34Marshal(ab34T, ab34Control))) {
				ab34T.Fatalf("ab34 distinct-axis lookup control failed: %v", ab34Error)
			}
		})
	}
}

func TestAb34OptionalAbsentVersusZero(ab34T *ab34testing.T) {
	for _, ab34Index := range []int{3, 4, 8, 10} {
		ab34Name := ab34reflect.TypeOf(Ab34Key{}).Field(ab34Index).Name
		ab34T.Run("ab34_"+ab34Name, func(ab34T *ab34testing.T) {
			ab34Absent, ab34Zero := ab34BaseKey(), ab34BaseKey()
			ab34AbsentField := ab34reflect.ValueOf(&ab34Absent).Elem().Field(ab34Index)
			ab34AbsentField.Set(ab34reflect.Zero(ab34AbsentField.Type()))
			ab34ZeroField := ab34reflect.ValueOf(&ab34Zero).Elem().Field(ab34Index)
			ab34ZeroField.Set(ab34reflect.New(ab34ZeroField.Type().Elem()))
			ab34Row := ab34Decode(ab34T, ab34Get(ab34T, ab34BaseKey()))
			ab34Other := ab34Decode(ab34T, ab34Get(ab34T, ab34BaseKey()))
			ab34Row["ab34_identity"], ab34Other["ab34_identity"] = ab34Absent, ab34Zero
			ab34Other["ab34_id"] = "ab34_zero_control"
			ab34Data := ab34Marshal(ab34T, []any{ab34Row, ab34Other})
			for _, ab34Case := range []struct {
				ab34Key Ab34Key
				ab34ID  string
			}{{ab34Absent, "ab34_legacy_request"}, {ab34Zero, "ab34_zero_control"}} {
				ab34Got, ab34Found, ab34Error := ab34Lookup(ab34Data, ab34Case.ab34Key)
				if ab34Error != nil || !ab34Found || ab34Decode(ab34T, ab34Got)["ab34_id"] != ab34Case.ab34ID {
					ab34T.Fatalf("ab34 absent/zero collision: %v", ab34Error)
				}
			}
		})
	}
}

func TestAb34ValidationErrors(ab34T *ab34testing.T) {
	ab34Row := ab34Get(ab34T, ab34BaseKey())
	ab34Duplicate := ab34Marshal(ab34T, []ab34json.RawMessage{ab34Row, ab34Row})
	for _, ab34Data := range [][]byte{[]byte(`[`), []byte(`{}`), []byte(`null`), []byte(`[null]`), []byte(`[{}]`), ab34Duplicate, []byte(`[{}] trailing`)} {
		if _, ab34Error := ab34Validate(ab34Data); ab34Error == nil {
			ab34T.Errorf("ab34 accepted invalid data %s", ab34Data)
		}
		if ab34Got, ab34Found, ab34Error := ab34Lookup(ab34Data, ab34BaseKey()); ab34Error == nil || ab34Found || ab34Got != nil {
			ab34T.Errorf("ab34 lookup failed open for %s", ab34Data)
		}
	}
	ab34Object := ab34Decode(ab34T, ab34Row)
	for ab34Axis := range ab34Object["ab34_identity"].(map[string]any) {
		ab34T.Run("ab34_missing_"+ab34Axis, func(ab34T *ab34testing.T) {
			ab34Control := ab34Decode(ab34T, ab34Row)
			delete(ab34Control["ab34_identity"].(map[string]any), ab34Axis)
			if _, ab34Error := ab34Validate(ab34Marshal(ab34T, []any{ab34Control})); ab34Error == nil {
				ab34T.Fatal("ab34 incomplete identity admitted")
			}
		})
	}
	for _, ab34Bad := range []string{"null", "-1", "256", `"B501"`} {
		ab34Data := ab34strings.Replace(string(ab34Row), `"ab34_primary": 181`, `"ab34_primary": `+ab34Bad, 1)
		if _, ab34Error := ab34Validate([]byte("[" + ab34Data + "]")); ab34Error == nil {
			ab34T.Errorf("ab34 invalid primary admitted: %s", ab34Bad)
		}
	}
}

func TestAb34IndependentBuffers(ab34T *ab34testing.T) {
	ab34Entries, ab34Error := Ab34Entries()
	if ab34Error != nil {
		ab34T.Fatal(ab34Error)
	}
	ab34Snapshot := ab34bytes.Clone(ab34Entries)
	ab34Row := ab34Get(ab34T, ab34BaseKey())
	ab34RowSnapshot := ab34bytes.Clone(ab34Row)
	ab34SecondRow := ab34Get(ab34T, ab34BaseKey())
	for ab34Index := range ab34Row {
		ab34Row[ab34Index] ^= 255
	}
	if !ab34bytes.Equal(ab34SecondRow, ab34RowSnapshot) || !ab34bytes.Equal(ab34Entries, ab34Snapshot) || !ab34bytes.Equal(ab34Get(ab34T, ab34BaseKey()), ab34RowSnapshot) {
		ab34T.Fatal("ab34 row memory shared")
	}
	for ab34Index := range ab34Entries {
		ab34Entries[ab34Index] ^= 255
	}
	ab34Later, ab34Error := Ab34Entries()
	if ab34Error != nil || !ab34bytes.Equal(ab34Later, ab34Snapshot) || !ab34bytes.Equal(ab34Get(ab34T, ab34BaseKey()), ab34RowSnapshot) {
		ab34T.Fatal("ab34 entries memory shared")
	}
	ab34Input := ab34bytes.Clone(ab34Snapshot)
	ab34Owned, ab34Found, ab34Error := ab34Lookup(ab34Input, ab34BaseKey())
	if ab34Error != nil || !ab34Found {
		ab34T.Fatal("ab34 seam control failed")
	}
	for ab34Index := range ab34Input {
		ab34Input[ab34Index] ^= 255
	}
	if !ab34bytes.Equal(ab34Owned, ab34RowSnapshot) {
		ab34T.Fatal("ab34 lookup retains caller input")
	}
}
