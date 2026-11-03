package system

import (
	ab62JSON "encoding/json"
	ab62OS "os"
	ab62Reflect "reflect"
	ab62Testing "testing"
)

type ab62Fixture struct {
	Ab62Label    string
	Ab62Family   string
	Ab62Register uint16
	Ab62Unit     string
	Ab62Raw      []byte
	Ab62Expected float32
}

func ab62Present[Ab62T any](ab62Value Ab62T) Ab62Field[Ab62T] {
	return Ab62Field[Ab62T]{Ab62Present: true, Ab62Value: ab62Value}
}

func ab62Fixtures(ab62T *ab62Testing.T) []ab62Fixture {
	ab62T.Helper()
	ab62Bytes, ab62Err := ab62OS.ReadFile("testdata/ab62_b524_circuit_scalars.json")
	if ab62Err != nil {
		ab62T.Fatal(ab62Err)
	}
	var ab62Rows []ab62Fixture
	if ab62Err = ab62JSON.Unmarshal(ab62Bytes, &ab62Rows); ab62Err != nil {
		ab62T.Fatal(ab62Err)
	}
	if len(ab62Rows) != 22 {
		ab62T.Fatalf("ab62 fixture count: %d", len(ab62Rows))
	}
	ab62Seen := map[string]map[uint16]bool{}
	for _, ab62Row := range ab62Rows {
		if ab62Seen[ab62Row.Ab62Family] == nil {
			ab62Seen[ab62Row.Ab62Family] = map[uint16]bool{}
		}
		if ab62Seen[ab62Row.Ab62Family][ab62Row.Ab62Register] {
			ab62T.Fatal("ab62 duplicate fixture")
		}
		ab62Seen[ab62Row.Ab62Family][ab62Row.Ab62Register] = true
		if ab62Row.Ab62Expected != 1.5 || !ab62Reflect.DeepEqual(ab62Row.Ab62Raw, []byte{0, 0, 192, 63}) {
			ab62T.Fatal("ab62 independent control changed")
		}
	}
	ab62Want := map[string][]uint16{
		"controller-720": {7, 8, 10, 11, 15, 16, 17, 18, 20, 26, 28, 31},
		"controller-700": {7, 8, 10, 11, 15, 16, 17, 18, 20, 31},
	}
	for ab62Family, ab62Registers := range ab62Want {
		for _, ab62Register := range ab62Registers {
			if !ab62Seen[ab62Family][ab62Register] {
				ab62T.Fatalf("ab62 missing tuple %s/%04X", ab62Family, ab62Register)
			}
		}
	}
	return ab62Rows
}

func ab62Control(ab62Row ab62Fixture) Ab62Input {
	return Ab62Input{
		Ab62Context: Ab62Context{
			Ab62Family:   ab62Present(ab62Row.Ab62Family),
			Ab62Software: ab62Present("ab62-software"), Ab62Hardware: ab62Present("ab62-hardware"),
			Ab62APIVersion: ab62Present("ab62-api"), Ab62APIRevision: ab62Present("ab62-revision"),
			Ab62Opcode: ab62Present(byte(2)), Ab62Group: ab62Present(byte(2)),
			Ab62Instance: ab62Present(byte(0)), Ab62Register: ab62Present(ab62Row.Ab62Register),
		},
		Ab62ReplyKind: ab62Present("echo_payload"), Ab62ReturnedGroup: ab62Present(byte(2)),
		Ab62ReturnedRegister: ab62Present(ab62Row.Ab62Register), Ab62CodecGrade: ab62Present(Ab62Hypothesis),
		Ab62Raw: append([]byte(nil), ab62Row.Ab62Raw...),
	}
}

func ab62Check(ab62T *ab62Testing.T, ab62Input Ab62Input, ab62Candidate bool, ab62Expected float32) Ab62Result {
	ab62T.Helper()
	ab62Result := Ab62InterpretCircuitScalar(ab62Input)
	if !ab62Reflect.DeepEqual(ab62Result.Ab62Input, ab62Input) {
		ab62T.Fatalf("ab62 context lost: %+v versus %+v", ab62Result.Ab62Input, ab62Input)
	}
	if ab62Candidate {
		if ab62Result.Ab62State != Ab62Candidate || ab62Result.Ab62Grade != Ab62Hypothesis || !ab62Result.Ab62Value.Ab62Present || ab62Result.Ab62Value.Ab62Value != ab62Expected || ab62Result.Ab62Codec != "float32_le" {
			ab62T.Fatalf("ab62 candidate: %+v", ab62Result)
		}
	} else if ab62Result.Ab62State != Ab62RawOnly || ab62Result.Ab62Grade != Ab62Unknown || ab62Result.Ab62Value != (Ab62Field[float32]{}) {
		ab62T.Fatalf("ab62 raw-only: %+v", ab62Result)
	}
	return ab62Result
}

