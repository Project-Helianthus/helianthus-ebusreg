package system

import (
	b524DHWVentilationScalarBinary "encoding/binary"
	b524DHWVentilationScalarMath "math"
)

type B524DHWVentilationScalarField[B524DHWVentilationScalarT any] struct {
	B524DHWVentilationScalarPresent bool                      `json:"Ab60Present"`
	B524DHWVentilationScalarValue   B524DHWVentilationScalarT `json:"Ab60Value"`
}

type B524DHWVentilationScalarContext struct {
	B524DHWVentilationScalarFamily      B524DHWVentilationScalarField[string] `json:"Ab60Family"`
	B524DHWVentilationScalarSoftware    B524DHWVentilationScalarField[string] `json:"Ab60Software"`
	B524DHWVentilationScalarHardware    B524DHWVentilationScalarField[string] `json:"Ab60Hardware"`
	B524DHWVentilationScalarAPIVersion  B524DHWVentilationScalarField[string] `json:"Ab60APIVersion"`
	B524DHWVentilationScalarAPIRevision B524DHWVentilationScalarField[string] `json:"Ab60APIRevision"`
	B524DHWVentilationScalarOpcode      B524DHWVentilationScalarField[byte]   `json:"Ab60Opcode"`
	B524DHWVentilationScalarGroup       B524DHWVentilationScalarField[byte]   `json:"Ab60Group"`
	B524DHWVentilationScalarInstance    B524DHWVentilationScalarField[byte]   `json:"Ab60Instance"`
	B524DHWVentilationScalarRegister    B524DHWVentilationScalarField[uint16] `json:"Ab60Register"`
}

type B524DHWVentilationScalarInput struct {
	B524DHWVentilationScalarContext          B524DHWVentilationScalarContext       `json:"Ab60Context"`
	B524DHWVentilationScalarReplyKind        B524DHWVentilationScalarField[string] `json:"Ab60ReplyKind"`
	B524DHWVentilationScalarReturnedGroup    B524DHWVentilationScalarField[byte]   `json:"Ab60ReturnedGroup"`
	B524DHWVentilationScalarReturnedRegister B524DHWVentilationScalarField[uint16] `json:"Ab60ReturnedRegister"`
	B524DHWVentilationScalarCodecGrade       B524DHWVentilationScalarField[string] `json:"Ab60CodecGrade"`
	B524DHWVentilationScalarRaw              []byte                                `json:"Ab60Raw"`
}

type B524DHWVentilationScalarState string

const (
	B524DHWVentilationScalarRawOnly   B524DHWVentilationScalarState = "raw_only"
	B524DHWVentilationScalarCandidate B524DHWVentilationScalarState = "candidate"
)

type B524DHWVentilationScalar struct {
	B524DHWVentilationScalarFloat32 float32 `json:"Ab60Float32"`
	B524DHWVentilationScalarUint16  uint16  `json:"Ab60Uint16"`
}

type B524DHWVentilationScalarResult struct {
	B524DHWVentilationScalarInput        B524DHWVentilationScalarInput `json:"Ab60Input"`
	B524DHWVentilationScalarState        B524DHWVentilationScalarState `json:"Ab60State"`
	B524DHWVentilationScalarCodec        string                        `json:"Ab60Codec"`
	B524DHWVentilationScalarDeclaredUnit string                        `json:"Ab60DeclaredUnit"`
	B524DHWVentilationScalarGrade        string                        `json:"Ab60Grade"`
	B524DHWVentilationScalarPresence     string                        `json:"Ab60Presence"`
	B524DHWVentilationScalarFreshness    string                        `json:"Ab60Freshness"`
	B524DHWVentilationScalarExecution    string                        `json:"Ab60Execution"`
	B524DHWVentilationScalarSupport      string                        `json:"Ab60Support"`
	B524DHWVentilationScalarEffects      string                        `json:"Ab60Effects"`
	B524DHWVentilationScalarPersistence  string                        `json:"Ab60Persistence"`
	B524DHWVentilationScalarDelivery     string                        `json:"Ab60Delivery"`
	B524DHWVentilationScalarValue        *B524DHWVentilationScalar     `json:"Ab60Value"`
}

