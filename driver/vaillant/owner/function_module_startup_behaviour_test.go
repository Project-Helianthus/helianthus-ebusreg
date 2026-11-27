package owner

import (
	"context"
	"errors"
	"fmt"
	"github.com/Project-Helianthus/helianthus-ebusgo/engine"
	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/output"
	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
	"testing"
	"time"
)

func TestVaillantOwnerFunctionModuleStartupMaxSlotFunctionModuleModelConvergesWithinBoundedIdentityBudget(t *testing.T) {
	poller, provider := newFunctionModuleDeadlineRecoveryTestPoller(nil)
	config := uint16(2)
	poller.system.ModuleConfigurationFunctionModuleController = &config
	seedAddresses := []byte{0x15, 0x16, 0x17, 0x18}
	for instance, address := range seedAddresses {
		poller.reg.Register(registry.DeviceInfo{
			Address:      address,
			Manufacturer: "Vaillant",
			DeviceID:     fmt.Sprintf("BASV%d", instance),
		})
	}
	highRegulatorKey := radioDeviceKey{Group: remoteRegulators.group, Instance: 0x03}
	if seed := poller.registryRadioDeviceSeeds()[highRegulatorKey]; seed == nil {
		t.Fatal("high-slot regulator registry seed fixture is missing")
	}

	recoveryCtx, cancelRecovery := context.WithCancel(context.Background())
	defer cancelRecovery()
	selectorReads := 0
	budgetExhausted := false
	poller.sendFrameFn = func(ctx context.Context, frame protocol.Frame) (*protocol.Frame, error) {
		return functionModuleHighSlotBudgetedResponse(ctx, frame, &selectorReads, &budgetExhausted, cancelRecovery)
	}

	poller.refreshRadioDevicesStartup(recoveryCtx)
	poller.refreshFunctionModuleSemanticStartup(recoveryCtx)

	functionModuleKey := radioDeviceKey{Group: remoteFunctionalModules.group, Instance: semanticStartupSlotFullMaxInstance}
	poller.mu.Lock()
	functionModuleSnapshot := cloneRadioSnapshot(poller.radioDevices[functionModuleKey])
	retainedRegulator := cloneRadioSnapshot(poller.radioDevices[highRegulatorKey])
	poller.mu.Unlock()
	verdict := provider.FunctionModuleInterpretation()
	if budgetExhausted || selectorReads > functionModuleHighSlotSelectorBudget || functionModuleSnapshot == nil ||
		functionModuleSnapshot.DeviceClassAddress == nil || *functionModuleSnapshot.DeviceClassAddress != circuitManagingDeviceFunctionModuleControllerAddress ||
		retainedRegulator == nil || retainedRegulator.SlotMode != "registry" ||
		verdict.Mode != output.FunctionModuleSemanticModeInterpreted || verdict.DegradedReason != "" || verdict.EvidenceRevision == "" {
		t.Fatalf("max-slot startup reads=%d exhausted=%t fm5=%#v retained_non_fm=%#v verdict=%#v; want slot 0x0A INTERPRETED within one bounded path with deferred non-FM seed retained", selectorReads, budgetExhausted, functionModuleSnapshot, retainedRegulator, verdict)
	}
}

