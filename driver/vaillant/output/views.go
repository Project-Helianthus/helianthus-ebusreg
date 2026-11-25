package output

import (
	"fmt"
	"github.com/Project-Helianthus/helianthus-ebusreg/vaillant/productids"
	"reflect"
	"strings"
)

// ZoneFromView maps the native values to the corresponding output record.
func ZoneFromView(view NativeZoneView) Zone {
	entry := view.Snapshot
	name := entry.Name
	if strings.TrimSpace(name) == "" {
		name = fmt.Sprintf("Zone %d", view.Instance+1)
	}

	quickVetoActive := entry.Preset == "quickveto"
	var qvExpiry string
	if quickVetoActive && entry.QuickVetoEndDate != "" && entry.QuickVetoEndTime != "" {
		qvExpiry = entry.QuickVetoEndDate + "T" + entry.QuickVetoEndTime
	}
	currentTempC := entry.CurrentTempC
	if currentTempC == nil {
		currentTempC = view.RadioRoomTempC
	}
	currentHumidityPct := entry.HumidityPct
	if currentHumidityPct == nil {
		currentHumidityPct = view.RadioRoomHumidityPct
	}
	roomMappingRaw := entry.ConfigurationRoomTemperatureZoneMappingRaw
	if roomMappingRaw == nil {
		roomMappingRaw = view.InferredRoomMappingRaw
	}
	associatedCircuitRaw := entry.ConfigurationAssociatedCircuitRaw
	if associatedCircuitRaw == nil && roomMappingRaw != nil {
		value := uint16(ResolveAssociatedCircuitInstance(roomMappingRaw, entry.Instance))
		associatedCircuitRaw = &value
	}
	return Zone{
		ID:   fmt.Sprintf("zone-%d", view.Instance+1),
		Name: name,
		State: ZoneState{
			CurrentTempC:       currentTempC,
			CurrentHumidityPct: currentHumidityPct,
			HvacAction:         entry.HvacAction,
			SpecialFunction:    entry.StateSpecialFunction,
			ValvePositionPct:   decodeValvePositionPct(entry.StateValveStatusRaw),
		},
		Config: ZoneConfig{
			OperatingMode:              entry.OperatingMode,
			Preset:                     entry.Preset,
			TargetTempC:                entry.TargetTempC,
			AllowedModes:               append([]string(nil), entry.AllowedModes...),
			CircuitType:                decodeCircuitType(entry.ConfigurationCircuitTypeRaw),
			AssociatedCircuit:          decodeAssociatedCircuit(associatedCircuitRaw),
			RoomTemperatureZoneMapping: decodeRoomTemperatureZoneMapping(roomMappingRaw),
			QuickVeto:                  quickVetoActive,
			QuickVetoSetpointC:         entry.QuickVetoTempC,
			QuickVetoDurationH:         entry.QuickVetoDurationH,
			QuickVetoExpiry:            qvExpiry,
			HolidayStartDate:           entry.HolidayStartDate,
			HolidayEndDate:             entry.HolidayEndDate,
			HolidaySetpointC:           entry.HolidaySetpointC,
			HolidayStartTime:           entry.HolidayStartTime,
			HolidayEndTime:             entry.HolidayEndTime,
		},
	}
}

// DhwStatusFromSnapshot maps the native values to the corresponding output record.
func DhwStatusFromSnapshot(snapshot *NativeDhwSnapshot) *DhwStatus {
	return &DhwStatus{
		State: DhwState{
			CurrentTempC:    snapshot.CurrentTempC,
			SpecialFunction: snapshot.StateSpecialFunction,
		},
		Config: DhwConfig{
			OperatingMode:    snapshot.OperatingMode,
			Preset:           snapshot.Preset,
			TargetTempC:      snapshot.TargetTempC,
			HolidayStartDate: snapshot.HolidayStartDate,
			HolidayEndDate:   snapshot.HolidayEndDate,
		},
	}
}