func B524DHWVentilationInterpretScalars(b524DHWVentilationScalarInput B524DHWVentilationScalarInput) B524DHWVentilationScalarResult {
	b524DHWVentilationScalarResult := B524DHWVentilationScalarResult{
		B524DHWVentilationScalarInput: b524DHWVentilationScalarInput, B524DHWVentilationScalarState: B524DHWVentilationScalarRawOnly,
		B524DHWVentilationScalarDeclaredUnit: "Unknown", B524DHWVentilationScalarGrade: "Unknown",
		B524DHWVentilationScalarPresence: "Unknown", B524DHWVentilationScalarFreshness: "Unknown", B524DHWVentilationScalarExecution: "Unknown",
		B524DHWVentilationScalarSupport: "Unknown", B524DHWVentilationScalarEffects: "Unknown",
		B524DHWVentilationScalarPersistence: "Unknown", B524DHWVentilationScalarDelivery: "Unknown",
	}
	if b524DHWVentilationScalarInput.B524DHWVentilationScalarRaw != nil {
		b524DHWVentilationScalarResult.B524DHWVentilationScalarInput.B524DHWVentilationScalarRaw = append([]byte{}, b524DHWVentilationScalarInput.B524DHWVentilationScalarRaw...)
	}
	b524DHWVentilationScalarContext := b524DHWVentilationScalarInput.B524DHWVentilationScalarContext
	if !b524DHWVentilationScalarContext.B524DHWVentilationScalarFamily.B524DHWVentilationScalarPresent || !b524DHWVentilationScalarContext.B524DHWVentilationScalarSoftware.B524DHWVentilationScalarPresent ||
		!b524DHWVentilationScalarContext.B524DHWVentilationScalarHardware.B524DHWVentilationScalarPresent || !b524DHWVentilationScalarContext.B524DHWVentilationScalarAPIVersion.B524DHWVentilationScalarPresent ||
		!b524DHWVentilationScalarContext.B524DHWVentilationScalarAPIRevision.B524DHWVentilationScalarPresent || !b524DHWVentilationScalarContext.B524DHWVentilationScalarOpcode.B524DHWVentilationScalarPresent ||
		!b524DHWVentilationScalarContext.B524DHWVentilationScalarGroup.B524DHWVentilationScalarPresent || !b524DHWVentilationScalarContext.B524DHWVentilationScalarInstance.B524DHWVentilationScalarPresent ||
		!b524DHWVentilationScalarContext.B524DHWVentilationScalarRegister.B524DHWVentilationScalarPresent {
		return b524DHWVentilationScalarResult
	}
	if b524DHWVentilationScalarContext.B524DHWVentilationScalarOpcode.B524DHWVentilationScalarValue != 0x02 || b524DHWVentilationScalarContext.B524DHWVentilationScalarInstance.B524DHWVentilationScalarValue != 0 {
		return b524DHWVentilationScalarResult
	}
	switch b524DHWVentilationScalarContext.B524DHWVentilationScalarFamily.B524DHWVentilationScalarValue {
	case "controller-720", "controller-700":
		if b524DHWVentilationScalarContext.B524DHWVentilationScalarGroup.B524DHWVentilationScalarValue == 0x01 &&
			(b524DHWVentilationScalarContext.B524DHWVentilationScalarRegister.B524DHWVentilationScalarValue == 0x0004 || b524DHWVentilationScalarContext.B524DHWVentilationScalarRegister.B524DHWVentilationScalarValue == 0x0005) {
			b524DHWVentilationScalarResult.B524DHWVentilationScalarCodec = "float32_le"
			b524DHWVentilationScalarResult.B524DHWVentilationScalarDeclaredUnit = "°C"
		} else if b524DHWVentilationScalarContext.B524DHWVentilationScalarFamily.B524DHWVentilationScalarValue == "controller-720" && b524DHWVentilationScalarContext.B524DHWVentilationScalarGroup.B524DHWVentilationScalarValue == 0x09 &&
			(b524DHWVentilationScalarContext.B524DHWVentilationScalarRegister.B524DHWVentilationScalarValue == 0x000D || b524DHWVentilationScalarContext.B524DHWVentilationScalarRegister.B524DHWVentilationScalarValue == 0x000E) {
			b524DHWVentilationScalarResult.B524DHWVentilationScalarCodec = "uint16_le"
		}
	}
	if b524DHWVentilationScalarResult.B524DHWVentilationScalarCodec == "" || !b524DHWVentilationScalarInput.B524DHWVentilationScalarReplyKind.B524DHWVentilationScalarPresent ||
		b524DHWVentilationScalarInput.B524DHWVentilationScalarReplyKind.B524DHWVentilationScalarValue != "echo_payload" ||
		!b524DHWVentilationScalarInput.B524DHWVentilationScalarReturnedGroup.B524DHWVentilationScalarPresent || !b524DHWVentilationScalarInput.B524DHWVentilationScalarReturnedRegister.B524DHWVentilationScalarPresent ||
		b524DHWVentilationScalarInput.B524DHWVentilationScalarReturnedGroup.B524DHWVentilationScalarValue != b524DHWVentilationScalarContext.B524DHWVentilationScalarGroup.B524DHWVentilationScalarValue ||
		b524DHWVentilationScalarInput.B524DHWVentilationScalarReturnedRegister.B524DHWVentilationScalarValue != b524DHWVentilationScalarContext.B524DHWVentilationScalarRegister.B524DHWVentilationScalarValue ||
		!b524DHWVentilationScalarInput.B524DHWVentilationScalarCodecGrade.B524DHWVentilationScalarPresent ||
		(b524DHWVentilationScalarInput.B524DHWVentilationScalarCodecGrade.B524DHWVentilationScalarValue != "Hypothesis" && b524DHWVentilationScalarInput.B524DHWVentilationScalarCodecGrade.B524DHWVentilationScalarValue != "Proven") {
		return b524DHWVentilationScalarResult
	}
	b524DHWVentilationScalarScalar := &B524DHWVentilationScalar{}
	switch b524DHWVentilationScalarResult.B524DHWVentilationScalarCodec {
	case "float32_le":
		if len(b524DHWVentilationScalarInput.B524DHWVentilationScalarRaw) != 4 {
			return b524DHWVentilationScalarResult
		}
		b524DHWVentilationScalarScalar.B524DHWVentilationScalarFloat32 = b524DHWVentilationScalarMath.Float32frombits(b524DHWVentilationScalarBinary.LittleEndian.Uint32(b524DHWVentilationScalarInput.B524DHWVentilationScalarRaw))
		if b524DHWVentilationScalarMath.IsNaN(float64(b524DHWVentilationScalarScalar.B524DHWVentilationScalarFloat32)) || b524DHWVentilationScalarMath.IsInf(float64(b524DHWVentilationScalarScalar.B524DHWVentilationScalarFloat32), 0) {
			return b524DHWVentilationScalarResult
		}
	case "uint16_le":
		if len(b524DHWVentilationScalarInput.B524DHWVentilationScalarRaw) != 2 {
			return b524DHWVentilationScalarResult
		}
		b524DHWVentilationScalarScalar.B524DHWVentilationScalarUint16 = b524DHWVentilationScalarBinary.LittleEndian.Uint16(b524DHWVentilationScalarInput.B524DHWVentilationScalarRaw)
	}
	b524DHWVentilationScalarResult.B524DHWVentilationScalarState = B524DHWVentilationScalarCandidate
	b524DHWVentilationScalarResult.B524DHWVentilationScalarGrade = "Hypothesis"
	b524DHWVentilationScalarResult.B524DHWVentilationScalarValue = b524DHWVentilationScalarScalar
	return b524DHWVentilationScalarResult
}
