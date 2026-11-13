package system

import (
	b524DescriptionSpanBinary "encoding/binary"
	b524DescriptionSpanErrors "errors"
	b524DescriptionSpanMath "math"
)

type B524DescriptionSpanContext struct {
	B524DescriptionSpanFamily      *string `json:"family"`
	B524DescriptionSpanSoftware    *uint16 `json:"software"`
	B524DescriptionSpanHardware    *uint16 `json:"hardware"`
	B524DescriptionSpanAPIVersion  *byte   `json:"api_version"`
	B524DescriptionSpanAPIRevision *byte   `json:"api_revision"`
	B524DescriptionSpanOP          *byte   `json:"op"`
	B524DescriptionSpanGG          *byte   `json:"gg"`
	B524DescriptionSpanII          *byte   `json:"ii"`
	B524DescriptionSpanRR          *uint16 `json:"rr"`
}

type B524DescriptionSpanGrade string

const (
	B524DescriptionSpanUnknown    B524DescriptionSpanGrade = "Unknown"
	B524DescriptionSpanHypothesis B524DescriptionSpanGrade = "Hypothesis"
	B524DescriptionSpanProven     B524DescriptionSpanGrade = "Proven"
)

type B524DescriptionSpanCodec string

const (
	B524DescriptionSpanUCH  B524DescriptionSpanCodec = "UCH"
	B524DescriptionSpanUIN  B524DescriptionSpanCodec = "UIN"
	B524DescriptionSpanEXP  B524DescriptionSpanCodec = "EXP"
	B524DescriptionSpanDate B524DescriptionSpanCodec = "DATE"
	B524DescriptionSpanTime B524DescriptionSpanCodec = "TIME"
)

type B524DescriptionSpanLimits struct {
	B524DescriptionSpanMinimum float64 `json:"minimum"`
	B524DescriptionSpanMaximum float64 `json:"maximum"`
	B524DescriptionSpanStep    float64 `json:"step"`
}

type B524Description struct {
	B524DescriptionSpanContext             B524DescriptionSpanContext `json:"context"`
	B524DescriptionSpanCodec               B524DescriptionSpanCodec   `json:"codec"`
	B524DescriptionSpanWidth               int                        `json:"width"`
	B524DescriptionSpanGrade               B524DescriptionSpanGrade   `json:"grade"`
	B524DescriptionSpanRaw                 []byte                     `json:"raw"`
	B524DescriptionSpanMinimumRaw          []byte                     `json:"minimum_raw"`
	B524DescriptionSpanMaximumRaw          []byte                     `json:"maximum_raw"`
	B524DescriptionSpanStepRaw             []byte                     `json:"step_raw"`
	B524DescriptionSpanMetadataOP          *byte                      `json:"metadata_op"`
	B524DescriptionSpanLimits              *B524DescriptionSpanLimits `json:"limits"`
	B524DescriptionSpanImplementationValid bool                       `json:"implementation_valid"`
	B524DescriptionSpanNumericValid        bool                       `json:"numeric_valid"`
	B524DescriptionSpanStepMeaning         B524DescriptionSpanGrade   `json:"step_meaning"`
	B524DescriptionSpanReason              string                     `json:"reason"`
}

func b524DescriptionSpanCopy[b524DescriptionSpanT any](b524DescriptionSpanValue *b524DescriptionSpanT) *b524DescriptionSpanT {
	if b524DescriptionSpanValue == nil {
		return nil
	}
	b524DescriptionSpanOwned := *b524DescriptionSpanValue
	return &b524DescriptionSpanOwned
}

func b524DescriptionSpanCloneContext(b524DescriptionSpanContext B524DescriptionSpanContext) B524DescriptionSpanContext {
	return B524DescriptionSpanContext{
		B524DescriptionSpanFamily: b524DescriptionSpanCopy(b524DescriptionSpanContext.B524DescriptionSpanFamily), B524DescriptionSpanSoftware: b524DescriptionSpanCopy(b524DescriptionSpanContext.B524DescriptionSpanSoftware),
		B524DescriptionSpanHardware: b524DescriptionSpanCopy(b524DescriptionSpanContext.B524DescriptionSpanHardware), B524DescriptionSpanAPIVersion: b524DescriptionSpanCopy(b524DescriptionSpanContext.B524DescriptionSpanAPIVersion),
		B524DescriptionSpanAPIRevision: b524DescriptionSpanCopy(b524DescriptionSpanContext.B524DescriptionSpanAPIRevision), B524DescriptionSpanOP: b524DescriptionSpanCopy(b524DescriptionSpanContext.B524DescriptionSpanOP),
		B524DescriptionSpanGG: b524DescriptionSpanCopy(b524DescriptionSpanContext.B524DescriptionSpanGG), B524DescriptionSpanII: b524DescriptionSpanCopy(b524DescriptionSpanContext.B524DescriptionSpanII), B524DescriptionSpanRR: b524DescriptionSpanCopy(b524DescriptionSpanContext.B524DescriptionSpanRR),
	}
}

