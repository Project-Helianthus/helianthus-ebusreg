package owner

import (
	"context"
	"github.com/Project-Helianthus/helianthus-ebusgo/engine"
	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/output"
	"github.com/Project-Helianthus/helianthus-ebusreg/vaillant/productids"
	"math"
	"slices"
	"testing"
	"time"
)

func TestVaillantOwnerRefreshCircuitsNoSuccessfulReadsPreservesSnapshot(t *testing.T) {
	t.Parallel()

	circuitType := uint16(1)
	flow := 35.0
	cooling := uint16(1)
	poller := &vaillantSemanticPoller{
		controller: 0x15,
		circuits: map[byte]*vaillantCircuitSnapshot{
			0x00: {
				Instance:          0x00,
				Active:            true,
				CircuitTypeRaw:    &circuitType,
				FlowTemperatureC:  &flow,
				CoolingEnabledRaw: &cooling,
			},
		},
	}

	poller.refreshCircuits(context.Background())

	entry, ok := poller.circuits[0x00]
	if !ok || entry == nil {
		t.Fatal("circuit snapshot missing after refresh with no successful reads")
	}
	if entry.FlowTemperatureC == nil || *entry.FlowTemperatureC != 35.0 {
		t.Fatalf("entry.FlowTemperatureC = %v; want preserved 35.0", entry.FlowTemperatureC)
	}
	if entry.CircuitTypeRaw == nil || *entry.CircuitTypeRaw != 1 {
		t.Fatalf("entry.CircuitTypeRaw = %v; want preserved 1", entry.CircuitTypeRaw)
	}
	if entry.CoolingEnabledRaw == nil || *entry.CoolingEnabledRaw != 1 {
		t.Fatalf("entry.CoolingEnabledRaw = %v; want preserved 1", entry.CoolingEnabledRaw)
	}
}

func TestVaillantOwnerRefreshCircuitsPromotesDhwPseudoCircuitWithTemperatureEvidence(t *testing.T) {
	t.Parallel()

	now := time.Unix(1000, 0)
	existingType := uint16(1)
	provider := newRetainedSemanticState()
	poller := &vaillantSemanticPoller{
		controller:                  0x15,
		source:                      0x7F,
		publication:                 newVaillantSemanticPublication(provider, nil, nil),
		scheduler:                   engine.NewSemanticReadScheduler(),
		requestTimeout:              50 * time.Millisecond,
		circuitFullScanInterval:     semanticCircuitFullScanInterval,
		lastCircuitFullScanAt:       now.Add(-5 * time.Minute),
		lastCircuitFullScanComplete: true,
		nowFn:                       func() time.Time { return now },
		circuits: map[byte]*vaillantCircuitSnapshot{
			0x09: {
				Instance:       0x09,
				Active:         true,
				CircuitTypeRaw: &existingType,
			},
		},
		sendFrameFn: func(_ context.Context, frame protocol.Frame) (*protocol.Frame, error) {
			if len(frame.Data) != 6 {
				return testB524ResponseForSelectorPayload(frame.Data), nil
			}
			group := frame.Data[2]
			instance := frame.Data[3]
			addr := uint16(frame.Data[4]) | uint16(frame.Data[5])<<8
			if group != localCircuits.group || instance != 0x09 {
				return testB524ResponseForSelectorPayload(frame.Data), nil
			}
			switch addr {
			case circuitTypeRegister:
				return testB524ResponseForSelectorPayload(frame.Data, 0x00, 0x00), nil
			case circuitFlowTemp:
				return testB524ResponseForSelectorPayload(frame.Data, testB524Float32Payload(47.5)...), nil
			case circuitCalcFlowTemp:
				return testB524ResponseForSelectorPayload(frame.Data, testB524Float32Payload(46.25)...), nil
			case circuitFlowSetpoint,
				circuitHeatingCurve,
				circuitFlowTempMax,
				circuitFlowTempMin,
				circuitSummerLimit,
				circuitMixerPosition,
				circuitHumidity,
				circuitDewPoint:
				return testB524ResponseForSelectorPayload(frame.Data, testB524Float32Payload(0)...), nil
			case circuitPumpHours,
				circuitPumpStarts:
				return testB524ResponseForSelectorPayload(frame.Data, 0x00, 0x00, 0x00, 0x00), nil
			default:
				return testB524ResponseForSelectorPayload(frame.Data, 0x00), nil
			}
		},
	}

	poller.refreshCircuits(context.Background())

	entry, ok := poller.circuits[0x09]
	if !ok || entry == nil {
		t.Fatal("DHW pseudo-circuit snapshot missing after refresh")
	}
	if entry.CircuitTypeRaw == nil || *entry.CircuitTypeRaw != 3 {
		t.Fatalf("entry.CircuitTypeRaw = %v; want 3 (dhw)", entry.CircuitTypeRaw)
	}
	if entry.FlowTemperatureC == nil || *entry.FlowTemperatureC != 47.5 {
		t.Fatalf("entry.FlowTemperatureC = %v; want 47.5", entry.FlowTemperatureC)
	}

	status := provider.Circuits()
	if len(status) != 1 {
		t.Fatalf("provider.Circuits() = %d entries; want 1", len(status))
	}
	if status[0].Index != 9 || status[0].CircuitType != "dhw" {
		t.Fatalf("provider.Circuits()[0] = %#v; want index 9 type dhw", status[0])
	}
}

