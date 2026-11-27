package driver

import (
	"context"
	"errors"
	"sync"

	semdriver "github.com/Project-Helianthus/helianthus-semreg/semreg/v1/driver"

	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/owner"
)

var _ semdriver.Instance = (*instance)(nil)

type instance struct {
	key  semdriver.InstanceKey
	cfg  owner.Config
	open openFunc

	mu          sync.Mutex
	started     bool
	holdPolling bool
	cancel      context.CancelFunc
	runCtx      context.Context
	acquired    *owner.Acquisition
	substrate   owner.Substrate
}

// Start opens the instance transport and starts its owner. It dials nothing itself.
func (i *instance) Start(ctx context.Context, _ semdriver.Session) error {
	if i == nil {
		return errors.New("nil instance")
	}
	i.mu.Lock()
	if i.started {
		i.mu.Unlock()
		return errors.New("instance already started")
	}
	if i.open == nil {
		i.mu.Unlock()
		return errors.New("ebus transport is not open")
	}
	cfg, open, hold := i.cfg, i.open, i.holdPolling
	i.mu.Unlock()

	if ctx == nil {
		ctx = context.Background()
	}
	substrate, err := open(ctx, cfg)
	if err != nil {
		return err
	}
	acquired := owner.NewAcquisition(cfg, substrate)
	if hold {
		acquired.HoldStartup()
	}
	runCtx, cancel := context.WithCancel(context.WithoutCancel(ctx))
	acquired.Begin(runCtx)

	i.mu.Lock()
	defer i.mu.Unlock()
	if i.started {
		cancel()
		closeSubstrate(substrate)
		return errors.New("instance already started")
	}
	i.started = true
	i.cancel = cancel
	i.runCtx = runCtx
	i.acquired = acquired
	i.substrate = substrate
	return nil
}

// FenceWithdrawal suspends the running owner without a supervisor callback.
func (i *instance) FenceWithdrawal() {
	if i == nil {
		return
	}
	i.mu.Lock()
	acquired := i.acquired
	i.mu.Unlock()
	if acquired != nil {
		acquired.Suspend()
	}
}

// Stop cancels the owner and closes its transport.
func (i *instance) Stop(context.Context) semdriver.CloseReport {
	if i == nil {
		return semdriver.CloseReport{Confirmed: false, Err: errors.New("nil instance")}
	}
	i.mu.Lock()
	cancel := i.cancel
	transport := i.substrate.Transport
	i.cancel = nil
	i.mu.Unlock()
	if cancel != nil {
		cancel()
	}
	if transport == nil {
		return semdriver.CloseReport{Confirmed: true}
	}
	if err := transport.Close(); err != nil {
		return semdriver.CloseReport{Err: err}
	}
	return semdriver.CloseReport{Confirmed: true}
}

func (i *instance) running() bool {
	if i == nil {
		return false
	}
	i.mu.Lock()
	defer i.mu.Unlock()
	return i.runCtx != nil && i.runCtx.Err() == nil
}

func closeSubstrate(substrate owner.Substrate) {
	if substrate.Transport != nil {
		_ = substrate.Transport.Close()
	}
}