func b524DescriptionSpanContextReason(b524DescriptionSpanContext B524DescriptionSpanContext) string {
	if b524DescriptionSpanContext.B524DescriptionSpanFamily == nil || *b524DescriptionSpanContext.B524DescriptionSpanFamily == "" || b524DescriptionSpanContext.B524DescriptionSpanSoftware == nil ||
		b524DescriptionSpanContext.B524DescriptionSpanHardware == nil || b524DescriptionSpanContext.B524DescriptionSpanAPIVersion == nil || b524DescriptionSpanContext.B524DescriptionSpanAPIRevision == nil ||
		b524DescriptionSpanContext.B524DescriptionSpanOP == nil || b524DescriptionSpanContext.B524DescriptionSpanGG == nil || b524DescriptionSpanContext.B524DescriptionSpanII == nil || b524DescriptionSpanContext.B524DescriptionSpanRR == nil {
		return "incomplete_context"
	}
	if *b524DescriptionSpanContext.B524DescriptionSpanOP != 0x01 && *b524DescriptionSpanContext.B524DescriptionSpanOP != 0x07 {
		return "unsupported_op"
	}
	return ""
}

func B524BuildDescription(b524DescriptionSpanContext B524DescriptionSpanContext) ([]byte, error) {
	if b524DescriptionSpanReason := b524DescriptionSpanContextReason(b524DescriptionSpanContext); b524DescriptionSpanReason != "" {
		return nil, b524DescriptionSpanErrors.New(b524DescriptionSpanReason)
	}
	return []byte{*b524DescriptionSpanContext.B524DescriptionSpanOP, *b524DescriptionSpanContext.B524DescriptionSpanGG, *b524DescriptionSpanContext.B524DescriptionSpanII, byte(*b524DescriptionSpanContext.B524DescriptionSpanRR), byte(*b524DescriptionSpanContext.B524DescriptionSpanRR >> 8)}, nil
}

