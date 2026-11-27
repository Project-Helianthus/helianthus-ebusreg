package owner

import (
	"context"
	"fmt"
	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
	"github.com/Project-Helianthus/helianthus-ebusreg/vaillant/productids"
	"testing"
	"time"
)

func TestVaillantOwnerVaillantSemanticPollerRegulatorDetectionCatalogError(t *testing.T) {
	t.Parallel()

	reg := newTestRegistry(
		registry.DeviceInfo{Address: 0x15, Manufacturer: "Vaillant", DeviceID: "BASV", SerialNumber: "21-22-09-0020028521-0082-005409-N4"},
	)

	poller := &vaillantSemanticPoller{
		reg:        reg,
		catalogErr: fmt.Errorf("simulated catalog load error"),
	}

	got := poller.findRegulatorCapability()
	if got != productids.ControllerUnknown {
		t.Fatalf("findRegulatorCapability() = %s; want ControllerUnknown on catalog error", got)
	}
}

func TestVaillantOwnerVaillantSemanticPollerRegulatorDetectionMixedKnownAndUnknown(t *testing.T) {
	t.Parallel()

	catalog, err := productids.LoadCatalog()
	if err != nil {
		t.Fatalf("LoadCatalog() error: %v", err)
	}

	// One device with a known boiler part number, one with no serial.
	// Any unknown device should make the overall result ControllerUnknown.
	reg := newTestRegistry(
		registry.DeviceInfo{Address: 0x15, Manufacturer: "Vaillant", DeviceID: "BASV", SerialNumber: "21-22-09-0010002460-0082-005409-N4"},
		registry.DeviceInfo{Address: 0x60, Manufacturer: "Vaillant", DeviceID: "UI"},
	)

	poller := &vaillantSemanticPoller{
		reg:     reg,
		catalog: catalog,
	}

	got := poller.findRegulatorCapability()
	if got != productids.ControllerUnknown {
		t.Fatalf("findRegulatorCapability() = %s; want ControllerUnknown when any device has unknown classification", got)
	}
}

func TestVaillantOwnerVaillantSemanticPollerRefreshDiscoverySetsRegulatorCapability(t *testing.T) {
	if testing.Short() {
		t.Skip("measured long-running coverage; exercised by full CI")
	}
	t.Parallel()

	catalog, err := productids.LoadCatalog()
	if err != nil {
		t.Fatalf("LoadCatalog() error: %v", err)
	}

	// Registry with a B524-capable controller and a regulator device.
	reg := newTestRegistry(
		registry.DeviceInfo{Address: 0x15, Manufacturer: "Vaillant", DeviceID: "BASV", SerialNumber: "21-22-09-0010002460-0082-005409-N4"},
		registry.DeviceInfo{Address: 0x60, Manufacturer: "Vaillant", DeviceID: "VRC430", SerialNumber: "21-22-09-0020028521-0082-005409-N4"},
	)

	provider := newRetainedSemanticState()
	poller := &vaillantSemanticPoller{
		reg:               reg,
		publication:       newVaillantSemanticPublication(provider, nil, nil),
		catalog:           catalog,
		zones:             make(map[byte]*vaillantZoneSnapshot),
		presence:          make(map[byte]*zonePresenceRecord),
		zoneMissThreshold: 3,
		zoneHitThreshold:  2,
		dhwStaleTTL:       10 * time.Minute,
		b524ProbeFn: func(ctx context.Context, target, opcode, group, instance byte, addr uint16) bool {
			return target == 0x15
		},
	}
	poller.nowFn = func() time.Time { return time.Date(2026, time.February, 26, 12, 0, 0, 0, time.UTC) }

	poller.refreshDiscovery(context.Background())

	poller.mu.Lock()
	gotController := poller.controller
	gotCap := poller.regulatorCapability
	poller.mu.Unlock()

	if gotController != 0x15 {
		t.Fatalf("controller = 0x%02x; want 0x15 (B524 root)", gotController)
	}
	if gotCap != productids.ControllerPresent {
		t.Fatalf("regulatorCapability = %s; want ControllerPresent", gotCap)
	}
}

