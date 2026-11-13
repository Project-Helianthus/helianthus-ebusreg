package b555

import (
	numericNarrowingJSON "encoding/json"
	numericNarrowingMath "math"
	numericNarrowingOS "os"
	numericNarrowingReflect "reflect"
	numericNarrowingStrings "strings"
	numericNarrowingTesting "testing"
)

type numericNarrowingVector struct {
	NumericNarrowingLabel      string  `json:"ab21an21Label"`
	NumericNarrowingPath       string  `json:"ab21an21Path"`
	NumericNarrowingInteger    int64   `json:"ab21an21Integer"`
	NumericNarrowingFloat      float64 `json:"ab21an21Float"`
	NumericNarrowingRounded    float64 `json:"ab21an21Rounded"`
	NumericNarrowingScaled     float64 `json:"ab21an21Scaled"`
	NumericNarrowingUnsigned32 uint32  `json:"ab21an21Unsigned32"`
	NumericNarrowingLow16      uint16  `json:"ab21an21Low16"`
}

func numericNarrowingCheckKnowledge(numericNarrowingT *numericNarrowingTesting.T, numericNarrowingResult NumericNarrowingResult) {
	numericNarrowingT.Helper()
	numericNarrowingWant := NumericNarrowingKnowledge{
		NumericNarrowingConversion: "Proven", NumericNarrowingLocalCarrier: "Hypothesis",
		NumericNarrowingUnits: "Unknown", NumericNarrowingTimerKindEffects: "Unknown",
		NumericNarrowingPersistence: "Unknown", NumericNarrowingAvailability: "Unknown",
		NumericNarrowingCountMaterialization: "Unknown", NumericNarrowingWriteAdmission: "Unknown",
		NumericNarrowingPhysicalQualification: "Unknown",
		NumericNarrowingPresence:              "Unknown", NumericNarrowingFreshness: "Unknown", NumericNarrowingExecution: "Unknown",
	}
	if numericNarrowingResult.NumericNarrowingKnowledge != numericNarrowingWant {
		numericNarrowingT.Fatalf("knowledge = %+v", numericNarrowingResult.NumericNarrowingKnowledge)
	}
}

