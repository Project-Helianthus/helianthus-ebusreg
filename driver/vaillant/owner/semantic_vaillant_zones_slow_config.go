package owner

import (
	"context"
	"fmt"
)

func (p *vaillantSemanticPoller) refreshZoneSlowConfigFields(ctx context.Context, instance byte) (*vaillantZoneSnapshot, bool) {
	incoming := &vaillantZoneSnapshot{}
	liveReadSuccess := false

	var targetPtr *float64
	if value, ok := p.readB524Float32LE(ctx, localZones.opcode, localZones.group, instance, zoneTargetTemp); ok {
		target := value
		targetPtr = &target
		liveReadSuccess = true
	}
	if targetPtr == nil {
		if value, ok := p.readB524Float32LE(ctx, localZones.opcode, localZones.group, instance, zoneFallbackManualTemp); ok {
			target := value
			targetPtr = &target
			liveReadSuccess = true
		}
	}

	var humidity *float64
	if value, ok := p.readB524Float32LE(ctx, localZones.opcode, localZones.group, instance, zoneCurrentHumidity); ok {
		currentHumidity := value
		humidity = &currentHumidity
		liveReadSuccess = true
	}

	zoneOpMode, zoneOpModeOK := p.readB524Uint16(ctx, localZones.opcode, localZones.group, instance, zoneHeatingOpMode)
	if zoneOpModeOK {
		liveReadSuccess = true
	}

	var qvTempPtr, qvDurPtr *float64
	if value, ok := p.readB524Float32LE(ctx, localZones.opcode, localZones.group, instance, zoneQuickVetoTemp); ok {
		v := value
		qvTempPtr = &v
		liveReadSuccess = true
	}
	if value, ok := p.readB524Float32LE(ctx, localZones.opcode, localZones.group, instance, zoneQuickVetoDuration); ok {
		v := value
		qvDurPtr = &v
		liveReadSuccess = true
	}
	var qvEndTime, qvEndDate string
	if raw, ok := p.readB524Value(ctx, localZones.opcode, localZones.group, instance, zoneQuickVetoEndTime); ok && len(raw) >= 2 {
		qvEndTime = fmt.Sprintf("%02d:%02d", raw[0], raw[1])
		liveReadSuccess = true
	}
	if raw, ok := p.readB524Value(ctx, localZones.opcode, localZones.group, instance, zoneQuickVetoEndDate); ok && len(raw) >= 3 {
		year := 2000 + int(raw[2])
		qvEndDate = fmt.Sprintf("%04d-%02d-%02d", year, raw[1], raw[0])
		liveReadSuccess = true
	}

	var holidayStartDate, holidayEndDate, holidayStartTime, holidayEndTime string
	var holidaySetpointPtr *float64
	if raw, ok := p.readB524Value(ctx, localZones.opcode, localZones.group, instance, zoneHolidayStartDate); ok && len(raw) >= 3 {
		if date := decodeB524DateSuppressSentinel(raw); date != "" {
			holidayStartDate = date
		}
		liveReadSuccess = true
	}
	if raw, ok := p.readB524Value(ctx, localZones.opcode, localZones.group, instance, zoneHolidayEndDate); ok && len(raw) >= 3 {
		if date := decodeB524DateSuppressSentinel(raw); date != "" {
			holidayEndDate = date
		}
		liveReadSuccess = true
	}
	if value, ok := p.readB524Float32LE(ctx, localZones.opcode, localZones.group, instance, zoneHolidaySetpoint); ok {
		v := value
		holidaySetpointPtr = &v
		liveReadSuccess = true
	}
	if raw, ok := p.readB524Value(ctx, localZones.opcode, localZones.group, instance, zoneHolidayEndTime); ok && len(raw) >= 2 {
		holidayEndTime = fmt.Sprintf("%02d:%02d", raw[0], raw[1])
		liveReadSuccess = true
	}
	if raw, ok := p.readB524Value(ctx, localZones.opcode, localZones.group, instance, zoneHolidayStartTime); ok && len(raw) >= 2 {
		holidayStartTime = fmt.Sprintf("%02d:%02d", raw[0], raw[1])
		liveReadSuccess = true
	}

	zoneRoomTemperatureZoneMappingRaw, zoneRoomTemperatureZoneMappingRawOK := p.readB524Uint16(ctx, localZones.opcode, localZones.group, instance, zoneRoomTemperatureZoneMappingRawRegister)
	if zoneRoomTemperatureZoneMappingRawOK {
		liveReadSuccess = true
	}
	circuitInstance := resolveAssociatedCircuitInstance(zoneRoomTemperatureZoneMappingRaw, instance)
	circuitType, hasCircuitType := p.readB524Uint16(ctx, localCircuits.opcode, localCircuits.group, circuitInstance, circuitTypeRegister)
	if hasCircuitType {
		liveReadSuccess = true
	}
	var associatedCircuitRaw *uint16
	if zoneRoomTemperatureZoneMappingRaw != nil {
		value := uint16(circuitInstance)
		associatedCircuitRaw = &value
	}

	_, cachedZoneSF, _, _ := p.cachedZoneDerivationInputs(instance)
	if zoneOpMode != nil || cachedZoneSF != nil || hasCircuitType {
		operatingMode, preset, allowedModes := deriveZoneModeAndPreset(zoneOpMode, cachedZoneSF, circuitType, hasCircuitType)
		if zoneOpMode != nil {
			incoming.OperatingMode = operatingMode
		}
		if zoneOpMode != nil || cachedZoneSF != nil {
			incoming.Preset = preset
		}
		if hasCircuitType {
			incoming.AllowedModes = allowedModes
		}
	}
	if zoneOpMode != nil {
		incoming.ConfigurationHeatingOperationMode = formatUintToken(*zoneOpMode)
	}

	incoming.TargetTempC = targetPtr
	incoming.HumidityPct = humidity
	incoming.QuickVetoTempC = qvTempPtr
	incoming.QuickVetoDurationH = qvDurPtr
	incoming.QuickVetoEndTime = qvEndTime
	incoming.QuickVetoEndDate = qvEndDate
	incoming.HolidayStartDate = holidayStartDate
	incoming.HolidayEndDate = holidayEndDate
	incoming.HolidaySetpointC = holidaySetpointPtr
	incoming.HolidayStartTime = holidayStartTime
	incoming.HolidayEndTime = holidayEndTime
	incoming.ConfigurationRoomTemperatureZoneMappingRaw = zoneRoomTemperatureZoneMappingRaw
	incoming.ConfigurationAssociatedCircuitRaw = associatedCircuitRaw
	incoming.ConfigurationCircuitTypeRaw = circuitType

	return incoming, liveReadSuccess
}

