package ebusstandard

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestCanonicalOutdoorNames(t *testing.T) {
	for _, secondary := range []uint8{0, 2} {
		key := standard07ResidualTestKey(secondary, "request", "outdoor_temperature")
		row := standard07ResidualTestRow(t, key)
		standard07ResidualAssert(t, row, "unit", `"degrees Celsius"`)
		standard07ResidualAssert(t, row, "width", `2`)
		standard07ResidualAssert(t, row, "encoding", `"signed LE16 /256"`)
		for _, old := range []string{"outside_temperature", "temperature"} {
			key.Standard07ResidualField = old
			if raw, found, err := Standard07ResidualLookup(key); err != nil || found || raw != nil {
				t.Fatalf("documentary key resolved: %s %v", old, err)
			}
		}
	}
	raw, err := Standard07ResidualEntries()
	if err != nil {
		t.Fatal(err)
	}
	var rows []struct {
		Key Standard07ResidualKey `json:"key"`
	}
	if err := json.Unmarshal(raw, &rows); err != nil {
		t.Fatal(err)
	}
	count := 0
	for _, row := range rows {
		if row.Key.Standard07ResidualSecondary != nil && (*row.Key.Standard07ResidualSecondary == 0 || *row.Key.Standard07ResidualSecondary == 2) && row.Key.Standard07ResidualField == "outdoor_temperature" {
			count++
		}
		if row.Key.Standard07ResidualField == "outside_temperature" {
			t.Fatal("old outdoor name emitted")
		}
	}
	if count != 2 {
		t.Fatalf("canonical quantity count: %d", count)
	}
	duration := standard07ResidualTestRow(t, standard07ResidualTestKey(2, "request", "duration"))
	standard07ResidualAssert(t, duration, "default", `0`)
	standard07ResidualAssert(t, duration, "zero_permanent_interpretation", `"Unknown"`)
}

func TestCanonicalOutdoorFieldViews(t *testing.T) {
	for _, fixture := range standard07FieldViewCases(t) {
		input := standard07FieldViewInput(t, fixture.Standard07FieldViewInput)
		context := input["context"].(map[string]any)
		if context["direction"] != "request" || (context["secondary"] != float64(0) && context["secondary"] != float64(2)) {
			continue
		}
		t.Run(fixture.Standard07FieldViewID, func(t *testing.T) {
			output := standard07FieldViewEvaluate(t, input)
			results := output["results"].(map[string]any)
			outdoor := results["outdoor_temperature"].(map[string]any)
			if output["state"] != "complete" || outdoor["state"] != "complete" || outdoor["number"] == nil || outdoor["grade"] != "Hypothesis" {
				t.Fatalf("canonical field not decoded: %v", output)
			}
			for _, old := range []string{"temperature", "outside_temperature"} {
				if _, found := results[old]; found {
					t.Fatalf("alternate field emitted: %s", old)
				}
				legacyInput := standard07FieldViewInput(t, fixture.Standard07FieldViewInput)
				fields := legacyInput["fields"].(map[string]any)
				fields[old] = fields["outdoor_temperature"]
				delete(fields, "outdoor_temperature")
				legacyOutput := standard07FieldViewEvaluate(t, legacyInput)
				legacyResults := legacyOutput["results"].(map[string]any)
				missing := legacyResults["outdoor_temperature"].(map[string]any)
				unknown := legacyResults[old].(map[string]any)
				if legacyOutput["state"] != "partial" || missing["absent"] != true || missing["state"] != "incomplete" || unknown["state"] != "Unknown" || len(unknown) != 2 {
					t.Fatalf("alternate field decoded: %v", legacyOutput)
				}
				for key, result := range results {
					if key != "outdoor_temperature" && !reflect.DeepEqual(legacyResults[key], result) {
						t.Fatalf("unaffected field changed: %s", key)
					}
				}
			}
		})
	}
}

