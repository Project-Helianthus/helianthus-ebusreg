package heat

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"reflect"
	"sort"
	"testing"

	semreg "github.com/Project-Helianthus/helianthus-semreg/semreg/v1"
	semdriver "github.com/Project-Helianthus/helianthus-semreg/semreg/v1/driver"
)

// Graph fields, and the production mutation each assertion catches:
//   - controller serial: the root slot and every child parent. A hash that drops the serial keeps the old root id. A hash that drops the parent keeps a child id.
//   - zone instance: the zone slot. A hash that uses a fixed slot keeps the old zone id.
//   - valve slot: the zone slot plus the valve suffix. Reusing the zone id removes the zone when the valve is cleared.
//   - circuit instance: the circuit slot. A hash that uses the zone instance keeps the old circuit id. A map keyed by slice index keeps the fluid supply after the named circuit is gone.
//   - cylinder instance: the store slot. A hash that uses another group's instance keeps the old store id.
//   - heat-pump serial: the refrigerant and compressor slots. A hash that uses the controller serial keeps those ids.
//   - generator serial: the generator slot. A fixed slot keeps the old generator id.
//   - hot water and generator presence: those resources. Always emitting one keeps its id when it is absent.
//   - assignment: the fluid-supply from endpoint. Using the zone instance as the circuit index keeps the old from endpoint.
//   - unmatched assignment, including a nil register: no fluid-supply row. Emitting one invents a circuit.
//   - relation definition, from, to and endpoint kinds: the expected graph. A reversed pair or one kind for every resource misses it.
//   - anchor endpoint: class, basis, parent and slot. Dropping the parent or the slot misses the expected endpoint.
//   - resource type: the kind on that endpoint.
//   - source link binding: a hash of the root id. A hash that drops the root keeps the old binding when the serial changes.
//   - source link resource and state: the new root, still a candidate. Keeping the old root, or marking the link qualified, treats the serial change as a rebind.
//   - source link asset: the session asset. Leaving it blank keeps the previous asset on the link.
//   - discovery order: reversed zones, circuits, stores and heat pumps. An id taken from the input index diverges.
//   - clock: the observation time and clock id. An id or link that includes the clock diverges.
//   - fresh and reopened publisher state: the publication after a prior one. A republish that supersedes or drops a resource diverges.
//   - person rebind: an operator link changes the link resource and leaves resource ids. A serial change that does the same is a replacement.

func TestThermalTwoLifecycles(t *testing.T) {
	spec := loadLifecycleGraph(t)
	forward := spec.statement(false)
	backward := spec.statement(true)
	if forward.Zones[0].Instance == backward.Zones[0].Instance || forward.Circuits[0].Instance == backward.Circuits[0].Instance || forward.Cylinders[0].Instance == backward.Cylinders[0].Instance || forward.HeatPumps[0].Serial == backward.HeatPumps[0].Serial {
		t.Fatal("reverse discovery order matched the forward order")
	}
	freshPublisher := &Publisher{}
	fresh := publishLifecycle(t, freshPublisher, forward, "clock.utc", "1000000000", spec.SourceLinks[0].Asset, nil)
	reopenedPublisher := &Publisher{}
	publishLifecycle(t, reopenedPublisher, forward, "clock.utc", "1000000000", spec.SourceLinks[0].Asset, nil)
	reopened := publishLifecycle(t, reopenedPublisher, backward, "clock.reopened", "2000000000", spec.SourceLinks[0].Asset, nil)
	if fresh.ObservedAt == reopened.ObservedAt || fresh.Sequence == reopened.Sequence {
		t.Fatal("the two lifecycles shared a clock or a sequence")
	}
	want := spec.view(t)
	freshView := publishedView(t, fresh, spec)
	reopenedView := publishedView(t, reopened, spec)
	if !reflect.DeepEqual(freshView, want) || !reflect.DeepEqual(reopenedView, want) {
		t.Fatalf("lifecycles\nfresh %s\nreopened %s\nwant %s", mustJSON(freshView), mustJSON(reopenedView), mustJSON(want))
	}
	changed := spec.withSerial(t, "serial-beta")
	third := publishedView(t, publishLifecycle(t, &Publisher{}, changed.statement(false), "clock.utc", "3000000000", changed.SourceLinks[0].Asset, nil), changed)
	if !reflect.DeepEqual(third, changed.view(t)) {
		t.Fatalf("changed serial\n%s\nwant %s", mustJSON(third), mustJSON(changed.view(t)))
	}
	if share := sharedIDs(freshView, third); len(share) != 0 {
		t.Fatalf("changed serial kept ids %v", share)
	}
	assertPersonRebinds(t, spec, freshView)
	assertLifecycleFields(t, spec, freshView)
}

