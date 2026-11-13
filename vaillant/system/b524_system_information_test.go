package system

import (
	b524SystemInformationBytes "bytes"
	b524SystemInformationBinary "encoding/binary"
	b524SystemInformationJSON "encoding/json"
	b524SystemInformationMath "math"
	b524SystemInformationOs "os"
	b524SystemInformationReflect "reflect"
	b524SystemInformationTesting "testing"
)

func b524SystemInformationTestPointer[B524SystemInformationType any](b524SystemInformationValue B524SystemInformationType) *B524SystemInformationType {
	return &b524SystemInformationValue
}

func TestAb42BuildInformation(b524SystemInformationT *b524SystemInformationTesting.T) {
	for _, b524SystemInformationCase := range []struct {
		b524SystemInformationIdentifier uint16
		b524SystemInformationExpected   []byte
	}{
		{0x1234, []byte{0, 0x34, 0x12}},
		{0, []byte{0, 0, 0}},
		{1, []byte{0, 1, 0}},
		{0xffff, []byte{0, 0xff, 0xff}},
	} {
		b524SystemInformationFirst := B524SystemBuildInformation(b524SystemInformationCase.b524SystemInformationIdentifier)
		if !b524SystemInformationBytes.Equal(b524SystemInformationFirst, b524SystemInformationCase.b524SystemInformationExpected) {
			b524SystemInformationT.Fatalf("ab42 request: %x", b524SystemInformationFirst)
		}
		b524SystemInformationFirst[0] = 99
		if !b524SystemInformationBytes.Equal(B524SystemBuildInformation(b524SystemInformationCase.b524SystemInformationIdentifier), b524SystemInformationCase.b524SystemInformationExpected) {
			b524SystemInformationT.Fatal("ab42 request aliases prior result")
		}
	}
}

