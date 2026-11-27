package owner

import (
	"context"
	"errors"
	"github.com/Project-Helianthus/helianthus-ebusgo/engine"
	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/output"
	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
	"testing"
	"time"
)

func TestVaillantOwnerFunctionModuleStartupIdentityClassificationUsesCompleteFunctionalModuleIdentityTuple(t *testing.T) {
	now := time.Date(2026, 8, 15, 16, 0, 0, 0, time.UTC)
	config := uint16(2)
	tests := []struct {
		name             string
		connectedPayload []byte
		firmwarePayload  []byte
		hardwarePayload  []byte
		firmwareTimeout  bool
		wantMode         output.FunctionModuleSemanticMode
	}{
		{
			name:            "firmware-only identity",
			firmwarePayload: []byte{0x00, 0x00, 0x00},
			hardwarePayload: []byte{0xFF, 0xFF},
			wantMode:        output.FunctionModuleSemanticModeInterpreted,
		},
		{
			name:            "hardware-only identity",
			firmwarePayload: []byte{0xFF, 0xFF, 0xFF},
			hardwarePayload: []byte{0x34, 0x12},
			wantMode:        output.FunctionModuleSemanticModeInterpreted,
		},
		{
			name:            "missing firmware field",
			firmwarePayload: []byte{0xFF, 0xFF, 0xFF},
			hardwarePayload: []byte{0xFF, 0xFF},
			firmwareTimeout: true,
		},
		{
			name:             "invalid connected sentinel",
			connectedPayload: []byte{0xFF},
			firmwarePayload:  []byte{0xFF, 0xFF, 0xFF},
			hardwarePayload:  []byte{0xFF, 0xFF},
		},
		{
			name:            "short hardware field",
			firmwarePayload: []byte{0xFF, 0xFF, 0xFF},
			hardwarePayload: []byte{0xFF},
		},
		{
			name:             "connected slot without identity",
			connectedPayload: []byte{0x01},
			firmwarePayload:  []byte{0xFF, 0xFF, 0xFF},
			hardwarePayload:  []byte{0xFF, 0xFF},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
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
			poller.sendFrameFn = func(_ context.Context, frame protocol.Frame) (*protocol.Frame, error) {
				if len(frame.Data) != 6 {
					return nil, errors.New("invalid B524 selector")
				}
				group := frame.Data[2]
				instance := frame.Data[3]
				addr := uint16(frame.Data[4]) | uint16(frame.Data[5])<<8
				if group == remoteFunctionalModules.group {
					switch addr {
					case deviceSlotConnected:
						payload := []byte{0x00}
						if instance == 0x04 && test.connectedPayload != nil {
							payload = test.connectedPayload
						}
						return testB524ResponseForSelectorPayload(frame.Data, payload...), nil
					case deviceSlotClassAddress:
						return testB524ResponseForSelectorPayload(frame.Data, 0xFF), nil
					case deviceSlotFirmware:
						if instance == 0x04 && test.firmwareTimeout {
							return nil, errors.New("firmware timeout")
						}
						payload := []byte{0xFF, 0xFF, 0xFF}
						if instance == 0x04 {
							payload = test.firmwarePayload
						}
						return testB524ResponseForSelectorPayload(frame.Data, payload...), nil
					case deviceSlotHardwareIdentifier:
						payload := []byte{0xFF, 0xFF}
						if instance == 0x04 {
							payload = test.hardwarePayload
						}
						return testB524ResponseForSelectorPayload(frame.Data, payload...), nil
					}
				}
				return testB524ResponseForSelectorPayload(frame.Data, 0x00, 0x00, 0x00, 0x00), nil
			}

			poller.refreshRadioDevicesStartup(context.Background())
			poller.refreshFunctionModuleSemantic(context.Background())
			got := provider.FunctionModuleInterpretation()
			if test.wantMode == "" {
				if got.Mode != "" || got.DegradedReason != "" || got.EvidenceRevision != "" {
					t.Fatalf("incomplete startup identity tuple verdict = %#v; want unavailable zero tuple", got)
				}
				return
			}
			if got.Mode != test.wantMode || got.DegradedReason != "" || got.EvidenceRevision == "" {
				t.Fatalf("startup identity tuple verdict = %#v; want healthy %s with revision", got, test.wantMode)
			}
		})
	}
}

