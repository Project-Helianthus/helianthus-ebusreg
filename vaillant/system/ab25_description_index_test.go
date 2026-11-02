package system

import (
	ab25Bytes "bytes"
	ab25JSON "encoding/json"
	ab25Errors "errors"
	ab25Fmt "fmt"
	ab25OS "os"
	ab25Reflect "reflect"
	ab25Testing "testing"
)

func ab25BaseRecord() Ab25DescriptionRecord {
	return Ab25DescriptionRecord{Ab25Key: Ab25DescriptionKey{"alpha", "ab25_software", "ab25_hardware", 1, 2, 2, 3, 4, 5}, Ab25State: Ab25SelectorQualified, Ab25RawDescription: []byte{17, 34}}
}

func TestAb25ExactKeyAxes(ab25T *ab25Testing.T) {
	ab25Base := ab25BaseRecord()
	ab25Axes := []struct {
		ab25Name   string
		ab25Change func(*Ab25DescriptionKey)
	}{
		{"ab25_family", func(ab25Key *Ab25DescriptionKey) { ab25Key.Ab25Family = "beta" }},
		{"ab25_software", func(ab25Key *Ab25DescriptionKey) { ab25Key.Ab25SoftwareIdentity = "ab25_other_software" }},
		{"ab25_hardware", func(ab25Key *Ab25DescriptionKey) { ab25Key.Ab25HardwareIdentity = "ab25_other_hardware" }},
		{"ab25_api_version", func(ab25Key *Ab25DescriptionKey) { ab25Key.Ab25APIVersion = 2 }},
		{"ab25_api_revision", func(ab25Key *Ab25DescriptionKey) { ab25Key.Ab25APIRevision = 1 }},
		{"ab25_opcode", func(ab25Key *Ab25DescriptionKey) { ab25Key.Ab25Opcode = 6 }},
		{"ab25_group", func(ab25Key *Ab25DescriptionKey) { ab25Key.Ab25Group = 4 }},
		{"ab25_instance", func(ab25Key *Ab25DescriptionKey) { ab25Key.Ab25RequestedInstance = 5 }},
		{"ab25_register", func(ab25Key *Ab25DescriptionKey) { ab25Key.Ab25Register = 261 }},
	}
	for _, ab25Axis := range ab25Axes {
		ab25T.Run(ab25Axis.ab25Name, func(ab25T *ab25Testing.T) {
			ab25Other := ab25BaseRecord()
			ab25Axis.ab25Change(&ab25Other.Ab25Key)
			ab25Index, ab25Err := Ab25NewDescriptionIndex([]Ab25DescriptionRecord{ab25Base})
			if ab25Err != nil {
				ab25T.Fatal(ab25Err)
			}
			if _, ab25OK, ab25Err := ab25Index.Ab25Lookup(ab25Other.Ab25Key); ab25OK || ab25Err != nil {
				ab25T.Fatal("inexact key matched")
			}
			ab25Other.Ab25RawDescription = []byte{99}
			ab25Index, ab25Err = Ab25NewDescriptionIndex([]Ab25DescriptionRecord{ab25Base, ab25Other})
			if ab25Err != nil {
				ab25T.Fatal(ab25Err)
			}
			for _, ab25Want := range []Ab25DescriptionRecord{ab25Base, ab25Other} {
				ab25Got, ab25OK, ab25Err := ab25Index.Ab25Lookup(ab25Want.Ab25Key)
				if !ab25OK || ab25Err != nil || !ab25Reflect.DeepEqual(ab25Got, ab25Want) {
					ab25T.Fatalf("exact control: %+v %v %v", ab25Got, ab25OK, ab25Err)
				}
			}
		})
	}
}