func TestAb42SyntheticInformation(b524SystemInformationT *b524SystemInformationTesting.T) {
	b524SystemInformationData, b524SystemInformationErr := b524SystemInformationOs.ReadFile("testdata/ab42_b524_system_information.json")
	if b524SystemInformationErr != nil {
		b524SystemInformationT.Fatal(b524SystemInformationErr)
	}
	var b524SystemInformationCases []struct {
		B524SystemInformationID            string  `json:"Ab42ID"`
		B524SystemInformationIdentifier    uint16  `json:"Ab42Identifier"`
		B524SystemInformationRaw           []byte  `json:"Ab42Raw"`
		B524SystemInformationValue         float32 `json:"Ab42Value"`
		B524SystemInformationCapacity      bool    `json:"Ab42Capacity"`
		B524SystemInformationHeatPumpCount bool    `json:"Ab42HeatPumpCount"`
		B524SystemInformationMeaning       string  `json:"Ab42Meaning"`
	}
	b524SystemInformationDecoder := b524SystemInformationJSON.NewDecoder(b524SystemInformationBytes.NewReader(b524SystemInformationData))
	b524SystemInformationDecoder.DisallowUnknownFields()
	if b524SystemInformationErr = b524SystemInformationDecoder.Decode(&b524SystemInformationCases); b524SystemInformationErr != nil {
		b524SystemInformationT.Fatal(b524SystemInformationErr)
	}
	if len(b524SystemInformationCases) != 7 {
		b524SystemInformationT.Fatalf("ab42 fixture count: %d", len(b524SystemInformationCases))
	}
	b524SystemInformationSeen := map[string]bool{}
	for _, b524SystemInformationCase := range b524SystemInformationCases {
		if b524SystemInformationCase.B524SystemInformationID == "" || b524SystemInformationSeen[b524SystemInformationCase.B524SystemInformationID] {
			b524SystemInformationT.Fatal("ab42 missing or duplicate fixture")
		}
		b524SystemInformationSeen[b524SystemInformationCase.B524SystemInformationID] = true
		b524SystemInformationT.Run(b524SystemInformationCase.B524SystemInformationID, func(b524SystemInformationT *b524SystemInformationTesting.T) {
			b524SystemInformationResult := B524SystemDecodeInformation(B524SystemInformationInformationContext{B524SystemInformationIdentifier: b524SystemInformationCase.B524SystemInformationIdentifier}, b524SystemInformationCase.B524SystemInformationRaw)
			if !b524SystemInformationResult.B524SystemInformationValid || b524SystemInformationResult.B524SystemInformationValue != b524SystemInformationCase.B524SystemInformationValue || b524SystemInformationResult.B524SystemInformationInvalidReason != "" {
				b524SystemInformationT.Fatalf("ab42 decode: %+v", b524SystemInformationResult)
			}
			if b524SystemInformationResult.B524SystemInformationContext.B524SystemInformationIdentifier != b524SystemInformationCase.B524SystemInformationIdentifier || !b524SystemInformationBytes.Equal(b524SystemInformationResult.B524SystemInformationRaw, b524SystemInformationCase.B524SystemInformationRaw) {
				b524SystemInformationT.Fatal("ab42 evidence lost")
			}
			if b524SystemInformationResult.B524SystemInformationMeaning != b524SystemInformationCase.B524SystemInformationMeaning {
				b524SystemInformationT.Fatalf("ab42 meaning: %s", b524SystemInformationResult.B524SystemInformationMeaning)
			}
			if (b524SystemInformationResult.B524SystemInformationCapacity != nil) != b524SystemInformationCase.B524SystemInformationCapacity || (b524SystemInformationResult.B524SystemInformationHeatPumpCount != nil) != b524SystemInformationCase.B524SystemInformationHeatPumpCount {
				b524SystemInformationT.Fatal("ab42 projection state")
			}
			if b524SystemInformationResult.B524SystemInformationCapacity != nil && *b524SystemInformationResult.B524SystemInformationCapacity != b524SystemInformationCase.B524SystemInformationValue {
				b524SystemInformationT.Fatal("ab42 capacity value")
			}
			if b524SystemInformationResult.B524SystemInformationHeatPumpCount != nil && *b524SystemInformationResult.B524SystemInformationHeatPumpCount != b524SystemInformationCase.B524SystemInformationValue {
				b524SystemInformationT.Fatal("ab42 count value")
			}
			if b524SystemInformationResult.B524SystemInformationRepresentationGrade != "Proven" || b524SystemInformationResult.B524SystemInformationLayoutGrade != "Hypothesis" || b524SystemInformationResult.B524SystemInformationQualificationGrade != "Unknown" || b524SystemInformationResult.B524SystemInformationPresenceGrade != "Unknown" || b524SystemInformationResult.B524SystemInformationEligibilityGrade != "Unknown" {
				b524SystemInformationT.Fatal("ab42 representation promoted to physical evidence")
			}
		})
	}
}

func TestAb42IndependentRepresentation(b524SystemInformationT *b524SystemInformationTesting.T) {
	for _, b524SystemInformationCase := range []struct {
		b524SystemInformationRaw      []byte
		b524SystemInformationValue    float32
		b524SystemInformationCapacity bool
	}{
		{[]byte{0, 0, 128, 63}, 1, true},
		{[]byte{0, 0, 192, 63}, 1.5, false},
		{[]byte{0, 0, 128, 191}, -1, false},
		{[]byte{0, 0, 0, 0}, 0, true},
		{[]byte{0, 0, 0, 128}, 0, true},
	} {
		for _, b524SystemInformationIdentifier := range []uint16{0, 1, 0x000d, 0xffff, 0x1234} {
			b524SystemInformationResult := B524SystemDecodeInformation(B524SystemInformationInformationContext{B524SystemInformationIdentifier: b524SystemInformationIdentifier}, b524SystemInformationCase.b524SystemInformationRaw)
			if !b524SystemInformationResult.B524SystemInformationValid || b524SystemInformationResult.B524SystemInformationValue != b524SystemInformationCase.b524SystemInformationValue {
				b524SystemInformationT.Fatal("ab42 independent numeric expectation")
			}
			if (b524SystemInformationResult.B524SystemInformationCapacity != nil) != (b524SystemInformationCase.b524SystemInformationCapacity && b524SystemInformationIdentifier <= 1) {
				b524SystemInformationT.Fatal("ab42 capacity is not an integral supported capacity")
			}
			if (b524SystemInformationResult.B524SystemInformationHeatPumpCount != nil) != (b524SystemInformationCase.b524SystemInformationCapacity && b524SystemInformationIdentifier == 0x000d) {
				b524SystemInformationT.Fatal("ab42 count identity conflated")
			}
			b524SystemInformationMeaning := "Unknown"
			switch b524SystemInformationIdentifier {
			case 0:
				b524SystemInformationMeaning = "supported circuit capacity"
			case 1:
				b524SystemInformationMeaning = "supported zone capacity"
			case 0x000d:
				b524SystemInformationMeaning = "heat-pump count"
			}
			if b524SystemInformationResult.B524SystemInformationMeaning != b524SystemInformationMeaning || b524SystemInformationResult.B524SystemInformationPresenceGrade != "Unknown" || b524SystemInformationResult.B524SystemInformationEligibilityGrade != "Unknown" {
				b524SystemInformationT.Fatal("ab42 capacity/count became presence or eligibility")
			}
		}
	}
}

