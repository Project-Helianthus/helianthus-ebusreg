package vaillant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

const b509SymbolicAnnotationExpectedTuples = `any|0D||read_symbolic_address||true
any|0E||write_symbolic_address||true
boiler|0D|0x1800|flow_temperature|°C|false
boiler|0D|0x4400|central_heating_pump_active||false
boiler|0D|0x0200|water_pressure|bar|false
boiler|0D|0x0500|flame_active||false
boiler|0D|0xBB00|gas_valve_active||false
boiler|0D|0x8300|fan_speed|rpm|false
boiler|0D|0x2E00|modulation_setpoint|%|true
boiler|0D|0xAB00|boiler_state_code||true
boiler|0D|0x5400|diverter_valve_position||true
boiler|0D|0x5800|dhw_demand_active||false
boiler|0D|0x5500|dhw_water_flow_rate|L/min|false
boiler|0D|0x0E04|maximum_heating_flow_temperature_setpoint|°C|true
boiler|0D|0xA500|maximum_heating_flow_temperature_setpoint|°C|true
boiler|0D|0x0F04|maximum_dhw_flow_temperature_setpoint|°C|false
boiler|0D|0x0704|heating_partial_load_power||true
boiler|0D|0x0804|dhw_partial_load_power||true
boiler|0D|0x9E00|storage_loading_pump_level|%|false
boiler|0D|0x3900|flow_temperature_setpoint|°C|true
boiler|0D|0xEA03|dhw_temperature_setpoint|°C|false
boiler|0D|0x7B00|circulation_pump_active||false
boiler|0D|0x3F00|external_pump_active||false
boiler|0D|0xF203|heating_switch_active||false
boiler|0D|0x2400|fan_speed_setpoint|rpm|false
boiler|0D|0xA400|flame_ionization_current||true
boiler|0D|0xFB00|primary_circuit_flow_rate||true
boiler|0D|0x2800|central_heating_operating_time|h|false
boiler|0D|0x2200|dhw_operating_time|h|false
boiler|0D|0x2900|central_heating_start_count|count|true
boiler|0D|0x2300|dhw_start_count|count|true
boiler|0D|0x1400|pump_operating_time|h|false
boiler|0D|0x1B00|fan_operating_time|h|false
boiler|0D|0x1F00|fire_control_deactivation_count|count|true
boiler|0D|0x2000|temperature_limiter_deactivation_count|count|false
boiler|0D|0x4904|central_heating_pump_mode||true
boiler|0D|0x8104|installer_phone_number||false
boiler|0D|0x2004|service_time_remaining|h|false
boiler|0D|0x0300|air_control_ok||false
boiler|0D|0x2F00|flame_sensing_asic_output||false
boiler|0D|0x3C00|flue_gas_valve_active||false
boiler|0D|0x4700|gas_valve_asic_feedback||false
boiler|0D|0x4800|gas_valve_microcontroller_feedback||false
boiler|0D|0xE400|external_gas_valve_active||false
boiler|0D|0xF200|predictive_air_pressure_switch_event_count|count|true
boiler|0D|0xF300|predictive_air_pressure_switch_accumulated_time|s|true
boiler|0D|0x9A00|device_specific_number||true
boiler|0D|0xDA00|display_mode||false
boiler|0D|0xDF00|minimum_operating_fan_speed|rpm|false
boiler|0D|0xE000|maximum_operating_fan_speed|rpm|false
boiler|0D|0xE500|dcf_time_date||true
heat_pump|0D|0xD000|heat_pump_status||false
heat_pump|0D|0xBC00|environmental_energy_yield_total|kWh|true
heat_pump|0D|0xBA00|environmental_power|kW|false
heat_pump|0D|0xBB00|environmental_power_ratio|%|false
heat_pump|0D|0x8D00|refrigerant_superheat|K|false
heat_pump|0D|0x8E00|refrigerant_subcooling|K|false
heat_pump|0D|0x0F00|source_inlet_temperature|°C|false
heat_pump|0D|0x1400|refrigerant_high_pressure|bar|false
heat_pump|0D|0x1500|refrigerant_low_pressure|bar|false
heat_pump|0D|0xE400|compressor_state||false
heat_pump|0D|0xE900|compressor_control_state||false
heat_pump|0D|0xE600|compressor_minimum_on_duration|s|false
heat_pump|0D|0xE700|compressor_minimum_off_duration|s|false
heat_pump|0D|0xE800|compressor_minimum_start_interval|s|false
heat_pump|0D|0xA800|compressor_dhw_start_count|count|false
heat_pump|0D|0xB400|compressor_high_pressure_cutout_count|count|false
heat_pump|0D|0xCE01|defrost_active||false
heat_pump|0D|0x5701|phase_order_status||false
heat_pump|0D|0x8800|phase_status_flags||false
heat_pump|0D|0x3F00|inrush_current_limiter_output_active||true
heat_pump|0D|0xA901|noise_reduction_active||false
heat_pump|0D|0x2401|noise_reduction_factor||true
heat_pump|0D|0xE201|environmental_power|W|true
heat_pump|0D|0xE301|bivalent_mode||false
heat_pump|0D|0xE401|source_temperature_limit_at_maximum_flow|°C|false
heat_pump|0D|0xE501|source_temperature_limit_slope||true
heat_pump|0D|0xE901|refrigerant_low_pressure_maximum||true
heat_pump|0D|0x8000|compressor_energy_integral||true
heat_pump_interface_unit|54|54 02 00 0D 0A|compressor_speed||true
heat_pump_interface_unit|54|54 02 00 28 0A|electronic_expansion_valve_absolute_position|steps|true
heat_pump_interface_unit|54|54 02 00 47 0A|fan_1_speed|rpm|false
heat_pump_interface_unit|54|54 02 00 49 0A|fan_2_speed|rpm|false
heat_pump_interface_unit|54|54 02 00 5B 09|electrical_input_power||true
heat_pump_interface_unit|54|54 02 00 BD 09|water_flow_rate|L/h|false
heat_pump_interface_unit|54|54 02 00 C5 09|building_pump_output_level||true
heat_pump_interface_unit|54|54 02 00 53 0D|operating_status_code||true
heat_pump_interface_unit|54|54 02 00 6A 09|refrigerant_high_pressure|bar|false
heat_pump_interface_unit|54|54 02 00 98 08|compressor_inlet_temperature|°C|false
heat_pump_interface_unit|54|54 02 00 A2 08|compressor_outlet_temperature|°C|false
heat_pump_interface_unit|54|54 02 00 AC 08|electronic_expansion_valve_outlet_temperature|°C|true
heat_pump_interface_unit|54|54 02 00 DE 08|air_inlet_temperature|°C|false
heat_pump_interface_unit|54|54 02 00 A6 0E|refrigerant_condensation_temperature|°C|true
heat_pump_interface_unit|54|54 02 00 A7 0E|refrigerant_evaporation_temperature|°C|false
heat_pump_interface_unit|54|54 02 00 A9 0E|refrigerant_superheat|K|true
heat_pump_interface_unit|54|54 02 00 B8 0B|total_operating_time|h|false
heat_pump_interface_unit|54|54 02 00 B9 0B|heating_operating_time|h|false
heat_pump_interface_unit|54|54 02 00 BC 0B|dhw_operating_time|h|false
heat_pump_interface_unit|54|54 02 00 C2 0B|compressor_operating_time|h|false
heat_pump_interface_unit|54|54 02 00 C3 0B|compressor_start_count|count|false
heat_pump_interface_unit|54|54 02 00 C4 0B|building_pump_operating_time|h|false
heat_pump_interface_unit|54|54 02 00 C5 0B|building_pump_start_count|count|false
heat_pump_interface_unit|54|54 02 00 C8 0B|four_way_valve_operating_time|h|false
heat_pump_interface_unit|54|54 02 00 C9 0B|four_way_valve_switch_count|count|false
heat_pump_interface_unit|54|54 02 00 D7 0B|fan_1_operating_time|h|false
heat_pump_interface_unit|54|54 02 00 D8 0B|fan_1_start_count|count|false
heat_pump_interface_unit|54|54 02 00 D9 0B|fan_2_operating_time|h|false
heat_pump_interface_unit|54|54 02 00 DA 0B|fan_2_start_count|count|false`

