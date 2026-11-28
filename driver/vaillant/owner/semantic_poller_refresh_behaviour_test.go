package owner

import (
	"context"
	"errors"
	"github.com/Project-Helianthus/helianthus-ebusgo/engine"
	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/output"
	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
	"slices"
	"testing"
	"time"
)

func TestVaillantOwnerRefreshDiscoveryPrimesControllerSemanticTasksWhenControllerAppears(t *testing.T) {
	t.Parallel()

	taskScheduler := newSemanticTaskScheduler()
	poller := &vaillantSemanticPoller{
		reg:            newTestRegistry(registry.DeviceInfo{Address: 0x15, Manufacturer: "Vaillant", DeviceID: "BASV"}),
		tasks:          taskScheduler,
		zones:          make(map[byte]*vaillantZoneSnapshot),
		presence:       make(map[byte]*zonePresenceRecord),
		circuits:       make(map[byte]*vaillantCircuitSnapshot),
		radioDevices:   make(map[radioDeviceKey]*vaillantRadioDeviceSnapshot),
		solarCylinders: make(map[byte]*vaillantCylinderSnapshot),
		b524ProbeFn: func(ctx context.Context, target, opcode, group, instance byte, addr uint16) bool {
			return target == 0x15
		},
	}

	poller.refreshDiscovery(context.Background())

	poller.mu.Lock()
	controller := poller.controller
	poller.mu.Unlock()
	if controller != 0x15 {
		t.Fatalf("controller = 0x%02x; want 0x15", controller)
	}

	taskScheduler.mu.Lock()
	queueLen := len(taskScheduler.queue)
	taskScheduler.mu.Unlock()
	if queueLen != 5 {
		t.Fatalf("queued controller tasks = %d; want 5", queueLen)
	}

	poller.refreshDiscovery(context.Background())

	taskScheduler.mu.Lock()
	queueLen = len(taskScheduler.queue)
	taskScheduler.mu.Unlock()
	if queueLen != 5 {
		t.Fatalf("queued controller tasks after unchanged discovery = %d; want unchanged 5", queueLen)
	}
}

func TestVaillantOwnerBoilerStatusRegisterDefinitionsForTierNoReturnTemperatureMapping(t *testing.T) {
	t.Parallel()

	for _, tier := range []boilerStatusTier{
		boilerStatusTierFast,
		boilerStatusTierMedium,
		boilerStatusTierSlow,
	} {
		for _, register := range boilerStatusRegisterDefinitionsForTier(tier) {
			if register.group == localCircuits.group && register.addr == circuitFlowTemp {
				t.Fatalf("tier %v maps GG=0x%02x RR=0x%04x; closed decision forbids using this as boiler return temperature", tier, localCircuits.group, circuitFlowTemp)
			}
		}
	}
}

func TestVaillantOwnerMergeBoilerSnapshotNonDestructivePartialUpdatePreservesLastKnown(t *testing.T) {
	t.Parallel()

	flow := 61.5
	returnTemp := 45.0
	pump := true
	heatingStatus := 3
	updatedPump := false

	merged := mergeBoilerSnapshotNonDestructive(
		&vaillantBoilerSnapshot{
			FlowTemperatureC:         &flow,
			ReturnTemperatureC:       &returnTemp,
			CentralHeatingPumpActive: &pump,
			HeatingStatusRaw:         &heatingStatus,
		},
		&vaillantBoilerSnapshot{
			CentralHeatingPumpActive: &updatedPump,
		},
	)

	if merged == nil {
		t.Fatal("merged snapshot is nil")
	}
	if merged.FlowTemperatureC == nil || *merged.FlowTemperatureC != 61.5 {
		t.Fatalf("merged.FlowTemperatureC = %v; want preserved 61.5", merged.FlowTemperatureC)
	}
	if merged.CentralHeatingPumpActive == nil || *merged.CentralHeatingPumpActive {
		t.Fatalf("merged.CentralHeatingPumpActive = %v; want updated false", merged.CentralHeatingPumpActive)
	}
	if merged.HeatingStatusRaw == nil || *merged.HeatingStatusRaw != 3 {
		t.Fatalf("merged.HeatingStatusRaw = %v; want preserved 3", merged.HeatingStatusRaw)
	}
	if merged.ReturnTemperatureC != nil {
		t.Fatalf("merged.ReturnTemperatureC = %v; want nil (closed decision)", merged.ReturnTemperatureC)
	}
}

