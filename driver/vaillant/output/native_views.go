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

// NativeFM5View represents the acquisition verdict and retained solar and cylinder values.
type NativeFM5View struct {
	Verdict   NativeFM5Interpretation
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

	SystemOff                    *bool
	SystemWaterPressure          *float64
	SystemFlowTemperature        *float64
	SystemFlowTemperatureLiveAt  time.Time
	OutdoorTemperature           *float64
	OutdoorTemperatureAvg24h     *float64
	MaintenanceDue               *bool
	HwcCylinderTemperatureTop    *float64
	HwcCylinderTemperatureBottom *float64

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

	SystemScheme            *uint16
	ModuleConfigurationVR71 *uint16
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

// NativeFM5Mode represents the function-module mode determined by acquisition.
type NativeFM5Mode string

// NativeFM5DegradedReason represents the reason acquisition could not establish or retain a verdict.
type NativeFM5DegradedReason string

// NativeFM5Interpretation represents the structural acquisition verdict and its evidence reference.
type NativeFM5Interpretation struct {
	Mode             NativeFM5Mode
	DegradedReason   NativeFM5DegradedReason
	EvidenceRevision string
}

// Structural modes used in acquisition.
const (
	NativeFM5ModeInterpreted NativeFM5Mode = "INTERPRETED"
	NativeFM5ModeGPIOOnly    NativeFM5Mode = "GPIO_ONLY"
	NativeFM5ModeAbsent      NativeFM5Mode = "ABSENT"
)

// Degradation reasons used in acquisition.
const (
	NativeFM5DegradedControllerUnreachable         NativeFM5DegradedReason = "CONTROLLER_UNREACHABLE"
	NativeFM5DegradedConfigurationUnavailable      NativeFM5DegradedReason = "CONFIGURATION_UNAVAILABLE"
	NativeFM5DegradedConfigurationNotInterpretable NativeFM5DegradedReason = "CONFIGURATION_NOT_INTERPRETABLE"
	NativeFM5DegradedSolarAcquisitionFailed        NativeFM5DegradedReason = "SOLAR_ACQUISITION_FAILED"
	NativeFM5DegradedCylinderAcquisitionFailed     NativeFM5DegradedReason = "CYLINDER_ACQUISITION_FAILED"
	NativeFM5DegradedEvidenceStale                 NativeFM5DegradedReason = "EVIDENCE_STALE"
	NativeFM5DegradedIncoherentAcquisition         NativeFM5DegradedReason = "INCOHERENT_ACQUISITION"
)

// Valid reports whether the native verdict can be published.
func (verdict NativeFM5Interpretation) Valid() bool {
	if strings.TrimSpace(verdict.EvidenceRevision) == "" {
		return false
	}
	switch verdict.Mode {
	case NativeFM5ModeInterpreted, NativeFM5ModeAbsent:
		return verdict.DegradedReason == "" || IsFM5TransientAcquisitionReason(verdict.DegradedReason)
	case NativeFM5ModeGPIOOnly:
		return verdict.DegradedReason == NativeFM5DegradedConfigurationNotInterpretable
	default:
		return false
	}
}

// IsFM5TransientAcquisitionReason reports whether the acquisition reason is transient.
func IsFM5TransientAcquisitionReason(reason NativeFM5DegradedReason) bool {
	switch reason {
	case NativeFM5DegradedControllerUnreachable,
		NativeFM5DegradedConfigurationUnavailable,
		NativeFM5DegradedSolarAcquisitionFailed,
		NativeFM5DegradedCylinderAcquisitionFailed,
		NativeFM5DegradedEvidenceStale,
		NativeFM5DegradedIncoherentAcquisition:
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
