#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$repo_root"

echo "==> terminology gate"
# Phase C M-C6: AddressByRole(SlotRole) test documents the eBUS-spec
# SlotRoleMaster/SlotRoleSlave enum values. Mirrors the exclusion in
# .github/workflows/ci.yml so local CI matches GH Actions semantics.
if git grep -nIwiE 'm[a]ster|s[l]ave' -- \
    ':!vendor/' \
    ':!registry/address_by_role_test.go'; then
  echo "Found legacy terminology."
  exit 1
fi

echo "==> gofmt"
unformatted="$(git ls-files '*.go' | xargs -n 50 gofmt -l || true)"
if [ -n "${unformatted}" ]; then
  echo "gofmt required for:"
  echo "${unformatted}"
  exit 1
fi

echo "==> go vet"
go vet ./...

echo "==> go build"
go build ./...

echo "==> go test (race)"
go test -race -count=1 ./...

if command -v golangci-lint >/dev/null 2>&1; then
  echo "==> golangci-lint"
  golangci-lint run ./...
else
  echo "==> golangci-lint not found; skipping"
fi

if command -v tinygo >/dev/null 2>&1; then
  echo "==> tinygo build (main packages)"
  toolchain="${TINYGO_GOTOOLCHAIN:-go1.25.0}"
  echo "Using GOTOOLCHAIN=${toolchain} for tinygo"
  mains="$(go list -f '{{if eq .Name "main"}}{{.ImportPath}}{{end}}' ./... | sed '/^$/d')"
  if [ -z "${mains}" ]; then
    echo "No main packages found; skipping TinyGo build."
    exit 0
  fi
  for pkg in ${mains}; do
    echo "tinygo build -target esp32 ${pkg}"
    GOTOOLCHAIN="${toolchain}" tinygo build -target esp32 "${pkg}"
  done
else
  echo "==> tinygo not found; skipping"
fi

