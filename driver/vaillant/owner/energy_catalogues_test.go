package owner

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"

	catalogvaillant "github.com/Project-Helianthus/helianthus-ebusreg/catalog/vaillant"
	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/output"
)

func TestVaillantEnergyCataloguesFillEveryDirection(t *testing.T) {
	if len(b524EnergyQueries) != 12 {
		t.Fatalf("len(b524EnergyQueries) = %d; want 12", len(b524EnergyQueries))
	}
	receipt := time.Date(2024, time.March, 15, 8, 0, 0, 123456789, time.UTC).UTC().Format(time.RFC3339Nano)
	noon := energyTestReference{}.Instant().UTC().Format(time.RFC3339Nano)

	t.Run("same date", func(t *testing.T) {
		bus := &energyCatalogueBus{}
		poller, provider := prepareEnergyPoller(bus)
		calendar := &energyTestCalendar{day: 15, generation: 1}
		session := &energyCatalogueSession{calendar: calendar}
		supplied := []byte{10, 20}
		noteEnergyResidual(poller, "setpoint", supplied)
		view, err := fillEnergyCatalogues(poller, context.Background(), session)
		if err != nil {
			t.Fatalf("fill energy catalogues: %v", err)
		}
		assertEnergyReads(t, bus.reads)
		assertAcceptedEnergy(t, view, energyTestReference{}, receipt, noon)
		assertEnergyProjection(t, provider.EnergyTotals())
		assertEnergyProjection(t, &view)
		published := provider.EnergyTotals()
		if published == nil || !reflect.DeepEqual(published.RegisterDirections, view.RegisterDirections) || !reflect.DeepEqual(published.CatalogueDirections, view.CatalogueDirections) || !reflect.DeepEqual(published.ResidualFields, view.ResidualFields) || !reflect.DeepEqual(published.CandidateFields, view.CandidateFields) {
			t.Fatalf("published energy view = %+v; want the fill view's directions, residuals and candidates", published)
		}
		assertEnergyCatalogue(t, view, supplied)
		mapped := output.MCPEnergyTotalsFromGraphQL(&view)
		if mapped == nil || len(mapped.RegisterDirections) != len(view.RegisterDirections) || mapped.RegisterDirections[0].ObservedAt != receipt || len(mapped.CatalogueDirections) != len(view.CatalogueDirections) || len(mapped.CandidateFields) != len(view.CandidateFields) || len(mapped.ResidualFields) != len(view.ResidualFields) {
			t.Fatalf("mcp energy = %+v; want the same directions, residuals and candidates", mapped)
		}
		reads, begins := len(bus.reads), calendar.begins
		again, err := fillEnergyCatalogues(poller, context.Background(), session)
		if err != nil {
			t.Fatalf("second fill: %v", err)
		}
		if len(bus.reads) != reads || calendar.begins != begins {
			t.Fatalf("second view reads=%d begins=%d; a view read must not poll or begin again", len(bus.reads), calendar.begins)
		}
		if again.RegisterDirections[4].ObservedAt != receipt || again.RegisterDirections[4].ValueKWh == nil || *again.RegisterDirections[4].ValueKWh != 55 {
			t.Fatalf("second view monthly = %+v; the receipt and value must stay", again.RegisterDirections[4])
		}
	})

	t.Run("civil day", func(t *testing.T) {
		bus := &energyCatalogueBus{}
		poller, provider := prepareEnergyPoller(bus)
		calendar := &energyTestCalendar{day: 15, generation: 1}
		bus.advance = func() { calendar.day++ }
		view, err := fillEnergyCatalogues(poller, context.Background(), &energyCatalogueSession{calendar: calendar})
		if err != nil {
			t.Fatalf("fill energy catalogues: %v", err)
		}
		assertEnergyReads(t, bus.reads)
		if calendar.begins != 1 || calendar.day == 15 {
			t.Fatalf("begins=%d day=%d; the guard must see the moved civil day", calendar.begins, calendar.day)
		}
		assertDatedEnergyRejected(t, view, provider, receipt)
	})

	t.Run("generation", func(t *testing.T) {
		bus := &energyCatalogueBus{}
		poller, provider := prepareEnergyPoller(bus)
		calendar := &energyTestCalendar{day: 15, generation: 1}
		bus.advance = func() { calendar.generation++ }
		view, err := fillEnergyCatalogues(poller, context.Background(), &energyCatalogueSession{calendar: calendar})
		if err != nil {
			t.Fatalf("fill energy catalogues: %v", err)
		}
		assertEnergyReads(t, bus.reads)
		if calendar.generation == 1 {
			t.Fatalf("generation stayed 1; the guard was not given a new generation")
		}
		assertDatedEnergyRejected(t, view, provider, receipt)
	})

	t.Run("other calendar", func(t *testing.T) {
		bus := &energyCatalogueBus{}
		poller, provider := prepareEnergyPoller(bus)
		first := &energyTestCalendar{day: 15, generation: 1}
		second := &energyTestCalendar{day: 15, generation: 1}
		view, err := fillEnergyCatalogues(poller, context.Background(), &energyCatalogueSession{calendar: first, later: second})
		if err != nil {
			t.Fatalf("fill energy catalogues: %v", err)
		}
		assertEnergyReads(t, bus.reads)
		if first.begins != 1 || second.begins != 0 {
			t.Fatalf("begins first=%d second=%d; one acquisition must stay on the first calendar", first.begins, second.begins)
		}
		assertDatedEnergyRejected(t, view, provider, receipt)
	})

	t.Run("no session", func(t *testing.T) {
		bus := &energyCatalogueBus{}
		poller, provider := prepareEnergyPoller(bus)
		view, err := fillEnergyCatalogues(poller, context.Background(), nil)
		if err != nil {
			t.Fatalf("fill energy catalogues: %v", err)
		}
		assertEnergyReads(t, bus.reads)
		assertDatedEnergyRejected(t, view, provider, receipt)
	})
}

