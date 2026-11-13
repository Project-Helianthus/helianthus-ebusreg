package system

import (
	b524SystemScalarHex "encoding/hex"
	b524SystemScalarJSON "encoding/json"
	b524SystemScalarMath "math"
	b524SystemScalarOs "os"
	b524SystemScalarReflect "reflect"
	b524SystemScalarTesting "testing"
)

type b524SystemScalarFixture struct {
	B524SystemScalarLabel    string `json:"ab58_label"`
	B524SystemScalarFamily   string `json:"ab58_family"`
	B524SystemScalarGroup    byte   `json:"ab58_group"`
	B524SystemScalarRegister uint16 `json:"ab58_register"`
	B524SystemScalarCodec    string `json:"ab58_codec"`
	B524SystemScalarUnit     string `json:"ab58_unit"`
	B524SystemScalarExamples []struct {
		B524SystemScalarLabel    string  `json:"ab58_label"`
		B524SystemScalarHex      string  `json:"ab58_hex"`
		B524SystemScalarExpected float64 `json:"ab58_expected"`
	} `json:"ab58_examples"`
}

func b524SystemScalarFixtures(b524SystemScalarT *b524SystemScalarTesting.T) []b524SystemScalarFixture {
	b524SystemScalarT.Helper()
	b524SystemScalarBytes, b524SystemScalarErr := b524SystemScalarOs.ReadFile("testdata/ab58_b524_system_scalars.json")
	if b524SystemScalarErr != nil {
		b524SystemScalarT.Fatal(b524SystemScalarErr)
	}
	var b524SystemScalarRows []b524SystemScalarFixture
	if b524SystemScalarErr = b524SystemScalarJSON.Unmarshal(b524SystemScalarBytes, &b524SystemScalarRows); b524SystemScalarErr != nil {
		b524SystemScalarT.Fatal(b524SystemScalarErr)
	}
	b524SystemScalarExpected := []b524SystemScalarFixture{
		{B524SystemScalarFamily: "controller-720", B524SystemScalarGroup: 0, B524SystemScalarRegister: 0x0002, B524SystemScalarCodec: "float32_le", B524SystemScalarUnit: "°C"},
		{B524SystemScalarFamily: "controller-720", B524SystemScalarGroup: 0, B524SystemScalarRegister: 0x0022, B524SystemScalarCodec: "float32_le", B524SystemScalarUnit: "°C"},
		{B524SystemScalarFamily: "controller-720", B524SystemScalarGroup: 0, B524SystemScalarRegister: 0x0023, B524SystemScalarCodec: "float32_le", B524SystemScalarUnit: "°C"},
		{B524SystemScalarFamily: "controller-720", B524SystemScalarGroup: 0, B524SystemScalarRegister: 0x0026, B524SystemScalarCodec: "float32_le", B524SystemScalarUnit: "°C"},
		{B524SystemScalarFamily: "controller-720", B524SystemScalarGroup: 0, B524SystemScalarRegister: 0x0038, B524SystemScalarCodec: "float32_le", B524SystemScalarUnit: "°C"},
		{B524SystemScalarFamily: "controller-720", B524SystemScalarGroup: 0, B524SystemScalarRegister: 0x0046, B524SystemScalarCodec: "float32_le", B524SystemScalarUnit: "°C"},
		{B524SystemScalarFamily: "controller-720", B524SystemScalarGroup: 0, B524SystemScalarRegister: 0x004B, B524SystemScalarCodec: "float32_le", B524SystemScalarUnit: "°C"},
		{B524SystemScalarFamily: "controller-720", B524SystemScalarGroup: 0, B524SystemScalarRegister: 0x0067, B524SystemScalarCodec: "float32_le", B524SystemScalarUnit: "°C"},
		{B524SystemScalarFamily: "controller-720", B524SystemScalarGroup: 0, B524SystemScalarRegister: 0x0068, B524SystemScalarCodec: "float32_le", B524SystemScalarUnit: "°C"},
		{B524SystemScalarFamily: "controller-720", B524SystemScalarGroup: 0, B524SystemScalarRegister: 0x0073, B524SystemScalarCodec: "float32_le", B524SystemScalarUnit: "°C"},
		{B524SystemScalarFamily: "controller-720", B524SystemScalarGroup: 0, B524SystemScalarRegister: 0x0095, B524SystemScalarCodec: "float32_le", B524SystemScalarUnit: "°C"},
		{B524SystemScalarFamily: "controller-720", B524SystemScalarGroup: 0, B524SystemScalarRegister: 0x009D, B524SystemScalarCodec: "float32_le", B524SystemScalarUnit: "°C"},
		{B524SystemScalarFamily: "controller-720", B524SystemScalarGroup: 0, B524SystemScalarRegister: 0x009E, B524SystemScalarCodec: "float32_le", B524SystemScalarUnit: "°C"},
		{B524SystemScalarFamily: "controller-720", B524SystemScalarGroup: 0, B524SystemScalarRegister: 0x009F, B524SystemScalarCodec: "float32_le", B524SystemScalarUnit: "°C"},
		{B524SystemScalarFamily: "controller-720", B524SystemScalarGroup: 0, B524SystemScalarRegister: 0x00A0, B524SystemScalarCodec: "float32_le", B524SystemScalarUnit: "°C"},
		{B524SystemScalarFamily: "controller-720", B524SystemScalarGroup: 0, B524SystemScalarRegister: 0x00EF, B524SystemScalarCodec: "float32_le", B524SystemScalarUnit: "°C"},
		{B524SystemScalarFamily: "controller-720", B524SystemScalarGroup: 0, B524SystemScalarRegister: 0x0004, B524SystemScalarCodec: "uint16_le", B524SystemScalarUnit: "min"},
		{B524SystemScalarFamily: "controller-720", B524SystemScalarGroup: 0, B524SystemScalarRegister: 0x0017, B524SystemScalarCodec: "uint16_le", B524SystemScalarUnit: "min"},
		{B524SystemScalarFamily: "controller-720", B524SystemScalarGroup: 0, B524SystemScalarRegister: 0x0018, B524SystemScalarCodec: "uint16_le", B524SystemScalarUnit: "min"},
		{B524SystemScalarFamily: "controller-720", B524SystemScalarGroup: 0, B524SystemScalarRegister: 0x0029, B524SystemScalarCodec: "float32_le", B524SystemScalarUnit: "K"},
		{B524SystemScalarFamily: "controller-720", B524SystemScalarGroup: 0, B524SystemScalarRegister: 0x002F, B524SystemScalarCodec: "uint16_le", B524SystemScalarUnit: "Unknown"},
		{B524SystemScalarFamily: "controller-720", B524SystemScalarGroup: 0, B524SystemScalarRegister: 0x0036, B524SystemScalarCodec: "uint16_le", B524SystemScalarUnit: "Unknown"},
		{B524SystemScalarFamily: "controller-720", B524SystemScalarGroup: 0, B524SystemScalarRegister: 0x0039, B524SystemScalarCodec: "float32_le", B524SystemScalarUnit: "bar"},
		{B524SystemScalarFamily: "controller-700", B524SystemScalarGroup: 0, B524SystemScalarRegister: 0x0002, B524SystemScalarCodec: "float32_le", B524SystemScalarUnit: "°C"},
		{B524SystemScalarFamily: "controller-700", B524SystemScalarGroup: 0, B524SystemScalarRegister: 0x0022, B524SystemScalarCodec: "float32_le", B524SystemScalarUnit: "°C"},
		{B524SystemScalarFamily: "controller-700", B524SystemScalarGroup: 0, B524SystemScalarRegister: 0x0023, B524SystemScalarCodec: "float32_le", B524SystemScalarUnit: "°C"},
		{B524SystemScalarFamily: "controller-700", B524SystemScalarGroup: 0, B524SystemScalarRegister: 0x0026, B524SystemScalarCodec: "float32_le", B524SystemScalarUnit: "°C"},
		{B524SystemScalarFamily: "controller-700", B524SystemScalarGroup: 0, B524SystemScalarRegister: 0x0038, B524SystemScalarCodec: "float32_le", B524SystemScalarUnit: "°C"},
		{B524SystemScalarFamily: "controller-700", B524SystemScalarGroup: 0, B524SystemScalarRegister: 0x0046, B524SystemScalarCodec: "float32_le", B524SystemScalarUnit: "°C"},
		{B524SystemScalarFamily: "controller-700", B524SystemScalarGroup: 0, B524SystemScalarRegister: 0x004B, B524SystemScalarCodec: "float32_le", B524SystemScalarUnit: "°C"},
		{B524SystemScalarFamily: "controller-700", B524SystemScalarGroup: 0, B524SystemScalarRegister: 0x0067, B524SystemScalarCodec: "float32_le", B524SystemScalarUnit: "°C"},
		{B524SystemScalarFamily: "controller-700", B524SystemScalarGroup: 0, B524SystemScalarRegister: 0x0068, B524SystemScalarCodec: "float32_le", B524SystemScalarUnit: "°C"},
		{B524SystemScalarFamily: "controller-700", B524SystemScalarGroup: 0, B524SystemScalarRegister: 0x0073, B524SystemScalarCodec: "float32_le", B524SystemScalarUnit: "°C"},
		{B524SystemScalarFamily: "controller-700", B524SystemScalarGroup: 0, B524SystemScalarRegister: 0x0095, B524SystemScalarCodec: "float32_le", B524SystemScalarUnit: "°C"},
		{B524SystemScalarFamily: "controller-700", B524SystemScalarGroup: 0, B524SystemScalarRegister: 0x009D, B524SystemScalarCodec: "float32_le", B524SystemScalarUnit: "°C"},
		{B524SystemScalarFamily: "controller-700", B524SystemScalarGroup: 0, B524SystemScalarRegister: 0x009E, B524SystemScalarCodec: "float32_le", B524SystemScalarUnit: "°C"},
		{B524SystemScalarFamily: "controller-700", B524SystemScalarGroup: 0, B524SystemScalarRegister: 0x009F, B524SystemScalarCodec: "float32_le", B524SystemScalarUnit: "°C"},
		{B524SystemScalarFamily: "controller-700", B524SystemScalarGroup: 0, B524SystemScalarRegister: 0x00A0, B524SystemScalarCodec: "float32_le", B524SystemScalarUnit: "°C"},
		{B524SystemScalarFamily: "controller-700", B524SystemScalarGroup: 0, B524SystemScalarRegister: 0x0006, B524SystemScalarCodec: "uint16_le", B524SystemScalarUnit: "Unknown"},
		{B524SystemScalarFamily: "controller-700", B524SystemScalarGroup: 0, B524SystemScalarRegister: 0x002F, B524SystemScalarCodec: "uint16_le", B524SystemScalarUnit: "Unknown"},
		{B524SystemScalarFamily: "controller-700", B524SystemScalarGroup: 0, B524SystemScalarRegister: 0x0036, B524SystemScalarCodec: "uint16_le", B524SystemScalarUnit: "Unknown"},
		{B524SystemScalarFamily: "controller-700", B524SystemScalarGroup: 0, B524SystemScalarRegister: 0x0017, B524SystemScalarCodec: "uint16_le", B524SystemScalarUnit: "min"},
		{B524SystemScalarFamily: "controller-700", B524SystemScalarGroup: 0, B524SystemScalarRegister: 0x0029, B524SystemScalarCodec: "float32_le", B524SystemScalarUnit: "K"},
		{B524SystemScalarFamily: "controller-700", B524SystemScalarGroup: 0, B524SystemScalarRegister: 0x0039, B524SystemScalarCodec: "float32_le", B524SystemScalarUnit: "bar"},
		{B524SystemScalarFamily: "controller-700", B524SystemScalarGroup: 0, B524SystemScalarRegister: 0x0042, B524SystemScalarCodec: "uint16_le", B524SystemScalarUnit: "Unknown"},
		{B524SystemScalarFamily: "controller-700", B524SystemScalarGroup: 0, B524SystemScalarRegister: 0x0043, B524SystemScalarCodec: "uint16_le", B524SystemScalarUnit: "Unknown"},
	}
	if len(b524SystemScalarRows) != len(b524SystemScalarExpected) {
		b524SystemScalarT.Fatal("ab58 incomplete fixture set")
	}
	for b524SystemScalarIndex, b524SystemScalarRow := range b524SystemScalarRows {
		b524SystemScalarWant := b524SystemScalarExpected[b524SystemScalarIndex]
		if b524SystemScalarRow.B524SystemScalarFamily != b524SystemScalarWant.B524SystemScalarFamily || b524SystemScalarRow.B524SystemScalarGroup != b524SystemScalarWant.B524SystemScalarGroup || b524SystemScalarRow.B524SystemScalarRegister != b524SystemScalarWant.B524SystemScalarRegister || b524SystemScalarRow.B524SystemScalarCodec != b524SystemScalarWant.B524SystemScalarCodec || b524SystemScalarRow.B524SystemScalarUnit != b524SystemScalarWant.B524SystemScalarUnit || len(b524SystemScalarRow.B524SystemScalarExamples) != 2 {
			b524SystemScalarT.Fatal("ab58 fixture selector or metadata drift")
		}
		b524SystemScalarHexes := []string{"0000c03f", "000080bf"}
		b524SystemScalarValues := []float64{1.5, -1}
		if b524SystemScalarWant.B524SystemScalarCodec == "uint16_le" {
			b524SystemScalarHexes = []string{"3412", "ffff"}
			b524SystemScalarValues = []float64{4660, 65535}
		}
		for b524SystemScalarIndex, b524SystemScalarExample := range b524SystemScalarRow.B524SystemScalarExamples {
			if b524SystemScalarExample.B524SystemScalarHex != b524SystemScalarHexes[b524SystemScalarIndex] || b524SystemScalarExample.B524SystemScalarExpected != b524SystemScalarValues[b524SystemScalarIndex] {
				b524SystemScalarT.Fatal("ab58 independent arithmetic control drift")
			}
		}
	}
	return b524SystemScalarRows
}

