package output

import (
	"errors"
	"strings"
)

// MCPZone represents one space-heating zone in the MCP output.
type MCPZone struct {
	ID     string        `json:"id"`
	Name   string        `json:"name"`
	State  MCPZoneState  `json:"state"`
	Config MCPZoneConfig `json:"config"`
}

// MCPZoneState represents the observed state of a zone in the MCP output.
type MCPZoneState struct {
	CurrentTempC       *float64 `json:"zone_room_temperature,omitempty"`
	CurrentHumidityPct *float64 `json:"zone_room_humidity,omitempty"`
	HvacAction         string   `json:"hvac_action,omitempty"`
	SpecialFunction    string   `json:"zone_special_function_mode,omitempty"`
	HeatingDemandPct   *float64 `json:"heating_demand_pct,omitempty"`
	ValvePositionPct   *float64 `json:"zone_valve_position,omitempty"`
	enumerations       []CatalogueEnumeration
}

// MCPZoneConfig represents the settings of a zone in the MCP output.
type MCPZoneConfig struct {
	OperatingMode              string   `json:"operating_mode"`
	Preset                     string   `json:"preset"`
	TargetTempC                *float64 `json:"target_temp_c,omitempty"`
	AllowedModes               []string `json:"allowed_modes,omitempty"`
	CircuitType                string   `json:"circuit_type,omitempty"`
	AssociatedCircuit          *int     `json:"associated_circuit,omitempty"`
	RoomTemperatureZoneMapping *int     `json:"room_temperature_zone_mapping,omitempty"`
	QuickVeto                  bool     `json:"quick_veto"`
	QuickVetoSetpointC         *float64 `json:"zone_room_temperature_setpoint_at_quick_veto,omitempty"`
	QuickVetoDurationH         *float64 `json:"zone_quick_veto_time,omitempty"`
	QuickVetoExpiry            string   `json:"quick_veto_expiry,omitempty"`
	HolidayStartDate           string   `json:"zone_holiday_start_date,omitempty"`
	HolidayEndDate             string   `json:"zone_holiday_end_date,omitempty"`
	HolidaySetpointC           *float64 `json:"zone_room_temperature_setpoint_at_holiday,omitempty"`
	HolidayStartTime           string   `json:"zone_holiday_start_clock_time,omitempty"`
	HolidayEndTime             string   `json:"zone_holiday_end_clock_time,omitempty"`
}

// MCPDhwStatus represents domestic hot water state and settings in the MCP output.
type MCPDhwStatus struct {
	State  MCPDhwState  `json:"state"`
	Config MCPDhwConfig `json:"config"`
}

// MCPDhwState represents the observed state of domestic hot water in the MCP output.
type MCPDhwState struct {
	CurrentTempC     *float64 `json:"dhw_temperature,omitempty"`
	SpecialFunction  string   `json:"dhw_special_function_mode,omitempty"`
	HeatingDemandPct *float64 `json:"heating_demand_pct,omitempty"`
	enumerations     []CatalogueEnumeration
}

// MCPDhwConfig represents domestic hot water settings in the MCP output.
type MCPDhwConfig struct {
	OperatingMode    string   `json:"dhw_operating_mode"`
	Preset           string   `json:"preset"`
	TargetTempC      *float64 `json:"dhw_temperature_setpoint,omitempty"`
	HolidayStartDate string   `json:"dhw_holiday_start_date,omitempty"`
	HolidayEndDate   string   `json:"dhw_holiday_end_date,omitempty"`
}

// MCPCircuitStatus represents the state, settings, and owner of a circuit in the MCP output.
type MCPCircuitStatus struct {
	Index          int               `json:"index"`
	CircuitType    string            `json:"circuit_type"`
	HasMixer       bool              `json:"has_mixer"`
	State          MCPCircuitState   `json:"state"`
	Config         MCPCircuitConfig  `json:"config"`
	ManagingDevice MCPManagingDevice `json:"managing_device"`
}

