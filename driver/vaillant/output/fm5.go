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

// PublishedFM5Interpretation maps an acquisition verdict into the published vocabulary.
func PublishedFM5Interpretation(verdict NativeFM5Interpretation) Fm5Interpretation {
	mode, modeOK := publishedFM5Mode(verdict.Mode)
	reason, reasonOK := publishedFM5DegradedReason(verdict.DegradedReason)
	if !modeOK || !reasonOK {
		return Fm5Interpretation{}
	}
	return Fm5Interpretation{
		Mode:             mode,
		DegradedReason:   reason,
		EvidenceRevision: verdict.EvidenceRevision,
	}
}

func publishedFM5Mode(mode NativeFM5Mode) (Fm5SemanticMode, bool) {
	switch mode {
	case "":
		return "", true
	case NativeFM5ModeInterpreted:
		return Fm5SemanticModeInterpreted, true
	case NativeFM5ModeGPIOOnly:
		return Fm5SemanticModeGPIOOnly, true
	case NativeFM5ModeAbsent:
		return Fm5SemanticModeAbsent, true
	default:
		return "", false
	}
}

func publishedFM5DegradedReason(reason NativeFM5DegradedReason) (Fm5SemanticDegradedReason, bool) {
	switch reason {
	case "":
		return "", true
	case NativeFM5DegradedControllerUnreachable:
		return Fm5SemanticDegradedReasonControllerUnreachable, true
	case NativeFM5DegradedConfigurationUnavailable:
		return Fm5SemanticDegradedReasonConfigurationUnavailable, true
	case NativeFM5DegradedConfigurationNotInterpretable:
		return Fm5SemanticDegradedReasonConfigurationNotInterpretable, true
	case NativeFM5DegradedSolarAcquisitionFailed:
		return Fm5SemanticDegradedReasonSolarAcquisitionFailed, true
	case NativeFM5DegradedCylinderAcquisitionFailed:
		return Fm5SemanticDegradedReasonCylinderAcquisitionFailed, true
	case NativeFM5DegradedEvidenceStale:
		return Fm5SemanticDegradedReasonEvidenceStale, true
	case NativeFM5DegradedIncoherentAcquisition:
		return Fm5SemanticDegradedReasonIncoherentAcquisition, true
	default:
		return "", false
	}
}
