package vaillant

import (
	ab36bytes "bytes"
	ab36json "encoding/json"
	ab36reflect "reflect"
	ab36strconv "strconv"
	ab36strings "strings"
	ab36testing "testing"
)

const ab36Expected = `
00|0001|dhw_outdoor_bivalence_temperature_limit|system_dhw_bivalence_point|false
00|0002|continuous_heating_outdoor_temperature_limit|continuous_heating_limit|false
00|0003|frost_protection_delay_time|frost_protection_delay|false
00|0004|preheating_time_maximum|max_preheating_time|false
00|0006|manual_cooling_time|manual_cooling_days|false
00|0007|unknown_11|system_off|true
00|0008|cooling_enabled_at_holiday|cooling_while_holiday|true
00|0009|automatic_cooling_enabled|automatic_cooling_enabled|false
00|000A|cylinder_parallel_loading_enabled|system_cylinder_parallel_loading_allowed|false
00|000B|solar_circuit_bleeding_time|system_solar_circuit_bleeding_time|true
00|000D|solar_frost_protection_limit|system_solar_frost_protection_limit|true
00|000E|room_humidity_limit_maximum|max_room_humidity|false
00|000F|hybrid_control_mode|hybrid_control_strategy|false
00|0010|backup_heater_tariff|backup_heater_tariff|true
00|0011|electricity_price_at_low_tariff|energy_rate_low_tariff|true
00|0012|electricity_price_at_high_tariff|energy_rate_high_tariff|true
00|0013|backup_heater_type|backup_heater_type|true
00|0014|adaptive_heating_curve_enabled|adaptive_heating_curve|false
00|0015|automatic_daylight_saving_time|automatic_daylight_saving_time|true
00|0016|ventilation_operating_mode|system_ventilation_operating_mode|false
00|0017|dhw_loading_time_maximum|system_dhw_max_loading_time|false
00|0018|dhw_blocking_time|system_dhw_blocking_time|false
00|0019|solar_flow_rate_setpoint|system_solar_flow_rate_setpoint|false
00|001B|dhw_pump_overrun_time|pump_overrun_time|false
00|001C|heat_generator_flow_temperature_limit_maximum|generator_max_flow_temperature|false
00|001E|automatic_sequence_reversal|automatic_sequence_reversal|true
00|0022|heating_outdoor_alternative_temperature_limit|alternative_point_for_central_heating|false
00|0023|heating_outdoor_bivalence_temperature_limit|heating_circuit_bivalence_point|false
00|0024|backup_heater_mode|backup_heater_allowed_for|true
00|0025|temporary_backup_heater_mode|backup_heater_allow_temporary|true
00|0026|heating_flow_temperature_setpoint_maximum_at_heat_pump_fault|heating_circuit_emergency_temperature|false
00|0027|cylinder_loading_start_hysteresis|system_cylinder_start_loading_hysteresis|false
00|0029|cylinder_loading_offset|system_cylinder_loading_offset|false
00|002A|dhw_legionella_protection_start_clock_time|system_cylinder_dhw_legionella_protection_start_time|false
00|002B|dhw_legionella_protection_start_weekday|system_cylinder_dhw_legionella_protection_start_weekday|false
00|002C|maintenance_date|maintenance_date|false
00|002D|outdoor_temperature_offset|outside_temperature_correction|false
00|002F|functional_module_configuration_code|module_configuration_vr71|true
00|0031|language|language|true
00|0034|date|system_date|false
00|0035|clock_time|system_time|false
00|0036|hydraulic_scheme_code|system_scheme|false
00|0038|cooling_start_temperature_limit|system_cooling_start_temperature|false
00|0039|water_pressure|system_water_pressure|false
00|003A|hybrid_unknown_2|tri_vai_calculation_result|true
00|003C|solar_yield_total_reset_command|system_yield_solar_reset|false
00|003D|solar_yield_total|system_yield_solar_total|false
00|003E|environmental_yield_total|system_yield_environmental_total|false
00|003F|environmental_yield_total_reset_command|system_yield_environmental_reset|true
00|0040|ventilation_volatile_organic_compound_sensor_1|system_ventilation_voc_sensor_1|true
00|0041|ventilation_volatile_organic_compound_sensor_2|system_ventilation_voc_sensor_2|true
00|0042|ventilation_fan_speed_limit_maximum_at_day|system_ventilation_day_max_fan_stage|false
00|0043|ventilation_fan_speed_limit_maximum_at_night|system_ventilation_night_max_fan_stage|false
00|0044|ventilation_heat_recovery|system_ventilation_heat_recovery|true
00|0045|utility_lockout_mode|power_cut_function_evu|true
00|0046|dhw_flow_temperature_setpoint_maximum|system_cylinder_mss_dhw_max_flow_setpoint|false
00|0047|emergency_operation_active|system_in_failure_mode|true
00|0048|energy_manager_state|energy_manager_state|false
00|0049|electricity_consumption_total_reset_command|system_consumption_electricity_reset|true
00|004A|settings_reset_command|reset_to_default|true
00|004B|flow_temperature|system_flow_temperature|false
00|004C|backup_heater_power_limit_maximum|vwz_ai_backup_heater_max_power|true
00|004D|multifunction_output_2_mode|vwz_ai_configuration_ma2|false
00|004E|heating_fuel_consumption_current_month|system_consumption_fuel_heating_current_month|false
00|004F|heating_electricity_consumption_current_month|system_consumption_electricity_heating_current_month|false
00|0050|dhw_electricity_consumption_current_month|system_consumption_electricity_dhw_current_month|false
00|0051|dhw_fuel_consumption_current_month|system_consumption_fuel_dhw_current_month|false
00|0052|heating_fuel_consumption_previous_month|system_consumption_fuel_heating_previous_month|false
00|0053|heating_electricity_consumption_previous_month|system_consumption_electricity_heating_previous_month|false
00|0054|dhw_electricity_consumption_previous_month|system_consumption_electricity_dhw_previous_month|false
00|0055|dhw_fuel_consumption_previous_month|system_consumption_fuel_dhw_previous_month|false
00|0056|heating_fuel_consumption_total|system_consumption_fuel_heating_total|false
00|0057|heating_electricity_consumption_total|system_consumption_electricity_heating_total|false
00|0058|dhw_electricity_consumption_total|system_consumption_electricity_dhw_total|false
00|0059|dhw_fuel_consumption_total|system_consumption_fuel_dhw_total|false
00|005C|electricity_consumption_1_total|system_consumption_electricity_total|true
00|005D|fuel_consumption_1_total|system_consumption_fuel_total|true
00|005E|heat_recovery_yield_total|system_yield_heatrecovery_total|false
00|005F|actuator_sensor_test_device_code|actor_sensor_test_device|true
00|0060|actuator_sensor_test_sensor_code|actor_sensor_test_sensor|true
00|0061|actuator_sensor_test_actuator_code|actor_sensor_test_actor|true
00|0062|ventilation_volatile_organic_compound_sensor_3|system_ventilation_voc_sensor_3|true
00|0063|ventilation_volatile_organic_compound_limit_maximum|system_ventilation_max_voc|true
00|0065|ventilation_fan_speed_setpoint|system_ventilation_displayed_wish_fan_stage|false
00|0066|home_screen_unknown_6|home_screen_status|true
00|0067|cylinder_bottom_temperature|system_cylinder_low_temperature|false
00|0068|cylinder_top_temperature|system_cylinder_top_temperature|false
00|0069|unknown_1|operation_mode_effect|true
00|006A|multifunction_input_mode|vwz_ai_configuration_me|true
00|006C|installer_name_prefix|installer_name_1|false
00|006D|installer_name_suffix|installer_name_2|false
00|006F|installer_phone_number_prefix|installer_phone_1|false
00|0070|installer_phone_number_suffix|installer_phone_2|false
00|0071|holiday_start_date|system_holiday_start|true
00|0072|holiday_end_date|system_holiday_end|true
00|0073|outdoor_temperature|outside_temperature|false
00|0074|special_function_mode|system_quick_mode|true
00|0076|installer_access_code|installer_menu_code|false
00|0077|installation_assistant_ready|install_assistant_ready|true
00|0078|home_screen_special_function|home_screen_special_function|true
00|0079|heating_zone_time_program_active|status_zone_period_heating|true
00|007A|cooling_zone_time_program_active|status_zone_period_cooling|true
00|007B|heating_zone_operating_mode|operating_mode_heating_for_zone|true
00|007C|cooling_zone_operating_mode|operating_mode_cooling_for_zone|true
00|007D|home_screen_operating_mode|home_screen_operation_mode|true
00|007E|unknown_3|fair_mode|true
00|007F|screed_drying_time|screed_drying_day|true
00|0080|cylinder_loading_start_offset|system_cylinder_loading_start_offset|false
00|0081|cylinder_loading_stop_offset|system_cylinder_loading_stop_offset|false
00|0082|room_temperature_setpoint_at_holiday|room_temperature_setpoint_for_holiday|true
00|0083|bank_holiday_start_date|bank_holiday_start_date|true
00|0084|bank_holiday_end_date|bank_holiday_end_date|true
00|0085|actuator_sensor_test_sensor_reading|actor_sensor_test_sensor_value|true
00|0086|backlight_on_time|backlight_switch_off_time|true
00|0087|drainback_start_time|drainback_start_time|true
00|0088|drainback_station|drainback_station|true
00|0089|solar_pump_kick_time|collector_pump_kick_time|false
00|008A|solar_pump_kick_gradient_limit_minimum|collector_min_kick_gradient|true
00|008B|heat_pump_flow_temperature_setpoint_maximum|heatpump_max_flow_setpoint|false
00|008C|boiler_remote_reset_available|boiler_with_remote_reset|true
00|008D|display_test_pattern|display_test_pattern|true
00|008E|backlight_mode|backlight_mode|false
00|008F|cylinder_type|system_cylinder_available|true
00|0090|unknown_4|plugged_in|true
00|0091|left_button_active|button_pressed_left|true
00|0092|right_button_active|button_pressed_right|true
00|0093|encoder_active|encoder_rotated|true
00|0094|external_crystal_status|external_crystal_state|true
00|0095|outdoor_temperature_daily_average|outside_temperature_average_24h|false
00|0096|maintenance_active|maintenance_mode|true
00|0097|eeprom_update_active|eeprom_update_in_progress|true
00|0098|boiler_reset_command|boiler_remote_reset|false
00|0099|rtc_crystal_ok|rtc_crystal_state|true
00|009A|unknown_5|green_iq|true
00|009B|time_program_reset_command|reset_to_default_time_program|false
00|009C|home_screen_headline|home_screen_headline|true
00|009D|dhw_cylinder_top_temperature|system_cylinder_mss_dhw_top_temperature|false
00|009E|dhw_cylinder_bottom_temperature|system_cylinder_mss_dhw_bottom_temperature|false
00|009F|heating_cylinder_top_temperature|system_cylinder_mss_heating_top_temperature|false
00|00A0|heating_cylinder_bottom_temperature|system_cylinder_mss_heating_bottom_temperature|false
00|00A1|pv_surplus_enabled|smart_pv_activation|true
00|00A2|heating_cylinder_offset_at_pv_surplus|system_cylinder_smart_pv_ch_offset|false
00|00A3|ventilation_time_program_active|system_ventilation_status_period|true
00|00A4|ventilation_special_function_state|system_ventilation_special_function_status|true
00|00A5|cascade_client_sequence|cascade_client_sequence|true
00|00A6|menu_button_active|button_pressed_menu|true
00|00A7|confirm_button_active|button_pressed_ok|true
00|00A8|back_button_active|button_pressed_back|true
00|00A9|up_button_active|button_pressed_up|true
00|00AA|down_button_active|button_pressed_down|true
00|00AB|isolation_switch_active|isolation_switch|true
00|00AF|cooling_electricity_consumption_current_month|system_consumption_electricity_cooling_current_month|false
00|00B0|cooling_electricity_consumption_previous_month|system_consumption_electricity_cooling_previous_month|false
00|00B1|cooling_electricity_consumption_total|system_consumption_electricity_cooling_total|false
00|00B2|network_identifier|net_id|true
00|00B3|multifunction_output_1_mode|vwz_ai_configuration_ma1|true
00|00B5|outdoor_temperature_sensor_ok|outside_temperature_sensor_visibility|false
00|00B6|external_energy_management_enabled|ext_energy_management_activation|true
00|00B8|solar_diverter_valve|system_solar_3wv|true
00|00B9|cooling_electricity_consumption_current_year|system_consumption_electricity_cooling_current_year|false
00|00BA|cooling_electricity_consumption_previous_year|system_consumption_electricity_cooling_previous_year|false
00|00BB|dhw_electricity_consumption_current_year|system_consumption_electricity_dhw_current_year|false
00|00BC|dhw_electricity_consumption_previous_year|system_consumption_electricity_dhw_previous_year|false
00|00BD|heating_fuel_consumption_previous_year|system_consumption_fuel_heating_previous_year|false
00|00BF|dhw_fuel_consumption_previous_year|system_consumption_fuel_dhw_previous_year|false
00|00C0|heating_electricity_consumption_previous_year|system_consumption_electricity_heating_previous_year|false
00|00C1|heating_electricity_consumption_current_year|system_consumption_electricity_heating_current_year|false
00|00C2|electricity_consumption_current_month|system_consumption_electricity_system_current_month|false
00|00C3|electricity_consumption_previous_month|system_consumption_electricity_system_previous_month|false
00|00C4|electricity_consumption_current_year|system_consumption_electricity_system_current_year|false
00|00C5|electricity_consumption_previous_year|system_consumption_electricity_system_previous_year|false
00|00C6|electricity_consumption_total|system_consumption_electricity_system_total|true
00|00C7|dhw_fuel_consumption_current_year|system_consumption_fuel_dhw_current_year|false
00|00C8|heating_fuel_consumption_current_year|system_consumption_fuel_heating_current_year|false
00|00C9|fuel_consumption_current_month|system_consumption_fuel_system_current_month|false
00|00CA|fuel_consumption_previous_month|system_consumption_fuel_system_previous_month|false
00|00CB|fuel_consumption_current_year|system_consumption_fuel_system_current_year|false
00|00CC|fuel_consumption_previous_year|system_consumption_fuel_system_previous_year|false
00|00CD|fuel_consumption_total|system_consumption_fuel_system_total|true
00|00CE|heat_recovery_yield_previous_month|system_yield_heatrecovery_previous_month|false
00|00CF|heat_recovery_yield_current_month|system_yield_heatrecovery_current_month|false
00|00D0|heat_recovery_yield_current_year|system_yield_heatrecovery_current_year|false
00|00D1|heat_recovery_yield_previous_year|system_yield_heatrecovery_previous_year|false
00|00D2|solar_yield_previous_month|system_yield_solar_previous_month|false
00|00D3|solar_yield_current_month|system_yield_solar_current_month|false
00|00D4|solar_yield_current_year|system_yield_solar_current_year|false
00|00D5|solar_yield_previous_year|system_yield_solar_previous_year|false
00|00D6|environmental_yield_previous_month|system_yield_environmental_previous_month|false
00|00D7|environmental_yield_previous_year|system_yield_environmental_previous_year|false
00|00D8|environmental_yield_current_year|system_yield_environmental_current_year|false
00|00D9|electricity_tariff_type|electricity_tariff_type|false
00|00DA|manual_cooling_start_date|manual_cooling_date_start|true
00|00DB|manual_cooling_end_date|manual_cooling_date_end|true
00|00DC|room_temperature_control_mode|room_temperature_control_strategy_type|true
00|00DD|screed_drying_profile_1_temperature_setpoint|screed_drying_profile_day_01_10|false
00|00DE|screed_drying_profile_2_temperature_setpoint|screed_drying_profile_day_11_20|false
00|00DF|screed_drying_profile_3_temperature_setpoint|screed_drying_profile_day_21_29|false
00|00E0|environmental_yield_current_month|system_yield_environmental_current_month|false
00|00E1|circuit_control_mode|circuit_control_strategy|true
00|00E2|automatic_control_active|automatic_mode_active|true
00|00E3|up_button_level|button_raw_value_up|true
00|00E4|down_button_level|button_raw_value_down|true
00|00E5|help_button_level|button_raw_value_question|true
00|00E6|confirm_button_level|button_raw_value_ok|true
00|00E7|menu_button_level|button_raw_value_menu|true
00|00E8|back_button_level|button_raw_value_back|true
00|00E9|slider_button_1_level|button_raw_value_slider_1|true
00|00EA|slider_button_2_level|button_raw_value_slider_2|true
00|00EB|slider_button_3_level|button_raw_value_slider_3|true
00|00EC|slider_button_4_level|button_raw_value_slider_4|true
00|00ED|slider_button_5_level|button_raw_value_slider_5|true
00|00EE|slider_button_6_level|button_raw_value_slider_6|true
00|00EF|continuous_heating_room_temperature_setpoint|continuous_heating_room_setpoint|false
00|00F0|priority_error_code|error_number_priority|true
00|00F1|priority_maintenance_code|maintenance_number_priority|true
00|00F2|holiday_start_date|system_holiday_start_date|true
00|00F3|holiday_end_date|system_holiday_end_date|true
00|00F4|holiday_start_clock_time|system_holiday_start_time|true
00|00F5|holiday_end_clock_time|system_holiday_end_time|false
00|00F6|room_temperature_setpoint_at_holiday|room_temperature_setpoint_for_holiday|true
00|00F7|holiday_abort_command|system_holiday_abort|false
00|00F8|manual_cooling_abort_command|manual_cooling_abort|false
00|00F9|ventilation_boost_active|system_ventilation_boost|false
00|00FA|external_demand_input_mode|external_heat_demand_input_configuration|false
00|00FB|heat_pump_temperature_difference_limit_minimum|heatpump_mss_min_difference|true
00|00FC|dhw_heat_pump_flow_temperature_setpoint_maximum|heatpump_max_flow_dhw_setpoint|false
00|00FD|boiler_flow_temperature_setpoint_maximum|boiler_max_flow_setpoint|true
01|0001|dhw_circuit_type|native_dhw_circuit_type|true
01|0002|dhw_circulation_pump_active|native_dhw_circulation_pump_status|false
01|0003|dhw_operating_mode|native_dhw_operation_mode|false
01|0004|dhw_temperature_setpoint|native_dhw_tapping_temperature_setpoint|false
01|0005|dhw_temperature|native_dhw_current_temperature|false
01|0006|dhw_pump_active|native_dhw_pump_status|false
01|0007|dhw_time_program_reset_command|native_dhw_program_reset_to_default|false
01|0008|dhw_flow_temperature_setpoint|native_dhw_current_target_flow_temperature|false
01|0009|dhw_holiday_start_date|native_dhw_holiday_start_date|false
01|000A|dhw_holiday_end_date|native_dhw_holiday_end_date|false
01|000B|dhw_bank_holiday_start_date|native_dhw_bank_holiday_start|false
01|000C|dhw_bank_holiday_end_date|native_dhw_bank_holiday_end|false
01|000D|dhw_special_function_mode|native_dhw_special_function_mode|false
01|000E|dhw_central_control_active|native_dhw_controlled_by_system|true
01|000F|dhw_special_function_state|native_dhw_status|false
01|0010|dhw_holiday_start_clock_time|native_dhw_holiday_start_time|false
01|0011|dhw_holiday_end_clock_time|native_dhw_holiday_end_time|false
01|0012|dhw_holiday_abort_command|native_dhw_holiday_abort|false
01|0014|dhw_preference|native_dhw_preference|true
01|0015|dhw_recharge_setpoint|native_dhw_recharge_tapping_setpoint|true
01|0019|dhw_loading_efficiency_state|native_dhw_loading_efficiency_state|true
02|0001|hydraulic_circuit_type|circuit_circuit_type|true
02|0002|circuit_type|circuit_mixer_type_external|false
02|0003|circuit_room_influence_mode|circuit_influence_type|false
02|0004|circuit_return_temperature_setpoint|circuit_backflow_temperature_setpoint|false
02|0005|circuit_condensate_sensitive_emitter_available|circuit_condensate_sensitive_emitter|false
02|0006|circuit_cooling_enabled|circuit_cooling_allowed|false
02|0007|circuit_flow_temperature_setpoint|circuit_current_target_flow_temperature|false
02|0008|circuit_flow_temperature|circuit_current_flow_temperature|false
02|0009|dhw_circuit_temperature_setpoint|circuit_dhw_tapping_setpoint|false
02|000A|circuit_unknown_7|circuit_epsilon|true
02|000B|heating_circuit_flow_temperature_offset|circuit_flow_temperature_setpoint_correction_heating|false
02|000C|circuit_flow_high_temperature_setpoint|circuit_flow_temperature_setpoint_high|false
02|000D|circuit_flow_low_temperature_setpoint|circuit_flow_temperature_setpoint_low|false
02|000E|circuit_frost_protection_mode|circuit_frost_protection_mode|true
02|000F|circuit_heating_curve_factor|circuit_heating_curve|false
02|0010|heating_circuit_flow_temperature_setpoint_maximum|circuit_heating_flow_temperature_max_setpoint|false
02|0011|cooling_circuit_flow_temperature_setpoint_minimum|circuit_cooling_flow_temperature_min_setpoint|false
02|0012|circuit_flow_temperature_setpoint_minimum|circuit_flow_temperature_min_setpoint|false
02|0013|dhw_circuit_operating_mode|circuit_dhw_operating_mode|true
02|0014|heating_circuit_outdoor_temperature_limit_maximum|circuit_maximum_outside_temperature_heating|false
02|0015|circuit_room_temperature_influence_mode|circuit_room_temperature_influence|false
02|0016|circuit_screed_drying_time|circuit_screed_drying_day|true
02|0017|screed_drying_circuit_temperature_setpoint|circuit_screed_drying_setpoint|false
02|0018|dhw_circuit_circulation_pump_active|circuit_dhw_circulation_pump_status|true
02|0019|circuit_external_demand_input|circuit_external_heat_demand|true
02|001A|circuit_mixer_movement|circuit_status_heating_circuit_mixer|false
02|001B|circuit_pump_unknown_8|circuit_pump_status|true
02|001C|circuit_adaptive_heating_curve_offset|circuit_adaptive_heating_curve_offset|false
02|001D|circuit_unknown_9|circuit_dhw_quick_mode|true
02|001E|circuit_state|circuit_status_circuit|true
02|001F|cooling_circuit_outdoor_temperature_limit_minimum|circuit_minimum_outside_temperature_cooling|false
02|0020|circuit_unknown_10|circuit_status_automatic_heating_cooling|true
02|0021|circuit_mixer_position|circuit_mixer_position_percentage|false
02|0022|circuit_room_humidity|circuit_current_room_humidity|false
02|0023|circuit_dew_point_temperature|circuit_dew_point_temperature|false
02|0024|circuit_pump_operating_time_total|circuit_pump_operating_hours|false
02|0025|circuit_pump_start_count_total|circuit_pump_starts_count|false
03|0001|cooling_zone_operating_mode|zone_cooling_operation_mode|false
03|0002|cooling_zone_room_temperature_setpoint|zone_room_temperature_setpoint_cooling|false
03|0003|zone_holiday_start_date|zone_holiday_start_date|false
03|0004|zone_holiday_end_date|zone_holiday_end_date|false
03|0005|zone_room_temperature_setpoint_at_holiday|zone_holiday_setpoint|false
03|0006|heating_zone_operating_mode|zone_heating_operation_mode|false
03|0007|heating_zone_room_comfort_temperature_setpoint|zone_comfort_room_temperature_setpoint_heating|true
03|0008|zone_room_temperature_setpoint_at_quick_veto|zone_quick_veto_setpoint|false
03|0009|heating_zone_room_setback_temperature_setpoint|zone_setback_room_temperature_setpoint_heating|false
03|000A|heating_zone_operating_mode_at_teleswitch|zone_operating_mode_heating_teleswitch|true
03|000B|cooling_zone_operating_mode_at_teleswitch|zone_operating_mode_cooling_teleswitch|true
03|000C|zone_bank_holiday_start_date|zone_bank_holiday_start|false
03|000D|zone_bank_holiday_end_date|zone_bank_holiday_end|false
03|000E|zone_special_function_mode|zone_quick_mode|false
03|000F|zone_room_temperature|zone_current_measured_temperature|false
03|0010|zone_time_program_reset_command|zone_schedule_reset_to_default|false
03|0011|zone_enabled|zone_is_active|false
03|0012|zone_valve_position|zone_valve_status|false
03|0013|zone_room_sensor_binding_mode|zone_binding_of_zone|false
03|0014|heating_zone_room_effective_temperature_setpoint|zone_heating_setpoint|false
03|0015|cooling_zone_room_effective_temperature_setpoint|zone_cooling_setpoint|false
03|0016|zone_name|zone_name|false
03|0017|zone_name_prefix|zone_name_prefix|false
03|0018|zone_name_suffix|zone_name_suffix|false
03|0019|heating_zone_time_program_active|zone_heating_schedule_status|false
03|001A|cooling_zone_time_program_active|zone_cooling_schedule_status|false
03|001B|zone_special_function_state|zone_special_function_status|false
03|001C|zone_circuit_assignment|zone_circuit_for_zone|false
03|001D|heating_zone_event_end_clock_time|zone_heating_event_end_time|true
03|001E|zone_quick_veto_end_clock_time|zone_quick_veto_end_time|false
03|0020|zone_holiday_end_clock_time|zone_holiday_end_time|false
03|0021|zone_holiday_start_clock_time|zone_holiday_start_time|false
03|0022|heating_zone_room_manual_temperature_setpoint|zone_heating_manual_setpoint|false
03|0023|cooling_zone_room_manual_temperature_setpoint|zone_cooling_manual_setpoint|false
03|0024|zone_quick_veto_end_date|zone_quick_veto_end_date|false
03|0025|zone_holiday_abort_command|zone_holiday_abort|true
03|0026|zone_quick_veto_time|zone_quick_veto_duration|false
03|0027|zone_manual_cooling_active|zone_manual_cooling_is_active_for_zone|false
03|0028|zone_room_humidity|zone_current_humidity|false
03|0029|zone_cooling_enabled|zone_cooling_allowed|true
03|002A|zone_summer_cutoff_active|zone_summer_cutoff_active|true
03|002B|zone_continuous_heating_active|zone_continuous_heating_active|true
03|002C|zone_frost_protection_active|zone_frost_protection_active|false
03|002D|heating_zone_thermostat_cutoff_active|zone_heating_roomthermostat_status|true
03|002E|cooling_zone_thermostat_cutoff_active|zone_cooling_roomthermostat_status|true
04|0001|solar_pump_operating_time_total_reset_command|solar_collector_runtime_reset|false
04|0002|solar_pump_kick_enabled|solar_collector_kick_function|false
04|0003|solar_collector_temperature|solar_collector_temperature|false
04|0004|solar_collector_temperature_limit_minimum|solar_collector_min_temperature|false
04|0005|solar_circuit_protection_temperature_limit|solar_circuit_protection_function|false
04|0006|solar_circuit_protection_stop_temperature_limit|solar_circuit_protection_function_setpoint|false
04|0007|solar_yield_temperature|yield_solar_temperature|false
04|0008|solar_pump_active|solar_pump_active|false
04|0009|solar_circuit_flow_rate|solar_flow_rate_current|false
04|000A|solar_circuit_flow_rate_setpoint|solar_flow_rate_setpoint|true
04|000B|solar_pump_operating_time_total|solar_pump_hours|false
04|000C|solar_pump_power_limit_minimum|solar_pump_min_power|true
05|0001|solar_cylinder_temperature_limit_maximum|solar_cylinder_max_temperature_setpoint|false
05|0002|solar_cylinder_start_temperature_difference_limit|solar_cylinder_start_temperature_difference|false
05|0003|solar_cylinder_stop_temperature_difference_limit|solar_cylinder_stop_temperature_difference|false
05|0004|solar_cylinder_bottom_temperature|solar_cylinder_bottom_temperature|false
08|0001|differential_control_temperature_limit_maximum|delta_temperature_control_max_temperature|false
08|0002|differential_control_temperature_limit_minimum|delta_temperature_control_min_temperature|false
08|0003|differential_control_start_temperature_difference_limit|delta_temperature_control_start_temperature_difference|false
08|0004|differential_control_stop_temperature_difference_limit|delta_temperature_control_stop_temperature_difference|false
08|0005|differential_control_sensor_1_temperature|delta_temperature_td1_temperature|false
08|0006|differential_control_sensor_2_temperature|delta_temperature_td2_temperature|false
08|0007|differential_control_output|delta_temperature_td_output|true
09|0001|ventilation_special_function_mode|ventilation_quick_mode|false
09|0002|ventilation_operating_mode|ventilation_operating_mode|false
09|0003|ventilation_time_program_active|ventilation_status_period|false
09|0004|ventilation_special_function_state|ventilation_status_special_operating_mode|true
09|0005|ventilation_volatile_organic_compound_sensor_1|ventilation_voc_sensor_1|true
09|0006|ventilation_volatile_organic_compound_sensor_2|ventilation_voc_sensor_2|true
09|0007|ventilation_holiday_end_date|ventilation_holiday_end|false
09|0008|ventilation_holiday_end_clock_time|ventilation_holiday_end_time|false
09|0009|ventilation_holiday_start_date|ventilation_holiday_start|false
09|000A|ventilation_holiday_start_clock_time|ventilation_holiday_start_time|false
09|000B|ventilation_heat_recovery|ventilation_heat_recovery_ventilation|true
09|000C|ventilation_volatile_organic_compound_limit_maximum|ventilation_max_voc|true
09|000D|ventilation_fan_speed_limit_maximum_at_day|ventilation_day_max_fan_stage|false
09|000E|ventilation_fan_speed_limit_maximum_at_night|ventilation_night_max_fan_stage|false
09|000F|ventilation_holiday_abort_command|ventilation_holiday_abort|false
`

