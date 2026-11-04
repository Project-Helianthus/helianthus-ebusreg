package vaillant

import (
	_ "embed"
	ab71JSON "encoding/json"
	ab71Reflect "reflect"
)

//go:embed data/ab71_b523_candidates.json
var ab71B523Data []byte

type Ab71Key struct {
	Ab71Family    string  `json:"family"`
	Ab71Variant   *string `json:"variant"`
	Ab71Primary   byte    `json:"primary"`
	Ab71Secondary *byte   `json:"secondary"`
	Ab71Operation *byte   `json:"operation"`
	Ab71Direction string  `json:"direction"`
	Ab71Field     string  `json:"field"`
}

type Ab71Rule struct {
	Ab71RawMin       byte     `json:"raw_min"`
	Ab71RawMax       byte     `json:"raw_max"`
	Ab71Result       string   `json:"result"`
	Ab71DivideBy     *float64 `json:"divide_by"`
	Ab71DurationUnit *string  `json:"duration_unit"`
	Ab71Grade        string   `json:"grade"`
}

type Ab71Alternative struct {
	Ab71Name  string     `json:"name"`
	Ab71Grade string     `json:"grade"`
	Ab71Rules []Ab71Rule `json:"rules"`
}

type Ab71Entry struct {
	Ab71Key             Ab71Key           `json:"key"`
	Ab71LogicalPosition *int              `json:"logical_position"`
	Ab71WidthBits       *int              `json:"width_bits"`
	Ab71RequestExtent   *int              `json:"request_extent"`
	Ab71PositiveExtent  *int              `json:"positive_extent"`
	Ab71DefaultRaw      *byte             `json:"default_raw"`
	Ab71DefaultResult   *string           `json:"default_result"`
	Ab71Grade           string            `json:"grade"`
	Ab71RawDomain       []byte            `json:"raw_domain"`
	Ab71Alternatives    []Ab71Alternative `json:"alternatives"`
}

type Ab71Catalog struct {
	Ab71RepresentationGrade string      `json:"representation_grade"`
	Ab71NativeSupport       string      `json:"native_support"`
	Ab71RuntimeEffects      string      `json:"runtime_effects"`
	Ab71OtherModuleLayouts  string      `json:"other_module_layouts"`
	Ab71Limitations         []string    `json:"limitations"`
	Ab71Entries             []Ab71Entry `json:"entries"`
}

func Ab71B523Candidates() (Ab71Catalog, error) {
	var ab71Catalog Ab71Catalog
	ab71Error := ab71JSON.Unmarshal(ab71B523Data, &ab71Catalog)
	return ab71Catalog, ab71Error
}

func Ab71Lookup(ab71Key Ab71Key) (Ab71Entry, bool, error) {
	ab71Catalog, ab71Error := Ab71B523Candidates()
	if ab71Error != nil {
		return Ab71Entry{}, false, ab71Error
	}
	for _, ab71Entry := range ab71Catalog.Ab71Entries {
		if ab71Reflect.DeepEqual(ab71Entry.Ab71Key, ab71Key) {
			return ab71Entry, true, nil
		}
	}
	return Ab71Entry{}, false, nil
}
