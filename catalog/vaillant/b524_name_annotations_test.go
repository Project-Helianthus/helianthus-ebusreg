package vaillant

import (
	b524NameAnnotationBytes "bytes"
	b524NameAnnotationJSON "encoding/json"
	b524NameAnnotationReflect "reflect"
	b524NameAnnotationStrconv "strconv"
	b524NameAnnotationStrings "strings"
	b524NameAnnotationTesting "testing"
)

const b524NameAnnotationExpected = `
00|0001|dhw_outdoor_bivalence_temperature_limit|false|
00|0002|continuous_heating_outdoor_temperature_limit|false|
00|0003|frost_protection_delay_time|false|
00|0004|preheating_time_maximum|false|
00|0006|manual_cooling_time|false|
00|0007|unknown_11|true|Polarity unresolved, so no Boolean name is given: the catalog row gives ` + "`" + `0=off 1=on` + "`" + ` (1 = switched on), while the candidate write sets the switch-off with 01 and cancels it with 00 (1 = switched off), unconfirmed on the wire.
00|0008|cooling_enabled_at_holiday|true|Candidate interpretation, unconfirmed on the wire: 1 means that cooling is permitted during a holiday; the polarity is not established and no value table is documented.
00|0009|automatic_cooling_enabled|false|
00|000A|cylinder_parallel_loading_enabled|false|
00|000B|solar_circuit_bleeding_time|true|Candidate interpretation, unconfirmed on the wire: a duration of the solar circuit bleeding; the value class and unit are not established.
00|000D|solar_frost_protection_limit|true|Candidate interpretation, unconfirmed on the wire: a limit of the solar frost protection; the limited quantity and the representation are not established.
00|000E|room_humidity_limit_maximum|false|
00|000F|hybrid_control_mode|false|
00|0010|backup_heater_tariff|true|Value class is not established: a price, a ratio or a selection are candidate readings, unconfirmed on the wire; the association with a tariff is not established.
00|0011|electricity_price_at_low_tariff|true|Candidate interpretation, unconfirmed on the wire: an electricity price at the low tariff; value class and unit are not established.
00|0012|electricity_price_at_high_tariff|true|Candidate interpretation, unconfirmed on the wire: an electricity price at the high tariff; the register map decodes a u16 in °C, which does not fit a price; unresolved.
00|0013|backup_heater_type|true|Candidate interpretation, unconfirmed on the wire: an enumeration of the backup heater kind (0 condensing, 1 non-condensing, 2 electrical); the representation and the applicability are not established.
00|0014|adaptive_heating_curve_enabled|false|
00|0015|automatic_daylight_saving_time|true|Value class and encoding are not established.
00|0016|ventilation_operating_mode|false|
00|0017|dhw_loading_time_maximum|false|
00|0018|dhw_blocking_time|false|
00|0019|solar_flow_rate_setpoint|false|
00|001B|dhw_pump_overrun_time|false|
00|001C|heat_generator_flow_temperature_limit_maximum|false|
00|001E|automatic_sequence_reversal|true|Value class, encoding and meaning are not established; candidate reading, unconfirmed on the wire: a pump or flag value related to the reversal of the generator sequence.
00|0022|heating_outdoor_alternative_temperature_limit|false|
00|0023|heating_outdoor_bivalence_temperature_limit|false|
00|0024|backup_heater_mode|true|Candidate interpretation, unconfirmed on the wire: an enumeration of the backup heater use (0 disabled, 1 heating, 2 DHW, 3 DHW and heating); the value table is not established.
00|0025|temporary_backup_heater_mode|true|Candidate interpretation, unconfirmed on the wire: an enumeration for a temporary permission of the backup heater; how it differs from the backup heater enumeration at RR=0x0024 is not established, and neither is the value table.
00|0026|heating_flow_temperature_setpoint_maximum_at_heat_pump_fault|false|
00|0027|cylinder_loading_start_hysteresis|false|
00|0029|cylinder_loading_offset|false|
00|002A|dhw_legionella_protection_start_clock_time|false|
00|002B|dhw_legionella_protection_start_weekday|false|
00|002C|maintenance_date|false|
00|002D|outdoor_temperature_offset|false|
00|002F|functional_module_configuration_code|true|Candidate interpretation, unconfirmed on the wire: a configuration code of a functional module kind (the VR71 module); the code list is not established.
00|0031|language|true|Value class and encoding are not established.
00|0034|date|false|
00|0035|clock_time|false|
00|0036|hydraulic_scheme_code|false|
00|0038|cooling_start_temperature_limit|false|
00|0039|water_pressure|false|
00|003A|hybrid_unknown_2|true|A kelvin value of the hybrid control; the meaning is not established.
00|003C|solar_yield_total_reset_command|false|
00|003D|solar_yield_total|false|
00|003E|environmental_yield_total|false|
00|003F|environmental_yield_total_reset_command|true|Candidate interpretation, unconfirmed on the wire: a command that resets the environmental yield total; representation and applicability are not established.
00|0040|ventilation_volatile_organic_compound_sensor_1|true|Representation is not established; the relation to ` + "`" + `OP=0x02, GG=0x09, RR=0x0005` + "`" + ` (same name) is not established.
00|0041|ventilation_volatile_organic_compound_sensor_2|true|Representation is not established; the relation to ` + "`" + `OP=0x02, GG=0x09, RR=0x0006` + "`" + ` (same name) is not established.
00|0042|ventilation_fan_speed_limit_maximum_at_day|false|
00|0043|ventilation_fan_speed_limit_maximum_at_night|false|
00|0044|ventilation_heat_recovery|true|Value class is not established; the relation to ` + "`" + `OP=0x02, GG=0x09, RR=0x000B` + "`" + ` (same name) is not established.
00|0045|utility_lockout_mode|true|Candidate interpretation, unconfirmed on the wire: an enumeration of utility lockouts (0 disable heat pump, 1 disable backup heater, 2 disable heat pump and backup heater, 3 disable heating, 4 disable cooling, 5 disable heating and cooling); the value table is not established.
00|0046|dhw_flow_temperature_setpoint_maximum|false|
00|0047|emergency_operation_active|true|Candidate interpretation, unconfirmed on the wire: 1 means that the device runs in an emergency (failure) mode; no value table is documented.
00|0048|energy_manager_state|false|
00|0049|electricity_consumption_total_reset_command|true|Candidate interpretation, unconfirmed on the wire: a command that resets an electricity consumption total; representation is not established, and which accumulated value is reset is not established.
00|004A|settings_reset_command|true|Candidate interpretation, unconfirmed on the wire: a command that resets settings to their defaults; which settings are reset, the representation and the writability are not established.
00|004B|flow_temperature|false|
00|004C|backup_heater_power_limit_maximum|true|Candidate interpretation, unconfirmed on the wire: the maximum power of the backup heater of the heat-pump interface unit; representation is not established.
00|004D|multifunction_output_2_mode|false|
00|004E|heating_fuel_consumption_current_month|false|
00|004F|heating_electricity_consumption_current_month|false|
00|0050|dhw_electricity_consumption_current_month|false|
00|0051|dhw_fuel_consumption_current_month|false|
00|0052|heating_fuel_consumption_previous_month|false|
00|0053|heating_electricity_consumption_previous_month|false|
00|0054|dhw_electricity_consumption_previous_month|false|
00|0055|dhw_fuel_consumption_previous_month|false|
00|0056|heating_fuel_consumption_total|false|
00|0057|heating_electricity_consumption_total|false|
00|0058|dhw_electricity_consumption_total|false|
00|0059|dhw_fuel_consumption_total|false|
00|005C|electricity_consumption_1_total|true|Coexists with ` + "`" + `RR=0x00C6` + "`" + ` (` + "`" + `electricity_consumption_total` + "`" + `); the documented terms differ (a consumption sum, and the consumption for the system history), so the index only keeps the two apart; whether the two registers carry the same total is not established. Candidate reading, unconfirmed on the wire: a total electricity consumption in kWh.
00|005D|fuel_consumption_1_total|true|Coexists with ` + "`" + `RR=0x00CD` + "`" + ` (` + "`" + `fuel_consumption_total` + "`" + `); the documented terms differ (a consumption sum, and the consumption for the system history), so the index only keeps the two apart; whether the two registers carry the same total is not established. Candidate reading, unconfirmed on the wire: a total fuel consumption in kWh.
00|005E|heat_recovery_yield_total|false|
00|005F|actuator_sensor_test_device_code|true|Candidate interpretation, unconfirmed on the wire: a code that selects the device of an actuator and sensor test; the value class and the code list are not established.
00|0060|actuator_sensor_test_sensor_code|true|Candidate interpretation, unconfirmed on the wire: a code that selects the sensor of an actuator and sensor test; the value class and the code list are not established.
00|0061|actuator_sensor_test_actuator_code|true|Candidate interpretation, unconfirmed on the wire: a code that selects the actuator of an actuator and sensor test; the value class and the code list are not established.
00|0062|ventilation_volatile_organic_compound_sensor_3|true|Representation is not established.
00|0063|ventilation_volatile_organic_compound_limit_maximum|true|Representation is not established; the relation to ` + "`" + `OP=0x02, GG=0x09, RR=0x000C` + "`" + ` (same name) is not established.
00|0065|ventilation_fan_speed_setpoint|false|
00|0066|home_screen_unknown_6|true|Value class and meaning are not established.
00|0067|cylinder_bottom_temperature|false|
00|0068|cylinder_top_temperature|false|
00|0069|unknown_1|true|Candidate interpretation, unconfirmed on the wire: an enumeration whose candidate labels are the controller and zones 1 to 3; the meaning is not established.
00|006A|multifunction_input_mode|true|Candidate interpretation, unconfirmed on the wire: an enumeration of the multifunction input configuration; the value table is not established.
00|006C|installer_name_prefix|false|
00|006D|installer_name_suffix|false|
00|006F|installer_phone_number_prefix|false|
00|0070|installer_phone_number_suffix|false|
00|0071|holiday_start_date|true|Family split: ` + "`" + `RR=0x0071` + "`" + ` is attested on the VRC700 family and ` + "`" + `RR=0x00F2` + "`" + ` on the VRC720 family (candidate, unconfirmed on the wire); whether the two selectors carry the same quantity is not established.
00|0072|holiday_end_date|true|Family split: ` + "`" + `RR=0x0072` + "`" + ` is attested on the VRC700 family and ` + "`" + `RR=0x00F3` + "`" + ` on the VRC720 family (candidate, unconfirmed on the wire); whether the two selectors carry the same quantity is not established.
00|0073|outdoor_temperature|false|
00|0074|special_function_mode|true|Candidate interpretation, unconfirmed on the wire: an enumeration of the system special function, dormant when none is active; the mode enumeration and its write path remain hypotheses.
00|0076|installer_access_code|false|
00|0077|installation_assistant_ready|true|Candidate interpretation, unconfirmed on the wire: 1 means that the installation assistant has completed; no value table is documented.
00|0078|home_screen_special_function|true|Representation and value class are not established.
00|0079|heating_zone_time_program_active|true|Representation is not established; the relation to ` + "`" + `OP=0x02, GG=0x03, RR=0x0019` + "`" + ` (same name) is not established.
00|007A|cooling_zone_time_program_active|true|Representation is not established; the relation to ` + "`" + `OP=0x02, GG=0x03, RR=0x001A` + "`" + ` (same name) is not established.
00|007B|heating_zone_operating_mode|true|Representation is not established; the relation to ` + "`" + `OP=0x02, GG=0x03, RR=0x0006` + "`" + ` (same name) is not established.
00|007C|cooling_zone_operating_mode|true|Representation is not established; the relation to ` + "`" + `OP=0x02, GG=0x03, RR=0x0001` + "`" + ` (same name) is not established.
00|007D|home_screen_operating_mode|true|Representation is not established.
00|007E|unknown_3|true|A float32 value; the meaning is not established.
00|007F|screed_drying_time|true|Two readings, unresolved (candidate interpretations, unconfirmed on the wire): the elapsed duration of the screed drying in days, or the index of the current day in the screed drying programme; representation is not established.
00|0080|cylinder_loading_start_offset|false|
00|0081|cylinder_loading_stop_offset|false|
00|0082|room_temperature_setpoint_at_holiday|true|Family split: ` + "`" + `RR=0x0082` + "`" + ` is attested on the VRC700 family and ` + "`" + `RR=0x00F6` + "`" + ` on the VRC720 family (candidate, unconfirmed on the wire); whether the two selectors carry the same quantity is not established. Representation of ` + "`" + `RR=0x0082` + "`" + ` is not established.
00|0083|bank_holiday_start_date|true|Representation is not established.
00|0084|bank_holiday_end_date|true|Representation is not established.
00|0085|actuator_sensor_test_sensor_reading|true|Representation is not established.
00|0086|backlight_on_time|true|Candidate interpretation, unconfirmed on the wire: the time until the backlight switches off; the unit is not established.
00|0087|drainback_start_time|true|Candidate interpretation, unconfirmed on the wire: a duration, an interval or a clock time of the drainback start; which of the three is not established.
00|0088|drainback_station|true|Representation and value class are not established.
00|0089|solar_pump_kick_time|false|
00|008A|solar_pump_kick_gradient_limit_minimum|true|Candidate interpretation, unconfirmed on the wire: a minimum gradient limit of the solar pump kick; the quantity of the gradient and its unit are not established.
00|008B|heat_pump_flow_temperature_setpoint_maximum|false|
00|008C|boiler_remote_reset_available|true|Candidate interpretation, unconfirmed on the wire: 1 means that a remote reset of the boiler is available; no value table is documented and the polarity is not established.
00|008D|display_test_pattern|true|Value class is not established.
00|008E|backlight_mode|false|
00|008F|cylinder_type|true|Two readings, unresolved: an enumeration (0 disabled, 1 heating, 2 DHW, 3 DHW and heating) or an availability flag (candidate interpretations, unconfirmed on the wire).
00|0090|unknown_4|true|Representation and meaning are not established.
00|0091|left_button_active|true|Candidate interpretation, unconfirmed on the wire: 1 means that the button is pressed; the value class, the representation and the polarity are not established.
00|0092|right_button_active|true|Candidate interpretation, unconfirmed on the wire: 1 means that the button is pressed; the value class, the representation and the polarity are not established.
00|0093|encoder_active|true|Representation is not established: whether the value is a flag or a rotation count is open. Candidate reading, unconfirmed on the wire: a flag that the encoder was rotated.
00|0094|external_crystal_status|true|Candidate interpretation, unconfirmed on the wire: a status of the external crystal; the representation and the value table are not established.
00|0095|outdoor_temperature_daily_average|false|
00|0096|maintenance_active|true|Two readings, unresolved (candidate interpretations, unconfirmed on the wire): 1 means that the device is in maintenance mode, or 1 means that maintenance is due; the semantic plane projects this selector as ` + "`" + `state.maintenance_due` + "`" + `.
00|0097|eeprom_update_active|true|Candidate interpretation, unconfirmed on the wire: 1 means that an EEPROM update is in progress; no value table is documented and the polarity is not established.
00|0098|boiler_reset_command|false|
00|0099|rtc_crystal_ok|true|Candidate interpretation, unconfirmed on the wire: 1 means that the real-time clock crystal is OK; the polarity and the representation are not established.
00|009A|unknown_5|true|A boolean value; the meaning is not established.
00|009B|time_program_reset_command|false|
00|009C|home_screen_headline|true|Representation is not established.
00|009D|dhw_cylinder_top_temperature|false|
00|009E|dhw_cylinder_bottom_temperature|false|
00|009F|heating_cylinder_top_temperature|false|
00|00A0|heating_cylinder_bottom_temperature|false|
00|00A1|pv_surplus_enabled|true|Candidate interpretation, unconfirmed on the wire: 1 means that the use of a photovoltaic surplus is enabled; no value table is documented and the polarity is not established.
00|00A2|heating_cylinder_offset_at_pv_surplus|false|
00|00A3|ventilation_time_program_active|true|Representation is not established; the relation to ` + "`" + `OP=0x02, GG=0x09, RR=0x0003` + "`" + ` (same name) is not established.
00|00A4|ventilation_special_function_state|true|Candidate interpretation, unconfirmed on the wire: an enumeration of the ventilation special function; the representation and the applicability are not established.
00|00A5|cascade_client_sequence|true|Candidate interpretation, unconfirmed on the wire: a position in the cascade client sequence; the value class and the meaning are not established.
00|00A6|menu_button_active|true|Candidate interpretation, unconfirmed on the wire: 1 means that the button is pressed; the value class, the representation and the polarity are not established.
00|00A7|confirm_button_active|true|Candidate interpretation, unconfirmed on the wire: 1 means that the button is pressed; the value class, the representation and the polarity are not established.
00|00A8|back_button_active|true|Candidate interpretation, unconfirmed on the wire: 1 means that the button is pressed; the value class, the representation and the polarity are not established.
00|00A9|up_button_active|true|Candidate interpretation, unconfirmed on the wire: 1 means that the button is pressed; the value class, the representation and the polarity are not established.
00|00AA|down_button_active|true|Candidate interpretation, unconfirmed on the wire: 1 means that the button is pressed; the value class, the representation and the polarity are not established.
00|00AB|isolation_switch_active|true|Candidate interpretation, unconfirmed on the wire: 1 means that the isolation switch is active; which switch position 1 denotes is not established.
00|00AF|cooling_electricity_consumption_current_month|false|
00|00B0|cooling_electricity_consumption_previous_month|false|
00|00B1|cooling_electricity_consumption_total|false|
00|00B2|network_identifier|true|Semantics are unknown; candidate reading, unconfirmed on the wire: an identifier of the network.
00|00B3|multifunction_output_1_mode|true|Candidate interpretation, unconfirmed on the wire: an enumeration of the first multifunction output configuration; the value table is not established.
00|00B5|outdoor_temperature_sensor_ok|false|
00|00B6|external_energy_management_enabled|true|Candidate interpretation, unconfirmed on the wire: a one-byte switch for an external energy management, 1 = enabled and 0 = disabled; the meanings of 0 and 1 and the controller reaction are not established.
00|00B8|solar_diverter_valve|true|Value class is not established.
00|00B9|cooling_electricity_consumption_current_year|false|
00|00BA|cooling_electricity_consumption_previous_year|false|
00|00BB|dhw_electricity_consumption_current_year|false|
00|00BC|dhw_electricity_consumption_previous_year|false|
00|00BD|heating_fuel_consumption_previous_year|false|
00|00BF|dhw_fuel_consumption_previous_year|false|
00|00C0|heating_electricity_consumption_previous_year|false|
00|00C1|heating_electricity_consumption_current_year|false|
00|00C2|electricity_consumption_current_month|false|
00|00C3|electricity_consumption_previous_month|false|
00|00C4|electricity_consumption_current_year|false|
00|00C5|electricity_consumption_previous_year|false|
00|00C6|electricity_consumption_total|true|Coexists with ` + "`" + `RR=0x005C` + "`" + ` (` + "`" + `electricity_consumption_1_total` + "`" + `); the documented terms differ, and whether the two registers carry the same total is not established. Candidate reading, unconfirmed on the wire: a total electricity consumption in kWh for the system history.
00|00C7|dhw_fuel_consumption_current_year|false|
00|00C8|heating_fuel_consumption_current_year|false|
00|00C9|fuel_consumption_current_month|false|
00|00CA|fuel_consumption_previous_month|false|
00|00CB|fuel_consumption_current_year|false|
00|00CC|fuel_consumption_previous_year|false|
00|00CD|fuel_consumption_total|true|Coexists with ` + "`" + `RR=0x005D` + "`" + ` (` + "`" + `fuel_consumption_1_total` + "`" + `); the documented terms differ, and whether the two registers carry the same total is not established. Candidate reading, unconfirmed on the wire: a total fuel consumption in kWh for the system history.
00|00CE|heat_recovery_yield_previous_month|false|
00|00CF|heat_recovery_yield_current_month|false|
00|00D0|heat_recovery_yield_current_year|false|
00|00D1|heat_recovery_yield_previous_year|false|
00|00D2|solar_yield_previous_month|false|
00|00D3|solar_yield_current_month|false|
00|00D4|solar_yield_current_year|false|
00|00D5|solar_yield_previous_year|false|
00|00D6|environmental_yield_previous_month|false|
00|00D7|environmental_yield_previous_year|false|
00|00D8|environmental_yield_current_year|false|
00|00D9|electricity_tariff_type|false|
00|00DA|manual_cooling_start_date|true|Candidate interpretation, unconfirmed on the wire: the start date of a manual cooling period on the VRC700 family; the interpretation, the writability and the family scope are not established.
00|00DB|manual_cooling_end_date|true|Candidate interpretation, unconfirmed on the wire: the end date of a manual cooling period on the VRC700 family; the interpretation, the writability and the family scope are not established.
00|00DC|room_temperature_control_mode|true|Representation is not established.
00|00DD|screed_drying_profile_1_temperature_setpoint|false|
00|00DE|screed_drying_profile_2_temperature_setpoint|false|
00|00DF|screed_drying_profile_3_temperature_setpoint|false|
00|00E0|environmental_yield_current_month|false|
00|00E1|circuit_control_mode|true|Representation is not established.
00|00E2|automatic_control_active|true|Representation is not established.
00|00E3|up_button_level|true|Candidate interpretation, unconfirmed on the wire: a raw value of the button input as a float32; only a zero has been observed, so whether it is a graded amount, a raw reading or a flag is not established.
00|00E4|down_button_level|true|Candidate interpretation, unconfirmed on the wire: a raw value of the button input as a float32; only a zero has been observed, so whether it is a graded amount, a raw reading or a flag is not established.
00|00E5|help_button_level|true|Candidate interpretation, unconfirmed on the wire: a raw value of the button input as a float32; only a zero has been observed, so whether it is a graded amount, a raw reading or a flag is not established.
00|00E6|confirm_button_level|true|Candidate interpretation, unconfirmed on the wire: a raw value of the button input as a float32; only a zero has been observed, so whether it is a graded amount, a raw reading or a flag is not established.
00|00E7|menu_button_level|true|Candidate interpretation, unconfirmed on the wire: a raw value of the button input as a float32; only a zero has been observed, so whether it is a graded amount, a raw reading or a flag is not established.
00|00E8|back_button_level|true|Candidate interpretation, unconfirmed on the wire: a raw value of the button input as a float32; only a zero has been observed, so whether it is a graded amount, a raw reading or a flag is not established.
00|00E9|slider_button_1_level|true|Candidate interpretation, unconfirmed on the wire: a raw value of the button input as a float32; only a zero has been observed, so whether it is a graded amount, a raw reading or a flag is not established.
00|00EA|slider_button_2_level|true|Candidate interpretation, unconfirmed on the wire: a raw value of the button input as a float32; only a zero has been observed, so whether it is a graded amount, a raw reading or a flag is not established.
00|00EB|slider_button_3_level|true|Candidate interpretation, unconfirmed on the wire: a raw value of the button input as a float32; only a zero has been observed, so whether it is a graded amount, a raw reading or a flag is not established.
00|00EC|slider_button_4_level|true|Candidate interpretation, unconfirmed on the wire: a raw value of the button input as a float32; only a zero has been observed, so whether it is a graded amount, a raw reading or a flag is not established.
00|00ED|slider_button_5_level|true|Candidate interpretation, unconfirmed on the wire: a raw value of the button input as a float32; only a zero has been observed, so whether it is a graded amount, a raw reading or a flag is not established.
00|00EE|slider_button_6_level|true|Candidate interpretation, unconfirmed on the wire: a raw value of the button input as a float32; only a zero has been observed, so whether it is a graded amount, a raw reading or a flag is not established.
00|00EF|continuous_heating_room_temperature_setpoint|false|
00|00F0|priority_error_code|true|Candidate interpretation, unconfirmed on the wire: the error number of highest priority; whether the value ranks errors by priority is not established.
00|00F1|priority_maintenance_code|true|Candidate interpretation, unconfirmed on the wire: the maintenance number of highest priority; whether the value ranks maintenance entries by priority is not established.
00|00F2|holiday_start_date|true|Family split: ` + "`" + `RR=0x0071` + "`" + ` is attested on the VRC700 family and ` + "`" + `RR=0x00F2` + "`" + ` on the VRC720 family (candidate, unconfirmed on the wire); whether the two selectors carry the same quantity is not established.
00|00F3|holiday_end_date|true|Family split: ` + "`" + `RR=0x0072` + "`" + ` is attested on the VRC700 family and ` + "`" + `RR=0x00F3` + "`" + ` on the VRC720 family (candidate, unconfirmed on the wire); whether the two selectors carry the same quantity is not established.
00|00F4|holiday_start_clock_time|true|Candidate interpretation, unconfirmed on the wire: the start time of the system holiday; the VRC720-family set sequence writes an end date here, and the disagreement is unresolved.
00|00F5|holiday_end_clock_time|false|
00|00F6|room_temperature_setpoint_at_holiday|true|Family split: ` + "`" + `RR=0x0082` + "`" + ` is attested on the VRC700 family and ` + "`" + `RR=0x00F6` + "`" + ` on the VRC720 family (candidate, unconfirmed on the wire); whether the two selectors carry the same quantity is not established.
00|00F7|holiday_abort_command|false|
00|00F8|manual_cooling_abort_command|false|
00|00F9|ventilation_boost_active|false|
00|00FA|external_demand_input_mode|false|
00|00FB|heat_pump_temperature_difference_limit_minimum|true|Candidate interpretation, unconfirmed on the wire: a minimum limit of a temperature difference; between which temperatures the difference is taken is not established.
00|00FC|dhw_heat_pump_flow_temperature_setpoint_maximum|false|
00|00FD|boiler_flow_temperature_setpoint_maximum|true|Representation is not established.
01|0001|dhw_circuit_type|true|Candidate interpretation, unconfirmed on the wire: an enumeration of the DHW circuit (0 deactivated, 1 DHW circuit); the availability test and a circuit-type reading remain distinct and the value table is not established.
01|0002|dhw_circulation_pump_active|false|
01|0003|dhw_operating_mode|false|
01|0004|dhw_temperature_setpoint|false|
01|0005|dhw_temperature|false|
01|0006|dhw_pump_active|false|
01|0007|dhw_time_program_reset_command|false|
01|0008|dhw_flow_temperature_setpoint|false|
01|0009|dhw_holiday_start_date|false|
01|000A|dhw_holiday_end_date|false|
01|000B|dhw_bank_holiday_start_date|false|
01|000C|dhw_bank_holiday_end_date|false|
01|000D|dhw_special_function_mode|false|
01|000E|dhw_central_control_active|true|Candidate interpretation, unconfirmed on the wire: 1 means that the DHW circuit is controlled by the system controller; no value table is documented.
01|000F|dhw_special_function_state|false|
01|0010|dhw_holiday_start_clock_time|false|
01|0011|dhw_holiday_end_clock_time|false|
01|0012|dhw_holiday_abort_command|false|
01|0014|dhw_preference|true|Representation and meaning are not established; candidate reading, unconfirmed on the wire: a preference of the DHW.
01|0015|dhw_recharge_setpoint|true|Candidate interpretation, unconfirmed on the wire: a setpoint of the DHW recharge; the regulated quantity and the representation are not established.
01|0019|dhw_loading_efficiency_state|true|Candidate interpretation, unconfirmed on the wire: an enumeration of the DHW loading efficiency; the representation and the value table are not established.
02|0001|hydraulic_circuit_type|true|Candidate interpretation, unconfirmed on the wire: an enumeration of the circuit kind (1 direct heating circuit, 2 external mixer circuit); the purpose is unverified.
02|0002|circuit_type|false|
02|0003|circuit_room_influence_mode|false|
02|0004|circuit_return_temperature_setpoint|false|
02|0005|circuit_condensate_sensitive_emitter_available|false|
02|0006|circuit_cooling_enabled|false|
02|0007|circuit_flow_temperature_setpoint|false|
02|0008|circuit_flow_temperature|false|
02|0009|dhw_circuit_temperature_setpoint|false|
02|000A|circuit_unknown_7|true|A kelvin value, present with cooling; the meaning is not established.
02|000B|heating_circuit_flow_temperature_offset|false|
02|000C|circuit_flow_high_temperature_setpoint|false|
02|000D|circuit_flow_low_temperature_setpoint|false|
02|000E|circuit_frost_protection_mode|true|Two readings, unresolved (candidate interpretations, unconfirmed on the wire): a frost protection mode, or a setback behaviour matching the value table (eco, night).
02|000F|circuit_heating_curve_factor|false|
02|0010|heating_circuit_flow_temperature_setpoint_maximum|false|
02|0011|cooling_circuit_flow_temperature_setpoint_minimum|false|
02|0012|circuit_flow_temperature_setpoint_minimum|false|
02|0013|dhw_circuit_operating_mode|true|Candidate interpretation, unconfirmed on the wire: an enumeration of the DHW circuit operation; the value table is not established.
02|0014|heating_circuit_outdoor_temperature_limit_maximum|false|
02|0015|circuit_room_temperature_influence_mode|false|
02|0016|circuit_screed_drying_time|true|Two readings, unresolved (candidate interpretations, unconfirmed on the wire): the elapsed duration of the screed drying in days, or the index of the current day in the screed drying programme.
02|0017|screed_drying_circuit_temperature_setpoint|false|
02|0018|dhw_circuit_circulation_pump_active|true|Two readings, unresolved (candidate interpretations, unconfirmed on the wire): the state of a DHW circulation pump, or that external DHW is active.
02|0019|circuit_external_demand_input|true|Candidate interpretation, unconfirmed on the wire: an external heat demand input of the circuit; the value class and the value table are not established.
02|001A|circuit_mixer_movement|false|
02|001B|circuit_pump_unknown_8|true|A raw two-byte value related to the circuit pump; value class and meaning are not established.
02|001C|circuit_adaptive_heating_curve_offset|false|
02|001D|circuit_unknown_9|true|Candidate interpretations, unconfirmed on the wire: a DHW quick mode, or a frost protection threshold; the catalog decodes a float32 in °C; unresolved.
02|001E|circuit_state|true|Candidate interpretation, unconfirmed on the wire: an enumeration of the circuit state (0 standby, 1 heating, 2 cooling, 3 DHW); the assignment rests on observed naming and the hardware confirmation is pending.
02|001F|cooling_circuit_outdoor_temperature_limit_minimum|false|
02|0020|circuit_unknown_10|true|An observed one-byte raw value related to automatic heating and cooling; the meaning is not established.
02|0021|circuit_mixer_position|false|
02|0022|circuit_room_humidity|false|
02|0023|circuit_dew_point_temperature|false|
02|0024|circuit_pump_operating_time_total|false|
02|0025|circuit_pump_start_count_total|false|
03|0001|cooling_zone_operating_mode|false|
03|0002|cooling_zone_room_temperature_setpoint|false|
03|0003|zone_holiday_start_date|false|
03|0004|zone_holiday_end_date|false|
03|0005|zone_room_temperature_setpoint_at_holiday|false|
03|0006|heating_zone_operating_mode|false|
03|0007|heating_zone_room_comfort_temperature_setpoint|true|Representation is not established.
03|0008|zone_room_temperature_setpoint_at_quick_veto|false|
03|0009|heating_zone_room_setback_temperature_setpoint|false|
03|000A|heating_zone_operating_mode_at_teleswitch|true|Representation is not established.
03|000B|cooling_zone_operating_mode_at_teleswitch|true|Representation is not established.
03|000C|zone_bank_holiday_start_date|false|
03|000D|zone_bank_holiday_end_date|false|
03|000E|zone_special_function_mode|false|
03|000F|zone_room_temperature|false|
03|0010|zone_time_program_reset_command|false|
03|0011|zone_enabled|false|
03|0012|zone_valve_position|false|
03|0013|zone_room_sensor_binding_mode|false|
03|0014|heating_zone_room_effective_temperature_setpoint|false|
03|0015|cooling_zone_room_effective_temperature_setpoint|false|
03|0016|zone_name|false|
03|0017|zone_name_prefix|false|
03|0018|zone_name_suffix|false|
03|0019|heating_zone_time_program_active|false|
03|001A|cooling_zone_time_program_active|false|
03|001B|zone_special_function_state|false|
03|001C|zone_circuit_assignment|false|
03|001D|heating_zone_event_end_clock_time|true|Representation is not established.
03|001E|zone_quick_veto_end_clock_time|false|
03|0020|zone_holiday_end_clock_time|false|
03|0021|zone_holiday_start_clock_time|false|
03|0022|heating_zone_room_manual_temperature_setpoint|false|
03|0023|cooling_zone_room_manual_temperature_setpoint|false|
03|0024|zone_quick_veto_end_date|false|
03|0025|zone_holiday_abort_command|true|Candidate interpretation, unconfirmed on the wire: a command that aborts the zone holiday; the representation, the applicability and the writability are not established.
03|0026|zone_quick_veto_time|false|
03|0027|zone_manual_cooling_active|false|
03|0028|zone_room_humidity|false|
03|0029|zone_cooling_enabled|true|Candidate interpretation, unconfirmed on the wire: 1 means that cooling of the zone is enabled; no value table is documented and the polarity is not established.
03|002A|zone_summer_cutoff_active|true|Candidate interpretation, unconfirmed on the wire: 1 means that the summer cutoff of the zone is active; no value table is documented and the polarity is not established.
03|002B|zone_continuous_heating_active|true|Candidate interpretation, unconfirmed on the wire: 1 means that continuous heating of the zone is active; no value table is documented and the polarity is not established.
03|002C|zone_frost_protection_active|false|
03|002D|heating_zone_thermostat_cutoff_active|true|Candidate interpretation, unconfirmed on the wire: 1 means that the room thermostat function has switched heating off; the polarity is not established and no value table is documented.
03|002E|cooling_zone_thermostat_cutoff_active|true|Candidate interpretation, unconfirmed on the wire: 1 means that the room thermostat function has switched cooling off; the polarity is not established and no value table is documented.
04|0001|solar_pump_operating_time_total_reset_command|false|
04|0002|solar_pump_kick_enabled|false|
04|0003|solar_collector_temperature|false|
04|0004|solar_collector_temperature_limit_minimum|false|
04|0005|solar_circuit_protection_temperature_limit|false|
04|0006|solar_circuit_protection_stop_temperature_limit|false|
04|0007|solar_yield_temperature|false|
04|0008|solar_pump_active|false|
04|0009|solar_circuit_flow_rate|false|
04|000A|solar_circuit_flow_rate_setpoint|true|Representation is not established.
04|000B|solar_pump_operating_time_total|false|
04|000C|solar_pump_power_limit_minimum|true|Representation is not established.
05|0001|solar_cylinder_temperature_limit_maximum|false|
05|0002|solar_cylinder_start_temperature_difference_limit|false|
05|0003|solar_cylinder_stop_temperature_difference_limit|false|
05|0004|solar_cylinder_bottom_temperature|false|
08|0001|differential_control_temperature_limit_maximum|false|
08|0002|differential_control_temperature_limit_minimum|false|
08|0003|differential_control_start_temperature_difference_limit|false|
08|0004|differential_control_stop_temperature_difference_limit|false|
08|0005|differential_control_sensor_1_temperature|false|
08|0006|differential_control_sensor_2_temperature|false|
08|0007|differential_control_output|true|Value class is not established.
09|0001|ventilation_special_function_mode|false|
09|0002|ventilation_operating_mode|false|
09|0003|ventilation_time_program_active|false|
09|0004|ventilation_special_function_state|true|Candidate interpretation, unconfirmed on the wire: an enumeration of the ventilation special function (0 regular, 1 boost, 7 holiday, 10 system off); the presentation values are not independently verified and the relation to ` + "`" + `OP=0x02, GG=0x00, RR=0x00A4` + "`" + ` is not established.
09|0005|ventilation_volatile_organic_compound_sensor_1|true|Value class and unit are not established.
09|0006|ventilation_volatile_organic_compound_sensor_2|true|Value class and unit are not established.
09|0007|ventilation_holiday_end_date|false|
09|0008|ventilation_holiday_end_clock_time|false|
09|0009|ventilation_holiday_start_date|false|
09|000A|ventilation_holiday_start_clock_time|false|
09|000B|ventilation_heat_recovery|true|Value class is not established.
09|000C|ventilation_volatile_organic_compound_limit_maximum|true|Unit is not established.
09|000D|ventilation_fan_speed_limit_maximum_at_day|false|
09|000E|ventilation_fan_speed_limit_maximum_at_night|false|
09|000F|ventilation_holiday_abort_command|false|
`

