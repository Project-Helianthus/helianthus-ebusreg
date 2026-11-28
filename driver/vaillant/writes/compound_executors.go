package writes

import (
	"context"
	"encoding/binary"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
	semreg "github.com/Project-Helianthus/helianthus-semreg/semreg/v1"
	semdriver "github.com/Project-Helianthus/helianthus-semreg/semreg/v1/driver"
	"github.com/Project-Helianthus/helianthus-semreg/semreg/v1/operation"
)

const (
	programWeekday semreg.DefinitionID = "ebusreg.time_program.weekday"
	programCount   semreg.DefinitionID = "ebusreg.time_program.slot_count"
)

// Compound is the setting-owner executor for one admitted compound operation.
// Execute is safe for concurrent use when Transport and the owner callback are.
type Compound struct {
	transport Transport
	owner     func() Owner
	limits    *Limits
}

// NewCompound returns an executor that sends each documented frame once.
// limits is consulted when it is non-nil. A nil limits value publishes no queried range.
func NewCompound(transport Transport, owner func() Owner, limits *Limits) *Compound {
	return &Compound{transport: transport, owner: owner, limits: limits}
}

// Execute sends the documented frames of one compound operation and reads them back.
// A dispatch that is not admitted sends nothing. A value outside a published range
// sends nothing.
func (c *Compound) Execute(ctx context.Context, dispatch semdriver.Dispatch) (semdriver.Completion, error) {
	if c == nil || c.transport == nil || c.owner == nil || ctx == nil || ctx.Err() != nil {
		return noContact(), nil
	}
	row, ok := Operation(string(dispatch.Request.Operation.ID))
	if !ok || row.Form != "compound" {
		return noContact(), nil
	}
	live, instance, ok := c.admit(row, dispatch)
	if !ok {
		return noContact(), nil
	}
	switch row.Selector.Family {
	case "timer":
		return c.executeTimer(ctx, dispatch, row, live, instance), nil
	case "local":
		if row.PairRegister != "" {
			return c.executePair(ctx, dispatch, row, live, instance), nil
		}
	}
	return noContact(), nil
}

func (c *Compound) admit(row Row, dispatch semdriver.Dispatch) (Owner, byte, bool) {
	live := c.owner()
	instance, ok := admittedInstance(row, dispatch, live)
	if !ok {
		return Owner{}, 0, false
	}
	checked := c.owner()
	checkedInstance, ok := admittedInstance(row, dispatch, checked)
	if !ok || checkedInstance != instance || checked.Origin != live.Origin || checked.Controller != live.Controller {
		return Owner{}, 0, false
	}
	return checked, instance, true
}

func (c *Compound) still(row Row, dispatch semdriver.Dispatch, live Owner, instance byte) bool {
	checked := c.owner()
	got, ok := admittedInstance(row, dispatch, checked)
	return ok && got == instance && checked.Origin == live.Origin && checked.Controller == live.Controller
}

func (c *Compound) transmit(ctx context.Context, row Row, dispatch semdriver.Dispatch, live Owner, instance byte, frame protocol.Frame) (bool, *protocol.Frame, error, bool) {
	if ctx.Err() != nil || !c.still(row, dispatch, live, instance) {
		return false, nil, nil, false
	}
	sent, response, err := c.transport.Send(ctx, frame)
	return sent, response, err, true
}

type decodeStatus int

const (
	decodeOK decodeStatus = iota
	decodeSilent
	decodeInvalid
)

type slotPlan struct {
	startHour   byte
	startMinute byte
	endHour     byte
	endMinute   byte
	temperature uint16
	explicit    bool
	celsius     *big.Rat
}

type programPlan struct {
	weekday byte
	count   int
	slots   []slotPlan
}

func (c *Compound) executeTimer(ctx context.Context, dispatch semdriver.Dispatch, row Row, live Owner, instance byte) semdriver.Completion {
	plan, status := decodeProgram(row, dispatch.Request.Arguments)
	switch status {
	case decodeSilent:
		return noContact()
	case decodeInvalid:
		return rejected()
	}
	if !c.programAdmitted(row.Operation, instance, plan) {
		return rejected()
	}
	return c.sendTimer(ctx, dispatch, row, live, instance, plan)
}

