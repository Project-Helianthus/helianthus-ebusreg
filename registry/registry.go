package registry

import (
	"sync"
	"time"

	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
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
	mu           sync.RWMutex
	providers    []PlaneProvider
	entries      map[byte]*deviceEntry
	addressTable [256]*AddressSlot
	identity     map[string]*deviceEntry
	order        []*deviceEntry
}

func NewDeviceRegistry(providers []PlaneProvider) *DeviceRegistry {
	providerCopy := make([]PlaneProvider, len(providers))
	copy(providerCopy, providers)
	return &DeviceRegistry{
		providers: providerCopy,
		entries:   make(map[byte]*deviceEntry),
		identity:  make(map[string]*deviceEntry),
	}
}

func (r *DeviceRegistry) RegisterProvider(provider PlaneProvider) {
	r.mu.Lock()
	r.providers = append(r.providers, provider)
	r.mu.Unlock()
}

func (r *DeviceRegistry) Register(info DeviceInfo) DeviceEntry {
	r.mu.Lock()
	defer r.mu.Unlock()

	physical := canonicalPhysicalIdentity(info)
	identityKey := physical.key()
	planes := make([]Plane, 0)
	matched := make([]PlaneProvider, 0, len(r.providers))

	existingByAddress := r.entries[info.Address]
	incomingHasStableIdentity := normalizeIdentityPart(info.SerialNumber) != "" || normalizeIdentityPart(info.MacAddress) != ""
	if !incomingHasStableIdentity && existingByAddress != nil && len(existingByAddress.addresses) > 1 && existingByAddress.identityKey != "" {
		identityKey = existingByAddress.identityKey
	}
	if identityKey == "" && existingByAddress != nil {
		identityKey = existingByAddress.identityKey
	}

	existingByIdentity := (*deviceEntry)(nil)
	if identityKey != "" {
		existingByIdentity = r.identity[identityKey]
	}

	entry := existingByIdentity
	if entry == nil {
		entry = existingByAddress
	}
	if entry == existingByAddress && existingByIdentity == nil && existingByAddress != nil && (!canMergeIdentity(info, existingByAddress.info) || hasConflictingModelSignature(info, existingByAddress.info)) {
		entry = nil
	}
	if entry == nil {
		if fallback, ok := r.lookupCompatibleBySignatureLocked(info); ok {
			entry = fallback
		}
	}

	if existingByAddress != nil && existingByAddress != entry {
		r.detachAddressLocked(existingByAddress, info.Address)
	}

	if entry == nil {
		entry = &deviceEntry{
			primaryAddress: info.Address,
			addresses:      []byte{info.Address},
		}
		r.order = append(r.order, entry)
	} else if !containsAddress(entry.addresses, info.Address) {
		entry.addresses = append(entry.addresses, info.Address)
	}

	storedInfo := info
	if storedInfo.Manufacturer == "" {
		storedInfo.Manufacturer = entry.info.Manufacturer
	}
	preserveExistingDeviceID := existingByIdentity != nil && existingByIdentity != existingByAddress
	if preserveExistingDeviceID && entry.info.DeviceID != "" && storedInfo.DeviceID != entry.info.DeviceID {
		storedInfo.DeviceID = entry.info.DeviceID
	} else if storedInfo.DeviceID == "" {
		storedInfo.DeviceID = entry.info.DeviceID
	}
	if storedInfo.SoftwareVersion == "" {
		storedInfo.SoftwareVersion = entry.info.SoftwareVersion
	}
	if storedInfo.HardwareVersion == "" {
		storedInfo.HardwareVersion = entry.info.HardwareVersion
	}
	if storedInfo.SerialNumber == "" {
		storedInfo.SerialNumber = entry.info.SerialNumber
	}
	if storedInfo.MacAddress == "" {
		storedInfo.MacAddress = entry.info.MacAddress
	}
	storedInfo.Address = entry.primaryAddress
	physical = canonicalPhysicalIdentity(storedInfo)
	identityKey = physical.key()
	if identityKey == "" && entry.identityKey != "" {
		identityKey = entry.identityKey
	}

	for _, provider := range r.providers {
		if provider.Match(storedInfo) {
			matched = append(matched, provider)
			planes = append(planes, provider.CreatePlanes(storedInfo)...)
		}
	}

	projections := make([]Projection, 0)
	for _, provider := range matched {
		projectionProvider, ok := provider.(ProjectionProvider)
		if !ok {
			continue
		}
		projections = append(projections, projectionProvider.CreateProjections(storedInfo, planes)...)
	}

	index, projectionErr := BuildCanonicalIndex(projections)
	if projectionErr != nil {
		projections = nil
	}

	if entry.identityKey != "" && entry.identityKey != identityKey {
		delete(r.identity, entry.identityKey)
	}
	entry.info = storedInfo
	entry.physical = physical
	entry.identityKey = identityKey
	entry.planes = planes
	entry.projections = projections
	entry.index = index
	entry.indexErr = projectionErr

	if identityKey != "" {
		r.identity[identityKey] = entry
	}
	r.entries[info.Address] = entry
	r.observeAddressSlotLocked(info.Address, entry, DiscoverySourceActiveConfirmed, VerificationStateIdentityConfirmed)
	r.syncEntryFacesLocked(entry)

	return entry
}

