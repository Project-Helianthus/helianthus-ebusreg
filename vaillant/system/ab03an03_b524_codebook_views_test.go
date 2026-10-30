package system

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

type ab03an03Vector struct {
	Ab03An03Name    string
	Ab03An03Family  string
	Ab03An03GG      byte
	Ab03An03RR      uint16
	Ab03An03Width   string
	Ab03An03Raw     []byte
	Ab03An03Integer int
	Ab03An03Label   string
}

func ab03an03Input() Ab03An03Input {
	ab03an03II := byte(0)
	return Ab03An03Input{Ab03An03Family: "controller-700", Ab03An03Software: "synthetic-software", Ab03An03Hardware: "synthetic-hardware", Ab03An03APIVersion: "synthetic-api", Ab03An03APIRevision: "synthetic-revision", Ab03An03OP: 2, Ab03An03GG: 2, Ab03An03II: &ab03an03II, Ab03An03RR: 2, Ab03An03Width: "u16-le", Ab03An03RequestedGG: 2, Ab03An03RequestedRR: 2, Ab03An03ReturnedGG: 2, Ab03An03ReturnedRR: 2, Ab03An03ResponseState: "echo-payload", Ab03An03Raw: []byte{5, 0}}
}

func TestAb03An03Vectors(ab03an03T *testing.T) {
	ab03an03Data, ab03an03Err := os.ReadFile("testdata/ab03an03_b524_codebook_views.json")
	if ab03an03Err != nil {
		ab03an03T.Fatal(ab03an03Err)
	}
	var ab03an03Vectors []ab03an03Vector
	if ab03an03Err = json.Unmarshal(ab03an03Data, &ab03an03Vectors); ab03an03Err != nil {
		ab03an03T.Fatal(ab03an03Err)
	}
	if len(ab03an03Vectors) != 88 {
		ab03an03T.Fatalf("ab03an03 vector count: %d", len(ab03an03Vectors))
	}
	for _, ab03an03Vector := range ab03an03Vectors {
		ab03an03T.Run(ab03an03Vector.Ab03An03Name, func(ab03an03T *testing.T) {
			ab03an03Input := ab03an03Input()
			ab03an03Input.Ab03An03Family = ab03an03Vector.Ab03An03Family
			ab03an03Input.Ab03An03GG = ab03an03Vector.Ab03An03GG
			ab03an03Input.Ab03An03RequestedGG = ab03an03Vector.Ab03An03GG
			ab03an03Input.Ab03An03ReturnedGG = ab03an03Vector.Ab03An03GG
			ab03an03Input.Ab03An03RR = ab03an03Vector.Ab03An03RR
			ab03an03Input.Ab03An03RequestedRR = ab03an03Vector.Ab03An03RR
			ab03an03Input.Ab03An03ReturnedRR = ab03an03Vector.Ab03An03RR
			ab03an03Input.Ab03An03Width = ab03an03Vector.Ab03An03Width
			ab03an03Input.Ab03An03Raw = ab03an03Vector.Ab03An03Raw
			ab03an03Result := Ab03An03Interpret(ab03an03Input)
			if !ab03an03Result.Ab03An03Applicable || ab03an03Result.Ab03An03MalformedWidth || ab03an03Result.Ab03An03Grade != "Hypothesis" || ab03an03Result.Ab03An03NativeQualification != "Unknown" || ab03an03Result.Ab03An03Integer == nil || *ab03an03Result.Ab03An03Integer != ab03an03Vector.Ab03An03Integer || ab03an03Result.Ab03An03CandidateLabel == nil || *ab03an03Result.Ab03An03CandidateLabel != ab03an03Vector.Ab03An03Label {
				ab03an03T.Fatalf("ab03an03 unexpected result: %+v", ab03an03Result)
			}
			if !reflect.DeepEqual(ab03an03Result.Ab03An03Context, ab03an03Input) {
				ab03an03T.Fatal("ab03an03 context changed")
			}
		})
	}
}

