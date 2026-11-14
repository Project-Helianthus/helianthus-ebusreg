package vaillant

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func holidayDecode(holidayT *testing.T) map[string]any {
	holidayT.Helper()
	var holidayDoc map[string]any
	if holidayErr := json.Unmarshal(HolidayCatalog(), &holidayDoc); holidayErr != nil {
		holidayT.Fatal(holidayErr)
	}
	return holidayDoc
}

func holidayExpectedRow(holidayID string, holidayChanges map[string]any) map[string]any {
	holidayRow := map[string]any{
		"id": "ab27_" + holidayID, "family": "Vaillant", "variant": nil, "direction": nil,
		"service": nil, "operation": nil, "alternative": "field-declaration",
		"logical_positions": nil, "extent_bytes": nil, "raw_domain": nil,
		"signedness": "Unknown", "byte_order": "Unknown", "conversion": "Unknown", "unit": "Unknown",
		"default": nil, "mappings": map[string]any{}, "reserved_regions": nil,
		"grade": "Hypothesis", "raw_fallback": "preserve raw bytes and unlabeled values",
	}
	for holidayKey, holidayValue := range holidayChanges {
		holidayRow[holidayKey] = holidayValue
	}
	return holidayRow
}

func holidayExpectedRows() []any {
	holidayRows := []any{
		holidayExpectedRow("holiday_number", map[string]any{"extent_bytes": float64(1), "raw_domain": []any{float64(0), float64(255)}, "signedness": "unsigned", "conversion": "identity", "default": float64(0)}),
	}
	for _, holidayEnd := range []string{"start", "end"} {
		for _, holidayPart := range []string{"day", "month", "year"} {
			var holidayDomain any
			if holidayPart == "day" {
				holidayDomain = []any{float64(1), float64(31)}
			}
			if holidayPart == "month" {
				holidayDomain = []any{float64(1), float64(12)}
			}
			holidayRows = append(holidayRows, holidayExpectedRow(holidayEnd+"_"+holidayPart, map[string]any{"default": float64(0), "calendar_domain": holidayDomain, "default_conflict": holidayPart != "year"}))
		}
	}
	holidayRows = append(holidayRows,
		holidayExpectedRow("op02_setpoint", map[string]any{"operation": "OP02", "extent_bytes": float64(1), "raw_domain": []any{float64(0), float64(255)}, "signedness": "unsigned", "conversion": "identity", "unit": "degrees C", "default": float64(0), "candidate_bounds": []any{float64(5), float64(30)}, "default_conflict": true, "off_frost_deletion": "Unknown"}),
		holidayExpectedRow("room_action", map[string]any{"default": float64(0), "mappings": map[string]any{"0": "no-action", "1": "switch-on", "2": "thermostat"}}),
		holidayExpectedRow("holiday_status", map[string]any{"default": float64(0), "mappings": map[string]any{"0": "0", "1": "1"}}),
	)
	for _, holidayDigit := range []string{"1", "2", "3", "4"} {
		holidayRows = append(holidayRows, holidayExpectedRow("keycode_"+holidayDigit, map[string]any{"extent_bytes": float64(1), "raw_domain": []any{float64(0), float64(255)}, "signedness": "unsigned", "conversion": "independent numeric byte", "unit": "numeric", "excluded_interpretations": []any{"ASCII", "packed decimal", "u32"}}))
	}
	return append(holidayRows,
		holidayExpectedRow("circuit_name", map[string]any{"extent_bytes": float64(10), "raw_domain": "ASCII characters", "conversion": "ASCII", "unit": "characters", "text_limitations": []any{"padding", "trimming", "NUL termination", "UTF-8", "truncation"}}),
		holidayExpectedRow("remote_control_write", map[string]any{"variant": "remote-control", "direction": "request", "alternative": "remote-control-write"}),
		holidayExpectedRow("b50d_ui_direction_write", map[string]any{"variant": "UI-direction", "direction": "request", "service": "B50D", "alternative": "UI-direction-write"}),
	)
}

