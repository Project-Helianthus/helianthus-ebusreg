package b555

import numericNarrowingMath "math"

type NumericNarrowingContext struct {
	NumericNarrowingTimerKind       string `json:"timer_kind"`
	NumericNarrowingFamily          string `json:"family"`
	NumericNarrowingSoftwareVariant string `json:"software_variant"`
}

type NumericNarrowingKnowledge struct {
	NumericNarrowingConversion            string `json:"conversion"`
	NumericNarrowingLocalCarrier          string `json:"local_carrier"`
	NumericNarrowingUnits                 string `json:"units"`
	NumericNarrowingTimerKindEffects      string `json:"timer_kind_effects"`
	NumericNarrowingPersistence           string `json:"persistence"`
	NumericNarrowingAvailability          string `json:"availability"`
	NumericNarrowingCountMaterialization  string `json:"count_materialization"`
	NumericNarrowingWriteAdmission        string `json:"write_admission"`
	NumericNarrowingPhysicalQualification string `json:"physical_qualification"`
	NumericNarrowingPresence              string `json:"presence"`
	NumericNarrowingFreshness             string `json:"freshness"`
	NumericNarrowingExecution             string `json:"execution"`
}

type NumericNarrowingResult struct {
	NumericNarrowingContext               NumericNarrowingContext   `json:"context"`
	NumericNarrowingPath                  string                    `json:"path"`
	NumericNarrowingOriginalInteger64     *int64                    `json:"original_integer64"`
	NumericNarrowingOriginalFloat64       *float64                  `json:"original_float64"`
	NumericNarrowingRoundedFloat64        *float64                  `json:"rounded_float64"`
	NumericNarrowingScaledFloat64         *float64                  `json:"scaled_float64"`
	NumericNarrowingRawClassification     string                    `json:"raw_classification"`
	NumericNarrowingRoundedClassification string                    `json:"rounded_classification"`
	NumericNarrowingIntermediateStatus    string                    `json:"intermediate_status"`
	NumericNarrowingUnsigned32            *uint32                   `json:"unsigned32"`
	NumericNarrowingLow16                 *uint16                   `json:"low16"`
	NumericNarrowingStatus                string                    `json:"status"`
	NumericNarrowingKnowledge             NumericNarrowingKnowledge `json:"knowledge"`
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
	numericNarrowingResult := numericNarrowingCandidate(numericNarrowingContext, "integer64")
	numericNarrowingResult.NumericNarrowingOriginalInteger64 = &numericNarrowingValue
	numericNarrowingResult.NumericNarrowingRawClassification = "finite-integer"
	numericNarrowingResult.NumericNarrowingIntermediateStatus = "not-applicable"
	numericNarrowingLow16 := uint16(numericNarrowingValue)
	numericNarrowingResult.NumericNarrowingLow16 = &numericNarrowingLow16
	numericNarrowingResult.NumericNarrowingStatus = "candidate"
	return numericNarrowingResult
}

func NumericNarrowingNarrowFloat64(numericNarrowingContext NumericNarrowingContext, numericNarrowingValue float64) NumericNarrowingResult {
	numericNarrowingResult := numericNarrowingCandidate(numericNarrowingContext, "float64")
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