func assertPersonRebinds(t *testing.T, spec lifecycleGraph, resources lifecycleView) {
	t.Helper()
	root := testID("controller", "", spec.Installation.ControllerSerial)
	binding := testBinding(root)
	selected := semreg.ResourceID("operator.selected.resource")
	pending := semreg.IdentityLink{
		AssetID: semreg.AssetID(spec.SourceLinks[0].Asset), BindingID: binding, State: semreg.LinkQualified,
		Revision: "4", Resource: &selected, Basis: []semreg.EvidenceRef{testEvidenceRef("operator-declaration")},
	}
	batch := publishLifecycle(t, &Publisher{}, spec.statement(false), "clock.utc", "1000000000", spec.SourceLinks[0].Asset, func(upserts []semreg.NativeBinding) []semreg.IdentityLink {
		for _, upsert := range upserts {
			if upsert.BindingID == binding {
				return []semreg.IdentityLink{pending}
			}
		}
		return nil
	})
	view := publishedView(t, batch, spec)
	if !reflect.DeepEqual(view.Resources, resources.Resources) || !reflect.DeepEqual(view.Relations, resources.Relations) {
		t.Fatal("a person rebind replaced a resource or a relation")
	}
	if len(view.Links) != 1 || view.Links[0].Binding != string(binding) || view.Links[0].Resource != string(selected) || view.Links[0].State != string(semreg.LinkQualified) {
		t.Fatalf("person link = %+v", view.Links)
	}
	got := batch.IdentityLinkUpserts[0]
	if got.Revision != pending.Revision || len(got.Basis) != 1 || got.Basis[0].Digest != pending.Basis[0].Digest {
		t.Fatalf("person link was rewritten: %+v", got)
	}
}

