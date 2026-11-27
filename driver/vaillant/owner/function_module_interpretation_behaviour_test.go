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

func TestVaillantOwnerDeriveFunctionModuleInterpretationCoherentClassificationsAndUnavailableBootstrap(t *testing.T) {
	configInterpretable := uint16(2)
	configUnsupported := uint16(3)
	tests := []struct {
		name                string
		controllerReachable bool
		moduleConfig        *uint16
		solarReadable       bool
		cylindersReadable   bool
		hasEvidence         bool
		evidenceStale       bool
		incoherent          bool
		wantMode            vaillantFunctionModuleMode
		wantReason          vaillantFunctionModuleDegradedReason
	}{
		{"incomplete bootstrap", false, nil, false, false, false, false, false, "", ""},
		{"no completed negative identity acquisition", true, &configInterpretable, false, false, false, false, false, "", ""},
		{"configuration deliberately unsupported", true, &configUnsupported, false, false, true, false, false, vaillantFunctionModuleModeGPIOOnly, vaillantFunctionModuleDegradedConfigurationNotInterpretable},
		{"coherent interpretation", true, &configInterpretable, true, true, true, false, false, vaillantFunctionModuleModeInterpreted, ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := deriveFunctionModuleInterpretation(
				test.controllerReachable,
				test.moduleConfig,
				test.solarReadable,
				test.cylindersReadable,
				test.hasEvidence,
				test.evidenceStale,
				test.incoherent,
				"acq-42",
			)
			if got.Mode != test.wantMode || got.DegradedReason != test.wantReason {
				t.Fatalf("deriveFM5Interpretation() = %#v; want mode=%s reason=%s", got, test.wantMode, test.wantReason)
			}
			if test.wantMode == "" && got.EvidenceRevision != "" {
				t.Fatalf("unavailable bootstrap evidence revision = %q; want empty", got.EvidenceRevision)
			}
			if test.wantMode != "" && got.EvidenceRevision != "acq-42" {
				t.Fatalf("evidence revision = %q; want acq-42", got.EvidenceRevision)
			}
		})
	}
}

func TestVaillantOwnerDeriveFunctionModuleInterpretationTransientReasonPrecedence(t *testing.T) {
	configInterpretable := uint16(2)
	tests := []struct {
		name                string
		controllerReachable bool
		moduleConfig        *uint16
		solarReadable       bool
		cylindersReadable   bool
		evidenceStale       bool
		incoherent          bool
		wantReason          vaillantFunctionModuleDegradedReason
	}{
		{"controller unavailable", false, nil, false, false, false, false, vaillantFunctionModuleDegradedControllerUnreachable},
		{"configuration unavailable", true, nil, false, false, false, false, vaillantFunctionModuleDegradedConfigurationUnavailable},
		{"stale identity before family reads", true, &configInterpretable, false, false, true, false, vaillantFunctionModuleDegradedEvidenceStale},
		{"solar read failed before cylinder", true, &configInterpretable, false, false, false, false, vaillantFunctionModuleDegradedSolarAcquisitionFailed},
		{"cylinder read failed", true, &configInterpretable, true, false, false, false, vaillantFunctionModuleDegradedCylinderAcquisitionFailed},
		{"generation changed during acquisition", true, &configInterpretable, true, true, false, true, vaillantFunctionModuleDegradedIncoherentAcquisition},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := deriveFunctionModuleInterpretation(
				test.controllerReachable,
				test.moduleConfig,
				test.solarReadable,
				test.cylindersReadable,
				true,
				test.evidenceStale,
				test.incoherent,
				"acq-42",
			)
			if got.DegradedReason != test.wantReason {
				t.Fatalf("deriveFM5Interpretation() = %#v; want reason=%s", got, test.wantReason)
			}
		})
	}
}

