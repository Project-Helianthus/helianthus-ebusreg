package owner

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	catalogvaillant "github.com/Project-Helianthus/helianthus-ebusreg/catalog/vaillant"
	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/output"
	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
)

func TestVaillantCatalogueStateFollowsItsInstance(t *testing.T) {
	left, leftView := newPublishedCatalogueOwner(0x71, 0x15, 0x08, "controller-700")
	right, rightView := newPublishedCatalogueOwner(0x72, 0x16, 0x09, "controller-720")
	if left.source == right.source || left.controller == right.controller || left.boilerAddress == right.boilerAddress || left.codebookFamily() == right.codebookFamily() {
		t.Fatal("the two instances do not differ in source, controller, boiler address and family")
	}
	runCatalogueRefresh(t, left)
	runCatalogueRefresh(t, right)
	assertPublishedCatalogueViews(t, leftView)
	if rightView.DHW() == nil || len(rightView.DHW().State.Enumerations) == 0 || rightView.BoilerStatus() == nil || len(rightView.BoilerStatus().Selectors) == 0 {
		t.Fatal("second instance published no catalogue rows")
	}
	if totals := rightView.EnergyTotals(); totals == nil || len(totals.CatalogueDirections) == 0 {
		t.Fatal("second instance published no energy directions")
	}
	survivor := catalogueProjectionJSON(t, rightView)
	if bytes.Equal(survivor, catalogueProjectionJSON(t, leftView)) {
		t.Fatal("distinct families published the same catalogue view")
	}

	left.releaseCatalogueState()
	if !left.catalogueStateReleased() || right.catalogueStateReleased() {
		t.Fatal("release did not stay on the stopped instance")
	}
	if !bytes.Equal(catalogueProjectionJSON(t, rightView), survivor) {
		t.Fatal("stopping one instance changed the other's catalogue view")
	}
	if totals := leftView.EnergyTotals(); totals != nil && (len(totals.RegisterDirections) != 0 || len(totals.CatalogueDirections) != 0) {
		t.Fatal("released instance still overlays its energy catalogue")
	}
	republishCatalogueViews(left)
	republishCatalogueViews(right)
	if zones := leftView.Zones(); len(zones) != 1 || len(zones[0].State.Enumerations) != 0 {
		t.Fatal("released instance reapplied codebook labels")
	}
	if status := leftView.BoilerStatus(); status != nil && len(status.Selectors) != 0 {
		t.Fatal("released instance reapplied boiler dictionary entries")
	}
	if !bytes.Equal(catalogueProjectionJSON(t, rightView), survivor) {
		t.Fatal("republishing the survivor changed its catalogue view")
	}
	runCatalogueRefresh(t, left)
	if !left.catalogueStateReleased() || !errors.Is(left.catalogueFillError(), errCatalogueStateReleased) {
		t.Fatal("a refresh after release restored catalogue state")
	}
	if totals := leftView.EnergyTotals(); totals != nil && len(totals.CatalogueDirections) != 0 {
		t.Fatal("a refresh after release restored energy directions")
	}
	if !bytes.Equal(catalogueProjectionJSON(t, rightView), survivor) {
		t.Fatal("a refresh of the stopped instance changed the survivor")
	}
}

func TestVaillantCatalogueFillFailureIsReported(t *testing.T) {
	poller, provider := newPublishedCatalogueOwner(0x71, 0x15, 0x08, "controller-700")
	runCatalogueRefresh(t, poller)
	assertPublishedCatalogueViews(t, provider)
	before := catalogueProjectionJSON(t, provider)
	if poller.catalogueFillError() != nil {
		t.Fatal("successful catalogue fill recorded an error")
	}
	previousEnergy := loadEnergyDirectionCatalogue
	previousBoiler := loadBoilerDictionaryEntries
	t.Cleanup(func() {
		loadEnergyDirectionCatalogue = previousEnergy
		loadBoilerDictionaryEntries = previousBoiler
	})
	loadEnergyDirectionCatalogue = func() ([]output.EnergyDirection, error) {
		return nil, errors.New("energy catalogue unavailable")
	}
	loadBoilerDictionaryEntries = func() ([]catalogvaillant.B509BoilerCodecDictionaryEntry, error) {
		return nil, errors.New("boiler dictionary unavailable")
	}
	runCatalogueRefresh(t, poller)
	err := poller.catalogueFillError()
	if err == nil || !bytes.Contains([]byte(err.Error()), []byte("energy catalogue unavailable")) || !bytes.Contains([]byte(err.Error()), []byte("boiler dictionary unavailable")) {
		t.Fatalf("scheduled catalogue pass recorded %v", err)
	}
	if !bytes.Equal(catalogueProjectionJSON(t, provider), before) {
		t.Fatal("failed catalogue fill changed the published view")
	}
	assertPublishedCatalogueViews(t, provider)
}

func TestVaillantCatalogueLabelsReachThePublishedView(t *testing.T) {
	poller, provider := newPublishedCatalogueOwner(0x71, 0x15, 0x08, "controller-700")
	runCatalogueRefresh(t, poller)
	raw := catalogueConsumerJSON(t, provider)
	for _, label := range []string{"ventilation-boost", "flame_active", "fuel_sum_heating", "layout"} {
		if !bytes.Contains(raw, []byte(label)) {
			t.Fatalf("published view %s lacks %s", raw, label)
		}
	}
	zones := output.MCPZonesFromGraphQL(provider.Zones())
	if len(zones) != 1 || zones[0].State.Enumerations[0].Label == "" {
		t.Fatal("MCP zone view has no enumeration label")
	}
	boiler := output.MCPBoilerStatusFromGraphQL(provider.BoilerStatus())
	if boiler == nil || len(boiler.Selectors) == 0 || boiler.Selectors[0].Label == "" {
		t.Fatal("MCP boiler view has no dictionary entry")
	}
	energy := output.MCPEnergyTotalsFromGraphQL(provider.EnergyTotals())
	if energy == nil || len(energy.RegisterDirections) == 0 || energy.RegisterDirections[0].Name == "" || len(energy.CatalogueDirections) == 0 || energy.CatalogueDirections[0].Field == "" {
		t.Fatal("MCP energy view has no catalogue direction")
	}
}

