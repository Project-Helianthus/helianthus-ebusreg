package system

import (
	ab44Binary "encoding/binary"
	ab44Errors "errors"
	ab44Math "math"
)

type Ab44Context struct {
	Ab44Family      *string
	Ab44Software    *uint16
	Ab44Hardware    *uint16
	Ab44APIVersion  *byte
	Ab44APIRevision *byte
	Ab44OP          *byte
	Ab44GG          *byte
	Ab44II          *byte
	Ab44RR          *uint16
}

type Ab44Grade string

const (
	Ab44Unknown    Ab44Grade = "Unknown"
	Ab44Hypothesis Ab44Grade = "Hypothesis"
	Ab44Proven     Ab44Grade = "Proven"
)

type Ab44Codec string

const (
	Ab44UCH  Ab44Codec = "UCH"
	Ab44UIN  Ab44Codec = "UIN"
	Ab44EXP  Ab44Codec = "EXP"
	Ab44Date Ab44Codec = "DATE"
	Ab44Time Ab44Codec = "TIME"
)

type Ab44Limits struct {
	Ab44Minimum float64
	Ab44Maximum float64
	Ab44Step    float64
}

type Ab44Description struct {
	Ab44Context             Ab44Context
	Ab44Codec               Ab44Codec
	Ab44Width               int
	Ab44Grade               Ab44Grade
	Ab44Raw                 []byte
	Ab44MinimumRaw          []byte
	Ab44MaximumRaw          []byte
	Ab44StepRaw             []byte
	Ab44MetadataOP          *byte
	Ab44Limits              *Ab44Limits
	Ab44ImplementationValid bool
	Ab44NumericValid        bool
	Ab44StepMeaning         Ab44Grade
	Ab44Reason              string
}

func ab44Copy[ab44T any](ab44Value *ab44T) *ab44T {
	if ab44Value == nil {
		return nil
	}
	ab44Owned := *ab44Value
	return &ab44Owned
}

func ab44CloneContext(ab44Context Ab44Context) Ab44Context {
	return Ab44Context{
		Ab44Family: ab44Copy(ab44Context.Ab44Family), Ab44Software: ab44Copy(ab44Context.Ab44Software),
		Ab44Hardware: ab44Copy(ab44Context.Ab44Hardware), Ab44APIVersion: ab44Copy(ab44Context.Ab44APIVersion),
		Ab44APIRevision: ab44Copy(ab44Context.Ab44APIRevision), Ab44OP: ab44Copy(ab44Context.Ab44OP),
		Ab44GG: ab44Copy(ab44Context.Ab44GG), Ab44II: ab44Copy(ab44Context.Ab44II), Ab44RR: ab44Copy(ab44Context.Ab44RR),
	}
}

func ab44ContextReason(ab44Context Ab44Context) string {
	if ab44Context.Ab44Family == nil || *ab44Context.Ab44Family == "" || ab44Context.Ab44Software == nil ||
		ab44Context.Ab44Hardware == nil || ab44Context.Ab44APIVersion == nil || ab44Context.Ab44APIRevision == nil ||
		ab44Context.Ab44OP == nil || ab44Context.Ab44GG == nil || ab44Context.Ab44II == nil || ab44Context.Ab44RR == nil {
		return "ab44_incomplete_context"
	}
	if *ab44Context.Ab44OP != 0x01 && *ab44Context.Ab44OP != 0x07 {
		return "ab44_unsupported_op"
	}
	return ""
}

func Ab44BuildDescription(ab44Context Ab44Context) ([]byte, error) {
	if ab44Reason := ab44ContextReason(ab44Context); ab44Reason != "" {
		return nil, ab44Errors.New(ab44Reason)
	}
	return []byte{*ab44Context.Ab44OP, *ab44Context.Ab44GG, *ab44Context.Ab44II, byte(*ab44Context.Ab44RR), byte(*ab44Context.Ab44RR >> 8)}, nil
}

