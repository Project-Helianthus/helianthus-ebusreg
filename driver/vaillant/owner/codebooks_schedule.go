package owner

import (
	"context"
	"strings"
)

// refreshPublishedCatalogues is the recurring acquisition that fills family
// codebook labels, boiler selector evidence and catalogued energy directions
// onto the views the owner already publishes.
func (p *vaillantSemanticPoller) refreshPublishedCatalogues(ctx context.Context) {
	if !p.publishing() {
		return
	}
	if ctx == nil {
		ctx = context.Background()
	}
	fillFamilyCodebooks(p, ctx, p.codebookFamily())
	p.mu.Lock()
	boiler := p.boilerAddress
	p.mu.Unlock()
	if boiler != 0 {
		_, _ = fillBoilerDictionary(p, ctx, boiler)
	}
	_, _ = refreshEnergyCatalogues(p, ctx, InstallationSessionFromContext(ctx))
}

// codebookFamily is the registered device id when that id is the family name.
// The codebook catalogue selects a family explicitly and does not infer one from a reply.
func (p *vaillantSemanticPoller) codebookFamily() string {
	if p == nil || p.reg == nil {
		return ""
	}
	p.mu.Lock()
	controller := p.controller
	p.mu.Unlock()
	if controller == 0 {
		return ""
	}
	snapshot, ok := p.reg.LookupEntrySnapshot(controller)
	if !ok {
		return ""
	}
	return strings.ToLower(strings.TrimSpace(snapshot.DeviceID))
}

func (p *vaillantSemanticPoller) reapplyFamilyCatalogue() {
	reapplyStoredFamilyCatalogue(p)
}

func (p *vaillantSemanticPoller) reapplyBoilerCatalogue() {
	reapplyStoredBoilerCatalogue(p)
}
