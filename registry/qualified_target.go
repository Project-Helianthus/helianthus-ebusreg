package registry

// TargetDeclaration is the answer of a thermal-target declaration rule for
// one identity. Declared is false with a named Reason when the rule does not
// declare the identity as a target. Flow fields name the declared flow row:
// its opcode, its selector, the publication element, and its evidence grade.
// FlowQualifiedBinding reports whether that row requires a qualified binding.
type TargetDeclaration struct {
	Declared             bool
	Reason               string
	Class                string
	Level                string
	RuleDigest           string
	FlowOpcode           byte
	FlowSelector         uint16
	FlowElement          string
	FlowEvidenceGrade    string
	FlowQualifiedBinding bool
}

// TargetDeclarationRule decides the thermal-target declaration of one
// identity. The caller supplies the rule; registry imports no vendor package.
// The rule must not call back into the registry.
type TargetDeclarationRule func(info DeviceInfo) TargetDeclaration

// QualifiedTarget is an immutable declared-target result for one exact
// current address. A declaration needs no serial, so the result can be
// current without a complete identity witness.
type QualifiedTarget struct {
	Address      byte
	Manufacturer string
	DeviceID     string
	Declaration  TargetDeclaration
}

// WithCurrentQualifiedTarget evaluates rule against the current identity of
// address, for the device the caller bound as declared (its manufacturer and
// device ID). It supplies the declared target to fn and returns true only
// when the address still resolves to that device, rule declares it, and any
// complete identity of the address has a current witness. It returns false
// and does not invoke fn otherwise, or when r, rule or fn is nil.
//
// A result retires when the address no longer resolves to the bound device:
// the entry is gone, its manufacturer or device ID changed, or its complete
// identity has no current witness. The evaluation runs under the registry
// read lock and is not stored.
func (r *DeviceRegistry) WithCurrentQualifiedTarget(address byte, declared DeviceInfo, rule TargetDeclarationRule, fn func(QualifiedTarget)) bool {
	if r == nil || rule == nil || fn == nil {
		return false
	}
	r.mu.RLock()
	defer r.mu.RUnlock()

	entry := r.entries[address]
	if entry == nil {
		return false
	}
	if normalizeIdentityPart(entry.info.Manufacturer) != normalizeIdentityPart(declared.Manufacturer) ||
		normalizeIdentityPart(entry.info.DeviceID) != normalizeIdentityPart(declared.DeviceID) {
		return false
	}
	if entry.identityKey != "" && !r.qualifiedIdentityWitnessCurrentLocked(address, r.qualifiedWitnesses[address].witness) {
		return false
	}
	declaration := rule(entry.info)
	if !declaration.Declared {
		return false
	}
	fn(QualifiedTarget{
		Address:      address,
		Manufacturer: entry.info.Manufacturer,
		DeviceID:     entry.info.DeviceID,
		Declaration:  declaration,
	})
	return true
}
