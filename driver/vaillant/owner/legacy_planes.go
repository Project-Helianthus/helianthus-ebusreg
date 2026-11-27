package owner

import (
	"context"
	"errors"
	"sync"

	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
	"github.com/Project-Helianthus/helianthus-ebusgo/types"

	catalogvaillant "github.com/Project-Helianthus/helianthus-ebusreg/catalog/vaillant"
	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/output"
	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
	"github.com/Project-Helianthus/helianthus-ebusreg/vaillant/dhw"
	"github.com/Project-Helianthus/helianthus-ebusreg/vaillant/heating"
	"github.com/Project-Helianthus/helianthus-ebusreg/vaillant/solar"
)

const (
	legacyNamespace      = "ebusreg"
	legacyLevelUnknown   = "Unknown"
	legacyAmbiguity      = "unknown"
	legacyUnresolvedBody = "unresolved"
	legacyEnergyMethod   = "get_energy_stats"
	legacyTargetMethod   = "set_target_temp"
)

var errLegacyLayout = errors.New("legacy plane operation layout is unresolved")

type legacyAdmission struct {
	ExtendedRegister bool
	AnswersHeating   bool
	AnswersHotWater  bool
	StandaloneSolar  bool
	HardwareVersion  string
	Manufacturer     string
	Target           byte
}

type legacyBuilder interface {
	Build(map[string]any) ([]byte, error)
}

func fillLegacyPlanes(ctx context.Context, poller *vaillantSemanticPoller, admission legacyAdmission, energyParams map[string]any) (record output.LegacyPlaneRecord, err error) {
	defer func() {
		retainLegacyPlanes(poller, record)
	}()
	record = output.LegacyPlaneRecord{
		ExtendedRegister: admission.ExtendedRegister,
		Catalogues:       legacyCataloguesLoaded(),
		Ambiguity:        legacyAmbiguity,
	}
	info := registry.DeviceInfo{Manufacturer: admission.Manufacturer, HardwareVersion: admission.HardwareVersion}
	if !admission.ExtendedRegister && admission.AnswersHeating {
		plane, operations, readErr := readLegacyProvider(ctx, poller, heating.NewProvider().CreatePlanes(info), admission, energyParams)
		record.Planes = append(record.Planes, plane)
		record.Operations = append(record.Operations, operations...)
		if readErr != nil {
			return record, readErr
		}
	}
	if !admission.ExtendedRegister && admission.AnswersHotWater {
		plane, operations, readErr := readLegacyProvider(ctx, poller, dhw.NewProvider().CreatePlanes(info), admission, energyParams)
		record.Planes = append(record.Planes, plane)
		record.Operations = append(record.Operations, operations...)
		if readErr != nil {
			return record, readErr
		}
	}
	if admission.StandaloneSolar {
		plane, operations, readErr := readLegacyProvider(ctx, poller, solar.NewProvider().CreatePlanes(info), admission, nil)
		record.Planes = append(record.Planes, plane)
		record.Operations = append(record.Operations, operations...)
		if readErr != nil {
			return record, readErr
		}
	}
	return record, nil
}

func readLegacyProvider(ctx context.Context, poller *vaillantSemanticPoller, planes []registry.Plane, admission legacyAdmission, energyParams map[string]any) (output.LegacyPlane, []output.LegacyOperation, error) {
	plane := output.LegacyPlane{}
	var operations []output.LegacyOperation
	if len(planes) == 0 {
		return plane, nil, nil
	}
	plane.Name = planes[0].Name()
	for _, method := range planes[0].Methods() {
		if !method.ReadOnly() {
			operations = append(operations, output.LegacyOperation{
				Namespace:  legacyNamespace,
				Plane:      plane.Name,
				Name:       method.Name(),
				Unresolved: legacyUnresolvedBody,
			})
			continue
		}
		read, err := readLegacyMethod(ctx, poller, method, admission, energyParams)
		plane.Reads = append(plane.Reads, read)
		if err != nil {
			return plane, operations, err
		}
	}
	return plane, operations, nil
}

