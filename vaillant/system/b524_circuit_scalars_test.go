package system

import (
	b524CircuitScalarJSON "encoding/json"
	b524CircuitScalarOS "os"
	b524CircuitScalarReflect "reflect"
	b524CircuitScalarTesting "testing"
)

type b524CircuitScalarFixture struct {
	B524CircuitScalarLabel    string  `json:"Ab62Label"`
	B524CircuitScalarFamily   string  `json:"Ab62Family"`
	B524CircuitScalarRegister uint16  `json:"Ab62Register"`
	B524CircuitScalarUnit     string  `json:"Ab62Unit"`
	B524CircuitScalarRaw      []byte  `json:"Ab62Raw"`
	B524CircuitScalarExpected float32 `json:"Ab62Expected"`
}

func b524CircuitScalarPresent[B524CircuitScalarT any](b524CircuitScalarValue B524CircuitScalarT) B524CircuitScalarField[B524CircuitScalarT] {
	return B524CircuitScalarField[B524CircuitScalarT]{B524CircuitScalarPresent: true, B524CircuitScalarValue: b524CircuitScalarValue}
}

func b524CircuitScalarFixtures(b524CircuitScalarT *b524CircuitScalarTesting.T) []b524CircuitScalarFixture {
	b524CircuitScalarT.Helper()
	b524CircuitScalarBytes, b524CircuitScalarErr := b524CircuitScalarOS.ReadFile("testdata/ab62_b524_circuit_scalars.json")
	if b524CircuitScalarErr != nil {
		b524CircuitScalarT.Fatal(b524CircuitScalarErr)
	}
	var b524CircuitScalarRows []b524CircuitScalarFixture
	if b524CircuitScalarErr = b524CircuitScalarJSON.Unmarshal(b524CircuitScalarBytes, &b524CircuitScalarRows); b524CircuitScalarErr != nil {
		b524CircuitScalarT.Fatal(b524CircuitScalarErr)
	}
	if len(b524CircuitScalarRows) != 22 {
		b524CircuitScalarT.Fatalf("ab62 fixture count: %d", len(b524CircuitScalarRows))
	}
	b524CircuitScalarSeen := map[string]map[uint16]bool{}
	for _, b524CircuitScalarRow := range b524CircuitScalarRows {
		if b524CircuitScalarSeen[b524CircuitScalarRow.B524CircuitScalarFamily] == nil {
			b524CircuitScalarSeen[b524CircuitScalarRow.B524CircuitScalarFamily] = map[uint16]bool{}
		}
		if b524CircuitScalarSeen[b524CircuitScalarRow.B524CircuitScalarFamily][b524CircuitScalarRow.B524CircuitScalarRegister] {
			b524CircuitScalarT.Fatal("ab62 duplicate fixture")
		}
		b524CircuitScalarSeen[b524CircuitScalarRow.B524CircuitScalarFamily][b524CircuitScalarRow.B524CircuitScalarRegister] = true
		if b524CircuitScalarRow.B524CircuitScalarExpected != 1.5 || !b524CircuitScalarReflect.DeepEqual(b524CircuitScalarRow.B524CircuitScalarRaw, []byte{0, 0, 192, 63}) {
			b524CircuitScalarT.Fatal("ab62 independent control changed")
		}
	}
	b524CircuitScalarWant := map[string][]uint16{
		"controller-720": {7, 8, 10, 11, 15, 16, 17, 18, 20, 26, 28, 31},
		"controller-700": {7, 8, 10, 11, 15, 16, 17, 18, 20, 31},
	}
	for b524CircuitScalarFamily, b524CircuitScalarRegisters := range b524CircuitScalarWant {
		for _, b524CircuitScalarRegister := range b524CircuitScalarRegisters {
			if !b524CircuitScalarSeen[b524CircuitScalarFamily][b524CircuitScalarRegister] {
				b524CircuitScalarT.Fatalf("ab62 missing tuple %s/%04X", b524CircuitScalarFamily, b524CircuitScalarRegister)
			}
		}
	}
	return b524CircuitScalarRows
}

