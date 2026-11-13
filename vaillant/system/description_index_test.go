package system

import (
	descriptionIndexBytes "bytes"
	descriptionIndexJSON "encoding/json"
	descriptionIndexErrors "errors"
	descriptionIndexFmt "fmt"
	descriptionIndexOS "os"
	descriptionIndexReflect "reflect"
	descriptionIndexTesting "testing"
)

func descriptionIndexBaseRecord() DescriptionRecord {
	return DescriptionRecord{DescriptionIndexKey: DescriptionKey{"alpha", "ab25_software", "ab25_hardware", 1, 2, 2, 3, 4, 5}, DescriptionIndexState: DescriptionIndexSelectorQualified, RawDescription: []byte{17, 34}}
}

func TestAb25ExactKeyAxes(descriptionIndexT *descriptionIndexTesting.T) {
	descriptionIndexBase := descriptionIndexBaseRecord()
	descriptionIndexAxes := []struct {
		descriptionIndexName   string
		descriptionIndexChange func(*DescriptionKey)
	}{
		{"ab25_family", func(descriptionIndexKey *DescriptionKey) { descriptionIndexKey.DescriptionIndexFamily = "beta" }},
		{"ab25_software", func(descriptionIndexKey *DescriptionKey) {
			descriptionIndexKey.DescriptionIndexSoftwareIdentity = "ab25_other_software"
		}},
		{"ab25_hardware", func(descriptionIndexKey *DescriptionKey) {
			descriptionIndexKey.DescriptionIndexHardwareIdentity = "ab25_other_hardware"
		}},
		{"ab25_api_version", func(descriptionIndexKey *DescriptionKey) { descriptionIndexKey.DescriptionIndexAPIVersion = 2 }},
		{"ab25_api_revision", func(descriptionIndexKey *DescriptionKey) { descriptionIndexKey.DescriptionIndexAPIRevision = 1 }},
		{"ab25_opcode", func(descriptionIndexKey *DescriptionKey) { descriptionIndexKey.DescriptionIndexOpcode = 6 }},
		{"ab25_group", func(descriptionIndexKey *DescriptionKey) { descriptionIndexKey.DescriptionIndexGroup = 4 }},
		{"ab25_instance", func(descriptionIndexKey *DescriptionKey) { descriptionIndexKey.DescriptionIndexRequestedInstance = 5 }},
		{"ab25_register", func(descriptionIndexKey *DescriptionKey) { descriptionIndexKey.DescriptionIndexRegister = 261 }},
	}
	for _, descriptionIndexAxis := range descriptionIndexAxes {
		descriptionIndexT.Run(descriptionIndexAxis.descriptionIndexName, func(descriptionIndexT *descriptionIndexTesting.T) {
			descriptionIndexOther := descriptionIndexBaseRecord()
			descriptionIndexAxis.descriptionIndexChange(&descriptionIndexOther.DescriptionIndexKey)
			descriptionIndexIndex, descriptionIndexErr := NewDescriptionIndex([]DescriptionRecord{descriptionIndexBase})
			if descriptionIndexErr != nil {
				descriptionIndexT.Fatal(descriptionIndexErr)
			}
			if _, descriptionIndexOK, descriptionIndexErr := descriptionIndexIndex.DescriptionIndexLookup(descriptionIndexOther.DescriptionIndexKey); descriptionIndexOK || descriptionIndexErr != nil {
				descriptionIndexT.Fatal("inexact key matched")
			}
			descriptionIndexOther.RawDescription = []byte{99}
			descriptionIndexIndex, descriptionIndexErr = NewDescriptionIndex([]DescriptionRecord{descriptionIndexBase, descriptionIndexOther})
			if descriptionIndexErr != nil {
				descriptionIndexT.Fatal(descriptionIndexErr)
			}
			for _, descriptionIndexWant := range []DescriptionRecord{descriptionIndexBase, descriptionIndexOther} {
				descriptionIndexGot, descriptionIndexOK, descriptionIndexErr := descriptionIndexIndex.DescriptionIndexLookup(descriptionIndexWant.DescriptionIndexKey)
				if !descriptionIndexOK || descriptionIndexErr != nil || !descriptionIndexReflect.DeepEqual(descriptionIndexGot, descriptionIndexWant) {
					descriptionIndexT.Fatalf("exact control: %+v %v %v", descriptionIndexGot, descriptionIndexOK, descriptionIndexErr)
				}
			}
		})
	}
}

