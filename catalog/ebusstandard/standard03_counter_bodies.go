package ebusstandard

import (
	_ "embed"
	standard03CounterBodyJSON "encoding/json"
)

type Standard03CounterBodyKey struct {
	Standard03CounterBodyFamily    string  `json:"family"`
	Standard03CounterBodyVariant   string  `json:"variant"`
	Standard03CounterBodyPrimary   *uint16 `json:"primary"`
	Standard03CounterBodySecondary *uint8  `json:"secondary"`
	Standard03CounterBodyOperation *uint8  `json:"operation"`
	Standard03CounterBodyDirection string  `json:"direction"`
	Standard03CounterBodyField     string  `json:"field"`
	Standard03CounterBodyOpcode    *uint8  `json:"opcode"`
	Standard03CounterBodyGroup     *uint8  `json:"group"`
	Standard03CounterBodyInstance  *uint8  `json:"instance"`
	Standard03CounterBodyRegister  *uint16 `json:"register"`
}

type Standard03CounterBodyField struct {
	Standard03CounterBodyIdentity          string            `json:"identity"`
	Standard03CounterBodyLogicalPosition   *int              `json:"logical_position"`
	Standard03CounterBodyNativePosition    *int              `json:"native_position"`
	Standard03CounterBodyWidthBits         *int              `json:"width_bits"`
	Standard03CounterBodySignedness        *string           `json:"signedness"`
	Standard03CounterBodyByteOrder         *string           `json:"byte_order"`
	Standard03CounterBodyEncoding          *string           `json:"encoding"`
	Standard03CounterBodyConversion        *string           `json:"conversion"`
	Standard03CounterBodyUnit              *string           `json:"unit"`
	Standard03CounterBodyDefault           *int              `json:"default"`
	Standard03CounterBodyMinimum           *int              `json:"minimum"`
	Standard03CounterBodyMaximum           *int              `json:"maximum"`
	Standard03CounterBodyEnumeration       map[string]string `json:"enumeration"`
	Standard03CounterBodyReplacementValues []int             `json:"replacement_values"`
	Standard03CounterBodyReservedValues    []int             `json:"reserved_values"`
	Standard03CounterBodyClaimGrade        string            `json:"claim_grade"`
	Standard03CounterBodyRawFallback       bool              `json:"raw_fallback"`
}

type Standard03CounterBody struct {
	Standard03CounterBodyKey               Standard03CounterBodyKey     `json:"key"`
	Standard03CounterBodyState             string                       `json:"body_state"`
	Standard03CounterBodyClaimGrade        string                       `json:"claim_grade"`
	Standard03CounterBodyPhysicalDirection string                       `json:"physical_direction"`
	Standard03CounterBodyCallable          bool                         `json:"callable"`
	Standard03CounterBodyRawFallback       bool                         `json:"raw_fallback"`
	Standard03CounterBodyFields            []Standard03CounterBodyField `json:"fields"`
}

//go:embed testdata/ab83_03_counter_bodies.json
var standard03CounterBodyCounterJSON string

// Standard03CounterBodies returns detached finite supplemental declarations, not invocation support.
func Standard03CounterBodies() []Standard03CounterBody {
	var standard03CounterBodyBodies []Standard03CounterBody
	if standard03CounterBodyErr := standard03CounterBodyJSON.Unmarshal([]byte(standard03CounterBodyCounterJSON), &standard03CounterBodyBodies); standard03CounterBodyErr != nil {
		panic("ab83 counter metadata: " + standard03CounterBodyErr.Error())
	}
	return standard03CounterBodyBodies
}

func standard03CounterBodyOptionalEqual[standard03CounterBodyT comparable](standard03CounterBodyLeft, standard03CounterBodyRight *standard03CounterBodyT) bool {
	if standard03CounterBodyLeft == nil || standard03CounterBodyRight == nil {
		return standard03CounterBodyLeft == nil && standard03CounterBodyRight == nil
	}
	return *standard03CounterBodyLeft == *standard03CounterBodyRight
}

func standard03CounterBodyEqualKey(standard03CounterBodyLeft, standard03CounterBodyRight Standard03CounterBodyKey) bool {
	return standard03CounterBodyLeft.Standard03CounterBodyFamily == standard03CounterBodyRight.Standard03CounterBodyFamily &&
		standard03CounterBodyLeft.Standard03CounterBodyVariant == standard03CounterBodyRight.Standard03CounterBodyVariant &&
		standard03CounterBodyOptionalEqual(standard03CounterBodyLeft.Standard03CounterBodyPrimary, standard03CounterBodyRight.Standard03CounterBodyPrimary) &&
		standard03CounterBodyOptionalEqual(standard03CounterBodyLeft.Standard03CounterBodySecondary, standard03CounterBodyRight.Standard03CounterBodySecondary) &&
		standard03CounterBodyOptionalEqual(standard03CounterBodyLeft.Standard03CounterBodyOperation, standard03CounterBodyRight.Standard03CounterBodyOperation) &&
		standard03CounterBodyLeft.Standard03CounterBodyDirection == standard03CounterBodyRight.Standard03CounterBodyDirection &&
		standard03CounterBodyLeft.Standard03CounterBodyField == standard03CounterBodyRight.Standard03CounterBodyField &&
		standard03CounterBodyOptionalEqual(standard03CounterBodyLeft.Standard03CounterBodyOpcode, standard03CounterBodyRight.Standard03CounterBodyOpcode) &&
		standard03CounterBodyOptionalEqual(standard03CounterBodyLeft.Standard03CounterBodyGroup, standard03CounterBodyRight.Standard03CounterBodyGroup) &&
		standard03CounterBodyOptionalEqual(standard03CounterBodyLeft.Standard03CounterBodyInstance, standard03CounterBodyRight.Standard03CounterBodyInstance) &&
		standard03CounterBodyOptionalEqual(standard03CounterBodyLeft.Standard03CounterBodyRegister, standard03CounterBodyRight.Standard03CounterBodyRegister)
}

// Standard03CounterLookupBody requires every axis to match, including an empty field axis.
func Standard03CounterLookupBody(standard03CounterBodyKey Standard03CounterBodyKey) (Standard03CounterBody, bool) {
	for _, standard03CounterBodyBody := range Standard03CounterBodies() {
		if standard03CounterBodyEqualKey(standard03CounterBodyBody.Standard03CounterBodyKey, standard03CounterBodyKey) {
			return standard03CounterBodyBody, true
		}
	}
	return Standard03CounterBody{}, false
}

// Standard03CounterBodyLookupField requires the complete body identity and the exact field identity.
func Standard03CounterBodyLookupField(standard03CounterBodyKey Standard03CounterBodyKey) (Standard03CounterBodyField, bool) {
	for _, standard03CounterBodyBody := range Standard03CounterBodies() {
		for _, standard03CounterBodyField := range standard03CounterBodyBody.Standard03CounterBodyFields {
			standard03CounterBodyFieldKey := standard03CounterBodyBody.Standard03CounterBodyKey
			standard03CounterBodyFieldKey.Standard03CounterBodyField = standard03CounterBodyField.Standard03CounterBodyIdentity
			if standard03CounterBodyEqualKey(standard03CounterBodyFieldKey, standard03CounterBodyKey) {
				return standard03CounterBodyField, true
			}
		}
	}
	return Standard03CounterBodyField{}, false
}
