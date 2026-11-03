package system

import (
	ab58hex "encoding/hex"
	ab58json "encoding/json"
	ab58math "math"
	ab58os "os"
	ab58reflect "reflect"
	ab58testing "testing"
)

type ab58Fixture struct {
	Ab58Label    string `json:"ab58_label"`
	Ab58Family   string `json:"ab58_family"`
	Ab58Group    byte   `json:"ab58_group"`
	Ab58Register uint16 `json:"ab58_register"`
	Ab58Codec    string `json:"ab58_codec"`
	Ab58Unit     string `json:"ab58_unit"`
	Ab58Examples []struct {
		Ab58Label    string  `json:"ab58_label"`
		Ab58Hex      string  `json:"ab58_hex"`
		Ab58Expected float64 `json:"ab58_expected"`
	} `json:"ab58_examples"`
}

func ab58Fixtures(ab58t *ab58testing.T) []ab58Fixture {
	ab58t.Helper()
	ab58bytes, ab58err := ab58os.ReadFile("testdata/ab58_b524_system_scalars.json")
	if ab58err != nil {
		ab58t.Fatal(ab58err)
	}
	var ab58rows []ab58Fixture
	if ab58err = ab58json.Unmarshal(ab58bytes, &ab58rows); ab58err != nil {
		ab58t.Fatal(ab58err)
	}
	ab58expected := []ab58Fixture{
		{Ab58Family: "controller-720", Ab58Group: 0, Ab58Register: 0x0002, Ab58Codec: "float32_le", Ab58Unit: "°C"},
		{Ab58Family: "controller-720", Ab58Group: 0, Ab58Register: 0x0022, Ab58Codec: "float32_le", Ab58Unit: "°C"},
		{Ab58Family: "controller-720", Ab58Group: 0, Ab58Register: 0x0023, Ab58Codec: "float32_le", Ab58Unit: "°C"},
		{Ab58Family: "controller-720", Ab58Group: 0, Ab58Register: 0x0026, Ab58Codec: "float32_le", Ab58Unit: "°C"},
		{Ab58Family: "controller-720", Ab58Group: 0, Ab58Register: 0x0038, Ab58Codec: "float32_le", Ab58Unit: "°C"},
		{Ab58Family: "controller-720", Ab58Group: 0, Ab58Register: 0x0046, Ab58Codec: "float32_le", Ab58Unit: "°C"},
		{Ab58Family: "controller-720", Ab58Group: 0, Ab58Register: 0x004B, Ab58Codec: "float32_le", Ab58Unit: "°C"},
		{Ab58Family: "controller-720", Ab58Group: 0, Ab58Register: 0x0067, Ab58Codec: "float32_le", Ab58Unit: "°C"},
		{Ab58Family: "controller-720", Ab58Group: 0, Ab58Register: 0x0068, Ab58Codec: "float32_le", Ab58Unit: "°C"},
		{Ab58Family: "controller-720", Ab58Group: 0, Ab58Register: 0x0073, Ab58Codec: "float32_le", Ab58Unit: "°C"},
		{Ab58Family: "controller-720", Ab58Group: 0, Ab58Register: 0x0095, Ab58Codec: "float32_le", Ab58Unit: "°C"},
		{Ab58Family: "controller-720", Ab58Group: 0, Ab58Register: 0x009D, Ab58Codec: "float32_le", Ab58Unit: "°C"},
		{Ab58Family: "controller-720", Ab58Group: 0, Ab58Register: 0x009E, Ab58Codec: "float32_le", Ab58Unit: "°C"},
		{Ab58Family: "controller-720", Ab58Group: 0, Ab58Register: 0x009F, Ab58Codec: "float32_le", Ab58Unit: "°C"},
		{Ab58Family: "controller-720", Ab58Group: 0, Ab58Register: 0x00A0, Ab58Codec: "float32_le", Ab58Unit: "°C"},
		{Ab58Family: "controller-720", Ab58Group: 0, Ab58Register: 0x00EF, Ab58Codec: "float32_le", Ab58Unit: "°C"},
		{Ab58Family: "controller-720", Ab58Group: 0, Ab58Register: 0x0004, Ab58Codec: "uint16_le", Ab58Unit: "min"},
		{Ab58Family: "controller-720", Ab58Group: 0, Ab58Register: 0x0017, Ab58Codec: "uint16_le", Ab58Unit: "min"},
		{Ab58Family: "controller-720", Ab58Group: 0, Ab58Register: 0x0018, Ab58Codec: "uint16_le", Ab58Unit: "min"},
		{Ab58Family: "controller-720", Ab58Group: 0, Ab58Register: 0x0029, Ab58Codec: "float32_le", Ab58Unit: "K"},
		{Ab58Family: "controller-720", Ab58Group: 0, Ab58Register: 0x002F, Ab58Codec: "uint16_le", Ab58Unit: "Unknown"},
		{Ab58Family: "controller-720", Ab58Group: 0, Ab58Register: 0x0036, Ab58Codec: "uint16_le", Ab58Unit: "Unknown"},
		{Ab58Family: "controller-720", Ab58Group: 0, Ab58Register: 0x0039, Ab58Codec: "float32_le", Ab58Unit: "bar"},
		{Ab58Family: "controller-700", Ab58Group: 0, Ab58Register: 0x0002, Ab58Codec: "float32_le", Ab58Unit: "°C"},
		{Ab58Family: "controller-700", Ab58Group: 0, Ab58Register: 0x0022, Ab58Codec: "float32_le", Ab58Unit: "°C"},
		{Ab58Family: "controller-700", Ab58Group: 0, Ab58Register: 0x0023, Ab58Codec: "float32_le", Ab58Unit: "°C"},
		{Ab58Family: "controller-700", Ab58Group: 0, Ab58Register: 0x0026, Ab58Codec: "float32_le", Ab58Unit: "°C"},
		{Ab58Family: "controller-700", Ab58Group: 0, Ab58Register: 0x0038, Ab58Codec: "float32_le", Ab58Unit: "°C"},
		{Ab58Family: "controller-700", Ab58Group: 0, Ab58Register: 0x0046, Ab58Codec: "float32_le", Ab58Unit: "°C"},
		{Ab58Family: "controller-700", Ab58Group: 0, Ab58Register: 0x004B, Ab58Codec: "float32_le", Ab58Unit: "°C"},
		{Ab58Family: "controller-700", Ab58Group: 0, Ab58Register: 0x0067, Ab58Codec: "float32_le", Ab58Unit: "°C"},
		{Ab58Family: "controller-700", Ab58Group: 0, Ab58Register: 0x0068, Ab58Codec: "float32_le", Ab58Unit: "°C"},
		{Ab58Family: "controller-700", Ab58Group: 0, Ab58Register: 0x0073, Ab58Codec: "float32_le", Ab58Unit: "°C"},
		{Ab58Family: "controller-700", Ab58Group: 0, Ab58Register: 0x0095, Ab58Codec: "float32_le", Ab58Unit: "°C"},
		{Ab58Family: "controller-700", Ab58Group: 0, Ab58Register: 0x009D, Ab58Codec: "float32_le", Ab58Unit: "°C"},
		{Ab58Family: "controller-700", Ab58Group: 0, Ab58Register: 0x009E, Ab58Codec: "float32_le", Ab58Unit: "°C"},
		{Ab58Family: "controller-700", Ab58Group: 0, Ab58Register: 0x009F, Ab58Codec: "float32_le", Ab58Unit: "°C"},
		{Ab58Family: "controller-700", Ab58Group: 0, Ab58Register: 0x00A0, Ab58Codec: "float32_le", Ab58Unit: "°C"},
		{Ab58Family: "controller-700", Ab58Group: 0, Ab58Register: 0x0006, Ab58Codec: "uint16_le", Ab58Unit: "Unknown"},
		{Ab58Family: "controller-700", Ab58Group: 0, Ab58Register: 0x002F, Ab58Codec: "uint16_le", Ab58Unit: "Unknown"},
		{Ab58Family: "controller-700", Ab58Group: 0, Ab58Register: 0x0036, Ab58Codec: "uint16_le", Ab58Unit: "Unknown"},
		{Ab58Family: "controller-700", Ab58Group: 0, Ab58Register: 0x0017, Ab58Codec: "uint16_le", Ab58Unit: "min"},
		{Ab58Family: "controller-700", Ab58Group: 0, Ab58Register: 0x0029, Ab58Codec: "float32_le", Ab58Unit: "K"},
		{Ab58Family: "controller-700", Ab58Group: 0, Ab58Register: 0x0039, Ab58Codec: "float32_le", Ab58Unit: "bar"},
		{Ab58Family: "controller-700", Ab58Group: 0, Ab58Register: 0x0042, Ab58Codec: "uint16_le", Ab58Unit: "Unknown"},
		{Ab58Family: "controller-700", Ab58Group: 0, Ab58Register: 0x0043, Ab58Codec: "uint16_le", Ab58Unit: "Unknown"},
	}
	if len(ab58rows) != len(ab58expected) {
		ab58t.Fatal("ab58 incomplete fixture set")
	}
	for ab58index, ab58row := range ab58rows {
		ab58want := ab58expected[ab58index]
		if ab58row.Ab58Family != ab58want.Ab58Family || ab58row.Ab58Group != ab58want.Ab58Group || ab58row.Ab58Register != ab58want.Ab58Register || ab58row.Ab58Codec != ab58want.Ab58Codec || ab58row.Ab58Unit != ab58want.Ab58Unit || len(ab58row.Ab58Examples) != 2 {
			ab58t.Fatal("ab58 fixture selector or metadata drift")
		}
		ab58hexes := []string{"0000c03f", "000080bf"}
		ab58values := []float64{1.5, -1}
		if ab58want.Ab58Codec == "uint16_le" {
			ab58hexes = []string{"3412", "ffff"}
			ab58values = []float64{4660, 65535}
		}
		for ab58index, ab58example := range ab58row.Ab58Examples {
			if ab58example.Ab58Hex != ab58hexes[ab58index] || ab58example.Ab58Expected != ab58values[ab58index] {
				ab58t.Fatal("ab58 independent arithmetic control drift")
			}
		}
	}
	return ab58rows
}