func TestVaillantOwnerVaillantSemanticPollerRefreshDiscoveryNoBasvSetsRegulatorCapability(t *testing.T) {
	t.Parallel()

	catalog, err := productids.LoadCatalog()
	if err != nil {
		t.Fatalf("LoadCatalog() error: %v", err)
	}

	// Registry with ONLY a regulator — no B524-capable device. refreshDiscovery should
	// still update regulatorCapability even when no B524 root is discovered.
	reg := newTestRegistry(
		registry.DeviceInfo{Address: 0x60, Manufacturer: "Vaillant", DeviceID: "VRC430", SerialNumber: "21-22-09-0020028521-0082-005409-N4"},
	)

	provider := newRetainedSemanticState()
	poller := &vaillantSemanticPoller{
		reg:               reg,
		publication:       newVaillantSemanticPublication(provider, nil, nil),
		catalog:           catalog,
		zones:             make(map[byte]*vaillantZoneSnapshot),
		presence:          make(map[byte]*zonePresenceRecord),
		zoneMissThreshold: 3,
		zoneHitThreshold:  2,
		dhwStaleTTL:       10 * time.Minute,
	}
	poller.nowFn = func() time.Time { return time.Date(2026, time.February, 26, 12, 0, 0, 0, time.UTC) }

	poller.refreshDiscovery(context.Background())

	poller.mu.Lock()
	gotController := poller.controller
	gotCap := poller.regulatorCapability
	poller.mu.Unlock()

	if gotController != 0 {
		t.Fatalf("controller = 0x%02x; want 0 (no B524 root found)", gotController)
	}
	if gotCap != productids.ControllerPresent {
		t.Fatalf("regulatorCapability = %s; want ControllerPresent (regulator in registry even without BASV)", gotCap)
	}
}

func TestVaillantOwnerExtractPartNumberFromSerial(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		serial string
		want   string
	}{
		{
			name:   "standard Vaillant serial",
			serial: "21-22-09-0020184848-0082-005409-N4",
			want:   "0020184848",
		},
		{
			name:   "empty serial",
			serial: "",
			want:   "",
		},
		{
			name:   "serial without dashes",
			serial: "212209002018484800820054",
			want:   "0020184848",
		},
		{
			name:   "short serial",
			serial: "21-22",
			want:   "",
		},
		{
			name:   "non-digit part number field",
			serial: "21-22-09-ABCDEFGHIJ-0082-005409-N4",
			want:   "",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := extractPartNumberFromSerial(test.serial)
			if got != test.want {
				t.Fatalf("extractPartNumberFromSerial(%q) = %q; want %q", test.serial, got, test.want)
			}
		})
	}
}

