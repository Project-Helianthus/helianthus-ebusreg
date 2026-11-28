package writes

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"math"
	"math/big"
	"sort"
	"sync"

	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
	semreg "github.com/Project-Helianthus/helianthus-semreg/semreg/v1"
	semdriver "github.com/Project-Helianthus/helianthus-semreg/semreg/v1/driver"
)

// Limits publishes operation ranges queried when SemReg is first populated.
// A queried range is kept. A register with no query and no documented numeric
// range publishes no quantity limit.
type Limits struct {
	mu        sync.Mutex
	transport Transport
	owner     func() Owner
	populated bool
	values    map[limitKey]semdriver.ValueLimits
	ceilings  map[limitKey]int
}

type limitKey struct {
	operation string
	instance  byte
}

type queryFate int

const (
	fateContinue queryFate = iota
	fateStop
)

type spanCodec struct {
	width  int
	unit   semreg.DefinitionID
	kelvin bool
	ratio  bool
}

// NewLimits returns the limit table filled by Populate.
func NewLimits(transport Transport, owner func() Owner) *Limits {
	return &Limits{transport: transport, owner: owner}
}

// Populate queries documented ranges once for an admitted owner.
// A second call sends nothing. An owner that is not admitted sends nothing
// and leaves the table unpopulated.
func (l *Limits) Populate(ctx context.Context) error {
	if l == nil || l.transport == nil || l.owner == nil || ctx == nil {
		return nil
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.populated {
		return nil
	}
	if ctx.Err() != nil {
		return ctx.Err()
	}
	live := l.owner()
	if !live.Admitted || live.Origin == 0 || live.Controller == 0 {
		return nil
	}
	rows, err := uniqueOperations()
	if err != nil {
		return nil
	}
	staged := map[limitKey]semdriver.ValueLimits{}
	ceilings := map[limitKey]int{}
	for _, row := range rows {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if !ownerSnapshotEqual(live, l.owner()) {
			return nil
		}
		if l.queryRow(ctx, live, row, staged, ceilings) == fateStop {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			return nil
		}
	}
	l.values = staged
	l.ceilings = ceilings
	l.populated = true
	return nil
}

// Published returns the queried limits for one operation and instance.
func (l *Limits) Published(operation string, instance byte) (semdriver.ValueLimits, bool) {
	if l == nil {
		return semdriver.ValueLimits{}, false
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	limits, ok := l.values[limitKey{operation, instance}]
	if !ok {
		return semdriver.ValueLimits{}, false
	}
	return cloneLimits(limits), true
}

// Admit reports whether value is inside a published range.
// No published range admits the value. A quantity is compared in the published
// unit, inclusive. A non-quantity value is admitted.
func (l *Limits) Admit(operation string, instance byte, value semreg.Value) bool {
	if l == nil {
		return true
	}
	l.mu.Lock()
	limits, ok := l.values[limitKey{operation, instance}]
	l.mu.Unlock()
	if !ok {
		return true
	}
	if value.Kind != semreg.ValueQuantity || value.Quantity == nil {
		return true
	}
	number, ok := decimalRat(value.Quantity.Number)
	if !ok || number == nil {
		return false
	}
	if limits.Min != nil {
		if value.Quantity.Unit != limits.Min.Unit || !atLeast(number, limits.Min.Number) {
			return false
		}
	}
	if limits.Max != nil {
		if value.Quantity.Unit != limits.Max.Unit || !atMost(number, limits.Max.Number) {
			return false
		}
	}
	return true
}

func (l *Limits) slotCeiling(operation string, instance byte) (int, bool) {
	if l == nil {
		return 0, false
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	value, ok := l.ceilings[limitKey{operation, instance}]
	return value, ok
}

func (l *Limits) queryRow(ctx context.Context, live Owner, row Row, staged map[limitKey]semdriver.ValueLimits, ceilings map[limitKey]int) queryFate {
	switch row.Selector.Family {
	case "timer":
		return l.queryTimer(ctx, live, row, staged, ceilings)
	case "local":
		return l.queryLocal(ctx, live, row, staged)
	default:
		return fateContinue
	}
}

func (l *Limits) queryLocal(ctx context.Context, live Owner, row Row, staged map[limitKey]semdriver.ValueLimits) queryFate {
	numbers, ok := parseSelector(row.Selector)
	if !ok {
		return fateContinue
	}
	instances, ok := instancesOf(row, live)
	if !ok {
		return fateContinue
	}
	registers := []uint16{numbers.register}
	if pair, paired := pairRegister(row); paired {
		registers = append(registers, pair)
	}
	codec, hasCodec := codecOf(row)
	for _, instance := range instances {
		for _, register := range registers {
			frame := describeFrame(live, numbers, instance, register)
			body, fate := l.sendQuery(ctx, live, frame)
			if fate == fateStop {
				return fateStop
			}
			if body == nil || !hasCodec || register != numbers.register {
				continue
			}
			storeDescription(staged, row.Operation, instance, numbers.group, register, codec, body)
		}
	}
	return fateContinue
}

func (l *Limits) queryTimer(ctx context.Context, live Owner, row Row, staged map[limitKey]semdriver.ValueLimits, ceilings map[limitKey]int) queryFate {
	numbers, ok := parseSelector(row.Selector)
	if !ok {
		return fateContinue
	}
	instances, ok := instancesOf(row, live)
	if !ok {
		return fateContinue
	}
	for _, instance := range instances {
		frame := configFrame(live, numbers, instance)
		body, fate := l.sendQuery(ctx, live, frame)
		if fate == fateStop {
			return fateStop
		}
		if body == nil {
			continue
		}
		storeTimer(staged, ceilings, row.Operation, instance, body)
	}
	return fateContinue
}

func (l *Limits) sendQuery(ctx context.Context, live Owner, frame protocol.Frame) ([]byte, queryFate) {
	if ctx.Err() != nil || !ownerSnapshotEqual(live, l.owner()) {
		return nil, fateStop
	}
	sent, response, err := l.transport.Send(ctx, frame)
	if ctx.Err() != nil {
		return nil, fateStop
	}
	if err != nil || !sent || response == nil {
		return nil, fateContinue
	}
	return append([]byte(nil), response.Data...), fateContinue
}

func describeFrame(live Owner, numbers selectorNumbers, instance byte, register uint16) protocol.Frame {
	return addressed(live, numbers, []byte{0x01, numbers.group, instance, byte(register), byte(register >> 8)})
}

func configFrame(live Owner, numbers selectorNumbers, zone byte) protocol.Frame {
	return addressed(live, numbers, []byte{0xa3, zone, numbers.group})
}

func codecOf(row Row) (spanCodec, bool) {
	if row.Selector.Family == "local" && row.Selector.Group == "01" && row.Selector.Register == "0004" {
		return spanCodec{width: 4, unit: "unit.kelvin", kelvin: true}, true
	}
	switch row.Kind {
	case "kelvin_celsius":
		return spanCodec{width: 4, unit: "unit.kelvin", kelvin: true}, true
	case "native_float32":
		return spanCodec{width: 4, unit: "unit.one"}, true
	case "percent_uint16":
		return spanCodec{width: 2, unit: "unit.ratio", ratio: true}, true
	case "native_uint16", "room_control":
		return spanCodec{width: 2, unit: "unit.one"}, true
	case "bool_byte", "system_switch":
		return spanCodec{width: 1, unit: "unit.one"}, true
	default:
		return spanCodec{}, false
	}
}

func storeDescription(staged map[limitKey]semdriver.ValueLimits, operation string, instance byte, group byte, register uint16, codec spanCodec, reply []byte) {
	minimum, maximum, step, ok := decodeSpans(reply, group, register, codec.width)
	if !ok {
		return
	}
	minQ, ok := absoluteQuantity(minimum, codec)
	if !ok {
		return
	}
	maxQ, ok := absoluteQuantity(maximum, codec)
	if !ok {
		return
	}
	limits := semdriver.ValueLimits{
		Min: &minQ, Max: &maxQ, Origin: semdriver.LimitDeviceQuery, Evidence: limitEvidence(reply),
	}
	if step > 0 {
		stepQ, ok := deltaQuantity(step, codec)
		if !ok {
			return
		}
		limits.Step = &stepQ
	}
	if limits.Validate() != nil {
		return
	}
	staged[limitKey{operation, instance}] = limits
}

func storeTimer(staged map[limitKey]semdriver.ValueLimits, ceilings map[limitKey]int, operation string, instance byte, reply []byte) {
	if len(reply) != 9 || reply[0] != 0 {
		return
	}
	key := limitKey{operation, instance}
	ceilings[key] = int(reply[1])
	var minimum, maximum *semreg.Quantity
	if reply[6] != 0xff {
		quantity, ok := wholeCelsiusKelvin(reply[6])
		if !ok {
			return
		}
		minimum = &quantity
	}
	if reply[7] != 0xff {
		quantity, ok := wholeCelsiusKelvin(reply[7])
		if !ok {
			return
		}
		maximum = &quantity
	}
	if minimum != nil && maximum != nil && !orderedQuantities(minimum, maximum) {
		return
	}
	if minimum == nil && maximum == nil {
		return
	}
	limits := semdriver.ValueLimits{
		Min: minimum, Max: maximum, Origin: semdriver.LimitDeviceQuery, Evidence: limitEvidence(reply),
	}
	if limits.Validate() != nil {
		return
	}
	staged[key] = limits
}

func decodeSpans(reply []byte, group byte, register uint16, width int) (min, max, step float64, ok bool) {
	if width != 1 && width != 2 && width != 4 {
		return 0, 0, 0, false
	}
	if len(reply) != 3+3*width || reply[0] != group || binary.LittleEndian.Uint16(reply[1:3]) != register {
		return 0, 0, 0, false
	}
	values := make([]float64, 3)
	for i := 0; i < 3; i++ {
		raw := reply[3+i*width : 3+(i+1)*width]
		switch width {
		case 4:
			value := math.Float32frombits(binary.LittleEndian.Uint32(raw))
			if math.IsNaN(float64(value)) || math.IsInf(float64(value), 0) {
				return 0, 0, 0, false
			}
			values[i] = float64(value)
		case 2:
			values[i] = float64(binary.LittleEndian.Uint16(raw))
		default:
			values[i] = float64(raw[0])
		}
	}
	if values[0] > values[1] || values[2] < 0 {
		return 0, 0, 0, false
	}
	return values[0], values[1], values[2], true
}

func absoluteQuantity(value float64, codec spanCodec) (semreg.Quantity, bool) {
	rat := float64ToRat(value)
	if codec.ratio {
		rat = new(big.Rat).Quo(rat, big.NewRat(100, 1))
	}
	if codec.kelvin {
		rat = new(big.Rat).Add(rat, big.NewRat(27315, 100))
	}
	return quantityFromRat(rat, codec.unit)
}

func deltaQuantity(value float64, codec spanCodec) (semreg.Quantity, bool) {
	rat := float64ToRat(value)
	if codec.ratio {
		rat = new(big.Rat).Quo(rat, big.NewRat(100, 1))
	}
	return quantityFromRat(rat, codec.unit)
}

func wholeCelsiusKelvin(value byte) (semreg.Quantity, bool) {
	rat := new(big.Rat).SetInt64(int64(value))
	rat.Add(rat, big.NewRat(27315, 100))
	return quantityFromRat(rat, "unit.kelvin")
}

func kelvinFromCelsiusRat(number *big.Rat) (semreg.Quantity, bool) {
	if number == nil {
		return semreg.Quantity{}, false
	}
	sum := new(big.Rat).Add(new(big.Rat).Set(number), big.NewRat(27315, 100))
	return quantityFromRat(sum, "unit.kelvin")
}

func quantityFromRat(number *big.Rat, unit semreg.DefinitionID) (semreg.Quantity, bool) {
	decimal, ok := ratToDecimal(number)
	if !ok {
		return semreg.Quantity{}, false
	}
	return semreg.Quantity{Number: decimal, Unit: unit}, true
}

func ratToDecimal(number *big.Rat) (semreg.Decimal, bool) {
	if number == nil {
		return semreg.Decimal{}, false
	}
	if number.Sign() == 0 {
		return semreg.Decimal{Coefficient: "0", Exponent10: 0}, true
	}
	coefficient := new(big.Int).Abs(number.Num())
	denominator := new(big.Int).Set(number.Denom())
	twos, fives := 0, 0
	for denominator.Bit(0) == 0 {
		denominator.Rsh(denominator, 1)
		twos++
	}
	five := big.NewInt(5)
	for new(big.Int).Mod(denominator, five).Sign() == 0 {
		denominator.Quo(denominator, five)
		fives++
	}
	if denominator.Cmp(big.NewInt(1)) != 0 {
		return semreg.Decimal{}, false
	}
	exp := twos
	if fives > twos {
		coefficient.Mul(coefficient, new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(fives-twos)), nil))
		exp = fives
	} else if twos > fives {
		coefficient.Mul(coefficient, new(big.Int).Exp(big.NewInt(5), big.NewInt(int64(twos-fives)), nil))
	}
	if exp > 18 {
		return semreg.Decimal{}, false
	}
	text := coefficient.String()
	for exp > 0 && len(text) > 1 && text[len(text)-1] == '0' {
		text = text[:len(text)-1]
		exp--
	}
	if number.Sign() < 0 {
		text = "-" + text
	}
	return semreg.Decimal{Coefficient: text, Exponent10: int32(-exp)}, true
}

func float64ToRat(value float64) *big.Rat {
	rat, _ := new(big.Float).SetPrec(64).SetMode(big.ToNearestEven).SetFloat64(value).Rat(nil)
	if rat == nil {
		return new(big.Rat)
	}
	return rat
}

func limitEvidence(raw []byte) []semreg.EvidenceRef {
	sum := sha256.Sum256(raw)
	return []semreg.EvidenceRef{{
		Owner:     "ebusreg.writes",
		Kind:      "evidence.device-query",
		Digest:    semreg.Digest("sha256:" + hex.EncodeToString(sum[:])),
		Contract:  "ebusreg.limit/v1",
		Access:    semreg.EvidenceAccessPublic,
		Redaction: semreg.RedactionNone,
	}}
}

func orderedQuantities(minimum, maximum *semreg.Quantity) bool {
	left, ok := decimalRat(minimum.Number)
	if !ok || left == nil {
		return false
	}
	right, ok := decimalRat(maximum.Number)
	return ok && right != nil && left.Cmp(right) <= 0
}

func atLeast(number *big.Rat, bound semreg.Decimal) bool {
	floor, ok := decimalRat(bound)
	return ok && floor != nil && number.Cmp(floor) >= 0
}

func atMost(number *big.Rat, bound semreg.Decimal) bool {
	ceiling, ok := decimalRat(bound)
	return ok && ceiling != nil && number.Cmp(ceiling) <= 0
}

func uniqueOperations() ([]Row, error) {
	rows, err := Inventory()
	if err != nil {
		return nil, err
	}
	seen := map[string]struct{}{}
	out := make([]Row, 0, len(rows))
	for _, row := range rows {
		if _, ok := seen[row.Operation]; ok {
			continue
		}
		seen[row.Operation] = struct{}{}
		out = append(out, row)
	}
	return out, nil
}

func instancesOf(row Row, live Owner) ([]byte, bool) {
	switch row.Selector.Instance {
	case "unknown", "":
		return nil, false
	case "binding":
		return bindingInstances(live), true
	default:
		instance, ok := parseByte(row.Selector.Instance)
		if !ok {
			return nil, false
		}
		return []byte{instance}, true
	}
}

func bindingInstances(live Owner) []byte {
	ids := make([]string, 0, len(live.Instances))
	for id := range live.Instances {
		ids = append(ids, string(id))
	}
	sort.Strings(ids)
	seen := map[byte]struct{}{}
	out := make([]byte, 0, len(ids))
	for _, id := range ids {
		instance := live.Instances[semreg.NativeBindingID(id)]
		if _, ok := seen[instance]; ok {
			continue
		}
		seen[instance] = struct{}{}
		out = append(out, instance)
	}
	return out
}

func ownerSnapshotEqual(left, right Owner) bool {
	if left.Admitted != right.Admitted || left.Origin != right.Origin || left.Controller != right.Controller {
		return false
	}
	if left.Generation != right.Generation || left.Epoch != right.Epoch || left.Source != right.Source {
		return false
	}
	if len(left.Instances) != len(right.Instances) {
		return false
	}
	for id, instance := range left.Instances {
		if right.Instances[id] != instance {
			return false
		}
	}
	return true
}

func cloneLimits(in semdriver.ValueLimits) semdriver.ValueLimits {
	out := in
	out.Evidence = append([]semreg.EvidenceRef(nil), in.Evidence...)
	if in.Min != nil {
		value := *in.Min
		out.Min = &value
	}
	if in.Max != nil {
		value := *in.Max
		out.Max = &value
	}
	if in.Step != nil {
		value := *in.Step
		out.Step = &value
	}
	return out
}
