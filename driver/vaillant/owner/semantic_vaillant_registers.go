package owner

const (
	vaillantExtRegisterPrimary   = byte(0xB5)
	vaillantExtRegisterSecondary = byte(0x24)
	vaillantB509Primary          = byte(0xB5)
	vaillantB509Secondary        = byte(0x09)
	vaillantB509OpcodeRead       = byte(0x0D)
	vaillantB509OpcodeWrite      = byte(0x0E)
	vaillantB524OpcodeRead       = byte(0x06)
	vaillantB524OpcodeLocal      = byte(0x02)
	vaillantB524OpRead           = byte(0x00)
	vaillantB524OpWrite          = byte(0x01)

	// --- Zone registers: GG=0x03 / OP=0x02 ---

	// STATE registers (read-only live values)
	zoneCurrentTemp      = uint16(0x000F) // state.current_room_temperature
	zoneSpecialFunction  = uint16(0x000E) // state.current_special_function
	zoneValveStatus      = uint16(0x0012) // state.valve_status
	zoneCurrentHumidity  = uint16(0x0028) // state.current_room_humidity
	zoneQuickVetoEndTime = uint16(0x001E) // state.quick_veto.end_time (HH:MM:SS)
	zoneQuickVetoEndDate = uint16(0x0024) // state.quick_veto.end_date (DD.MM.YY)

	// CONFIG registers (user-writable settings)
	zoneHeatingOpMode                         = uint16(0x0006) // configuration.heating.operation_mode
	zoneTargetTemp                            = uint16(0x0022) // configuration.heating.desired_setpoint
	zoneFallbackManualTemp                    = uint16(0x0014) // configuration.heating.manual_mode_setpoint
	zoneRoomTemperatureZoneMappingRawRegister = uint16(0x0013) // configuration.room_temperature_zone_mapping
	zoneQuickVetoTemp                         = uint16(0x0008) // configuration.quick_veto.temperature (f32 LE)
	zoneQuickVetoDuration                     = uint16(0x0026) // configuration.quick_veto.duration_hours (f32 LE)
	zoneHolidayStartDate                      = uint16(0x0003) // configuration.holiday.start_date (DD.MM.YY)
	zoneHolidayEndDate                        = uint16(0x0004) // configuration.holiday.end_date (DD.MM.YY)
	zoneHolidaySetpoint                       = uint16(0x0005) // configuration.holiday.setpoint (f32 LE)
	zoneHolidayEndTime                        = uint16(0x0020) // configuration.holiday.end_time (HH:MM)
	zoneHolidayStartTime                      = uint16(0x0021) // configuration.holiday.start_time (HH:MM)

	// PARAMS registers (identification / metadata)
	zoneName       = uint16(0x0016)
	zoneNamePrefix = uint16(0x0017)
	zoneNameSuffix = uint16(0x0018)
	zoneIndex      = uint16(0x001C)

	// --- Circuit registers: GG=0x02 / OP=0x02 ---

	// STATE registers (read-only live values)
	circuitFlowSetpoint  = uint16(0x0007) // flow_setpoint
	circuitFlowTemp      = uint16(0x0008) // flow_temperature (VF[x])
	circuitPumpStatus    = uint16(0x001B) // pump_status (observed Vaillant "Status Pump")
	circuitCircuitState  = uint16(0x001E) // circuit_state (observed Vaillant "Status circuit"; 0=standby 1=heating 2=cooling 3=dhw)
	circuitCalcFlowTemp  = uint16(0x0020) // legacy selector name; circuit_unknown_10 raw evidence, not a temperature
	circuitMixerPosition = uint16(0x0021) // mixer_position_pct
	circuitHumidity      = uint16(0x0022) // room_humidity_pct
	circuitDewPoint      = uint16(0x0023) // dew_point_temperature
	circuitPumpHours     = uint16(0x0024) // pump_operating_hours
	circuitPumpStarts    = uint16(0x0025) // pump_starts

	// CONFIG registers (user-writable settings)
	circuitTypeRegister    = uint16(0x0002) // configuration.heating_circuit_type / mixer_circuit_type_external
	circuitCoolingEnabled  = uint16(0x0006) // cooling_enabled
	circuitHeatingCurve    = uint16(0x000F) // heating_curve
	circuitFlowTempMax     = uint16(0x0010) // flow_temperature_max
	circuitFlowTempMin     = uint16(0x0012) // flow_temperature_min
	circuitSummerLimit     = uint16(0x0014) // summer_limit
	circuitRoomTempControl = uint16(0x0015) // room_temperature_control_mode
	// circuitSpecialFunctionMode is one unsigned byte. The name is
	// provisional; the meaning and value table are unconfirmed on the wire,
	// so the raw byte is published and no write is offered.
	circuitSpecialFunctionMode = uint16(0x001D)

	dhwPseudoCircuitInstance = byte(0x09)
	dhwPseudoCircuitTypeRaw  = uint16(0x0003)

	// CircuitStateStandby is the circuit state standby value of the Vaillant owner.
	CircuitStateStandby = "standby"
	// CircuitStateHeating is the circuit state heating value of the Vaillant owner.
	CircuitStateHeating = "heating"
	// CircuitStateCooling is the circuit state cooling value of the Vaillant owner.
	CircuitStateCooling = "cooling"
	// CircuitStateDHW is the circuit state dhw value of the Vaillant owner.
	CircuitStateDHW = "dhw"

	// --- DHW registers: GG=0x01 / OP=0x02 ---

	// STATE registers (read-only live values)
	dhwCurrentTemp     = uint16(0x0005) // state.current_dhw_temperature
	dhwSpecialFunction = uint16(0x000D) // state.current_special_function

	// CONFIG registers (user-writable settings)
	dhwOperationMode            = uint16(0x0003) // configuration.domestic_hot_water.operation_mode
	dhwTargetTemp               = uint16(0x0004) // configuration.domestic_hot_water.tapping_setpoint
	dhwHolidayStartDateRegister = uint16(0x0009) // configuration.holiday.start_date (DD.MM.YY)
	dhwHolidayEndDateRegister   = uint16(0x000A) // configuration.holiday.end_date (DD.MM.YY)

	dhwInstance = byte(0x00)

	// --- Device slot registers: OP=0x06 (remote) ---

	// STATE registers (read-only live values)
	deviceSlotConnected         = uint16(0x0001)
	deviceSlotRoomHumidity      = uint16(0x0007)
	deviceSlotRoomTemperature   = uint16(0x000F)
	deviceSlotPaired            = uint16(0x001E)
	deviceSlotReceptionStrength = uint16(0x001F)

	// PARAMS registers (identification / metadata)
	deviceSlotClassAddress         = uint16(0x0002)
	deviceSlotFirmware             = uint16(0x0004)
	deviceSlotRemoteControlAddress = uint16(0x0019)
	deviceSlotHardwareIdentifier   = uint16(0x0023)
	deviceSlotZoneAssignment       = uint16(0x0025)

	solarInstance = byte(0x00)

	// --- Solar registers: GG=0x04 / OP=0x02 ---

	// STATE registers (read-only live values)
	solarCollectorTemp = uint16(0x0003)
	solarReturnTemp    = uint16(0x0007)
	solarPumpActive    = uint16(0x0008)
	solarCurrentYield  = uint16(0x0009)
	solarPumpHours     = uint16(0x000B)

	// CONFIG registers (user-writable settings)
	solarEnabled      = uint16(0x0001)
	solarFunctionMode = uint16(0x0002)

	// --- Cylinder registers: GG=0x05 / OP=0x02 ---

	// STATE registers (read-only live values)
	cylinderTemperature = uint16(0x0004)

	// CONFIG registers (user-writable settings)
	cylinderMaxSetpoint      = uint16(0x0001)
	cylinderChargeHysteresis = uint16(0x0002)
	cylinderChargeOffset     = uint16(0x0003)
)