// CircuitStatusFromSnapshot maps the native values to the corresponding output record.
func CircuitStatusFromSnapshot(snapshot *NativeCircuitSnapshot, managedByFunctionModule bool) CircuitStatus {
	return CircuitStatus{
		Index:       int(snapshot.Instance),
		CircuitType: decodeHeatingCircuitTypeToken(snapshot.CircuitTypeRaw),
		HasMixer:    snapshot.MixerPositionPct != nil,
		State: CircuitState{
			PumpActive:       decodeUint16Bool(snapshot.PumpStatusRaw),
			MixerPositionPct: CloneFloat64Ptr(snapshot.MixerPositionPct),
			FlowTemperatureC: CloneFloat64Ptr(snapshot.FlowTemperatureC),
			FlowSetpointC:    CloneFloat64Ptr(snapshot.FlowSetpointC),
			CalcFlowTempC:    CloneFloat64Ptr(snapshot.CalcFlowTempC),
			CircuitState:     decodeCircuitStateToken(snapshot.CircuitStateRaw),
			Humidity:         CloneFloat64Ptr(snapshot.HumidityPct),
			DewPoint:         CloneFloat64Ptr(snapshot.DewPointC),
			PumpHours:        decodeUint32Float64(snapshot.PumpHoursRaw),
			PumpStarts:       decodeUint32Int(snapshot.PumpStartsRaw),
		},
		Config: CircuitConfig{
			HeatingCurve:        CloneFloat64Ptr(snapshot.HeatingCurve),
			FlowTempMaxC:        CloneFloat64Ptr(snapshot.FlowTempMaxC),
			FlowTempMinC:        CloneFloat64Ptr(snapshot.FlowTempMinC),
			SummerLimitC:        CloneFloat64Ptr(snapshot.SummerLimitC),
			SpecialFunctionMode: decodeUint8Int(snapshot.SpecialFunctionMode),
			RoomTempControl:     decodeRoomTempControlToken(snapshot.RoomTempControlRaw),
			CoolingEnabled:      decodeUint16Bool(snapshot.CoolingEnabledRaw),
		},
		ManagingDevice: circuitManagingDevice(managedByFunctionModule),
	}
}

// RadioDeviceFromSnapshot maps the native values to the corresponding output record.
func RadioDeviceFromSnapshot(snapshot *NativeRadioDeviceSnapshot) RadioDevice {
	return RadioDevice{
		Group:                int(snapshot.Group),
		Instance:             int(snapshot.Instance),
		SlotMode:             snapshot.SlotMode,
		DeviceConnected:      cloneBoilerBoolPtr(snapshot.DeviceConnected),
		DeviceClassAddress:   uint8ToIntPtr(snapshot.DeviceClassAddress),
		DeviceModel:          snapshot.DeviceModel,
		FirmwareVersion:      CloneStringPtr(snapshot.FirmwareVersion),
		HardwareIdentifier:   uint16ToIntPtr(snapshot.HardwareIdentifier),
		RemoteControlAddress: uint8ToIntPtr(snapshot.RemoteControlAddress),
		DevicePaired:         cloneBoilerBoolPtr(snapshot.DevicePaired),
		ReceptionStrength:    uint8ToIntPtr(snapshot.ReceptionStrength),
		ZoneAssignment:       uint8ToIntPtr(snapshot.ZoneAssignment),
		RoomTemperatureC:     CloneFloat64Ptr(snapshot.RoomTemperatureC),
		RoomHumidityPct:      CloneFloat64Ptr(snapshot.RoomHumidityPct),
	}
}