func TestAb62EveryBinding(ab62T *ab62Testing.T) {
	for _, ab62Row := range ab62Fixtures(ab62T) {
		ab62T.Run(ab62Row.Ab62Label, func(ab62T *ab62Testing.T) {
			ab62Input := ab62Control(ab62Row)
			ab62Result := ab62Check(ab62T, ab62Input, true, ab62Row.Ab62Expected)
			if ab62Result.Ab62DeclaredUnit != ab62Row.Ab62Unit {
				ab62T.Fatalf("ab62 unit: %q", ab62Result.Ab62DeclaredUnit)
			}
			for _, ab62Case := range []struct {
				ab62Label string
				ab62Raw   []byte
				ab62Valid bool
				ab62Value float32
			}{
				{"ab62-negative", []byte{0, 0, 128, 191}, true, -1},
				{"ab62-positive-infinity", []byte{0, 0, 128, 127}, false, 0},
				{"ab62-negative-infinity", []byte{0, 0, 128, 255}, false, 0},
				{"ab62-nan", []byte{1, 0, 128, 127}, false, 0},
				{"ab62-short", []byte{0, 0, 192}, false, 0},
				{"ab62-extra", []byte{0, 0, 192, 63, 0}, false, 0},
				{"ab62-empty", nil, false, 0},
				{"ab62-present-empty", []byte{}, false, 0},
				{"ab62-no-width-selection", []byte{52, 18}, false, 0},
			} {
				ab62T.Run(ab62Case.ab62Label, func(ab62T *ab62Testing.T) {
					ab62Input.Ab62Raw = ab62Case.ab62Raw
					ab62Result := ab62Check(ab62T, ab62Input, ab62Case.ab62Valid, ab62Case.ab62Value)
					if ab62Result.Ab62DeclaredUnit != ab62Row.Ab62Unit {
						ab62T.Fatal("ab62 declared unit lost")
					}
				})
			}
		})
	}
}

