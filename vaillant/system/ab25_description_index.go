package system

import (
	ab25Errors "errors"
	ab25Reflect "reflect"
)

type Ab25DescriptionState string

const (
	Ab25Unknown           Ab25DescriptionState = "unknown"
	Ab25SelectorQualified Ab25DescriptionState = "selector_qualified"
	Ab25LengthInferred    Ab25DescriptionState = "length_inferred"
	Ab25Undecoded         Ab25DescriptionState = "undecoded"
)

type Ab25DescriptionKey struct {
	Ab25Family            string
	Ab25SoftwareIdentity  string
	Ab25HardwareIdentity  string
	Ab25APIVersion        uint16
	Ab25APIRevision       uint16
	Ab25Opcode            byte
	Ab25Group             byte
	Ab25RequestedInstance byte
	Ab25Register          uint16
}

type Ab25DescriptionRecord struct {
	Ab25Key            Ab25DescriptionKey
	Ab25State          Ab25DescriptionState
	Ab25RawDescription []byte
	Ab25Codec          *string
	Ab25Width          *uint16
}

func (ab25Record Ab25DescriptionRecord) Ab25Presence() string { return "unknown" }

var ErrAb25DescriptionConflict = ab25Errors.New("ab25 conflicting description records")

type Ab25DescriptionIndex struct {
	ab25Records   map[Ab25DescriptionKey][]Ab25DescriptionRecord
	ab25Conflicts map[Ab25DescriptionKey]bool
}

func ab25CopyDescription(ab25Record Ab25DescriptionRecord) Ab25DescriptionRecord {
	if ab25Record.Ab25RawDescription != nil {
		ab25Record.Ab25RawDescription = append([]byte{}, ab25Record.Ab25RawDescription...)
	}
	if ab25Record.Ab25Codec != nil {
		ab25Codec := *ab25Record.Ab25Codec
		ab25Record.Ab25Codec = &ab25Codec
	}
	if ab25Record.Ab25Width != nil {
		ab25Width := *ab25Record.Ab25Width
		ab25Record.Ab25Width = &ab25Width
	}
	return ab25Record
}

func Ab25NewDescriptionIndex(ab25Records []Ab25DescriptionRecord) (*Ab25DescriptionIndex, error) {
	ab25Index := &Ab25DescriptionIndex{
		ab25Records:   make(map[Ab25DescriptionKey][]Ab25DescriptionRecord),
		ab25Conflicts: make(map[Ab25DescriptionKey]bool),
	}
	var ab25Err error
	for _, ab25Record := range ab25Records {
		switch ab25Record.Ab25State {
		case Ab25SelectorQualified, Ab25LengthInferred, Ab25Undecoded:
		default:
			ab25Record.Ab25State = Ab25Unknown
		}
		ab25Record = ab25CopyDescription(ab25Record)
		ab25Prior := ab25Index.ab25Records[ab25Record.Ab25Key]
		if len(ab25Prior) > 0 && !ab25Reflect.DeepEqual(ab25Prior[0], ab25Record) {
			ab25Index.ab25Conflicts[ab25Record.Ab25Key] = true
			ab25Err = ErrAb25DescriptionConflict
		}
		ab25Index.ab25Records[ab25Record.Ab25Key] = append(ab25Prior, ab25Record)
	}
	return ab25Index, ab25Err
}

func (ab25Index *Ab25DescriptionIndex) Ab25Lookup(ab25Key Ab25DescriptionKey) (Ab25DescriptionRecord, bool, error) {
	if ab25Index == nil {
		return Ab25DescriptionRecord{}, false, nil
	}
	if ab25Index.ab25Conflicts[ab25Key] {
		return Ab25DescriptionRecord{}, false, ErrAb25DescriptionConflict
	}
	ab25Records := ab25Index.ab25Records[ab25Key]
	if len(ab25Records) == 0 {
		return Ab25DescriptionRecord{}, false, nil
	}
	return ab25CopyDescription(ab25Records[0]), true, nil
}

func (ab25Index *Ab25DescriptionIndex) Ab25Records(ab25Key Ab25DescriptionKey) []Ab25DescriptionRecord {
	if ab25Index == nil {
		return nil
	}
	ab25Records := ab25Index.ab25Records[ab25Key]
	ab25Copies := make([]Ab25DescriptionRecord, len(ab25Records))
	for ab25Position, ab25Record := range ab25Records {
		ab25Copies[ab25Position] = ab25CopyDescription(ab25Record)
	}
	return ab25Copies
}
