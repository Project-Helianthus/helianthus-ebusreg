package system

// B524ControllerIdentity identifies the controller context a B524 bound was
// observed under. Matching is exact on every field: a bound observed under one
// identity is never reused for another software revision, hardware revision,
// API version or revision, and an empty field only matches an empty field
// (the observation did not record it).
type B524ControllerIdentity struct {
	EID         string // e.g. "BASV0"
	SoftwareRaw string // raw software value, e.g. "0217"
	HardwareRaw string // raw hardware value, e.g. "8803"
	APIVersion  string // B524 API version, e.g. "1"; empty when not recorded
	Revision    string // API revision, e.g. "1"; empty when not recorded
}

// B524RegisterBound is one profile-scoped bound for a selector.
type B524RegisterBound struct {
	Codec string // "EXP" = 32-bit float
	Min   float64
	Max   float64
	Step  float64
	// AccessFlags are the observed description flags. They are not write
	// acceptance evidence.
	AccessFlags      byte
	AccessFlagsKnown bool
	// Qualified is true only when the bound comes from a retained,
	// selector-qualified description record. It is false for statements
	// without a retained record.
	Qualified bool
	// Evidence names the published source of the bound.
	Evidence string
}

type b524BoundSelector struct {
	Opcode   byte
	Group    byte
	Register uint16 // little-endian register number (RRlo | RRhi<<8)
}

// B524ControllerProfile is the set of bounds observed under one controller
// identity.
type B524ControllerProfile struct {
	Identity B524ControllerIdentity
	Bounds   map[b524BoundSelector]B524RegisterBound
}

const (
	basv0DHWBivalenceEvidence = "helianthus-docs-ebus protocols/vaillant/fixtures/b524-basv0-description-profile-v1.json, record OP02/GG00/II00/RR0001 (observation basv0_sw0217_hw8803_api1)"
)

var b524ControllerProfiles = []B524ControllerProfile{
	{
		Identity: B524ControllerIdentity{EID: "BASV0", SoftwareRaw: "0217", HardwareRaw: "8803", APIVersion: "1", Revision: "1"},
		Bounds: map[b524BoundSelector]B524RegisterBound{
			// system_dhw_bivalence_point: raw 0001000000A0C10000A0410000803F.
			{Opcode: B524OpcodeLocal, Group: 0x00, Register: 0x0001}: {
				Codec: "EXP", Min: -20, Max: 20, Step: 1,
				AccessFlags: 0x03, AccessFlagsKnown: true,
				Qualified: true, Evidence: basv0DHWBivalenceEvidence,
			},
		},
	},
	// No BASV2 profile is seeded. The docs state a BASV2 observation of
	// -20..50 step 1 for OP02/GG00/RR0001 (the former unscoped catalog value),
	// but no published source binds that bound to a controller identity
	// (software, hardware, API version, revision). An unmatched identity
	// yields B524BoundUnknownProfile.
}

// B524BoundStatus is the outcome of LookupB524Bound.
type B524BoundStatus int

const (
	// B524BoundUnknownProfile: no seeded profile has exactly this identity.
	B524BoundUnknownProfile B524BoundStatus = iota
	// B524BoundUnknownRegister: the profile is known but has no bound for the selector.
	B524BoundUnknownRegister
	// B524BoundFound: Bound is valid.
	B524BoundFound
)

func (status B524BoundStatus) String() string {
	switch status {
	case B524BoundUnknownProfile:
		return "unknown_profile"
	case B524BoundUnknownRegister:
		return "unknown_register"
	case B524BoundFound:
		return "found"
	default:
		return "invalid"
	}
}

// B524BoundResult carries the bound only when Status is B524BoundFound; every
// other status returns the zero Bound, never a substituted default.
type B524BoundResult struct {
	Status B524BoundStatus
	Bound  B524RegisterBound
}

// LookupB524Bound returns the bound for (controller identity, opcode, group,
// register). register is the little-endian register number (RRlo | RRhi<<8).
func LookupB524Bound(identity B524ControllerIdentity, opcode, group byte, register uint16) B524BoundResult {
	for i := range b524ControllerProfiles {
		profile := &b524ControllerProfiles[i]
		if profile.Identity != identity {
			continue
		}
		bound, ok := profile.Bounds[b524BoundSelector{Opcode: opcode, Group: group, Register: register}]
		if !ok {
			return B524BoundResult{Status: B524BoundUnknownRegister}
		}
		return B524BoundResult{Status: B524BoundFound, Bound: bound}
	}
	return B524BoundResult{Status: B524BoundUnknownProfile}
}