func b509SymbolicAnnotationPointer[b509SymbolicAnnotationType any](b509SymbolicAnnotationValue b509SymbolicAnnotationType) *b509SymbolicAnnotationType {
	return &b509SymbolicAnnotationValue
}

func b509SymbolicAnnotationBaseKey() B509SymbolicAnnotationKey {
	return B509SymbolicAnnotationKey{B509SymbolicAnnotationFamily: "Vaillant", B509SymbolicAnnotationVariant: "none", B509SymbolicAnnotationPrimary: 181, B509SymbolicAnnotationSecondary: b509SymbolicAnnotationPointer(uint8(9)), B509SymbolicAnnotationOperation: b509SymbolicAnnotationPointer(uint8(13)), B509SymbolicAnnotationDirection: "none", B509SymbolicAnnotationField: "symbolic_annotation", B509SymbolicAnnotationAlternative: "none", B509SymbolicAnnotationInstanceScope: "none", B509SymbolicAnnotationClass: "boiler", B509SymbolicAnnotationSelector: b509SymbolicAnnotationPointer("0x1800")}
}

func b509SymbolicAnnotationMustJSON(b509SymbolicAnnotationT *testing.T, b509SymbolicAnnotationValue any) []byte {
	b509SymbolicAnnotationT.Helper()
	b509SymbolicAnnotationRaw, b509SymbolicAnnotationError := json.Marshal(b509SymbolicAnnotationValue)
	if b509SymbolicAnnotationError != nil {
		b509SymbolicAnnotationT.Fatal(b509SymbolicAnnotationError)
	}
	return b509SymbolicAnnotationRaw
}

func b509SymbolicAnnotationObject(b509SymbolicAnnotationT *testing.T, b509SymbolicAnnotationRaw []byte) map[string]any {
	b509SymbolicAnnotationT.Helper()
	var b509SymbolicAnnotationValue map[string]any
	if b509SymbolicAnnotationError := json.Unmarshal(b509SymbolicAnnotationRaw, &b509SymbolicAnnotationValue); b509SymbolicAnnotationError != nil {
		b509SymbolicAnnotationT.Fatal(b509SymbolicAnnotationError)
	}
	return b509SymbolicAnnotationValue
}

