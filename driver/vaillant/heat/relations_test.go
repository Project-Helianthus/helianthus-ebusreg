package heat

import (
	"testing"

	semreg "github.com/Project-Helianthus/helianthus-semreg/semreg/v1"
)

// Relation fields, and the production mutation each assertion catches:
//   - system member direction: from the thermal system. A reversed pair names the system as the target.
//   - system member targets: zone, circuit, hot water, store, generator, refrigerant circuit and valve.
//     Adding the controller or the compressor puts an id outside that set into the batch.
//   - fluid supply direction: from the circuit named by the assignment register to the zone.
//     Using the zone instance as the circuit index links zone 1 to circuit 1.
//   - assignment raw: the evidence text includes the register value. A different encoding changes the digest.
//   - unmatched assignment, including raw 255: no fluid-supply row. Emitting one invents a circuit.
//   - relation order: definition, from, to. Swapping the zone slice changes that order only if the sort is dropped.

func TestThermalGraphRelations(t *testing.T) {
	assignment := uint16(0)
	statement := baseStatement()
	statement.Circuits = []Circuit{{Instance: 1}, {Instance: 0}}
	statement.Zones = []Zone{{Instance: 1, Assignment: &assignment, ValveReported: true}}
	batch := publishStatement(t, statement, nil)
	parent := string(testID("controller", "", "serial-alpha"))
	system := testID("system", parent, testSlot(0, 0))
	members := map[semreg.ResourceID]bool{
		testID("zone", parent, testSlot(3, 1)):                   true,
		testID("valve", parent, testSlot(3, 1)+"/valve"):         true,
		testID("circuit", parent, testSlot(2, 0)):                true,
		testID("circuit", parent, testSlot(2, 1)):                true,
		testID("store", parent, testSlot(5, 4)):                  true,
		testID("hot-water", parent, testSlot(1, 0)):              true,
		testID("generator", parent, "generator-serial"):          true,
		testID("refrigerant", parent, "refrigerant/pump-serial"): true,
	}
	fluid := map[string]bool{}
	for _, relation := range batch.RelationUpserts {
		if err := relation.Validate(); err != nil {
			t.Fatalf("relation %s: %v", relation.Definition.ID, err)
		}
		if relation.BasisClass != semreg.RelationGateRule || len(relation.Basis) != 1 || relation.Definition.Pack.ID != "helianthus.pack.thermal" || relation.Definition.Version != "1.0.0" {
			t.Fatalf("relation basis = %+v", relation)
		}
		switch relation.Definition.ID {
		case "thermal.relation.system_member":
			if relation.From != system || !members[relation.To] {
				t.Fatalf("system member %s -> %s", relation.From, relation.To)
			}
			delete(members, relation.To)
		case "thermal.relation.fluid_supply":
			fluid[string(relation.From)+"->"+string(relation.To)] = true
			if relation.Basis[0].Digest != testEvidenceRef("zone-circuit-assignment instance=1 raw=0").Digest {
				t.Fatalf("assignment evidence = %s", relation.Basis[0].Digest)
			}
		default:
			t.Fatalf("unexpected relation %s", relation.Definition.ID)
		}
	}
	if len(members) != 0 {
		t.Fatalf("missing system members %#v", members)
	}
	circuit := testID("circuit", parent, testSlot(2, 0))
	otherCircuit := testID("circuit", parent, testSlot(2, 1))
	zone := testID("zone", parent, testSlot(3, 1))
	if !fluid[string(circuit)+"->"+string(zone)] || fluid[string(otherCircuit)+"->"+string(zone)] {
		t.Fatalf("fluid supply = %#v", fluid)
	}
	controller := testID("controller", "", "serial-alpha")
	compressor := testID("compressor", parent, "compressor/pump-serial")
	for _, relation := range batch.RelationUpserts {
		if relation.To == controller || relation.To == compressor || relation.From == controller {
			t.Fatalf("relation uses %s -> %s", relation.From, relation.To)
		}
	}

	for _, raw := range []*uint16{nil, uint16ptr(0x00FF), uint16ptr(9)} {
		unknown := baseStatement()
		unknown.Zones = []Zone{{Instance: 1, Assignment: raw, ValveReported: false}}
		unknown.Circuits = []Circuit{{Instance: 0}}
		for _, relation := range publishStatement(t, unknown, nil).RelationUpserts {
			if relation.Definition.ID == "thermal.relation.fluid_supply" {
				t.Fatalf("raw %v produced fluid supply %s -> %s", raw, relation.From, relation.To)
			}
		}
	}

	swapped := statement
	swapped.Zones = []Zone{{Instance: 2, Assignment: &assignment}, {Instance: 1, Assignment: &assignment, ValveReported: true}}
	swappedBack := swapped
	swappedBack.Zones = []Zone{swapped.Zones[1], swapped.Zones[0]}
	if relationKeys(publishStatement(t, swapped, nil)) != relationKeys(publishStatement(t, swappedBack, nil)) {
		t.Fatal("zone order changed relation keys")
	}
}

func uint16ptr(value uint16) *uint16 { return &value }

func relationKeys(batch semreg.PublicationBatch) string {
	out := ""
	for _, relation := range batch.RelationUpserts {
		out += string(relation.Definition.ID) + " " + string(relation.From) + " " + string(relation.To) + "\n"
	}
	return out
}
