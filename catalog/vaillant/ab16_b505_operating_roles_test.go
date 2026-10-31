package vaillant

import (
	ab16bytes "bytes"
	ab16json "encoding/json"
	ab16reflect "reflect"
	ab16testing "testing"
)

type ab16CatalogRow struct {
	Ab16ID        string            `json:"id"`
	Ab16Identity  Ab16Key           `json:"identity"`
	Ab16Role      string            `json:"ab16Role"`
	Ab16Relations map[string]any    `json:"ab16Relations"`
	Ab16Grade     string            `json:"grade"`
	Ab16Unknown   map[string]string `json:"ab16Unknown"`
	Ab16Raw       string            `json:"raw_fallback"`
}

func ab16Rows(ab16T *ab16testing.T) []ab16CatalogRow {
	ab16T.Helper()
	ab16Data, ab16Err := Ab16Entries()
	if ab16Err != nil {
		ab16T.Fatal(ab16Err)
	}
	var ab16Rows []ab16CatalogRow
	if ab16Err = ab16json.Unmarshal(ab16Data, &ab16Rows); ab16Err != nil {
		ab16T.Fatal(ab16Err)
	}
	for _, ab16Row := range ab16Rows {
		ab16Raw, ab16Found, ab16Err := Ab16Lookup(ab16Row.Ab16Identity)
		if ab16Err != nil || !ab16Found {
			ab16T.Fatalf("lookup: %v %v", ab16Found, ab16Err)
		}
		var ab16Actual ab16CatalogRow
		if ab16Err = ab16json.Unmarshal(ab16Raw, &ab16Actual); ab16Err != nil || !ab16reflect.DeepEqual(ab16Row, ab16Actual) {
			ab16T.Fatalf("whole row: %v", ab16Err)
		}
	}
	return ab16Rows
}

func ab16Find(ab16T *ab16testing.T, ab16ID string) ab16CatalogRow {
	ab16T.Helper()
	for _, ab16Row := range ab16Rows(ab16T) {
		if ab16Row.Ab16ID == "ab16_"+ab16ID {
			return ab16Row
		}
	}
	ab16T.Fatalf("missing %s", ab16ID)
	return ab16CatalogRow{}
}

