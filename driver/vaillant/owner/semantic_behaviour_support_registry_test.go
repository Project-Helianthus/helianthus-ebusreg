package owner

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"expvar"
	"fmt"
	"github.com/Project-Helianthus/helianthus-ebusgo/engine"
	ebuserrors "github.com/Project-Helianthus/helianthus-ebusgo/errors"
	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/output"
	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
	"github.com/Project-Helianthus/helianthus-ebusreg/vaillant/productids"
	"log"
	"strings"
	"sync"
	"testing"
	"time"
)

// mockB524Probe returns a probe function that records calls and responds
// based on a coherence map keyed by target address.
func mockB524Probe(coherent map[byte]bool, probed *[]byte) func(ctx context.Context, target, opcode, group, instance byte, addr uint16) bool {
	return func(ctx context.Context, target, opcode, group, instance byte, addr uint16) bool {
		if probed != nil {
			found := false
			for _, p := range *probed {
				if p == target {
					found = true
					break
				}
			}
			if !found {
				*probed = append(*probed, target)
			}
		}
		return coherent[target]
	}
}

func newTestPoller(reg *registry.DeviceRegistry) *vaillantSemanticPoller {
	return &vaillantSemanticPoller{
		reg:            reg,
		source:         0x71,
		requestTimeout: 2 * time.Second,
		nowFn:          time.Now,
	}
}

func newFallbackGatePoller(controller byte, zoneKnown bool, state regulatorAbsenceState, probes *int) *vaillantSemanticPoller {
	poller := &vaillantSemanticPoller{
		reg:               newTestRegistry(registry.DeviceInfo{Address: 0x08, Manufacturer: "Vaillant", DeviceID: "BAI00"}),
		publication:       newVaillantSemanticPublication(newRetainedSemanticState(), nil, nil),
		controller:        controller,
		regAbsenceState:   state,
		discoveryInterval: 10 * time.Minute,
		zones:             make(map[byte]*vaillantZoneSnapshot),
		presence:          make(map[byte]*zonePresenceRecord),
		circuits:          make(map[byte]*vaillantCircuitSnapshot),
		radioDevices:      make(map[radioDeviceKey]*vaillantRadioDeviceSnapshot),
		solarCylinders:    make(map[byte]*vaillantCylinderSnapshot),
		b524ProbeFn: func(context.Context, byte, byte, byte, byte, uint16) bool {
			*probes++
			return false
		},
	}
	if zoneKnown {
		poller.zones[0x00] = &vaillantZoneSnapshot{}
	}
	return poller
}

type fallbackLogWriter struct {
	mu  *sync.Mutex
	out *strings.Builder
}

func (w fallbackLogWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.out.Write(p)
}

func newFunctionModuleDeadlineRecoveryTestPoller(tasks *semanticTaskScheduler) (*vaillantSemanticPoller, *retainedSemanticState) {
	now := time.Date(2026, 8, 15, 21, 0, 0, 0, time.UTC)
	provider := newRetainedSemanticState()
	return &vaillantSemanticPoller{
		scheduler:                 engine.NewSemanticReadScheduler(),
		tasks:                     tasks,
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
	}, provider
}

func functionModuleInstanceOneSemanticTestResponse(frame protocol.Frame) (*protocol.Frame, error) {
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
			if instance == 0x01 {
				return testB524ResponseForSelectorPayload(frame.Data, 0x01), nil
			}
			return testB524ResponseForSelectorPayload(frame.Data, 0x00), nil
		case deviceSlotClassAddress:
			if instance == 0x01 {
				return testB524ResponseForSelectorPayload(frame.Data, circuitManagingDeviceFunctionModuleControllerAddress), nil
			}
			return testB524ResponseForSelectorPayload(frame.Data, 0xFF), nil
		case deviceSlotFirmware:
			if instance == 0x01 {
				return testB524ResponseForSelectorPayload(frame.Data, 0x01, 0x00, 0x00), nil
			}
			return testB524ResponseForSelectorPayload(frame.Data, 0xFF, 0xFF, 0xFF), nil
		case deviceSlotHardwareIdentifier:
			if instance == 0x01 {
				return testB524ResponseForSelectorPayload(frame.Data, 0x71, 0x00), nil
			}
			return testB524ResponseForSelectorPayload(frame.Data, 0xFF, 0xFF), nil
		}
	}
	if group == localSolar.group && (addr == solarEnabled || addr == solarPumpActive) {
		return testB524ResponseForSelectorPayload(frame.Data, 0x01), nil
	}
	if group == localSolar.group || group == localCylinders.group {
		return testB524ResponseForSelectorPayload(frame.Data, testB524Float32Payload(47.25)...), nil
	}
	return testB524ResponseForSelectorPayload(frame.Data, 0x00, 0x00, 0x00, 0x00), nil
}

