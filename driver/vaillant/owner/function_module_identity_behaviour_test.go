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

func TestVaillantOwnerFunctionModulePositiveIdentityDoesNotRequireNegativeNamespaceCompleteness(t *testing.T) {
	now := time.Date(2026, 8, 15, 19, 0, 0, 0, time.UTC)
	config := uint16(2)
	newPoller := func(provider *retainedSemanticState, reg *registry.DeviceRegistry) *vaillantSemanticPoller {
		return &vaillantSemanticPoller{
			scheduler:                 engine.NewSemanticReadScheduler(),
			reg:                       reg,
			publication:               newVaillantSemanticPublication(provider, nil, nil),
			source:                    0x7F,
			controller:                0x15,
			requestTimeout:            50 * time.Millisecond,
			system:                    &vaillantSystemSnapshot{Controller: 0x15, ModuleConfigurationFunctionModuleController: &config},
			radioDevices:              make(map[radioDeviceKey]*vaillantRadioDeviceSnapshot),
			solarCylinders:            make(map[byte]*vaillantCylinderSnapshot),
			functionModuleEvidenceTTL: 5 * time.Minute,
			deviceSlotCache:           make(map[deviceSlotKey]bool),
			nowFn:                     func() time.Time { return now },
		}
	}
	assertInterpreted := func(t *testing.T, provider *retainedSemanticState) {
		t.Helper()
		got := provider.FunctionModuleInterpretation()
		if got.Mode != output.FunctionModuleSemanticModeInterpreted || got.DegradedReason != "" || got.EvidenceRevision == "" {
			t.Fatalf("positive partial-namespace verdict = %#v; want healthy INTERPRETED with revision", got)
		}
	}
	responder := func(positiveInstance byte, timeoutInstance *byte) func(context.Context, protocol.Frame) (*protocol.Frame, error) {
		return func(_ context.Context, frame protocol.Frame) (*protocol.Frame, error) {
			if len(frame.Data) != 6 {
				return nil, errors.New("invalid B524 selector")
			}
			group := frame.Data[2]
			instance := frame.Data[3]
			addr := uint16(frame.Data[4]) | uint16(frame.Data[5])<<8
			if group == remoteFunctionalModules.group {
				if timeoutInstance != nil && instance == *timeoutInstance && addr == deviceSlotConnected {
					return nil, errors.New("unrelated functional-module slot timeout")
				}
				switch addr {
				case deviceSlotConnected:
					return testB524ResponseForSelectorPayload(frame.Data, 0x00), nil
				case deviceSlotClassAddress:
					value := byte(0xFF)
					if instance == positiveInstance {
						value = circuitManagingDeviceFunctionModuleControllerAddress
					}
					return testB524ResponseForSelectorPayload(frame.Data, value), nil
				case deviceSlotFirmware:
					return testB524ResponseForSelectorPayload(frame.Data, 0xFF, 0xFF, 0xFF), nil
				case deviceSlotHardwareIdentifier:
					return testB524ResponseForSelectorPayload(frame.Data, 0xFF, 0xFF), nil
				}
			}
			return testB524ResponseForSelectorPayload(frame.Data, 0x00, 0x00, 0x00, 0x00), nil
		}
	}

	t.Run("startup registry-seeded low-slot positive", func(t *testing.T) {
		provider := newRetainedSemanticState()
		reg := registry.NewDeviceRegistry(nil)
		reg.Register(registry.DeviceInfo{Address: circuitManagingDeviceFunctionModuleControllerAddress, Manufacturer: "Vaillant", DeviceID: circuitManagingDeviceFunctionModuleControllerID})
		poller := newPoller(provider, reg)
		poller.sendFrameFn = responder(0x00, nil)

		poller.refreshRadioDevicesStartup(context.Background())
		poller.refreshFunctionModuleSemanticStartup(context.Background())
		assertInterpreted(t, provider)
	})

	t.Run("startup positive with unrelated slot timeout", func(t *testing.T) {
		provider := newRetainedSemanticState()
		poller := newPoller(provider, registry.NewDeviceRegistry(nil))
		timeoutInstance := byte(0x07)
		poller.sendFrameFn = responder(0x04, &timeoutInstance)

		poller.refreshRadioDevicesStartup(context.Background())
		poller.refreshFunctionModuleSemanticStartup(context.Background())
		assertInterpreted(t, provider)
	})

	t.Run("steady positive with unrelated slot timeout", func(t *testing.T) {
		provider := newRetainedSemanticState()
		poller := newPoller(provider, registry.NewDeviceRegistry(nil))
		timeoutInstance := byte(0x07)
		poller.sendFrameFn = responder(0x04, &timeoutInstance)

		poller.refreshRadioDevices(context.Background())
		assertInterpreted(t, provider)
	})
}

