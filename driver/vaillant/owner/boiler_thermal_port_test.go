package owner

import (
	"errors"
	"sync"
	"testing"

	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
	"github.com/Project-Helianthus/helianthus-semreg/semreg/v1"
)

type withdrawalRecordingPublisher struct {
	vaillantSemanticPublisher
	mu        sync.Mutex
	lock      sync.Mutex
	withdrawn int
}

func (p *withdrawalRecordingPublisher) boilerPublishLock() sync.Locker { return &p.lock }

func (p *withdrawalRecordingPublisher) withdrawBoilerFlowTemperature() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.withdrawn++
}

type idleThermalHub struct{}

func (idleThermalHub) Declare(BoilerThermalAssetSpec) error { return nil }
func (idleThermalHub) Begin(BoilerThermalHandle, semreg.AssetID) (BoilerThermalTransaction, error) {
	return nil, errors.New("unused")
}
func (idleThermalHub) ReadAdmitted(semreg.AssetID) bool { return false }
func (idleThermalHub) Snapshot(semreg.AssetID) (semreg.Snapshot, error) {
	return semreg.Snapshot{}, errors.New("unused")
}

func TestVaillantOwnerThermalControlWithdrawsFlowTemperature(t *testing.T) {
	flow := 41.5
	publication := &withdrawalRecordingPublisher{}
	poller := &vaillantSemanticPoller{
		publication: publication,
		boiler:      &vaillantBoilerSnapshot{FlowTemperatureC: &flow},
	}
	poller.attachBoilerThermal(idleThermalHub{}, func() (BoilerThermalHandle, error) { return nil, errors.New("no handle") })
	if poller.boilerThermal == nil || poller.boilerThermal.onControl == nil {
		t.Fatal("attachBoilerThermal did not install the publisher and its control callback")
	}
	first := poller.boilerThermal
	poller.attachBoilerThermal(idleThermalHub{}, nil)
	if poller.boilerThermal != first {
		t.Fatal("a second attach replaced the installed publisher")
	}
	if err := poller.boilerThermal.onControl(); err != nil {
		t.Fatalf("onControl() = %v; want nil", err)
	}
	if poller.boiler.FlowTemperatureC == nil || *poller.boiler.FlowTemperatureC != flow {
		t.Fatalf("flow temperature after control = %v; want retained %.1f", poller.boiler.FlowTemperatureC, flow)
	}
	if publication.withdrawn != 0 {
		t.Fatalf("publication withdrawals = %d; want 0", publication.withdrawn)
	}
}

func TestVaillantOwnerThermalRetainsOnlyActualBoilerFlowReads(t *testing.T) {
	poller := &vaillantSemanticPoller{}
	poller.attachBoilerThermal(idleThermalHub{}, func() (BoilerThermalHandle, error) { return nil, errors.New("no handle") })
	request := protocol.Frame{
		FrameType: protocol.FrameTypeInitiatorTarget, Source: 0x71, Target: 0x08,
		Primary: 0xb5, Secondary: 0x09, Data: []byte{0x0d, 0x18, 0x00},
	}
	reply := &protocol.Frame{Data: []byte{0x00, 0x80}}
	poller.retainBoilerThermalRead(request, reply, nil)
	attempt := poller.boilerThermal.attempt
	if attempt == nil || !attempt.Valid || string(attempt.Request) != string(request.Data) || string(attempt.Reply) != string(reply.Data) {
		t.Fatalf("retained attempt = %+v; want a valid copy of the request and reply", attempt)
	}
	request.Data[0] = 0xff
	if attempt.Request[0] != 0x0d {
		t.Fatal("retained request aliases the caller frame")
	}
	poller.retainBoilerThermalRead(request, nil, errors.New("read failed"))
	if poller.boilerThermal.attempt.Valid {
		t.Fatal("a failed read was retained as valid evidence")
	}
}
