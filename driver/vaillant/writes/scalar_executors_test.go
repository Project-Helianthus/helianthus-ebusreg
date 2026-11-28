package writes

import (
	"context"
	"encoding/binary"
	"math"
	"sync"
	"testing"

	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
	semreg "github.com/Project-Helianthus/helianthus-semreg/semreg/v1"
	semdriver "github.com/Project-Helianthus/helianthus-semreg/semreg/v1/driver"
	"github.com/Project-Helianthus/helianthus-semreg/semreg/v1/operation"
)

func TestZoneTemperatureExecutor(t *testing.T) {
	bus := &scriptedBus{group: 0x03, register: 0x0022, readValue: float32Bytes(float32(21))}
	owner := liveOwner()
	executor := NewScalar(bus, func() Owner { return owner })
	completion, err := executor.Execute(context.Background(), zoneDispatch("zone-r", "4", "epoch-r", semreg.Decimal{Coefficient: "29415", Exponent10: -2}))
	if err != nil {
		t.Fatal(err)
	}
	if completion.Outcome != operation.OutcomeApplied || completion.Readback == nil || completion.Readback.Relation != operation.ReadbackConfirms {
		t.Fatalf("completion = %+v", completion)
	}
	if completion.Readback.BindingID != "zone-r" || completion.Readback.DriverGeneration != "4" || completion.Readback.SourceEpochID != "epoch-r" {
		t.Fatalf("readback identity = %+v", completion.Readback)
	}
	frames := bus.snapshot()
	if len(frames) != 2 {
		t.Fatalf("frames = %d", len(frames))
	}
	assertLocalFrame(t, frames[0], 0x10, 0x08, 0x01, 0x03, 0x17, 0x0022, float32Bytes(float32(21)))
	assertLocalFrame(t, frames[1], 0x10, 0x08, 0x00, 0x03, 0x17, 0x0022, nil)

	owner.Instances = map[semreg.NativeBindingID]byte{"zone-s": 0x2a}
	bus.group = 0x03
	bus.register = 0x0022
	bus.readValue = float32Bytes(float32(22))
	bus.reset()
	completion, err = executor.Execute(context.Background(), zoneDispatch("zone-s", "4", "epoch-r", semreg.Decimal{Coefficient: "29515", Exponent10: -2}))
	if err != nil || completion.Outcome != operation.OutcomeApplied {
		t.Fatalf("second completion = %+v err=%v", completion, err)
	}
	frames = bus.snapshot()
	if len(frames) != 2 {
		t.Fatalf("second frames = %d", len(frames))
	}
	assertLocalFrame(t, frames[0], 0x10, 0x08, 0x01, 0x03, 0x2a, 0x0022, float32Bytes(float32(22)))

	bus.reset()
	wrong := zoneDispatch("zone-s", "4", "epoch-r", semreg.Decimal{Coefficient: "29415", Exponent10: -2})
	wrong.Request.Arguments[0].Value.Quantity.Unit = "unit.celsius"
	completion, err = executor.Execute(context.Background(), wrong)
	if err != nil || completion.Outcome != operation.OutcomeFailedNoContact || len(bus.snapshot()) != 0 {
		t.Fatalf("celsius unit sent frames, completion %+v", completion)
	}
	wrong = zoneDispatch("zone-s", "4", "epoch-r", semreg.Decimal{Coefficient: "29415", Exponent10: -2})
	wrong.Request.Arguments[0].ID = "thermal.setpoint.other"
	if completion, err = executor.Execute(context.Background(), wrong); err != nil || completion.Outcome != operation.OutcomeFailedNoContact || len(bus.snapshot()) != 0 {
		t.Fatalf("wrong argument sent frames, completion %+v", completion)
	}
	wrong = zoneDispatch("zone-s", "4", "epoch-r", semreg.Decimal{Coefficient: "29415", Exponent10: -2})
	wrong.Request.Arguments = append(wrong.Request.Arguments, wrong.Request.Arguments[0])
	if completion, err = executor.Execute(context.Background(), wrong); err != nil || completion.Outcome != operation.OutcomeFailedNoContact || len(bus.snapshot()) != 0 {
		t.Fatalf("extra argument sent frames, completion %+v", completion)
	}
}

