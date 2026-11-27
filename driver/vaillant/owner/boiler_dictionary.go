package owner

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"

	catalogvaillant "github.com/Project-Helianthus/helianthus-ebusreg/catalog/vaillant"
	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/output"
)

const (
	rawZoneSpecialFunctionName = "zone_special_function_raw"
	catalogueWritePermission   = byte(0x02)
	boilerDictionaryCodecOnOff = "UCH (0x0F=on, 0xF0=off)"
)

type boilerCatalogueHold struct {
	selectors  []output.BoilerSelectorEvidence
	operations []output.RawOperation
}

var (
	boilerCatalogueMu sync.Mutex
	boilerCatalogues  = map[*vaillantSemanticPoller]*boilerCatalogueHold{}
)

func polledBoilerSelectors() map[uint16]string {
	return map[uint16]string{
		boilerB509WaterPressure:         "fast",
		boilerB509FlameActive:           "fast",
		boilerB509PartloadHcKw:          "medium",
		boilerB509PartloadHwcKw:         "medium",
		boilerB509FlowsetHcMaxC:         "medium",
		boilerB509FlowsetHcMaxCFallback: "medium",
		boilerB509FlowsetHwcMaxC:        "medium",
		boilerB509PumpHours:             "slow",
		boilerB509FlowTemperature:       "fast",
		boilerB509FanHours:              "slow",
		boilerB509DeactivationsIfc:      "slow",
		boilerB509HoursTillService:      "slow",
		boilerB509DeactivationsLimit:    "slow",
		boilerB509DhwHours:              "slow",
		boilerB509DhwStarts:             "slow",
		boilerB509TargetFanSpeedRpm:     "medium",
		boilerB509CentralHeatingHours:   "slow",
		boilerB509CentralHeatingStarts:  "slow",
		boilerB509ModulationPct:         "fast",
		boilerB509FlowTempDesiredC:      "medium",
		boilerB509ExternalPumpActive:    "medium",
		boilerB509InstallerMenuCode:     "slow",
		boilerB509CentralHeatingPump:    "fast",
		boilerB509DiverterValvePosition: "fast",
		boilerB509DhwWaterFlowLpm:       "fast",
		boilerB509DhwDemandActive:       "fast",
		boilerB509CirculationPumpActive: "medium",
		boilerB509PhoneNumber:           "slow",
		boilerB509FanSpeedRpm:           "fast",
		boilerB509StorageLoadPumpPct:    "medium",
		boilerB509IonisationVoltageUa:   "medium",
		boilerB509StateNumber:           "fast",
		boilerB509GasValveActive:        "fast",
		boilerB509HeatingSwitchActive:   "medium",
		boilerB509DhwTempDesiredC:       "medium",
		boilerB509PrimaryCircuitFlowLpm: "medium",
	}
}

func fillBoilerDictionary(p *vaillantSemanticPoller, ctx context.Context, boilerAddress byte) ([]output.BoilerSelectorEvidence, error) {
	if p == nil || boilerAddress == 0 {
		return nil, fmt.Errorf("boiler dictionary target is zero")
	}
	entries, err := catalogvaillant.B509BoilerCodecDictionaryEntries()
	if err != nil {
		return nil, err
	}
	polled := polledBoilerSelectors()
	seen := map[uint16]bool{}
	var selectors []output.BoilerSelectorEvidence
	for _, addr := range []uint16{boilerB509FlameActive, boilerB509GasValveActive, boilerB509StateNumber, boilerB509PumpHours} {
		item, itemErr := boilerSelectorEvidence(p, ctx, boilerAddress, addr, polled[addr])
		if itemErr != nil {
			return nil, itemErr
		}
		selectors = append(selectors, item)
		seen[addr] = true
	}
	for _, entry := range entries {
		addr, ok := dictionarySelectorAddress(entry.B509BoilerCodecDictionaryKey.B509BoilerCodecDictionarySelector)
		if !ok || seen[addr] {
			continue
		}
		seen[addr] = true
		if _, already := polled[addr]; already {
			continue
		}
		raw, readOK := p.readB509Value(ctx, boilerAddress, addr)
		item := output.BoilerSelectorEvidence{
			Selector: addr,
			Raw:      append([]byte(nil), raw...),
			Label:    catalogueUnresolved,
			Level:    entry.B509BoilerCodecDictionaryGrade,
			Tier:     "slow",
			Note:     "dictionary codec withheld",
		}
		if !readOK {
			item.Raw = nil
		}
		if dictionaryEstablishesExecutable(addr, entry.B509BoilerCodecDeclaration) {
			item.CodecUsed = entry.B509BoilerCodecDeclaration
			item.Note = ""
		}
		selectors = append(selectors, item)
	}
	storeBoilerSelectors(p, selectors)
	return append([]output.BoilerSelectorEvidence(nil), selectors...), nil
}