func b524SystemScalarControl(b524SystemScalarT *b524SystemScalarTesting.T, b524SystemScalarRow b524SystemScalarFixture, b524SystemScalarHexadecimal string) B524SystemScalarInput {
	b524SystemScalarT.Helper()
	b524SystemScalarRaw, b524SystemScalarErr := b524SystemScalarHex.DecodeString(b524SystemScalarHexadecimal)
	if b524SystemScalarErr != nil {
		b524SystemScalarT.Fatal(b524SystemScalarErr)
	}
	return B524SystemScalarInput{
		B524SystemScalarContext: B524SystemScalarContext{
			B524SystemScalarFamily:      B524SystemScalarField[string]{true, b524SystemScalarRow.B524SystemScalarFamily},
			B524SystemScalarSoftware:    B524SystemScalarField[string]{true, "ab58_software"},
			B524SystemScalarHardware:    B524SystemScalarField[string]{true, "ab58_hardware"},
			B524SystemScalarAPIVersion:  B524SystemScalarField[string]{true, "ab58_api"},
			B524SystemScalarAPIRevision: B524SystemScalarField[string]{true, "ab58_revision"},
			B524SystemScalarOpcode:      B524SystemScalarField[byte]{true, 2}, B524SystemScalarGroup: B524SystemScalarField[byte]{true, b524SystemScalarRow.B524SystemScalarGroup},
			B524SystemScalarInstance: B524SystemScalarField[byte]{true, 0}, B524SystemScalarRegister: B524SystemScalarField[uint16]{true, b524SystemScalarRow.B524SystemScalarRegister},
		},
		B524SystemScalarReplyKind:        B524SystemScalarField[string]{true, "echo_payload"},
		B524SystemScalarReturnedGroup:    B524SystemScalarField[byte]{true, b524SystemScalarRow.B524SystemScalarGroup},
		B524SystemScalarReturnedRegister: B524SystemScalarField[uint16]{true, b524SystemScalarRow.B524SystemScalarRegister},
		B524SystemScalarCodecGrade:       B524SystemScalarField[string]{true, "Hypothesis"}, B524SystemScalarRaw: b524SystemScalarRaw,
	}
}

