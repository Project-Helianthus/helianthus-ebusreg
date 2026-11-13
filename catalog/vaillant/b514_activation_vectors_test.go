package vaillant

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
)

type b514ActivationVectorCase struct {
	B514ActivationVectorID        string   `json:"ab80_id"`
	B514ActivationVectorOperation uint8    `json:"ab80_operation"`
	B514ActivationVector          []byte   `json:"ab80_vector"`
	B514ActivationVectorState     string   `json:"ab80_state"`
	B514ActivationVectorMask      uint16   `json:"ab80_mask"`
	B514ActivationVectorComplete  bool     `json:"ab80_complete"`
	B514ActivationVectorRemainder []byte   `json:"ab80_remainder"`
	B514ActivationVectorTrue      [][2]int `json:"ab80_true"`
}

func b514ActivationVectorIdent(b514ActivationVectorOperation uint8) B514ActivationVectorIdentity {
	return B514ActivationVectorIdentity{B514ActivationVectorFamily: B514ActivationVectorFamily, B514ActivationVectorVariant: B514ActivationVectorVariant, B514ActivationVectorOperation: b514ActivationVectorOperation, B514ActivationVectorDirection: B514ActivationVectorDirectionPositive}
}

func b514ActivationVectorLoad(b514ActivationVectorT *testing.T) []b514ActivationVectorCase {
	b514ActivationVectorT.Helper()
	b514ActivationVectorData, b514ActivationVectorErr := os.ReadFile("testdata/ab80_b514_activation_vectors.json")
	if b514ActivationVectorErr != nil {
		b514ActivationVectorT.Fatal(b514ActivationVectorErr)
	}
	var b514ActivationVectorCases []b514ActivationVectorCase
	if b514ActivationVectorErr = json.Unmarshal(b514ActivationVectorData, &b514ActivationVectorCases); b514ActivationVectorErr != nil {
		b514ActivationVectorT.Fatal(b514ActivationVectorErr)
	}
	if len(b514ActivationVectorCases) == 0 {
		b514ActivationVectorT.Fatal("no fixture cases")
	}
	return b514ActivationVectorCases
}

func TestAb80Fixtures(b514ActivationVectorT *testing.T) {
	for _, b514ActivationVectorCase := range b514ActivationVectorLoad(b514ActivationVectorT) {
		b514ActivationVectorT.Run(b514ActivationVectorCase.B514ActivationVectorID, func(b514ActivationVectorT *testing.T) {
			b514ActivationVectorView := B514ActivationVectorInterpret(b514ActivationVectorIdent(b514ActivationVectorCase.B514ActivationVectorOperation), b514ActivationVectorCase.B514ActivationVector)
			if b514ActivationVectorView.B514ActivationVectorState() != b514ActivationVectorCase.B514ActivationVectorState || b514ActivationVectorView.B514ActivationVectorGrade() != B514ActivationVectorGradeHypothesis {
				b514ActivationVectorT.Fatal(b514ActivationVectorView.B514ActivationVectorState(), b514ActivationVectorView.B514ActivationVectorGrade())
			}
			if b514ActivationVectorView.B514ActivationVectorPresentMask() != b514ActivationVectorCase.B514ActivationVectorMask || b514ActivationVectorView.B514ActivationVectorIsComplete() != b514ActivationVectorCase.B514ActivationVectorComplete {
				b514ActivationVectorT.Fatal(b514ActivationVectorView.B514ActivationVectorPresentMask(), b514ActivationVectorView.B514ActivationVectorIsComplete())
			}
			if !bytes.Equal(b514ActivationVectorView.B514ActivationVectorRawBytes(), b514ActivationVectorCase.B514ActivationVector) || !bytes.Equal(b514ActivationVectorView.B514ActivationVectorRemainderBytes(), b514ActivationVectorCase.B514ActivationVectorRemainder) {
				b514ActivationVectorT.Fatal("raw or remainder changed")
			}
			if b514ActivationVectorCase.B514ActivationVectorOperation == B514ActivationVectorOperationBits {
				b514ActivationVectorWant := map[[2]int]bool{}
				for _, b514ActivationVectorPos := range b514ActivationVectorCase.B514ActivationVectorTrue {
					b514ActivationVectorWant[b514ActivationVectorPos] = true
				}
				for b514ActivationVectorByte := 0; b514ActivationVectorByte < B514ActivationVectorBitFieldCount; b514ActivationVectorByte++ {
					for b514ActivationVectorBit := 0; b514ActivationVectorBit < 8; b514ActivationVectorBit++ {
						b514ActivationVectorGot := b514ActivationVectorView.B514ActivationVectorBitAt(b514ActivationVectorByte, b514ActivationVectorBit)
						if b514ActivationVectorGot.B514ActivationVectorPresent != (b514ActivationVectorByte < len(b514ActivationVectorCase.B514ActivationVector)) || b514ActivationVectorGot.B514ActivationVectorValue != b514ActivationVectorWant[[2]int{b514ActivationVectorByte, b514ActivationVectorBit}] {
							b514ActivationVectorT.Fatal(b514ActivationVectorByte, b514ActivationVectorBit, b514ActivationVectorGot)
						}
					}
				}
				return
			}
			for b514ActivationVectorIndex := 0; b514ActivationVectorIndex < B514ActivationVectorScalarFieldCount; b514ActivationVectorIndex++ {
				b514ActivationVectorGot := b514ActivationVectorView.B514ActivationVectorScalarAt(b514ActivationVectorIndex)
				b514ActivationVectorPresent := b514ActivationVectorIndex < len(b514ActivationVectorCase.B514ActivationVector)
				if b514ActivationVectorGot.B514ActivationVectorPresent != b514ActivationVectorPresent || b514ActivationVectorGot.B514ActivationVectorInterpretation != B514ActivationVectorGradeUnknown {
					b514ActivationVectorT.Fatal(b514ActivationVectorIndex, b514ActivationVectorGot)
				}
				if b514ActivationVectorPresent && b514ActivationVectorGot.B514ActivationVectorRaw != b514ActivationVectorCase.B514ActivationVector[b514ActivationVectorIndex] {
					b514ActivationVectorT.Fatal(b514ActivationVectorIndex, b514ActivationVectorGot)
				}
				if !b514ActivationVectorPresent && b514ActivationVectorGot.B514ActivationVectorRaw != 0 {
					b514ActivationVectorT.Fatal("absent carries value")
				}
			}
			if b514ActivationVectorView.B514ActivationVectorBitAt(0, 0).B514ActivationVectorPresent {
				b514ActivationVectorT.Fatal("scalars exposed as bits")
			}
		})
	}
}

