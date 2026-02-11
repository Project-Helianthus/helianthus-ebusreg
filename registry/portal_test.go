package registry

import (
	"errors"
	"testing"
)

func TestPortalIndex_QueryByCanonical(t *testing.T) {
	canonicalRoot := ProjectionPath{
		Plane: ServicePlane,
		Segments: []PathSegment{
			{Name: "devices"},
			{Name: "boiler"},
		},
	}
	canonicalChild := ProjectionPath{
		Plane: ServicePlane,
		Segments: []PathSegment{
			{Name: "devices"},
			{Name: "boiler"},
			{Name: "status"},
		},
	}
	rootPath := ProjectionPath{
		Plane: "Observability",
		Segments: []PathSegment{
			{Name: "devices"},
			{Name: "boiler"},
		},
	}
	childPath := ProjectionPath{
		Plane: "Observability",
		Segments: []PathSegment{
			{Name: "devices"},
			{Name: "boiler"},
			{Name: "status"},
		},
	}

	rootNode, err := NewNode(rootPath, canonicalRoot)
	if err != nil {
		t.Fatalf("unexpected root node error: %v", err)
	}
	childNode, err := NewNode(childPath, canonicalChild)
	if err != nil {
		t.Fatalf("unexpected child node error: %v", err)
	}

	edge, err := NewEdge("Observability", rootNode.ID, childNode.ID)
	if err != nil {
		t.Fatalf("unexpected edge error: %v", err)
	}

	projection, err := NewProjection("Observability", []Node{rootNode, childNode}, []Edge{edge})
	if err != nil {
		t.Fatalf("unexpected projection error: %v", err)
	}

	portal, err := NewPortalIndex([]Projection{projection})
	if err != nil {
		t.Fatalf("unexpected portal index error: %v", err)
	}

	found, ok, err := portal.NodeByCanonical("Observability", canonicalChild)
	if err != nil {
		t.Fatalf("unexpected canonical lookup error: %v", err)
	}
	if !ok || found.ID != childNode.ID {
		t.Fatalf("unexpected canonical lookup result: %v %v", ok, found.ID)
	}

	plane, ok := portal.PlaneIndex("Observability")
	if !ok || plane.Plane() != "Observability" {
		t.Fatalf("unexpected plane index lookup")
	}

	byPath, ok, err := plane.NodeByPath(childPath)
	if err != nil {
		t.Fatalf("unexpected path lookup error: %v", err)
	}
	if !ok || byPath.ID != childNode.ID {
		t.Fatalf("unexpected path lookup result: %v %v", ok, byPath.ID)
	}

	edgeOut, ok := portal.EdgeByID("Observability", edge.ID)
	if !ok || edgeOut.ID != edge.ID {
		t.Fatalf("unexpected edge lookup result")
	}

	fromEdges := plane.EdgesFrom(rootNode.ID)
	if len(fromEdges) != 1 || fromEdges[0].ID != edge.ID {
		t.Fatalf("unexpected edges from result: %v", fromEdges)
	}
	toEdges := plane.EdgesTo(childNode.ID)
	if len(toEdges) != 1 || toEdges[0].ID != edge.ID {
		t.Fatalf("unexpected edges to result: %v", toEdges)
	}
	if plane.EdgeCount() != 1 {
		t.Fatalf("unexpected edge count: %d", plane.EdgeCount())
	}
}

func TestPortalIndex_NodeByCanonicalRejectsNonServicePlane(t *testing.T) {
	canonical := ProjectionPath{
		Plane: ServicePlane,
		Segments: []PathSegment{
			{Name: "devices"},
		},
	}
	path := ProjectionPath{
		Plane: "Observability",
		Segments: []PathSegment{
			{Name: "devices"},
		},
	}
	node, err := NewNode(path, canonical)
	if err != nil {
		t.Fatalf("unexpected node error: %v", err)
	}
	projection, err := NewProjection("Observability", []Node{node}, nil)
	if err != nil {
		t.Fatalf("unexpected projection error: %v", err)
	}
	portal, err := NewPortalIndex([]Projection{projection})
	if err != nil {
		t.Fatalf("unexpected portal error: %v", err)
	}

	_, _, err = portal.NodeByCanonical("Observability", ProjectionPath{
		Plane: "Observability",
		Segments: []PathSegment{
			{Name: "devices"},
		},
	})
	if err == nil {
		t.Fatalf("expected error for non-service canonical plane")
	}
}

