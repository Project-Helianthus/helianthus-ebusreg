"""Synthetic boundary, identity mutation, parser and fail-closed coverage."""
import importlib.util
import json
from pathlib import Path
import tempfile
import unittest

SPEC = importlib.util.spec_from_file_location("size_ratchet", Path(__file__).parents[1] / "check_size_ratchet.py")
guard = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(guard)
SCRATCH = Path(__file__).resolve().parents[2] / ".tmp"
SCRATCH.mkdir(exist_ok=True)


def row(metric="function_lines", value=151, path="sample.go", symbol="Example"):
    return dict(path=path, symbol=symbol, metric=metric, value=value)


class RatchetTests(unittest.TestCase):
    def test_every_threshold(self):
        for kind, budgets in guard.BUDGETS.items():
            for metric, (soft, hard) in budgets.items():
                for boundary in (soft, hard):
                    for delta in (-1, 0, 1):
                        with self.subTest(kind=kind, metric=metric, boundary=boundary, delta=delta):
                            item = row(metric, boundary + delta,
                                       "sample_test.go" if kind == "test" else "sample.go",
                                       "Example" if metric == "function_lines" else "<file>")
                            reports, failures = guard.evaluate([item], [])
                            self.assertEqual(bool(reports), item["value"] > soft)
                            self.assertEqual(bool(failures), item["value"] > hard)

    def test_each_identity_field(self):
        # Removing path, symbol, metric or value from admission would admit one
        # of these independent mutations of existing hard debt.
        original = row(value=1601)
        mutations = {"path": "renamed.go", "symbol": "Renamed",
                     "metric": "file_lines", "value": 1602}
        for field, value in mutations.items():
            with self.subTest(field=field):
                changed = dict(original, **{field: value})
                self.assertEqual(guard.evaluate([changed], [original])[1], [changed])
        self.assertFalse(guard.evaluate([original], [original])[1])

    def test_deleted_reduced_and_renamed_debt(self):
        debt = row(value=200)
        for current in ([], [row(value=199)], [row(value=150)], [row(value=1)]):
            self.assertFalse(guard.evaluate(current, [debt])[1])
        self.assertTrue(guard.evaluate([row(value=200, path="renamed.go")], [debt])[1])
        self.assertTrue(guard.evaluate([row(value=200, symbol="Renamed")], [debt])[1])

    def load(self, text):
        with tempfile.TemporaryDirectory(dir=SCRATCH) as work:
            path = Path(work) / "baseline.json"
            path.write_text(text)
            return guard.load_baseline(path)

    def test_malformed_baselines(self):
        valid = {"version": 1, "debt": [row()]}
        self.assertEqual(self.load(json.dumps(valid)), valid["debt"])
        bad = [None, [], {}, {"version": True, "debt": []},
               {"version": 2, "debt": []}, {"version": 1, "debt": {}},
               dict(valid, extra=True), dict(valid, debt=[row(), row()])]
        for patch in ({"path": "*.go"}, {"path": "../sample.go"}, {"path": "/sample.go"},
                      {"path": "a//sample.go"}, {"path": "sample.txt"}, {"symbol": "*"},
                      {"symbol": ""}, {"symbol": "<file>"}, {"metric": "unknown"},
                      {"value": True}, {"value": "151"}, {"value": 150},
                      {"value": -1}, {"extra": 1}):
            bad.append(dict(valid, debt=[dict(row(), **patch)]))
        for value in bad:
            with self.subTest(value=value), self.assertRaises(ValueError):
                self.load(json.dumps(value))
        for text in ('{', '{"version":1,"version":1,"debt":[]}'):
            with self.assertRaises(ValueError):
                self.load(text)

    def test_parser_methods_literals_tags_platforms_and_spans(self):
        with tempfile.TemporaryDirectory(dir=SCRATCH) as work:
            root = Path(work)
            (root / "sample_windows.go").write_text(
                '//go:build windows && unavailable\n\npackage sample\n'
                'type Device[T any] struct{}\n'
                'func (d *Device[T]) Update() {\n'
                ' f := func() {\n'
                '  g := func() {}\n'
                '  _ = g\n'
                ' }\n'
                ' _ = f\n'
                '}\n'
                'var Hook = func() {}\n')
            (root / "sample_test.go").write_text('package sample\nfunc TestExample() {}')
            rows, _ = guard.measure(root)
            values = {(r["path"], r["symbol"], r["metric"]): r["value"] for r in rows}
            self.assertEqual(values[("sample_windows.go", "<file>", "file_lines")], 12)
            self.assertEqual(values[("sample_windows.go", "<file>", "functions_per_file")], 4)
            self.assertEqual(values[("sample_windows.go", "(*Device[T]).Update", "function_lines")], 7)
            self.assertEqual(values[("sample_windows.go", "(*Device[T]).Update/literal[1]", "function_lines")], 4)
            self.assertEqual(values[("sample_windows.go", "(*Device[T]).Update/literal[1]/literal[1]", "function_lines")], 1)
            self.assertEqual(values[("sample_windows.go", "Hook/value[1]/literal[1]", "function_lines")], 1)
            self.assertEqual(values[("sample_test.go", "<file>", "file_lines")], 2)
            (root / "sample_windows.go").write_text('\n' + (root / "sample_windows.go").read_text())
            moved, _ = guard.measure(root)
            self.assertEqual({guard.key(r) for r in rows}, {guard.key(r) for r in moved})

    def test_parse_failure_even_in_inactive_file(self):
        with tempfile.TemporaryDirectory(dir=SCRATCH) as work:
            root = Path(work)
            (root / "broken_plan9.go").write_text('//go:build never\n\npackage sample\nfunc broken( {')
            with self.assertRaisesRegex(ValueError, "Go parser failed"):
                guard.measure(root)

    def test_exclusions_and_generated_without_verified_owner(self):
        with tempfile.TemporaryDirectory(dir=SCRATCH) as work:
            root = Path(work)
            for directory in guard.EXCLUDED:
                (root / directory).mkdir()
                (root / directory / "bad.go").write_text('invalid')
            (root / "generated.go").write_text('// Code generated by unknown. DO NOT EDIT.\npackage sample\n')
            (root / "ordinary").mkdir()
            (root / "ordinary/plain.go").write_text('package sample\n')
            paths, excluded = guard.sources(root)
            self.assertEqual(paths, ["generated.go", "ordinary/plain.go"])
            self.assertEqual(dict(excluded), guard.EXCLUDED)
            rows, _ = guard.measure(root)
            self.assertEqual(sum(r["metric"] == "file_lines" for r in rows), 2)

    def test_empty_tree_is_unverified(self):
        with tempfile.TemporaryDirectory(dir=SCRATCH) as work:
            with self.assertRaisesRegex(ValueError, "no Go files"):
                guard.measure(Path(work))

    def test_symlink_is_fatal(self):
        with tempfile.TemporaryDirectory(dir=SCRATCH) as work:
            root = Path(work)
            (root / "real.go").write_text('package sample\n')
            (root / "linked.go").symlink_to(root / "real.go")
            with self.assertRaisesRegex(ValueError, "symlink source"):
                guard.sources(root)

    def test_repeated_init_and_local_variable_literal_identity(self):
        with tempfile.TemporaryDirectory(dir=SCRATCH) as work:
            root = Path(work)
            (root / "sample.go").write_text(
                'package sample\n'
                'func init() { var callback = func() {}; _ = callback }\n'
                'func init() { var callback = func() {}; _ = callback }\n'
                'func Other() { var callback = func() {}; _ = callback }\n')
            rows, _ = guard.measure(root)
            identities = [guard.key(r) for r in rows]
            self.assertEqual(len(identities), len(set(identities)))
            symbols = {r["symbol"] for r in rows}
            self.assertTrue({"init[1]/literal[1]", "init[2]/literal[1]",
                             "Other/literal[1]"} <= symbols)
