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

func TestVaillantB524CodebooksLabelEnumerationBesideRaw(t *testing.T) {
	bus := &catalogueReplyBus{payloads: cataloguePayloads()}
	poller, provider := newCataloguePoller(bus.send)
	provider.SetZones([]output.Zone{{ID: "zone-1"}})
	provider.SetDHW(&output.DhwStatus{})
	provider.SetCircuits([]output.CircuitStatus{{Index: 1}})
	provider.SetSystem(&output.SystemStatus{})

	view := fillFamilyCodebooks(poller, context.Background(), "controller-700")
	if bus.writes != 0 {
		t.Fatalf("writes = %d; codebook fill sends a write when no operation is called", bus.writes)
	}
	assertEnumeration(t, view.Enumerations, 1, 3, "u16", []byte{1, 0}, 1, "automatic", "Hypothesis")
	assertEnumeration(t, view.Enumerations, 1, 13, "u8", []byte{6}, 6, "start a one-time DHW load", "Hypothesis")
	assertEnumeration(t, view.Enumerations, 2, 2, "u16", []byte{5, 0}, 5, "pool", "Hypothesis")
	assertEnumeration(t, view.Enumerations, 2, 29, "u8", []byte{2}, 2, "party", "Hypothesis")
	assertEnumeration(t, view.Enumerations, 3, 6, "u16", []byte{3, 0}, 3, "setback", "Hypothesis")
	assertEnumeration(t, view.Enumerations, 3, 27, "u8", []byte{1, 0}, 1, "ventilation-boost", "Hypothesis")
	assertEnumeration(t, view.Enumerations, 3, 27, "u16", []byte{1, 0}, 1, "ventilation-boost", "Hypothesis")
	assertEnumeration(t, view.Enumerations, 3, 0x000E, "u16", []byte{4, 0}, 4, catalogueUnresolved, "Unknown")

	zones := provider.Zones()
	if len(zones) != 1 {
		t.Fatalf("zones = %d; want the pre-published zone", len(zones))
	}
	assertEnumeration(t, zones[0].State.Enumerations(), 3, 27, "u8", []byte{1, 0}, 1, "ventilation-boost", "Hypothesis")
	assertEnumeration(t, zones[0].State.Enumerations(), 3, 0x000E, "u16", []byte{4, 0}, 4, catalogueUnresolved, "Unknown")
	if enumerationAt(zones[0].State.Enumerations(), 1, 13, "u8").Label != "" {
		t.Fatalf("zone view carries the hot-water special-function label; group is not part of the published zone enumeration")
	}
	dhw := provider.DHW()
	assertEnumeration(t, dhw.State.Enumerations(), 1, 3, "u16", []byte{1, 0}, 1, "automatic", "Hypothesis")
	assertEnumeration(t, dhw.State.Enumerations(), 1, 13, "u8", []byte{6}, 6, "start a one-time DHW load", "Hypothesis")
	circuits := provider.Circuits()
	assertEnumeration(t, circuits[0].Config.Enumerations(), 2, 29, "u8", []byte{2}, 2, "party", "Hypothesis")
	assertEnumeration(t, circuits[0].Config.Enumerations(), 2, 2, "u16", []byte{5, 0}, 5, "pool", "Hypothesis")

	mapped := output.MCPZonesFromGraphQL(zones)
	assertEnumeration(t, mapped[0].State.Enumerations(), 3, 27, "u8", []byte{1, 0}, 1, "ventilation-boost", "Hypothesis")
	mappedDHW := output.MCPDHWFromGraphQL(dhw)
	assertEnumeration(t, mappedDHW.State.Enumerations(), 1, 13, "u8", []byte{6}, 6, "start a one-time DHW load", "Hypothesis")
	mappedCircuits := output.MCPCircuitsFromGraphQL(circuits)
	assertEnumeration(t, mappedCircuits[0].Config.Enumerations(), 2, 29, "u8", []byte{2}, 2, "party", "Hypothesis")

	day := scalarAt(view.Scalars, 0, 0x0002)
	if day.Value == nil || *day.Value != 21.5 || day.Unit != "°C" || day.Name != "" {
		t.Fatalf("system scalar 0x0002 = %+v; want raw float 21.5 °C and no invented name", day)
	}
	minutes := scalarAt(view.Scalars, 0, 0x0017)
	if minutes.Value == nil || *minutes.Value != 300 || minutes.Unit != "min" {
		t.Fatalf("system scalar 0x0017 = %+v; want uint16 300 min", minutes)
	}
	rawWord := scalarAt(view.Scalars, 0, 0x0004)
	if rawWord.Value != nil || !bytes.Equal(rawWord.Raw, []byte{0x2C, 0x01}) {
		t.Fatalf("controller-700 scalar 0x0004 = %+v; the family leaves that word raw", rawWord)
	}
	unnamed := scalarAt(view.Scalars, 1, 0x0015)
	if unnamed.Name != "" || unnamed.Value != nil || !bytes.Equal(unnamed.Raw, catalogueFloatRaw()) {
		t.Fatalf("unnamed scalar = %+v; want raw bytes and no name, unit or decoded value", unnamed)
	}
	fastSeen, slowSeen := false, false
	for _, item := range view.Scalars {
		_, profile, _ := semanticReadB524LocalDescriptorPolicy(item.Group, item.Register)
		if item.Tier != string(profile) {
			t.Fatalf("group %d register 0x%04x tier %q; want the existing policy %q", item.Group, item.Register, item.Tier, profile)
		}
		if item.Tier == string(engine.WatchFreshnessProfileStateFast) {
			fastSeen = true
		}
		if item.Tier == string(engine.WatchFreshnessProfileStateSlow) {
			slowSeen = true
		}
	}
	if !fastSeen {
		t.Fatalf("no scalar kept the fast tier; a constant tier would hide the policy")
	}
	if slowSeen {
		t.Fatalf("an unread scalar used the slow tier; the existing policy for these registers is not slow")
	}
	if got, want := len(view.Scalars), catalogueScalarCount("controller-700"); got != want {
		t.Fatalf("scalar count = %d; want %d", got, want)
	}
	if view.Ventilation == nil || view.Ventilation.DayFanSpeedLimit == nil || *view.Ventilation.DayFanSpeedLimit != 300 {
		t.Fatalf("controller-700 ventilation = %+v; want day word 300", view.Ventilation)
	}
	if view.Ventilation.NightFanSpeedLimit == nil || *view.Ventilation.NightFanSpeedLimit != 400 {
		t.Fatalf("controller-700 night ventilation = %+v; want 400", view.Ventilation)
	}
	named := scalarAt(view.Scalars, 0, 0x0042)
	if named.Name != "ventilation_fan_speed_limit_maximum_at_day" || named.Value == nil || *named.Value != 300 || named.Unit != "" {
		t.Fatalf("named ventilation scalar = %+v; want the native name, word 300 and no invented unit", named)
	}
	system := provider.System()
	mappedSystem := output.MCPSystemFromGraphQL(system)
	if system.Ventilation() == nil || mappedSystem.Ventilation() == nil || mappedSystem.Catalogue() == nil {
		t.Fatalf("system ventilation missing from provider=%+v mcp=%+v", system.Ventilation(), mappedSystem.Ventilation())
	}
	if mappedSystem.Ventilation().DayFanSpeedLimit == nil || *mappedSystem.Ventilation().DayFanSpeedLimit != 300 {
		t.Fatalf("mcp ventilation = %+v; want day word 300", mappedSystem.Ventilation())
	}

	other, _ := newCataloguePoller(bus.send)
	other.scheduler = engine.NewSemanticReadScheduler()
	otherView := fillFamilyCodebooks(other, context.Background(), "controller-720")
	assertEnumeration(t, otherView.Enumerations, 1, 3, "u16", []byte{1, 0}, 1, "time-controlled", "Hypothesis")
	assertEnumeration(t, otherView.Enumerations, 2, 2, "u16", []byte{5, 0}, 5, catalogueUnresolved, "Hypothesis")
	assertEnumeration(t, otherView.Enumerations, 3, 6, "u16", []byte{3, 0}, 3, catalogueUnresolved, "Hypothesis")
	assertEnumeration(t, otherView.Enumerations, 2, 29, "u8", []byte{2}, 2, catalogueUnresolved, "Hypothesis")
	assertEnumeration(t, otherView.Enumerations, 3, 0x000E, "u16", []byte{4, 0}, 4, catalogueUnresolved, "Unknown")
	if otherView.Ventilation == nil || otherView.Ventilation.NightFanSpeedLimit == nil || *otherView.Ventilation.NightFanSpeedLimit != 400 {
		t.Fatalf("controller-720 ventilation = %+v; want night word 400", otherView.Ventilation)
	}
	minutes720 := scalarAt(otherView.Scalars, 0, 0x0004)
	if minutes720.Value == nil || *minutes720.Value != 300 || minutes720.Unit != "min" {
		t.Fatalf("controller-720 scalar 0x0004 = %+v; want uint16 300 min", minutes720)
	}
	if got, want := len(otherView.Scalars), catalogueScalarCount("controller-720"); got != want {
		t.Fatalf("controller-720 scalar count = %d; want %d", got, want)
	}

	bus.payloads[[2]uint16{1, 13}] = []byte{9}
	unknownCode, _ := newCataloguePoller(bus.send)
	unknownCode.scheduler = engine.NewSemanticReadScheduler()
	unknownView := fillFamilyCodebooks(unknownCode, context.Background(), "controller-700")
	assertEnumeration(t, unknownView.Enumerations, 1, 13, "u8", []byte{9}, 9, catalogueUnresolved, "Hypothesis")

	unknownFamily, _ := newCataloguePoller(bus.send)
	unknownFamily.scheduler = engine.NewSemanticReadScheduler()
	absent := fillFamilyCodebooks(unknownFamily, context.Background(), "controller-absent")
	for _, item := range absent.Enumerations {
		if item.Label != catalogueUnresolved || item.Level != "Unknown" || len(item.Raw) == 0 {
			t.Fatalf("unknown family enumeration = %+v; want raw bytes labelled unresolved", item)
		}
	}
	if absent.Ventilation != nil {
		t.Fatalf("unknown family ventilation = %+v; want none", absent.Ventilation)
	}
	for _, item := range absent.Scalars {
		if item.Value != nil || len(item.Raw) == 0 {
			t.Fatalf("unknown family scalar = %+v; want raw bytes and no decoded value", item)
		}
	}
}

