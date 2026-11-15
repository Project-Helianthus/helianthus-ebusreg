package vaillant

import (
	_ "embed"
	"encoding/json"
)

//go:embed data/holiday.json
var holidayHolidayJSON []byte

// HolidayCatalog returns a caller-owned copy of the embedded holiday catalog JSON.
func HolidayCatalog() json.RawMessage {
	return append(json.RawMessage(nil), holidayHolidayJSON...)
}
