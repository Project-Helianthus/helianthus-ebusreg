package registry

import (
	stderrors "errors"
	"reflect"
	"testing"

	ebuserrors "github.com/Project-Helianthus/helianthus-ebusgo/errors"
	"github.com/Project-Helianthus/helianthus-ebusgo/types"
	"github.com/Project-Helianthus/helianthus-ebusreg/schema"
)

type projectionTestTemplate struct {
	primary   byte
	secondary byte
}

func (template projectionTestTemplate) Primary() byte   { return template.primary }
func (template projectionTestTemplate) Secondary() byte { return template.secondary }

type projectionTestMethod struct {
	name       string
	readOnly   bool
	template   FrameTemplate
	selector   schema.SchemaSelector
	mutability MethodMutability
	danger     MethodDanger
	routable   bool
}

func (method projectionTestMethod) Name() string                          { return method.name }
func (method projectionTestMethod) ReadOnly() bool                        { return method.readOnly }
func (method projectionTestMethod) Template() FrameTemplate               { return method.template }
func (method projectionTestMethod) ResponseSchema() schema.SchemaSelector { return method.selector }
func (method projectionTestMethod) Mutability() MethodMutability          { return method.mutability }
func (method projectionTestMethod) Danger() MethodDanger                  { return method.danger }
func (method projectionTestMethod) Routable() bool                        { return method.routable }

type projectionTestPlane struct {
	name    string
	methods []Method
}

func (plane projectionTestPlane) Name() string      { return plane.name }
func (plane projectionTestPlane) Methods() []Method { return plane.methods }

type projectionTestEntry struct {
	address         byte
	addresses       []byte
	manufacturer    string
	deviceID        string
	serialNumber    string
	macAddress      string
	softwareVersion string
	hardwareVersion string
	planes          []Plane
	projections     []Projection
}

func (entry projectionTestEntry) Address() byte             { return entry.address }
func (entry projectionTestEntry) Addresses() []byte         { return append([]byte(nil), entry.addresses...) }
func (entry projectionTestEntry) Manufacturer() string      { return entry.manufacturer }
func (entry projectionTestEntry) DeviceID() string          { return entry.deviceID }
func (entry projectionTestEntry) SerialNumber() string      { return entry.serialNumber }
func (entry projectionTestEntry) MacAddress() string        { return entry.macAddress }
func (entry projectionTestEntry) SoftwareVersion() string   { return entry.softwareVersion }
func (entry projectionTestEntry) HardwareVersion() string   { return entry.hardwareVersion }
func (entry projectionTestEntry) Planes() []Plane           { return entry.planes }
func (entry projectionTestEntry) Projections() []Projection { return entry.projections }

type projectionTestIterator struct {
	entries []DeviceEntry
}

func (iterator projectionTestIterator) Iterate(fn func(DeviceEntry) bool) {
	for _, entry := range iterator.entries {
		if !fn(entry) {
			return
		}
	}
}

func TestProjectMethod(t *testing.T) {
	t.Parallel()

	selector := schema.SchemaSelector{
		Default: schema.Schema{Fields: []schema.SchemaField{{Name: "temp_c", Type: types.DATA2b{}}}},
	}
	method := projectionTestMethod{
		name:     "set_target",
		readOnly: false,
		template: projectionTestTemplate{primary: 0xB5, secondary: 0x05},
		selector: selector,
		// Explicit metadata provider values take priority over ReadOnly fallback.
		mutability: MethodMutabilityUnknown,
		danger:     MethodDangerDangerous,
		routable:   false,
	}

	projected, err := ProjectMethod(method)
	if err != nil {
		t.Fatalf("ProjectMethod error = %v", err)
	}
	if projected.Name != "set_target" {
		t.Fatalf("Name = %q; want set_target", projected.Name)
	}
	if projected.ReadOnly {
		t.Fatalf("ReadOnly = true; want false")
	}
	if projected.Primary != 0xB5 || projected.Secondary != 0x05 {
		t.Fatalf("template = (%02x,%02x); want (b5,05)", projected.Primary, projected.Secondary)
	}
	if !reflect.DeepEqual(projected.ResponseSelector, selector) {
		t.Fatalf("ResponseSelector = %#v; want %#v", projected.ResponseSelector, selector)
	}
	if projected.Metadata.Mutability != MethodMutabilityUnknown {
		t.Fatalf("Mutability = %q; want unknown", projected.Metadata.Mutability)
	}
	if projected.Metadata.Danger != MethodDangerDangerous {
		t.Fatalf("Danger = %q; want dangerous", projected.Metadata.Danger)
	}
	if projected.Metadata.Routable {
		t.Fatalf("Routable = true; want false")
	}
}