func TestPortalIndex_DuplicatePlane(t *testing.T) {
	projection, err := NewProjection("Observability", nil, nil)
	if err != nil {
		t.Fatalf("unexpected projection error: %v", err)
	}
	if _, err := NewPortalIndex([]Projection{projection, projection}); err == nil || !errors.Is(err, ErrPortalDuplicatePlane) {
		t.Fatalf("expected duplicate plane error, got %v", err)
	}
}

func TestPlaneIndex_DuplicateNodeID(t *testing.T) {
	canonical := ProjectionPath{
		Plane: ServicePlane,
		Segments: []PathSegment{
			{Name: "devices"},
			{Name: "boiler"},
		},
	}
	nodeA, err := NewNode(ProjectionPath{
		Plane: "Observability",
		Segments: []PathSegment{
			{Name: "devices"},
			{Name: "boiler"},
		},
	}, canonical)
	if err != nil {
		t.Fatalf("unexpected nodeA error: %v", err)
	}
	nodeB, err := NewNode(ProjectionPath{
		Plane: "Observability",
		Segments: []PathSegment{
			{Name: "devices"},
			{Name: "boiler"},
			{Name: "alt"},
		},
	}, canonical)
	if err != nil {
		t.Fatalf("unexpected nodeB error: %v", err)
	}

	projection, err := NewProjection("Observability", []Node{nodeA, nodeB}, nil)
	if err != nil {
		t.Fatalf("unexpected projection error: %v", err)
	}
	if _, err := NewPlaneIndex(projection); err == nil || !errors.Is(err, ErrPortalDuplicateNode) {
		t.Fatalf("expected duplicate node error, got %v", err)
	}
}

func TestPlaneIndex_DuplicateEdgeID(t *testing.T) {
	canonicalA := ProjectionPath{
		Plane: ServicePlane,
		Segments: []PathSegment{
			{Name: "devices"},
			{Name: "boiler"},
			{Name: "a"},
		},
	}
	canonicalB := ProjectionPath{
		Plane: ServicePlane,
		Segments: []PathSegment{
			{Name: "devices"},
			{Name: "boiler"},
			{Name: "b"},
		},
	}
	nodeA, err := NewNode(ProjectionPath{
		Plane: "Observability",
		Segments: []PathSegment{
			{Name: "devices"},
			{Name: "boiler"},
			{Name: "a"},
		},
	}, canonicalA)
	if err != nil {
		t.Fatalf("unexpected nodeA error: %v", err)
	}
	nodeB, err := NewNode(ProjectionPath{
		Plane: "Observability",
		Segments: []PathSegment{
			{Name: "devices"},
			{Name: "boiler"},
			{Name: "b"},
		},
	}, canonicalB)
	if err != nil {
		t.Fatalf("unexpected nodeB error: %v", err)
	}

	edgeA, err := NewEdge("Observability", nodeA.ID, nodeB.ID)
	if err != nil {
		t.Fatalf("unexpected edgeA error: %v", err)
	}
	edgeB, err := NewEdge("Observability", nodeA.ID, nodeB.ID)
	if err != nil {
		t.Fatalf("unexpected edgeB error: %v", err)
	}

	projection, err := NewProjection("Observability", []Node{nodeA, nodeB}, []Edge{edgeA, edgeB})
	if err != nil {
		t.Fatalf("unexpected projection error: %v", err)
	}
	if _, err := NewPlaneIndex(projection); err == nil || !errors.Is(err, ErrPortalDuplicateEdge) {
		t.Fatalf("expected duplicate edge error, got %v", err)
	}
}

func TestPortalIndexForEntry_Nil(t *testing.T) {
	if _, err := PortalIndexForEntry(nil); err == nil || !errors.Is(err, ErrPortalInvalidEntry) {
		t.Fatalf("expected invalid entry error, got %v", err)
	}
}

func TestPortalIndexForEntry_TypedNil(t *testing.T) {
	var entry DeviceEntry = (*deviceEntry)(nil)
	if _, err := PortalIndexForEntry(entry); err == nil || !errors.Is(err, ErrPortalInvalidEntry) {
		t.Fatalf("expected invalid entry error, got %v", err)
	}
}
