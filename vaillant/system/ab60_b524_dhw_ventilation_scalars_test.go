package system

import (
	ab60hex "encoding/hex"
	ab60json "encoding/json"
	ab60os "os"
	ab60reflect "reflect"
	ab60testing "testing"
)

type ab60Fixture struct {
	Ab60Label    string `json:"ab60_label"`
	Ab60Family   string `json:"ab60_family"`
	Ab60Group    byte   `json:"ab60_group"`
	Ab60Register uint16 `json:"ab60_register"`
	Ab60Codec    string `json:"ab60_codec"`
	Ab60Unit     string `json:"ab60_unit"`
	Ab60Examples []struct {
		Ab60Label    string  `json:"ab60_label"`
		Ab60Hex      string  `json:"ab60_hex"`
		Ab60Expected float64 `json:"ab60_expected"`
	} `json:"ab60_examples"`
}

func ab60Fixtures(ab60t *ab60testing.T) []ab60Fixture {
	ab60t.Helper()
	ab60bytes, ab60err := ab60os.ReadFile("testdata/ab60_b524_dhw_ventilation_scalars.json")
	if ab60err != nil {
		ab60t.Fatal(ab60err)
	}
	var ab60rows []ab60Fixture
	if ab60err = ab60json.Unmarshal(ab60bytes, &ab60rows); ab60err != nil {
		ab60t.Fatal(ab60err)
	}
	// Independently fix the complete selector set and arithmetic controls.
	ab60expected := []ab60Fixture{
		{Ab60Family: "controller-720", Ab60Group: 1, Ab60Register: 4, Ab60Codec: "float32_le", Ab60Unit: "°C"},
		{Ab60Family: "controller-720", Ab60Group: 1, Ab60Register: 5, Ab60Codec: "float32_le", Ab60Unit: "°C"},
		{Ab60Family: "controller-720", Ab60Group: 9, Ab60Register: 13, Ab60Codec: "uint16_le", Ab60Unit: "Unknown"},
		{Ab60Family: "controller-720", Ab60Group: 9, Ab60Register: 14, Ab60Codec: "uint16_le", Ab60Unit: "Unknown"},
		{Ab60Family: "controller-700", Ab60Group: 1, Ab60Register: 4, Ab60Codec: "float32_le", Ab60Unit: "°C"},
		{Ab60Family: "controller-700", Ab60Group: 1, Ab60Register: 5, Ab60Codec: "float32_le", Ab60Unit: "°C"},
	}
	if len(ab60rows) != len(ab60expected) {
		ab60t.Fatal("ab60 incomplete fixture set")
	}
	for ab60index, ab60row := range ab60rows {
		ab60want := ab60expected[ab60index]
		if ab60row.Ab60Family != ab60want.Ab60Family || ab60row.Ab60Group != ab60want.Ab60Group || ab60row.Ab60Register != ab60want.Ab60Register || ab60row.Ab60Codec != ab60want.Ab60Codec || ab60row.Ab60Unit != ab60want.Ab60Unit || len(ab60row.Ab60Examples) != 2 {
			ab60t.Fatal("ab60 fixture selector or metadata drift")
		}
		ab60hexes := []string{"0000c03f", "000080bf"}
		ab60values := []float64{1.5, -1}
		if ab60want.Ab60Codec == "uint16_le" {
			ab60hexes = []string{"3412", "ffff"}
			ab60values = []float64{4660, 65535}
		}
		for ab60index, ab60example := range ab60row.Ab60Examples {
			if ab60example.Ab60Hex != ab60hexes[ab60index] || ab60example.Ab60Expected != ab60values[ab60index] {
				ab60t.Fatal("ab60 independent arithmetic control drift")
			}
		}
	}
	return ab60rows
}

