package registry

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	ebuserrors "github.com/Project-Helianthus/helianthus-ebusgo/errors"
	"github.com/Project-Helianthus/helianthus-ebusreg/schema"
)

// EntryIterator is the minimal registry surface required by projection helpers.
type EntryIterator interface {
	Iterate(func(DeviceEntry) bool)
}

// ServiceDeviceView is a normalized device projection for shared service-layer use.
type ServiceDeviceView struct {
	Address         byte
	Addresses       []byte
	Manufacturer    string
	DeviceID        string
	SerialNumber    string
	MacAddress      string
	SoftwareVersion string
	HardwareVersion string
	Planes          []ServicePlaneView
}

// ServicePlaneView is a normalized plane projection for shared service-layer use.
type ServicePlaneView struct {
	Name    string
	Methods []ServiceMethodView
}

// ServiceMethodView is a normalized method projection for shared service-layer use.
type ServiceMethodView struct {
	Name             string
	ReadOnly         bool
	Primary          byte
	Secondary        byte
	Metadata         MethodMetadata
	ResponseSelector schema.SchemaSelector
}

// ProjectRegistryDevices projects every entry of iter into service views,
// sorted by address and then by identity fields. The first entry that fails
// to project aborts the walk and its error is returned with no views. A nil
// iter wraps the ebusgo ErrInvalidPayload.
func ProjectRegistryDevices(iter EntryIterator) ([]ServiceDeviceView, error) {
	if iter == nil {
		return nil, fmt.Errorf("project registry devices: %w", ebuserrors.ErrInvalidPayload)
	}

	out := make([]ServiceDeviceView, 0)
	var projectionErr error
	iter.Iterate(func(entry DeviceEntry) bool {
		projected, err := ProjectDeviceEntry(entry)
		if err != nil {
			projectionErr = err
			return false
		}
		out = append(out, projected)
		return true
	})
	if projectionErr != nil {
		return nil, projectionErr
	}
	sortServiceDevices(out)
	return out, nil
}

// ProjectDeviceEntry copies the identity and planes of entry into a
// ServiceDeviceView. Address is the display address; Addresses lists it
// first, followed by the remaining addresses in ascending order. Planes are
// sorted by name trimmed of surrounding whitespace and case-folded, then by
// exact name; their methods are ordered as ProjectPlane orders them. A nil
// entry, plane or method, including an interface holding a nil pointer, or a
// method without a template, wraps ErrInvalidPayload.
func ProjectDeviceEntry(entry DeviceEntry) (ServiceDeviceView, error) {
	if isNilProjectionInput(entry) {
		return ServiceDeviceView{}, fmt.Errorf("project device entry: %w", ebuserrors.ErrInvalidPayload)
	}

	planes := entry.Planes()
	projectedPlanes := make([]ServicePlaneView, 0, len(planes))
	for index, plane := range planes {
		projected, err := ProjectPlane(plane)
		if err != nil {
			return ServiceDeviceView{}, fmt.Errorf("project device plane %d: %w", index, err)
		}
		projectedPlanes = append(projectedPlanes, projected)
	}
	sortServicePlanes(projectedPlanes)

	return ServiceDeviceView{
		Address:         entry.PrimaryDisplayAddress(),
		Addresses:       normalizeProjectedAddresses(entry.PrimaryDisplayAddress(), entry.Addresses()),
		Manufacturer:    entry.Manufacturer(),
		DeviceID:        entry.DeviceID(),
		SerialNumber:    entry.SerialNumber(),
		MacAddress:      entry.MacAddress(),
		SoftwareVersion: entry.SoftwareVersion(),
		HardwareVersion: entry.HardwareVersion(),
		Planes:          projectedPlanes,
	}, nil
}

// ProjectPlane copies plane into a ServicePlaneView with its methods sorted
// by name trimmed of surrounding whitespace and case-folded, then by exact
// name, then by command bytes (primary, then secondary). A nil plane,
// including an interface holding a nil pointer, wraps ErrInvalidPayload; a
// method that fails ProjectMethod returns its error.
func ProjectPlane(plane Plane) (ServicePlaneView, error) {
	if isNilProjectionInput(plane) {
		return ServicePlaneView{}, fmt.Errorf("project plane: %w", ebuserrors.ErrInvalidPayload)
	}

	methods := plane.Methods()
	projectedMethods := make([]ServiceMethodView, 0, len(methods))
	for index, method := range methods {
		projected, err := ProjectMethod(method)
		if err != nil {
			return ServicePlaneView{}, fmt.Errorf("project method %d: %w", index, err)
		}
		projectedMethods = append(projectedMethods, projected)
	}
	sortServiceMethods(projectedMethods)

	return ServicePlaneView{
		Name:    plane.Name(),
		Methods: projectedMethods,
	}, nil
}