func b524NameAnnotationPointer[b524NameAnnotationT any](b524NameAnnotationValue b524NameAnnotationT) *b524NameAnnotationT {
	return &b524NameAnnotationValue
}

func b524NameAnnotationExpectedRow(b524NameAnnotationTest *b524NameAnnotationTesting.T, b524NameAnnotationLine string) map[string]any {
	b524NameAnnotationTest.Helper()
	b524NameAnnotationParts := b524NameAnnotationStrings.Split(b524NameAnnotationLine, "|")
	b524NameAnnotationOperation := uint8(2)
	b524NameAnnotationGroupText, b524NameAnnotationRegisterText := b524NameAnnotationParts[0], b524NameAnnotationParts[1]
	if len(b524NameAnnotationParts) == 4 {
		b524NameAnnotationOperation = 6
		b524NameAnnotationGroupText, b524NameAnnotationRegisterText = b524NameAnnotationParts[1], b524NameAnnotationParts[2]
	}
	b524NameAnnotationGroup, b524NameAnnotationError := b524NameAnnotationStrconv.ParseUint(b524NameAnnotationGroupText, 16, 8)
	if b524NameAnnotationError != nil {
		b524NameAnnotationTest.Fatal(b524NameAnnotationError)
	}
	b524NameAnnotationRegister, b524NameAnnotationError := b524NameAnnotationStrconv.ParseUint(b524NameAnnotationRegisterText, 16, 16)
	if b524NameAnnotationError != nil {
		b524NameAnnotationTest.Fatal(b524NameAnnotationError)
	}
	b524NameAnnotationRow := map[string]any{
		"key": B524NameAnnotationKey{
			B524NameAnnotationFamily: "unqualified", B524NameAnnotationVariant: "unqualified", B524NameAnnotationPrimary: 181,
			B524NameAnnotationSecondary: b524NameAnnotationPointer(uint8(36)), B524NameAnnotationOperation: b524NameAnnotationPointer(b524NameAnnotationOperation),
			B524NameAnnotationDirection: "none", B524NameAnnotationField: "name", B524NameAnnotationAlternative: "none",
			B524NameAnnotationGroup: b524NameAnnotationPointer(uint8(b524NameAnnotationGroup)), B524NameAnnotationInstanceScope: "unqualified",
			B524NameAnnotationRegister: b524NameAnnotationPointer(uint16(b524NameAnnotationRegister)),
		},
		"canonical_name": nil, "unresolved": nil,
		"representation_grade": "Proven", "grade": "Hypothesis", "semantics_grade": "Unknown",
		"physical_qualification": "Unknown",
	}
	if b524NameAnnotationOperation == 2 {
		b524NameAnnotationRow["canonical_name"] = b524NameAnnotationParts[2]
		b524NameAnnotationRow["unresolved"] = b524NameAnnotationParts[3] == "true"
		b524NameAnnotationRow["open_question"] = b524NameAnnotationParts[4]
	} else {
		b524NameAnnotationRow["native_label"] = b524NameAnnotationParts[3]
		b524NameAnnotationRow["raw_form"] = b524NameAnnotationLine
		b524NameAnnotationRow["open_question"] = "No matching public OP06 canonical-name contract; native label is evidence only."
	}
	return b524NameAnnotationRow
}

