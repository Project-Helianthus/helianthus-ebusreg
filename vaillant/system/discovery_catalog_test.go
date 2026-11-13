package system

import (
	discoveryCatalogJSON "encoding/json"
	discoveryCatalogOs "os"
	discoveryCatalogReflect "reflect"
	discoveryCatalogTesting "testing"
)

func TestAb03ExactFixture(discoveryCatalogT *discoveryCatalogTesting.T) {
	discoveryCatalogData, discoveryCatalogErr := discoveryCatalogOs.ReadFile("testdata/ab03_discovery_windows.json")
	if discoveryCatalogErr != nil {
		discoveryCatalogT.Fatal(discoveryCatalogErr)
	}
	var discoveryCatalogFixture struct {
		DiscoveryCatalogWindows []DiscoveryWindow `json:"Ab03Windows"`
	}
	if discoveryCatalogErr = discoveryCatalogJSON.Unmarshal(discoveryCatalogData, &discoveryCatalogFixture); discoveryCatalogErr != nil {
		discoveryCatalogT.Fatal(discoveryCatalogErr)
	}
	if len(discoveryCatalogFixture.DiscoveryCatalogWindows) != 20 {
		discoveryCatalogT.Fatalf("fixture windows = %d", len(discoveryCatalogFixture.DiscoveryCatalogWindows))
	}
	if discoveryCatalogGot := DiscoveryWindows("720", nil); !discoveryCatalogReflect.DeepEqual(discoveryCatalogGot, discoveryCatalogFixture.DiscoveryCatalogWindows) {
		discoveryCatalogT.Fatalf("catalog differs from exact fixture: %+v", discoveryCatalogGot)
	}
	for _, discoveryCatalogWindow := range discoveryCatalogFixture.DiscoveryCatalogWindows {
		if discoveryCatalogWindow.DiscoveryCatalogRegisterCeiling == nil {
			continue
		}
		for _, discoveryCatalogInstance := range []byte{discoveryCatalogWindow.DiscoveryCatalogInstances.DiscoveryCatalogMinimum, discoveryCatalogWindow.DiscoveryCatalogInstances.DiscoveryCatalogMaximum} {
			discoveryCatalogGot, discoveryCatalogOK := LookupDiscoveryWindow(discoveryCatalogWindow.DiscoveryCatalogFamily, discoveryCatalogWindow.DiscoveryCatalogOpcode, discoveryCatalogWindow.DiscoveryCatalogGroup, discoveryCatalogInstance, *discoveryCatalogWindow.DiscoveryCatalogRegisterCeiling, nil)
			if !discoveryCatalogOK || !discoveryCatalogReflect.DeepEqual(discoveryCatalogGot, discoveryCatalogWindow) {
				discoveryCatalogT.Fatalf("boundary lookup differs: %+v", discoveryCatalogWindow)
			}
		}
	}
}

func TestAb03IndependentLookupDimensions(discoveryCatalogT *discoveryCatalogTesting.T) {
	for _, discoveryCatalogCase := range []struct {
		discoveryCatalogName                                       string
		discoveryCatalogFamily                                     string
		discoveryCatalogOP, discoveryCatalogGG, discoveryCatalogII byte
		discoveryCatalogRR                                         uint16
		discoveryCatalogWant                                       bool
	}{
		{"ab03_control", "720", 2, 3, 1, 0x2e, true},
		{"ab03_family", "unknown", 2, 3, 1, 0x2e, false},
		{"ab03_op_remote", "720", 6, 3, 1, 0x2f, true},
		{"ab03_op_local_collision", "720", 2, 3, 1, 0x2f, false},
		{"ab03_op_unknown", "720", 3, 3, 1, 0x2e, false},
		{"ab03_gg", "720", 2, 4, 1, 0x2e, false},
		{"ab03_min_control", "720", 6, 3, 1, 0x2e, true},
		{"ab03_min", "720", 6, 3, 0, 0x2e, false},
		{"ab03_max_control", "720", 6, 3, 8, 0x2e, true},
		{"ab03_max", "720", 6, 3, 9, 0x2e, false},
		{"ab03_rr", "720", 2, 3, 1, 0x2f, false},
		{"ab03_local_control", "720", 2, 3, 0, 0x2e, true},
		{"ab03_remote_collision", "720", 6, 3, 0, 0x2e, false},
		{"ab03_availability", "720", 6, 13, 1, 1, true},
		{"ab03_availability_zero", "720", 6, 13, 1, 0, false},
		{"ab03_availability_two", "720", 6, 13, 1, 2, false},
		{"ab03_absent_ceiling_zero", "720", 6, 4, 1, 0, false},
		{"ab03_absent_ceiling_ffff", "720", 6, 4, 1, 0xffff, false},
		{"ab03_virtual", "720", 2, 2, 9, 37, true},
		{"ab03_after_virtual", "720", 2, 2, 10, 37, false},
	} {
		discoveryCatalogT.Run(discoveryCatalogCase.discoveryCatalogName, func(discoveryCatalogT *discoveryCatalogTesting.T) {
			_, discoveryCatalogGot := LookupDiscoveryWindow(discoveryCatalogCase.discoveryCatalogFamily, discoveryCatalogCase.discoveryCatalogOP, discoveryCatalogCase.discoveryCatalogGG, discoveryCatalogCase.discoveryCatalogII, discoveryCatalogCase.discoveryCatalogRR, nil)
			if discoveryCatalogGot != discoveryCatalogCase.discoveryCatalogWant {
				discoveryCatalogT.Fatalf("match = %v, want %v", discoveryCatalogGot, discoveryCatalogCase.discoveryCatalogWant)
			}
		})
	}
}

