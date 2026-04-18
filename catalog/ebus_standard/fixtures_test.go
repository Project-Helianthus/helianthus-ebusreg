package ebus_standard_catalog

import (
	"testing"
)

// TestEmbeddedCatalog_Class2BroadcastFixture asserts the embedded catalog
// contains at least one class-2 broadcast entry (telegram_class=broadcast,
// no_answer). This is one of the axes explicitly called out in locked plan
// §3 risk 1.
func TestEmbeddedCatalog_Class2BroadcastFixture(t *testing.T) {
	cat := MustEmbeddedCatalog()
	for _, svc := range cat.Services {
		for _, cmd := range svc.Commands {
			if cmd.Identity.TelegramClass == TelegramClassBroadcast &&
				cmd.Identity.AnswerPolicy == AnswerNone {
				return
			}
		}
	}
	t.Fatalf("no class-2 broadcast (broadcast + no_answer) entry in embedded catalog")
}

// TestEmbeddedCatalog_OriginatorBroadcastModeFixture asserts the catalog
// contains at least one originator broadcast (controller-to-all broadcast
// mode, i.e. an originator broadcasting to every participant without
// expecting an addressed response). `identification_self_broadcast` meets
// this.
func TestEmbeddedCatalog_OriginatorBroadcastModeFixture(t *testing.T) {
	cat := MustEmbeddedCatalog()
	for _, svc := range cat.Services {
		for _, cmd := range svc.Commands {
			if cmd.Identity.RequestOrResponseRole == RoleOriginator &&
				cmd.Identity.BroadcastOrAddressed == AddressedBroadcast {
				return
			}
		}
	}
	t.Fatalf("no originator broadcast-mode entry (originator + broadcast)")
}

// TestEmbeddedCatalog_NoAnswerFormFixture asserts at least one addressed
// no-answer form exists (destructive or broadcast-initiated without
// responder reply). `ebus_standard.burner.barred_reserved` covers this.
func TestEmbeddedCatalog_NoAnswerFormFixture(t *testing.T) {
	cat := MustEmbeddedCatalog()
	for _, svc := range cat.Services {
		for _, cmd := range svc.Commands {
			if cmd.Identity.AnswerPolicy == AnswerNone {
				return
			}
		}
	}
	t.Fatalf("no no-answer-form entry")
}

// TestEmbeddedCatalog_TypedPayloadSelectorFixture asserts the catalog
// contains at least one typed-payload-selector entry (length_prefix_mode =
// typed_payload, selector_decoder != "none"). The 0x05 0x03 block-number
// variants cover this.
func TestEmbeddedCatalog_TypedPayloadSelectorFixture(t *testing.T) {
	cat := MustEmbeddedCatalog()
	for _, svc := range cat.Services {
		for _, cmd := range svc.Commands {
			if cmd.Identity.LengthPrefixMode == LengthPrefixTypedPayload &&
				cmd.Identity.SelectorDecoder != "none" &&
				cmd.Identity.SelectorDecoder != "" {
				return
			}
		}
	}
	t.Fatalf("no typed-payload-selector entry")
}

// TestLoadCatalog_IncompleteIdentityKey asserts that an otherwise valid
// YAML with a missing identity field is rejected with
// ErrIncompleteIdentityKey.
func TestLoadCatalog_IncompleteIdentityKey(t *testing.T) {
	yml := []byte(`
namespace: ebus_standard
version: v1.0-locked
plan_sha256: 9e0a29bb76d99f551904b05749e322aafd3972621858aa6d1acbe49b9ef37305
services:
  - pb: 0x07
    name: System Data
    commands:
      - id: ebus_standard.incomplete
        name: Incomplete entry
        identity:
          namespace: ebus_standard
          pb: 0x07
          sb: 0x04
          selector_path: ""
          telegram_class: addressed
          direction: request
          # request_or_response_role intentionally omitted
          broadcast_or_addressed: addressed
          answer_policy: answer_required
          length_prefix_mode: fixed
          selector_decoder: none
          service_variant: test_variant
          transport_capability_requirements: [master_slave]
          version: v1.0-locked
        safety_class: read_only_bus_load
`)
	_, err := LoadCatalog(yml)
	if err == nil {
		t.Fatalf("LoadCatalog: expected ErrIncompleteIdentityKey, got nil")
	}
}

// TestLoadCatalog_ContentSHARoundtrip asserts that SHA-256 of the raw bytes
// is deterministic and populated on load.
func TestLoadCatalog_ContentSHARoundtrip(t *testing.T) {
	cat1 := MustEmbeddedCatalog()
	cat2 := MustEmbeddedCatalog()
	if cat1.ContentSHA256 != cat2.ContentSHA256 {
		t.Fatalf("non-deterministic ContentSHA256: %q vs %q",
			cat1.ContentSHA256, cat2.ContentSHA256)
	}
	if cat1.ContentSHA256 == "" {
		t.Fatalf("ContentSHA256 empty")
	}
}