func b524NameAnnotationEncode(b524NameAnnotationTest *b524NameAnnotationTesting.T, b524NameAnnotationValue any) []byte {
	b524NameAnnotationTest.Helper()
	b524NameAnnotationEncoded, b524NameAnnotationError := b524NameAnnotationJSON.Marshal(b524NameAnnotationValue)
	if b524NameAnnotationError != nil {
		b524NameAnnotationTest.Fatal(b524NameAnnotationError)
	}
	return b524NameAnnotationEncoded
}

func b524NameAnnotationEqualJSON(b524NameAnnotationTest *b524NameAnnotationTesting.T, b524NameAnnotationActual, b524NameAnnotationExpectedJSON []byte) {
	b524NameAnnotationTest.Helper()
	var b524NameAnnotationActualValue, b524NameAnnotationExpectedValue any
	if b524NameAnnotationError := b524NameAnnotationJSON.Unmarshal(b524NameAnnotationActual, &b524NameAnnotationActualValue); b524NameAnnotationError != nil {
		b524NameAnnotationTest.Fatal(b524NameAnnotationError)
	}
	if b524NameAnnotationError := b524NameAnnotationJSON.Unmarshal(b524NameAnnotationExpectedJSON, &b524NameAnnotationExpectedValue); b524NameAnnotationError != nil {
		b524NameAnnotationTest.Fatal(b524NameAnnotationError)
	}
	if !b524NameAnnotationReflect.DeepEqual(b524NameAnnotationActualValue, b524NameAnnotationExpectedValue) {
		b524NameAnnotationTest.Fatalf("ab36 whole row mismatch:\n%s\nwant:\n%s", b524NameAnnotationActual, b524NameAnnotationExpectedJSON)
	}
}

