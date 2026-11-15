"""Generator check mode, join rules and the clear-interpretation rule."""
import importlib.util
import json
from pathlib import Path
import shutil
import subprocess
import sys
import tempfile
import unittest

SCRIPT = Path(__file__).parents[1] / "generate_native_points.py"
SPEC = importlib.util.spec_from_file_location("generate_native_points", SCRIPT)
generator = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(generator)
ROOT = Path(__file__).resolve().parents[2]
SCRATCH = ROOT / ".tmp"
SCRATCH.mkdir(exist_ok=True)


def location(**fields):
    base = dict(opcode=0x02, group=0, register=1, selector="", instance_rule="single",
                instance_min=0, instance_max=0, families=[], codec="float32_le",
                declaration="", divisor=1, unit="°C", unresolved=False)
    base.update(fields)
    return base


def run(*arguments):
    return subprocess.run([sys.executable, "-I", str(SCRIPT), *arguments],
                          capture_output=True, text=True, timeout=120)


class CheckModeTests(unittest.TestCase):
    def test_committed_table_matches(self):
        result = run("--check")
        self.assertEqual(result.returncode, 0, result.stdout + result.stderr)

    def test_drift_is_rejected_and_regenerated(self):
        with tempfile.TemporaryDirectory(dir=SCRATCH) as scratch:
            root = Path(scratch)
            for relative in [*generator.SOURCES.values(), generator.OUTPUT]:
                (root / relative).parent.mkdir(parents=True, exist_ok=True)
                shutil.copyfile(ROOT / relative, root / relative)
            output = root / generator.OUTPUT
            output.write_text(output.read_text().replace('"outdoor_temperature"', '"outdoor_temperature_x"'))
            self.assertEqual(run("--check", "--root", str(root)).returncode, 1)
            self.assertEqual(run("--root", str(root)).returncode, 0)
            self.assertEqual(run("--check", "--root", str(root)).returncode, 0)
            self.assertEqual(output.read_text(), (ROOT / generator.OUTPUT).read_text())

    def test_missing_source_fails_closed(self):
        with tempfile.TemporaryDirectory(dir=SCRATCH) as scratch:
            result = run("--check", "--root", scratch)
            self.assertEqual(result.returncode, 2)


