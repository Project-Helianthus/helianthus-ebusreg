package weishaupt

import (
	"bytes"
	"errors"
	"math"
	"testing"
)

func TestWeishauptWRSolIdentificationDecodesDocumentedFields(t *testing.T) {
	// Manufacturer 10h, the 2.40 device ID "WRSol", versions BCD 22.12 hardware 11.12.
	id, err := DecodeIdentification([]byte{0x10, 'W', 'R', 'S', 'o', 'l', 0x24, 0x22, 0x11, 0x12})
	if err != nil {
		t.Fatal(err)
	}
	if id.Manufacturer != 0x10 || id.DeviceID != "WRSol" || id.SoftwareVersion != 24 || id.SoftwareRev != 22 || id.HardwareVersion != 11 || id.HardwareRev != 12 {
		t.Fatalf("identification = %#v", id)
	}
	if _, err := DecodeIdentification([]byte{0x19, 'W', 'R', 'S', 'o', 'l', 0x24, 0x22, 0x11, 0x12}); !errors.Is(err, ErrManuf) {
		t.Fatalf("other manufacturer err = %v; want ErrManuf", err)
	}
	if _, err := DecodeIdentification([]byte{0x10, 'W', 'R', 'S', 'o', 'l', 0x2A, 0x22, 0x11, 0x12}); !errors.Is(err, ErrBCD) {
		t.Fatalf("invalid BCD err = %v; want ErrBCD", err)
	}
	if _, err := DecodeIdentification([]byte{0x10, 'W'}); !errors.Is(err, ErrWidth) {
		t.Fatalf("short identification err = %v; want ErrWidth", err)
	}
}

func TestWeishauptWRSolDateTimeEncodesBCDAndOutdoorTemperature(t *testing.T) {
	// Outdoor 20.0 degrees C is 0x1400 (DATA2b, 1/256). Time 12:34:56 on Monday 2 Feb 2026.
	got, err := EncodeDateTime(20.0, DateTime{Seconds: 56, Minutes: 34, Hours: 12, Day: 2, Month: 2, Weekday: 1, Year: 26})
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{0x00, 0x14, 0x56, 0x34, 0x12, 0x02, 0x02, 0x01, 0x26}
	if !bytes.Equal(got, want) {
		t.Fatalf("date/time = % X; want % X", got, want)
	}
}

func TestWeishauptWRSolDateTimeRejectsOutOfRange(t *testing.T) {
	ok := DateTime{Seconds: 0, Minutes: 0, Hours: 0, Day: 1, Month: 1, Weekday: 1, Year: 0}
	if _, err := EncodeDateTime(50.5, ok); !errors.Is(err, ErrRange) {
		t.Fatalf("temperature above 50 err = %v; want ErrRange", err)
	}
	bad := ok
	bad.Seconds = 60
	if _, err := EncodeDateTime(0, bad); !errors.Is(err, ErrRange) {
		t.Fatalf("seconds 60 err = %v; want ErrRange", err)
	}
	bad = ok
	bad.Weekday = 8
	if _, err := EncodeDateTime(0, bad); !errors.Is(err, ErrRange) {
		t.Fatalf("weekday 8 err = %v; want ErrRange", err)
	}
}

func TestWeishauptWRSolRAMReadRequestAndScales(t *testing.T) {
	if got := EncodeRAMRead(0xF484); !bytes.Equal(got, []byte{0x84, 0xF4, 0x02}) {
		t.Fatalf("RAM request = % X", got)
	}
	tso, ok := RAMByAddress(0xF484)
	if !ok || tso.Name != "TSO" || tso.Scale != "X10" {
		t.Fatalf("TSO row = %#v", tso)
	}
	// Integer, scale X10: wire 0x00D2 (210) is 21.0 C.
	f, err := DecodeRAMRead(tso, []byte{0xD2, 0x00})
	if err != nil || math.Abs(f.Value.(float64)-21.0) > 1e-9 || f.Unit != "C" {
		t.Fatalf("TSO value = %#v, err %v", f, err)
	}
	// Integer is signed: 0xFFFE is -2 tenths.
	f, _ = DecodeRAMRead(tso, []byte{0xFE, 0xFF})
	if math.Abs(f.Value.(float64)+0.2) > 1e-9 {
		t.Fatalf("signed TSO value = %v", f.Value)
	}
	// Word, scale X1: a solar power value in kW.
	sp, _ := RAMByAddress(0xF536)
	f, _ = DecodeRAMRead(sp, []byte{0x07, 0x00})
	if f.Value.(float64) != 7 || f.Unit != "kW" {
		t.Fatalf("solar power = %#v", f)
	}
	// ErrorBits is a Word with no scale and no unit: 0x8001 is 32769 as stated.
	eb, _ := RAMByAddress(0xF430)
	f, _ = DecodeRAMRead(eb, []byte{0x01, 0x80})
	if f.Value.(float64) != 32769 || f.Unit != "" {
		t.Fatalf("error bits = %#v", f)
	}
	// ActualTime's layout is not a numeric type on the page; it stays raw.
	at, _ := RAMByAddress(0xFC00)
	f, _ = DecodeRAMRead(at, []byte{0x01, 0x80})
	if !bytes.Equal(f.Value.([]byte), []byte{0x01, 0x80}) {
		t.Fatalf("actual time not kept raw: %#v", f)
	}
	if _, err := DecodeRAMRead(tso, []byte{0xD2}); !errors.Is(err, ErrWidth) {
		t.Fatalf("short RAM reply err = %v; want ErrWidth", err)
	}
}

func TestWeishauptWRSolRAMTableCount(t *testing.T) {
	if len(ramTable) != 43 {
		t.Fatalf("RAM rows = %d; the page's address map has 43", len(ramTable))
	}
}

func TestWeishauptWRSolBufferPumpRowDecodes(t *testing.T) {
	reg, ok := RAMByAddress(0xF4FC)
	if !ok {
		t.Fatal("RAMByAddress(0xF4FC) returned no row; the page's address map lists Drehzahl_Puf_WT")
	}
	if reg.Name != "Drehzahl_Puf_WT" || reg.Scale != "X10" || reg.Unit != "%" {
		t.Fatalf("buffer pump row = %#v; want Drehzahl_Puf_WT X10 %%", reg)
	}
	// 0x00FA as a little-endian Integer with X10 scale is 25.0 %.
	f, err := DecodeRAMRead(reg, []byte{0xFA, 0x00})
	if err != nil {
		t.Fatalf("decode buffer pump reply: %v", err)
	}
	if f.Name != "Drehzahl_Puf_WT" || f.Unit != "%" || f.Value.(float64) != 25 {
		t.Fatalf("buffer pump field = %#v; want Drehzahl_Puf_WT 25 %%", f)
	}
}