func TestVaillantOwnerRefreshCircuitsDoesNotPromoteDhwPseudoCircuitWithoutTemperatureEvidence(t *testing.T) {
	t.Parallel()

	now := time.Unix(1000, 0)
	existingType := uint16(1)
	poller := &vaillantSemanticPoller{
		controller:                  0x15,
		source:                      0x7F,
		publication:                 newVaillantSemanticPublication(newRetainedSemanticState(), nil, nil),
		scheduler:                   engine.NewSemanticReadScheduler(),
		requestTimeout:              50 * time.Millisecond,
		circuitFullScanInterval:     semanticCircuitFullScanInterval,
		lastCircuitFullScanAt:       now.Add(-5 * time.Minute),
		lastCircuitFullScanComplete: true,
		nowFn:                       func() time.Time { return now },
		circuits: map[byte]*vaillantCircuitSnapshot{
			0x09: {
				Instance:       0x09,
				Active:         true,
				CircuitTypeRaw: &existingType,
			},
		},
		sendFrameFn: func(_ context.Context, frame protocol.Frame) (*protocol.Frame, error) {
			if len(frame.Data) != 6 {
				return testB524ResponseForSelectorPayload(frame.Data), nil
			}
			addr := uint16(frame.Data[4]) | uint16(frame.Data[5])<<8
			switch addr {
			case circuitTypeRegister:
				return testB524ResponseForSelectorPayload(frame.Data, 0x00, 0x00), nil
			case circuitFlowTemp, circuitCalcFlowTemp:
				return testB524ResponseForSelectorPayload(frame.Data, testB524Float32Payload(math.NaN())...), nil
			default:
				return testB524ResponseForSelectorPayload(frame.Data, 0x00), nil
			}
		},
	}

	poller.refreshCircuits(context.Background())

	if entry := poller.circuits[0x09]; entry != nil {
		t.Fatalf("poller.circuits[0x09] = %#v; want removed without DHW temperature evidence", entry)
	}
}

