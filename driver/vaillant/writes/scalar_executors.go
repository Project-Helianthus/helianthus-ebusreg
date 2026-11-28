package writes

import (
	"context"
	"encoding/binary"
	"math"
	"math/big"

	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
	semreg "github.com/Project-Helianthus/helianthus-semreg/semreg/v1"
	semdriver "github.com/Project-Helianthus/helianthus-semreg/semreg/v1/driver"
	"github.com/Project-Helianthus/helianthus-semreg/semreg/v1/operation"
)

// Transport sends one frame. Sent reports whether the frame left the caller.
// The executor does not retry.
type Transport interface {
	Send(ctx context.Context, frame protocol.Frame) (sent bool, response *protocol.Frame, err error)
}

// Owner is the live setting owner rechecked immediately before a send.
type Owner struct {
	Admitted   bool
	Generation semreg.Uint64
	Epoch      semreg.SourceEpochID
	Source     semreg.SourceID
	Origin     byte
	Controller byte
	Instances  map[semreg.NativeBindingID]byte
}

// Scalar is the setting-owner executor for one admitted scalar operation.
// Execute is safe for concurrent use when Transport and the owner callback are.
type Scalar struct {
	transport Transport
	owner     func() Owner
}

// NewScalar returns an executor that sends through transport.
func NewScalar(transport Transport, owner func() Owner) *Scalar {
	return &Scalar{transport: transport, owner: owner}
}

// Execute sends one local write of an admitted scalar and reads that selector
// once. A dispatch that is not admitted, not scalar, or not encodable sends
// nothing.
func (s *Scalar) Execute(ctx context.Context, dispatch semdriver.Dispatch) (semdriver.Completion, error) {
	if s == nil || s.transport == nil || s.owner == nil || ctx == nil || ctx.Err() != nil {
		return noContact(), nil
	}
	row, payload, live, instance, ok := s.prepare(dispatch)
	if !ok {
		return noContact(), nil
	}
	return s.exchange(ctx, dispatch, row, payload, live, instance), nil
}

func (s *Scalar) prepare(dispatch semdriver.Dispatch) (Row, []byte, Owner, byte, bool) {
	row, ok := Operation(string(dispatch.Request.Operation.ID))
	if !ok || row.Form != "scalar" || row.Selector.Family != "local" {
		return Row{}, nil, Owner{}, 0, false
	}
	payload, ok := encodeArgument(row, dispatch.Request.Arguments)
	if !ok {
		return Row{}, nil, Owner{}, 0, false
	}
	live := s.owner()
	instance, ok := admittedInstance(row, dispatch, live)
	if !ok {
		return Row{}, nil, Owner{}, 0, false
	}
	checked := s.owner()
	checkedInstance, ok := admittedInstance(row, dispatch, checked)
	if !ok || checkedInstance != instance || checked.Origin != live.Origin || checked.Controller != live.Controller {
		return Row{}, nil, Owner{}, 0, false
	}
	return row, payload, checked, instance, true
}

func admittedInstance(row Row, dispatch semdriver.Dispatch, live Owner) (byte, bool) {
	if !live.Admitted || live.Origin == 0 || live.Controller == 0 || !sameFence(live, dispatch) {
		return 0, false
	}
	if row.Selector.Instance == "binding" {
		binding, ok := bindingOf(dispatch)
		if !ok {
			return 0, false
		}
		instance, present := live.Instances[binding]
		return instance, present
	}
	return parseByte(row.Selector.Instance)
}

func sameFence(live Owner, dispatch semdriver.Dispatch) bool {
	if live.Generation == "" || dispatch.Request.ExpectedDriverGeneration == "" || dispatch.Guard.DriverGeneration == "" {
		return false
	}
	if live.Generation != dispatch.Request.ExpectedDriverGeneration || live.Generation != dispatch.Guard.DriverGeneration {
		return false
	}
	if live.Epoch == "" || dispatch.Request.ExpectedSourceEpochID == "" || dispatch.Guard.SourceEpochID == "" {
		return false
	}
	if live.Epoch != dispatch.Request.ExpectedSourceEpochID || live.Epoch != dispatch.Guard.SourceEpochID {
		return false
	}
	return live.Source != "" && live.Source == dispatch.Guard.SourceID
}

