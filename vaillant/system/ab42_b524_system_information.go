package system

import (
	ab42binary "encoding/binary"
	ab42math "math"
)

type Ab42InformationContext struct {
	Ab42Family      *string
	Ab42Software    *string
	Ab42Hardware    *string
	Ab42APIVersion  *uint16
	Ab42APIRevision *uint16
	Ab42Identifier  uint16
}

type Ab42InformationResult struct {
	Ab42Context             Ab42InformationContext
	Ab42Raw                 []byte
	Ab42Valid               bool
	Ab42InvalidReason       string
	Ab42Value               float32
	Ab42Meaning             string
	Ab42Capacity            *float32
	Ab42HeatPumpCount       *float32
	Ab42RepresentationGrade string
	Ab42LayoutGrade         string
	Ab42QualificationGrade  string
	Ab42PresenceGrade       string
	Ab42EligibilityGrade    string
}

func Ab42BuildInformation(ab42Identifier uint16) []byte {
	return []byte{0, byte(ab42Identifier), byte(ab42Identifier >> 8)}
}

func ab42Copy[Ab42Type any](ab42Source *Ab42Type) *Ab42Type {
	if ab42Source == nil {
		return nil
	}
	ab42Value := *ab42Source
	return &ab42Value
}

func Ab42DecodeInformation(ab42RequestContext Ab42InformationContext, ab42Raw []byte) Ab42InformationResult {
	ab42Result := Ab42InformationResult{
		Ab42Context: Ab42InformationContext{
			Ab42Family:      ab42Copy(ab42RequestContext.Ab42Family),
			Ab42Software:    ab42Copy(ab42RequestContext.Ab42Software),
			Ab42Hardware:    ab42Copy(ab42RequestContext.Ab42Hardware),
			Ab42APIVersion:  ab42Copy(ab42RequestContext.Ab42APIVersion),
			Ab42APIRevision: ab42Copy(ab42RequestContext.Ab42APIRevision),
			Ab42Identifier:  ab42RequestContext.Ab42Identifier,
		},
		Ab42Raw:                 append([]byte(nil), ab42Raw...),
		Ab42Meaning:             "Unknown",
		Ab42RepresentationGrade: "Unknown",
		Ab42LayoutGrade:         "Hypothesis",
		Ab42QualificationGrade:  "Unknown",
		Ab42PresenceGrade:       "Unknown",
		Ab42EligibilityGrade:    "Unknown",
	}
	switch ab42RequestContext.Ab42Identifier {
	case 0:
		ab42Result.Ab42Meaning = "supported circuit capacity"
	case 1:
		ab42Result.Ab42Meaning = "supported zone capacity"
	case 0x000d:
		ab42Result.Ab42Meaning = "heat-pump count"
	}
	if len(ab42Raw) != 4 {
		ab42Result.Ab42InvalidReason = "ab42 requires exactly four normalized bytes"
		return ab42Result
	}
	ab42Value := ab42math.Float32frombits(ab42binary.LittleEndian.Uint32(ab42Raw))
	if ab42math.IsNaN(float64(ab42Value)) || ab42math.IsInf(float64(ab42Value), 0) {
		ab42Result.Ab42InvalidReason = "ab42 non-finite binary32"
		return ab42Result
	}
	ab42Result.Ab42Valid = true
	ab42Result.Ab42Value = ab42Value
	ab42Result.Ab42RepresentationGrade = "Proven"
	if ab42Value >= 0 && ab42math.Trunc(float64(ab42Value)) == float64(ab42Value) {
		switch ab42RequestContext.Ab42Identifier {
		case 0, 1:
			ab42Result.Ab42Capacity = ab42Copy(&ab42Value)
		case 0x000d:
			ab42Result.Ab42HeatPumpCount = ab42Copy(&ab42Value)
		}
	}
	return ab42Result
}
