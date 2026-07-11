#!/usr/bin/env python3
"""Generate a mixed UTF-8/UTF-16LE patch table for an external LBEE table.

Run this script from the directory containing LITBUS_WIN32.exe and patches2.py.
The supplied patches2.py format is accepted unchanged:

    (raw_offset, source_bytes, translated_text, context, note)

An optional six-field form can specify the encoding explicitly:

    (raw_offset, source_bytes, translated_text, encoding, context, note)
"""

from __future__ import annotations

import csv
import ast
import struct
import sys
from pathlib import Path


GAME_EXE = Path("LITBUS_WIN32.exe")
PATCH_FILE = Path("patches2.py")
PATCH_GAME_NAME = "Little Busters! English Edition"
PATCH_VERSION = "1.2-winmm-x86"

# Dolamroth's current Russian strings are test data and have not all been
# shortened manually yet. Keep this enabled for a runnable proof-of-concept;
# every automatic crop is listed in patches.csv and in the console report.
CROP_OVERSIZE = True

# The original fetch_from_exe.py preferred the first UTF-16 occurrence and
# collapsed duplicate source strings. That selects substrings such as
# "Close" in "CloseThreadpool" instead of the standalone menu label. Replace
# such entries with every standalone UTF-8/UTF-16LE occurrence in string-data
# sections, while retaining the supplied offset when no standalone match exists.
EXPAND_STANDALONE_OCCURRENCES = True


def load_patches(path: Path):
    tree = ast.parse(path.read_text(encoding="utf-8"), filename=str(path))
    for statement in tree.body:
        if (
            isinstance(statement, ast.Assign)
            and any(
                isinstance(target, ast.Name) and target.id == "PATCHES"
                for target in statement.targets
            )
        ):
            return ast.literal_eval(statement.value)
    raise ValueError(f"{path} does not define a literal PATCHES list")


def pe_sections(data: bytes):
    if data[:2] != b"MZ":
        raise ValueError("not a PE executable (missing MZ header)")
    pe_offset = struct.unpack_from("<I", data, 0x3C)[0]
    if data[pe_offset : pe_offset + 4] != b"PE\0\0":
        raise ValueError("not a PE executable (missing PE header)")
    machine, section_count = struct.unpack_from("<HH", data, pe_offset + 4)
    if machine != 0x014C:
        raise ValueError(
            f"LBEE requires a 32-bit x86 executable; PE machine is 0x{machine:04X}"
        )
    optional_size = struct.unpack_from("<H", data, pe_offset + 20)[0]
    section_table = pe_offset + 24 + optional_size
    sections = []
    for index in range(section_count):
        start = section_table + index * 40
        name = data[start : start + 8].split(b"\0", 1)[0].decode("ascii", "replace")
        virtual_size, virtual_address, raw_size, raw_offset = struct.unpack_from(
            "<IIII", data, start + 8
        )
        sections.append((name, raw_offset, raw_size, virtual_address, virtual_size))
    return sections


def raw_to_rva(raw_offset: int, sections) -> tuple[int, str]:
    for name, section_raw, raw_size, virtual_address, _ in sections:
        if section_raw <= raw_offset < section_raw + raw_size:
            return virtual_address + raw_offset - section_raw, name
    raise ValueError(f"raw offset 0x{raw_offset:X} is outside all PE sections")


def detect_encoding(source: bytes) -> str:
    if len(source) % 2:
        return "utf-8"
    try:
        source.decode("utf-8")
        valid_utf8 = True
    except UnicodeDecodeError:
        valid_utf8 = False
    odd = source[1::2]
    looks_wide = bool(odd) and odd.count(0) * 2 >= len(odd)
    return "utf-16-le" if not valid_utf8 or looks_wide else "utf-8"


def crop_to_budget(text: str, encoding: str, budget: int) -> tuple[str, bytes]:
    encoded = text.encode(encoding)
    while encoded and len(encoded) > budget:
        text = text[:-1]
        encoded = text.encode(encoding)
    return text, encoded


def c_bytes(value: bytes) -> str:
    return ",".join(f"0x{byte:02X}" for byte in value)


