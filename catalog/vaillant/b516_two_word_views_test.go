package vaillant

import (
	b516TwoWordViewBinary "encoding/binary"
	b516TwoWordViewJSON "encoding/json"
	b516TwoWordViewMath "math"
	b516TwoWordViewOS "os"
	b516TwoWordViewReflect "reflect"
	b516TwoWordViewTesting "testing"
)

func TestAb17An17FiniteVectors(b516TwoWordViewT *b516TwoWordViewTesting.T) {
	b516TwoWordViewData, b516TwoWordViewError := b516TwoWordViewOS.ReadFile("testdata/ab17an17_b516_two_word_views.json")
	if b516TwoWordViewError != nil {
		b516TwoWordViewT.Fatal(b516TwoWordViewError)
	}
	var b516TwoWordViewVectors []struct {
		B516TwoWordViewLabel     string `json:"Ab17An17Label"`
		B516TwoWordViewBody      []byte `json:"Ab17An17Body"`
		B516TwoWordViewDeclared  int    `json:"Ab17An17Declared"`
		B516TwoWordViewFamily    string `json:"Ab17An17Family"`
		B516TwoWordViewVariant   string `json:"Ab17An17Variant"`
		B516TwoWordViewOperation byte   `json:"Ab17An17Operation"`
		B516TwoWordViewDirection string `json:"Ab17An17Direction"`
		B516TwoWordViewExpected  any    `json:"Ab17An17Expected"`
	}
	if b516TwoWordViewError = b516TwoWordViewJSON.Unmarshal(b516TwoWordViewData, &b516TwoWordViewVectors); b516TwoWordViewError != nil {
		b516TwoWordViewT.Fatal(b516TwoWordViewError)
	}
	if len(b516TwoWordViewVectors) != 132 {
		b516TwoWordViewT.Fatalf("ab17an17 vector count: %d", len(b516TwoWordViewVectors))
	}
	for _, b516TwoWordViewVector := range b516TwoWordViewVectors {
		b516TwoWordViewT.Run(b516TwoWordViewVector.B516TwoWordViewLabel, func(b516TwoWordViewT *b516TwoWordViewTesting.T) {
			b516TwoWordViewBody := make([]byte, len(b516TwoWordViewVector.B516TwoWordViewBody), 10)
			copy(b516TwoWordViewBody, b516TwoWordViewVector.B516TwoWordViewBody)
			b516TwoWordViewResult := B516TwoWordViewInterpret(b516TwoWordViewVector.B516TwoWordViewFamily, b516TwoWordViewVector.B516TwoWordViewVariant, b516TwoWordViewVector.B516TwoWordViewOperation, b516TwoWordViewVector.B516TwoWordViewDirection, b516TwoWordViewVector.B516TwoWordViewDeclared, b516TwoWordViewBody)
			b516TwoWordViewEncoded, b516TwoWordViewError := b516TwoWordViewJSON.Marshal(b516TwoWordViewResult)
			if b516TwoWordViewError != nil {
				b516TwoWordViewT.Fatal(b516TwoWordViewError)
			}
			var b516TwoWordViewActual any
			if b516TwoWordViewError = b516TwoWordViewJSON.Unmarshal(b516TwoWordViewEncoded, &b516TwoWordViewActual); b516TwoWordViewError != nil {
				b516TwoWordViewT.Fatal(b516TwoWordViewError)
			}
			if !b516TwoWordViewReflect.DeepEqual(b516TwoWordViewActual, b516TwoWordViewVector.B516TwoWordViewExpected) {
				b516TwoWordViewT.Fatalf("ab17an17 actual %s; expected %v", b516TwoWordViewEncoded, b516TwoWordViewVector.B516TwoWordViewExpected)
			}
		})
	}
}

