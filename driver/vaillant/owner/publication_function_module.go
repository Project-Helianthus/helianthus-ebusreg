package owner

import (
	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/output"
)

func (o *vaillantSemanticPublication) publishFunctionModuleView(view vaillantFunctionModuleView, source semanticSnapshotSource) {
	verdict := output.PublishedFunctionModuleInterpretation(nativeFunctionModuleVerdict(view.verdict))
	switch source {
	case semanticSnapshotSourceCache:
		o.state.SetFunctionModuleInterpretationFromCache(verdict)
	default:
		o.state.SetFunctionModuleInterpretation(verdict)
	}

	var solar *output.SolarStatus
	var cylinders []output.CylinderStatus
	if isFunctionModuleStructuralWithdrawal(view.verdict) || (view.solar == nil && len(view.cylinders) == 0) {
		solar = &output.SolarStatus{}
		cylinders = []output.CylinderStatus{}
	} else {
		solar = output.SolarStatusFromSnapshot((*output.NativeSolarSnapshot)(view.solar))
		native := make(map[byte]*output.NativeCylinderSnapshot, len(view.cylinders))
		for instance, snapshot := range view.cylinders {
			native[instance] = (*output.NativeCylinderSnapshot)(snapshot)
		}
		cylinders = output.CylinderStatusesFromSnapshots(native)
	}

	switch source {
	case semanticSnapshotSourceCache:
		o.state.SetSolarFromCache(solar)
		o.state.SetCylindersFromCache(cylinders)
	default:
		o.state.SetSolar(solar)
		o.state.SetCylinders(cylinders)
	}
}

func (o *vaillantSemanticPublication) functionModulePublishedAsGPIOOnly() bool {
	return o.state.FunctionModuleSemanticMode() == output.FunctionModuleSemanticModeGPIOOnly
}