func ab58Control(ab58t *ab58testing.T, ab58row ab58Fixture, ab58hexadecimal string) Ab58Input {
	ab58t.Helper()
	ab58raw, ab58err := ab58hex.DecodeString(ab58hexadecimal)
	if ab58err != nil {
		ab58t.Fatal(ab58err)
	}
	return Ab58Input{
		Ab58Context: Ab58Context{
			Ab58Family:      Ab58Field[string]{true, ab58row.Ab58Family},
			Ab58Software:    Ab58Field[string]{true, "ab58_software"},
			Ab58Hardware:    Ab58Field[string]{true, "ab58_hardware"},
			Ab58APIVersion:  Ab58Field[string]{true, "ab58_api"},
			Ab58APIRevision: Ab58Field[string]{true, "ab58_revision"},
			Ab58Opcode:      Ab58Field[byte]{true, 2}, Ab58Group: Ab58Field[byte]{true, ab58row.Ab58Group},
			Ab58Instance: Ab58Field[byte]{true, 0}, Ab58Register: Ab58Field[uint16]{true, ab58row.Ab58Register},
		},
		Ab58ReplyKind:        Ab58Field[string]{true, "echo_payload"},
		Ab58ReturnedGroup:    Ab58Field[byte]{true, ab58row.Ab58Group},
		Ab58ReturnedRegister: Ab58Field[uint16]{true, ab58row.Ab58Register},
		Ab58CodecGrade:       Ab58Field[string]{true, "Hypothesis"}, Ab58Raw: ab58raw,
	}
}

