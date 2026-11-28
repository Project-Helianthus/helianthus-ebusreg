package owner

import (
	"context"
)

// refreshBoilerStatusB509 reads the B509 registers of one boiler status tier
// into snapshot and reports whether any read succeeded.
func (p *vaillantSemanticPoller) refreshBoilerStatusB509(ctx context.Context, boilerAddress byte, tier boilerStatusTier, snapshot *vaillantBoilerSnapshot) bool {
	if p == nil || snapshot == nil || boilerAddress == 0 {
		return false
	}

	switch tier {
	case boilerStatusTierFast:
		return p.refreshBoilerStatusB509Fast(ctx, boilerAddress, snapshot)
	case boilerStatusTierMedium:
		return p.refreshBoilerStatusB509Medium(ctx, boilerAddress, snapshot)
	case boilerStatusTierSlow:
		return p.refreshBoilerStatusB509Slow(ctx, boilerAddress, snapshot)
	}
	return false
}

// refreshBoilerStatusB509Fast reads the fast tier: flow evidence, pumps,
// flame, gas valve, fan, modulation, state and hot water demand.
func (p *vaillantSemanticPoller) refreshBoilerStatusB509Fast(ctx context.Context, boilerAddress byte, snapshot *vaillantBoilerSnapshot) bool {
	updated := false
	if raw, ok := p.readB509Value(ctx, boilerAddress, boilerB509FlowTemperature); ok {
		next, applied := applyOwnBoilerFlow(boilerB509FlowTemperature, p.boilerFlowReading(), raw, p.now())
		if applied {
			snapshot.FlowTemperatureC = cloneFloat64Ptr(next.celsius)
			p.mu.Lock()
			p.boilerFlowCached = next.cached
			p.boilerFlowObservedAt = next.observedAt
			p.mu.Unlock()
			updated = true
		}
	}
	if value := p.readB509OnOff(ctx, boilerAddress, boilerB509CentralHeatingPump); value != nil {
		snapshot.CentralHeatingPumpActive = value
		updated = true
	}
	if value := p.readB509Data2b(ctx, boilerAddress, boilerB509WaterPressure); value != nil {
		snapshot.WaterPressureBar = value
		updated = true
	}
	if value := p.readB509OnOff(ctx, boilerAddress, boilerB509FlameActive); value != nil {
		snapshot.FlameActive = value
		updated = true
	}
	if value := p.readB509OnOff(ctx, boilerAddress, boilerB509GasValveActive); value != nil {
		snapshot.GasValveActive = value
		updated = true
	}
	if value := p.readB509UINInt(ctx, boilerAddress, boilerB509FanSpeedRpm); value != nil {
		snapshot.FanSpeedRpm = value
		updated = true
	}
	if value := p.readB509SINScaled(ctx, boilerAddress, boilerB509ModulationPct, 10); value != nil {
		snapshot.ModulationPct = value
		updated = true
	}
	if value := p.readB509UCHInt(ctx, boilerAddress, boilerB509StateNumber); value != nil {
		snapshot.StateNumber = value
		updated = true
	}
	if value := p.readB509UCHFloat(ctx, boilerAddress, boilerB509DiverterValvePosition); value != nil {
		snapshot.DiverterValvePositionPct = value
		updated = true
	}
	if value := p.readB509OnOff(ctx, boilerAddress, boilerB509DhwDemandActive); value != nil {
		snapshot.DhwDemandActive = value
		updated = true
	}
	if value := p.readB509UnsignedHundredths(ctx, boilerAddress, boilerB509DhwWaterFlowLpm); value != nil {
		snapshot.DhwWaterFlowLpm = value
		updated = true
	}
	return updated
}

