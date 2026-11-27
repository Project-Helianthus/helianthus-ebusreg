package driver

import (
	"bytes"
	"encoding/json"
	"errors"
	"time"

	"github.com/Project-Helianthus/helianthus-ebusgo/engine"

	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/owner"
)

type durationText time.Duration

func (d *durationText) UnmarshalJSON(raw []byte) error {
	if string(raw) == "null" {
		return nil
	}
	var text string
	if err := json.Unmarshal(raw, &text); err != nil {
		return err
	}
	if text == "" {
		*d = 0
		return nil
	}
	parsed, err := time.ParseDuration(text)
	if err != nil {
		return err
	}
	*d = durationText(parsed)
	return nil
}

type transportDocument struct {
	ScanSource      byte         `json:"scan_source"`
	CompanionTarget byte         `json:"companion_target"`
	Protocol        string       `json:"protocol"`
	Network         string       `json:"network"`
	Address         string       `json:"address"`
	ReadTimeout     durationText `json:"read_timeout"`
	WriteTimeout    durationText `json:"write_timeout"`
	DialTimeout     durationText `json:"dial_timeout"`
}

type optionsDocument struct {
	StateInterval     durationText `json:"state_interval"`
	DiscoveryInterval durationText `json:"discovery_interval"`
	ConfigInterval    durationText `json:"config_interval"`
	EnergyInterval    durationText `json:"energy_interval"`
	RequestTimeout    durationText `json:"request_timeout"`
	HotWaterStale     durationText `json:"dhw_stale_ttl"`
	RegulatorRecheck  durationText `json:"regulator_recheck_interval"`
	RegulatorAbsence  durationText `json:"regulator_absence_grace"`
	BreakerCooldown   durationText `json:"read_breaker_open_cooldown"`
	ZoneMissThreshold int          `json:"zone_presence_miss_threshold"`
	ZoneHitThreshold  int          `json:"zone_presence_hit_threshold"`
	BreakerBudget     int          `json:"read_breaker_failure_budget"`
	BreakerProbes     int          `json:"read_breaker_half_open_probe_limit"`
	ObserveFirst      bool         `json:"observe_first"`
	PassiveState      bool         `json:"passive_state_direct_apply"`
	PassiveConfig     bool         `json:"passive_config_direct_apply"`
	ExternalWrite     string       `json:"external_write_policy"`
}

func decodeOwnerConfig(transportJSON, optionsJSON json.RawMessage) (owner.Config, error) {
	var transportDoc transportDocument
	if err := decodeObject(transportJSON, &transportDoc); err != nil {
		return owner.Config{}, err
	}
	var optionsDoc optionsDocument
	if err := decodeObject(optionsJSON, &optionsDoc); err != nil {
		return owner.Config{}, err
	}
	return ownerConfig(transportDoc, optionsDoc), nil
}

func decodeObject(raw json.RawMessage, dest any) error {
	if len(bytes.TrimSpace(raw)) == 0 || string(bytes.TrimSpace(raw)) == "null" {
		return nil
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	if err := dec.Decode(dest); err != nil {
		return err
	}
	if dec.More() {
		return errors.New("ebus config has trailing json")
	}
	return nil
}

func ownerConfig(transportDoc transportDocument, optionsDoc optionsDocument) owner.Config {
	cfg := owner.Config{
		ScanSource:             transportDoc.ScanSource,
		StartupCompanionTarget: transportDoc.CompanionTarget,
		TransportConfig: engine.TransportConfig{
			Protocol:     engine.TransportProtocol(transportDoc.Protocol),
			Network:      transportDoc.Network,
			Address:      transportDoc.Address,
			ReadTimeout:  time.Duration(transportDoc.ReadTimeout),
			WriteTimeout: time.Duration(transportDoc.WriteTimeout),
			DialTimeout:  time.Duration(transportDoc.DialTimeout),
		},
		SemanticStateInterval:                 time.Duration(optionsDoc.StateInterval),
		SemanticDiscoveryInterval:             time.Duration(optionsDoc.DiscoveryInterval),
		SemanticConfigInterval:                time.Duration(optionsDoc.ConfigInterval),
		SemanticEnergyInterval:                time.Duration(optionsDoc.EnergyInterval),
		SemanticRequestTimeout:                time.Duration(optionsDoc.RequestTimeout),
		SemanticDHWStaleTTL:                   time.Duration(optionsDoc.HotWaterStale),
		SemanticRegulatorRecheckInterval:      time.Duration(optionsDoc.RegulatorRecheck),
		SemanticRegulatorAbsenceGrace:         time.Duration(optionsDoc.RegulatorAbsence),
		SemanticReadBreakerOpenCooldown:       time.Duration(optionsDoc.BreakerCooldown),
		SemanticZonePresenceMissThreshold:     optionsDoc.ZoneMissThreshold,
		SemanticZonePresenceHitThreshold:      optionsDoc.ZoneHitThreshold,
		SemanticReadBreakerFailureBudget:      optionsDoc.BreakerBudget,
		SemanticReadBreakerHalfOpenProbeLimit: optionsDoc.BreakerProbes,
		ObserveFirstEnabled:                   optionsDoc.ObserveFirst,
		PassiveStateDirectApply:               optionsDoc.PassiveState,
		PassiveConfigDirectApply:              optionsDoc.PassiveConfig,
		ExternalWritePolicy:                   engine.ObserveFirstExternalWritePolicy(optionsDoc.ExternalWrite),
	}
	return cfg
}