func TestAb80WalkEveryBit(b514ActivationVectorT *testing.T) {
	for b514ActivationVectorByte := 0; b514ActivationVectorByte < 10; b514ActivationVectorByte++ {
		for b514ActivationVectorBit := 0; b514ActivationVectorBit < 8; b514ActivationVectorBit++ {
			b514ActivationVectorVector := make([]byte, 10)
			b514ActivationVectorVector[b514ActivationVectorByte] = 1 << uint(b514ActivationVectorBit)
			b514ActivationVectorView := B514ActivationVectorInterpret(b514ActivationVectorIdent(3), b514ActivationVectorVector)
			b514ActivationVectorCount := 0
			for b514ActivationVectorB := 0; b514ActivationVectorB < 10; b514ActivationVectorB++ {
				for b514ActivationVectorI := 0; b514ActivationVectorI < 8; b514ActivationVectorI++ {
					b514ActivationVectorGot := b514ActivationVectorView.B514ActivationVectorBitAt(b514ActivationVectorB, b514ActivationVectorI)
					if !b514ActivationVectorGot.B514ActivationVectorPresent {
						b514ActivationVectorT.Fatal("absent in full vector")
					}
					if b514ActivationVectorGot.B514ActivationVectorValue {
						b514ActivationVectorCount++
						if b514ActivationVectorB != b514ActivationVectorByte || b514ActivationVectorI != b514ActivationVectorBit {
							b514ActivationVectorT.Fatalf("set %d/%d reported at %d/%d", b514ActivationVectorByte, b514ActivationVectorBit, b514ActivationVectorB, b514ActivationVectorI)
						}
					}
				}
			}
			if b514ActivationVectorCount != 1 {
				b514ActivationVectorT.Fatal(b514ActivationVectorCount)
			}
		}
	}
}

func TestAb80EverySize(b514ActivationVectorT *testing.T) {
	for _, b514ActivationVectorOperation := range []uint8{3, 4} {
		b514ActivationVectorFull := 10
		if b514ActivationVectorOperation == 4 {
			b514ActivationVectorFull = 12
		}
		for b514ActivationVectorSize := 0; b514ActivationVectorSize <= b514ActivationVectorFull+1; b514ActivationVectorSize++ {
			b514ActivationVectorVector := bytes.Repeat([]byte{0xA5}, b514ActivationVectorSize)
			b514ActivationVectorView := B514ActivationVectorInterpret(b514ActivationVectorIdent(b514ActivationVectorOperation), b514ActivationVectorVector)
			b514ActivationVectorPresent := b514ActivationVectorSize
			if b514ActivationVectorPresent > b514ActivationVectorFull {
				b514ActivationVectorPresent = b514ActivationVectorFull
			}
			if b514ActivationVectorView.B514ActivationVectorPresentMask() != uint16(1<<uint(b514ActivationVectorPresent)-1) {
				b514ActivationVectorT.Fatal(b514ActivationVectorOperation, b514ActivationVectorSize, b514ActivationVectorView.B514ActivationVectorPresentMask())
			}
			if b514ActivationVectorView.B514ActivationVectorIsComplete() != (b514ActivationVectorSize == b514ActivationVectorFull) || len(b514ActivationVectorView.B514ActivationVectorRemainderBytes()) != b514ActivationVectorSize-b514ActivationVectorPresent {
				b514ActivationVectorT.Fatal(b514ActivationVectorOperation, b514ActivationVectorSize)
			}
		}
	}
}

