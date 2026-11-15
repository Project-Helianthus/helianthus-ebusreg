#!/usr/bin/env python3
"""Generate vaillant/nativepoints/points_data.go from the repository's tables.

Inputs, all read from the repository:
  * the controller register name annotations and the register-service symbolic
    annotations (canonical names, resolution flags, units);
  * the boiler and heat-pump register-service codec dictionaries (codec
    declarations, joined by device class and selector);
  * the controller discovery profiles (instance ranges per group and opcode);
  * the codec and declared unit of each controller register that a
    vaillant/system scalar interpreter decodes, per controller family, listed
    in INTERPRETED below; points_test.go probes those interpreters for every
    register and fails when this list drifts from them.

One row is emitted per (device class, canonical name). Annotations without a
canonical name and the two operation-level annotations of class "any" carry no
point. A name is final only when no location marks it unresolved and, when it
has several locations, every location names its controller families and no two
locations share one; otherwise it is provisional or ambiguous.

A row is exportable only when its name is final, its value class is numeric,
its unit converts to SI, its codec is one the decoder implements and it is not
a pressure without a documented reference. Otherwise the first failing rule in
WITHHELD_ORDER is its withheld reason. The evidence grade is Hypothesis on
every source row and is therefore not a gate.

Usage: generate_native_points.py [--check] [--root PATH]
Exit status: 0 written or up to date, 1 the committed table differs, 2 input error.
"""
import argparse
import json
from pathlib import Path
import re
import sys

ROOT = Path(__file__).resolve().parents[1]
OUTPUT = "vaillant/nativepoints/points_data.go"
SOURCES = {
    "controller_names": "catalog/vaillant/data/b524_name_annotations.json",
    "symbolic": "catalog/vaillant/data/b509_symbolic_annotations.json",
    "boiler_codecs": "catalog/vaillant/data/b509_boiler_codec_dictionary.json",
    "heat_pump_codecs": "catalog/vaillant/data/b509_heatpump_codec_dictionary.json",
    "profiles": "vaillant/system/b524_profile.go",
}
CONTROLLER_SERVICE = (0xB5, 0x24)
REGISTER_SERVICE = (0xB5, 0x09)
DICTIONARY_CLASSES = {"boiler": "boiler", "heat-pump": "heat_pump"}

F32, U16 = "float32_le", "uint16_le"
C700, C720 = "controller-700", "controller-720"
SYSTEM_CELSIUS = (0x0002, 0x0022, 0x0023, 0x0026, 0x0038, 0x0046, 0x004B, 0x0067,
                  0x0068, 0x0073, 0x0095, 0x009D, 0x009E, 0x009F, 0x00A0)
INTERPRETED = [
    ((C720,), 0x00, SYSTEM_CELSIUS + (0x00EF,), F32, "°C"),
    ((C700,), 0x00, SYSTEM_CELSIUS, F32, "°C"),
    ((C720,), 0x00, (0x0004, 0x0017, 0x0018), U16, "min"),
    ((C700,), 0x00, (0x0017,), U16, "min"),
    ((C700, C720), 0x00, (0x0029,), F32, "K"),
    ((C700, C720), 0x00, (0x0039,), F32, "bar"),
    ((C700, C720), 0x00, (0x002F, 0x0036), U16, ""),
    ((C700,), 0x00, (0x0006, 0x0042, 0x0043), U16, ""),
    ((C700, C720), 0x01, (0x0004, 0x0005), F32, "°C"),
    ((C720,), 0x09, (0x000D, 0x000E), U16, ""),
    ((C700, C720), 0x02, (0x0007, 0x0008, 0x0010, 0x0011, 0x0012, 0x0014, 0x001F), F32, "°C"),
    ((C700, C720), 0x02, (0x000A, 0x000B), F32, "K"),
    ((C700, C720), 0x02, (0x000F,), F32, ""),
    ((C720,), 0x02, (0x001C,), F32, ""),
    ((C720,), 0x02, (0x001A,), F32, "%"),
    ((C700, C720), 0x03, (0x0002, 0x0005, 0x0009, 0x000F, 0x0014, 0x0015), F32, "°C"),
    ((C720,), 0x03, (0x0023,), F32, "°C"),
    ((C700, C720), 0x03, (0x0028,), F32, "%"),
]

