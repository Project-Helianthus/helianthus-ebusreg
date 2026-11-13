package system

import (
	"bytes"
	"reflect"
	"testing"
)

func TestAb11Families(valueCodecT *testing.T) {
	valueCodecCases := []struct {
		valueCodecName  string
		valueCodecCodec ValueCodec
		valueCodecRaw   []byte
		valueCodecWant  any
	}{
		{"bool false", ValueCodecBOOL, []byte{0}, false},
		{"bool true", ValueCodecBOOL, []byte{1}, true},
		{"uch", ValueCodecUCH, []byte{255}, uint8(255)},
		{"uin endian", ValueCodecUIN, []byte{0x34, 0x12}, uint16(0x1234)},
		{"i8 minimum", ValueCodecI8, []byte{0x80}, int8(-128)},
		{"i8 maximum", ValueCodecI8, []byte{0x7f}, int8(127)},
		{"i8 minus one", ValueCodecI8, []byte{255}, int8(-1)},
		{"i16 minus one", ValueCodecI16, []byte{255, 255}, int16(-1)},
		{"i16 minimum", ValueCodecI16, []byte{0, 0x80}, int16(-32768)},
		{"i16 maximum", ValueCodecI16, []byte{255, 0x7f}, int16(32767)},
		{"i16 endian", ValueCodecI16, []byte{0x34, 0x12}, int16(0x1234)},
		{"u32 endian", ValueCodecU32, []byte{0x78, 0x56, 0x34, 0x12}, uint32(0x12345678)},
		{"u32 maximum", ValueCodecU32, []byte{255, 255, 255, 255}, uint32(4294967295)},
		{"i32 minimum", ValueCodecI32, []byte{0, 0, 0, 0x80}, int32(-2147483648)},
		{"i32 maximum", ValueCodecI32, []byte{255, 255, 255, 0x7f}, int32(2147483647)},
		{"i32 minus one", ValueCodecI32, []byte{255, 255, 255, 255}, int32(-1)},
		{"i32 endian", ValueCodecI32, []byte{0x78, 0x56, 0x34, 0x12}, int32(0x12345678)},
		{"exp", ValueCodecEXP, []byte{0, 0, 0xc0, 0x3f}, float32(1.5)},
		{"exp negative", ValueCodecEXP, []byte{0, 0, 0xc0, 0xbf}, float32(-1.5)},
		{"date numeric", ValueCodecDATE, []byte{0x1f, 0x0c, 0x63}, ValueCodecDate{31, 12, 99}},
		{"date no calendar inference", ValueCodecDATE, []byte{255, 0, 255}, ValueCodecDate{255, 0, 255}},
		{"time numeric", ValueCodecTIME, []byte{0x17, 0x3b, 0x3a}, ValueCodecTime{23, 59, 58}},
		{"time no constraints", ValueCodecTIME, []byte{255, 255, 255}, ValueCodecTime{255, 255, 255}},
		{"fw bcd", ValueCodecFW, []byte{0x12, 0x34, 0x56}, [3]uint8{12, 34, 56}},
		{"fw limits", ValueCodecFW, []byte{0, 0x99, 0x90}, [3]uint8{0, 99, 90}},
		{"fwu raw", ValueCodecFWU, []byte{0x12, 0x34, 0x56}, [3]uint8{18, 52, 86}},
		{"fwu no bcd", ValueCodecFWU, []byte{255, 0xaa, 0xfa}, [3]uint8{255, 170, 250}},
		{"str exact", ValueCodecSTR, []byte{'A', 0, 0x80, 255, 0}, ValueCodecString{[]byte{'A', 0, 0x80, 255, 0}, "A\x00\u0080\u00ff\x00"}},
		{"hex exact", ValueCodecHEX, []byte{0, 0x80, 255}, []byte{0, 0x80, 255}},
	}
	for _, valueCodecCase := range valueCodecCases {
		valueCodecT.Run(valueCodecCase.valueCodecName, func(valueCodecT *testing.T) {
			for _, valueCodecGrade := range []ValueCodecGrade{ValueCodecProven, ValueCodecHypothesis} {
				valueCodecGot := DecodeValue(valueCodecCase.valueCodecCodec, len(valueCodecCase.valueCodecRaw), valueCodecGrade, valueCodecCase.valueCodecRaw)
				if !reflect.DeepEqual(valueCodecGot.ValueCodecInterpretation, valueCodecCase.valueCodecWant) || valueCodecGot.ValueCodecQualified != (valueCodecGrade == ValueCodecProven) || valueCodecGot.ValueCodecGrade != valueCodecGrade || !bytes.Equal(valueCodecGot.ValueCodecRaw, valueCodecCase.valueCodecRaw) {
					valueCodecT.Fatalf("grade %s: got %#v, want %#v", valueCodecGrade, valueCodecGot, valueCodecCase.valueCodecWant)
				}
				valueCodecBefore := append([]byte(nil), valueCodecCase.valueCodecRaw...)
				valueCodecGot.ValueCodecRaw[0] ^= 255
				if !bytes.Equal(valueCodecCase.valueCodecRaw, valueCodecBefore) {
					valueCodecT.Fatal("successful raw result aliases input")
				}
			}
		})
	}
}

