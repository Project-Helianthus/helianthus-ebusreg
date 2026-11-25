package output

import (
	"github.com/Project-Helianthus/helianthus-ebusgo/transport"
	"strings"
	"time"
)

// NativeZoneSnapshot represents retained native zone values.
type NativeZoneSnapshot struct {
	Instance byte
	Present  bool

	Name string

	OperatingMode string
	Preset        string
	HvacAction    string
	AllowedModes  []string

	CurrentTempC *float64
	TargetTempC  *float64
	HumidityPct  *float64

	ConfigurationHeatingOperationMode          string
	StateSpecialFunction                       string
	ConfigurationRoomTemperatureZoneMappingRaw *uint16
	ConfigurationAssociatedCircuitRaw          *uint16
	ConfigurationCircuitTypeRaw                *uint16
	StateValveStatusRaw                        *uint16

	QuickVetoTempC     *float64
	QuickVetoDurationH *float64
	QuickVetoEndTime   string
	QuickVetoEndDate   string

	HolidayStartDate string
	HolidayEndDate   string
	HolidaySetpointC *float64
	HolidayStartTime string
	HolidayEndTime   string

	FieldFreshness map[SemanticFieldKey]SemanticFieldFreshness
}

// NativeDhwSnapshot represents retained native hot water values.
type NativeDhwSnapshot struct {
	OperatingMode string
	Preset        string

	CurrentTempC *float64
	TargetTempC  *float64

	ConfigurationDHWOperationMode string
	StateSpecialFunction          string

	HolidayStartDate string
	HolidayEndDate   string

	FieldFreshness map[SemanticFieldKey]SemanticFieldFreshness
}

// NativeCircuitSnapshot represents retained native circuit values.
type NativeCircuitSnapshot struct {
	Instance   byte
	Active     bool
	Controller byte

	CircuitTypeRaw     *uint16
	CoolingEnabledRaw  *uint16
	FlowSetpointC      *float64
	FlowTemperatureC   *float64
	HeatingCurve       *float64
	FlowTempMaxC       *float64
	FlowTempMinC       *float64
	SummerLimitC       *float64
	RoomTempControlRaw *uint16
	CircuitStateRaw    *uint16
	CircuitStateLiveAt time.Time

	SpecialFunctionMode *uint8
	PumpStatusRaw       *uint16
	PumpStatusLiveAt    time.Time
	CalcFlowTempC       *float64
	MixerPositionPct    *float64
	HumidityPct         *float64
	DewPointC           *float64
	PumpHoursRaw        *uint32
	PumpStartsRaw       *uint32
}

// NativeRadioDeviceSnapshot represents retained native radio device values.
type NativeRadioDeviceSnapshot struct {
	Group                byte
	Instance             byte
	SlotMode             string
	DeviceConnected      *bool
	DeviceClassAddress   *uint8
	DeviceModel          string
	FirmwareVersion      *string
	HardwareIdentifier   *uint16
	RemoteControlAddress *uint8
	DevicePaired         *bool
	ReceptionStrength    *uint8
	ZoneAssignment       *uint8
	RoomTemperatureC     *float64
	RoomHumidityPct      *float64
}

// NativeZoneView represents the zone snapshot and optional radio-derived values.
type NativeZoneView struct {
	Instance               byte
	Snapshot               *NativeZoneSnapshot
	RadioRoomTempC         *float64
	RadioRoomHumidityPct   *float64
	InferredRoomMappingRaw *uint16
}

// NativeDhwView represents the hot water snapshot and whether it expired.
type NativeDhwView struct {
	Snapshot *NativeDhwSnapshot
	Expired  bool
}

// NativeCircuitsView represents the active circuits and their management context.
type NativeCircuitsView struct {
	Circuits                []*NativeCircuitSnapshot
	ManagedByFunctionModule bool
}

// NativeSystemView represents the retained system snapshot and its hot water bound.
type NativeSystemView struct {
	Snapshot          *NativeSystemSnapshot
	DhwBivalenceBound ConfigBound
}