func TestVaillantOwnerRefreshCircuitsPreservesDhwPseudoCircuitWhenEvidenceReadsFail(t *testing.T) {
	t.Parallel()

	now := time.Unix(1000, 0)
	existingType := uint16(3)
	existingFlow := 47.5
	poller := &vaillantSemanticPoller{
		controller:                  0x15,
		source:                      0x7F,
		publication:                 newVaillantSemanticPublication(newRetainedSemanticState(), nil, nil),
		scheduler:                   engine.NewSemanticReadScheduler(),
		requestTimeout:              50 * time.Millisecond,
		circuitFullScanInterval:     semanticCircuitFullScanInterval,
		lastCircuitFullScanAt:       now.Add(-5 * time.Minute),
		lastCircuitFullScanComplete: true,
		nowFn:                       func() time.Time { return now },
		circuits: map[byte]*vaillantCircuitSnapshot{
			0x09: {
				Instance:         0x09,
				Active:           true,
				CircuitTypeRaw:   &existingType,
				FlowTemperatureC: &existingFlow,
			},
		},
		sendFrameFn: func(_ context.Context, frame protocol.Frame) (*protocol.Frame, error) {
			if len(frame.Data) != 6 {
				return testB524ResponseForSelectorPayload(frame.Data), nil
			}
			addr := uint16(frame.Data[4]) | uint16(frame.Data[5])<<8
			switch addr {
			case circuitTypeRegister:
				return testB524ResponseForSelectorPayload(frame.Data, 0x00, 0x00), nil
			case circuitFlowTemp, circuitCalcFlowTemp:
				return &protocol.Frame{Data: []byte{0x00}}, nil
			default:
				return testB524ResponseForSelectorPayload(frame.Data, 0x00), nil
			}
		},
	}

	poller.refreshCircuits(context.Background())

	entry, ok := poller.circuits[0x09]
	if !ok || entry == nil {
		t.Fatal("DHW pseudo-circuit snapshot missing after transient evidence read failures")
	}
	if entry.CircuitTypeRaw == nil || *entry.CircuitTypeRaw != 3 {
		t.Fatalf("entry.CircuitTypeRaw = %v; want preserved 3", entry.CircuitTypeRaw)
	}
	if entry.FlowTemperatureC == nil || *entry.FlowTemperatureC != 47.5 {
		t.Fatalf("entry.FlowTemperatureC = %v; want preserved 47.5", entry.FlowTemperatureC)
	}
}

func TestVaillantOwnerRefreshCircuitsTargetedProbeFailureForcesNextFullScan(t *testing.T) {
	t.Parallel()

	now := time.Unix(1000, 0)
	circuitType := uint16(1)
	poller := &vaillantSemanticPoller{
		controller:                  0x15,
		circuitFullScanInterval:     semanticCircuitFullScanInterval,
		lastCircuitFullScanAt:       now.Add(-5 * time.Minute),
		lastCircuitFullScanComplete: true,
		nowFn:                       func() time.Time { return now },
		circuits: map[byte]*vaillantCircuitSnapshot{
			0x01: {
				Instance:       0x01,
				Active:         true,
				CircuitTypeRaw: &circuitType,
			},
		},
	}

	poller.refreshCircuits(context.Background())

	poller.mu.Lock()
	got, fullScan := poller.circuitRefreshInstancesLocked(now.Add(time.Second))
	poller.mu.Unlock()
	if !fullScan {
		t.Fatal("circuitRefreshInstancesLocked fullScan = false; want true after targeted probe failure")
	}
	if want := allCircuitRefreshInstances(); !slices.Equal(got, want) {
		t.Fatalf("circuitRefreshInstancesLocked instances = %#v; want %#v", got, want)
	}
}

func TestVaillantOwnerCircuitRefreshInstancesUsesKnownActiveCircuitsBeforeFullScanDue(t *testing.T) {
	t.Parallel()

	now := time.Unix(1000, 0)
	poller := &vaillantSemanticPoller{
		circuitFullScanInterval:     semanticCircuitFullScanInterval,
		lastCircuitFullScanAt:       now.Add(-5 * time.Minute),
		lastCircuitFullScanComplete: true,
		circuits: map[byte]*vaillantCircuitSnapshot{
			0x09: {Instance: 0x09, Active: true},
			0x01: {Instance: 0x01, Active: true},
			0x02: {Instance: 0x02, Active: false},
		},
	}

	got, fullScan := poller.circuitRefreshInstancesLocked(now)
	if fullScan {
		t.Fatal("circuitRefreshInstancesLocked fullScan = true; want false")
	}
	if want := []byte{0x01, 0x09}; !slices.Equal(got, want) {
		t.Fatalf("circuitRefreshInstancesLocked instances = %#v; want %#v", got, want)
	}
}