func TestRenamedProjectionJSONKeys(t *testing.T) {
	cases := []struct {
		value any
		keys  map[string]string
	}{
		{Standard09BodyViewKey{}, map[string]string{"Standard09BodyViewFamily": "Ab25An25Family", "Standard09BodyViewProfile": "Ab25An25Profile", "Standard09BodyViewService": "Ab25An25Service", "Standard09BodyViewDirection": "Ab25An25Direction"}},
		{Standard09BodyViewState{}, map[string]string{"Standard09BodyViewKey": "Ab25An25Key", "Standard09BodyViewRaw": "Ab25An25Raw", "Standard09BodyViewState": "Ab25An25State", "Standard09BodyViewContractGrade": "Ab25An25ContractGrade", "Standard09BodyViewMeaningGrade": "Ab25An25MeaningGrade", "Standard09BodyViewExecutionGrade": "Ab25An25ExecutionGrade", "Standard09BodyViewPhysicalSupportGrade": "Ab25An25PhysicalSupportGrade", "Standard09BodyViewPhysicalOffsetsGrade": "Ab25An25PhysicalOffsetsGrade", "Standard09BodyViewPositivePrefixRelevanceGrade": "Ab25An25PositivePrefixRelevanceGrade"}},
		{Standard09BodyViewBody{}, map[string]string{"Standard09BodyViewState": "Ab25An25State", "Standard09BodyViewData": "Ab25An25Data", "Standard09BodyViewUnresolvedTail": "Ab25An25UnresolvedTail"}},
		{Standard09BodyViewAddress{}, map[string]string{"Standard09BodyViewState": "Ab25An25State", "Standard09BodyViewLE": "Ab25An25LE", "Standard09BodyViewLow": "Ab25An25Low", "Standard09BodyViewHigh": "Ab25An25High"}},
		{Standard09BodyViewCount{}, map[string]string{"Standard09BodyViewState": "Ab25An25State", "Standard09BodyViewValue": "Ab25An25Value", "Standard09BodyViewCandidateMaximum": "Ab25An25CandidateMaximum", "Standard09BodyViewConflict": "Ab25An25Conflict"}},
		{Standard05CandidateViewKey{}, map[string]string{"Standard05CandidateViewFamily": "Ab78Family", "Standard05CandidateViewVariant": "Ab78Variant", "Standard05CandidateViewPrimary": "Ab78Primary", "Standard05CandidateViewSecondary": "Ab78Secondary", "Standard05CandidateViewOperation": "Ab78Operation", "Standard05CandidateViewDirection": "Ab78Direction", "Standard05CandidateViewSelector": "Ab78Selector"}},
		{Standard05CandidateViewField{}, map[string]string{"Standard05CandidateViewLabel": "Ab78Label", "Standard05CandidateViewRaw": "Ab78Raw", "Standard05CandidateViewValue": "Ab78Value", "Standard05CandidateViewName": "Ab78Name", "Standard05CandidateViewInterpretation": "Ab78Interpretation", "Standard05CandidateViewDomain": "Ab78Domain", "Standard05CandidateViewDefault": "Ab78Default", "Standard05CandidateViewMinimum": "Ab78Minimum", "Standard05CandidateViewMaximum": "Ab78Maximum"}},
		{Standard05CandidateView{}, map[string]string{"Standard05CandidateViewKey": "Ab78Key", "Standard05CandidateViewRaw": "Ab78Raw", "Standard05CandidateViewRemainder": "Ab78Remainder", "Standard05CandidateViewState": "Ab78State", "Standard05CandidateViewComplete": "Ab78Complete", "Standard05CandidateViewFields": "Ab78Fields", "Standard05CandidateViewAPIGrade": "Ab78APIGrade", "Standard05CandidateViewMeaningGrade": "Ab78MeaningGrade", "Standard05CandidateViewNativeGrade": "Ab78NativeGrade", "Standard05CandidateViewUnknowns": "Ab78Unknowns"}},
		{Standard03CandidateViewSelection{}, map[string]string{"Standard03CandidateViewFamily": "Ab76Family", "Standard03CandidateViewVariant": "Ab76Variant", "Standard03CandidateViewFunction": "Ab76Function", "Standard03CandidateViewOperation": "Ab76Operation", "Standard03CandidateViewDirection": "Ab76Direction", "Standard03CandidateViewSelector": "Ab76Selector", "Standard03CandidateViewEncoding": "Ab76Encoding"}},
		{Standard03CandidateViewComponent{}, map[string]string{"Standard03CandidateViewName": "Ab76Name", "Standard03CandidateViewState": "Ab76State", "Standard03CandidateViewRaw": "Ab76Raw", "Standard03CandidateViewValue": "Ab76Value", "Standard03CandidateViewWidth": "Ab76Width", "Standard03CandidateViewMaximum": "Ab76Maximum", "Standard03CandidateViewDefault": "Ab76Default", "Standard03CandidateViewEncoding": "Ab76Encoding"}},
		{Standard03CandidateView{}, map[string]string{"Standard03CandidateViewSelection": "Ab76Selection", "Standard03CandidateViewState": "Ab76State", "Standard03CandidateViewRaw": "Ab76Raw", "Standard03CandidateViewRemainder": "Ab76Remainder", "Standard03CandidateViewComponents": "Ab76Components", "Standard03CandidateViewLocalCounter": "Ab76LocalCounter", "Standard03CandidateViewLocalHours": "Ab76LocalHours", "Standard03CandidateViewLocalMinute": "Ab76LocalMinute", "Standard03CandidateViewImplementationGrade": "Ab76ImplementationGrade", "Standard03CandidateViewMeaningGrade": "Ab76MeaningGrade", "Standard03CandidateViewNativeProperties": "Ab76NativeProperties"}},
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
		}
	}
	var selection Standard03CandidateViewSelection
	if err := json.Unmarshal([]byte(`{"Ab76Family":"standard-counter"}`), &selection); err != nil {
		t.Fatal(err)
	}
	if selection.Standard03CandidateViewFamily != "standard-counter" {
		t.Fatal("original family key was not decoded")
	}
}