func TestVaillantB524CodebooksRejectsMutatedIdentityAxis(t *testing.T) {
	base := familyCodebookKey("controller-700", "u16", "00", 1, 3)
	raw := []byte{1, 0}
	if labelFamilyEnumeration(base, raw).Label != "automatic" {
		t.Fatalf("unmutated codebook label = %q; want automatic", labelFamilyEnumeration(base, raw).Label)
	}
	mutations := []struct {
		field  string
		mutate func(*catalogvaillant.B524FamilyCodebookKey)
	}{
		{"family", func(key *catalogvaillant.B524FamilyCodebookKey) { key.B524Family = "controller-absent" }},
		{"variant", func(key *catalogvaillant.B524FamilyCodebookKey) { key.B524FamilyCodebookVariant = "absent" }},
		{"primary", func(key *catalogvaillant.B524FamilyCodebookKey) { key.B524FamilyCodebookPrimary = 1 }},
		{"secondary", func(key *catalogvaillant.B524FamilyCodebookKey) {
			value := uint8(1)
			key.B524FamilyCodebookSecondary = &value
		}},
		{"operation", func(key *catalogvaillant.B524FamilyCodebookKey) {
			value := uint8(1)
			key.B524FamilyCodebookOperation = &value
		}},
		{"direction", func(key *catalogvaillant.B524FamilyCodebookKey) { key.B524FamilyCodebookDirection = "absent" }},
		{"field", func(key *catalogvaillant.B524FamilyCodebookKey) { key.B524FamilyCodebookField = "absent" }},
		{"alternative", func(key *catalogvaillant.B524FamilyCodebookKey) { key.B524FamilyCodebookAlternative = "u8" }},
		{"group", func(key *catalogvaillant.B524FamilyCodebookKey) {
			value := uint8(9)
			key.B524FamilyCodebookGroup = &value
		}},
		{"instance_scope", func(key *catalogvaillant.B524FamilyCodebookKey) {
			key.B524FamilyCodebookInstanceScope = "absent"
		}},
		{"register", func(key *catalogvaillant.B524FamilyCodebookKey) {
			value := uint16(99)
			key.B524FamilyCodebookRegister = &value
		}},
	}
	for _, mutation := range mutations {
		key := cloneCodebookKey(base)
		mutation.mutate(&key)
		if labelFamilyEnumeration(key, raw).Label == "automatic" {
			t.Fatalf("%s mutation still labelled automatic; dropping %s from equality would keep the original row", mutation.field, mutation.field)
		}
	}
}

