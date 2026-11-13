package system

import (
	b524DHWVentilationScalarHex "encoding/hex"
	b524DHWVentilationScalarJSON "encoding/json"
	b524DHWVentilationScalarOs "os"
	b524DHWVentilationScalarReflect "reflect"
	b524DHWVentilationScalarTesting "testing"
)

type b524DHWVentilationScalarFixture struct {
	B524DHWVentilationScalarLabel    string `json:"ab60_label"`
	B524DHWVentilationScalarFamily   string `json:"ab60_family"`
	B524DHWVentilationScalarGroup    byte   `json:"ab60_group"`
	B524DHWVentilationScalarRegister uint16 `json:"ab60_register"`
	B524DHWVentilationScalarCodec    string `json:"ab60_codec"`
	B524DHWVentilationScalarUnit     string `json:"ab60_unit"`
	B524DHWVentilationScalarExamples []struct {
		B524DHWVentilationScalarLabel    string  `json:"ab60_label"`
		B524DHWVentilationScalarHex      string  `json:"ab60_hex"`
		B524DHWVentilationScalarExpected float64 `json:"ab60_expected"`
	} `json:"ab60_examples"`
}

func b524DHWVentilationScalarFixtures(b524DHWVentilationScalarT *b524DHWVentilationScalarTesting.T) []b524DHWVentilationScalarFixture {
	b524DHWVentilationScalarT.Helper()
	b524DHWVentilationScalarBytes, b524DHWVentilationScalarErr := b524DHWVentilationScalarOs.ReadFile("testdata/ab60_b524_dhw_ventilation_scalars.json")
	if b524DHWVentilationScalarErr != nil {
		b524DHWVentilationScalarT.Fatal(b524DHWVentilationScalarErr)
	}
	var b524DHWVentilationScalarRows []b524DHWVentilationScalarFixture
	if b524DHWVentilationScalarErr = b524DHWVentilationScalarJSON.Unmarshal(b524DHWVentilationScalarBytes, &b524DHWVentilationScalarRows); b524DHWVentilationScalarErr != nil {
		b524DHWVentilationScalarT.Fatal(b524DHWVentilationScalarErr)
	}
	// Independently fix the complete selector set and arithmetic controls.
	b524DHWVentilationScalarExpected := []b524DHWVentilationScalarFixture{
		{B524DHWVentilationScalarFamily: "controller-720", B524DHWVentilationScalarGroup: 1, B524DHWVentilationScalarRegister: 4, B524DHWVentilationScalarCodec: "float32_le", B524DHWVentilationScalarUnit: "°C"},
		{B524DHWVentilationScalarFamily: "controller-720", B524DHWVentilationScalarGroup: 1, B524DHWVentilationScalarRegister: 5, B524DHWVentilationScalarCodec: "float32_le", B524DHWVentilationScalarUnit: "°C"},
		{B524DHWVentilationScalarFamily: "controller-720", B524DHWVentilationScalarGroup: 9, B524DHWVentilationScalarRegister: 13, B524DHWVentilationScalarCodec: "uint16_le", B524DHWVentilationScalarUnit: "Unknown"},
		{B524DHWVentilationScalarFamily: "controller-720", B524DHWVentilationScalarGroup: 9, B524DHWVentilationScalarRegister: 14, B524DHWVentilationScalarCodec: "uint16_le", B524DHWVentilationScalarUnit: "Unknown"},
		{B524DHWVentilationScalarFamily: "controller-700", B524DHWVentilationScalarGroup: 1, B524DHWVentilationScalarRegister: 4, B524DHWVentilationScalarCodec: "float32_le", B524DHWVentilationScalarUnit: "°C"},
		{B524DHWVentilationScalarFamily: "controller-700", B524DHWVentilationScalarGroup: 1, B524DHWVentilationScalarRegister: 5, B524DHWVentilationScalarCodec: "float32_le", B524DHWVentilationScalarUnit: "°C"},
	}
	if len(b524DHWVentilationScalarRows) != len(b524DHWVentilationScalarExpected) {
		b524DHWVentilationScalarT.Fatal("ab60 incomplete fixture set")
	}
	for b524DHWVentilationScalarIndex, b524DHWVentilationScalarRow := range b524DHWVentilationScalarRows {
		b524DHWVentilationScalarWant := b524DHWVentilationScalarExpected[b524DHWVentilationScalarIndex]
		if b524DHWVentilationScalarRow.B524DHWVentilationScalarFamily != b524DHWVentilationScalarWant.B524DHWVentilationScalarFamily || b524DHWVentilationScalarRow.B524DHWVentilationScalarGroup != b524DHWVentilationScalarWant.B524DHWVentilationScalarGroup || b524DHWVentilationScalarRow.B524DHWVentilationScalarRegister != b524DHWVentilationScalarWant.B524DHWVentilationScalarRegister || b524DHWVentilationScalarRow.B524DHWVentilationScalarCodec != b524DHWVentilationScalarWant.B524DHWVentilationScalarCodec || b524DHWVentilationScalarRow.B524DHWVentilationScalarUnit != b524DHWVentilationScalarWant.B524DHWVentilationScalarUnit || len(b524DHWVentilationScalarRow.B524DHWVentilationScalarExamples) != 2 {
			b524DHWVentilationScalarT.Fatal("ab60 fixture selector or metadata drift")
		}
		b524DHWVentilationScalarHexes := []string{"0000c03f", "000080bf"}
		b524DHWVentilationScalarValues := []float64{1.5, -1}
		if b524DHWVentilationScalarWant.B524DHWVentilationScalarCodec == "uint16_le" {
			b524DHWVentilationScalarHexes = []string{"3412", "ffff"}
			b524DHWVentilationScalarValues = []float64{4660, 65535}
		}
		for b524DHWVentilationScalarIndex, b524DHWVentilationScalarExample := range b524DHWVentilationScalarRow.B524DHWVentilationScalarExamples {
			if b524DHWVentilationScalarExample.B524DHWVentilationScalarHex != b524DHWVentilationScalarHexes[b524DHWVentilationScalarIndex] || b524DHWVentilationScalarExample.B524DHWVentilationScalarExpected != b524DHWVentilationScalarValues[b524DHWVentilationScalarIndex] {
				b524DHWVentilationScalarT.Fatal("ab60 independent arithmetic control drift")
			}
		}
	}
	return b524DHWVentilationScalarRows
}