// refreshBoilerStatusB509Medium reads the medium tier: configured maxima
// and partial loads, desired temperatures, auxiliary pumps and flows.
func (p *vaillantSemanticPoller) refreshBoilerStatusB509Medium(ctx context.Context, boilerAddress byte, snapshot *vaillantBoilerSnapshot) bool {
	updated := false
	if value := p.readB509BoilerConfigFloat(ctx, boilerAddress, "maximum_heating_flow_temperature_setpoint"); value != nil {
		snapshot.FlowsetHcMaxC = value
		updated = true
	}
	if value := p.readB509BoilerConfigFloat(ctx, boilerAddress, "maximum_dhw_flow_temperature_setpoint"); value != nil {
		snapshot.FlowsetHwcMaxC = value
		updated = true
	}
	if value := p.readB509BoilerConfigFloat(ctx, boilerAddress, "heating_partial_load_power"); value != nil {
		snapshot.PartloadHcKW = value
		updated = true
	}
	if value := p.readB509BoilerConfigFloat(ctx, boilerAddress, "dhw_partial_load_power"); value != nil {
		snapshot.PartloadHwcKW = value
		updated = true
	}
	if value := p.readB509Percent0(ctx, boilerAddress, boilerB509StorageLoadPumpPct); value != nil {
		snapshot.StorageLoadPumpPct = value
		updated = true
	}
	if value := p.readB509Data2c(ctx, boilerAddress, boilerB509FlowTempDesiredC); value != nil {
		snapshot.FlowTempDesiredC = value
		updated = true
	}
	if value := p.readB509Data2c(ctx, boilerAddress, boilerB509DhwTempDesiredC); value != nil {
		snapshot.DhwTempDesiredC = value
		updated = true
	}
	if value := p.readB509OnOff(ctx, boilerAddress, boilerB509CirculationPumpActive); value != nil {
		snapshot.CirculationPumpActive = value
		updated = true
	}
	if value := p.readB509OnOff(ctx, boilerAddress, boilerB509ExternalPumpActive); value != nil {
		snapshot.ExternalPumpActive = value
		updated = true
	}
	if value := p.readB509OnOff(ctx, boilerAddress, boilerB509HeatingSwitchActive); value != nil {
		snapshot.HeatingSwitchActive = value
		updated = true
	}
	if value := p.readB509UINInt(ctx, boilerAddress, boilerB509TargetFanSpeedRpm); value != nil {
		snapshot.TargetFanSpeedRpm = value
		updated = true
	}
	if value := p.readB509SINScaled(ctx, boilerAddress, boilerB509IonisationVoltageUa, 10); value != nil {
		snapshot.IonisationVoltageUa = value
		updated = true
	}
	if value := p.readB509UnsignedHundredths(ctx, boilerAddress, boilerB509PrimaryCircuitFlowLpm); value != nil {
		snapshot.PrimaryCircuitFlowLpm = value
		updated = true
	}
	return updated
}

// refreshBoilerStatusB509Slow reads the slow tier: operating counters and
// the installer and maintenance configuration.
func (p *vaillantSemanticPoller) refreshBoilerStatusB509Slow(ctx context.Context, boilerAddress byte, snapshot *vaillantBoilerSnapshot) bool {
	updated := false
	if value := p.readB509Hoursum2(ctx, boilerAddress, boilerB509CentralHeatingHours); value != nil {
		snapshot.CentralHeatingHours = value
		updated = true
	}
	if value := p.readB509Hoursum2(ctx, boilerAddress, boilerB509DhwHours); value != nil {
		snapshot.DhwHours = value
		updated = true
	}
	if value := p.readB509UINInt(ctx, boilerAddress, boilerB509CentralHeatingStarts); value != nil {
		snapshot.CentralHeatingStarts = value
		updated = true
	}
	if value := p.readB509UINInt(ctx, boilerAddress, boilerB509DhwStarts); value != nil {
		snapshot.DhwStarts = value
		updated = true
	}
	if value := p.readB509Hoursum2(ctx, boilerAddress, boilerB509PumpHours); value != nil {
		snapshot.PumpHours = value
		updated = true
	}
	if value := p.readB509Hoursum2(ctx, boilerAddress, boilerB509FanHours); value != nil {
		snapshot.FanHours = value
		updated = true
	}
	if value := p.readB509UCHInt(ctx, boilerAddress, boilerB509DeactivationsIfc); value != nil {
		snapshot.DeactivationsIFC = value
		updated = true
	}
	if value := p.readB509UCHInt(ctx, boilerAddress, boilerB509DeactivationsLimit); value != nil {
		snapshot.DeactivationsTemplimiter = value
		updated = true
	}
	// Installer/maintenance config.
	if value := p.readB509UCHInt(ctx, boilerAddress, boilerB509InstallerMenuCode); value != nil {
		snapshot.InstallerMenuCode = value
		updated = true
	}
	if value := p.readB509PhoneBCD(ctx, boilerAddress, boilerB509PhoneNumber); value != nil {
		snapshot.PhoneNumber = value
		updated = true
	}
	if value := p.readB509Hoursum2Int(ctx, boilerAddress, boilerB509HoursTillService); value != nil {
		snapshot.HoursTillService = value
		updated = true
	}
	return updated
}

