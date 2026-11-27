package driver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
	semdriver "github.com/Project-Helianthus/helianthus-semreg/semreg/v1/driver"

	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/owner"
	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/owner/b503session"
	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
)

func TestEbusFactoryTwoInstances(t *testing.T) {
	factory := NewFactory()
	var opened atomic.Int32
	factory.open = func(context.Context, owner.Config) (owner.Substrate, error) {
		opened.Add(1)
		raw := newSimulatedTransport()
		return owner.Substrate{
			Registry:  registry.NewDeviceRegistry(nil),
			Bus:       protocol.NewBus(raw, protocol.BusConfig{}, 1),
			Transport: raw,
		}, nil
	}
	disabled, err := factory.Plan(semdriver.Config{Instance: "off", Enabled: false, SecretRef: "hidden-secret"})
	if err != nil || disabled != nil {
		t.Fatalf("disabled plan = %v, %v", disabled, err)
	}
	if _, err := factory.Plan(semdriver.Config{Instance: "a", Enabled: true, Options: json.RawMessage(`{"not_a_field":1}`)}); err == nil {
		t.Fatal("unknown option accepted")
	}
	if opened.Load() != 0 {
		t.Fatal("plan opened a transport")
	}
	left := planInstance(t, factory, "bus-a", 17, time.Second)
	right := planInstance(t, factory, "bus-b", 34, 2*time.Second)
	if strings.Contains(fmt.Sprintf("%+v", []semdriver.InstancePlan{left, right}), "hidden-secret") {
		t.Fatal("plan retained the secret reference")
	}
	if left.Instance == right.Instance || left.Transports[0].ConfigKey == right.Transports[0].ConfigKey {
		t.Fatal("instances share an identity")
	}
	if left.Transports[0].ConfigKey != "bus-a" || right.Transports[0].ConfigKey != "bus-b" {
		t.Fatal("config key is not the instance key")
	}
	if left.Transports[0].Provider != transportProvider || left.Transports[0].AccessDomain != transportAccessDomain || left.Transports[0].Share != semdriver.ShareSerialized {
		t.Fatalf("claim = %+v", left.Transports[0])
	}
	instA := constructHeld(t, factory, left)
	instB := constructHeld(t, factory, right)
	if opened.Load() != 0 {
		t.Fatal("construct opened a transport")
	}
	if err := instA.Start(context.Background(), nil); err != nil {
		t.Fatal(err)
	}
	if err := instB.Start(context.Background(), nil); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = instA.Stop(context.Background()) })
	t.Cleanup(func() { _ = instB.Stop(context.Background()) })
	if opened.Load() != 2 {
		t.Fatalf("opens = %d", opened.Load())
	}
	assertIndependentInstances(t, instA, instB)
}

func TestEbusFactoryB503Lifecycle(t *testing.T) {
	t.Run("cold-start", testEbusFactoryColdStartFence)
	t.Run("withdrawal", testEbusFactoryCurrentWithdrawal)
	t.Run("stale-source", testEbusFactoryStaleSource)
}

func testEbusFactoryColdStartFence(t *testing.T) {
	inst, _ := startHeld(t, "cold", 0x10)
	hist := inst.acquired.History()
	if inst.acquired.Runtime() == nil {
		t.Fatal("owner has no lifecycle runtime")
	}
	if !hist.TargetBlocked(owner.ConfiguredMonitorTarget) {
		t.Fatal("configured target is not fenced")
	}
	const later = byte(0x21)
	hist.MarkQualifiedTarget(later)
	if !hist.TargetBlocked(later) {
		t.Fatal("later qualified target is not fenced")
	}
	for _, target := range []byte{owner.ConfiguredMonitorTarget, later} {
		calls := 0
		_, err := hist.EnableOperation(context.Background(), target, func(context.Context, byte) b503session.DispatchOutcome {
			calls++
			return b503session.DispatchOutcome{}
		})
		if !errors.Is(err, b503session.ErrCleanupPending) || calls != 0 {
			t.Fatalf("target %d enable = %v calls=%d", target, err, calls)
		}
	}
}