func TestVaillantOwnerEnqueueAddressIdentityProbeGenericFunctionalModuleDoesNotPreemptLaterFunctionModuleModel(t *testing.T) {
	originalScanDirected := registryScanDirectedFn
	t.Cleanup(func() { registryScanDirectedFn = originalScanDirected })

	taskScheduler := newSemanticTaskScheduler()
	poller, provider := newFunctionModuleDeadlineRecoveryTestPoller(taskScheduler)
	classes := map[byte]byte{
		0x00: 0x2A,
		0x04: circuitManagingDeviceFunctionModuleControllerAddress,
	}
	registryScanDirectedFn = func(_ context.Context, _ registry.ScanBus, reg *registry.DeviceRegistry, _ byte, _ []byte) ([]registry.DeviceEntry, error) {
		entry := reg.Register(registry.DeviceInfo{
			Address:      circuitManagingDeviceFunctionModuleControllerAddress,
			Manufacturer: "Vaillant",
			DeviceID:     circuitManagingDeviceFunctionModuleControllerID,
		})
		return []registry.DeviceEntry{entry}, nil
	}
	poller.sendFrameFn = func(_ context.Context, frame protocol.Frame) (*protocol.Frame, error) {
		return functionModuleSemanticTestResponseForClasses(frame, classes)
	}

	poller.EnqueueAddressIdentityProbe(circuitManagingDeviceFunctionModuleControllerAddress)
	probeTask := taskScheduler.nextTask(context.Background())
	if probeTask == nil {
		t.Fatal("late identity probe task was not scheduled")
	}
	probeTask.run(context.Background())
	taskScheduler.taskDone(probeTask)

	genericKey := deviceSlotKey{Group: remoteFunctionalModules.group, Instance: 0x00}
	realKey := deviceSlotKey{Group: remoteFunctionalModules.group, Instance: 0x04}
	poller.mu.Lock()
	genericCached := poller.deviceSlotCache[genericKey]
	realCached := poller.deviceSlotCache[realKey]
	genericSnapshot := poller.radioDevices[radioDeviceKey(genericKey)]
	realSnapshot := poller.radioDevices[radioDeviceKey(realKey)]
	poller.mu.Unlock()
	verdict := provider.FunctionModuleInterpretation()
	if genericCached || genericSnapshot != nil || !realCached || realSnapshot == nil || realSnapshot.DeviceClassAddress == nil || *realSnapshot.DeviceClassAddress != circuitManagingDeviceFunctionModuleControllerAddress ||
		verdict.Mode != output.FunctionModuleSemanticModeInterpreted || verdict.DegradedReason != "" || verdict.EvidenceRevision == "" {
		t.Fatalf("generic-before-VR71 generic_cached=%t generic_snapshot=%#v real_cached=%t real_snapshot=%#v verdict=%#v; want generic ignored, later slot 4 VR71 merged, and healthy INTERPRETED", genericCached, genericSnapshot, realCached, realSnapshot, verdict)
	}
}

func TestVaillantOwnerEnqueueAddressIdentityProbeMaxSlotFunctionModuleModelConvergesWithinBoundedRecoveryBudget(t *testing.T) {
	originalScanDirected := registryScanDirectedFn
	t.Cleanup(func() { registryScanDirectedFn = originalScanDirected })

	taskScheduler := newSemanticTaskScheduler()
	poller, provider := newFunctionModuleDeadlineRecoveryTestPoller(taskScheduler)
	connected := true
	nonFMClass := uint8(0x18)
	highRegulatorKey := radioDeviceKey{Group: remoteRegulators.group, Instance: 0x03}
	poller.radioDevices[highRegulatorKey] = &vaillantRadioDeviceSnapshot{
		Group:              highRegulatorKey.Group,
		Instance:           highRegulatorKey.Instance,
		SlotMode:           "registry",
		DeviceConnected:    &connected,
		DeviceClassAddress: &nonFMClass,
	}

	scanSeamCalled := false
	registryScanDirectedFn = func(_ context.Context, _ registry.ScanBus, reg *registry.DeviceRegistry, _ byte, _ []byte) ([]registry.DeviceEntry, error) {
		scanSeamCalled = true
		entry := reg.Register(registry.DeviceInfo{
			Address:      circuitManagingDeviceFunctionModuleControllerAddress,
			Manufacturer: "Vaillant",
			DeviceID:     circuitManagingDeviceFunctionModuleControllerID,
		})
		return []registry.DeviceEntry{entry}, nil
	}

	recoveryCtx, cancelRecovery := context.WithCancel(context.Background())
	defer cancelRecovery()
	selectorReads := 0
	budgetExhausted := false
	poller.sendFrameFn = func(ctx context.Context, frame protocol.Frame) (*protocol.Frame, error) {
		return functionModuleHighSlotBudgetedResponse(ctx, frame, &selectorReads, &budgetExhausted, cancelRecovery)
	}

	poller.EnqueueAddressIdentityProbe(circuitManagingDeviceFunctionModuleControllerAddress)
	probeTask := taskScheduler.nextTask(context.Background())
	if probeTask == nil {
		t.Fatal("late identity probe task was not scheduled")
	}
	probeTask.run(recoveryCtx)
	taskScheduler.taskDone(probeTask)

	functionModuleDeviceKey := deviceSlotKey{Group: remoteFunctionalModules.group, Instance: semanticStartupSlotFullMaxInstance}
	poller.mu.Lock()
	functionModuleCached := poller.deviceSlotCache[functionModuleDeviceKey]
	functionModuleSnapshot := cloneRadioSnapshot(poller.radioDevices[radioDeviceKey(functionModuleDeviceKey)])
	retainedRegulator := cloneRadioSnapshot(poller.radioDevices[highRegulatorKey])
	poller.mu.Unlock()
	publishedRegulator := false
	for _, device := range provider.RadioDevices() {
		if device.Group == int(highRegulatorKey.Group) && device.Instance == int(highRegulatorKey.Instance) {
			publishedRegulator = true
			break
		}
	}
	verdict := provider.FunctionModuleInterpretation()
	if !scanSeamCalled || budgetExhausted || selectorReads > functionModuleHighSlotSelectorBudget || !functionModuleCached || functionModuleSnapshot == nil ||
		functionModuleSnapshot.DeviceClassAddress == nil || *functionModuleSnapshot.DeviceClassAddress != circuitManagingDeviceFunctionModuleControllerAddress ||
		retainedRegulator == nil || !publishedRegulator || verdict.Mode != output.FunctionModuleSemanticModeInterpreted ||
		verdict.DegradedReason != "" || verdict.EvidenceRevision == "" {
		t.Fatalf("max-slot late recovery seam=%t reads=%d exhausted=%t cached=%t fm5=%#v retained_non_fm=%#v published_non_fm=%t verdict=%#v; want one bounded directed recovery to publish slot 0x0A INTERPRETED without losing non-FM state", scanSeamCalled, selectorReads, budgetExhausted, functionModuleCached, functionModuleSnapshot, retainedRegulator, publishedRegulator, verdict)
	}
}

