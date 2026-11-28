package writes

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"testing"

	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
	semreg "github.com/Project-Helianthus/helianthus-semreg/semreg/v1"
	semdriver "github.com/Project-Helianthus/helianthus-semreg/semreg/v1/driver"
	"github.com/Project-Helianthus/helianthus-semreg/semreg/v1/operation"
)

func TestOperationLimits(t *testing.T) {
	t.Run("device query", func(t *testing.T) { assertDeviceQuery(t) })
	t.Run("instances stay apart", func(t *testing.T) { assertInstanceLimits(t) })
	t.Run("failed query", func(t *testing.T) { assertFailedQuery(t) })
	t.Run("mismatched echo", func(t *testing.T) { assertMismatchedEcho(t) })
	t.Run("timer range", func(t *testing.T) { assertTimerRange(t) })
	t.Run("no ceiling without a query", func(t *testing.T) { assertNoInventedCeiling(t) })
	t.Run("boiler has no query", func(t *testing.T) { assertBoilerHasNoLimit(t) })
	t.Run("second populate", func(t *testing.T) { assertSecondPopulate(t) })
	t.Run("refused owner", func(t *testing.T) { assertRefusedOwner(t) })
	t.Run("cancelled populate", func(t *testing.T) { assertCancelledPopulate(t) })
}

func assertDeviceQuery(t *testing.T) {
	t.Helper()
	bus := &scriptReply{respond: limitMode{}.respond}
	owner := compoundOwner(map[semreg.NativeBindingID]byte{"zone-r": 0x17})
	limits := NewLimits(bus, func() Owner { return owner })
	if err := limits.Populate(context.Background()); err != nil {
		t.Fatal(err)
	}
	frames := bus.snapshot()
	if countPrefix(frames, []byte{0x01, 0x03, 0x17, 0x22, 0x00}) != 1 {
		t.Fatalf("zone query count = %d", countPrefix(frames, []byte{0x01, 0x03, 0x17, 0x22, 0x00}))
	}
	if countPrefix(frames, []byte{0x01, 0x01, 0xff, 0x04, 0x00}) != 1 {
		t.Fatalf("hot water query count = %d", countPrefix(frames, []byte{0x01, 0x01, 0xff, 0x04, 0x00}))
	}
	if countPrefix(frames, []byte{0x01, 0x00, 0x00, 0x6c, 0x00}) != 1 || countPrefix(frames, []byte{0x01, 0x00, 0x00, 0x6f, 0x00}) != 1 {
		t.Fatal("string registers were not queried")
	}
	assertNoWriteDuringPopulate(t, frames)
	for _, frame := range frames {
		if frame.Source != 0x31 || frame.Target != 0x15 || frame.Primary != 0xb5 {
			t.Fatalf("header = %+v", frame)
		}
		if len(frame.Data) >= 5 && bytes.Equal(frame.Data[:5], []byte{0x01, 0x03, 0x17, 0x22, 0x00}) && frame.Secondary != 0x24 {
			t.Fatalf("zone query secondary = %x", frame.Secondary)
		}
	}
	zoneReply := floatSpan(0x03, 0x0022, 10, 25, 0.5)
	zone, ok := limits.Published("thermal.operation.set_temperature", 0x17)
	if !ok || zone.Origin != semdriver.LimitDeviceQuery || zone.Validate() != nil {
		t.Fatalf("zone limits = %+v ok=%v", zone, ok)
	}
	assertQuantity(t, zone.Min, "28315", -2, "unit.kelvin")
	assertQuantity(t, zone.Max, "29815", -2, "unit.kelvin")
	assertQuantity(t, zone.Step, "5", -1, "unit.kelvin")
	if zone.Evidence[0].Digest != digestOf(zoneReply) {
		t.Fatalf("zone digest = %s", zone.Evidence[0].Digest)
	}
	if !limits.Admit("thermal.operation.set_temperature", 0x17, kelvinValue("29415", -2)) {
		t.Fatal("21 degrees was refused")
	}
	if limits.Admit("thermal.operation.set_temperature", 0x17, kelvinValue("28015", -2)) || limits.Admit("thermal.operation.set_temperature", 0x17, kelvinValue("30015", -2)) {
		t.Fatal("a value outside 10 to 25 degrees was admitted")
	}
	if limits.Admit("thermal.operation.set_temperature", 0x17, textValue("21")) == false {
		t.Fatal("text was refused")
	}
	dhw, ok := limits.Published("thermal.operation.set_dhw_temperature", 0xff)
	if !ok {
		t.Fatal("hot water limits missing")
	}
	assertQuantity(t, dhw.Min, "30915", -2, "unit.kelvin")
	assertQuantity(t, dhw.Max, "33315", -2, "unit.kelvin")
	assertQuantity(t, dhw.Step, "1", 0, "unit.kelvin")
	if !limits.Admit("thermal.operation.set_dhw_temperature", 0xff, kelvinValue("30915", -2)) || !limits.Admit("thermal.operation.set_dhw_temperature", 0xff, kelvinValue("33315", -2)) {
		t.Fatal("boundary was refused")
	}
	if !limits.Admit("thermal.operation.set_dhw_temperature", 0xff, kelvinValue("31315", -2)) || !limits.Admit("thermal.operation.set_dhw_temperature", 0xff, kelvinValue("30965", -2)) {
		t.Fatal("an in-range value was refused")
	}
	if limits.Admit("thermal.operation.set_dhw_temperature", 0xff, kelvinValue("34315", -2)) || limits.Admit("thermal.operation.set_dhw_temperature", 0xff, kelvinValue("29315", -2)) {
		t.Fatal("a historical range admitted a value the query excludes")
	}
	if _, published := limits.Published("ebusreg.operation.installer_name", 0x00); published {
		t.Fatal("a string register was decoded as a quantity")
	}
}

