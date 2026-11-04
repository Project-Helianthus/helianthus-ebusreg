package vaillant

import (
	ab15embed "embed"
	ab15json "encoding/json"
	ab15fmt "fmt"
	ab15reflect "reflect"
	ab15strings "strings"
)

//go:embed data/ab15_b504_circuit_roles.json
var ab15Files ab15embed.FS

type Ab15Key struct {
	Ab15Family        string  `json:"family"`
	Ab15Variant       string  `json:"variant"`
	Ab15Primary       uint8   `json:"primary"`
	Ab15Secondary     *uint8  `json:"secondary"`
	Ab15Operation     *uint8  `json:"operation"`
	Ab15Direction     string  `json:"direction"`
	Ab15Field         string  `json:"field"`
	Ab15Alternative   string  `json:"alternative"`
	Ab15Group         *uint8  `json:"group"`
	Ab15InstanceScope string  `json:"instance_scope"`
	Ab15Register      *uint16 `json:"register"`
}

type ab15Row struct {
	Ab15ID                    string              `json:"id"`
	Ab15Key                   Ab15Key             `json:"key"`
	Ab15Grade                 string              `json:"grade"`
	Ab15MetadataPresence      string              `json:"metadata_presence"`
	Ab15NativeContract        string              `json:"native_contract"`
	Ab15PhysicalQualification string              `json:"physical_qualification"`
	Ab15Role                  string              `json:"role"`
	Ab15AppliesTo             []string            `json:"applies_to"`
	Ab15Relations             []string            `json:"relations"`
	Ab15Condition             string              `json:"condition"`
	Ab15Bit                   *uint8              `json:"bit"`
	Ab15Unknown               []string            `json:"unknown"`
	Ab15Raw                   ab15json.RawMessage `json:"raw"`
}

func Ab15Entries() (ab15json.RawMessage, error) {
	ab15Data, ab15Err := ab15Files.ReadFile("data/ab15_b504_circuit_roles.json")
	if ab15Err != nil {
		return nil, ab15Err
	}
	return ab15Entries(ab15Data)
}

func Ab15Lookup(ab15Key Ab15Key) (ab15json.RawMessage, bool, error) {
	ab15Data, ab15Err := ab15Files.ReadFile("data/ab15_b504_circuit_roles.json")
	if ab15Err != nil {
		return nil, false, ab15Err
	}
	return ab15Lookup(ab15Data, ab15Key)
}

func ab15Entries(ab15Data []byte) (ab15json.RawMessage, error) {
	var ab15Rows []ab15json.RawMessage
	if ab15Err := ab15json.Unmarshal(ab15Data, &ab15Rows); ab15Err != nil {
		return nil, ab15Err
	}
	if len(ab15Rows) == 0 {
		return nil, ab15fmt.Errorf("ab15: empty catalog")
	}
	ab15Keys := make([]Ab15Key, 0, len(ab15Rows))
	ab15IDs := map[string]bool{}
	for _, ab15Raw := range ab15Rows {
		var ab15Entry ab15Row
		if ab15Err := ab15json.Unmarshal(ab15Raw, &ab15Entry); ab15Err != nil {
			return nil, ab15Err
		}
		var ab15Object map[string]ab15json.RawMessage
		if ab15Err := ab15json.Unmarshal(ab15Raw, &ab15Object); ab15Err != nil {
			return nil, ab15Err
		}
		var ab15Identity map[string]ab15json.RawMessage
		if ab15Err := ab15json.Unmarshal(ab15Object["key"], &ab15Identity); ab15Err != nil {
			return nil, ab15Err
		}
		for _, ab15Axis := range []string{"family", "variant", "primary", "secondary", "operation", "direction", "field", "alternative", "group", "instance_scope", "register"} {
			ab15Value, ab15Present := ab15Identity[ab15Axis]
			if !ab15Present || (string(ab15Value) == "null" && ab15Axis != "secondary" && ab15Axis != "operation" && ab15Axis != "group" && ab15Axis != "register") {
				return nil, ab15fmt.Errorf("ab15: missing identity axis %s", ab15Axis)
			}
		}
		if len(ab15Identity) != 11 || ab15Entry.Ab15Key.Ab15Family == "" || ab15Entry.Ab15Key.Ab15Variant == "" || ab15Entry.Ab15Key.Ab15Direction == "" || ab15Entry.Ab15Key.Ab15Field == "" || ab15Entry.Ab15Key.Ab15Alternative == "" || ab15Entry.Ab15Key.Ab15InstanceScope == "" {
			return nil, ab15fmt.Errorf("ab15: incomplete identity")
		}
		if !ab15strings.HasPrefix(ab15Entry.Ab15ID, "ab15-") || ab15IDs[ab15Entry.Ab15ID] {
			return nil, ab15fmt.Errorf("ab15: invalid or duplicate id")
		}
		ab15IDs[ab15Entry.Ab15ID] = true
		if ab15Entry.Ab15Grade != "Hypothesis" || ab15Entry.Ab15MetadataPresence != "Proven" || ab15Entry.Ab15PhysicalQualification != "Unknown" || ab15Entry.Ab15NativeContract == "" || ab15Entry.Ab15Role == "" || len(ab15Entry.Ab15AppliesTo) == 0 || len(ab15Entry.Ab15Relations) == 0 || len(ab15Entry.Ab15Unknown) == 0 || len(ab15Entry.Ab15Raw) == 0 || string(ab15Entry.Ab15Raw) == "null" {
			return nil, ab15fmt.Errorf("ab15: incomplete evidence boundary")
		}
		for _, ab15Previous := range ab15Keys {
			if ab15reflect.DeepEqual(ab15Previous, ab15Entry.Ab15Key) {
				return nil, ab15fmt.Errorf("ab15: duplicate key")
			}
		}
		ab15Keys = append(ab15Keys, ab15Entry.Ab15Key)
	}
	return append(ab15json.RawMessage(nil), ab15Data...), nil
}

func ab15Lookup(ab15Data []byte, ab15Key Ab15Key) (ab15json.RawMessage, bool, error) {
	ab15Validated, ab15Err := ab15Entries(ab15Data)
	if ab15Err != nil {
		return nil, false, ab15Err
	}
	var ab15Rows []ab15json.RawMessage
	if ab15Err = ab15json.Unmarshal(ab15Validated, &ab15Rows); ab15Err != nil {
		return nil, false, ab15Err
	}
	for _, ab15Raw := range ab15Rows {
		var ab15Entry ab15Row
		if ab15Err = ab15json.Unmarshal(ab15Raw, &ab15Entry); ab15Err != nil {
			return nil, false, ab15Err
		}
		if ab15reflect.DeepEqual(ab15Key, ab15Entry.Ab15Key) {
			return append(ab15json.RawMessage(nil), ab15Raw...), true, nil
		}
	}
	return nil, false, nil
}