// MCPCircuitState represents the observed state of a circuit in the MCP output.
type MCPCircuitState struct {
	PumpActive       *bool    `json:"circuit_pump_unknown_8,omitempty"`
	MixerPositionPct *float64 `json:"circuit_mixer_position,omitempty"`
	FlowTemperatureC *float64 `json:"circuit_flow_temperature,omitempty"`
	FlowSetpointC    *float64 `json:"circuit_flow_temperature_setpoint,omitempty"`
	CalcFlowTempC    *float64 `json:"calc_flow_temp_c,omitempty"`
	CircuitState     string   `json:"circuit_state,omitempty"`
	Humidity         *float64 `json:"circuit_room_humidity,omitempty"`
	DewPoint         *float64 `json:"circuit_dew_point_temperature,omitempty"`
	PumpHours        *float64 `json:"circuit_pump_operating_time_total,omitempty"`
	PumpStarts       *int     `json:"circuit_pump_start_count_total,omitempty"`
}

// MCPCircuitConfig represents the settings of a circuit in the MCP output.
type MCPCircuitConfig struct {
	HeatingCurve        *float64 `json:"circuit_heating_curve_factor,omitempty"`
	FlowTempMaxC        *float64 `json:"heating_circuit_flow_temperature_setpoint_maximum,omitempty"`
	FlowTempMinC        *float64 `json:"circuit_flow_temperature_setpoint_minimum,omitempty"`
	SummerLimitC        *float64 `json:"heating_circuit_outdoor_temperature_limit_maximum,omitempty"`
	SpecialFunctionMode *int     `json:"circuit_special_function_mode,omitempty"`
	RoomTempControl     string   `json:"circuit_room_temperature_influence_mode,omitempty"`
	CoolingEnabled      *bool    `json:"circuit_cooling_enabled,omitempty"`
	enumerations        []CatalogueEnumeration
}

// MCPManagingDevice represents the device responsible for a circuit in the MCP output.
type MCPManagingDevice struct {
	Role     string  `json:"role"`
	DeviceID *string `json:"device_id,omitempty"`
	Address  *int    `json:"address,omitempty"`
}

// MCPRadioDevice represents one attached radio device in the MCP output.
type MCPRadioDevice struct {
	Group                int      `json:"group"`
	Instance             int      `json:"instance"`
	SlotMode             string   `json:"slot_mode"`
	DeviceConnected      *bool    `json:"device_connected,omitempty"`
	DeviceClassAddress   *int     `json:"device_class_address,omitempty"`
	DeviceModel          string   `json:"device_model,omitempty"`
	FirmwareVersion      *string  `json:"firmware_version,omitempty"`
	HardwareIdentifier   *int     `json:"hardware_identifier,omitempty"`
	RemoteControlAddress *int     `json:"remote_control_address,omitempty"`
	DevicePaired         *bool    `json:"device_paired,omitempty"`
	ReceptionStrength    *int     `json:"reception_strength,omitempty"`
	ZoneAssignment       *int     `json:"zone_assignment,omitempty"`
	RoomTemperatureC     *float64 `json:"room_temperature_c,omitempty"`
	RoomHumidityPct      *float64 `json:"room_humidity_pct,omitempty"`
}

// MCPFunctionModuleInterpretation represents the structural result and its evidence reference in the MCP output.
type MCPFunctionModuleInterpretation struct {
	Mode             MCPFunctionModuleSemanticMode            `json:"mode"`
	DegradedReason   *MCPFunctionModuleSemanticDegradedReason `json:"degraded_reason"`
	EvidenceRevision string                                   `json:"evidence_revision"`
}

// MCPFunctionModuleSemanticMode represents the structural mode of the function module in the MCP output.
type MCPFunctionModuleSemanticMode string

// MCPFunctionModuleSemanticDegradedReason represents the reason a function-module interpretation is degraded in the MCP output.
type MCPFunctionModuleSemanticDegradedReason string

// MCPSolarStatus represents the solar plane in the MCP output.
type MCPSolarStatus struct {
	CollectorTemperatureC *float64 `json:"collector_temperature_c,omitempty"`
	ReturnTemperatureC    *float64 `json:"return_temperature_c,omitempty"`
	PumpActive            *bool    `json:"pump_active,omitempty"`
	CurrentYield          *float64 `json:"current_yield,omitempty"`
	PumpHours             *float64 `json:"pump_hours,omitempty"`
	SolarEnabled          *bool    `json:"solar_enabled,omitempty"`
	FunctionMode          *bool    `json:"function_mode,omitempty"`
}