// NativeSolarSnapshot represents retained native solar values.
type NativeSolarSnapshot struct {
	CollectorTemperatureC *float64
	ReturnTemperatureC    *float64
	PumpActive            *bool
	CurrentYield          *float64
	PumpHours             *uint32
	SolarEnabled          *bool
	FunctionMode          *bool
}

// NativeCylinderSnapshot represents retained native cylinder values.
type NativeCylinderSnapshot struct {
	Instance         byte
	TemperatureC     *float64
	MaxSetpointC     *float64
	ChargeHysteresis *float64
	ChargeOffset     *float64
}

// NativeCacheHydration represents the native values recovered from the cache.
type NativeCacheHydration struct {
	Zones       []*NativeZoneSnapshot
	Dhw         *NativeDhwSnapshot
	PersistedAt time.Time
	Boiler      *NativeBoilerSnapshot
}

// NativeZoneCacheOverlay represents the retained zone values used in place of published fallbacks.
type NativeZoneCacheOverlay struct {
	CurrentTempC         *float64
	HumidityPct          *float64
	RoomMappingRaw       *uint16
	AssociatedCircuitRaw *uint16
}

// NativeFunctionModuleView represents the acquisition verdict and retained solar and cylinder values.
type NativeFunctionModuleView struct {
	Verdict   NativeFunctionModuleInterpretation
	Solar     *NativeSolarSnapshot
	Cylinders map[byte]*NativeCylinderSnapshot
}

// NativeBoilerSnapshot represents retained native boiler values.
type NativeBoilerSnapshot struct {
	FlowTemperatureC         *float64
	ReturnTemperatureC       *float64
	CentralHeatingPumpActive *bool
	WaterPressureBar         *float64
	ExternalPumpActive       *bool
	CirculationPumpActive    *bool
	GasValveActive           *bool
	FlameActive              *bool
	DiverterValvePositionPct *float64
	FanSpeedRpm              *int
	TargetFanSpeedRpm        *int
	IonisationVoltageUa      *float64
	DhwWaterFlowLpm          *float64
	DhwDemandActive          *bool
	HeatingSwitchActive      *bool
	StorageLoadPumpPct       *float64
	ModulationPct            *float64
	PrimaryCircuitFlowLpm    *float64
	FlowTempDesiredC         *float64
	DhwTempDesiredC          *float64
	StateNumber              *int
	DhwTemperatureC          *float64
	DhwTargetTemperatureC    *float64
	DhwOperatingMode         *int
	FlowsetHcMaxC            *float64
	FlowsetHwcMaxC           *float64
	PartloadHcKW             *float64
	PartloadHwcKW            *float64
	HeatingStatusRaw         *int
	DhwStatusRaw             *int
	CentralHeatingHours      *float64
	DhwHours                 *float64
	CentralHeatingStarts     *int
	DhwStarts                *int
	PumpHours                *float64
	FanHours                 *float64
	DeactivationsIFC         *int
	DeactivationsTemplimiter *int
	InstallerMenuCode        *int
	PhoneNumber              *string
	HoursTillService         *int
}

// NativeSystemSnapshot represents retained native system values.
type NativeSystemSnapshot struct {
	Controller byte

	SystemOff                      *bool
	SystemWaterPressure            *float64
	SystemFlowTemperature          *float64
	SystemFlowTemperatureLiveAt    time.Time
	OutdoorTemperature             *float64
	OutdoorTemperatureDailyAverage *float64
	MaintenanceDue                 *bool
	HwcCylinderTemperatureTop      *float64
	HwcCylinderTemperatureBottom   *float64

	AdaptiveHeatingCurve         *bool
	AlternativePoint             *float64
	HeatingCircuitBivalencePoint *float64
	DhwBivalencePoint            *float64
	HcEmergencyTemperature       *float64
	HwcMaxFlowTempDesired        *float64
	MaxRoomHumidity              *uint16
	MaintenanceDate              *string
	InstallerName                *string
	InstallerPhone               *string
	InstallerMenuCode            *uint16

	SystemScheme                    *uint16
	FunctionModuleConfigurationCode *uint16
}