func functionModuleSemanticTestResponseForClasses(frame protocol.Frame, classes map[byte]byte) (*protocol.Frame, error) {
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
		classAddress, populated := classes[instance]
		switch addr {
		case deviceSlotConnected:
			if populated {
				return testB524ResponseForSelectorPayload(frame.Data, 0x01), nil
			}
			return testB524ResponseForSelectorPayload(frame.Data, 0x00), nil
		case deviceSlotClassAddress:
			if populated {
				return testB524ResponseForSelectorPayload(frame.Data, classAddress), nil
			}
			return testB524ResponseForSelectorPayload(frame.Data, 0xFF), nil
		case deviceSlotFirmware:
			if populated {
				return testB524ResponseForSelectorPayload(frame.Data, 0x01, 0x00, 0x00), nil
			}
			return testB524ResponseForSelectorPayload(frame.Data, 0xFF, 0xFF, 0xFF), nil
		case deviceSlotHardwareIdentifier:
			if populated {
				return testB524ResponseForSelectorPayload(frame.Data, 0x71, 0x00), nil
			}
			return testB524ResponseForSelectorPayload(frame.Data, 0xFF, 0xFF), nil
		}
	}
	if group == localSolar.group && (addr == solarEnabled || addr == solarPumpActive) {
		return testB524ResponseForSelectorPayload(frame.Data, 0x01), nil
	}
	if group == localSolar.group || group == localCylinders.group {
		return testB524ResponseForSelectorPayload(frame.Data, testB524Float32Payload(47.25)...), nil
	}
	return testB524ResponseForSelectorPayload(frame.Data, 0x00, 0x00, 0x00, 0x00), nil
}

const functionModuleHighSlotSelectorBudget = 28

func functionModuleHighSlotBudgetedResponse(
	ctx context.Context,
	frame protocol.Frame,
	selectorReads *int,
	budgetExhausted *bool,
	cancel context.CancelFunc,
) (*protocol.Frame, error) {
	if len(frame.Data) == 6 && frame.Data[2] == remoteFunctionalModules.group {
		*selectorReads = *selectorReads + 1
		if *selectorReads > functionModuleHighSlotSelectorBudget {
			*budgetExhausted = true
			cancel()
			return nil, ctx.Err()
		}
	}
	return functionModuleSemanticTestResponseForClasses(frame, map[byte]byte{
		semanticStartupSlotFullMaxInstance: circuitManagingDeviceFunctionModuleControllerAddress,
	})
}

// M62 pins that circuit 0's pump status (RR=0x001B) and circuit state
// (RR=0x001E) stay on the circuit surface. Neither is a source for the boiler
// central-heating pump or the boiler heating-status diagnostic.

const boilerCircuitIsolationBoilerAddress = byte(0x08)

func boilerCircuitIsolationUint16(v uint16) *uint16 { return &v }

