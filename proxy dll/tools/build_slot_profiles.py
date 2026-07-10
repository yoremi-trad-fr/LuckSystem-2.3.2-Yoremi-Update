#!/usr/bin/env python3
"""Build the shared Luca menu catalog and JP/CN source-slot profiles."""

from __future__ import annotations

import argparse
import importlib.util
import json
import re
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]

DEFAULT_EXES = {
    "Kanon": Path(r"C:\Program Files (x86)\Steam\steamapps\common\Kanon\Kanon.exe"),
    "AIR": Path(r"C:\Program Files (x86)\Steam\steamapps\common\AIR\AIR.exe"),
    "HarmoniaHD": Path(r"C:\Program Files (x86)\Steam\steamapps\common\Harmonia Full HD Edition\HarmoniaFHD.exe"),
    "Loopers": Path(r"C:\Program Files (x86)\Steam\steamapps\common\LOOPERS\LOOPERS.exe"),
}


# Exact strings used by the Japanese Luca menu. Values beginning with '@' refer
# to the context column in Kanon's Arabic-B table, which keeps the long strings
# in one authoritative place.
TARGETS_JP = {
    "Close": "閉じる",
    "Defaults": "初期設定に戻す",
    "Skip": "早送り",
    "Previously Read Only": "既読のみ",
    "Voice": "ボイス",
    "Stop on New Message": "メッセージ送りで停止する",
    "      No Stops      ": "　　 停止しない 　　",
    "Positioned at ❝Yes❞": "「はい」に合わせる",
    "Positioned at ❝No❞": "「いいえ」に合わせる",
    "Controller Rumble Function": "コントローラの振動",
    "Disable": "なし",
    "Language": "言語",
    "Font": "フォント",
    "Mincho": " 明朝 ",
    "Modern": " モダン ",
    "Solid": "不透明",
    "Clear": "透明",
    "Window Transparency": "ウィンドウ透明度",
    "Previously Read Text": "既読の色表示",
    "Color of": "",
    "Green": "緑",
    "Blue": "青",
    "Purple": "紫",
    "Yellow": "黄",
    "Read Text Color": "既読色",
    "Text Speed": "文字表示速度",
    "Slow": "遅い",
    "Fast": "速い",
    "0 sec/char": "0秒/文字",
    "0.1 sec/char": "0.2秒/文字",
    "Wait Time Per Character": "文字数待機時間",
    "0 sec": "0秒",
    "1 sec": "1秒",
    "2 sec": "2秒",
    "3 sec": "3秒",
    "Base Wait Time": "固定待機時間",
    "Master Volume": "マスター音量",
    "System Sounds": "システム音",
    "While Pressed": "押している間",
    "Start/Stop": "開始／解除",
    "C (Skip)": "Cキー (早送り)",
    "Z (Rewind)": "Zキー (早戻し)",
    "  Disable  ": "　無効　",
    "Quick Save": "クイックセーブ",
    "Quick Load": "クイックロード",
    "Switch Language": "言語切り替え",
    "Up Arrow": "方向キー上",
    "System Menu": "システムメニュー",
    "Hide Window": "ウィンドウ非表示",
    "Right Click": "右ボタン",
    "Left+Right Click": "左ボタン＋右ボタン",
    "Mouse Wheel Button": "ホイールボタン",
    "Rewind Once": "１つ戻し",
    "Forward Once": "１つ送り",
    "Return/Proceed Button": "戻る／進むボタン",
    "Enable": "あり",
    "Gestures": "ジェスチャー操作",
    "Game Ver. ": "Game Ver. ",
    "Window": "ウィンドウ",
    "Full Screen": "フルスクリーン",
    "Screen Mode": "画面モード",
    "Auto": "自動",
    "Window Size": "ウィンドウサイズ",
    "Wait Time Per Character: In Auto Mode, this sets the wait time until the next message is displayed based on the number of characters in text.": "@Wait per character tooltip",
    "Base Wait Time: You can set a Base Wait Time to add to\n❝Wait Time Per Character❞.": "@Base wait tooltip",
    "Controller Rumble Function: Plug in a controller to use the controller Rumble function.": "@Controller rumble tooltip",
    "Gestures: Moving the cursor while holding the left button works the same way as Touch controls.": "@Gestures tooltip",
    "Left+Right Click: Hold left button then right click to switch between languages (English/Simplified Chinese/Japanese).": "@Left/right click tooltip",
    "Do you wish to load this save?": "ロードしますか？",
    "Are you sure you wish to overwrite data?": "上書きしますか？",
    "Do you wish to save?": "セーブしますか？",
    "Are you sure you wish to delete this save data?": "削除しますか？",
    "This cannot be deleted.": "削除できません。",
    "Delete": " 削除 ",
    "Latest": "最新へ",
    "Text preview.": "テキストプレビュー。",
    "Settings such as text speed are reflected.": "文字表示速度などの設定が反映されます。",
    "English": "日本語",
    "Return to the title screen?": "タイトル画面に戻りますか？",
    "Return to the menu?": "メニューに戻りますか？",
    "$A1There is unsaved data.\n$A1Are you sure you wish to quit the game?": "@Quit prompt unsaved",
    "Yes": "はい",
    "No": "いいえ",
    "Save completed.": "セーブが完了しました。",
    "Are you sure you wish to quit the game?": "ゲームを終了しますか？",
    "Changing the setting to ❝Auto❞ will open the window at a scale based on the Windows ❝Display❞ setting.\nChanging ❝%%❞ will scale the display, with the default resolution being %d×%d pixels.\n (Scale cannot be increased beyond the maximum resolution of your display.)": "@System scale tooltip",
}


