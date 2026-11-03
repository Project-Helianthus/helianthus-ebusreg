package system

import (
	ab01an01JSON "encoding/json"
	ab01an01OS "os"
	ab01an01Reflect "reflect"
	ab01an01Testing "testing"
)

type ab01an01Vector struct {
	Ab01An01Label       string
	Ab01An01Family      string
	Ab01An01GG          byte
	Ab01An01II          byte
	Ab01An01RR          uint16
	Ab01An01Alternative string
	Ab01An01Bytes       [3]byte
	Ab01An01Expected    [3]byte
}

func ab01an01Base() Ab01An01Input {
	ab01an01II := byte(0)
	return Ab01An01Input{Ab01An01Family: "controller-700", Ab01An01Software: "synthetic-software", Ab01An01Hardware: "synthetic-hardware", Ab01An01APIVersion: "synthetic-api", Ab01An01APIRevision: "synthetic-revision", Ab01An01OP: 2, Ab01An01II: &ab01an01II, Ab01An01RequestedRR: 0x34, Ab01An01ReturnedRR: 0x34, Ab01An01Alternative: "date", Ab01An01ResponseState: "echo-payload", Ab01An01Raw: []byte{31, 12, 99}, Ab01An01Cells: [3]Ab01An01Cell{{31, "present"}, {12, "present"}, {99, "present"}}}
}

func ab01an01Assert(ab01an01T *ab01an01Testing.T, ab01an01Input Ab01An01Input, ab01an01Expected [3]*byte, ab01an01States [3]string, ab01an01Complete bool, ab01an01Grade string) {
	ab01an01T.Helper()
	ab01an01Result := Ab01An01Interpret(ab01an01Input)
	if !ab01an01Reflect.DeepEqual(ab01an01Result.Ab01An01Context, ab01an01Input) {
		ab01an01T.Fatal("ab01an01 context lost")
	}
	if ab01an01Result.Ab01An01CompleteCandidate != ab01an01Complete || ab01an01Result.Ab01An01Grade != ab01an01Grade || ab01an01Result.Ab01An01ContractGrade != "Proven" || ab01an01Result.Ab01An01NativeQualification != "Unknown" || ab01an01Result.Ab01An01CombinedTimestamp != "Unknown" || ab01an01Result.Ab01An01EditValidation != "Unknown" {
		ab01an01T.Fatalf("ab01an01 result: %+v", ab01an01Result)
	}
	for ab01an01Index, ab01an01Component := range ab01an01Result.Ab01An01Components {
		if ab01an01Component.Ab01An01Cell != ab01an01Input.Ab01An01Cells[ab01an01Index] || ab01an01Component.Ab01An01State != ab01an01States[ab01an01Index] || !ab01an01Reflect.DeepEqual(ab01an01Component.Ab01An01Candidate, ab01an01Expected[ab01an01Index]) {
			ab01an01T.Fatalf("ab01an01 component %d: %+v", ab01an01Index, ab01an01Component)
		}
	}
}

func TestAb01An01Vectors(ab01an01T *ab01an01Testing.T) {
	ab01an01Data, ab01an01Err := ab01an01OS.ReadFile("testdata/ab01an01_b524_calendar_views.json")
	if ab01an01Err != nil {
		ab01an01T.Fatal(ab01an01Err)
	}
	var ab01an01Vectors []ab01an01Vector
	if ab01an01Err = ab01an01JSON.Unmarshal(ab01an01Data, &ab01an01Vectors); ab01an01Err != nil {
		ab01an01T.Fatal(ab01an01Err)
	}
	if len(ab01an01Vectors) != 26 {
		ab01an01T.Fatal("ab01an01 expected 26 rows")
	}
	ab01an01Seen := map[string]map[[3]uint16]bool{}
	for _, ab01an01Vector := range ab01an01Vectors {
		ab01an01T.Run(ab01an01Vector.Ab01An01Label, func(ab01an01T *ab01an01Testing.T) {
			ab01an01Input := ab01an01Base()
			ab01an01Input.Ab01An01Family = ab01an01Vector.Ab01An01Family
			ab01an01Input.Ab01An01RequestedGG = ab01an01Vector.Ab01An01GG
			ab01an01Input.Ab01An01ReturnedGG = ab01an01Vector.Ab01An01GG
			*ab01an01Input.Ab01An01II = ab01an01Vector.Ab01An01II
			ab01an01Input.Ab01An01RequestedRR = ab01an01Vector.Ab01An01RR
			ab01an01Input.Ab01An01ReturnedRR = ab01an01Vector.Ab01An01RR
			ab01an01Input.Ab01An01Alternative = ab01an01Vector.Ab01An01Alternative
			for ab01an01Index, ab01an01Byte := range ab01an01Vector.Ab01An01Bytes {
				ab01an01Input.Ab01An01Cells[ab01an01Index].Ab01An01Raw = ab01an01Byte
			}
			ab01an01Assert(ab01an01T, ab01an01Input, [3]*byte{&ab01an01Vector.Ab01An01Expected[0], &ab01an01Vector.Ab01An01Expected[1], &ab01an01Vector.Ab01An01Expected[2]}, [3]string{"candidate", "candidate", "candidate"}, true, "Hypothesis")
			if ab01an01Seen[ab01an01Vector.Ab01An01Family] == nil {
				ab01an01Seen[ab01an01Vector.Ab01An01Family] = map[[3]uint16]bool{}
			}
			ab01an01Key := [3]uint16{uint16(ab01an01Vector.Ab01An01GG), uint16(ab01an01Vector.Ab01An01II), ab01an01Vector.Ab01An01RR}
			if ab01an01Seen[ab01an01Vector.Ab01An01Family][ab01an01Key] {
				ab01an01T.Fatal("ab01an01 duplicate row")
			}
			ab01an01Seen[ab01an01Vector.Ab01An01Family][ab01an01Key] = true
		})
	}
}

