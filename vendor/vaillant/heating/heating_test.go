package heating

import (
	"testing"

	"github.com/d3vi1/helianthus-ebusreg/registry"
)

func TestProviderMatch(t *testing.T) {
	t.Parallel()

	provider := NewProvider()
	if !provider.Match(registry.DeviceInfo{Manufacturer: "Vaillant"}) {
		t.Fatalf("expected Vaillant to match")
	}
	if !provider.Match(registry.DeviceInfo{Manufacturer: "Saunier Duval"}) {
		t.Fatalf("expected Saunier Duval to match")
	}
	if provider.Match(registry.DeviceInfo{Manufacturer: "Other"}) {
		t.Fatalf("did not expect Other to match")
	}
}

func TestProviderMethods(t *testing.T) {
	t.Parallel()

	provider := NewProvider()
	planes := provider.CreatePlanes(registry.DeviceInfo{Manufacturer: "Vaillant"})
	if len(planes) != 1 {
		t.Fatalf("expected 1 plane, got %d", len(planes))
	}

	plane := planes[0]
	if plane.Name() != "heating" {
		t.Fatalf("expected plane name heating, got %s", plane.Name())
	}

	methods := plane.Methods()
	if len(methods) != 3 {
		t.Fatalf("expected 3 methods, got %d", len(methods))
	}

	if methods[0].Name() != methodGetStatus || !methods[0].ReadOnly() {
		t.Fatalf("unexpected get_status method definition")
	}
	if methods[1].Name() != methodSetTargetTemp || methods[1].ReadOnly() {
		t.Fatalf("unexpected set_target_temp method definition")
	}
	if methods[2].Name() != methodGetParameters || !methods[2].ReadOnly() {
		t.Fatalf("unexpected get_parameters method definition")
	}

	if methods[0].Template().Primary() != 0xB5 || methods[0].Template().Secondary() != 0x04 {
		t.Fatalf("unexpected get_status template")
	}
	if methods[1].Template().Primary() != 0xB5 || methods[1].Template().Secondary() != 0x05 {
		t.Fatalf("unexpected set_target_temp template")
	}
	if methods[2].Template().Primary() != 0xB5 || methods[2].Template().Secondary() != 0x04 {
		t.Fatalf("unexpected get_parameters template")
	}
}
