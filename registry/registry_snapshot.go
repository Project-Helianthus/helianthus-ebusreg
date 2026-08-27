package registry

import (
	"time"

	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
)

func (r *DeviceRegistry) Lookup(address byte) (DeviceEntry, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	entry, ok := r.entries[address]
	if !ok {
		return nil, false
	}
	return entry, true
}

// LookupSlot returns the AddressSlot for the requested address (M1
// address-table accessor), with the slot's own role/source/confidence
// metadata. When the address is aliased to a multi-address device, the
// returned slot.Device pointer is shared with the primary slot, but
// slot.Addr/Role/DiscoverySource/VerificationState describe the
// REQUESTED address — callers inspecting per-address metadata get the
// per-slot view (Codex P2: return-the-requested-address-slot).
func (r *DeviceRegistry) LookupSlot(address byte) (*AddressSlot, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	slot := r.addressTable[address]
	if slot == nil {
		return nil, false
	}
	return slot, true
}

// AddressSlotSnapshot is a value-typed copy of an AddressSlot's
// observable fields. Snapshots are taken under r.mu.RLock so callers
// can read the fields without holding any registry lock and without
// risking torn reads from concurrent writers.
//
// P8.1 — addresses the lock-free read advisory raised by the GitHub
// Codex bot on helianthus-ebusgateway PR #589: the gateway's
// AddressTable previously dereferenced a live AddressSlot pointer
// (returned by LookupSlot) outside the registry's RLock to read
// DiscoverySource / VerificationState. AddressSlotSnapshot eliminates
// that race surface — the value copy is immune to concurrent
// mutations because the writer must acquire r.mu.Lock() (which blocks
// behind the RLock taken in LookupSlotSnapshot below) before changing
// the underlying slot.
//
// Note: Device is reduced to a boolean (DeviceAttached) because
// returning the *deviceEntry pointer would re-introduce lock-free
// dereferencing of mutable identity fields downstream. Callers that
// need the entry's identity fields can fetch the entry via Lookup,
// but should be aware that the returned DeviceEntry interface still
// reads through to the registry's internal entry struct — those
// reads are not snapshot-isolated. A future entry-snapshot API may
// be added if the same pattern proves desirable for identity reads.
type AddressSlotSnapshot struct {
	Addr              byte
	Role              SlotRole
	DiscoverySource   DiscoverySource
	VerificationState VerificationState
	FirstObservedAt   time.Time
	LastObservedAt    time.Time
	DeviceAttached    bool
}

// LookupSlotSnapshot returns a value-typed snapshot of the AddressSlot
// at addr, taken under r.mu.RLock. Callers can read the snapshot
// fields without any registry-lock concerns. Returns (zero, false)
// when no slot exists for the address.
//
// P8.1 — the race-free counterpart to LookupSlot for callers that
// only need the slot's observable fields (gateway address-table
// projection, MCP/GraphQL surfaces). LookupSlot remains for callers
// that need the live pointer (e.g. registry-internal mutation paths
// holding the appropriate lock externally).
func (r *DeviceRegistry) LookupSlotSnapshot(address byte) (AddressSlotSnapshot, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	slot := r.addressTable[address]
	if slot == nil {
		return AddressSlotSnapshot{}, false
	}
	return AddressSlotSnapshot{
		Addr:              slot.Addr,
		Role:              slot.Role,
		DiscoverySource:   slot.DiscoverySource,
		VerificationState: slot.VerificationState,
		FirstObservedAt:   slot.FirstObservedAt,
		LastObservedAt:    slot.LastObservedAt,
		DeviceAttached:    slot.Device != nil,
	}, true
}

func (d *deviceEntry) PrimaryDisplayAddress() byte {
	if d.primaryAddress != 0 {
		return d.primaryAddress
	}
	return d.info.Address
}

