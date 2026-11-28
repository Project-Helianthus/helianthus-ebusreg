package writes

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"strconv"
	"sync"
	"testing"

	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
	semreg "github.com/Project-Helianthus/helianthus-semreg/semreg/v1"
	semdriver "github.com/Project-Helianthus/helianthus-semreg/semreg/v1/driver"
	"github.com/Project-Helianthus/helianthus-semreg/semreg/v1/operation"
)

func TestTimeProgramExecutorSlots(t *testing.T) {
	assertTwoSlotProgram(t, false)
	assertTwoSlotProgram(t, true)
	assertOmittedTemperature(t)
	assertHeatingSlot(t)
	assertTimerFence(t)
	assertTimerDrop(t)
	assertTimerRejected(t, 3, 25, 0, 1, 0)
	assertTimerRejected(t, 3, 6, 60, 12, 20)
	assertTimerRejected(t, 7, 6, 10, 12, 20)
	assertTimerSilent(t, false, true)
	assertTimerSilent(t, true, false)
}

func TestCStringPairExecutor(t *testing.T) {
	assertNamePair(t, false)
	assertNamePair(t, true)
	assertStatusFirstName(t)
	assertPhonePair(t)
	assertExactSix(t)
	assertTwelve(t)
	assertRejectedText(t, "ebusreg.operation.installer_name", "ABCDEFGHIJKLM")
	assertSilentText(t, "ebusreg.operation.installer_name", "Hélian")
	assertSilentText(t, "ebusreg.operation.installer_phone_number", "55\x000188")
	assertFirstRegisterDropped(t)
	assertNonzeroStatusWritesSecond(t)
	assertLongStatusWritesSecond(t)
}

func assertTwoSlotProgram(t *testing.T, contradict bool) {
	t.Helper()
	bus := &scriptReply{respond: dhwTimerReply(contradict)}
	owner := compoundOwner(map[semreg.NativeBindingID]byte{"hot-water": 0x01})
	executor := NewCompound(bus, ownerSource(&owner), nil)
	dispatch := programDispatch("ebusreg.operation.dhw_time_program", "hot-water", 3, 2, []slotSpec{
		{6, 10, 12, 20, decimalPtr("40", 0)},
		{12, 20, 23, 50, decimalPtr("555", -1)},
	})
	completion, err := executor.Execute(context.Background(), dispatch)
	if err != nil {
		t.Fatal(err)
	}
	frames := bus.snapshot()
	if len(frames) != 4 {
		t.Fatalf("frames = %d", len(frames))
	}
	assertTimerFrame(t, frames[0], []byte{0xa6, 0xff, 0x02, 0x03, 0x00, 0x02, 0x06, 0x0a, 0x0c, 0x14, 0x90, 0x01})
	assertTimerFrame(t, frames[1], []byte{0xa6, 0xff, 0x02, 0x03, 0x01, 0x02, 0x0c, 0x14, 0x17, 0x32, 0x2b, 0x02})
	assertTimerFrame(t, frames[2], []byte{0xa5, 0xff, 0x02, 0x03, 0x00})
	assertTimerFrame(t, frames[3], []byte{0xa5, 0xff, 0x02, 0x03, 0x01})
	want := operation.OutcomeApplied
	relation := operation.ReadbackConfirms
	if contradict {
		want = operation.OutcomeConflict
		relation = operation.ReadbackContradicts
	}
	assertCompoundOutcome(t, completion, want, relation, false, "hot-water")
	if completion.CompoundExecution == nil || len(completion.CompoundExecution.Steps) != 2 {
		t.Fatalf("steps = %+v", completion.CompoundExecution)
	}
	for index, step := range completion.CompoundExecution.Steps {
		if step.Point != "slot-"+strconv.Itoa(index) || !step.Attempted || step.Delivery != operation.DeliverySent || step.Acknowledgement != operation.AckAccepted {
			t.Fatalf("step %d = %+v", index, step)
		}
	}
	if err := completion.CompoundExecution.Validate(); err != nil {
		t.Fatal(err)
	}
}