DECODED_CODECS = {F32, U16, "DATA2c", "DATA2b", "UIN"}
SI_UNITS = {"°C", "K", "bar", "kWh", "kW", "W", "h", "min", "s", "%", "count", "steps",
            "L/min", "L/h", "rpm"}
PRESSURE_UNITS = {"bar"}
WITHHELD_ORDER = ("name_not_final", "identifier_or_code", "not_numeric", "unit_unknown",
                  "pressure_reference_unknown", "codec_unknown")
IDENTIFIER_MARKERS = ("phone_number", "serial_number", "device_specific_number")
TEXT_SUFFIXES = ("_name", "_name_prefix", "_name_suffix")
BOOLEAN_SUFFIXES = ("_active", "_enabled", "_ok", "_available")
ENUMERATION_SUFFIXES = ("_mode", "_type", "_state", "_status", "_weekday")
QUANTITY_TOKENS = {
    "temperature", "pressure", "count", "total", "yield", "consumption", "rate", "power",
    "speed", "humidity", "time", "duration", "interval", "offset", "hysteresis", "limit",
    "setpoint", "factor", "difference", "level", "position", "current", "movement",
    "superheat", "subcooling", "price", "integral", "ratio",
}
PROFILE = re.compile(r"\{Group: (0x[0-9A-Fa-f]{2}), Opcode: B524Opcode(Local|Remote), "
                     r"InstanceMax: (0x[0-9A-Fa-f]{2}), RegisterMax: 0x[0-9A-Fa-f]{4}\}")
DECLARATION = re.compile(r"(DATA2c|DATA2b|UIN)(?: / ([1-9][0-9]*))?")
PLAIN_TEXT = re.compile(r"[^\"\\\x00-\x1f]*")


class InputError(ValueError):
    pass


def parse_declaration(text):
    match = DECLARATION.fullmatch(text)
    if not match:
        return "", 1
    return match.group(1), int(match.group(2) or 1)


def value_class(name, declaration, unit, codec):
    lowered = declaration.lower()
    if name.endswith("_code") or "_code_" in name:
        return "code"
    if any(marker in name for marker in IDENTIFIER_MARKERS):
        return "identifier"
    if name.endswith(TEXT_SUFFIXES):
        return "text"
    if name.endswith("_command"):
        return "command"
    if name == "date" or name.endswith("_date"):
        return "date"
    if name == "clock_time" or name.endswith("_clock_time"):
        return "time"
    if "enum" in lowered:
        return "enumeration"
    if "flags" in lowered:
        return "flags"
    if lowered in ("on/off", "yesno") or "=on" in lowered or name.endswith(BOOLEAN_SUFFIXES):
        return "boolean"
    if name.endswith(ENUMERATION_SUFFIXES):
        return "enumeration"
    if codec in DECODED_CODECS or unit in SI_UNITS or QUANTITY_TOKENS & set(name.split("_")):
        return "numeric"
    return "unknown"


def withheld(point):
    failing = {
        "name_not_final": point["name_status"] != "final",
        "identifier_or_code": point["value_class"] in ("code", "identifier"),
        "not_numeric": point["value_class"] != "numeric",
        "unit_unknown": point["unit"] not in SI_UNITS,
        "pressure_reference_unknown": point["unit"] in PRESSURE_UNITS,
        "codec_unknown": point["codec"] not in DECODED_CODECS,
    }
    return next((reason for reason in WITHHELD_ORDER if failing[reason]), "")


def name_status(locations):
    if any(place["unresolved"] for place in locations):
        return "provisional"
    if len(locations) > 1:
        seen = set()
        for place in locations:
            families = set(place["families"])
            if not families or families & seen:
                return "ambiguous"
            seen |= families
    return "final"


def agreed(locations, field, default):
    values = {json.dumps(place[field]) for place in locations}
    return locations[0][field] if len(values) == 1 else default


