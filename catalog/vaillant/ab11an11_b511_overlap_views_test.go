package vaillant

import (
	ab11an11json "encoding/json"
	ab11an11os "os"
	ab11an11reflect "reflect"
	ab11an11testing "testing"
)

func TestAb11An11IndependentVectors(ab11an11T *ab11an11testing.T) {
	ab11an11Data, ab11an11Err := ab11an11os.ReadFile("testdata/ab11an11_b511_overlap_views.json")
	if ab11an11Err != nil {
		ab11an11T.Fatal(ab11an11Err)
	}
	var ab11an11Cases []struct {
		Ab11An11Label     string
		Ab11An11Selection Ab11An11Selection
		Ab11An11Image     []byte
		Ab11An11Expected  Ab11An11View
	}
	if ab11an11Err := ab11an11json.Unmarshal(ab11an11Data, &ab11an11Cases); ab11an11Err != nil {
		ab11an11T.Fatal(ab11an11Err)
	}
	if len(ab11an11Cases) != 11 {
		ab11an11T.Fatalf("ab11an11 expected 11 independent vectors, got %d", len(ab11an11Cases))
	}
	for _, ab11an11Case := range ab11an11Cases {
		ab11an11T.Run(ab11an11Case.Ab11An11Label, func(ab11an11T *ab11an11testing.T) {
			ab11an11Got := Ab11An11Interpret(ab11an11Case.Ab11An11Selection, ab11an11Case.Ab11An11Image)
			if !ab11an11reflect.DeepEqual(ab11an11Got, ab11an11Case.Ab11An11Expected) {
				ab11an11T.Fatalf("ab11an11 got %+v; want %+v", ab11an11Got, ab11an11Case.Ab11An11Expected)
			}
		})
	}
}

func TestAb11An11SelectionAndExtent(ab11an11T *ab11an11testing.T) {
	for _, ab11an11Variant := range []string{Ab11An11OP00Reserved, Ab11An11OP00Fan, Ab11An11OP07Base, Ab11An11OP07ExtendedRoom} {
		ab11an11Operation, ab11an11Required := byte(0), 12
		if ab11an11Variant == Ab11An11OP07Base {
			ab11an11Operation, ab11an11Required = 7, 8
		}
		if ab11an11Variant == Ab11An11OP07ExtendedRoom {
			ab11an11Operation, ab11an11Required = 7, 9
		}
		for ab11an11Length := 0; ab11an11Length <= 15; ab11an11Length++ {
			ab11an11Image := make([]byte, ab11an11Length)
			for ab11an11Index := range ab11an11Image {
				ab11an11Image[ab11an11Index] = byte(ab11an11Index + 121)
			}
			ab11an11Selection := Ab11An11Selection{Ab11An11Family, ab11an11Variant, ab11an11Operation, "positive"}
			ab11an11Got := Ab11An11Interpret(ab11an11Selection, ab11an11Image)
			if !ab11an11Got.Ab11An11Selected || ab11an11Got.Ab11An11RequiredExtent != ab11an11Required || ab11an11Got.Ab11An11CompleteLogicalExtent != (ab11an11Length >= ab11an11Required) || ab11an11Got.Ab11An11Selection != ab11an11Selection || ab11an11Got.Ab11An11AvailableExtent != ab11an11Length {
				ab11an11T.Fatalf("ab11an11 extent/selection: %+v", ab11an11Got)
			}
			if (ab11an11Got.Ab11An11RawFlag != nil) != (ab11an11Operation == 0 && ab11an11Length > 8) || (ab11an11Got.Ab11An11RawOverlay != nil) != (ab11an11Operation == 7 && ab11an11Length > 7) || (ab11an11Got.Ab11An11RawPressure != nil) != (ab11an11Variant == Ab11An11OP07ExtendedRoom && ab11an11Length > 8) {
				ab11an11T.Fatalf("ab11an11 availability: %+v", ab11an11Got)
			}
			for ab11an11OtherOp := 0; ab11an11OtherOp < 256; ab11an11OtherOp++ {
				if byte(ab11an11OtherOp) == ab11an11Operation {
					continue
				}
				ab11an11Wrong := ab11an11Selection
				ab11an11Wrong.Ab11An11Operation = byte(ab11an11OtherOp)
				ab11an11AssertUnknown(ab11an11T, ab11an11Wrong, ab11an11Image)
			}
			ab11an11Wrong := ab11an11Selection
			ab11an11Wrong.Ab11An11Family = "ab11an11-other-family"
			ab11an11AssertUnknown(ab11an11T, ab11an11Wrong, ab11an11Image)
			ab11an11Wrong = ab11an11Selection
			ab11an11Wrong.Ab11An11Variant = "ab11an11-unknown-variant"
			ab11an11AssertUnknown(ab11an11T, ab11an11Wrong, ab11an11Image)
			for _, ab11an11Direction := range []string{"", "request", "negative", "ab11an11-unknown-direction"} {
				ab11an11Wrong = ab11an11Selection
				ab11an11Wrong.Ab11An11Direction = ab11an11Direction
				ab11an11AssertUnknown(ab11an11T, ab11an11Wrong, ab11an11Image)
			}
		}
	}
}