func assertOmittedTemperature(t *testing.T) {
	t.Helper()
	bus := &scriptReply{respond: dhwTimerReply(false)}
	owner := compoundOwner(map[semreg.NativeBindingID]byte{"hot-water": 0x01})
	executor := NewCompound(bus, ownerSource(&owner), nil)
	dispatch := programDispatch("ebusreg.operation.dhw_time_program", "hot-water", 3, 2, []slotSpec{
		{6, 10, 12, 20, nil},
		{12, 20, 23, 50, nil},
	})
	completion, err := executor.Execute(context.Background(), dispatch)
	if err != nil {
		t.Fatal(err)
	}
	frames := bus.snapshot()
	if len(frames) != 4 {
		t.Fatalf("frames = %d", len(frames))
	}
	assertTimerFrame(t, frames[0], []byte{0xa6, 0xff, 0x02, 0x03, 0x00, 0x02, 0x06, 0x0a, 0x0c, 0x14, 0xff, 0xff})
	assertTimerFrame(t, frames[1], []byte{0xa6, 0xff, 0x02, 0x03, 0x01, 0x02, 0x0c, 0x14, 0x17, 0x32, 0xff, 0xff})
	assertCompoundOutcome(t, completion, operation.OutcomeApplied, operation.ReadbackConfirms, false, "hot-water")
}

func assertHeatingSlot(t *testing.T) {
	t.Helper()
	bus := &scriptReply{respond: heatingTimerReply()}
	owner := compoundOwner(map[semreg.NativeBindingID]byte{"zone-h": 0x01})
	executor := NewCompound(bus, ownerSource(&owner), nil)
	dispatch := programDispatch("ebusreg.operation.zone_heating_time_program", "zone-h", 4, 1, []slotSpec{
		{8, 15, 9, 45, decimalPtr("215", -1)},
	})
	completion, err := executor.Execute(context.Background(), dispatch)
	if err != nil {
		t.Fatal(err)
	}
	frames := bus.snapshot()
	if len(frames) != 2 {
		t.Fatalf("frames = %d", len(frames))
	}
	assertTimerFrame(t, frames[0], []byte{0xa6, 0x01, 0x00, 0x04, 0x00, 0x01, 0x08, 0x0f, 0x09, 0x2d, 0xd7, 0x00})
	assertTimerFrame(t, frames[1], []byte{0xa5, 0x01, 0x00, 0x04, 0x00})
	assertCompoundOutcome(t, completion, operation.OutcomeApplied, operation.ReadbackConfirms, false, "zone-h")
}

func assertTimerFence(t *testing.T) {
	t.Helper()
	owner := compoundOwner(map[semreg.NativeBindingID]byte{"hot-water": 0x01})
	bus := &scriptReply{respond: dhwTimerReply(false)}
	bus.onSend = func() { owner.Generation = "9" }
	executor := NewCompound(bus, ownerSource(&owner), nil)
	dispatch := programDispatch("ebusreg.operation.dhw_time_program", "hot-water", 3, 2, []slotSpec{
		{6, 10, 12, 20, decimalPtr("40", 0)},
		{12, 20, 23, 50, decimalPtr("555", -1)},
	})
	completion, err := executor.Execute(context.Background(), dispatch)
	if err != nil {
		t.Fatal(err)
	}
	frames := bus.snapshot()
	if len(frames) != 1 || frames[0].Data[0] != 0xa6 || frames[0].Data[4] != 0x00 {
		t.Fatalf("frames = %x", frameDatas(frames))
	}
	if completion.Outcome == operation.OutcomeApplied || completion.CompoundExecution == nil || !completion.CompoundExecution.Partial {
		t.Fatalf("completion = %+v", completion)
	}
	if completion.CompoundExecution.Steps[1].Attempted {
		t.Fatalf("second step = %+v", completion.CompoundExecution.Steps[1])
	}
}

