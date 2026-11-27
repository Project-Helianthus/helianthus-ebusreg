package owner

import (
	"context"
	"encoding/binary"
	"encoding/json"

	catalogvaillant "github.com/Project-Helianthus/helianthus-ebusreg/catalog/vaillant"
	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/output"
	"github.com/Project-Helianthus/helianthus-ebusreg/vaillant/system"
)

const catalogueUnresolved = "unresolved"

type enumerationSpec struct {
	group       byte
	register    uint16
	alternative string
	scope       string
}

var familyEnumerationSpecs = []enumerationSpec{
	{group: 1, register: 3, alternative: "u16", scope: "00"},
	{group: 1, register: 13, alternative: "u8", scope: "00"},
	{group: 2, register: 2, alternative: "u16", scope: "dynamic"},
	{group: 2, register: 29, alternative: "u8", scope: "dynamic"},
	{group: 3, register: 6, alternative: "u16", scope: "dynamic"},
	{group: 3, register: 27, alternative: "u8", scope: "dynamic"},
	{group: 3, register: 27, alternative: "u16", scope: "dynamic"},
	{group: 3, register: 0x000E, alternative: "u16", scope: "dynamic"},
}

type scalarSpec struct {
	group    byte
	register uint16
	families string
}

var unreadScalarSpecs = []scalarSpec{
	{0x00, 0x0002, "both"},
	{0x00, 0x0004, "both"},
	{0x00, 0x0006, "controller-700"},
	{0x00, 0x0017, "both"},
	{0x00, 0x0018, "controller-720"},
	{0x00, 0x0029, "both"},
	{0x00, 0x0038, "both"},
	{0x00, 0x0042, "controller-700"},
	{0x00, 0x0043, "controller-700"},
	{0x00, 0x0067, "both"},
	{0x00, 0x0068, "both"},
	{0x00, 0x009F, "both"},
	{0x00, 0x00A0, "both"},
	{0x00, 0x00EF, "controller-720"},
	{0x03, 0x0002, "both"},
	{0x03, 0x0009, "both"},
	{0x03, 0x0015, "both"},
	{0x03, 0x0023, "controller-720"},
	{0x02, 0x000A, "both"},
	{0x02, 0x000B, "both"},
	{0x02, 0x0011, "both"},
	{0x02, 0x001A, "controller-720"},
	{0x02, 0x001C, "controller-720"},
	{0x02, 0x001F, "both"},
	{0x03, 0x001B, "both"},
	{0x01, 0x0014, "both"},
	{0x01, 0x0019, "both"},
	{0x01, 0x0015, "both"},
}

func familyCodebookKey(family, alternative, scope string, group byte, register uint16) catalogvaillant.B524FamilyCodebookKey {
	secondary := uint8(36)
	operation := uint8(2)
	grouped := group
	address := register
	return catalogvaillant.B524FamilyCodebookKey{
		B524Family:                      family,
		B524FamilyCodebookVariant:       "none",
		B524FamilyCodebookPrimary:       181,
		B524FamilyCodebookSecondary:     &secondary,
		B524FamilyCodebookOperation:     &operation,
		B524FamilyCodebookDirection:     "read",
		B524FamilyCodebookField:         "value",
		B524FamilyCodebookAlternative:   alternative,
		B524FamilyCodebookGroup:         &grouped,
		B524FamilyCodebookInstanceScope: scope,
		B524FamilyCodebookRegister:      &address,
	}
}

func labelFamilyEnumeration(key catalogvaillant.B524FamilyCodebookKey, raw []byte) output.CatalogueEnumeration {
	item := output.CatalogueEnumeration{
		Alternative: key.B524FamilyCodebookAlternative,
		Raw:         append([]byte(nil), raw...),
		Label:       catalogueUnresolved,
		Level:       "Unknown",
	}
	if key.B524FamilyCodebookGroup != nil {
		item.Group = *key.B524FamilyCodebookGroup
	}
	if key.B524FamilyCodebookRegister != nil {
		item.Register = *key.B524FamilyCodebookRegister
	}
	if code, ok := codebookCode(key.B524FamilyCodebookAlternative, raw); ok {
		item.Code = &code
	}
	row, found, err := catalogvaillant.B524FamilyCodebookLookup(key)
	if err != nil || !found {
		return item
	}
	var decoded struct {
		Grade       string `json:"grade"`
		Enumeration map[string]struct {
			Value int    `json:"value"`
			Label string `json:"label"`
		} `json:"enumeration"`
	}
	if json.Unmarshal(row, &decoded) != nil || decoded.Grade == "" {
		return item
	}
	item.Level = decoded.Grade
	if item.Code == nil {
		return item
	}
	for _, entry := range decoded.Enumeration {
		if entry.Value == *item.Code && entry.Label != "" {
			item.Label = entry.Label
			return item
		}
	}
	return item
}

