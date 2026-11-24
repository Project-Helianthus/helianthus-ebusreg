"""Driver import direction for helianthus-ebusreg (DP01, issue #531).

Imports point down only (RULES rule 3): a driver may import the eBUS transport
helianthus-ebusgo, the SemReg driver contract, and this module. No repository
imports the gateway. Only production imports (`go list` field Imports) are checked.
Synthetic fixtures use import paths that `go list -e` reports unresolved, so no
helianthus-semreg package is compiled.
"""
import json
import os
from pathlib import Path
import subprocess
import tempfile
import unittest

# Module roots (RULES rule 3; OPERATOR item 7).
GATEWAY_MODULE = "github.com/Project-Helianthus/helianthus-gateway"
SEMREG_MODULE = "github.com/Project-Helianthus/helianthus-semreg"
EBUSGO_MODULE = "github.com/Project-Helianthus/helianthus-ebusgo"
ORG_PREFIX = "github.com/Project-Helianthus/"

# The six exact SemReg import paths a driver may use, each matched as a whole path.
# driver: rule 3 and the queue note "Allow scoped contract leaf" (SS02).
# roweval: SS02's approved row design; EX01 imports it.
# semreg/v1: the types those leaves expose (EX01).
# catalog, operation: the public Descriptor and SettingRoute types (SS01 items 2 and 5; SR20).
# packs/thermal: the exact public pack validator used by EE03.
SEMREG_CONTRACT_PATHS = (
    SEMREG_MODULE + "/semreg/v1/driver",
    SEMREG_MODULE + "/semreg/v1/roweval",
    SEMREG_MODULE + "/semreg/v1",
    SEMREG_MODULE + "/semreg/v1/catalog",
    SEMREG_MODULE + "/semreg/v1/operation",
    SEMREG_MODULE + "/semreg/v1/packs/thermal",
)

# Rule names printed in findings. Sources: rule 3 (gateway and SemReg direction) and
# rule 3 with OPERATOR item 7 (other Helianthus modules).
RULE_GATEWAY = "no repository imports the gateway (rule 3)"
RULE_SEMREG = "SemReg is imported only through its exact public contract paths (rule 3; SR20)"
RULE_OTHER_MODULE = "only helianthus-ebusgo and this module are imported from Helianthus (rule 3; OPERATOR item 7)"


def under(import_path, module):
    return import_path == module or import_path.startswith(module + "/")


def owning_module(import_path):
    """Return the Helianthus module of an import path, or None outside ORG_PREFIX."""
    if not import_path.startswith(ORG_PREFIX):
        return None
    return ORG_PREFIX + import_path[len(ORG_PREFIX):].split("/")[0]


def own_module(module_root):
    for line in (Path(module_root) / "go.mod").read_text().splitlines():
        fields = line.split()
        if len(fields) == 2 and fields[0] == "module":
            return fields[1]
    raise ValueError("go.mod has no module directive")


def violation(import_path, own):
    if under(import_path, GATEWAY_MODULE):
        return RULE_GATEWAY
    if under(import_path, SEMREG_MODULE):
        return None if import_path in SEMREG_CONTRACT_PATHS else RULE_SEMREG
    module = owning_module(import_path)
    if module is None or module in (EBUSGO_MODULE, own):
        return None
    return RULE_OTHER_MODULE


def list_packages(module_root):
    env = dict(os.environ, GOWORK="off", GOPROXY="off", GOSUMDB="off",
               GOTOOLCHAIN="local", GOFLAGS="-buildvcs=false")
    result = subprocess.run(["go", "list", "-e", "-json", "./..."], cwd=module_root, env=env,
                            text=True, capture_output=True, timeout=300)
    if result.returncode:
        raise RuntimeError("go list failed: " + result.stderr.strip())
    decoder = json.JSONDecoder()
    packages, position, text = [], 0, result.stdout
    while position < len(text):
        if text[position].isspace():
            position += 1
            continue
        package, position = decoder.raw_decode(text, position)
        packages.append(package)
    return packages


def findings(module_root):
    own = own_module(module_root)
    lines = []
    for package in list_packages(module_root):
        for import_path in package.get("Imports") or []:
            rule = violation(import_path, own)
            if rule:
                lines.append(f"{package['ImportPath']} imports {import_path}: {rule}")
    return sorted(lines)


REPOSITORY_ROOT = Path(__file__).resolve().parents[2]


def source(imports):
    lines = ["package sample", ""]
    if imports:
        lines.append("import (")
        lines.extend(f'\t"{path}"' for path in imports)
        lines.append(")")
    return "\n".join(lines) + "\n"


class DriverImportTests(unittest.TestCase):
    def fixture(self, imports, test_imports=()):
        directory = tempfile.TemporaryDirectory()
        self.addCleanup(directory.cleanup)
        root = Path(directory.name)
        (root / "go.mod").write_text("module github.com/Project-Helianthus/helianthus-ebusreg\n\ngo 1.22\n")
        package = root / "driver" / "sample"
        package.mkdir(parents=True)
        (package / "sample.go").write_text(source(imports))
        if test_imports:
            (package / "sample_test.go").write_text(source(test_imports))
        return findings(root)

    def assertOneFinding(self, lines, package, import_path):
        self.assertEqual(1, len(lines), lines)
        self.assertIn(package, lines[0])
        self.assertIn(import_path, lines[0])

    def test_repository_passes(self):
        self.assertEqual([], findings(REPOSITORY_ROOT))

    def test_gateway_import_is_rejected(self):
        imports = ["github.com/Project-Helianthus/helianthus-gateway/graphql"]
        self.assertOneFinding(self.fixture(imports), "driver/sample", imports[0])

    def test_semreg_host_import_is_rejected(self):
        imports = ["github.com/Project-Helianthus/helianthus-semreg/semreg/host"]
        self.assertOneFinding(self.fixture(imports), "driver/sample", imports[0])

    def test_other_driver_import_is_rejected(self):
        imports = ["github.com/Project-Helianthus/helianthus-otherdriver/api"]
        self.assertOneFinding(self.fixture(imports), "driver/sample", imports[0])

    def test_semreg_prefixed_path_is_rejected(self):
        imports = ["github.com/Project-Helianthus/helianthus-semreg/semreg/v1/driver/runtime"]
        self.assertOneFinding(self.fixture(imports), "driver/sample", imports[0])

    def test_semreg_non_contract_packages_are_rejected(self):
        for suffix in ("projection", "graph", "packs/pv"):
            with self.subTest(suffix=suffix):
                imports = [f"github.com/Project-Helianthus/helianthus-semreg/semreg/v1/{suffix}"]
                self.assertOneFinding(self.fixture(imports), "driver/sample", imports[0])

    def test_driver_leaf_with_contract_passes(self):
        imports = list(SEMREG_CONTRACT_PATHS) + ["github.com/Project-Helianthus/helianthus-ebusgo/protocol"]
        self.assertEqual([], self.fixture(imports))

    def test_standard_library_is_accepted(self):
        self.assertEqual([], self.fixture(["net", "net/http"]))

    def test_test_files_are_not_checked(self):
        self.assertEqual([], self.fixture(["net"], ["github.com/Project-Helianthus/helianthus-gateway/graphql"]))


if __name__ == "__main__":
    unittest.main()