// BoilerStatusFromSnapshot maps the native values to the corresponding output record.
func BoilerStatusFromSnapshot(snapshot *NativeBoilerSnapshot) *BoilerStatus {
	current := &BoilerStatus{
		State: BoilerState{
			FlowTemperatureC:         snapshot.FlowTemperatureC,
			ReturnTemperatureC:       snapshot.ReturnTemperatureC,
			CentralHeatingPumpActive: snapshot.CentralHeatingPumpActive,
			WaterPressureBar:         snapshot.WaterPressureBar,
			ExternalPumpActive:       snapshot.ExternalPumpActive,
			CirculationPumpActive:    snapshot.CirculationPumpActive,
			GasValveActive:           snapshot.GasValveActive,
			FlameActive:              snapshot.FlameActive,
			DiverterValvePositionPct: snapshot.DiverterValvePositionPct,
			FanSpeedRpm:              snapshot.FanSpeedRpm,
			TargetFanSpeedRpm:        snapshot.TargetFanSpeedRpm,
			IonisationVoltageUa:      snapshot.IonisationVoltageUa,
			DhwWaterFlowLpm:          snapshot.DhwWaterFlowLpm,
			DhwDemandActive:          snapshot.DhwDemandActive,
			HeatingSwitchActive:      snapshot.HeatingSwitchActive,
			StorageLoadPumpPct:       snapshot.StorageLoadPumpPct,
			ModulationPct:            snapshot.ModulationPct,
			PrimaryCircuitFlowLpm:    snapshot.PrimaryCircuitFlowLpm,
			FlowTempDesiredC:         snapshot.FlowTempDesiredC,
			DhwTempDesiredC:          snapshot.DhwTempDesiredC,
			StateNumber:              snapshot.StateNumber,
			DhwTemperatureC:          snapshot.DhwTemperatureC,
			DhwTargetTemperatureC:    snapshot.DhwTargetTemperatureC,
		},
		Config: BoilerConfig{
			FlowsetHcMaxC:     snapshot.FlowsetHcMaxC,
			FlowsetHwcMaxC:    snapshot.FlowsetHwcMaxC,
			PartloadHcKW:      snapshot.PartloadHcKW,
			PartloadHwcKW:     snapshot.PartloadHwcKW,
			InstallerMenuCode: snapshot.InstallerMenuCode,
			PhoneNumber:       CloneStringPtr(snapshot.PhoneNumber),
			HoursTillService:  snapshot.HoursTillService,
		},
		Diagnostics: BoilerDiagnostics{
			HeatingStatusRaw:         snapshot.HeatingStatusRaw,
			DhwStatusRaw:             snapshot.DhwStatusRaw,
			CentralHeatingHours:      snapshot.CentralHeatingHours,
			DhwHours:                 snapshot.DhwHours,
			CentralHeatingStarts:     snapshot.CentralHeatingStarts,
			DhwStarts:                snapshot.DhwStarts,
			PumpHours:                snapshot.PumpHours,
			FanHours:                 snapshot.FanHours,
			DeactivationsIFC:         snapshot.DeactivationsIFC,
			DeactivationsTemplimiter: snapshot.DeactivationsTemplimiter,
		},
	}
	if snapshot.DhwOperatingMode != nil {
		mode := DecodeDhwMode(*snapshot.DhwOperatingMode)
		current.Config.DhwOperatingMode = &mode
	}
	return current
}

// BoilerStatusEquals compares published boiler status values.
func BoilerStatusEquals(a, b *BoilerStatus) bool {
	return reflect.DeepEqual(a, b)
}

// SystemStatusEquals compares published system status values.
func SystemStatusEquals(a, b *SystemStatus) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return boolPtrEquals(a.State.SystemOff, b.State.SystemOff) &&
		floatPtrEquals(a.State.SystemWaterPressure, b.State.SystemWaterPressure) &&
		floatPtrEquals(a.State.SystemFlowTemperature, b.State.SystemFlowTemperature) &&
		floatPtrEquals(a.State.OutdoorTemperature, b.State.OutdoorTemperature) &&
		floatPtrEquals(a.State.OutdoorTemperatureDailyAverage, b.State.OutdoorTemperatureDailyAverage) &&
		boolPtrEquals(a.State.MaintenanceDue, b.State.MaintenanceDue) &&
		floatPtrEquals(a.State.HwcCylinderTemperatureTop, b.State.HwcCylinderTemperatureTop) &&
		floatPtrEquals(a.State.HwcCylinderTemperatureBottom, b.State.HwcCylinderTemperatureBottom) &&
		boolPtrEquals(a.Config.AdaptiveHeatingCurve, b.Config.AdaptiveHeatingCurve) &&
		floatPtrEquals(a.Config.AlternativePoint, b.Config.AlternativePoint) &&
		floatPtrEquals(a.Config.HeatingCircuitBivalencePoint, b.Config.HeatingCircuitBivalencePoint) &&
		floatPtrEquals(a.Config.DhwBivalencePoint, b.Config.DhwBivalencePoint) &&
		configBoundEquals(a.Config.DhwBivalenceBound, b.Config.DhwBivalenceBound) &&
		floatPtrEquals(a.Config.HcEmergencyTemperature, b.Config.HcEmergencyTemperature) &&
		floatPtrEquals(a.Config.HwcMaxFlowTempDesired, b.Config.HwcMaxFlowTempDesired) &&
		intPtrEquals(a.Config.MaxRoomHumidity, b.Config.MaxRoomHumidity) &&
		stringPtrEquals(a.Config.MaintenanceDate, b.Config.MaintenanceDate) &&
		stringPtrEquals(a.Config.InstallerName, b.Config.InstallerName) &&
		stringPtrEquals(a.Config.InstallerPhone, b.Config.InstallerPhone) &&
		intPtrEquals(a.Config.InstallerMenuCode, b.Config.InstallerMenuCode) &&
		intPtrEquals(a.Properties.SystemScheme, b.Properties.SystemScheme) &&
		intPtrEquals(a.Properties.FunctionModuleConfigurationCode, b.Properties.FunctionModuleConfigurationCode)
}

