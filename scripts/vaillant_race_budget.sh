#!/usr/bin/env bash
# Run each catalog/vaillant test alone under the race detector.
# A test that does not finish within 90 seconds is named and the script
# exits non-zero. The 20 minute suite timeout is not changed. No test is
# skipped and -short is not passed.
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$repo_root"

list_file="$(mktemp)"
log_file="$(mktemp)"
trap 'rm -f "${list_file}" "${log_file}"' EXIT

go test -list '^Test' ./catalog/vaillant >"${list_file}"

found=0
failed=0
while IFS= read -r name; do
  case "${name}" in
    Test*) ;;
    *) continue ;;
  esac
  case "${name}" in
    *[!A-Za-z0-9_]*)
      echo "vaillant race budget: refused a test name that is not an identifier: ${name}" >&2
      exit 1
      ;;
  esac
  found=1
  echo "==> ${name}"
  set +e
  go test -race -count=1 -timeout 90s -run "^(${name})$" ./catalog/vaillant >"${log_file}" 2>&1
  status=$?
  set -e
  cat "${log_file}"
  if [[ "${status}" -ne 0 ]]; then
    if grep -E -q 'test timed out|Test killed: ran too long' "${log_file}"; then
      echo "vaillant race budget timeout: ${name}"
      exit 1
    fi
    echo "vaillant race budget failed: ${name}"
    failed=1
  fi
done <"${list_file}"

if [[ "${found}" -eq 0 ]]; then
  echo "vaillant race budget: no tests listed in ./catalog/vaillant" >&2
  exit 1
fi

if [[ "${failed}" -ne 0 ]]; then
  echo "vaillant race budget: a listed test failed inside 90 seconds"
  exit 1
fi

echo "vaillant race budget: every listed test finished within 90 seconds"
exit 0
