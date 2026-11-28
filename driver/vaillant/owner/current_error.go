package owner

import (
	"context"

	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
	b503 "github.com/Project-Helianthus/helianthus-ebusgo/protocol/vaillant/b503"

	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/heat"
)

func (p *vaillantSemanticPoller) boilerFlowReading() boilerFlowMark {
	if p == nil {
		return boilerFlowMark{}
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.boiler == nil || p.boiler.FlowTemperatureC == nil {
		return boilerFlowMark{}
	}
	return boilerFlowMark{
		celsius:    cloneFloat64Ptr(p.boiler.FlowTemperatureC),
		cached:     p.boilerFlowCached,
		observedAt: p.boilerFlowObservedAt,
	}
}

func (p *vaillantSemanticPoller) refreshCurrentError(ctx context.Context) {
	if p == nil {
		return
	}
	p.mu.Lock()
	controller := p.controller
	boiler := p.boilerAddress
	p.mu.Unlock()
	controllerRead := p.readCurrentError(ctx, controller)
	boilerRead := p.readCurrentError(ctx, boiler)
	symbol, known := heat.FaultFromCurrentError(controllerRead, boilerRead)
	if !known {
		return
	}
	p.mu.Lock()
	p.currentErrorFault = symbol
	p.mu.Unlock()
}

func (p *vaillantSemanticPoller) readCurrentError(ctx context.Context, target byte) heat.CurrentErrorRead {
	if p == nil || target == 0 || (p.bus == nil && p.sendFrameFn == nil) {
		return heat.CurrentErrorRead{}
	}
	if ctx == nil {
		ctx = context.Background()
	}
	p.mu.Lock()
	source := p.source
	p.mu.Unlock()
	p.readMu.Lock()
	response, err := p.sendSemanticFrame(ctx, protocol.Frame{
		FrameType: protocol.FrameTypeInitiatorTarget,
		Source:    source,
		Target:    target,
		Primary:   b503PrimaryByte,
		Secondary: b503SecondaryByte,
		Data:      append([]byte(nil), b503.EncodeCurrentError()...),
	})
	p.readMu.Unlock()
	if err != nil || response == nil {
		return heat.CurrentErrorRead{Present: true}
	}
	slots, decErr := b503.DecodeCurrentError(response.Data)
	if decErr != nil {
		return heat.CurrentErrorRead{Present: true}
	}
	return heat.CurrentErrorRead{Present: true, Admitted: true, Slots: slots}
}
