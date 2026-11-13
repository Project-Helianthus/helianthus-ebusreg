package system

import (
	b524SystemInformationBinary "encoding/binary"
	b524SystemInformationMath "math"
)

type B524SystemInformationInformationContext struct {
	B524SystemInformationFamily      *string `json:"Ab42Family"`
	B524SystemInformationSoftware    *string `json:"Ab42Software"`
	B524SystemInformationHardware    *string `json:"Ab42Hardware"`
	B524SystemInformationAPIVersion  *uint16 `json:"Ab42APIVersion"`
	B524SystemInformationAPIRevision *uint16 `json:"Ab42APIRevision"`
	B524SystemInformationIdentifier  uint16  `json:"Ab42Identifier"`
}

type B524SystemInformationResult struct {
	B524SystemInformationContext             B524SystemInformationInformationContext `json:"Ab42Context"`
	B524SystemInformationRaw                 []byte                                  `json:"Ab42Raw"`
	B524SystemInformationValid               bool                                    `json:"Ab42Valid"`
	B524SystemInformationInvalidReason       string                                  `json:"Ab42InvalidReason"`
	B524SystemInformationValue               float32                                 `json:"Ab42Value"`
	B524SystemInformationMeaning             string                                  `json:"Ab42Meaning"`
	B524SystemInformationCapacity            *float32                                `json:"Ab42Capacity"`
	B524SystemInformationHeatPumpCount       *float32                                `json:"Ab42HeatPumpCount"`
	B524SystemInformationRepresentationGrade string                                  `json:"Ab42RepresentationGrade"`
	B524SystemInformationLayoutGrade         string                                  `json:"Ab42LayoutGrade"`
	B524SystemInformationQualificationGrade  string                                  `json:"Ab42QualificationGrade"`
	B524SystemInformationPresenceGrade       string                                  `json:"Ab42PresenceGrade"`
	B524SystemInformationEligibilityGrade    string                                  `json:"Ab42EligibilityGrade"`
}

func B524SystemBuildInformation(b524SystemInformationIdentifier uint16) []byte {
	return []byte{0, byte(b524SystemInformationIdentifier), byte(b524SystemInformationIdentifier >> 8)}
}

func b524SystemInformationCopy[B524SystemInformationType any](b524SystemInformationSource *B524SystemInformationType) *B524SystemInformationType {
	if b524SystemInformationSource == nil {
		return nil
	}
	b524SystemInformationValue := *b524SystemInformationSource
	return &b524SystemInformationValue
}

func B524SystemDecodeInformation(b524SystemInformationRequestContext B524SystemInformationInformationContext, b524SystemInformationRaw []byte) B524SystemInformationResult {
	b524SystemInformationResult := B524SystemInformationResult{
		B524SystemInformationContext: B524SystemInformationInformationContext{
			B524SystemInformationFamily:      b524SystemInformationCopy(b524SystemInformationRequestContext.B524SystemInformationFamily),
			B524SystemInformationSoftware:    b524SystemInformationCopy(b524SystemInformationRequestContext.B524SystemInformationSoftware),
			B524SystemInformationHardware:    b524SystemInformationCopy(b524SystemInformationRequestContext.B524SystemInformationHardware),
			B524SystemInformationAPIVersion:  b524SystemInformationCopy(b524SystemInformationRequestContext.B524SystemInformationAPIVersion),
			B524SystemInformationAPIRevision: b524SystemInformationCopy(b524SystemInformationRequestContext.B524SystemInformationAPIRevision),
			B524SystemInformationIdentifier:  b524SystemInformationRequestContext.B524SystemInformationIdentifier,
		},
		B524SystemInformationRaw:                 append([]byte(nil), b524SystemInformationRaw...),
		B524SystemInformationMeaning:             "Unknown",
		B524SystemInformationRepresentationGrade: "Unknown",
		B524SystemInformationLayoutGrade:         "Hypothesis",
		B524SystemInformationQualificationGrade:  "Unknown",
		B524SystemInformationPresenceGrade:       "Unknown",
		B524SystemInformationEligibilityGrade:    "Unknown",
	}
	switch b524SystemInformationRequestContext.B524SystemInformationIdentifier {
	case 0:
		b524SystemInformationResult.B524SystemInformationMeaning = "supported circuit capacity"
	case 1:
		b524SystemInformationResult.B524SystemInformationMeaning = "supported zone capacity"
	case 0x000d:
		b524SystemInformationResult.B524SystemInformationMeaning = "heat-pump count"
	}
	if len(b524SystemInformationRaw) != 4 {
		b524SystemInformationResult.B524SystemInformationInvalidReason = "ab42 requires exactly four normalized bytes"
		return b524SystemInformationResult
	}
	b524SystemInformationValue := b524SystemInformationMath.Float32frombits(b524SystemInformationBinary.LittleEndian.Uint32(b524SystemInformationRaw))
	if b524SystemInformationMath.IsNaN(float64(b524SystemInformationValue)) || b524SystemInformationMath.IsInf(float64(b524SystemInformationValue), 0) {
		b524SystemInformationResult.B524SystemInformationInvalidReason = "ab42 non-finite binary32"
		return b524SystemInformationResult
	}
	b524SystemInformationResult.B524SystemInformationValid = true
	b524SystemInformationResult.B524SystemInformationValue = b524SystemInformationValue
	b524SystemInformationResult.B524SystemInformationRepresentationGrade = "Proven"
	if b524SystemInformationValue >= 0 && b524SystemInformationMath.Trunc(float64(b524SystemInformationValue)) == float64(b524SystemInformationValue) {
		switch b524SystemInformationRequestContext.B524SystemInformationIdentifier {
		case 0, 1:
			b524SystemInformationResult.B524SystemInformationCapacity = b524SystemInformationCopy(&b524SystemInformationValue)
		case 0x000d:
			b524SystemInformationResult.B524SystemInformationHeatPumpCount = b524SystemInformationCopy(&b524SystemInformationValue)
		}
	}
	return b524SystemInformationResult
}
