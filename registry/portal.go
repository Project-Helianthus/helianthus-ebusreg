package registry

import (
	"errors"
	"fmt"
)

var (
	ErrPortalDuplicatePlane = errors.New("portal: duplicate plane")
	ErrPortalDuplicateNode  = errors.New("portal: duplicate node")
	ErrPortalDuplicateEdge  = errors.New("portal: duplicate edge")
	ErrPortalInvalidEntry   = errors.New("portal: invalid entry")
)

// PortalIndex provides plane-scoped projection lookups for a single device entry.
//
// Example:
//
//	portal, err := NewPortalIndex(entry.Projections())
//	if err != nil {
//	    return err
//	}
//	node, ok, err := portal.NodeByCanonical("Observability", ProjectionPath{
//	    Plane: ServicePlane,
//	    Segments: []PathSegment{
//	        {Name: "devices"},
//	        {Name: "boiler"},
//	    },
//	})
//	if err != nil {
//	    return err
//	}
//	if ok {
//	    _ = node.Path
//	}
type PortalIndex struct {
	planes map[string]*PlaneIndex
}

// PortalIndexForEntry builds a portal index from the entry projections.
func PortalIndexForEntry(entry DeviceEntry) (PortalIndex, error) {
	if entry == nil {
		return PortalIndex{}, ErrPortalInvalidEntry
	}
	return NewPortalIndex(entry.Projections())
}

// PlaneIndex returns the plane-specific index if present.
//
// Example:
//
//	plane, ok := portal.PlaneIndex("Observability")
//	if ok {
//	    _ = plane.Plane()
//	}
func (portal PortalIndex) PlaneIndex(plane string) (*PlaneIndex, bool) {
	if portal.planes == nil {
		return nil, false
	}
	index, ok := portal.planes[plane]
	return index, ok
}

// NodeByCanonical looks up a node by plane and canonical Service path.
//
// Example:
//
//	node, ok, err := portal.NodeByCanonical("Observability", ProjectionPath{
//	    Plane: ServicePlane,
//	    Segments: []PathSegment{
//	        {Name: "devices"},
//	        {Name: "boiler"},
//	    },
//	})
func (portal PortalIndex) NodeByCanonical(plane string, canonical ProjectionPath) (Node, bool, error) {
	index, ok := portal.PlaneIndex(plane)
	if !ok {
		return Node{}, false, nil
	}
	return index.NodeByCanonical(canonical)
}

// EdgeByID looks up an edge by plane and edge ID.
func (portal PortalIndex) EdgeByID(plane string, edgeID EdgeID) (Edge, bool) {
	index, ok := portal.PlaneIndex(plane)
	if !ok {
		return Edge{}, false
	}
	return index.EdgeByID(edgeID)
}

// NewPortalIndex builds a read-only portal index from projections.
//
// Example:
//
//	portal, err := NewPortalIndex(entry.Projections())
//	if err != nil {
//	    return err
//	}
func NewPortalIndex(projections []Projection) (PortalIndex, error) {
	if len(projections) == 0 {
		return PortalIndex{}, nil
	}
	planes := make(map[string]*PlaneIndex, len(projections))
	for _, projection := range projections {
		if _, exists := planes[projection.Plane]; exists {
			return PortalIndex{}, fmt.Errorf("portal plane %q: %w", projection.Plane, ErrPortalDuplicatePlane)
		}
		index, err := NewPlaneIndex(projection)
		if err != nil {
			return PortalIndex{}, fmt.Errorf("portal plane %q: %w", projection.Plane, err)
		}
		planes[projection.Plane] = &index
	}
	return PortalIndex{planes: planes}, nil
}

// PlaneIndex provides indexed access to nodes and edges within a projection plane.
//
// Example:
//
//	plane, ok := portal.PlaneIndex("Observability")
//	if ok {
//	    node, ok, err := plane.NodeByCanonical(ProjectionPath{
//	        Plane: ServicePlane,
//	        Segments: []PathSegment{
//	            {Name: "devices"},
//	            {Name: "boiler"},
//	        },
//	    })
//	    _ = node
//	    _ = err
//	}
type PlaneIndex struct {
	plane      string
	nodesByID  map[NodeID]Node
	paths      map[string]Node
	edgesByID  map[EdgeID]Edge
	edgesFrom  map[NodeID][]Edge
	edgesTo    map[NodeID][]Edge
	edgeCounts int
}

