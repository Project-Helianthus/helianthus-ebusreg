package system

import (
	b524ZoneScalarJSON "encoding/json"
	b524ZoneScalarOS "os"
	b524ZoneScalarReflect "reflect"
	b524ZoneScalarTesting "testing"
)

type b524ZoneScalarFixture struct {
	B524ZoneScalarLabel    string  `json:"ab64Label"`
	B524ZoneScalarFamily   string  `json:"ab64Family"`
	B524ZoneScalarRegister uint16  `json:"ab64Register"`
	B524ZoneScalarUnit     string  `json:"ab64Unit"`
	B524ZoneScalarRaw      []byte  `json:"ab64Raw"`
	B524ZoneScalarExpected float32 `json:"ab64Expected"`
}

func b524ZoneScalarValid(b524ZoneScalarRow b524ZoneScalarFixture) B524ZoneScalarInput {
	return B524ZoneScalarInput{
		B524ZoneScalarContext: B524ZoneScalarContext{
			B524ZoneScalarFamily:      B524ZoneScalarField[string]{true, b524ZoneScalarRow.B524ZoneScalarFamily},
			B524ZoneScalarSoftware:    B524ZoneScalarField[string]{true, "ab64-software"},
			B524ZoneScalarHardware:    B524ZoneScalarField[string]{true, "ab64-hardware"},
			B524ZoneScalarAPIVersion:  B524ZoneScalarField[string]{true, "ab64-api"},
			B524ZoneScalarAPIRevision: B524ZoneScalarField[string]{true, "ab64-revision"},
			B524ZoneScalarOpcode:      B524ZoneScalarField[byte]{true, 2}, B524ZoneScalarGroup: B524ZoneScalarField[byte]{true, 3},
			B524ZoneScalarInstance: B524ZoneScalarField[byte]{true, 0}, B524ZoneScalarRegister: B524ZoneScalarField[uint16]{true, b524ZoneScalarRow.B524ZoneScalarRegister},
		},
		B524ZoneScalarReplyKind: "echo_payload", B524ZoneScalarReturnedGroup: B524ZoneScalarField[byte]{true, 3},
		B524ZoneScalarReturnedRegister: B524ZoneScalarField[uint16]{true, b524ZoneScalarRow.B524ZoneScalarRegister},
		B524ZoneScalarCodecGrade:       B524ZoneScalarHypothesis, B524ZoneScalarRaw: append([]byte(nil), b524ZoneScalarRow.B524ZoneScalarRaw...),
	}
}

func b524ZoneScalarCheck(b524ZoneScalarT *b524ZoneScalarTesting.T, b524ZoneScalarInput B524ZoneScalarInput, b524ZoneScalarCandidate bool, b524ZoneScalarExpected float32, b524ZoneScalarUnit string) B524ZoneScalarResult {
	b524ZoneScalarT.Helper()
	b524ZoneScalarResult := B524ZoneScalarInterpret(b524ZoneScalarInput)
	if !b524ZoneScalarReflect.DeepEqual(b524ZoneScalarResult.B524ZoneScalarInput, b524ZoneScalarInput) {
		b524ZoneScalarT.Fatalf("ab64 lost context/raw: got %#v want %#v", b524ZoneScalarResult.B524ZoneScalarInput, b524ZoneScalarInput)
	}
	if b524ZoneScalarResult.B524ZoneScalarUnit != b524ZoneScalarUnit {
		b524ZoneScalarT.Fatalf("ab64 unit: %q want %q", b524ZoneScalarResult.B524ZoneScalarUnit, b524ZoneScalarUnit)
	}
	if b524ZoneScalarCandidate {
		if b524ZoneScalarResult.B524ZoneScalarState != B524ZoneScalarCandidate || !b524ZoneScalarResult.B524ZoneScalarValue.B524ZoneScalarPresent || b524ZoneScalarResult.B524ZoneScalarValue.B524ZoneScalarValue != b524ZoneScalarExpected || b524ZoneScalarResult.B524ZoneScalarGrade != B524ZoneScalarHypothesis || b524ZoneScalarResult.B524ZoneScalarCodec != "float32_le" || b524ZoneScalarResult.B524ZoneScalarInstanceRule != "dynamic" {
			b524ZoneScalarT.Fatalf("ab64 candidate: %#v", b524ZoneScalarResult)
		}
	} else if b524ZoneScalarResult.B524ZoneScalarState != B524ZoneScalarRawOnly || b524ZoneScalarResult.B524ZoneScalarValue != (B524ZoneScalarField[float32]{}) || b524ZoneScalarResult.B524ZoneScalarGrade != B524ZoneScalarUnknown {
		b524ZoneScalarT.Fatalf("ab64 raw-only: %#v", b524ZoneScalarResult)
	}
	for _, b524ZoneScalarGrade := range []B524ZoneScalarGrade{b524ZoneScalarResult.B524ZoneScalarNativeSupport, b524ZoneScalarResult.B524ZoneScalarQualification, b524ZoneScalarResult.B524ZoneScalarPresence, b524ZoneScalarResult.B524ZoneScalarFreshness, b524ZoneScalarResult.B524ZoneScalarExecution, b524ZoneScalarResult.B524ZoneScalarEffects, b524ZoneScalarResult.B524ZoneScalarPersistence, b524ZoneScalarResult.B524ZoneScalarDelivery} {
		if b524ZoneScalarGrade != B524ZoneScalarUnknown {
			b524ZoneScalarT.Fatal("ab64 native property promoted")
		}
	}
	return b524ZoneScalarResult
}

