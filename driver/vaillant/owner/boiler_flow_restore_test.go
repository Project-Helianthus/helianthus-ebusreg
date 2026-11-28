package owner

import (
	"context"
	"testing"
	"time"

	"github.com/Project-Helianthus/helianthus-ebusgo/engine"
	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
	b503 "github.com/Project-Helianthus/helianthus-ebusgo/protocol/vaillant/b503"
)

func TestVaillantOwnerBoilerFlowRestoredFromReplyAndCache(t *testing.T) {
	persisted := time.Date(2030, 4, 10, 8, 0, 0, 0, time.UTC)
	liveAt := time.Date(2030, 4, 10, 10, 0, 0, 0, time.UTC)
	cached := 36.0
	poller := newBoilerFlowPoller(liveAt, func(frame protocol.Frame) (*protocol.Frame, error) {
		if b509Address(frame) != boilerB509FlowTemperature {
			return nil, errBoilerFlowUnread
		}
		return b509ValueFrame([]byte{0x80, 0x02}), nil
	})
	poller.HydrateFromCache(VaillantCacheHydration{
		boiler:      &vaillantBoilerSnapshot{FlowTemperatureC: &cached},
		persistedAt: persisted,
	})
	if poller.boiler == nil || poller.boiler.FlowTemperatureC == nil || *poller.boiler.FlowTemperatureC != 36 || !poller.boilerFlowCached || !poller.boilerFlowObservedAt.Equal(persisted) {
		t.Fatalf("cache restore flow=%v cached=%v observed=%s", poller.boiler.FlowTemperatureC, poller.boilerFlowCached, poller.boilerFlowObservedAt)
	}

	poller.refreshBoilerStatusTier(context.Background(), boilerStatusTierFast)
	if poller.boiler.FlowTemperatureC == nil || *poller.boiler.FlowTemperatureC != 40 || poller.boilerFlowCached || !poller.boilerFlowObservedAt.Equal(liveAt) {
		t.Fatalf("live flow=%v cached=%v observed=%s", poller.boiler.FlowTemperatureC, poller.boilerFlowCached, poller.boilerFlowObservedAt)
	}

	invalidAt := liveAt.Add(time.Hour)
	invalid := newBoilerFlowPoller(invalidAt, func(frame protocol.Frame) (*protocol.Frame, error) {
		if b509Address(frame) != boilerB509FlowTemperature {
			return nil, errBoilerFlowUnread
		}
		return b509ValueFrame([]byte{0x00, 0x80}), nil
	})
	invalid.HydrateFromCache(VaillantCacheHydration{
		boiler:      &vaillantBoilerSnapshot{FlowTemperatureC: &cached},
		persistedAt: persisted,
	})
	invalid.refreshBoilerStatusTier(context.Background(), boilerStatusTierFast)
	if invalid.boiler.FlowTemperatureC == nil || *invalid.boiler.FlowTemperatureC != 36 || !invalid.boilerFlowCached || !invalid.boilerFlowObservedAt.Equal(persisted) {
		t.Fatalf("invalid marker flow=%v cached=%v observed=%s", invalid.boiler.FlowTemperatureC, invalid.boilerFlowCached, invalid.boilerFlowObservedAt)
	}

	system := 42.5
	mirror := &vaillantSemanticPoller{
		controller: 0x15,
		system:     &vaillantSystemSnapshot{Controller: 0x15, SystemFlowTemperature: &system},
	}
	empty := &vaillantBoilerSnapshot{}
	mirror.mergeBoilerMirrorFieldsFromSnapshots(empty)
	if empty.FlowTemperatureC != nil {
		t.Fatalf("system flow filled boiler flow with %v", *empty.FlowTemperatureC)
	}
	kept := 40.0
	seeded := &vaillantBoilerSnapshot{FlowTemperatureC: &kept}
	mirror.mergeBoilerMirrorFieldsFromSnapshots(seeded)
	if seeded.FlowTemperatureC == nil || *seeded.FlowTemperatureC != 40 {
		t.Fatalf("system flow replaced boiler flow with %v", seeded.FlowTemperatureC)
	}
}

