package owner

import (
	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/output"
	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/owner/systemconfig"
)

func nativeFieldFreshness(fields map[semanticFieldKey]semanticFieldFreshness) map[output.SemanticFieldKey]output.SemanticFieldFreshness {
	if fields == nil {
		return nil
	}
	out := make(map[output.SemanticFieldKey]output.SemanticFieldFreshness, len(fields))
	for key, value := range fields {
		out[output.SemanticFieldKey(key)] = output.SemanticFieldFreshness{Source: output.SemanticSnapshotSource(value.Source), Stale: value.Stale}
	}
	return out
}

func ownerFieldFreshness(fields map[output.SemanticFieldKey]output.SemanticFieldFreshness) map[semanticFieldKey]semanticFieldFreshness {
	if fields == nil {
		return nil
	}
	out := make(map[semanticFieldKey]semanticFieldFreshness, len(fields))
	for key, value := range fields {
		out[semanticFieldKey(key)] = semanticFieldFreshness{Source: semanticSnapshotSource(value.Source), Stale: value.Stale}
	}
	return out
}

func nativeZoneSnapshot(entry *vaillantZoneSnapshot) *output.NativeZoneSnapshot {
	if entry == nil {
		return nil
	}
	return &output.NativeZoneSnapshot{
		Instance:                          entry.Instance,
		Present:                           entry.Present,
		Name:                              entry.Name,
		OperatingMode:                     entry.OperatingMode,
		Preset:                            entry.Preset,
		HvacAction:                        entry.HvacAction,
		AllowedModes:                      entry.AllowedModes,
		CurrentTempC:                      entry.CurrentTempC,
		TargetTempC:                       entry.TargetTempC,
		HumidityPct:                       entry.HumidityPct,
		ConfigurationHeatingOperationMode: entry.ConfigurationHeatingOperationMode,
		StateSpecialFunction:              entry.StateSpecialFunction,
		ConfigurationRoomTemperatureZoneMappingRaw: entry.ConfigurationRoomTemperatureZoneMappingRaw,
		ConfigurationAssociatedCircuitRaw:          entry.ConfigurationAssociatedCircuitRaw,
		ConfigurationCircuitTypeRaw:                entry.ConfigurationCircuitTypeRaw,
		StateValveStatusRaw:                        entry.StateValveStatusRaw,
		QuickVetoTempC:                             entry.QuickVetoTempC,
		QuickVetoDurationH:                         entry.QuickVetoDurationH,
		QuickVetoEndTime:                           entry.QuickVetoEndTime,
		QuickVetoEndDate:                           entry.QuickVetoEndDate,
		HolidayStartDate:                           entry.HolidayStartDate,
		HolidayEndDate:                             entry.HolidayEndDate,
		HolidaySetpointC:                           entry.HolidaySetpointC,
		HolidayStartTime:                           entry.HolidayStartTime,
		HolidayEndTime:                             entry.HolidayEndTime,
		FieldFreshness:                             nativeFieldFreshness(entry.FieldFreshness),
	}
}

func ownerZoneSnapshot(entry *output.NativeZoneSnapshot) *vaillantZoneSnapshot {
	if entry == nil {
		return nil
	}
	return &vaillantZoneSnapshot{
		Instance:                          entry.Instance,
		Present:                           entry.Present,
		Name:                              entry.Name,
		OperatingMode:                     entry.OperatingMode,
		Preset:                            entry.Preset,
		HvacAction:                        entry.HvacAction,
		AllowedModes:                      entry.AllowedModes,
		CurrentTempC:                      entry.CurrentTempC,
		TargetTempC:                       entry.TargetTempC,
		HumidityPct:                       entry.HumidityPct,
		ConfigurationHeatingOperationMode: entry.ConfigurationHeatingOperationMode,
		StateSpecialFunction:              entry.StateSpecialFunction,
		ConfigurationRoomTemperatureZoneMappingRaw: entry.ConfigurationRoomTemperatureZoneMappingRaw,
		ConfigurationAssociatedCircuitRaw:          entry.ConfigurationAssociatedCircuitRaw,
		ConfigurationCircuitTypeRaw:                entry.ConfigurationCircuitTypeRaw,
		StateValveStatusRaw:                        entry.StateValveStatusRaw,
		QuickVetoTempC:                             entry.QuickVetoTempC,
		QuickVetoDurationH:                         entry.QuickVetoDurationH,
		QuickVetoEndTime:                           entry.QuickVetoEndTime,
		QuickVetoEndDate:                           entry.QuickVetoEndDate,
		HolidayStartDate:                           entry.HolidayStartDate,
		HolidayEndDate:                             entry.HolidayEndDate,
		HolidaySetpointC:                           entry.HolidaySetpointC,
		HolidayStartTime:                           entry.HolidayStartTime,
		HolidayEndTime:                             entry.HolidayEndTime,
		FieldFreshness:                             ownerFieldFreshness(entry.FieldFreshness),
	}
}

func nativeDhwSnapshot(snapshot *vaillantDhwSnapshot) *output.NativeDhwSnapshot {
	if snapshot == nil {
		return nil
	}
	return &output.NativeDhwSnapshot{
		OperatingMode:                 snapshot.OperatingMode,
		Preset:                        snapshot.Preset,
		CurrentTempC:                  snapshot.CurrentTempC,
		TargetTempC:                   snapshot.TargetTempC,
		ConfigurationDHWOperationMode: snapshot.ConfigurationDHWOperationMode,
		StateSpecialFunction:          snapshot.StateSpecialFunction,
		HolidayStartDate:              snapshot.HolidayStartDate,
		HolidayEndDate:                snapshot.HolidayEndDate,
		FieldFreshness:                nativeFieldFreshness(snapshot.FieldFreshness),
	}
}