func codebookCode(alternative string, raw []byte) (int, bool) {
	switch alternative {
	case "u8":
		if len(raw) < 1 {
			return 0, false
		}
		return int(raw[0]), true
	case "signed8":
		if len(raw) < 1 {
			return 0, false
		}
		return int(int8(raw[0])), true
	case "u16":
		if len(raw) < 2 {
			return 0, false
		}
		return int(binary.LittleEndian.Uint16(raw[:2])), true
	default:
		return 0, false
	}
}

func fillFamilyCodebooks(p *vaillantSemanticPoller, ctx context.Context, family string) output.FamilyCatalogue {
	var view output.FamilyCatalogue
	for _, spec := range familyEnumerationSpecs {
		raw, ok := readCatalogueRegister(p, ctx, spec.group, spec.register)
		if !ok {
			continue
		}
		key := familyCodebookKey(family, spec.alternative, spec.scope, spec.group, spec.register)
		view.Enumerations = append(view.Enumerations, labelFamilyEnumeration(key, raw))
	}
	for _, spec := range unreadScalarSpecs {
		if spec.families != "both" && spec.families != family {
			continue
		}
		raw, ok := readCatalogueRegister(p, ctx, spec.group, spec.register)
		if !ok {
			continue
		}
		view.Scalars = append(view.Scalars, interpretUnreadScalar(family, spec.group, spec.register, raw))
	}
	view.Ventilation = readVentilationView(p, ctx, family)
	publishFamilyCatalogue(p, view)
	return view
}

func publishFamilyCatalogue(p *vaillantSemanticPoller, view output.FamilyCatalogue) {
	rememberFamilyCatalogue(p, view)
	owner := publicationOwner(p)
	if owner == nil || owner.state == nil {
		return
	}
	if system := owner.state.System(); system != nil {
		stored := view
		system.SetVentilation(view.Ventilation)
		system.SetCatalogue(&stored)
		owner.state.SetSystem(system)
	}
	if zones := owner.state.Zones(); zones != nil {
		grouped := enumerationsForGroup(view.Enumerations, localZones.group)
		for i := range zones {
			zones[i].State.SetEnumerations(enumerationsForGroup(grouped, localZones.group))
		}
		owner.state.SetZones(zones)
	}
	if status := owner.state.DHW(); status != nil {
		status.State.SetEnumerations(enumerationsForGroup(view.Enumerations, localDHW.group))
		owner.state.SetDHW(status)
	}
	if circuits := owner.state.Circuits(); circuits != nil {
		grouped := enumerationsForGroup(view.Enumerations, localCircuits.group)
		for i := range circuits {
			circuits[i].Config.SetEnumerations(enumerationsForGroup(grouped, localCircuits.group))
		}
		owner.state.SetCircuits(circuits)
	}
}

func enumerationsForGroup(items []output.CatalogueEnumeration, group byte) []output.CatalogueEnumeration {
	var out []output.CatalogueEnumeration
	for _, item := range items {
		if item.Group != group {
			continue
		}
		item.Raw = append([]byte(nil), item.Raw...)
		if item.Code != nil {
			code := *item.Code
			item.Code = &code
		}
		out = append(out, item)
	}
	return out
}

func readCatalogueRegister(p *vaillantSemanticPoller, ctx context.Context, group byte, register uint16) ([]byte, bool) {
	if p == nil {
		return nil, false
	}
	raw, ok := p.readB524Value(ctx, vaillantB524OpcodeLocal, group, 0, register)
	if !ok {
		return nil, false
	}
	return append([]byte(nil), raw...), true
}

func catalogueTier(group byte, register uint16) string {
	_, profile, _ := semanticReadB524LocalDescriptorPolicy(group, register)
	return string(profile)
}