func b524SystemScalarCheck(b524SystemScalarT *b524SystemScalarTesting.T, b524SystemScalarInput B524SystemScalarInput, b524SystemScalarCandidate bool) B524SystemScalarResult {
	b524SystemScalarT.Helper()
	b524SystemScalarResult := B524SystemInterpretScalars(b524SystemScalarInput)
	if !b524SystemScalarReflect.DeepEqual(b524SystemScalarResult.B524SystemScalarInput, b524SystemScalarInput) {
		b524SystemScalarT.Fatalf("ab58 caller evidence lost: got %+v want %+v", b524SystemScalarResult.B524SystemScalarInput, b524SystemScalarInput)
	}
	if b524SystemScalarResult.B524SystemScalarPresence != "Unknown" || b524SystemScalarResult.B524SystemScalarFreshness != "Unknown" || b524SystemScalarResult.B524SystemScalarExecution != "Unknown" || b524SystemScalarResult.B524SystemScalarSupport != "Unknown" || b524SystemScalarResult.B524SystemScalarEffects != "Unknown" || b524SystemScalarResult.B524SystemScalarPersistence != "Unknown" || b524SystemScalarResult.B524SystemScalarDelivery != "Unknown" {
		b524SystemScalarT.Fatal("ab58 invented physical qualification")
	}
	if b524SystemScalarCandidate {
		if b524SystemScalarResult.B524SystemScalarState != B524SystemScalarCandidate || b524SystemScalarResult.B524SystemScalarValue == nil || b524SystemScalarResult.B524SystemScalarGrade != "Hypothesis" {
			b524SystemScalarT.Fatalf("ab58 missing hypothesis: %+v", b524SystemScalarResult)
		}
	} else if b524SystemScalarResult.B524SystemScalarState != B524SystemScalarRawOnly || b524SystemScalarResult.B524SystemScalarValue != nil || b524SystemScalarResult.B524SystemScalarGrade != "Unknown" {
		b524SystemScalarT.Fatalf("ab58 raw-only admission failed: %+v", b524SystemScalarResult)
	}
	return b524SystemScalarResult
}