func TestVaillantOwnerApplyFunctionModuleAcquisitionRetainsTransientAndWithdrawsStructural(t *testing.T) {
	collector := 61.5
	temperature := 47.25
	previousSolar := &vaillantSolarSnapshot{CollectorTemperatureC: &collector}
	previousCylinders := map[byte]*vaillantCylinderSnapshot{
		0: {Instance: 0, TemperatureC: &temperature},
	}

	retainedSolar, retainedCylinders := applyFunctionModuleAcquisition(
		previousSolar,
		previousCylinders,
		nil,
		nil,
		vaillantFunctionModuleInterpretation{
			Mode:             vaillantFunctionModuleModeInterpreted,
			DegradedReason:   vaillantFunctionModuleDegradedSolarAcquisitionFailed,
			EvidenceRevision: "acq-43",
		},
	)
	if retainedSolar == nil || retainedSolar.CollectorTemperatureC == nil || *retainedSolar.CollectorTemperatureC != collector {
		t.Fatalf("transient solar retention = %#v; want prior coherent snapshot", retainedSolar)
	}
	if retainedCylinders[0] == nil || retainedCylinders[0].TemperatureC == nil || *retainedCylinders[0].TemperatureC != temperature {
		t.Fatalf("transient cylinder retention = %#v; want prior coherent instance", retainedCylinders)
	}

	withdrawnSolar, withdrawnCylinders := applyFunctionModuleAcquisition(
		previousSolar,
		previousCylinders,
		nil,
		nil,
		vaillantFunctionModuleInterpretation{
			Mode:             vaillantFunctionModuleModeGPIOOnly,
			DegradedReason:   vaillantFunctionModuleDegradedConfigurationNotInterpretable,
			EvidenceRevision: "acq-44",
		},
	)
	if withdrawnSolar != nil || len(withdrawnCylinders) != 0 {
		t.Fatalf("structural withdrawal = solar %#v cylinders %#v; want empty", withdrawnSolar, withdrawnCylinders)
	}
}

func TestVaillantOwnerCommitFunctionModuleAcquisitionRetainsCoherentModeAcrossTransientFailures(t *testing.T) {
	config := uint16(2)
	tests := []struct {
		coherentMode vaillantFunctionModuleMode
		reason       vaillantFunctionModuleDegradedReason
	}{
		{vaillantFunctionModuleModeInterpreted, vaillantFunctionModuleDegradedControllerUnreachable},
		{vaillantFunctionModuleModeInterpreted, vaillantFunctionModuleDegradedConfigurationUnavailable},
		{vaillantFunctionModuleModeInterpreted, vaillantFunctionModuleDegradedSolarAcquisitionFailed},
		{vaillantFunctionModuleModeInterpreted, vaillantFunctionModuleDegradedCylinderAcquisitionFailed},
		{vaillantFunctionModuleModeInterpreted, vaillantFunctionModuleDegradedEvidenceStale},
		{vaillantFunctionModuleModeInterpreted, vaillantFunctionModuleDegradedIncoherentAcquisition},
		{vaillantFunctionModuleModeAbsent, vaillantFunctionModuleDegradedControllerUnreachable},
	}

	for _, test := range tests {
		t.Run(string(test.coherentMode)+"/"+string(test.reason), func(t *testing.T) {
			collector := 61.5
			temperature := 47.25
			poller := &vaillantSemanticPoller{
				controller:                       0x15,
				system:                           &vaillantSystemSnapshot{Controller: 0x15, ModuleConfigurationFunctionModuleController: &config},
				functionModuleMode:               test.coherentMode,
				functionModuleInterpretation:     vaillantFunctionModuleInterpretation{Mode: test.coherentMode, EvidenceRevision: "fm5-g7-a41"},
				functionModuleEvidenceGeneration: 7,
				solarCylinders:                   make(map[byte]*vaillantCylinderSnapshot),
			}
			if test.coherentMode == vaillantFunctionModuleModeInterpreted {
				poller.solar = &vaillantSolarSnapshot{CollectorTemperatureC: &collector}
				poller.solarCylinders[0] = &vaillantCylinderSnapshot{Instance: 0, TemperatureC: &temperature}
			}
			captured := functionModuleEvidenceCapture{
				controller:       0x15,
				moduleConfig:     &config,
				generation:       7,
				registryCoherent: true,
			}

			got := poller.commitFunctionModuleAcquisition(captured, vaillantFunctionModuleInterpretation{
				Mode:             vaillantFunctionModuleModeGPIOOnly,
				DegradedReason:   test.reason,
				EvidenceRevision: "fm5-g7-a42",
			}, nil, nil)

			if got.Mode != test.coherentMode || got.DegradedReason != test.reason {
				t.Fatalf("transient commit verdict = %#v; want retained %s/%s", got, test.coherentMode, test.reason)
			}
			if got.EvidenceRevision != "fm5-g7-a42" || got.EvidenceRevision == "fm5-g7-a41" {
				t.Fatalf("transient commit revision = %q; want advanced fm5-g7-a42", got.EvidenceRevision)
			}
			if poller.functionModuleMode != test.coherentMode {
				t.Fatalf("legacy FM5 scalar = %s; want retained %s", poller.functionModuleMode, test.coherentMode)
			}
			if test.coherentMode == vaillantFunctionModuleModeInterpreted {
				if poller.solar == nil || poller.solar.CollectorTemperatureC == nil || *poller.solar.CollectorTemperatureC != collector {
					t.Fatalf("retained solar = %#v; want prior coherent value", poller.solar)
				}
				if poller.solarCylinders[0] == nil || poller.solarCylinders[0].TemperatureC == nil || *poller.solarCylinders[0].TemperatureC != temperature {
					t.Fatalf("retained cylinders = %#v; want prior coherent value", poller.solarCylinders)
				}
			}
		})
	}
}

