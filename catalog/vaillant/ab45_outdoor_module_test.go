package vaillant

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
)

func ab45Decode(ab45T *testing.T, ab45Data []byte) map[string]any {
	ab45T.Helper()
	var ab45Value map[string]any
	if ab45Error := json.Unmarshal(ab45Data, &ab45Value); ab45Error != nil {
		ab45T.Fatal(ab45Error)
	}
	return ab45Value
}

func ab45ExpectedField(ab45Name string, ab45Position any, ab45Width int, ab45Domain []int, ab45Signedness, ab45Order string, ab45Numerator, ab45Denominator int, ab45Unit string, ab45Span any, ab45Mappings map[string]string) map[string]any {
	return map[string]any{"name": ab45Name, "logical_position": ab45Position, "width_bits": ab45Width, "raw_domain": ab45Domain, "signedness": ab45Signedness, "byte_order": ab45Order, "conversion": map[string]int{"numerator": ab45Numerator, "denominator": ab45Denominator, "offset": 0}, "unit": ab45Unit, "bit_span": ab45Span, "default": 0, "mappings": ab45Mappings, "grade": "Hypothesis", "raw_fallback": "preserve raw bytes and unresolved alternatives"}
}

func ab45ExpectedBit(ab45Name string, ab45Position, ab45Start, ab45End int, ab45Reserved bool) map[string]any {
	ab45Mappings := map[string]string{"0": "disabled", "1": "enabled"}
	if ab45Reserved {
		ab45Mappings = map[string]string{}
	}
	return ab45ExpectedField(ab45Name, ab45Position, ab45End-ab45Start+1, []int{0, (1 << (ab45End - ab45Start + 1)) - 1}, "unsigned", "not_applicable", 1, 1, "none", []int{ab45Start, ab45End}, ab45Mappings)
}

func ab45ExpectedFields() [][]map[string]any {
	return [][]map[string]any{
		{
			ab45ExpectedBit("raw_bit_0", 3, 0, 0, false), ab45ExpectedBit("raw_bit_1", 3, 1, 1, false),
			ab45ExpectedBit("raw_bit_2", 3, 2, 2, false), ab45ExpectedBit("raw_bit_3", 3, 3, 3, false),
			ab45ExpectedBit("raw_bit_4", 3, 4, 4, false), ab45ExpectedBit("raw_bit_5", 3, 5, 5, false),
			ab45ExpectedBit("reserved", 3, 6, 7, true),
			ab45ExpectedBit("appliance_code_bit_0", 4, 0, 0, false), ab45ExpectedBit("appliance_code_bit_1", 4, 1, 1, false),
			ab45ExpectedBit("appliance_code_bit_2", 4, 2, 2, false), ab45ExpectedBit("appliance_code_bit_3", 4, 3, 3, false),
			ab45ExpectedBit("appliance_code_bit_4", 4, 4, 4, false), ab45ExpectedBit("reserved", 4, 5, 7, true),
			ab45ExpectedField("noise_reduction", 5, 8, []int{0, 255}, "unsigned", "not_applicable", 1, 1, "percent", nil, map[string]string{}),
			ab45ExpectedField("frost_protection_temperature", 6, 8, []int{0, 255}, "unsigned", "not_applicable", 1, 1, "degrees_C", nil, map[string]string{}),
		},
		{
			ab45ExpectedBit("ready", 3, 0, 0, false), ab45ExpectedBit("deicing_request", 3, 1, 1, false),
			ab45ExpectedBit("deicing_active", 3, 2, 2, false), ab45ExpectedBit("deicing_timeout_error", 3, 3, 3, false),
			ab45ExpectedBit("deicing_with_heater", 3, 4, 4, false), ab45ExpectedBit("reserved", 3, 5, 6, true),
			ab45ExpectedBit("cleaning_request", 3, 7, 7, false),
			ab45ExpectedBit("heater", 4, 0, 0, false), ab45ExpectedBit("relay_output", 4, 1, 1, false),
			ab45ExpectedBit("fan_running", 4, 2, 2, false), ab45ExpectedBit("brine_flow_reverse", 4, 3, 3, false),
			ab45ExpectedBit("fan_error", 4, 4, 4, false), ab45ExpectedBit("safety_limiter_error", 4, 5, 5, false),
			ab45ExpectedBit("hot_summer", 4, 6, 6, false), ab45ExpectedBit("bad_weather", 4, 7, 7, false),
			ab45ExpectedField("brine_pump_speed", nil, 8, []int{0, 255}, "unsigned", "not_applicable", 1, 1, "percent", nil, map[string]string{}),
			ab45ExpectedField("fan_PWM", nil, 8, []int{0, 255}, "unsigned", "not_applicable", 1, 1, "percent", nil, map[string]string{}),
			ab45ExpectedField("fan_speed", nil, 8, []int{0, 255}, "unsigned", "not_applicable", 10, 1, "rpm", nil, map[string]string{}),
			ab45ExpectedField("temperature_1", nil, 16, []int{-32768, 32767}, "signed", "little_endian", 1, 16, "degrees_C", nil, map[string]string{}),
			ab45ExpectedField("temperature_2", nil, 16, []int{-32768, 32767}, "signed", "little_endian", 1, 16, "degrees_C", nil, map[string]string{}),
		}, {},
	}
}

