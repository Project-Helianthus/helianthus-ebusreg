package system

import (
	ab64JSON "encoding/json"
	ab64OS "os"
	ab64Reflect "reflect"
	ab64Testing "testing"
)

type ab64Fixture struct {
	Ab64Label    string  `json:"ab64Label"`
	Ab64Family   string  `json:"ab64Family"`
	Ab64Register uint16  `json:"ab64Register"`
	Ab64Unit     string  `json:"ab64Unit"`
	Ab64Raw      []byte  `json:"ab64Raw"`
	Ab64Expected float32 `json:"ab64Expected"`
}

func ab64Valid(ab64Row ab64Fixture) Ab64Input {
	return Ab64Input{
		Ab64Context: Ab64Context{
			Ab64Family:      Ab64Field[string]{true, ab64Row.Ab64Family},
			Ab64Software:    Ab64Field[string]{true, "ab64-software"},
			Ab64Hardware:    Ab64Field[string]{true, "ab64-hardware"},
			Ab64APIVersion:  Ab64Field[string]{true, "ab64-api"},
			Ab64APIRevision: Ab64Field[string]{true, "ab64-revision"},
			Ab64Opcode:      Ab64Field[byte]{true, 2}, Ab64Group: Ab64Field[byte]{true, 3},
			Ab64Instance: Ab64Field[byte]{true, 0}, Ab64Register: Ab64Field[uint16]{true, ab64Row.Ab64Register},
		},
		Ab64ReplyKind: "echo_payload", Ab64ReturnedGroup: Ab64Field[byte]{true, 3},
		Ab64ReturnedRegister: Ab64Field[uint16]{true, ab64Row.Ab64Register},
		Ab64CodecGrade:       Ab64Hypothesis, Ab64Raw: append([]byte(nil), ab64Row.Ab64Raw...),
	}
}

func ab64Check(ab64T *ab64Testing.T, ab64Input Ab64Input, ab64Candidate bool, ab64Expected float32, ab64Unit string) Ab64Result {
	ab64T.Helper()
	ab64Result := Ab64Interpret(ab64Input)
	if !ab64Reflect.DeepEqual(ab64Result.Ab64Input, ab64Input) {
		ab64T.Fatalf("ab64 lost context/raw: got %#v want %#v", ab64Result.Ab64Input, ab64Input)
	}
	if ab64Result.Ab64Unit != ab64Unit {
		ab64T.Fatalf("ab64 unit: %q want %q", ab64Result.Ab64Unit, ab64Unit)
	}
	if ab64Candidate {
		if ab64Result.Ab64State != Ab64Candidate || !ab64Result.Ab64Value.Ab64Present || ab64Result.Ab64Value.Ab64Value != ab64Expected || ab64Result.Ab64Grade != Ab64Hypothesis || ab64Result.Ab64Codec != "float32_le" || ab64Result.Ab64InstanceRule != "dynamic" {
			ab64T.Fatalf("ab64 candidate: %#v", ab64Result)
		}
	} else if ab64Result.Ab64State != Ab64RawOnly || ab64Result.Ab64Value != (Ab64Field[float32]{}) || ab64Result.Ab64Grade != Ab64Unknown {
		ab64T.Fatalf("ab64 raw-only: %#v", ab64Result)
	}
	for _, ab64Grade := range []Ab64Grade{ab64Result.Ab64NativeSupport, ab64Result.Ab64Qualification, ab64Result.Ab64Presence, ab64Result.Ab64Freshness, ab64Result.Ab64Execution, ab64Result.Ab64Effects, ab64Result.Ab64Persistence, ab64Result.Ab64Delivery} {
		if ab64Grade != Ab64Unknown {
			ab64T.Fatal("ab64 native property promoted")
		}
	}
	return ab64Result
}

