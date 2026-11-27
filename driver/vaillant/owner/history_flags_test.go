package owner

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
	b503 "github.com/Project-Helianthus/helianthus-ebusgo/protocol/vaillant/b503"

	catalogvaillant "github.com/Project-Helianthus/helianthus-ebusreg/catalog/vaillant"
	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/output"
	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/owner/b503session"
)

func TestVaillantB503HistoryFlagsStock(t *testing.T) {
	dispatcher, bus, _ := newTestDispatcher(t)
	bus.setResp([2]byte{0x01, 0x01}, historyBody(0, 0x0010))
	bare, err := readErrorHistory(context.Background(), dispatcher, 0x08, nil)
	if err != nil {
		t.Fatal(err)
	}
	bareFrame, ok := bus.lastCall()
	if !ok || !reflect.DeepEqual(bareFrame.Data, b503.EncodeErrorHistory()) {
		t.Fatalf("bare request = %x", bareFrame.Data)
	}
	zero := uint8(0)
	explicit, err := readErrorHistory(context.Background(), dispatcher, 0x08, &zero)
	if err != nil {
		t.Fatal(err)
	}
	explicitFrame, _ := bus.lastCall()
	if reflect.DeepEqual(explicitFrame.Data, bareFrame.Data) {
		t.Fatal("explicit zero matched the bare selector; dropping the index byte would hide the mutation")
	}
	if !reflect.DeepEqual(explicitFrame.Data, append(append([]byte(nil), b503.EncodeErrorHistory()...), 0)) {
		t.Fatalf("explicit zero = %x", explicitFrame.Data)
	}
	assertStockUnknown(t, bare)
	assertStockUnknown(t, explicit)
	if bare.ReturnedIndex == nil || *bare.ReturnedIndex != 0 || len(bare.Slots) != 5 || bare.Slots[0].Value != 0x0010 || bare.Slots[0].Empty {
		t.Fatalf("bare slots = %+v index %v", bare.Slots, bare.ReturnedIndex)
	}
	if !bare.Slots[1].Empty || bare.Slots[1].Value != b503.EmptySlot {
		t.Fatalf("empty slot = %+v", bare.Slots[1])
	}
	serviceBus := newB503DispatcherMockBus()
	serviceBus.setResp([2]byte{0x01, 0x02}, historyBody(0, 0x0020))
	service := historyDispatcher(t, serviceBus)
	read, err := readServiceHistory(context.Background(), service, 0x08, nil)
	if err != nil {
		t.Fatal(err)
	}
	frame, _ := serviceBus.lastCall()
	if !reflect.DeepEqual(frame.Data, b503.EncodeServiceHistory()) {
		t.Fatalf("service selector = %x", frame.Data)
	}
	indexed := uint8(0)
	if _, err := readServiceHistory(context.Background(), service, 0x08, &indexed); err != nil {
		t.Fatal(err)
	}
	indexedFrame, _ := serviceBus.lastCall()
	if reflect.DeepEqual(indexedFrame.Data, frame.Data) {
		t.Fatal("service index mutation was ignored")
	}
	assertStockUnknown(t, read)
	if read.Kind != historyKindService || read.Slots[0].Value != 0x0020 {
		t.Fatalf("service read = %+v", read)
	}
}

