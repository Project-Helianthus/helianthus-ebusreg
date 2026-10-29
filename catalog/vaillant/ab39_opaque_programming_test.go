package vaillant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func ab39Decode(ab39T *testing.T, ab39Raw []byte) map[string]any {
	ab39T.Helper()
	var ab39Value map[string]any
	if ab39Error := json.Unmarshal(ab39Raw, &ab39Value); ab39Error != nil {
		ab39T.Fatal(ab39Error)
	}
	return ab39Value
}

func ab39ExpectedField(ab39T *testing.T, ab39Name string, ab39Position any) map[string]any {
	ab39T.Helper()
	ab39Field := ab39Decode(ab39T, []byte(`{
		"name":"data", "logical_position":null, "width_bytes":null,
		"extent":"through-end-of-PDU", "minimum_bytes":0, "maximum_bytes":null,
		"raw_domain":"byte-sequence", "signedness":"not-applicable",
		"byte_order":"not-applicable", "conversion":"unchanged", "unit":null,
		"default":null, "mappings":{}, "reserved_regions":null,
		"grade":"Hypothesis", "raw_fallback":"preserve-raw-bytes"
	}`))
	ab39Field["name"] = ab39Name
	ab39Field["logical_position"] = ab39Position
	if ab39Name != "data" {
		ab39Field["width_bytes"] = float64(1)
		ab39Field["extent"] = nil
		ab39Field["minimum_bytes"] = nil
		ab39Field["raw_domain"] = map[string]any{"minimum": float64(0), "maximum": float64(255)}
		ab39Field["signedness"] = "unsigned"
		ab39Field["default"] = float64(0)
		ab39Field["raw_fallback"] = "preserve-raw-number"
	}
	return ab39Field
}

func TestAb39CatalogContract(ab39T *testing.T) {
	ab39Catalog := ab39Decode(ab39T, Ab39OpaqueProgramming())
	ab39Rows, ab39OK := ab39Catalog["rows"].([]any)
	if !ab39OK || len(ab39Rows) != 30 {
		ab39T.Fatalf("ab39 expected exactly 30 separate alternatives, got %v", ab39Catalog["rows"])
	}
	delete(ab39Catalog, "rows")
	ab39Expected := ab39Decode(ab39T, []byte(`{
		"catalog_id":"ab39_opaque_programming", "representation_grade":"Proven",
		"candidate_grade":"Hypothesis", "physical_qualification":"Unknown",
		"selection":"opt-in-metadata-only", "limitations":{
			"crc_algorithm":"Unknown", "crc_coverage":"Unknown", "crc_initialization":"Unknown",
			"persistence_reset":"Unknown", "correlation":"Unknown", "transport_crc_relation":"Unknown",
			"mode_dependent_reply_policy":"Unknown", "full_positive_layouts":"Unknown",
			"physical_native_effects":"Unknown", "reserved_regions":"Unknown", "data_maximum":"Unknown",
			"negative_layouts":"Unknown", "prefix_only_physical_response":"Unknown", "variant":"Unknown"
		}
	}`))
	if !reflect.DeepEqual(ab39Catalog, ab39Expected) {
		ab39T.Errorf("ab39 metadata or limitations changed: %#v", ab39Catalog)
	}
	ab39Cases := []struct {
		ab39Operation string
		ab39ID        string
		ab39Request   string
		ab39Positive  string
	}{
		{"00", "00", "none", "crc"},
		{"01", "01", "data", "crc"},
		{"02", "02", "none", "crc"},
		{"03", "03", "data", "data"},
		{"04", "04", "data", "crc"},
		{"05", "05", "data", "crc"},
		{"06", "06", "data", "prefix"},
		{"07", "07", "data", "crc"},
		{"08", "08", "address-data", "crc"},
		{"AA", "aa", "none", "data"},
	}
	for ab39Index, ab39Case := range ab39Cases {
		for ab39Offset, ab39Alternative := range []string{"request", "positive", "negative"} {
			ab39T.Run("ab39_"+ab39Case.ab39ID+"_"+ab39Alternative, func(ab39T *testing.T) {
				ab39Fields := []any{}
				ab39Layout := "Unknown"
				ab39Direction := "response"
				if ab39Alternative == "request" {
					ab39Direction = "request"
					ab39Layout = "declared-fields"
					switch ab39Case.ab39Request {
					case "none":
						ab39Layout = "no-declared-inputs"
					case "data":
						ab39Fields = append(ab39Fields, ab39ExpectedField(ab39T, "data", float64(3)))
					case "address-data":
						ab39Fields = append(ab39Fields, ab39ExpectedField(ab39T, "slave-address", float64(3)), ab39ExpectedField(ab39T, "data", float64(4)))
					}
				}
				if ab39Alternative == "positive" {
					switch ab39Case.ab39Positive {
					case "crc":
						ab39Layout = "partial"
						ab39Fields = append(ab39Fields, ab39ExpectedField(ab39T, "CRC", nil))
					case "data":
						ab39Layout = "declared-fields"
						ab39Fields = append(ab39Fields, ab39ExpectedField(ab39T, "data", float64(3)))
					case "prefix":
						ab39Layout = "prefix-only"
					}
				}
				ab39Want := map[string]any{
					"id":     "ab39_b519_" + ab39Case.ab39ID + "_" + ab39Alternative,
					"family": "Vaillant", "variant": "Unknown", "direction": ab39Direction,
					"service": "B519", "operation": ab39Case.ab39Operation, "alternative": ab39Alternative,
					"logical_positions": "include-prefix;not-physical-offsets", "extent": nil,
					"layout": ab39Layout, "grade": "Hypothesis", "raw_fallback": "preserve-raw-PDU",
					"reserved_regions": nil, "fields": ab39Fields,
				}
				ab39Got, ab39OK := ab39Rows[ab39Index*3+ab39Offset].(map[string]any)
				if !ab39OK {
					ab39T.Fatal("ab39 row is not an object")
				}
				for ab39Key, ab39Value := range ab39Want {
					if ab39Actual, ab39Present := ab39Got[ab39Key]; !ab39Present || !reflect.DeepEqual(ab39Actual, ab39Value) {
						ab39T.Errorf("ab39 %s: got %#v, want %#v", ab39Key, ab39Actual, ab39Value)
					}
				}
				if len(ab39Got) != len(ab39Want) {
					ab39T.Error("ab39 row has extra or missing properties")
				}
			})
		}
	}
}