func TestAb62IndependentAxes(ab62T *ab62Testing.T) {
	for _, ab62Row := range ab62Fixtures(ab62T) {
		ab62T.Run(ab62Row.Ab62Label, func(ab62T *ab62Testing.T) {
			for _, ab62Axis := range []struct {
				ab62Name      string
				ab62Value     any
				ab62Candidate bool
			}{
				{"Ab62Family", "ab62-unknown", false},
				{"Ab62Software", "ab62-other-software", true},
				{"Ab62Hardware", "ab62-other-hardware", true},
				{"Ab62APIVersion", "ab62-other-api", true},
				{"Ab62APIRevision", "ab62-other-revision", true},
				{"Ab62Opcode", byte(6), false},
				{"Ab62Group", byte(3), false},
				{"Ab62Instance", byte(1), true},
				{"Ab62Register", uint16(6), false},
				{"Ab62ReplyKind", "status_only", false},
				{"Ab62ReturnedGroup", byte(3), false},
				{"Ab62ReturnedRegister", uint16(6), false},
				{"Ab62CodecGrade", Ab62Proven, true},
			} {
				ab62T.Run(ab62Axis.ab62Name, func(ab62T *ab62Testing.T) {
					ab62Input := ab62Control(ab62Row)
					ab62Field := ab62Reflect.ValueOf(&ab62Input).Elem().FieldByName(ab62Axis.ab62Name)
					if !ab62Field.IsValid() {
						ab62Field = ab62Reflect.ValueOf(&ab62Input.Ab62Context).Elem().FieldByName(ab62Axis.ab62Name)
					}
					ab62Field.FieldByName("Ab62Value").Set(ab62Reflect.ValueOf(ab62Axis.ab62Value))
					ab62Check(ab62T, ab62Input, ab62Axis.ab62Candidate, 1.5)
					ab62Input = ab62Control(ab62Row)
					ab62Field = ab62Reflect.ValueOf(&ab62Input).Elem().FieldByName(ab62Axis.ab62Name)
					if !ab62Field.IsValid() {
						ab62Field = ab62Reflect.ValueOf(&ab62Input.Ab62Context).Elem().FieldByName(ab62Axis.ab62Name)
					}
					ab62Field.FieldByName("Ab62Present").SetBool(false)
					ab62Check(ab62T, ab62Input, false, 0)
				})
			}
			for _, ab62Name := range []string{"Ab62Software", "Ab62Hardware", "Ab62APIVersion", "Ab62APIRevision"} {
				ab62Input := ab62Control(ab62Row)
				ab62Field := ab62Reflect.ValueOf(&ab62Input.Ab62Context).Elem().FieldByName(ab62Name)
				ab62Field.FieldByName("Ab62Value").SetString("")
				ab62Check(ab62T, ab62Input, true, 1.5)
				ab62Field.FieldByName("Ab62Present").SetBool(false)
				ab62Check(ab62T, ab62Input, false, 0)
			}
			ab62Input := ab62Control(ab62Row)
			ab62Input.Ab62Context.Ab62Instance.Ab62Present = false
			ab62Check(ab62T, ab62Input, false, 0)
			for _, ab62Kind := range []string{"", "short", "failed", "unknown", "status_only", "ab62-unrecognized"} {
				ab62Input = ab62Control(ab62Row)
				ab62Input.Ab62ReplyKind = ab62Present(ab62Kind)
				ab62Check(ab62T, ab62Input, false, 0)
			}
			for _, ab62Grade := range []Ab62Grade{"", Ab62Unknown, "ab62-unrecognized"} {
				ab62Input = ab62Control(ab62Row)
				ab62Input.Ab62CodecGrade = ab62Present(ab62Grade)
				ab62Check(ab62T, ab62Input, false, 0)
			}
			for _, ab62Instance := range []byte{0, 1, 255} {
				ab62Input = ab62Control(ab62Row)
				ab62Input.Ab62Context.Ab62Instance = ab62Present(ab62Instance)
				ab62Check(ab62T, ab62Input, true, 1.5)
			}
			ab62Input = ab62Control(ab62Row)
			ab62Input.Ab62Context.Ab62Family.Ab62Value = "controller-700"
			if ab62Row.Ab62Family == "controller-700" {
				ab62Input.Ab62Context.Ab62Family.Ab62Value = "controller-720"
			}
			ab62Check(ab62T, ab62Input, ab62Row.Ab62Register != 0x1a && ab62Row.Ab62Register != 0x1c, 1.5)
			ab62Input = ab62Control(ab62Row)
			ab62Input.Ab62Context.Ab62Group = ab62Present(byte(3))
			ab62Input.Ab62ReturnedGroup = ab62Present(byte(3))
			ab62Check(ab62T, ab62Input, false, 0)
			ab62Input = ab62Control(ab62Row)
			ab62Input.Ab62Context.Ab62Opcode = ab62Present(byte(6))
			ab62Check(ab62T, ab62Input, false, 0)
			ab62Input = ab62Control(ab62Row)
			ab62Input.Ab62Context.Ab62Register = ab62Present(uint16(0x20))
			ab62Input.Ab62ReturnedRegister = ab62Present(uint16(0x20))
			ab62Check(ab62T, ab62Input, false, 0)
		})
	}
}

func TestAb62OwnedMemory(ab62T *ab62Testing.T) {
	for _, ab62Row := range ab62Fixtures(ab62T) {
		for _, ab62Kind := range []string{"echo_payload", "failed"} {
			ab62Input := ab62Control(ab62Row)
			ab62Input.Ab62ReplyKind = ab62Present(ab62Kind)
			ab62Result := ab62Check(ab62T, ab62Input, ab62Kind == "echo_payload", 1.5)
			ab62Input.Ab62Raw[0] = 255
			if ab62Result.Ab62Input.Ab62Raw[0] != 0 {
				ab62T.Fatal("ab62 aliases input")
			}
			ab62Result.Ab62Input.Ab62Raw[1] = 255
			if ab62Input.Ab62Raw[1] != 0 {
				ab62T.Fatal("ab62 aliases output")
			}
			ab62Other := Ab62InterpretCircuitScalar(ab62Control(ab62Row))
			if ab62Other.Ab62Input.Ab62Raw[1] != 0 {
				ab62T.Fatal("ab62 shared output")
			}
			if ab62Kind == "echo_payload" && ab62Result.Ab62Value.Ab62Value != 1.5 {
				ab62T.Fatal("ab62 value mutated")
			}
		}
	}
}