func TestVaillantEnergyCataloguesScheduledFillUsesInstallationSession(t *testing.T) {
	bus := &energyCatalogueBus{}
	poller, provider := prepareEnergyPoller(bus)
	calendar := &energyTestCalendar{day: 15, generation: 1}
	session := &energyCatalogueSession{calendar: calendar}
	ctx := ContextWithInstallationSession(context.Background(), session)
	stamp := time.Date(2024, time.March, 15, 8, 0, 0, 123456789, time.UTC)
	poller.nowFn = func() time.Time { return stamp }
	fire, stop := startCatalogueSchedule(t, poller, ctx)
	defer stop()

	fire()
	if calendar.begins != 1 {
		t.Fatalf("begins = %d; the scheduled fill must take one calendar acquisition", calendar.begins)
	}
	receipt := stamp.UTC().Format(time.RFC3339Nano)
	view := provider.EnergyTotals()
	monthly := energyDirectionNamed(view, energyRegisterName(energyFuelSumHcThisMonth))
	if monthly.ValueKWh == nil || *monthly.ValueKWh != 55 || monthly.ObservedAt != receipt || monthly.ScopeYear != 2024 || monthly.ScopeMonth != int(time.March) {
		t.Fatalf("scheduled monthly = %+v; want 55 at %s in March 2024", monthly, receipt)
	}
	if countEnergyReads(bus.reads, energyFuelSumHcThisMonth) != 1 {
		t.Fatalf("this-month reads = %d; want one acquisition", countEnergyReads(bus.reads, energyFuelSumHcThisMonth))
	}

	later := stamp.Add(time.Hour)
	poller.nowFn = func() time.Time { return later }
	bus.replace = map[uint16]uint32{energyFuelSumHcThisMonth: 56}
	bus.advance = func() { calendar.day++ }
	bus.advanced = false
	fire()
	if calendar.begins != 2 {
		t.Fatalf("begins = %d; the next scheduled fill must acquire the calendar again", calendar.begins)
	}
	laterReceipt := later.UTC().Format(time.RFC3339Nano)
	fresh := make([]byte, 4)
	binary.LittleEndian.PutUint32(fresh, 56)
	view = provider.EnergyTotals()
	monthly = energyDirectionNamed(view, energyRegisterName(energyFuelSumHcThisMonth))
	if monthly.ValueKWh != nil || monthly.ObservedAt != laterReceipt || !bytes.Equal(monthly.Raw, fresh) {
		t.Fatalf("second monthly = %+v; a moved civil day must reject the new reply %x at %s", monthly, fresh, laterReceipt)
	}
	lifetime := energyDirectionNamed(view, energyRegisterName(energyFuelSumHc))
	if lifetime.ValueKWh == nil || lifetime.ObservedAt != laterReceipt {
		t.Fatalf("second lifetime = %+v; want a new receipt %s", lifetime, laterReceipt)
	}
	if countEnergyReads(bus.reads, energyFuelSumHcThisMonth) != 2 {
		t.Fatalf("this-month reads = %d; the scheduled fill kept the first reply", countEnergyReads(bus.reads, energyFuelSumHcThisMonth))
	}

	reads, begins := len(bus.reads), calendar.begins
	if _, err := fillEnergyCatalogues(poller, ctx, session); err != nil {
		t.Fatalf("view fill: %v", err)
	}
	if len(bus.reads) != reads || calendar.begins != begins {
		t.Fatalf("view reads=%d begins=%d; a view read must not poll or begin again", len(bus.reads), calendar.begins)
	}

	otherBus := &energyCatalogueBus{}
	other, otherProvider := prepareEnergyPoller(otherBus)
	other.refreshPublishedCatalogues(context.Background())
	if calendar.begins != begins {
		t.Fatalf("begins = %d; another owner used this installation calendar", calendar.begins)
	}
	otherMonthly := energyDirectionNamed(otherProvider.EnergyTotals(), energyRegisterName(energyFuelSumHcThisMonth))
	if otherMonthly.ValueKWh != nil {
		t.Fatalf("other monthly = %+v; a missing session must not accept a dated reply", otherMonthly)
	}
}