func testEbusFactoryCurrentWithdrawal(t *testing.T) {
	inst, _ := startHeld(t, "withdraw", 0x10)
	hist := inst.acquired.History()
	b503session.ReleaseRestartFenceForTest(hist)
	runtime := inst.acquired.Runtime()
	runtime.EBusDriverActivated(1)
	if hist.TransportKey().TransportEpoch != 1 || inst.acquired.RefreshCalls() != 0 {
		t.Fatalf("epoch %d refresh %d", hist.TransportKey().TransportEpoch, inst.acquired.RefreshCalls())
	}
	_, err := hist.EnableOperation(context.Background(), 0x15, func(context.Context, byte) b503session.DispatchOutcome {
		return b503session.DispatchOutcome{Emitted: true, Native: b503session.NativeACK, Transport: hist.TransportKey()}
	})
	if err != nil {
		t.Fatal(err)
	}
	before, ok := hist.CleanupObligation()
	if !ok || before.GatewayCleanupAttemptID == "" || before.TransportEpoch != 1 || before.Target != 0x15 {
		t.Fatalf("cleanup %+v present %v", before, ok)
	}
	runtime.EBusDriverActivated(2)
	if hist.TransportKey().TransportEpoch != 1 || hist.LifecycleTransportKey().TransportEpoch != 2 || inst.acquired.RefreshCalls() != 0 {
		t.Fatalf("epoch %d lifecycle %d refresh %d", hist.TransportKey().TransportEpoch, hist.LifecycleTransportKey().TransportEpoch, inst.acquired.RefreshCalls())
	}
	runtime.EBusDriverWithdrawn(1)
	if !hist.IsOwned() {
		t.Fatal("stale withdrawal released the owner")
	}
	runtime.EBusDriverWithdrawn(2)
	if got := hist.StatusSnapshot(); got.State != b503session.Idle || got.Owned {
		t.Fatalf("snapshot %+v", got)
	}
	after, ok := hist.CleanupObligation()
	if !ok || after.GatewayCleanupAttemptID != before.GatewayCleanupAttemptID || after.Target != before.Target || after.TransportEpoch != before.TransportEpoch {
		t.Fatalf("cleanup drift before %+v after %+v", before, after)
	}
}

func testEbusFactoryStaleSource(t *testing.T) {
	inst, raw := startHeld(t, "stale", 0x10)
	hist := inst.acquired.History()
	b503session.ReleaseRestartFenceForTest(hist)
	runtime := inst.acquired.Runtime()
	runtime.EBusDriverActivated(7)
	if _, err := hist.Enable(context.Background()); err != nil {
		t.Fatal(err)
	}
	inst.acquired.LockDispatcher()
	unlocked := false
	unlock := func() {
		if !unlocked {
			unlocked = true
			inst.acquired.UnlockDispatcher()
		}
	}
	defer unlock()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	errCh := make(chan error, 1)
	go func() {
		_, err := inst.acquired.Invoke(ctx, 0x08, []byte{0x00, 0x01})
		errCh <- err
	}()
	waitForDispatcher(t, inst)
	runtime.EBusDriverWithdrawn(7)
	runtime.EBusDriverActivated(8)
	_, reenableErr := hist.EnableOperation(context.Background(), 0, func(context.Context, byte) b503session.DispatchOutcome {
		t.Error("cleanup-fenced re-enable reached the dispatcher")
		return b503session.DispatchOutcome{}
	})
	inst.acquired.SetSource(0)
	unlock()
	var invokeErr error
	select {
	case invokeErr = <-errCh:
	case <-time.After(5 * time.Second):
		t.Fatal("invoke did not return")
	}
	if !owner.SourceUnavailable(invokeErr) {
		t.Fatalf("invoke = %v", invokeErr)
	}
	if !errors.Is(reenableErr, b503session.ErrCleanupPending) || hist.IsOwned() {
		t.Fatalf("re-enable %v owned %v", reenableErr, hist.IsOwned())
	}
	if hist.TransportKey().TransportEpoch != 8 {
		t.Fatalf("epoch = %d", hist.TransportKey().TransportEpoch)
	}
	obligation, ok := hist.CleanupObligation()
	if !ok || obligation.TransportEpoch != 7 || obligation.LastNative != b503session.NativeAmbiguous || obligation.LastAttemptEpoch != 0 {
		t.Fatalf("cleanup %+v present %v", obligation, ok)
	}
	if raw.Writes() != 0 {
		t.Fatalf("writes = %d", raw.Writes())
	}
}

func planInstance(t *testing.T, factory *Factory, key string, scan byte, state time.Duration) semdriver.InstancePlan {
	t.Helper()
	plans, err := factory.Plan(semdriver.Config{
		Instance:  key,
		Enabled:   true,
		SecretRef: "hidden-secret",
		Transport: json.RawMessage(fmt.Sprintf(`{"scan_source":%d}`, scan)),
		Options:   json.RawMessage(fmt.Sprintf(`{"state_interval":%q}`, state.String())),
	})
	if err != nil || len(plans) != 1 {
		t.Fatalf("plan %s: %v len %d", key, err, len(plans))
	}
	return plans[0]
}

