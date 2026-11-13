package vaillant

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
)

func outdoorModuleDecode(outdoorModuleT *testing.T, outdoorModuleData []byte) map[string]any {
	outdoorModuleT.Helper()
	var outdoorModuleValue map[string]any
	if outdoorModuleError := json.Unmarshal(outdoorModuleData, &outdoorModuleValue); outdoorModuleError != nil {
		outdoorModuleT.Fatal(outdoorModuleError)
	}
	return outdoorModuleValue
}

func outdoorModuleExpectedField(outdoorModuleName string, outdoorModulePosition any, outdoorModuleWidth int, outdoorModuleDomain []int, outdoorModuleSignedness, outdoorModuleOrder string, outdoorModuleNumerator, outdoorModuleDenominator int, outdoorModuleUnit string, outdoorModuleSpan any, outdoorModuleMappings map[string]string) map[string]any {
	return map[string]any{"name": outdoorModuleName, "logical_position": outdoorModulePosition, "width_bits": outdoorModuleWidth, "raw_domain": outdoorModuleDomain, "signedness": outdoorModuleSignedness, "byte_order": outdoorModuleOrder, "conversion": map[string]int{"numerator": outdoorModuleNumerator, "denominator": outdoorModuleDenominator, "offset": 0}, "unit": outdoorModuleUnit, "bit_span": outdoorModuleSpan, "default": 0, "mappings": outdoorModuleMappings, "grade": "Hypothesis", "raw_fallback": "preserve raw bytes and unresolved alternatives"}
}

func outdoorModuleExpectedBit(outdoorModuleName string, outdoorModulePosition, outdoorModuleStart, outdoorModuleEnd int, outdoorModuleReserved bool) map[string]any {
	outdoorModuleMappings := map[string]string{"0": "disabled", "1": "enabled"}
	if outdoorModuleReserved {
		outdoorModuleMappings = map[string]string{}
	}
	return outdoorModuleExpectedField(outdoorModuleName, outdoorModulePosition, outdoorModuleEnd-outdoorModuleStart+1, []int{0, (1 << (outdoorModuleEnd - outdoorModuleStart + 1)) - 1}, "unsigned", "not_applicable", 1, 1, "none", []int{outdoorModuleStart, outdoorModuleEnd}, outdoorModuleMappings)
}

func outdoorModuleExpectedFields() [][]map[string]any {
	return [][]map[string]any{
		{
			outdoorModuleExpectedBit("raw_bit_0", 3, 0, 0, false), outdoorModuleExpectedBit("raw_bit_1", 3, 1, 1, false),
			outdoorModuleExpectedBit("raw_bit_2", 3, 2, 2, false), outdoorModuleExpectedBit("raw_bit_3", 3, 3, 3, false),
			outdoorModuleExpectedBit("raw_bit_4", 3, 4, 4, false), outdoorModuleExpectedBit("raw_bit_5", 3, 5, 5, false),
			outdoorModuleExpectedBit("reserved", 3, 6, 7, true),
			outdoorModuleExpectedBit("appliance_code_bit_0", 4, 0, 0, false), outdoorModuleExpectedBit("appliance_code_bit_1", 4, 1, 1, false),
			outdoorModuleExpectedBit("appliance_code_bit_2", 4, 2, 2, false), outdoorModuleExpectedBit("appliance_code_bit_3", 4, 3, 3, false),
			outdoorModuleExpectedBit("appliance_code_bit_4", 4, 4, 4, false), outdoorModuleExpectedBit("reserved", 4, 5, 7, true),
			outdoorModuleExpectedField("noise_reduction", 5, 8, []int{0, 255}, "unsigned", "not_applicable", 1, 1, "percent", nil, map[string]string{}),
			outdoorModuleExpectedField("frost_protection_temperature", 6, 8, []int{0, 255}, "unsigned", "not_applicable", 1, 1, "degrees_C", nil, map[string]string{}),
		},
		{
			outdoorModuleExpectedBit("ready", 3, 0, 0, false), outdoorModuleExpectedBit("deicing_request", 3, 1, 1, false),
			outdoorModuleExpectedBit("deicing_active", 3, 2, 2, false), outdoorModuleExpectedBit("deicing_timeout_error", 3, 3, 3, false),
			outdoorModuleExpectedBit("deicing_with_heater", 3, 4, 4, false), outdoorModuleExpectedBit("reserved", 3, 5, 6, true),
			outdoorModuleExpectedBit("cleaning_request", 3, 7, 7, false),
			outdoorModuleExpectedBit("heater", 4, 0, 0, false), outdoorModuleExpectedBit("relay_output", 4, 1, 1, false),
			outdoorModuleExpectedBit("fan_running", 4, 2, 2, false), outdoorModuleExpectedBit("brine_flow_reverse", 4, 3, 3, false),
			outdoorModuleExpectedBit("fan_error", 4, 4, 4, false), outdoorModuleExpectedBit("safety_limiter_error", 4, 5, 5, false),
			outdoorModuleExpectedBit("hot_summer", 4, 6, 6, false), outdoorModuleExpectedBit("bad_weather", 4, 7, 7, false),
			outdoorModuleExpectedField("brine_pump_speed", nil, 8, []int{0, 255}, "unsigned", "not_applicable", 1, 1, "percent", nil, map[string]string{}),
			outdoorModuleExpectedField("fan_PWM", nil, 8, []int{0, 255}, "unsigned", "not_applicable", 1, 1, "percent", nil, map[string]string{}),
			outdoorModuleExpectedField("fan_speed", nil, 8, []int{0, 255}, "unsigned", "not_applicable", 10, 1, "rpm", nil, map[string]string{}),
			outdoorModuleExpectedField("temperature_1", nil, 16, []int{-32768, 32767}, "signed", "little_endian", 1, 16, "degrees_C", nil, map[string]string{}),
			outdoorModuleExpectedField("temperature_2", nil, 16, []int{-32768, 32767}, "signed", "little_endian", 1, 16, "degrees_C", nil, map[string]string{}),
		}, {},
	}
}