func boilerCircuitIsolationPoller(boilerAddress byte, circuitRaw, stateRaw uint16) *vaillantSemanticPoller {
	observedAt := time.Unix(100, 0)
	// Fresh controller flow makes the fast refresh publish a boiler status, so
	// the assertions inspect a published payload rather than a missing one.
	flow := 42.5
	return &vaillantSemanticPoller{
		scheduler:          engine.NewSemanticReadScheduler(),
		publication:        newVaillantSemanticPublication(newRetainedSemanticState(), nil, nil),
		source:             0x7F,
		controller:         0x15,
		boilerAddress:      boilerAddress,
		configInterval:     5 * time.Minute,
		boilerFastInterval: 30 * time.Second,
		requestTimeout:     20 * time.Millisecond,
		nowFn:              func() time.Time { return observedAt.Add(10 * time.Second) },
		system: &vaillantSystemSnapshot{
			Controller:                  0x15,
			SystemFlowTemperature:       &flow,
			SystemFlowTemperatureLiveAt: observedAt,
		},
		circuits: map[byte]*vaillantCircuitSnapshot{
			0x00: {
				Instance:           0x00,
				Active:             true,
				Controller:         0x15,
				PumpStatusRaw:      boilerCircuitIsolationUint16(circuitRaw),
				PumpStatusLiveAt:   observedAt,
				CircuitStateRaw:    boilerCircuitIsolationUint16(stateRaw),
				CircuitStateLiveAt: observedAt,
			},
		},
	}
}

// boilerCircuitIsolationAssertSurfaces checks the provider record and the MCP
// JSON of that record. wantBoilerPump nil means the boiler pump must be absent.
func boilerCircuitIsolationAssertSurfaces(t *testing.T, poller *vaillantSemanticPoller, wantBoilerPump *bool, circuitRaw, stateRaw uint16) {
	t.Helper()
	wantCircuitPump := circuitRaw != 0
	wantCircuitState := DecodeCircuitStateToken(&stateRaw)

	status := publishedProvider(poller).BoilerStatus()
	if status == nil {
		t.Fatal("boiler status not published")
	}
	boilerCircuitIsolationAssertOptionalBool(t, "provider boiler pump", status.State.CentralHeatingPumpActive, wantBoilerPump)
	if got := status.Diagnostics.HeatingStatusRaw; got != nil {
		t.Errorf("provider boiler heating status = %d; want unset without a boiler-qualified source", *got)
	}
	circuits := publishedProvider(poller).Circuits()
	if len(circuits) != 1 || circuits[0].Index != 0 {
		t.Fatalf("provider circuits = %+v; want circuit 0", circuits)
	}
	boilerCircuitIsolationAssertOptionalBool(t, "provider circuit pump", circuits[0].State.PumpActive, &wantCircuitPump)
	if circuits[0].State.CircuitState != wantCircuitState {
		t.Errorf("provider circuit state = %q; want %q", circuits[0].State.CircuitState, wantCircuitState)
	}

	mcpBoiler := boilerCircuitIsolationJSONObject(t, output.MCPBoilerStatusFromGraphQL(status))
	boilerCircuitIsolationAssertJSONBool(t, "MCP boiler state.central_heating_pump_active", boilerCircuitIsolationObject(t, mcpBoiler, "state"), "central_heating_pump_active", false, wantBoilerPump)
	if diag, ok := mcpBoiler["diagnostics"].(map[string]any); ok {
		if value, present := diag["heating_status_raw"]; present {
			t.Errorf("MCP boiler diagnostics.heating_status_raw = %v; want absent", value)
		}
	}
	mcpCircuits := output.MCPCircuitsFromGraphQL(circuits)
	if len(mcpCircuits) != 1 {
		t.Fatalf("MCP circuits = %+v; want one circuit", mcpCircuits)
	}
	circuitState := boilerCircuitIsolationObject(t, boilerCircuitIsolationJSONObject(t, mcpCircuits[0]), "state")
	boilerCircuitIsolationAssertJSONBool(t, "MCP circuits[0].state.circuit_pump_unknown_8", circuitState, "circuit_pump_unknown_8", false, &wantCircuitPump)
	if circuitState["circuit_state"] != wantCircuitState {
		t.Errorf("MCP circuits[0].state.circuit_state = %v; want %q", circuitState["circuit_state"], wantCircuitState)
	}
}

func boilerCircuitIsolationAssertOptionalBool(t *testing.T, label string, got, want *bool) {
	t.Helper()
	switch {
	case want == nil && got != nil:
		t.Errorf("%s = %t; want unknown", label, *got)
	case want != nil && got == nil:
		t.Errorf("%s = unknown; want %t", label, *want)
	case want != nil && *got != *want:
		t.Errorf("%s = %t; want %t", label, *got, *want)
	}
}