func TestAb03An03Controls(ab03an03T *testing.T) {
	for _, ab03an03Case := range []struct {
		ab03an03Name       string
		ab03an03Change     func(*Ab03An03Input)
		ab03an03Integer    *int
		ab03an03Applicable bool
		ab03an03Malformed  bool
	}{
		{"ab03an03-family720-pool", func(ab03an03I *Ab03An03Input) { ab03an03I.Ab03An03Family = "controller-720" }, ab03an03Int(5), true, false},
		{"ab03an03-family-unknown", func(ab03an03I *Ab03An03Input) { ab03an03I.Ab03An03Family = "unknown" }, ab03an03Int(5), false, false},
		{"ab03an03-family-missing", func(ab03an03I *Ab03An03Input) { ab03an03I.Ab03An03Family = "" }, ab03an03Int(5), false, false},
		{"ab03an03-op06", func(ab03an03I *Ab03An03Input) { ab03an03I.Ab03An03OP = 6 }, ab03an03Int(5), false, false},
		{"ab03an03-gg", func(ab03an03I *Ab03An03Input) { ab03an03I.Ab03An03GG = 4 }, ab03an03Int(5), false, false},
		{"ab03an03-rr", func(ab03an03I *Ab03An03Input) { ab03an03I.Ab03An03RR = 4 }, ab03an03Int(5), false, false},
		{"ab03an03-request-gg", func(ab03an03I *Ab03An03Input) { ab03an03I.Ab03An03RequestedGG = 4 }, ab03an03Int(5), false, false},
		{"ab03an03-request-rr", func(ab03an03I *Ab03An03Input) { ab03an03I.Ab03An03RequestedRR = 4 }, ab03an03Int(5), false, false},
		{"ab03an03-return-gg", func(ab03an03I *Ab03An03Input) { ab03an03I.Ab03An03ReturnedGG = 4 }, ab03an03Int(5), false, false},
		{"ab03an03-return-rr", func(ab03an03I *Ab03An03Input) { ab03an03I.Ab03An03ReturnedRR = 4 }, ab03an03Int(5), false, false},
		{"ab03an03-ii-missing", func(ab03an03I *Ab03An03Input) { ab03an03I.Ab03An03II = nil }, ab03an03Int(5), false, false},
		{"ab03an03-width-missing", func(ab03an03I *Ab03An03Input) { ab03an03I.Ab03An03Width = "" }, nil, false, false},
		{"ab03an03-status00", func(ab03an03I *Ab03An03Input) { ab03an03I.Ab03An03ResponseState = "status00" }, ab03an03Int(5), false, false},
		{"ab03an03-statusFF", func(ab03an03I *Ab03An03Input) { ab03an03I.Ab03An03ResponseState = "statusFF" }, ab03an03Int(5), false, false},
		{"ab03an03-transport", func(ab03an03I *Ab03An03Input) { ab03an03I.Ab03An03ResponseState = "transport-failure" }, ab03an03Int(5), false, false},
		{"ab03an03-unknown-code", func(ab03an03I *Ab03An03Input) { ab03an03I.Ab03An03Raw = []byte{42, 0} }, ab03an03Int(42), true, false},
		{"ab03an03-max-u16", func(ab03an03I *Ab03An03Input) { ab03an03I.Ab03An03Raw = []byte{255, 255} }, ab03an03Int(65535), true, false},
		{"ab03an03-little-endian", func(ab03an03I *Ab03An03Input) { ab03an03I.Ab03An03Raw = []byte{0, 1} }, ab03an03Int(256), true, false},
		{"ab03an03-truncated", func(ab03an03I *Ab03An03Input) { ab03an03I.Ab03An03Raw = []byte{5} }, nil, false, true},
		{"ab03an03-extra", func(ab03an03I *Ab03An03Input) { ab03an03I.Ab03An03Raw = []byte{5, 0, 0} }, nil, false, true},
		{"ab03an03-empty", func(ab03an03I *Ab03An03Input) { ab03an03I.Ab03An03Raw = nil }, nil, false, true},
	} {
		ab03an03T.Run(ab03an03Case.ab03an03Name, func(ab03an03T *testing.T) {
			ab03an03Input := ab03an03Input()
			ab03an03Case.ab03an03Change(&ab03an03Input)
			ab03an03Result := Ab03An03Interpret(ab03an03Input)
			ab03an03Grade := "Unknown"
			if ab03an03Case.ab03an03Applicable {
				ab03an03Grade = "Hypothesis"
			}
			if !reflect.DeepEqual(ab03an03Result.Ab03An03Integer, ab03an03Case.ab03an03Integer) || ab03an03Result.Ab03An03Applicable != ab03an03Case.ab03an03Applicable || ab03an03Result.Ab03An03MalformedWidth != ab03an03Case.ab03an03Malformed || ab03an03Result.Ab03An03CandidateLabel != nil || ab03an03Result.Ab03An03Grade != ab03an03Grade || ab03an03Result.Ab03An03NativeQualification != "Unknown" || !reflect.DeepEqual(ab03an03Result.Ab03An03Context, ab03an03Input) {
				ab03an03T.Fatalf("ab03an03 result: %+v", ab03an03Result)
			}
		})
	}
}

