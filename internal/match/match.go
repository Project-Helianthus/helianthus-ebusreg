package match

// HardwareVersionAtLeast returns true when version parses to a number >= minimum.
func HardwareVersionAtLeast(version string, minimum uint16) bool {
	value, ok := parseVersionDigits(version)
	if !ok {
		return false
	}
	return value >= minimum
}

func parseVersionDigits(value string) (uint16, bool) {
	var number uint32
	found := false
	for index := 0; index < len(value); index++ {
		character := value[index]
		if character < '0' || character > '9' {
			continue
		}
		found = true
		number = number*10 + uint32(character-'0')
		if number > 0xFFFF {
			return 0, false
		}
	}
	if !found {
		return 0, false
	}
	return uint16(number), true
}
