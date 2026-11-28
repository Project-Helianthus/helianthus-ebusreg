package output

// MCPZonesFromGraphQL maps the Zones output record to its MCP form.
func MCPZonesFromGraphQL(zones []Zone) []MCPZone {
	if zones == nil {
		return nil
	}
	out := make([]MCPZone, len(zones))
	for i, zone := range zones {
		out[i] = MCPZone{
			ID:   zone.ID,
			Name: zone.Name,
			State: MCPZoneState{
				CurrentTempC:       CloneFloatPtr(zone.State.CurrentTempC),
				CurrentHumidityPct: CloneFloatPtr(zone.State.CurrentHumidityPct),
				HvacAction:         zone.State.HvacAction,
				SpecialFunction:    zone.State.SpecialFunction,
				HeatingDemandPct:   CloneFloatPtr(zone.State.HeatingDemandPct),
				ValvePositionPct:   CloneFloatPtr(zone.State.ValvePositionPct),
				Enumerations:       cloneEnumerations(zone.State.Enumerations),
			},
			Config: MCPZoneConfig{
				OperatingMode:              zone.Config.OperatingMode,
				Preset:                     zone.Config.Preset,
				TargetTempC:                CloneFloatPtr(zone.Config.TargetTempC),
				AllowedModes:               append([]string(nil), zone.Config.AllowedModes...),
				CircuitType:                zone.Config.CircuitType,
				AssociatedCircuit:          CloneIntPtr(zone.Config.AssociatedCircuit),
				RoomTemperatureZoneMapping: CloneIntPtr(zone.Config.RoomTemperatureZoneMapping),
				QuickVeto:                  zone.Config.QuickVeto,
				QuickVetoSetpointC:         CloneFloatPtr(zone.Config.QuickVetoSetpointC),
				QuickVetoDurationH:         CloneFloatPtr(zone.Config.QuickVetoDurationH),
				QuickVetoExpiry:            zone.Config.QuickVetoExpiry,
				HolidayStartDate:           zone.Config.HolidayStartDate,
				HolidayEndDate:             zone.Config.HolidayEndDate,
				HolidaySetpointC:           CloneFloatPtr(zone.Config.HolidaySetpointC),
				HolidayStartTime:           zone.Config.HolidayStartTime,
				HolidayEndTime:             zone.Config.HolidayEndTime,
			},
		}
	}
	return out
}

// MCPDHWFromGraphQL maps the DHW output record to its MCP form.
func MCPDHWFromGraphQL(status *DhwStatus) *MCPDhwStatus {
	if status == nil {
		return nil
	}
	return &MCPDhwStatus{
		State: MCPDhwState{
			CurrentTempC:     CloneFloatPtr(status.State.CurrentTempC),
			SpecialFunction:  status.State.SpecialFunction,
			HeatingDemandPct: CloneFloatPtr(status.State.HeatingDemandPct),
			Enumerations:     cloneEnumerations(status.State.Enumerations),
		},
		Config: MCPDhwConfig{
			OperatingMode:    status.Config.OperatingMode,
			Preset:           status.Config.Preset,
			TargetTempC:      CloneFloatPtr(status.Config.TargetTempC),
			HolidayStartDate: status.Config.HolidayStartDate,
			HolidayEndDate:   status.Config.HolidayEndDate,
		},
	}
}

// MCPCircuitsFromGraphQL maps the Circuits output record to its MCP form.
func MCPCircuitsFromGraphQL(circuits []CircuitStatus) []MCPCircuitStatus {
	if circuits == nil {
		return nil
	}
	out := make([]MCPCircuitStatus, len(circuits))
	for i, circuit := range circuits {
		out[i] = MCPCircuitStatus{
			Index:       circuit.Index,
			CircuitType: circuit.CircuitType,
			HasMixer:    circuit.HasMixer,
			State: MCPCircuitState{
				PumpActive:       CloneBoolPtr(circuit.State.PumpActive),
				MixerPositionPct: CloneFloatPtr(circuit.State.MixerPositionPct),
				FlowTemperatureC: CloneFloatPtr(circuit.State.FlowTemperatureC),
				FlowSetpointC:    CloneFloatPtr(circuit.State.FlowSetpointC),
				CalcFlowTempC:    CloneFloatPtr(circuit.State.CalcFlowTempC),
				CircuitState:     circuit.State.CircuitState,
				Humidity:         CloneFloatPtr(circuit.State.Humidity),
				DewPoint:         CloneFloatPtr(circuit.State.DewPoint),
				PumpHours:        CloneFloatPtr(circuit.State.PumpHours),
				PumpStarts:       CloneIntPtr(circuit.State.PumpStarts),
			},
			Config: MCPCircuitConfig{
				HeatingCurve:        CloneFloatPtr(circuit.Config.HeatingCurve),
				FlowTempMaxC:        CloneFloatPtr(circuit.Config.FlowTempMaxC),
				FlowTempMinC:        CloneFloatPtr(circuit.Config.FlowTempMinC),
				SummerLimitC:        CloneFloatPtr(circuit.Config.SummerLimitC),
				SpecialFunctionMode: CloneIntPtr(circuit.Config.SpecialFunctionMode),
				RoomTempControl:     circuit.Config.RoomTempControl,
				CoolingEnabled:      CloneBoolPtr(circuit.Config.CoolingEnabled),
				Enumerations:        cloneEnumerations(circuit.Config.Enumerations),
			},
			ManagingDevice: MCPManagingDevice{
				Role:     string(circuit.ManagingDevice.Role),
				DeviceID: CloneStringPtr(circuit.ManagingDevice.DeviceID),
				Address:  CloneIntPtr(circuit.ManagingDevice.Address),
			},
		}
	}
	return out
}

