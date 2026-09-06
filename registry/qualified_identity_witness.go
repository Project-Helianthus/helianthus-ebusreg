package registry

import "github.com/Project-Helianthus/helianthus-ebusgo/protocol"

// QualifiedIdentityAuthority is the complete normalized identity authority
// carried by a QualifiedIdentityWitness. It is a value, not a registry entry
// or identity-index handle.
type QualifiedIdentityAuthority struct {
	Manufacturer string
	DeviceID     string
	SerialNumber string
}

// QualifiedIdentityObservationProvenance identifies the only observation
// provenance that can produce a qualified identity witness.
type QualifiedIdentityObservationProvenance string

const (
	// QualifiedIdentityObservationProvenanceDirect means the authority came
	// from a direct complete observation of the witness's exact address.
	QualifiedIdentityObservationProvenanceDirect QualifiedIdentityObservationProvenance = "direct_observation"
)

// QualifiedIdentityWitness is an immutable-by-value, registry-produced
// authority for one exact eBUS address. Current and Immutable are always true
// for a value supplied by WithCurrentQualifiedIdentityWitness; they make the
// closed public-contract properties explicit without providing an attestation
// or a caller-controlled currentness API.
//
// A copied or mutated value is only a caller-local value. It cannot alter the
// registry and cannot be submitted back as authority: currentness is checked
// by the registry inside WithCurrentQualifiedIdentityWitness.
type QualifiedIdentityWitness struct {
	Address                       byte
	IdentityAuthority             QualifiedIdentityAuthority
	ObservationProvenance         QualifiedIdentityObservationProvenance
	Current                       bool
	Immutable                     bool
	RegistryObservationGeneration uint64
	RegistryProofGeneration       uint64
}

type qualifiedIdentityWitnessRecord struct {
	witness      QualifiedIdentityWitness
	authorityKey string
}

// WithCurrentQualifiedIdentityWitness atomically validates and supplies the
// current qualified identity witness for address. It returns false and does
// not invoke fn when the exact address has no current witness, or when r or fn
// is nil.
//
// fn runs while the registry read lock is held. It must remain bounded and
// must not call DeviceRegistry methods. The lock makes lookup, currentness
// validation, and the consumer's use one boundary: a replacement, retirement,
// or conflicting authority cannot interleave after validation and before use.
// The supplied witness is a value with no mutable registry internals.
func (r *DeviceRegistry) WithCurrentQualifiedIdentityWitness(address byte, fn func(QualifiedIdentityWitness)) bool {
	if r == nil || fn == nil {
		return false
	}
	r.mu.RLock()
	defer r.mu.RUnlock()

	witness := r.qualifiedWitnesses[address].witness
	if !r.qualifiedIdentityWitnessCurrentLocked(address, witness) {
		return false
	}
	fn(witness)
	return true
}

// recordDirectQualifiedIdentityWitnessLocked publishes a fresh witness only
// for Register's direct, complete observation at its exact address. Caller
// holds r.mu and has already advanced observationGeneration for this mutation.
func (r *DeviceRegistry) recordDirectQualifiedIdentityWitnessLocked(info DeviceInfo, entry *deviceEntry) {
	authority := canonicalPhysicalIdentity(info)
	authorityKey := authority.key()
	if !isQualifiedWitnessAddress(info.Address) || authorityKey == "" || entry == nil ||
		r.entries[info.Address] != entry || !isCurrentQualifiedIdentityKey(entry, authorityKey) ||
		r.identity[authorityKey] != entry {
		return
	}

	r.proofGeneration++
	r.qualifiedWitnesses[info.Address] = qualifiedIdentityWitnessRecord{
		authorityKey: authorityKey,
		witness: QualifiedIdentityWitness{
			Address: info.Address,
			IdentityAuthority: QualifiedIdentityAuthority{
				Manufacturer: authority.manufacturer,
				DeviceID:     authority.deviceID,
				SerialNumber: authority.serialNumber,
			},
			ObservationProvenance:         QualifiedIdentityObservationProvenanceDirect,
			Current:                       true,
			Immutable:                     true,
			RegistryObservationGeneration: r.observationGeneration,
			RegistryProofGeneration:       r.proofGeneration,
		},
	}
}

// reconcileQualifiedIdentityWitnessesLocked retires records whose exact
// address no longer resolves to their direct qualified authority. It is shared
// by lifecycle paths that can replace, split, merge, or otherwise invalidate
// an address binding; sparse observations with unchanged authority retain the
// original direct proof and its generations.
func (r *DeviceRegistry) reconcileQualifiedIdentityWitnessesLocked() {
	for address, record := range r.qualifiedWitnesses {
		if record.authorityKey == "" {
			continue
		}
		if !r.qualifiedIdentityWitnessCurrentLocked(byte(address), record.witness) {
			r.qualifiedWitnesses[address] = qualifiedIdentityWitnessRecord{}
		}
	}
}

// qualifiedIdentityWitnessCurrentLocked is the shared currentness rule for
// public lookup and lifecycle reconciliation. Caller holds r.mu for reading
// (or writing). It intentionally compares the whole value, exact address,
// authority, and both generations rather than trusting positive generations
// or the cached Current flag.
func (r *DeviceRegistry) qualifiedIdentityWitnessCurrentLocked(address byte, witness QualifiedIdentityWitness) bool {
	if witness.Address != address || !witness.Current || !witness.Immutable ||
		witness.ObservationProvenance != QualifiedIdentityObservationProvenanceDirect ||
		witness.RegistryObservationGeneration == 0 || witness.RegistryProofGeneration == 0 ||
		!isQualifiedWitnessAddress(address) {
		return false
	}

	authority := canonicalPhysicalIdentity(DeviceInfo{
		Manufacturer: witness.IdentityAuthority.Manufacturer,
		DeviceID:     witness.IdentityAuthority.DeviceID,
		SerialNumber: witness.IdentityAuthority.SerialNumber,
	})
	authorityKey := authority.key()
	if authorityKey == "" || authority.manufacturer != witness.IdentityAuthority.Manufacturer ||
		authority.deviceID != witness.IdentityAuthority.DeviceID || authority.serialNumber != witness.IdentityAuthority.SerialNumber {
		return false
	}

	record := r.qualifiedWitnesses[address]
	if record.authorityKey != authorityKey || record.witness != witness {
		return false
	}
	entry := r.entries[address]
	return isCurrentQualifiedIdentityKey(entry, authorityKey) && r.identity[authorityKey] == entry
}

func isQualifiedWitnessAddress(address byte) bool {
	class := protocol.AddressClassOf(address)
	return class == protocol.AddressClassMaster || class == protocol.AddressClassSlave
}