func ab58Check(ab58t *ab58testing.T, ab58input Ab58Input, ab58candidate bool) Ab58Result {
	ab58t.Helper()
	ab58result := Ab58InterpretB524Scalars(ab58input)
	if !ab58reflect.DeepEqual(ab58result.Ab58Input, ab58input) {
		ab58t.Fatalf("ab58 caller evidence lost: got %+v want %+v", ab58result.Ab58Input, ab58input)
	}
	if ab58result.Ab58Presence != "Unknown" || ab58result.Ab58Freshness != "Unknown" || ab58result.Ab58Execution != "Unknown" || ab58result.Ab58Support != "Unknown" || ab58result.Ab58Effects != "Unknown" || ab58result.Ab58Persistence != "Unknown" || ab58result.Ab58Delivery != "Unknown" {
		ab58t.Fatal("ab58 invented physical qualification")
	}
	if ab58candidate {
		if ab58result.Ab58State != Ab58Candidate || ab58result.Ab58Value == nil || ab58result.Ab58Grade != "Hypothesis" {
			ab58t.Fatalf("ab58 missing hypothesis: %+v", ab58result)
		}
	} else if ab58result.Ab58State != Ab58RawOnly || ab58result.Ab58Value != nil || ab58result.Ab58Grade != "Unknown" {
		ab58t.Fatalf("ab58 raw-only admission failed: %+v", ab58result)
	}
	return ab58result
}

