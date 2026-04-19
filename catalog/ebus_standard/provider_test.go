package ebus_standard_catalog

import (
	"context"
	"errors"
	"strings"
	"testing"
)

// TestProvider_ConstructsFromCatalog confirms that NewProvider accepts the
// embedded catalog and exposes Identification as a callable method. RED
// phase: the stub returns an error.
func TestProvider_ConstructsFromCatalog(t *testing.T) {
	cat := MustEmbeddedCatalog()
	p, err := NewProvider(cat, true)
	if err != nil {
		t.Fatalf("NewProvider: %v", err)
	}
	if p == nil {
		t.Fatal("NewProvider returned nil")
	}
	if !p.IsEnabled() {
		t.Fatal("expected provider to be enabled when constructed with enabled=true")
	}
}

// TestProvider_IdentificationReturnsDescriptor exercises the Identification
// entrypoint on the enabled provider. The descriptor MUST carry the
// ebus_standard.identification source label so downstream provenance code
// cannot accidentally overwrite DeviceInfo.
func TestProvider_IdentificationReturnsDescriptor(t *testing.T) {
	p, err := NewProvider(MustEmbeddedCatalog(), true)
	if err != nil {
		t.Fatalf("NewProvider: %v", err)
	}
	desc, err := p.Identification(context.Background())
	if err != nil {
		t.Fatalf("Identification returned error: %v", err)
	}
	if desc.Source != string(SourceEBUSStandardIdent) {
		t.Fatalf("descriptor.Source=%q, want %q", desc.Source, SourceEBUSStandardIdent)
	}
	if desc.CatalogVersion != CatalogVersion {
		t.Fatalf("descriptor.CatalogVersion=%q, want %q", desc.CatalogVersion, CatalogVersion)
	}
}

// TestProvider_DisabledReturnsSentinel confirms that a provider constructed
// in the disabled state refuses every entrypoint with ErrProviderDisabled.
func TestProvider_DisabledReturnsSentinel(t *testing.T) {
	p, err := NewProvider(MustEmbeddedCatalog(), false)
	if err != nil {
		t.Fatalf("NewProvider: %v", err)
	}
	if p.IsEnabled() {
		t.Fatal("expected IsEnabled() to be false")
	}
	if _, err := p.Identification(context.Background()); !errors.Is(err, ErrProviderDisabled) {
		t.Fatalf("Identification err=%v, want ErrProviderDisabled", err)
	}
	if _, err := p.Invoke(context.Background(), "ebus_standard.service_data.start_counts", nil, CallerContextUserFacing); !errors.Is(err, ErrProviderDisabled) {
		t.Fatalf("Invoke err=%v, want ErrProviderDisabled", err)
	}
}

// TestProvider_FromEnv_DefaultEnabled confirms the env-var default (unset =
// enabled).
func TestProvider_FromEnv_DefaultEnabled(t *testing.T) {
	t.Setenv(DisableEnvVar, "")
	// An empty string MUST be treated as "unset" => enabled by default.
	// Use Unsetenv path via Setenv to "" then explicit unset.
	p, err := NewProviderFromEnv(MustEmbeddedCatalog())
	if err != nil {
		t.Fatalf("NewProviderFromEnv: %v", err)
	}
	if !p.IsEnabled() {
		t.Fatal("empty env var should default to enabled")
	}
}

// TestProvider_FromEnv_DisabledByFalse confirms that common "false" spellings
// disable the provider.
func TestProvider_FromEnv_DisabledByFalse(t *testing.T) {
	for _, v := range []string{"0", "false", "False", "FALSE"} {
		t.Run(v, func(t *testing.T) {
			t.Setenv(DisableEnvVar, v)
			p, err := NewProviderFromEnv(MustEmbeddedCatalog())
			if err != nil {
				t.Fatalf("NewProviderFromEnv: %v", err)
			}
			if p.IsEnabled() {
				t.Fatalf("env=%q should produce disabled provider", v)
			}
		})
	}
}

// TestProvider_UnknownMethodID confirms Invoke distinguishes "not in
// catalog" from "safety denied".
// TestProvider_DuplicateMethodIDRejected asserts NewProvider surfaces
// ErrDuplicateMethodID when two catalog commands share the same ID. A
// previous implementation silently overwrote the first entry in the
// method index, which would cause Invoke to dispatch the wrong command.
// The error message must name BOTH colliding commands.
func TestProvider_DuplicateMethodIDRejected(t *testing.T) {
	pb := uint8(0x03)
	sbA := uint8(0xF0)
	sbB := uint8(0xF1)
	mkCmd := func(name string, sb uint8) Command {
		return Command{
			ID:   "ebus_standard.collision.method",
			Name: name,
			Identity: IdentityKey{
				Namespace:                       Namespace,
				PB:                              &pb,
				SB:                              &sb,
				TelegramClass:                   TelegramClassAddressed,
				Direction:                       DirectionRequest,
				RequestOrResponseRole:           RoleInitiator,
				BroadcastOrAddressed:            AddressedDirect,
				AnswerPolicy:                    AnswerRequired,
				LengthPrefixMode:                LengthPrefixNone,
				SelectorDecoder:                 "none",
				ServiceVariant:                  "synthetic",
				TransportCapabilityRequirements: []string{"master_slave"},
				Version:                         CatalogVersion,
			},
			SafetyClass: SafetyReadOnlySafe,
		}
	}
	cat := Catalog{
		Namespace:  Namespace,
		Version:    CatalogVersion,
		PlanSHA256: CanonicalPlanSHA256,
		Services: []Service{
			{
				PB:   &pb,
				Name: "synthetic",
				Commands: []Command{
					mkCmd("first", sbA),
					mkCmd("second", sbB),
				},
			},
		},
	}
	p, err := NewProvider(cat, true)
	if !errors.Is(err, ErrDuplicateMethodID) {
		t.Fatalf("err=%v, want ErrDuplicateMethodID", err)
	}
	if p != nil {
		t.Fatalf("expected nil provider on duplicate, got %+v", p)
	}
	msg := err.Error()
	if !strings.Contains(msg, "first") || !strings.Contains(msg, "second") {
		t.Fatalf("error message %q must name both colliding commands", msg)
	}
}

func TestProvider_UnknownMethodID(t *testing.T) {
	p, err := NewProvider(MustEmbeddedCatalog(), true)
	if err != nil {
		t.Fatalf("NewProvider: %v", err)
	}
	_, err = p.Invoke(context.Background(), "ebus_standard.does_not_exist", nil, CallerContextUserFacing)
	if !errors.Is(err, ErrUnknownMethod) {
		t.Fatalf("err=%v, want ErrUnknownMethod", err)
	}
}