func TestAb40FiniteCorrespondence(b509SymbolicAnnotationT *testing.T) {
	b509SymbolicAnnotationRaw, b509SymbolicAnnotationError := B509SymbolicAnnotationEntries()
	if b509SymbolicAnnotationError != nil {
		b509SymbolicAnnotationT.Fatal(b509SymbolicAnnotationError)
	}
	var b509SymbolicAnnotationRows []json.RawMessage
	if b509SymbolicAnnotationError = json.Unmarshal(b509SymbolicAnnotationRaw, &b509SymbolicAnnotationRows); b509SymbolicAnnotationError != nil {
		b509SymbolicAnnotationT.Fatal(b509SymbolicAnnotationError)
	}
	if len(b509SymbolicAnnotationRows) != 108 {
		b509SymbolicAnnotationT.Fatalf("ab40 rows: %d", len(b509SymbolicAnnotationRows))
	}
	b509SymbolicAnnotationCounts := map[string]int{}
	for b509SymbolicAnnotationPosition, b509SymbolicAnnotationLine := range strings.Split(b509SymbolicAnnotationExpectedTuples, "\n") {
		b509SymbolicAnnotationParts := strings.Split(b509SymbolicAnnotationLine, "|")
		b509SymbolicAnnotationOperation, b509SymbolicAnnotationError := strconv.ParseUint(b509SymbolicAnnotationParts[1], 16, 8)
		if b509SymbolicAnnotationError != nil {
			b509SymbolicAnnotationT.Fatal(b509SymbolicAnnotationError)
		}
		b509SymbolicAnnotationKey := b509SymbolicAnnotationBaseKey()
		b509SymbolicAnnotationKey.B509SymbolicAnnotationClass = b509SymbolicAnnotationParts[0]
		b509SymbolicAnnotationKey.B509SymbolicAnnotationOperation = b509SymbolicAnnotationPointer(uint8(b509SymbolicAnnotationOperation))
		b509SymbolicAnnotationKey.B509SymbolicAnnotationSelector = nil
		if b509SymbolicAnnotationParts[2] != "" {
			b509SymbolicAnnotationKey.B509SymbolicAnnotationSelector = b509SymbolicAnnotationPointer(b509SymbolicAnnotationParts[2])
		}
		b509SymbolicAnnotationWant := map[string]any{
			"family": "Vaillant", "variant": "none", "primary": float64(181), "secondary": float64(9), "operation": float64(b509SymbolicAnnotationOperation), "direction": "none", "field": "symbolic_annotation", "alternative": "none", "instance_scope": "none", "class": b509SymbolicAnnotationParts[0], "name": b509SymbolicAnnotationParts[3], "unresolved": b509SymbolicAnnotationParts[5] == "true", "representation_grade": "Proven", "grade": "Hypothesis", "physical_qualification": "Unknown", "width_bytes": nil, "signedness": "Unknown", "scale": nil, "native_effect": "Unknown", "raw_fallback": "preserve-raw-bytes", "raw_form": b509SymbolicAnnotationLine,
		}
		if b509SymbolicAnnotationParts[2] != "" {
			b509SymbolicAnnotationWant["selector"] = b509SymbolicAnnotationParts[2]
		}
		if b509SymbolicAnnotationParts[4] != "" {
			b509SymbolicAnnotationWant["unit"] = b509SymbolicAnnotationParts[4]
		}
		if b509SymbolicAnnotationGot := b509SymbolicAnnotationObject(b509SymbolicAnnotationT, b509SymbolicAnnotationRows[b509SymbolicAnnotationPosition]); !reflect.DeepEqual(b509SymbolicAnnotationGot, b509SymbolicAnnotationWant) {
			b509SymbolicAnnotationT.Fatalf("ab40 tuple %d changed: %#v", b509SymbolicAnnotationPosition, b509SymbolicAnnotationGot)
		}
		b509SymbolicAnnotationFound, b509SymbolicAnnotationPresent, b509SymbolicAnnotationError := B509SymbolicAnnotationLookup(b509SymbolicAnnotationKey)
		if b509SymbolicAnnotationError != nil || !b509SymbolicAnnotationPresent || !bytes.Equal(b509SymbolicAnnotationFound, b509SymbolicAnnotationRows[b509SymbolicAnnotationPosition]) {
			b509SymbolicAnnotationT.Fatalf("ab40 lookup tuple %d: %s %t %v", b509SymbolicAnnotationPosition, b509SymbolicAnnotationFound, b509SymbolicAnnotationPresent, b509SymbolicAnnotationError)
		}
		b509SymbolicAnnotationCounts["none|"+b509SymbolicAnnotationParts[1]]++
	}
	if !reflect.DeepEqual(b509SymbolicAnnotationCounts, map[string]int{"none|0D": 78, "none|0E": 1, "none|54": 29}) {
		b509SymbolicAnnotationT.Fatal(b509SymbolicAnnotationCounts)
	}
}