func assertLifecycleFields(t *testing.T, spec lifecycleGraph, baseline lifecycleView) {
	t.Helper()
	parent := string(testID("controller", "", spec.Installation.ControllerSerial))
	systemID := string(testID("system", parent, testSlot(0, 0)))
	member := func(id string) string { return relationKey("thermal.relation.system_member", systemID, id) }
	fluid := func(from, to string) string { return relationKey("thermal.relation.fluid_supply", from, to) }
	zoneSlot := testSlot(3, 1)
	zoneID := string(testID("zone", parent, zoneSlot))
	valveID := string(testID("valve", parent, zoneSlot+"/valve"))
	circuitID := string(testID("circuit", parent, testSlot(2, 0)))
	otherCircuit := string(testID("circuit", parent, testSlot(2, 1)))
	storeID := string(testID("store", parent, testSlot(5, 4)))
	generatorID := string(testID("generator", parent, spec.Installation.GeneratorSerial))
	hotWaterID := string(testID("hot-water", parent, testSlot(1, 0)))
	refrigerantID := string(testID("refrigerant", parent, "refrigerant/pump-serial"))
	compressorID := string(testID("compressor", parent, "compressor/pump-serial"))

	movedZone := spec.clone(t)
	movedZone.Installation.Zones[0].Instance = 7
	movedZoneID := string(testID("zone", parent, testSlot(3, 7)))
	movedValveID := string(testID("valve", parent, testSlot(3, 7)+"/valve"))
	swapped := publishedView(t, publishSpec(t, movedZone), spec)
	expectSwapped(t, baseline, swapped, []string{zoneID, valveID}, []string{movedZoneID, movedValveID})
	expectRelationSwaps(t, baseline, swapped, []string{member(zoneID), member(valveID), fluid(circuitID, zoneID)}, []string{member(movedZoneID), member(movedValveID), fluid(circuitID, movedZoneID)})

	cleared := spec.clone(t)
	cleared.Installation.Zones[0].Valve = false
	swapped = publishedView(t, publishSpec(t, cleared), spec)
	expectSwapped(t, baseline, swapped, []string{valveID}, nil)
	expectRelationSwaps(t, baseline, swapped, []string{member(valveID)}, nil)
	if !idSet(swapped)[zoneID] {
		t.Fatal("clearing the valve removed the zone")
	}

	movedCircuit := spec.clone(t)
	movedCircuit.Installation.Circuits[0].Instance = 5
	movedCircuitID := string(testID("circuit", parent, testSlot(2, 5)))
	swapped = publishedView(t, publishSpec(t, movedCircuit), spec)
	expectSwapped(t, baseline, swapped, []string{circuitID}, []string{movedCircuitID})
	expectRelationSwaps(t, baseline, swapped, []string{member(circuitID), fluid(circuitID, zoneID)}, []string{member(movedCircuitID)})

	movedStore := spec.clone(t)
	movedStore.Installation.Cylinders[1].Instance = 9
	movedStoreID := string(testID("store", parent, testSlot(5, 9)))
	swapped = publishedView(t, publishSpec(t, movedStore), spec)
	expectSwapped(t, baseline, swapped, []string{storeID}, []string{movedStoreID})
	expectRelationSwaps(t, baseline, swapped, []string{member(storeID)}, []string{member(movedStoreID)})

	movedPump := spec.clone(t)
	movedPump.Installation.HeatPumps[1] = "pump-third"
	movedRefrigerantID := string(testID("refrigerant", parent, "refrigerant/pump-third"))
	swapped = publishedView(t, publishSpec(t, movedPump), spec)
	expectSwapped(t, baseline, swapped, []string{refrigerantID, compressorID}, []string{movedRefrigerantID, string(testID("compressor", parent, "compressor/pump-third"))})
	expectRelationSwaps(t, baseline, swapped, []string{member(refrigerantID)}, []string{member(movedRefrigerantID)})

	movedGenerator := spec.clone(t)
	movedGenerator.Installation.GeneratorSerial = "generator-other"
	movedGeneratorID := string(testID("generator", parent, "generator-other"))
	swapped = publishedView(t, publishSpec(t, movedGenerator), spec)
	expectSwapped(t, baseline, swapped, []string{generatorID}, []string{movedGeneratorID})
	expectRelationSwaps(t, baseline, swapped, []string{member(generatorID)}, []string{member(movedGeneratorID)})

	noHotWater := spec.clone(t)
	noHotWater.Installation.HotWater = false
	swapped = publishedView(t, publishSpec(t, noHotWater), spec)
	expectSwapped(t, baseline, swapped, []string{hotWaterID}, nil)
	expectRelationSwaps(t, baseline, swapped, []string{member(hotWaterID)}, nil)
	noGenerator := spec.clone(t)
	noGenerator.Installation.GeneratorPresent = false
	swapped = publishedView(t, publishSpec(t, noGenerator), spec)
	expectSwapped(t, baseline, swapped, []string{generatorID}, nil)
	expectRelationSwaps(t, baseline, swapped, []string{member(generatorID)}, nil)

	reassigned := spec.clone(t)
	one := 1
	reassigned.Installation.Zones[0].Assignment = &one
	swapped = publishedView(t, publishSpec(t, reassigned), spec)
	expectSwapped(t, baseline, swapped, nil, nil)
	expectRelationSwaps(t, baseline, swapped, []string{fluid(circuitID, zoneID)}, []string{fluid(otherCircuit, zoneID)})
	unassigned := spec.clone(t)
	unassigned.Installation.Zones[1].Assignment = nil
	swapped = publishedView(t, publishSpec(t, unassigned), spec)
	expectSwapped(t, baseline, swapped, nil, nil)
	expectRelationSwaps(t, baseline, swapped, []string{fluid(otherCircuit, string(testID("zone", parent, testSlot(3, 2))))}, nil)

	otherAsset := publishedView(t, publishLifecycle(t, &Publisher{}, spec.statement(false), "clock.utc", "1000000000", "asset.thermal.other", nil), spec)
	if !reflect.DeepEqual(otherAsset.Resources, baseline.Resources) || !reflect.DeepEqual(otherAsset.Relations, baseline.Relations) {
		t.Fatal("the asset changed a resource id or a relation key")
	}
	if len(otherAsset.Links) != 1 || otherAsset.Links[0].Asset != "asset.thermal.other" || otherAsset.Links[0].Binding != baseline.Links[0].Binding || otherAsset.Links[0].Resource != baseline.Links[0].Resource {
		t.Fatalf("asset link = %+v", otherAsset.Links)
	}
}

