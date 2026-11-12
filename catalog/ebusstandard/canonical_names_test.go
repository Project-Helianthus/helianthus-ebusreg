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
