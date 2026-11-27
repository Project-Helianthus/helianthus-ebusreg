package owner

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/Project-Helianthus/helianthus-ebusgo/engine"
	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
	"github.com/Project-Helianthus/helianthus-ebusgo/transport"
	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/output"
	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
	"github.com/Project-Helianthus/helianthus-ebusreg/vaillant/productids"
	"math"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	semreg "github.com/Project-Helianthus/helianthus-semreg/semreg/v1"
)

type donorGateway struct {
	Registry  *registry.DeviceRegistry
	Bus       *protocol.Bus
	Transport transport.RawTransport
}

func donorDefaultConfig() Config {
	flags := engine.DefaultObserveFirstFeatureFlags()
	return Config{
		TransportConfig:                       engine.TransportConfig{Protocol: engine.TransportENH},
		ScanSource:                            0xF0,
		SemanticRequestTimeout:                2 * time.Second,
		SemanticDiscoveryInterval:             10 * time.Minute,
		SemanticConfigInterval:                5 * time.Minute,
		SemanticStateInterval:                 1 * time.Minute,
		SemanticEnergyInterval:                DefaultSemanticEnergyInterval,
		SemanticReadBreakerFailureBudget:      2,
		SemanticReadBreakerOpenCooldown:       15 * time.Second,
		SemanticReadBreakerHalfOpenProbeLimit: 1,
		SemanticZonePresenceMissThreshold:     DefaultSemanticZonePresenceMissThreshold,
		SemanticZonePresenceHitThreshold:      DefaultSemanticZonePresenceHitThreshold,
		SemanticDHWStaleTTL:                   DefaultSemanticDHWStaleTTL,
		SemanticRegulatorRecheckInterval:      DefaultSemanticRegulatorRecheckInterval,
		SemanticRegulatorAbsenceGrace:         DefaultSemanticRegulatorAbsenceGrace,
		ObserveFirstEnabled:                   flags.ObserveFirstEnabled(),
		PassiveStateDirectApply:               true,
		PassiveConfigDirectApply:              flags.PassiveConfigDirectApply(),
		ExternalWritePolicy:                   flags.ExternalWritePolicy(),
		ObserveFirstFlags:                     flags,
	}
}

func newVaillantSemanticPollerFromDonor(cfg Config, gateway *donorGateway, state *retainedSemanticState, events publicationEvents, cache output.SemanticCachePersister) *vaillantSemanticPoller {
	if gateway == nil {
		gateway = &donorGateway{}
	}
	return newVaillantSemanticPoller(cfg, Substrate{
		Registry:  gateway.Registry,
		Bus:       gateway.Bus,
		Transport: gateway.Transport,
	}, state, events, cache)
}

func publishedFunctionModuleInterpretation(verdict vaillantFunctionModuleInterpretation) output.FunctionModuleInterpretation {
	return output.PublishedFunctionModuleInterpretation(nativeFunctionModuleVerdict(verdict))
}

func publishedFunctionModuleMode(mode vaillantFunctionModuleMode) (output.FunctionModuleSemanticMode, bool) {
	published := output.PublishedFunctionModuleInterpretation(nativeFunctionModuleVerdict(vaillantFunctionModuleInterpretation{Mode: mode}))
	if mode != "" && published.Mode == "" {
		return "", false
	}
	return published.Mode, true
}

func controllerlessCatalogPart(t *testing.T, role string) string {
	t.Helper()
	return controllerlessCatalogPartWith(t, role, productids.ControllerNone)
}

func controllerlessCatalogPartWith(t *testing.T, role string, capability productids.ControllerCapability) string {
	t.Helper()
	catalog, err := productids.LoadCatalog()
	if err != nil {
		t.Fatalf("LoadCatalog() error = %v", err)
	}
	for _, record := range catalog.All {
		part := strings.TrimSpace(record.PartNumber)
		if !strings.EqualFold(strings.TrimSpace(record.Role), role) || len(part) != 10 || !isAllDigits(part) {
			continue
		}
		if catalog.ControllerCapability(part) == capability {
			return part
		}
	}
	t.Fatalf("product catalog has no %q part number with capability %v", role, capability)
	return ""
}