func TestAb01An01Context(ab01an01T *ab01an01Testing.T) {
	for ab01an01Axis := 0; ab01an01Axis < 20; ab01an01Axis++ {
		ab01an01Input := ab01an01Base()
		switch ab01an01Axis {
		case 0:
			ab01an01Input.Ab01An01Family = ""
		case 1:
			ab01an01Input.Ab01An01Family = "unknown-family"
		case 2:
			ab01an01Input.Ab01An01Software = ""
		case 3:
			ab01an01Input.Ab01An01Hardware = ""
		case 4:
			ab01an01Input.Ab01An01APIVersion = ""
		case 5:
			ab01an01Input.Ab01An01APIRevision = ""
		case 6:
			ab01an01Input.Ab01An01OP = 6
		case 7:
			ab01an01Input.Ab01An01RequestedGG = 2
			ab01an01Input.Ab01An01ReturnedGG = 2
		case 8:
			*ab01an01Input.Ab01An01II = 1
		case 9:
			ab01an01Input.Ab01An01II = nil
		case 10:
			ab01an01Input.Ab01An01RequestedRR = 0xffff
			ab01an01Input.Ab01An01ReturnedRR = 0xffff
		case 11:
			ab01an01Input.Ab01An01Alternative = "time"
		case 12:
			ab01an01Input.Ab01An01ReturnedGG = 1
		case 13:
			ab01an01Input.Ab01An01ReturnedRR = 0x35
		case 14:
			ab01an01Input.Ab01An01ResponseState = "one-byte-status"
		case 15:
			ab01an01Input.Ab01An01ResponseState = "transport-error"
		case 16:
			ab01an01Input.Ab01An01ResponseState = ""
		case 17:
			ab01an01Input.Ab01An01RequestedRR = 0xda
			ab01an01Input.Ab01An01ReturnedRR = 0xda
		case 18:
			ab01an01Input.Ab01An01Family = "controller-720"
			ab01an01Input.Ab01An01RequestedRR = 0x71
			ab01an01Input.Ab01An01ReturnedRR = 0x71
		case 19:
			ab01an01Input.Ab01An01Alternative = "timestamp"
		}
		ab01an01Assert(ab01an01T, ab01an01Input, [3]*byte{}, [3]string{"Unknown", "Unknown", "Unknown"}, false, "Unknown")
	}
	for ab01an01Axis := 0; ab01an01Axis < 5; ab01an01Axis++ {
		ab01an01Input := ab01an01Base()
		switch ab01an01Axis {
		case 0:
			ab01an01Input.Ab01An01Software = "synthetic-other"
		case 1:
			ab01an01Input.Ab01An01Hardware = "synthetic-other"
		case 2:
			ab01an01Input.Ab01An01APIVersion = "synthetic-other"
		case 3:
			ab01an01Input.Ab01An01APIRevision = "synthetic-other"
		case 4:
			ab01an01Input.Ab01An01RequestedGG = 3
			ab01an01Input.Ab01An01ReturnedGG = 3
			ab01an01Input.Ab01An01RequestedRR = 3
			ab01an01Input.Ab01An01ReturnedRR = 3
			*ab01an01Input.Ab01An01II = 255
		}
		ab01an01Expected := [3]byte{31, 12, 99}
		ab01an01Assert(ab01an01T, ab01an01Input, [3]*byte{&ab01an01Expected[0], &ab01an01Expected[1], &ab01an01Expected[2]}, [3]string{"candidate", "candidate", "candidate"}, true, "Hypothesis")
	}
}

