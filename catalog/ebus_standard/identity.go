package ebus_standard_catalog

// TelegramClass enumerates eBUS telegram classes supported by the catalog.
type TelegramClass string

// Telegram classes.
const (
	TelegramClassAddressed           TelegramClass = "addressed"
	TelegramClassBroadcast           TelegramClass = "broadcast"
	TelegramClassInitiatorInitiator  TelegramClass = "initiator_initiator"
	TelegramClassControllerBroadcast TelegramClass = "controller_broadcast"
)

// Direction enumerates the request/response axis.
type Direction string

// Directions.
const (
	DirectionRequest  Direction = "request"
	DirectionResponse Direction = "response"
)

// RequestOrResponseRole enumerates the initiator-or-responder role axis.
type RequestOrResponseRole string

// Roles.
const (
	RoleInitiator  RequestOrResponseRole = "initiator"
	RoleResponder  RequestOrResponseRole = "responder"
	RoleOriginator RequestOrResponseRole = "originator"
)

// BroadcastOrAddressed enumerates the addressing axis.
type BroadcastOrAddressed string

// Addressing axis values.
const (
	AddressedDirect    BroadcastOrAddressed = "addressed"
	AddressedBroadcast BroadcastOrAddressed = "broadcast"
)

// AnswerPolicy enumerates the answer-required vs no-answer axis.
type AnswerPolicy string

// Answer policy values.
const (
	AnswerRequired AnswerPolicy = "answer_required"
	AnswerNone     AnswerPolicy = "no_answer"
)

// LengthPrefixMode enumerates the length-prefix-mode axis for the identity
// key. Matches the locked plan's explicit length-selector disambiguation
// requirement.
type LengthPrefixMode string

// Length prefix modes.
const (
	LengthPrefixNone         LengthPrefixMode = "none"
	LengthPrefixFixed        LengthPrefixMode = "fixed"
	LengthPrefixByte         LengthPrefixMode = "byte_prefix"
	LengthPrefixTypedPayload LengthPrefixMode = "typed_payload"
)

// SafetyClass enumerates the execution-safety classes defined by the plan.
type SafetyClass string

// Safety classes.
const (
	SafetyReadOnlySafe    SafetyClass = "read_only_safe"
	SafetyReadOnlyBusLoad SafetyClass = "read_only_bus_load"
	SafetyMutating        SafetyClass = "mutating"
	SafetyDestructive     SafetyClass = "destructive"
	SafetyBroadcast       SafetyClass = "broadcast"
	SafetyMemoryWrite     SafetyClass = "memory_write"
)

// IdentityKey is the full 14-tuple catalog identity key specified in
// canonical §3 of the locked plan. Every catalog entry must populate every
// field; generation fails on duplicate identity keys.
type IdentityKey struct {
	Namespace string `yaml:"namespace"`
	// PB and SB are pointer-typed so the YAML loader can distinguish an
	// absent key (nil → ErrIncompleteIdentityKey) from an explicit zero
	// value (e.g. `pb: 0x00` is legitimate and must be accepted).
	PB                              *uint8                `yaml:"pb"`
	SB                              *uint8                `yaml:"sb"`
	SelectorPath                    string                `yaml:"selector_path"` // nullable ("" = no selector)
	TelegramClass                   TelegramClass         `yaml:"telegram_class"`
	Direction                       Direction             `yaml:"direction"`
	RequestOrResponseRole           RequestOrResponseRole `yaml:"request_or_response_role"`
	BroadcastOrAddressed            BroadcastOrAddressed  `yaml:"broadcast_or_addressed"`
	AnswerPolicy                    AnswerPolicy          `yaml:"answer_policy"`
	LengthPrefixMode                LengthPrefixMode      `yaml:"length_prefix_mode"`
	SelectorDecoder                 string                `yaml:"selector_decoder"`
	ServiceVariant                  string                `yaml:"service_variant"`
	TransportCapabilityRequirements []string              `yaml:"transport_capability_requirements"`
	Version                         string                `yaml:"version"`
}

// IsComplete returns true when every field of the 14-tuple identity key has
// been populated. Selector path is intentionally nullable; the empty string
// is a valid value and does NOT make the key incomplete.
//
// The test suite in identity_completeness_test.go enforces this contract for
// every catalog entry.
func (k IdentityKey) IsComplete() bool {
	if k.Namespace == "" {
		return false
	}
	// PB/SB must be present in the source document. The value 0x00 is
	// legitimate when explicitly set; only absence (nil) is rejected.
	if k.PB == nil {
		return false
	}
	if k.SB == nil {
		return false
	}
	if k.TelegramClass == "" {
		return false
	}
	if k.Direction == "" {
		return false
	}
	if k.RequestOrResponseRole == "" {
		return false
	}
	if k.BroadcastOrAddressed == "" {
		return false
	}
	if k.AnswerPolicy == "" {
		return false
	}
	if k.LengthPrefixMode == "" {
		return false
	}
	if k.SelectorDecoder == "" {
		return false
	}
	if k.ServiceVariant == "" {
		return false
	}
	if k.TransportCapabilityRequirements == nil {
		return false
	}
	if k.Version == "" {
		return false
	}
	// SelectorPath is value-typed; empty is explicitly nullable.
	return true
}

// PBValue returns the dereferenced PB byte, or 0 if PB is nil. Callers that
// need to distinguish "absent" from "explicit 0x00" must check the pointer
// field directly or use IsComplete().
func (k IdentityKey) PBValue() uint8 {
	if k.PB == nil {
		return 0
	}
	return *k.PB
}

// SBValue returns the dereferenced SB byte, or 0 if SB is nil.
func (k IdentityKey) SBValue() uint8 {
	if k.SB == nil {
		return 0
	}
	return *k.SB
}
