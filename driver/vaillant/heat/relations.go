package heat

import (
	"bytes"
	"fmt"
	"sort"

	semreg "github.com/Project-Helianthus/helianthus-semreg/semreg/v1"
	"github.com/Project-Helianthus/helianthus-semreg/semreg/v1/packs/thermal"
)

// thermalPack is the thermal pack whose relations this publisher can name.
// A core-owned relation has no pack definition, so it is not staged.
var thermalPack = thermal.New().Pack()

func buildRelations(statement Statement, index graphIndex, stamp originStamp) ([]semreg.RelationInstance, error) {
	var relations []semreg.RelationInstance
	seen := map[string]struct{}{}
	add := func(id string, from, to semreg.ResourceID, evidence string) {
		key := id + "\x00" + string(from) + "\x00" + string(to)
		if _, ok := seen[key]; ok {
			return
		}
		seen[key] = struct{}{}
		relations = append(relations, relation(id, from, to, evidence, stamp))
	}
	for _, member := range index.members {
		add("thermal.relation.system_member", index.system, member.id, "system-member "+member.role+" "+member.slot)
	}
	for _, zone := range statement.Zones {
		if zone.Assignment == nil {
			continue
		}
		raw := *zone.Assignment
		if raw == 0x00FF || raw > 0xFF {
			continue
		}
		circuitID, circuitOK := index.circuits[byte(raw)]
		zoneID, zoneOK := index.zones[zone.Instance]
		if !circuitOK || !zoneOK {
			continue
		}
		add("thermal.relation.fluid_supply", circuitID, zoneID, fmt.Sprintf("zone-circuit-assignment instance=%d raw=%d", zone.Instance, raw))
	}
	if err := sortRelations(relations); err != nil {
		return nil, err
	}
	if len(relations) == 0 {
		return nil, nil
	}
	return relations, nil
}

func relation(id string, from, to semreg.ResourceID, evidence string, stamp originStamp) semreg.RelationInstance {
	ref := evidenceRef(evidence)
	return semreg.RelationInstance{
		Definition: semreg.DefinitionRef{Pack: thermalPack, ID: semreg.DefinitionID(id), Version: thermalPack.Version},
		From:       from,
		To:         to,
		State:      semreg.LinkQualified,
		BasisClass: semreg.RelationGateRule,
		Basis:      []semreg.EvidenceRef{ref},
		Origin: semreg.OriginRef{
			OriginID:      semreg.OriginID("origin." + digestParts(id, string(from), string(to))),
			Kind:          semreg.OriginNativeObservation,
			SourceID:      stamp.source,
			SourceEpochID: stamp.epoch,
			BindingID:     stamp.binding,
			Evidence:      []semreg.EvidenceRef{ref},
		},
		Revision: "1",
	}
}

func sortRelations(relations []semreg.RelationInstance) error {
	var failed error
	sort.SliceStable(relations, func(i, j int) bool {
		left, leftErr := semreg.CanonicalJSON(relations[i].Definition)
		right, rightErr := semreg.CanonicalJSON(relations[j].Definition)
		if leftErr != nil || rightErr != nil {
			if failed == nil {
				failed = leftErr
			}
			if failed == nil {
				failed = rightErr
			}
			return i < j
		}
		if cmp := bytes.Compare(left, right); cmp != 0 {
			return cmp < 0
		}
		if relations[i].From != relations[j].From {
			return relations[i].From < relations[j].From
		}
		return relations[i].To < relations[j].To
	})
	return failed
}
