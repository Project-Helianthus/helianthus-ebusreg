package system

import (
	descriptionIndexErrors "errors"
	descriptionIndexReflect "reflect"
)

// DescriptionState says how far a retained description record was decoded.
type DescriptionState string

// Description states. NewDescriptionIndex stores any state outside this
// set as unknown.
const (
	DescriptionIndexUnknown           DescriptionState = "unknown"
	DescriptionIndexSelectorQualified DescriptionState = "selector_qualified"
	DescriptionIndexLengthInferred    DescriptionState = "length_inferred"
	DescriptionIndexUndecoded         DescriptionState = "undecoded"
)

// DescriptionKey identifies a description record exactly: the controller
// identity (family, software and hardware identity, API version and
// revision) and the opcode, group, requested instance and register. There
// is no wildcard or fallback matching.
type DescriptionKey struct {
	DescriptionIndexFamily            string `json:"family"`
	DescriptionIndexSoftwareIdentity  string `json:"software_identity"`
	DescriptionIndexHardwareIdentity  string `json:"hardware_identity"`
	DescriptionIndexAPIVersion        uint16 `json:"api_version"`
	DescriptionIndexAPIRevision       uint16 `json:"api_revision"`
	DescriptionIndexOpcode            byte   `json:"opcode"`
	DescriptionIndexGroup             byte   `json:"group"`
	DescriptionIndexRequestedInstance byte   `json:"requested_instance"`
	DescriptionIndexRegister          uint16 `json:"register"`
}

// DescriptionRecord is one retained description: its key, decoding state,
// raw description bytes and, when known, the codec name and width.
type DescriptionRecord struct {
	DescriptionIndexKey   DescriptionKey   `json:"key"`
	DescriptionIndexState DescriptionState `json:"state"`
	RawDescription        []byte           `json:"raw_description"`
	DescriptionIndexCodec *string          `json:"codec"`
	DescriptionIndexWidth *uint16          `json:"width"`
}

// DescriptionIndexPresence always returns "unknown": a description record
// says nothing about whether the register is present on a device.
func (descriptionIndexRecord DescriptionRecord) DescriptionIndexPresence() string { return "unknown" }

// ErrDescriptionConflict reports that one key has records that differ.
var ErrDescriptionConflict = descriptionIndexErrors.New("description index: conflicting description records")

// DescriptionIndex is an immutable index of description records by exact
// key. It holds its own copies of the records and hands out copies, so it is
// safe for concurrent reads; a nil index behaves as empty.
type DescriptionIndex struct {
	descriptionIndexRecords   map[DescriptionKey][]DescriptionRecord
	descriptionIndexConflicts map[DescriptionKey]bool
}

func descriptionIndexCopyDescription(descriptionIndexRecord DescriptionRecord) DescriptionRecord {
	if descriptionIndexRecord.RawDescription != nil {
		descriptionIndexRecord.RawDescription = append([]byte{}, descriptionIndexRecord.RawDescription...)
	}
	if descriptionIndexRecord.DescriptionIndexCodec != nil {
		descriptionIndexCodec := *descriptionIndexRecord.DescriptionIndexCodec
		descriptionIndexRecord.DescriptionIndexCodec = &descriptionIndexCodec
	}
	if descriptionIndexRecord.DescriptionIndexWidth != nil {
		descriptionIndexWidth := *descriptionIndexRecord.DescriptionIndexWidth
		descriptionIndexRecord.DescriptionIndexWidth = &descriptionIndexWidth
	}
	return descriptionIndexRecord
}

// NewDescriptionIndex builds an index from copies of the records. When two
// records for one key differ it still returns the index, with that key
// marked as conflicting, together with ErrDescriptionConflict; identical
// duplicates are kept without error.
func NewDescriptionIndex(descriptionIndexRecords []DescriptionRecord) (*DescriptionIndex, error) {
	descriptionIndexIndex := &DescriptionIndex{
		descriptionIndexRecords:   make(map[DescriptionKey][]DescriptionRecord),
		descriptionIndexConflicts: make(map[DescriptionKey]bool),
	}
	var descriptionIndexErr error
	for _, descriptionIndexRecord := range descriptionIndexRecords {
		switch descriptionIndexRecord.DescriptionIndexState {
		case DescriptionIndexSelectorQualified, DescriptionIndexLengthInferred, DescriptionIndexUndecoded:
		default:
			descriptionIndexRecord.DescriptionIndexState = DescriptionIndexUnknown
		}
		descriptionIndexRecord = descriptionIndexCopyDescription(descriptionIndexRecord)
		descriptionIndexPrior := descriptionIndexIndex.descriptionIndexRecords[descriptionIndexRecord.DescriptionIndexKey]
		if len(descriptionIndexPrior) > 0 && !descriptionIndexReflect.DeepEqual(descriptionIndexPrior[0], descriptionIndexRecord) {
			descriptionIndexIndex.descriptionIndexConflicts[descriptionIndexRecord.DescriptionIndexKey] = true
			descriptionIndexErr = ErrDescriptionConflict
		}
		descriptionIndexIndex.descriptionIndexRecords[descriptionIndexRecord.DescriptionIndexKey] = append(descriptionIndexPrior, descriptionIndexRecord)
	}
	return descriptionIndexIndex, descriptionIndexErr
}

// DescriptionIndexLookup returns a copy of the record for a key and true.
// It returns false for a missing key and ErrDescriptionConflict for a
// conflicting key, never picking one of the conflicting records.
func (descriptionIndexIndex *DescriptionIndex) DescriptionIndexLookup(descriptionIndexKey DescriptionKey) (DescriptionRecord, bool, error) {
	if descriptionIndexIndex == nil {
		return DescriptionRecord{}, false, nil
	}
	if descriptionIndexIndex.descriptionIndexConflicts[descriptionIndexKey] {
		return DescriptionRecord{}, false, ErrDescriptionConflict
	}
	descriptionIndexRecords := descriptionIndexIndex.descriptionIndexRecords[descriptionIndexKey]
	if len(descriptionIndexRecords) == 0 {
		return DescriptionRecord{}, false, nil
	}
	return descriptionIndexCopyDescription(descriptionIndexRecords[0]), true, nil
}

// DescriptionIndexRecords returns copies of every record for a key in input
// order, including conflicting ones.
func (descriptionIndexIndex *DescriptionIndex) DescriptionIndexRecords(descriptionIndexKey DescriptionKey) []DescriptionRecord {
	if descriptionIndexIndex == nil {
		return nil
	}
	descriptionIndexRecords := descriptionIndexIndex.descriptionIndexRecords[descriptionIndexKey]
	descriptionIndexCopies := make([]DescriptionRecord, len(descriptionIndexRecords))
	for descriptionIndexPosition, descriptionIndexRecord := range descriptionIndexRecords {
		descriptionIndexCopies[descriptionIndexPosition] = descriptionIndexCopyDescription(descriptionIndexRecord)
	}
	return descriptionIndexCopies
}