func b524CircuitScalarControl(b524CircuitScalarRow b524CircuitScalarFixture) B524CircuitScalarInput {
	return B524CircuitScalarInput{
		B524CircuitScalarContext: B524CircuitScalarContext{
			B524CircuitScalarFamily:   b524CircuitScalarPresent(b524CircuitScalarRow.B524CircuitScalarFamily),
			B524CircuitScalarSoftware: b524CircuitScalarPresent("ab62-software"), B524CircuitScalarHardware: b524CircuitScalarPresent("ab62-hardware"),
			B524CircuitScalarAPIVersion: b524CircuitScalarPresent("ab62-api"), B524CircuitScalarAPIRevision: b524CircuitScalarPresent("ab62-revision"),
			B524CircuitScalarOpcode: b524CircuitScalarPresent(byte(2)), B524CircuitScalarGroup: b524CircuitScalarPresent(byte(2)),
			B524CircuitScalarInstance: b524CircuitScalarPresent(byte(0)), B524CircuitScalarRegister: b524CircuitScalarPresent(b524CircuitScalarRow.B524CircuitScalarRegister),
		},
		B524CircuitScalarReplyKind: b524CircuitScalarPresent("echo_payload"), B524CircuitScalarReturnedGroup: b524CircuitScalarPresent(byte(2)),
		B524CircuitScalarReturnedRegister: b524CircuitScalarPresent(b524CircuitScalarRow.B524CircuitScalarRegister), B524CircuitScalarCodecGrade: b524CircuitScalarPresent(B524CircuitScalarHypothesis),
		B524CircuitScalarRaw: append([]byte(nil), b524CircuitScalarRow.B524CircuitScalarRaw...),
	}
}

func b524CircuitScalarCheck(b524CircuitScalarT *b524CircuitScalarTesting.T, b524CircuitScalarInput B524CircuitScalarInput, b524CircuitScalarCandidate bool, b524CircuitScalarExpected float32) B524CircuitScalarResult {
	b524CircuitScalarT.Helper()
	b524CircuitScalarResult := B524InterpretCircuitScalar(b524CircuitScalarInput)
	if !b524CircuitScalarReflect.DeepEqual(b524CircuitScalarResult.B524CircuitScalarInput, b524CircuitScalarInput) {
		b524CircuitScalarT.Fatalf("ab62 context lost: %+v versus %+v", b524CircuitScalarResult.B524CircuitScalarInput, b524CircuitScalarInput)
	}
	if b524CircuitScalarCandidate {
		if b524CircuitScalarResult.B524CircuitScalarState != B524CircuitScalarCandidate || b524CircuitScalarResult.B524CircuitScalarGrade != B524CircuitScalarHypothesis || !b524CircuitScalarResult.B524CircuitScalarValue.B524CircuitScalarPresent || b524CircuitScalarResult.B524CircuitScalarValue.B524CircuitScalarValue != b524CircuitScalarExpected || b524CircuitScalarResult.B524CircuitScalarCodec != "float32_le" {
			b524CircuitScalarT.Fatalf("ab62 candidate: %+v", b524CircuitScalarResult)
		}
	} else if b524CircuitScalarResult.B524CircuitScalarState != B524CircuitScalarRawOnly || b524CircuitScalarResult.B524CircuitScalarGrade != B524CircuitScalarUnknown || b524CircuitScalarResult.B524CircuitScalarValue != (B524CircuitScalarField[float32]{}) {
		b524CircuitScalarT.Fatalf("ab62 raw-only: %+v", b524CircuitScalarResult)
	}
	return b524CircuitScalarResult
}