func bindingOf(dispatch semdriver.Dispatch) (semreg.NativeBindingID, bool) {
	guard := dispatch.Guard.BindingID
	request := dispatch.Request.BindingID
	if guard != "" && request != "" && guard != request {
		return "", false
	}
	if guard != "" {
		return guard, true
	}
	return request, request != ""
}

func (s *Scalar) exchange(ctx context.Context, dispatch semdriver.Dispatch, row Row, payload []byte, live Owner, instance byte) semdriver.Completion {
	numbers, ok := parseSelector(row.Selector)
	if !ok || ctx.Err() != nil {
		return noContact()
	}
	write, ok := localFrame(live, numbers, instance, true, payload)
	if !ok {
		return noContact()
	}
	sent, response, err := s.transport.Send(ctx, write)
	if err != nil && !sent {
		return noContact()
	}
	if !sent {
		return noContact()
	}
	if ctx.Err() != nil {
		return afterWrite(response, operation.ReadbackInconclusive, false)
	}
	read, ok := localFrame(live, numbers, instance, false, nil)
	if !ok {
		return afterWrite(response, operation.ReadbackInconclusive, false)
	}
	readSent, readResponse, err := s.transport.Send(ctx, read)
	if err != nil && !readSent {
		return afterWrite(response, operation.ReadbackInconclusive, false)
	}
	if !readSent || readResponse == nil {
		return afterWrite(response, operation.ReadbackInconclusive, false)
	}
	relation, decisive := compareRead(numbers.group, numbers.register, payload, readResponse.Data)
	completion := afterWrite(response, relation, decisive)
	binding, _ := bindingOf(dispatch)
	completion.Readback = &operation.Readback{
		Relation:         relation,
		BindingID:        binding,
		SourceID:         live.Source,
		SourceEpochID:    live.Epoch,
		DriverGeneration: live.Generation,
	}
	return completion
}

func afterWrite(response *protocol.Frame, relation operation.ReadbackRelation, decisive bool) semdriver.Completion {
	switch {
	case decisive && relation == operation.ReadbackConfirms:
		return semdriver.Completion{Outcome: operation.OutcomeApplied}
	case decisive && relation == operation.ReadbackContradicts:
		return semdriver.Completion{Outcome: operation.OutcomeConflict}
	case response != nil:
		return semdriver.Completion{Outcome: operation.OutcomeAcknowledgedUnverified}
	default:
		return semdriver.Completion{Outcome: operation.OutcomeIndeterminate}
	}
}

func noContact() semdriver.Completion {
	return semdriver.Completion{Outcome: operation.OutcomeFailedNoContact}
}

func localFrame(live Owner, numbers selectorNumbers, instance byte, write bool, payload []byte) (protocol.Frame, bool) {
	access := byte(0x00)
	if write {
		access = 0x01
	}
	data := []byte{numbers.opcode, access, numbers.group, instance, byte(numbers.register), byte(numbers.register >> 8)}
	if write {
		data = append(data, payload...)
	}
	return protocol.Frame{
		Source:    live.Origin,
		Target:    live.Controller,
		Primary:   numbers.primary,
		Secondary: numbers.secondary,
		Data:      data,
		FrameType: protocol.FrameTypeInitiatorTarget,
	}, true
}

func compareRead(group byte, register uint16, payload, reply []byte) (operation.ReadbackRelation, bool) {
	if len(reply) < 4+len(payload) || reply[1] != group || binary.LittleEndian.Uint16(reply[2:4]) != register {
		return operation.ReadbackInconclusive, false
	}
	if bytesEqual(reply[4:4+len(payload)], payload) {
		return operation.ReadbackConfirms, true
	}
	return operation.ReadbackContradicts, true
}

func bytesEqual(left, right []byte) bool {
	if len(left) != len(right) {
		return false
	}
	for i := range left {
		if left[i] != right[i] {
			return false
		}
	}
	return true
}

func encodeArgument(row Row, arguments []semreg.TypedField) ([]byte, bool) {
	if len(arguments) != 1 || arguments[0].ID != semreg.DefinitionID(row.Argument) {
		return nil, false
	}
	field := arguments[0]
	switch row.Kind {
	case "kelvin_celsius":
		return encodeKelvin(field)
	case "native_float32":
		return encodeNativeFloat(field)
	case "percent_uint16":
		return encodePercent(field)
	case "native_uint16":
		return encodeNativeUint(field)
	case "bool_byte":
		return encodeBool(field)
	case "system_switch":
		return encodeSystemSwitch(field)
	case "room_control":
		return encodeRoomControl(field)
	default:
		return nil, false
	}
}

