package system

import (
	"testing"

	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
)

func TestCreateProjections_BASV2(t *testing.T) {
	t.Parallel()

	info := registry.DeviceInfo{
		Manufacturer: "Vaillant",
		DeviceID:     "BASV2",
		Address:      0x10,
	}
	planes := NewProvider().CreatePlanes(info)
	projections := NewProvider().CreateProjections(info, planes)

	if !hasNodePath(projections, registry.ServicePlane, "Service:/ebus/addr@10/device@BASV2/method@get_operational_data") {
		t.Fatalf("missing Service operational node")
	}
	if !hasNodePath(projections, registry.ServicePlane, "Service:/ebus/addr@10/device@BASV2/method@get_register") {
		t.Fatalf("missing Service register node")
	}
	if !hasNodePath(projections, registry.ServicePlane, "Service:/ebus/addr@10/device@BASV2/method@get_ext_register") {
		t.Fatalf("missing Service ext_register node")
	}
	if !hasNodePath(projections, "Observability", "Observability:/ebus/addr@10/device@BASV2/method@get_operational_data") {
		t.Fatalf("missing Observability operational node")
	}
	if !hasNodePath(projections, "Debug", "Debug:/ebus/addr@10/device@BASV2/register@b509") {
		t.Fatalf("missing Debug b509 node")
	}
	if !hasNodePath(projections, "Debug", "Debug:/ebus/addr@10/device@BASV2/register@b524") {
		t.Fatalf("missing Debug b524 node")
	}
}

func TestCreateProjections_BAI00(t *testing.T) {
	t.Parallel()

	info := registry.DeviceInfo{
		Manufacturer: "Vaillant",
		DeviceID:     "BAI00",
		Address:      0x08,
	}
	planes := NewProvider().CreatePlanes(info)
	projections := NewProvider().CreateProjections(info, planes)

	if !hasNodePath(projections, registry.ServicePlane, "Service:/ebus/addr@08/device@BAI00/method@get_operational_data") {
		t.Fatalf("missing Service operational node")
	}
	if !hasNodePath(projections, registry.ServicePlane, "Service:/ebus/addr@08/device@BAI00/method@get_register") {
		t.Fatalf("missing Service register node")
	}
	if !hasNodePath(projections, registry.ServicePlane, "Service:/ebus/addr@08/device@BAI00/method@get_ext_register") {
		t.Fatalf("missing Service ext_register node")
	}
	if !hasNodePath(projections, "Observability", "Observability:/ebus/addr@08/device@BAI00/method@get_operational_data") {
		t.Fatalf("missing Observability operational node")
	}
	if !hasNodePath(projections, "Debug", "Debug:/ebus/addr@08/device@BAI00/register@b509") {
		t.Fatalf("missing Debug b509 node")
	}
	if !hasNodePath(projections, "Debug", "Debug:/ebus/addr@08/device@BAI00/register@b524") {
		t.Fatalf("missing Debug b524 node")
	}
}

func TestCreateProjections_DeviceFilter(t *testing.T) {
	t.Parallel()

	info := registry.DeviceInfo{
		Manufacturer: "Vaillant",
		DeviceID:     "VRC720",
		Address:      0x10,
	}
	planes := NewProvider().CreatePlanes(info)
	projections := NewProvider().CreateProjections(info, planes)
	if len(projections) == 0 {
		t.Fatalf("expected projections for VRC720, got 0")
	}
	if !hasNodePath(projections, "Service", "Service:/ebus/addr@10/device@VRC720/method@get_operational_data") {
		t.Fatalf("missing Service operational node for VRC720")
	}
	if !hasNodePath(projections, "Debug", "Debug:/ebus/addr@10/device@VRC720/register@b509") {
		t.Fatalf("missing Debug b509 node for VRC720")
	}

	info.DeviceID = "VR_71"
	projections = NewProvider().CreateProjections(info, planes)
	if len(projections) == 0 {
		t.Fatalf("expected projections for VR_71")
	}
}

func TestCreateProjections_NoPlanes(t *testing.T) {
	t.Parallel()

	info := registry.DeviceInfo{
		Manufacturer: "Vaillant",
		DeviceID:     "BASV2",
		Address:      0x10,
	}
	projections := NewProvider().CreateProjections(info, nil)
	if projections != nil {
		t.Fatalf("expected nil projections for empty planes, got %d", len(projections))
	}
	projections = NewProvider().CreateProjections(info, []registry.Plane{})
	if projections != nil {
		t.Fatalf("expected nil projections for zero-length planes, got %d", len(projections))
	}
	projections = NewProvider().CreateProjections(info, []registry.Plane{nil, nil})
	if projections != nil {
		t.Fatalf("expected nil projections for nil-element planes, got %d", len(projections))
	}
}