def c_string(value: str) -> str:
    return (
        value.replace("\\", "\\\\")
        .replace('"', '\\"')
        .replace("\n", "\\n")
        .replace("\r", "\\r")
        .replace("\t", "\\t")
    )


def normalize(entry):
    if len(entry) == 5:
        raw_offset, source, target, context, note = entry
        encoding = detect_encoding(source)
    elif len(entry) == 6:
        raw_offset, source, target, encoding, context, note = entry
        encoding = encoding.lower().replace("utf16le", "utf-16-le")
        if encoding not in {"utf-8", "utf-16-le"}:
            raise ValueError(f"unsupported encoding {encoding!r} at 0x{raw_offset:X}")
    else:
        raise ValueError(f"patch entry must contain 5 or 6 values, got {len(entry)}")
    return raw_offset, source, target, encoding, context, note


def standalone_occurrences(data: bytes, sections, source_text: str):
    matches = []
    string_sections = {".rdata", ".data", ".rodata", "_RDATA"}
    for encoding in ("utf-8", "utf-16-le"):
        needle = source_text.encode(encoding)
        terminator_size = 2 if encoding == "utf-16-le" else 1
        terminator = bytes(terminator_size)
        for name, section_raw, raw_size, _, _ in sections:
            if name not in string_sections:
                continue
            section_end = section_raw + raw_size
            cursor = section_raw
            while True:
                offset = data.find(needle, cursor, section_end)
                if offset < 0:
                    break
                before_ok = (
                    offset == section_raw
                    or data[offset - terminator_size : offset] == terminator
                )
                after = offset + len(needle)
                after_ok = (
                    after + terminator_size <= section_end
                    and data[after : after + terminator_size] == terminator
                )
                if before_ok and after_ok:
                    matches.append((offset, needle, encoding))
                cursor = offset + 1
    return matches


def expand_standalone_patches(data: bytes, sections, patches):
    expanded = []
    replaced_offsets = 0
    skipped_noops = 0
    for entry in patches:
        raw_offset, source, target, encoding, context, note = normalize(entry)
        source_text = source.decode(encoding)
        if target == source_text:
            skipped_noops += 1
            continue
        matches = standalone_occurrences(data, sections, source_text)
        if matches:
            if not any(off == raw_offset and raw == source for off, raw, _ in matches):
                replaced_offsets += 1
            for off, raw, match_encoding in matches:
                expanded.append(
                    (off, raw, target, match_encoding, context, note)
                )
        else:
            expanded.append(
                (raw_offset, source, target, encoding, context, note)
            )

    unique = []
    seen = set()
    for entry in expanded:
        key = (entry[0], entry[1])
        if key in seen:
            continue
        seen.add(key)
        unique.append(entry)
    return unique, replaced_offsets, skipped_noops


def build_rows(data: bytes, sections, patches):
    rows = []
    errors = []
    for index, entry in enumerate(patches):
        raw_offset, source, target, encoding, context, note = normalize(entry)
        actual = data[raw_offset : raw_offset + len(source)]
        if actual != source:
            errors.append(
                f"#{index} 0x{raw_offset:X}: expected {source!r}, got {actual!r}"
            )
            continue

        rva, section = raw_to_rva(raw_offset, sections)
        source_text = source.decode(encoding)
        terminator_size = 2 if encoding == "utf-16-le" else 1
        after = raw_offset + len(source)
        zero_run = 0
        while after + zero_run < len(data) and data[after + zero_run] == 0:
            zero_run += 1

        # A complete terminator gives us a real string slot. Otherwise this is
        # a fixed-size field or substring, so never write beyond source bytes.
        terminated = zero_run >= terminator_size
        write_len = len(source) + zero_run if terminated else len(source)
        budget = write_len - terminator_size if terminated else write_len
        original_target = target
        target_bytes = target.encode(encoding)
        cropped = False
        if len(target_bytes) > budget:
            if CROP_OVERSIZE:
                target, target_bytes = crop_to_budget(target, encoding, budget)
                cropped = target != original_target
            else:
                errors.append(
                    f"#{index} 0x{raw_offset:X}: target {len(target_bytes)}B exceeds "
                    f"the {budget}B budget"
                )
                continue

        expected = source + bytes(write_len - len(source))
        replacement = target_bytes + bytes(write_len - len(target_bytes))
        rows.append(
            {
                "index": index,
                "raw_offset": raw_offset,
                "rva": rva,
                "section": section,
                "encoding": encoding,
                "source": source_text,
                "target": target,
                "original_target": original_target,
                "source_len": len(source),
                "target_len": len(target_bytes),
                "write_len": write_len,
                "budget": budget,
                "cropped": cropped,
                "expected": expected,
                "replacement": replacement,
                "context": context,
                "note": note,
            }
        )

    if errors:
        raise ValueError("\n".join(errors))
    return rows