func TestVaillantOwnerFunctionModuleEvidenceCaptureDetectsStaleAndGenerationChange(t *testing.T) {
	now := time.Date(2026, 8, 14, 1, 0, 0, 0, time.UTC)
	classAddress := uint8(circuitManagingDeviceFunctionModuleControllerAddress)
	config := uint16(2)
	poller := &vaillantSemanticPoller{
		controller:                       0x15,
		system:                           &vaillantSystemSnapshot{ModuleConfigurationFunctionModuleController: &config},
		radioDevices:                     map[radioDeviceKey]*vaillantRadioDeviceSnapshot{{Group: remoteFunctionalModules.group}: {DeviceClassAddress: &classAddress}},
		functionModuleEvidenceTTL:        5 * time.Minute,
		functionModuleIdentityObservedAt: now.Add(-6 * time.Minute),
		functionModuleEvidenceGeneration: 7,
		nowFn:                            func() time.Time { return now },
	}

	before := poller.captureFunctionModuleEvidence()
	if !before.staleAt(now, poller.functionModuleEvidenceTTL) {
		t.Fatal("retained FM5 identity outside TTL was not stale")
	}
	poller.mu.Lock()
	poller.functionModuleEvidenceGeneration++
	if before.matchesLockedPoller(poller) {
		poller.mu.Unlock()
		t.Fatal("final locked FM5 generation revalidation accepted a changed generation")
	}
	poller.mu.Unlock()
	after := poller.captureFunctionModuleEvidence()
	if before.sameGeneration(after) {
		t.Fatal("FM5 evidence generation change was not detected")
	}
}