func TestAb21An21FiniteVectors(numericNarrowingT *numericNarrowingTesting.T) {
	numericNarrowingData, numericNarrowingErr := numericNarrowingOS.ReadFile("testdata/ab21an21_b555_numeric_narrowing.json")
	if numericNarrowingErr != nil {
		numericNarrowingT.Fatal(numericNarrowingErr)
	}
	var numericNarrowingVectors []numericNarrowingVector
	if numericNarrowingErr = numericNarrowingJSON.Unmarshal(numericNarrowingData, &numericNarrowingVectors); numericNarrowingErr != nil {
		numericNarrowingT.Fatal(numericNarrowingErr)
	}
	if len(numericNarrowingVectors) != 15 {
		numericNarrowingT.Fatalf("vector count = %d", len(numericNarrowingVectors))
	}
	numericNarrowingContext := NumericNarrowingContext{"ab21an21-kind-a", "ab21an21-family-a", "ab21an21-software-a"}
	for numericNarrowingIndex := range numericNarrowingVectors {
		numericNarrowingVector := numericNarrowingVectors[numericNarrowingIndex]
		numericNarrowingT.Run(numericNarrowingVector.NumericNarrowingLabel, func(numericNarrowingT *numericNarrowingTesting.T) {
			if !numericNarrowingStrings.HasPrefix(numericNarrowingVector.NumericNarrowingLabel, "ab21an21") {
				numericNarrowingT.Fatal("unprefixed label")
			}
			var numericNarrowingResult NumericNarrowingResult
			switch numericNarrowingVector.NumericNarrowingPath {
			case "ab21an21-integer64":
				numericNarrowingResult = NumericNarrowingNarrowInteger64(numericNarrowingContext, numericNarrowingVector.NumericNarrowingInteger)
				if numericNarrowingResult.NumericNarrowingOriginalInteger64 == nil || *numericNarrowingResult.NumericNarrowingOriginalInteger64 != numericNarrowingVector.NumericNarrowingInteger || numericNarrowingResult.NumericNarrowingOriginalFloat64 != nil || numericNarrowingResult.NumericNarrowingRoundedFloat64 != nil || numericNarrowingResult.NumericNarrowingScaledFloat64 != nil || numericNarrowingResult.NumericNarrowingUnsigned32 != nil || numericNarrowingResult.NumericNarrowingRoundedClassification != "" || numericNarrowingResult.NumericNarrowingRawClassification != "finite-integer" || numericNarrowingResult.NumericNarrowingIntermediateStatus != "not-applicable" {
					numericNarrowingT.Fatalf("integer evidence = %+v", numericNarrowingResult)
				}
			case "ab21an21-float64":
				numericNarrowingResult = NumericNarrowingNarrowFloat64(numericNarrowingContext, numericNarrowingVector.NumericNarrowingFloat)
				if numericNarrowingResult.NumericNarrowingOriginalFloat64 == nil || *numericNarrowingResult.NumericNarrowingOriginalFloat64 != numericNarrowingVector.NumericNarrowingFloat || numericNarrowingResult.NumericNarrowingOriginalInteger64 != nil || numericNarrowingResult.NumericNarrowingRoundedFloat64 == nil || *numericNarrowingResult.NumericNarrowingRoundedFloat64 != numericNarrowingVector.NumericNarrowingRounded || numericNarrowingResult.NumericNarrowingScaledFloat64 == nil || *numericNarrowingResult.NumericNarrowingScaledFloat64 != numericNarrowingVector.NumericNarrowingScaled || numericNarrowingResult.NumericNarrowingUnsigned32 == nil || *numericNarrowingResult.NumericNarrowingUnsigned32 != numericNarrowingVector.NumericNarrowingUnsigned32 || numericNarrowingResult.NumericNarrowingIntermediateStatus != "converted" {
					numericNarrowingT.Fatalf("float evidence = %+v", numericNarrowingResult)
				}
				numericNarrowingClass := "positive-finite"
				if numericNarrowingVector.NumericNarrowingFloat == 0 {
					numericNarrowingClass = "zero"
				}
				if numericNarrowingResult.NumericNarrowingRawClassification != numericNarrowingClass || numericNarrowingResult.NumericNarrowingRoundedClassification != numericNarrowingClass {
					numericNarrowingT.Fatal("lost numeric classification")
				}
			default:
				numericNarrowingT.Fatal("unknown fixture path")
			}
			if numericNarrowingResult.NumericNarrowingContext != numericNarrowingContext || numericNarrowingResult.NumericNarrowingPath != numericNarrowingVector.NumericNarrowingPath || numericNarrowingResult.NumericNarrowingStatus != "candidate" || numericNarrowingResult.NumericNarrowingLow16 == nil || *numericNarrowingResult.NumericNarrowingLow16 != numericNarrowingVector.NumericNarrowingLow16 {
				numericNarrowingT.Fatalf("result = %+v, want low16 %d", numericNarrowingResult, numericNarrowingVector.NumericNarrowingLow16)
			}
			numericNarrowingCheckKnowledge(numericNarrowingT, numericNarrowingResult)
		})
	}
}