func TestAb36Correspondence(b524NameAnnotationTest *b524NameAnnotationTesting.T) {
	b524NameAnnotationEntries, b524NameAnnotationError := B524NameAnnotationEntries()
	if b524NameAnnotationError != nil {
		b524NameAnnotationTest.Fatal(b524NameAnnotationError)
	}
	var b524NameAnnotationRows []b524NameAnnotationJSON.RawMessage
	if b524NameAnnotationError := b524NameAnnotationJSON.Unmarshal(b524NameAnnotationEntries, &b524NameAnnotationRows); b524NameAnnotationError != nil {
		b524NameAnnotationTest.Fatal(b524NameAnnotationError)
	}
	b524NameAnnotationLines := b524NameAnnotationStrings.Split(b524NameAnnotationStrings.TrimSpace(b524NameAnnotationExpected), "\n")
	if len(b524NameAnnotationLines) != 368 || len(b524NameAnnotationRows) != 370 {
		b524NameAnnotationTest.Fatalf("ab36 counts: %d / %d", len(b524NameAnnotationLines), len(b524NameAnnotationRows))
	}
	b524NameAnnotationLines = append(b524NameAnnotationLines, "06|0A|0007|RoomHumidity", "06|0A|000F|RoomTemp")
	emittedRows := make(map[string]b524NameAnnotationJSON.RawMessage)
	for _, emittedRow := range b524NameAnnotationRows {
		var envelope struct {
			Key B524NameAnnotationKey `json:"key"`
		}
		if err := b524NameAnnotationJSON.Unmarshal(emittedRow, &envelope); err != nil {
			b524NameAnnotationTest.Fatal(err)
		}
		identity := string(b524NameAnnotationEncode(b524NameAnnotationTest, envelope.Key))
		if _, exists := emittedRows[identity]; exists {
			b524NameAnnotationTest.Fatal("duplicate emitted full identity")
		}
		emittedRows[identity] = emittedRow
	}
	b524NameAnnotationSeen := make(map[string]bool)
	provisionalCount := 0
	for _, b524NameAnnotationLine := range b524NameAnnotationLines {
		b524NameAnnotationExpectedRowValue := b524NameAnnotationExpectedRow(b524NameAnnotationTest, b524NameAnnotationLine)
		b524NameAnnotationKey := b524NameAnnotationExpectedRowValue["key"].(B524NameAnnotationKey)
		b524NameAnnotationKeyBytes := b524NameAnnotationEncode(b524NameAnnotationTest, b524NameAnnotationKey)
		if b524NameAnnotationSeen[string(b524NameAnnotationKeyBytes)] {
			b524NameAnnotationTest.Fatal("ab36 duplicate expected key")
		}
		b524NameAnnotationSeen[string(b524NameAnnotationKeyBytes)] = true
		b524NameAnnotationWant := b524NameAnnotationEncode(b524NameAnnotationTest, b524NameAnnotationExpectedRowValue)
		b524NameAnnotationEqualJSON(b524NameAnnotationTest, emittedRows[string(b524NameAnnotationKeyBytes)], b524NameAnnotationWant)
		b524NameAnnotationActual, b524NameAnnotationFound, b524NameAnnotationError := B524NameAnnotationLookup(b524NameAnnotationKey)
		if b524NameAnnotationError != nil || !b524NameAnnotationFound {
			b524NameAnnotationTest.Fatalf("ab36 missing %s: %v", b524NameAnnotationLine, b524NameAnnotationError)
		}
		b524NameAnnotationEqualJSON(b524NameAnnotationTest, b524NameAnnotationActual, b524NameAnnotationWant)
		if *b524NameAnnotationKey.B524NameAnnotationOperation == 2 && b524NameAnnotationExpectedRowValue["unresolved"] == true {
			provisionalCount++
		}
	}
	if len(b524NameAnnotationSeen) != 370 || provisionalCount != 142 {
		b524NameAnnotationTest.Fatalf("ab36 unique/provisional: %d / %d", len(b524NameAnnotationSeen), provisionalCount)
	}
	b524NameAnnotationTest.Logf("ab36 verified 368 OP02 keys, 142 provisional questions, 2 separate OP06 native labels; all whole rows")
}

