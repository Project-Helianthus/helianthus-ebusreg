package vaillant

import (
	_ "embed"
	ab90JSON "encoding/json"
	ab90Reflect "reflect"
)

//go:embed data/ab90_b51a_mirror_fields.json
var ab90MirrorJSON string

type Ab90MirrorKey struct {
	Ab90Family    string  `json:"family"`
	Ab90Variant   string  `json:"variant"`
	Ab90Primary   uint8   `json:"primary"`
	Ab90Secondary *uint8  `json:"secondary"`
	Ab90Operation *uint8  `json:"operation"`
	Ab90Direction string  `json:"direction"`
	Ab90Field     string  `json:"field"`
	Ab90Opcode    *uint8  `json:"opcode"`
	Ab90Group     *uint8  `json:"group"`
	Ab90Instance  *uint8  `json:"instance"`
	Ab90Register  *uint16 `json:"register"`
}

type Ab90MirrorBit struct {
	Ab90Index   int       `json:"index"`
	Ab90Name    *string   `json:"name"`
	Ab90Role    string    `json:"role"`
	Ab90Default *bool     `json:"default"`
	Ab90Grade   Ab05Grade `json:"grade"`
}

type Ab90MirrorAlternative struct {
	Ab90Name        string           `json:"name"`
	Ab90WidthBytes  *int             `json:"width_bytes"`
	Ab90Condition   string           `json:"condition"`
	Ab90Enumeration map[int64]string `json:"enumeration"`
	Ab90Bits        []Ab90MirrorBit  `json:"bits"`
	Ab90Grade       Ab05Grade        `json:"grade"`
}

type Ab90MirrorField struct {
	Ab90Key             Ab90MirrorKey           `json:"key"`
	Ab90LogicalPosition *int                    `json:"logical_position"`
	Ab90WidthBytes      *int                    `json:"width_bytes"`
	Ab90Signedness      string                  `json:"signedness"`
	Ab90ByteOrder       string                  `json:"byte_order"`
	Ab90Conversion      string                  `json:"conversion"`
	Ab90Unit            *string                 `json:"unit"`
	Ab90Default         *int64                  `json:"default"`
	Ab90FormalMin       *int64                  `json:"formal_min"`
	Ab90FormalMax       *int64                  `json:"formal_max"`
	Ab90LogicalMin      *int64                  `json:"logical_min"`
	Ab90LogicalMax      *int64                  `json:"logical_max"`
	Ab90LogicalLimit    *string                 `json:"logical_limit"`
	Ab90Enumeration     map[int64]string        `json:"enumeration"`
	Ab90Bits            []Ab90MirrorBit         `json:"bits"`
	Ab90Alternatives    []Ab90MirrorAlternative `json:"alternatives"`
	Ab90Replacement     *string                 `json:"replacement"`
	Ab90Reserved        []int                   `json:"reserved"`
	Ab90Grade           Ab05Grade               `json:"grade"`
	Ab90RawFallback     string                  `json:"raw_fallback"`
	Ab90Correlation     string                  `json:"correlation"`
	Ab90Unmapped        string                  `json:"unmapped"`
}

type Ab90MirrorMetadata struct {
	Ab90RepresentationGrade Ab05Grade         `json:"representation_grade"`
	Ab90SemanticsGrade      Ab05Grade         `json:"semantics_grade"`
	Ab90PhysicalAdaptation  Ab05Grade         `json:"physical_adaptation"`
	Ab90StorageDurability   Ab05Grade         `json:"storage_durability"`
	Ab90SuccessSemantics    Ab05Grade         `json:"success_semantics"`
	Ab90Complete            bool              `json:"complete"`
	Ab90NativeContract      string            `json:"native_contract"`
	Ab90Limitations         []string          `json:"limitations"`
	Ab90Fields              []Ab90MirrorField `json:"fields"`
}

func Ab90B51AMirrorFields() (Ab90MirrorMetadata, error) {
	var ab90Metadata Ab90MirrorMetadata
	ab90Error := ab90JSON.Unmarshal([]byte(ab90MirrorJSON), &ab90Metadata)
	return ab90Metadata, ab90Error
}

func ab90MirrorKeysEqual(ab90Left, ab90Right Ab90MirrorKey) bool {
	return ab90Reflect.DeepEqual(ab90Left, ab90Right)
}

func Ab90LookupB51AMirrorField(ab90Key Ab90MirrorKey) (Ab90MirrorField, bool, error) {
	ab90Metadata, ab90Error := Ab90B51AMirrorFields()
	if ab90Error != nil {
		return Ab90MirrorField{}, false, ab90Error
	}
	for _, ab90Field := range ab90Metadata.Ab90Fields {
		if ab90MirrorKeysEqual(ab90Field.Ab90Key, ab90Key) {
			return ab90Field, true, nil
		}
	}
	return Ab90MirrorField{}, false, nil
}
