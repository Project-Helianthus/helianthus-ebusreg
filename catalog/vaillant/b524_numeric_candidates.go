package vaillant

import (
	b524NumericCandidateBytes "bytes"
	_ "embed"
	b524NumericCandidateJSON "encoding/json"
	b524NumericCandidateFmt "fmt"
	b524NumericCandidateReflect "reflect"
)

type B524NumericCandidateContext struct {
	B524NumericCandidateFamily      string `json:"family,omitempty"`
	B524NumericCandidateSoftware    string `json:"software,omitempty"`
	B524NumericCandidateHardware    string `json:"hardware,omitempty"`
	B524NumericCandidateAPIVersion  string `json:"api_version,omitempty"`
	B524NumericCandidateAPIRevision string `json:"api_revision,omitempty"`
}

type B524NumericCandidateKey struct {
	B524NumericCandidateFamily        string  `json:"family,omitempty"`
	B524NumericCandidateSoftware      string  `json:"software,omitempty"`
	B524NumericCandidateHardware      string  `json:"hardware,omitempty"`
	B524NumericCandidateAPIVersion    string  `json:"api_version,omitempty"`
	B524NumericCandidateAPIRevision   string  `json:"api_revision,omitempty"`
	B524NumericCandidateVariant       string  `json:"variant"`
	B524NumericCandidateSet           string  `json:"candidate_set"`
	B524NumericCandidatePrimary       *uint16 `json:"primary"`
	B524NumericCandidateSecondary     *uint8  `json:"secondary"`
	B524NumericCandidateOperation     *uint8  `json:"operation"`
	B524NumericCandidateOpcode        *uint8  `json:"opcode"`
	B524NumericCandidateDirection     string  `json:"direction"`
	B524NumericCandidateField         string  `json:"field"`
	B524NumericCandidateAlternative   string  `json:"alternative"`
	B524NumericCandidateGroup         *uint8  `json:"group"`
	B524NumericCandidateInstanceScope string  `json:"instance_scope"`
	B524NumericCandidateInstance      *uint8  `json:"instance"`
	B524NumericCandidateRegister      *uint16 `json:"register"`
}

type b524NumericCandidateMetadata struct {
	B524NumericCandidateCodec         string   `json:"codec"`
	B524NumericCandidateWidth         *int     `json:"width_bytes"`
	B524NumericCandidateWidthGrade    string   `json:"width_grade"`
	B524NumericCandidateMinimum       *float64 `json:"minimum"`
	B524NumericCandidateMaximum       *float64 `json:"maximum"`
	B524NumericCandidateStep          *float64 `json:"step"`
	B524NumericCandidateGrade         string   `json:"grade"`
	B524NumericCandidateApplicability string   `json:"physical_applicability"`
	B524NumericCandidateQualification string   `json:"qualification"`
	B524NumericCandidateAdmitted      *bool    `json:"admitted"`
	B524NumericCandidatePresence      string   `json:"presence"`
	B524NumericCandidateFreshness     string   `json:"freshness"`
	B524NumericCandidateUnresolved    []string `json:"unresolved"`
	B524NumericCandidateRawForm       string   `json:"raw_form"`
	B524NumericCandidateRawFallback   string   `json:"raw_fallback"`
}

func b524NumericCandidateGrade(b524NumericCandidateValue string) bool {
	return b524NumericCandidateValue == "Proven" || b524NumericCandidateValue == "Hypothesis" || b524NumericCandidateValue == "Unknown"
}

type B524NumericCandidateView struct {
	b524NumericCandidateRows []b524NumericCandidateJSON.RawMessage
	b524NumericCandidateKeys []B524NumericCandidateKey
}

//go:embed data/ab56_b524_numeric_candidates.json
var b524NumericCandidateData string

