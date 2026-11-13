package system

import (
	b524CircuitScalarBinary "encoding/binary"
	b524CircuitScalarMath "math"
)

type B524CircuitScalarField[B524CircuitScalarT any] struct {
	B524CircuitScalarPresent bool               `json:"Ab62Present"`
	B524CircuitScalarValue   B524CircuitScalarT `json:"Ab62Value"`
}

type B524CircuitScalarGrade string

const (
	B524CircuitScalarUnknown    B524CircuitScalarGrade = "Unknown"
	B524CircuitScalarHypothesis B524CircuitScalarGrade = "Hypothesis"
	B524CircuitScalarProven     B524CircuitScalarGrade = "Proven"
)

type B524CircuitScalarContext struct {
	B524CircuitScalarFamily      B524CircuitScalarField[string] `json:"Ab62Family"`
	B524CircuitScalarSoftware    B524CircuitScalarField[string] `json:"Ab62Software"`
	B524CircuitScalarHardware    B524CircuitScalarField[string] `json:"Ab62Hardware"`
	B524CircuitScalarAPIVersion  B524CircuitScalarField[string] `json:"Ab62APIVersion"`
	B524CircuitScalarAPIRevision B524CircuitScalarField[string] `json:"Ab62APIRevision"`
	B524CircuitScalarOpcode      B524CircuitScalarField[byte]   `json:"Ab62Opcode"`
	B524CircuitScalarGroup       B524CircuitScalarField[byte]   `json:"Ab62Group"`
	B524CircuitScalarInstance    B524CircuitScalarField[byte]   `json:"Ab62Instance"`
	B524CircuitScalarRegister    B524CircuitScalarField[uint16] `json:"Ab62Register"`
}

type B524CircuitScalarInput struct {
	B524CircuitScalarContext          B524CircuitScalarContext                       `json:"Ab62Context"`
	B524CircuitScalarReplyKind        B524CircuitScalarField[string]                 `json:"Ab62ReplyKind"`
	B524CircuitScalarReturnedGroup    B524CircuitScalarField[byte]                   `json:"Ab62ReturnedGroup"`
	B524CircuitScalarReturnedRegister B524CircuitScalarField[uint16]                 `json:"Ab62ReturnedRegister"`
	B524CircuitScalarCodecGrade       B524CircuitScalarField[B524CircuitScalarGrade] `json:"Ab62CodecGrade"`
	B524CircuitScalarRaw              []byte                                         `json:"Ab62Raw"`
}

type B524CircuitScalarState string

const (
	B524CircuitScalarRawOnly   B524CircuitScalarState = "raw_only"
	B524CircuitScalarCandidate B524CircuitScalarState = "candidate"
)

type B524CircuitScalarResult struct {
	B524CircuitScalarInput        B524CircuitScalarInput          `json:"Ab62Input"`
	B524CircuitScalarState        B524CircuitScalarState          `json:"Ab62State"`
	B524CircuitScalarGrade        B524CircuitScalarGrade          `json:"Ab62Grade"`
	B524CircuitScalarDeclaredUnit string                          `json:"Ab62DeclaredUnit"`
	B524CircuitScalarCodec        string                          `json:"Ab62Codec"`
	B524CircuitScalarValue        B524CircuitScalarField[float32] `json:"Ab62Value"`
}

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