func TestAb40ExactAxes(b509SymbolicAnnotationT *testing.T) {
	b509SymbolicAnnotationCases := []struct {
		b509SymbolicAnnotationName   string
		b509SymbolicAnnotationChange func(*B509SymbolicAnnotationKey)
	}{
		{"ab40_family", func(b509SymbolicAnnotationKey *B509SymbolicAnnotationKey) {
			b509SymbolicAnnotationKey.B509SymbolicAnnotationFamily = "ab40_other"
		}},
		{"ab40_variant", func(b509SymbolicAnnotationKey *B509SymbolicAnnotationKey) {
			b509SymbolicAnnotationKey.B509SymbolicAnnotationVariant = "ab40_other"
		}},
		{"ab40_primary", func(b509SymbolicAnnotationKey *B509SymbolicAnnotationKey) {
			b509SymbolicAnnotationKey.B509SymbolicAnnotationPrimary = 0
		}},
		{"ab40_secondary_absent", func(b509SymbolicAnnotationKey *B509SymbolicAnnotationKey) {
			b509SymbolicAnnotationKey.B509SymbolicAnnotationSecondary = nil
		}},
		{"ab40_secondary_zero", func(b509SymbolicAnnotationKey *B509SymbolicAnnotationKey) {
			b509SymbolicAnnotationKey.B509SymbolicAnnotationSecondary = b509SymbolicAnnotationPointer(uint8(0))
		}},
		{"ab40_operation_absent", func(b509SymbolicAnnotationKey *B509SymbolicAnnotationKey) {
			b509SymbolicAnnotationKey.B509SymbolicAnnotationOperation = nil
		}},
		{"ab40_operation_zero", func(b509SymbolicAnnotationKey *B509SymbolicAnnotationKey) {
			b509SymbolicAnnotationKey.B509SymbolicAnnotationOperation = b509SymbolicAnnotationPointer(uint8(0))
		}},
		{"ab40_operation_other", func(b509SymbolicAnnotationKey *B509SymbolicAnnotationKey) {
			b509SymbolicAnnotationKey.B509SymbolicAnnotationOperation = b509SymbolicAnnotationPointer(uint8(14))
		}},
		{"ab40_direction", func(b509SymbolicAnnotationKey *B509SymbolicAnnotationKey) {
			b509SymbolicAnnotationKey.B509SymbolicAnnotationDirection = "response"
		}},
		{"ab40_field", func(b509SymbolicAnnotationKey *B509SymbolicAnnotationKey) {
			b509SymbolicAnnotationKey.B509SymbolicAnnotationField = "ab40_other"
		}},
		{"ab40_alternative", func(b509SymbolicAnnotationKey *B509SymbolicAnnotationKey) {
			b509SymbolicAnnotationKey.B509SymbolicAnnotationAlternative = "ab40_other"
		}},
		{"ab40_group_zero", func(b509SymbolicAnnotationKey *B509SymbolicAnnotationKey) {
			b509SymbolicAnnotationKey.B509SymbolicAnnotationGroup = b509SymbolicAnnotationPointer(uint8(0))
		}},
		{"ab40_group_other", func(b509SymbolicAnnotationKey *B509SymbolicAnnotationKey) {
			b509SymbolicAnnotationKey.B509SymbolicAnnotationGroup = b509SymbolicAnnotationPointer(uint8(1))
		}},
		{"ab40_instance_scope", func(b509SymbolicAnnotationKey *B509SymbolicAnnotationKey) {
			b509SymbolicAnnotationKey.B509SymbolicAnnotationInstanceScope = "ab40_other"
		}},
		{"ab40_register_zero", func(b509SymbolicAnnotationKey *B509SymbolicAnnotationKey) {
			b509SymbolicAnnotationKey.B509SymbolicAnnotationRegister = b509SymbolicAnnotationPointer(uint16(0))
		}},
		{"ab40_register_other", func(b509SymbolicAnnotationKey *B509SymbolicAnnotationKey) {
			b509SymbolicAnnotationKey.B509SymbolicAnnotationRegister = b509SymbolicAnnotationPointer(uint16(0x1800))
		}},
		{"ab40_class", func(b509SymbolicAnnotationKey *B509SymbolicAnnotationKey) {
			b509SymbolicAnnotationKey.B509SymbolicAnnotationClass = "heat_pump"
		}},
		{"ab40_class_any", func(b509SymbolicAnnotationKey *B509SymbolicAnnotationKey) {
			b509SymbolicAnnotationKey.B509SymbolicAnnotationClass = "any"
		}},
		{"ab40_selector_absent", func(b509SymbolicAnnotationKey *B509SymbolicAnnotationKey) {
			b509SymbolicAnnotationKey.B509SymbolicAnnotationSelector = nil
		}},
		{"ab40_selector_zero", func(b509SymbolicAnnotationKey *B509SymbolicAnnotationKey) {
			b509SymbolicAnnotationKey.B509SymbolicAnnotationSelector = b509SymbolicAnnotationPointer("0x0000")
		}},
	}
	b509SymbolicAnnotationOriginal := append([]byte(nil), b509SymbolicAnnotationData...)
	b509SymbolicAnnotationT.Cleanup(func() { b509SymbolicAnnotationData = b509SymbolicAnnotationOriginal })
	for _, b509SymbolicAnnotationCase := range b509SymbolicAnnotationCases {
		b509SymbolicAnnotationT.Run(b509SymbolicAnnotationCase.b509SymbolicAnnotationName, func(b509SymbolicAnnotationT *testing.T) {
			b509SymbolicAnnotationData = b509SymbolicAnnotationOriginal
			b509SymbolicAnnotationKey := b509SymbolicAnnotationBaseKey()
			if _, b509SymbolicAnnotationPresent, b509SymbolicAnnotationError := B509SymbolicAnnotationLookup(b509SymbolicAnnotationKey); b509SymbolicAnnotationError != nil || !b509SymbolicAnnotationPresent {
				b509SymbolicAnnotationT.Fatal("ab40 missing control", b509SymbolicAnnotationError)
			}
			b509SymbolicAnnotationCase.b509SymbolicAnnotationChange(&b509SymbolicAnnotationKey)
			if b509SymbolicAnnotationRaw, b509SymbolicAnnotationPresent, b509SymbolicAnnotationError := B509SymbolicAnnotationLookup(b509SymbolicAnnotationKey); b509SymbolicAnnotationError != nil || b509SymbolicAnnotationPresent || b509SymbolicAnnotationRaw != nil {
				b509SymbolicAnnotationT.Fatal("ab40 inexact match", b509SymbolicAnnotationError)
			}
			b509SymbolicAnnotationFirst := b509SymbolicAnnotationObject(b509SymbolicAnnotationT, b509SymbolicAnnotationMustJSON(b509SymbolicAnnotationT, b509SymbolicAnnotationBaseKey()))
			b509SymbolicAnnotationSecond := b509SymbolicAnnotationObject(b509SymbolicAnnotationT, b509SymbolicAnnotationMustJSON(b509SymbolicAnnotationT, b509SymbolicAnnotationKey))
			b509SymbolicAnnotationFirst["ab40_marker"] = "ab40_control"
			b509SymbolicAnnotationSecond["ab40_marker"] = "ab40_changed"
			b509SymbolicAnnotationData = b509SymbolicAnnotationMustJSON(b509SymbolicAnnotationT, []any{b509SymbolicAnnotationFirst, b509SymbolicAnnotationSecond})
			for b509SymbolicAnnotationPosition, b509SymbolicAnnotationWant := range []map[string]any{b509SymbolicAnnotationFirst, b509SymbolicAnnotationSecond} {
				b509SymbolicAnnotationLookupKey := b509SymbolicAnnotationBaseKey()
				if b509SymbolicAnnotationPosition == 1 {
					b509SymbolicAnnotationLookupKey = b509SymbolicAnnotationKey
				}
				b509SymbolicAnnotationRaw, b509SymbolicAnnotationPresent, b509SymbolicAnnotationError := B509SymbolicAnnotationLookup(b509SymbolicAnnotationLookupKey)
				if b509SymbolicAnnotationError != nil || !b509SymbolicAnnotationPresent || !reflect.DeepEqual(b509SymbolicAnnotationObject(b509SymbolicAnnotationT, b509SymbolicAnnotationRaw), b509SymbolicAnnotationWant) {
					b509SymbolicAnnotationT.Fatal("ab40 distinct valid controls collapsed", b509SymbolicAnnotationError)
				}
			}
		})
	}
}

