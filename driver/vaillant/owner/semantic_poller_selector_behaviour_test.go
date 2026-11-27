package owner

import (
	"context"
	"encoding/binary"
	"errors"
	"github.com/Project-Helianthus/helianthus-ebusgo/engine"
	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
	"math"
	"slices"
	"testing"
	"time"
)

func TestVaillantOwnerBuildB524ReadSelector(t *testing.T) {
	t.Parallel()

	got := buildB524ReadSelector(0x02, 0x03, 0x01, 0x001C)
	want := []byte{0x02, 0x00, 0x03, 0x01, 0x1C, 0x00}
	if len(got) != len(want) {
		t.Fatalf("len(got)=%d; want %d", len(got), len(want))
	}
	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("got[%d]=0x%02x; want 0x%02x", i, got[i], want[i])
		}
	}
}

func TestVaillantOwnerParseB524ReadPayload(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		payload  []byte
		opcode   byte
		group    byte
		instance byte
		addr     uint16
		want     []byte
		ok       bool
	}{
		{
			name:     "valid payload with float value",
			payload:  []byte{0x01, 0x03, 0x22, 0x00, 0x00, 0x00, 0x30, 0x41},
			opcode:   0x02,
			group:    0x03,
			instance: 0x01,
			addr:     0x0022,
			want:     []byte{0x00, 0x00, 0x30, 0x41},
			ok:       true,
		},
		{
			name:     "valid payload with reply instance tag",
			payload:  []byte{0x08, 0x01, 0x03, 0x22, 0x00, 0x00, 0x00, 0x30, 0x41},
			opcode:   0x02,
			group:    0x03,
			instance: 0x00,
			addr:     0x0022,
			want:     []byte{0x00, 0x00, 0x30, 0x41},
			ok:       true,
		},
		{
			name:     "rejects mismatched reply instance",
			payload:  []byte{0x08, 0x05, 0x03, 0x22, 0x00, 0x00, 0x00, 0x30, 0x41},
			opcode:   0x02,
			group:    0x03,
			instance: 0x00,
			addr:     0x0022,
			ok:       false,
		},
		{
			name:     "header only has no value",
			payload:  []byte{0x01, 0x03, 0x22, 0x00},
			opcode:   0x02,
			group:    0x03,
			instance: 0x01,
			addr:     0x0022,
			ok:       false,
		},
		{
			name:     "short payload",
			payload:  []byte{0x01, 0x03, 0x22},
			opcode:   0x02,
			group:    0x03,
			instance: 0x01,
			addr:     0x0022,
			ok:       false,
		},
		{
			name:     "single 0x00 has no value",
			payload:  []byte{0x00},
			opcode:   0x02,
			group:    0x03,
			instance: 0x01,
			addr:     0x0022,
			ok:       false,
		},
		{
			name:     "rejects mismatched group",
			payload:  []byte{0x01, 0x02, 0x22, 0x00, 0x01},
			opcode:   0x02,
			group:    0x03,
			instance: 0x01,
			addr:     0x0022,
			ok:       false,
		},
		{
			name:     "rejects mismatched addr",
			payload:  []byte{0x01, 0x03, 0x21, 0x00, 0x01},
			opcode:   0x02,
			group:    0x03,
			instance: 0x01,
			addr:     0x0022,
			ok:       false,
		},
		{
			name:     "allows kind 0x00 with value bytes",
			payload:  []byte{0x00, 0x03, 0x1c, 0x00, 0xff},
			opcode:   0x02,
			group:    0x03,
			instance: 0x02,
			addr:     0x001c,
			want:     []byte{0xff},
			ok:       true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got, ok := parseB524ReadPayload(test.payload, test.opcode, test.group, test.instance, test.addr)
			if ok != test.ok {
				t.Fatalf("ok = %v; want %v", ok, test.ok)
			}
			if !test.ok {
				return
			}
			if len(got) != len(test.want) {
				t.Fatalf("len(got) = %d; want %d", len(got), len(test.want))
			}
			for i := range got {
				if got[i] != test.want[i] {
					t.Fatalf("got[%d] = 0x%02x; want 0x%02x", i, got[i], test.want[i])
				}
			}
		})
	}
}