func TestVaillantOwnerEnqueueAddressIdentityProbeMaxSlotFunctionModuleModelRecoveryReadsOnlyConfigurationGate(t *testing.T) {
	originalScanDirected := registryScanDirectedFn
	t.Cleanup(func() { registryScanDirectedFn = originalScanDirected })

	taskScheduler := newSemanticTaskScheduler()
	poller, provider := newFunctionModuleDeadlineRecoveryTestPoller(taskScheduler)
	scanSeamCalled := false
	registryScanDirectedFn = func(_ context.Context, _ registry.ScanBus, reg *registry.DeviceRegistry, _ byte, _ []byte) ([]registry.DeviceEntry, error) {
		scanSeamCalled = true
		entry := reg.Register(registry.DeviceInfo{
			Address:      circuitManagingDeviceFunctionModuleControllerAddress,
			Manufacturer: "Vaillant",
			DeviceID:     circuitManagingDeviceFunctionModuleControllerID,
		})
		return []registry.DeviceEntry{entry}, nil
	}

	recoveryCtx, cancelRecovery := context.WithCancel(context.Background())
	defer cancelRecovery()
	systemSelectors := make([]uint16, 0, 1)
	selectorReads := 0
	budgetExhausted := false
	poller.sendFrameFn = func(ctx context.Context, frame protocol.Frame) (*protocol.Frame, error) {
		if len(frame.Data) != 6 {
			return nil, errors.New("invalid B524 selector")
		}
		group := frame.Data[2]
		addr := uint16(frame.Data[4]) | uint16(frame.Data[5])<<8
		if group == localRegulator.group {
			systemSelectors = append(systemSelectors, addr)
			if addr != systemModuleConfigurationFunctionModule {
				cancelRecovery()
				return nil, errors.New("late FM5 recovery attempted non-gate system selector")
			}
		}
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return functionModuleHighSlotBudgetedResponse(ctx, frame, &selectorReads, &budgetExhausted, cancelRecovery)
	}

	poller.EnqueueAddressIdentityProbe(circuitManagingDeviceFunctionModuleControllerAddress)
	probeTask := taskScheduler.nextTask(context.Background())
	if probeTask == nil {
		t.Fatal("late identity probe task was not scheduled")
	}
	probeTask.run(recoveryCtx)
	taskScheduler.taskDone(probeTask)

	functionModuleKey := radioDeviceKey{Group: remoteFunctionalModules.group, Instance: semanticStartupSlotFullMaxInstance}
	poller.mu.Lock()
	functionModuleSnapshot := cloneRadioSnapshot(poller.radioDevices[functionModuleKey])
	poller.mu.Unlock()
	verdict := provider.FunctionModuleInterpretation()
	onlyGateRead := len(systemSelectors) == 1 && systemSelectors[0] == systemModuleConfigurationFunctionModule
	if !scanSeamCalled || !onlyGateRead || budgetExhausted || selectorReads > functionModuleHighSlotSelectorBudget ||
		functionModuleSnapshot == nil || functionModuleSnapshot.DeviceClassAddress == nil || *functionModuleSnapshot.DeviceClassAddress != circuitManagingDeviceFunctionModuleControllerAddress ||
		verdict.Mode != output.FunctionModuleSemanticModeInterpreted || verdict.DegradedReason != "" || verdict.EvidenceRevision == "" {
		t.Fatalf("gate-only max-slot recovery seam=%t system_selectors=%#v reads=%d exhausted=%t fm5=%#v verdict=%#v; want only module_configuration_vr71 before bounded slot 0x0A INTERPRETED recovery", scanSeamCalled, systemSelectors, selectorReads, budgetExhausted, functionModuleSnapshot, verdict)
	}
}