func TestVaillantOwnerRefreshBoilerStatusTierNoSuccessfulReadsPreservesSnapshot(t *testing.T) {
	t.Parallel()

	for _, qualified := range []bool{false, true} {
		name := "unqualified"
		if qualified {
			name = "qualified"
		}
		t.Run(name, func(t *testing.T) {
			flow := 63.0
			pump := true
			heatingStatus := 4
			provider := newRetainedSemanticState()
			poller := &vaillantSemanticPoller{
				publication: newVaillantSemanticPublication(provider, nil, nil),
				controller:  0x15,
				boiler: &vaillantBoilerSnapshot{
					FlowTemperatureC:         &flow,
					CentralHeatingPumpActive: &pump,
					HeatingStatusRaw:         &heatingStatus,
				},
			}
			if qualified {
				pub, tick := qualifiedBoilerThermalFixture(t)
				pub.attempt.Reply = append([]byte{0x0d, 0x18, 0}, encodeTempData2c(flow)...)
				poller.boilerThermal = pub
				poller.boilerAddress = pub.binding.Identity.Target
				poller.publishBoilerStatus(semanticSnapshotSourceLive)
				// No new read succeeds; retained qualified evidence stays fresh.
				*tick = 50
				pub.attempt = nil
			}

			poller.refreshBoilerStatusTier(context.Background(), boilerStatusTierFast)

			if poller.boiler == nil {
				t.Fatal("poller.boiler = nil; want preserved last-known snapshot")
			}
			got := poller.boiler.FlowTemperatureC
			if qualified {
				if got == nil || *got != flow {
					t.Fatalf("qualified flow = %v; want preserved %.1f", got, flow)
				}
			} else if got == nil || *got != flow {
				t.Fatalf("unqualified flow = %v; want retained %.1f", got, flow)
			}
			if poller.boiler.CentralHeatingPumpActive == nil || !*poller.boiler.CentralHeatingPumpActive {
				t.Fatalf("poller.boiler.CentralHeatingPumpActive = %v; want preserved true", poller.boiler.CentralHeatingPumpActive)
			}
			if poller.boiler.HeatingStatusRaw == nil || *poller.boiler.HeatingStatusRaw != 4 {
				t.Fatalf("poller.boiler.HeatingStatusRaw = %v; want preserved 4", poller.boiler.HeatingStatusRaw)
			}
		})
	}
}

func TestVaillantOwnerRefreshBoilerStatusFastUsesCachedDHWWithoutB524HotWaterReads(t *testing.T) {
	t.Parallel()

	current := 47.5
	target := 50.0
	var selectors [][]byte
	poller := &vaillantSemanticPoller{
		scheduler:      engine.NewSemanticReadScheduler(),
		publication:    newVaillantSemanticPublication(newRetainedSemanticState(), nil, nil),
		source:         0x7F,
		controller:     0x15,
		requestTimeout: 50 * time.Millisecond,
		dhw: &vaillantDhwSnapshot{
			CurrentTempC:                  &current,
			TargetTempC:                   &target,
			ConfigurationDHWOperationMode: "1",
		},
	}
	poller.sendFrameFn = func(_ context.Context, frame protocol.Frame) (*protocol.Frame, error) {
		selectors = append(selectors, slices.Clone(frame.Data))
		return testB524ResponseForSelector(frame.Data), nil
	}

	poller.refreshBoilerStatusTier(context.Background(), boilerStatusTierFast)

	for _, selector := range selectors {
		if len(selector) == 6 && selector[2] == localDHW.group {
			t.Fatalf("boiler fast issued duplicate DHW B524 selector % x", selector)
		}
	}
	if poller.boiler == nil {
		t.Fatal("poller.boiler = nil; want merged boiler snapshot")
	}
	if poller.boiler.DhwTemperatureC == nil || *poller.boiler.DhwTemperatureC != current {
		t.Fatalf("boiler DhwTemperatureC = %v; want cached %.1f", poller.boiler.DhwTemperatureC, current)
	}
	if poller.boiler.DhwTargetTemperatureC == nil || *poller.boiler.DhwTargetTemperatureC != target {
		t.Fatalf("boiler DhwTargetTemperatureC = %v; want cached %.1f", poller.boiler.DhwTargetTemperatureC, target)
	}
	if poller.boiler.DhwOperatingMode == nil || *poller.boiler.DhwOperatingMode != 1 {
		t.Fatalf("boiler DhwOperatingMode = %v; want cached 1", poller.boiler.DhwOperatingMode)
	}
}

