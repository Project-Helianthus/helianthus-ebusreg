package heat

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"unicode"
	"unicode/utf8"

	semreg "github.com/Project-Helianthus/helianthus-semreg/semreg/v1"
	"golang.org/x/text/unicode/norm"
)

// Statement is the native thermal inventory one publisher owns. Positions are
// the documented operation, group and instance of each record. A serial or an
// operator label anchors the root. A bus address is not an identity.
type Statement struct {
	ControllerPresent bool
	ControllerSerial  string
	OperatorLabel     string
	GeneratorPresent  bool
	GeneratorSerial   string
	HotWater          bool
	Zones             []Zone
	Circuits          []Circuit
	Cylinders         []Cylinder
	HeatPumps         []Pump
}

// Zone is one present zone. Assignment is the circuit-assignment register.
// A nil assignment, or a value that names no active circuit, leaves the
// fluid supply unknown.
type Zone struct {
	Instance      byte
	Assignment    *uint16
	ValveReported bool
}

// Circuit is one active circuit.
type Circuit struct {
	Instance byte
}

// Cylinder is one reported store. Its presence is not a generator assignment.
type Cylinder struct {
	Instance byte
}

// Pump is one heat-pump identity carried by its transmitted serial.
type Pump struct {
	Serial string
}

type placed struct {
	role     string
	kind     semreg.DefinitionID
	class    semreg.AnchorClass
	basis    semreg.AnchorBasis
	parent   string
	slot     string
	evidence string
	operator bool
}

type memberRef struct {
	id   semreg.ResourceID
	role string
	slot string
}

type graphIndex struct {
	root         semreg.ResourceID
	rootEvidence string
	system       semreg.ResourceID
	zones        map[byte]semreg.ResourceID
	circuits     map[byte]semreg.ResourceID
	members      []memberRef
}

type originStamp struct {
	source  *semreg.SourceID
	epoch   *semreg.SourceEpochID
	binding *semreg.NativeBindingID
}

func slotOK(slot string) bool {
	if len(slot) < 1 || len(slot) > 128 || !utf8.ValidString(slot) || !norm.NFC.IsNormalString(slot) {
		return false
	}
	for _, r := range slot {
		if unicode.IsControl(r) {
			return false
		}
	}
	return true
}

func groupSlot(group, instance byte) string {
	return fmt.Sprintf("%02x/%02x/%02x", byte(2), group, instance)
}

func digestParts(parts ...string) string {
	hash := sha256.New()
	for i, part := range parts {
		if i > 0 {
			_, _ = hash.Write([]byte{0})
		}
		_, _ = hash.Write([]byte(part))
	}
	return hex.EncodeToString(hash.Sum(nil))
}

func resourceID(role, parent, slot string) semreg.ResourceID {
	return semreg.ResourceID("thermal." + role + "." + digestParts(role, parent, slot))
}

func evidenceRef(statement string) semreg.EvidenceRef {
	sum := sha256.Sum256([]byte(statement))
	return semreg.EvidenceRef{
		Owner:     "helianthus.ebus.native",
		Kind:      "evidence.native.statement",
		Digest:    semreg.Digest("sha256:" + hex.EncodeToString(sum[:])),
		Contract:  semreg.ContractKernelV1,
		Access:    semreg.EvidenceAccessPublic,
		Redaction: semreg.RedactionNone,
	}
}

type rootChoice struct {
	role     string
	kind     semreg.DefinitionID
	slot     string
	evidence string
	operator bool
	ok       bool
}

