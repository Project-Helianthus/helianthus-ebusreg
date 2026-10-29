package system

type Ab03InstanceInterval struct {
	Ab03Minimum byte
	Ab03Maximum byte
}

// Ab03DiscoveryWindow is a scheduling window, never a terminal register maximum.
// Characterization does not establish physical availability, codec qualification,
// or automatic admission; those remain Unknown. Nil ceilings remain absent.
type Ab03DiscoveryWindow struct {
	Ab03Family            string
	Ab03Opcode            byte
	Ab03Group             byte
	Ab03Instances         Ab03InstanceInterval
	Ab03RegisterCeiling   *uint16
	Ab03AvailabilityOnly  bool
	Ab03PredicateRegister *uint16
	Ab03Virtual           bool
}

func ab03Register(ab03Value uint16) *uint16 { return &ab03Value }

// Ab03DiscoveryWindows explicitly selects family 720 without modifying the
// generic exploration profiles. Local GG08 is omitted without a caller-selected
// bounded interval; its scalar ceiling is Unknown even when scope is selected.
func Ab03DiscoveryWindows(ab03Family string, ab03LocalGG08 *Ab03InstanceInterval) []Ab03DiscoveryWindow {
	if ab03Family != "720" {
		return nil
	}
	ab03Windows := []Ab03DiscoveryWindow{}
	ab03Add := func(ab03Opcode, ab03Group, ab03Minimum, ab03Maximum byte, ab03Ceiling *uint16) {
		ab03Windows = append(ab03Windows, Ab03DiscoveryWindow{
			Ab03Family: ab03Family, Ab03Opcode: ab03Opcode, Ab03Group: ab03Group,
			Ab03Instances: Ab03InstanceInterval{ab03Minimum, ab03Maximum}, Ab03RegisterCeiling: ab03Ceiling,
		})
	}
	ab03Add(B524OpcodeLocal, 0x02, 0, 8, ab03Register(0x25))
	ab03Add(B524OpcodeLocal, 0x02, 9, 9, ab03Register(0x25))
	ab03Windows[1].Ab03Virtual = true
	ab03Add(B524OpcodeLocal, 0x03, 0, 10, ab03Register(0x2e))
	ab03Add(B524OpcodeLocal, 0x04, 0, 1, ab03Register(0x0b))
	ab03Add(B524OpcodeLocal, 0x05, 0, 1, ab03Register(0x04))
	if ab03LocalGG08 != nil && ab03LocalGG08.Ab03Minimum <= ab03LocalGG08.Ab03Maximum {
		ab03Add(B524OpcodeLocal, 0x08, ab03LocalGG08.Ab03Minimum, ab03LocalGG08.Ab03Maximum, nil)
	}
	for _, ab03Group := range []byte{1, 2, 3, 5, 6, 7, 8, 0x0b, 0x0c} {
		ab03Add(B524OpcodeRemote, ab03Group, 1, 8, ab03Register(0x2f))
	}
	for _, ab03Group := range []byte{9, 10} {
		ab03Add(B524OpcodeRemote, ab03Group, 1, 8, ab03Register(0x35))
	}
	for _, ab03Group := range []byte{0x0e, 0x0f} {
		ab03Add(B524OpcodeRemote, ab03Group, 1, 8, ab03Register(0x33))
	}
	ab03Add(B524OpcodeRemote, 0x04, 1, 8, nil)
	ab03Add(B524OpcodeRemote, 0x0d, 1, 8, nil)
	ab03Windows[len(ab03Windows)-1].Ab03AvailabilityOnly = true
	ab03Windows[len(ab03Windows)-1].Ab03PredicateRegister = ab03Register(1)
	return ab03Windows
}

// Ab03LookupDiscoveryWindow matches a characterized scalar schedule or the
// separate availability predicate. Absent scalar ceilings admit no scalar RR.
// A match establishes neither presence nor authorization to transmit.
func Ab03LookupDiscoveryWindow(ab03Family string, ab03Opcode, ab03Group, ab03Instance byte, ab03RegisterID uint16, ab03LocalGG08 *Ab03InstanceInterval) (Ab03DiscoveryWindow, bool) {
	for _, ab03Window := range Ab03DiscoveryWindows(ab03Family, ab03LocalGG08) {
		if ab03Window.Ab03Opcode != ab03Opcode || ab03Window.Ab03Group != ab03Group || ab03Instance < ab03Window.Ab03Instances.Ab03Minimum || ab03Instance > ab03Window.Ab03Instances.Ab03Maximum {
			continue
		}
		if ab03Window.Ab03AvailabilityOnly {
			if ab03Window.Ab03PredicateRegister != nil && ab03RegisterID == *ab03Window.Ab03PredicateRegister {
				return ab03Window, true
			}
		} else if ab03Window.Ab03RegisterCeiling != nil && ab03RegisterID <= *ab03Window.Ab03RegisterCeiling {
			return ab03Window, true
		}
	}
	return Ab03DiscoveryWindow{}, false
}

// Ab03PresentDiscoveryInstances retains only explicit per-instance presence.
// OP00 capacity IDs cannot allocate slots or fill gaps in sparse presence.
func Ab03PresentDiscoveryInstances(ab03Family string, ab03Opcode, ab03Group byte, ab03RegisterID uint16, ab03Presence map[byte]bool) []byte {
	var ab03Instances []byte
	for ab03Index := 0; ab03Index <= 255; ab03Index++ {
		ab03Instance := byte(ab03Index)
		if !ab03Presence[ab03Instance] {
			continue
		}
		if _, ab03OK := Ab03LookupDiscoveryWindow(ab03Family, ab03Opcode, ab03Group, ab03Instance, ab03RegisterID, nil); ab03OK {
			ab03Instances = append(ab03Instances, ab03Instance)
		}
	}
	return ab03Instances
}
