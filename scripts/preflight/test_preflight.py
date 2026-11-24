"""Synthetic preflight, closure, comparison, and job-selection tests.

The manifest tool is the published docs-semantic checker. These tests copy it
into a temporary tree and never write the provisioned checkout.
"""

import hashlib
import json
import os
import shutil
import subprocess
import sys
import tempfile
import unittest
from pathlib import Path
from unittest import mock


REPO = Path(__file__).resolve().parents[2]
sys.path.insert(0, str(Path(__file__).resolve().parent))
sys.path.insert(0, str(REPO / "scripts" / "tests"))

import driver_imports_test
import schedule


TOOL_FILES = ("main.py", "manifestlib.py", "closure.py", "gitstore.py", "patterns.py")
BASELINE_PATHS = (
    REPO / "scripts/size-baseline.json",
    REPO / "scripts/behavior-name-baseline.json",
    REPO / "scripts/go-lint-baseline.json",
)
WORKFLOWS = REPO / ".github" / "workflows"
PREFLIGHT = REPO / "scripts" / "preflight" / "preflight.py"
READ_PATH = "driver/vaillant/read.go"
WRITE_PATH = "driver/vaillant/write.go"
OTHER_PATH = "driver/other/wide.go"
ACCEPT_ROOT = REPO / ".tmp" / "dp11-accept"
BASELINE_BYTES = {}


def tool_dir(root):
    root = Path(root)
    nested = root / "acceptance" / "tools" / "manifest" / "main.py"
    if nested.is_file():
        return nested.parent
    if (root / "main.py").is_file() and (root / "manifestlib.py").is_file():
        return root
    return None


def checkout_roots():
    """The published tool is not in this repository. The local gate sets no checkout variable."""
    roots = []
    override = os.environ.get("QQ06_CHECKOUT")
    if override:
        roots.append(Path(override))
    htrace = os.environ.get("HTRACE_DIR")
    if htrace:
        resolved = Path(htrace).resolve()
        if len(resolved.parents) > 1:
            roots.append(resolved.parents[1] / "mainviews" / "docs-semantic")
    scratch = os.environ.get("SC")
    if scratch:
        roots.append(Path(scratch) / "mainviews" / "docs-semantic")
    anchors = [REPO, Path(__file__).resolve()]
    for anchor in anchors:
        for parent in [anchor, *anchor.parents]:
            roots.append(parent / "pipeline-home" / "scratchpad" / "mainviews" / "docs-semantic")
    volumes = Path("/Volumes")
    if volumes.is_dir():
        try:
            mounted = list(volumes.iterdir())
        except OSError:
            mounted = []
        for volume in mounted:
            roots.append(volume / "pipeline-home" / "scratchpad" / "mainviews" / "docs-semantic")
    return roots


def discover_tool():
    seen = set()
    for root in checkout_roots():
        key = str(root)
        if key in seen:
            continue
        seen.add(key)
        try:
            found = tool_dir(root)
        except OSError:
            continue
        if found is not None:
            return found
    raise RuntimeError("manifest tool checkout is unavailable")


def setUpModule():
    cache = REPO / ".tmp" / "gocache"
    cache.mkdir(parents=True, exist_ok=True)
    os.environ["GOCACHE"] = str(cache)
    os.environ["GOWORK"] = "off"
    os.environ["GOPROXY"] = "off"
    os.environ["GOSUMDB"] = "off"
    os.environ["GOTOOLCHAIN"] = "local"
    os.environ["GOFLAGS"] = "-mod=mod"
    os.environ["PYTHONDONTWRITEBYTECODE"] = "1"
    discover_tool()
    for path in BASELINE_PATHS:
        BASELINE_BYTES[path] = path.read_bytes()


def tearDownModule():
    for path, data in BASELINE_BYTES.items():
        if path.read_bytes() != data:
            raise AssertionError("repository baseline changed: " + path.name)


def git(repo, *args, check=True):
    completed = subprocess.run(
        ["git", "-c", "safe.directory=*", "-C", str(repo), *args],
        text=True,
        capture_output=True,
    )
    if check and completed.returncode != 0:
        detail = (completed.stderr or completed.stdout).strip()
        raise RuntimeError(detail or "git " + " ".join(args))
    return completed


def init_repo(path):
    path.mkdir(parents=True, exist_ok=True)
    git(path, "init", "-b", "main")
    git(path, "config", "user.email", "preflight@example.com")
    git(path, "config", "user.name", "Preflight Fixture")
    git(path, "config", "commit.gpgsign", "false")
    git(path, "config", "core.autocrlf", "false")


def commit_all(path, message):
    git(path, "add", "-A")
    git(path, "commit", "-m", message)
    return git(path, "rev-parse", "HEAD").stdout.strip()


def physical_lines(text):
    if not text:
        return 0
    count = text.count("\n")
    if not text.endswith("\n"):
        count += 1
    return count


def letter_name(index):
    value = index + 1
    chars = []
    while value:
        value, remainder = divmod(value - 1, 26)
        chars.append(chr(ord("A") + remainder))
    return "Read" + "".join(reversed(chars))


def function_source(name, span):
    if span < 2:
        raise ValueError("span")
    if span == 2:
        body = "func %s() {\n}\n" % name
    else:
        inner = ["// keep"] * (span - 3) + ["\treturn 0"]
        body = "func %s() int {\n%s\n}\n" % (name, "\n".join(inner))
    return "// %s returns a synthetic reading.\n%s" % (name, body)