func TestAb40SelectorTokens(b509SymbolicAnnotationT *testing.T) {
	b509SymbolicAnnotationKey := b509SymbolicAnnotationBaseKey()
	b509SymbolicAnnotationKey.B509SymbolicAnnotationClass = "any"
	b509SymbolicAnnotationKey.B509SymbolicAnnotationSelector = nil
	if _, b509SymbolicAnnotationPresent, b509SymbolicAnnotationError := B509SymbolicAnnotationLookup(b509SymbolicAnnotationKey); b509SymbolicAnnotationError != nil || !b509SymbolicAnnotationPresent {
		b509SymbolicAnnotationT.Fatal("ab40 selectorless control")
	}
	for _, b509SymbolicAnnotationClass := range []string{"any", "boiler", "heat_pump"} {
		for _, b509SymbolicAnnotationSelector := range []*string{nil, b509SymbolicAnnotationPointer("0x0000"), b509SymbolicAnnotationPointer("0")} {
			b509SymbolicAnnotationKey.B509SymbolicAnnotationClass = b509SymbolicAnnotationClass
			b509SymbolicAnnotationKey.B509SymbolicAnnotationSelector = b509SymbolicAnnotationSelector
			_, b509SymbolicAnnotationPresent, b509SymbolicAnnotationError := B509SymbolicAnnotationLookup(b509SymbolicAnnotationKey)
			if b509SymbolicAnnotationError != nil || b509SymbolicAnnotationPresent != (b509SymbolicAnnotationClass == "any" && b509SymbolicAnnotationSelector == nil) {
				b509SymbolicAnnotationT.Fatal("ab40 wildcard or selectorless/zero collapse")
			}
		}
	}
	b509SymbolicAnnotationKey.B509SymbolicAnnotationClass = "heat_pump_interface_unit"
	b509SymbolicAnnotationKey.B509SymbolicAnnotationOperation = b509SymbolicAnnotationPointer(uint8(0x54))
	b509SymbolicAnnotationKey.B509SymbolicAnnotationSelector = b509SymbolicAnnotationPointer("54 02 00 0D 0A")
	if _, b509SymbolicAnnotationPresent, b509SymbolicAnnotationError := B509SymbolicAnnotationLookup(b509SymbolicAnnotationKey); b509SymbolicAnnotationError != nil || !b509SymbolicAnnotationPresent {
		b509SymbolicAnnotationT.Fatal("ab40 opaque control")
	}
	for b509SymbolicAnnotationPosition := 0; b509SymbolicAnnotationPosition < 5; b509SymbolicAnnotationPosition++ {
		b509SymbolicAnnotationTokens := strings.Fields("54 02 00 0D 0A")
		b509SymbolicAnnotationTokens[b509SymbolicAnnotationPosition] = "FF"
		b509SymbolicAnnotationKey.B509SymbolicAnnotationSelector = b509SymbolicAnnotationPointer(strings.Join(b509SymbolicAnnotationTokens, " "))
		if _, b509SymbolicAnnotationPresent, b509SymbolicAnnotationError := B509SymbolicAnnotationLookup(b509SymbolicAnnotationKey); b509SymbolicAnnotationError != nil || b509SymbolicAnnotationPresent {
			b509SymbolicAnnotationT.Fatalf("ab40 changed byte %d matched", b509SymbolicAnnotationPosition)
		}
	}
	for _, b509SymbolicAnnotationToken := range []string{"0x0D0A", "0D 0A", "02 00 0D 0A", "54 02 00 0D 0A 00"} {
		b509SymbolicAnnotationKey.B509SymbolicAnnotationSelector = b509SymbolicAnnotationPointer(b509SymbolicAnnotationToken)
		if _, b509SymbolicAnnotationPresent, b509SymbolicAnnotationError := B509SymbolicAnnotationLookup(b509SymbolicAnnotationKey); b509SymbolicAnnotationError != nil || b509SymbolicAnnotationPresent {
			b509SymbolicAnnotationT.Fatal("ab40 narrowed token matched")
		}
	}
}

