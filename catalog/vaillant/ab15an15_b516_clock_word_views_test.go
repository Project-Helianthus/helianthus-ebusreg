package vaillant

import (
	ab15an15JSON "encoding/json"
	ab15an15OS "os"
	ab15an15Reflect "reflect"
	ab15an15Testing "testing"
)

func ab15an15Selector(ab15an15Operation byte) Ab15An15Selector {
	return Ab15An15Selector{Ab15An15Family: Ab15An15Family, Ab15An15Variant: Ab15An15Variant, Ab15An15Operation: ab15an15Operation, Ab15An15Direction: "positive"}
}

func ab15an15AssertField(ab15an15T *ab15an15Testing.T, ab15an15Actual Ab15An15Field, ab15an15State string, ab15an15Value *uint16) {
	ab15an15T.Helper()
	if !ab15an15Reflect.DeepEqual(ab15an15Actual, Ab15An15Field{Ab15An15State: ab15an15State, Ab15An15Value: ab15an15Value}) {
		ab15an15T.Fatalf("ab15an15 field: %+v expected %s %v", ab15an15Actual, ab15an15State, ab15an15Value)
	}
}

func TestAb15An15ClockVectors(ab15an15T *ab15an15Testing.T) {
	ab15an15Data, ab15an15Error := ab15an15OS.ReadFile("testdata/ab15an15_b516_clock_word_views.json")
	if ab15an15Error != nil {
		ab15an15T.Fatal(ab15an15Error)
	}
	var ab15an15Vectors []struct {
		Ab15An15Label  string
		Ab15An15Clock  [5][]byte
		Ab15An15Failed [5]bool
		Ab15An15States [5]string
		Ab15An15Values [5]*uint16
	}
	if ab15an15Error = ab15an15JSON.Unmarshal(ab15an15Data, &ab15an15Vectors); ab15an15Error != nil {
		ab15an15T.Fatal(ab15an15Error)
	}
	if len(ab15an15Vectors) != 50 {
		ab15an15T.Fatalf("ab15an15 count %d", len(ab15an15Vectors))
	}
	for _, ab15an15Vector := range ab15an15Vectors {
		ab15an15T.Run(ab15an15Vector.Ab15An15Label, func(ab15an15T *ab15an15Testing.T) {
			ab15an15Input := Ab15An15ClockInput{}
			for ab15an15Index := range ab15an15Input.Ab15An15Clock {
				ab15an15Input.Ab15An15Clock[ab15an15Index] = Ab15An15Cell{Ab15An15Raw: ab15an15Vector.Ab15An15Clock[ab15an15Index], Ab15An15Failed: ab15an15Vector.Ab15An15Failed[ab15an15Index]}
			}
			ab15an15Result := Ab15An15ClockView(ab15an15Selector(8), ab15an15Input)
			if !ab15an15Result.Ab15An15Candidate.Ab15An15Admitted || ab15an15Result.Ab15An15Clock == nil || ab15an15Result.Ab15An15SetpointCelsius != nil || !ab15an15Reflect.DeepEqual(ab15an15Result.Ab15An15Raw, ab15an15Input) {
				ab15an15T.Fatal("ab15an15 clock admission/raw/missing setpoint")
			}
			for ab15an15Index := range ab15an15Input.Ab15An15Clock {
				ab15an15AssertField(ab15an15T, ab15an15Result.Ab15An15Clock[ab15an15Index], ab15an15Vector.Ab15An15States[ab15an15Index], ab15an15Vector.Ab15An15Values[ab15an15Index])
			}
		})
	}
}