func controllerlessSerial(part string) string {
	return "00-00-00-" + part + "-0000-000000-N0"
}

type semanticSnapshotCaptureSpy struct {
	calls int
	last  output.SemanticCacheSnapshot
}

func (spy *semanticSnapshotCaptureSpy) Save(snapshot output.SemanticCacheSnapshot) error {
	spy.calls++
	spy.last = snapshot
	return nil
}

type watchEfficiencyObserverSpy struct {
	mu                sync.Mutex
	readEvents        []engine.WatchEfficiencyReadEvent
	directApplyEvents []engine.WatchEfficiencyDirectApplyEvent
}

func (spy *watchEfficiencyObserverSpy) ObserveWatchRead(event engine.WatchEfficiencyReadEvent) {
	if spy == nil {
		return
	}
	spy.mu.Lock()
	defer spy.mu.Unlock()
	spy.readEvents = append(spy.readEvents, event)
}

func (spy *watchEfficiencyObserverSpy) ObserveWatchDirectApply(event engine.WatchEfficiencyDirectApplyEvent) {
	if spy == nil {
		return
	}
	spy.mu.Lock()
	defer spy.mu.Unlock()
	spy.directApplyEvents = append(spy.directApplyEvents, event)
}

func (spy *watchEfficiencyObserverSpy) latestReadEvent() (engine.WatchEfficiencyReadEvent, bool) {
	if spy == nil {
		return engine.WatchEfficiencyReadEvent{}, false
	}
	spy.mu.Lock()
	defer spy.mu.Unlock()
	if len(spy.readEvents) == 0 {
		return engine.WatchEfficiencyReadEvent{}, false
	}
	return spy.readEvents[len(spy.readEvents)-1], true
}

func (spy *watchEfficiencyObserverSpy) latestDirectApplyEvent() (engine.WatchEfficiencyDirectApplyEvent, bool) {
	if spy == nil {
		return engine.WatchEfficiencyDirectApplyEvent{}, false
	}
	spy.mu.Lock()
	defer spy.mu.Unlock()
	if len(spy.directApplyEvents) == 0 {
		return engine.WatchEfficiencyDirectApplyEvent{}, false
	}
	return spy.directApplyEvents[len(spy.directApplyEvents)-1], true
}

type b509MutationTestBus struct {
	mu     sync.Mutex
	addr   uint16
	value  []byte
	reads  int
	writes int
	ops    []byte
	onRead func()
}

func newB509MutationTestBus(addr uint16, initialValue []byte) *b509MutationTestBus {
	return &b509MutationTestBus{
		addr:  addr,
		value: append([]byte(nil), initialValue...),
	}
}

func (bus *b509MutationTestBus) Send(_ context.Context, frame protocol.Frame) (*protocol.Frame, error) {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	if len(frame.Data) < 3 {
		return nil, fmt.Errorf("unexpected short b509 frame: %x", frame.Data)
	}

	opcode := frame.Data[0]
	addr := uint16(frame.Data[1])<<8 | uint16(frame.Data[2])
	if addr != bus.addr {
		return nil, fmt.Errorf("unexpected b509 addr 0x%04x (want 0x%04x)", addr, bus.addr)
	}

	switch opcode {
	case vaillantB509OpcodeRead:
		bus.reads++
		bus.ops = append(bus.ops, 'R')
		if bus.onRead != nil {
			bus.onRead()
		}
		payload := append([]byte{vaillantB509OpcodeRead, byte(addr >> 8), byte(addr)}, bus.value...)
		return &protocol.Frame{Data: payload}, nil
	case vaillantB509OpcodeWrite:
		bus.writes++
		bus.ops = append(bus.ops, 'W')
		bus.value = append(bus.value[:0], frame.Data[3:]...)
		return &protocol.Frame{Data: []byte{vaillantB509OpcodeWrite, byte(addr >> 8), byte(addr)}}, nil
	default:
		return nil, fmt.Errorf("unexpected b509 opcode 0x%02x", opcode)
	}
}