// MCPCylinderStatus represents one hot water cylinder in the MCP output.
type MCPCylinderStatus struct {
	Index             int      `json:"index"`
	TemperatureC      *float64 `json:"temperature_c,omitempty"`
	MaxSetpointC      *float64 `json:"max_setpoint_c,omitempty"`
	ChargeHysteresisC *float64 `json:"charge_hysteresis_c,omitempty"`
	ChargeOffsetC     *float64 `json:"charge_offset_c,omitempty"`
}

// MCPEnergySeries represents the energy values and their metadata for one use in the MCP output.
type MCPEnergySeries struct {
	Today       float64              `json:"today"`
	Yearly      []float64            `json:"yearly"`
	Monthly     []float64            `json:"monthly,omitempty"`
	TodayMeta   MCPEnergyPointMeta   `json:"today_meta"`
	YearlyMeta  []MCPEnergyPointMeta `json:"yearly_meta,omitempty"`
	MonthlyMeta []MCPEnergyPointMeta `json:"monthly_meta,omitempty"`
}

// MCPEnergyChannel represents climate and hot water energy for one source in the MCP output.
type MCPEnergyChannel struct {
	DHW     MCPEnergySeries `json:"dhw"`
	Climate MCPEnergySeries `json:"climate"`
}

// MCPEnergyTotals represents the totals for the supported energy sources in the MCP output.
type MCPEnergyTotals struct {
	Gas                 MCPEnergyChannel `json:"gas"`
	Electric            MCPEnergyChannel `json:"electric"`
	Solar               MCPEnergyChannel `json:"solar"`
	registerDirections  []EnergyDirection
	catalogueDirections []EnergyDirection
	residualFields      []EnergyResidual
	candidateFields     []EnergyResidual
}

// MCPEnergyPointMeta represents the freshness and provenance of one energy value in the MCP output.
type MCPEnergyPointMeta struct {
	FreshnessState  string  `json:"freshness_state"`
	Provenance      string  `json:"provenance"`
	LastObservedUTC string  `json:"last_observed_utc,omitempty"`
	AgeSeconds      float64 `json:"age_seconds,omitempty"`
	Stale           bool    `json:"stale"`
}

// MCPBoilerStatus represents the full boiler status in the MCP output.
type MCPBoilerStatus struct {
	State       *MCPBoilerState       `json:"state,omitempty"`
	Config      *MCPBoilerConfig      `json:"config,omitempty"`
	Diagnostics *MCPBoilerDiagnostics `json:"diagnostics,omitempty"`
	selectors   []BoilerSelectorEvidence
	operations  []RawOperation
}

// MCPBoilerState represents the observed state of the boiler in the MCP output.
type MCPBoilerState struct {
	FlowTemperatureC         *float64 `json:"flow_temperature,omitempty"`
	ReturnTemperatureC       *float64 `json:"return_temperature_c,omitempty"`
	CentralHeatingPumpActive *bool    `json:"central_heating_pump_active,omitempty"`
	WaterPressureBar         *float64 `json:"water_pressure,omitempty"`
	ExternalPumpActive       *bool    `json:"external_pump_active,omitempty"`
	CirculationPumpActive    *bool    `json:"circulation_pump_active,omitempty"`
	GasValveActive           *bool    `json:"gas_valve_active,omitempty"`
	FlameActive              *bool    `json:"flame_active,omitempty"`
	DiverterValvePositionPct *float64 `json:"diverter_valve_position,omitempty"`
	FanSpeedRpm              *int     `json:"fan_speed,omitempty"`
	TargetFanSpeedRpm        *int     `json:"fan_speed_setpoint,omitempty"`
	IonisationVoltageUa      *float64 `json:"flame_ionization_current,omitempty"`
	DhwWaterFlowLpm          *float64 `json:"dhw_water_flow_rate,omitempty"`
	DhwDemandActive          *bool    `json:"dhw_demand_active,omitempty"`
	HeatingSwitchActive      *bool    `json:"heating_switch_active,omitempty"`
	StorageLoadPumpPct       *float64 `json:"storage_loading_pump_level,omitempty"`
	ModulationPct            *float64 `json:"modulation_setpoint,omitempty"`
	PrimaryCircuitFlowLpm    *float64 `json:"primary_circuit_flow_rate,omitempty"`
	FlowTempDesiredC         *float64 `json:"flow_temperature_setpoint,omitempty"`
	DhwTempDesiredC          *float64 `json:"dhw_temperature_setpoint,omitempty"`
	StateNumber              *int     `json:"boiler_state_code,omitempty"`
	DhwTemperatureC          *float64 `json:"dhw_temperature_c,omitempty"`
	DhwTargetTemperatureC    *float64 `json:"dhw_target_temperature_c,omitempty"`
}