// m62AssertJSONBool checks a decoded JSON boolean. nullWhenUnknown selects
// GraphQL semantics (key present with null) versus omitempty (key absent).
func boilerCircuitIsolationAssertJSONBool(t *testing.T, label string, object map[string]any, key string, nullWhenUnknown bool, want *bool) {
	t.Helper()
	value, present := object[key]
	if want == nil {
		if nullWhenUnknown && (!present || value != nil) {
			t.Errorf("%s = %v (present=%t); want null", label, value, present)
		}
		if !nullWhenUnknown && present {
			t.Errorf("%s = %v; want absent", label, value)
		}
		return
	}
	if got, ok := value.(bool); !ok || got != *want {
		t.Errorf("%s = %v; want %t", label, value, *want)
	}
}

func boilerCircuitIsolationJSONObject(t *testing.T, value any) map[string]any {
	t.Helper()
	payload, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("marshal %T: %v", value, err)
	}
	var out map[string]any
	if err := json.Unmarshal(payload, &out); err != nil {
		t.Fatalf("unmarshal %T: %v", value, err)
	}
	return out
}

func boilerCircuitIsolationObject(t *testing.T, object map[string]any, key string) map[string]any {
	t.Helper()
	child, ok := object[key].(map[string]any)
	if !ok {
		t.Fatalf("%s = %v; want object", key, object[key])
	}
	return child
}

// The tests in this file use synthetic replies only. The circuit register
// RR001D is addressed by its literal selector so the expectation does not
// depend on the production constant.
const (
	sfmTestCircuitGroup    = byte(0x02)
	sfmTestCircuitInstance = byte(0x00)
	sfmTestRegister        = uint16(0x001D)
	sfmTestFloatRegister   = uint16(0x0010)
	sfmTestCounterName     = "semantic_b524_width_mismatch_total"
	sfmTestLogMarker       = "semantic: b524 reply width mismatch"
)

// sfmReplyBus answers B524 reads from a selector-keyed table. A selector
// without an entry gets a header-only reply, which the reader treats as no
// value (not as a width mismatch).
type sfmReplyBus struct {
	mu      sync.Mutex
	replies map[sfmSelector][]byte
}

type sfmSelector struct {
	opcode   byte
	group    byte
	instance byte
	addr     uint16
}

func (bus *sfmReplyBus) send(_ context.Context, frame protocol.Frame) (*protocol.Frame, error) {
	bus.mu.Lock()
	defer bus.mu.Unlock()
	if len(frame.Data) != 6 {
		return &protocol.Frame{Data: []byte{0x00}}, nil
	}
	key := sfmSelector{
		opcode:   frame.Data[0],
		group:    frame.Data[2],
		instance: frame.Data[3],
		addr:     uint16(frame.Data[4]) | uint16(frame.Data[5])<<8,
	}
	return testB524ResponseForSelectorPayload(frame.Data, bus.replies[key]...), nil
}

// sfmCircuitPoller returns a poller that refreshes one known active circuit
// (instance 0, heating type) and answers the given circuit registers.
func sfmCircuitPoller(now *time.Time, circuitReplies map[uint16][]byte) (*vaillantSemanticPoller, *retainedSemanticState) {
	replies := map[sfmSelector][]byte{
		{opcode: 0x02, group: sfmTestCircuitGroup, instance: sfmTestCircuitInstance, addr: 0x0002}: {0x01, 0x00},
	}
	for addr, payload := range circuitReplies {
		replies[sfmSelector{opcode: 0x02, group: sfmTestCircuitGroup, instance: sfmTestCircuitInstance, addr: addr}] = payload
	}
	bus := &sfmReplyBus{replies: replies}
	provider := newRetainedSemanticState()
	existingType := uint16(1)
	poller := &vaillantSemanticPoller{
		controller:                  0x15,
		source:                      0x7F,
		publication:                 newVaillantSemanticPublication(provider, nil, nil),
		scheduler:                   engine.NewSemanticReadScheduler(),
		requestTimeout:              50 * time.Millisecond,
		circuitFullScanInterval:     semanticCircuitFullScanInterval,
		lastCircuitFullScanAt:       now.Add(-5 * time.Minute),
		lastCircuitFullScanComplete: true,
		nowFn:                       func() time.Time { return *now },
		pauseFn:                     sfmNoPause,
		circuits: map[byte]*vaillantCircuitSnapshot{
			sfmTestCircuitInstance: {Instance: sfmTestCircuitInstance, Active: true, CircuitTypeRaw: &existingType},
		},
		sendFrameFn: bus.send,
	}
	return poller, provider
}

