# helianthus-ebusreg

`helianthus-ebusreg` is the registry and projection layer for Helianthus eBUS integrations. It turns discovered devices into typed planes, exposes method/schema metadata, and provides routing helpers for invoke and broadcast flows.

## Purpose and Scope

### What belongs in this repository

- Device discovery, scan flow, and registry indexing (`registry/`).
- Method/schema modeling and selector decoding (`schema/`).
- Invoke and broadcast routing (`router/`).
- Default provider wiring and Vaillant plane implementations (`providers/`, `vaillant/`).

### What does not belong in this repository

- Byte-level transport, framing, and low-level bus primitives (use `helianthus-ebusgo`).
- Gateway runtime/API process wiring and live deployment orchestration (use `helianthus-ebusgateway`).

## Status and Maturity

- Active, CI-validated library with stable package-level tests.
- Suitable for contributor onboarding and issue-focused feature work.
- Live adapter smoke runs are executed from `helianthus-ebusgateway`; this repo provides the registry/projection core used there.

## Helianthus Dependency Chain

```text
helianthus-ebusgo  ->  helianthus-ebusreg  ->  helianthus-ebusgateway  ->  integrations/add-ons
  (transport)          (registry/schema)        (runtime/api/smoke)
```

## Method Metadata Contract

Method safety/routing metadata is exposed via `registry.ResolveMethodMetadata(method)`.

- `mutability`: `read_only`, `mutating`, or `unknown`
- `danger`: `safe`, `dangerous`, or `unknown`
- `routable`: `true` or `false`

Backward-compatible defaults for legacy `registry.Method` implementations:

- `mutability`: inferred from `ReadOnly()` (`true` -> `read_only`, `false` -> `mutating`)
- `danger`: inferred from mutability (`read_only` -> `safe`, otherwise -> `dangerous`)
- `routable`: defaults to `true`

Methods can override defaults by implementing optional interfaces in `registry`:
`MethodMutabilityProvider`, `MethodDangerProvider`, and `MethodRoutableProvider`.

## Shared Service Projections

For MCP and GraphQL service-layer reuse, `registry` exposes projection helpers:

- `ProjectRegistryDevices(iter EntryIterator)`
- `ProjectDeviceEntry(entry DeviceEntry)`
- `ProjectPlane(plane Plane)`
- `ProjectMethod(method Method)`

Projected method data includes frame template bytes plus normalized method metadata (`mutability`, `danger`, `routable`) through `ResolveMethodMetadata`.

## Quickstart (copy/paste)

### 1) Clone and baseline checks

```bash
git clone https://github.com/d3vi1/helianthus-ebusreg.git
cd helianthus-ebusreg
go test ./...
go vet ./...
go build ./...
```

### 2) CI-parity test run

```bash
./scripts/ci_local.sh
go test -race -count=1 ./...
```

### 3) Repository notes

```bash
# library-only repo; no main packages are expected
go list -f '{{if eq .Name "main"}}{{.ImportPath}}{{end}}' ./... | sed '/^$/d'
```

## Local Smoke-Test Configuration Examples

`helianthus-ebusreg` is consumed by a runtime service, but these are the core values you typically set during local smoke wiring:

```yaml
scan:
  initiator: 0x10
  targets: [0x08, 0x09, 0x15]
providers:
  - vaillant_system
  - vaillant_heating
  - vaillant_dhw
  - vaillant_solar
invoke_defaults:
  source: 0x10
  target: 0x08
```

Equivalent Go wiring in this repository:

```go
deviceRegistry := registry.NewDeviceRegistry(vaillantproviders.Default())
entries, err := registry.Scan(ctx, bus, deviceRegistry, 0x10, []byte{0x08, 0x09, 0x15})
if err != nil {
	return err
}

eventRouter := router.NewBusEventRouter(bus)
eventRouter.SetPlanes(planesFromEntries(entries))
```

## Focused Validation Commands

| Area | Command |
|---|---|
| terminology gate (CI parity) | `if grep -RInwi --exclude-dir=.git -E 'm[a]ster|s[l]ave' .; then echo "Found legacy terminology."; exit 1; fi` |
| compile | `go build ./...` |
| vet | `go vet ./...` |
| all tests (CI parity) | `go test -race -count=1 ./...` |
| registry scan/indexes | `go test ./registry -count=1` |
| router invoke/broadcast | `go test ./router -count=1` |
| schema selectors/loaders | `go test ./schema -count=1` |
| Vaillant providers/planes | `go test ./vaillant/... -count=1` |
| lint (if installed locally) | `golangci-lint run` |
| TinyGo CI parity | `mains=$(go list -f '{{if eq .Name "main"}}{{.ImportPath}}{{end}}' ./... | sed '/^$/d'); if [ -z "$mains" ]; then echo "No main packages found; skipping TinyGo build."; else for pkg in $mains; do tinygo build -target esp32 "$pkg"; done; fi` |

## Link Map

### Core repos

- `helianthus-ebusgo`: https://github.com/d3vi1/helianthus-ebusgo
- `helianthus-ebusgateway`: https://github.com/d3vi1/helianthus-ebusgateway

### Architecture and smoke docs

- Architecture overview: https://github.com/d3vi1/helianthus-docs-ebus/blob/main/architecture/overview.md
- Architecture decisions: https://github.com/d3vi1/helianthus-docs-ebus/blob/main/architecture/decisions.md
- Smoke test flow: https://github.com/d3vi1/helianthus-docs-ebus/blob/main/development/smoke-test.md
- Protocol overview: https://github.com/d3vi1/helianthus-docs-ebus/blob/main/protocols/ebus-overview.md

### Issues and workflow conventions

- Issue tracker: https://github.com/d3vi1/helianthus-ebusreg/issues
- Keep one issue-focused branch (example: `issue-70-readme-refresh`).
- Keep PR scope aligned to issue acceptance criteria.
- Include closing keyword in PR body (example: `Fixes #70`).
