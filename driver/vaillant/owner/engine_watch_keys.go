package owner

import (
	"fmt"

	"github.com/Project-Helianthus/helianthus-ebusgo/engine"
	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
)

// B555WatchKey is the b555 watch key type of the Vaillant owner.
type B555WatchKey struct {
	Target byte
	Opcode byte
	Zone   byte
	HC     byte

	Weekday    byte
	Slot       byte
	HasWeekday bool
	HasSlot    bool
}

// NewB555ProgramWatchKey implements the new b555 program watch key step of the Vaillant owner.
func NewB555ProgramWatchKey(target, opcode, zone, hc byte) B555WatchKey {
	return B555WatchKey{
		Target: target,
		Opcode: opcode,
		Zone:   zone,
		HC:     hc,
	}
}

// NewB555WatchKey implements the new b555 watch key step of the Vaillant owner.
func NewB555WatchKey(target, opcode, zone, hc, weekday, slot byte) B555WatchKey {
	return B555WatchKey{
		Target:     target,
		Opcode:     opcode,
		Zone:       zone,
		HC:         hc,
		Weekday:    weekday,
		Slot:       slot,
		HasWeekday: true,
		HasSlot:    true,
	}
}

// Canonical implements the canonical step of the Vaillant owner.
func (key B555WatchKey) Canonical() string {
	canonical := fmt.Sprintf(
		"b555:%02x:%02x:%02x:%02x:%02x:%02x",
		key.Target,
		key.Opcode,
		key.Zone,
		key.HC,
		b555SelectorByte(key.HasWeekday, key.Weekday),
		b555SelectorByte(key.HasSlot, key.Slot),
	)
	if b555WatchKeySelectorsValid(key) {
		return canonical
	}
	return fmt.Sprintf(
		"%s:raw:%d:%02x:%d:%02x",
		canonical,
		b555BoolByte(key.HasWeekday),
		key.Weekday,
		b555BoolByte(key.HasSlot),
		key.Slot,
	)
}

func (key B555WatchKey) String() string {
	return key.Canonical()
}

// Family implements the family step of the Vaillant owner.
func (key B555WatchKey) Family() engine.WatchFamily {
	return engine.WatchFamilyB555
}

// PassiveWatchKeyFromEvent implements the passive watch key from event step of the Vaillant owner.
func PassiveWatchKeyFromEvent(event engine.PassiveClassifiedEvent) (engine.WatchKey, bool) {
	if !event.HasRequest {
		return nil, false
	}
	return PassiveWatchKeyFromFrame(event.Request)
}

// PassiveWatchKeyFromFrame implements the passive watch key from frame step of the Vaillant owner.
func PassiveWatchKeyFromFrame(frame protocol.Frame) (engine.WatchKey, bool) {
	switch {
	case frame.Primary == 0xB5 && frame.Secondary == 0x24:
		return parsePassiveB524WatchKey(frame)
	case frame.Primary == 0xB5 && frame.Secondary == 0x09:
		return parsePassiveB509WatchKey(frame)
	default:
		return nil, false
	}
}

func parsePassiveB524WatchKey(frame protocol.Frame) (engine.WatchKey, bool) {
	if len(frame.Data) < 6 {
		return nil, false
	}
	opcode := frame.Data[0]
	if opcode != 0x02 && opcode != 0x06 {
		return nil, false
	}
	addr := uint16(frame.Data[4]) | uint16(frame.Data[5])<<8
	return engine.NewB524WatchKey(frame.Target, opcode, frame.Data[2], frame.Data[3], addr), true
}

func parsePassiveB509WatchKey(frame protocol.Frame) (engine.WatchKey, bool) {
	if len(frame.Data) < 3 {
		return nil, false
	}
	switch frame.Data[0] {
	case 0x0D, 0x29, 0x0E:
	default:
		return nil, false
	}
	addr := uint16(frame.Data[1])<<8 | uint16(frame.Data[2])
	return engine.NewB509WatchKey(frame.Target, addr), true
}

func b555SelectorByte(set bool, value byte) byte {
	if !set {
		return 0xFF
	}
	return value
}

func b555BoolByte(value bool) byte {
	if value {
		return 1
	}
	return 0
}

func b555WatchKeySelectorsValid(key B555WatchKey) bool {
	if key.HasWeekday != key.HasSlot {
		return false
	}
	if !key.HasWeekday {
		return key.Weekday == 0x00 && key.Slot == 0x00
	}
	return key.Weekday != 0xFF && key.Slot != 0xFF
}

func cloneWatchKey(key engine.WatchKey) (engine.WatchKey, error) {
	switch typed := key.(type) {
	case engine.B509WatchKey:
		return typed, nil
	case *engine.B509WatchKey:
		if typed == nil {
			return nil, fmt.Errorf("watch descriptor contains nil *engine.B509WatchKey")
		}
		return *typed, nil
	case engine.B516WatchKey:
		return typed, nil
	case *engine.B516WatchKey:
		if typed == nil {
			return nil, fmt.Errorf("watch descriptor contains nil *engine.B516WatchKey")
		}
		return *typed, nil
	case engine.B524WatchKey:
		return typed, nil
	case *engine.B524WatchKey:
		if typed == nil {
			return nil, fmt.Errorf("watch descriptor contains nil *engine.B524WatchKey")
		}
		return *typed, nil
	case B555WatchKey:
		return typed, nil
	case *B555WatchKey:
		if typed == nil {
			return nil, fmt.Errorf("watch descriptor contains nil *B555WatchKey")
		}
		return *typed, nil
	default:
		return nil, fmt.Errorf("watch descriptor key type %T is not supported for immutable catalog storage", key)
	}
}

func validateB555WatchKey(key engine.WatchKey) error {
	b555Key, ok := key.(B555WatchKey)
	if !ok {
		return nil
	}
	if b555Key.HasWeekday != b555Key.HasSlot {
		return fmt.Errorf("must set both weekday and slot selectors together")
	}
	if !b555Key.HasWeekday {
		if b555Key.Weekday != 0x00 || b555Key.Slot != 0x00 {
			return fmt.Errorf("program selectors must keep raw weekday/slot zero when unset")
		}
		return nil
	}
	if b555Key.Weekday == 0xFF || b555Key.Slot == 0xFF {
		return fmt.Errorf("timer selectors may not use 0xff sentinel")
	}
	return nil
}
