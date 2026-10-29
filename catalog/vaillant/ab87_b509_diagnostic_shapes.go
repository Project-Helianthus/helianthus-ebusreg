package vaillant

import (
	_ "embed"
	ab87JSON "encoding/json"
	ab87Reflect "reflect"
)

type Ab87DiagnosticKey struct {
	Ab87Family    string  `json:"family"`
	Ab87Variant   string  `json:"variant"`
	Ab87Primary   *uint16 `json:"primary"`
	Ab87Secondary *uint16 `json:"secondary"`
	Ab87Operation *uint16 `json:"operation"`
	Ab87Direction string  `json:"direction"`
	Ab87Field     string  `json:"field"`
	Ab87Opcode    *uint16 `json:"opcode"`
	Ab87Group     *uint16 `json:"group"`
	Ab87Instance  *uint16 `json:"instance"`
	Ab87Register  *uint16 `json:"register"`
}

type Ab87DiagnosticField struct {
	Ab87Key              Ab87DiagnosticKey  `json:"key"`
	Ab87LogicalPosition  *int               `json:"logical_position"`
	Ab87WidthBits        *int               `json:"width_bits"`
	Ab87Signedness       string             `json:"signedness"`
	Ab87ByteOrder        string             `json:"byte_order"`
	Ab87Conversion       string             `json:"conversion"`
	Ab87Unit             *string            `json:"unit"`
	Ab87Default          *int64             `json:"default"`
	Ab87Minimum          *int64             `json:"minimum"`
	Ab87Maximum          *int64             `json:"maximum"`
	Ab87Enumeration      map[string]*string `json:"enumeration"`
	Ab87Replacement      []int64            `json:"replacement"`
	Ab87Reserved         []int64            `json:"reserved"`
	Ab87Grade            string             `json:"grade"`
	Ab87MappingGrade     string             `json:"mapping_grade"`
	Ab87RawFallback      string             `json:"raw_fallback"`
	Ab87RecommendedBytes *int               `json:"recommended_bytes"`
	Ab87MaximumBytes     *int               `json:"maximum_bytes"`
	Ab87Padding          *string            `json:"padding"`
}

type Ab87DiagnosticCatalog struct {
	Ab87Representation      string                `json:"representation"`
	Ab87Qualification       string                `json:"qualification"`
	Ab87PhysicalEffects     string                `json:"physical_effects"`
	Ab87Execution           string                `json:"execution"`
	Ab87ResponseCorrelation string                `json:"response_correlation"`
	Ab87Limitations         []string              `json:"limitations"`
	Ab87Fields              []Ab87DiagnosticField `json:"fields"`
}

//go:embed data/ab87_b509_diagnostic_shapes.json
var ab87DiagnosticJSON string

func Ab87B509DiagnosticShapes() (Ab87DiagnosticCatalog, error) {
	var ab87Catalog Ab87DiagnosticCatalog
	ab87Error := ab87JSON.Unmarshal([]byte(ab87DiagnosticJSON), &ab87Catalog)
	return ab87Catalog, ab87Error
}

func ab87LookupField(ab87Fields []Ab87DiagnosticField, ab87Key Ab87DiagnosticKey) (Ab87DiagnosticField, bool, error) {
	for _, ab87Field := range ab87Fields {
		if ab87Reflect.DeepEqual(ab87Field.Ab87Key, ab87Key) {
			ab87Data, ab87Error := ab87JSON.Marshal(ab87Field)
			if ab87Error != nil {
				return Ab87DiagnosticField{}, false, ab87Error
			}
			var ab87Copy Ab87DiagnosticField
			ab87Error = ab87JSON.Unmarshal(ab87Data, &ab87Copy)
			return ab87Copy, ab87Error == nil, ab87Error
		}
	}
	return Ab87DiagnosticField{}, false, nil
}

func Ab87LookupB509DiagnosticField(ab87Key Ab87DiagnosticKey) (Ab87DiagnosticField, bool, error) {
	ab87Catalog, ab87Error := Ab87B509DiagnosticShapes()
	if ab87Error != nil {
		return Ab87DiagnosticField{}, false, ab87Error
	}
	return ab87LookupField(ab87Catalog.Ab87Fields, ab87Key)
}
