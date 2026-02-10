package registry

import (
	"errors"
	"fmt"
	"strings"
)

const ServicePlane = "Service"

var (
	ErrProjectionInvalidPlane   = errors.New("projection: invalid plane")
	ErrProjectionInvalidSegment = errors.New("projection: invalid path segment")
	ErrProjectionInvalidPath    = errors.New("projection: invalid path")
	ErrProjectionInvalidNode    = errors.New("projection: invalid node")
	ErrProjectionInvalidEdge    = errors.New("projection: invalid edge")
)

type PathSegment struct {
	Name     string
	Location bool
}

func (segment PathSegment) String() string {
	if segment.Location {
		return "@" + segment.Name
	}
	return segment.Name
}

func (segment PathSegment) validate() error {
	if strings.TrimSpace(segment.Name) == "" {
		return ErrProjectionInvalidSegment
	}
	if strings.ContainsAny(segment.Name, "/:") {
		return ErrProjectionInvalidSegment
	}
	if strings.HasPrefix(segment.Name, "@") {
		return ErrProjectionInvalidSegment
	}
	return nil
}

type ProjectionPath struct {
	Plane    string
	Segments []PathSegment
}

func (path ProjectionPath) String() string {
	if path.Plane == "" {
		return ""
	}
	builder := strings.Builder{}
	builder.WriteString(path.Plane)
	builder.WriteString(":/")
	for index, segment := range path.Segments {
		if index > 0 {
			builder.WriteString("/")
		}
		builder.WriteString(segment.String())
	}
	return builder.String()
}

func (path ProjectionPath) Validate() error {
	if strings.TrimSpace(path.Plane) == "" {
		return ErrProjectionInvalidPlane
	}
	if strings.ContainsAny(path.Plane, "/:") {
		return ErrProjectionInvalidPlane
	}
	for index, segment := range path.Segments {
		if err := segment.validate(); err != nil {
			return fmt.Errorf("segment %d: %w", index, errors.Join(err, ErrProjectionInvalidPath))
		}
	}
	return nil
}

type NodeID string
type EdgeID string

type Node struct {
	ID            NodeID
	Path          ProjectionPath
	CanonicalPath ProjectionPath
}

func StableNodeID(canonical ProjectionPath) (NodeID, error) {
	if err := canonical.Validate(); err != nil {
		return "", fmt.Errorf("canonical path: %w", err)
	}
	if canonical.Plane != ServicePlane {
		return "", fmt.Errorf("canonical plane %q: %w", canonical.Plane, ErrProjectionInvalidNode)
	}
	return NodeID(canonical.String()), nil
}

func NewNode(path ProjectionPath, canonical ProjectionPath) (Node, error) {
	if err := path.Validate(); err != nil {
		return Node{}, fmt.Errorf("path: %w", err)
	}
	id, err := StableNodeID(canonical)
	if err != nil {
		return Node{}, err
	}
	if path.Plane == ServicePlane && path.String() != canonical.String() {
		return Node{}, fmt.Errorf("service path mismatch: %w", ErrProjectionInvalidNode)
	}
	return Node{
		ID:            id,
		Path:          path,
		CanonicalPath: canonical,
	}, nil
}

type Edge struct {
	ID   EdgeID
	From NodeID
	To   NodeID
}

func StableEdgeID(plane string, from NodeID, to NodeID) (EdgeID, error) {
	if strings.TrimSpace(plane) == "" || strings.ContainsAny(plane, "/:") {
		return "", ErrProjectionInvalidPlane
	}
	if from == "" || to == "" {
		return "", ErrProjectionInvalidEdge
	}
	return EdgeID(fmt.Sprintf("%s:%s->%s", plane, from, to)), nil
}

func NewEdge(plane string, from NodeID, to NodeID) (Edge, error) {
	id, err := StableEdgeID(plane, from, to)
	if err != nil {
		return Edge{}, err
	}
	return Edge{ID: id, From: from, To: to}, nil
}

type Projection struct {
	Plane string
	Nodes []Node
	Edges []Edge
}

type CanonicalIndex struct {
	canonicalByID map[NodeID]ProjectionPath
	planePaths    map[string]map[NodeID]ProjectionPath
}

func (index CanonicalIndex) Canonical(id NodeID) (ProjectionPath, bool) {
	if index.canonicalByID == nil {
		return ProjectionPath{}, false
	}
	path, ok := index.canonicalByID[id]
	return path, ok
}

func (index CanonicalIndex) PlanePath(plane string, id NodeID) (ProjectionPath, bool) {
	if index.planePaths == nil {
		return ProjectionPath{}, false
	}
	planeMap, ok := index.planePaths[plane]
	if !ok {
		return ProjectionPath{}, false
	}
	path, ok := planeMap[id]
	return path, ok
}