func assertInstanceLimits(t *testing.T) {
	t.Helper()
	bus := &scriptReply{respond: limitMode{}.respond}
	owner := compoundOwner(map[semreg.NativeBindingID]byte{"zone-r": 0x17, "zone-s": 0x2a})
	limits := NewLimits(bus, func() Owner { return owner })
	if err := limits.Populate(context.Background()); err != nil {
		t.Fatal(err)
	}
	if !limits.Admit("thermal.operation.set_temperature", 0x17, kelvinValue("29415", -2)) {
		t.Fatal("instance 17 refused 21 degrees")
	}
	if limits.Admit("thermal.operation.set_temperature", 0x2a, kelvinValue("29415", -2)) {
		t.Fatal("instance 2a reused the other description")
	}
	if !limits.Admit("thermal.operation.set_temperature", 0x2a, kelvinValue("31815", -2)) {
		t.Fatal("instance 2a refused 45 degrees")
	}
}

func assertFailedQuery(t *testing.T) {
	t.Helper()
	bus := &scriptReply{respond: limitMode{failHotWater: true}.respond}
	owner := compoundOwner(map[semreg.NativeBindingID]byte{"zone-r": 0x17})
	limits := NewLimits(bus, func() Owner { return owner })
	if err := limits.Populate(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, ok := limits.Published("thermal.operation.set_dhw_temperature", 0xff); ok {
		t.Fatal("a failed query published a range")
	}
	if !limits.Admit("thermal.operation.set_dhw_temperature", 0xff, kelvinValue("29315", -2)) || !limits.Admit("thermal.operation.set_dhw_temperature", 0xff, kelvinValue("35315", -2)) {
		t.Fatal("a failed query refused a value")
	}
}

func assertMismatchedEcho(t *testing.T) {
	t.Helper()
	bus := &scriptReply{respond: limitMode{badZoneEcho: true}.respond}
	owner := compoundOwner(map[semreg.NativeBindingID]byte{"zone-r": 0x17})
	limits := NewLimits(bus, func() Owner { return owner })
	if err := limits.Populate(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, ok := limits.Published("thermal.operation.set_temperature", 0x17); ok {
		t.Fatal("a mismatched description was published")
	}
	if !limits.Admit("thermal.operation.set_temperature", 0x17, kelvinValue("28015", -2)) || !limits.Admit("thermal.operation.set_temperature", 0x17, kelvinValue("31315", -2)) {
		t.Fatal("a mismatched description refused a value")
	}
}

func assertTimerRange(t *testing.T) {
	t.Helper()
	bus := &scriptReply{respond: limitMode{}.respond}
	owner := compoundOwner(map[semreg.NativeBindingID]byte{"hot-water": 0x01})
	limits := NewLimits(bus, func() Owner { return owner })
	if err := limits.Populate(context.Background()); err != nil {
		t.Fatal(err)
	}
	timer, ok := limits.Published("ebusreg.operation.dhw_time_program", 0xff)
	if !ok || timer.Step != nil || timer.Origin != semdriver.LimitDeviceQuery {
		t.Fatalf("timer limits = %+v ok=%v", timer, ok)
	}
	assertQuantity(t, timer.Min, "30815", -2, "unit.kelvin")
	assertQuantity(t, timer.Max, "33815", -2, "unit.kelvin")
	body := []byte{0x00, 3, 10, 10, 1, 1, 35, 65, 0}
	if timer.Evidence[0].Digest != digestOf(body) || timer.Validate() != nil {
		t.Fatalf("timer evidence = %+v", timer.Evidence)
	}
	bus.reset()
	executor := NewCompound(bus, func() Owner { return owner }, limits)
	low := programDispatch("ebusreg.operation.dhw_time_program", "hot-water", 2, 1, []slotSpec{{0, 0, 1, 0, decimalPtr("34", 0)}})
	completion, err := executor.Execute(context.Background(), low)
	if err != nil || completion.Outcome != operation.OutcomeRejected || completion.CompoundExecution != nil || countOpcode(bus.snapshot(), 0xa6) != 0 {
		t.Fatalf("low completion=%+v frames=%d err=%v", completion, len(bus.snapshot()), err)
	}
	bus.reset()
	high := programDispatch("ebusreg.operation.dhw_time_program", "hot-water", 2, 1, []slotSpec{{0, 0, 1, 0, decimalPtr("35", 0)}})
	completion, err = executor.Execute(context.Background(), high)
	if err != nil || completion.Outcome == operation.OutcomeRejected || countOpcode(bus.snapshot(), 0xa6) != 1 {
		t.Fatalf("boundary completion=%+v frames=%x err=%v", completion, frameDatas(bus.snapshot()), err)
	}
	if !bytes.Contains(bus.snapshot()[0].Data, []byte{0x5e, 0x01}) {
		t.Fatalf("boundary frame = %x", bus.snapshot()[0].Data)
	}
	bus.reset()
	late := programDispatch("ebusreg.operation.dhw_time_program", "hot-water", 2, 1, []slotSpec{{25, 0, 1, 0, decimalPtr("40", 0)}})
	completion, err = executor.Execute(context.Background(), late)
	if err != nil || completion.Outcome != operation.OutcomeRejected || len(bus.snapshot()) != 0 {
		t.Fatalf("hour completion=%+v frames=%d err=%v", completion, len(bus.snapshot()), err)
	}
	bus.reset()
	wide := fourSlotProgram()
	completion, err = executor.Execute(context.Background(), wide)
	if err != nil || completion.Outcome != operation.OutcomeRejected || countOpcode(bus.snapshot(), 0xa6) != 0 {
		t.Fatalf("wide completion=%+v frames=%d err=%v", completion, len(bus.snapshot()), err)
	}
}

func assertNoInventedCeiling(t *testing.T) {
	t.Helper()
	bus := &scriptReply{respond: func(frame protocol.Frame) (bool, *protocol.Frame, error) {
		if len(frame.Data) > 0 && frame.Data[0] == 0xa5 {
			return true, timerBody(0, frame.Data[4], 0, byte(int(frame.Data[4])+1), 0, 400), nil
		}
		return true, &protocol.Frame{Data: []byte{0x00}}, nil
	}}
	owner := compoundOwner(map[semreg.NativeBindingID]byte{"hot-water": 0x01})
	limits := NewLimits(bus, func() Owner { return owner })
	executor := NewCompound(bus, func() Owner { return owner }, limits)
	completion, err := executor.Execute(context.Background(), fourSlotProgram())
	if err != nil {
		t.Fatal(err)
	}
	frames := bus.snapshot()
	if countOpcode(frames, 0xa6) != 4 {
		t.Fatalf("writes = %d data = %x", countOpcode(frames, 0xa6), frameDatas(frames))
	}
	seen := map[byte]int{}
	for _, frame := range frames {
		if len(frame.Data) > 5 && frame.Data[0] == 0xa6 {
			seen[frame.Data[4]]++
			if frame.Data[5] != 4 || frame.Data[6] != frame.Data[4] {
				t.Fatalf("slot frame = %x", frame.Data)
			}
		}
	}
	for index := byte(0); index < 4; index++ {
		if seen[index] != 1 {
			t.Fatalf("slot counts = %+v", seen)
		}
	}
	if completion.Outcome == operation.OutcomeRejected {
		t.Fatalf("completion = %+v", completion)
	}
}

func assertBoilerHasNoLimit(t *testing.T) {
	t.Helper()
	bus := &scriptReply{respond: limitMode{}.respond}
	owner := compoundOwner(map[semreg.NativeBindingID]byte{"zone-r": 0x17})
	limits := NewLimits(bus, func() Owner { return owner })
	if err := limits.Populate(context.Background()); err != nil {
		t.Fatal(err)
	}
	for _, frame := range bus.snapshot() {
		if frame.Secondary == 0x09 || (len(frame.Data) > 0 && frame.Data[0] == 0x0e) {
			t.Fatalf("boiler frame = %+v", frame)
		}
	}
	if _, ok := limits.Published("ebusreg.operation.maximum_heating_flow_temperature_setpoint", 0); ok {
		t.Fatal("boiler limits were published")
	}
	if !limits.Admit("ebusreg.operation.maximum_heating_flow_temperature_setpoint", 0, kelvinValue("0", 0)) || !limits.Admit("ebusreg.operation.maximum_heating_flow_temperature_setpoint", 0, kelvinValue("1000", 0)) {
		t.Fatal("a boiler value was refused")
	}
}

func assertSecondPopulate(t *testing.T) {
	t.Helper()
	bus := &scriptReply{respond: limitMode{}.respond}
	owner := compoundOwner(map[semreg.NativeBindingID]byte{"zone-r": 0x17})
	limits := NewLimits(bus, func() Owner { return owner })
	if err := limits.Populate(context.Background()); err != nil {
		t.Fatal(err)
	}
	before := len(bus.snapshot())
	if before == 0 {
		t.Fatal("first populate sent nothing")
	}
	if err := limits.Populate(context.Background()); err != nil {
		t.Fatal(err)
	}
	if len(bus.snapshot()) != before {
		t.Fatalf("second populate sent %d more frames", len(bus.snapshot())-before)
	}
}

func assertRefusedOwner(t *testing.T) {
	t.Helper()
	bus := &scriptReply{respond: limitMode{}.respond}
	owner := compoundOwner(map[semreg.NativeBindingID]byte{"zone-r": 0x17})
	owner.Admitted = false
	limits := NewLimits(bus, func() Owner { return owner })
	if err := limits.Populate(context.Background()); err != nil || len(bus.snapshot()) != 0 {
		t.Fatalf("refused populate err=%v frames=%d", err, len(bus.snapshot()))
	}
	owner.Admitted = true
	if err := limits.Populate(context.Background()); err != nil || len(bus.snapshot()) == 0 {
		t.Fatalf("later populate err=%v frames=%d", err, len(bus.snapshot()))
	}
}

func assertCancelledPopulate(t *testing.T) {
	t.Helper()
	bus := &scriptReply{respond: limitMode{}.respond}
	owner := compoundOwner(map[semreg.NativeBindingID]byte{"zone-r": 0x17})
	limits := NewLimits(bus, func() Owner { return owner })
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := limits.Populate(ctx); err == nil || len(bus.snapshot()) != 0 {
		t.Fatalf("cancelled err=%v frames=%d", err, len(bus.snapshot()))
	}
	if err := limits.Populate(context.Background()); err != nil || len(bus.snapshot()) == 0 {
		t.Fatalf("retry err=%v frames=%d", err, len(bus.snapshot()))
	}
}

type limitMode struct {
	failHotWater bool
	badZoneEcho  bool
}

func (m limitMode) respond(frame protocol.Frame) (bool, *protocol.Frame, error) {
	data := frame.Data
	if len(data) >= 5 && data[0] == 0x01 {
		group := data[1]
		instance := data[2]
		register := binary.LittleEndian.Uint16(data[3:5])
		if m.failHotWater && group == 0x01 && instance == 0xff && register == 0x0004 {
			return false, nil, errors.New("absent")
		}
		min, max, step := float32(0), float32(1), float32(1)
		echo := group
		if group == 0x03 && register == 0x0022 && instance == 0x17 {
			min, max, step = 10, 25, 0.5
			if m.badZoneEcho {
				echo = 0
			}
		}
		if group == 0x03 && register == 0x0022 && instance == 0x2a {
			min, max, step = 40, 50, 1
		}
		if group == 0x01 && register == 0x0004 && instance == 0xff {
			min, max, step = 36, 60, 1
		}
		return true, &protocol.Frame{Data: floatSpan(echo, register, min, max, step)}, nil
	}
	if len(data) >= 3 && data[0] == 0xa3 {
		body := []byte{0x00, 12, 10, 0, 1, 1, 0xff, 0xff, 0}
		if data[1] == 0xff && data[2] == 0x02 {
			body = []byte{0x00, 3, 10, 10, 1, 1, 35, 65, 0}
		}
		return true, &protocol.Frame{Data: append([]byte(nil), body...)}, nil
	}
	if len(data) > 0 && data[0] == 0xa5 && len(data) >= 5 {
		return true, timerBody(0, 0, 0, 1, 0, 350), nil
	}
	return true, &protocol.Frame{Data: []byte{0x00}}, nil
}

func floatSpan(group byte, register uint16, min, max, step float32) []byte {
	body := []byte{group, byte(register), byte(register >> 8)}
	body = append(body, float32Bytes(min)...)
	body = append(body, float32Bytes(max)...)
	body = append(body, float32Bytes(step)...)
	return body
}

func fourSlotProgram() semdriver.Dispatch {
	slots := make([]slotSpec, 4)
	for index := range slots {
		slots[index] = slotSpec{index, 0, index + 1, 0, decimalPtr("40", 0)}
	}
	return programDispatch("ebusreg.operation.dhw_time_program", "hot-water", 1, 4, slots)
}

func assertNoWriteDuringPopulate(t *testing.T, frames []protocol.Frame) {
	t.Helper()
	for _, frame := range frames {
		if frame.Secondary == 0x09 || (len(frame.Data) > 0 && (frame.Data[0] == 0xa6 || frame.Data[0] == 0x0e)) {
			t.Fatalf("unexpected frame %+v", frame)
		}
		if len(frame.Data) >= 2 && frame.Data[0] == 0x02 && frame.Data[1] == 0x01 {
			t.Fatalf("write during populate %x", frame.Data)
		}
	}
}

func countPrefix(frames []protocol.Frame, prefix []byte) int {
	count := 0
	for _, frame := range frames {
		if len(frame.Data) >= len(prefix) && bytes.Equal(frame.Data[:len(prefix)], prefix) {
			count++
		}
	}
	return count
}

func countOpcode(frames []protocol.Frame, opcode byte) int {
	count := 0
	for _, frame := range frames {
		if len(frame.Data) > 0 && frame.Data[0] == opcode {
			count++
		}
	}
	return count
}

func assertQuantity(t *testing.T, got *semreg.Quantity, coefficient string, exponent int32, unit semreg.DefinitionID) {
	t.Helper()
	if got == nil || got.Unit != unit || got.Number.Coefficient != coefficient || got.Number.Exponent10 != exponent {
		if got == nil {
			t.Fatal("quantity is absent")
		}
		t.Fatalf("quantity = %s e%d %s", got.Number.Coefficient, got.Number.Exponent10, got.Unit)
	}
}

func kelvinValue(coefficient string, exponent int32) semreg.Value {
	return quantityValue("unit.kelvin", coefficient, exponent)
}

func textValue(text string) semreg.Value {
	value := text
	return semreg.Value{Kind: semreg.ValueText, Text: &value}
}

func digestOf(raw []byte) semreg.Digest {
	sum := sha256.Sum256(raw)
	return semreg.Digest("sha256:" + hex.EncodeToString(sum[:]))
}
