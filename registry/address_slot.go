package registry

import "time"

// SlotRole is the bus role of one address face: an initiator address that
// may start frames, or a target address that only answers them.
type SlotRole int

// SlotRole values. SlotRoleUnknown is the zero value and means the role has
// not been established; it is not a guard value, and AddressByRole called
// with it matches faces whose role is still unknown.
const (
	SlotRoleUnknown SlotRole = iota
	SlotRoleMaster
	SlotRoleSlave
)

// DiscoverySource identifies how one address face first entered the registry.
// Its numeric order is not a confidence order; once non-unknown it is
// immutable for that face. VerificationState carries advancing confidence.
type DiscoverySource int

// DiscoverySource values, from the unset zero value to the three native
// sources a face can be stamped with.
const (
	DiscoverySourceUnknown DiscoverySource = iota
	DiscoverySourcePassiveObserved
	DiscoverySourceStaticSeed
	DiscoverySourceActiveConfirmed
)

// VerificationState is the advancing confidence that an address face belongs
// to its device entry. The numeric order is a confidence order: registry
// paths only raise it, from Candidate through Corroborated to
// IdentityConfirmed, and never lower it for a retained face.
type VerificationState int

// VerificationState values in ascending confidence; the zero value is
// VerificationStateUnknown.
const (
	VerificationStateUnknown VerificationState = iota
	VerificationStateCandidate
	VerificationStateCorroborated
	VerificationStateIdentityConfirmed
)

// AddressSlot is the registry's per-address record: the face's own role,
// provenance and observation times plus the device entry it currently
// resolves to. Device is nil while the address is known without an attached
// identity, and several slots share one Device when addresses are aliased.
//
// A slot returned by DeviceRegistry.LookupSlot is live registry state that
// writers mutate under the registry lock; readers outside the registry use
// LookupSlotSnapshot for a race-free copy. FirstObservedAt and
// LastObservedAt are zero until an observation or seed time is recorded.
type AddressSlot struct {
	Addr              byte
	Role              SlotRole
	DiscoverySource   DiscoverySource
	VerificationState VerificationState
	Device            *deviceEntry
	FirstObservedAt   time.Time
	LastObservedAt    time.Time
}

// BusFace is a value copy of one address face of a device entry: the
// address with its slot role, discovery source and verification state. The
// registry rebuilds an entry's faces whenever one of its slots or aliases
// changes. The registry itself leaves AccessProtocols empty; snapshots copy
// whatever a caller recorded there.
type BusFace struct {
	Addr              byte
	Role              SlotRole
	DiscoverySource   DiscoverySource
	VerificationState VerificationState
	AccessProtocols   []string
}

// PrimaryDisplayAddress returns the slot's display address —
// the wrapped Device's PrimaryDisplayAddress when present, otherwise
// the slot's own Addr. Phase C M-C6c: replaces the removed
// AddressSlot.Address() method (which had identical semantics but
// conflated display vs. routing intent).
func (s *AddressSlot) PrimaryDisplayAddress() byte {
	if s == nil {
		return 0
	}
	if s.Device != nil {
		return s.Device.PrimaryDisplayAddress()
	}
	return s.Addr
}

// AddressByRole forwards to the wrapped Device's AddressByRole; falls
// back to s.Addr matching when Device is nil and the slot's own Role
// matches the requested role.
func (s *AddressSlot) AddressByRole(role SlotRole) (byte, bool) {
	if s == nil {
		return 0, false
	}
	if s.Device != nil {
		return s.Device.AddressByRole(role)
	}
	if s.Role == role {
		return s.Addr, true
	}
	return 0, false
}

// Addresses returns every address of the attached device entry, or nil when
// the slot or its Device is nil.
func (s *AddressSlot) Addresses() []byte {
	if s == nil || s.Device == nil {
		return nil
	}
	return s.Device.Addresses()
}

// Manufacturer returns the attached device's manufacturer, or "" when no
// device is attached.
func (s *AddressSlot) Manufacturer() string {
	if s == nil || s.Device == nil {
		return ""
	}
	return s.Device.Manufacturer()
}

// DeviceID returns the attached device's identifier, or "" when no device is
// attached.
func (s *AddressSlot) DeviceID() string {
	if s == nil || s.Device == nil {
		return ""
	}
	return s.Device.DeviceID()
}

// SerialNumber returns the attached device's serial number, or "" when no
// device is attached.
func (s *AddressSlot) SerialNumber() string {
	if s == nil || s.Device == nil {
		return ""
	}
	return s.Device.SerialNumber()
}

// MacAddress returns the attached device's MAC address, or "" when no device
// is attached.
func (s *AddressSlot) MacAddress() string {
	if s == nil || s.Device == nil {
		return ""
	}
	return s.Device.MacAddress()
}

// SoftwareVersion returns the attached device's software version, or "" when
// no device is attached.
func (s *AddressSlot) SoftwareVersion() string {
	if s == nil || s.Device == nil {
		return ""
	}
	return s.Device.SoftwareVersion()
}

// HardwareVersion returns the attached device's hardware version, or "" when
// no device is attached.
func (s *AddressSlot) HardwareVersion() string {
	if s == nil || s.Device == nil {
		return ""
	}
	return s.Device.HardwareVersion()
}

// Planes returns the attached device's planes, or nil when no device is
// attached.
func (s *AddressSlot) Planes() []Plane {
	if s == nil || s.Device == nil {
		return nil
	}
	return s.Device.Planes()
}

// Projections returns the attached device's projections, or nil when no
// device is attached.
func (s *AddressSlot) Projections() []Projection {
	if s == nil || s.Device == nil {
		return nil
	}
	return s.Device.Projections()
}
