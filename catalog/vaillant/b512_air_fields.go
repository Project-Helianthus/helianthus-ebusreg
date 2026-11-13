package vaillant

import (
	_ "embed"
	"encoding/json"
)

type B512AirFieldAxis struct {
	B512AirFieldPresent bool   `json:"present"`
	B512AirFieldValue   uint16 `json:"value"`
}

type B512AirKey struct {
	B512AirFieldFamily    string           `json:"family"`
	B512AirFieldVariant   string           `json:"variant"`
	B512AirFieldPrimary   uint16           `json:"primary"`
	B512AirFieldSecondary B512AirFieldAxis `json:"secondary"`
	B512AirFieldOperation B512AirFieldAxis `json:"operation"`
	B512AirFieldDirection string           `json:"direction"`
	B512AirFieldField     string           `json:"field"`
	B512AirFieldOpcode    B512AirFieldAxis `json:"opcode"`
	B512AirFieldGroup     B512AirFieldAxis `json:"group"`
	B512AirFieldInstance  B512AirFieldAxis `json:"instance"`
	B512AirFieldRegister  B512AirFieldAxis `json:"register"`
}

type B512AirField struct {
	B512AirFieldKey             B512AirKey                   `json:"key"`
	B512AirFieldLogicalPosition *int                         `json:"logical_position"`
	B512AirFieldWidth           *int                         `json:"width_bytes"`
	B512AirFieldSignedness      string                       `json:"signedness"`
	B512AirFieldByteOrder       string                       `json:"byte_order"`
	B512AirFieldConversion      string                       `json:"conversion"`
	B512AirFieldDivisor         float64                      `json:"divisor"`
	B512AirFieldUnit            *string                      `json:"unit"`
	B512AirFieldDefault         *float64                     `json:"default"`
	B512AirFieldDefaultLabel    *string                      `json:"default_label"`
	B512AirFieldBounds          []float64                    `json:"bounds"`
	B512AirFieldEnumeration     map[string]string            `json:"enumeration"`
	B512AirFieldOtherwise       *string                      `json:"otherwise"`
	B512AirFieldBits            map[string]map[string]string `json:"bits"`
	B512AirFieldReservedBits    []int                        `json:"reserved_bits"`
	B512AirFieldReplacement     *string                      `json:"replacement"`
	B512AirFieldGrade           string                       `json:"grade"`
	B512AirFieldRawFallback     string                       `json:"raw_fallback"`
}

type B512AirOperation struct {
	B512AirFieldOperation        uint16 `json:"operation"`
	B512AirFieldRequestArguments bool   `json:"request_arguments"`
	B512AirFieldGrade            string `json:"grade"`
}

type B512AirMetadata struct {
	B512AirFieldRepresentation  string             `json:"representation"`
	B512AirFieldPhysicalSupport string             `json:"physical_support"`
	B512AirFieldAcknowledgement string             `json:"acknowledgement"`
	B512AirFieldComplete        bool               `json:"complete"`
	B512AirFieldCoordinates     string             `json:"coordinates"`
	B512AirFieldLimitations     []string           `json:"limitations"`
	B512AirFieldOperations      []B512AirOperation `json:"operations"`
	B512AirFieldFields          []B512AirField     `json:"fields"`
}

//go:embed data/ab88_b512_air_fields.json
var b512AirFieldAirJSON []byte

func B512AirFields() B512AirMetadata {
	var b512AirFieldMetadata B512AirMetadata
	if b512AirFieldError := json.Unmarshal(b512AirFieldAirJSON, &b512AirFieldMetadata); b512AirFieldError != nil {
		panic(b512AirFieldError)
	}
	return b512AirFieldMetadata
}

func b512AirFieldFindAirField(b512AirFieldFields []B512AirField, b512AirFieldKey B512AirKey) (B512AirField, bool) {
	for _, b512AirFieldField := range b512AirFieldFields {
		if b512AirFieldField.B512AirFieldKey == b512AirFieldKey {
			return b512AirFieldField, true
		}
	}
	return B512AirField{}, false
}

func B512LookupAirField(b512AirFieldKey B512AirKey) (B512AirField, bool) {
	return b512AirFieldFindAirField(B512AirFields().B512AirFieldFields, b512AirFieldKey)
}