def file_source(functions, file_lines):
    text = "// Package vaillant is a synthetic candidate.\npackage vaillant\n\n"
    for name, span in functions:
        text += function_source(name, span)
    current = physical_lines(text)
    if current > file_lines:
        raise ValueError("file already has %s lines" % current)
    text += "".join("// pad\n" for _ in range(file_lines - current))
    if physical_lines(text) != file_lines:
        raise ValueError("line count")
    return text


def small_source(name):
    return file_source([(name, 2)], 6)


def run_tool(main, args, cwd):
    return subprocess.run(
        [sys.executable, str(main), *args],
        cwd=str(cwd),
        text=True,
        capture_output=True,
        env=os.environ.copy(),
    )


def copy_tool(source, root):
    dest = root / "acceptance" / "tools" / "manifest"
    dest.mkdir(parents=True)
    for name in TOOL_FILES:
        shutil.copy(source / name, dest / name)
    return dest / "main.py"


def write_json(path, payload):
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(json.dumps(payload) + "\n", encoding="utf-8")


class World:
    def __init__(self, root):
        self.root = Path(root)
        self.tool = self.root / "tool"
        self.donor = self.root / "donor"
        source = discover_tool()
        init_repo(self.tool)
        self.main = copy_tool(source, self.tool)
        write_json(self.tool / "vec" / "reading.json", {
            "fields": {"reading": "1"},
            "inputs": {},
            "selectors": {"command": "0200"},
        })
        write_json(self.tool / "vec" / "delta.json", {
            "approval": "approved-delta",
            "case": "read-case",
            "field": "reading",
            "new": "2",
            "old": "1",
            "old_present": True,
            "source": "pinned-reading",
        })
        self.tool_sha = commit_all(self.tool, "fixture references")
        init_repo(self.donor)
        (self.donor / "poller").mkdir()
        (self.donor / "poller" / "read.go").write_text(
            "package sample\n\nfunc Read() int { return 1 }\n", encoding="utf-8"
        )
        (self.donor / "poller" / "write.go").write_text(
            "package sample\n\nfunc Write() int { return 1 }\n", encoding="utf-8"
        )
        exemplar = self.donor / "exemplar"
        exemplar.mkdir()
        (exemplar / "rows.go").write_text(
            "package exemplar\n\n"
            "// QQ06ROW {\"id\":\"read\",\"signature\":{\"family\":\"ebus\",\"inverse\":\"none\",\"selector_fields\":[\"command\"]}}\n",
            encoding="utf-8",
        )
        self.donor_sha = commit_all(self.donor, "fixture donor")
        self.limit = self.generate("limit-row", [self.move("poller/read.go", READ_PATH)])
        self.poller = self.generate("poller-row", [
            self.move("poller/read.go", READ_PATH),
            self.move("poller/write.go", WRITE_PATH),
        ])
        self.compare = self.generate(
            "read-row",
            [self.move("poller/read.go", READ_PATH)],
            mapping_rows=[self.pattern()],
            expectations=[self.expectation()],
            expected_differences=[self.difference()],
        )
        self.bytes = {}
        for path in (self.tool / "acceptance" / "manifests").iterdir():
            self.bytes[path] = path.read_bytes()

    def move(self, donor_path, destination):
        return {
            "path": donor_path,
            "destination_repository": "ebusreg",
            "destination_path": destination,
        }

    def pattern(self):
        return {
            "row": "read-pattern",
            "kind": "specification",
            "cases": ["read-case"],
            "assertion": "read-assert",
            "old_symbol": "",
            "bom_entry": "",
            "signature": {"family": "ebus", "inverse": "none", "selector_fields": ["command"]},
            "exemplar": {
                "repository": "eebusreg",
                "path": "exemplar/rows.go",
                "revision": self.donor_sha,
            },
            "candidate": {"form": "inline", "operands": {"command": "0200"}},
        }

    def expectation(self):
        return {
            "case": "read-case",
            "kind": "specification",
            "rows": [{
                "row": "read-pattern",
                "assertion": "read-assert",
                "destination": {
                    "repository": "ebusreg",
                    "path": "driver/vaillant/read_test.go",
                    "test": "TestRead",
                },
            }],
            "selectors": {"command": "0200"},
            "inputs": {"raw": "0200", "qualifiers": "header", "clock": "t0", "state": "idle"},
            "references": [{
                "role": "specification",
                "repository": "docs-semantic",
                "revision": self.tool_sha,
                "path": "vec/reading.json",
            }],
        }

    def difference(self):
        return {
            "approval": "approved-delta",
            "case": "read-case",
            "field": "reading",
            "new": "2",
            "old": "1",
            "old_present": True,
            "source": "pinned-reading",
            "repository": "docs-semantic",
            "revision": self.tool_sha,
            "path": "vec/delta.json",
        }

    def generate(self, row, moves, **extra):
        spec = {
            "row": row,
            "donor_repository": "gateway",
            "donor_sha": self.donor_sha,
            "moves": moves,
            "new_files": [],
            "shared": [],
            "specimens": [],
            "assertions": [],
            "mapping_rows": [],
            "expectations": [],
            "expected_differences": [],
        }
        spec.update(extra)
        path = self.tool / "specs" / (row + ".json")
        write_json(path, spec)
        completed = run_tool(self.main, ["generate", "--input", str(path), "--donor", str(self.donor)], self.tool)
        if completed.returncode != 0:
            raise RuntimeError(completed.stdout + completed.stderr)
        manifest = self.tool / "acceptance" / "manifests" / (row + ".json")
        body = json.loads(manifest.read_text(encoding="utf-8"))
        if body.get("receipts") != []:
            raise RuntimeError("candidate receipt array is not empty")
        return manifest

    def manifest_bytes(self):
        return dict(self.bytes)

    def restore(self):
        for path, data in self.bytes.items():
            path.write_bytes(data)


