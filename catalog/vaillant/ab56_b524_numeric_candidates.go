package vaillant

import (
	ab56Bytes "bytes"
	_ "embed"
	ab56JSON "encoding/json"
	ab56Fmt "fmt"
	ab56Reflect "reflect"
)

type Ab56Context struct {
	Ab56Family      string `json:"family,omitempty"`
	Ab56Software    string `json:"software,omitempty"`
	Ab56Hardware    string `json:"hardware,omitempty"`
	Ab56APIVersion  string `json:"api_version,omitempty"`
	Ab56APIRevision string `json:"api_revision,omitempty"`
}

type Ab56Key struct {
	Ab56Family        string  `json:"family,omitempty"`
	Ab56Software      string  `json:"software,omitempty"`
	Ab56Hardware      string  `json:"hardware,omitempty"`
	Ab56APIVersion    string  `json:"api_version,omitempty"`
	Ab56APIRevision   string  `json:"api_revision,omitempty"`
	Ab56Variant       string  `json:"variant"`
	Ab56CandidateSet  string  `json:"candidate_set"`
	Ab56Primary       *uint16 `json:"primary"`
	Ab56Secondary     *uint8  `json:"secondary"`
	Ab56Operation     *uint8  `json:"operation"`
	Ab56Opcode        *uint8  `json:"opcode"`
	Ab56Direction     string  `json:"direction"`
	Ab56Field         string  `json:"field"`
	Ab56Alternative   string  `json:"alternative"`
	Ab56Group         *uint8  `json:"group"`
	Ab56InstanceScope string  `json:"instance_scope"`
	Ab56Instance      *uint8  `json:"instance"`
	Ab56Register      *uint16 `json:"register"`
}

type ab56Metadata struct {
	Ab56Codec         string   `json:"codec"`
	Ab56Width         *int     `json:"width_bytes"`
	Ab56WidthGrade    string   `json:"width_grade"`
	Ab56Minimum       *float64 `json:"minimum"`
	Ab56Maximum       *float64 `json:"maximum"`
	Ab56Step          *float64 `json:"step"`
	Ab56Grade         string   `json:"grade"`
	Ab56Applicability string   `json:"physical_applicability"`
	Ab56Qualification string   `json:"qualification"`
	Ab56Admitted      *bool    `json:"admitted"`
	Ab56Presence      string   `json:"presence"`
	Ab56Freshness     string   `json:"freshness"`
	Ab56Unresolved    []string `json:"unresolved"`
	Ab56RawForm       string   `json:"raw_form"`
	Ab56RawFallback   string   `json:"raw_fallback"`
}

func ab56Grade(ab56Value string) bool {
	return ab56Value == "Proven" || ab56Value == "Hypothesis" || ab56Value == "Unknown"
}

type Ab56View struct {
	ab56Rows []ab56JSON.RawMessage
	ab56Keys []Ab56Key
}

//go:embed data/ab56_b524_numeric_candidates.json
var ab56Data string

func ab56Identity(ab56Key Ab56Key) Ab56Context {
	return Ab56Context{ab56Key.Ab56Family, ab56Key.Ab56Software, ab56Key.Ab56Hardware, ab56Key.Ab56APIVersion, ab56Key.Ab56APIRevision}
}

func ab56Complete(ab56Context Ab56Context) bool {
	return ab56Context.Ab56Family != "" && ab56Context.Ab56Software != "" && ab56Context.Ab56Hardware != "" && ab56Context.Ab56APIVersion != "" && ab56Context.Ab56APIRevision != ""
}

func ab56UniqueJSON(ab56Source []byte) error {
	ab56Decoder := ab56JSON.NewDecoder(ab56Bytes.NewReader(ab56Source))
	var ab56Visit func() error
	ab56Visit = func() error {
		ab56Token, ab56Error := ab56Decoder.Token()
		if ab56Error != nil {
			return ab56Error
		}
		ab56Delimiter, ab56OK := ab56Token.(ab56JSON.Delim)
		if !ab56OK {
			return nil
		}
		if ab56Delimiter != '{' && ab56Delimiter != '[' {
			return ab56Fmt.Errorf("ab56: unexpected delimiter")
		}
		ab56Seen := make(map[string]bool)
		for ab56Decoder.More() {
			if ab56Delimiter == '{' {
				ab56NameToken, ab56Error := ab56Decoder.Token()
				if ab56Error != nil {
					return ab56Error
				}
				ab56Name, ab56OK := ab56NameToken.(string)
				if !ab56OK || ab56Seen[ab56Name] {
					return ab56Fmt.Errorf("ab56: duplicate or invalid JSON member")
				}
				ab56Seen[ab56Name] = true
			}
			if ab56Error := ab56Visit(); ab56Error != nil {
				return ab56Error
			}
		}
		_, ab56Error = ab56Decoder.Token()
		return ab56Error
	}
	return ab56Visit()
}

