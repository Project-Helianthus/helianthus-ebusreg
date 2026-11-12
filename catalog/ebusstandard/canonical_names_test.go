package ebusstandard

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestCanonicalOutdoorNames(t *testing.T) {
	for _, secondary := range []uint8{0, 2} {
		key := ab28TestKey(secondary, "request", "outdoor_temperature")
		row := ab28TestRow(t, key)
		ab28Assert(t, row, "unit", `"degrees Celsius"`)
		ab28Assert(t, row, "width", `2`)
		ab28Assert(t, row, "encoding", `"signed LE16 /256"`)
		for _, old := range []string{"outside_temperature", "temperature"} {
			key.Ab28Field = old
			if raw, found, err := Ab28Lookup(key); err != nil || found || raw != nil {
				t.Fatalf("documentary key resolved: %s %v", old, err)
			}
		}
	}
	raw, err := Ab28Entries()
	if err != nil {
		t.Fatal(err)
	}
	var rows []struct {
		Key Ab28Key `json:"key"`
	}
	if err := json.Unmarshal(raw, &rows); err != nil {
		t.Fatal(err)
	}
	count := 0
	for _, row := range rows {
		if row.Key.Ab28Secondary != nil && (*row.Key.Ab28Secondary == 0 || *row.Key.Ab28Secondary == 2) && row.Key.Ab28Field == "outdoor_temperature" {
			count++
		}
		if row.Key.Ab28Field == "outside_temperature" {
			t.Fatal("old outdoor name emitted")
		}
	}
	if count != 2 {
		t.Fatalf("canonical quantity count: %d", count)
	}
	duration := ab28TestRow(t, ab28TestKey(2, "request", "duration"))
	ab28Assert(t, duration, "default", `0`)
	ab28Assert(t, duration, "zero_permanent_interpretation", `"Unknown"`)
}

func TestCanonicalOutdoorFieldViews(t *testing.T) {
	for _, fixture := range ab15aq15Cases(t) {
		input := ab15aq15Input(t, fixture.Ab15Aq15Input)
		context := input["context"].(map[string]any)
		if context["direction"] != "request" || (context["secondary"] != float64(0) && context["secondary"] != float64(2)) {
			continue
		}
		t.Run(fixture.Ab15Aq15ID, func(t *testing.T) {
			output := ab15aq15Evaluate(t, input)
			results := output["results"].(map[string]any)
			outdoor := results["outdoor_temperature"].(map[string]any)
			if output["state"] != "complete" || outdoor["state"] != "complete" || outdoor["number"] == nil || outdoor["grade"] != "Hypothesis" {
				t.Fatalf("canonical field not decoded: %v", output)
			}
			for _, old := range []string{"temperature", "outside_temperature"} {
				if _, found := results[old]; found {
					t.Fatalf("alternate field emitted: %s", old)
				}
				legacyInput := ab15aq15Input(t, fixture.Ab15Aq15Input)
				fields := legacyInput["fields"].(map[string]any)
				fields[old] = fields["outdoor_temperature"]
				delete(fields, "outdoor_temperature")
				legacyOutput := ab15aq15Evaluate(t, legacyInput)
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