func TestAb40ValidationAndRetention(b509SymbolicAnnotationT *testing.T) {
	b509SymbolicAnnotationOriginal := append([]byte(nil), b509SymbolicAnnotationData...)
	b509SymbolicAnnotationT.Cleanup(func() { b509SymbolicAnnotationData = b509SymbolicAnnotationOriginal })
	b509SymbolicAnnotationRow := b509SymbolicAnnotationObject(b509SymbolicAnnotationT, b509SymbolicAnnotationMustJSON(b509SymbolicAnnotationT, b509SymbolicAnnotationBaseKey()))
	b509SymbolicAnnotationRow["unresolved"] = true
	b509SymbolicAnnotationRow["grade"] = "Hypothesis"
	b509SymbolicAnnotationRow["width_bytes"] = nil
	b509SymbolicAnnotationRow["ab40_unknown"] = map[string]any{"ab40_raw": []any{float64(0), float64(255)}}
	b509SymbolicAnnotationSource := b509SymbolicAnnotationMustJSON(b509SymbolicAnnotationT, []any{b509SymbolicAnnotationRow})
	b509SymbolicAnnotationRows, b509SymbolicAnnotationIndex, b509SymbolicAnnotationError := b509SymbolicAnnotationValidate(b509SymbolicAnnotationSource)
	if b509SymbolicAnnotationError != nil {
		b509SymbolicAnnotationT.Fatal(b509SymbolicAnnotationError)
	}
	for b509SymbolicAnnotationPosition := range b509SymbolicAnnotationSource {
		b509SymbolicAnnotationSource[b509SymbolicAnnotationPosition] = ' '
	}
	if !reflect.DeepEqual(b509SymbolicAnnotationObject(b509SymbolicAnnotationT, b509SymbolicAnnotationRows[0]), b509SymbolicAnnotationRow) || !reflect.DeepEqual(b509SymbolicAnnotationObject(b509SymbolicAnnotationT, b509SymbolicAnnotationIndex[b509SymbolicAnnotationIdentity(b509SymbolicAnnotationBaseKey())]), b509SymbolicAnnotationRow) {
		b509SymbolicAnnotationT.Fatal("ab40 source aliases parsed rows")
	}
	b509SymbolicAnnotationData = b509SymbolicAnnotationMustJSON(b509SymbolicAnnotationT, []any{b509SymbolicAnnotationRow})
	b509SymbolicAnnotationRaw, b509SymbolicAnnotationPresent, b509SymbolicAnnotationError := B509SymbolicAnnotationLookup(b509SymbolicAnnotationBaseKey())
	if b509SymbolicAnnotationError != nil || !b509SymbolicAnnotationPresent || !reflect.DeepEqual(b509SymbolicAnnotationObject(b509SymbolicAnnotationT, b509SymbolicAnnotationRaw), b509SymbolicAnnotationRow) {
		b509SymbolicAnnotationT.Fatal("ab40 absent name/unit or unresolved lost")
	}
	for _, b509SymbolicAnnotationAxis := range []string{"grade", "width_bytes"} {
		b509SymbolicAnnotationChanged := b509SymbolicAnnotationObject(b509SymbolicAnnotationT, b509SymbolicAnnotationMustJSON(b509SymbolicAnnotationT, b509SymbolicAnnotationRow))
		if b509SymbolicAnnotationAxis == "grade" {
			b509SymbolicAnnotationChanged[b509SymbolicAnnotationAxis] = "Unknown"
		} else {
			b509SymbolicAnnotationChanged[b509SymbolicAnnotationAxis] = float64(2)
		}
		b509SymbolicAnnotationData = b509SymbolicAnnotationMustJSON(b509SymbolicAnnotationT, []any{b509SymbolicAnnotationChanged})
		b509SymbolicAnnotationRaw, b509SymbolicAnnotationPresent, b509SymbolicAnnotationError = B509SymbolicAnnotationLookup(b509SymbolicAnnotationBaseKey())
		if b509SymbolicAnnotationError != nil || !b509SymbolicAnnotationPresent || !reflect.DeepEqual(b509SymbolicAnnotationObject(b509SymbolicAnnotationT, b509SymbolicAnnotationRaw), b509SymbolicAnnotationChanged) {
			b509SymbolicAnnotationT.Fatal("ab40 nonidentity metadata not retained")
		}
		b509SymbolicAnnotationData = b509SymbolicAnnotationMustJSON(b509SymbolicAnnotationT, []any{b509SymbolicAnnotationRow, b509SymbolicAnnotationChanged})
		if _, _, b509SymbolicAnnotationError = b509SymbolicAnnotationValidate(b509SymbolicAnnotationData); b509SymbolicAnnotationError == nil {
			b509SymbolicAnnotationT.Fatal("ab40 grade/width cannot distinguish duplicate identities")
		}
	}
	b509SymbolicAnnotationBad := [][]byte{[]byte(`null`), []byte(`{}`), []byte(`[`), []byte(`[null]`), []byte(`[1]`), []byte(`[{}]`), []byte(`[] true`), b509SymbolicAnnotationMustJSON(b509SymbolicAnnotationT, []any{b509SymbolicAnnotationRow, b509SymbolicAnnotationRow})}
	for _, b509SymbolicAnnotationAxis := range []string{"primary", "operation", "secondary", "group", "register", "selector"} {
		b509SymbolicAnnotationChanged := b509SymbolicAnnotationObject(b509SymbolicAnnotationT, b509SymbolicAnnotationMustJSON(b509SymbolicAnnotationT, b509SymbolicAnnotationRow))
		b509SymbolicAnnotationChanged[b509SymbolicAnnotationAxis] = nil
		b509SymbolicAnnotationBad = append(b509SymbolicAnnotationBad, b509SymbolicAnnotationMustJSON(b509SymbolicAnnotationT, []any{b509SymbolicAnnotationChanged}))
		b509SymbolicAnnotationChanged[b509SymbolicAnnotationAxis] = map[string]any{}
		b509SymbolicAnnotationBad = append(b509SymbolicAnnotationBad, b509SymbolicAnnotationMustJSON(b509SymbolicAnnotationT, []any{b509SymbolicAnnotationChanged}))
	}
	for b509SymbolicAnnotationPosition, b509SymbolicAnnotationInput := range b509SymbolicAnnotationBad {
		b509SymbolicAnnotationT.Run(fmt.Sprintf("ab40_invalid_%d", b509SymbolicAnnotationPosition), func(b509SymbolicAnnotationT *testing.T) {
			if _, _, b509SymbolicAnnotationError := b509SymbolicAnnotationValidate(b509SymbolicAnnotationInput); b509SymbolicAnnotationError == nil {
				b509SymbolicAnnotationT.Fatal("ab40 invalid accepted")
			}
			b509SymbolicAnnotationData = b509SymbolicAnnotationInput
			if b509SymbolicAnnotationRaw, b509SymbolicAnnotationError := B509SymbolicAnnotationEntries(); b509SymbolicAnnotationError == nil || b509SymbolicAnnotationRaw != nil {
				b509SymbolicAnnotationT.Fatal("ab40 entries accepted invalid data")
			}
			if b509SymbolicAnnotationRaw, b509SymbolicAnnotationPresent, b509SymbolicAnnotationError := B509SymbolicAnnotationLookup(b509SymbolicAnnotationBaseKey()); b509SymbolicAnnotationError == nil || b509SymbolicAnnotationPresent || b509SymbolicAnnotationRaw != nil {
				b509SymbolicAnnotationT.Fatal("ab40 lookup accepted invalid data")
			}
		})
	}
}

