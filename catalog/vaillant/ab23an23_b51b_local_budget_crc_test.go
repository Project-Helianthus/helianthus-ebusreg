package vaillant

import (
	ab23an23Bytes "bytes"
	ab23an23JSON "encoding/json"
	ab23an23OS "os"
	ab23an23Testing "testing"
)

type ab23an23Vectors struct {
	Ab23An23Budgets []struct {
		Ab23An23Label     string
		Ab23An23Procedure Ab23An23Procedure
		Ab23An23Length    int
		Ab23An23Chunks    []int
		Ab23An23Consumed  int
		Ab23An23Unplanned int
	}
	Ab23An23CRCs []struct {
		Ab23An23Label    string
		Ab23An23Bytes    []byte
		Ab23An23Initial  uint16
		Ab23An23Expected uint16
	}
}

func ab23an23Load(ab23an23T *ab23an23Testing.T) ab23an23Vectors {
	ab23an23T.Helper()
	ab23an23Data, ab23an23Err := ab23an23OS.ReadFile("testdata/ab23an23_b51b_local_budget_crc.json")
	if ab23an23Err != nil {
		ab23an23T.Fatal(ab23an23Err)
	}
	var ab23an23Fixture ab23an23Vectors
	if ab23an23Err = ab23an23JSON.Unmarshal(ab23an23Data, &ab23an23Fixture); ab23an23Err != nil {
		ab23an23T.Fatal(ab23an23Err)
	}
	if len(ab23an23Fixture.Ab23An23Budgets) != 11 || len(ab23an23Fixture.Ab23An23CRCs) != 16 {
		ab23an23T.Fatal("ab23an23: incomplete vectors")
	}
	return ab23an23Fixture
}

func TestAb23An23Budgets(ab23an23T *ab23an23Testing.T) {
	for _, ab23an23Vector := range ab23an23Load(ab23an23T).Ab23An23Budgets {
		ab23an23T.Run(ab23an23Vector.Ab23An23Label, func(ab23an23T *ab23an23Testing.T) {
			ab23an23Input := make([]byte, ab23an23Vector.Ab23An23Length)
			for ab23an23Index := range ab23an23Input {
				ab23an23Input[ab23an23Index] = byte(ab23an23Index % 251)
			}
			ab23an23Original := append([]byte(nil), ab23an23Input...)
			ab23an23Result, ab23an23Err := Ab23An23Plan(ab23an23Input, ab23an23Vector.Ab23An23Procedure)
			if ab23an23Err != nil {
				ab23an23T.Fatal(ab23an23Err)
			}
			if ab23an23Result.Ab23An23Procedure != ab23an23Vector.Ab23An23Procedure || ab23an23Result.Ab23An23Consumed != ab23an23Vector.Ab23An23Consumed || ab23an23Result.Ab23An23Unplanned != ab23an23Vector.Ab23An23Unplanned || len(ab23an23Result.Ab23An23Chunks) != len(ab23an23Vector.Ab23An23Chunks) {
				ab23an23T.Fatalf("ab23an23: wrong budget: %+v", ab23an23Result)
			}
			ab23an23Offset := 0
			for ab23an23Index, ab23an23Chunk := range ab23an23Result.Ab23An23Chunks {
				ab23an23Length := ab23an23Vector.Ab23An23Chunks[ab23an23Index]
				if len(ab23an23Chunk) != ab23an23Length || !ab23an23Bytes.Equal(ab23an23Chunk, ab23an23Original[ab23an23Offset:ab23an23Offset+ab23an23Length]) {
					ab23an23T.Fatal("ab23an23: wrong chunk bytes or length")
				}
				ab23an23Offset += ab23an23Length
			}
			for _, ab23an23Chunk := range ab23an23Result.Ab23An23Chunks {
				for ab23an23Index := range ab23an23Chunk {
					ab23an23Chunk[ab23an23Index] ^= 255
				}
			}
			if !ab23an23Bytes.Equal(ab23an23Input, ab23an23Original) {
				ab23an23T.Fatal("ab23an23: output aliases input")
			}
			ab23an23Again, ab23an23Err := Ab23An23Plan(ab23an23Input, ab23an23Vector.Ab23An23Procedure)
			if ab23an23Err != nil || !ab23an23Bytes.Equal(ab23an23Bytes.Join(ab23an23Again.Ab23An23Chunks, nil), ab23an23Original[:ab23an23Vector.Ab23An23Consumed]) {
				ab23an23T.Fatal("ab23an23: later call changed")
			}
			for ab23an23Index := range ab23an23Input {
				ab23an23Input[ab23an23Index] ^= 1
			}
			if !ab23an23Bytes.Equal(ab23an23Bytes.Join(ab23an23Again.Ab23An23Chunks, nil), ab23an23Original[:ab23an23Vector.Ab23An23Consumed]) {
				ab23an23T.Fatal("ab23an23: input aliases output")
			}
		})
	}
	ab23an23Result, ab23an23Err := Ab23An23Plan([]byte{1, 2, 3}, Ab23An23Procedure("ab23an23_unknown"))
	if ab23an23Err == nil || ab23an23Result.Ab23An23Procedure != "ab23an23_unknown" || ab23an23Result.Ab23An23Consumed != 0 || ab23an23Result.Ab23An23Unplanned != 3 || len(ab23an23Result.Ab23An23Chunks) != 0 {
		ab23an23T.Fatal("ab23an23: unknown choice admitted")
	}
}

