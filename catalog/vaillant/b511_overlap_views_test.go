package vaillant

import (
	b511OverlapViewJSON "encoding/json"
	b511OverlapViewOs "os"
	b511OverlapViewReflect "reflect"
	b511OverlapViewTesting "testing"
)

func TestAb11An11IndependentVectors(b511OverlapViewT *b511OverlapViewTesting.T) {
	b511OverlapViewData, b511OverlapViewErr := b511OverlapViewOs.ReadFile("testdata/ab11an11_b511_overlap_views.json")
	if b511OverlapViewErr != nil {
		b511OverlapViewT.Fatal(b511OverlapViewErr)
	}
	var b511OverlapViewCases []struct {
		B511OverlapViewLabel     string                   `json:"Ab11An11Label"`
		B511OverlapViewSelection B511OverlapViewSelection `json:"Ab11An11Selection"`
		B511OverlapViewImage     []byte                   `json:"Ab11An11Image"`
		B511OverlapViewExpected  B511OverlapView          `json:"Ab11An11Expected"`
	}
	if b511OverlapViewErr := b511OverlapViewJSON.Unmarshal(b511OverlapViewData, &b511OverlapViewCases); b511OverlapViewErr != nil {
		b511OverlapViewT.Fatal(b511OverlapViewErr)
	}
	if len(b511OverlapViewCases) != 11 {
		b511OverlapViewT.Fatalf("ab11an11 expected 11 independent vectors, got %d", len(b511OverlapViewCases))
	}
	for _, b511OverlapViewCase := range b511OverlapViewCases {
		b511OverlapViewT.Run(b511OverlapViewCase.B511OverlapViewLabel, func(b511OverlapViewT *b511OverlapViewTesting.T) {
			b511OverlapViewGot := B511OverlapViewInterpret(b511OverlapViewCase.B511OverlapViewSelection, b511OverlapViewCase.B511OverlapViewImage)
			if !b511OverlapViewReflect.DeepEqual(b511OverlapViewGot, b511OverlapViewCase.B511OverlapViewExpected) {
				b511OverlapViewT.Fatalf("ab11an11 got %+v; want %+v", b511OverlapViewGot, b511OverlapViewCase.B511OverlapViewExpected)
			}
		})
	}
}