func TestAb01An01Components(ab01an01T *ab01an01Testing.T) {
	for ab01an01Index := 0; ab01an01Index < 3; ab01an01Index++ {
		for _, ab01an01State := range []string{"missing", "failed", "invalid"} {
			ab01an01Input := ab01an01Base()
			ab01an01Input.Ab01An01Cells[ab01an01Index].Ab01An01State = ab01an01State
			ab01an01Expected := [3]*byte{&ab01an01Input.Ab01An01Cells[0].Ab01An01Raw, &ab01an01Input.Ab01An01Cells[1].Ab01An01Raw, &ab01an01Input.Ab01An01Cells[2].Ab01An01Raw}
			ab01an01Expected[ab01an01Index] = nil
			ab01an01States := [3]string{"candidate", "candidate", "candidate"}
			if ab01an01State == "invalid" {
				ab01an01State = "Unknown"
			}
			ab01an01States[ab01an01Index] = ab01an01State
			ab01an01Assert(ab01an01T, ab01an01Input, ab01an01Expected, ab01an01States, false, "Hypothesis")
		}
		for ab01an01Value := 0; ab01an01Value < 256; ab01an01Value++ {
			ab01an01Input := ab01an01Base()
			ab01an01Input.Ab01An01Cells[ab01an01Index].Ab01An01Raw = byte(ab01an01Value)
			ab01an01Expected := [3]*byte{&ab01an01Input.Ab01An01Cells[0].Ab01An01Raw, &ab01an01Input.Ab01An01Cells[1].Ab01An01Raw, &ab01an01Input.Ab01An01Cells[2].Ab01An01Raw}
			ab01an01States := [3]string{"candidate", "candidate", "candidate"}
			ab01an01Valid := (ab01an01Index == 0 && ab01an01Value >= 1 && ab01an01Value <= 31) || (ab01an01Index == 1 && ab01an01Value >= 1 && ab01an01Value <= 12) || (ab01an01Index == 2 && ab01an01Value <= 99)
			if !ab01an01Valid {
				ab01an01Expected[ab01an01Index] = nil
				ab01an01States[ab01an01Index] = "out-of-domain"
			}
			ab01an01Assert(ab01an01T, ab01an01Input, ab01an01Expected, ab01an01States, ab01an01Valid, "Hypothesis")
			ab01an01Input.Ab01An01RequestedRR = 0x35
			ab01an01Input.Ab01An01ReturnedRR = 0x35
			ab01an01Input.Ab01An01Alternative = "time"
			ab01an01Expected[ab01an01Index] = &ab01an01Input.Ab01An01Cells[ab01an01Index].Ab01An01Raw
			ab01an01Assert(ab01an01T, ab01an01Input, ab01an01Expected, [3]string{"candidate", "candidate", "candidate"}, true, "Hypothesis")
		}
	}
}

func TestAb01An01Detached(ab01an01T *ab01an01Testing.T) {
	ab01an01Input := ab01an01Base()
	ab01an01Result := Ab01An01Interpret(ab01an01Input)
	ab01an01Input.Ab01An01Raw[0] = 0
	*ab01an01Input.Ab01An01II = 9
	ab01an01Input.Ab01An01Cells[0].Ab01An01Raw = 0
	if ab01an01Result.Ab01An01Context.Ab01An01Raw[0] != 31 || *ab01an01Result.Ab01An01Context.Ab01An01II != 0 || *ab01an01Result.Ab01An01Components[0].Ab01An01Candidate != 31 {
		ab01an01T.Fatal("ab01an01 input aliases result")
	}
	ab01an01Other := Ab01An01Interpret(ab01an01Result.Ab01An01Context)
	ab01an01Result.Ab01An01Context.Ab01An01Raw[1] = 0
	*ab01an01Result.Ab01An01Context.Ab01An01II = 8
	*ab01an01Result.Ab01An01Components[1].Ab01An01Candidate = 0
	if ab01an01Input.Ab01An01Raw[1] != 12 || *ab01an01Input.Ab01An01II != 9 || ab01an01Other.Ab01An01Context.Ab01An01Raw[1] != 12 || *ab01an01Other.Ab01An01Context.Ab01An01II != 0 || *ab01an01Other.Ab01An01Components[1].Ab01An01Candidate != 12 || ab01an01Result.Ab01An01Context.Ab01An01Cells[1].Ab01An01Raw != 12 {
		ab01an01T.Fatal("ab01an01 returned storage aliases")
	}
}