func TestReadbackContradiction(t *testing.T) {
	bus := &scriptedBus{group: 0x03, register: 0x0022, readValue: float32Bytes(float32(21))}
	executor := NewScalar(bus, func() Owner { return liveOwner() })
	completion, err := executor.Execute(context.Background(), zoneDispatch("zone-r", "4", "epoch-r", semreg.Decimal{Coefficient: "30015", Exponent10: -2}))
	if err != nil {
		t.Fatal(err)
	}
	if completion.Outcome != operation.OutcomeConflict || completion.Readback == nil || completion.Readback.Relation != operation.ReadbackContradicts {
		t.Fatalf("completion = %+v", completion)
	}
	frames := bus.snapshot()
	if len(frames) != 2 {
		t.Fatalf("frames = %d", len(frames))
	}
	assertLocalFrame(t, frames[0], 0x10, 0x08, 0x01, 0x03, 0x17, 0x0022, float32Bytes(float32(27)))
	assertLocalFrame(t, frames[1], 0x10, 0x08, 0x00, 0x03, 0x17, 0x0022, nil)

	bus.reset()
	bus.group = 0x04
	completion, err = executor.Execute(context.Background(), zoneDispatch("zone-r", "4", "epoch-r", semreg.Decimal{Coefficient: "29415", Exponent10: -2}))
	if err != nil || completion.Outcome == operation.OutcomeApplied || completion.Outcome == operation.OutcomeConflict {
		t.Fatalf("mismatched readback completion = %+v", completion)
	}
	if len(bus.snapshot()) != 2 {
		t.Fatalf("mismatched readback frames = %d", len(bus.snapshot()))
	}
}

