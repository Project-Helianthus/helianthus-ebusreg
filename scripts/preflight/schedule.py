#!/usr/bin/env python3
"""Stack position and additive job selection for one pull request."""

import argparse
import json
import os
import subprocess
import sys
import urllib.error
import urllib.request
from pathlib import Path


DESTINATION_REPOSITORY = "ebusreg"
DOCS_SEMANTIC_URL = "https://github.com/Project-Helianthus/helianthus-docs-semantic.git"
FULL_JOBS = (
    "terminology",
    "build",
    "test",
    "tinygo",
    "lint",
    "size-ratchet",
    "behavior-names",
    "readability",
)
FAST_MODES = ("fast-additive", "fast-compare")
ADDITIVE_JOB_MODES = ("fast-additive", "fast-compare", "check-failed")

# Public names from the manifest tool's acquisition recipe, including the
# gateway rename to helianthus-gateway.
PUBLIC_REPOSITORIES = {
    "gateway": "helianthus-gateway",
    "helianthus-ebusgateway": "helianthus-gateway",
    "semreg": "helianthus-semreg",
    "ebusreg": "helianthus-ebusreg",
    "ebusgo": "helianthus-ebusgo",
    "modbusreg": "helianthus-modbusreg",
    "modbus": "helianthus-modbus",
    "canbusreg": "helianthus-canbusreg",
    "canbus": "helianthus-canbus",
    "eebusreg": "helianthus-eebusreg",
    "eebus-go": "helianthus-eebus-go",
    "matterreg": "helianthus-matterreg",
    "helianthus-matter-binding": "helianthus-matterreg",
    "ha-integration": "helianthus-ha-integration",
    "vrc-explorer": "helianthus-vrc-explorer",
    "helianthus-vrc-explorer": "helianthus-vrc-explorer",
    "docs-semantic": "helianthus-docs-semantic",
    "addon": "helianthus-ha-addon",
    "helianthus-ha-addon": "helianthus-ha-addon",
}

RETAINED_COMMANDS = {
    "test": "go test -race -count=1 -timeout 20m ./...",
    "size-ratchet": "python3 scripts/check_size_ratchet.py",
    "behavior-names": "python3 scripts/check_behavior_names.py --check",
    "readability": "python3 scripts/check_go_lint_baseline.py --check",
    "driver-imports": "python3 -m unittest discover -s scripts/tests -p driver_imports_test.py",
    "vet": "go vet ./...",
    "build": "go build ./...",
}


def public_repository(name):
    if name not in PUBLIC_REPOSITORIES:
        raise KeyError(name)
    return PUBLIC_REPOSITORIES[name]


def repository_url(name):
    return "https://github.com/Project-Helianthus/" + public_repository(name) + ".git"


def stack_position(pulls, current):
    """Parent head matches child base by repository id and ref.

    A member with a child is lower. A member with a parent and no child is
    the tip. A member with neither is unstacked. Missing identity fields or
    an unavailable pull list are unresolved.
    """
    if pulls is None or current is None:
        return "unresolved"
    fields = ("number", "base_repo_id", "base_ref", "head_repo_id", "head_ref")
    if any(current.get(field) in (None, "") for field in fields):
        return "unresolved"
    parent = False
    child = False
    for pull in pulls:
        if pull is current:
            continue
        if any(pull.get(field) in (None, "") for field in fields):
            return "unresolved"
        if pull.get("number") == current.get("number"):
            continue
        if pull["head_repo_id"] == current["base_repo_id"] and pull["head_ref"] == current["base_ref"]:
            parent = True
        if pull["base_repo_id"] == current["head_repo_id"] and pull["base_ref"] == current["head_ref"]:
            child = True
    if child:
        return "lower"
    if parent:
        return "tip"
    return "unstacked"


def sealed_paths(manifest, repository):
    for entry in manifest.get("paths") or []:
        if entry.get("repository") == repository:
            return list(entry.get("paths") or [])
    return None


def select_manifest(manifests, changed_paths, repository=DESTINATION_REPOSITORY):
    """The manifest whose sealed destination list covers every changed path.

    Donor sealed paths are not destination paths. An equal-tightness tie is
    not a match. No match selects full CI.
    """
    if not changed_paths:
        return None
    covered = []
    for manifest in manifests:
        sealed = sealed_paths(manifest, repository)
        if sealed is None:
            continue
        if all(path in sealed for path in changed_paths):
            covered.append((len(sealed), str(manifest.get("row") or ""), manifest))
    if not covered:
        return None
    covered.sort(key=lambda item: (item[0], item[1]))
    tightest = covered[0][0]
    same = [item for item in covered if item[0] == tightest]
    if len(same) > 1:
        return None
    return covered[0][2]


