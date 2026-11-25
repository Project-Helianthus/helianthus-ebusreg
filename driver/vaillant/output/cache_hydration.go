package output

import (
	"fmt"
	"strconv"
	"strings"
)

// CacheHydrationFromSnapshot maps the loaded cache records back into native snapshots.
func CacheHydrationFromSnapshot(snapshot SemanticCacheSnapshot) NativeCacheHydration {
	hydration := NativeCacheHydration{PersistedAt: snapshot.PersistedAt}
	if len(snapshot.Zones) > 0 {
		hydration.Zones = make([]*NativeZoneSnapshot, 0, len(snapshot.Zones))
	}
	for idx, zone := range snapshot.Zones {
		instance := zoneInstanceFromSemanticID(zone.ID, idx)
		hydration.Zones = append(hydration.Zones, zoneSnapshotFromSemanticZone(instance, zone))
	}
	if snapshot.DHW != nil {
		hydration.Dhw = dhwSnapshotFromSemanticStatus(snapshot.DHW)
	}
	if snapshot.Boiler != nil {
		hydration.Boiler = BoilerSnapshotFromGraphQL(snapshot.Boiler)
		hydration.Boiler.ReturnTemperatureC = CloneFloat64Ptr(snapshot.Boiler.State.ReturnTemperatureC)
	}
	return hydration
}

func zoneInstanceFromSemanticID(id string, fallbackIndex int) byte {
	trimmed := strings.TrimSpace(strings.ToLower(id))
	if strings.HasPrefix(trimmed, "zone-") {
		ordinalRaw := strings.TrimPrefix(trimmed, "zone-")
		if ordinal, err := strconv.Atoi(ordinalRaw); err == nil && ordinal > 0 && ordinal <= 256 {
			return byte(ordinal - 1)
		}
	}
	if fallbackIndex < 0 {
		return 0
	}
	if fallbackIndex > 255 {
		return 255
	}
	return byte(fallbackIndex)
}

func zoneSnapshotFromSemanticZone(instance byte, zone Zone) *NativeZoneSnapshot {
	snapshot := &NativeZoneSnapshot{
		Instance:             instance,
		Present:              true,
		Name:                 zone.Name,
		OperatingMode:        zone.Config.OperatingMode,
		Preset:               zone.Config.Preset,
		HvacAction:           zone.State.HvacAction,
		AllowedModes:         append([]string(nil), zone.Config.AllowedModes...),
		CurrentTempC:         CloneFloat64Ptr(zone.State.CurrentTempC),
		TargetTempC:          CloneFloat64Ptr(zone.Config.TargetTempC),
		HumidityPct:          CloneFloat64Ptr(zone.State.CurrentHumidityPct),
		StateSpecialFunction: zone.State.SpecialFunction,
	}
	if zone.Config.AssociatedCircuit != nil {
		v := uint16(*zone.Config.AssociatedCircuit)
		snapshot.ConfigurationAssociatedCircuitRaw = &v
	}
	if zone.Config.RoomTemperatureZoneMapping != nil {
		v := uint16(*zone.Config.RoomTemperatureZoneMapping)
		snapshot.ConfigurationRoomTemperatureZoneMappingRaw = &v
	}
	if zone.Config.CircuitType != "" {
		v := encodeCircuitTypeRaw(zone.Config.CircuitType)
		snapshot.ConfigurationCircuitTypeRaw = &v
	}
	if zone.State.ValvePositionPct != nil {
		v := uint16(*zone.State.ValvePositionPct * 655.35)
		snapshot.StateValveStatusRaw = &v
	}
	snapshot.QuickVetoTempC = CloneFloat64Ptr(zone.Config.QuickVetoSetpointC)
	snapshot.QuickVetoDurationH = CloneFloat64Ptr(zone.Config.QuickVetoDurationH)
	snapshot.HolidayStartDate = zone.Config.HolidayStartDate
	snapshot.HolidayEndDate = zone.Config.HolidayEndDate
	snapshot.HolidaySetpointC = CloneFloat64Ptr(zone.Config.HolidaySetpointC)
	snapshot.HolidayStartTime = zone.Config.HolidayStartTime
	snapshot.HolidayEndTime = zone.Config.HolidayEndTime
	SeedZoneFreshness(snapshot, SemanticSnapshotSourceCache, true)
	return snapshot
}

func encodeCircuitTypeRaw(circuitType string) uint16 {
	switch circuitType {
	case "radiator":
		return 0
	case "underfloor":
		return 1
	case "mixed":
		return 2
	default:
		var n uint16
		if _, err := fmt.Sscanf(circuitType, "unknown_%d", &n); err == nil {
			return n
		}
		return 0
	}
}

func dhwSnapshotFromSemanticStatus(status *DhwStatus) *NativeDhwSnapshot {
	if status == nil {
		return nil
	}
	snapshot := &NativeDhwSnapshot{
		OperatingMode:        status.Config.OperatingMode,
		Preset:               status.Config.Preset,
		CurrentTempC:         CloneFloat64Ptr(status.State.CurrentTempC),
		TargetTempC:          CloneFloat64Ptr(status.Config.TargetTempC),
		StateSpecialFunction: status.State.SpecialFunction,
		HolidayStartDate:     status.Config.HolidayStartDate,
		HolidayEndDate:       status.Config.HolidayEndDate,
	}
	SeedDhwFreshness(snapshot, SemanticSnapshotSourceCache, true)
	return snapshot
}