# These are menu-ready Simplified Chinese labels. Exact matches are also used
# to locate the CN source slot; labels absent from a game are simply skipped.
TARGETS_CN = {
    "Close": "关闭",
    "Defaults": "恢复默认设置",
    "Skip": "快进",
    "Previously Read Only": "仅限已读内容",
    "Voice": "语音",
    "Stop on New Message": "跳过对白时停止配音",
    "      No Stops      ": "　　　不停止　　　",
    "Positioned at ❝Yes❞": "位于“是”",
    "Positioned at ❝No❞": "位于“否”",
    "Controller Rumble Function": "控制器震动",
    "Disable": "无效",
    "Language": "语言",
    "Font": "字体",
    "Mincho": "宋体",
    "Modern": "现代",
    "Solid": "不透明",
    "Clear": "透明",
    "Window Transparency": "窗口透明度",
    "Previously Read Text": "已读内容颜色标记",
    "Color of": "颜色",
    "Green": "绿",
    "Blue": "蓝",
    "Purple": "紫",
    "Yellow": "黄",
    "Read Text Color": "已读色",
    "Text Speed": "文字速度",
    "Slow": "慢",
    "Fast": "快",
    "0 sec/char": "0秒/字",
    "0.1 sec/char": "0.2秒/字",
    "Wait Time Per Character": "按字数翻页",
    "0 sec": "0秒",
    "1 sec": "1秒",
    "2 sec": "2秒",
    "3 sec": "3秒",
    "Base Wait Time": "固定等待时间",
    "Master Volume": "整体音量",
    "System Sounds": "系统音",
    "While Pressed": "长按中持续",
    "Start/Stop": "开始／结束",
    "C (Skip)": "C键 (快进)",
    "Z (Rewind)": "Z键 (快退)",
    "  Disable  ": "  无效  ",
    "Quick Save": "快速保存",
    "Quick Load": "快速读取",
    "Switch Language": "切换语言",
    "Up Arrow": "上方向键",
    "System Menu": "系统菜单",
    "Hide Window": "隐藏窗口",
    "Right Click": "右键",
    "Left+Right Click": "左键＋右键",
    "Mouse Wheel Button": "鼠标滚轮",
    "Rewind Once": "后退一项",
    "Forward Once": "前进一项",
    "Return/Proceed Button": "返回/前进键",
    "Enable": "有效",
    "Gestures": "手势操作",
    "Game Ver. ": "Game Ver. ",
    "Window": "窗口模式",
    "Full Screen": "全屏模式",
    "Screen Mode": "显示模式",
    "Auto": "自动",
    "Window Size": "窗口大小",
    "Wait Time Per Character: In Auto Mode, this sets the wait time until the next message is displayed based on the number of characters in text.": "自动模式中，按文字数设置显示下一条消息前的等待时间。",
    "Base Wait Time: You can set a Base Wait Time to add to\n❝Wait Time Per Character❞.": "可设置叠加在按字数翻页时间上的固定等待时间。",
    "Controller Rumble Function: Plug in a controller to use the controller Rumble function.": "连接控制器后可使用震动功能。",
    "Gestures: Moving the cursor while holding the left button works the same way as Touch controls.": "按住左键移动鼠标可进行触摸操作。",
    "Left+Right Click: Hold left button then right click to switch between languages (English/Simplified Chinese/Japanese).": "按住左键后点击右键可切换语言。",
    "Do you wish to load this save?": "是否读取？",
    "Are you sure you wish to overwrite data?": "是否覆盖？",
    "Do you wish to save?": "是否保存？",
    "Are you sure you wish to delete this save data?": "是否删除？",
    "This cannot be deleted.": "无法删除。",
    "Delete": " 删除 ",
    "Latest": "跳至最近一次的存档",
    "Text preview.": "文字预览。",
    "Settings such as text speed are reflected.": "此文反映了文字速度等的设置结果。",
    "English": "简体中文",
    "Return to the title screen?": "是否返回标题画面？",
    "Return to the menu?": "是否返回菜单？",
    "$A1There is unsaved data.\n$A1Are you sure you wish to quit the game?": "$A1有数据尚未保存。\n$A1是否结束游戏？",
    "Yes": "是",
    "No": "否",
    "Save completed.": "已保存。",
    "Are you sure you wish to quit the game?": "是否结束游戏？",
    "Changing the setting to ❝Auto❞ will open the window at a scale based on the Windows ❝Display❞ setting.\nChanging ❝%%❞ will scale the display, with the default resolution being %d×%d pixels.\n (Scale cannot be increased beyond the maximum resolution of your display.)": "自动模式会依据 Windows 显示设置缩放窗口。百分比以 %d×%d 为标准，且不会超过屏幕分辨率。",
}