func TestVaillantOwnerRefreshBoilerStatusFastCachedDHWOnlyDoesNotPublishLive(t *testing.T) {
	t.Parallel()

	current := 47.5
	target := 50.0
	provider := newRetainedSemanticState()
	poller := &vaillantSemanticPoller{
		scheduler:      engine.NewSemanticReadScheduler(),
		publication:    newVaillantSemanticPublication(provider, nil, nil),
		source:         0x7F,
		controller:     0x15,
		requestTimeout: 50 * time.Millisecond,
		dhw: &vaillantDhwSnapshot{
			CurrentTempC:                  &current,
			TargetTempC:                   &target,
			ConfigurationDHWOperationMode: "1",
		},
	}
	poller.sendFrameFn = func(context.Context, protocol.Frame) (*protocol.Frame, error) {
		return nil, errors.New("boiler fast unavailable")
	}

	poller.refreshBoilerStatusTier(context.Background(), boilerStatusTierFast)

	if got := provider.BoilerStatus(); got != nil {
		t.Fatalf("BoilerStatus() = %+v; want nil without live boiler/B524 success", got)
	}
}

func TestVaillantOwnerRefreshBoilerStatusFastMergesCachedDHWWhenB509UpdatesWithoutB524Mirrors(t *testing.T) {
	t.Parallel()

	current := 47.5
	target := 50.0
	flow := 42.5
	provider := newRetainedSemanticState()
	poller := &vaillantSemanticPoller{
		scheduler:      engine.NewSemanticReadScheduler(),
		publication:    newVaillantSemanticPublication(provider, nil, nil),
		source:         0x7F,
		controller:     0x15,
		boilerAddress:  0x08,
		requestTimeout: 50 * time.Millisecond,
		dhw: &vaillantDhwSnapshot{
			CurrentTempC:                  &current,
			TargetTempC:                   &target,
			ConfigurationDHWOperationMode: "1",
		},
	}
	pub, _ := qualifiedBoilerThermalFixture(t)
	pub.attempt = nil // Only the actual B509 acquisition may seed flow.
	poller.boilerThermal = pub
	poller.source = pub.binding.Identity.Initiator
	poller.boilerAddress = pub.binding.Identity.Target

	poller.sendFrameFn = func(_ context.Context, frame protocol.Frame) (*protocol.Frame, error) {
		if frame.Primary != vaillantB509Primary || frame.Secondary != vaillantB509Secondary {
			t.Fatalf("unexpected non-B509 active read: primary=0x%02x secondary=0x%02x data=% x", frame.Primary, frame.Secondary, frame.Data)
		}
		if len(frame.Data) < 3 || frame.Data[0] != vaillantB509OpcodeRead {
			t.Fatalf("unexpected B509 request data % x", frame.Data)
		}
		addr := uint16(frame.Data[1])<<8 | uint16(frame.Data[2])
		payload := append([]byte{vaillantB509OpcodeRead, byte(addr >> 8), byte(addr)}, encodeTempData2c(flow)...)
		return &protocol.Frame{Data: payload}, nil
	}

	poller.refreshBoilerStatusTier(context.Background(), boilerStatusTierFast)

	if poller.boiler == nil {
		t.Fatal("poller.boiler = nil; want merged boiler snapshot")
	}
	if poller.boiler.FlowTemperatureC == nil || *poller.boiler.FlowTemperatureC != flow {
		t.Fatalf("boiler FlowTemperatureC = %v; want B509 %.1f", poller.boiler.FlowTemperatureC, flow)
	}
	if poller.boiler.DhwTemperatureC == nil || *poller.boiler.DhwTemperatureC != current {
		t.Fatalf("boiler DhwTemperatureC = %v; want cached %.1f", poller.boiler.DhwTemperatureC, current)
	}
	if poller.boiler.DhwTargetTemperatureC == nil || *poller.boiler.DhwTargetTemperatureC != target {
		t.Fatalf("boiler DhwTargetTemperatureC = %v; want cached %.1f", poller.boiler.DhwTargetTemperatureC, target)
	}
	if poller.boiler.DhwOperatingMode == nil || *poller.boiler.DhwOperatingMode != 1 {
		t.Fatalf("boiler DhwOperatingMode = %v; want cached 1", poller.boiler.DhwOperatingMode)
	}
	status := provider.BoilerStatus()
	if status == nil || status.State.DhwTemperatureC == nil || *status.State.DhwTemperatureC != current {
		t.Fatalf("published boiler status = %+v; want cached DHW on live B509 update", status)
	}
}