func TestAb58Bindings(ab58t *ab58testing.T) {
	for _, ab58row := range ab58Fixtures(ab58t) {
		ab58t.Run(ab58row.Ab58Label, func(ab58t *ab58testing.T) {
			for _, ab58example := range ab58row.Ab58Examples {
				ab58t.Run(ab58example.Ab58Label, func(ab58t *ab58testing.T) {
					ab58result := ab58Check(ab58t, ab58Control(ab58t, ab58row, ab58example.Ab58Hex), true)
					if ab58result.Ab58Codec != ab58row.Ab58Codec || ab58result.Ab58DeclaredUnit != ab58row.Ab58Unit {
						ab58t.Fatal("ab58 metadata mismatch")
					}
					ab58value := float64(ab58result.Ab58Value.Ab58Float32)
					if ab58row.Ab58Codec == "uint16_le" {
						ab58value = float64(ab58result.Ab58Value.Ab58Uint16)
					}
					if ab58value != ab58example.Ab58Expected {
						ab58t.Fatalf("ab58 value %v want %v", ab58value, ab58example.Ab58Expected)
					}
				})
			}
		})
	}
}

func TestAb58IndependentAxes(ab58t *ab58testing.T) {
	for _, ab58row := range ab58Fixtures(ab58t) {
		ab58t.Run(ab58row.Ab58Label, func(ab58t *ab58testing.T) {
			for _, ab58axis := range []struct {
				ab58name      string
				ab58candidate bool
				ab58mutate    func(*Ab58Input)
			}{
				{"ab58_family_unknown", false, func(ab58input *Ab58Input) { ab58input.Ab58Context.Ab58Family.Ab58Value = "ab58_unknown" }},
				{"ab58_family_other", (ab58row.Ab58Register != 4 && ab58row.Ab58Register != 0x18 && ab58row.Ab58Register != 0xef && ab58row.Ab58Register != 6 && ab58row.Ab58Register != 0x42 && ab58row.Ab58Register != 0x43), func(ab58input *Ab58Input) {
					if ab58row.Ab58Family == "controller-720" {
						ab58input.Ab58Context.Ab58Family.Ab58Value = "controller-700"
					} else {
						ab58input.Ab58Context.Ab58Family.Ab58Value = "controller-720"
					}
				}},
				{"ab58_software", true, func(ab58input *Ab58Input) { ab58input.Ab58Context.Ab58Software.Ab58Value = "ab58_changed" }},
				{"ab58_hardware", true, func(ab58input *Ab58Input) { ab58input.Ab58Context.Ab58Hardware.Ab58Value = "ab58_changed" }},
				{"ab58_api_version", true, func(ab58input *Ab58Input) { ab58input.Ab58Context.Ab58APIVersion.Ab58Value = "ab58_changed" }},
				{"ab58_api_revision", true, func(ab58input *Ab58Input) { ab58input.Ab58Context.Ab58APIRevision.Ab58Value = "ab58_changed" }},
				{"ab58_opcode", false, func(ab58input *Ab58Input) { ab58input.Ab58Context.Ab58Opcode.Ab58Value = 6 }},
				{"ab58_group", false, func(ab58input *Ab58Input) { ab58input.Ab58Context.Ab58Group.Ab58Value = 0x7f }},
				{"ab58_instance", ab58row.Ab58Family == "controller-700" && (ab58row.Ab58Register == 0x42 || ab58row.Ab58Register == 0x43), func(ab58input *Ab58Input) { ab58input.Ab58Context.Ab58Instance.Ab58Value = 1 }},
				{"ab58_register", false, func(ab58input *Ab58Input) { ab58input.Ab58Context.Ab58Register.Ab58Value = ab58row.Ab58Register + 2 }},
				{"ab58_reply_kind", false, func(ab58input *Ab58Input) { ab58input.Ab58ReplyKind.Ab58Value = "status_only" }},
				{"ab58_returned_group", false, func(ab58input *Ab58Input) { ab58input.Ab58ReturnedGroup.Ab58Value = 0x7f }},
				{"ab58_returned_register", false, func(ab58input *Ab58Input) { ab58input.Ab58ReturnedRegister.Ab58Value++ }},
				{"ab58_grade_proven", true, func(ab58input *Ab58Input) { ab58input.Ab58CodecGrade.Ab58Value = "Proven" }},
				{"ab58_grade_unknown", false, func(ab58input *Ab58Input) { ab58input.Ab58CodecGrade.Ab58Value = "Unknown" }},
			} {
				ab58t.Run(ab58axis.ab58name, func(ab58t *ab58testing.T) {
					ab58input := ab58Control(ab58t, ab58row, ab58row.Ab58Examples[0].Ab58Hex)
					ab58axis.ab58mutate(&ab58input)
					ab58result := ab58Check(ab58t, ab58input, ab58axis.ab58candidate)
					if ab58axis.ab58candidate {
						if ab58result.Ab58Value.Ab58Float32 != 1.5 && ab58row.Ab58Codec == "float32_le" {
							ab58t.Fatal("ab58 context changed arithmetic")
						}
						if ab58result.Ab58Value.Ab58Uint16 != 4660 && ab58row.Ab58Codec == "uint16_le" {
							ab58t.Fatal("ab58 context changed arithmetic")
						}
					}
				})
			}
			ab58input := ab58Control(ab58t, ab58row, ab58row.Ab58Examples[0].Ab58Hex)
			for _, ab58container := range []string{"Ab58Context", "ab58_reply"} {
				ab58fields := ab58reflect.ValueOf(&ab58input).Elem()
				ab58start := 1
				ab58end := 5
				if ab58container == "Ab58Context" {
					ab58fields = ab58fields.FieldByName("Ab58Context")
					ab58start = 0
					ab58end = ab58fields.NumField()
				}
				for ab58index := ab58start; ab58index < ab58end; ab58index++ {
					ab58t.Run("ab58_missing_"+ab58fields.Type().Field(ab58index).Name, func(ab58t *ab58testing.T) {
						ab58presence := ab58fields.Field(ab58index).FieldByName("Ab58Present")
						ab58presence.SetBool(false)
						ab58Check(ab58t, ab58input, false)
						ab58presence.SetBool(true)
					})
				}
			}
		})
	}
}

