package pv

import (
	"fmt"
	"net"
	"regexp"
	"sort"
	"strings"
)

const (
	FactACActivePower           FactID = "pv.ac.power.active"
	FactACApparentPower         FactID = "pv.ac.power.apparent"
	FactACReactivePower         FactID = "pv.ac.power.reactive"
	FactACPowerFactor           FactID = "pv.ac.power_factor"
	FactACFrequency             FactID = "pv.ac.frequency"
	FactACCurrent               FactID = "pv.ac.current"
	FactACVoltageLineToNeutral  FactID = "pv.ac.voltage.line_to_neutral"
	FactACVoltageLineToLine     FactID = "pv.ac.voltage.line_to_line"
	FactEnergyActiveExportTotal FactID = "pv.energy.active_export_total"
	FactDCCurrent               FactID = "pv.dc.current"
	FactDCVoltage               FactID = "pv.dc.voltage"
	FactDCActivePower           FactID = "pv.dc.power.active"
	FactDCEnergyActiveTotal     FactID = "pv.dc.energy.active_total"
	FactTemperature             FactID = "pv.temperature"
	FactOperatingState          FactID = "pv.operating.state"
	FactEventFlags              FactID = "pv.event.flags"
	FactRatingACActivePower     FactID = "pv.rating.ac.active_power"

	OperatingStateUnknown      = "UNKNOWN"
	OperatingStateOff          = "OFF"
	OperatingStateStandby      = "STANDBY"
	OperatingStateStarting     = "STARTING"
	OperatingStateOperating    = "OPERATING"
	OperatingStateDerated      = "DERATED"
	OperatingStateFault        = "FAULT"
	OperatingStateShuttingDown = "SHUTTING_DOWN"
)

type dimensionKind string

const (
	dimensionScope     dimensionKind = "scope"
	dimensionPhase     dimensionKind = "phase"
	dimensionPhasePair dimensionKind = "phase_pair"
	dimensionInputID   dimensionKind = "input_id"
	dimensionSensorID  dimensionKind = "sensor_id"
)

type FactDefinition struct {
	ID          FactID
	Kind        ValueKind
	Unit        Unit
	Dimension   dimensionKind
	Policy      PolicyID
	Accumulator bool
	Domain      map[string]struct{}
}

type Catalog struct {
	ContractID           string
	Facts                map[FactID]FactDefinition
	AdditiveFactsAllowed bool
	WriteAuthority       bool
	SourceCompatibility  map[string]string
}

var identifierPattern = regexp.MustCompile(`^[A-Za-z0-9._:-]+$`)

func CatalogV1() Catalog {
	definitions := []FactDefinition{
		decimalDefinition(FactACActivePower, UnitWatt, dimensionScope, PolicyTelemetryFastV1, false),
		decimalDefinition(FactACApparentPower, UnitVoltAmp, dimensionScope, PolicyTelemetryFastV1, false),
		decimalDefinition(FactACReactivePower, UnitVar, dimensionScope, PolicyTelemetryFastV1, false),
		decimalDefinition(FactACPowerFactor, UnitOne, dimensionScope, PolicyTelemetryFastV1, false),
		decimalDefinition(FactACFrequency, UnitHertz, dimensionScope, PolicyTelemetryFastV1, false),
		decimalDefinition(FactACCurrent, UnitAmpere, dimensionPhase, PolicyTelemetryFastV1, false),
		decimalDefinition(FactACVoltageLineToNeutral, UnitVolt, dimensionPhase, PolicyTelemetryFastV1, false),
		decimalDefinition(FactACVoltageLineToLine, UnitVolt, dimensionPhasePair, PolicyTelemetryFastV1, false),
		decimalDefinition(FactEnergyActiveExportTotal, UnitWattHour, dimensionScope, PolicyAccumulatorV1, true),
		decimalDefinition(FactDCCurrent, UnitAmpere, dimensionInputID, PolicyTelemetryFastV1, false),
		decimalDefinition(FactDCVoltage, UnitVolt, dimensionInputID, PolicyTelemetryFastV1, false),
		decimalDefinition(FactDCActivePower, UnitWatt, dimensionInputID, PolicyTelemetryFastV1, false),
		decimalDefinition(FactDCEnergyActiveTotal, UnitWattHour, dimensionInputID, PolicyAccumulatorV1, true),
		decimalDefinition(FactTemperature, UnitCelsius, dimensionSensorID, PolicyTelemetryFastV1, false),
		{ID: FactOperatingState, Kind: ValueKindEnum, Unit: UnitOne, Dimension: dimensionScope, Policy: PolicyStatusV1, Domain: stringSet(
			OperatingStateUnknown, OperatingStateOff, OperatingStateStandby, OperatingStateStarting,
			OperatingStateOperating, OperatingStateDerated, OperatingStateFault, OperatingStateShuttingDown,
		)},
		{ID: FactEventFlags, Kind: ValueKindBitfield, Unit: UnitOne, Dimension: dimensionScope, Policy: PolicyStatusV1, Domain: stringSet(
			"GROUND_FAULT", "DC_OVER_VOLTAGE", "AC_DISCONNECT", "DC_DISCONNECT", "GRID_DISCONNECT",
			"CABINET_OPEN", "MANUAL_SHUTDOWN", "OVER_TEMPERATURE", "FREQUENCY_OUT_OF_RANGE",
			"VOLTAGE_OUT_OF_RANGE", "COMMUNICATION_FAULT", "INTERNAL_FAULT",
		)},
		decimalDefinition(FactRatingACActivePower, UnitWatt, dimensionScope, PolicyRatingV1, false),
	}
	facts := make(map[FactID]FactDefinition, len(definitions))
	for _, definition := range definitions {
		facts[definition.ID] = definition
	}
	return Catalog{
		ContractID:           ContractV1,
		Facts:                facts,
		AdditiveFactsAllowed: false,
		WriteAuthority:       false,
		SourceCompatibility: map[string]string{
			"sunspec.phase1@1.0.0":                          "",
			"sunspec.inverter.three_phase.monitoring@1.0.0": "",
		},
	}
}

