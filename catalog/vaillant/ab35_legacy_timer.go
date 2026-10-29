package vaillant

import (
	_ "embed"
	ab35json "encoding/json"
)

//go:embed data/ab35_legacy_timer.json
var ab35LegacyTimer []byte

// Ab35LegacyTimerCatalog returns an independently owned copy of the opt-in B515
// candidate catalog. It grants no runtime registration or write admission.
// Public contract: https://github.com/Project-Helianthus/helianthus-docs-ebus/blob/main/development/ab35-b515-catalog.md
func Ab35LegacyTimerCatalog() ab35json.RawMessage {
	return append(ab35json.RawMessage(nil), ab35LegacyTimer...)
}
