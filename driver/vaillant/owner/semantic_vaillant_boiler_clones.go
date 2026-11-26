package owner

import (
	"log"
	"strings"

	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
	"github.com/Project-Helianthus/helianthus-ebusreg/vaillant/productids"
)

func cloneBoilerSnapshot(snapshot *vaillantBoilerSnapshot) *vaillantBoilerSnapshot {
	if snapshot == nil {
		return nil
	}
	return &vaillantBoilerSnapshot{
		FlowTemperatureC:         cloneFloat64Ptr(snapshot.FlowTemperatureC),
		ReturnTemperatureC:       cloneFloat64Ptr(snapshot.ReturnTemperatureC),
		CentralHeatingPumpActive: cloneBoilerBoolPtr(snapshot.CentralHeatingPumpActive),
		WaterPressureBar:         cloneFloat64Ptr(snapshot.WaterPressureBar),
		ExternalPumpActive:       cloneBoilerBoolPtr(snapshot.ExternalPumpActive),
		CirculationPumpActive:    cloneBoilerBoolPtr(snapshot.CirculationPumpActive),
		GasValveActive:           cloneBoilerBoolPtr(snapshot.GasValveActive),
		FlameActive:              cloneBoilerBoolPtr(snapshot.FlameActive),
		DiverterValvePositionPct: cloneFloat64Ptr(snapshot.DiverterValvePositionPct),
		FanSpeedRpm:              cloneBoilerIntPtr(snapshot.FanSpeedRpm),
		TargetFanSpeedRpm:        cloneBoilerIntPtr(snapshot.TargetFanSpeedRpm),
		IonisationVoltageUa:      cloneFloat64Ptr(snapshot.IonisationVoltageUa),
		DhwWaterFlowLpm:          cloneFloat64Ptr(snapshot.DhwWaterFlowLpm),
		DhwDemandActive:          cloneBoilerBoolPtr(snapshot.DhwDemandActive),
		HeatingSwitchActive:      cloneBoilerBoolPtr(snapshot.HeatingSwitchActive),
		StorageLoadPumpPct:       cloneFloat64Ptr(snapshot.StorageLoadPumpPct),
		ModulationPct:            cloneFloat64Ptr(snapshot.ModulationPct),
		PrimaryCircuitFlowLpm:    cloneFloat64Ptr(snapshot.PrimaryCircuitFlowLpm),
		FlowTempDesiredC:         cloneFloat64Ptr(snapshot.FlowTempDesiredC),
		DhwTempDesiredC:          cloneFloat64Ptr(snapshot.DhwTempDesiredC),
		StateNumber:              cloneBoilerIntPtr(snapshot.StateNumber),
		DhwTemperatureC:          cloneFloat64Ptr(snapshot.DhwTemperatureC),
		DhwTargetTemperatureC:    cloneFloat64Ptr(snapshot.DhwTargetTemperatureC),
		DhwOperatingMode:         cloneBoilerIntPtr(snapshot.DhwOperatingMode),
		FlowsetHcMaxC:            cloneFloat64Ptr(snapshot.FlowsetHcMaxC),
		FlowsetHwcMaxC:           cloneFloat64Ptr(snapshot.FlowsetHwcMaxC),
		PartloadHcKW:             cloneFloat64Ptr(snapshot.PartloadHcKW),
		PartloadHwcKW:            cloneFloat64Ptr(snapshot.PartloadHwcKW),
		HeatingStatusRaw:         cloneBoilerIntPtr(snapshot.HeatingStatusRaw),
		DhwStatusRaw:             cloneBoilerIntPtr(snapshot.DhwStatusRaw),
		CentralHeatingHours:      cloneFloat64Ptr(snapshot.CentralHeatingHours),
		DhwHours:                 cloneFloat64Ptr(snapshot.DhwHours),
		CentralHeatingStarts:     cloneBoilerIntPtr(snapshot.CentralHeatingStarts),
		DhwStarts:                cloneBoilerIntPtr(snapshot.DhwStarts),
		PumpHours:                cloneFloat64Ptr(snapshot.PumpHours),
		FanHours:                 cloneFloat64Ptr(snapshot.FanHours),
		DeactivationsIFC:         cloneBoilerIntPtr(snapshot.DeactivationsIFC),
		DeactivationsTemplimiter: cloneBoilerIntPtr(snapshot.DeactivationsTemplimiter),
		InstallerMenuCode:        cloneBoilerIntPtr(snapshot.InstallerMenuCode),
		PhoneNumber:              cloneStringPtr(snapshot.PhoneNumber),
		HoursTillService:         cloneBoilerIntPtr(snapshot.HoursTillService),
	}
}

func cloneBoilerBoolPtr(value *bool) *bool {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}

func cloneBoilerIntPtr(value *int) *int {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}

// CloneBoilerModeIntPtr implements the clone boiler mode int ptr step of the Vaillant owner.
func CloneBoilerModeIntPtr(mode *string) *int {
	if mode == nil {
		return nil
	}
	switch strings.ToUpper(strings.TrimSpace(*mode)) {
	case "OFF":
		value := 0
		return &value
	case "ON":
		value := 1
		return &value
	case "AUTO":
		value := 2
		return &value
	default:
		return nil
	}
}

