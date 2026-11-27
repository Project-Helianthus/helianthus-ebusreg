// Package driver is the public constructor of the eBUS driver.
package driver

import (
	"context"
	"errors"
	"fmt"
	"sync"

	semdriver "github.com/Project-Helianthus/helianthus-semreg/semreg/v1/driver"

	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/owner"
)

// Kind is the driver kind of every instance this factory plans.
const Kind semdriver.DriverKind = "ebusreg"

const (
	transportProvider     = "helianthus-ebusgo"
	transportAccessDomain = "ebus"
)

// Factory is the SemReg constructor of the eBUS driver.
type Factory struct {
	mu      sync.Mutex
	planned map[semdriver.InstanceKey]owner.Config
	open    openFunc
}

type openFunc func(context.Context, owner.Config) (owner.Substrate, error)

var _ semdriver.Factory = (*Factory)(nil)

// NewFactory returns the eBUS driver factory. Start fails closed until a transport is open.
func NewFactory() *Factory {
	return &Factory{
		planned: map[semdriver.InstanceKey]owner.Config{},
		open:    closedTransport,
	}
}

func closedTransport(context.Context, owner.Config) (owner.Substrate, error) {
	return owner.Substrate{}, errors.New("ebus transport is not open")
}

// Descriptor returns the families the factory declares.
func (f *Factory) Descriptor() semdriver.Descriptor { return semdriver.Descriptor{} }

// Plan plans one enabled instance. It dials nothing and does not read SecretRef.
func (f *Factory) Plan(cfg semdriver.Config) ([]semdriver.InstancePlan, error) {
	if f == nil {
		return nil, errors.New("nil factory")
	}
	if err := semdriver.InstanceKey(cfg.Instance).Validate(); err != nil {
		return nil, err
	}
	if !cfg.Enabled {
		return nil, nil
	}
	decoded, err := decodeOwnerConfig(cfg.Transport, cfg.Options)
	if err != nil {
		return nil, err
	}
	plan := semdriver.InstancePlan{
		Driver:   Kind,
		Instance: semdriver.InstanceKey(cfg.Instance),
		Transports: []semdriver.TransportClaim{{
			Provider:     transportProvider,
			ConfigKey:    cfg.Instance,
			AccessDomain: transportAccessDomain,
			Share:        semdriver.ShareSerialized,
		}},
	}
	if err := plan.Validate(); err != nil {
		return nil, err
	}
	f.mu.Lock()
	f.planned[plan.Instance] = decoded
	f.mu.Unlock()
	return []semdriver.InstancePlan{plan}, nil
}

// Construct builds one planned instance. It opens nothing.
func (f *Factory) Construct(_ context.Context, plan semdriver.InstancePlan) (semdriver.Instance, error) {
	if f == nil {
		return nil, errors.New("nil factory")
	}
	if err := plan.Validate(); err != nil {
		return nil, err
	}
	if plan.Driver != Kind {
		return nil, fmt.Errorf("plan is for driver %q, not %q", plan.Driver, Kind)
	}
	f.mu.Lock()
	cfg, ok := f.planned[plan.Instance]
	open := f.open
	f.mu.Unlock()
	if !ok {
		return nil, fmt.Errorf("instance %q was not planned by this factory", plan.Instance)
	}
	return &instance{key: plan.Instance, cfg: cfg, open: open}, nil
}