func assertTimerDrop(t *testing.T) {
	t.Helper()
	bus := &scriptReply{respond: func(frame protocol.Frame) (bool, *protocol.Frame, error) {
		return false, nil, errors.New("drop")
	}}
	owner := compoundOwner(map[semreg.NativeBindingID]byte{"hot-water": 0x01})
	executor := NewCompound(bus, ownerSource(&owner), nil)
	dispatch := programDispatch("ebusreg.operation.dhw_time_program", "hot-water", 3, 2, []slotSpec{
		{6, 10, 12, 20, decimalPtr("40", 0)},
		{12, 20, 23, 50, decimalPtr("555", -1)},
	})
	completion, err := executor.Execute(context.Background(), dispatch)
	if err != nil {
		t.Fatal(err)
	}
	if len(bus.snapshot()) != 1 || completion.Outcome != operation.OutcomeFailedNoContact || completion.CompoundExecution == nil {
		t.Fatalf("completion=%+v frames=%d", completion, len(bus.snapshot()))
	}
	first := completion.CompoundExecution.Steps[0]
	second := completion.CompoundExecution.Steps[1]
	if !first.Attempted || first.Delivery != operation.DeliveryNotSent || first.Reason != "" || second.Attempted || second.Delivery != "" {
		t.Fatalf("steps = %+v", completion.CompoundExecution.Steps)
	}
}

func assertTimerRejected(t *testing.T, weekday, startHour, startMinute, endHour, endMinute int) {
	t.Helper()
	bus := &scriptReply{respond: dhwTimerReply(false)}
	owner := compoundOwner(map[semreg.NativeBindingID]byte{"hot-water": 0x01})
	executor := NewCompound(bus, ownerSource(&owner), nil)
	dispatch := programDispatch("ebusreg.operation.dhw_time_program", "hot-water", weekday, 1, []slotSpec{
		{startHour, startMinute, endHour, endMinute, decimalPtr("40", 0)},
	})
	completion, err := executor.Execute(context.Background(), dispatch)
	if err != nil || completion.Outcome != operation.OutcomeRejected || completion.CompoundExecution != nil || len(bus.snapshot()) != 0 {
		t.Fatalf("completion=%+v frames=%d err=%v", completion, len(bus.snapshot()), err)
	}
}

func assertTimerSilent(t *testing.T, cancel, block bool) {
	t.Helper()
	bus := &scriptReply{respond: dhwTimerReply(false)}
	owner := compoundOwner(map[semreg.NativeBindingID]byte{"hot-water": 0x01})
	if block {
		owner.Admitted = false
	}
	executor := NewCompound(bus, ownerSource(&owner), nil)
	ctx := context.Background()
	if cancel {
		var stop context.CancelFunc
		ctx, stop = context.WithCancel(ctx)
		stop()
	}
	dispatch := programDispatch("ebusreg.operation.dhw_time_program", "hot-water", 3, 1, []slotSpec{
		{6, 10, 12, 20, decimalPtr("40", 0)},
	})
	completion, err := executor.Execute(ctx, dispatch)
	if err != nil || completion.Outcome != operation.OutcomeFailedNoContact || completion.CompoundExecution != nil || len(bus.snapshot()) != 0 {
		t.Fatalf("completion=%+v frames=%d err=%v", completion, len(bus.snapshot()), err)
	}
}