// NativeAdapterInfo represents the parsed adapter version and retained telemetry.
type NativeAdapterInfo struct {
	Identity       *transport.AdapterVersion
	HardwareID     string
	HardwareConfig string

	TemperatureC *float64
	SupplyMV     *int
	BusMaxDV     *int
	BusMinDV     *int
	ResetCause   *string
	ResetCode    *byte
	RestartCount *byte
	WifiRSSI     *int

	LastIdentity  time.Time
	LastTelemetry time.Time
}

// NativeFunctionModuleMode represents the function-module mode determined by acquisition.
type NativeFunctionModuleMode string

// NativeFunctionModuleDegradedReason represents the reason acquisition could not establish or retain a verdict.
type NativeFunctionModuleDegradedReason string

// NativeFunctionModuleInterpretation represents the structural acquisition verdict and its evidence reference.
type NativeFunctionModuleInterpretation struct {
	Mode             NativeFunctionModuleMode
	DegradedReason   NativeFunctionModuleDegradedReason
	EvidenceRevision string
}

// Structural modes used in acquisition.
const (
	NativeFunctionModuleModeInterpreted NativeFunctionModuleMode = "INTERPRETED"
	NativeFunctionModuleModeGPIOOnly    NativeFunctionModuleMode = "GPIO_ONLY"
	NativeFunctionModuleModeAbsent      NativeFunctionModuleMode = "ABSENT"
)

// Degradation reasons used in acquisition.
const (
	NativeFunctionModuleDegradedControllerUnreachable         NativeFunctionModuleDegradedReason = "CONTROLLER_UNREACHABLE"
	NativeFunctionModuleDegradedConfigurationUnavailable      NativeFunctionModuleDegradedReason = "CONFIGURATION_UNAVAILABLE"
	NativeFunctionModuleDegradedConfigurationNotInterpretable NativeFunctionModuleDegradedReason = "CONFIGURATION_NOT_INTERPRETABLE"
	NativeFunctionModuleDegradedSolarAcquisitionFailed        NativeFunctionModuleDegradedReason = "SOLAR_ACQUISITION_FAILED"
	NativeFunctionModuleDegradedCylinderAcquisitionFailed     NativeFunctionModuleDegradedReason = "CYLINDER_ACQUISITION_FAILED"
	NativeFunctionModuleDegradedEvidenceStale                 NativeFunctionModuleDegradedReason = "EVIDENCE_STALE"
	NativeFunctionModuleDegradedIncoherentAcquisition         NativeFunctionModuleDegradedReason = "INCOHERENT_ACQUISITION"
)

// Valid reports whether the native verdict can be published.
func (verdict NativeFunctionModuleInterpretation) Valid() bool {
	if strings.TrimSpace(verdict.EvidenceRevision) == "" {
		return false
	}
	switch verdict.Mode {
	case NativeFunctionModuleModeInterpreted, NativeFunctionModuleModeAbsent:
		return verdict.DegradedReason == "" || IsFunctionModuleTransientAcquisitionReason(verdict.DegradedReason)
	case NativeFunctionModuleModeGPIOOnly:
		return verdict.DegradedReason == NativeFunctionModuleDegradedConfigurationNotInterpretable
	default:
		return false
	}
}

// IsFunctionModuleTransientAcquisitionReason reports whether the acquisition reason is transient.
func IsFunctionModuleTransientAcquisitionReason(reason NativeFunctionModuleDegradedReason) bool {
	switch reason {
	case NativeFunctionModuleDegradedControllerUnreachable,
		NativeFunctionModuleDegradedConfigurationUnavailable,
		NativeFunctionModuleDegradedSolarAcquisitionFailed,
		NativeFunctionModuleDegradedCylinderAcquisitionFailed,
		NativeFunctionModuleDegradedEvidenceStale,
		NativeFunctionModuleDegradedIncoherentAcquisition:
		return true
	default:
		return false
	}
}