func TestVaillantOwnerFunctionModuleStartupUnseededLowSlotPositiveDoesNotLosePrimingBudget(t *testing.T) {
	now := time.Date(2026, 8, 15, 20, 0, 0, 0, time.UTC)
	config := uint16(2)
	const (
		collectorTemperature = 61.5
		cylinderTemperatureC = 47.25
	)
	provider := newRetainedSemanticState()
	poller := &vaillantSemanticPoller{
		scheduler:                 engine.NewSemanticReadScheduler(),
		reg:                       registry.NewDeviceRegistry(nil),
		publication:               newVaillantSemanticPublication(provider, nil, nil),
		source:                    0x7F,
		controller:                0x15,
		requestTimeout:            50 * time.Millisecond,
		system:                    &vaillantSystemSnapshot{Controller: 0x15, ModuleConfigurationFunctionModuleController: &config},
		radioDevices:              make(map[radioDeviceKey]*vaillantRadioDeviceSnapshot),
		solarCylinders:            make(map[byte]*vaillantCylinderSnapshot),
		functionModuleEvidenceTTL: 5 * time.Minute,
		deviceSlotCache:           make(map[deviceSlotKey]bool),
		nowFn:                     func() time.Time { return now },
	}

	startupCtx, cancelStartup := context.WithCancel(context.Background())
	defer cancelStartup()
	tailEntered := false
	blockedFamilyReads := 0
	poller.sendFrameFn = func(ctx context.Context, frame protocol.Frame) (*protocol.Frame, error) {
		if len(frame.Data) != 6 {
			return nil, errors.New("invalid B524 selector")
		}
		group := frame.Data[2]
		instance := frame.Data[3]
		addr := uint16(frame.Data[4]) | uint16(frame.Data[5])<<8

		if group == remoteFunctionalModules.group && instance >= semanticStartupSlotFastMaxInstance+1 {
			tailEntered = true
			cancelStartup()
		}
		if (group == localSolar.group || group == localCylinders.group) && ctx.Err() != nil {
			blockedFamilyReads++
			return nil, ctx.Err()
		}

		if group == remoteFunctionalModules.group {
			switch addr {
			case deviceSlotConnected:
				connected := byte(0x00)
				if instance == 0x01 {
					connected = 0x01
				}
				return testB524ResponseForSelectorPayload(frame.Data, connected), nil
			case deviceSlotClassAddress:
				classAddress := byte(0xFF)
				if instance == 0x01 {
					classAddress = circuitManagingDeviceFunctionModuleControllerAddress
				}
				return testB524ResponseForSelectorPayload(frame.Data, classAddress), nil
			case deviceSlotFirmware:
				firmware := []byte{0xFF, 0xFF, 0xFF}
				if instance == 0x01 {
					firmware = []byte{0x01, 0x00, 0x00}
				}
				return testB524ResponseForSelectorPayload(frame.Data, firmware...), nil
			case deviceSlotHardwareIdentifier:
				hardware := []byte{0xFF, 0xFF}
				if instance == 0x01 {
					hardware = []byte{0x71, 0x00}
				}
				return testB524ResponseForSelectorPayload(frame.Data, hardware...), nil
			}
		}
		if group == localSolar.group {
			switch addr {
			case solarCollectorTemp:
				return testB524ResponseForSelectorPayload(frame.Data, testB524Float32Payload(collectorTemperature)...), nil
			case solarEnabled, solarPumpActive:
				return testB524ResponseForSelectorPayload(frame.Data, 0x01), nil
			}
		}
		if group == localCylinders.group && addr == cylinderTemperature {
			return testB524ResponseForSelectorPayload(frame.Data, testB524Float32Payload(cylinderTemperatureC)...), nil
		}
		return testB524ResponseForSelectorPayload(frame.Data, 0x00, 0x00, 0x00, 0x00), nil
	}

	poller.refreshRadioDevicesStartup(startupCtx)
	poller.refreshFunctionModuleSemanticStartup(startupCtx)

	got := provider.FunctionModuleInterpretation()
	if got.Mode != output.FunctionModuleSemanticModeInterpreted || got.DegradedReason != "" || got.EvidenceRevision == "" {
		t.Fatalf("unseeded low-slot startup verdict = %#v; want healthy INTERPRETED with revision (tail_entered=%t blocked_family_reads=%d)", got, tailEntered, blockedFamilyReads)
	}
	if solar := provider.Solar(); solar == nil || solar.CollectorTemperatureC == nil || *solar.CollectorTemperatureC != collectorTemperature {
		t.Fatalf("unseeded low-slot startup solar = %#v; want collector temperature %v", solar, collectorTemperature)
	}
	if cylinders := provider.Cylinders(); len(cylinders) == 0 || cylinders[0].TemperatureC == nil || *cylinders[0].TemperatureC != cylinderTemperatureC {
		t.Fatalf("unseeded low-slot startup cylinders = %#v; want temperature %v", cylinders, cylinderTemperatureC)
	}
}

