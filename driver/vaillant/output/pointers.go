package output

// CloneBoolPtr copies the pointed-to value, preserving absence.
func CloneBoolPtr(value *bool) *bool {
	if value == nil {
		return nil
	}
	cp := *value
	return &cp
}

// CloneStringPtr copies the pointed-to value, preserving absence.
func CloneStringPtr(value *string) *string {
	if value == nil {
		return nil
	}
	cp := *value
	return &cp
}

// CloneIntPtr copies the pointed-to value, preserving absence.
func CloneIntPtr(value *int) *int {
	if value == nil {
		return nil
	}
	cp := *value
	return &cp
}

// CloneFloatPtr copies the pointed-to value, preserving absence.
func CloneFloatPtr(value *float64) *float64 {
	if value == nil {
		return nil
	}
	cp := *value
	return &cp
}
