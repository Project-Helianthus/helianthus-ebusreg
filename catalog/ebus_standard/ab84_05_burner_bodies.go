package ebus_standard_catalog

import (
	_ "embed"
	ab84json "encoding/json"
)

type Ab84Axis struct {
	Ab84Present bool   `json:"present"`
	Ab84Value   uint16 `json:"value"`
}

type Ab84Key struct {
	Ab84Family    string   `json:"family"`
	Ab84Variant   string   `json:"variant"`
	Ab84Primary   Ab84Axis `json:"primary"`
	Ab84Secondary Ab84Axis `json:"secondary"`
	Ab84Operation Ab84Axis `json:"operation"`
	Ab84Direction string   `json:"direction"`
	Ab84Field     string   `json:"field"`
	Ab84Opcode    Ab84Axis `json:"opcode"`
	Ab84Group     Ab84Axis `json:"group"`
	Ab84Instance  Ab84Axis `json:"instance"`
	Ab84Register  Ab84Axis `json:"register"`
}

type Ab84Bounds struct {
	Ab84Min int `json:"min"`
	Ab84Max int `json:"max"`
}

type Ab84BurnerField struct {
	Ab84Key             Ab84Key        `json:"key"`
	Ab84LogicalPosition *int           `json:"logical_position"`
	Ab84WidthBits       *int           `json:"width_bits"`
	Ab84Signedness      string         `json:"signedness"`
	Ab84ByteOrder       string         `json:"byte_order"`
	Ab84Conversion      string         `json:"conversion"`
	Ab84Unit            string         `json:"unit"`
	Ab84Default         *int           `json:"default"`
	Ab84Bounds          *Ab84Bounds    `json:"bounds"`
	Ab84Enumeration     map[int]string `json:"enumeration"`
	Ab84Replacement     []int          `json:"replacement"`
	Ab84Reserved        []int          `json:"reserved"`
	Ab84Constant        *int           `json:"constant"`
	Ab84Grade           string         `json:"grade"`
	Ab84RawFallback     string         `json:"raw_fallback"`
}

type Ab84BurnerBodies struct {
	Ab84RepresentationGrade string            `json:"representation_grade"`
	Ab84Complete            bool              `json:"complete"`
	Ab84Coordinates         string            `json:"coordinates"`
	Ab84NativeSupport       string            `json:"native_support"`
	Ab84Qualification       string            `json:"qualification"`
	Ab84Freshness           string            `json:"freshness"`
	Ab84Persistence         string            `json:"persistence"`
	Ab84ExecutionPermission string            `json:"execution_permission"`
	Ab84Fields              []Ab84BurnerField `json:"fields"`
}

//go:embed testdata/ab84_05_burner_bodies.json
var ab84BurnerJSON []byte

func Ab84BurnerMetadata() Ab84BurnerBodies {
	var ab84Bodies Ab84BurnerBodies
	if ab84Error := ab84json.Unmarshal(ab84BurnerJSON, &ab84Bodies); ab84Error != nil {
		panic(ab84Error)
	}
	return ab84Bodies
}

func Ab84LookupBurnerField(ab84Key Ab84Key) (Ab84BurnerField, bool) {
	return ab84LookupBurnerField(Ab84BurnerMetadata().Ab84Fields, ab84Key)
}

func ab84LookupBurnerField(ab84Fields []Ab84BurnerField, ab84Key Ab84Key) (Ab84BurnerField, bool) {
	for _, ab84Field := range ab84Fields {
		if ab84Field.Ab84Key == ab84Key {
			return ab84Field, true
		}
	}
	return Ab84BurnerField{}, false
}
