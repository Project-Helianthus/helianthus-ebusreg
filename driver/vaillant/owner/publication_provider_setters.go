package owner

import (
	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/output"
)

func (provider *retainedSemanticState) SetZones(zones []output.Zone) {
	provider.setZonesWithSource(zones, semanticDataSourceLive, "zones_live_update")
}

func (provider *retainedSemanticState) SetZonesFromCache(zones []output.Zone) {
	provider.setZonesWithSource(zones, semanticDataSourceCache, "zones_cache_update")
}

func (provider *retainedSemanticState) setZonesWithSource(zones []output.Zone, source semanticDataSource, reason string) {
	if provider == nil {
		return
	}
	receivedZones := zones != nil
	zonesCopy := make([]output.Zone, len(zones))
	for i, z := range zones {
		zonesCopy[i] = cloneZone(z)
	}
	published := receivedZones && (source == semanticDataSourceLive || len(zonesCopy) > 0)
	var transition *phaseTransitionLog
	provider.mu.Lock()
	provider.zones = zonesCopy
	if published {
		provider.zonePublished = true
	}
	if published {
		if source == semanticDataSourceLive {
			provider.zoneLiveSeen = true
		}
		transition = provider.recordEpochUpdateLocked(source, semanticStreamZones, reason)
	}
	provider.mu.Unlock()
	provider.logPhaseTransition(transition)
}

func (provider *retainedSemanticState) SetDHW(status *output.DhwStatus) {
	provider.setDHWWithSource(status, semanticDataSourceLive, "dhw_live_update")
}

func (provider *retainedSemanticState) SetDHWFromCache(status *output.DhwStatus) {
	provider.setDHWWithSource(status, semanticDataSourceCache, "dhw_cache_update")
}

func (provider *retainedSemanticState) setDHWWithSource(status *output.DhwStatus, source semanticDataSource, reason string) {
	if provider == nil {
		return
	}
	var transition *phaseTransitionLog
	provider.mu.Lock()
	if status == nil {
		provider.dhw = nil
		provider.mu.Unlock()
		return
	}
	cp := cloneDhwStatus(status)
	provider.dhw = cp
	provider.dhwPublished = true
	if source == semanticDataSourceLive {
		provider.dhwLiveSeen = true
	}
	transition = provider.recordEpochUpdateLocked(source, semanticStreamDHW, reason)
	provider.mu.Unlock()
	provider.logPhaseTransition(transition)
}

func (provider *retainedSemanticState) SetCircuits(circuits []output.CircuitStatus) {
	provider.setCircuitsWithSource(circuits, semanticDataSourceLive, "circuits_live_update")
}

func (provider *retainedSemanticState) SetCircuitsFromCache(circuits []output.CircuitStatus) {
	provider.setCircuitsWithSource(circuits, semanticDataSourceCache, "circuits_cache_update")
}

func (provider *retainedSemanticState) setCircuitsWithSource(circuits []output.CircuitStatus, source semanticDataSource, reason string) {
	if provider == nil {
		return
	}

	receivedCircuits := circuits != nil
	circuitsCopy := make([]output.CircuitStatus, len(circuits))
	for i, status := range circuits {
		circuitsCopy[i] = cloneCircuitStatus(status)
	}
	published := receivedCircuits && (source == semanticDataSourceLive || len(circuitsCopy) > 0)

	var transition *phaseTransitionLog
	provider.mu.Lock()
	if output.CircuitsEqual(provider.circuits, circuitsCopy) && (!published || provider.circuitPublished) {
		if source == semanticDataSourceLive && published {
			provider.circuitLiveSeen = true
		}
		provider.mu.Unlock()
		return
	}
	provider.circuits = circuitsCopy
	if published {
		provider.circuitPublished = true
		if source == semanticDataSourceLive {
			provider.circuitLiveSeen = true
		}
		transition = provider.recordEpochUpdateLocked(source, semanticStreamCircuits, reason)
	}
	provider.mu.Unlock()
	provider.logPhaseTransition(transition)
}

func (provider *retainedSemanticState) SetRadioDevices(devices []output.RadioDevice) {
	provider.setRadioDevicesWithSource(devices, semanticDataSourceLive, "radio_devices_live_update")
}

func (provider *retainedSemanticState) SetRadioDevicesFromCache(devices []output.RadioDevice) {
	provider.setRadioDevicesWithSource(devices, semanticDataSourceCache, "radio_devices_cache_update")
}

func (provider *retainedSemanticState) setRadioDevicesWithSource(devices []output.RadioDevice, source semanticDataSource, reason string) {
	if provider == nil {
		return
	}

	var devicesCopy []output.RadioDevice
	if devices != nil {
		devicesCopy = make([]output.RadioDevice, len(devices))
		for i, device := range devices {
			devicesCopy[i] = cloneRadioDevice(device)
		}
	}

	var transition *phaseTransitionLog
	provider.mu.Lock()
	samePublishedShape := (provider.radio == nil) == (devicesCopy == nil)
	if samePublishedShape && output.RadioDevicesEqual(provider.radio, devicesCopy) {
		if source == semanticDataSourceLive && len(devicesCopy) > 0 {
			provider.radioLiveSeen = true
		}
		provider.mu.Unlock()
		return
	}
	provider.radio = devicesCopy
	if len(devicesCopy) > 0 {
		provider.radioPublished = true
		if source == semanticDataSourceLive {
			provider.radioLiveSeen = true
		}
		transition = provider.recordEpochUpdateLocked(source, semanticStreamRadioDevices, reason)
	} else if source == semanticDataSourceLive {
		provider.radioLiveSeen = true
	}
	provider.mu.Unlock()
	provider.logPhaseTransition(transition)
}

func (provider *retainedSemanticState) BoilerStatus() *output.BoilerStatus {
	if provider == nil {
		return nil
	}
	provider.mu.RLock()
	defer provider.mu.RUnlock()
	if provider.boiler == nil {
		return nil
	}
	cp := *provider.boiler
	cp.State = cloneBoilerState(cp.State)
	cp.Config = cloneBoilerConfig(cp.Config)
	cp.Diagnostics = cloneBoilerDiagnostics(cp.Diagnostics)
	return &cp
}

func (provider *retainedSemanticState) SetBoilerStatus(status *output.BoilerStatus) {
	provider.setBoilerStatusWithSource(status, semanticDataSourceLive, "boiler_live_update")
}

func (provider *retainedSemanticState) SetBoilerStatusFromCache(status *output.BoilerStatus) {
	provider.setBoilerStatusWithSource(status, semanticDataSourceCache, "boiler_cache_update")
}

func (provider *retainedSemanticState) setBoilerStatusWithSource(status *output.BoilerStatus, source semanticDataSource, reason string) {
	if provider == nil {
		return
	}
	var transition *phaseTransitionLog
	provider.mu.Lock()
	if status == nil {
		provider.boiler = nil
		provider.mu.Unlock()
		return
	}
	cp := *status
	cp.State = cloneBoilerState(cp.State)
	cp.Config = cloneBoilerConfig(cp.Config)
	cp.Diagnostics = cloneBoilerDiagnostics(cp.Diagnostics)
	provider.boiler = &cp
	provider.boilerPublished = true
	if source == semanticDataSourceLive {
		provider.boilerLiveSeen = true
	}
	transition = provider.recordEpochUpdateLocked(source, semanticStreamBoiler, reason)
	provider.mu.Unlock()
	provider.logPhaseTransition(transition)
}