func TestVaillantOwnerEnqueueAddressIdentityProbeLateFunctionModuleSuccessTriggersBoundedSemanticReacquisition(t *testing.T) {
	originalScanDirected := registryScanDirectedFn
	t.Cleanup(func() { registryScanDirectedFn = originalScanDirected })

	now := time.Date(2026, 8, 15, 20, 30, 0, 0, time.UTC)
	provider := newRetainedSemanticState()
	taskScheduler := newSemanticTaskScheduler()
	poller := &vaillantSemanticPoller{
		scheduler:                 engine.NewSemanticReadScheduler(),
		tasks:                     taskScheduler,
		reg:                       registry.NewDeviceRegistry(nil),
		bus:                       &protocol.Bus{},
		publication:               newVaillantSemanticPublication(provider, nil, nil),
		source:                    0x7F,
		controller:                0x15,
		requestTimeout:            50 * time.Millisecond,
		system:                    &vaillantSystemSnapshot{Controller: 0x15},
		radioDevices:              make(map[radioDeviceKey]*vaillantRadioDeviceSnapshot),
		solarCylinders:            make(map[byte]*vaillantCylinderSnapshot),
		functionModuleEvidenceTTL: 5 * time.Minute,
		deviceSlotRediscoveryTTL:  30 * time.Minute,
		deviceSlotCache:           make(map[deviceSlotKey]bool),
		nowFn:                     func() time.Time { return now },
	}

	scanSeamCalled := false
	registryScanDirectedFn = func(_ context.Context, _ registry.ScanBus, reg *registry.DeviceRegistry, source byte, targets []byte) ([]registry.DeviceEntry, error) {
		scanSeamCalled = true
		if source != poller.source || len(targets) != 1 || targets[0] != circuitManagingDeviceFunctionModuleControllerAddress {
			return nil, errors.New("unexpected directed identity target")
		}
		entry := reg.Register(registry.DeviceInfo{
			Address:      circuitManagingDeviceFunctionModuleControllerAddress,
			Manufacturer: "Vaillant",
			DeviceID:     circuitManagingDeviceFunctionModuleControllerID,
		})
		return []registry.DeviceEntry{entry}, nil
	}

	poller.sendFrameFn = func(_ context.Context, frame protocol.Frame) (*protocol.Frame, error) {
		if len(frame.Data) != 6 {
			return nil, errors.New("invalid B524 selector")
		}
		group := frame.Data[2]
		instance := frame.Data[3]
		addr := uint16(frame.Data[4]) | uint16(frame.Data[5])<<8
		if group == localRegulator.group && addr == systemModuleConfigurationFunctionModule {
			return testB524ResponseForSelectorPayload(frame.Data, 0x02, 0x00), nil
		}
		if group == remoteFunctionalModules.group {
			switch addr {
			case deviceSlotConnected:
				connected := byte(0x00)
				if instance == 0x01 {
					connected = 0x01
				}
				return testB524ResponseForSelectorPayload(frame.Data, connected), nil
			case deviceSlotClassAddress:
				classAddress := byte(0xFF)
				if instance == 0x01 {
					classAddress = circuitManagingDeviceFunctionModuleControllerAddress
				}
				return testB524ResponseForSelectorPayload(frame.Data, classAddress), nil
			case deviceSlotFirmware:
				firmware := []byte{0xFF, 0xFF, 0xFF}
				if instance == 0x01 {
					firmware = []byte{0x01, 0x00, 0x00}
				}
				return testB524ResponseForSelectorPayload(frame.Data, firmware...), nil
			case deviceSlotHardwareIdentifier:
				hardware := []byte{0xFF, 0xFF}
				if instance == 0x01 {
					hardware = []byte{0x71, 0x00}
				}
				return testB524ResponseForSelectorPayload(frame.Data, hardware...), nil
			}
		}
		if group == localSolar.group {
			switch addr {
			case solarCollectorTemp:
				return testB524ResponseForSelectorPayload(frame.Data, testB524Float32Payload(61.5)...), nil
			case solarEnabled, solarPumpActive:
				return testB524ResponseForSelectorPayload(frame.Data, 0x01), nil
			}
		}
		if group == localCylinders.group && addr == cylinderTemperature {
			return testB524ResponseForSelectorPayload(frame.Data, testB524Float32Payload(47.25)...), nil
		}
		return testB524ResponseForSelectorPayload(frame.Data, 0x00, 0x00, 0x00, 0x00), nil
	}

	poller.EnqueueAddressIdentityProbe(circuitManagingDeviceFunctionModuleControllerAddress)
	probeTask := taskScheduler.nextTask(context.Background())
	if probeTask == nil {
		t.Fatal("late identity probe was not scheduled")
	}
	canceledProbeCtx, cancelProbe := context.WithCancel(context.Background())
	cancelProbe()
	probeTask.run(canceledProbeCtx)
	taskScheduler.taskDone(probeTask)

	scheduled := make([]semanticTaskKey, 0, 2)
	for len(scheduled) < 2 {
		taskScheduler.mu.Lock()
		pending := len(taskScheduler.queue)
		taskScheduler.mu.Unlock()
		if pending == 0 {
			break
		}
		task := taskScheduler.nextTask(context.Background())
		scheduled = append(scheduled, task.key)
		task.run(context.Background())
		taskScheduler.taskDone(task)
	}
	verdict := provider.FunctionModuleInterpretation()
	if !scanSeamCalled || len(scheduled) != 0 ||
		verdict.Mode != output.FunctionModuleSemanticModeInterpreted || verdict.DegradedReason != "" || verdict.EvidenceRevision == "" {
		t.Fatalf("late FM5 recovery seam_called=%t scheduled=%v verdict=%#v; want directed success with inline bounded recovery, no extra queued tasks, and healthy INTERPRETED without ticker", scanSeamCalled, scheduled, verdict)
	}
}

