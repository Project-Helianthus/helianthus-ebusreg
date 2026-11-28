package output

import (
	"errors"
	"strings"
	"time"
)

// Zone represents one space-heating zone.
type Zone struct {
	ID     string
	Name   string
	State  ZoneState
	Config ZoneConfig
}

// ZoneState represents the observed state of a zone.
type ZoneState struct {
	CurrentTempC       *float64
	CurrentHumidityPct *float64
	HvacAction         string
	SpecialFunction    string
	HeatingDemandPct   *float64
	ValvePositionPct   *float64
	Enumerations       []CatalogueEnumeration
}

// ZoneConfig represents the settings of a zone.
type ZoneConfig struct {
	OperatingMode              string
	OperationModeChangeable    *bool
	SourceLabel                *string
	Preset                     string
	TargetTempC                *float64
	AllowedModes               []string
	CircuitType                string
	AssociatedCircuit          *int
	RoomTemperatureZoneMapping *int
	QuickVeto                  bool
	QuickVetoSetpointC         *float64
	QuickVetoDurationH         *float64
	QuickVetoExpiry            string
	HolidayStartDate           string
	HolidayEndDate             string
	HolidaySetpointC           *float64
	HolidayStartTime           string
	HolidayEndTime             string
}

// DhwStatus represents domestic hot water state and settings.
type DhwStatus struct {
	State  DhwState
	Config DhwConfig
}

// DhwState represents the observed state of domestic hot water.
type DhwState struct {
	CurrentTempC     *float64
	OverrunActive    *bool
	SpecialFunction  string
	HeatingDemandPct *float64
	Enumerations     []CatalogueEnumeration
}

// DhwConfig represents domestic hot water settings.
type DhwConfig struct {
	OperatingMode           string
	OperationModeChangeable *bool
	Preset                  string
	TargetTempC             *float64
	HolidayStartDate        string
	HolidayEndDate          string
}

// CircuitStatus represents the state, settings, and owner of a circuit.
type CircuitStatus struct {
	Index          int
	CircuitType    string
	HasMixer       bool
	State          CircuitState
	Config         CircuitConfig
	ManagingDevice ManagingDevice
}

// CircuitState represents the observed state of a circuit.
type CircuitState struct {
	PumpActive       *bool
	MixerPositionPct *float64
	FlowTemperatureC *float64
	FlowSetpointC    *float64
	CalcFlowTempC    *float64
	CircuitState     string
	Humidity         *float64
	DewPoint         *float64
	PumpHours        *float64
	PumpStarts       *int
}

// CircuitConfig represents the settings of a circuit.
type CircuitConfig struct {
	HeatingCurve *float64
	FlowTempMaxC *float64
	FlowTempMinC *float64
	SummerLimitC *float64

	SpecialFunctionMode *int
	RoomTempControl     string
	CoolingEnabled      *bool
	Enumerations        []CatalogueEnumeration
}

// ManagingDevice represents the device responsible for a circuit.
type ManagingDevice struct {
	Role     ManagingDeviceRole
	DeviceID *string
	Address  *int
}

// ManagingDeviceRole represents the role of the device responsible for a circuit.
type ManagingDeviceRole string

// RadioDevice represents one attached radio device.
type RadioDevice struct {
	Group                int
	Instance             int
	SlotMode             string
	DeviceConnected      *bool
	DeviceClassAddress   *int
	DeviceModel          string
	FirmwareVersion      *string
	HardwareIdentifier   *int
	RemoteControlAddress *int
	DevicePaired         *bool
	ReceptionStrength    *int
	ZoneAssignment       *int
	RoomTemperatureC     *float64
	RoomHumidityPct      *float64
}

// FunctionModuleInterpretation represents the structural result and its evidence reference.
type FunctionModuleInterpretation struct {
	Mode             FunctionModuleSemanticMode
	DegradedReason   FunctionModuleSemanticDegradedReason
	EvidenceRevision string
}

// FunctionModuleSemanticMode represents the structural mode of the function module.
type FunctionModuleSemanticMode string

// FunctionModuleSemanticDegradedReason represents the reason a function-module interpretation is degraded.
type FunctionModuleSemanticDegradedReason string