func boilerSelectorEvidence(p *vaillantSemanticPoller, ctx context.Context, boilerAddress byte, addr uint16, tier string) (output.BoilerSelectorEvidence, error) {
	raw, _ := p.readB509Value(ctx, boilerAddress, addr)
	item := output.BoilerSelectorEvidence{
		Selector:  addr,
		Raw:       append([]byte(nil), raw...),
		Tier:      tier,
		CodecUsed: executableBoilerCodec(addr),
		Label:     catalogueUnresolved,
		Level:     "Unknown",
	}
	declaration, declared := dictionaryDeclaration(addr)
	if declared && dictionaryEstablishesExecutable(addr, declaration) {
		item.CodecUsed = declaration
	}
	name, unresolved, grade, note, annotated := boilerSymbolicAnnotation(addr)
	if !annotated {
		item.Note = "annotation absent"
		return item, nil
	}
	item.Name = name
	item.Level = grade
	item.Note = note
	if unresolved {
		item.Label = catalogueUnresolved
		return item, nil
	}
	if name != "" {
		item.Label = name
	}
	return item, nil
}

func dictionarySelectorAddress(selector []uint16) (uint16, bool) {
	if len(selector) != 2 || selector[0] > 255 || selector[1] > 255 {
		return 0, false
	}
	return uint16(selector[0])<<8 | uint16(selector[1]), true
}

func dictionaryDeclaration(addr uint16) (string, bool) {
	entry, found, err := catalogvaillant.B509BoilerCodecDictionaryLookup(boilerDictionaryKey(addr))
	if err != nil || !found {
		return "", false
	}
	return entry.B509BoilerCodecDeclaration, true
}

func boilerDictionaryKey(addr uint16) catalogvaillant.B509BoilerCodecDictionaryKey {
	primary := byte(181)
	secondary := byte(9)
	operation := byte(13)
	return catalogvaillant.B509BoilerCodecDictionaryKey{
		B509BoilerCodecDictionaryFamily:    "boiler",
		B509BoilerCodecDictionaryVariant:   "direct-symbolic-candidate",
		B509BoilerCodecDictionaryPrimary:   &primary,
		B509BoilerCodecDictionarySecondary: &secondary,
		B509BoilerCodecDictionaryOperation: &operation,
		B509BoilerCodecDictionaryDirection: "read",
		B509BoilerCodecDictionarySelector:  []uint16{uint16(addr >> 8), uint16(addr & 0xFF)},
	}
}

func dictionaryEstablishesExecutable(addr uint16, declaration string) bool {
	switch addr {
	case boilerB509FlameActive, boilerB509GasValveActive:
		on, ok := decodeOnOff([]byte{0x00})
		off, offOK := decodeOnOff([]byte{0xF0})
		if !ok || !offOK || !on || off {
			return false
		}
		return declaration != "" && declaration != boilerDictionaryCodecOnOff && declaration == boilerDictionaryCodecOnOff
	default:
		return false
	}
}

func executableBoilerCodec(addr uint16) string {
	switch addr {
	case boilerB509FlameActive, boilerB509GasValveActive, boilerB509CentralHeatingPump, boilerB509DhwDemandActive, boilerB509CirculationPumpActive, boilerB509ExternalPumpActive, boilerB509HeatingSwitchActive:
		return "onoff"
	case boilerB509PumpHours, boilerB509FanHours, boilerB509CentralHeatingHours, boilerB509DhwHours, boilerB509HoursTillService:
		return "hoursum2"
	case boilerB509StateNumber, boilerB509DeactivationsIfc, boilerB509DeactivationsLimit, boilerB509InstallerMenuCode:
		return "uch"
	default:
		return ""
	}
}

func boilerSymbolicAnnotation(addr uint16) (name string, unresolved bool, grade, note string, ok bool) {
	selector := fmt.Sprintf("0x%04X", addr)
	secondary := uint8(9)
	operation := uint8(13)
	key := catalogvaillant.B509SymbolicAnnotationKey{
		B509SymbolicAnnotationFamily:        "Vaillant",
		B509SymbolicAnnotationVariant:       "none",
		B509SymbolicAnnotationPrimary:       181,
		B509SymbolicAnnotationSecondary:     &secondary,
		B509SymbolicAnnotationOperation:     &operation,
		B509SymbolicAnnotationDirection:     "none",
		B509SymbolicAnnotationField:         "symbolic_annotation",
		B509SymbolicAnnotationAlternative:   "none",
		B509SymbolicAnnotationInstanceScope: "none",
		B509SymbolicAnnotationClass:         "boiler",
		B509SymbolicAnnotationSelector:      &selector,
	}
	raw, found, err := catalogvaillant.B509SymbolicAnnotationLookup(key)
	if err != nil || !found {
		return "", true, "Unknown", "", false
	}
	var decoded struct {
		Name        string `json:"name"`
		Unresolved  bool   `json:"unresolved"`
		Grade       string `json:"grade"`
		RawFallback string `json:"raw_fallback"`
	}
	if json.Unmarshal(raw, &decoded) != nil {
		return "", true, "Unknown", "", false
	}
	note = ""
	if decoded.Unresolved {
		note = decoded.RawFallback
	}
	return decoded.Name, decoded.Unresolved, decoded.Grade, note, true
}

