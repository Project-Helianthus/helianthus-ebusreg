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

func (s *AddressSlot) Address() byte {
	if s == nil {
		return 0
	}
	if s.Device != nil {
		return s.Device.Address()
	}
	return s.Addr
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