// SolarStatus represents the solar plane.
type SolarStatus struct {
	CollectorTemperatureC *float64
	ReturnTemperatureC    *float64
	PumpActive            *bool
	CurrentYield          *float64
	PumpHours             *float64
	SolarEnabled          *bool
	FunctionMode          *bool
}

// CylinderStatus represents one hot water cylinder.
type CylinderStatus struct {
	Index             int
	TemperatureC      *float64
	MaxSetpointC      *float64
	ChargeHysteresisC *float64
	ChargeOffsetC     *float64
}

// EnergySeries represents the energy values and their metadata for one use.
type EnergySeries struct {
	Today       float64
	Yearly      []float64
	Monthly     []float64
	TodayMeta   EnergyPointMeta
	YearlyMeta  []EnergyPointMeta
	MonthlyMeta []EnergyPointMeta
}

// EnergyChannel represents climate and hot water energy for one source.
type EnergyChannel struct {
	DHW     EnergySeries
	Climate EnergySeries
}

// EnergyTotals represents the totals for the supported energy sources.
type EnergyTotals struct {
	Gas                 EnergyChannel
	Electric            EnergyChannel
	Solar               EnergyChannel
	RegisterDirections  []EnergyDirection
	CatalogueDirections []EnergyDirection
	ResidualFields      []EnergyResidual
	CandidateFields     []EnergyResidual
}

// EnergyFreshnessState represents the freshness classification of an energy value.
type EnergyFreshnessState string

// EnergyProvenance represents the provenance classification of an energy value.
type EnergyProvenance string

// EnergyPointMeta represents the freshness and provenance of one energy value.
type EnergyPointMeta struct {
	FreshnessState  EnergyFreshnessState
	Provenance      EnergyProvenance
	LastObservedUTC string
	AgeSeconds      float64
	Stale           bool
}

// BoilerStatus represents the full boiler status.
type BoilerStatus struct {
	State       BoilerState
	Config      BoilerConfig
	Diagnostics BoilerDiagnostics
	Selectors   []BoilerSelectorEvidence
	Operations  []RawOperation
}

// BoilerState represents the observed state of the boiler.
type BoilerState struct {
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
}

// BoilerConfig represents the boiler settings.
type BoilerConfig struct {
	DhwOperatingMode  *string
	FlowsetHcMaxC     *float64
	FlowsetHwcMaxC    *float64
	PartloadHcKW      *float64
	PartloadHwcKW     *float64
	InstallerMenuCode *int
	PhoneNumber       *string
	HoursTillService  *int
}

// BoilerDiagnostics represents the boiler operating counters and diagnostic values.
type BoilerDiagnostics struct {
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
}

// SystemStatus represents the system status.
type SystemStatus struct {
	State          SystemState
	Config         SystemConfig
	Properties     SystemProperties
	GatewayBrand   *string
	GatewayVendor  *string
	Ventilation    *VentilationView
	Catalogue      *FamilyCatalogue
	history        *HistoryOwnerRecord
	functionModule *FunctionModuleRecord
	legacyPlanes   *LegacyPlaneRecord
}

// SystemState represents the observed state of the system.
type SystemState struct {
	SystemOff                      *bool
	SystemWaterPressure            *float64
	SystemFlowTemperature          *float64
	OutdoorTemperature             *float64
	OutdoorTemperatureDailyAverage *float64
	MaintenanceDue                 *bool
	HwcCylinderTemperatureTop      *float64
	HwcCylinderTemperatureBottom   *float64
}

// SystemConfig represents the system settings and active bound.
type SystemConfig struct {
	AdaptiveHeatingCurve         *bool
	AlternativePoint             *float64
	HeatingCircuitBivalencePoint *float64
	DhwBivalencePoint            *float64

	DhwBivalenceBound      *ConfigBound
	HcEmergencyTemperature *float64
	HwcMaxFlowTempDesired  *float64
	MaxRoomHumidity        *int
	MaintenanceDate        *string
	InstallerName          *string
	InstallerPhone         *string
	InstallerMenuCode      *int
}

