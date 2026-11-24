#!/usr/bin/env python3
"""Destination limits, names, lint, and manifest closure for one candidate."""

import argparse
import json
import subprocess
import sys
from pathlib import Path


SCRIPTS = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(SCRIPTS))

import check_size_ratchet as size


DESTINATION_REPOSITORY = "ebusreg"


def manifest_tool(manifest):
    return Path(manifest).resolve().parent.parent / "tools" / "manifest" / "main.py"


def head_sha(donor):
    completed = subprocess.run(
        ["git", "-C", str(donor), "rev-parse", "HEAD"],
        text=True,
        capture_output=True,
    )
    if completed.returncode != 0:
        raise RuntimeError((completed.stderr or completed.stdout).strip() or "donor rev-parse failed")
    return completed.stdout.strip()


def destination_entries(manifest):
    entries = []
    for entry in manifest.get("destinations") or []:
        if not isinstance(entry, dict):
            continue
        entries.append(entry)
    return entries


def destination_paths(manifest):
    paths = []
    for entry in destination_entries(manifest):
        if entry.get("repository") != DESTINATION_REPOSITORY:
            continue
        path = entry.get("path")
        if isinstance(path, str) and path and path not in paths:
            paths.append(path)
    return paths


def closure_lines(manifest_path, donor):
    tool = manifest_tool(manifest_path)
    if not tool.is_file():
        return ["manifest tool is unavailable"]
    completed = subprocess.run(
        [sys.executable, str(tool), "check", "--manifest", str(manifest_path), "--donor", str(donor)],
        cwd=str(tool.resolve().parents[3]),
        text=True,
        capture_output=True,
    )
    lines = [line for line in completed.stdout.splitlines() if line and line != "check: ok"]
    if completed.returncode and not lines:
        lines.append((completed.stderr or "").strip() or "check failed")
    return lines


def allowance(row, baseline):
    allowed = {size.key(item): item["value"] for item in baseline}
    _soft, hard = size.limits(row)
    return allowed.get(size.key(row), hard)


def behavior_lines(candidate, destinations):
    script = SCRIPTS / "check_behavior_names.py"
    baseline = candidate / "scripts" / "behavior-name-baseline.json"
    completed = subprocess.run(
        [
            sys.executable,
            str(script),
            "--check",
            "--root",
            str(candidate),
            "--baseline",
            str(baseline),
        ],
        text=True,
        capture_output=True,
    )
    lines = []
    for line in completed.stdout.splitlines():
        if not line or line.startswith("behavior-name guard:"):
            continue
        path, separator, _rest = line.partition(": ")
        if separator and path in destinations:
            lines.append("rename or fix required: " + line)
    if completed.returncode and not lines:
        detail = (completed.stderr or completed.stdout).strip() or "behavior-name scan failed"
        lines.append("rename or fix required: " + detail)
    return lines


def lint_script(candidate):
    local = candidate / "scripts" / "check_go_lint_baseline.py"
    if local.is_file():
        return local.resolve()
    return (SCRIPTS / "check_go_lint_baseline.py").resolve()


def outside_module(candidate, destinations):
    module = lint_script(candidate).parents[1]
    outside = []
    for path in destinations:
        try:
            (candidate / path).resolve().relative_to(module)
        except ValueError:
            outside.append(path)
    return outside


def lint_lines(candidate, destinations):
    script = lint_script(candidate)
    completed = subprocess.run(
        [sys.executable, str(script), "--check"],
        cwd=str(script.parents[1]),
        text=True,
        capture_output=True,
    )
    lines = []
    for line in (completed.stdout + "\n" + completed.stderr).splitlines():
        if not line:
            continue
        if "finite exemptions" in line:
            continue
        if any(path in line for path in destinations):
            lines.append("rename or fix required: " + line)
            continue
        if completed.returncode and line.startswith("lint baseline failed:"):
            lines.append("rename or fix required: " + line)
    return lines


def split_line(row, baseline):
    return "split required: {path} {symbol} {metric} {value} > {limit}".format(
        path=row["path"],
        symbol=row["symbol"],
        metric=row["metric"],
        value=row["value"],
        limit=allowance(row, baseline),
    )


def run(manifest_path, donor, candidate):
    manifest_path = Path(manifest_path)
    candidate = Path(candidate)
    manifest = json.loads(manifest_path.read_text(encoding="utf-8"))
    paths = destination_paths(manifest)
    lines = []
    if donor is None:
        lines.append("donor checkout is required")
        return lines
    donor = Path(donor)
    expected = manifest.get("donor", {}).get("sha")
    actual = head_sha(donor)
    lines.extend(closure_lines(manifest_path, donor))
    if actual != expected:
        if "donor-sha" not in lines:
            lines.append("donor-sha")
        return lines
    for entry in destination_entries(manifest):
        repository = entry.get("repository")
        path = entry.get("path")
        if repository != DESTINATION_REPOSITORY and path:
            lines.append("stop: destination file is not in this module: " + str(path))
            return lines
    missing = [path for path in paths if not (candidate / path).is_file()]
    for path in missing:
        lines.append("absent destination: " + path)
    outside = outside_module(candidate, paths)
    if outside:
        for path in outside:
            lines.append("stop: destination file is not in this module: " + path)
        return lines
    baseline = size.load_baseline(candidate / "scripts" / "size-baseline.json")
    try:
        rows, _excluded = size.measure(candidate)
    except ValueError as error:
        lines.append(str(error))
        return lines
    soft, failures = size.evaluate(rows, baseline)
    for row in soft:
        lines.append("SOFT " + json.dumps(row, sort_keys=True))
    for row in failures:
        lines.append(split_line(row, baseline))
    present = set(paths) - set(missing)
    lines.extend(behavior_lines(candidate, present))
    lines.extend(lint_lines(candidate, present))
    return lines


def main(argv=None):
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--manifest", required=True, type=Path)
    parser.add_argument("--donor", type=Path)
    args = parser.parse_args(argv)
    try:
        lines = run(args.manifest, args.donor, Path.cwd())
    except (OSError, RuntimeError, ValueError, json.JSONDecodeError, subprocess.SubprocessError) as error:
        print("preflight: " + str(error))
        return 1
    findings = 0
    for line in lines:
        print(line)
        if not line.startswith("SOFT "):
            findings += 1
    return 1 if findings else 0


if __name__ == "__main__":
    sys.exit(main())