func (c *Compound) programAdmitted(operation string, instance byte, plan programPlan) bool {
	if c.limits == nil {
		return true
	}
	if ceiling, ok := c.limits.slotCeiling(operation, instance); ok && plan.count > ceiling {
		return false
	}
	if _, ok := c.limits.Published(operation, instance); !ok {
		return true
	}
	for _, slot := range plan.slots {
		if !slot.explicit {
			continue
		}
		quantity, ok := kelvinFromCelsiusRat(slot.celsius)
		if !ok {
			return false
		}
		value := semreg.Value{Kind: semreg.ValueQuantity, Quantity: &quantity}
		if !c.limits.Admit(operation, instance, value) {
			return false
		}
	}
	return true
}

func (c *Compound) sendTimer(ctx context.Context, dispatch semdriver.Dispatch, row Row, live Owner, instance byte, plan programPlan) semdriver.Completion {
	numbers, ok := parseSelector(row.Selector)
	if !ok {
		return noContact()
	}
	total := len(plan.slots)
	steps := make([]semdriver.CompoundStep, total)
	stop := false
	for index := 0; index < total; index++ {
		point := fmt.Sprintf("slot-%d", index)
		steps[index].Point = point
		if stop {
			continue
		}
		frame := timerWriteFrame(live, numbers, instance, plan.weekday, byte(index), byte(plan.count), plan.slots[index])
		sent, response, err, attempted := c.transmit(ctx, row, dispatch, live, instance, frame)
		if !attempted {
			stop = true
			continue
		}
		steps[index].Attempted = true
		steps[index].Delivery = operation.DeliveryNotSent
		if sent {
			steps[index].Delivery = operation.DeliverySent
			steps[index].Acknowledgement = acknowledgement(response)
		}
		if !sent || err != nil || steps[index].Acknowledgement != operation.AckAccepted {
			stop = true
		}
	}
	c.readTimer(ctx, dispatch, row, live, instance, numbers, plan, steps)
	return finishCompound(dispatch, live, steps)
}

func (c *Compound) readTimer(ctx context.Context, dispatch semdriver.Dispatch, row Row, live Owner, instance byte, numbers selectorNumbers, plan programPlan, steps []semdriver.CompoundStep) {
	readOpcode, ok := parseByte(row.Selector.ReadOpcode)
	if !ok {
		return
	}
	for index := range steps {
		if !steps[index].Attempted || steps[index].Delivery != operation.DeliverySent || steps[index].Acknowledgement != operation.AckAccepted {
			continue
		}
		frame := timerReadFrame(live, numbers, readOpcode, instance, plan.weekday, byte(index))
		sent, response, err, attempted := c.transmit(ctx, row, dispatch, live, instance, frame)
		if !attempted || !sent || err != nil || response == nil {
			steps[index].ReadbackRelation = operation.ReadbackInconclusive
			continue
		}
		steps[index].ReadbackRelation = relateTimer(plan.slots[index], response.Data)
	}
}

func timerWriteFrame(live Owner, numbers selectorNumbers, zone, weekday, index, count byte, slot slotPlan) protocol.Frame {
	data := []byte{
		numbers.opcode, zone, numbers.group, weekday, index, count,
		slot.startHour, slot.startMinute, slot.endHour, slot.endMinute,
		byte(slot.temperature), byte(slot.temperature >> 8),
	}
	return addressed(live, numbers, data)
}

func timerReadFrame(live Owner, numbers selectorNumbers, opcode, zone, weekday, index byte) protocol.Frame {
	return addressed(live, numbers, []byte{opcode, zone, numbers.group, weekday, index})
}

func addressed(live Owner, numbers selectorNumbers, data []byte) protocol.Frame {
	return protocol.Frame{
		Source:    live.Origin,
		Target:    live.Controller,
		Primary:   numbers.primary,
		Secondary: numbers.secondary,
		Data:      data,
		FrameType: protocol.FrameTypeInitiatorTarget,
	}
}