// b524GroupDef binds a B524 group number to its owning opcode.
// This eliminates the possibility of accidentally mixing OP/GG pairs.
// Public contract:
// https://github.com/Project-Helianthus/helianthus-docs-ebus/blob/main/architecture/b524-namespace-invariants.md
type b524GroupDef struct {
	group  byte
	opcode byte
	name   string
}

// OP=0x02 (local/controller) groups.
var (
	localRegulator = b524GroupDef{group: 0x00, opcode: vaillantB524OpcodeLocal, name: "regulator_parameters"}
	localDHW       = b524GroupDef{group: 0x01, opcode: vaillantB524OpcodeLocal, name: "hot_water_circuit"}
	localCircuits  = b524GroupDef{group: 0x02, opcode: vaillantB524OpcodeLocal, name: "heating_circuits"}
	localZones     = b524GroupDef{group: 0x03, opcode: vaillantB524OpcodeLocal, name: "zones"}
	localSolar     = b524GroupDef{group: 0x04, opcode: vaillantB524OpcodeLocal, name: "solar_circuit"}
	localCylinders = b524GroupDef{group: 0x05, opcode: vaillantB524OpcodeLocal, name: "hot_water_cylinder"}
)

// OP=0x06 (remote/device) groups — all gated on device_connected (RR=0x0001).
var (
	remoteRegulators        = b524GroupDef{group: 0x09, opcode: vaillantB524OpcodeRead, name: "regulators"}
	remoteThermostats       = b524GroupDef{group: 0x0A, opcode: vaillantB524OpcodeRead, name: "thermostats"}
	remoteFunctionalModules = b524GroupDef{group: 0x0C, opcode: vaillantB524OpcodeRead, name: "functional_modules"}
)

// remoteDeviceGroups lists the OP=0x06 groups that the gateway actively polls.
var remoteDeviceGroups = []b524GroupDef{
	remoteRegulators,
	remoteThermostats,
	remoteFunctionalModules,
}

