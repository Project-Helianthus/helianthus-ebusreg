package b555

import ab21an21Math "math"

type Ab21An21Context struct {
	Ab21An21TimerKind       string
	Ab21An21Family          string
	Ab21An21SoftwareVariant string
}

type Ab21An21Knowledge struct {
	Ab21An21Conversion            string
	Ab21An21LocalCarrier          string
	Ab21An21Units                 string
	Ab21An21TimerKindEffects      string
	Ab21An21Persistence           string
	Ab21An21Availability          string
	Ab21An21CountMaterialization  string
	Ab21An21WriteAdmission        string
	Ab21An21PhysicalQualification string
	Ab21An21Presence              string
	Ab21An21Freshness             string
	Ab21An21Execution             string
}

type Ab21An21Result struct {
	Ab21An21Context               Ab21An21Context
	Ab21An21Path                  string
	Ab21An21OriginalInteger64     *int64
	Ab21An21OriginalFloat64       *float64
	Ab21An21RoundedFloat64        *float64
	Ab21An21ScaledFloat64         *float64
	Ab21An21RawClassification     string
	Ab21An21RoundedClassification string
	Ab21An21IntermediateStatus    string
	Ab21An21Unsigned32            *uint32
	Ab21An21Low16                 *uint16
	Ab21An21Status                string
	Ab21An21Knowledge             Ab21An21Knowledge
}

func ab21an21Candidate(ab21an21Context Ab21An21Context, ab21an21Path string) Ab21An21Result {
	return Ab21An21Result{
		Ab21An21Context: ab21an21Context,
		Ab21An21Path:    ab21an21Path,
		Ab21An21Status:  "unresolved",
		Ab21An21Knowledge: Ab21An21Knowledge{
			Ab21An21Conversion:            "Proven",
			Ab21An21LocalCarrier:          "Hypothesis",
			Ab21An21Units:                 "Unknown",
			Ab21An21TimerKindEffects:      "Unknown",
			Ab21An21Persistence:           "Unknown",
			Ab21An21Availability:          "Unknown",
			Ab21An21CountMaterialization:  "Unknown",
			Ab21An21WriteAdmission:        "Unknown",
			Ab21An21PhysicalQualification: "Unknown",
			Ab21An21Presence:              "Unknown",
			Ab21An21Freshness:             "Unknown",
			Ab21An21Execution:             "Unknown",
		},
	}
}

func Ab21An21NarrowInteger64(ab21an21Context Ab21An21Context, ab21an21Value int64) Ab21An21Result {
	ab21an21Result := ab21an21Candidate(ab21an21Context, "ab21an21-integer64")
	ab21an21Result.Ab21An21OriginalInteger64 = &ab21an21Value
	ab21an21Result.Ab21An21RawClassification = "finite-integer"
	ab21an21Result.Ab21An21IntermediateStatus = "not-applicable"
	ab21an21Low16 := uint16(ab21an21Value)
	ab21an21Result.Ab21An21Low16 = &ab21an21Low16
	ab21an21Result.Ab21An21Status = "candidate"
	return ab21an21Result
}

func Ab21An21NarrowFloat64(ab21an21Context Ab21An21Context, ab21an21Value float64) Ab21An21Result {
	ab21an21Result := ab21an21Candidate(ab21an21Context, "ab21an21-float64")
	ab21an21Result.Ab21An21OriginalFloat64 = &ab21an21Value
	ab21an21Result.Ab21An21RawClassification = ab21an21Classify(ab21an21Value)
	ab21an21Rounded := float64(float32(ab21an21Value))
	ab21an21Scaled := ab21an21Rounded * 10
	ab21an21Result.Ab21An21RoundedFloat64 = &ab21an21Rounded
	ab21an21Result.Ab21An21ScaledFloat64 = &ab21an21Scaled
	ab21an21Result.Ab21An21RoundedClassification = ab21an21Classify(ab21an21Rounded)
	switch {
	case ab21an21Math.IsNaN(ab21an21Value) || ab21an21Math.IsInf(ab21an21Value, 0):
		ab21an21Result.Ab21An21IntermediateStatus = "unsupported-nonfinite-input"
	case ab21an21Value < 0:
		ab21an21Result.Ab21An21IntermediateStatus = "unsupported-negative-input"
	case ab21an21Math.IsNaN(ab21an21Scaled) || ab21an21Math.IsInf(ab21an21Scaled, 0):
		ab21an21Result.Ab21An21IntermediateStatus = "unsupported-nonfinite-scaled"
	case ab21an21Scaled < 0 || ab21an21Scaled > 4294967295:
		ab21an21Result.Ab21An21IntermediateStatus = "unsupported-scaled-range"
	default:
		ab21an21Unsigned32 := uint32(ab21an21Math.Trunc(ab21an21Scaled))
		ab21an21Low16 := uint16(ab21an21Unsigned32)
		ab21an21Result.Ab21An21Unsigned32 = &ab21an21Unsigned32
		ab21an21Result.Ab21An21Low16 = &ab21an21Low16
		ab21an21Result.Ab21An21IntermediateStatus = "converted"
		ab21an21Result.Ab21An21Status = "candidate"
	}
	return ab21an21Result
}

func ab21an21Classify(ab21an21Value float64) string {
	switch {
	case ab21an21Math.IsNaN(ab21an21Value):
		return "nan"
	case ab21an21Math.IsInf(ab21an21Value, 1):
		return "positive-infinity"
	case ab21an21Math.IsInf(ab21an21Value, -1):
		return "negative-infinity"
	case ab21an21Value == 0 && ab21an21Math.Signbit(ab21an21Value):
		return "negative-zero"
	case ab21an21Value == 0:
		return "zero"
	case ab21an21Value < 0:
		return "negative-finite"
	default:
		return "positive-finite"
	}
}