// ProjectMethod copies method into a ServiceMethodView, including its
// resolved MethodMetadata. A nil method, including an interface holding a
// nil pointer, or a nil Template wraps ErrInvalidPayload.
func ProjectMethod(method Method) (ServiceMethodView, error) {
	if isNilProjectionInput(method) {
		return ServiceMethodView{}, fmt.Errorf("project method: %w", ebuserrors.ErrInvalidPayload)
	}
	template := method.Template()
	if template == nil {
		return ServiceMethodView{}, fmt.Errorf("project method %q missing template: %w", method.Name(), ebuserrors.ErrInvalidPayload)
	}

	return ServiceMethodView{
		Name:             method.Name(),
		ReadOnly:         method.ReadOnly(),
		Primary:          template.Primary(),
		Secondary:        template.Secondary(),
		Metadata:         ResolveMethodMetadata(method),
		ResponseSelector: method.ResponseSchema(),
	}, nil
}

func isNilProjectionInput(value any) bool {
	if value == nil {
		return true
	}
	reflected := reflect.ValueOf(value)
	switch reflected.Kind() {
	case reflect.Pointer, reflect.Map, reflect.Slice, reflect.Func, reflect.Chan, reflect.Interface:
		return reflected.IsNil()
	default:
		return false
	}
}

func normalizeProjectedAddresses(primary byte, aliases []byte) []byte {
	out := make([]byte, 0, len(aliases)+1)
	appendUniqueAddress := func(address byte) {
		for _, existing := range out {
			if existing == address {
				return
			}
		}
		out = append(out, address)
	}

	appendUniqueAddress(primary)
	aliasesCopy := append([]byte(nil), aliases...)
	sort.Slice(aliasesCopy, func(i, j int) bool {
		return aliasesCopy[i] < aliasesCopy[j]
	})
	for _, address := range aliasesCopy {
		appendUniqueAddress(address)
	}

	if len(out) == 0 {
		return nil
	}
	return out
}

func sortServiceDevices(devices []ServiceDeviceView) {
	sort.Slice(devices, func(i, j int) bool {
		left := devices[i]
		right := devices[j]
		if left.Address != right.Address {
			return left.Address < right.Address
		}
		if compareFolded(left.Manufacturer, right.Manufacturer) != 0 {
			return compareFolded(left.Manufacturer, right.Manufacturer) < 0
		}
		if compareFolded(left.DeviceID, right.DeviceID) != 0 {
			return compareFolded(left.DeviceID, right.DeviceID) < 0
		}
		if compareFolded(left.HardwareVersion, right.HardwareVersion) != 0 {
			return compareFolded(left.HardwareVersion, right.HardwareVersion) < 0
		}
		return compareFolded(left.SerialNumber, right.SerialNumber) < 0
	})
	for index := range devices {
		sortServicePlanes(devices[index].Planes)
	}
}

func sortServicePlanes(planes []ServicePlaneView) {
	sort.Slice(planes, func(i, j int) bool {
		if compareFolded(planes[i].Name, planes[j].Name) != 0 {
			return compareFolded(planes[i].Name, planes[j].Name) < 0
		}
		return planes[i].Name < planes[j].Name
	})
	for index := range planes {
		sortServiceMethods(planes[index].Methods)
	}
}

func sortServiceMethods(methods []ServiceMethodView) {
	sort.Slice(methods, func(i, j int) bool {
		if compareFolded(methods[i].Name, methods[j].Name) != 0 {
			return compareFolded(methods[i].Name, methods[j].Name) < 0
		}
		if methods[i].Name != methods[j].Name {
			return methods[i].Name < methods[j].Name
		}
		if methods[i].Primary != methods[j].Primary {
			return methods[i].Primary < methods[j].Primary
		}
		return methods[i].Secondary < methods[j].Secondary
	})
}

func compareFolded(left string, right string) int {
	leftFolded := strings.ToLower(strings.TrimSpace(left))
	rightFolded := strings.ToLower(strings.TrimSpace(right))
	if leftFolded < rightFolded {
		return -1
	}
	if leftFolded > rightFolded {
		return 1
	}
	return 0
}