func TestVaillantOwnerReadB524StartupUsesLivePathWithoutScheduler(t *testing.T) {
	t.Parallel()

	const (
		opcode   = byte(0x02)
		group    = byte(0x00)
		instance = byte(0x00)
		addr     = uint16(0x0036)
	)

	calls := 0
	poller := &vaillantSemanticPoller{
		sendFrameFn: func(_ context.Context, frame protocol.Frame) (*protocol.Frame, error) {
			calls++
			if !slices.Equal(frame.Data, buildB524ReadSelector(opcode, group, instance, addr)) {
				t.Fatalf("request selector = % x; want B524 selector", frame.Data)
			}
			return &protocol.Frame{
				Data: []byte{0x01, instance, group, byte(addr), byte(addr >> 8), 0x02, 0x00},
			}, nil
		},
		source:         0x7F,
		controller:     0x15,
		requestTimeout: 50 * time.Millisecond,
	}

	got, ok := poller.readB524Startup(context.Background(), opcode, group, instance, addr)
	if !ok {
		t.Fatal("readB524Startup() ok = false; want live direct read success")
	}
	if !slices.Equal(got, []byte{0x02, 0x00}) {
		t.Fatalf("readB524Startup() = % x; want 02 00", got)
	}
	if calls != 1 {
		t.Fatalf("send calls = %d; want 1", calls)
	}
}

func TestVaillantOwnerRefreshStateSkipsSlowZoneAndHotWaterB524Selectors(t *testing.T) {
	t.Parallel()

	var selectors [][]byte
	poller := &vaillantSemanticPoller{
		scheduler:      engine.NewSemanticReadScheduler(),
		publication:    newVaillantSemanticPublication(newRetainedSemanticState(), nil, nil),
		source:         0x7F,
		controller:     0x15,
		requestTimeout: 50 * time.Millisecond,
		zones: map[byte]*vaillantZoneSnapshot{
			0x00: {
				Instance:                    0x00,
				Present:                     true,
				Preset:                      "manual",
				ConfigurationCircuitTypeRaw: uint16Ptr(1),
			},
		},
		circuits: map[byte]*vaillantCircuitSnapshot{
			0x00: {Instance: 0x00, Active: true, CircuitTypeRaw: uint16Ptr(1)},
		},
		dhw: &vaillantDhwSnapshot{Preset: "manual"},
	}
	poller.sendFrameFn = func(_ context.Context, frame protocol.Frame) (*protocol.Frame, error) {
		selectors = append(selectors, slices.Clone(frame.Data))
		return testB524ResponseForSelector(frame.Data), nil
	}

	poller.refreshState(context.Background())

	requireB524Selector(t, selectors, localZones.group, 0x00, zoneCurrentTemp)
	requireB524Selector(t, selectors, localZones.group, 0x00, zoneSpecialFunction)
	requireB524Selector(t, selectors, localZones.group, 0x00, zoneValveStatus)
	requireB524Selector(t, selectors, localDHW.group, dhwInstance, dhwCurrentTemp)
	requireB524Selector(t, selectors, localDHW.group, dhwInstance, dhwSpecialFunction)

	for _, forbidden := range []struct {
		group    byte
		instance byte
		addr     uint16
	}{
		{localZones.group, 0x00, zoneTargetTemp},
		{localZones.group, 0x00, zoneFallbackManualTemp},
		{localZones.group, 0x00, zoneCurrentHumidity},
		{localZones.group, 0x00, zoneHeatingOpMode},
		{localZones.group, 0x00, zoneQuickVetoTemp},
		{localZones.group, 0x00, zoneQuickVetoEndTime},
		{localZones.group, 0x00, zoneHolidayStartDate},
		{localZones.group, 0x00, zoneRoomTemperatureZoneMappingRawRegister},
		{localCircuits.group, 0x00, circuitTypeRegister},
		{localDHW.group, dhwInstance, dhwTargetTemp},
		{localDHW.group, dhwInstance, dhwOperationMode},
		{localDHW.group, dhwInstance, dhwHolidayStartDateRegister},
		{localDHW.group, dhwInstance, dhwHolidayEndDateRegister},
	} {
		if hasB524Selector(selectors, forbidden.group, forbidden.instance, forbidden.addr) {
			t.Fatalf("refreshState read slow selector group=0x%02x instance=0x%02x addr=0x%04x", forbidden.group, forbidden.instance, forbidden.addr)
		}
	}

	if got := poller.zones[0x00].Preset; got != "manual" {
		t.Fatalf("zone preset = %q; want preserved manual without cached op-mode", got)
	}
	if got := poller.dhw.Preset; got != "manual" {
		t.Fatalf("DHW preset = %q; want preserved manual without cached op-mode", got)
	}
}

