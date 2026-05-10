package registry

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/Project-Helianthus/helianthus-ebusreg/schema"
)

// P9.x — DeviceEntrySnapshot extended with Planes + Projections.
//
// The graphql.BuildSchema hot path used to iterate via reg.Iterate
// (live *deviceEntry pointer) and dereference entry.Planes() /
// entry.Projections() after the Iterate RLock had been released.
// During an identity-merge in mergeEntries, dst.planes /
// dst.projections are reassigned — readers iterating the live entry
// could observe a torn slice header. Snapshot now copies the slice
// header under RLock; downstream readers see a stable view.

// TestDeviceEntrySnapshot_PlanesPopulated verifies that Planes
// captured from registerLocked are present on the snapshot.
func TestDeviceEntrySnapshot_PlanesPopulated(t *testing.T) {
	t.Parallel()

	reg := NewDeviceRegistry([]PlaneProvider{stubPlaneProvider{name: "stub-plane"}})
	reg.Register(DeviceInfo{Address: 0x10, Manufacturer: "Vaillant", DeviceID: "X"})

	snap, ok := reg.LookupEntrySnapshot(0x10)
	if !ok {
		t.Fatalf("LookupEntrySnapshot(0x10) ok=false")
	}
	if len(snap.Planes) != 1 {
		t.Fatalf("snap.Planes len = %d; want 1", len(snap.Planes))
	}
	if got := snap.Planes[0].Name(); got != "stub-plane" {
		t.Errorf("snap.Planes[0].Name = %q; want stub-plane", got)
	}
}

// TestDeviceEntrySnapshot_PlanesIsolatedFromRegistryRebind verifies the
// snapshot's Planes slice is unaffected by a subsequent identity merge
// that rebinds entry.planes. This is the regression test for the
// graphql.BuildSchema race.
func TestDeviceEntrySnapshot_PlanesIsolatedFromRegistryRebind(t *testing.T) {
	t.Parallel()

	reg := NewDeviceRegistry([]PlaneProvider{stubPlaneProvider{name: "first"}})
	reg.Register(DeviceInfo{Address: 0x10, Manufacturer: "Vaillant", DeviceID: "Y"})

	snap, ok := reg.LookupEntrySnapshot(0x10)
	if !ok {
		t.Fatalf("LookupEntrySnapshot(0x10) ok=false")
	}

	// Re-register with a new identity that triggers a merge path. Then
	// also wire a new provider for the next register so planes change.
	reg.RegisterProvider(stubPlaneProvider{name: "second"})
	reg.Register(DeviceInfo{Address: 0x10, Manufacturer: "Vaillant", DeviceID: "Y"})

	// snap captured BEFORE the second register — must remain "first".
	if len(snap.Planes) != 1 {
		t.Fatalf("snap.Planes len = %d; want 1 (snapshot captured before re-register)", len(snap.Planes))
	}
	if got := snap.Planes[0].Name(); got != "first" {
		t.Errorf("snap.Planes[0].Name = %q; want first (snapshot must be isolated from registry rebind)", got)
	}
}

// TestIterateSnapshots_PlanesRaceFree models the production race:
// concurrent BuildSchema-style iteration vs identity-merge writes.
// With the live-pointer Iterate path, this would race on entry.planes
// reassignment. With IterateSnapshots + Plane snapshot, no race.
func TestIterateSnapshots_PlanesRaceFree(t *testing.T) {
	t.Parallel()

	reg := NewDeviceRegistry([]PlaneProvider{stubPlaneProvider{name: "plane-1"}})
	reg.Register(DeviceInfo{Address: 0x10, Manufacturer: "Vaillant", DeviceID: "Z"})

	stop := atomic.Bool{}
	var wg sync.WaitGroup

	// Reader: iterates snapshots and reads Planes. Must not race.
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 200; i++ {
			reg.IterateSnapshots(func(snap DeviceEntrySnapshot) bool {
				for _, plane := range snap.Planes {
					_ = plane.Name()
					for _, method := range plane.Methods() {
						_ = method.Name()
					}
				}
				return true
			})
		}
		stop.Store(true)
	}()

	// Writer: continuously re-registers until reader signals stop.
	wg.Add(1)
	go func() {
		defer wg.Done()
		for !stop.Load() {
			reg.Register(DeviceInfo{Address: 0x10, Manufacturer: "Vaillant", DeviceID: "Z"})
		}
	}()

	wg.Wait()

	// Final snapshot must still have planes populated.
	snap, ok := reg.LookupEntrySnapshot(0x10)
	if !ok {
		t.Fatalf("LookupEntrySnapshot(0x10) ok=false")
	}
	if len(snap.Planes) == 0 {
		t.Errorf("final snap.Planes empty; want populated")
	}
}

// stubPlaneProvider creates immutable single-method planes for tests.
type stubPlaneProvider struct {
	name string
}

func (p stubPlaneProvider) Name() string                    { return p.name }
func (p stubPlaneProvider) Match(info DeviceInfo) bool      { return true }
func (p stubPlaneProvider) CreatePlanes(info DeviceInfo) []Plane {
	return []Plane{&stubPlane{name: p.name}}
}

type stubPlane struct {
	name string
}

func (p *stubPlane) Name() string             { return p.name }
func (p *stubPlane) Methods() []Method        { return []Method{stubMethod{name: "noop"}} }

type stubMethod struct {
	name string
}

func (m stubMethod) Name() string                   { return m.name }
func (m stubMethod) ReadOnly() bool                 { return true }
func (m stubMethod) Template() FrameTemplate        { return stubTemplate{} }
func (m stubMethod) ResponseSchema() schema.SchemaSelector {
	return schema.SchemaSelector{}
}

type stubTemplate struct{}

func (stubTemplate) Primary() byte   { return 0xB5 }
func (stubTemplate) Secondary() byte { return 0x09 }