func TestVaillantOwnerRegulatorRedetectionPresentToAbsentGrace(t *testing.T) {
	t.Parallel()

	catalog, err := productids.LoadCatalog()
	if err != nil {
		t.Fatalf("LoadCatalog() error: %v", err)
	}

	// Start with a regulator present (VRC430).
	reg := newTestRegistry(
		registry.DeviceInfo{Address: 0x15, Manufacturer: "Vaillant", DeviceID: "BASV", SerialNumber: "21-22-09-0010002460-0082-005409-N4"},
		registry.DeviceInfo{Address: 0x60, Manufacturer: "Vaillant", DeviceID: "VRC430", SerialNumber: "21-22-09-0020028521-0082-005409-N4"},
	)

	now := time.Date(2026, time.February, 26, 12, 0, 0, 0, time.UTC)
	poller := &vaillantSemanticPoller{
		reg:                      reg,
		catalog:                  catalog,
		regulatorRecheckInterval: 60 * time.Second,
		regulatorAbsenceGrace:    5 * time.Minute,
		regAbsenceState:          regulatorPresent,
		regulatorCapability:      productids.ControllerPresent,
		registryDeviceCount:      2,
		zones:                    make(map[byte]*vaillantZoneSnapshot),
		presence:                 make(map[byte]*zonePresenceRecord),
		nowFn:                    func() time.Time { return now },
	}

	// Confirm initial state is present.
	poller.refreshRegulatorCapability(context.Background())
	poller.mu.Lock()
	if poller.regAbsenceState != regulatorPresent {
		t.Fatalf("initial state = %s; want PRESENT", poller.regAbsenceState)
	}
	poller.mu.Unlock()

	// Remove the regulator from registry (simulate device disappearance).
	regNew := newTestRegistry(
		registry.DeviceInfo{Address: 0x15, Manufacturer: "Vaillant", DeviceID: "BASV", SerialNumber: "21-22-09-0010002460-0082-005409-N4"},
	)
	poller.reg = regNew

	poller.refreshRegulatorCapability(context.Background())
	poller.mu.Lock()
	state := poller.regAbsenceState
	absenceSince := poller.regAbsenceSince
	poller.mu.Unlock()

	if state != regulatorAbsenceGrace {
		t.Fatalf("state after regulator removal = %s; want ABSENCE_GRACE", state)
	}
	if absenceSince.IsZero() {
		t.Fatal("regAbsenceSince should be set after entering grace")
	}
}

func TestVaillantOwnerRegulatorRedetectionGraceToAbsent(t *testing.T) {
	t.Parallel()

	catalog, err := productids.LoadCatalog()
	if err != nil {
		t.Fatalf("LoadCatalog() error: %v", err)
	}

	// Registry with only the boiler (no regulator).
	reg := newTestRegistry(
		registry.DeviceInfo{Address: 0x15, Manufacturer: "Vaillant", DeviceID: "BASV", SerialNumber: "21-22-09-0010002460-0082-005409-N4"},
	)

	graceDuration := 5 * time.Minute
	now := time.Date(2026, time.February, 26, 12, 0, 0, 0, time.UTC)
	poller := &vaillantSemanticPoller{
		reg:                      reg,
		catalog:                  catalog,
		regulatorRecheckInterval: 60 * time.Second,
		regulatorAbsenceGrace:    graceDuration,
		regAbsenceState:          regulatorAbsenceGrace,
		regAbsenceSince:          now.Add(-4 * time.Minute), // 4 min into grace
		regulatorCapability:      productids.ControllerNone,
		registryDeviceCount:      1,
		zones:                    make(map[byte]*vaillantZoneSnapshot),
		presence:                 make(map[byte]*zonePresenceRecord),
		nowFn:                    func() time.Time { return now },
	}

	// Still within grace window — should stay in grace.
	poller.refreshRegulatorCapability(context.Background())
	poller.mu.Lock()
	if poller.regAbsenceState != regulatorAbsenceGrace {
		t.Fatalf("state at 4min = %s; want ABSENCE_GRACE", poller.regAbsenceState)
	}
	poller.mu.Unlock()

	// Advance time to exactly grace boundary (5 min).
	now = now.Add(1 * time.Minute) // now exactly 5 min past absence start
	poller.nowFn = func() time.Time { return now }

	poller.refreshRegulatorCapability(context.Background())
	poller.mu.Lock()
	state := poller.regAbsenceState
	poller.mu.Unlock()

	if state != regulatorAbsent {
		t.Fatalf("state at exact grace boundary = %s; want ABSENT", state)
	}
}

