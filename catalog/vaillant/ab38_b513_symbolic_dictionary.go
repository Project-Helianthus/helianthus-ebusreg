package vaillant

import (
	ab38bytes "bytes"
	_ "embed"
	ab38json "encoding/json"
	ab38fmt "fmt"
	ab38io "io"
)

//go:embed data/ab38_b513_symbolic_dictionary.json
var ab38Data []byte

type Ab38Key struct {
	Ab38Family              string  `json:"family"`
	Ab38Variant             string  `json:"variant"`
	Ab38Primary             uint8   `json:"primary"`
	Ab38Secondary           *uint8  `json:"secondary"`
	Ab38Operation           *uint8  `json:"operation"`
	Ab38Direction           string  `json:"direction"`
	Ab38Field               string  `json:"field"`
	Ab38Alternative         string  `json:"alternative"`
	Ab38Group               *uint8  `json:"group"`
	Ab38InstanceScope       string  `json:"instance_scope"`
	Ab38Register            *uint16 `json:"register"`
	Ab38DictionaryNamespace string  `json:"dictionary_namespace"`
	Ab38DictionaryClass     string  `json:"dictionary_class"`
	Ab38RawAddress          *uint16 `json:"raw_address"`
	Ab38BusFamily           string  `json:"bus_family,omitempty"`
}

type ab38Row struct {
	Ab38Key
	Ab38NativeLabel           string   `json:"native_label"`
	Ab38RawForm               string   `json:"raw_form"`
	Ab38BusFamilies           []string `json:"bus_families"`
	Ab38Grade                 string   `json:"grade"`
	Ab38PhysicalQualification string   `json:"physical_qualification"`
	Ab38ValueCodec            *string  `json:"value_codec"`
	Ab38ExecutableCapability  bool     `json:"executable_capability"`
	Ab38AddressWidthBits      uint8    `json:"address_width_bits"`
	Ab38ValueWidthBits        *uint8   `json:"value_width_bits"`
	Ab38Unknown               []string `json:"unknown"`
	Ab38UnknownAddress        string   `json:"unknown_address"`
	Ab38LabelMeaning          string   `json:"label_meaning"`
	Ab38NativeContract        string   `json:"native_contract"`
}

func ab38Validate(ab38Source []byte) ([]ab38json.RawMessage, map[string]ab38json.RawMessage, error) {
	var ab38Rows []ab38json.RawMessage
	if ab38Err := ab38json.Unmarshal(ab38Source, &ab38Rows); ab38Err != nil {
		return nil, nil, ab38fmt.Errorf("ab38 dictionary: %w", ab38Err)
	}
	if len(ab38Rows) == 0 {
		return nil, nil, ab38fmt.Errorf("ab38 dictionary: expected nonempty array")
	}
	ab38Index := make(map[string]ab38json.RawMessage)
	for _, ab38Raw := range ab38Rows {
		if ab38Err := ab38ValidateObject(ab38Raw); ab38Err != nil {
			return nil, nil, ab38Err
		}
		var ab38Entry ab38Row
		ab38Decoder := ab38json.NewDecoder(ab38bytes.NewReader(ab38Raw))
		ab38Decoder.DisallowUnknownFields()
		if ab38Err := ab38Decoder.Decode(&ab38Entry); ab38Err != nil {
			return nil, nil, ab38fmt.Errorf("ab38 row: %w", ab38Err)
		}
		if ab38Err := ab38Decoder.Decode(new(any)); ab38Err != ab38io.EOF {
			return nil, nil, ab38fmt.Errorf("ab38 row: trailing data")
		}
		if ab38Entry.Ab38Family == "" || ab38Entry.Ab38Variant == "" || ab38Entry.Ab38Direction == "" || ab38Entry.Ab38Field == "" || ab38Entry.Ab38Alternative == "" || ab38Entry.Ab38InstanceScope == "" || ab38Entry.Ab38DictionaryNamespace == "" || ab38Entry.Ab38DictionaryClass == "" || ab38Entry.Ab38RawAddress == nil || ab38Entry.Ab38BusFamily != "" || ab38Entry.Ab38NativeLabel == "" {
			return nil, nil, ab38fmt.Errorf("ab38 row: incomplete dictionary identity")
		}
		if ab38Entry.Ab38Grade != "Hypothesis" || ab38Entry.Ab38PhysicalQualification != "Unknown" || ab38Entry.Ab38ValueCodec != nil || ab38Entry.Ab38ValueWidthBits != nil || ab38Entry.Ab38ExecutableCapability || ab38Entry.Ab38AddressWidthBits != 16 || len(ab38Entry.Ab38BusFamilies) != 2 || ab38Entry.Ab38BusFamilies[0] != "eBUS" || ab38Entry.Ab38BusFamilies[1] != "PeBUS" || ab38Entry.Ab38RawForm != ab38fmt.Sprintf("%04X|%s", *ab38Entry.Ab38RawAddress, ab38Entry.Ab38NativeLabel) {
			return nil, nil, ab38fmt.Errorf("ab38 row: invalid candidate boundary or raw form")
		}
		for _, ab38BusFamily := range ab38Entry.Ab38BusFamilies {
			ab38Entry.Ab38BusFamily = ab38BusFamily
			ab38Encoded, ab38Err := ab38json.Marshal(ab38Entry.Ab38Key)
			if ab38Err != nil {
				return nil, nil, ab38Err
			}
			if _, ab38Exists := ab38Index[string(ab38Encoded)]; ab38Exists {
				return nil, nil, ab38fmt.Errorf("ab38 dictionary: duplicate key")
			}
			ab38Index[string(ab38Encoded)] = ab38bytes.Clone(ab38Raw)
		}
	}
	return ab38Rows, ab38Index, nil
}