func TestAb27ExactDeclarations(holidayT *testing.T) {
	holidayDoc := holidayDecode(holidayT)
	holidayRows := holidayDoc["entries"].([]any)
	holidayExpected := holidayExpectedRows()
	if len(holidayRows) != 17 {
		holidayT.Fatalf("alternative count: %d", len(holidayRows))
	}
	for holidayIndex, holidayWant := range holidayExpected {
		holidayGot := holidayRows[holidayIndex].(map[string]any)
		if len(holidayGot) != len(holidayWant.(map[string]any)) {
			holidayT.Errorf("row %d property count", holidayIndex)
		}
		for holidayKey, holidayValue := range holidayWant.(map[string]any) {
			if !reflect.DeepEqual(holidayGot[holidayKey], holidayValue) {
				holidayT.Errorf("row %d %s: got %#v want %#v", holidayIndex, holidayKey, holidayGot[holidayKey], holidayValue)
			}
			if _, holidayOk := holidayGot[holidayKey]; !holidayOk {
				holidayT.Errorf("row %d missing %s", holidayIndex, holidayKey)
			}
		}
	}
}

func TestAb27Boundaries(holidayT *testing.T) {
	holidayDoc := holidayDecode(holidayT)
	holidayExpected := map[string]any{
		"catalog": "ab27_holiday", "representation_grade": "Proven",
		"coordinate_convention": "logical positions include prefixes; not physical offsets",
		"qualification":         "Unknown", "runtime_registration": false, "wire_builder": false,
		"limitations": []any{"calendar century", "timezone", "valid-date rules", "start-before-end rules", "atomicity", "persistent assignment", "physical response", "readback", "pairing of B50B writes with B50C reads", "request/positive/negative field layouts Unknown", "reserved regions Unknown", "zero day/month conflicts with calendar domains; no delete, disable or reset meaning"},
	}
	if len(holidayDoc) != len(holidayExpected)+1 {
		holidayT.Fatal("unexpected catalog properties")
	}
	for holidayKey, holidayWant := range holidayExpected {
		if !reflect.DeepEqual(holidayDoc[holidayKey], holidayWant) {
			holidayT.Errorf("%s: %#v", holidayKey, holidayDoc[holidayKey])
		}
	}
}

func TestAb27Copy(holidayT *testing.T) {
	holidayFirst := HolidayCatalog()
	holidayControl := string(holidayFirst)
	for holidayIndex := range holidayFirst {
		holidayFirst[holidayIndex] = 'x'
	}
	if string(HolidayCatalog()) != holidayControl {
		holidayT.Fatal("returned bytes alias embedded catalog")
	}
}

func TestAb27IdentityAxes(holidayT *testing.T) {
	holidayRows := holidayDecode(holidayT)["entries"].([]any)
	for _, holidayAxis := range []string{"family", "variant", "direction", "service", "operation", "alternative"} {
		for holidayIndex, holidayEntry := range holidayRows {
			holidayOriginal := holidayEntry.(map[string]any)
			for _, holidayReplacement := range []any{"ab27_changed", nil} {
				holidayCopy := map[string]any{}
				for holidayKey, holidayValue := range holidayOriginal {
					holidayCopy[holidayKey] = holidayValue
				}
				holidayCopy[holidayAxis] = holidayReplacement
				holidayBytes, holidayErr := json.Marshal(holidayCopy)
				if holidayErr != nil {
					holidayT.Fatal(holidayErr)
				}
				var holidayRoundtrip map[string]any
				if holidayErr = json.Unmarshal(holidayBytes, &holidayRoundtrip); holidayErr != nil {
					holidayT.Fatal(holidayErr)
				}
				if !reflect.DeepEqual(holidayRoundtrip, holidayCopy) {
					holidayT.Fatalf("row %d lost %s", holidayIndex, holidayAxis)
				}
				if !reflect.DeepEqual(holidayReplacement, holidayOriginal[holidayAxis]) && reflect.DeepEqual(holidayRoundtrip, holidayOriginal) {
					holidayT.Fatalf("row %d collapsed %s", holidayIndex, holidayAxis)
				}
			}
		}
	}
}

