package pv

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
)

const ContractV1 = "helianthus.canonical-pv/v1"

var (
	ErrInvalidDecimal       = errors.New("pv: invalid decimal")
	ErrInvalidFact          = errors.New("pv: invalid fact")
	ErrInvalidMonotonicTime = errors.New("pv: invalid monotonic time")
	ErrInvalidCounter       = errors.New("pv: invalid counter evidence")
	ErrSourceNotAdmitted    = errors.New("pv: source not admitted")
	ErrAssetNotFound        = errors.New("pv: asset not found")
)

type FactID string
type Unit string
type ValueKind string
type PolicyID string
type Quality string
type Availability string
type Freshness string
type ContinuityState string
type CounterEventKind string
type Digest string
type MonotonicNanos int64
type Phase string
type PhasePair string
type Scope string
type FactKey string

const (
	ValueKindDecimal  ValueKind = "decimal"
	ValueKindEnum     ValueKind = "enum"
	ValueKindBitfield ValueKind = "bitfield"

	UnitWatt     Unit = "W"
	UnitVoltAmp  Unit = "VA"
	UnitVar      Unit = "var"
	UnitVolt     Unit = "V"
	UnitAmpere   Unit = "A"
	UnitHertz    Unit = "Hz"
	UnitWattHour Unit = "Wh"
	UnitCelsius  Unit = "Cel"
	UnitOne      Unit = "1"

	QualityGood    Quality = "GOOD"
	QualitySuspect Quality = "SUSPECT"
	QualityBad     Quality = "BAD"

	AvailabilityAvailable   Availability = "AVAILABLE"
	AvailabilityUnavailable Availability = "UNAVAILABLE"
	AvailabilityUnsupported Availability = "UNSUPPORTED"

	FreshnessFresh   Freshness = "FRESH"
	FreshnessStale   Freshness = "STALE"
	FreshnessExpired Freshness = "EXPIRED"

	ContinuityBaseline      ContinuityState = "BASELINE"
	ContinuityContiguous    ContinuityState = "CONTIGUOUS"
	ContinuityRollover      ContinuityState = "ROLLOVER"
	ContinuityReset         ContinuityState = "RESET"
	ContinuityDiscontinuity ContinuityState = "DISCONTINUITY"

	CounterEventNone     CounterEventKind = ""
	CounterEventReset    CounterEventKind = "RESET"
	CounterEventRollover CounterEventKind = "ROLLOVER"

	PhaseL1 Phase = "L1"
	PhaseL2 Phase = "L2"
	PhaseL3 Phase = "L3"

	PhasePairL1L2 PhasePair = "L1_L2"
	PhasePairL2L3 PhasePair = "L2_L3"
	PhasePairL3L1 PhasePair = "L3_L1"

	ScopeTotal Scope = "total"

	SourceTerminalVerified = "terminal_verified"
)

var (
	decimalPattern = regexp.MustCompile(`^-?(0|[1-9][0-9]*)$`)
	digestPattern  = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)
)

type Decimal struct {
	Coefficient string
	Scale       int
}

func NewDecimal(coefficient string, scale int) (Decimal, error) {
	if !decimalPattern.MatchString(coefficient) || coefficient == "-0" || scale < -18 || scale > 18 {
		return Decimal{}, ErrInvalidDecimal
	}
	return Decimal{Coefficient: coefficient, Scale: scale}, nil
}

func MustDecimal(coefficient string, scale int) Decimal {
	decimal, err := NewDecimal(coefficient, scale)
	if err != nil {
		panic(err)
	}
	return decimal
}

func (decimal Decimal) Validate() error {
	_, err := NewDecimal(decimal.Coefficient, decimal.Scale)
	return err
}

type FactValue struct {
	Kind    ValueKind
	Decimal *Decimal
	Symbol  string
	Symbols []string
}

func DecimalFactValue(decimal Decimal) FactValue {
	copy := decimal
	return FactValue{Kind: ValueKindDecimal, Decimal: &copy}
}

func EnumFactValue(symbol string) FactValue {
	return FactValue{Kind: ValueKindEnum, Symbol: symbol}
}

func BitfieldFactValue(symbols ...string) FactValue {
	return FactValue{Kind: ValueKindBitfield, Symbols: append([]string(nil), symbols...)}
}

type Dimensions struct {
	Scope     Scope
	Phase     Phase
	PhasePair PhasePair
	InputID   string
	SensorID  string
}

func (dimensions Dimensions) identity() string {
	parts := make([]string, 0, 5)
	if dimensions.Scope != "" {
		parts = append(parts, "scope="+string(dimensions.Scope))
	}
	if dimensions.Phase != "" {
		parts = append(parts, "phase="+string(dimensions.Phase))
	}
	if dimensions.PhasePair != "" {
		parts = append(parts, "phase_pair="+string(dimensions.PhasePair))
	}
	if dimensions.InputID != "" {
		parts = append(parts, "input_id="+dimensions.InputID)
	}
	if dimensions.SensorID != "" {
		parts = append(parts, "sensor_id="+dimensions.SensorID)
	}
	sort.Strings(parts)
	return strings.Join(parts, ",")
}

func NewFactKey(id FactID, dimensions Dimensions) FactKey {
	return FactKey(string(id) + "|" + dimensions.identity())
}

type FactCandidate struct {
	ID         FactID
	Dimensions Dimensions
	Value      FactValue
	Unit       Unit
}

type Temporal struct {
	Receipt     MonotonicNanos
	FreshUntil  MonotonicNanos
	RetainUntil MonotonicNanos
	Policy      PolicyID
}

type TemporalState struct {
	Availability Availability
	Freshness    Freshness
	Temporal     Temporal
}

type Continuity struct {
	State       ContinuityState
	Delta       *Decimal
	Modulus     *Decimal
	EvidenceRef Digest
}

type CounterEvidence struct {
	Kind             CounterEventKind
	Modulus          *Decimal
	BoundaryVerified bool
	EvidenceRef      Digest
}

func (evidence CounterEvidence) empty() bool {
	return evidence.Kind == CounterEventNone && evidence.Modulus == nil && !evidence.BoundaryVerified && evidence.EvidenceRef == ""
}

type SourceIdentity struct {
	Protocol       string
	ProfileID      string
	ProfileVersion string
	Validity       string
}

type Provenance struct {
	SourceIdentity
	SourceRegistryRef    Digest
	SourceObservationRef Digest
	SourceShadowRef      Digest
	EvidenceRef          Digest
}

type SourceResolver interface {
	ResolveSource(SourceIdentity) (Digest, bool)
}

type FactInput struct {
	Candidate FactCandidate
	Quality   Quality
	Receipt   MonotonicNanos
	Policy    PolicyID
	Counter   CounterEvidence
}

type Fact struct {
	ID           FactID
	Dimensions   Dimensions
	Value        FactValue
	Unit         Unit
	Quality      Quality
	Availability Availability
	Freshness    Freshness
	Temporal     Temporal
	OriginRef    Digest
	Continuity   *Continuity
}

type Update struct {
	AssetRef  string
	Evaluated MonotonicNanos
	Source    Provenance
	Facts     []FactInput
}

type Snapshot struct {
	ContractID            string
	AssetRef              string
	Generation            uint64
	Evaluated             MonotonicNanos
	Source                Provenance
	Facts                 map[FactKey]Fact
	Origins               map[Digest]Provenance
	ThreePhaseTelemetryV1 bool
}

func (digest Digest) Validate() error {
	if !digestPattern.MatchString(string(digest)) {
		return fmt.Errorf("digest %q: %w", digest, ErrSourceNotAdmitted)
	}
	return nil
}
