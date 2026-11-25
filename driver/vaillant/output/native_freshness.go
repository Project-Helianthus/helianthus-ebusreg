package output

import (
	"strings"
)

// CloneFloat64Ptr copies the pointed-to value, preserving absence.
func CloneFloat64Ptr(value *float64) *float64 {
	if value == nil {
		return nil
	}
	copyValue := *value
	return &copyValue
}

// EnsureFieldFreshnessMap initializes the freshness map if absent.
func EnsureFieldFreshnessMap(entry *map[SemanticFieldKey]SemanticFieldFreshness) map[SemanticFieldKey]SemanticFieldFreshness {
	if *entry == nil {
		*entry = make(map[SemanticFieldKey]SemanticFieldFreshness)
	}
	return *entry
}

// SeedZoneFreshness sets freshness for available zone fields.
func SeedZoneFreshness(snapshot *NativeZoneSnapshot, source SemanticSnapshotSource, stale bool) {
	if snapshot == nil {
		return
	}
	fields := EnsureFieldFreshnessMap(&snapshot.FieldFreshness)
	if strings.TrimSpace(snapshot.Name) != "" {
		SetFieldFreshness(fields, ZoneFieldName, source, stale)
	}
	if strings.TrimSpace(snapshot.OperatingMode) != "" {
		SetFieldFreshness(fields, ZoneFieldOperatingMode, source, stale)
	}
	if strings.TrimSpace(snapshot.Preset) != "" {
		SetFieldFreshness(fields, ZoneFieldPreset, source, stale)
	}
	if strings.TrimSpace(snapshot.HvacAction) != "" {
		SetFieldFreshness(fields, ZoneFieldHvacAction, source, stale)
	}
	if len(snapshot.AllowedModes) > 0 {
		SetFieldFreshness(fields, ZoneFieldAllowedModes, source, stale)
	}
	if snapshot.CurrentTempC != nil {
		SetFieldFreshness(fields, ZoneFieldCurrentTempC, source, stale)
	}
	if snapshot.TargetTempC != nil {
		SetFieldFreshness(fields, ZoneFieldTargetTempC, source, stale)
	}
	if snapshot.HumidityPct != nil {
		SetFieldFreshness(fields, ZoneFieldCurrentHumidityPct, source, stale)
	}
	if strings.TrimSpace(snapshot.ConfigurationHeatingOperationMode) != "" {
		SetFieldFreshness(fields, ZoneFieldZoneOperationModeRaw, source, stale)
	}
	if strings.TrimSpace(snapshot.StateSpecialFunction) != "" {
		SetFieldFreshness(fields, ZoneFieldSpecialFunctionRaw, source, stale)
	}
	if snapshot.ConfigurationRoomTemperatureZoneMappingRaw != nil {
		SetFieldFreshness(fields, ZoneFieldRoomTemperatureZoneMappingRaw, source, stale)
	}
	if snapshot.ConfigurationAssociatedCircuitRaw != nil {
		SetFieldFreshness(fields, ZoneFieldZoneCircuitIndexRaw, source, stale)
	}
	if snapshot.ConfigurationCircuitTypeRaw != nil {
		SetFieldFreshness(fields, ZoneFieldCircuitTypeRaw, source, stale)
	}
	if snapshot.StateValveStatusRaw != nil {
		SetFieldFreshness(fields, ZoneFieldZoneValveStatusRaw, source, stale)
	}
	if snapshot.QuickVetoTempC != nil {
		SetFieldFreshness(fields, ZoneFieldQuickVetoTempC, source, stale)
	}
	if snapshot.QuickVetoDurationH != nil {
		SetFieldFreshness(fields, ZoneFieldQuickVetoDurationH, source, stale)
	}
	if snapshot.QuickVetoEndTime != "" {
		SetFieldFreshness(fields, ZoneFieldQuickVetoEndTime, source, stale)
	}
	if snapshot.QuickVetoEndDate != "" {
		SetFieldFreshness(fields, ZoneFieldQuickVetoEndDate, source, stale)
	}
	if snapshot.HolidayStartDate != "" {
		SetFieldFreshness(fields, ZoneFieldHolidayStartDate, source, stale)
	}
	if snapshot.HolidayEndDate != "" {
		SetFieldFreshness(fields, ZoneFieldHolidayEndDate, source, stale)
	}
	if snapshot.HolidaySetpointC != nil {
		SetFieldFreshness(fields, ZoneFieldHolidaySetpointC, source, stale)
	}
	if snapshot.HolidayStartTime != "" {
		SetFieldFreshness(fields, ZoneFieldHolidayStartTime, source, stale)
	}
	if snapshot.HolidayEndTime != "" {
		SetFieldFreshness(fields, ZoneFieldHolidayEndTime, source, stale)
	}
}

// SeedDhwFreshness sets freshness for available hot water fields.
func SeedDhwFreshness(snapshot *NativeDhwSnapshot, source SemanticSnapshotSource, stale bool) {
	if snapshot == nil {
		return
	}
	fields := EnsureFieldFreshnessMap(&snapshot.FieldFreshness)
	if strings.TrimSpace(snapshot.OperatingMode) != "" {
		SetFieldFreshness(fields, DhwFieldOperatingMode, source, stale)
	}
	if strings.TrimSpace(snapshot.Preset) != "" {
		SetFieldFreshness(fields, DhwFieldPreset, source, stale)
	}
	if snapshot.CurrentTempC != nil {
		SetFieldFreshness(fields, DhwFieldCurrentTempC, source, stale)
	}
	if snapshot.TargetTempC != nil {
		SetFieldFreshness(fields, DhwFieldTargetTempC, source, stale)
	}
	if strings.TrimSpace(snapshot.ConfigurationDHWOperationMode) != "" {
		SetFieldFreshness(fields, DhwFieldDhwOperationModeRaw, source, stale)
	}
	if strings.TrimSpace(snapshot.StateSpecialFunction) != "" {
		SetFieldFreshness(fields, DhwFieldSpecialFunctionRaw, source, stale)
	}
	if snapshot.HolidayStartDate != "" {
		SetFieldFreshness(fields, DhwFieldHolidayStartDate, source, stale)
	}
	if snapshot.HolidayEndDate != "" {
		SetFieldFreshness(fields, DhwFieldHolidayEndDate, source, stale)
	}
}

// SetFieldFreshness updates the retained source and stale flag for a field.
func SetFieldFreshness(fields map[SemanticFieldKey]SemanticFieldFreshness, field SemanticFieldKey, source SemanticSnapshotSource, stale bool) {
	if fields == nil {
		return
	}
	fields[field] = SemanticFieldFreshness{
		Source: source,
		Stale:  stale,
	}
}