func publishSpec(t *testing.T, spec lifecycleGraph) semreg.PublicationBatch {
	t.Helper()
	return publishLifecycle(t, &Publisher{}, spec.statement(false), "clock.utc", "1000000000", spec.SourceLinks[0].Asset, nil)
}

func expectSwapped(t *testing.T, before, after lifecycleView, removed, added []string) {
	t.Helper()
	left, right := idSet(before), idSet(after)
	for _, id := range removed {
		if !left[id] || right[id] {
			t.Fatalf("removed %s before=%v after=%v", id, left[id], right[id])
		}
		delete(left, id)
	}
	for _, id := range added {
		if left[id] || !right[id] {
			t.Fatalf("added %s before=%v after=%v", id, left[id], right[id])
		}
		delete(right, id)
	}
	if !reflect.DeepEqual(left, right) {
		t.Fatalf("other ids changed\n%v\n%v", left, right)
	}
	for _, relation := range after.Relations {
		for _, id := range removed {
			if relation.From == id || relation.To == id {
				t.Fatalf("relation still uses %s", id)
			}
		}
	}
}

func expectRelationSwaps(t *testing.T, before, after lifecycleView, removed, added []string) {
	t.Helper()
	left, right := relationSet(before), relationSet(after)
	for _, key := range removed {
		if !left[key] || right[key] {
			t.Fatalf("removed relation %s\nbefore %v\nafter %v", key, left, right)
		}
		delete(left, key)
	}
	for _, key := range added {
		if left[key] || !right[key] {
			t.Fatalf("added relation %s\nbefore %v\nafter %v", key, left, right)
		}
		delete(right, key)
	}
	if !reflect.DeepEqual(left, right) {
		t.Fatalf("other relations changed\n%v\n%v", left, right)
	}
}

func relationKey(id, from, to string) string { return id + " " + from + " " + to }

func relationSet(view lifecycleView) map[string]bool {
	out := make(map[string]bool, len(view.Relations))
	for _, relation := range view.Relations {
		out[relationKey(relation.ID, relation.From, relation.To)] = true
	}
	return out
}

func idSet(view lifecycleView) map[string]bool {
	out := make(map[string]bool, len(view.Resources))
	for _, resource := range view.Resources {
		out[resource.ID] = true
	}
	return out
}

func sharedIDs(left, right lifecycleView) []string {
	rightIDs := idSet(right)
	var out []string
	for id := range idSet(left) {
		if rightIDs[id] {
			out = append(out, id)
		}
	}
	sort.Strings(out)
	return out
}

type lifecycleGraph struct {
	Installation    lifecycleInstallation `json:"installation"`
	Resources       []lifecycleResource   `json:"resources"`
	Relations       []lifecycleRelation   `json:"relations"`
	SourceLinks     []lifecycleLink       `json:"source_links"`
	RelationPack    string                `json:"relation_pack"`
	RelationVersion string                `json:"relation_version"`
	RelationBasis   string                `json:"relation_basis"`
	RelationState   string                `json:"relation_state"`
	ResourceState   string                `json:"resource_state"`
}

type lifecycleInstallation struct {
	ControllerSerial string              `json:"controller_serial"`
	GeneratorPresent bool                `json:"generator_present"`
	GeneratorSerial  string              `json:"generator_serial"`
	HotWater         bool                `json:"hot_water"`
	Zones            []lifecycleZone     `json:"zones"`
	Circuits         []lifecycleCircuit  `json:"circuits"`
	Cylinders        []lifecycleCylinder `json:"cylinders"`
	HeatPumps        []string            `json:"heat_pumps"`
}

