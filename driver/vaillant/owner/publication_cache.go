package owner

import (
	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/output"
)

func preloadSemanticCache(state *retainedSemanticState, cacheStore *output.SemanticCacheStore) (VaillantCacheHydration, bool) {
	if state == nil || cacheStore == nil {
		return VaillantCacheHydration{}, false
	}
	snapshot, ok := cacheStore.Load()
	if !ok {
		return VaillantCacheHydration{}, false
	}
	if len(snapshot.Zones) > 0 {
		state.SetZonesFromCache(snapshot.Zones)
	}
	if snapshot.DHW != nil {
		state.SetDHWFromCache(snapshot.DHW)
	}
	if snapshot.Boiler != nil {
		state.SetBoilerStatusFromCache(snapshot.Boiler)
	}
	return cacheHydrationFromSnapshot(snapshot), true
}

func cacheHydrationFromSnapshot(snapshot output.SemanticCacheSnapshot) VaillantCacheHydration {
	native := output.CacheHydrationFromSnapshot(snapshot)
	hydration := VaillantCacheHydration{persistedAt: native.PersistedAt}
	if native.Zones != nil {
		hydration.zones = make([]*vaillantZoneSnapshot, 0, len(native.Zones))
		for _, zone := range native.Zones {
			hydration.zones = append(hydration.zones, ownerZoneSnapshot(zone))
		}
	}
	hydration.dhw = ownerDhwSnapshot(native.Dhw)
	if native.Boiler != nil {
		hydration.boiler = (*vaillantBoilerSnapshot)(native.Boiler)
	}
	return hydration
}

func (o *vaillantSemanticPublication) persistsCache() bool {
	return o.cache != nil
}

func (o *vaillantSemanticPublication) persistCache(overlays map[byte]vaillantZoneCacheOverlay) {
	native := make(map[byte]output.NativeZoneCacheOverlay, len(overlays))
	for instance, overlay := range overlays {
		native[instance] = output.NativeZoneCacheOverlay{
			CurrentTempC:         overlay.currentTempC,
			HumidityPct:          overlay.humidityPct,
			RoomMappingRaw:       overlay.roomMappingRaw,
			AssociatedCircuitRaw: overlay.associatedCircuitRaw,
		}
	}
	_ = o.cache.Save(output.SemanticCacheSnapshot{
		Zones:  output.CacheZones(o.state.Zones(), native),
		DHW:    o.state.DHW(),
		Boiler: o.state.BoilerStatus(),
	})
}