func ownerDhwSnapshot(snapshot *output.NativeDhwSnapshot) *vaillantDhwSnapshot {
	if snapshot == nil {
		return nil
	}
	return &vaillantDhwSnapshot{
		OperatingMode:                 snapshot.OperatingMode,
		Preset:                        snapshot.Preset,
		CurrentTempC:                  snapshot.CurrentTempC,
		TargetTempC:                   snapshot.TargetTempC,
		ConfigurationDHWOperationMode: snapshot.ConfigurationDHWOperationMode,
		StateSpecialFunction:          snapshot.StateSpecialFunction,
		HolidayStartDate:              snapshot.HolidayStartDate,
		HolidayEndDate:                snapshot.HolidayEndDate,
		FieldFreshness:                ownerFieldFreshness(snapshot.FieldFreshness),
	}
}

func nativeSystemSnapshot(snapshot *vaillantSystemSnapshot) *output.NativeSystemSnapshot {
	if snapshot == nil {
		return nil
	}
	return &output.NativeSystemSnapshot{
		Controller:                      snapshot.Controller,
		SystemOff:                       snapshot.SystemOff,
		SystemWaterPressure:             snapshot.SystemWaterPressure,
		SystemFlowTemperature:           snapshot.SystemFlowTemperature,
		SystemFlowTemperatureLiveAt:     snapshot.SystemFlowTemperatureLiveAt,
		OutdoorTemperature:              snapshot.OutdoorTemperature,
		OutdoorTemperatureDailyAverage:  snapshot.OutdoorTemperatureAverage24Hours,
		MaintenanceDue:                  snapshot.MaintenanceDue,
		HwcCylinderTemperatureTop:       snapshot.HwcCylinderTemperatureTop,
		HwcCylinderTemperatureBottom:    snapshot.HwcCylinderTemperatureBottom,
		AdaptiveHeatingCurve:            snapshot.AdaptiveHeatingCurve,
		AlternativePoint:                snapshot.AlternativePoint,
		HeatingCircuitBivalencePoint:    snapshot.HeatingCircuitBivalencePoint,
		DhwBivalencePoint:               snapshot.DhwBivalencePoint,
		HcEmergencyTemperature:          snapshot.HcEmergencyTemperature,
		HwcMaxFlowTempDesired:           snapshot.HwcMaxFlowTempDesired,
		MaxRoomHumidity:                 snapshot.MaxRoomHumidity,
		MaintenanceDate:                 snapshot.MaintenanceDate,
		InstallerName:                   snapshot.InstallerName,
		InstallerPhone:                  snapshot.InstallerPhone,
		InstallerMenuCode:               snapshot.InstallerMenuCode,
		SystemScheme:                    snapshot.SystemScheme,
		FunctionModuleConfigurationCode: snapshot.ModuleConfigurationFunctionModuleController,
	}
}

func outputConfigBound(bound systemconfig.ConfigBound) output.ConfigBound {
	return output.ConfigBound{
		ControllerAddress: bound.ControllerAddress,
		ControllerRegistryIdentity: output.ControllerRegistryIdentity{
			Valid:           bound.ControllerRegistryIdentity.Valid,
			Address:         bound.ControllerRegistryIdentity.Address,
			Generation:      bound.ControllerRegistryIdentity.Generation,
			Manufacturer:    bound.ControllerRegistryIdentity.Manufacturer,
			DeviceID:        bound.ControllerRegistryIdentity.DeviceID,
			SoftwareVersion: bound.ControllerRegistryIdentity.SoftwareVersion,
			HardwareVersion: bound.ControllerRegistryIdentity.HardwareVersion,
			SerialNumber:    bound.ControllerRegistryIdentity.SerialNumber,
			MacAddress:      bound.ControllerRegistryIdentity.MacAddress,
		},
		MinC:         bound.MinC,
		MaxC:         bound.MaxC,
		Source:       bound.Source,
		Reason:       bound.Reason,
		Profile:      bound.Profile,
		NarrowedOnly: bound.NarrowedOnly,
	}
}

func nativeFunctionModuleVerdict(verdict vaillantFunctionModuleInterpretation) output.NativeFunctionModuleInterpretation {
	return output.NativeFunctionModuleInterpretation{
		Mode:             output.NativeFunctionModuleMode(verdict.Mode),
		DegradedReason:   output.NativeFunctionModuleDegradedReason(verdict.DegradedReason),
		EvidenceRevision: verdict.EvidenceRevision,
	}
}

func nativeScheduleDay(day vaillantScheduleDay) output.NativeScheduleDay {
	out := output.NativeScheduleDay{Weekday: day.Weekday}
	if day.Slots != nil {
		out.Slots = make([]output.NativeScheduleTimerSlot, len(day.Slots))
		for i, slot := range day.Slots {
			out.Slots[i] = output.NativeScheduleTimerSlot(slot)
		}
	}
	return out
}

func nativeSchedulePrograms(programs []vaillantScheduleProgram) []output.NativeScheduleProgram {
	if programs == nil {
		return nil
	}
	out := make([]output.NativeScheduleProgram, len(programs))
	for i, program := range programs {
		out[i] = output.NativeScheduleProgram{Zone: program.Zone, HC: program.HC, SlotsUsed: program.SlotsUsed}
		if program.Config != nil {
			config := output.NativeScheduleConfig(*program.Config)
			out[i].Config = &config
		}
		if program.Days != nil {
			out[i].Days = make([]output.NativeScheduleDay, len(program.Days))
			for d, day := range program.Days {
				out[i].Days[d] = nativeScheduleDay(day)
			}
		}
	}
	return out
}
