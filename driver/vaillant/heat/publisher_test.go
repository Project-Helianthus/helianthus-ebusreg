package heat

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"reflect"
	"testing"

	semreg "github.com/Project-Helianthus/helianthus-semreg/semreg/v1"
	semdriver "github.com/Project-Helianthus/helianthus-semreg/semreg/v1/driver"
)

// Identity fields, and the production mutation each assertion catches:
//   - controller serial: the root slot. A hash that drops the serial keeps the old root id.
//   - parent: every child id includes the root id. A hash that drops the parent keeps child ids when the serial changes.
//   - zone instance: the zone slot. A hash that uses a fixed slot keeps the old zone id.
//   - valve slot: the zone slot plus the valve suffix. A hash that reuses the zone id keeps the valve id equal to the zone id.
//   - circuit instance: the circuit slot. A hash that uses the zone instance keeps the old circuit id.
//   - cylinder instance: the store slot. A hash that uses another instance keeps the old store id.
//   - heat-pump serial: the refrigerant and compressor slots. A hash that uses the controller serial keeps those ids.
//   - pump type: Statement.HeatPumps elements are Pump, so other packages see heat.Pump.
//   - operator label: the operator slot and operator origin. A native origin on that root keeps a source pointer.
//   - generator serial: the root slot when no controller is present. A second generator resource is a duplicate root.
//   - binding: "binding", NUL, root id. A pending link for another formula is not applied.
//   - committed operator link: a later batch leaves it. A candidate upsert for that binding is refused.
//   - supersedes: stays empty. A replacement record sets it.

func TestThermalPublisher(t *testing.T) {
	t.Run("graph", testPublisherGraph)
	t.Run("identity fields", testPublisherIdentityFields)
	t.Run("operator links", testPublisherOperatorLinks)
	t.Run("pump type", testPumpTypeName)
}

func testPumpTypeName(t *testing.T) {
	field, ok := reflect.TypeOf(Statement{}).FieldByName("HeatPumps")
	if !ok {
		t.Fatal("Statement.HeatPumps missing")
	}
	if field.Type.Kind() != reflect.Slice {
		t.Fatalf("HeatPumps kind = %s", field.Type.Kind())
	}
	if name := field.Type.Elem().Name(); name != "Pump" {
		t.Fatalf("HeatPumps element type = %s, want Pump", name)
	}
}

func testPublisherGraph(t *testing.T) {
	batch := publishStatement(t, baseStatement(), nil)
	byID := resourcesByID(t, batch)
	root := testID("controller", "", "serial-alpha")
	parent := string(root)
	expectResource(t, byID, root, "resource.equipment", semreg.AnchorEquipment, semreg.AnchorNativeSlot, "", "serial-alpha")
	expectResource(t, byID, testID("system", parent, testSlot(0, 0)), "resource.thermal_system", semreg.AnchorFunction, semreg.AnchorNativeSlot, parent, testSlot(0, 0))
	expectResource(t, byID, testID("zone", parent, testSlot(3, 1)), "resource.zone", semreg.AnchorFunction, semreg.AnchorNativeSlot, parent, testSlot(3, 1))
	expectResource(t, byID, testID("valve", parent, testSlot(3, 1)+"/valve"), "resource.terminal_unit", semreg.AnchorEquipment, semreg.AnchorNativeSlot, parent, testSlot(3, 1)+"/valve")
	expectResource(t, byID, testID("circuit", parent, testSlot(2, 0)), "resource.thermal_circuit", semreg.AnchorFunction, semreg.AnchorNativeSlot, parent, testSlot(2, 0))
	expectResource(t, byID, testID("store", parent, testSlot(5, 4)), "resource.thermal_store", semreg.AnchorEquipment, semreg.AnchorNativeSlot, parent, testSlot(5, 4))
	expectResource(t, byID, testID("hot-water", parent, testSlot(1, 0)), "resource.domestic_hot_water", semreg.AnchorFunction, semreg.AnchorNativeSlot, parent, testSlot(1, 0))
	expectResource(t, byID, testID("generator", parent, "generator-serial"), "resource.heat_generator", semreg.AnchorEquipment, semreg.AnchorNativeSlot, parent, "generator-serial")
	expectResource(t, byID, testID("refrigerant", parent, "refrigerant/pump-serial"), "resource.refrigerant_circuit", semreg.AnchorFunction, semreg.AnchorNativeSlot, parent, "refrigerant/pump-serial")
	compressor := testID("compressor", parent, "compressor/pump-serial")
	expectResource(t, byID, compressor, "resource.compressor", semreg.AnchorComponent, semreg.AnchorNativeSlot, parent, "compressor/pump-serial")
	if len(byID) != 10 {
		t.Fatalf("resource count = %d, want 10", len(byID))
	}
	assignment := uint16ptr(0)
	forward := baseStatement()
	forward.Zones = []Zone{{Instance: 1, Assignment: assignment, ValveReported: true}, {Instance: 2}}
	forward.Circuits = []Circuit{{Instance: 0}, {Instance: 4}}
	forward.Cylinders = []Cylinder{{Instance: 4}, {Instance: 1}}
	forward.HeatPumps = []Pump{{Serial: "pump-a"}, {Serial: "pump-serial"}}
	backward := forward
	backward.Zones = []Zone{forward.Zones[1], forward.Zones[0]}
	backward.Circuits = []Circuit{forward.Circuits[1], forward.Circuits[0]}
	backward.Cylinders = []Cylinder{forward.Cylinders[1], forward.Cylinders[0]}
	backward.HeatPumps = []Pump{forward.HeatPumps[1], forward.HeatPumps[0]}
	if !reflect.DeepEqual(resourceIDs(publishStatement(t, forward, nil)), resourceIDs(publishStatement(t, backward, nil))) {
		t.Fatal("reversed inventory changed resource ids")
	}
}