func TestVaillantOwnerRegulatorRedetectionAbsentToPresent(t *testing.T) {
	t.Parallel()

	catalog, err := productids.LoadCatalog()
	if err != nil {
		t.Fatalf("LoadCatalog() error: %v", err)
	}

	// Start with absent state, then re-add regulator.
	reg := newTestRegistry(
		registry.DeviceInfo{Address: 0x15, Manufacturer: "Vaillant", DeviceID: "BASV", SerialNumber: "21-22-09-0010002460-0082-005409-N4"},
	)

	now := time.Date(2026, time.February, 26, 12, 10, 0, 0, time.UTC)
	poller := &vaillantSemanticPoller{
		reg:                      reg,
		catalog:                  catalog,
		regulatorRecheckInterval: 60 * time.Second,
		regulatorAbsenceGrace:    5 * time.Minute,
		regAbsenceState:          regulatorAbsent,
		regAbsenceSince:          now.Add(-10 * time.Minute),
		regulatorCapability:      productids.ControllerNone,
		registryDeviceCount:      1,
		zones:                    make(map[byte]*vaillantZoneSnapshot),
		presence:                 make(map[byte]*zonePresenceRecord),
		nowFn:                    func() time.Time { return now },
	}

	// Verify it stays absent when no regulator.
	poller.refreshRegulatorCapability(context.Background())
	poller.mu.Lock()
	if poller.regAbsenceState != regulatorAbsent {
		t.Fatalf("state without regulator = %s; want ABSENT", poller.regAbsenceState)
	}
	poller.mu.Unlock()

	// Re-add the regulator.
	poller.reg = newTestRegistry(
		registry.DeviceInfo{Address: 0x15, Manufacturer: "Vaillant", DeviceID: "BASV", SerialNumber: "21-22-09-0010002460-0082-005409-N4"},
		registry.DeviceInfo{Address: 0x60, Manufacturer: "Vaillant", DeviceID: "VRC430", SerialNumber: "21-22-09-0020028521-0082-005409-N4"},
	)

	poller.refreshRegulatorCapability(context.Background())
	poller.mu.Lock()
	state := poller.regAbsenceState
	absenceSince := poller.regAbsenceSince
	poller.mu.Unlock()

	if state != regulatorPresent {
		t.Fatalf("state after re-adding regulator = %s; want PRESENT", state)
	}
	if !absenceSince.IsZero() {
		t.Fatal("regAbsenceSince should be cleared after recovery to present")
	}
}

func TestVaillantOwnerRegulatorRedetectionInventoryTrigger(t *testing.T) {
	t.Parallel()

	catalog, err := productids.LoadCatalog()
	if err != nil {
		t.Fatalf("LoadCatalog() error: %v", err)
	}

	reg := newTestRegistry(
		registry.DeviceInfo{Address: 0x15, Manufacturer: "Vaillant", DeviceID: "BASV", SerialNumber: "21-22-09-0010002460-0082-005409-N4"},
		registry.DeviceInfo{Address: 0x60, Manufacturer: "Vaillant", DeviceID: "VRC430", SerialNumber: "21-22-09-0020028521-0082-005409-N4"},
	)

	now := time.Date(2026, time.February, 26, 12, 0, 0, 0, time.UTC)
	poller := &vaillantSemanticPoller{
		reg:                      reg,
		catalog:                  catalog,
		regulatorRecheckInterval: 60 * time.Second,
		regulatorAbsenceGrace:    5 * time.Minute,
		regAbsenceState:          regulatorPresent,
		regulatorCapability:      productids.ControllerPresent,
		registryDeviceCount:      2,
		zones:                    make(map[byte]*vaillantZoneSnapshot),
		presence:                 make(map[byte]*zonePresenceRecord),
		nowFn:                    func() time.Time { return now },
	}

	// No change — same count. Device count should remain 2.
	poller.refreshRegulatorCapability(context.Background())
	poller.mu.Lock()
	if poller.registryDeviceCount != 2 {
		t.Fatalf("registryDeviceCount = %d; want 2", poller.registryDeviceCount)
	}
	poller.mu.Unlock()

	// Now add a third device — count changes from 2 to 3.
	regNew := newTestRegistry(
		registry.DeviceInfo{Address: 0x15, Manufacturer: "Vaillant", DeviceID: "BASV", SerialNumber: "21-22-09-0010002460-0082-005409-N4"},
		registry.DeviceInfo{Address: 0x60, Manufacturer: "Vaillant", DeviceID: "VRC430", SerialNumber: "21-22-09-0020028521-0082-005409-N4"},
		registry.DeviceInfo{Address: 0x70, Manufacturer: "Vaillant", DeviceID: "UI", SerialNumber: "21-22-09-0020099999-0082-005409-N4"},
	)
	poller.reg = regNew

	// Use a task scheduler so enqueueTask doesn't panic.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	taskScheduler := newSemanticTaskScheduler()
	go taskScheduler.run(ctx)
	poller.tasks = taskScheduler

	poller.refreshRegulatorCapability(context.Background())

	poller.mu.Lock()
	newCount := poller.registryDeviceCount
	poller.mu.Unlock()

	if newCount != 3 {
		t.Fatalf("registryDeviceCount after add = %d; want 3", newCount)
	}

	// Verify inventory tracking works — count went from 2 to 3.
	// The enqueueTask call to refreshDiscovery is triggered by the
	// deviceCount != prevDeviceCount condition. We verify the count
	// tracking which is the necessary precondition for the trigger.
}