func TestVaillantB503HistoryFlagsLists(t *testing.T) {
	bus := &recordingHistoryBus{respond: echoHistory}
	list, err := readErrorHistoryList(context.Background(), historyDispatcher(t, bus), 0x08, nil)
	if err != nil || list.Limit != historyListDefault || len(list.Items) != historyListDefault || list.Stop != nil {
		t.Fatalf("default list items=%d limit=%d err=%v stop=%+v", len(list.Items), list.Limit, err, list.Stop)
	}
	payloads := bus.payloads()
	if len(payloads) != historyListDefault {
		t.Fatalf("default calls = %d", len(payloads))
	}
	for index, payload := range payloads {
		want := append(append([]byte(nil), b503.EncodeErrorHistory()...), byte(index))
		if !reflect.DeepEqual(payload, want) {
			t.Fatalf("index %d payload = %x", index, payload)
		}
		if list.Items[index].FlagsKnown || list.Items[index].CategoryKnown {
			t.Fatal("list item inferred qualifiers from the stock reply")
		}
	}
	bound := 1
	bus.calls = nil
	bounded, err := readErrorHistoryList(context.Background(), historyDispatcher(t, bus), 0x08, &bound)
	if err != nil || len(bounded.Items) != 1 || len(bus.payloads()) != 1 {
		t.Fatalf("bound list items=%d calls=%d err=%v", len(bounded.Items), len(bus.payloads()), err)
	}
	for _, limit := range []int{0, 17, -1} {
		before := len(bus.payloads())
		_, err := readErrorHistoryList(context.Background(), historyDispatcher(t, bus), 0x08, &limit)
		if !errors.Is(err, errHistoryLimit) || len(bus.payloads()) != before {
			t.Fatalf("limit %d err=%v calls changed", limit, err)
		}
	}
}

func TestVaillantB503HistoryFlagsFailure(t *testing.T) {
	bus := &recordingHistoryBus{respond: func(payload []byte) ([]byte, error) {
		if len(payload) >= 3 && payload[2] == 2 {
			return nil, errors.New("history transport down")
		}
		return historyBody(payload[2], 0x0010), nil
	}}
	list, err := readErrorHistoryList(context.Background(), historyDispatcher(t, bus), 0x08, nil)
	if err == nil || len(list.Items) != 2 || list.Stop == nil || list.Stop.Index != 2 || list.Stop.Code != historyCodeTransport {
		t.Fatalf("failure list items=%d stop=%+v err=%v", len(list.Items), list.Stop, err)
	}
	short := &recordingHistoryBus{respond: func(payload []byte) ([]byte, error) {
		if len(payload) >= 3 && payload[2] == 1 {
			return []byte{0x01}, nil
		}
		return historyBody(payload[2], 0x0011), nil
	}}
	stopped, err := readErrorHistoryList(context.Background(), historyDispatcher(t, short), 0x08, nil)
	if !errors.Is(err, b503.ErrShortPayload) || len(stopped.Items) != 1 || stopped.Stop == nil || stopped.Stop.Code != historyCodeDecode || len(stopped.Stop.Raw) != 1 {
		t.Fatalf("short stop = %+v err=%v", stopped.Stop, err)
	}
	echo := &recordingHistoryBus{respond: func([]byte) ([]byte, error) {
		return historyBody(0, 0x0012), nil
	}}
	mismatched, err := readErrorHistoryList(context.Background(), historyDispatcher(t, echo), 0x08, nil)
	if !errors.Is(err, errHistoryEcho) || len(mismatched.Items) != 1 || mismatched.Stop == nil || mismatched.Stop.Code != historyCodeEcho || mismatched.Stop.Index != 1 {
		t.Fatalf("echo stop = %+v items=%d err=%v", mismatched.Stop, len(mismatched.Items), err)
	}
	if mismatched.Stop.DecodedSlot == nil || *mismatched.Stop.DecodedSlot != 0x0012 || len(mismatched.Stop.Slots) != 5 {
		t.Fatalf("echo evidence = %+v", mismatched.Stop)
	}
}

func TestVaillantB503HistoryFlagsMalformed(t *testing.T) {
	dispatcher, bus, _ := newTestDispatcher(t)
	long := make([]byte, 12)
	bus.setResp([2]byte{0x01, 0x01}, long)
	read, err := readErrorHistory(context.Background(), dispatcher, 0x08, nil)
	if !errors.Is(err, b503.ErrWrongPayloadLen) || !reflect.DeepEqual(read.Raw, long) || len(read.Slots) != 0 || read.FlagsKnown || read.CategoryKnown || read.Ambiguity != historyAmbiguityUnknown {
		t.Fatalf("long read err=%v read=%+v", err, read)
	}
	record := historyOwnerRecord([]output.HistoryRead{read}, nil, nil, nil, nil)
	if len(record.Flags) != 0 || len(record.Categories) != 0 || record.Ambiguity != historyAmbiguityUnknown {
		t.Fatal("malformed stock gained a qualifier")
	}
}