func testPublisherIdentityFields(t *testing.T) {
	original := resourcesByID(t, publishStatement(t, baseStatement(), nil))
	zoneChanged := baseStatement()
	zoneChanged.Zones[0].Instance = 2
	changed := resourcesByID(t, publishStatement(t, zoneChanged, nil))
	expectReplaced(t, cloneResources(original), changed, []semreg.ResourceID{
		testID("zone", string(testID("controller", "", "serial-alpha")), testSlot(3, 1)),
		testID("valve", string(testID("controller", "", "serial-alpha")), testSlot(3, 1)+"/valve"),
	}, []semreg.ResourceID{
		testID("zone", string(testID("controller", "", "serial-alpha")), testSlot(3, 2)),
		testID("valve", string(testID("controller", "", "serial-alpha")), testSlot(3, 2)+"/valve"),
	})

	circuitChanged := baseStatement()
	circuitChanged.Circuits[0].Instance = 3
	changed = resourcesByID(t, publishStatement(t, circuitChanged, nil))
	parent := string(testID("controller", "", "serial-alpha"))
	expectReplaced(t, cloneResources(original), changed, []semreg.ResourceID{testID("circuit", parent, testSlot(2, 0))}, []semreg.ResourceID{testID("circuit", parent, testSlot(2, 3))})

	storeChanged := baseStatement()
	storeChanged.Cylinders[0].Instance = 9
	changed = resourcesByID(t, publishStatement(t, storeChanged, nil))
	expectReplaced(t, cloneResources(original), changed, []semreg.ResourceID{testID("store", parent, testSlot(5, 4))}, []semreg.ResourceID{testID("store", parent, testSlot(5, 9))})

	valveGone := baseStatement()
	valveGone.Zones[0].ValveReported = false
	changed = resourcesByID(t, publishStatement(t, valveGone, nil))
	expectReplaced(t, cloneResources(original), changed, []semreg.ResourceID{testID("valve", parent, testSlot(3, 1)+"/valve")}, nil)
	if _, ok := changed[testID("zone", parent, testSlot(3, 1))]; !ok {
		t.Fatal("clearing the valve removed the zone")
	}

	pumpChanged := baseStatement()
	pumpChanged.HeatPumps[0].Serial = "pump-other"
	changed = resourcesByID(t, publishStatement(t, pumpChanged, nil))
	expectReplaced(t, cloneResources(original), changed, []semreg.ResourceID{
		testID("refrigerant", parent, "refrigerant/pump-serial"),
		testID("compressor", parent, "compressor/pump-serial"),
	}, []semreg.ResourceID{
		testID("refrigerant", parent, "refrigerant/pump-other"),
		testID("compressor", parent, "compressor/pump-other"),
	})

	moved := baseStatement()
	moved.ControllerSerial = "serial-beta"
	changed = resourcesByID(t, publishStatement(t, moved, nil))
	if len(intersection(original, changed)) != 0 {
		t.Fatal("changed controller serial kept a resource id")
	}
	for _, resource := range changed {
		if resource.Supersedes != nil {
			t.Fatal("changed controller serial set a supersession")
		}
	}

	label := Statement{ControllerPresent: true, OperatorLabel: "operator-label", HotWater: true}
	labeled := resourcesByID(t, publishStatement(t, label, nil))
	root := labeled[testID("controller", "", "operator-label")]
	if root.Anchor.Basis != semreg.AnchorOperatorSlot || root.Origin.Kind != semreg.OriginOperator || root.Origin.SourceID != nil || root.Origin.SourceEpochID != nil || root.Origin.BindingID != nil {
		t.Fatalf("operator root origin = %+v anchor %+v", root.Origin, root.Anchor)
	}
	hotWater := labeled[testID("hot-water", string(root.ResourceID), testSlot(1, 0))]
	if hotWater.Origin.Kind != semreg.OriginNativeObservation || hotWater.Origin.SourceID == nil || hotWater.Origin.BindingID == nil {
		t.Fatalf("child of operator root origin = %+v", hotWater.Origin)
	}

	generated := Statement{GeneratorPresent: true, GeneratorSerial: "generator-serial", HotWater: true}
	generatedBatch := resourcesByID(t, publishStatement(t, generated, nil))
	generatorID := testID("generator", "", "generator-serial")
	generator := generatedBatch[generatorID]
	if generator.Kind != "resource.heat_generator" || generator.Anchor.Parent != nil {
		t.Fatalf("generator root = %+v", generator)
	}
	if countKind(generatedBatch, "resource.heat_generator") != 1 {
		t.Fatal("generator root was published twice")
	}

	session := &recordingSession{}
	if err := (&Publisher{}).Publish(context.Background(), session, Statement{}); err != nil || session.begins != 0 {
		t.Fatalf("empty statement begin=%d err=%v", session.begins, err)
	}
	session = &recordingSession{}
	if err := (&Publisher{}).Publish(context.Background(), session, Statement{ControllerPresent: true, ControllerSerial: "bad\nserial"}); err != nil || session.begins != 0 {
		t.Fatalf("unusable serial begin=%d err=%v", session.begins, err)
	}
}

