package system

import (
	b524CircuitScalarBinary "encoding/binary"
	b524CircuitScalarMath "math"
)

// B524CircuitScalarField is an optional value; Present distinguishes an absent
// value from a zero value.
type B524CircuitScalarField[B524CircuitScalarT any] struct {
	B524CircuitScalarPresent bool               `json:"present"`
	B524CircuitScalarValue   B524CircuitScalarT `json:"value"`
}

// B524CircuitScalarGrade is the qualification grade of a circuit scalar.
type B524CircuitScalarGrade string

// Circuit scalar grades, from unqualified to proven.
const (
	B524CircuitScalarUnknown    B524CircuitScalarGrade = "Unknown"
	B524CircuitScalarHypothesis B524CircuitScalarGrade = "Hypothesis"
	B524CircuitScalarProven     B524CircuitScalarGrade = "Proven"
)

// B524CircuitScalarContext is the controller identity (family, software,
// hardware, API version and revision) and the opcode, group, instance and
// register of a circuit scalar read; every field must be present.
type B524CircuitScalarContext struct {
	B524CircuitScalarFamily      B524CircuitScalarField[string] `json:"family"`
	B524CircuitScalarSoftware    B524CircuitScalarField[string] `json:"software"`
	B524CircuitScalarHardware    B524CircuitScalarField[string] `json:"hardware"`
	B524CircuitScalarAPIVersion  B524CircuitScalarField[string] `json:"api_version"`
	B524CircuitScalarAPIRevision B524CircuitScalarField[string] `json:"api_revision"`
	B524CircuitScalarOpcode      B524CircuitScalarField[byte]   `json:"opcode"`
	B524CircuitScalarGroup       B524CircuitScalarField[byte]   `json:"group"`
	B524CircuitScalarInstance    B524CircuitScalarField[byte]   `json:"instance"`
	B524CircuitScalarRegister    B524CircuitScalarField[uint16] `json:"register"`
}

// B524CircuitScalarInput is one observed circuit scalar read: its context,
// the reply kind, the returned group and register, the caller's grade of
// the codec and the raw reply value.
type B524CircuitScalarInput struct {
	B524CircuitScalarContext          B524CircuitScalarContext                       `json:"context"`
	B524CircuitScalarReplyKind        B524CircuitScalarField[string]                 `json:"reply_kind"`
	B524CircuitScalarReturnedGroup    B524CircuitScalarField[byte]                   `json:"returned_group"`
	B524CircuitScalarReturnedRegister B524CircuitScalarField[uint16]                 `json:"returned_register"`
	B524CircuitScalarCodecGrade       B524CircuitScalarField[B524CircuitScalarGrade] `json:"codec_grade"`
	B524CircuitScalarRaw              []byte                                         `json:"raw"`
}

// B524CircuitScalarState says whether a circuit scalar was decoded.
type B524CircuitScalarState string

// Circuit scalar states: raw-only keeps just the input, candidate adds a
// decoded value.
const (
	B524CircuitScalarRawOnly   B524CircuitScalarState = "raw_only"
	B524CircuitScalarCandidate B524CircuitScalarState = "candidate"
)

// B524CircuitScalarResult is the graded interpretation of a circuit scalar:
// the copied input, the state, the grade, the declared unit ("Unknown" when
// none is declared), the codec and the value, which is present only for a
// candidate.
type B524CircuitScalarResult struct {
	B524CircuitScalarInput        B524CircuitScalarInput          `json:"input"`
	B524CircuitScalarState        B524CircuitScalarState          `json:"state"`
	B524CircuitScalarGrade        B524CircuitScalarGrade          `json:"grade"`
	B524CircuitScalarDeclaredUnit string                          `json:"declared_unit"`
	B524CircuitScalarCodec        string                          `json:"codec"`
	B524CircuitScalarValue        B524CircuitScalarField[float32] `json:"value"`
}

