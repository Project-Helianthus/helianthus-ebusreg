package owner

import (
	"fmt"
	"log"
	"slices"
	"strings"
)

func (p *vaillantSemanticPoller) snapshotZones() (byte, []byte) {
	p.mu.Lock()
	defer p.mu.Unlock()
	controller := p.controller
	if controller == 0 || len(p.zones) == 0 {
		return controller, nil
	}
	instances := make([]byte, 0, len(p.zones))
	for instance := range p.zones {
		instances = append(instances, instance)
	}
	slices.Sort(instances)
	return controller, instances
}

// zoneViewsLocked builds the native views of the retained zones in instance
// order, with the radio room state and the inferred room mapping derived from
// the retained radio devices. The caller holds p.mu.
func (p *vaillantSemanticPoller) zoneViewsLocked() []vaillantZoneView {
	instances := make([]byte, 0, len(p.zones))
	for instance := range p.zones {
		instances = append(instances, instance)
	}
	slices.Sort(instances)
	views := make([]vaillantZoneView, 0, len(instances))
	for _, instance := range instances {
		entry := p.zones[instance]
		if entry == nil {
			continue
		}
		view := vaillantZoneView{instance: instance, snapshot: cloneZoneSnapshot(entry)}
		if entry.CurrentTempC == nil || entry.HumidityPct == nil {
			view.radioRoomTempC, view.radioRoomHumidityPct = p.mappedRadioRoomStateForZoneLocked(entry)
		}
		if entry.ConfigurationRoomTemperatureZoneMappingRaw == nil {
			view.inferredRoomMappingRaw = p.inferRadioRoomMappingForZoneLocked(entry.Instance)
		}
		views = append(views, view)
	}
	return views
}

func (p *vaillantSemanticPoller) mappedRadioRoomStateForZoneLocked(entry *vaillantZoneSnapshot) (*float64, *float64) {
	if p == nil || entry == nil {
		return nil, nil
	}
	mappingRaw := entry.ConfigurationRoomTemperatureZoneMappingRaw
	if mappingRaw == nil {
		mappingRaw = p.inferRadioRoomMappingForZoneLocked(entry.Instance)
	}
	if mappingRaw == nil {
		return nil, nil
	}
	mapping := *mappingRaw
	if mapping == 0 || mapping == 0xFFFF {
		return nil, nil
	}
	if decodeRoomTemperatureZoneMapping(&mapping) == nil {
		return nil, nil
	}

	var snapshot *vaillantRadioDeviceSnapshot
	switch mapping {
	case 1:
		snapshot = p.firstConnectedRadioDeviceLocked(remoteRegulators.group)
	default:
		if mapping > 0x100 {
			return nil, nil
		}
		snapshot = p.radioDevices[radioDeviceKey{Group: remoteThermostats.group, Instance: byte(mapping - 1)}]
	}
	if snapshot == nil || !radioDeviceConnected(snapshot) {
		return nil, nil
	}
	return cloneFloat64Ptr(snapshot.RoomTemperatureC), cloneFloat64Ptr(snapshot.RoomHumidityPct)
}

func (p *vaillantSemanticPoller) inferRadioRoomMappingForZoneLocked(instance byte) *uint16 {
	if p == nil {
		return nil
	}
	zoneNumber := int(instance) + 1
	matches := make([]byte, 0, 1)
	for key, snapshot := range p.radioDevices {
		if key.Group != remoteThermostats.group || snapshot == nil || !radioDeviceConnected(snapshot) || snapshot.ZoneAssignment == nil {
			continue
		}
		if int(*snapshot.ZoneAssignment) != zoneNumber {
			continue
		}
		matches = append(matches, key.Instance)
	}
	if len(matches) != 1 {
		return nil
	}
	mapping := uint16(matches[0]) + 1
	return &mapping
}

func (p *vaillantSemanticPoller) firstConnectedRadioDeviceLocked(group byte) *vaillantRadioDeviceSnapshot {
	if p == nil {
		return nil
	}
	var selected *vaillantRadioDeviceSnapshot
	var selectedInstance byte
	for key, snapshot := range p.radioDevices {
		if key.Group != group || snapshot == nil || !radioDeviceConnected(snapshot) {
			continue
		}
		if selected == nil || key.Instance < selectedInstance {
			selected = snapshot
			selectedInstance = key.Instance
		}
	}
	return selected
}

func radioDeviceConnected(snapshot *vaillantRadioDeviceSnapshot) bool {
	return snapshot != nil && snapshot.DeviceConnected != nil && *snapshot.DeviceConnected
}

func (p *vaillantSemanticPoller) recordZonePresenceTransitionLocked(instance byte, previous, next zonePresenceState, hitStreak, missStreak int) {
	if previous == next {
		return
	}
	semanticZonePresenceTransitionsTotal.Add(fmt.Sprintf("%s->%s", previous, next), 1)
	log.Printf(
		"semantic_zone_presence_transition instance=%d from=%s to=%s hit_streak=%d miss_streak=%d zone_count=%d",
		instance,
		previous,
		next,
		hitStreak,
		missStreak,
		len(p.zones),
	)
}

func resolveAssociatedCircuitInstance(roomTemperatureZoneMapping *uint16, zoneInstance byte) byte {
	if roomTemperatureZoneMapping == nil {
		return zoneInstance
	}
	value := *roomTemperatureZoneMapping
	switch value {
	case 0xFF, 0xFFFF, 0:
		return zoneInstance
	default:
		if value >= 1 && value <= 0x20 {
			return byte(value - 1)
		}
		return zoneInstance
	}
}

func deriveZoneModeAndPreset(opMode, specialFunction, circuitType *uint16, hasCircuitType bool) (string, string, []string) {
	allowed := deriveZoneAllowedModes(circuitType, hasCircuitType)
	var mode string
	switch value := uint16Value(opMode, 1); value {
	case 0:
		mode = "off"
	case 1:
		mode = "auto"
	case 2:
		mode = deriveManualOperatingMode(circuitType, hasCircuitType)
	default:
		mode = "auto"
	}

	preset := derivePresetFromSpecialFunctionAndOpMode(specialFunction, opMode)
	if preset == "" {
		switch mode {
		case "auto":
			preset = "schedule"
		default:
			preset = "manual"
		}
	}
	return mode, preset, allowed
}

func deriveZoneAllowedModes(circuitType *uint16, hasCircuitType bool) []string {
	if !hasCircuitType || circuitType == nil {
		return []string{"off", "auto", "heat"}
	}
	switch *circuitType {
	case 0:
		return []string{"off"}
	case 2:
		return []string{"off", "auto", "cool"}
	default:
		return []string{"off", "auto", "heat"}
	}
}

func deriveZoneHvacAction(valveStatus, circuitType *uint16, hasCircuitType bool) string {
	if valveStatus == nil {
		return ""
	}
	switch *valveStatus {
	case 0:
		return "idle"
	case 1:
		if hasCircuitType && circuitType != nil && *circuitType == 2 {
			return "cooling"
		}
		return "heating"
	default:
		return ""
	}
}

func composeZoneName(primary, prefix, suffix string) string {
	if primary = strings.TrimSpace(primary); primary != "" {
		return primary
	}
	parts := make([]string, 0, 2)
	if prefix = strings.TrimSpace(prefix); prefix != "" {
		parts = append(parts, prefix)
	}
	if suffix = strings.TrimSpace(suffix); suffix != "" {
		parts = append(parts, suffix)
	}
	return strings.TrimSpace(strings.Join(parts, " "))
}
