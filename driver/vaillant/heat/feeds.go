package heat

import (
	"math/big"
	"strconv"
	"strings"

	semreg "github.com/Project-Helianthus/helianthus-semreg/semreg/v1"
	"github.com/Project-Helianthus/helianthus-semreg/semreg/v1/packs/thermal"
)

// Source ranks are an encoding of which device produces a fact. The device
// that produces the state ranks above a second copy of it. They are not a
// score and they do not select a core value by themselves.
const (
	sourceRankPrimary   = 3
	sourceRankSecondary = 2
	sourceRankTertiary  = 1
)

// Status is the result of applying one feed row to one reading.
const (
	StatusFed        = "fed"
	StatusUnknown    = "unknown"
	StatusWithdrawn  = "withdrawn"
	StatusNamespace  = "namespace"
	StatusNone       = "none"
	StatusUnmatched  = "unmatched"
	StatusProjection = "projection"
	StatusUnchanged  = "unchanged"
	StatusConflict   = "conflict"
)

// NamespaceEbusreg is the driver namespace for a value the thermal pack
// cannot express.
const NamespaceEbusreg = "ebusreg"

// Feed is one table row. It binds a native device row to a thermal fact, or
// names the gap that keeps the value in the driver namespace.
type Feed struct {
	ID          string
	Subject     string
	Group       byte
	Register    uint16
	Selector    string
	InstanceMin byte
	InstanceMax byte
	NativeName  string
	NativeUnit  string
	Element     string
	Conversion  string
	Grade       string
	Families    []string
	Codes       map[string][]FeedCode
	SourceRank  int
	KeepNative  bool
	Gap         string
	Hold        string
	Derived     bool
	NoMatch     bool
	Poll        string
}

// FeedCode is one documented code of a family value table. Element, when set,
// is the pack field this code feeds and Token is that field's symbol.
type FeedCode struct {
	Value   int
	Label   string
	Token   string
	Element string
	None    bool
}

// Reading is one decoded device row. A failed or absent reply sets Failed or
// leaves Present false. Missing is a zone that left the inventory.
// Deactivated is a circuit the controller reports off. Code is set for an
// enumeration; Text is the decoded native quantity or label.
type Reading struct {
	Present     bool
	Failed      bool
	Deactivated bool
	Missing     bool
	Stale       bool
	Text        string
	Code        *int
	Family      string
	Instance    byte
	Group       byte
	Register    uint16
	Selector    string
	Resource    string
	Poll        string
}

// Outcome is the fact, the native evidence, or the gap for one row.
type Outcome struct {
	Feed       Feed
	Status     string
	Value      *semreg.Value
	Key        semreg.FactKey
	Evidence   string
	Namespace  string
	NativeText string
	Gap        string
	Freshness  string
	Grade      string
	Rank       int
}

// ReadingOf returns a present reading of text at the feed's own locator.
func ReadingOf(feed Feed, text, resource string) Reading {
	return Reading{
		Present: true, Text: text, Family: firstFamily(feed), Instance: feed.InstanceMin,
		Group: feed.Group, Register: feed.Register, Selector: feed.Selector, Resource: resource,
		Poll: feed.Poll,
	}
}

// CodeReading returns a present enumeration reading.
func CodeReading(feed Feed, code int, resource string) Reading {
	reading := ReadingOf(feed, "", resource)
	reading.Code = &code
	return reading
}

func firstFamily(feed Feed) string {
	if len(feed.Families) == 0 {
		return ""
	}
	return feed.Families[0]
}

// Matches reports whether the reading is this row's group, register or
// selector, instance and family. A row with no locator never matches.
func (f Feed) Matches(reading Reading) bool {
	if f.NoMatch || f.Derived {
		return false
	}
	if f.Selector != "" {
		if reading.Selector != f.Selector {
			return false
		}
	} else if reading.Group != f.Group || reading.Register != f.Register {
		return false
	}
	if reading.Instance < f.InstanceMin || reading.Instance > f.InstanceMax {
		return false
	}
	if len(f.Families) > 0 && !containsString(f.Families, reading.Family) {
		return false
	}
	return true
}

