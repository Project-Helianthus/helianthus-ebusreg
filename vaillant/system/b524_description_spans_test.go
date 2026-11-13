package system

import (
	b524DescriptionSpanBinary "encoding/binary"
	b524DescriptionSpanJSON "encoding/json"
	b524DescriptionSpanMath "math"
	b524DescriptionSpanOS "os"
	b524DescriptionSpanReflect "reflect"
	b524DescriptionSpanTesting "testing"
)

func b524DescriptionSpanPtr[b524DescriptionSpanT any](b524DescriptionSpanValue b524DescriptionSpanT) *b524DescriptionSpanT {
	return &b524DescriptionSpanValue
}

func b524DescriptionSpanTestContext() B524DescriptionSpanContext {
	return B524DescriptionSpanContext{
		B524DescriptionSpanFamily: b524DescriptionSpanPtr("ab44-family-a"), B524DescriptionSpanSoftware: b524DescriptionSpanPtr(uint16(1)), B524DescriptionSpanHardware: b524DescriptionSpanPtr(uint16(2)),
		B524DescriptionSpanAPIVersion: b524DescriptionSpanPtr(byte(1)), B524DescriptionSpanAPIRevision: b524DescriptionSpanPtr(byte(0)),
		B524DescriptionSpanOP: b524DescriptionSpanPtr(byte(1)), B524DescriptionSpanGG: b524DescriptionSpanPtr(byte(2)), B524DescriptionSpanII: b524DescriptionSpanPtr(byte(1)), B524DescriptionSpanRR: b524DescriptionSpanPtr(uint16(2)),
	}
}