// B524InterpretCircuitScalar decodes a known circuit register (group 0x02,
// read opcode) as a little-endian float32 for a supported controller
// family. It yields a Hypothesis candidate only when the identity is
// complete, the codec grade is Hypothesis or Proven, the reply is an echo
// of the requested group and register with exactly four bytes, and the
// value is finite; otherwise the result stays raw-only and Unknown. It
// never fails and copies the raw bytes.
func B524InterpretCircuitScalar(b524CircuitScalarInput B524CircuitScalarInput) B524CircuitScalarResult {
	if b524CircuitScalarInput.B524CircuitScalarRaw != nil {
		b524CircuitScalarRaw := make([]byte, len(b524CircuitScalarInput.B524CircuitScalarRaw))
		copy(b524CircuitScalarRaw, b524CircuitScalarInput.B524CircuitScalarRaw)
		b524CircuitScalarInput.B524CircuitScalarRaw = b524CircuitScalarRaw
	}
	b524CircuitScalarResult := B524CircuitScalarResult{B524CircuitScalarInput: b524CircuitScalarInput, B524CircuitScalarState: B524CircuitScalarRawOnly, B524CircuitScalarGrade: B524CircuitScalarUnknown, B524CircuitScalarDeclaredUnit: "Unknown"}
	b524CircuitScalarContext := b524CircuitScalarInput.B524CircuitScalarContext
	if !b524CircuitScalarContext.B524CircuitScalarFamily.B524CircuitScalarPresent || !b524CircuitScalarContext.B524CircuitScalarSoftware.B524CircuitScalarPresent || !b524CircuitScalarContext.B524CircuitScalarHardware.B524CircuitScalarPresent || !b524CircuitScalarContext.B524CircuitScalarAPIVersion.B524CircuitScalarPresent || !b524CircuitScalarContext.B524CircuitScalarAPIRevision.B524CircuitScalarPresent || !b524CircuitScalarContext.B524CircuitScalarOpcode.B524CircuitScalarPresent || !b524CircuitScalarContext.B524CircuitScalarGroup.B524CircuitScalarPresent || !b524CircuitScalarContext.B524CircuitScalarInstance.B524CircuitScalarPresent || !b524CircuitScalarContext.B524CircuitScalarRegister.B524CircuitScalarPresent {
		return b524CircuitScalarResult
	}
	if b524CircuitScalarContext.B524CircuitScalarOpcode.B524CircuitScalarValue != 0x02 || b524CircuitScalarContext.B524CircuitScalarGroup.B524CircuitScalarValue != 0x02 {
		return b524CircuitScalarResult
	}
	b524CircuitScalarFamily := b524CircuitScalarContext.B524CircuitScalarFamily.B524CircuitScalarValue
	if b524CircuitScalarFamily != "controller-720" && b524CircuitScalarFamily != "controller-700" {
		return b524CircuitScalarResult
	}
	switch b524CircuitScalarContext.B524CircuitScalarRegister.B524CircuitScalarValue {
	case 0x0007, 0x0008, 0x0010, 0x0011, 0x0012, 0x0014, 0x001f:
		b524CircuitScalarResult.B524CircuitScalarDeclaredUnit = "°C"
	case 0x000a, 0x000b:
		b524CircuitScalarResult.B524CircuitScalarDeclaredUnit = "K"
	case 0x000f:
	case 0x001c:
		if b524CircuitScalarFamily != "controller-720" {
			return b524CircuitScalarResult
		}
	case 0x001a:
		if b524CircuitScalarFamily != "controller-720" {
			return b524CircuitScalarResult
		}
		b524CircuitScalarResult.B524CircuitScalarDeclaredUnit = "%"
	default:
		return b524CircuitScalarResult
	}
	b524CircuitScalarResult.B524CircuitScalarCodec = "float32_le"
	if !b524CircuitScalarInput.B524CircuitScalarCodecGrade.B524CircuitScalarPresent || (b524CircuitScalarInput.B524CircuitScalarCodecGrade.B524CircuitScalarValue != B524CircuitScalarHypothesis && b524CircuitScalarInput.B524CircuitScalarCodecGrade.B524CircuitScalarValue != B524CircuitScalarProven) || !b524CircuitScalarInput.B524CircuitScalarReplyKind.B524CircuitScalarPresent || b524CircuitScalarInput.B524CircuitScalarReplyKind.B524CircuitScalarValue != "echo_payload" || !b524CircuitScalarInput.B524CircuitScalarReturnedGroup.B524CircuitScalarPresent || b524CircuitScalarInput.B524CircuitScalarReturnedGroup.B524CircuitScalarValue != b524CircuitScalarContext.B524CircuitScalarGroup.B524CircuitScalarValue || !b524CircuitScalarInput.B524CircuitScalarReturnedRegister.B524CircuitScalarPresent || b524CircuitScalarInput.B524CircuitScalarReturnedRegister.B524CircuitScalarValue != b524CircuitScalarContext.B524CircuitScalarRegister.B524CircuitScalarValue || len(b524CircuitScalarInput.B524CircuitScalarRaw) != 4 {
		return b524CircuitScalarResult
	}
	b524CircuitScalarValue := b524CircuitScalarMath.Float32frombits(b524CircuitScalarBinary.LittleEndian.Uint32(b524CircuitScalarInput.B524CircuitScalarRaw))
	if b524CircuitScalarMath.IsNaN(float64(b524CircuitScalarValue)) || b524CircuitScalarMath.IsInf(float64(b524CircuitScalarValue), 0) {
		return b524CircuitScalarResult
	}
	b524CircuitScalarResult.B524CircuitScalarState = B524CircuitScalarCandidate
	b524CircuitScalarResult.B524CircuitScalarGrade = B524CircuitScalarHypothesis
	b524CircuitScalarResult.B524CircuitScalarValue = B524CircuitScalarField[float32]{B524CircuitScalarPresent: true, B524CircuitScalarValue: b524CircuitScalarValue}
	return b524CircuitScalarResult
}
