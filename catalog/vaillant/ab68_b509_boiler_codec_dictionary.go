package vaillant

import (
	ab68bytes "bytes"
	_ "embed"
	ab68json "encoding/json"
	ab68fmt "fmt"
	ab68io "io"
	ab68reflect "reflect"
)

type Ab68Key struct {
	Ab68Family    string   `json:"family"`
	Ab68Variant   string   `json:"variant"`
	Ab68Primary   *byte    `json:"primary"`
	Ab68Secondary *byte    `json:"secondary"`
	Ab68Operation *byte    `json:"operation"`
	Ab68Direction string   `json:"direction"`
	Ab68Selector  []uint16 `json:"selector"`
}

type Ab68Entry struct {
	Ab68Key              Ab68Key           `json:"key"`
	Ab68Quantity         string            `json:"quantity"`
	Ab68CodecDeclaration string            `json:"codec_declaration"`
	Ab68WidthBits        *int              `json:"width_bits"`
	Ab68ByteOrder        *string           `json:"byte_order"`
	Ab68Unit             *string           `json:"unit"`
	Ab68Converter        *string           `json:"converter"`
	Ab68Mapping          map[string]string `json:"mapping"`
	Ab68Sentinel         *int64            `json:"sentinel"`
	Ab68Grade            string            `json:"grade"`
	Ab68Conflicts        []string          `json:"conflicts"`
	Ab68Evidence         map[string]string `json:"evidence"`
}

//go:embed data/ab68_b509_boiler_codec_dictionary.json
var ab68JSON []byte

func ab68Parse(ab68Data []byte) ([]Ab68Entry, error) {
	var ab68Entries []Ab68Entry
	ab68Decoder := ab68json.NewDecoder(ab68bytes.NewReader(ab68Data))
	ab68Decoder.DisallowUnknownFields()
	if ab68Error := ab68Decoder.Decode(&ab68Entries); ab68Error != nil {
		return nil, ab68Error
	}
	if ab68Error := ab68Decoder.Decode(new(any)); ab68Error != ab68io.EOF {
		return nil, ab68fmt.Errorf("ab68: trailing JSON: %v", ab68Error)
	}
	if len(ab68Entries) == 0 {
		return nil, ab68fmt.Errorf("ab68: empty dictionary")
	}
	ab68Seen := make(map[string]bool)
	for ab68Index, ab68Entry := range ab68Entries {
		ab68Key := ab68Entry.Ab68Key
		if ab68Key.Ab68Family == "" || ab68Key.Ab68Variant == "" || ab68Key.Ab68Primary == nil || ab68Key.Ab68Secondary == nil || ab68Key.Ab68Operation == nil || ab68Key.Ab68Direction == "" || len(ab68Key.Ab68Selector) != 2 {
			return nil, ab68fmt.Errorf("ab68: incomplete key at row %d", ab68Index)
		}
		for _, ab68Byte := range ab68Key.Ab68Selector {
			if ab68Byte > 255 {
				return nil, ab68fmt.Errorf("ab68: selector byte out of range at row %d", ab68Index)
			}
		}
		if ab68Entry.Ab68Quantity == "" || ab68Entry.Ab68CodecDeclaration == "" || ab68Entry.Ab68Grade != "Hypothesis" {
			return nil, ab68fmt.Errorf("ab68: incomplete declaration at row %d", ab68Index)
		}
		if ab68Entry.Ab68Evidence["metadata"] != "Proven" {
			return nil, ab68fmt.Errorf("ab68: missing metadata grade at row %d", ab68Index)
		}
		for _, ab68Axis := range []string{"qualification", "physical_support", "presence", "freshness", "effects", "persistence", "delivery", "execution"} {
			if ab68Entry.Ab68Evidence[ab68Axis] != "Unknown" {
				return nil, ab68fmt.Errorf("ab68: invalid %s grade at row %d", ab68Axis, ab68Index)
			}
		}
		ab68Encoded, ab68Error := ab68json.Marshal(ab68Key)
		if ab68Error != nil {
			return nil, ab68Error
		}
		if ab68Seen[string(ab68Encoded)] {
			return nil, ab68fmt.Errorf("ab68: duplicate key at row %d", ab68Index)
		}
		ab68Seen[string(ab68Encoded)] = true
	}
	return ab68Entries, nil
}

func Ab68Entries() ([]Ab68Entry, error) {
	return ab68Parse(ab68JSON)
}

func ab68Lookup(ab68Data []byte, ab68Key Ab68Key) (Ab68Entry, bool, error) {
	ab68Entries, ab68Error := ab68Parse(ab68Data)
	if ab68Error != nil {
		return Ab68Entry{}, false, ab68Error
	}
	for _, ab68Entry := range ab68Entries {
		if ab68reflect.DeepEqual(ab68Entry.Ab68Key, ab68Key) {
			return ab68Entry, true, nil
		}
	}
	return Ab68Entry{}, false, nil
}

func Ab68Lookup(ab68Key Ab68Key) (Ab68Entry, bool, error) {
	return ab68Lookup(ab68JSON, ab68Key)
}