func TestVaillantOwnerRefreshRegulatorCapabilityNoChangeNoLog(t *testing.T) {
	t.Parallel()

	catalog, err := productids.LoadCatalog()
	if err != nil {
		t.Fatalf("LoadCatalog() error: %v", err)
	}

	reg := newTestRegistry(
		registry.DeviceInfo{Address: 0x15, Manufacturer: "Vaillant", DeviceID: "BASV", SerialNumber: "21-22-09-0010002460-0082-005409-N4"},
		registry.DeviceInfo{Address: 0x60, Manufacturer: "Vaillant", DeviceID: "VRC430", SerialNumber: "21-22-09-0020028521-0082-005409-N4"},
	)

	now := time.Date(2026, time.February, 26, 12, 0, 0, 0, time.UTC)
	poller := &vaillantSemanticPoller{
		reg:                      reg,
		catalog:                  catalog,
		regulatorRecheckInterval: 60 * time.Second,
		regulatorAbsenceGrace:    5 * time.Minute,
		regAbsenceState:          regulatorPresent,
		regulatorCapability:      productids.ControllerPresent,
		registryDeviceCount:      2,
		zones:                    make(map[byte]*vaillantZoneSnapshot),
		presence:                 make(map[byte]*zonePresenceRecord),
		nowFn:                    func() time.Time { return now },
	}

	// Call multiple times — no state transition should occur.
	for i := 0; i < 5; i++ {
		poller.refreshRegulatorCapability(context.Background())
	}

	poller.mu.Lock()
	state := poller.regAbsenceState
	cap := poller.regulatorCapability
	count := poller.registryDeviceCount
	poller.mu.Unlock()

	if state != regulatorPresent {
		t.Fatalf("state after repeated calls = %s; want PRESENT", state)
	}
	if cap != productids.ControllerPresent {
		t.Fatalf("capability after repeated calls = %s; want ControllerPresent", cap)
	}
	if count != 2 {
		t.Fatalf("registryDeviceCount after repeated calls = %d; want 2", count)
	}
}