func ab36Pointer[ab36T any](ab36Value ab36T) *ab36T { return &ab36Value }

func ab36ExpectedRow(ab36Test *ab36testing.T, ab36Line string) map[string]any {
	ab36Test.Helper()
	ab36Parts := ab36strings.Split(ab36Line, "|")
	ab36Operation := uint8(2)
	ab36GroupText, ab36RegisterText := ab36Parts[0], ab36Parts[1]
	if len(ab36Parts) == 4 {
		ab36Operation = 6
		ab36GroupText, ab36RegisterText = ab36Parts[1], ab36Parts[2]
	}
	ab36Group, ab36Error := ab36strconv.ParseUint(ab36GroupText, 16, 8)
	if ab36Error != nil {
		ab36Test.Fatal(ab36Error)
	}
	ab36Register, ab36Error := ab36strconv.ParseUint(ab36RegisterText, 16, 16)
	if ab36Error != nil {
		ab36Test.Fatal(ab36Error)
	}
	ab36Row := map[string]any{
		"key": Ab36Key{
			Ab36Family: "unqualified", Ab36Variant: "unqualified", Ab36Primary: 181,
			Ab36Secondary: ab36Pointer(uint8(36)), Ab36Operation: ab36Pointer(ab36Operation),
			Ab36Direction: "none", Ab36Field: "name", Ab36Alternative: "none",
			Ab36Group: ab36Pointer(uint8(ab36Group)), Ab36InstanceScope: "unqualified",
			Ab36Register: ab36Pointer(uint16(ab36Register)),
		},
		"canonical_name": nil, "legacy_name": nil, "display_alias": nil, "unresolved": nil,
		"representation_grade": "Proven", "grade": "Hypothesis", "semantics_grade": "Unknown",
		"physical_qualification": "Unknown", "raw_form": ab36Line,
	}
	if ab36Operation == 2 {
		ab36Row["canonical_name"], ab36Row["legacy_name"] = ab36Parts[2], ab36Parts[3]
		ab36Row["unresolved"] = ab36Parts[4] == "true"
	} else {
		ab36Row["display_alias"] = ab36Parts[3]
	}
	return ab36Row
}

