package system

import (
	"bytes"
	"encoding/binary"
	"math"
	"testing"

	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
)

func TestSystemPlane_DecodeBroadcast_EnergyStatsSelectorRequest(t *testing.T) {
	t.Parallel()

	planes := NewProvider().CreatePlanes(registry.DeviceInfo{Manufacturer: "Vaillant", Address: 0x08})
	systemPlane := planes[0].(*plane)

	frame := protocol.Frame{
		Primary:   0xB5,
		Secondary: 0x16,
		Data:      []byte{0x10, 0x03, 0xFF, 0xFF, 0x04, 0x04, 0x00, 0x32},
	}

	values, handled, err := systemPlane.DecodeBroadcast(frame)
	if err != nil {
		t.Fatalf("DecodeBroadcast error = %v", err)
	}
	if !handled {
		t.Fatal("expected handled=true")
	}

	if got := values["period"]; !got.Valid || got.Value != "year" {
		t.Fatalf("period = %+v; want year valid", got)
	}
	if got := values["year"]; !got.Valid || got.Value != uint16(2025) {
		t.Fatalf("year = %+v; want 2025 valid", got)
	}
	if got := values["source"]; !got.Valid || got.Value != "gas" {
		t.Fatalf("source = %+v; want gas valid", got)
	}
	if got := values["usage"]; !got.Valid || got.Value != "hot_water" {
		t.Fatalf("usage = %+v; want hot_water valid", got)
	}
	if got := values["month"]; got.Valid {
		t.Fatalf("month = %+v; want invalid", got)
	}
	if got := values["day"]; got.Valid {
		t.Fatalf("day = %+v; want invalid", got)
	}
}

func TestSystemPlane_DecodeBroadcast_EnergyStatsResponseWh(t *testing.T) {
	t.Parallel()

	planes := NewProvider().CreatePlanes(registry.DeviceInfo{Manufacturer: "Vaillant", Address: 0x08})
	systemPlane := planes[0].(*plane)

	wantWh := float32(12345.0)
	whBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(whBytes, math.Float32bits(wantWh))

	payload := append([]byte{0x03, 0x00, 0xFF, 0x04, 0x04, 0x00, 0x32}, whBytes...)
	frame := protocol.Frame{
		Primary:   0xB5,
		Secondary: 0x16,
		Data:      payload,
	}

	values, handled, err := systemPlane.DecodeBroadcast(frame)
	if err != nil {
		t.Fatalf("DecodeBroadcast error = %v", err)
	}
	if !handled {
		t.Fatal("expected handled=true")
	}

	if got := values["wh"]; !got.Valid || got.Value != wantWh {
		t.Fatalf("wh = %+v; want %v valid", got, wantWh)
	}
}