// sfmReaderPoller returns a poller whose bus answers exactly the given
// selector table.
func sfmReaderPoller(now *time.Time, replies map[sfmSelector][]byte) *vaillantSemanticPoller {
	bus := &sfmReplyBus{replies: replies}
	return &vaillantSemanticPoller{
		controller:     0x15,
		source:         0x7F,
		publication:    newVaillantSemanticPublication(newRetainedSemanticState(), nil, nil),
		scheduler:      engine.NewSemanticReadScheduler(),
		requestTimeout: 50 * time.Millisecond,
		nowFn:          func() time.Time { return *now },
		pauseFn:        sfmNoPause,
		sendFrameFn:    bus.send,
	}
}

// sfmNoPause skips the waits between read attempts.
func sfmNoPause(ctx context.Context, _ time.Duration) bool { return ctx.Err() == nil }

// sfmCounter returns the current value of one reader kind in the width
// mismatch counter map, failing the test when the map is not published.
func sfmCounter(t *testing.T, kind string) int64 {
	t.Helper()
	published, ok := expvar.Get(sfmTestCounterName).(*expvar.Map)
	if !ok {
		t.Fatalf("expvar %s is not published as a map", sfmTestCounterName)
	}
	value, ok := published.Get(kind).(*expvar.Int)
	if !ok || value == nil {
		return 0
	}
	return value.Value()
}

// sfmCaptureLog redirects the standard logger for the duration of the test
// and returns a function yielding the width mismatch lines logged so far.
// Callers must not run in parallel.
func sfmCaptureLog(t *testing.T) func() []string {
	t.Helper()
	var mu sync.Mutex
	var buf bytes.Buffer
	previous := log.Writer()
	log.SetOutput(writerFunc(func(p []byte) (int, error) {
		mu.Lock()
		defer mu.Unlock()
		return buf.Write(p)
	}))
	t.Cleanup(func() { log.SetOutput(previous) })
	return func() []string {
		mu.Lock()
		defer mu.Unlock()
		var lines []string
		for _, line := range strings.Split(buf.String(), "\n") {
			if strings.Contains(line, sfmTestLogMarker) {
				lines = append(lines, line[strings.Index(line, sfmTestLogMarker):])
			}
		}
		return lines
	}
}

type writerFunc func([]byte) (int, error)

func (fn writerFunc) Write(p []byte) (int, error) { return fn(p) }

func sfmMismatchLine(opcode, group byte, addr uint16, want, got int) string {
	return fmt.Sprintf("%s opcode=0x%02X group=0x%02X register=0x%04X want=%d got=%d", sfmTestLogMarker, opcode, group, addr, want, got)
}

// sfmPublishedConfig returns the published circuit configuration and the MCP
// JSON of that configuration.
func sfmPublishedConfig(t *testing.T, provider *retainedSemanticState) (output.CircuitConfig, map[string]any) {
	t.Helper()
	circuits := provider.Circuits()
	if len(circuits) != 1 {
		t.Fatalf("published circuits = %d; want 1", len(circuits))
	}
	mcpCircuits := output.MCPCircuitsFromGraphQL(circuits)
	if len(mcpCircuits) != 1 {
		t.Fatalf("MCP circuits = %d; want 1", len(mcpCircuits))
	}
	raw, err := json.Marshal(mcpCircuits[0].Config)
	if err != nil {
		t.Fatal(err)
	}
	mcpConfig := map[string]any{}
	if err := json.Unmarshal(raw, &mcpConfig); err != nil {
		t.Fatal(err)
	}
	return circuits[0].Config, mcpConfig
}

