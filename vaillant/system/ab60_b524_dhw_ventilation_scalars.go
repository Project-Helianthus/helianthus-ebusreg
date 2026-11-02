package system

import (
	ab60binary "encoding/binary"
	ab60math "math"
)

type Ab60Field[Ab60T any] struct {
	Ab60Present bool
	Ab60Value   Ab60T
}

type Ab60Context struct {
	Ab60Family      Ab60Field[string]
	Ab60Software    Ab60Field[string]
	Ab60Hardware    Ab60Field[string]
	Ab60APIVersion  Ab60Field[string]
	Ab60APIRevision Ab60Field[string]
	Ab60Opcode      Ab60Field[byte]
	Ab60Group       Ab60Field[byte]
	Ab60Instance    Ab60Field[byte]
	Ab60Register    Ab60Field[uint16]
}

type Ab60Input struct {
	Ab60Context          Ab60Context
	Ab60ReplyKind        Ab60Field[string]
	Ab60ReturnedGroup    Ab60Field[byte]
	Ab60ReturnedRegister Ab60Field[uint16]
	Ab60CodecGrade       Ab60Field[string]
	Ab60Raw              []byte
}

type Ab60State string

const (
	Ab60RawOnly   Ab60State = "raw_only"
	Ab60Candidate Ab60State = "candidate"
)

type Ab60Scalar struct {
	Ab60Float32 float32
	Ab60Uint16  uint16
}

type Ab60Result struct {
	Ab60Input        Ab60Input
	Ab60State        Ab60State
	Ab60Codec        string
	Ab60DeclaredUnit string
	Ab60Grade        string
	Ab60Presence     string
	Ab60Freshness    string
	Ab60Execution    string
	Ab60Support      string
	Ab60Effects      string
	Ab60Persistence  string
	Ab60Delivery     string
	Ab60Value        *Ab60Scalar
}

func Ab60InterpretB524Scalars(ab60input Ab60Input) Ab60Result {
	ab60result := Ab60Result{
		Ab60Input: ab60input, Ab60State: Ab60RawOnly,
		Ab60DeclaredUnit: "Unknown", Ab60Grade: "Unknown",
		Ab60Presence: "Unknown", Ab60Freshness: "Unknown", Ab60Execution: "Unknown",
		Ab60Support: "Unknown", Ab60Effects: "Unknown",
		Ab60Persistence: "Unknown", Ab60Delivery: "Unknown",
	}
	if ab60input.Ab60Raw != nil {
		ab60result.Ab60Input.Ab60Raw = append([]byte{}, ab60input.Ab60Raw...)
	}
	ab60context := ab60input.Ab60Context
	if !ab60context.Ab60Family.Ab60Present || !ab60context.Ab60Software.Ab60Present ||
		!ab60context.Ab60Hardware.Ab60Present || !ab60context.Ab60APIVersion.Ab60Present ||
		!ab60context.Ab60APIRevision.Ab60Present || !ab60context.Ab60Opcode.Ab60Present ||
		!ab60context.Ab60Group.Ab60Present || !ab60context.Ab60Instance.Ab60Present ||
		!ab60context.Ab60Register.Ab60Present {
		return ab60result
	}
	if ab60context.Ab60Opcode.Ab60Value != 0x02 || ab60context.Ab60Instance.Ab60Value != 0 {
		return ab60result
	}
	switch ab60context.Ab60Family.Ab60Value {
	case "controller-720", "controller-700":
		if ab60context.Ab60Group.Ab60Value == 0x01 &&
			(ab60context.Ab60Register.Ab60Value == 0x0004 || ab60context.Ab60Register.Ab60Value == 0x0005) {
			ab60result.Ab60Codec = "float32_le"
			ab60result.Ab60DeclaredUnit = "°C"
		} else if ab60context.Ab60Family.Ab60Value == "controller-720" && ab60context.Ab60Group.Ab60Value == 0x09 &&
			(ab60context.Ab60Register.Ab60Value == 0x000D || ab60context.Ab60Register.Ab60Value == 0x000E) {
			ab60result.Ab60Codec = "uint16_le"
		}
	}
	if ab60result.Ab60Codec == "" || !ab60input.Ab60ReplyKind.Ab60Present ||
		ab60input.Ab60ReplyKind.Ab60Value != "echo_payload" ||
		!ab60input.Ab60ReturnedGroup.Ab60Present || !ab60input.Ab60ReturnedRegister.Ab60Present ||
		ab60input.Ab60ReturnedGroup.Ab60Value != ab60context.Ab60Group.Ab60Value ||
		ab60input.Ab60ReturnedRegister.Ab60Value != ab60context.Ab60Register.Ab60Value ||
		!ab60input.Ab60CodecGrade.Ab60Present ||
		(ab60input.Ab60CodecGrade.Ab60Value != "Hypothesis" && ab60input.Ab60CodecGrade.Ab60Value != "Proven") {
		return ab60result
	}
	ab60scalar := &Ab60Scalar{}
	switch ab60result.Ab60Codec {
	case "float32_le":
		if len(ab60input.Ab60Raw) != 4 {
			return ab60result
		}
		ab60scalar.Ab60Float32 = ab60math.Float32frombits(ab60binary.LittleEndian.Uint32(ab60input.Ab60Raw))
		if ab60math.IsNaN(float64(ab60scalar.Ab60Float32)) || ab60math.IsInf(float64(ab60scalar.Ab60Float32), 0) {
			return ab60result
		}
	case "uint16_le":
		if len(ab60input.Ab60Raw) != 2 {
			return ab60result
		}
		ab60scalar.Ab60Uint16 = ab60binary.LittleEndian.Uint16(ab60input.Ab60Raw)
	}
	ab60result.Ab60State = Ab60Candidate
	ab60result.Ab60Grade = "Hypothesis"
	ab60result.Ab60Value = ab60scalar
	return ab60result
}
