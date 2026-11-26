package owner

import (
	"context"
	"fmt"
	"log"
	"time"
)

func (p *vaillantSemanticPoller) boilerStatusTierSchedules() []boilerStatusTierSchedule {
	if p == nil {
		return nil
	}
	return []boilerStatusTierSchedule{
		{tier: boilerStatusTierFast, interval: p.boilerFastInterval, priority: semanticTaskPriorityHigh},
		{tier: boilerStatusTierMedium, interval: p.boilerMediumInterval, priority: semanticTaskPriorityMedium},
		{tier: boilerStatusTierSlow, interval: p.boilerSlowInterval, priority: semanticTaskPriorityLow},
	}
}

func (p *vaillantSemanticPoller) boilerStatusTierTask(tier boilerStatusTier) func(context.Context) {
	return func(ctx context.Context) {
		if tier == boilerStatusTierFast {
			p.refreshBoilerStatus(ctx)
			return
		}
		p.refreshBoilerStatusTier(ctx, tier)
	}
}

func boilerStatusTierTaskKey(tier boilerStatusTier) semanticTaskKey {
	switch tier {
	case boilerStatusTierFast:
		return semanticTaskRefreshBoilerFast
	case boilerStatusTierMedium:
		return semanticTaskRefreshBoilerMedium
	case boilerStatusTierSlow:
		return semanticTaskRefreshBoilerSlow
	default:
		return semanticTaskKey("")
	}
}