def candidate_tree(parent, files, debt):
    root = Path(tempfile.mkdtemp(prefix="cand-", dir=str(parent)))
    (root / "go.mod").write_text(
        "module github.com/Project-Helianthus/helianthus-ebusreg\n\ngo 1.22\n",
        encoding="utf-8",
    )
    shutil.copy(REPO / ".golangci.yml", root / ".golangci.yml")
    scripts = root / "scripts"
    scripts.mkdir()
    shutil.copy(REPO / "scripts" / "check_go_lint_baseline.py", scripts / "check_go_lint_baseline.py")
    write_json(scripts / "behavior-name-baseline.json", {"version": 1, "exceptions": [], "baseline": []})
    write_json(scripts / "go-lint-baseline.json", {"version": 1, "golangci_lint": "2.11.4", "entries": []})
    write_json(scripts / "size-baseline.json", {"version": 1, "debt": debt})
    for rel, text in files.items():
        dest = root / rel
        dest.parent.mkdir(parents=True, exist_ok=True)
        dest.write_text(text, encoding="utf-8")
    return root


READ_TEST = """package vaillant

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
)

// TestRead compares the pinned reading.
func TestRead(t *testing.T) {
	fmt.Printf("MANIFEST_RESULT case=%s row=%s assertion=%s\\n", os.Getenv("MANIFEST_CASE"), os.Getenv("MANIFEST_ROW"), os.Getenv("MANIFEST_ASSERTION"))
	fmt.Printf("MANIFEST_INPUT raw=%s qualifiers=%s clock=%s state=%s\\n", os.Getenv("MANIFEST_RAW"), os.Getenv("MANIFEST_QUALIFIERS"), os.Getenv("MANIFEST_CLOCK"), os.Getenv("MANIFEST_STATE"))
	var pins map[string]string
	if err := json.Unmarshal([]byte(os.Getenv("MANIFEST_PINNED_FIELDS")), &pins); err != nil {
		t.Fatal(err)
	}
	replacement := ""
	if text, err := os.ReadFile("observed"); err == nil {
		replacement = strings.TrimSpace(string(text))
	}
	for name, value := range pins {
		if strings.HasPrefix(name, "template.") {
			continue
		}
		observed := value
		if name == "reading" && replacement != "" {
			observed = replacement
		}
		fmt.Printf("MANIFEST_FIELD field=%s expected=%s observed=%s\\n", name, value, observed)
	}
}
"""


def compare_files():
    return {
        READ_PATH: small_source("Read"),
        "driver/vaillant/read_test.go": READ_TEST,
    }


def soft_line(path, symbol, metric, value):
    row = {"metric": metric, "path": path, "symbol": symbol, "value": value}
    return "SOFT " + json.dumps(row, sort_keys=True)


def split_line(path, symbol, metric, value, limit):
    return "split required: %s %s %s %s > %s" % (path, symbol, metric, value, limit)


def debt_row(path, symbol, metric, value):
    return {"path": path, "symbol": symbol, "metric": metric, "value": value}


def job_blocks(text):
    blocks = {}
    current = None
    lines = []
    for line in text.splitlines(True):
        if line.startswith("  ") and not line.startswith("   ") and line.rstrip().endswith(":"):
            if current is not None:
                blocks[current] = "".join(lines)
            current = line.strip()[:-1]
            lines = [line]
        elif current is not None:
            lines.append(line)
    if current is not None:
        blocks[current] = "".join(lines)
    return blocks


def pull(number, base_repo, base_ref, head_repo, head_ref):
    return {
        "number": number,
        "base_repo_id": base_repo,
        "base_ref": base_ref,
        "head_repo_id": head_repo,
        "head_ref": head_ref,
    }