func testPublisherOperatorLinks(t *testing.T) {
	built := []semreg.IdentityLink{{BindingID: "binding.thermal.one", Revision: "1"}}
	if got := carryOperatorLinks(built, nil); len(got) != 1 || &got[0] != &built[0] {
		t.Fatal("nil pending set did not return the built slice")
	}
	if got := carryOperatorLinks(built, []semreg.IdentityLink{}); len(got) != 1 || &got[0] != &built[0] {
		t.Fatal("empty pending set did not return the built slice")
	}
	otherResource := semreg.ResourceID("resource.kept")
	other := semreg.IdentityLink{BindingID: "binding.thermal.other", State: semreg.LinkCandidate, Revision: "7", Resource: &otherResource, Basis: []semreg.EvidenceRef{testEvidenceRef("kept")}}
	root := testID("controller", "", "serial-alpha")
	binding := testBinding(root)
	pendingResource := semreg.ResourceID("operator.selected.resource")
	pending := semreg.IdentityLink{
		AssetID: "asset.thermal.fixed", BindingID: binding, State: semreg.LinkQualified, Revision: "4",
		Resource: &pendingResource, Basis: []semreg.EvidenceRef{testEvidenceRef("operator-declaration")},
	}
	publisherLink := semreg.IdentityLink{BindingID: binding, State: semreg.LinkCandidate, Revision: "1"}
	carried := carryOperatorLinks([]semreg.IdentityLink{publisherLink, other}, []semreg.IdentityLink{pending})
	if got := linkByBinding(t, carried, other.BindingID); !reflect.DeepEqual(got, other) {
		t.Fatalf("carry dropped another binding: %+v", got)
	}
	if got := linkByBinding(t, carried, binding); !reflect.DeepEqual(got, pending) {
		t.Fatalf("carry kept the publisher link: %+v", got)
	}

	first := publishStatement(t, baseStatement(), nil)
	second := publishStatement(t, baseStatement(), nil)
	left, err := semreg.CanonicalJSON(first)
	if err != nil {
		t.Fatal(err)
	}
	right, err := semreg.CanonicalJSON(second)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(left, right) {
		t.Fatal("empty pending set changed the batch bytes")
	}

	var seen []semreg.NativeBinding
	alpha := publishStatement(t, baseStatement(), func(upserts []semreg.NativeBinding) []semreg.IdentityLink {
		seen = append(seen, upserts...)
		for _, upsert := range upserts {
			if upsert.BindingID == binding {
				return []semreg.IdentityLink{pending}
			}
		}
		return nil
	})
	if len(seen) != 1 || seen[0].BindingID != binding {
		t.Fatalf("operator links saw %#v", seen)
	}
	if got := linkByBinding(t, alpha.IdentityLinkUpserts, binding); !reflect.DeepEqual(got, pending) {
		t.Fatalf("pending link = %+v", got)
	}
	withdrawnResource := semreg.ResourceID("operator.withdrawn.resource")
	withdrawal := pending
	withdrawal.State = semreg.LinkWithdrawn
	withdrawal.Revision = "2"
	withdrawal.Resource = &withdrawnResource
	withdrawal.Basis = []semreg.EvidenceRef{testEvidenceRef("operator-withdrawal")}
	withdrawn := publishStatement(t, baseStatement(), func(upserts []semreg.NativeBinding) []semreg.IdentityLink {
		for _, upsert := range upserts {
			if upsert.BindingID == binding {
				return []semreg.IdentityLink{withdrawal}
			}
		}
		return nil
	})
	if got := linkByBinding(t, withdrawn.IdentityLinkUpserts, binding); !reflect.DeepEqual(got, withdrawal) {
		t.Fatalf("withdrawal = %+v", got)
	}

	betaStatement := baseStatement()
	betaStatement.ControllerSerial = "serial-beta"
	betaBinding := testBinding(testID("controller", "", "serial-beta"))
	beta := publishStatement(t, betaStatement, func(upserts []semreg.NativeBinding) []semreg.IdentityLink {
		for _, upsert := range upserts {
			if upsert.BindingID == binding {
				return []semreg.IdentityLink{pending}
			}
		}
		return nil
	})
	got := linkByBinding(t, beta.IdentityLinkUpserts, betaBinding)
	if got.State != semreg.LinkCandidate || got.Resource == nil || *got.Resource != testID("controller", "", "serial-beta") {
		t.Fatalf("unrelated batch carried the pending link: %+v", got)
	}
	for _, link := range beta.IdentityLinkUpserts {
		if link.BindingID == binding {
			t.Fatal("beta batch carried the alpha binding")
		}
	}

	publisher := &Publisher{}
	carrying := &recordingSession{links: func(upserts []semreg.NativeBinding) []semreg.IdentityLink {
		for _, upsert := range upserts {
			if upsert.BindingID == binding {
				return []semreg.IdentityLink{pending}
			}
		}
		return nil
	}}
	if err := publisher.Publish(context.Background(), carrying, baseStatement()); err != nil {
		t.Fatal(err)
	}
	following := &recordingSession{}
	if err := publisher.Publish(context.Background(), following, baseStatement()); err != nil {
		t.Fatal(err)
	}
	if len(following.staged) != 1 {
		t.Fatalf("staged %d batches after the committed operator link", len(following.staged))
	}
	if following.staged[0].IdentityLinkUpserts == nil || hasBindingLink(following.staged[0].IdentityLinkUpserts, binding) {
		t.Fatal("the batch after the committed operator link replaced it with the candidate")
	}
	if len(following.staged[0].ResourceUpserts) == 0 {
		t.Fatal("leaving the operator link dropped the thermal graph")
	}
	withdrawing := &recordingSession{links: func(upserts []semreg.NativeBinding) []semreg.IdentityLink {
		for _, upsert := range upserts {
			if upsert.BindingID == binding {
				return []semreg.IdentityLink{withdrawal}
			}
		}
		return nil
	}}
	if err := publisher.Publish(context.Background(), withdrawing, baseStatement()); err != nil {
		t.Fatal(err)
	}
	if got := linkByBinding(t, withdrawing.staged[0].IdentityLinkUpserts, binding); got.State != semreg.LinkWithdrawn {
		t.Fatalf("withdrawal after a committed operator link = %+v", got)
	}
	restored := &recordingSession{}
	if err := publisher.Publish(context.Background(), restored, baseStatement()); err != nil {
		t.Fatal(err)
	}
	if got := linkByBinding(t, restored.staged[0].IdentityLinkUpserts, binding); got.State != semreg.LinkCandidate {
		t.Fatalf("candidate after the withdrawal = %+v", got)
	}

	fresh := &Publisher{}
	refused := &recordingSession{committedOperator: binding}
	if err := fresh.Publish(context.Background(), refused, baseStatement()); err != nil {
		t.Fatalf("a committed operator link with an empty pending set rejected the batch: %v", err)
	}
	if len(refused.staged) != 2 {
		t.Fatalf("staged %d batches, want the refused candidate and the batch that leaves the operator link", len(refused.staged))
	}
	if !hasBindingLink(refused.staged[0].IdentityLinkUpserts, binding) {
		t.Fatal("the refused batch did not upsert the candidate")
	}
	if refused.staged[1].IdentityLinkUpserts == nil || hasBindingLink(refused.staged[1].IdentityLinkUpserts, binding) {
		t.Fatal("the batch that follows a committed operator link replaced it")
	}
	if len(refused.staged[1].ResourceUpserts) == 0 {
		t.Fatal("leaving the operator link dropped the thermal graph")
	}
}

