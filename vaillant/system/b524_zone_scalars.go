package system

import (
	b524ZoneScalarBinary "encoding/binary"
	b524ZoneScalarMath "math"
)

// B524ZoneScalarField is an optional value; Present distinguishes an absent
// value from a zero value.
type B524ZoneScalarField[B524ZoneScalarT any] struct {
	B524ZoneScalarPresent bool            `json:"present"`
	B524ZoneScalarValue   B524ZoneScalarT `json:"value"`
}

// B524ZoneScalarGrade is the qualification grade of a zone scalar aspect.
type B524ZoneScalarGrade string

// Zone scalar grades, from unqualified to proven.
const (
	B524ZoneScalarUnknown    B524ZoneScalarGrade = "Unknown"
	B524ZoneScalarHypothesis B524ZoneScalarGrade = "Hypothesis"
	B524ZoneScalarProven     B524ZoneScalarGrade = "Proven"
)

// B524ZoneScalarContext is the controller identity and the opcode, group,
// instance and register of a zone scalar read; every field must be
// present.
type B524ZoneScalarContext struct {
	B524ZoneScalarFamily      B524ZoneScalarField[string] `json:"family"`
	B524ZoneScalarSoftware    B524ZoneScalarField[string] `json:"software"`
	B524ZoneScalarHardware    B524ZoneScalarField[string] `json:"hardware"`
	B524ZoneScalarAPIVersion  B524ZoneScalarField[string] `json:"api_version"`
	B524ZoneScalarAPIRevision B524ZoneScalarField[string] `json:"api_revision"`
	B524ZoneScalarOpcode      B524ZoneScalarField[byte]   `json:"opcode"`
	B524ZoneScalarGroup       B524ZoneScalarField[byte]   `json:"group"`
	B524ZoneScalarInstance    B524ZoneScalarField[byte]   `json:"instance"`
	B524ZoneScalarRegister    B524ZoneScalarField[uint16] `json:"register"`
}

// B524ZoneScalarInput is one observed zone scalar read: its context, the
// reply kind, the returned group and register, the caller's codec grade and
// the raw reply value.
type B524ZoneScalarInput struct {
	B524ZoneScalarContext          B524ZoneScalarContext       `json:"context"`
	B524ZoneScalarReplyKind        string                      `json:"reply_kind"`
	B524ZoneScalarReturnedGroup    B524ZoneScalarField[byte]   `json:"returned_group"`
	B524ZoneScalarReturnedRegister B524ZoneScalarField[uint16] `json:"returned_register"`
	B524ZoneScalarCodecGrade       B524ZoneScalarGrade         `json:"codec_grade"`
	B524ZoneScalarRaw              []byte                      `json:"raw"`
}

// B524ZoneScalarState says whether a zone scalar was decoded.
type B524ZoneScalarState string

// Zone scalar states: raw-only keeps just the input, candidate adds a
// decoded value.
const (
	B524ZoneScalarRawOnly   B524ZoneScalarState = "raw_only"
	B524ZoneScalarCandidate B524ZoneScalarState = "candidate"
)

// B524ZoneScalarResult is the graded interpretation of a zone scalar. The
// grade is at most Hypothesis, the unit is degrees Celsius, percent or
// "Unknown", and native support, qualification, presence, freshness,
// execution, effects, persistence and delivery stay Unknown. Value is
// present only for a candidate.
type B524ZoneScalarResult struct {
	B524ZoneScalarInput         B524ZoneScalarInput          `json:"input"`
	B524ZoneScalarState         B524ZoneScalarState          `json:"state"`
	B524ZoneScalarValue         B524ZoneScalarField[float32] `json:"value"`
	B524ZoneScalarGrade         B524ZoneScalarGrade          `json:"grade"`
	B524ZoneScalarUnit          string                       `json:"unit"`
	B524ZoneScalarCodec         string                       `json:"codec"`
	B524ZoneScalarInstanceRule  string                       `json:"instance_rule"`
	B524ZoneScalarNativeSupport B524ZoneScalarGrade          `json:"native_support"`
	B524ZoneScalarQualification B524ZoneScalarGrade          `json:"qualification"`
	B524ZoneScalarPresence      B524ZoneScalarGrade          `json:"presence"`
	B524ZoneScalarFreshness     B524ZoneScalarGrade          `json:"freshness"`
	B524ZoneScalarExecution     B524ZoneScalarGrade          `json:"execution"`
	B524ZoneScalarEffects       B524ZoneScalarGrade          `json:"effects"`
	B524ZoneScalarPersistence   B524ZoneScalarGrade          `json:"persistence"`
	B524ZoneScalarDelivery      B524ZoneScalarGrade          `json:"delivery"`
}

// B524ZoneScalarInterpret decodes a known zone register (group 0x03, local
// opcode, any instance) as a little-endian float32 for a supported
// controller family. It yields a Hypothesis candidate only for a complete
// identity, a codec grade of Hypothesis or Proven, an echo of the requested
// group and register with exactly four bytes and a finite value; otherwise
// the result stays raw-only. It never fails and copies the raw bytes.
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
