package vaillant

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func ab27Decode(ab27t *testing.T) map[string]any {
	ab27t.Helper()
	var ab27doc map[string]any
	if ab27err := json.Unmarshal(Ab27HolidayCatalog(), &ab27doc); ab27err != nil {
		ab27t.Fatal(ab27err)
	}
	return ab27doc
}

func ab27ExpectedRow(ab27id string, ab27changes map[string]any) map[string]any {
	ab27row := map[string]any{
		"id": "ab27_" + ab27id, "family": "Vaillant", "variant": nil, "direction": nil,
		"service": nil, "operation": nil, "alternative": "field-declaration",
		"logical_positions": nil, "extent_bytes": nil, "raw_domain": nil,
		"signedness": "Unknown", "byte_order": "Unknown", "conversion": "Unknown", "unit": "Unknown",
		"default": nil, "mappings": map[string]any{}, "reserved_regions": nil,
		"grade": "Hypothesis", "raw_fallback": "preserve raw bytes and unlabeled values",
	}
	for ab27key, ab27value := range ab27changes {
		ab27row[ab27key] = ab27value
	}
	return ab27row
}

func ab27ExpectedRows() []any {
	ab27rows := []any{
		ab27ExpectedRow("holiday_number", map[string]any{"extent_bytes": float64(1), "raw_domain": []any{float64(0), float64(255)}, "signedness": "unsigned", "conversion": "identity", "default": float64(0)}),
	}
	for _, ab27end := range []string{"start", "end"} {
		for _, ab27part := range []string{"day", "month", "year"} {
			var ab27domain any
			if ab27part == "day" {
				ab27domain = []any{float64(1), float64(31)}
			}
			if ab27part == "month" {
				ab27domain = []any{float64(1), float64(12)}
			}
			ab27rows = append(ab27rows, ab27ExpectedRow(ab27end+"_"+ab27part, map[string]any{"default": float64(0), "calendar_domain": ab27domain, "default_conflict": ab27part != "year"}))
		}
	}
	ab27rows = append(ab27rows,
		ab27ExpectedRow("op02_setpoint", map[string]any{"operation": "OP02", "extent_bytes": float64(1), "raw_domain": []any{float64(0), float64(255)}, "signedness": "unsigned", "conversion": "identity", "unit": "degrees C", "default": float64(0), "candidate_bounds": []any{float64(5), float64(30)}, "default_conflict": true, "off_frost_deletion": "Unknown"}),
		ab27ExpectedRow("room_action", map[string]any{"default": float64(0), "mappings": map[string]any{"0": "no-action", "1": "switch-on", "2": "thermostat"}}),
		ab27ExpectedRow("holiday_status", map[string]any{"default": float64(0), "mappings": map[string]any{"0": "0", "1": "1"}}),
	)
	for _, ab27digit := range []string{"1", "2", "3", "4"} {
		ab27rows = append(ab27rows, ab27ExpectedRow("keycode_"+ab27digit, map[string]any{"extent_bytes": float64(1), "raw_domain": []any{float64(0), float64(255)}, "signedness": "unsigned", "conversion": "independent numeric byte", "unit": "numeric", "excluded_interpretations": []any{"ASCII", "packed decimal", "u32"}}))
	}
	return append(ab27rows,
		ab27ExpectedRow("circuit_name", map[string]any{"extent_bytes": float64(10), "raw_domain": "ASCII characters", "conversion": "ASCII", "unit": "characters", "text_limitations": []any{"padding", "trimming", "NUL termination", "UTF-8", "truncation"}}),
		ab27ExpectedRow("remote_control_write", map[string]any{"variant": "remote-control", "direction": "request", "alternative": "remote-control-write"}),
		ab27ExpectedRow("b50d_ui_direction_write", map[string]any{"variant": "UI-direction", "direction": "request", "service": "B50D", "alternative": "UI-direction-write"}),
	)
}

func TestAb27ExactDeclarations(ab27t *testing.T) {
	ab27doc := ab27Decode(ab27t)
	ab27rows := ab27doc["entries"].([]any)
	ab27expected := ab27ExpectedRows()
	if len(ab27rows) != 17 {
		ab27t.Fatalf("alternative count: %d", len(ab27rows))
	}
	for ab27index, ab27want := range ab27expected {
		ab27got := ab27rows[ab27index].(map[string]any)
		if len(ab27got) != len(ab27want.(map[string]any)) {
			ab27t.Errorf("row %d property count", ab27index)
		}
		for ab27key, ab27value := range ab27want.(map[string]any) {
			if !reflect.DeepEqual(ab27got[ab27key], ab27value) {
				ab27t.Errorf("row %d %s: got %#v want %#v", ab27index, ab27key, ab27got[ab27key], ab27value)
			}
			if _, ab27ok := ab27got[ab27key]; !ab27ok {
				ab27t.Errorf("row %d missing %s", ab27index, ab27key)
			}
		}
	}
}

