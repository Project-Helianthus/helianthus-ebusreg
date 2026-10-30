package b555

import (
	ab21an21JSON "encoding/json"
	ab21an21Math "math"
	ab21an21OS "os"
	ab21an21Reflect "reflect"
	ab21an21Strings "strings"
	ab21an21Testing "testing"
)

type ab21an21Vector struct {
	Ab21An21Label      string  `json:"ab21an21Label"`
	Ab21An21Path       string  `json:"ab21an21Path"`
	Ab21An21Integer    int64   `json:"ab21an21Integer"`
	Ab21An21Float      float64 `json:"ab21an21Float"`
	Ab21An21Rounded    float64 `json:"ab21an21Rounded"`
	Ab21An21Scaled     float64 `json:"ab21an21Scaled"`
	Ab21An21Unsigned32 uint32  `json:"ab21an21Unsigned32"`
	Ab21An21Low16      uint16  `json:"ab21an21Low16"`
}

func ab21an21CheckKnowledge(ab21an21T *ab21an21Testing.T, ab21an21Result Ab21An21Result) {
	ab21an21T.Helper()
	ab21an21Want := Ab21An21Knowledge{
		Ab21An21Conversion: "Proven", Ab21An21LocalCarrier: "Hypothesis",
		Ab21An21Units: "Unknown", Ab21An21TimerKindEffects: "Unknown",
		Ab21An21Persistence: "Unknown", Ab21An21Availability: "Unknown",
		Ab21An21CountMaterialization: "Unknown", Ab21An21WriteAdmission: "Unknown",
		Ab21An21PhysicalQualification: "Unknown",
		Ab21An21Presence:              "Unknown", Ab21An21Freshness: "Unknown", Ab21An21Execution: "Unknown",
	}
	if ab21an21Result.Ab21An21Knowledge != ab21an21Want {
		ab21an21T.Fatalf("knowledge = %+v", ab21an21Result.Ab21An21Knowledge)
	}
}

func TestAb21An21FiniteVectors(ab21an21T *ab21an21Testing.T) {
	ab21an21Data, ab21an21Err := ab21an21OS.ReadFile("testdata/ab21an21_b555_numeric_narrowing.json")
	if ab21an21Err != nil {
		ab21an21T.Fatal(ab21an21Err)
	}
	var ab21an21Vectors []ab21an21Vector
	if ab21an21Err = ab21an21JSON.Unmarshal(ab21an21Data, &ab21an21Vectors); ab21an21Err != nil {
		ab21an21T.Fatal(ab21an21Err)
	}
	if len(ab21an21Vectors) != 15 {
		ab21an21T.Fatalf("vector count = %d", len(ab21an21Vectors))
	}
	ab21an21Context := Ab21An21Context{"ab21an21-kind-a", "ab21an21-family-a", "ab21an21-software-a"}
	for ab21an21Index := range ab21an21Vectors {
		ab21an21Vector := ab21an21Vectors[ab21an21Index]
		ab21an21T.Run(ab21an21Vector.Ab21An21Label, func(ab21an21T *ab21an21Testing.T) {
			if !ab21an21Strings.HasPrefix(ab21an21Vector.Ab21An21Label, "ab21an21") {
				ab21an21T.Fatal("unprefixed label")
			}
			var ab21an21Result Ab21An21Result
			switch ab21an21Vector.Ab21An21Path {
			case "ab21an21-integer64":
				ab21an21Result = Ab21An21NarrowInteger64(ab21an21Context, ab21an21Vector.Ab21An21Integer)
				if ab21an21Result.Ab21An21OriginalInteger64 == nil || *ab21an21Result.Ab21An21OriginalInteger64 != ab21an21Vector.Ab21An21Integer || ab21an21Result.Ab21An21OriginalFloat64 != nil || ab21an21Result.Ab21An21RoundedFloat64 != nil || ab21an21Result.Ab21An21ScaledFloat64 != nil || ab21an21Result.Ab21An21Unsigned32 != nil || ab21an21Result.Ab21An21RoundedClassification != "" || ab21an21Result.Ab21An21RawClassification != "finite-integer" || ab21an21Result.Ab21An21IntermediateStatus != "not-applicable" {
					ab21an21T.Fatalf("integer evidence = %+v", ab21an21Result)
				}
			case "ab21an21-float64":
				ab21an21Result = Ab21An21NarrowFloat64(ab21an21Context, ab21an21Vector.Ab21An21Float)
				if ab21an21Result.Ab21An21OriginalFloat64 == nil || *ab21an21Result.Ab21An21OriginalFloat64 != ab21an21Vector.Ab21An21Float || ab21an21Result.Ab21An21OriginalInteger64 != nil || ab21an21Result.Ab21An21RoundedFloat64 == nil || *ab21an21Result.Ab21An21RoundedFloat64 != ab21an21Vector.Ab21An21Rounded || ab21an21Result.Ab21An21ScaledFloat64 == nil || *ab21an21Result.Ab21An21ScaledFloat64 != ab21an21Vector.Ab21An21Scaled || ab21an21Result.Ab21An21Unsigned32 == nil || *ab21an21Result.Ab21An21Unsigned32 != ab21an21Vector.Ab21An21Unsigned32 || ab21an21Result.Ab21An21IntermediateStatus != "converted" {
					ab21an21T.Fatalf("float evidence = %+v", ab21an21Result)
				}
				ab21an21Class := "positive-finite"
				if ab21an21Vector.Ab21An21Float == 0 {
					ab21an21Class = "zero"
				}
				if ab21an21Result.Ab21An21RawClassification != ab21an21Class || ab21an21Result.Ab21An21RoundedClassification != ab21an21Class {
					ab21an21T.Fatal("lost numeric classification")
				}
			default:
				ab21an21T.Fatal("unknown fixture path")
			}
			if ab21an21Result.Ab21An21Context != ab21an21Context || ab21an21Result.Ab21An21Path != ab21an21Vector.Ab21An21Path || ab21an21Result.Ab21An21Status != "candidate" || ab21an21Result.Ab21An21Low16 == nil || *ab21an21Result.Ab21An21Low16 != ab21an21Vector.Ab21An21Low16 {
				ab21an21T.Fatalf("result = %+v, want low16 %d", ab21an21Result, ab21an21Vector.Ab21An21Low16)
			}
			ab21an21CheckKnowledge(ab21an21T, ab21an21Result)
		})
	}
}