func ab03an03Int(ab03an03Value int) *int { return &ab03an03Value }

func TestAb03An03Alternatives(ab03an03T *testing.T) {
	for _, ab03an03Case := range []struct {
		ab03an03Name       string
		ab03an03Family     string
		ab03an03RR         uint16
		ab03an03Width      string
		ab03an03Raw        []byte
		ab03an03Integer    *int
		ab03an03Applicable bool
		ab03an03Malformed  bool
	}{
		{"ab03an03-720-remote-eight", "controller-720", 19, "u16-le", []byte{9, 0}, ab03an03Int(9), true, false},
		{"ab03an03-720-party", "controller-720", 27, "u8", []byte{3}, ab03an03Int(3), true, false},
		{"ab03an03-no-width-inference", "controller-700", 27, "u16-le", []byte{3}, nil, false, true},
		{"ab03an03-u8-extra", "controller-700", 27, "u8", []byte{3, 0}, nil, false, true},
		{"ab03an03-u8-max", "controller-700", 27, "u8", []byte{255}, ab03an03Int(255), true, false},
		{"ab03an03-signed-max", "controller-700", 28, "signed8", []byte{127}, ab03an03Int(127), true, false},
		{"ab03an03-signed-min", "controller-700", 28, "signed8", []byte{128}, ab03an03Int(-128), true, false},
		{"ab03an03-no-unsigned255", "controller-720", 28, "u8", []byte{255}, ab03an03Int(255), false, false},
		{"ab03an03-signed-extra", "controller-720", 28, "signed8", []byte{255, 0}, nil, false, true},
		{"ab03an03-unknown-selector", "controller-700", 99, "u16-le", []byte{0, 0}, ab03an03Int(0), false, false},
	} {
		ab03an03T.Run(ab03an03Case.ab03an03Name, func(ab03an03T *testing.T) {
			ab03an03Input := ab03an03Input()
			ab03an03Input.Ab03An03Family = ab03an03Case.ab03an03Family
			ab03an03Input.Ab03An03GG = 3
			ab03an03Input.Ab03An03RequestedGG = 3
			ab03an03Input.Ab03An03ReturnedGG = 3
			ab03an03Input.Ab03An03RR = ab03an03Case.ab03an03RR
			ab03an03Input.Ab03An03RequestedRR = ab03an03Case.ab03an03RR
			ab03an03Input.Ab03An03ReturnedRR = ab03an03Case.ab03an03RR
			ab03an03Input.Ab03An03Width = ab03an03Case.ab03an03Width
			ab03an03Input.Ab03An03Raw = ab03an03Case.ab03an03Raw
			ab03an03Result := Ab03An03Interpret(ab03an03Input)
			if !reflect.DeepEqual(ab03an03Result.Ab03An03Integer, ab03an03Case.ab03an03Integer) || ab03an03Result.Ab03An03Applicable != ab03an03Case.ab03an03Applicable || ab03an03Result.Ab03An03MalformedWidth != ab03an03Case.ab03an03Malformed || ab03an03Result.Ab03An03CandidateLabel != nil {
				ab03an03T.Fatalf("ab03an03 result: %+v", ab03an03Result)
			}
		})
	}
}