func readLegacyMethod(ctx context.Context, poller *vaillantSemanticPoller, method registry.Method, admission legacyAdmission, energyParams map[string]any) (output.LegacyRead, error) {
	read := output.LegacyRead{Method: method.Name()}
	builder, ok := method.Template().(legacyBuilder)
	if !ok {
		read.Unresolved = legacyUnresolvedBody
		return read, nil
	}
	params := map[string]any(nil)
	if method.Name() == legacyEnergyMethod {
		if len(energyParams) == 0 {
			read.Unresolved = legacyUnresolvedBody
			return read, nil
		}
		params = energyParams
	}
	payload, err := builder.Build(params)
	if err != nil {
		read.Unresolved = legacyUnresolvedBody
		return read, err
	}
	if poller == nil {
		return read, errors.New("legacy plane poller unavailable")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	source, _ := poller.sourceAndCompanion()
	response, err := poller.sendSemanticFrame(ctx, protocol.Frame{
		FrameType: protocol.FrameTypeInitiatorTarget,
		Source:    source,
		Target:    admission.Target,
		Primary:   method.Template().Primary(),
		Secondary: method.Template().Secondary(),
		Data:      payload,
	})
	if err != nil {
		return read, err
	}
	read.Sent = true
	if response != nil {
		read.Raw = append([]byte(nil), response.Data...)
	}
	fields, decErr := decodeLegacyFields(method, admission.Target, admission.HardwareVersion, read.Raw)
	read.Fields = fields
	if decErr != nil {
		read.Unresolved = legacyUnresolvedBody
		return read, decErr
	}
	return read, nil
}

func decodeLegacyPlane(planeName, methodName string, target byte, hardware string, raw []byte) (output.LegacyRead, error) {
	method, ok := legacyMethod(planeName, methodName, hardware)
	if !ok {
		return output.LegacyRead{}, errLegacyLayout
	}
	read := output.LegacyRead{Method: methodName, Raw: append([]byte(nil), raw...)}
	fields, err := decodeLegacyFields(method, target, hardware, raw)
	read.Fields = fields
	if err != nil {
		read.Unresolved = legacyUnresolvedBody
		return read, err
	}
	return read, nil
}

func legacyMethod(planeName, methodName, hardware string) (registry.Method, bool) {
	info := registry.DeviceInfo{Manufacturer: "Vaillant", HardwareVersion: hardware}
	var planes []registry.Plane
	switch planeName {
	case "heating":
		planes = heating.NewProvider().CreatePlanes(info)
	case "dhw":
		planes = dhw.NewProvider().CreatePlanes(info)
	case "solar":
		planes = solar.NewProvider().CreatePlanes(info)
	default:
		return nil, false
	}
	if len(planes) == 0 {
		return nil, false
	}
	for _, method := range planes[0].Methods() {
		if method.Name() == methodName {
			return method, true
		}
	}
	return nil, false
}

func decodeLegacyFields(method registry.Method, target byte, hardware string, raw []byte) ([]output.LegacyField, error) {
	selected := method.ResponseSchema().Select(target, hardware)
	if len(selected.Fields) == 0 {
		return nil, nil
	}
	values, err := selected.Decode(raw)
	if err != nil {
		return nil, err
	}
	fields := make([]output.LegacyField, 0, len(selected.Fields))
	offset := 0
	for _, field := range selected.Fields {
		size := 0
		if field.Type != nil {
			size = field.Type.Size()
		}
		item := output.LegacyField{Name: field.Name, Level: legacyLevelUnknown}
		if size > 0 && offset+size <= len(raw) {
			item.Raw = append([]byte(nil), raw[offset:offset+size]...)
		}
		offset += size
		if value, ok := values[field.Name]; ok {
			item.Valid = value.Valid
			if number, ok := legacyNumber(value); ok && value.Valid {
				item.Value = &number
			}
		}
		fields = append(fields, item)
	}
	return fields, nil
}

func legacyNumber(value types.Value) (float64, bool) {
	switch typed := value.Value.(type) {
	case float64:
		return typed, true
	case int8:
		return float64(typed), true
	case uint16:
		return float64(typed), true
	case int:
		return float64(typed), true
	default:
		return 0, false
	}
}

func callLegacyOperation(ctx context.Context, poller *vaillantSemanticPoller, planeName, methodName string, target byte, hardware string, raw []byte) error {
	method, ok := legacyMethod(planeName, methodName, hardware)
	if !ok || method.ReadOnly() || poller == nil || ctx == nil {
		return errLegacyLayout
	}
	if _, built := method.Template().(legacyBuilder); built && len(raw) == 0 && len(method.ResponseSchema().Select(target, hardware).Fields) > 0 {
		return errLegacyLayout
	}
	return errLegacyLayout
}

func legacyCataloguesLoaded() []string {
	var names []string
	if len(catalogvaillant.HolidayCatalog()) > 0 {
		names = append(names, "holiday")
	}
	if len(catalogvaillant.ProgramBlocks()) > 0 {
		names = append(names, "program_blocks")
	}
	if len(catalogvaillant.OpaqueProgramming()) > 0 {
		names = append(names, "opaque_programming")
	}
	if len(catalogvaillant.LegacyTimerCatalog()) > 0 {
		names = append(names, "legacy_timer")
	}
	return names
}

func gradeHeatingCandidate(key catalogvaillant.B504CandidateKey) (string, bool, error) {
	entry, found, err := catalogvaillant.B504LookupCandidate(key)
	if err != nil || !found {
		return "", found, err
	}
	return string(entry.B504CandidateGrade), true, nil
}

func gradeOperatingCandidate(key catalogvaillant.B505CandidateKey) (string, bool) {
	entry, found := catalogvaillant.B505CandidateLookup(key)
	if !found {
		return "", false
	}
	return entry.B505CandidateGrade, true
}

func lookupCircuitRole(key catalogvaillant.B504CircuitRoleKey) (bool, error) {
	_, found, err := catalogvaillant.B504CircuitRoleLookup(key)
	return found, err
}

func lookupOperatingRole(key catalogvaillant.B505OperatingRoleKey) (bool, error) {
	_, found, err := catalogvaillant.B505OperatingRoleLookup(key)
	return found, err
}

func lookupLegacySequence(key catalogvaillant.B501LegacySequenceKey) (bool, error) {
	_, found, err := catalogvaillant.B501LegacySequenceLookup(key)
	return found, err
}

func annotateLegacyFields(fields []output.LegacyField, name, level, unit string, found bool) []output.LegacyField {
	if !found || name == "" {
		return fields
	}
	out := make([]output.LegacyField, len(fields))
	copy(out, fields)
	for i := range out {
		out[i].Raw = append([]byte(nil), fields[i].Raw...)
		if fields[i].Value != nil {
			value := *fields[i].Value
			out[i].Value = &value
		}
		if out[i].Name == name {
			out[i].Level = level
			if unit != "" {
				out[i].Unit = unit
			}
		}
	}
	return out
}

func retainLegacyPlanes(p *vaillantSemanticPoller, record output.LegacyPlaneRecord) {
	if p == nil {
		return
	}
	legacyMu.Lock()
	defer legacyMu.Unlock()
	cloned := record.GraphQL()
	legacyHeld[p] = &cloned
}

func retainedLegacyPlanes(p *vaillantSemanticPoller) *output.LegacyPlaneRecord {
	if p == nil {
		return nil
	}
	legacyMu.Lock()
	defer legacyMu.Unlock()
	record := legacyHeld[p]
	if record == nil {
		return nil
	}
	cloned := record.GraphQL()
	return &cloned
}

var (
	legacyMu   sync.Mutex
	legacyHeld = map[*vaillantSemanticPoller]*output.LegacyPlaneRecord{}
)