func TestVaillantOwnerCurrentErrorOnSlowTier(t *testing.T) {
	poller := &vaillantSemanticPoller{boilerSlowInterval: 17 * time.Minute}
	var task semanticPollTask
	found := false
	for _, item := range poller.pollSchedule().recurring {
		if item.key == semanticTaskRefreshCurrentError {
			task = item
			found = true
		}
	}
	if !found || task.interval != 17*time.Minute || task.priority != semanticTaskPriorityLow || task.refresh == nil {
		t.Fatalf("current-error task found=%v interval=%v priority=%v", found, task.interval, task.priority)
	}

	active := newCurrentErrorPoller(func(target byte) ([]byte, error) {
		if target == 0x15 {
			return currentErrorPayload(4, 25), nil
		}
		return currentErrorPayload(-1, 0), nil
	})
	active.boilerSlowInterval = 17 * time.Minute
	for _, item := range active.pollSchedule().recurring {
		if item.key == semanticTaskRefreshCurrentError {
			item.refresh(context.Background())
		}
	}
	if active.currentErrorFault != "present" || active.currentErrorFrames != 2 {
		t.Fatalf("active fault=%q frames=%d", active.currentErrorFault, active.currentErrorFrames)
	}

	clear := newCurrentErrorPoller(func(byte) ([]byte, error) {
		return currentErrorPayload(-1, 0), nil
	})
	clear.refreshCurrentError(context.Background())
	if clear.currentErrorFault != "clear" || clear.currentErrorFrames != 2 {
		t.Fatalf("clear fault=%q frames=%d", clear.currentErrorFault, clear.currentErrorFrames)
	}

	failed := newCurrentErrorPoller(func(byte) ([]byte, error) {
		return nil, errBoilerFlowUnread
	})
	failed.currentErrorFault = "present"
	failed.refreshCurrentError(context.Background())
	if failed.currentErrorFault != "present" || failed.currentErrorFrames != 2 {
		t.Fatalf("failed read fault=%q frames=%d", failed.currentErrorFault, failed.currentErrorFrames)
	}

	one := newCurrentErrorPoller(func(target byte) ([]byte, error) {
		if target != 0x15 {
			t.Fatalf("addressed absent boiler target %#x", target)
		}
		return currentErrorPayload(-1, 0), nil
	})
	one.boilerAddress = 0
	one.refreshCurrentError(context.Background())
	if one.currentErrorFault != "clear" || one.currentErrorFrames != 1 {
		t.Fatalf("absent boiler fault=%q frames=%d", one.currentErrorFault, one.currentErrorFrames)
	}
}

var errBoilerFlowUnread = context.DeadlineExceeded

func newBoilerFlowPoller(now time.Time, reply func(protocol.Frame) (*protocol.Frame, error)) *vaillantSemanticPoller {
	return &vaillantSemanticPoller{
		scheduler:      engine.NewSemanticReadScheduler(),
		source:         0x7F,
		boilerAddress:  0x08,
		requestTimeout: 50 * time.Millisecond,
		nowFn:          func() time.Time { return now },
		pauseFn:        func(context.Context, time.Duration) bool { return true },
		sendFrameFn: func(_ context.Context, frame protocol.Frame) (*protocol.Frame, error) {
			return reply(frame)
		},
	}
}

func b509Address(frame protocol.Frame) uint16 {
	if len(frame.Data) < 3 {
		return 0
	}
	return uint16(frame.Data[1])<<8 | uint16(frame.Data[2])
}

func b509ValueFrame(payload []byte) *protocol.Frame {
	data := []byte{vaillantB509OpcodeRead, 0x18, 0x00}
	return &protocol.Frame{Data: append(data, payload...)}
}

type currentErrorReply func(target byte) ([]byte, error)

type currentErrorPoller struct {
	*vaillantSemanticPoller
	currentErrorFrames int
}

func newCurrentErrorPoller(reply currentErrorReply) *currentErrorPoller {
	poller := &currentErrorPoller{vaillantSemanticPoller: &vaillantSemanticPoller{
		source: 0x7F, controller: 0x15, boilerAddress: 0x08, requestTimeout: 50 * time.Millisecond,
	}}
	poller.sendFrameFn = func(_ context.Context, frame protocol.Frame) (*protocol.Frame, error) {
		poller.currentErrorFrames++
		if frame.Source != 0x7F || frame.Primary != 0xB5 || frame.Secondary != 0x03 || string(frame.Data) != string(b503.EncodeCurrentError()) {
			return nil, errBoilerFlowUnread
		}
		payload, err := reply(frame.Target)
		if err != nil {
			return nil, err
		}
		return &protocol.Frame{Data: payload}, nil
	}
	return poller
}

func currentErrorPayload(active int, code uint16) []byte {
	payload := []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	if active >= 0 {
		payload[active*2] = byte(code)
		payload[active*2+1] = byte(code >> 8)
	}
	return payload
}
