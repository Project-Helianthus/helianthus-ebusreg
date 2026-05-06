package productids

type StaticSeedEntry struct {
	Manufacturer string
	DeviceID     string
	Addresses    []SeedAddressEntry
	Source       string
}

type SeedAddressEntry struct {
	Addr       byte
	Role       string
	Confidence string
}

const staticSeedSource = "vaillant_static_seed_w19_26"

func LoadSeedTable(enabled bool) []StaticSeedEntry {
	if !enabled {
		return nil
	}

	return []StaticSeedEntry{
		{
			Manufacturer: "Vaillant",
			DeviceID:     "NETX3",
			Source:       staticSeedSource,
			Addresses: []SeedAddressEntry{
				{Addr: 0xF1, Role: "initiator", Confidence: "candidate"},
				{Addr: 0xF6, Role: "target", Confidence: "candidate"},
				{Addr: 0x04, Role: "target", Confidence: "candidate"},
			},
		},
		{
			Manufacturer: "Vaillant",
			DeviceID:     "BASV2",
			Source:       staticSeedSource,
			Addresses: []SeedAddressEntry{
				{Addr: 0x15, Role: "target", Confidence: "candidate"},
				{Addr: 0xEC, Role: "target", Confidence: "candidate"},
			},
		},
	}
}