// Apply binds one reading to one row. An unmatched reading publishes nothing.
// A failed or absent reply publishes no fact and leaves the element unknown.
// Withdrawal is a zone missing from the inventory or a deactivated circuit.
// The outcome key names the fact and the resource that leave, so two
// withdrawals of the same element stay distinct. A value the pack cannot
// express stays in the driver namespace with its native value and unit.
func Apply(feed Feed, reading Reading) Outcome {
	out := Outcome{Feed: feed, Status: StatusUnmatched, Grade: feed.Grade, Rank: feed.SourceRank, Gap: feed.Gap}
	if !feed.Matches(reading) {
		return out
	}
	if (reading.Missing && feed.Subject == "zone") || (reading.Deactivated && feed.Subject == "circuit") {
		out.Status = StatusWithdrawn
		out.Key = factKeyIfElement(feed, reading.Resource)
		out.Evidence = evidenceOf(feed, reading, "", "")
		return out
	}
	if reading.Failed || !reading.Present {
		out.Status = StatusUnknown
		out.Freshness = freshnessOf(reading)
		out.Key = factKeyIfElement(feed, reading.Resource)
		out.Evidence = evidenceOf(feed, reading, reading.Text, "")
		out.NativeText = reading.Text
		if feed.Hold != "" {
			out.Gap = feed.Hold
		}
		return out
	}
	if feed.Hold != "" {
		out.Status = StatusNamespace
		out.Namespace = NamespaceEbusreg
		out.Gap = feed.Hold
		out.Freshness = freshnessOf(reading)
		out.Key = factKeyIfElement(feed, reading.Resource)
		out.Evidence = evidenceOf(feed, reading, reading.Text, "")
		out.NativeText = reading.Text
		return out
	}
	out.Freshness = freshnessOf(reading)
	switch feed.Conversion {
	case "celsius", "percent", "bar", "valve":
		return applyQuantity(feed, reading, out)
	case "symbol":
		return applySymbol(feed, reading, out)
	case "maintenance":
		return applyMaintenance(feed, reading, out)
	case "capability":
		out.Status = StatusFed
		out.NativeText = feed.Element
		out.Evidence = evidenceOf(feed, reading, feed.Element, "")
		return out
	default:
		return applyNamespace(feed, reading, out, reading.Text, feed.Gap)
	}
}

func applyQuantity(feed Feed, reading Reading, out Outcome) Outcome {
	number, ok := parseNativeDecimal(reading.Text)
	if !ok {
		out.Status = StatusUnchanged
		out.Key = factKeyIfElement(feed, reading.Resource)
		out.Evidence = evidenceOf(feed, reading, reading.Text, "")
		return out
	}
	converted, unit, ok := convertNative(feed.Conversion, number)
	if !ok {
		out.Status = StatusUnchanged
		out.Key = factKeyIfElement(feed, reading.Resource)
		out.Evidence = evidenceOf(feed, reading, reading.Text, "")
		return out
	}
	if feed.Element == "" || !subjectAccepts(feed) {
		return applyNamespace(feed, reading, out, reading.Text, gapOr(feed, "catalog subject does not express this row"))
	}
	value := &semreg.Value{Kind: semreg.ValueQuantity, Quantity: &semreg.Quantity{Number: converted, Unit: unit}}
	key := factKey(feed, reading.Resource)
	if err := thermal.New().ValidateFact(key, value); err != nil {
		return applyNamespace(feed, reading, out, reading.Text, "value outside the pack range")
	}
	out.Status = StatusFed
	out.Value = value
	out.Key = key
	out.NativeText = reading.Text
	out.Evidence = evidenceOf(feed, reading, reading.Text, "")
	if feed.KeepNative {
		out.Namespace = NamespaceEbusreg
	}
	return out
}

