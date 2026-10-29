package vaillant

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Ab05Grade string

const (
	Ab05Unknown    Ab05Grade = "Unknown"
	Ab05Hypothesis Ab05Grade = "Hypothesis"
	Ab05Proven     Ab05Grade = "Proven"
)

type Ab05Key struct {
	Ab05Family      string `json:"ab05_family"`
	Ab05Variant     string `json:"ab05_variant"`
	Ab05Direction   string `json:"ab05_direction"`
	Ab05Opcode      uint8  `json:"ab05_opcode"`
	Ab05Group       uint8  `json:"ab05_group"`
	Ab05Instance    uint8  `json:"ab05_instance"`
	Ab05Register    uint16 `json:"ab05_register"`
	Ab05Alternative string `json:"ab05_alternative"`
}

// UnmarshalJSON requires explicit numeric selectors, including explicit zero,
// while keeping Ab05Key comparable for exact lookup and duplicate detection.
func (ab05Key *Ab05Key) UnmarshalJSON(ab05Data []byte) error {
	type ab05KeyValue Ab05Key
	var ab05Wire struct {
		ab05KeyValue
		Ab05Opcode   *uint8  `json:"ab05_opcode"`
		Ab05Group    *uint8  `json:"ab05_group"`
		Ab05Instance *uint8  `json:"ab05_instance"`
		Ab05Register *uint16 `json:"ab05_register"`
	}
	ab05Decoder := json.NewDecoder(bytes.NewReader(ab05Data))
	ab05Decoder.DisallowUnknownFields()
	if ab05Error := ab05Decoder.Decode(&ab05Wire); ab05Error != nil {
		return ab05Error
	}
	if ab05Wire.Ab05Opcode == nil || ab05Wire.Ab05Group == nil || ab05Wire.Ab05Instance == nil || ab05Wire.Ab05Register == nil {
		return fmt.Errorf("ab05 key: opcode, group, instance and register must be supplied and non-null")
	}
	ab05Value := Ab05Key(ab05Wire.ab05KeyValue)
	ab05Value.Ab05Opcode = *ab05Wire.Ab05Opcode
	ab05Value.Ab05Group = *ab05Wire.Ab05Group
	ab05Value.Ab05Instance = *ab05Wire.Ab05Instance
	ab05Value.Ab05Register = *ab05Wire.Ab05Register
	*ab05Key = ab05Value
	return nil
}

type Ab05Field struct {
	Ab05Name        string           `json:"ab05_name"`
	Ab05Position    *int             `json:"ab05_position,omitempty"`
	Ab05Width       *int             `json:"ab05_width,omitempty"`
	Ab05Codec       string           `json:"ab05_codec"`
	Ab05Unit        string           `json:"ab05_unit,omitempty"`
	Ab05Enum        map[int64]string `json:"ab05_enum,omitempty"`
	Ab05Sentinel    []byte           `json:"ab05_sentinel,omitempty"`
	Ab05Grade       Ab05Grade        `json:"ab05_grade,omitempty"`
	Ab05RawFallback []byte           `json:"ab05_raw_fallback"`
}

type Ab05Record struct {
	Ab05Key         Ab05Key     `json:"ab05_key"`
	Ab05BodyWidth   *int        `json:"ab05_body_width,omitempty"`
	Ab05Grade       Ab05Grade   `json:"ab05_grade,omitempty"`
	Ab05RawFallback []byte      `json:"ab05_raw_fallback"`
	Ab05Fields      []Ab05Field `json:"ab05_fields,omitempty"`
}

func ab05ValidGrade(ab05Grade Ab05Grade) bool {
	return ab05Grade == "" || ab05Grade == Ab05Unknown || ab05Grade == Ab05Hypothesis || ab05Grade == Ab05Proven
}