// SystemProperties represents the hydraulic scheme and module settings.
type SystemProperties struct {
	SystemScheme                    *int
	FunctionModuleConfigurationCode *int
}

// ScheduleTimerSlot represents one scheduled interval.
type ScheduleTimerSlot struct {
	StartHour      int
	StartMinute    int
	EndHour        int
	EndMinute      int
	TemperatureC   *float64
	TemperatureRaw *int
}

// ScheduleDayProgram represents the intervals for one day.
type ScheduleDayProgram struct {
	Weekday string
	Slots   []ScheduleTimerSlot
}

// ScheduleConfig represents the native limits and temperature support of a program.
type ScheduleConfig struct {
	MaxSlots       int
	TimeResolution int
	MinDuration    int
	HasTemperature bool
	TempSlots      int
	MinTempC       *float64
	MaxTempC       *float64
}

// ScheduleProgram represents the program for one zone and circuit kind.
type ScheduleProgram struct {
	Zone      int
	HC        string
	Config    *ScheduleConfig
	SlotsUsed []int
	Days      []ScheduleDayProgram
}

// ScheduleStatus represents the set of programs.
type ScheduleStatus struct {
	Programs []ScheduleProgram
}

// AdapterHardwareInfo represents adapter identity and telemetry.
type AdapterHardwareInfo struct {
	FirmwareVersion    string
	FirmwareChecksum   string
	BootloaderVersion  string
	BootloaderChecksum string
	HardwareID         string
	HardwareConfig     string
	Features           byte
	Jumpers            byte
	JumperFlags        []string
	IsWiFi             bool
	IsEthernet         bool
	TemperatureC       *float64
	SupplyVoltageMV    *int
	BusVoltageMaxDV    *int
	BusVoltageMinDV    *int
	ResetCause         *string
	ResetCauseCode     *byte
	RestartCount       *byte
	WiFiRSSIDBm        *int
	LastIdentityQuery  *time.Time
	LastTelemetryQuery *time.Time
	VersionResponseLen int
	InfoSupported      bool
}

// VaillantRegulatorCapability represents the tri-state regulator capability.
type VaillantRegulatorCapability string

// Managing device roles in the circuit output.
const (
	ManagingDeviceRoleRegulator      ManagingDeviceRole = "REGULATOR"
	ManagingDeviceRoleFunctionModule ManagingDeviceRole = "FUNCTION_MODULE"
	ManagingDeviceRoleUnknown        ManagingDeviceRole = "UNKNOWN"
)

// Published regulator capability outcomes.
const (
	VaillantRegulatorCapabilityUnknown VaillantRegulatorCapability = "UNKNOWN"
	VaillantRegulatorCapabilityNone    VaillantRegulatorCapability = "NONE"
	VaillantRegulatorCapabilityPresent VaillantRegulatorCapability = "PRESENT"
)

// NormalizeVaillantRegulatorCapability returns unknown for values outside the capability vocabulary.
func NormalizeVaillantRegulatorCapability(capability VaillantRegulatorCapability) VaillantRegulatorCapability {
	switch capability {
	case VaillantRegulatorCapabilityNone, VaillantRegulatorCapabilityPresent, VaillantRegulatorCapabilityUnknown:
		return capability
	default:
		return VaillantRegulatorCapabilityUnknown
	}
}

// Published function-module modes.
const (
	FunctionModuleSemanticModeInterpreted FunctionModuleSemanticMode = "INTERPRETED"
	FunctionModuleSemanticModeGPIOOnly    FunctionModuleSemanticMode = "GPIO_ONLY"
	FunctionModuleSemanticModeAbsent      FunctionModuleSemanticMode = "ABSENT"
)

