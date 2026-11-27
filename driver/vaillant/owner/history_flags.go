package owner

import (
	"context"
	"errors"
	"sync"

	b503 "github.com/Project-Helianthus/helianthus-ebusgo/protocol/vaillant/b503"

	catalogvaillant "github.com/Project-Helianthus/helianthus-ebusreg/catalog/vaillant"
	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/output"
)

const (
	historyKindError          = "error"
	historyKindService        = "service"
	historyAmbiguityUnknown   = "unknown"
	historyCodeDecode         = "decode"
	historyCodeTransport      = "transport"
	historyCodeEcho           = "echo"
	historyListDefault        = 5
	historyListMinimum        = 1
	historyListMaximum        = 16
	historyQualificationLevel = "Unknown"
)

var (
	errHistoryLimit = errors.New("error history limit outside the closed range")
	errHistoryEcho  = errors.New("history returned index mismatch")
)

func readErrorHistory(ctx context.Context, dispatcher *RawFrameDispatcher, target byte, index *uint8) (output.HistoryRead, error) {
	return readStockHistory(ctx, dispatcher, target, historyKindError, index)
}

func readServiceHistory(ctx context.Context, dispatcher *RawFrameDispatcher, target byte, index *uint8) (output.HistoryRead, error) {
	return readStockHistory(ctx, dispatcher, target, historyKindService, index)
}

func readErrorHistoryList(ctx context.Context, dispatcher *RawFrameDispatcher, target byte, limit *int) (output.HistoryList, error) {
	count, err := historyLimit(limit)
	list := output.HistoryList{Kind: historyKindError, Limit: count}
	if err != nil {
		return output.HistoryList{}, err
	}
	for index := 0; index < count; index++ {
		requested := uint8(index)
		read, readErr := readErrorHistory(ctx, dispatcher, target, &requested)
		if readErr != nil || read.ReturnedIndex == nil || *read.ReturnedIndex != requested {
			list.Stop = historyStop(requested, read, readErr)
			if readErr == nil {
				readErr = errHistoryEcho
			}
			return list, readErr
		}
		list.Items = append(list.Items, read)
	}
	return list, nil
}

func historyLimit(limit *int) (int, error) {
	if limit == nil {
		return historyListDefault, nil
	}
	if *limit < historyListMinimum || *limit > historyListMaximum {
		return 0, errHistoryLimit
	}
	return *limit, nil
}

func readStockHistory(ctx context.Context, dispatcher *RawFrameDispatcher, target byte, kind string, index *uint8) (output.HistoryRead, error) {
	read := output.HistoryRead{Kind: kind, Ambiguity: historyAmbiguityUnknown}
	if index != nil {
		requested := *index
		read.RequestedIndex = &requested
	}
	if ctx == nil {
		ctx = context.Background()
	}
	payload := historyRequest(kind, index)
	response, err := dispatcher.Invoke(ctx, target, payload)
	if err != nil {
		return read, err
	}
	read.Raw = append([]byte(nil), response...)
	record, decErr := decodeStockHistory(kind, response)
	if decErr != nil {
		read.DecodeError = decErr.Error()
		return read, decErr
	}
	returned := record.Index
	read.ReturnedIndex = &returned
	read.Slots = make([]output.HistorySlot, len(record.Slots))
	for slot := range record.Slots {
		value := record.Slots[slot]
		read.Slots[slot] = output.HistorySlot{Value: value, Empty: value == b503.EmptySlot}
	}
	if index != nil && returned != *index {
		return read, errHistoryEcho
	}
	return read, nil
}

func historyRequest(kind string, index *uint8) []byte {
	var payload []byte
	if kind == historyKindService {
		payload = append([]byte(nil), b503.EncodeServiceHistory()...)
	} else {
		payload = append([]byte(nil), b503.EncodeErrorHistory()...)
	}
	if index != nil {
		payload = append(payload, *index)
	}
	return payload
}

func decodeStockHistory(kind string, payload []byte) (b503.ErrorHistoryRecord, error) {
	if kind == historyKindService {
		return b503.DecodeServiceHistory(payload)
	}
	return b503.DecodeErrorHistory(payload)
}

func historyStop(index uint8, read output.HistoryRead, err error) *output.HistoryStop {
	stop := &output.HistoryStop{
		Index:         index,
		Raw:           append([]byte(nil), read.Raw...),
		ReturnedIndex: cloneHistoryIndex(read.ReturnedIndex),
		Slots:         append([]output.HistorySlot(nil), read.Slots...),
	}
	switch {
	case err == nil || errors.Is(err, errHistoryEcho):
		stop.Code = historyCodeEcho
		stop.Detail = errHistoryEcho.Error()
		if len(read.Slots) > 0 {
			value := read.Slots[0].Value
			stop.DecodedSlot = &value
		}
	case read.DecodeError != "":
		stop.Code = historyCodeDecode
		stop.Detail = read.DecodeError
	default:
		stop.Code = historyCodeTransport
		if err != nil {
			stop.Detail = err.Error()
		}
	}
	return stop
}

