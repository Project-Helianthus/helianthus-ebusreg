package system

type DiscoveryCatalogInstanceInterval struct {
	DiscoveryCatalogMinimum byte `json:"Ab03Minimum"`
	DiscoveryCatalogMaximum byte `json:"Ab03Maximum"`
}

// DiscoveryWindow is a scheduling window, never a terminal register maximum.
// Characterization does not establish physical availability, codec qualification,
// or automatic admission; those remain Unknown. Nil ceilings remain absent.
type DiscoveryWindow struct {
	DiscoveryCatalogFamily            string                           `json:"Ab03Family"`
	DiscoveryCatalogOpcode            byte                             `json:"Ab03Opcode"`
	DiscoveryCatalogGroup             byte                             `json:"Ab03Group"`
	DiscoveryCatalogInstances         DiscoveryCatalogInstanceInterval `json:"Ab03Instances"`
	DiscoveryCatalogRegisterCeiling   *uint16                          `json:"Ab03RegisterCeiling"`
	DiscoveryCatalogAvailabilityOnly  bool                             `json:"Ab03AvailabilityOnly"`
	DiscoveryCatalogPredicateRegister *uint16                          `json:"Ab03PredicateRegister"`
	DiscoveryCatalogVirtual           bool                             `json:"Ab03Virtual"`
}

func discoveryCatalogRegister(discoveryCatalogValue uint16) *uint16 { return &discoveryCatalogValue }

// DiscoveryWindows explicitly selects family 720 without modifying the
// generic exploration profiles. Local GG08 is omitted without a caller-selected
// bounded interval; its scalar ceiling is Unknown even when scope is selected.
func DiscoveryWindows(discoveryCatalogFamily string, discoveryCatalogLocalGG08 *DiscoveryCatalogInstanceInterval) []DiscoveryWindow {
	if discoveryCatalogFamily != "720" {
		return nil
	}
	discoveryCatalogWindows := []DiscoveryWindow{}
	discoveryCatalogAdd := func(discoveryCatalogOpcode, discoveryCatalogGroup, discoveryCatalogMinimum, discoveryCatalogMaximum byte, discoveryCatalogCeiling *uint16) {
		discoveryCatalogWindows = append(discoveryCatalogWindows, DiscoveryWindow{
			DiscoveryCatalogFamily: discoveryCatalogFamily, DiscoveryCatalogOpcode: discoveryCatalogOpcode, DiscoveryCatalogGroup: discoveryCatalogGroup,
			DiscoveryCatalogInstances: DiscoveryCatalogInstanceInterval{discoveryCatalogMinimum, discoveryCatalogMaximum}, DiscoveryCatalogRegisterCeiling: discoveryCatalogCeiling,
		})
	}
	discoveryCatalogAdd(B524OpcodeLocal, 0x02, 0, 8, discoveryCatalogRegister(0x25))
	discoveryCatalogAdd(B524OpcodeLocal, 0x02, 9, 9, discoveryCatalogRegister(0x25))
	discoveryCatalogWindows[1].DiscoveryCatalogVirtual = true
	discoveryCatalogAdd(B524OpcodeLocal, 0x03, 0, 10, discoveryCatalogRegister(0x2e))
	discoveryCatalogAdd(B524OpcodeLocal, 0x04, 0, 1, discoveryCatalogRegister(0x0b))
	discoveryCatalogAdd(B524OpcodeLocal, 0x05, 0, 1, discoveryCatalogRegister(0x04))
	if discoveryCatalogLocalGG08 != nil && discoveryCatalogLocalGG08.DiscoveryCatalogMinimum <= discoveryCatalogLocalGG08.DiscoveryCatalogMaximum {
		discoveryCatalogAdd(B524OpcodeLocal, 0x08, discoveryCatalogLocalGG08.DiscoveryCatalogMinimum, discoveryCatalogLocalGG08.DiscoveryCatalogMaximum, nil)
	}
	for _, discoveryCatalogGroup := range []byte{1, 2, 3, 5, 6, 7, 8, 0x0b, 0x0c} {
		discoveryCatalogAdd(B524OpcodeRemote, discoveryCatalogGroup, 1, 8, discoveryCatalogRegister(0x2f))
	}
	for _, discoveryCatalogGroup := range []byte{9, 10} {
		discoveryCatalogAdd(B524OpcodeRemote, discoveryCatalogGroup, 1, 8, discoveryCatalogRegister(0x35))
	}
	for _, discoveryCatalogGroup := range []byte{0x0e, 0x0f} {
		discoveryCatalogAdd(B524OpcodeRemote, discoveryCatalogGroup, 1, 8, discoveryCatalogRegister(0x33))
	}
	discoveryCatalogAdd(B524OpcodeRemote, 0x04, 1, 8, nil)
	discoveryCatalogAdd(B524OpcodeRemote, 0x0d, 1, 8, nil)
	discoveryCatalogWindows[len(discoveryCatalogWindows)-1].DiscoveryCatalogAvailabilityOnly = true
	discoveryCatalogWindows[len(discoveryCatalogWindows)-1].DiscoveryCatalogPredicateRegister = discoveryCatalogRegister(1)
	return discoveryCatalogWindows
}