func (p *vaillantSemanticPoller) enqueueBoilerStatusPriming(ctx context.Context) {
	if p == nil {
		return
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if p.tasks == nil {
		for _, schedule := range p.boilerStatusTierSchedules() {
			p.boilerStatusTierTask(schedule.tier)(ctx)
		}
		return
	}
	for _, schedule := range p.boilerStatusTierSchedules() {
		p.enqueueTask(boilerStatusTierTaskKey(schedule.tier), schedule.priority, p.boilerStatusTierTask(schedule.tier))
	}
}

func (p *vaillantSemanticPoller) enqueueControllerSemanticPriming(ctx context.Context) {
	if p == nil {
		return
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if p.tasks == nil {
		p.refreshConfig(ctx)
		p.refreshCircuits(ctx)
		p.refreshSystem(ctx)
		p.refreshRadioDevices(ctx)
		p.refreshEnergy(ctx)
		return
	}
	p.enqueueTask(semanticTaskRefreshConfig, semanticTaskPriorityHigh, p.refreshConfig)
	p.enqueueTask(semanticTaskRefreshCircuits, semanticTaskPriorityMedium, p.refreshCircuits)
	p.enqueueTask(semanticTaskRefreshSystem, semanticTaskPriorityMedium, p.refreshSystem)
	p.enqueueTask(semanticTaskRefreshRadioDevices, semanticTaskPriorityMedium, p.refreshRadioDevices)
	p.enqueueTask(semanticTaskRefreshEnergy, semanticTaskPriorityMedium, p.refreshEnergy)
}

func (p *vaillantSemanticPoller) refreshStartupSemanticPlanes(ctx context.Context) {
	if p == nil {
		return
	}
	if ctx == nil {
		ctx = context.Background()
	}

	p.mu.Lock()
	controller := p.controller
	p.mu.Unlock()
	if controller == 0 {
		p.refreshDiscovery(ctx)
		p.mu.Lock()
		controller = p.controller
		p.mu.Unlock()
	}
	if controller == 0 {
		return
	}
	if !p.publishing() {
		return
	}

	p.publishStartupSchedules()
	p.publishStartupPlaneDefaults()
	primingCtx := ctx
	cancel := func() {}
	if semanticStartupPrimingBudget > 0 {
		primingCtx, cancel = context.WithTimeout(ctx, semanticStartupPrimingBudget)
	}
	defer cancel()

	attempts := 0
	for {
		attempts++
		status := p.startupLevelOnePrimingStatus()
		if !status.dhw {
			p.refreshDHWStartupUntilReady(primingCtx, 1)
			status = p.startupLevelOnePrimingStatus()
		}
		if !status.zones {
			p.refreshZoneDiscovery(primingCtx, true)
			status = p.startupLevelOnePrimingStatus()
		}
		if !status.dhw {
			p.refreshDHWStartupUntilReady(primingCtx, 1)
		}
		if !status.system {
			p.refreshSystemStartup(primingCtx)
		}
		if !status.boilerStatus {
			p.refreshBoilerStatusStartup(primingCtx)
		}
		status = p.startupLevelOnePrimingStatus()
		if !status.circuits {
			p.refreshCircuitsStartup(primingCtx)
		}
		if !status.radioDevices {
			p.refreshRadioDevicesStartup(primingCtx)
		}
		status = p.startupLevelOnePrimingStatus()
		if !status.functionModuleSatisfied && status.functionModuleEvidence && !status.functionModuleGateKnown {
			p.refreshSystemStartup(primingCtx)
			status = p.startupLevelOnePrimingStatus()
		}
		if status.system && (!status.functionModuleSatisfied || !status.solar || !status.cylinders) {
			p.refreshFunctionModuleSemanticStartup(primingCtx)
		}
		status = p.startupLevelOnePrimingStatus()
		if status.ready() {
			log.Printf("semantic_startup_l1_priming result=ready attempts=%d %s", attempts, status.String())
			return
		}
		if semanticStartupPrimingBudget <= 0 {
			log.Printf("semantic_startup_l1_priming result=incomplete attempts=%d %s", attempts, status.String())
			return
		}

		delay := semanticStartupPrimingRetryDelay
		if delay <= 0 {
			delay = 50 * time.Millisecond
		}
		timer := time.NewTimer(delay)
		select {
		case <-primingCtx.Done():
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
			log.Printf("semantic_startup_l1_priming result=deadline attempts=%d %s", attempts, status.String())
			return
		case <-timer.C:
		}
	}
}

// publishStartupPlaneDefaults has the publication owner publish the empty
// circuits, solar and cylinders planes that are not published yet.
func (p *vaillantSemanticPoller) publishStartupPlaneDefaults() {
	if !p.publishing() {
		return
	}
	p.publication.publishStartupPlaneDefaults()
}

func (p *vaillantSemanticPoller) refreshStartupCriticalSemanticPlanes(ctx context.Context) {
	if !p.publishing() {
		return
	}

	p.publishStartupSchedules()
	status := p.startupLevelOnePrimingStatus()
	if !status.dhw {
		p.refreshDHWStartupUntilReady(ctx, semanticStartupCriticalDHWAttempts)
	}
	status = p.startupLevelOnePrimingStatus()
	if !status.system {
		p.refreshSystemStartup(ctx)
	}
	status = p.startupLevelOnePrimingStatus()
	if !status.boilerStatus {
		p.refreshBoilerStatusStartup(ctx)
	}
}

func (p *vaillantSemanticPoller) refreshDHWStartupUntilReady(ctx context.Context, attempts int) bool {
	if p == nil {
		return false
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if attempts <= 0 {
		attempts = 1
	}
	for attempt := 0; attempt < attempts; attempt++ {
		if p.startupLevelOnePrimingStatus().dhw {
			return true
		}
		p.refreshDHWStartup(ctx)
		if p.startupLevelOnePrimingStatus().dhw {
			return true
		}
		if attempt == attempts-1 {
			break
		}
		delay := semanticStartupPrimingRetryDelay
		if delay <= 0 {
			delay = 50 * time.Millisecond
		}
		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
			return false
		case <-timer.C:
		}
	}
	return p.startupLevelOnePrimingStatus().dhw
}

type startupLevelOnePrimingStatus struct {
	zones                   bool
	dhw                     bool
	circuits                bool
	system                  bool
	radioDevices            bool
	functionModuleGateKnown bool
	functionModuleEvidence  bool
	functionModuleRequired  bool
	functionModuleSatisfied bool
	solar                   bool
	cylinders               bool
	boilerStatus            bool
}

func (status startupLevelOnePrimingStatus) ready() bool {
	return status.zones &&
		status.dhw &&
		status.circuits &&
		status.system &&
		status.radioDevices &&
		status.functionModuleSatisfied &&
		status.boilerStatus
}

func (status startupLevelOnePrimingStatus) String() string {
	return fmt.Sprintf(
		"zones=%t dhw=%t circuits=%t system=%t radio_devices=%t fm5_gate_known=%t fm5_evidence=%t fm5_required=%t fm5_satisfied=%t solar=%t cylinders=%t boiler_status=%t",
		status.zones,
		status.dhw,
		status.circuits,
		status.system,
		status.radioDevices,
		status.functionModuleGateKnown,
		status.functionModuleEvidence,
		status.functionModuleRequired,
		status.functionModuleSatisfied,
		status.solar,
		status.cylinders,
		status.boilerStatus,
	)
}

func (p *vaillantSemanticPoller) startupLevelOnePrimingStatus() startupLevelOnePrimingStatus {
	if !p.publishing() {
		return startupLevelOnePrimingStatus{}
	}
	p.mu.Lock()
	moduleConfig := (*uint16)(nil)
	if p.system != nil {
		moduleConfig = cloneUint16Ptr(p.system.ModuleConfigurationFunctionModuleController)
	}
	functionModuleGateKnown := p.system != nil && p.system.ModuleConfigurationFunctionModuleController != nil
	radioSnapshots := make([]*vaillantRadioDeviceSnapshot, 0, len(p.radioDevices))
	for _, snapshot := range p.radioDevices {
		if snapshot != nil {
			radioSnapshots = append(radioSnapshots, cloneRadioSnapshot(snapshot))
		}
	}
	radioProbed := p.startupRadioDevicesProbed
	functionModuleRegistryEvidenceIgnored := p.functionModuleRegistryEvidenceIgnored
	p.mu.Unlock()
	liveFunctionModuleEvidence := hasFunctionModuleEvidenceFromRadioSnapshots(radioSnapshots)
	functionModuleEvidence := liveFunctionModuleEvidence || (!functionModuleRegistryEvidenceIgnored && p.hasFunctionModuleRegistryEvidence())
	functionModuleRequired := functionModuleEvidence && moduleConfig != nil && *moduleConfig <= 2
	planes := p.publication.publishedPlanes()
	zonesReady := planes.zones
	solarReady := planes.solar
	cylindersPublished := planes.cylinders
	interpretedCylindersReady := planes.cylinderCount > 0
	functionModuleNonInterpretedPublished := functionModuleEvidence &&
		p.publication.functionModulePublishedAsGPIOOnly() &&
		solarReady &&
		cylindersPublished
	functionModuleSatisfied := !functionModuleEvidence ||
		functionModuleNonInterpretedPublished ||
		(functionModuleGateKnown && !functionModuleRequired) ||
		(functionModuleRequired && solarReady && interpretedCylindersReady)
	return startupLevelOnePrimingStatus{
		zones:                   zonesReady,
		dhw:                     planes.dhw,
		circuits:                planes.circuits,
		system:                  planes.system,
		radioDevices:            radioProbed || planes.radioDeviceCount > 0,
		functionModuleGateKnown: functionModuleGateKnown,
		functionModuleEvidence:  functionModuleEvidence,
		functionModuleRequired:  functionModuleRequired,
		functionModuleSatisfied: functionModuleSatisfied,
		solar:                   solarReady,
		cylinders:               cylindersPublished,
		boilerStatus:            planes.boilerStatus,
	}
}

func (p *vaillantSemanticPoller) startupProbeContext(ctx context.Context) (context.Context, context.CancelFunc) {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithTimeout(ctx, semanticStartupProbeTimeout)
}

func (p *vaillantSemanticPoller) readB524Startup(ctx context.Context, opcode, group, instance byte, addr uint16) ([]byte, bool) {
	probeCtx, cancel := p.startupProbeContext(ctx)
	defer cancel()
	return p.readB524ValueLive(probeCtx, opcode, group, instance, addr)
}

func (p *vaillantSemanticPoller) readB524Uint16Startup(ctx context.Context, opcode, group, instance byte, addr uint16) (*uint16, bool) {
	raw, ok := p.readB524Startup(ctx, opcode, group, instance, addr)
	if !ok {
		return nil, false
	}
	value, ok := decodeB524Uint16(raw)
	if !ok {
		return nil, false
	}
	return &value, true
}

// publishStartupSchedules has the publication owner publish the empty
// schedule placeholder when no schedules are published yet.
func (p *vaillantSemanticPoller) publishStartupSchedules() {
	if !p.publishing() {
		return
	}
	p.publication.publishStartupSchedulePlaceholder()
}

func (p *vaillantSemanticPoller) refreshDHWStartup(ctx context.Context) {
	if p == nil {
		return
	}

	attempted := make(semanticFieldSet)
	status := &vaillantDhwSnapshot{}
	readAny := false
	if raw, ok := p.readB524Startup(ctx, localDHW.opcode, localDHW.group, dhwInstance, dhwCurrentTemp); ok {
		if value := decodeB524Float32FromRaw(raw); value != nil {
			status.CurrentTempC = value
			attempted[dhwFieldCurrentTempC] = struct{}{}
			readAny = true
		}
	}
	opModeRaw, opModeOK := p.readB524Uint16Startup(ctx, localDHW.opcode, localDHW.group, dhwInstance, dhwOperationMode)
	if opModeOK {
		attempted[dhwFieldOperatingMode] = struct{}{}
		attempted[dhwFieldPreset] = struct{}{}
		attempted[dhwFieldDhwOperationModeRaw] = struct{}{}
		readAny = true
	}
	if opModeOK {
		status.OperatingMode, status.Preset = deriveDhwModeAndPreset(opModeRaw, nil)
		status.ConfigurationDHWOperationMode = formatUintToken(*opModeRaw)
	}

	p.mu.Lock()
	source := semanticSnapshotSourceCache
	if readAny {
		if p.dhw == nil {
			p.dhw = &vaillantDhwSnapshot{}
		}
		mergeDhwSnapshotFields(p.dhw, status, semanticSnapshotSourceLive, attempted)
		p.markDHWUpdatedNowLocked()
		source = semanticSnapshotSourceLive
	}
	hasSnapshot := p.dhw != nil
	p.mu.Unlock()
	if hasSnapshot {
		p.publishDHW(source)
	}
}

func (p *vaillantSemanticPoller) refreshCircuitsStartup(ctx context.Context) {
	if p == nil {
		return
	}

	p.mu.Lock()
	controller := p.controller
	p.mu.Unlock()
	instances := allCircuitRefreshInstances()
	updates := make(map[byte]*vaillantCircuitSnapshot)
	inactive := make(map[byte]struct{})
	for _, instance := range instances {
		circuitTypeRaw, ok := p.readB524Uint16Startup(ctx, localCircuits.opcode, localCircuits.group, instance, circuitTypeRegister)
		if !ok || circuitTypeRaw == nil {
			continue
		}
		switch *circuitTypeRaw {
		case 0x0000, 0x00FF, 0xFFFF:
			if snapshot, known := p.readDhwPseudoCircuitStartupEvidence(ctx, instance); snapshot != nil {
				snapshot.Controller = controller
				updates[instance] = snapshot
				continue
			} else if !known && instance == dhwPseudoCircuitInstance {
				continue
			}
			inactive[instance] = struct{}{}
			continue
		default:
			updates[instance] = &vaillantCircuitSnapshot{
				Instance:       instance,
				Active:         true,
				Controller:     controller,
				CircuitTypeRaw: cloneUint16Ptr(circuitTypeRaw),
			}
		}
	}
	if len(updates) == 0 && len(inactive) == 0 {
		p.mu.Lock()
		hasExisting := len(p.circuits) > 0
		p.mu.Unlock()
		if hasExisting {
			p.publishCircuits(semanticSnapshotSourceCache)
		}
		return
	}

	p.mu.Lock()
	if p.circuits == nil {
		p.circuits = make(map[byte]*vaillantCircuitSnapshot)
	}
	for instance := range inactive {
		delete(p.circuits, instance)
	}
	for instance, incoming := range updates {
		if incoming.CircuitStateRaw != nil {
			incoming.CircuitStateLiveAt = p.now()
		}
		if incoming.PumpStatusRaw != nil {
			incoming.PumpStatusLiveAt = p.now()
		}
		p.circuits[instance] = mergeCircuitSnapshotNonDestructive(p.circuits[instance], incoming)
	}
	p.lastCircuitFullScanAt = p.now()
	p.lastCircuitFullScanComplete = len(updates)+len(inactive) == len(instances)
	p.mu.Unlock()
	p.publishCircuits(semanticSnapshotSourceLive)
}