func TestAb25StatesAndUnknowns(descriptionIndexT *descriptionIndexTesting.T) {
	descriptionIndexData, descriptionIndexErr := descriptionIndexOS.ReadFile("testdata/ab25_description_states.json")
	if descriptionIndexErr != nil {
		descriptionIndexT.Fatal(descriptionIndexErr)
	}
	var descriptionIndexRows []struct {
		DescriptionIndexFamily string           `json:"family"`
		DescriptionIndexState  DescriptionState `json:"state"`
		DescriptionIndexRaw    []byte           `json:"raw"`
		DescriptionIndexCodec  *string          `json:"codec"`
		DescriptionIndexWidth  *uint16          `json:"width"`
	}
	if descriptionIndexErr = descriptionIndexJSON.Unmarshal(descriptionIndexData, &descriptionIndexRows); descriptionIndexErr != nil {
		descriptionIndexT.Fatal(descriptionIndexErr)
	}
	descriptionIndexExpected := []DescriptionState{DescriptionIndexSelectorQualified, DescriptionIndexLengthInferred, DescriptionIndexUndecoded, DescriptionIndexUnknown, DescriptionIndexUnknown}
	if len(descriptionIndexRows) != len(descriptionIndexExpected) {
		descriptionIndexT.Fatal("fixture size")
	}
	for descriptionIndexPosition, descriptionIndexRow := range descriptionIndexRows {
		descriptionIndexRecord := descriptionIndexBaseRecord()
		descriptionIndexRecord.DescriptionIndexKey.DescriptionIndexFamily = descriptionIndexRow.DescriptionIndexFamily
		descriptionIndexRecord.DescriptionIndexState = descriptionIndexRow.DescriptionIndexState
		descriptionIndexRecord.RawDescription = descriptionIndexRow.DescriptionIndexRaw
		descriptionIndexRecord.DescriptionIndexCodec = descriptionIndexRow.DescriptionIndexCodec
		descriptionIndexRecord.DescriptionIndexWidth = descriptionIndexRow.DescriptionIndexWidth
		descriptionIndexIndex, descriptionIndexErr := NewDescriptionIndex([]DescriptionRecord{descriptionIndexRecord})
		if descriptionIndexErr != nil {
			descriptionIndexT.Fatal(descriptionIndexErr)
		}
		descriptionIndexGot, descriptionIndexOK, descriptionIndexErr := descriptionIndexIndex.DescriptionIndexLookup(descriptionIndexRecord.DescriptionIndexKey)
		if !descriptionIndexOK || descriptionIndexErr != nil || descriptionIndexGot.DescriptionIndexState != descriptionIndexExpected[descriptionIndexPosition] {
			descriptionIndexT.Fatalf("row %d: %+v %v", descriptionIndexPosition, descriptionIndexGot, descriptionIndexErr)
		}
		if !descriptionIndexBytes.Equal(descriptionIndexGot.RawDescription, descriptionIndexRow.DescriptionIndexRaw) || !descriptionIndexReflect.DeepEqual(descriptionIndexGot.DescriptionIndexCodec, descriptionIndexRow.DescriptionIndexCodec) || !descriptionIndexReflect.DeepEqual(descriptionIndexGot.DescriptionIndexWidth, descriptionIndexRow.DescriptionIndexWidth) || descriptionIndexGot.DescriptionIndexPresence() != "unknown" {
			descriptionIndexT.Fatal("metadata inferred or evidence lost")
		}
	}
}

