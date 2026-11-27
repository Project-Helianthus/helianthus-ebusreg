package owner

import (
	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/output"
)

func cloneEnergyChannel(channel output.EnergyChannel) output.EnergyChannel {
	channel.DHW = cloneEnergySeries(channel.DHW)
	channel.Climate = cloneEnergySeries(channel.Climate)
	return channel
}

func cloneEnergySeries(series output.EnergySeries) output.EnergySeries {
	if len(series.Yearly) > 0 {
		copySlice := make([]float64, len(series.Yearly))
		copy(copySlice, series.Yearly)
		series.Yearly = copySlice
	}
	if len(series.Monthly) > 0 {
		copySlice := make([]float64, len(series.Monthly))
		copy(copySlice, series.Monthly)
		series.Monthly = copySlice
	}
	if len(series.YearlyMeta) > 0 {
		copySlice := make([]output.EnergyPointMeta, len(series.YearlyMeta))
		copy(copySlice, series.YearlyMeta)
		series.YearlyMeta = copySlice
	}
	if len(series.MonthlyMeta) > 0 {
		copySlice := make([]output.EnergyPointMeta, len(series.MonthlyMeta))
		copy(copySlice, series.MonthlyMeta)
		series.MonthlyMeta = copySlice
	}
	return series
}

func cloneZone(z output.Zone) output.Zone {
	out := z
	out.State = cloneZoneState(z.State)
	out.Config = cloneZoneConfig(z.Config)
	return out
}

func cloneZoneState(s output.ZoneState) output.ZoneState {
	out := s
	if s.CurrentTempC != nil {
		v := *s.CurrentTempC
		out.CurrentTempC = &v
	}
	if s.CurrentHumidityPct != nil {
		v := *s.CurrentHumidityPct
		out.CurrentHumidityPct = &v
	}
	if s.HeatingDemandPct != nil {
		v := *s.HeatingDemandPct
		out.HeatingDemandPct = &v
	}
	if s.ValvePositionPct != nil {
		v := *s.ValvePositionPct
		out.ValvePositionPct = &v
	}
	return out
}

func cloneZoneConfig(c output.ZoneConfig) output.ZoneConfig {
	out := c
	if c.TargetTempC != nil {
		v := *c.TargetTempC
		out.TargetTempC = &v
	}
	if c.OperationModeChangeable != nil {
		v := *c.OperationModeChangeable
		out.OperationModeChangeable = &v
	}
	if c.SourceLabel != nil {
		v := *c.SourceLabel
		out.SourceLabel = &v
	}
	if len(c.AllowedModes) > 0 {
		out.AllowedModes = make([]string, len(c.AllowedModes))
		copy(out.AllowedModes, c.AllowedModes)
	}
	if c.AssociatedCircuit != nil {
		v := *c.AssociatedCircuit
		out.AssociatedCircuit = &v
	}
	if c.RoomTemperatureZoneMapping != nil {
		v := *c.RoomTemperatureZoneMapping
		out.RoomTemperatureZoneMapping = &v
	}
	return out
}

func cloneDhwStatus(s *output.DhwStatus) *output.DhwStatus {
	if s == nil {
		return nil
	}
	out := *s
	out.State = cloneDhwState(s.State)
	out.Config = cloneDhwConfig(s.Config)
	return &out
}

func cloneDhwState(s output.DhwState) output.DhwState {
	out := s
	if s.CurrentTempC != nil {
		v := *s.CurrentTempC
		out.CurrentTempC = &v
	}
	if s.OverrunActive != nil {
		v := *s.OverrunActive
		out.OverrunActive = &v
	}
	if s.HeatingDemandPct != nil {
		v := *s.HeatingDemandPct
		out.HeatingDemandPct = &v
	}
	return out
}

func cloneDhwConfig(c output.DhwConfig) output.DhwConfig {
	out := c
	if c.TargetTempC != nil {
		v := *c.TargetTempC
		out.TargetTempC = &v
	}
	if c.OperationModeChangeable != nil {
		v := *c.OperationModeChangeable
		out.OperationModeChangeable = &v
	}
	return out
}

func cloneCircuitStatus(status output.CircuitStatus) output.CircuitStatus {
	out := status
	out.State = cloneCircuitState(status.State)
	out.Config = cloneCircuitConfig(status.Config)
	out.ManagingDevice = cloneManagingDevice(status.ManagingDevice)
	return out
}

func cloneManagingDevice(device output.ManagingDevice) output.ManagingDevice {
	out := device
	if device.Role == "" {
		out.Role = output.ManagingDeviceRoleUnknown
	}
	if device.DeviceID != nil {
		v := *device.DeviceID
		out.DeviceID = &v
	}
	if device.Address != nil {
		v := *device.Address
		out.Address = &v
	}
	return out
}