func outdoorModuleNormalize(outdoorModuleT *testing.T, outdoorModuleValue any) any {
	outdoorModuleT.Helper()
	outdoorModuleData, outdoorModuleError := json.Marshal(outdoorModuleValue)
	if outdoorModuleError != nil {
		outdoorModuleT.Fatal(outdoorModuleError)
	}
	var outdoorModuleResult any
	if outdoorModuleError = json.Unmarshal(outdoorModuleData, &outdoorModuleResult); outdoorModuleError != nil {
		outdoorModuleT.Fatal(outdoorModuleError)
	}
	return outdoorModuleResult
}

func TestAb45CatalogContract(outdoorModuleT *testing.T) {
	outdoorModuleCatalog := outdoorModuleDecode(outdoorModuleT, OutdoorModuleCatalog())
	outdoorModuleEntries := outdoorModuleCatalog["entries"].([]any)
	if len(outdoorModuleEntries) != 3 {
		outdoorModuleT.Fatalf("ab45 alternative count: %d", len(outdoorModuleEntries))
	}
	outdoorModuleExpected := outdoorModuleExpectedFields()
	for outdoorModuleIndex, outdoorModuleAlternative := range []string{"request", "positive", "negative"} {
		outdoorModuleT.Run("ab45_"+outdoorModuleAlternative, func(outdoorModuleT *testing.T) {
			outdoorModuleDirection := "response"
			if outdoorModuleIndex == 0 {
				outdoorModuleDirection = "request"
			}
			var outdoorModuleExtent any
			if outdoorModuleIndex == 0 {
				outdoorModuleExtent = 7
			}
			if outdoorModuleIndex == 1 {
				outdoorModuleExtent = 12
			}
			outdoorModuleWant := map[string]any{"id": "ab45_b521_" + outdoorModuleAlternative, "family": "outdoor_module", "variant": nil, "direction": outdoorModuleDirection, "service": "B521", "operation": nil, "alternative": outdoorModuleAlternative, "logical_extent": outdoorModuleExtent, "physical_extent": nil, "coordinate_system": "logical positions include prefix; not physical offsets", "grade": "Hypothesis", "physical_qualification": "Unknown", "raw_fallback": "preserve raw bytes and unresolved alternatives", "fields": outdoorModuleExpected[outdoorModuleIndex]}
			outdoorModuleGot := outdoorModuleEntries[outdoorModuleIndex].(map[string]any)
			outdoorModuleFields := outdoorModuleGot["fields"].([]any)
			if len(outdoorModuleFields) != len(outdoorModuleExpected[outdoorModuleIndex]) {
				outdoorModuleT.Fatal("ab45 field count changed")
			}
			for outdoorModuleFieldIndex, outdoorModuleExpectedFieldValue := range outdoorModuleExpected[outdoorModuleIndex] {
				outdoorModuleFieldValue := outdoorModuleFields[outdoorModuleFieldIndex].(map[string]any)
				if len(outdoorModuleFieldValue) != len(outdoorModuleExpectedFieldValue) {
					outdoorModuleT.Error("ab45 field property count changed")
				}
				for outdoorModuleProperty, outdoorModuleExpectedValue := range outdoorModuleExpectedFieldValue {
					outdoorModuleT.Run("ab45_"+outdoorModuleExpectedFieldValue["name"].(string)+"_"+outdoorModuleProperty, func(outdoorModuleT *testing.T) {
						outdoorModuleActualValue, outdoorModulePresent := outdoorModuleFieldValue[outdoorModuleProperty]
						if !outdoorModulePresent || !reflect.DeepEqual(outdoorModuleActualValue, outdoorModuleNormalize(outdoorModuleT, outdoorModuleExpectedValue)) {
							outdoorModuleT.Fatalf("ab45 field %d %s: got %#v want %#v", outdoorModuleFieldIndex, outdoorModuleProperty, outdoorModuleActualValue, outdoorModuleExpectedValue)
						}
					})
				}
			}
			delete(outdoorModuleGot, "fields")
			delete(outdoorModuleWant, "fields")
			if !reflect.DeepEqual(outdoorModuleGot, outdoorModuleNormalize(outdoorModuleT, outdoorModuleWant)) {
				outdoorModuleT.Errorf("ab45 %s row contract mismatch: %#v", outdoorModuleAlternative, outdoorModuleGot)
			}
		})
	}
	delete(outdoorModuleCatalog, "entries")
	outdoorModuleWant := map[string]any{"representation_grade": "Proven", "interpretation_grade": "Hypothesis", "physical_qualification": "Unknown", "opt_in": true, "runtime_registration": false, "wire_builder": false, "unknown": []string{"complete request bit names", "positive scalar positions", "variant", "operation", "ventilation role", "negative forms", "acknowledgement", "retries/deadlines", "atomicity", "endpoint effects", "relationship to other registers", "physical lengths"}}
	if !reflect.DeepEqual(outdoorModuleCatalog, outdoorModuleNormalize(outdoorModuleT, outdoorModuleWant)) {
		outdoorModuleT.Fatal("ab45 evidence boundary mismatch")
	}
}

