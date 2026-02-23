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
	getMetadata := registry.ResolveMethodMetadata(getMethod)
	if getMetadata.Mutability != registry.MethodMutabilityReadOnly {
		t.Fatalf("get_operational_data mutability = %q; want %q", getMetadata.Mutability, registry.MethodMutabilityReadOnly)
	}
	if getMetadata.Danger != registry.MethodDangerSafe {
		t.Fatalf("get_operational_data danger = %q; want %q", getMetadata.Danger, registry.MethodDangerSafe)
	}
	if !getMetadata.Routable {
		t.Fatalf("get_operational_data routable = false; want true")
	}

	setMethod, ok := findMethod(methods, methodSetOperationalData)
	if !ok || setMethod.ReadOnly() {
		t.Fatalf("expected set_operational_data method")
	}
	if setMethod.Template().Primary() != 0xB5 || setMethod.Template().Secondary() != 0x05 {
		t.Fatalf("unexpected set_operational_data template")
	}
	setMetadata := registry.ResolveMethodMetadata(setMethod)
	if setMetadata.Mutability != registry.MethodMutabilityMutating {
		t.Fatalf("set_operational_data mutability = %q; want %q", setMetadata.Mutability, registry.MethodMutabilityMutating)
	}
	if setMetadata.Danger != registry.MethodDangerDangerous {
		t.Fatalf("set_operational_data danger = %q; want %q", setMetadata.Danger, registry.MethodDangerDangerous)
	}
	if !setMetadata.Routable {
		t.Fatalf("set_operational_data routable = false; want true")
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