// SemanticSnapshotSource represents whether values came from live acquisition or cache.
type SemanticSnapshotSource uint8

// Origins of retained semantic snapshots.
const (
	SemanticSnapshotSourceCache SemanticSnapshotSource = iota
	semanticSnapshotSourceLive
)

// SemanticFieldKey represents one native field in the freshness map.
type SemanticFieldKey string

// Native zone and domestic hot water freshness keys.
const (
	ZoneFieldName                          SemanticFieldKey = "zone.name"
	ZoneFieldOperatingMode                 SemanticFieldKey = "zone.operating_mode"
	ZoneFieldPreset                        SemanticFieldKey = "zone.preset"
	ZoneFieldHvacAction                    SemanticFieldKey = "zone.hvac_action"
	ZoneFieldAllowedModes                  SemanticFieldKey = "zone.allowed_modes"
	ZoneFieldCurrentTempC                  SemanticFieldKey = "zone.current_temp_c"
	ZoneFieldTargetTempC                   SemanticFieldKey = "zone.target_temp_c"
	ZoneFieldCurrentHumidityPct            SemanticFieldKey = "zone.current_humidity_pct"
	ZoneFieldSpecialFunctionRaw            SemanticFieldKey = "zone.special_function_raw"
	ZoneFieldZoneOperationModeRaw          SemanticFieldKey = "zone.operation_mode_raw"
	ZoneFieldRoomTemperatureZoneMappingRaw SemanticFieldKey = "zone.room_temperature_zone_mapping_raw"
	ZoneFieldZoneCircuitIndexRaw           SemanticFieldKey = "zone.circuit_index_raw"
	ZoneFieldCircuitTypeRaw                SemanticFieldKey = "zone.circuit_type_raw"
	ZoneFieldZoneValveStatusRaw            SemanticFieldKey = "zone.valve_status_raw"
	ZoneFieldQuickVetoTempC                SemanticFieldKey = "zone.quick_veto_temp_c"
	ZoneFieldQuickVetoDurationH            SemanticFieldKey = "zone.quick_veto_duration_h"
	ZoneFieldQuickVetoEndTime              SemanticFieldKey = "zone.quick_veto_end_time"
	ZoneFieldQuickVetoEndDate              SemanticFieldKey = "zone.quick_veto_end_date"
	ZoneFieldHolidayStartDate              SemanticFieldKey = "zone.holiday_start_date"
	ZoneFieldHolidayEndDate                SemanticFieldKey = "zone.holiday_end_date"
	ZoneFieldHolidaySetpointC              SemanticFieldKey = "zone.holiday_setpoint_c"
	ZoneFieldHolidayStartTime              SemanticFieldKey = "zone.holiday_start_time"
	ZoneFieldHolidayEndTime                SemanticFieldKey = "zone.holiday_end_time"
	DhwFieldOperatingMode                  SemanticFieldKey = "dhw.operating_mode"
	DhwFieldPreset                         SemanticFieldKey = "dhw.preset"
	DhwFieldCurrentTempC                   SemanticFieldKey = "dhw.current_temp_c"
	DhwFieldTargetTempC                    SemanticFieldKey = "dhw.target_temp_c"
	DhwFieldSpecialFunctionRaw             SemanticFieldKey = "dhw.special_function_raw"
	DhwFieldDhwOperationModeRaw            SemanticFieldKey = "dhw.operation_mode_raw"
	DhwFieldHolidayStartDate               SemanticFieldKey = "dhw.holiday_start_date"
	DhwFieldHolidayEndDate                 SemanticFieldKey = "dhw.holiday_end_date"
)

// SemanticFieldFreshness represents the source and stale flag of a native field.
type SemanticFieldFreshness struct {
	Source SemanticSnapshotSource
	Stale  bool
}