func TestVaillantOwnerRefreshFunctionModuleSemanticStaleEvidenceSkipsReadsAndRetainsCoherentFamily(t *testing.T) {
	now := time.Date(2026, 8, 14, 1, 0, 0, 0, time.UTC)
	classAddress := uint8(circuitManagingDeviceFunctionModuleControllerAddress)
	config := uint16(2)
	collector := 61.5
	temperature := 47.25
	provider := newRetainedSemanticState()
	poller := &vaillantSemanticPoller{
		reg:                              registry.NewDeviceRegistry(nil),
		publication:                      newVaillantSemanticPublication(provider, nil, nil),
		controller:                       0x15,
		system:                           &vaillantSystemSnapshot{Controller: 0x15, ModuleConfigurationFunctionModuleController: &config},
		radioDevices:                     map[radioDeviceKey]*vaillantRadioDeviceSnapshot{{Group: remoteFunctionalModules.group}: {DeviceClassAddress: &classAddress}},
		solar:                            &vaillantSolarSnapshot{CollectorTemperatureC: &collector},
		solarCylinders:                   map[byte]*vaillantCylinderSnapshot{0: {Instance: 0, TemperatureC: &temperature}},
		functionModuleEvidenceTTL:        5 * time.Minute,
		functionModuleIdentityObservedAt: now.Add(-6 * time.Minute),
		functionModuleEvidenceGeneration: 7,
		functionModuleMode:               vaillantFunctionModuleModeInterpreted,
		functionModuleInterpretation: vaillantFunctionModuleInterpretation{
			Mode:             vaillantFunctionModuleModeInterpreted,
			EvidenceRevision: "fm5-g7-a41",
		},
		functionModuleEvidenceRevision: 41,
		nowFn:                          func() time.Time { return now },
	}

	poller.refreshFunctionModuleSemantic(context.Background())
	verdict := provider.FunctionModuleInterpretation()
	if verdict.Mode != output.FunctionModuleSemanticModeInterpreted || verdict.DegradedReason != output.FunctionModuleSemanticDegradedReasonEvidenceStale {
		t.Fatalf("stale refresh verdict = %#v; want INTERPRETED/EVIDENCE_STALE", verdict)
	}
	if verdict.EvidenceRevision == "fm5-g7-a41" {
		t.Fatalf("stale refresh evidence revision = %q; want advanced revision", verdict.EvidenceRevision)
	}
	if solar := provider.Solar(); solar == nil || solar.CollectorTemperatureC == nil || *solar.CollectorTemperatureC != collector {
		t.Fatalf("stale refresh solar = %#v; want retained coherent snapshot", solar)
	}
	if cylinders := provider.Cylinders(); len(cylinders) != 1 || cylinders[0].TemperatureC == nil || *cylinders[0].TemperatureC != temperature {
		t.Fatalf("stale refresh cylinders = %#v; want retained coherent instance", cylinders)
	}
}

func TestVaillantOwnerCommitFunctionModuleAcquisitionRegistryMutationAfterPostReadCaptureIsIncoherent(t *testing.T) {
	classAddress := uint8(circuitManagingDeviceFunctionModuleControllerAddress)
	config := uint16(2)
	reg := registry.NewDeviceRegistry(nil)
	poller := &vaillantSemanticPoller{
		reg:                reg,
		controller:         0x15,
		system:             &vaillantSystemSnapshot{Controller: 0x15, ModuleConfigurationFunctionModuleController: &config},
		radioDevices:       map[radioDeviceKey]*vaillantRadioDeviceSnapshot{{Group: remoteFunctionalModules.group}: {DeviceClassAddress: &classAddress}},
		functionModuleMode: vaillantFunctionModuleModeInterpreted,
		functionModuleInterpretation: vaillantFunctionModuleInterpretation{
			Mode:             vaillantFunctionModuleModeInterpreted,
			EvidenceRevision: "fm5-g0-a0",
		},
		solarCylinders: make(map[byte]*vaillantCylinderSnapshot),
		nowFn:          time.Now,
	}
	captured := poller.captureFunctionModuleEvidence()
	if !captured.registryCoherent {
		t.Fatal("initial registry capture is incoherent")
	}
	reg.Register(registry.DeviceInfo{
		Address:      circuitManagingDeviceFunctionModuleControllerAddress,
		Manufacturer: "Vaillant",
		DeviceID:     circuitManagingDeviceFunctionModuleControllerID,
	})

	verdict := poller.commitFunctionModuleAcquisition(captured, vaillantFunctionModuleInterpretation{
		Mode:             vaillantFunctionModuleModeInterpreted,
		EvidenceRevision: "fm5-g0-a1",
	}, &vaillantSolarSnapshot{}, map[byte]*vaillantCylinderSnapshot{})
	if verdict.Mode != vaillantFunctionModuleModeInterpreted || verdict.DegradedReason != vaillantFunctionModuleDegradedIncoherentAcquisition {
		t.Fatalf("registry interleaving verdict = %#v; want INTERPRETED/INCOHERENT_ACQUISITION", verdict)
	}
	if verdict.EvidenceRevision == "fm5-g0-a0" {
		t.Fatalf("registry interleaving revision = %q; want advanced revision", verdict.EvidenceRevision)
	}
}

