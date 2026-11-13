package system

import (
	b524ZoneScalarBinary "encoding/binary"
	b524ZoneScalarMath "math"
)

type B524ZoneScalarField[B524ZoneScalarT any] struct {
	B524ZoneScalarPresent bool            `json:"Ab64Present"`
	B524ZoneScalarValue   B524ZoneScalarT `json:"Ab64Value"`
}

type B524ZoneScalarGrade string

const (
	B524ZoneScalarUnknown    B524ZoneScalarGrade = "Unknown"
	B524ZoneScalarHypothesis B524ZoneScalarGrade = "Hypothesis"
	B524ZoneScalarProven     B524ZoneScalarGrade = "Proven"
)

type B524ZoneScalarContext struct {
	B524ZoneScalarFamily      B524ZoneScalarField[string] `json:"Ab64Family"`
	B524ZoneScalarSoftware    B524ZoneScalarField[string] `json:"Ab64Software"`
	B524ZoneScalarHardware    B524ZoneScalarField[string] `json:"Ab64Hardware"`
	B524ZoneScalarAPIVersion  B524ZoneScalarField[string] `json:"Ab64APIVersion"`
	B524ZoneScalarAPIRevision B524ZoneScalarField[string] `json:"Ab64APIRevision"`
	B524ZoneScalarOpcode      B524ZoneScalarField[byte]   `json:"Ab64Opcode"`
	B524ZoneScalarGroup       B524ZoneScalarField[byte]   `json:"Ab64Group"`
	B524ZoneScalarInstance    B524ZoneScalarField[byte]   `json:"Ab64Instance"`
	B524ZoneScalarRegister    B524ZoneScalarField[uint16] `json:"Ab64Register"`
}

type B524ZoneScalarInput struct {
	B524ZoneScalarContext          B524ZoneScalarContext       `json:"Ab64Context"`
	B524ZoneScalarReplyKind        string                      `json:"Ab64ReplyKind"`
	B524ZoneScalarReturnedGroup    B524ZoneScalarField[byte]   `json:"Ab64ReturnedGroup"`
	B524ZoneScalarReturnedRegister B524ZoneScalarField[uint16] `json:"Ab64ReturnedRegister"`
	B524ZoneScalarCodecGrade       B524ZoneScalarGrade         `json:"Ab64CodecGrade"`
	B524ZoneScalarRaw              []byte                      `json:"Ab64Raw"`
}

type B524ZoneScalarState string

const (
	B524ZoneScalarRawOnly   B524ZoneScalarState = "raw_only"
	B524ZoneScalarCandidate B524ZoneScalarState = "candidate"
)

type B524ZoneScalarResult struct {
	B524ZoneScalarInput         B524ZoneScalarInput          `json:"Ab64Input"`
	B524ZoneScalarState         B524ZoneScalarState          `json:"Ab64State"`
	B524ZoneScalarValue         B524ZoneScalarField[float32] `json:"Ab64Value"`
	B524ZoneScalarGrade         B524ZoneScalarGrade          `json:"Ab64Grade"`
	B524ZoneScalarUnit          string                       `json:"Ab64Unit"`
	B524ZoneScalarCodec         string                       `json:"Ab64Codec"`
	B524ZoneScalarInstanceRule  string                       `json:"Ab64InstanceRule"`
	B524ZoneScalarNativeSupport B524ZoneScalarGrade          `json:"Ab64NativeSupport"`
	B524ZoneScalarQualification B524ZoneScalarGrade          `json:"Ab64Qualification"`
	B524ZoneScalarPresence      B524ZoneScalarGrade          `json:"Ab64Presence"`
	B524ZoneScalarFreshness     B524ZoneScalarGrade          `json:"Ab64Freshness"`
	B524ZoneScalarExecution     B524ZoneScalarGrade          `json:"Ab64Execution"`
	B524ZoneScalarEffects       B524ZoneScalarGrade          `json:"Ab64Effects"`
	B524ZoneScalarPersistence   B524ZoneScalarGrade          `json:"Ab64Persistence"`
	B524ZoneScalarDelivery      B524ZoneScalarGrade          `json:"Ab64Delivery"`
}