func TestVaillantOwnerDeriveCircuitManagingDevice(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		system             *vaillantSystemSnapshot
		functionModuleMode vaillantFunctionModuleMode
		wantRole           output.ManagingDeviceRole
		wantID             *string
		wantAddr           *int
	}{
		{
			name:               "unknown without proven topology",
			system:             nil,
			functionModuleMode: vaillantFunctionModuleModeAbsent,
			wantRole:           output.ManagingDeviceRoleUnknown,
		},
		{
			name: "proven topology emits VR71 function module ownership",
			system: &vaillantSystemSnapshot{
				SystemScheme: uint16Ptr(1),
				ModuleConfigurationFunctionModuleController: uint16Ptr(2),
			},
			functionModuleMode: vaillantFunctionModuleModeInterpreted,
			wantRole:           output.ManagingDeviceRoleFunctionModule,
			wantID:             stringPtr(circuitManagingDeviceFunctionModuleControllerID),
			wantAddr:           intPtr(circuitManagingDeviceFunctionModuleControllerAddress),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := output.DeriveCircuitManagingDevice(nativeSystemSnapshot(test.system), output.NativeFunctionModuleMode(test.functionModuleMode))
			if got.Role != test.wantRole {
				t.Fatalf("managing device role = %q; want %q", got.Role, test.wantRole)
			}
			if !StringPtrEquals(got.DeviceID, test.wantID) {
				t.Fatalf("managing device id = %v; want %v", got.DeviceID, test.wantID)
			}
			if !IntPtrEquals(got.Address, test.wantAddr) {
				t.Fatalf("managing device address = %v; want %v", got.Address, test.wantAddr)
			}
		})
	}
}

func TestVaillantOwnerFunctionModuleInventoryRegistryInfosMaterializesFunctionModuleModelForInterpretedTopology(t *testing.T) {
	t.Parallel()

	got := functionModuleInventoryRegistryInfos(&vaillantSystemSnapshot{
		SystemScheme: uint16Ptr(1),
		ModuleConfigurationFunctionModuleController: uint16Ptr(2),
	}, vaillantFunctionModuleModeInterpreted)
	if len(got) != 1 {
		t.Fatalf("fm5InventoryRegistryInfos(...) len = %d; want 1", len(got))
	}
	if got[0].Address != 0x26 || got[0].Manufacturer != "Vaillant" || got[0].DeviceID != "VR_71" {
		t.Fatalf("fm5InventoryRegistryInfos(...) = %+v; want VR_71 at 0x26", got[0])
	}
}

