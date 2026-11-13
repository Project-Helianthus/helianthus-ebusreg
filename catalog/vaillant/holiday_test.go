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
		{B503HistoryBranchContext{}, map[string]string{"B503HistoryBranchFamily": "Ab05An05Family", "B503HistoryBranchVariant": "Ab05An05Variant", "B503HistoryBranchOperation": "Ab05An05Operation", "B503HistoryBranchDirection": "Ab05An05Direction"}},
		{B503HistoryBranchCell{}, map[string]string{"B503HistoryBranchRawPresent": "Ab05An05RawPresent", "B503HistoryBranchRaw": "Ab05An05Raw", "B503HistoryBranchNormalizedPresent": "Ab05An05NormalizedPresent", "B503HistoryBranchValue": "Ab05An05Value"}},
		{B503HistoryBranchFields{}, map[string]string{"B503HistoryBranchContext": "Ab05An05Context", "B503HistoryBranchMarkerPresent": "Ab05An05MarkerPresent", "B503HistoryBranchMarker": "Ab05An05Marker", "B503HistoryBranchRemainderPresent": "Ab05An05RemainderPresent", "B503HistoryBranchRemainder": "Ab05An05Remainder", "B503HistoryBranchErrorPresent": "Ab05An05ErrorPresent", "B503HistoryBranchError": "Ab05An05Error", "B503HistoryBranchCalendar": "Ab05An05Calendar"}},
		{B503HistoryBranchComponent{}, map[string]string{"B503HistoryBranchCell": "Ab05An05Cell", "B503HistoryBranchState": "Ab05An05State"}},
		{B503HistoryBranchView{}, map[string]string{"B503HistoryBranchFields": "Ab05An05Fields", "B503HistoryBranchState": "Ab05An05State", "B503HistoryBranchComponents": "Ab05An05Components", "B503HistoryBranchComplete": "Ab05An05Complete", "B503HistoryBranchContract": "Ab05An05Contract", "B503HistoryBranchMeaning": "Ab05An05Meaning", "B503HistoryBranchQualification": "Ab05An05Qualification"}},
		{B511OverlapViewSelection{}, map[string]string{"B511OverlapViewFamily": "Ab11An11Family", "B511OverlapViewVariant": "Ab11An11Variant", "B511OverlapViewOperation": "Ab11An11Operation", "B511OverlapViewDirection": "Ab11An11Direction"}},
		{B511OverlapView{}, map[string]string{"B511OverlapViewSelection": "Ab11An11Selection", "B511OverlapViewRaw": "Ab11An11Raw", "B511OverlapViewAvailableExtent": "Ab11An11AvailableExtent", "B511OverlapViewRequiredExtent": "Ab11An11RequiredExtent", "B511OverlapViewCompleteLogicalExtent": "Ab11An11CompleteLogicalExtent", "B511OverlapViewSelected": "Ab11An11Selected", "B511OverlapViewRawFlag": "Ab11An11RawFlag", "B511OverlapViewFanRunning": "Ab11An11FanRunning", "B511OverlapViewRawOverlay": "Ab11An11RawOverlay", "B511OverlapViewRoomSetpointCelsius": "Ab11An11RoomSetpointCelsius", "B511OverlapViewRawPressure": "Ab11An11RawPressure", "B511OverlapViewPressurePercent": "Ab11An11PressurePercent", "B511OverlapViewPressureDomain": "Ab11An11PressureDomain", "B511OverlapViewCandidateGrade": "Ab11An11CandidateGrade", "B511OverlapViewPhysicalQualification": "Ab11An11PhysicalQualification"}},
		{B512DHWViewScalar{}, map[string]string{"B512DHWViewRaw": "Ab13An13Raw", "B512DHWViewValue": "Ab13An13Value", "B512DHWViewUnit": "Ab13An13Unit", "B512DHWViewUnitConflict": "Ab13An13UnitConflict", "B512DHWViewNamedDomain": "Ab13An13NamedDomain"}},
		{B512DHWViewRequest{}, map[string]string{"B512DHWViewTaskActive": "Ab13An13TaskActive", "B512DHWViewHeatDemand": "Ab13An13HeatDemand", "B512DHWViewLoadingActive": "Ab13An13LoadingActive", "B512DHWViewLoadingPumpFault": "Ab13An13LoadingPumpFault", "B512DHWViewSensorConditionBits": "Ab13An13SensorConditionBits", "B512DHWViewModulation": "Ab13An13Modulation", "B512DHWViewFlow": "Ab13An13Flow", "B512DHWViewTemperatures": "Ab13An13Temperatures", "B512DHWViewInletConditionBits": "Ab13An13InletConditionBits"}},
		{B512DHWViewPositive{}, map[string]string{"B512DHWViewLoadingDisabled": "Ab13An13LoadingDisabled", "B512DHWViewInternalTestActive": "Ab13An13InternalTestActive", "B512DHWViewReady": "Ab13An13Ready", "B512DHWViewFrostProtectionActive": "Ab13An13FrostProtectionActive", "B512DHWViewLegionellaDemand": "Ab13An13LegionellaDemand", "B512DHWViewReset": "Ab13An13Reset", "B512DHWViewSetpoint": "Ab13An13Setpoint", "B512DHWViewReturnTemperature": "Ab13An13ReturnTemperature", "B512DHWViewAdjustment": "Ab13An13Adjustment"}},
		{B512DHWViewResult{}, map[string]string{"B512DHWViewRaw": "Ab13An13Raw", "B512DHWViewUnmodelled": "Ab13An13Unmodelled", "B512DHWViewAdmitted": "Ab13An13Admitted", "B512DHWViewComplete": "Ab13An13Complete", "B512DHWViewContractGrade": "Ab13An13ContractGrade", "B512DHWViewFieldGrade": "Ab13An13FieldGrade", "B512DHWViewQualification": "Ab13An13Qualification", "B512DHWViewPresence": "Ab13An13Presence", "B512DHWViewFreshness": "Ab13An13Freshness", "B512DHWViewExecution": "Ab13An13Execution", "B512DHWViewDefaults": "Ab13An13Defaults", "B512DHWViewRequest": "Ab13An13Request", "B512DHWViewPositive": "Ab13An13Positive"}},
		{B516ClockWordViewSelector{}, map[string]string{"B516ClockWordViewFamily": "Ab15An15Family", "B516ClockWordViewVariant": "Ab15An15Variant", "B516ClockWordViewOperation": "Ab15An15Operation", "B516ClockWordViewDirection": "Ab15An15Direction"}},
		{B516ClockWordViewCell{}, map[string]string{"B516ClockWordViewRaw": "Ab15An15Raw", "B516ClockWordViewFailed": "Ab15An15Failed"}},
		{B516ClockWordViewField{}, map[string]string{"B516ClockWordViewState": "Ab15An15State", "B516ClockWordViewValue": "Ab15An15Value"}},
		{B516ClockWordViewCandidate{}, map[string]string{"B516ClockWordViewAdmitted": "Ab15An15Admitted", "B516ClockWordViewContractGrade": "Ab15An15ContractGrade", "B516ClockWordViewFieldGrade": "Ab15An15FieldGrade", "B516ClockWordViewNativeSupport": "Ab15An15NativeSupport", "B516ClockWordViewPhysicalAdaptation": "Ab15An15PhysicalAdaptation", "B516ClockWordViewPresence": "Ab15An15Presence", "B516ClockWordViewFreshness": "Ab15An15Freshness", "B516ClockWordViewExecution": "Ab15An15Execution"}},
		{B516ClockWordViewClockInput{}, map[string]string{"B516ClockWordViewClock": "Ab15An15Clock", "B516ClockWordViewSetpoint": "Ab15An15Setpoint", "B516ClockWordViewMode": "Ab15An15Mode", "B516ClockWordViewSpecialFunction": "Ab15An15SpecialFunction", "B516ClockWordViewMessages": "Ab15An15Messages"}},
		{B516ClockWordViewClockResult{}, map[string]string{"B516ClockWordViewCandidate": "Ab15An15Candidate", "B516ClockWordViewRaw": "Ab15An15Raw", "B516ClockWordViewClock": "Ab15An15Clock", "B516ClockWordViewSetpointCelsius": "Ab15An15SetpointCelsius", "B516ClockWordViewTablesGrade": "Ab15An15TablesGrade"}},
		{B516ClockWordViewWordResult{}, map[string]string{"B516ClockWordViewCandidate": "Ab15An15Candidate", "B516ClockWordViewRaw": "Ab15An15Raw", "B516ClockWordViewWords": "Ab15An15Words", "B516ClockWordViewPositionsGrade": "Ab15An15PositionsGrade"}},
		{B505OperatingRoleKey{}, map[string]string{"B505OperatingRoleFamily": "Ab16Family", "B505OperatingRoleVariant": "Ab16Variant", "B505OperatingRolePrimary": "Ab16Primary", "B505OperatingRoleSecondary": "Ab16Secondary", "B505OperatingRoleOperation": "Ab16Operation", "B505OperatingRoleDirection": "Ab16Direction", "B505OperatingRoleField": "Ab16Field", "B505OperatingRoleAlternative": "Ab16Alternative", "B505OperatingRoleGroup": "Ab16Group", "B505OperatingRoleInstanceScope": "Ab16InstanceScope", "B505OperatingRoleRegister": "Ab16Register"}},
		{B516TwoWordViewResult{}, map[string]string{"B516TwoWordViewDeclaredCount": "Ab17An17DeclaredCount", "B516TwoWordViewAvailableExtent": "Ab17An17AvailableExtent", "B516TwoWordViewRaw": "Ab17An17Raw", "B516TwoWordViewTrailing": "Ab17An17Trailing", "B516TwoWordViewAdmitted": "Ab17An17Admitted", "B516TwoWordViewStatusAvailable": "Ab17An17StatusAvailable", "B516TwoWordViewStatus": "Ab17An17Status", "B516TwoWordViewWord0Available": "Ab17An17Word0Available", "B516TwoWordViewWord0Raw": "Ab17An17Word0Raw", "B516TwoWordViewWord0Bits": "Ab17An17Word0Bits", "B516TwoWordViewWord0Finite": "Ab17An17Word0Finite", "B516TwoWordViewWord0Value": "Ab17An17Word0Value", "B516TwoWordViewWord1Available": "Ab17An17Word1Available", "B516TwoWordViewWord1Raw": "Ab17An17Word1Raw", "B516TwoWordViewCompleteTwoWord": "Ab17An17CompleteTwoWord", "B516TwoWordViewContractGrade": "Ab17An17ContractGrade", "B516TwoWordViewInterpretationGrade": "Ab17An17InterpretationGrade", "B516TwoWordViewStatusMeaning": "Ab17An17StatusMeaning", "B516TwoWordViewWord0NativeUnit": "Ab17An17Word0NativeUnit", "B516TwoWordViewWord0NativeSign": "Ab17An17Word0NativeSign", "B516TwoWordViewWord1Type": "Ab17An17Word1Type", "B516TwoWordViewWord1Scale": "Ab17An17Word1Scale", "B516TwoWordViewWord1Sign": "Ab17An17Word1Sign", "B516TwoWordViewQualification": "Ab17An17Qualification"}},
		{B531LimitNumericViewField{}, map[string]string{"B531LimitNumericViewRaw": "Ab19An19Raw", "B531LimitNumericViewState": "Ab19An19State", "B531LimitNumericViewAvailable": "Ab19An19Available", "B531LimitNumericViewValue": "Ab19An19Value", "B531LimitNumericViewFresh": "Ab19An19Fresh"}},
		{B531LimitNumericViewPair{}, map[string]string{"B531LimitNumericViewMinimum": "Ab19An19Minimum", "B531LimitNumericViewMaximum": "Ab19An19Maximum", "B531LimitNumericViewFresh": "Ab19An19Fresh"}},
		{B531LimitNumericViewResult{}, map[string]string{"B531LimitNumericViewAdmitted": "Ab19An19Admitted", "B531LimitNumericViewAttempt": "Ab19An19Attempt", "B531LimitNumericViewPrior": "Ab19An19Prior", "B531LimitNumericViewLastKnown": "Ab19An19LastKnown", "B531LimitNumericViewContractGrade": "Ab19An19ContractGrade", "B531LimitNumericViewInterpretationGrade": "Ab19An19InterpretationGrade", "B531LimitNumericViewNativeUnit": "Ab19An19NativeUnit", "B531LimitNumericViewQualification": "Ab19An19Qualification", "B531LimitNumericViewPresence": "Ab19An19Presence", "B531LimitNumericViewExecution": "Ab19An19Execution"}},
		{B51BLocalBudgetCRCBudget{}, map[string]string{"B51BLocalBudgetCRCProcedure": "Ab23An23Procedure", "B51BLocalBudgetCRCChunks": "Ab23An23Chunks", "B51BLocalBudgetCRCConsumed": "Ab23An23Consumed", "B51BLocalBudgetCRCUnplanned": "Ab23An23Unplanned"}},
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