func TestNotAdmittedSendsNothing(t *testing.T) {
	cases := []struct {
		name    string
		owner   func() Owner
		build   func() semdriver.Dispatch
		cancel  bool
		nilExec bool
	}{
		{name: "not admitted", owner: func() Owner { owner := liveOwner(); owner.Admitted = false; return owner }},
		{name: "epoch", owner: func() Owner { owner := liveOwner(); owner.Epoch = "other-epoch"; return owner }},
		{name: "generation", owner: func() Owner { owner := liveOwner(); owner.Generation = "8"; return owner }},
		{name: "source", owner: func() Owner { owner := liveOwner(); owner.Source = "other-source"; return owner }},
		{name: "empty generation", build: func() semdriver.Dispatch {
			dispatch := zoneDispatch("zone-r", "", "epoch-r", semreg.Decimal{Coefficient: "29415", Exponent10: -2})
			dispatch.Guard.DriverGeneration = ""
			return dispatch
		}},
		{name: "empty epoch", build: func() semdriver.Dispatch {
			dispatch := zoneDispatch("zone-r", "4", "", semreg.Decimal{Coefficient: "29415", Exponent10: -2})
			dispatch.Guard.SourceEpochID = ""
			return dispatch
		}},
		{name: "cancelled", cancel: true},
		{name: "compound", build: func() semdriver.Dispatch {
			dispatch := zoneDispatch("zone-r", "4", "epoch-r", semreg.Decimal{Coefficient: "29415", Exponent10: -2})
			dispatch.Request.Operation.ID = "ebusreg.operation.zone_heating_time_program"
			return dispatch
		}},
		{name: "unknown codec", build: func() semdriver.Dispatch {
			dispatch := zoneDispatch("zone-r", "4", "epoch-r", semreg.Decimal{Coefficient: "29415", Exponent10: -2})
			dispatch.Request.Operation.ID = "thermal.operation.set_zone_mode"
			dispatch.Request.Arguments[0].ID = "thermal.mode.zone"
			return dispatch
		}},
		{name: "dhw setpoint", build: func() semdriver.Dispatch {
			dispatch := zoneDispatch("zone-r", "4", "epoch-r", semreg.Decimal{Coefficient: "32315", Exponent10: -2})
			dispatch.Request.Operation.ID = "thermal.operation.set_dhw_temperature"
			dispatch.Request.Arguments[0].ID = "thermal.setpoint.dhw_temperature_setpoint"
			return dispatch
		}},
		{name: "missing binding", owner: func() Owner { owner := liveOwner(); owner.Instances = map[semreg.NativeBindingID]byte{}; return owner }},
		{name: "origin", owner: func() Owner { owner := liveOwner(); owner.Origin = 0; return owner }},
		{name: "nil executor", nilExec: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			bus := &scriptedBus{}
			ownerFn := tc.owner
			if ownerFn == nil {
				ownerFn = func() Owner { return liveOwner() }
			}
			var executor *Scalar
			if !tc.nilExec {
				executor = NewScalar(bus, ownerFn)
			}
			ctx := context.Background()
			if tc.cancel {
				var cancel context.CancelFunc
				ctx, cancel = context.WithCancel(ctx)
				cancel()
			}
			dispatch := zoneDispatch("zone-r", "4", "epoch-r", semreg.Decimal{Coefficient: "29415", Exponent10: -2})
			if tc.build != nil {
				dispatch = tc.build()
			}
			completion, err := executor.Execute(ctx, dispatch)
			if err != nil || completion.Outcome != operation.OutcomeFailedNoContact || completion.Readback != nil || len(bus.snapshot()) != 0 {
				t.Fatalf("completion=%+v err=%v frames=%d", completion, err, len(bus.snapshot()))
			}
		})
	}

	t.Run("recheck generation", func(t *testing.T) {
		bus := &scriptedBus{}
		calls := 0
		executor := NewScalar(bus, func() Owner {
			calls++
			owner := liveOwner()
			if calls >= 2 {
				owner.Generation = "9"
			}
			return owner
		})
		completion, err := executor.Execute(context.Background(), zoneDispatch("zone-r", "4", "epoch-r", semreg.Decimal{Coefficient: "29415", Exponent10: -2}))
		if err != nil || completion.Outcome != operation.OutcomeFailedNoContact || len(bus.snapshot()) != 0 || calls < 2 {
			t.Fatalf("completion=%+v err=%v frames=%d calls=%d", completion, err, len(bus.snapshot()), calls)
		}
	})
	t.Run("recheck instance", func(t *testing.T) {
		bus := &scriptedBus{}
		calls := 0
		executor := NewScalar(bus, func() Owner {
			calls++
			owner := liveOwner()
			if calls >= 2 {
				owner.Instances = map[semreg.NativeBindingID]byte{"zone-r": 0x18}
			}
			return owner
		})
		completion, err := executor.Execute(context.Background(), zoneDispatch("zone-r", "4", "epoch-r", semreg.Decimal{Coefficient: "29415", Exponent10: -2}))
		if err != nil || completion.Outcome != operation.OutcomeFailedNoContact || len(bus.snapshot()) != 0 {
			t.Fatalf("completion=%+v err=%v frames=%d", completion, err, len(bus.snapshot()))
		}
	})
}