func TestAb25StatesAndUnknowns(ab25T *ab25Testing.T) {
	ab25Data, ab25Err := ab25OS.ReadFile("testdata/ab25_description_states.json")
	if ab25Err != nil {
		ab25T.Fatal(ab25Err)
	}
	var ab25Rows []struct {
		Ab25Family string               `json:"family"`
		Ab25State  Ab25DescriptionState `json:"state"`
		Ab25Raw    []byte               `json:"raw"`
		Ab25Codec  *string              `json:"codec"`
		Ab25Width  *uint16              `json:"width"`
	}
	if ab25Err = ab25JSON.Unmarshal(ab25Data, &ab25Rows); ab25Err != nil {
		ab25T.Fatal(ab25Err)
	}
	ab25Expected := []Ab25DescriptionState{Ab25SelectorQualified, Ab25LengthInferred, Ab25Undecoded, Ab25Unknown, Ab25Unknown}
	if len(ab25Rows) != len(ab25Expected) {
		ab25T.Fatal("fixture size")
	}
	for ab25Position, ab25Row := range ab25Rows {
		ab25Record := ab25BaseRecord()
		ab25Record.Ab25Key.Ab25Family = ab25Row.Ab25Family
		ab25Record.Ab25State = ab25Row.Ab25State
		ab25Record.Ab25RawDescription = ab25Row.Ab25Raw
		ab25Record.Ab25Codec = ab25Row.Ab25Codec
		ab25Record.Ab25Width = ab25Row.Ab25Width
		ab25Index, ab25Err := Ab25NewDescriptionIndex([]Ab25DescriptionRecord{ab25Record})
		if ab25Err != nil {
			ab25T.Fatal(ab25Err)
		}
		ab25Got, ab25OK, ab25Err := ab25Index.Ab25Lookup(ab25Record.Ab25Key)
		if !ab25OK || ab25Err != nil || ab25Got.Ab25State != ab25Expected[ab25Position] {
			ab25T.Fatalf("row %d: %+v %v", ab25Position, ab25Got, ab25Err)
		}
		if !ab25Bytes.Equal(ab25Got.Ab25RawDescription, ab25Row.Ab25Raw) || !ab25Reflect.DeepEqual(ab25Got.Ab25Codec, ab25Row.Ab25Codec) || !ab25Reflect.DeepEqual(ab25Got.Ab25Width, ab25Row.Ab25Width) || ab25Got.Ab25Presence() != "unknown" {
			ab25T.Fatal("metadata inferred or evidence lost")
		}
	}
}

func TestAb25Population(ab25T *ab25Testing.T) {
	var ab25Records []Ab25DescriptionRecord
	for ab25Position := 0; ab25Position < 2093; ab25Position++ {
		ab25Record := ab25BaseRecord()
		ab25Record.Ab25Key.Ab25Register = uint16(ab25Position)
		switch {
		case ab25Position < 139:
			ab25Record.Ab25State = Ab25SelectorQualified
		case ab25Position < 2024:
			ab25Record.Ab25State = Ab25LengthInferred
		default:
			ab25Record.Ab25State = Ab25Undecoded
		}
		ab25Records = append(ab25Records, ab25Record)
	}
	ab25Index, ab25Err := Ab25NewDescriptionIndex(ab25Records)
	if ab25Err != nil {
		ab25T.Fatal(ab25Err)
	}
	ab25Counts := map[Ab25DescriptionState]int{}
	for _, ab25Record := range ab25Records {
		ab25Got, ab25OK, ab25Err := ab25Index.Ab25Lookup(ab25Record.Ab25Key)
		if !ab25OK || ab25Err != nil {
			ab25T.Fatal("population record lost")
		}
		ab25Counts[ab25Got.Ab25State]++
	}
	if ab25Counts[Ab25SelectorQualified] != 139 {
		ab25T.Fatal(ab25Counts)
	}
	if ab25Counts[Ab25LengthInferred] != 1885 {
		ab25T.Fatal(ab25Counts)
	}
	if ab25Counts[Ab25Undecoded] != 69 {
		ab25T.Fatal(ab25Counts)
	}
	if len(ab25Records) != 2093 || len(ab25Counts) != 3 {
		ab25T.Fatal("population total or states")
	}
}

func TestAb25Ownership(ab25T *ab25Testing.T) {
	ab25Codec, ab25Width := "ab25_opaque", uint16(2)
	ab25Source := ab25BaseRecord()
	ab25Source.Ab25Codec, ab25Source.Ab25Width = &ab25Codec, &ab25Width
	ab25Want := ab25BaseRecord()
	ab25WantCodec, ab25WantWidth := "ab25_opaque", uint16(2)
	ab25Want.Ab25Codec, ab25Want.Ab25Width = &ab25WantCodec, &ab25WantWidth
	ab25Index, ab25Err := Ab25NewDescriptionIndex([]Ab25DescriptionRecord{ab25Source})
	if ab25Err != nil {
		ab25T.Fatal(ab25Err)
	}
	ab25Source.Ab25RawDescription[0] = 99
	ab25Codec, ab25Width = "ab25_changed", 9
	for ab25Position := 0; ab25Position < 2; ab25Position++ {
		ab25Got, ab25OK, ab25Err := ab25Index.Ab25Lookup(ab25Want.Ab25Key)
		if !ab25OK || ab25Err != nil || !ab25Reflect.DeepEqual(ab25Got, ab25Want) {
			ab25T.Fatal("stored record mutated")
		}
		ab25Got.Ab25RawDescription[0] = 88
		*ab25Got.Ab25Codec = "ab25_changed"
		*ab25Got.Ab25Width = 8
		ab25Copies := ab25Index.Ab25Records(ab25Want.Ab25Key)
		if len(ab25Copies) != 1 || !ab25Reflect.DeepEqual(ab25Copies[0], ab25Want) {
			ab25T.Fatal("record accessor mutated storage")
		}
		ab25Copies[0].Ab25RawDescription[0] = 77
		*ab25Copies[0].Ab25Codec = "ab25_changed"
		*ab25Copies[0].Ab25Width = 7
	}
}

