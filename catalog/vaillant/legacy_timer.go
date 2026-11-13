package vaillant

import (
	_ "embed"
	legacyTimerJSON "encoding/json"
)

//go:embed data/ab35_legacy_timer.json
var legacyTimerLegacyTimer []byte

// LegacyTimerCatalog returns an independently owned copy of the opt-in B515
// candidate catalog. It grants no runtime registration or write admission.
// Public contract: https://github.com/Project-Helianthus/helianthus-docs-ebus/blob/main/development/ab35-b515-catalog.md
func LegacyTimerCatalog() legacyTimerJSON.RawMessage {
	return append(legacyTimerJSON.RawMessage(nil), legacyTimerLegacyTimer...)
}
