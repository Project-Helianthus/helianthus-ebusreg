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