func TestSystemPlane_DecodeBroadcast_EnergyStatsPackedDate(t *testing.T) {
	t.Parallel()

	// Literal wire bytes exercise both offsets independently of an encoder.
	cases := []struct {
		name       string
		period     byte
		lo, hi     byte
		year       uint16
		month, day uint8
		valid      bool
	}{
		{name: "year 2024", period: 3, lo: 0x00, hi: 0x30, year: 2024, valid: true},
		{name: "year 2025", period: 3, lo: 0x00, hi: 0x32, year: 2025, valid: true},
		{name: "year 2026", period: 3, lo: 0x00, hi: 0x34, year: 2026, valid: true},
		{name: "year 2032", period: 3, lo: 0x00, hi: 0x40, year: 2032, valid: true},
		{name: "year 2000", period: 3, lo: 0x00, hi: 0x00, year: 2000, valid: true},
		{name: "year 2127", period: 3, lo: 0x00, hi: 0xFE, year: 2127, valid: true},
		{name: "January 1 2024", period: 1, lo: 0x21, hi: 0x30, year: 2024, month: 1, day: 1, valid: true},
		{name: "December 31 2024", period: 1, lo: 0x9F, hi: 0x31, year: 2024, month: 12, day: 31, valid: true},
		{name: "January 1 2025", period: 1, lo: 0x21, hi: 0x32, year: 2025, month: 1, day: 1, valid: true},
		{name: "December 31 2025", period: 1, lo: 0x9F, hi: 0x33, year: 2025, month: 12, day: 31, valid: true},
		{name: "January 1 2026", period: 1, lo: 0x21, hi: 0x34, year: 2026, month: 1, day: 1, valid: true},
		{name: "August 1 2026", period: 1, lo: 0x01, hi: 0x35, year: 2026, month: 8, day: 1, valid: true},
		{name: "December 31 2026", period: 1, lo: 0x9F, hi: 0x35, year: 2026, month: 12, day: 31, valid: true},
		{name: "January 1 2032", period: 1, lo: 0x21, hi: 0x40, year: 2032, month: 1, day: 1, valid: true},
		{name: "December 31 2127", period: 1, lo: 0x9F, hi: 0xFF, year: 2127, month: 12, day: 31, valid: true},
		{name: "leap day 2024", period: 1, lo: 0x5D, hi: 0x30, year: 2024, month: 2, day: 29, valid: true},
		{name: "leap day 2000", period: 1, lo: 0x5D, hi: 0x00, year: 2000, month: 2, day: 29, valid: true},
		{name: "March 2025", period: 2, lo: 0x60, hi: 0x32, year: 2025, month: 3, valid: true},
		{name: "January 2026", period: 2, lo: 0x20, hi: 0x34, year: 2026, month: 1, valid: true},
		{name: "August 2026", period: 2, lo: 0x00, hi: 0x35, year: 2026, month: 8, valid: true},
		{name: "December 2032", period: 2, lo: 0x80, hi: 0x41, year: 2032, month: 12, valid: true},
		{name: "total zero date", period: 0, lo: 0x00, hi: 0x00},
		{name: "total echoed date", period: 0, lo: 0x87, hi: 0x34},
		{name: "total ignores date", period: 0, lo: 0xFF, hi: 0xFF},
		{name: "daily zero month", period: 1, lo: 0x01, hi: 0x34},
		{name: "daily zero day", period: 1, lo: 0x20, hi: 0x34},
		{name: "daily month 13", period: 1, lo: 0xA1, hi: 0x35},
		{name: "April 31", period: 1, lo: 0x9F, hi: 0x34},
		{name: "non leap day 2026", period: 1, lo: 0x5D, hi: 0x34},
		{name: "non leap century 2100", period: 1, lo: 0x5D, hi: 0xC8},
		{name: "monthly zero month", period: 2, lo: 0x00, hi: 0x34},
		{name: "monthly month 13", period: 2, lo: 0xA0, hi: 0x35},
		{name: "monthly nonzero day", period: 2, lo: 0x21, hi: 0x34},
		{name: "yearly nonzero month", period: 3, lo: 0x20, hi: 0x34},
		{name: "yearly nonzero day", period: 3, lo: 0x01, hi: 0x34},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for _, response := range []bool{false, true} {
				name := "request"
				if response {
					name = "response"
				}
				t.Run(name, func(t *testing.T) {
					systemPlane := newSystemPlane(registry.DeviceInfo{})
					payload := []byte{0x10, tc.period, 0xFF, 0xFF, 0x04, 0x03, tc.lo, tc.hi}
					selectorLen := 8
					if response {
						payload = []byte{tc.period, 0x00, 0xFF, 0x04, 0x03, tc.lo, tc.hi, 0x00, 0x54, 0x5D, 0x46}
						selectorLen = 7
					}
					values, handled, err := systemPlane.DecodeBroadcast(protocol.Frame{Primary: 0xB5, Secondary: 0x16, Data: payload})
					if err != nil || !handled {
						t.Fatalf("DecodeBroadcast = handled %v, error %v", handled, err)
					}
					if _, exists := values["year_kind"]; exists {
						t.Fatal("year_kind must not be published without a reference date")
					}
					for _, field := range []struct {
						name  string
						want  any
						valid bool
					}{
						{"year", tc.year, tc.valid},
						{"month", tc.month, tc.valid && (tc.period == 1 || tc.period == 2)},
						{"day", tc.day, tc.valid && tc.period == 1},
					} {
						got, exists := values[field.name]
						if !exists || got.Valid != field.valid || (field.valid && got.Value != field.want) {
							t.Errorf("%s = %+v; want %v, valid %v", field.name, got, field.want, field.valid)
						}
					}
					if got := values["source"]; !got.Valid || got.Value != "gas" {
						t.Errorf("source = %+v; want gas valid", got)
					}
					if got := values["usage"]; !got.Valid || got.Value != "heating" {
						t.Errorf("usage = %+v; want heating valid", got)
					}
					if response {
						if got := values["wh"]; !got.Valid || got.Value != float32(14165) {
							t.Errorf("wh = %+v; want 14165 valid", got)
						}
					} else if _, exists := values["wh"]; exists {
						t.Error("requests must not publish wh")
					}
					if got := values["selector_is_rsp"]; !got.Valid || got.Value != response {
						t.Errorf("selector_is_rsp = %+v; want %v valid", got, response)
					}
					for _, field := range []struct {
						name string
						want uint8
					}{{"selector_w", tc.lo >> 4}, {"selector_v", tc.lo & 0x0F}, {"selector_q", tc.hi & 0x0F}} {
						if got := values[field.name]; !got.Valid || got.Value != field.want {
							t.Errorf("%s = %+v; want %v valid", field.name, got, field.want)
						}
					}
					// Preserve full wire bytes and ownership after the caller reuses its buffer.
					wantPayload := append([]byte(nil), payload...)
					payload[0] ^= 0xFF
					for _, field := range []struct {
						name string
						want []byte
					}{{"payload", wantPayload}, {"selector", wantPayload[:selectorLen]}} {
						got := values[field.name]
						if raw, ok := got.Value.([]byte); !ok || !got.Valid || !bytes.Equal(raw, field.want) {
							t.Errorf("%s = %+v; want %x valid", field.name, got, field.want)
						}
					}
				})
			}
		})
	}
}

