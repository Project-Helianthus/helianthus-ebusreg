package system

import (
	b524DHWVentilationScalarBinary "encoding/binary"
	b524DHWVentilationScalarMath "math"
)

type B524DHWVentilationScalarField[B524DHWVentilationScalarT any] struct {
	B524DHWVentilationScalarPresent bool                      `json:"present"`
	B524DHWVentilationScalarValue   B524DHWVentilationScalarT `json:"value"`
}

type B524DHWVentilationScalarContext struct {
	B524DHWVentilationScalarFamily      B524DHWVentilationScalarField[string] `json:"family"`
	B524DHWVentilationScalarSoftware    B524DHWVentilationScalarField[string] `json:"software"`
	B524DHWVentilationScalarHardware    B524DHWVentilationScalarField[string] `json:"hardware"`
	B524DHWVentilationScalarAPIVersion  B524DHWVentilationScalarField[string] `json:"api_version"`
	B524DHWVentilationScalarAPIRevision B524DHWVentilationScalarField[string] `json:"api_revision"`
	B524DHWVentilationScalarOpcode      B524DHWVentilationScalarField[byte]   `json:"opcode"`
	B524DHWVentilationScalarGroup       B524DHWVentilationScalarField[byte]   `json:"group"`
	B524DHWVentilationScalarInstance    B524DHWVentilationScalarField[byte]   `json:"instance"`
	B524DHWVentilationScalarRegister    B524DHWVentilationScalarField[uint16] `json:"register"`
}

type B524DHWVentilationScalarInput struct {
	B524DHWVentilationScalarContext          B524DHWVentilationScalarContext       `json:"context"`
	B524DHWVentilationScalarReplyKind        B524DHWVentilationScalarField[string] `json:"reply_kind"`
	B524DHWVentilationScalarReturnedGroup    B524DHWVentilationScalarField[byte]   `json:"returned_group"`
	B524DHWVentilationScalarReturnedRegister B524DHWVentilationScalarField[uint16] `json:"returned_register"`
	B524DHWVentilationScalarCodecGrade       B524DHWVentilationScalarField[string] `json:"codec_grade"`
	B524DHWVentilationScalarRaw              []byte                                `json:"raw"`
}

type B524DHWVentilationScalarState string

const (
	B524DHWVentilationScalarRawOnly   B524DHWVentilationScalarState = "raw_only"
	B524DHWVentilationScalarCandidate B524DHWVentilationScalarState = "candidate"
)

type B524DHWVentilationScalar struct {
	B524DHWVentilationScalarFloat32 float32 `json:"float32"`
	B524DHWVentilationScalarUint16  uint16  `json:"uint16"`
}

type B524DHWVentilationScalarResult struct {
	B524DHWVentilationScalarInput        B524DHWVentilationScalarInput `json:"input"`
	B524DHWVentilationScalarState        B524DHWVentilationScalarState `json:"state"`
	B524DHWVentilationScalarCodec        string                        `json:"codec"`
	B524DHWVentilationScalarDeclaredUnit string                        `json:"declared_unit"`
	B524DHWVentilationScalarGrade        string                        `json:"grade"`
	B524DHWVentilationScalarPresence     string                        `json:"presence"`
	B524DHWVentilationScalarFreshness    string                        `json:"freshness"`
	B524DHWVentilationScalarExecution    string                        `json:"execution"`
	B524DHWVentilationScalarSupport      string                        `json:"support"`
	B524DHWVentilationScalarEffects      string                        `json:"effects"`
	B524DHWVentilationScalarPersistence  string                        `json:"persistence"`
	B524DHWVentilationScalarDelivery     string                        `json:"delivery"`
	B524DHWVentilationScalarValue        *B524DHWVentilationScalar     `json:"value"`
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
