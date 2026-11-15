package vaillant

import (
	_ "embed"
	"encoding/json"
)

//go:embed data/holiday.json
var holidayHolidayJSON []byte

func HolidayCatalog() json.RawMessage {
	return append(json.RawMessage(nil), holidayHolidayJSON...)
}