func TestAb80IdentityAxes(b514ActivationVectorT *testing.T) {
	b514ActivationVectorVector := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	b514ActivationVectorBase := b514ActivationVectorIdent(3)
	b514ActivationVectorValid := B514ActivationVectorInterpret(b514ActivationVectorBase, b514ActivationVectorVector)
	if b514ActivationVectorValid.B514ActivationVectorState() != B514ActivationVectorStateBits {
		b514ActivationVectorT.Fatal("control")
	}
	b514ActivationVectorMutants := map[string]B514ActivationVectorIdentity{
		"family":       {B514ActivationVectorFamily: "other", B514ActivationVectorVariant: B514ActivationVectorVariant, B514ActivationVectorOperation: 3, B514ActivationVectorDirection: B514ActivationVectorDirectionPositive},
		"variant":      {B514ActivationVectorFamily: B514ActivationVectorFamily, B514ActivationVectorVariant: "other", B514ActivationVectorOperation: 3, B514ActivationVectorDirection: B514ActivationVectorDirectionPositive},
		"operation2":   {B514ActivationVectorFamily: B514ActivationVectorFamily, B514ActivationVectorVariant: B514ActivationVectorVariant, B514ActivationVectorOperation: 2, B514ActivationVectorDirection: B514ActivationVectorDirectionPositive},
		"operation5":   {B514ActivationVectorFamily: B514ActivationVectorFamily, B514ActivationVectorVariant: B514ActivationVectorVariant, B514ActivationVectorOperation: 5, B514ActivationVectorDirection: B514ActivationVectorDirectionPositive},
		"direction":    {B514ActivationVectorFamily: B514ActivationVectorFamily, B514ActivationVectorVariant: B514ActivationVectorVariant, B514ActivationVectorOperation: 3, B514ActivationVectorDirection: "request"},
		"directionBad": {B514ActivationVectorFamily: B514ActivationVectorFamily, B514ActivationVectorVariant: B514ActivationVectorVariant, B514ActivationVectorOperation: 4, B514ActivationVectorDirection: ""},
		"empty":        {},
	}
	for b514ActivationVectorName, b514ActivationVectorIdentity := range b514ActivationVectorMutants {
		b514ActivationVectorView := B514ActivationVectorInterpret(b514ActivationVectorIdentity, b514ActivationVectorVector)
		if b514ActivationVectorView.B514ActivationVectorState() != B514ActivationVectorStateRawOnly || b514ActivationVectorView.B514ActivationVectorGrade() != B514ActivationVectorGradeUnknown || b514ActivationVectorView.B514ActivationVectorPresentMask() != 0 || b514ActivationVectorView.B514ActivationVectorIsComplete() || b514ActivationVectorView.B514ActivationVectorBitAt(0, 0).B514ActivationVectorPresent || b514ActivationVectorView.B514ActivationVectorScalarAt(0).B514ActivationVectorPresent || !bytes.Equal(b514ActivationVectorView.B514ActivationVectorRawBytes(), b514ActivationVectorVector) {
			b514ActivationVectorT.Fatal(b514ActivationVectorName)
		}
		if b514ActivationVectorView.B514ActivationVectorFieldCount() != 0 || len(b514ActivationVectorView.B514ActivationVectorRemainderBytes()) != 0 {
			b514ActivationVectorT.Fatal(b514ActivationVectorName, "raw-only exposes fields or remainder")
		}
	}
	if B514ActivationVectorInterpret(b514ActivationVectorIdent(4), b514ActivationVectorVector).B514ActivationVectorState() != B514ActivationVectorStateScalars {
		b514ActivationVectorT.Fatal("operation 4 control")
	}
	if b514ActivationVectorValid.B514ActivationVectorFieldCount() != 10 || B514ActivationVectorInterpret(b514ActivationVectorIdent(4), b514ActivationVectorVector).B514ActivationVectorFieldCount() != 12 {
		b514ActivationVectorT.Fatal("field count control")
	}
}