func b524NumericCandidateIdentity(b524NumericCandidateKey B524NumericCandidateKey) B524NumericCandidateContext {
	return B524NumericCandidateContext{b524NumericCandidateKey.B524NumericCandidateFamily, b524NumericCandidateKey.B524NumericCandidateSoftware, b524NumericCandidateKey.B524NumericCandidateHardware, b524NumericCandidateKey.B524NumericCandidateAPIVersion, b524NumericCandidateKey.B524NumericCandidateAPIRevision}
}

func b524NumericCandidateComplete(b524NumericCandidateContext B524NumericCandidateContext) bool {
	return b524NumericCandidateContext.B524NumericCandidateFamily != "" && b524NumericCandidateContext.B524NumericCandidateSoftware != "" && b524NumericCandidateContext.B524NumericCandidateHardware != "" && b524NumericCandidateContext.B524NumericCandidateAPIVersion != "" && b524NumericCandidateContext.B524NumericCandidateAPIRevision != ""
}

func b524NumericCandidateUniqueJSON(b524NumericCandidateSource []byte) error {
	b524NumericCandidateDecoder := b524NumericCandidateJSON.NewDecoder(b524NumericCandidateBytes.NewReader(b524NumericCandidateSource))
	var b524NumericCandidateVisit func() error
	b524NumericCandidateVisit = func() error {
		b524NumericCandidateToken, b524NumericCandidateError := b524NumericCandidateDecoder.Token()
		if b524NumericCandidateError != nil {
			return b524NumericCandidateError
		}
		b524NumericCandidateDelimiter, b524NumericCandidateOK := b524NumericCandidateToken.(b524NumericCandidateJSON.Delim)
		if !b524NumericCandidateOK {
			return nil
		}
		if b524NumericCandidateDelimiter != '{' && b524NumericCandidateDelimiter != '[' {
			return b524NumericCandidateFmt.Errorf("ab56: unexpected delimiter")
		}
		b524NumericCandidateSeen := make(map[string]bool)
		for b524NumericCandidateDecoder.More() {
			if b524NumericCandidateDelimiter == '{' {
				b524NumericCandidateNameToken, b524NumericCandidateError := b524NumericCandidateDecoder.Token()
				if b524NumericCandidateError != nil {
					return b524NumericCandidateError
				}
				b524NumericCandidateName, b524NumericCandidateOK := b524NumericCandidateNameToken.(string)
				if !b524NumericCandidateOK || b524NumericCandidateSeen[b524NumericCandidateName] {
					return b524NumericCandidateFmt.Errorf("ab56: duplicate or invalid JSON member")
				}
				b524NumericCandidateSeen[b524NumericCandidateName] = true
			}
			if b524NumericCandidateError := b524NumericCandidateVisit(); b524NumericCandidateError != nil {
				return b524NumericCandidateError
			}
		}
		_, b524NumericCandidateError = b524NumericCandidateDecoder.Token()
		return b524NumericCandidateError
	}
	return b524NumericCandidateVisit()
}

