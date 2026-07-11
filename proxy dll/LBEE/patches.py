#!/usr/bin/env python3
"""LBEE English-slot inventory for the LuckSystem Luca DLL GUI."""

GAME_EXE = 'LITBUS_WIN32.exe'
RVA_DELTA = 0x1000
RVA_MODE = 'pe'
PROXY_DLL = 'winmm'
ARCHITECTURE = 'x86'
PATCH_GAME_NAME = 'Little Busters! English Edition'
PATCH_VERSION = '3.29-lbee'

# Each entry: (raw_offset, source_text, target, context, note, encoding)
PATCHES = [
    (0x36350C, 'Close', '', 'bottom-right button', 'catalog=close; EN source slot; encoding=utf-8; budget 5', 'utf-8'),
    (0x395BCC, 'Defaults', '', 'bottom-left button', 'catalog=defaults; EN source slot; encoding=utf-8; budget 11', 'utf-8'),
    (0x39B904, 'Skip', '', 'Basic/Skip', 'catalog=skip; EN source slot; encoding=utf-8; budget 7', 'utf-8'),
    (0x398714, 'Voice', '', 'Basic/Voice', 'catalog=voice; EN source slot; encoding=utf-16-le; budget 10', 'utf-16-le'),
    (0x39C38C, 'Disable', '', 'Basic/Rumble value', 'catalog=disable; EN source slot; encoding=utf-8; budget 7', 'utf-8'),
    (0x3992D6, 'Language', '', 'Text1/Language', 'catalog=language; EN source slot; encoding=utf-8; budget 8', 'utf-8'),
    (0x392860, 'Font', '', 'Text1/Font', 'catalog=font; EN source slot; encoding=utf-16-le; budget 8', 'utf-16-le'),
    (0x39BC24, 'Mincho', '', 'Text1/Font value', 'catalog=mincho; EN source slot; encoding=utf-8; budget 7', 'utf-8'),
    (0x39BC1C, 'Modern', '', 'Text1/Font value', 'catalog=modern; EN source slot; encoding=utf-8; budget 7', 'utf-8'),
    (0x39BCE0, 'Green', '', 'Text1/Color value', 'catalog=green; EN source slot; encoding=utf-8; budget 7', 'utf-8'),
    (0x39BCD8, 'Blue', '', 'Text1/Color value', 'catalog=blue; EN source slot; encoding=utf-8; budget 7', 'utf-8'),
    (0x39BCD0, 'Purple', '', 'Text1/Color value', 'catalog=purple; EN source slot; encoding=utf-8; budget 7', 'utf-8'),
    (0x39BCC4, 'Yellow', '', 'Text1/Color value', 'catalog=yellow; EN source slot; encoding=utf-8; budget 7', 'utf-8'),
    (0x39C020, 'Text Speed', '', 'Text2/Speed', 'catalog=text-speed; EN source slot; encoding=utf-8; budget 11', 'utf-8'),
    (0x39ADB8, 'Slow', '', 'Text2/Speed value', 'catalog=slow; EN source slot; encoding=utf-16-le; budget 8', 'utf-16-le'),
    (0x39C010, 'Fast', '', 'Text2/Speed value', 'catalog=fast; EN source slot; encoding=utf-8; budget 7', 'utf-8'),
    (0x39C03C, '0 sec/char', '', 'Text2/Speed value', 'catalog=0-sec-char; EN source slot; encoding=utf-8; budget 11', 'utf-8'),
    (0x39C07C, '1 sec', '', 'Text2/Wait value', 'catalog=1-sec; EN source slot; encoding=utf-8; budget 7', 'utf-8'),
    (0x39C02E, '2 sec', '', 'Text2/Wait value', 'catalog=2-sec; EN source slot; encoding=utf-8; budget 5', 'utf-8'),
    (0x39C06C, '3 sec', '', 'Text2/Wait value', 'catalog=3-sec; EN source slot; encoding=utf-8; budget 7', 'utf-8'),
    (0x39BEE2, 'Base Wait Time', '', 'Text2/Base', 'catalog=base-wait-time; EN source slot; encoding=utf-16-le; budget 28', 'utf-16-le'),
    (0x39C0D8, 'Master Volume', '', 'Sound/Master', 'catalog=master-volume; EN source slot; encoding=utf-8; budget 15', 'utf-8'),
    (0x39C374, 'C (Skip)', '', 'Keyboard/key label', 'catalog=c-skip; EN source slot; encoding=utf-8; budget 11', 'utf-8'),
    (0x39C380, 'Z (Rewind)', '', 'Keyboard/key label', 'catalog=z-rewind; EN source slot; encoding=utf-8; budget 11', 'utf-8'),
    (0x39C3F0, 'Quick Save', '', 'Keyboard/key label', 'catalog=quick-save; EN source slot; encoding=utf-8; budget 11', 'utf-8'),
    (0x39C3D0, 'Up Arrow', '', 'Keyboard/key label', 'catalog=up-arrow; EN source slot; encoding=utf-8; budget 11', 'utf-8'),
    (0x39C3DC, 'Hide Window', '', 'Mouse/target', 'catalog=hide-window; EN source slot; encoding=utf-8; budget 11', 'utf-8'),
    (0x39C3B0, 'Rewind Once', '', 'Mouse/target', 'catalog=rewind-once; EN source slot; encoding=utf-8; budget 11', 'utf-8'),
    (0x39C394, 'Forward Once', '', 'Mouse/target', 'catalog=forward-once; EN source slot; encoding=utf-8; budget 15', 'utf-8'),
    (0x3962DC, 'Window', '', 'System/Window', 'catalog=window; EN source slot; encoding=utf-16-le; budget 12', 'utf-16-le'),
    (0x39CB84, 'Full Screen', '', 'System/FullScreen', 'catalog=full-screen; EN source slot; encoding=utf-8; budget 11', 'utf-8'),
    (0x39CB98, 'Screen Mode', '', 'System/ScreenMode', 'catalog=screen-mode; EN source slot; encoding=utf-8; budget 11', 'utf-8'),
    (0x39861A, 'Auto', '', 'System/value', 'catalog=auto; EN source slot; encoding=utf-8; budget 4', 'utf-8'),
    (0x39CBB4, 'Window Size', '', 'System/WindowSize', 'catalog=window-size; EN source slot; encoding=utf-8; budget 11', 'utf-8'),
    (0x3949A8, 'Do you wish to save?', '', 'Save prompt', 'catalog=do-you-wish-to-save; EN source slot; encoding=utf-8; budget 23', 'utf-8'),
    (0x394820, 'Delete', '', 'Save menu button', 'catalog=delete; EN source slot; encoding=utf-8; budget 7', 'utf-8'),
    (0x394834, 'Latest', '', 'Save menu button', 'catalog=latest; EN source slot; encoding=utf-8; budget 7', 'utf-8'),
    (0x3943E4, 'Return to the title screen?', '', 'Title return prompt', 'catalog=return-to-the-title-screen; EN source slot; encoding=utf-8; budget 27', 'utf-8'),
    (0x3B06C8, 'Yes', '', 'Global dialog button', 'catalog=yes; EN source slot; encoding=utf-16-le; budget 6', 'utf-16-le'),
    (0x363A6C, 'No', '', 'Global dialog button', 'catalog=no; EN source slot; encoding=utf-16-le; budget 4', 'utf-16-le'),
    (0x398E1C, 'Are you sure you wish to quit the game?', '', 'Quit prompt (saved)', 'catalog=are-you-sure-you-wish-to-quit-the-game; EN source slot; encoding=utf-8; budget 39', 'utf-8'),
]