func ab60Control(ab60t *ab60testing.T, ab60row ab60Fixture, ab60hexadecimal string) Ab60Input {
	ab60t.Helper()
	ab60raw, ab60err := ab60hex.DecodeString(ab60hexadecimal)
	if ab60err != nil {
		ab60t.Fatal(ab60err)
	}
	return Ab60Input{
		Ab60Context: Ab60Context{
			Ab60Family:      Ab60Field[string]{true, ab60row.Ab60Family},
			Ab60Software:    Ab60Field[string]{true, "ab60_software"},
			Ab60Hardware:    Ab60Field[string]{true, "ab60_hardware"},
			Ab60APIVersion:  Ab60Field[string]{true, "ab60_api"},
			Ab60APIRevision: Ab60Field[string]{true, "ab60_revision"},
			Ab60Opcode:      Ab60Field[byte]{true, 2}, Ab60Group: Ab60Field[byte]{true, ab60row.Ab60Group},
			Ab60Instance: Ab60Field[byte]{true, 0}, Ab60Register: Ab60Field[uint16]{true, ab60row.Ab60Register},
		},
		Ab60ReplyKind:        Ab60Field[string]{true, "echo_payload"},
		Ab60ReturnedGroup:    Ab60Field[byte]{true, ab60row.Ab60Group},
		Ab60ReturnedRegister: Ab60Field[uint16]{true, ab60row.Ab60Register},
		Ab60CodecGrade:       Ab60Field[string]{true, "Hypothesis"}, Ab60Raw: ab60raw,
	}
}

func ab60Check(ab60t *ab60testing.T, ab60input Ab60Input, ab60candidate bool) Ab60Result {
	ab60t.Helper()
	ab60result := Ab60InterpretB524Scalars(ab60input)
	if !ab60reflect.DeepEqual(ab60result.Ab60Input, ab60input) {
		ab60t.Fatalf("ab60 caller evidence lost: got %+v want %+v", ab60result.Ab60Input, ab60input)
	}
	if ab60result.Ab60Presence != "Unknown" || ab60result.Ab60Freshness != "Unknown" || ab60result.Ab60Execution != "Unknown" || ab60result.Ab60Support != "Unknown" || ab60result.Ab60Effects != "Unknown" || ab60result.Ab60Persistence != "Unknown" || ab60result.Ab60Delivery != "Unknown" {
		ab60t.Fatal("ab60 invented physical qualification")
	}
	if ab60candidate {
		if ab60result.Ab60State != Ab60Candidate || ab60result.Ab60Value == nil || ab60result.Ab60Grade != "Hypothesis" {
			ab60t.Fatalf("ab60 missing hypothesis: %+v", ab60result)
		}
	} else if ab60result.Ab60State != Ab60RawOnly || ab60result.Ab60Value != nil || ab60result.Ab60Grade != "Unknown" {
		ab60t.Fatalf("ab60 raw-only admission failed: %+v", ab60result)
	}
	return ab60result
}

func TestAb60Bindings(ab60t *ab60testing.T) {
	for _, ab60row := range ab60Fixtures(ab60t) {
		ab60t.Run(ab60row.Ab60Label, func(ab60t *ab60testing.T) {
			for _, ab60example := range ab60row.Ab60Examples {
				ab60t.Run(ab60example.Ab60Label, func(ab60t *ab60testing.T) {
					ab60result := ab60Check(ab60t, ab60Control(ab60t, ab60row, ab60example.Ab60Hex), true)
					if ab60result.Ab60Codec != ab60row.Ab60Codec || ab60result.Ab60DeclaredUnit != ab60row.Ab60Unit {
						ab60t.Fatal("ab60 metadata mismatch")
					}
					ab60value := float64(ab60result.Ab60Value.Ab60Float32)
					if ab60row.Ab60Codec == "uint16_le" {
						ab60value = float64(ab60result.Ab60Value.Ab60Uint16)
					}
					if ab60value != ab60example.Ab60Expected {
						ab60t.Fatalf("ab60 value %v want %v", ab60value, ab60example.Ab60Expected)
					}
				})
			}
		})
	}
}

