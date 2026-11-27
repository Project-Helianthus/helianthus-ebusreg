package owner

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"sync"
	"time"

	semdriver "github.com/Project-Helianthus/helianthus-semreg/semreg/v1/driver"

	catalogvaillant "github.com/Project-Helianthus/helianthus-ebusreg/catalog/vaillant"
	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/output"
)

type energyCatalogueHold struct {
	totals output.EnergyTotals
	ready  bool
}

var (
	energyCatalogueMu sync.Mutex
	energyCatalogues  = map[*vaillantSemanticPoller]*energyCatalogueHold{}
	energyResidualMu  sync.Mutex
	energyResiduals   = map[*vaillantSemanticPoller]map[string][]byte{}
)

func noteEnergyResidual(p *vaillantSemanticPoller, field string, raw []byte) {
	energyResidualMu.Lock()
	defer energyResidualMu.Unlock()
	noted := energyResiduals[p]
	if noted == nil {
		noted = map[string][]byte{}
		energyResiduals[p] = noted
	}
	noted[field] = append([]byte(nil), raw...)
}

func fillEnergyCatalogues(p *vaillantSemanticPoller, ctx context.Context, session semdriver.Session) (output.EnergyTotals, error) {
	energyCatalogueMu.Lock()
	if hold := energyCatalogues[p]; hold != nil && hold.ready {
		totals := hold.totals
		energyCatalogueMu.Unlock()
		return totals, nil
	}
	energyCatalogueMu.Unlock()

	var acquisition semdriver.CalendarAcquisition
	var reference semdriver.CalendarReference
	if session != nil && session.Calendar() != nil {
		acquisition = session.Calendar().BeginAcquisition()
		if acquisition != nil {
			reference = acquisition.Reference()
		}
	}
	ledger := p.energyLedger()
	directions := make([]output.EnergyDirection, 0, len(b524EnergyQueries))
	for _, query := range b524EnergyQueries {
		value, ok := p.readEnergyRegister(ctx, query)
		direction := output.EnergyDirection{
			Name:      energyRegisterName(query.addr),
			Period:    query.period,
			Source:    query.channel,
			Usage:     query.usage,
			Freshness: "Unknown",
			Level:     "Unknown",
		}
		if reference != nil {
			direction.ScopeYear, direction.ScopeMonth = energyScope(reference, query.yearKind)
		}
		if !ok {
			directions = append(directions, direction)
			continue
		}
		raw := make([]byte, 4)
		binary.LittleEndian.PutUint32(raw, value)
		direction.Raw = raw
		direction.ObservedAt = p.now().UTC().Format(time.RFC3339Nano)
		if query.period == "month" {
			if !energyResponseAdmitted(session, acquisition, ledger, reference) {
				directions = append(directions, direction)
				if ledger != nil && reference != nil {
					now := ledger.energyReferenceTime()
					if now.Year() != reference.Instant().Year() || now.Month() != reference.Instant().Month() {
						break
					}
				}
				continue
			}
		}
		kwh := float64(value)
		direction.ValueKWh = &kwh
		if ledger != nil {
			if query.period == "all" {
				ledger.applyLifetimeEnergy(query.channel, query.usage, kwh)
			} else if reference != nil {
				ledger.applyWindowEnergy(query.channel, query.usage, query.period, query.yearKind, reference.Instant(), kwh)
			}
		}
		directions = append(directions, direction)
	}
	catalogue, err := catalogueEnergyDirections()
	if err != nil {
		return output.EnergyTotals{}, err
	}
	residuals, err := catalogueEnergyResiduals(p)
	if err != nil {
		return output.EnergyTotals{}, err
	}
	candidates, err := catalogueEnergyCandidates()
	if err != nil {
		return output.EnergyTotals{}, err
	}
	totals := output.EnergyTotals{}
	if publication := publicationOwner(p); publication != nil && publication.state != nil {
		if published := publication.state.EnergyTotals(); published != nil {
			totals.Gas = published.Gas
			totals.Electric = published.Electric
			totals.Solar = published.Solar
		}
	}
	totals.SetRegisterDirections(directions)
	totals.SetCatalogueDirections(catalogue)
	totals.SetResidualFields(residuals)
	totals.SetCandidateFields(candidates)
	energyCatalogueMu.Lock()
	energyCatalogues[p] = &energyCatalogueHold{totals: totals, ready: true}
	energyCatalogueMu.Unlock()
	if ledger != nil {
		ledger.publishEnergyTotals()
	}
	return totals, nil
}

// overlayEnergyCatalogue copies the stored catalogue rows onto a rebuilt
// energy view. Channel totals stay the values the merge store just projected.
func overlayEnergyCatalogue(provider *retainedSemanticState, totals *output.EnergyTotals) {
	if provider == nil || totals == nil {
		return
	}
	energyCatalogueMu.Lock()
	defer energyCatalogueMu.Unlock()
	for poller, hold := range energyCatalogues {
		if hold == nil || !hold.ready {
			continue
		}
		owner := publicationOwner(poller)
		if owner == nil || owner.state != provider {
			continue
		}
		totals.SetRegisterDirections(hold.totals.RegisterDirections())
		totals.SetCatalogueDirections(hold.totals.CatalogueDirections())
		totals.SetResidualFields(hold.totals.ResidualFields())
		totals.SetCandidateFields(hold.totals.CandidateFields())
		return
	}
}