func TestAb36ExactAxes(b524NameAnnotationTest *b524NameAnnotationTesting.T) {
	b524NameAnnotationBase := b524NameAnnotationExpectedRow(b524NameAnnotationTest, "06|0A|000F|RoomTemp")["key"].(B524NameAnnotationKey)
	b524NameAnnotationControls, b524NameAnnotationFound, b524NameAnnotationError := B524NameAnnotationLookup(b524NameAnnotationBase)
	if b524NameAnnotationError != nil || !b524NameAnnotationFound || len(b524NameAnnotationControls) == 0 {
		b524NameAnnotationTest.Fatal("ab36 missing positive control")
	}
	for b524NameAnnotationField := 0; b524NameAnnotationField < b524NameAnnotationReflect.TypeOf(b524NameAnnotationBase).NumField(); b524NameAnnotationField++ {
		b524NameAnnotationName := b524NameAnnotationReflect.TypeOf(b524NameAnnotationBase).Field(b524NameAnnotationField).Name
		b524NameAnnotationTest.Run(b524NameAnnotationName, func(b524NameAnnotationTest *b524NameAnnotationTesting.T) {
			b524NameAnnotationQuery := b524NameAnnotationBase
			b524NameAnnotationValue := b524NameAnnotationReflect.ValueOf(&b524NameAnnotationQuery).Elem().Field(b524NameAnnotationField)
			switch b524NameAnnotationValue.Kind() {
			case b524NameAnnotationReflect.String:
				b524NameAnnotationValue.SetString("ab36_other")
			case b524NameAnnotationReflect.Uint8:
				b524NameAnnotationValue.SetUint(254)
			case b524NameAnnotationReflect.Ptr:
				b524NameAnnotationReplacement := b524NameAnnotationReflect.New(b524NameAnnotationValue.Type().Elem())
				b524NameAnnotationReplacement.Elem().SetUint(254)
				b524NameAnnotationValue.Set(b524NameAnnotationReplacement)
			default:
				b524NameAnnotationTest.Fatal("ab36 untested axis")
			}
			b524NameAnnotationActual, b524NameAnnotationFound, b524NameAnnotationError := B524NameAnnotationLookup(b524NameAnnotationQuery)
			if b524NameAnnotationError != nil || b524NameAnnotationFound || b524NameAnnotationActual != nil {
				b524NameAnnotationTest.Fatalf("ab36 fallback on %s: %s %v", b524NameAnnotationName, b524NameAnnotationActual, b524NameAnnotationError)
			}
			if b524NameAnnotationValue.Kind() == b524NameAnnotationReflect.Ptr {
				for _, b524NameAnnotationZero := range []bool{false, true} {
					if b524NameAnnotationZero {
						b524NameAnnotationValue.Set(b524NameAnnotationReflect.New(b524NameAnnotationValue.Type().Elem()))
					} else {
						b524NameAnnotationValue.SetZero()
					}
					b524NameAnnotationActual, b524NameAnnotationFound, b524NameAnnotationError := B524NameAnnotationLookup(b524NameAnnotationQuery)
					if b524NameAnnotationError != nil || b524NameAnnotationFound || b524NameAnnotationActual != nil {
						b524NameAnnotationTest.Fatalf("ab36 absent/zero fallback on %s", b524NameAnnotationName)
					}
				}
			}
		})
	}
	b524NameAnnotationQuery := b524NameAnnotationBase
	b524NameAnnotationQuery.B524NameAnnotationOperation = b524NameAnnotationPointer(uint8(2))
	b524NameAnnotationActual, b524NameAnnotationFound, b524NameAnnotationError := B524NameAnnotationLookup(b524NameAnnotationQuery)
	if b524NameAnnotationError != nil || b524NameAnnotationFound || b524NameAnnotationActual != nil {
		b524NameAnnotationTest.Fatal("ab36 OP06 merged into OP02")
	}
}

