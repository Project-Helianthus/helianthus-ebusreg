package vaillant

import (
	_ "embed"
	ab85json "encoding/json"
	ab85reflect "reflect"
)

//go:embed data/ab85_b524_scalar_shapes.json
var ab85ScalarShapesJSON []byte

type Ab85Key struct {
	Ab85Family    *string `json:"family"`
	Ab85Variant   *string `json:"variant"`
	Ab85Primary   *uint8  `json:"primary"`
	Ab85Secondary *uint8  `json:"secondary"`
	Ab85Operation *uint8  `json:"operation"`
	Ab85Direction *string `json:"direction"`
	Ab85Field     *string `json:"field"`
	Ab85Opcode    *uint8  `json:"opcode"`
	Ab85Group     *uint8  `json:"group"`
	Ab85Instance  *uint8  `json:"instance"`
	Ab85Register  *uint16 `json:"register"`
}

type Ab85Conversion struct {
	Ab85Encoding string  `json:"encoding"`
	Ab85Factor   float64 `json:"factor"`
	Ab85Offset   float64 `json:"offset"`
}

type Ab85Bounds struct {
	Ab85Minimum *float64 `json:"minimum"`
	Ab85Maximum *float64 `json:"maximum"`
}

type Ab85ScalarShape struct {
	Ab85Key               Ab85Key           `json:"key"`
	Ab85Shape             string            `json:"shape"`
	Ab85LogicalPosition   *int              `json:"logical_position"`
	Ab85WidthBytes        *int              `json:"width_bytes"`
	Ab85Signedness        string            `json:"signedness"`
	Ab85ByteOrder         string            `json:"byte_order"`
	Ab85Conversion        *Ab85Conversion   `json:"conversion"`
	Ab85Unit              *string           `json:"unit"`
	Ab85Default           *float64          `json:"default"`
	Ab85Bounds            *Ab85Bounds       `json:"bounds"`
	Ab85Enumeration       map[string]string `json:"enumeration"`
	Ab85ReplacementValues []string          `json:"replacement_values"`
	Ab85ReservedValues    []string          `json:"reserved_values"`
	Ab85Grade             string            `json:"grade"`
	Ab85RawFallback       string            `json:"raw_fallback"`
}

type Ab85Catalog struct {
	Ab85Catalog               string            `json:"catalog"`
	Ab85RepresentationGrade   string            `json:"representation_grade"`
	Ab85PhysicalQualification string            `json:"physical_qualification"`
	Ab85Complete              bool              `json:"complete"`
	Ab85Limitations           []string          `json:"limitations"`
	Ab85Entries               []Ab85ScalarShape `json:"entries"`
}

func Ab85ScalarShapes() (Ab85Catalog, error) {
	var ab85Catalog Ab85Catalog
	ab85Error := ab85json.Unmarshal(ab85ScalarShapesJSON, &ab85Catalog)
	return ab85Catalog, ab85Error
}

func Ab85Lookup(ab85Key Ab85Key) ([]Ab85ScalarShape, error) {
	ab85Catalog, ab85Error := Ab85ScalarShapes()
	if ab85Error != nil {
		return nil, ab85Error
	}
	var ab85Matches []Ab85ScalarShape
	for _, ab85Entry := range ab85Catalog.Ab85Entries {
		if ab85reflect.DeepEqual(ab85Entry.Ab85Key, ab85Key) {
			ab85Matches = append(ab85Matches, ab85Entry)
		}
	}
	return ab85Matches, nil
}
