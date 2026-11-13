package vaillant

import (
	_ "embed"
	b509DiagnosticShapeJSON "encoding/json"
	b509DiagnosticShapeReflect "reflect"
)

type B509DiagnosticKey struct {
	B509DiagnosticShapeFamily    string  `json:"family"`
	B509DiagnosticShapeVariant   string  `json:"variant"`
	B509DiagnosticShapePrimary   *uint16 `json:"primary"`
	B509DiagnosticShapeSecondary *uint16 `json:"secondary"`
	B509DiagnosticShapeOperation *uint16 `json:"operation"`
	B509DiagnosticShapeDirection string  `json:"direction"`
	B509DiagnosticShapeField     string  `json:"field"`
	B509DiagnosticShapeOpcode    *uint16 `json:"opcode"`
	B509DiagnosticShapeGroup     *uint16 `json:"group"`
	B509DiagnosticShapeInstance  *uint16 `json:"instance"`
	B509DiagnosticShapeRegister  *uint16 `json:"register"`
}

type B509DiagnosticField struct {
	B509DiagnosticShapeKey              B509DiagnosticKey  `json:"key"`
	B509DiagnosticShapeLogicalPosition  *int               `json:"logical_position"`
	B509DiagnosticShapeWidthBits        *int               `json:"width_bits"`
	B509DiagnosticShapeSignedness       string             `json:"signedness"`
	B509DiagnosticShapeByteOrder        string             `json:"byte_order"`
	B509DiagnosticShapeConversion       string             `json:"conversion"`
	B509DiagnosticShapeUnit             *string            `json:"unit"`
	B509DiagnosticShapeDefault          *int64             `json:"default"`
	B509DiagnosticShapeMinimum          *int64             `json:"minimum"`
	B509DiagnosticShapeMaximum          *int64             `json:"maximum"`
	B509DiagnosticShapeEnumeration      map[string]*string `json:"enumeration"`
	B509DiagnosticShapeReplacement      []int64            `json:"replacement"`
	B509DiagnosticShapeReserved         []int64            `json:"reserved"`
	B509DiagnosticShapeGrade            string             `json:"grade"`
	B509DiagnosticShapeMappingGrade     string             `json:"mapping_grade"`
	B509DiagnosticShapeRawFallback      string             `json:"raw_fallback"`
	B509DiagnosticShapeRecommendedBytes *int               `json:"recommended_bytes"`
	B509DiagnosticShapeMaximumBytes     *int               `json:"maximum_bytes"`
	B509DiagnosticShapePadding          *string            `json:"padding"`
}

type B509DiagnosticCatalog struct {
	B509DiagnosticShapeRepresentation      string                `json:"representation"`
	B509DiagnosticShapeQualification       string                `json:"qualification"`
	B509DiagnosticShapePhysicalEffects     string                `json:"physical_effects"`
	B509DiagnosticShapeExecution           string                `json:"execution"`
	B509DiagnosticShapeResponseCorrelation string                `json:"response_correlation"`
	B509DiagnosticShapeLimitations         []string              `json:"limitations"`
	B509DiagnosticShapeFields              []B509DiagnosticField `json:"fields"`
}

//go:embed data/ab87_b509_diagnostic_shapes.json
var b509DiagnosticShapeDiagnosticJSON string

func B509DiagnosticShapes() (B509DiagnosticCatalog, error) {
	var b509DiagnosticShapeCatalog B509DiagnosticCatalog
	b509DiagnosticShapeError := b509DiagnosticShapeJSON.Unmarshal([]byte(b509DiagnosticShapeDiagnosticJSON), &b509DiagnosticShapeCatalog)
	return b509DiagnosticShapeCatalog, b509DiagnosticShapeError
}

func b509DiagnosticShapeLookupField(b509DiagnosticShapeFields []B509DiagnosticField, b509DiagnosticShapeKey B509DiagnosticKey) (B509DiagnosticField, bool, error) {
	for _, b509DiagnosticShapeField := range b509DiagnosticShapeFields {
		if b509DiagnosticShapeReflect.DeepEqual(b509DiagnosticShapeField.B509DiagnosticShapeKey, b509DiagnosticShapeKey) {
			b509DiagnosticShapeData, b509DiagnosticShapeError := b509DiagnosticShapeJSON.Marshal(b509DiagnosticShapeField)
			if b509DiagnosticShapeError != nil {
				return B509DiagnosticField{}, false, b509DiagnosticShapeError
			}
			var b509DiagnosticShapeCopy B509DiagnosticField
			b509DiagnosticShapeError = b509DiagnosticShapeJSON.Unmarshal(b509DiagnosticShapeData, &b509DiagnosticShapeCopy)
			return b509DiagnosticShapeCopy, b509DiagnosticShapeError == nil, b509DiagnosticShapeError
		}
	}
	return B509DiagnosticField{}, false, nil
}

func B509LookupDiagnosticField(b509DiagnosticShapeKey B509DiagnosticKey) (B509DiagnosticField, bool, error) {
	b509DiagnosticShapeCatalog, b509DiagnosticShapeError := B509DiagnosticShapes()
	if b509DiagnosticShapeError != nil {
		return B509DiagnosticField{}, false, b509DiagnosticShapeError
	}
	return b509DiagnosticShapeLookupField(b509DiagnosticShapeCatalog.B509DiagnosticShapeFields, b509DiagnosticShapeKey)
}