func TestVaillantOwnerRefreshRadioDevicesIncompleteCurrentFunctionModuleIdentityTupleRetainsPriorClassification(t *testing.T) {
	now := time.Date(2026, 8, 15, 17, 0, 0, 0, time.UTC)
	config := uint16(2)
	collector := 61.5
	temperature := 47.25
	previous := vaillantFunctionModuleInterpretation{
		Mode:             vaillantFunctionModuleModeInterpreted,
		EvidenceRevision: "fm5-g7-a41",
	}
	zeroFreshnessObserver := staticSemanticReadWatchObserver{
		observation: engine.WatchObservation{
			State:         engine.WatchObservationStateActive,
			HasDescriptor: true,
			Descriptor: engine.WatchDescriptor{
				SemanticClass:     engine.WatchSemanticClassDebug,
				FreshnessProfile:  engine.WatchFreshnessProfileDebug,
				DecoderID:         "test.fm5.current.identity",
				CorrelationPolicy: engine.WatchCorrelationPolicyRequestResponse,
				DirectApplyPolicy: engine.WatchDirectApplyPolicyNever,
			},
			Sources: []engine.WatchActivationSource{engine.WatchActivationSourcePoller},
		},
	}

	newPoller := func(provider *retainedSemanticState) *vaillantSemanticPoller {
		return &vaillantSemanticPoller{
			scheduler:                          engine.NewSemanticReadScheduler(),
			watchObserver:                      zeroFreshnessObserver,
			reg:                                registry.NewDeviceRegistry(nil),
			publication:                        newVaillantSemanticPublication(provider, nil, nil),
			source:                             0x7F,
			controller:                         0x15,
			requestTimeout:                     50 * time.Millisecond,
			system:                             &vaillantSystemSnapshot{Controller: 0x15, ModuleConfigurationFunctionModuleController: &config},
			deviceSlotCache:                    make(map[deviceSlotKey]bool),
			functionModuleMode:                 vaillantFunctionModuleModeInterpreted,
			functionModuleInterpretation:       previous,
			functionModuleEvidenceRevision:     41,
			functionModuleEvidenceGeneration:   7,
			functionModuleIdentityScanComplete: true,
			functionModuleIdentityObservedAt:   now,
			functionModuleEvidenceTTL:          5 * time.Minute,
			solar:                              &vaillantSolarSnapshot{CollectorTemperatureC: &collector},
			solarCylinders:                     map[byte]*vaillantCylinderSnapshot{0: {Instance: 0, TemperatureC: &temperature}},
			nowFn:                              func() time.Time { return now },
		}
	}
	newProvider := func() *retainedSemanticState {
		provider := newRetainedSemanticState()
		provider.SetFunctionModuleInterpretation(publishedFunctionModuleInterpretation(previous))
		provider.SetSolar(&output.SolarStatus{CollectorTemperatureC: &collector})
		provider.SetCylinders([]output.CylinderStatus{{Index: 0, TemperatureC: &temperature}})
		return provider
	}
	assertRetainedMode := func(t *testing.T, provider *retainedSemanticState, wantMode vaillantFunctionModuleMode) {
		t.Helper()
		published, _ := publishedFunctionModuleMode(wantMode)
		got := provider.FunctionModuleInterpretation()
		if got.Mode != published || got.DegradedReason != output.FunctionModuleSemanticDegradedReasonIncoherentAcquisition {
			t.Errorf("incomplete current identity verdict = %#v; want retained %s/INCOHERENT_ACQUISITION", got, wantMode)
		}
		if got.EvidenceRevision == "" || got.EvidenceRevision == previous.EvidenceRevision {
			t.Errorf("incomplete current identity revision = %q; want advancement from %q", got.EvidenceRevision, previous.EvidenceRevision)
		}
		if gotMode := provider.FunctionModuleSemanticMode(); gotMode != published {
			t.Errorf("incomplete current identity legacy scalar = %s; want retained %s", gotMode, wantMode)
		}
		if solar := provider.Solar(); solar == nil || solar.CollectorTemperatureC == nil || *solar.CollectorTemperatureC != collector {
			t.Errorf("incomplete current identity solar = %#v; want retained collector temperature", solar)
		}
		if cylinders := provider.Cylinders(); len(cylinders) != 1 || cylinders[0].TemperatureC == nil || *cylinders[0].TemperatureC != temperature {
			t.Errorf("incomplete current identity cylinders = %#v; want retained temperature", cylinders)
		}
	}
	assertRetained := func(t *testing.T, provider *retainedSemanticState) {
		t.Helper()
		assertRetainedMode(t, provider, vaillantFunctionModuleModeInterpreted)
	}

	t.Run("same-call firmware identity then detail timeout", func(t *testing.T) {
		functionModuleIdentitySameCallFirmwareDetailTimeout(t, newPoller, newProvider, assertRetained)
	})

	t.Run("full scan invalid connected sentinel", func(t *testing.T) {
		functionModuleIdentityFullScanInvalidConnectedSentinel(t, newPoller, newProvider, assertRetained)
	})

	t.Run("full scan connected slot without identity", func(t *testing.T) {
		functionModuleIdentityFullScanConnectedSlotWithoutIdentity(t, newPoller, newProvider, assertRetainedMode, previous)
	})

	t.Run("cached hardware identity then later detail timeout", func(t *testing.T) {
		functionModuleIdentityCachedHardwareThenDetailTimeout(t, newPoller, newProvider, assertRetained, now)
	})

	t.Run("cached identity then first connected read timeout", func(t *testing.T) {
		functionModuleIdentityCachedIdentityFirstConnectedReadTimeout(t, newPoller, newProvider, assertRetained, now)
	})

	t.Run("complete FM negative survives non-FM detail timeout", func(t *testing.T) {
		functionModuleIdentityCompleteNegativeSurvivesOtherDetailTimeout(t, newPoller, newProvider, previous)
	})
}