func relateTimer(slot slotPlan, reply []byte) operation.ReadbackRelation {
	if len(reply) != 7 || reply[0] != 0 {
		return operation.ReadbackInconclusive
	}
	if reply[1] != slot.startHour || reply[2] != slot.startMinute || reply[3] != slot.endHour || reply[4] != slot.endMinute {
		return operation.ReadbackContradicts
	}
	if !slot.explicit {
		return operation.ReadbackConfirms
	}
	if binary.LittleEndian.Uint16(reply[5:7]) != slot.temperature {
		return operation.ReadbackContradicts
	}
	return operation.ReadbackConfirms
}

func (c *Compound) executePair(ctx context.Context, dispatch semdriver.Dispatch, row Row, live Owner, instance byte) semdriver.Completion {
	text, ok := textArgument(row, dispatch.Request.Arguments)
	if !ok {
		return noContact()
	}
	first, second, reject, ok := encodePairText(text)
	if reject {
		return rejected()
	}
	if !ok {
		return noContact()
	}
	numbers, ok := parseSelector(row.Selector)
	if !ok {
		return noContact()
	}
	pair, ok := pairRegister(row)
	if !ok {
		return noContact()
	}
	steps := []semdriver.CompoundStep{{Point: "first"}, {Point: "second"}}
	payloads := [][]byte{first, second}
	registers := []uint16{numbers.register, pair}
	for index := range steps {
		if index > 0 && steps[0].Delivery != operation.DeliverySent {
			break
		}
		current := numbers
		current.register = registers[index]
		frame, ok := localFrame(live, current, instance, true, payloads[index])
		if !ok {
			break
		}
		sent, response, _, attempted := c.transmit(ctx, row, dispatch, live, instance, frame)
		if !attempted {
			break
		}
		steps[index].Attempted = true
		steps[index].Delivery = operation.DeliveryNotSent
		if sent {
			steps[index].Delivery = operation.DeliverySent
			steps[index].Acknowledgement = registerAcknowledgement(response)
		}
		if !sent {
			break
		}
	}
	if steps[0].Delivery == operation.DeliverySent && steps[1].Delivery == operation.DeliverySent {
		for index := range steps {
			current := numbers
			current.register = registers[index]
			steps[index].ReadbackRelation = c.readRegister(ctx, dispatch, row, live, instance, current, payloads[index])
		}
	}
	return finishCompound(dispatch, live, steps)
}

func (c *Compound) readRegister(ctx context.Context, dispatch semdriver.Dispatch, row Row, live Owner, instance byte, numbers selectorNumbers, payload []byte) operation.ReadbackRelation {
	frame, ok := localFrame(live, numbers, instance, false, nil)
	if !ok {
		return operation.ReadbackInconclusive
	}
	sent, response, err, attempted := c.transmit(ctx, row, dispatch, live, instance, frame)
	if !attempted || !sent || err != nil || response == nil {
		return operation.ReadbackInconclusive
	}
	relation, decisive := compareRead(numbers.group, numbers.register, payload, response.Data)
	if !decisive {
		return operation.ReadbackInconclusive
	}
	return relation
}

func textArgument(row Row, arguments []semreg.TypedField) (string, bool) {
	if len(arguments) != 1 || arguments[0].ID != semreg.DefinitionID(row.Operation) {
		return "", false
	}
	if arguments[0].Value.Kind != semreg.ValueText || arguments[0].Value.Text == nil {
		return "", false
	}
	return *arguments[0].Value.Text, true
}

func encodePairText(raw string) (first, second []byte, reject, ok bool) {
	raw = strings.TrimSpace(raw)
	if strings.IndexByte(raw, 0) >= 0 || !utf8.ValidString(raw) {
		return nil, nil, false, false
	}
	for _, char := range raw {
		if char > 127 {
			return nil, nil, false, false
		}
	}
	if len(raw) > 12 {
		return nil, nil, true, false
	}
	first = make([]byte, 7)
	second = make([]byte, 7)
	if len(raw) > 6 {
		copy(first, raw[:6])
		copy(second, raw[6:])
		return first, second, false, true
	}
	copy(first, raw)
	return first, second, false, true
}

