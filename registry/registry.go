package registry

import (
	"sync"

	"github.com/d3vi1/helianthus-ebusreg/schema"
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
	Address() byte
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
	mu        sync.RWMutex
	providers []PlaneProvider
	entries   map[byte]*deviceEntry
	identity  map[string]*deviceEntry
	order     []*deviceEntry
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
	if identityKey == "" && existingByAddress != nil {
		physical = existingByAddress.physical
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
	if storedInfo.DeviceID == "" {
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
	if info.SerialNumber == "" && info.MacAddress == "" && entry.identityKey != "" {
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

func (r *DeviceRegistry) Lookup(address byte) (DeviceEntry, bool) {
	r.mu.RLock()
	entry, ok := r.entries[address]
	r.mu.RUnlock()
	if !ok {
		return nil, false
	}
	return entry, true
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
}

func (d *deviceEntry) Address() byte {
	if d.primaryAddress != 0 {
		return d.primaryAddress
	}
	return d.info.Address
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
