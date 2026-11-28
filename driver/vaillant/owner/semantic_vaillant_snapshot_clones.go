package owner

import (
	"maps"
	"strings"
	"time"
)

func cloneZoneSnapshot(snapshot *vaillantZoneSnapshot) *vaillantZoneSnapshot {
	if snapshot == nil {
		return nil
	}
	cloned := *snapshot
	cloned.AllowedModes = append([]string(nil), snapshot.AllowedModes...)
	cloned.CurrentTempC = cloneFloat64Ptr(snapshot.CurrentTempC)
	cloned.TargetTempC = cloneFloat64Ptr(snapshot.TargetTempC)
	cloned.HumidityPct = cloneFloat64Ptr(snapshot.HumidityPct)
	cloned.ConfigurationRoomTemperatureZoneMappingRaw = cloneUint16Ptr(snapshot.ConfigurationRoomTemperatureZoneMappingRaw)
	cloned.ConfigurationAssociatedCircuitRaw = cloneUint16Ptr(snapshot.ConfigurationAssociatedCircuitRaw)
	cloned.ConfigurationCircuitTypeRaw = cloneUint16Ptr(snapshot.ConfigurationCircuitTypeRaw)
	cloned.StateValveStatusRaw = cloneUint16Ptr(snapshot.StateValveStatusRaw)
	cloned.QuickVetoTempC = cloneFloat64Ptr(snapshot.QuickVetoTempC)
	cloned.QuickVetoDurationH = cloneFloat64Ptr(snapshot.QuickVetoDurationH)
	cloned.HolidaySetpointC = cloneFloat64Ptr(snapshot.HolidaySetpointC)
	cloned.FieldFreshness = maps.Clone(snapshot.FieldFreshness)
	return &cloned
}

func cloneDhwSnapshot(snapshot *vaillantDhwSnapshot) *vaillantDhwSnapshot {
	if snapshot == nil {
		return nil
	}
	cloned := *snapshot
	cloned.CurrentTempC = cloneFloat64Ptr(snapshot.CurrentTempC)
	cloned.TargetTempC = cloneFloat64Ptr(snapshot.TargetTempC)
	cloned.FieldFreshness = maps.Clone(snapshot.FieldFreshness)
	return &cloned
}

func cloneFloat64Ptr(value *float64) *float64 {
	if value == nil {
		return nil
	}
	copyValue := *value
	return &copyValue
}

func cloneUint16Ptr(value *uint16) *uint16 {
	if value == nil {
		return nil
	}
	copyValue := *value
	return &copyValue
}

func cloneUint8Ptr(value *uint8) *uint8 {
	if value == nil {
		return nil
	}
	copyValue := *value
	return &copyValue
}

func cloneUint32Ptr(value *uint32) *uint32 {
	if value == nil {
		return nil
	}
	copyValue := *value
	return &copyValue
}

func ensureFieldFreshnessMap(entry *map[semanticFieldKey]semanticFieldFreshness) map[semanticFieldKey]semanticFieldFreshness {
	if *entry == nil {
		*entry = make(map[semanticFieldKey]semanticFieldFreshness)
	}
	return *entry
}

func markFieldStale(fields map[semanticFieldKey]semanticFieldFreshness, field semanticFieldKey) {
	if fields == nil {
		return
	}
	freshness, ok := fields[field]
	if !ok {
		return
	}
	freshness.Stale = true
	fields[field] = freshness
}

// SeedZoneFreshness implements the seed zone freshness step of the Vaillant owner.
func SeedZoneFreshness(snapshot *vaillantZoneSnapshot, source semanticSnapshotSource, stale bool) {
	if snapshot == nil {
		return
	}
	fields := ensureFieldFreshnessMap(&snapshot.FieldFreshness)
	if strings.TrimSpace(snapshot.Name) != "" {
		setFieldFreshness(fields, zoneFieldName, source, stale)
	}
	if strings.TrimSpace(snapshot.OperatingMode) != "" {
		setFieldFreshness(fields, zoneFieldOperatingMode, source, stale)
	}
	if strings.TrimSpace(snapshot.Preset) != "" {
		setFieldFreshness(fields, zoneFieldPreset, source, stale)
	}
	if strings.TrimSpace(snapshot.HvacAction) != "" {
		setFieldFreshness(fields, zoneFieldHvacAction, source, stale)
	}
	if len(snapshot.AllowedModes) > 0 {
		setFieldFreshness(fields, zoneFieldAllowedModes, source, stale)
	}
	if snapshot.CurrentTempC != nil {
		setFieldFreshness(fields, zoneFieldCurrentTempC, source, stale)
	}
	if snapshot.TargetTempC != nil {
		setFieldFreshness(fields, zoneFieldTargetTempC, source, stale)
	}
	if snapshot.HumidityPct != nil {
		setFieldFreshness(fields, zoneFieldCurrentHumidityPct, source, stale)
	}
	if strings.TrimSpace(snapshot.ConfigurationHeatingOperationMode) != "" {
		setFieldFreshness(fields, zoneFieldZoneOperationModeRaw, source, stale)
	}
	if strings.TrimSpace(snapshot.StateSpecialFunction) != "" {
		setFieldFreshness(fields, zoneFieldSpecialFunctionRaw, source, stale)
	}
	if snapshot.ConfigurationRoomTemperatureZoneMappingRaw != nil {
		setFieldFreshness(fields, zoneFieldRoomTemperatureZoneMappingRaw, source, stale)
	}
	if snapshot.ConfigurationAssociatedCircuitRaw != nil {
		setFieldFreshness(fields, zoneFieldZoneCircuitIndexRaw, source, stale)
	}
	if snapshot.ConfigurationCircuitTypeRaw != nil {
		setFieldFreshness(fields, zoneFieldCircuitTypeRaw, source, stale)
	}
	if snapshot.StateValveStatusRaw != nil {
		setFieldFreshness(fields, zoneFieldZoneValveStatusRaw, source, stale)
	}
	if snapshot.QuickVetoTempC != nil {
		setFieldFreshness(fields, zoneFieldQuickVetoTempC, source, stale)
	}
	if snapshot.QuickVetoDurationH != nil {
		setFieldFreshness(fields, zoneFieldQuickVetoDurationH, source, stale)
	}
	if snapshot.QuickVetoEndTime != "" {
		setFieldFreshness(fields, zoneFieldQuickVetoEndTime, source, stale)
	}
	if snapshot.QuickVetoEndDate != "" {
		setFieldFreshness(fields, zoneFieldQuickVetoEndDate, source, stale)
	}
	if snapshot.HolidayStartDate != "" {
		setFieldFreshness(fields, zoneFieldHolidayStartDate, source, stale)
	}
	if snapshot.HolidayEndDate != "" {
		setFieldFreshness(fields, zoneFieldHolidayEndDate, source, stale)
	}
	if snapshot.HolidaySetpointC != nil {
		setFieldFreshness(fields, zoneFieldHolidaySetpointC, source, stale)
	}
	if snapshot.HolidayStartTime != "" {
		setFieldFreshness(fields, zoneFieldHolidayStartTime, source, stale)
	}
	if snapshot.HolidayEndTime != "" {
		setFieldFreshness(fields, zoneFieldHolidayEndTime, source, stale)
	}
}