func TestAb21An21Unresolved(numericNarrowingT *numericNarrowingTesting.T) {
	numericNarrowingCases := []struct {
		NumericNarrowingLabel        string
		NumericNarrowingValue        float64
		NumericNarrowingClass        string
		NumericNarrowingRoundedClass string
		NumericNarrowingIntermediate string
	}{
		{"ab21an21-negative", -1, "negative-finite", "negative-finite", "unsupported-negative-input"},
		{"ab21an21-negative-underflow", -numericNarrowingMath.SmallestNonzeroFloat64, "negative-finite", "negative-zero", "unsupported-negative-input"},
		{"ab21an21-nan", numericNarrowingMath.NaN(), "nan", "nan", "unsupported-nonfinite-input"},
		{"ab21an21-positive-infinity", numericNarrowingMath.Inf(1), "positive-infinity", "positive-infinity", "unsupported-nonfinite-input"},
		{"ab21an21-negative-infinity", numericNarrowingMath.Inf(-1), "negative-infinity", "negative-infinity", "unsupported-nonfinite-input"},
		{"ab21an21-scaled-overflow", 429496729.5, "positive-finite", "positive-finite", "unsupported-scaled-range"},
		{"ab21an21-rounded-overflow", numericNarrowingMath.MaxFloat64, "positive-finite", "positive-infinity", "unsupported-nonfinite-scaled"},
	}
	numericNarrowingContext := NumericNarrowingContext{"ab21an21-kind-a", "ab21an21-family-a", "ab21an21-software-a"}
	numericNarrowingPrevious := NumericNarrowingNarrowInteger64(numericNarrowingContext, 65537)
	for numericNarrowingIndex := range numericNarrowingCases {
		numericNarrowingCase := numericNarrowingCases[numericNarrowingIndex]
		numericNarrowingT.Run(numericNarrowingCase.NumericNarrowingLabel, func(numericNarrowingT *numericNarrowingTesting.T) {
			numericNarrowingResult := NumericNarrowingNarrowFloat64(numericNarrowingContext, numericNarrowingCase.NumericNarrowingValue)
			if numericNarrowingResult.NumericNarrowingStatus != "unresolved" || numericNarrowingResult.NumericNarrowingLow16 != nil || numericNarrowingResult.NumericNarrowingUnsigned32 != nil || numericNarrowingResult.NumericNarrowingOriginalInteger64 != nil || numericNarrowingResult.NumericNarrowingOriginalFloat64 == nil || numericNarrowingResult.NumericNarrowingRoundedFloat64 == nil || numericNarrowingResult.NumericNarrowingScaledFloat64 == nil {
				numericNarrowingT.Fatalf("fabricated or lost evidence: %+v", numericNarrowingResult)
			}
			if numericNarrowingMath.Float64bits(*numericNarrowingResult.NumericNarrowingOriginalFloat64) != numericNarrowingMath.Float64bits(numericNarrowingCase.NumericNarrowingValue) || numericNarrowingResult.NumericNarrowingRawClassification != numericNarrowingCase.NumericNarrowingClass || numericNarrowingResult.NumericNarrowingRoundedClassification != numericNarrowingCase.NumericNarrowingRoundedClass || numericNarrowingResult.NumericNarrowingIntermediateStatus != numericNarrowingCase.NumericNarrowingIntermediate || numericNarrowingResult.NumericNarrowingContext != numericNarrowingContext || numericNarrowingResult.NumericNarrowingPath != "ab21an21-float64" {
				numericNarrowingT.Fatalf("unresolved evidence = %+v", numericNarrowingResult)
			}
			numericNarrowingCheckKnowledge(numericNarrowingT, numericNarrowingResult)
			if *numericNarrowingPrevious.NumericNarrowingLow16 != 1 || numericNarrowingPrevious.NumericNarrowingStatus != "candidate" || *numericNarrowingPrevious.NumericNarrowingOriginalInteger64 != 65537 {
				numericNarrowingT.Fatal("previous useful result changed")
			}
		})
	}
}