func TestAb36IndependentStorage(b524NameAnnotationTest *b524NameAnnotationTesting.T) {
	b524NameAnnotationKey := b524NameAnnotationExpectedRow(b524NameAnnotationTest, "06|0A|0007|RoomHumidity")["key"].(B524NameAnnotationKey)
	b524NameAnnotationRow, b524NameAnnotationFound, b524NameAnnotationError := B524NameAnnotationLookup(b524NameAnnotationKey)
	if b524NameAnnotationError != nil || !b524NameAnnotationFound {
		b524NameAnnotationTest.Fatal("ab36 missing row")
	}
	b524NameAnnotationWant := append([]byte(nil), b524NameAnnotationRow...)
	b524NameAnnotationEntries, b524NameAnnotationError := B524NameAnnotationEntries()
	if b524NameAnnotationError != nil {
		b524NameAnnotationTest.Fatal(b524NameAnnotationError)
	}
	b524NameAnnotationWantEntries := append([]byte(nil), b524NameAnnotationEntries...)
	for b524NameAnnotationPosition := range b524NameAnnotationRow {
		b524NameAnnotationRow[b524NameAnnotationPosition] = 'x'
	}
	for b524NameAnnotationPosition := range b524NameAnnotationEntries {
		b524NameAnnotationEntries[b524NameAnnotationPosition] = 'x'
	}
	b524NameAnnotationAgain, b524NameAnnotationFound, b524NameAnnotationError := B524NameAnnotationLookup(b524NameAnnotationKey)
	if b524NameAnnotationError != nil || !b524NameAnnotationFound || !b524NameAnnotationBytes.Equal(b524NameAnnotationAgain, b524NameAnnotationWant) {
		b524NameAnnotationTest.Fatal("ab36 row storage shared")
	}
	b524NameAnnotationAgainEntries, b524NameAnnotationError := B524NameAnnotationEntries()
	if b524NameAnnotationError != nil || !b524NameAnnotationBytes.Equal(b524NameAnnotationAgainEntries, b524NameAnnotationWantEntries) {
		b524NameAnnotationTest.Fatal("ab36 entries storage shared")
	}
	b524NameAnnotationSource := append([]byte(nil), b524NameAnnotationWantEntries...)
	b524NameAnnotationRows, b524NameAnnotationIndex, b524NameAnnotationError := b524NameAnnotationValidate(b524NameAnnotationSource)
	if b524NameAnnotationError != nil {
		b524NameAnnotationTest.Fatal(b524NameAnnotationError)
	}
	for b524NameAnnotationPosition := range b524NameAnnotationSource {
		b524NameAnnotationSource[b524NameAnnotationPosition] = 'x'
	}
	b524NameAnnotationEqualJSON(b524NameAnnotationTest, b524NameAnnotationRows[368], b524NameAnnotationWant)
	b524NameAnnotationEqualJSON(b524NameAnnotationTest, b524NameAnnotationIndex[string(b524NameAnnotationEncode(b524NameAnnotationTest, b524NameAnnotationKey))], b524NameAnnotationWant)
	b524NameAnnotationAgain, b524NameAnnotationFound, b524NameAnnotationError = B524NameAnnotationLookup(b524NameAnnotationKey)
	if b524NameAnnotationError != nil || !b524NameAnnotationFound || !b524NameAnnotationBytes.Equal(b524NameAnnotationAgain, b524NameAnnotationWant) {
		b524NameAnnotationTest.Fatal("ab36 source mutation affected API")
	}
}