func TestVaillantOwnerRefreshDiscoveryTransientRootFailureRetainsFunctionModuleClassificationAndValues(t *testing.T) {
	config := uint16(2)
	classAddress := uint8(circuitManagingDeviceFunctionModuleControllerAddress)
	collector := 61.5
	temperature := 47.25
	provider := newRetainedSemanticState()
	previous := vaillantFunctionModuleInterpretation{
		Mode:             vaillantFunctionModuleModeInterpreted,
		EvidenceRevision: "fm5-g7-a41",
	}
	provider.SetFunctionModuleInterpretation(publishedFunctionModuleInterpretation(previous))
	provider.SetSolar(&output.SolarStatus{CollectorTemperatureC: &collector})
	provider.SetCylinders([]output.CylinderStatus{{Index: 0, TemperatureC: &temperature}})

	reg := registry.NewDeviceRegistry(nil)
	reg.Register(registry.DeviceInfo{Address: 0x15, Manufacturer: "Vaillant", DeviceID: "BASV2"})
	poller := &vaillantSemanticPoller{
		reg:                              reg,
		publication:                      newVaillantSemanticPublication(provider, nil, nil),
		controller:                       0x15,
		system:                           &vaillantSystemSnapshot{Controller: 0x15, ModuleConfigurationFunctionModuleController: &config},
		radioDevices:                     map[radioDeviceKey]*vaillantRadioDeviceSnapshot{{Group: remoteFunctionalModules.group}: {DeviceClassAddress: &classAddress}},
		zones:                            make(map[byte]*vaillantZoneSnapshot),
		presence:                         make(map[byte]*zonePresenceRecord),
		circuits:                         make(map[byte]*vaillantCircuitSnapshot),
		functionModuleMode:               vaillantFunctionModuleModeInterpreted,
		functionModuleInterpretation:     previous,
		functionModuleEvidenceRevision:   41,
		functionModuleEvidenceGeneration: 7,
		functionModuleIdentityObservedAt: time.Now(),
		solar:                            &vaillantSolarSnapshot{CollectorTemperatureC: &collector},
		solarCylinders:                   map[byte]*vaillantCylinderSnapshot{0: {Instance: 0, TemperatureC: &temperature}},
		b524ProbeFn:                      mockB524Probe(map[byte]bool{}, nil),
		nowFn:                            time.Now,
	}

	poller.refreshDiscovery(context.Background())

	verdict := provider.FunctionModuleInterpretation()
	if verdict.Mode != output.FunctionModuleSemanticModeInterpreted || verdict.DegradedReason != output.FunctionModuleSemanticDegradedReasonControllerUnreachable {
		t.Errorf("root-discovery failure verdict = %#v; want retained INTERPRETED/CONTROLLER_UNREACHABLE", verdict)
	}
	if verdict.EvidenceRevision == "" || verdict.EvidenceRevision == previous.EvidenceRevision {
		t.Errorf("root-discovery failure revision = %q; want advancement from %q", verdict.EvidenceRevision, previous.EvidenceRevision)
	}
	if got := provider.FunctionModuleSemanticMode(); got != output.FunctionModuleSemanticModeInterpreted {
		t.Errorf("root-discovery failure legacy scalar = %s; want retained INTERPRETED", got)
	}
	if solar := provider.Solar(); solar == nil || solar.CollectorTemperatureC == nil || *solar.CollectorTemperatureC != collector {
		t.Errorf("root-discovery failure solar = %#v; want retained collector temperature %.2f", solar, collector)
	}
	if cylinders := provider.Cylinders(); len(cylinders) != 1 || cylinders[0].TemperatureC == nil || *cylinders[0].TemperatureC != temperature {
		t.Errorf("root-discovery failure cylinders = %#v; want retained temperature %.2f", cylinders, temperature)
	}
}