func TestCreateProjections_BuildCanonicalIndex(t *testing.T) {
	t.Parallel()

	devices := []registry.DeviceInfo{
		{Manufacturer: "Vaillant", DeviceID: "BAI00", Address: 0x08},
		{Manufacturer: "Vaillant", DeviceID: "BASV2", Address: 0x15},
		{Manufacturer: "Vaillant", DeviceID: "VRC720", Address: 0x10},
	}
	provider := NewProvider()
	for _, info := range devices {
		planes := provider.CreatePlanes(info)
		projections := provider.CreateProjections(info, planes)
		if len(projections) == 0 {
			t.Fatalf("expected projections for %s", info.DeviceID)
		}
		_, err := registry.BuildCanonicalIndex(projections)
		if err != nil {
			t.Fatalf("BuildCanonicalIndex failed for %s: %v", info.DeviceID, err)
		}
	}
}

func TestCreateProjectionsLocatorIsNotResourceIdentity(t *testing.T) {
	t.Parallel()

	original := registry.DeviceInfo{
		Manufacturer:    "ExampleMaker",
		DeviceID:        "UNIT0",
		SerialNumber:    "SYNTHETIC-SERIAL-A",
		MacAddress:      "SYNTHETIC-MAC-A",
		SoftwareVersion: "0100",
		HardwareVersion: "0200",
		Address:         0x15,
	}
	originalLocators := serviceLocators(t, original)

	t.Run("re-addressing moves every locator of the same device", func(t *testing.T) {
		moved := original
		moved.Address = 0x26
		movedLocators := serviceLocators(t, moved)
		if len(movedLocators) != len(originalLocators) {
			t.Fatalf("re-addressed device has %d service nodes, want %d", len(movedLocators), len(originalLocators))
		}
		for method, locator := range originalLocators {
			if movedLocators[method] == locator {
				t.Fatalf("re-addressing kept locator %q for %s; the locator follows the address, so it is not resource identity", locator, method)
			}
		}
	})

	swaps := map[string]func(*registry.DeviceInfo){
		"serial":           func(info *registry.DeviceInfo) { info.SerialNumber = "SYNTHETIC-SERIAL-B" },
		"mac":              func(info *registry.DeviceInfo) { info.MacAddress = "SYNTHETIC-MAC-B" },
		"manufacturer":     func(info *registry.DeviceInfo) { info.Manufacturer = "OtherMaker" },
		"software version": func(info *registry.DeviceInfo) { info.SoftwareVersion = "0101" },
		"hardware version": func(info *registry.DeviceInfo) { info.HardwareVersion = "0201" },
	}
	for field, mutate := range swaps {
		t.Run("same-model swap changing "+field+" repeats every locator", func(t *testing.T) {
			replacement := original
			mutate(&replacement)
			if replacement == original {
				t.Fatalf("fixture must change %s", field)
			}
			replacementLocators := serviceLocators(t, replacement)
			if len(replacementLocators) != len(originalLocators) {
				t.Fatalf("replacement has %d service nodes, want %d", len(replacementLocators), len(originalLocators))
			}
			for method, locator := range originalLocators {
				if replacementLocators[method] != locator {
					t.Fatalf("swap changing %s gave locator %q for %s, want %q; the locator carries no identity", field, replacementLocators[method], method, locator)
				}
			}
		})
	}
}

func serviceLocators(t *testing.T, info registry.DeviceInfo) map[string]registry.NodeID {
	t.Helper()
	provider := NewProvider()
	projections := provider.CreateProjections(info, provider.CreatePlanes(info))
	if _, err := registry.BuildCanonicalIndex(projections); err != nil {
		t.Fatalf("canonical index: %v", err)
	}
	locators := map[string]registry.NodeID{}
	for _, projection := range projections {
		if projection.Plane != registry.ServicePlane {
			continue
		}
		for _, node := range projection.Nodes {
			segments := node.CanonicalPath.Segments
			locators[segments[len(segments)-1].Name] = node.ID
		}
	}
	if len(locators) == 0 {
		t.Fatalf("no service nodes for %+v", info)
	}
	return locators
}

func hasNodePath(projections []registry.Projection, plane string, path string) bool {
	for _, projection := range projections {
		if projection.Plane != plane {
			continue
		}
		for _, node := range projection.Nodes {
			if node.Path.String() == path {
				return true
			}
		}
	}
	return false
}