func TestVaillantOwnerDeriveFunctionModuleSemanticModeTransitions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                        string
		controllerReachable         bool
		functionModuleGateSatisfied bool
		solarReadable               bool
		cylindersReadable           bool
		hasEvidence                 bool
		want                        vaillantFunctionModuleMode
	}{
		{
			name:                        "interpreted",
			controllerReachable:         true,
			functionModuleGateSatisfied: true,
			solarReadable:               true,
			cylindersReadable:           true,
			hasEvidence:                 true,
			want:                        vaillantFunctionModuleModeInterpreted,
		},
		{
			name:                        "gpio only",
			controllerReachable:         false,
			functionModuleGateSatisfied: false,
			solarReadable:               false,
			cylindersReadable:           false,
			hasEvidence:                 true,
			want:                        vaillantFunctionModuleModeGPIOOnly,
		},
		{
			name:                        "partial solar evidence is not interpreted",
			controllerReachable:         true,
			functionModuleGateSatisfied: true,
			solarReadable:               true,
			cylindersReadable:           false,
			hasEvidence:                 true,
			want:                        vaillantFunctionModuleModeGPIOOnly,
		},
		{
			name:                        "partial cylinder evidence is not interpreted",
			controllerReachable:         true,
			functionModuleGateSatisfied: true,
			solarReadable:               false,
			cylindersReadable:           true,
			hasEvidence:                 true,
			want:                        vaillantFunctionModuleModeGPIOOnly,
		},
		{
			name:                        "gate failure prevents interpretation",
			controllerReachable:         true,
			functionModuleGateSatisfied: false,
			solarReadable:               true,
			cylindersReadable:           true,
			hasEvidence:                 true,
			want:                        vaillantFunctionModuleModeGPIOOnly,
		},
		{
			name:                        "absent",
			controllerReachable:         false,
			functionModuleGateSatisfied: false,
			solarReadable:               false,
			cylindersReadable:           false,
			hasEvidence:                 false,
			want:                        vaillantFunctionModuleModeAbsent,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := DeriveFunctionModuleSemanticMode(
				test.controllerReachable,
				test.functionModuleGateSatisfied,
				test.solarReadable,
				test.cylindersReadable,
				test.hasEvidence,
			)
			if got != test.want {
				t.Fatalf("deriveFM5SemanticMode(...) = %s; want %s", got, test.want)
			}
		})
	}
}

func TestVaillantOwnerMergeSystemSnapshotNonDestructivePartialUpdatePreservesLastKnown(t *testing.T) {
	t.Parallel()

	pressure := 1.4
	flow := 35.2
	scheme := uint16(8)
	module := uint16(2)
	maxHumidity := uint16(70)
	updatedOutdoor := 4.5

	merged := mergeSystemSnapshotNonDestructive(
		&vaillantSystemSnapshot{
			SystemWaterPressure:   &pressure,
			SystemFlowTemperature: &flow,
			SystemScheme:          &scheme,
			ModuleConfigurationFunctionModuleController: &module,
			MaxRoomHumidity: &maxHumidity,
		},
		&vaillantSystemSnapshot{
			OutdoorTemperature: &updatedOutdoor,
		},
	)

	if merged == nil {
		t.Fatal("merged snapshot is nil")
	}
	if merged.OutdoorTemperature == nil || *merged.OutdoorTemperature != 4.5 {
		t.Fatalf("merged.OutdoorTemperature = %v; want updated 4.5", merged.OutdoorTemperature)
	}
	if merged.SystemWaterPressure == nil || *merged.SystemWaterPressure != 1.4 {
		t.Fatalf("merged.SystemWaterPressure = %v; want preserved 1.4", merged.SystemWaterPressure)
	}
	if merged.SystemFlowTemperature == nil || *merged.SystemFlowTemperature != 35.2 {
		t.Fatalf("merged.SystemFlowTemperature = %v; want preserved 35.2", merged.SystemFlowTemperature)
	}
	if merged.SystemScheme == nil || *merged.SystemScheme != 8 {
		t.Fatalf("merged.SystemScheme = %v; want preserved 8", merged.SystemScheme)
	}
	if merged.ModuleConfigurationFunctionModuleController == nil || *merged.ModuleConfigurationFunctionModuleController != 2 {
		t.Fatalf("merged.ModuleConfigurationVR71 = %v; want preserved 2", merged.ModuleConfigurationFunctionModuleController)
	}
	if merged.MaxRoomHumidity == nil || *merged.MaxRoomHumidity != 70 {
		t.Fatalf("merged.MaxRoomHumidity = %v; want preserved 70", merged.MaxRoomHumidity)
	}
}

func TestVaillantOwnerRefreshSystemNoSuccessfulReadsPreservesSnapshot(t *testing.T) {
	t.Parallel()

	pressure := 1.7
	scheme := uint16(8)
	module := uint16(2)
	poller := &vaillantSemanticPoller{
		controller: 0x15,
		system: &vaillantSystemSnapshot{
			SystemWaterPressure: &pressure,
			SystemScheme:        &scheme,
			ModuleConfigurationFunctionModuleController: &module,
		},
	}

	poller.refreshSystem(context.Background())

	if poller.system == nil {
		t.Fatal("poller.system = nil; want preserved last-known snapshot")
	}
	if poller.system.SystemWaterPressure == nil || *poller.system.SystemWaterPressure != 1.7 {
		t.Fatalf("poller.system.SystemWaterPressure = %v; want preserved 1.7", poller.system.SystemWaterPressure)
	}
	if poller.system.SystemScheme == nil || *poller.system.SystemScheme != 8 {
		t.Fatalf("poller.system.SystemScheme = %v; want preserved 8", poller.system.SystemScheme)
	}
	if poller.system.ModuleConfigurationFunctionModuleController == nil || *poller.system.ModuleConfigurationFunctionModuleController != 2 {
		t.Fatalf("poller.system.ModuleConfigurationVR71 = %v; want preserved 2", poller.system.ModuleConfigurationFunctionModuleController)
	}
}

