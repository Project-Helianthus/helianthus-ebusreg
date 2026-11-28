package owner

import (
	"errors"
	"sync"

	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/output"
)

var errCatalogueStateReleased = errors.New("catalogue state is released")

type instanceCatalogues struct {
	mu        sync.Mutex
	released  bool
	family    output.FamilyCatalogue
	familySet bool
	boiler    boilerCatalogueHold
	boilerSet bool
	energy    energyCatalogueHold
	residuals map[string][]byte
	fillErr   error
}

func (provider *retainedSemanticState) catalogueState() *instanceCatalogues {
	if provider == nil {
		return nil
	}
	provider.mu.Lock()
	defer provider.mu.Unlock()
	if provider.catalogues != nil && provider.catalogues.released {
		return provider.catalogues
	}
	if provider.catalogues == nil {
		provider.catalogues = &instanceCatalogues{}
	}
	return provider.catalogues
}

func (provider *retainedSemanticState) releaseCatalogueState() {
	if provider == nil {
		return
	}
	provider.mu.Lock()
	defer provider.mu.Unlock()
	if provider.catalogues == nil {
		provider.catalogues = &instanceCatalogues{released: true}
		return
	}
	provider.catalogues.release()
}

func (provider *retainedSemanticState) catalogueStateReleased() bool {
	if provider == nil {
		return false
	}
	provider.mu.RLock()
	state := provider.catalogues
	provider.mu.RUnlock()
	return state.isReleased()
}

func (s *instanceCatalogues) isReleased() bool {
	if s == nil {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.released
}

func (s *instanceCatalogues) release() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.released = true
	s.family = output.FamilyCatalogue{}
	s.familySet = false
	s.boiler = boilerCatalogueHold{}
	s.boilerSet = false
	s.energy = energyCatalogueHold{}
	s.residuals = nil
	s.fillErr = nil
}

func (s *instanceCatalogues) recordFillError(err error) {
	if s == nil {
		return
	}
	s.mu.Lock()
	s.fillErr = err
	s.mu.Unlock()
}

func (s *instanceCatalogues) rememberFamily(view output.FamilyCatalogue) bool {
	if s == nil {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.released {
		return false
	}
	s.family = view
	s.familySet = true
	return true
}

func (s *instanceCatalogues) storedFamily() (output.FamilyCatalogue, bool) {
	if s == nil {
		return output.FamilyCatalogue{}, false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.released || !s.familySet {
		return output.FamilyCatalogue{}, false
	}
	return s.family, true
}

func (s *instanceCatalogues) noteResidual(field string, raw []byte) {
	if s == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.released {
		return
	}
	if s.residuals == nil {
		s.residuals = map[string][]byte{}
	}
	s.residuals[field] = append([]byte(nil), raw...)
}

func (s *instanceCatalogues) residual(field string) ([]byte, bool) {
	if s == nil {
		return nil, false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	raw, ok := s.residuals[field]
	if !ok {
		return nil, false
	}
	return append([]byte(nil), raw...), true
}

func (s *instanceCatalogues) readyEnergy() (output.EnergyTotals, bool) {
	if s == nil {
		return output.EnergyTotals{}, false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.released || !s.energy.ready {
		return output.EnergyTotals{}, false
	}
	return s.energy.totals, true
}

func (s *instanceCatalogues) storeEnergy(totals output.EnergyTotals) bool {
	if s == nil {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.released {
		return false
	}
	s.energy = energyCatalogueHold{totals: totals, ready: true}
	return true
}

func (s *instanceCatalogues) overlayEnergy(totals *output.EnergyTotals) {
	if s == nil || totals == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.released || !s.energy.ready {
		return
	}
	totals.SetRegisterDirections(s.energy.totals.RegisterDirections)
	totals.SetCatalogueDirections(s.energy.totals.CatalogueDirections)
	totals.SetResidualFields(s.energy.totals.ResidualFields)
	totals.SetCandidateFields(s.energy.totals.CandidateFields)
}

func (s *instanceCatalogues) storeBoilerSelectors(selectors []output.BoilerSelectorEvidence) ([]output.RawOperation, bool) {
	if s == nil {
		return nil, false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.released {
		return nil, false
	}
	s.boiler.selectors = append([]output.BoilerSelectorEvidence(nil), selectors...)
	s.boilerSet = true
	return append([]output.RawOperation(nil), s.boiler.operations...), true
}

func (s *instanceCatalogues) storedBoiler() ([]output.BoilerSelectorEvidence, []output.RawOperation, bool) {
	if s == nil {
		return nil, nil, false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.released || (!s.boilerSet && len(s.boiler.operations) == 0 && len(s.boiler.selectors) == 0) {
		return nil, nil, false
	}
	return append([]output.BoilerSelectorEvidence(nil), s.boiler.selectors...), append([]output.RawOperation(nil), s.boiler.operations...), true
}

func (s *instanceCatalogues) admitOperation(operation output.RawOperation, admitted bool) ([]output.BoilerSelectorEvidence, []output.RawOperation) {
	if s == nil {
		return nil, nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.released {
		return nil, nil
	}
	if admitted {
		s.boiler.operations = []output.RawOperation{operation}
	} else {
		s.boiler.operations = nil
	}
	s.boilerSet = true
	return append([]output.BoilerSelectorEvidence(nil), s.boiler.selectors...), append([]output.RawOperation(nil), s.boiler.operations...)
}

func (s *instanceCatalogues) admittedOperation(name string) (output.RawOperation, bool) {
	if s == nil {
		return output.RawOperation{}, false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.released {
		return output.RawOperation{}, false
	}
	for i := range s.boiler.operations {
		if s.boiler.operations[i].Name == name && s.boiler.operations[i].Namespace == "ebusreg" {
			return s.boiler.operations[i], true
		}
	}
	return output.RawOperation{}, false
}