// MCPRadioDevicesFromGraphQL maps the RadioDevices output record to its MCP form.
func MCPRadioDevicesFromGraphQL(devices []RadioDevice) []MCPRadioDevice {
	if devices == nil {
		return nil
	}
	out := make([]MCPRadioDevice, len(devices))
	for i, device := range devices {
		out[i] = MCPRadioDevice{
			Group:                device.Group,
			Instance:             device.Instance,
			SlotMode:             device.SlotMode,
			DeviceConnected:      CloneBoolPtr(device.DeviceConnected),
			DeviceClassAddress:   CloneIntPtr(device.DeviceClassAddress),
			DeviceModel:          device.DeviceModel,
			FirmwareVersion:      CloneStringPtr(device.FirmwareVersion),
			HardwareIdentifier:   CloneIntPtr(device.HardwareIdentifier),
			RemoteControlAddress: CloneIntPtr(device.RemoteControlAddress),
			DevicePaired:         CloneBoolPtr(device.DevicePaired),
			ReceptionStrength:    CloneIntPtr(device.ReceptionStrength),
			ZoneAssignment:       CloneIntPtr(device.ZoneAssignment),
			RoomTemperatureC:     CloneFloatPtr(device.RoomTemperatureC),
			RoomHumidityPct:      CloneFloatPtr(device.RoomHumidityPct),
		}
	}
	return out
}

// MCPFunctionModuleSemanticModeFromGraphQL maps the FunctionModuleSemanticMode output record to its MCP form.
func MCPFunctionModuleSemanticModeFromGraphQL(mode FunctionModuleSemanticMode) MCPFunctionModuleSemanticMode {
	switch mode {
	case FunctionModuleSemanticModeInterpreted:
		return MCPFunctionModuleSemanticModeInterpreted
	case FunctionModuleSemanticModeGPIOOnly:
		return MCPFunctionModuleSemanticModeGPIOOnly
	default:
		return MCPFunctionModuleSemanticModeAbsent
	}
}

// MCPFunctionModuleInterpretationFromGraphQL maps the FunctionModuleInterpretation output record to its MCP form.
func MCPFunctionModuleInterpretationFromGraphQL(verdict FunctionModuleInterpretation) MCPFunctionModuleInterpretation {
	out := MCPFunctionModuleInterpretation{
		Mode:             MCPFunctionModuleSemanticMode(verdict.Mode),
		EvidenceRevision: verdict.EvidenceRevision,
	}
	if verdict.DegradedReason != "" {
		reason := MCPFunctionModuleSemanticDegradedReason(verdict.DegradedReason)
		out.DegradedReason = &reason
	}
	return out
}

// MCPSolarFromGraphQL maps the Solar output record to its MCP form.
func MCPSolarFromGraphQL(status *SolarStatus) *MCPSolarStatus {
	if status == nil {
		return nil
	}
	return &MCPSolarStatus{
		CollectorTemperatureC: CloneFloatPtr(status.CollectorTemperatureC),
		ReturnTemperatureC:    CloneFloatPtr(status.ReturnTemperatureC),
		PumpActive:            CloneBoolPtr(status.PumpActive),
		CurrentYield:          CloneFloatPtr(status.CurrentYield),
		PumpHours:             CloneFloatPtr(status.PumpHours),
		SolarEnabled:          CloneBoolPtr(status.SolarEnabled),
		FunctionMode:          CloneBoolPtr(status.FunctionMode),
	}
}

// MCPCylindersFromGraphQL maps the Cylinders output record to its MCP form.
func MCPCylindersFromGraphQL(cylinders []CylinderStatus) []MCPCylinderStatus {
	if cylinders == nil {
		return nil
	}
	out := make([]MCPCylinderStatus, len(cylinders))
	for i, cylinder := range cylinders {
		out[i] = MCPCylinderStatus{
			Index:             cylinder.Index,
			TemperatureC:      CloneFloatPtr(cylinder.TemperatureC),
			MaxSetpointC:      CloneFloatPtr(cylinder.MaxSetpointC),
			ChargeHysteresisC: CloneFloatPtr(cylinder.ChargeHysteresisC),
			ChargeOffsetC:     CloneFloatPtr(cylinder.ChargeOffsetC),
		}
	}
	return out
}