func TestAb42InvalidInformation(b524SystemInformationT *b524SystemInformationTesting.T) {
	b524SystemInformationInputs := [][]byte{nil, {}, {0}, {0, 0}, {0, 0, 128}, {0, 0, 128, 63, 0}, {0, 0, 128, 63, 0, 0, 128, 63}}
	for _, b524SystemInformationBits := range []uint32{0x7fc00001, 0x7f800001, 0x7f800000, 0xff800000} {
		b524SystemInformationRaw := make([]byte, 4)
		b524SystemInformationBinary.LittleEndian.PutUint32(b524SystemInformationRaw, b524SystemInformationBits)
		b524SystemInformationInputs = append(b524SystemInformationInputs, b524SystemInformationRaw)
	}
	for _, b524SystemInformationRaw := range b524SystemInformationInputs {
		for _, b524SystemInformationIdentifier := range []uint16{0, 1, 0x000d, 0xffff} {
			b524SystemInformationContext := B524SystemInformationInformationContext{B524SystemInformationIdentifier: b524SystemInformationIdentifier, B524SystemInformationAPIVersion: b524SystemInformationTestPointer(uint16(2)), B524SystemInformationAPIRevision: b524SystemInformationTestPointer(uint16(7))}
			b524SystemInformationResult := B524SystemDecodeInformation(b524SystemInformationContext, b524SystemInformationRaw)
			if b524SystemInformationResult.B524SystemInformationValid || b524SystemInformationResult.B524SystemInformationInvalidReason == "" || b524SystemInformationResult.B524SystemInformationCapacity != nil || b524SystemInformationResult.B524SystemInformationHeatPumpCount != nil || b524SystemInformationResult.B524SystemInformationRepresentationGrade != "Unknown" {
				b524SystemInformationT.Fatalf("ab42 invalid input admitted: %x", b524SystemInformationRaw)
			}
			if !b524SystemInformationBytes.Equal(b524SystemInformationResult.B524SystemInformationRaw, b524SystemInformationRaw) || !b524SystemInformationReflect.DeepEqual(b524SystemInformationResult.B524SystemInformationContext, b524SystemInformationContext) {
				b524SystemInformationT.Fatal("ab42 invalid response discarded independent evidence")
			}
			if len(b524SystemInformationRaw) > 0 {
				b524SystemInformationExpected := append([]byte(nil), b524SystemInformationRaw...)
				b524SystemInformationRaw[0] ^= 0xff
				if !b524SystemInformationBytes.Equal(b524SystemInformationResult.B524SystemInformationRaw, b524SystemInformationExpected) {
					b524SystemInformationT.Fatal("ab42 invalid raw aliases caller")
				}
				b524SystemInformationRaw[0] ^= 0xff
				b524SystemInformationResult.B524SystemInformationRaw[0] ^= 0xff
				if !b524SystemInformationBytes.Equal(b524SystemInformationRaw, b524SystemInformationExpected) {
					b524SystemInformationT.Fatal("ab42 caller raw aliases invalid result")
				}
			}
		}
	}
}

