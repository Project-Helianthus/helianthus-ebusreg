# helianthus-ebusreg

`helianthus-ebusreg` is the registry/projection layer of Helianthus. It matches discovered eBUS devices to semantic planes, exposes method metadata + schema selection, and routes request/response + broadcast handling through plane implementations.

## Purpose and Scope

This repository is for:

- engineers building registry, provider, and routing behavior on top of eBUS transport primitives;
- operators validating registry/projection behavior in gateway integrations.

What this repo includes:

- device registry + scan flow (`registry`)
- schema and selector decoding/encoding (`schema`)
- bus event routing/invocation (`router`)
- default provider wiring (`providers/vaillant`)
- Vaillant plane implementations and projection mappings (`vaillant/*`)

What this repo does **not** include:

- low-level transport/framing state machines (that is `helianthus-ebusgo`)
- runnable gateway/API server binaries (that is `helianthus-ebusgateway`)

## Architecture and Dependencies

Dependency direction in the current stack:

`helianthus-ebusgateway` → `helianthus-ebusreg` → `helianthus-ebusgo`

`helianthus-ebusreg` depends on `helianthus-ebusgo` for:

- eBUS frame model (`protocol.Frame`)
- typed bus errors (`errors`)
- data type codecs used by schema/routing (`types`)

Runtime flow (using current terminology):

1. construct `registry.DeviceRegistry` with providers (for example `providers/vaillant.Default()`),
2. scan devices with an initiator address (`protocol.Frame.Source`) and target addresses (`protocol.Frame.Target`) via `registry.Scan(...)`,
3. select planes from discovered entries, wire them into `router.BusEventRouter`, and invoke methods.

## Prerequisites

- Go `1.22+`
- Git
- Optional (for local CI parity): TinyGo `0.40.1` (matches CI setup)
- Optional (for live bus/smoke integration): sibling `helianthus-ebusgateway` + reachable eBUS backend

## Quick Start (Clean Machine)

```bash
git clone https://github.com/d3vi1/helianthus-ebusreg.git
cd helianthus-ebusreg

go mod download
go build ./...
go test ./...
```

CI-parity local checks:

```bash
go vet ./...
go test -race -count=1 ./...
```

## Smoke / Integration Context and Limits

- This repo is a library only (no `main` packages currently).
- End-to-end smoke runs are executed from `helianthus-ebusgateway` (`cmd/smoke`) where real transport and bus wiring exist.
- Docs for integration behavior:
  - https://github.com/d3vi1/helianthus-docs-ebus/blob/main/development/smoke-test.md
  - https://github.com/d3vi1/helianthus-docs-ebus/blob/main/architecture/overview.md

TinyGo note:

- CI includes a TinyGo job, but it skips when no main packages are present.
- On TinyGo targets, JSON schema loaders in `schema` intentionally return `ErrInvalidPayload` (`schema/loader_json_tinygo.go`).

## Package Overview

| Package | Responsibility | Notes |
|---|---|---|
| `registry` | device model, provider matching, scan flow, projection/canonical indexes, portal indexes | central entrypoint for discovery + registration |
| `schema` | field schemas, conditional selector matching, JSON loader helpers | tinygo build disables JSON loaders |
| `router` | invoke path + broadcast fan-out/event stream (`BusEventRouter`) | planes must implement `router.Plane` for invocation |
| `providers/vaillant` | default provider constructors (`System/Heating/DHW/Solar`) | convenience wiring for standard setup |
| `vaillant/*` | concrete provider/plane behavior, schema selectors, templates | includes B509/B524/B516 logic |
| `internal/match` | internal version matching helpers | not public API |

## Provider / Plane Coverage Matrix

Current provider constructors in `providers/vaillant` map to the following capabilities:

| Constructor | Provider `Name()` | Plane name | Method coverage | Router coverage | Projection coverage |
|---|---|---|---|---|---|
| `System()` | `vaillant_system` | `system` | `get_operational_data`, `set_operational_data`, `get_register`, `set_register`, `get_ext_register`, `set_ext_register` | ✅ Implements `router.Plane` and `DecodeBroadcast` (`B5 16`) | ✅ Creates `Service`, `Observability`, `Debug` projections for device IDs `BASV2`, `BAI00`, `VR71`; also emits root projections for present planes (`System`, `Heating`, `DHW`, `Solar`) |
| `Heating()` | `vaillant_heating` | `heating` | `get_status`, `set_target_temp`, `get_parameters`, `get_energy_stats` (only when hardware version is `>=7603`) | ❌ Registry metadata only (does not implement `router.Plane`) | ➖ No direct `ProjectionProvider`; can appear as root `Heating` projection when `vaillant_system` also matches |
| `DHW()` | `vaillant_dhw` | `dhw` | `get_status`, `set_target_temp`, `get_parameters`, `get_energy_stats` (only when hardware version is `>=7603`) | ❌ Registry metadata only (does not implement `router.Plane`) | ➖ No direct `ProjectionProvider`; can appear as root `DHW` projection when `vaillant_system` also matches |
| `Solar()` | `vaillant_solar` | `solar` | `get_status`, `get_solar_yield`, `get_parameters` | ❌ Registry metadata only (does not implement `router.Plane`) | ➖ No direct `ProjectionProvider`; can appear as root `Solar` projection when `vaillant_system` also matches |

`Default()` returns all four providers in this order: `System`, `Heating`, `DHW`, `Solar`.

### System plane invoke parameter quick-reference

For `router.Invoke`, keep using `source` as the map key for the initiator address (this maps to `protocol.Frame.Source`).

| Method | Required params | Optional params | Notes |
|---|---|---|---|
| `get_operational_data` | `source`, `op` | `target` | `target` is only needed when plane address is unknown (for example address `0x00`) |
| `set_operational_data` | `source`, `op` | `data` or `payload`, `target` | payload becomes `[op,...data]` |
| `get_register` | `source`, `addr` | `target` | `addr` accepts `uint16`, hex string (`F600`/`0xF600`), byte pairs |
| `set_register` | `source`, `addr` | `data` or `payload`, `target` | write payload starts with register write opcode |
| `get_ext_register` | `source`, `group`, `instance`, `addr` | `opcode`, `target` | default opcode is local (`0x02`) |
| `set_ext_register` | `source`, `group`, `instance`, `addr` | `opcode`, `data` or `payload`, `target` | write payload includes ext-register header + data |

## End-to-End Cookbook (scan → invoke)

The example below uses one bus implementation for both registry scan and router invoke.

```go
import (
	"context"
	"fmt"

	"github.com/d3vi1/helianthus-ebusgo/protocol"
	vaillantproviders "github.com/d3vi1/helianthus-ebusreg/providers/vaillant"
	"github.com/d3vi1/helianthus-ebusreg/registry"
	"github.com/d3vi1/helianthus-ebusreg/router"
)

func scanThenInvoke(ctx context.Context, bus interface {
	Send(context.Context, protocol.Frame) (*protocol.Frame, error)
}) error {
	initiator := byte(0x10)
	targets := []byte{0x08, 0x09}

	deviceRegistry := registry.NewDeviceRegistry(vaillantproviders.Default())
	entries, err := registry.Scan(ctx, bus, deviceRegistry, initiator, targets)
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		return fmt.Errorf("no devices discovered")
	}

	var plane router.Plane
	for _, candidate := range entries[0].Planes() {
		rp, ok := candidate.(router.Plane)
		if ok && rp.Name() == "system" {
			plane = rp
			break
		}
	}
	if plane == nil {
		return fmt.Errorf("no invokable plane found")
	}

	eventRouter := router.NewBusEventRouter(bus)
	eventRouter.SetPlanes([]router.Plane{plane})

	_, err = eventRouter.Invoke(ctx, plane, "get_operational_data", map[string]any{
		"source": initiator, // initiator address (current API key name)
		"op":     byte(0x00),
	})
	return err
}
```

Cookbook notes:

- Use `registry.DefaultScanTargets()` when you want the full target sweep (`0x01..0xFD`) instead of explicit targets.
- `registry.Scan(...)` skips initiator-capable/invalid targets and retries transient collisions/timeouts.
- Today, only the `system` plane is invokable through `router.BusEventRouter` in this repository.
- If you create a plane without a concrete address (`DeviceInfo.Address == 0`), pass `target` in `router.Invoke` params.

## Extension Checklist (new provider / schema / plane)

- [ ] Add a provider package implementing `registry.PlaneProvider` (`Name`, `Match`, `CreatePlanes`).
- [ ] Define plane method metadata (`registry.Method`) with explicit templates + `schema.SchemaSelector` decoding.
- [ ] Keep initiator/target semantics explicit: initiator is `protocol.Frame.Source`, target is `protocol.Frame.Target`.
- [ ] If method invocation is required, implement `router.Plane` (`BuildRequest`, `DecodeResponse`, `Subscriptions`, `OnBroadcast`).
- [ ] If projection support is required, implement `registry.ProjectionProvider` and maintain `Service` canonical-path invariants.
- [ ] Wire constructors in `providers/vaillant/providers.go` and include in `Default()` when broadly applicable.
- [ ] Add focused tests for matching, method lists, template encoding, schema decode paths, router invoke/broadcast, and projections.
- [ ] Run the validation matrix below before opening a PR.

## Validation Command Matrix

| Scope | Command | Why |
|---|---|---|
| terminology gate (CI parity) | `grep -RInwi --exclude-dir=.git -E 'm[a]ster|s[l]ave' .` | mirrors CI terminology check |
| compile | `go build ./...` | catches compile-time API drift |
| vet | `go vet ./...` | static diagnostics used by CI |
| all tests (CI parity) | `go test -race -count=1 ./...` | same race-enabled suite as CI |
| registry scan/projection focus | `go test ./registry -count=1` | scan logic, registry matching, projection indexes |
| router invoke/broadcast focus | `go test ./router -count=1` | invoke + broadcast fanout behavior |
| schema selectors/loaders focus | `go test ./schema -count=1` | schema parsing/selector behavior |
| Vaillant plane/provider focus | `go test ./vaillant/... -count=1` | provider matching, methods, templates, projection behavior |
| lint (CI parity) | `golangci-lint run` | mirrors `golangci-lint-action` |
| TinyGo parity | `mains=$(go list -f '{{if eq .Name "main"}}{{.ImportPath}}{{end}}' ./... | sed '/^$/d'); if [ -z "$mains" ]; then echo "No main packages found; skipping TinyGo build."; else for pkg in $mains; do tinygo build -target esp32 "$pkg"; done; fi` | mirrors CI TinyGo behavior |

## Troubleshooting

| Symptom / error text | Likely cause | Practical fix |
|---|---|---|
| `scan missing bus` / `scan missing registry` | nil dependency passed to `registry.Scan` | ensure both bus and registry are constructed before scanning |
| `short device info payload` | incomplete `07 04` response payload | verify transport reliability and response framing from bus layer |
| `router.Invoke missing method` | wrong method name for selected plane | inspect `plane.Methods()` names and use exact method ID |
| projection validation errors (`ErrProjectionInvalid*`) | invalid plane/path/node/edge composition | enforce `Service` canonical path invariants and unique node/edge IDs |
| `schema json unavailable` on TinyGo | JSON loader called on TinyGo build | embed/construct schemas in code for TinyGo targets |
| CI TinyGo job prints “No main packages found; skipping” | expected in current repo state | no action needed unless adding runnable binaries |

## Docs, CI, and Issue Tracking

- CI workflow: `.github/workflows/ci.yml`
- Cross-repo architecture/protocol docs:
  - https://github.com/d3vi1/helianthus-docs-ebus/blob/main/architecture/overview.md
  - https://github.com/d3vi1/helianthus-docs-ebus/blob/main/architecture/decisions.md
  - https://github.com/d3vi1/helianthus-docs-ebus/blob/main/protocols/vaillant.md
- Issue tracker: https://github.com/d3vi1/helianthus-ebusreg/issues
- Releases: no published GitHub releases at this time.