func b524NumericCandidateValidate(b524NumericCandidateSource []byte) ([]b524NumericCandidateJSON.RawMessage, []B524NumericCandidateKey, error) {
	if b524NumericCandidateError := b524NumericCandidateUniqueJSON(b524NumericCandidateSource); b524NumericCandidateError != nil {
		return nil, nil, b524NumericCandidateError
	}
	var b524NumericCandidateRows []b524NumericCandidateJSON.RawMessage
	if b524NumericCandidateError := b524NumericCandidateJSON.Unmarshal(b524NumericCandidateSource, &b524NumericCandidateRows); b524NumericCandidateError != nil {
		return nil, nil, b524NumericCandidateError
	}
	if b524NumericCandidateRows == nil {
		return nil, nil, b524NumericCandidateFmt.Errorf("ab56: expected a tuple array")
	}
	b524NumericCandidateKeys := make([]B524NumericCandidateKey, 0, len(b524NumericCandidateRows))
	for b524NumericCandidateIndex, b524NumericCandidateRow := range b524NumericCandidateRows {
		var b524NumericCandidateObject map[string]b524NumericCandidateJSON.RawMessage
		if b524NumericCandidateError := b524NumericCandidateJSON.Unmarshal(b524NumericCandidateRow, &b524NumericCandidateObject); b524NumericCandidateError != nil || b524NumericCandidateObject == nil {
			return nil, nil, b524NumericCandidateFmt.Errorf("ab56: row %d must be an object", b524NumericCandidateIndex)
		}
		var b524NumericCandidateMetadataValue b524NumericCandidateMetadata
		if b524NumericCandidateError := b524NumericCandidateJSON.Unmarshal(b524NumericCandidateRow, &b524NumericCandidateMetadataValue); b524NumericCandidateError != nil {
			return nil, nil, b524NumericCandidateError
		}
		if b524NumericCandidateMetadataValue.B524NumericCandidateCodec == "" ||
			b524NumericCandidateObject["width_bytes"] == nil ||
			!b524NumericCandidateGrade(b524NumericCandidateMetadataValue.B524NumericCandidateWidthGrade) ||
			(b524NumericCandidateMetadataValue.B524NumericCandidateWidth != nil && *b524NumericCandidateMetadataValue.B524NumericCandidateWidth < 0) ||
			b524NumericCandidateMetadataValue.B524NumericCandidateMinimum == nil ||
			b524NumericCandidateMetadataValue.B524NumericCandidateMaximum == nil ||
			b524NumericCandidateMetadataValue.B524NumericCandidateStep == nil ||
			*b524NumericCandidateMetadataValue.B524NumericCandidateStep <= 0 ||
			*b524NumericCandidateMetadataValue.B524NumericCandidateMinimum > *b524NumericCandidateMetadataValue.B524NumericCandidateMaximum ||
			!b524NumericCandidateGrade(b524NumericCandidateMetadataValue.B524NumericCandidateGrade) ||
			b524NumericCandidateMetadataValue.B524NumericCandidateApplicability != "Unknown" ||
			b524NumericCandidateMetadataValue.B524NumericCandidateQualification != "Unknown" ||
			b524NumericCandidateMetadataValue.B524NumericCandidateAdmitted == nil ||
			*b524NumericCandidateMetadataValue.B524NumericCandidateAdmitted ||
			b524NumericCandidateMetadataValue.B524NumericCandidatePresence != "Unknown" ||
			b524NumericCandidateMetadataValue.B524NumericCandidateFreshness != "Unknown" ||
			len(b524NumericCandidateMetadataValue.B524NumericCandidateUnresolved) == 0 ||
			b524NumericCandidateMetadataValue.B524NumericCandidateRawForm == "" ||
			b524NumericCandidateMetadataValue.B524NumericCandidateRawFallback == "" ||
			b524NumericCandidateObject["writable"] != nil {
			return nil, nil, b524NumericCandidateFmt.Errorf("ab56: invalid or promoted metadata at row %d", b524NumericCandidateIndex)
		}
		var b524NumericCandidateKey B524NumericCandidateKey
		b524NumericCandidateDecoder := b524NumericCandidateJSON.NewDecoder(b524NumericCandidateBytes.NewReader(b524NumericCandidateObject["key"]))
		b524NumericCandidateDecoder.DisallowUnknownFields()
		if b524NumericCandidateError := b524NumericCandidateDecoder.Decode(&b524NumericCandidateKey); b524NumericCandidateError != nil {
			return nil, nil, b524NumericCandidateError
		}
		if b524NumericCandidateIdentity(b524NumericCandidateKey) != (B524NumericCandidateContext{}) || b524NumericCandidateKey.B524NumericCandidateVariant == "" || b524NumericCandidateKey.B524NumericCandidateSet == "" || b524NumericCandidateKey.B524NumericCandidatePrimary == nil || b524NumericCandidateKey.B524NumericCandidateDirection == "" || b524NumericCandidateKey.B524NumericCandidateField == "" || b524NumericCandidateKey.B524NumericCandidateAlternative == "" || b524NumericCandidateKey.B524NumericCandidateInstanceScope == "" {
			return nil, nil, b524NumericCandidateFmt.Errorf("ab56: row %d has an incomplete or identity-qualified tuple key", b524NumericCandidateIndex)
		}
		for _, b524NumericCandidatePrevious := range b524NumericCandidateKeys {
			if b524NumericCandidateReflect.DeepEqual(b524NumericCandidatePrevious, b524NumericCandidateKey) {
				return nil, nil, b524NumericCandidateFmt.Errorf("ab56: duplicate tuple key at row %d", b524NumericCandidateIndex)
			}
		}
		b524NumericCandidateKeys = append(b524NumericCandidateKeys, b524NumericCandidateKey)
	}
	return b524NumericCandidateRows, b524NumericCandidateKeys, nil
}