func TestVaillantOwnerRefreshRadioDevicesCompleteNegativeFunctionModuleScanSupersedesRetainedRegistryIdentity(t *testing.T) {
	now := time.Date(2026, 8, 15, 15, 0, 0, 0, time.UTC)
	config := uint16(2)
	classAddress := uint8(circuitManagingDeviceFunctionModuleControllerAddress)
	collector := 61.5
	temperature := 47.25
	previous := vaillantFunctionModuleInterpretation{
		Mode:             vaillantFunctionModuleModeInterpreted,
		EvidenceRevision: "fm5-g7-a41",
	}
	provider := newRetainedSemanticState()
	provider.SetFunctionModuleInterpretation(publishedFunctionModuleInterpretation(previous))
	provider.SetSolar(&output.SolarStatus{CollectorTemperatureC: &collector})
	provider.SetCylinders([]output.CylinderStatus{{Index: 0, TemperatureC: &temperature}})
	reg := registry.NewDeviceRegistry(nil)
	reg.Register(registry.DeviceInfo{
		Address:      circuitManagingDeviceFunctionModuleControllerAddress,
		Manufacturer: "Vaillant",
		DeviceID:     circuitManagingDeviceFunctionModuleControllerID,
	})
	poller := &vaillantSemanticPoller{
		scheduler:                        engine.NewSemanticReadScheduler(),
		reg:                              reg,
		publication:                      newVaillantSemanticPublication(provider, nil, nil),
		source:                           0x7F,
		controller:                       0x15,
		requestTimeout:                   50 * time.Millisecond,
		system:                           &vaillantSystemSnapshot{Controller: 0x15, ModuleConfigurationFunctionModuleController: &config},
		radioDevices:                     map[radioDeviceKey]*vaillantRadioDeviceSnapshot{{Group: remoteFunctionalModules.group}: {DeviceClassAddress: &classAddress}},
		deviceSlotCache:                  make(map[deviceSlotKey]bool),
		functionModuleMode:               vaillantFunctionModuleModeInterpreted,
		functionModuleInterpretation:     previous,
		functionModuleEvidenceRevision:   41,
		functionModuleEvidenceGeneration: 7,
		functionModuleIdentityObservedAt: now,
		functionModuleEvidenceTTL:        5 * time.Minute,
		solar:                            &vaillantSolarSnapshot{CollectorTemperatureC: &collector},
		solarCylinders:                   map[byte]*vaillantCylinderSnapshot{0: {Instance: 0, TemperatureC: &temperature}},
		nowFn:                            func() time.Time { return now },
	}
	poller.sendFrameFn = func(_ context.Context, frame protocol.Frame) (*protocol.Frame, error) {
		if len(frame.Data) != 6 {
			return nil, errors.New("invalid B524 selector")
		}
		group := frame.Data[2]
		addr := uint16(frame.Data[4]) | uint16(frame.Data[5])<<8
		if group == remoteFunctionalModules.group {
			switch addr {
			case deviceSlotConnected:
				return testB524ResponseForSelectorPayload(frame.Data, 0x00), nil
			case deviceSlotClassAddress:
				return testB524ResponseForSelectorPayload(frame.Data, 0xFF), nil
			case deviceSlotFirmware:
				return testB524ResponseForSelectorPayload(frame.Data, 0xFF, 0xFF, 0xFF), nil
			case deviceSlotHardwareIdentifier:
				return testB524ResponseForSelectorPayload(frame.Data, 0xFF, 0xFF), nil
			}
		}
		return testB524ResponseForSelectorPayload(frame.Data, 0x00, 0x00, 0x00, 0x00), nil
	}

	poller.refreshRadioDevices(context.Background())

	got := provider.FunctionModuleInterpretation()
	if got.Mode != output.FunctionModuleSemanticModeAbsent || got.DegradedReason != "" {
		t.Errorf("complete negative scan verdict = %#v; want healthy ABSENT despite retained registry entry", got)
	}
	if got.EvidenceRevision == "" || got.EvidenceRevision == previous.EvidenceRevision {
		t.Errorf("complete negative scan revision = %q; want advancement from %q", got.EvidenceRevision, previous.EvidenceRevision)
	}
	if gotMode := provider.FunctionModuleSemanticMode(); gotMode != output.FunctionModuleSemanticModeAbsent {
		t.Errorf("complete negative scan legacy scalar = %s; want ABSENT", gotMode)
	}
	if solar := provider.Solar(); solar == nil || solar.CollectorTemperatureC != nil {
		t.Errorf("complete negative scan solar = %#v; want cleared structural plane", solar)
	}
	if cylinders := provider.Cylinders(); len(cylinders) != 0 {
		t.Errorf("complete negative scan cylinders = %#v; want cleared structural plane", cylinders)
	}
}