func ab36Encode(ab36Test *ab36testing.T, ab36Value any) []byte {
	ab36Test.Helper()
	ab36Encoded, ab36Error := ab36json.Marshal(ab36Value)
	if ab36Error != nil {
		ab36Test.Fatal(ab36Error)
	}
	return ab36Encoded
}

func ab36EqualJSON(ab36Test *ab36testing.T, ab36Actual, ab36ExpectedJSON []byte) {
	ab36Test.Helper()
	var ab36ActualValue, ab36ExpectedValue any
	if ab36Error := ab36json.Unmarshal(ab36Actual, &ab36ActualValue); ab36Error != nil {
		ab36Test.Fatal(ab36Error)
	}
	if ab36Error := ab36json.Unmarshal(ab36ExpectedJSON, &ab36ExpectedValue); ab36Error != nil {
		ab36Test.Fatal(ab36Error)
	}
	if !ab36reflect.DeepEqual(ab36ActualValue, ab36ExpectedValue) {
		ab36Test.Fatalf("ab36 whole row mismatch:\n%s\nwant:\n%s", ab36Actual, ab36ExpectedJSON)
	}
}

func TestAb36Correspondence(ab36Test *ab36testing.T) {
	ab36Entries, ab36Error := Ab36Entries()
	if ab36Error != nil {
		ab36Test.Fatal(ab36Error)
	}
	var ab36Rows []ab36json.RawMessage
	if ab36Error := ab36json.Unmarshal(ab36Entries, &ab36Rows); ab36Error != nil {
		ab36Test.Fatal(ab36Error)
	}
	ab36Lines := ab36strings.Split(ab36strings.TrimSpace(ab36Expected), "\n")
	if len(ab36Lines) != 368 || len(ab36Rows) != 370 {
		ab36Test.Fatalf("ab36 counts: %d / %d", len(ab36Lines), len(ab36Rows))
	}
	ab36Lines = append(ab36Lines, "06|0A|0007|RoomHumidity", "06|0A|000F|RoomTemp")
	ab36Seen := make(map[string]bool)
	ab36Aliases := 0
	for ab36Position, ab36Line := range ab36Lines {
		ab36ExpectedRowValue := ab36ExpectedRow(ab36Test, ab36Line)
		ab36Key := ab36ExpectedRowValue["key"].(Ab36Key)
		ab36KeyBytes := ab36Encode(ab36Test, ab36Key)
		if ab36Seen[string(ab36KeyBytes)] {
			ab36Test.Fatal("ab36 duplicate expected key")
		}
		ab36Seen[string(ab36KeyBytes)] = true
		ab36Want := ab36Encode(ab36Test, ab36ExpectedRowValue)
		ab36EqualJSON(ab36Test, ab36Rows[ab36Position], ab36Want)
		ab36Actual, ab36Found, ab36Error := Ab36Lookup(ab36Key)
		if ab36Error != nil || !ab36Found {
			ab36Test.Fatalf("ab36 missing %s: %v", ab36Line, ab36Error)
		}
		ab36EqualJSON(ab36Test, ab36Actual, ab36Want)
		if *ab36Key.Ab36Operation == 2 && ab36ExpectedRowValue["canonical_name"] != ab36ExpectedRowValue["legacy_name"] {
			ab36Aliases++
		}
	}
	if len(ab36Seen) != 370 || ab36Aliases != 334 {
		ab36Test.Fatalf("ab36 unique/aliases: %d / %d", len(ab36Seen), ab36Aliases)
	}
	ab36Test.Logf("ab36 verified 368 OP02 keys, 334 differing aliases, 2 OP06 keys; all whole rows")
}

