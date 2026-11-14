package registry

import (
	"sync"

	"github.com/Project-Helianthus/helianthus-ebusreg/schema"
)

// DeviceInfo is one identification observation for a bus address, as parsed
// from an identification response or supplied by a seed. An empty string
// means the field was not observed: Register keeps the stored value for that
// field instead of clearing it. Only a complete manufacturer, device and
// serial identity may join or split entries across addresses.
type DeviceInfo struct {
	Address         byte
	Manufacturer    string
	DeviceID        string
	SerialNumber    string
	MacAddress      string
	SoftwareVersion string
	HardwareVersion string
}

// DeviceEntry is the read view of one device known to the registry, which
// may answer on several aliased addresses. Entries returned by a
// DeviceRegistry read through to live registry state, so consecutive calls
// are not a consistent snapshot; use the registry snapshot APIs for that.
// Identity accessors return "" for fields not yet observed.
type DeviceEntry interface {
	// AddressByRole returns the first BusFace address whose Role
	// matches the requested SlotRole. Returns (0, false) when no face
	// matches. Used by routing code to address the correct byte for
	// the intended frame type (per AddressClass taxonomy).
	//
	// Phase C M-C6c: replaces the previous ambiguous Address() method,
	// which conflated the "show me a representative byte" use case
	// (now PrimaryDisplayAddress) with the "give me the routing-
	// correct byte for this frame type" use case (this method). The
	// removed method silently returned the initiator byte for an
	// aliased canonical pair (e.g. BAI 0x03↔0x08), causing M2S writes
	// to mis-route to the initiator side. AddressByRole forces
	// callers to declare their intent.
	AddressByRole(role SlotRole) (byte, bool)
	// PrimaryDisplayAddress returns a representative address for log /
	// UI display. May be initiator OR target for aliased pairs; do
	// NOT use for wire routing. For routing, use AddressByRole.
	PrimaryDisplayAddress() byte
	Addresses() []byte
	Manufacturer() string
	DeviceID() string
	SerialNumber() string
	MacAddress() string
	SoftwareVersion() string
	HardwareVersion() string
	Planes() []Plane
	Projections() []Projection
}

// Plane is a named group of methods that a provider exposes for one device.
type Plane interface {
	Name() string
	Methods() []Method
}

// PlaneProvider attaches planes to the devices it recognises. For each
// registration that updates an entry, the registry calls Match and then
// CreatePlanes while holding its write lock, with the merged DeviceInfo of
// the entry; implementations must be fast and must not call back into the
// registry. A registration whose qualified identity conflicts with one bound
// elsewhere, at an address that already has an entry, leaves that entry
// unchanged and calls neither. Name identifies the provider.
type PlaneProvider interface {
	Name() string
	Match(info DeviceInfo) bool
	CreatePlanes(info DeviceInfo) []Plane
}

// ProjectionProvider is an optional interface of a PlaneProvider. For every
// matched provider that implements it, the registry passes the planes of all
// matched providers and keeps the returned projections only when together
// they build a valid canonical index; otherwise the entry has no projections
// and CanonicalIndexForEntry reports the index error.
type ProjectionProvider interface {
	CreateProjections(info DeviceInfo, planes []Plane) []Projection
}

// Method is one invocable operation of a plane. ReadOnly reports whether it
// leaves device state unchanged; Template gives the primary and secondary
// command bytes of its request, and ResponseSchema selects the schema that
// decodes its response. See ResolveMethodMetadata for the optional safety
// and routing overrides.
type Method interface {
	Name() string
	ReadOnly() bool
	Template() FrameTemplate
	ResponseSchema() schema.SchemaSelector
}

// FrameTemplate carries the primary and secondary command bytes of a request.
type FrameTemplate interface {
	Primary() byte
	Secondary() byte
}

// DeviceRegistry stores the devices seen on one bus, keyed by address, with
// per-address slots, alias groups and the identity index that joins
// addresses of the same physical device. Its methods may be called from
// several goroutines, but a live DeviceEntry returned by Lookup or passed
// to Iterate reads its identity fields without the registry lock: while
// another goroutine may register or alias devices, read entries through
// LookupEntrySnapshot or IterateSnapshots instead. The zero value is not
// usable; create one with NewDeviceRegistry.
type DeviceRegistry struct {
	mu                    sync.RWMutex
	observationGeneration uint64
	proofGeneration       uint64
	providers             []PlaneProvider
	entries               map[byte]*deviceEntry
	addressTable          [256]*AddressSlot
	identity              map[string]*deviceEntry
	qualifiedWitnesses    [256]qualifiedIdentityWitnessRecord
	// topology records explicit source-target or canonical-companion
	// relationships independently from qualified identity membership.
	// A complete identity may join independent entries, but only these
	// edges may carry current-session confirmation to another face.
	topology map[byte]map[byte]struct{}
	order    []*deviceEntry
}

// NewDeviceRegistry returns an empty registry that matches registrations
// against providers in slice order. The slice is copied, so the caller may
// reuse it.
func NewDeviceRegistry(providers []PlaneProvider) *DeviceRegistry {
	providerCopy := make([]PlaneProvider, len(providers))
	copy(providerCopy, providers)
	return &DeviceRegistry{
		providers: providerCopy,
		entries:   make(map[byte]*deviceEntry),
		identity:  make(map[string]*deviceEntry),
		topology:  make(map[byte]map[byte]struct{}),
	}
}

// RegisterProvider appends a provider and advances the observation
// generation. The provider applies to later registrations only: entries
// already stored keep their planes until their address is registered again.
func (r *DeviceRegistry) RegisterProvider(provider PlaneProvider) {
	r.mu.Lock()
	r.providers = append(r.providers, provider)
	r.observationGeneration++
	r.mu.Unlock()
}

// WithObservationGeneration executes fn while holding the registry read lock
// and supplies the current monotonic observation generation. Public registry
// mutations advance the generation under the corresponding write lock, so a
// writer linearizes either before the callback begins or after it returns.
//
// fn must remain bounded and must not call another DeviceRegistry method: the
// callback deliberately executes inside the read critical section so callers
// can atomically compare a captured generation and commit derived state before
// a relevant registry mutation can interleave.
//
// The method returns false without invoking fn when either receiver or callback
// is nil.
func (r *DeviceRegistry) WithObservationGeneration(fn func(uint64)) bool {
	if r == nil || fn == nil {
		return false
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	fn(r.observationGeneration)
	return true
}
