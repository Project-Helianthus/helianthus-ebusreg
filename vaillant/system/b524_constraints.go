package system

type b524ConstraintSelector struct {
	Group  byte
	Record uint16
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

var b524StaticConstraintsCatalog = map[b524ConstraintSelector]b524Constraint{
	{Group: 0x00, Record: 0x0100}: {Type: "f32_range", Min: -20, Max: 50, Step: 1},
	{Group: 0x00, Record: 0x0200}: {Type: "f32_range", Min: -26, Max: 10, Step: 1},
	{Group: 0x00, Record: 0x0300}: {Type: "u16_range", Min: 0, Max: 12, Step: 1},
	{Group: 0x00, Record: 0x0400}: {Type: "u16_range", Min: 0, Max: 300, Step: 10},
	{Group: 0x00, Record: 0x8000}: {Type: "f32_range", Min: -10, Max: 10, Step: 1},

	{Group: 0x01, Record: 0x0100}: {Type: "u16_range", Min: 0, Max: 1, Step: 1},
	{Group: 0x01, Record: 0x0200}: {Type: "u8_range", Min: 0, Max: 1, Step: 1},
	{Group: 0x01, Record: 0x0300}: {Type: "u16_range", Min: 0, Max: 2, Step: 1},
	{Group: 0x01, Record: 0x0400}: {Type: "f32_range", Min: 35, Max: 70, Step: 1},
	{Group: 0x01, Record: 0x0500}: {Type: "f32_range", Min: 0, Max: 99, Step: 1},
	{Group: 0x01, Record: 0x0600}: {Type: "u8_range", Min: 0, Max: 1, Step: 1},

	{Group: 0x02, Record: 0x0100}: {Type: "u16_range", Min: 1, Max: 2, Step: 1},
	{Group: 0x02, Record: 0x0200}: {Type: "u16_range", Min: 0, Max: 4, Step: 1},
	{Group: 0x02, Record: 0x0400}: {Type: "f32_range", Min: 15, Max: 80, Step: 1},
	{Group: 0x02, Record: 0x0500}: {Type: "u8_range", Min: 0, Max: 1, Step: 1},
	{Group: 0x02, Record: 0x0600}: {Type: "u8_range", Min: 0, Max: 1, Step: 1},

	{Group: 0x03, Record: 0x0100}: {Type: "u16_range", Min: 0, Max: 2, Step: 1},
	{Group: 0x03, Record: 0x0200}: {Type: "f32_range", Min: 15, Max: 30, Step: 0.5},
	{Group: 0x03, Record: 0x0500}: {Type: "f32_range", Min: 5, Max: 30, Step: 1},
	{Group: 0x03, Record: 0x0600}: {Type: "u16_range", Min: 0, Max: 2, Step: 1},

	{Group: 0x04, Record: 0x0100}: {Type: "u8_range", Min: 0, Max: 1, Step: 1},
	{Group: 0x04, Record: 0x0200}: {Type: "u8_range", Min: 0, Max: 1, Step: 1},
	{Group: 0x04, Record: 0x0300}: {Type: "f32_range", Min: -40, Max: 155, Step: 1},
	{Group: 0x04, Record: 0x0400}: {Type: "f32_range", Min: 0, Max: 99, Step: 1},
	{Group: 0x04, Record: 0x0500}: {Type: "f32_range", Min: 110, Max: 150, Step: 1},
	{Group: 0x04, Record: 0x0600}: {Type: "f32_range", Min: 75, Max: 115, Step: 1},

	{Group: 0x05, Record: 0x0100}: {Type: "f32_range", Min: 0, Max: 99, Step: 1},
	{Group: 0x05, Record: 0x0200}: {Type: "f32_range", Min: 2, Max: 25, Step: 1},
	{Group: 0x05, Record: 0x0300}: {Type: "f32_range", Min: 1, Max: 20, Step: 1},
	{Group: 0x05, Record: 0x0400}: {Type: "f32_range", Min: -10, Max: 110, Step: 1},

	{Group: 0x08, Record: 0x0100}: {Type: "f32_range", Min: 0, Max: 99, Step: 1},
	{Group: 0x08, Record: 0x0200}: {Type: "f32_range", Min: 0, Max: 99, Step: 1},
	{Group: 0x08, Record: 0x0300}: {Type: "f32_range", Min: 2, Max: 25, Step: 1},
	{Group: 0x08, Record: 0x0400}: {Type: "f32_range", Min: 1, Max: 20, Step: 1},
	{Group: 0x08, Record: 0x0500}: {Type: "f32_range", Min: -10, Max: 110, Step: 1},
	{Group: 0x08, Record: 0x0600}: {Type: "f32_range", Min: -10, Max: 110, Step: 1},

	{Group: 0x09, Record: 0x0100}: {Type: "u16_range", Min: 0, Max: 255, Step: 1},
	{Group: 0x09, Record: 0x0200}: {Type: "u16_range", Min: 1, Max: 3, Step: 1},
	{Group: 0x09, Record: 0x0300}: {Type: "u8_range", Min: 0, Max: 1, Step: 1},
	{Group: 0x09, Record: 0x0400}: {Type: "u16_range", Min: 0, Max: 10, Step: 1},
	{Group: 0x09, Record: 0x0500}: {Type: "u16_range", Min: 0, Max: 32768, Step: 1},
	{Group: 0x09, Record: 0x0600}: {Type: "u16_range", Min: 0, Max: 32768, Step: 1},

	{Group: 0x0A, Record: 0x0100}: {Type: "u8_range", Min: 0, Max: 3, Step: 1},
	{Group: 0x0A, Record: 0x0200}: {Type: "u8_range", Min: 1, Max: 2, Step: 1},
	{Group: 0x0A, Record: 0x0300}: {Type: "u8_range", Min: 1, Max: 2, Step: 1},
	{Group: 0x0A, Record: 0x0500}: {Type: "u8_range", Min: 0, Max: 3, Step: 1},
	{Group: 0x0A, Record: 0x0600}: {Type: "u8_range", Min: 0, Max: 1, Step: 1},
}

func lookupB524Constraint(group byte, record uint16) (b524Constraint, bool) {
	constraint, ok := b524StaticConstraintsCatalog[b524ConstraintSelector{
		Group:  group,
		Record: record,
	}]
	return constraint, ok
}