func TestAb36ExactAxes(ab36Test *ab36testing.T) {
	ab36Base := ab36ExpectedRow(ab36Test, "06|0A|000F|RoomTemp")["key"].(Ab36Key)
	ab36Controls, ab36Found, ab36Error := Ab36Lookup(ab36Base)
	if ab36Error != nil || !ab36Found || len(ab36Controls) == 0 {
		ab36Test.Fatal("ab36 missing positive control")
	}
	for ab36Field := 0; ab36Field < ab36reflect.TypeOf(ab36Base).NumField(); ab36Field++ {
		ab36Name := ab36reflect.TypeOf(ab36Base).Field(ab36Field).Name
		ab36Test.Run(ab36Name, func(ab36Test *ab36testing.T) {
			ab36Query := ab36Base
			ab36Value := ab36reflect.ValueOf(&ab36Query).Elem().Field(ab36Field)
			switch ab36Value.Kind() {
			case ab36reflect.String:
				ab36Value.SetString("ab36_other")
			case ab36reflect.Uint8:
				ab36Value.SetUint(254)
			case ab36reflect.Ptr:
				ab36Replacement := ab36reflect.New(ab36Value.Type().Elem())
				ab36Replacement.Elem().SetUint(254)
				ab36Value.Set(ab36Replacement)
			default:
				ab36Test.Fatal("ab36 untested axis")
			}
			ab36Actual, ab36Found, ab36Error := Ab36Lookup(ab36Query)
			if ab36Error != nil || ab36Found || ab36Actual != nil {
				ab36Test.Fatalf("ab36 fallback on %s: %s %v", ab36Name, ab36Actual, ab36Error)
			}
			if ab36Value.Kind() == ab36reflect.Ptr {
				for _, ab36Zero := range []bool{false, true} {
					if ab36Zero {
						ab36Value.Set(ab36reflect.New(ab36Value.Type().Elem()))
					} else {
						ab36Value.SetZero()
					}
					ab36Actual, ab36Found, ab36Error := Ab36Lookup(ab36Query)
					if ab36Error != nil || ab36Found || ab36Actual != nil {
						ab36Test.Fatalf("ab36 absent/zero fallback on %s", ab36Name)
					}
				}
			}
		})
	}
	ab36Query := ab36Base
	ab36Query.Ab36Operation = ab36Pointer(uint8(2))
	ab36Actual, ab36Found, ab36Error := Ab36Lookup(ab36Query)
	if ab36Error != nil || ab36Found || ab36Actual != nil {
		ab36Test.Fatal("ab36 OP06 merged into OP02")
	}
}