func TestVaillantOwnerRefreshConfigReadsSlowZoneAndHotWaterB524Selectors(t *testing.T) {
	t.Parallel()

	var selectors [][]byte
	poller := &vaillantSemanticPoller{
		scheduler:      engine.NewSemanticReadScheduler(),
		publication:    newVaillantSemanticPublication(newRetainedSemanticState(), nil, nil),
		source:         0x7F,
		controller:     0x15,
		requestTimeout: 50 * time.Millisecond,
		zones: map[byte]*vaillantZoneSnapshot{
			0x00: {Instance: 0x00, Present: true},
		},
		circuits: make(map[byte]*vaillantCircuitSnapshot),
	}
	poller.sendFrameFn = func(_ context.Context, frame protocol.Frame) (*protocol.Frame, error) {
		selectors = append(selectors, slices.Clone(frame.Data))
		return testB524ResponseForSelector(frame.Data), nil
	}

	poller.refreshConfig(context.Background())

	requireB524Selector(t, selectors, localZones.group, 0x00, zoneTargetTemp)
	requireB524Selector(t, selectors, localZones.group, 0x00, zoneCurrentHumidity)
	requireB524Selector(t, selectors, localZones.group, 0x00, zoneHeatingOpMode)
	requireB524Selector(t, selectors, localZones.group, 0x00, zoneQuickVetoTemp)
	requireB524Selector(t, selectors, localZones.group, 0x00, zoneQuickVetoEndTime)
	requireB524Selector(t, selectors, localZones.group, 0x00, zoneHolidayStartDate)
	requireB524Selector(t, selectors, localZones.group, 0x00, zoneRoomTemperatureZoneMappingRawRegister)
	requireB524Selector(t, selectors, localCircuits.group, 0x00, circuitTypeRegister)
	requireB524Selector(t, selectors, localDHW.group, dhwInstance, dhwTargetTemp)
	requireB524Selector(t, selectors, localDHW.group, dhwInstance, dhwOperationMode)
	requireB524Selector(t, selectors, localDHW.group, dhwInstance, dhwHolidayStartDateRegister)
	requireB524Selector(t, selectors, localDHW.group, dhwInstance, dhwHolidayEndDateRegister)

	if hasB524Selector(selectors, localZones.group, 0x00, zoneCurrentTemp) {
		t.Fatal("refreshConfig read fast zone current temperature selector")
	}
	if hasB524Selector(selectors, localZones.group, 0x00, zoneSpecialFunction) {
		t.Fatal("refreshConfig read fast zone special function selector")
	}
	if hasB524Selector(selectors, localZones.group, 0x00, zoneValveStatus) {
		t.Fatal("refreshConfig read fast zone valve status selector")
	}
	if hasB524Selector(selectors, localDHW.group, dhwInstance, dhwCurrentTemp) {
		t.Fatal("refreshConfig read fast DHW current temperature selector")
	}
	if hasB524Selector(selectors, localDHW.group, dhwInstance, dhwSpecialFunction) {
		t.Fatal("refreshConfig read fast DHW special function selector")
	}
}

