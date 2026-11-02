package vaillant

import (
	ab48bytes "bytes"
	ab48json "encoding/json"
	ab48reflect "reflect"
	ab48testing "testing"
)

func ab48Decode(ab48T *ab48testing.T, ab48Raw []byte) map[string]any {
	ab48T.Helper()
	var ab48Result map[string]any
	if ab48Err := ab48json.Unmarshal(ab48Raw, &ab48Result); ab48Err != nil {
		ab48T.Fatal(ab48Err)
	}
	return ab48Result
}

func ab48Marshal(ab48T *ab48testing.T, ab48Value any) []byte {
	ab48T.Helper()
	ab48Raw, ab48Err := ab48json.Marshal(ab48Value)
	if ab48Err != nil {
		ab48T.Fatal(ab48Err)
	}
	return ab48Raw
}

func ab48KeyFor(ab48Field string) Ab48Key {
	ab48Secondary, ab48Operation := uint8(16), uint8(0)
	return Ab48Key{Ab48Family: "vaillant", Ab48Variant: "B510", Ab48Primary: 181, Ab48Secondary: &ab48Secondary, Ab48Operation: &ab48Operation, Ab48Direction: "request", Ab48Field: ab48Field, Ab48Alternative: "none", Ab48InstanceScope: "none"}
}

func TestAb48FiniteMetadata(ab48T *ab48testing.T) {
	ab48Expected := map[string]string{
		"mode":            `{"default":"forced-heating","numeric_mode_code":null,"symbolic_contexts":["forced-heating","forced-DHW","forced-off"],"mode_encodings":"Unknown"}`,
		"modulation":      `{"declared_default":67,"width_bits":null,"domain":null}`,
		"options":         `{"bit_defaults":[1,0,0,0,0,0,0,0],"bit2_context":{"field":"CH_setpoint","clear_role":"heating","set_role":"cooling","raw_value_changes":false},"width_bits":null}`,
		"loading_pump":    `{"width_bits":8,"signedness":"unsigned","conversion":"raw / 2","unit":"percent","formal_domain":[0,127.5],"declared_default":255.0,"default_conflict":true,"raw_example":{"bytes":"FF","value":127.5},"forced_mode_dependency":true,"contexts":[{"mode":"forced-heating","numeric_mode_code":null,"meaning":"Unknown","activation":"Unknown"},{"mode":"forced-DHW","numeric_mode_code":null,"meaning":"Unknown","activation":"Unknown"},{"mode":"forced-off","numeric_mode_code":null,"meaning":"Unknown","activation":"Unknown"}]}`,
		"CH_setpoint":     `{"width_bits":null,"zero_degrees_role":["stop","pump-stop"],"transition":{"from":0,"to":"positive","candidate_minimum_seconds":300,"execution_guarantee":false,"enforcement":"Unknown"},"other_transitions":"Unknown"}`,
		"pump_activation": `{"width_bits":null,"forced_mode_dependency":true,"contexts":[{"mode":"forced-heating","numeric_mode_code":null,"activation":"Unknown"},{"mode":"forced-DHW","numeric_mode_code":null,"activation":"Unknown"},{"mode":"forced-off","numeric_mode_code":null,"activation":"Unknown"}]}`,
	}
	ab48Entries, ab48Err := Ab48Entries()
	if ab48Err != nil {
		ab48T.Fatal(ab48Err)
	}
	var ab48Rows []ab48json.RawMessage
	if ab48Err = ab48json.Unmarshal(ab48Entries, &ab48Rows); ab48Err != nil || len(ab48Rows) != 6 {
		ab48T.Fatalf("ab48 finite array: %d %v", len(ab48Rows), ab48Err)
	}
	for ab48Field, ab48Details := range ab48Expected {
		ab48T.Run("ab48_"+ab48Field, func(ab48T *ab48testing.T) {
			ab48Key := ab48KeyFor(ab48Field)
			ab48Raw, ab48Found, ab48Err := Ab48Lookup(ab48Key)
			if ab48Err != nil || !ab48Found {
				ab48T.Fatalf("ab48 lookup: %v %v", ab48Found, ab48Err)
			}
			ab48ExpectedRow := map[string]any{"id": "ab48_" + ab48Field, "key": ab48Key, "grade": "Hypothesis", "metadata_presence": "Proven", "physical_qualification": "Unknown", "presence": "Unknown", "freshness": "Unknown", "execution": "Unknown", "native_contract": "https://github.com/Project-Helianthus/helianthus-docs-ebus/blob/main/protocols/vaillant/ebus-vaillant-B510.md", "position": nil, "raw_form": "preserve original bytes", "unknown": []string{"other positions", "peer acceptance", "effects"}, "details": ab48Decode(ab48T, []byte(ab48Details))}
			if !ab48reflect.DeepEqual(ab48Decode(ab48T, ab48Raw), ab48Decode(ab48T, ab48Marshal(ab48T, ab48ExpectedRow))) {
				ab48T.Fatalf("ab48 whole row mismatch: %s", ab48Raw)
			}
			ab48Matches := 0
			for _, ab48Row := range ab48Rows {
				if ab48bytes.Equal(ab48Row, ab48Raw) {
					ab48Matches++
				}
			}
			if ab48Matches != 1 {
				ab48T.Fatal("ab48 lookup must retain whole exact raw row")
			}
		})
	}
}