func baseStatement() Statement {
	assignment := uint16(0)
	return Statement{
		ControllerPresent: true,
		ControllerSerial:  "serial-alpha",
		GeneratorPresent:  true,
		GeneratorSerial:   "generator-serial",
		HotWater:          true,
		Zones:             []Zone{{Instance: 1, Assignment: &assignment, ValveReported: true}},
		Circuits:          []Circuit{{Instance: 0}},
		Cylinders:         []Cylinder{{Instance: 4}},
		HeatPumps:         []Pump{{Serial: "pump-serial"}},
	}
}

func publishStatement(t *testing.T, statement Statement, links func([]semreg.NativeBinding) []semreg.IdentityLink) semreg.PublicationBatch {
	t.Helper()
	session := &recordingSession{links: links}
	if err := (&Publisher{}).Publish(context.Background(), session, statement); err != nil {
		t.Fatal(err)
	}
	if len(session.staged) != 1 {
		t.Fatalf("staged %d batches", len(session.staged))
	}
	return session.staged[0]
}

func resourcesByID(t *testing.T, batch semreg.PublicationBatch) map[semreg.ResourceID]semreg.Resource {
	t.Helper()
	out := map[semreg.ResourceID]semreg.Resource{}
	for _, resource := range batch.ResourceUpserts {
		if err := resource.Validate(); err != nil {
			t.Fatalf("resource %s: %v", resource.ResourceID, err)
		}
		if resource.Supersedes != nil {
			t.Fatalf("resource %s supersedes %+v", resource.ResourceID, resource.Supersedes)
		}
		out[resource.ResourceID] = resource
	}
	return out
}