func ab11an11AssertUnknown(ab11an11T *ab11an11testing.T, ab11an11Selection Ab11An11Selection, ab11an11Image []byte) {
	ab11an11T.Helper()
	ab11an11Expected := Ab11An11View{Ab11An11Selection: ab11an11Selection, Ab11An11Raw: append([]byte(nil), ab11an11Image...), Ab11An11AvailableExtent: len(ab11an11Image), Ab11An11PressureDomain: "Unknown", Ab11An11CandidateGrade: "Unknown", Ab11An11PhysicalQualification: "Unknown"}
	if ab11an11Got := Ab11An11Interpret(ab11an11Selection, ab11an11Image); !ab11an11reflect.DeepEqual(ab11an11Got, ab11an11Expected) {
		ab11an11T.Fatalf("ab11an11 unknown selection invented fields: %+v", ab11an11Got)
	}
}

func TestAb11An11EveryBitAndOverlap(ab11an11T *ab11an11testing.T) {
	for ab11an11Value := 0; ab11an11Value < 256; ab11an11Value++ {
		ab11an11Image := []byte{17, 34, 51, 222, 223, 224, 225, 64, byte(ab11an11Value), 226, 227, 228}
		for _, ab11an11Variant := range []string{Ab11An11OP00Reserved, Ab11An11OP00Fan} {
			ab11an11Got := Ab11An11Interpret(Ab11An11Selection{Ab11An11Family, ab11an11Variant, 0, "positive"}, ab11an11Image)
			if ab11an11Got.Ab11An11RawFlag == nil || *ab11an11Got.Ab11An11RawFlag != byte(ab11an11Value) || !ab11an11reflect.DeepEqual(ab11an11Got.Ab11An11Raw, ab11an11Image) {
				ab11an11T.Fatal("ab11an11 lost flag bits")
			}
			if ab11an11Variant == Ab11An11OP00Reserved {
				if ab11an11Got.Ab11An11FanRunning != nil {
					ab11an11T.Fatal("ab11an11 reserved bit named fan")
				}
			} else if ab11an11Got.Ab11An11FanRunning == nil || *ab11an11Got.Ab11An11FanRunning != (ab11an11Value >= 128) {
				ab11an11T.Fatal("ab11an11 fan bit mismatch")
			}
		}
		ab11an11Selection := Ab11An11Selection{Ab11An11Family, Ab11An11OP07ExtendedRoom, 7, "positive"}
		ab11an11Got := Ab11An11Interpret(ab11an11Selection, ab11an11Image)
		if ab11an11Got.Ab11An11RawPressure == nil || *ab11an11Got.Ab11An11RawPressure != byte(ab11an11Value) || ab11an11Got.Ab11An11RoomSetpointCelsius == nil || *ab11an11Got.Ab11An11RoomSetpointCelsius != 32 || ab11an11Got.Ab11An11RawOverlay == nil || *ab11an11Got.Ab11An11RawOverlay != 64 {
			ab11an11T.Fatal("ab11an11 pressure changed independent overlay/room")
		}
		if ab11an11Value <= 100 {
			if ab11an11Got.Ab11An11PressurePercent == nil || *ab11an11Got.Ab11An11PressurePercent != byte(ab11an11Value) || ab11an11Got.Ab11An11PressureDomain != "in-domain" {
				ab11an11T.Fatal("ab11an11 percent codec mismatch")
			}
		} else if ab11an11Got.Ab11An11PressurePercent != nil || ab11an11Got.Ab11An11PressureDomain != "out-of-domain" {
			ab11an11T.Fatal("ab11an11 invalid pressure clipped")
		}
		ab11an11Image[7], ab11an11Image[8] = byte(ab11an11Value), 100
		ab11an11Got = Ab11An11Interpret(ab11an11Selection, ab11an11Image)
		if ab11an11Got.Ab11An11RawOverlay == nil || *ab11an11Got.Ab11An11RawOverlay != byte(ab11an11Value) || ab11an11Got.Ab11An11RoomSetpointCelsius == nil || *ab11an11Got.Ab11An11RoomSetpointCelsius != float64(ab11an11Value)*0.5 || ab11an11Got.Ab11An11PressurePercent == nil || *ab11an11Got.Ab11An11PressurePercent != 100 {
			ab11an11T.Fatal("ab11an11 overlapping room/pressure offsets mismatch")
		}
	}
}