func TestVaillantOwnerCircuitRefreshInstancesFullScanAfterRediscoveryInterval(t *testing.T) {
	t.Parallel()

	now := time.Unix(1000, 0)
	poller := &vaillantSemanticPoller{
		circuitFullScanInterval: semanticCircuitFullScanInterval,
		lastCircuitFullScanAt:   now.Add(-semanticCircuitFullScanInterval),
		circuits: map[byte]*vaillantCircuitSnapshot{
			0x01: {Instance: 0x01, Active: true},
		},
	}

	got, fullScan := poller.circuitRefreshInstancesLocked(now)
	if !fullScan {
		t.Fatal("circuitRefreshInstancesLocked fullScan = false; want true")
	}
	if want := allCircuitRefreshInstances(); !slices.Equal(got, want) {
		t.Fatalf("circuitRefreshInstancesLocked instances = %#v; want %#v", got, want)
	}
}

func TestVaillantOwnerCircuitRefreshInstancesRetriesPartialFullScanOnConfigCadence(t *testing.T) {
	t.Parallel()

	now := time.Unix(1000, 0)
	poller := &vaillantSemanticPoller{
		circuitFullScanInterval:     semanticCircuitFullScanInterval,
		lastCircuitFullScanAt:       now.Add(-semanticCircuitPartialScanInterval),
		lastCircuitFullScanComplete: false,
		circuits: map[byte]*vaillantCircuitSnapshot{
			0x01: {Instance: 0x01, Active: true},
		},
	}

	got, fullScan := poller.circuitRefreshInstancesLocked(now)
	if !fullScan {
		t.Fatal("circuitRefreshInstancesLocked fullScan = false; want true after partial scan retry interval")
	}
	if want := allCircuitRefreshInstances(); !slices.Equal(got, want) {
		t.Fatalf("circuitRefreshInstancesLocked instances = %#v; want %#v", got, want)
	}
}

func TestVaillantOwnerDecodeHeatingCircuitTypeToken(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		raw  *uint16
		want string
	}{
		{name: "nil", raw: nil, want: ""},
		{name: "heating", raw: uint16Ptr(1), want: "heating"},
		{name: "fixed_value", raw: uint16Ptr(2), want: "fixed_value"},
		{name: "dhw", raw: uint16Ptr(3), want: "dhw"},
		{name: "return_increase", raw: uint16Ptr(4), want: "return_increase"},
		{name: "unknown", raw: uint16Ptr(9), want: "unknown_9"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			if got := DecodeHeatingCircuitTypeToken(test.raw); got != test.want {
				t.Fatalf("DecodeHeatingCircuitTypeToken(%v) = %q; want %q", test.raw, got, test.want)
			}
		})
	}
}

func TestVaillantOwnerDecodeRoomTempControlToken(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		raw  *uint16
		want string
	}{
		{name: "nil", raw: nil, want: ""},
		{name: "off", raw: uint16Ptr(0), want: "off"},
		{name: "modulating", raw: uint16Ptr(1), want: "modulating"},
		{name: "thermostat", raw: uint16Ptr(2), want: "thermostat"},
		{name: "unknown", raw: uint16Ptr(7), want: "unknown_7"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			if got := DecodeRoomTempControlToken(test.raw); got != test.want {
				t.Fatalf("DecodeRoomTempControlToken(%v) = %q; want %q", test.raw, got, test.want)
			}
		})
	}
}

func TestVaillantOwnerCircuitStateDecodesKnownValues(t *testing.T) {
	tests := []struct {
		raw  uint16
		want string
	}{
		{0, "standby"},
		{1, "heating"},
		{2, "cooling"},
		{3, "dhw"},
	}
	for _, tt := range tests {
		raw := tt.raw
		got := DecodeCircuitStateToken(&raw)
		if got != tt.want {
			t.Errorf("DecodeCircuitStateToken(%d) = %q; want %q", tt.raw, got, tt.want)
		}
	}
}

func TestVaillantOwnerCircuitRegisterSelectorsPinCircuitStateAndPumpStatus(t *testing.T) {
	t.Parallel()

	if circuitCircuitState != uint16(0x001E) {
		t.Fatalf("circuit_circuit_state = 0x%04x; want 0x001E (Status circuit)", circuitCircuitState)
	}
	if circuitPumpStatus != uint16(0x001B) {
		t.Fatalf("circuit_pump_status = 0x%04x; want 0x001B (Status Pump)", circuitPumpStatus)
	}
}