func decimalDefinition(id FactID, unit Unit, dimension dimensionKind, policy PolicyID, accumulator bool) FactDefinition {
	return FactDefinition{ID: id, Kind: ValueKindDecimal, Unit: unit, Dimension: dimension, Policy: policy, Accumulator: accumulator}
}

func stringSet(values ...string) map[string]struct{} {
	set := make(map[string]struct{}, len(values))
	for _, value := range values {
		set[value] = struct{}{}
	}
	return set
}

func (catalog Catalog) ValidateCandidate(candidate FactCandidate) error {
	definition, ok := catalog.Facts[candidate.ID]
	if !ok || definition.Unit != candidate.Unit || definition.Kind != candidate.Value.Kind {
		return ErrInvalidFact
	}
	if err := validateDimensions(definition.Dimension, candidate.Dimensions); err != nil {
		return err
	}
	switch definition.Kind {
	case ValueKindDecimal:
		if candidate.Value.Decimal == nil || candidate.Value.Symbol != "" || len(candidate.Value.Symbols) != 0 || candidate.Value.Decimal.Validate() != nil {
			return ErrInvalidFact
		}
	case ValueKindEnum:
		if candidate.Value.Decimal != nil || len(candidate.Value.Symbols) != 0 {
			return ErrInvalidFact
		}
		if _, ok := definition.Domain[candidate.Value.Symbol]; !ok {
			return ErrInvalidFact
		}
	case ValueKindBitfield:
		if candidate.Value.Decimal != nil || candidate.Value.Symbol != "" {
			return ErrInvalidFact
		}
		seen := make(map[string]struct{}, len(candidate.Value.Symbols))
		for _, symbol := range candidate.Value.Symbols {
			if _, ok := definition.Domain[symbol]; !ok {
				return ErrInvalidFact
			}
			if _, duplicate := seen[symbol]; duplicate {
				return ErrInvalidFact
			}
			seen[symbol] = struct{}{}
		}
	default:
		return ErrInvalidFact
	}
	return nil
}

func validateDimensions(kind dimensionKind, dimensions Dimensions) error {
	nonempty := 0
	for _, present := range []bool{
		dimensions.Scope != "", dimensions.Phase != "", dimensions.PhasePair != "",
		dimensions.InputID != "", dimensions.SensorID != "",
	} {
		if present {
			nonempty++
		}
	}
	if nonempty != 1 {
		return ErrInvalidFact
	}
	switch kind {
	case dimensionScope:
		if dimensions.Scope != ScopeTotal {
			return ErrInvalidFact
		}
	case dimensionPhase:
		if dimensions.Phase != PhaseL1 && dimensions.Phase != PhaseL2 && dimensions.Phase != PhaseL3 {
			return ErrInvalidFact
		}
	case dimensionPhasePair:
		if dimensions.PhasePair != PhasePairL1L2 && dimensions.PhasePair != PhasePairL2L3 && dimensions.PhasePair != PhasePairL3L1 {
			return ErrInvalidFact
		}
	case dimensionInputID:
		if !validOpaqueIdentifier(dimensions.InputID) {
			return ErrInvalidFact
		}
	case dimensionSensorID:
		if !validOpaqueIdentifier(dimensions.SensorID) {
			return ErrInvalidFact
		}
	default:
		return ErrInvalidFact
	}
	return nil
}

func validOpaqueIdentifier(value string) bool {
	if len(value) == 0 || len(value) > 64 || !identifierPattern.MatchString(value) || net.ParseIP(value) != nil {
		return false
	}
	if strings.Contains(value, "://") {
		return false
	}
	if _, _, err := net.SplitHostPort(value); err == nil {
		return false
	}
	return true
}

func (catalog Catalog) ThreePhaseTelemetrySatisfied(facts map[FactKey]Fact) bool {
	for _, key := range threePhaseRequiredKeys() {
		fact, ok := facts[key]
		if !ok || (fact.Availability != AvailabilityAvailable && fact.Availability != AvailabilityUnavailable) {
			return false
		}
	}
	return true
}

func threePhaseRequiredKeys() []FactKey {
	keys := []FactKey{
		NewFactKey(FactACActivePower, Dimensions{Scope: ScopeTotal}),
		NewFactKey(FactACFrequency, Dimensions{Scope: ScopeTotal}),
		NewFactKey(FactEnergyActiveExportTotal, Dimensions{Scope: ScopeTotal}),
		NewFactKey(FactOperatingState, Dimensions{Scope: ScopeTotal}),
	}
	for _, phase := range []Phase{PhaseL1, PhaseL2, PhaseL3} {
		keys = append(keys,
			NewFactKey(FactACCurrent, Dimensions{Phase: phase}),
			NewFactKey(FactACVoltageLineToNeutral, Dimensions{Phase: phase}),
		)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	return keys
}

func (definition FactDefinition) validatePolicy(policy PolicyID) error {
	if definition.Policy != policy {
		return fmt.Errorf("policy %q for %q: %w", policy, definition.ID, ErrInvalidFact)
	}
	return nil
}