func (p *vaillantSemanticPoller) cachedZoneDerivationInputs(instance byte) (opMode, specialFunction, circuitType *uint16, hasCircuitType bool) {
	if p == nil {
		return nil, nil, nil, false
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	circuitInstance := instance
	if zone := p.zones[instance]; zone != nil {
		opMode, _ = parseUint16Token(zone.ConfigurationHeatingOperationMode)
		specialFunction, _ = parseUint16Token(zone.StateSpecialFunction)
		if zone.ConfigurationCircuitTypeRaw != nil {
			return opMode, specialFunction, cloneUint16Ptr(zone.ConfigurationCircuitTypeRaw), true
		}
		if zone.ConfigurationAssociatedCircuitRaw != nil && *zone.ConfigurationAssociatedCircuitRaw <= 0xFF {
			circuitInstance = byte(*zone.ConfigurationAssociatedCircuitRaw)
		}
	}
	if circuit := p.circuits[circuitInstance]; circuit != nil && circuit.CircuitTypeRaw != nil {
		return opMode, specialFunction, cloneUint16Ptr(circuit.CircuitTypeRaw), true
	}
	return opMode, specialFunction, nil, false
}

func (p *vaillantSemanticPoller) refreshState(ctx context.Context) {
	controller, zones := p.snapshotZones()
	grabHydrated := false
	if controller == 0 || len(zones) == 0 {
		if p.fallbackDiscoveryAllowed(controller) {
			p.refreshDiscovery(ctx)
			controller, zones = p.snapshotZones()
		} else {
			p.refreshBoilerAddress(ctx)
		}
	}
	if controller == 0 || len(zones) == 0 {
		if p.tryRefreshFromEbusdGrab(ctx) {
			controller, zones = p.snapshotZones()
			grabHydrated = true
		}
	}
	if controller == 0 || len(zones) == 0 {
		dhwSource := p.refreshDHW(ctx)
		p.publishDHW(dhwSource)
		return
	}

	liveReadSuccess := false
	for _, instance := range zones {
		var currentPtr *float64
		if value, ok := p.readB524Float32LE(ctx, localZones.opcode, localZones.group, instance, zoneCurrentTemp); ok {
			current := value
			currentPtr = &current
			liveReadSuccess = true
		}

		zoneSF, zoneSFOK := p.readB524Uint16(ctx, localZones.opcode, localZones.group, instance, zoneSpecialFunction)
		zoneValve, zoneValveOK := p.readB524Uint16(ctx, localZones.opcode, localZones.group, instance, zoneValveStatus)
		if zoneSFOK || zoneValveOK {
			liveReadSuccess = true
		}
		cachedOpMode, cachedZoneSF, circuitType, hasCircuitType := p.cachedZoneDerivationInputs(instance)
		zoneSFForDerive := zoneSF
		if zoneSFForDerive == nil {
			zoneSFForDerive = cachedZoneSF
		}

		allowedModes := deriveZoneAllowedModes(circuitType, hasCircuitType)
		hvacAction := deriveZoneHvacAction(zoneValve, circuitType, hasCircuitType)
		incoming := &vaillantZoneSnapshot{
			HvacAction:          hvacAction,
			CurrentTempC:        currentPtr,
			StateValveStatusRaw: zoneValve,
		}
		if cachedOpMode != nil || hasMeaningfulSpecialFunction(zoneSFForDerive) {
			operatingMode, preset, _ := deriveZoneModeAndPreset(cachedOpMode, zoneSFForDerive, circuitType, hasCircuitType)
			if cachedOpMode != nil {
				incoming.OperatingMode = operatingMode
			}
			incoming.Preset = preset
		}
		if hasCircuitType {
			incoming.AllowedModes = allowedModes
		}
		if zoneSF != nil {
			incoming.StateSpecialFunction = formatUintToken(*zoneSF)
		}

		p.mu.Lock()
		if entry := p.zones[instance]; entry != nil {
			mergeZoneSnapshotFields(entry, incoming, semanticSnapshotSourceLive, zoneFastStateFieldSet)
		}
		p.mu.Unlock()
	}
	if !liveReadSuccess && p.tryRefreshFromEbusdGrab(ctx) {
		grabHydrated = true
	}

	zoneSource := semanticSnapshotSourceCache
	if liveReadSuccess || grabHydrated {
		zoneSource = semanticSnapshotSourceLive
	}

	dhwSource := p.refreshDHW(ctx)
	p.publishZones(zoneSource)
	p.publishDHW(dhwSource)
}