// RR=0x001E carries the circuit state enumeration and RR=0x001B the raw pump
// status. Distinct payloads per register prove each field follows its own
// selector, including DHW (3) and a raw pump value that is not a state code.
func TestVaillantOwnerRefreshCircuitsReadsCircuitStateFromCircuitStateRegisterAndPumpFromPumpStatusRegister(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		stateRaw  byte
		pumpRaw   byte
		wantState string
		wantPump  bool
	}{
		{"standby_pump_on", 0, 1, "standby", true},
		{"heating_pump_off", 1, 0, "heating", false},
		{"cooling", 2, 1, "cooling", true},
		{"dhw", 3, 1, "dhw", true},
		{"unmapped_stays_raw", 7, 0, "unknown_7", false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			now := time.Unix(1000, 0)
			circuitType := uint16(1)
			provider := newRetainedSemanticState()
			poller := &vaillantSemanticPoller{
				controller:                  0x15,
				source:                      0x7F,
				publication:                 newVaillantSemanticPublication(provider, nil, nil),
				scheduler:                   engine.NewSemanticReadScheduler(),
				requestTimeout:              50 * time.Millisecond,
				circuitFullScanInterval:     semanticCircuitFullScanInterval,
				lastCircuitFullScanAt:       now.Add(-5 * time.Minute),
				lastCircuitFullScanComplete: true,
				nowFn:                       func() time.Time { return now },
				circuits: map[byte]*vaillantCircuitSnapshot{
					0x00: {Instance: 0x00, Active: true, CircuitTypeRaw: &circuitType},
				},
				sendFrameFn: func(_ context.Context, frame protocol.Frame) (*protocol.Frame, error) {
					if len(frame.Data) != 6 {
						return testB524ResponseForSelectorPayload(frame.Data), nil
					}
					group := frame.Data[2]
					instance := frame.Data[3]
					addr := uint16(frame.Data[4]) | uint16(frame.Data[5])<<8
					if group != localCircuits.group || instance != 0x00 {
						return testB524ResponseForSelectorPayload(frame.Data), nil
					}
					switch addr {
					case 0x001E:
						return testB524ResponseForSelectorPayload(frame.Data, tt.stateRaw, 0x00), nil
					case 0x001B:
						return testB524ResponseForSelectorPayload(frame.Data, tt.pumpRaw, 0x00), nil
					case circuitTypeRegister:
						return testB524ResponseForSelectorPayload(frame.Data, 0x01, 0x00), nil
					case circuitFlowTemp:
						return testB524ResponseForSelectorPayload(frame.Data, testB524Float32Payload(35)...), nil
					default:
						return testB524ResponseForSelectorPayload(frame.Data, 0x00), nil
					}
				},
			}

			poller.refreshCircuits(context.Background())

			entry := poller.circuits[0x00]
			if entry == nil || entry.CircuitStateRaw == nil || entry.PumpStatusRaw == nil {
				t.Fatalf("snapshot = %#v; want circuit state and pump status read", entry)
			}
			if *entry.CircuitStateRaw != uint16(tt.stateRaw) {
				t.Fatalf("CircuitStateRaw = %d; want %d (from RR=0x001E)", *entry.CircuitStateRaw, tt.stateRaw)
			}
			if *entry.PumpStatusRaw != uint16(tt.pumpRaw) {
				t.Fatalf("PumpStatusRaw = %d; want %d (from RR=0x001B)", *entry.PumpStatusRaw, tt.pumpRaw)
			}
			circuits := provider.Circuits()
			if len(circuits) != 1 {
				t.Fatalf("provider.Circuits() = %d entries; want 1", len(circuits))
			}
			if got := circuits[0].State.CircuitState; got != tt.wantState {
				t.Fatalf("published circuit_state = %q; want %q", got, tt.wantState)
			}
			if got := circuits[0].State.PumpActive; got == nil || *got != tt.wantPump {
				t.Fatalf("published pump_active = %v; want %v", got, tt.wantPump)
			}
		})
	}
}