func applySymbol(feed Feed, reading Reading, out Outcome) Outcome {
	label, token, element, none, known := lookupCode(feed, reading)
	if !known {
		raw := reading.Text
		if reading.Code != nil {
			raw = strconv.Itoa(*reading.Code)
		}
		out.Status = StatusProjection
		out.Namespace = NamespaceEbusreg
		out.NativeText = raw
		out.Gap = "code has no symbol in the documented table"
		out.Key = factKeyIfElement(feed, reading.Resource)
		out.Evidence = evidenceOf(feed, reading, raw, "")
		return out
	}
	if element != "" {
		feed.Element = element
	}
	if none || token == "" {
		status := StatusProjection
		gap := "documented label has no pack token"
		if none {
			status = StatusNone
			gap = ""
		}
		out.Status = status
		out.Namespace = NamespaceEbusreg
		out.NativeText = label
		out.Gap = gapOr(feed, gap)
		out.Key = factKeyIfElement(feed, reading.Resource)
		out.Evidence = evidenceOf(feed, reading, label, "")
		return out
	}
	if feed.Element == "" || !subjectAccepts(feed) {
		return applyNamespace(feed, reading, out, label, gapOr(feed, "catalog subject does not express this row"))
	}
	value := &semreg.Value{Kind: semreg.ValueSymbol, Symbol: &semreg.Symbol{Namespace: semreg.DefinitionID(feed.Element), Token: token, Known: true}}
	key := factKey(feed, reading.Resource)
	if err := thermal.New().ValidateFact(key, value); err != nil {
		return applyNamespace(feed, reading, out, label, "pack rejects the token")
	}
	out.Status = StatusFed
	out.Value = value
	out.Key = key
	out.NativeText = label
	out.Evidence = evidenceOf(feed, reading, label, token)
	if feed.KeepNative {
		out.Namespace = NamespaceEbusreg
	}
	return out
}

func applyMaintenance(feed Feed, reading Reading, out Outcome) Outcome {
	var token string
	switch reading.Text {
	case "1", "true":
		token = "due"
	case "0", "false":
		token = "clear"
	default:
		if reading.Code == nil {
			out.Status = StatusUnchanged
			out.Key = factKeyIfElement(feed, reading.Resource)
			out.Evidence = evidenceOf(feed, reading, reading.Text, "")
			return out
		}
		switch *reading.Code {
		case 1:
			token = "due"
		case 0:
			token = "clear"
		default:
			out.Status = StatusProjection
			out.Namespace = NamespaceEbusreg
			out.NativeText = strconv.Itoa(*reading.Code)
			out.Gap = "code has no symbol in the documented table"
			out.Key = factKeyIfElement(feed, reading.Resource)
			out.Evidence = evidenceOf(feed, reading, out.NativeText, "")
			return out
		}
	}
	value := &semreg.Value{Kind: semreg.ValueSymbol, Symbol: &semreg.Symbol{Namespace: semreg.DefinitionID(feed.Element), Token: token, Known: true}}
	key := factKey(feed, reading.Resource)
	if err := thermal.New().ValidateFact(key, value); err != nil {
		return applyNamespace(feed, reading, out, reading.Text, "pack rejects the token")
	}
	out.Status = StatusFed
	out.Value = value
	out.Key = key
	out.NativeText = reading.Text
	out.Evidence = evidenceOf(feed, reading, reading.Text, token)
	return out
}

func applyNamespace(feed Feed, reading Reading, out Outcome, text, gap string) Outcome {
	out.Status = StatusNamespace
	out.Namespace = NamespaceEbusreg
	out.NativeText = text
	out.Gap = gap
	out.Evidence = evidenceOf(feed, reading, text, "")
	out.Value = nil
	out.Key = factKeyIfElement(feed, reading.Resource)
	return out
}

func factKeyIfElement(feed Feed, resource string) semreg.FactKey {
	if feed.Element == "" || !subjectAccepts(feed) {
		return semreg.FactKey{}
	}
	return factKey(feed, resource)
}

func fedSymbol(out Outcome, element, token string) Outcome {
	out.Status = StatusFed
	out.Value = symbolValue(element, token)
	if err := thermal.New().ValidateFact(out.Key, out.Value); err != nil {
		out.Status = StatusNamespace
		out.Namespace = NamespaceEbusreg
		out.Value = nil
		out.Gap = "pack rejects the token"
	}
	return out
}

