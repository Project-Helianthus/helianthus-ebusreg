package vaillant

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type SchemaGrade string

const (
	SchemaUnknown    SchemaGrade = "Unknown"
	SchemaHypothesis SchemaGrade = "Hypothesis"
	SchemaProven     SchemaGrade = "Proven"
)

type SchemaKey struct {
	SchemaFamily      string `json:"family"`
	SchemaVariant     string `json:"variant"`
	SchemaDirection   string `json:"direction"`
	SchemaOpcode      uint8  `json:"opcode"`
	SchemaGroup       uint8  `json:"group"`
	SchemaInstance    uint8  `json:"instance"`
	SchemaRegister    uint16 `json:"register"`
	SchemaAlternative string `json:"alternative"`
}

// UnmarshalJSON requires explicit numeric selectors, including explicit zero,
// while keeping SchemaKey comparable for exact lookup and duplicate detection.
func (schemaKey *SchemaKey) UnmarshalJSON(schemaData []byte) error {
	type schemaKeyValue SchemaKey
	var schemaWire struct {
		schemaKeyValue
		SchemaOpcode   *uint8  `json:"opcode"`
		SchemaGroup    *uint8  `json:"group"`
		SchemaInstance *uint8  `json:"instance"`
		SchemaRegister *uint16 `json:"register"`
	}
	schemaDecoder := json.NewDecoder(bytes.NewReader(schemaData))
	schemaDecoder.DisallowUnknownFields()
	if schemaError := schemaDecoder.Decode(&schemaWire); schemaError != nil {
		return schemaError
	}
	if schemaWire.SchemaOpcode == nil || schemaWire.SchemaGroup == nil || schemaWire.SchemaInstance == nil || schemaWire.SchemaRegister == nil {
		return fmt.Errorf("schema key: opcode, group, instance and register must be supplied and non-null")
	}
	schemaValue := SchemaKey(schemaWire.schemaKeyValue)
	schemaValue.SchemaOpcode = *schemaWire.SchemaOpcode
	schemaValue.SchemaGroup = *schemaWire.SchemaGroup
	schemaValue.SchemaInstance = *schemaWire.SchemaInstance
	schemaValue.SchemaRegister = *schemaWire.SchemaRegister
	*schemaKey = schemaValue
	return nil
}

type SchemaField struct {
	SchemaName        string           `json:"name"`
	SchemaPosition    *int             `json:"position,omitempty"`
	SchemaWidth       *int             `json:"width,omitempty"`
	SchemaCodec       string           `json:"codec"`
	SchemaUnit        string           `json:"unit,omitempty"`
	SchemaEnum        map[int64]string `json:"enum,omitempty"`
	SchemaSentinel    []byte           `json:"sentinel,omitempty"`
	SchemaGrade       SchemaGrade      `json:"grade,omitempty"`
	SchemaRawFallback []byte           `json:"raw_fallback"`
}

type SchemaRecord struct {
	SchemaKey         SchemaKey     `json:"key"`
	SchemaBodyWidth   *int          `json:"body_width,omitempty"`
	SchemaGrade       SchemaGrade   `json:"grade,omitempty"`
	SchemaRawFallback []byte        `json:"raw_fallback"`
	SchemaFields      []SchemaField `json:"fields,omitempty"`
}

func schemaValidGrade(schemaGrade SchemaGrade) bool {
	return schemaGrade == "" || schemaGrade == SchemaUnknown || schemaGrade == SchemaHypothesis || schemaGrade == SchemaProven
}

