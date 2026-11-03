package vaillant

import (
	ab17Bytes "bytes"
	_ "embed"
	ab17JSON "encoding/json"
	ab17Fmt "fmt"
	ab17Strings "strings"
)

//go:embed data/ab17_b506_rf_sequence.json
var ab17Catalog string

type Ab17Key struct {
	Ab17Family        string  `json:"family"`
	Ab17Variant       string  `json:"variant"`
	Ab17Primary       uint8   `json:"primary"`
	Ab17Secondary     *uint8  `json:"secondary"`
	Ab17Operation     *uint8  `json:"operation"`
	Ab17Direction     string  `json:"direction"`
	Ab17Field         string  `json:"field"`
	Ab17Alternative   string  `json:"alternative"`
	Ab17Group         *uint8  `json:"group"`
	Ab17InstanceScope string  `json:"instance_scope"`
	Ab17Register      *uint16 `json:"register"`
}

type ab17Row struct {
	Ab17ID            string            `json:"ab17_id"`
	Ab17Key           Ab17Key           `json:"key"`
	Ab17Grade         string            `json:"grade"`
	Ab17Qualification string            `json:"qualification"`
	Ab17RawForm       string            `json:"ab17_raw_form"`
	Ab17Unknown       map[string]string `json:"ab17_unknown"`
}

func ab17EntriesFrom(ab17Data string) ([]ab17JSON.RawMessage, error) {
	var ab17Rows []ab17JSON.RawMessage
	if ab17Error := ab17JSON.Unmarshal([]byte(ab17Data), &ab17Rows); ab17Error != nil {
		return nil, ab17Error
	}
	if ab17Rows == nil {
		return nil, ab17Fmt.Errorf("ab17: expected array")
	}
	ab17Seen := make(map[string]bool)
	for _, ab17Raw := range ab17Rows {
		var ab17Entry ab17Row
		if ab17Error := ab17JSON.Unmarshal(ab17Raw, &ab17Entry); ab17Error != nil {
			return nil, ab17Error
		}
		var ab17Object map[string]ab17JSON.RawMessage
		if ab17Error := ab17JSON.Unmarshal(ab17Raw, &ab17Object); ab17Error != nil {
			return nil, ab17Error
		}
		var ab17Axes map[string]ab17JSON.RawMessage
		if ab17Error := ab17JSON.Unmarshal(ab17Object["key"], &ab17Axes); ab17Error != nil {
			return nil, ab17Error
		}
		for _, ab17Axis := range []string{"family", "variant", "primary", "secondary", "operation", "direction", "field", "alternative", "group", "instance_scope", "register"} {
			if _, ab17Present := ab17Axes[ab17Axis]; !ab17Present {
				return nil, ab17Fmt.Errorf("ab17: missing axis %s", ab17Axis)
			}
			if ab17Axis != "secondary" && ab17Axis != "operation" && ab17Axis != "group" && ab17Axis != "register" && ab17Bytes.Equal(ab17Axes[ab17Axis], []byte("null")) {
				return nil, ab17Fmt.Errorf("ab17: null axis %s", ab17Axis)
			}
		}
		ab17Decoder := ab17JSON.NewDecoder(ab17Bytes.NewReader(ab17Object["key"]))
		ab17Decoder.DisallowUnknownFields()
		if ab17Error := ab17Decoder.Decode(&ab17Entry.Ab17Key); ab17Error != nil {
			return nil, ab17Error
		}
		if !ab17Strings.HasPrefix(ab17Entry.Ab17ID, "ab17-") || ab17Entry.Ab17Grade != "Hypothesis" || ab17Entry.Ab17Qualification != "Unknown" || ab17Entry.Ab17RawForm == "" || len(ab17Entry.Ab17Unknown) == 0 || ab17Entry.Ab17Key.Ab17Family == "" || ab17Entry.Ab17Key.Ab17Variant == "" || ab17Entry.Ab17Key.Ab17Direction == "" || ab17Entry.Ab17Key.Ab17Field == "" || ab17Entry.Ab17Key.Ab17Alternative == "" || ab17Entry.Ab17Key.Ab17InstanceScope == "" {
			return nil, ab17Fmt.Errorf("ab17: incomplete metadata")
		}
		ab17Identity, ab17Error := ab17JSON.Marshal(ab17Entry.Ab17Key)
		if ab17Error != nil {
			return nil, ab17Error
		}
		if ab17Seen[string(ab17Identity)] {
			return nil, ab17Fmt.Errorf("ab17: duplicate key")
		}
		ab17Seen[string(ab17Identity)] = true
	}
	return ab17Rows, nil
}

func Ab17Entries() (ab17JSON.RawMessage, error) {
	if _, ab17Error := ab17EntriesFrom(ab17Catalog); ab17Error != nil {
		return nil, ab17Error
	}
	return ab17JSON.RawMessage([]byte(ab17Catalog)), nil
}

func ab17LookupFrom(ab17Data string, ab17Key Ab17Key) (ab17JSON.RawMessage, bool, error) {
	ab17Rows, ab17Error := ab17EntriesFrom(ab17Data)
	if ab17Error != nil {
		return nil, false, ab17Error
	}
	ab17Wanted, ab17Error := ab17JSON.Marshal(ab17Key)
	if ab17Error != nil {
		return nil, false, ab17Error
	}
	for _, ab17Raw := range ab17Rows {
		var ab17Entry ab17Row
		if ab17Error := ab17JSON.Unmarshal(ab17Raw, &ab17Entry); ab17Error != nil {
			return nil, false, ab17Error
		}
		ab17Identity, ab17Error := ab17JSON.Marshal(ab17Entry.Ab17Key)
		if ab17Error != nil {
			return nil, false, ab17Error
		}
		if ab17Bytes.Equal(ab17Wanted, ab17Identity) {
			return append(ab17JSON.RawMessage(nil), ab17Raw...), true, nil
		}
	}
	return nil, false, nil
}

func Ab17Lookup(ab17Key Ab17Key) (ab17JSON.RawMessage, bool, error) {
	return ab17LookupFrom(ab17Catalog, ab17Key)
}
