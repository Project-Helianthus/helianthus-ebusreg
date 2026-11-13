package system

import (
	b524SystemScalarBinary "encoding/binary"
	b524SystemScalarMath "math"
)

type B524SystemScalarField[B524SystemScalarT any] struct {
	B524SystemScalarPresent bool              `json:"Ab58Present"`
	B524SystemScalarValue   B524SystemScalarT `json:"Ab58Value"`
}

type B524SystemScalarContext struct {
	B524SystemScalarFamily      B524SystemScalarField[string] `json:"Ab58Family"`
	B524SystemScalarSoftware    B524SystemScalarField[string] `json:"Ab58Software"`
	B524SystemScalarHardware    B524SystemScalarField[string] `json:"Ab58Hardware"`
	B524SystemScalarAPIVersion  B524SystemScalarField[string] `json:"Ab58APIVersion"`
	B524SystemScalarAPIRevision B524SystemScalarField[string] `json:"Ab58APIRevision"`
	B524SystemScalarOpcode      B524SystemScalarField[byte]   `json:"Ab58Opcode"`
	B524SystemScalarGroup       B524SystemScalarField[byte]   `json:"Ab58Group"`
	B524SystemScalarInstance    B524SystemScalarField[byte]   `json:"Ab58Instance"`
	B524SystemScalarRegister    B524SystemScalarField[uint16] `json:"Ab58Register"`
}

type B524SystemScalarInput struct {
	B524SystemScalarContext          B524SystemScalarContext       `json:"Ab58Context"`
	B524SystemScalarReplyKind        B524SystemScalarField[string] `json:"Ab58ReplyKind"`
	B524SystemScalarReturnedGroup    B524SystemScalarField[byte]   `json:"Ab58ReturnedGroup"`
	B524SystemScalarReturnedRegister B524SystemScalarField[uint16] `json:"Ab58ReturnedRegister"`
	B524SystemScalarCodecGrade       B524SystemScalarField[string] `json:"Ab58CodecGrade"`
	B524SystemScalarRaw              []byte                        `json:"Ab58Raw"`
}

type B524SystemScalarState string

const (
	B524SystemScalarRawOnly   B524SystemScalarState = "raw_only"
	B524SystemScalarCandidate B524SystemScalarState = "candidate"
)

type B524SystemScalar struct {
	B524SystemScalarFloat32 float32 `json:"Ab58Float32"`
	B524SystemScalarUint16  uint16  `json:"Ab58Uint16"`
}

type B524SystemScalarResult struct {
	B524SystemScalarInput        B524SystemScalarInput `json:"Ab58Input"`
	B524SystemScalarState        B524SystemScalarState `json:"Ab58State"`
	B524SystemScalarCodec        string                `json:"Ab58Codec"`
	B524SystemScalarDeclaredUnit string                `json:"Ab58DeclaredUnit"`
	B524SystemScalarGrade        string                `json:"Ab58Grade"`
	B524SystemScalarPresence     string                `json:"Ab58Presence"`
	B524SystemScalarFreshness    string                `json:"Ab58Freshness"`
	B524SystemScalarExecution    string                `json:"Ab58Execution"`
	B524SystemScalarSupport      string                `json:"Ab58Support"`
	B524SystemScalarEffects      string                `json:"Ab58Effects"`
	B524SystemScalarPersistence  string                `json:"Ab58Persistence"`
	B524SystemScalarDelivery     string                `json:"Ab58Delivery"`
	B524SystemScalarValue        *B524SystemScalar     `json:"Ab58Value"`
}