// BoilerSnapshotFromGraphQL maps a published boiler status back into a native snapshot.
func BoilerSnapshotFromGraphQL(status *BoilerStatus) *NativeBoilerSnapshot {
	if status == nil {
		return nil
	}
	return MergeBoilerSnapshotNonDestructive(nil, &NativeBoilerSnapshot{
		FlowTemperatureC:         CloneFloat64Ptr(status.State.FlowTemperatureC),
		ReturnTemperatureC:       CloneFloat64Ptr(status.State.ReturnTemperatureC),
		CentralHeatingPumpActive: cloneBoilerBoolPtr(status.State.CentralHeatingPumpActive),
		WaterPressureBar:         CloneFloat64Ptr(status.State.WaterPressureBar),
		ExternalPumpActive:       cloneBoilerBoolPtr(status.State.ExternalPumpActive),
		CirculationPumpActive:    cloneBoilerBoolPtr(status.State.CirculationPumpActive),
		GasValveActive:           cloneBoilerBoolPtr(status.State.GasValveActive),
		FlameActive:              cloneBoilerBoolPtr(status.State.FlameActive),
		DiverterValvePositionPct: CloneFloat64Ptr(status.State.DiverterValvePositionPct),
		FanSpeedRpm:              cloneBoilerIntPtr(status.State.FanSpeedRpm),
		TargetFanSpeedRpm:        cloneBoilerIntPtr(status.State.TargetFanSpeedRpm),
		IonisationVoltageUa:      CloneFloat64Ptr(status.State.IonisationVoltageUa),
		DhwWaterFlowLpm:          CloneFloat64Ptr(status.State.DhwWaterFlowLpm),
		DhwDemandActive:          cloneBoilerBoolPtr(status.State.DhwDemandActive),
		HeatingSwitchActive:      cloneBoilerBoolPtr(status.State.HeatingSwitchActive),
		StorageLoadPumpPct:       CloneFloat64Ptr(status.State.StorageLoadPumpPct),
		ModulationPct:            CloneFloat64Ptr(status.State.ModulationPct),
		PrimaryCircuitFlowLpm:    CloneFloat64Ptr(status.State.PrimaryCircuitFlowLpm),
		FlowTempDesiredC:         CloneFloat64Ptr(status.State.FlowTempDesiredC),
		DhwTempDesiredC:          CloneFloat64Ptr(status.State.DhwTempDesiredC),
		StateNumber:              cloneBoilerIntPtr(status.State.StateNumber),
		DhwTemperatureC:          CloneFloat64Ptr(status.State.DhwTemperatureC),
		DhwTargetTemperatureC:    CloneFloat64Ptr(status.State.DhwTargetTemperatureC),
		DhwOperatingMode:         cloneBoilerModeIntPtr(status.Config.DhwOperatingMode),
		FlowsetHcMaxC:            CloneFloat64Ptr(status.Config.FlowsetHcMaxC),
		FlowsetHwcMaxC:           CloneFloat64Ptr(status.Config.FlowsetHwcMaxC),
		PartloadHcKW:             CloneFloat64Ptr(status.Config.PartloadHcKW),
		PartloadHwcKW:            CloneFloat64Ptr(status.Config.PartloadHwcKW),
		HeatingStatusRaw:         cloneBoilerIntPtr(status.Diagnostics.HeatingStatusRaw),
		DhwStatusRaw:             cloneBoilerIntPtr(status.Diagnostics.DhwStatusRaw),
		CentralHeatingHours:      CloneFloat64Ptr(status.Diagnostics.CentralHeatingHours),
		DhwHours:                 CloneFloat64Ptr(status.Diagnostics.DhwHours),
		CentralHeatingStarts:     cloneBoilerIntPtr(status.Diagnostics.CentralHeatingStarts),
		DhwStarts:                cloneBoilerIntPtr(status.Diagnostics.DhwStarts),
		PumpHours:                CloneFloat64Ptr(status.Diagnostics.PumpHours),
		FanHours:                 CloneFloat64Ptr(status.Diagnostics.FanHours),
		DeactivationsIFC:         cloneBoilerIntPtr(status.Diagnostics.DeactivationsIFC),
		DeactivationsTemplimiter: cloneBoilerIntPtr(status.Diagnostics.DeactivationsTemplimiter),
	})
}

// CacheZones overlays retained zone values on the published records for persistence.
func CacheZones(published []Zone, overlays map[byte]NativeZoneCacheOverlay) []Zone {
	out := make([]Zone, 0, len(published))
	for idx, zone := range published {
		instance := zoneInstanceFromSemanticID(zone.ID, idx)
		if overlay, ok := overlays[instance]; ok {
			zone.State.CurrentTempC = CloneFloat64Ptr(overlay.CurrentTempC)
			zone.State.CurrentHumidityPct = CloneFloat64Ptr(overlay.HumidityPct)
			zone.Config.RoomTemperatureZoneMapping = decodeRoomTemperatureZoneMapping(overlay.RoomMappingRaw)
			zone.Config.AssociatedCircuit = decodeAssociatedCircuit(overlay.AssociatedCircuitRaw)
		}
		out = append(out, zone)
	}
	return out
}