func TestAb27SyntheticLimits(holidayT *testing.T) {
	holidayRows := holidayDecode(holidayT)["entries"].([]any)
	for _, holidayIndex := range []int{1, 2, 4, 5} {
		holidayRow := holidayRows[holidayIndex].(map[string]any)
		holidayDomain := holidayRow["calendar_domain"].([]any)
		if holidayRow["default"] != float64(0) || holidayRow["default_conflict"] != true || holidayDomain[0].(float64) <= 0 {
			holidayT.Fatal("zero calendar conflict lost")
		}
	}
	holidayBounds := holidayRows[7].(map[string]any)["candidate_bounds"].([]any)
	for holidayIndex, holidayValue := range []float64{0, 5, 30, 31} {
		holidayInside := holidayValue >= holidayBounds[0].(float64) && holidayValue <= holidayBounds[1].(float64)
		if holidayInside != []bool{false, true, true, false}[holidayIndex] {
			holidayT.Errorf("setpoint %v", holidayValue)
		}
	}
	for holidayIndex, holidayValue := range []string{"0", "1", "2", "3"} {
		_, holidayKnown := holidayRows[8].(map[string]any)["mappings"].(map[string]any)[holidayValue]
		if holidayKnown != (holidayIndex < 3) {
			holidayT.Errorf("room action %s", holidayValue)
		}
	}
	for holidayIndex, holidayValue := range []string{"0", "1", "2"} {
		_, holidayKnown := holidayRows[9].(map[string]any)["mappings"].(map[string]any)[holidayValue]
		if holidayKnown != (holidayIndex < 2) {
			holidayT.Errorf("holiday status %s", holidayValue)
		}
	}
	for holidayIndex, holidayValue := range []byte{1, 2, 3, 4} {
		holidayRow := holidayRows[10+holidayIndex].(map[string]any)
		holidayDomain := holidayRow["raw_domain"].([]any)
		if float64(holidayValue) < holidayDomain[0].(float64) || float64(holidayValue) > holidayDomain[1].(float64) || holidayRow["conversion"] != "independent numeric byte" || holidayRow["extent_bytes"] != float64(1) {
			holidayT.Fatalf("key byte %d", holidayIndex)
		}
	}
	holidayWidth := int(holidayRows[14].(map[string]any)["extent_bytes"].(float64))
	for holidayIndex, holidayText := range []string{"ABCDEFGHI", "ABCDEFGHIJ", "ABCDEFGHIJK", "ABCDEFGHIé"} {
		holidayASCII := strings.IndexFunc(holidayText, func(holidayRune rune) bool { return holidayRune > 127 }) < 0
		holidayMatches := len(holidayText) == holidayWidth && holidayASCII
		if holidayMatches != (holidayIndex == 1) {
			holidayT.Errorf("ASCII extent case %d", holidayIndex)
		}
	}
}

