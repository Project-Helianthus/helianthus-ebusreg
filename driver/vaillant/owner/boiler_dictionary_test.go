package owner

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/Project-Helianthus/helianthus-ebusgo/engine"
	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"

	catalogvaillant "github.com/Project-Helianthus/helianthus-ebusreg/catalog/vaillant"
	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/output"
)

func TestVaillantB509DictionaryKeepsExecutableDecode(t *testing.T) {
	if len(b524EnergyQueries) != 12 {
		t.Fatalf("len(b524EnergyQueries) = %d; want 12", len(b524EnergyQueries))
	}
	bus := &boilerReplyBus{values: map[uint16][]byte{
		boilerB509FlameActive:    {0x00},
		boilerB509GasValveActive: {0x00},
		boilerB509StateNumber:    {0x07},
		boilerB509PumpHours:      {0x05, 0x00},
		boilerB509FanHours:       {0x03, 0x00},
	}}
	poller, provider := newCataloguePoller(bus.send)
	poller.boilerAddress = 0x08
	poller.boilerFastInterval = 30 * time.Second
	poller.boilerMediumInterval = 5 * time.Minute
	poller.boilerSlowInterval = 10 * time.Minute
	schedules := poller.boilerStatusTierSchedules()
	if len(schedules) != 3 || schedules[0].interval != 30*time.Second || schedules[1].interval != 5*time.Minute || schedules[2].interval != 10*time.Minute {
		t.Fatalf("boiler schedules = %+v; want the existing fast, medium and slow intervals", schedules)
	}

	bus.reads = nil
	poller.refreshBoilerStatusTier(context.Background(), boilerStatusTierFast)
	if containsAddress(bus.reads, boilerB509PumpHours) {
		t.Fatalf("fast tier read pump hours %x; the slow tier owns that register", bus.reads)
	}
	if !containsAddress(bus.reads, boilerB509FlameActive) {
		t.Fatalf("fast tier reads = %x; want flame", bus.reads)
	}
	status := provider.BoilerStatus()
	if status == nil || status.State.FlameActive == nil || !*status.State.FlameActive {
		t.Fatalf("flame after 0x00 = %+v; want true", status)
	}
	if status.Diagnostics.PumpHours != nil {
		t.Fatalf("pump hours after fast tier = %v; want unset", *status.Diagnostics.PumpHours)
	}

	bus.reads = nil
	poller.refreshBoilerStatusTier(context.Background(), boilerStatusTierSlow)
	if !containsAddress(bus.reads, boilerB509PumpHours) || !containsAddress(bus.reads, boilerB509FanHours) {
		t.Fatalf("slow tier reads = %x; want pump hours and fan hours", bus.reads)
	}
	status = provider.BoilerStatus()
	if status == nil || status.Diagnostics.PumpHours == nil || *status.Diagnostics.PumpHours != 10 {
		t.Fatalf("PumpHours = %+v; want 10 from bytes 05 00", status.Diagnostics)
	}
	if status.Diagnostics.FanHours == nil || *status.Diagnostics.FanHours != 6 {
		t.Fatalf("FanHours = %+v; want 6 from the existing hoursum decode", status.Diagnostics)
	}
	if status.State.StateNumber == nil || *status.State.StateNumber != 7 {
		t.Fatalf("StateNumber = %+v; want 7", status.State.StateNumber)
	}

	selectors, err := fillBoilerDictionary(poller, context.Background(), poller.boilerAddress)
	if err != nil {
		t.Fatalf("fill boiler dictionary: %v", err)
	}
	flame := boilerSelector(selectors, boilerB509FlameActive)
	if !bytes.Equal(flame.Raw, []byte{0x00}) || flame.Label != "flame_active" || flame.Name != "flame_active" || flame.CodecUsed != "onoff" || flame.CodecUsed == boilerDictionaryCodecOnOff {
		t.Fatalf("flame evidence = %+v; want raw 00, separate annotation flame_active and executable codec onoff", flame)
	}
	gas := boilerSelector(selectors, boilerB509GasValveActive)
	if !bytes.Equal(gas.Raw, []byte{0x00}) || gas.Label != "gas_valve_active" || gas.CodecUsed != "onoff" {
		t.Fatalf("gas valve evidence = %+v; want raw 00, annotation gas_valve_active and codec onoff", gas)
	}
	pump := boilerSelector(selectors, boilerB509PumpHours)
	if !bytes.Equal(pump.Raw, []byte{0x05, 0x00}) || pump.CodecUsed != "hoursum2" || pump.Name != "pump_operating_time" || pump.Label != "pump_operating_time" {
		t.Fatalf("pump evidence = %+v; want raw 05 00, hoursum2 and the annotation beside the executable decode", pump)
	}
	state := boilerSelector(selectors, boilerB509StateNumber)
	if !bytes.Equal(state.Raw, []byte{0x07}) || state.Label != catalogueUnresolved || state.Name != "boiler_state_code" || state.Note != "preserve-raw-bytes" {
		t.Fatalf("state evidence = %+v; want unresolved annotation beside raw 07", state)
	}
	declaration, found := dictionaryDeclaration(boilerB509FlameActive)
	if !found || declaration != boilerDictionaryCodecOnOff || dictionaryEstablishesExecutable(boilerB509FlameActive, declaration) {
		t.Fatalf("flame declaration %q found=%v establishes executable; want the catalog declaration withheld", declaration, found)
	}
	if _, foundPump := dictionaryDeclaration(boilerB509PumpHours); foundPump {
		t.Fatalf("pump hours has a dictionary row; the executable hoursum decode must stay without one")
	}
	polled := polledBoilerSelectors()
	entries, err := catalogvaillant.B509BoilerCodecDictionaryEntries()
	if err != nil {
		t.Fatalf("dictionary entries: %v", err)
	}
	outside := 0
	for _, entry := range entries {
		addr, ok := dictionarySelectorAddress(entry.B509BoilerCodecDictionaryKey.B509BoilerCodecDictionarySelector)
		if !ok {
			continue
		}
		if _, already := polled[addr]; already {
			continue
		}
		outside++
		item := boilerSelector(selectors, addr)
		if item.Selector != addr || item.Tier != "slow" || item.CodecUsed != "" || item.Note != "dictionary codec withheld" || item.Label != catalogueUnresolved {
			t.Fatalf("unread dictionary selector 0x%04x = %+v; want raw slow evidence with the codec withheld", addr, item)
		}
	}
	if outside == 0 {
		for _, item := range selectors {
			if _, already := polled[item.Selector]; !already {
				t.Fatalf("fill published selector 0x%04x; the dictionary has no selector outside the polled set", item.Selector)
			}
		}
	}

	bus.values[boilerB509FlameActive] = []byte{0xF0}
	poller.scheduler = engine.NewSemanticReadScheduler()
	poller.refreshBoilerStatusTier(context.Background(), boilerStatusTierFast)
	status = provider.BoilerStatus()
	if status.State.FlameActive == nil || *status.State.FlameActive {
		t.Fatalf("flame after F0 = %+v; want false", status.State.FlameActive)
	}
	if status.Diagnostics.PumpHours == nil || *status.Diagnostics.PumpHours != 10 {
		t.Fatalf("PumpHours after flame flip = %+v; want 10", status.Diagnostics.PumpHours)
	}
	poller.scheduler = engine.NewSemanticReadScheduler()
	selectors, err = fillBoilerDictionary(poller, context.Background(), poller.boilerAddress)
	if err != nil {
		t.Fatalf("refill boiler dictionary: %v", err)
	}
	flame = boilerSelector(selectors, boilerB509FlameActive)
	if !bytes.Equal(flame.Raw, []byte{0xF0}) || flame.CodecUsed != "onoff" || flame.Label != "flame_active" {
		t.Fatalf("flame evidence after F0 = %+v; want raw F0, annotation flame_active and codec onoff", flame)
	}
	published := provider.BoilerStatus()
	publishedFlame := boilerSelector(published.Selectors(), boilerB509FlameActive)
	if !bytes.Equal(publishedFlame.Raw, []byte{0xF0}) || published.State.FlameActive == nil || *published.State.FlameActive {
		t.Fatalf("published flame = %+v evidence %+v; want false beside raw F0", published.State.FlameActive, publishedFlame)
	}
	mapped := output.MCPBoilerStatusFromGraphQL(published)
	if mapped.Selectors() == nil || !bytes.Equal(boilerSelector(mapped.Selectors(), boilerB509FlameActive).Raw, []byte{0xF0}) {
		t.Fatalf("mcp flame selectors = %+v; want raw F0", mapped.Selectors())
	}
}