func TestVaillantOwnerFunctionModuleStartupPositiveFastSlotClassifiesBeforeUnrelatedRadioTailDeadline(t *testing.T) {
	poller, provider := newFunctionModuleDeadlineRecoveryTestPoller(nil)
	config := uint16(2)
	poller.system.ModuleConfigurationFunctionModuleController = &config
	startupCtx, cancelStartup := context.WithCancel(context.Background())
	defer cancelStartup()
	tailEntered := false
	poller.sendFrameFn = func(ctx context.Context, frame protocol.Frame) (*protocol.Frame, error) {
		if len(frame.Data) == 6 &&
			(frame.Data[2] == remoteRegulators.group || frame.Data[2] == remoteThermostats.group) &&
			frame.Data[3] > semanticStartupSlotFastMaxInstance {
			tailEntered = true
			cancelStartup()
		}
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return functionModuleInstanceOneSemanticTestResponse(frame)
	}

	poller.refreshRadioDevicesStartup(startupCtx)
	poller.refreshFunctionModuleSemanticStartup(startupCtx)

	verdict := provider.FunctionModuleInterpretation()
	if verdict.Mode != output.FunctionModuleSemanticModeInterpreted || verdict.DegradedReason != "" || verdict.EvidenceRevision == "" {
		t.Fatalf("deadline-aware startup verdict=%#v tail_entered=%t; want healthy INTERPRETED before unrelated regulator/thermostat tail", verdict, tailEntered)
	}
}

