package vaillant

import (
	_ "embed"
	ab89json "encoding/json"
)

type Ab89Axis struct {
	Ab89Present bool   `json:"present"`
	Ab89Value   uint16 `json:"value"`
}

type Ab89Key struct {
	Ab89Family    string   `json:"family"`
	Ab89Variant   string   `json:"variant"`
	Ab89Primary   uint16   `json:"primary"`
	Ab89Secondary Ab89Axis `json:"secondary"`
	Ab89Operation Ab89Axis `json:"operation"`
	Ab89Direction string   `json:"direction"`
	Ab89Field     string   `json:"field"`
	Ab89Opcode    Ab89Axis `json:"opcode"`
	Ab89Group     Ab89Axis `json:"group"`
	Ab89Instance  Ab89Axis `json:"instance"`
	Ab89Register  Ab89Axis `json:"register"`
}

type Ab89Conversion struct {
	Ab89Factor float64 `json:"factor"`
	Ab89Offset float64 `json:"offset"`
}

type Ab89Range struct {
	Ab89Minimum int    `json:"minimum"`
	Ab89Maximum int    `json:"maximum"`
	Ab89Meaning string `json:"meaning"`
}

type Ab89Field struct {
	Ab89Key         Ab89Key         `json:"key"`
	Ab89Position    *int            `json:"logical_position"`
	Ab89Width       *int            `json:"width"`
	Ab89Extent      *int            `json:"logical_extent"`
	Ab89Signedness  string          `json:"signedness"`
	Ab89ByteOrder   string          `json:"byte_order"`
	Ab89Conversion  *Ab89Conversion `json:"conversion"`
	Ab89Unit        *string         `json:"unit"`
	Ab89Default     *int            `json:"default"`
	Ab89Bounds      []int           `json:"bounds"`
	Ab89Enumeration map[int]string  `json:"enumeration"`
	Ab89Ranges      []Ab89Range     `json:"ranges"`
	Ab89Replacement *string         `json:"replacement"`
	Ab89Reserved    []Ab89Range     `json:"reserved"`
	Ab89Constant    *int            `json:"constant"`
	Ab89Conflict    *string         `json:"conflict"`
	Ab89Grade       string          `json:"grade"`
	Ab89RawFallback string          `json:"raw_fallback"`
}

type Ab89Catalog struct {
	Ab89RepresentationGrade string      `json:"representation_grade"`
	Ab89Qualification       string      `json:"physical_qualification"`
	Ab89Limitations         []string    `json:"limitations"`
	Ab89Fields              []Ab89Field `json:"fields"`
}

//go:embed data/ab89_b514_menu_feedback.json
var ab89MenuFeedbackJSON string

func Ab89MenuFeedback() (Ab89Catalog, error) {
	var ab89Catalog Ab89Catalog
	ab89Error := ab89json.Unmarshal([]byte(ab89MenuFeedbackJSON), &ab89Catalog)
	return ab89Catalog, ab89Error
}

func Ab89Lookup(ab89Key Ab89Key) (Ab89Field, bool, error) {
	ab89Catalog, ab89Error := Ab89MenuFeedback()
	if ab89Error != nil {
		return Ab89Field{}, false, ab89Error
	}
	return ab89LookupFields(ab89Catalog.Ab89Fields, ab89Key)
}

func ab89LookupFields(ab89Fields []Ab89Field, ab89Key Ab89Key) (Ab89Field, bool, error) {
	for _, ab89Field := range ab89Fields {
		if ab89Field.Ab89Key == ab89Key {
			return ab89Field, true, nil
		}
	}
	return Ab89Field{}, false, nil
}