type lifecycleZone struct {
	Instance   byte `json:"instance"`
	Assignment *int `json:"assignment"`
	Valve      bool `json:"valve"`
}

type lifecycleCircuit struct {
	Instance byte `json:"instance"`
}

type lifecycleCylinder struct {
	Instance byte `json:"instance"`
}

type lifecycleResource struct {
	Role       string `json:"role"`
	Kind       string `json:"kind"`
	Class      string `json:"class"`
	Basis      string `json:"basis"`
	ParentRole string `json:"parent_role,omitempty"`
	ParentSlot string `json:"parent_slot,omitempty"`
	Slot       string `json:"slot"`
}

type lifecycleRelation struct {
	ID       string `json:"id"`
	FromRole string `json:"from_role"`
	FromSlot string `json:"from_slot"`
	ToRole   string `json:"to_role"`
	ToSlot   string `json:"to_slot"`
	FromKind string `json:"from_kind"`
	ToKind   string `json:"to_kind"`
}

type lifecycleLink struct {
	Asset        string `json:"asset"`
	BindingRole  string `json:"binding_role"`
	BindingSlot  string `json:"binding_slot"`
	ResourceRole string `json:"resource_role"`
	ResourceSlot string `json:"resource_slot"`
	State        string `json:"state"`
	Revision     string `json:"revision"`
	Evidence     string `json:"evidence"`
}

type lifecycleView struct {
	Resources []lifecycleResourceView `json:"resources"`
	Relations []lifecycleRelationView `json:"relations"`
	Links     []lifecycleLinkView     `json:"source_links"`
}

type lifecycleResourceView struct {
	ID     string `json:"id"`
	Kind   string `json:"kind"`
	Class  string `json:"class"`
	Basis  string `json:"basis"`
	Parent string `json:"parent"`
	Slot   string `json:"slot"`
	State  string `json:"state"`
}

type lifecycleRelationView struct {
	ID       string `json:"id"`
	Pack     string `json:"pack"`
	Version  string `json:"version"`
	From     string `json:"from"`
	To       string `json:"to"`
	FromKind string `json:"from_kind"`
	ToKind   string `json:"to_kind"`
	Basis    string `json:"basis"`
	State    string `json:"state"`
}

type lifecycleLinkView struct {
	Asset    string `json:"asset"`
	Binding  string `json:"binding"`
	Resource string `json:"resource"`
	State    string `json:"state"`
	Revision string `json:"revision"`
}

func loadLifecycleGraph(t *testing.T) lifecycleGraph {
	t.Helper()
	body, err := os.ReadFile("testdata/expected_graph.json")
	if err != nil {
		t.Fatal(err)
	}
	var spec lifecycleGraph
	dec := json.NewDecoder(bytes.NewReader(body))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&spec); err != nil {
		t.Fatal(err)
	}
	if len(spec.Resources) == 0 || len(spec.Relations) == 0 || len(spec.SourceLinks) != 1 {
		t.Fatal("expected graph is incomplete")
	}
	return spec
}

func (g lifecycleGraph) clone(t *testing.T) lifecycleGraph {
	t.Helper()
	body, err := json.Marshal(g)
	if err != nil {
		t.Fatal(err)
	}
	var out lifecycleGraph
	if err := json.Unmarshal(body, &out); err != nil {
		t.Fatal(err)
	}
	return out
}

func (g lifecycleGraph) withSerial(t *testing.T, serial string) lifecycleGraph {
	t.Helper()
	next := g.clone(t)
	old := next.Installation.ControllerSerial
	next.Installation.ControllerSerial = serial
	for i := range next.Resources {
		if next.Resources[i].Slot == old {
			next.Resources[i].Slot = serial
		}
		if next.Resources[i].ParentSlot == old {
			next.Resources[i].ParentSlot = serial
		}
	}
	for i := range next.SourceLinks {
		link := &next.SourceLinks[i]
		if link.BindingSlot == old {
			link.BindingSlot = serial
		}
		if link.ResourceSlot == old {
			link.ResourceSlot = serial
		}
		link.Evidence = "controller-serial " + serial
	}
	return next
}

