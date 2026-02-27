package system

import (
	"testing"

	"github.com/d3vi1/helianthus-ebusreg/registry"
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