func TestAb36InvalidData(b524NameAnnotationTest *b524NameAnnotationTesting.T) {
	b524NameAnnotationRow := b524NameAnnotationExpectedRow(b524NameAnnotationTest, "06|0A|0007|RoomHumidity")
	b524NameAnnotationValid := string(b524NameAnnotationEncode(b524NameAnnotationTest, []any{b524NameAnnotationRow}))
	b524NameAnnotationDuplicate := string(b524NameAnnotationEncode(b524NameAnnotationTest, []any{b524NameAnnotationRow, b524NameAnnotationRow}))
	b524NameAnnotationOtherRow := b524NameAnnotationExpectedRow(b524NameAnnotationTest, "06|0A|0007|RoomHumidity")
	b524NameAnnotationOtherRow["grade"] = "Unknown"
	b524NameAnnotationOtherRow["ab36_width"] = 2
	b524NameAnnotationDifferentDuplicate := string(b524NameAnnotationEncode(b524NameAnnotationTest, []any{b524NameAnnotationRow, b524NameAnnotationOtherRow}))
	b524NameAnnotationOriginal := b524NameAnnotationData
	b524NameAnnotationOriginalCatalog := b524NameAnnotationCachedCatalog
	b524NameAnnotationTest.Cleanup(func() {
		b524NameAnnotationData = b524NameAnnotationOriginal
		b524NameAnnotationCachedCatalog = b524NameAnnotationOriginalCatalog
	})
	for b524NameAnnotationPosition, b524NameAnnotationBad := range []string{
		"", "[", "null", "{}", "[null]", "[42]", "[{}]", "[{\"key\":null}]", b524NameAnnotationDuplicate, b524NameAnnotationDifferentDuplicate,
		b524NameAnnotationStrings.Replace(b524NameAnnotationValid, `"primary":181,`, ``, 1),
		b524NameAnnotationStrings.Replace(b524NameAnnotationValid, `"primary":181`, `"primary":null`, 1),
		b524NameAnnotationStrings.Replace(b524NameAnnotationValid, `"primary":181`, `"primary":256`, 1),
		b524NameAnnotationStrings.Replace(b524NameAnnotationValid, `"register":7`, `"register":65536`, 1),
		b524NameAnnotationStrings.Replace(b524NameAnnotationValid, `"secondary":36`, `"secondary":-1`, 1),
		b524NameAnnotationStrings.Replace(b524NameAnnotationValid, `"operation":6`, `"operation":"6"`, 1),
		b524NameAnnotationStrings.Replace(b524NameAnnotationValid, `"family":"unqualified"`, `"family":""`, 1),
		b524NameAnnotationStrings.Replace(b524NameAnnotationValid, `"variant":"unqualified",`, ``, 1),
		b524NameAnnotationStrings.Replace(b524NameAnnotationValid, `"primary":181`, `"ab36_unknown":181`, 1),
	} {
		b524NameAnnotationTest.Run("ab36_invalid_"+b524NameAnnotationStrconv.Itoa(b524NameAnnotationPosition), func(b524NameAnnotationTest *b524NameAnnotationTesting.T) {
			b524NameAnnotationRows, b524NameAnnotationIndex, b524NameAnnotationError := b524NameAnnotationValidate([]byte(b524NameAnnotationBad))
			if b524NameAnnotationError == nil || b524NameAnnotationRows != nil || b524NameAnnotationIndex != nil {
				b524NameAnnotationTest.Fatal("ab36 invalid metadata accepted")
			}
			b524NameAnnotationData = b524NameAnnotationBad
			b524NameAnnotationCachedCatalog = &b524NameAnnotationCatalog{b524NameAnnotationSource: b524NameAnnotationData}
			b524NameAnnotationEntries, b524NameAnnotationError := B524NameAnnotationEntries()
			if b524NameAnnotationError == nil || b524NameAnnotationEntries != nil {
				b524NameAnnotationTest.Fatal("ab36 Entries swallowed invalid metadata")
			}
			b524NameAnnotationActual, b524NameAnnotationFound, b524NameAnnotationError := B524NameAnnotationLookup(b524NameAnnotationRow["key"].(B524NameAnnotationKey))
			if b524NameAnnotationError == nil || b524NameAnnotationFound || b524NameAnnotationActual != nil {
				b524NameAnnotationTest.Fatal("ab36 Lookup swallowed invalid metadata")
			}
		})
	}
	b524NameAnnotationData = "[]"
	b524NameAnnotationCachedCatalog = &b524NameAnnotationCatalog{b524NameAnnotationSource: b524NameAnnotationData}
	b524NameAnnotationEntries, b524NameAnnotationError := B524NameAnnotationEntries()
	if b524NameAnnotationError != nil || string(b524NameAnnotationEntries) != "[]" {
		b524NameAnnotationTest.Fatal("ab36 valid empty array rejected")
	}
}

func TestAb36OptionalAxesAndRawMetadata(b524NameAnnotationTest *b524NameAnnotationTesting.T) {
	b524NameAnnotationOriginal := b524NameAnnotationData
	b524NameAnnotationOriginalCatalog := b524NameAnnotationCachedCatalog
	b524NameAnnotationTest.Cleanup(func() {
		b524NameAnnotationData = b524NameAnnotationOriginal
		b524NameAnnotationCachedCatalog = b524NameAnnotationOriginalCatalog
	})
	b524NameAnnotationBase := b524NameAnnotationExpectedRow(b524NameAnnotationTest, "06|0A|0007|RoomHumidity")
	for _, b524NameAnnotationField := range []string{"B524NameAnnotationSecondary", "B524NameAnnotationOperation", "B524NameAnnotationGroup", "B524NameAnnotationRegister"} {
		b524NameAnnotationTest.Run("ab36_"+b524NameAnnotationField, func(b524NameAnnotationTest *b524NameAnnotationTesting.T) {
			b524NameAnnotationAbsent := b524NameAnnotationBase["key"].(B524NameAnnotationKey)
			b524NameAnnotationReflect.ValueOf(&b524NameAnnotationAbsent).Elem().FieldByName(b524NameAnnotationField).SetZero()
			b524NameAnnotationZero := b524NameAnnotationAbsent
			b524NameAnnotationValue := b524NameAnnotationReflect.ValueOf(&b524NameAnnotationZero).Elem().FieldByName(b524NameAnnotationField)
			b524NameAnnotationValue.Set(b524NameAnnotationReflect.New(b524NameAnnotationValue.Type().Elem()))
			b524NameAnnotationRows := []map[string]any{}
			for b524NameAnnotationPosition, b524NameAnnotationKey := range []B524NameAnnotationKey{b524NameAnnotationAbsent, b524NameAnnotationZero} {
				b524NameAnnotationRow := b524NameAnnotationExpectedRow(b524NameAnnotationTest, "06|0A|0007|RoomHumidity")
				b524NameAnnotationRow["key"] = b524NameAnnotationKey
				b524NameAnnotationRow["raw_form"] = "ab36_synthetic_" + b524NameAnnotationStrconv.Itoa(b524NameAnnotationPosition)
				b524NameAnnotationRow["ab36_extra"] = map[string]any{"ab36_bytes": []int{0, 255}, "ab36_state": "Unknown"}
				b524NameAnnotationRows = append(b524NameAnnotationRows, b524NameAnnotationRow)
			}
			b524NameAnnotationData = string(b524NameAnnotationEncode(b524NameAnnotationTest, b524NameAnnotationRows))
			b524NameAnnotationCachedCatalog = &b524NameAnnotationCatalog{b524NameAnnotationSource: b524NameAnnotationData}
			for _, b524NameAnnotationRow := range b524NameAnnotationRows {
				b524NameAnnotationActual, b524NameAnnotationFound, b524NameAnnotationError := B524NameAnnotationLookup(b524NameAnnotationRow["key"].(B524NameAnnotationKey))
				if b524NameAnnotationError != nil || !b524NameAnnotationFound {
					b524NameAnnotationTest.Fatal("ab36 absent/zero collapsed")
				}
				b524NameAnnotationEqualJSON(b524NameAnnotationTest, b524NameAnnotationActual, b524NameAnnotationEncode(b524NameAnnotationTest, b524NameAnnotationRow))
			}
		})
	}
}