func expectResource(t *testing.T, byID map[semreg.ResourceID]semreg.Resource, id semreg.ResourceID, kind semreg.DefinitionID, class semreg.AnchorClass, basis semreg.AnchorBasis, parent, slot string) {
	t.Helper()
	resource, ok := byID[id]
	if !ok {
		t.Fatalf("missing %s %s", kind, id)
	}
	if resource.Kind != kind || resource.Anchor.Class != class || resource.Anchor.Basis != basis || resource.Anchor.Slot == nil || *resource.Anchor.Slot != slot {
		t.Fatalf("%s anchor = %+v", kind, resource.Anchor)
	}
	if parent == "" {
		if resource.Anchor.Parent != nil {
			t.Fatalf("%s parent = %s", kind, *resource.Anchor.Parent)
		}
		return
	}
	if resource.Anchor.Parent == nil || string(*resource.Anchor.Parent) != parent {
		t.Fatalf("%s parent = %v, want %s", kind, resource.Anchor.Parent, parent)
	}
}

func expectReplaced(t *testing.T, before, after map[semreg.ResourceID]semreg.Resource, removed, added []semreg.ResourceID) {
	t.Helper()
	for _, id := range removed {
		if _, ok := after[id]; ok {
			t.Fatalf("removed id still present: %s", id)
		}
		delete(before, id)
	}
	for _, id := range added {
		if _, ok := after[id]; !ok {
			t.Fatalf("added id missing: %s", id)
		}
		delete(after, id)
	}
	if !reflect.DeepEqual(resourceIDKeys(before), resourceIDKeys(after)) {
		t.Fatalf("unchanged ids = %v, want %v", resourceIDKeys(after), resourceIDKeys(before))
	}
}