// MCPBoilerConfig represents the boiler settings in the MCP output.
type MCPBoilerConfig struct {
	DhwOperatingMode  *string  `json:"dhw_operating_mode,omitempty"`
	FlowsetHcMaxC     *float64 `json:"maximum_heating_flow_temperature_setpoint,omitempty"`
	FlowsetHwcMaxC    *float64 `json:"maximum_dhw_flow_temperature_setpoint,omitempty"`
	PartloadHcKW      *float64 `json:"heating_partial_load_power,omitempty"`
	PartloadHwcKW     *float64 `json:"dhw_partial_load_power,omitempty"`
	InstallerMenuCode *int     `json:"central_heating_pump_mode,omitempty"`
	PhoneNumber       *string  `json:"installer_phone_number,omitempty"`
	HoursTillService  *int     `json:"service_time_remaining,omitempty"`
}

// MCPBoilerDiagnostics represents the boiler operating counters and diagnostic values in the MCP output.
type MCPBoilerDiagnostics struct {
	HeatingStatusRaw         *int     `json:"heating_status_raw,omitempty"`
	DhwStatusRaw             *int     `json:"dhw_status_raw,omitempty"`
	CentralHeatingHours      *float64 `json:"central_heating_operating_time,omitempty"`
	DhwHours                 *float64 `json:"dhw_operating_time,omitempty"`
	CentralHeatingStarts     *int     `json:"central_heating_start_count,omitempty"`
	DhwStarts                *int     `json:"dhw_start_count,omitempty"`
	PumpHours                *float64 `json:"pump_operating_time,omitempty"`
	FanHours                 *float64 `json:"fan_operating_time,omitempty"`
	DeactivationsIFC         *int     `json:"fire_control_deactivation_count,omitempty"`
	DeactivationsTemplimiter *int     `json:"temperature_limiter_deactivation_count,omitempty"`
}

// MCPSystemStatus represents the system status in the MCP output.
type MCPSystemStatus struct {
	State          *MCPSystemState      `json:"state,omitempty"`
	Config         *MCPSystemConfig     `json:"config,omitempty"`
	Properties     *MCPSystemProperties `json:"properties,omitempty"`
	ventilation    *VentilationView
	catalogue      *FamilyCatalogue
	history        *HistoryOwnerRecord
	functionModule *FunctionModuleRecord
	legacyPlanes   *LegacyPlaneRecord
}

// MCPSystemState represents the observed state of the system in the MCP output.
type MCPSystemState struct {
	SystemOff                      *bool    `json:"system_off,omitempty"`
	SystemWaterPressure            *float64 `json:"water_pressure,omitempty"`
	SystemFlowTemperature          *float64 `json:"flow_temperature,omitempty"`
	OutdoorTemperature             *float64 `json:"outdoor_temperature,omitempty"`
	OutdoorTemperatureDailyAverage *float64 `json:"outdoor_temperature_daily_average,omitempty"`
	MaintenanceDue                 *bool    `json:"maintenance_active,omitempty"`
	HwcCylinderTemperatureTop      *float64 `json:"dhw_cylinder_top_temperature,omitempty"`
	HwcCylinderTemperatureBottom   *float64 `json:"dhw_cylinder_bottom_temperature,omitempty"`
}