func TestProjectMethodErrors(t *testing.T) {
	t.Parallel()

	if _, err := ProjectMethod(nil); !stderrors.Is(err, ebuserrors.ErrInvalidPayload) {
		t.Fatalf("ProjectMethod(nil) error = %v; want ErrInvalidPayload", err)
	}

	invalidMethod := projectionTestMethod{name: "broken", template: nil}
	if _, err := ProjectMethod(invalidMethod); !stderrors.Is(err, ebuserrors.ErrInvalidPayload) {
		t.Fatalf("ProjectMethod(missing template) error = %v; want ErrInvalidPayload", err)
	}
}

func TestProjectDeviceEntryAndPlane(t *testing.T) {
	t.Parallel()

	readMethod := projectionTestMethod{
		name:       "get_status",
		readOnly:   true,
		template:   projectionTestTemplate{primary: 0xB5, secondary: 0x04},
		selector:   schema.SchemaSelector{},
		mutability: MethodMutabilityReadOnly,
		danger:     MethodDangerSafe,
		routable:   true,
	}
	writeMethod := projectionTestMethod{
		name:       "set_target",
		readOnly:   false,
		template:   projectionTestTemplate{primary: 0xB5, secondary: 0x05},
		selector:   schema.SchemaSelector{},
		mutability: MethodMutabilityMutating,
		danger:     MethodDangerDangerous,
		routable:   true,
	}
	plane := projectionTestPlane{name: "Heating", methods: []Method{readMethod, writeMethod}}
	entry := projectionTestEntry{
		address:         0x08,
		addresses:       []byte{0x10, 0x08, 0x10},
		manufacturer:    "Vaillant",
		deviceID:        "BAI00",
		serialNumber:    "ABC",
		macAddress:      "DEF",
		softwareVersion: "123",
		hardwareVersion: "456",
		planes:          []Plane{plane},
	}

	projected, err := ProjectDeviceEntry(entry)
	if err != nil {
		t.Fatalf("ProjectDeviceEntry error = %v", err)
	}
	if projected.Address != 0x08 {
		t.Fatalf("Address = %02x; want 08", projected.Address)
	}
	if !reflect.DeepEqual(projected.Addresses, []byte{0x08, 0x10}) {
		t.Fatalf("Addresses = %v; want [8 16]", projected.Addresses)
	}
	if projected.Manufacturer != "Vaillant" || projected.DeviceID != "BAI00" {
		t.Fatalf("identity projection mismatch: %#v", projected)
	}
	if len(projected.Planes) != 1 {
		t.Fatalf("Planes len = %d; want 1", len(projected.Planes))
	}
	if projected.Planes[0].Name != "Heating" {
		t.Fatalf("Plane name = %q; want Heating", projected.Planes[0].Name)
	}
	if len(projected.Planes[0].Methods) != 2 {
		t.Fatalf("Methods len = %d; want 2", len(projected.Planes[0].Methods))
	}

	projected.Addresses[0] = 0xFF
	again, err := ProjectDeviceEntry(entry)
	if err != nil {
		t.Fatalf("ProjectDeviceEntry(second) error = %v", err)
	}
	if !reflect.DeepEqual(again.Addresses, []byte{0x08, 0x10}) {
		t.Fatalf("Addresses leaked mutation: %v", again.Addresses)
	}
}

