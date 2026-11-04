package vaillant

import (
	_ "embed"
	"encoding/json"
)

type Ab88Axis struct {
	Ab88Present bool   `json:"present"`
	Ab88Value   uint16 `json:"value"`
}

type Ab88AirKey struct {
	Ab88Family    string   `json:"family"`
	Ab88Variant   string   `json:"variant"`
	Ab88Primary   uint16   `json:"primary"`
	Ab88Secondary Ab88Axis `json:"secondary"`
	Ab88Operation Ab88Axis `json:"operation"`
	Ab88Direction string   `json:"direction"`
	Ab88Field     string   `json:"field"`
	Ab88Opcode    Ab88Axis `json:"opcode"`
	Ab88Group     Ab88Axis `json:"group"`
	Ab88Instance  Ab88Axis `json:"instance"`
	Ab88Register  Ab88Axis `json:"register"`
}

type Ab88AirField struct {
	Ab88Key             Ab88AirKey                   `json:"key"`
	Ab88LogicalPosition *int                         `json:"logical_position"`
	Ab88Width           *int                         `json:"width_bytes"`
	Ab88Signedness      string                       `json:"signedness"`
	Ab88ByteOrder       string                       `json:"byte_order"`
	Ab88Conversion      string                       `json:"conversion"`
	Ab88Divisor         float64                      `json:"divisor"`
	Ab88Unit            *string                      `json:"unit"`
	Ab88Default         *float64                     `json:"default"`
	Ab88DefaultLabel    *string                      `json:"default_label"`
	Ab88Bounds          []float64                    `json:"bounds"`
	Ab88Enumeration     map[string]string            `json:"enumeration"`
	Ab88Otherwise       *string                      `json:"otherwise"`
	Ab88Bits            map[string]map[string]string `json:"bits"`
	Ab88ReservedBits    []int                        `json:"reserved_bits"`
	Ab88Replacement     *string                      `json:"replacement"`
	Ab88Grade           string                       `json:"grade"`
	Ab88RawFallback     string                       `json:"raw_fallback"`
}

type Ab88AirOperation struct {
	Ab88Operation        uint16 `json:"operation"`
	Ab88RequestArguments bool   `json:"request_arguments"`
	Ab88Grade            string `json:"grade"`
}

type Ab88AirMetadata struct {
	Ab88Representation  string             `json:"representation"`
	Ab88PhysicalSupport string             `json:"physical_support"`
	Ab88Acknowledgement string             `json:"acknowledgement"`
	Ab88Complete        bool               `json:"complete"`
	Ab88Coordinates     string             `json:"coordinates"`
	Ab88Limitations     []string           `json:"limitations"`
	Ab88Operations      []Ab88AirOperation `json:"operations"`
	Ab88Fields          []Ab88AirField     `json:"fields"`
}

//go:embed data/ab88_b512_air_fields.json
var ab88AirJSON []byte

func Ab88B512AirFields() Ab88AirMetadata {
	var ab88Metadata Ab88AirMetadata
	if ab88Error := json.Unmarshal(ab88AirJSON, &ab88Metadata); ab88Error != nil {
		panic(ab88Error)
	}
	return ab88Metadata
}

func ab88FindAirField(ab88Fields []Ab88AirField, ab88Key Ab88AirKey) (Ab88AirField, bool) {
	for _, ab88Field := range ab88Fields {
		if ab88Field.Ab88Key == ab88Key {
			return ab88Field, true
		}
	}
	return Ab88AirField{}, false
}

func Ab88LookupB512AirField(ab88Key Ab88AirKey) (Ab88AirField, bool) {
	return ab88FindAirField(Ab88B512AirFields().Ab88Fields, ab88Key)
}