func (p *vaillantSemanticPoller) findBoilerAddress() byte {
	if p == nil || p.reg == nil {
		return 0
	}

	// Phase C M-C6b: select the routing target byte for B5.09
	// boiler reads/writes. An aliased canonical pair (e.g. BAI
	// 0x03↔0x08) whose display is the initiator side (0x03) must
	// resolve to 0x08 here, otherwise refreshBoilerStatusB509
	// builds frames with Target: 0x03 and the read/write fails.
	// TargetAddressForRouting prefers AddressByRole(SlotRoleSlave)
	// and falls back to PrimaryDisplayAddress when no target-role
	// face exists.
	// P9.3 — race-free boiler-candidate enumeration via
	// IterateSnapshots + isBoilerDeviceSnapshotCandidate +
	// SnapshotTargetAddressForRouting.
	selected := byte(0)
	p.reg.IterateSnapshots(func(snap registry.DeviceEntrySnapshot) bool {
		if !isBoilerDeviceSnapshotCandidate(snap) {
			return true
		}
		addr := snapshotTargetAddressForRouting(snap)
		if addr == 0x08 {
			selected = addr
			return false
		}
		if selected == 0 || addr < selected {
			selected = addr
		}
		return true
	})
	return selected
}

// isBoilerDeviceSnapshotCandidate is the snapshot-based counterpart of
// isBoilerDeviceCandidate (P9.3). Reads the snapshot's DeviceID
// (race-free) and applies the same prefix-based BAI/BMU classification.
func isBoilerDeviceSnapshotCandidate(snap registry.DeviceEntrySnapshot) bool {
	normalizedID := normalizeDeviceID(snap.DeviceID)
	if strings.HasPrefix(normalizedID, "BAI") {
		return true
	}
	return strings.HasPrefix(normalizedID, "BMU")
}

// heatPumpDeviceIDFamilies are the identification prefixes that the public
// device-class table lists for heat-pump-class devices. A prefix is an
// identity hint only; it authorizes no register read or decoding.
var heatPumpDeviceIDFamilies = []string{"EHP", "HMU", "VWL", "VWS", "VWZ"}

// isHeatPumpDeviceSnapshotCandidate reports whether snap carries a
// heat-pump-class identity: a heat-pump device-ID family or a Vaillant serial
// whose part number the product catalog classifies with role "Heat Pump" in
// every row. A boiler-class identity is never heat-pump class, and an
// identity that is neither is unknown, not heat-pump class.
func isHeatPumpDeviceSnapshotCandidate(snap registry.DeviceEntrySnapshot, catalog productids.Catalog, catalogErr error) bool {
	if isBoilerDeviceSnapshotCandidate(snap) {
		return false
	}
	normalizedID := normalizeDeviceID(snap.DeviceID)
	for _, family := range heatPumpDeviceIDFamilies {
		if strings.HasPrefix(normalizedID, family) {
			return true
		}
	}
	if catalogErr != nil || !strings.EqualFold(snap.Manufacturer, "Vaillant") {
		return false
	}
	partNumber := strings.TrimSpace(extractPartNumberFromSerial(snap.SerialNumber))
	if partNumber == "" {
		return false
	}
	found := false
	for _, record := range catalog.All {
		if strings.TrimSpace(record.PartNumber) != partNumber {
			continue
		}
		if !strings.EqualFold(strings.TrimSpace(record.Role), "Heat Pump") {
			return false
		}
		found = true
	}
	return found
}

// findHeatPumpClassAddress returns the routing address of the lowest-addressed
// heat-pump-class identity in the registry, or 0 when there is none.
func (p *vaillantSemanticPoller) findHeatPumpClassAddress() byte {
	if p == nil || p.reg == nil {
		return 0
	}
	selected := byte(0)
	p.reg.IterateSnapshots(func(snap registry.DeviceEntrySnapshot) bool {
		if !isHeatPumpDeviceSnapshotCandidate(snap, p.catalog, p.catalogErr) {
			return true
		}
		if addr := snapshotTargetAddressForRouting(snap); addr != 0 && (selected == 0 || addr < selected) {
			selected = addr
		}
		return true
	})
	return selected
}

// refreshHeatSourcePathReason sets the startup reason
// no_supported_heat_source_path while the poller has no supported heat-source
// data path: no B524 root, no boiler-class heat source, and a heat-pump-class
// identity in the inventory. It clears the reason as soon as a root or a
// boiler-class heat source exists. A process that has discovered a root once,
// or an inventory the catalog classifies as holding a controller, keeps the
// reason clear: there the missing root is a failing read, not a bus without a
// supported path. It reads the registry only and sends nothing.
func (p *vaillantSemanticPoller) refreshHeatSourcePathReason() {
	if !p.publishing() {
		return
	}
	p.mu.Lock()
	controller, boilerAddress := p.controller, p.boilerAddress
	everDiscovered, regCap := p.controllerEverDiscovered, p.regulatorCapability
	p.mu.Unlock()
	heatPump := byte(0)
	if controller == 0 && boilerAddress == 0 && !everDiscovered && regCap != productids.ControllerPresent {
		heatPump = p.findHeatPumpClassAddress()
	}
	reason, changed := p.publication.setNoSupportedHeatSourcePath(heatPump != 0)
	if !changed {
		return
	}
	if heatPump != 0 {
		log.Printf("semantic_startup_reason reason=%s state=set heat_source=0x%02x", reason, heatPump)
		return
	}
	log.Printf("semantic_startup_reason reason=%s state=cleared controller=0x%02x boiler=0x%02x", reason, controller, boilerAddress)
}