func TestVaillantOwnerEnqueueAddressIdentityProbeSaturatedQueueCoalescesGuaranteedOrderedFunctionModuleRecovery(t *testing.T) {
	originalScanDirected := registryScanDirectedFn
	t.Cleanup(func() { registryScanDirectedFn = originalScanDirected })

	taskScheduler := newSemanticTaskScheduler()
	taskScheduler.maxDepth = 4
	poller, provider := newFunctionModuleDeadlineRecoveryTestPoller(taskScheduler)
	scanSeamCalled := false
	registryScanDirectedFn = func(_ context.Context, _ registry.ScanBus, reg *registry.DeviceRegistry, _ byte, targets []byte) ([]registry.DeviceEntry, error) {
		scanSeamCalled = true
		if len(targets) != 1 || targets[0] != circuitManagingDeviceFunctionModuleControllerAddress {
			return nil, errors.New("unexpected directed identity target")
		}
		entry := reg.Register(registry.DeviceInfo{
			Address:      circuitManagingDeviceFunctionModuleControllerAddress,
			Manufacturer: "Vaillant",
			DeviceID:     circuitManagingDeviceFunctionModuleControllerID,
		})
		return []registry.DeviceEntry{entry}, nil
	}

	readOrder := make([]string, 0, 2)
	systemRecorded := false
	radioRecorded := false
	poller.sendFrameFn = func(_ context.Context, frame protocol.Frame) (*protocol.Frame, error) {
		if len(frame.Data) == 6 {
			group := frame.Data[2]
			addr := uint16(frame.Data[4]) | uint16(frame.Data[5])<<8
			if group == localRegulator.group && addr == systemModuleConfigurationFunctionModule && !systemRecorded {
				readOrder = append(readOrder, "system")
				systemRecorded = true
			} else if (group == remoteRegulators.group || group == remoteThermostats.group || group == remoteFunctionalModules.group) && !radioRecorded {
				readOrder = append(readOrder, "radio")
				radioRecorded = true
			}
		}
		return functionModuleInstanceOneSemanticTestResponse(frame)
	}

	if err := taskScheduler.submitCoalesced(semanticTaskRefreshRadioDevices, semanticTaskPriorityLow, func(ctx context.Context) {
		poller.refreshRadioDevices(ctx)
	}); err != nil {
		t.Fatalf("queue older radio task: %v", err)
	}
	poller.EnqueueAddressIdentityProbe(circuitManagingDeviceFunctionModuleControllerAddress)
	for _, key := range []semanticTaskKey{"fm5-p2-fill-1", "fm5-p2-fill-2"} {
		if err := taskScheduler.submitCoalesced(key, semanticTaskPriorityLow, func(context.Context) {}); err != nil {
			t.Fatalf("fill queue before probe: %v", err)
		}
	}

	var probeTask *semanticTask
	taskScheduler.mu.Lock()
	for i, task := range taskScheduler.queue {
		if task.key == "" {
			probeTask = task
			taskScheduler.queue = append(taskScheduler.queue[:i], taskScheduler.queue[i+1:]...)
			break
		}
	}
	taskScheduler.mu.Unlock()
	if probeTask == nil {
		t.Fatal("late identity probe task missing from saturated queue fixture")
	}
	if err := taskScheduler.submitCoalesced("fm5-p2-fill-3", semanticTaskPriorityLow, func(context.Context) {}); err != nil {
		t.Fatalf("refill queue while identity probe runs: %v", err)
	}
	probeTask.run(context.Background())
	taskScheduler.taskDone(probeTask)

	recoveryTasks := 0
	taskScheduler.mu.Lock()
	for _, task := range taskScheduler.queue {
		if task.key == semanticTaskRefreshSystem || task.key == semanticTaskRefreshRadioDevices {
			recoveryTasks++
		}
	}
	taskScheduler.mu.Unlock()
	for {
		taskScheduler.mu.Lock()
		pending := len(taskScheduler.queue)
		taskScheduler.mu.Unlock()
		if pending == 0 {
			break
		}
		task := taskScheduler.nextTask(context.Background())
		task.run(context.Background())
		taskScheduler.taskDone(task)
	}

	verdict := provider.FunctionModuleInterpretation()
	ordered := len(readOrder) >= 2 && readOrder[0] == "system" && readOrder[1] == "radio"
	if !scanSeamCalled || recoveryTasks != 1 || !ordered ||
		verdict.Mode != output.FunctionModuleSemanticModeInterpreted || verdict.DegradedReason != "" || verdict.EvidenceRevision == "" {
		t.Fatalf("saturated late recovery seam_called=%t recovery_tasks=%d read_order=%v verdict=%#v; want one guaranteed coalesced system-then-radio recovery and healthy INTERPRETED", scanSeamCalled, recoveryTasks, readOrder, verdict)
	}
}