func b524SystemInformationTestContext() B524SystemInformationInformationContext {
	return B524SystemInformationInformationContext{
		B524SystemInformationFamily: b524SystemInformationTestPointer("ab42-family"), B524SystemInformationSoftware: b524SystemInformationTestPointer("ab42-software"), B524SystemInformationHardware: b524SystemInformationTestPointer("ab42-hardware"),
		B524SystemInformationAPIVersion: b524SystemInformationTestPointer(uint16(2)), B524SystemInformationAPIRevision: b524SystemInformationTestPointer(uint16(7)), B524SystemInformationIdentifier: 0x1234,
	}
}

func TestAb42ContextAxes(b524SystemInformationT *b524SystemInformationTesting.T) {
	for _, b524SystemInformationAxis := range []string{"B524SystemInformationFamily", "B524SystemInformationSoftware", "B524SystemInformationHardware", "B524SystemInformationAPIVersion", "B524SystemInformationAPIRevision", "B524SystemInformationIdentifier"} {
		b524SystemInformationT.Run(b524SystemInformationAxis, func(b524SystemInformationT *b524SystemInformationTesting.T) {
			for _, b524SystemInformationAlternative := range []string{"ab42-absent", "ab42-zero", "ab42-alternative"} {
				b524SystemInformationContext := b524SystemInformationTestContext()
				b524SystemInformationField := b524SystemInformationReflect.ValueOf(&b524SystemInformationContext).Elem().FieldByName(b524SystemInformationAxis)
				if b524SystemInformationField.Kind() == b524SystemInformationReflect.Uint16 {
					if b524SystemInformationAlternative == "ab42-alternative" {
						b524SystemInformationField.SetUint(0xffff)
					} else {
						b524SystemInformationField.SetUint(0)
					}
				} else if b524SystemInformationAlternative == "ab42-absent" {
					b524SystemInformationField.SetZero()
				} else {
					b524SystemInformationField.Set(b524SystemInformationReflect.New(b524SystemInformationField.Type().Elem()))
					if b524SystemInformationAlternative == "ab42-alternative" {
						if b524SystemInformationField.Elem().Kind() == b524SystemInformationReflect.String {
							b524SystemInformationField.Elem().SetString("ab42-alternative")
						} else {
							b524SystemInformationField.Elem().SetUint(19)
						}
					}
				}
				b524SystemInformationResult := B524SystemDecodeInformation(b524SystemInformationContext, []byte{0, 0, 128, 63})
				if !b524SystemInformationResult.B524SystemInformationValid || !b524SystemInformationReflect.DeepEqual(b524SystemInformationResult.B524SystemInformationContext, b524SystemInformationContext) {
					b524SystemInformationT.Fatalf("ab42 context axis %s %s lost", b524SystemInformationAxis, b524SystemInformationAlternative)
				}
				b524SystemInformationInvalid := B524SystemDecodeInformation(b524SystemInformationContext, []byte{0})
				if b524SystemInformationInvalid.B524SystemInformationValid || !b524SystemInformationReflect.DeepEqual(b524SystemInformationInvalid.B524SystemInformationContext, b524SystemInformationContext) {
					b524SystemInformationT.Fatal("ab42 partial failure lost context")
				}
			}
		})
	}
}