def write_header(rows):
    with Path("patches.h").open("w", encoding="utf-8", newline="\n") as output:
        output.write("/* Auto-generated from patches2.py. Do not edit. */\n")
        output.write("#ifndef LUCKPROXY_PATCHES_H\n#define LUCKPROXY_PATCHES_H\n\n")
        output.write(f'#define PATCH_GAME_NAME "{c_string(PATCH_GAME_NAME)}"\n')
        output.write(f'#define PATCH_VERSION   "{c_string(PATCH_VERSION)}"\n\n')
        for number, row in enumerate(rows):
            output.write(
                f"static const BYTE s_src_{number:04d}[] = "
                f"{{ {c_bytes(row['expected'])} }};\n"
            )
            output.write(
                f"static const BYTE s_tgt_{number:04d}[] = "
                f"{{ {c_bytes(row['replacement'])} }};\n"
            )
        output.write("\nstatic const LuckPatch g_patches[] = {\n")
        for number, row in enumerate(rows):
            label = row["context"] or f"LBEE {row['encoding']}"
            comment = c_string(f"{label}: {row['source'][:60]}")
            output.write(
                f"    {{ 0x{row['rva']:08X}, {row['write_len']:4}, "
                f"s_src_{number:04d}, s_tgt_{number:04d}, \"{comment}\" }},\n"
            )
        output.write("};\n\n")
        output.write("#define N_PATCHES (sizeof(g_patches)/sizeof(g_patches[0]))\n\n")
        output.write("#endif\n")


def write_csv(rows):
    fields = [
        "index",
        "raw_offset",
        "rva",
        "section",
        "encoding",
        "write_len",
        "budget",
        "source_len",
        "target_len",
        "cropped",
        "source",
        "target",
        "original_target",
        "context",
        "note",
    ]
    with Path("patches.csv").open("w", encoding="utf-8-sig", newline="") as output:
        writer = csv.DictWriter(output, fieldnames=fields, extrasaction="ignore")
        writer.writeheader()
        for row in rows:
            serial = dict(row)
            serial["raw_offset"] = f"0x{row['raw_offset']:X}"
            serial["rva"] = f"0x{row['rva']:X}"
            writer.writerow(serial)


def main() -> int:
    try:
        data = GAME_EXE.read_bytes()
        patches = load_patches(PATCH_FILE)
        sections = pe_sections(data)
        original_count = len(patches)
        replaced_offsets = 0
        skipped_noops = 0
        if EXPAND_STANDALONE_OCCURRENCES:
            patches, replaced_offsets, skipped_noops = expand_standalone_patches(
                data, sections, patches
            )
        rows = build_rows(data, sections, patches)
    except (OSError, AttributeError, ImportError, ValueError) as error:
        print(f"ERROR: {error}", file=sys.stderr)
        return 1

    write_header(rows)
    write_csv(rows)
    utf16 = sum(row["encoding"] == "utf-16-le" for row in rows)
    cropped = sum(row["cropped"] for row in rows)
    sections_used = ", ".join(sorted({row["section"] for row in rows}))
    print(f"Generated patches.h and patches.csv with {len(rows)} entries.")
    if EXPAND_STANDALONE_OCCURRENCES:
        print(
            f"Standalone expansion: {original_count} -> {len(rows)} entries; "
            f"replaced {replaced_offsets} non-standalone first-match offsets; "
            f"removed {skipped_noops} no-op source=target rows."
        )
    print(f"Encodings: {utf16} UTF-16LE, {len(rows) - utf16} UTF-8.")
    print(f"PE sections: {sections_used}; automatically cropped: {cropped}.")
    if cropped:
        print("Review cropped=True rows in patches.csv before a release build.")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