func TestAb16Associations(ab16T *ab16testing.T) {
	ab16Expected := []struct{ ab16ID, ab16Role, ab16Relations string }{
		{"5A_unchanged_setpoint", "room_day_night", `{"ab16Replacement":true,"ab16ReplacementRaw":"Unknown","ab16ReplacementAction":"retain_prior_token","ab16NumericBranch":"numeric_candidate_without_scale","ab16BranchOrder":"replacement_before_numeric"}`},
		{"3D_unchanged_setpoint", "room_day_night", `{"ab16Replacement":true,"ab16ReplacementRaw":"Unknown","ab16ReplacementAction":"retain_prior_token","ab16NumericBranch":"numeric_candidate_without_scale","ab16BranchOrder":"replacement_before_numeric"}`},
		{"47_heating_night", "inverter_indoor_circuit", `{"ab16Mode":"heating","ab16Period":"night","ab16Quantity":"room_setpoint"}`},
		{"47_cooling_night", "inverter_indoor_circuit", `{"ab16Mode":"cooling","ab16Period":"night","ab16Quantity":"room_setpoint"}`},
		{"48_heating_holiday", "inverter_indoor_circuit", `{"ab16Mode":"heating","ab16Period":"holiday","ab16Quantity":"room_setpoint"}`},
		{"48_cooling_holiday", "inverter_indoor_circuit", `{"ab16Mode":"cooling","ab16Period":"holiday","ab16Quantity":"room_setpoint"}`},
		{"41_frost_protection", "solar", `{"ab16ReplacementAction":"off","ab16ReplacementRaw":"Unknown"}`},
		{"40_minimum_temperature", "second_solar_circuit_or_solid_fuel_boiler", `{}`},
		{"10_maximum_preheat_time", "heating", `{"ab16Action":"advance_before_first_period"}`},
		{"3F_maximum_preheat_time", "heating", `{"ab16Action":"advance_before_first_period"}`},
		{"15_early_cutoff", "heating", `{"ab16Action":"stop_demand_before_reduction_period"}`},
		{"2B_early_cutoff", "heating", `{"ab16Action":"stop_demand_before_reduction_period"}`},
		{"3C_uncorrected_room_temperature", "room", `{"ab16SelfHeating":true,"ab16UserCorrection":false}`},
		{"3C_corrected_room_temperature", "room", `{"ab16SelfHeating":true,"ab16UserCorrection":true}`},
		{"1B_switch_on_cutoff_difference", "solar_tank_one", `{"ab16Tank":"one"}`},
		{"1E_switch_on_cutoff_difference", "solar_tank_two", `{"ab16Tank":"two"}`},
		{"38_switch_on_cutoff_difference", "solar_tank_three", `{"ab16Tank":"three"}`},
		{"1A_maximum_loading_temperature", "solar_tank_one", `{"ab16Tank":"one"}`},
		{"1D_maximum_loading_temperature", "solar_tank_two", `{"ab16Tank":"two"}`},
		{"37_maximum_loading_temperature", "solar_tank_three", `{"ab16Tank":"three"}`},
		{"31_ordinary_dhw", "dhw", `{"ab16EnabledValue":1,"ab16EnabledEligibility":"ordinary_dhw","ab16BlockedEligibility":"frost_only"}`},
		{"19_charging_enable", "dhw_tank_charging", `{"ab16EnabledValue":1}`},
		{"27_charging_enable", "dhw_tank_charging", `{"ab16EnabledValue":1}`},
		{"18_charging_disable", "circuit_charging", `{"ab16DisabledValue":1,"ab16Action":"suppress_without_parallel","ab16ExceptionRole":"mixer"}`},
		{"27_charging_disable", "circuit_charging", `{"ab16DisabledValue":1,"ab16Action":"suppress_without_parallel","ab16ExceptionRole":"mixer"}`},
		{"14_parallel_charging", "mixer", `{"ab16ExceptionRole":"mixer"}`},
		{"27_tank_flow_temperature", "dhw_flow", `{}`},
		{"2F_external_storage_sensor_temperature", "external_storage_sensor", `{}`},
		{"1C_outside_correction", "outside", `{}`},
		{"26_hydraulic_schema", "solar_controller", `{}`},
		{"23_duty_control", "collector_pumps", `{}`},
		{"21_flow", "collector_water", `{}`},
		{"22_reset", "solar", `{"ab16RequestValue":1,"ab16Intent":"clear_solar_pump_hours","ab16Accepted":"Unknown","ab16Persistent":"Unknown"}`},
		{"20_reset", "solar", `{"ab16RequestValue":1,"ab16Intent":"clear_yield","ab16Accepted":"Unknown","ab16Persistent":"Unknown"}`},
		{"17_flow_excess", "flow_demand", `{"ab16Action":"add_to_demanded_flow","ab16Encoding":"Unknown"}`},
		{"0F_maximum_flow", "flow_demand", `{"ab16Action":"maximum"}`},
		{"0E_minimum_flow", "flow_demand", `{"ab16Zero":"off","ab16Nonzero":"minimum_candidate"}`},
		{"12_overrun", "charging", `{"ab16Action":"retain_pump_after_charging","ab16Boiler":"inhibited"}`},
		{"11_recharge_delay", "solar_charging", `{"ab16Action":"limit_recharge_delay"}`},
		{"0D_pump_idle", "pump", `{"ab16Action":"energy_saving_interval_after_flow_exceeds_desired"}`},
		{"0C_outside_cutoff", "heating", `{"ab16Action":"select_summer_heating_operation"}`},
		{"01_day_setpoint_heating", "room", `{"ab16CircuitClass":"heating"}`},
		{"01_day_setpoint_fixed_value", "flow", `{"ab16CircuitClass":"fixed_value"}`},
		{"01_day_setpoint_dhw", "tank", `{"ab16CircuitClass":"dhw"}`},
		{"0A_reduction_heating", "room", `{"ab16CircuitClass":"heating"}`},
		{"0A_reduction_fixed_value", "fixed_flow", `{"ab16CircuitClass":"fixed_value"}`},
		{"09_first_period_start", "programming", `{"ab16ProgrammingMode":0,"ab16OtherModes":"ignored"}`},
		{"07_saving", "heating", `{"ab16Zero":"off"}`},
		{"03_drying_day", "drying", `{"ab16Zero":"disabled"}`},
	}
	if len(ab16Rows(ab16T)) != len(ab16Expected) {
		ab16T.Fatal("finite area changed")
	}
	for _, ab16Case := range ab16Expected {
		ab16T.Run("ab16_"+ab16Case.ab16ID, func(ab16T *ab16testing.T) {
			ab16Actual := ab16Find(ab16T, ab16Case.ab16ID)
			var ab16Relations map[string]any
			if ab16Err := ab16json.Unmarshal([]byte(ab16Case.ab16Relations), &ab16Relations); ab16Err != nil {
				ab16T.Fatal(ab16Err)
			}
			if ab16Actual.Ab16Role != ab16Case.ab16Role || !ab16reflect.DeepEqual(ab16Relations, ab16Actual.Ab16Relations) {
				ab16T.Fatalf("association: %+v", ab16Actual)
			}
			if ab16Actual.Ab16Grade != "Hypothesis" || ab16Actual.Ab16Raw != "retain native bytes and last-known-good unaffected fields; no automatic admission" || len(ab16Actual.Ab16Unknown) != 9 {
				ab16T.Fatal("evidence lost")
			}
			for _, ab16Value := range ab16Actual.Ab16Unknown {
				if ab16Value != "Unknown" {
					ab16T.Fatal("qualification invented")
				}
			}
			ab16Operation := *ab16Actual.Ab16Identity.Ab16Operation
			if ab16Actual.Ab16Identity.Ab16Family != "B505" || ab16Actual.Ab16Identity.Ab16Variant != "operating_roles" || ab16Actual.Ab16Identity.Ab16Primary != 0xB5 || ab16Actual.Ab16Identity.Ab16Secondary == nil || *ab16Actual.Ab16Identity.Ab16Secondary != 5 || ab16Actual.Ab16Identity.Ab16Direction != "Unknown" || ab16Actual.Ab16Identity.Ab16Alternative != "none" || ab16Actual.Ab16Identity.Ab16Group != nil || ab16Actual.Ab16Identity.Ab16Register != nil || ab16Actual.Ab16Identity.Ab16InstanceScope != "dynamic" {
				ab16T.Fatal("identity changed")
			}
			const ab16Hex = "0123456789ABCDEF"
			if ab16Case.ab16ID[:2] != string([]byte{ab16Hex[ab16Operation>>4], ab16Hex[ab16Operation&15]}) || ab16Actual.Ab16Identity.Ab16Field != ab16Case.ab16ID[3:] {
				ab16T.Fatal("operation/field changed")
			}
		})
	}
}

