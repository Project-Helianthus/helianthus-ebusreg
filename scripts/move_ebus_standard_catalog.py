#!/usr/bin/env python3
"""Reproduce the move of catalog/ebus_standard to catalog/ebusstandard.

The Go import path and package name change; the ebus_standard service
namespace, catalog data, fixtures and error texts are left byte-identical.
Run once against a tree that still holds catalog/ebus_standard.
"""

import argparse
from pathlib import Path

OLD_DIR = "catalog/ebus_standard"
NEW_DIR = "catalog/ebusstandard"
OLD_PACKAGE = "ebus_standard_catalog"
NEW_PACKAGE = "ebusstandard"
GOLDEN = NEW_DIR + "/testdata/provider_abi.golden.txt"
PATH_KEYED = (
    "naming-inventory.json",
    "scripts/behavior-name-baseline.json",
    "scripts/go-lint-baseline.json",
    "scripts/size-baseline.json",
)
CHANGELOG_ENTRY = (
    "- Move `catalog/ebus_standard` (package `ebus_standard_catalog`) to\n"
    "  `catalog/ebusstandard` (package `ebusstandard`), together with its\n"
    "  `internal/safetypolicy` package. Consumers must update imports; the old\n"
    "  path has no forwarding package. The `ebus_standard` service namespace,\n"
    "  method identifiers, catalog data and error texts are unchanged. Gateway\n"
    "  adoption and release bill of materials reconciliation remain required.\n"
)


def rewrite(path, *pairs, line_prefix=None):
    before = path.read_bytes()
    after = before
    for old, new in pairs:
        after = after.replace(old.encode(), new.encode())
    if line_prefix is not None:
        old, new = (part.encode() for part in line_prefix)
        lines = after.split(b"\n")
        after = b"\n".join(new + line[len(old):] if line.startswith(old) else line
                           for line in lines)
    if after != before:
        path.write_bytes(after)


def go_sources(root):
    for path in sorted(root.rglob("*.go")):
        if any(part.startswith(".") for part in path.relative_to(root).parts):
            continue
        yield path


def move(root):
    source, target = root / OLD_DIR, root / NEW_DIR
    if target.exists():
        raise SystemExit(f"target already exists: {NEW_DIR}")
    if not source.is_dir():
        raise SystemExit(f"source directory missing: {OLD_DIR}")
    source.rename(target)
    for path in go_sources(root):
        rewrite(path,
                (OLD_DIR, NEW_DIR),
                ("package " + OLD_PACKAGE, "package " + NEW_PACKAGE),
                ("Package " + OLD_PACKAGE, "Package " + NEW_PACKAGE))
    rewrite(root / GOLDEN, line_prefix=(OLD_PACKAGE + " ", NEW_PACKAGE + " "))
    for name in PATH_KEYED:
        rewrite(root / name, ('"' + OLD_DIR + "/", '"' + NEW_DIR + "/"))
    changelog = root / "CHANGELOG.md"
    text = changelog.read_bytes()
    heading = b"## Unreleased\n\n"
    if heading not in text:
        raise SystemExit("CHANGELOG.md has no Unreleased section")
    changelog.write_bytes(text.replace(heading, heading + CHANGELOG_ENTRY.encode() + b"\n", 1))


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("root", nargs="?", type=Path,
                        default=Path(__file__).resolve().parent.parent)
    move(parser.parse_args().root.resolve())
