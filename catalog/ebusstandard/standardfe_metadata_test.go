package ebusstandard

import (
	standardFEMetadataBytes "bytes"
	standardFEMetadataJSON "encoding/json"
	standardFEMetadataOs "os"
	standardFEMetadataReflect "reflect"
	standardFEMetadataTesting "testing"
)

func standardFEMetadataReadMetadata(standardFEMetadataT *standardFEMetadataTesting.T) map[string]any {
	standardFEMetadataT.Helper()
	var standardFEMetadataDocument map[string]any
	if standardFEMetadataErr := standardFEMetadataJSON.Unmarshal(StandardFEMetadata(), &standardFEMetadataDocument); standardFEMetadataErr != nil {
		standardFEMetadataT.Fatal(standardFEMetadataErr)
	}
	return standardFEMetadataDocument
}

func TestAb49Alternatives(standardFEMetadataT *standardFEMetadataTesting.T) {
	standardFEMetadataDocument := standardFEMetadataReadMetadata(standardFEMetadataT)
	standardFEMetadataAlternatives, standardFEMetadataOk := standardFEMetadataDocument["alternatives"].([]any)
	if !standardFEMetadataOk || len(standardFEMetadataAlternatives) != 4 {
		standardFEMetadataT.Fatalf("want four direction/encoding alternatives, got %#v", standardFEMetadataDocument["alternatives"])
	}
	for standardFEMetadataIndex, standardFEMetadataAlternative := range standardFEMetadataAlternatives {
		standardFEMetadataDirection := []string{"request", "positive"}[standardFEMetadataIndex/2]
		standardFEMetadataEncoding := []string{"explicit", "unspecified"}[standardFEMetadataIndex%2]
		var standardFEMetadataPebus any = "Unknown"
		if standardFEMetadataDirection == "positive" {
			standardFEMetadataPebus = true
		}
		standardFEMetadataWant := map[string]any{
			"id":                "ab49_fe01_" + standardFEMetadataDirection + "_" + standardFEMetadataEncoding,
			"direction":         standardFEMetadataDirection,
			"constant_encoding": standardFEMetadataEncoding,
			"grade":             "Hypothesis",
			"pebus_relevant":    standardFEMetadataPebus,
			"constants": []any{
				map[string]any{"logical_position": float64(0), "bit": float64(0), "width_bits": float64(8), "value_hex": "FE"},
				map[string]any{"logical_position": float64(1), "bit": float64(0), "width_bits": float64(8), "value_hex": "01"},
			},
		}
		if !standardFEMetadataReflect.DeepEqual(standardFEMetadataAlternative, standardFEMetadataWant) {
			standardFEMetadataT.Errorf("alternative %d: got %#v, want %#v", standardFEMetadataIndex, standardFEMetadataAlternative, standardFEMetadataWant)
		}
	}
}

func TestAb49UnknownBodyAndQualification(standardFEMetadataT *standardFEMetadataTesting.T) {
	standardFEMetadataDocument := standardFEMetadataReadMetadata(standardFEMetadataT)
	delete(standardFEMetadataDocument, "alternatives")
	standardFEMetadataWant := map[string]any{
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
	if !standardFEMetadataReflect.DeepEqual(standardFEMetadataDocument, standardFEMetadataWant) {
		standardFEMetadataT.Fatalf("unknown properties must remain unknown and body/physical extent absent: got %#v, want %#v", standardFEMetadataDocument, standardFEMetadataWant)
	}
}

func TestAb49EmbeddedCopy(standardFEMetadataT *standardFEMetadataTesting.T) {
	standardFEMetadataFixture, standardFEMetadataErr := standardFEMetadataOs.ReadFile("testdata/ab49_fe_metadata.json")
	if standardFEMetadataErr != nil {
		standardFEMetadataT.Fatal(standardFEMetadataErr)
	}
	standardFEMetadataFirst := StandardFEMetadata()
	if !standardFEMetadataBytes.Equal(standardFEMetadataFirst, standardFEMetadataFixture) {
		standardFEMetadataT.Fatal("accessor must return the embedded fixture")
	}
	if len(standardFEMetadataFirst) == 0 {
		standardFEMetadataT.Fatal("empty metadata")
	}
	for standardFEMetadataIndex := range standardFEMetadataFirst {
		standardFEMetadataFirst[standardFEMetadataIndex] ^= 0xff
	}
	if !standardFEMetadataBytes.Equal(StandardFEMetadata(), standardFEMetadataFixture) {
		standardFEMetadataT.Fatal("returned bytes changed a subsequent result")
	}
}