func Ab38Entries() (ab38json.RawMessage, error) {
	if _, _, ab38Err := ab38Validate(ab38Data); ab38Err != nil {
		return nil, ab38Err
	}
	return ab38bytes.Clone(ab38Data), nil
}

func Ab38Lookup(ab38Key Ab38Key) (ab38json.RawMessage, bool, error) {
	_, ab38Index, ab38Err := ab38Validate(ab38Data)
	if ab38Err != nil {
		return nil, false, ab38Err
	}
	ab38Encoded, ab38Err := ab38json.Marshal(ab38Key)
	if ab38Err != nil {
		return nil, false, ab38Err
	}
	ab38Raw, ab38Found := ab38Index[string(ab38Encoded)]
	return ab38bytes.Clone(ab38Raw), ab38Found, nil
}

func ab38ValidateObject(ab38Raw []byte) error {
	ab38Decoder := ab38json.NewDecoder(ab38bytes.NewReader(ab38Raw))
	ab38Token, ab38Err := ab38Decoder.Token()
	if ab38Err != nil || ab38Token != ab38json.Delim('{') {
		return ab38fmt.Errorf("ab38 row: expected object")
	}
	ab38Seen := map[string]bool{}
	for ab38Decoder.More() {
		ab38Token, ab38Err = ab38Decoder.Token()
		if ab38Err != nil {
			return ab38Err
		}
		ab38Name, ab38OK := ab38Token.(string)
		if !ab38OK || ab38Seen[ab38Name] {
			return ab38fmt.Errorf("ab38 row: duplicate member")
		}
		ab38Seen[ab38Name] = true
		var ab38Value ab38json.RawMessage
		if ab38Err = ab38Decoder.Decode(&ab38Value); ab38Err != nil {
			return ab38Err
		}
		if ab38bytes.Equal(ab38Value, []byte("null")) && ab38Name != "secondary" && ab38Name != "operation" && ab38Name != "group" && ab38Name != "register" && ab38Name != "value_codec" && ab38Name != "value_width_bits" {
			return ab38fmt.Errorf("ab38 row: null required member")
		}
	}
	for _, ab38Name := range []string{"family", "variant", "primary", "direction", "field", "alternative", "instance_scope", "dictionary_namespace", "dictionary_class", "raw_address", "native_label", "raw_form", "bus_families", "grade", "physical_qualification", "value_codec", "executable_capability", "address_width_bits", "value_width_bits", "unknown", "unknown_address", "label_meaning", "native_contract"} {
		if !ab38Seen[ab38Name] {
			return ab38fmt.Errorf("ab38 row: missing member %s", ab38Name)
		}
	}
	return nil
}