def changed_statuses(changed):
    return [item["status"] for item in changed]


def is_additive(changed, sealed):
    if not changed:
        return False
    paths = [item["path"] for item in changed]
    if any(path not in sealed for path in paths):
        return False
    return all(item["status"] == "added" for item in changed)


def classify_batch(position, manifest, changed, check_ok, repository=DESTINATION_REPOSITORY):
    selected = None
    if manifest is not None:
        sealed = sealed_paths(manifest, repository)
        paths = [item["path"] for item in changed]
        if sealed is not None and paths and all(path in sealed for path in paths):
            selected = manifest
    if position not in {"lower", "tip", "unstacked", "unresolved"}:
        position = "unresolved"
    if selected is None:
        return {"mode": "full", "position": position, "additive": False}
    if not check_ok:
        return {"mode": "check-failed", "position": position, "additive": False}
    if position == "unresolved":
        return {"mode": "full", "position": position, "additive": False}
    additive = is_additive(changed, sealed_paths(selected, repository))
    if additive and position in {"unstacked", "lower"}:
        mode = "fast-additive"
    elif (not additive) and position == "lower":
        mode = "fast-compare"
    else:
        mode = "full"
    return {"mode": mode, "position": position, "additive": additive}


def job_plan(mode):
    skip_full = mode in FAST_MODES
    jobs = {name: ("skip" if skip_full else "run") for name in FULL_JOBS}
    jobs["driver-imports"] = "run"
    jobs["additive-rows"] = "run" if mode in ADDITIVE_JOB_MODES else "skip"
    if mode == "fast-additive":
        steps = ["package-tests", "check", "compare", "preflight", "size", "behavior-names", "lint"]
    elif mode == "fast-compare":
        steps = ["package-tests", "check", "compare"]
    elif mode == "check-failed":
        steps = ["check", "compare"]
    else:
        steps = []
    return {"mode": mode, "jobs": jobs, "steps": steps, "commands": dict(RETAINED_COMMANDS)}


def packages(changed):
    found = []
    for item in changed:
        path = item["path"] if isinstance(item, dict) else item
        if not str(path).endswith(".go"):
            continue
        parent = Path(path).parent.as_posix()
        package = "." if parent in ("", ".") else "./" + parent
        if package not in found:
            found.append(package)
    return found


def parse_name_status(text):
    changed = []
    for line in text.splitlines():
        if not line.strip():
            continue
        parts = line.split("\t")
        status = parts[0]
        if status.startswith("R") or status.startswith("C"):
            changed.append({"path": parts[-1], "status": "renamed"})
        elif status == "A":
            changed.append({"path": parts[1], "status": "added"})
        elif status == "D":
            changed.append({"path": parts[1], "status": "deleted"})
        elif status == "M" or status == "T":
            changed.append({"path": parts[1], "status": "modified"})
        else:
            changed.append({"path": parts[-1], "status": status})
    return changed


def pull_record(item):
    base = item.get("base") or {}
    head = item.get("head") or {}
    base_repo = base.get("repo") or {}
    head_repo = head.get("repo") or {}
    return {
        "number": item.get("number"),
        "base_repo_id": base_repo.get("id"),
        "base_ref": base.get("ref"),
        "head_repo_id": head_repo.get("id"),
        "head_ref": head.get("ref"),
    }


def next_link(header):
    if not header:
        return None
    for part in header.split(","):
        if 'rel="next"' not in part:
            continue
        url = part.split(";", 1)[0].strip()
        if url.startswith("<") and url.endswith(">"):
            return url[1:-1]
    return None


def list_pulls(repository):
    url = "https://api.github.com/repos/" + repository + "/pulls?state=open&per_page=100"
    found = []
    while url:
        request = urllib.request.Request(
            url,
            headers={
                "Accept": "application/vnd.github+json",
                "User-Agent": "ebusreg-preflight",
            },
        )
        with urllib.request.urlopen(request, timeout=60) as response:
            body = json.loads(response.read().decode())
            link = response.headers.get("Link")
        if not isinstance(body, list):
            raise RuntimeError("pull list is not an array")
        found.extend(pull_record(item) for item in body)
        url = next_link(link)
    return found