func (bus *b509MutationTestBus) Counters() (reads, writes int) {
	bus.mu.Lock()
	defer bus.mu.Unlock()
	return bus.reads, bus.writes
}

func (bus *b509MutationTestBus) Value() []byte {
	bus.mu.Lock()
	defer bus.mu.Unlock()
	return append([]byte(nil), bus.value...)
}

func (bus *b509MutationTestBus) OperationTrace() string {
	bus.mu.Lock()
	defer bus.mu.Unlock()
	return string(append([]byte(nil), bus.ops...))
}

func hasB524Selector(selectors [][]byte, group, instance byte, addr uint16) bool {
	for _, selector := range selectors {
		if len(selector) != 6 {
			continue
		}
		gotAddr := uint16(selector[4]) | uint16(selector[5])<<8
		if selector[2] == group && selector[3] == instance && gotAddr == addr {
			return true
		}
	}
	return false
}

func requireB524Selector(t *testing.T, selectors [][]byte, group, instance byte, addr uint16) {
	t.Helper()
	if !hasB524Selector(selectors, group, instance, addr) {
		t.Fatalf("missing selector group=0x%02x instance=0x%02x addr=0x%04x in % x", group, instance, addr, selectors)
	}
}

func testB524ResponseForSelector(selector []byte) *protocol.Frame {
	if len(selector) != 6 {
		return &protocol.Frame{Data: []byte{0x00}}
	}
	group := selector[2]
	instance := selector[3]
	addr := uint16(selector[4]) | uint16(selector[5])<<8
	payload := testB524PayloadForSelector(group, addr)
	data := []byte{0x01, instance, group, byte(addr), byte(addr >> 8)}
	data = append(data, payload...)
	return &protocol.Frame{Data: data}
}

func testB524ResponseForSelectorPayload(selector []byte, payload ...byte) *protocol.Frame {
	if len(selector) != 6 {
		return &protocol.Frame{Data: []byte{0x00}}
	}
	group := selector[2]
	instance := selector[3]
	addr := uint16(selector[4]) | uint16(selector[5])<<8
	data := []byte{0x01, instance, group, byte(addr), byte(addr >> 8)}
	data = append(data, payload...)
	return &protocol.Frame{Data: data}
}

func testB524PayloadForSelector(group byte, addr uint16) []byte {
	switch {
	case group == localZones.group && (addr == zoneCurrentTemp ||
		addr == zoneTargetTemp ||
		addr == zoneFallbackManualTemp ||
		addr == zoneCurrentHumidity ||
		addr == zoneQuickVetoTemp ||
		addr == zoneQuickVetoDuration ||
		addr == zoneHolidaySetpoint):
		return testB524Float32Payload(21.5)
	case group == localDHW.group && (addr == dhwCurrentTemp || addr == dhwTargetTemp):
		return testB524Float32Payload(48.5)
	case group == localRegulator.group && addr == systemFlowTemperature:
		return testB524Float32Payload(42.5)
	case group == localZones.group && (addr == zoneQuickVetoEndTime ||
		addr == zoneHolidayStartTime ||
		addr == zoneHolidayEndTime):
		return []byte{6, 30}
	case group == localZones.group && (addr == zoneQuickVetoEndDate ||
		addr == zoneHolidayStartDate ||
		addr == zoneHolidayEndDate):
		return []byte{2, 1, 24}
	case group == localDHW.group && (addr == dhwHolidayStartDateRegister || addr == dhwHolidayEndDateRegister):
		return []byte{2, 1, 24}
	case group == localZones.group && (addr == zoneHeatingOpMode ||
		addr == zoneSpecialFunction ||
		addr == zoneValveStatus ||
		addr == zoneRoomTemperatureZoneMappingRawRegister):
		return []byte{1, 0}
	case group == localDHW.group && (addr == dhwOperationMode ||
		addr == dhwSpecialFunction):
		return []byte{1, 0}
	case group == localCircuits.group && (addr == circuitTypeRegister ||
		addr == circuitPumpStatus ||
		addr == circuitCircuitState):
		return []byte{1, 0}
	default:
		return []byte{'T', 'e', 's', 't', 0, 0}
	}
}

