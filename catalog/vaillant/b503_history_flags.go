package vaillant

import (
	_ "embed"
	b503HistoryFlagJSON "encoding/json"
	b503HistoryFlagReflect "reflect"
)

type B503HistoryFlagKey struct {
	B503HistoryFlagFamily    string  `json:"family"`
	B503HistoryFlagVariant   *string `json:"variant"`
	B503HistoryFlagPrimary   uint8   `json:"primary"`
	B503HistoryFlagSecondary *uint8  `json:"secondary"`
	B503HistoryFlagOperation *uint8  `json:"operation"`
	B503HistoryFlagDirection *string `json:"direction"`
	B503HistoryFlagField     string  `json:"field"`
	B503HistoryFlagOpcode    *uint8  `json:"opcode"`
	B503HistoryFlagGroup     *uint8  `json:"group"`
	B503HistoryFlagInstance  *uint8  `json:"instance"`
	B503HistoryFlagRegister  *uint16 `json:"register"`
}

type B503HistoryField struct {
	B503HistoryFlagKey             B503HistoryFlagKey `json:"key"`
	B503HistoryFlagLogicalPosition *int               `json:"logical_position"`
	B503HistoryFlagWidthBits       *int               `json:"width_bits"`
	B503HistoryFlagSignedness      string             `json:"signedness"`
	B503HistoryFlagByteOrder       string             `json:"byte_order"`
	B503HistoryFlagConversion      string             `json:"conversion"`
	B503HistoryFlagUnit            *string            `json:"unit"`
	B503HistoryFlagDefault         *int               `json:"default"`
	B503HistoryFlagBounds          []int              `json:"bounds"`
	B503HistoryFlagEnumeration     map[int]string     `json:"enumeration"`
	B503HistoryFlagReplacement     []int              `json:"replacement"`
	B503HistoryFlagReserved        []int              `json:"reserved"`
	B503HistoryFlagBits            map[int]string     `json:"bits"`
	B503HistoryFlagGrade           string             `json:"grade"`
	B503HistoryFlagRawFallback     string             `json:"raw_fallback"`
	B503HistoryFlagBranches        map[int]string     `json:"branches"`
	B503HistoryFlagNoItemTail      []int              `json:"no_item_tail"`
	B503HistoryFlagOtherwise       *string            `json:"otherwise"`
}

type B503HistoryMetadata struct {
	B503HistoryFlagRepresentation string             `json:"representation"`
	B503HistoryFlagQualification  string             `json:"qualification"`
	B503HistoryFlagUnknown        []string           `json:"unknown"`
	B503HistoryFlagFields         []B503HistoryField `json:"fields"`
}

//go:embed data/ab86_b503_history_flags.json
var b503HistoryFlagHistoryJSON []byte

func B503HistoryCatalog() (B503HistoryMetadata, error) {
	var b503HistoryFlagCatalog B503HistoryMetadata
	b503HistoryFlagError := b503HistoryFlagJSON.Unmarshal(b503HistoryFlagHistoryJSON, &b503HistoryFlagCatalog)
	return b503HistoryFlagCatalog, b503HistoryFlagError
}

func b503HistoryFlagLookup(b503HistoryFlagFields []B503HistoryField, b503HistoryFlagKey B503HistoryFlagKey) (B503HistoryField, bool) {
	for _, b503HistoryFlagField := range b503HistoryFlagFields {
		if b503HistoryFlagReflect.DeepEqual(b503HistoryFlagField.B503HistoryFlagKey, b503HistoryFlagKey) {
			return b503HistoryFlagField, true
		}
	}
	return B503HistoryField{}, false
}

func B503HistoryLookup(b503HistoryFlagKey B503HistoryFlagKey) (B503HistoryField, bool, error) {
	b503HistoryFlagCatalog, b503HistoryFlagError := B503HistoryCatalog()
	if b503HistoryFlagError != nil {
		return B503HistoryField{}, false, b503HistoryFlagError
	}
	b503HistoryFlagField, b503HistoryFlagFound := b503HistoryFlagLookup(b503HistoryFlagCatalog.B503HistoryFlagFields, b503HistoryFlagKey)
	return b503HistoryFlagField, b503HistoryFlagFound, nil
}
