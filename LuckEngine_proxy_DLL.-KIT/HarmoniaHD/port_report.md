# Harmonia HD proxy DLL port report

- Source table: `C:\Users\jeuxpc\Documents\GitHub\LuckSystem-2.3.2-Yoremi-Update\LuckEngine_proxy_DLL.-KIT\Kanon\patches.py`
- Target exe: `C:\Program Files (x86)\Steam\steamapps\common\Harmonia Full HD Edition\HarmoniaFHD.exe`
- RVA delta: `0xE00`
- Selected entries: `101`
- Skipped entries: `11`

## Selected Ambiguous Entries

| Kanon raw | Harmonia raw | Hits | Distance | Context | Source |
|---:|---:|---:|---:|---|---|
| `0x4874A4` | `0x48C4D4` | 6 | `0x0` | bottom-right button | `Close` |
| `0x499208` | `0x49DBC0` | 5 | `0x18` | Basic/Skip | `Skip` |
| `0x499320` | `0x49DC88` | 18 | `0x1C` | Basic/Voice | `Voice` |
| `0x499450` | `0x49DD58` | 2 | `0x1` | Basic/Cursor | `Initial Cursor Position` |
| `0x499490` | `0x49DD98` | 2 | `0x2` | Basic/Rumble | `Controller Rumble Function` |
| `0x499470` | `0x49DD78` | 3 | `0x1` | Basic/Rumble value | `Disable` |
| `0x49A4D8` | `0x49EE00` | 9 | `0x2C` | Text1/Language | `Language` |
| `0x49A598` | `0x49EE94` | 3 | `0x5` | Text1/Font | `Font` |
| `0x49A5FC` | `0x49EEEC` | 5 | `0x2` | Text1/Window Transp value | `Clear` |
| `0x49A6D8` | `0x49EFEC` | 2 | `0x10` | Text1/Color value | `Blue` |
| `0x49ABB0` | `0x49F4A0` | 3 | `0x11` | Text2/Wait | `Wait Time Per Character` |
| `0x49ABC8` | `0x49F4E8` | 2 | `0x11` | Text2/Wait value | `0 sec` |
| `0x49ABD8` | `0x49F4F8` | 2 | `0x9` | Text2/Wait value | `1 sec` |
| `0x49AC08` | `0x49F548` | 3 | `0x12` | Text2/Base | `Base Wait Time` |
| `0x49AC5C` | `0x49F558` | 23 | `0x3B` | Sound tab | `Sound` |
| `0x49B02C` | `0x49F8EC` | 2 | `0xA` | Mouse tab | `Mouse` |
| `0x49B0E8` | `0x49F99D` | 3 | `0xD` | Mouse/binding | `Right Click` |
| `0x49B108` | `0x49F998` | 2 | `0x2F` | Mouse/binding | `Left+Right Click` |
| `0x49B310` | `0x49FBA0` | 2 | `0x4` | Mouse/Gestures | `Gestures` |
| `0x49B76C` | `0x4A001C` | 23 | `0xE` | System tab | `System` |
| `0x49B7B0` | `0x4A0024` | 29 | `0x47` | System/Window | `Window` |
| `0x49B7B8` | `0x4A0078` | 2 | `0x6` | System/FullScreen | `Full Screen` |
| `0x49B820` | `0x4A00A4` | 7 | `0x1B` | System/value | `Auto` |
| `0x488B2C` | `0x48DB5C` | 10 | `0x0` | Save menu button | `Delete` |
| `0x487458` | `0x48C488` | 2 | `0x0` | Language switcher | `English` |
| `0x487478` | `0x48C4A8` | 3 | `0x0` | Global dialog button | `Yes` |
| `0x48748C` | `0x48C4BC` | 96 | `0x0` | Global dialog button | `No` |
| `0x495C58` | `0x49A7C8` | 2 | `0xBC` | Quit prompt (saved) | `Are you sure you wish to quit the game?` |

## Skipped Entries

| Kanon raw | Hits | Context | Source | Reason |
|---:|---:|---|---|---|
| `0x4991A0` | 0 | Basic/Shortcut | `Shortcut Menu` | not found in HarmoniaFHD.exe |
| `0x499144` | 3 | Basic/Shortcut value | `Hide` | ambiguous, no local match near predicted offset; hits=0x488d78, 0x48dea8, 0x49f8f8 |
| `0x499150` | 5 | Basic/Shortcut value | `Display` | ambiguous, no local match near predicted offset; hits=0x499bd3, 0x4a03eb, 0x4e9cf3, 0x4e9d13, 0x4e9d28 |
| `0x499248` | 0 | Basic/Position | `Position of Choices` | not found in HarmoniaFHD.exe |
| `0x499230` | 0 | Basic/Position value | `Bottom` | not found in HarmoniaFHD.exe |
| `0x499238` | 0 | Basic/Position value | `Center` | not found in HarmoniaFHD.exe |
| `0x499398` | 0 | Basic/Date | `Display Date` | not found in HarmoniaFHD.exe |
| `0x49A630` | 0 | Text1/Transp target | `Only Choices` | not found in HarmoniaFHD.exe |
| `0x49AD30` | 0 | Voice tab young label | `	(Young)	` | not found in HarmoniaFHD.exe |
| `0x49B350` | 0 | Mouse/Snap target | `Dialog and Choices` | not found in HarmoniaFHD.exe |
| `0x499358` | 75 | Read text color toggle | `On` | ambiguous, no local match near predicted offset; hits=0x4967, 0x1ac2d, 0x32a4f, 0x47fec, 0x49f25, 0x52f1c, 0x707aa, 0x84ad4 |