func TestAb62EveryBinding(b524CircuitScalarT *b524CircuitScalarTesting.T) {
	for _, b524CircuitScalarRow := range b524CircuitScalarFixtures(b524CircuitScalarT) {
		b524CircuitScalarT.Run(b524CircuitScalarRow.B524CircuitScalarLabel, func(b524CircuitScalarT *b524CircuitScalarTesting.T) {
			b524CircuitScalarInput := b524CircuitScalarControl(b524CircuitScalarRow)
			b524CircuitScalarResult := b524CircuitScalarCheck(b524CircuitScalarT, b524CircuitScalarInput, true, b524CircuitScalarRow.B524CircuitScalarExpected)
			if b524CircuitScalarResult.B524CircuitScalarDeclaredUnit != b524CircuitScalarRow.B524CircuitScalarUnit {
				b524CircuitScalarT.Fatalf("ab62 unit: %q", b524CircuitScalarResult.B524CircuitScalarDeclaredUnit)
			}
			for _, b524CircuitScalarCase := range []struct {
				b524CircuitScalarLabel string
				b524CircuitScalarRaw   []byte
				b524CircuitScalarValid bool
				b524CircuitScalarValue float32
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
				b524CircuitScalarT.Run(b524CircuitScalarCase.b524CircuitScalarLabel, func(b524CircuitScalarT *b524CircuitScalarTesting.T) {
					b524CircuitScalarInput.B524CircuitScalarRaw = b524CircuitScalarCase.b524CircuitScalarRaw
					b524CircuitScalarResult := b524CircuitScalarCheck(b524CircuitScalarT, b524CircuitScalarInput, b524CircuitScalarCase.b524CircuitScalarValid, b524CircuitScalarCase.b524CircuitScalarValue)
					if b524CircuitScalarResult.B524CircuitScalarDeclaredUnit != b524CircuitScalarRow.B524CircuitScalarUnit {
						b524CircuitScalarT.Fatal("ab62 declared unit lost")
					}
				})
			}
		})
	}
}