func SchemaValidate(schemaRecords []SchemaRecord) error {
	schemaSeen := make(map[SchemaKey]bool)
	for schemaIndex, schemaRecord := range schemaRecords {
		schemaKey := schemaRecord.SchemaKey
		if schemaKey.SchemaFamily == "" || schemaKey.SchemaVariant == "" || schemaKey.SchemaAlternative == "" || (schemaKey.SchemaDirection != "request" && schemaKey.SchemaDirection != "response") {
			return fmt.Errorf("schema record %d: invalid key dimension", schemaIndex)
		}
		if schemaSeen[schemaRecord.SchemaKey] {
			return fmt.Errorf("schema record %d: duplicate alternative key", schemaIndex)
		}
		schemaSeen[schemaRecord.SchemaKey] = true
		if !schemaValidGrade(schemaRecord.SchemaGrade) {
			return fmt.Errorf("schema record %d: invalid grade", schemaIndex)
		}
		if schemaRecord.SchemaBodyWidth != nil && *schemaRecord.SchemaBodyWidth < 0 {
			return fmt.Errorf("schema record %d: negative body width", schemaIndex)
		}
		for schemaFieldIndex, schemaField := range schemaRecord.SchemaFields {
			if !schemaValidGrade(schemaField.SchemaGrade) {
				return fmt.Errorf("schema record %d field %d: invalid grade", schemaIndex, schemaFieldIndex)
			}
			if (schemaField.SchemaPosition != nil && *schemaField.SchemaPosition < 0) || (schemaField.SchemaWidth != nil && *schemaField.SchemaWidth <= 0) {
				return fmt.Errorf("schema record %d field %d: invalid extent", schemaIndex, schemaFieldIndex)
			}
			if schemaRecord.SchemaBodyWidth != nil && schemaField.SchemaPosition != nil {
				schemaRemaining := *schemaRecord.SchemaBodyWidth - *schemaField.SchemaPosition
				if schemaRemaining <= 0 || (schemaField.SchemaWidth != nil && *schemaField.SchemaWidth > schemaRemaining) {
					return fmt.Errorf("schema record %d field %d: extent exceeds body", schemaIndex, schemaFieldIndex)
				}
			} else if schemaRecord.SchemaBodyWidth != nil && schemaField.SchemaWidth != nil && *schemaField.SchemaWidth > *schemaRecord.SchemaBodyWidth {
				return fmt.Errorf("schema record %d field %d: width exceeds body", schemaIndex, schemaFieldIndex)
			}
		}
	}
	return nil
}

func schemaCopyInt(schemaValue *int) *int {
	if schemaValue == nil {
		return nil
	}
	schemaCopy := *schemaValue
	return &schemaCopy
}

func schemaClaimGrade(schemaGrade SchemaGrade) SchemaGrade {
	if schemaGrade == "" {
		return SchemaUnknown
	}
	return schemaGrade
}

func SchemaQuery(schemaRecords []SchemaRecord, schemaKey SchemaKey) (SchemaRecord, bool, error) {
	if schemaError := SchemaValidate(schemaRecords); schemaError != nil {
		return SchemaRecord{}, false, schemaError
	}
	for _, schemaRecord := range schemaRecords {
		if schemaRecord.SchemaKey != schemaKey {
			continue
		}
		schemaRecord.SchemaBodyWidth = schemaCopyInt(schemaRecord.SchemaBodyWidth)
		schemaRecord.SchemaGrade = schemaClaimGrade(schemaRecord.SchemaGrade)
		schemaRecord.SchemaRawFallback = bytes.Clone(schemaRecord.SchemaRawFallback)
		schemaRecord.SchemaFields = append([]SchemaField(nil), schemaRecord.SchemaFields...)
		for schemaIndex := range schemaRecord.SchemaFields {
			schemaField := &schemaRecord.SchemaFields[schemaIndex]
			schemaField.SchemaPosition = schemaCopyInt(schemaField.SchemaPosition)
			schemaField.SchemaWidth = schemaCopyInt(schemaField.SchemaWidth)
			schemaField.SchemaGrade = schemaClaimGrade(schemaField.SchemaGrade)
			schemaField.SchemaRawFallback = bytes.Clone(schemaField.SchemaRawFallback)
			schemaField.SchemaSentinel = bytes.Clone(schemaField.SchemaSentinel)
			if schemaField.SchemaEnum != nil {
				schemaEnum := make(map[int64]string, len(schemaField.SchemaEnum))
				for schemaValue, schemaLabel := range schemaField.SchemaEnum {
					schemaEnum[schemaValue] = schemaLabel
				}
				schemaField.SchemaEnum = schemaEnum
			}
		}
		return schemaRecord, true, nil
	}
	return SchemaRecord{}, false, nil
}