func (r *DeviceRegistry) lookupByIdentity(info DeviceInfo) (DeviceEntry, bool) {
	identity := canonicalPhysicalIdentity(info).key()
	if identity == "" {
		return r.lookupBySignature(info)
	}

	r.mu.RLock()
	entry, ok := r.identity[identity]
	r.mu.RUnlock()
	if !ok {
		return r.lookupBySignature(info)
	}
	return entry, true
}

func (r *DeviceRegistry) lookupBySignature(info DeviceInfo) (DeviceEntry, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	entry, ok := r.lookupCompatibleBySignatureLocked(info)
	if !ok {
		return nil, false
	}
	return entry, true
}

func (r *DeviceRegistry) lookupCompatibleBySignatureLocked(info DeviceInfo) (*deviceEntry, bool) {
	signature := canonicalPhysicalIdentity(info).withFallbackModelSignature()
	if signature == "" {
		return nil, false
	}
	var match *deviceEntry
	for _, candidate := range r.order {
		if candidate == nil {
			continue
		}
		if candidate.physical.withFallbackModelSignature() != signature {
			continue
		}
		if !canMergeIdentity(info, candidate.info) {
			continue
		}
		if match != nil && match != candidate {
			return nil, false
		}
		match = candidate
	}
	if match == nil {
		return nil, false
	}
	return match, true
}

func canMergeIdentity(incoming DeviceInfo, existing DeviceInfo) bool {
	normalizedIncomingSerial := normalizeIdentityPart(incoming.SerialNumber)
	normalizedExistingSerial := normalizeIdentityPart(existing.SerialNumber)
	if normalizedIncomingSerial != "" && normalizedExistingSerial != "" && normalizedIncomingSerial != normalizedExistingSerial {
		return false
	}
	normalizedIncomingMAC := normalizeIdentityPart(incoming.MacAddress)
	normalizedExistingMAC := normalizeIdentityPart(existing.MacAddress)
	if normalizedIncomingMAC != "" && normalizedExistingMAC != "" && normalizedIncomingMAC != normalizedExistingMAC {
		return false
	}
	return true
}

func hasConflictingModelSignature(incoming DeviceInfo, existing DeviceInfo) bool {
	incomingDeviceID := normalizeIdentityPart(incoming.DeviceID)
	existingDeviceID := normalizeIdentityPart(existing.DeviceID)
	incomingSoftware := normalizeIdentityPart(incoming.SoftwareVersion)
	existingSoftware := normalizeIdentityPart(existing.SoftwareVersion)
	incomingHardware := normalizeIdentityPart(incoming.HardwareVersion)
	existingHardware := normalizeIdentityPart(existing.HardwareVersion)

	if incomingDeviceID == "" || existingDeviceID == "" || incomingSoftware == "" || existingSoftware == "" || incomingHardware == "" || existingHardware == "" {
		return false
	}

	return incomingDeviceID != existingDeviceID || incomingSoftware != existingSoftware || incomingHardware != existingHardware
}