func TestVaillantOwnerCircuitStateReturnsUnknownForUnmappedValue(t *testing.T) {
	raw := uint16(99)
	got := DecodeCircuitStateToken(&raw)
	want := "unknown_99"
	if got != want {
		t.Errorf("DecodeCircuitStateToken(%d) = %q; want %q", raw, got, want)
	}
}

func TestVaillantOwnerCircuitStateNilReturnsEmpty(t *testing.T) {
	got := DecodeCircuitStateToken(nil)
	if got != "" {
		t.Errorf("DecodeCircuitStateToken(nil) = %q; want empty", got)
	}
}

// A refresh that accepts no energy point leaves the energy plane null.
func TestVaillantOwnerRefreshEnergyNilControllerKeepsEnergyNull(t *testing.T) {
	t.Parallel()

	provider := newRetainedSemanticState()
	poller := &vaillantSemanticPoller{
		publication:         newVaillantSemanticPublication(provider, nil, nil),
		regulatorCapability: productids.ControllerPresent,
	}

	poller.refreshEnergy(context.Background())

	if totals := provider.EnergyTotals(); totals != nil {
		t.Fatalf("EnergyTotals() = %+v; want nil while no energy point was accepted", totals)
	}
}

func TestVaillantOwnerRefreshEnergyNoRegulatorKeepsEnergyNull(t *testing.T) {
	t.Parallel()

	provider := newRetainedSemanticState()
	provider.SetEnergyClock(func() time.Time { return time.Date(2026, time.February, 26, 12, 0, 0, 0, time.UTC) })
	poller := &vaillantSemanticPoller{
		publication:         newVaillantSemanticPublication(provider, nil, nil),
		controller:          0x15,
		regulatorCapability: productids.ControllerNone,
	}

	poller.refreshEnergy(context.Background())

	if totals := provider.EnergyTotals(); totals != nil {
		t.Fatalf("EnergyTotals() = %+v; want nil while no energy point was accepted", totals)
	}
}

func TestVaillantOwnerRefreshEnergyReadFailurePreservesLastKnown(t *testing.T) {
	t.Parallel()

	provider := newRetainedSemanticState()
	provider.SetEnergyClock(func() time.Time { return time.Date(2026, time.February, 26, 12, 0, 0, 0, time.UTC) })
	poller := &vaillantSemanticPoller{
		publication: newVaillantSemanticPublication(provider, nil, nil),
		controller:  0x15,
	}

	provider.ApplyEnergyFromRegister(provider.EnergyWindowKey("gas", "hot_water", "day", ""), 3.5)
	provider.ApplyEnergyFromRegister(provider.EnergyWindowKey("gas", "hot_water", "year", "current"), 240.0)

	before := provider.EnergyTotals()
	if before == nil {
		t.Fatal("EnergyTotals() before refresh = nil; want seeded totals")
	}

	poller.refreshEnergy(context.Background())

	after := provider.EnergyTotals()
	if after == nil {
		t.Fatal("EnergyTotals() after failed refresh = nil; want preserved last-known totals")
	}
	if after.Gas.DHW.Today != before.Gas.DHW.Today {
		t.Fatalf("Gas.DHW.Today after failed refresh = %v; want preserved %v", after.Gas.DHW.Today, before.Gas.DHW.Today)
	}
	if len(after.Gas.DHW.Yearly) != len(before.Gas.DHW.Yearly) || after.Gas.DHW.Yearly[1] != before.Gas.DHW.Yearly[1] {
		t.Fatalf("Gas.DHW.Yearly after failed refresh = %#v; want preserved %#v", after.Gas.DHW.Yearly, before.Gas.DHW.Yearly)
	}
}