func B524SystemInterpretScalars(b524SystemScalarInput B524SystemScalarInput) B524SystemScalarResult {
	b524SystemScalarResult := B524SystemScalarResult{
		B524SystemScalarInput: b524SystemScalarInput, B524SystemScalarState: B524SystemScalarRawOnly,
		B524SystemScalarDeclaredUnit: "Unknown", B524SystemScalarGrade: "Unknown",
		B524SystemScalarPresence: "Unknown", B524SystemScalarFreshness: "Unknown", B524SystemScalarExecution: "Unknown",
		B524SystemScalarSupport: "Unknown", B524SystemScalarEffects: "Unknown",
		B524SystemScalarPersistence: "Unknown", B524SystemScalarDelivery: "Unknown",
	}
	if b524SystemScalarInput.B524SystemScalarRaw != nil {
		b524SystemScalarResult.B524SystemScalarInput.B524SystemScalarRaw = append([]byte{}, b524SystemScalarInput.B524SystemScalarRaw...)
	}
	b524SystemScalarContext := b524SystemScalarInput.B524SystemScalarContext
	if !b524SystemScalarContext.B524SystemScalarFamily.B524SystemScalarPresent || !b524SystemScalarContext.B524SystemScalarSoftware.B524SystemScalarPresent ||
		!b524SystemScalarContext.B524SystemScalarHardware.B524SystemScalarPresent || !b524SystemScalarContext.B524SystemScalarAPIVersion.B524SystemScalarPresent ||
		!b524SystemScalarContext.B524SystemScalarAPIRevision.B524SystemScalarPresent || !b524SystemScalarContext.B524SystemScalarOpcode.B524SystemScalarPresent ||
		!b524SystemScalarContext.B524SystemScalarGroup.B524SystemScalarPresent || !b524SystemScalarContext.B524SystemScalarInstance.B524SystemScalarPresent ||
		!b524SystemScalarContext.B524SystemScalarRegister.B524SystemScalarPresent {
		return b524SystemScalarResult
	}
	if b524SystemScalarContext.B524SystemScalarOpcode.B524SystemScalarValue != 0x02 || b524SystemScalarContext.B524SystemScalarGroup.B524SystemScalarValue != 0 {
		return b524SystemScalarResult
	}
	b524SystemScalarDynamic := false
	switch b524SystemScalarContext.B524SystemScalarFamily.B524SystemScalarValue {
	case "controller-720":
		switch b524SystemScalarContext.B524SystemScalarRegister.B524SystemScalarValue {
		case 0x0002, 0x0022, 0x0023, 0x0026, 0x0038, 0x0046, 0x004B, 0x0067, 0x0068, 0x0073, 0x0095, 0x009D, 0x009E, 0x009F, 0x00A0, 0x00EF:
			b524SystemScalarResult.B524SystemScalarCodec, b524SystemScalarResult.B524SystemScalarDeclaredUnit = "float32_le", "°C"
		case 0x0004, 0x0017, 0x0018:
			b524SystemScalarResult.B524SystemScalarCodec, b524SystemScalarResult.B524SystemScalarDeclaredUnit = "uint16_le", "min"
		case 0x0029:
			b524SystemScalarResult.B524SystemScalarCodec, b524SystemScalarResult.B524SystemScalarDeclaredUnit = "float32_le", "K"
		case 0x002F, 0x0036:
			b524SystemScalarResult.B524SystemScalarCodec = "uint16_le"
		case 0x0039:
			b524SystemScalarResult.B524SystemScalarCodec, b524SystemScalarResult.B524SystemScalarDeclaredUnit = "float32_le", "bar"
		}
	case "controller-700":
		switch b524SystemScalarContext.B524SystemScalarRegister.B524SystemScalarValue {
		case 0x0002, 0x0022, 0x0023, 0x0026, 0x0038, 0x0046, 0x004B, 0x0067, 0x0068, 0x0073, 0x0095, 0x009D, 0x009E, 0x009F, 0x00A0:
			b524SystemScalarResult.B524SystemScalarCodec, b524SystemScalarResult.B524SystemScalarDeclaredUnit = "float32_le", "°C"
		case 0x0006, 0x002F, 0x0036:
			b524SystemScalarResult.B524SystemScalarCodec = "uint16_le"
		case 0x0017:
			b524SystemScalarResult.B524SystemScalarCodec, b524SystemScalarResult.B524SystemScalarDeclaredUnit = "uint16_le", "min"
		case 0x0029:
			b524SystemScalarResult.B524SystemScalarCodec, b524SystemScalarResult.B524SystemScalarDeclaredUnit = "float32_le", "K"
		case 0x0039:
			b524SystemScalarResult.B524SystemScalarCodec, b524SystemScalarResult.B524SystemScalarDeclaredUnit = "float32_le", "bar"
		case 0x0042, 0x0043:
			b524SystemScalarResult.B524SystemScalarCodec = "uint16_le"
			b524SystemScalarDynamic = true
		}
	}
	if !b524SystemScalarDynamic && b524SystemScalarContext.B524SystemScalarInstance.B524SystemScalarValue != 0 {
		return b524SystemScalarResult
	}
	if b524SystemScalarResult.B524SystemScalarCodec == "" || !b524SystemScalarInput.B524SystemScalarReplyKind.B524SystemScalarPresent ||
		b524SystemScalarInput.B524SystemScalarReplyKind.B524SystemScalarValue != "echo_payload" ||
		!b524SystemScalarInput.B524SystemScalarReturnedGroup.B524SystemScalarPresent || !b524SystemScalarInput.B524SystemScalarReturnedRegister.B524SystemScalarPresent ||
		b524SystemScalarInput.B524SystemScalarReturnedGroup.B524SystemScalarValue != b524SystemScalarContext.B524SystemScalarGroup.B524SystemScalarValue ||
		b524SystemScalarInput.B524SystemScalarReturnedRegister.B524SystemScalarValue != b524SystemScalarContext.B524SystemScalarRegister.B524SystemScalarValue ||
		!b524SystemScalarInput.B524SystemScalarCodecGrade.B524SystemScalarPresent ||
		(b524SystemScalarInput.B524SystemScalarCodecGrade.B524SystemScalarValue != "Hypothesis" && b524SystemScalarInput.B524SystemScalarCodecGrade.B524SystemScalarValue != "Proven") {
		return b524SystemScalarResult
	}
	b524SystemScalarScalar := &B524SystemScalar{}
	switch b524SystemScalarResult.B524SystemScalarCodec {
	case "float32_le":
		if len(b524SystemScalarInput.B524SystemScalarRaw) != 4 {
			return b524SystemScalarResult
		}
		b524SystemScalarScalar.B524SystemScalarFloat32 = b524SystemScalarMath.Float32frombits(b524SystemScalarBinary.LittleEndian.Uint32(b524SystemScalarInput.B524SystemScalarRaw))
		if b524SystemScalarMath.IsNaN(float64(b524SystemScalarScalar.B524SystemScalarFloat32)) || b524SystemScalarMath.IsInf(float64(b524SystemScalarScalar.B524SystemScalarFloat32), 0) {
			return b524SystemScalarResult
		}
	case "uint16_le":
		if len(b524SystemScalarInput.B524SystemScalarRaw) != 2 {
			return b524SystemScalarResult
		}
		b524SystemScalarScalar.B524SystemScalarUint16 = b524SystemScalarBinary.LittleEndian.Uint16(b524SystemScalarInput.B524SystemScalarRaw)
	}
	b524SystemScalarResult.B524SystemScalarState = B524SystemScalarCandidate
	b524SystemScalarResult.B524SystemScalarGrade = "Hypothesis"
	b524SystemScalarResult.B524SystemScalarValue = b524SystemScalarScalar
	return b524SystemScalarResult
}
