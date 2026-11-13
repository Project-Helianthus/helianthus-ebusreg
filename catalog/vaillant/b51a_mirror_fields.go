package vaillant

import (
	_ "embed"
	b51AMirrorFieldJSON "encoding/json"
	b51AMirrorFieldReflect "reflect"
)

//go:embed data/ab90_b51a_mirror_fields.json
var b51AMirrorFieldMirrorJSON string

type B51AMirrorKey struct {
	B51AMirrorFieldFamily    string  `json:"family"`
	B51AMirrorFieldVariant   string  `json:"variant"`
	B51AMirrorFieldPrimary   uint8   `json:"primary"`
	B51AMirrorFieldSecondary *uint8  `json:"secondary"`
	B51AMirrorFieldOperation *uint8  `json:"operation"`
	B51AMirrorFieldDirection string  `json:"direction"`
	B51AMirrorFieldField     string  `json:"field"`
	B51AMirrorFieldOpcode    *uint8  `json:"opcode"`
	B51AMirrorFieldGroup     *uint8  `json:"group"`
	B51AMirrorFieldInstance  *uint8  `json:"instance"`
	B51AMirrorFieldRegister  *uint16 `json:"register"`
}

type B51AMirrorBit struct {
	B51AMirrorFieldIndex   int         `json:"index"`
	B51AMirrorFieldName    *string     `json:"name"`
	B51AMirrorFieldRole    string      `json:"role"`
	B51AMirrorFieldDefault *bool       `json:"default"`
	B51AMirrorFieldGrade   SchemaGrade `json:"grade"`
}

type B51AMirrorAlternative struct {
	B51AMirrorFieldName        string           `json:"name"`
	B51AMirrorFieldWidthBytes  *int             `json:"width_bytes"`
	B51AMirrorFieldCondition   string           `json:"condition"`
	B51AMirrorFieldEnumeration map[int64]string `json:"enumeration"`
	B51AMirrorFieldBits        []B51AMirrorBit  `json:"bits"`
	B51AMirrorFieldGrade       SchemaGrade      `json:"grade"`
}

type B51AMirrorField struct {
	B51AMirrorFieldKey             B51AMirrorKey           `json:"key"`
	B51AMirrorFieldLogicalPosition *int                    `json:"logical_position"`
	B51AMirrorFieldWidthBytes      *int                    `json:"width_bytes"`
	B51AMirrorFieldSignedness      string                  `json:"signedness"`
	B51AMirrorFieldByteOrder       string                  `json:"byte_order"`
	B51AMirrorFieldConversion      string                  `json:"conversion"`
	B51AMirrorFieldUnit            *string                 `json:"unit"`
	B51AMirrorFieldDefault         *int64                  `json:"default"`
	B51AMirrorFieldFormalMin       *int64                  `json:"formal_min"`
	B51AMirrorFieldFormalMax       *int64                  `json:"formal_max"`
	B51AMirrorFieldLogicalMin      *int64                  `json:"logical_min"`
	B51AMirrorFieldLogicalMax      *int64                  `json:"logical_max"`
	B51AMirrorFieldLogicalLimit    *string                 `json:"logical_limit"`
	B51AMirrorFieldEnumeration     map[int64]string        `json:"enumeration"`
	B51AMirrorFieldBits            []B51AMirrorBit         `json:"bits"`
	B51AMirrorFieldAlternatives    []B51AMirrorAlternative `json:"alternatives"`
	B51AMirrorFieldReplacement     *string                 `json:"replacement"`
	B51AMirrorFieldReserved        []int                   `json:"reserved"`
	B51AMirrorFieldGrade           SchemaGrade             `json:"grade"`
	B51AMirrorFieldRawFallback     string                  `json:"raw_fallback"`
	B51AMirrorFieldCorrelation     string                  `json:"correlation"`
	B51AMirrorFieldUnmapped        string                  `json:"unmapped"`
}

type B51AMirrorMetadata struct {
	B51AMirrorFieldRepresentationGrade SchemaGrade       `json:"representation_grade"`
	B51AMirrorFieldSemanticsGrade      SchemaGrade       `json:"semantics_grade"`
	B51AMirrorFieldPhysicalAdaptation  SchemaGrade       `json:"physical_adaptation"`
	B51AMirrorFieldStorageDurability   SchemaGrade       `json:"storage_durability"`
	B51AMirrorFieldSuccessSemantics    SchemaGrade       `json:"success_semantics"`
	B51AMirrorFieldComplete            bool              `json:"complete"`
	B51AMirrorFieldNativeContract      string            `json:"native_contract"`
	B51AMirrorFieldLimitations         []string          `json:"limitations"`
	B51AMirrorFieldFields              []B51AMirrorField `json:"fields"`
}

func B51AMirrorFields() (B51AMirrorMetadata, error) {
	var b51AMirrorFieldMetadata B51AMirrorMetadata
	b51AMirrorFieldError := b51AMirrorFieldJSON.Unmarshal([]byte(b51AMirrorFieldMirrorJSON), &b51AMirrorFieldMetadata)
	return b51AMirrorFieldMetadata, b51AMirrorFieldError
}

func b51AMirrorFieldMirrorKeysEqual(b51AMirrorFieldLeft, b51AMirrorFieldRight B51AMirrorKey) bool {
	return b51AMirrorFieldReflect.DeepEqual(b51AMirrorFieldLeft, b51AMirrorFieldRight)
}

func B51ALookupMirrorField(b51AMirrorFieldKey B51AMirrorKey) (B51AMirrorField, bool, error) {
	b51AMirrorFieldMetadata, b51AMirrorFieldError := B51AMirrorFields()
	if b51AMirrorFieldError != nil {
		return B51AMirrorField{}, false, b51AMirrorFieldError
	}
	for _, b51AMirrorFieldField := range b51AMirrorFieldMetadata.B51AMirrorFieldFields {
		if b51AMirrorFieldMirrorKeysEqual(b51AMirrorFieldField.B51AMirrorFieldKey, b51AMirrorFieldKey) {
			return b51AMirrorFieldField, true, nil
		}
	}
	return B51AMirrorField{}, false, nil
}
