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

	DefaultSemanticZonePresenceMissThreshold = 3
	DefaultSemanticZonePresenceHitThreshold  = 2
	DefaultSemanticDHWStaleTTL               = 15 * time.Minute
	DefaultSemanticEnergyInterval            = 5 * time.Minute
	DefaultSemanticRegulatorRecheckInterval  = 60 * time.Second
	DefaultSemanticRegulatorAbsenceGrace     = 5 * time.Minute
)

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

type Substrate struct {
	Registry            *registry.DeviceRegistry
	Bus                 *protocol.Bus
	Transport           transport.RawTransport
	RefreshRouterPlanes func() int
}