class RuleTests(unittest.TestCase):
    def test_withheld_reasons_in_order(self):
        cases = [
            (dict(name_status="provisional", value_class="code", unit="", codec=""), "name_not_final"),
            (dict(name_status="ambiguous"), "name_not_final"),
            (dict(value_class="code", unit="", codec=""), "identifier_or_code"),
            (dict(value_class="identifier"), "identifier_or_code"),
            (dict(value_class="boolean"), "not_numeric"),
            (dict(value_class="unknown"), "not_numeric"),
            (dict(unit="", codec=""), "unit_unknown"),
            (dict(unit="furlong"), "unit_unknown"),
            (dict(unit="bar", codec=""), "pressure_reference_unknown"),
            (dict(codec=""), "codec_unknown"),
            (dict(), ""),
        ]
        for fields, reason in cases:
            point = dict(name_status="final", value_class="numeric", unit="°C", codec="DATA2c")
            point.update(fields)
            with self.subTest(fields=fields):
                self.assertEqual(generator.withheld(point), reason)

    def test_value_classes(self):
        cases = [
            ("installer_access_code", "", "", ""), ("installer_menu_code", "", "", ""),
            ("hydraulic_scheme_code", "", "", "uint16_le"), ("installer_phone_number", "", "", ""),
            ("device_specific_number", "", "", ""), ("zone_name_prefix", "", "", ""),
            ("settings_reset_command", "", "", ""), ("maintenance_date", "", "", ""),
            ("holiday_end_clock_time", "", "", ""), ("heat_pump_status", "hpstatus enum", "", ""),
            ("phase_status_flags", "phase flags", "", ""), ("flame_active", "UCH (0x0F=on, 0xF0=off)", "", ""),
            ("defrost_active", "yesno", "", ""), ("zone_enabled", "", "", ""),
            ("dhw_operating_mode", "", "", ""), ("outdoor_temperature", "", "°C", "float32_le"),
            ("manual_cooling_time", "", "", "uint16_le"), ("fan_1_start_count", "", "count", ""),
            ("language", "", "", ""),
        ]
        expected = ["code", "code", "code", "identifier", "identifier", "text", "command", "date",
                    "time", "enumeration", "flags", "boolean", "boolean", "boolean", "enumeration",
                    "numeric", "numeric", "numeric", "unknown"]
        for (name, declaration, unit, codec), value_class in zip(cases, expected):
            with self.subTest(name=name):
                self.assertEqual(generator.value_class(name, declaration, unit, codec), value_class)

    def test_declarations(self):
        self.assertEqual(generator.parse_declaration("DATA2c"), ("DATA2c", 1))
        self.assertEqual(generator.parse_declaration("DATA2b"), ("DATA2b", 1))
        self.assertEqual(generator.parse_declaration("UIN"), ("UIN", 1))
        self.assertEqual(generator.parse_declaration("UIN / 100"), ("UIN", 100))
        for literal in ("SIN / 10", "UCH", "percent0", "temp", "on/off", "UIN / 0", ""):
            with self.subTest(literal=literal):
                self.assertEqual(generator.parse_declaration(literal), ("", 1))

    def test_selector_bytes(self):
        self.assertEqual(generator.selector_bytes("0x1800"), "18 00")
        self.assertEqual(generator.selector_bytes("54 02 00 0d 0A"), "54 02 00 0D 0A")
        for malformed in ("", "0x", "0x180", "0x18G0"):
            with self.subTest(selector=malformed):
                with self.assertRaises(generator.InputError):
                    generator.selector_bytes(malformed)

    def test_point_key_varies_device_class_and_name(self):
        # Keying only by name merges the first pair; keying only by class merges the second.
        pairs = {
            "device_class": [("boiler", "flow_temperature"), ("heat_pump", "flow_temperature")],
            "name": [("boiler", "flow_temperature"), ("boiler", "fan_speed")],
        }
        for field, keys in pairs.items():
            with self.subTest(field=field):
                grouped = {key: [location(register=index)] for index, key in enumerate(keys)}
                points = generator.assemble(grouped, {key: (0xB5, 0x09) for key in keys})
                self.assertEqual([(p["device_class"], p["name"]) for p in points], sorted(keys))
        merged = generator.assemble({("boiler", "x"): [location(register=1), location(register=2)]},
                                    {("boiler", "x"): (0xB5, 0x09)})
        self.assertEqual(len(merged), 1)
        self.assertEqual(len(merged[0]["locations"]), 2)

    def test_controller_codec_key_varies_opcode_group_and_register(self):
        # Ignoring the opcode, group or register in the interpreter lookup would
        # give the mutated row the codec of the base row.
        base = dict(operation=0x02, group=0x00, register=0x0073)
        mutations = {"base": {}, "opcode": dict(operation=0x06), "group": dict(group=0x01),
                     "register": dict(register=0x0074)}
        rows = []
        for label, change in mutations.items():
            key = dict(base, **change)
            rows.append({"key": key, "canonical_name": label, "unresolved": False})
        with tempfile.TemporaryDirectory(dir=SCRATCH) as scratch:
            root = Path(scratch)
            names = root / generator.SOURCES["controller_names"]
            names.parent.mkdir(parents=True)
            names.write_text(json.dumps(rows))
            profiles = root / generator.SOURCES["profiles"]
            profiles.parent.mkdir(parents=True)
            shutil.copyfile(ROOT / generator.SOURCES["profiles"], profiles)
            grouped, services = {}, {}
            generator.controller_locations(root, grouped, services)
        codecs = {name: (places[0]["codec"], places[0]["unit"], places[0]["families"])
                  for (_, name), places in grouped.items()}
        self.assertEqual(codecs["base"], ("float32_le", "°C", ["controller-700", "controller-720"]))
        for label in ("opcode", "group", "register"):
            with self.subTest(field=label):
                self.assertEqual(codecs[label], ("", "", []))

    def test_name_status(self):
        disjoint = [location(families=["controller-700"]), location(register=2, families=["controller-720"])]
        overlapping = [location(families=["controller-700"]), location(register=2, families=["controller-700"])]
        unknown = [location(), location(register=2, families=["controller-720"])]
        self.assertEqual(generator.name_status(disjoint), "final")
        self.assertEqual(generator.name_status(overlapping), "ambiguous")
        self.assertEqual(generator.name_status(unknown), "ambiguous")
        self.assertEqual(generator.name_status([location(unresolved=True)]), "provisional")
        self.assertEqual(generator.name_status(disjoint + [location(register=3, unresolved=True)]), "provisional")
        self.assertEqual(generator.name_status([location()]), "final")

    def test_conflicting_locations_lose_codec_and_unit(self):
        point = generator.assemble({("controller", "x_temperature"): [
            location(families=["controller-700"]),
            location(register=2, families=["controller-720"], codec="uint16_le", unit="min")]},
            {("controller", "x_temperature"): (0xB5, 0x24)})[0]
        self.assertEqual((point["codec"], point["unit"], point["divisor"]), ("", "", 1))
        self.assertEqual(point["families"], ["controller-700", "controller-720"])
        self.assertEqual(point["withheld"], "unit_unknown")


if __name__ == "__main__":
    unittest.main()