def assemble(grouped, services):
    points = []
    for (device_class, name), locations in sorted(grouped.items()):
        locations = sorted(locations, key=lambda p: (p["opcode"], p["group"], p["register"], p["selector"]))
        consistent = len({(p["codec"], p["divisor"], p["unit"]) for p in locations}) == 1
        point = dict(
            device_class=device_class, name=name, service=services[(device_class, name)],
            name_status=name_status(locations), locations=locations,
            codec=locations[0]["codec"] if consistent else "",
            divisor=locations[0]["divisor"] if consistent else 1,
            unit=locations[0]["unit"] if consistent else "",
            declaration=agreed(locations, "declaration", ""),
            families=sorted({f for place in locations for f in place["families"]}),
        )
        point["value_class"] = value_class(name, point["declaration"], point["unit"], point["codec"])
        point["withheld"] = withheld(point)
        points.append(point)
    return points


def load_json(root, key):
    path = root / SOURCES[key]
    try:
        data = json.loads(path.read_text(encoding="utf-8"))
    except (OSError, ValueError) as error:
        raise InputError(f"{SOURCES[key]}: {error}") from error
    if not isinstance(data, list) or not data:
        raise InputError(f"{SOURCES[key]}: expected a non-empty array")
    return data


def instance_rules(root):
    try:
        text = (root / SOURCES["profiles"]).read_text(encoding="utf-8")
    except OSError as error:
        raise InputError(f"{SOURCES['profiles']}: {error}") from error
    rules = {}
    for group, opcode, instance_max in PROFILE.findall(text):
        key = (int(group, 16), 0x02 if opcode == "Local" else 0x06)
        if key in rules:
            raise InputError("duplicate discovery profile")
        maximum = int(instance_max, 16)
        rules[key] = ("single", 0, 0) if maximum == 0 else ("indexed", 0, maximum)
    if not rules:
        raise InputError("no discovery profile found")
    return rules


def interpreted():
    table = {}
    for families, group, registers, codec, unit in INTERPRETED:
        for register in registers:
            for family in families:
                slot = table.setdefault((group, register), {})
                if family in slot:
                    raise InputError("duplicate interpreter declaration")
                slot[family] = (codec, unit)
    return table


def controller_locations(root, grouped, services):
    rules, codecs = instance_rules(root), interpreted()
    for row in load_json(root, "controller_names"):
        key, name = row["key"], row["canonical_name"]
        if name is None:
            continue
        opcode, group, register = key["operation"], key["group"], key["register"]
        declared = codecs.get((group, register), {}) if opcode == 0x02 else {}
        pairs = set(declared.values())
        if len(pairs) > 1:
            raise InputError("controller families disagree on one register")
        codec, unit = pairs.pop() if pairs else ("", "")
        rule, low, high = rules.get((group, opcode), ("unknown", 0, 0))
        place = dict(opcode=opcode, group=group, register=register, selector="",
                     instance_rule=rule, instance_min=low, instance_max=high,
                     families=sorted(declared), codec=codec, declaration="", divisor=1,
                     unit=unit, unresolved=row["unresolved"] is not False)
        grouped.setdefault(("controller", name), []).append(place)
        services[("controller", name)] = CONTROLLER_SERVICE


def register_service_locations(root, grouped, services):
    declarations = {}
    for source in ("boiler_codecs", "heat_pump_codecs"):
        for entry in load_json(root, source):
            key = entry["key"]
            identity = (DICTIONARY_CLASSES[key["family"]], key["operation"],
                        " ".join("%02X" % byte for byte in key["selector"]))
            if identity in declarations:
                raise InputError("duplicate codec dictionary key")
            declarations[identity] = entry["codec_declaration"]
    joined = set()
    for row in load_json(root, "symbolic"):
        if row["class"] == "any":
            continue
        selector = selector_bytes(row["selector"])
        identity = (row["class"], row["operation"], selector)
        declaration = declarations.get(identity, "")
        if declaration:
            joined.add(identity)
        codec, divisor = parse_declaration(declaration)
        place = dict(opcode=row["operation"], group=0, register=0, selector=selector,
                     instance_rule="single", instance_min=0, instance_max=0, families=[],
                     codec=codec, declaration=declaration, divisor=divisor,
                     unit=row.get("unit") or "", unresolved=row["unresolved"] is not False)
        grouped.setdefault((row["class"], row["name"]), []).append(place)
        services[(row["class"], row["name"])] = REGISTER_SERVICE
    if joined != set(declarations):
        raise InputError("codec dictionary entry without a symbolic annotation")


