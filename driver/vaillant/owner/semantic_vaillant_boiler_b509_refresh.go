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
	// Acquisition retains native evidence; only the hub result supplies flow.
	_, updated := p.readB509Value(ctx, boilerAddress, boiler_b509_flow_temperature)
	if value := p.readB509OnOff(ctx, boilerAddress, boiler_b509_central_heating_pump); value != nil {
		snapshot.CentralHeatingPumpActive = value
		updated = true
	}
	if value := p.readB509Data2b(ctx, boilerAddress, boiler_b509_water_pressure); value != nil {
		snapshot.WaterPressureBar = value
		updated = true
	}
	if value := p.readB509OnOff(ctx, boilerAddress, boiler_b509_flame_active); value != nil {
		snapshot.FlameActive = value
		updated = true
	}
	if value := p.readB509OnOff(ctx, boilerAddress, boiler_b509_gas_valve_active); value != nil {
		snapshot.GasValveActive = value
		updated = true
	}
	if value := p.readB509UINInt(ctx, boilerAddress, boiler_b509_fan_speed_rpm); value != nil {
		snapshot.FanSpeedRpm = value
		updated = true
	}
	if value := p.readB509SINScaled(ctx, boilerAddress, boiler_b509_modulation_pct, 10); value != nil {
		snapshot.ModulationPct = value
		updated = true
	}
	if value := p.readB509UCHInt(ctx, boilerAddress, boiler_b509_state_number); value != nil {
		snapshot.StateNumber = value
		updated = true
	}
	if value := p.readB509UCHFloat(ctx, boilerAddress, boiler_b509_diverter_valve_position); value != nil {
		snapshot.DiverterValvePositionPct = value
		updated = true
	}
	if value := p.readB509OnOff(ctx, boilerAddress, boiler_b509_dhw_demand_active); value != nil {
		snapshot.DhwDemandActive = value
		updated = true
	}
	if value := p.readB509UnsignedHundredths(ctx, boilerAddress, boiler_b509_dhw_water_flow_lpm); value != nil {
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
	if value := p.readB509Percent0(ctx, boilerAddress, boiler_b509_storage_load_pump_pct); value != nil {
		snapshot.StorageLoadPumpPct = value
		updated = true
	}
	if value := p.readB509Data2c(ctx, boilerAddress, boiler_b509_flow_temp_desired_c); value != nil {
		snapshot.FlowTempDesiredC = value
		updated = true
	}
	if value := p.readB509Data2c(ctx, boilerAddress, boiler_b509_dhw_temp_desired_c); value != nil {
		snapshot.DhwTempDesiredC = value
		updated = true
	}
	if value := p.readB509OnOff(ctx, boilerAddress, boiler_b509_circulation_pump_active); value != nil {
		snapshot.CirculationPumpActive = value
		updated = true
	}
	if value := p.readB509OnOff(ctx, boilerAddress, boiler_b509_external_pump_active); value != nil {
		snapshot.ExternalPumpActive = value
		updated = true
	}
	if value := p.readB509OnOff(ctx, boilerAddress, boiler_b509_heating_switch_active); value != nil {
		snapshot.HeatingSwitchActive = value
		updated = true
	}
	if value := p.readB509UINInt(ctx, boilerAddress, boiler_b509_target_fan_speed_rpm); value != nil {
		snapshot.TargetFanSpeedRpm = value
		updated = true
	}
	if value := p.readB509SINScaled(ctx, boilerAddress, boiler_b509_ionisation_voltage_ua, 10); value != nil {
		snapshot.IonisationVoltageUa = value
		updated = true
	}
	if value := p.readB509UnsignedHundredths(ctx, boilerAddress, boiler_b509_primary_circuit_flow_lpm); value != nil {
		snapshot.PrimaryCircuitFlowLpm = value
		updated = true
	}
	return updated
}

// refreshBoilerStatusB509Slow reads the slow tier: operating counters and
// the installer and maintenance configuration.
func (p *vaillantSemanticPoller) refreshBoilerStatusB509Slow(ctx context.Context, boilerAddress byte, snapshot *vaillantBoilerSnapshot) bool {
	updated := false
	if value := p.readB509Hoursum2(ctx, boilerAddress, boiler_b509_central_heating_hours); value != nil {
		snapshot.CentralHeatingHours = value
		updated = true
	}
	if value := p.readB509Hoursum2(ctx, boilerAddress, boiler_b509_dhw_hours); value != nil {
		snapshot.DhwHours = value
		updated = true
	}
	if value := p.readB509UINInt(ctx, boilerAddress, boiler_b509_central_heating_starts); value != nil {
		snapshot.CentralHeatingStarts = value
		updated = true
	}
	if value := p.readB509UINInt(ctx, boilerAddress, boiler_b509_dhw_starts); value != nil {
		snapshot.DhwStarts = value
		updated = true
	}
	if value := p.readB509Hoursum2(ctx, boilerAddress, boiler_b509_pump_hours); value != nil {
		snapshot.PumpHours = value
		updated = true
	}
	if value := p.readB509Hoursum2(ctx, boilerAddress, boiler_b509_fan_hours); value != nil {
		snapshot.FanHours = value
		updated = true
	}
	if value := p.readB509UCHInt(ctx, boilerAddress, boiler_b509_deactivations_ifc); value != nil {
		snapshot.DeactivationsIFC = value
		updated = true
	}
	if value := p.readB509UCHInt(ctx, boilerAddress, boiler_b509_deactivations_limit); value != nil {
		snapshot.DeactivationsTemplimiter = value
		updated = true
	}
	// Installer/maintenance config.
	if value := p.readB509UCHInt(ctx, boilerAddress, boiler_b509_installer_menu_code); value != nil {
		snapshot.InstallerMenuCode = value
		updated = true
	}
	if value := p.readB509PhoneBCD(ctx, boilerAddress, boiler_b509_phone_number); value != nil {
		snapshot.PhoneNumber = value
		updated = true
	}
	if value := p.readB509Hoursum2Int(ctx, boilerAddress, boiler_b509_hours_till_service); value != nil {
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
