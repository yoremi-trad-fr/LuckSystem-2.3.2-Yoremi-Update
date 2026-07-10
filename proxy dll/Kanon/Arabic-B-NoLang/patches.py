#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Kanon Steam Arabic-B proxy string table.

This build targets the Japanese language slot only. English strings are left
untouched so the English slot remains English.

Arabic-B text is stored as Arabic Presentation Forms and reversed per line so
Luck Engine's left-to-right renderer can display it.
"""

from pathlib import Path
import importlib.util
import re

HERE = Path(__file__).resolve().parent
BASE_PATH = HERE.parent / "patches.py"

spec = importlib.util.spec_from_file_location("kanon_base_patches", BASE_PATH)
base = importlib.util.module_from_spec(spec)
spec.loader.exec_module(base)


# tuple order: isolated, final, initial, medial
ARABIC_FORMS = {
    "ء": ("ﺀ", None, None, None),
    "آ": ("ﺁ", "ﺂ", None, None),
    "أ": ("ﺃ", "ﺄ", None, None),
    "ؤ": ("ﺅ", "ﺆ", None, None),
    "إ": ("ﺇ", "ﺈ", None, None),
    "ئ": ("ﺉ", "ﺊ", "ﺋ", "ﺌ"),
    "ا": ("ﺍ", "ﺎ", None, None),
    "ب": ("ﺏ", "ﺐ", "ﺑ", "ﺒ"),
    "ة": ("ﺓ", "ﺔ", None, None),
    "ت": ("ﺕ", "ﺖ", "ﺗ", "ﺘ"),
    "ث": ("ﺙ", "ﺚ", "ﺛ", "ﺜ"),
    "ج": ("ﺝ", "ﺞ", "ﺟ", "ﺠ"),
    "ح": ("ﺡ", "ﺢ", "ﺣ", "ﺤ"),
    "خ": ("ﺥ", "ﺦ", "ﺧ", "ﺨ"),
    "د": ("ﺩ", "ﺪ", None, None),
    "ذ": ("ﺫ", "ﺬ", None, None),
    "ر": ("ﺭ", "ﺮ", None, None),
    "ز": ("ﺯ", "ﺰ", None, None),
    "س": ("ﺱ", "ﺲ", "ﺳ", "ﺴ"),
    "ش": ("ﺵ", "ﺶ", "ﺷ", "ﺸ"),
    "ص": ("ﺹ", "ﺺ", "ﺻ", "ﺼ"),
    "ض": ("ﺽ", "ﺾ", "ﺿ", "ﻀ"),
    "ط": ("ﻁ", "ﻂ", "ﻃ", "ﻄ"),
    "ظ": ("ﻅ", "ﻆ", "ﻇ", "ﻈ"),
    "ع": ("ﻉ", "ﻊ", "ﻋ", "ﻌ"),
    "غ": ("ﻍ", "ﻎ", "ﻏ", "ﻐ"),
    "ف": ("ﻑ", "ﻒ", "ﻓ", "ﻔ"),
    "ق": ("ﻕ", "ﻖ", "ﻗ", "ﻘ"),
    "ك": ("ﻙ", "ﻚ", "ﻛ", "ﻜ"),
    "ل": ("ﻝ", "ﻞ", "ﻟ", "ﻠ"),
    "م": ("ﻡ", "ﻢ", "ﻣ", "ﻤ"),
    "ن": ("ﻥ", "ﻦ", "ﻧ", "ﻨ"),
    "ه": ("ﻩ", "ﻪ", "ﻫ", "ﻬ"),
    "و": ("ﻭ", "ﻮ", None, None),
    "ى": ("ﻯ", "ﻰ", None, None),
    "ي": ("ﻱ", "ﻲ", "ﻳ", "ﻴ"),
}


def _joins_next(ch):
    forms = ARABIC_FORMS.get(ch)
    return bool(forms and forms[2])


def _joins_prev(ch):
    forms = ARABIC_FORMS.get(ch)
    return bool(forms and forms[1])


def _shape_arabic(text):
    out = []
    chars = list(text)
    for i, ch in enumerate(chars):
        forms = ARABIC_FORMS.get(ch)
        if not forms:
            out.append(ch)
            continue

        prev_ch = chars[i - 1] if i > 0 else ""
        next_ch = chars[i + 1] if i + 1 < len(chars) else ""
        join_prev = _joins_prev(ch) and _joins_next(prev_ch)
        join_next = _joins_next(ch) and _joins_prev(next_ch)

        isolated, final, initial, medial = forms
        if join_prev and join_next and medial:
            out.append(medial)
        elif join_prev and final:
            out.append(final)
        elif join_next and initial:
            out.append(initial)
        else:
            out.append(isolated)
    return "".join(out)


def _visual_line(line):
    tag = ""
    match = re.match(r"^(\$[A-Z]\d+)", line)
    if match:
        tag = match.group(1)
        line = line[len(tag):]
    return tag + "".join(reversed(_shape_arabic(line)))


def A(text):
    """Arabic-B: shaped Arabic Presentation Forms, reversed per line."""
    return "\n".join(_visual_line(line) for line in text.split("\n"))


def J(text):
    return text.encode("utf-8")


# Japanese-slot source strings only. Do not add English source offsets here.
JP_SLOT_PATCHES = [
    # Global language and dialog buttons.
    (0x487470, J("はい"), "ن", "Global dialog yes", "Short yes because slot budget is 7 bytes"),
    (0x487480, J("いいえ"), "لا", "Global dialog no", "Japanese slot"),
    (0x487498, J("閉じる"), "إغ", "Close button", "Japanese slot"),
    (0x488AC8, J("戻る"), "رج", "Back button", "Japanese slot"),
    (0x488B20, J(" 削除 "), "حذف", "Delete button", "Japanese slot"),
    (0x488B48, J("最新へ"), "آخر", "Latest button", "Japanese slot"),
    (0x488B80, J("初期設定に戻す"), "ضبط", "Defaults button", "Japanese slot"),

    # Misc system menu pages and popups.
    (0x488C08, J("基本操作説明"), "أوامر", "Controls page", "Japanese slot"),
    (0x488C48, J("切り替え言語一覧"), "لغات", "Languages page", "Japanese slot"),
    (0x488C90, J("ボイス再生"), "صوت", "Play voice", "Japanese slot"),
    (0x488CC8, J("ジャンプ"), "قفز", "Jump command", "Japanese slot"),
    (0x488CE8, J("ボイスコレクション"), "أصوات", "Voice collection", "Japanese slot"),
    (0x488D20, J("BGM変更"), "لحن", "BGM select", "Japanese slot"),
    (0x488D50, J("初期化"), "ضبط", "Reset", "Japanese slot"),
    (0x488D68, J("メッセージログ"), "سجل", "Message log", "Japanese slot"),
    (0x488D8C, J("削除"), "حذف", "Delete", "Japanese slot"),
    (0x488D98, J("プレイリストへ"), "قائمة", "To playlist", "Japanese slot"),
    (0x488DC0, J("ボイス一覧へ"), "أصوات", "To voice list", "Japanese slot"),
    (0x488DF0, J("連続再生"), "تشغيل", "Continuous play", "Japanese slot"),
    (0x488E10, J("BGM選択"), "لحن", "BGM select", "Japanese slot"),
    (0x488E1C, J("編集"), "حر", "Edit", "Japanese slot"),
    (0x488E30, J("プレビュー表示"), "عرض", "Show preview", "Japanese slot"),
    (0x488E58, J("プレビュー非表示"), "إخفاء", "Hide preview", "Japanese slot"),
    (0x488E88, J("【長押し】連続再生"), "تشغيل مستمر", "Hold continuous play", "Japanese slot"),
    (0x488EC8, J("テキストプレビュー。"), "عرض النص.", "Text preview", "Japanese slot"),

    # Basic tab.
    (0x4987D0, J("なし"), "لا", "None value", "Japanese slot"),
    (0x4990B8, J(" ダイアログの初期カーソル位置 "), "المؤشر", "Initial cursor label", "Japanese slot"),
    (0x4990E8, J("コントローラの振動"), "اهتزاز", "Controller rumble label", "Japanese slot"),
    (0x49910C, J("基本"), "أس", "Basic tab", "Japanese slot"),
    (0x499118, J("非表示"), "إخف", "Shortcut menu hide", "Japanese slot"),
    (0x499158, J("表示"), "عر", "Shortcut menu display", "Japanese slot"),
    (0x499160, J("ショートカットメニュー"), "قائمة", "Shortcut menu", "Japanese slot"),
    (0x4991C8, J("既読のみ"), "مقروء", "Skip read only", "Japanese slot"),
    (0x4991D8, J("すべて"), "كل", "Skip all", "Japanese slot"),
    (0x499210, J("早送り"), "تخ", "Skip", "Japanese slot"),
    (0x49921C, J("下"), "ت", "Down value", "Japanese slot"),
    (0x499228, J("下部"), "أس", "Choices bottom", "Japanese slot"),
    (0x499240, J("中央"), "وس", "Choices center", "Japanese slot"),
    (0x499260, J("選択肢の位置"), "مكان", "Choices position", "Japanese slot"),
    (0x499278, J("メッセージ送りで停止する"), "قف", "Voice stop on message", "Japanese slot"),
    (0x499300, J("　　 停止しない 　　"), "لا توقف", "Voice no stops", "Japanese slot"),
    (0x499348, J("　オフ　"), "لا", "Off value", "Japanese slot"),
    (0x499360, J("　オン　"), "نعم", "On value", "Japanese slot"),
    (0x499370, J("日付表示"), "تاريخ", "Display date", "Japanese slot"),
    (0x4993C0, J("「はい」に合わせる"), "نعم", "Initial cursor yes", "Japanese slot"),
    (0x4993E0, J("「いいえ」に合わせる"), "لا", "Initial cursor no", "Japanese slot"),
    (0x49947C, J("弱"), "ض", "Rumble min", "Japanese slot"),
    (0x499484, J("中"), "و", "Rumble mid", "Japanese slot"),
    (0x49948C, J("強"), "ق", "Rumble max", "Japanese slot"),

    # Touch tab.
    (0x4999C0, J("タッチ"), "لمس", "Touch tab", "Japanese slot"),
    (0x4999D8, J("早戻し"), "رج", "Touch rewind", "Japanese slot"),
    (0x499A08, J("ジャンプ（送り）"), "تقدم", "Touch jump forward", "Japanese slot"),
    (0x499A28, J("ジャンプ（戻し）"), "رجوع", "Touch jump backward", "Japanese slot"),

    # Text tabs.
    (0x49A4B0, J("テキスト1"), "نص", "Text1 tab", "Japanese slot"),
    (0x49A49C, J("橙"), "ب", "Orange color value", "Japanese slot"),
    (0x49A4A0, J("黄"), "ص", "Yellow color value", "Japanese slot"),
    (0x49A4E4, J("言語"), "لغة", "Language", "Japanese slot"),
    (0x49A4F0, J(" ゴシック "), "خط1", "Font family value", "Japanese slot"),
    (0x49A508, J(" 丸ゴシック "), "خط2", "Font family value", "Japanese slot"),
    (0x49A528, J(" 太丸ゴシック "), "خط3", "Font family value", "Japanese slot"),
    (0x49A548, J(" 明朝 "), "خط4", "Font family value", "Japanese slot"),
    (0x49A560, J(" モダン "), "خط5", "Font family value", "Japanese slot"),
    (0x49A578, J("フォント"), "خط", "Font", "Japanese slot"),
    (0x49A5B8, J("不透明"), "صلب", "Opaque value", "Japanese slot"),
    (0x49A5CC, J("透明"), "شف", "Transparent value", "Japanese slot"),
    (0x49A5E0, J("ウィンドウ透明度"), "شفافية", "Window transparency", "Japanese slot"),
    (0x49A640, J("選択肢のみオン"), "خيار", "Only choices", "Japanese slot"),
    (0x49A698, J("既読の色表示"), "لون", "Read text color", "Japanese slot"),
    (0x49A6AC, J("緑"), "خ", "Green color value", "Japanese slot"),
    (0x49A6E0, J("青"), "ز", "Blue color value", "Japanese slot"),
    (0x49A6E4, J("紫"), "ب", "Purple color value", "Japanese slot"),
    (0x49A6EC, J("赤"), "ح", "Red color value", "Japanese slot"),
    (0x49A710, J("既読色"), "لون", "Read color short", "Japanese slot"),
    (0x49AA38, J("文字数待機時間"), "انتظار", "Wait per character", "Japanese slot"),
    (0x49AAC8, J("テキスト2"), "نص", "Text2 tab", "Japanese slot"),
    (0x49AAD8, J("遅い"), "بط", "Slow", "Japanese slot"),
    (0x49AAFC, J("速い"), "سر", "Fast", "Japanese slot"),
    (0x49AB08, J("文字表示速度"), "سرعة", "Text speed", "Japanese slot"),
    (0x49AB50, J("0秒/文字"), "0ث/حرف", "Text speed value", "Japanese slot"),
    (0x49AB60, J("0.2秒/文字"), "0.2ث/ح", "Text speed value", "Japanese slot"),
    (0x49ABD0, J("0秒"), "0ث", "Wait value", "Japanese slot"),
    (0x49ABE0, J("1秒"), "1ث", "Wait value", "Japanese slot"),
    (0x49ABF0, J("2秒"), "2ث", "Wait value", "Japanese slot"),
    (0x49AC00, J("3秒"), "3ث", "Wait value", "Japanese slot"),
    (0x49AC18, J("固定待機時間"), "مهلة", "Base wait time", "Japanese slot"),

    # Sound, voice, keyboard.
    (0x49AC30, J("サウンド"), "صوت", "Sound tab", "Japanese slot"),
    (0x49AC54, J("音量"), "صت", "Volume", "Japanese slot"),
    (0x49AC78, J("マスター音量"), "عام", "Master volume", "Japanese slot"),
    (0x49AC90, J("ボイス"), "صت", "Voice volume", "Japanese slot"),
    (0x49ACC8, J("効果音"), "اثر", "SFX value", "Japanese slot"),
    (0x49ACD8, J("システム音"), "نظام", "System sound", "Japanese slot"),
    (0x49ADB0, J("キーボード"), "لوح", "Keyboard tab", "Japanese slot"),
    (0x49ADC0, J("押している間"), "ضغط", "While pressed", "Japanese slot"),
    (0x49AE10, J("開始／解除"), "بدء", "Start stop", "Japanese slot"),
    (0x49AE20, J("Cキー (早送り)"), "تخطي", "C skip", "Japanese slot"),
    (0x49AE78, J("Zキー (早戻し)"), "رجوع", "Z rewind", "Japanese slot"),
    (0x49AE90, J("　無効　"), "لا", "Disable", "Japanese slot"),
    (0x49AEE0, J("クイックセーブ"), "حفظ", "Quick save", "Japanese slot"),
    (0x49AEF8, J("Sキー"), "S", "S key value", "Japanese slot"),
    (0x49AF30, J("クイックロード"), "تحميل", "Quick load", "Japanese slot"),
    (0x49AF48, J("言語切り替え"), "لغة", "Switch language", "Japanese slot"),
    (0x49AF98, J("Lキー"), "L", "L key value", "Japanese slot"),
    (0x49AFA0, J("無効"), "لا", "Disable value", "Japanese slot"),
    (0x49AFB0, J("方向キー上"), "فوق", "Up arrow", "Japanese slot"),

    # Mouse and system tabs.
    (0x49AFF0, J("左ボタン＋右ボタン"), "يسار يم", "Left and right button", "Japanese slot"),
    (0x49B010, J("ジェスチャー操作"), "لمس", "Gestures", "Japanese slot"),
    (0x49B038, J("マウス"), "فأر", "Mouse tab", "Japanese slot"),
    (0x49B048, J("システムメニュー"), "نظام", "System menu", "Japanese slot"),
    (0x49B0A0, J("ウィンドウ非表示"), "إخفاء", "Hide window", "Japanese slot"),
    (0x49B0C0, J("右ボタン"), "يمين", "Right button", "Japanese slot"),
    (0x49B138, J("ホイールボタン"), "عجلة", "Wheel button", "Japanese slot"),
    (0x49B150, J("１つ戻し"), "رجوع", "Rewind once", "Japanese slot"),
    (0x49B1A0, J("ホイール上方向"), "فوق", "Wheel up", "Japanese slot"),
    (0x49B1B8, J("１つ送り"), "تقدم", "Forward once", "Japanese slot"),
    (0x49B210, J("ホイール下方向"), "تحت", "Wheel down", "Japanese slot"),
    (0x49B228, J("ジャンプとページ切り替え"), "صفحات", "Jump and switch pages", "Japanese slot"),
    (0x49B2B0, J("戻る／進むボタン"), "موافقة", "Return proceed button", "Japanese slot"),
    (0x49B2CC, J("あり"), "فعل", "Enable", "Japanese slot"),
    (0x49B338, J("     なし     "), "لا", "Disable padded", "Japanese slot"),
    (0x49B368, J("ダイアログと選択肢"), "حوار", "Dialog and choices", "Japanese slot"),
    (0x49B388, J("ポインターの自動移動"), "مؤشر", "Snap pointer", "Japanese slot"),
    (0x49B718, J("ウィンドウサイズ"), "حجم", "Window size", "Japanese slot"),
    (0x49B738, J("画面モード"), "وضع", "Screen mode", "Japanese slot"),
    (0x49B778, J("システム"), "نظم", "System tab", "Japanese slot"),
    (0x49B788, J("ウィンドウ"), "نافذة", "Window", "Japanese slot"),
    (0x49B7C8, J("フルスクリーン"), "ملء", "Full screen", "Japanese slot"),
    (0x49B800, J("自動"), "آل", "Auto", "Japanese slot"),
]


ENGLISH_FALLBACK_PATCHES = [
    # Confirmation windows: Arabic-B was blank here, so use short English.
    (0x487E90, J("ロードしますか？"), "Load this save?", "Load prompt", "ASCII fallback"),
    (0x488060, J("セーブしますか？"), "Save data?", "Save prompt", "ASCII fallback"),
    (0x4880A8, J("上書きしますか？"), "Overwrite?", "Overwrite prompt", "ASCII fallback"),
    (0x488290, J("セーブが完了しました。"), "Save complete.", "Save completed prompt", "ASCII fallback"),
    (0x4882D8, J("削除しますか？"), "Delete data?", "Delete prompt", "ASCII fallback"),
    (0x488330, J("削除できません。"), "Cannot delete.", "Delete error", "ASCII fallback"),
    (0x488378, J("メニューに戻りますか？"), "Return to menu?", "Return menu prompt", "ASCII fallback"),
    (0x4883D0, J("タイトル画面に戻りますか？"), "Return to title?", "Return title prompt", "ASCII fallback"),
    (0x488438, J("このページを初期設定に戻しますか？"), "Reset this page?", "Reset page prompt", "ASCII fallback"),
    (0x488540, J("次の選択肢へジャンプしますか？"), "Jump to next choice?", "Jump choice prompt", "ASCII fallback"),
    (0x4885B8, J("次のチャプターへジャンプしますか？"), "Jump to next chapter?", "Jump chapter prompt", "ASCII fallback"),
    (0x488640, J("次のチャプターまたは選択肢までジャンプしますか？"), "Jump to next chapter/choice?", "Jump chapter choice prompt", "ASCII fallback"),
    (0x4886F0, J("次の選択肢へジャンプ..."), "Jumping to next choice...", "Jumping choice status", "ASCII fallback"),
    (0x488758, J("次のチャプターへジャンプ..."), "Jumping to next chapter...", "Jumping chapter status", "ASCII fallback"),
    (0x488798, J("＿＿ジャンプ中...＿＿"), "Jumping...", "Jumping status", "ASCII fallback"),
    (0x4887C0, J("$A1未読のテキストがあるため、\nその位置でジャンプを停止しました。"),
     "$A1Unread text found.\n$A1Jump stopped there.", "Unread jump stop", "ASCII fallback"),
    (0x488880, J("1つ前の選択肢にジャンプしますか？"), "Jump to previous choice?", "Previous choice prompt", "ASCII fallback"),
    (0x488908, J("このチャプターの先頭にジャンプしますか？"), "Jump to chapter start?", "Chapter start prompt", "ASCII fallback"),
    (0x4889A8, J("1つ前のチャプターにジャンプしますか？"), "Jump to previous chapter?", "Previous chapter prompt", "ASCII fallback"),
    (0x488A30, J("ゲームの先頭にジャンプしますか？"), "Jump to game start?", "Game start prompt", "ASCII fallback"),
    (0x495B80, J("$A1セーブしていないデータがあります。\n$A1ゲームを終了しますか？"),
     "$A1Unsaved data exists.\n$A1Quit the game?", "Quit prompt unsaved", "ASCII fallback"),
    (0x495BE0, J("ゲームを終了しますか？"), "Quit the game?", "Quit prompt", "ASCII fallback"),

    # Cloud/save maintenance prompts.
    (0x4969B0, J("全セーブデータの削除\n＿終了＿"), "Delete all save data\n_Quit_", "Cloud delete menu", "ASCII fallback"),
    (0x496AC0, J("$A1システムデータを含む全セーブデータをクラウド上から削除し、\n$A1ゲームを初期状態に戻しますか？（削除すると元に戻せません）"),
     "$A1Delete all save data, including system data, from the cloud?\n$A1This resets the game and cannot be undone.",
     "Cloud delete confirmation", "ASCII fallback"),
    (0x496B78, J("セーブデータの削除中…"), "Deleting save data...", "Save delete progress", "ASCII fallback"),
    (0x496C88, J("削除が完了しました。"), "Delete complete.", "Delete complete", "ASCII fallback"),
    (0x496CA8, J("システムデータを削除しますか？"), "Delete system data?", "System data delete prompt", "ASCII fallback"),
    (0x496D58, J("システムデータの削除中…"), "Deleting system data...", "System data delete progress", "ASCII fallback"),
    (0x496DF0, J("全ての実績をリセットします。よろしいですか？"), "Reset all achievements?", "Achievement reset prompt", "ASCII fallback"),
    (0x496E68, J("スクリプトを再読み込みしました。"), "Script reloaded.", "Script reload message", "ASCII fallback"),

    # Help/tooltips at the bottom of the options screen.
    (0x488F10, J("$A1言語切り替えで使用する言語にチェックを入れてください。"),
     "$A1Select languages for switching.", "Language tooltip", "ASCII fallback"),
    (0x488FC8, J("文字表示速度などの設定が反映されます。"),
     "Text settings are previewed.", "Text preview tooltip", "ASCII fallback"),
    (0x489070, J("サウンドAPIはゲームをいったん終了し、再度開始すると切り替わります。"),
     "Restart the game to apply Sound API changes.", "Sound API tooltip", "ASCII fallback"),
    (0x4994B0, J("【ボイス再生】「停止しない」にすると、音声再生中にメッセージを送っても音声が流れ続けます。（次のメッセージに音声がある場合は停止します）"),
     "Voice: with No stop, voice keeps playing when text advances. It stops if the next line has voice.",
     "Voice tooltip", "ASCII fallback"),
    (0x499760, J("【ダイアログの初期カーソル位置】確認用の小ウィンドウが開いた時のカーソルの初期位置を設定します。"),
     "Initial cursor: sets the default choice in Yes/No windows.", "Initial cursor tooltip", "ASCII fallback"),
    (0x499800, J("【コントローラの振動】コントローラを接続すると、コントローラの振動機能が使用できます。"),
     "Controller rumble: connect a controller to use vibration.", "Controller rumble tooltip", "ASCII fallback"),
    (0x499AE0, J("メッセージウィンドウ内でドラッグすると早送り。\nフリックすると自動早送り（タップで解除）。"),
     "Drag inside the message window to skip forward.\nFlick for auto-skip; tap to cancel.",
     "Touch skip tooltip", "ASCII fallback"),
    (0x499B70, J("メッセージウィンドウ内でドラッグすると早戻し。\nフリックすると自動早戻し（タップで解除）。"),
     "Drag inside the message window to rewind.\nFlick for auto-rewind; tap to cancel.",
     "Touch rewind tooltip", "ASCII fallback"),
    (0x499D80, J("メッセージウィンドウ外でフリックすると「次の選択肢／チャプター」もしくは「既読部分の最後」までジャンプ。"),
     "Flick outside the message window to jump to the next choice/chapter or the end of read text.",
     "Touch jump next tooltip", "ASCII fallback"),
    (0x499E20, J("メッセージウィンドウ外でフリックすると「直前の選択肢／チャプター」までジャンプ。"),
     "Flick outside the message window to jump to the previous choice/chapter.",
     "Touch jump previous tooltip", "ASCII fallback"),
    (0x49A030, J("メッセージウィンドウ外でフリックすると「次のチャプター」もしくは「既読部分の最後」までジャンプ。"),
     "Flick outside the message window to jump to the next chapter or end of read text.",
     "Touch next chapter tooltip", "ASCII fallback"),
    (0x49A0D0, J("メッセージウィンドウ外でフリックすると「直前のチャプター」までジャンプ。"),
     "Flick outside the message window to jump to the previous chapter.",
     "Touch previous chapter tooltip", "ASCII fallback"),
    (0x49A2B0, J("メッセージウィンドウ外でフリックすると「次の選択肢」もしくは「既読部分の最後」までジャンプ。"),
     "Flick outside the message window to jump to the next choice or end of read text.",
     "Touch next choice tooltip", "ASCII fallback"),
    (0x49A340, J("メッセージウィンドウ外でフリックすると「直前の選択肢」までジャンプ。"),
     "Flick outside the message window to jump to the previous choice.",
     "Touch previous choice tooltip", "ASCII fallback"),
    (0x49A720, J("【文字数待機時間】オートモードの時、全メッセージが表示されてから、次のメッセージに進むまでの待ち時間を、メッセージの文字数に合わせて設定できます。"),
     "Wait per character: in auto mode, sets the wait before the next message based on text length.",
     "Wait per character tooltip", "ASCII fallback"),
    (0x49A9C0, J("【固定待機時間】「文字数待機時間」に追加する固定的な待ち時間を設定できます。"),
     "Base wait time: adds a fixed wait to wait per character.", "Base wait tooltip", "ASCII fallback"),
    (0x49B450, J("【ジェスチャー操作】左ボタンを押しながら、ゲーム画面上でマウスポインターを動かすと「タッチ」と同じ操作ができます。"),
     "Gestures: move the mouse while holding left click to use touch-style controls.",
     "Gestures tooltip", "ASCII fallback"),
    (0x49B500, J("【左ボタン＋右ボタン】左ボタンを押しながら右ボタンを押すと、テキストの表示言語が切り替わります（日本語／英語／簡体中文）。"),
     "Left+Right Click: hold left click, then right click to switch text language.",
     "Left/right click tooltip", "ASCII fallback"),
    (0x49B9A0, J("「自動」にすると、Windowsの「ディスプレイ」設定に基づいた表示スケールでウィンドウを開きます。\n「％」指定は、標準サイズ（%d×%dピクセル）に対しての拡大・縮小表示を行います。（ディスプレイの解像度を超えた拡大はできません）"),
     "Auto uses the Windows display scale.\nPercent values resize from the default %d x %d pixels. Scaling cannot exceed your display resolution.",
     "System scale tooltip", "ASCII fallback"),
]


ASCII_FALLBACK_TARGETS = {
    # These strings are drawn with another menu font style/size. Arabic-B glyphs
    # can become invisible there, so keep short ASCII labels for this test build.
    0x4987D0: "None",
    0x499118: "Hide",
    0x499158: "Show",
    0x4991C8: "Read",
    0x4991D8: "All",
    0x49921C: "DN",
    0x499228: "Bottom",
    0x499240: "Center",
    0x499278: "Stop",
    0x499300: "No stop",
    0x499348: "Off",
    0x499360: "On",
    0x4993C0: "Yes",
    0x4993E0: "No",
    0x49947C: "Low",
    0x499484: "Mid",
    0x49948C: "Hi",
    0x49A49C: "O",
    0x49A4A0: "Y",
    0x49A4F0: "Font1",
    0x49A508: "Font2",
    0x49A528: "Font3",
    0x49A548: "Font4",
    0x49A560: "Font5",
    0x49A5B8: "Opaque",
    0x49A5CC: "Trans",
    0x49A6AC: "G",
    0x49A6E0: "B",
    0x49A6E4: "P",
    0x49A6EC: "R",
    0x49AAD8: "Slow",
    0x49AAFC: "Fast",
    0x49AB50: "0s/char",
    0x49AB60: "0.2s/ch",
    0x49ABD0: "0s",
    0x49ABE0: "1s",
    0x49ABF0: "2s",
    0x49AC00: "3s",
    0x49ACC8: "SFX",
    0x49ADC0: "Hold",
    0x49AE10: "Toggle",
    0x49AE20: "Skip",
    0x49AE78: "Rewind",
    0x49AE90: "Off",
    0x49AEE0: "Q.Save",
    0x49AF30: "Q.Load",
    0x49AF48: "Lang",
    0x49AFA0: "Off",
    0x49AFB0: "Up",
    0x49AFF0: "L+R",
    0x49B048: "System",
    0x49B0A0: "Hide",
    0x49B0C0: "Right",
    0x49B138: "Wheel",
    0x49B150: "Back",
    0x49B1A0: "Wheel Up",
    0x49B1B8: "Forward",
    0x49B210: "Wheel Down",
    0x49B228: "Jump/Page",
    0x49B2B0: "Return",
    0x49B2CC: "On",
    0x49B338: "None",
    0x49B368: "Dialog",
    0x49B388: "Pointer",
    0x49B788: "Window",
    0x49B7C8: "Full",
    0x49B800: "Auto",
}


def _resolved_patches():
    rows = []
    for off, src, logical_target, context, note in JP_SLOT_PATCHES:
        if off in ASCII_FALLBACK_TARGETS:
            target = ASCII_FALLBACK_TARGETS[off]
            render_note = f"ASCII fallback for menu font-size test; Arabic-B logical='{logical_target}'"
        else:
            target = A(logical_target)
            render_note = f"Arabic-B logical='{logical_target}'"
        rows.append((off, src, target, context, f"{note}; {render_note}"))
    rows.extend(ENGLISH_FALLBACK_PATCHES)
    return rows


def main():
    base.GAME_EXE = r"C:\Program Files (x86)\Steam\steamapps\common\Kanon\Kanon.exe"
    base.PATCH_GAME_NAME = "Kanon Arabic-B JP slot no language labels"
    base.PATCH_VERSION = "0.6-ar-en-mixed-nolang"
    base.PATCHES = _resolved_patches()
    base.main()


if __name__ == "__main__":
    main()