func TestAb45IdentityAxes(outdoorModuleT *testing.T) {
	outdoorModuleCatalog := outdoorModuleDecode(outdoorModuleT, OutdoorModuleCatalog())
	for outdoorModuleIndex, outdoorModuleEntry := range outdoorModuleCatalog["entries"].([]any) {
		for _, outdoorModuleAxis := range []string{"id", "family", "variant", "direction", "service", "operation", "alternative"} {
			outdoorModuleT.Run(outdoorModuleEntry.(map[string]any)["id"].(string)+"_"+outdoorModuleAxis, func(outdoorModuleT *testing.T) {
				outdoorModuleOriginal := outdoorModuleEntry.(map[string]any)
				outdoorModuleData, outdoorModuleError := json.Marshal(outdoorModuleOriginal)
				if outdoorModuleError != nil {
					outdoorModuleT.Fatal(outdoorModuleError)
				}
				outdoorModuleControl := outdoorModuleDecode(outdoorModuleT, outdoorModuleData)
				if !reflect.DeepEqual(outdoorModuleControl, outdoorModuleOriginal) {
					outdoorModuleT.Fatal("ab45 valid identity did not round trip")
				}
				outdoorModuleControl[outdoorModuleAxis] = "ab45_distinct"
				outdoorModuleData, outdoorModuleError = json.Marshal(outdoorModuleControl)
				if outdoorModuleError != nil {
					outdoorModuleT.Fatal(outdoorModuleError)
				}
				outdoorModuleChanged := outdoorModuleDecode(outdoorModuleT, outdoorModuleData)
				if reflect.DeepEqual(outdoorModuleChanged, outdoorModuleOriginal) || outdoorModuleChanged[outdoorModuleAxis] != "ab45_distinct" {
					outdoorModuleT.Fatalf("ab45 row %d lost %s", outdoorModuleIndex, outdoorModuleAxis)
				}
				delete(outdoorModuleControl, outdoorModuleAxis)
				if reflect.DeepEqual(outdoorModuleControl, outdoorModuleOriginal) {
					outdoorModuleT.Fatal("ab45 missing identity axis accepted")
				}
			})
		}
	}
}