// LookupDiscoveryWindow matches a characterized scalar schedule or the
// separate availability predicate. Absent scalar ceilings admit no scalar RR.
// A match establishes neither presence nor authorization to transmit.
func LookupDiscoveryWindow(discoveryCatalogFamily string, discoveryCatalogOpcode, discoveryCatalogGroup, discoveryCatalogInstance byte, discoveryCatalogRegisterID uint16, discoveryCatalogLocalGG08 *DiscoveryCatalogInstanceInterval) (DiscoveryWindow, bool) {
	for _, discoveryCatalogWindow := range DiscoveryWindows(discoveryCatalogFamily, discoveryCatalogLocalGG08) {
		if discoveryCatalogWindow.DiscoveryCatalogOpcode != discoveryCatalogOpcode || discoveryCatalogWindow.DiscoveryCatalogGroup != discoveryCatalogGroup || discoveryCatalogInstance < discoveryCatalogWindow.DiscoveryCatalogInstances.DiscoveryCatalogMinimum || discoveryCatalogInstance > discoveryCatalogWindow.DiscoveryCatalogInstances.DiscoveryCatalogMaximum {
			continue
		}
		if discoveryCatalogWindow.DiscoveryCatalogAvailabilityOnly {
			if discoveryCatalogWindow.DiscoveryCatalogPredicateRegister != nil && discoveryCatalogRegisterID == *discoveryCatalogWindow.DiscoveryCatalogPredicateRegister {
				return discoveryCatalogWindow, true
			}
		} else if discoveryCatalogWindow.DiscoveryCatalogRegisterCeiling != nil && discoveryCatalogRegisterID <= *discoveryCatalogWindow.DiscoveryCatalogRegisterCeiling {
			return discoveryCatalogWindow, true
		}
	}
	return DiscoveryWindow{}, false
}

// PresentDiscoveryInstances retains only explicit per-instance presence.
// OP00 capacity IDs cannot allocate slots or fill gaps in sparse presence.
func PresentDiscoveryInstances(discoveryCatalogFamily string, discoveryCatalogOpcode, discoveryCatalogGroup byte, discoveryCatalogRegisterID uint16, discoveryCatalogPresence map[byte]bool) []byte {
	var discoveryCatalogInstances []byte
	for discoveryCatalogIndex := 0; discoveryCatalogIndex <= 255; discoveryCatalogIndex++ {
		discoveryCatalogInstance := byte(discoveryCatalogIndex)
		if !discoveryCatalogPresence[discoveryCatalogInstance] {
			continue
		}
		if _, discoveryCatalogOK := LookupDiscoveryWindow(discoveryCatalogFamily, discoveryCatalogOpcode, discoveryCatalogGroup, discoveryCatalogInstance, discoveryCatalogRegisterID, nil); discoveryCatalogOK {
			discoveryCatalogInstances = append(discoveryCatalogInstances, discoveryCatalogInstance)
		}
	}
	return discoveryCatalogInstances
}
