package owner

import (
	"context"
	"strings"
	"time"
)

func (p *vaillantSemanticPoller) readFunctionalModuleIdentitySnapshot(ctx context.Context, instance byte) (*vaillantRadioDeviceSnapshot, bool, bool) {
	return p.readFunctionalModuleIdentitySnapshotWith(ctx, instance, p.readB524Value)
}

type functionalModuleIdentityReadFunc func(context.Context, byte, byte, byte, uint16) ([]byte, bool)

type functionalModuleIdentityProbe struct {
	instance     byte
	connectedRaw []byte
	connectedOK  bool
	classRaw     []byte
	classOK      bool
	completed    bool
}

type functionalModuleIdentityScanResult struct {
	snapshots         map[byte]*vaillantRadioDeviceSnapshot
	observedAny       bool
	namespaceComplete bool
	positiveInstance  byte
	positiveSnapshot  *vaillantRadioDeviceSnapshot
}

func (p *vaillantSemanticPoller) scanFunctionalModuleIdentityBudgetAwareWith(
	ctx context.Context,
	firstInstance byte,
	lastInstance byte,
	read functionalModuleIdentityReadFunc,
) functionalModuleIdentityScanResult {
	result := functionalModuleIdentityScanResult{
		snapshots:         make(map[byte]*vaillantRadioDeviceSnapshot),
		namespaceComplete: true,
	}
	if lastInstance < firstInstance {
		result.namespaceComplete = false
		return result
	}
	probes := make([]functionalModuleIdentityProbe, 0, int(lastInstance-firstInstance)+1)
	preflight := func(first byte, last byte) {
		for instance := first; instance <= last; instance++ {
			connectedRaw, connectedOK := read(ctx, vaillantB524OpcodeRead, remoteFunctionalModules.group, instance, device_slot_connected)
			classRaw, classOK := read(ctx, vaillantB524OpcodeRead, remoteFunctionalModules.group, instance, device_slot_class_address)
			if connectedOK && len(connectedRaw) > 0 {
				result.observedAny = true
			}
			probes = append(probes, functionalModuleIdentityProbe{
				instance:     instance,
				connectedRaw: connectedRaw,
				connectedOK:  connectedOK,
				classRaw:     classRaw,
				classOK:      classOK,
			})
		}
	}

	completeProbe := func(probe *functionalModuleIdentityProbe) functionalModuleIdentityTuple {
		firmwareRaw, firmwareOK := read(ctx, vaillantB524OpcodeRead, remoteFunctionalModules.group, probe.instance, device_slot_firmware)
		hardwareRaw, hardwareOK := read(ctx, vaillantB524OpcodeRead, remoteFunctionalModules.group, probe.instance, device_slot_hardware_identifier)
		probe.completed = true
		tuple := decodeFunctionalModuleIdentityTuple(
			probe.connectedRaw, probe.connectedOK,
			probe.classRaw, probe.classOK,
			firmwareRaw, firmwareOK,
			hardwareRaw, hardwareOK,
		)
		if !tuple.complete {
			result.namespaceComplete = false
		}
		if snapshot := functionalModuleIdentitySnapshot(probe.instance, tuple); snapshot != nil {
			result.snapshots[probe.instance] = snapshot
		}
		return tuple
	}

	completeCandidates := func(firstProbe int) bool {
		for i := firstProbe; i < len(probes); i++ {
			probe := &probes[i]
			if !probe.classOK || len(probe.classRaw) == 0 || probe.classRaw[0] != circuitManagingDeviceFunctionModuleControllerAddress {
				continue
			}
			tuple := completeProbe(probe)
			if completeFunctionModuleControllerIdentityTuple(tuple) {
				result.namespaceComplete = false
				result.positiveInstance = probe.instance
				result.positiveSnapshot = result.snapshots[probe.instance]
				return true
			}
		}
		return false
	}

	fastLast := lastInstance
	if fastLast > semanticStartupSlotFastMaxInstance {
		fastLast = semanticStartupSlotFastMaxInstance
	}
	preflight(firstInstance, fastLast)
	if completeCandidates(0) {
		return result
	}
	if fastLast < lastInstance {
		firstTailProbe := len(probes)
		preflight(fastLast+1, lastInstance)
		if completeCandidates(firstTailProbe) {
			return result
		}
	}

	for i := range probes {
		if !probes[i].completed {
			completeProbe(&probes[i])
		}
	}
	return result
}

