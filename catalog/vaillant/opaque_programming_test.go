package vaillant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func opaqueProgrammingDecode(opaqueProgrammingT *testing.T, opaqueProgrammingRaw []byte) map[string]any {
	opaqueProgrammingT.Helper()
	var opaqueProgrammingValue map[string]any
	if opaqueProgrammingError := json.Unmarshal(opaqueProgrammingRaw, &opaqueProgrammingValue); opaqueProgrammingError != nil {
		opaqueProgrammingT.Fatal(opaqueProgrammingError)
	}
	return opaqueProgrammingValue
}

func opaqueProgrammingExpectedField(opaqueProgrammingT *testing.T, opaqueProgrammingName string, opaqueProgrammingPosition any) map[string]any {
	opaqueProgrammingT.Helper()
	opaqueProgrammingField := opaqueProgrammingDecode(opaqueProgrammingT, []byte(`{
		"name":"data", "logical_position":null, "width_bytes":null,
		"extent":"through-end-of-PDU", "minimum_bytes":0, "maximum_bytes":null,
		"raw_domain":"byte-sequence", "signedness":"not-applicable",
		"byte_order":"not-applicable", "conversion":"unchanged", "unit":null,
		"default":null, "mappings":{}, "reserved_regions":null,
		"grade":"Hypothesis", "raw_fallback":"preserve-raw-bytes"
	}`))
	opaqueProgrammingField["name"] = opaqueProgrammingName
	opaqueProgrammingField["logical_position"] = opaqueProgrammingPosition
	if opaqueProgrammingName != "data" {
		opaqueProgrammingField["width_bytes"] = float64(1)
		opaqueProgrammingField["extent"] = nil
		opaqueProgrammingField["minimum_bytes"] = nil
		opaqueProgrammingField["raw_domain"] = map[string]any{"minimum": float64(0), "maximum": float64(255)}
		opaqueProgrammingField["signedness"] = "unsigned"
		opaqueProgrammingField["default"] = float64(0)
		opaqueProgrammingField["raw_fallback"] = "preserve-raw-number"
	}
	return opaqueProgrammingField
}