func TestAb62IndependentAxes(b524CircuitScalarT *b524CircuitScalarTesting.T) {
	for _, b524CircuitScalarRow := range b524CircuitScalarFixtures(b524CircuitScalarT) {
		b524CircuitScalarT.Run(b524CircuitScalarRow.B524CircuitScalarLabel, func(b524CircuitScalarT *b524CircuitScalarTesting.T) {
			for _, b524CircuitScalarAxis := range []struct {
				b524CircuitScalarName      string
				b524CircuitScalarValue     any
				b524CircuitScalarCandidate bool
			}{
				{"B524CircuitScalarFamily", "ab62-unknown", false},
				{"B524CircuitScalarSoftware", "ab62-other-software", true},
				{"B524CircuitScalarHardware", "ab62-other-hardware", true},
				{"B524CircuitScalarAPIVersion", "ab62-other-api", true},
				{"B524CircuitScalarAPIRevision", "ab62-other-revision", true},
				{"B524CircuitScalarOpcode", byte(6), false},
				{"B524CircuitScalarGroup", byte(3), false},
				{"B524CircuitScalarInstance", byte(1), true},
				{"B524CircuitScalarRegister", uint16(6), false},
				{"B524CircuitScalarReplyKind", "status_only", false},
				{"B524CircuitScalarReturnedGroup", byte(3), false},
				{"B524CircuitScalarReturnedRegister", uint16(6), false},
				{"B524CircuitScalarCodecGrade", B524CircuitScalarProven, true},
			} {
				b524CircuitScalarT.Run(b524CircuitScalarAxis.b524CircuitScalarName, func(b524CircuitScalarT *b524CircuitScalarTesting.T) {
					b524CircuitScalarInput := b524CircuitScalarControl(b524CircuitScalarRow)
					b524CircuitScalarField := b524CircuitScalarReflect.ValueOf(&b524CircuitScalarInput).Elem().FieldByName(b524CircuitScalarAxis.b524CircuitScalarName)
					if !b524CircuitScalarField.IsValid() {
						b524CircuitScalarField = b524CircuitScalarReflect.ValueOf(&b524CircuitScalarInput.B524CircuitScalarContext).Elem().FieldByName(b524CircuitScalarAxis.b524CircuitScalarName)
					}
					b524CircuitScalarField.FieldByName("B524CircuitScalarValue").Set(b524CircuitScalarReflect.ValueOf(b524CircuitScalarAxis.b524CircuitScalarValue))
					b524CircuitScalarCheck(b524CircuitScalarT, b524CircuitScalarInput, b524CircuitScalarAxis.b524CircuitScalarCandidate, 1.5)
					b524CircuitScalarInput = b524CircuitScalarControl(b524CircuitScalarRow)
					b524CircuitScalarField = b524CircuitScalarReflect.ValueOf(&b524CircuitScalarInput).Elem().FieldByName(b524CircuitScalarAxis.b524CircuitScalarName)
					if !b524CircuitScalarField.IsValid() {
						b524CircuitScalarField = b524CircuitScalarReflect.ValueOf(&b524CircuitScalarInput.B524CircuitScalarContext).Elem().FieldByName(b524CircuitScalarAxis.b524CircuitScalarName)
					}
					b524CircuitScalarField.FieldByName("B524CircuitScalarPresent").SetBool(false)
					b524CircuitScalarCheck(b524CircuitScalarT, b524CircuitScalarInput, false, 0)
				})
			}
			for _, b524CircuitScalarName := range []string{"B524CircuitScalarSoftware", "B524CircuitScalarHardware", "B524CircuitScalarAPIVersion", "B524CircuitScalarAPIRevision"} {
				b524CircuitScalarInput := b524CircuitScalarControl(b524CircuitScalarRow)
				b524CircuitScalarField := b524CircuitScalarReflect.ValueOf(&b524CircuitScalarInput.B524CircuitScalarContext).Elem().FieldByName(b524CircuitScalarName)
				b524CircuitScalarField.FieldByName("B524CircuitScalarValue").SetString("")
				b524CircuitScalarCheck(b524CircuitScalarT, b524CircuitScalarInput, true, 1.5)
				b524CircuitScalarField.FieldByName("B524CircuitScalarPresent").SetBool(false)
				b524CircuitScalarCheck(b524CircuitScalarT, b524CircuitScalarInput, false, 0)
			}
			b524CircuitScalarInput := b524CircuitScalarControl(b524CircuitScalarRow)
			b524CircuitScalarInput.B524CircuitScalarContext.B524CircuitScalarInstance.B524CircuitScalarPresent = false
			b524CircuitScalarCheck(b524CircuitScalarT, b524CircuitScalarInput, false, 0)
			for _, b524CircuitScalarKind := range []string{"", "short", "failed", "unknown", "status_only", "ab62-unrecognized"} {
				b524CircuitScalarInput = b524CircuitScalarControl(b524CircuitScalarRow)
				b524CircuitScalarInput.B524CircuitScalarReplyKind = b524CircuitScalarPresent(b524CircuitScalarKind)
				b524CircuitScalarCheck(b524CircuitScalarT, b524CircuitScalarInput, false, 0)
			}
			for _, b524CircuitScalarGrade := range []B524CircuitScalarGrade{"", B524CircuitScalarUnknown, "ab62-unrecognized"} {
				b524CircuitScalarInput = b524CircuitScalarControl(b524CircuitScalarRow)
				b524CircuitScalarInput.B524CircuitScalarCodecGrade = b524CircuitScalarPresent(b524CircuitScalarGrade)
				b524CircuitScalarCheck(b524CircuitScalarT, b524CircuitScalarInput, false, 0)
			}
			for _, b524CircuitScalarInstance := range []byte{0, 1, 255} {
				b524CircuitScalarInput = b524CircuitScalarControl(b524CircuitScalarRow)
				b524CircuitScalarInput.B524CircuitScalarContext.B524CircuitScalarInstance = b524CircuitScalarPresent(b524CircuitScalarInstance)
				b524CircuitScalarCheck(b524CircuitScalarT, b524CircuitScalarInput, true, 1.5)
			}
			b524CircuitScalarInput = b524CircuitScalarControl(b524CircuitScalarRow)
			b524CircuitScalarInput.B524CircuitScalarContext.B524CircuitScalarFamily.B524CircuitScalarValue = "controller-700"
			if b524CircuitScalarRow.B524CircuitScalarFamily == "controller-700" {
				b524CircuitScalarInput.B524CircuitScalarContext.B524CircuitScalarFamily.B524CircuitScalarValue = "controller-720"
			}
			b524CircuitScalarCheck(b524CircuitScalarT, b524CircuitScalarInput, b524CircuitScalarRow.B524CircuitScalarRegister != 0x1a && b524CircuitScalarRow.B524CircuitScalarRegister != 0x1c, 1.5)
			b524CircuitScalarInput = b524CircuitScalarControl(b524CircuitScalarRow)
			b524CircuitScalarInput.B524CircuitScalarContext.B524CircuitScalarGroup = b524CircuitScalarPresent(byte(3))
			b524CircuitScalarInput.B524CircuitScalarReturnedGroup = b524CircuitScalarPresent(byte(3))
			b524CircuitScalarCheck(b524CircuitScalarT, b524CircuitScalarInput, false, 0)
			b524CircuitScalarInput = b524CircuitScalarControl(b524CircuitScalarRow)
			b524CircuitScalarInput.B524CircuitScalarContext.B524CircuitScalarOpcode = b524CircuitScalarPresent(byte(6))
			b524CircuitScalarCheck(b524CircuitScalarT, b524CircuitScalarInput, false, 0)
			b524CircuitScalarInput = b524CircuitScalarControl(b524CircuitScalarRow)
			b524CircuitScalarInput.B524CircuitScalarContext.B524CircuitScalarRegister = b524CircuitScalarPresent(uint16(0x20))
			b524CircuitScalarInput.B524CircuitScalarReturnedRegister = b524CircuitScalarPresent(uint16(0x20))
			b524CircuitScalarCheck(b524CircuitScalarT, b524CircuitScalarInput, false, 0)
		})
	}
}