func pairRegister(row Row) (uint16, bool) {
	if row.PairRegister == "" {
		return 0, false
	}
	parsed, err := strconv.ParseUint(row.PairRegister, 16, 16)
	if err != nil {
		return 0, false
	}
	return uint16(parsed), true
}

func decodeProgram(row Row, arguments []semreg.TypedField) (programPlan, decodeStatus) {
	fields, ok := indexFields(arguments)
	if !ok {
		return programPlan{}, decodeSilent
	}
	weekday, status := clockByte(fields, programWeekday, 6)
	if status != decodeOK {
		return programPlan{}, status
	}
	count, status := clockByte(fields, programCount, 255)
	if status != decodeOK {
		return programPlan{}, status
	}
	used := int(count)
	if used == 0 {
		used = 1
	}
	slots := make([]slotPlan, used)
	expected := map[semreg.DefinitionID]struct{}{programWeekday: {}, programCount: {}}
	for index := 0; index < used; index++ {
		slot, ids, status := parseSlot(row, fields, index)
		if status != decodeOK {
			return programPlan{}, status
		}
		slots[index] = slot
		for id := range ids {
			expected[id] = struct{}{}
		}
	}
	if len(fields) != len(expected) {
		return programPlan{}, decodeSilent
	}
	for id := range fields {
		if _, known := expected[id]; !known {
			return programPlan{}, decodeSilent
		}
	}
	return programPlan{weekday: weekday, count: int(count), slots: slots}, decodeOK
}

func parseSlot(row Row, fields map[semreg.DefinitionID]semreg.TypedField, index int) (slotPlan, map[semreg.DefinitionID]struct{}, decodeStatus) {
	names := []string{"start_hour", "start_minute", "end_hour", "end_minute"}
	maxima := []int{24, 59, 24, 59}
	ids := map[semreg.DefinitionID]struct{}{}
	values := make([]byte, len(names))
	for i, name := range names {
		id := programSlotID(index, name)
		ids[id] = struct{}{}
		value, status := clockByte(fields, id, maxima[i])
		if status != decodeOK {
			return slotPlan{}, nil, status
		}
		values[i] = value
	}
	slot := slotPlan{
		startHour: values[0], startMinute: values[1], endHour: values[2], endMinute: values[3],
		temperature: 0xffff,
	}
	tempID := programSlotID(index, "temperature")
	field, present := fields[tempID]
	switch {
	case row.SlotTemperature == "required" && !present:
		return slotPlan{}, nil, decodeSilent
	case present:
		number, raw, status := slotTemperature(field)
		if status != decodeOK {
			return slotPlan{}, nil, status
		}
		ids[tempID] = struct{}{}
		slot.explicit = true
		slot.temperature = raw
		slot.celsius = number
	case row.SlotTemperature != "optional":
		return slotPlan{}, nil, decodeSilent
	}
	return slot, ids, decodeOK
}

func programSlotID(index int, name string) semreg.DefinitionID {
	return semreg.DefinitionID(fmt.Sprintf("ebusreg.time_program.slot.%d.%s", index, name))
}

func indexFields(arguments []semreg.TypedField) (map[semreg.DefinitionID]semreg.TypedField, bool) {
	fields := make(map[semreg.DefinitionID]semreg.TypedField, len(arguments))
	for _, argument := range arguments {
		if _, exists := fields[argument.ID]; exists {
			return nil, false
		}
		fields[argument.ID] = argument
	}
	return fields, true
}

func clockByte(fields map[semreg.DefinitionID]semreg.TypedField, id semreg.DefinitionID, max int) (byte, decodeStatus) {
	field, ok := fields[id]
	if !ok {
		return 0, decodeSilent
	}
	number, ok := quantity(field, "unit.one")
	if !ok || number == nil || !number.IsInt() || !number.Num().IsInt64() {
		return 0, decodeSilent
	}
	value := number.Num().Int64()
	if value < 0 || value > int64(max) {
		return 0, decodeInvalid
	}
	return byte(value), decodeOK
}

