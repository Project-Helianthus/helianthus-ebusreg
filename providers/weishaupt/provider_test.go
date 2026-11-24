package weishaupt

import (
	"testing"

	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
)

func TestWeishauptMatchesStatedManufacturerBytes(t *testing.T) {
	p := NewProvider()
	for _, m := range []string{"0x10", "0x50", "0xC5"} {
		if !p.Match(registry.DeviceInfo{Address: 0x08, Manufacturer: m}) {
			t.Fatalf("manufacturer %s did not match", m)
		}
	}
	for _, m := range []string{"0x19", "Vaillant", ""} {
		if p.Match(registry.DeviceInfo{Address: 0x08, Manufacturer: m}) {
			t.Fatalf("manufacturer %q matched the Weishaupt provider", m)
		}
	}
}

func TestWeishauptPlaneMethodsAndFlags(t *testing.T) {
	reg := registry.NewDeviceRegistry([]registry.PlaneProvider{NewProvider()})
	reg.Register(registry.DeviceInfo{Address: 0xF7, Manufacturer: "0x10", DeviceID: "WRSol"})
	entry, ok := reg.Lookup(0xF7)
	if !ok || len(entry.Planes()) != 1 {
		t.Fatalf("entry planes = %v, ok %v", entry, ok)
	}
	want := map[string]struct {
		readOnly  bool
		primary   byte
		secondary byte
	}{
		methodIdentificationRead: {true, 0x07, 0x04},
		methodRAMRead:            {true, 0x09, 0x00},
		methodDateTimeSet:        {false, 0x07, 0x00},
		methodMemoryRead:         {true, 0x50, 0x00},
		methodMemoryWrite:        {false, 0x50, 0x01},
		methodSecondaryRead:      {true, 0x09, 0x02},
		methodSecondaryWrite:     {false, 0x09, 0x03},
	}
	got := map[string]registry.Method{}
	for _, m := range entry.Planes()[0].Methods() {
		got[m.Name()] = m
	}
	if len(got) != len(want) {
		t.Fatalf("methods = %d; want %d", len(got), len(want))
	}
	for name, w := range want {
		m, ok := got[name]
		if !ok || m.ReadOnly() != w.readOnly || m.Template().Primary() != w.primary || m.Template().Secondary() != w.secondary {
			t.Fatalf("%s = %#v; want readOnly %v template %02X %02X", name, m, w.readOnly, w.primary, w.secondary)
		}
	}
}
