package owner

import (
	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/output"
	"github.com/Project-Helianthus/helianthus-ebusreg/vaillant/productids"
)

func (o *vaillantSemanticPublication) publishZoneViews(views []vaillantZoneView, source semanticSnapshotSource) bool {
	zones := make([]output.Zone, 0, len(views))
	for _, view := range views {
		if view.snapshot == nil {
			continue
		}
		zones = append(zones, output.ZoneFromView(output.NativeZoneView{
			Instance:               view.instance,
			Snapshot:               nativeZoneSnapshot(view.snapshot),
			RadioRoomTempC:         view.radioRoomTempC,
			RadioRoomHumidityPct:   view.radioRoomHumidityPct,
			InferredRoomMappingRaw: view.inferredRoomMappingRaw,
		}))
	}

	previous := o.state.Zones()
	if source == semanticSnapshotSourceCache && len(zones) == 0 && len(previous) > 0 {
		return false
	}
	switch source {
	case semanticSnapshotSourceCache:
		o.state.SetZonesFromCache(zones)
	default:
		o.state.SetZones(zones)
	}

	if o.events != nil {
		prevByID := make(map[string]output.Zone, len(previous))
		for _, z := range previous {
			prevByID[z.ID] = z
		}
		for _, zone := range zones {
			if !output.ZoneEquals(prevByID[zone.ID], zone) {
				o.events.PublishZoneUpdate(zone)
			}
		}
	}
	return true
}

func (o *vaillantSemanticPublication) publishDhwView(view vaillantDhwView, source semanticSnapshotSource) bool {
	previous := o.state.DHW()
	if view.snapshot == nil {
		if source == semanticSnapshotSourceCache && previous != nil && !view.expired {
			return false
		}
		switch source {
		case semanticSnapshotSourceCache:
			o.state.SetDHWFromCache(nil)
		default:
			o.state.SetDHW(nil)
		}
		if previous != nil && o.events != nil {
			o.events.PublishDHWUpdate(nil)
		}
		return true
	}

	current := output.DhwStatusFromSnapshot(nativeDhwSnapshot(view.snapshot))
	switch source {
	case semanticSnapshotSourceCache:
		o.state.SetDHWFromCache(current)
	default:
		o.state.SetDHW(current)
	}
	if o.events != nil && !output.DhwEquals(previous, current) {
		o.events.PublishDHWUpdate(current)
	}
	return true
}

func (o *vaillantSemanticPublication) publishCircuitsView(view vaillantCircuitsView, source semanticSnapshotSource) bool {
	out := make([]output.CircuitStatus, 0, len(view.circuits))
	for _, snapshot := range view.circuits {
		if snapshot == nil {
			continue
		}
		out = append(out, output.CircuitStatusFromSnapshot((*output.NativeCircuitSnapshot)(snapshot), view.managedByFunctionModule))
	}

	previous := o.state.Circuits()
	if len(out) == 0 && previous == nil {
		return false
	}
	if source == semanticSnapshotSourceCache && len(out) == 0 && len(previous) > 0 {
		return false
	}
	switch source {
	case semanticSnapshotSourceCache:
		o.state.SetCircuitsFromCache(out)
	default:
		o.state.SetCircuits(out)
	}
	return true
}

func (o *vaillantSemanticPublication) publishRadioDeviceViews(views []*vaillantRadioDeviceSnapshot, source semanticSnapshotSource) bool {
	out := make([]output.RadioDevice, 0, len(views))
	for _, snapshot := range views {
		if snapshot == nil {
			continue
		}
		out = append(out, output.RadioDeviceFromSnapshot((*output.NativeRadioDeviceSnapshot)(snapshot)))
	}

	previous := o.state.RadioDevices()
	if source == semanticSnapshotSourceCache && len(out) == 0 && len(previous) > 0 {
		return false
	}
	switch source {
	case semanticSnapshotSourceCache:
		o.state.SetRadioDevicesFromCache(out)
	default:
		o.state.SetRadioDevices(out)
	}
	if o.events != nil && !output.RadioDevicesEqual(previous, out) {
		o.events.PublishRadioDevicesUpdate(out)
	}
	return true
}

func (o *vaillantSemanticPublication) publishRegulatorCapability(capability productids.ControllerCapability) {
	if o == nil || o.state == nil {
		return
	}
	o.state.SetVaillantRegulatorCapability(output.RegulatorCapabilityFromNative(capability))
}

func (o *vaillantSemanticPublication) publishBoilerStatusView(snapshot *vaillantBoilerSnapshot, source semanticSnapshotSource) bool {
	previous := o.state.BoilerStatus()
	if snapshot == nil {
		if previous != nil {
			switch source {
			case semanticSnapshotSourceCache:
				o.state.SetBoilerStatusFromCache(nil)
			default:
				o.state.SetBoilerStatus(nil)
			}
		}
		return false
	}

	current := output.BoilerStatusFromSnapshot((*output.NativeBoilerSnapshot)(snapshot))
	switch source {
	case semanticSnapshotSourceCache:
		o.state.SetBoilerStatusFromCache(current)
	default:
		o.state.SetBoilerStatus(current)
	}
	if o.events != nil && !output.BoilerStatusEquals(previous, current) {
		o.events.PublishBoilerStatusUpdate(current)
	}
	return true
}

func (o *vaillantSemanticPublication) withdrawBoilerFlowTemperature() {
	if o == nil || o.state == nil {
		return
	}
	status := o.state.BoilerStatus()
	if status == nil || status.State.FlowTemperatureC == nil {
		return
	}
	status.State.FlowTemperatureC = nil
	o.state.SetBoilerStatus(status)
	if o.events != nil {
		o.events.PublishBoilerStatusUpdate(status)
	}
}

func (o *vaillantSemanticPublication) publishSystemView(view vaillantSystemView, source semanticSnapshotSource) bool {
	if view.snapshot == nil {
		switch source {
		case semanticSnapshotSourceCache:
			o.state.SetSystemFromCache(nil)
		default:
			o.state.SetSystem(nil)
		}
		return false
	}

	current := output.SystemStatusFromView(output.NativeSystemView{
		Snapshot:          nativeSystemSnapshot(view.snapshot),
		DhwBivalenceBound: outputConfigBound(view.dhwBivalenceBound),
	})
	switch source {
	case semanticSnapshotSourceCache:
		o.state.SetSystemFromCache(current)
	default:
		o.state.SetSystem(current)
	}
	return true
}
