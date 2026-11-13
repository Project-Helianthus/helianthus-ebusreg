package vaillant

import (
	_ "embed"
	b524ScalarShapeJSON "encoding/json"
	b524ScalarShapeReflect "reflect"
)

//go:embed data/ab85_b524_scalar_shapes.json
var b524ScalarShapeScalarShapesJSON []byte

type B524ScalarShapeKey struct {
	B524ScalarShapeFamily    *string `json:"family"`
	B524ScalarShapeVariant   *string `json:"variant"`
	B524ScalarShapePrimary   *uint8  `json:"primary"`
	B524ScalarShapeSecondary *uint8  `json:"secondary"`
	B524ScalarShapeOperation *uint8  `json:"operation"`
	B524ScalarShapeDirection *string `json:"direction"`
	B524ScalarShapeField     *string `json:"field"`
	B524ScalarShapeOpcode    *uint8  `json:"opcode"`
	B524ScalarShapeGroup     *uint8  `json:"group"`
	B524ScalarShapeInstance  *uint8  `json:"instance"`
	B524ScalarShapeRegister  *uint16 `json:"register"`
}

type B524ScalarShapeConversion struct {
	B524ScalarShapeEncoding string  `json:"encoding"`
	B524ScalarShapeFactor   float64 `json:"factor"`
	B524ScalarShapeOffset   float64 `json:"offset"`
}

type B524ScalarShapeBounds struct {
	B524ScalarShapeMinimum *float64 `json:"minimum"`
	B524ScalarShapeMaximum *float64 `json:"maximum"`
}

type B524ScalarShape struct {
	B524ScalarShapeKey               B524ScalarShapeKey         `json:"key"`
	B524ScalarShapeShape             string                     `json:"shape"`
	B524ScalarShapeLogicalPosition   *int                       `json:"logical_position"`
	B524ScalarShapeWidthBytes        *int                       `json:"width_bytes"`
	B524ScalarShapeSignedness        string                     `json:"signedness"`
	B524ScalarShapeByteOrder         string                     `json:"byte_order"`
	B524ScalarShapeConversion        *B524ScalarShapeConversion `json:"conversion"`
	B524ScalarShapeUnit              *string                    `json:"unit"`
	B524ScalarShapeDefault           *float64                   `json:"default"`
	B524ScalarShapeBounds            *B524ScalarShapeBounds     `json:"bounds"`
	B524ScalarShapeEnumeration       map[string]string          `json:"enumeration"`
	B524ScalarShapeReplacementValues []string                   `json:"replacement_values"`
	B524ScalarShapeReservedValues    []string                   `json:"reserved_values"`
	B524ScalarShapeGrade             string                     `json:"grade"`
	B524ScalarShapeRawFallback       string                     `json:"raw_fallback"`
}

type B524ScalarShapeCatalog struct {
	B524ScalarShapeCatalog               string            `json:"catalog"`
	B524ScalarShapeRepresentationGrade   string            `json:"representation_grade"`
	B524ScalarShapePhysicalQualification string            `json:"physical_qualification"`
	B524ScalarShapeComplete              bool              `json:"complete"`
	B524ScalarShapeLimitations           []string          `json:"limitations"`
	B524ScalarShapeEntries               []B524ScalarShape `json:"entries"`
}

func B524ScalarShapes() (B524ScalarShapeCatalog, error) {
	var b524ScalarShapeCatalog B524ScalarShapeCatalog
	b524ScalarShapeError := b524ScalarShapeJSON.Unmarshal(b524ScalarShapeScalarShapesJSON, &b524ScalarShapeCatalog)
	return b524ScalarShapeCatalog, b524ScalarShapeError
}

func B524ScalarShapeLookup(b524ScalarShapeKey B524ScalarShapeKey) ([]B524ScalarShape, error) {
	b524ScalarShapeCatalog, b524ScalarShapeError := B524ScalarShapes()
	if b524ScalarShapeError != nil {
		return nil, b524ScalarShapeError
	}
	var b524ScalarShapeMatches []B524ScalarShape
	for _, b524ScalarShapeEntry := range b524ScalarShapeCatalog.B524ScalarShapeEntries {
		if b524ScalarShapeReflect.DeepEqual(b524ScalarShapeEntry.B524ScalarShapeKey, b524ScalarShapeKey) {
			b524ScalarShapeMatches = append(b524ScalarShapeMatches, b524ScalarShapeEntry)
		}
	}
	return b524ScalarShapeMatches, nil
}