func TestAb11An11SelectionAndExtent(b511OverlapViewT *b511OverlapViewTesting.T) {
	for _, b511OverlapViewVariant := range []string{B511OverlapViewOP00Reserved, B511OverlapViewOP00Fan, B511OverlapViewOP07Base, B511OverlapViewOP07ExtendedRoom} {
		b511OverlapViewOperation, b511OverlapViewRequired := byte(0), 12
		if b511OverlapViewVariant == B511OverlapViewOP07Base {
			b511OverlapViewOperation, b511OverlapViewRequired = 7, 8
		}
		if b511OverlapViewVariant == B511OverlapViewOP07ExtendedRoom {
			b511OverlapViewOperation, b511OverlapViewRequired = 7, 9
		}
		for b511OverlapViewLength := 0; b511OverlapViewLength <= 15; b511OverlapViewLength++ {
			b511OverlapViewImage := make([]byte, b511OverlapViewLength)
			for b511OverlapViewIndex := range b511OverlapViewImage {
				b511OverlapViewImage[b511OverlapViewIndex] = byte(b511OverlapViewIndex + 121)
			}
			b511OverlapViewSelection := B511OverlapViewSelection{B511OverlapViewFamily, b511OverlapViewVariant, &b511OverlapViewOperation, "positive"}
			b511OverlapViewGot := B511OverlapViewInterpret(b511OverlapViewSelection, b511OverlapViewImage)
			if !b511OverlapViewGot.B511OverlapViewSelected || b511OverlapViewGot.B511OverlapViewRequiredExtent != b511OverlapViewRequired || b511OverlapViewGot.B511OverlapViewCompleteLogicalExtent != (b511OverlapViewLength >= b511OverlapViewRequired) || !b511OverlapViewReflect.DeepEqual(b511OverlapViewGot.B511OverlapViewSelection, b511OverlapViewSelection) || b511OverlapViewGot.B511OverlapViewAvailableExtent != b511OverlapViewLength {
				b511OverlapViewT.Fatalf("ab11an11 extent/selection: %+v", b511OverlapViewGot)
			}
			if (b511OverlapViewGot.B511OverlapViewRawFlag != nil) != (b511OverlapViewOperation == 0 && b511OverlapViewLength > 8) || (b511OverlapViewGot.B511OverlapViewRawOverlay != nil) != (b511OverlapViewOperation == 7 && b511OverlapViewLength > 7) || (b511OverlapViewGot.B511OverlapViewRawPressure != nil) != (b511OverlapViewVariant == B511OverlapViewOP07ExtendedRoom && b511OverlapViewLength > 8) {
				b511OverlapViewT.Fatalf("ab11an11 availability: %+v", b511OverlapViewGot)
			}
			for b511OverlapViewOtherOp := 0; b511OverlapViewOtherOp < 256; b511OverlapViewOtherOp++ {
				if byte(b511OverlapViewOtherOp) == b511OverlapViewOperation {
					continue
				}
				b511OverlapViewWrong := b511OverlapViewSelection
				b511OverlapViewOtherByte := byte(b511OverlapViewOtherOp)
				b511OverlapViewWrong.B511OverlapViewOperation = &b511OverlapViewOtherByte
				b511OverlapViewAssertUnknown(b511OverlapViewT, b511OverlapViewWrong, b511OverlapViewImage)
			}
			b511OverlapViewWrong := b511OverlapViewSelection
			b511OverlapViewWrong.B511OverlapViewFamily = "ab11an11-other-family"
			b511OverlapViewAssertUnknown(b511OverlapViewT, b511OverlapViewWrong, b511OverlapViewImage)
			b511OverlapViewWrong = b511OverlapViewSelection
			b511OverlapViewWrong.B511OverlapViewVariant = "ab11an11-unknown-variant"
			b511OverlapViewAssertUnknown(b511OverlapViewT, b511OverlapViewWrong, b511OverlapViewImage)
			for _, b511OverlapViewDirection := range []string{"", "request", "negative", "ab11an11-unknown-direction"} {
				b511OverlapViewWrong = b511OverlapViewSelection
				b511OverlapViewWrong.B511OverlapViewDirection = b511OverlapViewDirection
				b511OverlapViewAssertUnknown(b511OverlapViewT, b511OverlapViewWrong, b511OverlapViewImage)
			}
		}
	}
}

func b511OverlapViewOpPtr(b511OverlapViewValue byte) *byte { return &b511OverlapViewValue }

func TestAb11An11AbsentOperationFailsClosed(b511OverlapViewT *b511OverlapViewTesting.T) {
	b511OverlapViewImage := make([]byte, 12)
	for _, b511OverlapViewVariant := range []string{B511OverlapViewOP00Reserved, B511OverlapViewOP00Fan, B511OverlapViewOP07Base, B511OverlapViewOP07ExtendedRoom} {
		b511OverlapViewAssertUnknown(b511OverlapViewT, B511OverlapViewSelection{B511OverlapViewFamily: B511OverlapViewFamily, B511OverlapViewVariant: b511OverlapViewVariant, B511OverlapViewDirection: "positive"}, b511OverlapViewImage)
	}
}