func TestAb58MalformedAndOwnership(ab58t *ab58testing.T) {
	for _, ab58row := range ab58Fixtures(ab58t) {
		ab58t.Run(ab58row.Ab58Label, func(ab58t *ab58testing.T) {
			for _, ab58kind := range []string{"status_only", "short", "failed", "mismatched", "unknown", "", "error", "ack", "request"} {
				ab58t.Run("ab58_kind_"+ab58kind, func(ab58t *ab58testing.T) {
					ab58input := ab58Control(ab58t, ab58row, ab58row.Ab58Examples[0].Ab58Hex)
					ab58input.Ab58ReplyKind.Ab58Value = ab58kind
					ab58Check(ab58t, ab58input, false)
				})
			}
			for _, ab58grade := range []string{"", "ab58_invalid", "hypothesis"} {
				ab58t.Run("ab58_grade_"+ab58grade, func(ab58t *ab58testing.T) {
					ab58input := ab58Control(ab58t, ab58row, ab58row.Ab58Examples[0].Ab58Hex)
					ab58input.Ab58CodecGrade.Ab58Value = ab58grade
					ab58Check(ab58t, ab58input, false)
				})
			}
			for _, ab58variant := range []string{"ab58_minus_one", "ab58_plus_one", "ab58_nil", "ab58_empty", "ab58_wrong_codec_width", "ab58_unlisted_group", "ab58_neighbor_low", "ab58_neighbor_high"} {
				ab58t.Run(ab58variant, func(ab58t *ab58testing.T) {
					ab58input := ab58Control(ab58t, ab58row, ab58row.Ab58Examples[0].Ab58Hex)
					switch ab58variant {
					case "ab58_minus_one":
						ab58input.Ab58Raw = ab58input.Ab58Raw[:len(ab58input.Ab58Raw)-1]
					case "ab58_plus_one":
						ab58input.Ab58Raw = append(ab58input.Ab58Raw, 0)
					case "ab58_nil":
						ab58input.Ab58Raw = nil
					case "ab58_empty":
						ab58input.Ab58Raw = []byte{}
					case "ab58_wrong_codec_width":
						if ab58row.Ab58Codec == "float32_le" {
							ab58input.Ab58Raw = []byte{0x34, 0x12}
						} else {
							ab58input.Ab58Raw = []byte{0, 0, 0xc0, 0x3f}
						}
					case "ab58_unlisted_group":
						ab58input.Ab58Context.Ab58Group.Ab58Value = 0x7f
						ab58input.Ab58ReturnedGroup.Ab58Value = 0x7f
					case "ab58_neighbor_low":
						ab58input.Ab58Context.Ab58Register.Ab58Value = 1
						ab58input.Ab58ReturnedRegister = ab58input.Ab58Context.Ab58Register
					case "ab58_neighbor_high":
						ab58input.Ab58Context.Ab58Register.Ab58Value = 0xf0
						ab58input.Ab58ReturnedRegister = ab58input.Ab58Context.Ab58Register
					}
					ab58result := ab58Check(ab58t, ab58input, false)
					if ab58variant != "ab58_unlisted_group" && ab58variant != "ab58_neighbor_low" && ab58variant != "ab58_neighbor_high" && ab58result.Ab58DeclaredUnit != ab58row.Ab58Unit {
						ab58t.Fatal("ab58 declared unit lost on partial failure")
					}
				})
			}
			if ab58row.Ab58Codec == "float32_le" {
				for _, ab58bits := range []string{"0000807f", "000080ff", "0000c07f"} {
					ab58t.Run("ab58_nonfinite_"+ab58bits, func(ab58t *ab58testing.T) { ab58Check(ab58t, ab58Control(ab58t, ab58row, ab58bits), false) })
				}
			}
			for _, ab58candidate := range []bool{true, false} {
				ab58name := "ab58_memory_candidate"
				if !ab58candidate {
					ab58name = "ab58_memory_raw_only"
				}
				ab58t.Run(ab58name, func(ab58t *ab58testing.T) {
					ab58input := ab58Control(ab58t, ab58row, ab58row.Ab58Examples[0].Ab58Hex)
					if !ab58candidate {
						ab58input.Ab58ReplyKind.Ab58Value = "failed"
					}
					ab58first := ab58Check(ab58t, ab58input, ab58candidate)
					ab58second := ab58Check(ab58t, ab58input, ab58candidate)
					ab58original := ab58input.Ab58Raw[0]
					ab58input.Ab58Raw[0] ^= 0xff
					ab58input.Ab58Context.Ab58Software.Ab58Value = "ab58_mutated"
					if ab58first.Ab58Input.Ab58Raw[0] != ab58original || ab58first.Ab58Input.Ab58Context.Ab58Software.Ab58Value != "ab58_software" {
						ab58t.Fatal("ab58 input aliases result")
					}
					ab58first.Ab58Input.Ab58Raw[0] ^= 0x55
					if ab58input.Ab58Raw[0] != ab58original^0xff || ab58second.Ab58Input.Ab58Raw[0] != ab58original {
						ab58t.Fatal("ab58 returned memory aliases caller or another result")
					}
					if ab58candidate {
						ab58first.Ab58Value.Ab58Float32 = 99
						ab58first.Ab58Value.Ab58Uint16 = 99
						if ab58second.Ab58Value.Ab58Float32 == 99 || ab58second.Ab58Value.Ab58Uint16 == 99 {
							ab58t.Fatal("ab58 scalar aliases another result")
						}
					}
				})
			}
		})
	}
}

