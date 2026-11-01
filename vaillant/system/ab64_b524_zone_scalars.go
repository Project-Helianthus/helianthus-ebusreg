package system

import (
	ab64Binary "encoding/binary"
	ab64Math "math"
)

type Ab64Field[Ab64T any] struct {
	Ab64Present bool
	Ab64Value   Ab64T
}

type Ab64Grade string

const (
	Ab64Unknown    Ab64Grade = "Unknown"
	Ab64Hypothesis Ab64Grade = "Hypothesis"
	Ab64Proven     Ab64Grade = "Proven"
)

type Ab64Context struct {
	Ab64Family      Ab64Field[string]
	Ab64Software    Ab64Field[string]
	Ab64Hardware    Ab64Field[string]
	Ab64APIVersion  Ab64Field[string]
	Ab64APIRevision Ab64Field[string]
	Ab64Opcode      Ab64Field[byte]
	Ab64Group       Ab64Field[byte]
	Ab64Instance    Ab64Field[byte]
	Ab64Register    Ab64Field[uint16]
}

type Ab64Input struct {
	Ab64Context          Ab64Context
	Ab64ReplyKind        string
	Ab64ReturnedGroup    Ab64Field[byte]
	Ab64ReturnedRegister Ab64Field[uint16]
	Ab64CodecGrade       Ab64Grade
	Ab64Raw              []byte
}

type Ab64State string

const (
	Ab64RawOnly   Ab64State = "raw_only"
	Ab64Candidate Ab64State = "candidate"
)

type Ab64Result struct {
	Ab64Input         Ab64Input
	Ab64State         Ab64State
	Ab64Value         Ab64Field[float32]
	Ab64Grade         Ab64Grade
	Ab64Unit          string
	Ab64Codec         string
	Ab64InstanceRule  string
	Ab64NativeSupport Ab64Grade
	Ab64Qualification Ab64Grade
	Ab64Presence      Ab64Grade
	Ab64Freshness     Ab64Grade
	Ab64Execution     Ab64Grade
	Ab64Effects       Ab64Grade
	Ab64Persistence   Ab64Grade
	Ab64Delivery      Ab64Grade
}

func Ab64Interpret(ab64Input Ab64Input) Ab64Result {
	ab64Input.Ab64Raw = append([]byte(nil), ab64Input.Ab64Raw...)
	ab64Result := Ab64Result{
		Ab64Input: ab64Input, Ab64State: Ab64RawOnly, Ab64Grade: Ab64Unknown,
		Ab64Unit: "Unknown", Ab64NativeSupport: Ab64Unknown, Ab64Qualification: Ab64Unknown,
		Ab64Presence: Ab64Unknown, Ab64Freshness: Ab64Unknown, Ab64Execution: Ab64Unknown,
		Ab64Effects: Ab64Unknown, Ab64Persistence: Ab64Unknown, Ab64Delivery: Ab64Unknown,
	}
	ab64Context := ab64Input.Ab64Context
	if !ab64Context.Ab64Family.Ab64Present || !ab64Context.Ab64Software.Ab64Present ||
		!ab64Context.Ab64Hardware.Ab64Present || !ab64Context.Ab64APIVersion.Ab64Present ||
		!ab64Context.Ab64APIRevision.Ab64Present || !ab64Context.Ab64Opcode.Ab64Present ||
		!ab64Context.Ab64Group.Ab64Present || !ab64Context.Ab64Instance.Ab64Present ||
		!ab64Context.Ab64Register.Ab64Present {
		return ab64Result
	}
	if ab64Context.Ab64Family.Ab64Value != "controller-720" && ab64Context.Ab64Family.Ab64Value != "controller-700" {
		return ab64Result
	}
	if ab64Context.Ab64Opcode.Ab64Value != 0x02 || ab64Context.Ab64Group.Ab64Value != 0x03 {
		return ab64Result
	}
	switch ab64Context.Ab64Register.Ab64Value {
	case 0x0002, 0x0005, 0x0009, 0x000f, 0x0014, 0x0015:
		ab64Result.Ab64Unit = "°C"
	case 0x0023:
		if ab64Context.Ab64Family.Ab64Value != "controller-720" {
			return ab64Result
		}
		ab64Result.Ab64Unit = "°C"
	case 0x0028:
		ab64Result.Ab64Unit = "%"
	default:
		return ab64Result
	}
	ab64Result.Ab64Codec = "float32_le"
	ab64Result.Ab64InstanceRule = "dynamic"
	if (ab64Input.Ab64CodecGrade != Ab64Hypothesis && ab64Input.Ab64CodecGrade != Ab64Proven) ||
		ab64Input.Ab64ReplyKind != "echo_payload" || !ab64Input.Ab64ReturnedGroup.Ab64Present ||
		!ab64Input.Ab64ReturnedRegister.Ab64Present ||
		ab64Input.Ab64ReturnedGroup.Ab64Value != ab64Context.Ab64Group.Ab64Value ||
		ab64Input.Ab64ReturnedRegister.Ab64Value != ab64Context.Ab64Register.Ab64Value || len(ab64Input.Ab64Raw) != 4 {
		return ab64Result
	}
	ab64Value := ab64Math.Float32frombits(ab64Binary.LittleEndian.Uint32(ab64Input.Ab64Raw))
	if ab64Math.IsNaN(float64(ab64Value)) || ab64Math.IsInf(float64(ab64Value), 0) {
		return ab64Result
	}
	ab64Result.Ab64Value = Ab64Field[float32]{Ab64Present: true, Ab64Value: ab64Value}
	ab64Result.Ab64State = Ab64Candidate
	ab64Result.Ab64Grade = Ab64Hypothesis
	return ab64Result
}