func TestVaillantOwnerRefreshBoilerStatusFastProjectsB524MirrorsWithoutActiveB524Reads(t *testing.T) {
	t.Parallel()

	var selectors [][]byte
	observedAt := time.Unix(100, 0)
	flowTemp := 42.5
	dhwTemp := 48.5
	pumpRaw := uint16(1)
	heatingRaw := uint16(2)
	poller := &vaillantSemanticPoller{
		scheduler:      engine.NewSemanticReadScheduler(),
		source:         0x7F,
		controller:     0x15,
		requestTimeout: 50 * time.Millisecond,
		nowFn:          func() time.Time { return observedAt.Add(30 * time.Second) },
		system: &vaillantSystemSnapshot{
			Controller:                  0x15,
			SystemFlowTemperature:       &flowTemp,
			SystemFlowTemperatureLiveAt: observedAt,
		},
		circuits: map[byte]*vaillantCircuitSnapshot{
			0x00: {
				Instance:           0x00,
				Active:             true,
				Controller:         0x15,
				PumpStatusRaw:      &pumpRaw,
				PumpStatusLiveAt:   observedAt,
				CircuitStateRaw:    &heatingRaw,
				CircuitStateLiveAt: observedAt,
			},
		},
		dhw: &vaillantDhwSnapshot{CurrentTempC: &dhwTemp},
	}
	poller.sendFrameFn = func(_ context.Context, frame protocol.Frame) (*protocol.Frame, error) {
		selectors = append(selectors, slices.Clone(frame.Data))
		return testB524ResponseForSelector(frame.Data), nil
	}

	snapshot := &vaillantBoilerSnapshot{}
	if !poller.refreshBoilerStatusB524(context.Background(), boilerStatusTierFast, snapshot) {
		t.Fatal("refreshBoilerStatusB524(fast) updated = false; want projected snapshot update")
	}
	if len(selectors) != 0 {
		t.Fatalf("refreshBoilerStatusB524(fast) performed %d active reads; want zero", len(selectors))
	}
	if snapshot.FlowTemperatureC != nil {
		t.Fatalf("flow temperature = %#v; want withheld without qualified boiler evidence", snapshot.FlowTemperatureC)
	}
	if got := poller.system.SystemFlowTemperature; got == nil || *got != flowTemp {
		t.Fatalf("native system flow = %v; want preserved %.1f", got, flowTemp)
	}
	if snapshot.CentralHeatingPumpActive != nil {
		t.Fatalf("boiler pump active = %#v; want unknown without independent boiler evidence", snapshot.CentralHeatingPumpActive)
	}
	if got := poller.circuits[0].PumpStatusRaw; got == nil || *got != pumpRaw {
		t.Fatalf("circuit pump = %v; want preserved %d", got, pumpRaw)
	}
	if snapshot.HeatingStatusRaw != nil {
		t.Fatalf("boiler heating status = %d; want unset without a boiler-qualified source", *snapshot.HeatingStatusRaw)
	}
	if got := poller.circuits[0].CircuitStateRaw; got == nil || *got != heatingRaw {
		t.Fatalf("circuit state = %v; want preserved %d", got, heatingRaw)
	}
	if snapshot.DhwTemperatureC == nil || *snapshot.DhwTemperatureC != dhwTemp {
		t.Fatalf("dhw temperature = %#v; want %.1f", snapshot.DhwTemperatureC, dhwTemp)
	}
}

func TestVaillantOwnerRefreshBoilerStatusFastIgnoresStaleB524Mirrors(t *testing.T) {
	t.Parallel()

	observedAt := time.Unix(100, 0)
	flowTemp := 42.5
	pumpRaw := uint16(1)
	heatingRaw := uint16(2)
	poller := &vaillantSemanticPoller{
		scheduler:          engine.NewSemanticReadScheduler(),
		publication:        newVaillantSemanticPublication(newRetainedSemanticState(), nil, nil),
		source:             0x7F,
		controller:         0x15,
		configInterval:     5 * time.Minute,
		boilerFastInterval: 30 * time.Second,
		requestTimeout:     50 * time.Millisecond,
		nowFn:              func() time.Time { return observedAt.Add(6 * time.Minute) },
		system: &vaillantSystemSnapshot{
			Controller:                  0x15,
			SystemFlowTemperature:       &flowTemp,
			SystemFlowTemperatureLiveAt: observedAt,
		},
		circuits: map[byte]*vaillantCircuitSnapshot{
			0x00: {
				Instance:           0x00,
				Active:             true,
				Controller:         0x15,
				PumpStatusRaw:      &pumpRaw,
				PumpStatusLiveAt:   observedAt,
				CircuitStateRaw:    &heatingRaw,
				CircuitStateLiveAt: observedAt,
			},
		},
	}

	snapshot := &vaillantBoilerSnapshot{}
	if poller.refreshBoilerStatusB524(context.Background(), boilerStatusTierFast, snapshot) {
		t.Fatal("refreshBoilerStatusB524(fast) updated = true; want stale mirrors ignored")
	}
	if snapshot.FlowTemperatureC != nil || snapshot.CentralHeatingPumpActive != nil || snapshot.HeatingStatusRaw != nil {
		t.Fatalf("stale mirror snapshot = %+v; want no projected fields", snapshot)
	}
}