func b524DHWVentilationScalarControl(b524DHWVentilationScalarT *b524DHWVentilationScalarTesting.T, b524DHWVentilationScalarRow b524DHWVentilationScalarFixture, b524DHWVentilationScalarHexadecimal string) B524DHWVentilationScalarInput {
	b524DHWVentilationScalarT.Helper()
	b524DHWVentilationScalarRaw, b524DHWVentilationScalarErr := b524DHWVentilationScalarHex.DecodeString(b524DHWVentilationScalarHexadecimal)
	if b524DHWVentilationScalarErr != nil {
		b524DHWVentilationScalarT.Fatal(b524DHWVentilationScalarErr)
	}
	return B524DHWVentilationScalarInput{
		B524DHWVentilationScalarContext: B524DHWVentilationScalarContext{
			B524DHWVentilationScalarFamily:      B524DHWVentilationScalarField[string]{true, b524DHWVentilationScalarRow.B524DHWVentilationScalarFamily},
			B524DHWVentilationScalarSoftware:    B524DHWVentilationScalarField[string]{true, "ab60_software"},
			B524DHWVentilationScalarHardware:    B524DHWVentilationScalarField[string]{true, "ab60_hardware"},
			B524DHWVentilationScalarAPIVersion:  B524DHWVentilationScalarField[string]{true, "ab60_api"},
			B524DHWVentilationScalarAPIRevision: B524DHWVentilationScalarField[string]{true, "ab60_revision"},
			B524DHWVentilationScalarOpcode:      B524DHWVentilationScalarField[byte]{true, 2}, B524DHWVentilationScalarGroup: B524DHWVentilationScalarField[byte]{true, b524DHWVentilationScalarRow.B524DHWVentilationScalarGroup},
			B524DHWVentilationScalarInstance: B524DHWVentilationScalarField[byte]{true, 0}, B524DHWVentilationScalarRegister: B524DHWVentilationScalarField[uint16]{true, b524DHWVentilationScalarRow.B524DHWVentilationScalarRegister},
		},
		B524DHWVentilationScalarReplyKind:        B524DHWVentilationScalarField[string]{true, "echo_payload"},
		B524DHWVentilationScalarReturnedGroup:    B524DHWVentilationScalarField[byte]{true, b524DHWVentilationScalarRow.B524DHWVentilationScalarGroup},
		B524DHWVentilationScalarReturnedRegister: B524DHWVentilationScalarField[uint16]{true, b524DHWVentilationScalarRow.B524DHWVentilationScalarRegister},
		B524DHWVentilationScalarCodecGrade:       B524DHWVentilationScalarField[string]{true, "Hypothesis"}, B524DHWVentilationScalarRaw: b524DHWVentilationScalarRaw,
	}
}