func gapOr(feed Feed, fallback string) string {
	if feed.Gap != "" {
		return feed.Gap
	}
	return fallback
}

func lookupCode(feed Feed, reading Reading) (label, token, element string, none, known bool) {
	if reading.Code == nil {
		if reading.Text == "" {
			return "", "", "", false, false
		}
		for _, family := range codeFamilies(feed, reading.Family) {
			for _, code := range feed.Codes[family] {
				if code.Label == reading.Text {
					return code.Label, code.Token, code.Element, code.None, true
				}
			}
		}
		return "", "", "", false, false
	}
	families := codeFamilies(feed, reading.Family)
	if len(families) == 0 {
		return "", "", "", false, false
	}
	for _, code := range feed.Codes[families[0]] {
		if code.Value == *reading.Code {
			return code.Label, code.Token, code.Element, code.None, true
		}
	}
	return "", "", "", false, false
}

func codeFamilies(feed Feed, family string) []string {
	if family == "" {
		return nil
	}
	if _, ok := feed.Codes[family]; ok {
		return []string{family}
	}
	return nil
}

func subjectAccepts(feed Feed) bool {
	want, ok := subjectDimensions[feed.Subject]
	if !ok || feed.Element == "" {
		return false
	}
	got, known := elementDimension(semreg.DefinitionID(feed.Element))
	return known && got == want
}

var subjectDimensions = map[string]semreg.DefinitionID{
	"zone":        "thermal.dimension.zone",
	"dhw":         "thermal.dimension.dhw",
	"circuit":     "thermal.dimension.circuit",
	"system":      "thermal.dimension.system",
	"terminal":    "thermal.dimension.terminal_unit",
	"refrigerant": "thermal.dimension.refrigerant_circuit",
	"compressor":  "thermal.dimension.compressor",
	"generator":   "thermal.dimension.heat_generator",
	"store":       "thermal.dimension.thermal_store",
}

func elementDimension(id semreg.DefinitionID) (semreg.DefinitionID, bool) {
	for _, field := range thermal.Metadata().Fields {
		if field.Ref.ID == id {
			return field.Dimension.ID, true
		}
	}
	return "", false
}

func factKey(feed Feed, resource string) semreg.FactKey {
	pack := thermal.New().Pack()
	dimension, _ := elementDimension(semreg.DefinitionID(feed.Element))
	if dimension == "" {
		dimension = subjectDimensions[feed.Subject]
	}
	text := resource
	return semreg.FactKey{
		PackID: pack.ID, PackVersion: pack.Version, FactID: semreg.DefinitionID(feed.Element),
		Dimensions: []semreg.Dimension{{
			ID: dimension, Value: semreg.Value{Kind: semreg.ValueText, Text: &text},
		}},
	}
}

func evidenceOf(feed Feed, reading Reading, native, token string) string {
	poll := reading.Poll
	if poll == "" {
		poll = "unknown"
	}
	unit := feed.NativeUnit
	if unit == "" {
		unit = "none"
	}
	parts := []string{
		"native=" + native,
		"unit=" + unit,
		"name=" + feed.NativeName,
		"grade=" + feed.Grade,
		"rank=" + strconv.Itoa(feed.SourceRank),
		"poll=" + poll,
		"freshness=" + freshnessOf(reading),
	}
	if token != "" {
		parts = append(parts, "token="+token)
	}
	if reading.Code != nil {
		parts = append(parts, "code="+strconv.Itoa(*reading.Code))
	}
	return strings.Join(parts, " ")
}

func freshnessOf(reading Reading) string {
	if reading.Stale {
		return "stale"
	}
	if reading.Failed || !reading.Present {
		return "unknown"
	}
	return "fresh"
}

func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