class ScheduleTests(unittest.TestCase):
    def test_each_stack_identity_field_changes_the_position(self):
        current = pull(2, 10, "main", 10, "feature")
        parent = pull(1, 10, "older", 10, "main")
        child = pull(3, 10, "feature", 10, "newer")
        self.assertEqual(schedule.stack_position([parent], current), "tip")
        self.assertEqual(schedule.stack_position([child], current), "lower")
        self.assertEqual(schedule.stack_position([parent, child], current), "lower")
        self.assertEqual(schedule.stack_position([], current), "unstacked")
        self.assertEqual(schedule.stack_position(None, current), "unresolved")
        self.assertEqual(schedule.stack_position([parent], None), "unresolved")
        for field in ("number", "base_repo_id", "base_ref", "head_repo_id", "head_ref"):
            broken = dict(current)
            broken[field] = ""
            self.assertEqual(schedule.stack_position([parent, child], broken), "unresolved", field)
        incomplete = dict(parent)
        incomplete["base_repo_id"] = ""
        self.assertEqual(schedule.stack_position([incomplete], current), "unresolved")
        cases = (
            ("head_repo_id", parent, 11, "tip"),
            ("head_ref", parent, "other", "tip"),
            ("base_repo_id", child, 11, "lower"),
            ("base_ref", child, "other", "lower"),
        )
        for field, other, value, absent in cases:
            changed = dict(other)
            changed[field] = value
            self.assertEqual(schedule.stack_position([changed], current), "unstacked", field)
            self.assertNotEqual(schedule.stack_position([changed], current), absent, field)
        for field, value in (("base_repo_id", 11), ("base_ref", "trunk")):
            changed = dict(current)
            changed[field] = value
            self.assertEqual(schedule.stack_position([parent], changed), "unstacked", field)
        for field, value in (("head_repo_id", 11), ("head_ref", "other")):
            changed = dict(current)
            changed[field] = value
            self.assertEqual(schedule.stack_position([child], changed), "unstacked", field)
        copy = dict(current)
        self.assertEqual(schedule.stack_position([copy], current), "unstacked")

    def test_manifest_coverage_uses_destination_paths_and_statuses(self):
        granted = {
            "row": "bb",
            "paths": [{"repository": "ebusreg", "paths": ["driver/rows/one.go", "scripts/ci_local.sh"]}],
        }
        added = [{"path": "driver/rows/one.go", "status": "added"}]
        self.assertIs(schedule.select_manifest([granted], ["driver/rows/one.go"]), granted)
        self.assertTrue(schedule.is_additive(added, ["driver/rows/one.go", "scripts/ci_local.sh"]))
        self.assertEqual(schedule.classify_batch("unstacked", granted, added, True)["mode"], "fast-additive")
        self.assertEqual(schedule.classify_batch("lower", granted, added, True)["mode"], "fast-additive")
        uncovered = added + [{"path": "README.md", "status": "added"}]
        self.assertIsNone(schedule.select_manifest([granted], ["driver/rows/one.go", "README.md"]))
        self.assertEqual(schedule.classify_batch("unstacked", granted, uncovered, True)["mode"], "full")
        self.assertEqual(
            schedule.classify_batch("unstacked", None, added, True)["mode"],
            "full",
        )
        for status in ("modified", "deleted", "renamed"):
            changed = [{"path": "driver/rows/one.go", "status": status}]
            self.assertFalse(schedule.is_additive(changed, ["driver/rows/one.go", "scripts/ci_local.sh"]))
            self.assertEqual(schedule.classify_batch("unstacked", granted, changed, True)["mode"], "full", status)
            self.assertEqual(schedule.classify_batch("lower", granted, changed, True)["mode"], "fast-compare", status)
            self.assertEqual(schedule.classify_batch("tip", granted, changed, True)["mode"], "full", status)
        self.assertEqual(schedule.classify_batch("tip", granted, added, True)["mode"], "full")
        self.assertEqual(schedule.classify_batch("unresolved", granted, added, True)["mode"], "full")
        self.assertEqual(schedule.classify_batch("unresolved", granted, added, False)["mode"], "check-failed")
        self.assertEqual(schedule.classify_batch("unstacked", granted, added, False)["mode"], "check-failed")
        donor_only = {"row": "donor", "paths": [{"repository": "gateway", "paths": ["driver/rows/one.go"]}]}
        self.assertIsNone(schedule.select_manifest([donor_only], ["driver/rows/one.go"]))
        wide_a = {"row": "a", "paths": [{"repository": "ebusreg", "paths": ["a.go", "b.go"]}]}
        wide_b = {"row": "b", "paths": [{"repository": "ebusreg", "paths": ["a.go", "c.go"]}]}
        tight = {"row": "c", "paths": [{"repository": "ebusreg", "paths": ["a.go"]}]}
        self.assertIsNone(schedule.select_manifest([wide_a, wide_b], ["a.go"]))
        self.assertEqual(schedule.select_manifest([wide_a, tight], ["a.go"])["row"], "c")
        self.assertIsNone(schedule.select_manifest([granted], []))
        parsed = schedule.parse_name_status("A\tdriver/rows/one.go\nM\tscripts/ci_local.sh\nD\told.go\nR100\told.go\tdriver/rows/one.go\n")
        self.assertEqual(
            parsed,
            [
                {"path": "driver/rows/one.go", "status": "added"},
                {"path": "scripts/ci_local.sh", "status": "modified"},
                {"path": "old.go", "status": "deleted"},
                {"path": "driver/rows/one.go", "status": "renamed"},
            ],
        )

    def test_job_plan_keeps_driver_imports_and_full_commands(self):
        fast = schedule.job_plan("fast-additive")
        self.assertEqual(fast["jobs"]["driver-imports"], "run")
        self.assertEqual(fast["jobs"]["test"], "skip")
        self.assertEqual(fast["jobs"]["additive-rows"], "run")
        self.assertEqual(
            fast["steps"],
            ["package-tests", "check", "compare", "preflight", "size", "behavior-names", "lint"],
        )
        compared = schedule.job_plan("fast-compare")
        self.assertEqual(compared["jobs"]["driver-imports"], "run")
        self.assertEqual(compared["jobs"]["test"], "skip")
        self.assertEqual(compared["steps"], ["package-tests", "check", "compare"])
        self.assertNotEqual(compared["steps"], [])
        failed = schedule.job_plan("check-failed")
        self.assertEqual(failed["jobs"]["test"], "run")
        self.assertEqual(failed["jobs"]["additive-rows"], "run")
        self.assertEqual(failed["steps"], ["check", "compare"])
        full = schedule.job_plan("full")
        self.assertEqual(full["jobs"]["test"], "run")
        self.assertEqual(full["jobs"]["driver-imports"], "run")
        self.assertEqual(full["jobs"]["additive-rows"], "skip")
        self.assertEqual(full["steps"], [])
        self.assertEqual(full["commands"]["test"], "go test -race -count=1 -timeout 20m ./...")
        self.assertEqual(full["commands"]["size-ratchet"], "python3 scripts/check_size_ratchet.py")
        self.assertEqual(full["commands"]["behavior-names"], "python3 scripts/check_behavior_names.py --check")
        self.assertEqual(full["commands"]["readability"], "python3 scripts/check_go_lint_baseline.py --check")
        self.assertEqual(
            full["commands"]["driver-imports"],
            "python3 -m unittest discover -s scripts/tests -p driver_imports_test.py",
        )
        self.assertEqual(schedule.packages([
            {"path": "driver/vaillant/read.go"},
            {"path": "driver/vaillant/write.go"},
            {"path": "README.md"},
            {"path": "doc.go"},
        ]), ["./driver/vaillant", "."])

    def test_forbidden_production_import_fails_on_the_fast_path(self):
        self.assertEqual(schedule.job_plan("fast-additive")["jobs"]["driver-imports"], "run")
        self.assertEqual(schedule.job_plan("fast-compare")["jobs"]["driver-imports"], "run")
        directory = tempfile.TemporaryDirectory()
        self.addCleanup(directory.cleanup)
        root = Path(directory.name)
        (root / "go.mod").write_text("module github.com/Project-Helianthus/helianthus-ebusreg\n\ngo 1.22\n")
        package = root / "driver" / "rows"
        package.mkdir(parents=True)
        forbidden = "github.com/Project-Helianthus/helianthus-semreg/semreg/v1/graph"
        allowed = "github.com/Project-Helianthus/helianthus-semreg/semreg/v1/driver"
        (package / "row.go").write_text('package rows\n\nimport "%s"\n' % forbidden)
        lines = driver_imports_test.findings(root)
        self.assertTrue(lines)
        self.assertIn(forbidden, lines[0])
        (package / "row.go").write_text('package rows\n\nimport "%s"\n' % allowed)
        self.assertEqual([], driver_imports_test.findings(root))

    def test_unresolved_relationships_still_run_check(self):
        sha = "%040x" % 1
        directory = tempfile.TemporaryDirectory()
        self.addCleanup(directory.cleanup)
        tool = Path(directory.name) / "docs-semantic"
        manifest = {
            "row": "row",
            "donor": {"repository": "gateway", "sha": sha},
            "paths": [{"repository": "ebusreg", "paths": ["driver/rows/one.go"]}],
            "patterns": [],
            "expectations": [],
            "expected_differences": [],
            "_path": tool / "acceptance" / "manifests" / "row.json",
        }
        output = Path(directory.name) / "github-output"
        previous = os.environ.get("GITHUB_OUTPUT")
        os.environ["GITHUB_OUTPUT"] = str(output)
        env = {
            "BASE_SHA": "base",
            "HEAD_SHA": "head",
            "GITHUB_REPOSITORY": "Project-Helianthus/helianthus-ebusreg",
            "PR_NUMBER": "2",
            "PR_BASE_REF": "main",
            "PR_BASE_REPO_ID": "10",
            "PR_HEAD_REF": "feature",
            "PR_HEAD_REPO_ID": "10",
            "RUNNER_TEMP": directory.name,
        }
        try:
            for check_ok, mode in ((True, "full"), (False, "check-failed")):
                with self.subTest(check_ok=check_ok):
                    output.write_text("", encoding="utf-8")
                    with mock.patch.object(schedule, "resolve_ref", return_value=sha), \
                            mock.patch.object(schedule, "checkout_commit"), \
                            mock.patch.object(schedule, "load_manifests", return_value=[manifest]), \
                            mock.patch.object(schedule, "changed_files", return_value=[{"path": "driver/rows/one.go", "status": "added"}]), \
                            mock.patch.object(schedule, "list_pulls", side_effect=OSError("unavailable")), \
                            mock.patch.object(schedule, "fetch_revision"), \
                            mock.patch.object(schedule, "run_check", return_value=check_ok) as checked:
                        self.assertEqual(schedule.gate(env), 0)
                    self.assertTrue(checked.called)
                    recorded = output.read_text(encoding="utf-8")
                    self.assertIn("mode=" + mode, recorded)
                    if mode == "full":
                        self.assertIn("skip_full=false", recorded)
                    else:
                        self.assertIn("skip_full=false", recorded)
        finally:
            if previous is None:
                os.environ.pop("GITHUB_OUTPUT", None)
            else:
                os.environ["GITHUB_OUTPUT"] = previous

    def test_public_gateway_name_is_the_renamed_repository(self):
        self.assertEqual(schedule.public_repository("gateway"), "helianthus-gateway")
        self.assertEqual(schedule.public_repository("helianthus-ebusgateway"), "helianthus-gateway")
        self.assertEqual(
            schedule.repository_url("gateway"),
            "https://github.com/Project-Helianthus/helianthus-gateway.git",
        )

    def test_workflows_select_jobs_without_replacing_full_commands(self):
        retained = {
            "ci.yml": (
                "go test -race -count=1 -timeout 20m ./...",
                "go vet ./...",
                "go build ./...",
                "python3 -m unittest discover -s scripts/tests -p driver_imports_test.py",
            ),
            "size-ratchet.yml": ("python3 scripts/check_size_ratchet.py",),
            "behavior-names.yml": ("python3 scripts/check_behavior_names.py --check",),
            "go-readability.yml": (
                "version: v2.11.4",
                "install-only: true",
                "python3 scripts/check_go_lint_baseline.py --check",
            ),
        }
        for name, snippets in retained.items():
            text = (WORKFLOWS / name).read_text(encoding="utf-8")
            blocks = job_blocks(text)
            self.assertNotIn("full suite applies", text, name)
            self.assertNotIn("cron:", text, name)
            self.assertIn("python3 scripts/preflight/schedule.py gate", text, name)
            self.assertIn("needs.classify.outputs.skip_full", text, name)
            for snippet in snippets:
                self.assertIn(snippet, text, name)
            if name == "ci.yml":
                self.assertNotIn("skip_full", blocks["driver-imports"])
                self.assertNotIn("needs:", blocks["driver-imports"])
                for job in ("terminology", "build", "test", "tinygo", "lint"):
                    self.assertIn("skip_full", blocks[job], job)
        additive = (WORKFLOWS / "additive-rows.yml").read_text(encoding="utf-8")
        self.assertNotIn("full suite applies", additive)
        self.assertNotIn("cron:", additive)
        self.assertIn("go test -race -count=1 -timeout 20m", additive)
        self.assertIn("acceptance/tools/manifest/main.py", additive)
        self.assertIn("compare --manifest", additive)
        self.assertIn("check --manifest", additive)
        self.assertIn("continue-on-error: true", additive)
        self.assertIn('echo "check=${{ steps.check.outcome }}"', additive)
        self.assertIn('echo "compare=${{ steps.compare.outcome }}"', additive)
        self.assertIn("scripts/preflight/preflight.py", additive)
        self.assertIn("python3 scripts/check_size_ratchet.py", additive)
        self.assertIn("python3 scripts/check_behavior_names.py --check", additive)
        self.assertIn("python3 scripts/check_go_lint_baseline.py --check", additive)
        self.assertIn("actions/setup-python@v5", additive)
        self.assertIn("install-only: true", additive)
        self.assertIn("fast-additive", additive)
        self.assertIn("fast-compare", additive)
        self.assertIn("check-failed", additive)
        blocks = job_blocks(additive)
        self.assertIn("synology", blocks["additive-rows"])
        self.assertIn("GIT_CONFIG_GLOBAL", blocks["additive-rows"])


