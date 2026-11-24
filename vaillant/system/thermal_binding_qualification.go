package system

import (
	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
)

// ControllerRowRef names one native controller row by group, instance and
// register, under the Vaillant B524 identity of the record.
type ControllerRowRef struct {
	Group    byte
	Instance byte
	Register uint16
}

// ThermalBinding is the driver's answer for one native row bound to a current
// declaration. Declared reports whether the row is declared; Reason names why
// it is not. RuleDigest, EvidenceGrade and Level are the digest of the rule
// document, the grade of the row, and the recorded hardware test level. The
// binding decides nothing beyond the declaration.
type ThermalBinding struct {
	Declared      bool
	Reason        string
	Family        string
	Element       string
	EvidenceGrade string
	Level         string
	RuleDigest    string
}

// Reasons named by a thermal binding that is not declared.
const (
	BindingReasonNotDeclared       = "not_declared"
	BindingReasonInstanceNotListed = "instance_not_listed"
	BindingReasonRowNotListed      = "row_not_listed"
)

// controllerThermalInstance is the instance the record lists for every
// declared controller thermal row (II=0x00).
const controllerThermalInstance byte = 0x00

// BindControllerRow answers whether row is declared for the current
// controller result. A row is declared when its instance is the record's
// thermal instance and the record lists its group and register for the
// controller's family.
func BindControllerRow(controller registry.QualifiedController, row ControllerRowRef) ThermalBinding {
	if !controller.Declaration.Declared {
		return ThermalBinding{Reason: BindingReasonNotDeclared}
	}
	base := ThermalBinding{
		Family:     controller.Declaration.Family,
		Level:      controller.Declaration.Level,
		RuleDigest: controller.Declaration.RuleDigest,
	}
	if row.Instance != controllerThermalInstance {
		base.Reason = BindingReasonInstanceNotListed
		return base
	}
	for _, declared := range controller.Declaration.Rows {
		if declared.Group == row.Group && declared.Register == row.Register {
			base.Declared = true
			base.Element = declared.Element
			base.EvidenceGrade = declared.EvidenceGrade
			base.Level = declared.Level
			return base
		}
	}
	base.Reason = BindingReasonRowNotListed
	return base
}

// BindBoilerFlow answers whether the boiler flow row (opcode and selector) is
// declared for the current boiler-target result. The flow row needs no
// qualified binding.
func BindBoilerFlow(target registry.QualifiedTarget, opcode byte, selector uint16) ThermalBinding {
	if !target.Declaration.Declared {
		return ThermalBinding{Reason: BindingReasonNotDeclared}
	}
	base := ThermalBinding{
		Family:     target.DeviceID,
		Level:      target.Declaration.Level,
		RuleDigest: target.Declaration.RuleDigest,
	}
	if opcode != target.Declaration.FlowOpcode || selector != target.Declaration.FlowSelector {
		base.Reason = BindingReasonRowNotListed
		return base
	}
	base.Declared = true
	base.Element = target.Declaration.FlowElement
	base.EvidenceGrade = target.Declaration.FlowEvidenceGrade
	return base
}
