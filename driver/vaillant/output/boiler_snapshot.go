package output

import (
	"strings"
)

// MergeBoilerSnapshotNonDestructive merges the present fields without discarding retained values.
func MergeBoilerSnapshotNonDestructive(existing, incoming *NativeBoilerSnapshot) *NativeBoilerSnapshot {
	merged := CloneBoilerSnapshot(existing)
	if merged == nil {
		merged = &NativeBoilerSnapshot{}
	}
	if incoming != nil {
		if incoming.FlowTemperatureC != nil {
			merged.FlowTemperatureC = CloneFloat64Ptr(incoming.FlowTemperatureC)
		}
		if incoming.CentralHeatingPumpActive != nil {
			merged.CentralHeatingPumpActive = cloneBoilerBoolPtr(incoming.CentralHeatingPumpActive)
		}
		if incoming.WaterPressureBar != nil {
			merged.WaterPressureBar = CloneFloat64Ptr(incoming.WaterPressureBar)
		}
		if incoming.ExternalPumpActive != nil {
			merged.ExternalPumpActive = cloneBoilerBoolPtr(incoming.ExternalPumpActive)
		}
		if incoming.CirculationPumpActive != nil {
			merged.CirculationPumpActive = cloneBoilerBoolPtr(incoming.CirculationPumpActive)
		}
		if incoming.GasValveActive != nil {
			merged.GasValveActive = cloneBoilerBoolPtr(incoming.GasValveActive)
		}
		if incoming.FlameActive != nil {
			merged.FlameActive = cloneBoilerBoolPtr(incoming.FlameActive)
		}
		if incoming.DiverterValvePositionPct != nil {
			merged.DiverterValvePositionPct = CloneFloat64Ptr(incoming.DiverterValvePositionPct)
		}
		if incoming.FanSpeedRpm != nil {
			merged.FanSpeedRpm = cloneBoilerIntPtr(incoming.FanSpeedRpm)
		}
		if incoming.TargetFanSpeedRpm != nil {
			merged.TargetFanSpeedRpm = cloneBoilerIntPtr(incoming.TargetFanSpeedRpm)
		}
		if incoming.IonisationVoltageUa != nil {
			merged.IonisationVoltageUa = CloneFloat64Ptr(incoming.IonisationVoltageUa)
		}
		if incoming.DhwWaterFlowLpm != nil {
			merged.DhwWaterFlowLpm = CloneFloat64Ptr(incoming.DhwWaterFlowLpm)
		}
		if incoming.DhwDemandActive != nil {
			merged.DhwDemandActive = cloneBoilerBoolPtr(incoming.DhwDemandActive)
		}
		if incoming.HeatingSwitchActive != nil {
			merged.HeatingSwitchActive = cloneBoilerBoolPtr(incoming.HeatingSwitchActive)
		}
		if incoming.StorageLoadPumpPct != nil {
			merged.StorageLoadPumpPct = CloneFloat64Ptr(incoming.StorageLoadPumpPct)
		}
		if incoming.ModulationPct != nil {
			merged.ModulationPct = CloneFloat64Ptr(incoming.ModulationPct)
		}
		if incoming.PrimaryCircuitFlowLpm != nil {
			merged.PrimaryCircuitFlowLpm = CloneFloat64Ptr(incoming.PrimaryCircuitFlowLpm)
		}
		if incoming.FlowTempDesiredC != nil {
			merged.FlowTempDesiredC = CloneFloat64Ptr(incoming.FlowTempDesiredC)
		}
		if incoming.DhwTempDesiredC != nil {
			merged.DhwTempDesiredC = CloneFloat64Ptr(incoming.DhwTempDesiredC)
		}
		if incoming.StateNumber != nil {
			merged.StateNumber = cloneBoilerIntPtr(incoming.StateNumber)
		}
		if incoming.DhwTemperatureC != nil {
			merged.DhwTemperatureC = CloneFloat64Ptr(incoming.DhwTemperatureC)
		}
		if incoming.DhwTargetTemperatureC != nil {
			merged.DhwTargetTemperatureC = CloneFloat64Ptr(incoming.DhwTargetTemperatureC)
		}
		if incoming.DhwOperatingMode != nil {
			merged.DhwOperatingMode = cloneBoilerIntPtr(incoming.DhwOperatingMode)
		}
		if incoming.FlowsetHcMaxC != nil {
			merged.FlowsetHcMaxC = CloneFloat64Ptr(incoming.FlowsetHcMaxC)
		}
		if incoming.FlowsetHwcMaxC != nil {
			merged.FlowsetHwcMaxC = CloneFloat64Ptr(incoming.FlowsetHwcMaxC)
		}
		if incoming.PartloadHcKW != nil {
			merged.PartloadHcKW = CloneFloat64Ptr(incoming.PartloadHcKW)
		}
		if incoming.PartloadHwcKW != nil {
			merged.PartloadHwcKW = CloneFloat64Ptr(incoming.PartloadHwcKW)
		}
		if incoming.HeatingStatusRaw != nil {
			merged.HeatingStatusRaw = cloneBoilerIntPtr(incoming.HeatingStatusRaw)
		}
		if incoming.DhwStatusRaw != nil {
			merged.DhwStatusRaw = cloneBoilerIntPtr(incoming.DhwStatusRaw)
		}
		if incoming.CentralHeatingHours != nil {
			merged.CentralHeatingHours = CloneFloat64Ptr(incoming.CentralHeatingHours)
		}
		if incoming.DhwHours != nil {
			merged.DhwHours = CloneFloat64Ptr(incoming.DhwHours)
		}
		if incoming.CentralHeatingStarts != nil {
			merged.CentralHeatingStarts = cloneBoilerIntPtr(incoming.CentralHeatingStarts)
		}
		if incoming.DhwStarts != nil {
			merged.DhwStarts = cloneBoilerIntPtr(incoming.DhwStarts)
		}
		if incoming.PumpHours != nil {
			merged.PumpHours = CloneFloat64Ptr(incoming.PumpHours)
		}
		if incoming.FanHours != nil {
			merged.FanHours = CloneFloat64Ptr(incoming.FanHours)
		}
		if incoming.DeactivationsIFC != nil {
			merged.DeactivationsIFC = cloneBoilerIntPtr(incoming.DeactivationsIFC)
		}
		if incoming.DeactivationsTemplimiter != nil {
			merged.DeactivationsTemplimiter = cloneBoilerIntPtr(incoming.DeactivationsTemplimiter)
		}
		if incoming.InstallerMenuCode != nil {
			merged.InstallerMenuCode = cloneBoilerIntPtr(incoming.InstallerMenuCode)
		}
		if incoming.PhoneNumber != nil {
			merged.PhoneNumber = CloneStringPtr(incoming.PhoneNumber)
		}
		if incoming.HoursTillService != nil {
			merged.HoursTillService = cloneBoilerIntPtr(incoming.HoursTillService)
		}
	}

	merged.ReturnTemperatureC = nil
	return merged
}

