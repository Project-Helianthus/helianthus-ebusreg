package owner

import (
	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/output"
)

const semanticTaskRefreshCatalogues semanticTaskKey = "refresh_catalogues"

func rememberFamilyCatalogue(p *vaillantSemanticPoller, view output.FamilyCatalogue) bool {
	return p.catalogueState().rememberFamily(view)
}

func reapplyStoredFamilyCatalogue(p *vaillantSemanticPoller) {
	view, ok := p.catalogueState().storedFamily()
	if !ok {
		return
	}
	publishFamilyCatalogue(p, view)
}