func ab56Validate(ab56Source []byte) ([]ab56JSON.RawMessage, []Ab56Key, error) {
	if ab56Error := ab56UniqueJSON(ab56Source); ab56Error != nil {
		return nil, nil, ab56Error
	}
	var ab56Rows []ab56JSON.RawMessage
	if ab56Error := ab56JSON.Unmarshal(ab56Source, &ab56Rows); ab56Error != nil {
		return nil, nil, ab56Error
	}
	if ab56Rows == nil {
		return nil, nil, ab56Fmt.Errorf("ab56: expected a tuple array")
	}
	ab56Keys := make([]Ab56Key, 0, len(ab56Rows))
	for ab56Index, ab56Row := range ab56Rows {
		var ab56Object map[string]ab56JSON.RawMessage
		if ab56Error := ab56JSON.Unmarshal(ab56Row, &ab56Object); ab56Error != nil || ab56Object == nil {
			return nil, nil, ab56Fmt.Errorf("ab56: row %d must be an object", ab56Index)
		}
		var ab56MetadataValue ab56Metadata
		if ab56Error := ab56JSON.Unmarshal(ab56Row, &ab56MetadataValue); ab56Error != nil {
			return nil, nil, ab56Error
		}
		if ab56MetadataValue.Ab56Codec == "" ||
			ab56Object["width_bytes"] == nil ||
			!ab56Grade(ab56MetadataValue.Ab56WidthGrade) ||
			(ab56MetadataValue.Ab56Width != nil && *ab56MetadataValue.Ab56Width < 0) ||
			ab56MetadataValue.Ab56Minimum == nil ||
			ab56MetadataValue.Ab56Maximum == nil ||
			ab56MetadataValue.Ab56Step == nil ||
			*ab56MetadataValue.Ab56Step <= 0 ||
			*ab56MetadataValue.Ab56Minimum > *ab56MetadataValue.Ab56Maximum ||
			!ab56Grade(ab56MetadataValue.Ab56Grade) ||
			ab56MetadataValue.Ab56Applicability != "Unknown" ||
			ab56MetadataValue.Ab56Qualification != "Unknown" ||
			ab56MetadataValue.Ab56Admitted == nil ||
			*ab56MetadataValue.Ab56Admitted ||
			ab56MetadataValue.Ab56Presence != "Unknown" ||
			ab56MetadataValue.Ab56Freshness != "Unknown" ||
			len(ab56MetadataValue.Ab56Unresolved) == 0 ||
			ab56MetadataValue.Ab56RawForm == "" ||
			ab56MetadataValue.Ab56RawFallback == "" ||
			ab56Object["writable"] != nil {
			return nil, nil, ab56Fmt.Errorf("ab56: invalid or promoted metadata at row %d", ab56Index)
		}
		var ab56Key Ab56Key
		ab56Decoder := ab56JSON.NewDecoder(ab56Bytes.NewReader(ab56Object["key"]))
		ab56Decoder.DisallowUnknownFields()
		if ab56Error := ab56Decoder.Decode(&ab56Key); ab56Error != nil {
			return nil, nil, ab56Error
		}
		if ab56Identity(ab56Key) != (Ab56Context{}) || ab56Key.Ab56Variant == "" || ab56Key.Ab56CandidateSet == "" || ab56Key.Ab56Primary == nil || ab56Key.Ab56Direction == "" || ab56Key.Ab56Field == "" || ab56Key.Ab56Alternative == "" || ab56Key.Ab56InstanceScope == "" {
			return nil, nil, ab56Fmt.Errorf("ab56: row %d has an incomplete or identity-qualified tuple key", ab56Index)
		}
		for _, ab56Previous := range ab56Keys {
			if ab56Reflect.DeepEqual(ab56Previous, ab56Key) {
				return nil, nil, ab56Fmt.Errorf("ab56: duplicate tuple key at row %d", ab56Index)
			}
		}
		ab56Keys = append(ab56Keys, ab56Key)
	}
	return ab56Rows, ab56Keys, nil
}

func Ab56Entries() (ab56JSON.RawMessage, error) {
	if _, _, ab56Error := ab56Validate([]byte(ab56Data)); ab56Error != nil {
		return nil, ab56Error
	}
	return ab56JSON.RawMessage([]byte(ab56Data)), nil
}

func ab56NewView(ab56Context Ab56Context, ab56Source []byte) (Ab56View, error) {
	if !ab56Complete(ab56Context) {
		return Ab56View{}, ab56Fmt.Errorf("ab56: all caller identity strings are required")
	}
	ab56Rows, ab56Keys, ab56Error := ab56Validate(ab56Source)
	if ab56Error != nil {
		return Ab56View{}, ab56Error
	}
	ab56View := Ab56View{ab56Rows: ab56Rows, ab56Keys: ab56Keys}
	for ab56Index, ab56Key := range ab56Keys {
		ab56Key.Ab56Family = ab56Context.Ab56Family
		ab56Key.Ab56Software = ab56Context.Ab56Software
		ab56Key.Ab56Hardware = ab56Context.Ab56Hardware
		ab56Key.Ab56APIVersion = ab56Context.Ab56APIVersion
		ab56Key.Ab56APIRevision = ab56Context.Ab56APIRevision
		ab56View.ab56Keys[ab56Index] = ab56Key
	}
	return ab56View, nil
}

func Ab56NewView(ab56Context Ab56Context) (Ab56View, error) {
	return ab56NewView(ab56Context, []byte(ab56Data))
}

func (ab56View Ab56View) Ab56Lookup(ab56Key Ab56Key) (ab56JSON.RawMessage, bool, error) {
	if !ab56Complete(ab56Identity(ab56Key)) {
		return nil, false, nil
	}
	for ab56Index, ab56Candidate := range ab56View.ab56Keys {
		if ab56Reflect.DeepEqual(ab56Candidate, ab56Key) {
			return append(ab56JSON.RawMessage(nil), ab56View.ab56Rows[ab56Index]...), true, nil
		}
	}
	return nil, false, nil
}