func startCatalogueSchedule(t *testing.T, poller *vaillantSemanticPoller, ctx context.Context) (fire func(), stop func()) {
	t.Helper()
	ticks := make(chan time.Time)
	ready := make(chan struct{})
	stopped := make(chan struct{})
	done := make(chan struct{}, 1)
	schedule := &semanticPollSchedule{
		tasks:  newSemanticTaskScheduler(),
		pollMu: &poller.pollMu,
		newTicker: func(time.Duration) *semanticPollTicker {
			close(ready)
			return &semanticPollTicker{C: ticks, stop: func() { close(stopped) }}
		},
		recurring: []semanticPollTask{{
			key:      semanticTaskRefreshCatalogues,
			interval: time.Second,
			refresh: func(ctx context.Context) {
				poller.refreshPublishedCatalogues(ctx)
				done <- struct{}{}
			},
		}},
	}
	runCtx, cancel := context.WithCancel(ctx)
	schedule.Start(runCtx)
	scheduleWait(t, ready)
	return func() {
			t.Helper()
			ticks <- time.Time{}
			scheduleWait(t, done)
		}, func() {
			cancel()
			scheduleWait(t, stopped)
		}
}

func TestVaillantEnergyCataloguesRejectsMutatedIdentityAxis(t *testing.T) {
	statistics, err := catalogvaillant.B516StatisticsDirectionEntries()
	if err != nil || len(statistics) == 0 {
		t.Fatalf("statistics directions: %v", err)
	}
	statisticsKey := statistics[0].B516StatisticsDirectionKey
	for _, entry := range statistics {
		if entry.B516StatisticsDirectionKey.B516StatisticsDirectionSelector != nil {
			statisticsKey = entry.B516StatisticsDirectionKey
			break
		}
	}
	if _, found, err := catalogvaillant.B516StatisticsDirectionLookup(statisticsKey); err != nil || !found {
		t.Fatalf("unmutated statistics lookup found=%v err=%v", found, err)
	}
	statisticsMutations := []struct {
		field  string
		mutate func(*catalogvaillant.B516StatisticsDirectionKey)
	}{
		{"family", func(key *catalogvaillant.B516StatisticsDirectionKey) {
			key.B516StatisticsDirectionFamily = "absent-family"
		}},
		{"variant", func(key *catalogvaillant.B516StatisticsDirectionKey) {
			key.B516StatisticsDirectionVariant = "absent-variant"
		}},
		{"operation", func(key *catalogvaillant.B516StatisticsDirectionKey) {
			value := uint8(255)
			key.B516StatisticsDirectionOperation = &value
		}},
		{"direction", func(key *catalogvaillant.B516StatisticsDirectionKey) {
			key.B516StatisticsDirection = "sideways"
		}},
		{"selector", func(key *catalogvaillant.B516StatisticsDirectionKey) {
			value := uint8(1)
			key.B516StatisticsDirectionSelector = &value
		}},
		{"field", func(key *catalogvaillant.B516StatisticsDirectionKey) {
			key.B516StatisticsDirectionField = "absent-field"
		}},
	}
	for _, mutation := range statisticsMutations {
		key := cloneStatisticsKey(statisticsKey)
		mutation.mutate(&key)
		if _, found, err := catalogvaillant.B516StatisticsDirectionLookup(key); err != nil || found {
			t.Fatalf("%s mutation still matched; dropping %s from equality would keep the original row", mutation.field, mutation.field)
		}
	}

	candidates, err := catalogvaillant.B507EnergyCandidateEntries()
	if err != nil || len(candidates) == 0 {
		t.Fatalf("energy candidates: %v", err)
	}
	candidateKey := candidates[0].B507EnergyCandidateKey
	if _, found, err := catalogvaillant.B507EnergyCandidateLookup(candidateKey); err != nil || !found {
		t.Fatalf("unmutated candidate lookup found=%v err=%v", found, err)
	}
	sideways := "sideways"
	role := "absent-role"
	candidateMutations := []struct {
		field  string
		mutate func(*catalogvaillant.B507EnergyCandidateKey)
	}{
		{"family", func(key *catalogvaillant.B507EnergyCandidateKey) {
			key.B507EnergyCandidateFamily = "absent-family"
		}},
		{"variant", func(key *catalogvaillant.B507EnergyCandidateKey) {
			key.B507EnergyCandidateVariant = "absent-variant"
		}},
		{"primary", func(key *catalogvaillant.B507EnergyCandidateKey) { key.B507EnergyCandidatePrimary = 1 }},
		{"secondary", func(key *catalogvaillant.B507EnergyCandidateKey) {
			value := uint8(255)
			key.B507EnergyCandidateSecondary = &value
		}},
		{"operation", func(key *catalogvaillant.B507EnergyCandidateKey) {
			value := uint8(255)
			key.B507EnergyCandidateOperation = &value
		}},
		{"direction", func(key *catalogvaillant.B507EnergyCandidateKey) {
			key.B507EnergyCandidateDirection = &sideways
		}},
		{"field", func(key *catalogvaillant.B507EnergyCandidateKey) {
			key.B507EnergyCandidateField = "absent-field"
		}},
		{"group", func(key *catalogvaillant.B507EnergyCandidateKey) {
			value := uint8(255)
			key.B507EnergyCandidateGroup = &value
		}},
		{"requested_instance", func(key *catalogvaillant.B507EnergyCandidateKey) {
			value := uint8(255)
			key.B507EnergyCandidateRequestedInstance = &value
		}},
		{"register", func(key *catalogvaillant.B507EnergyCandidateKey) {
			value := uint16(65535)
			key.B507EnergyCandidateRegister = &value
		}},
		{"command_role", func(key *catalogvaillant.B507EnergyCandidateKey) {
			key.B507EnergyCandidateCommandRole = &role
		}},
	}
	for _, mutation := range candidateMutations {
		key := cloneCandidateKey(candidateKey)
		mutation.mutate(&key)
		if _, found, err := catalogvaillant.B507EnergyCandidateLookup(key); err != nil || found {
			t.Fatalf("%s mutation still matched; dropping %s from equality would keep the original row", mutation.field, mutation.field)
		}
	}
}