func functionModuleIdentitySameCallFirmwareDetailTimeout(t *testing.T, newPoller func(*retainedSemanticState) *vaillantSemanticPoller, newProvider func() *retainedSemanticState, assertRetained func(*testing.T, *retainedSemanticState)) {
	provider := newProvider()
	poller := newPoller(provider)
	firmwareReads := 0
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	poller.sendFrameFn = func(_ context.Context, frame protocol.Frame) (*protocol.Frame, error) {
		if len(frame.Data) != 6 {
			return nil, errors.New("invalid B524 selector")
		}
		group := frame.Data[2]
		instance := frame.Data[3]
		addr := uint16(frame.Data[4]) | uint16(frame.Data[5])<<8
		if group == remoteFunctionalModules.group {
			switch addr {
			case deviceSlotConnected:
				return testB524ResponseForSelectorPayload(frame.Data, 0x00), nil
			case deviceSlotClassAddress:
				return testB524ResponseForSelectorPayload(frame.Data, 0xFF), nil
			case deviceSlotFirmware:
				if instance == 0x04 {
					firmwareReads++
					if firmwareReads > 1 {
						cancel()
						return nil, errors.New("current firmware detail timeout")
					}
					return testB524ResponseForSelectorPayload(frame.Data, 0x00, 0x00, 0x00), nil
				}
				return testB524ResponseForSelectorPayload(frame.Data, 0xFF, 0xFF, 0xFF), nil
			case deviceSlotHardwareIdentifier:
				return testB524ResponseForSelectorPayload(frame.Data, 0xFF, 0xFF), nil
			}
		}
		return testB524ResponseForSelectorPayload(frame.Data, 0x00, 0x00, 0x00, 0x00), nil
	}

	poller.refreshRadioDevices(ctx)
	if firmwareReads < 2 {
		t.Fatalf("firmware reads = %d; want discovery identity followed by current detail read", firmwareReads)
	}
	assertRetained(t, provider)
}