// captureHeatSourcePathLog returns the semantic_startup_reason lines logged
// while fn runs.
func captureHeatSourcePathLog(t *testing.T, fn func()) []string {
	t.Helper()
	var buf bytes.Buffer
	writer, flags, prefix := log.Writer(), log.Flags(), log.Prefix()
	log.SetOutput(&buf)
	log.SetFlags(0)
	log.SetPrefix("")
	defer func() {
		log.SetOutput(writer)
		log.SetFlags(flags)
		log.SetPrefix(prefix)
	}()
	fn()
	var lines []string
	for _, line := range strings.Split(buf.String(), "\n") {
		if strings.HasPrefix(line, "semantic_startup_reason ") {
			lines = append(lines, line)
		}
	}
	return lines
}

func newHeatSourcePathPoller(t *testing.T, devices ...registry.DeviceInfo) *vaillantSemanticPoller {
	t.Helper()
	catalog, err := productids.LoadCatalog()
	if err != nil {
		t.Fatalf("LoadCatalog() error = %v", err)
	}
	return &vaillantSemanticPoller{
		reg:               newTestRegistry(devices...),
		publication:       newVaillantSemanticPublication(newRetainedSemanticState(), nil, nil),
		catalog:           catalog,
		tasks:             newSemanticTaskScheduler(),
		zones:             make(map[byte]*vaillantZoneSnapshot),
		presence:          make(map[byte]*zonePresenceRecord),
		circuits:          make(map[byte]*vaillantCircuitSnapshot),
		zoneMissThreshold: 3,
		zoneHitThreshold:  2,
		// Every frame fails at once and pauses take no time: no device answers.
		sendFrameFn: func(context.Context, protocol.Frame) (*protocol.Frame, error) {
			return nil, ebuserrors.ErrTimeout
		},
		pauseFn: func(ctx context.Context, _ time.Duration) bool { return ctx.Err() == nil },
	}
}

const heatSourcePathSetLine = "semantic_startup_reason reason=no_supported_heat_source_path state=set heat_source=0x%02x"

func rawFlowPoller(payload []byte, requests *int) *vaillantSemanticPoller {
	return &vaillantSemanticPoller{
		controller: 0x15, source: 0x7f,
		sendFrameFn: func(_ context.Context, frame protocol.Frame) (*protocol.Frame, error) {
			(*requests)++
			return testB524ResponseForSelectorPayload(frame.Data, payload...), nil
		},
	}
}

const (
	cfgKeysSystemAllowed = "adaptive_heating_curve_enabled, dhw_flow_temperature_setpoint_maximum, dhw_outdoor_bivalence_temperature_limit, heating_flow_temperature_setpoint_maximum_at_heat_pump_fault, heating_outdoor_bivalence_temperature_limit, installer_access_code, installer_name, installer_phone_number, maintenance_date, room_humidity_limit_maximum"
	cfgKeysBoilerAllowed = "central_heating_pump_mode, dhw_partial_load_power, heating_partial_load_power, installer_phone_number, maximum_dhw_flow_temperature_setpoint, maximum_heating_flow_temperature_setpoint"
	cfgKeysBoilerAddress = byte(0x08)
)

// cfgKeysWrite is one acknowledged write seen by cfgKeysFakeBus.
type cfgKeysWrite struct {
	service  string // "B524" or "B509"
	group    byte
	instance byte
	addr     uint16
	payload  []byte
}

// cfgKeysFakeBus answers B524 OP02 and B509 reads and writes from one register
// store per service. Reads of an unseeded register fail, so the B509 fallback
// selector path can be exercised.
type cfgKeysFakeBus struct {
	mu       sync.Mutex
	requests int
	b524     map[uint16][]byte
	b509     map[uint16][]byte
	writes   []cfgKeysWrite
	ops      []string // "W<service>:<addr>" or "R<service>:<addr>" in order
}