func (p *vaillantSemanticPoller) readFunctionalModuleIdentitySnapshotWith(
	ctx context.Context,
	instance byte,
	read functionalModuleIdentityReadFunc,
) (*vaillantRadioDeviceSnapshot, bool, bool) {
	connectedRaw, connectedOK := read(ctx, vaillantB524OpcodeRead, remoteFunctionalModules.group, instance, device_slot_connected)
	if !connectedOK || len(connectedRaw) == 0 {
		return nil, false, false
	}
	classRaw, classOK := read(ctx, vaillantB524OpcodeRead, remoteFunctionalModules.group, instance, device_slot_class_address)
	firmwareRaw, firmwareOK := read(ctx, vaillantB524OpcodeRead, remoteFunctionalModules.group, instance, device_slot_firmware)
	hardwareRaw, hardwareOK := read(ctx, vaillantB524OpcodeRead, remoteFunctionalModules.group, instance, device_slot_hardware_identifier)
	tuple := decodeFunctionalModuleIdentityTuple(
		connectedRaw, connectedOK,
		classRaw, classOK,
		firmwareRaw, firmwareOK,
		hardwareRaw, hardwareOK,
	)
	return functionalModuleIdentitySnapshot(instance, tuple), true, tuple.complete
}

func functionalModuleIdentitySnapshot(instance byte, tuple functionalModuleIdentityTuple) *vaillantRadioDeviceSnapshot {
	if tuple.classAddress != nil && *tuple.classAddress != circuitManagingDeviceFunctionModuleControllerAddress {
		return nil
	}
	if !tuple.connected && !hasRemoteIdentityEvidence(tuple.classAddress, tuple.firmware, tuple.hardware) {
		return nil
	}
	slotMode := "active"
	if !tuple.connected {
		slotMode = "inventory"
	}
	return &vaillantRadioDeviceSnapshot{
		Group:              remoteFunctionalModules.group,
		Instance:           instance,
		SlotMode:           slotMode,
		DeviceConnected:    cloneBoilerBoolPtr(&tuple.connected),
		DeviceClassAddress: cloneUint8Ptr(tuple.classAddress),
		DeviceModel:        decodeRadioDeviceModel(tuple.classAddress),
		FirmwareVersion:    cloneStringPtr(tuple.firmware),
		HardwareIdentifier: cloneUint16Ptr(tuple.hardware),
	}
}

func (p *vaillantSemanticPoller) refreshFunctionModuleAfterLateIdentity(ctx context.Context) {
	if p == nil {
		return
	}
	scan := p.scanFunctionalModuleIdentityBudgetAwareWith(ctx, 0x00, semanticStartupSlotFullMaxInstance, p.readB524Startup)
	if scan.positiveSnapshot == nil {
		return
	}
	instance := scan.positiveInstance
	snapshot := scan.positiveSnapshot

	deviceKey := deviceSlotKey{Group: remoteFunctionalModules.group, Instance: instance}
	radioKey := radioDeviceKey{Group: remoteFunctionalModules.group, Instance: instance}
	p.mu.Lock()
	if p.deviceSlotCache == nil {
		p.deviceSlotCache = make(map[deviceSlotKey]bool)
	}
	if p.radioDevices == nil {
		p.radioDevices = make(map[radioDeviceKey]*vaillantRadioDeviceSnapshot)
	}
	p.deviceSlotCache[deviceKey] = true
	p.radioDevices[radioKey] = snapshot
	p.functionModuleEvidenceGeneration++
	p.functionModuleIdentityScanComplete = false
	p.functionModuleIdentityIncoherent = false
	p.functionModuleRegistryEvidenceIgnored = false
	p.functionModuleIdentityObservedAt = p.now()
	p.mu.Unlock()

	p.publishRadioDevices(semanticSnapshotSourceLive)
	p.refreshFunctionModuleSemantic(ctx)
}

type functionalModuleIdentityTuple struct {
	connected    bool
	classAddress *uint8
	firmware     *string
	hardware     *uint16
	complete     bool
}

func completeFunctionModuleControllerIdentityTuple(tuple functionalModuleIdentityTuple) bool {
	return tuple.complete && tuple.classAddress != nil && *tuple.classAddress == circuitManagingDeviceFunctionModuleControllerAddress
}

func decodeFunctionalModuleIdentityTuple(
	connectedRaw []byte,
	connectedOK bool,
	classRaw []byte,
	classOK bool,
	firmwareRaw []byte,
	firmwareOK bool,
	hardwareRaw []byte,
	hardwareOK bool,
) functionalModuleIdentityTuple {
	connected, connectedValid := decodeExactB524Bool(connectedRaw)
	tuple := functionalModuleIdentityTuple{
		connected: connected,
		complete: connectedOK && connectedValid &&
			classOK && len(classRaw) >= 1 &&
			firmwareOK && len(firmwareRaw) >= 3 &&
			hardwareOK && len(hardwareRaw) >= 2,
	}
	if classOK && len(classRaw) >= 1 && classRaw[0] != 0xFF {
		value := classRaw[0]
		tuple.classAddress = &value
	}
	if firmwareOK && len(firmwareRaw) >= 3 {
		tuple.firmware = decodeB524FirmwareVersion(firmwareRaw)
	}
	if hardwareOK && len(hardwareRaw) >= 2 {
		if value, ok := decodeB524Uint16(hardwareRaw); ok && value != 0xFFFF {
			tuple.hardware = &value
		}
	}
	if tuple.connected && !hasRemoteIdentityEvidence(tuple.classAddress, tuple.firmware, tuple.hardware) {
		tuple.complete = false
	}
	return tuple
}

