package owner

import (
	"sync"

	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/output"
)

const semanticTaskRefreshCatalogues semanticTaskKey = "refresh_catalogues"

var (
	familyCatalogueMu sync.Mutex
	familyCatalogues  = map[*vaillantSemanticPoller]output.FamilyCatalogue{}
)

func rememberFamilyCatalogue(p *vaillantSemanticPoller, view output.FamilyCatalogue) {
	if p == nil {
		return
	}
	familyCatalogueMu.Lock()
	familyCatalogues[p] = view
	familyCatalogueMu.Unlock()
}

func reapplyStoredFamilyCatalogue(p *vaillantSemanticPoller) {
	if p == nil {
		return
	}
	familyCatalogueMu.Lock()
	view, ok := familyCatalogues[p]
	familyCatalogueMu.Unlock()
	if !ok {
		return
	}
	publishFamilyCatalogue(p, view)
}