class PreflightTests(unittest.TestCase):
    world = None

    @classmethod
    def setUpClass(cls):
        cls.scratch = tempfile.TemporaryDirectory(prefix="dp11-", dir=str(REPO / ".tmp"))
        cls.world = World(Path(cls.scratch.name) / "world")
        cls.addClassCleanup(cls.scratch.cleanup)

    def setUp(self):
        self.world.restore()

    def candidate(self, files, debt=None):
        root = candidate_tree(self.scratch.name, files, debt or [])
        self.addCleanup(shutil.rmtree, root, True)
        return root

    def tracked(self, root):
        return [
            root / "scripts" / "size-baseline.json",
            root / "scripts" / "behavior-name-baseline.json",
            root / "scripts" / "go-lint-baseline.json",
        ]

    def invoke(self, root, manifest, donor=None, include_donor=True):
        before = {path: path.read_bytes() for path in self.tracked(root)}
        published = self.world.manifest_bytes()
        command = [sys.executable, str(PREFLIGHT), "--manifest", str(manifest)]
        if include_donor:
            command.extend(["--donor", str(self.world.donor if donor is None else donor)])
        completed = subprocess.run(command, cwd=str(root), text=True, capture_output=True, env=os.environ.copy())
        for path, data in before.items():
            self.assertEqual(path.read_bytes(), data, path)
        for path, data in published.items():
            self.assertEqual(path.read_bytes(), data, path)
        return completed

    def test_limits_follow_native_evaluation(self):
        cases = (
            ("function-151", file_source([("ReadLong", 151)], 155), [], "function_lines", "ReadLong", 151, 150, True),
            ("file-1001", file_source([("ReadWide", 3)], 1001), [], "file_lines", "<file>", 1001, 1000, True),
            (
                "file-held",
                file_source([("ReadWide", 3)], 1002),
                [debt_row(READ_PATH, "<file>", "file_lines", 1002)],
                "file_lines", "<file>", 1002, 1002, False,
            ),
            (
                "file-grown",
                file_source([("ReadWide", 3)], 1002),
                [debt_row(READ_PATH, "<file>", "file_lines", 1001)],
                "file_lines", "<file>", 1002, 1001, True,
            ),
            (
                "function-held",
                file_source([("ReadLong", 152)], 156),
                [debt_row(READ_PATH, "ReadLong", "function_lines", 152)],
                "function_lines", "ReadLong", 152, 152, False,
            ),
            (
                "function-grown",
                file_source([("ReadLong", 152)], 156),
                [debt_row(READ_PATH, "ReadLong", "function_lines", 151)],
                "function_lines", "ReadLong", 152, 151, True,
            ),
            (
                "count-held",
                file_source([(letter_name(index), 2) for index in range(62)], 189),
                [debt_row(READ_PATH, "<file>", "functions_per_file", 62)],
                "functions_per_file", "<file>", 62, 62, False,
            ),
            (
                "count-grown",
                file_source([(letter_name(index), 2) for index in range(62)], 189),
                [debt_row(READ_PATH, "<file>", "functions_per_file", 61)],
                "functions_per_file", "<file>", 62, 61, True,
            ),
        )
        for name, source, debt, metric, symbol, value, limit, splits in cases:
            with self.subTest(name=name):
                root = self.candidate({READ_PATH: source}, debt)
                completed = self.invoke(root, self.world.limit)
                lines = [line for line in completed.stdout.splitlines() if line]
                self.assertIn(soft_line(READ_PATH, symbol, metric, value), lines)
                wanted = split_line(READ_PATH, symbol, metric, value, limit)
                if splits:
                    self.assertIn(wanted, lines)
                    self.assertNotEqual(completed.returncode, 0)
                else:
                    self.assertFalse([line for line in lines if line.startswith("split required:")], lines)
                    self.assertEqual(completed.returncode, 0, completed.stdout + completed.stderr)
                self.assertNotIn("receipt", completed.stdout)

    def test_file_under_the_soft_budget_prints_nothing(self):
        root = self.candidate({READ_PATH: small_source("Read")})
        completed = self.invoke(root, self.world.limit)
        self.assertEqual(completed.returncode, 0, completed.stdout + completed.stderr)
        self.assertEqual(completed.stdout, "")

    def test_another_file_outside_the_manifest_still_fails(self):
        root = self.candidate({
            READ_PATH: small_source("Read"),
            OTHER_PATH: file_source([("ReadWide", 3)], 1001),
        })
        completed = self.invoke(root, self.world.limit)
        self.assertNotEqual(completed.returncode, 0)
        self.assertIn(split_line(OTHER_PATH, "<file>", "file_lines", 1001, 1000), completed.stdout.splitlines())

    def test_poller_destinations_under_the_soft_budget_print_nothing(self):
        root = self.candidate({
            READ_PATH: small_source("Read"),
            WRITE_PATH: small_source("Write"),
        })
        completed = self.invoke(root, self.world.poller)
        self.assertEqual(completed.returncode, 0, completed.stdout + completed.stderr)
        self.assertEqual(completed.stdout, "")
        body = json.loads(self.world.poller.read_text(encoding="utf-8"))
        self.assertEqual(body["receipts"], [])

    def test_incomplete_closure_is_reported_unchanged(self):
        mutated = self.world.limit.parent / "limit-open.json"
        body = json.loads(self.world.limit.read_text(encoding="utf-8"))
        body["symbols"][0]["name"] = body["symbols"][0]["name"] + "Other"
        mutated.write_text(json.dumps(body), encoding="utf-8")
        root = self.candidate({READ_PATH: small_source("Read")})
        completed = self.invoke(root, mutated)
        self.assertNotEqual(completed.returncode, 0)
        self.assertIn("symbol-set", completed.stdout.splitlines())

    def test_pattern_signature_mismatch_is_reported_unchanged(self):
        mutated = self.world.compare.parent / "read-stale.json"
        body = json.loads(self.world.compare.read_text(encoding="utf-8"))
        body["patterns"][0]["signature"]["selector_fields"] = ["command", "extra"]
        mutated.write_text(json.dumps(body), encoding="utf-8")
        root = self.candidate(compare_files())
        completed = self.invoke(root, mutated)
        self.assertNotEqual(completed.returncode, 0)
        self.assertIn("signature-mismatch row: read-pattern", completed.stdout.splitlines())

    def test_wrong_donor_stops_before_limits(self):
        later = Path(self.scratch.name) / "donor-later"
        git(self.world.root, "clone", str(self.world.donor), str(later))
        git(later, "config", "user.email", "preflight@example.com")
        git(later, "config", "user.name", "Preflight Fixture")
        git(later, "commit", "--allow-empty", "-m", "later")
        root = self.candidate({READ_PATH: file_source([("ReadLong", 151)], 155)})
        completed = self.invoke(root, self.world.limit, donor=later)
        self.assertNotEqual(completed.returncode, 0)
        self.assertIn("donor-sha", completed.stdout.splitlines())
        self.assertFalse([line for line in completed.stdout.splitlines() if line.startswith("split required:")])

    def test_omitted_donor_stops(self):
        root = self.candidate({READ_PATH: small_source("Read")})
        completed = self.invoke(root, self.world.limit, include_donor=False)
        self.assertNotEqual(completed.returncode, 0)
        self.assertIn("donor checkout is required", completed.stdout.splitlines())

    def test_destination_lint_finding_uses_the_native_collector(self):
        text = "package vaillant\n\n// Read returns a synthetic reading.\nfunc Read() int {\n\treturn 1\n}\n"
        root = self.candidate({READ_PATH: text})
        completed = self.invoke(root, self.world.limit)
        self.assertNotEqual(completed.returncode, 0)
        self.assertTrue(
            any(line.startswith("rename or fix required:") and READ_PATH in line for line in completed.stdout.splitlines()),
            completed.stdout + completed.stderr,
        )

    def test_check_compare_and_preflight_are_separate(self):
        root = self.candidate(compare_files())
        check = run_tool(
            self.world.main,
            ["check", "--manifest", str(self.world.compare), "--donor", str(self.world.donor)],
            self.world.tool,
        )
        self.assertEqual(check.returncode, 0, check.stdout + check.stderr)
        self.assertIn("check: ok", check.stdout)
        equal = run_tool(
            self.world.main,
            ["compare", "--manifest", str(self.world.compare), "--destination", str(root)],
            self.world.tool,
        )
        self.assertEqual(equal.returncode, 0, equal.stdout + equal.stderr)
        self.assertIn("compare: ok", equal.stdout)
        self.assertIn("field=reading", equal.stdout)
        self.assertIn("expected=1", equal.stdout)
        self.assertIn("observed=1", equal.stdout)
        self.assertIn("result=pass", equal.stdout)
        preflight = self.invoke(root, self.world.compare)
        self.assertEqual(preflight.returncode, 0, preflight.stdout + preflight.stderr)
        self.assertEqual(preflight.stdout, "")
        self.assertEqual(check.returncode, 0)
        self.assertEqual(equal.returncode, 0)
        self.assertEqual(preflight.returncode, 0)
        approved = self.candidate(compare_files())
        (approved / "driver" / "vaillant" / "observed").write_text("2\n", encoding="utf-8")
        approved_run = run_tool(
            self.world.main,
            ["compare", "--manifest", str(self.world.compare), "--destination", str(approved)],
            self.world.tool,
        )
        self.assertEqual(approved_run.returncode, 0, approved_run.stdout + approved_run.stderr)
        self.assertIn("observed=2", approved_run.stdout)
        self.assertIn("compare: ok", approved_run.stdout)
        mismatch = self.candidate(compare_files())
        (mismatch / "driver" / "vaillant" / "observed").write_text("9\n", encoding="utf-8")
        mismatch_run = run_tool(
            self.world.main,
            ["compare", "--manifest", str(self.world.compare), "--destination", str(mismatch)],
            self.world.tool,
        )
        self.assertNotEqual(mismatch_run.returncode, 0, mismatch_run.stdout)
        self.assertEqual(check.returncode, 0)
        self.assertIn("observed=9", mismatch_run.stdout)
        self.assertIn("result=fail", mismatch_run.stdout)
        self.assertNotIn("compare: ok", mismatch_run.stdout)
        missing = self.world.compare.parent / "read-missing.json"
        body = json.loads(self.world.compare.read_text(encoding="utf-8"))
        body["expectations"][0]["references"][0]["path"] = "vec/missing.json"
        missing.write_text(json.dumps(body), encoding="utf-8")
        missing_check = run_tool(
            self.world.main,
            ["check", "--manifest", str(missing), "--donor", str(self.world.donor)],
            self.world.tool,
        )
        missing_compare = run_tool(
            self.world.main,
            ["compare", "--manifest", str(missing), "--destination", str(root)],
            self.world.tool,
        )
        self.assertNotEqual(missing_check.returncode, 0)
        self.assertNotEqual(missing_compare.returncode, 0)
        self.assertIn("unavailable:", missing_check.stdout)
        self.assertIn("unavailable:", missing_compare.stdout)
        mutated = self.world.compare.parent / "read-mutated.json"
        body = json.loads(self.world.compare.read_text(encoding="utf-8"))
        body["expectations"][0]["references"][0]["blob_sha256"] = hashlib.sha256(b"other").hexdigest()
        mutated.write_text(json.dumps(body), encoding="utf-8")
        mutated_check = run_tool(
            self.world.main,
            ["check", "--manifest", str(mutated), "--donor", str(self.world.donor)],
            self.world.tool,
        )
        mutated_compare = run_tool(
            self.world.main,
            ["compare", "--manifest", str(mutated), "--destination", str(root)],
            self.world.tool,
        )
        self.assertNotEqual(mutated_check.returncode, 0)
        self.assertNotEqual(mutated_compare.returncode, 0)
        self.assertIn("mutated:", mutated_check.stdout)
        self.assertIn("mutated:", mutated_compare.stdout)


def acceptance(command):
    setUpModule()
    root = Path(tempfile.mkdtemp(prefix="dp11-accept-", dir=str(REPO / ".tmp")))
    world = World(root / "world")
    candidate = candidate_tree(root, compare_files(), [])
    if command == "check":
        completed = run_tool(
            world.main,
            ["check", "--manifest", str(world.compare), "--donor", str(world.donor)],
            world.tool,
        )
    elif command == "compare":
        completed = run_tool(
            world.main,
            ["compare", "--manifest", str(world.compare), "--destination", str(candidate)],
            world.tool,
        )
    elif command == "preflight":
        completed = subprocess.run(
            [sys.executable, str(PREFLIGHT), "--manifest", str(world.compare), "--donor", str(world.donor)],
            cwd=str(candidate),
            text=True,
            capture_output=True,
            env=os.environ.copy(),
        )
    else:
        raise SystemExit("unknown acceptance command")
    sys.stdout.write(completed.stdout)
    sys.stderr.write(completed.stderr)
    return completed.returncode


if __name__ == "__main__":
    if len(sys.argv) > 1 and sys.argv[1] == "acceptance":
        raise SystemExit(acceptance(sys.argv[2]))
    unittest.main()