func assertNamePair(t *testing.T, tweak bool) {
	t.Helper()
	payloads := map[uint16][]byte{
		0x006c: {'H', 'e', 'l', 'i', 'a', 'n', 0},
		0x006d: {'t', 'h', 'u', 's', 0, 0, 0},
	}
	tweakRegister := uint16(0)
	if tweak {
		tweakRegister = 0x006d
	}
	bus := &scriptReply{respond: registerReply(payloads, tweakRegister, 0)}
	owner := compoundOwner(nil)
	executor := NewCompound(bus, ownerSource(&owner), nil)
	completion, err := executor.Execute(context.Background(), textDispatch("ebusreg.operation.installer_name", "  Helianthus  "))
	if err != nil {
		t.Fatal(err)
	}
	frames := bus.snapshot()
	if len(frames) != 4 {
		t.Fatalf("frames = %d", len(frames))
	}
	assertLocalFrame(t, frames[0], 0x31, 0x15, 0x01, 0x00, 0x00, 0x006c, payloads[0x006c])
	assertLocalFrame(t, frames[1], 0x31, 0x15, 0x01, 0x00, 0x00, 0x006d, payloads[0x006d])
	assertLocalFrame(t, frames[2], 0x31, 0x15, 0x00, 0x00, 0x00, 0x006c, nil)
	assertLocalFrame(t, frames[3], 0x31, 0x15, 0x00, 0x00, 0x00, 0x006d, nil)
	want := operation.OutcomeApplied
	relation := operation.ReadbackConfirms
	if tweak {
		want = operation.OutcomeConflict
		relation = operation.ReadbackContradicts
	}
	assertCompoundOutcome(t, completion, want, relation, false, "system")
	if len(completion.CompoundExecution.Steps) != 2 || completion.CompoundExecution.Steps[0].Point != "first" || completion.CompoundExecution.Steps[1].Point != "second" {
		t.Fatalf("steps = %+v", completion.CompoundExecution.Steps)
	}
	for _, step := range completion.CompoundExecution.Steps {
		if step.Acknowledgement != operation.AckAccepted {
			t.Fatalf("acknowledgement = %+v", step)
		}
	}
}

func assertStatusFirstName(t *testing.T) {
	t.Helper()
	payloads := map[uint16][]byte{
		0x006c: {'H', 'e', 'l', 'i', 'a', 'n', 0},
		0x006d: {'t', 'h', 'u', 's', 0, 0, 0},
	}
	bus := &scriptReply{respond: registerReplyShape(payloads, 0, 0, false)}
	owner := compoundOwner(nil)
	executor := NewCompound(bus, ownerSource(&owner), nil)
	completion, err := executor.Execute(context.Background(), textDispatch("ebusreg.operation.installer_name", "Helianthus"))
	if err != nil {
		t.Fatal(err)
	}
	if len(bus.snapshot()) != 4 {
		t.Fatalf("frames = %d", len(bus.snapshot()))
	}
	assertCompoundOutcome(t, completion, operation.OutcomeApplied, operation.ReadbackConfirms, false, "system")
	for _, step := range completion.CompoundExecution.Steps {
		if step.Acknowledgement != operation.AckAccepted {
			t.Fatalf("acknowledgement = %+v", step)
		}
	}
}

func assertPhonePair(t *testing.T) {
	t.Helper()
	payloads := map[uint16][]byte{
		0x006f: {'+', '1', '5', '5', '5', '1', 0},
		0x0070: {'2', 0, 0, 0, 0, 0, 0},
	}
	bus := &scriptReply{respond: registerReply(payloads, 0, 0)}
	owner := compoundOwner(nil)
	executor := NewCompound(bus, ownerSource(&owner), nil)
	completion, err := executor.Execute(context.Background(), textDispatch("ebusreg.operation.installer_phone_number", "+155512"))
	if err != nil {
		t.Fatal(err)
	}
	frames := bus.snapshot()
	if len(frames) != 4 {
		t.Fatalf("frames = %d", len(frames))
	}
	assertLocalFrame(t, frames[0], 0x31, 0x15, 0x01, 0x00, 0x00, 0x006f, payloads[0x006f])
	assertLocalFrame(t, frames[1], 0x31, 0x15, 0x01, 0x00, 0x00, 0x0070, payloads[0x0070])
	assertLocalFrame(t, frames[2], 0x31, 0x15, 0x00, 0x00, 0x00, 0x006f, nil)
	assertLocalFrame(t, frames[3], 0x31, 0x15, 0x00, 0x00, 0x00, 0x0070, nil)
	assertCompoundOutcome(t, completion, operation.OutcomeApplied, operation.ReadbackConfirms, false, "system")
}

