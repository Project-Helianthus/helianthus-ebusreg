package system

import (
	"bytes"
	"reflect"
	"testing"
)

func TestAb11Families(ab11T *testing.T) {
	ab11Cases := []struct {
		ab11Name  string
		ab11Codec Ab11Codec
		ab11Raw   []byte
		ab11Want  any
	}{
		{"bool false", Ab11BOOL, []byte{0}, false},
		{"bool true", Ab11BOOL, []byte{1}, true},
		{"uch", Ab11UCH, []byte{255}, uint8(255)},
		{"uin endian", Ab11UIN, []byte{0x34, 0x12}, uint16(0x1234)},
		{"i8 minimum", Ab11I8, []byte{0x80}, int8(-128)},
		{"i8 maximum", Ab11I8, []byte{0x7f}, int8(127)},
		{"i8 minus one", Ab11I8, []byte{255}, int8(-1)},
		{"i16 minus one", Ab11I16, []byte{255, 255}, int16(-1)},
		{"i16 minimum", Ab11I16, []byte{0, 0x80}, int16(-32768)},
		{"i16 maximum", Ab11I16, []byte{255, 0x7f}, int16(32767)},
		{"i16 endian", Ab11I16, []byte{0x34, 0x12}, int16(0x1234)},
		{"u32 endian", Ab11U32, []byte{0x78, 0x56, 0x34, 0x12}, uint32(0x12345678)},
		{"u32 maximum", Ab11U32, []byte{255, 255, 255, 255}, uint32(4294967295)},
		{"i32 minimum", Ab11I32, []byte{0, 0, 0, 0x80}, int32(-2147483648)},
		{"i32 maximum", Ab11I32, []byte{255, 255, 255, 0x7f}, int32(2147483647)},
		{"i32 minus one", Ab11I32, []byte{255, 255, 255, 255}, int32(-1)},
		{"i32 endian", Ab11I32, []byte{0x78, 0x56, 0x34, 0x12}, int32(0x12345678)},
		{"exp", Ab11EXP, []byte{0, 0, 0xc0, 0x3f}, float32(1.5)},
		{"exp negative", Ab11EXP, []byte{0, 0, 0xc0, 0xbf}, float32(-1.5)},
		{"date numeric", Ab11DATE, []byte{0x1f, 0x0c, 0x63}, Ab11Date{31, 12, 99}},
		{"date no calendar inference", Ab11DATE, []byte{255, 0, 255}, Ab11Date{255, 0, 255}},
		{"time numeric", Ab11TIME, []byte{0x17, 0x3b, 0x3a}, Ab11Time{23, 59, 58}},
		{"time no constraints", Ab11TIME, []byte{255, 255, 255}, Ab11Time{255, 255, 255}},
		{"fw bcd", Ab11FW, []byte{0x12, 0x34, 0x56}, [3]uint8{12, 34, 56}},
		{"fw limits", Ab11FW, []byte{0, 0x99, 0x90}, [3]uint8{0, 99, 90}},
		{"fwu raw", Ab11FWU, []byte{0x12, 0x34, 0x56}, [3]uint8{18, 52, 86}},
		{"fwu no bcd", Ab11FWU, []byte{255, 0xaa, 0xfa}, [3]uint8{255, 170, 250}},
		{"str exact", Ab11STR, []byte{'A', 0, 0x80, 255, 0}, Ab11String{[]byte{'A', 0, 0x80, 255, 0}, "A\x00\u0080\u00ff\x00"}},
		{"hex exact", Ab11HEX, []byte{0, 0x80, 255}, []byte{0, 0x80, 255}},
	}
	for _, ab11Case := range ab11Cases {
		ab11T.Run(ab11Case.ab11Name, func(ab11T *testing.T) {
			for _, ab11Grade := range []Ab11Grade{Ab11Proven, Ab11Hypothesis} {
				ab11Got := Ab11DecodeValue(ab11Case.ab11Codec, len(ab11Case.ab11Raw), ab11Grade, ab11Case.ab11Raw)
				if !reflect.DeepEqual(ab11Got.Ab11Interpretation, ab11Case.ab11Want) || ab11Got.Ab11Qualified != (ab11Grade == Ab11Proven) || ab11Got.Ab11Grade != ab11Grade || !bytes.Equal(ab11Got.Ab11Raw, ab11Case.ab11Raw) {
					ab11T.Fatalf("grade %s: got %#v, want %#v", ab11Grade, ab11Got, ab11Case.ab11Want)
				}
				ab11Before := append([]byte(nil), ab11Case.ab11Raw...)
				ab11Got.Ab11Raw[0] ^= 255
				if !bytes.Equal(ab11Case.ab11Raw, ab11Before) {
					ab11T.Fatal("successful raw result aliases input")
				}
			}
		})
	}
}

func ab11AssertFallback(ab11T *testing.T, ab11Codec Ab11Codec, ab11Width int, ab11Grade Ab11Grade, ab11Raw []byte) {
	ab11T.Helper()
	ab11Before := append([]byte(nil), ab11Raw...)
	ab11Got := Ab11DecodeValue(ab11Codec, ab11Width, ab11Grade, ab11Raw)
	if ab11Got.Ab11Interpretation != nil || ab11Got.Ab11Qualified || ab11Got.Ab11Grade != ab11Grade || !bytes.Equal(ab11Got.Ab11Raw, ab11Before) {
		ab11T.Fatalf("fallback codec=%s width=%d grade=%s raw=%x: %#v", ab11Codec, ab11Width, ab11Grade, ab11Raw, ab11Got)
	}
	if len(ab11Got.Ab11Raw) > 0 {
		ab11Got.Ab11Raw[0] ^= 255
		if !bytes.Equal(ab11Raw, ab11Before) {
			ab11T.Fatal("fallback aliases input")
		}
	}
}