func TestAb42Ownership(b524SystemInformationT *b524SystemInformationTesting.T) {
	b524SystemInformationContext := b524SystemInformationTestContext()
	b524SystemInformationRaw := []byte{0, 0, 128, 63}
	b524SystemInformationFirst := B524SystemDecodeInformation(b524SystemInformationContext, b524SystemInformationRaw)
	b524SystemInformationSecond := B524SystemDecodeInformation(b524SystemInformationContext, b524SystemInformationRaw)
	for _, b524SystemInformationAxis := range []string{"B524SystemInformationFamily", "B524SystemInformationSoftware", "B524SystemInformationHardware", "B524SystemInformationAPIVersion", "B524SystemInformationAPIRevision"} {
		b524SystemInformationSource := b524SystemInformationReflect.ValueOf(&b524SystemInformationContext).Elem().FieldByName(b524SystemInformationAxis).Elem()
		b524SystemInformationSaved := b524SystemInformationReflect.ValueOf(b524SystemInformationFirst.B524SystemInformationContext).FieldByName(b524SystemInformationAxis).Elem().Interface()
		if b524SystemInformationSource.Kind() == b524SystemInformationReflect.String {
			b524SystemInformationSource.SetString("ab42-mutated")
		} else {
			b524SystemInformationSource.SetUint(99)
		}
		if !b524SystemInformationReflect.DeepEqual(b524SystemInformationReflect.ValueOf(b524SystemInformationFirst.B524SystemInformationContext).FieldByName(b524SystemInformationAxis).Elem().Interface(), b524SystemInformationSaved) {
			b524SystemInformationT.Fatalf("ab42 caller context aliases %s", b524SystemInformationAxis)
		}
		b524SystemInformationReturned := b524SystemInformationReflect.ValueOf(b524SystemInformationFirst.B524SystemInformationContext).FieldByName(b524SystemInformationAxis).Elem()
		if b524SystemInformationReturned.Kind() == b524SystemInformationReflect.String {
			b524SystemInformationReturned.SetString("ab42-result")
		} else {
			b524SystemInformationReturned.SetUint(100)
		}
		if !b524SystemInformationReflect.DeepEqual(b524SystemInformationReflect.ValueOf(b524SystemInformationSecond.B524SystemInformationContext).FieldByName(b524SystemInformationAxis).Elem().Interface(), b524SystemInformationSaved) {
			b524SystemInformationT.Fatalf("ab42 returned contexts alias %s", b524SystemInformationAxis)
		}
		if b524SystemInformationReflect.DeepEqual(b524SystemInformationSource.Interface(), b524SystemInformationReturned.Interface()) {
			b524SystemInformationT.Fatal("ab42 returned context aliases source")
		}
	}
	b524SystemInformationRaw[0] = 1
	if !b524SystemInformationBytes.Equal(b524SystemInformationFirst.B524SystemInformationRaw, []byte{0, 0, 128, 63}) {
		b524SystemInformationT.Fatal("ab42 caller raw aliases result")
	}
	b524SystemInformationFirst.B524SystemInformationRaw[1] = 1
	if !b524SystemInformationBytes.Equal(b524SystemInformationSecond.B524SystemInformationRaw, []byte{0, 0, 128, 63}) || b524SystemInformationRaw[1] != 0 {
		b524SystemInformationT.Fatal("ab42 returned raw aliases caller or later result")
	}
	for _, b524SystemInformationIdentifier := range []uint16{0, 0x000d} {
		b524SystemInformationOne := B524SystemDecodeInformation(B524SystemInformationInformationContext{B524SystemInformationIdentifier: b524SystemInformationIdentifier}, []byte{0, 0, 128, 63})
		b524SystemInformationTwo := B524SystemDecodeInformation(B524SystemInformationInformationContext{B524SystemInformationIdentifier: b524SystemInformationIdentifier}, []byte{0, 0, 128, 63})
		b524SystemInformationProjection, b524SystemInformationOther := b524SystemInformationOne.B524SystemInformationCapacity, b524SystemInformationTwo.B524SystemInformationCapacity
		if b524SystemInformationIdentifier == 0x000d {
			b524SystemInformationProjection, b524SystemInformationOther = b524SystemInformationOne.B524SystemInformationHeatPumpCount, b524SystemInformationTwo.B524SystemInformationHeatPumpCount
		}
		*b524SystemInformationProjection = 9
		if *b524SystemInformationOther != 1 || b524SystemInformationOne.B524SystemInformationValue != 1 {
			b524SystemInformationT.Fatal("ab42 projection aliases independent state")
		}
	}
	b524SystemInformationLater := B524SystemDecodeInformation(b524SystemInformationTestContext(), []byte{0, 0, 128, 63})
	if !b524SystemInformationReflect.DeepEqual(b524SystemInformationLater.B524SystemInformationContext, b524SystemInformationTestContext()) || b524SystemInformationMath.Float32bits(b524SystemInformationLater.B524SystemInformationValue) != 0x3f800000 {
		b524SystemInformationT.Fatal("ab42 mutations changed later decode")
	}
}
