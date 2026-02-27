package system

import (
	"fmt"
	"strings"

	"github.com/d3vi1/helianthus-ebusreg/registry"
)

const (
	projectionPlaneObservability = "Observability"
	projectionPlaneDebug         = "Debug"
	projectionPlaneSystem        = "System"
	projectionPlaneHeating       = "Heating"
	projectionPlaneDHW           = "DHW"
	projectionPlaneSolar         = "Solar"
)

func (Provider) CreateProjections(info registry.DeviceInfo, planes []registry.Plane) []registry.Projection {
	if !shouldCreateProjections(planes) {
		return nil
	}

	base := projectionBaseSegments(info)
	canonicalRoot := registry.ProjectionPath{Plane: registry.ServicePlane, Segments: base}

	rootService, ok := newNode(registry.ServicePlane, base, canonicalRoot)
	if !ok {
		return nil
	}
	rootObservability, ok := newNode(projectionPlaneObservability, base, canonicalRoot)
	if !ok {
		return nil
	}
	rootDebug, ok := newNode(projectionPlaneDebug, base, canonicalRoot)
	if !ok {
		return nil
	}

	operationalCanonical := registry.ProjectionPath{
		Plane:    registry.ServicePlane,
		Segments: appendSegment(base, methodSegment(methodGetOperationalData)),
	}
	operationalService, ok := newNode(registry.ServicePlane, operationalCanonical.Segments, operationalCanonical)
	if !ok {
		return nil
	}
	operationalObservability, ok := newNode(projectionPlaneObservability, operationalCanonical.Segments, operationalCanonical)
	if !ok {
		return nil
	}

	registerCanonical := registry.ProjectionPath{
		Plane:    registry.ServicePlane,
		Segments: appendSegment(base, methodSegment(methodGetRegister)),
	}
	registerDebug, ok := newNode(projectionPlaneDebug, appendSegment(base, registry.PathSegment{Name: "register@b509"}), registerCanonical)
	if !ok {
		return nil
	}

	extRegisterCanonical := registry.ProjectionPath{
		Plane:    registry.ServicePlane,
		Segments: appendSegment(base, methodSegment(methodGetExtRegister)),
	}
	extRegisterDebug, ok := newNode(projectionPlaneDebug, appendSegment(base, registry.PathSegment{Name: "register@b524"}), extRegisterCanonical)
	if !ok {
		return nil
	}

	projections := make([]registry.Projection, 0, 7)

	serviceEdges := make([]registry.Edge, 0, 1)
	if edge, err := registry.NewEdge(registry.ServicePlane, rootService.ID, operationalService.ID); err == nil {
		serviceEdges = append(serviceEdges, edge)
	}
	if projection, err := registry.NewProjection(registry.ServicePlane, []registry.Node{rootService, operationalService}, serviceEdges); err == nil {
		projections = append(projections, projection)
	}

	observabilityEdges := make([]registry.Edge, 0, 1)
	if edge, err := registry.NewEdge(projectionPlaneObservability, rootObservability.ID, operationalObservability.ID); err == nil {
		observabilityEdges = append(observabilityEdges, edge)
	}
	if projection, err := registry.NewProjection(projectionPlaneObservability, []registry.Node{rootObservability, operationalObservability}, observabilityEdges); err == nil {
		projections = append(projections, projection)
	}

	debugEdges := make([]registry.Edge, 0, 2)
	if edge, err := registry.NewEdge(projectionPlaneDebug, rootDebug.ID, registerDebug.ID); err == nil {
		debugEdges = append(debugEdges, edge)
	}
	if edge, err := registry.NewEdge(projectionPlaneDebug, rootDebug.ID, extRegisterDebug.ID); err == nil {
		debugEdges = append(debugEdges, edge)
	}
	if projection, err := registry.NewProjection(projectionPlaneDebug, []registry.Node{rootDebug, registerDebug, extRegisterDebug}, debugEdges); err == nil {
		projections = append(projections, projection)
	}

	planeSet := planeNameSet(planes)
	projections = append(projections, rootProjectionIfPresent(planeSet, projectionPlaneSystem, canonicalRoot, base)...)
	projections = append(projections, rootProjectionIfPresent(planeSet, projectionPlaneHeating, canonicalRoot, base)...)
	projections = append(projections, rootProjectionIfPresent(planeSet, projectionPlaneDHW, canonicalRoot, base)...)
	projections = append(projections, rootProjectionIfPresent(planeSet, projectionPlaneSolar, canonicalRoot, base)...)

	return projections
}

func rootProjectionIfPresent(planes map[string]struct{}, plane string, canonical registry.ProjectionPath, base []registry.PathSegment) []registry.Projection {
	if _, ok := planes[strings.ToLower(plane)]; !ok {
		return nil
	}

	node, ok := newNode(plane, base, canonical)
	if !ok {
		return nil
	}
	projection, err := registry.NewProjection(plane, []registry.Node{node}, nil)
	if err != nil {
		return nil
	}
	return []registry.Projection{projection}
}

func planeNameSet(planes []registry.Plane) map[string]struct{} {
	out := make(map[string]struct{}, len(planes))
	for _, plane := range planes {
		if plane == nil {
			continue
		}
		name := strings.TrimSpace(plane.Name())
		if name == "" {
			continue
		}
		out[strings.ToLower(name)] = struct{}{}
	}
	return out
}

func shouldCreateProjections(planes []registry.Plane) bool {
	return len(planes) > 0
}

func normalizeDeviceID(id string) string {
	if strings.TrimSpace(id) == "" {
		return ""
	}
	normalized := strings.ToUpper(strings.TrimSpace(id))
	normalized = strings.ReplaceAll(normalized, "_", "")
	normalized = strings.ReplaceAll(normalized, " ", "")
	normalized = strings.ReplaceAll(normalized, "-", "")
	return normalized
}

func projectionBaseSegments(info registry.DeviceInfo) []registry.PathSegment {
	addr := fmt.Sprintf("%02X", info.Address)
	deviceID := strings.TrimSpace(info.DeviceID)
	if deviceID == "" {
		deviceID = addr
	}
	return []registry.PathSegment{
		{Name: "ebus"},
		{Name: fmt.Sprintf("addr@%s", addr)},
		{Name: fmt.Sprintf("device@%s", deviceID)},
	}
}

func methodSegment(name string) registry.PathSegment {
	return registry.PathSegment{Name: "method@" + name}
}

func appendSegment(segments []registry.PathSegment, segment registry.PathSegment) []registry.PathSegment {
	out := make([]registry.PathSegment, 0, len(segments)+1)
	out = append(out, segments...)
	out = append(out, segment)
	return out
}

func newNode(plane string, segments []registry.PathSegment, canonical registry.ProjectionPath) (registry.Node, bool) {
	path := registry.ProjectionPath{Plane: plane, Segments: segments}
	node, err := registry.NewNode(path, canonical)
	if err != nil {
		return registry.Node{}, false
	}
	return node, true
}