// B5.55 timer/schedule protocol constants.
const (
	vaillantB555Primary   = byte(0xB5)
	vaillantB555Secondary = byte(0x55)

	b555OpcodeConfigRead = byte(0xA3)
	b555OpcodeSlotsRead  = byte(0xA4)
	b555OpcodeTimerRead  = byte(0xA5)
	b555OpcodeTimerWrite = byte(0xA6)

	b555HCHeating = byte(0x00)
	b555HCCooling = byte(0x01)
	b555HCDHW     = byte(0x02)
	b555HCCC      = byte(0x03)
	b555HCSilent  = byte(0x04)
	b555ZoneDHW   = byte(0xFF) // zone-agnostic selector for DHW/CC/Silent
)

var b555WeekdayNames = [7]string{"monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"}

var b555HCNames = map[byte]string{
	b555HCHeating: "heating",
	b555HCCooling: "cooling",
	b555HCDHW:     "dhw",
	b555HCCC:      "circulation",
	b555HCSilent:  "silent",
}

// --- Boiler status ---

// B524 registers on the controller that mirror boiler operational data.
// Group 0x00 (regulator parameters), instance 0x00.
const (
	regulatorInstance = byte(0x00)

	// --- System registers: GG=0x00 / OP=0x02 ---

	// STATE registers (read-only live values)
	systemOff                            = uint16(0x0007)
	systemWaterPressure                  = uint16(0x0039)
	systemFlowTemperature                = uint16(0x004B)
	systemOutdoorTemperature             = uint16(0x0073)
	systemOutdoorTemperatureDailyAverage = uint16(0x0095)
	systemMaintenanceDue                 = uint16(0x0096)
	systemHwcCylinderTemperatureTop      = uint16(0x009D)
	systemHwcCylinderTemperatureBottom   = uint16(0x009E)

	// CONFIG registers (user-writable settings)
	systemAdaptiveHeatingCurve         = uint16(0x0014)
	systemAlternativePoint             = uint16(0x0022)
	systemHeatingCircuitBivalencePoint = uint16(0x0023)
	systemDhwBivalencePoint            = uint16(0x0001)
	systemHcEmergencyTemperature       = uint16(0x0026)
	systemHwcMaxFlowTempDesired        = uint16(0x0046)
	systemMaxRoomHumidity              = uint16(0x000E)
	systemMaintenanceDate              = uint16(0x002C)
	systemInstallerName1               = uint16(0x006C)
	systemInstallerName2               = uint16(0x006D)
	systemInstallerPhone1              = uint16(0x006F)
	systemInstallerPhone2              = uint16(0x0070)
	systemInstallerMenuCode            = uint16(0x0076)

	// PARAMS registers (system topology / identification)
	systemScheme                            = uint16(0x0036)
	systemModuleConfigurationFunctionModule = uint16(0x002F)

	// --- Boiler B509 direct registers on BAI00 ---
	boilerB509WaterPressure         = uint16(0x0200)
	boilerB509FlameActive           = uint16(0x0500)
	boilerB509PartloadHcKw          = uint16(0x0704)
	boilerB509PartloadHwcKw         = uint16(0x0804)
	boilerB509FlowsetHcMaxC         = uint16(0x0E04)
	boilerB509FlowsetHcMaxCFallback = uint16(0xA500)
	boilerB509FlowsetHwcMaxC        = uint16(0x0F04)
	boilerB509PumpHours             = uint16(0x1400)
	boilerB509FlowTemperature       = uint16(0x1800)
	boilerB509FanHours              = uint16(0x1B00)
	boilerB509DeactivationsIfc      = uint16(0x1F00)
	boilerB509HoursTillService      = uint16(0x2004)
	boilerB509DeactivationsLimit    = uint16(0x2000)
	boilerB509DhwHours              = uint16(0x2200)
	boilerB509DhwStarts             = uint16(0x2300)
	boilerB509TargetFanSpeedRpm     = uint16(0x2400)
	boilerB509CentralHeatingHours   = uint16(0x2800)
	boilerB509CentralHeatingStarts  = uint16(0x2900)
	boilerB509ModulationPct         = uint16(0x2E00)
	boilerB509FlowTempDesiredC      = uint16(0x3900)
	boilerB509ExternalPumpActive    = uint16(0x3F00)
	boilerB509InstallerMenuCode     = uint16(0x4904)
	boilerB509CentralHeatingPump    = uint16(0x4400)
	boilerB509DiverterValvePosition = uint16(0x5400)
	boilerB509DhwWaterFlowLpm       = uint16(0x5500)
	boilerB509DhwDemandActive       = uint16(0x5800)
	boilerB509CirculationPumpActive = uint16(0x7B00)
	boilerB509PhoneNumber           = uint16(0x8104)
	boilerB509FanSpeedRpm           = uint16(0x8300)
	boilerB509StorageLoadPumpPct    = uint16(0x9E00)
	boilerB509IonisationVoltageUa   = uint16(0xA400)
	boilerB509StateNumber           = uint16(0xAB00)
	boilerB509GasValveActive        = uint16(0xBB00)
	boilerB509HeatingSwitchActive   = uint16(0xF203)
	boilerB509DhwTempDesiredC       = uint16(0xEA03)
	boilerB509PrimaryCircuitFlowLpm = uint16(0xFB00)
)
