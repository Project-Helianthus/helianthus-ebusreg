package vaillant

import b51BLocalBudgetCRCErrors "errors"

type B51BLocalBudgetCRCProcedure string

const (
	B51BLocalBudgetCRCGrouped1024 B51BLocalBudgetCRCProcedure = "ab23an23_grouped_1024"
	B51BLocalBudgetCRCSingle205   B51BLocalBudgetCRCProcedure = "ab23an23_single_205"
)

type B51BLocalBudgetCRCBudget struct {
	B51BLocalBudgetCRCProcedure B51BLocalBudgetCRCProcedure `json:"Ab23An23Procedure"`
	B51BLocalBudgetCRCChunks    [][]byte                    `json:"Ab23An23Chunks"`
	B51BLocalBudgetCRCConsumed  int                         `json:"Ab23An23Consumed"`
	B51BLocalBudgetCRCUnplanned int                         `json:"Ab23An23Unplanned"`
}

func B51BLocalBudgetCRCPlan(b51BLocalBudgetCRCInput []byte, b51BLocalBudgetCRCProcedure B51BLocalBudgetCRCProcedure) (B51BLocalBudgetCRCBudget, error) {
	b51BLocalBudgetCRCResult := B51BLocalBudgetCRCBudget{B51BLocalBudgetCRCProcedure: b51BLocalBudgetCRCProcedure, B51BLocalBudgetCRCUnplanned: len(b51BLocalBudgetCRCInput)}
	var b51BLocalBudgetCRCCaps []int
	switch b51BLocalBudgetCRCProcedure {
	case B51BLocalBudgetCRCGrouped1024:
		b51BLocalBudgetCRCCaps = []int{206, 206, 206, 206, 200}
	case B51BLocalBudgetCRCSingle205:
		b51BLocalBudgetCRCCaps = []int{205}
	default:
		return b51BLocalBudgetCRCResult, b51BLocalBudgetCRCErrors.New("ab23an23: unknown procedure")
	}
	for _, b51BLocalBudgetCRCCap := range b51BLocalBudgetCRCCaps {
		if b51BLocalBudgetCRCResult.B51BLocalBudgetCRCUnplanned == 0 {
			break
		}
		b51BLocalBudgetCRCLength := b51BLocalBudgetCRCCap
		if b51BLocalBudgetCRCLength > b51BLocalBudgetCRCResult.B51BLocalBudgetCRCUnplanned {
			b51BLocalBudgetCRCLength = b51BLocalBudgetCRCResult.B51BLocalBudgetCRCUnplanned
		}
		b51BLocalBudgetCRCEnd := b51BLocalBudgetCRCResult.B51BLocalBudgetCRCConsumed + b51BLocalBudgetCRCLength
		b51BLocalBudgetCRCChunk := append([]byte(nil), b51BLocalBudgetCRCInput[b51BLocalBudgetCRCResult.B51BLocalBudgetCRCConsumed:b51BLocalBudgetCRCEnd]...)
		b51BLocalBudgetCRCResult.B51BLocalBudgetCRCChunks = append(b51BLocalBudgetCRCResult.B51BLocalBudgetCRCChunks, b51BLocalBudgetCRCChunk)
		b51BLocalBudgetCRCResult.B51BLocalBudgetCRCConsumed = b51BLocalBudgetCRCEnd
		b51BLocalBudgetCRCResult.B51BLocalBudgetCRCUnplanned -= b51BLocalBudgetCRCLength
	}
	return b51BLocalBudgetCRCResult, nil
}

func B51BLocalBudgetCRCProgress(b51BLocalBudgetCRCPrior int, b51BLocalBudgetCRCAccepted bool, b51BLocalBudgetCRCActual int) (int, error) {
	if b51BLocalBudgetCRCPrior < 0 || b51BLocalBudgetCRCActual < 0 {
		return b51BLocalBudgetCRCPrior, b51BLocalBudgetCRCErrors.New("ab23an23: negative progress or length")
	}
	if !b51BLocalBudgetCRCAccepted {
		return b51BLocalBudgetCRCPrior, nil
	}
	if b51BLocalBudgetCRCActual > int(^uint(0)>>1)-b51BLocalBudgetCRCPrior {
		return b51BLocalBudgetCRCPrior, b51BLocalBudgetCRCErrors.New("ab23an23: progress overflow")
	}
	return b51BLocalBudgetCRCPrior + b51BLocalBudgetCRCActual, nil
}

func B51BLocalBudgetCRC(b51BLocalBudgetCRCBytes []byte) uint16 {
	return B51BLocalBudgetCRCWithInitial(b51BLocalBudgetCRCBytes, 0xffff)
}

func B51BLocalBudgetCRCWithInitial(b51BLocalBudgetCRCBytes []byte, b51BLocalBudgetCRCInitial uint16) uint16 {
	b51BLocalBudgetCRCCRC := b51BLocalBudgetCRCInitial
	for _, b51BLocalBudgetCRCByte := range b51BLocalBudgetCRCBytes {
		b51BLocalBudgetCRCCRC ^= uint16(b51BLocalBudgetCRCByte) << 8
		for b51BLocalBudgetCRCBit := 0; b51BLocalBudgetCRCBit < 8; b51BLocalBudgetCRCBit++ {
			if b51BLocalBudgetCRCCRC&0x8000 != 0 {
				b51BLocalBudgetCRCCRC = (b51BLocalBudgetCRCCRC << 1) ^ 0x1021
			} else {
				b51BLocalBudgetCRCCRC <<= 1
			}
		}
	}
	return b51BLocalBudgetCRCCRC
}