// MCPSystemConfig represents the system settings and active bound in the MCP output.
type MCPSystemConfig struct {
	AdaptiveHeatingCurve         *bool           `json:"adaptive_heating_curve_enabled,omitempty"`
	AlternativePoint             *float64        `json:"heating_outdoor_alternative_temperature_limit,omitempty"`
	HeatingCircuitBivalencePoint *float64        `json:"heating_outdoor_bivalence_temperature_limit,omitempty"`
	DhwBivalencePoint            *float64        `json:"dhw_outdoor_bivalence_temperature_limit,omitempty"`
	DhwBivalenceBound            *MCPConfigBound `json:"dhw_outdoor_bivalence_temperature_limit_bound,omitempty"`
	HcEmergencyTemperature       *float64        `json:"heating_flow_temperature_setpoint_maximum_at_heat_pump_fault,omitempty"`
	HwcMaxFlowTempDesired        *float64        `json:"dhw_flow_temperature_setpoint_maximum,omitempty"`
	MaxRoomHumidity              *int            `json:"room_humidity_limit_maximum,omitempty"`
	MaintenanceDate              *string         `json:"maintenance_date,omitempty"`
	InstallerName                *string         `json:"installer_name,omitempty"`
	InstallerPhone               *string         `json:"installer_phone_number,omitempty"`
	InstallerMenuCode            *int            `json:"installer_access_code,omitempty"`
}

// MCPSystemProperties represents the hydraulic scheme and module settings in the MCP output.
type MCPSystemProperties struct {
	SystemScheme                    *int `json:"hydraulic_scheme_code,omitempty"`
	FunctionModuleConfigurationCode *int `json:"functional_module_configuration_code,omitempty"`
}

// MCPScheduleTimerSlot represents one scheduled interval in the MCP output.
type MCPScheduleTimerSlot struct {
	StartHour      int      `json:"start_hour"`
	StartMinute    int      `json:"start_minute"`
	EndHour        int      `json:"end_hour"`
	EndMinute      int      `json:"end_minute"`
	TemperatureC   *float64 `json:"temperature_c,omitempty"`
	TemperatureRaw *int     `json:"temperature_raw,omitempty"`
}

// MCPScheduleDayProgram represents the intervals for one day in the MCP output.
type MCPScheduleDayProgram struct {
	Weekday string                 `json:"weekday"`
	Slots   []MCPScheduleTimerSlot `json:"slots"`
}

// MCPScheduleConfig represents the native limits and temperature support of a program in the MCP output.
type MCPScheduleConfig struct {
	MaxSlots       int      `json:"max_slots"`
	TimeResolution int      `json:"time_resolution"`
	MinDuration    int      `json:"min_duration"`
	HasTemperature bool     `json:"has_temperature"`
	TempSlots      int      `json:"temp_slots"`
	MinTempC       *float64 `json:"min_temp_c,omitempty"`
	MaxTempC       *float64 `json:"max_temp_c,omitempty"`
}

// MCPScheduleProgram represents the program for one zone and circuit kind in the MCP output.
type MCPScheduleProgram struct {
	Zone      int                     `json:"zone"`
	HC        string                  `json:"hc"`
	Config    *MCPScheduleConfig      `json:"config,omitempty"`
	SlotsUsed []int                   `json:"slots_used,omitempty"`
	Days      []MCPScheduleDayProgram `json:"days,omitempty"`
}

// MCPScheduleStatus represents the set of programs in the MCP output.
type MCPScheduleStatus struct {
	Programs []MCPScheduleProgram `json:"programs"`
}

