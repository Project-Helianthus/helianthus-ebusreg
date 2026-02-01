package schema

// SchemaCondition defines a conditional schema match.
type SchemaCondition struct {
	Target    byte
	HasTarget bool
	MinHW     uint16
	HasMinHW  bool
	MaxHW     uint16
	HasMaxHW  bool
	Schema    Schema
}

// SchemaSelector selects schemas based on target address and hardware version.
type SchemaSelector struct {
	Default    Schema
	Conditions []SchemaCondition
}

// Select returns the first matching schema or the default.
func (s SchemaSelector) Select(target byte, hwVersion string) Schema {
	hwValue, hasHW := parseHWVersion(hwVersion)
	for _, condition := range s.Conditions {
		if condition.HasTarget && condition.Target != target {
			continue
		}
		if (condition.HasMinHW || condition.HasMaxHW) && !hasHW {
			continue
		}
		if condition.HasMinHW && hwValue < condition.MinHW {
			continue
		}
		if condition.HasMaxHW && hwValue > condition.MaxHW {
			continue
		}
		return condition.Schema
	}
	return s.Default
}

func parseHWVersion(value string) (uint16, bool) {
	var number uint32
	found := false
	for i := 0; i < len(value); i++ {
		ch := value[i]
		if ch < '0' || ch > '9' {
			continue
		}
		found = true
		number = number*10 + uint32(ch-'0')
		if number > 0xFFFF {
			return 0, false
		}
	}
	if !found {
		return 0, false
	}
	return uint16(number), true
}