func (g lifecycleGraph) statement(reverse bool) Statement {
	zones := append([]lifecycleZone(nil), g.Installation.Zones...)
	circuits := append([]lifecycleCircuit(nil), g.Installation.Circuits...)
	cylinders := append([]lifecycleCylinder(nil), g.Installation.Cylinders...)
	pumps := append([]string(nil), g.Installation.HeatPumps...)
	if reverse {
		reverseSlice(zones)
		reverseSlice(circuits)
		reverseSlice(cylinders)
		reverseSlice(pumps)
	}
	out := Statement{
		ControllerPresent: true,
		ControllerSerial:  g.Installation.ControllerSerial,
		GeneratorPresent:  g.Installation.GeneratorPresent,
		GeneratorSerial:   g.Installation.GeneratorSerial,
		HotWater:          g.Installation.HotWater,
	}
	for _, zone := range zones {
		item := Zone{Instance: zone.Instance, ValveReported: zone.Valve}
		if zone.Assignment != nil {
			raw := uint16(*zone.Assignment)
			item.Assignment = &raw
		}
		out.Zones = append(out.Zones, item)
	}
	for _, circuit := range circuits {
		out.Circuits = append(out.Circuits, Circuit(circuit))
	}
	for _, cylinder := range cylinders {
		out.Cylinders = append(out.Cylinders, Cylinder(cylinder))
	}
	for _, serial := range pumps {
		out.HeatPumps = append(out.HeatPumps, Pump{Serial: serial})
	}
	return out
}

func reverseSlice[T any](values []T) {
	for i, j := 0, len(values)-1; i < j; i, j = i+1, j-1 {
		values[i], values[j] = values[j], values[i]
	}
}

func (g lifecycleGraph) resourceID(t *testing.T, role, slot string) string {
	t.Helper()
	var match *lifecycleResource
	for i := range g.Resources {
		if g.Resources[i].Role == role && g.Resources[i].Slot == slot {
			if match != nil {
				t.Fatalf("duplicate resource %s %s", role, slot)
			}
			match = &g.Resources[i]
		}
	}
	if match == nil {
		t.Fatalf("missing resource %s %s", role, slot)
	}
	parent := ""
	if match.ParentRole != "" {
		parent = g.resourceID(t, match.ParentRole, match.ParentSlot)
	}
	return string(testID(role, parent, slot))
}

func (g lifecycleGraph) view(t *testing.T) lifecycleView {
	t.Helper()
	view := lifecycleView{}
	for _, resource := range g.Resources {
		parent := ""
		if resource.ParentRole != "" {
			parent = g.resourceID(t, resource.ParentRole, resource.ParentSlot)
		}
		view.Resources = append(view.Resources, lifecycleResourceView{
			ID: string(testID(resource.Role, parent, resource.Slot)), Kind: resource.Kind,
			Class: resource.Class, Basis: resource.Basis, Parent: parent, Slot: resource.Slot, State: g.ResourceState,
		})
	}
	for _, relation := range g.Relations {
		view.Relations = append(view.Relations, lifecycleRelationView{
			ID: relation.ID, Pack: g.RelationPack, Version: g.RelationVersion,
			From: g.resourceID(t, relation.FromRole, relation.FromSlot), To: g.resourceID(t, relation.ToRole, relation.ToSlot),
			FromKind: relation.FromKind, ToKind: relation.ToKind, Basis: g.RelationBasis, State: g.RelationState,
		})
	}
	for _, link := range g.SourceLinks {
		view.Links = append(view.Links, lifecycleLinkView{
			Asset: link.Asset, Binding: string(testBinding(semreg.ResourceID(g.resourceID(t, link.BindingRole, link.BindingSlot)))),
			Resource: g.resourceID(t, link.ResourceRole, link.ResourceSlot), State: link.State, Revision: link.Revision,
		})
	}
	sortLifecycleView(&view)
	return view
}

