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

At runtime, the gateway usually:

1. constructs `registry.DeviceRegistry` with providers (for example `providers/vaillant.Default()`),
2. discovers devices via `registry.Scan(...)`,
3. builds router planes from matched providers and invokes via `router.BusEventRouter`.

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
| `router` | invoke path + broadcast fan-out/event stream (`BusEventRouter`) | planes must implement `router.Plane` |
| `providers/vaillant` | default provider constructors (`System/Heating/DHW/Solar`) | convenience wiring for standard setup |
| `vaillant/*` | concrete plane/provider behavior and templates | includes B509/B524/B516 behavior |
| `internal/match` | internal matching helpers | not public API |

## Common Workflows

### 1) Full local verification (same core checks as CI)

```bash
go vet ./...
go build ./...
go test -race -count=1 ./...
```

### 2) Iterate on registry/projection logic

```bash
go test ./registry -count=1
```

### 3) Iterate on router invocation/broadcast behavior

```bash
go test ./router -count=1
```

### 4) Iterate on schema loaders/selectors and Vaillant providers

```bash
go test ./schema -count=1
go test ./vaillant/... -count=1
```

## Troubleshooting

| Symptom / error text | Likely cause | Practical fix |
|---|---|---|
| `scan missing bus` / `scan missing registry` | nil dependency passed to `registry.Scan` | ensure both bus and registry are constructed before scanning |
| `short device info payload` | incomplete 07/04 response payload | verify transport reliability and response framing from bus layer |
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
