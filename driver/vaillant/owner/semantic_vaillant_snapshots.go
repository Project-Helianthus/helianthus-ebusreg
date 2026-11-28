package owner

import (
	"time"

	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/owner/systemconfig"
)

type zonePresenceState string

const (
	zonePresenceAbsent           zonePresenceState = "ABSENT"
	zonePresenceSuspectResurrect zonePresenceState = "SUSPECT_RESURRECT"
	zonePresencePresent          zonePresenceState = "PRESENT"
	zonePresenceSuspectMissing   zonePresenceState = "SUSPECT_MISSING"
)

type zonePresenceRecord struct {
	State      zonePresenceState
	HitStreak  int
	MissStreak int
}

type vaillantZoneSnapshot struct {
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
	QuickVetoEndTime   string // "HH:MM" or "" if no veto
	QuickVetoEndDate   string // "YYYY-MM-DD" or "" if no veto

	HolidayStartDate string   // "YYYY-MM-DD" or ""
	HolidayEndDate   string   // "YYYY-MM-DD" or ""
	HolidaySetpointC *float64 // configured setpoint during holiday
	HolidayStartTime string   // "HH:MM" or ""
	HolidayEndTime   string   // "HH:MM" or ""

	FieldFreshness map[semanticFieldKey]semanticFieldFreshness
}

type vaillantDhwSnapshot struct {
	OperatingMode string
	Preset        string

	CurrentTempC *float64
	TargetTempC  *float64

	ConfigurationDHWOperationMode string
	StateSpecialFunction          string

	HolidayStartDate string // "YYYY-MM-DD" or ""
	HolidayEndDate   string // "YYYY-MM-DD" or ""

	FieldFreshness map[semanticFieldKey]semanticFieldFreshness
}

