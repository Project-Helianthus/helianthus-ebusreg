package vaillant

import (
	_ "embed"
	ab93JSON "encoding/json"
	ab93Reflect "reflect"
)

//go:embed data/ab93_b523_sensor_fields.json
var ab93B523SensorJSON string

type Ab93SensorKey struct {
	Ab93Family    string  `json:"family"`
	Ab93Variant   string  `json:"variant"`
	Ab93Primary   *uint16 `json:"primary"`
	Ab93Secondary *uint16 `json:"secondary"`
	Ab93Operation *uint16 `json:"operation"`
	Ab93Direction *string `json:"direction"`
	Ab93Field     string  `json:"field"`
	Ab93Opcode    *uint16 `json:"opcode"`
	Ab93Group     *uint16 `json:"group"`
	Ab93Instance  *uint16 `json:"instance"`
	Ab93Register  *uint16 `json:"register"`
}

type Ab93SensorBranch struct {
	Ab93Name       string    `json:"name"`
	Ab93Signedness string    `json:"signedness"`
	Ab93ByteOrder  string    `json:"byte_order"`
	Ab93DivideBy   *float64  `json:"divide_by"`
	Ab93Unit       *string   `json:"unit"`
	Ab93Default    *float64  `json:"default"`
	Ab93Bounds     []float64 `json:"bounds"`
	Ab93Grade      string    `json:"grade"`
}

type Ab93SensorField struct {
	Ab93Key            Ab93SensorKey      `json:"key"`
	Ab93Position       *int               `json:"logical_position"`
	Ab93Width          *int               `json:"width_bits"`
	Ab93BitOffset      *int               `json:"bit_offset"`
	Ab93Signedness     string             `json:"signedness"`
	Ab93ByteOrder      string             `json:"byte_order"`
	Ab93DivideBy       *float64           `json:"divide_by"`
	Ab93Unit           *string            `json:"unit"`
	Ab93Default        *float64           `json:"default"`
	Ab93Bounds         []float64          `json:"bounds"`
	Ab93Enumeration    map[string]string  `json:"enumeration"`
	Ab93ReplacementRaw *uint16            `json:"replacement_raw"`
	Ab93Selection      string             `json:"selection"`
	Ab93Branches       []Ab93SensorBranch `json:"branches"`
	Ab93Reserved       bool               `json:"reserved"`
	Ab93Grade          string             `json:"grade"`
	Ab93RawFallback    string             `json:"raw_fallback"`
}

type Ab93SensorCatalog struct {
	Ab93RepresentationGrade string            `json:"representation_grade"`
	Ab93Qualification       string            `json:"physical_qualification"`
	Ab93Complete            bool              `json:"complete"`
	Ab93Limitations         []string          `json:"limitations"`
	Ab93Fields              []Ab93SensorField `json:"fields"`
}

func Ab93B523SensorFields() (Ab93SensorCatalog, error) {
	var ab93Catalog Ab93SensorCatalog
	ab93Error := ab93JSON.Unmarshal([]byte(ab93B523SensorJSON), &ab93Catalog)
	return ab93Catalog, ab93Error
}

func ab93FindSensorField(ab93Fields []Ab93SensorField, ab93Key Ab93SensorKey) (Ab93SensorField, bool, error) {
	for _, ab93Field := range ab93Fields {
		if !ab93Reflect.DeepEqual(ab93Field.Ab93Key, ab93Key) {
			continue
		}
		ab93Data, ab93Error := ab93JSON.Marshal(ab93Field)
		if ab93Error != nil {
			return Ab93SensorField{}, false, ab93Error
		}
		var ab93Copy Ab93SensorField
		ab93Error = ab93JSON.Unmarshal(ab93Data, &ab93Copy)
		return ab93Copy, ab93Error == nil, ab93Error
	}
	return Ab93SensorField{}, false, nil
}

func Ab93LookupB523SensorField(ab93Key Ab93SensorKey) (Ab93SensorField, bool, error) {
	ab93Catalog, ab93Error := Ab93B523SensorFields()
	if ab93Error != nil {
		return Ab93SensorField{}, false, ab93Error
	}
	return ab93FindSensorField(ab93Catalog.Ab93Fields, ab93Key)
}
