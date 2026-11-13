package system

import (
	"bytes"
	"encoding/binary"
	"math"
)

type ValueCodec string

const (
	ValueCodecBOOL ValueCodec = "BOOL"
	ValueCodecUCH  ValueCodec = "UCH"
	ValueCodecUIN  ValueCodec = "UIN"
	ValueCodecI8   ValueCodec = "I8"
	ValueCodecI16  ValueCodec = "I16"
	ValueCodecU32  ValueCodec = "U32"
	ValueCodecI32  ValueCodec = "I32"
	ValueCodecEXP  ValueCodec = "EXP"
	ValueCodecDATE ValueCodec = "DATE_DMY_U8X3"
	ValueCodecTIME ValueCodec = "TIME_HMS_U8X3"
	ValueCodecFW   ValueCodec = "FW"
	ValueCodecFWU  ValueCodec = "FWU"
	ValueCodecSTR  ValueCodec = "STR"
	ValueCodecHEX  ValueCodec = "HEX"
)

type ValueCodecGrade string

const (
	ValueCodecUnknown    ValueCodecGrade = "Unknown"
	ValueCodecHypothesis ValueCodecGrade = "Hypothesis"
	ValueCodecProven     ValueCodecGrade = "Proven"
)

// ValueCodecValue keeps native evidence independently of an explicitly requested
// interpretation. ValueCodecQualified describes software decoding with a Proven
// caller-supplied grade, never physical support or register qualification.
type ValueCodecValue struct {
	ValueCodecRaw            []byte          `json:"Ab11Raw"`
	ValueCodecGrade          ValueCodecGrade `json:"Ab11Grade"`
	ValueCodecInterpretation any             `json:"Ab11Interpretation"`
	ValueCodecQualified      bool            `json:"Ab11Qualified"`
}

type ValueCodecDate struct {
	ValueCodecDay   uint8 `json:"Ab11Day"`
	ValueCodecMonth uint8 `json:"Ab11Month"`
	ValueCodecYear  uint8 `json:"Ab11Year"`
}

type ValueCodecTime struct {
	ValueCodecHour   uint8 `json:"Ab11Hour"`
	ValueCodecMinute uint8 `json:"Ab11Minute"`
	ValueCodecSecond uint8 `json:"Ab11Second"`
}

// ValueCodecString retains every byte; ValueCodecLatin1 maps each byte to its Unicode
// code point without trimming or interpreting termination.
type ValueCodecString struct {
	ValueCodecBytes  []byte `json:"Ab11Bytes"`
	ValueCodecLatin1 string `json:"Ab11Latin1"`
}

// DecodeValue requires separate codec, byte width and evidence grade.
// For STR:n and HEX:n select ValueCodecSTR or ValueCodecHEX and pass n as valueCodecWidth.
// Failed admission returns owned raw evidence and no interpretation.
// Hypothesis interpretations remain Hypothesis and are never qualified.
func DecodeValue(valueCodecCodec ValueCodec, valueCodecWidth int, valueCodecGrade ValueCodecGrade, valueCodecRaw []byte) ValueCodecValue {
	valueCodecResult := ValueCodecValue{ValueCodecRaw: append([]byte(nil), valueCodecRaw...), ValueCodecGrade: valueCodecGrade}
	if valueCodecGrade != ValueCodecProven && valueCodecGrade != ValueCodecHypothesis {
		return valueCodecResult
	}
	valueCodecExpected := 0
	switch valueCodecCodec {
	case ValueCodecBOOL, ValueCodecUCH, ValueCodecI8:
		valueCodecExpected = 1
	case ValueCodecUIN, ValueCodecI16:
		valueCodecExpected = 2
	case ValueCodecU32, ValueCodecI32, ValueCodecEXP:
		valueCodecExpected = 4
	case ValueCodecDATE, ValueCodecTIME, ValueCodecFW, ValueCodecFWU:
		valueCodecExpected = 3
	case ValueCodecSTR, ValueCodecHEX:
		valueCodecExpected = valueCodecWidth
	default:
		return valueCodecResult
	}
	if valueCodecWidth <= 0 || valueCodecWidth != valueCodecExpected || len(valueCodecRaw) != valueCodecWidth {
		return valueCodecResult
	}
	switch valueCodecCodec {
	case ValueCodecBOOL:
		if valueCodecRaw[0] > 1 {
			return valueCodecResult
		}
		valueCodecResult.ValueCodecInterpretation = valueCodecRaw[0] == 1
	case ValueCodecUCH:
		valueCodecResult.ValueCodecInterpretation = valueCodecRaw[0]
	case ValueCodecUIN:
		valueCodecResult.ValueCodecInterpretation = binary.LittleEndian.Uint16(valueCodecRaw)
	case ValueCodecI8:
		valueCodecResult.ValueCodecInterpretation = int8(valueCodecRaw[0])
	case ValueCodecI16:
		valueCodecResult.ValueCodecInterpretation = int16(binary.LittleEndian.Uint16(valueCodecRaw))
	case ValueCodecU32:
		valueCodecResult.ValueCodecInterpretation = binary.LittleEndian.Uint32(valueCodecRaw)
	case ValueCodecI32:
		valueCodecResult.ValueCodecInterpretation = int32(binary.LittleEndian.Uint32(valueCodecRaw))
	case ValueCodecEXP:
		valueCodecFloat := math.Float32frombits(binary.LittleEndian.Uint32(valueCodecRaw))
		if math.IsNaN(float64(valueCodecFloat)) || math.IsInf(float64(valueCodecFloat), 0) {
			return valueCodecResult
		}
		valueCodecResult.ValueCodecInterpretation = valueCodecFloat
	case ValueCodecDATE:
		valueCodecResult.ValueCodecInterpretation = ValueCodecDate{valueCodecRaw[0], valueCodecRaw[1], valueCodecRaw[2]}
	case ValueCodecTIME:
		valueCodecResult.ValueCodecInterpretation = ValueCodecTime{valueCodecRaw[0], valueCodecRaw[1], valueCodecRaw[2]}
	case ValueCodecFW, ValueCodecFWU:
		if valueCodecCodec == ValueCodecFWU && bytes.Equal(valueCodecRaw, []byte{255, 255, 255}) {
			// All-0xFF is the unavailable-firmware sentinel; keep raw only.
			return valueCodecResult
		}
		valueCodecParts := [3]uint8{}
		for valueCodecIndex, valueCodecByte := range valueCodecRaw {
			if valueCodecCodec == ValueCodecFW {
				if valueCodecByte>>4 > 9 || valueCodecByte&15 > 9 {
					return valueCodecResult
				}
				valueCodecParts[valueCodecIndex] = (valueCodecByte>>4)*10 + valueCodecByte&15
			} else {
				valueCodecParts[valueCodecIndex] = valueCodecByte
			}
		}
		valueCodecResult.ValueCodecInterpretation = valueCodecParts
	case ValueCodecSTR:
		valueCodecRunes := make([]rune, len(valueCodecRaw))
		for valueCodecIndex, valueCodecByte := range valueCodecRaw {
			valueCodecRunes[valueCodecIndex] = rune(valueCodecByte)
		}
		valueCodecResult.ValueCodecInterpretation = ValueCodecString{append([]byte(nil), valueCodecRaw...), string(valueCodecRunes)}
	case ValueCodecHEX:
		valueCodecResult.ValueCodecInterpretation = append([]byte(nil), valueCodecRaw...)
	}
	valueCodecResult.ValueCodecQualified = valueCodecGrade == ValueCodecProven
	return valueCodecResult
}
