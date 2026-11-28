package heat

import (
	"encoding/binary"
	"reflect"
	"testing"

	b503 "github.com/Project-Helianthus/helianthus-ebusgo/protocol/vaillant/b503"
	semreg "github.com/Project-Helianthus/helianthus-semreg/semreg/v1"
	"github.com/Project-Helianthus/helianthus-semreg/semreg/v1/packs/thermal"
)

// Identity fields, and the production mutation each assertion catches:
//   - namespace: every member is ebusreg. Another namespace keeps the keys off the driver namespace.
//   - error member name: vaillant_error_history. Reusing the service name serves the error record under the service name.
//   - service member name: vaillant_service_history. Reusing the error name drops the service record.
//   - list member name: vaillant_errors_history. Giving it the record keys serves the list as one record.
//   - MCP get names: ebus_v1_vaillant_errors_history_get and ebus_v1_vaillant_service_history_get. Swapping them serves each record on the other tool.
//   - MCP list name: ebus_v1_vaillant_errors_history_list. Giving it record keys drops failure and records.
//   - record keys: index, first_active_error, slots. A missing or renamed key changes the served record.
//   - list keys: failure and records. Using the record keys on the list drops the failure.
//   - failure keys: code, index, message. Swapping code and message serves the text as the code.
//   - history index: the error record index is 3 and the service record index is 7. A fixed index keeps 0.
//   - first active slot: an empty first slot is nil and the later slot 281 is first. Reading slot 0 as active serves 65535.
//   - service first slot: slot 0 value 42 is first. Using the error record keeps index 3 and value 281.
//   - empty failure: the list failure stays nil. An invented failure map adds code, index and message.
//   - current-error slot: an active slot after the first is present. Reading only slot 0 reports clear.
//   - admitted clear: both empty replies are clear. Treating an empty slot as active reports present.
//   - failed reply: a reply that did not decode is not clear. Reporting clear replaces the last fault.
//   - absent device: a device that was not addressed is not a failed read. Requiring both devices leaves a single clear reply unknown.
//   - fault events: the pinned thermal events publish no history event. Publishing valve_fault or detector_sensor_fault adds one event per entry.

func TestErrorHistoryNamespace(t *testing.T) {
	record := []string{"index", "first_active_error", "slots"}
	list := []string{"failure", "records"}
	want := []NamespaceMember{
		{Namespace: "ebusreg", Name: "vaillant_error_history", Keys: record},
		{Namespace: "ebusreg", Name: "vaillant_service_history", Keys: record},
		{Namespace: "ebusreg", Name: "vaillant_errors_history", Keys: list},
		{Namespace: "ebusreg", Name: "ebus_v1_vaillant_errors_history_get", Keys: record},
		{Namespace: "ebusreg", Name: "ebus_v1_vaillant_service_history_get", Keys: record},
		{Namespace: "ebusreg", Name: "ebus_v1_vaillant_errors_history_list", Keys: list},
	}
	if !reflect.DeepEqual(ErrorHistoryNamespace(), want) {
		t.Fatalf("namespace = %#v", ErrorHistoryNamespace())
	}

	errorEntry := decodeHistory(t, HistoryEntryFromPayload, 3, [5]uint16{b503.EmptySlot, 281, b503.EmptySlot, b503.EmptySlot, b503.EmptySlot}, 281, 1)
	serviceEntry := decodeHistory(t, ServiceHistoryEntryFromPayload, 7, [5]uint16{42, b503.EmptySlot, b503.EmptySlot, b503.EmptySlot, b503.EmptySlot}, 42, 0)
	if reflect.DeepEqual(errorEntry.GraphQL(), serviceEntry.GraphQL()) {
		t.Fatal("error and service records are the same value")
	}
	if _, err := HistoryEntryFromPayload([]byte{3}); err == nil {
		t.Fatal("short error history decoded")
	}
	if _, err := ServiceHistoryEntryFromPayload([]byte{7}); err == nil {
		t.Fatal("short service history decoded")
	}

	historyList := HistoryList{
		Failure: &HistoryFailure{Code: "unavailable", Index: 4, Message: "history index unread"},
		Records: []HistoryEntry{errorEntry},
	}
	assertSameSurface(t, historyList.GraphQL(), historyList.MCP(), []string{"failure", "records"})
	failure, ok := historyList.GraphQL()["failure"].(map[string]any)
	if !ok || failure["code"] != "unavailable" || failure["index"] != 4 || failure["message"] != "history index unread" || len(failure) != 3 {
		t.Fatalf("failure = %#v", historyList.GraphQL()["failure"])
	}
	empty := HistoryList{Records: []HistoryEntry{serviceEntry}}
	if empty.GraphQL()["failure"] != nil || !reflect.DeepEqual(empty.GraphQL(), empty.MCP()) {
		t.Fatalf("empty failure = %#v", empty.GraphQL())
	}

	active := admittedError(3, 25)
	clear := admittedError(-1, 0)
	failed := CurrentErrorRead{Present: true}
	absent := CurrentErrorRead{}
	assertFault(t, active, absent, "present", true)
	assertFault(t, clear, clear, "clear", true)
	assertFault(t, failed, clear, "", false)
	assertFault(t, failed, active, "present", true)
	assertFault(t, absent, absent, "", false)
	assertFault(t, absent, clear, "clear", true)
	assertFault(t, failed, failed, "", false)
	ignored := admittedError(0, 11)
	ignored.Present = false
	assertFault(t, ignored, absent, "", false)

	ids := thermalEventIDs(t)
	if got := HistoryFaultEvents(ids, 4); len(got) != 0 {
		t.Fatalf("history events = %v, want none", got)
	}
	if got := HistoryFaultEvents(ids, 0); len(got) != 0 {
		t.Fatalf("empty history events = %v", got)
	}
}