func valueCodecAssertFallback(valueCodecT *testing.T, valueCodecCodec ValueCodec, valueCodecWidth int, valueCodecGrade ValueCodecGrade, valueCodecRaw []byte) {
	valueCodecT.Helper()
	valueCodecBefore := append([]byte(nil), valueCodecRaw...)
	valueCodecGot := DecodeValue(valueCodecCodec, valueCodecWidth, valueCodecGrade, valueCodecRaw)
	if valueCodecGot.ValueCodecInterpretation != nil || valueCodecGot.ValueCodecQualified || valueCodecGot.ValueCodecGrade != valueCodecGrade || !bytes.Equal(valueCodecGot.ValueCodecRaw, valueCodecBefore) {
		valueCodecT.Fatalf("fallback codec=%s width=%d grade=%s raw=%x: %#v", valueCodecCodec, valueCodecWidth, valueCodecGrade, valueCodecRaw, valueCodecGot)
	}
	if len(valueCodecGot.ValueCodecRaw) > 0 {
		valueCodecGot.ValueCodecRaw[0] ^= 255
		if !bytes.Equal(valueCodecRaw, valueCodecBefore) {
			valueCodecT.Fatal("fallback aliases input")
		}
	}
}

func TestAb11AdmissionDimensions(valueCodecT *testing.T) {
	for _, valueCodecCodec := range []ValueCodec{ValueCodecBOOL, ValueCodecUCH, ValueCodecUIN, ValueCodecI8, ValueCodecI16, ValueCodecU32, ValueCodecI32, ValueCodecEXP, ValueCodecDATE, ValueCodecTIME, ValueCodecFW, ValueCodecFWU, ValueCodecSTR, ValueCodecHEX} {
		valueCodecWidth := 3
		switch valueCodecCodec {
		case ValueCodecBOOL, ValueCodecUCH, ValueCodecI8:
			valueCodecWidth = 1
		case ValueCodecUIN, ValueCodecI16:
			valueCodecWidth = 2
		case ValueCodecU32, ValueCodecI32, ValueCodecEXP:
			valueCodecWidth = 4
		}
		valueCodecT.Run(string(valueCodecCodec), func(valueCodecT *testing.T) {
			valueCodecRaw := make([]byte, valueCodecWidth)
			valueCodecControl := DecodeValue(valueCodecCodec, valueCodecWidth, ValueCodecProven, valueCodecRaw)
			if !valueCodecControl.ValueCodecQualified || valueCodecControl.ValueCodecInterpretation == nil {
				valueCodecT.Fatal("valid control rejected")
			}
			for _, valueCodecLength := range []int{valueCodecWidth - 1, valueCodecWidth + 1} {
				valueCodecAssertFallback(valueCodecT, valueCodecCodec, valueCodecWidth, ValueCodecProven, make([]byte, valueCodecLength))
			}
			for _, valueCodecOtherWidth := range []int{-1, 0, valueCodecWidth - 1, valueCodecWidth + 1} {
				valueCodecAssertFallback(valueCodecT, valueCodecCodec, valueCodecOtherWidth, ValueCodecProven, valueCodecRaw)
				if valueCodecCodec != ValueCodecSTR && valueCodecCodec != ValueCodecHEX && valueCodecOtherWidth > 0 {
					valueCodecAssertFallback(valueCodecT, valueCodecCodec, valueCodecOtherWidth, ValueCodecProven, make([]byte, valueCodecOtherWidth))
				}
			}
			for _, valueCodecGrade := range []ValueCodecGrade{"", ValueCodecUnknown, "unrecognized"} {
				valueCodecAssertFallback(valueCodecT, valueCodecCodec, valueCodecWidth, valueCodecGrade, valueCodecRaw)
			}
			valueCodecAssertFallback(valueCodecT, "unrecognized", valueCodecWidth, ValueCodecProven, valueCodecRaw)
		})
	}
}