// MCPAdapterHardwareInfo represents adapter identity and telemetry in the MCP output.
type MCPAdapterHardwareInfo struct {
	FirmwareVersion    string   `json:"firmware_version"`
	FirmwareChecksum   string   `json:"firmware_checksum"`
	BootloaderVersion  string   `json:"bootloader_version"`
	BootloaderChecksum string   `json:"bootloader_checksum"`
	HardwareID         string   `json:"hardware_id"`
	HardwareConfig     string   `json:"hardware_config"`
	Features           byte     `json:"features"`
	Jumpers            byte     `json:"jumpers"`
	JumperFlags        []string `json:"jumper_flags"`
	IsWiFi             bool     `json:"is_wifi"`
	IsEthernet         bool     `json:"is_ethernet"`
	TemperatureC       *float64 `json:"temperature_c,omitempty"`
	SupplyVoltageMV    *int     `json:"supply_voltage_mv,omitempty"`
	BusVoltageMaxDV    *int     `json:"bus_voltage_max_dv,omitempty"`
	BusVoltageMinDV    *int     `json:"bus_voltage_min_dv,omitempty"`
	ResetCause         *string  `json:"reset_cause,omitempty"`
	ResetCauseCode     *byte    `json:"reset_cause_code,omitempty"`
	RestartCount       *byte    `json:"restart_count,omitempty"`
	WiFiRSSIDBm        *int     `json:"wifi_rssi_dbm,omitempty"`
	LastIdentityQuery  *string  `json:"last_identity_query,omitempty"`
	LastTelemetryQuery *string  `json:"last_telemetry_query,omitempty"`
	VersionResponseLen int      `json:"version_response_len"`
	InfoSupported      bool     `json:"info_supported"`
}

// MCPConfigBound represents the active bound and its target identity in the MCP output.
type MCPConfigBound struct {
	NarrowedOnly bool    `json:"narrowed_only"`
	MinC         float64 `json:"min_c"`
	MaxC         float64 `json:"max_c"`
	Source       string  `json:"source"`
	Reason       string  `json:"reason,omitempty"`
	Profile      string  `json:"profile,omitempty"`
}

// MCPTimeProgramWriteResult represents the result of a time-program write in the MCP output.
type MCPTimeProgramWriteResult struct {
	Success     bool                       `json:"success"`
	SlotResults []MCPTimeProgramSlotResult `json:"slot_results"`
	Error       string                     `json:"error,omitempty"`
}

// MCPTimeProgramSlot represents one requested interval of a time-program write in the MCP output.
type MCPTimeProgramSlot struct {
	StartHour    int      `json:"start_hour"`
	StartMinute  int      `json:"start_minute"`
	EndHour      int      `json:"end_hour"`
	EndMinute    int      `json:"end_minute"`
	TemperatureC *float64 `json:"temperature_c,omitempty"`
}

// MCPTimeProgramSlotResult represents the write outcome for one interval in the MCP output.
type MCPTimeProgramSlotResult struct {
	SlotIndex int    `json:"slot_index"`
	Accepted  bool   `json:"accepted"`
	ErrorCode int    `json:"error_code"`
	ErrorDesc string `json:"error_description,omitempty"`
}

// Function-module modes in MCP output.
const (
	MCPFunctionModuleSemanticModeInterpreted MCPFunctionModuleSemanticMode = "INTERPRETED"
	MCPFunctionModuleSemanticModeGPIOOnly    MCPFunctionModuleSemanticMode = "GPIO_ONLY"
	MCPFunctionModuleSemanticModeAbsent      MCPFunctionModuleSemanticMode = "ABSENT"
)

func (value MCPFunctionModuleInterpretation) validate() error {
	if strings.TrimSpace(value.EvidenceRevision) == "" {
		return errors.New("fm5 evidence revision is required")
	}
	if value.Mode == MCPFunctionModuleSemanticModeGPIOOnly {
		if value.DegradedReason == nil || *value.DegradedReason != "CONFIGURATION_NOT_INTERPRETABLE" {
			return errors.New("fm5 GPIO_ONLY requires CONFIGURATION_NOT_INTERPRETABLE")
		}
		return nil
	}
	if value.Mode != MCPFunctionModuleSemanticModeInterpreted && value.Mode != MCPFunctionModuleSemanticModeAbsent {
		return errors.New("invalid fm5 semantic mode")
	}
	if value.DegradedReason == nil {
		return nil
	}
	switch *value.DegradedReason {
	case "CONTROLLER_UNREACHABLE",
		"CONFIGURATION_UNAVAILABLE",
		"SOLAR_ACQUISITION_FAILED",
		"CYLINDER_ACQUISITION_FAILED",
		"EVIDENCE_STALE",
		"INCOHERENT_ACQUISITION":
		return nil
	default:
		return errors.New("fm5 retained structural mode requires a transient degraded reason")
	}
}