func TestAb48EveryKeyAxis(ab48T *ab48testing.T) {
	ab48Zero, ab48One := uint8(0), uint8(1)
	ab48RegisterZero, ab48RegisterOne := uint16(0), uint16(1)
	ab48Changes := map[string]func(*Ab48Key){
		"ab48_family":           func(ab48Key *Ab48Key) { ab48Key.Ab48Family = "ab48_other" },
		"ab48_variant":          func(ab48Key *Ab48Key) { ab48Key.Ab48Variant = "ab48_other" },
		"ab48_primary":          func(ab48Key *Ab48Key) { ab48Key.Ab48Primary = 0 },
		"ab48_secondary_absent": func(ab48Key *Ab48Key) { ab48Key.Ab48Secondary = nil },
		"ab48_secondary_zero":   func(ab48Key *Ab48Key) { ab48Key.Ab48Secondary = &ab48Zero },
		"ab48_operation_absent": func(ab48Key *Ab48Key) { ab48Key.Ab48Operation = nil },
		"ab48_operation_one":    func(ab48Key *Ab48Key) { ab48Key.Ab48Operation = &ab48One },
		"ab48_direction":        func(ab48Key *Ab48Key) { ab48Key.Ab48Direction = "response" },
		"ab48_field":            func(ab48Key *Ab48Key) { ab48Key.Ab48Field = "ab48_other" },
		"ab48_alternative":      func(ab48Key *Ab48Key) { ab48Key.Ab48Alternative = "ab48_other" },
		"ab48_group_zero":       func(ab48Key *Ab48Key) { ab48Key.Ab48Group = &ab48Zero },
		"ab48_group_one":        func(ab48Key *Ab48Key) { ab48Key.Ab48Group = &ab48One },
		"ab48_instance_scope":   func(ab48Key *Ab48Key) { ab48Key.Ab48InstanceScope = "ab48_other" },
		"ab48_register_zero":    func(ab48Key *Ab48Key) { ab48Key.Ab48Register = &ab48RegisterZero },
		"ab48_register_one":     func(ab48Key *Ab48Key) { ab48Key.Ab48Register = &ab48RegisterOne },
	}
	for ab48Name, ab48Change := range ab48Changes {
		ab48T.Run(ab48Name, func(ab48T *ab48testing.T) {
			ab48Base := ab48KeyFor("loading_pump")
			ab48Original, ab48Found, ab48Err := Ab48Lookup(ab48Base)
			if ab48Err != nil || !ab48Found {
				ab48T.Fatal("ab48 invalid control")
			}
			ab48Changed := ab48KeyFor("loading_pump")
			ab48Change(&ab48Changed)
			if ab48Raw, ab48Found, ab48Err := Ab48Lookup(ab48Changed); ab48Err != nil || ab48Found || ab48Raw != nil {
				ab48T.Fatal("ab48 namespace axis ignored")
			}
			ab48Second := ab48Decode(ab48T, ab48Original)
			ab48Second["key"] = ab48Changed
			ab48Source := ab48Marshal(ab48T, []any{ab48json.RawMessage(ab48Original), ab48Second})
			ab48Want := ab48Decode(ab48T, ab48Marshal(ab48T, ab48Second))
			ab48Raw, ab48Found, ab48Err := ab48LookupSource(ab48Source, ab48Changed)
			if ab48Err != nil || !ab48Found || !ab48reflect.DeepEqual(ab48Decode(ab48T, ab48Raw), ab48Want) {
				ab48T.Fatalf("ab48 distinct valid control: %v", ab48Err)
			}
		})
	}
}