func b511OverlapViewAssertUnknown(b511OverlapViewT *b511OverlapViewTesting.T, b511OverlapViewSelection B511OverlapViewSelection, b511OverlapViewImage []byte) {
	b511OverlapViewT.Helper()
	b511OverlapViewExpected := B511OverlapView{B511OverlapViewSelection: b511OverlapViewSelection, B511OverlapViewRaw: append([]byte(nil), b511OverlapViewImage...), B511OverlapViewAvailableExtent: len(b511OverlapViewImage), B511OverlapViewPressureDomain: "Unknown", B511OverlapViewCandidateGrade: "Unknown", B511OverlapViewPhysicalQualification: "Unknown"}
	if b511OverlapViewGot := B511OverlapViewInterpret(b511OverlapViewSelection, b511OverlapViewImage); !b511OverlapViewReflect.DeepEqual(b511OverlapViewGot, b511OverlapViewExpected) {
		b511OverlapViewT.Fatalf("ab11an11 unknown selection invented fields: %+v", b511OverlapViewGot)
	}
}

func TestAb11An11EveryBitAndOverlap(b511OverlapViewT *b511OverlapViewTesting.T) {
	for b511OverlapViewValue := 0; b511OverlapViewValue < 256; b511OverlapViewValue++ {
		b511OverlapViewImage := []byte{17, 34, 51, 222, 223, 224, 225, 64, byte(b511OverlapViewValue), 226, 227, 228}
		for _, b511OverlapViewVariant := range []string{B511OverlapViewOP00Reserved, B511OverlapViewOP00Fan} {
			b511OverlapViewGot := B511OverlapViewInterpret(B511OverlapViewSelection{B511OverlapViewFamily, b511OverlapViewVariant, b511OverlapViewOpPtr(0), "positive"}, b511OverlapViewImage)
			if b511OverlapViewGot.B511OverlapViewRawFlag == nil || *b511OverlapViewGot.B511OverlapViewRawFlag != byte(b511OverlapViewValue) || !b511OverlapViewReflect.DeepEqual(b511OverlapViewGot.B511OverlapViewRaw, b511OverlapViewImage) {
				b511OverlapViewT.Fatal("ab11an11 lost flag bits")
			}
			if b511OverlapViewVariant == B511OverlapViewOP00Reserved {
				if b511OverlapViewGot.B511OverlapViewFanRunning != nil {
					b511OverlapViewT.Fatal("ab11an11 reserved bit named fan")
				}
			} else if b511OverlapViewGot.B511OverlapViewFanRunning == nil || *b511OverlapViewGot.B511OverlapViewFanRunning != (b511OverlapViewValue >= 128) {
				b511OverlapViewT.Fatal("ab11an11 fan bit mismatch")
			}
		}
		b511OverlapViewSelection := B511OverlapViewSelection{B511OverlapViewFamily, B511OverlapViewOP07ExtendedRoom, b511OverlapViewOpPtr(7), "positive"}
		b511OverlapViewGot := B511OverlapViewInterpret(b511OverlapViewSelection, b511OverlapViewImage)
		if b511OverlapViewGot.B511OverlapViewRawPressure == nil || *b511OverlapViewGot.B511OverlapViewRawPressure != byte(b511OverlapViewValue) || b511OverlapViewGot.B511OverlapViewRoomSetpointCelsius == nil || *b511OverlapViewGot.B511OverlapViewRoomSetpointCelsius != 32 || b511OverlapViewGot.B511OverlapViewRawOverlay == nil || *b511OverlapViewGot.B511OverlapViewRawOverlay != 64 {
			b511OverlapViewT.Fatal("ab11an11 pressure changed independent overlay/room")
		}
		if b511OverlapViewValue <= 100 {
			if b511OverlapViewGot.B511OverlapViewPressurePercent == nil || *b511OverlapViewGot.B511OverlapViewPressurePercent != byte(b511OverlapViewValue) || b511OverlapViewGot.B511OverlapViewPressureDomain != "in-domain" {
				b511OverlapViewT.Fatal("ab11an11 percent codec mismatch")
			}
		} else if b511OverlapViewGot.B511OverlapViewPressurePercent != nil || b511OverlapViewGot.B511OverlapViewPressureDomain != "out-of-domain" {
			b511OverlapViewT.Fatal("ab11an11 invalid pressure clipped")
		}
		b511OverlapViewImage[7], b511OverlapViewImage[8] = byte(b511OverlapViewValue), 100
		b511OverlapViewGot = B511OverlapViewInterpret(b511OverlapViewSelection, b511OverlapViewImage)
		if b511OverlapViewGot.B511OverlapViewRawOverlay == nil || *b511OverlapViewGot.B511OverlapViewRawOverlay != byte(b511OverlapViewValue) || b511OverlapViewGot.B511OverlapViewRoomSetpointCelsius == nil || *b511OverlapViewGot.B511OverlapViewRoomSetpointCelsius != float64(b511OverlapViewValue)*0.5 || b511OverlapViewGot.B511OverlapViewPressurePercent == nil || *b511OverlapViewGot.B511OverlapViewPressurePercent != 100 {
			b511OverlapViewT.Fatal("ab11an11 overlapping room/pressure offsets mismatch")
		}
	}
}