func TestVaillantB503HistoryFlagsAlternatives(t *testing.T) {
	noItem := historyBranchInput(0)
	noItem.B503HistoryBranchRemainderPresent = true
	noItem.B503HistoryBranchRemainder = [7]uint8{255, 255, 255, 255, 255, 255, 255}
	if state := qualifyHistoryBranch(noItem).State; state != "no-item" {
		t.Fatalf("seven ff state = %s", state)
	}
	codeOnly := historyBranchInput(1)
	codeOnly.B503HistoryBranchErrorPresent = true
	codeOnly.B503HistoryBranchError = 42
	for index := range codeOnly.B503HistoryBranchCalendar {
		codeOnly.B503HistoryBranchCalendar[index].B503HistoryBranchRawPresent = true
		codeOnly.B503HistoryBranchCalendar[index].B503HistoryBranchRaw = 255
	}
	codeReading := qualifyHistoryBranch(codeOnly)
	if codeReading.State != "error-without-complete-time" {
		t.Fatalf("code-only state = %s", codeReading.State)
	}
	for _, component := range codeReading.Components {
		if component.State != "replacement" {
			t.Fatalf("replacement component = %s", component.State)
		}
	}
	ranged := historyBranchInput(2)
	values := []uint64{10, 24, 15, 6, 26}
	for index, value := range values {
		ranged.B503HistoryBranchCalendar[index] = catalogvaillant.B503HistoryBranchCell{
			B503HistoryBranchRawPresent:        true,
			B503HistoryBranchRaw:               byte(value),
			B503HistoryBranchNormalizedPresent: true,
			B503HistoryBranchValue:             value,
		}
	}
	rangedReading := qualifyHistoryBranch(ranged)
	if rangedReading.State != "candidate-components" || rangedReading.Complete || rangedReading.Components[1].State != "out-of-domain" || rangedReading.Components[0].State != "candidate" {
		t.Fatalf("out-of-range branch = %+v", rangedReading)
	}
	flags, found, err := qualifyHistoryField(historyFlagKey(t, "history_flags"), 1<<0|1<<5)
	if err != nil || !found {
		t.Fatalf("flags found=%v err=%v", found, err)
	}
	if !flags.Bits[0].Set || flags.Bits[0].Name != "history_flag_0" || flags.Bits[0].Reserved || flags.Bits[1].Set || !flags.Bits[5].Set || !flags.Bits[5].Reserved || flags.Bits[5].Name != "" {
		t.Fatalf("flag bits = %+v", flags.Bits)
	}
	common, found, err := qualifyHistoryCategory(historyCategoryKey(t, "common"), 2)
	if err != nil || !found || common.Label != "predictive maintenance" || !common.Labelled || common.Level != "Hypothesis" {
		t.Fatalf("common category = %+v found=%v err=%v", common, found, err)
	}
	reset, found, err := qualifyHistoryCategory(historyCategoryKey(t, "reset"), 2)
	if err != nil || !found || reset.Label != "" || reset.Labelled || reset.Raw != 2 || !reset.Unresolved {
		t.Fatalf("reset category = %+v found=%v err=%v", reset, found, err)
	}
	if common.Ambiguity != reset.Ambiguity || common.Ambiguity != historyAmbiguityUnknown {
		t.Fatal("category surfaces would need a second disposition")
	}
}

