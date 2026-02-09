package registry

import "testing"

func TestProjectionPath_StringAndValidate(t *testing.T) {
	path := ProjectionPath{
		Plane: "Observability",
		Segments: []PathSegment{
			{Name: "devices"},
			{Name: "boiler"},
			{Name: "zone1", Location: true},
		},
	}

	if got := path.String(); got != "Observability:/devices/boiler/@zone1" {
		t.Fatalf("unexpected path string: %s", got)
	}
	if err := path.Validate(); err != nil {
		t.Fatalf("unexpected validate error: %v", err)
	}

	invalid := []ProjectionPath{
		{Plane: "", Segments: []PathSegment{{Name: "devices"}}},
		{Plane: "Bad/Plane", Segments: []PathSegment{{Name: "devices"}}},
		{Plane: "Bad:Plane", Segments: []PathSegment{{Name: "devices"}}},
		{Plane: "Service", Segments: []PathSegment{{Name: ""}}},
		{Plane: "Service", Segments: []PathSegment{{Name: "bad/seg"}}},
		{Plane: "Service", Segments: []PathSegment{{Name: "bad:seg"}}},
		{Plane: "Service", Segments: []PathSegment{{Name: "@loc"}}},
	}

	for index, candidate := range invalid {
		if err := candidate.Validate(); err == nil {
			t.Fatalf("expected error for invalid path %d", index)
		}
	}
}

func TestStableNodeID_UsesServicePlane(t *testing.T) {
	service := ProjectionPath{
		Plane: ServicePlane,
		Segments: []PathSegment{
			{Name: "devices"},
			{Name: "boiler"},
		},
	}
	id, err := StableNodeID(service)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != NodeID("Service:/devices/boiler") {
		t.Fatalf("unexpected id: %s", id)
	}

	nonService := ProjectionPath{
		Plane: "Observability",
		Segments: []PathSegment{
			{Name: "devices"},
		},
	}
	if _, err := StableNodeID(nonService); err == nil {
		t.Fatalf("expected error for non-service canonical path")
	}
}

func TestNewNode_ServicePathInvariant(t *testing.T) {
	canonical := ProjectionPath{
		Plane: ServicePlane,
		Segments: []PathSegment{
			{Name: "devices"},
			{Name: "boiler"},
		},
	}
	path := ProjectionPath{
		Plane: "Observability",
		Segments: []PathSegment{
			{Name: "devices"},
			{Name: "boiler"},
		},
	}
	node, err := NewNode(path, canonical)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if node.ID != NodeID(canonical.String()) {
		t.Fatalf("unexpected id: %s", node.ID)
	}

	if _, err := NewNode(ProjectionPath{Plane: ServicePlane, Segments: []PathSegment{{Name: "other"}}}, canonical); err == nil {
		t.Fatalf("expected error for service path mismatch")
	}
}

func TestProjection_ValidatesPathsAndEdges(t *testing.T) {
	canonicalA := ProjectionPath{
		Plane: ServicePlane,
		Segments: []PathSegment{
			{Name: "devices"},
			{Name: "boiler"},
			{Name: "temp"},
		},
	}
	canonicalB := ProjectionPath{
		Plane: ServicePlane,
		Segments: []PathSegment{
			{Name: "devices"},
			{Name: "boiler"},
			{Name: "status"},
		},
	}

	nodeA, err := NewNode(ProjectionPath{
		Plane: "Observability",
		Segments: []PathSegment{
			{Name: "devices"},
			{Name: "boiler"},
			{Name: "temp"},
		},
	}, canonicalA)
	if err != nil {
		t.Fatalf("unexpected node error: %v", err)
	}

	nodeB, err := NewNode(ProjectionPath{
		Plane: "Observability",
		Segments: []PathSegment{
			{Name: "devices"},
			{Name: "boiler"},
			{Name: "status"},
		},
	}, canonicalB)
	if err != nil {
		t.Fatalf("unexpected node error: %v", err)
	}

	edge, err := NewEdge("Observability", nodeA.ID, nodeB.ID)
	if err != nil {
		t.Fatalf("unexpected edge error: %v", err)
	}

	projection, err := NewProjection("Observability", []Node{nodeA, nodeB}, []Edge{edge})
	if err != nil {
		t.Fatalf("unexpected projection error: %v", err)
	}
	if projection.Plane != "Observability" {
		t.Fatalf("unexpected projection plane: %s", projection.Plane)
	}

	duplicatePath := projection
	duplicatePath.Nodes = append(duplicatePath.Nodes, nodeA)
	if err := duplicatePath.Validate(); err == nil {
		t.Fatalf("expected error for duplicate path")
	}

	missingEdge := Projection{
		Plane: "Observability",
		Nodes: []Node{nodeA},
		Edges: []Edge{
			{ID: "Observability:missing->missing", From: "missing", To: "missing"},
		},
	}
	if err := missingEdge.Validate(); err == nil {
		t.Fatalf("expected error for missing edge nodes")
	}
}
