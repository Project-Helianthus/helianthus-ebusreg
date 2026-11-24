# AGENTS

## Scope

`helianthus-ebusreg` owns eBUS registry, device identity, typed projection planes, and routing helpers. Keep eBUS-specific registry details here; transport/framing is not added here, and consumer, MCP, GraphQL, or Home Assistant policy stays out as well.

`helianthus-ebusreg` is the eBUS driver. It owns eBUS device knowledge and its mapping to SemReg in both directions. Hosts import this driver and SemReg; this repository never imports a host.

Among the Helianthus modules, the packages of this repository import only `helianthus-ebusgo` (the eBUS transport), the SemReg driver contract `github.com/Project-Helianthus/helianthus-semreg/semreg/v1/driver`, the row evaluator `github.com/Project-Helianthus/helianthus-semreg/semreg/v1/roweval`, the types exposed at `github.com/Project-Helianthus/helianthus-semreg/semreg/v1`, the exact public packages `github.com/Project-Helianthus/helianthus-semreg/semreg/v1/catalog`, `github.com/Project-Helianthus/helianthus-semreg/semreg/v1/operation` and `github.com/Project-Helianthus/helianthus-semreg/semreg/v1/packs/thermal`, and this module. The standard library and `gopkg.in/yaml.v3`, which `go.mod` already requires, are not restricted. Importing the gateway or any other `helianthus-semreg` package is not allowed. The repository opens no socket, serves no API, and holds no host code. `scripts/tests/driver_imports_test.py` checks this import direction.

This repository is not the universal semantic layer. The planned `helianthus-semreg` repository will own universal, cross-protocol semantics when that work is active.

Public references:

- https://github.com/Project-Helianthus/helianthus-docs-ebus/blob/main/architecture/overview.md
- https://github.com/Project-Helianthus/helianthus-docs-ebus/blob/main/architecture/decisions.md
- https://github.com/Project-Helianthus/helianthus-docs-ebus/blob/main/protocols/ebus-services/ebus-overview.md

## Working rules

- One focused issue and one PR at a time. Branches use `issue/<number>-<slug>` and start from a clean `origin/main` worktree.
- Preserve field-level merge behavior: on partial failures, retain last-known-good values for unaffected fields; never replace useful state wholesale.
- Run `./scripts/ci_local.sh` before pushing.
- Review the exact PR HEAD in a fresh context. Fix valid P0-P2 findings and re-review the new HEAD; P3-P4 are triaged without blocking.
- Use squash merge only after CI and fresh exact-HEAD review are clear. Verify remote `main` and stop at the requested boundary.
- Stop for explicit action-time confirmation before credential handling, real installation writes, live-device mutation, or destructive/irreversible operations.

## Documentation

When a change establishes or changes public eBUS registry, identity, or projection behavior, update the public documentation in `helianthus-docs-ebus` in the same delivery cycle. Documentation-only instruction changes do not establish a protocol or semantic claim.