func TestAb36IndependentStorage(ab36Test *ab36testing.T) {
	ab36Key := ab36ExpectedRow(ab36Test, "06|0A|0007|RoomHumidity")["key"].(Ab36Key)
	ab36Row, ab36Found, ab36Error := Ab36Lookup(ab36Key)
	if ab36Error != nil || !ab36Found {
		ab36Test.Fatal("ab36 missing row")
	}
	ab36Want := append([]byte(nil), ab36Row...)
	ab36Entries, ab36Error := Ab36Entries()
	if ab36Error != nil {
		ab36Test.Fatal(ab36Error)
	}
	ab36WantEntries := append([]byte(nil), ab36Entries...)
	for ab36Position := range ab36Row {
		ab36Row[ab36Position] = 'x'
	}
	for ab36Position := range ab36Entries {
		ab36Entries[ab36Position] = 'x'
	}
	ab36Again, ab36Found, ab36Error := Ab36Lookup(ab36Key)
	if ab36Error != nil || !ab36Found || !ab36bytes.Equal(ab36Again, ab36Want) {
		ab36Test.Fatal("ab36 row storage shared")
	}
	ab36AgainEntries, ab36Error := Ab36Entries()
	if ab36Error != nil || !ab36bytes.Equal(ab36AgainEntries, ab36WantEntries) {
		ab36Test.Fatal("ab36 entries storage shared")
	}
	ab36Source := append([]byte(nil), ab36WantEntries...)
	ab36Rows, ab36Index, ab36Error := ab36Validate(ab36Source)
	if ab36Error != nil {
		ab36Test.Fatal(ab36Error)
	}
	for ab36Position := range ab36Source {
		ab36Source[ab36Position] = 'x'
	}
	ab36EqualJSON(ab36Test, ab36Rows[368], ab36Want)
	ab36EqualJSON(ab36Test, ab36Index[string(ab36Encode(ab36Test, ab36Key))], ab36Want)
	ab36Again, ab36Found, ab36Error = Ab36Lookup(ab36Key)
	if ab36Error != nil || !ab36Found || !ab36bytes.Equal(ab36Again, ab36Want) {
		ab36Test.Fatal("ab36 source mutation affected API")
	}
}