func testB524Float32Payload(value float64) []byte {
	payload := make([]byte, 4)
	binary.LittleEndian.PutUint32(payload, math.Float32bits(float32(value)))
	return payload
}

func uint16Ptr(value uint16) *uint16 {
	v := value
	return &v
}

func uint8Ptr(value uint8) *uint8 {
	v := value
	return &v
}

func stringPtr(value string) *string {
	v := value
	return &v
}

func intPtr(value int) *int {
	v := value
	return &v
}

// --- Regulator detection tests (issue #193) ---

func newTestRegistry(devices ...registry.DeviceInfo) *registry.DeviceRegistry {
	reg := registry.NewDeviceRegistry(nil)
	for _, d := range devices {
		reg.Register(d)
	}
	return reg
}

func seedZoneFreshness(snapshot *vaillantZoneSnapshot, source semanticSnapshotSource, stale bool) {
	if snapshot == nil {
		return
	}
	fields := ensureFieldFreshnessMap(&snapshot.FieldFreshness)
	if strings.TrimSpace(snapshot.Name) != "" {
		setFieldFreshness(fields, zoneFieldName, source, stale)
	}
	if strings.TrimSpace(snapshot.OperatingMode) != "" {
		setFieldFreshness(fields, zoneFieldOperatingMode, source, stale)
	}
	if strings.TrimSpace(snapshot.Preset) != "" {
		setFieldFreshness(fields, zoneFieldPreset, source, stale)
	}
	if strings.TrimSpace(snapshot.HvacAction) != "" {
		setFieldFreshness(fields, zoneFieldHvacAction, source, stale)
	}
	if len(snapshot.AllowedModes) > 0 {
		setFieldFreshness(fields, zoneFieldAllowedModes, source, stale)
	}
	if snapshot.CurrentTempC != nil {
		setFieldFreshness(fields, zoneFieldCurrentTempC, source, stale)
	}
	if snapshot.TargetTempC != nil {
		setFieldFreshness(fields, zoneFieldTargetTempC, source, stale)
	}
	if snapshot.HumidityPct != nil {
		setFieldFreshness(fields, zoneFieldCurrentHumidityPct, source, stale)
	}
	if strings.TrimSpace(snapshot.ConfigurationHeatingOperationMode) != "" {
		setFieldFreshness(fields, zoneFieldZoneOperationModeRaw, source, stale)
	}
	if strings.TrimSpace(snapshot.StateSpecialFunction) != "" {
		setFieldFreshness(fields, zoneFieldSpecialFunctionRaw, source, stale)
	}
	if snapshot.ConfigurationRoomTemperatureZoneMappingRaw != nil {
		setFieldFreshness(fields, zoneFieldRoomTemperatureZoneMappingRaw, source, stale)
	}
	if snapshot.ConfigurationAssociatedCircuitRaw != nil {
		setFieldFreshness(fields, zoneFieldZoneCircuitIndexRaw, source, stale)
	}
	if snapshot.ConfigurationCircuitTypeRaw != nil {
		setFieldFreshness(fields, zoneFieldCircuitTypeRaw, source, stale)
	}
	if snapshot.StateValveStatusRaw != nil {
		setFieldFreshness(fields, zoneFieldZoneValveStatusRaw, source, stale)
	}
	if snapshot.QuickVetoTempC != nil {
		setFieldFreshness(fields, zoneFieldQuickVetoTempC, source, stale)
	}
	if snapshot.QuickVetoDurationH != nil {
		setFieldFreshness(fields, zoneFieldQuickVetoDurationH, source, stale)
	}
	if snapshot.QuickVetoEndTime != "" {
		setFieldFreshness(fields, zoneFieldQuickVetoEndTime, source, stale)
	}
	if snapshot.QuickVetoEndDate != "" {
		setFieldFreshness(fields, zoneFieldQuickVetoEndDate, source, stale)
	}
	if snapshot.HolidayStartDate != "" {
		setFieldFreshness(fields, zoneFieldHolidayStartDate, source, stale)
	}
	if snapshot.HolidayEndDate != "" {
		setFieldFreshness(fields, zoneFieldHolidayEndDate, source, stale)
	}
	if snapshot.HolidaySetpointC != nil {
		setFieldFreshness(fields, zoneFieldHolidaySetpointC, source, stale)
	}
	if snapshot.HolidayStartTime != "" {
		setFieldFreshness(fields, zoneFieldHolidayStartTime, source, stale)
	}
	if snapshot.HolidayEndTime != "" {
		setFieldFreshness(fields, zoneFieldHolidayEndTime, source, stale)
	}
}