func TestAb25Population(descriptionIndexT *descriptionIndexTesting.T) {
	var descriptionIndexRecords []DescriptionRecord
	for descriptionIndexPosition := 0; descriptionIndexPosition < 2093; descriptionIndexPosition++ {
		descriptionIndexRecord := descriptionIndexBaseRecord()
		descriptionIndexRecord.DescriptionIndexKey.DescriptionIndexRegister = uint16(descriptionIndexPosition)
		switch {
		case descriptionIndexPosition < 139:
			descriptionIndexRecord.DescriptionIndexState = DescriptionIndexSelectorQualified
		case descriptionIndexPosition < 2024:
			descriptionIndexRecord.DescriptionIndexState = DescriptionIndexLengthInferred
		default:
			descriptionIndexRecord.DescriptionIndexState = DescriptionIndexUndecoded
		}
		descriptionIndexRecords = append(descriptionIndexRecords, descriptionIndexRecord)
	}
	descriptionIndexIndex, descriptionIndexErr := NewDescriptionIndex(descriptionIndexRecords)
	if descriptionIndexErr != nil {
		descriptionIndexT.Fatal(descriptionIndexErr)
	}
	descriptionIndexCounts := map[DescriptionState]int{}
	for _, descriptionIndexRecord := range descriptionIndexRecords {
		descriptionIndexGot, descriptionIndexOK, descriptionIndexErr := descriptionIndexIndex.DescriptionIndexLookup(descriptionIndexRecord.DescriptionIndexKey)
		if !descriptionIndexOK || descriptionIndexErr != nil {
			descriptionIndexT.Fatal("population record lost")
		}
		descriptionIndexCounts[descriptionIndexGot.DescriptionIndexState]++
	}
	if descriptionIndexCounts[DescriptionIndexSelectorQualified] != 139 {
		descriptionIndexT.Fatal(descriptionIndexCounts)
	}
	if descriptionIndexCounts[DescriptionIndexLengthInferred] != 1885 {
		descriptionIndexT.Fatal(descriptionIndexCounts)
	}
	if descriptionIndexCounts[DescriptionIndexUndecoded] != 69 {
		descriptionIndexT.Fatal(descriptionIndexCounts)
	}
	if len(descriptionIndexRecords) != 2093 || len(descriptionIndexCounts) != 3 {
		descriptionIndexT.Fatal("population total or states")
	}
}

func TestAb25Ownership(descriptionIndexT *descriptionIndexTesting.T) {
	descriptionIndexCodec, descriptionIndexWidth := "ab25_opaque", uint16(2)
	descriptionIndexSource := descriptionIndexBaseRecord()
	descriptionIndexSource.DescriptionIndexCodec, descriptionIndexSource.DescriptionIndexWidth = &descriptionIndexCodec, &descriptionIndexWidth
	descriptionIndexWant := descriptionIndexBaseRecord()
	descriptionIndexWantCodec, descriptionIndexWantWidth := "ab25_opaque", uint16(2)
	descriptionIndexWant.DescriptionIndexCodec, descriptionIndexWant.DescriptionIndexWidth = &descriptionIndexWantCodec, &descriptionIndexWantWidth
	descriptionIndexIndex, descriptionIndexErr := NewDescriptionIndex([]DescriptionRecord{descriptionIndexSource})
	if descriptionIndexErr != nil {
		descriptionIndexT.Fatal(descriptionIndexErr)
	}
	descriptionIndexSource.RawDescription[0] = 99
	descriptionIndexCodec, descriptionIndexWidth = "ab25_changed", 9
	for descriptionIndexPosition := 0; descriptionIndexPosition < 2; descriptionIndexPosition++ {
		descriptionIndexGot, descriptionIndexOK, descriptionIndexErr := descriptionIndexIndex.DescriptionIndexLookup(descriptionIndexWant.DescriptionIndexKey)
		if !descriptionIndexOK || descriptionIndexErr != nil || !descriptionIndexReflect.DeepEqual(descriptionIndexGot, descriptionIndexWant) {
			descriptionIndexT.Fatal("stored record mutated")
		}
		descriptionIndexGot.RawDescription[0] = 88
		*descriptionIndexGot.DescriptionIndexCodec = "ab25_changed"
		*descriptionIndexGot.DescriptionIndexWidth = 8
		descriptionIndexCopies := descriptionIndexIndex.DescriptionIndexRecords(descriptionIndexWant.DescriptionIndexKey)
		if len(descriptionIndexCopies) != 1 || !descriptionIndexReflect.DeepEqual(descriptionIndexCopies[0], descriptionIndexWant) {
			descriptionIndexT.Fatal("record accessor mutated storage")
		}
		descriptionIndexCopies[0].RawDescription[0] = 77
		*descriptionIndexCopies[0].DescriptionIndexCodec = "ab25_changed"
		*descriptionIndexCopies[0].DescriptionIndexWidth = 7
	}
}