func TestAb36InvalidData(ab36Test *ab36testing.T) {
	ab36Row := ab36ExpectedRow(ab36Test, "06|0A|0007|RoomHumidity")
	ab36Valid := string(ab36Encode(ab36Test, []any{ab36Row}))
	ab36Duplicate := string(ab36Encode(ab36Test, []any{ab36Row, ab36Row}))
	ab36OtherRow := ab36ExpectedRow(ab36Test, "06|0A|0007|RoomHumidity")
	ab36OtherRow["grade"] = "Unknown"
	ab36OtherRow["ab36_width"] = 2
	ab36DifferentDuplicate := string(ab36Encode(ab36Test, []any{ab36Row, ab36OtherRow}))
	ab36Original := ab36Data
	ab36Test.Cleanup(func() { ab36Data = ab36Original })
	for ab36Position, ab36Bad := range []string{
		"", "[", "null", "{}", "[null]", "[42]", "[{}]", "[{\"key\":null}]", ab36Duplicate, ab36DifferentDuplicate,
		ab36strings.Replace(ab36Valid, `"primary":181,`, ``, 1),
		ab36strings.Replace(ab36Valid, `"primary":181`, `"primary":null`, 1),
		ab36strings.Replace(ab36Valid, `"primary":181`, `"primary":256`, 1),
		ab36strings.Replace(ab36Valid, `"register":7`, `"register":65536`, 1),
		ab36strings.Replace(ab36Valid, `"secondary":36`, `"secondary":-1`, 1),
		ab36strings.Replace(ab36Valid, `"operation":6`, `"operation":"6"`, 1),
		ab36strings.Replace(ab36Valid, `"family":"unqualified"`, `"family":""`, 1),
		ab36strings.Replace(ab36Valid, `"variant":"unqualified",`, ``, 1),
		ab36strings.Replace(ab36Valid, `"primary":181`, `"ab36_unknown":181`, 1),
	} {
		ab36Test.Run("ab36_invalid_"+ab36strconv.Itoa(ab36Position), func(ab36Test *ab36testing.T) {
			ab36Rows, ab36Index, ab36Error := ab36Validate([]byte(ab36Bad))
			if ab36Error == nil || ab36Rows != nil || ab36Index != nil {
				ab36Test.Fatal("ab36 invalid metadata accepted")
			}
			ab36Data = ab36Bad
			ab36Entries, ab36Error := Ab36Entries()
			if ab36Error == nil || ab36Entries != nil {
				ab36Test.Fatal("ab36 Entries swallowed invalid metadata")
			}
			ab36Actual, ab36Found, ab36Error := Ab36Lookup(ab36Row["key"].(Ab36Key))
			if ab36Error == nil || ab36Found || ab36Actual != nil {
				ab36Test.Fatal("ab36 Lookup swallowed invalid metadata")
			}
		})
	}
	ab36Data = "[]"
	ab36Entries, ab36Error := Ab36Entries()
	if ab36Error != nil || string(ab36Entries) != "[]" {
		ab36Test.Fatal("ab36 valid empty array rejected")
	}
}

