package vaillant

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
)

type ab80Case struct {
	Ab80ID        string   `json:"ab80_id"`
	Ab80Operation uint8    `json:"ab80_operation"`
	Ab80Vector    []byte   `json:"ab80_vector"`
	Ab80State     string   `json:"ab80_state"`
	Ab80Mask      uint16   `json:"ab80_mask"`
	Ab80Complete  bool     `json:"ab80_complete"`
	Ab80Remainder []byte   `json:"ab80_remainder"`
	Ab80True      [][2]int `json:"ab80_true"`
}

func ab80Ident(ab80Operation uint8) Ab80Identity {
	return Ab80Identity{Ab80Family: Ab80Family, Ab80Variant: Ab80Variant, Ab80Operation: ab80Operation, Ab80Direction: Ab80DirectionPositive}
}

func ab80Load(ab80T *testing.T) []ab80Case {
	ab80T.Helper()
	ab80Data, ab80Err := os.ReadFile("testdata/ab80_b514_activation_vectors.json")
	if ab80Err != nil {
		ab80T.Fatal(ab80Err)
	}
	var ab80Cases []ab80Case
	if ab80Err = json.Unmarshal(ab80Data, &ab80Cases); ab80Err != nil {
		ab80T.Fatal(ab80Err)
	}
	if len(ab80Cases) == 0 {
		ab80T.Fatal("no fixture cases")
	}
	return ab80Cases
}

func TestAb80Fixtures(ab80T *testing.T) {
	for _, ab80Case := range ab80Load(ab80T) {
		ab80T.Run(ab80Case.Ab80ID, func(ab80T *testing.T) {
			ab80View := Ab80Interpret(ab80Ident(ab80Case.Ab80Operation), ab80Case.Ab80Vector)
			if ab80View.Ab80State() != ab80Case.Ab80State || ab80View.Ab80Grade() != Ab80GradeHypothesis {
				ab80T.Fatal(ab80View.Ab80State(), ab80View.Ab80Grade())
			}
			if ab80View.Ab80PresentMask() != ab80Case.Ab80Mask || ab80View.Ab80IsComplete() != ab80Case.Ab80Complete {
				ab80T.Fatal(ab80View.Ab80PresentMask(), ab80View.Ab80IsComplete())
			}
			if !bytes.Equal(ab80View.Ab80RawBytes(), ab80Case.Ab80Vector) || !bytes.Equal(ab80View.Ab80RemainderBytes(), ab80Case.Ab80Remainder) {
				ab80T.Fatal("raw or remainder changed")
			}
			if ab80Case.Ab80Operation == Ab80OperationBits {
				ab80Want := map[[2]int]bool{}
				for _, ab80Pos := range ab80Case.Ab80True {
					ab80Want[ab80Pos] = true
				}
				for ab80Byte := 0; ab80Byte < Ab80BitFieldCount; ab80Byte++ {
					for ab80Bit := 0; ab80Bit < 8; ab80Bit++ {
						ab80Got := ab80View.Ab80BitAt(ab80Byte, ab80Bit)
						if ab80Got.Ab80Present != (ab80Byte < len(ab80Case.Ab80Vector)) || ab80Got.Ab80Value != ab80Want[[2]int{ab80Byte, ab80Bit}] {
							ab80T.Fatal(ab80Byte, ab80Bit, ab80Got)
						}
					}
				}
				return
			}
			for ab80Index := 0; ab80Index < Ab80ScalarFieldCount; ab80Index++ {
				ab80Got := ab80View.Ab80ScalarAt(ab80Index)
				ab80Present := ab80Index < len(ab80Case.Ab80Vector)
				if ab80Got.Ab80Present != ab80Present || ab80Got.Ab80Interpretation != Ab80GradeUnknown {
					ab80T.Fatal(ab80Index, ab80Got)
				}
				if ab80Present && ab80Got.Ab80Raw != ab80Case.Ab80Vector[ab80Index] {
					ab80T.Fatal(ab80Index, ab80Got)
				}
				if !ab80Present && ab80Got.Ab80Raw != 0 {
					ab80T.Fatal("absent carries value")
				}
			}
			if ab80View.Ab80BitAt(0, 0).Ab80Present {
				ab80T.Fatal("scalars exposed as bits")
			}
		})
	}
}