func TestVaillantOwnerRefreshBoilerStatusFastIgnoresDifferentControllerB524Mirrors(t *testing.T) {
	t.Parallel()

	observedAt := time.Unix(100, 0)
	flowTemp := 42.5
	pumpRaw := uint16(1)
	heatingRaw := uint16(2)
	poller := &vaillantSemanticPoller{
		scheduler:          engine.NewSemanticReadScheduler(),
		publication:        newVaillantSemanticPublication(newRetainedSemanticState(), nil, nil),
		source:             0x7F,
		controller:         0x26,
		configInterval:     5 * time.Minute,
		boilerFastInterval: 30 * time.Second,
		requestTimeout:     50 * time.Millisecond,
		nowFn:              func() time.Time { return observedAt.Add(30 * time.Second) },
		system: &vaillantSystemSnapshot{
			Controller:                  0x15,
			SystemFlowTemperature:       &flowTemp,
			SystemFlowTemperatureLiveAt: observedAt,
		},
		circuits: map[byte]*vaillantCircuitSnapshot{
			0x00: {
				Instance:           0x00,
				Active:             true,
				Controller:         0x15,
				PumpStatusRaw:      &pumpRaw,
				PumpStatusLiveAt:   observedAt,
				CircuitStateRaw:    &heatingRaw,
				CircuitStateLiveAt: observedAt,
			},
		},
	}

	snapshot := &vaillantBoilerSnapshot{}
	if poller.refreshBoilerStatusB524(context.Background(), boilerStatusTierFast, snapshot) {
		t.Fatal("refreshBoilerStatusB524(fast) updated = true; want different-controller mirrors ignored")
	}
	if snapshot.FlowTemperatureC != nil || snapshot.CentralHeatingPumpActive != nil || snapshot.HeatingStatusRaw != nil {
		t.Fatalf("different-controller mirror snapshot = %+v; want no projected fields", snapshot)
	}
}

func TestVaillantOwnerRefreshBoilerStatusFastIgnoresRestampedPartialMirrorSnapshots(t *testing.T) {
	t.Parallel()

	observedAt := time.Unix(100, 0)
	flowTemp := 42.5
	pressure := 1.4
	pumpRaw := uint16(1)
	heatingRaw := uint16(2)
	circuitType := uint16(1)

	system := mergeSystemSnapshotNonDestructive(
		&vaillantSystemSnapshot{
			Controller:                  0x15,
			SystemFlowTemperature:       &flowTemp,
			SystemFlowTemperatureLiveAt: observedAt,
		},
		&vaillantSystemSnapshot{
			Controller:          0x26,
			SystemWaterPressure: &pressure,
		},
	)
	circuit := mergeCircuitSnapshotNonDestructive(
		&vaillantCircuitSnapshot{
			Instance:           0x00,
			Active:             true,
			Controller:         0x15,
			PumpStatusRaw:      &pumpRaw,
			PumpStatusLiveAt:   observedAt,
			CircuitStateRaw:    &heatingRaw,
			CircuitStateLiveAt: observedAt,
		},
		&vaillantCircuitSnapshot{
			Instance:       0x00,
			Active:         true,
			Controller:     0x26,
			CircuitTypeRaw: &circuitType,
		},
	)
	if system.SystemFlowTemperature != nil || !system.SystemFlowTemperatureLiveAt.IsZero() {
		t.Fatalf("merged system mirror = (%v, %s); want cleared on controller change", system.SystemFlowTemperature, system.SystemFlowTemperatureLiveAt)
	}
	if circuit.PumpStatusRaw != nil || !circuit.PumpStatusLiveAt.IsZero() {
		t.Fatalf("merged pump mirror = (%v, %s); want cleared on controller change", circuit.PumpStatusRaw, circuit.PumpStatusLiveAt)
	}
	if circuit.CircuitStateRaw != nil || !circuit.CircuitStateLiveAt.IsZero() {
		t.Fatalf("merged state mirror = (%v, %s); want cleared on controller change", circuit.CircuitStateRaw, circuit.CircuitStateLiveAt)
	}

	poller := &vaillantSemanticPoller{
		scheduler:          engine.NewSemanticReadScheduler(),
		publication:        newVaillantSemanticPublication(newRetainedSemanticState(), nil, nil),
		source:             0x7F,
		controller:         0x26,
		configInterval:     5 * time.Minute,
		boilerFastInterval: 30 * time.Second,
		requestTimeout:     50 * time.Millisecond,
		nowFn:              func() time.Time { return observedAt.Add(30 * time.Second) },
		system:             system,
		circuits: map[byte]*vaillantCircuitSnapshot{
			0x00: circuit,
		},
	}

	snapshot := &vaillantBoilerSnapshot{}
	if poller.refreshBoilerStatusB524(context.Background(), boilerStatusTierFast, snapshot) {
		t.Fatal("refreshBoilerStatusB524(fast) updated = true; want partial restamped mirrors ignored")
	}
}

