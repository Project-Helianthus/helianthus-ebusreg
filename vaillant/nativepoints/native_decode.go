package nativepoints

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"unicode/utf8"
)

// NativeValue is the exact native value of a reply. Native is a decimal
// integer or terminating decimal (numerator/denominator otherwise) for a
// numeric value and for the code of an enumeration, true or false for a
// Boolean, DD.MM.YY for a date, HH:MM:SS for a time, the major, minor and
// patch byte-decimal numbers padded to two digits and joined by dots (08.05.00)
// for a version and the text up to its first NUL byte for a text. Label is the documented label of the code,
// empty when the documentation states none. ValueClass is the class the value
// was decoded as: boolean, enumeration, date, time, version, text or numeric.
type NativeValue struct {
	ValueClass ValueClass `json:"value_class"`
	Native     string     `json:"native"`
	NativeUnit string     `json:"native_unit"`
	Label      string     `json:"label"`
}

// DecodeNative failures beyond those of Decode: a code outside the documented
// Boolean range or date range, and a text that is not UTF-8.
const (
	DecodeOutOfRange  DecodeFailure = "out_of_range"
	DecodeInvalidText DecodeFailure = "invalid_text"
)

var nativeWidths = map[Codec][2]int{
	CodecFloat32LE:        {4, 4},
	CodecUint16LE:         {2, 2},
	CodecUint32LE:         {4, 4},
	CodecUint8:            {1, 1},
	CodecSignedByte:       {1, 1},
	CodecUint8OrUint16LE:  {1, 2},
	CodecDayMonthYear:     {3, 3},
	CodecHourMinuteSecond: {3, 3},
	CodecMajorMinorPatch:  {3, 3},
	CodecData2c:           {2, 2},
	CodecData2b:           {2, 2},
	CodecUIN:              {2, 2},
}

// DecodeNative decodes a retained raw reply of any point whose codec is known
// into its exact native value; it does not apply the SI export rule, so a
// withheld point decodes too. An integer of a Boolean point must be 0 or 1, and
// a date must have day 1 to 31, month 1 to 12 and year 0 to 99. The documented
// labels of an enumeration are candidates and a code without one is not shown
// to be invalid, so such a code decodes with an empty label. A version of three
// 0xFF bytes is the empty-slot sentinel. A numeric value is divided by
// ScaleDivisor. It refuses an unknown codec or a zero scale, a reply of the
// wrong width, NaN, infinities, codec sentinels, a code out of range and a text
// that is not UTF-8 with a typed DecodeError. The reply is not modified.
func DecodeNative(point Point, raw []byte) (NativeValue, error) {
	refuse := func(failure DecodeFailure, width int) (NativeValue, error) {
		return NativeValue{}, &DecodeError{DeviceClass: point.DeviceClass, Name: point.Name, Failure: failure, Length: len(raw), Width: width}
	}
	if point.Codec == CodecText {
		text := raw
		if end := bytes.IndexByte(raw, 0); end >= 0 {
			text = raw[:end]
		}
		if !utf8.Valid(text) {
			return refuse(DecodeInvalidText, 0)
		}
		return NativeValue{ValueClass: ValueText, Native: string(text), NativeUnit: point.NativeUnit}, nil
	}
	widths, known := nativeWidths[point.Codec]
	if !known || point.ScaleDivisor == 0 {
		return refuse(DecodeCodecUnknown, 0)
	}
	if len(raw) < widths[0] {
		return refuse(DecodeShortReply, widths[0])
	}
	if len(raw) > widths[1] {
		return refuse(DecodeLongReply, widths[1])
	}
	switch point.Codec {
	case CodecDayMonthYear:
		if raw[0] < 1 || raw[0] > 31 || raw[1] < 1 || raw[1] > 12 || raw[2] > 99 {
			return refuse(DecodeOutOfRange, widths[1])
		}
		return NativeValue{ValueClass: ValueDate, Native: fmt.Sprintf("%02d.%02d.%02d", raw[0], raw[1], raw[2])}, nil
	case CodecHourMinuteSecond:
		return NativeValue{ValueClass: ValueTime, Native: fmt.Sprintf("%02d:%02d:%02d", raw[0], raw[1], raw[2])}, nil
	case CodecMajorMinorPatch:
		if raw[0] == 0xFF && raw[1] == 0xFF && raw[2] == 0xFF {
			return refuse(DecodeSentinel, widths[1])
		}
		return NativeValue{ValueClass: ValueVersion, Native: fmt.Sprintf("%02d.%02d.%02d", raw[0], raw[1], raw[2])}, nil
	case CodecFloat32LE:
		value := math.Float32frombits(binary.LittleEndian.Uint32(raw))
		if math.IsNaN(float64(value)) {
			return refuse(DecodeNotANumber, widths[1])
		}
		if math.IsInf(float64(value), 0) {
			return refuse(DecodeInfinite, widths[1])
		}
		native, _ := new(big.Rat).SetString(strconv.FormatFloat(float64(value), 'f', -1, 32))
		return point.numeric(native), nil
	}
	code, sentinel := nativeInteger(point.Codec, raw)
	if sentinel {
		return refuse(DecodeSentinel, widths[1])
	}
	if point.ValueClass == ValueBoolean {
		if code != 0 && code != 1 {
			return refuse(DecodeOutOfRange, widths[1])
		}
		label, _ := point.label(code)
		return NativeValue{ValueClass: ValueBoolean, Native: strconv.FormatBool(code == 1), NativeUnit: point.NativeUnit, Label: label}, nil
	}
	if len(point.Enumeration) > 0 || point.ValueClass == ValueEnumeration {
		label, _ := point.label(code)
		return NativeValue{ValueClass: ValueEnumeration, Native: strconv.FormatInt(code, 10), NativeUnit: point.NativeUnit, Label: label}, nil
	}
	if point.Codec == CodecData2c || point.Codec == CodecData2b {
		divisor := int64(16)
		if point.Codec == CodecData2b {
			divisor = 256
		}
		return point.numeric(new(big.Rat).SetFrac64(code, divisor)), nil
	}
	return point.numeric(new(big.Rat).SetInt64(code)), nil
}

func nativeInteger(codec Codec, raw []byte) (int64, bool) {
	switch codec {
	case CodecSignedByte:
		return int64(int8(raw[0])), false
	case CodecUint32LE:
		return int64(binary.LittleEndian.Uint32(raw)), false
	case CodecData2c, CodecData2b:
		word := binary.LittleEndian.Uint16(raw)
		return int64(int16(word)), word == 0x8000
	case CodecUIN:
		word := binary.LittleEndian.Uint16(raw)
		return int64(word), word == 0xFFFF
	}
	if len(raw) == 1 {
		return int64(raw[0]), false
	}
	return int64(binary.LittleEndian.Uint16(raw)), false
}

func (point Point) numeric(native *big.Rat) NativeValue {
	native.Quo(native, new(big.Rat).SetInt64(int64(point.ScaleDivisor)))
	return NativeValue{ValueClass: ValueNumeric, Native: exactText(native), NativeUnit: point.NativeUnit}
}

func (point Point) label(code int64) (string, bool) {
	for _, entry := range point.Enumeration {
		if entry.Value == code {
			return entry.Label, true
		}
	}
	return "", false
}
