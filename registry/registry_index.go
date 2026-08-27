package registry

func CanonicalIndexForEntry(entry DeviceEntry) (CanonicalIndex, error) {
	if entry == nil {
		return CanonicalIndex{}, ErrProjectionInvalidNode
	}
	if internal, ok := entry.(*deviceEntry); ok {
		return internal.index, internal.indexErr
	}
	return BuildCanonicalIndex(entry.Projections())
}

func containsAddress(addresses []byte, address byte) bool {
	for _, existing := range addresses {
		if existing == address {
			return true
		}
	}
	return false
}

func removeAddress(addresses []byte, address byte) []byte {
	for index, existing := range addresses {
		if existing != address {
			continue
		}
		return append(addresses[:index], addresses[index+1:]...)
	}
	return addresses
}

func removeEntry(entries []*deviceEntry, entry *deviceEntry) []*deviceEntry {
	for index, existing := range entries {
		if existing != entry {
			continue
		}
		return append(entries[:index], entries[index+1:]...)
	}
	return entries
}
