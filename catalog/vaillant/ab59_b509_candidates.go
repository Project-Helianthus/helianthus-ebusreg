package vaillant

import (
	_ "embed"
	ab59json "encoding/json"
	ab59reflect "reflect"
)

//go:embed data/ab59_b509_candidates.json
var ab59B509CandidatesJSON string

type Ab59Key struct {
	Ab59Family    string  `json:"family"`
	Ab59Variant   string  `json:"variant"`
	Ab59Primary   byte    `json:"primary"`
	Ab59Secondary *byte   `json:"secondary"`
	Ab59Operation *byte   `json:"operation"`
	Ab59Direction *string `json:"direction"`
	Ab59Field     string  `json:"field"`
}

type Ab59Entry struct {
	Ab59Key           Ab59Key           `json:"key"`
	Ab59Position      *int              `json:"position"`
	Ab59WidthBits     *int              `json:"width_bits"`
	Ab59Signedness    string            `json:"signedness"`
	Ab59ByteOrder     string            `json:"byte_order"`
	Ab59Conversion    string            `json:"conversion"`
	Ab59Divisor       *int              `json:"divisor"`
	Ab59Unit          *string           `json:"unit"`
	Ab59Default       *int64            `json:"default"`
	Ab59Lower         *int64            `json:"lower"`
	Ab59Upper         *int64            `json:"upper"`
	Ab59Mappings      map[string]string `json:"mappings"`
	Ab59OtherFeedback *string           `json:"other_feedback"`
	Ab59State         string            `json:"state"`
	Ab59Grade         Ab05Grade         `json:"grade"`
	Ab59Notes         []string          `json:"notes"`
}

type Ab59Catalog struct {
	Ab59MetadataPresence Ab05Grade   `json:"metadata_presence"`
	Ab59WireAdaptation   Ab05Grade   `json:"wire_adaptation"`
	Ab59InstalledSupport Ab05Grade   `json:"installed_support"`
	Ab59Invocation       Ab05Grade   `json:"invocation"`
	Ab59Limitations      []string    `json:"limitations"`
	Ab59Entries          []Ab59Entry `json:"entries"`
}

func Ab59B509Candidates() (Ab59Catalog, error) {
	var ab59Catalog Ab59Catalog
	ab59Error := ab59json.Unmarshal([]byte(ab59B509CandidatesJSON), &ab59Catalog)
	return ab59Catalog, ab59Error
}

func Ab59Lookup(ab59Key Ab59Key) (Ab59Entry, bool, error) {
	ab59Catalog, ab59Error := Ab59B509Candidates()
	if ab59Error != nil {
		return Ab59Entry{}, false, ab59Error
	}
	ab59Entry, ab59Found := ab59Find(ab59Catalog.Ab59Entries, ab59Key)
	return ab59Entry, ab59Found, nil
}

func ab59Find(ab59Entries []Ab59Entry, ab59Key Ab59Key) (Ab59Entry, bool) {
	for _, ab59Entry := range ab59Entries {
		if ab59reflect.DeepEqual(ab59Entry.Ab59Key, ab59Key) {
			return ab59Entry, true
		}
	}
	return Ab59Entry{}, false
}
