package system

import (
	"encoding/binary"
	"math"
	"testing"

	"github.com/d3vi1/helianthus-ebusgo/protocol"
	"github.com/d3vi1/helianthus-ebusreg/registry"
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
	if got := values["year_kind"]; !got.Valid || got.Value != "current" {
		t.Fatalf("year_kind = %+v; want current valid", got)
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