func mergeBoilerSnapshotNonDestructive(existing, incoming *vaillantBoilerSnapshot) *vaillantBoilerSnapshot {
	merged := cloneBoilerSnapshot(existing)
	if merged == nil {
		merged = &vaillantBoilerSnapshot{}
	}
	if incoming != nil {
		if incoming.FlowTemperatureC != nil {
			merged.FlowTemperatureC = cloneFloat64Ptr(incoming.FlowTemperatureC)
		}
		if incoming.CentralHeatingPumpActive != nil {
			merged.CentralHeatingPumpActive = cloneBoilerBoolPtr(incoming.CentralHeatingPumpActive)
		}
		if incoming.WaterPressureBar != nil {
			merged.WaterPressureBar = cloneFloat64Ptr(incoming.WaterPressureBar)
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
			merged.DiverterValvePositionPct = cloneFloat64Ptr(incoming.DiverterValvePositionPct)
		}
		if incoming.FanSpeedRpm != nil {
			merged.FanSpeedRpm = cloneBoilerIntPtr(incoming.FanSpeedRpm)
		}
		if incoming.TargetFanSpeedRpm != nil {
			merged.TargetFanSpeedRpm = cloneBoilerIntPtr(incoming.TargetFanSpeedRpm)
		}
		if incoming.IonisationVoltageUa != nil {
			merged.IonisationVoltageUa = cloneFloat64Ptr(incoming.IonisationVoltageUa)
		}
		if incoming.DhwWaterFlowLpm != nil {
			merged.DhwWaterFlowLpm = cloneFloat64Ptr(incoming.DhwWaterFlowLpm)
		}
		if incoming.DhwDemandActive != nil {
			merged.DhwDemandActive = cloneBoilerBoolPtr(incoming.DhwDemandActive)
		}
		if incoming.HeatingSwitchActive != nil {
			merged.HeatingSwitchActive = cloneBoilerBoolPtr(incoming.HeatingSwitchActive)
		}
		if incoming.StorageLoadPumpPct != nil {
			merged.StorageLoadPumpPct = cloneFloat64Ptr(incoming.StorageLoadPumpPct)
		}
		if incoming.ModulationPct != nil {
			merged.ModulationPct = cloneFloat64Ptr(incoming.ModulationPct)
		}
		if incoming.PrimaryCircuitFlowLpm != nil {
			merged.PrimaryCircuitFlowLpm = cloneFloat64Ptr(incoming.PrimaryCircuitFlowLpm)
		}
		if incoming.FlowTempDesiredC != nil {
			merged.FlowTempDesiredC = cloneFloat64Ptr(incoming.FlowTempDesiredC)
		}
		if incoming.DhwTempDesiredC != nil {
			merged.DhwTempDesiredC = cloneFloat64Ptr(incoming.DhwTempDesiredC)
		}
		if incoming.StateNumber != nil {
			merged.StateNumber = cloneBoilerIntPtr(incoming.StateNumber)
		}
		if incoming.DhwTemperatureC != nil {
			merged.DhwTemperatureC = cloneFloat64Ptr(incoming.DhwTemperatureC)
		}
		if incoming.DhwTargetTemperatureC != nil {
			merged.DhwTargetTemperatureC = cloneFloat64Ptr(incoming.DhwTargetTemperatureC)
		}
		if incoming.DhwOperatingMode != nil {
			merged.DhwOperatingMode = cloneBoilerIntPtr(incoming.DhwOperatingMode)
		}
		if incoming.FlowsetHcMaxC != nil {
			merged.FlowsetHcMaxC = cloneFloat64Ptr(incoming.FlowsetHcMaxC)
		}
		if incoming.FlowsetHwcMaxC != nil {
			merged.FlowsetHwcMaxC = cloneFloat64Ptr(incoming.FlowsetHwcMaxC)
		}
		if incoming.PartloadHcKW != nil {
			merged.PartloadHcKW = cloneFloat64Ptr(incoming.PartloadHcKW)
		}
		if incoming.PartloadHwcKW != nil {
			merged.PartloadHwcKW = cloneFloat64Ptr(incoming.PartloadHwcKW)
		}
		if incoming.HeatingStatusRaw != nil {
			merged.HeatingStatusRaw = cloneBoilerIntPtr(incoming.HeatingStatusRaw)
		}
		if incoming.DhwStatusRaw != nil {
			merged.DhwStatusRaw = cloneBoilerIntPtr(incoming.DhwStatusRaw)
		}
		if incoming.CentralHeatingHours != nil {
			merged.CentralHeatingHours = cloneFloat64Ptr(incoming.CentralHeatingHours)
		}
		if incoming.DhwHours != nil {
			merged.DhwHours = cloneFloat64Ptr(incoming.DhwHours)
		}
		if incoming.CentralHeatingStarts != nil {
			merged.CentralHeatingStarts = cloneBoilerIntPtr(incoming.CentralHeatingStarts)
		}
		if incoming.DhwStarts != nil {
			merged.DhwStarts = cloneBoilerIntPtr(incoming.DhwStarts)
		}
		if incoming.PumpHours != nil {
			merged.PumpHours = cloneFloat64Ptr(incoming.PumpHours)
		}
		if incoming.FanHours != nil {
			merged.FanHours = cloneFloat64Ptr(incoming.FanHours)
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
			merged.PhoneNumber = cloneStringPtr(incoming.PhoneNumber)
		}
		if incoming.HoursTillService != nil {
			merged.HoursTillService = cloneBoilerIntPtr(incoming.HoursTillService)
		}
	}

	// Closed decision: do not map GG=0x02 RR=0x0008 as boiler return temperature.
	// Keep return temperature unset in B524-derived snapshots.
	merged.ReturnTemperatureC = nil
	return merged
}