func TestProjectRegistryDevices(t *testing.T) {
	t.Parallel()

	entryA := projectionTestEntry{
		address:   0x08,
		addresses: []byte{0x08},
		planes:    []Plane{projectionTestPlane{name: "System", methods: []Method{projectionTestMethod{name: "a", template: projectionTestTemplate{primary: 0xB5, secondary: 0x04}, selector: schema.SchemaSelector{}, mutability: MethodMutabilityReadOnly, danger: MethodDangerSafe, routable: true, readOnly: true}}}},
	}
	entryB := projectionTestEntry{
		address:   0x10,
		addresses: []byte{0x10},
		planes:    []Plane{projectionTestPlane{name: "Heating", methods: []Method{projectionTestMethod{name: "b", template: projectionTestTemplate{primary: 0xB5, secondary: 0x05}, selector: schema.SchemaSelector{}, mutability: MethodMutabilityMutating, danger: MethodDangerDangerous, routable: true, readOnly: false}}}},
	}
	iterator := projectionTestIterator{entries: []DeviceEntry{entryA, entryB}}

	projected, err := ProjectRegistryDevices(iterator)
	if err != nil {
		t.Fatalf("ProjectRegistryDevices error = %v", err)
	}
	if len(projected) != 2 {
		t.Fatalf("len = %d; want 2", len(projected))
	}
	if projected[0].Address != 0x08 || projected[1].Address != 0x10 {
		t.Fatalf("order mismatch: %02x %02x", projected[0].Address, projected[1].Address)
	}

	if _, err := ProjectRegistryDevices(nil); !stderrors.Is(err, ebuserrors.ErrInvalidPayload) {
		t.Fatalf("ProjectRegistryDevices(nil) error = %v; want ErrInvalidPayload", err)
	}

	badIterator := projectionTestIterator{entries: []DeviceEntry{entryA, nil, entryB}}
	if _, err := ProjectRegistryDevices(badIterator); !stderrors.Is(err, ebuserrors.ErrInvalidPayload) {
		t.Fatalf("ProjectRegistryDevices(bad iterator) error = %v; want ErrInvalidPayload", err)
	}
}