func TestAb48WidthAndGradeRetention(ab48T *ab48testing.T) {
	ab48Raw, _, ab48Err := Ab48Lookup(ab48KeyFor("loading_pump"))
	if ab48Err != nil {
		ab48T.Fatal(ab48Err)
	}
	for ab48Name, ab48Value := range map[string]any{"ab48_width": float64(16), "ab48_grade": "Unknown"} {
		ab48T.Run(ab48Name, func(ab48T *ab48testing.T) {
			ab48Row := ab48Decode(ab48T, ab48Raw)
			if ab48Name == "ab48_width" {
				ab48Row["details"].(map[string]any)["width_bits"] = ab48Value
			} else {
				ab48Row["grade"] = ab48Value
			}
			ab48Got, ab48Found, ab48Err := ab48LookupSource(ab48Marshal(ab48T, []any{ab48Row}), ab48KeyFor("loading_pump"))
			if ab48Err != nil || !ab48Found || !ab48reflect.DeepEqual(ab48Decode(ab48T, ab48Got), ab48Row) {
				ab48T.Fatal("ab48 row metadata normalized or lost")
			}
			if _, _, ab48Err = ab48Validate(ab48Marshal(ab48T, []any{ab48json.RawMessage(ab48Raw), ab48Row})); ab48Err == nil {
				ab48T.Fatal("ab48 metadata must not disambiguate duplicate identity")
			}
		})
	}
}

func TestAb48OptionalAbsentVersusZero(ab48T *ab48testing.T) {
	ab48Zero, ab48RegisterZero := uint8(0), uint16(0)
	for _, ab48Axis := range []string{"secondary", "operation", "group", "register"} {
		ab48T.Run("ab48_"+ab48Axis, func(ab48T *ab48testing.T) {
			ab48Absent, ab48Explicit := ab48KeyFor("mode"), ab48KeyFor("mode")
			switch ab48Axis {
			case "secondary":
				ab48Absent.Ab48Secondary, ab48Explicit.Ab48Secondary = nil, &ab48Zero
			case "operation":
				ab48Absent.Ab48Operation, ab48Explicit.Ab48Operation = nil, &ab48Zero
			case "group":
				ab48Absent.Ab48Group, ab48Explicit.Ab48Group = nil, &ab48Zero
			case "register":
				ab48Absent.Ab48Register, ab48Explicit.Ab48Register = nil, &ab48RegisterZero
			}
			ab48Source := ab48Marshal(ab48T, []any{map[string]any{"key": ab48Absent, "id": "ab48_absent"}, map[string]any{"key": ab48Explicit, "id": "ab48_zero"}})
			for ab48Index, ab48Key := range []Ab48Key{ab48Absent, ab48Explicit} {
				ab48Raw, ab48Found, ab48Err := ab48LookupSource(ab48Source, ab48Key)
				if ab48Err != nil || !ab48Found || ab48Decode(ab48T, ab48Raw)["id"] != []string{"ab48_absent", "ab48_zero"}[ab48Index] {
					ab48T.Fatalf("ab48 absent and zero conflated: %v", ab48Err)
				}
			}
		})
	}
}

func TestAb48InvalidAndDuplicate(ab48T *ab48testing.T) {
	ab48Row, _, _ := Ab48Lookup(ab48KeyFor("mode"))
	ab48Cases := map[string][]byte{"ab48_syntax": []byte("["), "ab48_null": []byte("null"), "ab48_object": []byte("{}"), "ab48_null_row": []byte("[null]"), "ab48_missing_key": []byte("[{}]"), "ab48_trailing": []byte("[] []"), "ab48_duplicate": ab48Marshal(ab48T, []any{ab48json.RawMessage(ab48Row), ab48json.RawMessage(ab48Row)}), "ab48_invalid_later": ab48Marshal(ab48T, []any{ab48json.RawMessage(ab48Row), nil})}
	for _, ab48Axis := range []string{"family", "variant", "primary", "direction", "field", "alternative", "instance_scope"} {
		ab48Changed := ab48Decode(ab48T, ab48Row)
		delete(ab48Changed["key"].(map[string]any), ab48Axis)
		ab48Cases["ab48_missing_"+ab48Axis] = ab48Marshal(ab48T, []any{ab48Changed})
	}
	ab48Saved := ab48Data
	defer func() { ab48Data = ab48Saved }()
	for ab48Name, ab48Source := range ab48Cases {
		ab48T.Run(ab48Name, func(ab48T *ab48testing.T) {
			ab48Data = string(ab48Source)
			if ab48Raw, ab48Err := Ab48Entries(); ab48Err == nil || ab48Raw != nil {
				ab48T.Fatal("ab48 invalid entries accepted")
			}
			if ab48Raw, ab48Found, ab48Err := Ab48Lookup(ab48KeyFor("mode")); ab48Err == nil || ab48Found || ab48Raw != nil {
				ab48T.Fatal("ab48 invalid lookup accepted")
			}
		})
	}
}

