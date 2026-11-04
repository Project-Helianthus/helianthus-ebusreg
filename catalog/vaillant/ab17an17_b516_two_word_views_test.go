package vaillant

import (
	ab17an17Binary "encoding/binary"
	ab17an17JSON "encoding/json"
	ab17an17Math "math"
	ab17an17OS "os"
	ab17an17Reflect "reflect"
	ab17an17Testing "testing"
)

func TestAb17An17FiniteVectors(ab17an17T *ab17an17Testing.T) {
	ab17an17Data, ab17an17Error := ab17an17OS.ReadFile("testdata/ab17an17_b516_two_word_views.json")
	if ab17an17Error != nil {
		ab17an17T.Fatal(ab17an17Error)
	}
	var ab17an17Vectors []struct {
		Ab17An17Label     string
		Ab17An17Body      []byte
		Ab17An17Declared  int
		Ab17An17Family    string
		Ab17An17Variant   string
		Ab17An17Operation byte
		Ab17An17Direction string
		Ab17An17Expected  any
	}
	if ab17an17Error = ab17an17JSON.Unmarshal(ab17an17Data, &ab17an17Vectors); ab17an17Error != nil {
		ab17an17T.Fatal(ab17an17Error)
	}
	if len(ab17an17Vectors) != 132 {
		ab17an17T.Fatalf("ab17an17 vector count: %d", len(ab17an17Vectors))
	}
	for _, ab17an17Vector := range ab17an17Vectors {
		ab17an17T.Run(ab17an17Vector.Ab17An17Label, func(ab17an17T *ab17an17Testing.T) {
			ab17an17Body := make([]byte, len(ab17an17Vector.Ab17An17Body), 10)
			copy(ab17an17Body, ab17an17Vector.Ab17An17Body)
			ab17an17Result := Ab17An17Interpret(ab17an17Vector.Ab17An17Family, ab17an17Vector.Ab17An17Variant, ab17an17Vector.Ab17An17Operation, ab17an17Vector.Ab17An17Direction, ab17an17Vector.Ab17An17Declared, ab17an17Body)
			ab17an17Encoded, ab17an17Error := ab17an17JSON.Marshal(ab17an17Result)
			if ab17an17Error != nil {
				ab17an17T.Fatal(ab17an17Error)
			}
			var ab17an17Actual any
			if ab17an17Error = ab17an17JSON.Unmarshal(ab17an17Encoded, &ab17an17Actual); ab17an17Error != nil {
				ab17an17T.Fatal(ab17an17Error)
			}
			if !ab17an17Reflect.DeepEqual(ab17an17Actual, ab17an17Vector.Ab17An17Expected) {
				ab17an17T.Fatalf("ab17an17 actual %s; expected %v", ab17an17Encoded, ab17an17Vector.Ab17An17Expected)
			}
		})
	}
}

func TestAb17An17GeneratedNonfinite(ab17an17T *ab17an17Testing.T) {
	for _, ab17an17Bits := range []uint32{
		ab17an17Math.Float32bits(float32(ab17an17Math.NaN())),
		ab17an17Math.Float32bits(float32(ab17an17Math.Inf(1))),
		ab17an17Math.Float32bits(float32(ab17an17Math.Inf(-1))),
		uint32(255)<<23 | 123,
		uint32(511)<<23 | 456,
	} {
		ab17an17Body := []byte{0, 0, 0, 0, 0, 1, 2, 3, 4}
		ab17an17Binary.LittleEndian.PutUint32(ab17an17Body[1:5], ab17an17Bits)
		ab17an17Result := Ab17An17Interpret(Ab17An17Family, Ab17An17Variant, 0x14, "positive", 9, ab17an17Body)
		if !ab17an17Result.Ab17An17Word0Available || !ab17an17Result.Ab17An17CompleteTwoWord || ab17an17Result.Ab17An17Word0Finite || ab17an17Result.Ab17An17Word0Value != nil || ab17an17Result.Ab17An17Word0Bits != ab17an17Bits || !ab17an17Reflect.DeepEqual(ab17an17Result.Ab17An17Raw, ab17an17Body) || !ab17an17Reflect.DeepEqual(ab17an17Result.Ab17An17Word0Raw[:], ab17an17Body[1:5]) {
			ab17an17T.Fatalf("ab17an17 nonfinite bits %x: %+v", ab17an17Bits, ab17an17Result)
		}
	}
}

func TestAb17An17Ownership(ab17an17T *ab17an17Testing.T) {
	ab17an17Body := []byte{0, 0, 0, 128, 63, 1, 2, 3, 4, 170}
	ab17an17Expected := append([]byte(nil), ab17an17Body...)
	ab17an17Result := Ab17An17Interpret(Ab17An17Family, Ab17An17Variant, 0x14, "positive", 9, ab17an17Body)
	for ab17an17Index := range ab17an17Body {
		ab17an17Body[ab17an17Index] ^= 255
	}
	if !ab17an17Reflect.DeepEqual(ab17an17Result.Ab17An17Raw, ab17an17Expected) || ab17an17Result.Ab17An17Trailing[0] != 170 || ab17an17Result.Ab17An17Word0Raw != [4]byte{0, 0, 128, 63} || ab17an17Result.Ab17An17Word1Raw != [4]byte{1, 2, 3, 4} || *ab17an17Result.Ab17An17Word0Value != 1 {
		ab17an17T.Fatal("ab17an17 input mutation changed result")
	}
	ab17an17Second := Ab17An17Interpret(Ab17An17Family, Ab17An17Variant, 0x14, "positive", 9, ab17an17Expected)
	ab17an17Result.Ab17An17Raw[9] = 0
	if ab17an17Result.Ab17An17Trailing[0] != 170 {
		ab17an17T.Fatal("ab17an17 trailing aliases raw")
	}
	ab17an17Result.Ab17An17Raw[1] = 255
	ab17an17Result.Ab17An17Trailing[0] = 0
	ab17an17Result.Ab17An17Word0Raw[0] = 255
	ab17an17Result.Ab17An17Word1Raw[0] = 255
	*ab17an17Result.Ab17An17Word0Value = -9
	if !ab17an17Reflect.DeepEqual(ab17an17Second.Ab17An17Raw, ab17an17Expected) || !ab17an17Reflect.DeepEqual(ab17an17Expected, []byte{0, 0, 0, 128, 63, 1, 2, 3, 4, 170}) || ab17an17Second.Ab17An17Trailing[0] != 170 || ab17an17Second.Ab17An17Word0Raw != [4]byte{0, 0, 128, 63} || ab17an17Second.Ab17An17Word1Raw != [4]byte{1, 2, 3, 4} || *ab17an17Second.Ab17An17Word0Value != 1 {
		ab17an17T.Fatal("ab17an17 returned mutation leaked")
	}
}