func TestAb23An23Progress(ab23an23T *ab23an23Testing.T) {
	ab23an23Max := int(^uint(0) >> 1)
	for _, ab23an23Case := range []struct {
		ab23an23Prior    int
		ab23an23Accepted bool
		ab23an23Actual   int
		ab23an23Want     int
		ab23an23Error    bool
	}{
		{10, true, 17, 27, false}, {10, false, 17, 10, false}, {11, true, 17, 28, false}, {10, true, 18, 28, false},
		{10, true, 0, 10, false}, {0, true, 17, 17, false}, {ab23an23Max, true, 0, ab23an23Max, false},
		{ab23an23Max - 17, true, 17, ab23an23Max, false}, {ab23an23Max, true, 1, ab23an23Max, true},
		{ab23an23Max, false, 17, ab23an23Max, false}, {-1, true, 17, -1, true}, {10, true, -1, 10, true},
		{-1, false, 17, -1, true}, {10, false, -1, 10, true},
	} {
		ab23an23Got, ab23an23Err := Ab23An23Progress(ab23an23Case.ab23an23Prior, ab23an23Case.ab23an23Accepted, ab23an23Case.ab23an23Actual)
		if ab23an23Got != ab23an23Case.ab23an23Want || (ab23an23Err != nil) != ab23an23Case.ab23an23Error {
			ab23an23T.Fatalf("ab23an23: transition %+v got %d, %v", ab23an23Case, ab23an23Got, ab23an23Err)
		}
	}
}

func TestAb23An23CRC(ab23an23T *ab23an23Testing.T) {
	for _, ab23an23Vector := range ab23an23Load(ab23an23T).Ab23An23CRCs {
		ab23an23T.Run(ab23an23Vector.Ab23An23Label, func(ab23an23T *ab23an23Testing.T) {
			ab23an23Original := append([]byte(nil), ab23an23Vector.Ab23An23Bytes...)
			if ab23an23Got := Ab23An23CRCWithInitial(ab23an23Vector.Ab23An23Bytes, ab23an23Vector.Ab23An23Initial); ab23an23Got != ab23an23Vector.Ab23An23Expected {
				ab23an23T.Fatalf("ab23an23: CRC got %04X want %04X", ab23an23Got, ab23an23Vector.Ab23An23Expected)
			}
			if ab23an23Vector.Ab23An23Initial == 65535 && Ab23An23CRC(ab23an23Vector.Ab23An23Bytes) != ab23an23Vector.Ab23An23Expected {
				ab23an23T.Fatal("ab23an23: wrong default")
			}
			if !ab23an23Bytes.Equal(ab23an23Vector.Ab23An23Bytes, ab23an23Original) {
				ab23an23T.Fatal("ab23an23: CRC mutated input")
			}
		})
	}
	ab23an23Span := []byte("x123456789y")
	if Ab23An23CRC(ab23an23Span[1:10]) != 0x29b1 || Ab23An23CRC(ab23an23Span[1:10:11]) != 0x29b1 || Ab23An23CRC(nil) != 0xffff || Ab23An23CRCWithInitial([]byte("123456789"), 0) != 0x31c3 {
		ab23an23T.Fatal("ab23an23: caller span or canonical CRC mismatch")
	}
}
