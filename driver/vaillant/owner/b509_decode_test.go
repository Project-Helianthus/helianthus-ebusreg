package owner

import (
	"bytes"
	"math"
	"testing"
)

func TestVaillantOwnerDecodeData2c(t *testing.T) {
	tests := []struct {
		name string
		in   []byte
		want float64
		ok   bool
	}{
		{name: "normal", in: []byte{0x20, 0x03}, want: 50.0, ok: true},
		{name: "negative", in: []byte{0xF0, 0xFF}, want: -1.0, ok: true},
		{name: "short", in: []byte{0x20}, want: 0, ok: false},
	}
	for _, tc := range tests {
		got, ok := decodeData2c(tc.in)
		if ok != tc.ok {
			t.Fatalf("%s: ok=%v want %v", tc.name, ok, tc.ok)
		}
		if tc.ok && math.Abs(got-tc.want) > 1e-9 {
			t.Fatalf("%s: got=%v want=%v", tc.name, got, tc.want)
		}
	}
}

func TestVaillantOwnerDecodeData2b(t *testing.T) {
	tests := []struct {
		name string
		in   []byte
		want float64
		ok   bool
	}{
		{name: "normal", in: []byte{0x00, 0x01}, want: 1.0, ok: true},
		{name: "boundary", in: []byte{0xFF, 0xFF}, want: 65535.0 / 256.0, ok: true},
		{name: "short", in: []byte{}, want: 0, ok: false},
	}
	for _, tc := range tests {
		got, ok := decodeData2b(tc.in)
		if ok != tc.ok {
			t.Fatalf("%s: ok=%v want %v", tc.name, ok, tc.ok)
		}
		if tc.ok && math.Abs(got-tc.want) > 1e-9 {
			t.Fatalf("%s: got=%v want=%v", tc.name, got, tc.want)
		}
	}
}

func TestVaillantOwnerDecodeUCH(t *testing.T) {
	tests := []struct {
		name string
		in   []byte
		want uint8
		ok   bool
	}{
		{name: "normal", in: []byte{0x7F}, want: 0x7F, ok: true},
		{name: "zero", in: []byte{0x00}, want: 0x00, ok: true},
		{name: "short", in: []byte{}, want: 0, ok: false},
	}
	for _, tc := range tests {
		got, ok := decodeUCH(tc.in)
		if ok != tc.ok || got != tc.want {
			t.Fatalf("%s: got=(%v,%v) want=(%v,%v)", tc.name, got, ok, tc.want, tc.ok)
		}
	}
}

func TestVaillantOwnerDecodeUIN(t *testing.T) {
	tests := []struct {
		name string
		in   []byte
		want uint16
		ok   bool
	}{
		{name: "normal", in: []byte{0x34, 0x12}, want: 0x1234, ok: true},
		{name: "zero", in: []byte{0x00, 0x00}, want: 0x0000, ok: true},
		{name: "short", in: []byte{0x34}, want: 0, ok: false},
	}
	for _, tc := range tests {
		got, ok := decodeUIN(tc.in)
		if ok != tc.ok || got != tc.want {
			t.Fatalf("%s: got=(%v,%v) want=(%v,%v)", tc.name, got, ok, tc.want, tc.ok)
		}
	}
}

func TestVaillantOwnerDecodeSIN(t *testing.T) {
	tests := []struct {
		name string
		in   []byte
		want int16
		ok   bool
	}{
		{name: "positive", in: []byte{0x34, 0x12}, want: 0x1234, ok: true},
		{name: "negative", in: []byte{0xCC, 0xFF}, want: -52, ok: true},
		{name: "short", in: []byte{0xCC}, want: 0, ok: false},
	}
	for _, tc := range tests {
		got, ok := decodeSIN(tc.in)
		if ok != tc.ok || got != tc.want {
			t.Fatalf("%s: got=(%v,%v) want=(%v,%v)", tc.name, got, ok, tc.want, tc.ok)
		}
	}
}