func chooseRoot(statement Statement) rootChoice {
	if statement.ControllerPresent {
		if slotOK(statement.ControllerSerial) {
			return rootChoice{role: "controller", kind: "resource.equipment", slot: statement.ControllerSerial, evidence: "controller-serial " + statement.ControllerSerial, ok: true}
		}
		if slotOK(statement.OperatorLabel) {
			return rootChoice{role: "controller", kind: "resource.equipment", slot: statement.OperatorLabel, evidence: "operator-label " + statement.OperatorLabel, operator: true, ok: true}
		}
		return rootChoice{}
	}
	if !statement.GeneratorPresent {
		return rootChoice{}
	}
	if slotOK(statement.GeneratorSerial) {
		return rootChoice{role: "generator", kind: "resource.heat_generator", slot: statement.GeneratorSerial, evidence: "generator-serial " + statement.GeneratorSerial, ok: true}
	}
	if slotOK(statement.OperatorLabel) {
		return rootChoice{role: "generator", kind: "resource.heat_generator", slot: statement.OperatorLabel, evidence: "operator-label " + statement.OperatorLabel, operator: true, ok: true}
	}
	return rootChoice{}
}

func (stamp originStamp) native(id semreg.ResourceID, evidence semreg.EvidenceRef) semreg.OriginRef {
	return semreg.OriginRef{
		OriginID:      semreg.OriginID("origin." + string(id)),
		Kind:          semreg.OriginNativeObservation,
		SourceID:      stamp.source,
		SourceEpochID: stamp.epoch,
		BindingID:     stamp.binding,
		Evidence:      []semreg.EvidenceRef{evidence},
	}
}

func operatorOrigin(id semreg.ResourceID, evidence semreg.EvidenceRef) semreg.OriginRef {
	return semreg.OriginRef{
		OriginID: semreg.OriginID("origin." + string(id)),
		Kind:     semreg.OriginOperator,
		Evidence: []semreg.EvidenceRef{evidence},
	}
}

func resourceRecord(item placed, id semreg.ResourceID, stamp originStamp) semreg.Resource {
	evidence := evidenceRef(item.evidence)
	origin := stamp.native(id, evidence)
	if item.operator {
		origin = operatorOrigin(id, evidence)
	}
	slot := item.slot
	anchor := semreg.ResourceAnchor{
		Class:    item.class,
		Basis:    item.basis,
		Slot:     &slot,
		Evidence: []semreg.EvidenceRef{evidence},
	}
	if item.parent != "" {
		parent := semreg.ResourceID(item.parent)
		anchor.Parent = &parent
	}
	return semreg.Resource{
		ResourceID: id,
		Kind:       item.kind,
		Anchor:     anchor,
		State:      semreg.ResourceCurrent,
		Origin:     origin,
		Revision:   "1",
	}
}

func positioned(role string, kind semreg.DefinitionID, class semreg.AnchorClass, parent, slot string) placed {
	return placed{
		role: role, kind: kind, class: class, basis: semreg.AnchorNativeSlot,
		parent: parent, slot: slot, evidence: "native-position " + slot,
	}
}

type resourceList struct {
	items []semreg.Resource
	seen  map[semreg.ResourceID]struct{}
	stamp originStamp
}

func (list *resourceList) add(item placed) semreg.ResourceID {
	id := resourceID(item.role, item.parent, item.slot)
	if _, ok := list.seen[id]; ok {
		return id
	}
	list.seen[id] = struct{}{}
	list.items = append(list.items, resourceRecord(item, id, list.stamp))
	return id
}