func TestVaillantOwnerRefreshFunctionModuleSemanticStaleIdentityOutranksUnsupportedConfiguration(t *testing.T) {
	now := time.Date(2026, 8, 15, 12, 0, 0, 0, time.UTC)
	configUnsupported := uint16(3)
	classAddress := uint8(circuitManagingDeviceFunctionModuleControllerAddress)
	newPoller := func(provider *retainedSemanticState) *vaillantSemanticPoller {
		return &vaillantSemanticPoller{
			reg:                              registry.NewDeviceRegistry(nil),
			publication:                      newVaillantSemanticPublication(provider, nil, nil),
			controller:                       0x15,
			system:                           &vaillantSystemSnapshot{Controller: 0x15, ModuleConfigurationFunctionModuleController: &configUnsupported},
			radioDevices:                     map[radioDeviceKey]*vaillantRadioDeviceSnapshot{{Group: remoteFunctionalModules.group}: {DeviceClassAddress: &classAddress}},
			solarCylinders:                   make(map[byte]*vaillantCylinderSnapshot),
			functionModuleEvidenceTTL:        5 * time.Minute,
			functionModuleIdentityObservedAt: now.Add(-6 * time.Minute),
			functionModuleEvidenceGeneration: 7,
			functionModuleEvidenceRevision:   41,
			nowFn:                            func() time.Time { return now },
		}
	}

	t.Run("bootstrap remains unavailable", func(t *testing.T) {
		provider := newRetainedSemanticState()
		poller := newPoller(provider)
		poller.refreshFunctionModuleSemantic(context.Background())
		if got := provider.FunctionModuleInterpretation(); got.Mode != "" || got.DegradedReason != "" || got.EvidenceRevision != "" {
			t.Fatalf("stale unsupported bootstrap verdict = %#v; want unavailable zero tuple", got)
		}
	})

	t.Run("coherent family retains prior classification and values", func(t *testing.T) {
		collector := 61.5
		temperature := 47.25
		provider := newRetainedSemanticState()
		poller := newPoller(provider)
		poller.functionModuleMode = vaillantFunctionModuleModeInterpreted
		poller.functionModuleInterpretation = vaillantFunctionModuleInterpretation{
			Mode:             vaillantFunctionModuleModeInterpreted,
			EvidenceRevision: "fm5-g7-a41",
		}
		poller.solar = &vaillantSolarSnapshot{CollectorTemperatureC: &collector}
		poller.solarCylinders[0] = &vaillantCylinderSnapshot{Instance: 0, TemperatureC: &temperature}

		poller.refreshFunctionModuleSemantic(context.Background())
		got := provider.FunctionModuleInterpretation()
		if got.Mode != output.FunctionModuleSemanticModeInterpreted || got.DegradedReason != output.FunctionModuleSemanticDegradedReasonEvidenceStale {
			t.Errorf("stale unsupported retained verdict = %#v; want INTERPRETED/EVIDENCE_STALE", got)
		}
		if got.EvidenceRevision == "" || got.EvidenceRevision == "fm5-g7-a41" {
			t.Errorf("stale unsupported revision = %q; want advancement", got.EvidenceRevision)
		}
		if solar := provider.Solar(); solar == nil || solar.CollectorTemperatureC == nil || *solar.CollectorTemperatureC != collector {
			t.Errorf("stale unsupported solar = %#v; want retained coherent value", solar)
		}
		if cylinders := provider.Cylinders(); len(cylinders) != 1 || cylinders[0].TemperatureC == nil || *cylinders[0].TemperatureC != temperature {
			t.Errorf("stale unsupported cylinders = %#v; want retained coherent value", cylinders)
		}
	})
}

