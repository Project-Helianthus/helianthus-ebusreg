package b555

import numericNarrowingMath "math"

type NumericNarrowingContext struct {
	NumericNarrowingTimerKind       string `json:"Ab21An21TimerKind"`
	NumericNarrowingFamily          string `json:"Ab21An21Family"`
	NumericNarrowingSoftwareVariant string `json:"Ab21An21SoftwareVariant"`
}

type NumericNarrowingKnowledge struct {
	NumericNarrowingConversion            string `json:"Ab21An21Conversion"`
	NumericNarrowingLocalCarrier          string `json:"Ab21An21LocalCarrier"`
	NumericNarrowingUnits                 string `json:"Ab21An21Units"`
	NumericNarrowingTimerKindEffects      string `json:"Ab21An21TimerKindEffects"`
	NumericNarrowingPersistence           string `json:"Ab21An21Persistence"`
	NumericNarrowingAvailability          string `json:"Ab21An21Availability"`
	NumericNarrowingCountMaterialization  string `json:"Ab21An21CountMaterialization"`
	NumericNarrowingWriteAdmission        string `json:"Ab21An21WriteAdmission"`
	NumericNarrowingPhysicalQualification string `json:"Ab21An21PhysicalQualification"`
	NumericNarrowingPresence              string `json:"Ab21An21Presence"`
	NumericNarrowingFreshness             string `json:"Ab21An21Freshness"`
	NumericNarrowingExecution             string `json:"Ab21An21Execution"`
}

type NumericNarrowingResult struct {
	NumericNarrowingContext               NumericNarrowingContext   `json:"Ab21An21Context"`
	NumericNarrowingPath                  string                    `json:"Ab21An21Path"`
	NumericNarrowingOriginalInteger64     *int64                    `json:"Ab21An21OriginalInteger64"`
	NumericNarrowingOriginalFloat64       *float64                  `json:"Ab21An21OriginalFloat64"`
	NumericNarrowingRoundedFloat64        *float64                  `json:"Ab21An21RoundedFloat64"`
	NumericNarrowingScaledFloat64         *float64                  `json:"Ab21An21ScaledFloat64"`
	NumericNarrowingRawClassification     string                    `json:"Ab21An21RawClassification"`
	NumericNarrowingRoundedClassification string                    `json:"Ab21An21RoundedClassification"`
	NumericNarrowingIntermediateStatus    string                    `json:"Ab21An21IntermediateStatus"`
	NumericNarrowingUnsigned32            *uint32                   `json:"Ab21An21Unsigned32"`
	NumericNarrowingLow16                 *uint16                   `json:"Ab21An21Low16"`
	NumericNarrowingStatus                string                    `json:"Ab21An21Status"`
	NumericNarrowingKnowledge             NumericNarrowingKnowledge `json:"Ab21An21Knowledge"`
}

func numericNarrowingCandidate(numericNarrowingContext NumericNarrowingContext, numericNarrowingPath string) NumericNarrowingResult {
	return NumericNarrowingResult{
		NumericNarrowingContext: numericNarrowingContext,
		NumericNarrowingPath:    numericNarrowingPath,
		NumericNarrowingStatus:  "unresolved",
		NumericNarrowingKnowledge: NumericNarrowingKnowledge{
			NumericNarrowingConversion:            "Proven",
			NumericNarrowingLocalCarrier:          "Hypothesis",
			NumericNarrowingUnits:                 "Unknown",
			NumericNarrowingTimerKindEffects:      "Unknown",
			NumericNarrowingPersistence:           "Unknown",
			NumericNarrowingAvailability:          "Unknown",
			NumericNarrowingCountMaterialization:  "Unknown",
			NumericNarrowingWriteAdmission:        "Unknown",
			NumericNarrowingPhysicalQualification: "Unknown",
			NumericNarrowingPresence:              "Unknown",
			NumericNarrowingFreshness:             "Unknown",
			NumericNarrowingExecution:             "Unknown",
		},
	}
}