func TestAb58ExplicitZeroContext(ab58t *ab58testing.T) {
	for _, ab58row := range ab58Fixtures(ab58t) {
		ab58t.Run(ab58row.Ab58Label, func(ab58t *ab58testing.T) {
			for _, ab58field := range []string{"Ab58Software", "Ab58Hardware", "Ab58APIVersion", "Ab58APIRevision"} {
				ab58t.Run("ab58_zero_"+ab58field, func(ab58t *ab58testing.T) {
					ab58input := ab58Control(ab58t, ab58row, ab58row.Ab58Examples[0].Ab58Hex)
					ab58reflect.ValueOf(&ab58input.Ab58Context).Elem().FieldByName(ab58field).FieldByName("Ab58Value").SetString("")
					ab58Check(ab58t, ab58input, true)
				})
			}
		})
	}
}

func TestAb58AllConcreteInstances(ab58t *ab58testing.T) {
	for _, ab58row := range ab58Fixtures(ab58t) {
		ab58t.Run(ab58row.Ab58Label, func(ab58t *ab58testing.T) {
			for ab58instance := 0; ab58instance < 256; ab58instance++ {
				ab58input := ab58Control(ab58t, ab58row, ab58row.Ab58Examples[0].Ab58Hex)
				ab58input.Ab58Context.Ab58Instance.Ab58Value = byte(ab58instance)
				ab58Check(ab58t, ab58input, ab58instance == 0 || (ab58row.Ab58Family == "controller-700" && (ab58row.Ab58Register == 0x42 || ab58row.Ab58Register == 0x43)))
			}
		})
	}
}