func TestAb25Conflicts(descriptionIndexT *descriptionIndexTesting.T) {
	for _, descriptionIndexAxis := range []string{"ab25_raw", "ab25_state", "ab25_codec", "ab25_width"} {
		descriptionIndexT.Run(descriptionIndexAxis, func(descriptionIndexT *descriptionIndexTesting.T) {
			descriptionIndexBase := descriptionIndexBaseRecord()
			descriptionIndexOther := descriptionIndexBaseRecord()
			switch descriptionIndexAxis {
			case "ab25_raw":
				descriptionIndexOther.RawDescription = []byte{9}
			case "ab25_state":
				descriptionIndexOther.DescriptionIndexState = DescriptionIndexUndecoded
			case "ab25_codec":
				descriptionIndexCodec := "ab25_opaque"
				descriptionIndexOther.DescriptionIndexCodec = &descriptionIndexCodec
			case "ab25_width":
				descriptionIndexWidth := uint16(2)
				descriptionIndexOther.DescriptionIndexWidth = &descriptionIndexWidth
			}
			descriptionIndexGood := descriptionIndexBaseRecord()
			descriptionIndexGood.DescriptionIndexKey.DescriptionIndexFamily = "beta"
			descriptionIndexIndex, descriptionIndexErr := NewDescriptionIndex([]DescriptionRecord{descriptionIndexBase, descriptionIndexOther, descriptionIndexBase, descriptionIndexGood})
			if !descriptionIndexErrors.Is(descriptionIndexErr, ErrDescriptionConflict) {
				descriptionIndexT.Fatal("missing construction conflict")
			}
			descriptionIndexGot, descriptionIndexOK, descriptionIndexErr := descriptionIndexIndex.DescriptionIndexLookup(descriptionIndexBase.DescriptionIndexKey)
			if descriptionIndexOK || !descriptionIndexErrors.Is(descriptionIndexErr, ErrDescriptionConflict) || !descriptionIndexReflect.DeepEqual(descriptionIndexGot, DescriptionRecord{}) {
				descriptionIndexT.Fatal("conflict returned usable record")
			}
			descriptionIndexRaw := descriptionIndexIndex.DescriptionIndexRecords(descriptionIndexBase.DescriptionIndexKey)
			if !descriptionIndexReflect.DeepEqual(descriptionIndexRaw, []DescriptionRecord{descriptionIndexBase, descriptionIndexOther, descriptionIndexBase}) {
				descriptionIndexT.Fatal("conflicting evidence lost")
			}
			descriptionIndexRaw[1].RawDescription[0] = 100
			if !descriptionIndexReflect.DeepEqual(descriptionIndexIndex.DescriptionIndexRecords(descriptionIndexBase.DescriptionIndexKey)[1], descriptionIndexOther) {
				descriptionIndexT.Fatal("conflict evidence aliased")
			}
			descriptionIndexGot, descriptionIndexOK, descriptionIndexErr = descriptionIndexIndex.DescriptionIndexLookup(descriptionIndexGood.DescriptionIndexKey)
			if !descriptionIndexOK || descriptionIndexErr != nil || !descriptionIndexReflect.DeepEqual(descriptionIndexGot, descriptionIndexGood) {
				descriptionIndexT.Fatal("unaffected record lost")
			}
		})
	}
	descriptionIndexBase := descriptionIndexBaseRecord()
	descriptionIndexIndex, descriptionIndexErr := NewDescriptionIndex([]DescriptionRecord{descriptionIndexBase, descriptionIndexBase})
	if descriptionIndexErr != nil {
		descriptionIndexT.Fatal(descriptionIndexErr)
	}
	if _, descriptionIndexOK, descriptionIndexErr := descriptionIndexIndex.DescriptionIndexLookup(descriptionIndexBase.DescriptionIndexKey); !descriptionIndexOK || descriptionIndexErr != nil {
		descriptionIndexT.Fatal("identical duplicate rejected")
	}
}

func TestAb25Empty(descriptionIndexT *descriptionIndexTesting.T) {
	for descriptionIndexPosition, descriptionIndexIndex := range []*DescriptionIndex{nil, {}} {
		descriptionIndexT.Run(descriptionIndexFmt.Sprintf("ab25_%d", descriptionIndexPosition), func(descriptionIndexT *descriptionIndexTesting.T) {
			if _, descriptionIndexOK, descriptionIndexErr := descriptionIndexIndex.DescriptionIndexLookup(descriptionIndexBaseRecord().DescriptionIndexKey); descriptionIndexOK || descriptionIndexErr != nil {
				descriptionIndexT.Fatal("empty index matched")
			}
			if len(descriptionIndexIndex.DescriptionIndexRecords(descriptionIndexBaseRecord().DescriptionIndexKey)) != 0 {
				descriptionIndexT.Fatal("empty index records")
			}
		})
	}
}