func Ab05Validate(ab05Records []Ab05Record) error {
	ab05Seen := make(map[Ab05Key]bool)
	for ab05Index, ab05Record := range ab05Records {
		ab05Key := ab05Record.Ab05Key
		if ab05Key.Ab05Family == "" || ab05Key.Ab05Variant == "" || ab05Key.Ab05Alternative == "" || (ab05Key.Ab05Direction != "request" && ab05Key.Ab05Direction != "response") {
			return fmt.Errorf("ab05 record %d: invalid key dimension", ab05Index)
		}
		if ab05Seen[ab05Record.Ab05Key] {
			return fmt.Errorf("ab05 record %d: duplicate alternative key", ab05Index)
		}
		ab05Seen[ab05Record.Ab05Key] = true
		if !ab05ValidGrade(ab05Record.Ab05Grade) {
			return fmt.Errorf("ab05 record %d: invalid grade", ab05Index)
		}
		if ab05Record.Ab05BodyWidth != nil && *ab05Record.Ab05BodyWidth < 0 {
			return fmt.Errorf("ab05 record %d: negative body width", ab05Index)
		}
		for ab05FieldIndex, ab05Field := range ab05Record.Ab05Fields {
			if !ab05ValidGrade(ab05Field.Ab05Grade) {
				return fmt.Errorf("ab05 record %d field %d: invalid grade", ab05Index, ab05FieldIndex)
			}
			if (ab05Field.Ab05Position != nil && *ab05Field.Ab05Position < 0) || (ab05Field.Ab05Width != nil && *ab05Field.Ab05Width <= 0) {
				return fmt.Errorf("ab05 record %d field %d: invalid extent", ab05Index, ab05FieldIndex)
			}
			if ab05Record.Ab05BodyWidth != nil && ab05Field.Ab05Position != nil {
				ab05Remaining := *ab05Record.Ab05BodyWidth - *ab05Field.Ab05Position
				if ab05Remaining <= 0 || (ab05Field.Ab05Width != nil && *ab05Field.Ab05Width > ab05Remaining) {
					return fmt.Errorf("ab05 record %d field %d: extent exceeds body", ab05Index, ab05FieldIndex)
				}
			} else if ab05Record.Ab05BodyWidth != nil && ab05Field.Ab05Width != nil && *ab05Field.Ab05Width > *ab05Record.Ab05BodyWidth {
				return fmt.Errorf("ab05 record %d field %d: width exceeds body", ab05Index, ab05FieldIndex)
			}
		}
	}
	return nil
}

func ab05CopyInt(ab05Value *int) *int {
	if ab05Value == nil {
		return nil
	}
	ab05Copy := *ab05Value
	return &ab05Copy
}

func ab05ClaimGrade(ab05Grade Ab05Grade) Ab05Grade {
	if ab05Grade == "" {
		return Ab05Unknown
	}
	return ab05Grade
}

func Ab05Query(ab05Records []Ab05Record, ab05Key Ab05Key) (Ab05Record, bool, error) {
	if ab05Error := Ab05Validate(ab05Records); ab05Error != nil {
		return Ab05Record{}, false, ab05Error
	}
	for _, ab05Record := range ab05Records {
		if ab05Record.Ab05Key != ab05Key {
			continue
		}
		ab05Record.Ab05BodyWidth = ab05CopyInt(ab05Record.Ab05BodyWidth)
		ab05Record.Ab05Grade = ab05ClaimGrade(ab05Record.Ab05Grade)
		ab05Record.Ab05RawFallback = bytes.Clone(ab05Record.Ab05RawFallback)
		ab05Record.Ab05Fields = append([]Ab05Field(nil), ab05Record.Ab05Fields...)
		for ab05Index := range ab05Record.Ab05Fields {
			ab05Field := &ab05Record.Ab05Fields[ab05Index]
			ab05Field.Ab05Position = ab05CopyInt(ab05Field.Ab05Position)
			ab05Field.Ab05Width = ab05CopyInt(ab05Field.Ab05Width)
			ab05Field.Ab05Grade = ab05ClaimGrade(ab05Field.Ab05Grade)
			ab05Field.Ab05RawFallback = bytes.Clone(ab05Field.Ab05RawFallback)
			ab05Field.Ab05Sentinel = bytes.Clone(ab05Field.Ab05Sentinel)
			if ab05Field.Ab05Enum != nil {
				ab05Enum := make(map[int64]string, len(ab05Field.Ab05Enum))
				for ab05Value, ab05Label := range ab05Field.Ab05Enum {
					ab05Enum[ab05Value] = ab05Label
				}
				ab05Field.Ab05Enum = ab05Enum
			}
		}
		return ab05Record, true, nil
	}
	return Ab05Record{}, false, nil
}