func TestAb21An21Unresolved(ab21an21T *ab21an21Testing.T) {
	ab21an21Cases := []struct {
		Ab21An21Label        string
		Ab21An21Value        float64
		Ab21An21Class        string
		Ab21An21RoundedClass string
		Ab21An21Intermediate string
	}{
		{"ab21an21-negative", -1, "negative-finite", "negative-finite", "unsupported-negative-input"},
		{"ab21an21-negative-underflow", -ab21an21Math.SmallestNonzeroFloat64, "negative-finite", "negative-zero", "unsupported-negative-input"},
		{"ab21an21-nan", ab21an21Math.NaN(), "nan", "nan", "unsupported-nonfinite-input"},
		{"ab21an21-positive-infinity", ab21an21Math.Inf(1), "positive-infinity", "positive-infinity", "unsupported-nonfinite-input"},
		{"ab21an21-negative-infinity", ab21an21Math.Inf(-1), "negative-infinity", "negative-infinity", "unsupported-nonfinite-input"},
		{"ab21an21-scaled-overflow", 429496729.5, "positive-finite", "positive-finite", "unsupported-scaled-range"},
		{"ab21an21-rounded-overflow", ab21an21Math.MaxFloat64, "positive-finite", "positive-infinity", "unsupported-nonfinite-scaled"},
	}
	ab21an21Context := Ab21An21Context{"ab21an21-kind-a", "ab21an21-family-a", "ab21an21-software-a"}
	ab21an21Previous := Ab21An21NarrowInteger64(ab21an21Context, 65537)
	for ab21an21Index := range ab21an21Cases {
		ab21an21Case := ab21an21Cases[ab21an21Index]
		ab21an21T.Run(ab21an21Case.Ab21An21Label, func(ab21an21T *ab21an21Testing.T) {
			ab21an21Result := Ab21An21NarrowFloat64(ab21an21Context, ab21an21Case.Ab21An21Value)
			if ab21an21Result.Ab21An21Status != "unresolved" || ab21an21Result.Ab21An21Low16 != nil || ab21an21Result.Ab21An21Unsigned32 != nil || ab21an21Result.Ab21An21OriginalInteger64 != nil || ab21an21Result.Ab21An21OriginalFloat64 == nil || ab21an21Result.Ab21An21RoundedFloat64 == nil || ab21an21Result.Ab21An21ScaledFloat64 == nil {
				ab21an21T.Fatalf("fabricated or lost evidence: %+v", ab21an21Result)
			}
			if ab21an21Math.Float64bits(*ab21an21Result.Ab21An21OriginalFloat64) != ab21an21Math.Float64bits(ab21an21Case.Ab21An21Value) || ab21an21Result.Ab21An21RawClassification != ab21an21Case.Ab21An21Class || ab21an21Result.Ab21An21RoundedClassification != ab21an21Case.Ab21An21RoundedClass || ab21an21Result.Ab21An21IntermediateStatus != ab21an21Case.Ab21An21Intermediate || ab21an21Result.Ab21An21Context != ab21an21Context || ab21an21Result.Ab21An21Path != "ab21an21-float64" {
				ab21an21T.Fatalf("unresolved evidence = %+v", ab21an21Result)
			}
			ab21an21CheckKnowledge(ab21an21T, ab21an21Result)
			if *ab21an21Previous.Ab21An21Low16 != 1 || ab21an21Previous.Ab21An21Status != "candidate" || *ab21an21Previous.Ab21An21OriginalInteger64 != 65537 {
				ab21an21T.Fatal("previous useful result changed")
			}
		})
	}
}