func b524DHWVentilationScalarCheck(b524DHWVentilationScalarT *b524DHWVentilationScalarTesting.T, b524DHWVentilationScalarInput B524DHWVentilationScalarInput, b524DHWVentilationScalarCandidate bool) B524DHWVentilationScalarResult {
	b524DHWVentilationScalarT.Helper()
	b524DHWVentilationScalarResult := B524DHWVentilationInterpretScalars(b524DHWVentilationScalarInput)
	if !b524DHWVentilationScalarReflect.DeepEqual(b524DHWVentilationScalarResult.B524DHWVentilationScalarInput, b524DHWVentilationScalarInput) {
		b524DHWVentilationScalarT.Fatalf("ab60 caller evidence lost: got %+v want %+v", b524DHWVentilationScalarResult.B524DHWVentilationScalarInput, b524DHWVentilationScalarInput)
	}
	if b524DHWVentilationScalarResult.B524DHWVentilationScalarPresence != "Unknown" || b524DHWVentilationScalarResult.B524DHWVentilationScalarFreshness != "Unknown" || b524DHWVentilationScalarResult.B524DHWVentilationScalarExecution != "Unknown" || b524DHWVentilationScalarResult.B524DHWVentilationScalarSupport != "Unknown" || b524DHWVentilationScalarResult.B524DHWVentilationScalarEffects != "Unknown" || b524DHWVentilationScalarResult.B524DHWVentilationScalarPersistence != "Unknown" || b524DHWVentilationScalarResult.B524DHWVentilationScalarDelivery != "Unknown" {
		b524DHWVentilationScalarT.Fatal("ab60 invented physical qualification")
	}
	if b524DHWVentilationScalarCandidate {
		if b524DHWVentilationScalarResult.B524DHWVentilationScalarState != B524DHWVentilationScalarCandidate || b524DHWVentilationScalarResult.B524DHWVentilationScalarValue == nil || b524DHWVentilationScalarResult.B524DHWVentilationScalarGrade != "Hypothesis" {
			b524DHWVentilationScalarT.Fatalf("ab60 missing hypothesis: %+v", b524DHWVentilationScalarResult)
		}
	} else if b524DHWVentilationScalarResult.B524DHWVentilationScalarState != B524DHWVentilationScalarRawOnly || b524DHWVentilationScalarResult.B524DHWVentilationScalarValue != nil || b524DHWVentilationScalarResult.B524DHWVentilationScalarGrade != "Unknown" {
		b524DHWVentilationScalarT.Fatalf("ab60 raw-only admission failed: %+v", b524DHWVentilationScalarResult)
	}
	return b524DHWVentilationScalarResult
}

func TestAb60Bindings(b524DHWVentilationScalarT *b524DHWVentilationScalarTesting.T) {
	for _, b524DHWVentilationScalarRow := range b524DHWVentilationScalarFixtures(b524DHWVentilationScalarT) {
		b524DHWVentilationScalarT.Run(b524DHWVentilationScalarRow.B524DHWVentilationScalarLabel, func(b524DHWVentilationScalarT *b524DHWVentilationScalarTesting.T) {
			for _, b524DHWVentilationScalarExample := range b524DHWVentilationScalarRow.B524DHWVentilationScalarExamples {
				b524DHWVentilationScalarT.Run(b524DHWVentilationScalarExample.B524DHWVentilationScalarLabel, func(b524DHWVentilationScalarT *b524DHWVentilationScalarTesting.T) {
					b524DHWVentilationScalarResult := b524DHWVentilationScalarCheck(b524DHWVentilationScalarT, b524DHWVentilationScalarControl(b524DHWVentilationScalarT, b524DHWVentilationScalarRow, b524DHWVentilationScalarExample.B524DHWVentilationScalarHex), true)
					if b524DHWVentilationScalarResult.B524DHWVentilationScalarCodec != b524DHWVentilationScalarRow.B524DHWVentilationScalarCodec || b524DHWVentilationScalarResult.B524DHWVentilationScalarDeclaredUnit != b524DHWVentilationScalarRow.B524DHWVentilationScalarUnit {
						b524DHWVentilationScalarT.Fatal("ab60 metadata mismatch")
					}
					b524DHWVentilationScalarValue := float64(b524DHWVentilationScalarResult.B524DHWVentilationScalarValue.B524DHWVentilationScalarFloat32)
					if b524DHWVentilationScalarRow.B524DHWVentilationScalarCodec == "uint16_le" {
						b524DHWVentilationScalarValue = float64(b524DHWVentilationScalarResult.B524DHWVentilationScalarValue.B524DHWVentilationScalarUint16)
					}
					if b524DHWVentilationScalarValue != b524DHWVentilationScalarExample.B524DHWVentilationScalarExpected {
						b524DHWVentilationScalarT.Fatalf("ab60 value %v want %v", b524DHWVentilationScalarValue, b524DHWVentilationScalarExample.B524DHWVentilationScalarExpected)
					}
				})
			}
		})
	}
}