func TestAb80WalkEveryBit(ab80T *testing.T) {
	for ab80Byte := 0; ab80Byte < 10; ab80Byte++ {
		for ab80Bit := 0; ab80Bit < 8; ab80Bit++ {
			ab80Vector := make([]byte, 10)
			ab80Vector[ab80Byte] = 1 << uint(ab80Bit)
			ab80View := Ab80Interpret(ab80Ident(3), ab80Vector)
			ab80Count := 0
			for ab80B := 0; ab80B < 10; ab80B++ {
				for ab80I := 0; ab80I < 8; ab80I++ {
					ab80Got := ab80View.Ab80BitAt(ab80B, ab80I)
					if !ab80Got.Ab80Present {
						ab80T.Fatal("absent in full vector")
					}
					if ab80Got.Ab80Value {
						ab80Count++
						if ab80B != ab80Byte || ab80I != ab80Bit {
							ab80T.Fatalf("set %d/%d reported at %d/%d", ab80Byte, ab80Bit, ab80B, ab80I)
						}
					}
				}
			}
			if ab80Count != 1 {
				ab80T.Fatal(ab80Count)
			}
		}
	}
}

func TestAb80EverySize(ab80T *testing.T) {
	for _, ab80Operation := range []uint8{3, 4} {
		ab80Full := 10
		if ab80Operation == 4 {
			ab80Full = 12
		}
		for ab80Size := 0; ab80Size <= ab80Full+1; ab80Size++ {
			ab80Vector := bytes.Repeat([]byte{0xA5}, ab80Size)
			ab80View := Ab80Interpret(ab80Ident(ab80Operation), ab80Vector)
			ab80Present := ab80Size
			if ab80Present > ab80Full {
				ab80Present = ab80Full
			}
			if ab80View.Ab80PresentMask() != uint16(1<<uint(ab80Present)-1) {
				ab80T.Fatal(ab80Operation, ab80Size, ab80View.Ab80PresentMask())
			}
			if ab80View.Ab80IsComplete() != (ab80Size == ab80Full) || len(ab80View.Ab80RemainderBytes()) != ab80Size-ab80Present {
				ab80T.Fatal(ab80Operation, ab80Size)
			}
		}
	}
}

func TestAb80IdentityAxes(ab80T *testing.T) {
	ab80Vector := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	ab80Base := ab80Ident(3)
	ab80Valid := Ab80Interpret(ab80Base, ab80Vector)
	if ab80Valid.Ab80State() != Ab80StateBits {
		ab80T.Fatal("control")
	}
	ab80Mutants := map[string]Ab80Identity{
		"family":       {Ab80Family: "other", Ab80Variant: Ab80Variant, Ab80Operation: 3, Ab80Direction: Ab80DirectionPositive},
		"variant":      {Ab80Family: Ab80Family, Ab80Variant: "other", Ab80Operation: 3, Ab80Direction: Ab80DirectionPositive},
		"operation2":   {Ab80Family: Ab80Family, Ab80Variant: Ab80Variant, Ab80Operation: 2, Ab80Direction: Ab80DirectionPositive},
		"operation5":   {Ab80Family: Ab80Family, Ab80Variant: Ab80Variant, Ab80Operation: 5, Ab80Direction: Ab80DirectionPositive},
		"direction":    {Ab80Family: Ab80Family, Ab80Variant: Ab80Variant, Ab80Operation: 3, Ab80Direction: "request"},
		"directionBad": {Ab80Family: Ab80Family, Ab80Variant: Ab80Variant, Ab80Operation: 4, Ab80Direction: ""},
		"empty":        {},
	}
	for ab80Name, ab80Identity := range ab80Mutants {
		ab80View := Ab80Interpret(ab80Identity, ab80Vector)
		if ab80View.Ab80State() != Ab80StateRawOnly || ab80View.Ab80Grade() != Ab80GradeUnknown || ab80View.Ab80PresentMask() != 0 || ab80View.Ab80IsComplete() || ab80View.Ab80BitAt(0, 0).Ab80Present || ab80View.Ab80ScalarAt(0).Ab80Present || !bytes.Equal(ab80View.Ab80RawBytes(), ab80Vector) {
			ab80T.Fatal(ab80Name)
		}
		if ab80View.Ab80FieldCount() != 0 || len(ab80View.Ab80RemainderBytes()) != 0 {
			ab80T.Fatal(ab80Name, "raw-only exposes fields or remainder")
		}
	}
	if Ab80Interpret(ab80Ident(4), ab80Vector).Ab80State() != Ab80StateScalars {
		ab80T.Fatal("operation 4 control")
	}
	if ab80Valid.Ab80FieldCount() != 10 || Ab80Interpret(ab80Ident(4), ab80Vector).Ab80FieldCount() != 12 {
		ab80T.Fatal("field count control")
	}
}

