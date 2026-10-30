package vaillant

import (
	ab12Bytes "bytes"
	_ "embed"
	ab12JSON "encoding/json"
	ab12Fmt "fmt"
	ab12IO "io"
)

//go:embed data/ab12_b50f_status_spans.json
var ab12StatusJSON string

type Ab12Key struct {
	Ab12Family            string  `json:"family"`
	Ab12Variant           string  `json:"variant"`
	Ab12Primary           byte    `json:"primary"`
	Ab12Secondary         *byte   `json:"secondary"`
	Ab12Operation         *byte   `json:"operation"`
	Ab12Direction         *string `json:"direction"`
	Ab12Field             string  `json:"field"`
	Ab12Group             *byte   `json:"group"`
	Ab12RequestedInstance *byte   `json:"requested_instance"`
	Ab12Register          *uint16 `json:"register"`
	Ab12CommandRole       *string `json:"command_role"`
}

type Ab12Bit struct {
	Ab12Bit      byte    `json:"bit"`
	Ab12Label    *string `json:"label"`
	Ab12Category *string `json:"category"`
	Ab12Reserved bool    `json:"reserved"`
}

type Ab12Entry struct {
	Ab12Key             Ab12Key   `json:"key"`
	Ab12Position        *int      `json:"position"`
	Ab12WidthBits       *int      `json:"width_bits"`
	Ab12Grade           string    `json:"grade"`
	Ab12Qualification   string    `json:"qualification"`
	Ab12CoordinateGrade string    `json:"coordinate_grade"`
	Ab12AllocatedMask   *byte     `json:"allocated_mask"`
	Ab12ReservedMask    *byte     `json:"reserved_mask"`
	Ab12Bits            []Ab12Bit `json:"bits"`
	Ab12RawFallback     string    `json:"raw_fallback"`
	Ab12Unknown         []string  `json:"unknown"`
}

func ab12OptionalEqual[ab12Value comparable](ab12Left, ab12Right *ab12Value) bool {
	if ab12Left == nil || ab12Right == nil {
		return ab12Left == nil && ab12Right == nil
	}
	return *ab12Left == *ab12Right
}

func ab12KeyEqual(ab12Left, ab12Right Ab12Key) bool {
	return ab12Left.Ab12Family == ab12Right.Ab12Family &&
		ab12Left.Ab12Variant == ab12Right.Ab12Variant &&
		ab12Left.Ab12Primary == ab12Right.Ab12Primary &&
		ab12OptionalEqual(ab12Left.Ab12Secondary, ab12Right.Ab12Secondary) &&
		ab12OptionalEqual(ab12Left.Ab12Operation, ab12Right.Ab12Operation) &&
		ab12OptionalEqual(ab12Left.Ab12Direction, ab12Right.Ab12Direction) &&
		ab12Left.Ab12Field == ab12Right.Ab12Field &&
		ab12OptionalEqual(ab12Left.Ab12Group, ab12Right.Ab12Group) &&
		ab12OptionalEqual(ab12Left.Ab12RequestedInstance, ab12Right.Ab12RequestedInstance) &&
		ab12OptionalEqual(ab12Left.Ab12Register, ab12Right.Ab12Register) &&
		ab12OptionalEqual(ab12Left.Ab12CommandRole, ab12Right.Ab12CommandRole)
}

func Ab12Entries() ([]Ab12Entry, error) {
	var ab12Entries []Ab12Entry
	ab12Decoder := ab12JSON.NewDecoder(ab12Bytes.NewBufferString(ab12StatusJSON))
	ab12Decoder.DisallowUnknownFields()
	if ab12Error := ab12Decoder.Decode(&ab12Entries); ab12Error != nil {
		return nil, ab12Error
	}
	if ab12Error := ab12Decoder.Decode(new(any)); ab12Error != ab12IO.EOF {
		return nil, ab12Fmt.Errorf("ab12: trailing data")
	}
	if len(ab12Entries) != 8 {
		return nil, ab12Fmt.Errorf("ab12: incomplete finite catalog")
	}
	for ab12Index, ab12Entry := range ab12Entries {
		if ab12Entry.Ab12Key.Ab12Primary == 0 || ab12Entry.Ab12Key.Ab12Family == "" || ab12Entry.Ab12Key.Ab12Variant == "" || ab12Entry.Ab12Key.Ab12Field == "" || ab12Entry.Ab12WidthBits == nil || *ab12Entry.Ab12WidthBits != 8 || ab12Entry.Ab12Grade != "Hypothesis" || ab12Entry.Ab12Qualification != "Unknown" || ab12Entry.Ab12CoordinateGrade != "Unknown" || ab12Entry.Ab12Position != nil || ab12Entry.Ab12RawFallback == "" || len(ab12Entry.Ab12Unknown) == 0 {
			return nil, ab12Fmt.Errorf("ab12: invalid entry %d", ab12Index)
		}
		for ab12Previous := 0; ab12Previous < ab12Index; ab12Previous++ {
			if ab12KeyEqual(ab12Entries[ab12Previous].Ab12Key, ab12Entry.Ab12Key) {
				return nil, ab12Fmt.Errorf("ab12: duplicate key")
			}
		}
		if (ab12Index < 6) != (ab12Entry.Ab12AllocatedMask != nil && ab12Entry.Ab12ReservedMask != nil) {
			return nil, ab12Fmt.Errorf("ab12: invalid container/raw partition")
		}
		if ab12Entry.Ab12AllocatedMask == nil || ab12Entry.Ab12ReservedMask == nil {
			if ab12Entry.Ab12AllocatedMask != nil || ab12Entry.Ab12ReservedMask != nil || len(ab12Entry.Ab12Bits) != 0 {
				return nil, ab12Fmt.Errorf("ab12: inconsistent raw byte")
			}
			continue
		}
		if len(ab12Entry.Ab12Bits) != 8 || *ab12Entry.Ab12AllocatedMask&*ab12Entry.Ab12ReservedMask != 0 || *ab12Entry.Ab12AllocatedMask|*ab12Entry.Ab12ReservedMask != 255 {
			return nil, ab12Fmt.Errorf("ab12: incomplete bit partition")
		}
		for ab12Index, ab12Bit := range ab12Entry.Ab12Bits {
			if int(ab12Bit.Ab12Bit) != ab12Index || ab12Bit.Ab12Reserved != (*ab12Entry.Ab12ReservedMask&(1<<ab12Index) != 0) || (ab12Bit.Ab12Reserved && (ab12Bit.Ab12Label != nil || ab12Bit.Ab12Category != nil)) {
				return nil, ab12Fmt.Errorf("ab12: invalid bit")
			}
		}
	}
	return ab12Entries, nil
}

func Ab12Lookup(ab12Key Ab12Key) (Ab12Entry, bool, error) {
	ab12Entries, ab12Error := Ab12Entries()
	if ab12Error != nil {
		return Ab12Entry{}, false, ab12Error
	}
	for _, ab12Entry := range ab12Entries {
		if ab12KeyEqual(ab12Entry.Ab12Key, ab12Key) {
			return ab12Entry, true, nil
		}
	}
	return Ab12Entry{}, false, nil
}
