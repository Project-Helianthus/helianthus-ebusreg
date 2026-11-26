package owner

import (
	"fmt"
	"slices"
	"strings"

	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
)

func (p *vaillantSemanticPoller) commitFunctionModuleAcquisition(
	captured functionModuleEvidenceCapture,
	verdict vaillantFunctionModuleInterpretation,
	incomingSolar *vaillantSolarSnapshot,
	incomingCylinders map[byte]*vaillantCylinderSnapshot,
) vaillantFunctionModuleInterpretation {
	commit := func(registryGeneration uint64, registryCoherent bool) {
		p.mu.Lock()
		defer p.mu.Unlock()
		if !registryCoherent || !captured.registryCoherent ||
			captured.registryGeneration != registryGeneration ||
			!captured.matchesLockedPoller(p) {
			p.functionModuleEvidenceRevision++
			verdict = vaillantFunctionModuleInterpretation{
				Mode:             vaillantFunctionModuleModeGPIOOnly,
				DegradedReason:   vaillantFunctionModuleDegradedIncoherentAcquisition,
				EvidenceRevision: formatFunctionModuleEvidenceRevision(captured.generation, p.functionModuleEvidenceGeneration, p.functionModuleEvidenceRevision),
			}
		}
		verdict = retainFunctionModuleStructuralMode(p.functionModuleInterpretation, verdict)
		if verdict.Mode != "" {
			p.functionModuleMode = verdict.Mode
		}
		p.functionModuleInterpretation = verdict
		p.solar, p.solarCylinders = applyFunctionModuleAcquisition(
			p.solar, p.solarCylinders, incomingSolar, incomingCylinders, verdict,
		)
	}

	withGeneration := p.withFunctionModuleObservationGeneration
	if withGeneration == nil && p.reg != nil {
		withGeneration = p.reg.WithObservationGeneration
	}
	if withGeneration == nil {
		commit(0, true)
		return verdict
	}
	committed := false
	if withGeneration(func(current uint64) {
		// Do not call any DeviceRegistry method from this callback. The
		// registry read lock is deliberately held across the generation
		// comparison and semantic commit.
		commit(current, true)
		committed = true
	}) && committed {
		return verdict
	}
	commit(captured.registryGeneration, false)
	return verdict
}

func (p *vaillantSemanticPoller) updateFunctionModuleConfigurationAcquisitionLocked(acquired bool) {
	failed := !acquired
	if p.functionModuleConfigurationFailed != failed {
		p.functionModuleConfigurationFailed = failed
		p.functionModuleEvidenceGeneration++
	}
}

func uint16PointersEqual(left, right *uint16) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}
	return *left == *right
}

func deriveFunctionModuleInterpretation(
	controllerReachable bool,
	moduleConfig *uint16,
	solarReadable bool,
	cylindersReadable bool,
	hasEvidence bool,
	evidenceStale bool,
	incoherent bool,
	evidenceRevision string,
) vaillantFunctionModuleInterpretation {
	verdict := vaillantFunctionModuleInterpretation{EvidenceRevision: evidenceRevision}
	verdict.Mode = vaillantFunctionModuleModeGPIOOnly
	switch {
	case incoherent:
		verdict.DegradedReason = vaillantFunctionModuleDegradedIncoherentAcquisition
	case !hasEvidence:
		return vaillantFunctionModuleInterpretation{}
	case !controllerReachable:
		verdict.DegradedReason = vaillantFunctionModuleDegradedControllerUnreachable
	case moduleConfig == nil:
		verdict.DegradedReason = vaillantFunctionModuleDegradedConfigurationUnavailable
	case evidenceStale:
		verdict.DegradedReason = vaillantFunctionModuleDegradedEvidenceStale
	case *moduleConfig > 2:
		verdict.DegradedReason = vaillantFunctionModuleDegradedConfigurationNotInterpretable
	case !solarReadable:
		verdict.DegradedReason = vaillantFunctionModuleDegradedSolarAcquisitionFailed
	case !cylindersReadable:
		verdict.DegradedReason = vaillantFunctionModuleDegradedCylinderAcquisitionFailed
	default:
		verdict.Mode = vaillantFunctionModuleModeInterpreted
	}
	return verdict
}

func retainFunctionModuleStructuralMode(previous, candidate vaillantFunctionModuleInterpretation) vaillantFunctionModuleInterpretation {
	if !isFunctionModuleTransientAcquisitionReason(candidate.DegradedReason) {
		return candidate
	}
	if !previous.valid() {
		return vaillantFunctionModuleInterpretation{}
	}
	switch previous.Mode {
	case vaillantFunctionModuleModeInterpreted, vaillantFunctionModuleModeAbsent:
		candidate.Mode = previous.Mode
		return candidate
	case vaillantFunctionModuleModeGPIOOnly:
		previous.EvidenceRevision = candidate.EvidenceRevision
		return previous
	default:
		return vaillantFunctionModuleInterpretation{}
	}
}

func isFunctionModuleTransientAcquisitionReason(reason vaillantFunctionModuleDegradedReason) bool {
	switch reason {
	case vaillantFunctionModuleDegradedControllerUnreachable,
		vaillantFunctionModuleDegradedConfigurationUnavailable,
		vaillantFunctionModuleDegradedSolarAcquisitionFailed,
		vaillantFunctionModuleDegradedCylinderAcquisitionFailed,
		vaillantFunctionModuleDegradedEvidenceStale,
		vaillantFunctionModuleDegradedIncoherentAcquisition:
		return true
	default:
		return false
	}
}