func cloneResources(in map[semreg.ResourceID]semreg.Resource) map[semreg.ResourceID]semreg.Resource {
	out := make(map[semreg.ResourceID]semreg.Resource, len(in))
	for id, resource := range in {
		out[id] = resource
	}
	return out
}

func intersection(left, right map[semreg.ResourceID]semreg.Resource) []semreg.ResourceID {
	var out []semreg.ResourceID
	for id := range left {
		if _, ok := right[id]; ok {
			out = append(out, id)
		}
	}
	return out
}

func resourceIDs(batch semreg.PublicationBatch) []semreg.ResourceID {
	out := make([]semreg.ResourceID, 0, len(batch.ResourceUpserts))
	for _, resource := range batch.ResourceUpserts {
		out = append(out, resource.ResourceID)
	}
	return out
}

func resourceIDKeys(resources map[semreg.ResourceID]semreg.Resource) []string {
	out := make([]string, 0, len(resources))
	for id := range resources {
		out = append(out, string(id))
	}
	sortStrings(out)
	return out
}

func sortStrings(values []string) {
	for i := 1; i < len(values); i++ {
		for j := i; j > 0 && values[j] < values[j-1]; j-- {
			values[j], values[j-1] = values[j-1], values[j]
		}
	}
}

func countKind(resources map[semreg.ResourceID]semreg.Resource, kind semreg.DefinitionID) int {
	count := 0
	for _, resource := range resources {
		if resource.Kind == kind {
			count++
		}
	}
	return count
}

func linkByBinding(t *testing.T, links []semreg.IdentityLink, id semreg.NativeBindingID) semreg.IdentityLink {
	t.Helper()
	for _, link := range links {
		if link.BindingID == id {
			return link
		}
	}
	t.Fatalf("missing link %s", id)
	return semreg.IdentityLink{}
}

func hasBindingLink(links []semreg.IdentityLink, id semreg.NativeBindingID) bool {
	for _, link := range links {
		if link.BindingID == id {
			return true
		}
	}
	return false
}

func testDigest(parts ...string) string {
	hash := sha256.New()
	for i, part := range parts {
		if i > 0 {
			_, _ = hash.Write([]byte{0})
		}
		_, _ = hash.Write([]byte(part))
	}
	return hex.EncodeToString(hash.Sum(nil))
}

func testID(role, parent, slot string) semreg.ResourceID {
	return semreg.ResourceID("thermal." + role + "." + testDigest(role, parent, slot))
}

func testSlot(group, instance byte) string {
	return fmt.Sprintf("%02x/%02x/%02x", byte(2), group, instance)
}

func testBinding(root semreg.ResourceID) semreg.NativeBindingID {
	return semreg.NativeBindingID("binding.thermal." + testDigest("binding", string(root)))
}

func testEvidenceRef(statement string) semreg.EvidenceRef {
	sum := sha256.Sum256([]byte(statement))
	return semreg.EvidenceRef{
		Owner: "helianthus.ebus.native", Kind: "evidence.native.statement",
		Digest: semreg.Digest("sha256:" + hex.EncodeToString(sum[:])), Contract: semreg.ContractKernelV1,
		Access: semreg.EvidenceAccessPublic, Redaction: semreg.RedactionNone,
	}
}

