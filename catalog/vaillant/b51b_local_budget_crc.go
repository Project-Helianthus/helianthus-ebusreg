package vaillant

import b51BLocalBudgetCRCErrors "errors"

// B51BLocalBudgetCRCProcedure names a local chunking procedure for a B51B transfer budget.
type B51BLocalBudgetCRCProcedure string

// Known procedures: grouped_1024 plans chunks of at most 206, 206, 206, 206 and 200 bytes (1,024 in all),
// and single_205 plans one chunk of at most 205 bytes.
const (
	B51BLocalBudgetCRCGrouped1024 B51BLocalBudgetCRCProcedure = "grouped_1024"
	B51BLocalBudgetCRCSingle205   B51BLocalBudgetCRCProcedure = "single_205"
)

// B51BLocalBudgetCRCBudget is a chunk plan: detached chunks in order, the bytes they consume and the input
// bytes left unplanned once the procedure's budget is spent.
type B51BLocalBudgetCRCBudget struct {
	B51BLocalBudgetCRCProcedure B51BLocalBudgetCRCProcedure `json:"procedure"`
	B51BLocalBudgetCRCChunks    [][]byte                    `json:"chunks"`
	B51BLocalBudgetCRCConsumed  int                         `json:"consumed"`
	B51BLocalBudgetCRCUnplanned int                         `json:"unplanned"`
}

// B51BLocalBudgetCRCPlan splits the input into detached chunks for the given procedure; input beyond the
// budget is counted as unplanned, and an unknown procedure is an error with everything unplanned.
func B51BLocalBudgetCRCPlan(b51BLocalBudgetCRCInput []byte, b51BLocalBudgetCRCProcedure B51BLocalBudgetCRCProcedure) (B51BLocalBudgetCRCBudget, error) {
	b51BLocalBudgetCRCResult := B51BLocalBudgetCRCBudget{B51BLocalBudgetCRCProcedure: b51BLocalBudgetCRCProcedure, B51BLocalBudgetCRCUnplanned: len(b51BLocalBudgetCRCInput)}
	var b51BLocalBudgetCRCCaps []int
	switch b51BLocalBudgetCRCProcedure {
	case B51BLocalBudgetCRCGrouped1024:
		b51BLocalBudgetCRCCaps = []int{206, 206, 206, 206, 200}
	case B51BLocalBudgetCRCSingle205:
		b51BLocalBudgetCRCCaps = []int{205}
	default:
		return b51BLocalBudgetCRCResult, b51BLocalBudgetCRCErrors.New("b51b local budget crc: unknown procedure")
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

// B51BLocalBudgetCRCProgress advances the progress by the actual length only when the chunk was accepted;
// negative values and overflow are errors and return the prior progress.
func B51BLocalBudgetCRCProgress(b51BLocalBudgetCRCPrior int, b51BLocalBudgetCRCAccepted bool, b51BLocalBudgetCRCActual int) (int, error) {
	if b51BLocalBudgetCRCPrior < 0 || b51BLocalBudgetCRCActual < 0 {
		return b51BLocalBudgetCRCPrior, b51BLocalBudgetCRCErrors.New("b51b local budget crc: negative progress or length")
	}
	if !b51BLocalBudgetCRCAccepted {
		return b51BLocalBudgetCRCPrior, nil
	}
	if b51BLocalBudgetCRCActual > int(^uint(0)>>1)-b51BLocalBudgetCRCPrior {
		return b51BLocalBudgetCRCPrior, b51BLocalBudgetCRCErrors.New("b51b local budget crc: progress overflow")
	}
	return b51BLocalBudgetCRCPrior + b51BLocalBudgetCRCActual, nil
}

// B51BLocalBudgetCRC returns the CRC-16 of the bytes with the initial value 0xFFFF.
func B51BLocalBudgetCRC(b51BLocalBudgetCRCBytes []byte) uint16 {
	return B51BLocalBudgetCRCWithInitial(b51BLocalBudgetCRCBytes, 0xffff)
}

// B51BLocalBudgetCRCWithInitial returns the CRC-16 of the bytes with polynomial 0x1021, most significant
// bit first, no reflection and no final XOR, starting from the given value; passing a previous result
// continues the checksum across chunks.
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
