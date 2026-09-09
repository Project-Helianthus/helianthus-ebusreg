package registry

import "time"

const (
	qualifiedVaillantControllerRuleRevision = "vaillant-controller-qualification/v1"
	qualifiedVaillantControllerEvidenceRef  = "https://github.com/Project-Helianthus/helianthus-docs-ebus/blob/055738bbad31f5cfe7fcd87bffc16bc09e021571/protocols/vaillant/ebus-vaillant-b555-timer-protocol.md"
	qualifiedVaillantControllerEvidenceSHA  = "2eff947174e4eb163b8e38f5eb93deb795491ec084e3e89ebc1dcdb7431e12ff"
)

// NativeEvidenceReference names immutable public protocol evidence. It is a
// value so callers cannot mutate registry-owned evidence.
type NativeEvidenceReference struct {
	Kind      string
	Reference string
	SHA256    string
}

// QualifiedVaillantController is the immutable-by-value, exact-address
// authority for the one published BASV2 thermal-map applicability tuple. It
// is native registry evidence, not a semantic or operation authority.
type QualifiedVaillantController struct {
	Address                       byte
	IdentityAuthority             QualifiedIdentityAuthority
	Product                       string
	SoftwareVersion               string
	HardwareVersion               string
	ObservationProvenance         QualifiedIdentityObservationProvenance
	ObservedAt                    time.Time
	RuleRevision                  string
	NativeEvidence                NativeEvidenceReference
	Current                       bool
	Immutable                     bool
	RegistryObservationGeneration uint64
	RegistryProofGeneration       uint64
}

type qualifiedVaillantControllerRecord struct {
	controller QualifiedVaillantController
	witness    QualifiedIdentityWitness
}

// CurrentQualifiedVaillantController returns a value copy only while the
// exact address has current direct identity authority and the documented
// BASV2/SW0507/HW1704 applicability tuple. It does not run a callback under a
// registry lock, so callers cannot mutate or re-enter the registry through a
// read-lock callback.
func (r *DeviceRegistry) CurrentQualifiedVaillantController(address byte) (QualifiedVaillantController, bool) {
	if r == nil {
		return QualifiedVaillantController{}, false
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	controller := r.qualifiedVaillantControllers[address].controller
	if !r.qualifiedVaillantControllerCurrentLocked(address, controller) {
		return QualifiedVaillantController{}, false
	}
	return controller, true
}

// recordDirectQualifiedVaillantControllerLocked records only an incoming
// complete direct observation. Stored fields and sparse LKG refreshes cannot
// synthesize a product/firmware proof. Caller holds r.mu.
func (r *DeviceRegistry) recordDirectQualifiedVaillantControllerLocked(info DeviceInfo) {
	if !matchesQualifiedVaillantControllerTuple(info) {
		return
	}
	witness := r.qualifiedWitnesses[info.Address].witness
	if !r.qualifiedIdentityWitnessCurrentLocked(info.Address, witness) {
		return
	}
	r.proofGeneration++
	r.qualifiedVaillantControllers[info.Address] = qualifiedVaillantControllerRecord{
		witness: witness,
		controller: QualifiedVaillantController{
			Address:                       info.Address,
			IdentityAuthority:             witness.IdentityAuthority,
			Product:                       "BASV2",
			SoftwareVersion:               "0507",
			HardwareVersion:               "1704",
			ObservationProvenance:         QualifiedIdentityObservationProvenanceDirect,
			ObservedAt:                    time.Now().UTC(),
			RuleRevision:                  qualifiedVaillantControllerRuleRevision,
			NativeEvidence:                NativeEvidenceReference{Kind: "protocol.observation", Reference: qualifiedVaillantControllerEvidenceRef, SHA256: qualifiedVaillantControllerEvidenceSHA},
			Current:                       true,
			Immutable:                     true,
			RegistryObservationGeneration: r.observationGeneration,
			RegistryProofGeneration:       r.proofGeneration,
		},
	}
}

func matchesQualifiedVaillantControllerTuple(info DeviceInfo) bool {
	identity := canonicalPhysicalIdentity(info)
	return isQualifiedWitnessAddress(info.Address) && identity.isQualified() &&
		identity.manufacturer == "VAILLANT" && identity.deviceID == "BASV2" &&
		identity.softwareVersion == "0507" && identity.hardwareVersion == "1704"
}

func (r *DeviceRegistry) reconcileQualifiedVaillantControllersLocked() {
	for address, record := range r.qualifiedVaillantControllers {
		if !r.qualifiedVaillantControllerCurrentLocked(byte(address), record.controller) {
			r.qualifiedVaillantControllers[address] = qualifiedVaillantControllerRecord{}
		}
	}
}

func (r *DeviceRegistry) retireQualifiedVaillantControllerLocked(address byte) {
	if r.qualifiedVaillantControllers[address].controller.Current {
		r.proofGeneration++
		r.qualifiedVaillantControllers[address] = qualifiedVaillantControllerRecord{}
	}
}

func (r *DeviceRegistry) qualifiedVaillantControllerCurrentLocked(address byte, controller QualifiedVaillantController) bool {
	if controller.Address != address || !controller.Current || !controller.Immutable ||
		controller.ObservationProvenance != QualifiedIdentityObservationProvenanceDirect ||
		controller.ObservedAt.IsZero() || controller.RuleRevision != qualifiedVaillantControllerRuleRevision ||
		controller.NativeEvidence != (NativeEvidenceReference{Kind: "protocol.observation", Reference: qualifiedVaillantControllerEvidenceRef, SHA256: qualifiedVaillantControllerEvidenceSHA}) ||
		controller.RegistryObservationGeneration == 0 || controller.RegistryProofGeneration == 0 {
		return false
	}
	record := r.qualifiedVaillantControllers[address]
	if record.controller != controller || !r.qualifiedIdentityWitnessCurrentLocked(address, record.witness) ||
		record.witness.IdentityAuthority != controller.IdentityAuthority {
		return false
	}
	entry := r.entries[address]
	if entry == nil {
		return false
	}
	return matchesQualifiedVaillantControllerTuple(DeviceInfo{
		Address: address, Manufacturer: entry.info.Manufacturer, DeviceID: entry.info.DeviceID,
		SerialNumber: entry.info.SerialNumber, SoftwareVersion: entry.info.SoftwareVersion,
		HardwareVersion: entry.info.HardwareVersion,
	}) && controller.Product == "BASV2" && controller.SoftwareVersion == "0507" && controller.HardwareVersion == "1704"
}