// Lookup returns the canonical *DeviceEntry for the given address, or
// (nil, false) if no device occupies that slot. Preserved (signature
// unchanged) so existing callers continue to compile.
func (r *DeviceRegistry) Lookup(address byte) (DeviceEntry, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	entry, ok := r.entries[address]
	if !ok {
		return nil, false
	}
	return entry, true
}

// LookupSlot returns the AddressSlot for the requested address (M1
// address-table accessor), with the slot's own role/source/confidence
// metadata. When the address is aliased to a multi-address device, the
// returned slot.Device pointer is shared with the primary slot, but
// slot.Addr/Role/DiscoverySource/VerificationState describe the
// REQUESTED address — callers inspecting per-address metadata get the
// per-slot view (Codex P2: return-the-requested-address-slot).
func (r *DeviceRegistry) LookupSlot(address byte) (*AddressSlot, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	slot := r.addressTable[address]
	if slot == nil {
		return nil, false
	}
	return slot, true
}

func (r *DeviceRegistry) AliasAddresses(a, b byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	slotA := r.ensureAddressSlotLocked(a)
	slotB := r.ensureAddressSlotLocked(b)

	switch {
	case slotA.Device != nil:
		canonical := slotA.Device
		if secondary := slotB.Device; secondary != nil && secondary != canonical {
			secondary.addresses = removeAddress(secondary.addresses, b)
			if len(secondary.addresses) == 0 {
				// Secondary is fully removed: this is the BASV2 /
				// NETX3 scenario where the identity-bearing entry
				// becomes orphaned. Promote secondary's identity
				// fields onto canonical (only when canonical's
				// fields are empty) and transfer the r.identity
				// binding. (Codex P2 round-2 finding 2026-05-08 on
				// PR #136: absorb must NOT fire when secondary
				// survives, otherwise canonical and the surviving
				// secondary expose duplicate identity.)
				r.absorbIdentityLocked(canonical, secondary)
				if secondary.identityKey != "" {
					if canonical.identityKey == "" {
						canonical.identityKey = secondary.identityKey
						r.identity[canonical.identityKey] = canonical
					} else {
						delete(r.identity, secondary.identityKey)
					}
				}
				r.order = removeEntry(r.order, secondary)
			} else {
				// Secondary survives at remaining addresses (e.g.
				// 0x15 + 0x16 with same identity, then alias 0x10
				// to 0x15 → secondary still owns 0x16 with the
				// original identity). DO NOT absorb identity onto
				// canonical and DO NOT transfer identityKey:
				// secondary keeps its r.identity row + manufacturer
				// fields, canonical gets identity via Register's
				// identity-merge path on a future write.
				if secondary.primaryAddress == b {
					secondary.primaryAddress = secondary.addresses[0]
					secondary.info.Address = secondary.primaryAddress
				}
				r.syncEntryFacesLocked(secondary)
			}
		}
		slotB.Device = canonical
		if !containsAddress(canonical.addresses, b) {
			canonical.addresses = append(canonical.addresses, b)
		}
		r.entries[b] = canonical
		r.syncEntryFacesLocked(canonical)
	case slotB.Device != nil:
		canonical := slotB.Device
		if secondary := slotA.Device; secondary != nil && secondary != canonical {
			// Symmetric case. See slotA branch.
			secondary.addresses = removeAddress(secondary.addresses, a)
			if len(secondary.addresses) == 0 {
				r.absorbIdentityLocked(canonical, secondary)
				if secondary.identityKey != "" {
					if canonical.identityKey == "" {
						canonical.identityKey = secondary.identityKey
						r.identity[canonical.identityKey] = canonical
					} else {
						delete(r.identity, secondary.identityKey)
					}
				}
				r.order = removeEntry(r.order, secondary)
			} else {
				if secondary.primaryAddress == a {
					secondary.primaryAddress = secondary.addresses[0]
					secondary.info.Address = secondary.primaryAddress
				}
				r.syncEntryFacesLocked(secondary)
			}
		}
		slotA.Device = canonical
		if !containsAddress(canonical.addresses, a) {
			canonical.addresses = append(canonical.addresses, a)
		}
		r.entries[a] = canonical
		r.syncEntryFacesLocked(canonical)
	}

	return nil
}

