package registry

import "time"

type SlotRole int

const (
	SlotRoleUnknown SlotRole = iota
	SlotRoleMaster
	SlotRoleSlave
)

type DiscoverySource int

const (
	DiscoverySourceUnknown DiscoverySource = iota
	DiscoverySourcePassiveObserved
	DiscoverySourceStaticSeed
	DiscoverySourceActiveConfirmed
)

type VerificationState int

const (
	VerificationStateUnknown VerificationState = iota
	VerificationStateCandidate
	VerificationStateCorroborated
	VerificationStateIdentityConfirmed
)

type AddressSlot struct {
	Addr              byte
	Role              SlotRole
	DiscoverySource   DiscoverySource
	VerificationState VerificationState
	Device            *deviceEntry
	FirstObservedAt   time.Time
	LastObservedAt    time.Time
}

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

func (s *AddressSlot) Addresses() []byte {
	if s == nil || s.Device == nil {
		return nil
	}
	return s.Device.Addresses()
}

func (s *AddressSlot) Manufacturer() string {
	if s == nil || s.Device == nil {
		return ""
	}
	return s.Device.Manufacturer()
}

func (s *AddressSlot) DeviceID() string {
	if s == nil || s.Device == nil {
		return ""
	}
	return s.Device.DeviceID()
}

func (s *AddressSlot) SerialNumber() string {
	if s == nil || s.Device == nil {
		return ""
	}
	return s.Device.SerialNumber()
}

func (s *AddressSlot) MacAddress() string {
	if s == nil || s.Device == nil {
		return ""
	}
	return s.Device.MacAddress()
}

func (s *AddressSlot) SoftwareVersion() string {
	if s == nil || s.Device == nil {
		return ""
	}
	return s.Device.SoftwareVersion()
}

func (s *AddressSlot) HardwareVersion() string {
	if s == nil || s.Device == nil {
		return ""
	}
	return s.Device.HardwareVersion()
}

func (s *AddressSlot) Planes() []Plane {
	if s == nil || s.Device == nil {
		return nil
	}
	return s.Device.Planes()
}

func (s *AddressSlot) Projections() []Projection {
	if s == nil || s.Device == nil {
		return nil
	}
	return s.Device.Projections()
}