func newCataloguePoller(send func(context.Context, protocol.Frame) (*protocol.Frame, error)) (*vaillantSemanticPoller, *retainedSemanticState) {
	provider := newRetainedSemanticState()
	poller := &vaillantSemanticPoller{
		controller:     0x15,
		source:         0x7F,
		publication:    newVaillantSemanticPublication(provider, nil, nil),
		scheduler:      engine.NewSemanticReadScheduler(),
		requestTimeout: 50 * time.Millisecond,
		nowFn:          func() time.Time { return time.Date(2024, 3, 15, 8, 0, 0, 123456789, time.UTC) },
		pauseFn:        func(context.Context, time.Duration) bool { return true },
		sendFrameFn:    send,
	}
	return poller, provider
}

type catalogueReplyBus struct {
	payloads map[[2]uint16][]byte
	writes   int
}

func (bus *catalogueReplyBus) send(_ context.Context, frame protocol.Frame) (*protocol.Frame, error) {
	if len(frame.Data) >= 2 && frame.Data[1] == vaillantB524OpWrite {
		bus.writes++
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
	return testB524ResponseForSelectorPayload(frame.Data, payload...), nil
}

func cataloguePayloads() map[[2]uint16][]byte {
	payloads := map[[2]uint16][]byte{
		{1, 3}:      {1, 0},
		{1, 13}:     {6},
		{2, 2}:      {5, 0},
		{2, 29}:     {2},
		{3, 6}:      {3, 0},
		{3, 27}:     {1, 0},
		{3, 0x000E}: {4, 0},
		{9, 0x000D}: {0x2C, 0x01},
		{9, 0x000E}: {0x90, 0x01},
		{0, 0x0004}: {0x2C, 0x01},
		{0, 0x0006}: {0x2C, 0x01},
		{0, 0x0017}: {0x2C, 0x01},
		{0, 0x0018}: {0x2C, 0x01},
		{0, 0x0042}: {0x2C, 0x01},
		{0, 0x0043}: {0x90, 0x01},
	}
	for _, spec := range unreadScalarSpecs {
		key := [2]uint16{uint16(spec.group), spec.register}
		if _, ok := payloads[key]; !ok {
			payloads[key] = catalogueFloatRaw()
		}
	}
	return payloads
}

func catalogueFloatRaw() []byte {
	return testB524Float32Payload(21.5)
}

func catalogueScalarCount(family string) int {
	count := 0
	for _, spec := range unreadScalarSpecs {
		if spec.families == "both" || spec.families == family {
			count++
		}
	}
	return count
}

func assertEnumeration(t *testing.T, items []output.CatalogueEnumeration, group byte, register uint16, alternative string, raw []byte, code int, label, level string) {
	t.Helper()
	item := enumerationAt(items, group, register, alternative)
	if item.Code == nil || *item.Code != code || item.Label != label || item.Level != level || !bytes.Equal(item.Raw, raw) {
		t.Fatalf("group %d register 0x%04x %s = %+v; want code %d label %q level %q raw %x", group, register, alternative, item, code, label, level, raw)
	}
}

func enumerationAt(items []output.CatalogueEnumeration, group byte, register uint16, alternative string) output.CatalogueEnumeration {
	for _, item := range items {
		if item.Group == group && item.Register == register && item.Alternative == alternative {
			return item
		}
	}
	return output.CatalogueEnumeration{}
}

func scalarAt(items []output.CatalogueScalar, group byte, register uint16) output.CatalogueScalar {
	for _, item := range items {
		if item.Group == group && item.Register == register {
			return item
		}
	}
	return output.CatalogueScalar{}
}

func cloneCodebookKey(key catalogvaillant.B524FamilyCodebookKey) catalogvaillant.B524FamilyCodebookKey {
	if key.B524FamilyCodebookSecondary != nil {
		value := *key.B524FamilyCodebookSecondary
		key.B524FamilyCodebookSecondary = &value
	}
	if key.B524FamilyCodebookOperation != nil {
		value := *key.B524FamilyCodebookOperation
		key.B524FamilyCodebookOperation = &value
	}
	if key.B524FamilyCodebookGroup != nil {
		value := *key.B524FamilyCodebookGroup
		key.B524FamilyCodebookGroup = &value
	}
	if key.B524FamilyCodebookRegister != nil {
		value := *key.B524FamilyCodebookRegister
		key.B524FamilyCodebookRegister = &value
	}
	return key
}