func publishedView(t *testing.T, batch semreg.PublicationBatch, spec lifecycleGraph) lifecycleView {
	t.Helper()
	if len(batch.BindingUpserts) != 1 || len(batch.IdentityLinkUpserts) != 1 {
		t.Fatalf("bindings %d links %d", len(batch.BindingUpserts), len(batch.IdentityLinkUpserts))
	}
	link := batch.IdentityLinkUpserts[0]
	binding := batch.BindingUpserts[0]
	if binding.BindingID != link.BindingID || binding.State != semreg.BindingCurrent || binding.Revision != "1" {
		t.Fatalf("binding %+v link %+v", binding, link)
	}
	if binding.NativeResource.Digest != testEvidenceRef(spec.SourceLinks[0].Evidence).Digest {
		t.Fatal("binding evidence does not match the controller statement")
	}
	if link.State == semreg.LinkCandidate && (len(link.Basis) != 1 || link.Basis[0].Digest != binding.NativeResource.Digest) {
		t.Fatal("candidate source link evidence does not match the controller statement")
	}
	kinds := map[string]string{}
	view := lifecycleView{}
	for _, resource := range batch.ResourceUpserts {
		if err := resource.Validate(); err != nil {
			t.Fatalf("resource %s: %v", resource.ResourceID, err)
		}
		if resource.Supersedes != nil || resource.Anchor.Slot == nil {
			t.Fatalf("resource %s supersedes or has no slot", resource.ResourceID)
		}
		if resource.Origin.Kind != semreg.OriginNativeObservation || resource.Origin.BindingID == nil || *resource.Origin.BindingID != link.BindingID || resource.Origin.SourceID == nil || *resource.Origin.SourceID != batch.SourceID {
			t.Fatalf("resource %s origin %+v", resource.ResourceID, resource.Origin)
		}
		parent := ""
		if resource.Anchor.Parent != nil {
			parent = string(*resource.Anchor.Parent)
		}
		id := string(resource.ResourceID)
		if kinds[id] != "" {
			t.Fatalf("duplicate resource %s", id)
		}
		kinds[id] = string(resource.Kind)
		view.Resources = append(view.Resources, lifecycleResourceView{
			ID: id, Kind: string(resource.Kind), Class: string(resource.Anchor.Class), Basis: string(resource.Anchor.Basis),
			Parent: parent, Slot: *resource.Anchor.Slot, State: string(resource.State),
		})
	}
	for _, relation := range batch.RelationUpserts {
		if err := relation.Validate(); err != nil {
			t.Fatalf("relation %s: %v", relation.Definition.ID, err)
		}
		from, to := string(relation.From), string(relation.To)
		if kinds[from] == "" || kinds[to] == "" {
			t.Fatalf("relation endpoint missing %s -> %s", from, to)
		}
		view.Relations = append(view.Relations, lifecycleRelationView{
			ID: string(relation.Definition.ID), Pack: string(relation.Definition.Pack.ID), Version: string(relation.Definition.Version),
			From: from, To: to, FromKind: kinds[from], ToKind: kinds[to], Basis: string(relation.BasisClass), State: string(relation.State),
		})
	}
	resource := ""
	if link.Resource != nil {
		resource = string(*link.Resource)
	}
	view.Links = append(view.Links, lifecycleLinkView{
		Asset: string(link.AssetID), Binding: string(link.BindingID), Resource: resource, State: string(link.State), Revision: string(link.Revision),
	})
	sortLifecycleView(&view)
	return view
}

func sortLifecycleView(view *lifecycleView) {
	sort.Slice(view.Resources, func(i, j int) bool { return view.Resources[i].ID < view.Resources[j].ID })
	sort.Slice(view.Relations, func(i, j int) bool {
		if view.Relations[i].ID != view.Relations[j].ID {
			return view.Relations[i].ID < view.Relations[j].ID
		}
		if view.Relations[i].From != view.Relations[j].From {
			return view.Relations[i].From < view.Relations[j].From
		}
		return view.Relations[i].To < view.Relations[j].To
	})
	sort.Slice(view.Links, func(i, j int) bool { return view.Links[i].Binding < view.Links[j].Binding })
}

func mustJSON(value any) string {
	body, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(body)
}

