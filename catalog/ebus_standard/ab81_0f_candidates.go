package ebus_standard_catalog

import (
	_ "embed"
	ab81json "encoding/json"
)

//go:embed testdata/ab81_0f_candidates.json
var ab81CandidatesJSON []byte

type Ab81Key struct {
	Ab81Family    string `json:"family"`
	Ab81Variant   string `json:"variant"`
	Ab81Primary   byte   `json:"primary"`
	Ab81Secondary *byte  `json:"secondary"`
	Ab81Operation *byte  `json:"operation"`
	Ab81Direction string `json:"direction"`
	Ab81Field     string `json:"field"`
}

type Ab81Positive struct {
	Ab81Prefix           string `json:"prefix"`
	Ab81TypedDataPresent *bool  `json:"typed_data_present"`
	Ab81Grade            string `json:"grade"`
}

type Ab81Candidate struct {
	Ab81Key           Ab81Key       `json:"key"`
	Ab81MetadataGrade string        `json:"metadata_grade"`
	Ab81Grade         string        `json:"grade"`
	Ab81Body          string        `json:"body"`
	Ab81Offset        *int          `json:"offset"`
	Ab81Width         *int          `json:"width"`
	Ab81PhysicalNN    *int          `json:"physical_nn"`
	Ab81ACK           string        `json:"ack"`
	Ab81Support       string        `json:"support"`
	Ab81Effects       string        `json:"effects"`
	Ab81Executable    bool          `json:"executable"`
	Ab81Sendable      bool          `json:"sendable"`
	Ab81Raw           []byte        `json:"raw"`
	Ab81Positive      *Ab81Positive `json:"positive"`
	Ab81Limits        []string      `json:"limits"`
}

func Ab81Candidates() []Ab81Candidate {
	var ab81Entries []Ab81Candidate
	if ab81Err := ab81json.Unmarshal(ab81CandidatesJSON, &ab81Entries); ab81Err != nil {
		panic("ab81 embedded candidate metadata: " + ab81Err.Error())
	}
	return ab81Entries
}

func ab81OptionalEqual(ab81Left, ab81Right *byte) bool {
	return ab81Left == nil && ab81Right == nil || ab81Left != nil && ab81Right != nil && *ab81Left == *ab81Right
}

func Ab81Lookup(ab81Key Ab81Key) (Ab81Candidate, bool) {
	for _, ab81Entry := range Ab81Candidates() {
		ab81Stored := ab81Entry.Ab81Key
		if ab81Stored.Ab81Family == ab81Key.Ab81Family &&
			ab81Stored.Ab81Variant == ab81Key.Ab81Variant &&
			ab81Stored.Ab81Primary == ab81Key.Ab81Primary &&
			ab81OptionalEqual(ab81Stored.Ab81Secondary, ab81Key.Ab81Secondary) &&
			ab81OptionalEqual(ab81Stored.Ab81Operation, ab81Key.Ab81Operation) &&
			ab81Stored.Ab81Direction == ab81Key.Ab81Direction &&
			ab81Stored.Ab81Field == ab81Key.Ab81Field {
			return ab81Entry, true
		}
	}
	return Ab81Candidate{}, false
}