func constructHeld(t *testing.T, factory *Factory, plan semdriver.InstancePlan) *instance {
	t.Helper()
	built, err := factory.Construct(context.Background(), plan)
	if err != nil {
		t.Fatal(err)
	}
	inst := built.(*instance)
	inst.holdPolling = true
	return inst
}

func startHeld(t *testing.T, key string, scan byte) (*instance, *simulatedTransport) {
	t.Helper()
	factory := NewFactory()
	raw := newSimulatedTransport()
	factory.open = func(context.Context, owner.Config) (owner.Substrate, error) {
		return owner.Substrate{
			Registry:  registry.NewDeviceRegistry(nil),
			Bus:       protocol.NewBus(raw, protocol.BusConfig{}, 1),
			Transport: raw,
		}, nil
	}
	inst := constructHeld(t, factory, planInstance(t, factory, key, scan, time.Second))
	if err := inst.Start(context.Background(), nil); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = inst.Stop(context.Background()) })
	t.Cleanup(func() { inst.acquired.History().ResetForRestart() })
	if inst.acquired.Transport() != raw || inst.acquired.Source() != scan {
		t.Fatal("started owner did not receive the factory transport and source")
	}
	return inst, raw
}

func assertIndependentInstances(t *testing.T, instA, instB *instance) {
	t.Helper()
	if instA.acquired.History() == instB.acquired.History() || instA.acquired.Runtime() == instB.acquired.Runtime() {
		t.Fatal("instances share a history session")
	}
	if instA.acquired.Schedule()["refresh_state"] != time.Second || instB.acquired.Schedule()["refresh_state"] != 2*time.Second {
		t.Fatalf("schedules %v %v", instA.acquired.Schedule(), instB.acquired.Schedule())
	}
	copied := instA.acquired.Schedule()
	copied["refresh_state"] = 9 * time.Second
	if instA.acquired.Schedule()["refresh_state"] != time.Second || instB.acquired.Schedule()["refresh_state"] != 2*time.Second {
		t.Fatal("schedule view is shared")
	}
	if instA.acquired.Views().Source != 17 || instB.acquired.Views().Source != 34 {
		t.Fatalf("views %v %v", instA.acquired.Views(), instB.acquired.Views())
	}
	instA.acquired.SetSource(19)
	if instA.acquired.Source() != 19 || instB.acquired.Source() != 34 || instB.acquired.Views().Source != 34 {
		t.Fatal("source is shared")
	}
	if instA.acquired.Transport() == instB.acquired.Transport() {
		t.Fatal("transport is shared")
	}
	instA.FenceWithdrawal()
	beforeSchedule := instB.acquired.Schedule()
	beforeView := instB.acquired.Views()
	report := instA.Stop(context.Background())
	if !report.Confirmed || instA.running() {
		t.Fatalf("stop report %+v running %v", report, instA.running())
	}
	if !transportClosed(instA) || transportClosed(instB) || !instB.running() {
		t.Fatal("stopping one instance changed the other")
	}
	if !reflect.DeepEqual(beforeSchedule, instB.acquired.Schedule()) || instB.acquired.Views() != beforeView {
		t.Fatal("survivor schedule or view changed")
	}
	if err := instB.Start(context.Background(), nil); err == nil {
		t.Fatal("second start succeeded")
	}
}

func transportClosed(inst *instance) bool {
	raw, _ := inst.substrate.Transport.(*simulatedTransport)
	return raw != nil && raw.closed()
}

func waitForDispatcher(t *testing.T, inst *instance) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for inst.acquired.DispatcherWaiting() != 1 {
		if time.Now().After(deadline) {
			t.Fatal("invoke did not reach the dispatcher lock")
		}
		runtime.Gosched()
	}
}

type simulatedTransport struct {
	mu     sync.Mutex
	isDone bool
	writes atomic.Int32
	done   chan struct{}
}

func newSimulatedTransport() *simulatedTransport {
	return &simulatedTransport{done: make(chan struct{})}
}

func (s *simulatedTransport) ReadByte() (byte, error) {
	<-s.done
	return 0, io.EOF
}

func (s *simulatedTransport) Write(payload []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.isDone {
		return 0, io.ErrClosedPipe
	}
	s.writes.Add(1)
	return len(payload), nil
}

func (s *simulatedTransport) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.isDone {
		return nil
	}
	s.isDone = true
	close(s.done)
	return nil
}

func (s *simulatedTransport) Writes() int { return int(s.writes.Load()) }

func (s *simulatedTransport) closed() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.isDone
}