func TestScalarCodecFrames(t *testing.T) {
	cases := []struct {
		name      string
		operation string
		argument  semreg.DefinitionID
		value     semreg.Value
		group     byte
		register  uint16
		instance  byte
		payload   []byte
		binding   semreg.NativeBindingID
	}{
		{
			name: "system off", operation: "thermal.operation.set_mode", argument: "thermal.mode.system",
			value: symbolValue("thermal.mode.system", "off"),
			group: 0x00, register: 0x0007, instance: 0x00, payload: []byte{0x01},
			binding: "system",
		},
		{
			name: "system auto", operation: "thermal.operation.set_mode", argument: "thermal.mode.system",
			value: symbolValue("thermal.mode.system", "auto"),
			group: 0x00, register: 0x0007, instance: 0x00, payload: []byte{0x00},
			binding: "system",
		},
		{
			name: "bivalence", operation: "ebusreg.operation.heating_outdoor_bivalence_temperature_limit", argument: "ebusreg.operation.heating_outdoor_bivalence_temperature_limit",
			value: quantityValue("unit.kelvin", "25315", -2),
			group: 0x00, register: 0x0023, instance: 0x00, payload: float32Bytes(float32(-20)),
			binding: "system",
		},
		{
			name: "curve", operation: "ebusreg.operation.circuit_heating_curve_factor", argument: "ebusreg.operation.circuit_heating_curve_factor",
			value: quantityValue("unit.one", "15", -1),
			group: 0x02, register: 0x000f, instance: 0x09, payload: float32Bytes(float32(1.5)),
			binding: "circuit-r",
		},
		{
			name: "humidity", operation: "ebusreg.operation.room_humidity_limit_maximum", argument: "ebusreg.operation.room_humidity_limit_maximum",
			value: quantityValue("unit.ratio", "8", -1),
			group: 0x00, register: 0x000e, instance: 0x00, payload: []byte{80, 0},
			binding: "system",
		},
		{
			name: "room control", operation: "ebusreg.operation.circuit_room_temperature_influence_mode", argument: "ebusreg.operation.circuit_room_temperature_influence_mode",
			value: symbolValue("ebusreg.circuit_room_temperature_influence_mode", "thermostat"),
			group: 0x02, register: 0x0015, instance: 0x09, payload: []byte{2, 0},
			binding: "circuit-r",
		},
		{
			name: "access code", operation: "ebusreg.operation.installer_access_code", argument: "ebusreg.operation.installer_access_code",
			value: quantityValue("unit.one", "42", 0),
			group: 0x00, register: 0x0076, instance: 0x00, payload: []byte{42, 0},
			binding: "system",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			bus := &scriptedBus{group: tc.group, register: tc.register, readValue: append([]byte(nil), tc.payload...)}
			owner := liveOwner()
			owner.Instances = map[semreg.NativeBindingID]byte{tc.binding: tc.instance, "zone-r": 0x17}
			executor := NewScalar(bus, func() Owner { return owner })
			dispatch := zoneDispatch(tc.binding, "4", "epoch-r", semreg.Decimal{Coefficient: "0", Exponent10: 0})
			dispatch.Request.Operation.ID = semreg.DefinitionID(tc.operation)
			dispatch.Request.Arguments = []semreg.TypedField{{ID: tc.argument, Value: tc.value}}
			completion, err := executor.Execute(context.Background(), dispatch)
			if err != nil || completion.Outcome != operation.OutcomeApplied {
				t.Fatalf("completion=%+v err=%v", completion, err)
			}
			frames := bus.snapshot()
			if len(frames) != 2 {
				t.Fatalf("frames = %d", len(frames))
			}
			assertLocalFrame(t, frames[0], 0x10, 0x08, 0x01, tc.group, tc.instance, tc.register, tc.payload)
			assertLocalFrame(t, frames[1], 0x10, 0x08, 0x00, tc.group, tc.instance, tc.register, nil)
		})
	}

	bus := &scriptedBus{}
	owner := liveOwner()
	executor := NewScalar(bus, func() Owner { return owner })
	dispatch := zoneDispatch("circuit-r", "4", "epoch-r", semreg.Decimal{Coefficient: "29415", Exponent10: -2})
	dispatch.Request.Operation.ID = "ebusreg.operation.circuit_heating_curve_factor"
	dispatch.Request.Arguments[0].ID = "ebusreg.operation.circuit_heating_curve_factor"
	completion, err := executor.Execute(context.Background(), dispatch)
	if err != nil || completion.Outcome != operation.OutcomeFailedNoContact || len(bus.snapshot()) != 0 {
		t.Fatalf("kelvin curve completion=%+v frames=%d", completion, len(bus.snapshot()))
	}
	mode := zoneDispatch("system", "4", "epoch-r", semreg.Decimal{})
	mode.Request.Operation.ID = "thermal.operation.set_mode"
	mode.Request.Arguments = []semreg.TypedField{{ID: "thermal.mode.system", Value: symbolValue("thermal.mode.system", "heat")}}
	if completion, err = executor.Execute(context.Background(), mode); err != nil || completion.Outcome != operation.OutcomeFailedNoContact || len(bus.snapshot()) != 0 {
		t.Fatalf("undocumented system mode sent frames")
	}
}

