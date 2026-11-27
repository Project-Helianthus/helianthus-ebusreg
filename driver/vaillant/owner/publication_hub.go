package owner

import (
	"context"
	"sync"

	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/output"
)

const publicationHubBufferSize = 16

type publicationHubSubscriber struct {
	ch        chan interface{}
	closeOnce sync.Once
}

func (sub *publicationHubSubscriber) close() {
	sub.closeOnce.Do(func() { close(sub.ch) })
}

// publicationEvents is the neutral sink the publication owner emits its
// published updates to; a consumer connects it to its own broadcast path.
type publicationEvents interface {
	PublishZoneUpdate(zone output.Zone)
	PublishDHWUpdate(status *output.DhwStatus)
	PublishEnergyUpdate(totals *output.EnergyTotals)
	PublishBoilerStatusUpdate(status *output.BoilerStatus)
	PublishRadioDevicesUpdate(devices []output.RadioDevice)
}

// publicationHub fans the published updates out to buffered subscribers; a
// subscriber whose buffer is full misses the update.
type publicationHub struct {
	mu     sync.RWMutex
	zones  map[uint64]*publicationHubSubscriber
	dhw    map[uint64]*publicationHubSubscriber
	energy map[uint64]*publicationHubSubscriber
	boiler map[uint64]*publicationHubSubscriber
	radio  map[uint64]*publicationHubSubscriber
	nextID uint64
}

func newPublicationHub() *publicationHub {
	return &publicationHub{
		zones:  make(map[uint64]*publicationHubSubscriber),
		dhw:    make(map[uint64]*publicationHubSubscriber),
		energy: make(map[uint64]*publicationHubSubscriber),
		boiler: make(map[uint64]*publicationHubSubscriber),
		radio:  make(map[uint64]*publicationHubSubscriber),
	}
}

func (hub *publicationHub) SubscribeZones(ctx context.Context) (chan interface{}, error) {
	return hub.subscribe(ctx, hub.zones)
}

func (hub *publicationHub) SubscribeDHW(ctx context.Context) (chan interface{}, error) {
	return hub.subscribe(ctx, hub.dhw)
}

func (hub *publicationHub) SubscribeEnergy(ctx context.Context) (chan interface{}, error) {
	return hub.subscribe(ctx, hub.energy)
}

func (hub *publicationHub) SubscribeBoiler(ctx context.Context) (chan interface{}, error) {
	return hub.subscribe(ctx, hub.boiler)
}

func (hub *publicationHub) SubscribeRadioDevices(ctx context.Context) (chan interface{}, error) {
	return hub.subscribe(ctx, hub.radio)
}

func (hub *publicationHub) PublishZoneUpdate(zone output.Zone) {
	hub.publish(hub.zones, zone)
}

func (hub *publicationHub) PublishDHWUpdate(status *output.DhwStatus) {
	if status == nil {
		return
	}
	hub.publish(hub.dhw, status)
}

func (hub *publicationHub) PublishEnergyUpdate(totals *output.EnergyTotals) {
	if totals == nil {
		return
	}
	hub.publish(hub.energy, totals)
}

func (hub *publicationHub) PublishBoilerStatusUpdate(status *output.BoilerStatus) {
	if status == nil {
		return
	}
	hub.publish(hub.boiler, status)
}

func (hub *publicationHub) PublishRadioDevicesUpdate(devices []output.RadioDevice) {
	out := make([]output.RadioDevice, len(devices))
	for i, device := range devices {
		out[i] = cloneRadioDevice(device)
	}
	hub.publish(hub.radio, out)
}

func (hub *publicationHub) subscribe(ctx context.Context, subscribers map[uint64]*publicationHubSubscriber) (chan interface{}, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	sub := &publicationHubSubscriber{ch: make(chan interface{}, publicationHubBufferSize)}
	hub.mu.Lock()
	hub.nextID++
	id := hub.nextID
	subscribers[id] = sub
	hub.mu.Unlock()

	go func() {
		<-ctx.Done()
		hub.mu.Lock()
		removed := subscribers[id]
		delete(subscribers, id)
		hub.mu.Unlock()
		if removed != nil {
			removed.close()
		}
	}()
	return sub.ch, nil
}

func (hub *publicationHub) publish(subscribers map[uint64]*publicationHubSubscriber, payload any) {
	hub.mu.RLock()
	list := make([]*publicationHubSubscriber, 0, len(subscribers))
	for _, sub := range subscribers {
		list = append(list, sub)
	}
	hub.mu.RUnlock()
	for _, sub := range list {
		select {
		case sub.ch <- payload:
		default:
		}
	}
}