func TestAb58Bindings(b524SystemScalarT *b524SystemScalarTesting.T) {
	for _, b524SystemScalarRow := range b524SystemScalarFixtures(b524SystemScalarT) {
		b524SystemScalarT.Run(b524SystemScalarRow.B524SystemScalarLabel, func(b524SystemScalarT *b524SystemScalarTesting.T) {
			for _, b524SystemScalarExample := range b524SystemScalarRow.B524SystemScalarExamples {
				b524SystemScalarT.Run(b524SystemScalarExample.B524SystemScalarLabel, func(b524SystemScalarT *b524SystemScalarTesting.T) {
					b524SystemScalarResult := b524SystemScalarCheck(b524SystemScalarT, b524SystemScalarControl(b524SystemScalarT, b524SystemScalarRow, b524SystemScalarExample.B524SystemScalarHex), true)
					if b524SystemScalarResult.B524SystemScalarCodec != b524SystemScalarRow.B524SystemScalarCodec || b524SystemScalarResult.B524SystemScalarDeclaredUnit != b524SystemScalarRow.B524SystemScalarUnit {
						b524SystemScalarT.Fatal("ab58 metadata mismatch")
					}
					b524SystemScalarValue := float64(b524SystemScalarResult.B524SystemScalarValue.B524SystemScalarFloat32)
					if b524SystemScalarRow.B524SystemScalarCodec == "uint16_le" {
						b524SystemScalarValue = float64(b524SystemScalarResult.B524SystemScalarValue.B524SystemScalarUint16)
					}
					if b524SystemScalarValue != b524SystemScalarExample.B524SystemScalarExpected {
						b524SystemScalarT.Fatalf("ab58 value %v want %v", b524SystemScalarValue, b524SystemScalarExample.B524SystemScalarExpected)
					}
				})
			}
		})
	}
}

