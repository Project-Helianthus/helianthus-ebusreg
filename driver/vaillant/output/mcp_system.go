package output

// MCPEnergyTotalsFromGraphQL maps the EnergyTotals output record to its MCP form.
func MCPEnergyTotalsFromGraphQL(totals *EnergyTotals) *MCPEnergyTotals {
	if totals == nil {
		return nil
	}
	return &MCPEnergyTotals{
		Gas:                 mapEnergyChannel(totals.Gas),
		Electric:            mapEnergyChannel(totals.Electric),
		Solar:               mapEnergyChannel(totals.Solar),
		registerDirections:  cloneEnergyDirections(totals.registerDirections),
		catalogueDirections: cloneEnergyDirections(totals.catalogueDirections),
		residualFields:      cloneEnergyResiduals(totals.residualFields),
		candidateFields:     cloneEnergyResiduals(totals.candidateFields),
	}
}

// MCPBoilerStatusFromGraphQL maps the BoilerStatus output record to its MCP form.
func MCPBoilerStatusFromGraphQL(status *BoilerStatus) *MCPBoilerStatus {
	if status == nil {
		return nil
	}
	out := &MCPBoilerStatus{
		State: &MCPBoilerState{
			FlowTemperatureC:         CloneFloatPtr(status.State.FlowTemperatureC),
			ReturnTemperatureC:       CloneFloatPtr(status.State.ReturnTemperatureC),
			CentralHeatingPumpActive: CloneBoolPtr(status.State.CentralHeatingPumpActive),
			WaterPressureBar:         CloneFloatPtr(status.State.WaterPressureBar),
			ExternalPumpActive:       CloneBoolPtr(status.State.ExternalPumpActive),
			CirculationPumpActive:    CloneBoolPtr(status.State.CirculationPumpActive),
			GasValveActive:           CloneBoolPtr(status.State.GasValveActive),
			FlameActive:              CloneBoolPtr(status.State.FlameActive),
			DiverterValvePositionPct: CloneFloatPtr(status.State.DiverterValvePositionPct),
			FanSpeedRpm:              CloneIntPtr(status.State.FanSpeedRpm),
			TargetFanSpeedRpm:        CloneIntPtr(status.State.TargetFanSpeedRpm),
			IonisationVoltageUa:      CloneFloatPtr(status.State.IonisationVoltageUa),
			DhwWaterFlowLpm:          CloneFloatPtr(status.State.DhwWaterFlowLpm),
			DhwDemandActive:          CloneBoolPtr(status.State.DhwDemandActive),
			HeatingSwitchActive:      CloneBoolPtr(status.State.HeatingSwitchActive),
			StorageLoadPumpPct:       CloneFloatPtr(status.State.StorageLoadPumpPct),
			ModulationPct:            CloneFloatPtr(status.State.ModulationPct),
			PrimaryCircuitFlowLpm:    CloneFloatPtr(status.State.PrimaryCircuitFlowLpm),
			FlowTempDesiredC:         CloneFloatPtr(status.State.FlowTempDesiredC),
			DhwTempDesiredC:          CloneFloatPtr(status.State.DhwTempDesiredC),
			StateNumber:              CloneIntPtr(status.State.StateNumber),
			DhwTemperatureC:          CloneFloatPtr(status.State.DhwTemperatureC),
			DhwTargetTemperatureC:    CloneFloatPtr(status.State.DhwTargetTemperatureC),
		},
		Config: &MCPBoilerConfig{
			DhwOperatingMode:  CloneStringPtr(status.Config.DhwOperatingMode),
			FlowsetHcMaxC:     CloneFloatPtr(status.Config.FlowsetHcMaxC),
			FlowsetHwcMaxC:    CloneFloatPtr(status.Config.FlowsetHwcMaxC),
			PartloadHcKW:      CloneFloatPtr(status.Config.PartloadHcKW),
			PartloadHwcKW:     CloneFloatPtr(status.Config.PartloadHwcKW),
			InstallerMenuCode: CloneIntPtr(status.Config.InstallerMenuCode),
			PhoneNumber:       CloneStringPtr(status.Config.PhoneNumber),
			HoursTillService:  CloneIntPtr(status.Config.HoursTillService),
		},
		Diagnostics: &MCPBoilerDiagnostics{
			HeatingStatusRaw:         CloneIntPtr(status.Diagnostics.HeatingStatusRaw),
			DhwStatusRaw:             CloneIntPtr(status.Diagnostics.DhwStatusRaw),
			CentralHeatingHours:      CloneFloatPtr(status.Diagnostics.CentralHeatingHours),
			DhwHours:                 CloneFloatPtr(status.Diagnostics.DhwHours),
			CentralHeatingStarts:     CloneIntPtr(status.Diagnostics.CentralHeatingStarts),
			DhwStarts:                CloneIntPtr(status.Diagnostics.DhwStarts),
			PumpHours:                CloneFloatPtr(status.Diagnostics.PumpHours),
			FanHours:                 CloneFloatPtr(status.Diagnostics.FanHours),
			DeactivationsIFC:         CloneIntPtr(status.Diagnostics.DeactivationsIFC),
			DeactivationsTemplimiter: CloneIntPtr(status.Diagnostics.DeactivationsTemplimiter),
		},
		selectors:  cloneBoilerSelectors(status.selectors),
		operations: cloneRawOperations(status.operations),
	}
	return out
}

