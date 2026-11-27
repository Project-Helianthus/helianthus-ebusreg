package owner

import (
	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/output"
)

func cloneBoilerState(state output.BoilerState) output.BoilerState {
	if state.FlowTemperatureC != nil {
		v := *state.FlowTemperatureC
		state.FlowTemperatureC = &v
	}
	if state.ReturnTemperatureC != nil {
		v := *state.ReturnTemperatureC
		state.ReturnTemperatureC = &v
	}
	if state.CentralHeatingPumpActive != nil {
		v := *state.CentralHeatingPumpActive
		state.CentralHeatingPumpActive = &v
	}
	if state.WaterPressureBar != nil {
		v := *state.WaterPressureBar
		state.WaterPressureBar = &v
	}
	if state.ExternalPumpActive != nil {
		v := *state.ExternalPumpActive
		state.ExternalPumpActive = &v
	}
	if state.CirculationPumpActive != nil {
		v := *state.CirculationPumpActive
		state.CirculationPumpActive = &v
	}
	if state.GasValveActive != nil {
		v := *state.GasValveActive
		state.GasValveActive = &v
	}
	if state.FlameActive != nil {
		v := *state.FlameActive
		state.FlameActive = &v
	}
	if state.DiverterValvePositionPct != nil {
		v := *state.DiverterValvePositionPct
		state.DiverterValvePositionPct = &v
	}
	if state.FanSpeedRpm != nil {
		v := *state.FanSpeedRpm
		state.FanSpeedRpm = &v
	}
	if state.TargetFanSpeedRpm != nil {
		v := *state.TargetFanSpeedRpm
		state.TargetFanSpeedRpm = &v
	}
	if state.IonisationVoltageUa != nil {
		v := *state.IonisationVoltageUa
		state.IonisationVoltageUa = &v
	}
	if state.DhwWaterFlowLpm != nil {
		v := *state.DhwWaterFlowLpm
		state.DhwWaterFlowLpm = &v
	}
	if state.DhwDemandActive != nil {
		v := *state.DhwDemandActive
		state.DhwDemandActive = &v
	}
	if state.HeatingSwitchActive != nil {
		v := *state.HeatingSwitchActive
		state.HeatingSwitchActive = &v
	}
	if state.StorageLoadPumpPct != nil {
		v := *state.StorageLoadPumpPct
		state.StorageLoadPumpPct = &v
	}
	if state.ModulationPct != nil {
		v := *state.ModulationPct
		state.ModulationPct = &v
	}
	if state.PrimaryCircuitFlowLpm != nil {
		v := *state.PrimaryCircuitFlowLpm
		state.PrimaryCircuitFlowLpm = &v
	}
	if state.FlowTempDesiredC != nil {
		v := *state.FlowTempDesiredC
		state.FlowTempDesiredC = &v
	}
	if state.DhwTempDesiredC != nil {
		v := *state.DhwTempDesiredC
		state.DhwTempDesiredC = &v
	}
	if state.StateNumber != nil {
		v := *state.StateNumber
		state.StateNumber = &v
	}
	if state.DhwTemperatureC != nil {
		v := *state.DhwTemperatureC
		state.DhwTemperatureC = &v
	}
	if state.DhwTargetTemperatureC != nil {
		v := *state.DhwTargetTemperatureC
		state.DhwTargetTemperatureC = &v
	}
	return state
}

func cloneBoilerConfig(config output.BoilerConfig) output.BoilerConfig {
	if config.DhwOperatingMode != nil {
		v := *config.DhwOperatingMode
		config.DhwOperatingMode = &v
	}
	if config.FlowsetHcMaxC != nil {
		v := *config.FlowsetHcMaxC
		config.FlowsetHcMaxC = &v
	}
	if config.FlowsetHwcMaxC != nil {
		v := *config.FlowsetHwcMaxC
		config.FlowsetHwcMaxC = &v
	}
	if config.PartloadHcKW != nil {
		v := *config.PartloadHcKW
		config.PartloadHcKW = &v
	}
	if config.PartloadHwcKW != nil {
		v := *config.PartloadHwcKW
		config.PartloadHwcKW = &v
	}
	if config.InstallerMenuCode != nil {
		v := *config.InstallerMenuCode
		config.InstallerMenuCode = &v
	}
	if config.PhoneNumber != nil {
		v := *config.PhoneNumber
		config.PhoneNumber = &v
	}
	if config.HoursTillService != nil {
		v := *config.HoursTillService
		config.HoursTillService = &v
	}
	return config
}

func cloneBoilerDiagnostics(diag output.BoilerDiagnostics) output.BoilerDiagnostics {
	if diag.HeatingStatusRaw != nil {
		v := *diag.HeatingStatusRaw
		diag.HeatingStatusRaw = &v
	}
	if diag.DhwStatusRaw != nil {
		v := *diag.DhwStatusRaw
		diag.DhwStatusRaw = &v
	}
	if diag.CentralHeatingHours != nil {
		v := *diag.CentralHeatingHours
		diag.CentralHeatingHours = &v
	}
	if diag.DhwHours != nil {
		v := *diag.DhwHours
		diag.DhwHours = &v
	}
	if diag.CentralHeatingStarts != nil {
		v := *diag.CentralHeatingStarts
		diag.CentralHeatingStarts = &v
	}
	if diag.DhwStarts != nil {
		v := *diag.DhwStarts
		diag.DhwStarts = &v
	}
	if diag.PumpHours != nil {
		v := *diag.PumpHours
		diag.PumpHours = &v
	}
	if diag.FanHours != nil {
		v := *diag.FanHours
		diag.FanHours = &v
	}
	if diag.DeactivationsIFC != nil {
		v := *diag.DeactivationsIFC
		diag.DeactivationsIFC = &v
	}
	if diag.DeactivationsTemplimiter != nil {
		v := *diag.DeactivationsTemplimiter
		diag.DeactivationsTemplimiter = &v
	}
	return diag
}