func TestAb36OptionalAxesAndRawMetadata(ab36Test *ab36testing.T) {
	ab36Original := ab36Data
	ab36Test.Cleanup(func() { ab36Data = ab36Original })
	ab36Base := ab36ExpectedRow(ab36Test, "06|0A|0007|RoomHumidity")
	for _, ab36Field := range []string{"Ab36Secondary", "Ab36Operation", "Ab36Group", "Ab36Register"} {
		ab36Test.Run("ab36_"+ab36Field, func(ab36Test *ab36testing.T) {
			ab36Absent := ab36Base["key"].(Ab36Key)
			ab36reflect.ValueOf(&ab36Absent).Elem().FieldByName(ab36Field).SetZero()
			ab36Zero := ab36Absent
			ab36Value := ab36reflect.ValueOf(&ab36Zero).Elem().FieldByName(ab36Field)
			ab36Value.Set(ab36reflect.New(ab36Value.Type().Elem()))
			ab36Rows := []map[string]any{}
			for ab36Position, ab36Key := range []Ab36Key{ab36Absent, ab36Zero} {
				ab36Row := ab36ExpectedRow(ab36Test, "06|0A|0007|RoomHumidity")
				ab36Row["key"] = ab36Key
				ab36Row["raw_form"] = "ab36_synthetic_" + ab36strconv.Itoa(ab36Position)
				ab36Row["ab36_extra"] = map[string]any{"ab36_bytes": []int{0, 255}, "ab36_state": "Unknown"}
				ab36Rows = append(ab36Rows, ab36Row)
			}
			ab36Data = string(ab36Encode(ab36Test, ab36Rows))
			for _, ab36Row := range ab36Rows {
				ab36Actual, ab36Found, ab36Error := Ab36Lookup(ab36Row["key"].(Ab36Key))
				if ab36Error != nil || !ab36Found {
					ab36Test.Fatal("ab36 absent/zero collapsed")
				}
				ab36EqualJSON(ab36Test, ab36Actual, ab36Encode(ab36Test, ab36Row))
			}
		})
	}
}