// AddressByRole returns the first BusFace address whose Role matches.
// Routing-correct alternative to Address: M2S writers pass
// SlotRoleSlave to get the target byte; M2M arbitration logic passes
// SlotRoleMaster for the initiator byte. Returns (0, false) when no
// face matches the requested role.
//
// Decision references: Phase C AD30 (entry.Address ambiguity fix);
// uses the existing BusFace.Role machinery populated by
// syncEntryFacesLocked.
func (d *deviceEntry) AddressByRole(role SlotRole) (byte, bool) {
	if d == nil {
		return 0, false
	}
	// Pass 1: explicit Role match (set via MarkSlotPassiveObserved or
	// other role-aware paths).
	for _, face := range d.Faces {
		if face.Role == role {
			return face.Addr, true
		}
	}
	// Pass 2: SlotRoleUnknown fallback — active scan registers entries
	// without populating Role (Codex P2 from PR #134). Infer the role
	// from the address class so callers migrating from Address() to
	// AddressByRole get a useful answer for actively-scanned devices.
	for _, face := range d.Faces {
		if face.Role != SlotRoleUnknown {
			continue
		}
		switch protocol.AddressClassOf(face.Addr) {
		case protocol.AddressClassMaster:
			if role == SlotRoleMaster {
				return face.Addr, true
			}
		case protocol.AddressClassSlave:
			if role == SlotRoleSlave {
				return face.Addr, true
			}
		}
	}
	return 0, false
}

func (d *deviceEntry) Addresses() []byte {
	if len(d.addresses) == 0 {
		if d.info.Address == 0 {
			return nil
		}
		return []byte{d.info.Address}
	}
	out := make([]byte, len(d.addresses))
	copy(out, d.addresses)
	return out
}

func (d *deviceEntry) Manufacturer() string {
	return d.info.Manufacturer
}

func (d *deviceEntry) DeviceID() string {
	return d.info.DeviceID
}

func (d *deviceEntry) SerialNumber() string {
	return d.info.SerialNumber
}

func (d *deviceEntry) MacAddress() string {
	return d.info.MacAddress
}

func (d *deviceEntry) SoftwareVersion() string {
	return d.info.SoftwareVersion
}

func (d *deviceEntry) HardwareVersion() string {
	return d.info.HardwareVersion
}

func (d *deviceEntry) Planes() []Plane {
	return d.planes
}

func (d *deviceEntry) Projections() []Projection {
	return d.projections
}

// DeviceEntrySnapshot is a value-typed copy of a DeviceEntry's
// observable identity fields. Snapshots are taken under r.mu.RLock
// so callers can read the fields without holding any registry lock
// and without risking torn reads from concurrent writers
// (Register / RegisterStaticSeed / RegisterPassiveObserved /
// AliasAddresses / detachAddressLocked).
//
// P9 — addresses Codex post-P8.3 audit: the DeviceEntry interface
// methods (Manufacturer / DeviceID / SerialNumber / etc.) read
// `d.info.<Field>` lock-free. Concurrent Register replaces
// `entry.info` with a new DeviceInfo struct (line 240
// `entry.info = storedInfo`); a reader holding the *deviceEntry
// pointer can observe a torn read of the string fields
// (string is a 16-byte ptr+len header).
//
// DeviceEntrySnapshot copies all string and slice fields under
// RLock; the snapshot is disconnected from registry storage and
// safe to read concurrently. Slice copies (Addresses, Faces) prevent
// callers from mutating registry state through the snapshot.
//
// SCOPE (P9.x): Planes and Projections were originally omitted on the
// theory that their interface trees transitively exposed registry-
// mutable state. In practice the *plane and *method implementations
// (vaillant providers + similar) are constructed once in
// PlaneProvider.CreatePlanes / ProjectionProvider.CreateProjections
// and never mutated afterward; the only registry-side write to the
// `entry.planes` / `entry.projections` slice headers happens during
// the identity-merge path (mergeEntries dst.planes = src.planes /
// dst.projections = src.projections). Capturing those slice headers
// under RLock therefore produces a stable view: readers iterating
// the snapshot's Planes / Projections see element references that
// remain valid for the lifetime of the snapshot, with no risk of a
// mid-iteration slice reassignment.
//
// P9.x adds Planes + Projections to DeviceEntrySnapshot so the
// graphql.BuildSchema hot path can drop its live-pointer Iterate
// usage. Callers that mutate the registry (registerLocked path) must
// continue to use the live `*deviceEntry` pointer; the snapshot is
// READ-ONLY by construction.
type DeviceEntrySnapshot struct {
	PrimaryAddress  byte
	Addresses       []byte
	Faces           []BusFace
	Manufacturer    string
	DeviceID        string
	SerialNumber    string
	MacAddress      string
	SoftwareVersion string
	HardwareVersion string
	// Planes + Projections — slice headers captured under RLock at
	// snapshot time. The interface elements (Plane, Method, etc.) are
	// immutable after PlaneProvider.CreatePlanes returns; safe to
	// read after the lock has been released.
	Planes      []Plane
	Projections []Projection
}