func assertExactSix(t *testing.T) {
	t.Helper()
	first := []byte{'H', 'e', 'l', 'i', 'a', 'n', 0}
	second := []byte{0, 0, 0, 0, 0, 0, 0}
	bus := &scriptReply{respond: registerReply(map[uint16][]byte{0x006c: first, 0x006d: second}, 0, 0)}
	owner := compoundOwner(nil)
	executor := NewCompound(bus, ownerSource(&owner), nil)
	completion, err := executor.Execute(context.Background(), textDispatch("ebusreg.operation.installer_name", "Helian"))
	if err != nil {
		t.Fatal(err)
	}
	frames := bus.snapshot()
	if len(frames) != 4 {
		t.Fatalf("frames = %d", len(frames))
	}
	assertLocalFrame(t, frames[1], 0x31, 0x15, 0x01, 0x00, 0x00, 0x006d, second)
	assertCompoundOutcome(t, completion, operation.OutcomeApplied, operation.ReadbackConfirms, false, "system")
}

func assertTwelve(t *testing.T) {
	t.Helper()
	first := []byte{'A', 'B', 'C', 'D', 'E', 'F', 0}
	second := []byte{'G', 'H', 'I', 'J', 'K', 'L', 0}
	bus := &scriptReply{respond: registerReply(map[uint16][]byte{0x006c: first, 0x006d: second}, 0, 0)}
	owner := compoundOwner(nil)
	executor := NewCompound(bus, ownerSource(&owner), nil)
	_, err := executor.Execute(context.Background(), textDispatch("ebusreg.operation.installer_name", "ABCDEFGHIJKL"))
	if err != nil {
		t.Fatal(err)
	}
	frames := bus.snapshot()
	if len(frames) != 4 {
		t.Fatalf("frames = %d", len(frames))
	}
	assertLocalFrame(t, frames[0], 0x31, 0x15, 0x01, 0x00, 0x00, 0x006c, first)
	assertLocalFrame(t, frames[1], 0x31, 0x15, 0x01, 0x00, 0x00, 0x006d, second)
}

func assertRejectedText(t *testing.T, id, text string) {
	t.Helper()
	bus := &scriptReply{}
	owner := compoundOwner(nil)
	executor := NewCompound(bus, ownerSource(&owner), nil)
	completion, err := executor.Execute(context.Background(), textDispatch(id, text))
	if err != nil || completion.Outcome != operation.OutcomeRejected || completion.CompoundExecution != nil || len(bus.snapshot()) != 0 {
		t.Fatalf("completion=%+v frames=%d err=%v", completion, len(bus.snapshot()), err)
	}
}

func assertSilentText(t *testing.T, id, text string) {
	t.Helper()
	bus := &scriptReply{}
	owner := compoundOwner(nil)
	executor := NewCompound(bus, ownerSource(&owner), nil)
	completion, err := executor.Execute(context.Background(), textDispatch(id, text))
	if err != nil || completion.Outcome != operation.OutcomeFailedNoContact || len(bus.snapshot()) != 0 {
		t.Fatalf("completion=%+v frames=%d err=%v", completion, len(bus.snapshot()), err)
	}
}

func assertFirstRegisterDropped(t *testing.T) {
	t.Helper()
	bus := &scriptReply{respond: func(protocol.Frame) (bool, *protocol.Frame, error) {
		return false, nil, errors.New("drop")
	}}
	owner := compoundOwner(nil)
	executor := NewCompound(bus, ownerSource(&owner), nil)
	completion, err := executor.Execute(context.Background(), textDispatch("ebusreg.operation.installer_name", "Helianthus"))
	if err != nil || len(bus.snapshot()) != 1 || completion.CompoundExecution == nil {
		t.Fatalf("completion=%+v frames=%d err=%v", completion, len(bus.snapshot()), err)
	}
	if completion.CompoundExecution.Steps[1].Attempted {
		t.Fatalf("second step = %+v", completion.CompoundExecution.Steps[1])
	}
}