func TestAb60IndependentAxes(b524DHWVentilationScalarT *b524DHWVentilationScalarTesting.T) {
	for _, b524DHWVentilationScalarRow := range b524DHWVentilationScalarFixtures(b524DHWVentilationScalarT) {
		b524DHWVentilationScalarT.Run(b524DHWVentilationScalarRow.B524DHWVentilationScalarLabel, func(b524DHWVentilationScalarT *b524DHWVentilationScalarTesting.T) {
			for _, b524DHWVentilationScalarAxis := range []struct {
				b524DHWVentilationScalarName      string
				b524DHWVentilationScalarCandidate bool
				b524DHWVentilationScalarMutate    func(*B524DHWVentilationScalarInput)
			}{
				{"ab60_family_unknown", false, func(b524DHWVentilationScalarInput *B524DHWVentilationScalarInput) {
					b524DHWVentilationScalarInput.B524DHWVentilationScalarContext.B524DHWVentilationScalarFamily.B524DHWVentilationScalarValue = "ab60_unknown"
				}},
				{"ab60_family_other", b524DHWVentilationScalarRow.B524DHWVentilationScalarGroup == 1, func(b524DHWVentilationScalarInput *B524DHWVentilationScalarInput) {
					if b524DHWVentilationScalarRow.B524DHWVentilationScalarFamily == "controller-720" {
						b524DHWVentilationScalarInput.B524DHWVentilationScalarContext.B524DHWVentilationScalarFamily.B524DHWVentilationScalarValue = "controller-700"
					} else {
						b524DHWVentilationScalarInput.B524DHWVentilationScalarContext.B524DHWVentilationScalarFamily.B524DHWVentilationScalarValue = "controller-720"
					}
				}},
				{"ab60_software", true, func(b524DHWVentilationScalarInput *B524DHWVentilationScalarInput) {
					b524DHWVentilationScalarInput.B524DHWVentilationScalarContext.B524DHWVentilationScalarSoftware.B524DHWVentilationScalarValue = "ab60_changed"
				}},
				{"ab60_hardware", true, func(b524DHWVentilationScalarInput *B524DHWVentilationScalarInput) {
					b524DHWVentilationScalarInput.B524DHWVentilationScalarContext.B524DHWVentilationScalarHardware.B524DHWVentilationScalarValue = "ab60_changed"
				}},
				{"ab60_api_version", true, func(b524DHWVentilationScalarInput *B524DHWVentilationScalarInput) {
					b524DHWVentilationScalarInput.B524DHWVentilationScalarContext.B524DHWVentilationScalarAPIVersion.B524DHWVentilationScalarValue = "ab60_changed"
				}},
				{"ab60_api_revision", true, func(b524DHWVentilationScalarInput *B524DHWVentilationScalarInput) {
					b524DHWVentilationScalarInput.B524DHWVentilationScalarContext.B524DHWVentilationScalarAPIRevision.B524DHWVentilationScalarValue = "ab60_changed"
				}},
				{"ab60_opcode", false, func(b524DHWVentilationScalarInput *B524DHWVentilationScalarInput) {
					b524DHWVentilationScalarInput.B524DHWVentilationScalarContext.B524DHWVentilationScalarOpcode.B524DHWVentilationScalarValue = 6
				}},
				{"ab60_group", false, func(b524DHWVentilationScalarInput *B524DHWVentilationScalarInput) {
					b524DHWVentilationScalarInput.B524DHWVentilationScalarContext.B524DHWVentilationScalarGroup.B524DHWVentilationScalarValue = 0x7f
				}},
				{"ab60_instance", false, func(b524DHWVentilationScalarInput *B524DHWVentilationScalarInput) {
					b524DHWVentilationScalarInput.B524DHWVentilationScalarContext.B524DHWVentilationScalarInstance.B524DHWVentilationScalarValue = 1
				}},
				{"ab60_register", false, func(b524DHWVentilationScalarInput *B524DHWVentilationScalarInput) {
					b524DHWVentilationScalarInput.B524DHWVentilationScalarContext.B524DHWVentilationScalarRegister.B524DHWVentilationScalarValue = b524DHWVentilationScalarRow.B524DHWVentilationScalarRegister + 2
				}},
				{"ab60_reply_kind", false, func(b524DHWVentilationScalarInput *B524DHWVentilationScalarInput) {
					b524DHWVentilationScalarInput.B524DHWVentilationScalarReplyKind.B524DHWVentilationScalarValue = "status_only"
				}},
				{"ab60_returned_group", false, func(b524DHWVentilationScalarInput *B524DHWVentilationScalarInput) {
					b524DHWVentilationScalarInput.B524DHWVentilationScalarReturnedGroup.B524DHWVentilationScalarValue = 0x7f
				}},
				{"ab60_returned_register", false, func(b524DHWVentilationScalarInput *B524DHWVentilationScalarInput) {
					b524DHWVentilationScalarInput.B524DHWVentilationScalarReturnedRegister.B524DHWVentilationScalarValue++
				}},
				{"ab60_grade_proven", true, func(b524DHWVentilationScalarInput *B524DHWVentilationScalarInput) {
					b524DHWVentilationScalarInput.B524DHWVentilationScalarCodecGrade.B524DHWVentilationScalarValue = "Proven"
				}},
				{"ab60_grade_unknown", false, func(b524DHWVentilationScalarInput *B524DHWVentilationScalarInput) {
					b524DHWVentilationScalarInput.B524DHWVentilationScalarCodecGrade.B524DHWVentilationScalarValue = "Unknown"
				}},
			} {
				b524DHWVentilationScalarT.Run(b524DHWVentilationScalarAxis.b524DHWVentilationScalarName, func(b524DHWVentilationScalarT *b524DHWVentilationScalarTesting.T) {
					b524DHWVentilationScalarInput := b524DHWVentilationScalarControl(b524DHWVentilationScalarT, b524DHWVentilationScalarRow, b524DHWVentilationScalarRow.B524DHWVentilationScalarExamples[0].B524DHWVentilationScalarHex)
					b524DHWVentilationScalarAxis.b524DHWVentilationScalarMutate(&b524DHWVentilationScalarInput)
					b524DHWVentilationScalarResult := b524DHWVentilationScalarCheck(b524DHWVentilationScalarT, b524DHWVentilationScalarInput, b524DHWVentilationScalarAxis.b524DHWVentilationScalarCandidate)
					if b524DHWVentilationScalarAxis.b524DHWVentilationScalarCandidate {
						if b524DHWVentilationScalarResult.B524DHWVentilationScalarValue.B524DHWVentilationScalarFloat32 != 1.5 && b524DHWVentilationScalarRow.B524DHWVentilationScalarCodec == "float32_le" {
							b524DHWVentilationScalarT.Fatal("ab60 context changed arithmetic")
						}
						if b524DHWVentilationScalarResult.B524DHWVentilationScalarValue.B524DHWVentilationScalarUint16 != 4660 && b524DHWVentilationScalarRow.B524DHWVentilationScalarCodec == "uint16_le" {
							b524DHWVentilationScalarT.Fatal("ab60 context changed arithmetic")
						}
					}
				})
			}
			// Every context and reply field has an independent presence bit, including II00.
			b524DHWVentilationScalarInput := b524DHWVentilationScalarControl(b524DHWVentilationScalarT, b524DHWVentilationScalarRow, b524DHWVentilationScalarRow.B524DHWVentilationScalarExamples[0].B524DHWVentilationScalarHex)
			for _, b524DHWVentilationScalarContainer := range []string{"B524DHWVentilationScalarContext", "ab60_reply"} {
				b524DHWVentilationScalarFields := b524DHWVentilationScalarReflect.ValueOf(&b524DHWVentilationScalarInput).Elem()
				b524DHWVentilationScalarStart := 1
				b524DHWVentilationScalarEnd := 5
				if b524DHWVentilationScalarContainer == "B524DHWVentilationScalarContext" {
					b524DHWVentilationScalarFields = b524DHWVentilationScalarFields.FieldByName("B524DHWVentilationScalarContext")
					b524DHWVentilationScalarStart = 0
					b524DHWVentilationScalarEnd = b524DHWVentilationScalarFields.NumField()
				}
				for b524DHWVentilationScalarIndex := b524DHWVentilationScalarStart; b524DHWVentilationScalarIndex < b524DHWVentilationScalarEnd; b524DHWVentilationScalarIndex++ {
					b524DHWVentilationScalarT.Run("ab60_missing_"+b524DHWVentilationScalarFields.Type().Field(b524DHWVentilationScalarIndex).Name, func(b524DHWVentilationScalarT *b524DHWVentilationScalarTesting.T) {
						b524DHWVentilationScalarPresence := b524DHWVentilationScalarFields.Field(b524DHWVentilationScalarIndex).FieldByName("B524DHWVentilationScalarPresent")
						b524DHWVentilationScalarPresence.SetBool(false)
						b524DHWVentilationScalarCheck(b524DHWVentilationScalarT, b524DHWVentilationScalarInput, false)
						b524DHWVentilationScalarPresence.SetBool(true)
					})
				}
			}
		})
	}
}