func TestAb36OmittedOptionalKeys(ab36Test *ab36testing.T) {
	ab36Original := ab36Data
	ab36Test.Cleanup(func() { ab36Data = ab36Original })
	for _, ab36Field := range []string{"Ab36Secondary", "Ab36Operation", "Ab36Group", "Ab36Register"} {
		ab36Test.Run("ab36_"+ab36Field, func(ab36Test *ab36testing.T) {
			ab36Row := ab36ExpectedRow(ab36Test, "06|0A|0007|RoomHumidity")
			ab36Key := ab36Row["key"].(Ab36Key)
			ab36Value := ab36reflect.ValueOf(&ab36Key).Elem().FieldByName(ab36Field)
			ab36Value.SetZero()
			var ab36KeyObject map[string]any
			if ab36Error := ab36json.Unmarshal(ab36Encode(ab36Test, ab36Key), &ab36KeyObject); ab36Error != nil {
				ab36Test.Fatal(ab36Error)
			}
			ab36Definition, ab36Present := ab36reflect.TypeOf(ab36Key).FieldByName(ab36Field)
			if !ab36Present {
				ab36Test.Fatal("ab36 missing axis")
			}
			delete(ab36KeyObject, ab36Definition.Tag.Get("json"))
			ab36Row["key"] = ab36KeyObject
			ab36Data = string(ab36Encode(ab36Test, []any{ab36Row}))
			ab36Actual, ab36Found, ab36Error := Ab36Lookup(ab36Key)
			if ab36Error != nil || !ab36Found {
				ab36Test.Fatalf("ab36 omitted optional key rejected: %v", ab36Error)
			}
			ab36EqualJSON(ab36Test, ab36Actual, ab36Encode(ab36Test, ab36Row))
			ab36Value.Set(ab36reflect.New(ab36Value.Type().Elem()))
			ab36Actual, ab36Found, ab36Error = Ab36Lookup(ab36Key)
			if ab36Error != nil || ab36Found || ab36Actual != nil {
				ab36Test.Fatal("ab36 omitted key matched explicit zero")
			}
		})
	}
}
