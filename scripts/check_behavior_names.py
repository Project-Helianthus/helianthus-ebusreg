#!/usr/bin/env python3
"""Ratchet exact declaration/file debt; never interpret payloads or prose.

The Go parser helper includes every platform/build-tagged source file. Baseline
entries identify one path, symbol and category; moving or renaming debt therefore
requires review. Removing debt is allowed without regenerating the baseline.
Protocol/external exceptions are exact vocabulary tokens with a nonempty reason,
not file patterns. Exceptions must cite a public native protocol or stable conformance contract.
"""
import argparse
import json
import os
from pathlib import Path
import re
import subprocess
import sys
import tempfile

ROOT = Path(__file__).resolve().parents[1]
BASELINE = ROOT / "scripts/behavior-name-baseline.json"
CATEGORIES = {"filename", "identifier", "test"}
TOKENS = re.compile(
    r"(?i:batch|wave|issue|ticket|phase|step|milestone|sprint)[_-]*[0-9]+"
    r"|(?<![A-Z])[A-Z]+[0-9]+(?:[A-F](?![a-z]))?"
    r"|(?<![a-z])[a-z]{1,2}[0-9]+(?:[a-f](?![a-z]))?"
)


def load_policy(path):
    policy = json.loads(Path(path).read_text())
    if (not isinstance(policy, dict) or set(policy) != {"version", "exceptions", "baseline"}
            or type(policy["version"]) is not int or policy["version"] != 1
            or not isinstance(policy["exceptions"], list)
            or not isinstance(policy["baseline"], list)):
        raise ValueError("unsupported policy schema")
    vocabulary = set()
    for entry in policy["exceptions"]:
        if not isinstance(entry, dict) or set(entry) != {"token", "reason"}:
            raise ValueError("exceptions require an exact token and reason")
        token, reason = entry["token"], entry["reason"]
        if (not isinstance(token, str) or not TOKENS.fullmatch(token)
                or token != token.upper() or token in vocabulary
                or not isinstance(reason, str) or not reason.strip()):
            raise ValueError("invalid exact vocabulary exception")
        vocabulary.add(token)
    baseline = set()
    for entry in policy["baseline"]:
        if not isinstance(entry, dict) or set(entry) != {"path", "symbol", "category"}:
            raise ValueError("baseline requires path, symbol and category")
        path, symbol, category = (entry[k] for k in ("path", "symbol", "category"))
        if (not isinstance(path, str) or not path or path.startswith("/")
                or ".." in path.split("/") or "\\" in path
                or not isinstance(symbol, str) or not symbol
                or not isinstance(category, str) or category not in CATEGORIES
                or any(c in path + symbol for c in "*?[]")):
            raise ValueError("invalid exact baseline entry")
        key = (path, symbol, category)
        if key in baseline:
            raise ValueError("duplicate baseline entry")
        baseline.add(key)
    return vocabulary, baseline


def inventory(root):
    env = dict(os.environ, GOWORK="off", GOPROXY="off", GOSUMDB="off",
               GOTOOLCHAIN="local", GOFLAGS="-buildvcs=false")
    # The parser source is kept as a text template so the product Go-file set
    # is unchanged; it is materialized in a scratch directory to run.
    template = (ROOT / "scripts/behaviornames/parser.go.txt").read_text()
    with tempfile.TemporaryDirectory() as scratch:
        source = Path(scratch) / "main.go"
        source.write_text(template)
        completed = subprocess.run(
            ["go", "run", str(source), str(root)],
            env=env, text=True, capture_output=True, timeout=120, check=True)
    entries = json.loads(completed.stdout)
    for base, dirs, files in os.walk(root, followlinks=False):
        dirs[:] = sorted(d for d in dirs if d not in {".git", ".tmp"})
        for directory in dirs:
            if (Path(base) / directory).is_symlink():
                raise ValueError("symlink directory input")
        for name in sorted(files):
            path = Path(base) / name
            rel = path.relative_to(root).as_posix()
            if rel in {".REPORT-492.md", ".pr-body-492.md"}:
                continue
            if path.is_symlink():
                raise ValueError("symlink input: " + rel)
            entries.append({"path": rel, "symbol": name, "category": "filename"})
    return entries


TEST_PREFIX = re.compile(r"^(?:Test|Benchmark|Fuzz|Example)")


def scan_text(entry):
    # The Go test prefix is not part of the name, so a lowercase code directly
    # after it (Testxy73Read) is a token start just like TestXY73Read.
    symbol = entry["symbol"]
    if entry["category"] == "test":
        symbol = TEST_PREFIX.sub("", symbol, count=1)
    return symbol


def violations(entries, vocabulary):
    return sorted({(entry["path"], entry["symbol"], entry["category"])
                   for entry in entries
                   if any(m.group().upper() not in vocabulary
                          for m in TOKENS.finditer(scan_text(entry)))})


def check(root, policy):
    vocabulary, baseline = load_policy(policy)
    return [entry for entry in violations(inventory(root), vocabulary)
            if entry not in baseline]


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--check", action="store_true", required=True)
    parser.add_argument("--root", type=Path, default=ROOT)
    parser.add_argument("--baseline", type=Path, default=BASELINE)
    args = parser.parse_args()
    try:
        new = check(args.root.resolve(), args.baseline)
    except (ValueError, OSError, subprocess.SubprocessError) as exc:
        print("behavior-name scan failed: " + str(exc), file=sys.stderr)
        return 2
    for path, symbol, category in new:
        print(f"{path}: {category} {symbol}")
    print(f"behavior-name guard: {len(new)} new entries")
    return 1 if new else 0


if __name__ == "__main__":
    sys.exit(main())
