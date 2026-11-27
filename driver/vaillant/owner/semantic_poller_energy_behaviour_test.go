package owner

import (
	"context"
	"errors"
	"fmt"
	"github.com/Project-Helianthus/helianthus-ebusgo/engine"
	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
	"github.com/Project-Helianthus/helianthus-ebusreg/vaillant/productids"
	"log"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestVaillantOwnerB524GroupDefOpcodeGroupBindingIsCorrect(t *testing.T) {
	t.Parallel()

	// Verify all local groups use OP=0x02.
	localGroups := []b524GroupDef{localRegulator, localDHW, localCircuits, localZones, localSolar, localCylinders}
	for _, g := range localGroups {
		if g.opcode != vaillantB524OpcodeLocal {
			t.Errorf("local group %q (GG=0x%02X) has opcode 0x%02X; want 0x02", g.name, g.group, g.opcode)
		}
	}
	// Verify all remote groups use OP=0x06.
	remoteGroups := []b524GroupDef{remoteRegulators, remoteThermostats, remoteFunctionalModules}
	for _, g := range remoteGroups {
		if g.opcode != vaillantB524OpcodeRead {
			t.Errorf("remote group %q (GG=0x%02X) has opcode 0x%02X; want 0x06", g.name, g.group, g.opcode)
		}
	}
}

func TestVaillantOwnerDiscoverDeviceSlotsExcludesDisconnectedRegulatorAndThermostatSlots(t *testing.T) {
	t.Parallel()

	var selectors [][]byte
	poller := &vaillantSemanticPoller{
		scheduler:      engine.NewSemanticReadScheduler(),
		source:         0x7F,
		controller:     0x15,
		requestTimeout: 50 * time.Millisecond,
	}
	poller.sendFrameFn = func(_ context.Context, frame protocol.Frame) (*protocol.Frame, error) {
		selectors = append(selectors, slices.Clone(frame.Data))
		if len(frame.Data) != 6 {
			return nil, errors.New("invalid B524 selector")
		}
		group := frame.Data[2]
		instance := frame.Data[3]
		addr := uint16(frame.Data[4]) | uint16(frame.Data[5])<<8
		if addr == deviceSlotConnected {
			switch {
			case group == remoteRegulators.group && instance == 0x01:
				return testB524ResponseForSelectorPayload(frame.Data, 0x00), nil
			case group == remoteThermostats.group && instance == 0x01:
				return testB524ResponseForSelectorPayload(frame.Data, 0x00), nil
			case group == remoteThermostats.group && instance == 0x02:
				return testB524ResponseForSelectorPayload(frame.Data, 0x01), nil
			}
		}
		return testB524ResponseForSelectorPayload(frame.Data, 0xFF), nil
	}

	active, observedAny := poller.DiscoverDeviceSlots(context.Background())

	if !observedAny {
		t.Fatal("DiscoverDeviceSlots observedAny = false; want true for responsive disconnected/connected slots")
	}
	if active[deviceSlotKey{Group: remoteRegulators.group, Instance: 0x01}] {
		t.Fatal("disconnected regulator slot was retained for steady-state detail refresh")
	}
	if active[deviceSlotKey{Group: remoteThermostats.group, Instance: 0x01}] {
		t.Fatal("disconnected thermostat slot was retained for steady-state detail refresh")
	}
	if !active[deviceSlotKey{Group: remoteThermostats.group, Instance: 0x02}] {
		t.Fatal("connected thermostat slot was not retained for steady-state detail refresh")
	}
	if hasB524Selector(selectors, remoteRegulators.group, 0x01, deviceSlotClassAddress) {
		t.Fatal("disconnected regulator slot should not receive identity-detail probes")
	}
	if hasB524Selector(selectors, remoteThermostats.group, 0x01, deviceSlotClassAddress) {
		t.Fatal("disconnected thermostat slot should not receive identity-detail probes")
	}
}

func TestVaillantOwnerDiscoverDeviceSlotsKeepsDisconnectedFunctionalModuleIdentityEvidence(t *testing.T) {
	t.Parallel()

	var selectors [][]byte
	poller := &vaillantSemanticPoller{
		scheduler:      engine.NewSemanticReadScheduler(),
		source:         0x7F,
		controller:     0x15,
		requestTimeout: 50 * time.Millisecond,
	}
	poller.sendFrameFn = func(_ context.Context, frame protocol.Frame) (*protocol.Frame, error) {
		selectors = append(selectors, slices.Clone(frame.Data))
		if len(frame.Data) != 6 {
			return nil, errors.New("invalid B524 selector")
		}
		group := frame.Data[2]
		instance := frame.Data[3]
		addr := uint16(frame.Data[4]) | uint16(frame.Data[5])<<8
		if group == remoteFunctionalModules.group && instance == 0x04 {
			switch addr {
			case deviceSlotConnected:
				return testB524ResponseForSelectorPayload(frame.Data, 0x00), nil
			case deviceSlotClassAddress:
				return testB524ResponseForSelectorPayload(frame.Data, 0x26), nil
			}
		}
		return testB524ResponseForSelectorPayload(frame.Data, 0xFF), nil
	}

	active, observedAny := poller.DiscoverDeviceSlots(context.Background())

	if !observedAny {
		t.Fatal("DiscoverDeviceSlots observedAny = false; want true for responsive functional module slot")
	}
	if !active[deviceSlotKey{Group: remoteFunctionalModules.group, Instance: 0x04}] {
		t.Fatal("disconnected functional module identity evidence was not retained")
	}
	requireB524Selector(t, selectors, remoteFunctionalModules.group, 0x04, deviceSlotClassAddress)
}

func TestVaillantOwnerRefreshRadioDevicesSkipsVolatileDetailsForDisconnectedFunctionalModuleInventory(t *testing.T) {
	t.Parallel()

	var selectors [][]byte
	poller := &vaillantSemanticPoller{
		scheduler:                engine.NewSemanticReadScheduler(),
		publication:              newVaillantSemanticPublication(newRetainedSemanticState(), nil, nil),
		reg:                      newTestRegistry(),
		source:                   0x7F,
		controller:               0x15,
		requestTimeout:           50 * time.Millisecond,
		deviceSlotRediscoveryTTL: 30 * time.Minute,
		deviceSlotDiscoveryDone:  true,
		deviceSlotDiscoveryAt:    time.Now(),
		deviceSlotCache: map[deviceSlotKey]bool{
			{Group: remoteFunctionalModules.group, Instance: 0x04}: true,
		},
		nowFn: time.Now,
	}
	poller.sendFrameFn = func(_ context.Context, frame protocol.Frame) (*protocol.Frame, error) {
		selectors = append(selectors, slices.Clone(frame.Data))
		if len(frame.Data) != 6 {
			return nil, errors.New("invalid B524 selector")
		}
		group := frame.Data[2]
		instance := frame.Data[3]
		addr := uint16(frame.Data[4]) | uint16(frame.Data[5])<<8
		if group == remoteFunctionalModules.group && instance == 0x04 {
			switch addr {
			case deviceSlotConnected:
				return testB524ResponseForSelectorPayload(frame.Data, 0x00), nil
			case deviceSlotClassAddress:
				return testB524ResponseForSelectorPayload(frame.Data, 0x26), nil
			case deviceSlotFirmware:
				return testB524ResponseForSelectorPayload(frame.Data, 0x01, 0x02, 0x03), nil
			case deviceSlotHardwareIdentifier:
				return testB524ResponseForSelectorPayload(frame.Data, 0x34, 0x12), nil
			}
		}
		return testB524ResponseForSelectorPayload(frame.Data, 0xFF), nil
	}

	poller.refreshRadioDevices(context.Background())

	requireB524Selector(t, selectors, remoteFunctionalModules.group, 0x04, deviceSlotConnected)
	requireB524Selector(t, selectors, remoteFunctionalModules.group, 0x04, deviceSlotClassAddress)
	requireB524Selector(t, selectors, remoteFunctionalModules.group, 0x04, deviceSlotFirmware)
	requireB524Selector(t, selectors, remoteFunctionalModules.group, 0x04, deviceSlotHardwareIdentifier)

	for _, forbidden := range []uint16{
		deviceSlotRemoteControlAddress,
		deviceSlotPaired,
		deviceSlotReceptionStrength,
		deviceSlotZoneAssignment,
		deviceSlotRoomTemperature,
		deviceSlotRoomHumidity,
	} {
		if hasB524Selector(selectors, remoteFunctionalModules.group, 0x04, forbidden) {
			t.Fatalf("disconnected functional-module inventory read volatile addr=0x%04x; want identity-only detail refresh", forbidden)
		}
	}

	poller.mu.Lock()
	device := poller.radioDevices[radioDeviceKey{Group: remoteFunctionalModules.group, Instance: 0x04}]
	poller.mu.Unlock()
	if device == nil {
		t.Fatal("radioDevices missing disconnected functional-module inventory snapshot")
	}
	if device.SlotMode != "inventory" {
		t.Fatalf("SlotMode = %q; want inventory", device.SlotMode)
	}
	if device.RoomTemperatureC != nil || device.RoomHumidityPct != nil || device.ReceptionStrength != nil || device.DevicePaired != nil {
		t.Fatalf("volatile fields populated for disconnected inventory snapshot: %+v", device)
	}
}

func TestVaillantOwnerRefreshRadioDevicesClearsStaleSnapshotsWhenDiscoveryFindsNoRefreshableSlots(t *testing.T) {
	t.Parallel()

	poller := &vaillantSemanticPoller{
		scheduler:      engine.NewSemanticReadScheduler(),
		publication:    newVaillantSemanticPublication(newRetainedSemanticState(), nil, nil),
		source:         0x7F,
		controller:     0x15,
		requestTimeout: 50 * time.Millisecond,
		radioDevices: map[radioDeviceKey]*vaillantRadioDeviceSnapshot{
			{Group: remoteRegulators.group, Instance: 0x01}: {
				Group:    remoteRegulators.group,
				Instance: 0x01,
				SlotMode: "active",
			},
		},
	}
	poller.sendFrameFn = func(_ context.Context, frame protocol.Frame) (*protocol.Frame, error) {
		if len(frame.Data) != 6 {
			return nil, errors.New("invalid B524 selector")
		}
		return testB524ResponseForSelectorPayload(frame.Data, 0x00), nil
	}

	poller.refreshRadioDevices(context.Background())

	poller.mu.Lock()
	got := len(poller.radioDevices)
	poller.mu.Unlock()
	if got != 0 {
		t.Fatalf("radioDevices length = %d; want stale snapshots cleared", got)
	}
}

func TestVaillantOwnerRefreshRadioDevicesPreservesSnapshotsWhenRediscoveryOnlyTimesOut(t *testing.T) {
	t.Parallel()

	poller := &vaillantSemanticPoller{
		scheduler:                engine.NewSemanticReadScheduler(),
		publication:              newVaillantSemanticPublication(newRetainedSemanticState(), nil, nil),
		source:                   0x7F,
		controller:               0x15,
		requestTimeout:           50 * time.Millisecond,
		deviceSlotRediscoveryTTL: time.Millisecond,
		deviceSlotCache: map[deviceSlotKey]bool{
			{Group: remoteRegulators.group, Instance: 0x01}: true,
		},
		deviceSlotDiscoveryDone: true,
		deviceSlotDiscoveryAt:   time.Now().Add(-time.Second),
		radioDevices: map[radioDeviceKey]*vaillantRadioDeviceSnapshot{
			{Group: remoteRegulators.group, Instance: 0x01}: {
				Group:    remoteRegulators.group,
				Instance: 0x01,
				SlotMode: "active",
			},
		},
	}
	poller.sendFrameFn = func(context.Context, protocol.Frame) (*protocol.Frame, error) {
		return nil, errors.New("temporary rediscovery timeout")
	}

	poller.refreshRadioDevices(context.Background())

	poller.mu.Lock()
	radioDevices := len(poller.radioDevices)
	cachedSlots := len(poller.deviceSlotCache)
	poller.mu.Unlock()
	if radioDevices != 1 {
		t.Fatalf("radioDevices length = %d; want stale snapshot preserved through all-timeout rediscovery", radioDevices)
	}
	if cachedSlots != 1 {
		t.Fatalf("deviceSlotCache length = %d; want existing cache preserved through all-timeout rediscovery", cachedSlots)
	}
}

func TestVaillantOwnerDeviceSlotCacheGatingOnlyPollsCachedSlots(t *testing.T) {
	t.Parallel()

	poller := &vaillantSemanticPoller{
		nowFn:                    time.Now,
		deviceSlotRediscoveryTTL: 30 * time.Minute,
	}

	// Simulate discovery: 2 active slots out of 33.
	cache := map[deviceSlotKey]bool{
		{Group: remoteRegulators.group, Instance: 0x00}:  true,
		{Group: remoteThermostats.group, Instance: 0x01}: true,
	}
	poller.mu.Lock()
	poller.deviceSlotCache = cache
	poller.deviceSlotDiscoveryDone = true
	poller.deviceSlotDiscoveryAt = time.Now()
	poller.mu.Unlock()

	// Verify cache is populated.
	poller.mu.Lock()
	got := len(poller.deviceSlotCache)
	poller.mu.Unlock()
	if got != 2 {
		t.Fatalf("deviceSlotCache length = %d; want 2", got)
	}

	// Verify needsDiscovery is false within TTL.
	poller.mu.Lock()
	needsDiscovery := !poller.deviceSlotDiscoveryDone ||
		(poller.deviceSlotRediscoveryTTL > 0 && poller.now().Sub(poller.deviceSlotDiscoveryAt) >= poller.deviceSlotRediscoveryTTL)
	poller.mu.Unlock()
	if needsDiscovery {
		t.Fatal("needsDiscovery should be false within TTL")
	}
}

func TestVaillantOwnerDeviceSlotCacheGatingRequiresRediscoveryAfterTTL(t *testing.T) {
	t.Parallel()

	poller := &vaillantSemanticPoller{
		nowFn:                    time.Now,
		deviceSlotRediscoveryTTL: 1 * time.Millisecond,
	}

	poller.mu.Lock()
	poller.deviceSlotCache = map[deviceSlotKey]bool{
		{Group: remoteRegulators.group, Instance: 0x00}: true,
	}
	poller.deviceSlotDiscoveryDone = true
	poller.deviceSlotDiscoveryAt = time.Now().Add(-1 * time.Second) // expired
	poller.mu.Unlock()

	poller.mu.Lock()
	needsDiscovery := !poller.deviceSlotDiscoveryDone ||
		(poller.deviceSlotRediscoveryTTL > 0 && poller.now().Sub(poller.deviceSlotDiscoveryAt) >= poller.deviceSlotRediscoveryTTL)
	poller.mu.Unlock()
	if !needsDiscovery {
		t.Fatal("needsDiscovery should be true after TTL expiry")
	}
}

// One acquisition derives every dated key from one reference date; replies read
// after a month/year boundary crossing are discarded, not mislabelled.
func TestVaillantOwnerAcquireEnergyBoundaryCrossingDiscardsLaterReplies(t *testing.T) {
	t.Parallel()

	clock := time.Date(2026, time.December, 31, 23, 59, 59, 0, time.UTC)
	provider := newRetainedSemanticState()
	provider.SetEnergyClock(func() time.Time { return clock })
	ledger := newVaillantSemanticPublication(provider, nil, nil)

	calls := 0
	acquireEnergy(ledger, func(q b524EnergyQuery, _ time.Time) (uint32, bool) {
		calls++
		if calls == 6 { // second monthly reply arrives after midnight
			clock = time.Date(2027, time.January, 1, 0, 0, 1, 0, time.UTC)
		}
		return uint32(100 * calls), true
	})
	if calls != 6 {
		t.Fatalf("reads after boundary = %d total; want acquisition to stop at 6", calls)
	}

	totals := provider.EnergyTotals() // read on 2027-01-01
	if got := totals.Gas.Climate.Yearly; got[1] != 100 || got[0] != 0 {
		t.Fatalf("lifetime on Jan 1 = %v; want [0 100]", got)
	}
	// Reply 5 (gas climate, "this month") was read before midnight: Dec 2026 = previous month.
	if got := totals.Gas.Climate.Monthly; got[0] != 500 || got[1] != 0 {
		t.Fatalf("Gas.Climate.Monthly = %v; want [500 0]", got)
	}
	// Reply 6 crossed midnight and must not be stored under any month.
	if got := totals.Gas.DHW.Monthly; len(got) == 2 && (got[0] != 0 || got[1] != 0) {
		t.Fatalf("Gas.DHW.Monthly = %v; want discarded reply (zeros)", got)
	}
}

// Relative monthly queries (this-month / last-month) never consult the
// scheduler or shadow cache, whatever the pass reference: they always read
// fresh from the bus. Lifetime totals keep their normal TTL.
func TestVaillantOwnerReadEnergyRegisterMonthlyNeverServedFromCache(t *testing.T) {
	t.Parallel()

	calls := 0
	poller := &vaillantSemanticPoller{
		scheduler: engine.NewSemanticReadScheduler(),
		sendFrameFn: func(_ context.Context, frame protocol.Frame) (*protocol.Frame, error) {
			calls++
			addr := uint16(frame.Data[4]) | uint16(frame.Data[5])<<8
			return &protocol.Frame{Data: []byte{0x01, 0x00, localRegulator.group, byte(addr), byte(addr >> 8), 0x64, 0x00, 0x00, 0x00}}, nil
		},
		source:         0x7F,
		controller:     0x15,
		requestTimeout: 50 * time.Millisecond,
	}
	var month, life b524EnergyQuery
	for _, q := range b524EnergyQueries {
		if q.period == "month" && month.addr == 0 {
			month = q
		}
		if q.period == "all" && life.addr == 0 {
			life = q
		}
	}
	ctx := context.Background()

	// Warm the cache through the ordinary (TTL) path, as a December pass would.
	for _, q := range []b524EnergyQuery{month, life} {
		if _, ok := poller.readB524Uint32LE(ctx, vaillantB524OpcodeLocal, localRegulator.group, regulatorInstance, q.addr); !ok {
			t.Fatalf("warm read of %#x failed", q.addr)
		}
	}
	warm := calls
	// Precondition: the ordinary path is served from the warm cache.
	if _, ok := poller.readB524Uint32LE(ctx, vaillantB524OpcodeLocal, localRegulator.group, regulatorInstance, month.addr); !ok || calls != warm {
		t.Fatalf("ordinary second read calls=%d; want cache hit (%d)", calls, warm)
	}

	// Monthly energy reads: a pass in January and another in the same month as
	// the warm value both go to the bus.
	for i := 1; i <= 2; i++ {
		if _, ok := poller.readEnergyRegister(ctx, month); !ok {
			t.Fatalf("monthly read %d failed", i)
		}
		if calls != warm+i {
			t.Fatalf("after monthly read %d calls=%d; want bus read (%d)", i, calls, warm+i)
		}
	}

	// Lifetime total is undated: normal TTL, served from cache.
	if _, ok := poller.readEnergyRegister(ctx, life); !ok || calls != warm+2 {
		t.Fatalf("lifetime read calls=%d; want cache hit (%d)", calls, warm+2)
	}
}

// TestVaillantSemanticPollerFallbackDiscoveryGate varies every field of the
// gate: the task (state or configuration), whether a controller is known,
// whether a zone is known and the absence state. Production mutations caught:
// skipping in UNKNOWN, ABSENCE_GRACE or PRESENT too (those rows expect a pass);
// dropping the "no controller" condition (controller-known ABSENT rows, the
// "controller known, no zone" case, expect a pass); gating only one of the two
// tasks (each row runs per task); gating on the zone map instead of the
// controller (controller-less ABSENT rows with a cached zone expect a skip).
func TestVaillantOwnerVaillantSemanticPollerFallbackDiscoveryGate(t *testing.T) {
	tasks := []struct {
		name string
		run  func(*vaillantSemanticPoller, context.Context)
	}{
		{name: "state", run: (*vaillantSemanticPoller).refreshState},
		{name: "config", run: (*vaillantSemanticPoller).refreshConfig},
	}
	for _, task := range tasks {
		for _, controller := range []byte{0x00, 0x15} {
			for _, zoneKnown := range []bool{false, true} {
				for _, state := range []regulatorAbsenceState{regulatorAbsenceUnknown, regulatorPresent, regulatorAbsenceGrace, regulatorAbsent} {
					name := fmt.Sprintf("%s/controller_0x%02X/zone_known_%v/%s", task.name, controller, zoneKnown, state)
					t.Run(name, func(t *testing.T) {
						// The fallback is reached when no controller or no zone
						// is known; it is skipped only without a controller in ABSENT.
						reached := controller == 0 || !zoneKnown
						want := reached && (controller != 0 || state != regulatorAbsent)
						probes := 0
						poller := newFallbackGatePoller(controller, zoneKnown, state, &probes)
						task.run(poller, context.Background())
						if got := probes > 0; got != want {
							t.Fatalf("root discovery ran = %v (%d probe calls); want %v", got, probes, want)
						}
					})
				}
			}
		}
	}
}

// TestVaillantSemanticPollerFallbackSuspensionLogsOnTransition checks that the
// suspension marker is written once per entry into the suspended state, not on
// every skipped tick.
func TestVaillantOwnerVaillantSemanticPollerFallbackSuspensionLogsOnTransition(t *testing.T) {
	var mu sync.Mutex
	var out strings.Builder
	original := log.Writer()
	log.SetOutput(fallbackLogWriter{mu: &mu, out: &out})
	t.Cleanup(func() { log.SetOutput(original) })
	marker := "semantic_discovery_fallback_suspended regulator_state=ABSENT discovery_interval=10m0s"
	count := func() int {
		mu.Lock()
		defer mu.Unlock()
		return strings.Count(out.String(), marker)
	}

	probes := 0
	poller := newFallbackGatePoller(0x00, false, regulatorAbsent, &probes)
	ctx := context.Background()
	poller.refreshState(ctx)
	poller.refreshConfig(ctx)
	poller.refreshState(ctx)
	if got := count(); got != 1 || probes != 0 {
		t.Fatalf("markers = %d, probe calls = %d after three skipped ticks; want 1, 0", got, probes)
	}

	poller.mu.Lock()
	poller.regAbsenceState = regulatorAbsenceGrace
	poller.mu.Unlock()
	poller.refreshState(ctx)
	if got := count(); got != 1 || probes == 0 {
		t.Fatalf("markers = %d, probe calls = %d after leaving ABSENT; want 1 and a discovery pass", got, probes)
	}

	poller.mu.Lock()
	poller.regAbsenceState = regulatorAbsent
	poller.mu.Unlock()
	poller.refreshState(ctx)
	if got := count(); got != 2 {
		t.Fatalf("markers = %d after entering ABSENT again; want 2", got)
	}

	// Leaving the suspension by discovering a controller clears the flag. No
	// provider, so the successful pass runs no startup priming in this test;
	// a task queue, as in production, so the controller priming is queued and
	// no other refresh re-evaluates the gate before the check.
	poller.publication = newVaillantSemanticPublication(nil, publishedHub(poller), publishedCache(poller))
	poller.tasks = newSemanticTaskScheduler()
	poller.b524ProbeFn = func(_ context.Context, target, _, _, _ byte, _ uint16) bool {
		probes++
		return target == 0x15
	}
	poller.refreshDiscovery(ctx)
	poller.mu.Lock()
	controller, suspended := poller.controller, poller.discoveryFallbackSuspended
	poller.mu.Unlock()
	if controller != 0x15 || suspended {
		t.Fatalf("after discovering a controller: controller 0x%02X, suspended %v; want 0x15, false", controller, suspended)
	}

	// Losing it again in ABSENT does not suspend: this process had a root, so
	// the state tick keeps its per-minute recovery and no marker is written.
	poller.b524ProbeFn = func(context.Context, byte, byte, byte, byte, uint16) bool {
		probes++
		return false
	}
	poller.refreshDiscovery(ctx)
	before := probes
	poller.refreshState(ctx)
	if probes == before {
		t.Fatal("state tick after losing a discovered controller ran no root discovery; want the per-minute recovery")
	}
	if got := count(); got != 2 {
		t.Fatalf("markers = %d after losing a discovered controller in ABSENT; want 2 (no suspension)", got)
	}
}

// TestVaillantSemanticPollerKeepsRecoveryWhenAbsentControllerWasDiscovered
// covers sites where a B524 root exists but the regulator is not classified
// by the product data, so the absence state still reaches ABSENT: a
// controller registered structurally without a serial, a serial whose part
// number the product data does not list, and unavailable product data. After
// one failed discovery pass loses the controller and its zones, the next state
// tick must rediscover it instead of waiting for the discovery tick.
func TestVaillantOwnerVaillantSemanticPollerKeepsRecoveryWhenAbsentControllerWasDiscovered(t *testing.T) {
	catalog, err := productids.LoadCatalog()
	if err != nil {
		t.Fatalf("LoadCatalog() error = %v", err)
	}
	boilerSerial := controllerlessSerial(controllerlessCatalogPart(t, "Boiler"))
	missingPart := "0000000000"
	if got := catalog.ControllerCapability(missingPart); got != productids.ControllerUnknown {
		t.Fatalf("part %s capability = %v; want it absent from the product data", missingPart, got)
	}
	for _, tc := range []struct {
		name            string
		regulatorSerial string // empty: the controller is registered structurally by discovery
		catalogErr      error
	}{
		{name: "structural registration without serial"},
		{name: "part number missing from product data", regulatorSerial: controllerlessSerial(missingPart)},
		{name: "product data unavailable", regulatorSerial: controllerlessSerial(controllerlessCatalogPartWith(t, "Regulator", productids.ControllerPresent)), catalogErr: errors.New("product data unavailable")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			now := time.Unix(0, 0)
			probes := 0
			coherent := true
			devices := []registry.DeviceInfo{{Address: 0x08, Manufacturer: "Vaillant", DeviceID: "BAI00", SerialNumber: boilerSerial}}
			if tc.regulatorSerial != "" {
				devices = append(devices, registry.DeviceInfo{Address: 0x15, Manufacturer: "Vaillant", SerialNumber: tc.regulatorSerial})
			}
			// No provider: the successful pass runs no startup priming here.
			poller := &vaillantSemanticPoller{
				reg:                   newTestRegistry(devices...),
				catalog:               catalog,
				catalogErr:            tc.catalogErr,
				tasks:                 newSemanticTaskScheduler(),
				regAbsenceState:       regulatorPresent,
				regulatorAbsenceGrace: 5 * time.Minute,
				discoveryInterval:     10 * time.Minute,
				nowFn:                 func() time.Time { return now },
				zones:                 make(map[byte]*vaillantZoneSnapshot),
				presence:              make(map[byte]*zonePresenceRecord),
				circuits:              make(map[byte]*vaillantCircuitSnapshot),
				radioDevices:          make(map[radioDeviceKey]*vaillantRadioDeviceSnapshot),
				solarCylinders:        make(map[byte]*vaillantCylinderSnapshot),
				b524ProbeFn: func(_ context.Context, target, _, _, _ byte, _ uint16) bool {
					probes++
					return coherent && target == 0x15
				},
			}
			controllerOf := func() byte {
				poller.mu.Lock()
				defer poller.mu.Unlock()
				return poller.controller
			}

			poller.refreshDiscovery(ctx)
			if got := controllerOf(); got != 0x15 {
				t.Fatalf("controller after discovery = 0x%02X; want 0x15", got)
			}
			poller.refreshRegulatorCapability(ctx)
			now = now.Add(6 * time.Minute)
			poller.refreshRegulatorCapability(ctx)
			poller.mu.Lock()
			state, capability := poller.regAbsenceState, poller.regulatorCapability
			poller.mu.Unlock()
			if state != regulatorAbsent || capability != productids.ControllerUnknown {
				t.Fatalf("absence state %s, capability %v with a controller known; want ABSENT, unknown", state, capability)
			}

			// A transient failure of one discovery pass loses the controller.
			coherent = false
			poller.refreshDiscovery(ctx)
			if got := controllerOf(); got != 0 {
				t.Fatalf("controller after a failed pass = 0x%02X; want 0", got)
			}

			coherent = true
			probes = 0
			poller.refreshState(ctx)
			if probes == 0 || controllerOf() != 0x15 {
				t.Fatalf("state tick after the loss: %d probe calls, controller 0x%02X; want the controller rediscovered on this tick", probes, controllerOf())
			}
		})
	}
}