func TestAb58IndependentAxes(b524SystemScalarT *b524SystemScalarTesting.T) {
	for _, b524SystemScalarRow := range b524SystemScalarFixtures(b524SystemScalarT) {
		b524SystemScalarT.Run(b524SystemScalarRow.B524SystemScalarLabel, func(b524SystemScalarT *b524SystemScalarTesting.T) {
			for _, b524SystemScalarAxis := range []struct {
				b524SystemScalarName      string
				b524SystemScalarCandidate bool
				b524SystemScalarMutate    func(*B524SystemScalarInput)
			}{
				{"ab58_family_unknown", false, func(b524SystemScalarInput *B524SystemScalarInput) {
					b524SystemScalarInput.B524SystemScalarContext.B524SystemScalarFamily.B524SystemScalarValue = "ab58_unknown"
				}},
				{"ab58_family_other", (b524SystemScalarRow.B524SystemScalarRegister != 4 && b524SystemScalarRow.B524SystemScalarRegister != 0x18 && b524SystemScalarRow.B524SystemScalarRegister != 0xef && b524SystemScalarRow.B524SystemScalarRegister != 6 && b524SystemScalarRow.B524SystemScalarRegister != 0x42 && b524SystemScalarRow.B524SystemScalarRegister != 0x43), func(b524SystemScalarInput *B524SystemScalarInput) {
					if b524SystemScalarRow.B524SystemScalarFamily == "controller-720" {
						b524SystemScalarInput.B524SystemScalarContext.B524SystemScalarFamily.B524SystemScalarValue = "controller-700"
					} else {
						b524SystemScalarInput.B524SystemScalarContext.B524SystemScalarFamily.B524SystemScalarValue = "controller-720"
					}
				}},
				{"ab58_software", true, func(b524SystemScalarInput *B524SystemScalarInput) {
					b524SystemScalarInput.B524SystemScalarContext.B524SystemScalarSoftware.B524SystemScalarValue = "ab58_changed"
				}},
				{"ab58_hardware", true, func(b524SystemScalarInput *B524SystemScalarInput) {
					b524SystemScalarInput.B524SystemScalarContext.B524SystemScalarHardware.B524SystemScalarValue = "ab58_changed"
				}},
				{"ab58_api_version", true, func(b524SystemScalarInput *B524SystemScalarInput) {
					b524SystemScalarInput.B524SystemScalarContext.B524SystemScalarAPIVersion.B524SystemScalarValue = "ab58_changed"
				}},
				{"ab58_api_revision", true, func(b524SystemScalarInput *B524SystemScalarInput) {
					b524SystemScalarInput.B524SystemScalarContext.B524SystemScalarAPIRevision.B524SystemScalarValue = "ab58_changed"
				}},
				{"ab58_opcode", false, func(b524SystemScalarInput *B524SystemScalarInput) {
					b524SystemScalarInput.B524SystemScalarContext.B524SystemScalarOpcode.B524SystemScalarValue = 6
				}},
				{"ab58_group", false, func(b524SystemScalarInput *B524SystemScalarInput) {
					b524SystemScalarInput.B524SystemScalarContext.B524SystemScalarGroup.B524SystemScalarValue = 0x7f
				}},
				{"ab58_instance", b524SystemScalarRow.B524SystemScalarFamily == "controller-700" && (b524SystemScalarRow.B524SystemScalarRegister == 0x42 || b524SystemScalarRow.B524SystemScalarRegister == 0x43), func(b524SystemScalarInput *B524SystemScalarInput) {
					b524SystemScalarInput.B524SystemScalarContext.B524SystemScalarInstance.B524SystemScalarValue = 1
				}},
				{"ab58_register", false, func(b524SystemScalarInput *B524SystemScalarInput) {
					b524SystemScalarInput.B524SystemScalarContext.B524SystemScalarRegister.B524SystemScalarValue = b524SystemScalarRow.B524SystemScalarRegister + 2
				}},
				{"ab58_reply_kind", false, func(b524SystemScalarInput *B524SystemScalarInput) {
					b524SystemScalarInput.B524SystemScalarReplyKind.B524SystemScalarValue = "status_only"
				}},
				{"ab58_returned_group", false, func(b524SystemScalarInput *B524SystemScalarInput) {
					b524SystemScalarInput.B524SystemScalarReturnedGroup.B524SystemScalarValue = 0x7f
				}},
				{"ab58_returned_register", false, func(b524SystemScalarInput *B524SystemScalarInput) {
					b524SystemScalarInput.B524SystemScalarReturnedRegister.B524SystemScalarValue++
				}},
				{"ab58_grade_proven", true, func(b524SystemScalarInput *B524SystemScalarInput) {
					b524SystemScalarInput.B524SystemScalarCodecGrade.B524SystemScalarValue = "Proven"
				}},
				{"ab58_grade_unknown", false, func(b524SystemScalarInput *B524SystemScalarInput) {
					b524SystemScalarInput.B524SystemScalarCodecGrade.B524SystemScalarValue = "Unknown"
				}},
			} {
				b524SystemScalarT.Run(b524SystemScalarAxis.b524SystemScalarName, func(b524SystemScalarT *b524SystemScalarTesting.T) {
					b524SystemScalarInput := b524SystemScalarControl(b524SystemScalarT, b524SystemScalarRow, b524SystemScalarRow.B524SystemScalarExamples[0].B524SystemScalarHex)
					b524SystemScalarAxis.b524SystemScalarMutate(&b524SystemScalarInput)
					b524SystemScalarResult := b524SystemScalarCheck(b524SystemScalarT, b524SystemScalarInput, b524SystemScalarAxis.b524SystemScalarCandidate)
					if b524SystemScalarAxis.b524SystemScalarCandidate {
						if b524SystemScalarResult.B524SystemScalarValue.B524SystemScalarFloat32 != 1.5 && b524SystemScalarRow.B524SystemScalarCodec == "float32_le" {
							b524SystemScalarT.Fatal("ab58 context changed arithmetic")
						}
						if b524SystemScalarResult.B524SystemScalarValue.B524SystemScalarUint16 != 4660 && b524SystemScalarRow.B524SystemScalarCodec == "uint16_le" {
							b524SystemScalarT.Fatal("ab58 context changed arithmetic")
						}
					}
				})
			}
			b524SystemScalarInput := b524SystemScalarControl(b524SystemScalarT, b524SystemScalarRow, b524SystemScalarRow.B524SystemScalarExamples[0].B524SystemScalarHex)
			for _, b524SystemScalarContainer := range []string{"B524SystemScalarContext", "ab58_reply"} {
				b524SystemScalarFields := b524SystemScalarReflect.ValueOf(&b524SystemScalarInput).Elem()
				b524SystemScalarStart := 1
				b524SystemScalarEnd := 5
				if b524SystemScalarContainer == "B524SystemScalarContext" {
					b524SystemScalarFields = b524SystemScalarFields.FieldByName("B524SystemScalarContext")
					b524SystemScalarStart = 0
					b524SystemScalarEnd = b524SystemScalarFields.NumField()
				}
				for b524SystemScalarIndex := b524SystemScalarStart; b524SystemScalarIndex < b524SystemScalarEnd; b524SystemScalarIndex++ {
					b524SystemScalarT.Run("ab58_missing_"+b524SystemScalarFields.Type().Field(b524SystemScalarIndex).Name, func(b524SystemScalarT *b524SystemScalarTesting.T) {
						b524SystemScalarPresence := b524SystemScalarFields.Field(b524SystemScalarIndex).FieldByName("B524SystemScalarPresent")
						b524SystemScalarPresence.SetBool(false)
						b524SystemScalarCheck(b524SystemScalarT, b524SystemScalarInput, false)
						b524SystemScalarPresence.SetBool(true)
					})
				}
			}
		})
	}
}

