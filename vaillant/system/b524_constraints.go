package system

// b524ConstraintSelector keys the unscoped constraint catalog. Register is the
// little-endian register number (RRlo | RRhi<<8), the same value the request
// builder and response decoder use. Previous revisions keyed on the
// byte-swapped value (RR 0x0001 stored as 0x0100) and ignored the opcode.
type b524ConstraintSelector struct {
	Opcode   byte
	Group    byte
	Register uint16
}

type b524Constraint struct {
	Type string
	Min  float64
	Max  float64
	Step float64
}

func (constraint b524Constraint) mapValue() map[string]any {
	return map[string]any{
		"type": constraint.Type,
		"min":  constraint.Min,
		"max":  constraint.Max,
		"step": constraint.Step,
	}
}

// b524StaticConstraintsCatalog holds controller-unscoped, OP02 (local) bounds.
// Every entry is an observation from one controller family and must not be
// treated as a universal bound; controller-scoped bounds live in
// b524_controller_profile.go. Entries for OP06 (remote) are not carried: no
// evidence ties these ranges to the remote namespace.
var b524StaticConstraintsCatalog = map[b524ConstraintSelector]b524Constraint{
	// OP02/GG00/RR0001 (system_dhw_bivalence_point) is deliberately absent: its
	// maximum is controller-profile scoped (BASV0 SW0217/HW8803/API1 = 20, the
	// earlier unscoped -20..50 is a BASV2 observation). Use LookupB524Bound.
	{Opcode: B524OpcodeLocal, Group: 0x00, Register: 0x0002}: {Type: "f32_range", Min: -26, Max: 10, Step: 1},
	{Opcode: B524OpcodeLocal, Group: 0x00, Register: 0x0003}: {Type: "u16_range", Min: 0, Max: 12, Step: 1},
	{Opcode: B524OpcodeLocal, Group: 0x00, Register: 0x0004}: {Type: "u16_range", Min: 0, Max: 300, Step: 10},
	{Opcode: B524OpcodeLocal, Group: 0x00, Register: 0x0080}: {Type: "f32_range", Min: -10, Max: 10, Step: 1},

	{Opcode: B524OpcodeLocal, Group: 0x01, Register: 0x0001}: {Type: "u16_range", Min: 0, Max: 1, Step: 1},
	{Opcode: B524OpcodeLocal, Group: 0x01, Register: 0x0002}: {Type: "u8_range", Min: 0, Max: 1, Step: 1},
	{Opcode: B524OpcodeLocal, Group: 0x01, Register: 0x0003}: {Type: "u16_range", Min: 0, Max: 2, Step: 1},
	{Opcode: B524OpcodeLocal, Group: 0x01, Register: 0x0004}: {Type: "f32_range", Min: 35, Max: 70, Step: 1},
	{Opcode: B524OpcodeLocal, Group: 0x01, Register: 0x0005}: {Type: "f32_range", Min: 0, Max: 99, Step: 1},
	{Opcode: B524OpcodeLocal, Group: 0x01, Register: 0x0006}: {Type: "u8_range", Min: 0, Max: 1, Step: 1},

	{Opcode: B524OpcodeLocal, Group: 0x02, Register: 0x0001}: {Type: "u16_range", Min: 1, Max: 2, Step: 1},
	{Opcode: B524OpcodeLocal, Group: 0x02, Register: 0x0002}: {Type: "u16_range", Min: 0, Max: 4, Step: 1},
	{Opcode: B524OpcodeLocal, Group: 0x02, Register: 0x0004}: {Type: "f32_range", Min: 15, Max: 80, Step: 1},
	{Opcode: B524OpcodeLocal, Group: 0x02, Register: 0x0005}: {Type: "u8_range", Min: 0, Max: 1, Step: 1},
	{Opcode: B524OpcodeLocal, Group: 0x02, Register: 0x0006}: {Type: "u8_range", Min: 0, Max: 1, Step: 1},

	{Opcode: B524OpcodeLocal, Group: 0x03, Register: 0x0001}: {Type: "u16_range", Min: 0, Max: 2, Step: 1},
	{Opcode: B524OpcodeLocal, Group: 0x03, Register: 0x0002}: {Type: "f32_range", Min: 15, Max: 30, Step: 0.5},
	{Opcode: B524OpcodeLocal, Group: 0x03, Register: 0x0005}: {Type: "f32_range", Min: 5, Max: 30, Step: 1},
	{Opcode: B524OpcodeLocal, Group: 0x03, Register: 0x0006}: {Type: "u16_range", Min: 0, Max: 2, Step: 1},

	{Opcode: B524OpcodeLocal, Group: 0x04, Register: 0x0001}: {Type: "u8_range", Min: 0, Max: 1, Step: 1},
	{Opcode: B524OpcodeLocal, Group: 0x04, Register: 0x0002}: {Type: "u8_range", Min: 0, Max: 1, Step: 1},
	{Opcode: B524OpcodeLocal, Group: 0x04, Register: 0x0003}: {Type: "f32_range", Min: -40, Max: 155, Step: 1},
	{Opcode: B524OpcodeLocal, Group: 0x04, Register: 0x0004}: {Type: "f32_range", Min: 0, Max: 99, Step: 1},
	{Opcode: B524OpcodeLocal, Group: 0x04, Register: 0x0005}: {Type: "f32_range", Min: 110, Max: 150, Step: 1},
	{Opcode: B524OpcodeLocal, Group: 0x04, Register: 0x0006}: {Type: "f32_range", Min: 75, Max: 115, Step: 1},

	{Opcode: B524OpcodeLocal, Group: 0x05, Register: 0x0001}: {Type: "f32_range", Min: 0, Max: 99, Step: 1},
	{Opcode: B524OpcodeLocal, Group: 0x05, Register: 0x0002}: {Type: "f32_range", Min: 2, Max: 25, Step: 1},
	{Opcode: B524OpcodeLocal, Group: 0x05, Register: 0x0003}: {Type: "f32_range", Min: 1, Max: 20, Step: 1},
	{Opcode: B524OpcodeLocal, Group: 0x05, Register: 0x0004}: {Type: "f32_range", Min: -10, Max: 110, Step: 1},

	{Opcode: B524OpcodeLocal, Group: 0x08, Register: 0x0001}: {Type: "f32_range", Min: 0, Max: 99, Step: 1},
	{Opcode: B524OpcodeLocal, Group: 0x08, Register: 0x0002}: {Type: "f32_range", Min: 0, Max: 99, Step: 1},
	{Opcode: B524OpcodeLocal, Group: 0x08, Register: 0x0003}: {Type: "f32_range", Min: 2, Max: 25, Step: 1},
	{Opcode: B524OpcodeLocal, Group: 0x08, Register: 0x0004}: {Type: "f32_range", Min: 1, Max: 20, Step: 1},
	{Opcode: B524OpcodeLocal, Group: 0x08, Register: 0x0005}: {Type: "f32_range", Min: -10, Max: 110, Step: 1},
	{Opcode: B524OpcodeLocal, Group: 0x08, Register: 0x0006}: {Type: "f32_range", Min: -10, Max: 110, Step: 1},

	{Opcode: B524OpcodeLocal, Group: 0x09, Register: 0x0001}: {Type: "u16_range", Min: 0, Max: 255, Step: 1},
	{Opcode: B524OpcodeLocal, Group: 0x09, Register: 0x0002}: {Type: "u16_range", Min: 1, Max: 3, Step: 1},
	{Opcode: B524OpcodeLocal, Group: 0x09, Register: 0x0003}: {Type: "u8_range", Min: 0, Max: 1, Step: 1},
	{Opcode: B524OpcodeLocal, Group: 0x09, Register: 0x0004}: {Type: "u16_range", Min: 0, Max: 10, Step: 1},
	{Opcode: B524OpcodeLocal, Group: 0x09, Register: 0x0005}: {Type: "u16_range", Min: 0, Max: 32768, Step: 1},
	{Opcode: B524OpcodeLocal, Group: 0x09, Register: 0x0006}: {Type: "u16_range", Min: 0, Max: 32768, Step: 1},

	{Opcode: B524OpcodeLocal, Group: 0x0A, Register: 0x0001}: {Type: "u8_range", Min: 0, Max: 3, Step: 1},
	{Opcode: B524OpcodeLocal, Group: 0x0A, Register: 0x0002}: {Type: "u8_range", Min: 1, Max: 2, Step: 1},
	{Opcode: B524OpcodeLocal, Group: 0x0A, Register: 0x0003}: {Type: "u8_range", Min: 1, Max: 2, Step: 1},
	{Opcode: B524OpcodeLocal, Group: 0x0A, Register: 0x0005}: {Type: "u8_range", Min: 0, Max: 3, Step: 1},
	{Opcode: B524OpcodeLocal, Group: 0x0A, Register: 0x0006}: {Type: "u8_range", Min: 0, Max: 1, Step: 1},
}

func lookupB524Constraint(opcode, group byte, register uint16) (b524Constraint, bool) {
	constraint, ok := b524StaticConstraintsCatalog[b524ConstraintSelector{
		Opcode:   opcode,
		Group:    group,
		Register: register,
	}]
	return constraint, ok
}
