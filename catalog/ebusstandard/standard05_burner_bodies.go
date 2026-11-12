package ebusstandard

import (
	_ "embed"
	standard05BurnerBodyJSON "encoding/json"
)

type Standard05BurnerBodyAxis struct {
	Standard05BurnerBodyPresent bool   `json:"present"`
	Standard05BurnerBodyValue   uint16 `json:"value"`
}

type Standard05BurnerBodyKey struct {
	Standard05BurnerBodyFamily    string                   `json:"family"`
	Standard05BurnerBodyVariant   string                   `json:"variant"`
	Standard05BurnerBodyPrimary   Standard05BurnerBodyAxis `json:"primary"`
	Standard05BurnerBodySecondary Standard05BurnerBodyAxis `json:"secondary"`
	Standard05BurnerBodyOperation Standard05BurnerBodyAxis `json:"operation"`
	Standard05BurnerBodyDirection string                   `json:"direction"`
	Standard05BurnerBodyField     string                   `json:"field"`
	Standard05BurnerBodyOpcode    Standard05BurnerBodyAxis `json:"opcode"`
	Standard05BurnerBodyGroup     Standard05BurnerBodyAxis `json:"group"`
	Standard05BurnerBodyInstance  Standard05BurnerBodyAxis `json:"instance"`
	Standard05BurnerBodyRegister  Standard05BurnerBodyAxis `json:"register"`
}

type Standard05BurnerBodyBounds struct {
	Standard05BurnerBodyMin int `json:"min"`
	Standard05BurnerBodyMax int `json:"max"`
}

type Standard05BurnerField struct {
	Standard05BurnerBodyKey             Standard05BurnerBodyKey     `json:"key"`
	Standard05BurnerBodyLogicalPosition *int                        `json:"logical_position"`
	Standard05BurnerBodyWidthBits       *int                        `json:"width_bits"`
	Standard05BurnerBodySignedness      string                      `json:"signedness"`
	Standard05BurnerBodyByteOrder       string                      `json:"byte_order"`
	Standard05BurnerBodyConversion      string                      `json:"conversion"`
	Standard05BurnerBodyUnit            string                      `json:"unit"`
	Standard05BurnerBodyDefault         *int                        `json:"default"`
	Standard05BurnerBodyBounds          *Standard05BurnerBodyBounds `json:"bounds"`
	Standard05BurnerBodyEnumeration     map[int]string              `json:"enumeration"`
	Standard05BurnerBodyReplacement     []int                       `json:"replacement"`
	Standard05BurnerBodyReserved        []int                       `json:"reserved"`
	Standard05BurnerBodyConstant        *int                        `json:"constant"`
	Standard05BurnerBodyGrade           string                      `json:"grade"`
	Standard05BurnerBodyRawFallback     string                      `json:"raw_fallback"`
}

type Standard05BurnerBodies struct {
	Standard05BurnerBodyRepresentationGrade string                  `json:"representation_grade"`
	Standard05BurnerBodyComplete            bool                    `json:"complete"`
	Standard05BurnerBodyCoordinates         string                  `json:"coordinates"`
	Standard05BurnerBodyNativeSupport       string                  `json:"native_support"`
	Standard05BurnerBodyQualification       string                  `json:"qualification"`
	Standard05BurnerBodyFreshness           string                  `json:"freshness"`
	Standard05BurnerBodyPersistence         string                  `json:"persistence"`
	Standard05BurnerBodyExecutionPermission string                  `json:"execution_permission"`
	Standard05BurnerBodyFields              []Standard05BurnerField `json:"fields"`
}

//go:embed testdata/ab84_05_burner_bodies.json
var standard05BurnerBodyBurnerJSON []byte

func Standard05BurnerMetadata() Standard05BurnerBodies {
	var standard05BurnerBodyBodies Standard05BurnerBodies
	if standard05BurnerBodyError := standard05BurnerBodyJSON.Unmarshal(standard05BurnerBodyBurnerJSON, &standard05BurnerBodyBodies); standard05BurnerBodyError != nil {
		panic(standard05BurnerBodyError)
	}
	return standard05BurnerBodyBodies
}

func Standard05LookupBurnerField(standard05BurnerBodyKey Standard05BurnerBodyKey) (Standard05BurnerField, bool) {
	return standard05BurnerBodyLookupBurnerField(Standard05BurnerMetadata().Standard05BurnerBodyFields, standard05BurnerBodyKey)
}

func standard05BurnerBodyLookupBurnerField(standard05BurnerBodyFields []Standard05BurnerField, standard05BurnerBodyKey Standard05BurnerBodyKey) (Standard05BurnerField, bool) {
	for _, standard05BurnerBodyField := range standard05BurnerBodyFields {
		if standard05BurnerBodyField.Standard05BurnerBodyKey == standard05BurnerBodyKey {
			return standard05BurnerBodyField, true
		}
	}
	return Standard05BurnerField{}, false
}
