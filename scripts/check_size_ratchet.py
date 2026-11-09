#!/usr/bin/env python3
"""Measure all Go syntax, including inactive tags; baseline debt still blocks release.

Lines are inclusive token spans for functions and physical lines for files.
Literals use lexical ordinals within their enclosing declaration or literal.
Generated Go is measured too: no generator ownership exemption is assumed.
"""
import argparse
import json
import os
from pathlib import Path, PurePosixPath
import shutil
import subprocess
import sys
import tempfile

BUDGETS = {
    "production": {"file_lines": (500, 1000), "function_lines": (80, 150),
                   "functions_per_file": (40, 60)},
    "test": {"file_lines": (800, 1500), "function_lines": (120, 200),
             "functions_per_file": (40, 60)},
}
EXCLUDED = {".git": "repository metadata", ".tmp": "tool scratch",
            "vendor": "dependencies", "node_modules": "dependencies",
            "references": "references"}


def sources(root):
    paths, excluded = [], []
    for directory, dirs, files in os.walk(root, followlinks=False):
        for name in sorted(dirs[:]):
            path = Path(directory) / name
            if path.is_symlink():
                raise ValueError(f"symlink directory cannot be inspected: {path.relative_to(root)}")
            if name in EXCLUDED:
                dirs.remove(name)
                excluded.append((path.relative_to(root).as_posix(), EXCLUDED[name]))
        for name in sorted(files):
            if name.endswith(".go"):
                path = Path(directory) / name
                if path.is_symlink():
                    raise ValueError(f"symlink source cannot be inspected: {path.relative_to(root)}")
                paths.append(path.relative_to(root).as_posix())
    return sorted(paths), sorted(excluded)


def measure(root):
    paths, excluded = sources(root)
    if not paths:
        raise ValueError("no Go files measured")
    scratch = root / ".tmp"
    scratch.mkdir(exist_ok=True)
    with tempfile.TemporaryDirectory(prefix="sizeratchet-", dir=scratch) as work:
        helper = Path(work) / "parser.go"
        shutil.copyfile(Path(__file__).parent / "sizeratchet/parser.go.txt", helper)
        env = dict(os.environ, GOWORK="off", GOPROXY="off", GOSUMDB="off",
                   GOTOOLCHAIN="local", GOFLAGS="-buildvcs=false", TMPDIR=str(scratch))
        result = subprocess.run(["go", "run", str(helper)], cwd=root, env=env,
                                input=json.dumps(paths), text=True, capture_output=True,
                                timeout=150)
        if result.returncode:
            raise ValueError("Go parser failed: " + result.stderr.strip())
        return json.loads(result.stdout), excluded


def limits(row):
    kind = "test" if row["path"].endswith("_test.go") else "production"
    return BUDGETS[kind][row["metric"]]


def key(row):
    return row["path"], row["symbol"], row["metric"]


def unique_object(pairs):
    obj = {}
    for name, value in pairs:
        if name in obj:
            raise ValueError("duplicate JSON key: " + name)
        obj[name] = value
    return obj


def load_baseline(path):
    data = json.loads(path.read_text(), object_pairs_hook=unique_object)
    if not isinstance(data, dict) or set(data) != {"version", "debt"} or type(data["version"]) is not int or data["version"] != 1:
        raise ValueError("invalid baseline schema/version")
    if not isinstance(data["debt"], list):
        raise ValueError("debt must be a list")
    seen = set()
    for row in data["debt"]:
        if not isinstance(row, dict) or set(row) != {"path", "symbol", "metric", "value"}:
            raise ValueError("invalid debt row")
        if any(not isinstance(row[x], str) or not row[x] for x in ("path", "symbol", "metric")):
            raise ValueError("invalid row identity")
        path = PurePosixPath(row["path"])
        if path.is_absolute() or ".." in path.parts or path.as_posix() != row["path"] or not row["path"].endswith(".go") or any(c in row["path"] + row["symbol"] for c in "*?\\\n\r"):
            raise ValueError("invalid exact path/symbol")
        if row["metric"] not in BUDGETS["production"] or type(row["value"]) is not int:
            raise ValueError("invalid metric/value")
        if (row["metric"] == "function_lines") == (row["symbol"] == "<file>"):
            raise ValueError("metric/symbol mismatch")
        if row["value"] <= limits(row)[1] or key(row) in seen:
            raise ValueError("baseline must contain unique hard violations only")
        seen.add(key(row))
    return data["debt"]


def evaluate(rows, baseline):
    allowed = {key(row): row["value"] for row in baseline}
    soft, failures = [], []
    for row in rows:
        low, high = limits(row)
        if row["value"] > low:
            soft.append(row)
        if row["value"] > high and row["value"] > allowed.get(key(row), high):
            failures.append(row)
    return soft, failures


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--root", type=Path, default=Path(__file__).resolve().parents[1])
    args = parser.parse_args()
    try:
        root = args.root.resolve()
        baseline = load_baseline(root / "scripts/size-baseline.json")
        rows, excluded = measure(root)
        for path, reason in excluded:
            print(f"EXCLUDED {path}: {reason}")
        soft, failures = evaluate(rows, baseline)
        for row in soft:
            print("SOFT " + json.dumps(row, sort_keys=True))
        for row in failures:
            print("HARD NEW/GROWN " + json.dumps(row, sort_keys=True))
        print(f"Measured {sum(r['metric'] == 'file_lines' for r in rows)} files; "
              f"{len(soft)} soft violations; {len(failures)} new/grown hard violations. "
              "Baselined debt must still be cleared before release.")
        return int(bool(failures))
    except (ValueError, OSError, subprocess.SubprocessError) as error:
        print(f"size ratchet: {error}", file=sys.stderr)
        return 2


if __name__ == "__main__":
    sys.exit(main())
