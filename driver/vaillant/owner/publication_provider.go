package owner

import (
	"expvar"
	"sync"
	"time"

	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/output"
)

var (
	semanticStartupPhaseTransitionsTotal = expvar.NewMap("semantic_startup_phase_transitions_total")
	semanticStartupCurrentPhase          = expvar.NewString("semantic_startup_current_phase")
	semanticCacheEpoch                   = expvar.NewInt("semantic_cache_epoch")
	semanticLiveEpoch                    = expvar.NewInt("semantic_live_epoch")
)

type semanticStartupPhase string

const (
	semanticStartupPhaseBootInit         semanticStartupPhase = "BOOT_INIT"
	semanticStartupPhaseCacheLoadedStale semanticStartupPhase = "CACHE_LOADED_STALE"
	semanticStartupPhaseLiveWarmup       semanticStartupPhase = "LIVE_WARMUP"
	semanticStartupPhaseLiveReady        semanticStartupPhase = "LIVE_READY"
	semanticStartupPhaseDegraded         semanticStartupPhase = "DEGRADED"
)

// semanticStartupReasonNoSupportedHeatSourcePath is the startup reason set
// while the bus offers no data path the semantic poller supports.
const semanticStartupReasonNoSupportedHeatSourcePath = "no_supported_heat_source_path"

type semanticDataSource uint8

const (
	semanticDataSourceCache semanticDataSource = iota
	semanticDataSourceLive
)

type semanticStreamKind uint8

const (
	semanticStreamZones semanticStreamKind = iota
	semanticStreamDHW
	semanticStreamBoiler
	semanticStreamCircuits
	semanticStreamRadioDevices
)

type phaseTransitionLog struct {
	previous   semanticStartupPhase
	next       semanticStartupPhase
	reason     string
	cacheEpoch uint64
	liveEpoch  uint64
}

// retainedSemanticState holds the published records of one owner.
type retainedSemanticState struct {
	mu                          sync.RWMutex
	zones                       []output.Zone
	dhw                         *output.DhwStatus
	circuits                    []output.CircuitStatus
	radio                       []output.RadioDevice
	functionModuleMode          output.FunctionModuleSemanticMode
	functionModuleVerdict       output.FunctionModuleInterpretation
	solar                       *output.SolarStatus
	cylinders                   []output.CylinderStatus
	energy                      *output.EnergyTotals
	boiler                      *output.BoilerStatus
	system                      *output.SystemStatus
	schedules                   *output.ScheduleStatus
	adapter                     *output.AdapterHardwareInfo
	vaillantRegulatorCapability output.VaillantRegulatorCapability

	energyMerge                *energyMergeStore
	energyRevision             uint64
	energyClock                func() time.Time
	energyPassiveState         string
	energyPassiveStateProvider func() string

	phase            semanticStartupPhase
	startupReason    string
	startupUpdatedAt time.Time
	cacheEpoch       uint64
	liveEpoch        uint64
	zonePublished    bool
	dhwPublished     bool
	boilerPublished  bool
	circuitPublished bool
	radioPublished   bool
	zoneLiveSeen     bool
	dhwLiveSeen      bool
	boilerLiveSeen   bool
	circuitLiveSeen  bool
	radioLiveSeen    bool
	phaseLogger      func(string, ...any)
}

func newRetainedSemanticState() *retainedSemanticState {
	semanticStartupCurrentPhase.Set(string(semanticStartupPhaseBootInit))
	semanticCacheEpoch.Set(0)
	semanticLiveEpoch.Set(0)

	return &retainedSemanticState{
		phase:                       semanticStartupPhaseBootInit,
		startupUpdatedAt:            time.Now().UTC(),
		functionModuleMode:          output.FunctionModuleSemanticModeAbsent,
		vaillantRegulatorCapability: output.VaillantRegulatorCapabilityUnknown,
		energyMerge:                 newEnergyMergeStore(),
	}
}

func (provider *retainedSemanticState) StartupPhase() semanticStartupPhase {
	if provider == nil {
		return semanticStartupPhaseBootInit
	}
	provider.mu.RLock()
	defer provider.mu.RUnlock()
	return provider.phase
}

// StartupReason returns the reason shown next to the startup phase, or "" when
// none applies.
func (provider *retainedSemanticState) StartupReason() string {
	if provider == nil {
		return ""
	}
	provider.mu.RLock()
	defer provider.mu.RUnlock()
	return provider.startupReason
}

// SetNoSupportedHeatSourcePath sets or clears the startup reason
// semanticStartupReasonNoSupportedHeatSourcePath and reports whether it
// changed. It changes neither the phase nor the epochs.
func (provider *retainedSemanticState) SetNoSupportedHeatSourcePath(active bool) bool {
	if provider == nil {
		return false
	}
	next := ""
	if active {
		next = semanticStartupReasonNoSupportedHeatSourcePath
	}
	provider.mu.Lock()
	defer provider.mu.Unlock()
	if provider.startupReason == next {
		return false
	}
	provider.startupReason = next
	provider.touchStartupUpdatedAtLocked()
	return true
}

func (provider *retainedSemanticState) StartupEpochs() (cacheEpoch uint64, liveEpoch uint64) {
	if provider == nil {
		return 0, 0
	}
	provider.mu.RLock()
	defer provider.mu.RUnlock()
	return provider.cacheEpoch, provider.liveEpoch
}

