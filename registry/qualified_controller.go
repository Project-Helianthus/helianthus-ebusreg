package registry

// DeclaredRow is one native row that a declaration rule lists for a device
// identity. Element is the publication name the documentation gives the row
// and is empty when the documentation names none.
type DeclaredRow struct {
	Group         byte
	Register      uint16
	Element       string
	EvidenceGrade string
	Level         string
}

// ControllerDeclaration is the answer of a controller declaration rule for
// one identity. Declared is false with a named Reason when the rule does not
// list the identity. Family is the device identity the rule declares,
// Level is the recorded hardware test level, and RuleDigest identifies the
// rule document that produced the answer.
type ControllerDeclaration struct {
	Declared      bool
	Reason        string
	Family        string
	ProductNumber string
	Level         string
	EvidenceGrade string
	RuleDigest    string
	Rows          []DeclaredRow
}

// ControllerDeclarationRule decides the controller declaration of one
// identity. The caller supplies the rule; registry imports no vendor package.
// The rule must not call back into the registry.
type ControllerDeclarationRule func(info DeviceInfo) ControllerDeclaration

// QualifiedController is an immutable declared-controller result for one
// exact current address. Its identity fields are the raw values the
// registry holds: SoftwareVersion and ProductNumberRaw are the hexadecimal
// form of identification bytes 6-7 and 8-9. The result carries the
// generations of the current identity witness and retires with that witness.
type QualifiedController struct {
	Address                       byte
	Manufacturer                  string
	DeviceID                      string
	SoftwareVersion               string
	ProductNumberRaw              string
	Declaration                   ControllerDeclaration
	RegistryObservationGeneration uint64
	RegistryProofGeneration       uint64
}

// WithCurrentQualifiedController evaluates rule against the current identity
// of address and supplies the declared result to fn. It returns false and
// does not invoke fn when the address has no entry, when the rule does not
// declare the identity, when the address has a complete identity whose
// current witness has ended, or when r, rule or fn is nil.
//
// An address whose identity is not complete (no serial) has no witness; the
// family rule needs no serial, so the result is declared with zero witness
// generations. When a current witness exists, the result carries its
// generations and retires with it.
//
// The evaluation runs under the registry read lock, so a replacement or
// retirement cannot interleave between validation and use. A result is
// computed from the registry state at call time and is not stored.
func (r *DeviceRegistry) WithCurrentQualifiedController(address byte, rule ControllerDeclarationRule, fn func(QualifiedController)) bool {
	if r == nil || rule == nil || fn == nil {
		return false
	}
	r.mu.RLock()
	defer r.mu.RUnlock()

	entry := r.entries[address]
	if entry == nil {
		return false
	}
	var witness QualifiedIdentityWitness
	if entry.identityKey != "" {
		witness = r.qualifiedWitnesses[address].witness
		if !r.qualifiedIdentityWitnessCurrentLocked(address, witness) {
			return false
		}
	}
	declaration := rule(entry.info)
	if !declaration.Declared {
		return false
	}
	declaration.Rows = append([]DeclaredRow(nil), declaration.Rows...)
	fn(QualifiedController{
		Address:                       address,
		Manufacturer:                  entry.info.Manufacturer,
		DeviceID:                      entry.info.DeviceID,
		SoftwareVersion:               entry.info.SoftwareVersion,
		ProductNumberRaw:              entry.info.HardwareVersion,
		Declaration:                   declaration,
		RegistryObservationGeneration: witness.RegistryObservationGeneration,
		RegistryProofGeneration:       witness.RegistryProofGeneration,
	})
	return true
}
