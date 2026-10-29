package vaillant

import (
	_ "embed"
	ab35json "encoding/json"
)

//go:embed data/ab35_legacy_timer.json
var ab35LegacyTimer []byte

func Ab35LegacyTimerCatalog() ab35json.RawMessage {
	return append(ab35json.RawMessage(nil), ab35LegacyTimer...)
}