func TestVaillantOwnerPublishCircuitsUsesExplicitManagingDeviceForProvenTopology(t *testing.T) {
	t.Parallel()

	scheme := uint16(1)
	module := uint16(2)
	circuitType := uint16(1)
	provider := newRetainedSemanticState()
	poller := &vaillantSemanticPoller{
		publication:        newVaillantSemanticPublication(provider, nil, nil),
		functionModuleMode: vaillantFunctionModuleModeInterpreted,
		system: &vaillantSystemSnapshot{
			SystemScheme: &scheme,
			ModuleConfigurationFunctionModuleController: &module,
		},
		circuits: map[byte]*vaillantCircuitSnapshot{
			0x00: {
				Instance:       0x00,
				Active:         true,
				CircuitTypeRaw: &circuitType,
			},
		},
	}

	poller.publishCircuits(semanticSnapshotSourceLive)

	status := provider.Circuits()
	if len(status) != 1 {
		t.Fatalf("provider.Circuits() = %d entries; want 1", len(status))
	}
	if status[0].ManagingDevice.Role != output.ManagingDeviceRoleFunctionModule {
		t.Fatalf("provider.Circuits()[0].ManagingDevice.Role = %q; want %q", status[0].ManagingDevice.Role, output.ManagingDeviceRoleFunctionModule)
	}
	if !StringPtrEquals(status[0].ManagingDevice.DeviceID, stringPtr(circuitManagingDeviceFunctionModuleControllerID)) {
		t.Fatalf("provider.Circuits()[0].ManagingDevice.DeviceID = %v; want %q", status[0].ManagingDevice.DeviceID, circuitManagingDeviceFunctionModuleControllerID)
	}
	if !IntPtrEquals(status[0].ManagingDevice.Address, intPtr(circuitManagingDeviceFunctionModuleControllerAddress)) {
		t.Fatalf("provider.Circuits()[0].ManagingDevice.Address = %v; want %d", status[0].ManagingDevice.Address, circuitManagingDeviceFunctionModuleControllerAddress)
	}
}

func TestVaillantOwnerPublishCircuitsUsesUnknownManagingDeviceForUnprovenTopology(t *testing.T) {
	t.Parallel()

	scheme := uint16(8)
	module := uint16(2)
	circuitType := uint16(1)
	provider := newRetainedSemanticState()
	poller := &vaillantSemanticPoller{
		publication:        newVaillantSemanticPublication(provider, nil, nil),
		functionModuleMode: vaillantFunctionModuleModeInterpreted,
		system: &vaillantSystemSnapshot{
			SystemScheme: &scheme,
			ModuleConfigurationFunctionModuleController: &module,
		},
		circuits: map[byte]*vaillantCircuitSnapshot{
			0x00: {
				Instance:       0x00,
				Active:         true,
				CircuitTypeRaw: &circuitType,
			},
		},
	}

	poller.publishCircuits(semanticSnapshotSourceLive)

	status := provider.Circuits()
	if len(status) != 1 {
		t.Fatalf("provider.Circuits() = %d entries; want 1", len(status))
	}
	if status[0].ManagingDevice.Role != output.ManagingDeviceRoleUnknown {
		t.Fatalf("provider.Circuits()[0].ManagingDevice.Role = %q; want %q", status[0].ManagingDevice.Role, output.ManagingDeviceRoleUnknown)
	}
	if status[0].ManagingDevice.DeviceID != nil {
		t.Fatalf("provider.Circuits()[0].ManagingDevice.DeviceID = %v; want nil", status[0].ManagingDevice.DeviceID)
	}
	if status[0].ManagingDevice.Address != nil {
		t.Fatalf("provider.Circuits()[0].ManagingDevice.Address = %v; want nil", status[0].ManagingDevice.Address)
	}
}