// Select keeps the highest-ranked source of each fact. An unknown or stale
// top source is returned as that source and is not replaced by a lower one.
func Select(outcomes []Outcome) []Outcome {
	type group struct {
		order []Outcome
	}
	groups := map[string]*group{}
	var keys []string
	var rest []Outcome
	for _, outcome := range outcomes {
		if outcome.Status == StatusUnmatched || outcome.Key.FactID == "" {
			rest = append(rest, outcome)
			continue
		}
		id := string(outcome.Key.FactID) + "\x00" + string(outcome.Key.Dimensions[0].ID) + "\x00" + dimensionText(outcome.Key)
		item := groups[id]
		if item == nil {
			item = &group{}
			groups[id] = item
			keys = append(keys, id)
		}
		item.order = append(item.order, outcome)
	}
	for _, id := range keys {
		rest = append(rest, selectGroup(groups[id].order)...)
	}
	return rest
}

func selectGroup(outcomes []Outcome) []Outcome {
	ranked := map[int][]Outcome{}
	var ranks []int
	for _, outcome := range outcomes {
		if _, ok := ranked[outcome.Rank]; !ok {
			ranks = append(ranks, outcome.Rank)
		}
		ranked[outcome.Rank] = append(ranked[outcome.Rank], outcome)
	}
	for i := 0; i < len(ranks); i++ {
		for j := i + 1; j < len(ranks); j++ {
			if ranks[j] > ranks[i] {
				ranks[i], ranks[j] = ranks[j], ranks[i]
			}
		}
	}
	for _, rank := range ranks {
		bucket := ranked[rank]
		if blocking(bucket) {
			return bucket
		}
		fed := fedOnly(bucket)
		if len(fed) > 0 {
			return uniqueFed(fed)
		}
	}
	if len(ranks) == 0 {
		return nil
	}
	return ranked[ranks[len(ranks)-1]]
}

func blocking(outcomes []Outcome) bool {
	sawFed := false
	sawBlock := false
	for _, outcome := range outcomes {
		switch outcome.Status {
		case StatusFed:
			sawFed = true
		case StatusUnknown, StatusProjection, StatusUnchanged:
			sawBlock = true
		case StatusNamespace:
			if outcome.Value == nil {
				sawBlock = true
			}
		}
		if outcome.Freshness == "stale" && outcome.Status == StatusFed {
			sawBlock = true
			sawFed = true
		}
	}
	if sawFed && !sawBlock {
		return false
	}
	return sawBlock
}

func fedOnly(outcomes []Outcome) []Outcome {
	var fed []Outcome
	for _, outcome := range outcomes {
		if outcome.Status == StatusFed && outcome.Freshness != "stale" {
			fed = append(fed, outcome)
		}
	}
	return fed
}

func uniqueFed(outcomes []Outcome) []Outcome {
	if len(outcomes) < 2 {
		return outcomes
	}
	first := outcomes[0]
	for _, outcome := range outcomes[1:] {
		if !sameValue(first.Value, outcome.Value) {
			conflict := first
			conflict.Status = StatusConflict
			conflict.Value = nil
			return []Outcome{conflict}
		}
	}
	return []Outcome{first}
}

func sameValue(left, right *semreg.Value) bool {
	if left == nil || right == nil || left.Kind != right.Kind {
		return false
	}
	if left.Kind == semreg.ValueQuantity && left.Quantity != nil && right.Quantity != nil {
		return left.Quantity.Unit == right.Quantity.Unit && left.Quantity.Number == right.Quantity.Number
	}
	if left.Kind == semreg.ValueSymbol && left.Symbol != nil && right.Symbol != nil {
		return left.Symbol.Namespace == right.Symbol.Namespace && left.Symbol.Token == right.Symbol.Token
	}
	return false
}

func dimensionText(key semreg.FactKey) string {
	if len(key.Dimensions) == 0 || key.Dimensions[0].Value.Text == nil {
		return ""
	}
	return *key.Dimensions[0].Value.Text
}