func cloneCircuitState(state output.CircuitState) output.CircuitState {
	out := state
	if state.PumpActive != nil {
		v := *state.PumpActive
		out.PumpActive = &v
	}
	if state.MixerPositionPct != nil {
		v := *state.MixerPositionPct
		out.MixerPositionPct = &v
	}
	if state.FlowTemperatureC != nil {
		v := *state.FlowTemperatureC
		out.FlowTemperatureC = &v
	}
	if state.FlowSetpointC != nil {
		v := *state.FlowSetpointC
		out.FlowSetpointC = &v
	}
	if state.CalcFlowTempC != nil {
		v := *state.CalcFlowTempC
		out.CalcFlowTempC = &v
	}
	if state.Humidity != nil {
		v := *state.Humidity
		out.Humidity = &v
	}
	if state.DewPoint != nil {
		v := *state.DewPoint
		out.DewPoint = &v
	}
	if state.PumpHours != nil {
		v := *state.PumpHours
		out.PumpHours = &v
	}
	if state.PumpStarts != nil {
		v := *state.PumpStarts
		out.PumpStarts = &v
	}
	return out
}

func cloneCircuitConfig(config output.CircuitConfig) output.CircuitConfig {
	out := config
	if config.HeatingCurve != nil {
		v := *config.HeatingCurve
		out.HeatingCurve = &v
	}
	if config.FlowTempMaxC != nil {
		v := *config.FlowTempMaxC
		out.FlowTempMaxC = &v
	}
	if config.FlowTempMinC != nil {
		v := *config.FlowTempMinC
		out.FlowTempMinC = &v
	}
	if config.SummerLimitC != nil {
		v := *config.SummerLimitC
		out.SummerLimitC = &v
	}
	if config.SpecialFunctionMode != nil {
		v := *config.SpecialFunctionMode
		out.SpecialFunctionMode = &v
	}
	if config.CoolingEnabled != nil {
		v := *config.CoolingEnabled
		out.CoolingEnabled = &v
	}
	return out
}

func cloneRadioDevice(device output.RadioDevice) output.RadioDevice {
	out := device
	if device.DeviceConnected != nil {
		v := *device.DeviceConnected
		out.DeviceConnected = &v
	}
	if device.DeviceClassAddress != nil {
		v := *device.DeviceClassAddress
		out.DeviceClassAddress = &v
	}
	if device.FirmwareVersion != nil {
		v := *device.FirmwareVersion
		out.FirmwareVersion = &v
	}
	if device.HardwareIdentifier != nil {
		v := *device.HardwareIdentifier
		out.HardwareIdentifier = &v
	}
	if device.RemoteControlAddress != nil {
		v := *device.RemoteControlAddress
		out.RemoteControlAddress = &v
	}
	if device.DevicePaired != nil {
		v := *device.DevicePaired
		out.DevicePaired = &v
	}
	if device.ReceptionStrength != nil {
		v := *device.ReceptionStrength
		out.ReceptionStrength = &v
	}
	if device.ZoneAssignment != nil {
		v := *device.ZoneAssignment
		out.ZoneAssignment = &v
	}
	if device.RoomTemperatureC != nil {
		v := *device.RoomTemperatureC
		out.RoomTemperatureC = &v
	}
	if device.RoomHumidityPct != nil {
		v := *device.RoomHumidityPct
		out.RoomHumidityPct = &v
	}
	return out
}

func cloneSolarStatus(status *output.SolarStatus) *output.SolarStatus {
	if status == nil {
		return nil
	}
	out := *status
	if status.CollectorTemperatureC != nil {
		v := *status.CollectorTemperatureC
		out.CollectorTemperatureC = &v
	}
	if status.ReturnTemperatureC != nil {
		v := *status.ReturnTemperatureC
		out.ReturnTemperatureC = &v
	}
	if status.PumpActive != nil {
		v := *status.PumpActive
		out.PumpActive = &v
	}
	if status.CurrentYield != nil {
		v := *status.CurrentYield
		out.CurrentYield = &v
	}
	if status.PumpHours != nil {
		v := *status.PumpHours
		out.PumpHours = &v
	}
	if status.SolarEnabled != nil {
		v := *status.SolarEnabled
		out.SolarEnabled = &v
	}
	if status.FunctionMode != nil {
		v := *status.FunctionMode
		out.FunctionMode = &v
	}
	return &out
}

func cloneCylinderStatus(status output.CylinderStatus) output.CylinderStatus {
	out := status
	if status.TemperatureC != nil {
		v := *status.TemperatureC
		out.TemperatureC = &v
	}
	if status.MaxSetpointC != nil {
		v := *status.MaxSetpointC
		out.MaxSetpointC = &v
	}
	if status.ChargeHysteresisC != nil {
		v := *status.ChargeHysteresisC
		out.ChargeHysteresisC = &v
	}
	if status.ChargeOffsetC != nil {
		v := *status.ChargeOffsetC
		out.ChargeOffsetC = &v
	}
	return out
}

func cloneCylinderStatuses(cylinders []output.CylinderStatus) []output.CylinderStatus {
	if cylinders == nil {
		return nil
	}
	out := make([]output.CylinderStatus, len(cylinders))
	for i := range cylinders {
		out[i] = cloneCylinderStatus(cylinders[i])
	}
	return out
}