// PrimaryDisplayAddress mirrors deviceEntry.PrimaryDisplayAddress for
// callers that work with the value-typed snapshot.
func (s DeviceEntrySnapshot) PrimaryDisplayAddress() byte {
	return s.PrimaryAddress
}

// AddressByRole mirrors deviceEntry.AddressByRole using the snapshot's
// Faces slice. Returns (0, false) when no face matches the requested
// role (after the same role-class fallback rules used by the live
// implementation).
func (s DeviceEntrySnapshot) AddressByRole(role SlotRole) (byte, bool) {
	for _, face := range s.Faces {
		if face.Role == role {
			return face.Addr, true
		}
	}
	for _, face := range s.Faces {
		if face.Role != SlotRoleUnknown {
			continue
		}
		switch protocol.AddressClassOf(face.Addr) {
		case protocol.AddressClassMaster:
			if role == SlotRoleMaster {
				return face.Addr, true
			}
		case protocol.AddressClassSlave:
			if role == SlotRoleSlave {
				return face.Addr, true
			}
		}
	}
	return 0, false
}

// LookupEntrySnapshot returns a value-typed snapshot of the
// DeviceEntry registered for addr, taken under r.mu.RLock. The
// race-free counterpart to Lookup for callers that only need the
// entry's observable identity fields (gateway MCP / GraphQL /
// observability projections). Returns (zero, false) when no entry
// exists for the address.
//
// P9 — Lookup remains for callers that need the live entry pointer
// (e.g. registry-internal mutation paths or callers that need
// Planes/Projections).
func (r *DeviceRegistry) LookupEntrySnapshot(address byte) (DeviceEntrySnapshot, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	entry := r.entries[address]
	if entry == nil {
		return DeviceEntrySnapshot{}, false
	}
	return r.snapshotEntryLocked(entry), true
}

// IterateSnapshots visits each registered device entry as a
// value-typed snapshot. Snapshots are built under r.mu.RLock; the
// lock is released BEFORE the callback runs. Callers can therefore
// safely call any DeviceRegistry method from within the callback
// (no deadlock risk).
//
// This matches the existing Iterate API's lock-then-snapshot-then-
// unlock contract. The only behavioural difference vs Iterate is
// the value-typed snapshot vs live entry pointer (Codex P9 review
// pass 1 MINOR FINDING_2).
//
// P9 — Iterate remains for callers that need live entry pointers
// (or Planes / Projections, which the snapshot intentionally
// omits). New consumers SHOULD prefer IterateSnapshots.
func (r *DeviceRegistry) IterateSnapshots(fn func(DeviceEntrySnapshot) bool) {
	r.mu.RLock()
	snapshots := make([]DeviceEntrySnapshot, 0, len(r.order))
	for _, entry := range r.order {
		if entry == nil {
			continue
		}
		snapshots = append(snapshots, r.snapshotEntryLocked(entry))
	}
	r.mu.RUnlock()

	for _, snap := range snapshots {
		if !fn(snap) {
			return
		}
	}
}

// snapshotEntryLocked builds a DeviceEntrySnapshot from a live
// *deviceEntry. Caller MUST hold r.mu.RLock or r.mu.Lock.
//
// Each BusFace is deep-copied: the AccessProtocols slice is
// copied per face, not aliased (Codex P9 review pass 1 MINOR
// FINDING_1). This guarantees the snapshot is fully disconnected
// from registry storage — consumer mutations of any nested slice
// do NOT leak through.
// deepCopyProjectionPathLocked returns a copy of `src` whose Segments
// slice header points to a fresh backing array. Caller does NOT need
// to hold any registry lock — operates on local copies only — but
// snapshot construction holds RLock to ensure atomicity with the rest
// of the snapshot. P9.x.
func deepCopyProjectionPathLocked(src ProjectionPath) ProjectionPath {
	out := ProjectionPath{Plane: src.Plane}
	if len(src.Segments) > 0 {
		out.Segments = make([]PathSegment, len(src.Segments))
		copy(out.Segments, src.Segments)
	}
	return out
}