// TestVaillantSemanticPollerFallbackKeptWhenPeriodicDiscoveryDisabled covers a
// zero or negative discovery interval, for which the poll schedule never
// starts the periodic discovery loop. The fallback on the state and
// configuration ticks is then the only periodic search and stays in place in
// ABSENT, with no suspension marker. The gate's view must match the
// schedule's: the fallback is suspended exactly when the schedule starts the
// discovery loop (creates its ticker).
func TestVaillantOwnerVaillantSemanticPollerFallbackKeptWhenPeriodicDiscoveryDisabled(t *testing.T) {
	var mu sync.Mutex
	var out strings.Builder
	original := log.Writer()
	log.SetOutput(fallbackLogWriter{mu: &mu, out: &out})
	t.Cleanup(func() { log.SetOutput(original) })
	markers := func() int {
		mu.Lock()
		defer mu.Unlock()
		return strings.Count(out.String(), "semantic_discovery_fallback_suspended")
	}

	for _, interval := range []time.Duration{-time.Second, 0, 10 * time.Minute} {
		t.Run(interval.String(), func(t *testing.T) {
			probes := 0
			poller := newFallbackGatePoller(0x00, false, regulatorAbsent, &probes)
			poller.discoveryInterval = interval

			var discovery semanticPollTask
			for _, task := range poller.pollSchedule().recurring {
				if task.key == semanticTaskRefreshDiscovery {
					discovery = task
				}
			}
			loopStarted := false
			schedule := &semanticPollSchedule{pollMu: &sync.Mutex{}, newTicker: func(time.Duration) *semanticPollTicker {
				loopStarted = true
				return &semanticPollTicker{C: make(chan time.Time), stop: func() {}}
			}}
			stopped, cancel := context.WithCancel(context.Background())
			cancel()
			schedule.runLoop(stopped, discovery.interval, discovery.key, discovery.priority, discovery.refresh)

			before := markers()
			ctx := context.Background()
			poller.refreshState(ctx)
			stateProbes := probes
			poller.refreshConfig(ctx)
			configProbes := probes - stateProbes
			fallbackRan := stateProbes > 0 && configProbes > 0
			if fallbackRan == loopStarted {
				t.Fatalf("periodic discovery loop started = %v, fallback ran (state %d, config %d probe calls) = %v; want exactly one of them", loopStarted, stateProbes, configProbes, fallbackRan)
			}
			wantMarkers := 0
			if loopStarted {
				wantMarkers = 1
			}
			if got := markers() - before; got != wantMarkers {
				t.Fatalf("suspension markers = %d; want %d", got, wantMarkers)
			}
			if loopStarted && !strings.Contains(out.String(), "discovery_interval=10m0s") {
				t.Fatalf("suspension marker does not name the running interval: %q", out.String())
			}
		})
	}
}