func buildResources(statement Statement, stamp originStamp) ([]semreg.Resource, graphIndex, bool) {
	root := chooseRoot(statement)
	if !root.ok {
		return nil, graphIndex{}, false
	}
	list := resourceList{seen: map[semreg.ResourceID]struct{}{}, stamp: stamp}
	basis := semreg.AnchorNativeSlot
	if root.operator {
		basis = semreg.AnchorOperatorSlot
	}
	rootID := list.add(placed{
		role: root.role, kind: root.kind, class: semreg.AnchorEquipment, basis: basis,
		slot: root.slot, evidence: root.evidence, operator: root.operator,
	})
	parent := string(rootID)
	index := graphIndex{
		root: rootID, rootEvidence: root.evidence,
		zones: map[byte]semreg.ResourceID{}, circuits: map[byte]semreg.ResourceID{},
	}
	index.system = list.add(positioned("system", "resource.thermal_system", semreg.AnchorFunction, parent, groupSlot(0, 0)))
	if root.role == "generator" {
		index.members = append(index.members, memberRef{id: rootID, role: "generator", slot: root.slot})
	} else if statement.GeneratorPresent {
		slot := "06/01/00"
		if slotOK(statement.GeneratorSerial) {
			slot = statement.GeneratorSerial
		}
		id := list.add(positioned("generator", "resource.heat_generator", semreg.AnchorEquipment, parent, slot))
		index.members = append(index.members, memberRef{id: id, role: "generator", slot: slot})
	}
	addZones(&list, &index, statement.Zones, parent)
	addCircuits(&list, &index, statement.Circuits, parent)
	addCylinders(&list, &index, statement.Cylinders, parent)
	if statement.HotWater {
		slot := groupSlot(1, 0)
		id := list.add(positioned("hot-water", "resource.domestic_hot_water", semreg.AnchorFunction, parent, slot))
		index.members = append(index.members, memberRef{id: id, role: "hot-water", slot: slot})
	}
	addHeatPumps(&list, &index, statement.HeatPumps, parent)
	sort.Slice(list.items, func(i, j int) bool { return list.items[i].ResourceID < list.items[j].ResourceID })
	return list.items, index, true
}

func addZones(list *resourceList, index *graphIndex, zones []Zone, parent string) {
	copied := append([]Zone(nil), zones...)
	sort.Slice(copied, func(i, j int) bool { return copied[i].Instance < copied[j].Instance })
	for _, zone := range copied {
		slot := groupSlot(3, zone.Instance)
		id := list.add(positioned("zone", "resource.zone", semreg.AnchorFunction, parent, slot))
		index.zones[zone.Instance] = id
		index.members = append(index.members, memberRef{id: id, role: "zone", slot: slot})
		if !zone.ValveReported {
			continue
		}
		valveSlot := slot + "/valve"
		valveID := list.add(positioned("valve", "resource.terminal_unit", semreg.AnchorEquipment, parent, valveSlot))
		index.members = append(index.members, memberRef{id: valveID, role: "valve", slot: valveSlot})
	}
}

func addCircuits(list *resourceList, index *graphIndex, circuits []Circuit, parent string) {
	copied := append([]Circuit(nil), circuits...)
	sort.Slice(copied, func(i, j int) bool { return copied[i].Instance < copied[j].Instance })
	for _, circuit := range copied {
		slot := groupSlot(2, circuit.Instance)
		id := list.add(positioned("circuit", "resource.thermal_circuit", semreg.AnchorFunction, parent, slot))
		index.circuits[circuit.Instance] = id
		index.members = append(index.members, memberRef{id: id, role: "circuit", slot: slot})
	}
}

func addCylinders(list *resourceList, index *graphIndex, cylinders []Cylinder, parent string) {
	copied := append([]Cylinder(nil), cylinders...)
	sort.Slice(copied, func(i, j int) bool { return copied[i].Instance < copied[j].Instance })
	for _, cylinder := range copied {
		slot := groupSlot(5, cylinder.Instance)
		id := list.add(positioned("store", "resource.thermal_store", semreg.AnchorEquipment, parent, slot))
		index.members = append(index.members, memberRef{id: id, role: "store", slot: slot})
	}
}

func addHeatPumps(list *resourceList, index *graphIndex, pumps []Pump, parent string) {
	copied := append([]Pump(nil), pumps...)
	sort.Slice(copied, func(i, j int) bool { return copied[i].Serial < copied[j].Serial })
	for _, pump := range copied {
		if !slotOK(pump.Serial) {
			continue
		}
		refrigerantSlot := "refrigerant/" + pump.Serial
		compressorSlot := "compressor/" + pump.Serial
		if !slotOK(refrigerantSlot) || !slotOK(compressorSlot) {
			continue
		}
		refrigerantID := list.add(positioned("refrigerant", "resource.refrigerant_circuit", semreg.AnchorFunction, parent, refrigerantSlot))
		index.members = append(index.members, memberRef{id: refrigerantID, role: "refrigerant", slot: refrigerantSlot})
		list.add(positioned("compressor", "resource.compressor", semreg.AnchorComponent, parent, compressorSlot))
	}
}