func TestAb11An11DetachedBuffers(b511OverlapViewT *b511OverlapViewTesting.T) {
	for _, b511OverlapViewVariant := range []string{B511OverlapViewOP00Fan, B511OverlapViewOP07ExtendedRoom, "ab11an11-unknown"} {
		b511OverlapViewImage := []byte{17, 34, 51, 0, 0, 0, 0, 64, 100, 0, 0, 0}
		b511OverlapViewOperation := byte(0)
		if b511OverlapViewVariant == B511OverlapViewOP07ExtendedRoom {
			b511OverlapViewOperation = 7
		}
		b511OverlapViewSelection := B511OverlapViewSelection{B511OverlapViewFamily, b511OverlapViewVariant, &b511OverlapViewOperation, "positive"}
		b511OverlapViewGot := B511OverlapViewInterpret(b511OverlapViewSelection, b511OverlapViewImage)
		b511OverlapViewBefore := B511OverlapViewInterpret(b511OverlapViewSelection, append([]byte(nil), b511OverlapViewImage...))
		for b511OverlapViewIndex := range b511OverlapViewImage {
			b511OverlapViewImage[b511OverlapViewIndex] = 255
		}
		if !b511OverlapViewReflect.DeepEqual(b511OverlapViewGot, b511OverlapViewBefore) {
			b511OverlapViewT.Fatal("ab11an11 input buffer aliases result")
		}
		b511OverlapViewGot.B511OverlapViewRaw[7], b511OverlapViewGot.B511OverlapViewRaw[8] = 0, 0
		if b511OverlapViewImage[7] != 255 || (b511OverlapViewGot.B511OverlapViewRawFlag != nil && *b511OverlapViewGot.B511OverlapViewRawFlag != 100) || (b511OverlapViewGot.B511OverlapViewRawOverlay != nil && *b511OverlapViewGot.B511OverlapViewRawOverlay != 64) || (b511OverlapViewGot.B511OverlapViewRawPressure != nil && *b511OverlapViewGot.B511OverlapViewRawPressure != 100) {
			b511OverlapViewT.Fatal("ab11an11 raw buffer aliases input/fields")
		}
		if b511OverlapViewGot.B511OverlapViewRawOverlay != nil {
			*b511OverlapViewGot.B511OverlapViewRawOverlay = 0
		}
		if b511OverlapViewGot.B511OverlapViewRawPressure != nil {
			*b511OverlapViewGot.B511OverlapViewRawPressure = 0
		}
		if (b511OverlapViewGot.B511OverlapViewRoomSetpointCelsius != nil && *b511OverlapViewGot.B511OverlapViewRoomSetpointCelsius != 32) || (b511OverlapViewGot.B511OverlapViewPressurePercent != nil && *b511OverlapViewGot.B511OverlapViewPressurePercent != 100) {
			b511OverlapViewT.Fatal("ab11an11 independent candidate fields alias")
		}
	}
}