func TestVaillantOwnerDecodeOnOff(t *testing.T) {
	tests := []struct {
		name string
		in   []byte
		want bool
		ok   bool
	}{
		{name: "off", in: []byte{0xF0}, want: false, ok: true},
		{name: "on", in: []byte{0x00}, want: true, ok: true},
		{name: "short", in: []byte{}, want: false, ok: false},
	}
	for _, tc := range tests {
		got, ok := decodeOnOff(tc.in)
		if ok != tc.ok || got != tc.want {
			t.Fatalf("%s: got=(%v,%v) want=(%v,%v)", tc.name, got, ok, tc.want, tc.ok)
		}
	}
}

func TestVaillantOwnerDecodeHoursum2(t *testing.T) {
	tests := []struct {
		name string
		in   []byte
		want float64
		ok   bool
	}{
		{name: "normal", in: []byte{0x05, 0x00}, want: 10.0, ok: true},
		{name: "boundary", in: []byte{0xFF, 0xFF}, want: 131070.0, ok: true},
		{name: "short", in: []byte{0x05}, want: 0, ok: false},
	}
	for _, tc := range tests {
		got, ok := decodeHoursum2(tc.in)
		if ok != tc.ok {
			t.Fatalf("%s: ok=%v want %v", tc.name, ok, tc.ok)
		}
		if tc.ok && math.Abs(got-tc.want) > 1e-9 {
			t.Fatalf("%s: got=%v want=%v", tc.name, got, tc.want)
		}
	}
}

func TestVaillantOwnerDecodeUnsignedHundredths(t *testing.T) {
	tests := []struct {
		name string
		in   []byte
		want float64
		ok   bool
	}{
		{name: "normal", in: []byte{0x10, 0x27}, want: 100.0, ok: true},
		{name: "zero", in: []byte{0x00, 0x00}, want: 0.0, ok: true},
		{name: "short", in: []byte{0x10}, want: 0, ok: false},
	}
	for _, tc := range tests {
		got, ok := decodeUnsignedHundredths(tc.in)
		if ok != tc.ok {
			t.Fatalf("%s: ok=%v want %v", tc.name, ok, tc.ok)
		}
		if tc.ok && math.Abs(got-tc.want) > 1e-9 {
			t.Fatalf("%s: got=%v want=%v", tc.name, got, tc.want)
		}
	}
}

func TestVaillantOwnerDecodePercent0(t *testing.T) {
	tests := []struct {
		name string
		in   []byte
		want uint8
		ok   bool
	}{
		{name: "normal", in: []byte{50}, want: 50, ok: true},
		{name: "boundary", in: []byte{100}, want: 100, ok: true},
		{name: "short", in: []byte{}, want: 0, ok: false},
	}
	for _, tc := range tests {
		got, ok := decodePercent0(tc.in)
		if ok != tc.ok || got != tc.want {
			t.Fatalf("%s: got=(%v,%v) want=(%v,%v)", tc.name, got, ok, tc.want, tc.ok)
		}
	}
}

func TestVaillantOwnerEncodeTempData2c(t *testing.T) {
	tests := []struct {
		name string
		in   float64
		want []byte
	}{
		{name: "positive", in: 50.0, want: []byte{0x20, 0x03}},
		{name: "negative", in: -1.0, want: []byte{0xF0, 0xFF}},
		{name: "zero", in: 0.0, want: []byte{0x00, 0x00}},
	}
	for _, tc := range tests {
		got := encodeTempData2c(tc.in)
		if !bytes.Equal(got, tc.want) {
			t.Fatalf("%s: got=%v want=%v", tc.name, got, tc.want)
		}
	}
}

func TestVaillantOwnerEncodeUCH(t *testing.T) {
	tests := []struct {
		name string
		in   float64
		want []byte
		ok   bool
	}{
		{name: "whole number", in: 18, want: []byte{0x12}, ok: true},
		{name: "zero", in: 0, want: []byte{0x00}, ok: true},
		{name: "fraction rejected", in: 18.5, ok: false},
		{name: "negative rejected", in: -1, ok: false},
		{name: "overflow rejected", in: 256, ok: false},
	}
	for _, tc := range tests {
		got, ok := encodeUCH(tc.in)
		if ok != tc.ok {
			t.Fatalf("%s: ok=%v want %v", tc.name, ok, tc.ok)
		}
		if tc.ok && !bytes.Equal(got, tc.want) {
			t.Fatalf("%s: got=%v want=%v", tc.name, got, tc.want)
		}
	}
}