func TestAb39IndependentCopies(ab39T *testing.T) {
	ab39First := Ab39OpaqueProgramming()
	ab39Expected := bytes.Clone(ab39First)
	ab39Second := Ab39OpaqueProgramming()
	for ab39Index := range ab39First {
		ab39First[ab39Index] ^= 0xff
	}
	if !bytes.Equal(ab39Second, ab39Expected) || !bytes.Equal(Ab39OpaqueProgramming(), ab39Expected) {
		ab39T.Fatal("ab39 accessor shares mutable storage")
	}
}

func TestAb39SyntheticRawDomains(ab39T *testing.T) {
	ab39Rows := ab39Decode(ab39T, Ab39OpaqueProgramming())["rows"].([]any)
	ab39Seen := map[string]int{}
	for _, ab39Row := range ab39Rows {
		for _, ab39Value := range ab39Row.(map[string]any)["fields"].([]any) {
			ab39Field := ab39Value.(map[string]any)
			ab39Name := ab39Field["name"].(string)
			ab39Seen[ab39Name]++
			if ab39Field["conversion"] != "unchanged" || len(ab39Field["mappings"].(map[string]any)) != 0 {
				ab39T.Fatal("ab39 raw evidence acquired a conversion or sentinel mapping")
			}
			if ab39Name == "data" {
				for _, ab39Synthetic := range [][]int{{0, 255, 65}, {}} {
					ab39Raw, ab39Error := json.Marshal(ab39Synthetic)
					if ab39Error != nil {
						ab39T.Fatal(ab39Error)
					}
					var ab39RoundTrip []int
					if ab39Error = json.Unmarshal(ab39Raw, &ab39RoundTrip); ab39Error != nil || !reflect.DeepEqual(ab39Synthetic, ab39RoundTrip) {
						ab39T.Fatalf("ab39 synthetic opaque bytes changed: %s", ab39Raw)
					}
					if ab39Field["minimum_bytes"] != float64(0) || ab39Field["maximum_bytes"] != nil || ab39Field["width_bytes"] != nil || ab39Field["raw_domain"] != "byte-sequence" {
						ab39T.Fatal("ab39 data acquired scalar, text or fixed buffer constraints")
					}
				}
			} else {
				ab39Domain := ab39Field["raw_domain"].(map[string]any)
				for _, ab39Raw := range []float64{0, 255} {
					if ab39Raw < ab39Domain["minimum"].(float64) || ab39Raw > ab39Domain["maximum"].(float64) {
						ab39T.Fatal("ab39 raw boundary excluded")
					}
					ab39Encoded := json.RawMessage(fmt.Sprintf("%.0f", ab39Raw))
					var ab39Number float64
					if ab39Error := json.Unmarshal(ab39Encoded, &ab39Number); ab39Error != nil || ab39Number != ab39Raw {
						ab39T.Fatal("ab39 raw number lost")
					}
				}
			}
		}
	}
	if !reflect.DeepEqual(ab39Seen, map[string]int{"data": 9, "CRC": 7, "slave-address": 1}) {
		ab39T.Fatalf("ab39 synthetic checks did not cover all declared fields: %v", ab39Seen)
	}
}