func storeBoilerSelectors(p *vaillantSemanticPoller, selectors []output.BoilerSelectorEvidence) {
	boilerCatalogueMu.Lock()
	hold := boilerCatalogueHoldFor(p)
	hold.selectors = append([]output.BoilerSelectorEvidence(nil), selectors...)
	operations := append([]output.RawOperation(nil), hold.operations...)
	boilerCatalogueMu.Unlock()
	overlayBoilerCatalogue(p, selectors, operations)
}

func reapplyStoredBoilerCatalogue(p *vaillantSemanticPoller) {
	if p == nil {
		return
	}
	boilerCatalogueMu.Lock()
	hold := boilerCatalogues[p]
	if hold == nil || (len(hold.selectors) == 0 && len(hold.operations) == 0) {
		boilerCatalogueMu.Unlock()
		return
	}
	selectors := append([]output.BoilerSelectorEvidence(nil), hold.selectors...)
	operations := append([]output.RawOperation(nil), hold.operations...)
	boilerCatalogueMu.Unlock()
	overlayBoilerCatalogue(p, selectors, operations)
}

func overlayBoilerCatalogue(p *vaillantSemanticPoller, selectors []output.BoilerSelectorEvidence, operations []output.RawOperation) {
	owner := publicationOwner(p)
	if owner == nil || owner.state == nil {
		return
	}
	status := owner.state.BoilerStatus()
	if status == nil {
		return
	}
	status.SetSelectors(selectors)
	status.SetOperations(operations)
	owner.state.SetBoilerStatus(status)
}

func boilerCatalogueHoldFor(p *vaillantSemanticPoller) *boilerCatalogueHold {
	hold := boilerCatalogues[p]
	if hold == nil {
		hold = &boilerCatalogueHold{}
		boilerCatalogues[p] = hold
	}
	return hold
}

func admitOpenSelector(p *vaillantSemanticPoller, ctx context.Context) (output.RawOperation, bool) {
	operation := output.RawOperation{
		Namespace: "ebusreg",
		Name:      rawZoneSpecialFunctionName,
		Group:     localZones.group,
		Register:  zoneSpecialFunction,
	}
	reply, ok := readOpenSelectorReply(p, ctx, localZones.group, 0, zoneSpecialFunction)
	admitted := ok && len(reply) > 0 && reply[0]&catalogueWritePermission != 0
	boilerCatalogueMu.Lock()
	hold := boilerCatalogueHoldFor(p)
	if admitted {
		hold.operations = []output.RawOperation{operation}
	} else {
		hold.operations = nil
		operation = output.RawOperation{}
	}
	selectors := append([]output.BoilerSelectorEvidence(nil), hold.selectors...)
	operations := append([]output.RawOperation(nil), hold.operations...)
	boilerCatalogueMu.Unlock()
	overlayBoilerCatalogue(p, selectors, operations)
	return operation, admitted
}

func readOpenSelectorReply(p *vaillantSemanticPoller, ctx context.Context, group, instance byte, addr uint16) ([]byte, bool) {
	if p == nil || (p.bus == nil && p.sendFrameFn == nil) || p.controller == 0 {
		return nil, false
	}
	if ctx == nil {
		ctx = context.Background()
	}
	p.mu.Lock()
	source := p.source
	target := p.controller
	p.mu.Unlock()
	response, err := p.sendSemanticFrame(ctx, protocol.Frame{
		FrameType: protocol.FrameTypeInitiatorTarget,
		Source:    source,
		Target:    target,
		Primary:   vaillantExtRegisterPrimary,
		Secondary: vaillantExtRegisterSecondary,
		Data:      buildB524ReadSelector(vaillantB524OpcodeLocal, group, instance, addr),
	})
	if err != nil || response == nil || len(response.Data) == 0 {
		return nil, false
	}
	if _, parsed := parseB524ReadPayload(response.Data, vaillantB524OpcodeLocal, group, instance, addr); !parsed {
		return append([]byte(nil), response.Data...), false
	}
	return append([]byte(nil), response.Data...), true
}

func callRawSelectorOperation(p *vaillantSemanticPoller, ctx context.Context, name string, payload []byte) error {
	boilerCatalogueMu.Lock()
	hold := boilerCatalogues[p]
	var admitted *output.RawOperation
	if hold != nil {
		for i := range hold.operations {
			if hold.operations[i].Name == name && hold.operations[i].Namespace == "ebusreg" {
				copied := hold.operations[i]
				admitted = &copied
				break
			}
		}
	}
	boilerCatalogueMu.Unlock()
	if admitted == nil {
		return fmt.Errorf("raw operation is not admitted")
	}
	return p.writeB524Value(ctx, vaillantB524OpcodeLocal, admitted.Group, 0, admitted.Register, append([]byte(nil), payload...))
}