func TestAb58MalformedAndOwnership(b524SystemScalarT *b524SystemScalarTesting.T) {
	for _, b524SystemScalarRow := range b524SystemScalarFixtures(b524SystemScalarT) {
		b524SystemScalarT.Run(b524SystemScalarRow.B524SystemScalarLabel, func(b524SystemScalarT *b524SystemScalarTesting.T) {
			for _, b524SystemScalarKind := range []string{"status_only", "short", "failed", "mismatched", "unknown", "", "error", "ack", "request"} {
				b524SystemScalarT.Run("ab58_kind_"+b524SystemScalarKind, func(b524SystemScalarT *b524SystemScalarTesting.T) {
					b524SystemScalarInput := b524SystemScalarControl(b524SystemScalarT, b524SystemScalarRow, b524SystemScalarRow.B524SystemScalarExamples[0].B524SystemScalarHex)
					b524SystemScalarInput.B524SystemScalarReplyKind.B524SystemScalarValue = b524SystemScalarKind
					b524SystemScalarCheck(b524SystemScalarT, b524SystemScalarInput, false)
				})
			}
			for _, b524SystemScalarGrade := range []string{"", "ab58_invalid", "hypothesis"} {
				b524SystemScalarT.Run("ab58_grade_"+b524SystemScalarGrade, func(b524SystemScalarT *b524SystemScalarTesting.T) {
					b524SystemScalarInput := b524SystemScalarControl(b524SystemScalarT, b524SystemScalarRow, b524SystemScalarRow.B524SystemScalarExamples[0].B524SystemScalarHex)
					b524SystemScalarInput.B524SystemScalarCodecGrade.B524SystemScalarValue = b524SystemScalarGrade
					b524SystemScalarCheck(b524SystemScalarT, b524SystemScalarInput, false)
				})
			}
			for _, b524SystemScalarVariant := range []string{"ab58_minus_one", "ab58_plus_one", "ab58_nil", "ab58_empty", "ab58_wrong_codec_width", "ab58_unlisted_group", "ab58_neighbor_low", "ab58_neighbor_high"} {
				b524SystemScalarT.Run(b524SystemScalarVariant, func(b524SystemScalarT *b524SystemScalarTesting.T) {
					b524SystemScalarInput := b524SystemScalarControl(b524SystemScalarT, b524SystemScalarRow, b524SystemScalarRow.B524SystemScalarExamples[0].B524SystemScalarHex)
					switch b524SystemScalarVariant {
					case "ab58_minus_one":
						b524SystemScalarInput.B524SystemScalarRaw = b524SystemScalarInput.B524SystemScalarRaw[:len(b524SystemScalarInput.B524SystemScalarRaw)-1]
					case "ab58_plus_one":
						b524SystemScalarInput.B524SystemScalarRaw = append(b524SystemScalarInput.B524SystemScalarRaw, 0)
					case "ab58_nil":
						b524SystemScalarInput.B524SystemScalarRaw = nil
					case "ab58_empty":
						b524SystemScalarInput.B524SystemScalarRaw = []byte{}
					case "ab58_wrong_codec_width":
						if b524SystemScalarRow.B524SystemScalarCodec == "float32_le" {
							b524SystemScalarInput.B524SystemScalarRaw = []byte{0x34, 0x12}
						} else {
							b524SystemScalarInput.B524SystemScalarRaw = []byte{0, 0, 0xc0, 0x3f}
						}
					case "ab58_unlisted_group":
						b524SystemScalarInput.B524SystemScalarContext.B524SystemScalarGroup.B524SystemScalarValue = 0x7f
						b524SystemScalarInput.B524SystemScalarReturnedGroup.B524SystemScalarValue = 0x7f
					case "ab58_neighbor_low":
						b524SystemScalarInput.B524SystemScalarContext.B524SystemScalarRegister.B524SystemScalarValue = 1
						b524SystemScalarInput.B524SystemScalarReturnedRegister = b524SystemScalarInput.B524SystemScalarContext.B524SystemScalarRegister
					case "ab58_neighbor_high":
						b524SystemScalarInput.B524SystemScalarContext.B524SystemScalarRegister.B524SystemScalarValue = 0xf0
						b524SystemScalarInput.B524SystemScalarReturnedRegister = b524SystemScalarInput.B524SystemScalarContext.B524SystemScalarRegister
					}
					b524SystemScalarResult := b524SystemScalarCheck(b524SystemScalarT, b524SystemScalarInput, false)
					if b524SystemScalarVariant != "ab58_unlisted_group" && b524SystemScalarVariant != "ab58_neighbor_low" && b524SystemScalarVariant != "ab58_neighbor_high" && b524SystemScalarResult.B524SystemScalarDeclaredUnit != b524SystemScalarRow.B524SystemScalarUnit {
						b524SystemScalarT.Fatal("ab58 declared unit lost on partial failure")
					}
				})
			}
			if b524SystemScalarRow.B524SystemScalarCodec == "float32_le" {
				for _, b524SystemScalarBits := range []string{"0000807f", "000080ff", "0000c07f"} {
					b524SystemScalarT.Run("ab58_nonfinite_"+b524SystemScalarBits, func(b524SystemScalarT *b524SystemScalarTesting.T) {
						b524SystemScalarCheck(b524SystemScalarT, b524SystemScalarControl(b524SystemScalarT, b524SystemScalarRow, b524SystemScalarBits), false)
					})
				}
			}
			for _, b524SystemScalarCandidate := range []bool{true, false} {
				b524SystemScalarName := "ab58_memory_candidate"
				if !b524SystemScalarCandidate {
					b524SystemScalarName = "ab58_memory_raw_only"
				}
				b524SystemScalarT.Run(b524SystemScalarName, func(b524SystemScalarT *b524SystemScalarTesting.T) {
					b524SystemScalarInput := b524SystemScalarControl(b524SystemScalarT, b524SystemScalarRow, b524SystemScalarRow.B524SystemScalarExamples[0].B524SystemScalarHex)
					if !b524SystemScalarCandidate {
						b524SystemScalarInput.B524SystemScalarReplyKind.B524SystemScalarValue = "failed"
					}
					b524SystemScalarFirst := b524SystemScalarCheck(b524SystemScalarT, b524SystemScalarInput, b524SystemScalarCandidate)
					b524SystemScalarSecond := b524SystemScalarCheck(b524SystemScalarT, b524SystemScalarInput, b524SystemScalarCandidate)
					b524SystemScalarOriginal := b524SystemScalarInput.B524SystemScalarRaw[0]
					b524SystemScalarInput.B524SystemScalarRaw[0] ^= 0xff
					b524SystemScalarInput.B524SystemScalarContext.B524SystemScalarSoftware.B524SystemScalarValue = "ab58_mutated"
					if b524SystemScalarFirst.B524SystemScalarInput.B524SystemScalarRaw[0] != b524SystemScalarOriginal || b524SystemScalarFirst.B524SystemScalarInput.B524SystemScalarContext.B524SystemScalarSoftware.B524SystemScalarValue != "ab58_software" {
						b524SystemScalarT.Fatal("ab58 input aliases result")
					}
					b524SystemScalarFirst.B524SystemScalarInput.B524SystemScalarRaw[0] ^= 0x55
					if b524SystemScalarInput.B524SystemScalarRaw[0] != b524SystemScalarOriginal^0xff || b524SystemScalarSecond.B524SystemScalarInput.B524SystemScalarRaw[0] != b524SystemScalarOriginal {
						b524SystemScalarT.Fatal("ab58 returned memory aliases caller or another result")
					}
					if b524SystemScalarCandidate {
						b524SystemScalarFirst.B524SystemScalarValue.B524SystemScalarFloat32 = 99
						b524SystemScalarFirst.B524SystemScalarValue.B524SystemScalarUint16 = 99
						if b524SystemScalarSecond.B524SystemScalarValue.B524SystemScalarFloat32 == 99 || b524SystemScalarSecond.B524SystemScalarValue.B524SystemScalarUint16 == 99 {
							b524SystemScalarT.Fatal("ab58 scalar aliases another result")
						}
					}
				})
			}
		})
	}
}