func TestAb48StorageIsolation(ab48T *ab48testing.T) {
	ab48Before, ab48Err := Ab48Entries()
	if ab48Err != nil {
		ab48T.Fatal(ab48Err)
	}
	ab48Source := ab48bytes.Clone(ab48Before)
	ab48Rows, _, ab48Err := ab48Validate(ab48Source)
	if ab48Err != nil {
		ab48T.Fatal(ab48Err)
	}
	ab48RowBefore := ab48bytes.Clone(ab48Rows[0])
	ab48LookupBefore, _, _ := Ab48Lookup(ab48KeyFor("mode"))
	ab48LookupCopy, _, _ := ab48LookupSource(ab48Source, ab48KeyFor("mode"))
	for ab48Index := range ab48Source {
		ab48Source[ab48Index] = 'x'
	}
	if !ab48bytes.Equal(ab48Rows[0], ab48RowBefore) || !ab48bytes.Equal(ab48LookupCopy, ab48LookupBefore) {
		ab48T.Fatal("ab48 source mutation aliases results")
	}
	ab48Returned, _ := Ab48Entries()
	for ab48Index := range ab48Returned {
		ab48Returned[ab48Index] = 'x'
	}
	for ab48Index := range ab48LookupCopy {
		ab48LookupCopy[ab48Index] = 'x'
	}
	ab48PublicCopy, _, _ := Ab48Lookup(ab48KeyFor("mode"))
	for ab48Index := range ab48PublicCopy {
		ab48PublicCopy[ab48Index] = 'x'
	}
	ab48After, _ := Ab48Entries()
	ab48LookupAfter, _, _ := Ab48Lookup(ab48KeyFor("mode"))
	if !ab48bytes.Equal(ab48Before, ab48After) || !ab48bytes.Equal(ab48LookupBefore, ab48LookupAfter) {
		ab48T.Fatal("ab48 returned memory aliases catalog")
	}
}

func TestAb48SyntheticScenarios(ab48T *ab48testing.T) {
	ab48PumpRaw, _, _ := Ab48Lookup(ab48KeyFor("loading_pump"))
	ab48Pump := ab48Decode(ab48T, ab48PumpRaw)["details"].(map[string]any)
	if ab48Pump["conversion"] != "raw / 2" || float64(byte(0xff))/2 != 127.5 || ab48Pump["raw_example"].(map[string]any)["value"] != 127.5 || ab48Pump["declared_default"] != 255.0 {
		ab48T.Fatal("ab48 raw FF and conflicting literal lost")
	}
	ab48OptionsRaw, _, _ := Ab48Lookup(ab48KeyFor("options"))
	ab48Context := ab48Decode(ab48T, ab48OptionsRaw)["details"].(map[string]any)["bit2_context"].(map[string]any)
	ab48Raw := byte(42)
	for ab48Index, ab48Option := range []byte{0, 4} {
		ab48Role := ab48Context["clear_role"]
		if ab48Option&4 != 0 {
			ab48Role = ab48Context["set_role"]
		}
		if ab48Role != []string{"heating", "cooling"}[ab48Index] || ab48Raw != 42 || ab48Context["raw_value_changes"] != false {
			ab48T.Fatal("ab48 option context changed raw")
		}
	}
	ab48SetpointRaw, _, _ := Ab48Lookup(ab48KeyFor("CH_setpoint"))
	ab48Details := ab48Decode(ab48T, ab48SetpointRaw)["details"].(map[string]any)
	ab48Transition := ab48Details["transition"].(map[string]any)
	for ab48Index, ab48Pair := range [][2]int{{0, 1}, {1, 2}} {
		ab48Candidate := ab48Pair[0] == int(ab48Transition["from"].(float64)) && ab48Pair[1] > 0
		if ab48Candidate != (ab48Index == 0) {
			ab48T.Fatal("ab48 distinct transition collapsed")
		}
		if ab48Candidate && ab48Transition["candidate_minimum_seconds"] != float64(300) {
			ab48T.Fatal("ab48 candidate interval")
		}
		if !ab48Candidate && ab48Details["other_transitions"] != "Unknown" {
			ab48T.Fatal("ab48 invented interval")
		}
	}
	if ab48Transition["execution_guarantee"] != false || ab48Transition["enforcement"] != "Unknown" {
		ab48T.Fatal("ab48 candidate promoted to guarantee")
	}
}