func prepareEnergyPoller(bus *energyCatalogueBus) (*vaillantSemanticPoller, *retainedSemanticState) {
	poller, provider := newCataloguePoller(bus.send)
	instant := energyTestReference{}.Instant()
	provider.SetEnergyClock(func() time.Time { return instant })
	return poller, provider
}

func assertEnergyReads(t *testing.T, reads []uint16) {
	t.Helper()
	if len(reads) != len(b524EnergyQueries) {
		t.Fatalf("reads = %d; want one read for each catalogued register and no candidate read", len(reads))
	}
	seen := map[uint16]int{}
	for _, addr := range reads {
		if _, ok := energyCatalogueValue(addr); !ok {
			t.Fatalf("read 0x%04x; a candidate or residual must not invent an address", addr)
		}
		seen[addr]++
	}
	for _, query := range b524EnergyQueries {
		if seen[query.addr] != 1 {
			t.Fatalf("register 0x%04x reads = %d; want one", query.addr, seen[query.addr])
		}
	}
}

func assertAcceptedEnergy(t *testing.T, view output.EnergyTotals, ref energyTestReference, receipt, noon string) {
	t.Helper()
	if len(view.RegisterDirections) != len(b524EnergyQueries) {
		t.Fatalf("register directions = %d; want %d", len(view.RegisterDirections), len(b524EnergyQueries))
	}
	for i, query := range b524EnergyQueries {
		item := view.RegisterDirections[i]
		value, _ := energyCatalogueValue(query.addr)
		raw := make([]byte, 4)
		binary.LittleEndian.PutUint32(raw, value)
		wantYear, wantMonth := ref.Year(), 0
		switch query.yearKind {
		case "previous":
			month := ref.PreviousMonth()
			wantYear, wantMonth = month.Year, int(month.Month)
		case "current":
			month := ref.Month()
			wantYear, wantMonth = month.Year, int(month.Month)
		}
		if item.Name != energyRegisterName(query.addr) || item.Period != query.period || item.Source != query.channel || item.Usage != query.usage || item.Level != "Unknown" || item.Freshness != "Unknown" || item.ScopeYear != wantYear || item.ScopeMonth != wantMonth || item.ValueKWh == nil || *item.ValueKWh != float64(value) || !bytes.Equal(item.Raw, raw) || item.ObservedAt != receipt || item.ObservedAt == noon {
			t.Fatalf("direction %d = %+v; want %s %s %s %g scoped %d-%d at the wire receipt %s", i, item, query.period, query.channel, query.usage, float64(value), wantYear, wantMonth, receipt)
		}
	}
}