func parseNativeDecimal(text string) (semreg.Decimal, bool) {
	text = strings.TrimSpace(text)
	if text == "" || strings.ContainsAny(text, "eE+ ") {
		return semreg.Decimal{}, false
	}
	negative := false
	if strings.HasPrefix(text, "-") {
		negative = true
		text = text[1:]
		if text == "" {
			return semreg.Decimal{}, false
		}
	}
	whole, frac, found := strings.Cut(text, ".")
	if strings.Contains(frac, ".") || !digitsOnly(whole) || !digitsOnly(frac) {
		return semreg.Decimal{}, false
	}
	if whole == "" {
		whole = "0"
	}
	if !found {
		frac = ""
	}
	coefficient := strings.TrimLeft(whole+frac, "0")
	if coefficient == "" {
		return semreg.Decimal{Coefficient: "0"}, true
	}
	exponent := -int32(len(frac))
	if negative {
		coefficient = "-" + coefficient
	}
	return canonicalDecimal(coefficient, exponent)
}

func digitsOnly(text string) bool {
	if text == "" {
		return true
	}
	for _, r := range text {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func canonicalDecimal(coefficient string, exponent int32) (semreg.Decimal, bool) {
	for strings.HasSuffix(coefficient, "0") && coefficient != "-0" && coefficient != "0" {
		coefficient = coefficient[:len(coefficient)-1]
		exponent++
	}
	if coefficient == "-0" {
		coefficient = "0"
		exponent = 0
	}
	if coefficient == "0" {
		exponent = 0
	}
	decimal := semreg.Decimal{Coefficient: coefficient, Exponent10: exponent}
	if decimal.Validate() != nil {
		return semreg.Decimal{}, false
	}
	return decimal, true
}

func convertNative(conversion string, number semreg.Decimal) (semreg.Decimal, semreg.DefinitionID, bool) {
	switch conversion {
	case "celsius":
		sum, ok := decimalShift(number, 27315, -2)
		return sum, "unit.kelvin", ok
	case "percent":
		ratio, ok := decimalShiftExponent(number, -2)
		return ratio, "unit.ratio", ok
	case "bar":
		product, ok := decimalShiftExponent(number, 5)
		return product, "unit.pascal", ok
	case "valve":
		if number.Coefficient != "0" && number.Coefficient != "1" {
			return semreg.Decimal{}, "", false
		}
		if number.Exponent10 != 0 && number.Coefficient != "0" {
			return semreg.Decimal{}, "", false
		}
		return number, "unit.ratio", number.Validate() == nil
	default:
		return semreg.Decimal{}, "", false
	}
}

func decimalShiftExponent(number semreg.Decimal, delta int32) (semreg.Decimal, bool) {
	if number.Coefficient == "0" {
		return number, true
	}
	return canonicalDecimal(number.Coefficient, number.Exponent10+delta)
}

func decimalShift(number semreg.Decimal, offset int64, offsetExponent int32) (semreg.Decimal, bool) {
	coefficient, ok := new(big.Int).SetString(number.Coefficient, 10)
	if !ok || number.Validate() != nil {
		return semreg.Decimal{}, false
	}
	exponent := number.Exponent10
	if offsetExponent < exponent {
		exponent = offsetExponent
	}
	scale := func(value *big.Int, from int32) *big.Int {
		return new(big.Int).Mul(value, new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(from-exponent)), nil))
	}
	sum := new(big.Int).Add(scale(coefficient, number.Exponent10), scale(big.NewInt(offset), offsetExponent))
	return canonicalDecimal(sum.String(), exponent)
}

func families700720() []string {
	return []string{"controller-700", "controller-720"}
}

func indexedInstance() (byte, byte) { return 0, 10 }

func singleInstance() (byte, byte) { return 0, 0 }

// Gaps lists every row the pack does not express, with the reason or the element it would need.
func Gaps(feeds []Feed) []Feed {
	var gaps []Feed
	for _, feed := range feeds {
		if feed.Gap != "" || feed.Hold != "" {
			gaps = append(gaps, feed)
		}
	}
	return gaps
}

// FindFeed returns the row with the given id.
func FindFeed(feeds []Feed, id string) (Feed, bool) {
	for _, feed := range feeds {
		if feed.ID == id {
			return feed, true
		}
	}
	return Feed{}, false
}
