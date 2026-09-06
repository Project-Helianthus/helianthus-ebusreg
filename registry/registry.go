package registry

import (
	"sync"

	"github.com/Project-Helianthus/helianthus-ebusreg/schema"
)

type DeviceInfo struct {
	Address         byte
	Manufacturer    string
	DeviceID        string
	SerialNumber    string
	MacAddress      string
	SoftwareVersion string
	HardwareVersion string
}

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

type Plane interface {
	Name() string
	Methods() []Method
}

type PlaneProvider interface {
	Name() string
	Match(info DeviceInfo) bool
	CreatePlanes(info DeviceInfo) []Plane
}

type ProjectionProvider interface {
	CreateProjections(info DeviceInfo, planes []Plane) []Projection
}

type Method interface {
	Name() string
	ReadOnly() bool
	Template() FrameTemplate
	ResponseSchema() schema.SchemaSelector
}

type FrameTemplate interface {
	Primary() byte
	Secondary() byte
}

type DeviceRegistry struct {
	mu                    sync.RWMutex
	observationGeneration uint64
	providers             []PlaneProvider
	entries               map[byte]*deviceEntry
	addressTable          [256]*AddressSlot
	identity              map[string]*deviceEntry
	// topology records explicit source-target or canonical-companion
	// relationships independently from qualified identity membership.
	// A complete identity may join independent entries, but only these
	// edges may carry current-session confirmation to another face.
	topology map[byte]map[byte]struct{}
	order    []*deviceEntry
}

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