func TestAb17An17GeneratedNonfinite(b516TwoWordViewT *b516TwoWordViewTesting.T) {
	for _, b516TwoWordViewBits := range []uint32{
		b516TwoWordViewMath.Float32bits(float32(b516TwoWordViewMath.NaN())),
		b516TwoWordViewMath.Float32bits(float32(b516TwoWordViewMath.Inf(1))),
		b516TwoWordViewMath.Float32bits(float32(b516TwoWordViewMath.Inf(-1))),
		uint32(255)<<23 | 123,
		uint32(511)<<23 | 456,
	} {
		b516TwoWordViewBody := []byte{0, 0, 0, 0, 0, 1, 2, 3, 4}
		b516TwoWordViewBinary.LittleEndian.PutUint32(b516TwoWordViewBody[1:5], b516TwoWordViewBits)
		b516TwoWordViewResult := B516TwoWordViewInterpret(B516TwoWordViewFamily, B516TwoWordViewVariant, 0x14, "positive", 9, b516TwoWordViewBody)
		if !b516TwoWordViewResult.B516TwoWordViewWord0Available || !b516TwoWordViewResult.B516TwoWordViewCompleteTwoWord || b516TwoWordViewResult.B516TwoWordViewWord0Finite || b516TwoWordViewResult.B516TwoWordViewWord0Value != nil || b516TwoWordViewResult.B516TwoWordViewWord0Bits != b516TwoWordViewBits || !b516TwoWordViewReflect.DeepEqual(b516TwoWordViewResult.B516TwoWordViewRaw, b516TwoWordViewBody) || !b516TwoWordViewReflect.DeepEqual(b516TwoWordViewResult.B516TwoWordViewWord0Raw[:], b516TwoWordViewBody[1:5]) {
			b516TwoWordViewT.Fatalf("ab17an17 nonfinite bits %x: %+v", b516TwoWordViewBits, b516TwoWordViewResult)
		}
	}
}

func TestAb17An17Ownership(b516TwoWordViewT *b516TwoWordViewTesting.T) {
	b516TwoWordViewBody := []byte{0, 0, 0, 128, 63, 1, 2, 3, 4, 170}
	b516TwoWordViewExpected := append([]byte(nil), b516TwoWordViewBody...)
	b516TwoWordViewResult := B516TwoWordViewInterpret(B516TwoWordViewFamily, B516TwoWordViewVariant, 0x14, "positive", 9, b516TwoWordViewBody)
	for b516TwoWordViewIndex := range b516TwoWordViewBody {
		b516TwoWordViewBody[b516TwoWordViewIndex] ^= 255
	}
	if !b516TwoWordViewReflect.DeepEqual(b516TwoWordViewResult.B516TwoWordViewRaw, b516TwoWordViewExpected) || b516TwoWordViewResult.B516TwoWordViewTrailing[0] != 170 || b516TwoWordViewResult.B516TwoWordViewWord0Raw != [4]byte{0, 0, 128, 63} || b516TwoWordViewResult.B516TwoWordViewWord1Raw != [4]byte{1, 2, 3, 4} || *b516TwoWordViewResult.B516TwoWordViewWord0Value != 1 {
		b516TwoWordViewT.Fatal("ab17an17 input mutation changed result")
	}
	b516TwoWordViewSecond := B516TwoWordViewInterpret(B516TwoWordViewFamily, B516TwoWordViewVariant, 0x14, "positive", 9, b516TwoWordViewExpected)
	b516TwoWordViewResult.B516TwoWordViewRaw[9] = 0
	if b516TwoWordViewResult.B516TwoWordViewTrailing[0] != 170 {
		b516TwoWordViewT.Fatal("ab17an17 trailing aliases raw")
	}
	b516TwoWordViewResult.B516TwoWordViewRaw[1] = 255
	b516TwoWordViewResult.B516TwoWordViewTrailing[0] = 0
	b516TwoWordViewResult.B516TwoWordViewWord0Raw[0] = 255
	b516TwoWordViewResult.B516TwoWordViewWord1Raw[0] = 255
	*b516TwoWordViewResult.B516TwoWordViewWord0Value = -9
	if !b516TwoWordViewReflect.DeepEqual(b516TwoWordViewSecond.B516TwoWordViewRaw, b516TwoWordViewExpected) || !b516TwoWordViewReflect.DeepEqual(b516TwoWordViewExpected, []byte{0, 0, 0, 128, 63, 1, 2, 3, 4, 170}) || b516TwoWordViewSecond.B516TwoWordViewTrailing[0] != 170 || b516TwoWordViewSecond.B516TwoWordViewWord0Raw != [4]byte{0, 0, 128, 63} || b516TwoWordViewSecond.B516TwoWordViewWord1Raw != [4]byte{1, 2, 3, 4} || *b516TwoWordViewSecond.B516TwoWordViewWord0Value != 1 {
		b516TwoWordViewT.Fatal("ab17an17 returned mutation leaked")
	}
}