func NumericNarrowingNarrowInteger64(numericNarrowingContext NumericNarrowingContext, numericNarrowingValue int64) NumericNarrowingResult {
	numericNarrowingResult := numericNarrowingCandidate(numericNarrowingContext, "ab21an21-integer64")
	numericNarrowingResult.NumericNarrowingOriginalInteger64 = &numericNarrowingValue
	numericNarrowingResult.NumericNarrowingRawClassification = "finite-integer"
	numericNarrowingResult.NumericNarrowingIntermediateStatus = "not-applicable"
	numericNarrowingLow16 := uint16(numericNarrowingValue)
	numericNarrowingResult.NumericNarrowingLow16 = &numericNarrowingLow16
	numericNarrowingResult.NumericNarrowingStatus = "candidate"
	return numericNarrowingResult
}

func NumericNarrowingNarrowFloat64(numericNarrowingContext NumericNarrowingContext, numericNarrowingValue float64) NumericNarrowingResult {
	numericNarrowingResult := numericNarrowingCandidate(numericNarrowingContext, "ab21an21-float64")
	numericNarrowingResult.NumericNarrowingOriginalFloat64 = &numericNarrowingValue
	numericNarrowingResult.NumericNarrowingRawClassification = numericNarrowingClassify(numericNarrowingValue)
	numericNarrowingRounded := float64(float32(numericNarrowingValue))
	numericNarrowingScaled := numericNarrowingRounded * 10
	numericNarrowingResult.NumericNarrowingRoundedFloat64 = &numericNarrowingRounded
	numericNarrowingResult.NumericNarrowingScaledFloat64 = &numericNarrowingScaled
	numericNarrowingResult.NumericNarrowingRoundedClassification = numericNarrowingClassify(numericNarrowingRounded)
	switch {
	case numericNarrowingMath.IsNaN(numericNarrowingValue) || numericNarrowingMath.IsInf(numericNarrowingValue, 0):
		numericNarrowingResult.NumericNarrowingIntermediateStatus = "unsupported-nonfinite-input"
	case numericNarrowingValue < 0:
		numericNarrowingResult.NumericNarrowingIntermediateStatus = "unsupported-negative-input"
	case numericNarrowingMath.IsNaN(numericNarrowingScaled) || numericNarrowingMath.IsInf(numericNarrowingScaled, 0):
		numericNarrowingResult.NumericNarrowingIntermediateStatus = "unsupported-nonfinite-scaled"
	case numericNarrowingScaled < 0 || numericNarrowingScaled > 4294967295:
		numericNarrowingResult.NumericNarrowingIntermediateStatus = "unsupported-scaled-range"
	default:
		numericNarrowingUnsigned32 := uint32(numericNarrowingMath.Trunc(numericNarrowingScaled))
		numericNarrowingLow16 := uint16(numericNarrowingUnsigned32)
		numericNarrowingResult.NumericNarrowingUnsigned32 = &numericNarrowingUnsigned32
		numericNarrowingResult.NumericNarrowingLow16 = &numericNarrowingLow16
		numericNarrowingResult.NumericNarrowingIntermediateStatus = "converted"
		numericNarrowingResult.NumericNarrowingStatus = "candidate"
	}
	return numericNarrowingResult
}

func numericNarrowingClassify(numericNarrowingValue float64) string {
	switch {
	case numericNarrowingMath.IsNaN(numericNarrowingValue):
		return "nan"
	case numericNarrowingMath.IsInf(numericNarrowingValue, 1):
		return "positive-infinity"
	case numericNarrowingMath.IsInf(numericNarrowingValue, -1):
		return "negative-infinity"
	case numericNarrowingValue == 0 && numericNarrowingMath.Signbit(numericNarrowingValue):
		return "negative-zero"
	case numericNarrowingValue == 0:
		return "zero"
	case numericNarrowingValue < 0:
		return "negative-finite"
	default:
		return "positive-finite"
	}
}