func TestVaillantOwnerEnqueueAddressIdentityProbeLateFunctionModuleRecoveryReadsStartupSubsetBeforeOtherSystem(t *testing.T) {
	originalScanDirected := registryScanDirectedFn
	t.Cleanup(func() { registryScanDirectedFn = originalScanDirected })

	taskScheduler := newSemanticTaskScheduler()
	poller, provider := newFunctionModuleDeadlineRecoveryTestPoller(taskScheduler)
	scanSeamCalled := false
	registryScanDirectedFn = func(_ context.Context, _ registry.ScanBus, reg *registry.DeviceRegistry, _ byte, targets []byte) ([]registry.DeviceEntry, error) {
		scanSeamCalled = true
		if len(targets) != 1 || targets[0] != circuitManagingDeviceFunctionModuleControllerAddress {
			return nil, errors.New("unexpected directed identity target")
		}
		entry := reg.Register(registry.DeviceInfo{
			Address:      circuitManagingDeviceFunctionModuleControllerAddress,
			Manufacturer: "Vaillant",
			DeviceID:     circuitManagingDeviceFunctionModuleControllerID,
		})
		return []registry.DeviceEntry{entry}, nil
	}

	taskCtx, cancelTask := context.WithCancel(context.Background())
	defer cancelTask()
	readOrder := make([]string, 0, 2)
	blockedOutsideStartupSubset := false
	gateRecorded := false
	radioRecorded := false
	poller.sendFrameFn = func(ctx context.Context, frame protocol.Frame) (*protocol.Frame, error) {
		if len(frame.Data) != 6 {
			return nil, errors.New("invalid B524 selector")
		}
		group := frame.Data[2]
		addr := uint16(frame.Data[4]) | uint16(frame.Data[5])<<8
		if group == localRegulator.group && addr == systemModuleConfigurationFunctionModule {
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
			if !gateRecorded {
				readOrder = append(readOrder, "gate")
				gateRecorded = true
			}
			return functionModuleInstanceOneSemanticTestResponse(frame)
		}
		startupSubsetRead := addr == systemScheme || addr == systemWaterPressure
		if group == localRegulator.group && !startupSubsetRead && !radioRecorded {
			blockedOutsideStartupSubset = true
			cancelTask()
			return nil, context.DeadlineExceeded
		}
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		if (group == remoteRegulators.group || group == remoteThermostats.group || group == remoteFunctionalModules.group) && !radioRecorded {
			readOrder = append(readOrder, "radio")
			radioRecorded = true
		}
		return functionModuleInstanceOneSemanticTestResponse(frame)
	}

	poller.EnqueueAddressIdentityProbe(circuitManagingDeviceFunctionModuleControllerAddress)
	probeTask := taskScheduler.nextTask(context.Background())
	if probeTask == nil {
		t.Fatal("late identity probe task was not scheduled")
	}
	probeTask.run(taskCtx)
	taskScheduler.taskDone(probeTask)

	verdict := provider.FunctionModuleInterpretation()
	ordered := len(readOrder) == 2 && readOrder[0] == "gate" && readOrder[1] == "radio"
	if !scanSeamCalled || blockedOutsideStartupSubset || !ordered ||
		verdict.Mode != output.FunctionModuleSemanticModeInterpreted || verdict.DegradedReason != "" || verdict.EvidenceRevision == "" {
		t.Fatalf("deadline-aware late recovery seam_called=%t blocked_outside_startup_subset=%t read_order=%v verdict=%#v; want startup subset gate path then radio and healthy INTERPRETED without ordinary ticker", scanSeamCalled, blockedOutsideStartupSubset, readOrder, verdict)
	}
}

