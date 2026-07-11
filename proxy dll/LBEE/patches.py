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
    (0x39457C, 'Close', 'Наз', 'bottom-right button', 'catalog=close; EN source slot; encoding=utf-8; budget 7', 'utf-8'),
    (0x395BCC, 'Defaults', 'Сброс', 'bottom-left button', 'catalog=defaults; EN source slot; encoding=utf-8; budget 11', 'utf-8'),
    (0x39B904, 'Skip', 'Про', 'Basic/Skip', 'catalog=skip; EN source slot; encoding=utf-8; budget 7', 'utf-8'),
    (0x39B968, 'Voice', 'Гол', 'Basic/Voice', 'catalog=voice; EN source slot; encoding=utf-8; budget 7', 'utf-8'),
    (0x39C38C, 'Disable', 'Отк', 'Basic/Rumble value', 'catalog=disable; EN source slot; encoding=utf-8; budget 7', 'utf-8'),
    (0x39BC10, 'Language', 'Язык', 'Text1/Language', 'catalog=language; EN source slot; encoding=utf-8; budget 11', 'utf-8'),
    (0x39BC44, 'Font', 'Шри', 'Text1/Font', 'catalog=font; EN source slot; encoding=utf-8; budget 7', 'utf-8'),
    (0x39BC24, 'Mincho', 'Мин', 'Text1/Font value', 'catalog=mincho; EN source slot; encoding=utf-8; budget 7', 'utf-8'),
    (0x39BC1C, 'Modern', 'Мод', 'Text1/Font value', 'catalog=modern; EN source slot; encoding=utf-8; budget 7', 'utf-8'),
    (0x39BCE0, 'Green', 'Зел', 'Text1/Color value', 'catalog=green; EN source slot; encoding=utf-8; budget 7', 'utf-8'),
    (0x39BCD8, 'Blue', 'Гол', 'Text1/Color value', 'catalog=blue; EN source slot; encoding=utf-8; budget 7', 'utf-8'),
    (0x39BCD0, 'Purple', 'Фио', 'Text1/Color value', 'catalog=purple; EN source slot; encoding=utf-8; budget 7', 'utf-8'),
    (0x39BCC4, 'Yellow', 'Жёл', 'Text1/Color value', 'catalog=yellow; EN source slot; encoding=utf-8; budget 7', 'utf-8'),
    (0x39C020, 'Text Speed', 'Скоро', 'Text2/Speed', 'catalog=text-speed; EN source slot; encoding=utf-8; budget 11', 'utf-8'),
    (0x39C018, 'Slow', 'Мал', 'Text2/Speed value', 'catalog=slow; EN source slot; encoding=utf-8; budget 7', 'utf-8'),
    (0x39C010, 'Fast', 'Бол', 'Text2/Speed value', 'catalog=fast; EN source slot; encoding=utf-8; budget 7', 'utf-8'),
    (0x39C03C, '0 sec/char', '0 сек./', 'Text2/Speed value', 'catalog=0-sec-char; EN source slot; encoding=utf-8; budget 11', 'utf-8'),
    (0x39C07C, '1 sec', '1 се', 'Text2/Wait value', 'catalog=1-sec; EN source slot; encoding=utf-8; budget 7', 'utf-8'),
    (0x39C074, '2 sec', '2 се', 'Text2/Wait value', 'catalog=2-sec; EN source slot; encoding=utf-8; budget 7', 'utf-8'),
    (0x39C06C, '3 sec', '3 се', 'Text2/Wait value', 'catalog=3-sec; EN source slot; encoding=utf-8; budget 7', 'utf-8'),
    (0x39C08C, 'Base Wait Time', 'Базовая ', 'Text2/Base', 'catalog=base-wait-time; EN source slot; encoding=utf-8; budget 15', 'utf-8'),
    (0x39C0D8, 'Master Volume', 'Общая гр', 'Sound/Master', 'catalog=master-volume; EN source slot; encoding=utf-8; budget 15', 'utf-8'),
    (0x39C374, 'C (Skip)', 'C (Проп', 'Keyboard/key label', 'catalog=c-skip; EN source slot; encoding=utf-8; budget 11', 'utf-8'),
    (0x39C380, 'Z (Rewind)', 'Z (Пере', 'Keyboard/key label', 'catalog=z-rewind; EN source slot; encoding=utf-8; budget 11', 'utf-8'),
    (0x39C3F0, 'Quick Save', 'Быстр', 'Keyboard/key label', 'catalog=quick-save; EN source slot; encoding=utf-8; budget 11', 'utf-8'),
    (0x39C3D0, 'Up Arrow', 'Стрел', 'Keyboard/key label', 'catalog=up-arrow; EN source slot; encoding=utf-8; budget 11', 'utf-8'),
    (0x39C3DC, 'Hide Window', 'Скрыт', 'Mouse/target', 'catalog=hide-window; EN source slot; encoding=utf-8; budget 11', 'utf-8'),
    (0x39C3B0, 'Rewind Once', 'Шаг на', 'Mouse/target', 'catalog=rewind-once; EN source slot; encoding=utf-8; budget 11', 'utf-8'),
    (0x39C394, 'Forward Once', 'Шаг впер', 'Mouse/target', 'catalog=forward-once; EN source slot; encoding=utf-8; budget 15', 'utf-8'),
    (0x39CB90, 'Window', 'Око', 'System/Window', 'catalog=window; EN source slot; encoding=utf-8; budget 7', 'utf-8'),
    (0x39CB84, 'Full Screen', 'Полно', 'System/FullScreen', 'catalog=full-screen; EN source slot; encoding=utf-8; budget 11', 'utf-8'),
    (0x39CB98, 'Screen Mode', 'Режим ', 'System/ScreenMode', 'catalog=screen-mode; EN source slot; encoding=utf-8; budget 11', 'utf-8'),
    (0x39CBAC, 'Auto', 'Авт', 'System/value', 'catalog=auto; EN source slot; encoding=utf-8; budget 7', 'utf-8'),
    (0x39CBB4, 'Window Size', 'Разме', 'System/WindowSize', 'catalog=window-size; EN source slot; encoding=utf-8; budget 11', 'utf-8'),
    (0x3949A8, 'Do you wish to save?', 'Сохранить иг', 'Save prompt', 'catalog=do-you-wish-to-save; EN source slot; encoding=utf-8; budget 23', 'utf-8'),
    (0x394820, 'Delete', 'Уда', 'Save menu button', 'catalog=delete; EN source slot; encoding=utf-8; budget 7', 'utf-8'),
    (0x394834, 'Latest', 'Пос', 'Save menu button', 'catalog=latest; EN source slot; encoding=utf-8; budget 7', 'utf-8'),
    (0x3943E4, 'Return to the title screen?', 'Вернуться в гл', 'Title return prompt', 'catalog=return-to-the-title-screen; EN source slot; encoding=utf-8; budget 27', 'utf-8'),
    (0x3941E0, 'Yes', 'Д', 'Global dialog button', 'catalog=yes; EN source slot; encoding=utf-8; budget 3', 'utf-8'),
    (0x3941E4, 'No', 'Н', 'Global dialog button', 'catalog=no; EN source slot; encoding=utf-8; budget 3', 'utf-8'),
    (0x398E1C, 'Are you sure you wish to quit the game?', 'Вы уверены, что хотит', 'Quit prompt (saved)', 'catalog=are-you-sure-you-wish-to-quit-the-game; EN source slot; encoding=utf-8; budget 39', 'utf-8'),
]