func seedDhwFreshness(snapshot *vaillantDhwSnapshot, source semanticSnapshotSource, stale bool) {
	if snapshot == nil {
		return
	}
	fields := ensureFieldFreshnessMap(&snapshot.FieldFreshness)
	if strings.TrimSpace(snapshot.OperatingMode) != "" {
		setFieldFreshness(fields, dhwFieldOperatingMode, source, stale)
	}
	if strings.TrimSpace(snapshot.Preset) != "" {
		setFieldFreshness(fields, dhwFieldPreset, source, stale)
	}
	if snapshot.CurrentTempC != nil {
		setFieldFreshness(fields, dhwFieldCurrentTempC, source, stale)
	}
	if snapshot.TargetTempC != nil {
		setFieldFreshness(fields, dhwFieldTargetTempC, source, stale)
	}
	if strings.TrimSpace(snapshot.ConfigurationDHWOperationMode) != "" {
		setFieldFreshness(fields, dhwFieldDhwOperationModeRaw, source, stale)
	}
	if strings.TrimSpace(snapshot.StateSpecialFunction) != "" {
		setFieldFreshness(fields, dhwFieldSpecialFunctionRaw, source, stale)
	}
	if snapshot.HolidayStartDate != "" {
		setFieldFreshness(fields, dhwFieldHolidayStartDate, source, stale)
	}
	if snapshot.HolidayEndDate != "" {
		setFieldFreshness(fields, dhwFieldHolidayEndDate, source, stale)
	}
}

type recordingThermalHandle struct {
	source     semreg.SourceID
	epoch      semreg.SourceEpochID
	generation semreg.Uint64
}

func (handle recordingThermalHandle) Source() semreg.SourceID     { return handle.source }
func (handle recordingThermalHandle) Epoch() semreg.SourceEpochID { return handle.epoch }
func (handle recordingThermalHandle) GenerationID() semreg.Uint64 { return handle.generation }

type recordingThermalHub struct {
	mu       sync.Mutex
	snapshot *semreg.Snapshot
	pending  *semreg.Snapshot
	admitted bool
	sequence int
}

func (hub *recordingThermalHub) Declare(BoilerThermalAssetSpec) error { return nil }

func (hub *recordingThermalHub) Begin(BoilerThermalHandle, semreg.AssetID) (BoilerThermalTransaction, error) {
	return &recordingThermalTransaction{hub: hub}, nil
}

func (hub *recordingThermalHub) ReadAdmitted(semreg.AssetID) bool {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	return hub.admitted
}

func (hub *recordingThermalHub) Snapshot(semreg.AssetID) (semreg.Snapshot, error) {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	if hub.snapshot == nil {
		return semreg.Snapshot{}, nil
	}
	return *hub.snapshot, nil
}

type recordingThermalTransaction struct {
	hub *recordingThermalHub
}

func (tx *recordingThermalTransaction) Current() (semreg.Snapshot, bool) {
	tx.hub.mu.Lock()
	defer tx.hub.mu.Unlock()
	if tx.hub.snapshot == nil {
		return semreg.Snapshot{}, false
	}
	return *tx.hub.snapshot, true
}

func (tx *recordingThermalTransaction) NextSequence() semreg.Uint64 {
	tx.hub.mu.Lock()
	defer tx.hub.mu.Unlock()
	tx.hub.sequence++
	return semreg.Uint64(strconv.Itoa(tx.hub.sequence))
}

func (tx *recordingThermalTransaction) SourceRetirements() []semreg.SourceEpochID {
	return []semreg.SourceEpochID{}
}