func TestAb36OmittedOptionalKeys(b524NameAnnotationTest *b524NameAnnotationTesting.T) {
	b524NameAnnotationOriginal := b524NameAnnotationData
	b524NameAnnotationOriginalCatalog := b524NameAnnotationCachedCatalog
	b524NameAnnotationTest.Cleanup(func() {
		b524NameAnnotationData = b524NameAnnotationOriginal
		b524NameAnnotationCachedCatalog = b524NameAnnotationOriginalCatalog
	})
	for _, b524NameAnnotationField := range []string{"B524NameAnnotationSecondary", "B524NameAnnotationOperation", "B524NameAnnotationGroup", "B524NameAnnotationRegister"} {
		b524NameAnnotationTest.Run("ab36_"+b524NameAnnotationField, func(b524NameAnnotationTest *b524NameAnnotationTesting.T) {
			b524NameAnnotationRow := b524NameAnnotationExpectedRow(b524NameAnnotationTest, "06|0A|0007|RoomHumidity")
			b524NameAnnotationKey := b524NameAnnotationRow["key"].(B524NameAnnotationKey)
			b524NameAnnotationValue := b524NameAnnotationReflect.ValueOf(&b524NameAnnotationKey).Elem().FieldByName(b524NameAnnotationField)
			b524NameAnnotationValue.SetZero()
			var b524NameAnnotationKeyObject map[string]any
			if b524NameAnnotationError := b524NameAnnotationJSON.Unmarshal(b524NameAnnotationEncode(b524NameAnnotationTest, b524NameAnnotationKey), &b524NameAnnotationKeyObject); b524NameAnnotationError != nil {
				b524NameAnnotationTest.Fatal(b524NameAnnotationError)
			}
			b524NameAnnotationDefinition, b524NameAnnotationPresent := b524NameAnnotationReflect.TypeOf(b524NameAnnotationKey).FieldByName(b524NameAnnotationField)
			if !b524NameAnnotationPresent {
				b524NameAnnotationTest.Fatal("ab36 missing axis")
			}
			delete(b524NameAnnotationKeyObject, b524NameAnnotationDefinition.Tag.Get("json"))
			b524NameAnnotationRow["key"] = b524NameAnnotationKeyObject
			b524NameAnnotationData = string(b524NameAnnotationEncode(b524NameAnnotationTest, []any{b524NameAnnotationRow}))
			b524NameAnnotationCachedCatalog = &b524NameAnnotationCatalog{b524NameAnnotationSource: b524NameAnnotationData}
			b524NameAnnotationActual, b524NameAnnotationFound, b524NameAnnotationError := B524NameAnnotationLookup(b524NameAnnotationKey)
			if b524NameAnnotationError != nil || !b524NameAnnotationFound {
				b524NameAnnotationTest.Fatalf("ab36 omitted optional key rejected: %v", b524NameAnnotationError)
			}
			b524NameAnnotationEqualJSON(b524NameAnnotationTest, b524NameAnnotationActual, b524NameAnnotationEncode(b524NameAnnotationTest, b524NameAnnotationRow))
			b524NameAnnotationValue.Set(b524NameAnnotationReflect.New(b524NameAnnotationValue.Type().Elem()))
			b524NameAnnotationActual, b524NameAnnotationFound, b524NameAnnotationError = B524NameAnnotationLookup(b524NameAnnotationKey)
			if b524NameAnnotationError != nil || b524NameAnnotationFound || b524NameAnnotationActual != nil {
				b524NameAnnotationTest.Fatal("ab36 omitted key matched explicit zero")
			}
		})
	}
}

func TestAb36LookupCachedIndex(b524NameAnnotationTest *b524NameAnnotationTesting.T) {
	b524NameAnnotationKey := b524NameAnnotationExpectedRow(b524NameAnnotationTest, "06|0A|0007|RoomHumidity")["key"].(B524NameAnnotationKey)
	// AllocsPerRun warms the cache before measuring repeated public lookups.
	b524NameAnnotationAllocations := b524NameAnnotationTesting.AllocsPerRun(10, func() {
		b524NameAnnotationRow, b524NameAnnotationFound, b524NameAnnotationError := B524NameAnnotationLookup(b524NameAnnotationKey)
		if b524NameAnnotationError != nil || !b524NameAnnotationFound || len(b524NameAnnotationRow) == 0 {
			b524NameAnnotationTest.Fatal("ab36 cached lookup failed")
		}
	})
	// Allow ample room for key encoding and the independent returned row,
	// while rejecting full-catalog parsing and indexing on every lookup.
	if b524NameAnnotationAllocations > 100 {
		b524NameAnnotationTest.Fatalf("ab36 repeated lookup allocated %.0f objects; catalog index rebuilt", b524NameAnnotationAllocations)
	}
}

func TestAb36CatalogValidationOnce(b524NameAnnotationTest *b524NameAnnotationTesting.T) {
	b524NameAnnotationKey := b524NameAnnotationExpectedRow(b524NameAnnotationTest, "06|0A|0007|RoomHumidity")["key"].(B524NameAnnotationKey)
	b524NameAnnotationEncoded := string(b524NameAnnotationEncode(b524NameAnnotationTest, b524NameAnnotationKey))
	b524NameAnnotationCatalogValue := &b524NameAnnotationCatalog{b524NameAnnotationSource: b524NameAnnotationData}
	b524NameAnnotationIndex, b524NameAnnotationError := b524NameAnnotationCatalogValue.b524NameAnnotationLoad()
	if b524NameAnnotationError != nil || len(b524NameAnnotationIndex) != 370 {
		b524NameAnnotationTest.Fatalf("ab36 initial validation failed: %v", b524NameAnnotationError)
	}
	b524NameAnnotationWant := append([]byte(nil), b524NameAnnotationIndex[b524NameAnnotationEncoded]...)
	// Changing the private source after initialization must not trigger validation again.
	b524NameAnnotationCatalogValue.b524NameAnnotationSource = "["
	b524NameAnnotationAgain, b524NameAnnotationError := b524NameAnnotationCatalogValue.b524NameAnnotationLoad()
	if b524NameAnnotationError != nil || !b524NameAnnotationBytes.Equal(b524NameAnnotationAgain[b524NameAnnotationEncoded], b524NameAnnotationWant) {
		b524NameAnnotationTest.Fatal("ab36 validated index was rebuilt")
	}
	b524NameAnnotationInvalid := &b524NameAnnotationCatalog{b524NameAnnotationSource: "["}
	b524NameAnnotationIndex, b524NameAnnotationFirstError := b524NameAnnotationInvalid.b524NameAnnotationLoad()
	if b524NameAnnotationFirstError == nil || b524NameAnnotationIndex != nil {
		b524NameAnnotationTest.Fatal("ab36 invalid source accepted")
	}
	b524NameAnnotationInvalid.b524NameAnnotationSource = b524NameAnnotationData
	b524NameAnnotationIndex, b524NameAnnotationError = b524NameAnnotationInvalid.b524NameAnnotationLoad()
	if b524NameAnnotationError != b524NameAnnotationFirstError || b524NameAnnotationIndex != nil {
		b524NameAnnotationTest.Fatal("ab36 validation error was not cached")
	}
}
