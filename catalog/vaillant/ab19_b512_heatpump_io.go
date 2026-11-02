package vaillant

import (
	_ "embed"
	ab19json "encoding/json"
	ab19fmt "fmt"
	ab19reflect "reflect"
)

//go:embed data/ab19_b512_heatpump_io.json
var ab19Data string

type Ab19Key struct {
	Ab19Family        string  `json:"family"`
	Ab19Variant       string  `json:"variant"`
	Ab19Primary       uint8   `json:"primary"`
	Ab19Secondary     *uint8  `json:"secondary"`
	Ab19Operation     *uint8  `json:"operation"`
	Ab19Direction     string  `json:"direction"`
	Ab19Field         string  `json:"field"`
	Ab19Alternative   string  `json:"alternative"`
	Ab19Group         *uint8  `json:"group"`
	Ab19InstanceScope string  `json:"instance_scope"`
	Ab19Register      *uint16 `json:"register"`
}

func ab19Validate(ab19Input string) ([]ab19json.RawMessage, []Ab19Key, error) {
	var ab19Rows []ab19json.RawMessage
	if ab19Err := ab19json.Unmarshal([]byte(ab19Input), &ab19Rows); ab19Err != nil {
		return nil, nil, ab19fmt.Errorf("ab19 catalog: %w", ab19Err)
	}
	if len(ab19Rows) == 0 {
		return nil, nil, ab19fmt.Errorf("ab19 catalog: empty array")
	}
	ab19Keys := make([]Ab19Key, 0, len(ab19Rows))
	ab19IDs := make(map[string]bool)
	for _, ab19Row := range ab19Rows {
		var ab19Fields map[string]ab19json.RawMessage
		if ab19Err := ab19json.Unmarshal(ab19Row, &ab19Fields); ab19Err != nil {
			return nil, nil, ab19Err
		}
		var ab19KeyFields map[string]ab19json.RawMessage
		if ab19Err := ab19json.Unmarshal(ab19Fields["key"], &ab19KeyFields); ab19Err != nil {
			return nil, nil, ab19Err
		}
		for _, ab19Axis := range []string{"family", "variant", "primary", "secondary", "operation", "direction", "field", "alternative", "group", "instance_scope", "register"} {
			if _, ab19Present := ab19KeyFields[ab19Axis]; !ab19Present {
				return nil, nil, ab19fmt.Errorf("ab19 catalog: missing %s", ab19Axis)
			}
		}
		var ab19Key Ab19Key
		if ab19Err := ab19json.Unmarshal(ab19Fields["key"], &ab19Key); ab19Err != nil {
			return nil, nil, ab19Err
		}
		if len(ab19KeyFields) != 11 || ab19Key.Ab19Family == "" || ab19Key.Ab19Variant == "" || ab19Key.Ab19Direction == "" || ab19Key.Ab19Field == "" || ab19Key.Ab19Alternative == "" || ab19Key.Ab19InstanceScope == "" || string(ab19KeyFields["primary"]) == "null" {
			return nil, nil, ab19fmt.Errorf("ab19 catalog: incomplete identity")
		}
		for _, ab19Previous := range ab19Keys {
			if ab19reflect.DeepEqual(ab19Previous, ab19Key) {
				return nil, nil, ab19fmt.Errorf("ab19 catalog: duplicate key")
			}
		}
		var ab19ID string
		if ab19Err := ab19json.Unmarshal(ab19Fields["id"], &ab19ID); ab19Err != nil {
			return nil, nil, ab19Err
		}
		if len(ab19ID) < 5 || ab19ID[:5] != "ab19_" || ab19IDs[ab19ID] {
			return nil, nil, ab19fmt.Errorf("ab19 catalog: invalid ID")
		}
		ab19IDs[ab19ID] = true
		ab19Keys = append(ab19Keys, ab19Key)
	}
	return ab19Rows, ab19Keys, nil
}

func ab19Entries(ab19Input string) (ab19json.RawMessage, error) {
	if _, _, ab19Err := ab19Validate(ab19Input); ab19Err != nil {
		return nil, ab19Err
	}
	return append(ab19json.RawMessage(nil), ab19Input...), nil
}

func Ab19Entries() (ab19json.RawMessage, error) { return ab19Entries(ab19Data) }

func ab19Lookup(ab19Input string, ab19Key Ab19Key) (ab19json.RawMessage, bool, error) {
	ab19Rows, ab19Keys, ab19Err := ab19Validate(ab19Input)
	if ab19Err != nil {
		return nil, false, ab19Err
	}
	for ab19Index, ab19Candidate := range ab19Keys {
		if ab19reflect.DeepEqual(ab19Candidate, ab19Key) {
			return append(ab19json.RawMessage(nil), ab19Rows[ab19Index]...), true, nil
		}
	}
	return nil, false, nil
}

func Ab19Lookup(ab19Key Ab19Key) (ab19json.RawMessage, bool, error) {
	return ab19Lookup(ab19Data, ab19Key)
}