func TestAb11InvalidDomains(valueCodecT *testing.T) {
	for valueCodecByte := 2; valueCodecByte <= 255; valueCodecByte++ {
		valueCodecAssertFallback(valueCodecT, ValueCodecBOOL, 1, ValueCodecProven, []byte{byte(valueCodecByte)})
	}
	for valueCodecIndex := 0; valueCodecIndex < 3; valueCodecIndex++ {
		for valueCodecByte := 0; valueCodecByte <= 255; valueCodecByte++ {
			if valueCodecByte>>4 <= 9 && valueCodecByte&15 <= 9 {
				continue
			}
			valueCodecRaw := []byte{0x12, 0x34, 0x56}
			valueCodecRaw[valueCodecIndex] = byte(valueCodecByte)
			valueCodecAssertFallback(valueCodecT, ValueCodecFW, 3, ValueCodecProven, valueCodecRaw)
		}
	}
	for _, valueCodecGrade := range []ValueCodecGrade{ValueCodecProven, ValueCodecHypothesis} {
		valueCodecAssertFallback(valueCodecT, ValueCodecFWU, 3, valueCodecGrade, []byte{255, 255, 255})
	}
	for _, valueCodecRaw := range [][]byte{{0, 0, 0x80, 0x7f}, {0, 0, 0x80, 0xff}, {0, 0, 0xc0, 0x7f}, {1, 0, 0x80, 0x7f}} {
		for _, valueCodecGrade := range []ValueCodecGrade{ValueCodecProven, ValueCodecHypothesis} {
			valueCodecAssertFallback(valueCodecT, ValueCodecEXP, 4, valueCodecGrade, valueCodecRaw)
		}
	}
}

func TestAb11VariableWidthsAndOwnership(valueCodecT *testing.T) {
	for _, valueCodecCodec := range []ValueCodec{ValueCodecSTR, ValueCodecHEX} {
		for _, valueCodecWidth := range []int{1, 2, 5} {
			valueCodecRaw := bytes.Repeat([]byte{0xff}, valueCodecWidth)
			valueCodecGot := DecodeValue(valueCodecCodec, valueCodecWidth, ValueCodecProven, valueCodecRaw)
			if !valueCodecGot.ValueCodecQualified {
				valueCodecT.Fatal("variable width rejected")
			}
			valueCodecGot.ValueCodecRaw[0] = 1
			if valueCodecRaw[0] != 255 {
				valueCodecT.Fatal("raw aliases input")
			}
			var valueCodecBytes []byte
			if valueCodecCodec == ValueCodecSTR {
				valueCodecBytes = valueCodecGot.ValueCodecInterpretation.(ValueCodecString).ValueCodecBytes
			} else {
				valueCodecBytes = valueCodecGot.ValueCodecInterpretation.([]byte)
			}
			if valueCodecBytes[0] != 255 {
				valueCodecT.Fatal("interpretation aliases raw")
			}
			valueCodecBytes[0] = 2
			if valueCodecRaw[0] != 255 {
				valueCodecT.Fatal("interpretation aliases input")
			}
			valueCodecRaw[0] = 3
			if valueCodecGot.ValueCodecRaw[0] != 1 || valueCodecBytes[0] != 2 {
				valueCodecT.Fatal("input aliases result")
			}
		}
	}
}