func TestVaillantOwnerRefreshFunctionModuleSemanticAbsenceRequiresFreshCompletedNegativeIdentityAcquisition(t *testing.T) {
	now := time.Date(2026, 8, 15, 12, 0, 0, 0, time.UTC)
	config := uint16(2)
	newPoller := func(provider *retainedSemanticState) *vaillantSemanticPoller {
		return &vaillantSemanticPoller{
			reg:                       registry.NewDeviceRegistry(nil),
			publication:               newVaillantSemanticPublication(provider, nil, nil),
			controller:                0x15,
			system:                    &vaillantSystemSnapshot{Controller: 0x15, ModuleConfigurationFunctionModuleController: &config},
			radioDevices:              make(map[radioDeviceKey]*vaillantRadioDeviceSnapshot),
			solarCylinders:            make(map[byte]*vaillantCylinderSnapshot),
			functionModuleEvidenceTTL: 5 * time.Minute,
			nowFn:                     func() time.Time { return now },
		}
	}

	provider := newRetainedSemanticState()
	poller := newPoller(provider)
	poller.refreshFunctionModuleSemantic(context.Background())
	if got := provider.FunctionModuleInterpretation(); got.Mode != "" || got.DegradedReason != "" || got.EvidenceRevision != "" {
		t.Errorf("unobserved startup verdict = %#v; want unavailable zero tuple", got)
	}

	provider = newRetainedSemanticState()
	poller = newPoller(provider)
	poller.startupRadioDevicesProbed = true
	poller.functionModuleIdentityScanComplete = true
	poller.functionModuleIdentityObservedAt = now
	poller.refreshFunctionModuleSemantic(context.Background())
	if got := provider.FunctionModuleInterpretation(); got.Mode != output.FunctionModuleSemanticModeAbsent || got.DegradedReason != "" || got.EvidenceRevision == "" {
		t.Errorf("fresh completed negative verdict = %#v; want healthy ABSENT with revision", got)
	}
}

func TestVaillantOwnerFunctionModuleAcquisitionEmptyEvidenceGenerationInterleavingRetainsPriorClassification(t *testing.T) {
	configUnsupported := uint16(3)
	collector := 61.5
	temperature := 47.25
	reg := registry.NewDeviceRegistry(nil)
	poller := &vaillantSemanticPoller{
		reg:                reg,
		controller:         0x15,
		system:             &vaillantSystemSnapshot{Controller: 0x15, ModuleConfigurationFunctionModuleController: &configUnsupported},
		radioDevices:       make(map[radioDeviceKey]*vaillantRadioDeviceSnapshot),
		functionModuleMode: vaillantFunctionModuleModeInterpreted,
		functionModuleInterpretation: vaillantFunctionModuleInterpretation{
			Mode:             vaillantFunctionModuleModeInterpreted,
			EvidenceRevision: "fm5-g7-a41",
		},
		functionModuleEvidenceGeneration: 7,
		functionModuleEvidenceRevision:   41,
		solar:                            &vaillantSolarSnapshot{CollectorTemperatureC: &collector},
		solarCylinders:                   map[byte]*vaillantCylinderSnapshot{0: {Instance: 0, TemperatureC: &temperature}},
		nowFn:                            time.Now,
	}

	before := poller.captureFunctionModuleEvidence()
	reg.Register(registry.DeviceInfo{Address: 0x08, Manufacturer: "Vaillant", DeviceID: "BAI00"})
	after := poller.captureFunctionModuleEvidence()
	if before.hasEvidence() || after.hasEvidence() {
		t.Fatal("non-FM5 registry mutation unexpectedly fabricated FM5 identity evidence")
	}
	if before.sameGeneration(after) {
		t.Fatal("registry observation-generation interleaving was not detected")
	}

	verdict := deriveFunctionModuleInterpretation(
		before.controller != 0,
		before.moduleConfig,
		false,
		false,
		false,
		false,
		true,
		poller.nextFunctionModuleEvidenceRevision(before.generation, after.generation),
	)
	verdict = poller.commitFunctionModuleAcquisition(after, verdict, nil, nil)

	if verdict.Mode != vaillantFunctionModuleModeInterpreted || verdict.DegradedReason != vaillantFunctionModuleDegradedIncoherentAcquisition {
		t.Errorf("empty-evidence interleaving verdict = %#v; want retained INTERPRETED/INCOHERENT_ACQUISITION", verdict)
	}
	if verdict.EvidenceRevision == "" || verdict.EvidenceRevision == "fm5-g7-a41" {
		t.Errorf("empty-evidence interleaving revision = %q; want advancement", verdict.EvidenceRevision)
	}
	if poller.functionModuleMode != vaillantFunctionModuleModeInterpreted {
		t.Errorf("empty-evidence interleaving legacy scalar = %s; want retained INTERPRETED", poller.functionModuleMode)
	}
	if poller.solar == nil || poller.solar.CollectorTemperatureC == nil || *poller.solar.CollectorTemperatureC != collector {
		t.Errorf("empty-evidence interleaving solar = %#v; want retained coherent value", poller.solar)
	}
	if poller.solarCylinders[0] == nil || poller.solarCylinders[0].TemperatureC == nil || *poller.solarCylinders[0].TemperatureC != temperature {
		t.Errorf("empty-evidence interleaving cylinders = %#v; want retained coherent value", poller.solarCylinders)
	}
}