func TestVaillantB503HistoryFlagsIdentity(t *testing.T) {
	catalog, err := catalogvaillant.B503HistoryCatalog()
	if err != nil || len(catalog.B503HistoryFlagFields) == 0 {
		t.Fatal(err)
	}
	key := catalog.B503HistoryFlagFields[0].B503HistoryFlagKey
	if _, found, err := historyField(key); err != nil || !found {
		t.Fatalf("unmutated history field found=%v err=%v", found, err)
	}
	mutations := []struct {
		field  string
		mutate func(*catalogvaillant.B503HistoryFlagKey)
	}{
		{"family", func(key *catalogvaillant.B503HistoryFlagKey) { key.B503HistoryFlagFamily = "absent-family" }},
		{"variant", func(key *catalogvaillant.B503HistoryFlagKey) {
			key.B503HistoryFlagVariant = mutateString(key.B503HistoryFlagVariant)
		}},
		{"primary", func(key *catalogvaillant.B503HistoryFlagKey) { key.B503HistoryFlagPrimary++ }},
		{"secondary", func(key *catalogvaillant.B503HistoryFlagKey) {
			key.B503HistoryFlagSecondary = mutateOptionalOctet(key.B503HistoryFlagSecondary)
		}},
		{"operation", func(key *catalogvaillant.B503HistoryFlagKey) {
			key.B503HistoryFlagOperation = mutateOptionalOctet(key.B503HistoryFlagOperation)
		}},
		{"direction", func(key *catalogvaillant.B503HistoryFlagKey) {
			key.B503HistoryFlagDirection = mutateString(key.B503HistoryFlagDirection)
		}},
		{"field", func(key *catalogvaillant.B503HistoryFlagKey) { key.B503HistoryFlagField = "absent-field" }},
		{"opcode", func(key *catalogvaillant.B503HistoryFlagKey) {
			key.B503HistoryFlagOpcode = mutateOptionalOctet(key.B503HistoryFlagOpcode)
		}},
		{"group", func(key *catalogvaillant.B503HistoryFlagKey) {
			key.B503HistoryFlagGroup = mutateOptionalOctet(key.B503HistoryFlagGroup)
		}},
		{"instance", func(key *catalogvaillant.B503HistoryFlagKey) {
			key.B503HistoryFlagInstance = mutateOptionalOctet(key.B503HistoryFlagInstance)
		}},
		{"register", func(key *catalogvaillant.B503HistoryFlagKey) {
			key.B503HistoryFlagRegister = mutateOptionalWide(key.B503HistoryFlagRegister)
		}},
	}
	for _, mutation := range mutations {
		mutated := cloneJSON(t, key)
		mutation.mutate(&mutated)
		if _, found, err := historyField(mutated); err != nil || found {
			t.Fatalf("%s mutation still matched; dropping %s from equality would keep the original row", mutation.field, mutation.field)
		}
	}
	entries, err := catalogvaillant.B503TextCategoryEntries()
	if err != nil || len(entries) == 0 {
		t.Fatal(err)
	}
	category := entries[0].B503TextCategoryKey
	if _, found, err := qualifyHistoryCategory(category, 0); err != nil || !found {
		t.Fatalf("unmutated category found=%v err=%v", found, err)
	}
	categoryMutations := []struct {
		field  string
		mutate func(*catalogvaillant.B503TextCategoryKey)
	}{
		{"family", func(key *catalogvaillant.B503TextCategoryKey) { key.B503TextCategoryFamily = "absent-family" }},
		{"variant", func(key *catalogvaillant.B503TextCategoryKey) { key.B503TextCategoryVariant = "absent-variant" }},
		{"primary", func(key *catalogvaillant.B503TextCategoryKey) { key.B503TextCategoryPrimary++ }},
		{"secondary", func(key *catalogvaillant.B503TextCategoryKey) {
			key.B503TextCategorySecondary = mutateOptionalOctet(key.B503TextCategorySecondary)
		}},
		{"operation", func(key *catalogvaillant.B503TextCategoryKey) {
			key.B503TextCategoryOperation = mutateOptionalOctet(key.B503TextCategoryOperation)
		}},
		{"direction", func(key *catalogvaillant.B503TextCategoryKey) {
			key.B503TextCategoryDirection = mutateString(key.B503TextCategoryDirection)
		}},
		{"field", func(key *catalogvaillant.B503TextCategoryKey) { key.B503TextCategoryField = "absent-field" }},
		{"group", func(key *catalogvaillant.B503TextCategoryKey) {
			key.B503TextCategoryGroup = mutateOptionalOctet(key.B503TextCategoryGroup)
		}},
		{"requested_instance", func(key *catalogvaillant.B503TextCategoryKey) {
			key.B503TextCategoryRequestedInstance = mutateOptionalOctet(key.B503TextCategoryRequestedInstance)
		}},
		{"register", func(key *catalogvaillant.B503TextCategoryKey) {
			key.B503TextCategoryRegister = mutateOptionalWide(key.B503TextCategoryRegister)
		}},
		{"command_role", func(key *catalogvaillant.B503TextCategoryKey) {
			key.B503TextCategoryCommandRole = mutateString(key.B503TextCategoryCommandRole)
		}},
	}
	for _, mutation := range categoryMutations {
		mutated := cloneJSON(t, category)
		mutation.mutate(&mutated)
		if _, found, err := qualifyHistoryCategory(mutated, 0); err != nil || found {
			t.Fatalf("%s mutation still matched; dropping %s from equality would keep the original row", mutation.field, mutation.field)
		}
	}
	noItem := historyBranchInput(0)
	noItem.B503HistoryBranchRemainderPresent = true
	noItem.B503HistoryBranchRemainder = [7]uint8{255, 255, 255, 255, 255, 255, 255}
	cleared := noItem
	cleared.B503HistoryBranchContext.B503HistoryBranchFamily = ""
	if qualifyHistoryBranch(cleared).State != "unknown-branch" {
		t.Fatal("family mutation still selected the no-item branch")
	}
	cleared = noItem
	cleared.B503HistoryBranchContext.B503HistoryBranchVariant = ""
	if qualifyHistoryBranch(cleared).State != "unknown-branch" {
		t.Fatal("variant mutation still selected the no-item branch")
	}
	cleared = noItem
	cleared.B503HistoryBranchContext.B503HistoryBranchOperation = 0
	if qualifyHistoryBranch(cleared).State != "unknown-branch" {
		t.Fatal("operation mutation still selected the no-item branch")
	}
	cleared = noItem
	cleared.B503HistoryBranchContext.B503HistoryBranchDirection = ""
	if qualifyHistoryBranch(cleared).State != "unknown-branch" {
		t.Fatal("direction mutation still selected the no-item branch")
	}
}

