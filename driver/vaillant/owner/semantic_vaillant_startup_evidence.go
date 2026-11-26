package owner

import (
	"context"
	"strings"

	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
)

func (p *vaillantSemanticPoller) readDhwPseudoCircuitStartupEvidence(ctx context.Context, instance byte) (*vaillantCircuitSnapshot, bool) {
	if p == nil || instance != dhwPseudoCircuitInstance {
		return nil, false
	}
	snapshot := newDhwPseudoCircuitSnapshot(instance)
	flow, flowKnown := p.readDhwPseudoCircuitStartupTemperatureEvidence(ctx, instance, circuitFlowTemp)
	if isDhwPseudoCircuitTemperatureEvidence(flow) {
		snapshot.FlowTemperatureC = flow
	}
	if snapshot.FlowTemperatureC == nil {
		return nil, flowKnown
	}
	return snapshot, true
}

func (p *vaillantSemanticPoller) readDhwPseudoCircuitStartupTemperatureEvidence(ctx context.Context, instance byte, addr uint16) (*float64, bool) {
	raw, ok := p.readB524Startup(ctx, localCircuits.opcode, localCircuits.group, instance, addr)
	if !ok {
		return nil, false
	}
	return decodeB524Float32FromRaw(raw), true
}

func (p *vaillantSemanticPoller) refreshSystemStartup(ctx context.Context) {
	if p == nil {
		return
	}

	p.mu.Lock()
	controller := p.controller
	p.mu.Unlock()
	snapshot := &vaillantSystemSnapshot{Controller: controller}
	readAny := false
	configurationAcquired := false
	if raw, ok := p.readB524Uint16Startup(ctx, localRegulator.opcode, localRegulator.group, regulatorInstance, systemScheme); ok && raw != nil {
		snapshot.SystemScheme = cloneUint16Ptr(raw)
		readAny = true
	}
	if raw, ok := p.readB524Uint16Startup(ctx, localRegulator.opcode, localRegulator.group, regulatorInstance, systemModuleConfigurationFunctionModule); ok && raw != nil {
		snapshot.ModuleConfigurationFunctionModuleController = cloneUint16Ptr(raw)
		readAny = true
		configurationAcquired = true
	}
	if raw, ok := p.readB524Startup(ctx, localRegulator.opcode, localRegulator.group, regulatorInstance, systemWaterPressure); ok {
		if value := decodeB524Float32FromRaw(raw); value != nil {
			snapshot.SystemWaterPressure = value
			readAny = true
		}
	}

	p.mu.Lock()
	source := semanticSnapshotSourceCache
	p.updateFunctionModuleConfigurationAcquisitionLocked(configurationAcquired)
	if readAny {
		if snapshot.SystemFlowTemperature != nil {
			snapshot.SystemFlowTemperatureLiveAt = p.now()
		}
		p.updateSystemSnapshotLocked(mergeSystemSnapshotNonDestructive(p.system, snapshot))
		source = semanticSnapshotSourceLive
	}
	hasSnapshot := p.system != nil
	p.mu.Unlock()
	if hasSnapshot {
		p.publishSystem(source)
	}
}

func (p *vaillantSemanticPoller) refreshFunctionModuleConfigurationForLateIdentity(ctx context.Context) {
	if p == nil {
		return
	}

	p.mu.Lock()
	controller := p.controller
	p.mu.Unlock()
	snapshot := &vaillantSystemSnapshot{Controller: controller}
	configurationAcquired := false
	if raw, ok := p.readB524Uint16Startup(ctx, localRegulator.opcode, localRegulator.group, regulatorInstance, systemModuleConfigurationFunctionModule); ok && raw != nil {
		snapshot.ModuleConfigurationFunctionModuleController = cloneUint16Ptr(raw)
		configurationAcquired = true
	}

	p.mu.Lock()
	source := semanticSnapshotSourceCache
	p.updateFunctionModuleConfigurationAcquisitionLocked(configurationAcquired)
	if configurationAcquired {
		p.updateSystemSnapshotLocked(mergeSystemSnapshotNonDestructive(p.system, snapshot))
		source = semanticSnapshotSourceLive
	}
	hasSnapshot := p.system != nil
	p.mu.Unlock()
	if hasSnapshot {
		p.publishSystem(source)
	}
}

