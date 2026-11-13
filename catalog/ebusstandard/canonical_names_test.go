package ebusstandard

import (
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
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
		{Standard09BodyViewKey{}, map[string]string{"Standard09BodyViewFamily": "family", "Standard09BodyViewProfile": "profile", "Standard09BodyViewService": "service", "Standard09BodyViewDirection": "direction"}},
		{Standard09BodyViewState{}, map[string]string{"Standard09BodyViewKey": "key", "Standard09BodyViewRaw": "raw", "Standard09BodyViewState": "state", "Standard09BodyViewContractGrade": "contract_grade", "Standard09BodyViewMeaningGrade": "meaning_grade", "Standard09BodyViewExecutionGrade": "execution_grade", "Standard09BodyViewPhysicalSupportGrade": "physical_support_grade", "Standard09BodyViewPhysicalOffsetsGrade": "physical_offsets_grade", "Standard09BodyViewPositivePrefixRelevanceGrade": "positive_prefix_relevance_grade"}},
		{Standard09BodyViewBody{}, map[string]string{"Standard09BodyViewState": "state", "Standard09BodyViewData": "data", "Standard09BodyViewUnresolvedTail": "unresolved_tail"}},
		{Standard09BodyViewAddress{}, map[string]string{"Standard09BodyViewState": "state", "Standard09BodyViewLE": "le", "Standard09BodyViewLow": "low", "Standard09BodyViewHigh": "high"}},
		{Standard09BodyViewCount{}, map[string]string{"Standard09BodyViewState": "state", "Standard09BodyViewValue": "value", "Standard09BodyViewCandidateMaximum": "candidate_maximum", "Standard09BodyViewConflict": "conflict"}},
		{Standard05CandidateViewKey{}, map[string]string{"Standard05CandidateViewFamily": "family", "Standard05CandidateViewVariant": "variant", "Standard05CandidateViewPrimary": "primary", "Standard05CandidateViewSecondary": "secondary", "Standard05CandidateViewOperation": "operation", "Standard05CandidateViewDirection": "direction", "Standard05CandidateViewSelector": "selector"}},
		{Standard05CandidateViewField{}, map[string]string{"Standard05CandidateViewLabel": "label", "Standard05CandidateViewRaw": "raw", "Standard05CandidateViewValue": "value", "Standard05CandidateViewName": "name", "Standard05CandidateViewInterpretation": "interpretation", "Standard05CandidateViewDomain": "domain", "Standard05CandidateViewDefault": "default", "Standard05CandidateViewMinimum": "minimum", "Standard05CandidateViewMaximum": "maximum"}},
		{Standard05CandidateView{}, map[string]string{"Standard05CandidateViewKey": "key", "Standard05CandidateViewRaw": "raw", "Standard05CandidateViewRemainder": "remainder", "Standard05CandidateViewState": "state", "Standard05CandidateViewComplete": "complete", "Standard05CandidateViewFields": "fields", "Standard05CandidateViewAPIGrade": "api_grade", "Standard05CandidateViewMeaningGrade": "meaning_grade", "Standard05CandidateViewNativeGrade": "native_grade", "Standard05CandidateViewUnknowns": "unknowns"}},
		{Standard03CandidateViewSelection{}, map[string]string{"Standard03CandidateViewFamily": "family", "Standard03CandidateViewVariant": "variant", "Standard03CandidateViewFunction": "function", "Standard03CandidateViewOperation": "operation", "Standard03CandidateViewDirection": "direction", "Standard03CandidateViewSelector": "selector", "Standard03CandidateViewEncoding": "encoding"}},
		{Standard03CandidateViewComponent{}, map[string]string{"Standard03CandidateViewName": "name", "Standard03CandidateViewState": "state", "Standard03CandidateViewRaw": "raw", "Standard03CandidateViewValue": "value", "Standard03CandidateViewWidth": "width", "Standard03CandidateViewMaximum": "maximum", "Standard03CandidateViewDefault": "default", "Standard03CandidateViewEncoding": "encoding"}},
		{Standard03CandidateView{}, map[string]string{"Standard03CandidateViewSelection": "selection", "Standard03CandidateViewState": "state", "Standard03CandidateViewRaw": "raw", "Standard03CandidateViewRemainder": "remainder", "Standard03CandidateViewComponents": "components", "Standard03CandidateViewLocalCounter": "local_counter", "Standard03CandidateViewLocalHours": "local_hours", "Standard03CandidateViewLocalMinute": "local_minute", "Standard03CandidateViewImplementationGrade": "implementation_grade", "Standard03CandidateViewMeaningGrade": "meaning_grade", "Standard03CandidateViewNativeProperties": "native_properties"}},
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
	if err := json.Unmarshal([]byte(`{"family":"standard-counter"}`), &selection); err != nil {
		t.Fatal(err)
	}
	if selection.Standard03CandidateViewFamily != "standard-counter" {
		t.Fatal("family key was not decoded")
	}
	var coded Standard03CandidateViewSelection
	if err := json.Unmarshal([]byte(`{"Ab76Family":"standard-counter"}`), &coded); err != nil {
		t.Fatal(err)
	}
	if coded.Standard03CandidateViewFamily != "" {
		t.Fatal("coded family key was decoded")
	}
}

func TestProductionStringsCarryNoWorkItemCode(t *testing.T) {
	coded := regexp.MustCompile(`(?i)(?:^|[^a-z0-9])ab[0-9]{2}`)
	paths, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	checked := 0
	for _, path := range paths {
		if strings.HasSuffix(path, "_test.go") {
			continue
		}
		file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		checked++
		ast.Inspect(file, func(node ast.Node) bool {
			literal, ok := node.(*ast.BasicLit)
			if !ok || literal.Kind != token.STRING {
				return true
			}
			value, err := strconv.Unquote(literal.Value)
			if err != nil {
				t.Fatalf("%s: %v", path, err)
			}
			if coded.MatchString(value) {
				t.Errorf("%s: string %q carries a work-item code", path, value)
			}
			return true
		})
	}
	if checked == 0 {
		t.Fatal("no production file checked")
	}
}