func TestAb58ExplicitZeroContext(b524SystemScalarT *b524SystemScalarTesting.T) {
	for _, b524SystemScalarRow := range b524SystemScalarFixtures(b524SystemScalarT) {
		b524SystemScalarT.Run(b524SystemScalarRow.B524SystemScalarLabel, func(b524SystemScalarT *b524SystemScalarTesting.T) {
			for _, b524SystemScalarField := range []string{"B524SystemScalarSoftware", "B524SystemScalarHardware", "B524SystemScalarAPIVersion", "B524SystemScalarAPIRevision"} {
				b524SystemScalarT.Run("ab58_zero_"+b524SystemScalarField, func(b524SystemScalarT *b524SystemScalarTesting.T) {
					b524SystemScalarInput := b524SystemScalarControl(b524SystemScalarT, b524SystemScalarRow, b524SystemScalarRow.B524SystemScalarExamples[0].B524SystemScalarHex)
					b524SystemScalarReflect.ValueOf(&b524SystemScalarInput.B524SystemScalarContext).Elem().FieldByName(b524SystemScalarField).FieldByName("B524SystemScalarValue").SetString("")
					b524SystemScalarCheck(b524SystemScalarT, b524SystemScalarInput, true)
				})
			}
		})
	}
}

func TestAb58AllConcreteInstances(b524SystemScalarT *b524SystemScalarTesting.T) {
	for _, b524SystemScalarRow := range b524SystemScalarFixtures(b524SystemScalarT) {
		b524SystemScalarT.Run(b524SystemScalarRow.B524SystemScalarLabel, func(b524SystemScalarT *b524SystemScalarTesting.T) {
			for b524SystemScalarInstance := 0; b524SystemScalarInstance < 256; b524SystemScalarInstance++ {
				b524SystemScalarInput := b524SystemScalarControl(b524SystemScalarT, b524SystemScalarRow, b524SystemScalarRow.B524SystemScalarExamples[0].B524SystemScalarHex)
				b524SystemScalarInput.B524SystemScalarContext.B524SystemScalarInstance.B524SystemScalarValue = byte(b524SystemScalarInstance)
				b524SystemScalarCheck(b524SystemScalarT, b524SystemScalarInput, b524SystemScalarInstance == 0 || (b524SystemScalarRow.B524SystemScalarFamily == "controller-700" && (b524SystemScalarRow.B524SystemScalarRegister == 0x42 || b524SystemScalarRow.B524SystemScalarRegister == 0x43)))
			}
		})
	}
}

