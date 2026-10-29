package vaillant

import (
	_ "embed"
	ab67json "encoding/json"
)

//go:embed data/ab67_b511_candidates.json
var ab67B511Data []byte

type Ab67B511Key struct {
	Ab67Family    string `json:"family"`
	Ab67Variant   string `json:"variant"`
	Ab67Primary   *uint8 `json:"primary"`
	Ab67Secondary *uint8 `json:"secondary"`
	Ab67Operation *uint8 `json:"operation"`
	Ab67Direction string `json:"direction"`
	Ab67Field     string `json:"field"`
}

type Ab67B511Entry struct {
	Ab67Key              Ab67B511Key      `json:"key"`
	Ab67Position         *int             `json:"position"`
	Ab67WidthBits        *int             `json:"width_bits"`
	Ab67Extent           *int             `json:"extent"`
	Ab67Grade            Ab05Grade        `json:"grade"`
	Ab67Default          *int             `json:"default"`
	Ab67Values           map[uint8]string `json:"values"`
	Ab67Bits             map[uint8]string `json:"bits"`
	Ab67ReservedBits     []uint8          `json:"reserved_bits"`
	Ab67Membership       []uint8          `json:"membership"`
	Ab67RawDomain        []int            `json:"raw_domain"`
	Ab67Scale            *float64         `json:"scale"`
	Ab67Unit             *string          `json:"unit"`
	Ab67ClockValidation  Ab05Grade        `json:"clock_validation"`
	Ab67PhysicalPrefix   Ab05Grade        `json:"physical_prefix"`
	Ab67InstalledSupport Ab05Grade        `json:"installed_support"`
	Ab67RawPolicy        string           `json:"raw_policy"`
}

type Ab67B511Declarations struct {
	Ab67MetadataPresence      Ab05Grade       `json:"metadata_presence"`
	Ab67PhysicalQualification Ab05Grade       `json:"physical_qualification"`
	Ab67NativeContract        string          `json:"native_contract"`
	Ab67Limitations           []string        `json:"limitations"`
	Ab67Entries               []Ab67B511Entry `json:"entries"`
}

func Ab67B511Catalog() (Ab67B511Declarations, error) {
	var ab67Catalog Ab67B511Declarations
	ab67Error := ab67json.Unmarshal(ab67B511Data, &ab67Catalog)
	return ab67Catalog, ab67Error
}

func ab67EqualByte(ab67Left, ab67Right *uint8) bool {
	if ab67Left == nil || ab67Right == nil {
		return ab67Left == nil && ab67Right == nil
	}
	return *ab67Left == *ab67Right
}

func Ab67B511Lookup(ab67Key Ab67B511Key) (Ab67B511Entry, bool, error) {
	ab67Catalog, ab67Error := Ab67B511Catalog()
	if ab67Error != nil {
		return Ab67B511Entry{}, false, ab67Error
	}
	for _, ab67Entry := range ab67Catalog.Ab67Entries {
		ab67Candidate := ab67Entry.Ab67Key
		if ab67Candidate.Ab67Family == ab67Key.Ab67Family &&
			ab67Candidate.Ab67Variant == ab67Key.Ab67Variant &&
			ab67EqualByte(ab67Candidate.Ab67Primary, ab67Key.Ab67Primary) &&
			ab67EqualByte(ab67Candidate.Ab67Secondary, ab67Key.Ab67Secondary) &&
			ab67EqualByte(ab67Candidate.Ab67Operation, ab67Key.Ab67Operation) &&
			ab67Candidate.Ab67Direction == ab67Key.Ab67Direction &&
			ab67Candidate.Ab67Field == ab67Key.Ab67Field {
			return ab67Entry, true, nil
		}
	}
	return Ab67B511Entry{}, false, nil
}