func newPublishedCatalogueOwner(source, controller, boiler byte, family string) (*vaillantSemanticPoller, *retainedSemanticState) {
	bus := &publishedCatalogueBus{
		b524: &catalogueReplyBus{payloads: publishedCataloguePayloads()},
		boiler: &boilerReplyBus{values: map[uint16][]byte{
			boilerB509FlameActive:    {0x00},
			boilerB509GasValveActive: {0x00},
			boilerB509StateNumber:    {0x07},
			boilerB509PumpHours:      {0x05, 0x00},
		}},
	}
	poller, provider := newCataloguePoller(bus.send)
	poller.source = source
	poller.controller = controller
	poller.boilerAddress = boiler
	poller.configInterval = time.Minute
	reg := registry.NewDeviceRegistry(nil)
	reg.Register(registry.DeviceInfo{Address: controller, Manufacturer: "Vaillant", DeviceID: family})
	poller.reg = reg
	poller.zones = map[byte]*vaillantZoneSnapshot{0: {}}
	poller.dhw = &vaillantDhwSnapshot{}
	poller.circuits = map[byte]*vaillantCircuitSnapshot{0: {Instance: 1, Active: true}}
	poller.system = &vaillantSystemSnapshot{Controller: controller}
	poller.boiler = &vaillantBoilerSnapshot{}
	republishCatalogueViews(poller)
	return poller, provider
}

func republishCatalogueViews(poller *vaillantSemanticPoller) {
	poller.publishZones(semanticSnapshotSourceLive)
	poller.publishDHW(semanticSnapshotSourceLive)
	poller.publishCircuits(semanticSnapshotSourceLive)
	poller.publishSystem(semanticSnapshotSourceLive)
	poller.publishBoilerStatus(semanticSnapshotSourceLive)
}

func runCatalogueRefresh(t *testing.T, poller *vaillantSemanticPoller) {
	t.Helper()
	for _, item := range poller.pollSchedule().recurring {
		if item.key == semanticTaskRefreshCatalogues {
			item.refresh(context.Background())
			return
		}
	}
	t.Fatal("catalogue refresh task missing")
}

func catalogueProjectionJSON(t *testing.T, provider *retainedSemanticState) []byte {
	t.Helper()
	var zones, dhw, circuits any
	if published := provider.Zones(); len(published) == 1 {
		zones = published[0].State.Enumerations
	}
	if published := provider.DHW(); published != nil {
		dhw = published.State.Enumerations
	}
	if published := provider.Circuits(); len(published) == 1 {
		circuits = published[0].Config.Enumerations
	}
	var ventilation any
	var catalogue any
	if published := provider.System(); published != nil {
		ventilation = published.Ventilation
		catalogue = published.Catalogue
	}
	var selectors, operations any
	if published := provider.BoilerStatus(); published != nil {
		selectors = published.Selectors
		operations = published.Operations
	}
	var registers, directions, residuals, candidates any
	if published := provider.EnergyTotals(); published != nil {
		registers = published.RegisterDirections
		directions = published.CatalogueDirections
		residuals = published.ResidualFields
		candidates = published.CandidateFields
	}
	raw, err := json.Marshal(struct {
		Zones       any `json:"zones"`
		DHW         any `json:"dhw"`
		Circuits    any `json:"circuits"`
		Ventilation any `json:"ventilation"`
		Catalogue   any `json:"catalogue"`
		Selectors   any `json:"selectors"`
		Operations  any `json:"operations"`
		Registers   any `json:"registers"`
		Directions  any `json:"directions"`
		Residuals   any `json:"residuals"`
		Candidates  any `json:"candidates"`
	}{zones, dhw, circuits, ventilation, catalogue, selectors, operations, registers, directions, residuals, candidates})
	if err != nil {
		t.Fatal(err)
	}
	return raw
}

func catalogueConsumerJSON(t *testing.T, provider *retainedSemanticState) []byte {
	t.Helper()
	raw, err := json.Marshal(struct {
		Zones    []output.MCPZone          `json:"zones"`
		DHW      *output.MCPDhwStatus      `json:"dhw"`
		Circuits []output.MCPCircuitStatus `json:"circuits"`
		System   *output.MCPSystemStatus   `json:"system"`
		Boiler   *output.MCPBoilerStatus   `json:"boiler"`
		Energy   *output.MCPEnergyTotals   `json:"energy"`
	}{
		Zones:    output.MCPZonesFromGraphQL(provider.Zones()),
		DHW:      output.MCPDHWFromGraphQL(provider.DHW()),
		Circuits: output.MCPCircuitsFromGraphQL(provider.Circuits()),
		System:   output.MCPSystemFromGraphQL(provider.System()),
		Boiler:   output.MCPBoilerStatusFromGraphQL(provider.BoilerStatus()),
		Energy:   output.MCPEnergyTotalsFromGraphQL(provider.EnergyTotals()),
	})
	if err != nil {
		t.Fatal(err)
	}
	return raw
}