func TestProjectRegistryDevicesDeterministicOrdering(t *testing.T) {
	t.Parallel()

	methodZ := projectionTestMethod{
		name:       "zeta",
		readOnly:   true,
		template:   projectionTestTemplate{primary: 0xB5, secondary: 0x09},
		selector:   schema.SchemaSelector{},
		mutability: MethodMutabilityReadOnly,
		danger:     MethodDangerSafe,
		routable:   true,
	}
	methodA := projectionTestMethod{
		name:       "alpha",
		readOnly:   true,
		template:   projectionTestTemplate{primary: 0xB5, secondary: 0x01},
		selector:   schema.SchemaSelector{},
		mutability: MethodMutabilityReadOnly,
		danger:     MethodDangerSafe,
		routable:   true,
	}

	entryB := projectionTestEntry{
		address:         0x10,
		addresses:       []byte{0x15, 0x10, 0x12},
		manufacturer:    "vaillant",
		deviceID:        "B",
		hardwareVersion: "2",
		serialNumber:    "2",
		planes: []Plane{
			projectionTestPlane{name: "zPlane", methods: []Method{methodZ, methodA}},
			projectionTestPlane{name: "APlane", methods: []Method{methodZ}},
		},
	}
	entryA := projectionTestEntry{
		address:         0x08,
		addresses:       []byte{0x09, 0x08, 0x0A},
		manufacturer:    "Vaillant",
		deviceID:        "A",
		hardwareVersion: "1",
		serialNumber:    "1",
		planes: []Plane{
			projectionTestPlane{name: "System", methods: []Method{methodZ, methodA}},
		},
	}

	projected, err := ProjectRegistryDevices(projectionTestIterator{entries: []DeviceEntry{entryB, entryA}})
	if err != nil {
		t.Fatalf("ProjectRegistryDevices error = %v", err)
	}

	if len(projected) != 2 {
		t.Fatalf("devices len = %d; want 2", len(projected))
	}
	if projected[0].Address != 0x08 || projected[1].Address != 0x10 {
		t.Fatalf("device order = [%02x,%02x]; want [08,10]", projected[0].Address, projected[1].Address)
	}
	if !reflect.DeepEqual(projected[1].Addresses, []byte{0x10, 0x12, 0x15}) {
		t.Fatalf("addresses order = %v; want [16 18 21]", projected[1].Addresses)
	}
	if len(projected[1].Planes) != 2 {
		t.Fatalf("plane len = %d; want 2", len(projected[1].Planes))
	}
	if projected[1].Planes[0].Name != "APlane" || projected[1].Planes[1].Name != "zPlane" {
		t.Fatalf("plane order = [%s,%s]; want [APlane,zPlane]", projected[1].Planes[0].Name, projected[1].Planes[1].Name)
	}
	if len(projected[1].Planes[1].Methods) != 2 {
		t.Fatalf("method len = %d; want 2", len(projected[1].Planes[1].Methods))
	}
	if projected[1].Planes[1].Methods[0].Name != "alpha" || projected[1].Planes[1].Methods[1].Name != "zeta" {
		t.Fatalf(
			"method order = [%s,%s]; want [alpha,zeta]",
			projected[1].Planes[1].Methods[0].Name,
			projected[1].Planes[1].Methods[1].Name,
		)
	}
}

func TestProjectPlaneErrors(t *testing.T) {
	t.Parallel()

	if _, err := ProjectPlane(nil); !stderrors.Is(err, ebuserrors.ErrInvalidPayload) {
		t.Fatalf("ProjectPlane(nil) error = %v; want ErrInvalidPayload", err)
	}

	badPlane := projectionTestPlane{name: "Bad", methods: []Method{projectionTestMethod{name: "broken", template: nil}}}
	if _, err := ProjectPlane(badPlane); !stderrors.Is(err, ebuserrors.ErrInvalidPayload) {
		t.Fatalf("ProjectPlane(bad method) error = %v; want ErrInvalidPayload", err)
	}

	if _, err := ProjectDeviceEntry(nil); !stderrors.Is(err, ebuserrors.ErrInvalidPayload) {
		t.Fatalf("ProjectDeviceEntry(nil) error = %v; want ErrInvalidPayload", err)
	}
}

func TestProjectPlaneDeterministicOrdering(t *testing.T) {
	t.Parallel()

	methods := []Method{
		projectionTestMethod{
			name:       "zeta",
			readOnly:   true,
			template:   projectionTestTemplate{primary: 0xB5, secondary: 0x02},
			selector:   schema.SchemaSelector{},
			mutability: MethodMutabilityReadOnly,
			danger:     MethodDangerSafe,
			routable:   true,
		},
		projectionTestMethod{
			name:       "alpha",
			readOnly:   true,
			template:   projectionTestTemplate{primary: 0xB5, secondary: 0x01},
			selector:   schema.SchemaSelector{},
			mutability: MethodMutabilityReadOnly,
			danger:     MethodDangerSafe,
			routable:   true,
		},
	}

	projected, err := ProjectPlane(projectionTestPlane{name: "Heating", methods: methods})
	if err != nil {
		t.Fatalf("ProjectPlane error = %v", err)
	}
	if len(projected.Methods) != 2 {
		t.Fatalf("methods len = %d; want 2", len(projected.Methods))
	}
	if projected.Methods[0].Name != "alpha" || projected.Methods[1].Name != "zeta" {
		t.Fatalf("method order = [%s,%s]; want [alpha,zeta]", projected.Methods[0].Name, projected.Methods[1].Name)
	}
}