func TestVaillantOwnerPublishSystemRehydratesCircuitOwnershipAfterSystemArrives(t *testing.T) {
	t.Parallel()

	circuitType := uint16(1)
	scheme := uint16(1)
	module := uint16(2)
	provider := newRetainedSemanticState()
	poller := &vaillantSemanticPoller{
		publication:        newVaillantSemanticPublication(provider, nil, nil),
		functionModuleMode: vaillantFunctionModuleModeInterpreted,
		circuits: map[byte]*vaillantCircuitSnapshot{
			0x00: {
				Instance:       0x00,
				Active:         true,
				CircuitTypeRaw: &circuitType,
			},
		},
	}

	poller.publishCircuits(semanticSnapshotSourceLive)

	initial := provider.Circuits()
	if len(initial) != 1 || initial[0].ManagingDevice.Role != output.ManagingDeviceRoleUnknown {
		t.Fatalf("initial circuit ownership = %#v; want UNKNOWN before system arrives", initial)
	}

	poller.system = &vaillantSystemSnapshot{
		SystemScheme: &scheme,
		ModuleConfigurationFunctionModuleController: &module,
	}
	poller.publishSystem(semanticSnapshotSourceLive)

	status := provider.Circuits()
	if len(status) != 1 {
		t.Fatalf("provider.Circuits() = %d entries; want 1", len(status))
	}
	if status[0].ManagingDevice.Role != output.ManagingDeviceRoleFunctionModule {
		t.Fatalf("provider.Circuits()[0].ManagingDevice.Role = %q; want %q", status[0].ManagingDevice.Role, output.ManagingDeviceRoleFunctionModule)
	}
}

func TestVaillantOwnerMergeCircuitSnapshotNonDestructivePartialUpdatePreservesLastKnown(t *testing.T) {
	t.Parallel()

	circuitType := uint16(1)
	coolingEnabled := uint16(1)
	flowTemp := 41.5
	flowSetpoint := 44.0
	mixer := 67.0
	pumpStarts := uint32(120)
	updatedFlow := 42.25
	updatedPump := uint16(0)

	merged := mergeCircuitSnapshotNonDestructive(
		&vaillantCircuitSnapshot{
			Instance:          0x01,
			Active:            true,
			CircuitTypeRaw:    &circuitType,
			CoolingEnabledRaw: &coolingEnabled,
			FlowTemperatureC:  &flowTemp,
			FlowSetpointC:     &flowSetpoint,
			MixerPositionPct:  &mixer,
			PumpStartsRaw:     &pumpStarts,
		},
		&vaillantCircuitSnapshot{
			Instance:         0x01,
			Active:           true,
			FlowTemperatureC: &updatedFlow,
			PumpStatusRaw:    &updatedPump,
		},
	)

	if merged == nil {
		t.Fatal("merged snapshot is nil")
	}
	if merged.FlowTemperatureC == nil || *merged.FlowTemperatureC != 42.25 {
		t.Fatalf("merged.FlowTemperatureC = %v; want updated 42.25", merged.FlowTemperatureC)
	}
	if merged.FlowSetpointC == nil || *merged.FlowSetpointC != 44.0 {
		t.Fatalf("merged.FlowSetpointC = %v; want preserved 44.0", merged.FlowSetpointC)
	}
	if merged.PumpStatusRaw == nil || *merged.PumpStatusRaw != 0 {
		t.Fatalf("merged.PumpStatusRaw = %v; want updated 0", merged.PumpStatusRaw)
	}
	if merged.CoolingEnabledRaw == nil || *merged.CoolingEnabledRaw != 1 {
		t.Fatalf("merged.CoolingEnabledRaw = %v; want preserved 1", merged.CoolingEnabledRaw)
	}
	if merged.PumpStartsRaw == nil || *merged.PumpStartsRaw != 120 {
		t.Fatalf("merged.PumpStartsRaw = %v; want preserved 120", merged.PumpStartsRaw)
	}
	if merged.CircuitTypeRaw == nil || *merged.CircuitTypeRaw != 1 {
		t.Fatalf("merged.CircuitTypeRaw = %v; want preserved 1", merged.CircuitTypeRaw)
	}
}
