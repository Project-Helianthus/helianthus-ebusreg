# AGENTS

## Scope

`helianthus-ebusreg` owns eBUS registry, device identity, typed projection planes, and routing helpers. Keep eBUS-specific registry details here; do not add transport/framing, consumer, MCP, GraphQL, or Home Assistant policy.

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