func B524DecodeDescription(b524DescriptionSpanContext B524DescriptionSpanContext, b524DescriptionSpanExplicitCodec B524DescriptionSpanCodec, b524DescriptionSpanExplicitWidth int, b524DescriptionSpanGrade B524DescriptionSpanGrade, b524DescriptionSpanRaw []byte) B524Description {
	b524DescriptionSpanResult := B524Description{
		B524DescriptionSpanContext: b524DescriptionSpanCloneContext(b524DescriptionSpanContext), B524DescriptionSpanCodec: b524DescriptionSpanExplicitCodec, B524DescriptionSpanWidth: b524DescriptionSpanExplicitWidth,
		B524DescriptionSpanGrade: b524DescriptionSpanGrade, B524DescriptionSpanRaw: append([]byte(nil), b524DescriptionSpanRaw...), B524DescriptionSpanStepMeaning: B524DescriptionSpanUnknown,
	}
	if b524DescriptionSpanExplicitWidth < 1 || b524DescriptionSpanExplicitWidth > 4 {
		b524DescriptionSpanResult.B524DescriptionSpanReason = "bad_width"
		return b524DescriptionSpanResult
	}
	if len(b524DescriptionSpanRaw) != 3+3*b524DescriptionSpanExplicitWidth {
		b524DescriptionSpanResult.B524DescriptionSpanReason = "malformed_length"
		return b524DescriptionSpanResult
	}
	b524DescriptionSpanResult.B524DescriptionSpanMinimumRaw = append([]byte(nil), b524DescriptionSpanRaw[3:3+b524DescriptionSpanExplicitWidth]...)
	b524DescriptionSpanResult.B524DescriptionSpanMaximumRaw = append([]byte(nil), b524DescriptionSpanRaw[3+b524DescriptionSpanExplicitWidth:3+2*b524DescriptionSpanExplicitWidth]...)
	b524DescriptionSpanResult.B524DescriptionSpanStepRaw = append([]byte(nil), b524DescriptionSpanRaw[3+2*b524DescriptionSpanExplicitWidth:]...)
	if b524DescriptionSpanReason := b524DescriptionSpanContextReason(b524DescriptionSpanContext); b524DescriptionSpanReason != "" {
		b524DescriptionSpanResult.B524DescriptionSpanReason = b524DescriptionSpanReason
		return b524DescriptionSpanResult
	}
	if b524DescriptionSpanRaw[0] != *b524DescriptionSpanContext.B524DescriptionSpanGG || b524DescriptionSpanBinary.LittleEndian.Uint16(b524DescriptionSpanRaw[1:3]) != *b524DescriptionSpanContext.B524DescriptionSpanRR {
		b524DescriptionSpanResult.B524DescriptionSpanReason = "echo_mismatch"
		return b524DescriptionSpanResult
	}
	b524DescriptionSpanMetadataOP := byte(0x02)
	if *b524DescriptionSpanContext.B524DescriptionSpanOP == 0x07 {
		b524DescriptionSpanMetadataOP = 0x06
	}
	b524DescriptionSpanResult.B524DescriptionSpanMetadataOP = &b524DescriptionSpanMetadataOP
	b524DescriptionSpanExpectedWidth := 0
	switch b524DescriptionSpanExplicitCodec {
	case B524DescriptionSpanUCH:
		b524DescriptionSpanExpectedWidth = 1
	case B524DescriptionSpanUIN:
		b524DescriptionSpanExpectedWidth = 2
	case B524DescriptionSpanEXP:
		b524DescriptionSpanExpectedWidth = 4
	case B524DescriptionSpanDate, B524DescriptionSpanTime:
		b524DescriptionSpanExpectedWidth = 3
	default:
		b524DescriptionSpanResult.B524DescriptionSpanReason = "unsupported_codec"
		return b524DescriptionSpanResult
	}
	if b524DescriptionSpanExplicitWidth != b524DescriptionSpanExpectedWidth {
		b524DescriptionSpanResult.B524DescriptionSpanReason = "codec_width_mismatch"
		return b524DescriptionSpanResult
	}
	if b524DescriptionSpanGrade != B524DescriptionSpanHypothesis && b524DescriptionSpanGrade != B524DescriptionSpanProven {
		b524DescriptionSpanResult.B524DescriptionSpanReason = "codec_grade_required"
		return b524DescriptionSpanResult
	}
	if b524DescriptionSpanExplicitCodec == B524DescriptionSpanDate || b524DescriptionSpanExplicitCodec == B524DescriptionSpanTime {
		b524DescriptionSpanResult.B524DescriptionSpanImplementationValid = true
		b524DescriptionSpanResult.B524DescriptionSpanReason = "date_time_raw_only"
		return b524DescriptionSpanResult
	}
	b524DescriptionSpanNumber := func(b524DescriptionSpanSpan []byte) float64 {
		switch b524DescriptionSpanExplicitCodec {
		case B524DescriptionSpanUCH:
			return float64(b524DescriptionSpanSpan[0])
		case B524DescriptionSpanUIN:
			return float64(b524DescriptionSpanBinary.LittleEndian.Uint16(b524DescriptionSpanSpan))
		default:
			return float64(b524DescriptionSpanMath.Float32frombits(b524DescriptionSpanBinary.LittleEndian.Uint32(b524DescriptionSpanSpan)))
		}
	}
	b524DescriptionSpanLimits := B524DescriptionSpanLimits{b524DescriptionSpanNumber(b524DescriptionSpanResult.B524DescriptionSpanMinimumRaw), b524DescriptionSpanNumber(b524DescriptionSpanResult.B524DescriptionSpanMaximumRaw), b524DescriptionSpanNumber(b524DescriptionSpanResult.B524DescriptionSpanStepRaw)}
	for _, b524DescriptionSpanValue := range []float64{b524DescriptionSpanLimits.B524DescriptionSpanMinimum, b524DescriptionSpanLimits.B524DescriptionSpanMaximum, b524DescriptionSpanLimits.B524DescriptionSpanStep} {
		if b524DescriptionSpanMath.IsNaN(b524DescriptionSpanValue) || b524DescriptionSpanMath.IsInf(b524DescriptionSpanValue, 0) {
			b524DescriptionSpanResult.B524DescriptionSpanReason = "nonfinite_limits"
			return b524DescriptionSpanResult
		}
	}
	if b524DescriptionSpanLimits.B524DescriptionSpanMinimum > b524DescriptionSpanLimits.B524DescriptionSpanMaximum {
		b524DescriptionSpanResult.B524DescriptionSpanReason = "reversed_limits"
		return b524DescriptionSpanResult
	}
	if b524DescriptionSpanLimits.B524DescriptionSpanStep < 0 {
		b524DescriptionSpanResult.B524DescriptionSpanReason = "negative_step"
		return b524DescriptionSpanResult
	}
	b524DescriptionSpanResult.B524DescriptionSpanLimits = &b524DescriptionSpanLimits
	b524DescriptionSpanResult.B524DescriptionSpanImplementationValid = true
	b524DescriptionSpanResult.B524DescriptionSpanNumericValid = true
	b524DescriptionSpanResult.B524DescriptionSpanReason = "numeric_limits"
	return b524DescriptionSpanResult
}