// absorbIdentityLocked copies non-empty identity-bearing fields and
// derived state from src onto dst when dst's corresponding fields are
// empty. Re-keys r.identity[dst.identityKey] = dst when a new
// identityKey is adopted from src. Holds r.mu (caller's responsibility).
//
// Phase post-C P0: introduced to fix AliasAddresses identity loss
// (see live observation on 2026-05-08: BASV2 0x10↔0x15 + NETX3
// 0xF1↔0xF6 pairs aliased correctly but with manufacturer="" because
// the identity-bearing target-face entry was the secondary in the
// merge).
//
// Fields absorbed (only when dst's value is empty/zero AND src's value
// is non-empty):
//   - info.Manufacturer, info.DeviceID, info.SerialNumber, info.MacAddress
//   - info.SoftwareVersion, info.HardwareVersion
//   - physical (only if dst's physicalIdentity is zero)
//   - identityKey (re-keyed in r.identity)
//   - planes, projections, index, indexErr (only when dst has none)
//
// This function does NOT touch addresses or primaryAddress — those are
// owned by the caller (AliasAddresses) since the merge's address-graph
// semantics are independent of identity.
func (r *DeviceRegistry) absorbIdentityLocked(dst, src *deviceEntry) {
	if dst == nil || src == nil || dst == src {
		return
	}
	if dst.info.Manufacturer == "" && src.info.Manufacturer != "" {
		dst.info.Manufacturer = src.info.Manufacturer
	}
	if dst.info.DeviceID == "" && src.info.DeviceID != "" {
		dst.info.DeviceID = src.info.DeviceID
	}
	if dst.info.SerialNumber == "" && src.info.SerialNumber != "" {
		dst.info.SerialNumber = src.info.SerialNumber
	}
	if dst.info.MacAddress == "" && src.info.MacAddress != "" {
		dst.info.MacAddress = src.info.MacAddress
	}
	if dst.info.SoftwareVersion == "" && src.info.SoftwareVersion != "" {
		dst.info.SoftwareVersion = src.info.SoftwareVersion
	}
	if dst.info.HardwareVersion == "" && src.info.HardwareVersion != "" {
		dst.info.HardwareVersion = src.info.HardwareVersion
	}
	// Adopt physicalIdentity + identityKey if dst has neither and src has both.
	emptyPhysical := physicalIdentity{}
	if dst.physical == emptyPhysical && src.physical != emptyPhysical {
		dst.physical = src.physical
	}
	// NOTE: identityKey transfer is intentionally NOT done here.
	// absorbIdentityLocked runs BEFORE AliasAddresses removes
	// `b`/`a` from src.addresses, so we don't yet know whether src
	// will survive (multiple addresses → survives at remaining face)
	// or be fully removed (only the aliased address → removed). The
	// caller (AliasAddresses) handles identityKey transfer in the
	// "secondary fully removed" branch only — see the comment at
	// the call site. This avoids the live-validation P2 finding
	// where surviving multi-address secondaries lost their identity
	// map binding because absorb prematurely moved identityKey.
	//
	// Adopt planes / projections / index if dst has none. These
	// derive from info; absorbing them avoids re-running providers.
	if len(dst.planes) == 0 && len(src.planes) > 0 {
		dst.planes = src.planes
	}
	if len(dst.projections) == 0 && len(src.projections) > 0 {
		dst.projections = src.projections
	}
	if dst.index.canonicalByID == nil && src.index.canonicalByID != nil {
		dst.index = src.index
	}
	if dst.indexErr == nil && src.indexErr != nil {
		dst.indexErr = src.indexErr
	}
}

func (r *DeviceRegistry) Iterate(fn func(DeviceEntry) bool) {
	r.mu.RLock()
	order := make([]*deviceEntry, len(r.order))
	copy(order, r.order)
	r.mu.RUnlock()

	for _, entry := range order {
		if entry == nil {
			continue
		}
		if !fn(entry) {
			return
		}
	}
}

