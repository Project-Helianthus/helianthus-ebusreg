package vaillant

import (
	_ "embed"
	ab86json "encoding/json"
	ab86reflect "reflect"
)

type Ab86Key struct {
	Ab86Family    string  `json:"family"`
	Ab86Variant   *string `json:"variant"`
	Ab86Primary   uint8   `json:"primary"`
	Ab86Secondary *uint8  `json:"secondary"`
	Ab86Operation *uint8  `json:"operation"`
	Ab86Direction *string `json:"direction"`
	Ab86Field     string  `json:"field"`
	Ab86Opcode    *uint8  `json:"opcode"`
	Ab86Group     *uint8  `json:"group"`
	Ab86Instance  *uint8  `json:"instance"`
	Ab86Register  *uint16 `json:"register"`
}

type Ab86HistoryField struct {
	Ab86Key             Ab86Key        `json:"key"`
	Ab86LogicalPosition *int           `json:"logical_position"`
	Ab86WidthBits       *int           `json:"width_bits"`
	Ab86Signedness      string         `json:"signedness"`
	Ab86ByteOrder       string         `json:"byte_order"`
	Ab86Conversion      string         `json:"conversion"`
	Ab86Unit            *string        `json:"unit"`
	Ab86Default         *int           `json:"default"`
	Ab86Bounds          []int          `json:"bounds"`
	Ab86Enumeration     map[int]string `json:"enumeration"`
	Ab86Replacement     []int          `json:"replacement"`
	Ab86Reserved        []int          `json:"reserved"`
	Ab86Bits            map[int]string `json:"bits"`
	Ab86Grade           string         `json:"grade"`
	Ab86RawFallback     string         `json:"raw_fallback"`
	Ab86Branches        map[int]string `json:"branches"`
	Ab86NoItemTail      []int          `json:"no_item_tail"`
	Ab86Otherwise       *string        `json:"otherwise"`
}

type Ab86HistoryMetadata struct {
	Ab86Representation string             `json:"representation"`
	Ab86Qualification  string             `json:"qualification"`
	Ab86Unknown        []string           `json:"unknown"`
	Ab86Fields         []Ab86HistoryField `json:"fields"`
}

//go:embed data/ab86_b503_history_flags.json
var ab86HistoryJSON []byte

func Ab86HistoryCatalog() (Ab86HistoryMetadata, error) {
	var ab86Catalog Ab86HistoryMetadata
	ab86Error := ab86json.Unmarshal(ab86HistoryJSON, &ab86Catalog)
	return ab86Catalog, ab86Error
}

func ab86Lookup(ab86Fields []Ab86HistoryField, ab86Key Ab86Key) (Ab86HistoryField, bool) {
	for _, ab86Field := range ab86Fields {
		if ab86reflect.DeepEqual(ab86Field.Ab86Key, ab86Key) {
			return ab86Field, true
		}
	}
	return Ab86HistoryField{}, false
}

func Ab86HistoryLookup(ab86Key Ab86Key) (Ab86HistoryField, bool, error) {
	ab86Catalog, ab86Error := Ab86HistoryCatalog()
	if ab86Error != nil {
		return Ab86HistoryField{}, false, ab86Error
	}
	ab86Field, ab86Found := ab86Lookup(ab86Catalog.Ab86Fields, ab86Key)
	return ab86Field, ab86Found, nil
}
