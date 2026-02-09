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
	if len(projections) != 0 {
		t.Fatalf("expected no projections, got %d", len(projections))
	}

	info.DeviceID = "VR_71"
	projections = NewProvider().CreateProjections(info, planes)
	if len(projections) == 0 {
		t.Fatalf("expected projections for VR_71")
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