func energyResponseAdmitted(session semdriver.Session, acquisition semdriver.CalendarAcquisition, ledger vaillantEnergyLedger, reference semdriver.CalendarReference) bool {
	if session == nil || session.Calendar() == nil || acquisition == nil || reference == nil {
		return false
	}
	if acquisition.Check(session.Calendar()) != nil {
		return false
	}
	if ledger != nil {
		now := ledger.energyReferenceTime()
		if now.Year() != reference.Instant().Year() || now.Month() != reference.Instant().Month() {
			return false
		}
	}
	return true
}

func energyScope(reference semdriver.CalendarReference, label string) (int, int) {
	if label == "previous" {
		month := reference.PreviousMonth()
		return month.Year, int(month.Month)
	}
	if label == "current" {
		month := reference.Month()
		return month.Year, int(month.Month)
	}
	return reference.Year(), 0
}

func energyRegisterName(addr uint16) string {
	switch addr {
	case energyFuelSumHc:
		return "fuel_sum_heating"
	case energyFuelSumHwc:
		return "fuel_sum_hot_water"
	case energyElectricitySumHc:
		return "electricity_sum_heating"
	case energyElectricitySumHwc:
		return "electricity_sum_hot_water"
	case energyFuelSumHcThisMonth:
		return "fuel_sum_heating_this_month"
	case energyFuelSumHwcThisMonth:
		return "fuel_sum_hot_water_this_month"
	case energyElectricitySumHcThisMonth:
		return "electricity_sum_heating_this_month"
	case energyElectricitySumHwcThisMonth:
		return "electricity_sum_hot_water_this_month"
	case energyFuelSumHcLastMonth:
		return "fuel_sum_heating_last_month"
	case energyFuelSumHwcLastMonth:
		return "fuel_sum_hot_water_last_month"
	case energyElectricitySumHcLastMonth:
		return "electricity_sum_heating_last_month"
	case energyElectricitySumHwcLastMonth:
		return "electricity_sum_hot_water_last_month"
	default:
		return ""
	}
}

func catalogueEnergyDirections() ([]output.EnergyDirection, error) {
	entries, err := catalogvaillant.B516StatisticsDirectionEntries()
	if err != nil {
		return nil, err
	}
	directions := make([]output.EnergyDirection, 0, len(entries))
	for _, entry := range entries {
		key := entry.B516StatisticsDirectionKey
		directions = append(directions, output.EnergyDirection{
			Name:      key.B516StatisticsDirectionField,
			Family:    key.B516StatisticsDirectionFamily,
			Variant:   key.B516StatisticsDirectionVariant,
			Field:     key.B516StatisticsDirectionField,
			Direction: key.B516StatisticsDirection,
			Freshness: entry.B516StatisticsDirectionFreshness,
			Level:     entry.B516StatisticsDirectionGrade,
		})
	}
	raw, err := catalogvaillant.B516DirectionalTypeEntries()
	if err != nil {
		return nil, err
	}
	var typed []struct {
		Key struct {
			Family    string `json:"family"`
			Variant   string `json:"variant"`
			Field     string `json:"field"`
			Direction string `json:"direction"`
		} `json:"key"`
		Grade string `json:"grade"`
	}
	if err := json.Unmarshal(raw, &typed); err != nil {
		return nil, err
	}
	for _, entry := range typed {
		directions = append(directions, output.EnergyDirection{
			Name:      entry.Key.Field,
			Family:    entry.Key.Family,
			Variant:   entry.Key.Variant,
			Field:     entry.Key.Field,
			Direction: entry.Key.Direction,
			Freshness: "Unknown",
			Level:     entry.Grade,
			Label:     catalogueUnresolved,
		})
	}
	return directions, nil
}

func catalogueEnergyResiduals(p *vaillantSemanticPoller) ([]output.EnergyResidual, error) {
	raw, err := catalogvaillant.B516ResidualFieldEntries()
	if err != nil {
		return nil, err
	}
	var rows []struct {
		Key struct {
			Family  string `json:"family"`
			Variant string `json:"variant"`
			Field   string `json:"field"`
		} `json:"key"`
		Grade string `json:"grade"`
	}
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, err
	}
	energyResidualMu.Lock()
	noted := energyResiduals[p]
	energyResidualMu.Unlock()
	residuals := make([]output.EnergyResidual, 0, len(rows))
	for _, row := range rows {
		item := output.EnergyResidual{
			Name:    row.Key.Field,
			Family:  row.Key.Family,
			Variant: row.Key.Variant,
			Field:   row.Key.Field,
			Level:   row.Grade,
		}
		if raw, ok := noted[row.Key.Field]; ok {
			item.Raw = append([]byte(nil), raw...)
		}
		residuals = append(residuals, item)
	}
	return residuals, nil
}

func catalogueEnergyCandidates() ([]output.EnergyResidual, error) {
	entries, err := catalogvaillant.B507EnergyCandidateEntries()
	if err != nil {
		return nil, err
	}
	candidates := make([]output.EnergyResidual, 0, len(entries))
	for _, entry := range entries {
		key := entry.B507EnergyCandidateKey
		item := output.EnergyResidual{
			Name:     key.B507EnergyCandidateField,
			Family:   key.B507EnergyCandidateFamily,
			Variant:  key.B507EnergyCandidateVariant,
			Field:    key.B507EnergyCandidateField,
			Level:    entry.B507EnergyCandidateGrade,
			Register: key.B507EnergyCandidateRegister,
			Note:     "register absent",
		}
		if key.B507EnergyCandidateRegister != nil {
			item.Note = ""
		}
		candidates = append(candidates, item)
	}
	return candidates, nil
}
