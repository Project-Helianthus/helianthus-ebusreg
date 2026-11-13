package vaillant

import (
	b505OperatingRoleBytes "bytes"
	b505OperatingRoleJSON "encoding/json"
	b505OperatingRoleReflect "reflect"
	b505OperatingRoleTesting "testing"
)

type b505OperatingRoleCatalogRow struct {
	B505OperatingRoleID        string               `json:"id"`
	B505OperatingRoleIdentity  B505OperatingRoleKey `json:"identity"`
	B505OperatingRole          string               `json:"ab16Role"`
	B505OperatingRoleRelations map[string]any       `json:"ab16Relations"`
	B505OperatingRoleGrade     string               `json:"grade"`
	B505OperatingRoleUnknown   map[string]string    `json:"ab16Unknown"`
	B505OperatingRoleRaw       string               `json:"raw_fallback"`
}

func b505OperatingRoleRows(b505OperatingRoleT *b505OperatingRoleTesting.T) []b505OperatingRoleCatalogRow {
	b505OperatingRoleT.Helper()
	b505OperatingRoleData, b505OperatingRoleErr := B505OperatingRoleEntries()
	if b505OperatingRoleErr != nil {
		b505OperatingRoleT.Fatal(b505OperatingRoleErr)
	}
	var b505OperatingRoleRows []b505OperatingRoleCatalogRow
	if b505OperatingRoleErr = b505OperatingRoleJSON.Unmarshal(b505OperatingRoleData, &b505OperatingRoleRows); b505OperatingRoleErr != nil {
		b505OperatingRoleT.Fatal(b505OperatingRoleErr)
	}
	for _, b505OperatingRoleRow := range b505OperatingRoleRows {
		b505OperatingRoleRaw, b505OperatingRoleFound, b505OperatingRoleErr := B505OperatingRoleLookup(b505OperatingRoleRow.B505OperatingRoleIdentity)
		if b505OperatingRoleErr != nil || !b505OperatingRoleFound {
			b505OperatingRoleT.Fatalf("lookup: %v %v", b505OperatingRoleFound, b505OperatingRoleErr)
		}
		var b505OperatingRoleActual b505OperatingRoleCatalogRow
		if b505OperatingRoleErr = b505OperatingRoleJSON.Unmarshal(b505OperatingRoleRaw, &b505OperatingRoleActual); b505OperatingRoleErr != nil || !b505OperatingRoleReflect.DeepEqual(b505OperatingRoleRow, b505OperatingRoleActual) {
			b505OperatingRoleT.Fatalf("whole row: %v", b505OperatingRoleErr)
		}
	}
	return b505OperatingRoleRows
}

func b505OperatingRoleFind(b505OperatingRoleT *b505OperatingRoleTesting.T, b505OperatingRoleID string) b505OperatingRoleCatalogRow {
	b505OperatingRoleT.Helper()
	for _, b505OperatingRoleRow := range b505OperatingRoleRows(b505OperatingRoleT) {
		if b505OperatingRoleRow.B505OperatingRoleID == "ab16_"+b505OperatingRoleID {
			return b505OperatingRoleRow
		}
	}
	b505OperatingRoleT.Fatalf("missing %s", b505OperatingRoleID)
	return b505OperatingRoleCatalogRow{}
}

