package system

import (
	ab03json "encoding/json"
	ab03os "os"
	ab03reflect "reflect"
	ab03testing "testing"
)

func TestAb03ExactFixture(ab03T *ab03testing.T) {
	ab03Data, ab03Err := ab03os.ReadFile("testdata/ab03_discovery_windows.json")
	if ab03Err != nil {
		ab03T.Fatal(ab03Err)
	}
	var ab03Fixture struct{ Ab03Windows []Ab03DiscoveryWindow }
	if ab03Err = ab03json.Unmarshal(ab03Data, &ab03Fixture); ab03Err != nil {
		ab03T.Fatal(ab03Err)
	}
	if len(ab03Fixture.Ab03Windows) != 20 {
		ab03T.Fatalf("fixture windows = %d", len(ab03Fixture.Ab03Windows))
	}
	if ab03Got := Ab03DiscoveryWindows("720", nil); !ab03reflect.DeepEqual(ab03Got, ab03Fixture.Ab03Windows) {
		ab03T.Fatalf("catalog differs from exact fixture: %+v", ab03Got)
	}
	for _, ab03Window := range ab03Fixture.Ab03Windows {
		if ab03Window.Ab03RegisterCeiling == nil {
			continue
		}
		for _, ab03Instance := range []byte{ab03Window.Ab03Instances.Ab03Minimum, ab03Window.Ab03Instances.Ab03Maximum} {
			ab03Got, ab03OK := Ab03LookupDiscoveryWindow(ab03Window.Ab03Family, ab03Window.Ab03Opcode, ab03Window.Ab03Group, ab03Instance, *ab03Window.Ab03RegisterCeiling, nil)
			if !ab03OK || !ab03reflect.DeepEqual(ab03Got, ab03Window) {
				ab03T.Fatalf("boundary lookup differs: %+v", ab03Window)
			}
		}
	}
}

func TestAb03IndependentLookupDimensions(ab03T *ab03testing.T) {
	for _, ab03Case := range []struct {
		ab03Name               string
		ab03Family             string
		ab03OP, ab03GG, ab03II byte
		ab03RR                 uint16
		ab03Want               bool
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
		ab03T.Run(ab03Case.ab03Name, func(ab03T *ab03testing.T) {
			_, ab03Got := Ab03LookupDiscoveryWindow(ab03Case.ab03Family, ab03Case.ab03OP, ab03Case.ab03GG, ab03Case.ab03II, ab03Case.ab03RR, nil)
			if ab03Got != ab03Case.ab03Want {
				ab03T.Fatalf("match = %v, want %v", ab03Got, ab03Case.ab03Want)
			}
		})
	}
}

func TestAb03ExplicitGG08AndIsolation(ab03T *ab03testing.T) {
	for _, ab03Scope := range []Ab03InstanceInterval{{2, 5}, {3, 5}, {2, 6}, {0, 0}, {255, 255}} {
		ab03Found := false
		for _, ab03Window := range Ab03DiscoveryWindows("720", &ab03Scope) {
			if ab03Window.Ab03Opcode == 2 && ab03Window.Ab03Group == 8 {
				ab03Found = true
				if ab03Window.Ab03Instances != ab03Scope || ab03Window.Ab03RegisterCeiling != nil {
					ab03T.Fatal("GG08 borrowed scope or ceiling")
				}
			}
		}
		if !ab03Found {
			ab03T.Fatal("missing explicit GG08")
		}
	}
	for _, ab03Scope := range []*Ab03InstanceInterval{nil, {5, 2}} {
		for _, ab03Window := range Ab03DiscoveryWindows("720", ab03Scope) {
			if ab03Window.Ab03Opcode == 2 && ab03Window.Ab03Group == 8 {
				ab03T.Fatal("implicit or invalid GG08 scope")
			}
		}
	}
	for _, ab03Family := range []string{"", "unknown", "700"} {
		if len(Ab03DiscoveryWindows(ab03Family, &Ab03InstanceInterval{2, 5})) != 0 {
			ab03T.Fatal("unknown family fallback")
		}
	}
	ab03Windows := Ab03DiscoveryWindows("720", nil)
	*ab03Windows[0].Ab03RegisterCeiling = 0
	*ab03Windows[len(ab03Windows)-1].Ab03PredicateRegister = 0
	if _, ab03OK := Ab03LookupDiscoveryWindow("720", 2, 2, 1, 37, nil); !ab03OK {
		ab03T.Fatal("caller mutated catalog")
	}
	if _, ab03OK := Ab03LookupDiscoveryWindow("720", 6, 13, 1, 1, nil); !ab03OK {
		ab03T.Fatal("caller mutated predicate")
	}
}

func TestAb03SparsePresenceAndCapacity(ab03T *ab03testing.T) {
	ab03Presence := map[byte]bool{1: true, 8: true, 2: false}
	if ab03Got := Ab03PresentDiscoveryInstances("720", 6, 3, 1, ab03Presence); !ab03reflect.DeepEqual(ab03Got, []byte{1, 8}) {
		ab03T.Fatalf("sparse presence replaced by count: %v", ab03Got)
	}
	for _, ab03ID := range []uint16{0, 1} {
		if ab03Got := Ab03PresentDiscoveryInstances("720", 0, 3, ab03ID, ab03Presence); len(ab03Got) != 0 {
			ab03T.Fatalf("capacity allocated instances: %v", ab03Got)
		}
	}
	for _, ab03Case := range []struct {
		ab03Family     string
		ab03OP, ab03GG byte
		ab03RR         uint16
	}{
		{"unknown", 6, 3, 1},
		{"720", 0, 3, 1},
		{"720", 6, 4, 1},
		{"720", 6, 3, 0x30},
	} {
		if ab03Got := Ab03PresentDiscoveryInstances(ab03Case.ab03Family, ab03Case.ab03OP, ab03Case.ab03GG, ab03Case.ab03RR, ab03Presence); len(ab03Got) != 0 {
			ab03T.Fatalf("presence bypassed selector scope: %v", ab03Got)
		}
	}
	if ab03Got := Ab03PresentDiscoveryInstances("720", 6, 3, 1, nil); len(ab03Got) != 0 {
		ab03T.Fatal("absent presence allocated instances")
	}
	if !ab03reflect.DeepEqual(ab03Presence, map[byte]bool{1: true, 8: true, 2: false}) {
		ab03T.Fatal("presence mutated")
	}
}

func TestAb03GenericProfilesUnchanged(ab03T *ab03testing.T) {
	ab03Want := []B524GroupProfile{
		{0, 2, 0, 255}, {1, 2, 0, 19}, {2, 2, 10, 37}, {3, 2, 10, 47},
		{8, 2, 0, 7}, {8, 6, 10, 4}, {9, 2, 10, 15}, {9, 6, 10, 47},
		{10, 2, 10, 77}, {10, 6, 10, 63}, {12, 6, 10, 63},
	}
	Ab03DiscoveryWindows("720", &Ab03InstanceInterval{2, 5})
	if !ab03reflect.DeepEqual(B524DiscoveryProfiles, ab03Want) {
		ab03T.Fatal("generic exploration defaults changed")
	}
}
