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

const ab40ExpectedTuples = `any|0D||read_symbolic_address||true
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

func ab40Pointer[ab40Type any](ab40Value ab40Type) *ab40Type { return &ab40Value }

func ab40BaseKey() Ab40Key {
	return Ab40Key{Ab40Family: "Vaillant", Ab40Variant: "none", Ab40Primary: 181, Ab40Secondary: ab40Pointer(uint8(9)), Ab40Operation: ab40Pointer(uint8(13)), Ab40Direction: "none", Ab40Field: "symbolic_annotation", Ab40Alternative: "none", Ab40InstanceScope: "none", Ab40Class: "boiler", Ab40Selector: ab40Pointer("0x1800")}
}

func ab40MustJSON(ab40T *testing.T, ab40Value any) []byte {
	ab40T.Helper()
	ab40Raw, ab40Error := json.Marshal(ab40Value)
	if ab40Error != nil {
		ab40T.Fatal(ab40Error)
	}
	return ab40Raw
}

func ab40Object(ab40T *testing.T, ab40Raw []byte) map[string]any {
	ab40T.Helper()
	var ab40Value map[string]any
	if ab40Error := json.Unmarshal(ab40Raw, &ab40Value); ab40Error != nil {
		ab40T.Fatal(ab40Error)
	}
	return ab40Value
}

func TestAb40FiniteCorrespondence(ab40T *testing.T) {
	ab40Raw, ab40Error := Ab40Entries()
	if ab40Error != nil {
		ab40T.Fatal(ab40Error)
	}
	var ab40Rows []json.RawMessage
	if ab40Error = json.Unmarshal(ab40Raw, &ab40Rows); ab40Error != nil {
		ab40T.Fatal(ab40Error)
	}
	if len(ab40Rows) != 108 {
		ab40T.Fatalf("ab40 rows: %d", len(ab40Rows))
	}
	ab40Counts := map[string]int{}
	for ab40Position, ab40Line := range strings.Split(ab40ExpectedTuples, "\n") {
		ab40Parts := strings.Split(ab40Line, "|")
		ab40Operation, ab40Error := strconv.ParseUint(ab40Parts[1], 16, 8)
		if ab40Error != nil {
			ab40T.Fatal(ab40Error)
		}
		ab40Key := ab40BaseKey()
		ab40Key.Ab40Class = ab40Parts[0]
		ab40Key.Ab40Operation = ab40Pointer(uint8(ab40Operation))
		ab40Key.Ab40Selector = nil
		if ab40Parts[2] != "" {
			ab40Key.Ab40Selector = ab40Pointer(ab40Parts[2])
		}
		ab40Want := map[string]any{
			"family": "Vaillant", "variant": "none", "primary": float64(181), "secondary": float64(9), "operation": float64(ab40Operation), "direction": "none", "field": "symbolic_annotation", "alternative": "none", "instance_scope": "none", "class": ab40Parts[0], "name": ab40Parts[3], "unresolved": ab40Parts[5] == "true", "representation_grade": "Proven", "grade": "Hypothesis", "physical_qualification": "Unknown", "width_bytes": nil, "signedness": "Unknown", "scale": nil, "native_effect": "Unknown", "raw_fallback": "preserve-raw-bytes", "raw_form": ab40Line,
		}
		if ab40Parts[2] != "" {
			ab40Want["selector"] = ab40Parts[2]
		}
		if ab40Parts[4] != "" {
			ab40Want["unit"] = ab40Parts[4]
		}
		if ab40Got := ab40Object(ab40T, ab40Rows[ab40Position]); !reflect.DeepEqual(ab40Got, ab40Want) {
			ab40T.Fatalf("ab40 tuple %d changed: %#v", ab40Position, ab40Got)
		}
		ab40Found, ab40Present, ab40Error := Ab40Lookup(ab40Key)
		if ab40Error != nil || !ab40Present || !bytes.Equal(ab40Found, ab40Rows[ab40Position]) {
			ab40T.Fatalf("ab40 lookup tuple %d: %s %t %v", ab40Position, ab40Found, ab40Present, ab40Error)
		}
		ab40Counts["none|"+ab40Parts[1]]++
	}
	if !reflect.DeepEqual(ab40Counts, map[string]int{"none|0D": 78, "none|0E": 1, "none|54": 29}) {
		ab40T.Fatal(ab40Counts)
	}
}

func TestAb40ExactAxes(ab40T *testing.T) {
	ab40Cases := []struct {
		ab40Name   string
		ab40Change func(*Ab40Key)
	}{
		{"ab40_family", func(ab40Key *Ab40Key) { ab40Key.Ab40Family = "ab40_other" }},
		{"ab40_variant", func(ab40Key *Ab40Key) { ab40Key.Ab40Variant = "ab40_other" }},
		{"ab40_primary", func(ab40Key *Ab40Key) { ab40Key.Ab40Primary = 0 }},
		{"ab40_secondary_absent", func(ab40Key *Ab40Key) { ab40Key.Ab40Secondary = nil }},
		{"ab40_secondary_zero", func(ab40Key *Ab40Key) { ab40Key.Ab40Secondary = ab40Pointer(uint8(0)) }},
		{"ab40_operation_absent", func(ab40Key *Ab40Key) { ab40Key.Ab40Operation = nil }},
		{"ab40_operation_zero", func(ab40Key *Ab40Key) { ab40Key.Ab40Operation = ab40Pointer(uint8(0)) }},
		{"ab40_operation_other", func(ab40Key *Ab40Key) { ab40Key.Ab40Operation = ab40Pointer(uint8(14)) }},
		{"ab40_direction", func(ab40Key *Ab40Key) { ab40Key.Ab40Direction = "response" }},
		{"ab40_field", func(ab40Key *Ab40Key) { ab40Key.Ab40Field = "ab40_other" }},
		{"ab40_alternative", func(ab40Key *Ab40Key) { ab40Key.Ab40Alternative = "ab40_other" }},
		{"ab40_group_zero", func(ab40Key *Ab40Key) { ab40Key.Ab40Group = ab40Pointer(uint8(0)) }},
		{"ab40_group_other", func(ab40Key *Ab40Key) { ab40Key.Ab40Group = ab40Pointer(uint8(1)) }},
		{"ab40_instance_scope", func(ab40Key *Ab40Key) { ab40Key.Ab40InstanceScope = "ab40_other" }},
		{"ab40_register_zero", func(ab40Key *Ab40Key) { ab40Key.Ab40Register = ab40Pointer(uint16(0)) }},
		{"ab40_register_other", func(ab40Key *Ab40Key) { ab40Key.Ab40Register = ab40Pointer(uint16(0x1800)) }},
		{"ab40_class", func(ab40Key *Ab40Key) { ab40Key.Ab40Class = "heat_pump" }},
		{"ab40_class_any", func(ab40Key *Ab40Key) { ab40Key.Ab40Class = "any" }},
		{"ab40_selector_absent", func(ab40Key *Ab40Key) { ab40Key.Ab40Selector = nil }},
		{"ab40_selector_zero", func(ab40Key *Ab40Key) { ab40Key.Ab40Selector = ab40Pointer("0x0000") }},
	}
	ab40Original := append([]byte(nil), ab40Data...)
	ab40T.Cleanup(func() { ab40Data = ab40Original })
	for _, ab40Case := range ab40Cases {
		ab40T.Run(ab40Case.ab40Name, func(ab40T *testing.T) {
			ab40Data = ab40Original
			ab40Key := ab40BaseKey()
			if _, ab40Present, ab40Error := Ab40Lookup(ab40Key); ab40Error != nil || !ab40Present {
				ab40T.Fatal("ab40 missing control", ab40Error)
			}
			ab40Case.ab40Change(&ab40Key)
			if ab40Raw, ab40Present, ab40Error := Ab40Lookup(ab40Key); ab40Error != nil || ab40Present || ab40Raw != nil {
				ab40T.Fatal("ab40 inexact match", ab40Error)
			}
			ab40First := ab40Object(ab40T, ab40MustJSON(ab40T, ab40BaseKey()))
			ab40Second := ab40Object(ab40T, ab40MustJSON(ab40T, ab40Key))
			ab40First["ab40_marker"] = "ab40_control"
			ab40Second["ab40_marker"] = "ab40_changed"
			ab40Data = ab40MustJSON(ab40T, []any{ab40First, ab40Second})
			for ab40Position, ab40Want := range []map[string]any{ab40First, ab40Second} {
				ab40LookupKey := ab40BaseKey()
				if ab40Position == 1 {
					ab40LookupKey = ab40Key
				}
				ab40Raw, ab40Present, ab40Error := Ab40Lookup(ab40LookupKey)
				if ab40Error != nil || !ab40Present || !reflect.DeepEqual(ab40Object(ab40T, ab40Raw), ab40Want) {
					ab40T.Fatal("ab40 distinct valid controls collapsed", ab40Error)
				}
			}
		})
	}
}

func TestAb40SelectorTokens(ab40T *testing.T) {
	ab40Key := ab40BaseKey()
	ab40Key.Ab40Class = "any"
	ab40Key.Ab40Selector = nil
	if _, ab40Present, ab40Error := Ab40Lookup(ab40Key); ab40Error != nil || !ab40Present {
		ab40T.Fatal("ab40 selectorless control")
	}
	for _, ab40Class := range []string{"any", "boiler", "heat_pump"} {
		for _, ab40Selector := range []*string{nil, ab40Pointer("0x0000"), ab40Pointer("0")} {
			ab40Key.Ab40Class = ab40Class
			ab40Key.Ab40Selector = ab40Selector
			_, ab40Present, ab40Error := Ab40Lookup(ab40Key)
			if ab40Error != nil || ab40Present != (ab40Class == "any" && ab40Selector == nil) {
				ab40T.Fatal("ab40 wildcard or selectorless/zero collapse")
			}
		}
	}
	ab40Key.Ab40Class = "heat_pump_interface_unit"
	ab40Key.Ab40Operation = ab40Pointer(uint8(0x54))
	ab40Key.Ab40Selector = ab40Pointer("54 02 00 0D 0A")
	if _, ab40Present, ab40Error := Ab40Lookup(ab40Key); ab40Error != nil || !ab40Present {
		ab40T.Fatal("ab40 opaque control")
	}
	for ab40Position := 0; ab40Position < 5; ab40Position++ {
		ab40Tokens := strings.Fields("54 02 00 0D 0A")
		ab40Tokens[ab40Position] = "FF"
		ab40Key.Ab40Selector = ab40Pointer(strings.Join(ab40Tokens, " "))
		if _, ab40Present, ab40Error := Ab40Lookup(ab40Key); ab40Error != nil || ab40Present {
			ab40T.Fatalf("ab40 changed byte %d matched", ab40Position)
		}
	}
	for _, ab40Token := range []string{"0x0D0A", "0D 0A", "02 00 0D 0A", "54 02 00 0D 0A 00"} {
		ab40Key.Ab40Selector = ab40Pointer(ab40Token)
		if _, ab40Present, ab40Error := Ab40Lookup(ab40Key); ab40Error != nil || ab40Present {
			ab40T.Fatal("ab40 narrowed token matched")
		}
	}
}

func TestAb40ValidationAndRetention(ab40T *testing.T) {
	ab40Original := append([]byte(nil), ab40Data...)
	ab40T.Cleanup(func() { ab40Data = ab40Original })
	ab40Row := ab40Object(ab40T, ab40MustJSON(ab40T, ab40BaseKey()))
	ab40Row["unresolved"] = true
	ab40Row["grade"] = "Hypothesis"
	ab40Row["width_bytes"] = nil
	ab40Row["ab40_unknown"] = map[string]any{"ab40_raw": []any{float64(0), float64(255)}}
	ab40Source := ab40MustJSON(ab40T, []any{ab40Row})
	ab40Rows, ab40Index, ab40Error := ab40Validate(ab40Source)
	if ab40Error != nil {
		ab40T.Fatal(ab40Error)
	}
	for ab40Position := range ab40Source {
		ab40Source[ab40Position] = ' '
	}
	if !reflect.DeepEqual(ab40Object(ab40T, ab40Rows[0]), ab40Row) || !reflect.DeepEqual(ab40Object(ab40T, ab40Index[ab40Identity(ab40BaseKey())]), ab40Row) {
		ab40T.Fatal("ab40 source aliases parsed rows")
	}
	ab40Data = ab40MustJSON(ab40T, []any{ab40Row})
	ab40Raw, ab40Present, ab40Error := Ab40Lookup(ab40BaseKey())
	if ab40Error != nil || !ab40Present || !reflect.DeepEqual(ab40Object(ab40T, ab40Raw), ab40Row) {
		ab40T.Fatal("ab40 absent name/unit or unresolved lost")
	}
	for _, ab40Axis := range []string{"grade", "width_bytes"} {
		ab40Changed := ab40Object(ab40T, ab40MustJSON(ab40T, ab40Row))
		if ab40Axis == "grade" {
			ab40Changed[ab40Axis] = "Unknown"
		} else {
			ab40Changed[ab40Axis] = float64(2)
		}
		ab40Data = ab40MustJSON(ab40T, []any{ab40Changed})
		ab40Raw, ab40Present, ab40Error = Ab40Lookup(ab40BaseKey())
		if ab40Error != nil || !ab40Present || !reflect.DeepEqual(ab40Object(ab40T, ab40Raw), ab40Changed) {
			ab40T.Fatal("ab40 nonidentity metadata not retained")
		}
		ab40Data = ab40MustJSON(ab40T, []any{ab40Row, ab40Changed})
		if _, _, ab40Error = ab40Validate(ab40Data); ab40Error == nil {
			ab40T.Fatal("ab40 grade/width cannot distinguish duplicate identities")
		}
	}
	ab40Bad := [][]byte{[]byte(`null`), []byte(`{}`), []byte(`[`), []byte(`[null]`), []byte(`[1]`), []byte(`[{}]`), []byte(`[] true`), ab40MustJSON(ab40T, []any{ab40Row, ab40Row})}
	for _, ab40Axis := range []string{"primary", "operation", "secondary", "group", "register", "selector"} {
		ab40Changed := ab40Object(ab40T, ab40MustJSON(ab40T, ab40Row))
		ab40Changed[ab40Axis] = nil
		ab40Bad = append(ab40Bad, ab40MustJSON(ab40T, []any{ab40Changed}))
		ab40Changed[ab40Axis] = map[string]any{}
		ab40Bad = append(ab40Bad, ab40MustJSON(ab40T, []any{ab40Changed}))
	}
	for ab40Position, ab40Input := range ab40Bad {
		ab40T.Run(fmt.Sprintf("ab40_invalid_%d", ab40Position), func(ab40T *testing.T) {
			if _, _, ab40Error := ab40Validate(ab40Input); ab40Error == nil {
				ab40T.Fatal("ab40 invalid accepted")
			}
			ab40Data = ab40Input
			if ab40Raw, ab40Error := Ab40Entries(); ab40Error == nil || ab40Raw != nil {
				ab40T.Fatal("ab40 entries accepted invalid data")
			}
			if ab40Raw, ab40Present, ab40Error := Ab40Lookup(ab40BaseKey()); ab40Error == nil || ab40Present || ab40Raw != nil {
				ab40T.Fatal("ab40 lookup accepted invalid data")
			}
		})
	}
}

func TestAb40IndependentStorage(ab40T *testing.T) {
	ab40First, ab40Error := Ab40Entries()
	if ab40Error != nil {
		ab40T.Fatal(ab40Error)
	}
	ab40Expected := bytes.Clone(ab40First)
	ab40Second, ab40Error := Ab40Entries()
	if ab40Error != nil {
		ab40T.Fatal(ab40Error)
	}
	for ab40Position := range ab40First {
		ab40First[ab40Position] ^= 0xff
	}
	ab40Third, ab40Error := Ab40Entries()
	if ab40Error != nil || !bytes.Equal(ab40Expected, ab40Second) || !bytes.Equal(ab40Expected, ab40Third) {
		ab40T.Fatal("ab40 entries share storage")
	}
	ab40First, ab40Present, ab40Error := Ab40Lookup(ab40BaseKey())
	if ab40Error != nil || !ab40Present {
		ab40T.Fatal("ab40 lookup missing")
	}
	ab40Expected = bytes.Clone(ab40First)
	ab40Second, ab40Present, ab40Error = Ab40Lookup(ab40BaseKey())
	if ab40Error != nil || !ab40Present {
		ab40T.Fatal("ab40 second lookup missing")
	}
	for ab40Position := range ab40First {
		ab40First[ab40Position] ^= 0xff
	}
	ab40Third, ab40Present, ab40Error = Ab40Lookup(ab40BaseKey())
	if ab40Error != nil || !ab40Present || !bytes.Equal(ab40Expected, ab40Second) || !bytes.Equal(ab40Expected, ab40Third) {
		ab40T.Fatal("ab40 lookups share storage")
	}
}

func TestAb40OptionalAbsentZero(ab40T *testing.T) {
	ab40Original := append([]byte(nil), ab40Data...)
	ab40T.Cleanup(func() { ab40Data = ab40Original })
	for _, ab40Axis := range []string{"secondary", "operation", "group", "register", "selector"} {
		ab40T.Run("ab40_"+ab40Axis, func(ab40T *testing.T) {
			ab40Absent := ab40BaseKey()
			ab40Zero := ab40BaseKey()
			switch ab40Axis {
			case "secondary":
				ab40Absent.Ab40Secondary = nil
				ab40Zero.Ab40Secondary = ab40Pointer(uint8(0))
			case "operation":
				ab40Absent.Ab40Operation = nil
				ab40Zero.Ab40Operation = ab40Pointer(uint8(0))
			case "group":
				ab40Absent.Ab40Group = nil
				ab40Zero.Ab40Group = ab40Pointer(uint8(0))
			case "register":
				ab40Absent.Ab40Register = nil
				ab40Zero.Ab40Register = ab40Pointer(uint16(0))
			case "selector":
				ab40Absent.Ab40Selector = nil
				ab40Zero.Ab40Selector = ab40Pointer("0x0000")
			}
			ab40AbsentRow := ab40Object(ab40T, ab40MustJSON(ab40T, ab40Absent))
			ab40ZeroRow := ab40Object(ab40T, ab40MustJSON(ab40T, ab40Zero))
			ab40AbsentRow["ab40_marker"] = "ab40_absent"
			ab40ZeroRow["ab40_marker"] = "ab40_zero"
			ab40Data = ab40MustJSON(ab40T, []any{ab40AbsentRow, ab40ZeroRow})
			for ab40Position, ab40Key := range []Ab40Key{ab40Absent, ab40Zero} {
				ab40Raw, ab40Present, ab40Error := Ab40Lookup(ab40Key)
				ab40Want := []map[string]any{ab40AbsentRow, ab40ZeroRow}[ab40Position]
				if ab40Error != nil || !ab40Present || !reflect.DeepEqual(ab40Object(ab40T, ab40Raw), ab40Want) {
					ab40T.Fatal("ab40 absent/zero collapsed", ab40Error)
				}
			}
		})
	}
}