func TestVaillantB509DictionaryRejectsMutatedIdentityAxis(t *testing.T) {
	base := boilerDictionaryKey(boilerB509FlameActive)
	if _, found, err := catalogvaillant.B509BoilerCodecDictionaryLookup(base); err != nil || !found {
		t.Fatalf("unmutated dictionary lookup found=%v err=%v", found, err)
	}
	selector := append([]uint16(nil), base.B509BoilerCodecDictionarySelector...)
	mutations := []struct {
		field  string
		mutate func(*catalogvaillant.B509BoilerCodecDictionaryKey)
	}{
		{"family", func(key *catalogvaillant.B509BoilerCodecDictionaryKey) {
			key.B509BoilerCodecDictionaryFamily = "absent"
		}},
		{"variant", func(key *catalogvaillant.B509BoilerCodecDictionaryKey) {
			key.B509BoilerCodecDictionaryVariant = "absent"
		}},
		{"primary", func(key *catalogvaillant.B509BoilerCodecDictionaryKey) {
			value := byte(1)
			key.B509BoilerCodecDictionaryPrimary = &value
		}},
		{"secondary", func(key *catalogvaillant.B509BoilerCodecDictionaryKey) {
			value := byte(1)
			key.B509BoilerCodecDictionarySecondary = &value
		}},
		{"operation", func(key *catalogvaillant.B509BoilerCodecDictionaryKey) {
			value := byte(1)
			key.B509BoilerCodecDictionaryOperation = &value
		}},
		{"direction", func(key *catalogvaillant.B509BoilerCodecDictionaryKey) {
			key.B509BoilerCodecDictionaryDirection = "absent"
		}},
		{"selector", func(key *catalogvaillant.B509BoilerCodecDictionaryKey) {
			key.B509BoilerCodecDictionarySelector = []uint16{selector[0] + 1, selector[1]}
		}},
	}
	for _, mutation := range mutations {
		key := cloneDictionaryKey(base)
		mutation.mutate(&key)
		_, found, err := catalogvaillant.B509BoilerCodecDictionaryLookup(key)
		if err != nil || found {
			t.Fatalf("%s mutation still matched; dropping %s from equality would keep the original row", mutation.field, mutation.field)
		}
	}

	annotation := boilerAnnotationKey(boilerB509FlameActive)
	if _, found, err := catalogvaillant.B509SymbolicAnnotationLookup(annotation); err != nil || !found {
		t.Fatalf("unmutated annotation lookup found=%v err=%v", found, err)
	}
	annotationMutations := []struct {
		field  string
		mutate func(*catalogvaillant.B509SymbolicAnnotationKey)
	}{
		{"family", func(key *catalogvaillant.B509SymbolicAnnotationKey) { key.B509SymbolicAnnotationFamily = "absent" }},
		{"variant", func(key *catalogvaillant.B509SymbolicAnnotationKey) { key.B509SymbolicAnnotationVariant = "absent" }},
		{"primary", func(key *catalogvaillant.B509SymbolicAnnotationKey) { key.B509SymbolicAnnotationPrimary = 1 }},
		{"secondary", func(key *catalogvaillant.B509SymbolicAnnotationKey) {
			value := uint8(1)
			key.B509SymbolicAnnotationSecondary = &value
		}},
		{"operation", func(key *catalogvaillant.B509SymbolicAnnotationKey) {
			value := uint8(1)
			key.B509SymbolicAnnotationOperation = &value
		}},
		{"direction", func(key *catalogvaillant.B509SymbolicAnnotationKey) { key.B509SymbolicAnnotationDirection = "absent" }},
		{"field", func(key *catalogvaillant.B509SymbolicAnnotationKey) { key.B509SymbolicAnnotationField = "absent" }},
		{"alternative", func(key *catalogvaillant.B509SymbolicAnnotationKey) {
			key.B509SymbolicAnnotationAlternative = "absent"
		}},
		{"group", func(key *catalogvaillant.B509SymbolicAnnotationKey) {
			value := uint8(1)
			key.B509SymbolicAnnotationGroup = &value
		}},
		{"instance_scope", func(key *catalogvaillant.B509SymbolicAnnotationKey) {
			key.B509SymbolicAnnotationInstanceScope = "absent"
		}},
		{"register", func(key *catalogvaillant.B509SymbolicAnnotationKey) {
			value := uint16(1)
			key.B509SymbolicAnnotationRegister = &value
		}},
		{"class", func(key *catalogvaillant.B509SymbolicAnnotationKey) { key.B509SymbolicAnnotationClass = "absent" }},
		{"selector", func(key *catalogvaillant.B509SymbolicAnnotationKey) {
			value := "0x0001"
			key.B509SymbolicAnnotationSelector = &value
		}},
	}
	for _, mutation := range annotationMutations {
		key := cloneAnnotationKey(annotation)
		mutation.mutate(&key)
		_, found, err := catalogvaillant.B509SymbolicAnnotationLookup(key)
		if err != nil || found {
			t.Fatalf("%s mutation still matched; dropping %s from equality would keep the original annotation", mutation.field, mutation.field)
		}
	}
}

