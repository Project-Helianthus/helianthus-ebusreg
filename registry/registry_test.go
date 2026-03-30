package registry

import (
	"slices"
	"testing"

	"github.com/Project-Helianthus/helianthus-ebusreg/schema"
)

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
	response schema.SchemaSelector
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

func (m mockMethod) ResponseSchema() schema.SchemaSelector {
	return m.response
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

type mockProjectionProvider struct {
	mockProvider
	projections []Projection
}

func (p mockProjectionProvider) CreateProjections(info DeviceInfo, planes []Plane) []Projection {
	return p.projections
}

type countingProvider struct {
	name        string
	matchFn     func(DeviceInfo) bool
	createFn    func(DeviceInfo) []Plane
	matchCalls  int
	createCalls int
}

func (p *countingProvider) Name() string {
	return p.name
}

func (p *countingProvider) Match(info DeviceInfo) bool {
	p.matchCalls++
	return p.matchFn(info)
}

func (p *countingProvider) CreatePlanes(info DeviceInfo) []Plane {
	p.createCalls++
	return p.createFn(info)
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

func TestDeviceRegistry_Projections(t *testing.T) {
	path := ProjectionPath{
		Plane: ServicePlane,
		Segments: []PathSegment{
			{Name: "devices"},
			{Name: "boiler"},
		},
	}
	node, err := NewNode(path, path)
	if err != nil {
		t.Fatalf("unexpected node error: %v", err)
	}
	projection, err := NewProjection(ServicePlane, []Node{node}, nil)
	if err != nil {
		t.Fatalf("unexpected projection error: %v", err)
	}

	provider := mockProjectionProvider{
		mockProvider: mockProvider{
			name:  "projection",
			match: func(info DeviceInfo) bool { return info.Address == 0x08 },
			planes: []Plane{
				mockPlane{name: "heating"},
			},
		},
		projections: []Projection{projection},
	}

	registry := NewDeviceRegistry([]PlaneProvider{provider})
	entry := registry.Register(DeviceInfo{Address: 0x08})

	projections := entry.Projections()
	if len(projections) != 1 {
		t.Fatalf("expected 1 projection, got %d", len(projections))
	}
	if projections[0].Plane != ServicePlane {
		t.Fatalf("unexpected projection plane: %s", projections[0].Plane)
	}
}

func TestDeviceRegistry_IterateStops(t *testing.T) {
	registry := NewDeviceRegistry(nil)
	registry.Register(DeviceInfo{Address: 0x08})
	registry.Register(DeviceInfo{Address: 0x10})
	registry.Register(DeviceInfo{Address: 0x30})

	addresses := make([]byte, 0)
	registry.Iterate(func(entry DeviceEntry) bool {
		addresses = append(addresses, entry.Address())
		return len(addresses) < 2
	})

	if len(addresses) != 2 {
		t.Fatalf("expected early stop after 2 entries, got %d", len(addresses))
	}
	if addresses[0] != 0x08 || addresses[1] != 0x10 {
		t.Fatalf("unexpected iteration order: %v", addresses)
	}
}

func TestDeviceRegistry_ProviderMatching(t *testing.T) {
	planeHeating := mockPlane{name: "heating"}
	planeSystem := mockPlane{name: "system"}

	providerA := &countingProvider{
		name:    "vaillant",
		matchFn: func(info DeviceInfo) bool { return info.Manufacturer == "vaillant" },
		createFn: func(info DeviceInfo) []Plane {
			return []Plane{planeHeating}
		},
	}
	providerB := &countingProvider{
		name:    "noop",
		matchFn: func(info DeviceInfo) bool { return false },
		createFn: func(info DeviceInfo) []Plane {
			return []Plane{mockPlane{name: "noop"}}
		},
	}
	providerC := &countingProvider{
		name:    "system",
		matchFn: func(info DeviceInfo) bool { return info.Address == 0x10 },
		createFn: func(info DeviceInfo) []Plane {
			return []Plane{planeSystem}
		},
	}

	registry := NewDeviceRegistry([]PlaneProvider{providerA, providerB, providerC})

	registry.Register(DeviceInfo{Address: 0x08, Manufacturer: "vaillant"})
	registry.Register(DeviceInfo{Address: 0x10, Manufacturer: "other"})

	entry, ok := registry.Lookup(0x08)
	if !ok {
		t.Fatalf("expected device 0x08 to be registered")
	}
	if len(entry.Planes()) != 1 || entry.Planes()[0].Name() != "heating" {
		t.Fatalf("unexpected planes for 0x08: %v", entry.Planes())
	}

	entry, ok = registry.Lookup(0x10)
	if !ok {
		t.Fatalf("expected device 0x10 to be registered")
	}
	if len(entry.Planes()) != 1 || entry.Planes()[0].Name() != "system" {
		t.Fatalf("unexpected planes for 0x10: %v", entry.Planes())
	}

	if providerA.matchCalls != 2 || providerB.matchCalls != 2 || providerC.matchCalls != 2 {
		t.Fatalf("unexpected match call counts: A=%d B=%d C=%d", providerA.matchCalls, providerB.matchCalls, providerC.matchCalls)
	}
	if providerA.createCalls != 1 || providerB.createCalls != 0 || providerC.createCalls != 1 {
		t.Fatalf("unexpected create call counts: A=%d B=%d C=%d", providerA.createCalls, providerB.createCalls, providerC.createCalls)
	}
}

func TestDeviceRegistry_AliasAddressesShareCanonicalEntry(t *testing.T) {
	t.Parallel()

	registry := NewDeviceRegistry(nil)
	registry.Register(DeviceInfo{
		Address:         0x08,
		Manufacturer:    "Vaillant",
		DeviceID:        "BAI00",
		SoftwareVersion: "0704",
		HardwareVersion: "7603",
	})
	registry.Register(DeviceInfo{
		Address:         0x09,
		Manufacturer:    "vaillant",
		DeviceID:        "bai00",
		SoftwareVersion: "0704",
		HardwareVersion: "7603",
	})

	entry08, ok := registry.Lookup(0x08)
	if !ok {
		t.Fatalf("expected alias primary address lookup")
	}
	entry09, ok := registry.Lookup(0x09)
	if !ok {
		t.Fatalf("expected alias secondary address lookup")
	}
	if entry08 != entry09 {
		t.Fatalf("expected aliases to resolve to the same device entry")
	}
	if entry08.Address() != 0x08 {
		t.Fatalf("canonical address = %02x; want 08", entry08.Address())
	}
	if !slices.Equal(entry08.Addresses(), []byte{0x08, 0x09}) {
		t.Fatalf("addresses = %v; want [8 9]", entry08.Addresses())
	}

	addresses := make([]byte, 0)
	registry.Iterate(func(entry DeviceEntry) bool {
		addresses = append(addresses, entry.Address())
		return true
	})
	if !slices.Equal(addresses, []byte{0x08}) {
		t.Fatalf("iterate canonical addresses = %v; want [8]", addresses)
	}
}

func TestDeviceRegistry_DoesNotMergeDifferentSerials(t *testing.T) {
	t.Parallel()

	registry := NewDeviceRegistry(nil)
	registry.Register(DeviceInfo{
		Address:         0x08,
		Manufacturer:    "Vaillant",
		DeviceID:        "VR92",
		SoftwareVersion: "0514",
		HardwareVersion: "1204",
		SerialNumber:    "SN-1",
	})
	registry.Register(DeviceInfo{
		Address:         0x09,
		Manufacturer:    "Vaillant",
		DeviceID:        "VR92",
		SoftwareVersion: "0514",
		HardwareVersion: "1204",
		SerialNumber:    "SN-2",
	})

	entry08, ok := registry.Lookup(0x08)
	if !ok {
		t.Fatalf("expected address 0x08 to exist")
	}
	entry09, ok := registry.Lookup(0x09)
	if !ok {
		t.Fatalf("expected address 0x09 to exist")
	}
	if entry08 == entry09 {
		t.Fatalf("different serial numbers must not merge")
	}
	if !slices.Equal(entry08.Addresses(), []byte{0x08}) {
		t.Fatalf("entry08 addresses = %v; want [8]", entry08.Addresses())
	}
	if !slices.Equal(entry09.Addresses(), []byte{0x09}) {
		t.Fatalf("entry09 addresses = %v; want [9]", entry09.Addresses())
	}
}

func TestDeviceRegistry_SplitsAliasWhenConflictingSerialArrives(t *testing.T) {
	t.Parallel()

	registry := NewDeviceRegistry(nil)
	registry.Register(DeviceInfo{
		Address:         0x08,
		Manufacturer:    "Vaillant",
		DeviceID:        "VR92",
		SoftwareVersion: "0514",
		HardwareVersion: "1204",
	})
	registry.Register(DeviceInfo{
		Address:         0x09,
		Manufacturer:    "Vaillant",
		DeviceID:        "VR92",
		SoftwareVersion: "0514",
		HardwareVersion: "1204",
	})
	registry.Register(DeviceInfo{
		Address:         0x09,
		Manufacturer:    "Vaillant",
		DeviceID:        "VR92",
		SoftwareVersion: "0514",
		HardwareVersion: "1204",
		SerialNumber:    "SN-B",
	})
	registry.Register(DeviceInfo{
		Address:         0x08,
		Manufacturer:    "Vaillant",
		DeviceID:        "VR92",
		SoftwareVersion: "0514",
		HardwareVersion: "1204",
		SerialNumber:    "SN-A",
	})

	entry08, ok := registry.Lookup(0x08)
	if !ok {
		t.Fatalf("expected address 0x08 to exist")
	}
	entry09, ok := registry.Lookup(0x09)
	if !ok {
		t.Fatalf("expected address 0x09 to exist")
	}
	if entry08 == entry09 {
		t.Fatalf("conflicting serial numbers must split aliases")
	}
	if entry08.SerialNumber() != "SN-A" {
		t.Fatalf("entry08 serial = %q; want SN-A", entry08.SerialNumber())
	}
	if entry09.SerialNumber() != "SN-B" {
		t.Fatalf("entry09 serial = %q; want SN-B", entry09.SerialNumber())
	}
}

func TestDeviceRegistry_MergesSameSerialDifferentDeviceID(t *testing.T) {
	t.Parallel()

	registry := NewDeviceRegistry(nil)

	// Step 1-2: scan discovers two addresses with different DeviceIDs
	registry.Register(DeviceInfo{
		Address:         0x15,
		Manufacturer:    "Vaillant",
		DeviceID:        "BASV2",
		SoftwareVersion: "0507",
		HardwareVersion: "1704",
	})
	registry.Register(DeviceInfo{
		Address:         0xEC,
		Manufacturer:    "Vaillant",
		DeviceID:        "SOL00",
		SoftwareVersion: "0507",
		HardwareVersion: "1704",
	})

	// Before enrichment: separate entries
	entry15, _ := registry.Lookup(0x15)
	entryEC, _ := registry.Lookup(0xEC)
	if entry15 == entryEC {
		t.Fatalf("expected separate entries before serial enrichment")
	}

	// Step 3-4: enrichment discovers same serial on both
	registry.Register(DeviceInfo{
		Address:         0x15,
		Manufacturer:    "Vaillant",
		DeviceID:        "BASV2",
		SoftwareVersion: "0507",
		HardwareVersion: "1704",
		SerialNumber:    "21-22-09-0020184848-0082-005409-N4",
	})
	registry.Register(DeviceInfo{
		Address:         0xEC,
		Manufacturer:    "Vaillant",
		DeviceID:        "SOL00",
		SoftwareVersion: "0507",
		HardwareVersion: "1704",
		SerialNumber:    "21-22-09-0020184848-0082-005409-N4",
	})

	// After enrichment: merged into single entry
	entry15, ok := registry.Lookup(0x15)
	if !ok {
		t.Fatalf("expected address 0x15 to exist")
	}
	entryEC, ok = registry.Lookup(0xEC)
	if !ok {
		t.Fatalf("expected address 0xEC to exist")
	}
	if entry15 != entryEC {
		t.Fatalf("same serial number must merge entries regardless of DeviceID")
	}
	if entry15.DeviceID() != "BASV2" {
		t.Fatalf("merged DeviceID = %q; want BASV2 (primary)", entry15.DeviceID())
	}
	if !slices.Equal(entry15.Addresses(), []byte{0x15, 0xEC}) {
		t.Fatalf("merged addresses = %v; want [21 236]", entry15.Addresses())
	}

	// Iterate should return single canonical entry
	var count int
	registry.Iterate(func(entry DeviceEntry) bool {
		count++
		return true
	})
	if count != 1 {
		t.Fatalf("iterate count = %d; want 1 (merged)", count)
	}
}

func TestDeviceRegistry_PreservesKnownIdentityOnSparseUpdate(t *testing.T) {
	t.Parallel()

	registry := NewDeviceRegistry(nil)
	registry.Register(DeviceInfo{
		Address:         0x08,
		Manufacturer:    "Vaillant",
		DeviceID:        "BASV2",
		SoftwareVersion: "0507",
		HardwareVersion: "1704",
		SerialNumber:    "SN-KEEP",
		MacAddress:      "AA:BB:CC:DD:EE:FF",
	})
	registry.Register(DeviceInfo{
		Address:         0x08,
		Manufacturer:    "Vaillant",
		DeviceID:        "BASV2",
		SoftwareVersion: "0507",
		HardwareVersion: "1704",
	})

	entry, ok := registry.Lookup(0x08)
	if !ok {
		t.Fatalf("expected address 0x08 to exist")
	}
	if entry.SerialNumber() != "SN-KEEP" {
		t.Fatalf("serial = %q; want SN-KEEP", entry.SerialNumber())
	}
	if entry.MacAddress() != "AA:BB:CC:DD:EE:FF" {
		t.Fatalf("mac = %q; want preserved MAC", entry.MacAddress())
	}
}