// Published reasons for degraded function-module interpretations.
const (
	FunctionModuleSemanticDegradedReasonControllerUnreachable         FunctionModuleSemanticDegradedReason = "CONTROLLER_UNREACHABLE"
	FunctionModuleSemanticDegradedReasonConfigurationUnavailable      FunctionModuleSemanticDegradedReason = "CONFIGURATION_UNAVAILABLE"
	FunctionModuleSemanticDegradedReasonConfigurationNotInterpretable FunctionModuleSemanticDegradedReason = "CONFIGURATION_NOT_INTERPRETABLE"
	FunctionModuleSemanticDegradedReasonSolarAcquisitionFailed        FunctionModuleSemanticDegradedReason = "SOLAR_ACQUISITION_FAILED"
	FunctionModuleSemanticDegradedReasonCylinderAcquisitionFailed     FunctionModuleSemanticDegradedReason = "CYLINDER_ACQUISITION_FAILED"
	FunctionModuleSemanticDegradedReasonEvidenceStale                 FunctionModuleSemanticDegradedReason = "EVIDENCE_STALE"
	FunctionModuleSemanticDegradedReasonIncoherentAcquisition         FunctionModuleSemanticDegradedReason = "INCOHERENT_ACQUISITION"
)

// Validate reports an incomplete or inconsistent interpretation.
func (value FunctionModuleInterpretation) Validate() error {
	if strings.TrimSpace(value.EvidenceRevision) == "" {
		return errors.New("fm5 evidence revision is required")
	}
	switch value.Mode {
	case FunctionModuleSemanticModeInterpreted, FunctionModuleSemanticModeAbsent:
		if value.DegradedReason != "" && !validFunctionModuleTransientDegradedReason(value.DegradedReason) {
			return errors.New("fm5 retained structural mode requires a transient degraded reason")
		}
	case FunctionModuleSemanticModeGPIOOnly:
		if value.DegradedReason != FunctionModuleSemanticDegradedReasonConfigurationNotInterpretable {
			return errors.New("fm5 GPIO_ONLY requires CONFIGURATION_NOT_INTERPRETABLE")
		}
	default:
		return errors.New("invalid fm5 semantic mode")
	}
	return nil
}

func validFunctionModuleSemanticDegradedReason(reason FunctionModuleSemanticDegradedReason) bool {
	switch reason {
	case FunctionModuleSemanticDegradedReasonControllerUnreachable,
		FunctionModuleSemanticDegradedReasonConfigurationUnavailable,
		FunctionModuleSemanticDegradedReasonConfigurationNotInterpretable,
		FunctionModuleSemanticDegradedReasonSolarAcquisitionFailed,
		FunctionModuleSemanticDegradedReasonCylinderAcquisitionFailed,
		FunctionModuleSemanticDegradedReasonEvidenceStale,
		FunctionModuleSemanticDegradedReasonIncoherentAcquisition:
		return true
	default:
		return false
	}
}

func validFunctionModuleTransientDegradedReason(reason FunctionModuleSemanticDegradedReason) bool {
	return validFunctionModuleSemanticDegradedReason(reason) &&
		reason != FunctionModuleSemanticDegradedReasonConfigurationNotInterpretable
}

// Freshness classifications for energy fields.
const (
	EnergyFreshnessStateNeverSeen   EnergyFreshnessState = "never_seen"
	EnergyFreshnessStateFresh       EnergyFreshnessState = "fresh"
	EnergyFreshnessStateWarmingUp   EnergyFreshnessState = "warming_up"
	EnergyFreshnessStateStale       EnergyFreshnessState = "stale"
	EnergyFreshnessStateUnavailable EnergyFreshnessState = "unavailable"
)

// Provenance classifications for energy fields.
const (
	EnergyProvenanceNone      EnergyProvenance = "none"
	EnergyProvenanceRegister  EnergyProvenance = "register"
	EnergyProvenanceBroadcast EnergyProvenance = "broadcast"
)

// ConfigBound represents the active bound and its target identity.
type ConfigBound struct {
	ControllerAddress          byte                       `json:"-"`
	ControllerRegistryIdentity ControllerRegistryIdentity `json:"-"`
	MinC                       float64
	MaxC                       float64
	Source                     string

	Reason string

	Profile string

	NarrowedOnly bool
}

// ControllerRegistryIdentity represents the controller identity bound to a registry observation.
type ControllerRegistryIdentity struct {
	Valid                                                    bool
	Address                                                  byte
	Generation                                               uint64
	Manufacturer, DeviceID, SoftwareVersion, HardwareVersion string
	SerialNumber, MacAddress                                 string
}
