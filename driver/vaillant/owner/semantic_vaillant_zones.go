package owner

import (
	"context"
)

func (p *vaillantSemanticPoller) refreshZoneDiscovery(ctx context.Context, startup bool) {
	present := make(map[byte]bool, 4)
	checked := make(map[byte]bool, 9)
	for instance := byte(0x00); instance <= 0x08; instance++ { // II_MAX=0x08 per Vaillant regulator spec
		var indexBytes []byte
		var ok bool
		if startup {
			indexBytes, ok = p.readB524Startup(ctx, localZones.opcode, localZones.group, instance, zoneIndex)
		} else {
			indexBytes, ok = p.readB524Value(ctx, localZones.opcode, localZones.group, instance, zoneIndex)
		}
		if !ok {
			continue
		}
		checked[instance] = true
		if len(indexBytes) < 1 || indexBytes[0] == 0xFF {
			continue
		}
		present[instance] = true
	}

	source := p.reconcileDiscoveryPresence(ctx, checked, present)
	p.publishZones(source)
}

func (p *vaillantSemanticPoller) reconcileDiscoveryPresence(ctx context.Context, checked, present map[byte]bool) semanticSnapshotSource {
	if len(present) == 0 {
		grabPresent, ok := p.tryRefreshFromEbusdGrabWithInstances(ctx)
		if ok {
			if checked == nil {
				checked = make(map[byte]bool, len(grabPresent))
			}
			if present == nil {
				present = make(map[byte]bool, len(grabPresent))
			}
			for instance := range grabPresent {
				checked[instance] = true
				present[instance] = true
			}
			p.applyZonePresenceMissesOnly(checked, present)
			return semanticSnapshotSourceLive
		}
	}
	p.applyZonePresenceProbes(checked, present)
	if len(present) > 0 && p.startupZonesNeedImmediatePresence() {
		p.applyZonePresenceImmediate(checked, present)
		return semanticSnapshotSourceLive
	}
	if len(present) == 0 {
		return semanticSnapshotSourceCache
	}
	return semanticSnapshotSourceLive
}

func (p *vaillantSemanticPoller) startupZonesNeedImmediatePresence() bool {
	if !p.publishing() {
		return false
	}
	if p.publication.publishedZoneCount() > 0 {
		return false
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.zones) == 0
}