func assertNonzeroStatusWritesSecond(t *testing.T) {
	t.Helper()
	writes := 0
	bus := &scriptReply{respond: func(frame protocol.Frame) (bool, *protocol.Frame, error) {
		if len(frame.Data) >= 2 && frame.Data[1] == 0x01 {
			writes++
			if writes == 1 {
				return true, &protocol.Frame{Data: []byte{0x02}}, nil
			}
			return true, &protocol.Frame{Data: []byte{0x00}}, nil
		}
		return registerReply(map[uint16][]byte{
			0x006c: {'H', 'e', 'l', 'i', 'a', 'n', 0},
			0x006d: {'t', 'h', 'u', 's', 0, 0, 0},
		}, 0, 0)(frame)
	}}
	owner := compoundOwner(nil)
	executor := NewCompound(bus, ownerSource(&owner), nil)
	_, err := executor.Execute(context.Background(), textDispatch("ebusreg.operation.installer_name", "Helianthus"))
	if err != nil {
		t.Fatal(err)
	}
	frames := bus.snapshot()
	if len(frames) != 4 || writes != 2 {
		t.Fatalf("frames=%d writes=%d", len(frames), writes)
	}
	assertLocalFrame(t, frames[1], 0x31, 0x15, 0x01, 0x00, 0x00, 0x006d, []byte{'t', 'h', 'u', 's', 0, 0, 0})
}

func assertLongStatusWritesSecond(t *testing.T) {
	t.Helper()
	payloads := map[uint16][]byte{
		0x006c: {'H', 'e', 'l', 'i', 'a', 'n', 0},
		0x006d: {'t', 'h', 'u', 's', 0, 0, 0},
	}
	bus := &scriptReply{respond: registerReply(payloads, 0, 1)}
	owner := compoundOwner(nil)
	executor := NewCompound(bus, ownerSource(&owner), nil)
	completion, err := executor.Execute(context.Background(), textDispatch("ebusreg.operation.installer_name", "Helianthus"))
	if err != nil {
		t.Fatal(err)
	}
	frames := bus.snapshot()
	if len(frames) != 4 {
		t.Fatalf("frames = %d", len(frames))
	}
	assertLocalFrame(t, frames[1], 0x31, 0x15, 0x01, 0x00, 0x00, 0x006d, payloads[0x006d])
	if completion.Outcome == operation.OutcomeApplied || completion.CompoundExecution == nil {
		t.Fatalf("completion = %+v", completion)
	}
	steps := completion.CompoundExecution.Steps
	if len(steps) != 2 || steps[0].Acknowledgement != operation.AckRejected || steps[1].Acknowledgement != operation.AckAccepted {
		t.Fatalf("steps = %+v", steps)
	}
}

func localWriteReply(request []byte, status byte, opcodeEcho bool) []byte {
	if opcodeEcho {
		return []byte{request[0], status, request[2], request[4], request[5]}
	}
	return []byte{status, request[2], request[4], request[5]}
}

type scriptReply struct {
	mu      sync.Mutex
	frames  []protocol.Frame
	respond func(protocol.Frame) (bool, *protocol.Frame, error)
	onSend  func()
}

func (s *scriptReply) Send(ctx context.Context, frame protocol.Frame) (bool, *protocol.Frame, error) {
	if ctx.Err() != nil {
		return false, nil, ctx.Err()
	}
	s.mu.Lock()
	s.frames = append(s.frames, frame)
	respond := s.respond
	hook := s.onSend
	s.mu.Unlock()
	if hook != nil {
		hook()
	}
	if respond == nil {
		return true, &protocol.Frame{Data: []byte{0x00}}, nil
	}
	return respond(frame)
}