func TestVaillantB503HistoryFlagsSurfaces(t *testing.T) {
	read, err := readErrorHistory(context.Background(), mustHistoryDispatcher(t, historyBody(0, 0x0010)), 0x08, nil)
	if err != nil {
		t.Fatal(err)
	}
	flag, _, err := qualifyHistoryField(historyFlagKey(t, "history_flags"), 1)
	if err != nil {
		t.Fatal(err)
	}
	record := historyOwnerRecord([]output.HistoryRead{read}, nil, []output.HistoryFlagReading{flag}, nil, nil)
	if record.Reads[0].FlagsKnown || len(record.Flags) != 1 || record.Ambiguity != historyAmbiguityUnknown {
		t.Fatal("stock read absorbed the separate flag qualification")
	}
	if !reflect.DeepEqual(record.GraphQL(), record.MCP()) || record.GraphQL().Ambiguity != record.MCP().Ambiguity {
		t.Fatal("history surfaces diverged")
	}
	status := &output.SystemStatus{}
	mcp := bindHistorySurfaces(status, record)
	if !reflect.DeepEqual(status.History(), mcp.History()) {
		t.Fatal("system history diverged")
	}
	poller := &vaillantSemanticPoller{}
	retainHistory(poller, record)
	if !reflect.DeepEqual(retainedHistory(poller), status.History()) {
		t.Fatal("retained history diverged")
	}
}

func assertStockUnknown(t *testing.T, read output.HistoryRead) {
	t.Helper()
	if read.FlagsKnown || read.CategoryKnown || read.Ambiguity != historyAmbiguityUnknown || len(read.Raw) != 11 {
		t.Fatalf("stock qualifier changed: %+v", read)
	}
}

func historyBody(index byte, first uint16) []byte {
	body := []byte{index, byte(first), byte(first >> 8)}
	for slot := 0; slot < 4; slot++ {
		body = append(body, 0xFF, 0xFF)
	}
	return body
}

