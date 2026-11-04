package vaillant

import (
	_ "embed"
	ab92JSON "encoding/json"
	ab92Reflect "reflect"
)

//go:embed data/ab92_b520_blocks.json
var ab92B520Data []byte

type Ab92Key struct {
	Ab92Family      string  `json:"family"`
	Ab92Variant     string  `json:"variant"`
	Ab92Primary     uint16  `json:"primary"`
	Ab92Secondary   *uint16 `json:"secondary"`
	Ab92Operation   *uint16 `json:"operation"`
	Ab92Direction   string  `json:"direction"`
	Ab92Alternative string  `json:"alternative"`
	Ab92Field       string  `json:"field"`
	Ab92Opcode      *uint16 `json:"opcode"`
	Ab92Group       *uint16 `json:"group"`
	Ab92Instance    *uint16 `json:"instance"`
	Ab92Register    *uint16 `json:"register"`
}

type Ab92B520Field struct {
	Ab92Key             Ab92Key           `json:"key"`
	Ab92Prefix          []int             `json:"prefix"`
	Ab92LogicalPosition *int              `json:"logical_position"`
	Ab92WidthBytes      *int              `json:"width_bytes"`
	Ab92Signedness      string            `json:"signedness"`
	Ab92ByteOrder       string            `json:"byte_order"`
	Ab92Conversion      string            `json:"conversion"`
	Ab92Unit            *string           `json:"unit"`
	Ab92Default         *uint64           `json:"default"`
	Ab92Minimum         *uint64           `json:"minimum"`
	Ab92Maximum         *uint64           `json:"maximum"`
	Ab92MinimumBytes    *int              `json:"minimum_bytes"`
	Ab92MaximumBytes    *int              `json:"maximum_bytes"`
	Ab92Extent          *string           `json:"extent"`
	Ab92Enumeration     map[string]string `json:"enumeration"`
	Ab92Replacement     *string           `json:"replacement"`
	Ab92Reserved        bool              `json:"reserved"`
	Ab92Grade           string            `json:"grade"`
	Ab92Semantics       string            `json:"semantics"`
	Ab92RawFallback     string            `json:"raw_fallback"`
}

type Ab92B520Catalog struct {
	Ab92RepresentationGrade   string            `json:"representation_grade"`
	Ab92PhysicalQualification string            `json:"physical_qualification"`
	Ab92FamilyCompleteness    string            `json:"family_completeness"`
	Ab92Invocation            string            `json:"invocation"`
	Ab92Selection             string            `json:"selection"`
	Ab92NegativeForm          string            `json:"negative_form"`
	Ab92LogicalPositions      string            `json:"logical_positions"`
	Ab92Unknowns              map[string]string `json:"unknowns"`
	Ab92Fields                []Ab92B520Field   `json:"fields"`
}

func Ab92B520Blocks() (Ab92B520Catalog, error) {
	var ab92Catalog Ab92B520Catalog
	ab92Error := ab92JSON.Unmarshal(ab92B520Data, &ab92Catalog)
	return ab92Catalog, ab92Error
}

func ab92SameKey(ab92Left, ab92Right Ab92Key) bool {
	return ab92Reflect.DeepEqual(ab92Left, ab92Right)
}

func Ab92LookupB520(ab92Key Ab92Key) (Ab92B520Field, bool, error) {
	ab92Catalog, ab92Error := Ab92B520Blocks()
	if ab92Error != nil {
		return Ab92B520Field{}, false, ab92Error
	}
	for _, ab92Field := range ab92Catalog.Ab92Fields {
		if ab92SameKey(ab92Field.Ab92Key, ab92Key) {
			return ab92Field, true, nil
		}
	}
	return Ab92B520Field{}, false, nil
}