func interpretUnreadScalar(family string, group byte, register uint16, raw []byte) output.CatalogueScalar {
	item := output.CatalogueScalar{
		Name:     ventilationScalarName(family, group, register),
		Group:    group,
		Register: register,
		Raw:      append([]byte(nil), raw...),
		Level:    "Unknown",
		Tier:     catalogueTier(group, register),
	}
	switch group {
	case 0x00:
		result := system.B524SystemInterpretScalars(systemScalarInput(family, group, register, raw))
		if result.B524SystemScalarState == system.B524SystemScalarCandidate && result.B524SystemScalarValue != nil {
			item.Level = result.B524SystemScalarGrade
			item.Unit = withheldUnit(result.B524SystemScalarDeclaredUnit)
			value := scalarNumber(result.B524SystemScalarCodec, result.B524SystemScalarValue.B524SystemScalarFloat32, result.B524SystemScalarValue.B524SystemScalarUint16)
			item.Value = &value
		}
	case 0x03:
		result := system.B524ZoneScalarInterpret(zoneScalarInput(family, register, raw))
		if result.B524ZoneScalarState == system.B524ZoneScalarCandidate && result.B524ZoneScalarValue.B524ZoneScalarPresent {
			item.Level = string(result.B524ZoneScalarGrade)
			item.Unit = withheldUnit(result.B524ZoneScalarUnit)
			value := float64(result.B524ZoneScalarValue.B524ZoneScalarValue)
			item.Value = &value
		}
	case 0x02:
		result := system.B524InterpretCircuitScalar(circuitScalarInput(family, register, raw))
		if result.B524CircuitScalarState == system.B524CircuitScalarCandidate && result.B524CircuitScalarValue.B524CircuitScalarPresent {
			item.Level = string(result.B524CircuitScalarGrade)
			item.Unit = withheldUnit(result.B524CircuitScalarDeclaredUnit)
			value := float64(result.B524CircuitScalarValue.B524CircuitScalarValue)
			item.Value = &value
		}
	}
	return item
}

func withheldUnit(unit string) string {
	if unit == "" || unit == "Unknown" {
		return ""
	}
	return unit
}

func scalarNumber(codec string, floatValue float32, word uint16) float64 {
	if codec == "uint16_le" {
		return float64(word)
	}
	return float64(floatValue)
}

func readVentilationView(p *vaillantSemanticPoller, ctx context.Context, family string) *output.VentilationView {
	var group byte
	var dayRegister, nightRegister uint16
	switch family {
	case "controller-720":
		group, dayRegister, nightRegister = 0x09, 0x000D, 0x000E
	case "controller-700":
		group, dayRegister, nightRegister = 0x00, 0x0042, 0x0043
	default:
		return nil
	}
	day, dayOK := readCatalogueRegister(p, ctx, group, dayRegister)
	night, nightOK := readCatalogueRegister(p, ctx, group, nightRegister)
	if !dayOK && !nightOK {
		return nil
	}
	view := &output.VentilationView{Level: "Unknown"}
	if dayOK {
		view.DayFanSpeedLimitRaw = day
		if value, grade, ok := decodeVentilationWord(family, group, dayRegister, day); ok {
			view.DayFanSpeedLimit = &value
			view.Level = grade
		}
	}
	if nightOK {
		view.NightFanSpeedLimitRaw = night
		if value, grade, ok := decodeVentilationWord(family, group, nightRegister, night); ok {
			view.NightFanSpeedLimit = &value
			view.Level = grade
		}
	}
	return view
}

func decodeVentilationWord(family string, group byte, register uint16, raw []byte) (uint16, string, bool) {
	if family == "controller-720" && group == 0x09 {
		result := system.B524DHWVentilationInterpretScalars(ventilationScalarInput(family, group, register, raw))
		if result.B524DHWVentilationScalarState != system.B524DHWVentilationScalarCandidate || result.B524DHWVentilationScalarValue == nil {
			return 0, "", false
		}
		return result.B524DHWVentilationScalarValue.B524DHWVentilationScalarUint16, result.B524DHWVentilationScalarGrade, true
	}
	result := system.B524SystemInterpretScalars(systemScalarInput(family, group, register, raw))
	if result.B524SystemScalarState != system.B524SystemScalarCandidate || result.B524SystemScalarValue == nil || result.B524SystemScalarCodec != "uint16_le" {
		return 0, "", false
	}
	return result.B524SystemScalarValue.B524SystemScalarUint16, result.B524SystemScalarGrade, true
}

func ventilationScalarName(family string, group byte, register uint16) string {
	switch {
	case family == "controller-700" && group == 0x00 && register == 0x0042:
		return "ventilation_fan_speed_limit_maximum_at_day"
	case family == "controller-700" && group == 0x00 && register == 0x0043:
		return "ventilation_fan_speed_limit_maximum_at_night"
	case family == "controller-720" && group == 0x09 && register == 0x000D:
		return "ventilation_fan_speed_limit_maximum_at_day"
	case family == "controller-720" && group == 0x09 && register == 0x000E:
		return "ventilation_fan_speed_limit_maximum_at_night"
	default:
		return ""
	}
}