func TestAb11An11DetachedBuffers(ab11an11T *ab11an11testing.T) {
	for _, ab11an11Variant := range []string{Ab11An11OP00Fan, Ab11An11OP07ExtendedRoom, "ab11an11-unknown"} {
		ab11an11Image := []byte{17, 34, 51, 0, 0, 0, 0, 64, 100, 0, 0, 0}
		ab11an11Operation := byte(0)
		if ab11an11Variant == Ab11An11OP07ExtendedRoom {
			ab11an11Operation = 7
		}
		ab11an11Selection := Ab11An11Selection{Ab11An11Family, ab11an11Variant, ab11an11Operation, "positive"}
		ab11an11Got := Ab11An11Interpret(ab11an11Selection, ab11an11Image)
		ab11an11Before := Ab11An11Interpret(ab11an11Selection, append([]byte(nil), ab11an11Image...))
		for ab11an11Index := range ab11an11Image {
			ab11an11Image[ab11an11Index] = 255
		}
		if !ab11an11reflect.DeepEqual(ab11an11Got, ab11an11Before) {
			ab11an11T.Fatal("ab11an11 input buffer aliases result")
		}
		ab11an11Got.Ab11An11Raw[7], ab11an11Got.Ab11An11Raw[8] = 0, 0
		if ab11an11Image[7] != 255 || (ab11an11Got.Ab11An11RawFlag != nil && *ab11an11Got.Ab11An11RawFlag != 100) || (ab11an11Got.Ab11An11RawOverlay != nil && *ab11an11Got.Ab11An11RawOverlay != 64) || (ab11an11Got.Ab11An11RawPressure != nil && *ab11an11Got.Ab11An11RawPressure != 100) {
			ab11an11T.Fatal("ab11an11 raw buffer aliases input/fields")
		}
		if ab11an11Got.Ab11An11RawOverlay != nil {
			*ab11an11Got.Ab11An11RawOverlay = 0
		}
		if ab11an11Got.Ab11An11RawPressure != nil {
			*ab11an11Got.Ab11An11RawPressure = 0
		}
		if (ab11an11Got.Ab11An11RoomSetpointCelsius != nil && *ab11an11Got.Ab11An11RoomSetpointCelsius != 32) || (ab11an11Got.Ab11An11PressurePercent != nil && *ab11an11Got.Ab11An11PressurePercent != 100) {
			ab11an11T.Fatal("ab11an11 independent candidate fields alias")
		}
	}
}