func (r *DeviceRegistry) detachAddressLocked(entry *deviceEntry, address byte) {
	if entry == nil {
		return
	}
	delete(r.entries, address)
	if slot := r.addressTable[address]; slot != nil && slot.Device == entry {
		r.addressTable[address] = nil
	}
	if !containsAddress(entry.addresses, address) {
		return
	}

	entry.addresses = removeAddress(entry.addresses, address)
	if len(entry.addresses) == 0 {
		if entry.identityKey != "" {
			delete(r.identity, entry.identityKey)
		}
		r.order = removeEntry(r.order, entry)
		return
	}
	if entry.primaryAddress == address {
		entry.primaryAddress = entry.addresses[0]
		entry.info.Address = entry.primaryAddress
	}
	r.syncEntryFacesLocked(entry)
}

func (r *DeviceRegistry) ensureAddressSlotLocked(address byte) *AddressSlot {
	slot := r.addressTable[address]
	if slot == nil {
		slot = &AddressSlot{Addr: address}
		r.addressTable[address] = slot
	}
	return slot
}

// MarkSlotPassiveObserved updates an AddressSlot for an address that was
// passively observed by the gateway (e.g. by AddressTableInserter on
// positive ACK following a complete request). Writes Role / Discovery
// Source / VerificationState / FirstObservedAt / LastObservedAt under the
// registry write lock so concurrent readers via LookupSlot / Lookup do
// not see torn state.
//
// This API replaces direct *AddressSlot field mutation by the gateway
// inserter, which was racy with other readers (Codex P2 follow-up from
// PR #565). Idempotent: re-marking the same slot only advances
// DiscoverySource / VerificationState monotonically (the slot retains
// the higher of the existing and new value, matching
// observeAddressSlotLocked's monotonic semantics).
func (r *DeviceRegistry) MarkSlotPassiveObserved(address byte, role SlotRole, observedAt time.Time) {
	if r == nil {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	slot := r.ensureAddressSlotLocked(address)
	if slot.DiscoverySource < DiscoverySourcePassiveObserved {
		slot.DiscoverySource = DiscoverySourcePassiveObserved
	}
	if slot.VerificationState < VerificationStateCorroborated {
		slot.VerificationState = VerificationStateCorroborated
	}
	if role != SlotRoleUnknown && slot.Role == SlotRoleUnknown {
		slot.Role = role
	}
	if slot.FirstObservedAt.IsZero() && !observedAt.IsZero() {
		slot.FirstObservedAt = observedAt
	}
	if !observedAt.IsZero() {
		slot.LastObservedAt = observedAt
	}
	// Phase C M-C6a: refresh entry.Faces so AddressByRole sees the
	// updated SlotRole. Without this sync, MarkSlotPassiveObserved
	// would leave Faces stale and AddressByRole(SlotRoleSlave) on a
	// just-passively-observed slot would return (0, false).
	if slot.Device != nil {
		r.syncEntryFacesLocked(slot.Device)
	}
}

func (r *DeviceRegistry) observeAddressSlotLocked(address byte, entry *deviceEntry, source DiscoverySource, state VerificationState) {
	now := time.Now()
	slot := r.ensureAddressSlotLocked(address)
	slot.Device = entry
	if slot.DiscoverySource < source {
		slot.DiscoverySource = source
	}
	if slot.VerificationState < state {
		slot.VerificationState = state
	}
	if slot.FirstObservedAt.IsZero() {
		slot.FirstObservedAt = now
	}
	slot.LastObservedAt = now
}

func (r *DeviceRegistry) syncEntryFacesLocked(entry *deviceEntry) {
	if entry == nil {
		return
	}
	faces := make([]BusFace, 0, len(entry.addresses))
	for _, address := range entry.addresses {
		slot := r.ensureAddressSlotLocked(address)
		if slot.Device == nil {
			slot.Device = entry
		}
		faces = append(faces, BusFace{
			Addr:              address,
			Role:              slot.Role,
			DiscoverySource:   slot.DiscoverySource,
			VerificationState: slot.VerificationState,
		})
	}
	entry.Faces = faces
}

type deviceEntry struct {
	primaryAddress byte
	addresses      []byte
	physical       physicalIdentity
	identityKey    string
	info           DeviceInfo
	planes         []Plane
	projections    []Projection
	index          CanonicalIndex
	indexErr       error
	Faces          []BusFace
}

// PrimaryDisplayAddress returns a representative address for log/UI
// display. Returns the canonical primary if set, otherwise the
// originally registered info.Address. Use this for log lines,
// MCP/GraphQL device.address fields, UI labels — anywhere the value
// is shown to humans rather than written to the wire. For wire
// routing, use AddressByRole(SlotRole) which is class-aware.
//
// Phase C M-C6c: replaces deviceEntry.Address(), whose name conflated
// display and routing semantics for aliased canonical pairs.
func (d *deviceEntry) PrimaryDisplayAddress() byte {
	if d.primaryAddress != 0 {
		return d.primaryAddress
	}
	return d.info.Address
}

// AddressByRole returns the first BusFace address whose Role matches.
// Routing-correct alternative to Address: M2S writers pass
// SlotRoleSlave to get the target byte; M2M arbitration logic passes
// SlotRoleMaster for the initiator byte. Returns (0, false) when no
// face matches the requested role.
//
// Decision references: Phase C AD30 (entry.Address ambiguity fix);
// uses the existing BusFace.Role machinery populated by
// syncEntryFacesLocked.
func (d *deviceEntry) AddressByRole(role SlotRole) (byte, bool) {
	if d == nil {
		return 0, false
	}
	// Pass 1: explicit Role match (set via MarkSlotPassiveObserved or
	// other role-aware paths).
	for _, face := range d.Faces {
		if face.Role == role {
			return face.Addr, true
		}
	}
	// Pass 2: SlotRoleUnknown fallback — active scan registers entries
	// without populating Role (Codex P2 from PR #134). Infer the role
	// from the address class so callers migrating from Address() to
	// AddressByRole get a useful answer for actively-scanned devices.
	for _, face := range d.Faces {
		if face.Role != SlotRoleUnknown {
			continue
		}
		switch protocol.AddressClassOf(face.Addr) {
		case protocol.AddressClassMaster:
			if role == SlotRoleMaster {
				return face.Addr, true
			}
		case protocol.AddressClassSlave:
			if role == SlotRoleSlave {
				return face.Addr, true
			}
		}
	}
	return 0, false
}

func (d *deviceEntry) Addresses() []byte {
	if len(d.addresses) == 0 {
		if d.info.Address == 0 {
			return nil
		}
		return []byte{d.info.Address}
	}
	out := make([]byte, len(d.addresses))
	copy(out, d.addresses)
	return out
}

func (d *deviceEntry) Manufacturer() string {
	return d.info.Manufacturer
}

func (d *deviceEntry) DeviceID() string {
	return d.info.DeviceID
}

func (d *deviceEntry) SerialNumber() string {
	return d.info.SerialNumber
}

func (d *deviceEntry) MacAddress() string {
	return d.info.MacAddress
}

func (d *deviceEntry) SoftwareVersion() string {
	return d.info.SoftwareVersion
}

func (d *deviceEntry) HardwareVersion() string {
	return d.info.HardwareVersion
}

func (d *deviceEntry) Planes() []Plane {
	return d.planes
}

func (d *deviceEntry) Projections() []Projection {
	return d.projections
}

func CanonicalIndexForEntry(entry DeviceEntry) (CanonicalIndex, error) {
	if entry == nil {
		return CanonicalIndex{}, ErrProjectionInvalidNode
	}
	if internal, ok := entry.(*deviceEntry); ok {
		return internal.index, internal.indexErr
	}
	return BuildCanonicalIndex(entry.Projections())
}

func containsAddress(addresses []byte, address byte) bool {
	for _, existing := range addresses {
		if existing == address {
			return true
		}
	}
	return false
}

func removeAddress(addresses []byte, address byte) []byte {
	for index, existing := range addresses {
		if existing != address {
			continue
		}
		return append(addresses[:index], addresses[index+1:]...)
	}
	return addresses
}

func removeEntry(entries []*deviceEntry, entry *deviceEntry) []*deviceEntry {
	for index, existing := range entries {
		if existing != entry {
			continue
		}
		return append(entries[:index], entries[index+1:]...)
	}
	return entries
}
