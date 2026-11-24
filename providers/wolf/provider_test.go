package wolf

import (
	"testing"

	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
)

func TestWolfMatchesManufacturerByteOnly(t *testing.T) {
	p := NewProvider()
	if !p.Match(registry.DeviceInfo{Address: 0x08, Manufacturer: "0x19"}) {
		t.Fatal("manufacturer 0x19 did not match")
	}
	for _, m := range []string{"0x08", "Vaillant", "0x50", ""} {
		if p.Match(registry.DeviceInfo{Address: 0x08, Manufacturer: m}) {
			t.Fatalf("manufacturer %q matched the Wolf provider", m)
		}
	}
}

func TestWolfProviderRegistersOnRegistry(t *testing.T) {
	reg := registry.NewDeviceRegistry([]registry.PlaneProvider{NewProvider()})
	reg.Register(registry.DeviceInfo{Address: 0x08, Manufacturer: "0x19", DeviceID: "WOLF0"})
	entry, ok := reg.Lookup(0x08)
	if !ok {
		t.Fatal("registration produced no entry")
	}
	planes := entry.Planes()
	if len(planes) != 1 || planes[0].Name() != "wolf" {
		t.Fatalf("planes = %v; want the wolf plane", planes)
	}
}

func TestWolfPlaneMethodsAndReadOnlyFlags(t *testing.T) {
	planes := NewProvider().CreatePlanes(registry.DeviceInfo{Manufacturer: "0x19"})
	want := map[string]struct {
		readOnly  bool
		primary   byte
		secondary byte
	}{
		methodBurnerStatusSend:   {false, 0x05, 0x03},
		methodModulationSend:     {false, 0x05, 0x04},
		methodHeatRequestSend:    {false, 0x05, 0x07},
		methodBoilerSetpointSend: {false, 0x08, 0x00},
		methodSolarStatusSend:    {false, 0x50, 0x14},
		methodSolarPumpSend:      {false, 0x50, 0x17},
		methodSolarPowerSend:     {false, 0x50, 0x18},
		methodServiceRead:        {true, 0x50, 0x22},
		methodServiceWrite:       {false, 0x50, 0x23},
	}
	got := map[string]registry.Method{}
	for _, m := range planes[0].Methods() {
		got[m.Name()] = m
	}
	if len(got) != len(want) {
		t.Fatalf("methods = %d; want %d", len(got), len(want))
	}
	for name, w := range want {
		m, ok := got[name]
		if !ok {
			t.Fatalf("method %s missing", name)
		}
		if m.ReadOnly() != w.readOnly || m.Template().Primary() != w.primary || m.Template().Secondary() != w.secondary {
			t.Fatalf("%s = readOnly %v template %02X %02X; want %v %02X %02X", name, m.ReadOnly(), m.Template().Primary(), m.Template().Secondary(), w.readOnly, w.primary, w.secondary)
		}
	}
}
