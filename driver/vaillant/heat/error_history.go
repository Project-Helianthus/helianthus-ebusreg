package heat

import (
	b503 "github.com/Project-Helianthus/helianthus-ebusgo/protocol/vaillant/b503"
	semreg "github.com/Project-Helianthus/helianthus-semreg/semreg/v1"
)

// NamespaceMember is one ebusreg name served with the same keys on GraphQL and MCP.
type NamespaceMember struct {
	Namespace string
	Name      string
	Keys      []string
}

// HistoryEntry is one error-history or service-history record. Empty slots are nil.
type HistoryEntry struct {
	Index            int
	FirstActiveError *int
	Slots            []*int
}

// HistoryFailure is the first failing index of an error-history list.
type HistoryFailure struct {
	Code    string
	Index   int
	Message string
}

// HistoryList is the error-history list served on both surfaces.
type HistoryList struct {
	Failure *HistoryFailure
	Records []HistoryEntry
}

// CurrentErrorRead is one current-error reply. Present means the device was
// addressed. Admitted means the reply decoded.
type CurrentErrorRead struct {
	Present  bool
	Admitted bool
	Slots    b503.ErrorSlots
}

// ErrorHistoryNamespace lists the error and service history names. Record
// members use index, first_active_error and slots. The list uses failure and records.
func ErrorHistoryNamespace() []NamespaceMember {
	record := []string{"index", "first_active_error", "slots"}
	list := []string{"failure", "records"}
	return []NamespaceMember{
		{Namespace: "ebusreg", Name: "vaillant_error_history", Keys: record},
		{Namespace: "ebusreg", Name: "vaillant_service_history", Keys: record},
		{Namespace: "ebusreg", Name: "vaillant_errors_history", Keys: list},
		{Namespace: "ebusreg", Name: "ebus_v1_vaillant_errors_history_get", Keys: record},
		{Namespace: "ebusreg", Name: "ebus_v1_vaillant_service_history_get", Keys: record},
		{Namespace: "ebusreg", Name: "ebus_v1_vaillant_errors_history_list", Keys: list},
	}
}

// HistoryEntryFromPayload decodes an error-history reply into the shared record.
func HistoryEntryFromPayload(payload []byte) (HistoryEntry, error) {
	record, err := b503.DecodeErrorHistory(payload)
	if err != nil {
		return HistoryEntry{}, err
	}
	return historyEntry(record), nil
}

// ServiceHistoryEntryFromPayload decodes a service-history reply into the shared record.
func ServiceHistoryEntryFromPayload(payload []byte) (HistoryEntry, error) {
	record, err := b503.DecodeServiceHistory(payload)
	if err != nil {
		return HistoryEntry{}, err
	}
	return historyEntry(record), nil
}

func historyEntry(record b503.ErrorHistoryRecord) HistoryEntry {
	entry := HistoryEntry{Index: int(record.Index), Slots: make([]*int, len(record.Slots))}
	for i, value := range record.Slots {
		if value == b503.EmptySlot {
			continue
		}
		decoded := int(value)
		entry.Slots[i] = &decoded
		if entry.FirstActiveError == nil {
			first := decoded
			entry.FirstActiveError = &first
		}
	}
	return entry
}

// GraphQL returns the history record keys.
func (e HistoryEntry) GraphQL() map[string]any {
	slots := make([]any, len(e.Slots))
	for i, value := range e.Slots {
		if value != nil {
			slots[i] = *value
		}
	}
	return map[string]any{"index": e.Index, "first_active_error": e.FirstActiveError, "slots": slots}
}

// MCP returns the same history record keys as GraphQL.
func (e HistoryEntry) MCP() map[string]any {
	return e.GraphQL()
}

// GraphQL returns the history list keys.
func (l HistoryList) GraphQL() map[string]any {
	records := make([]any, len(l.Records))
	for i, record := range l.Records {
		records[i] = record.GraphQL()
	}
	var failure any
	if l.Failure != nil {
		failure = map[string]any{"code": l.Failure.Code, "index": l.Failure.Index, "message": l.Failure.Message}
	}
	return map[string]any{"failure": failure, "records": records}
}

// MCP returns the same history list keys as GraphQL.
func (l HistoryList) MCP() map[string]any {
	return l.GraphQL()
}

// historyFaultEventIDs is empty: the pinned thermal pack defines no fault
// event record whose subject is a history entry.
var historyFaultEventIDs []semreg.DefinitionID

// HistoryFaultEvents publishes one fault event per history entry when the
// pinned pack defines that event. It publishes none otherwise.
func HistoryFaultEvents(events []semreg.DefinitionID, entries int) []semreg.DefinitionID {
	var matched semreg.DefinitionID
	found := false
	for _, id := range events {
		for _, allowed := range historyFaultEventIDs {
			if id == allowed {
				matched = id
				found = true
			}
		}
	}
	if !found || entries <= 0 {
		return nil
	}
	out := make([]semreg.DefinitionID, entries)
	for i := range out {
		out[i] = matched
	}
	return out
}

// FaultFromCurrentError reports present when an admitted current-error reply
// has an active slot, and clear when every addressed reply has none. A missing
// reply is not a clear reading.
func FaultFromCurrentError(controller, boiler CurrentErrorRead) (string, bool) {
	failed := false
	active := false
	seen := false
	for _, read := range []CurrentErrorRead{controller, boiler} {
		if !read.Present {
			continue
		}
		seen = true
		if !read.Admitted {
			failed = true
			continue
		}
		if _, ok := read.Slots.FirstActive(); ok {
			active = true
		}
	}
	if !seen {
		return "", false
	}
	if active {
		return "present", true
	}
	if failed {
		return "", false
	}
	return "clear", true
}