func TestVaillantB509DictionaryOpenSelectorPermissionBit(t *testing.T) {
	bus := &openSelectorBus{flags: catalogueWritePermission, payloads: cataloguePayloads()}
	poller, provider := newCataloguePoller(bus.send)
	provider.SetBoilerStatus(&output.BoilerStatus{})
	raw := []byte{4, 0}
	if labelFamilyEnumeration(familyCodebookKey("controller-700", "u16", "dynamic", localZones.group, zoneSpecialFunction), raw).Label != catalogueUnresolved {
		t.Fatalf("open selector documentation changed while the permission bit was set")
	}
	operation, admitted := admitOpenSelector(poller, context.Background())
	if !admitted || operation.Namespace != "ebusreg" || operation.Name != rawZoneSpecialFunctionName || operation.Register != zoneSpecialFunction {
		t.Fatalf("admitted operation = %+v admitted=%v; want the raw zone operation", operation, admitted)
	}
	if err := callRawSelectorOperation(poller, context.Background(), rawZoneSpecialFunctionName, []byte{0x01}); err != nil {
		t.Fatalf("admitted call: %v", err)
	}
	if bus.writes != 1 || !bytes.Equal(bus.lastWrite, []byte{0x01}) {
		t.Fatalf("writes = %d payload %x; want one write of 01", bus.writes, bus.lastWrite)
	}
	status := provider.BoilerStatus()
	if len(status.Operations()) != 1 || status.Operations()[0].Name != rawZoneSpecialFunctionName {
		t.Fatalf("published operations = %+v; want the admitted raw operation", status.Operations())
	}

	bus.flags = 0x00
	if _, admitted = admitOpenSelector(poller, context.Background()); admitted {
		t.Fatalf("cleared permission bit still admitted an operation")
	}
	writes := bus.writes
	if err := callRawSelectorOperation(poller, context.Background(), rawZoneSpecialFunctionName, []byte{0x01}); err == nil || bus.writes != writes {
		t.Fatalf("cleared call err=%v writes=%d; want no send", err, bus.writes)
	}
	if labelFamilyEnumeration(familyCodebookKey("controller-720", "u16", "dynamic", localZones.group, zoneSpecialFunction), raw).Label != catalogueUnresolved {
		t.Fatalf("open selector documentation changed while the permission bit was clear")
	}
	if err := poller.writeB524Value(context.Background(), vaillantB524OpcodeLocal, localZones.group, 0, zoneSpecialFunction, []byte{0x01}); err != nil {
		t.Fatalf("existing write: %v", err)
	}
	if bus.writes != writes+1 {
		t.Fatalf("existing write count = %d; want %d", bus.writes, writes+1)
	}
}

