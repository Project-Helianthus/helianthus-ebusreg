// Package owner is the Vaillant acquisition owner: poller, schedule, discovery and tier refreshes.
package owner

import (
	"time"

	"github.com/Project-Helianthus/helianthus-ebusgo/engine"
	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
	"github.com/Project-Helianthus/helianthus-ebusgo/transport"
	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
)

const (
	defaultSemanticConfigInterval = 5 * time.Minute

	// DefaultSemanticZonePresenceMissThreshold is the default semantic zone presence miss threshold value of the Vaillant owner.
	DefaultSemanticZonePresenceMissThreshold = 3
	// DefaultSemanticZonePresenceHitThreshold is the default semantic zone presence hit threshold value of the Vaillant owner.
	DefaultSemanticZonePresenceHitThreshold = 2
	// DefaultSemanticDHWStaleTTL is the default semantic dhwstale ttl value of the Vaillant owner.
	DefaultSemanticDHWStaleTTL = 15 * time.Minute
	// DefaultSemanticEnergyInterval is the default semantic energy interval value of the Vaillant owner.
	DefaultSemanticEnergyInterval = 5 * time.Minute
	// DefaultSemanticRegulatorRecheckInterval is the default semantic regulator recheck interval value of the Vaillant owner.
	DefaultSemanticRegulatorRecheckInterval = 60 * time.Second
	// DefaultSemanticRegulatorAbsenceGrace is the default semantic regulator absence grace value of the Vaillant owner.
	DefaultSemanticRegulatorAbsenceGrace = 5 * time.Minute
)

// Config is the config type of the Vaillant owner.
type Config struct {
	TransportConfig         engine.TransportConfig
	WatchObserver           engine.WatchObserver
	WatchEfficiencyObserver engine.WatchEfficiencyObserver
	ScanSource              byte
	StartupCompanionTarget  byte

	ObserveFirstEnabled      bool
	PassiveStateDirectApply  bool
	PassiveConfigDirectApply bool
	ExternalWritePolicy      engine.ObserveFirstExternalWritePolicy
	ObserveFirstFlags        engine.ObserveFirstFeatureFlags

	SemanticRequestTimeout                time.Duration
	SemanticDiscoveryInterval             time.Duration
	SemanticConfigInterval                time.Duration
	SemanticStateInterval                 time.Duration
	SemanticEnergyInterval                time.Duration
	SemanticReadBreakerFailureBudget      int
	SemanticReadBreakerOpenCooldown       time.Duration
	SemanticReadBreakerHalfOpenProbeLimit int
	SemanticZonePresenceMissThreshold     int
	SemanticZonePresenceHitThreshold      int
	SemanticDHWStaleTTL                   time.Duration
	SemanticRegulatorRecheckInterval      time.Duration
	SemanticRegulatorAbsenceGrace         time.Duration
}

// Substrate is the substrate type of the Vaillant owner.
type Substrate struct {
	Registry            *registry.DeviceRegistry
	Bus                 *protocol.Bus
	Transport           transport.RawTransport
	RefreshRouterPlanes func() int
}
