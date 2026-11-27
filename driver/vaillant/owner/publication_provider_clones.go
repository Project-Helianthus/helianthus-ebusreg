package owner

import "github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/output"

func cloneAdapterHardwareInfo(info *output.AdapterHardwareInfo) *output.AdapterHardwareInfo {
	if info == nil {
		return nil
	}
	cp := *info
	if info.JumperFlags != nil {
		cp.JumperFlags = make([]string, len(info.JumperFlags))
		copy(cp.JumperFlags, info.JumperFlags)
	}
	if info.TemperatureC != nil {
		v := *info.TemperatureC
		cp.TemperatureC = &v
	}
	if info.SupplyVoltageMV != nil {
		v := *info.SupplyVoltageMV
		cp.SupplyVoltageMV = &v
	}
	if info.BusVoltageMaxDV != nil {
		v := *info.BusVoltageMaxDV
		cp.BusVoltageMaxDV = &v
	}
	if info.BusVoltageMinDV != nil {
		v := *info.BusVoltageMinDV
		cp.BusVoltageMinDV = &v
	}
	if info.ResetCause != nil {
		v := *info.ResetCause
		cp.ResetCause = &v
	}
	if info.ResetCauseCode != nil {
		v := *info.ResetCauseCode
		cp.ResetCauseCode = &v
	}
	if info.RestartCount != nil {
		v := *info.RestartCount
		cp.RestartCount = &v
	}
	if info.WiFiRSSIDBm != nil {
		v := *info.WiFiRSSIDBm
		cp.WiFiRSSIDBm = &v
	}
	if info.LastIdentityQuery != nil {
		v := *info.LastIdentityQuery
		cp.LastIdentityQuery = &v
	}
	if info.LastTelemetryQuery != nil {
		v := *info.LastTelemetryQuery
		cp.LastTelemetryQuery = &v
	}
	return &cp
}

func cloneSystemStatus(status *output.SystemStatus) *output.SystemStatus {
	if status == nil {
		return nil
	}
	cp := *status
	cp.State = cloneSystemState(cp.State)
	cp.Config = cloneSystemConfig(cp.Config)
	cp.Properties = cloneSystemProperties(cp.Properties)
	if status.GatewayBrand != nil {
		v := *status.GatewayBrand
		cp.GatewayBrand = &v
	}
	if status.GatewayVendor != nil {
		v := *status.GatewayVendor
		cp.GatewayVendor = &v
	}
	return &cp
}

func cloneSystemState(state output.SystemState) output.SystemState {
	if state.SystemOff != nil {
		v := *state.SystemOff
		state.SystemOff = &v
	}
	if state.SystemWaterPressure != nil {
		v := *state.SystemWaterPressure
		state.SystemWaterPressure = &v
	}
	if state.SystemFlowTemperature != nil {
		v := *state.SystemFlowTemperature
		state.SystemFlowTemperature = &v
	}
	if state.OutdoorTemperature != nil {
		v := *state.OutdoorTemperature
		state.OutdoorTemperature = &v
	}
	if state.OutdoorTemperatureDailyAverage != nil {
		v := *state.OutdoorTemperatureDailyAverage
		state.OutdoorTemperatureDailyAverage = &v
	}
	if state.MaintenanceDue != nil {
		v := *state.MaintenanceDue
		state.MaintenanceDue = &v
	}
	if state.HwcCylinderTemperatureTop != nil {
		v := *state.HwcCylinderTemperatureTop
		state.HwcCylinderTemperatureTop = &v
	}
	if state.HwcCylinderTemperatureBottom != nil {
		v := *state.HwcCylinderTemperatureBottom
		state.HwcCylinderTemperatureBottom = &v
	}
	return state
}