// Plane returns the plane name for this index.
func (index *PlaneIndex) Plane() string {
	if index == nil {
		return ""
	}
	return index.plane
}

// NodeByCanonical looks up a node by canonical Service path.
func (index *PlaneIndex) NodeByCanonical(canonical ProjectionPath) (Node, bool, error) {
	if index == nil {
		return Node{}, false, nil
	}
	id, err := StableNodeID(canonical)
	if err != nil {
		return Node{}, false, err
	}
	node, ok := index.NodeByID(id)
	return node, ok, nil
}

// NodeByID looks up a node by stable node ID.
func (index *PlaneIndex) NodeByID(id NodeID) (Node, bool) {
	if index == nil {
		return Node{}, false
	}
	node, ok := index.nodesByID[id]
	return node, ok
}

// NodeByPath looks up a node by its plane-specific path.
func (index *PlaneIndex) NodeByPath(path ProjectionPath) (Node, bool, error) {
	if index == nil {
		return Node{}, false, nil
	}
	if err := path.Validate(); err != nil {
		return Node{}, false, err
	}
	if path.Plane != index.plane {
		return Node{}, false, ErrProjectionInvalidPlane
	}
	node, ok := index.paths[path.String()]
	return node, ok, nil
}

// EdgeByID looks up an edge by its stable ID.
func (index *PlaneIndex) EdgeByID(id EdgeID) (Edge, bool) {
	if index == nil {
		return Edge{}, false
	}
	edge, ok := index.edgesByID[id]
	return edge, ok
}

// EdgesFrom returns edges originating from the provided node ID.
func (index *PlaneIndex) EdgesFrom(id NodeID) []Edge {
	if index == nil {
		return nil
	}
	edges := index.edgesFrom[id]
	if len(edges) == 0 {
		return nil
	}
	out := make([]Edge, len(edges))
	copy(out, edges)
	return out
}

// EdgesTo returns edges targeting the provided node ID.
func (index *PlaneIndex) EdgesTo(id NodeID) []Edge {
	if index == nil {
		return nil
	}
	edges := index.edgesTo[id]
	if len(edges) == 0 {
		return nil
	}
	out := make([]Edge, len(edges))
	copy(out, edges)
	return out
}

// EdgeCount returns the number of edges indexed for the plane.
func (index *PlaneIndex) EdgeCount() int {
	if index == nil {
		return 0
	}
	return index.edgeCounts
}

// NewPlaneIndex builds an index from a single projection.
func NewPlaneIndex(projection Projection) (PlaneIndex, error) {
	if err := projection.Validate(); err != nil {
		return PlaneIndex{}, err
	}
	nodesByID := make(map[NodeID]Node, len(projection.Nodes))
	paths := make(map[string]Node, len(projection.Nodes))
	for _, node := range projection.Nodes {
		if _, exists := nodesByID[node.ID]; exists {
			return PlaneIndex{}, fmt.Errorf("node %s: %w", node.ID, ErrPortalDuplicateNode)
		}
		nodesByID[node.ID] = node
		paths[node.Path.String()] = node
	}

	edgesByID := make(map[EdgeID]Edge, len(projection.Edges))
	edgesFrom := make(map[NodeID][]Edge)
	edgesTo := make(map[NodeID][]Edge)
	for _, edge := range projection.Edges {
		if _, exists := edgesByID[edge.ID]; exists {
			return PlaneIndex{}, fmt.Errorf("edge %s: %w", edge.ID, ErrPortalDuplicateEdge)
		}
		edgesByID[edge.ID] = edge
		edgesFrom[edge.From] = append(edgesFrom[edge.From], edge)
		edgesTo[edge.To] = append(edgesTo[edge.To], edge)
	}

	return PlaneIndex{
		plane:      projection.Plane,
		nodesByID:  nodesByID,
		paths:      paths,
		edgesByID:  edgesByID,
		edgesFrom:  edgesFrom,
		edgesTo:    edgesTo,
		edgeCounts: len(projection.Edges),
	}, nil
}