func TestAb58ExactSelectorAdmission(b524SystemScalarT *b524SystemScalarTesting.T) {
	b524SystemScalarRows := b524SystemScalarFixtures(b524SystemScalarT)
	for _, b524SystemScalarFamily := range []string{"controller-720", "controller-700", "ab58_unknown"} {
		b524SystemScalarT.Run("ab58_"+b524SystemScalarFamily, func(b524SystemScalarT *b524SystemScalarTesting.T) {
			for b524SystemScalarRegister := 0; b524SystemScalarRegister <= 0xffff; b524SystemScalarRegister++ {
				b524SystemScalarListed := false
				for _, b524SystemScalarRow := range b524SystemScalarRows {
					if b524SystemScalarRow.B524SystemScalarFamily == b524SystemScalarFamily && b524SystemScalarRow.B524SystemScalarRegister == uint16(b524SystemScalarRegister) {
						b524SystemScalarListed = true
					}
				}
				if b524SystemScalarListed {
					continue
				}
				for _, b524SystemScalarRaw := range []string{"3412", "0000c03f"} {
					b524SystemScalarInput := b524SystemScalarControl(b524SystemScalarT, b524SystemScalarRows[0], b524SystemScalarRaw)
					b524SystemScalarInput.B524SystemScalarContext.B524SystemScalarFamily.B524SystemScalarValue = b524SystemScalarFamily
					b524SystemScalarInput.B524SystemScalarContext.B524SystemScalarRegister.B524SystemScalarValue = uint16(b524SystemScalarRegister)
					b524SystemScalarInput.B524SystemScalarReturnedRegister = b524SystemScalarInput.B524SystemScalarContext.B524SystemScalarRegister
					b524SystemScalarCheck(b524SystemScalarT, b524SystemScalarInput, false)
				}
			}
		})
	}
}

func TestAb58FiniteValues(b524SystemScalarT *b524SystemScalarTesting.T) {
	for _, b524SystemScalarRow := range b524SystemScalarFixtures(b524SystemScalarT) {
		if b524SystemScalarRow.B524SystemScalarCodec != "float32_le" {
			continue
		}
		b524SystemScalarT.Run(b524SystemScalarRow.B524SystemScalarLabel, func(b524SystemScalarT *b524SystemScalarTesting.T) {
			for _, b524SystemScalarCase := range []struct {
				b524SystemScalarLabel string
				b524SystemScalarHex   string
				b524SystemScalarValue float32
				b524SystemScalarBits  uint32
			}{
				{"ab58_max_finite", "ffff7f7f", 3.4028234663852886e38, 0x7f7fffff},
				{"ab58_min_finite", "ffff7fff", -3.4028234663852886e38, 0xff7fffff},
				{"ab58_subnormal", "01000000", 1.401298464324817e-45, 1},
				{"ab58_zero", "00000000", 0, 0},
				{"ab58_negative_zero", "00000080", 0, 0x80000000},
			} {
				b524SystemScalarT.Run(b524SystemScalarCase.b524SystemScalarLabel, func(b524SystemScalarT *b524SystemScalarTesting.T) {
					b524SystemScalarResult := b524SystemScalarCheck(b524SystemScalarT, b524SystemScalarControl(b524SystemScalarT, b524SystemScalarRow, b524SystemScalarCase.b524SystemScalarHex), true)
					if b524SystemScalarResult.B524SystemScalarValue.B524SystemScalarFloat32 != b524SystemScalarCase.b524SystemScalarValue || b524SystemScalarMath.Float32bits(b524SystemScalarResult.B524SystemScalarValue.B524SystemScalarFloat32) != b524SystemScalarCase.b524SystemScalarBits {
						b524SystemScalarT.Fatal("ab58 finite value clipped, offset or normalized")
					}
				})
			}
		})
	}
}