func (p *vaillantSemanticPoller) refreshFunctionModuleSemantic(ctx context.Context) {
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
		incomingSolar, solarReadable = p.readSolarSnapshot(ctx)
		incomingCylinders, cylindersReadable = p.readCylinderSnapshots(ctx)
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

// vaillantFM5Mode is the acquisition verdict on the function module: its
// solar and cylinder planes are interpreted, it is known only at GPIO level,
// or it is absent. The zero value means that no verdict was reached yet. The
// publication owner maps it to the published mode; the values are spelled
// like the published ones so that diagnostics read the same.
type vaillantFunctionModuleMode string

const (
	vaillantFunctionModuleModeInterpreted vaillantFunctionModuleMode = "INTERPRETED"
	vaillantFunctionModuleModeGPIOOnly    vaillantFunctionModuleMode = "GPIO_ONLY"
	vaillantFunctionModuleModeAbsent      vaillantFunctionModuleMode = "ABSENT"
)

// vaillantFM5DegradedReason says why an acquisition did not reach, or did
// not keep, an interpreted verdict. Every reason except a configuration that
// is not interpretable is transient. The zero value means not degraded.
type vaillantFunctionModuleDegradedReason string

const (
	vaillantFunctionModuleDegradedControllerUnreachable         vaillantFunctionModuleDegradedReason = "CONTROLLER_UNREACHABLE"
	vaillantFunctionModuleDegradedConfigurationUnavailable      vaillantFunctionModuleDegradedReason = "CONFIGURATION_UNAVAILABLE"
	vaillantFunctionModuleDegradedConfigurationNotInterpretable vaillantFunctionModuleDegradedReason = "CONFIGURATION_NOT_INTERPRETABLE"
	vaillantFunctionModuleDegradedSolarAcquisitionFailed        vaillantFunctionModuleDegradedReason = "SOLAR_ACQUISITION_FAILED"
	vaillantFunctionModuleDegradedCylinderAcquisitionFailed     vaillantFunctionModuleDegradedReason = "CYLINDER_ACQUISITION_FAILED"
	vaillantFunctionModuleDegradedEvidenceStale                 vaillantFunctionModuleDegradedReason = "EVIDENCE_STALE"
	vaillantFunctionModuleDegradedIncoherentAcquisition         vaillantFunctionModuleDegradedReason = "INCOHERENT_ACQUISITION"
)

// vaillantFM5Interpretation is one acquisition verdict: the mode, the reason
// it is degraded, and the revision of the identity evidence it was derived
// from.
type vaillantFunctionModuleInterpretation struct {
	Mode             vaillantFunctionModuleMode
	DegradedReason   vaillantFunctionModuleDegradedReason
	EvidenceRevision string
}

// valid reports whether the verdict is complete: it names its evidence
// revision, a retained interpreted or absent mode carries at most a transient
// reason, and a GPIO-only mode carries the not-interpretable configuration
// reason. The publication accepts exactly the valid verdicts.
func (verdict vaillantFunctionModuleInterpretation) valid() bool {
	if strings.TrimSpace(verdict.EvidenceRevision) == "" {
		return false
	}
	switch verdict.Mode {
	case vaillantFunctionModuleModeInterpreted, vaillantFunctionModuleModeAbsent:
		return verdict.DegradedReason == "" || isFunctionModuleTransientAcquisitionReason(verdict.DegradedReason)
	case vaillantFunctionModuleModeGPIOOnly:
		return verdict.DegradedReason == vaillantFunctionModuleDegradedConfigurationNotInterpretable
	default:
		return false
	}
}

func DeriveFunctionModuleSemanticMode(controllerReachable, functionModuleGateSatisfied, solarReadable, cylindersReadable, hasEvidence bool) vaillantFunctionModuleMode {
	if controllerReachable && functionModuleGateSatisfied && solarReadable && cylindersReadable {
		return vaillantFunctionModuleModeInterpreted
	}
	if hasEvidence {
		return vaillantFunctionModuleModeGPIOOnly
	}
	return vaillantFunctionModuleModeAbsent
}

type functionModuleEvidenceCapture struct {
	controller                  byte
	moduleConfig                *uint16
	systemSnapshot              *vaillantSystemSnapshot
	radioSnapshots              []*vaillantRadioDeviceSnapshot
	registryEvidence            string
	registryGeneration          uint64
	registryCoherent            bool
	registryEvidenceIgnored     bool
	identityAcquisitionComplete bool
	identityIncoherent          bool
	configurationFailed         bool
	identityObservedAt          time.Time
	generation                  uint64
}

func (p *vaillantSemanticPoller) captureFunctionModuleEvidence() functionModuleEvidenceCapture {
	if p == nil {
		return functionModuleEvidenceCapture{}
	}
	p.mu.Lock()
	capture := functionModuleEvidenceCapture{
		controller:                  p.controller,
		systemSnapshot:              cloneSystemSnapshot(p.system),
		registryEvidenceIgnored:     p.functionModuleRegistryEvidenceIgnored,
		identityAcquisitionComplete: p.functionModuleIdentityScanComplete,
		identityIncoherent:          p.functionModuleIdentityIncoherent,
		configurationFailed:         p.functionModuleConfigurationFailed,
		identityObservedAt:          p.functionModuleIdentityObservedAt,
		generation:                  p.functionModuleEvidenceGeneration,
		radioSnapshots:              make([]*vaillantRadioDeviceSnapshot, 0, len(p.radioDevices)),
	}
	if p.system != nil {
		capture.moduleConfig = cloneUint16Ptr(p.system.ModuleConfigurationFunctionModuleController)
	}
	for _, snapshot := range p.radioDevices {
		if snapshot != nil {
			capture.radioSnapshots = append(capture.radioSnapshots, cloneRadioSnapshot(snapshot))
		}
	}
	p.mu.Unlock()

	capture.registryEvidence, capture.registryGeneration, capture.registryCoherent = p.captureFunctionModuleRegistryEvidence()
	if capture.hasEvidence() && capture.identityObservedAt.IsZero() {
		observedAt := p.now()
		p.mu.Lock()
		if p.functionModuleIdentityObservedAt.IsZero() {
			p.functionModuleIdentityObservedAt = observedAt
		}
		capture.identityObservedAt = p.functionModuleIdentityObservedAt
		capture.generation = p.functionModuleEvidenceGeneration
		p.mu.Unlock()
	}
	return capture
}

func (capture functionModuleEvidenceCapture) hasEvidence() bool {
	if hasFunctionModuleEvidenceFromRadioSnapshots(capture.radioSnapshots) {
		return true
	}
	return !capture.registryEvidenceIgnored && capture.registryEvidence != ""
}

func (capture functionModuleEvidenceCapture) staleAt(now time.Time, ttl time.Duration) bool {
	if (!capture.hasEvidence() && !capture.identityAcquisitionComplete) || capture.identityObservedAt.IsZero() || ttl <= 0 {
		return false
	}
	return now.Sub(capture.identityObservedAt) > ttl
}

func (capture functionModuleEvidenceCapture) hasFreshNegativeObservation(now time.Time, ttl time.Duration) bool {
	return capture.hasNegativeObservation() && !capture.staleAt(now, ttl)
}

func (capture functionModuleEvidenceCapture) hasNegativeObservation() bool {
	return !capture.hasEvidence() && capture.identityAcquisitionComplete && !capture.identityObservedAt.IsZero()
}

func (capture functionModuleEvidenceCapture) sameGeneration(other functionModuleEvidenceCapture) bool {
	return capture.registryCoherent && other.registryCoherent &&
		capture.generation == other.generation &&
		capture.controller == other.controller &&
		uint16PointersEqual(capture.moduleConfig, other.moduleConfig) &&
		capture.registryEvidenceIgnored == other.registryEvidenceIgnored &&
		capture.identityAcquisitionComplete == other.identityAcquisitionComplete &&
		capture.identityIncoherent == other.identityIncoherent &&
		capture.configurationFailed == other.configurationFailed &&
		capture.identityObservedAt.Equal(other.identityObservedAt) &&
		capture.registryGeneration == other.registryGeneration &&
		capture.registryEvidence == other.registryEvidence
}

func (capture functionModuleEvidenceCapture) matchesLockedPoller(p *vaillantSemanticPoller) bool {
	if p == nil {
		return false
	}
	var moduleConfig *uint16
	if p.system != nil {
		moduleConfig = p.system.ModuleConfigurationFunctionModuleController
	}
	return capture.generation == p.functionModuleEvidenceGeneration &&
		capture.controller == p.controller &&
		uint16PointersEqual(capture.moduleConfig, moduleConfig) &&
		capture.registryEvidenceIgnored == p.functionModuleRegistryEvidenceIgnored &&
		capture.identityAcquisitionComplete == p.functionModuleIdentityScanComplete &&
		capture.identityIncoherent == p.functionModuleIdentityIncoherent &&
		capture.configurationFailed == p.functionModuleConfigurationFailed &&
		capture.identityObservedAt.Equal(p.functionModuleIdentityObservedAt)
}