func assertEnergyProjection(t *testing.T, totals *output.EnergyTotals) {
	t.Helper()
	if totals == nil {
		t.Fatalf("energy totals are nil")
	}
	checks := []struct {
		name    string
		yearly  []float64
		monthly []float64
		year    float64
		prev    float64
		cur     float64
	}{
		{"gas climate", totals.Gas.Climate.Yearly, totals.Gas.Climate.Monthly, 11, 99, 55},
		{"gas hot water", totals.Gas.DHW.Yearly, totals.Gas.DHW.Monthly, 22, 111, 66},
		{"electricity climate", totals.Electric.Climate.Yearly, totals.Electric.Climate.Monthly, 33, 122, 77},
		{"electricity hot water", totals.Electric.DHW.Yearly, totals.Electric.DHW.Monthly, 44, 133, 88},
	}
	for _, check := range checks {
		if len(check.yearly) != 2 || check.yearly[0] != 0 || check.yearly[1] != check.year || len(check.monthly) != 2 || check.monthly[0] != check.prev || check.monthly[1] != check.cur {
			t.Fatalf("%s yearly %v monthly %v; want lifetime %g in the current year and previous/current %g/%g", check.name, check.yearly, check.monthly, check.year, check.prev, check.cur)
		}
	}
}

func assertEnergyCatalogue(t *testing.T, view output.EnergyTotals, supplied []byte) {
	t.Helper()
	statistics, err := catalogvaillant.B516StatisticsDirectionEntries()
	if err != nil {
		t.Fatalf("statistics directions: %v", err)
	}
	rawTypes, err := catalogvaillant.B516DirectionalTypeEntries()
	if err != nil {
		t.Fatalf("directional types: %v", err)
	}
	var types []json.RawMessage
	if err := json.Unmarshal(rawTypes, &types); err != nil {
		t.Fatalf("directional types: %v", err)
	}
	if len(view.CatalogueDirections) != len(statistics)+len(types) {
		t.Fatalf("catalogue directions = %d; want %d statistics and %d directional types", len(view.CatalogueDirections), len(statistics), len(types))
	}
	unresolved := 0
	for _, item := range view.CatalogueDirections {
		if item.ValueKWh != nil || len(item.Raw) != 0 {
			t.Fatalf("catalogue direction %s copied a register value %v raw %x", item.Name, item.ValueKWh, item.Raw)
		}
		if item.Label == catalogueUnresolved {
			unresolved++
		} else if item.Label != "" {
			t.Fatalf("catalogue label %q; an unknown enum stays unresolved", item.Label)
		}
	}
	if unresolved != len(types) {
		t.Fatalf("unresolved catalogue labels = %d; want one per directional type", unresolved)
	}
	setpoints := 0
	for _, item := range view.ResidualFields {
		if item.Field != "setpoint" {
			if len(item.Raw) != 0 {
				t.Fatalf("residual %s raw %x; only the supplied field keeps bytes", item.Field, item.Raw)
			}
			continue
		}
		setpoints++
		if !bytes.Equal(item.Raw, supplied) {
			t.Fatalf("setpoint residual raw %x; divide_by would replace the supplied bytes %x", item.Raw, supplied)
		}
	}
	if setpoints == 0 {
		t.Fatalf("no setpoint residual kept the supplied bytes")
	}
	candidates, err := catalogvaillant.B507EnergyCandidateEntries()
	if err != nil {
		t.Fatalf("energy candidates: %v", err)
	}
	if len(view.CandidateFields) != len(candidates) {
		t.Fatalf("candidates = %d; want %d", len(view.CandidateFields), len(candidates))
	}
	for _, item := range view.CandidateFields {
		if item.Register != nil || item.Note != "register absent" {
			t.Fatalf("candidate %s register %v note %q; dropping the absent register would still publish an address", item.Field, item.Register, item.Note)
		}
	}
}