func TestAb40IndependentStorage(b509SymbolicAnnotationT *testing.T) {
	b509SymbolicAnnotationFirst, b509SymbolicAnnotationError := B509SymbolicAnnotationEntries()
	if b509SymbolicAnnotationError != nil {
		b509SymbolicAnnotationT.Fatal(b509SymbolicAnnotationError)
	}
	b509SymbolicAnnotationExpected := bytes.Clone(b509SymbolicAnnotationFirst)
	b509SymbolicAnnotationSecond, b509SymbolicAnnotationError := B509SymbolicAnnotationEntries()
	if b509SymbolicAnnotationError != nil {
		b509SymbolicAnnotationT.Fatal(b509SymbolicAnnotationError)
	}
	for b509SymbolicAnnotationPosition := range b509SymbolicAnnotationFirst {
		b509SymbolicAnnotationFirst[b509SymbolicAnnotationPosition] ^= 0xff
	}
	b509SymbolicAnnotationThird, b509SymbolicAnnotationError := B509SymbolicAnnotationEntries()
	if b509SymbolicAnnotationError != nil || !bytes.Equal(b509SymbolicAnnotationExpected, b509SymbolicAnnotationSecond) || !bytes.Equal(b509SymbolicAnnotationExpected, b509SymbolicAnnotationThird) {
		b509SymbolicAnnotationT.Fatal("ab40 entries share storage")
	}
	b509SymbolicAnnotationFirst, b509SymbolicAnnotationPresent, b509SymbolicAnnotationError := B509SymbolicAnnotationLookup(b509SymbolicAnnotationBaseKey())
	if b509SymbolicAnnotationError != nil || !b509SymbolicAnnotationPresent {
		b509SymbolicAnnotationT.Fatal("ab40 lookup missing")
	}
	b509SymbolicAnnotationExpected = bytes.Clone(b509SymbolicAnnotationFirst)
	b509SymbolicAnnotationSecond, b509SymbolicAnnotationPresent, b509SymbolicAnnotationError = B509SymbolicAnnotationLookup(b509SymbolicAnnotationBaseKey())
	if b509SymbolicAnnotationError != nil || !b509SymbolicAnnotationPresent {
		b509SymbolicAnnotationT.Fatal("ab40 second lookup missing")
	}
	for b509SymbolicAnnotationPosition := range b509SymbolicAnnotationFirst {
		b509SymbolicAnnotationFirst[b509SymbolicAnnotationPosition] ^= 0xff
	}
	b509SymbolicAnnotationThird, b509SymbolicAnnotationPresent, b509SymbolicAnnotationError = B509SymbolicAnnotationLookup(b509SymbolicAnnotationBaseKey())
	if b509SymbolicAnnotationError != nil || !b509SymbolicAnnotationPresent || !bytes.Equal(b509SymbolicAnnotationExpected, b509SymbolicAnnotationSecond) || !bytes.Equal(b509SymbolicAnnotationExpected, b509SymbolicAnnotationThird) {
		b509SymbolicAnnotationT.Fatal("ab40 lookups share storage")
	}
}