func (provider *retainedSemanticState) StartupUpdatedAt() time.Time {
	if provider == nil {
		return time.Time{}
	}
	provider.mu.RLock()
	defer provider.mu.RUnlock()
	return provider.startupUpdatedAt
}

func (provider *retainedSemanticState) Zones() []output.Zone {
	if provider == nil {
		return nil
	}
	provider.mu.RLock()
	defer provider.mu.RUnlock()
	if len(provider.zones) == 0 {
		if provider.zonePublished {
			return []output.Zone{}
		}
		return nil
	}
	zones := make([]output.Zone, len(provider.zones))
	for i, z := range provider.zones {
		zones[i] = cloneZone(z)
	}
	return zones
}

func (provider *retainedSemanticState) DHW() *output.DhwStatus {
	if provider == nil {
		return nil
	}
	provider.mu.RLock()
	defer provider.mu.RUnlock()
	if provider.dhw == nil {
		return nil
	}
	cp := cloneDhwStatus(provider.dhw)
	return cp
}

func (provider *retainedSemanticState) Circuits() []output.CircuitStatus {
	if provider == nil {
		return nil
	}
	provider.mu.RLock()
	defer provider.mu.RUnlock()
	if len(provider.circuits) == 0 {
		if provider.circuitPublished {
			return []output.CircuitStatus{}
		}
		return nil
	}
	out := make([]output.CircuitStatus, len(provider.circuits))
	for i, status := range provider.circuits {
		out[i] = cloneCircuitStatus(status)
	}
	return out
}

func (provider *retainedSemanticState) RadioDevices() []output.RadioDevice {
	if provider == nil {
		return nil
	}
	provider.mu.RLock()
	defer provider.mu.RUnlock()
	if provider.radio == nil {
		return nil
	}
	out := make([]output.RadioDevice, len(provider.radio))
	for i, device := range provider.radio {
		out[i] = cloneRadioDevice(device)
	}
	return out
}

func (provider *retainedSemanticState) FunctionModuleSemanticMode() output.FunctionModuleSemanticMode {
	if provider == nil {
		return output.FunctionModuleSemanticModeAbsent
	}
	provider.mu.RLock()
	defer provider.mu.RUnlock()
	if provider.functionModuleMode == "" {
		return output.FunctionModuleSemanticModeAbsent
	}
	return provider.functionModuleMode
}

func (provider *retainedSemanticState) FunctionModuleInterpretation() output.FunctionModuleInterpretation {
	if provider == nil {
		return output.FunctionModuleInterpretation{}
	}
	provider.mu.RLock()
	defer provider.mu.RUnlock()
	return provider.functionModuleVerdict
}

func (provider *retainedSemanticState) VaillantRegulatorCapability() output.VaillantRegulatorCapability {
	if provider == nil {
		return output.VaillantRegulatorCapabilityUnknown
	}
	provider.mu.RLock()
	defer provider.mu.RUnlock()
	return output.NormalizeVaillantRegulatorCapability(provider.vaillantRegulatorCapability)
}

func (provider *retainedSemanticState) SetVaillantRegulatorCapability(capability output.VaillantRegulatorCapability) {
	if provider == nil {
		return
	}
	capability = output.NormalizeVaillantRegulatorCapability(capability)
	provider.mu.Lock()
	provider.vaillantRegulatorCapability = capability
	provider.mu.Unlock()
}

func (provider *retainedSemanticState) Solar() *output.SolarStatus {
	if provider == nil {
		return nil
	}
	provider.mu.RLock()
	defer provider.mu.RUnlock()
	return cloneSolarStatus(provider.solar)
}

func (provider *retainedSemanticState) Cylinders() []output.CylinderStatus {
	if provider == nil {
		return nil
	}
	provider.mu.RLock()
	defer provider.mu.RUnlock()
	if provider.cylinders == nil {
		return nil
	}
	return cloneCylinderStatuses(provider.cylinders)
}

func (provider *retainedSemanticState) SetFunctionModuleInterpretation(verdict output.FunctionModuleInterpretation) {
	if provider == nil {
		return
	}
	if err := verdict.Validate(); err != nil {
		return
	}
	provider.mu.Lock()
	defer provider.mu.Unlock()
	provider.functionModuleMode = verdict.Mode
	provider.functionModuleVerdict = verdict
}

func (provider *retainedSemanticState) SetFunctionModuleInterpretationFromCache(verdict output.FunctionModuleInterpretation) {
	provider.SetFunctionModuleInterpretation(verdict)
}

func (provider *retainedSemanticState) SetSolar(status *output.SolarStatus) {
	if provider == nil {
		return
	}
	provider.mu.Lock()
	defer provider.mu.Unlock()
	provider.solar = cloneSolarStatus(status)
}

func (provider *retainedSemanticState) SetSolarFromCache(status *output.SolarStatus) {
	provider.SetSolar(status)
}

func (provider *retainedSemanticState) SetCylinders(cylinders []output.CylinderStatus) {
	if provider == nil {
		return
	}
	provider.mu.Lock()
	defer provider.mu.Unlock()
	if cylinders == nil {
		provider.cylinders = nil
		return
	}
	provider.cylinders = cloneCylinderStatuses(cylinders)
}

func (provider *retainedSemanticState) SetCylindersFromCache(cylinders []output.CylinderStatus) {
	provider.SetCylinders(cylinders)
}