func applyFunctionModuleAcquisition(
	previousSolar *vaillantSolarSnapshot,
	previousCylinders map[byte]*vaillantCylinderSnapshot,
	incomingSolar *vaillantSolarSnapshot,
	incomingCylinders map[byte]*vaillantCylinderSnapshot,
	verdict vaillantFunctionModuleInterpretation,
) (*vaillantSolarSnapshot, map[byte]*vaillantCylinderSnapshot) {
	switch {
	case isFunctionModuleTransientAcquisitionReason(verdict.DegradedReason):
		return cloneSolarSnapshot(previousSolar), cloneCylinderSnapshotsMap(previousCylinders)
	case verdict.Mode == vaillantFunctionModuleModeInterpreted:
		return mergeSolarSnapshotNonDestructive(previousSolar, incomingSolar),
			mergeCylinderSnapshotMapNonDestructive(previousCylinders, incomingCylinders)
	case isFunctionModuleStructuralWithdrawal(verdict):
		return nil, make(map[byte]*vaillantCylinderSnapshot)
	default:
		return cloneSolarSnapshot(previousSolar), cloneCylinderSnapshotsMap(previousCylinders)
	}
}

func isFunctionModuleStructuralWithdrawal(verdict vaillantFunctionModuleInterpretation) bool {
	return (verdict.Mode == vaillantFunctionModuleModeAbsent && verdict.DegradedReason == "") ||
		(verdict.Mode == vaillantFunctionModuleModeGPIOOnly &&
			verdict.DegradedReason == vaillantFunctionModuleDegradedConfigurationNotInterpretable)
}

func (p *vaillantSemanticPoller) nextFunctionModuleEvidenceRevision(startGeneration, endGeneration uint64) string {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.functionModuleEvidenceRevision++
	return formatFunctionModuleEvidenceRevision(startGeneration, endGeneration, p.functionModuleEvidenceRevision)
}

func formatFunctionModuleEvidenceRevision(startGeneration, endGeneration, acquisition uint64) string {
	if startGeneration != endGeneration {
		return fmt.Sprintf("fm5-g%d-%d-a%d", startGeneration, endGeneration, acquisition)
	}
	return fmt.Sprintf("fm5-g%d-a%d", endGeneration, acquisition)
}

func hasFunctionModuleEvidenceFromRadioSnapshots(snapshots []*vaillantRadioDeviceSnapshot) bool {
	for _, snapshot := range snapshots {
		if snapshot == nil {
			continue
		}
		if hasRemoteIdentityEvidence(snapshot.DeviceClassAddress, snapshot.FirmwareVersion, snapshot.HardwareIdentifier) {
			return true
		}
		if snapshot.DeviceClassAddress != nil && *snapshot.DeviceClassAddress == 0x26 {
			return true
		}
		if strings.EqualFold(strings.TrimSpace(snapshot.DeviceModel), "VR71/FM5") {
			return true
		}
	}
	return false
}

func hasFunctionModuleEvidenceFromRadioMap(snapshots map[radioDeviceKey]*vaillantRadioDeviceSnapshot) bool {
	for _, snapshot := range snapshots {
		if hasFunctionModuleEvidenceFromRadioSnapshots([]*vaillantRadioDeviceSnapshot{snapshot}) {
			return true
		}
	}
	return false
}

func (p *vaillantSemanticPoller) hasFunctionModuleRegistryEvidence() bool {
	return p.functionModuleRegistryEvidenceFingerprint() != ""
}

func (p *vaillantSemanticPoller) captureFunctionModuleRegistryEvidence() (string, uint64, bool) {
	if p == nil || p.reg == nil {
		return "", 0, true
	}
	for range 3 {
		var before uint64
		if !p.reg.WithObservationGeneration(func(current uint64) { before = current }) {
			return "", 0, false
		}
		fingerprint := p.functionModuleRegistryEvidenceFingerprint()
		var after uint64
		if !p.reg.WithObservationGeneration(func(current uint64) { after = current }) {
			return fingerprint, before, false
		}
		if before == after {
			return fingerprint, after, true
		}
	}
	return p.functionModuleRegistryEvidenceFingerprint(), 0, false
}

func (p *vaillantSemanticPoller) functionModuleRegistryEvidenceFingerprint() string {
	if p == nil || p.reg == nil {
		return ""
	}
	// P9.3 — race-free DeviceID read via snapshot.
	identities := make([]string, 0, 1)
	p.reg.IterateSnapshots(func(snap registry.DeviceEntrySnapshot) bool {
		deviceID := normalizeDeviceID(snap.DeviceID)
		if strings.HasPrefix(deviceID, "VR71") || strings.HasPrefix(deviceID, "FM5") {
			identities = append(identities, fmt.Sprintf("%02x:%s", snap.PrimaryAddress, deviceID))
		}
		return true
	})
	slices.Sort(identities)
	return strings.Join(identities, "|")
}

func functionModuleInventoryRegistryInfos(system *vaillantSystemSnapshot, functionModuleMode vaillantFunctionModuleMode) []registry.DeviceInfo {
	if system != nil &&
		system.SystemScheme != nil && *system.SystemScheme == 1 &&
		system.ModuleConfigurationFunctionModuleController != nil && *system.ModuleConfigurationFunctionModuleController == 2 &&
		functionModuleMode == vaillantFunctionModuleModeInterpreted {
		return []registry.DeviceInfo{{
			Address:      circuitManagingDeviceFunctionModuleControllerAddress,
			Manufacturer: "Vaillant",
			DeviceID:     circuitManagingDeviceFunctionModuleControllerID,
		}}
	}
	return nil
}
