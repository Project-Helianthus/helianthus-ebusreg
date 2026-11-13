package system

import (
	descriptionIndexErrors "errors"
	descriptionIndexReflect "reflect"
)

type DescriptionState string

const (
	DescriptionIndexUnknown           DescriptionState = "unknown"
	DescriptionIndexSelectorQualified DescriptionState = "selector_qualified"
	DescriptionIndexLengthInferred    DescriptionState = "length_inferred"
	DescriptionIndexUndecoded         DescriptionState = "undecoded"
)

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

type DescriptionRecord struct {
	DescriptionIndexKey   DescriptionKey   `json:"key"`
	DescriptionIndexState DescriptionState `json:"state"`
	RawDescription        []byte           `json:"raw_description"`
	DescriptionIndexCodec *string          `json:"codec"`
	DescriptionIndexWidth *uint16          `json:"width"`
}

func (descriptionIndexRecord DescriptionRecord) DescriptionIndexPresence() string { return "unknown" }

var ErrDescriptionConflict = descriptionIndexErrors.New("description index: conflicting description records")

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