func TestVaillantOwnerRefreshConfigDHWOnlySuccessKeepsZonesCacheSource(t *testing.T) {
	t.Parallel()

	provider := newRetainedSemanticState()
	poller := &vaillantSemanticPoller{
		scheduler:      engine.NewSemanticReadScheduler(),
		publication:    newVaillantSemanticPublication(provider, nil, nil),
		source:         0x7F,
		controller:     0x15,
		requestTimeout: 50 * time.Millisecond,
		zones: map[byte]*vaillantZoneSnapshot{
			0x00: {Instance: 0x00, Present: true, Name: "Zone 1"},
		},
		circuits: make(map[byte]*vaillantCircuitSnapshot),
	}
	poller.sendFrameFn = func(_ context.Context, frame protocol.Frame) (*protocol.Frame, error) {
		if len(frame.Data) == 6 && frame.Data[2] == localDHW.group {
			return testB524ResponseForSelector(frame.Data), nil
		}
		return nil, errors.New("zone read unavailable")
	}

	poller.refreshConfig(context.Background())

	if dhw := provider.DHW(); dhw == nil {
		t.Fatal("DHW = nil; want DHW slow config to publish")
	}
	if cacheEpoch, liveEpoch := provider.StartupEpochs(); liveEpoch != 1 || cacheEpoch != 1 {
		t.Fatalf("StartupEpochs() = cache=%d live=%d; want zone cache=1 and DHW live=1", cacheEpoch, liveEpoch)
	}
}