func (p *vaillantSemanticPoller) applyZonePresenceImmediate(checked, present map[byte]bool) {
	if p == nil {
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if p.zones == nil {
		p.zones = make(map[byte]*vaillantZoneSnapshot)
	}
	if p.presence == nil {
		p.presence = make(map[byte]*zonePresenceRecord)
	}

	hitThreshold := p.zoneHitThresholdValue()
	for instance := range checked {
		if !present[instance] {
			p.markZoneMissingLocked(instance)
			continue
		}
		record := p.ensureZonePresenceRecordLocked(instance)
		previousState := record.State
		record.State = zonePresencePresent
		record.HitStreak = hitThreshold
		record.MissStreak = 0
		p.ensureZoneEntryLocked(instance)
		p.recordZonePresenceTransitionLocked(instance, previousState, record.State, record.HitStreak, record.MissStreak)
	}
	semanticZoneCount.Set(int64(len(p.zones)))
}

func (p *vaillantSemanticPoller) applyZonePresenceProbes(checked, present map[byte]bool) {
	if p == nil {
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if p.zones == nil {
		p.zones = make(map[byte]*vaillantZoneSnapshot)
	}
	if p.presence == nil {
		p.presence = make(map[byte]*zonePresenceRecord)
	}

	for instance := range checked {
		if present[instance] {
			p.markZonePresentLocked(instance)
			continue
		}
		p.markZoneMissingLocked(instance)
	}
}

func (p *vaillantSemanticPoller) applyZonePresenceMissesOnly(checked, present map[byte]bool) {
	if p == nil {
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if p.zones == nil {
		p.zones = make(map[byte]*vaillantZoneSnapshot)
	}
	if p.presence == nil {
		p.presence = make(map[byte]*zonePresenceRecord)
	}

	for instance := range checked {
		if present[instance] {
			continue
		}
		p.markZoneMissingLocked(instance)
	}
}

func (p *vaillantSemanticPoller) markZonePresentLocked(instance byte) {
	record := p.ensureZonePresenceRecordLocked(instance)
	previousState := record.State
	record.MissStreak = 0
	hitThreshold := p.zoneHitThresholdValue()

	switch record.State {
	case zonePresencePresent, zonePresenceSuspectMissing:
		record.State = zonePresencePresent
		if record.HitStreak < hitThreshold {
			record.HitStreak = hitThreshold
		}
		p.ensureZoneEntryLocked(instance)
	default:
		record.HitStreak++
		if record.HitStreak >= hitThreshold {
			record.State = zonePresencePresent
			record.HitStreak = hitThreshold
			p.ensureZoneEntryLocked(instance)
		} else {
			record.State = zonePresenceSuspectResurrect
		}
	}
	p.recordZonePresenceTransitionLocked(instance, previousState, record.State, record.HitStreak, record.MissStreak)
	semanticZoneCount.Set(int64(len(p.zones)))
}

func (p *vaillantSemanticPoller) markZoneMissingLocked(instance byte) {
	record := p.ensureZonePresenceRecordLocked(instance)
	previousState := record.State
	record.HitStreak = 0
	missThreshold := p.zoneMissThresholdValue()

	switch record.State {
	case zonePresencePresent, zonePresenceSuspectMissing:
		record.MissStreak++
		if record.MissStreak >= missThreshold {
			record.MissStreak = missThreshold
			record.State = zonePresenceAbsent
			delete(p.zones, instance)
			p.recordZonePresenceTransitionLocked(instance, previousState, record.State, record.HitStreak, record.MissStreak)
			semanticZoneCount.Set(int64(len(p.zones)))
			return
		}
		record.State = zonePresenceSuspectMissing
	default:
		record.MissStreak = 0
		record.State = zonePresenceAbsent
		delete(p.zones, instance)
	}
	p.recordZonePresenceTransitionLocked(instance, previousState, record.State, record.HitStreak, record.MissStreak)
	semanticZoneCount.Set(int64(len(p.zones)))
}

func (p *vaillantSemanticPoller) ensureZonePresenceRecordLocked(instance byte) *zonePresenceRecord {
	if p.presence == nil {
		p.presence = make(map[byte]*zonePresenceRecord)
	}
	record := p.presence[instance]
	if record == nil {
		record = &zonePresenceRecord{State: zonePresenceAbsent}
		p.presence[instance] = record
	}
	return record
}

func (p *vaillantSemanticPoller) ensureZoneEntryLocked(instance byte) {
	entry := p.zones[instance]
	if entry == nil {
		entry = &vaillantZoneSnapshot{Instance: instance}
		p.zones[instance] = entry
	}
	entry.Present = true
}

func (p *vaillantSemanticPoller) zoneMissThresholdValue() int {
	if p == nil || p.zoneMissThreshold <= 0 {
		return DefaultSemanticZonePresenceMissThreshold
	}
	return p.zoneMissThreshold
}

func (p *vaillantSemanticPoller) zoneHitThresholdValue() int {
	if p == nil || p.zoneHitThreshold <= 0 {
		return DefaultSemanticZonePresenceHitThreshold
	}
	return p.zoneHitThreshold
}

func (p *vaillantSemanticPoller) refreshConfig(ctx context.Context) {
	controller, zones := p.snapshotZones()
	grabHydrated := false
	if controller == 0 || len(zones) == 0 {
		if p.fallbackDiscoveryAllowed(controller) {
			p.refreshDiscovery(ctx)
			controller, zones = p.snapshotZones()
		} else {
			p.refreshBoilerAddress(ctx)
		}
	}
	if controller == 0 || len(zones) == 0 {
		if p.tryRefreshFromEbusdGrab(ctx) {
			controller, zones = p.snapshotZones()
			grabHydrated = true
		}
	}
	if controller == 0 || len(zones) == 0 {
		if controller != 0 {
			p.publishDHW(p.refreshDHWSlowConfig(ctx))
		}
		return
	}

	zoneLiveReadSuccess := false
	for _, instance := range zones {
		primaryName, primaryOK := p.readB524ZoneNamePart(ctx, instance, zoneName)
		prefix, prefixOK := p.readB524ZoneNamePart(ctx, instance, zoneNamePrefix)
		suffix, suffixOK := p.readB524ZoneNamePart(ctx, instance, zoneNameSuffix)
		if primaryOK || prefixOK || suffixOK {
			zoneLiveReadSuccess = true
		}

		incoming, slowOK := p.refreshZoneSlowConfigFields(ctx, instance)
		if slowOK {
			zoneLiveReadSuccess = true
		}
		if incoming == nil {
			incoming = &vaillantZoneSnapshot{}
		}
		incoming.Name = composeZoneName(primaryName, prefix, suffix)

		p.mu.Lock()
		if entry := p.zones[instance]; entry != nil {
			mergeZoneSnapshotFields(entry, incoming, semanticSnapshotSourceLive, zoneConfigRefreshFieldSet)
		}
		p.mu.Unlock()
	}

	dhwSource := p.refreshDHWSlowConfig(ctx)
	if !zoneLiveReadSuccess && p.tryRefreshFromEbusdGrab(ctx) {
		grabHydrated = true
	}

	zoneSource := semanticSnapshotSourceCache
	if zoneLiveReadSuccess || grabHydrated {
		zoneSource = semanticSnapshotSourceLive
	}
	p.publishZones(zoneSource)
	p.publishDHW(dhwSource)
}