func TestAb16Associations(b505OperatingRoleT *b505OperatingRoleTesting.T) {
	b505OperatingRoleExpected := []struct{ b505OperatingRoleID, b505OperatingRoleRole, b505OperatingRoleRelations string }{
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
	if len(b505OperatingRoleRows(b505OperatingRoleT)) != len(b505OperatingRoleExpected) {
		b505OperatingRoleT.Fatal("finite area changed")
	}
	for _, b505OperatingRoleCase := range b505OperatingRoleExpected {
		b505OperatingRoleT.Run("ab16_"+b505OperatingRoleCase.b505OperatingRoleID, func(b505OperatingRoleT *b505OperatingRoleTesting.T) {
			b505OperatingRoleActual := b505OperatingRoleFind(b505OperatingRoleT, b505OperatingRoleCase.b505OperatingRoleID)
			var b505OperatingRoleRelations map[string]any
			if b505OperatingRoleErr := b505OperatingRoleJSON.Unmarshal([]byte(b505OperatingRoleCase.b505OperatingRoleRelations), &b505OperatingRoleRelations); b505OperatingRoleErr != nil {
				b505OperatingRoleT.Fatal(b505OperatingRoleErr)
			}
			if b505OperatingRoleActual.B505OperatingRole != b505OperatingRoleCase.b505OperatingRoleRole || !b505OperatingRoleReflect.DeepEqual(b505OperatingRoleRelations, b505OperatingRoleActual.B505OperatingRoleRelations) {
				b505OperatingRoleT.Fatalf("association: %+v", b505OperatingRoleActual)
			}
			if b505OperatingRoleActual.B505OperatingRoleGrade != "Hypothesis" || b505OperatingRoleActual.B505OperatingRoleRaw != "retain native bytes and last-known-good unaffected fields; no automatic admission" || len(b505OperatingRoleActual.B505OperatingRoleUnknown) != 9 {
				b505OperatingRoleT.Fatal("evidence lost")
			}
			for _, b505OperatingRoleValue := range b505OperatingRoleActual.B505OperatingRoleUnknown {
				if b505OperatingRoleValue != "Unknown" {
					b505OperatingRoleT.Fatal("qualification invented")
				}
			}
			b505OperatingRoleOperation := *b505OperatingRoleActual.B505OperatingRoleIdentity.B505OperatingRoleOperation
			if b505OperatingRoleActual.B505OperatingRoleIdentity.B505OperatingRoleFamily != "B505" || b505OperatingRoleActual.B505OperatingRoleIdentity.B505OperatingRoleVariant != "operating_roles" || b505OperatingRoleActual.B505OperatingRoleIdentity.B505OperatingRolePrimary != 0xB5 || b505OperatingRoleActual.B505OperatingRoleIdentity.B505OperatingRoleSecondary == nil || *b505OperatingRoleActual.B505OperatingRoleIdentity.B505OperatingRoleSecondary != 5 || b505OperatingRoleActual.B505OperatingRoleIdentity.B505OperatingRoleDirection != "Unknown" || b505OperatingRoleActual.B505OperatingRoleIdentity.B505OperatingRoleAlternative != "none" || b505OperatingRoleActual.B505OperatingRoleIdentity.B505OperatingRoleGroup != nil || b505OperatingRoleActual.B505OperatingRoleIdentity.B505OperatingRoleRegister != nil || b505OperatingRoleActual.B505OperatingRoleIdentity.B505OperatingRoleInstanceScope != "dynamic" {
				b505OperatingRoleT.Fatal("identity changed")
			}
			const b505OperatingRoleHex = "0123456789ABCDEF"
			if b505OperatingRoleCase.b505OperatingRoleID[:2] != string([]byte{b505OperatingRoleHex[b505OperatingRoleOperation>>4], b505OperatingRoleHex[b505OperatingRoleOperation&15]}) || b505OperatingRoleActual.B505OperatingRoleIdentity.B505OperatingRoleField != b505OperatingRoleCase.b505OperatingRoleID[3:] {
				b505OperatingRoleT.Fatal("operation/field changed")
			}
		})
	}
}

func TestAb16ExactAxes(b505OperatingRoleT *b505OperatingRoleTesting.T) {
	b505OperatingRoleBase := b505OperatingRoleRows(b505OperatingRoleT)[0].B505OperatingRoleIdentity
	b505OperatingRoleChanges := []struct {
		b505OperatingRoleName   string
		b505OperatingRoleChange func(*B505OperatingRoleKey)
	}{
		{"ab16Family", func(b505OperatingRoleKey *B505OperatingRoleKey) {
			b505OperatingRoleKey.B505OperatingRoleFamily = "synthetic_family"
		}},
		{"ab16Variant", func(b505OperatingRoleKey *B505OperatingRoleKey) {
			b505OperatingRoleKey.B505OperatingRoleVariant = "synthetic_variant"
		}},
		{"ab16Primary", func(b505OperatingRoleKey *B505OperatingRoleKey) { b505OperatingRoleKey.B505OperatingRolePrimary = 0 }},
		{"ab16Secondary", func(b505OperatingRoleKey *B505OperatingRoleKey) {
			b505OperatingRoleValue := uint8(0)
			b505OperatingRoleKey.B505OperatingRoleSecondary = &b505OperatingRoleValue
		}},
		{"ab16Operation", func(b505OperatingRoleKey *B505OperatingRoleKey) {
			b505OperatingRoleValue := uint8(0)
			b505OperatingRoleKey.B505OperatingRoleOperation = &b505OperatingRoleValue
		}},
		{"ab16Direction", func(b505OperatingRoleKey *B505OperatingRoleKey) {
			b505OperatingRoleKey.B505OperatingRoleDirection = "synthetic_direction"
		}},
		{"ab16Field", func(b505OperatingRoleKey *B505OperatingRoleKey) {
			b505OperatingRoleKey.B505OperatingRoleField = "synthetic_field"
		}},
		{"ab16Alternative", func(b505OperatingRoleKey *B505OperatingRoleKey) {
			b505OperatingRoleKey.B505OperatingRoleAlternative = "synthetic_alternative"
		}},
		{"ab16Group", func(b505OperatingRoleKey *B505OperatingRoleKey) {
			b505OperatingRoleValue := uint8(0)
			b505OperatingRoleKey.B505OperatingRoleGroup = &b505OperatingRoleValue
		}},
		{"ab16InstanceScope", func(b505OperatingRoleKey *B505OperatingRoleKey) {
			b505OperatingRoleKey.B505OperatingRoleInstanceScope = "00"
		}},
		{"ab16Register", func(b505OperatingRoleKey *B505OperatingRoleKey) {
			b505OperatingRoleValue := uint16(0)
			b505OperatingRoleKey.B505OperatingRoleRegister = &b505OperatingRoleValue
		}},
	}
	for _, b505OperatingRoleCase := range b505OperatingRoleChanges {
		b505OperatingRoleT.Run(b505OperatingRoleCase.b505OperatingRoleName, func(b505OperatingRoleT *b505OperatingRoleTesting.T) {
			b505OperatingRoleChanged := b505OperatingRoleBase
			b505OperatingRoleCase.b505OperatingRoleChange(&b505OperatingRoleChanged)
			if b505OperatingRoleRaw, b505OperatingRoleFound, b505OperatingRoleErr := B505OperatingRoleLookup(b505OperatingRoleChanged); b505OperatingRoleErr != nil || b505OperatingRoleFound || b505OperatingRoleRaw != nil {
				b505OperatingRoleT.Fatal("fallback admitted")
			}
			b505OperatingRoleFixture, b505OperatingRoleErr := b505OperatingRoleJSON.Marshal([]any{map[string]any{"identity": b505OperatingRoleBase, "id": "ab16_base"}, map[string]any{"identity": b505OperatingRoleChanged, "id": "ab16_changed"}})
			if b505OperatingRoleErr != nil {
				b505OperatingRoleT.Fatal(b505OperatingRoleErr)
			}
			for _, b505OperatingRoleKey := range []B505OperatingRoleKey{b505OperatingRoleBase, b505OperatingRoleChanged} {
				b505OperatingRoleRaw, b505OperatingRoleFound, b505OperatingRoleErr := b505OperatingRoleLookup(b505OperatingRoleFixture, b505OperatingRoleKey)
				if b505OperatingRoleErr != nil || !b505OperatingRoleFound {
					b505OperatingRoleT.Fatalf("valid independent key: %v", b505OperatingRoleErr)
				}
				var b505OperatingRoleActual struct {
					B505OperatingRoleID string `json:"id"`
				}
				if b505OperatingRoleErr := b505OperatingRoleJSON.Unmarshal(b505OperatingRoleRaw, &b505OperatingRoleActual); b505OperatingRoleErr != nil {
					b505OperatingRoleT.Fatal(b505OperatingRoleErr)
				}
				b505OperatingRoleWant := "ab16_base"
				if b505OperatingRoleReflect.DeepEqual(b505OperatingRoleKey, b505OperatingRoleChanged) {
					b505OperatingRoleWant = "ab16_changed"
				}
				if b505OperatingRoleActual.B505OperatingRoleID != b505OperatingRoleWant {
					b505OperatingRoleT.Fatal("axis conflated")
				}
			}
		})
	}
	for _, b505OperatingRoleAxis := range []string{"ab16Secondary", "ab16Operation", "ab16Group", "ab16Register"} {
		b505OperatingRoleT.Run(b505OperatingRoleAxis+"_absent_zero", func(b505OperatingRoleT *b505OperatingRoleTesting.T) {
			b505OperatingRoleAbsent := b505OperatingRoleBase
			b505OperatingRoleAbsent.B505OperatingRoleSecondary = nil
			b505OperatingRoleAbsent.B505OperatingRoleOperation = nil
			b505OperatingRoleAbsent.B505OperatingRoleGroup = nil
			b505OperatingRoleAbsent.B505OperatingRoleRegister = nil
			b505OperatingRoleZero := b505OperatingRoleAbsent
			for _, b505OperatingRoleChange := range b505OperatingRoleChanges {
				if b505OperatingRoleChange.b505OperatingRoleName == b505OperatingRoleAxis {
					b505OperatingRoleChange.b505OperatingRoleChange(&b505OperatingRoleZero)
				}
			}
			b505OperatingRoleFixture, b505OperatingRoleErr := b505OperatingRoleJSON.Marshal([]any{map[string]any{"identity": b505OperatingRoleAbsent}, map[string]any{"identity": b505OperatingRoleZero}})
			if b505OperatingRoleErr != nil {
				b505OperatingRoleT.Fatal(b505OperatingRoleErr)
			}
			b505OperatingRoleIndex, b505OperatingRoleErr := b505OperatingRoleValidate(b505OperatingRoleFixture)
			if b505OperatingRoleErr != nil || len(b505OperatingRoleIndex) != 2 {
				b505OperatingRoleT.Fatalf("absent/zero conflated: %v", b505OperatingRoleErr)
			}
			for _, b505OperatingRoleKey := range []B505OperatingRoleKey{b505OperatingRoleAbsent, b505OperatingRoleZero} {
				if _, b505OperatingRoleFound, b505OperatingRoleErr := b505OperatingRoleLookup(b505OperatingRoleFixture, b505OperatingRoleKey); b505OperatingRoleErr != nil || !b505OperatingRoleFound {
					b505OperatingRoleT.Fatal("optional control missing")
				}
			}
		})
	}
}

func TestAb16ValidationAndCopies(b505OperatingRoleT *b505OperatingRoleTesting.T) {
	b505OperatingRoleEntries, b505OperatingRoleErr := B505OperatingRoleEntries()
	if b505OperatingRoleErr != nil {
		b505OperatingRoleT.Fatal(b505OperatingRoleErr)
	}
	b505OperatingRoleKey := b505OperatingRoleRows(b505OperatingRoleT)[0].B505OperatingRoleIdentity
	b505OperatingRoleRow, b505OperatingRoleFound, b505OperatingRoleErr := B505OperatingRoleLookup(b505OperatingRoleKey)
	if b505OperatingRoleErr != nil || !b505OperatingRoleFound {
		b505OperatingRoleT.Fatal("lookup")
	}
	b505OperatingRoleDuplicate := append([]byte("["), b505OperatingRoleRow...)
	b505OperatingRoleDuplicate = append(b505OperatingRoleDuplicate, ',')
	b505OperatingRoleDuplicate = append(b505OperatingRoleDuplicate, b505OperatingRoleRow...)
	b505OperatingRoleDuplicate = append(b505OperatingRoleDuplicate, ']')
	for _, b505OperatingRoleInvalid := range [][]byte{[]byte("{"), []byte("null"), []byte("{}"), []byte("[null]"), []byte("[{}]"), []byte(`[{"identity":{}}]`), b505OperatingRoleDuplicate} {
		if _, b505OperatingRoleErr := b505OperatingRoleValidate(b505OperatingRoleInvalid); b505OperatingRoleErr == nil {
			b505OperatingRoleT.Fatal("invalid data accepted")
		}
		if b505OperatingRoleRaw, b505OperatingRoleFound, b505OperatingRoleErr := b505OperatingRoleLookup(b505OperatingRoleInvalid, b505OperatingRoleKey); b505OperatingRoleErr == nil || b505OperatingRoleFound || b505OperatingRoleRaw != nil {
			b505OperatingRoleT.Fatal("invalid data admitted")
		}
	}
	b505OperatingRoleSavedEntries := append([]byte(nil), b505OperatingRoleEntries...)
	b505OperatingRoleSavedRow := append([]byte(nil), b505OperatingRoleRow...)
	for b505OperatingRoleIndex := range b505OperatingRoleEntries {
		b505OperatingRoleEntries[b505OperatingRoleIndex] = 'x'
	}
	for b505OperatingRoleIndex := range b505OperatingRoleRow {
		b505OperatingRoleRow[b505OperatingRoleIndex] = 'x'
	}
	b505OperatingRoleLaterEntries, b505OperatingRoleErr := B505OperatingRoleEntries()
	if b505OperatingRoleErr != nil || !b505OperatingRoleBytes.Equal(b505OperatingRoleLaterEntries, b505OperatingRoleSavedEntries) {
		b505OperatingRoleT.Fatal("entries shared")
	}
	b505OperatingRoleLaterRow, b505OperatingRoleFound, b505OperatingRoleErr := B505OperatingRoleLookup(b505OperatingRoleKey)
	if b505OperatingRoleErr != nil || !b505OperatingRoleFound || !b505OperatingRoleBytes.Equal(b505OperatingRoleLaterRow, b505OperatingRoleSavedRow) {
		b505OperatingRoleT.Fatal("row shared")
	}
}

func TestAb16SyntheticCandidates(b505OperatingRoleT *b505OperatingRoleTesting.T) {
	b505OperatingRoleTanks := map[string]string{"one": "ab16_alpha", "two": "ab16_beta", "three": "ab16_gamma"}
	for _, b505OperatingRoleCase := range []struct{ b505OperatingRoleID, b505OperatingRoleExpected string }{
		{"1B_switch_on_cutoff_difference", "ab16_alpha"}, {"1E_switch_on_cutoff_difference", "ab16_beta"}, {"38_switch_on_cutoff_difference", "ab16_gamma"},
		{"1A_maximum_loading_temperature", "ab16_alpha"}, {"1D_maximum_loading_temperature", "ab16_beta"}, {"37_maximum_loading_temperature", "ab16_gamma"},
	} {
		b505OperatingRoleRow := b505OperatingRoleFind(b505OperatingRoleT, b505OperatingRoleCase.b505OperatingRoleID)
		if b505OperatingRoleTanks[b505OperatingRoleRow.B505OperatingRoleRelations["ab16Tank"].(string)] != b505OperatingRoleCase.b505OperatingRoleExpected {
			b505OperatingRoleT.Fatal("tank token misrouted")
		}
		for _, b505OperatingRoleOther := range []string{"one", "two", "three"} {
			if b505OperatingRoleOther != b505OperatingRoleRow.B505OperatingRoleRelations["ab16Tank"] && b505OperatingRoleTanks[b505OperatingRoleOther] == b505OperatingRoleCase.b505OperatingRoleExpected {
				b505OperatingRoleT.Fatal("tank controls not distinct")
			}
		}
	}
	b505OperatingRoleInputs := map[string]float64{"heating_night": 17, "cooling_night": 23, "heating_holiday": 13, "cooling_holiday": 29}
	for _, b505OperatingRoleCase := range []struct {
		b505OperatingRoleID   string
		b505OperatingRoleWant float64
	}{{"47_heating_night", 17}, {"47_cooling_night", 23}, {"48_heating_holiday", 13}, {"48_cooling_holiday", 29}} {
		b505OperatingRoleRow := b505OperatingRoleFind(b505OperatingRoleT, b505OperatingRoleCase.b505OperatingRoleID)
		b505OperatingRoleField := b505OperatingRoleRow.B505OperatingRoleRelations["ab16Mode"].(string) + "_" + b505OperatingRoleRow.B505OperatingRoleRelations["ab16Period"].(string)
		if b505OperatingRoleInputs[b505OperatingRoleField] != b505OperatingRoleCase.b505OperatingRoleWant {
			b505OperatingRoleT.Fatal("setpoints merged")
		}
	}
	for _, b505OperatingRoleID := range []string{"5A_unchanged_setpoint", "3D_unchanged_setpoint"} {
		b505OperatingRoleRow := b505OperatingRoleFind(b505OperatingRoleT, b505OperatingRoleID)
		for _, b505OperatingRoleReplacement := range []bool{true, false} {
			b505OperatingRoleCandidate := any(float64(21))
			b505OperatingRoleWant := any(float64(21))
			if b505OperatingRoleReplacement {
				b505OperatingRoleWant = "ab16_prior"
				if b505OperatingRoleRow.B505OperatingRoleRelations["ab16ReplacementAction"] == "retain_prior_token" && b505OperatingRoleRow.B505OperatingRoleRelations["ab16BranchOrder"] == "replacement_before_numeric" {
					b505OperatingRoleCandidate = "ab16_prior"
				}
			}
			if b505OperatingRoleCandidate != b505OperatingRoleWant || b505OperatingRoleRow.B505OperatingRoleRelations["ab16ReplacementRaw"] != "Unknown" || b505OperatingRoleRow.B505OperatingRoleRelations["ab16NumericBranch"] != "numeric_candidate_without_scale" {
				b505OperatingRoleT.Fatal("replacement or numeric branch changed")
			}
		}
	}
	b505OperatingRoleDHW := b505OperatingRoleFind(b505OperatingRoleT, "31_ordinary_dhw")
	for _, b505OperatingRoleEnabled := range []bool{false, true} {
		b505OperatingRoleAxis := "ab16BlockedEligibility"
		b505OperatingRoleWant := "frost_only"
		if b505OperatingRoleEnabled {
			b505OperatingRoleAxis = "ab16EnabledEligibility"
			b505OperatingRoleWant = "ordinary_dhw"
		}
		if b505OperatingRoleDHW.B505OperatingRoleRelations[b505OperatingRoleAxis] != b505OperatingRoleWant {
			b505OperatingRoleT.Fatal("eligibility changed")
		}
	}
	for _, b505OperatingRoleID := range []string{"18_charging_disable", "27_charging_disable"} {
		b505OperatingRoleRow := b505OperatingRoleFind(b505OperatingRoleT, b505OperatingRoleID)
		for _, b505OperatingRoleCase := range []struct {
			b505OperatingRoleRole                                  string
			b505OperatingRoleParallel, b505OperatingRoleSuppressed bool
		}{{"mixer", false, true}, {"mixer", true, false}, {"heating", false, true}, {"heating", true, true}} {
			b505OperatingRoleSuppressed := !b505OperatingRoleCase.b505OperatingRoleParallel || b505OperatingRoleCase.b505OperatingRoleRole != b505OperatingRoleRow.B505OperatingRoleRelations["ab16ExceptionRole"]
			if b505OperatingRoleSuppressed != b505OperatingRoleCase.b505OperatingRoleSuppressed || b505OperatingRoleFind(b505OperatingRoleT, "14_parallel_charging").B505OperatingRole != "mixer" {
				b505OperatingRoleT.Fatal("parallel exception broadened")
			}
		}
	}
	b505OperatingRoleExcess := b505OperatingRoleFind(b505OperatingRoleT, "17_flow_excess")
	b505OperatingRoleCandidate := 40
	if b505OperatingRoleExcess.B505OperatingRoleRelations["ab16Action"] == "add_to_demanded_flow" {
		b505OperatingRoleCandidate += 5
	}
	if b505OperatingRoleCandidate != 45 || b505OperatingRoleExcess.B505OperatingRoleRelations["ab16Encoding"] != "Unknown" {
		b505OperatingRoleT.Fatal("semantic addition changed")
	}
	for _, b505OperatingRoleID := range []string{"22_reset", "20_reset"} {
		b505OperatingRoleRow := b505OperatingRoleFind(b505OperatingRoleT, b505OperatingRoleID)
		if b505OperatingRoleRow.B505OperatingRoleRelations["ab16Accepted"] != "Unknown" || b505OperatingRoleRow.B505OperatingRoleRelations["ab16Persistent"] != "Unknown" {
			b505OperatingRoleT.Fatal("request promoted to acceptance")
		}
	}
	if b505OperatingRoleFind(b505OperatingRoleT, "0E_minimum_flow").B505OperatingRoleRelations["ab16Zero"] != "off" {
		b505OperatingRoleT.Fatal("zero became positive minimum")
	}
	for _, b505OperatingRoleCase := range []struct{ b505OperatingRoleClass, b505OperatingRoleWant string }{{"heating", "room"}, {"fixed_value", "flow"}, {"dhw", "tank"}} {
		if b505OperatingRoleFind(b505OperatingRoleT, "01_day_setpoint_"+b505OperatingRoleCase.b505OperatingRoleClass).B505OperatingRole != b505OperatingRoleCase.b505OperatingRoleWant {
			b505OperatingRoleT.Fatal("day role changed")
		}
		for _, b505OperatingRoleMode := range []float64{0, 1} {
			b505OperatingRoleRow := b505OperatingRoleFind(b505OperatingRoleT, "09_first_period_start")
			b505OperatingRoleResult := "ignored"
			if b505OperatingRoleMode == b505OperatingRoleRow.B505OperatingRoleRelations["ab16ProgrammingMode"] {
				b505OperatingRoleResult = "ab16_same_start"
			}
			b505OperatingRoleWant := "ignored"
			if b505OperatingRoleMode == 0 {
				b505OperatingRoleWant = "ab16_same_start"
			}
			if b505OperatingRoleResult != b505OperatingRoleWant || b505OperatingRoleRow.B505OperatingRoleRelations["ab16OtherModes"] != "ignored" {
				b505OperatingRoleT.Fatal("programming mode changed")
			}
		}
	}
}