type recordingSession struct {
	begins int
	staged []semreg.PublicationBatch
	links  func([]semreg.NativeBinding) []semreg.IdentityLink
	// committedOperator is a binding whose published operator link already
	// matches its declaration, so OperatorLinks stays empty. A batch that
	// upserts a link for it is the replacement the runtime refuses.
	committedOperator semreg.NativeBindingID
}

func (s *recordingSession) Native() semdriver.NativeRegistry             { return nil }
func (s *recordingSession) SemanticInput() semdriver.SemanticInput       { return nil }
func (s *recordingSession) ScopedOperations() semdriver.ScopedOperations { return nil }
func (s *recordingSession) Calendar() semdriver.Calendar                 { return nil }
func (s *recordingSession) Publisher() semdriver.Publisher               { return s }
func (s *recordingSession) RegisterExecutor(semdriver.SettingRoute, semdriver.Executor) (func(), error) {
	return func() {}, nil
}
func (s *recordingSession) Prepare(mutate func(*semreg.PublicationBatch)) (semreg.PublicationBatch, error) {
	batch := semreg.PublicationBatch{
		Contract: semreg.ContractKernelV1, BatchID: "batch.thermal.fixed", AssetID: "asset.thermal.fixed",
		SourceID: "source.thermal.fixed", SourceEpochID: "epoch.thermal.fixed", DriverGeneration: "1",
		Sequence: "1", ExpectedSemanticRevision: "0",
		ObservedAt:            semreg.TimePoint{UnixNanoseconds: "0", ClockID: "clock.utc", UncertaintyNS: "0"},
		SourceUpserts:         []semreg.SourceDescriptor{},
		SourceRetirements:     []semreg.SourceEpochID{},
		BindingUpserts:        []semreg.NativeBinding{},
		IdentityLinkUpserts:   []semreg.IdentityLink{},
		FactUpserts:           []semreg.FactCandidate{},
		FactWithdrawals:       []semreg.CandidateID{},
		ServiceUpserts:        []semreg.ServiceInstance{},
		ServiceWithdrawals:    []semreg.ServiceInstanceID{},
		CapabilityUpserts:     []semreg.CapabilityInstance{},
		CapabilityWithdrawals: []semreg.CapabilityInstanceID{},
		GenerationFences:      []semreg.GenerationFence{},
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
func (s *recordingSession) Begin(context.Context) (semdriver.PublicationTransaction, error) {
	s.begins++
	return &recordingTransaction{session: s}, nil
}

type recordingTransaction struct {
	session *recordingSession
}

func (tx *recordingTransaction) OperatorLinks(upserts []semreg.NativeBinding) []semreg.IdentityLink {
	if tx.session.links == nil {
		return nil
	}
	return tx.session.links(upserts)
}
func (tx *recordingTransaction) AdmitManufacturerSerial(semdriver.ManufacturerSerialClaim) error {
	return nil
}
func (tx *recordingTransaction) Stage(batch semreg.PublicationBatch) error {
	tx.session.staged = append(tx.session.staged, batch)
	return nil
}
func (tx *recordingTransaction) DeclareSource(semdriver.Source) error { return nil }
func (tx *recordingTransaction) DeclareSettingRoute(semdriver.SettingRoute) error {
	return nil
}
func (tx *recordingTransaction) StageSample(semdriver.Sample) error  { return nil }
func (tx *recordingTransaction) StageLink(semdriver.LinkClaim) error { return nil }
func (tx *recordingTransaction) Commit(context.Context) error {
	if tx.session.committedOperator == "" || len(tx.session.staged) == 0 {
		return nil
	}
	batch := tx.session.staged[len(tx.session.staged)-1]
	for _, link := range batch.IdentityLinkUpserts {
		if link.BindingID == tx.session.committedOperator {
			return errors.New(operatorLinkRefusal)
		}
	}
	return nil
}
func (tx *recordingTransaction) Abort(context.Context) error { return nil }