func publishLifecycle(t *testing.T, publisher *Publisher, statement Statement, clock, nanos, asset string, links func([]semreg.NativeBinding) []semreg.IdentityLink) semreg.PublicationBatch {
	t.Helper()
	session := &lifecycleSession{clock: clock, nanos: nanos, asset: asset, links: links}
	if err := publisher.Publish(context.Background(), session, statement); err != nil {
		t.Fatal(err)
	}
	if len(session.staged) != 1 {
		t.Fatalf("staged %d batches", len(session.staged))
	}
	return session.staged[0]
}

type lifecycleSession struct {
	clock  string
	nanos  string
	asset  string
	links  func([]semreg.NativeBinding) []semreg.IdentityLink
	staged []semreg.PublicationBatch
}

func (s *lifecycleSession) Native() semdriver.NativeRegistry             { return nil }
func (s *lifecycleSession) SemanticInput() semdriver.SemanticInput       { return nil }
func (s *lifecycleSession) ScopedOperations() semdriver.ScopedOperations { return nil }
func (s *lifecycleSession) Calendar() semdriver.Calendar                 { return nil }
func (s *lifecycleSession) Publisher() semdriver.Publisher               { return s }
func (s *lifecycleSession) RegisterExecutor(semdriver.SettingRoute, semdriver.Executor) (func(), error) {
	return func() {}, nil
}
func (s *lifecycleSession) Prepare(mutate func(*semreg.PublicationBatch)) (semreg.PublicationBatch, error) {
	batch := semreg.PublicationBatch{
		Contract: semreg.ContractKernelV1, BatchID: "batch.thermal.synthetic", AssetID: semreg.AssetID(s.asset),
		SourceID: "source.thermal.synthetic", SourceEpochID: "epoch.thermal.synthetic", DriverGeneration: "1",
		ExpectedSemanticRevision: "0",
		ObservedAt:               semreg.TimePoint{UnixNanoseconds: semreg.Int64(s.nanos), ClockID: semreg.DefinitionID(s.clock), UncertaintyNS: "0"},
		SourceUpserts:            []semreg.SourceDescriptor{},
		SourceRetirements:        []semreg.SourceEpochID{},
		BindingUpserts:           []semreg.NativeBinding{},
		IdentityLinkUpserts:      []semreg.IdentityLink{},
		FactUpserts:              []semreg.FactCandidate{},
		FactWithdrawals:          []semreg.CandidateID{},
		ServiceUpserts:           []semreg.ServiceInstance{},
		ServiceWithdrawals:       []semreg.ServiceInstanceID{},
		CapabilityUpserts:        []semreg.CapabilityInstance{},
		CapabilityWithdrawals:    []semreg.CapabilityInstanceID{},
		GenerationFences:         []semreg.GenerationFence{},
	}
	if mutate != nil {
		mutate(&batch)
	}
	if batch.BatchDigest == "" {
		digest, err := batch.ComputedDigest()
		if err != nil {
			return semreg.PublicationBatch{}, err
		}
		batch.BatchDigest = digest
	}
	return batch, nil
}
func (s *lifecycleSession) Begin(context.Context) (semdriver.PublicationTransaction, error) {
	return &lifecycleTransaction{session: s}, nil
}

type lifecycleTransaction struct {
	session *lifecycleSession
}

func (tx *lifecycleTransaction) OperatorLinks(upserts []semreg.NativeBinding) []semreg.IdentityLink {
	if tx.session.links == nil {
		return nil
	}
	return tx.session.links(upserts)
}
func (tx *lifecycleTransaction) AdmitManufacturerSerial(semdriver.ManufacturerSerialClaim) error {
	return nil
}
func (tx *lifecycleTransaction) Stage(batch semreg.PublicationBatch) error {
	tx.session.staged = append(tx.session.staged, batch)
	return nil
}
func (tx *lifecycleTransaction) DeclareSource(semdriver.Source) error { return nil }
func (tx *lifecycleTransaction) DeclareSettingRoute(semdriver.SettingRoute) error {
	return nil
}
func (tx *lifecycleTransaction) StageSample(semdriver.Sample) error  { return nil }
func (tx *lifecycleTransaction) StageLink(semdriver.LinkClaim) error { return nil }
func (tx *lifecycleTransaction) Commit(context.Context) error        { return nil }
func (tx *lifecycleTransaction) Abort(context.Context) error         { return nil }