func (tx *recordingThermalTransaction) Apply(batch semreg.PublicationBatch, _ semreg.MonotonicPoint) error {
	tx.hub.mu.Lock()
	defer tx.hub.mu.Unlock()
	snapshot := semreg.Snapshot{Contract: semreg.ContractKernelV1, AssetID: batch.AssetID}
	if tx.hub.snapshot != nil {
		snapshot = *tx.hub.snapshot
	}
	snapshot.Sources = append(snapshot.Sources, batch.SourceUpserts...)
	snapshot.Bindings = append(snapshot.Bindings, batch.BindingUpserts...)
	snapshot.IdentityLinks = append(snapshot.IdentityLinks, batch.IdentityLinkUpserts...)
	if len(batch.FactUpserts) > 0 {
		snapshot.Facts = []semreg.FactEnvelope{{
			AssetID:    batch.AssetID,
			Key:        batch.FactUpserts[0].Key,
			Candidates: append([]semreg.FactCandidate(nil), batch.FactUpserts...),
			Conflicts:  []semreg.Conflict{},
			Revision:   "1",
		}}
	}
	if len(batch.FactWithdrawals) > 0 {
		snapshot.Facts = nil
	}
	tx.hub.pending = &snapshot
	return nil
}

func (tx *recordingThermalTransaction) Commit() error {
	tx.hub.mu.Lock()
	defer tx.hub.mu.Unlock()
	if tx.hub.pending != nil {
		tx.hub.snapshot = tx.hub.pending
		tx.hub.admitted = true
		tx.hub.pending = nil
	}
	return nil
}

func (tx *recordingThermalTransaction) Abort() {
	tx.hub.mu.Lock()
	tx.hub.pending = nil
	tx.hub.mu.Unlock()
}

func qualifiedBoilerThermalFixture(t *testing.T) (*boilerThermalPublisher, *int64) {
	t.Helper()
	tick := int64(0)
	source := semreg.SourceID("source:synthetic-boiler")
	epoch := semreg.SourceEpochID("epoch:synthetic-boiler")
	generation := semreg.Uint64("1")
	handle := recordingThermalHandle{source: source, epoch: epoch, generation: generation}
	identity := boilerThermalIdentity{
		Asset: "asset:synthetic-boiler", Source: source, Epoch: epoch, Generation: generation,
		Binding: "binding:synthetic-boiler", Service: "service:synthetic-boiler",
		Product: "synthetic-boiler-family", Profile: "boiler.synthetic", Revision: "1",
		Codec: "DATA2c", ResponseForm: "operation_selector",
		Initiator: 0x31, Target: 0x36, PB: 0xb5, SB: 9, Operation: 0x0d, Selector: [2]byte{0x18, 0},
	}
	publisher := &boilerThermalPublisher{
		hub:    &recordingThermalHub{},
		handle: func() (BoilerThermalHandle, error) { return handle, nil },
		binding: &boilerThermalBinding{
			Identity:  identity,
			Qualified: true,
			Policy: semreg.FreshnessPolicy{
				PolicyID: "policy:synthetic-boiler", Version: "1.0.0",
				FreshForNS: "100", RetainForNS: "200", MaxWallUncertaintyNS: "0",
			},
		},
	}
	publisher.now = func() (semreg.TimePoint, semreg.MonotonicPoint) {
		return semreg.TimePoint{
				UnixNanoseconds: semreg.Int64(strconv.FormatInt(1800000000000000000+tick, 10)),
				ClockID:         "clock.utc",
				UncertaintyNS:   "0",
			}, semreg.MonotonicPoint{
				ClockEpochID: "clock:synthetic-boiler",
				Nanoseconds:  semreg.Uint64(strconv.FormatInt(tick, 10)),
			}
	}
	wall, mono := publisher.now()
	observation := &boilerThermalObservation{
		Identity: identity,
		Request:  []byte{0x0d, 0x18, 0},
		Reply:    []byte{0x0d, 0x18, 0, 0x20, 3},
		Received: wall,
		Receipt:  mono,
		Valid:    true,
	}
	publisher.binding.Basis = []semreg.EvidenceRef{boilerThermalEvidence(observation)}
	publisher.attempt = observation
	return publisher, &tick
}