func TestAb39CatalogContract(opaqueProgrammingT *testing.T) {
	opaqueProgrammingCatalog := opaqueProgrammingDecode(opaqueProgrammingT, OpaqueProgramming())
	opaqueProgrammingRows, opaqueProgrammingOK := opaqueProgrammingCatalog["rows"].([]any)
	if !opaqueProgrammingOK || len(opaqueProgrammingRows) != 30 {
		opaqueProgrammingT.Fatalf("ab39 expected exactly 30 separate alternatives, got %v", opaqueProgrammingCatalog["rows"])
	}
	delete(opaqueProgrammingCatalog, "rows")
	opaqueProgrammingExpected := opaqueProgrammingDecode(opaqueProgrammingT, []byte(`{
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
	if !reflect.DeepEqual(opaqueProgrammingCatalog, opaqueProgrammingExpected) {
		opaqueProgrammingT.Errorf("ab39 metadata or limitations changed: %#v", opaqueProgrammingCatalog)
	}
	opaqueProgrammingCases := []struct {
		opaqueProgrammingOperation string
		opaqueProgrammingID        string
		opaqueProgrammingRequest   string
		opaqueProgrammingPositive  string
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
	for opaqueProgrammingIndex, opaqueProgrammingCase := range opaqueProgrammingCases {
		for opaqueProgrammingOffset, opaqueProgrammingAlternative := range []string{"request", "positive", "negative"} {
			opaqueProgrammingT.Run("ab39_"+opaqueProgrammingCase.opaqueProgrammingID+"_"+opaqueProgrammingAlternative, func(opaqueProgrammingT *testing.T) {
				opaqueProgrammingFields := []any{}
				opaqueProgrammingLayout := "Unknown"
				opaqueProgrammingDirection := "response"
				if opaqueProgrammingAlternative == "request" {
					opaqueProgrammingDirection = "request"
					opaqueProgrammingLayout = "declared-fields"
					switch opaqueProgrammingCase.opaqueProgrammingRequest {
					case "none":
						opaqueProgrammingLayout = "no-declared-inputs"
					case "data":
						opaqueProgrammingFields = append(opaqueProgrammingFields, opaqueProgrammingExpectedField(opaqueProgrammingT, "data", float64(3)))
					case "address-data":
						opaqueProgrammingFields = append(opaqueProgrammingFields, opaqueProgrammingExpectedField(opaqueProgrammingT, "responder-address", float64(3)), opaqueProgrammingExpectedField(opaqueProgrammingT, "data", float64(4)))
					}
				}
				if opaqueProgrammingAlternative == "positive" {
					switch opaqueProgrammingCase.opaqueProgrammingPositive {
					case "crc":
						opaqueProgrammingLayout = "partial"
						opaqueProgrammingFields = append(opaqueProgrammingFields, opaqueProgrammingExpectedField(opaqueProgrammingT, "CRC", nil))
					case "data":
						opaqueProgrammingLayout = "declared-fields"
						opaqueProgrammingFields = append(opaqueProgrammingFields, opaqueProgrammingExpectedField(opaqueProgrammingT, "data", float64(3)))
					case "prefix":
						opaqueProgrammingLayout = "prefix-only"
					}
				}
				opaqueProgrammingWant := map[string]any{
					"id":     "ab39_b519_" + opaqueProgrammingCase.opaqueProgrammingID + "_" + opaqueProgrammingAlternative,
					"family": "Vaillant", "variant": "Unknown", "direction": opaqueProgrammingDirection,
					"service": "B519", "operation": opaqueProgrammingCase.opaqueProgrammingOperation, "alternative": opaqueProgrammingAlternative,
					"logical_positions": "include-prefix;not-physical-offsets", "extent": nil,
					"layout": opaqueProgrammingLayout, "grade": "Hypothesis", "raw_fallback": "preserve-raw-PDU",
					"reserved_regions": nil, "fields": opaqueProgrammingFields,
				}
				opaqueProgrammingGot, opaqueProgrammingOK := opaqueProgrammingRows[opaqueProgrammingIndex*3+opaqueProgrammingOffset].(map[string]any)
				if !opaqueProgrammingOK {
					opaqueProgrammingT.Fatal("ab39 row is not an object")
				}
				for opaqueProgrammingKey, opaqueProgrammingValue := range opaqueProgrammingWant {
					if opaqueProgrammingActual, opaqueProgrammingPresent := opaqueProgrammingGot[opaqueProgrammingKey]; !opaqueProgrammingPresent || !reflect.DeepEqual(opaqueProgrammingActual, opaqueProgrammingValue) {
						opaqueProgrammingT.Errorf("ab39 %s: got %#v, want %#v", opaqueProgrammingKey, opaqueProgrammingActual, opaqueProgrammingValue)
					}
				}
				if len(opaqueProgrammingGot) != len(opaqueProgrammingWant) {
					opaqueProgrammingT.Error("ab39 row has extra or missing properties")
				}
			})
		}
	}
}

func TestAb39IndependentCopies(opaqueProgrammingT *testing.T) {
	opaqueProgrammingFirst := OpaqueProgramming()
	opaqueProgrammingExpected := bytes.Clone(opaqueProgrammingFirst)
	opaqueProgrammingSecond := OpaqueProgramming()
	for opaqueProgrammingIndex := range opaqueProgrammingFirst {
		opaqueProgrammingFirst[opaqueProgrammingIndex] ^= 0xff
	}
	if !bytes.Equal(opaqueProgrammingSecond, opaqueProgrammingExpected) || !bytes.Equal(OpaqueProgramming(), opaqueProgrammingExpected) {
		opaqueProgrammingT.Fatal("ab39 accessor shares mutable storage")
	}
}

func TestAb39SyntheticRawDomains(opaqueProgrammingT *testing.T) {
	opaqueProgrammingRows := opaqueProgrammingDecode(opaqueProgrammingT, OpaqueProgramming())["rows"].([]any)
	opaqueProgrammingSeen := map[string]int{}
	for _, opaqueProgrammingRow := range opaqueProgrammingRows {
		for _, opaqueProgrammingValue := range opaqueProgrammingRow.(map[string]any)["fields"].([]any) {
			opaqueProgrammingField := opaqueProgrammingValue.(map[string]any)
			opaqueProgrammingName := opaqueProgrammingField["name"].(string)
			opaqueProgrammingSeen[opaqueProgrammingName]++
			if opaqueProgrammingField["conversion"] != "unchanged" || len(opaqueProgrammingField["mappings"].(map[string]any)) != 0 {
				opaqueProgrammingT.Fatal("ab39 raw evidence acquired a conversion or sentinel mapping")
			}
			if opaqueProgrammingName == "data" {
				for _, opaqueProgrammingSynthetic := range [][]int{{0, 255, 65}, {}} {
					opaqueProgrammingRaw, opaqueProgrammingError := json.Marshal(opaqueProgrammingSynthetic)
					if opaqueProgrammingError != nil {
						opaqueProgrammingT.Fatal(opaqueProgrammingError)
					}
					var opaqueProgrammingRoundTrip []int
					if opaqueProgrammingError = json.Unmarshal(opaqueProgrammingRaw, &opaqueProgrammingRoundTrip); opaqueProgrammingError != nil || !reflect.DeepEqual(opaqueProgrammingSynthetic, opaqueProgrammingRoundTrip) {
						opaqueProgrammingT.Fatalf("ab39 synthetic opaque bytes changed: %s", opaqueProgrammingRaw)
					}
					if opaqueProgrammingField["minimum_bytes"] != float64(0) || opaqueProgrammingField["maximum_bytes"] != nil || opaqueProgrammingField["width_bytes"] != nil || opaqueProgrammingField["raw_domain"] != "byte-sequence" {
						opaqueProgrammingT.Fatal("ab39 data acquired scalar, text or fixed buffer constraints")
					}
				}
			} else {
				opaqueProgrammingDomain := opaqueProgrammingField["raw_domain"].(map[string]any)
				for _, opaqueProgrammingRaw := range []float64{0, 255} {
					if opaqueProgrammingRaw < opaqueProgrammingDomain["minimum"].(float64) || opaqueProgrammingRaw > opaqueProgrammingDomain["maximum"].(float64) {
						opaqueProgrammingT.Fatal("ab39 raw boundary excluded")
					}
					opaqueProgrammingEncoded := json.RawMessage(fmt.Sprintf("%.0f", opaqueProgrammingRaw))
					var opaqueProgrammingNumber float64
					if opaqueProgrammingError := json.Unmarshal(opaqueProgrammingEncoded, &opaqueProgrammingNumber); opaqueProgrammingError != nil || opaqueProgrammingNumber != opaqueProgrammingRaw {
						opaqueProgrammingT.Fatal("ab39 raw number lost")
					}
				}
			}
		}
	}
	if !reflect.DeepEqual(opaqueProgrammingSeen, map[string]int{"data": 9, "CRC": 7, "responder-address": 1}) {
		opaqueProgrammingT.Fatalf("ab39 synthetic checks did not cover all declared fields: %v", opaqueProgrammingSeen)
	}
}