func cloneSystemConfig(config output.SystemConfig) output.SystemConfig {
	if config.AdaptiveHeatingCurve != nil {
		v := *config.AdaptiveHeatingCurve
		config.AdaptiveHeatingCurve = &v
	}
	if config.AlternativePoint != nil {
		v := *config.AlternativePoint
		config.AlternativePoint = &v
	}
	if config.HeatingCircuitBivalencePoint != nil {
		v := *config.HeatingCircuitBivalencePoint
		config.HeatingCircuitBivalencePoint = &v
	}
	if config.DhwBivalencePoint != nil {
		v := *config.DhwBivalencePoint
		config.DhwBivalencePoint = &v
	}
	if config.DhwBivalenceBound != nil {
		v := *config.DhwBivalenceBound
		config.DhwBivalenceBound = &v
	}
	if config.HcEmergencyTemperature != nil {
		v := *config.HcEmergencyTemperature
		config.HcEmergencyTemperature = &v
	}
	if config.HwcMaxFlowTempDesired != nil {
		v := *config.HwcMaxFlowTempDesired
		config.HwcMaxFlowTempDesired = &v
	}
	if config.MaxRoomHumidity != nil {
		v := *config.MaxRoomHumidity
		config.MaxRoomHumidity = &v
	}
	if config.MaintenanceDate != nil {
		v := *config.MaintenanceDate
		config.MaintenanceDate = &v
	}
	if config.InstallerName != nil {
		v := *config.InstallerName
		config.InstallerName = &v
	}
	if config.InstallerPhone != nil {
		v := *config.InstallerPhone
		config.InstallerPhone = &v
	}
	if config.InstallerMenuCode != nil {
		v := *config.InstallerMenuCode
		config.InstallerMenuCode = &v
	}
	return config
}

func cloneScheduleStatus(status *output.ScheduleStatus) *output.ScheduleStatus {
	if status == nil {
		return nil
	}
	cp := output.ScheduleStatus{}
	if len(status.Programs) > 0 {
		cp.Programs = make([]output.ScheduleProgram, len(status.Programs))
		for i, prog := range status.Programs {
			cp.Programs[i] = cloneScheduleProgram(prog)
		}
	}
	return &cp
}

func cloneScheduleProgram(prog output.ScheduleProgram) output.ScheduleProgram {
	out := prog
	if prog.Config != nil {
		cfgCopy := *prog.Config
		if prog.Config.MinTempC != nil {
			v := *prog.Config.MinTempC
			cfgCopy.MinTempC = &v
		}
		if prog.Config.MaxTempC != nil {
			v := *prog.Config.MaxTempC
			cfgCopy.MaxTempC = &v
		}
		out.Config = &cfgCopy
	}
	if len(prog.SlotsUsed) > 0 {
		out.SlotsUsed = make([]int, len(prog.SlotsUsed))
		copy(out.SlotsUsed, prog.SlotsUsed)
	}
	if len(prog.Days) > 0 {
		out.Days = make([]output.ScheduleDayProgram, len(prog.Days))
		for i, day := range prog.Days {
			out.Days[i] = cloneScheduleDayProgram(day)
		}
	}
	return out
}

func cloneScheduleDayProgram(day output.ScheduleDayProgram) output.ScheduleDayProgram {
	out := day
	if len(day.Slots) > 0 {
		out.Slots = make([]output.ScheduleTimerSlot, len(day.Slots))
		for i, slot := range day.Slots {
			out.Slots[i] = cloneScheduleTimerSlot(slot)
		}
	}
	return out
}

func cloneScheduleTimerSlot(slot output.ScheduleTimerSlot) output.ScheduleTimerSlot {
	out := slot
	if slot.TemperatureC != nil {
		v := *slot.TemperatureC
		out.TemperatureC = &v
	}
	if slot.TemperatureRaw != nil {
		v := *slot.TemperatureRaw
		out.TemperatureRaw = &v
	}
	return out
}

func cloneSystemProperties(properties output.SystemProperties) output.SystemProperties {
	if properties.SystemScheme != nil {
		v := *properties.SystemScheme
		properties.SystemScheme = &v
	}
	if properties.FunctionModuleConfigurationCode != nil {
		v := *properties.FunctionModuleConfigurationCode
		properties.FunctionModuleConfigurationCode = &v
	}
	return properties
}
