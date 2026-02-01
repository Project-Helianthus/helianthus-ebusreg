package registry

import "testing"

type mockTemplate struct {
	primary   byte
	secondary byte
}

func (m mockTemplate) Primary() byte {
	return m.primary
}

func (m mockTemplate) Secondary() byte {
	return m.secondary
}

type mockMethod struct {
	name     string
	readOnly bool
	template FrameTemplate
}

func (m mockMethod) Name() string {
	return m.name
}

func (m mockMethod) ReadOnly() bool {
	return m.readOnly
}

func (m mockMethod) Template() FrameTemplate {
	return m.template
}

type mockPlane struct {
	name    string
	methods []Method
}

func (p mockPlane) Name() string {
	return p.name
}

func (p mockPlane) Methods() []Method {
	return p.methods
}

type mockProvider struct {
	name   string
	match  func(DeviceInfo) bool
	planes []Plane
}

func (p mockProvider) Name() string {
	return p.name
}

func (p mockProvider) Match(info DeviceInfo) bool {
	return p.match(info)
}

func (p mockProvider) CreatePlanes(info DeviceInfo) []Plane {
	return p.planes
}

func TestDeviceRegistry_RegisterLookupIterate(t *testing.T) {
	planeA := mockPlane{
		name: "heating",
		methods: []Method{
			mockMethod{
				name:     "get_status",
				readOnly: true,
				template: mockTemplate{primary: 0xB5, secondary: 0x04},
			},
		},
	}
	planeB := mockPlane{name: "dhw"}

	providerA := mockProvider{
		name:  "vaillant",
		match: func(info DeviceInfo) bool { return info.Manufacturer == "vaillant" },
		planes: []Plane{
			planeA,
		},
	}
	providerB := mockProvider{
		name:  "universal",
		match: func(info DeviceInfo) bool { return info.Address == 0x08 },
		planes: []Plane{
			planeB,
		},
	}

	registry := NewDeviceRegistry([]PlaneProvider{providerA, providerB})

	info1 := DeviceInfo{
		Address:         0x08,
		Manufacturer:    "vaillant",
		DeviceID:        "bai",
		SoftwareVersion: "1.0",
		HardwareVersion: "2.0",
	}
	info2 := DeviceInfo{
		Address:         0x10,
		Manufacturer:    "other",
		DeviceID:        "vrc",
		SoftwareVersion: "3.0",
		HardwareVersion: "4.0",
	}

	registry.Register(info1)
	registry.Register(info2)

	entry, ok := registry.Lookup(0x08)
	if !ok {
		t.Fatalf("expected device 0x08 to be registered")
	}
	if entry.Manufacturer() != "vaillant" || entry.DeviceID() != "bai" {
		t.Fatalf("unexpected device info: %+v", entry)
	}
	planes := entry.Planes()
	if len(planes) != 2 {
		t.Fatalf("expected 2 planes, got %d", len(planes))
	}
	if planes[0].Name() != "heating" || planes[1].Name() != "dhw" {
		t.Fatalf("unexpected planes order: %s, %s", planes[0].Name(), planes[1].Name())
	}

	info1Updated := info1
	info1Updated.Manufacturer = "vaillant-updated"
	registry.Register(info1Updated)

	entry, ok = registry.Lookup(0x08)
	if !ok || entry.Manufacturer() != "vaillant-updated" {
		t.Fatalf("expected updated device info")
	}

	addresses := make([]byte, 0)
	registry.Iterate(func(entry DeviceEntry) bool {
		addresses = append(addresses, entry.Address())
		return true
	})
	if len(addresses) != 2 {
		t.Fatalf("expected 2 devices, got %d", len(addresses))
	}
	if addresses[0] != 0x08 || addresses[1] != 0x10 {
		t.Fatalf("unexpected iteration order: %v", addresses)
	}
}