func TestAb15An15Words(ab15an15T *ab15an15Testing.T) {
	for ab15an15Slot := 0; ab15an15Slot < 4; ab15an15Slot++ {
		for _, ab15an15Case := range []struct {
			Ab15An15Raw    []byte
			Ab15An15Failed bool
			Ab15An15State  string
			Ab15An15Value  uint16
		}{
			{[]byte{0, 0}, false, "numeric", 0}, {[]byte{1, 0}, false, "numeric", 1}, {[]byte{0, 1}, false, "numeric", 256}, {[]byte{255, 255}, false, "numeric", 65535}, {[]byte{52, 18}, false, "numeric", 4660}, {nil, false, "missing", 0}, {[]byte{1}, false, "failed", 0}, {[]byte{1, 2, 3}, false, "failed", 0}, {[]byte{1, 0}, true, "failed", 0},
		} {
			ab15an15Words := [4]Ab15An15Cell{{Ab15An15Raw: []byte{0, 0}}, {Ab15An15Raw: []byte{1, 0}}, {Ab15An15Raw: []byte{0, 1}}, {Ab15An15Raw: []byte{255, 255}}}
			ab15an15Expected := [4]uint16{0, 1, 256, 65535}
			ab15an15Words[ab15an15Slot] = Ab15An15Cell{Ab15An15Raw: ab15an15Case.Ab15An15Raw, Ab15An15Failed: ab15an15Case.Ab15An15Failed}
			ab15an15Result := Ab15An15WordView(ab15an15Selector(13), ab15an15Words)
			if !ab15an15Result.Ab15An15Candidate.Ab15An15Admitted || ab15an15Result.Ab15An15Words == nil || !ab15an15Reflect.DeepEqual(ab15an15Result.Ab15An15Raw, ab15an15Words) || ab15an15Result.Ab15An15PositionsGrade != "Unknown" {
				ab15an15T.Fatal("ab15an15 word admission/raw/positions")
			}
			for ab15an15Index := range ab15an15Words {
				ab15an15State := "numeric"
				ab15an15Value := &ab15an15Expected[ab15an15Index]
				if ab15an15Index == ab15an15Slot {
					ab15an15State = ab15an15Case.Ab15An15State
					ab15an15Value = &ab15an15Case.Ab15An15Value
					if ab15an15State != "numeric" {
						ab15an15Value = nil
					}
				}
				ab15an15AssertField(ab15an15T, ab15an15Result.Ab15An15Words[ab15an15Index], ab15an15State, ab15an15Value)
			}
		}
	}
}

func TestAb15An15AdmissionAndDetachedRaw(ab15an15T *ab15an15Testing.T) {
	for ab15an15Axis := 0; ab15an15Axis < 6; ab15an15Axis++ {
		ab15an15ClockSelector := ab15an15Selector(8)
		ab15an15WordSelector := ab15an15Selector(13)
		for _, ab15an15Selection := range []*Ab15An15Selector{&ab15an15ClockSelector, &ab15an15WordSelector} {
			switch ab15an15Axis {
			case 1:
				ab15an15Selection.Ab15An15Family = "ab15an15_other_family"
			case 2:
				ab15an15Selection.Ab15An15Variant = "ab15an15_other_variant"
			case 3:
				ab15an15Selection.Ab15An15Operation = 0
			case 4:
				ab15an15Selection.Ab15An15Direction = "request"
			case 5:
				if ab15an15Selection.Ab15An15Operation == 8 {
					ab15an15Selection.Ab15An15Operation = 13
				} else {
					ab15an15Selection.Ab15An15Operation = 8
				}
			}
		}
		ab15an15Input := Ab15An15ClockInput{}
		ab15an15Words := [4]Ab15An15Cell{}
		for ab15an15Index := range ab15an15Input.Ab15An15Clock {
			ab15an15Input.Ab15An15Clock[ab15an15Index].Ab15An15Raw = []byte{18}
		}
		for ab15an15Index := range ab15an15Words {
			ab15an15Words[ab15an15Index].Ab15An15Raw = []byte{1, 2}
		}
		ab15an15Bytes := [8]byte{1, 255, 254, 0, 1, 2, 3, 255}
		ab15an15Input.Ab15An15Setpoint = &ab15an15Bytes[0]
		ab15an15Input.Ab15An15Mode = &ab15an15Bytes[1]
		ab15an15Input.Ab15An15SpecialFunction = &ab15an15Bytes[2]
		for ab15an15Index := range ab15an15Input.Ab15An15Messages {
			ab15an15Input.Ab15An15Messages[ab15an15Index] = &ab15an15Bytes[ab15an15Index+3]
		}
		ab15an15Clock := Ab15An15ClockView(ab15an15ClockSelector, ab15an15Input)
		ab15an15Word := Ab15An15WordView(ab15an15WordSelector, ab15an15Words)
		if !ab15an15Reflect.DeepEqual(ab15an15Clock.Ab15An15Raw, ab15an15Input) || !ab15an15Reflect.DeepEqual(ab15an15Word.Ab15An15Raw, ab15an15Words) {
			ab15an15T.Fatal("ab15an15 raw retention")
		}
		if ab15an15Axis == 0 {
			if ab15an15Clock.Ab15An15SetpointCelsius == nil || *ab15an15Clock.Ab15An15SetpointCelsius != 0.5 {
				ab15an15T.Fatal("ab15an15 half Celsius")
			}
		} else if ab15an15Clock.Ab15An15Candidate.Ab15An15Admitted || ab15an15Word.Ab15An15Candidate.Ab15An15Admitted || ab15an15Clock.Ab15An15Clock != nil || ab15an15Word.Ab15An15Words != nil || ab15an15Clock.Ab15An15SetpointCelsius != nil {
			ab15an15T.Fatal("ab15an15 raw-only admission")
		}
		ab15an15ExpectedCandidate := Ab15An15Candidate{Ab15An15Admitted: ab15an15Axis == 0, Ab15An15ContractGrade: "Proven", Ab15An15FieldGrade: "Hypothesis", Ab15An15NativeSupport: "Unknown", Ab15An15PhysicalAdaptation: "Unknown", Ab15An15Presence: "Unknown", Ab15An15Freshness: "Unknown", Ab15An15Execution: "Unknown"}
		if ab15an15Clock.Ab15An15Candidate != ab15an15ExpectedCandidate || ab15an15Word.Ab15An15Candidate != ab15an15ExpectedCandidate || ab15an15Clock.Ab15An15TablesGrade != "Unknown" {
			ab15an15T.Fatal("ab15an15 evidence grades")
		}
		for ab15an15Index := range ab15an15Bytes {
			ab15an15Bytes[ab15an15Index] ^= 128
		}
		if *ab15an15Clock.Ab15An15Raw.Ab15An15Setpoint != 1 || *ab15an15Clock.Ab15An15Raw.Ab15An15Mode != 255 || *ab15an15Clock.Ab15An15Raw.Ab15An15SpecialFunction != 254 {
			ab15an15T.Fatal("ab15an15 optional aliases")
		}
		for ab15an15Index, ab15an15Expected := range [5]byte{0, 1, 2, 3, 255} {
			if *ab15an15Clock.Ab15An15Raw.Ab15An15Messages[ab15an15Index] != ab15an15Expected {
				ab15an15T.Fatal("ab15an15 message aliases/index")
			}
		}
		for ab15an15Index := range ab15an15Input.Ab15An15Clock {
			ab15an15Input.Ab15An15Clock[ab15an15Index].Ab15An15Raw[0] = 99
			if ab15an15Clock.Ab15An15Raw.Ab15An15Clock[ab15an15Index].Ab15An15Raw[0] != 18 {
				ab15an15T.Fatal("ab15an15 clock aliases")
			}
		}
		for ab15an15Index := range ab15an15Words {
			ab15an15Words[ab15an15Index].Ab15An15Raw[0] = 99
			ab15an15Words[ab15an15Index].Ab15An15Raw[1] = 99
			if !ab15an15Reflect.DeepEqual(ab15an15Word.Ab15An15Raw[ab15an15Index].Ab15An15Raw, []byte{1, 2}) {
				ab15an15T.Fatal("ab15an15 word aliases")
			}
		}
	}
}