func TestAb45IndependentBytes(outdoorModuleT *testing.T) {
	outdoorModuleFirst := OutdoorModuleCatalog()
	outdoorModuleControl := bytes.Clone(outdoorModuleFirst)
	for outdoorModuleIndex := range outdoorModuleFirst {
		outdoorModuleFirst[outdoorModuleIndex] = 'x'
	}
	if !bytes.Equal(OutdoorModuleCatalog(), outdoorModuleControl) {
		outdoorModuleT.Fatal("ab45 accessor aliases embedded bytes")
	}
}

func TestAb45WalkingBitsAndConversions(outdoorModuleT *testing.T) {
	outdoorModuleEntries := outdoorModuleDecode(outdoorModuleT, OutdoorModuleCatalog())["entries"].([]any)
	outdoorModuleFields := outdoorModuleEntries[1].(map[string]any)["fields"].([]any)
	for outdoorModuleBit := 0; outdoorModuleBit < 8; outdoorModuleBit++ {
		outdoorModuleRaw := 1 << outdoorModuleBit
		outdoorModuleSeen := []string{}
		for _, outdoorModuleValue := range outdoorModuleFields {
			outdoorModuleField := outdoorModuleValue.(map[string]any)
			if outdoorModuleField["logical_position"] != float64(3) {
				continue
			}
			outdoorModuleSpan := outdoorModuleField["bit_span"].([]any)
			outdoorModuleMask := ((1 << int(outdoorModuleSpan[1].(float64)-outdoorModuleSpan[0].(float64)+1)) - 1) << int(outdoorModuleSpan[0].(float64))
			if outdoorModuleRaw&outdoorModuleMask != 0 {
				outdoorModuleSeen = append(outdoorModuleSeen, outdoorModuleField["name"].(string))
			}
		}
		outdoorModuleWant := []string{"ready", "deicing_request", "deicing_active", "deicing_timeout_error", "deicing_with_heater", "reserved", "reserved", "cleaning_request"}[outdoorModuleBit]
		if !reflect.DeepEqual(outdoorModuleSeen, []string{outdoorModuleWant}) {
			outdoorModuleT.Fatalf("ab45 walking bit %d: %v", outdoorModuleBit, outdoorModuleSeen)
		}
	}
	for _, outdoorModuleCase := range []struct {
		outdoorModuleIndex int
		outdoorModuleRaw   []byte
		outdoorModuleWant  float64
	}{
		{15, []byte{255}, 255}, {16, []byte{255}, 255}, {17, []byte{255}, 2550}, {18, []byte{254, 255}, -0.125}, {19, []byte{254, 255}, -0.125},
	} {
		outdoorModuleField := outdoorModuleFields[outdoorModuleCase.outdoorModuleIndex].(map[string]any)
		outdoorModuleCarrier := int(outdoorModuleCase.outdoorModuleRaw[0])
		if outdoorModuleField["width_bits"] == float64(16) {
			if outdoorModuleField["byte_order"] != "little_endian" {
				outdoorModuleT.Fatal("ab45 byte order")
			}
			outdoorModuleCarrier |= int(outdoorModuleCase.outdoorModuleRaw[1]) << 8
			if outdoorModuleField["signedness"] == "signed" {
				outdoorModuleCarrier = int(int16(outdoorModuleCarrier))
			}
		}
		outdoorModuleConversion := outdoorModuleField["conversion"].(map[string]any)
		outdoorModuleGot := float64(outdoorModuleCarrier)*outdoorModuleConversion["numerator"].(float64)/outdoorModuleConversion["denominator"].(float64) + outdoorModuleConversion["offset"].(float64)
		if outdoorModuleGot != outdoorModuleCase.outdoorModuleWant {
			outdoorModuleT.Fatalf("ab45 %v = %v, want %v", outdoorModuleField["name"], outdoorModuleGot, outdoorModuleCase.outdoorModuleWant)
		}
	}
}
