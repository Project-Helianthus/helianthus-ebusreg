package ebus_standard_catalog

import (
	ab49bytes "bytes"
	ab49json "encoding/json"
	ab49os "os"
	ab49reflect "reflect"
	ab49testing "testing"
)

func ab49ReadMetadata(ab49t *ab49testing.T) map[string]any {
	ab49t.Helper()
	var ab49document map[string]any
	if ab49err := ab49json.Unmarshal(Ab49FEMetadata(), &ab49document); ab49err != nil {
		ab49t.Fatal(ab49err)
	}
	return ab49document
}

func TestAb49Alternatives(ab49t *ab49testing.T) {
	ab49document := ab49ReadMetadata(ab49t)
	ab49alternatives, ab49ok := ab49document["alternatives"].([]any)
	if !ab49ok || len(ab49alternatives) != 4 {
		ab49t.Fatalf("want four direction/encoding alternatives, got %#v", ab49document["alternatives"])
	}
	for ab49index, ab49alternative := range ab49alternatives {
		ab49direction := []string{"request", "positive"}[ab49index/2]
		ab49encoding := []string{"explicit", "unspecified"}[ab49index%2]
		var ab49pebus any = "Unknown"
		if ab49direction == "positive" {
			ab49pebus = true
		}
		ab49want := map[string]any{
			"id":                "ab49_fe01_" + ab49direction + "_" + ab49encoding,
			"direction":         ab49direction,
			"constant_encoding": ab49encoding,
			"grade":             "Hypothesis",
			"pebus_relevant":    ab49pebus,
			"constants": []any{
				map[string]any{"logical_position": float64(0), "bit": float64(0), "width_bits": float64(8), "value_hex": "FE"},
				map[string]any{"logical_position": float64(1), "bit": float64(0), "width_bits": float64(8), "value_hex": "01"},
			},
		}
		if !ab49reflect.DeepEqual(ab49alternative, ab49want) {
			ab49t.Errorf("alternative %d: got %#v, want %#v", ab49index, ab49alternative, ab49want)
		}
	}
}

func TestAb49UnknownBodyAndQualification(ab49t *ab49testing.T) {
	ab49document := ab49ReadMetadata(ab49t)
	delete(ab49document, "alternatives")
	ab49want := map[string]any{
		"id":        "ab49_fe01_partial_body",
		"namespace": "ebus_standard",
		"grade":     "Hypothesis",
		"body": map[string]any{
			"grade": "Unknown", "unit": "Unknown", "enum": "Unknown",
			"selector": "Unknown", "sentinel": "Unknown", "text": "Unknown",
			"termination": "Unknown",
		},
		"family_availability":    "Unknown",
		"native_support":         "Unknown",
		"lifecycle":              "Unknown",
		"physical_effects":       "Unknown",
		"physical_qualification": "Unknown",
	}
	if !ab49reflect.DeepEqual(ab49document, ab49want) {
		ab49t.Fatalf("unknown properties must remain unknown and body/physical extent absent: got %#v, want %#v", ab49document, ab49want)
	}
}

func TestAb49EmbeddedCopy(ab49t *ab49testing.T) {
	ab49fixture, ab49err := ab49os.ReadFile("testdata/ab49_fe_metadata.json")
	if ab49err != nil {
		ab49t.Fatal(ab49err)
	}
	ab49first := Ab49FEMetadata()
	if !ab49bytes.Equal(ab49first, ab49fixture) {
		ab49t.Fatal("accessor must return the embedded fixture")
	}
	if len(ab49first) == 0 {
		ab49t.Fatal("empty metadata")
	}
	for ab49index := range ab49first {
		ab49first[ab49index] ^= 0xff
	}
	if !ab49bytes.Equal(Ab49FEMetadata(), ab49fixture) {
		ab49t.Fatal("returned bytes changed a subsequent result")
	}
}
