package vaillant

import (
	_ "embed"
	"encoding/json"
)

//go:embed data/ab27_holiday.json
var ab27HolidayJSON []byte

func Ab27HolidayCatalog() json.RawMessage {
	return append(json.RawMessage(nil), ab27HolidayJSON...)
}