func cloneHistoryIndex(index *uint8) *uint8 {
	if index == nil {
		return nil
	}
	value := *index
	return &value
}

func historyOwnerRecord(reads []output.HistoryRead, list *output.HistoryList, flags []output.HistoryFlagReading, categories []output.HistoryCategoryReading, branches []output.HistoryBranchReading) output.HistoryOwnerRecord {
	return output.HistoryOwnerRecord{
		Reads:      reads,
		List:       list,
		Flags:      flags,
		Categories: categories,
		Branches:   branches,
		Ambiguity:  historyAmbiguityUnknown,
	}
}

func historyField(key catalogvaillant.B503HistoryFlagKey) (catalogvaillant.B503HistoryField, bool, error) {
	return catalogvaillant.B503HistoryLookup(key)
}

func qualifyHistoryField(key catalogvaillant.B503HistoryFlagKey, raw byte) (output.HistoryFlagReading, bool, error) {
	field, found, err := historyField(key)
	reading := output.HistoryFlagReading{
		Field:         key.B503HistoryFlagField,
		Raw:           []byte{raw},
		Qualification: historyQualificationLevel,
	}
	if err != nil || !found {
		reading.Level = historyQualificationLevel
		reading.Unresolved = true
		return reading, found, err
	}
	reading.Level = field.B503HistoryFlagGrade
	if label, ok := field.B503HistoryFlagEnumeration[int(raw)]; ok {
		reading.Label = label
	} else if field.B503HistoryFlagOtherwise != nil {
		reading.Label = *field.B503HistoryFlagOtherwise
	} else if len(field.B503HistoryFlagEnumeration) > 0 {
		reading.Unresolved = true
	}
	if len(field.B503HistoryFlagBits) > 0 || len(field.B503HistoryFlagReserved) > 0 {
		reserved := map[int]bool{}
		for _, bit := range field.B503HistoryFlagReserved {
			reserved[bit] = true
		}
		for bit := 0; bit < 8; bit++ {
			name := field.B503HistoryFlagBits[bit]
			reading.Bits = append(reading.Bits, output.HistoryFlagBit{
				Index:    bit,
				Name:     name,
				Set:      raw&(1<<uint(bit)) != 0,
				Reserved: reserved[bit],
			})
		}
	}
	return reading, true, nil
}

func qualifyHistoryCategory(key catalogvaillant.B503TextCategoryKey, raw byte) (output.HistoryCategoryReading, bool, error) {
	entry, found, err := catalogvaillant.B503TextCategoryLookup(key)
	reading := output.HistoryCategoryReading{
		Variant:   key.B503TextCategoryVariant,
		Raw:       raw,
		Ambiguity: historyAmbiguityUnknown,
	}
	if err != nil || !found {
		reading.Level = historyQualificationLevel
		reading.Unresolved = true
		return reading, found, err
	}
	reading.Level = entry.B503TextCategoryGrade
	label, named := entry.B503TextCategoryLabels[raw]
	reading.Named = named
	if named && label != nil {
		reading.Label = *label
		reading.Labelled = true
	}
	if !reading.Labelled {
		reading.Unresolved = true
	}
	return reading, true, nil
}

func qualifyHistoryBranch(fields catalogvaillant.B503HistoryBranchFields) output.HistoryBranchReading {
	view := catalogvaillant.B503HistoryBranchHistory(fields)
	reading := output.HistoryBranchReading{
		Marker:   fields.B503HistoryBranchMarker,
		State:    view.B503HistoryBranchState,
		Meaning:  view.B503HistoryBranchMeaning,
		Complete: view.B503HistoryBranchComplete,
	}
	for _, component := range view.B503HistoryBranchComponents {
		reading.Components = append(reading.Components, output.HistoryBranchComponent{State: component.B503HistoryBranchState})
	}
	return reading
}

func retainHistory(p *vaillantSemanticPoller, record output.HistoryOwnerRecord) {
	if p == nil {
		return
	}
	historyMu.Lock()
	defer historyMu.Unlock()
	cloned := record.GraphQL()
	historyHeld[p] = &cloned
}

func retainedHistory(p *vaillantSemanticPoller) *output.HistoryOwnerRecord {
	if p == nil {
		return nil
	}
	historyMu.Lock()
	defer historyMu.Unlock()
	record := historyHeld[p]
	if record == nil {
		return nil
	}
	cloned := record.GraphQL()
	return &cloned
}

var (
	historyMu   sync.Mutex
	historyHeld = map[*vaillantSemanticPoller]*output.HistoryOwnerRecord{}
)

func bindHistorySurfaces(status *output.SystemStatus, record output.HistoryOwnerRecord) *output.MCPSystemStatus {
	if status == nil {
		return nil
	}
	status.SetHistory(&record)
	return output.MCPSystemFromGraphQL(status)
}