func assertDatedEnergyRejected(t *testing.T, view output.EnergyTotals, provider *retainedSemanticState, receipt string) {
	t.Helper()
	if len(view.RegisterDirections) != len(b524EnergyQueries) {
		t.Fatalf("directions = %d; want %d", len(view.RegisterDirections), len(b524EnergyQueries))
	}
	for i, query := range b524EnergyQueries {
		item := view.RegisterDirections[i]
		if query.period == "month" {
			if item.ValueKWh != nil || item.ObservedAt != receipt || len(item.Raw) != 4 {
				t.Fatalf("%s = %+v; a rejected dated reply keeps its wire receipt and raw bytes", item.Name, item)
			}
			continue
		}
		value, _ := energyCatalogueValue(query.addr)
		if item.ValueKWh == nil || *item.ValueKWh != float64(value) || item.ObservedAt != receipt {
			t.Fatalf("lifetime %s = %+v; want %d at %s", item.Name, item, value, receipt)
		}
	}
	published := provider.EnergyTotals()
	if published == nil || len(published.Gas.Climate.Yearly) < 2 || published.Gas.Climate.Yearly[1] != 11 {
		t.Fatalf("lifetime projection = %+v; want 11", published)
	}
	if len(published.Gas.Climate.Monthly) >= 2 && (published.Gas.Climate.Monthly[0] != 0 || published.Gas.Climate.Monthly[1] != 0) {
		t.Fatalf("monthly projection = %v; the guard must leave it empty", published.Gas.Climate.Monthly)
	}
}

func energyDirectionNamed(view *output.EnergyTotals, name string) output.EnergyDirection {
	if view == nil {
		return output.EnergyDirection{}
	}
	for _, item := range view.RegisterDirections {
		if item.Name == name {
			return item
		}
	}
	return output.EnergyDirection{}
}