func TestAb21An21ContextAndCallerIsolation(numericNarrowingT *numericNarrowingTesting.T) {
	numericNarrowingContexts := []NumericNarrowingContext{
		{"ab21an21-kind-a", "ab21an21-family-a", "ab21an21-software-a"},
		{"ab21an21-kind-b", "ab21an21-family-a", "ab21an21-software-a"},
		{"ab21an21-kind-a", "ab21an21-family-b", "ab21an21-software-a"},
		{"ab21an21-kind-a", "ab21an21-family-a", "ab21an21-software-b"},
		{"", "ab21an21-family-a", "ab21an21-software-a"},
		{"ab21an21-kind-a", "", "ab21an21-software-a"},
		{"ab21an21-kind-a", "ab21an21-family-a", ""},
		{},
	}
	for numericNarrowingIndex := range numericNarrowingContexts {
		numericNarrowingContext := numericNarrowingContexts[numericNarrowingIndex]
		numericNarrowingSaved := numericNarrowingContext
		numericNarrowingInteger := int64(1)
		numericNarrowingFloat := float64(1)
		numericNarrowingIntegerResult := NumericNarrowingNarrowInteger64(numericNarrowingContext, numericNarrowingInteger)
		numericNarrowingFloatResult := NumericNarrowingNarrowFloat64(numericNarrowingContext, numericNarrowingFloat)
		numericNarrowingRepeat := NumericNarrowingNarrowFloat64(numericNarrowingContext, numericNarrowingFloat)
		if !numericNarrowingReflect.DeepEqual(numericNarrowingFloatResult, numericNarrowingRepeat) {
			numericNarrowingT.Fatal("caller-independent repeat differs")
		}
		numericNarrowingContext = NumericNarrowingContext{"ab21an21-changed", "ab21an21-changed", "ab21an21-changed"}
		numericNarrowingInteger, numericNarrowingFloat = 99, 99
		numericNarrowingChangedIntegerResult := NumericNarrowingNarrowInteger64(numericNarrowingContext, numericNarrowingInteger)
		numericNarrowingChangedFloatResult := NumericNarrowingNarrowFloat64(numericNarrowingContext, numericNarrowingFloat)
		if numericNarrowingChangedIntegerResult.NumericNarrowingContext != numericNarrowingContext || numericNarrowingChangedFloatResult.NumericNarrowingContext != numericNarrowingContext || *numericNarrowingChangedIntegerResult.NumericNarrowingOriginalInteger64 != 99 || *numericNarrowingChangedFloatResult.NumericNarrowingOriginalFloat64 != 99 || *numericNarrowingChangedIntegerResult.NumericNarrowingLow16 != 99 || *numericNarrowingChangedFloatResult.NumericNarrowingLow16 != 990 {
			numericNarrowingT.Fatal("changed caller context or input not reflected in new results")
		}
		if numericNarrowingIntegerResult.NumericNarrowingContext != numericNarrowingSaved || numericNarrowingFloatResult.NumericNarrowingContext != numericNarrowingSaved || *numericNarrowingIntegerResult.NumericNarrowingOriginalInteger64 != 1 || *numericNarrowingFloatResult.NumericNarrowingOriginalFloat64 != 1 || *numericNarrowingIntegerResult.NumericNarrowingLow16 != 1 || *numericNarrowingFloatResult.NumericNarrowingLow16 != 10 || numericNarrowingIntegerResult.NumericNarrowingPath != "ab21an21-integer64" || numericNarrowingFloatResult.NumericNarrowingPath != "ab21an21-float64" {
			numericNarrowingT.Fatal("context, input or selection not retained independently")
		}
		*numericNarrowingFloatResult.NumericNarrowingLow16 = 999
		*numericNarrowingFloatResult.NumericNarrowingRoundedFloat64 = 999
		if *numericNarrowingRepeat.NumericNarrowingLow16 != 10 || *numericNarrowingRepeat.NumericNarrowingRoundedFloat64 != 1 || *numericNarrowingIntegerResult.NumericNarrowingLow16 != 1 {
			numericNarrowingT.Fatal("results alias")
		}
		numericNarrowingCheckKnowledge(numericNarrowingT, numericNarrowingIntegerResult)
		numericNarrowingCheckKnowledge(numericNarrowingT, numericNarrowingFloatResult)
	}
}

func TestAb21An21NegativeControlsAndZero(numericNarrowingT *numericNarrowingTesting.T) {
	numericNarrowingValue := 0.7
	if uint16(numericNarrowingValue*10) != 7 || *NumericNarrowingNarrowFloat64(NumericNarrowingContext{}, numericNarrowingValue).NumericNarrowingLow16 != 6 {
		numericNarrowingT.Fatal("binary32 negative control")
	}
	if numericNarrowingMath.Round(12.5) != 13 || *NumericNarrowingNarrowFloat64(NumericNarrowingContext{}, 1.25).NumericNarrowingLow16 != 12 {
		numericNarrowingT.Fatal("rounding negative control")
	}
	numericNarrowingValues := []float64{numericNarrowingMath.Copysign(0, -1), numericNarrowingMath.SmallestNonzeroFloat64}
	for numericNarrowingIndex := range numericNarrowingValues {
		numericNarrowingValue := numericNarrowingValues[numericNarrowingIndex]
		numericNarrowingResult := NumericNarrowingNarrowFloat64(NumericNarrowingContext{}, numericNarrowingValue)
		if numericNarrowingResult.NumericNarrowingStatus != "candidate" || numericNarrowingResult.NumericNarrowingLow16 == nil || *numericNarrowingResult.NumericNarrowingLow16 != 0 || numericNarrowingResult.NumericNarrowingUnsigned32 == nil || *numericNarrowingResult.NumericNarrowingUnsigned32 != 0 || numericNarrowingMath.Float64bits(*numericNarrowingResult.NumericNarrowingOriginalFloat64) != numericNarrowingMath.Float64bits(numericNarrowingValue) {
			numericNarrowingT.Fatal("zero lost")
		}
	}
}