type vaillantCircuitSnapshot struct {
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
	// SpecialFunctionMode is the raw RR001D byte (provisional name).
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

type radioDeviceKey struct {
	Group    byte
	Instance byte
}

type vaillantRadioDeviceSnapshot struct {
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

// vaillantZoneView is the native view of one zone that the acquisition side
// hands to the publication owner. snapshot is a detached copy of the retained
// zone. The radio room values are taken from the mapped radio device only when
// the zone lacks its own temperature or humidity, and inferredRoomMappingRaw
// only when the zone carries no room temperature mapping.
type vaillantZoneView struct {
	instance               byte
	snapshot               *vaillantZoneSnapshot
	radioRoomTempC         *float64
	radioRoomHumidityPct   *float64
	inferredRoomMappingRaw *uint16
}

// vaillantDhwView is the native view of hot water handed to the publication
// owner: a detached copy of the retained snapshot, nil when none is retained,
// and whether this publication expired a stale snapshot.
type vaillantDhwView struct {
	snapshot *vaillantDhwSnapshot
	expired  bool
}

// vaillantCircuitsView is the native view of the active circuits handed to
// the publication owner, in instance order, with whether the function module
// manages them.
type vaillantCircuitsView struct {
	circuits                []*vaillantCircuitSnapshot
	managedByFunctionModule bool
}

// vaillantSystemView is the native view of the system handed to the
// publication owner: a detached copy of the retained snapshot, nil when none
// is retained, and the hot water bivalence bound resolved for the discovered
// controller, set only with a snapshot.
type vaillantSystemView struct {
	snapshot          *vaillantSystemSnapshot
	dhwBivalenceBound systemconfig.ConfigBound
}

type vaillantSolarSnapshot struct {
	CollectorTemperatureC *float64
	ReturnTemperatureC    *float64
	PumpActive            *bool
	CurrentYield          *float64
	PumpHours             *uint32
	SolarEnabled          *bool
	FunctionMode          *bool
}

type vaillantCylinderSnapshot struct {
	Instance         byte
	TemperatureC     *float64
	MaxSetpointC     *float64
	ChargeHysteresis *float64
	ChargeOffset     *float64
}

type semanticSnapshotSource uint8

const (
	semanticSnapshotSourceCache semanticSnapshotSource = iota
	semanticSnapshotSourceLive
)

type semanticFieldKey string

const (
	zoneFieldName                          semanticFieldKey = "zone.name"
	zoneFieldOperatingMode                 semanticFieldKey = "zone.operating_mode"
	zoneFieldPreset                        semanticFieldKey = "zone.preset"
	zoneFieldHvacAction                    semanticFieldKey = "zone.hvac_action"
	zoneFieldAllowedModes                  semanticFieldKey = "zone.allowed_modes"
	zoneFieldCurrentTempC                  semanticFieldKey = "zone.current_temp_c"
	zoneFieldTargetTempC                   semanticFieldKey = "zone.target_temp_c"
	zoneFieldCurrentHumidityPct            semanticFieldKey = "zone.current_humidity_pct"
	zoneFieldSpecialFunctionRaw            semanticFieldKey = "zone.special_function_raw"
	zoneFieldZoneOperationModeRaw          semanticFieldKey = "zone.operation_mode_raw"
	zoneFieldRoomTemperatureZoneMappingRaw semanticFieldKey = "zone.room_temperature_zone_mapping_raw"
	zoneFieldZoneCircuitIndexRaw           semanticFieldKey = "zone.circuit_index_raw"
	zoneFieldCircuitTypeRaw                semanticFieldKey = "zone.circuit_type_raw"
	zoneFieldZoneValveStatusRaw            semanticFieldKey = "zone.valve_status_raw"
	zoneFieldQuickVetoTempC                semanticFieldKey = "zone.quick_veto_temp_c"
	zoneFieldQuickVetoDurationH            semanticFieldKey = "zone.quick_veto_duration_h"
	zoneFieldQuickVetoEndTime              semanticFieldKey = "zone.quick_veto_end_time"
	zoneFieldQuickVetoEndDate              semanticFieldKey = "zone.quick_veto_end_date"
	zoneFieldHolidayStartDate              semanticFieldKey = "zone.holiday_start_date"
	zoneFieldHolidayEndDate                semanticFieldKey = "zone.holiday_end_date"
	zoneFieldHolidaySetpointC              semanticFieldKey = "zone.holiday_setpoint_c"
	zoneFieldHolidayStartTime              semanticFieldKey = "zone.holiday_start_time"
	zoneFieldHolidayEndTime                semanticFieldKey = "zone.holiday_end_time"
	dhwFieldOperatingMode                  semanticFieldKey = "dhw.operating_mode"
	dhwFieldPreset                         semanticFieldKey = "dhw.preset"
	dhwFieldCurrentTempC                   semanticFieldKey = "dhw.current_temp_c"
	dhwFieldTargetTempC                    semanticFieldKey = "dhw.target_temp_c"
	dhwFieldSpecialFunctionRaw             semanticFieldKey = "dhw.special_function_raw"
	dhwFieldDhwOperationModeRaw            semanticFieldKey = "dhw.operation_mode_raw"
	dhwFieldHolidayStartDate               semanticFieldKey = "dhw.holiday_start_date"
	dhwFieldHolidayEndDate                 semanticFieldKey = "dhw.holiday_end_date"
)

type semanticFieldFreshness struct {
	Source semanticSnapshotSource
	Stale  bool
}

type semanticFieldSet map[semanticFieldKey]struct{}

func newSemanticFieldSet(keys ...semanticFieldKey) semanticFieldSet {
	set := make(semanticFieldSet, len(keys))
	for _, key := range keys {
		set[key] = struct{}{}
	}
	return set
}

func (set semanticFieldSet) has(key semanticFieldKey) bool {
	if len(set) == 0 {
		return false
	}
	_, ok := set[key]
	return ok
}

type boilerStatusTier uint8

const (
	boilerStatusTierFast boilerStatusTier = iota
	boilerStatusTierMedium
	boilerStatusTierSlow
)

type boilerStatusTierSchedule struct {
	tier     boilerStatusTier
	interval time.Duration
	priority semanticTaskPriority
}

const (
	semanticTaskRefreshRegulatorCapability semanticTaskKey = "refresh_regulator_capability"
	semanticTaskRefreshDiscovery           semanticTaskKey = "refresh_discovery"
	semanticTaskRefreshConfig              semanticTaskKey = "refresh_config"
	semanticTaskRefreshState               semanticTaskKey = "refresh_state"
	semanticTaskRefreshCircuits            semanticTaskKey = "refresh_circuits"
	semanticTaskRefreshSystem              semanticTaskKey = "refresh_system"
	semanticTaskRefreshRadioDevices        semanticTaskKey = "refresh_radio_devices"
	semanticTaskRefreshEnergy              semanticTaskKey = "refresh_energy"
	semanticTaskRefreshSchedules           semanticTaskKey = "refresh_schedules"
	semanticTaskRefreshBoilerFast          semanticTaskKey = "refresh_boiler_fast"
	semanticTaskRefreshBoilerMedium        semanticTaskKey = "refresh_boiler_medium"
	semanticTaskRefreshBoilerSlow          semanticTaskKey = "refresh_boiler_slow"
	semanticTaskRefreshCurrentError        semanticTaskKey = "refresh_b503_current_error"
)

var (
	zoneConfigRefreshFieldSet = newSemanticFieldSet(
		zoneFieldName,
		zoneFieldOperatingMode,
		zoneFieldPreset,
		zoneFieldAllowedModes,
		zoneFieldTargetTempC,
		zoneFieldCurrentHumidityPct,
		zoneFieldZoneOperationModeRaw,
		zoneFieldRoomTemperatureZoneMappingRaw,
		zoneFieldZoneCircuitIndexRaw,
		zoneFieldCircuitTypeRaw,
		zoneFieldQuickVetoTempC,
		zoneFieldQuickVetoDurationH,
		zoneFieldQuickVetoEndTime,
		zoneFieldQuickVetoEndDate,
		zoneFieldHolidayStartDate,
		zoneFieldHolidayEndDate,
		zoneFieldHolidaySetpointC,
		zoneFieldHolidayStartTime,
		zoneFieldHolidayEndTime,
	)
	zoneFastStateFieldSet = newSemanticFieldSet(
		zoneFieldOperatingMode,
		zoneFieldPreset,
		zoneFieldHvacAction,
		zoneFieldAllowedModes,
		zoneFieldCurrentTempC,
		zoneFieldSpecialFunctionRaw,
		zoneFieldZoneValveStatusRaw,
	)
	zoneGrabFieldSet = newSemanticFieldSet(
		zoneFieldName,
		zoneFieldOperatingMode,
		zoneFieldPreset,
		zoneFieldHvacAction,
		zoneFieldAllowedModes,
		zoneFieldCurrentTempC,
		zoneFieldTargetTempC,
		zoneFieldCurrentHumidityPct,
		zoneFieldZoneOperationModeRaw,
		zoneFieldSpecialFunctionRaw,
		zoneFieldRoomTemperatureZoneMappingRaw,
		zoneFieldZoneCircuitIndexRaw,
		zoneFieldZoneValveStatusRaw,
		zoneFieldCircuitTypeRaw,
	)
	dhwFieldSet = newSemanticFieldSet(
		dhwFieldOperatingMode,
		dhwFieldPreset,
		dhwFieldCurrentTempC,
		dhwFieldTargetTempC,
		dhwFieldDhwOperationModeRaw,
		dhwFieldSpecialFunctionRaw,
		dhwFieldHolidayStartDate,
		dhwFieldHolidayEndDate,
	)
)
