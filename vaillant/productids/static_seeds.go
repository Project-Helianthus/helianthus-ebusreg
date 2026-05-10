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
				// NETX3 exposes four faces on the wire (operator
				// observation 2026-05-10): 0xF1 (master / initiator),
				// 0xF6 (slave / target companion of 0xF1), and the
				// pair 0x04 / 0xFF (target / broadcast faces). The
				// pair 0x04+0xFF must be seeded together so the
				// registry's identity-merge collapses all four
				// addresses into one DeviceEntry at boot — without
				// 0xFF in the seed, passive observations of NETX3's
				// 0xFF broadcasts land in a separate unidentified
				// entry that never gets enriched (NETX3 does not
				// respond to active identify probes on 0xFF).
				{Addr: 0xF1, Role: "initiator", Confidence: "candidate"},
				{Addr: 0xF6, Role: "target", Confidence: "candidate"},
				{Addr: 0x04, Role: "target", Confidence: "candidate"},
				{Addr: 0xFF, Role: "target", Confidence: "candidate"},
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