func ab45Normalize(ab45T *testing.T, ab45Value any) any {
	ab45T.Helper()
	ab45Data, ab45Error := json.Marshal(ab45Value)
	if ab45Error != nil {
		ab45T.Fatal(ab45Error)
	}
	var ab45Result any
	if ab45Error = json.Unmarshal(ab45Data, &ab45Result); ab45Error != nil {
		ab45T.Fatal(ab45Error)
	}
	return ab45Result
}

func TestAb45CatalogContract(ab45T *testing.T) {
	ab45Catalog := ab45Decode(ab45T, Ab45OutdoorModuleCatalog())
	ab45Entries := ab45Catalog["entries"].([]any)
	if len(ab45Entries) != 3 {
		ab45T.Fatalf("ab45 alternative count: %d", len(ab45Entries))
	}
	ab45Expected := ab45ExpectedFields()
	for ab45Index, ab45Alternative := range []string{"request", "positive", "negative"} {
		ab45T.Run("ab45_"+ab45Alternative, func(ab45T *testing.T) {
			ab45Direction := "response"
			if ab45Index == 0 {
				ab45Direction = "request"
			}
			var ab45Extent any
			if ab45Index == 0 {
				ab45Extent = 7
			}
			if ab45Index == 1 {
				ab45Extent = 12
			}
			ab45Want := map[string]any{"id": "ab45_b521_" + ab45Alternative, "family": "outdoor_module", "variant": nil, "direction": ab45Direction, "service": "B521", "operation": nil, "alternative": ab45Alternative, "logical_extent": ab45Extent, "physical_extent": nil, "coordinate_system": "logical positions include prefix; not physical offsets", "grade": "Hypothesis", "physical_qualification": "Unknown", "raw_fallback": "preserve raw bytes and unresolved alternatives", "fields": ab45Expected[ab45Index]}
			ab45Got := ab45Entries[ab45Index].(map[string]any)
			ab45Fields := ab45Got["fields"].([]any)
			if len(ab45Fields) != len(ab45Expected[ab45Index]) {
				ab45T.Fatal("ab45 field count changed")
			}
			for ab45FieldIndex, ab45ExpectedFieldValue := range ab45Expected[ab45Index] {
				ab45FieldValue := ab45Fields[ab45FieldIndex].(map[string]any)
				if len(ab45FieldValue) != len(ab45ExpectedFieldValue) {
					ab45T.Error("ab45 field property count changed")
				}
				for ab45Property, ab45ExpectedValue := range ab45ExpectedFieldValue {
					ab45T.Run("ab45_"+ab45ExpectedFieldValue["name"].(string)+"_"+ab45Property, func(ab45T *testing.T) {
						ab45ActualValue, ab45Present := ab45FieldValue[ab45Property]
						if !ab45Present || !reflect.DeepEqual(ab45ActualValue, ab45Normalize(ab45T, ab45ExpectedValue)) {
							ab45T.Fatalf("ab45 field %d %s: got %#v want %#v", ab45FieldIndex, ab45Property, ab45ActualValue, ab45ExpectedValue)
						}
					})
				}
			}
			delete(ab45Got, "fields")
			delete(ab45Want, "fields")
			if !reflect.DeepEqual(ab45Got, ab45Normalize(ab45T, ab45Want)) {
				ab45T.Errorf("ab45 %s row contract mismatch: %#v", ab45Alternative, ab45Got)
			}
		})
	}
	delete(ab45Catalog, "entries")
	ab45Want := map[string]any{"representation_grade": "Proven", "interpretation_grade": "Hypothesis", "physical_qualification": "Unknown", "opt_in": true, "runtime_registration": false, "wire_builder": false, "unknown": []string{"complete request bit names", "positive scalar positions", "variant", "operation", "ventilation role", "negative forms", "acknowledgement", "retries/deadlines", "atomicity", "endpoint effects", "relationship to other registers", "physical lengths"}}
	if !reflect.DeepEqual(ab45Catalog, ab45Normalize(ab45T, ab45Want)) {
		ab45T.Fatal("ab45 evidence boundary mismatch")
	}
}