func TestAb62OwnedMemory(b524CircuitScalarT *b524CircuitScalarTesting.T) {
	for _, b524CircuitScalarRow := range b524CircuitScalarFixtures(b524CircuitScalarT) {
		for _, b524CircuitScalarKind := range []string{"echo_payload", "failed"} {
			b524CircuitScalarInput := b524CircuitScalarControl(b524CircuitScalarRow)
			b524CircuitScalarInput.B524CircuitScalarReplyKind = b524CircuitScalarPresent(b524CircuitScalarKind)
			b524CircuitScalarResult := b524CircuitScalarCheck(b524CircuitScalarT, b524CircuitScalarInput, b524CircuitScalarKind == "echo_payload", 1.5)
			b524CircuitScalarInput.B524CircuitScalarRaw[0] = 255
			if b524CircuitScalarResult.B524CircuitScalarInput.B524CircuitScalarRaw[0] != 0 {
				b524CircuitScalarT.Fatal("ab62 aliases input")
			}
			b524CircuitScalarResult.B524CircuitScalarInput.B524CircuitScalarRaw[1] = 255
			if b524CircuitScalarInput.B524CircuitScalarRaw[1] != 0 {
				b524CircuitScalarT.Fatal("ab62 aliases output")
			}
			b524CircuitScalarOther := B524InterpretCircuitScalar(b524CircuitScalarControl(b524CircuitScalarRow))
			if b524CircuitScalarOther.B524CircuitScalarInput.B524CircuitScalarRaw[1] != 0 {
				b524CircuitScalarT.Fatal("ab62 shared output")
			}
			if b524CircuitScalarKind == "echo_payload" && b524CircuitScalarResult.B524CircuitScalarValue.B524CircuitScalarValue != 1.5 {
				b524CircuitScalarT.Fatal("ab62 value mutated")
			}
		}
	}
}
