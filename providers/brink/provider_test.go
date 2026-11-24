package brink

import (
	"testing"

	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
)

func TestBrinkProviderRegistersOnRegistry(t *testing.T) {
	var provider registry.PlaneProvider = NewProvider()
	reg := registry.NewDeviceRegistry([]registry.PlaneProvider{provider})
	reg.Register(registry.DeviceInfo{Address: 0x3C, Manufacturer: "0x40", DeviceID: "HRU00", SerialNumber: "SN-B"})
	entry, ok := reg.Lookup(0x3C)
	if !ok {
		t.Fatal("registration produced no entry")
	}
	if len(entry.Planes()) != 0 {
		t.Fatalf("Match is false, yet the entry carries planes %v", entry.Planes())
	}
}

// Match is false for every identity, including one whose address and
// manufacturer the page would name: the page states no family key on DeviceInfo.
func TestBrinkMatchFalseWithoutFamilyKey(t *testing.T) {
	p := NewProvider()
	for _, info := range []registry.DeviceInfo{
		{Address: 0x3C, Manufacturer: "0x40", DeviceID: "HRU00"},
		{Address: 0x3C, Manufacturer: "Brink", DeviceID: "HRU01"},
		{Address: 0x05},
	} {
		if p.Match(info) {
			t.Fatalf("Match(%#v) = true; the page states no family key on DeviceInfo", info)
		}
	}
}

func TestBrinkCreatePlanesCarriesEveryDocumentedMethod(t *testing.T) {
	planes := NewProvider().CreatePlanes(registry.DeviceInfo{Address: 0x21, Manufacturer: "Acme", DeviceID: "SYN00"})
	if len(planes) != 1 {
		t.Fatalf("planes = %d; want one family plane", len(planes))
	}
	if planes[0].Name() != "brink_40xx" {
		t.Fatalf("plane name = %s", planes[0].Name())
	}
	want := map[string]struct {
		readOnly bool
		sb       byte
	}{
		methodSensorRead:         {true, 0x22},
		methodParameterRead:      {true, 0x50},
		methodErrorHistoryRead:   {true, 0x90},
		methodFactoryReset:       {false, 0xFF},
		methodResetNotifications: {false, 0x91},
		methodFanModeTemporary:   {false, 0xA1},
		methodFanModePersistent:  {false, 0xCB},
		methodParameterWrite:     {false, 0x80},
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
		if m.ReadOnly() != w.readOnly {
			t.Fatalf("%s ReadOnly = %v; want %v", name, m.ReadOnly(), w.readOnly)
		}
		if m.Template().Primary() != 0x40 || m.Template().Secondary() != w.sb {
			t.Fatalf("%s template = %02X %02X; want 40 %02X", name, m.Template().Primary(), m.Template().Secondary(), w.sb)
		}
	}
}

// A caller-supplied identity with any address gets the same plane.
func TestBrinkCreatePlanesIndependentOfIdentity(t *testing.T) {
	a := NewProvider().CreatePlanes(registry.DeviceInfo{Address: 0x3C})
	b := NewProvider().CreatePlanes(registry.DeviceInfo{Address: 0x7E, Manufacturer: "0x19"})
	if len(a) != 1 || len(b) != 1 || len(a[0].Methods()) != len(b[0].Methods()) {
		t.Fatal("plane depends on the identity")
	}
}
