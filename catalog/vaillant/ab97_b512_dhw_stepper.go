package vaillant

import (
	ab97bytes "bytes"
	_ "embed"
	ab97json "encoding/json"
	ab97fmt "fmt"
	ab97io "io"
	ab97maps "maps"
	ab97reflect "reflect"
	ab97slices "slices"
	ab97sync "sync"
)

//go:embed data/ab97_b512_dhw_stepper.json
var ab97Data []byte

type Ab97Key struct {
	Ab97Family            string  `json:"family"`
	Ab97Variant           string  `json:"variant"`
	Ab97Primary           uint8   `json:"primary"`
	Ab97Secondary         *uint8  `json:"secondary"`
	Ab97Operation         *uint8  `json:"operation"`
	Ab97Direction         string  `json:"direction"`
	Ab97Field             string  `json:"field"`
	Ab97Group             *uint8  `json:"group"`
	Ab97RequestedInstance *uint8  `json:"requested_instance"`
	Ab97Register          *uint16 `json:"register"`
	Ab97CommandRole       *string `json:"command_role"`
}

type Ab97Entry struct {
	Ab97CoordinateBasis *string           `json:"coordinate_basis"`
	Ab97Key             Ab97Key           `json:"key"`
	Ab97Position        *uint8            `json:"position"`
	Ab97WidthBits       *uint8            `json:"width_bits"`
	Ab97Bit             *uint8            `json:"bit"`
	Ab97Divisor         *float64          `json:"divisor"`
	Ab97Unit            *string           `json:"unit"`
	Ab97FormalUnit      *string           `json:"formal_unit"`
	Ab97BooleanDefault  *bool             `json:"boolean_default"`
	Ab97Default         *int              `json:"default"`
	Ab97Values          map[string]string `json:"values"`
	Ab97Label           *string           `json:"label"`
	Ab97LabelGrade      string            `json:"label_grade"`
	Ab97Signed          *bool             `json:"signed"`
	Ab97ByteOrder       *string           `json:"byte_order"`
	Ab97Reserved        bool              `json:"reserved"`
	Ab97Grade           string            `json:"grade"`
	Ab97Qualification   string            `json:"qualification"`
	Ab97RawPolicy       string            `json:"raw_policy"`
	Ab97Unknown         []string          `json:"unknown"`
	Ab97NativeContract  string            `json:"native_contract"`
}

// The embedded catalog is validated once and never exposed to callers.
var ab97Catalog = ab97sync.OnceValues(ab97LoadEntries)

func Ab97Entries() ([]Ab97Entry, error) {
	ab97Entries, ab97Err := ab97Catalog()
	if ab97Err != nil {
		return nil, ab97Err
	}
	ab97Result := make([]Ab97Entry, len(ab97Entries))
	for ab97Index, ab97Entry := range ab97Entries {
		ab97Result[ab97Index] = ab97CloneEntry(ab97Entry)
	}
	return ab97Result, nil
}

func ab97LoadEntries() ([]Ab97Entry, error) {
	ab97Decoder := ab97json.NewDecoder(ab97bytes.NewReader(ab97Data))
	ab97Decoder.DisallowUnknownFields()
	var ab97Entries []Ab97Entry
	if ab97Err := ab97Decoder.Decode(&ab97Entries); ab97Err != nil {
		return nil, ab97Err
	}
	var ab97Tail any
	if ab97Err := ab97Decoder.Decode(&ab97Tail); ab97Err != ab97io.EOF {
		return nil, ab97fmt.Errorf("ab97 trailing data")
	}
	if len(ab97Entries) == 0 {
		return nil, ab97fmt.Errorf("ab97 empty catalog")
	}
	for ab97Index, ab97Entry := range ab97Entries {
		if ab97Entry.Ab97Key.Ab97Family == "" || ab97Entry.Ab97Key.Ab97Variant == "" || ab97Entry.Ab97Key.Ab97Field == "" || (ab97Entry.Ab97Key.Ab97Direction != "request" && ab97Entry.Ab97Key.Ab97Direction != "positive") || ab97Entry.Ab97Grade != "Hypothesis" || ab97Entry.Ab97Qualification != "Unknown" || ab97Entry.Ab97RawPolicy == "" || ab97Entry.Ab97NativeContract == "" || len(ab97Entry.Ab97Unknown) == 0 {
			return nil, ab97fmt.Errorf("ab97 invalid entry %d", ab97Index)
		}
		if ab97Entry.Ab97Divisor != nil && *ab97Entry.Ab97Divisor <= 0 {
			return nil, ab97fmt.Errorf("ab97 invalid divisor")
		}
		if ab97Entry.Ab97Bit != nil && (*ab97Entry.Ab97Bit > 7 || ab97Entry.Ab97WidthBits == nil || *ab97Entry.Ab97WidthBits != 1) {
			return nil, ab97fmt.Errorf("ab97 invalid bit")
		}
		for ab97Previous := 0; ab97Previous < ab97Index; ab97Previous++ {
			if ab97reflect.DeepEqual(ab97Entry.Ab97Key, ab97Entries[ab97Previous].Ab97Key) {
				return nil, ab97fmt.Errorf("ab97 duplicate key")
			}
		}
	}
	return ab97Entries, nil
}