func TestAb25Conflicts(ab25T *ab25Testing.T) {
	for _, ab25Axis := range []string{"ab25_raw", "ab25_state", "ab25_codec", "ab25_width"} {
		ab25T.Run(ab25Axis, func(ab25T *ab25Testing.T) {
			ab25Base := ab25BaseRecord()
			ab25Other := ab25BaseRecord()
			switch ab25Axis {
			case "ab25_raw":
				ab25Other.Ab25RawDescription = []byte{9}
			case "ab25_state":
				ab25Other.Ab25State = Ab25Undecoded
			case "ab25_codec":
				ab25Codec := "ab25_opaque"
				ab25Other.Ab25Codec = &ab25Codec
			case "ab25_width":
				ab25Width := uint16(2)
				ab25Other.Ab25Width = &ab25Width
			}
			ab25Good := ab25BaseRecord()
			ab25Good.Ab25Key.Ab25Family = "beta"
			ab25Index, ab25Err := Ab25NewDescriptionIndex([]Ab25DescriptionRecord{ab25Base, ab25Other, ab25Base, ab25Good})
			if !ab25Errors.Is(ab25Err, ErrAb25DescriptionConflict) {
				ab25T.Fatal("missing construction conflict")
			}
			ab25Got, ab25OK, ab25Err := ab25Index.Ab25Lookup(ab25Base.Ab25Key)
			if ab25OK || !ab25Errors.Is(ab25Err, ErrAb25DescriptionConflict) || !ab25Reflect.DeepEqual(ab25Got, Ab25DescriptionRecord{}) {
				ab25T.Fatal("conflict returned usable record")
			}
			ab25Raw := ab25Index.Ab25Records(ab25Base.Ab25Key)
			if !ab25Reflect.DeepEqual(ab25Raw, []Ab25DescriptionRecord{ab25Base, ab25Other, ab25Base}) {
				ab25T.Fatal("conflicting evidence lost")
			}
			ab25Raw[1].Ab25RawDescription[0] = 100
			if !ab25Reflect.DeepEqual(ab25Index.Ab25Records(ab25Base.Ab25Key)[1], ab25Other) {
				ab25T.Fatal("conflict evidence aliased")
			}
			ab25Got, ab25OK, ab25Err = ab25Index.Ab25Lookup(ab25Good.Ab25Key)
			if !ab25OK || ab25Err != nil || !ab25Reflect.DeepEqual(ab25Got, ab25Good) {
				ab25T.Fatal("unaffected record lost")
			}
		})
	}
	ab25Base := ab25BaseRecord()
	ab25Index, ab25Err := Ab25NewDescriptionIndex([]Ab25DescriptionRecord{ab25Base, ab25Base})
	if ab25Err != nil {
		ab25T.Fatal(ab25Err)
	}
	if _, ab25OK, ab25Err := ab25Index.Ab25Lookup(ab25Base.Ab25Key); !ab25OK || ab25Err != nil {
		ab25T.Fatal("identical duplicate rejected")
	}
}

func TestAb25Empty(ab25T *ab25Testing.T) {
	for ab25Position, ab25Index := range []*Ab25DescriptionIndex{nil, {}} {
		ab25T.Run(ab25Fmt.Sprintf("ab25_%d", ab25Position), func(ab25T *ab25Testing.T) {
			if _, ab25OK, ab25Err := ab25Index.Ab25Lookup(ab25BaseRecord().Ab25Key); ab25OK || ab25Err != nil {
				ab25T.Fatal("empty index matched")
			}
			if len(ab25Index.Ab25Records(ab25BaseRecord().Ab25Key)) != 0 {
				ab25T.Fatal("empty index records")
			}
		})
	}
}