func decodeHistory(t *testing.T, decode func([]byte) (HistoryEntry, error), index byte, slots [5]uint16, first int, firstAt int) HistoryEntry {
	t.Helper()
	payload := make([]byte, 11)
	payload[0] = index
	for i, value := range slots {
		binary.LittleEndian.PutUint16(payload[1+2*i:], value)
	}
	entry, err := decode(payload)
	if err != nil {
		t.Fatal(err)
	}
	if entry.Index != int(index) || entry.FirstActiveError == nil || *entry.FirstActiveError != first {
		t.Fatalf("entry index=%d first=%v, want index %d first %d", entry.Index, entry.FirstActiveError, index, first)
	}
	if len(entry.Slots) != 5 || entry.Slots[firstAt] == nil || *entry.Slots[firstAt] != first {
		t.Fatalf("slots = %v", entry.Slots)
	}
	for i, slot := range entry.Slots {
		if i != firstAt && slot != nil {
			t.Fatalf("slot %d = %v, want empty", i, *slot)
		}
	}
	assertSameSurface(t, entry.GraphQL(), entry.MCP(), []string{"index", "first_active_error", "slots"})
	gotSlots, ok := entry.GraphQL()["slots"].([]any)
	if !ok || len(gotSlots) != 5 || !reflect.DeepEqual(gotSlots[firstAt], first) {
		t.Fatalf("served slots = %#v", entry.GraphQL()["slots"])
	}
	for i, slot := range gotSlots {
		if i != firstAt && slot != nil {
			t.Fatalf("served slot %d = %#v, want empty", i, slot)
		}
	}
	return entry
}

func assertSameSurface(t *testing.T, graphql, mcp map[string]any, keys []string) {
	t.Helper()
	if !reflect.DeepEqual(graphql, mcp) {
		t.Fatalf("GraphQL %#v, MCP %#v", graphql, mcp)
	}
	if len(graphql) != len(keys) {
		t.Fatalf("keys = %#v, want %v", graphql, keys)
	}
	for _, key := range keys {
		if _, ok := graphql[key]; !ok {
			t.Fatalf("missing key %s in %#v", key, graphql)
		}
	}
}

func admittedError(active int, code uint16) CurrentErrorRead {
	read := CurrentErrorRead{Present: true, Admitted: true}
	for i := range read.Slots.Slots {
		read.Slots.Slots[i] = b503.EmptySlot
	}
	if active >= 0 {
		read.Slots.Slots[active] = code
	}
	return read
}

func assertFault(t *testing.T, controller, boiler CurrentErrorRead, symbol string, known bool) {
	t.Helper()
	got, ok := FaultFromCurrentError(controller, boiler)
	if ok != known || got != symbol {
		t.Fatalf("fault = %q known=%v, want %q known=%v", got, ok, symbol, known)
	}
}

func thermalEventIDs(t *testing.T) []semreg.DefinitionID {
	t.Helper()
	var ids []semreg.DefinitionID
	sawValve := false
	sawDetector := false
	for _, event := range thermal.Metadata().Events {
		ids = append(ids, event.ID)
		if event.ID == "thermal.event.valve_fault" {
			sawValve = true
		}
		if event.ID == "thermal.event.detector_sensor_fault" {
			sawDetector = true
		}
	}
	if !sawValve || !sawDetector {
		t.Fatal("pinned thermal events do not include valve_fault and detector_sensor_fault")
	}
	return ids
}