func TestVaillantOwnerB524EnergyQueriesCoverage(t *testing.T) {
	t.Parallel()

	type tuple struct {
		channel  string
		usage    string
		period   string
		yearKind string
	}
	expected := map[tuple]bool{
		// All-time (lifetime) totals: undated window.
		{"gas", "climate", "all", ""}:           false,
		{"gas", "hot_water", "all", ""}:         false,
		{"electricity", "climate", "all", ""}:   false,
		{"electricity", "hot_water", "all", ""}: false,
		// Monthly: this month.
		{"gas", "climate", "month", "current"}:           false,
		{"gas", "hot_water", "month", "current"}:         false,
		{"electricity", "climate", "month", "current"}:   false,
		{"electricity", "hot_water", "month", "current"}: false,
		// Monthly: last month.
		{"gas", "climate", "month", "previous"}:           false,
		{"gas", "hot_water", "month", "previous"}:         false,
		{"electricity", "climate", "month", "previous"}:   false,
		{"electricity", "hot_water", "month", "previous"}: false,
	}
	for _, q := range b524EnergyQueries {
		tup := tuple{q.channel, q.usage, q.period, q.yearKind}
		if _, ok := expected[tup]; !ok {
			t.Fatalf("unexpected query in b524EnergyQueries: %+v", tup)
		}
		expected[tup] = true
	}
	for tup, seen := range expected {
		if !seen {
			t.Fatalf("missing query in b524EnergyQueries: %+v", tup)
		}
	}
	if len(b524EnergyQueries) != 12 {
		t.Fatalf("len(b524EnergyQueries) = %d; want 12", len(b524EnergyQueries))
	}
}

func TestVaillantOwnerIsB524ProbeCoherent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		payload []byte
		group   byte
		addr    uint16
		want    bool
	}{
		{
			name:    "4-byte header matching group and addr",
			payload: []byte{0x01, 0x03, 0x1C, 0x00},
			group:   0x03,
			addr:    0x001C,
			want:    true,
		},
		{
			name:    "5-byte header matching group and addr",
			payload: []byte{0x01, 0x01, 0x03, 0x1C, 0x00},
			group:   0x03,
			addr:    0x001C,
			want:    true,
		},
		{
			name:    "full response with value bytes",
			payload: []byte{0x01, 0x01, 0x03, 0x1C, 0x00, 0x42, 0x43},
			group:   0x03,
			addr:    0x001C,
			want:    true,
		},
		{
			name:    "short payload rejected",
			payload: []byte{0x01, 0x03, 0x1C},
			group:   0x03,
			addr:    0x001C,
			want:    false,
		},
		{
			name:    "empty payload rejected",
			payload: nil,
			group:   0x03,
			addr:    0x001C,
			want:    false,
		},
		{
			name:    "group mismatch rejected",
			payload: []byte{0x01, 0x02, 0x1C, 0x00},
			group:   0x03,
			addr:    0x001C,
			want:    false,
		},
		{
			name:    "addr mismatch rejected",
			payload: []byte{0x01, 0x03, 0x1D, 0x00},
			group:   0x03,
			addr:    0x001C,
			want:    false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			if got := isB524ProbeCoherent(test.payload, test.group, test.addr); got != test.want {
				t.Fatalf("isB524ProbeCoherent(%x, 0x%02x, 0x%04x) = %v; want %v", test.payload, test.group, test.addr, got, test.want)
			}
		})
	}
}

func TestVaillantOwnerRadioInventoryRegistryInfoMaterializesFunctionModuleModelPhysicalDevice(t *testing.T) {
	t.Parallel()

	classAddress := uint8(0x26)
	info, ok := radioInventoryRegistryInfo(&vaillantRadioDeviceSnapshot{
		DeviceClassAddress: &classAddress,
		DeviceModel:        "VR71/FM5",
	})
	if !ok {
		t.Fatal("radioInventoryRegistryInfo should materialize VR_71 inventory evidence")
	}
	if info.Address != 0x26 || info.Manufacturer != "Vaillant" || info.DeviceID != "VR_71" {
		t.Fatalf("radioInventoryRegistryInfo = %+v; want VR_71 at 0x26", info)
	}
}