type boilerReplyBus struct {
	values map[uint16][]byte
	reads  []uint16
}

func (bus *boilerReplyBus) send(_ context.Context, frame protocol.Frame) (*protocol.Frame, error) {
	if frame.Primary != vaillantB509Primary || len(frame.Data) < 3 || frame.Data[0] != vaillantB509OpcodeRead {
		return &protocol.Frame{Data: []byte{0x00}}, nil
	}
	addr := uint16(frame.Data[1])<<8 | uint16(frame.Data[2])
	bus.reads = append(bus.reads, addr)
	value := bus.values[addr]
	if value == nil {
		value = []byte{0x00}
	}
	data := append([]byte{vaillantB509OpcodeRead, frame.Data[1], frame.Data[2]}, value...)
	return &protocol.Frame{Data: data}, nil
}

type openSelectorBus struct {
	flags     byte
	payloads  map[[2]uint16][]byte
	writes    int
	lastWrite []byte
}

func (bus *openSelectorBus) send(_ context.Context, frame protocol.Frame) (*protocol.Frame, error) {
	if len(frame.Data) >= 2 && frame.Data[1] == vaillantB524OpWrite {
		bus.writes++
		if len(frame.Data) > 6 {
			bus.lastWrite = append([]byte(nil), frame.Data[6:]...)
		}
		return &protocol.Frame{Data: []byte{0x00}}, nil
	}
	if len(frame.Data) != 6 {
		return &protocol.Frame{Data: []byte{0x00}}, nil
	}
	group := frame.Data[2]
	addr := uint16(frame.Data[4]) | uint16(frame.Data[5])<<8
	payload := bus.payloads[[2]uint16{uint16(group), addr}]
	if payload == nil {
		payload = []byte{0}
	}
	if group == localZones.group && addr == zoneSpecialFunction {
		data := []byte{bus.flags, frame.Data[3], group, byte(addr), byte(addr >> 8)}
		data = append(data, payload...)
		return &protocol.Frame{Data: data}, nil
	}
	return testB524ResponseForSelectorPayload(frame.Data, payload...), nil
}