func slotTemperature(field semreg.TypedField) (*big.Rat, uint16, decodeStatus) {
	if field.Value.Kind != semreg.ValueQuantity || field.Value.Quantity == nil || field.Value.Quantity.Unit != "unit.celsius" {
		return nil, 0, decodeSilent
	}
	number, ok := decimalRat(field.Value.Quantity.Number)
	if !ok || number == nil {
		return nil, 0, decodeSilent
	}
	// The native word is tenths of a degree. A finer value has no exact encoding.
	tenths := new(big.Rat).Mul(number, big.NewRat(10, 1))
	if !tenths.IsInt() {
		return nil, 0, decodeSilent
	}
	scaled := tenths.Num()
	if scaled.Sign() < 0 || !scaled.IsUint64() || scaled.Uint64() >= 0xffff {
		return nil, 0, decodeSilent
	}
	return number, uint16(scaled.Uint64()), decodeOK
}

func acknowledgement(response *protocol.Frame) operation.AckState {
	if response == nil || len(response.Data) != 1 {
		return ""
	}
	if response.Data[0] == 0 {
		return operation.AckAccepted
	}
	return operation.AckRejected
}

// registerAcknowledgement reads the status of one local register write.
// A one-byte reply is that status. A longer reply carries the status after
// opcode 02 or 06, and otherwise in the first byte. Status 0 accepts the write.
func registerAcknowledgement(response *protocol.Frame) operation.AckState {
	status, ok := registerWriteStatus(response)
	if !ok {
		return ""
	}
	if status == 0 {
		return operation.AckAccepted
	}
	return operation.AckRejected
}

func registerWriteStatus(response *protocol.Frame) (byte, bool) {
	if response == nil || len(response.Data) == 0 {
		return 0, false
	}
	data := response.Data
	if len(data) == 1 {
		return data[0], true
	}
	if data[0] == 0x02 || data[0] == 0x06 {
		return data[1], true
	}
	return data[0], true
}

func rejected() semdriver.Completion {
	return semdriver.Completion{Outcome: operation.OutcomeRejected}
}

func finishCompound(dispatch semdriver.Dispatch, live Owner, steps []semdriver.CompoundStep) semdriver.Completion {
	partial := false
	sent, accepted, rejectedWrites, confirms, contradicts, attempted := 0, 0, 0, 0, 0, 0
	for _, step := range steps {
		if !step.Attempted {
			partial = true
			continue
		}
		attempted++
		if step.Delivery != operation.DeliverySent {
			partial = true
			continue
		}
		sent++
		switch step.Acknowledgement {
		case operation.AckAccepted:
			accepted++
		case operation.AckRejected:
			rejectedWrites++
		}
		switch step.ReadbackRelation {
		case operation.ReadbackConfirms:
			confirms++
		case operation.ReadbackContradicts:
			contradicts++
		}
	}
	outcome := operation.OutcomeIndeterminate
	relation := operation.ReadbackInconclusive
	switch {
	case contradicts > 0:
		outcome = operation.OutcomeConflict
		relation = operation.ReadbackContradicts
	case sent == 0:
		outcome = operation.OutcomeFailedNoContact
		relation = ""
	case !partial && accepted == len(steps) && confirms == len(steps):
		outcome = operation.OutcomeApplied
		relation = operation.ReadbackConfirms
	case accepted == 0 && rejectedWrites == sent && sent == attempted:
		outcome = operation.OutcomeNoEffect
		relation = ""
	case !partial && accepted == sent:
		outcome = operation.OutcomeAcknowledgedUnverified
	}
	completion := semdriver.Completion{
		Outcome:           outcome,
		CompoundExecution: &semdriver.CompoundExecution{Partial: partial, Steps: steps},
	}
	if relation != "" {
		binding, _ := bindingOf(dispatch)
		completion.Readback = &operation.Readback{
			Relation:         relation,
			BindingID:        binding,
			SourceID:         live.Source,
			SourceEpochID:    live.Epoch,
			DriverGeneration: live.Generation,
		}
	}
	return completion
}