func TestAb60IndependentAxes(ab60t *ab60testing.T) {
	for _, ab60row := range ab60Fixtures(ab60t) {
		ab60t.Run(ab60row.Ab60Label, func(ab60t *ab60testing.T) {
			for _, ab60axis := range []struct {
				ab60name      string
				ab60candidate bool
				ab60mutate    func(*Ab60Input)
			}{
				{"ab60_family_unknown", false, func(ab60input *Ab60Input) { ab60input.Ab60Context.Ab60Family.Ab60Value = "ab60_unknown" }},
				{"ab60_family_other", ab60row.Ab60Group == 1, func(ab60input *Ab60Input) {
					if ab60row.Ab60Family == "controller-720" {
						ab60input.Ab60Context.Ab60Family.Ab60Value = "controller-700"
					} else {
						ab60input.Ab60Context.Ab60Family.Ab60Value = "controller-720"
					}
				}},
				{"ab60_software", true, func(ab60input *Ab60Input) { ab60input.Ab60Context.Ab60Software.Ab60Value = "ab60_changed" }},
				{"ab60_hardware", true, func(ab60input *Ab60Input) { ab60input.Ab60Context.Ab60Hardware.Ab60Value = "ab60_changed" }},
				{"ab60_api_version", true, func(ab60input *Ab60Input) { ab60input.Ab60Context.Ab60APIVersion.Ab60Value = "ab60_changed" }},
				{"ab60_api_revision", true, func(ab60input *Ab60Input) { ab60input.Ab60Context.Ab60APIRevision.Ab60Value = "ab60_changed" }},
				{"ab60_opcode", false, func(ab60input *Ab60Input) { ab60input.Ab60Context.Ab60Opcode.Ab60Value = 6 }},
				{"ab60_group", false, func(ab60input *Ab60Input) { ab60input.Ab60Context.Ab60Group.Ab60Value = 0x7f }},
				{"ab60_instance", false, func(ab60input *Ab60Input) { ab60input.Ab60Context.Ab60Instance.Ab60Value = 1 }},
				{"ab60_register", false, func(ab60input *Ab60Input) { ab60input.Ab60Context.Ab60Register.Ab60Value = ab60row.Ab60Register + 2 }},
				{"ab60_reply_kind", false, func(ab60input *Ab60Input) { ab60input.Ab60ReplyKind.Ab60Value = "status_only" }},
				{"ab60_returned_group", false, func(ab60input *Ab60Input) { ab60input.Ab60ReturnedGroup.Ab60Value = 0x7f }},
				{"ab60_returned_register", false, func(ab60input *Ab60Input) { ab60input.Ab60ReturnedRegister.Ab60Value++ }},
				{"ab60_grade_proven", true, func(ab60input *Ab60Input) { ab60input.Ab60CodecGrade.Ab60Value = "Proven" }},
				{"ab60_grade_unknown", false, func(ab60input *Ab60Input) { ab60input.Ab60CodecGrade.Ab60Value = "Unknown" }},
			} {
				ab60t.Run(ab60axis.ab60name, func(ab60t *ab60testing.T) {
					ab60input := ab60Control(ab60t, ab60row, ab60row.Ab60Examples[0].Ab60Hex)
					ab60axis.ab60mutate(&ab60input)
					ab60result := ab60Check(ab60t, ab60input, ab60axis.ab60candidate)
					if ab60axis.ab60candidate {
						if ab60result.Ab60Value.Ab60Float32 != 1.5 && ab60row.Ab60Codec == "float32_le" {
							ab60t.Fatal("ab60 context changed arithmetic")
						}
						if ab60result.Ab60Value.Ab60Uint16 != 4660 && ab60row.Ab60Codec == "uint16_le" {
							ab60t.Fatal("ab60 context changed arithmetic")
						}
					}
				})
			}
			// Every context and reply field has an independent presence bit, including II00.
			ab60input := ab60Control(ab60t, ab60row, ab60row.Ab60Examples[0].Ab60Hex)
			for _, ab60container := range []string{"Ab60Context", "ab60_reply"} {
				ab60fields := ab60reflect.ValueOf(&ab60input).Elem()
				ab60start := 1
				ab60end := 5
				if ab60container == "Ab60Context" {
					ab60fields = ab60fields.FieldByName("Ab60Context")
					ab60start = 0
					ab60end = ab60fields.NumField()
				}
				for ab60index := ab60start; ab60index < ab60end; ab60index++ {
					ab60t.Run("ab60_missing_"+ab60fields.Type().Field(ab60index).Name, func(ab60t *ab60testing.T) {
						ab60presence := ab60fields.Field(ab60index).FieldByName("Ab60Present")
						ab60presence.SetBool(false)
						ab60Check(ab60t, ab60input, false)
						ab60presence.SetBool(true)
					})
				}
			}
		})
	}
}

