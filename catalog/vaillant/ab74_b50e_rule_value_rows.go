package vaillant

import (
	ab74bytes "bytes"
	_ "embed"
	ab74json "encoding/json"
	ab74fmt "fmt"
	ab74io "io"
	ab74reflect "reflect"
)

//go:embed data/ab74_b50e_rule_value_rows.json
var ab74RowsJSON []byte

type Ab74Key struct {
	Ab74Family    *string `json:"family"`
	Ab74Variant   *string `json:"variant"`
	Ab74Primary   *byte   `json:"primary"`
	Ab74Secondary *byte   `json:"secondary"`
	Ab74Operation *byte   `json:"operation"`
	Ab74Direction *string `json:"direction"`
	Ab74Selector  *byte   `json:"selector"`
	Ab74Field     *string `json:"field"`
}

type Ab74Field struct {
	Ab74Offset         *int              `json:"offset"`
	Ab74Width          *int              `json:"width"`
	Ab74Encoding       string            `json:"encoding"`
	Ab74ByteOrder      string            `json:"byte_order"`
	Ab74Default        *float64          `json:"default"`
	Ab74Enumeration    map[string]string `json:"enumeration"`
	Ab74Scale          *float64          `json:"scale"`
	Ab74Unit           *string           `json:"unit"`
	Ab74CarrierMinimum *int              `json:"carrier_minimum"`
	Ab74CarrierMaximum *int              `json:"carrier_maximum"`
	Ab74ScaledMinimum  *float64          `json:"scaled_minimum"`
	Ab74ScaledMaximum  *float64          `json:"scaled_maximum"`
}

type Ab74Entry struct {
	Ab74Key                  Ab74Key    `json:"key"`
	Ab74Value                Ab74Field  `json:"value"`
	Ab74Companion            *Ab74Field `json:"companion"`
	Ab74Base                 string     `json:"base"`
	Ab74Grade                string     `json:"grade"`
	Ab74MetadataPresence     string     `json:"metadata_presence"`
	Ab74PhysicalSupport      string     `json:"physical_support"`
	Ab74Effects              string     `json:"effects"`
	Ab74Persistence          string     `json:"persistence"`
	Ab74Delivery             string     `json:"delivery"`
	Ab74DefaultQualification string     `json:"default_qualification"`
	Ab74RawFallback          string     `json:"raw_fallback"`
}