// MCPSchedulesFromGraphQL maps the Schedules output record to its MCP form.
func MCPSchedulesFromGraphQL(status *ScheduleStatus) *MCPScheduleStatus {
	if status == nil {
		return nil
	}
	out := &MCPScheduleStatus{}
	if len(status.Programs) > 0 {
		out.Programs = make([]MCPScheduleProgram, len(status.Programs))
		for i, prog := range status.Programs {
			mp := MCPScheduleProgram{
				Zone: prog.Zone,
				HC:   prog.HC,
			}
			if prog.Config != nil {
				mp.Config = &MCPScheduleConfig{
					MaxSlots:       prog.Config.MaxSlots,
					TimeResolution: prog.Config.TimeResolution,
					MinDuration:    prog.Config.MinDuration,
					HasTemperature: prog.Config.HasTemperature,
					TempSlots:      prog.Config.TempSlots,
					MinTempC:       CloneFloatPtr(prog.Config.MinTempC),
					MaxTempC:       CloneFloatPtr(prog.Config.MaxTempC),
				}
			}
			if len(prog.SlotsUsed) > 0 {
				mp.SlotsUsed = make([]int, len(prog.SlotsUsed))
				copy(mp.SlotsUsed, prog.SlotsUsed)
			}
			if len(prog.Days) > 0 {
				mp.Days = make([]MCPScheduleDayProgram, len(prog.Days))
				for j, day := range prog.Days {
					md := MCPScheduleDayProgram{Weekday: day.Weekday}
					if len(day.Slots) > 0 {
						md.Slots = make([]MCPScheduleTimerSlot, len(day.Slots))
						for k, slot := range day.Slots {
							md.Slots[k] = MCPScheduleTimerSlot{
								StartHour:      slot.StartHour,
								StartMinute:    slot.StartMinute,
								EndHour:        slot.EndHour,
								EndMinute:      slot.EndMinute,
								TemperatureC:   CloneFloatPtr(slot.TemperatureC),
								TemperatureRaw: CloneIntPtr(slot.TemperatureRaw),
							}
						}
					}
					mp.Days[j] = md
				}
			}
			out.Programs[i] = mp
		}
	}
	return out
}

// MCPAdapterHardwareInfoFromGraphQL maps the AdapterHardwareInfo output record to its MCP form.
func MCPAdapterHardwareInfoFromGraphQL(info *AdapterHardwareInfo) *MCPAdapterHardwareInfo {
	if info == nil {
		return nil
	}
	out := &MCPAdapterHardwareInfo{
		FirmwareVersion:    info.FirmwareVersion,
		FirmwareChecksum:   info.FirmwareChecksum,
		BootloaderVersion:  info.BootloaderVersion,
		BootloaderChecksum: info.BootloaderChecksum,
		HardwareID:         info.HardwareID,
		HardwareConfig:     info.HardwareConfig,
		Features:           info.Features,
		Jumpers:            info.Jumpers,
		IsWiFi:             info.IsWiFi,
		IsEthernet:         info.IsEthernet,
		TemperatureC:       CloneFloatPtr(info.TemperatureC),
		SupplyVoltageMV:    CloneIntPtr(info.SupplyVoltageMV),
		BusVoltageMaxDV:    CloneIntPtr(info.BusVoltageMaxDV),
		BusVoltageMinDV:    CloneIntPtr(info.BusVoltageMinDV),
		ResetCause:         CloneStringPtr(info.ResetCause),
		WiFiRSSIDBm:        CloneIntPtr(info.WiFiRSSIDBm),
		VersionResponseLen: info.VersionResponseLen,
		InfoSupported:      info.InfoSupported,
	}
	if info.JumperFlags != nil {
		out.JumperFlags = make([]string, len(info.JumperFlags))
		copy(out.JumperFlags, info.JumperFlags)
	}
	if info.ResetCauseCode != nil {
		v := *info.ResetCauseCode
		out.ResetCauseCode = &v
	}
	if info.RestartCount != nil {
		v := *info.RestartCount
		out.RestartCount = &v
	}
	if info.LastIdentityQuery != nil {
		s := info.LastIdentityQuery.UTC().Format("2006-01-02T15:04:05Z")
		out.LastIdentityQuery = &s
	}
	if info.LastTelemetryQuery != nil {
		s := info.LastTelemetryQuery.UTC().Format("2006-01-02T15:04:05Z")
		out.LastTelemetryQuery = &s
	}
	return out
}