func b524DescriptionSpanTestEXP(b524DescriptionSpanMinimum, b524DescriptionSpanMaximum, b524DescriptionSpanStep float32) []byte {
	b524DescriptionSpanRaw := []byte{2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	b524DescriptionSpanBinary.LittleEndian.PutUint32(b524DescriptionSpanRaw[3:7], b524DescriptionSpanMath.Float32bits(b524DescriptionSpanMinimum))
	b524DescriptionSpanBinary.LittleEndian.PutUint32(b524DescriptionSpanRaw[7:11], b524DescriptionSpanMath.Float32bits(b524DescriptionSpanMaximum))
	b524DescriptionSpanBinary.LittleEndian.PutUint32(b524DescriptionSpanRaw[11:15], b524DescriptionSpanMath.Float32bits(b524DescriptionSpanStep))
	return b524DescriptionSpanRaw
}

func TestAb44Fixtures(b524DescriptionSpanT *b524DescriptionSpanTesting.T) {
	b524DescriptionSpanData, b524DescriptionSpanErr := b524DescriptionSpanOS.ReadFile("testdata/ab44_b524_description_spans.json")
	if b524DescriptionSpanErr != nil {
		b524DescriptionSpanT.Fatal(b524DescriptionSpanErr)
	}
	var b524DescriptionSpanFixtures []struct {
		B524DescriptionSpanID      string                   `json:"Ab44ID"`
		B524DescriptionSpanCodec   B524DescriptionSpanCodec `json:"Ab44Codec"`
		B524DescriptionSpanWidth   int                      `json:"Ab44Width"`
		B524DescriptionSpanRR      uint16                   `json:"Ab44RR"`
		B524DescriptionSpanRaw     []byte                   `json:"Ab44Raw"`
		B524DescriptionSpanMinimum float64                  `json:"Ab44Minimum"`
		B524DescriptionSpanMaximum float64                  `json:"Ab44Maximum"`
		B524DescriptionSpanStep    float64                  `json:"Ab44Step"`
		B524DescriptionSpanNumeric bool                     `json:"Ab44Numeric"`
	}
	if b524DescriptionSpanErr = b524DescriptionSpanJSON.Unmarshal(b524DescriptionSpanData, &b524DescriptionSpanFixtures); b524DescriptionSpanErr != nil {
		b524DescriptionSpanT.Fatal(b524DescriptionSpanErr)
	}
	if len(b524DescriptionSpanFixtures) != 6 {
		b524DescriptionSpanT.Fatal("ab44 fixture count", len(b524DescriptionSpanFixtures))
	}
	b524DescriptionSpanSeen := map[string]bool{}
	for _, b524DescriptionSpanFixture := range b524DescriptionSpanFixtures {
		if b524DescriptionSpanSeen[b524DescriptionSpanFixture.B524DescriptionSpanID] {
			b524DescriptionSpanT.Fatal("ab44 duplicate fixture")
		}
		b524DescriptionSpanSeen[b524DescriptionSpanFixture.B524DescriptionSpanID] = true
		b524DescriptionSpanT.Run(b524DescriptionSpanFixture.B524DescriptionSpanID, func(b524DescriptionSpanT *b524DescriptionSpanTesting.T) {
			b524DescriptionSpanContext := b524DescriptionSpanTestContext()
			b524DescriptionSpanContext.B524DescriptionSpanRR = b524DescriptionSpanPtr(b524DescriptionSpanFixture.B524DescriptionSpanRR)
			for _, b524DescriptionSpanGrade := range []B524DescriptionSpanGrade{B524DescriptionSpanHypothesis, B524DescriptionSpanProven} {
				b524DescriptionSpanResult := B524DecodeDescription(b524DescriptionSpanContext, b524DescriptionSpanFixture.B524DescriptionSpanCodec, b524DescriptionSpanFixture.B524DescriptionSpanWidth, b524DescriptionSpanGrade, b524DescriptionSpanFixture.B524DescriptionSpanRaw)
				if !b524DescriptionSpanResult.B524DescriptionSpanImplementationValid || b524DescriptionSpanResult.B524DescriptionSpanNumericValid != b524DescriptionSpanFixture.B524DescriptionSpanNumeric || b524DescriptionSpanResult.B524DescriptionSpanGrade != b524DescriptionSpanGrade || !b524DescriptionSpanReflect.DeepEqual(b524DescriptionSpanResult.B524DescriptionSpanRaw, b524DescriptionSpanFixture.B524DescriptionSpanRaw) || !b524DescriptionSpanReflect.DeepEqual(b524DescriptionSpanResult.B524DescriptionSpanContext, b524DescriptionSpanContext) || b524DescriptionSpanResult.B524DescriptionSpanMetadataOP == nil || *b524DescriptionSpanResult.B524DescriptionSpanMetadataOP != 2 {
					b524DescriptionSpanT.Fatalf("ab44 unexpected result: %+v", b524DescriptionSpanResult)
				}
				if !b524DescriptionSpanReflect.DeepEqual(b524DescriptionSpanResult.B524DescriptionSpanMinimumRaw, b524DescriptionSpanFixture.B524DescriptionSpanRaw[3:3+b524DescriptionSpanFixture.B524DescriptionSpanWidth]) || !b524DescriptionSpanReflect.DeepEqual(b524DescriptionSpanResult.B524DescriptionSpanMaximumRaw, b524DescriptionSpanFixture.B524DescriptionSpanRaw[3+b524DescriptionSpanFixture.B524DescriptionSpanWidth:3+2*b524DescriptionSpanFixture.B524DescriptionSpanWidth]) || !b524DescriptionSpanReflect.DeepEqual(b524DescriptionSpanResult.B524DescriptionSpanStepRaw, b524DescriptionSpanFixture.B524DescriptionSpanRaw[3+2*b524DescriptionSpanFixture.B524DescriptionSpanWidth:]) {
					b524DescriptionSpanT.Fatal("ab44 spans")
				}
				if b524DescriptionSpanFixture.B524DescriptionSpanNumeric {
					b524DescriptionSpanExpected := B524DescriptionSpanLimits{b524DescriptionSpanFixture.B524DescriptionSpanMinimum, b524DescriptionSpanFixture.B524DescriptionSpanMaximum, b524DescriptionSpanFixture.B524DescriptionSpanStep}
					if b524DescriptionSpanResult.B524DescriptionSpanLimits == nil || *b524DescriptionSpanResult.B524DescriptionSpanLimits != b524DescriptionSpanExpected {
						b524DescriptionSpanT.Fatalf("ab44 limits: %+v", b524DescriptionSpanResult.B524DescriptionSpanLimits)
					}
				} else if b524DescriptionSpanResult.B524DescriptionSpanLimits != nil || b524DescriptionSpanResult.B524DescriptionSpanStepMeaning != B524DescriptionSpanUnknown || b524DescriptionSpanResult.B524DescriptionSpanReason != "date_time_raw_only" {
					b524DescriptionSpanT.Fatal("ab44 date/time numeric permission")
				}
			}
		})
	}
}

func TestAb44Build(b524DescriptionSpanT *b524DescriptionSpanTesting.T) {
	b524DescriptionSpanContext := b524DescriptionSpanTestContext()
	b524DescriptionSpanRequest, b524DescriptionSpanErr := B524BuildDescription(b524DescriptionSpanContext)
	if b524DescriptionSpanErr != nil || !b524DescriptionSpanReflect.DeepEqual(b524DescriptionSpanRequest, []byte{1, 2, 1, 2, 0}) {
		b524DescriptionSpanT.Fatal("ab44 request", b524DescriptionSpanRequest, b524DescriptionSpanErr)
	}
	*b524DescriptionSpanContext.B524DescriptionSpanOP = 7
	*b524DescriptionSpanContext.B524DescriptionSpanRR = 0x1234
	b524DescriptionSpanRequest, b524DescriptionSpanErr = B524BuildDescription(b524DescriptionSpanContext)
	if b524DescriptionSpanErr != nil || !b524DescriptionSpanReflect.DeepEqual(b524DescriptionSpanRequest, []byte{7, 2, 1, 0x34, 0x12}) {
		b524DescriptionSpanT.Fatal("ab44 remote high register", b524DescriptionSpanRequest, b524DescriptionSpanErr)
	}
	for b524DescriptionSpanOP := 0; b524DescriptionSpanOP < 256; b524DescriptionSpanOP++ {
		*b524DescriptionSpanContext.B524DescriptionSpanOP = byte(b524DescriptionSpanOP)
		b524DescriptionSpanRequest, b524DescriptionSpanErr = B524BuildDescription(b524DescriptionSpanContext)
		if b524DescriptionSpanOP == 1 || b524DescriptionSpanOP == 7 {
			continue
		}
		if b524DescriptionSpanErr == nil || b524DescriptionSpanRequest != nil {
			b524DescriptionSpanT.Fatal("ab44 unsupported op", b524DescriptionSpanOP)
		}
	}
}

func TestAb44ContextAxes(b524DescriptionSpanT *b524DescriptionSpanTesting.T) {
	b524DescriptionSpanAxes := []string{"ab44Family", "ab44Software", "ab44Hardware", "ab44APIVersion", "ab44APIRevision", "ab44OP", "ab44GG", "ab44II", "ab44RR"}
	for b524DescriptionSpanIndex, b524DescriptionSpanAxis := range b524DescriptionSpanAxes {
		b524DescriptionSpanT.Run(b524DescriptionSpanAxis, func(b524DescriptionSpanT *b524DescriptionSpanTesting.T) {
			for b524DescriptionSpanVariant := 0; b524DescriptionSpanVariant < 3; b524DescriptionSpanVariant++ {
				b524DescriptionSpanContext := b524DescriptionSpanTestContext()
				b524DescriptionSpanField := b524DescriptionSpanReflect.ValueOf(&b524DescriptionSpanContext).Elem().Field(b524DescriptionSpanIndex)
				if b524DescriptionSpanVariant == 0 {
					b524DescriptionSpanField.SetZero()
				} else if b524DescriptionSpanIndex == 0 {
					b524DescriptionSpanAlternative := "ab44-family-b"
					if b524DescriptionSpanVariant == 1 {
						b524DescriptionSpanAlternative = ""
					}
					b524DescriptionSpanField.Set(b524DescriptionSpanReflect.ValueOf(b524DescriptionSpanPtr(b524DescriptionSpanAlternative)))
				} else {
					b524DescriptionSpanAlternative := uint64(9)
					if b524DescriptionSpanVariant == 1 {
						b524DescriptionSpanAlternative = 0
					}
					if b524DescriptionSpanIndex == 5 && b524DescriptionSpanVariant == 2 {
						b524DescriptionSpanAlternative = 7
					}
					b524DescriptionSpanField.Elem().SetUint(b524DescriptionSpanAlternative)
				}
				b524DescriptionSpanRaw := b524DescriptionSpanTestEXP(-20, 20, 1)
				if b524DescriptionSpanContext.B524DescriptionSpanGG != nil {
					b524DescriptionSpanRaw[0] = *b524DescriptionSpanContext.B524DescriptionSpanGG
				}
				if b524DescriptionSpanContext.B524DescriptionSpanRR != nil {
					b524DescriptionSpanBinary.LittleEndian.PutUint16(b524DescriptionSpanRaw[1:3], *b524DescriptionSpanContext.B524DescriptionSpanRR)
				}
				b524DescriptionSpanExpectedValid := b524DescriptionSpanVariant != 0 && (b524DescriptionSpanIndex != 0 || b524DescriptionSpanVariant != 1) && (b524DescriptionSpanIndex != 5 || b524DescriptionSpanVariant != 1)
				b524DescriptionSpanResult := B524DecodeDescription(b524DescriptionSpanContext, B524DescriptionSpanEXP, 4, B524DescriptionSpanHypothesis, b524DescriptionSpanRaw)
				if !b524DescriptionSpanReflect.DeepEqual(b524DescriptionSpanResult.B524DescriptionSpanContext, b524DescriptionSpanContext) || b524DescriptionSpanResult.B524DescriptionSpanNumericValid != b524DescriptionSpanExpectedValid {
					b524DescriptionSpanT.Fatalf("ab44 context variant %d: %+v", b524DescriptionSpanVariant, b524DescriptionSpanResult)
				}
				b524DescriptionSpanRequest, b524DescriptionSpanErr := B524BuildDescription(b524DescriptionSpanContext)
				if (b524DescriptionSpanErr == nil) != b524DescriptionSpanExpectedValid {
					b524DescriptionSpanT.Fatal("ab44 build context", b524DescriptionSpanErr)
				}
				if b524DescriptionSpanExpectedValid {
					if !b524DescriptionSpanReflect.DeepEqual(b524DescriptionSpanRequest, []byte{*b524DescriptionSpanContext.B524DescriptionSpanOP, *b524DescriptionSpanContext.B524DescriptionSpanGG, *b524DescriptionSpanContext.B524DescriptionSpanII, byte(*b524DescriptionSpanContext.B524DescriptionSpanRR), byte(*b524DescriptionSpanContext.B524DescriptionSpanRR >> 8)}) {
						b524DescriptionSpanT.Fatal("ab44 build identity")
					}
					b524DescriptionSpanExpectedOP := byte(2)
					if *b524DescriptionSpanContext.B524DescriptionSpanOP == 7 {
						b524DescriptionSpanExpectedOP = 6
					}
					if b524DescriptionSpanResult.B524DescriptionSpanMetadataOP == nil || *b524DescriptionSpanResult.B524DescriptionSpanMetadataOP != b524DescriptionSpanExpectedOP {
						b524DescriptionSpanT.Fatal("ab44 direction")
					}
				} else if b524DescriptionSpanResult.B524DescriptionSpanLimits != nil || b524DescriptionSpanResult.B524DescriptionSpanImplementationValid || b524DescriptionSpanResult.B524DescriptionSpanReason == "" {
					b524DescriptionSpanT.Fatal("ab44 incomplete context admitted")
				}
			}
		})
	}
}

func TestAb44RejectRawOnly(b524DescriptionSpanT *b524DescriptionSpanTesting.T) {
	b524DescriptionSpanRaw := []byte{2, 2, 0, 0, 0, 4, 0, 1, 0}
	b524DescriptionSpanCases := []struct {
		B524DescriptionSpanID     string
		B524DescriptionSpanCodec  B524DescriptionSpanCodec
		B524DescriptionSpanWidth  int
		B524DescriptionSpanGrade  B524DescriptionSpanGrade
		B524DescriptionSpanReason string
	}{
		{"ab44_absent_codec", "", 2, B524DescriptionSpanProven, "unsupported_codec"},
		{"ab44_unknown_codec", "ab44_future", 2, B524DescriptionSpanProven, "unsupported_codec"},
		{"ab44_wrong_codec", B524DescriptionSpanEXP, 2, B524DescriptionSpanProven, "codec_width_mismatch"},
		{"ab44_uch_width", B524DescriptionSpanUCH, 2, B524DescriptionSpanProven, "codec_width_mismatch"},
		{"ab44_wrong_width", B524DescriptionSpanUIN, 1, B524DescriptionSpanProven, "malformed_length"},
		{"ab44_absent_width", B524DescriptionSpanUIN, 0, B524DescriptionSpanProven, "bad_width"},
		{"ab44_negative_width", B524DescriptionSpanUIN, -1, B524DescriptionSpanProven, "bad_width"},
		{"ab44_large_width", B524DescriptionSpanUIN, int(^uint(0) >> 1), B524DescriptionSpanProven, "bad_width"},
		{"ab44_absent_grade", B524DescriptionSpanUIN, 2, "", "codec_grade_required"},
		{"ab44_unknown_grade", B524DescriptionSpanUIN, 2, B524DescriptionSpanUnknown, "codec_grade_required"},
		{"ab44_other_grade", B524DescriptionSpanUIN, 2, "ab44_other", "codec_grade_required"},
	}
	for _, b524DescriptionSpanCase := range b524DescriptionSpanCases {
		b524DescriptionSpanT.Run(b524DescriptionSpanCase.B524DescriptionSpanID, func(b524DescriptionSpanT *b524DescriptionSpanTesting.T) {
			b524DescriptionSpanResult := B524DecodeDescription(b524DescriptionSpanTestContext(), b524DescriptionSpanCase.B524DescriptionSpanCodec, b524DescriptionSpanCase.B524DescriptionSpanWidth, b524DescriptionSpanCase.B524DescriptionSpanGrade, b524DescriptionSpanRaw)
			if b524DescriptionSpanResult.B524DescriptionSpanImplementationValid || b524DescriptionSpanResult.B524DescriptionSpanNumericValid || b524DescriptionSpanResult.B524DescriptionSpanLimits != nil || b524DescriptionSpanResult.B524DescriptionSpanReason != b524DescriptionSpanCase.B524DescriptionSpanReason || b524DescriptionSpanResult.B524DescriptionSpanCodec != b524DescriptionSpanCase.B524DescriptionSpanCodec || b524DescriptionSpanResult.B524DescriptionSpanWidth != b524DescriptionSpanCase.B524DescriptionSpanWidth || b524DescriptionSpanResult.B524DescriptionSpanGrade != b524DescriptionSpanCase.B524DescriptionSpanGrade || !b524DescriptionSpanReflect.DeepEqual(b524DescriptionSpanResult.B524DescriptionSpanRaw, b524DescriptionSpanRaw) {
				b524DescriptionSpanT.Fatalf("ab44 rejection: %+v", b524DescriptionSpanResult)
			}
			if b524DescriptionSpanCase.B524DescriptionSpanWidth == 2 && (!b524DescriptionSpanReflect.DeepEqual(b524DescriptionSpanResult.B524DescriptionSpanMinimumRaw, []byte{0, 0}) || !b524DescriptionSpanReflect.DeepEqual(b524DescriptionSpanResult.B524DescriptionSpanMaximumRaw, []byte{4, 0}) || !b524DescriptionSpanReflect.DeepEqual(b524DescriptionSpanResult.B524DescriptionSpanStepRaw, []byte{1, 0})) {
				b524DescriptionSpanT.Fatal("ab44 rejected interpretation lost raw spans")
			}
		})
	}
	for _, b524DescriptionSpanCodec := range []B524DescriptionSpanCodec{B524DescriptionSpanUCH, B524DescriptionSpanUIN, B524DescriptionSpanEXP, B524DescriptionSpanDate, B524DescriptionSpanTime} {
		for b524DescriptionSpanWidth := 1; b524DescriptionSpanWidth <= 4; b524DescriptionSpanWidth++ {
			b524DescriptionSpanExpectedWidth := map[B524DescriptionSpanCodec]int{B524DescriptionSpanUCH: 1, B524DescriptionSpanUIN: 2, B524DescriptionSpanEXP: 4, B524DescriptionSpanDate: 3, B524DescriptionSpanTime: 3}[b524DescriptionSpanCodec]
			if b524DescriptionSpanWidth == b524DescriptionSpanExpectedWidth {
				continue
			}
			b524DescriptionSpanWrong := make([]byte, 3+3*b524DescriptionSpanWidth)
			copy(b524DescriptionSpanWrong, []byte{2, 2, 0})
			b524DescriptionSpanResult := B524DecodeDescription(b524DescriptionSpanTestContext(), b524DescriptionSpanCodec, b524DescriptionSpanWidth, B524DescriptionSpanProven, b524DescriptionSpanWrong)
			if b524DescriptionSpanResult.B524DescriptionSpanReason != "codec_width_mismatch" || b524DescriptionSpanResult.B524DescriptionSpanImplementationValid {
				b524DescriptionSpanT.Fatal("ab44 codec width inferred", b524DescriptionSpanCodec, b524DescriptionSpanWidth)
			}
		}
	}
}

func TestAb44MalformedAndEcho(b524DescriptionSpanT *b524DescriptionSpanTesting.T) {
	b524DescriptionSpanRaw := b524DescriptionSpanTestEXP(-20, 20, 1)
	for b524DescriptionSpanIndex := 0; b524DescriptionSpanIndex < 3; b524DescriptionSpanIndex++ {
		b524DescriptionSpanChanged := append([]byte(nil), b524DescriptionSpanRaw...)
		b524DescriptionSpanChanged[b524DescriptionSpanIndex] ^= 1
		b524DescriptionSpanResult := B524DecodeDescription(b524DescriptionSpanTestContext(), B524DescriptionSpanEXP, 4, B524DescriptionSpanProven, b524DescriptionSpanChanged)
		if b524DescriptionSpanResult.B524DescriptionSpanReason != "echo_mismatch" || b524DescriptionSpanResult.B524DescriptionSpanNumericValid || b524DescriptionSpanResult.B524DescriptionSpanLimits != nil || !b524DescriptionSpanReflect.DeepEqual(b524DescriptionSpanResult.B524DescriptionSpanRaw, b524DescriptionSpanChanged) || !b524DescriptionSpanReflect.DeepEqual(b524DescriptionSpanResult.B524DescriptionSpanStepRaw, b524DescriptionSpanRaw[11:]) {
			b524DescriptionSpanT.Fatal("ab44 echo", b524DescriptionSpanIndex, b524DescriptionSpanResult)
		}
	}
	for b524DescriptionSpanLength := 0; b524DescriptionSpanLength <= 30; b524DescriptionSpanLength++ {
		if b524DescriptionSpanLength == 15 {
			continue
		}
		b524DescriptionSpanMalformed := make([]byte, b524DescriptionSpanLength)
		copy(b524DescriptionSpanMalformed, b524DescriptionSpanRaw)
		b524DescriptionSpanResult := B524DecodeDescription(b524DescriptionSpanTestContext(), B524DescriptionSpanEXP, 4, B524DescriptionSpanProven, b524DescriptionSpanMalformed)
		if b524DescriptionSpanResult.B524DescriptionSpanReason != "malformed_length" || b524DescriptionSpanResult.B524DescriptionSpanImplementationValid || b524DescriptionSpanResult.B524DescriptionSpanLimits != nil || len(b524DescriptionSpanResult.B524DescriptionSpanRaw) != b524DescriptionSpanLength {
			b524DescriptionSpanT.Fatal("ab44 malformed", b524DescriptionSpanLength)
		}
	}
	b524DescriptionSpanDuplicate := append(append([]byte(nil), b524DescriptionSpanRaw...), b524DescriptionSpanRaw...)
	b524DescriptionSpanLegacy := []byte{1, 2, 2, 0}
	for _, b524DescriptionSpanMalformed := range [][]byte{nil, b524DescriptionSpanDuplicate, b524DescriptionSpanLegacy} {
		b524DescriptionSpanResult := B524DecodeDescription(b524DescriptionSpanTestContext(), B524DescriptionSpanEXP, 4, B524DescriptionSpanProven, b524DescriptionSpanMalformed)
		if b524DescriptionSpanResult.B524DescriptionSpanImplementationValid || b524DescriptionSpanResult.B524DescriptionSpanReason != "malformed_length" {
			b524DescriptionSpanT.Fatal("ab44 duplicate/legacy")
		}
	}
	for _, b524DescriptionSpanOP := range []byte{0, 2, 6, 8} {
		b524DescriptionSpanContext := b524DescriptionSpanTestContext()
		*b524DescriptionSpanContext.B524DescriptionSpanOP = b524DescriptionSpanOP
		b524DescriptionSpanResult := B524DecodeDescription(b524DescriptionSpanContext, B524DescriptionSpanEXP, 4, B524DescriptionSpanProven, b524DescriptionSpanRaw)
		if b524DescriptionSpanResult.B524DescriptionSpanReason != "unsupported_op" || b524DescriptionSpanResult.B524DescriptionSpanImplementationValid {
			b524DescriptionSpanT.Fatal("ab44 reply direction", b524DescriptionSpanOP)
		}
	}
}

func TestAb44NumericSpans(b524DescriptionSpanT *b524DescriptionSpanTesting.T) {
	for b524DescriptionSpanIndex, b524DescriptionSpanExpected := range []B524DescriptionSpanLimits{{-19, 20, 1}, {-20, 21, 1}, {-20, 20, 2}, {-20, 20, 0}} {
		b524DescriptionSpanRaw := b524DescriptionSpanTestEXP(float32(b524DescriptionSpanExpected.B524DescriptionSpanMinimum), float32(b524DescriptionSpanExpected.B524DescriptionSpanMaximum), float32(b524DescriptionSpanExpected.B524DescriptionSpanStep))
		b524DescriptionSpanResult := B524DecodeDescription(b524DescriptionSpanTestContext(), B524DescriptionSpanEXP, 4, B524DescriptionSpanHypothesis, b524DescriptionSpanRaw)
		if !b524DescriptionSpanResult.B524DescriptionSpanNumericValid || b524DescriptionSpanResult.B524DescriptionSpanLimits == nil || *b524DescriptionSpanResult.B524DescriptionSpanLimits != b524DescriptionSpanExpected {
			b524DescriptionSpanT.Fatal("ab44 independent span", b524DescriptionSpanIndex)
		}
	}
	b524DescriptionSpanCases := []struct {
		B524DescriptionSpanRaw    []byte
		B524DescriptionSpanReason string
	}{
		{b524DescriptionSpanTestEXP(-20, 20, -1), "negative_step"},
		{b524DescriptionSpanTestEXP(21, 20, 1), "reversed_limits"},
	}
	for b524DescriptionSpanIndex := 0; b524DescriptionSpanIndex < 3; b524DescriptionSpanIndex++ {
		for _, b524DescriptionSpanNonfinite := range []float32{float32(b524DescriptionSpanMath.NaN()), float32(b524DescriptionSpanMath.Inf(1)), float32(b524DescriptionSpanMath.Inf(-1))} {
			b524DescriptionSpanRaw := b524DescriptionSpanTestEXP(-20, 20, 1)
			b524DescriptionSpanBinary.LittleEndian.PutUint32(b524DescriptionSpanRaw[3+4*b524DescriptionSpanIndex:], b524DescriptionSpanMath.Float32bits(b524DescriptionSpanNonfinite))
			b524DescriptionSpanCases = append(b524DescriptionSpanCases, struct {
				B524DescriptionSpanRaw    []byte
				B524DescriptionSpanReason string
			}{b524DescriptionSpanRaw, "nonfinite_limits"})
		}
	}
	for _, b524DescriptionSpanCase := range b524DescriptionSpanCases {
		b524DescriptionSpanResult := B524DecodeDescription(b524DescriptionSpanTestContext(), B524DescriptionSpanEXP, 4, B524DescriptionSpanProven, b524DescriptionSpanCase.B524DescriptionSpanRaw)
		if b524DescriptionSpanResult.B524DescriptionSpanNumericValid || b524DescriptionSpanResult.B524DescriptionSpanImplementationValid || b524DescriptionSpanResult.B524DescriptionSpanLimits != nil || b524DescriptionSpanResult.B524DescriptionSpanGrade != B524DescriptionSpanProven || b524DescriptionSpanResult.B524DescriptionSpanReason != b524DescriptionSpanCase.B524DescriptionSpanReason || !b524DescriptionSpanReflect.DeepEqual(b524DescriptionSpanResult.B524DescriptionSpanMinimumRaw, b524DescriptionSpanCase.B524DescriptionSpanRaw[3:7]) || !b524DescriptionSpanReflect.DeepEqual(b524DescriptionSpanResult.B524DescriptionSpanMaximumRaw, b524DescriptionSpanCase.B524DescriptionSpanRaw[7:11]) || !b524DescriptionSpanReflect.DeepEqual(b524DescriptionSpanResult.B524DescriptionSpanStepRaw, b524DescriptionSpanCase.B524DescriptionSpanRaw[11:15]) {
			b524DescriptionSpanT.Fatalf("ab44 invalid numeric: %+v", b524DescriptionSpanResult)
		}
	}
	for _, b524DescriptionSpanCodec := range []B524DescriptionSpanCodec{B524DescriptionSpanDate, B524DescriptionSpanTime} {
		for b524DescriptionSpanIndex := 3; b524DescriptionSpanIndex < 12; b524DescriptionSpanIndex++ {
			b524DescriptionSpanRaw := []byte{2, 2, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
			b524DescriptionSpanRaw[b524DescriptionSpanIndex] = 255
			b524DescriptionSpanResult := B524DecodeDescription(b524DescriptionSpanTestContext(), b524DescriptionSpanCodec, 3, B524DescriptionSpanProven, b524DescriptionSpanRaw)
			if !b524DescriptionSpanResult.B524DescriptionSpanImplementationValid || b524DescriptionSpanResult.B524DescriptionSpanNumericValid || b524DescriptionSpanResult.B524DescriptionSpanLimits != nil || b524DescriptionSpanResult.B524DescriptionSpanStepMeaning != B524DescriptionSpanUnknown || !b524DescriptionSpanReflect.DeepEqual(b524DescriptionSpanResult.B524DescriptionSpanStepRaw, b524DescriptionSpanRaw[9:]) {
				b524DescriptionSpanT.Fatal("ab44 date raw mutation")
			}
		}
	}
}

func TestAb44Ownership(b524DescriptionSpanT *b524DescriptionSpanTesting.T) {
	b524DescriptionSpanContext := b524DescriptionSpanTestContext()
	b524DescriptionSpanExpectedContext := b524DescriptionSpanTestContext()
	b524DescriptionSpanRaw := b524DescriptionSpanTestEXP(-20, 20, 1)
	b524DescriptionSpanExpectedRaw := b524DescriptionSpanTestEXP(-20, 20, 1)
	b524DescriptionSpanResult := B524DecodeDescription(b524DescriptionSpanContext, B524DescriptionSpanEXP, 4, B524DescriptionSpanHypothesis, b524DescriptionSpanRaw)
	for b524DescriptionSpanIndex := range b524DescriptionSpanRaw {
		b524DescriptionSpanRaw[b524DescriptionSpanIndex] = 255
	}
	for b524DescriptionSpanIndex := 0; b524DescriptionSpanIndex < 9; b524DescriptionSpanIndex++ {
		b524DescriptionSpanReflect.ValueOf(b524DescriptionSpanContext).Field(b524DescriptionSpanIndex).Elem().SetZero()
	}
	if !b524DescriptionSpanReflect.DeepEqual(b524DescriptionSpanResult.B524DescriptionSpanContext, b524DescriptionSpanExpectedContext) || !b524DescriptionSpanReflect.DeepEqual(b524DescriptionSpanResult.B524DescriptionSpanRaw, b524DescriptionSpanExpectedRaw) || !b524DescriptionSpanReflect.DeepEqual(b524DescriptionSpanResult.B524DescriptionSpanMinimumRaw, b524DescriptionSpanExpectedRaw[3:7]) {
		b524DescriptionSpanT.Fatal("ab44 input alias")
	}
	b524DescriptionSpanResult.B524DescriptionSpanRaw[3] = 255
	if !b524DescriptionSpanReflect.DeepEqual(b524DescriptionSpanResult.B524DescriptionSpanMinimumRaw, b524DescriptionSpanExpectedRaw[3:7]) {
		b524DescriptionSpanT.Fatal("ab44 raw aliases spans")
	}
	b524DescriptionSpanResult.B524DescriptionSpanMinimumRaw[0] = 254
	b524DescriptionSpanResult.B524DescriptionSpanMaximumRaw[0] = 253
	b524DescriptionSpanResult.B524DescriptionSpanStepRaw[0] = 252
	*b524DescriptionSpanResult.B524DescriptionSpanMetadataOP = 99
	b524DescriptionSpanResult.B524DescriptionSpanLimits.B524DescriptionSpanMinimum = 99
	for b524DescriptionSpanIndex := 0; b524DescriptionSpanIndex < 9; b524DescriptionSpanIndex++ {
		b524DescriptionSpanReflect.ValueOf(b524DescriptionSpanResult.B524DescriptionSpanContext).Field(b524DescriptionSpanIndex).Elem().SetZero()
	}
	b524DescriptionSpanNext := B524DecodeDescription(b524DescriptionSpanExpectedContext, B524DescriptionSpanEXP, 4, B524DescriptionSpanHypothesis, b524DescriptionSpanExpectedRaw)
	if !b524DescriptionSpanReflect.DeepEqual(b524DescriptionSpanNext.B524DescriptionSpanContext, b524DescriptionSpanExpectedContext) || !b524DescriptionSpanReflect.DeepEqual(b524DescriptionSpanNext.B524DescriptionSpanRaw, b524DescriptionSpanExpectedRaw) || !b524DescriptionSpanReflect.DeepEqual(b524DescriptionSpanNext.B524DescriptionSpanMinimumRaw, b524DescriptionSpanExpectedRaw[3:7]) || !b524DescriptionSpanReflect.DeepEqual(b524DescriptionSpanNext.B524DescriptionSpanMaximumRaw, b524DescriptionSpanExpectedRaw[7:11]) || !b524DescriptionSpanReflect.DeepEqual(b524DescriptionSpanNext.B524DescriptionSpanStepRaw, b524DescriptionSpanExpectedRaw[11:]) || *b524DescriptionSpanNext.B524DescriptionSpanLimits != (B524DescriptionSpanLimits{-20, 20, 1}) || *b524DescriptionSpanNext.B524DescriptionSpanMetadataOP != 2 {
		b524DescriptionSpanT.Fatal("ab44 later result alias")
	}
	b524DescriptionSpanRequest, b524DescriptionSpanErr := B524BuildDescription(b524DescriptionSpanExpectedContext)
	if b524DescriptionSpanErr != nil {
		b524DescriptionSpanT.Fatal(b524DescriptionSpanErr)
	}
	b524DescriptionSpanRequest[0] = 99
	b524DescriptionSpanNextRequest, b524DescriptionSpanErr := B524BuildDescription(b524DescriptionSpanExpectedContext)
	if b524DescriptionSpanErr != nil || !b524DescriptionSpanReflect.DeepEqual(b524DescriptionSpanNextRequest, []byte{1, 2, 1, 2, 0}) {
		b524DescriptionSpanT.Fatal("ab44 request alias")
	}
}

func TestAb44ZeroAndRawOwnership(b524DescriptionSpanT *b524DescriptionSpanTesting.T) {
	for _, b524DescriptionSpanCase := range []struct {
		B524DescriptionSpanCodec B524DescriptionSpanCodec `json:"Ab44Codec"`
		B524DescriptionSpanWidth int                      `json:"Ab44Width"`
		B524DescriptionSpanRaw   []byte                   `json:"Ab44Raw"`
	}{
		{B524DescriptionSpanUCH, 1, []byte{2, 2, 0, 0, 0, 0}},
		{B524DescriptionSpanUIN, 2, []byte{2, 2, 0, 0, 0, 0, 0, 0, 0}},
		{B524DescriptionSpanEXP, 4, []byte{2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
	} {
		b524DescriptionSpanResult := B524DecodeDescription(b524DescriptionSpanTestContext(), b524DescriptionSpanCase.B524DescriptionSpanCodec, b524DescriptionSpanCase.B524DescriptionSpanWidth, B524DescriptionSpanHypothesis, b524DescriptionSpanCase.B524DescriptionSpanRaw)
		if !b524DescriptionSpanResult.B524DescriptionSpanNumericValid || b524DescriptionSpanResult.B524DescriptionSpanLimits == nil || *b524DescriptionSpanResult.B524DescriptionSpanLimits != (B524DescriptionSpanLimits{0, 0, 0}) {
			b524DescriptionSpanT.Fatal("ab44 literal zero/equal bounds", b524DescriptionSpanCase.B524DescriptionSpanCodec)
		}
	}
	for _, b524DescriptionSpanGrade := range []B524DescriptionSpanGrade{"", B524DescriptionSpanUnknown} {
		b524DescriptionSpanContext := b524DescriptionSpanTestContext()
		b524DescriptionSpanRaw := []byte{2, 2, 0, 0, 0, 4, 0, 1, 0}
		b524DescriptionSpanResult := B524DecodeDescription(b524DescriptionSpanContext, B524DescriptionSpanUIN, 2, b524DescriptionSpanGrade, b524DescriptionSpanRaw)
		b524DescriptionSpanRaw[3] = 99
		*b524DescriptionSpanContext.B524DescriptionSpanII = 99
		if b524DescriptionSpanResult.B524DescriptionSpanRaw[3] != 0 || *b524DescriptionSpanResult.B524DescriptionSpanContext.B524DescriptionSpanII != 1 || !b524DescriptionSpanReflect.DeepEqual(b524DescriptionSpanResult.B524DescriptionSpanMinimumRaw, []byte{0, 0}) {
			b524DescriptionSpanT.Fatal("ab44 raw-only input alias")
		}
		b524DescriptionSpanResult.B524DescriptionSpanMinimumRaw[0] = 88
		if b524DescriptionSpanResult.B524DescriptionSpanRaw[3] != 0 || b524DescriptionSpanResult.B524DescriptionSpanMaximumRaw[0] != 4 || b524DescriptionSpanResult.B524DescriptionSpanStepRaw[0] != 1 {
			b524DescriptionSpanT.Fatal("ab44 raw-only span alias")
		}
	}
}