type scriptedBus struct {
	mu        sync.Mutex
	frames    []protocol.Frame
	group     byte
	register  uint16
	readValue []byte
}

func (s *scriptedBus) Send(_ context.Context, frame protocol.Frame) (bool, *protocol.Frame, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	copied := frame
	copied.Data = append([]byte(nil), frame.Data...)
	s.frames = append(s.frames, copied)
	if len(frame.Data) > 1 && frame.Data[1] == 0x00 {
		payload := []byte{0x03, s.group, byte(s.register), byte(s.register >> 8)}
		payload = append(payload, s.readValue...)
		return true, &protocol.Frame{Data: payload}, nil
	}
	return true, &protocol.Frame{Data: []byte{0x00}}, nil
}

func (s *scriptedBus) snapshot() []protocol.Frame {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]protocol.Frame(nil), s.frames...)
}

func (s *scriptedBus) reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.frames = nil
}

func liveOwner() Owner {
	return Owner{
		Admitted:   true,
		Generation: "4",
		Epoch:      "epoch-r",
		Source:     "source-r",
		Origin:     0x10,
		Controller: 0x08,
		Instances:  map[semreg.NativeBindingID]byte{"zone-r": 0x17},
	}
}

func zoneDispatch(binding semreg.NativeBindingID, generation, epoch string, number semreg.Decimal) semdriver.Dispatch {
	return semdriver.Dispatch{
		Request: semdriver.OperationRequest{
			BindingID: binding,
			Operation: semreg.DefinitionRef{ID: "thermal.operation.set_temperature"},
			Arguments: []semreg.TypedField{{
				ID: "thermal.setpoint.temperature_setpoint",
				Value: semreg.Value{
					Kind:     semreg.ValueQuantity,
					Quantity: &semreg.Quantity{Number: number, Unit: "unit.kelvin"},
				},
			}},
			ExpectedDriverGeneration: semreg.Uint64(generation),
			ExpectedSourceEpochID:    semreg.SourceEpochID(epoch),
		},
		Guard: operation.GuardClaims{
			BindingID:        binding,
			SourceID:         "source-r",
			SourceEpochID:    semreg.SourceEpochID(epoch),
			DriverGeneration: semreg.Uint64(generation),
		},
	}
}

func quantityValue(unit semreg.DefinitionID, coefficient string, exponent int32) semreg.Value {
	return semreg.Value{
		Kind: semreg.ValueQuantity,
		Quantity: &semreg.Quantity{
			Number: semreg.Decimal{Coefficient: coefficient, Exponent10: exponent},
			Unit:   unit,
		},
	}
}

func symbolValue(namespace semreg.DefinitionID, token string) semreg.Value {
	return semreg.Value{
		Kind:   semreg.ValueSymbol,
		Symbol: &semreg.Symbol{Namespace: namespace, Token: token, Known: true},
	}
}

func float32Bytes(value float32) []byte {
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], math.Float32bits(value))
	return buf[:]
}

func assertLocalFrame(t *testing.T, frame protocol.Frame, source, target, access, group, instance byte, register uint16, payload []byte) {
	t.Helper()
	if frame.Source != source || frame.Target != target || frame.Primary != 0xb5 || frame.Secondary != 0x24 || frame.FrameType != protocol.FrameTypeInitiatorTarget {
		t.Fatalf("frame header = %+v", frame)
	}
	want := []byte{0x02, access, group, instance, byte(register), byte(register >> 8)}
	want = append(want, payload...)
	if len(frame.Data) != len(want) {
		t.Fatalf("data = %x want %x", frame.Data, want)
	}
	for i := range want {
		if frame.Data[i] != want[i] {
			t.Fatalf("data = %x want %x", frame.Data, want)
		}
	}
}