// SeedDhwFreshness implements the seed dhw freshness step of the Vaillant owner.
func SeedDhwFreshness(snapshot *vaillantDhwSnapshot, source semanticSnapshotSource, stale bool) {
	if snapshot == nil {
		return
	}
	fields := ensureFieldFreshnessMap(&snapshot.FieldFreshness)
	if strings.TrimSpace(snapshot.OperatingMode) != "" {
		setFieldFreshness(fields, dhwFieldOperatingMode, source, stale)
	}
	if strings.TrimSpace(snapshot.Preset) != "" {
		setFieldFreshness(fields, dhwFieldPreset, source, stale)
	}
	if snapshot.CurrentTempC != nil {
		setFieldFreshness(fields, dhwFieldCurrentTempC, source, stale)
	}
	if snapshot.TargetTempC != nil {
		setFieldFreshness(fields, dhwFieldTargetTempC, source, stale)
	}
	if strings.TrimSpace(snapshot.ConfigurationDHWOperationMode) != "" {
		setFieldFreshness(fields, dhwFieldDhwOperationModeRaw, source, stale)
	}
	if strings.TrimSpace(snapshot.StateSpecialFunction) != "" {
		setFieldFreshness(fields, dhwFieldSpecialFunctionRaw, source, stale)
	}
	if snapshot.HolidayStartDate != "" {
		setFieldFreshness(fields, dhwFieldHolidayStartDate, source, stale)
	}
	if snapshot.HolidayEndDate != "" {
		setFieldFreshness(fields, dhwFieldHolidayEndDate, source, stale)
	}
}

// mergePublishedBoilerFlow copies a qualified thermal scalar into the native
// field and clears its cached mark. A withheld scalar leaves the native value
// and its age. System flow does not seed this field.
func (p *vaillantSemanticPoller) mergePublishedBoilerFlow() bool {
	p.mu.Lock()
	publisher, target := p.boilerThermal, p.boilerAddress
	p.mu.Unlock()
	var flow *float64
	if publisher != nil {
		flow = publisher.scalar(target)
	}
	if flow == nil {
		return false
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	var previous *float64
	if p.boiler != nil {
		previous = p.boiler.FlowTemperatureC
	}
	changed := !floatPtrEquals(previous, flow) || p.boilerFlowCached
	if p.boiler == nil {
		p.boiler = &vaillantBoilerSnapshot{}
	}
	p.boiler.FlowTemperatureC = flow
	p.boilerFlowCached = false
	p.boilerFlowObservedAt = p.now()
	return changed
}

// VaillantCacheHydration is the native value that a loaded semantic cache
// hydrates the poller with. zones holds one snapshot per cached zone in cache
// order, each keyed by its Instance; a later zone with the same instance
// replaces an earlier one. dhw and boiler are nil when the cache holds none.
// persistedAt is when the cache was written, zero when unknown. The poller
// takes ownership of the snapshots when it hydrates.
type VaillantCacheHydration struct {
	zones       []*vaillantZoneSnapshot
	dhw         *vaillantDhwSnapshot
	persistedAt time.Time
	boiler      *vaillantBoilerSnapshot
}

// vaillantZoneCacheOverlay holds the retained native zone values that the
// persisted cache takes over the published zone: the current temperature and
// humidity, and the raw room temperature mapping and associated circuit.
type vaillantZoneCacheOverlay struct {
	currentTempC         *float64
	humidityPct          *float64
	roomMappingRaw       *uint16
	associatedCircuitRaw *uint16
}

// vaillantFM5View is the native FM5 state that the poller hands to the
// publication owner: the acquisition verdict and detached copies of the
// retained solar and cylinder snapshots.
type vaillantFunctionModuleView struct {
	verdict   vaillantFunctionModuleInterpretation
	solar     *vaillantSolarSnapshot
	cylinders map[byte]*vaillantCylinderSnapshot
}