// SystemStatusFromView maps the native values to the corresponding output record.
func SystemStatusFromView(view NativeSystemView) *SystemStatus {
	snapshot := view.Snapshot
	dhwBound := view.DhwBivalenceBound
	return &SystemStatus{
		State: SystemState{
			SystemOff:                      cloneBoilerBoolPtr(snapshot.SystemOff),
			SystemWaterPressure:            CloneFloat64Ptr(snapshot.SystemWaterPressure),
			SystemFlowTemperature:          CloneFloat64Ptr(snapshot.SystemFlowTemperature),
			OutdoorTemperature:             CloneFloat64Ptr(snapshot.OutdoorTemperature),
			OutdoorTemperatureDailyAverage: CloneFloat64Ptr(snapshot.OutdoorTemperatureDailyAverage),
			MaintenanceDue:                 cloneBoilerBoolPtr(snapshot.MaintenanceDue),
			HwcCylinderTemperatureTop:      CloneFloat64Ptr(snapshot.HwcCylinderTemperatureTop),
			HwcCylinderTemperatureBottom:   CloneFloat64Ptr(snapshot.HwcCylinderTemperatureBottom),
		},
		Config: SystemConfig{
			AdaptiveHeatingCurve:         cloneBoilerBoolPtr(snapshot.AdaptiveHeatingCurve),
			AlternativePoint:             CloneFloat64Ptr(snapshot.AlternativePoint),
			HeatingCircuitBivalencePoint: CloneFloat64Ptr(snapshot.HeatingCircuitBivalencePoint),
			DhwBivalencePoint:            CloneFloat64Ptr(snapshot.DhwBivalencePoint),
			DhwBivalenceBound:            &dhwBound,
			HcEmergencyTemperature:       CloneFloat64Ptr(snapshot.HcEmergencyTemperature),
			HwcMaxFlowTempDesired:        CloneFloat64Ptr(snapshot.HwcMaxFlowTempDesired),
			MaxRoomHumidity:              uint16ToIntPtr(snapshot.MaxRoomHumidity),
			MaintenanceDate:              CloneStringPtr(snapshot.MaintenanceDate),
			InstallerName:                CloneStringPtr(snapshot.InstallerName),
			InstallerPhone:               CloneStringPtr(snapshot.InstallerPhone),
			InstallerMenuCode:            uint16ToIntPtr(snapshot.InstallerMenuCode),
		},
		Properties: SystemProperties{
			SystemScheme:                    uint16ToIntPtr(snapshot.SystemScheme),
			FunctionModuleConfigurationCode: uint16ToIntPtr(snapshot.FunctionModuleConfigurationCode),
		},
	}
}

// RegulatorCapabilityFromNative maps the catalog regulator capability into the output vocabulary.
func RegulatorCapabilityFromNative(capability productids.ControllerCapability) VaillantRegulatorCapability {
	switch capability {
	case productids.ControllerPresent:
		return VaillantRegulatorCapabilityPresent
	case productids.ControllerNone:
		return VaillantRegulatorCapabilityNone
	default:
		return VaillantRegulatorCapabilityUnknown
	}
}