func (s *scriptReply) snapshot() []protocol.Frame {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]protocol.Frame(nil), s.frames...)
}

func (s *scriptReply) reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.frames = nil
}

type slotSpec struct {
	startHour, startMinute, endHour, endMinute int
	temperature                                *semreg.Decimal
}

func ownerSource(owner *Owner) func() Owner {
	return func() Owner { return *owner }
}

func compoundOwner(instances map[semreg.NativeBindingID]byte) Owner {
	owner := liveOwner()
	owner.Origin = 0x31
	owner.Controller = 0x15
	if instances != nil {
		owner.Instances = instances
	}
	return owner
}

func programDispatch(id, binding string, weekday, count int, slots []slotSpec) semdriver.Dispatch {
	arguments := []semreg.TypedField{
		{ID: "ebusreg.time_program.weekday", Value: quantityValue("unit.one", strconv.Itoa(weekday), 0)},
		{ID: "ebusreg.time_program.slot_count", Value: quantityValue("unit.one", strconv.Itoa(count), 0)},
	}
	for index, slot := range slots {
		arguments = append(arguments,
			semreg.TypedField{ID: slotArgument(index, "start_hour"), Value: quantityValue("unit.one", strconv.Itoa(slot.startHour), 0)},
			semreg.TypedField{ID: slotArgument(index, "start_minute"), Value: quantityValue("unit.one", strconv.Itoa(slot.startMinute), 0)},
			semreg.TypedField{ID: slotArgument(index, "end_hour"), Value: quantityValue("unit.one", strconv.Itoa(slot.endHour), 0)},
			semreg.TypedField{ID: slotArgument(index, "end_minute"), Value: quantityValue("unit.one", strconv.Itoa(slot.endMinute), 0)},
		)
		if slot.temperature != nil {
			arguments = append(arguments, semreg.TypedField{
				ID:    slotArgument(index, "temperature"),
				Value: semreg.Value{Kind: semreg.ValueQuantity, Quantity: &semreg.Quantity{Number: *slot.temperature, Unit: "unit.celsius"}},
			})
		}
	}
	return fencedDispatch(id, binding, arguments)
}

func textDispatch(id, text string) semdriver.Dispatch {
	value := text
	return fencedDispatch(id, "system", []semreg.TypedField{{
		ID:    semreg.DefinitionID(id),
		Value: semreg.Value{Kind: semreg.ValueText, Text: &value},
	}})
}

func fencedDispatch(id, binding string, arguments []semreg.TypedField) semdriver.Dispatch {
	return semdriver.Dispatch{
		Request: semdriver.OperationRequest{
			BindingID:                semreg.NativeBindingID(binding),
			Operation:                semreg.DefinitionRef{ID: semreg.DefinitionID(id)},
			Arguments:                arguments,
			ExpectedDriverGeneration: "4",
			ExpectedSourceEpochID:    "epoch-r",
		},
		Guard: operation.GuardClaims{
			BindingID:        semreg.NativeBindingID(binding),
			SourceID:         "source-r",
			SourceEpochID:    "epoch-r",
			DriverGeneration: "4",
		},
	}
}

func slotArgument(index int, name string) semreg.DefinitionID {
	return semreg.DefinitionID("ebusreg.time_program.slot." + strconv.Itoa(index) + "." + name)
}

func decimalPtr(coefficient string, exponent int32) *semreg.Decimal {
	return &semreg.Decimal{Coefficient: coefficient, Exponent10: exponent}
}

func dhwTimerReply(contradict bool) func(protocol.Frame) (bool, *protocol.Frame, error) {
	return func(frame protocol.Frame) (bool, *protocol.Frame, error) {
		if len(frame.Data) == 0 {
			return true, &protocol.Frame{Data: []byte{0x00}}, nil
		}
		if frame.Data[0] == 0xa6 {
			return true, &protocol.Frame{Data: []byte{0x00}}, nil
		}
		if frame.Data[0] != 0xa5 || len(frame.Data) < 5 {
			return true, &protocol.Frame{Data: []byte{0x00}}, nil
		}
		if frame.Data[4] == 0 {
			return true, timerBody(0, 6, 10, 12, 20, 400), nil
		}
		temp := uint16(555)
		if contradict {
			temp = 556
		}
		return true, timerBody(0, 12, 20, 23, 50, temp), nil
	}
}