func TestAb64Bindings(ab64T *ab64Testing.T) {
	ab64Data, ab64Err := ab64OS.ReadFile("testdata/ab64_b524_zone_scalars.json")
	if ab64Err != nil {
		ab64T.Fatal(ab64Err)
	}
	var ab64Rows []ab64Fixture
	if ab64Err = ab64JSON.Unmarshal(ab64Data, &ab64Rows); ab64Err != nil {
		ab64T.Fatal(ab64Err)
	}
	ab64Manifest := map[string][]uint16{
		"controller-720": {2, 5, 9, 15, 20, 21, 35, 40}, "controller-700": {2, 5, 9, 15, 20, 21, 40},
	}
	ab64Seen := make(map[string]map[uint16]bool)
	if len(ab64Rows) != 15 {
		ab64T.Fatal("ab64 fixture must contain all 15 tuples")
	}
	for _, ab64Row := range ab64Rows {
		ab64Found := false
		for _, ab64Register := range ab64Manifest[ab64Row.Ab64Family] {
			if ab64Register == ab64Row.Ab64Register {
				ab64Found = true
			}
		}
		if ab64Seen[ab64Row.Ab64Family] == nil {
			ab64Seen[ab64Row.Ab64Family] = make(map[uint16]bool)
		}
		if !ab64Found || ab64Seen[ab64Row.Ab64Family][ab64Row.Ab64Register] {
			ab64T.Fatal("ab64 unexpected/duplicate tuple")
		}
		ab64Seen[ab64Row.Ab64Family][ab64Row.Ab64Register] = true
		ab64Unit := "°C"
		if ab64Row.Ab64Register == 40 {
			ab64Unit = "%"
		}
		if ab64Row.Ab64Expected != 1.5 || ab64Row.Ab64Unit != ab64Unit || !ab64Reflect.DeepEqual(ab64Row.Ab64Raw, []byte{0, 0, 192, 63}) {
			ab64T.Fatal("ab64 independent fixture expectation")
		}
		ab64T.Run(ab64Row.Ab64Label, func(ab64T *ab64Testing.T) {
			ab64Check(ab64T, ab64Valid(ab64Row), true, 1.5, ab64Unit)
			ab64Cases := []struct {
				ab64Name      string
				ab64Change    func(*Ab64Input)
				ab64Candidate bool
				ab64Expected  float32
				ab64Unit      string
			}{
				{"ab64-negative", func(ab64I *Ab64Input) { ab64I.Ab64Raw = []byte{0, 0, 128, 191} }, true, -1, ab64Unit},
				{"ab64-short", func(ab64I *Ab64Input) { ab64I.Ab64Raw = []byte{0, 0, 192} }, false, 0, ab64Unit},
				{"ab64-extra", func(ab64I *Ab64Input) { ab64I.Ab64Raw = append(ab64I.Ab64Raw, 0) }, false, 0, ab64Unit},
				{"ab64-width-two", func(ab64I *Ab64Input) { ab64I.Ab64Raw = []byte{52, 18} }, false, 0, ab64Unit},
				{"ab64-empty", func(ab64I *Ab64Input) { ab64I.Ab64Raw = nil }, false, 0, ab64Unit},
				{"ab64-inf", func(ab64I *Ab64Input) { ab64I.Ab64Raw = []byte{0, 0, 128, 127} }, false, 0, ab64Unit},
				{"ab64-negative-inf", func(ab64I *Ab64Input) { ab64I.Ab64Raw = []byte{0, 0, 128, 255} }, false, 0, ab64Unit},
				{"ab64-nan", func(ab64I *Ab64Input) { ab64I.Ab64Raw = []byte{0, 0, 192, 127} }, false, 0, ab64Unit},
				{"ab64-software", func(ab64I *Ab64Input) { ab64I.Ab64Context.Ab64Software.Ab64Value = "ab64-other" }, true, 1.5, ab64Unit},
				{"ab64-hardware", func(ab64I *Ab64Input) { ab64I.Ab64Context.Ab64Hardware.Ab64Value = "ab64-other" }, true, 1.5, ab64Unit},
				{"ab64-api-version", func(ab64I *Ab64Input) { ab64I.Ab64Context.Ab64APIVersion.Ab64Value = "ab64-other" }, true, 1.5, ab64Unit},
				{"ab64-api-revision", func(ab64I *Ab64Input) { ab64I.Ab64Context.Ab64APIRevision.Ab64Value = "ab64-other" }, true, 1.5, ab64Unit},
				{"ab64-instance-one", func(ab64I *Ab64Input) { ab64I.Ab64Context.Ab64Instance.Ab64Value = 1 }, true, 1.5, ab64Unit},
				{"ab64-instance-max", func(ab64I *Ab64Input) { ab64I.Ab64Context.Ab64Instance.Ab64Value = 255 }, true, 1.5, ab64Unit},
				{"ab64-family", func(ab64I *Ab64Input) { ab64I.Ab64Context.Ab64Family.Ab64Value = "ab64-unknown" }, false, 0, "Unknown"},
				{"ab64-op06", func(ab64I *Ab64Input) { ab64I.Ab64Context.Ab64Opcode.Ab64Value = 6 }, false, 0, "Unknown"},
				{"ab64-group", func(ab64I *Ab64Input) { ab64I.Ab64Context.Ab64Group.Ab64Value = 4 }, false, 0, "Unknown"},
				{"ab64-register-neighbor", func(ab64I *Ab64Input) {
					ab64I.Ab64Context.Ab64Register.Ab64Value = ab64Row.Ab64Register + 1
					if ab64Row.Ab64Register == 20 {
						ab64I.Ab64Context.Ab64Register.Ab64Value = 19
					}
				}, false, 0, "Unknown"},
				{"ab64-returned-group", func(ab64I *Ab64Input) { ab64I.Ab64ReturnedGroup.Ab64Value = 4 }, false, 0, ab64Unit},
				{"ab64-returned-register", func(ab64I *Ab64Input) { ab64I.Ab64ReturnedRegister.Ab64Value++ }, false, 0, ab64Unit},
				{"ab64-returned-group-absent", func(ab64I *Ab64Input) { ab64I.Ab64ReturnedGroup.Ab64Present = false }, false, 0, ab64Unit},
				{"ab64-returned-register-absent", func(ab64I *Ab64Input) { ab64I.Ab64ReturnedRegister.Ab64Present = false }, false, 0, ab64Unit},
				{"ab64-grade-proven", func(ab64I *Ab64Input) { ab64I.Ab64CodecGrade = Ab64Proven }, true, 1.5, ab64Unit},
				{"ab64-grade-unknown", func(ab64I *Ab64Input) { ab64I.Ab64CodecGrade = Ab64Unknown }, false, 0, ab64Unit},
				{"ab64-grade-missing", func(ab64I *Ab64Input) { ab64I.Ab64CodecGrade = "" }, false, 0, ab64Unit},
				{"ab64-grade-invalid", func(ab64I *Ab64Input) { ab64I.Ab64CodecGrade = "ab64-invalid" }, false, 0, ab64Unit},
			}
			for _, ab64Case := range ab64Cases {
				ab64T.Run(ab64Case.ab64Name, func(ab64T *ab64Testing.T) {
					ab64Input := ab64Valid(ab64Row)
					ab64Case.ab64Change(&ab64Input)
					ab64Check(ab64T, ab64Input, ab64Case.ab64Candidate, ab64Case.ab64Expected, ab64Case.ab64Unit)
				})
			}
			for _, ab64Kind := range []string{"", "status_only", "short", "failed", "mismatched", "unknown", "ab64-other"} {
				ab64T.Run("ab64-kind-"+ab64Kind, func(ab64T *ab64Testing.T) {
					ab64Input := ab64Valid(ab64Row)
					ab64Input.Ab64ReplyKind = ab64Kind
					ab64Check(ab64T, ab64Input, false, 0, ab64Unit)
				})
			}
			ab64T.Run("ab64-missing-context", func(ab64T *ab64Testing.T) {
				for ab64Index := 0; ab64Index < 9; ab64Index++ {
					ab64Input := ab64Valid(ab64Row)
					ab64Field := ab64Reflect.ValueOf(&ab64Input.Ab64Context).Elem().Field(ab64Index)
					ab64Field.FieldByName("Ab64Present").SetBool(false)
					ab64Check(ab64T, ab64Input, false, 0, "Unknown")
				}
			})
			ab64T.Run("ab64-explicit-zero-context", func(ab64T *ab64Testing.T) {
				for ab64Index := 0; ab64Index < 9; ab64Index++ {
					ab64Input := ab64Valid(ab64Row)
					ab64Field := ab64Reflect.ValueOf(&ab64Input.Ab64Context).Elem().Field(ab64Index).FieldByName("Ab64Value")
					ab64Field.Set(ab64Reflect.Zero(ab64Field.Type()))
					ab64Candidate := ab64Index >= 1 && ab64Index <= 4 || ab64Index == 7
					ab64ExpectedUnit := "Unknown"
					if ab64Candidate {
						ab64ExpectedUnit = ab64Unit
					}
					ab64Check(ab64T, ab64Input, ab64Candidate, 1.5, ab64ExpectedUnit)
				}
			})
			ab64T.Run("ab64-other-family", func(ab64T *ab64Testing.T) {
				ab64Input := ab64Valid(ab64Row)
				if ab64Row.Ab64Family == "controller-720" {
					ab64Input.Ab64Context.Ab64Family.Ab64Value = "controller-700"
				} else {
					ab64Input.Ab64Context.Ab64Family.Ab64Value = "controller-720"
				}
				ab64ExpectedUnit := ab64Unit
				if ab64Row.Ab64Register == 35 {
					ab64ExpectedUnit = "Unknown"
				}
				ab64Check(ab64T, ab64Input, ab64Row.Ab64Register != 35, 1.5, ab64ExpectedUnit)
			})
			ab64T.Run("ab64-owned-memory", func(ab64T *ab64Testing.T) {
				for _, ab64Kind := range []string{"echo_payload", "status_only"} {
					ab64Input := ab64Valid(ab64Row)
					ab64Input.Ab64ReplyKind = ab64Kind
					ab64Result := ab64Check(ab64T, ab64Input, ab64Kind == "echo_payload", 1.5, ab64Unit)
					ab64Input.Ab64Raw[0] = 42
					if ab64Result.Ab64Input.Ab64Raw[0] != 0 {
						ab64T.Fatal("ab64 input aliases result")
					}
					ab64Result.Ab64Input.Ab64Raw[1] = 43
					if ab64Input.Ab64Raw[1] != 0 {
						ab64T.Fatal("ab64 result aliases input")
					}
					ab64Again := Ab64Interpret(ab64Valid(ab64Row))
					if ab64Again.Ab64Input.Ab64Raw[1] != 0 {
						ab64T.Fatal("ab64 results share storage")
					}
				}
			})
		})
	}
}