func TestVaillantOwnerFunctionModuleIdentityClassificationRequiresCompleteFunctionalModuleNamespaceScan(t *testing.T) {
	now := time.Date(2026, 8, 15, 15, 0, 0, 0, time.UTC)
	config := uint16(2)
	newPoller := func(provider *retainedSemanticState) *vaillantSemanticPoller {
		return &vaillantSemanticPoller{
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
			deviceSlotDiscoveryDone:   false,
			nowFn:                     func() time.Time { return now },
		}
	}
	assertUnavailable := func(t *testing.T, got output.FunctionModuleInterpretation) {
		t.Helper()
		if got.Mode != "" || got.DegradedReason != "" || got.EvidenceRevision != "" {
			t.Fatalf("partial FM5 namespace verdict = %#v; want unavailable zero tuple", got)
		}
	}
	partialResponder := func(cancel context.CancelFunc) func(context.Context, protocol.Frame) (*protocol.Frame, error) {
		return func(_ context.Context, frame protocol.Frame) (*protocol.Frame, error) {
			if len(frame.Data) != 6 {
				return nil, errors.New("invalid B524 selector")
			}
			if frame.Data[2] == remoteFunctionalModules.group {
				cancel()
				return nil, errors.New("functional-module namespace timeout")
			}
			return testB524ResponseForSelectorPayload(frame.Data, 0x00, 0x00, 0x00, 0x00), nil
		}
	}

	t.Run("startup non-FM5 success does not classify absence", func(t *testing.T) {
		provider := newRetainedSemanticState()
		poller := newPoller(provider)
		ctx, cancel := context.WithCancel(context.Background())
		poller.sendFrameFn = partialResponder(cancel)
		poller.refreshRadioDevicesStartup(ctx)
		poller.refreshFunctionModuleSemantic(context.Background())
		assertUnavailable(t, provider.FunctionModuleInterpretation())
	})

	t.Run("steady non-FM5 success does not classify absence", func(t *testing.T) {
		provider := newRetainedSemanticState()
		poller := newPoller(provider)
		ctx, cancel := context.WithCancel(context.Background())
		poller.sendFrameFn = partialResponder(cancel)
		poller.refreshRadioDevices(ctx)
		assertUnavailable(t, provider.FunctionModuleInterpretation())
	})

	t.Run("complete negative functional-module scan permits absence", func(t *testing.T) {
		provider := newRetainedSemanticState()
		poller := newPoller(provider)
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
		poller.refreshRadioDevicesStartup(context.Background())
		poller.refreshFunctionModuleSemantic(context.Background())
		got := provider.FunctionModuleInterpretation()
		if got.Mode != output.FunctionModuleSemanticModeAbsent || got.DegradedReason != "" || got.EvidenceRevision == "" {
			t.Fatalf("complete negative FM5 namespace verdict = %#v; want healthy ABSENT with revision", got)
		}
	})
}