func TestAb16ExactAxes(ab16T *ab16testing.T) {
	ab16Base := ab16Rows(ab16T)[0].Ab16Identity
	ab16Changes := []struct {
		ab16Name   string
		ab16Change func(*Ab16Key)
	}{
		{"ab16Family", func(ab16Key *Ab16Key) { ab16Key.Ab16Family = "synthetic_family" }},
		{"ab16Variant", func(ab16Key *Ab16Key) { ab16Key.Ab16Variant = "synthetic_variant" }},
		{"ab16Primary", func(ab16Key *Ab16Key) { ab16Key.Ab16Primary = 0 }},
		{"ab16Secondary", func(ab16Key *Ab16Key) { ab16Value := uint8(0); ab16Key.Ab16Secondary = &ab16Value }},
		{"ab16Operation", func(ab16Key *Ab16Key) { ab16Value := uint8(0); ab16Key.Ab16Operation = &ab16Value }},
		{"ab16Direction", func(ab16Key *Ab16Key) { ab16Key.Ab16Direction = "synthetic_direction" }},
		{"ab16Field", func(ab16Key *Ab16Key) { ab16Key.Ab16Field = "synthetic_field" }},
		{"ab16Alternative", func(ab16Key *Ab16Key) { ab16Key.Ab16Alternative = "synthetic_alternative" }},
		{"ab16Group", func(ab16Key *Ab16Key) { ab16Value := uint8(0); ab16Key.Ab16Group = &ab16Value }},
		{"ab16InstanceScope", func(ab16Key *Ab16Key) { ab16Key.Ab16InstanceScope = "00" }},
		{"ab16Register", func(ab16Key *Ab16Key) { ab16Value := uint16(0); ab16Key.Ab16Register = &ab16Value }},
	}
	for _, ab16Case := range ab16Changes {
		ab16T.Run(ab16Case.ab16Name, func(ab16T *ab16testing.T) {
			ab16Changed := ab16Base
			ab16Case.ab16Change(&ab16Changed)
			if ab16Raw, ab16Found, ab16Err := Ab16Lookup(ab16Changed); ab16Err != nil || ab16Found || ab16Raw != nil {
				ab16T.Fatal("fallback admitted")
			}
			ab16Fixture, ab16Err := ab16json.Marshal([]any{map[string]any{"identity": ab16Base, "id": "ab16_base"}, map[string]any{"identity": ab16Changed, "id": "ab16_changed"}})
			if ab16Err != nil {
				ab16T.Fatal(ab16Err)
			}
			for _, ab16Key := range []Ab16Key{ab16Base, ab16Changed} {
				ab16Raw, ab16Found, ab16Err := ab16Lookup(ab16Fixture, ab16Key)
				if ab16Err != nil || !ab16Found {
					ab16T.Fatalf("valid independent key: %v", ab16Err)
				}
				var ab16Actual struct {
					Ab16ID string `json:"id"`
				}
				if ab16Err := ab16json.Unmarshal(ab16Raw, &ab16Actual); ab16Err != nil {
					ab16T.Fatal(ab16Err)
				}
				ab16Want := "ab16_base"
				if ab16reflect.DeepEqual(ab16Key, ab16Changed) {
					ab16Want = "ab16_changed"
				}
				if ab16Actual.Ab16ID != ab16Want {
					ab16T.Fatal("axis conflated")
				}
			}
		})
	}
	for _, ab16Axis := range []string{"ab16Secondary", "ab16Operation", "ab16Group", "ab16Register"} {
		ab16T.Run(ab16Axis+"_absent_zero", func(ab16T *ab16testing.T) {
			ab16Absent := ab16Base
			ab16Absent.Ab16Secondary = nil
			ab16Absent.Ab16Operation = nil
			ab16Absent.Ab16Group = nil
			ab16Absent.Ab16Register = nil
			ab16Zero := ab16Absent
			for _, ab16Change := range ab16Changes {
				if ab16Change.ab16Name == ab16Axis {
					ab16Change.ab16Change(&ab16Zero)
				}
			}
			ab16Fixture, ab16Err := ab16json.Marshal([]any{map[string]any{"identity": ab16Absent}, map[string]any{"identity": ab16Zero}})
			if ab16Err != nil {
				ab16T.Fatal(ab16Err)
			}
			ab16Index, ab16Err := ab16Validate(ab16Fixture)
			if ab16Err != nil || len(ab16Index) != 2 {
				ab16T.Fatalf("absent/zero conflated: %v", ab16Err)
			}
			for _, ab16Key := range []Ab16Key{ab16Absent, ab16Zero} {
				if _, ab16Found, ab16Err := ab16Lookup(ab16Fixture, ab16Key); ab16Err != nil || !ab16Found {
					ab16T.Fatal("optional control missing")
				}
			}
		})
	}
}

