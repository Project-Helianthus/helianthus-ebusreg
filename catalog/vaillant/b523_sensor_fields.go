package vaillant

import (
	_ "embed"
	b523SensorFieldJSON "encoding/json"
	b523SensorFieldReflect "reflect"
)

//go:embed data/ab93_b523_sensor_fields.json
var b523SensorFieldB523SensorJSON string

type B523SensorKey struct {
	B523SensorFieldFamily    string  `json:"family"`
	B523SensorFieldVariant   string  `json:"variant"`
	B523SensorFieldPrimary   *uint16 `json:"primary"`
	B523SensorFieldSecondary *uint16 `json:"secondary"`
	B523SensorFieldOperation *uint16 `json:"operation"`
	B523SensorFieldDirection *string `json:"direction"`
	B523SensorFieldField     string  `json:"field"`
	B523SensorFieldOpcode    *uint16 `json:"opcode"`
	B523SensorFieldGroup     *uint16 `json:"group"`
	B523SensorFieldInstance  *uint16 `json:"instance"`
	B523SensorFieldRegister  *uint16 `json:"register"`
}

type B523SensorBranch struct {
	B523SensorFieldName       string    `json:"name"`
	B523SensorFieldSignedness string    `json:"signedness"`
	B523SensorFieldByteOrder  string    `json:"byte_order"`
	B523SensorFieldDivideBy   *float64  `json:"divide_by"`
	B523SensorFieldUnit       *string   `json:"unit"`
	B523SensorFieldDefault    *float64  `json:"default"`
	B523SensorFieldBounds     []float64 `json:"bounds"`
	B523SensorFieldGrade      string    `json:"grade"`
}

type B523SensorField struct {
	B523SensorFieldKey            B523SensorKey      `json:"key"`
	B523SensorFieldPosition       *int               `json:"logical_position"`
	B523SensorFieldWidth          *int               `json:"width_bits"`
	B523SensorFieldBitOffset      *int               `json:"bit_offset"`
	B523SensorFieldSignedness     string             `json:"signedness"`
	B523SensorFieldByteOrder      string             `json:"byte_order"`
	B523SensorFieldDivideBy       *float64           `json:"divide_by"`
	B523SensorFieldUnit           *string            `json:"unit"`
	B523SensorFieldDefault        *float64           `json:"default"`
	B523SensorFieldBounds         []float64          `json:"bounds"`
	B523SensorFieldEnumeration    map[string]string  `json:"enumeration"`
	B523SensorFieldReplacementRaw *uint16            `json:"replacement_raw"`
	B523SensorFieldSelection      string             `json:"selection"`
	B523SensorFieldBranches       []B523SensorBranch `json:"branches"`
	B523SensorFieldReserved       bool               `json:"reserved"`
	B523SensorFieldGrade          string             `json:"grade"`
	B523SensorFieldRawFallback    string             `json:"raw_fallback"`
}

type B523SensorCatalog struct {
	B523SensorFieldRepresentationGrade string            `json:"representation_grade"`
	B523SensorFieldQualification       string            `json:"physical_qualification"`
	B523SensorFieldComplete            bool              `json:"complete"`
	B523SensorFieldLimitations         []string          `json:"limitations"`
	B523SensorFieldFields              []B523SensorField `json:"fields"`
}

func B523SensorFields() (B523SensorCatalog, error) {
	var b523SensorFieldCatalog B523SensorCatalog
	b523SensorFieldError := b523SensorFieldJSON.Unmarshal([]byte(b523SensorFieldB523SensorJSON), &b523SensorFieldCatalog)
	return b523SensorFieldCatalog, b523SensorFieldError
}

func b523SensorFieldFindSensorField(b523SensorFieldFields []B523SensorField, b523SensorFieldKey B523SensorKey) (B523SensorField, bool, error) {
	for _, b523SensorFieldField := range b523SensorFieldFields {
		if !b523SensorFieldReflect.DeepEqual(b523SensorFieldField.B523SensorFieldKey, b523SensorFieldKey) {
			continue
		}
		b523SensorFieldData, b523SensorFieldError := b523SensorFieldJSON.Marshal(b523SensorFieldField)
		if b523SensorFieldError != nil {
			return B523SensorField{}, false, b523SensorFieldError
		}
		var b523SensorFieldCopy B523SensorField
		b523SensorFieldError = b523SensorFieldJSON.Unmarshal(b523SensorFieldData, &b523SensorFieldCopy)
		return b523SensorFieldCopy, b523SensorFieldError == nil, b523SensorFieldError
	}
	return B523SensorField{}, false, nil
}

func B523LookupSensorField(b523SensorFieldKey B523SensorKey) (B523SensorField, bool, error) {
	b523SensorFieldCatalog, b523SensorFieldError := B523SensorFields()
	if b523SensorFieldError != nil {
		return B523SensorField{}, false, b523SensorFieldError
	}
	return b523SensorFieldFindSensorField(b523SensorFieldCatalog.B523SensorFieldFields, b523SensorFieldKey)
}