def load_module(name: str, path: Path):
    spec = importlib.util.spec_from_file_location(name, path)
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module


def all_hits(data: bytes, needle: bytes) -> list[int]:
    hits = []
    pos = data.find(needle)
    while pos >= 0:
        hits.append(pos)
        pos = data.find(needle, pos + 1)
    return hits


def string_hits(data: bytes, needle: bytes) -> list[int]:
    return [
        offset
        for offset in all_hits(data, needle)
        if (offset == 0 or data[offset - 1] == 0)
        and offset + len(needle) < len(data)
        and data[offset + len(needle)] == 0
    ]


def nearest_hit(data: bytes, text: str, reference: int) -> int | None:
    if not text:
        return None
    hits = string_hits(data, text.encode("utf-8"))
    return min(hits, key=lambda off: abs(off - reference)) if hits else None


def slot_size(data: bytes, offset: int, source: bytes) -> int:
    pos = offset + len(source)
    while pos < len(data) and data[pos] == 0:
        pos += 1
    return pos - offset


def py_string(value: str) -> str:
    return repr(value)


def slug(value: str) -> str:
    value = re.sub(r"[^a-z0-9]+", "-", value.lower()).strip("-")
    return value[:48] or "menu-string"


def resolve_japanese_targets(arabic_module) -> dict[str, str]:
    by_context = {}
    for _, source, _, context, _ in list(arabic_module.JP_SLOT_PATCHES) + list(arabic_module.ENGLISH_FALLBACK_PATCHES):
        by_context.setdefault(context, source.decode("utf-8"))
    resolved = {}
    for source_en, value in TARGETS_JP.items():
        if value.startswith("@"):
            value = by_context.get(value[1:], "")
        resolved[source_en] = value
    return resolved


