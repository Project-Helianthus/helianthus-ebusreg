package system

import (
	ab62Binary "encoding/binary"
	ab62Math "math"
)

type Ab62Field[Ab62T any] struct {
	Ab62Present bool
	Ab62Value   Ab62T
}

type Ab62Grade string

const (
	Ab62Unknown    Ab62Grade = "Unknown"
	Ab62Hypothesis Ab62Grade = "Hypothesis"
	Ab62Proven     Ab62Grade = "Proven"
)

type Ab62Context struct {
	Ab62Family      Ab62Field[string]
	Ab62Software    Ab62Field[string]
	Ab62Hardware    Ab62Field[string]
	Ab62APIVersion  Ab62Field[string]
	Ab62APIRevision Ab62Field[string]
	Ab62Opcode      Ab62Field[byte]
	Ab62Group       Ab62Field[byte]
	Ab62Instance    Ab62Field[byte]
	Ab62Register    Ab62Field[uint16]
}

type Ab62Input struct {
	Ab62Context          Ab62Context
	Ab62ReplyKind        Ab62Field[string]
	Ab62ReturnedGroup    Ab62Field[byte]
	Ab62ReturnedRegister Ab62Field[uint16]
	Ab62CodecGrade       Ab62Field[Ab62Grade]
	Ab62Raw              []byte
}

type Ab62State string

const (
	Ab62RawOnly   Ab62State = "raw_only"
	Ab62Candidate Ab62State = "candidate"
)

type Ab62Result struct {
	Ab62Input        Ab62Input
	Ab62State        Ab62State
	Ab62Grade        Ab62Grade
	Ab62DeclaredUnit string
	Ab62Codec        string
	Ab62Value        Ab62Field[float32]
}

func Ab62InterpretCircuitScalar(ab62Input Ab62Input) Ab62Result {
	if ab62Input.Ab62Raw != nil {
		ab62Raw := make([]byte, len(ab62Input.Ab62Raw))
		copy(ab62Raw, ab62Input.Ab62Raw)
		ab62Input.Ab62Raw = ab62Raw
	}
	ab62Result := Ab62Result{Ab62Input: ab62Input, Ab62State: Ab62RawOnly, Ab62Grade: Ab62Unknown, Ab62DeclaredUnit: "Unknown"}
	ab62Context := ab62Input.Ab62Context
	if !ab62Context.Ab62Family.Ab62Present || !ab62Context.Ab62Software.Ab62Present || !ab62Context.Ab62Hardware.Ab62Present || !ab62Context.Ab62APIVersion.Ab62Present || !ab62Context.Ab62APIRevision.Ab62Present || !ab62Context.Ab62Opcode.Ab62Present || !ab62Context.Ab62Group.Ab62Present || !ab62Context.Ab62Instance.Ab62Present || !ab62Context.Ab62Register.Ab62Present {
		return ab62Result
	}
	if ab62Context.Ab62Opcode.Ab62Value != 0x02 || ab62Context.Ab62Group.Ab62Value != 0x02 {
		return ab62Result
	}
	ab62Family := ab62Context.Ab62Family.Ab62Value
	if ab62Family != "controller-720" && ab62Family != "controller-700" {
		return ab62Result
	}
	switch ab62Context.Ab62Register.Ab62Value {
	case 0x0007, 0x0008, 0x0010, 0x0011, 0x0012, 0x0014, 0x001f:
		ab62Result.Ab62DeclaredUnit = "°C"
	case 0x000a, 0x000b:
		ab62Result.Ab62DeclaredUnit = "K"
	case 0x000f:
	case 0x001c:
		if ab62Family != "controller-720" {
			return ab62Result
		}
	case 0x001a:
		if ab62Family != "controller-720" {
			return ab62Result
		}
		ab62Result.Ab62DeclaredUnit = "%"
	default:
		return ab62Result
	}
	ab62Result.Ab62Codec = "float32_le"
	if !ab62Input.Ab62CodecGrade.Ab62Present || (ab62Input.Ab62CodecGrade.Ab62Value != Ab62Hypothesis && ab62Input.Ab62CodecGrade.Ab62Value != Ab62Proven) || !ab62Input.Ab62ReplyKind.Ab62Present || ab62Input.Ab62ReplyKind.Ab62Value != "echo_payload" || !ab62Input.Ab62ReturnedGroup.Ab62Present || ab62Input.Ab62ReturnedGroup.Ab62Value != ab62Context.Ab62Group.Ab62Value || !ab62Input.Ab62ReturnedRegister.Ab62Present || ab62Input.Ab62ReturnedRegister.Ab62Value != ab62Context.Ab62Register.Ab62Value || len(ab62Input.Ab62Raw) != 4 {
		return ab62Result
	}
	ab62Value := ab62Math.Float32frombits(ab62Binary.LittleEndian.Uint32(ab62Input.Ab62Raw))
	if ab62Math.IsNaN(float64(ab62Value)) || ab62Math.IsInf(float64(ab62Value), 0) {
		return ab62Result
	}
	ab62Result.Ab62State = Ab62Candidate
	ab62Result.Ab62Grade = Ab62Hypothesis
	ab62Result.Ab62Value = Ab62Field[float32]{Ab62Present: true, Ab62Value: ab62Value}
	return ab62Result
}
