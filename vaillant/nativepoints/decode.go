package nativepoints

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
	"strconv"
)

// Value is an exact decoded value. Native and SI are exact rationals written
// as a decimal when the decimal terminates and as numerator/denominator
// otherwise; both parse with math/big.Rat.SetString. SIUnit is K, Pa, J, W,
// s, m3/s, 1/s or 1 for a dimensionless ratio or count.
type Value struct {
	Native     string `json:"native"`
	NativeUnit string `json:"native_unit"`
	SI         string `json:"si"`
	SIUnit     string `json:"si_unit"`
}

// DecodeFailure says why Decode refused a reply.
type DecodeFailure string

// Decode failures: the row is withheld, its codec or scale is not decodable,
// its unit has no SI conversion, the reply is shorter or longer than the codec
// width, or the value is NaN, infinite or the codec sentinel.
const (
	DecodeWithheld     DecodeFailure = "withheld"
	DecodeCodecUnknown DecodeFailure = "codec_unknown"
	DecodeUnitUnknown  DecodeFailure = "unit_unknown"
	DecodeShortReply   DecodeFailure = "short_reply"
	DecodeLongReply    DecodeFailure = "long_reply"
	DecodeNotANumber   DecodeFailure = "not_a_number"
	DecodeInfinite     DecodeFailure = "infinite"
	DecodeSentinel     DecodeFailure = "sentinel"
)

// DecodeError is the typed refusal of Decode; Length is the reply length and
// Width the codec width, zero when the codec is not known.
type DecodeError struct {
	DeviceClass DeviceClass
	Name        string
	Failure     DecodeFailure
	Length      int
	Width       int
}

func (decodeError *DecodeError) Error() string {
	return fmt.Sprintf("native point %s/%s: %s (reply %d bytes, codec width %d)",
		decodeError.DeviceClass, decodeError.Name, decodeError.Failure, decodeError.Length, decodeError.Width)
}

type siConversion struct {
	unit   string
	factor string
	offset string
}

var siConversions = map[string]siConversion{
	"°C":    {"K", "1", "273.15"},
	"K":     {"K", "1", "0"},
	"bar":   {"Pa", "100000", "0"},
	"kWh":   {"J", "3600000", "0"},
	"kW":    {"W", "1000", "0"},
	"W":     {"W", "1", "0"},
	"h":     {"s", "3600", "0"},
	"min":   {"s", "60", "0"},
	"s":     {"s", "1", "0"},
	"%":     {"1", "1/100", "0"},
	"count": {"1", "1", "0"},
	"steps": {"1", "1", "0"},
	"L/min": {"m3/s", "1/60000", "0"},
	"L/h":   {"m3/s", "1/3600000", "0"},
	"rpm":   {"1/s", "1/60", "0"},
}

var codecWidths = map[Codec]int{
	CodecFloat32LE: 4,
	CodecUint16LE:  2,
	CodecData2c:    2,
	CodecData2b:    2,
	CodecUIN:       2,
}

// Decode decodes a retained raw reply of an exportable point into its exact
// native value and converts it exactly into SI. A float32 is taken at the
// shortest decimal that reads back as the same float32, not at its widened
// float64 expansion. The reply is not modified.
func Decode(point Point, raw []byte) (Value, error) {
	refuse := func(failure DecodeFailure, width int) (Value, error) {
		return Value{}, &DecodeError{DeviceClass: point.DeviceClass, Name: point.Name, Failure: failure, Length: len(raw), Width: width}
	}
	if !point.Exportable || point.Withheld != WithheldNone {
		return refuse(DecodeWithheld, 0)
	}
	width, known := codecWidths[point.Codec]
	if !known || point.ScaleDivisor == 0 {
		return refuse(DecodeCodecUnknown, 0)
	}
	conversion, known := siConversions[point.NativeUnit]
	if !known {
		return refuse(DecodeUnitUnknown, width)
	}
	if len(raw) < width {
		return refuse(DecodeShortReply, width)
	}
	if len(raw) > width {
		return refuse(DecodeLongReply, width)
	}
	native := new(big.Rat)
	switch point.Codec {
	case CodecFloat32LE:
		value := math.Float32frombits(binary.LittleEndian.Uint32(raw))
		if math.IsNaN(float64(value)) {
			return refuse(DecodeNotANumber, width)
		}
		if math.IsInf(float64(value), 0) {
			return refuse(DecodeInfinite, width)
		}
		native.SetString(strconv.FormatFloat(float64(value), 'f', -1, 32))
	case CodecUint16LE:
		native.SetInt64(int64(binary.LittleEndian.Uint16(raw)))
	case CodecUIN:
		word := binary.LittleEndian.Uint16(raw)
		if word == 0xFFFF {
			return refuse(DecodeSentinel, width)
		}
		native.SetInt64(int64(word))
	case CodecData2c, CodecData2b:
		word := binary.LittleEndian.Uint16(raw)
		if word == 0x8000 {
			return refuse(DecodeSentinel, width)
		}
		divisor := int64(16)
		if point.Codec == CodecData2b {
			divisor = 256
		}
		native.SetFrac64(int64(int16(word)), divisor)
	}
	native.Quo(native, new(big.Rat).SetInt64(int64(point.ScaleDivisor)))
	factor, _ := new(big.Rat).SetString(conversion.factor)
	offset, _ := new(big.Rat).SetString(conversion.offset)
	si := new(big.Rat).Mul(native, factor)
	si.Add(si, offset)
	return Value{Native: exactText(native), NativeUnit: point.NativeUnit, SI: exactText(si), SIUnit: conversion.unit}, nil
}

func exactText(value *big.Rat) string {
	if value.IsInt() {
		return value.Num().String()
	}
	remainder := new(big.Int).Set(value.Denom())
	two, five, zero := big.NewInt(2), big.NewInt(5), new(big.Int)
	modulus := new(big.Int)
	twos, fives := 0, 0
	for modulus.Mod(remainder, two).Cmp(zero) == 0 {
		remainder.Quo(remainder, two)
		twos++
	}
	for modulus.Mod(remainder, five).Cmp(zero) == 0 {
		remainder.Quo(remainder, five)
		fives++
	}
	if remainder.Cmp(big.NewInt(1)) != 0 {
		return value.RatString()
	}
	return value.FloatString(max(twos, fives))
}