func TestVaillantOwnerCommitFunctionModuleAcquisitionSerializesWriterAfterCommitWithoutDeadlock(t *testing.T) {
	classAddress := uint8(circuitManagingDeviceFunctionModuleControllerAddress)
	config := uint16(2)
	reg := registry.NewDeviceRegistry(nil)
	reg.Register(registry.DeviceInfo{
		Address:      circuitManagingDeviceFunctionModuleControllerAddress,
		Manufacturer: "Vaillant",
		DeviceID:     circuitManagingDeviceFunctionModuleControllerID,
	})
	poller := &vaillantSemanticPoller{
		reg:            reg,
		controller:     0x15,
		system:         &vaillantSystemSnapshot{Controller: 0x15, ModuleConfigurationFunctionModuleController: &config},
		radioDevices:   map[radioDeviceKey]*vaillantRadioDeviceSnapshot{{Group: remoteFunctionalModules.group}: {DeviceClassAddress: &classAddress}},
		solarCylinders: make(map[byte]*vaillantCylinderSnapshot),
		nowFn:          time.Now,
	}
	captured := poller.captureFunctionModuleEvidence()
	if !captured.registryCoherent {
		t.Fatal("initial registry capture is incoherent")
	}

	readLocked := make(chan struct{})
	poller.withFunctionModuleObservationGeneration = func(fn func(uint64)) bool {
		return reg.WithObservationGeneration(func(current uint64) {
			close(readLocked)
			fn(current)
		})
	}
	poller.mu.Lock()
	commitDone := make(chan vaillantFunctionModuleInterpretation, 1)
	go func() {
		commitDone <- poller.commitFunctionModuleAcquisition(captured, vaillantFunctionModuleInterpretation{
			Mode:             vaillantFunctionModuleModeInterpreted,
			EvidenceRevision: "fm5-g0-a1",
		}, &vaillantSolarSnapshot{}, map[byte]*vaillantCylinderSnapshot{})
	}()
	<-readLocked

	writerDone := make(chan struct{})
	go func() {
		reg.Register(registry.DeviceInfo{
			Address:      circuitManagingDeviceFunctionModuleControllerAddress + 1,
			Manufacturer: "Vaillant",
			DeviceID:     "VR71-WRITER-AFTER",
		})
		close(writerDone)
	}()
	poller.mu.Unlock()

	deadline := time.NewTimer(time.Second)
	defer deadline.Stop()
	var verdict vaillantFunctionModuleInterpretation
	select {
	case verdict = <-commitDone:
	case <-deadline.C:
		t.Fatal("FM5 commit deadlocked with queued registry writer")
	}
	select {
	case <-writerDone:
	case <-deadline.C:
		t.Fatal("registry writer did not complete after FM5 read critical section")
	}
	if verdict.Mode != vaillantFunctionModuleModeInterpreted || verdict.DegradedReason != "" {
		t.Fatalf("serialized commit verdict = %#v; want INTERPRETED", verdict)
	}
	_, afterGeneration, coherent := poller.captureFunctionModuleRegistryEvidence()
	if !coherent || afterGeneration != captured.registryGeneration+1 {
		t.Fatalf("writer-after generation = %d coherent=%t; want %d true", afterGeneration, coherent, captured.registryGeneration+1)
	}
}