def selector_bytes(text):
    digits = text.removeprefix("0x").replace(" ", "").upper()
    if not digits or len(digits) % 2 or not re.fullmatch(r"[0-9A-F]+", digits):
        raise InputError(f"malformed selector {text!r}")
    return " ".join(digits[index:index + 2] for index in range(0, len(digits), 2))


def quote(text):
    if not PLAIN_TEXT.fullmatch(text):
        raise InputError(f"unexpected character in {text!r}")
    return f'"{text}"'


def families_literal(families):
    return "[]string{" + ", ".join(quote(f) for f in families) + "}"


def location_literal(place):
    fields = [f"Opcode: 0x{place['opcode']:02x}"]
    if place["selector"]:
        fields.append(f"Selector: {quote(place['selector'])}")
    else:
        fields.append(f"Group: 0x{place['group']:02x}, Register: 0x{place['register']:04x}")
    fields.append(f"InstanceRule: {quote(place['instance_rule'])}")
    if place["instance_rule"] == "indexed":
        fields.append(f"InstanceMin: {place['instance_min']}, InstanceMax: {place['instance_max']}")
    if place["families"]:
        fields.append(f"ControllerFamilies: {families_literal(place['families'])}")
    if place["codec"]:
        fields.append(f"Codec: {quote(place['codec'])}")
    if place["declaration"]:
        fields.append(f"CodecDeclaration: {quote(place['declaration'])}")
    fields.append(f"ScaleDivisor: {place['divisor']}")
    if place["unit"]:
        fields.append(f"NativeUnit: {quote(place['unit'])}")
    return "{" + ", ".join(fields) + "}"


def point_literal(point):
    primary, secondary = point["service"]
    fields = [f"DeviceClass: {quote(point['device_class'])}", f"Name: {quote(point['name'])}",
              f"NameStatus: {quote(point['name_status'])}",
              f"Primary: 0x{primary:02x}, Secondary: 0x{secondary:02x}",
              f"ValueClass: {quote(point['value_class'])}"]
    if point["codec"]:
        fields.append(f"Codec: {quote(point['codec'])}")
    if point["declaration"]:
        fields.append(f"CodecDeclaration: {quote(point['declaration'])}")
    fields.append(f"ScaleDivisor: {point['divisor']}")
    if point["unit"]:
        fields.append(f"NativeUnit: {quote(point['unit'])}")
    if point["families"]:
        fields.append(f"ControllerFamilies: {families_literal(point['families'])}")
    fields.append("Locations: []Location{" + ", ".join(location_literal(p) for p in point["locations"]) + "}")
    if point["withheld"]:
        fields.append(f"Withheld: {quote(point['withheld'])}")
    else:
        fields.append("Exportable: true")
    return "\t{" + ", ".join(fields) + "},"


def build(root):
    grouped, services = {}, {}
    try:
        controller_locations(root, grouped, services)
        register_service_locations(root, grouped, services)
    except (KeyError, TypeError, AttributeError) as error:
        raise InputError(f"malformed source row: {error!r}") from error
    points = assemble(grouped, services)
    lines = ["// Code generated by scripts/generate_native_points.py; DO NOT EDIT.", "",
             "package nativepoints", "", "var generatedPoints = []Point{"]
    lines += [point_literal(point) for point in points]
    lines += ["}", ""]
    return "\n".join(lines)


def main(argv=None):
    parser = argparse.ArgumentParser(description=__doc__.splitlines()[0])
    parser.add_argument("--check", action="store_true", help="compare instead of writing")
    parser.add_argument("--root", type=Path, default=ROOT)
    args = parser.parse_args(argv)
    root = args.root.resolve()
    try:
        content = build(root)
    except InputError as error:
        print(f"native points: {error}", file=sys.stderr)
        return 2
    target = root / OUTPUT
    if args.check:
        try:
            current = target.read_text(encoding="utf-8")
        except OSError:
            current = None
        if current != content:
            print(f"native points: {OUTPUT} differs from the generated table", file=sys.stderr)
            return 1
        print(f"native points: {OUTPUT} is up to date")
        return 0
    target.parent.mkdir(parents=True, exist_ok=True)
    target.write_text(content, encoding="utf-8")
    print(f"native points: wrote {OUTPUT}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