func systemScalarInput(family string, group byte, register uint16, raw []byte) system.B524SystemScalarInput {
	return system.B524SystemScalarInput{
		B524SystemScalarContext: system.B524SystemScalarContext{
			B524SystemScalarFamily:      systemField(family),
			B524SystemScalarSoftware:    systemField("none"),
			B524SystemScalarHardware:    systemField("none"),
			B524SystemScalarAPIVersion:  systemField("none"),
			B524SystemScalarAPIRevision: systemField("none"),
			B524SystemScalarOpcode:      system.B524SystemScalarField[byte]{B524SystemScalarPresent: true, B524SystemScalarValue: vaillantB524OpcodeLocal},
			B524SystemScalarGroup:       system.B524SystemScalarField[byte]{B524SystemScalarPresent: true, B524SystemScalarValue: group},
			B524SystemScalarInstance:    system.B524SystemScalarField[byte]{B524SystemScalarPresent: true, B524SystemScalarValue: 0},
			B524SystemScalarRegister:    system.B524SystemScalarField[uint16]{B524SystemScalarPresent: true, B524SystemScalarValue: register},
		},
		B524SystemScalarReplyKind:        systemField("echo_payload"),
		B524SystemScalarReturnedGroup:    system.B524SystemScalarField[byte]{B524SystemScalarPresent: true, B524SystemScalarValue: group},
		B524SystemScalarReturnedRegister: system.B524SystemScalarField[uint16]{B524SystemScalarPresent: true, B524SystemScalarValue: register},
		B524SystemScalarCodecGrade:       systemField("Hypothesis"),
		B524SystemScalarRaw:              append([]byte(nil), raw...),
	}
}

func systemField(value string) system.B524SystemScalarField[string] {
	return system.B524SystemScalarField[string]{B524SystemScalarPresent: true, B524SystemScalarValue: value}
}

func zoneScalarInput(family string, register uint16, raw []byte) system.B524ZoneScalarInput {
	return system.B524ZoneScalarInput{
		B524ZoneScalarContext: system.B524ZoneScalarContext{
			B524ZoneScalarFamily:      system.B524ZoneScalarField[string]{B524ZoneScalarPresent: true, B524ZoneScalarValue: family},
			B524ZoneScalarSoftware:    system.B524ZoneScalarField[string]{B524ZoneScalarPresent: true, B524ZoneScalarValue: "none"},
			B524ZoneScalarHardware:    system.B524ZoneScalarField[string]{B524ZoneScalarPresent: true, B524ZoneScalarValue: "none"},
			B524ZoneScalarAPIVersion:  system.B524ZoneScalarField[string]{B524ZoneScalarPresent: true, B524ZoneScalarValue: "none"},
			B524ZoneScalarAPIRevision: system.B524ZoneScalarField[string]{B524ZoneScalarPresent: true, B524ZoneScalarValue: "none"},
			B524ZoneScalarOpcode:      system.B524ZoneScalarField[byte]{B524ZoneScalarPresent: true, B524ZoneScalarValue: vaillantB524OpcodeLocal},
			B524ZoneScalarGroup:       system.B524ZoneScalarField[byte]{B524ZoneScalarPresent: true, B524ZoneScalarValue: 0x03},
			B524ZoneScalarInstance:    system.B524ZoneScalarField[byte]{B524ZoneScalarPresent: true, B524ZoneScalarValue: 0},
			B524ZoneScalarRegister:    system.B524ZoneScalarField[uint16]{B524ZoneScalarPresent: true, B524ZoneScalarValue: register},
		},
		B524ZoneScalarReplyKind:        "echo_payload",
		B524ZoneScalarReturnedGroup:    system.B524ZoneScalarField[byte]{B524ZoneScalarPresent: true, B524ZoneScalarValue: 0x03},
		B524ZoneScalarReturnedRegister: system.B524ZoneScalarField[uint16]{B524ZoneScalarPresent: true, B524ZoneScalarValue: register},
		B524ZoneScalarCodecGrade:       system.B524ZoneScalarHypothesis,
		B524ZoneScalarRaw:              append([]byte(nil), raw...),
	}
}

