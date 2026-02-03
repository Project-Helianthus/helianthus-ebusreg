package system

import (
	"testing"

	"github.com/d3vi1/helianthus-ebusreg/registry"
	"github.com/d3vi1/helianthus-ebusreg/router"
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

func TestSubscriptions(t *testing.T) {
	t.Parallel()

	planes := NewProvider().CreatePlanes(registry.DeviceInfo{Manufacturer: "Vaillant"})
	if len(planes) != 1 {
		t.Fatalf("expected 1 plane, got %d", len(planes))
	}

	systemPlane, ok := planes[0].(*plane)
	if !ok {
		t.Fatalf("expected system plane type")
	}

	got := systemPlane.Subscriptions()
	want := []router.Subscription{
		{Primary: 0xB5, Secondary: 0x16},
		{Primary: 0xB5, Secondary: 0x16},
		{Primary: 0xFE, Secondary: 0x01},
	}
	if len(got) != len(want) {
		t.Fatalf("expected %d subscriptions, got %d", len(want), len(got))
	}
	for index := range want {
		if got[index] != want[index] {
			t.Fatalf("subscription %d = %+v; want %+v", index, got[index], want[index])
		}
	}
}

func TestMethods(t *testing.T) {
	t.Parallel()

	planes := NewProvider().CreatePlanes(registry.DeviceInfo{Manufacturer: "Vaillant", Address: 0x08})
	if len(planes) != 1 {
		t.Fatalf("expected 1 plane, got %d", len(planes))
	}

	systemPlane, ok := planes[0].(*plane)
	if !ok {
		t.Fatalf("expected system plane type")
	}

	methods := systemPlane.Methods()
	getMethod, ok := findMethod(methods, methodGetOperationalData)
	if !ok || !getMethod.ReadOnly() {
		t.Fatalf("expected get_operational_data method")
	}
	if getMethod.Template().Primary() != 0xB5 || getMethod.Template().Secondary() != 0x04 {
		t.Fatalf("unexpected get_operational_data template")
	}

	setMethod, ok := findMethod(methods, methodSetOperationalData)
	if !ok || setMethod.ReadOnly() {
		t.Fatalf("expected set_operational_data method")
	}
	if setMethod.Template().Primary() != 0xB5 || setMethod.Template().Secondary() != 0x05 {
		t.Fatalf("unexpected set_operational_data template")
	}
}

func findMethod(methods []registry.Method, name string) (registry.Method, bool) {
	for _, method := range methods {
		if method.Name() == name {
			return method, true
		}
	}
	return nil, false
}