func TestVaillantOwnerEnqueueAddressIdentityProbeFreshEmptySlotCacheForcesBoundedFunctionModuleDiscovery(t *testing.T) {
	originalScanDirected := registryScanDirectedFn
	t.Cleanup(func() { registryScanDirectedFn = originalScanDirected })

	taskScheduler := newSemanticTaskScheduler()
	poller, provider := newFunctionModuleDeadlineRecoveryTestPoller(taskScheduler)
	freshDiscoveryAt := time.Date(2026, 8, 15, 20, 59, 0, 0, time.UTC)
	poller.deviceSlotDiscoveryDone = true
	poller.deviceSlotDiscoveryAt = freshDiscoveryAt
	poller.deviceSlotCache = make(map[deviceSlotKey]bool)

	scanSeamCalled := false
	registryScanDirectedFn = func(_ context.Context, _ registry.ScanBus, reg *registry.DeviceRegistry, _ byte, targets []byte) ([]registry.DeviceEntry, error) {
		scanSeamCalled = true
		if len(targets) != 1 || targets[0] != circuitManagingDeviceFunctionModuleControllerAddress {
			return nil, errors.New("unexpected directed identity target")
		}
		entry := reg.Register(registry.DeviceInfo{
			Address:      circuitManagingDeviceFunctionModuleControllerAddress,
			Manufacturer: "Vaillant",
			DeviceID:     circuitManagingDeviceFunctionModuleControllerID,
		})
		return []registry.DeviceEntry{entry}, nil
	}

	gateReads := 0
	functionModuleSlotReads := 0
	nonFMRadioReads := 0
	poller.sendFrameFn = func(_ context.Context, frame protocol.Frame) (*protocol.Frame, error) {
		if len(frame.Data) != 6 {
			return nil, errors.New("invalid B524 selector")
		}
		group := frame.Data[2]
		addr := uint16(frame.Data[4]) | uint16(frame.Data[5])<<8
		if group == localRegulator.group && addr == systemModuleConfigurationFunctionModule {
			gateReads++
		}
		switch group {
		case remoteFunctionalModules.group:
			functionModuleSlotReads++
		case remoteRegulators.group, remoteThermostats.group:
			nonFMRadioReads++
		}
		return functionModuleInstanceOneSemanticTestResponse(frame)
	}

	poller.EnqueueAddressIdentityProbe(circuitManagingDeviceFunctionModuleControllerAddress)
	probeTask := taskScheduler.nextTask(context.Background())
	if probeTask == nil {
		t.Fatal("late identity probe task was not scheduled")
	}
	probeTask.run(context.Background())
	taskScheduler.taskDone(probeTask)

	poller.mu.Lock()
	globalDiscoveryDone := poller.deviceSlotDiscoveryDone
	globalDiscoveryAt := poller.deviceSlotDiscoveryAt
	nonFMCacheEntries := 0
	for key := range poller.deviceSlotCache {
		if key.Group != remoteFunctionalModules.group {
			nonFMCacheEntries++
		}
	}
	poller.mu.Unlock()
	verdict := provider.FunctionModuleInterpretation()
	if !scanSeamCalled || gateReads == 0 || functionModuleSlotReads == 0 || nonFMRadioReads != 0 ||
		!globalDiscoveryDone || !globalDiscoveryAt.Equal(freshDiscoveryAt) || nonFMCacheEntries != 0 ||
		verdict.Mode != output.FunctionModuleSemanticModeInterpreted || verdict.DegradedReason != "" || verdict.EvidenceRevision == "" {
		t.Fatalf("fresh-cache late recovery seam_called=%t gate_reads=%d fm5_reads=%d non_fm_reads=%d discovery_done=%t discovery_at=%s non_fm_cache=%d verdict=%#v; want bounded forced FM5-only discovery and healthy INTERPRETED without ticker/TTL mutation", scanSeamCalled, gateReads, functionModuleSlotReads, nonFMRadioReads, globalDiscoveryDone, globalDiscoveryAt.Format(time.RFC3339), nonFMCacheEntries, verdict)
	}
}