func (index CanonicalIndex) PlanePaths(id NodeID) map[string]ProjectionPath {
	if index.planePaths == nil {
		return nil
	}
	out := make(map[string]ProjectionPath)
	for plane, planeMap := range index.planePaths {
		if path, ok := planeMap[id]; ok {
			out[plane] = path
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func NewProjection(plane string, nodes []Node, edges []Edge) (Projection, error) {
	projection := Projection{Plane: plane, Nodes: nodes, Edges: edges}
	if err := projection.Validate(); err != nil {
		return Projection{}, err
	}
	return projection, nil
}

func BuildCanonicalIndex(projections []Projection) (CanonicalIndex, error) {
	index := CanonicalIndex{
		canonicalByID: make(map[NodeID]ProjectionPath),
		planePaths:    make(map[string]map[NodeID]ProjectionPath),
	}
	if len(projections) == 0 {
		return index, nil
	}

	hasService := false
	serviceIDs := make(map[NodeID]struct{})

	for projectionIndex, projection := range projections {
		if err := projection.Validate(); err != nil {
			return CanonicalIndex{}, fmt.Errorf("projection %d: %w", projectionIndex, err)
		}

		plane := projection.Plane
		if plane == ServicePlane {
			hasService = true
		}

		planeMap, ok := index.planePaths[plane]
		if !ok {
			planeMap = make(map[NodeID]ProjectionPath, len(projection.Nodes))
			index.planePaths[plane] = planeMap
		}

		for nodeIndex, node := range projection.Nodes {
			expectedID, err := StableNodeID(node.CanonicalPath)
			if err != nil {
				return CanonicalIndex{}, fmt.Errorf("projection %d node %d canonical: %w", projectionIndex, nodeIndex, err)
			}
			if node.ID != expectedID {
				return CanonicalIndex{}, fmt.Errorf("projection %d node %d id mismatch: %w", projectionIndex, nodeIndex, ErrProjectionInvalidNode)
			}
			if plane == ServicePlane && node.Path.String() != node.CanonicalPath.String() {
				return CanonicalIndex{}, fmt.Errorf("projection %d node %d service path mismatch: %w", projectionIndex, nodeIndex, ErrProjectionInvalidNode)
			}

			if existing, ok := index.canonicalByID[node.ID]; ok {
				if existing.String() != node.CanonicalPath.String() {
					return CanonicalIndex{}, fmt.Errorf("projection %d node %d id collision: %w", projectionIndex, nodeIndex, ErrProjectionInvalidNode)
				}
			} else {
				index.canonicalByID[node.ID] = node.CanonicalPath
			}

			if existing, ok := planeMap[node.ID]; ok {
				if existing.String() != node.Path.String() {
					return CanonicalIndex{}, fmt.Errorf("projection %d node %d plane path mismatch: %w", projectionIndex, nodeIndex, ErrProjectionInvalidNode)
				}
			} else {
				planeMap[node.ID] = node.Path
			}

			if plane == ServicePlane {
				serviceIDs[node.ID] = struct{}{}
			}
		}
	}

	if hasService {
		for plane, planeMap := range index.planePaths {
			if plane == ServicePlane {
				continue
			}
			for id := range planeMap {
				if _, ok := serviceIDs[id]; !ok {
					return CanonicalIndex{}, fmt.Errorf("plane %q node %s missing service canonical: %w", plane, id, ErrProjectionInvalidNode)
				}
			}
		}
	}

	return index, nil
}

func (projection Projection) Validate() error {
	if strings.TrimSpace(projection.Plane) == "" || strings.ContainsAny(projection.Plane, "/:") {
		return ErrProjectionInvalidPlane
	}

	paths := make(map[string]struct{}, len(projection.Nodes))
	nodesByID := make(map[NodeID]ProjectionPath, len(projection.Nodes))
	for index, node := range projection.Nodes {
		if node.ID == "" {
			return fmt.Errorf("node %d missing id: %w", index, ErrProjectionInvalidNode)
		}
		if err := node.Path.Validate(); err != nil {
			return fmt.Errorf("node %d path: %w", index, err)
		}
		if node.Path.Plane != projection.Plane {
			return fmt.Errorf("node %d plane %q: %w", index, node.Path.Plane, ErrProjectionInvalidNode)
		}
		if err := node.CanonicalPath.Validate(); err != nil {
			return fmt.Errorf("node %d canonical: %w", index, err)
		}
		if node.CanonicalPath.Plane != ServicePlane {
			return fmt.Errorf("node %d canonical plane %q: %w", index, node.CanonicalPath.Plane, ErrProjectionInvalidNode)
		}
		pathKey := node.Path.String()
		if _, ok := paths[pathKey]; ok {
			return fmt.Errorf("duplicate path %q: %w", pathKey, ErrProjectionInvalidNode)
		}
		paths[pathKey] = struct{}{}
		if existing, ok := nodesByID[node.ID]; ok {
			if existing.String() != node.CanonicalPath.String() {
				return fmt.Errorf("node %d id collision: %w", index, ErrProjectionInvalidNode)
			}
		} else {
			nodesByID[node.ID] = node.CanonicalPath
		}
	}

	for index, edge := range projection.Edges {
		if edge.ID == "" {
			return fmt.Errorf("edge %d missing id: %w", index, ErrProjectionInvalidEdge)
		}
		if edge.From == "" || edge.To == "" {
			return fmt.Errorf("edge %d missing endpoints: %w", index, ErrProjectionInvalidEdge)
		}
		if _, ok := nodesByID[edge.From]; !ok {
			return fmt.Errorf("edge %d missing from node: %w", index, ErrProjectionInvalidEdge)
		}
		if _, ok := nodesByID[edge.To]; !ok {
			return fmt.Errorf("edge %d missing to node: %w", index, ErrProjectionInvalidEdge)
		}
	}
	return nil
}