func TestAb45IdentityAxes(ab45T *testing.T) {
	ab45Catalog := ab45Decode(ab45T, Ab45OutdoorModuleCatalog())
	for ab45Index, ab45Entry := range ab45Catalog["entries"].([]any) {
		for _, ab45Axis := range []string{"id", "family", "variant", "direction", "service", "operation", "alternative"} {
			ab45T.Run(ab45Entry.(map[string]any)["id"].(string)+"_"+ab45Axis, func(ab45T *testing.T) {
				ab45Original := ab45Entry.(map[string]any)
				ab45Data, ab45Error := json.Marshal(ab45Original)
				if ab45Error != nil {
					ab45T.Fatal(ab45Error)
				}
				ab45Control := ab45Decode(ab45T, ab45Data)
				if !reflect.DeepEqual(ab45Control, ab45Original) {
					ab45T.Fatal("ab45 valid identity did not round trip")
				}
				ab45Control[ab45Axis] = "ab45_distinct"
				ab45Data, ab45Error = json.Marshal(ab45Control)
				if ab45Error != nil {
					ab45T.Fatal(ab45Error)
				}
				ab45Changed := ab45Decode(ab45T, ab45Data)
				if reflect.DeepEqual(ab45Changed, ab45Original) || ab45Changed[ab45Axis] != "ab45_distinct" {
					ab45T.Fatalf("ab45 row %d lost %s", ab45Index, ab45Axis)
				}
				delete(ab45Control, ab45Axis)
				if reflect.DeepEqual(ab45Control, ab45Original) {
					ab45T.Fatal("ab45 missing identity axis accepted")
				}
			})
		}
	}
}

func TestAb45IndependentBytes(ab45T *testing.T) {
	ab45First := Ab45OutdoorModuleCatalog()
	ab45Control := bytes.Clone(ab45First)
	for ab45Index := range ab45First {
		ab45First[ab45Index] = 'x'
	}
	if !bytes.Equal(Ab45OutdoorModuleCatalog(), ab45Control) {
		ab45T.Fatal("ab45 accessor aliases embedded bytes")
	}
}

func TestAb45WalkingBitsAndConversions(ab45T *testing.T) {
	ab45Entries := ab45Decode(ab45T, Ab45OutdoorModuleCatalog())["entries"].([]any)
	ab45Fields := ab45Entries[1].(map[string]any)["fields"].([]any)
	for ab45Bit := 0; ab45Bit < 8; ab45Bit++ {
		ab45Raw := 1 << ab45Bit
		ab45Seen := []string{}
		for _, ab45Value := range ab45Fields {
			ab45Field := ab45Value.(map[string]any)
			if ab45Field["logical_position"] != float64(3) {
				continue
			}
			ab45Span := ab45Field["bit_span"].([]any)
			ab45Mask := ((1 << int(ab45Span[1].(float64)-ab45Span[0].(float64)+1)) - 1) << int(ab45Span[0].(float64))
			if ab45Raw&ab45Mask != 0 {
				ab45Seen = append(ab45Seen, ab45Field["name"].(string))
			}
		}
		ab45Want := []string{"ready", "deicing_request", "deicing_active", "deicing_timeout_error", "deicing_with_heater", "reserved", "reserved", "cleaning_request"}[ab45Bit]
		if !reflect.DeepEqual(ab45Seen, []string{ab45Want}) {
			ab45T.Fatalf("ab45 walking bit %d: %v", ab45Bit, ab45Seen)
		}
	}
	for _, ab45Case := range []struct {
		ab45Index int
		ab45Raw   []byte
		ab45Want  float64
	}{
		{15, []byte{255}, 255}, {16, []byte{255}, 255}, {17, []byte{255}, 2550}, {18, []byte{254, 255}, -0.125}, {19, []byte{254, 255}, -0.125},
	} {
		ab45Field := ab45Fields[ab45Case.ab45Index].(map[string]any)
		ab45Carrier := int(ab45Case.ab45Raw[0])
		if ab45Field["width_bits"] == float64(16) {
			if ab45Field["byte_order"] != "little_endian" {
				ab45T.Fatal("ab45 byte order")
			}
			ab45Carrier |= int(ab45Case.ab45Raw[1]) << 8
			if ab45Field["signedness"] == "signed" {
				ab45Carrier = int(int16(ab45Carrier))
			}
		}
		ab45Conversion := ab45Field["conversion"].(map[string]any)
		ab45Got := float64(ab45Carrier)*ab45Conversion["numerator"].(float64)/ab45Conversion["denominator"].(float64) + ab45Conversion["offset"].(float64)
		if ab45Got != ab45Case.ab45Want {
			ab45T.Fatalf("ab45 %v = %v, want %v", ab45Field["name"], ab45Got, ab45Case.ab45Want)
		}
	}
}
