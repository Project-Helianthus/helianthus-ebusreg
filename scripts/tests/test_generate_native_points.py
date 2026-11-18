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
                declaration="", divisor=1, unit="°C", enumeration=[], kind="", unresolved=False)
    base.update(fields)
    return base


def matrix_row(family, register, value_type, widths, unit="", enumeration="", group=0x00):
    return dict(family=family, opcode=0x02, group=group, register=register, name="synthetic",
                reply_widths=widths, value_type=value_type, unit=unit, enumeration=enumeration)


def names_for(base, mutations, **columns):
    rows = []
    for label, change in mutations.items():
        rows.append({"key": dict(base, **change), "canonical_name": label, "unresolved": False, **columns})
    return rows


def controller_join(names, matrix):
    with tempfile.TemporaryDirectory(dir=SCRATCH) as scratch:
        root = Path(scratch)
        for key, rows in (("controller_names", names), ("register_matrix", matrix)):
            path = root / generator.SOURCES[key]
            path.parent.mkdir(parents=True, exist_ok=True)
            path.write_text(json.dumps(rows), encoding="utf-8")
        profiles = root / generator.SOURCES["profiles"]
        profiles.parent.mkdir(parents=True, exist_ok=True)
        shutil.copyfile(ROOT / generator.SOURCES["profiles"], profiles)
        grouped, services = {}, {}
        generator.controller_locations(root, grouped, services)
    return grouped


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
        matrix = [matrix_row(family, 0x0073, "float32 LE", [4], unit="°C")
                  for family in ("controller-700", "controller-720")]
        grouped = controller_join(names_for(base, mutations), matrix)
        codecs = {name: (places[0]["codec"], places[0]["unit"], places[0]["families"])
                  for (_, name), places in grouped.items()}
        self.assertEqual(codecs["base"], ("float32_le", "°C", ["controller-700", "controller-720"]))
        for label in ("opcode", "group", "register"):
            with self.subTest(field=label):
                self.assertEqual(codecs[label], ("", "", []))

    def test_matrix_key_varies_opcode_group_register_and_family(self):
        # Ignoring the opcode, group or register of a matrix row would give the
        # mutated annotation the base layout; ignoring the family would merge
        # the two family rows of the base register into one location.
        base = dict(operation=0x02, group=0x03, register=0x0200)
        mutations = {"base": {}, "opcode": dict(operation=0x06), "group": dict(group=0x01),
                     "register": dict(register=0x0201)}
        matrix = [matrix_row("controller-720", 0x0200, "uint8", [1], enumeration="0=Off; 1=On", group=0x03),
                  matrix_row("controller-700", 0x0200, "uint16 LE", [2], unit="min", group=0x03)]
        grouped = controller_join(names_for(base, mutations), matrix)
        layouts = {name: [(p["codec"], p["unit"], p["families"], p["enumeration"]) for p in places]
                   for (_, name), places in grouped.items()}
        self.assertEqual(layouts["base"], [
            ("uint16_le", "min", ["controller-700"], []),
            ("uint8", "", ["controller-720"], [{"value": 0, "label": "Off"}, {"value": 1, "label": "On"}])])
        for label in ("opcode", "group", "register"):
            with self.subTest(field=label):
                self.assertEqual(layouts[label], [("", "", [], [])])
        slots = {name: (places[0]["instance_rule"], places[0]["instance_min"], places[0]["instance_max"])
                 for (_, name), places in grouped.items()}
        self.assertEqual(slots["opcode"], ("indexed", 1, 8))
        self.assertEqual(slots["base"], ("indexed", 0, 10))
        groups = {"slot group": dict(operation=0x06, group=0x09), "group above": dict(operation=0x06, group=0x10),
                  "singleton": dict(group=0x00)}
        grouped = controller_join(names_for(base, groups), matrix)
        slots = {name: (places[0]["instance_rule"], places[0]["instance_min"], places[0]["instance_max"])
                 for (_, name), places in grouped.items()}
        self.assertEqual(slots, {"slot group": ("indexed", 1, 8), "group above": ("unknown", 0, 0),
                                 "singleton": ("single", 0, 0)})

    def test_interpreter_disagreement_fails_generation(self):
        base = dict(operation=0x02, group=0x00, register=0x0073)
        for family in ("controller-700", "controller-720"):
            with self.subTest(family=family):
                matrix = [matrix_row(other, 0x0073, "float32 LE", [4], unit="°C")
                          for other in ("controller-700", "controller-720") if other != family]
                matrix.append(matrix_row(family, 0x0073, "uint16 LE", [2], unit="°C"))
                with self.assertRaisesRegex(generator.InputError, family):
                    controller_join(names_for(base, {"base": {}}), matrix)
        with self.assertRaisesRegex(generator.InputError, "interpreter"):
            controller_join(names_for(base, {"base": {}}, wire="u16", decode="°C"),
                            [matrix_row("controller-720", 0x0073, "float32 LE", [4], unit="°C")])

    def test_matrix_rows_fail_closed(self):
        cases = {
            "width": matrix_row("controller-720", 0x0200, "float32 LE", [2]),
            "type": matrix_row("controller-720", 0x0200, "float64 LE", [8]),
            "unit": matrix_row("controller-720", 0x0200, "uint8", [1], unit="furlong"),
            "family": matrix_row("controller-999", 0x0200, "uint8", [1]),
            "calendar": matrix_row("controller-720", 0x0200, "3 x uint8", [3], enumeration="0=Off"),
            "labels": matrix_row("controller-720", 0x0200, "uint8", [1], enumeration="0=Off; 0=On"),
            "date width": matrix_row("controller-720", 0x0200, "uint16 LE", [2],
                                     enumeration="date: day 1..31; month 1..12; year 0..99"),
        }
        base = dict(operation=0x02, group=0x00, register=0x0200)
        for label, row in cases.items():
            with self.subTest(case=label):
                with self.assertRaises(generator.InputError):
                    controller_join(names_for(base, {"base": {}}), [row])
        duplicate = [matrix_row("controller-720", 0x0200, "uint8", [1])] * 2
        with self.assertRaises(generator.InputError):
            controller_join(names_for(base, {"base": {}}), duplicate)

    def test_matrix_calendar_and_labels(self):
        self.assertEqual(generator.enumeration_labels("candidate labels: 0=regulator"),
                         [{"value": 0, "label": "regulator"}])
        self.assertEqual(generator.enumeration_labels("1=b; -1=notBound; 0=a"),
                         [{"value": -1, "label": "notBound"}, {"value": 0, "label": "a"},
                          {"value": 1, "label": "b"}])
        base = dict(operation=0x02, group=0x00, register=0x0200)
        for text, codec, kind in (("date: day 1..31; month 1..12; year 0..99", "day_month_year", "date"),
                                  ("time: hour; minute; second", "hour_minute_second", "time")):
            with self.subTest(codec=codec):
                grouped = controller_join(names_for(base, {"base": {}}),
                                          [matrix_row("controller-720", 0x0200, "3 x uint8", [3], enumeration=text)])
                place = grouped[("controller", "base")][0]
                self.assertEqual((place["codec"], place["kind"], place["enumeration"]), (codec, kind, []))

    def test_map_columns(self):
        cases = [
            (dict(wire="u8", decode="bool"), ("uint8", "", 1, "boolean")),
            (dict(wire="u16", decode="enum"), ("uint16_le", "", 1, "enumeration")),
            (dict(wire="u16", decode="state"), ("uint16_le", "", 1, "enumeration")),
            (dict(wire="f32", decode="°C"), ("float32_le", "°C", 1, "")),
            (dict(wire="u8", decode="°C/2"), ("uint8", "°C", 2, "")),
            (dict(wire="u32", decode="hrs"), ("uint32_le", "h", 1, "")),
            (dict(wire="i8", decode="—"), ("int8", "", 1, "")),
            (dict(wire="string", decode="text"), ("text", "", 1, "text")),
            (dict(wire="date", decode="date"), ("", "", 1, "date")),
            (dict(wire="time", decode="version"), ("", "", 1, "")),
            (dict(wire="bytes", decode="raw"), ("", "", 1, "")),
            (dict(), ("", "", 1, "")),
        ]
        for columns, expected in cases:
            with self.subTest(columns=columns):
                place = generator.map_layout(columns)
                self.assertEqual((place["codec"], place["unit"], place["divisor"], place["kind"]), expected)
        for columns in (dict(wire="f64", decode="°C"), dict(wire="u8", decode="furlong"), dict(decode="°C")):
            with self.subTest(columns=columns):
                with self.assertRaises(generator.InputError):
                    generator.map_layout(columns)

    def test_stated_kind_precedes_the_name_rules(self):
        self.assertEqual(generator.value_class("device_connected", "", "", "uint8", "boolean"), "boolean")
        self.assertEqual(generator.value_class("device_connected", "", "", "uint8"), "unknown")
        self.assertEqual(generator.value_class("hydraulic_scheme_code", "", "", "uint16_le", "enumeration"), "code")
        self.assertEqual(generator.value_class("serial_number_prefix", "", "", "text", "text"), "identifier")
        self.assertEqual(generator.value_class("zone_valve_position", "", "", "uint16_le", "boolean"), "boolean")

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

    def test_labels_and_kinds_need_agreeing_locations(self):
        def assembled(*locations):
            key = ("controller", "x_mode")
            return generator.assemble({key: list(locations)}, {key: (0xB5, 0x24)})[0]
        off = [{"value": 0, "label": "Off"}]
        on = [{"value": 0, "label": "On"}]
        same = assembled(location(codec="uint8", unit="", enumeration=off, kind="enumeration"),
                         location(register=2, codec="uint8", unit="", enumeration=off))
        self.assertEqual((same["codec"], same["enumeration"], same["value_class"]),
                         ("uint8", off, "enumeration"))
        split = assembled(location(codec="uint8", unit="", enumeration=off),
                          location(register=2, codec="uint8", unit="", enumeration=on))
        self.assertEqual((split["codec"], split["enumeration"]), ("uint8", []))
        conflicting = assembled(location(codec="uint8", unit="", kind="boolean"),
                                location(register=2, codec="uint8", unit="", kind="date"))
        self.assertEqual(conflicting["value_class"], "enumeration")
        recoded = assembled(location(codec="uint8", unit="", enumeration=off),
                            location(register=2, codec="uint16_le", unit="", enumeration=off))
        self.assertEqual((recoded["codec"], recoded["enumeration"]), ("", []))


if __name__ == "__main__":
    unittest.main()