func (p *vaillantSemanticPoller) refreshRadioDevicesStartup(ctx context.Context) {
	if p == nil {
		return
	}

	discovered := p.registryRadioDeviceSeeds()
	fullScanGroups := startupRadioFullScanGroups(discovered)
	functionModulePositiveObserved := false
	verified := make(map[radioDeviceKey]bool)
	readAny := false
	probeSlot := func(grp b524GroupDef, instance byte) {
		connectedRaw, ok := p.readB524Startup(ctx, grp.opcode, grp.group, instance, deviceSlotConnected)
		if !ok || len(connectedRaw) == 0 {
			return
		}
		readAny = true
		connected := connectedRaw[0] == 1
		classAddress := p.readB524Uint8Startup(ctx, grp.opcode, grp.group, instance, deviceSlotClassAddress)
		key := radioDeviceKey{Group: grp.group, Instance: instance}
		verified[key] = true
		include, slotMode := startupRadioDeviceInclude(grp.group, connected, classAddress)
		if !include {
			delete(discovered, key)
			return
		}
		discovered[key] = &vaillantRadioDeviceSnapshot{
			Group:              grp.group,
			Instance:           instance,
			SlotMode:           slotMode,
			DeviceConnected:    &connected,
			DeviceClassAddress: cloneUint8Ptr(classAddress),
			DeviceModel:        decodeRadioDeviceModel(classAddress),
		}
	}

	for _, grp := range remoteDeviceGroups {
		if grp.group == remoteFunctionalModules.group {
			continue
		}
		for instance := byte(0x00); instance <= semanticStartupSlotFastMaxInstance; instance++ {
			probeSlot(grp, instance)
		}
	}
	functionModuleLastInstance := semanticStartupSlotFastMaxInstance
	if fullScanGroups[remoteFunctionalModules.group] {
		functionModuleLastInstance = semanticStartupSlotFullMaxInstance
	}
	functionModuleScan := p.scanFunctionalModuleIdentityBudgetAwareWith(ctx, 0x00, functionModuleLastInstance, p.readB524Startup)
	if functionModuleScan.observedAny {
		readAny = true
	}
	for instance := byte(0x00); instance <= functionModuleLastInstance; instance++ {
		key := radioDeviceKey{Group: remoteFunctionalModules.group, Instance: instance}
		verified[key] = true
		if snapshot := functionModuleScan.snapshots[instance]; snapshot != nil {
			discovered[key] = snapshot
		} else {
			delete(discovered, key)
		}
	}
	functionModulePositiveObserved = functionModuleScan.positiveSnapshot != nil
	functionModuleNamespaceComplete := fullScanGroups[remoteFunctionalModules.group] && functionModuleScan.namespaceComplete
	if !functionModulePositiveObserved {
		for _, grp := range remoteDeviceGroups {
			if grp.group == remoteFunctionalModules.group || !fullScanGroups[grp.group] {
				continue
			}
			for instance := semanticStartupSlotFastMaxInstance + 1; instance <= semanticStartupSlotFullMaxInstance; instance++ {
				probeSlot(grp, instance)
			}
		}
	}
	if readAny {
		for key, snapshot := range discovered {
			if snapshot != nil && snapshot.SlotMode == "registry" && !verified[key] {
				if functionModulePositiveObserved && key.Group != remoteFunctionalModules.group {
					continue
				}
				delete(discovered, key)
			}
		}
	}

	liveFunctionModuleEvidence := hasFunctionModuleEvidenceFromRadioMap(discovered)
	p.mu.Lock()
	p.startupRadioDevicesProbed = readAny
	p.functionModuleIdentityScanComplete = functionModuleNamespaceComplete
	p.functionModuleIdentityIncoherent = !functionModuleNamespaceComplete && !liveFunctionModuleEvidence
	if functionModuleNamespaceComplete {
		p.functionModuleRegistryEvidenceIgnored = !liveFunctionModuleEvidence
	} else if liveFunctionModuleEvidence {
		p.functionModuleRegistryEvidenceIgnored = false
	}
	if len(discovered) == 0 {
		if readAny {
			p.radioDevices = make(map[radioDeviceKey]*vaillantRadioDeviceSnapshot)
			p.functionModuleEvidenceGeneration++
			if functionModuleNamespaceComplete {
				p.functionModuleIdentityObservedAt = p.now()
			}
		}
		p.mu.Unlock()
		if readAny {
			p.publishRadioDevices(semanticSnapshotSourceLive)
		}
		return
	}
	if p.radioDevices == nil || readAny {
		p.radioDevices = make(map[radioDeviceKey]*vaillantRadioDeviceSnapshot)
	}
	for key, snapshot := range discovered {
		p.radioDevices[key] = snapshot
	}
	if readAny {
		p.functionModuleEvidenceGeneration++
		if functionModuleNamespaceComplete || liveFunctionModuleEvidence {
			p.functionModuleIdentityObservedAt = p.now()
		}
	}
	p.mu.Unlock()
	source := semanticSnapshotSourceCache
	if readAny {
		source = semanticSnapshotSourceLive
	}
	p.publishRadioDevices(source)
}