func countEnergyReads(reads []uint16, addr uint16) int {
	count := 0
	for _, item := range reads {
		if item == addr {
			count++
		}
	}
	return count
}

func energyCatalogueValue(addr uint16) (uint32, bool) {
	switch addr {
	case energyFuelSumHc:
		return 11, true
	case energyFuelSumHwc:
		return 22, true
	case energyElectricitySumHc:
		return 33, true
	case energyElectricitySumHwc:
		return 44, true
	case energyFuelSumHcThisMonth:
		return 55, true
	case energyFuelSumHwcThisMonth:
		return 66, true
	case energyElectricitySumHcThisMonth:
		return 77, true
	case energyElectricitySumHwcThisMonth:
		return 88, true
	case energyFuelSumHcLastMonth:
		return 99, true
	case energyFuelSumHwcLastMonth:
		return 111, true
	case energyElectricitySumHcLastMonth:
		return 122, true
	case energyElectricitySumHwcLastMonth:
		return 133, true
	default:
		return 0, false
	}
}

type energyCatalogueBus struct {
	reads    []uint16
	advance  func()
	advanced bool
	replace  map[uint16]uint32
}

func (bus *energyCatalogueBus) send(_ context.Context, frame protocol.Frame) (*protocol.Frame, error) {
	if len(frame.Data) != 6 {
		return &protocol.Frame{Data: []byte{0x00}}, nil
	}
	addr := uint16(frame.Data[4]) | uint16(frame.Data[5])<<8
	bus.reads = append(bus.reads, addr)
	if addr == energyFuelSumHcThisMonth && bus.advance != nil && !bus.advanced {
		bus.advanced = true
		bus.advance()
	}
	value, ok := energyCatalogueValue(addr)
	if replaced, found := bus.replace[addr]; found {
		value, ok = replaced, true
	}
	if !ok {
		return &protocol.Frame{Data: []byte{0x00}}, nil
	}
	raw := []byte{byte(value), byte(value >> 8), byte(value >> 16), byte(value >> 24)}
	return testB524ResponseForSelectorPayload(frame.Data, raw...), nil
}

func cloneStatisticsKey(key catalogvaillant.B516StatisticsDirectionKey) catalogvaillant.B516StatisticsDirectionKey {
	if key.B516StatisticsDirectionOperation != nil {
		value := *key.B516StatisticsDirectionOperation
		key.B516StatisticsDirectionOperation = &value
	}
	if key.B516StatisticsDirectionSelector != nil {
		value := *key.B516StatisticsDirectionSelector
		key.B516StatisticsDirectionSelector = &value
	}
	return key
}

func cloneCandidateKey(key catalogvaillant.B507EnergyCandidateKey) catalogvaillant.B507EnergyCandidateKey {
	if key.B507EnergyCandidateSecondary != nil {
		value := *key.B507EnergyCandidateSecondary
		key.B507EnergyCandidateSecondary = &value
	}
	if key.B507EnergyCandidateOperation != nil {
		value := *key.B507EnergyCandidateOperation
		key.B507EnergyCandidateOperation = &value
	}
	if key.B507EnergyCandidateDirection != nil {
		value := *key.B507EnergyCandidateDirection
		key.B507EnergyCandidateDirection = &value
	}
	if key.B507EnergyCandidateGroup != nil {
		value := *key.B507EnergyCandidateGroup
		key.B507EnergyCandidateGroup = &value
	}
	if key.B507EnergyCandidateRequestedInstance != nil {
		value := *key.B507EnergyCandidateRequestedInstance
		key.B507EnergyCandidateRequestedInstance = &value
	}
	if key.B507EnergyCandidateRegister != nil {
		value := *key.B507EnergyCandidateRegister
		key.B507EnergyCandidateRegister = &value
	}
	if key.B507EnergyCandidateCommandRole != nil {
		value := *key.B507EnergyCandidateCommandRole
		key.B507EnergyCandidateCommandRole = &value
	}
	return key
}
