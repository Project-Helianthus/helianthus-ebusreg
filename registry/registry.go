package registry

import (
	"sync"

	"github.com/d3vi1/helianthus-ebusreg/schema"
)

type DeviceInfo struct {
	Address         byte
	Manufacturer    string
	DeviceID        string
	SoftwareVersion string
	HardwareVersion string
}

type DeviceEntry interface {
	Address() byte
	Manufacturer() string
	DeviceID() string
	SoftwareVersion() string
	HardwareVersion() string
	Planes() []Plane
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
	order     []byte
}

func NewDeviceRegistry(providers []PlaneProvider) *DeviceRegistry {
	providerCopy := make([]PlaneProvider, len(providers))
	copy(providerCopy, providers)
	return &DeviceRegistry{
		providers: providerCopy,
		entries:   make(map[byte]*deviceEntry),
	}
}

func (r *DeviceRegistry) RegisterProvider(provider PlaneProvider) {
	r.mu.Lock()
	r.providers = append(r.providers, provider)
	r.mu.Unlock()
}

func (r *DeviceRegistry) Register(info DeviceInfo) DeviceEntry {
	r.mu.RLock()
	providers := make([]PlaneProvider, len(r.providers))
	copy(providers, r.providers)
	r.mu.RUnlock()

	planes := make([]Plane, 0)
	for _, provider := range providers {
		if provider.Match(info) {
			planes = append(planes, provider.CreatePlanes(info)...)
		}
	}

	entry := &deviceEntry{
		info:   info,
		planes: planes,
	}

	r.mu.Lock()
	if _, exists := r.entries[info.Address]; !exists {
		r.order = append(r.order, info.Address)
	}
	r.entries[info.Address] = entry
	r.mu.Unlock()

	return entry
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
	order := make([]byte, len(r.order))
	copy(order, r.order)
	entries := make(map[byte]*deviceEntry, len(r.entries))
	for address, entry := range r.entries {
		entries[address] = entry
	}
	r.mu.RUnlock()

	for _, address := range order {
		entry, ok := entries[address]
		if !ok {
			continue
		}
		if !fn(entry) {
			return
		}
	}
}

type deviceEntry struct {
	info   DeviceInfo
	planes []Plane
}

func (d *deviceEntry) Address() byte {
	return d.info.Address
}

func (d *deviceEntry) Manufacturer() string {
	return d.info.Manufacturer
}

func (d *deviceEntry) DeviceID() string {
	return d.info.DeviceID
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