func startupRadioFullScanGroups(discovered map[radioDeviceKey]*vaillantRadioDeviceSnapshot) map[byte]bool {
	seeded := make(map[byte]bool)
	highSeeded := make(map[byte]bool)
	for key := range discovered {
		seeded[key.Group] = true
		if key.Instance > semanticStartupSlotFastMaxInstance {
			highSeeded[key.Group] = true
		}
	}

	out := make(map[byte]bool)
	for _, grp := range remoteDeviceGroups {
		if !seeded[grp.group] || highSeeded[grp.group] {
			out[grp.group] = true
		}
	}
	return out
}

func startupRadioDeviceInclude(group byte, connected bool, classAddress *uint8) (bool, string) {
	switch group {
	case remoteRegulators.group, remoteThermostats.group:
		return connected, "startup"
	case remoteFunctionalModules.group:
		if connected {
			return true, "startup"
		}
		if classAddress != nil && *classAddress == circuitManagingDeviceFunctionModuleControllerAddress {
			return true, "inventory"
		}
		return false, "startup"
	default:
		return connected, "startup"
	}
}

func (p *vaillantSemanticPoller) registryRadioDeviceSeeds() map[radioDeviceKey]*vaillantRadioDeviceSnapshot {
	out := make(map[radioDeviceKey]*vaillantRadioDeviceSnapshot)
	if p == nil || p.reg == nil {
		return out
	}
	nextInstance := map[byte]byte{}
	p.reg.IterateSnapshots(func(snap registry.DeviceEntrySnapshot) bool {
		addr := snapshotTargetAddressForRouting(snap)
		deviceID := normalizeDeviceID(snap.DeviceID)
		var group byte
		switch {
		case strings.HasPrefix(deviceID, "BASV") || addr == 0x15:
			group = remoteRegulators.group
		case strings.HasPrefix(deviceID, "VR71") || strings.HasPrefix(deviceID, "FM5") || addr == circuitManagingDeviceFunctionModuleControllerAddress:
			group = remoteFunctionalModules.group
		case strings.HasPrefix(deviceID, "VR92"):
			group = remoteThermostats.group
		default:
			return true
		}
		instance := nextInstance[group]
		nextInstance[group] = instance + 1
		classAddress := uint8(addr)
		connected := true
		out[radioDeviceKey{Group: group, Instance: instance}] = &vaillantRadioDeviceSnapshot{
			Group:              group,
			Instance:           instance,
			SlotMode:           "registry",
			DeviceConnected:    &connected,
			DeviceClassAddress: &classAddress,
			DeviceModel:        decodeRadioDeviceModel(&classAddress),
		}
		return true
	})
	return out
}

func (p *vaillantSemanticPoller) readB524Uint8Startup(ctx context.Context, opcode, group, instance byte, addr uint16) *uint8 {
	raw, ok := p.readB524Startup(ctx, opcode, group, instance, addr)
	if !ok || len(raw) == 0 {
		return nil
	}
	value := raw[0]
	return &value
}