func heatingTimerReply() func(protocol.Frame) (bool, *protocol.Frame, error) {
	return func(frame protocol.Frame) (bool, *protocol.Frame, error) {
		if len(frame.Data) > 0 && frame.Data[0] == 0xa5 {
			return true, timerBody(0, 8, 15, 9, 45, 215), nil
		}
		return true, &protocol.Frame{Data: []byte{0x00}}, nil
	}
}

func timerBody(status, startHour, startMinute, endHour, endMinute byte, temperature uint16) *protocol.Frame {
	return &protocol.Frame{Data: []byte{status, startHour, startMinute, endHour, endMinute, byte(temperature), byte(temperature >> 8)}}
}

func registerReply(payloads map[uint16][]byte, tweak uint16, nonzeroWrites int) func(protocol.Frame) (bool, *protocol.Frame, error) {
	return registerReplyShape(payloads, tweak, nonzeroWrites, true)
}

func registerReplyShape(payloads map[uint16][]byte, tweak uint16, nonzeroWrites int, opcodeEcho bool) func(protocol.Frame) (bool, *protocol.Frame, error) {
	seenWrites := 0
	return func(frame protocol.Frame) (bool, *protocol.Frame, error) {
		data := frame.Data
		if len(data) < 6 || data[0] != 0x02 {
			return true, &protocol.Frame{Data: []byte{0x00}}, nil
		}
		if data[1] == 0x01 {
			seenWrites++
			status := byte(0)
			if seenWrites <= nonzeroWrites {
				status = 0x02
			}
			return true, &protocol.Frame{Data: localWriteReply(data, status, opcodeEcho)}, nil
		}
		register := binary.LittleEndian.Uint16(data[4:6])
		payload := append([]byte(nil), payloads[register]...)
		if tweak == register && len(payload) > 0 {
			payload[0] ^= 0x01
		}
		body := []byte{0x03, data[2], data[4], data[5]}
		body = append(body, payload...)
		return true, &protocol.Frame{Data: body}, nil
	}
}

func assertTimerFrame(t *testing.T, frame protocol.Frame, data []byte) {
	t.Helper()
	if frame.Source != 0x31 || frame.Target != 0x15 || frame.Primary != 0xb5 || frame.Secondary != 0x55 || frame.FrameType != protocol.FrameTypeInitiatorTarget {
		t.Fatalf("header = %+v", frame)
	}
	if !bytes.Equal(frame.Data, data) {
		t.Fatalf("data = %x want %x", frame.Data, data)
	}
}

func assertCompoundOutcome(t *testing.T, completion semdriver.Completion, outcome operation.Outcome, relation operation.ReadbackRelation, partial bool, binding semreg.NativeBindingID) {
	t.Helper()
	if completion.Outcome != outcome || completion.CompoundExecution == nil || completion.CompoundExecution.Partial != partial {
		t.Fatalf("completion = %+v", completion)
	}
	if completion.Readback == nil || completion.Readback.Relation != relation || completion.Readback.BindingID != binding {
		t.Fatalf("readback = %+v", completion.Readback)
	}
	if completion.Readback.DriverGeneration != "4" || completion.Readback.SourceEpochID != "epoch-r" || completion.Readback.SourceID != "source-r" {
		t.Fatalf("readback identity = %+v", completion.Readback)
	}
}

func frameDatas(frames []protocol.Frame) [][]byte {
	out := make([][]byte, len(frames))
	for i, frame := range frames {
		out[i] = frame.Data
	}
	return out
}