func TestAb15An15OptionalFields(ab15an15T *ab15an15Testing.T) {
	for ab15an15Index := 0; ab15an15Index < 8; ab15an15Index++ {
		for _, ab15an15Case := range []struct {
			Ab15An15Raw     byte
			Ab15An15Celsius float64
		}{{0, 0}, {1, 0.5}, {127, 63.5}, {255, 127.5}} {
			ab15an15Value := ab15an15Case.Ab15An15Raw
			ab15an15Input := Ab15An15ClockInput{}
			switch ab15an15Index {
			case 0:
				ab15an15Input.Ab15An15Setpoint = &ab15an15Value
			case 1:
				ab15an15Input.Ab15An15Mode = &ab15an15Value
			case 2:
				ab15an15Input.Ab15An15SpecialFunction = &ab15an15Value
			default:
				ab15an15Input.Ab15An15Messages[ab15an15Index-3] = &ab15an15Value
			}
			ab15an15Result := Ab15An15ClockView(ab15an15Selector(8), ab15an15Input)
			if !ab15an15Reflect.DeepEqual(ab15an15Result.Ab15An15Raw, ab15an15Input) {
				ab15an15T.Fatal("ab15an15 optional cross-contamination")
			}
			if ab15an15Index == 0 {
				if ab15an15Result.Ab15An15SetpointCelsius == nil || *ab15an15Result.Ab15An15SetpointCelsius != ab15an15Case.Ab15An15Celsius {
					ab15an15T.Fatal("ab15an15 optional setpoint")
				}
			} else if ab15an15Result.Ab15An15SetpointCelsius != nil {
				ab15an15T.Fatal("ab15an15 missing setpoint")
			}
		}
	}
}
