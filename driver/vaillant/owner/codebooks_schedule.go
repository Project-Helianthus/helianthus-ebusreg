package owner

import (
	"context"
	"errors"
	"log"
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
	var errs []error
	if boiler != 0 {
		if _, err := fillBoilerDictionary(p, ctx, boiler); err != nil {
			errs = append(errs, err)
		}
	}
	if _, err := refreshEnergyCatalogues(p, ctx, InstallationSessionFromContext(ctx)); err != nil {
		errs = append(errs, err)
	}
	joined := errors.Join(errs...)
	if state := p.catalogueState(); state != nil {
		state.recordFillError(joined)
	}
	if joined != nil {
		log.Printf("semantic catalogue refresh failed key=%s err=%v", semanticTaskRefreshCatalogues, joined)
	}
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

func (p *vaillantSemanticPoller) catalogueState() *instanceCatalogues {
	owner := publicationOwner(p)
	if owner == nil || owner.state == nil {
		return nil
	}
	return owner.state.catalogueState()
}

func (p *vaillantSemanticPoller) releaseCatalogueState() {
	owner := publicationOwner(p)
	if owner == nil || owner.state == nil {
		return
	}
	owner.state.releaseCatalogueState()
}

func (p *vaillantSemanticPoller) catalogueStateReleased() bool {
	owner := publicationOwner(p)
	if owner == nil || owner.state == nil {
		return false
	}
	return owner.state.catalogueStateReleased()
}

func (p *vaillantSemanticPoller) catalogueFillError() error {
	state := p.catalogueState()
	if state == nil {
		return nil
	}
	state.mu.Lock()
	defer state.mu.Unlock()
	return state.fillErr
}

func (p *vaillantSemanticPoller) reapplyFamilyCatalogue() {
	reapplyStoredFamilyCatalogue(p)
}

func (p *vaillantSemanticPoller) reapplyBoilerCatalogue() {
	reapplyStoredBoilerCatalogue(p)
}