func TestVaillantOwnerPublishEnergyTotalsPublishesProviderSnapshot(t *testing.T) {
	t.Parallel()

	provider := newRetainedSemanticState()
	provider.SetEnergyClock(func() time.Time { return time.Date(2026, time.February, 26, 12, 0, 0, 0, time.UTC) })
	if !provider.ApplyEnergyFromRegister(provider.EnergyWindowKey("gas", "climate", "day", ""), 4.25) {
		t.Fatal("ApplyEnergyFromRegister() = false; want true")
	}

	hub := newPublicationHub()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	subscription, err := hub.SubscribeEnergy(ctx)
	if err != nil {
		t.Fatalf("SubscribeEnergy() error = %v", err)
	}

	newVaillantSemanticPublication(provider, hub, nil).publishEnergyTotals()

	select {
	case raw := <-subscription:
		totals, ok := raw.(*output.EnergyTotals)
		if !ok {
			t.Fatalf("energy payload type = %T; want *output.EnergyTotals", raw)
		}
		if totals.Gas.Climate.Today != 4.25 {
			t.Fatalf("energy payload gas.climate.today = %v; want 4.25", totals.Gas.Climate.Today)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for energy publish")
	}
}

func TestVaillantOwnerPublishFunctionModuleSemanticSkipsConfigOnlyCylinderWithoutTemperature(t *testing.T) {
	t.Parallel()

	provider := newRetainedSemanticState()
	temp := 48.0
	maxSetpoint := 60.0
	chargeOffset := 5.0
	poller := &vaillantSemanticPoller{
		publication:        newVaillantSemanticPublication(provider, nil, nil),
		functionModuleMode: vaillantFunctionModuleModeInterpreted,
		solarCylinders: map[byte]*vaillantCylinderSnapshot{
			0x00: {
				Instance:         0x00,
				TemperatureC:     &temp,
				MaxSetpointC:     &maxSetpoint,
				ChargeHysteresis: &chargeOffset,
			},
			0x01: {
				Instance:     0x01,
				MaxSetpointC: &maxSetpoint,
			},
		},
	}

	poller.publishFunctionModuleSemantic(semanticSnapshotSourceLive)

	cylinders := provider.Cylinders()
	if len(cylinders) != 1 {
		t.Fatalf("provider.Cylinders() len = %d; want 1 live-evidenced cylinder", len(cylinders))
	}
	if cylinders[0].Index != 0 {
		t.Fatalf("provider.Cylinders()[0].Index = %d; want 0", cylinders[0].Index)
	}
	if cylinders[0].TemperatureC == nil || *cylinders[0].TemperatureC != temp {
		t.Fatalf("provider.Cylinders()[0].TemperatureC = %#v; want %v", cylinders[0].TemperatureC, temp)
	}
}

func TestVaillantOwnerPublishFunctionModuleSemanticDowngradeRehydratesCircuitOwnershipToUnknown(t *testing.T) {
	t.Parallel()

	circuitType := uint16(1)
	scheme := uint16(1)
	module := uint16(2)
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

	initial := provider.Circuits()
	if len(initial) != 1 || initial[0].ManagingDevice.Role != output.ManagingDeviceRoleFunctionModule {
		t.Fatalf("initial circuit ownership = %#v; want FUNCTION_MODULE before downgrade", initial)
	}

	poller.mu.Lock()
	poller.functionModuleMode = vaillantFunctionModuleModeGPIOOnly
	poller.mu.Unlock()
	poller.publishFunctionModuleSemantic(semanticSnapshotSourceLive)

	if provider.Solar() == nil {
		t.Fatal("provider.Solar() = nil after GPIO-only FM5 publish; want empty non-null plane")
	}
	if cylinders := provider.Cylinders(); cylinders == nil || len(cylinders) != 0 {
		t.Fatalf("provider.Cylinders() = %#v after GPIO-only FM5 publish; want empty non-null plane", cylinders)
	}

	status := provider.Circuits()
	if len(status) != 1 {
		t.Fatalf("provider.Circuits() = %d entries; want 1", len(status))
	}
	if status[0].ManagingDevice.Role != output.ManagingDeviceRoleUnknown {
		t.Fatalf("provider.Circuits()[0].ManagingDevice.Role = %q; want %q", status[0].ManagingDevice.Role, output.ManagingDeviceRoleUnknown)
	}
	if status[0].ManagingDevice.DeviceID != nil || status[0].ManagingDevice.Address != nil {
		t.Fatalf("provider.Circuits()[0].ManagingDevice = %#v; want nil identity on downgrade", status[0].ManagingDevice)
	}
}