func Ab44DecodeDescription(ab44Context Ab44Context, ab44ExplicitCodec Ab44Codec, ab44ExplicitWidth int, ab44Grade Ab44Grade, ab44Raw []byte) Ab44Description {
	ab44Result := Ab44Description{
		Ab44Context: ab44CloneContext(ab44Context), Ab44Codec: ab44ExplicitCodec, Ab44Width: ab44ExplicitWidth,
		Ab44Grade: ab44Grade, Ab44Raw: append([]byte(nil), ab44Raw...), Ab44StepMeaning: Ab44Unknown,
	}
	if ab44ExplicitWidth < 1 || ab44ExplicitWidth > 4 {
		ab44Result.Ab44Reason = "ab44_bad_width"
		return ab44Result
	}
	if len(ab44Raw) != 3+3*ab44ExplicitWidth {
		ab44Result.Ab44Reason = "ab44_malformed_length"
		return ab44Result
	}
	ab44Result.Ab44MinimumRaw = append([]byte(nil), ab44Raw[3:3+ab44ExplicitWidth]...)
	ab44Result.Ab44MaximumRaw = append([]byte(nil), ab44Raw[3+ab44ExplicitWidth:3+2*ab44ExplicitWidth]...)
	ab44Result.Ab44StepRaw = append([]byte(nil), ab44Raw[3+2*ab44ExplicitWidth:]...)
	if ab44Reason := ab44ContextReason(ab44Context); ab44Reason != "" {
		ab44Result.Ab44Reason = ab44Reason
		return ab44Result
	}
	if ab44Raw[0] != *ab44Context.Ab44GG || ab44Binary.LittleEndian.Uint16(ab44Raw[1:3]) != *ab44Context.Ab44RR {
		ab44Result.Ab44Reason = "ab44_echo_mismatch"
		return ab44Result
	}
	ab44MetadataOP := byte(0x02)
	if *ab44Context.Ab44OP == 0x07 {
		ab44MetadataOP = 0x06
	}
	ab44Result.Ab44MetadataOP = &ab44MetadataOP
	ab44ExpectedWidth := 0
	switch ab44ExplicitCodec {
	case Ab44UCH:
		ab44ExpectedWidth = 1
	case Ab44UIN:
		ab44ExpectedWidth = 2
	case Ab44EXP:
		ab44ExpectedWidth = 4
	case Ab44Date, Ab44Time:
		ab44ExpectedWidth = 3
	default:
		ab44Result.Ab44Reason = "ab44_unsupported_codec"
		return ab44Result
	}
	if ab44ExplicitWidth != ab44ExpectedWidth {
		ab44Result.Ab44Reason = "ab44_codec_width_mismatch"
		return ab44Result
	}
	if ab44Grade != Ab44Hypothesis && ab44Grade != Ab44Proven {
		ab44Result.Ab44Reason = "ab44_codec_grade_required"
		return ab44Result
	}
	if ab44ExplicitCodec == Ab44Date || ab44ExplicitCodec == Ab44Time {
		ab44Result.Ab44ImplementationValid = true
		ab44Result.Ab44Reason = "ab44_date_time_raw_only"
		return ab44Result
	}
	ab44Number := func(ab44Span []byte) float64 {
		switch ab44ExplicitCodec {
		case Ab44UCH:
			return float64(ab44Span[0])
		case Ab44UIN:
			return float64(ab44Binary.LittleEndian.Uint16(ab44Span))
		default:
			return float64(ab44Math.Float32frombits(ab44Binary.LittleEndian.Uint32(ab44Span)))
		}
	}
	ab44Limits := Ab44Limits{ab44Number(ab44Result.Ab44MinimumRaw), ab44Number(ab44Result.Ab44MaximumRaw), ab44Number(ab44Result.Ab44StepRaw)}
	for _, ab44Value := range []float64{ab44Limits.Ab44Minimum, ab44Limits.Ab44Maximum, ab44Limits.Ab44Step} {
		if ab44Math.IsNaN(ab44Value) || ab44Math.IsInf(ab44Value, 0) {
			ab44Result.Ab44Reason = "ab44_nonfinite_limits"
			return ab44Result
		}
	}
	if ab44Limits.Ab44Minimum > ab44Limits.Ab44Maximum {
		ab44Result.Ab44Reason = "ab44_reversed_limits"
		return ab44Result
	}
	if ab44Limits.Ab44Step < 0 {
		ab44Result.Ab44Reason = "ab44_negative_step"
		return ab44Result
	}
	ab44Result.Ab44Limits = &ab44Limits
	ab44Result.Ab44ImplementationValid = true
	ab44Result.Ab44NumericValid = true
	ab44Result.Ab44Reason = "ab44_numeric_limits"
	return ab44Result
}