func TestSystemPlane_DecodeBroadcast_EnergyStatsDocumentedTotalReply(t *testing.T) {
	t.Parallel()

	// docs-ebus B516 section 6.1, Query 1: 2026-04-07, 14165 Wh.
	payload := []byte{0x00, 0x00, 0x00, 0x03, 0x03, 0x87, 0x34, 0x00, 0x54, 0x5D, 0x46}
	values, handled, err := newSystemPlane(registry.DeviceInfo{}).DecodeBroadcast(protocol.Frame{Primary: 0xB5, Secondary: 0x16, Data: payload})
	if err != nil || !handled {
		t.Fatalf("DecodeBroadcast = handled %v, error %v", handled, err)
	}
	if got := values["wh"]; !got.Valid || got.Value != float32(14165) {
		t.Fatalf("wh = %+v; want 14165 valid", got)
	}
	if got := values["period"]; !got.Valid || got.Value != "all" {
		t.Fatalf("period = %+v; want all valid", got)
	}
	for _, name := range []string{"year", "month", "day"} {
		if values[name].Valid {
			t.Errorf("%s must not identify a window for a total reply", name)
		}
	}
}

func TestSystemPlane_DecodeBroadcast_EnergyStatsRejectedFrames(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		primary byte
		payload []byte
	}{
		{"wrong service", 0xB4, []byte{0x10, 0x03, 0xFF, 0xFF, 0x04, 0x03, 0x00, 0x34}},
		{"wrong subcommand", 0xB5, []byte{0x11, 0x03, 0xFF, 0xFF, 0x04, 0x03, 0x00, 0x34}},
		{"indexed period", 0xB5, []byte{0x10, 0x03, 0x00, 0xFF, 0x04, 0x03, 0x00, 0x34}},
		{"indexed energy", 0xB5, []byte{0x10, 0x03, 0xFF, 0x00, 0x04, 0x03, 0x00, 0x34}},
		{"not OK total", 0xB5, []byte{0x10, 0x00, 0xFF, 0x03, 0x04, 0x87, 0x34, 0x00, 0x00, 0x00, 0x00}},
		{"not OK year", 0xB5, []byte{0x13, 0xFF, 0xFF, 0x03, 0x03, 0x00, 0x32, 0x00, 0x00, 0x00, 0x00}},
		{"access flag", 0xB5, []byte{0x04, 0x00, 0xFF, 0x04, 0x03, 0x00, 0x34, 0x00, 0x00, 0x00, 0x00}},
		{"truncated", 0xB5, []byte{0x10, 0x03, 0xFF, 0xFF, 0x04, 0x03, 0x00}},
		{"empty", 0xB5, nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			values, handled, err := newSystemPlane(registry.DeviceInfo{}).DecodeBroadcast(protocol.Frame{Primary: tc.primary, Secondary: 0x16, Data: tc.payload})
			if err != nil || handled || values != nil {
				t.Fatalf("DecodeBroadcast = %v, handled %v, error %v; want nil, false, nil", values, handled, err)
			}
		})
	}
}