func TestVaillantOwnerRegistryRadioDeviceSeedsUsesKnownRegulatorAndFunctionModuleIdentities(t *testing.T) {
	t.Parallel()

	reg := newTestRegistry(
		registry.DeviceInfo{Address: 0x15, Manufacturer: "Vaillant", DeviceID: "BASV2"},
		registry.DeviceInfo{Address: 0x26, Manufacturer: "Vaillant", DeviceID: "VR_71"},
	)
	poller := newTestPoller(reg)

	seeds := poller.registryRadioDeviceSeeds()

	regulator := seeds[radioDeviceKey{Group: remoteRegulators.group, Instance: 0}]
	if regulator == nil || regulator.DeviceClassAddress == nil || *regulator.DeviceClassAddress != 0x15 {
		t.Fatalf("regulator seed = %+v; want BASV2 class address 0x15", regulator)
	}
	functionModule := seeds[radioDeviceKey{Group: remoteFunctionalModules.group, Instance: 0}]
	if functionModule == nil || functionModule.DeviceClassAddress == nil || *functionModule.DeviceClassAddress != 0x26 {
		t.Fatalf("functional module seed = %+v; want VR_71 class address 0x26", functionModule)
	}
}

func TestVaillantOwnerReconcileDiscoveryPresencePublishesStartupZonesOnFirstHit(t *testing.T) {
	t.Parallel()

	provider := newRetainedSemanticState()
	poller := &vaillantSemanticPoller{
		publication:      newVaillantSemanticPublication(provider, nil, nil),
		zones:            make(map[byte]*vaillantZoneSnapshot),
		presence:         make(map[byte]*zonePresenceRecord),
		zoneHitThreshold: 2,
		nowFn:            time.Now,
	}

	source := poller.reconcileDiscoveryPresence(
		context.Background(),
		map[byte]bool{0: true, 1: true},
		map[byte]bool{0: true, 1: true},
	)
	poller.publishZones(source)

	zones := provider.Zones()
	if len(zones) != 2 {
		t.Fatalf("published zones = %d; want 2 after first startup hit", len(zones))
	}
}

func TestVaillantOwnerPreserveExistingRegistryMetadataKeepsRicherIdentityFields(t *testing.T) {
	t.Parallel()

	reg := newTestRegistry(registry.DeviceInfo{
		Address:         0x26,
		Manufacturer:    "Vaillant",
		DeviceID:        "VR_71",
		SerialNumber:    "21-21-34-0020262148-0082-014267-N7",
		MacAddress:      "00:11:22:33:44:55",
		SoftwareVersion: "0507",
		HardwareVersion: "1704",
	})

	got := preserveExistingRegistryMetadata(reg, registry.DeviceInfo{
		Address:      0x26,
		Manufacturer: "Vaillant",
		DeviceID:     "VR_71",
	})

	if got.SerialNumber != "21-21-34-0020262148-0082-014267-N7" {
		t.Fatalf("SerialNumber = %q; want preserved serial", got.SerialNumber)
	}
	if got.MacAddress != "00:11:22:33:44:55" {
		t.Fatalf("MacAddress = %q; want preserved MAC", got.MacAddress)
	}
	if got.SoftwareVersion != "0507" {
		t.Fatalf("SoftwareVersion = %q; want preserved software version", got.SoftwareVersion)
	}
	if got.HardwareVersion != "1704" {
		t.Fatalf("HardwareVersion = %q; want preserved hardware version", got.HardwareVersion)
	}
}

func TestVaillantOwnerCapabilityFirstDiscoveryFindsRootWithSingleCandidate(t *testing.T) {
	t.Parallel()

	reg := registry.NewDeviceRegistry(nil)
	reg.Register(registry.DeviceInfo{Address: 0x15, Manufacturer: "Vaillant", DeviceID: "BASV2"})

	poller := newTestPoller(reg)
	poller.b524ProbeFn = mockB524Probe(map[byte]bool{0x15: true}, nil)

	addr, err := poller.discoverB524Root(context.Background())
	if err != nil {
		t.Fatalf("discoverB524Root error = %v", err)
	}
	if addr != 0x15 {
		t.Fatalf("discoverB524Root = 0x%02x; want 0x15", addr)
	}
}
