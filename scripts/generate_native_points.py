#!/usr/bin/env python3
"""Generate vaillant/nativepoints/points_data.go from the repository's tables.

Inputs, all read from the repository:
  * the controller register name annotations, which also carry the Wire and
    Decode columns the register map prints for each register, and the
    register-service symbolic annotations (canonical names, resolution flags,
    units);
  * the controller register matrix by family (reply width, value type, unit
    and enumeration of each opcode 0x02 register per controller family);
  * the boiler and heat-pump register-service codec dictionaries (codec
    declarations, joined by device class and selector);
  * the controller discovery profiles (instance ranges of the opcode 0x02
    groups) and the slot interval of the opcode 0x06 groups, OP06_SLOTS below,
    from the register map's group topology;
  * the codec and declared unit of each controller register that a
    vaillant/system scalar interpreter decodes, per controller family, listed
    in INTERPRETED below; points_test.go probes those interpreters for every
    register and fails when this list drifts from them.

A controller location takes its codec, unit and enumeration from the register
matrix row of each controller family that states a value type, and from the
register map's Wire and Decode columns when no family row states one. When the
matrix states a value type for one family only, the other family takes the
map's layout, except on a register the documentation reads on one family only
(ONE_FAMILY below). Families whose rows disagree get one location each. An
INTERPRETED declaration must agree with the documented codec and unit of its
family; a disagreement fails generation. Map wire types whose layout the map
leaves open (date, time, bytes, unknown) and decode columns without a unit or
value class give no codec and no unit; the map's time wire with the version
decode is the three-byte major, minor and patch triple. An enumeration
or Boolean register of the map takes the label list its Values column states
outright (MAP_LABELS below). A location on a register the map reads on one
controller family only names that family even when no matrix row lists it.

One row is emitted per (device class, canonical name). Annotations without a
canonical name and the two operation-level annotations of class "any" carry no
point. A name is final only when no location marks it unresolved and, when it
has several locations, every location names its controller families and no two
locations share one; otherwise it is provisional or ambiguous.

A point's value class follows the decode column of the map (Boolean,
enumeration, date, time or text) and otherwise the name rules below. Each
controller location also carries the value class its own documentation states:
the decode column of the map, the matrix calendar layout, the matrix Boolean
label pair 0=False 1=True, or enumeration for any other label list; it is
empty when the location's documentation states none, so a name rule or a
sibling location never makes DecodeNative treat its integer as a Boolean or
an enumeration code. A row is
exportable only when its name is final, its value class is numeric,
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
    "register_matrix": "catalog/vaillant/data/b524_register_matrix.json",
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

U8, I8, U8_OR_U16, U32 = "uint8", "int8", "uint8_or_uint16_le", "uint32_le"
DAY_MONTH_YEAR, HOUR_MINUTE_SECOND, TEXT = "day_month_year", "hour_minute_second", "text"
FAMILIES = (C700, C720)
MATRIX_TYPES = {"float32 LE": (F32, [4]), "uint16 LE": (U16, [2]), "uint8": (U8, [1]),
                "int8": (I8, [1]), "uint8 or uint16 LE": (U8_OR_U16, [1, 2]), "3 x uint8": ("", [3])}
MATRIX_UNITS = {"", "°C", "K", "min", "%", "bar"}
MATRIX_CALENDAR = {"date: day 1..31; month 1..12; year 0..99": (DAY_MONTH_YEAR, "date"),
                   "time: hour; minute; second": (HOUR_MINUTE_SECOND, "time")}
LABEL_PREFIX = "candidate labels: "
LABEL = re.compile(r"(-?[0-9]+)=([^;=]+)")
WIRE_CODECS = {"u8": U8, "i8": I8, "u16": U16, "u32": U32, "f32": F32, "string": TEXT,
               "date": "", "time": "", "bytes": "", "unknown": ""}
DECODE_UNITS = {"°C": ("°C", 1), "°C/2": ("°C", 2), "K": ("K", 1), "bar": ("bar", 1),
                "%": ("%", 1), "kWh": ("kWh", 1), "hrs": ("h", 1), "min": ("min", 1),
                "count": ("count", 1), "days": ("d", 1)}
DECODE_KINDS = {"bool": "boolean", "enum": "enumeration", "state": "enumeration",
                "date": "date", "time": "time", "text": "text"}
UNSTATED_DECODES = {"—", "unknown", "raw", "schedule"}
MAJOR_MINOR_PATCH = "major_minor_patch"
OP06_SLOTS = {group: (1, 8) for group in range(0x01, 0x10)}
# Label lists that the register map's Values column states outright for an
# enumeration or Boolean register, verbatim; a reference to the enum reference
# section is not a list and is not copied.
MAP_LABELS = {
    (0x02, 0x00, 0x0016): "1=AUTO 2=DAY 3=SET_BACK",
    (0x02, 0x00, 0x002B): "0=off 1=Mon 2=Tue 3=Wed 4=Thu 5=Fri 6=Sat 7=Sun",
    (0x02, 0x00, 0x0048): "0=standby 1=heating 2=cooling 3=dhw",
    (0x02, 0x02, 0x0003): "0=inactive 1=active 2=extended",
    (0x02, 0x09, 0x0002): "1=TIME_CONTROLLED 2=NORMAL 3=REDUCED",
    (0x06, 0x09, 0x0001): "0=empty 1=paired",
    (0x06, 0x09, 0x0002): "0x15=VRC720 0x35=VR92 0x26=VR71",
    (0x06, 0x0A, 0x0001): "0=empty 1=paired",
    (0x06, 0x0A, 0x0002): "0x15=VRC720 0x35=VR92 0x26=VR71",
    (0x06, 0x0C, 0x0001): "0=not_connected 1=connected",
}
LABELLED_DECODES = ("bool", "enum", "state")
BOOLEAN_LABELS = [{"value": 0, "label": "False"}, {"value": 1, "label": "True"}]
MAP_LABEL = re.compile(r"(0x[0-9A-Fa-f]+|[0-9]+)=(\S+)")
# Opcode 0x02 registers that the register map reads on one controller family
# only: group 0x09 and the continuous-heating set point on the controller-720
# family; the controller-700 ventilation items of group 0x00 and the family
# split of the holiday dates on the controller-700 family.
ONE_FAMILY_GROUPS = {0x09: C720}
ONE_FAMILY = {(0x00, 0x00EF): C720, (0x00, 0x0016): C700, (0x00, 0x0042): C700,
              (0x00, 0x0043): C700, (0x00, 0x00A4): C700, (0x00, 0x0071): C700,
              (0x00, 0x0072): C700}

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


def value_class(name, declaration, unit, codec, kind=""):
    lowered = declaration.lower()
    if name.endswith("_code") or "_code_" in name:
        return "code"
    if any(marker in name for marker in IDENTIFIER_MARKERS):
        return "identifier"
    if name.endswith(TEXT_SUFFIXES):
        return "text"
    if name.endswith("_command"):
        return "command"
    if kind:
        return kind
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
        labelled = consistent and len({json.dumps(p["enumeration"]) for p in locations}) == 1
        kinds = {p["kind"] for p in locations if p["kind"]}
        point = dict(
            device_class=device_class, name=name, service=services[(device_class, name)],
            name_status=name_status(locations), locations=locations,
            codec=locations[0]["codec"] if consistent else "",
            divisor=locations[0]["divisor"] if consistent else 1,
            unit=locations[0]["unit"] if consistent else "",
            enumeration=locations[0]["enumeration"] if labelled else [],
            declaration=agreed(locations, "declaration", ""),
            families=sorted({f for place in locations for f in place["families"]}),
        )
        kind = kinds.pop() if len(kinds) == 1 else ""
        point["value_class"] = value_class(name, point["declaration"], point["unit"], point["codec"], kind)
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


def layout(codec="", unit="", divisor=1, enumeration=(), kind=""):
    return dict(codec=codec, unit=unit, divisor=divisor, enumeration=list(enumeration), kind=kind)


def enumeration_labels(text):
    if not text:
        return []
    labels = {}
    for part in text.removeprefix(LABEL_PREFIX).split("; "):
        match = LABEL.fullmatch(part)
        if not match or int(match.group(1)) in labels:
            raise InputError(f"malformed enumeration {text!r}")
        labels[int(match.group(1))] = match.group(2)
    return [{"value": value, "label": labels[value]} for value in sorted(labels)]


def register_matrix(root):
    table = {}
    for row in load_json(root, "register_matrix"):
        family, value_type, unit = row["family"], row["value_type"], row["unit"]
        key = (row["opcode"], row["group"], row["register"])
        if family not in FAMILIES or family in table.get(key, {}) or unit not in MATRIX_UNITS:
            raise InputError(f"malformed or duplicate register matrix row {key} {family}")
        if not value_type:
            continue
        if value_type not in MATRIX_TYPES or MATRIX_TYPES[value_type][1] != row["reply_widths"]:
            raise InputError(f"register matrix {key} {family}: width {row['reply_widths']} "
                             f"does not fit value type {value_type!r}")
        codec, kind, enumeration = MATRIX_TYPES[value_type][0], "", row["enumeration"]
        if (enumeration in MATRIX_CALENDAR) != (value_type == "3 x uint8"):
            raise InputError(f"register matrix {key} {family}: {value_type!r} with {enumeration!r}")
        if enumeration in MATRIX_CALENDAR:
            codec, kind = MATRIX_CALENDAR[enumeration]
            enumeration = ""
        table.setdefault(key, {})[family] = layout(codec, unit, 1, enumeration_labels(enumeration), kind)
    return table


def map_labels(text):
    labels = {}
    for part in text.split(" "):
        match = MAP_LABEL.fullmatch(part)
        if not match or int(match.group(1), 0) in labels:
            raise InputError(f"malformed register map label list {text!r}")
        labels[int(match.group(1), 0)] = match.group(2)
    return [{"value": value, "label": labels[value]} for value in sorted(labels)]


def map_layout(row, key=None):
    wire, decode = row.get("wire"), row.get("decode")
    labels = MAP_LABELS.get(key)
    if labels is not None and decode not in LABELLED_DECODES:
        raise InputError(f"register map labels on {key} without an enumeration or Boolean decode")
    if wire is None and decode is None:
        return layout()
    if wire not in WIRE_CODECS:
        raise InputError(f"unknown register map wire type {wire!r}")
    if decode == "version":
        if wire != "time":
            raise InputError(f"register map version decode on wire type {wire!r}")
        return layout(MAJOR_MINOR_PATCH, kind="version")
    if decode in DECODE_UNITS:
        unit, divisor = DECODE_UNITS[decode]
        return layout(WIRE_CODECS[wire], unit, divisor)
    if decode in DECODE_KINDS:
        return layout(WIRE_CODECS[wire], enumeration=map_labels(labels) if labels else (),
                      kind=DECODE_KINDS[decode])
    if decode not in UNSTATED_DECODES:
        raise InputError(f"unknown register map decode {decode!r}")
    return layout(WIRE_CODECS[wire])


def one_family(group, register):
    return ONE_FAMILY.get((group, register), ONE_FAMILY_GROUPS.get(group))


def documented_layouts(name, documented, mapped, declared, only=None):
    families = dict(documented)
    if only is not None and set(families) - {only}:
        raise InputError(f"register matrix lists {sorted(families)} for {name}, which the "
                         f"register map reads on {only} only")
    for family, (codec, unit) in sorted(declared.items()):
        expected = families.get(family, mapped)
        if (expected["codec"], expected["unit"]) != (codec, unit):
            raise InputError(f"interpreter of {family} reads {name} as {codec!r} {unit!r}; the "
                             f"documentation states {expected['codec']!r} {expected['unit']!r}")
        families.setdefault(family, mapped)
    if not families:
        return [(mapped, [only] if only is not None else [])]
    if only is None:
        for family in FAMILIES:
            families.setdefault(family, mapped)
    split = {}
    for family, place in sorted(families.items()):
        split.setdefault(json.dumps(place, sort_keys=True), (place, []))[1].append(family)
    return list(split.values())


def location_class(place, mapped):
    kind = place["kind"] or mapped["kind"]
    if kind:
        return kind
    if place["enumeration"] == BOOLEAN_LABELS:
        return "boolean"
    return "enumeration" if place["enumeration"] else ""


def controller_locations(root, grouped, services):
    rules, codecs, matrix = instance_rules(root), interpreted(), register_matrix(root)
    labelled = set()
    for row in load_json(root, "controller_names"):
        key, name = row["key"], row["canonical_name"]
        if name is None:
            continue
        opcode, group, register = key["operation"], key["group"], key["register"]
        labelled.add((opcode, group, register))
        declared = codecs.get((group, register), {}) if opcode == 0x02 else {}
        documented = matrix.get((opcode, group, register), {})
        if opcode == 0x06:
            low, high = OP06_SLOTS.get(group, (0, 0))
            rule = "indexed" if group in OP06_SLOTS else "unknown"
        else:
            rule, low, high = rules.get((group, opcode), ("unknown", 0, 0))
        only = one_family(group, register) if opcode == 0x02 else None
        mapped = map_layout(row, (opcode, group, register))
        for place, families in documented_layouts(name, documented, mapped, declared, only):
            place = dict(place, opcode=opcode, group=group, register=register, selector="",
                         instance_rule=rule, instance_min=low, instance_max=high,
                         families=families, declaration="",
                         value_class=location_class(place, mapped),
                         unresolved=row["unresolved"] is not False)
            grouped.setdefault(("controller", name), []).append(place)
        services[("controller", name)] = CONTROLLER_SERVICE
    return labelled


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
                     unit=row.get("unit") or "", enumeration=[], kind="", value_class="",
                     unresolved=row["unresolved"] is not False)
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


def enumeration_literal(enumeration):
    return "[]EnumerationLabel{" + ", ".join(
        f"{{Value: {entry['value']}, Label: {quote(entry['label'])}}}" for entry in enumeration) + "}"


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
    if place["enumeration"]:
        fields.append(f"Enumeration: {enumeration_literal(place['enumeration'])}")
    if place["value_class"]:
        fields.append(f"ValueClass: {quote(place['value_class'])}")
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
    if point["enumeration"]:
        fields.append(f"Enumeration: {enumeration_literal(point['enumeration'])}")
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
        named = controller_locations(root, grouped, services)
        register_service_locations(root, grouped, services)
    except (KeyError, TypeError, AttributeError) as error:
        raise InputError(f"malformed source row: {error!r}") from error
    if set(MAP_LABELS) - named:
        raise InputError(f"register map labels without a named register: {sorted(set(MAP_LABELS) - named)}")
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
