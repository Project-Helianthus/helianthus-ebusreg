package vaillant

import (
	b51BLocalBudgetCRCBytes "bytes"
	b51BLocalBudgetCRCJSON "encoding/json"
	b51BLocalBudgetCRCOS "os"
	b51BLocalBudgetCRCTesting "testing"
)

type b51BLocalBudgetCRCVectors struct {
	B51BLocalBudgetCRCBudgets []struct {
		B51BLocalBudgetCRCLabel     string                      `json:"Ab23An23Label"`
		B51BLocalBudgetCRCProcedure B51BLocalBudgetCRCProcedure `json:"Ab23An23Procedure"`
		B51BLocalBudgetCRCLength    int                         `json:"Ab23An23Length"`
		B51BLocalBudgetCRCChunks    []int                       `json:"Ab23An23Chunks"`
		B51BLocalBudgetCRCConsumed  int                         `json:"Ab23An23Consumed"`
		B51BLocalBudgetCRCUnplanned int                         `json:"Ab23An23Unplanned"`
	} `json:"Ab23An23Budgets"`
	B51BLocalBudgetCRCCRCs []struct {
		B51BLocalBudgetCRCLabel    string `json:"Ab23An23Label"`
		B51BLocalBudgetCRCBytes    []byte `json:"Ab23An23Bytes"`
		B51BLocalBudgetCRCInitial  uint16 `json:"Ab23An23Initial"`
		B51BLocalBudgetCRCExpected uint16 `json:"Ab23An23Expected"`
	} `json:"Ab23An23CRCs"`
}

func b51BLocalBudgetCRCLoad(b51BLocalBudgetCRCT *b51BLocalBudgetCRCTesting.T) b51BLocalBudgetCRCVectors {
	b51BLocalBudgetCRCT.Helper()
	b51BLocalBudgetCRCData, b51BLocalBudgetCRCErr := b51BLocalBudgetCRCOS.ReadFile("testdata/ab23an23_b51b_local_budget_crc.json")
	if b51BLocalBudgetCRCErr != nil {
		b51BLocalBudgetCRCT.Fatal(b51BLocalBudgetCRCErr)
	}
	var b51BLocalBudgetCRCFixture b51BLocalBudgetCRCVectors
	if b51BLocalBudgetCRCErr = b51BLocalBudgetCRCJSON.Unmarshal(b51BLocalBudgetCRCData, &b51BLocalBudgetCRCFixture); b51BLocalBudgetCRCErr != nil {
		b51BLocalBudgetCRCT.Fatal(b51BLocalBudgetCRCErr)
	}
	if len(b51BLocalBudgetCRCFixture.B51BLocalBudgetCRCBudgets) != 11 || len(b51BLocalBudgetCRCFixture.B51BLocalBudgetCRCCRCs) != 16 {
		b51BLocalBudgetCRCT.Fatal("ab23an23: incomplete vectors")
	}
	return b51BLocalBudgetCRCFixture
}