def git(args, cwd=None):
    completed = subprocess.run(args, cwd=cwd, text=True, capture_output=True)
    if completed.returncode != 0:
        detail = (completed.stderr or completed.stdout).strip()
        raise RuntimeError(detail or " ".join(args))
    return completed.stdout


def resolve_ref(url, ref):
    text = git(["git", "ls-remote", url, ref])
    line = text.splitlines()[0] if text.strip() else ""
    sha = line.split()[0] if line else ""
    if len(sha) != 40 or any(character not in "0123456789abcdef" for character in sha):
        raise RuntimeError("ls-remote did not return one commit")
    return sha


def checkout_commit(url, sha, dest):
    dest = Path(dest)
    if not dest.exists():
        git(["git", "clone", url, str(dest)])
    git(["git", "-C", str(dest), "fetch", "--depth", "1", "origin", sha])
    git(["git", "-C", str(dest), "checkout", "--detach", sha])


def commit_exists(repo, revision):
    if not revision:
        return False
    completed = subprocess.run(
        ["git", "-C", str(repo), "cat-file", "-e", revision + "^{commit}"],
        text=True,
        capture_output=True,
    )
    return completed.returncode == 0


def fetch_revision(repo, repository_name, revision):
    if commit_exists(repo, revision):
        return
    url = repository_url(repository_name)
    git(["git", "-C", str(repo), "fetch", "--depth", "1", url, revision])


def load_manifests(checkout):
    directory = Path(checkout) / "acceptance" / "manifests"
    manifests = []
    if not directory.is_dir():
        return manifests
    for path in sorted(directory.glob("*.json")):
        try:
            manifest = json.loads(path.read_text(encoding="utf-8"))
        except json.JSONDecodeError:
            continue
        if isinstance(manifest, dict) and "paths" in manifest:
            manifest["_path"] = path
            manifests.append(manifest)
    return manifests


def changed_files(root, base_sha, head_sha):
    text = git(["git", "-C", str(root), "diff", "--name-status", base_sha, head_sha])
    return parse_name_status(text)


def reference_fetches(manifest):
    donor_fetches = []
    tool_fetches = []
    for pattern in manifest.get("patterns") or []:
        exemplar = pattern.get("exemplar") or {}
        if exemplar.get("repository") and exemplar.get("revision"):
            donor_fetches.append((exemplar["repository"], exemplar["revision"]))
        candidate = pattern.get("candidate") or {}
        reference = candidate.get("reference") or {}
        if candidate.get("form") == "reference" and reference.get("repository"):
            tool_fetches.append((reference["repository"], reference["revision"]))
    for expectation in manifest.get("expectations") or []:
        for reference in expectation.get("references") or []:
            if reference.get("repository") and reference.get("revision"):
                tool_fetches.append((reference["repository"], reference["revision"]))
    for difference in manifest.get("expected_differences") or []:
        reference = difference.get("reference") or {}
        if reference.get("repository") and reference.get("revision"):
            tool_fetches.append((reference["repository"], reference["revision"]))
    return donor_fetches, tool_fetches


def prepare_inputs(source_sha, manifest_rel, donor_repository, donor_sha, tool, donor):
    checkout_commit(DOCS_SEMANTIC_URL, source_sha, tool)
    manifest_path = Path(tool) / manifest_rel
    manifest = json.loads(manifest_path.read_text(encoding="utf-8"))
    checkout_commit(repository_url(donor_repository), donor_sha, donor)
    donor_fetches, tool_fetches = reference_fetches(manifest)
    for repository_name, revision in donor_fetches:
        fetch_revision(donor, repository_name, revision)
    for repository_name, revision in tool_fetches:
        fetch_revision(tool, repository_name, revision)
    return manifest_path


def integer(value):
    if value in (None, ""):
        return None
    try:
        return int(value)
    except (TypeError, ValueError):
        return None


def current_pull(env):
    return {
        "number": integer(env.get("PR_NUMBER")),
        "base_repo_id": integer(env.get("PR_BASE_REPO_ID")),
        "base_ref": env.get("PR_BASE_REF") or "",
        "head_repo_id": integer(env.get("PR_HEAD_REPO_ID")),
        "head_ref": env.get("PR_HEAD_REF") or "",
    }