func TestAb27Boundaries(ab27t *testing.T) {
	ab27doc := ab27Decode(ab27t)
	ab27expected := map[string]any{
		"catalog": "ab27_holiday", "representation_grade": "Proven",
		"coordinate_convention": "logical positions include prefixes; not physical offsets",
		"qualification":         "Unknown", "runtime_registration": false, "wire_builder": false,
		"limitations": []any{"calendar century", "timezone", "valid-date rules", "start-before-end rules", "atomicity", "persistent assignment", "physical response", "readback", "pairing of B50B writes with B50C reads", "request/positive/negative field layouts Unknown", "reserved regions Unknown", "zero day/month conflicts with calendar domains; no delete, disable or reset meaning"},
	}
	if len(ab27doc) != len(ab27expected)+1 {
		ab27t.Fatal("unexpected catalog properties")
	}
	for ab27key, ab27want := range ab27expected {
		if !reflect.DeepEqual(ab27doc[ab27key], ab27want) {
			ab27t.Errorf("%s: %#v", ab27key, ab27doc[ab27key])
		}
	}
}

func TestAb27Copy(ab27t *testing.T) {
	ab27first := Ab27HolidayCatalog()
	ab27control := string(ab27first)
	for ab27index := range ab27first {
		ab27first[ab27index] = 'x'
	}
	if string(Ab27HolidayCatalog()) != ab27control {
		ab27t.Fatal("returned bytes alias embedded catalog")
	}
}

func TestAb27IdentityAxes(ab27t *testing.T) {
	ab27rows := ab27Decode(ab27t)["entries"].([]any)
	for _, ab27axis := range []string{"family", "variant", "direction", "service", "operation", "alternative"} {
		for ab27index, ab27entry := range ab27rows {
			ab27original := ab27entry.(map[string]any)
			for _, ab27replacement := range []any{"ab27_changed", nil} {
				ab27copy := map[string]any{}
				for ab27key, ab27value := range ab27original {
					ab27copy[ab27key] = ab27value
				}
				ab27copy[ab27axis] = ab27replacement
				ab27bytes, ab27err := json.Marshal(ab27copy)
				if ab27err != nil {
					ab27t.Fatal(ab27err)
				}
				var ab27roundtrip map[string]any
				if ab27err = json.Unmarshal(ab27bytes, &ab27roundtrip); ab27err != nil {
					ab27t.Fatal(ab27err)
				}
				if !reflect.DeepEqual(ab27roundtrip, ab27copy) {
					ab27t.Fatalf("row %d lost %s", ab27index, ab27axis)
				}
				if !reflect.DeepEqual(ab27replacement, ab27original[ab27axis]) && reflect.DeepEqual(ab27roundtrip, ab27original) {
					ab27t.Fatalf("row %d collapsed %s", ab27index, ab27axis)
				}
			}
		}
	}
}

func TestAb27SyntheticLimits(ab27t *testing.T) {
	ab27rows := ab27Decode(ab27t)["entries"].([]any)
	for _, ab27index := range []int{1, 2, 4, 5} {
		ab27row := ab27rows[ab27index].(map[string]any)
		ab27domain := ab27row["calendar_domain"].([]any)
		if ab27row["default"] != float64(0) || ab27row["default_conflict"] != true || ab27domain[0].(float64) <= 0 {
			ab27t.Fatal("zero calendar conflict lost")
		}
	}
	ab27bounds := ab27rows[7].(map[string]any)["candidate_bounds"].([]any)
	for ab27index, ab27value := range []float64{0, 5, 30, 31} {
		ab27inside := ab27value >= ab27bounds[0].(float64) && ab27value <= ab27bounds[1].(float64)
		if ab27inside != []bool{false, true, true, false}[ab27index] {
			ab27t.Errorf("setpoint %v", ab27value)
		}
	}
	for ab27index, ab27value := range []string{"0", "1", "2", "3"} {
		_, ab27known := ab27rows[8].(map[string]any)["mappings"].(map[string]any)[ab27value]
		if ab27known != (ab27index < 3) {
			ab27t.Errorf("room action %s", ab27value)
		}
	}
	for ab27index, ab27value := range []string{"0", "1", "2"} {
		_, ab27known := ab27rows[9].(map[string]any)["mappings"].(map[string]any)[ab27value]
		if ab27known != (ab27index < 2) {
			ab27t.Errorf("holiday status %s", ab27value)
		}
	}
	for ab27index, ab27value := range []byte{1, 2, 3, 4} {
		ab27row := ab27rows[10+ab27index].(map[string]any)
		ab27domain := ab27row["raw_domain"].([]any)
		if float64(ab27value) < ab27domain[0].(float64) || float64(ab27value) > ab27domain[1].(float64) || ab27row["conversion"] != "independent numeric byte" || ab27row["extent_bytes"] != float64(1) {
			ab27t.Fatalf("key byte %d", ab27index)
		}
	}
	ab27width := int(ab27rows[14].(map[string]any)["extent_bytes"].(float64))
	for ab27index, ab27text := range []string{"ABCDEFGHI", "ABCDEFGHIJ", "ABCDEFGHIJK", "ABCDEFGHIé"} {
		ab27ascii := strings.IndexFunc(ab27text, func(ab27rune rune) bool { return ab27rune > 127 }) < 0
		ab27matches := len(ab27text) == ab27width && ab27ascii
		if ab27matches != (ab27index == 1) {
			ab27t.Errorf("ASCII extent case %d", ab27index)
		}
	}
}