func circuitScalarInput(family string, register uint16, raw []byte) system.B524CircuitScalarInput {
	return system.B524CircuitScalarInput{
		B524CircuitScalarContext: system.B524CircuitScalarContext{
			B524CircuitScalarFamily:      system.B524CircuitScalarField[string]{B524CircuitScalarPresent: true, B524CircuitScalarValue: family},
			B524CircuitScalarSoftware:    system.B524CircuitScalarField[string]{B524CircuitScalarPresent: true, B524CircuitScalarValue: "none"},
			B524CircuitScalarHardware:    system.B524CircuitScalarField[string]{B524CircuitScalarPresent: true, B524CircuitScalarValue: "none"},
			B524CircuitScalarAPIVersion:  system.B524CircuitScalarField[string]{B524CircuitScalarPresent: true, B524CircuitScalarValue: "none"},
			B524CircuitScalarAPIRevision: system.B524CircuitScalarField[string]{B524CircuitScalarPresent: true, B524CircuitScalarValue: "none"},
			B524CircuitScalarOpcode:      system.B524CircuitScalarField[byte]{B524CircuitScalarPresent: true, B524CircuitScalarValue: vaillantB524OpcodeLocal},
			B524CircuitScalarGroup:       system.B524CircuitScalarField[byte]{B524CircuitScalarPresent: true, B524CircuitScalarValue: 0x02},
			B524CircuitScalarInstance:    system.B524CircuitScalarField[byte]{B524CircuitScalarPresent: true, B524CircuitScalarValue: 0},
			B524CircuitScalarRegister:    system.B524CircuitScalarField[uint16]{B524CircuitScalarPresent: true, B524CircuitScalarValue: register},
		},
		B524CircuitScalarReplyKind:        system.B524CircuitScalarField[string]{B524CircuitScalarPresent: true, B524CircuitScalarValue: "echo_payload"},
		B524CircuitScalarReturnedGroup:    system.B524CircuitScalarField[byte]{B524CircuitScalarPresent: true, B524CircuitScalarValue: 0x02},
		B524CircuitScalarReturnedRegister: system.B524CircuitScalarField[uint16]{B524CircuitScalarPresent: true, B524CircuitScalarValue: register},
		B524CircuitScalarCodecGrade:       system.B524CircuitScalarField[system.B524CircuitScalarGrade]{B524CircuitScalarPresent: true, B524CircuitScalarValue: system.B524CircuitScalarHypothesis},
		B524CircuitScalarRaw:              append([]byte(nil), raw...),
	}
}

func ventilationScalarInput(family string, group byte, register uint16, raw []byte) system.B524DHWVentilationScalarInput {
	return system.B524DHWVentilationScalarInput{
		B524DHWVentilationScalarContext: system.B524DHWVentilationScalarContext{
			B524DHWVentilationScalarFamily:      ventilationField(family),
			B524DHWVentilationScalarSoftware:    ventilationField("none"),
			B524DHWVentilationScalarHardware:    ventilationField("none"),
			B524DHWVentilationScalarAPIVersion:  ventilationField("none"),
			B524DHWVentilationScalarAPIRevision: ventilationField("none"),
			B524DHWVentilationScalarOpcode:      system.B524DHWVentilationScalarField[byte]{B524DHWVentilationScalarPresent: true, B524DHWVentilationScalarValue: vaillantB524OpcodeLocal},
			B524DHWVentilationScalarGroup:       system.B524DHWVentilationScalarField[byte]{B524DHWVentilationScalarPresent: true, B524DHWVentilationScalarValue: group},
			B524DHWVentilationScalarInstance:    system.B524DHWVentilationScalarField[byte]{B524DHWVentilationScalarPresent: true, B524DHWVentilationScalarValue: 0},
			B524DHWVentilationScalarRegister:    system.B524DHWVentilationScalarField[uint16]{B524DHWVentilationScalarPresent: true, B524DHWVentilationScalarValue: register},
		},
		B524DHWVentilationScalarReplyKind:        ventilationField("echo_payload"),
		B524DHWVentilationScalarReturnedGroup:    system.B524DHWVentilationScalarField[byte]{B524DHWVentilationScalarPresent: true, B524DHWVentilationScalarValue: group},
		B524DHWVentilationScalarReturnedRegister: system.B524DHWVentilationScalarField[uint16]{B524DHWVentilationScalarPresent: true, B524DHWVentilationScalarValue: register},
		B524DHWVentilationScalarCodecGrade:       ventilationField("Hypothesis"),
		B524DHWVentilationScalarRaw:              append([]byte(nil), raw...),
	}
}

func ventilationField(value string) system.B524DHWVentilationScalarField[string] {
	return system.B524DHWVentilationScalarField[string]{B524DHWVentilationScalarPresent: true, B524DHWVentilationScalarValue: value}
}
