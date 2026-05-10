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

// TestDeviceEntrySnapshot_PlaneMethodsIsolated verifies that the
// Method slice returned by snap.Planes[i].Methods() is snapshot-owned:
// mutating it must NOT corrupt registry-visible state on a subsequent
// LookupEntrySnapshot. Codex P9.x pass 2 GitHub-bot finding — vaillant
// providers' Plane.Methods() returns the live methods slice directly,
// so a snapshot caller writing to the result would leak into registry
// storage without the snapshotPlane wrapper.
func TestDeviceEntrySnapshot_PlaneMethodsIsolated(t *testing.T) {
	t.Parallel()

	reg := NewDeviceRegistry([]PlaneProvider{stubPlaneProvider{name: "stub"}})
	reg.Register(DeviceInfo{Address: 0x10, Manufacturer: "Vaillant", DeviceID: "M"})

	snap1, ok := reg.LookupEntrySnapshot(0x10)
	if !ok {
		t.Fatalf("LookupEntrySnapshot(0x10) ok=false")
	}
	if len(snap1.Planes) != 1 {
		t.Fatalf("snap1.Planes len = %d; want 1", len(snap1.Planes))
	}

	// Mutate the slice returned by Methods() on the FIRST snapshot.
	methods1 := snap1.Planes[0].Methods()
	if len(methods1) == 0 {
		t.Fatalf("snap1 Methods() empty; want >=1")
	}
	methods1[0] = nil // try to corrupt registry storage via snapshot

	// A subsequent snapshot must observe the registry's original
	// methods, not the nil we wrote.
	snap2, _ := reg.LookupEntrySnapshot(0x10)
	methods2 := snap2.Planes[0].Methods()
	if len(methods2) == 0 {
		t.Fatalf("snap2 Methods() empty; want >=1")
	}
	if methods2[0] == nil {
		t.Errorf("snap2 Methods()[0] = nil — snapshot mutation leaked into registry")
	}

	// Two calls to the SAME snapshot's Methods() must return
	// independent slices (no cross-mutation between callers).
	methodsA := snap1.Planes[0].Methods()
	methodsB := snap1.Planes[0].Methods()
	if len(methodsA) > 0 && len(methodsB) > 0 {
		methodsA[0] = nil
		if methodsB[0] == nil {
			t.Errorf("Methods() returned aliased slice — caller A's mutation visible to caller B")
		}
	}
}

// TestDeviceEntrySnapshot_ProjectionsDeepCopied verifies that mutating
// the snapshot's nested Projection.Nodes[i].Path.Segments slice does
// NOT affect the registry's state — Codex P9.x pass 1 finding.
func TestDeviceEntrySnapshot_ProjectionsDeepCopied(t *testing.T) {
	t.Parallel()

	reg := NewDeviceRegistry([]PlaneProvider{stubPlaneProviderWithProjection{name: "stub"}})
	reg.Register(DeviceInfo{Address: 0x10, Manufacturer: "Vaillant", DeviceID: "P"})

	snap, ok := reg.LookupEntrySnapshot(0x10)
	if !ok {
		t.Fatalf("LookupEntrySnapshot(0x10) ok=false")
	}
	if len(snap.Projections) != 1 {
		t.Fatalf("snap.Projections len = %d; want 1", len(snap.Projections))
	}
	if len(snap.Projections[0].Nodes) == 0 {
		t.Fatalf("snap.Projections[0].Nodes is empty; want >=1")
	}

	// Mutate the snapshot's Path.Segments — must NOT leak into registry.
	snap.Projections[0].Nodes[0].Path.Segments[0].Name = "MUTATED"

	snap2, ok := reg.LookupEntrySnapshot(0x10)
	if !ok {
		t.Fatalf("LookupEntrySnapshot(0x10) #2 ok=false")
	}
	if got := snap2.Projections[0].Nodes[0].Path.Segments[0].Name; got == "MUTATED" {
		t.Errorf("registry observed mutation through snapshot: Path.Segments[0].Name = %q (want unmodified)", got)
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

// stubPlaneProviderWithProjection adds a projection so the snapshot
// has Nodes/Edges/Segments to verify the deep-copy.
type stubPlaneProviderWithProjection struct {
	name string
}

func (p stubPlaneProviderWithProjection) Name() string               { return p.name }
func (p stubPlaneProviderWithProjection) Match(info DeviceInfo) bool { return true }
func (p stubPlaneProviderWithProjection) CreatePlanes(info DeviceInfo) []Plane {
	return []Plane{&stubPlane{name: p.name}}
}
func (p stubPlaneProviderWithProjection) CreateProjections(info DeviceInfo, planes []Plane) []Projection {
	return []Projection{{
		Plane: ServicePlane,
		Nodes: []Node{{
			ID:            NodeID(ServicePlane + ":/seg"),
			Path:          ProjectionPath{Plane: ServicePlane, Segments: []PathSegment{{Name: "seg"}}},
			CanonicalPath: ProjectionPath{Plane: ServicePlane, Segments: []PathSegment{{Name: "seg"}}},
		}},
	}}
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