func Ab97Lookup(ab97Key Ab97Key) (Ab97Entry, bool, error) {
	ab97Entries, ab97Err := ab97Catalog()
	if ab97Err != nil {
		return Ab97Entry{}, false, ab97Err
	}
	for _, ab97Entry := range ab97Entries {
		if ab97reflect.DeepEqual(ab97Key, ab97Entry.Ab97Key) {
			return ab97CloneEntry(ab97Entry), true, nil
		}
	}
	return Ab97Entry{}, false, nil
}

func ab97ClonePointer[T any](ab97Value *T) *T {
	if ab97Value == nil {
		return nil
	}
	ab97Copy := *ab97Value
	return &ab97Copy
}

func ab97CloneEntry(ab97Entry Ab97Entry) Ab97Entry {
	ab97Entry.Ab97Key.Ab97Secondary = ab97ClonePointer(ab97Entry.Ab97Key.Ab97Secondary)
	ab97Entry.Ab97Key.Ab97Operation = ab97ClonePointer(ab97Entry.Ab97Key.Ab97Operation)
	ab97Entry.Ab97Key.Ab97Group = ab97ClonePointer(ab97Entry.Ab97Key.Ab97Group)
	ab97Entry.Ab97Key.Ab97RequestedInstance = ab97ClonePointer(ab97Entry.Ab97Key.Ab97RequestedInstance)
	ab97Entry.Ab97Key.Ab97Register = ab97ClonePointer(ab97Entry.Ab97Key.Ab97Register)
	ab97Entry.Ab97Key.Ab97CommandRole = ab97ClonePointer(ab97Entry.Ab97Key.Ab97CommandRole)
	ab97Entry.Ab97CoordinateBasis = ab97ClonePointer(ab97Entry.Ab97CoordinateBasis)
	ab97Entry.Ab97Position = ab97ClonePointer(ab97Entry.Ab97Position)
	ab97Entry.Ab97WidthBits = ab97ClonePointer(ab97Entry.Ab97WidthBits)
	ab97Entry.Ab97Bit = ab97ClonePointer(ab97Entry.Ab97Bit)
	ab97Entry.Ab97Divisor = ab97ClonePointer(ab97Entry.Ab97Divisor)
	ab97Entry.Ab97Unit = ab97ClonePointer(ab97Entry.Ab97Unit)
	ab97Entry.Ab97FormalUnit = ab97ClonePointer(ab97Entry.Ab97FormalUnit)
	ab97Entry.Ab97BooleanDefault = ab97ClonePointer(ab97Entry.Ab97BooleanDefault)
	ab97Entry.Ab97Default = ab97ClonePointer(ab97Entry.Ab97Default)
	ab97Entry.Ab97Label = ab97ClonePointer(ab97Entry.Ab97Label)
	ab97Entry.Ab97Signed = ab97ClonePointer(ab97Entry.Ab97Signed)
	ab97Entry.Ab97ByteOrder = ab97ClonePointer(ab97Entry.Ab97ByteOrder)
	ab97Entry.Ab97Values = ab97maps.Clone(ab97Entry.Ab97Values)
	ab97Entry.Ab97Unknown = ab97slices.Clone(ab97Entry.Ab97Unknown)
	return ab97Entry
}