func TestAb23An23Budgets(b51BLocalBudgetCRCT *b51BLocalBudgetCRCTesting.T) {
	for _, b51BLocalBudgetCRCVector := range b51BLocalBudgetCRCLoad(b51BLocalBudgetCRCT).B51BLocalBudgetCRCBudgets {
		b51BLocalBudgetCRCT.Run(b51BLocalBudgetCRCVector.B51BLocalBudgetCRCLabel, func(b51BLocalBudgetCRCT *b51BLocalBudgetCRCTesting.T) {
			b51BLocalBudgetCRCInput := make([]byte, b51BLocalBudgetCRCVector.B51BLocalBudgetCRCLength)
			for b51BLocalBudgetCRCIndex := range b51BLocalBudgetCRCInput {
				b51BLocalBudgetCRCInput[b51BLocalBudgetCRCIndex] = byte(b51BLocalBudgetCRCIndex % 251)
			}
			b51BLocalBudgetCRCOriginal := append([]byte(nil), b51BLocalBudgetCRCInput...)
			b51BLocalBudgetCRCResult, b51BLocalBudgetCRCErr := B51BLocalBudgetCRCPlan(b51BLocalBudgetCRCInput, b51BLocalBudgetCRCVector.B51BLocalBudgetCRCProcedure)
			if b51BLocalBudgetCRCErr != nil {
				b51BLocalBudgetCRCT.Fatal(b51BLocalBudgetCRCErr)
			}
			if b51BLocalBudgetCRCResult.B51BLocalBudgetCRCProcedure != b51BLocalBudgetCRCVector.B51BLocalBudgetCRCProcedure || b51BLocalBudgetCRCResult.B51BLocalBudgetCRCConsumed != b51BLocalBudgetCRCVector.B51BLocalBudgetCRCConsumed || b51BLocalBudgetCRCResult.B51BLocalBudgetCRCUnplanned != b51BLocalBudgetCRCVector.B51BLocalBudgetCRCUnplanned || len(b51BLocalBudgetCRCResult.B51BLocalBudgetCRCChunks) != len(b51BLocalBudgetCRCVector.B51BLocalBudgetCRCChunks) {
				b51BLocalBudgetCRCT.Fatalf("ab23an23: wrong budget: %+v", b51BLocalBudgetCRCResult)
			}
			b51BLocalBudgetCRCOffset := 0
			for b51BLocalBudgetCRCIndex, b51BLocalBudgetCRCChunk := range b51BLocalBudgetCRCResult.B51BLocalBudgetCRCChunks {
				b51BLocalBudgetCRCLength := b51BLocalBudgetCRCVector.B51BLocalBudgetCRCChunks[b51BLocalBudgetCRCIndex]
				if len(b51BLocalBudgetCRCChunk) != b51BLocalBudgetCRCLength || !b51BLocalBudgetCRCBytes.Equal(b51BLocalBudgetCRCChunk, b51BLocalBudgetCRCOriginal[b51BLocalBudgetCRCOffset:b51BLocalBudgetCRCOffset+b51BLocalBudgetCRCLength]) {
					b51BLocalBudgetCRCT.Fatal("ab23an23: wrong chunk bytes or length")
				}
				b51BLocalBudgetCRCOffset += b51BLocalBudgetCRCLength
			}
			for _, b51BLocalBudgetCRCChunk := range b51BLocalBudgetCRCResult.B51BLocalBudgetCRCChunks {
				for b51BLocalBudgetCRCIndex := range b51BLocalBudgetCRCChunk {
					b51BLocalBudgetCRCChunk[b51BLocalBudgetCRCIndex] ^= 255
				}
			}
			if !b51BLocalBudgetCRCBytes.Equal(b51BLocalBudgetCRCInput, b51BLocalBudgetCRCOriginal) {
				b51BLocalBudgetCRCT.Fatal("ab23an23: output aliases input")
			}
			b51BLocalBudgetCRCAgain, b51BLocalBudgetCRCErr := B51BLocalBudgetCRCPlan(b51BLocalBudgetCRCInput, b51BLocalBudgetCRCVector.B51BLocalBudgetCRCProcedure)
			if b51BLocalBudgetCRCErr != nil || !b51BLocalBudgetCRCBytes.Equal(b51BLocalBudgetCRCBytes.Join(b51BLocalBudgetCRCAgain.B51BLocalBudgetCRCChunks, nil), b51BLocalBudgetCRCOriginal[:b51BLocalBudgetCRCVector.B51BLocalBudgetCRCConsumed]) {
				b51BLocalBudgetCRCT.Fatal("ab23an23: later call changed")
			}
			for b51BLocalBudgetCRCIndex := range b51BLocalBudgetCRCInput {
				b51BLocalBudgetCRCInput[b51BLocalBudgetCRCIndex] ^= 1
			}
			if !b51BLocalBudgetCRCBytes.Equal(b51BLocalBudgetCRCBytes.Join(b51BLocalBudgetCRCAgain.B51BLocalBudgetCRCChunks, nil), b51BLocalBudgetCRCOriginal[:b51BLocalBudgetCRCVector.B51BLocalBudgetCRCConsumed]) {
				b51BLocalBudgetCRCT.Fatal("ab23an23: input aliases output")
			}
		})
	}
	b51BLocalBudgetCRCResult, b51BLocalBudgetCRCErr := B51BLocalBudgetCRCPlan([]byte{1, 2, 3}, B51BLocalBudgetCRCProcedure("ab23an23_unknown"))
	if b51BLocalBudgetCRCErr == nil || b51BLocalBudgetCRCResult.B51BLocalBudgetCRCProcedure != "ab23an23_unknown" || b51BLocalBudgetCRCResult.B51BLocalBudgetCRCConsumed != 0 || b51BLocalBudgetCRCResult.B51BLocalBudgetCRCUnplanned != 3 || len(b51BLocalBudgetCRCResult.B51BLocalBudgetCRCChunks) != 0 {
		b51BLocalBudgetCRCT.Fatal("ab23an23: unknown choice admitted")
	}
}

