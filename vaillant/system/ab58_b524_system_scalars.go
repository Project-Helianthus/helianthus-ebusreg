package system

import (
	ab58binary "encoding/binary"
	ab58math "math"
)

type Ab58Field[Ab58T any] struct {
	Ab58Present bool
	Ab58Value   Ab58T
}

type Ab58Context struct {
	Ab58Family      Ab58Field[string]
	Ab58Software    Ab58Field[string]
	Ab58Hardware    Ab58Field[string]
	Ab58APIVersion  Ab58Field[string]
	Ab58APIRevision Ab58Field[string]
	Ab58Opcode      Ab58Field[byte]
	Ab58Group       Ab58Field[byte]
	Ab58Instance    Ab58Field[byte]
	Ab58Register    Ab58Field[uint16]
}

type Ab58Input struct {
	Ab58Context          Ab58Context
	Ab58ReplyKind        Ab58Field[string]
	Ab58ReturnedGroup    Ab58Field[byte]
	Ab58ReturnedRegister Ab58Field[uint16]
	Ab58CodecGrade       Ab58Field[string]
	Ab58Raw              []byte
}

type Ab58State string

const (
	Ab58RawOnly   Ab58State = "raw_only"
	Ab58Candidate Ab58State = "candidate"
)

type Ab58Scalar struct {
	Ab58Float32 float32
	Ab58Uint16  uint16
}

type Ab58Result struct {
	Ab58Input        Ab58Input
	Ab58State        Ab58State
	Ab58Codec        string
	Ab58DeclaredUnit string
	Ab58Grade        string
	Ab58Presence     string
	Ab58Freshness    string
	Ab58Execution    string
	Ab58Support      string
	Ab58Effects      string
	Ab58Persistence  string
	Ab58Delivery     string
	Ab58Value        *Ab58Scalar
}

func Ab58InterpretB524Scalars(ab58input Ab58Input) Ab58Result {
	ab58result := Ab58Result{
		Ab58Input: ab58input, Ab58State: Ab58RawOnly,
		Ab58DeclaredUnit: "Unknown", Ab58Grade: "Unknown",
		Ab58Presence: "Unknown", Ab58Freshness: "Unknown", Ab58Execution: "Unknown",
		Ab58Support: "Unknown", Ab58Effects: "Unknown",
		Ab58Persistence: "Unknown", Ab58Delivery: "Unknown",
	}
	if ab58input.Ab58Raw != nil {
		ab58result.Ab58Input.Ab58Raw = append([]byte{}, ab58input.Ab58Raw...)
	}
	ab58context := ab58input.Ab58Context
	if !ab58context.Ab58Family.Ab58Present || !ab58context.Ab58Software.Ab58Present ||
		!ab58context.Ab58Hardware.Ab58Present || !ab58context.Ab58APIVersion.Ab58Present ||
		!ab58context.Ab58APIRevision.Ab58Present || !ab58context.Ab58Opcode.Ab58Present ||
		!ab58context.Ab58Group.Ab58Present || !ab58context.Ab58Instance.Ab58Present ||
		!ab58context.Ab58Register.Ab58Present {
		return ab58result
	}
	if ab58context.Ab58Opcode.Ab58Value != 0x02 || ab58context.Ab58Group.Ab58Value != 0 {
		return ab58result
	}
	ab58dynamic := false
	switch ab58context.Ab58Family.Ab58Value {
	case "controller-720":
		switch ab58context.Ab58Register.Ab58Value {
		case 0x0002, 0x0022, 0x0023, 0x0026, 0x0038, 0x0046, 0x004B, 0x0067, 0x0068, 0x0073, 0x0095, 0x009D, 0x009E, 0x009F, 0x00A0, 0x00EF:
			ab58result.Ab58Codec, ab58result.Ab58DeclaredUnit = "float32_le", "°C"
		case 0x0004, 0x0017, 0x0018:
			ab58result.Ab58Codec, ab58result.Ab58DeclaredUnit = "uint16_le", "min"
		case 0x0029:
			ab58result.Ab58Codec, ab58result.Ab58DeclaredUnit = "float32_le", "K"
		case 0x002F, 0x0036:
			ab58result.Ab58Codec = "uint16_le"
		case 0x0039:
			ab58result.Ab58Codec, ab58result.Ab58DeclaredUnit = "float32_le", "bar"
		}
	case "controller-700":
		switch ab58context.Ab58Register.Ab58Value {
		case 0x0002, 0x0022, 0x0023, 0x0026, 0x0038, 0x0046, 0x004B, 0x0067, 0x0068, 0x0073, 0x0095, 0x009D, 0x009E, 0x009F, 0x00A0:
			ab58result.Ab58Codec, ab58result.Ab58DeclaredUnit = "float32_le", "°C"
		case 0x0006, 0x002F, 0x0036:
			ab58result.Ab58Codec = "uint16_le"
		case 0x0017:
			ab58result.Ab58Codec, ab58result.Ab58DeclaredUnit = "uint16_le", "min"
		case 0x0029:
			ab58result.Ab58Codec, ab58result.Ab58DeclaredUnit = "float32_le", "K"
		case 0x0039:
			ab58result.Ab58Codec, ab58result.Ab58DeclaredUnit = "float32_le", "bar"
		case 0x0042, 0x0043:
			ab58result.Ab58Codec = "uint16_le"
			ab58dynamic = true
		}
	}
	if !ab58dynamic && ab58context.Ab58Instance.Ab58Value != 0 {
		return ab58result
	}
	if ab58result.Ab58Codec == "" || !ab58input.Ab58ReplyKind.Ab58Present ||
		ab58input.Ab58ReplyKind.Ab58Value != "echo_payload" ||
		!ab58input.Ab58ReturnedGroup.Ab58Present || !ab58input.Ab58ReturnedRegister.Ab58Present ||
		ab58input.Ab58ReturnedGroup.Ab58Value != ab58context.Ab58Group.Ab58Value ||
		ab58input.Ab58ReturnedRegister.Ab58Value != ab58context.Ab58Register.Ab58Value ||
		!ab58input.Ab58CodecGrade.Ab58Present ||
		(ab58input.Ab58CodecGrade.Ab58Value != "Hypothesis" && ab58input.Ab58CodecGrade.Ab58Value != "Proven") {
		return ab58result
	}
	ab58scalar := &Ab58Scalar{}
	switch ab58result.Ab58Codec {
	case "float32_le":
		if len(ab58input.Ab58Raw) != 4 {
			return ab58result
		}
		ab58scalar.Ab58Float32 = ab58math.Float32frombits(ab58binary.LittleEndian.Uint32(ab58input.Ab58Raw))
		if ab58math.IsNaN(float64(ab58scalar.Ab58Float32)) || ab58math.IsInf(float64(ab58scalar.Ab58Float32), 0) {
			return ab58result
		}
	case "uint16_le":
		if len(ab58input.Ab58Raw) != 2 {
			return ab58result
		}
		ab58scalar.Ab58Uint16 = ab58binary.LittleEndian.Uint16(ab58input.Ab58Raw)
	}
	ab58result.Ab58State = Ab58Candidate
	ab58result.Ab58Grade = "Hypothesis"
	ab58result.Ab58Value = ab58scalar
	return ab58result
}
