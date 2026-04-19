package ebus_standard_catalog

import (
	"context"
	"errors"
	"testing"
)

// TestProvider_ConstructsFromCatalog confirms that NewProvider accepts the
// embedded catalog and exposes Identification as a callable method. RED
// phase: the stub returns an error.
func TestProvider_ConstructsFromCatalog(t *testing.T) {
	cat := MustEmbeddedCatalog()
	p := NewProvider(cat, true)
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
	p := NewProvider(MustEmbeddedCatalog(), true)
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
	p := NewProvider(MustEmbeddedCatalog(), false)
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
	p := NewProviderFromEnv(MustEmbeddedCatalog())
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
			p := NewProviderFromEnv(MustEmbeddedCatalog())
			if p.IsEnabled() {
				t.Fatalf("env=%q should produce disabled provider", v)
			}
		})
	}
}

// TestProvider_UnknownMethodID confirms Invoke distinguishes "not in
// catalog" from "safety denied".
func TestProvider_UnknownMethodID(t *testing.T) {
	p := NewProvider(MustEmbeddedCatalog(), true)
	_, err := p.Invoke(context.Background(), "ebus_standard.does_not_exist", nil, CallerContextUserFacing)
	if !errors.Is(err, ErrUnknownMethod) {
		t.Fatalf("err=%v, want ErrUnknownMethod", err)
	}
}