func functionModuleIdentityFullScanInvalidConnectedSentinel(t *testing.T, newPoller func(*retainedSemanticState) *vaillantSemanticPoller, newProvider func() *retainedSemanticState, assertRetained func(*testing.T, *retainedSemanticState)) {
	provider := newProvider()
	poller := newPoller(provider)
	poller.sendFrameFn = func(_ context.Context, frame protocol.Frame) (*protocol.Frame, error) {
		if len(frame.Data) != 6 {
			return nil, errors.New("invalid B524 selector")
		}
		group := frame.Data[2]
		instance := frame.Data[3]
		addr := uint16(frame.Data[4]) | uint16(frame.Data[5])<<8
		if group == remoteFunctionalModules.group {
			switch addr {
			case deviceSlotConnected:
				value := byte(0x00)
				if instance == 0x04 {
					value = 0xFF
				}
				return testB524ResponseForSelectorPayload(frame.Data, value), nil
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
	assertRetained(t, provider)
}

func functionModuleIdentityFullScanConnectedSlotWithoutIdentity(t *testing.T, newPoller func(*retainedSemanticState) *vaillantSemanticPoller, newProvider func() *retainedSemanticState, assertRetainedMode func(*testing.T, *retainedSemanticState, vaillantFunctionModuleMode), previous vaillantFunctionModuleInterpretation) {
	for _, wantMode := range []vaillantFunctionModuleMode{
		vaillantFunctionModuleModeInterpreted,
		vaillantFunctionModuleModeAbsent,
	} {
		t.Run(string(wantMode), func(t *testing.T) {
			provider := newProvider()
			poller := newPoller(provider)
			prior := previous
			prior.Mode = wantMode
			provider.SetFunctionModuleInterpretation(publishedFunctionModuleInterpretation(prior))
			poller.functionModuleMode = wantMode
			poller.functionModuleInterpretation = prior
			poller.sendFrameFn = func(_ context.Context, frame protocol.Frame) (*protocol.Frame, error) {
				if len(frame.Data) != 6 {
					return nil, errors.New("invalid B524 selector")
				}
				group := frame.Data[2]
				instance := frame.Data[3]
				addr := uint16(frame.Data[4]) | uint16(frame.Data[5])<<8
				if group == remoteFunctionalModules.group {
					switch addr {
					case deviceSlotConnected:
						value := byte(0x00)
						if instance == 0x04 {
							value = 0x01
						}
						return testB524ResponseForSelectorPayload(frame.Data, value), nil
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
			assertRetainedMode(t, provider, wantMode)
		})
	}
}

func functionModuleIdentityCachedHardwareThenDetailTimeout(t *testing.T, newPoller func(*retainedSemanticState) *vaillantSemanticPoller, newProvider func() *retainedSemanticState, assertRetained func(*testing.T, *retainedSemanticState), now time.Time) {
	provider := newProvider()
	poller := newPoller(provider)
	hardware := uint16(0x1234)
	poller.radioDevices = map[radioDeviceKey]*vaillantRadioDeviceSnapshot{
		{Group: remoteFunctionalModules.group, Instance: 0x04}: {
			Group:              remoteFunctionalModules.group,
			Instance:           0x04,
			SlotMode:           "inventory",
			HardwareIdentifier: &hardware,
		},
	}
	poller.deviceSlotCache[deviceSlotKey{Group: remoteFunctionalModules.group, Instance: 0x04}] = true
	poller.deviceSlotDiscoveryDone = true
	poller.deviceSlotDiscoveryAt = now
	poller.deviceSlotRediscoveryTTL = time.Hour
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
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
				cancel()
				return nil, errors.New("current hardware detail timeout")
			}
		}
		return testB524ResponseForSelectorPayload(frame.Data, 0x00, 0x00, 0x00, 0x00), nil
	}

	poller.refreshRadioDevices(ctx)
	assertRetained(t, provider)
}

func functionModuleIdentityCachedIdentityFirstConnectedReadTimeout(t *testing.T, newPoller func(*retainedSemanticState) *vaillantSemanticPoller, newProvider func() *retainedSemanticState, assertRetained func(*testing.T, *retainedSemanticState), now time.Time) {
	provider := newProvider()
	poller := newPoller(provider)
	hardware := uint16(0x1234)
	poller.radioDevices = map[radioDeviceKey]*vaillantRadioDeviceSnapshot{
		{Group: remoteFunctionalModules.group, Instance: 0x04}: {
			Group:              remoteFunctionalModules.group,
			Instance:           0x04,
			SlotMode:           "inventory",
			HardwareIdentifier: &hardware,
		},
	}
	poller.deviceSlotCache[deviceSlotKey{Group: remoteFunctionalModules.group, Instance: 0x04}] = true
	poller.deviceSlotDiscoveryDone = true
	poller.deviceSlotDiscoveryAt = now
	poller.deviceSlotRediscoveryTTL = time.Hour
	connectedReads := 0
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	poller.sendFrameFn = func(_ context.Context, frame protocol.Frame) (*protocol.Frame, error) {
		if len(frame.Data) != 6 {
			return nil, errors.New("invalid B524 selector")
		}
		if frame.Data[2] == remoteFunctionalModules.group {
			addr := uint16(frame.Data[4]) | uint16(frame.Data[5])<<8
			if addr == deviceSlotConnected {
				connectedReads++
				cancel()
				return nil, errors.New("current connected timeout")
			}
		}
		return testB524ResponseForSelectorPayload(frame.Data, 0x00, 0x00, 0x00, 0x00), nil
	}

	poller.refreshRadioDevices(ctx)
	if connectedReads == 0 {
		t.Fatal("current device_connected read was not attempted")
	}
	assertRetained(t, provider)
}

func functionModuleIdentityCompleteNegativeSurvivesOtherDetailTimeout(t *testing.T, newPoller func(*retainedSemanticState) *vaillantSemanticPoller, newProvider func() *retainedSemanticState, previous vaillantFunctionModuleInterpretation) {
	provider := newProvider()
	poller := newPoller(provider)
	connected := true
	classAddress := uint8(0x15)
	roomTemperature := 21.375
	roomHumidity := 48.0
	roomMapping := uint16(1)
	poller.radioDevices = map[radioDeviceKey]*vaillantRadioDeviceSnapshot{
		{Group: remoteRegulators.group, Instance: 0x01}: {
			Group:              remoteRegulators.group,
			Instance:           0x01,
			SlotMode:           "active",
			DeviceConnected:    &connected,
			DeviceClassAddress: &classAddress,
			DeviceModel:        "retained regulator",
			RoomTemperatureC:   &roomTemperature,
			RoomHumidityPct:    &roomHumidity,
		},
	}
	poller.zones = map[byte]*vaillantZoneSnapshot{
		0x00: {
			Instance: 0x00,
			Present:  true,
			Name:     "Living",
			ConfigurationRoomTemperatureZoneMappingRaw: &roomMapping,
		},
	}
	poller.publishRadioDevices(semanticSnapshotSourceLive)
	poller.publishZones(semanticSnapshotSourceLive)
	retainedRadioDevices := provider.RadioDevices()
	if len(retainedRadioDevices) != 1 {
		t.Fatalf("retained radio fixture = %#v; want one regulator", retainedRadioDevices)
	}
	retainedZones := provider.Zones()
	if len(retainedZones) != 1 || retainedZones[0].State.CurrentTempC == nil || *retainedZones[0].State.CurrentTempC != roomTemperature || retainedZones[0].State.CurrentHumidityPct == nil || *retainedZones[0].State.CurrentHumidityPct != roomHumidity {
		t.Fatalf("retained zone fixture = %#v; want regulator temperature/humidity fallback", retainedZones)
	}
	regulatorConnectedReads := 0
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	poller.sendFrameFn = func(_ context.Context, frame protocol.Frame) (*protocol.Frame, error) {
		if len(frame.Data) != 6 {
			return nil, errors.New("invalid B524 selector")
		}
		group := frame.Data[2]
		instance := frame.Data[3]
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
		if group == remoteRegulators.group && instance == 0x01 && addr == deviceSlotConnected {
			regulatorConnectedReads++
			if regulatorConnectedReads > 1 {
				cancel()
				return nil, errors.New("current regulator detail timeout")
			}
			return testB524ResponseForSelectorPayload(frame.Data, 0x01), nil
		}
		return testB524ResponseForSelectorPayload(frame.Data, 0x00, 0x00, 0x00, 0x00), nil
	}

	poller.refreshRadioDevices(ctx)
	if regulatorConnectedReads < 2 {
		t.Fatalf("regulator connected reads = %d; want discovery success then detail timeout", regulatorConnectedReads)
	}
	got := provider.FunctionModuleInterpretation()
	if got.Mode != output.FunctionModuleSemanticModeAbsent || got.DegradedReason != "" {
		t.Errorf("complete negative with non-FM detail timeout verdict = %#v; want healthy ABSENT", got)
	}
	if got.EvidenceRevision == "" || got.EvidenceRevision == previous.EvidenceRevision {
		t.Errorf("complete negative with non-FM detail timeout revision = %q; want advancement", got.EvidenceRevision)
	}
	if solar := provider.Solar(); solar == nil || solar.CollectorTemperatureC != nil {
		t.Errorf("complete negative with non-FM detail timeout solar = %#v; want cleared structural plane", solar)
	}
	if cylinders := provider.Cylinders(); len(cylinders) != 0 {
		t.Errorf("complete negative with non-FM detail timeout cylinders = %#v; want cleared structural plane", cylinders)
	}

	poller.publishZones(semanticSnapshotSourceLive)
	if radioDevices := provider.RadioDevices(); !output.RadioDevicesEqual(radioDevices, retainedRadioDevices) {
		t.Errorf("complete negative with non-FM detail timeout radio devices = %#v; want retained %#v", radioDevices, retainedRadioDevices)
	}
	zones := provider.Zones()
	if len(zones) != 1 || zones[0].State.CurrentTempC == nil || *zones[0].State.CurrentTempC != roomTemperature || zones[0].State.CurrentHumidityPct == nil || *zones[0].State.CurrentHumidityPct != roomHumidity {
		t.Errorf("complete negative with non-FM detail timeout zone fallback = %#v; want retained regulator temperature/humidity", zones)
	}
}

func TestVaillantOwnerRefreshSystemFunctionModuleConfigurationRequiresCurrentAcquisition(t *testing.T) {
	now := time.Date(2026, 8, 15, 18, 0, 0, 0, time.UTC)
	configInterpretable := uint16(2)
	configUnsupported := uint16(3)
	classAddress := uint8(circuitManagingDeviceFunctionModuleControllerAddress)
	collector := 61.5
	temperature := 47.25
	previous := vaillantFunctionModuleInterpretation{
		Mode:             vaillantFunctionModuleModeInterpreted,
		EvidenceRevision: "fm5-g7-a41",
	}
	newPoller := func(provider *retainedSemanticState, config *uint16) *vaillantSemanticPoller {
		return &vaillantSemanticPoller{
			scheduler:                        engine.NewSemanticReadScheduler(),
			reg:                              registry.NewDeviceRegistry(nil),
			publication:                      newVaillantSemanticPublication(provider, nil, nil),
			source:                           0x7F,
			controller:                       0x15,
			requestTimeout:                   50 * time.Millisecond,
			system:                           &vaillantSystemSnapshot{Controller: 0x15, ModuleConfigurationFunctionModuleController: config},
			radioDevices:                     map[radioDeviceKey]*vaillantRadioDeviceSnapshot{{Group: remoteFunctionalModules.group}: {DeviceClassAddress: &classAddress}},
			functionModuleMode:               previous.Mode,
			functionModuleInterpretation:     previous,
			functionModuleEvidenceRevision:   41,
			functionModuleEvidenceGeneration: 7,
			functionModuleIdentityObservedAt: now,
			functionModuleEvidenceTTL:        5 * time.Minute,
			solar:                            &vaillantSolarSnapshot{CollectorTemperatureC: &collector},
			solarCylinders:                   map[byte]*vaillantCylinderSnapshot{0: {Instance: 0, TemperatureC: &temperature}},
			nowFn:                            func() time.Time { return now },
		}
	}
	newProvider := func() *retainedSemanticState {
		provider := newRetainedSemanticState()
		provider.SetFunctionModuleInterpretation(publishedFunctionModuleInterpretation(previous))
		provider.SetSolar(&output.SolarStatus{CollectorTemperatureC: &collector})
		provider.SetCylinders([]output.CylinderStatus{{Index: 0, TemperatureC: &temperature}})
		return provider
	}
	assertRetained := func(t *testing.T, provider *retainedSemanticState, familyReads int) {
		t.Helper()
		got := provider.FunctionModuleInterpretation()
		if got.Mode != output.FunctionModuleSemanticModeInterpreted || got.DegradedReason != output.FunctionModuleSemanticDegradedReasonConfigurationUnavailable {
			t.Errorf("configuration timeout verdict = %#v; want retained INTERPRETED/CONFIGURATION_UNAVAILABLE", got)
		}
		if got.EvidenceRevision == "" || got.EvidenceRevision == previous.EvidenceRevision {
			t.Errorf("configuration timeout revision = %q; want advancement from %q", got.EvidenceRevision, previous.EvidenceRevision)
		}
		if familyReads != 0 {
			t.Errorf("configuration timeout family reads = %d; want 0", familyReads)
		}
		if solar := provider.Solar(); solar == nil || solar.CollectorTemperatureC == nil || *solar.CollectorTemperatureC != collector {
			t.Errorf("configuration timeout solar = %#v; want retained collector temperature", solar)
		}
		if cylinders := provider.Cylinders(); len(cylinders) != 1 || cylinders[0].TemperatureC == nil || *cylinders[0].TemperatureC != temperature {
			t.Errorf("configuration timeout cylinders = %#v; want retained temperature", cylinders)
		}
	}

	t.Run("all system fields fail", func(t *testing.T) {
		provider := newProvider()
		poller := newPoller(provider, &configInterpretable)
		familyReads := 0
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		poller.sendFrameFn = func(_ context.Context, frame protocol.Frame) (*protocol.Frame, error) {
			if len(frame.Data) == 6 && (frame.Data[2] == localSolar.group || frame.Data[2] == localCylinders.group) {
				familyReads++
			}
			cancel()
			return nil, errors.New("system acquisition timeout")
		}

		poller.refreshSystem(ctx)
		assertRetained(t, provider, familyReads)
	})

	t.Run("another system field succeeds", func(t *testing.T) {
		provider := newProvider()
		poller := newPoller(provider, &configInterpretable)
		familyReads := 0
		configReads := 0
		poller.sendFrameFn = func(_ context.Context, frame protocol.Frame) (*protocol.Frame, error) {
			if len(frame.Data) != 6 {
				return nil, errors.New("invalid B524 selector")
			}
			group := frame.Data[2]
			addr := uint16(frame.Data[4]) | uint16(frame.Data[5])<<8
			if group == localRegulator.group && addr == systemModuleConfigurationFunctionModule {
				configReads++
				return nil, errors.New("module configuration timeout")
			}
			if group == localSolar.group || group == localCylinders.group {
				familyReads++
			}
			return testB524ResponseForSelectorPayload(frame.Data, 0x00, 0x00, 0x00, 0x00), nil
		}

		poller.refreshSystem(context.Background())
		if configReads == 0 {
			t.Fatal("module_configuration_vr71 read was not attempted")
		}
		assertRetained(t, provider, familyReads)
	})

	t.Run("stale unsupported config cannot classify bootstrap", func(t *testing.T) {
		provider := newRetainedSemanticState()
		poller := newPoller(provider, &configUnsupported)
		poller.functionModuleMode = ""
		poller.functionModuleInterpretation = vaillantFunctionModuleInterpretation{}
		poller.solar = nil
		poller.solarCylinders = make(map[byte]*vaillantCylinderSnapshot)
		poller.sendFrameFn = func(_ context.Context, frame protocol.Frame) (*protocol.Frame, error) {
			if len(frame.Data) != 6 {
				return nil, errors.New("invalid B524 selector")
			}
			group := frame.Data[2]
			addr := uint16(frame.Data[4]) | uint16(frame.Data[5])<<8
			if group == localRegulator.group && addr == systemModuleConfigurationFunctionModule {
				return nil, errors.New("module configuration timeout")
			}
			return testB524ResponseForSelectorPayload(frame.Data, 0x00, 0x00, 0x00, 0x00), nil
		}

		poller.refreshSystem(context.Background())
		if got := provider.FunctionModuleInterpretation(); got != (output.FunctionModuleInterpretation{}) {
			t.Errorf("stale unsupported config bootstrap verdict = %#v; want unavailable zero tuple", got)
		}
	})
}