func B524NumericCandidateEntries() (b524NumericCandidateJSON.RawMessage, error) {
	if _, _, b524NumericCandidateError := b524NumericCandidateValidate([]byte(b524NumericCandidateData)); b524NumericCandidateError != nil {
		return nil, b524NumericCandidateError
	}
	return b524NumericCandidateJSON.RawMessage([]byte(b524NumericCandidateData)), nil
}

func b524NumericCandidateNewView(b524NumericCandidateContext B524NumericCandidateContext, b524NumericCandidateSource []byte) (B524NumericCandidateView, error) {
	if !b524NumericCandidateComplete(b524NumericCandidateContext) {
		return B524NumericCandidateView{}, b524NumericCandidateFmt.Errorf("ab56: all caller identity strings are required")
	}
	b524NumericCandidateRows, b524NumericCandidateKeys, b524NumericCandidateError := b524NumericCandidateValidate(b524NumericCandidateSource)
	if b524NumericCandidateError != nil {
		return B524NumericCandidateView{}, b524NumericCandidateError
	}
	b524NumericCandidateView := B524NumericCandidateView{b524NumericCandidateRows: b524NumericCandidateRows, b524NumericCandidateKeys: b524NumericCandidateKeys}
	for b524NumericCandidateIndex, b524NumericCandidateKey := range b524NumericCandidateKeys {
		b524NumericCandidateKey.B524NumericCandidateFamily = b524NumericCandidateContext.B524NumericCandidateFamily
		b524NumericCandidateKey.B524NumericCandidateSoftware = b524NumericCandidateContext.B524NumericCandidateSoftware
		b524NumericCandidateKey.B524NumericCandidateHardware = b524NumericCandidateContext.B524NumericCandidateHardware
		b524NumericCandidateKey.B524NumericCandidateAPIVersion = b524NumericCandidateContext.B524NumericCandidateAPIVersion
		b524NumericCandidateKey.B524NumericCandidateAPIRevision = b524NumericCandidateContext.B524NumericCandidateAPIRevision
		b524NumericCandidateView.b524NumericCandidateKeys[b524NumericCandidateIndex] = b524NumericCandidateKey
	}
	return b524NumericCandidateView, nil
}

func B524NumericCandidateNewView(b524NumericCandidateContext B524NumericCandidateContext) (B524NumericCandidateView, error) {
	return b524NumericCandidateNewView(b524NumericCandidateContext, []byte(b524NumericCandidateData))
}

func (b524NumericCandidateView B524NumericCandidateView) B524NumericCandidateLookup(b524NumericCandidateKey B524NumericCandidateKey) (b524NumericCandidateJSON.RawMessage, bool, error) {
	if !b524NumericCandidateComplete(b524NumericCandidateIdentity(b524NumericCandidateKey)) {
		return nil, false, nil
	}
	for b524NumericCandidateIndex, b524NumericCandidateCandidate := range b524NumericCandidateView.b524NumericCandidateKeys {
		if b524NumericCandidateReflect.DeepEqual(b524NumericCandidateCandidate, b524NumericCandidateKey) {
			return append(b524NumericCandidateJSON.RawMessage(nil), b524NumericCandidateView.b524NumericCandidateRows[b524NumericCandidateIndex]...), true, nil
		}
	}
	return nil, false, nil
}