// snapshotPlane wraps a registry-owned Plane so that Methods() returns
// a snapshot-owned copy of the underlying methods slice. Without this
// wrapper, the vaillant providers' plane.Methods() returns
// plane.methods directly — a snapshot caller could write to the
// returned slice and corrupt the live registry plane (Codex P9.x
// review pass 2 GitHub-bot finding). The snapshot wrapper preserves
// the Plane interface contract while shielding registry storage.
//
// The Method interface values in `methods` are shared with the
// registry; vaillant Method implementations are immutable struct
// values (verified in vaillant/system/system.go and analogous
// providers), so sharing them is safe. Only the SLICE itself needs
// to be snapshot-owned to prevent index-write corruption.
type snapshotPlane struct {
	name    string
	methods []Method
}

// Name returns the wrapped plane's name, captured at snapshot time.
func (p *snapshotPlane) Name() string { return p.name }

// Methods returns a fresh copy of the snapshot-owned method slice on
// every call. Callers therefore cannot mutate the slice another caller
// holds (or the registry's underlying slice).
func (p *snapshotPlane) Methods() []Method {
	out := make([]Method, len(p.methods))
	copy(out, p.methods)
	return out
}

// snapshotPlaneFromLocked builds a snapshotPlane from a registry-owned
// Plane. Caller MUST hold r.mu (read or write) so the underlying
// plane's Methods() result is captured atomically with the rest of the
// snapshot.
func snapshotPlaneFromLocked(p Plane) *snapshotPlane {
	src := p.Methods()
	methods := make([]Method, len(src))
	copy(methods, src)
	return &snapshotPlane{name: p.Name(), methods: methods}
}

func (r *DeviceRegistry) snapshotEntryLocked(entry *deviceEntry) DeviceEntrySnapshot {
	addresses := make([]byte, len(entry.addresses))
	copy(addresses, entry.addresses)
	faces := make([]BusFace, len(entry.Faces))
	for i, src := range entry.Faces {
		face := src
		if len(src.AccessProtocols) > 0 {
			face.AccessProtocols = make([]string, len(src.AccessProtocols))
			copy(face.AccessProtocols, src.AccessProtocols)
		}
		faces[i] = face
	}
	primary := entry.primaryAddress
	if primary == 0 {
		primary = entry.info.Address
	}
	// P9.x — capture Planes + Projections under RLock. The Plane
	// interface implementations are immutable after CreatePlanes
	// returns (vaillant providers verified), so copying the slice
	// header is sufficient there. Projections are concrete structs
	// (Plane string, Nodes []Node, Edges []Edge) where each Node has
	// a ProjectionPath.Segments []PathSegment slice; deep-copy those
	// nested slices so callers cannot mutate registry-visible state
	// through the snapshot (Codex P9.x review pass 1).
	var planes []Plane
	if len(entry.planes) > 0 {
		planes = make([]Plane, len(entry.planes))
		for i, src := range entry.planes {
			planes[i] = snapshotPlaneFromLocked(src)
		}
	}
	var projections []Projection
	if len(entry.projections) > 0 {
		projections = make([]Projection, len(entry.projections))
		for i, src := range entry.projections {
			proj := Projection{Plane: src.Plane}
			if len(src.Nodes) > 0 {
				proj.Nodes = make([]Node, len(src.Nodes))
				for j, srcNode := range src.Nodes {
					proj.Nodes[j] = Node{
						ID:            srcNode.ID,
						Path:          deepCopyProjectionPathLocked(srcNode.Path),
						CanonicalPath: deepCopyProjectionPathLocked(srcNode.CanonicalPath),
					}
				}
			}
			if len(src.Edges) > 0 {
				proj.Edges = make([]Edge, len(src.Edges))
				copy(proj.Edges, src.Edges) // Edge has only string-typed fields
			}
			projections[i] = proj
		}
	}
	return DeviceEntrySnapshot{
		PrimaryAddress:  primary,
		Addresses:       addresses,
		Faces:           faces,
		Manufacturer:    entry.info.Manufacturer,
		DeviceID:        entry.info.DeviceID,
		SerialNumber:    entry.info.SerialNumber,
		MacAddress:      entry.info.MacAddress,
		SoftwareVersion: entry.info.SoftwareVersion,
		HardwareVersion: entry.info.HardwareVersion,
		Planes:          planes,
		Projections:     projections,
	}
}