func B524ZoneScalarInterpret(b524ZoneScalarInput B524ZoneScalarInput) B524ZoneScalarResult {
	b524ZoneScalarInput.B524ZoneScalarRaw = append([]byte(nil), b524ZoneScalarInput.B524ZoneScalarRaw...)
	b524ZoneScalarResult := B524ZoneScalarResult{
		B524ZoneScalarInput: b524ZoneScalarInput, B524ZoneScalarState: B524ZoneScalarRawOnly, B524ZoneScalarGrade: B524ZoneScalarUnknown,
		B524ZoneScalarUnit: "Unknown", B524ZoneScalarNativeSupport: B524ZoneScalarUnknown, B524ZoneScalarQualification: B524ZoneScalarUnknown,
		B524ZoneScalarPresence: B524ZoneScalarUnknown, B524ZoneScalarFreshness: B524ZoneScalarUnknown, B524ZoneScalarExecution: B524ZoneScalarUnknown,
		B524ZoneScalarEffects: B524ZoneScalarUnknown, B524ZoneScalarPersistence: B524ZoneScalarUnknown, B524ZoneScalarDelivery: B524ZoneScalarUnknown,
	}
	b524ZoneScalarContext := b524ZoneScalarInput.B524ZoneScalarContext
	if !b524ZoneScalarContext.B524ZoneScalarFamily.B524ZoneScalarPresent || !b524ZoneScalarContext.B524ZoneScalarSoftware.B524ZoneScalarPresent ||
		!b524ZoneScalarContext.B524ZoneScalarHardware.B524ZoneScalarPresent || !b524ZoneScalarContext.B524ZoneScalarAPIVersion.B524ZoneScalarPresent ||
		!b524ZoneScalarContext.B524ZoneScalarAPIRevision.B524ZoneScalarPresent || !b524ZoneScalarContext.B524ZoneScalarOpcode.B524ZoneScalarPresent ||
		!b524ZoneScalarContext.B524ZoneScalarGroup.B524ZoneScalarPresent || !b524ZoneScalarContext.B524ZoneScalarInstance.B524ZoneScalarPresent ||
		!b524ZoneScalarContext.B524ZoneScalarRegister.B524ZoneScalarPresent {
		return b524ZoneScalarResult
	}
	if b524ZoneScalarContext.B524ZoneScalarFamily.B524ZoneScalarValue != "controller-720" && b524ZoneScalarContext.B524ZoneScalarFamily.B524ZoneScalarValue != "controller-700" {
		return b524ZoneScalarResult
	}
	if b524ZoneScalarContext.B524ZoneScalarOpcode.B524ZoneScalarValue != 0x02 || b524ZoneScalarContext.B524ZoneScalarGroup.B524ZoneScalarValue != 0x03 {
		return b524ZoneScalarResult
	}
	switch b524ZoneScalarContext.B524ZoneScalarRegister.B524ZoneScalarValue {
	case 0x0002, 0x0005, 0x0009, 0x000f, 0x0014, 0x0015:
		b524ZoneScalarResult.B524ZoneScalarUnit = "°C"
	case 0x0023:
		if b524ZoneScalarContext.B524ZoneScalarFamily.B524ZoneScalarValue != "controller-720" {
			return b524ZoneScalarResult
		}
		b524ZoneScalarResult.B524ZoneScalarUnit = "°C"
	case 0x0028:
		b524ZoneScalarResult.B524ZoneScalarUnit = "%"
	default:
		return b524ZoneScalarResult
	}
	b524ZoneScalarResult.B524ZoneScalarCodec = "float32_le"
	b524ZoneScalarResult.B524ZoneScalarInstanceRule = "dynamic"
	if (b524ZoneScalarInput.B524ZoneScalarCodecGrade != B524ZoneScalarHypothesis && b524ZoneScalarInput.B524ZoneScalarCodecGrade != B524ZoneScalarProven) ||
		b524ZoneScalarInput.B524ZoneScalarReplyKind != "echo_payload" || !b524ZoneScalarInput.B524ZoneScalarReturnedGroup.B524ZoneScalarPresent ||
		!b524ZoneScalarInput.B524ZoneScalarReturnedRegister.B524ZoneScalarPresent ||
		b524ZoneScalarInput.B524ZoneScalarReturnedGroup.B524ZoneScalarValue != b524ZoneScalarContext.B524ZoneScalarGroup.B524ZoneScalarValue ||
		b524ZoneScalarInput.B524ZoneScalarReturnedRegister.B524ZoneScalarValue != b524ZoneScalarContext.B524ZoneScalarRegister.B524ZoneScalarValue || len(b524ZoneScalarInput.B524ZoneScalarRaw) != 4 {
		return b524ZoneScalarResult
	}
	b524ZoneScalarValue := b524ZoneScalarMath.Float32frombits(b524ZoneScalarBinary.LittleEndian.Uint32(b524ZoneScalarInput.B524ZoneScalarRaw))
	if b524ZoneScalarMath.IsNaN(float64(b524ZoneScalarValue)) || b524ZoneScalarMath.IsInf(float64(b524ZoneScalarValue), 0) {
		return b524ZoneScalarResult
	}
	b524ZoneScalarResult.B524ZoneScalarValue = B524ZoneScalarField[float32]{B524ZoneScalarPresent: true, B524ZoneScalarValue: b524ZoneScalarValue}
	b524ZoneScalarResult.B524ZoneScalarState = B524ZoneScalarCandidate
	b524ZoneScalarResult.B524ZoneScalarGrade = B524ZoneScalarHypothesis
	return b524ZoneScalarResult
}