func TestAb16ValidationAndCopies(ab16T *ab16testing.T) {
	ab16Entries, ab16Err := Ab16Entries()
	if ab16Err != nil {
		ab16T.Fatal(ab16Err)
	}
	ab16Key := ab16Rows(ab16T)[0].Ab16Identity
	ab16Row, ab16Found, ab16Err := Ab16Lookup(ab16Key)
	if ab16Err != nil || !ab16Found {
		ab16T.Fatal("lookup")
	}
	ab16Duplicate := append([]byte("["), ab16Row...)
	ab16Duplicate = append(ab16Duplicate, ',')
	ab16Duplicate = append(ab16Duplicate, ab16Row...)
	ab16Duplicate = append(ab16Duplicate, ']')
	for _, ab16Invalid := range [][]byte{[]byte("{"), []byte("null"), []byte("{}"), []byte("[null]"), []byte("[{}]"), []byte(`[{"identity":{}}]`), ab16Duplicate} {
		if _, ab16Err := ab16Validate(ab16Invalid); ab16Err == nil {
			ab16T.Fatal("invalid data accepted")
		}
		if ab16Raw, ab16Found, ab16Err := ab16Lookup(ab16Invalid, ab16Key); ab16Err == nil || ab16Found || ab16Raw != nil {
			ab16T.Fatal("invalid data admitted")
		}
	}
	ab16SavedEntries := append([]byte(nil), ab16Entries...)
	ab16SavedRow := append([]byte(nil), ab16Row...)
	for ab16Index := range ab16Entries {
		ab16Entries[ab16Index] = 'x'
	}
	for ab16Index := range ab16Row {
		ab16Row[ab16Index] = 'x'
	}
	ab16LaterEntries, ab16Err := Ab16Entries()
	if ab16Err != nil || !ab16bytes.Equal(ab16LaterEntries, ab16SavedEntries) {
		ab16T.Fatal("entries shared")
	}
	ab16LaterRow, ab16Found, ab16Err := Ab16Lookup(ab16Key)
	if ab16Err != nil || !ab16Found || !ab16bytes.Equal(ab16LaterRow, ab16SavedRow) {
		ab16T.Fatal("row shared")
	}
}