func TestAb11AdmissionDimensions(ab11T *testing.T) {
	for _, ab11Codec := range []Ab11Codec{Ab11BOOL, Ab11UCH, Ab11UIN, Ab11I8, Ab11I16, Ab11U32, Ab11I32, Ab11EXP, Ab11DATE, Ab11TIME, Ab11FW, Ab11FWU, Ab11STR, Ab11HEX} {
		ab11Width := 3
		switch ab11Codec {
		case Ab11BOOL, Ab11UCH, Ab11I8:
			ab11Width = 1
		case Ab11UIN, Ab11I16:
			ab11Width = 2
		case Ab11U32, Ab11I32, Ab11EXP:
			ab11Width = 4
		}
		ab11T.Run(string(ab11Codec), func(ab11T *testing.T) {
			ab11Raw := make([]byte, ab11Width)
			ab11Control := Ab11DecodeValue(ab11Codec, ab11Width, Ab11Proven, ab11Raw)
			if !ab11Control.Ab11Qualified || ab11Control.Ab11Interpretation == nil {
				ab11T.Fatal("valid control rejected")
			}
			for _, ab11Length := range []int{ab11Width - 1, ab11Width + 1} {
				ab11AssertFallback(ab11T, ab11Codec, ab11Width, Ab11Proven, make([]byte, ab11Length))
			}
			for _, ab11OtherWidth := range []int{-1, 0, ab11Width - 1, ab11Width + 1} {
				ab11AssertFallback(ab11T, ab11Codec, ab11OtherWidth, Ab11Proven, ab11Raw)
				if ab11Codec != Ab11STR && ab11Codec != Ab11HEX && ab11OtherWidth > 0 {
					ab11AssertFallback(ab11T, ab11Codec, ab11OtherWidth, Ab11Proven, make([]byte, ab11OtherWidth))
				}
			}
			for _, ab11Grade := range []Ab11Grade{"", Ab11Unknown, "unrecognized"} {
				ab11AssertFallback(ab11T, ab11Codec, ab11Width, ab11Grade, ab11Raw)
			}
			ab11AssertFallback(ab11T, "unrecognized", ab11Width, Ab11Proven, ab11Raw)
		})
	}
}

func TestAb11InvalidDomains(ab11T *testing.T) {
	for ab11Byte := 2; ab11Byte <= 255; ab11Byte++ {
		ab11AssertFallback(ab11T, Ab11BOOL, 1, Ab11Proven, []byte{byte(ab11Byte)})
	}
	for ab11Index := 0; ab11Index < 3; ab11Index++ {
		for ab11Byte := 0; ab11Byte <= 255; ab11Byte++ {
			if ab11Byte>>4 <= 9 && ab11Byte&15 <= 9 {
				continue
			}
			ab11Raw := []byte{0x12, 0x34, 0x56}
			ab11Raw[ab11Index] = byte(ab11Byte)
			ab11AssertFallback(ab11T, Ab11FW, 3, Ab11Proven, ab11Raw)
		}
	}
	for _, ab11Raw := range [][]byte{{0, 0, 0x80, 0x7f}, {0, 0, 0x80, 0xff}, {0, 0, 0xc0, 0x7f}, {1, 0, 0x80, 0x7f}} {
		for _, ab11Grade := range []Ab11Grade{Ab11Proven, Ab11Hypothesis} {
			ab11AssertFallback(ab11T, Ab11EXP, 4, ab11Grade, ab11Raw)
		}
	}
}

func TestAb11VariableWidthsAndOwnership(ab11T *testing.T) {
	for _, ab11Codec := range []Ab11Codec{Ab11STR, Ab11HEX} {
		for _, ab11Width := range []int{1, 2, 5} {
			ab11Raw := bytes.Repeat([]byte{0xff}, ab11Width)
			ab11Got := Ab11DecodeValue(ab11Codec, ab11Width, Ab11Proven, ab11Raw)
			if !ab11Got.Ab11Qualified {
				ab11T.Fatal("variable width rejected")
			}
			ab11Got.Ab11Raw[0] = 1
			if ab11Raw[0] != 255 {
				ab11T.Fatal("raw aliases input")
			}
			var ab11Bytes []byte
			if ab11Codec == Ab11STR {
				ab11Bytes = ab11Got.Ab11Interpretation.(Ab11String).Ab11Bytes
			} else {
				ab11Bytes = ab11Got.Ab11Interpretation.([]byte)
			}
			if ab11Bytes[0] != 255 {
				ab11T.Fatal("interpretation aliases raw")
			}
			ab11Bytes[0] = 2
			if ab11Raw[0] != 255 {
				ab11T.Fatal("interpretation aliases input")
			}
			ab11Raw[0] = 3
			if ab11Got.Ab11Raw[0] != 1 || ab11Bytes[0] != 2 {
				ab11T.Fatal("input aliases result")
			}
		}
	}
}