func TestVaillantOwnerFunctionModuleStartupHighSlotPositivePreemptsUnrelatedRadioTails(t *testing.T) {
	poller, provider := newFunctionModuleDeadlineRecoveryTestPoller(nil)
	config := uint16(2)
	poller.system.ModuleConfigurationFunctionModuleController = &config
	classes := map[byte]byte{0x03: circuitManagingDeviceFunctionModuleControllerAddress}
	startupCtx, cancelStartup := context.WithCancel(context.Background())
	defer cancelStartup()
	unrelatedTailEntered := false
	poller.sendFrameFn = func(ctx context.Context, frame protocol.Frame) (*protocol.Frame, error) {
		if len(frame.Data) == 6 && frame.Data[3] > semanticStartupSlotFastMaxInstance &&
			(frame.Data[2] == remoteRegulators.group || frame.Data[2] == remoteThermostats.group) {
			unrelatedTailEntered = true
			cancelStartup()
		}
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return functionModuleSemanticTestResponseForClasses(frame, classes)
	}

	poller.refreshRadioDevicesStartup(startupCtx)
	poller.refreshFunctionModuleSemanticStartup(startupCtx)

	poller.mu.Lock()
	realSnapshot := poller.radioDevices[radioDeviceKey{Group: remoteFunctionalModules.group, Instance: 0x03}]
	poller.mu.Unlock()
	verdict := provider.FunctionModuleInterpretation()
	if unrelatedTailEntered || realSnapshot == nil || realSnapshot.DeviceClassAddress == nil || *realSnapshot.DeviceClassAddress != circuitManagingDeviceFunctionModuleControllerAddress ||
		verdict.Mode != output.FunctionModuleSemanticModeInterpreted || verdict.DegradedReason != "" || verdict.EvidenceRevision == "" {
		t.Fatalf("high-slot startup unrelated_tail=%t snapshot=%#v verdict=%#v; want FM tail first, real slot 3 VR71, and healthy INTERPRETED", unrelatedTailEntered, realSnapshot, verdict)
	}
}

func TestVaillantOwnerFunctionModuleStartupFastPositiveRetainsDeferredHighSlotRegistrySeed(t *testing.T) {
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
	highKey := radioDeviceKey{Group: remoteRegulators.group, Instance: 0x03}
	seededHigh := poller.registryRadioDeviceSeeds()[highKey]
	if seededHigh == nil || seededHigh.DeviceClassAddress == nil || *seededHigh.DeviceClassAddress != seedAddresses[3] {
		t.Fatalf("high-slot registry seed fixture = %#v; want regulator instance 3 at 0x%02X", seededHigh, seedAddresses[3])
	}

	unrelatedTailEntered := false
	poller.sendFrameFn = func(_ context.Context, frame protocol.Frame) (*protocol.Frame, error) {
		if len(frame.Data) == 6 && frame.Data[3] > semanticStartupSlotFastMaxInstance &&
			(frame.Data[2] == remoteRegulators.group || frame.Data[2] == remoteThermostats.group) {
			unrelatedTailEntered = true
		}
		return functionModuleSemanticTestResponseForClasses(frame, map[byte]byte{0x01: circuitManagingDeviceFunctionModuleControllerAddress})
	}

	poller.refreshRadioDevicesStartup(context.Background())
	poller.refreshFunctionModuleSemanticStartup(context.Background())

	poller.mu.Lock()
	highSnapshot := cloneRadioSnapshot(poller.radioDevices[highKey])
	poller.mu.Unlock()
	publishedHighSeed := false
	for _, device := range provider.RadioDevices() {
		if device.Group == int(highKey.Group) && device.Instance == int(highKey.Instance) &&
			device.DeviceClassAddress != nil && *device.DeviceClassAddress == int(seedAddresses[3]) {
			publishedHighSeed = true
			break
		}
	}
	verdict := provider.FunctionModuleInterpretation()
	if unrelatedTailEntered || highSnapshot == nil || highSnapshot.SlotMode != "registry" ||
		highSnapshot.DeviceClassAddress == nil || *highSnapshot.DeviceClassAddress != seedAddresses[3] || !publishedHighSeed ||
		verdict.Mode != output.FunctionModuleSemanticModeInterpreted {
		t.Fatalf("fast-positive deferred seed tail_entered=%t snapshot=%#v published=%t verdict=%#v; want no non-FM tail I/O and retained high-slot registry seed until deferred scan", unrelatedTailEntered, highSnapshot, publishedHighSeed, verdict)
	}
}