func TestAb60MalformedAndOwnership(ab60t *ab60testing.T) {
	for _, ab60row := range ab60Fixtures(ab60t) {
		ab60t.Run(ab60row.Ab60Label, func(ab60t *ab60testing.T) {
			for _, ab60kind := range []string{"status_only", "short", "failed", "mismatched", "unknown", "", "error", "ack", "request"} {
				ab60t.Run("ab60_kind_"+ab60kind, func(ab60t *ab60testing.T) {
					ab60input := ab60Control(ab60t, ab60row, ab60row.Ab60Examples[0].Ab60Hex)
					ab60input.Ab60ReplyKind.Ab60Value = ab60kind
					ab60Check(ab60t, ab60input, false)
				})
			}
			for _, ab60grade := range []string{"", "ab60_invalid", "hypothesis"} {
				ab60t.Run("ab60_grade_"+ab60grade, func(ab60t *ab60testing.T) {
					ab60input := ab60Control(ab60t, ab60row, ab60row.Ab60Examples[0].Ab60Hex)
					ab60input.Ab60CodecGrade.Ab60Value = ab60grade
					ab60Check(ab60t, ab60input, false)
				})
			}
			for _, ab60variant := range []string{"ab60_minus_one", "ab60_plus_one", "ab60_nil", "ab60_empty", "ab60_wrong_codec_width", "ab60_unlisted_group", "ab60_neighbor_low", "ab60_neighbor_high"} {
				ab60t.Run(ab60variant, func(ab60t *ab60testing.T) {
					ab60input := ab60Control(ab60t, ab60row, ab60row.Ab60Examples[0].Ab60Hex)
					switch ab60variant {
					case "ab60_minus_one":
						ab60input.Ab60Raw = ab60input.Ab60Raw[:len(ab60input.Ab60Raw)-1]
					case "ab60_plus_one":
						ab60input.Ab60Raw = append(ab60input.Ab60Raw, 0)
					case "ab60_nil":
						ab60input.Ab60Raw = nil
					case "ab60_empty":
						ab60input.Ab60Raw = []byte{}
					case "ab60_wrong_codec_width":
						if ab60row.Ab60Codec == "float32_le" {
							ab60input.Ab60Raw = []byte{0x34, 0x12}
						} else {
							ab60input.Ab60Raw = []byte{0, 0, 0xc0, 0x3f}
						}
					case "ab60_unlisted_group":
						ab60input.Ab60Context.Ab60Group.Ab60Value = 0x7f
						ab60input.Ab60ReturnedGroup.Ab60Value = 0x7f
					case "ab60_neighbor_low":
						if ab60row.Ab60Group == 1 {
							ab60input.Ab60Context.Ab60Register.Ab60Value = 3
						} else {
							ab60input.Ab60Context.Ab60Register.Ab60Value = 12
						}
						ab60input.Ab60ReturnedRegister = ab60input.Ab60Context.Ab60Register
					case "ab60_neighbor_high":
						if ab60row.Ab60Group == 1 {
							ab60input.Ab60Context.Ab60Register.Ab60Value = 6
						} else {
							ab60input.Ab60Context.Ab60Register.Ab60Value = 15
						}
						ab60input.Ab60ReturnedRegister = ab60input.Ab60Context.Ab60Register
					}
					ab60result := ab60Check(ab60t, ab60input, false)
					if ab60variant != "ab60_unlisted_group" && ab60variant != "ab60_neighbor_low" && ab60variant != "ab60_neighbor_high" && ab60result.Ab60DeclaredUnit != ab60row.Ab60Unit {
						ab60t.Fatal("ab60 declared unit lost on partial failure")
					}
				})
			}
			if ab60row.Ab60Codec == "float32_le" {
				for _, ab60bits := range []string{"0000807f", "000080ff", "0000c07f"} {
					ab60t.Run("ab60_nonfinite_"+ab60bits, func(ab60t *ab60testing.T) { ab60Check(ab60t, ab60Control(ab60t, ab60row, ab60bits), false) })
				}
			}
			for _, ab60candidate := range []bool{true, false} {
				ab60name := "ab60_memory_candidate"
				if !ab60candidate {
					ab60name = "ab60_memory_raw_only"
				}
				ab60t.Run(ab60name, func(ab60t *ab60testing.T) {
					ab60input := ab60Control(ab60t, ab60row, ab60row.Ab60Examples[0].Ab60Hex)
					if !ab60candidate {
						ab60input.Ab60ReplyKind.Ab60Value = "failed"
					}
					ab60first := ab60Check(ab60t, ab60input, ab60candidate)
					ab60second := ab60Check(ab60t, ab60input, ab60candidate)
					ab60original := ab60input.Ab60Raw[0]
					ab60input.Ab60Raw[0] ^= 0xff
					ab60input.Ab60Context.Ab60Software.Ab60Value = "ab60_mutated"
					if ab60first.Ab60Input.Ab60Raw[0] != ab60original || ab60first.Ab60Input.Ab60Context.Ab60Software.Ab60Value != "ab60_software" {
						ab60t.Fatal("ab60 input aliases result")
					}
					ab60first.Ab60Input.Ab60Raw[0] ^= 0x55
					if ab60input.Ab60Raw[0] != ab60original^0xff || ab60second.Ab60Input.Ab60Raw[0] != ab60original {
						ab60t.Fatal("ab60 returned memory aliases caller or another result")
					}
					if ab60candidate {
						ab60first.Ab60Value.Ab60Float32 = 99
						ab60first.Ab60Value.Ab60Uint16 = 99
						if ab60second.Ab60Value.Ab60Float32 == 99 || ab60second.Ab60Value.Ab60Uint16 == 99 {
							ab60t.Fatal("ab60 scalar aliases another result")
						}
					}
				})
			}
		})
	}
}

func TestAb60ExplicitZeroContext(ab60t *ab60testing.T) {
	for _, ab60row := range ab60Fixtures(ab60t) {
		ab60t.Run(ab60row.Ab60Label, func(ab60t *ab60testing.T) {
			for _, ab60field := range []string{"Ab60Software", "Ab60Hardware", "Ab60APIVersion", "Ab60APIRevision"} {
				ab60t.Run("ab60_zero_"+ab60field, func(ab60t *ab60testing.T) {
					ab60input := ab60Control(ab60t, ab60row, ab60row.Ab60Examples[0].Ab60Hex)
					ab60reflect.ValueOf(&ab60input.Ab60Context).Elem().FieldByName(ab60field).FieldByName("Ab60Value").SetString("")
					ab60Check(ab60t, ab60input, true)
				})
			}
		})
	}
}