func echoHistory(payload []byte) ([]byte, error) {
	index := byte(0)
	if len(payload) >= 3 {
		index = payload[2]
	}
	return historyBody(index, 0x0010), nil
}

func historyBranchInput(marker uint8) catalogvaillant.B503HistoryBranchFields {
	return catalogvaillant.B503HistoryBranchFields{
		B503HistoryBranchContext: catalogvaillant.B503HistoryBranchContext{
			B503HistoryBranchFamily:    "B503",
			B503HistoryBranchVariant:   "history",
			B503HistoryBranchOperation: 1,
			B503HistoryBranchDirection: "field",
		},
		B503HistoryBranchMarkerPresent: true,
		B503HistoryBranchMarker:        marker,
	}
}

func historyFlagKey(t *testing.T, field string) catalogvaillant.B503HistoryFlagKey {
	t.Helper()
	catalog, err := catalogvaillant.B503HistoryCatalog()
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range catalog.B503HistoryFlagFields {
		if item.B503HistoryFlagKey.B503HistoryFlagField == field {
			return item.B503HistoryFlagKey
		}
	}
	t.Fatalf("missing history field %s", field)
	return catalogvaillant.B503HistoryFlagKey{}
}

func historyCategoryKey(t *testing.T, variant string) catalogvaillant.B503TextCategoryKey {
	t.Helper()
	entries, err := catalogvaillant.B503TextCategoryEntries()
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if entry.B503TextCategoryKey.B503TextCategoryField == "categories" && entry.B503TextCategoryKey.B503TextCategoryVariant == variant {
			return entry.B503TextCategoryKey
		}
	}
	t.Fatalf("missing category variant %s", variant)
	return catalogvaillant.B503TextCategoryKey{}
}

func historyDispatcher(t *testing.T, bus b503Bus) *RawFrameDispatcher {
	t.Helper()
	manager := b503session.New(
		b503session.TransportKey{AdapterInstanceID: "test", TransportEpoch: 1},
		30*time.Second,
		func(context.Context) (b503session.TransportKey, error) {
			return b503session.TransportKey{}, b503session.ErrTransportDown
		},
	)
	var readMu sync.Mutex
	return NewRawFrameDispatcher(bus, adapterSource, &readMu, manager, time.Second)
}

func mustHistoryDispatcher(t *testing.T, body []byte) *RawFrameDispatcher {
	t.Helper()
	dispatcher, bus, _ := newTestDispatcher(t)
	bus.setResp([2]byte{0x01, 0x01}, body)
	return dispatcher
}

type recordingHistoryBus struct {
	mu      sync.Mutex
	calls   [][]byte
	respond func([]byte) ([]byte, error)
}

func (b *recordingHistoryBus) Send(ctx context.Context, frame protocol.Frame) (*protocol.Frame, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	b.mu.Lock()
	b.calls = append(b.calls, append([]byte(nil), frame.Data...))
	respond := b.respond
	b.mu.Unlock()
	if respond == nil {
		return &protocol.Frame{}, nil
	}
	data, err := respond(frame.Data)
	if err != nil {
		return nil, err
	}
	return &protocol.Frame{Data: data}, nil
}

func (b *recordingHistoryBus) payloads() [][]byte {
	b.mu.Lock()
	defer b.mu.Unlock()
	out := make([][]byte, len(b.calls))
	for i, call := range b.calls {
		out[i] = append([]byte(nil), call...)
	}
	return out
}

func cloneJSON[T any](t *testing.T, value T) T {
	t.Helper()
	raw, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	var out T
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatal(err)
	}
	return out
}

func mutateOptionalOctet(current *uint8) *uint8 {
	if current == nil {
		value := uint8(0)
		return &value
	}
	value := *current + 1
	return &value
}

func mutateOptionalWide(current *uint16) *uint16 {
	if current == nil {
		value := uint16(0)
		return &value
	}
	value := *current + 1
	return &value
}

func mutateString(current *string) *string {
	value := "absent"
	if current != nil && *current == value {
		value = "other-absent"
	}
	return &value
}