func TestAb64Bindings(b524ZoneScalarT *b524ZoneScalarTesting.T) {
	b524ZoneScalarData, b524ZoneScalarErr := b524ZoneScalarOS.ReadFile("testdata/ab64_b524_zone_scalars.json")
	if b524ZoneScalarErr != nil {
		b524ZoneScalarT.Fatal(b524ZoneScalarErr)
	}
	var b524ZoneScalarRows []b524ZoneScalarFixture
	if b524ZoneScalarErr = b524ZoneScalarJSON.Unmarshal(b524ZoneScalarData, &b524ZoneScalarRows); b524ZoneScalarErr != nil {
		b524ZoneScalarT.Fatal(b524ZoneScalarErr)
	}
	b524ZoneScalarManifest := map[string][]uint16{
		"controller-720": {2, 5, 9, 15, 20, 21, 35, 40}, "controller-700": {2, 5, 9, 15, 20, 21, 40},
	}
	b524ZoneScalarSeen := make(map[string]map[uint16]bool)
	if len(b524ZoneScalarRows) != 15 {
		b524ZoneScalarT.Fatal("ab64 fixture must contain all 15 tuples")
	}
	for _, b524ZoneScalarRow := range b524ZoneScalarRows {
		b524ZoneScalarFound := false
		for _, b524ZoneScalarRegister := range b524ZoneScalarManifest[b524ZoneScalarRow.B524ZoneScalarFamily] {
			if b524ZoneScalarRegister == b524ZoneScalarRow.B524ZoneScalarRegister {
				b524ZoneScalarFound = true
			}
		}
		if b524ZoneScalarSeen[b524ZoneScalarRow.B524ZoneScalarFamily] == nil {
			b524ZoneScalarSeen[b524ZoneScalarRow.B524ZoneScalarFamily] = make(map[uint16]bool)
		}
		if !b524ZoneScalarFound || b524ZoneScalarSeen[b524ZoneScalarRow.B524ZoneScalarFamily][b524ZoneScalarRow.B524ZoneScalarRegister] {
			b524ZoneScalarT.Fatal("ab64 unexpected/duplicate tuple")
		}
		b524ZoneScalarSeen[b524ZoneScalarRow.B524ZoneScalarFamily][b524ZoneScalarRow.B524ZoneScalarRegister] = true
		b524ZoneScalarUnit := "°C"
		if b524ZoneScalarRow.B524ZoneScalarRegister == 40 {
			b524ZoneScalarUnit = "%"
		}
		if b524ZoneScalarRow.B524ZoneScalarExpected != 1.5 || b524ZoneScalarRow.B524ZoneScalarUnit != b524ZoneScalarUnit || !b524ZoneScalarReflect.DeepEqual(b524ZoneScalarRow.B524ZoneScalarRaw, []byte{0, 0, 192, 63}) {
			b524ZoneScalarT.Fatal("ab64 independent fixture expectation")
		}
		b524ZoneScalarT.Run(b524ZoneScalarRow.B524ZoneScalarLabel, func(b524ZoneScalarT *b524ZoneScalarTesting.T) {
			b524ZoneScalarCheck(b524ZoneScalarT, b524ZoneScalarValid(b524ZoneScalarRow), true, 1.5, b524ZoneScalarUnit)
			b524ZoneScalarCases := []struct {
				b524ZoneScalarName      string
				b524ZoneScalarChange    func(*B524ZoneScalarInput)
				b524ZoneScalarCandidate bool
				b524ZoneScalarExpected  float32
				b524ZoneScalarUnit      string
			}{
				{"ab64-negative", func(b524ZoneScalarI *B524ZoneScalarInput) { b524ZoneScalarI.B524ZoneScalarRaw = []byte{0, 0, 128, 191} }, true, -1, b524ZoneScalarUnit},
				{"ab64-short", func(b524ZoneScalarI *B524ZoneScalarInput) { b524ZoneScalarI.B524ZoneScalarRaw = []byte{0, 0, 192} }, false, 0, b524ZoneScalarUnit},
				{"ab64-extra", func(b524ZoneScalarI *B524ZoneScalarInput) {
					b524ZoneScalarI.B524ZoneScalarRaw = append(b524ZoneScalarI.B524ZoneScalarRaw, 0)
				}, false, 0, b524ZoneScalarUnit},
				{"ab64-width-two", func(b524ZoneScalarI *B524ZoneScalarInput) { b524ZoneScalarI.B524ZoneScalarRaw = []byte{52, 18} }, false, 0, b524ZoneScalarUnit},
				{"ab64-empty", func(b524ZoneScalarI *B524ZoneScalarInput) { b524ZoneScalarI.B524ZoneScalarRaw = nil }, false, 0, b524ZoneScalarUnit},
				{"ab64-inf", func(b524ZoneScalarI *B524ZoneScalarInput) { b524ZoneScalarI.B524ZoneScalarRaw = []byte{0, 0, 128, 127} }, false, 0, b524ZoneScalarUnit},
				{"ab64-negative-inf", func(b524ZoneScalarI *B524ZoneScalarInput) { b524ZoneScalarI.B524ZoneScalarRaw = []byte{0, 0, 128, 255} }, false, 0, b524ZoneScalarUnit},
				{"ab64-nan", func(b524ZoneScalarI *B524ZoneScalarInput) { b524ZoneScalarI.B524ZoneScalarRaw = []byte{0, 0, 192, 127} }, false, 0, b524ZoneScalarUnit},
				{"ab64-software", func(b524ZoneScalarI *B524ZoneScalarInput) {
					b524ZoneScalarI.B524ZoneScalarContext.B524ZoneScalarSoftware.B524ZoneScalarValue = "ab64-other"
				}, true, 1.5, b524ZoneScalarUnit},
				{"ab64-hardware", func(b524ZoneScalarI *B524ZoneScalarInput) {
					b524ZoneScalarI.B524ZoneScalarContext.B524ZoneScalarHardware.B524ZoneScalarValue = "ab64-other"
				}, true, 1.5, b524ZoneScalarUnit},
				{"ab64-api-version", func(b524ZoneScalarI *B524ZoneScalarInput) {
					b524ZoneScalarI.B524ZoneScalarContext.B524ZoneScalarAPIVersion.B524ZoneScalarValue = "ab64-other"
				}, true, 1.5, b524ZoneScalarUnit},
				{"ab64-api-revision", func(b524ZoneScalarI *B524ZoneScalarInput) {
					b524ZoneScalarI.B524ZoneScalarContext.B524ZoneScalarAPIRevision.B524ZoneScalarValue = "ab64-other"
				}, true, 1.5, b524ZoneScalarUnit},
				{"ab64-instance-one", func(b524ZoneScalarI *B524ZoneScalarInput) {
					b524ZoneScalarI.B524ZoneScalarContext.B524ZoneScalarInstance.B524ZoneScalarValue = 1
				}, true, 1.5, b524ZoneScalarUnit},
				{"ab64-instance-max", func(b524ZoneScalarI *B524ZoneScalarInput) {
					b524ZoneScalarI.B524ZoneScalarContext.B524ZoneScalarInstance.B524ZoneScalarValue = 255
				}, true, 1.5, b524ZoneScalarUnit},
				{"ab64-family", func(b524ZoneScalarI *B524ZoneScalarInput) {
					b524ZoneScalarI.B524ZoneScalarContext.B524ZoneScalarFamily.B524ZoneScalarValue = "ab64-unknown"
				}, false, 0, "Unknown"},
				{"ab64-op06", func(b524ZoneScalarI *B524ZoneScalarInput) {
					b524ZoneScalarI.B524ZoneScalarContext.B524ZoneScalarOpcode.B524ZoneScalarValue = 6
				}, false, 0, "Unknown"},
				{"ab64-group", func(b524ZoneScalarI *B524ZoneScalarInput) {
					b524ZoneScalarI.B524ZoneScalarContext.B524ZoneScalarGroup.B524ZoneScalarValue = 4
				}, false, 0, "Unknown"},
				{"ab64-register-neighbor", func(b524ZoneScalarI *B524ZoneScalarInput) {
					b524ZoneScalarI.B524ZoneScalarContext.B524ZoneScalarRegister.B524ZoneScalarValue = b524ZoneScalarRow.B524ZoneScalarRegister + 1
					if b524ZoneScalarRow.B524ZoneScalarRegister == 20 {
						b524ZoneScalarI.B524ZoneScalarContext.B524ZoneScalarRegister.B524ZoneScalarValue = 19
					}
				}, false, 0, "Unknown"},
				{"ab64-returned-group", func(b524ZoneScalarI *B524ZoneScalarInput) {
					b524ZoneScalarI.B524ZoneScalarReturnedGroup.B524ZoneScalarValue = 4
				}, false, 0, b524ZoneScalarUnit},
				{"ab64-returned-register", func(b524ZoneScalarI *B524ZoneScalarInput) {
					b524ZoneScalarI.B524ZoneScalarReturnedRegister.B524ZoneScalarValue++
				}, false, 0, b524ZoneScalarUnit},
				{"ab64-returned-group-absent", func(b524ZoneScalarI *B524ZoneScalarInput) {
					b524ZoneScalarI.B524ZoneScalarReturnedGroup.B524ZoneScalarPresent = false
				}, false, 0, b524ZoneScalarUnit},
				{"ab64-returned-register-absent", func(b524ZoneScalarI *B524ZoneScalarInput) {
					b524ZoneScalarI.B524ZoneScalarReturnedRegister.B524ZoneScalarPresent = false
				}, false, 0, b524ZoneScalarUnit},
				{"ab64-grade-proven", func(b524ZoneScalarI *B524ZoneScalarInput) {
					b524ZoneScalarI.B524ZoneScalarCodecGrade = B524ZoneScalarProven
				}, true, 1.5, b524ZoneScalarUnit},
				{"ab64-grade-unknown", func(b524ZoneScalarI *B524ZoneScalarInput) {
					b524ZoneScalarI.B524ZoneScalarCodecGrade = B524ZoneScalarUnknown
				}, false, 0, b524ZoneScalarUnit},
				{"ab64-grade-missing", func(b524ZoneScalarI *B524ZoneScalarInput) { b524ZoneScalarI.B524ZoneScalarCodecGrade = "" }, false, 0, b524ZoneScalarUnit},
				{"ab64-grade-invalid", func(b524ZoneScalarI *B524ZoneScalarInput) { b524ZoneScalarI.B524ZoneScalarCodecGrade = "ab64-invalid" }, false, 0, b524ZoneScalarUnit},
			}
			for _, b524ZoneScalarCase := range b524ZoneScalarCases {
				b524ZoneScalarT.Run(b524ZoneScalarCase.b524ZoneScalarName, func(b524ZoneScalarT *b524ZoneScalarTesting.T) {
					b524ZoneScalarInput := b524ZoneScalarValid(b524ZoneScalarRow)
					b524ZoneScalarCase.b524ZoneScalarChange(&b524ZoneScalarInput)
					b524ZoneScalarCheck(b524ZoneScalarT, b524ZoneScalarInput, b524ZoneScalarCase.b524ZoneScalarCandidate, b524ZoneScalarCase.b524ZoneScalarExpected, b524ZoneScalarCase.b524ZoneScalarUnit)
				})
			}
			for _, b524ZoneScalarKind := range []string{"", "status_only", "short", "failed", "mismatched", "unknown", "ab64-other"} {
				b524ZoneScalarT.Run("ab64-kind-"+b524ZoneScalarKind, func(b524ZoneScalarT *b524ZoneScalarTesting.T) {
					b524ZoneScalarInput := b524ZoneScalarValid(b524ZoneScalarRow)
					b524ZoneScalarInput.B524ZoneScalarReplyKind = b524ZoneScalarKind
					b524ZoneScalarCheck(b524ZoneScalarT, b524ZoneScalarInput, false, 0, b524ZoneScalarUnit)
				})
			}
			b524ZoneScalarT.Run("ab64-missing-context", func(b524ZoneScalarT *b524ZoneScalarTesting.T) {
				for b524ZoneScalarIndex := 0; b524ZoneScalarIndex < 9; b524ZoneScalarIndex++ {
					b524ZoneScalarInput := b524ZoneScalarValid(b524ZoneScalarRow)
					b524ZoneScalarField := b524ZoneScalarReflect.ValueOf(&b524ZoneScalarInput.B524ZoneScalarContext).Elem().Field(b524ZoneScalarIndex)
					b524ZoneScalarField.FieldByName("B524ZoneScalarPresent").SetBool(false)
					b524ZoneScalarCheck(b524ZoneScalarT, b524ZoneScalarInput, false, 0, "Unknown")
				}
			})
			b524ZoneScalarT.Run("ab64-explicit-zero-context", func(b524ZoneScalarT *b524ZoneScalarTesting.T) {
				for b524ZoneScalarIndex := 0; b524ZoneScalarIndex < 9; b524ZoneScalarIndex++ {
					b524ZoneScalarInput := b524ZoneScalarValid(b524ZoneScalarRow)
					b524ZoneScalarField := b524ZoneScalarReflect.ValueOf(&b524ZoneScalarInput.B524ZoneScalarContext).Elem().Field(b524ZoneScalarIndex).FieldByName("B524ZoneScalarValue")
					b524ZoneScalarField.Set(b524ZoneScalarReflect.Zero(b524ZoneScalarField.Type()))
					b524ZoneScalarCandidate := b524ZoneScalarIndex >= 1 && b524ZoneScalarIndex <= 4 || b524ZoneScalarIndex == 7
					b524ZoneScalarExpectedUnit := "Unknown"
					if b524ZoneScalarCandidate {
						b524ZoneScalarExpectedUnit = b524ZoneScalarUnit
					}
					b524ZoneScalarCheck(b524ZoneScalarT, b524ZoneScalarInput, b524ZoneScalarCandidate, 1.5, b524ZoneScalarExpectedUnit)
				}
			})
			b524ZoneScalarT.Run("ab64-other-family", func(b524ZoneScalarT *b524ZoneScalarTesting.T) {
				b524ZoneScalarInput := b524ZoneScalarValid(b524ZoneScalarRow)
				if b524ZoneScalarRow.B524ZoneScalarFamily == "controller-720" {
					b524ZoneScalarInput.B524ZoneScalarContext.B524ZoneScalarFamily.B524ZoneScalarValue = "controller-700"
				} else {
					b524ZoneScalarInput.B524ZoneScalarContext.B524ZoneScalarFamily.B524ZoneScalarValue = "controller-720"
				}
				b524ZoneScalarExpectedUnit := b524ZoneScalarUnit
				if b524ZoneScalarRow.B524ZoneScalarRegister == 35 {
					b524ZoneScalarExpectedUnit = "Unknown"
				}
				b524ZoneScalarCheck(b524ZoneScalarT, b524ZoneScalarInput, b524ZoneScalarRow.B524ZoneScalarRegister != 35, 1.5, b524ZoneScalarExpectedUnit)
			})
			b524ZoneScalarT.Run("ab64-owned-memory", func(b524ZoneScalarT *b524ZoneScalarTesting.T) {
				for _, b524ZoneScalarKind := range []string{"echo_payload", "status_only"} {
					b524ZoneScalarInput := b524ZoneScalarValid(b524ZoneScalarRow)
					b524ZoneScalarInput.B524ZoneScalarReplyKind = b524ZoneScalarKind
					b524ZoneScalarResult := b524ZoneScalarCheck(b524ZoneScalarT, b524ZoneScalarInput, b524ZoneScalarKind == "echo_payload", 1.5, b524ZoneScalarUnit)
					b524ZoneScalarInput.B524ZoneScalarRaw[0] = 42
					if b524ZoneScalarResult.B524ZoneScalarInput.B524ZoneScalarRaw[0] != 0 {
						b524ZoneScalarT.Fatal("ab64 input aliases result")
					}
					b524ZoneScalarResult.B524ZoneScalarInput.B524ZoneScalarRaw[1] = 43
					if b524ZoneScalarInput.B524ZoneScalarRaw[1] != 0 {
						b524ZoneScalarT.Fatal("ab64 result aliases input")
					}
					b524ZoneScalarAgain := B524ZoneScalarInterpret(b524ZoneScalarValid(b524ZoneScalarRow))
					if b524ZoneScalarAgain.B524ZoneScalarInput.B524ZoneScalarRaw[1] != 0 {
						b524ZoneScalarT.Fatal("ab64 results share storage")
					}
				}
			})
		})
	}
}