// CloneBoilerSnapshot detaches the boiler snapshot from its input.
func CloneBoilerSnapshot(snapshot *NativeBoilerSnapshot) *NativeBoilerSnapshot {
	if snapshot == nil {
		return nil
	}
	return &NativeBoilerSnapshot{
		FlowTemperatureC:         CloneFloat64Ptr(snapshot.FlowTemperatureC),
		ReturnTemperatureC:       CloneFloat64Ptr(snapshot.ReturnTemperatureC),
		CentralHeatingPumpActive: cloneBoilerBoolPtr(snapshot.CentralHeatingPumpActive),
		WaterPressureBar:         CloneFloat64Ptr(snapshot.WaterPressureBar),
		ExternalPumpActive:       cloneBoilerBoolPtr(snapshot.ExternalPumpActive),
		CirculationPumpActive:    cloneBoilerBoolPtr(snapshot.CirculationPumpActive),
		GasValveActive:           cloneBoilerBoolPtr(snapshot.GasValveActive),
		FlameActive:              cloneBoilerBoolPtr(snapshot.FlameActive),
		DiverterValvePositionPct: CloneFloat64Ptr(snapshot.DiverterValvePositionPct),
		FanSpeedRpm:              cloneBoilerIntPtr(snapshot.FanSpeedRpm),
		TargetFanSpeedRpm:        cloneBoilerIntPtr(snapshot.TargetFanSpeedRpm),
		IonisationVoltageUa:      CloneFloat64Ptr(snapshot.IonisationVoltageUa),
		DhwWaterFlowLpm:          CloneFloat64Ptr(snapshot.DhwWaterFlowLpm),
		DhwDemandActive:          cloneBoilerBoolPtr(snapshot.DhwDemandActive),
		HeatingSwitchActive:      cloneBoilerBoolPtr(snapshot.HeatingSwitchActive),
		StorageLoadPumpPct:       CloneFloat64Ptr(snapshot.StorageLoadPumpPct),
		ModulationPct:            CloneFloat64Ptr(snapshot.ModulationPct),
		PrimaryCircuitFlowLpm:    CloneFloat64Ptr(snapshot.PrimaryCircuitFlowLpm),
		FlowTempDesiredC:         CloneFloat64Ptr(snapshot.FlowTempDesiredC),
		DhwTempDesiredC:          CloneFloat64Ptr(snapshot.DhwTempDesiredC),
		StateNumber:              cloneBoilerIntPtr(snapshot.StateNumber),
		DhwTemperatureC:          CloneFloat64Ptr(snapshot.DhwTemperatureC),
		DhwTargetTemperatureC:    CloneFloat64Ptr(snapshot.DhwTargetTemperatureC),
		DhwOperatingMode:         cloneBoilerIntPtr(snapshot.DhwOperatingMode),
		FlowsetHcMaxC:            CloneFloat64Ptr(snapshot.FlowsetHcMaxC),
		FlowsetHwcMaxC:           CloneFloat64Ptr(snapshot.FlowsetHwcMaxC),
		PartloadHcKW:             CloneFloat64Ptr(snapshot.PartloadHcKW),
		PartloadHwcKW:            CloneFloat64Ptr(snapshot.PartloadHwcKW),
		HeatingStatusRaw:         cloneBoilerIntPtr(snapshot.HeatingStatusRaw),
		DhwStatusRaw:             cloneBoilerIntPtr(snapshot.DhwStatusRaw),
		CentralHeatingHours:      CloneFloat64Ptr(snapshot.CentralHeatingHours),
		DhwHours:                 CloneFloat64Ptr(snapshot.DhwHours),
		CentralHeatingStarts:     cloneBoilerIntPtr(snapshot.CentralHeatingStarts),
		DhwStarts:                cloneBoilerIntPtr(snapshot.DhwStarts),
		PumpHours:                CloneFloat64Ptr(snapshot.PumpHours),
		FanHours:                 CloneFloat64Ptr(snapshot.FanHours),
		DeactivationsIFC:         cloneBoilerIntPtr(snapshot.DeactivationsIFC),
		DeactivationsTemplimiter: cloneBoilerIntPtr(snapshot.DeactivationsTemplimiter),
		InstallerMenuCode:        cloneBoilerIntPtr(snapshot.InstallerMenuCode),
		PhoneNumber:              CloneStringPtr(snapshot.PhoneNumber),
		HoursTillService:         cloneBoilerIntPtr(snapshot.HoursTillService),
	}
}

func cloneBoilerBoolPtr(value *bool) *bool {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}

func cloneBoilerIntPtr(value *int) *int {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}

func cloneBoilerModeIntPtr(mode *string) *int {
	if mode == nil {
		return nil
	}
	switch strings.ToUpper(strings.TrimSpace(*mode)) {
	case "OFF":
		value := 0
		return &value
	case "ON":
		value := 1
		return &value
	case "AUTO":
		value := 2
		return &value
	default:
		return nil
	}
}