func TestAb60MalformedAndOwnership(b524DHWVentilationScalarT *b524DHWVentilationScalarTesting.T) {
	for _, b524DHWVentilationScalarRow := range b524DHWVentilationScalarFixtures(b524DHWVentilationScalarT) {
		b524DHWVentilationScalarT.Run(b524DHWVentilationScalarRow.B524DHWVentilationScalarLabel, func(b524DHWVentilationScalarT *b524DHWVentilationScalarTesting.T) {
			for _, b524DHWVentilationScalarKind := range []string{"status_only", "short", "failed", "mismatched", "unknown", "", "error", "ack", "request"} {
				b524DHWVentilationScalarT.Run("ab60_kind_"+b524DHWVentilationScalarKind, func(b524DHWVentilationScalarT *b524DHWVentilationScalarTesting.T) {
					b524DHWVentilationScalarInput := b524DHWVentilationScalarControl(b524DHWVentilationScalarT, b524DHWVentilationScalarRow, b524DHWVentilationScalarRow.B524DHWVentilationScalarExamples[0].B524DHWVentilationScalarHex)
					b524DHWVentilationScalarInput.B524DHWVentilationScalarReplyKind.B524DHWVentilationScalarValue = b524DHWVentilationScalarKind
					b524DHWVentilationScalarCheck(b524DHWVentilationScalarT, b524DHWVentilationScalarInput, false)
				})
			}
			for _, b524DHWVentilationScalarGrade := range []string{"", "ab60_invalid", "hypothesis"} {
				b524DHWVentilationScalarT.Run("ab60_grade_"+b524DHWVentilationScalarGrade, func(b524DHWVentilationScalarT *b524DHWVentilationScalarTesting.T) {
					b524DHWVentilationScalarInput := b524DHWVentilationScalarControl(b524DHWVentilationScalarT, b524DHWVentilationScalarRow, b524DHWVentilationScalarRow.B524DHWVentilationScalarExamples[0].B524DHWVentilationScalarHex)
					b524DHWVentilationScalarInput.B524DHWVentilationScalarCodecGrade.B524DHWVentilationScalarValue = b524DHWVentilationScalarGrade
					b524DHWVentilationScalarCheck(b524DHWVentilationScalarT, b524DHWVentilationScalarInput, false)
				})
			}
			for _, b524DHWVentilationScalarVariant := range []string{"ab60_minus_one", "ab60_plus_one", "ab60_nil", "ab60_empty", "ab60_wrong_codec_width", "ab60_unlisted_group", "ab60_neighbor_low", "ab60_neighbor_high"} {
				b524DHWVentilationScalarT.Run(b524DHWVentilationScalarVariant, func(b524DHWVentilationScalarT *b524DHWVentilationScalarTesting.T) {
					b524DHWVentilationScalarInput := b524DHWVentilationScalarControl(b524DHWVentilationScalarT, b524DHWVentilationScalarRow, b524DHWVentilationScalarRow.B524DHWVentilationScalarExamples[0].B524DHWVentilationScalarHex)
					switch b524DHWVentilationScalarVariant {
					case "ab60_minus_one":
						b524DHWVentilationScalarInput.B524DHWVentilationScalarRaw = b524DHWVentilationScalarInput.B524DHWVentilationScalarRaw[:len(b524DHWVentilationScalarInput.B524DHWVentilationScalarRaw)-1]
					case "ab60_plus_one":
						b524DHWVentilationScalarInput.B524DHWVentilationScalarRaw = append(b524DHWVentilationScalarInput.B524DHWVentilationScalarRaw, 0)
					case "ab60_nil":
						b524DHWVentilationScalarInput.B524DHWVentilationScalarRaw = nil
					case "ab60_empty":
						b524DHWVentilationScalarInput.B524DHWVentilationScalarRaw = []byte{}
					case "ab60_wrong_codec_width":
						if b524DHWVentilationScalarRow.B524DHWVentilationScalarCodec == "float32_le" {
							b524DHWVentilationScalarInput.B524DHWVentilationScalarRaw = []byte{0x34, 0x12}
						} else {
							b524DHWVentilationScalarInput.B524DHWVentilationScalarRaw = []byte{0, 0, 0xc0, 0x3f}
						}
					case "ab60_unlisted_group":
						b524DHWVentilationScalarInput.B524DHWVentilationScalarContext.B524DHWVentilationScalarGroup.B524DHWVentilationScalarValue = 0x7f
						b524DHWVentilationScalarInput.B524DHWVentilationScalarReturnedGroup.B524DHWVentilationScalarValue = 0x7f
					case "ab60_neighbor_low":
						if b524DHWVentilationScalarRow.B524DHWVentilationScalarGroup == 1 {
							b524DHWVentilationScalarInput.B524DHWVentilationScalarContext.B524DHWVentilationScalarRegister.B524DHWVentilationScalarValue = 3
						} else {
							b524DHWVentilationScalarInput.B524DHWVentilationScalarContext.B524DHWVentilationScalarRegister.B524DHWVentilationScalarValue = 12
						}
						b524DHWVentilationScalarInput.B524DHWVentilationScalarReturnedRegister = b524DHWVentilationScalarInput.B524DHWVentilationScalarContext.B524DHWVentilationScalarRegister
					case "ab60_neighbor_high":
						if b524DHWVentilationScalarRow.B524DHWVentilationScalarGroup == 1 {
							b524DHWVentilationScalarInput.B524DHWVentilationScalarContext.B524DHWVentilationScalarRegister.B524DHWVentilationScalarValue = 6
						} else {
							b524DHWVentilationScalarInput.B524DHWVentilationScalarContext.B524DHWVentilationScalarRegister.B524DHWVentilationScalarValue = 15
						}
						b524DHWVentilationScalarInput.B524DHWVentilationScalarReturnedRegister = b524DHWVentilationScalarInput.B524DHWVentilationScalarContext.B524DHWVentilationScalarRegister
					}
					b524DHWVentilationScalarResult := b524DHWVentilationScalarCheck(b524DHWVentilationScalarT, b524DHWVentilationScalarInput, false)
					if b524DHWVentilationScalarVariant != "ab60_unlisted_group" && b524DHWVentilationScalarVariant != "ab60_neighbor_low" && b524DHWVentilationScalarVariant != "ab60_neighbor_high" && b524DHWVentilationScalarResult.B524DHWVentilationScalarDeclaredUnit != b524DHWVentilationScalarRow.B524DHWVentilationScalarUnit {
						b524DHWVentilationScalarT.Fatal("ab60 declared unit lost on partial failure")
					}
				})
			}
			if b524DHWVentilationScalarRow.B524DHWVentilationScalarCodec == "float32_le" {
				for _, b524DHWVentilationScalarBits := range []string{"0000807f", "000080ff", "0000c07f"} {
					b524DHWVentilationScalarT.Run("ab60_nonfinite_"+b524DHWVentilationScalarBits, func(b524DHWVentilationScalarT *b524DHWVentilationScalarTesting.T) {
						b524DHWVentilationScalarCheck(b524DHWVentilationScalarT, b524DHWVentilationScalarControl(b524DHWVentilationScalarT, b524DHWVentilationScalarRow, b524DHWVentilationScalarBits), false)
					})
				}
			}
			for _, b524DHWVentilationScalarCandidate := range []bool{true, false} {
				b524DHWVentilationScalarName := "ab60_memory_candidate"
				if !b524DHWVentilationScalarCandidate {
					b524DHWVentilationScalarName = "ab60_memory_raw_only"
				}
				b524DHWVentilationScalarT.Run(b524DHWVentilationScalarName, func(b524DHWVentilationScalarT *b524DHWVentilationScalarTesting.T) {
					b524DHWVentilationScalarInput := b524DHWVentilationScalarControl(b524DHWVentilationScalarT, b524DHWVentilationScalarRow, b524DHWVentilationScalarRow.B524DHWVentilationScalarExamples[0].B524DHWVentilationScalarHex)
					if !b524DHWVentilationScalarCandidate {
						b524DHWVentilationScalarInput.B524DHWVentilationScalarReplyKind.B524DHWVentilationScalarValue = "failed"
					}
					b524DHWVentilationScalarFirst := b524DHWVentilationScalarCheck(b524DHWVentilationScalarT, b524DHWVentilationScalarInput, b524DHWVentilationScalarCandidate)
					b524DHWVentilationScalarSecond := b524DHWVentilationScalarCheck(b524DHWVentilationScalarT, b524DHWVentilationScalarInput, b524DHWVentilationScalarCandidate)
					b524DHWVentilationScalarOriginal := b524DHWVentilationScalarInput.B524DHWVentilationScalarRaw[0]
					b524DHWVentilationScalarInput.B524DHWVentilationScalarRaw[0] ^= 0xff
					b524DHWVentilationScalarInput.B524DHWVentilationScalarContext.B524DHWVentilationScalarSoftware.B524DHWVentilationScalarValue = "ab60_mutated"
					if b524DHWVentilationScalarFirst.B524DHWVentilationScalarInput.B524DHWVentilationScalarRaw[0] != b524DHWVentilationScalarOriginal || b524DHWVentilationScalarFirst.B524DHWVentilationScalarInput.B524DHWVentilationScalarContext.B524DHWVentilationScalarSoftware.B524DHWVentilationScalarValue != "ab60_software" {
						b524DHWVentilationScalarT.Fatal("ab60 input aliases result")
					}
					b524DHWVentilationScalarFirst.B524DHWVentilationScalarInput.B524DHWVentilationScalarRaw[0] ^= 0x55
					if b524DHWVentilationScalarInput.B524DHWVentilationScalarRaw[0] != b524DHWVentilationScalarOriginal^0xff || b524DHWVentilationScalarSecond.B524DHWVentilationScalarInput.B524DHWVentilationScalarRaw[0] != b524DHWVentilationScalarOriginal {
						b524DHWVentilationScalarT.Fatal("ab60 returned memory aliases caller or another result")
					}
					if b524DHWVentilationScalarCandidate {
						b524DHWVentilationScalarFirst.B524DHWVentilationScalarValue.B524DHWVentilationScalarFloat32 = 99
						b524DHWVentilationScalarFirst.B524DHWVentilationScalarValue.B524DHWVentilationScalarUint16 = 99
						if b524DHWVentilationScalarSecond.B524DHWVentilationScalarValue.B524DHWVentilationScalarFloat32 == 99 || b524DHWVentilationScalarSecond.B524DHWVentilationScalarValue.B524DHWVentilationScalarUint16 == 99 {
							b524DHWVentilationScalarT.Fatal("ab60 scalar aliases another result")
						}
					}
				})
			}
		})
	}
}