func encodeKelvin(field semreg.TypedField) ([]byte, bool) {
	number, ok := quantity(field, "unit.kelvin")
	if !ok {
		return nil, false
	}
	celsius := new(big.Rat).Sub(number, big.NewRat(27315, 100))
	return putFloat32(celsius)
}

func encodeNativeFloat(field semreg.TypedField) ([]byte, bool) {
	number, ok := quantity(field, "unit.one")
	if !ok {
		return nil, false
	}
	return putFloat32(number)
}

func encodePercent(field semreg.TypedField) ([]byte, bool) {
	number, ok := quantity(field, "unit.ratio")
	if !ok {
		return nil, false
	}
	return putUint16(new(big.Rat).Mul(number, big.NewRat(100, 1)))
}

func encodeNativeUint(field semreg.TypedField) ([]byte, bool) {
	number, ok := quantity(field, "unit.one")
	if !ok {
		return nil, false
	}
	return putUint16(number)
}

func encodeBool(field semreg.TypedField) ([]byte, bool) {
	if field.Value.Kind != semreg.ValueBoolean || field.Value.Boolean == nil {
		return nil, false
	}
	if *field.Value.Boolean {
		return []byte{1}, true
	}
	return []byte{0}, true
}

func encodeSystemSwitch(field semreg.TypedField) ([]byte, bool) {
	token, ok := symbolToken(field, "thermal.mode.system")
	if !ok {
		return nil, false
	}
	switch token {
	case "off":
		return []byte{1}, true
	case "auto":
		return []byte{0}, true
	default:
		return nil, false
	}
}

func encodeRoomControl(field semreg.TypedField) ([]byte, bool) {
	token, ok := symbolToken(field, "ebusreg.circuit_room_temperature_influence_mode")
	if !ok {
		return nil, false
	}
	switch token {
	case "off":
		return []byte{0, 0}, true
	case "modulating":
		return []byte{1, 0}, true
	case "thermostat":
		return []byte{2, 0}, true
	default:
		return nil, false
	}
}

func symbolToken(field semreg.TypedField, namespace semreg.DefinitionID) (string, bool) {
	if field.Value.Kind != semreg.ValueSymbol || field.Value.Symbol == nil || !field.Value.Symbol.Known {
		return "", false
	}
	if field.Value.Symbol.Namespace != namespace {
		return "", false
	}
	return field.Value.Symbol.Token, true
}

func quantity(field semreg.TypedField, unit semreg.DefinitionID) (*big.Rat, bool) {
	if field.Value.Kind != semreg.ValueQuantity || field.Value.Quantity == nil || field.Value.Quantity.Unit != unit {
		return nil, false
	}
	return decimalRat(field.Value.Quantity.Number)
}

func decimalRat(number semreg.Decimal) (*big.Rat, bool) {
	if number.Exponent10 < -18 || number.Exponent10 > 18 || number.Coefficient == "" {
		return nil, false
	}
	coefficient := new(big.Int)
	if _, ok := coefficient.SetString(number.Coefficient, 10); !ok {
		return nil, false
	}
	rat := new(big.Rat).SetInt(coefficient)
	exponent := int64(number.Exponent10)
	if exponent < 0 {
		exponent = -exponent
	}
	scale := new(big.Int).Exp(big.NewInt(10), big.NewInt(exponent), nil)
	if number.Exponent10 >= 0 {
		return rat.Mul(rat, new(big.Rat).SetInt(scale)), true
	}
	return rat.Quo(rat, new(big.Rat).SetInt(scale)), true
}

func putFloat32(number *big.Rat) ([]byte, bool) {
	value, _ := new(big.Float).SetPrec(128).SetRat(number).Float32()
	if math.IsNaN(float64(value)) || math.IsInf(float64(value), 0) {
		return nil, false
	}
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], math.Float32bits(value))
	return buf[:], true
}

func putUint16(number *big.Rat) ([]byte, bool) {
	if number == nil || !number.IsInt() || number.Sign() < 0 {
		return nil, false
	}
	width := number.Num()
	if !width.IsUint64() || width.Uint64() > math.MaxUint16 {
		return nil, false
	}
	var buf [2]byte
	binary.LittleEndian.PutUint16(buf[:], uint16(width.Uint64()))
	return buf[:], true
}