func ab74Decode(ab74Input []byte) ([]Ab74Entry, error) {
	var ab74Entries []Ab74Entry
	ab74Decoder := ab74json.NewDecoder(ab74bytes.NewReader(ab74Input))
	ab74Decoder.DisallowUnknownFields()
	if ab74Error := ab74Decoder.Decode(&ab74Entries); ab74Error != nil {
		return nil, ab74Error
	}
	if ab74Error := ab74Decoder.Decode(new(any)); ab74Error != ab74io.EOF {
		return nil, ab74fmt.Errorf("ab74 trailing JSON: %v", ab74Error)
	}
	ab74Canonical, ab74Error := ab74json.Marshal(ab74Entries)
	if ab74Error != nil {
		return nil, ab74Error
	}
	var ab74RawShape, ab74TypedShape any
	if ab74Error = ab74json.Unmarshal(ab74Input, &ab74RawShape); ab74Error != nil {
		return nil, ab74Error
	}
	if ab74Error = ab74json.Unmarshal(ab74Canonical, &ab74TypedShape); ab74Error != nil {
		return nil, ab74Error
	}
	if !ab74CompleteShape(ab74RawShape, ab74TypedShape) {
		return nil, ab74fmt.Errorf("ab74 missing explicit metadata field")
	}
	if len(ab74Entries) != 19 {
		return nil, ab74fmt.Errorf("ab74 expected nineteen rows")
	}
	ab74Seen := make(map[byte]bool)
	for _, ab74Entry := range ab74Entries {
		ab74Key := ab74Entry.Ab74Key
		if ab74Key.Ab74Family == nil || *ab74Key.Ab74Family != "vaillant" || ab74Key.Ab74Variant == nil || *ab74Key.Ab74Variant != "B50E" || ab74Key.Ab74Primary == nil || *ab74Key.Ab74Primary != 0xB5 || ab74Key.Ab74Secondary == nil || *ab74Key.Ab74Secondary != 0x0E || ab74Key.Ab74Operation == nil || *ab74Key.Ab74Operation != 2 || ab74Key.Ab74Direction == nil || *ab74Key.Ab74Direction != "positive" || ab74Key.Ab74Field == nil || *ab74Key.Ab74Field != "rule_value_row" || ab74Key.Ab74Selector == nil || *ab74Key.Ab74Selector > 18 {
			return nil, ab74fmt.Errorf("ab74 invalid row identity")
		}
		ab74Selector := *ab74Key.Ab74Selector
		if ab74Seen[ab74Selector] {
			return nil, ab74fmt.Errorf("ab74 duplicate row key")
		}
		ab74Seen[ab74Selector] = true
		if ab74Entry.Ab74Base != "row-relative" || ab74Entry.Ab74Grade != "Hypothesis" || ab74Entry.Ab74MetadataPresence != "Proven" || ab74Entry.Ab74PhysicalSupport != "Unknown" || ab74Entry.Ab74Effects != "Unknown" || ab74Entry.Ab74Persistence != "Unknown" || ab74Entry.Ab74Delivery != "Unknown" || ab74Entry.Ab74DefaultQualification != "Unknown" || ab74Entry.Ab74RawFallback == "" {
			return nil, ab74fmt.Errorf("ab74 incomplete qualification")
		}
		if ab74Entry.Ab74Value.Ab74Offset == nil || *ab74Entry.Ab74Value.Ab74Offset != 0 {
			return nil, ab74fmt.Errorf("ab74 invalid value base")
		}
		for _, ab74Field := range []*Ab74Field{&ab74Entry.Ab74Value, ab74Entry.Ab74Companion} {
			if ab74Field == nil {
				continue
			}
			if ab74Field.Ab74Offset == nil || *ab74Field.Ab74Offset < 0 || ab74Field.Ab74Encoding == "" || ab74Field.Ab74ByteOrder != "Unknown" || ab74Field.Ab74Default != nil || (ab74Field.Ab74Width != nil && *ab74Field.Ab74Width <= 0) {
				return nil, ab74fmt.Errorf("ab74 invalid field declaration")
			}
		}
		if ab74Selector == 0 || ab74Selector == 15 || ab74Selector == 16 {
			if ab74Entry.Ab74Value.Ab74Width != nil || ab74Entry.Ab74Value.Ab74Encoding != "opaque" {
				return nil, ab74fmt.Errorf("ab74 opaque width must remain absent")
			}
		} else if ab74Entry.Ab74Value.Ab74Width == nil {
			return nil, ab74fmt.Errorf("ab74 missing known width")
		}
		if (ab74Selector >= 9 && ab74Selector <= 14 || ab74Selector >= 17) != (ab74Entry.Ab74Companion != nil) {
			return nil, ab74fmt.Errorf("ab74 missing or extra companion")
		}
	}
	return ab74Entries, nil
}

func Ab74Entries() ([]Ab74Entry, error) { return ab74Decode(ab74RowsJSON) }

func Ab74Lookup(ab74Key Ab74Key) (Ab74Entry, bool, error) {
	ab74Entries, ab74Error := Ab74Entries()
	if ab74Error != nil {
		return Ab74Entry{}, false, ab74Error
	}
	for _, ab74Entry := range ab74Entries {
		if ab74reflect.DeepEqual(ab74Entry.Ab74Key, ab74Key) {
			return ab74Entry, true, nil
		}
	}
	return Ab74Entry{}, false, nil
}

func ab74CompleteShape(ab74Raw, ab74Typed any) bool {
	switch ab74Expected := ab74Typed.(type) {
	case map[string]any:
		ab74Actual, ab74OK := ab74Raw.(map[string]any)
		if !ab74OK || len(ab74Actual) != len(ab74Expected) {
			return false
		}
		for ab74Name, ab74Value := range ab74Expected {
			ab74Input, ab74Present := ab74Actual[ab74Name]
			if !ab74Present || !ab74CompleteShape(ab74Input, ab74Value) {
				return false
			}
		}
	case []any:
		ab74Actual, ab74OK := ab74Raw.([]any)
		if !ab74OK || len(ab74Actual) != len(ab74Expected) {
			return false
		}
		for ab74Index := range ab74Expected {
			if !ab74CompleteShape(ab74Actual[ab74Index], ab74Expected[ab74Index]) {
				return false
			}
		}
	}
	return true
}