func TestAb03ExplicitGG08AndIsolation(discoveryCatalogT *discoveryCatalogTesting.T) {
	for _, discoveryCatalogScope := range []DiscoveryCatalogInstanceInterval{{2, 5}, {3, 5}, {2, 6}, {0, 0}, {255, 255}} {
		discoveryCatalogFound := false
		for _, discoveryCatalogWindow := range DiscoveryWindows("720", &discoveryCatalogScope) {
			if discoveryCatalogWindow.DiscoveryCatalogOpcode == 2 && discoveryCatalogWindow.DiscoveryCatalogGroup == 8 {
				discoveryCatalogFound = true
				if discoveryCatalogWindow.DiscoveryCatalogInstances != discoveryCatalogScope || discoveryCatalogWindow.DiscoveryCatalogRegisterCeiling != nil {
					discoveryCatalogT.Fatal("GG08 borrowed scope or ceiling")
				}
			}
		}
		if !discoveryCatalogFound {
			discoveryCatalogT.Fatal("missing explicit GG08")
		}
	}
	for _, discoveryCatalogScope := range []*DiscoveryCatalogInstanceInterval{nil, {5, 2}} {
		for _, discoveryCatalogWindow := range DiscoveryWindows("720", discoveryCatalogScope) {
			if discoveryCatalogWindow.DiscoveryCatalogOpcode == 2 && discoveryCatalogWindow.DiscoveryCatalogGroup == 8 {
				discoveryCatalogT.Fatal("implicit or invalid GG08 scope")
			}
		}
	}
	for _, discoveryCatalogFamily := range []string{"", "unknown", "700"} {
		if len(DiscoveryWindows(discoveryCatalogFamily, &DiscoveryCatalogInstanceInterval{2, 5})) != 0 {
			discoveryCatalogT.Fatal("unknown family fallback")
		}
	}
	discoveryCatalogWindows := DiscoveryWindows("720", nil)
	*discoveryCatalogWindows[0].DiscoveryCatalogRegisterCeiling = 0
	*discoveryCatalogWindows[len(discoveryCatalogWindows)-1].DiscoveryCatalogPredicateRegister = 0
	if _, discoveryCatalogOK := LookupDiscoveryWindow("720", 2, 2, 1, 37, nil); !discoveryCatalogOK {
		discoveryCatalogT.Fatal("caller mutated catalog")
	}
	if _, discoveryCatalogOK := LookupDiscoveryWindow("720", 6, 13, 1, 1, nil); !discoveryCatalogOK {
		discoveryCatalogT.Fatal("caller mutated predicate")
	}
}

func TestAb03SparsePresenceAndCapacity(discoveryCatalogT *discoveryCatalogTesting.T) {
	discoveryCatalogPresence := map[byte]bool{1: true, 8: true, 2: false}
	if discoveryCatalogGot := PresentDiscoveryInstances("720", 6, 3, 1, discoveryCatalogPresence); !discoveryCatalogReflect.DeepEqual(discoveryCatalogGot, []byte{1, 8}) {
		discoveryCatalogT.Fatalf("sparse presence replaced by count: %v", discoveryCatalogGot)
	}
	for _, discoveryCatalogID := range []uint16{0, 1} {
		if discoveryCatalogGot := PresentDiscoveryInstances("720", 0, 3, discoveryCatalogID, discoveryCatalogPresence); len(discoveryCatalogGot) != 0 {
			discoveryCatalogT.Fatalf("capacity allocated instances: %v", discoveryCatalogGot)
		}
	}
	for _, discoveryCatalogCase := range []struct {
		discoveryCatalogFamily                 string
		discoveryCatalogOP, discoveryCatalogGG byte
		discoveryCatalogRR                     uint16
	}{
		{"unknown", 6, 3, 1},
		{"720", 0, 3, 1},
		{"720", 6, 4, 1},
		{"720", 6, 3, 0x30},
	} {
		if discoveryCatalogGot := PresentDiscoveryInstances(discoveryCatalogCase.discoveryCatalogFamily, discoveryCatalogCase.discoveryCatalogOP, discoveryCatalogCase.discoveryCatalogGG, discoveryCatalogCase.discoveryCatalogRR, discoveryCatalogPresence); len(discoveryCatalogGot) != 0 {
			discoveryCatalogT.Fatalf("presence bypassed selector scope: %v", discoveryCatalogGot)
		}
	}
	if discoveryCatalogGot := PresentDiscoveryInstances("720", 6, 3, 1, nil); len(discoveryCatalogGot) != 0 {
		discoveryCatalogT.Fatal("absent presence allocated instances")
	}
	if !discoveryCatalogReflect.DeepEqual(discoveryCatalogPresence, map[byte]bool{1: true, 8: true, 2: false}) {
		discoveryCatalogT.Fatal("presence mutated")
	}
}

func TestAb03GenericProfilesUnchanged(discoveryCatalogT *discoveryCatalogTesting.T) {
	discoveryCatalogWant := []B524GroupProfile{
		{0, 2, 0, 255}, {1, 2, 0, 19}, {2, 2, 10, 37}, {3, 2, 10, 47},
		{8, 2, 0, 7}, {8, 6, 10, 4}, {9, 2, 10, 15}, {9, 6, 10, 47},
		{10, 2, 10, 77}, {10, 6, 10, 63}, {12, 6, 10, 63},
	}
	DiscoveryWindows("720", &DiscoveryCatalogInstanceInterval{2, 5})
	if !discoveryCatalogReflect.DeepEqual(B524DiscoveryProfiles, discoveryCatalogWant) {
		discoveryCatalogT.Fatal("generic exploration defaults changed")
	}
}
