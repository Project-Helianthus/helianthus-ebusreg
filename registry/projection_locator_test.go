package registry

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

// locatorFixtureDevice is a synthetic device used by the locator tests. The
// first four fields form the qualified identity the registry admits on; the
// address is where the device currently answers on the bus.
type locatorFixtureDevice struct {
	manufacturer string
	deviceID     string
	serial       string
	mac          string
	address      byte
}

func (device locatorFixtureDevice) qualifiedIdentity() string {
	return strings.Join([]string{device.manufacturer, device.deviceID, device.serial, device.mac}, "|")
}

// nativeLocatorPaths builds the per-device canonical path shape that native
// providers anchor their projections at (bus, address, device segment, then a
// method), and the matching path on a second plane.
func nativeLocatorPaths(device locatorFixtureDevice) (ProjectionPath, ProjectionPath) {
	segments := []PathSegment{
		{Name: "ebus"},
		{Name: fmt.Sprintf("addr@%02X", device.address)},
		{Name: "device@" + device.deviceID},
		{Name: "method@get_operational_data"},
	}
	return ProjectionPath{Plane: ServicePlane, Segments: segments},
		ProjectionPath{Plane: "Observability", Segments: segments}
}

func locatorFor(t *testing.T, device locatorFixtureDevice) NodeID {
	t.Helper()
	canonical, observed := nativeLocatorPaths(device)
	service, err := NewNode(canonical, canonical)
	if err != nil {
		t.Fatalf("service node: %v", err)
	}
	other, err := NewNode(observed, canonical)
	if err != nil {
		t.Fatalf("observability node: %v", err)
	}
	if service.ID != other.ID {
		t.Fatalf("locator differs across planes inside one snapshot: %q vs %q", service.ID, other.ID)
	}
	if _, err := BuildCanonicalIndex([]Projection{
		{Plane: ServicePlane, Nodes: []Node{service}},
		{Plane: "Observability", Nodes: []Node{other}},
	}); err != nil {
		t.Fatalf("canonical index: %v", err)
	}
	stable, err := StableNodeID(canonical)
	if err != nil {
		t.Fatalf("stable node id: %v", err)
	}
	if stable != service.ID {
		t.Fatalf("node id %q differs from StableNodeID %q", service.ID, stable)
	}
	return service.ID
}

// The locator only finds a node inside one registry snapshot. Moving the same
// qualified device to another bus address moves its locator, so a consumer
// that stored the locator as resource identity would split one device into
// two resources.
func TestProjectionLocatorChangesWhenTheSameDeviceIsReaddressed(t *testing.T) {
	before := locatorFixtureDevice{manufacturer: "ExampleMaker", deviceID: "UNIT0", serial: "SYNTHETIC-SERIAL-A", mac: "SYNTHETIC-MAC-A", address: 0x15}
	after := before
	after.address = 0x26

	if before.qualifiedIdentity() != after.qualifiedIdentity() {
		t.Fatalf("fixture must keep the qualified identity: %q vs %q", before.qualifiedIdentity(), after.qualifiedIdentity())
	}
	beforeLocator := locatorFor(t, before)
	afterLocator := locatorFor(t, after)
	if beforeLocator == afterLocator {
		t.Fatalf("re-addressing kept locator %q; the locator must follow the native path, so it cannot be resource identity", beforeLocator)
	}
}

// A same-model replacement at the same address reuses the locator of the old
// unit although the qualified identity is different, so a consumer that stored
// the locator as resource identity would merge two devices into one resource.
func TestProjectionLocatorRepeatsWhenTheSameModelIsSwappedInAtOneAddress(t *testing.T) {
	original := locatorFixtureDevice{manufacturer: "ExampleMaker", deviceID: "UNIT0", serial: "SYNTHETIC-SERIAL-A", mac: "SYNTHETIC-MAC-A", address: 0x15}
	replacements := map[string]locatorFixtureDevice{
		"serial":       {manufacturer: original.manufacturer, deviceID: original.deviceID, serial: "SYNTHETIC-SERIAL-B", mac: original.mac, address: original.address},
		"mac":          {manufacturer: original.manufacturer, deviceID: original.deviceID, serial: original.serial, mac: "SYNTHETIC-MAC-B", address: original.address},
		"manufacturer": {manufacturer: "OtherMaker", deviceID: original.deviceID, serial: original.serial, mac: original.mac, address: original.address},
	}

	originalLocator := locatorFor(t, original)
	for field, replacement := range replacements {
		t.Run(field, func(t *testing.T) {
			if replacement.qualifiedIdentity() == original.qualifiedIdentity() {
				t.Fatalf("fixture must change the qualified identity through %s", field)
			}
			if got := locatorFor(t, replacement); got != originalLocator {
				t.Fatalf("swap changing %s produced locator %q, want %q: the locator carries no identity, so it repeats for a different unit", field, got, originalLocator)
			}
		})
	}
}

// The exported locator API must say what the value is, so that the gateway,
// the Portal and the semantic registry do not persist it as resource identity.
func TestProjectionLocatorDocumentationNamesASnapshotLocator(t *testing.T) {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, "projection.go", nil, parser.ParseComments)
	if err != nil {
		t.Fatalf("parse projection.go: %v", err)
	}
	docs := map[string]string{}
	for _, declaration := range file.Decls {
		switch typed := declaration.(type) {
		case *ast.FuncDecl:
			if typed.Recv == nil && typed.Doc != nil {
				docs[typed.Name.Name] = typed.Doc.Text()
			}
		case *ast.GenDecl:
			for _, spec := range typed.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				doc := typeSpec.Doc
				if doc == nil {
					doc = typed.Doc
				}
				if doc != nil {
					docs[typeSpec.Name.Name] = doc.Text()
				}
			}
		}
	}
	for _, name := range []string{"NodeID", "StableNodeID"} {
		text := strings.Join(strings.Fields(docs[name]), " ")
		for _, phrase := range []string{"snapshot locator", "native-only", "resource identity"} {
			if !strings.Contains(text, phrase) {
				t.Errorf("doc comment of %s must mention %q; got %q", name, phrase, text)
			}
		}
	}
}