func TestAb03An03ContextOwnership(ab03an03T *testing.T) {
	for _, ab03an03Field := range []string{"Ab03An03Software", "Ab03An03Hardware", "Ab03An03APIVersion", "Ab03An03APIRevision"} {
		ab03an03T.Run("ab03an03-"+ab03an03Field, func(ab03an03T *testing.T) {
			ab03an03Input := ab03an03Input()
			reflect.ValueOf(&ab03an03Input).Elem().FieldByName(ab03an03Field).SetString("synthetic-alternative")
			ab03an03Result := Ab03An03Interpret(ab03an03Input)
			if !reflect.DeepEqual(ab03an03Result.Ab03An03Context, ab03an03Input) || !ab03an03Result.Ab03An03Applicable || ab03an03Result.Ab03An03CandidateLabel == nil || *ab03an03Result.Ab03An03CandidateLabel != "pool" {
				ab03an03T.Fatal("ab03an03 independent identity lost or used as inferred selector")
			}
		})
	}
	for _, ab03an03II := range []byte{1, 255} {
		ab03an03Input := ab03an03Input()
		ab03an03Input.Ab03An03II = &ab03an03II
		if ab03an03Result := Ab03An03Interpret(ab03an03Input); !ab03an03Result.Ab03An03Applicable || *ab03an03Result.Ab03An03Context.Ab03An03II != ab03an03II {
			ab03an03T.Fatal("ab03an03 explicit dynamic II lost")
		}
		ab03an03Input.Ab03An03GG = 1
		ab03an03Input.Ab03An03RequestedGG = 1
		ab03an03Input.Ab03An03ReturnedGG = 1
		ab03an03Input.Ab03An03RR = 3
		ab03an03Input.Ab03An03RequestedRR = 3
		ab03an03Input.Ab03An03ReturnedRR = 3
		if ab03an03Result := Ab03An03Interpret(ab03an03Input); ab03an03Result.Ab03An03Applicable {
			ab03an03T.Fatal("ab03an03 fixed II00 broadened")
		}
	}
	for _, ab03an03State := range []string{"echo-payload", "transport-failure"} {
		ab03an03Input := ab03an03Input()
		ab03an03Input.Ab03An03ResponseState = ab03an03State
		ab03an03First := Ab03An03Interpret(ab03an03Input)
		ab03an03Second := Ab03An03Interpret(ab03an03Input)
		ab03an03Input.Ab03An03Raw[0] = 99
		*ab03an03Input.Ab03An03II = 99
		if ab03an03First.Ab03An03Context.Ab03An03Raw[0] != 5 || *ab03an03First.Ab03An03Context.Ab03An03II != 0 {
			ab03an03T.Fatal("ab03an03 input aliases output")
		}
		ab03an03First.Ab03An03Context.Ab03An03Raw[0] = 88
		*ab03an03First.Ab03An03Context.Ab03An03II = 88
		*ab03an03First.Ab03An03Integer = 88
		if ab03an03First.Ab03An03CandidateLabel != nil {
			*ab03an03First.Ab03An03CandidateLabel = "ab03an03-mutated"
		}
		if ab03an03Input.Ab03An03Raw[0] != 99 || *ab03an03Input.Ab03An03II != 99 {
			ab03an03T.Fatal("ab03an03 output aliases input")
		}
		if ab03an03Second.Ab03An03Context.Ab03An03Raw[0] != 5 || *ab03an03Second.Ab03An03Context.Ab03An03II != 0 || *ab03an03Second.Ab03An03Integer != 5 {
			ab03an03T.Fatal("ab03an03 outputs alias")
		}
		if ab03an03Second.Ab03An03CandidateLabel != nil && *ab03an03Second.Ab03An03CandidateLabel != "pool" {
			ab03an03T.Fatal("ab03an03 label aliases")
		}
	}
}

func TestAb03An03CandidateSurface(ab03an03T *testing.T) {
	ab03an03Type := reflect.TypeOf(Ab03An03Result{})
	ab03an03Fields := []string{"Ab03An03Context", "Ab03An03Integer", "Ab03An03CandidateLabel", "Ab03An03Applicable", "Ab03An03MalformedWidth", "Ab03An03Grade", "Ab03An03NativeQualification"}
	if ab03an03Type.NumField() != len(ab03an03Fields) {
		ab03an03T.Fatal("ab03an03 candidate surface gained an unsupported policy or permission field")
	}
	for ab03an03Index, ab03an03Field := range ab03an03Fields {
		if ab03an03Type.Field(ab03an03Index).Name != ab03an03Field {
			ab03an03T.Fatal("ab03an03 candidate surface changed")
		}
	}
}