func TestAb58ExactSelectorAdmission(ab58t *ab58testing.T) {
	ab58rows := ab58Fixtures(ab58t)
	for _, ab58family := range []string{"controller-720", "controller-700", "ab58_unknown"} {
		ab58t.Run("ab58_"+ab58family, func(ab58t *ab58testing.T) {
			for ab58register := 0; ab58register <= 0xffff; ab58register++ {
				ab58listed := false
				for _, ab58row := range ab58rows {
					if ab58row.Ab58Family == ab58family && ab58row.Ab58Register == uint16(ab58register) {
						ab58listed = true
					}
				}
				if ab58listed {
					continue
				}
				for _, ab58raw := range []string{"3412", "0000c03f"} {
					ab58input := ab58Control(ab58t, ab58rows[0], ab58raw)
					ab58input.Ab58Context.Ab58Family.Ab58Value = ab58family
					ab58input.Ab58Context.Ab58Register.Ab58Value = uint16(ab58register)
					ab58input.Ab58ReturnedRegister = ab58input.Ab58Context.Ab58Register
					ab58Check(ab58t, ab58input, false)
				}
			}
		})
	}
}

func TestAb58FiniteValues(ab58t *ab58testing.T) {
	for _, ab58row := range ab58Fixtures(ab58t) {
		if ab58row.Ab58Codec != "float32_le" {
			continue
		}
		ab58t.Run(ab58row.Ab58Label, func(ab58t *ab58testing.T) {
			for _, ab58case := range []struct {
				ab58label string
				ab58hex   string
				ab58value float32
				ab58bits  uint32
			}{
				{"ab58_max_finite", "ffff7f7f", 3.4028234663852886e38, 0x7f7fffff},
				{"ab58_min_finite", "ffff7fff", -3.4028234663852886e38, 0xff7fffff},
				{"ab58_subnormal", "01000000", 1.401298464324817e-45, 1},
				{"ab58_zero", "00000000", 0, 0},
				{"ab58_negative_zero", "00000080", 0, 0x80000000},
			} {
				ab58t.Run(ab58case.ab58label, func(ab58t *ab58testing.T) {
					ab58result := ab58Check(ab58t, ab58Control(ab58t, ab58row, ab58case.ab58hex), true)
					if ab58result.Ab58Value.Ab58Float32 != ab58case.ab58value || ab58math.Float32bits(ab58result.Ab58Value.Ab58Float32) != ab58case.ab58bits {
						ab58t.Fatal("ab58 finite value clipped, offset or normalized")
					}
				})
			}
		})
	}
}