func TestAb40OptionalAbsentZero(b509SymbolicAnnotationT *testing.T) {
	b509SymbolicAnnotationOriginal := append([]byte(nil), b509SymbolicAnnotationData...)
	b509SymbolicAnnotationT.Cleanup(func() { b509SymbolicAnnotationData = b509SymbolicAnnotationOriginal })
	for _, b509SymbolicAnnotationAxis := range []string{"secondary", "operation", "group", "register", "selector"} {
		b509SymbolicAnnotationT.Run("ab40_"+b509SymbolicAnnotationAxis, func(b509SymbolicAnnotationT *testing.T) {
			b509SymbolicAnnotationAbsent := b509SymbolicAnnotationBaseKey()
			b509SymbolicAnnotationZero := b509SymbolicAnnotationBaseKey()
			switch b509SymbolicAnnotationAxis {
			case "secondary":
				b509SymbolicAnnotationAbsent.B509SymbolicAnnotationSecondary = nil
				b509SymbolicAnnotationZero.B509SymbolicAnnotationSecondary = b509SymbolicAnnotationPointer(uint8(0))
			case "operation":
				b509SymbolicAnnotationAbsent.B509SymbolicAnnotationOperation = nil
				b509SymbolicAnnotationZero.B509SymbolicAnnotationOperation = b509SymbolicAnnotationPointer(uint8(0))
			case "group":
				b509SymbolicAnnotationAbsent.B509SymbolicAnnotationGroup = nil
				b509SymbolicAnnotationZero.B509SymbolicAnnotationGroup = b509SymbolicAnnotationPointer(uint8(0))
			case "register":
				b509SymbolicAnnotationAbsent.B509SymbolicAnnotationRegister = nil
				b509SymbolicAnnotationZero.B509SymbolicAnnotationRegister = b509SymbolicAnnotationPointer(uint16(0))
			case "selector":
				b509SymbolicAnnotationAbsent.B509SymbolicAnnotationSelector = nil
				b509SymbolicAnnotationZero.B509SymbolicAnnotationSelector = b509SymbolicAnnotationPointer("0x0000")
			}
			b509SymbolicAnnotationAbsentRow := b509SymbolicAnnotationObject(b509SymbolicAnnotationT, b509SymbolicAnnotationMustJSON(b509SymbolicAnnotationT, b509SymbolicAnnotationAbsent))
			b509SymbolicAnnotationZeroRow := b509SymbolicAnnotationObject(b509SymbolicAnnotationT, b509SymbolicAnnotationMustJSON(b509SymbolicAnnotationT, b509SymbolicAnnotationZero))
			b509SymbolicAnnotationAbsentRow["ab40_marker"] = "ab40_absent"
			b509SymbolicAnnotationZeroRow["ab40_marker"] = "ab40_zero"
			b509SymbolicAnnotationData = b509SymbolicAnnotationMustJSON(b509SymbolicAnnotationT, []any{b509SymbolicAnnotationAbsentRow, b509SymbolicAnnotationZeroRow})
			for b509SymbolicAnnotationPosition, b509SymbolicAnnotationKey := range []B509SymbolicAnnotationKey{b509SymbolicAnnotationAbsent, b509SymbolicAnnotationZero} {
				b509SymbolicAnnotationRaw, b509SymbolicAnnotationPresent, b509SymbolicAnnotationError := B509SymbolicAnnotationLookup(b509SymbolicAnnotationKey)
				b509SymbolicAnnotationWant := []map[string]any{b509SymbolicAnnotationAbsentRow, b509SymbolicAnnotationZeroRow}[b509SymbolicAnnotationPosition]
				if b509SymbolicAnnotationError != nil || !b509SymbolicAnnotationPresent || !reflect.DeepEqual(b509SymbolicAnnotationObject(b509SymbolicAnnotationT, b509SymbolicAnnotationRaw), b509SymbolicAnnotationWant) {
					b509SymbolicAnnotationT.Fatal("ab40 absent/zero collapsed", b509SymbolicAnnotationError)
				}
			}
		})
	}
}
