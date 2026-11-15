package vaillant

import (
	_ "embed"
	b520BlockJSON "encoding/json"
	b520BlockReflect "reflect"
)

//go:embed data/b520_blocks.json
var b520BlockB520Data []byte

// B520BlockKey identifies one B520 block field; nil optional axes are absent and lookups compare every
// axis exactly.
type B520BlockKey struct {
	B520BlockFamily      string  `json:"family"`
	B520BlockVariant     string  `json:"variant"`
	B520BlockPrimary     uint16  `json:"primary"`
	B520BlockSecondary   *uint16 `json:"secondary"`
	B520BlockOperation   *uint16 `json:"operation"`
	B520BlockDirection   string  `json:"direction"`
	B520BlockAlternative string  `json:"alternative"`
	B520BlockField       string  `json:"field"`
	B520BlockOpcode      *uint16 `json:"opcode"`
	B520BlockGroup       *uint16 `json:"group"`
	B520BlockInstance    *uint16 `json:"instance"`
	B520BlockRegister    *uint16 `json:"register"`
}

// B520BlockB520Field declares one block field with its prefix, layout, conversion, unit, limits, labels,
// grade and the raw fallback; nil members are not declared.
type B520BlockB520Field struct {
	B520BlockKey             B520BlockKey      `json:"key"`
	B520BlockPrefix          []int             `json:"prefix"`
	B520BlockLogicalPosition *int              `json:"logical_position"`
	B520BlockWidthBytes      *int              `json:"width_bytes"`
	B520BlockSignedness      string            `json:"signedness"`
	B520BlockByteOrder       string            `json:"byte_order"`
	B520BlockConversion      string            `json:"conversion"`
	B520BlockUnit            *string           `json:"unit"`
	B520BlockDefault         *uint64           `json:"default"`
	B520BlockMinimum         *uint64           `json:"minimum"`
	B520BlockMaximum         *uint64           `json:"maximum"`
	B520BlockMinimumBytes    *int              `json:"minimum_bytes"`
	B520BlockMaximumBytes    *int              `json:"maximum_bytes"`
	B520BlockExtent          *string           `json:"extent"`
	B520BlockEnumeration     map[string]string `json:"enumeration"`
	B520BlockReplacement     *string           `json:"replacement"`
	B520BlockReserved        bool              `json:"reserved"`
	B520BlockGrade           string            `json:"grade"`
	B520BlockSemantics       string            `json:"semantics"`
	B520BlockRawFallback     string            `json:"raw_fallback"`
}

// B520BlockCatalog is the B520 block catalog with its catalog-wide grades and declared unknowns.
type B520BlockCatalog struct {
	B520BlockRepresentationGrade   string               `json:"representation_grade"`
	B520BlockPhysicalQualification string               `json:"physical_qualification"`
	B520BlockFamilyCompleteness    string               `json:"family_completeness"`
	B520BlockInvocation            string               `json:"invocation"`
	B520BlockSelection             string               `json:"selection"`
	B520BlockNegativeForm          string               `json:"negative_form"`
	B520BlockLogicalPositions      string               `json:"logical_positions"`
	B520BlockUnknowns              map[string]string    `json:"unknowns"`
	B520BlockFields                []B520BlockB520Field `json:"fields"`
}

// B520Blocks decodes the embedded catalog on each call and returns a detached copy or the decoding error.
func B520Blocks() (B520BlockCatalog, error) {
	var b520BlockCatalog B520BlockCatalog
	b520BlockError := b520BlockJSON.Unmarshal(b520BlockB520Data, &b520BlockCatalog)
	return b520BlockCatalog, b520BlockError
}

func b520BlockSameKey(b520BlockLeft, b520BlockRight B520BlockKey) bool {
	return b520BlockReflect.DeepEqual(b520BlockLeft, b520BlockRight)
}

// B520BlockLookup returns the field whose key equals the given key on every axis, false when none matches,
// or the catalog decoding error.
func B520BlockLookup(b520BlockKey B520BlockKey) (B520BlockB520Field, bool, error) {
	b520BlockCatalog, b520BlockError := B520Blocks()
	if b520BlockError != nil {
		return B520BlockB520Field{}, false, b520BlockError
	}
	for _, b520BlockField := range b520BlockCatalog.B520BlockFields {
		if b520BlockSameKey(b520BlockField.B520BlockKey, b520BlockKey) {
			return b520BlockField, true, nil
		}
	}
	return B520BlockB520Field{}, false, nil
}