func TestVaillantOwnerRefreshConfigNoZonesStillRefreshesDHWSlowConfig(t *testing.T) {
	t.Parallel()

	provider := newRetainedSemanticState()
	poller := &vaillantSemanticPoller{
		scheduler:             engine.NewSemanticReadScheduler(),
		publication:           newVaillantSemanticPublication(provider, nil, nil),
		reg:                   newTestRegistry(registry.DeviceInfo{Address: 0x15, Manufacturer: "Vaillant", DeviceID: "BASV"}),
		source:                0x7F,
		controller:            0x15,
		requestTimeout:        50 * time.Millisecond,
		zones:                 make(map[byte]*vaillantZoneSnapshot),
		presence:              make(map[byte]*zonePresenceRecord),
		circuits:              make(map[byte]*vaillantCircuitSnapshot),
		radioDevices:          make(map[radioDeviceKey]*vaillantRadioDeviceSnapshot),
		solarCylinders:        make(map[byte]*vaillantCylinderSnapshot),
		startupSemanticPrimed: true,
		b524ProbeFn: func(context.Context, byte, byte, byte, byte, uint16) bool {
			return true
		},
	}
	poller.sendFrameFn = func(_ context.Context, frame protocol.Frame) (*protocol.Frame, error) {
		if len(frame.Data) == 6 && frame.Data[2] == localDHW.group {
			return testB524ResponseForSelector(frame.Data), nil
		}
		if len(frame.Data) == 6 && frame.Data[2] == localZones.group && (uint16(frame.Data[4])|uint16(frame.Data[5])<<8) == zoneIndex {
			group := frame.Data[2]
			instance := frame.Data[3]
			addr := uint16(frame.Data[4]) | uint16(frame.Data[5])<<8
			return &protocol.Frame{Data: []byte{0x01, instance, group, byte(addr), byte(addr >> 8), 0xFF}}, nil
		}
		return nil, errors.New("non-DHW read unavailable")
	}

	poller.refreshConfig(context.Background())

	dhw := provider.DHW()
	if dhw == nil {
		t.Fatal("DHW = nil; want no-zone config refresh to publish DHW slow fields")
	}
	if dhw.Config.HolidayStartDate == "" || dhw.Config.HolidayEndDate == "" {
		t.Fatalf("DHW holiday dates = %q/%q; want slow config values", dhw.Config.HolidayStartDate, dhw.Config.HolidayEndDate)
	}
}

func TestVaillantOwnerRefreshDiscoveryPrimesDHWBeforeStartupZoneScan(t *testing.T) {
	if testing.Short() {
		t.Skip("measured long-running coverage; exercised by full CI")
	}
	t.Parallel()

	var selectors [][]byte
	poller := &vaillantSemanticPoller{
		publication:    newVaillantSemanticPublication(newRetainedSemanticState(), nil, nil),
		reg:            newTestRegistry(registry.DeviceInfo{Address: 0x15, Manufacturer: "Vaillant", DeviceID: "BASV"}),
		tasks:          newSemanticTaskScheduler(),
		zones:          make(map[byte]*vaillantZoneSnapshot),
		presence:       make(map[byte]*zonePresenceRecord),
		circuits:       make(map[byte]*vaillantCircuitSnapshot),
		radioDevices:   make(map[radioDeviceKey]*vaillantRadioDeviceSnapshot),
		solarCylinders: make(map[byte]*vaillantCylinderSnapshot),
		source:         0x7F,
		requestTimeout: 50 * time.Millisecond,
		b524ProbeFn: func(ctx context.Context, target, opcode, group, instance byte, addr uint16) bool {
			return target == 0x15
		},
	}
	poller.sendFrameFn = func(_ context.Context, frame protocol.Frame) (*protocol.Frame, error) {
		if frame.Primary != vaillantExtRegisterPrimary || frame.Secondary != vaillantExtRegisterSecondary || len(frame.Data) != 6 {
			return &protocol.Frame{Data: []byte{0x00}}, nil
		}
		selector := slices.Clone(frame.Data)
		selectors = append(selectors, selector)
		group := selector[2]
		instance := selector[3]
		addr := uint16(selector[4]) | uint16(selector[5])<<8
		payload := []byte{0x00}
		switch {
		case group == localDHW.group && instance == dhwInstance && addr == dhwCurrentTemp:
			payload = make([]byte, 4)
			binary.LittleEndian.PutUint32(payload, math.Float32bits(48.5))
		case group == localDHW.group && instance == dhwInstance && addr == dhwOperationMode:
			payload = []byte{0x01, 0x00}
		case group == localZones.group && addr == zoneIndex:
			payload = []byte{0xFF}
		}
		data := []byte{0x01, instance, group, byte(addr), byte(addr >> 8)}
		data = append(data, payload...)
		return &protocol.Frame{Data: data}, nil
	}

	poller.refreshDiscovery(context.Background())

	if len(selectors) == 0 {
		t.Fatal("refreshDiscovery sent no B524 startup reads")
	}
	first := selectors[0]
	if first[2] != localDHW.group || first[3] != dhwInstance {
		t.Fatalf("first startup selector = % x; want DHW before zone scan", first)
	}
	if dhw := publishedProvider(poller).DHW(); dhw == nil {
		t.Fatal("provider.DHW() = nil; want critical DHW singleton primed during discovery")
	}
}