func TestRenamedFieldJSONKeys(t *testing.T) {
	cases := []struct {
		value any
		keys  map[string]string
	}{
		{B503HistoryBranchContext{}, map[string]string{"B503HistoryBranchFamily": "family", "B503HistoryBranchVariant": "variant", "B503HistoryBranchOperation": "operation", "B503HistoryBranchDirection": "direction"}},
		{B503HistoryBranchCell{}, map[string]string{"B503HistoryBranchRawPresent": "raw_present", "B503HistoryBranchRaw": "raw", "B503HistoryBranchNormalizedPresent": "normalized_present", "B503HistoryBranchValue": "value"}},
		{B503HistoryBranchFields{}, map[string]string{"B503HistoryBranchContext": "context", "B503HistoryBranchMarkerPresent": "marker_present", "B503HistoryBranchMarker": "marker", "B503HistoryBranchRemainderPresent": "remainder_present", "B503HistoryBranchRemainder": "remainder", "B503HistoryBranchErrorPresent": "error_present", "B503HistoryBranchError": "error", "B503HistoryBranchCalendar": "calendar"}},
		{B503HistoryBranchComponent{}, map[string]string{"B503HistoryBranchCell": "cell", "B503HistoryBranchState": "state"}},
		{B503HistoryBranchView{}, map[string]string{"B503HistoryBranchFields": "fields", "B503HistoryBranchState": "state", "B503HistoryBranchComponents": "components", "B503HistoryBranchComplete": "complete", "B503HistoryBranchContract": "contract", "B503HistoryBranchMeaning": "meaning", "B503HistoryBranchQualification": "qualification"}},
		{B511OverlapViewSelection{}, map[string]string{"B511OverlapViewFamily": "family", "B511OverlapViewVariant": "variant", "B511OverlapViewOperation": "operation", "B511OverlapViewDirection": "direction"}},
		{B511OverlapView{}, map[string]string{"B511OverlapViewSelection": "selection", "B511OverlapViewRaw": "raw", "B511OverlapViewAvailableExtent": "available_extent", "B511OverlapViewRequiredExtent": "required_extent", "B511OverlapViewCompleteLogicalExtent": "complete_logical_extent", "B511OverlapViewSelected": "selected", "B511OverlapViewRawFlag": "raw_flag", "B511OverlapViewFanRunning": "fan_running", "B511OverlapViewRawOverlay": "raw_overlay", "B511OverlapViewRoomSetpointCelsius": "room_setpoint_celsius", "B511OverlapViewRawPressure": "raw_pressure", "B511OverlapViewPressurePercent": "pressure_percent", "B511OverlapViewPressureDomain": "pressure_domain", "B511OverlapViewCandidateGrade": "candidate_grade", "B511OverlapViewPhysicalQualification": "physical_qualification"}},
		{B512DHWViewScalar{}, map[string]string{"B512DHWViewRaw": "raw", "B512DHWViewValue": "value", "B512DHWViewUnit": "unit", "B512DHWViewUnitConflict": "unit_conflict", "B512DHWViewNamedDomain": "named_domain"}},
		{B512DHWViewRequest{}, map[string]string{"B512DHWViewTaskActive": "task_active", "B512DHWViewHeatDemand": "heat_demand", "B512DHWViewLoadingActive": "loading_active", "B512DHWViewLoadingPumpFault": "loading_pump_fault", "B512DHWViewSensorConditionBits": "sensor_condition_bits", "B512DHWViewModulation": "modulation", "B512DHWViewFlow": "flow", "B512DHWViewTemperatures": "temperatures", "B512DHWViewInletConditionBits": "inlet_condition_bits"}},
		{B512DHWViewPositive{}, map[string]string{"B512DHWViewLoadingDisabled": "loading_disabled", "B512DHWViewInternalTestActive": "internal_test_active", "B512DHWViewReady": "ready", "B512DHWViewFrostProtectionActive": "frost_protection_active", "B512DHWViewLegionellaDemand": "legionella_demand", "B512DHWViewReset": "reset", "B512DHWViewSetpoint": "setpoint", "B512DHWViewReturnTemperature": "return_temperature", "B512DHWViewAdjustment": "adjustment"}},
		{B512DHWViewResult{}, map[string]string{"B512DHWViewRaw": "raw", "B512DHWViewUnmodelled": "unmodelled", "B512DHWViewAdmitted": "admitted", "B512DHWViewComplete": "complete", "B512DHWViewContractGrade": "contract_grade", "B512DHWViewFieldGrade": "field_grade", "B512DHWViewQualification": "qualification", "B512DHWViewPresence": "presence", "B512DHWViewFreshness": "freshness", "B512DHWViewExecution": "execution", "B512DHWViewDefaults": "defaults", "B512DHWViewRequest": "request", "B512DHWViewPositive": "positive"}},
		{B516ClockWordViewSelector{}, map[string]string{"B516ClockWordViewFamily": "family", "B516ClockWordViewVariant": "variant", "B516ClockWordViewOperation": "operation", "B516ClockWordViewDirection": "direction"}},
		{B516ClockWordViewCell{}, map[string]string{"B516ClockWordViewRaw": "raw", "B516ClockWordViewFailed": "failed"}},
		{B516ClockWordViewField{}, map[string]string{"B516ClockWordViewState": "state", "B516ClockWordViewValue": "value"}},
		{B516ClockWordViewCandidate{}, map[string]string{"B516ClockWordViewAdmitted": "admitted", "B516ClockWordViewContractGrade": "contract_grade", "B516ClockWordViewFieldGrade": "field_grade", "B516ClockWordViewNativeSupport": "native_support", "B516ClockWordViewPhysicalAdaptation": "physical_adaptation", "B516ClockWordViewPresence": "presence", "B516ClockWordViewFreshness": "freshness", "B516ClockWordViewExecution": "execution"}},
		{B516ClockWordViewClockInput{}, map[string]string{"B516ClockWordViewClock": "clock", "B516ClockWordViewSetpoint": "setpoint", "B516ClockWordViewMode": "mode", "B516ClockWordViewSpecialFunction": "special_function", "B516ClockWordViewMessages": "messages"}},
		{B516ClockWordViewClockResult{}, map[string]string{"B516ClockWordViewCandidate": "candidate", "B516ClockWordViewRaw": "raw", "B516ClockWordViewClock": "clock", "B516ClockWordViewSetpointCelsius": "setpoint_celsius", "B516ClockWordViewTablesGrade": "tables_grade"}},
		{B516ClockWordViewWordResult{}, map[string]string{"B516ClockWordViewCandidate": "candidate", "B516ClockWordViewRaw": "raw", "B516ClockWordViewWords": "words", "B516ClockWordViewPositionsGrade": "positions_grade"}},
		{B505OperatingRoleKey{}, map[string]string{"B505OperatingRoleFamily": "family", "B505OperatingRoleVariant": "variant", "B505OperatingRolePrimary": "primary", "B505OperatingRoleSecondary": "secondary", "B505OperatingRoleOperation": "operation", "B505OperatingRoleDirection": "direction", "B505OperatingRoleField": "field", "B505OperatingRoleAlternative": "alternative", "B505OperatingRoleGroup": "group", "B505OperatingRoleInstanceScope": "instance_scope", "B505OperatingRoleRegister": "register"}},
		{B516TwoWordViewResult{}, map[string]string{"B516TwoWordViewDeclaredCount": "declared_count", "B516TwoWordViewAvailableExtent": "available_extent", "B516TwoWordViewRaw": "raw", "B516TwoWordViewTrailing": "trailing", "B516TwoWordViewAdmitted": "admitted", "B516TwoWordViewStatusAvailable": "status_available", "B516TwoWordViewStatus": "status", "B516TwoWordViewWord0Available": "word0_available", "B516TwoWordViewWord0Raw": "word0_raw", "B516TwoWordViewWord0Bits": "word0_bits", "B516TwoWordViewWord0Finite": "word0_finite", "B516TwoWordViewWord0Value": "word0_value", "B516TwoWordViewWord1Available": "word1_available", "B516TwoWordViewWord1Raw": "word1_raw", "B516TwoWordViewCompleteTwoWord": "complete_two_word", "B516TwoWordViewContractGrade": "contract_grade", "B516TwoWordViewInterpretationGrade": "interpretation_grade", "B516TwoWordViewStatusMeaning": "status_meaning", "B516TwoWordViewWord0NativeUnit": "word0_native_unit", "B516TwoWordViewWord0NativeSign": "word0_native_sign", "B516TwoWordViewWord1Type": "word1_type", "B516TwoWordViewWord1Scale": "word1_scale", "B516TwoWordViewWord1Sign": "word1_sign", "B516TwoWordViewQualification": "qualification"}},
		{B531LimitNumericViewField{}, map[string]string{"B531LimitNumericViewRaw": "raw", "B531LimitNumericViewState": "state", "B531LimitNumericViewAvailable": "available", "B531LimitNumericViewValue": "value", "B531LimitNumericViewFresh": "fresh"}},
		{B531LimitNumericViewPair{}, map[string]string{"B531LimitNumericViewMinimum": "minimum", "B531LimitNumericViewMaximum": "maximum", "B531LimitNumericViewFresh": "fresh"}},
		{B531LimitNumericViewResult{}, map[string]string{"B531LimitNumericViewAdmitted": "admitted", "B531LimitNumericViewAttempt": "attempt", "B531LimitNumericViewPrior": "prior", "B531LimitNumericViewLastKnown": "last_known", "B531LimitNumericViewContractGrade": "contract_grade", "B531LimitNumericViewInterpretationGrade": "interpretation_grade", "B531LimitNumericViewNativeUnit": "native_unit", "B531LimitNumericViewQualification": "qualification", "B531LimitNumericViewPresence": "presence", "B531LimitNumericViewExecution": "execution"}},
		{B51BLocalBudgetCRCBudget{}, map[string]string{"B51BLocalBudgetCRCProcedure": "procedure", "B51BLocalBudgetCRCChunks": "chunks", "B51BLocalBudgetCRCConsumed": "consumed", "B51BLocalBudgetCRCUnplanned": "unplanned"}},
	}
	for _, tc := range cases {
		data, err := json.Marshal(tc.value)
		if err != nil {
			t.Fatal(err)
		}
		var object map[string]json.RawMessage
		if err := json.Unmarshal(data, &object); err != nil {
			t.Fatal(err)
		}
		typ := reflect.TypeOf(tc.value)
		for field, key := range tc.keys {
			if _, ok := object[key]; !ok {
				t.Errorf("%s.%s missing JSON key %s", typ, field, key)
			}
			if _, ok := object[field]; ok {
				t.Errorf("%s exposes renamed key %s", typ, field)
			}
			if member, ok := typ.FieldByName(field); !ok || member.Tag.Get("json") != key {
				t.Errorf("%s.%s does not decode JSON key %s", typ, field, key)
			}
		}
	}
}