func TestAb60ExplicitZeroContext(b524DHWVentilationScalarT *b524DHWVentilationScalarTesting.T) {
	for _, b524DHWVentilationScalarRow := range b524DHWVentilationScalarFixtures(b524DHWVentilationScalarT) {
		b524DHWVentilationScalarT.Run(b524DHWVentilationScalarRow.B524DHWVentilationScalarLabel, func(b524DHWVentilationScalarT *b524DHWVentilationScalarTesting.T) {
			for _, b524DHWVentilationScalarField := range []string{"B524DHWVentilationScalarSoftware", "B524DHWVentilationScalarHardware", "B524DHWVentilationScalarAPIVersion", "B524DHWVentilationScalarAPIRevision"} {
				b524DHWVentilationScalarT.Run("ab60_zero_"+b524DHWVentilationScalarField, func(b524DHWVentilationScalarT *b524DHWVentilationScalarTesting.T) {
					b524DHWVentilationScalarInput := b524DHWVentilationScalarControl(b524DHWVentilationScalarT, b524DHWVentilationScalarRow, b524DHWVentilationScalarRow.B524DHWVentilationScalarExamples[0].B524DHWVentilationScalarHex)
					b524DHWVentilationScalarReflect.ValueOf(&b524DHWVentilationScalarInput.B524DHWVentilationScalarContext).Elem().FieldByName(b524DHWVentilationScalarField).FieldByName("B524DHWVentilationScalarValue").SetString("")
					b524DHWVentilationScalarCheck(b524DHWVentilationScalarT, b524DHWVentilationScalarInput, true)
				})
			}
		})
	}
}