func TestAb80ScalarsAreNotBits(b514ActivationVectorT *testing.T) {
	b514ActivationVectorVector := []byte{0, 3, 128, 255, 1, 2, 4, 8, 16, 32, 64, 7}
	b514ActivationVectorView := B514ActivationVectorInterpret(b514ActivationVectorIdent(4), b514ActivationVectorVector)
	for b514ActivationVectorIndex, b514ActivationVectorWant := range b514ActivationVectorVector {
		b514ActivationVectorGot := b514ActivationVectorView.B514ActivationVectorScalarAt(b514ActivationVectorIndex)
		if !b514ActivationVectorGot.B514ActivationVectorPresent || b514ActivationVectorGot.B514ActivationVectorRaw != b514ActivationVectorWant || b514ActivationVectorGot.B514ActivationVectorInterpretation != "Unknown" {
			b514ActivationVectorT.Fatal(b514ActivationVectorIndex, b514ActivationVectorGot)
		}
	}
	if b514ActivationVectorView.B514ActivationVectorScalarAt(12).B514ActivationVectorPresent || b514ActivationVectorView.B514ActivationVectorScalarAt(-1).B514ActivationVectorPresent {
		b514ActivationVectorT.Fatal("out of range present")
	}
	b514ActivationVectorBits := B514ActivationVectorInterpret(b514ActivationVectorIdent(3), b514ActivationVectorVector)
	if b514ActivationVectorBits.B514ActivationVectorScalarAt(0).B514ActivationVectorPresent || b514ActivationVectorBits.B514ActivationVectorBitAt(10, 0).B514ActivationVectorPresent || b514ActivationVectorBits.B514ActivationVectorBitAt(0, 8).B514ActivationVectorPresent {
		b514ActivationVectorT.Fatal("cross-type or out-of-range exposure")
	}
}

func TestAb80OwnedMemory(b514ActivationVectorT *testing.T) {
	b514ActivationVectorInput := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 99}
	b514ActivationVectorView := B514ActivationVectorInterpret(b514ActivationVectorIdent(3), b514ActivationVectorInput)
	b514ActivationVectorInput[0] = 0xFF
	b514ActivationVectorInput[10] = 0
	if !b514ActivationVectorView.B514ActivationVectorBitAt(0, 0).B514ActivationVectorValue || b514ActivationVectorView.B514ActivationVectorRawBytes()[0] != 1 || b514ActivationVectorView.B514ActivationVectorRemainderBytes()[0] != 99 {
		b514ActivationVectorT.Fatal("input aliased")
	}
	b514ActivationVectorView.B514ActivationVectorRawBytes()[0] = 0
	b514ActivationVectorView.B514ActivationVectorRemainderBytes()[0] = 0
	if b514ActivationVectorView.B514ActivationVectorRawBytes()[0] != 1 || b514ActivationVectorView.B514ActivationVectorRemainderBytes()[0] != 99 || !b514ActivationVectorView.B514ActivationVectorBitAt(0, 0).B514ActivationVectorValue {
		b514ActivationVectorT.Fatal("returned memory aliased")
	}
	b514ActivationVectorRaw := B514ActivationVectorInterpret(B514ActivationVectorIdentity{}, b514ActivationVectorInput)
	b514ActivationVectorInput[1] = 0x77
	if b514ActivationVectorRaw.B514ActivationVectorRawBytes()[1] != 2 {
		b514ActivationVectorT.Fatal("raw-only aliased")
	}
}

func TestAb80OwnedMemoryScalars(b514ActivationVectorT *testing.T) {
	b514ActivationVectorInput := []byte{0, 3, 128, 255, 1, 2, 4, 8, 16, 32, 64, 7, 9}
	b514ActivationVectorView := B514ActivationVectorInterpret(b514ActivationVectorIdent(4), b514ActivationVectorInput)
	b514ActivationVectorInput[3] = 0x11
	b514ActivationVectorInput[12] = 0x22
	if b514ActivationVectorView.B514ActivationVectorRawBytes()[3] != 255 || b514ActivationVectorView.B514ActivationVectorScalarAt(3).B514ActivationVectorRaw != 255 || b514ActivationVectorView.B514ActivationVectorRemainderBytes()[0] != 9 {
		b514ActivationVectorT.Fatal("scalar input aliased")
	}
	b514ActivationVectorView.B514ActivationVectorRawBytes()[3] = 0
	b514ActivationVectorView.B514ActivationVectorRemainderBytes()[0] = 0
	if b514ActivationVectorView.B514ActivationVectorRawBytes()[3] != 255 || b514ActivationVectorView.B514ActivationVectorScalarAt(3).B514ActivationVectorRaw != 255 || b514ActivationVectorView.B514ActivationVectorRemainderBytes()[0] != 9 {
		b514ActivationVectorT.Fatal("scalar returned memory aliased")
	}
}