func TestAb80ScalarsAreNotBits(ab80T *testing.T) {
	ab80Vector := []byte{0, 3, 128, 255, 1, 2, 4, 8, 16, 32, 64, 7}
	ab80View := Ab80Interpret(ab80Ident(4), ab80Vector)
	for ab80Index, ab80Want := range ab80Vector {
		ab80Got := ab80View.Ab80ScalarAt(ab80Index)
		if !ab80Got.Ab80Present || ab80Got.Ab80Raw != ab80Want || ab80Got.Ab80Interpretation != "Unknown" {
			ab80T.Fatal(ab80Index, ab80Got)
		}
	}
	if ab80View.Ab80ScalarAt(12).Ab80Present || ab80View.Ab80ScalarAt(-1).Ab80Present {
		ab80T.Fatal("out of range present")
	}
	ab80Bits := Ab80Interpret(ab80Ident(3), ab80Vector)
	if ab80Bits.Ab80ScalarAt(0).Ab80Present || ab80Bits.Ab80BitAt(10, 0).Ab80Present || ab80Bits.Ab80BitAt(0, 8).Ab80Present {
		ab80T.Fatal("cross-type or out-of-range exposure")
	}
}

func TestAb80OwnedMemory(ab80T *testing.T) {
	ab80Input := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 99}
	ab80View := Ab80Interpret(ab80Ident(3), ab80Input)
	ab80Input[0] = 0xFF
	ab80Input[10] = 0
	if !ab80View.Ab80BitAt(0, 0).Ab80Value || ab80View.Ab80RawBytes()[0] != 1 || ab80View.Ab80RemainderBytes()[0] != 99 {
		ab80T.Fatal("input aliased")
	}
	ab80View.Ab80RawBytes()[0] = 0
	ab80View.Ab80RemainderBytes()[0] = 0
	if ab80View.Ab80RawBytes()[0] != 1 || ab80View.Ab80RemainderBytes()[0] != 99 || !ab80View.Ab80BitAt(0, 0).Ab80Value {
		ab80T.Fatal("returned memory aliased")
	}
	ab80Raw := Ab80Interpret(Ab80Identity{}, ab80Input)
	ab80Input[1] = 0x77
	if ab80Raw.Ab80RawBytes()[1] != 2 {
		ab80T.Fatal("raw-only aliased")
	}
}

func TestAb80OwnedMemoryScalars(ab80T *testing.T) {
	ab80Input := []byte{0, 3, 128, 255, 1, 2, 4, 8, 16, 32, 64, 7, 9}
	ab80View := Ab80Interpret(ab80Ident(4), ab80Input)
	ab80Input[3] = 0x11
	ab80Input[12] = 0x22
	if ab80View.Ab80RawBytes()[3] != 255 || ab80View.Ab80ScalarAt(3).Ab80Raw != 255 || ab80View.Ab80RemainderBytes()[0] != 9 {
		ab80T.Fatal("scalar input aliased")
	}
	ab80View.Ab80RawBytes()[3] = 0
	ab80View.Ab80RemainderBytes()[0] = 0
	if ab80View.Ab80RawBytes()[3] != 255 || ab80View.Ab80ScalarAt(3).Ab80Raw != 255 || ab80View.Ab80RemainderBytes()[0] != 9 {
		ab80T.Fatal("scalar returned memory aliased")
	}
}