func TestAb16SyntheticCandidates(ab16T *ab16testing.T) {
	ab16Tanks := map[string]string{"one": "ab16_alpha", "two": "ab16_beta", "three": "ab16_gamma"}
	for _, ab16Case := range []struct{ ab16ID, ab16Expected string }{
		{"1B_switch_on_cutoff_difference", "ab16_alpha"}, {"1E_switch_on_cutoff_difference", "ab16_beta"}, {"38_switch_on_cutoff_difference", "ab16_gamma"},
		{"1A_maximum_loading_temperature", "ab16_alpha"}, {"1D_maximum_loading_temperature", "ab16_beta"}, {"37_maximum_loading_temperature", "ab16_gamma"},
	} {
		ab16Row := ab16Find(ab16T, ab16Case.ab16ID)
		if ab16Tanks[ab16Row.Ab16Relations["ab16Tank"].(string)] != ab16Case.ab16Expected {
			ab16T.Fatal("tank token misrouted")
		}
		for _, ab16Other := range []string{"one", "two", "three"} {
			if ab16Other != ab16Row.Ab16Relations["ab16Tank"] && ab16Tanks[ab16Other] == ab16Case.ab16Expected {
				ab16T.Fatal("tank controls not distinct")
			}
		}
	}
	ab16Inputs := map[string]float64{"heating_night": 17, "cooling_night": 23, "heating_holiday": 13, "cooling_holiday": 29}
	for _, ab16Case := range []struct {
		ab16ID   string
		ab16Want float64
	}{{"47_heating_night", 17}, {"47_cooling_night", 23}, {"48_heating_holiday", 13}, {"48_cooling_holiday", 29}} {
		ab16Row := ab16Find(ab16T, ab16Case.ab16ID)
		ab16Field := ab16Row.Ab16Relations["ab16Mode"].(string) + "_" + ab16Row.Ab16Relations["ab16Period"].(string)
		if ab16Inputs[ab16Field] != ab16Case.ab16Want {
			ab16T.Fatal("setpoints merged")
		}
	}
	for _, ab16ID := range []string{"5A_unchanged_setpoint", "3D_unchanged_setpoint"} {
		ab16Row := ab16Find(ab16T, ab16ID)
		for _, ab16Replacement := range []bool{true, false} {
			ab16Candidate := any(float64(21))
			ab16Want := any(float64(21))
			if ab16Replacement {
				ab16Want = "ab16_prior"
				if ab16Row.Ab16Relations["ab16ReplacementAction"] == "retain_prior_token" && ab16Row.Ab16Relations["ab16BranchOrder"] == "replacement_before_numeric" {
					ab16Candidate = "ab16_prior"
				}
			}
			if ab16Candidate != ab16Want || ab16Row.Ab16Relations["ab16ReplacementRaw"] != "Unknown" || ab16Row.Ab16Relations["ab16NumericBranch"] != "numeric_candidate_without_scale" {
				ab16T.Fatal("replacement or numeric branch changed")
			}
		}
	}
	ab16DHW := ab16Find(ab16T, "31_ordinary_dhw")
	for _, ab16Enabled := range []bool{false, true} {
		ab16Axis := "ab16BlockedEligibility"
		ab16Want := "frost_only"
		if ab16Enabled {
			ab16Axis = "ab16EnabledEligibility"
			ab16Want = "ordinary_dhw"
		}
		if ab16DHW.Ab16Relations[ab16Axis] != ab16Want {
			ab16T.Fatal("eligibility changed")
		}
	}
	for _, ab16ID := range []string{"18_charging_disable", "27_charging_disable"} {
		ab16Row := ab16Find(ab16T, ab16ID)
		for _, ab16Case := range []struct {
			ab16Role                     string
			ab16Parallel, ab16Suppressed bool
		}{{"mixer", false, true}, {"mixer", true, false}, {"heating", false, true}, {"heating", true, true}} {
			ab16Suppressed := !ab16Case.ab16Parallel || ab16Case.ab16Role != ab16Row.Ab16Relations["ab16ExceptionRole"]
			if ab16Suppressed != ab16Case.ab16Suppressed || ab16Find(ab16T, "14_parallel_charging").Ab16Role != "mixer" {
				ab16T.Fatal("parallel exception broadened")
			}
		}
	}
	ab16Excess := ab16Find(ab16T, "17_flow_excess")
	ab16Candidate := 40
	if ab16Excess.Ab16Relations["ab16Action"] == "add_to_demanded_flow" {
		ab16Candidate += 5
	}
	if ab16Candidate != 45 || ab16Excess.Ab16Relations["ab16Encoding"] != "Unknown" {
		ab16T.Fatal("semantic addition changed")
	}
	for _, ab16ID := range []string{"22_reset", "20_reset"} {
		ab16Row := ab16Find(ab16T, ab16ID)
		if ab16Row.Ab16Relations["ab16Accepted"] != "Unknown" || ab16Row.Ab16Relations["ab16Persistent"] != "Unknown" {
			ab16T.Fatal("request promoted to acceptance")
		}
	}
	if ab16Find(ab16T, "0E_minimum_flow").Ab16Relations["ab16Zero"] != "off" {
		ab16T.Fatal("zero became positive minimum")
	}
	for _, ab16Case := range []struct{ ab16Class, ab16Want string }{{"heating", "room"}, {"fixed_value", "flow"}, {"dhw", "tank"}} {
		if ab16Find(ab16T, "01_day_setpoint_"+ab16Case.ab16Class).Ab16Role != ab16Case.ab16Want {
			ab16T.Fatal("day role changed")
		}
		for _, ab16Mode := range []float64{0, 1} {
			ab16Row := ab16Find(ab16T, "09_first_period_start")
			ab16Result := "ignored"
			if ab16Mode == ab16Row.Ab16Relations["ab16ProgrammingMode"] {
				ab16Result = "ab16_same_start"
			}
			ab16Want := "ignored"
			if ab16Mode == 0 {
				ab16Want = "ab16_same_start"
			}
			if ab16Result != ab16Want || ab16Row.Ab16Relations["ab16OtherModes"] != "ignored" {
				ab16T.Fatal("programming mode changed")
			}
		}
	}
}
