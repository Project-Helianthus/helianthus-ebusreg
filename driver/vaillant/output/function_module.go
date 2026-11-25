package output

import (
	"slices"
)

// SolarStatusFromSnapshot maps native solar values to the output record.
func SolarStatusFromSnapshot(snapshot *NativeSolarSnapshot) *SolarStatus {
	if snapshot == nil {
		return nil
	}
	return &SolarStatus{
		CollectorTemperatureC: CloneFloat64Ptr(snapshot.CollectorTemperatureC),
		ReturnTemperatureC:    CloneFloat64Ptr(snapshot.ReturnTemperatureC),
		PumpActive:            cloneBoilerBoolPtr(snapshot.PumpActive),
		CurrentYield:          CloneFloat64Ptr(snapshot.CurrentYield),
		PumpHours:             decodeUint32Float64(snapshot.PumpHours),
		SolarEnabled:          cloneBoilerBoolPtr(snapshot.SolarEnabled),
		FunctionMode:          cloneBoilerBoolPtr(snapshot.FunctionMode),
	}
}

// CylinderStatusesFromSnapshots maps cylinders with evidence in instance order.
func CylinderStatusesFromSnapshots(snapshots map[byte]*NativeCylinderSnapshot) []CylinderStatus {
	instances := make([]byte, 0, len(snapshots))
	for instance := range snapshots {
		instances = append(instances, instance)
	}
	slices.Sort(instances)
	cylinders := make([]CylinderStatus, 0, len(instances))
	for _, instance := range instances {
		snapshot := snapshots[instance]
		if snapshot == nil || !HasLiveCylinderEvidence(snapshot) {
			continue
		}
		cylinders = append(cylinders, CylinderStatus{
			Index:             int(snapshot.Instance),
			TemperatureC:      CloneFloat64Ptr(snapshot.TemperatureC),
			MaxSetpointC:      CloneFloat64Ptr(snapshot.MaxSetpointC),
			ChargeHysteresisC: CloneFloat64Ptr(snapshot.ChargeHysteresis),
			ChargeOffsetC:     CloneFloat64Ptr(snapshot.ChargeOffset),
		})
	}
	return cylinders
}

// PublishedFunctionModuleInterpretation maps an acquisition verdict into the published vocabulary.
func PublishedFunctionModuleInterpretation(verdict NativeFunctionModuleInterpretation) FunctionModuleInterpretation {
	mode, modeOK := publishedFunctionModuleMode(verdict.Mode)
	reason, reasonOK := publishedFunctionModuleDegradedReason(verdict.DegradedReason)
	if !modeOK || !reasonOK {
		return FunctionModuleInterpretation{}
	}
	return FunctionModuleInterpretation{
		Mode:             mode,
		DegradedReason:   reason,
		EvidenceRevision: verdict.EvidenceRevision,
	}
}

func publishedFunctionModuleMode(mode NativeFunctionModuleMode) (FunctionModuleSemanticMode, bool) {
	switch mode {
	case "":
		return "", true
	case NativeFunctionModuleModeInterpreted:
		return FunctionModuleSemanticModeInterpreted, true
	case NativeFunctionModuleModeGPIOOnly:
		return FunctionModuleSemanticModeGPIOOnly, true
	case NativeFunctionModuleModeAbsent:
		return FunctionModuleSemanticModeAbsent, true
	default:
		return "", false
	}
}

func publishedFunctionModuleDegradedReason(reason NativeFunctionModuleDegradedReason) (FunctionModuleSemanticDegradedReason, bool) {
	switch reason {
	case "":
		return "", true
	case NativeFunctionModuleDegradedControllerUnreachable:
		return FunctionModuleSemanticDegradedReasonControllerUnreachable, true
	case NativeFunctionModuleDegradedConfigurationUnavailable:
		return FunctionModuleSemanticDegradedReasonConfigurationUnavailable, true
	case NativeFunctionModuleDegradedConfigurationNotInterpretable:
		return FunctionModuleSemanticDegradedReasonConfigurationNotInterpretable, true
	case NativeFunctionModuleDegradedSolarAcquisitionFailed:
		return FunctionModuleSemanticDegradedReasonSolarAcquisitionFailed, true
	case NativeFunctionModuleDegradedCylinderAcquisitionFailed:
		return FunctionModuleSemanticDegradedReasonCylinderAcquisitionFailed, true
	case NativeFunctionModuleDegradedEvidenceStale:
		return FunctionModuleSemanticDegradedReasonEvidenceStale, true
	case NativeFunctionModuleDegradedIncoherentAcquisition:
		return FunctionModuleSemanticDegradedReasonIncoherentAcquisition, true
	default:
		return "", false
	}
}