func TestAb21An21ContextAndCallerIsolation(ab21an21T *ab21an21Testing.T) {
	ab21an21Contexts := []Ab21An21Context{
		{"ab21an21-kind-a", "ab21an21-family-a", "ab21an21-software-a"},
		{"ab21an21-kind-b", "ab21an21-family-a", "ab21an21-software-a"},
		{"ab21an21-kind-a", "ab21an21-family-b", "ab21an21-software-a"},
		{"ab21an21-kind-a", "ab21an21-family-a", "ab21an21-software-b"},
		{"", "ab21an21-family-a", "ab21an21-software-a"},
		{"ab21an21-kind-a", "", "ab21an21-software-a"},
		{"ab21an21-kind-a", "ab21an21-family-a", ""},
		{},
	}
	for ab21an21Index := range ab21an21Contexts {
		ab21an21Context := ab21an21Contexts[ab21an21Index]
		ab21an21Saved := ab21an21Context
		ab21an21Integer := int64(1)
		ab21an21Float := float64(1)
		ab21an21IntegerResult := Ab21An21NarrowInteger64(ab21an21Context, ab21an21Integer)
		ab21an21FloatResult := Ab21An21NarrowFloat64(ab21an21Context, ab21an21Float)
		ab21an21Repeat := Ab21An21NarrowFloat64(ab21an21Context, ab21an21Float)
		if !ab21an21Reflect.DeepEqual(ab21an21FloatResult, ab21an21Repeat) {
			ab21an21T.Fatal("caller-independent repeat differs")
		}
		ab21an21Context = Ab21An21Context{"ab21an21-changed", "ab21an21-changed", "ab21an21-changed"}
		ab21an21Integer, ab21an21Float = 99, 99
		ab21an21ChangedIntegerResult := Ab21An21NarrowInteger64(ab21an21Context, ab21an21Integer)
		ab21an21ChangedFloatResult := Ab21An21NarrowFloat64(ab21an21Context, ab21an21Float)
		if ab21an21ChangedIntegerResult.Ab21An21Context != ab21an21Context || ab21an21ChangedFloatResult.Ab21An21Context != ab21an21Context || *ab21an21ChangedIntegerResult.Ab21An21OriginalInteger64 != 99 || *ab21an21ChangedFloatResult.Ab21An21OriginalFloat64 != 99 || *ab21an21ChangedIntegerResult.Ab21An21Low16 != 99 || *ab21an21ChangedFloatResult.Ab21An21Low16 != 990 {
			ab21an21T.Fatal("changed caller context or input not reflected in new results")
		}
		if ab21an21IntegerResult.Ab21An21Context != ab21an21Saved || ab21an21FloatResult.Ab21An21Context != ab21an21Saved || *ab21an21IntegerResult.Ab21An21OriginalInteger64 != 1 || *ab21an21FloatResult.Ab21An21OriginalFloat64 != 1 || *ab21an21IntegerResult.Ab21An21Low16 != 1 || *ab21an21FloatResult.Ab21An21Low16 != 10 || ab21an21IntegerResult.Ab21An21Path != "ab21an21-integer64" || ab21an21FloatResult.Ab21An21Path != "ab21an21-float64" {
			ab21an21T.Fatal("context, input or selection not retained independently")
		}
		*ab21an21FloatResult.Ab21An21Low16 = 999
		*ab21an21FloatResult.Ab21An21RoundedFloat64 = 999
		if *ab21an21Repeat.Ab21An21Low16 != 10 || *ab21an21Repeat.Ab21An21RoundedFloat64 != 1 || *ab21an21IntegerResult.Ab21An21Low16 != 1 {
			ab21an21T.Fatal("results alias")
		}
		ab21an21CheckKnowledge(ab21an21T, ab21an21IntegerResult)
		ab21an21CheckKnowledge(ab21an21T, ab21an21FloatResult)
	}
}

func TestAb21An21NegativeControlsAndZero(ab21an21T *ab21an21Testing.T) {
	ab21an21Value := 0.7
	if uint16(ab21an21Value*10) != 7 || *Ab21An21NarrowFloat64(Ab21An21Context{}, ab21an21Value).Ab21An21Low16 != 6 {
		ab21an21T.Fatal("binary32 negative control")
	}
	if ab21an21Math.Round(12.5) != 13 || *Ab21An21NarrowFloat64(Ab21An21Context{}, 1.25).Ab21An21Low16 != 12 {
		ab21an21T.Fatal("rounding negative control")
	}
	ab21an21Values := []float64{ab21an21Math.Copysign(0, -1), ab21an21Math.SmallestNonzeroFloat64}
	for ab21an21Index := range ab21an21Values {
		ab21an21Value := ab21an21Values[ab21an21Index]
		ab21an21Result := Ab21An21NarrowFloat64(Ab21An21Context{}, ab21an21Value)
		if ab21an21Result.Ab21An21Status != "candidate" || ab21an21Result.Ab21An21Low16 == nil || *ab21an21Result.Ab21An21Low16 != 0 || ab21an21Result.Ab21An21Unsigned32 == nil || *ab21an21Result.Ab21An21Unsigned32 != 0 || ab21an21Math.Float64bits(*ab21an21Result.Ab21An21OriginalFloat64) != ab21an21Math.Float64bits(ab21an21Value) {
			ab21an21T.Fatal("zero lost")
		}
	}
}