def common_catalog(exes: dict[str, Path]):
    modules = {
        game: load_module(f"slot_base_{game}", ROOT / game / "patches.py")
        for game in DEFAULT_EXES
    }
    source_sets = {
        game: {source.decode("utf-8", errors="replace") for _, source, _, _, _ in module.PATCHES}
        for game, module in modules.items()
    }
    common = set.intersection(*source_sets.values())

    french = {}
    for game in ("HarmoniaHD", "Loopers"):
        for _, source, target, _, _ in modules[game].PATCHES:
            source = source.decode("utf-8", errors="replace")
            if source in common and target:
                french.setdefault(source, []).append(target)

    arabic_module = load_module("slot_arabic", ROOT / "Kanon" / "Arabic-B" / "patches.py")
    jp_targets = resolve_japanese_targets(arabic_module)
    resolved_arabic = arabic_module._resolved_patches()
    kanon_data = exes["Kanon"].read_bytes()
    resolved_by_source = {}
    for offset, source, target, _, _ in resolved_arabic:
        resolved_by_source.setdefault(source.decode("utf-8"), []).append((offset, target))

    rows = []
    used_ids = set()
    for offset, source_bytes, _, context, note in modules["Kanon"].PATCHES:
        source = source_bytes.decode("utf-8", errors="replace")
        if source not in common:
            continue
        fr_values = french.get(source, [])
        fr = min(fr_values, key=lambda value: len(value.encode("utf-8"))) if fr_values else ""
        jp = jp_targets.get(source, "")
        cn = TARGETS_CN.get(source, "")
        jp_offset = nearest_hit(kanon_data, jp, offset)
        cn_offset = nearest_hit(kanon_data, cn, offset)
        ar = ""
        if jp and jp in resolved_by_source:
            candidates = resolved_by_source[jp]
            _, ar = min(candidates, key=lambda item: abs(item[0] - (jp_offset or offset)))
        size = slot_size(kanon_data, offset, source_bytes)
        safe = bool(fr and fr != source and len(fr.encode("utf-8")) <= size - 1)
        catalog_id = slug(source)
        if catalog_id in used_ids:
            catalog_id = f"{catalog_id}-{slug(context)}"
        used_ids.add(catalog_id)
        rows.append({
            "id": catalog_id,
            "context": context,
            "en": source,
            "fr": fr,
            "ar": ar,
            "jp": jp,
            "cn": cn,
            "safe": safe,
            "kanonOffsets": {
                "en": f"0x{offset:X}",
                "jp": f"0x{jp_offset:X}" if jp_offset is not None else "",
                "cn": f"0x{cn_offset:X}" if cn_offset is not None else "",
            },
            "note": note,
        })
    return rows, modules


def write_profile(game: str, slot: str, rows: list[dict], module, exe: Path):
    data = exe.read_bytes()
    english_offsets = {
        source.decode("utf-8", errors="replace"): offset
        for offset, source, _, _, _ in module.PATCHES
    }
    key = "jp" if slot == "jp" else "cn"
    found = []
    used = set()
    for row in rows:
        source_text = row[key]
        en_offset = english_offsets.get(row["en"])
        if not source_text or en_offset is None:
            continue
        source = source_text.encode("utf-8")
        hits = [off for off in string_hits(data, source) if off not in used]
        if not hits:
            continue
        offset = min(hits, key=lambda off: abs(off - en_offset))
        used.add(offset)
        budget = slot_size(data, offset, source) - 1
        found.append((offset, source, row, budget))

    folder = ROOT / game / ("Slots-JP" if slot == "jp" else "Slots-CN")
    folder.mkdir(parents=True, exist_ok=True)
    output = folder / "patches.py"
    lines = [
        "#!/usr/bin/env python3",
        '"""Generated source-slot inventory for the LuckSystem Luca GUI."""',
        "",
        f"GAME_EXE = {str(exe.name)!r}",
        f"RVA_DELTA = {module.RVA_DELTA:#x}",
        f"PATCH_GAME_NAME = {game!r}",
        "PATCH_VERSION = 'slot-catalog-1'",
        "",
        "PATCHES = [",
    ]
    for offset, source, row, budget in sorted(found):
        note = f"catalog={row['id']}; {slot.upper()} source slot; budget {budget}"
        lines.append(
            f"    (0x{offset:X}, {source!r}, '', {py_string(row['context'])}, {py_string(note)}),"
        )
    lines.extend(["]", ""])
    output.write_text("\n".join(lines), encoding="utf-8")
    return len(found), output


def main():
    parser = argparse.ArgumentParser()
    for game, default in DEFAULT_EXES.items():
        parser.add_argument(f"--{game.lower()}", type=Path, default=default)
    args = parser.parse_args()
    exes = {game: getattr(args, game.lower()) for game in DEFAULT_EXES}
    missing = [f"{game}: {path}" for game, path in exes.items() if not path.is_file()]
    if missing:
        raise SystemExit("Missing game executable(s):\n" + "\n".join(missing))

    rows, modules = common_catalog(exes)
    catalog_path = ROOT / "menu_catalog.json"
    catalog_path.write_text(json.dumps({"version": 1, "entries": rows}, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    print(f"Catalog: {len(rows)} entries -> {catalog_path}")
    for game in DEFAULT_EXES:
        for slot in ("jp", "cn"):
            count, output = write_profile(game, slot, rows, modules[game], exes[game])
            print(f"{game} {slot.upper()}: {count} entries -> {output}")


if __name__ == "__main__":
    main()