def write_output(values):
    lines = []
    for key, value in values.items():
        text = "" if value is None else str(value)
        if "\n" in text or "\r" in text:
            raise RuntimeError("output value is not one line")
        lines.append(key + "=" + text)
    payload = "\n".join(lines) + "\n"
    path = os.environ.get("GITHUB_OUTPUT")
    if path:
        with open(path, "a", encoding="utf-8") as handle:
            handle.write(payload)
    sys.stdout.write(payload)


def run_check(tool_checkout, manifest_path, donor):
    tool = Path(tool_checkout) / "acceptance" / "tools" / "manifest" / "main.py"
    completed = subprocess.run(
        [sys.executable, str(tool), "check", "--manifest", str(manifest_path), "--donor", str(donor)],
        cwd=str(tool_checkout),
        text=True,
        capture_output=True,
    )
    return completed.returncode == 0


def gate(env=None):
    env = os.environ if env is None else env
    values = {
        "mode": "full",
        "skip_full": "false",
        "source_sha": "",
        "manifest": "",
        "donor_repository": "",
        "donor_sha": "",
        "packages": "",
    }
    try:
        source_sha = resolve_ref(DOCS_SEMANTIC_URL, "refs/heads/main")
        values["source_sha"] = source_sha
        tool = Path(env.get("RUNNER_TEMP") or ".tmp") / "docs-semantic"
        donor_root = Path(env.get("RUNNER_TEMP") or ".tmp") / "donor"
        checkout_commit(DOCS_SEMANTIC_URL, source_sha, tool)
        manifests = load_manifests(tool)
        changed = changed_files(Path.cwd(), env["BASE_SHA"], env["HEAD_SHA"])
        manifest = select_manifest(manifests, [item["path"] for item in changed])
        try:
            pulls = list_pulls(env["GITHUB_REPOSITORY"])
        except (OSError, RuntimeError, urllib.error.URLError, json.JSONDecodeError, KeyError):
            pulls = None
        current = current_pull(env)
        position = stack_position(pulls, current)
        check_ok = False
        if manifest is not None:
            donor_name = manifest["donor"]["repository"]
            donor_sha = manifest["donor"]["sha"]
            checkout_commit(repository_url(donor_name), donor_sha, donor_root)
            donor_fetches, tool_fetches = reference_fetches(manifest)
            for repository_name, revision in donor_fetches:
                fetch_revision(donor_root, repository_name, revision)
            for repository_name, revision in tool_fetches:
                fetch_revision(tool, repository_name, revision)
            manifest_path = manifest["_path"]
            check_ok = run_check(tool, manifest_path, donor_root)
            values["manifest"] = manifest_path.relative_to(tool).as_posix()
            values["donor_repository"] = donor_name
            values["donor_sha"] = donor_sha
        decision = classify_batch(position, manifest, changed, check_ok)
        values["mode"] = decision["mode"]
        values["skip_full"] = "true" if decision["mode"] in FAST_MODES else "false"
        values["packages"] = " ".join(packages(changed))
    except (OSError, RuntimeError, KeyError, subprocess.SubprocessError, urllib.error.URLError, json.JSONDecodeError) as error:
        print("unresolved: " + str(error), file=sys.stderr)
        values["mode"] = "full"
        values["skip_full"] = "false"
    write_output(values)
    print("mode=" + values["mode"], file=sys.stderr)
    return 0


def main(argv=None):
    parser = argparse.ArgumentParser(description=__doc__)
    commands = parser.add_subparsers(dest="command", required=True)
    commands.add_parser("gate")
    prepare = commands.add_parser("prepare")
    prepare.add_argument("--source-sha", required=True)
    prepare.add_argument("--manifest", required=True)
    prepare.add_argument("--donor-repository", required=True)
    prepare.add_argument("--donor-sha", required=True)
    prepare.add_argument("--tool", required=True)
    prepare.add_argument("--donor", required=True)
    args = parser.parse_args(argv)
    if args.command == "gate":
        return gate()
    prepare_inputs(
        args.source_sha,
        args.manifest,
        args.donor_repository,
        args.donor_sha,
        args.tool,
        args.donor,
    )
    return 0


if __name__ == "__main__":
    sys.exit(main())