func TestAb23An23Progress(b51BLocalBudgetCRCT *b51BLocalBudgetCRCTesting.T) {
	b51BLocalBudgetCRCMax := int(^uint(0) >> 1)
	for _, b51BLocalBudgetCRCCase := range []struct {
		b51BLocalBudgetCRCPrior    int
		b51BLocalBudgetCRCAccepted bool
		b51BLocalBudgetCRCActual   int
		b51BLocalBudgetCRCWant     int
		b51BLocalBudgetCRCError    bool
	}{
		{10, true, 17, 27, false}, {10, false, 17, 10, false}, {11, true, 17, 28, false}, {10, true, 18, 28, false},
		{10, true, 0, 10, false}, {0, true, 17, 17, false}, {b51BLocalBudgetCRCMax, true, 0, b51BLocalBudgetCRCMax, false},
		{b51BLocalBudgetCRCMax - 17, true, 17, b51BLocalBudgetCRCMax, false}, {b51BLocalBudgetCRCMax, true, 1, b51BLocalBudgetCRCMax, true},
		{b51BLocalBudgetCRCMax, false, 17, b51BLocalBudgetCRCMax, false}, {-1, true, 17, -1, true}, {10, true, -1, 10, true},
		{-1, false, 17, -1, true}, {10, false, -1, 10, true},
	} {
		b51BLocalBudgetCRCGot, b51BLocalBudgetCRCErr := B51BLocalBudgetCRCProgress(b51BLocalBudgetCRCCase.b51BLocalBudgetCRCPrior, b51BLocalBudgetCRCCase.b51BLocalBudgetCRCAccepted, b51BLocalBudgetCRCCase.b51BLocalBudgetCRCActual)
		if b51BLocalBudgetCRCGot != b51BLocalBudgetCRCCase.b51BLocalBudgetCRCWant || (b51BLocalBudgetCRCErr != nil) != b51BLocalBudgetCRCCase.b51BLocalBudgetCRCError {
			b51BLocalBudgetCRCT.Fatalf("ab23an23: transition %+v got %d, %v", b51BLocalBudgetCRCCase, b51BLocalBudgetCRCGot, b51BLocalBudgetCRCErr)
		}
	}
}

func TestAb23An23CRC(b51BLocalBudgetCRCT *b51BLocalBudgetCRCTesting.T) {
	for _, b51BLocalBudgetCRCVector := range b51BLocalBudgetCRCLoad(b51BLocalBudgetCRCT).B51BLocalBudgetCRCCRCs {
		b51BLocalBudgetCRCT.Run(b51BLocalBudgetCRCVector.B51BLocalBudgetCRCLabel, func(b51BLocalBudgetCRCT *b51BLocalBudgetCRCTesting.T) {
			b51BLocalBudgetCRCOriginal := append([]byte(nil), b51BLocalBudgetCRCVector.B51BLocalBudgetCRCBytes...)
			if b51BLocalBudgetCRCGot := B51BLocalBudgetCRCWithInitial(b51BLocalBudgetCRCVector.B51BLocalBudgetCRCBytes, b51BLocalBudgetCRCVector.B51BLocalBudgetCRCInitial); b51BLocalBudgetCRCGot != b51BLocalBudgetCRCVector.B51BLocalBudgetCRCExpected {
				b51BLocalBudgetCRCT.Fatalf("ab23an23: CRC got %04X want %04X", b51BLocalBudgetCRCGot, b51BLocalBudgetCRCVector.B51BLocalBudgetCRCExpected)
			}
			if b51BLocalBudgetCRCVector.B51BLocalBudgetCRCInitial == 65535 && B51BLocalBudgetCRC(b51BLocalBudgetCRCVector.B51BLocalBudgetCRCBytes) != b51BLocalBudgetCRCVector.B51BLocalBudgetCRCExpected {
				b51BLocalBudgetCRCT.Fatal("ab23an23: wrong default")
			}
			if !b51BLocalBudgetCRCBytes.Equal(b51BLocalBudgetCRCVector.B51BLocalBudgetCRCBytes, b51BLocalBudgetCRCOriginal) {
				b51BLocalBudgetCRCT.Fatal("ab23an23: CRC mutated input")
			}
		})
	}
	b51BLocalBudgetCRCSpan := []byte("x123456789y")
	if B51BLocalBudgetCRC(b51BLocalBudgetCRCSpan[1:10]) != 0x29b1 || B51BLocalBudgetCRC(b51BLocalBudgetCRCSpan[1:10:11]) != 0x29b1 || B51BLocalBudgetCRC(nil) != 0xffff || B51BLocalBudgetCRCWithInitial([]byte("123456789"), 0) != 0x31c3 {
		b51BLocalBudgetCRCT.Fatal("ab23an23: caller span or canonical CRC mismatch")
	}
}