// MCPSystemFromGraphQL maps the System output record to its MCP form.
func MCPSystemFromGraphQL(status *SystemStatus) *MCPSystemStatus {
	if status == nil {
		return nil
	}
	return &MCPSystemStatus{
		ventilation:    cloneVentilation(status.ventilation),
		catalogue:      cloneFamilyCatalogue(status.catalogue),
		history:        cloneHistoryPtr(status.history),
		functionModule: cloneFunctionModulePtr(status.functionModule),
		legacyPlanes:   cloneLegacyPlanePtr(status.legacyPlanes),
		State: &MCPSystemState{
			SystemOff:                      CloneBoolPtr(status.State.SystemOff),
			SystemWaterPressure:            CloneFloatPtr(status.State.SystemWaterPressure),
			SystemFlowTemperature:          CloneFloatPtr(status.State.SystemFlowTemperature),
			OutdoorTemperature:             CloneFloatPtr(status.State.OutdoorTemperature),
			OutdoorTemperatureDailyAverage: CloneFloatPtr(status.State.OutdoorTemperatureDailyAverage),
			MaintenanceDue:                 CloneBoolPtr(status.State.MaintenanceDue),
			HwcCylinderTemperatureTop:      CloneFloatPtr(status.State.HwcCylinderTemperatureTop),
			HwcCylinderTemperatureBottom:   CloneFloatPtr(status.State.HwcCylinderTemperatureBottom),
		},
		Config: &MCPSystemConfig{
			AdaptiveHeatingCurve:         CloneBoolPtr(status.Config.AdaptiveHeatingCurve),
			AlternativePoint:             CloneFloatPtr(status.Config.AlternativePoint),
			HeatingCircuitBivalencePoint: CloneFloatPtr(status.Config.HeatingCircuitBivalencePoint),
			DhwBivalencePoint:            CloneFloatPtr(status.Config.DhwBivalencePoint),
			DhwBivalenceBound:            mcpConfigBound(status.Config.DhwBivalenceBound),
			HcEmergencyTemperature:       CloneFloatPtr(status.Config.HcEmergencyTemperature),
			HwcMaxFlowTempDesired:        CloneFloatPtr(status.Config.HwcMaxFlowTempDesired),
			MaxRoomHumidity:              CloneIntPtr(status.Config.MaxRoomHumidity),
			MaintenanceDate:              CloneStringPtr(status.Config.MaintenanceDate),
			InstallerName:                CloneStringPtr(status.Config.InstallerName),
			InstallerPhone:               CloneStringPtr(status.Config.InstallerPhone),
			InstallerMenuCode:            CloneIntPtr(status.Config.InstallerMenuCode),
		},
		Properties: &MCPSystemProperties{
			SystemScheme:                    CloneIntPtr(status.Properties.SystemScheme),
			FunctionModuleConfigurationCode: CloneIntPtr(status.Properties.FunctionModuleConfigurationCode),
		},
	}
}

func mcpConfigBound(bound *ConfigBound) *MCPConfigBound {
	if bound == nil {
		return nil
	}
	return &MCPConfigBound{
		NarrowedOnly: bound.NarrowedOnly,
		MinC:         bound.MinC,
		MaxC:         bound.MaxC,
		Source:       bound.Source,
		Reason:       bound.Reason,
		Profile:      bound.Profile,
	}
}

func mapEnergyChannel(channel EnergyChannel) MCPEnergyChannel {
	return MCPEnergyChannel{
		DHW:     mapEnergySeries(channel.DHW),
		Climate: mapEnergySeries(channel.Climate),
	}
}

func mapEnergySeries(series EnergySeries) MCPEnergySeries {
	out := MCPEnergySeries{
		Today:     series.Today,
		TodayMeta: mapEnergyPointMeta(series.TodayMeta),
	}
	if len(series.Yearly) > 0 {
		out.Yearly = make([]float64, len(series.Yearly))
		copy(out.Yearly, series.Yearly)
	}
	if len(series.Monthly) > 0 {
		out.Monthly = make([]float64, len(series.Monthly))
		copy(out.Monthly, series.Monthly)
	}
	if len(series.YearlyMeta) > 0 {
		out.YearlyMeta = make([]MCPEnergyPointMeta, len(series.YearlyMeta))
		for i, meta := range series.YearlyMeta {
			out.YearlyMeta[i] = mapEnergyPointMeta(meta)
		}
	}
	if len(series.MonthlyMeta) > 0 {
		out.MonthlyMeta = make([]MCPEnergyPointMeta, len(series.MonthlyMeta))
		for i, meta := range series.MonthlyMeta {
			out.MonthlyMeta[i] = mapEnergyPointMeta(meta)
		}
	}
	return out
}

func mapEnergyPointMeta(meta EnergyPointMeta) MCPEnergyPointMeta {
	return MCPEnergyPointMeta{
		FreshnessState:  string(meta.FreshnessState),
		Provenance:      string(meta.Provenance),
		LastObservedUTC: meta.LastObservedUTC,
		AgeSeconds:      meta.AgeSeconds,
		Stale:           meta.Stale,
	}
}