func (p *vaillantSemanticPoller) refreshFunctionModuleSemanticStartup(ctx context.Context) {
	if !p.publishing() {
		return
	}

	evidence := p.captureFunctionModuleEvidence()
	moduleConfig := evidence.moduleConfig
	if evidence.configurationFailed {
		moduleConfig = nil
	}
	functionModuleGateSatisfied := moduleConfig != nil && *moduleConfig <= 2
	observedAt := p.now()
	evidenceStale := evidence.staleAt(observedAt, p.functionModuleEvidenceTTL)
	var incomingSolar *vaillantSolarSnapshot
	incomingCylinders := make(map[byte]*vaillantCylinderSnapshot)
	solarReadable := false
	cylindersReadable := false
	if evidence.controller != 0 && functionModuleGateSatisfied && !evidenceStale && !evidence.identityIncoherent {
		incomingSolar, solarReadable = p.readSolarSnapshotStartup(ctx)
		incomingCylinders, cylindersReadable = p.readCylinderSnapshotsStartup(ctx)
	}
	currentEvidence := p.captureFunctionModuleEvidence()
	incoherent := evidence.identityIncoherent || currentEvidence.identityIncoherent || !evidence.sameGeneration(currentEvidence)
	evidenceRevision := p.nextFunctionModuleEvidenceRevision(evidence.generation, currentEvidence.generation)
	negativeIdentityObserved := evidence.hasNegativeObservation() && currentEvidence.hasNegativeObservation()
	freshNegativeIdentityObserved := evidence.hasFreshNegativeObservation(observedAt, p.functionModuleEvidenceTTL) &&
		currentEvidence.hasFreshNegativeObservation(observedAt, p.functionModuleEvidenceTTL)
	var verdict vaillantFunctionModuleInterpretation
	if freshNegativeIdentityObserved && !incoherent && evidence.controller != 0 && moduleConfig != nil {
		verdict = vaillantFunctionModuleInterpretation{
			Mode:             vaillantFunctionModuleModeAbsent,
			EvidenceRevision: evidenceRevision,
		}
	} else {
		verdict = deriveFunctionModuleInterpretation(
			evidence.controller != 0,
			moduleConfig,
			solarReadable,
			cylindersReadable,
			evidence.hasEvidence() || currentEvidence.hasEvidence() || negativeIdentityObserved,
			evidenceStale,
			incoherent,
			evidenceRevision,
		)
	}
	verdict = p.commitFunctionModuleAcquisition(currentEvidence, verdict, incomingSolar, incomingCylinders)
	for _, info := range functionModuleInventoryRegistryInfos(evidence.systemSnapshot, verdict.Mode) {
		p.reg.Register(preserveExistingRegistryMetadata(p.reg, info))
	}
	p.publishFunctionModuleSemantic(semanticSnapshotSourceLive)
}

func (p *vaillantSemanticPoller) readSolarSnapshotStartup(ctx context.Context) (*vaillantSolarSnapshot, bool) {
	incoming := &vaillantSolarSnapshot{}
	readAny := false
	if raw, ok := p.readB524Startup(ctx, localSolar.opcode, localSolar.group, solarInstance, solarEnabled); ok {
		incoming.SolarEnabled = decodeB524BoolFromRaw(raw)
		readAny = true
	}
	if raw, ok := p.readB524Startup(ctx, localSolar.opcode, localSolar.group, solarInstance, solarCollectorTemp); ok {
		incoming.CollectorTemperatureC = decodeB524Float32FromRaw(raw)
		readAny = true
	}
	if raw, ok := p.readB524Startup(ctx, localSolar.opcode, localSolar.group, solarInstance, solarPumpActive); ok {
		incoming.PumpActive = decodeB524BoolFromRaw(raw)
		readAny = true
	}
	if !readAny {
		return nil, false
	}
	return incoming, true
}

func (p *vaillantSemanticPoller) readCylinderSnapshotsStartup(ctx context.Context) (map[byte]*vaillantCylinderSnapshot, bool) {
	out := make(map[byte]*vaillantCylinderSnapshot, 2)
	for instance := byte(0x00); instance <= 0x01; instance++ {
		raw, ok := p.readB524Startup(ctx, localCylinders.opcode, localCylinders.group, instance, cylinderTemperature)
		if !ok {
			continue
		}
		out[instance] = &vaillantCylinderSnapshot{
			Instance:     instance,
			TemperatureC: decodeB524Float32FromRaw(raw),
		}
	}
	if len(out) == 0 {
		return nil, false
	}
	return out, true
}

func (p *vaillantSemanticPoller) refreshBoilerStatusStartup(ctx context.Context) {
	if p == nil {
		return
	}

	p.mu.Lock()
	boilerAddress := p.boilerAddress
	p.mu.Unlock()
	if boilerAddress == 0 {
		boilerAddress = p.findBoilerAddress()
	}

	snapshot := &vaillantBoilerSnapshot{}
	updated := false
	if raw, ok := p.readB524Startup(ctx, localDHW.opcode, localDHW.group, dhwInstance, dhwCurrentTemp); ok {
		if value := decodeB524Float32FromRaw(raw); value != nil {
			snapshot.DhwTemperatureC = value
			updated = true
		}
	}

	p.mu.Lock()
	if boilerAddress != 0 {
		p.boilerAddress = boilerAddress
	}
	source := semanticSnapshotSourceCache
	if updated {
		p.boiler = mergeBoilerSnapshotNonDestructive(p.boiler, snapshot)
		source = semanticSnapshotSourceLive
	}
	hasSnapshot := p.boiler != nil
	p.mu.Unlock()
	if hasSnapshot {
		p.publishBoilerStatus(source)
	}
}