func boilerSelector(items []output.BoilerSelectorEvidence, addr uint16) output.BoilerSelectorEvidence {
	for _, item := range items {
		if item.Selector == addr {
			return item
		}
	}
	return output.BoilerSelectorEvidence{}
}

func containsAddress(items []uint16, addr uint16) bool {
	for _, item := range items {
		if item == addr {
			return true
		}
	}
	return false
}

func boilerAnnotationKey(addr uint16) catalogvaillant.B509SymbolicAnnotationKey {
	selector := "0x0500"
	if addr != boilerB509FlameActive {
		selector = "0x0000"
	}
	secondary := uint8(9)
	operation := uint8(13)
	return catalogvaillant.B509SymbolicAnnotationKey{
		B509SymbolicAnnotationFamily:        "Vaillant",
		B509SymbolicAnnotationVariant:       "none",
		B509SymbolicAnnotationPrimary:       181,
		B509SymbolicAnnotationSecondary:     &secondary,
		B509SymbolicAnnotationOperation:     &operation,
		B509SymbolicAnnotationDirection:     "none",
		B509SymbolicAnnotationField:         "symbolic_annotation",
		B509SymbolicAnnotationAlternative:   "none",
		B509SymbolicAnnotationInstanceScope: "none",
		B509SymbolicAnnotationClass:         "boiler",
		B509SymbolicAnnotationSelector:      &selector,
	}
}

func cloneDictionaryKey(key catalogvaillant.B509BoilerCodecDictionaryKey) catalogvaillant.B509BoilerCodecDictionaryKey {
	if key.B509BoilerCodecDictionaryPrimary != nil {
		value := *key.B509BoilerCodecDictionaryPrimary
		key.B509BoilerCodecDictionaryPrimary = &value
	}
	if key.B509BoilerCodecDictionarySecondary != nil {
		value := *key.B509BoilerCodecDictionarySecondary
		key.B509BoilerCodecDictionarySecondary = &value
	}
	if key.B509BoilerCodecDictionaryOperation != nil {
		value := *key.B509BoilerCodecDictionaryOperation
		key.B509BoilerCodecDictionaryOperation = &value
	}
	key.B509BoilerCodecDictionarySelector = append([]uint16(nil), key.B509BoilerCodecDictionarySelector...)
	return key
}

func cloneAnnotationKey(key catalogvaillant.B509SymbolicAnnotationKey) catalogvaillant.B509SymbolicAnnotationKey {
	if key.B509SymbolicAnnotationSecondary != nil {
		value := *key.B509SymbolicAnnotationSecondary
		key.B509SymbolicAnnotationSecondary = &value
	}
	if key.B509SymbolicAnnotationOperation != nil {
		value := *key.B509SymbolicAnnotationOperation
		key.B509SymbolicAnnotationOperation = &value
	}
	if key.B509SymbolicAnnotationGroup != nil {
		value := *key.B509SymbolicAnnotationGroup
		key.B509SymbolicAnnotationGroup = &value
	}
	if key.B509SymbolicAnnotationRegister != nil {
		value := *key.B509SymbolicAnnotationRegister
		key.B509SymbolicAnnotationRegister = &value
	}
	if key.B509SymbolicAnnotationSelector != nil {
		value := *key.B509SymbolicAnnotationSelector
		key.B509SymbolicAnnotationSelector = &value
	}
	return key
}
