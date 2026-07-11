<script>
  import { onMount, onDestroy, tick } from 'svelte';
  import { EventsOn, EventsOff, ClipboardGetText, ClipboardSetText } from '../wailsjs/runtime/runtime.js';
  import {
    GetLuckSystemPath,
    SetLuckSystemPath,
    ScanGameData,
    SelectPakFile,
    SelectFile,
    SelectDirectory,
    SelectSaveFile,
    StopProcess,
    ScriptDecompile,
    ScriptCompile,
    SiglusLucaBridge,
    PakExtract,
    BGMOVIEExtract,
    MusicPakExtract,
    VoicePakExtract,
    AudioConvert,
    PakReplace,
    PakFontExtract,
    PakFontReplace,
    FontExtract,
    FontEdit,
    ImageExport,
    ImageImport,
    ImageBatchExport,
    ImageBatchImport,
    DialogueDetectFormat,
    DialogueExtractFile,
    DialogueExtractBatch,
    DialogueImportFile,
    DialogueImportBatch,
    VietnameseFontPatch,
    SupportsLucaMenuDLL,
    ScanLucaMenuKit,
    LucaMenuGenerate,
    SelectScriptTxtFile,
    SelectTsvFile,
    SelectSaveTsvFile,
    SelectSaveScriptFile
  } from '../wailsjs/go/main/App.js';

  // ===== State =====
  let selectedOp = 'decompile';
  let running = false;
  let consoleLines = [];
  let consoleEl;
  let consoleMenuVisible = false;
  let consoleMenuX = 0;
  let consoleMenuY = 0;
  let lsPath = '';
  let lucaMenuDllAvailable = false;
  let uiLanguage = 'fr';

  function t(fr, en, language = uiLanguage) {
    return language === 'en' ? en : fr;
  }

  function setUiLanguage(language) {
    uiLanguage = language === 'en' ? 'en' : 'fr';
    try { localStorage.setItem('lucksystem-ui-language', uiLanguage); } catch (_) {}
  }

  // --- Script fields ---
  let pakFile = '';
  let opcodeFile = '';
  let pluginFile = '';
  let charsetVal = 'UTF-8';
  let gameName = '';
  let gamePresets = [];
  let selectedPreset = '';
  let outputDir = '';
  let importDir = '';
  let outputPak = '';

  // --- Siglus -> Luca bridge fields ---
  let siglusLucaLucaDir = '';
  let siglusLucaSiglusDir = '';
  let siglusLucaOutput = '';
  let siglusLucaTargetCol = 2;

  // --- PAK fields ---
  let pakExtSource = '';
  let pakExtOutput = '';
  let pakRepSource = '';
  let pakRepListFile = '';
  let pakRepInput = '';
  let pakRepOutput = '';
  let pakRepUseList = true; // mode par défaut : fichier liste

  // --- BGMOVIE / Video fields ---
  let bgMoviePak = '';
  let bgMovieOutput = '';

  // --- PAK Audio fields ---
  let musicPak = '';
  let musicOutput = '';
  let musicToMp3 = false;
  let voicePak = '';
  let voiceOutput = '';
  let voiceToMp3 = false;
  let audioConvInput = '';
  let audioConvOutput = '';
  let audioConvDirection = 'mp3'; // 'mp3' | 'native'

  // --- PAK Font fields ---
  let pakFontExtSource = '';
  let pakFontExtCharset = 'UTF-8';
  let pakFontExtOutput = '';
  let pakFontRepSource = '';
  let pakFontRepCharset = 'UTF-8';
  let pakFontRepListFile = '';
  let pakFontRepInput = '';
  let pakFontRepSingleFile = '';
  let pakFontRepSingleName = '';
  let pakFontRepAliasFrom = 'info30';
  let pakFontRepAliasTo = 'info32';
  let pakFontRepOutput = '';
  let pakFontRepMode = 'list'; // 'list' | 'dir' | 'single' | 'alias'

  // --- Font Extract ---
  let fontExtCz = '';
  let fontExtInfo = '';
  let fontExtPng = '';
  let fontExtCharset = '';

  // --- Font Edit ---
  let fontEditCz = '';
  let fontEditInfo = '';
  let fontEditTtf = '';
  let fontEditOutCz = '';
  let fontEditOutInfo = '';
  let fontEditCharsetFile = '';
  let fontEditMode = 'append'; // 'redraw' | 'append' | 'insert'
  let fontEditIndex = 0;
  let fontEditArabicMetrics = false;
  let fontEditMetricSetYEnabled = false;
  let fontEditMetricSetY = 0;
  let fontEditMetricYOffset = 0;
  let fontEditMetricXOffset = 0;
  let fontEditMetricWOffset = 0;
  let fontEditArabicConnectorBleed = 0;

  // --- Vietnamese Font Patch ---
  let vietFontRoot = '';
  let vietCharsetFile = '';
  let vietTtfFile = '';
  let vietOutputDir = '';
  let vietSlot = 'en';
  let vietFamily = 'GOTHIC1';
  let vietYMinus2 = false;
  let vietYMinus1 = false;
  let vietY0 = false;
  let vietY1 = false;
  let vietY2 = true;
  let vietY3 = false;
  let vietRedrawLatin = false;

  // --- Luca Menu DLL ---
  let lucaInventory = null;
  let lucaGame = '';
  let lucaSlot = 'en';
  let lucaPatchName = '';
  let lucaPatchVersion = '0.1-gui';
  let lucaExe = '';
  let lucaOutputDir = '';
  let lucaBuildDll = true;
  let lucaProxyChoice = 'version';
  let lucaCustomPatch = '';
  let lucaFillMode = 'fr-safe';
  let lucaSearch = '';
  let lucaEntries = [];

  // --- Image Export ---
  let imgExpBatch = false;
  let imgExpInput = '';
  let imgExpOutput = '';

  // --- Image Import ---
  let imgImpBatch = false;
  let imgImpSource = '';
  let imgImpInput = '';
  let imgImpOutput = '';
  let imgImpFill = false;

  // --- Dialogue Extract ---
  let dlgExtBatch = false;
  let dlgExtInput = '';
  let dlgExtOutput = '';
  let dlgExtLang1 = false;
  let dlgExtLang2 = true;   // default: Lang 2 (typically ENG in AIR)
  let dlgExtLang3 = false;
  let dlgExtLang4 = false;
  let dlgExtDetectedFmt = '';
  let dlgExtMaxCols = 0;

  // --- Dialogue Import ---
  let dlgImpBatch = false;
  let dlgImpScript = '';
  let dlgImpTsv = '';
  let dlgImpOutput = '';
  let dlgImpTargetCol = 2;  // default: Lang 2

  // ===== Operations list =====
  const operations = [
    { id: '_s1', label: 'SCRIPT', section: true },
    { id: 'decompile', label: 'Script Decompile', labelFr: 'Décompiler les scripts' },
    { id: 'compile', label: 'Script Compile', labelFr: 'Compiler les scripts' },
    { id: 'siglus_luca', label: 'Siglus -> Luca' },
    { id: '_s2', label: 'PAK (CG)', section: true },
    { id: 'pak_cg_extract', label: 'CG Extract', labelFr: 'Extraire les CG' },
    { id: 'pak_cg_replace', label: 'CG Replace', labelFr: 'Remplacer les CG' },
    { id: '_s2v', label: 'PAK (Video)', labelFr: 'PAK (Vidéo)', section: true },
    { id: 'bgmovie_extract', label: 'BGMOVIE Extract', labelFr: 'Extraire BGMOVIE' },
    { id: '_s2a', label: 'PAK (Audio)', section: true },
    { id: 'music_extract', label: 'Music Extract', labelFr: 'Extraire la musique' },
    { id: 'voice_extract', label: 'Voice Extract', labelFr: 'Extraire les voix' },
    { id: 'audio_convert', label: 'Ogg / MP3 Convert', labelFr: 'Convertir Ogg / MP3' },
    { id: '_s2b', label: 'PAK (Font)', labelFr: 'PAK (Police)', section: true },
    { id: 'pak_font_extract', label: 'Font Extract', labelFr: 'Extraire la police' },
    { id: 'pak_font_replace', label: 'Font Replace', labelFr: 'Remplacer la police' },
    { id: '_s3', label: 'FONT', labelFr: 'POLICE', section: true },
    { id: 'font_extract', label: 'Font Extract', labelFr: 'Extraire la police' },
    { id: 'font_edit', label: 'Font Edit', labelFr: 'Modifier la police' },
    { id: '_s3b', label: 'VIET FONT', labelFr: 'POLICE VIET', section: true },
    { id: 'viet_font_patch', label: 'AIR / SG Patch' },
    { id: '_s3c', label: 'DLL HOOK', section: true },
    { id: 'luca_menu_dll', label: 'Luca Menu DLL' },
    { id: '_s4', label: 'IMAGE', section: true },
    { id: 'image_export', label: 'Image Export', labelFr: 'Exporter les images' },
    { id: 'image_import', label: 'Image Import', labelFr: 'Importer les images' },
    { id: '_s6', label: 'DIALOGUE', labelFr: 'DIALOGUES', section: true },
    { id: 'dlg_extract', label: 'Extract Dialogues', labelFr: 'Extraire les dialogues' },
    { id: 'dlg_import', label: 'Import Dialogues', labelFr: 'Importer les dialogues' },
    { id: '_s5', label: '', section: true },
    { id: 'about', label: 'About', labelFr: 'À propos' },
  ];

  // ===== Console =====
  // Batched console updates for performance (flush every 80ms instead of per-line)
  let pendingLines = [];
  let flushTimer = null;

  function addLine(text) {
    let cls = '';
    if (text.includes('[OK]')) cls = 'line-ok';
    else if (text.includes('[ERROR]') || text.includes('Panic') || text.includes('Error')) cls = 'line-err';
    else if (text.startsWith('═') || text.startsWith('─')) cls = 'line-sep';
    else if (text.startsWith('>')) cls = 'line-cmd';
    pendingLines.push({ text, cls });
    if (!flushTimer) {
      flushTimer = setTimeout(flushConsole, 80);
    }
  }
  function flushConsole() {
    if (pendingLines.length > 0) {
      consoleLines = [...consoleLines, ...pendingLines];
      pendingLines = [];
      // Cap at 2000 lines to prevent memory bloat
      if (consoleLines.length > 2000) consoleLines = consoleLines.slice(-1500);
      tick().then(() => { if (consoleEl) consoleEl.scrollTop = consoleEl.scrollHeight; });
    }
    flushTimer = null;
  }
  function clearConsole() { consoleLines = []; pendingLines = []; }

  function getConsoleText() {
    return [...consoleLines, ...pendingLines].map(line => line.text).join('\n');
  }

  function openConsoleMenu(event) {
    event.preventDefault();
    const menuWidth = 190;
    const menuHeight = 116;
    consoleMenuX = Math.min(event.clientX, window.innerWidth - menuWidth - 8);
    consoleMenuY = Math.min(event.clientY, window.innerHeight - menuHeight - 8);
    consoleMenuVisible = true;
  }

  function closeConsoleMenu() {
    consoleMenuVisible = false;
  }

  async function setClipboardText(text) {
    if (!text) return false;
    try {
      const ok = await ClipboardSetText(text);
      if (ok) return true;
    } catch (e) {
      // Browser fallback below.
    }
    if (typeof navigator !== 'undefined' && navigator.clipboard?.writeText) {
      await navigator.clipboard.writeText(text);
      return true;
    }
    return false;
  }

  async function getClipboardText() {
    try {
      return await ClipboardGetText();
    } catch (e) {
      if (typeof navigator !== 'undefined' && navigator.clipboard?.readText) {
        return await navigator.clipboard.readText();
      }
    }
    return '';
  }

  async function copyConsoleSelection() {
    const selection = typeof window !== 'undefined' ? window.getSelection()?.toString() || '' : '';
    await setClipboardText(selection || getConsoleText());
    closeConsoleMenu();
  }

  async function copyConsoleAll() {
    await setClipboardText(getConsoleText());
    closeConsoleMenu();
  }

  async function pasteConsoleClipboard() {
    const text = await getClipboardText();
    if (text) {
      text.replace(/\r\n/g, '\n').replace(/\r/g, '\n').split('\n').forEach(line => addLine(line));
    }
    closeConsoleMenu();
  }

  function handleConsoleKeydown(event) {
    if (!(event.ctrlKey || event.metaKey)) return;
    const key = event.key.toLowerCase();
    if (key === 'c') {
      event.preventDefault();
      copyConsoleSelection();
    } else if (key === 'v') {
      event.preventDefault();
      pasteConsoleClipboard();
    }
  }

  function handleWindowKeydown(event) {
    if (event.key === 'Escape') closeConsoleMenu();
  }

  onMount(async () => {
    try { setUiLanguage(localStorage.getItem('lucksystem-ui-language') || 'fr'); } catch (_) {}
    window.addEventListener('click', closeConsoleMenu);
    window.addEventListener('keydown', handleWindowKeydown);
    EventsOn('log', (msg) => addLine(msg));
    lucaMenuDllAvailable = await SupportsLucaMenuDLL();
    lsPath = await GetLuckSystemPath();
    if (lsPath) {
      addLine('LuckSystem 2.3.2 - Yoremi fork v3.30');
      addLine('Executable: ' + lsPath);
      // Scan data/ folder for game presets
      gamePresets = (await ScanGameData()) || [];
      if (gamePresets.length > 0) {
        addLine('Found ' + gamePresets.length + ' game preset(s): ' + gamePresets.map(p => p.name).join(', '));
      }
    } else {
      addLine('[ERROR] lucksystem.exe not found!');
      addLine('Place lucksystem.exe next to the GUI, or click "Locate" below.');
    }
    addLine('Ready.');
  });
  onDestroy(() => {
    EventsOff('log');
    window.removeEventListener('click', closeConsoleMenu);
    window.removeEventListener('keydown', handleWindowKeydown);
  });

  // ===== Browse helpers =====
  async function browsePak() { const f = await SelectPakFile(); if (f) pakFile = f; }
  async function browseOpcode() { const f = await SelectFile('Select Opcode (.txt)', '*.txt', 'Opcode files'); if (f) { opcodeFile = f; selectedPreset = ''; } }
  async function browsePlugin() { const f = await SelectFile('Select Plugin (.py)', '*.py', 'Python plugins'); if (f) { pluginFile = f; selectedPreset = ''; } }
  async function browseOutputDir() { const d = await SelectDirectory('Select output directory'); if (d) outputDir = d; }
  async function browseImportDir() { const d = await SelectDirectory('Select translated scripts directory'); if (d) importDir = d; }
  async function browseOutputPak() { const f = await SelectSaveFile('Save output PAK', 'SCRIPT_FR.PAK', '*.PAK;*.pak', 'PAK files'); if (f) outputPak = f; }
  async function browseSiglusLucaLucaDir() { const d = await SelectDirectory('Select Luca decompiled scripts folder'); if (d) siglusLucaLucaDir = d; }
  async function browseSiglusLucaSiglusDir() { const d = await SelectDirectory('Select Siglus Full folder'); if (d) siglusLucaSiglusDir = d; }
  async function browseSiglusLucaOutput() { const d = await SelectDirectory('Select patched Luca output folder'); if (d) siglusLucaOutput = d; }

  function applyPreset(presetName) {
    selectedPreset = presetName;
    if (!presetName) { opcodeFile = ''; pluginFile = ''; gameName = ''; return; }
    const p = gamePresets.find(g => g.name === presetName);
    if (p) { opcodeFile = p.opcodeFile; pluginFile = p.pluginFile || ''; gameName = p.gameFlag || ''; }
  }

  async function browsePakExtSource() { const f = await SelectPakFile(); if (f) pakExtSource = f; }
  async function browsePakExtOutput() { const d = await SelectDirectory('Select extraction output'); if (d) pakExtOutput = d; }
  async function browsePakRepSource() { const f = await SelectPakFile(); if (f) pakRepSource = f; }
  async function browsePakRepListFile() { const f = await SelectFile(t('Sélectionner le fichier liste (_list.txt)', 'Select list file (_list.txt)'), '*.txt', t('Fichiers liste', 'List files')); if (f) pakRepListFile = f; }
  async function browsePakRepInput() { const d = await SelectDirectory('Select folder with modified files'); if (d) pakRepInput = d; }
  async function browsePakRepOutput() { const f = await SelectSaveFile('Save output PAK', 'FONT.out.PAK', '*.PAK;*.pak', 'PAK files'); if (f) pakRepOutput = f; }

  async function browseBgMoviePak() { const f = await SelectPakFile(); if (f) bgMoviePak = f; }
  async function browseBgMovieOutput() { const d = await SelectDirectory('Select BGMOVIE output folder'); if (d) bgMovieOutput = d; }

  async function browseMusicPak() { const f = await SelectPakFile(); if (f) musicPak = f; }
  async function browseMusicOutput() { const d = await SelectDirectory('Select MUSIC output folder'); if (d) musicOutput = d; }
  async function browseVoicePak() { const f = await SelectPakFile(); if (f) voicePak = f; }
  async function browseVoiceOutput() { const d = await SelectDirectory('Select VOICE output folder'); if (d) voiceOutput = d; }
  async function browseAudioConvInput() {
    const d = await SelectDirectory(audioConvDirection === 'mp3' ? 'Select native Ogg folder' : 'Select MP3 folder');
    if (d) audioConvInput = d;
  }
  async function browseAudioConvOutput() { const d = await SelectDirectory('Select converted audio output folder'); if (d) audioConvOutput = d; }

  async function browsePakFontExtSource() { const f = await SelectPakFile(); if (f) pakFontExtSource = f; }
  async function browsePakFontExtOutput() { const d = await SelectDirectory(t('Dossier d\'extraction', 'Extraction folder')); if (d) pakFontExtOutput = d; }
  async function browsePakFontRepSource() { const f = await SelectPakFile(); if (f) pakFontRepSource = f; }
  async function browsePakFontRepListFile() { const f = await SelectFile(t('Sélectionner le fichier liste (_list.txt)', 'Select list file (_list.txt)'), '*.txt', t('Fichiers liste', 'List files')); if (f) pakFontRepListFile = f; }
  async function browsePakFontRepInput() { const d = await SelectDirectory(t('Dossier des fichiers modifiés', 'Modified files folder')); if (d) pakFontRepInput = d; }
  async function browsePakFontRepSingleFile() {
    const f = await SelectFile(t('Sélectionner le fichier à remplacer', 'Select file to replace'), '*.*', t('Tous les fichiers', 'All files'));
    if (f) {
      pakFontRepSingleFile = f;
      if (!pakFontRepSingleName) pakFontRepSingleName = f.split(/[\\/]/).pop();
    }
  }
  async function browsePakFontRepOutput() { const f = await SelectSaveFile('Save output PAK', 'FONT.out.PAK', '*.PAK;*.pak', 'PAK files'); if (f) pakFontRepOutput = f; }

  async function browseFontExtCz() { const f = await SelectFile('Select font CZ file', '*.*', 'Font CZ files'); if (f) fontExtCz = f; }
  async function browseFontExtInfo() { const f = await SelectFile('Select info file', '*.*', 'Info files'); if (f) fontExtInfo = f; }
  async function browseFontExtPng() { const f = await SelectSaveFile('Save font PNG', 'font.png', '*.png', 'PNG Images'); if (f) fontExtPng = f; }
  async function browseFontExtCharset() { const f = await SelectSaveFile('Save charset TXT', 'charset.txt', '*.txt', 'Text files'); if (f) fontExtCharset = f; }
  async function browseFontEditCz() { const f = await SelectFile('Select source CZ', '*.*', 'Font CZ files'); if (f) fontEditCz = f; }
  async function browseFontEditInfo() { const f = await SelectFile('Select source info', '*.*', 'Info files'); if (f) fontEditInfo = f; }
  async function browseFontEditTtf() { const f = await SelectFile('Select TTF font', '*.ttf;*.otf', 'Font files'); if (f) fontEditTtf = f; }
  async function browseFontEditOutCz() { const d = await SelectDirectory(t('Dossier de sortie pour le CZ modifié', 'Output folder for modified CZ')); if (d) fontEditOutCz = d + '\\'; }
  async function browseFontEditOutInfo() { const d = await SelectDirectory(t('Dossier de sortie pour le fichier info', 'Output folder for info file')); if (d) fontEditOutInfo = d + '\\'; }
  async function browseFontEditCharset() { const f = await SelectFile('Select charset file', '*.txt', 'Text files'); if (f) fontEditCharsetFile = f; }

  async function browseVietFontRoot() { const d = await SelectDirectory('Select AIR / Planetarian SG files folder'); if (d) vietFontRoot = d; }
  async function browseVietCharset() { const f = await SelectFile('Select full Vietnamese charset (134 chars)', '*.txt', 'Text files'); if (f) vietCharsetFile = f; }
  async function browseVietTtf() { const f = await SelectFile('Select Vietnamese-capable TTF/OTF', '*.ttf;*.otf', 'Font files'); if (f) vietTtfFile = f; }
  async function browseVietOutput() { const d = await SelectDirectory('Select output folder'); if (d) vietOutputDir = d; }

  async function browseImgExpInput() {
    if (imgExpBatch) { const d = await SelectDirectory('Select CZ folder'); if (d) imgExpInput = d; }
    else { const f = await SelectFile('Select CZ file', '*.*', 'CZ image files'); if (f) imgExpInput = f; }
  }
  async function browseImgExpOutput() {
    if (imgExpBatch) { const d = await SelectDirectory('Select PNG output folder'); if (d) imgExpOutput = d; }
    else { const f = await SelectSaveFile('Save PNG', 'output.png', '*.png', 'PNG Images'); if (f) imgExpOutput = f; }
  }
  async function browseImgImpSource() {
    if (imgImpBatch) { const d = await SelectDirectory('Select original CZ folder'); if (d) imgImpSource = d; }
    else { const f = await SelectFile('Select original CZ file', '*.*', 'CZ files'); if (f) imgImpSource = f; }
  }
  async function browseImgImpInput() {
    if (imgImpBatch) { const d = await SelectDirectory('Select PNG folder'); if (d) imgImpInput = d; }
    else { const f = await SelectFile('Select PNG file', '*.png', 'PNG Images'); if (f) imgImpInput = f; }
  }
  async function browseImgImpOutput() {
    if (imgImpBatch) { const d = await SelectDirectory('Select output CZ folder'); if (d) imgImpOutput = d; }
    else { const f = await SelectSaveFile('Save CZ', 'output.cz', '*.*', 'All files'); if (f) imgImpOutput = f; }
  }

  async function locateLuckSystem() {
    lsPath = await SetLuckSystemPath();
    if (lsPath) {
      addLine('Executable set: ' + lsPath);
      gamePresets = (await ScanGameData()) || [];
      if (gamePresets.length > 0) addLine('Found ' + gamePresets.length + ' game preset(s): ' + gamePresets.map(p => p.name).join(', '));
    }
  }

  async function stopProcess() {
    await StopProcess();
  }

  // ===== Actions =====
  async function run(fn) {
    if (running) return;
    running = true;
    try { await fn(); } catch (e) { addLine('[ERROR] ' + e); }
    running = false;
  }

  function startDecompile() { run(() => ScriptDecompile(pakFile, opcodeFile, pluginFile, charsetVal, outputDir, gameName)); }
  function startCompile() { run(() => ScriptCompile(pakFile, opcodeFile, pluginFile, charsetVal, importDir, outputPak, gameName)); }
  function startSiglusLucaBridge() { run(() => SiglusLucaBridge(siglusLucaLucaDir, siglusLucaSiglusDir, siglusLucaOutput, siglusLucaTargetCol)); }
  function startPakExtract() { run(() => PakExtract(pakExtSource, pakExtOutput)); }
  function startBgMovieExtract() { run(() => BGMOVIEExtract(bgMoviePak, bgMovieOutput)); }
  function startMusicExtract() { run(() => MusicPakExtract(musicPak, musicOutput, musicToMp3)); }
  function startVoiceExtract() { run(() => VoicePakExtract(voicePak, voiceOutput, voiceToMp3)); }
  function startAudioConvert() { run(() => AudioConvert(audioConvInput, audioConvOutput, audioConvDirection)); }
  function startPakReplace() {
    const listArg = pakRepUseList ? pakRepListFile : '';
    const dirArg  = pakRepUseList ? '' : pakRepInput;
    run(() => PakReplace(pakRepSource, dirArg, listArg, pakRepOutput));
  }
  function startPakFontExtract() { run(() => PakFontExtract(pakFontExtSource, pakFontExtCharset, pakFontExtOutput)); }
  function startPakFontReplace() {
    const listArg = pakFontRepMode === 'list' ? pakFontRepListFile : '';
    const dirArg  = pakFontRepMode === 'dir' ? pakFontRepInput : '';
    const fileArg = pakFontRepMode === 'single' ? pakFontRepSingleFile : '';
    const nameArg = pakFontRepMode === 'single' ? pakFontRepSingleName : '';
    const aliasFromArg = pakFontRepMode === 'alias' ? pakFontRepAliasFrom : '';
    const aliasToArg = pakFontRepMode === 'alias' ? pakFontRepAliasTo : '';
    run(() => PakFontReplace(pakFontRepSource, pakFontRepCharset, dirArg, listArg, fileArg, nameArg, aliasFromArg, aliasToArg, pakFontRepOutput));
  }
  function startFontExtract() { run(() => FontExtract(fontExtCz, fontExtInfo, fontExtPng, fontExtCharset)); }
  function setFontEditArabicPreset(checked) {
    fontEditArabicMetrics = checked;
  }
  function startFontEdit() {
    const redraw  = fontEditMode === 'redraw';
    const append  = fontEditMode === 'append';
    const index   = (fontEditMode === 'insert') ? fontEditIndex : 0;
    run(() => FontEdit(
      fontEditCz, fontEditInfo, fontEditTtf, fontEditOutCz, fontEditOutInfo, fontEditCharsetFile,
      redraw, append, index,
      fontEditArabicMetrics, fontEditMetricSetYEnabled, fontEditMetricSetY,
      fontEditMetricYOffset, fontEditMetricXOffset, fontEditMetricWOffset,
      fontEditArabicConnectorBleed
    ));
  }

  function getVietYOffsets() {
    const values = [];
    if (vietYMinus2) values.push(-2);
    if (vietYMinus1) values.push(-1);
    if (vietY0) values.push(0);
    if (vietY1) values.push(1);
    if (vietY2) values.push(2);
    if (vietY3) values.push(3);
    return values;
  }
  function startVietnameseFontPatch() {
    run(() => VietnameseFontPatch(vietFontRoot, vietCharsetFile, vietTtfFile, vietOutputDir, vietSlot, vietFamily, getVietYOffsets(), vietRedrawLatin));
  }

  // --- Luca Menu DLL helpers ---
  function lucaProfiles() {
    return lucaInventory?.profiles || [];
  }

  function currentLucaProfile() {
    return lucaProfiles().find(p => p.id === lucaGame) || null;
  }

  function lucaProxyName() {
    return lucaProxyChoice;
  }

  function lucaGameProfiles() {
    return lucaProfiles().filter(p => !p.id.includes('/'));
  }

  function lucaFamilyProfiles() {
    const profile = currentLucaProfile();
    if (!profile) return [];
    const root = (profile.id || '').split('/')[0];
    return lucaProfiles().filter(p => p.id === root || p.id.startsWith(root + '/'));
  }

  function lucaSlotEntryCount(profile, slot = lucaSlot) {
    return (profile?.entries || []).filter(e => e.slot === slot).length;
  }

  function lucaSlotSourceProfile(slot = lucaSlot) {
    const current = currentLucaProfile();
    return lucaFamilyProfiles()
      .filter(p => lucaSlotEntryCount(p, slot) > 0)
      .sort((a, b) => {
        const expected = slot === 'jp' ? '/slots-jp' : slot === 'cn' ? '/slots-cn' : '';
        const aPreferred = expected && a.id.toLowerCase().endsWith(expected) ? 1 : 0;
        const bPreferred = expected && b.id.toLowerCase().endsWith(expected) ? 1 : 0;
        if (aPreferred !== bPreferred) return bPreferred - aPreferred;
        const countDiff = lucaSlotEntryCount(b, slot) - lucaSlotEntryCount(a, slot);
        if (countDiff) return countDiff;
        if (a.id === current?.id) return -1;
        if (b.id === current?.id) return 1;
        return a.id.length - b.id.length || a.id.localeCompare(b.id);
      })[0] || null;
  }

  function lucaAvailableSlotCount(slot) {
    const source = lucaSlotSourceProfile(slot);
    return lucaSlotEntryCount(source, slot);
  }

  function lucaSafeFrenchCount(slot = lucaSlot) {
    const source = lucaSlotSourceProfile(slot);
    return (source?.entries || []).filter(e => e.slot === slot && e.safeAuto && e.suggestedFr).length;
  }

  function slotLabel(slot) {
    if (slot === 'jp') return t('Japonais', 'Japanese');
    if (slot === 'cn') return t('Chinois', 'Chinese');
    return t('Anglais', 'English');
  }

  function byteLen(value) {
    return new TextEncoder().encode(value || '').length;
  }

  function lucaEncodedLen(entry, value = entry?.target || '') {
    return entry?.encoding === 'utf-16-le' ? String(value).length * 2 : byteLen(value);
  }

  function entryTooLong(entry) {
    return entry.budget >= 0 && lucaEncodedLen(entry) > entry.budget;
  }

  async function loadLucaInventory() {
    lucaInventory = await ScanLucaMenuKit();
    const profiles = lucaProfiles();
    if (!profiles.length) {
      addLine('[ERROR] No Luca DLL patch profiles found.');
      return;
    }
    if (!lucaGame || !profiles.find(p => p.id === lucaGame)) {
      const preferred = profiles.find(p => p.id === 'Kanon') || profiles.find(p => p.id === 'AIR') || profiles[0];
      lucaGame = preferred.id;
    }
    syncLucaProfileDefaults();
    refreshLucaEntries();
    addLine('Luca kit: ' + lucaInventory.kitDir);
    addLine('Luca games: ' + lucaGameProfiles().map(p => p.id).join(', '));
  }

  function syncLucaProfileDefaults() {
    const profile = currentLucaProfile();
    if (!profile) return;
    lucaPatchName = profile.patchGameName || profile.name || lucaGame;
    lucaPatchVersion = profile.patchVersion || '0.1-gui';
    const identity = [profile.id, profile.name, profile.patchGameName, profile.gameExe, lucaGame]
      .filter(Boolean).join(' ').toUpperCase();
    lucaProxyChoice = identity.includes('LBEE') || identity.includes('LITTLE BUSTERS') || identity.includes('LITBUS_WIN32')
      ? 'winmm'
      : (profile.proxyDll === 'winmm' ? 'winmm' : 'version');
  }

  function selectLucaProxy(choice, checked) {
    if (checked) {
      lucaProxyChoice = choice;
      lucaBuildDll = true;
    } else if (lucaProxyChoice === choice) {
      lucaBuildDll = false;
    }
  }

  function refreshLucaEntries(mode = '') {
    const profile = lucaSlotSourceProfile();
    if (!profile) {
      lucaEntries = [];
      return;
    }
    lucaEntries = profile.entries
      .filter(e => e.slot === lucaSlot)
      .map(e => ({ ...e, target: e.target || '', include: false }));
    applyLucaFillMode(mode || lucaFillMode);
  }

  function setLucaGame(value) {
    lucaGame = value;
    lucaExe = '';
    lucaBuildDll = true;
    lucaCustomPatch = '';
    syncLucaProfileDefaults();
    refreshLucaEntries();
  }

  function setLucaSlot(value) {
    lucaSlot = value;
    refreshLucaEntries();
  }

  function lucaPresetTarget(entry, mode) {
    if (mode === 'fr' || mode === 'fr-safe') return entry.suggestedFr || '';
    if (mode === 'en' || mode === 'en-safe') return entry.suggestedEn || '';
    if (mode === 'ar') return entry.suggestedAr || '';
    if (mode === 'ru') return entry.suggestedRu || '';
    if (mode === 'jp') return entry.suggestedJp || '';
    if (mode === 'cn') return entry.suggestedCn || '';
    return entry.target || '';
  }

  function applyLucaFillMode(mode = lucaFillMode) {
    lucaFillMode = mode;
    const safeMode = mode === 'fr-safe' || mode === 'en-safe';
    lucaEntries = lucaEntries.map(entry => {
      const target = lucaPresetTarget(entry, mode);
      const fits = entry.budget < 0 || lucaEncodedLen(entry, target) <= entry.budget;
      const safe = entry.commonCount >= 4 && (mode !== 'fr-safe' || entry.safeAuto);
      const include = target !== '' && target !== entry.source && fits && (!safeMode || safe);
      return { ...entry, target, include };
    });
  }

  function setLucaFillMode(value) {
    if (value === 'ru') {
      lucaProxyChoice = 'winmm';
      lucaBuildDll = true;
    }
    applyLucaFillMode(value);
  }

  function clearLucaTargets() {
    lucaEntries = lucaEntries.map(e => ({ ...e, target: '', include: false }));
  }

  function lucaVisibleEntries() {
    const q = lucaSearch.trim().toLowerCase();
    return lucaEntries.filter(e => {
      if (!q) return true;
      return [e.source, e.target, e.context, e.note, e.category, e.textKind]
        .some(v => (v || '').toLowerCase().includes(q));
    });
  }

  function lucaSelectedEntries() {
    return lucaEntries.filter(e => e.include && e.target);
  }

  async function browseLucaExe() {
    const f = await SelectFile(t("Sélectionner l'EXE du jeu Luca", 'Select Luca game EXE'), '*.exe', t('Fichiers exécutables', 'Executable files'));
    if (f) lucaExe = f;
  }

  async function browseLucaOutput() {
    const d = await SelectDirectory(t('Sélectionner le dossier de sortie DLL Luca', 'Select Luca DLL output folder'));
    if (d) lucaOutputDir = d;
  }

  async function browseLucaCustomPatch() {
    const f = await SelectFile(t('Sélectionner un fichier PATCHES personnalisé', 'Select custom PATCHES file'), '*.py', t('Fichiers Python PATCHES', 'Python PATCHES files'));
    if (f) {
      lucaCustomPatch = f;
      lucaProxyChoice = 'winmm';
      lucaBuildDll = true;
    }
  }

  function startLucaGenerate() {
    const entries = lucaSelectedEntries().map(e => ({
      rawOffset: e.rawOffset,
      source: e.source,
      target: e.target,
      context: e.context,
      note: e.note,
      encoding: e.encoding,
      include: e.include,
      budget: e.budget
    }));
    if (!entries.length && !lucaCustomPatch) {
      addLine(t('[ERROR] Aucune chaîne remplie et sélectionnée pour le slot ', '[ERROR] No filled and selected string for the ') + slotLabel(lucaSlot) + t('.', ' slot.'));
      return;
    }
    run(() => LucaMenuGenerate({
      profileId: lucaGame,
      gameExe: lucaExe,
      outputDir: lucaOutputDir,
      patchGameName: lucaPatchName,
      patchVersion: lucaPatchVersion,
      slot: lucaSlot,
      buildDll: lucaBuildDll,
      proxyDll: lucaProxyChoice,
      preset: lucaFillMode,
      customPatch: lucaCustomPatch,
      entries
    }));
  }

  function startImageExport() {
    if (imgExpBatch) run(() => ImageBatchExport(imgExpInput, imgExpOutput));
    else run(() => ImageExport(imgExpInput, imgExpOutput));
  }
  function startImageImport() {
    if (imgImpBatch) run(() => ImageBatchImport(imgImpSource, imgImpInput, imgImpOutput, imgImpFill));
    else run(() => ImageImport(imgImpSource, imgImpInput, imgImpOutput, imgImpFill));
  }

  function selectOp(op) {
    if (op.disabled || op.section) return;
    selectedOp = op.id;
    if (selectedOp === 'luca_menu_dll' && !lucaInventory) loadLucaInventory();
  }

  // Reset fields when switching batch mode
  function toggleExpBatch() { imgExpInput = ''; imgExpOutput = ''; }
  function toggleImpBatch() { imgImpSource = ''; imgImpInput = ''; imgImpOutput = ''; }

  // --- Dialogue helpers ---
  async function browseDlgExtInput() {
    if (dlgExtBatch) { const d = await SelectDirectory('Select scripts folder'); if (d) dlgExtInput = d; }
    else { const f = await SelectScriptTxtFile(); if (f) dlgExtInput = f; }
    if (dlgExtInput) await detectDlgFormat();
  }
  async function browseDlgExtOutput() {
    if (dlgExtBatch) { const d = await SelectDirectory('Select output folder'); if (d) dlgExtOutput = d; }
    else {
      const defName = dlgExtInput ? dlgExtInput.replace(/\.txt$/i, '.ext.txt').split(/[\\/]/).pop() : 'dialogues.ext.txt';
      const f = await SelectSaveTsvFile(defName);
      if (f) dlgExtOutput = f;
    }
  }
  async function detectDlgFormat() {
    if (!dlgExtInput) return;
    const target = dlgExtBatch ? '' : dlgExtInput;
    if (!target) return;
    const info = await DialogueDetectFormat(target);
    dlgExtDetectedFmt = info.format || 'Unknown';
    dlgExtMaxCols = info.maxCols || 0;
  }
  function toggleDlgExtBatch() { dlgExtInput = ''; dlgExtOutput = ''; dlgExtDetectedFmt = ''; dlgExtMaxCols = 0; }

  async function browseDlgImpScript() {
    if (dlgImpBatch) { const d = await SelectDirectory('Select original scripts folder'); if (d) dlgImpScript = d; }
    else { const f = await SelectScriptTxtFile(); if (f) dlgImpScript = f; }
  }
  async function browseDlgImpTsv() {
    if (dlgImpBatch) { const d = await SelectDirectory('Select TSV folder'); if (d) dlgImpTsv = d; }
    else { const f = await SelectTsvFile(); if (f) dlgImpTsv = f; }
  }
  async function browseDlgImpOutput() {
    if (dlgImpBatch) { const d = await SelectDirectory('Select output folder'); if (d) dlgImpOutput = d; }
    else {
      const defName = dlgImpScript ? dlgImpScript.split(/[\\/]/).pop() : 'patched.txt';
      const f = await SelectSaveScriptFile(defName);
      if (f) dlgImpOutput = f;
    }
  }
  function toggleDlgImpBatch() { dlgImpScript = ''; dlgImpTsv = ''; dlgImpOutput = ''; }

  function getDlgExtCols() {
    const cols = [];
    if (dlgExtLang1) cols.push(1);
    if (dlgExtLang2) cols.push(2);
    if (dlgExtLang3) cols.push(3);
    if (dlgExtLang4) cols.push(4);
    return cols;
  }

  function startDlgExtract() {
    const cols = getDlgExtCols();
    if (dlgExtBatch) run(() => DialogueExtractBatch(dlgExtInput, dlgExtOutput, cols));
    else run(() => DialogueExtractFile(dlgExtInput, dlgExtOutput, cols));
  }
  function startDlgImport() {
    if (dlgImpBatch) run(() => DialogueImportBatch(dlgImpScript, dlgImpTsv, dlgImpTargetCol, dlgImpOutput));
    else run(() => DialogueImportFile(dlgImpScript, dlgImpTsv, dlgImpTargetCol, dlgImpOutput));
  }
</script>

<div id="app">
  <div class="titlebar">
    <span>LuckSystem 2.3.2 - Yoremi fork v3.30</span>
    <div class="titlebar-tools">
      <label class="ui-language">
        <span>{t('Interface', 'Interface', uiLanguage)}</span>
        <select value={uiLanguage} on:change={(e) => setUiLanguage(e.target.value)}>
          <option value="fr">Français</option>
          <option value="en">English</option>
        </select>
      </label>
      <span class="titlebar-path" on:click={locateLuckSystem} title={t('Cliquer pour modifier', 'Click to change', uiLanguage)}>
        {#if lsPath}📁 {lsPath}{:else}⚠ {t('lucksystem.exe introuvable — cliquer pour le localiser', 'lucksystem.exe not found — click to locate', uiLanguage)}{/if}
      </span>
    </div>
  </div>

  <div class="content">
    <!-- LEFT SIDEBAR -->
    <div class="sidebar">
      <div class="sidebar-title">{t('Choisir une option :', 'Select option:', uiLanguage)}</div>
      <div class="sidebar-list">
        {#each operations as op}
          {#if lucaMenuDllAvailable || (op.id !== '_s3c' && op.id !== 'luca_menu_dll')}
            {#if op.section}
              <div class="sidebar-section">{uiLanguage === 'fr' && op.labelFr ? op.labelFr : op.label}</div>
            {:else}
              <div class="sidebar-item" class:active={selectedOp === op.id} class:disabled={op.disabled} on:click={() => selectOp(op)}>
                {uiLanguage === 'fr' && op.labelFr ? op.labelFr : op.label}
              </div>
            {/if}
          {/if}
        {/each}
      </div>
    </div>

    <!-- RIGHT FORM PANEL -->
    <div class="form-panel">

      <!-- SCRIPT DECOMPILE -->
      {#if selectedOp === 'decompile'}
        <div class="form-title">Script Decompile</div>
        <div class="form-group"><label>SCRIPT.PAK file:</label><div class="form-row"><input type="text" bind:value={pakFile} readonly /><button class="btn" on:click={browsePak}>Select</button></div></div>
        {#if gamePresets.length > 0}
        <div class="form-group"><label>Game preset:</label><div class="form-row"><select value={selectedPreset} on:change={(e) => applyPreset(e.target.value)}><option value="">— Manual —</option>{#each gamePresets as p}<option value={p.name}>{p.name}{p.pluginFile ? ' (plugin)' : ''}</option>{/each}</select></div><div class="form-hint">Auto-fills Opcode, Plugin and Game from data/ folder</div></div>
        {/if}
        <div class="form-group"><label>Opcode file (.txt):</label><div class="form-row"><input type="text" bind:value={opcodeFile} placeholder="e.g. data/AIR.txt" readonly /><button class="btn" on:click={browseOpcode}>Select</button></div></div>
        <div class="form-group"><label>Plugin file (.py):</label><div class="form-row"><input type="text" bind:value={pluginFile} placeholder="e.g. data/AIR.py" readonly /><button class="btn" on:click={browsePlugin}>Select</button></div></div>
        <div class="form-group"><label>Charset:</label><div class="form-row"><select bind:value={charsetVal}><option value="UTF-8">UTF-8</option><option value="ShiftJIS">Shift-JIS</option><option value="GBK">GBK</option></select></div></div>
        <div class="form-group"><label>Output folder:</label><div class="form-row"><input type="text" bind:value={outputDir} readonly /><button class="btn" on:click={browseOutputDir}>Select</button></div><div class="form-hint">A SCRIPT.PAK subfolder will be created automatically inside</div></div>
        <div class="form-actions">{#if running}<span class="running-indicator"></span> Running...{:else}<button class="btn btn-primary" on:click={startDecompile} disabled={!pakFile || !outputDir}>Start Decompile</button>{/if}</div>

      <!-- SCRIPT COMPILE -->
      {:else if selectedOp === 'compile'}
        <div class="form-title">Script Compile</div>
        <div class="form-group"><label>Original SCRIPT.PAK:</label><div class="form-row"><input type="text" bind:value={pakFile} readonly /><button class="btn" on:click={browsePak}>Select</button></div><div class="form-hint">The original unmodified SCRIPT.PAK</div></div>
        {#if gamePresets.length > 0}
        <div class="form-group"><label>Game preset:</label><div class="form-row"><select value={selectedPreset} on:change={(e) => applyPreset(e.target.value)}><option value="">— Manual —</option>{#each gamePresets as p}<option value={p.name}>{p.name}{p.pluginFile ? ' (plugin)' : ''}</option>{/each}</select></div><div class="form-hint">Auto-fills Opcode, Plugin and Game from data/ folder</div></div>
        {/if}
        <div class="form-group"><label>Opcode file (.txt):</label><div class="form-row"><input type="text" bind:value={opcodeFile} readonly /><button class="btn" on:click={browseOpcode}>Select</button></div></div>
        <div class="form-group"><label>Plugin file (.py):</label><div class="form-row"><input type="text" bind:value={pluginFile} readonly /><button class="btn" on:click={browsePlugin}>Select</button></div></div>
        <div class="form-group"><label>Charset:</label><div class="form-row"><select bind:value={charsetVal}><option value="UTF-8">UTF-8</option><option value="ShiftJIS">Shift-JIS</option><option value="GBK">GBK</option></select></div></div>
        <div class="form-group"><label>Translated scripts folder:</label><div class="form-row"><input type="text" bind:value={importDir} readonly /><button class="btn" on:click={browseImportDir}>Select</button></div><div class="form-hint">Select the parent folder containing SCRIPT.PAK (e.g. TRAD), not SCRIPT.PAK itself</div></div>
        <div class="form-group"><label>Output PAK file:</label><div class="form-row"><input type="text" bind:value={outputPak} readonly /><button class="btn" on:click={browseOutputPak}>Select</button></div></div>
        <div class="form-actions">{#if running}<span class="running-indicator"></span> Running...{:else}<button class="btn btn-primary" on:click={startCompile} disabled={!pakFile || !importDir || !outputPak}>Start Compile</button>{/if}</div>

      <!-- SIGLUS -> LUCA -->
      {:else if selectedOp === 'siglus_luca'}
        <div class="form-title">Siglus -> Luca Script Bridge</div>
        <div class="form-hint" style="margin-bottom:10px">
          {t('Importe des lignes traduites depuis des exports Siglus dans des scripts Luca décompilés. Les lignes Luca-only et les découpages à vérifier sont exportés en TSV dans le dossier de sortie.', 'Imports translated lines from Siglus exports into decompiled Luca scripts. Luca-only lines and segments requiring review are exported as TSV files in the output folder.', uiLanguage)}
        </div>
        <div class="form-group"><label>Luca scripts folder:</label><div class="form-row"><input type="text" bind:value={siglusLucaLucaDir} readonly placeholder="SCRIPT.PAK decompiled folder" /><button class="btn" on:click={browseSiglusLucaLucaDir}>Select</button></div><div class="form-hint">{t('Dossier contenant les scripts Luca .txt à patcher.', 'Folder containing the Luca .txt scripts to patch.', uiLanguage)}</div></div>
        <div class="form-group"><label>Siglus Full folder:</label><div class="form-row"><input type="text" bind:value={siglusLucaSiglusDir} readonly placeholder="TRAD-silgus\Full" /><button class="btn" on:click={browseSiglusLucaSiglusDir}>Select</button></div><div class="form-hint">{t('Dossier contenant les exports Siglus .ss.txt avec source et traduction.', 'Folder containing Siglus .ss.txt exports with source and translated text.', uiLanguage)}</div></div>
        <div class="form-group"><label>Output folder:</label><div class="form-row"><input type="text" bind:value={siglusLucaOutput} readonly placeholder="Luca_from_Siglus_FR" /><button class="btn" on:click={browseSiglusLucaOutput}>Select</button></div><div class="form-hint">{t('Les scripts patchés,', 'The patched scripts,', uiLanguage)} <code>hd_candidates.tsv</code> {t('et', 'and', uiLanguage)} <code>review.tsv</code> {t('seront écrits ici.', 'will be written here.', uiLanguage)}</div></div>
        <div class="form-group">
          <label>Target column:</label>
          <div class="form-row">
            <select bind:value={siglusLucaTargetCol}>
              <option value={1}>Lang 1</option>
              <option value={2}>Lang 2 (English slot)</option>
              <option value={3}>Lang 3</option>
              <option value={4}>Lang 4</option>
            </select>
          </div>
          <div class="form-hint">{t('Garder Lang 2 pour les scripts Luca dont le deuxième slot texte est la langue à remplacer.', 'Keep Lang 2 for Luca scripts whose second text slot is the language to replace.', uiLanguage)}</div>
        </div>
        <div class="form-actions">
          {#if running}<span class="running-indicator"></span> Running...
          {:else}<button class="btn btn-primary" on:click={startSiglusLucaBridge}
            disabled={!siglusLucaLucaDir || !siglusLucaSiglusDir || !siglusLucaOutput}>
            Import Siglus FR
          </button>{/if}
        </div>

      <!-- PAK CG EXTRACT -->
      {:else if selectedOp === 'pak_cg_extract'}
        <div class="form-title">PAK (CG) — Extract</div>
        <div class="form-group"><label>PAK file (CG) :</label><div class="form-row"><input type="text" bind:value={pakExtSource} readonly /><button class="btn" on:click={browsePakExtSource}>Select</button></div></div>
        <div class="form-group"><label>Output folder:</label><div class="form-row"><input type="text" bind:value={pakExtOutput} readonly /><button class="btn" on:click={browsePakExtOutput}>Select</button></div><div class="form-hint">{t('Le fichier liste', 'The list file', uiLanguage)} <code>&lt;NAME&gt;_list.txt</code> {t('sera généré automatiquement dans ce dossier', 'will be generated automatically in this folder', uiLanguage)}</div></div>
        <div class="form-actions">{#if running}<span class="running-indicator"></span> Running...{:else}<button class="btn btn-primary" on:click={startPakExtract} disabled={!pakExtSource || !pakExtOutput}>Start Extract</button>{/if}</div>

      <!-- BGMOVIE EXTRACT -->
      {:else if selectedOp === 'bgmovie_extract'}
        <div class="form-title">BGMOVIE.PAK — Video Extract</div>
        <div class="form-group"><label>BGMOVIE.PAK file:</label><div class="form-row"><input type="text" bind:value={bgMoviePak} readonly /><button class="btn" on:click={browseBgMoviePak}>Select</button></div></div>
        <div class="form-group"><label>Output folder:</label><div class="form-row"><input type="text" bind:value={bgMovieOutput} readonly /><button class="btn" on:click={browseBgMovieOutput}>Select</button></div><div class="form-hint">Creates a folder named after the PAK with raw MVT files and a <code>webm</code> subfolder.</div></div>
        <div class="form-actions">{#if running}<span class="running-indicator"></span> Running...{:else}<button class="btn btn-primary" on:click={startBgMovieExtract} disabled={!bgMoviePak || !bgMovieOutput}>Extract Videos</button>{/if}</div>

      <!-- MUSIC EXTRACT -->
      {:else if selectedOp === 'music_extract'}
        <div class="form-title">MUSIC.PAK — Audio Extract</div>
        <div class="form-group"><label>MUSIC.PAK file:</label><div class="form-row"><input type="text" bind:value={musicPak} readonly /><button class="btn" on:click={browseMusicPak}>Select</button></div></div>
        <div class="form-group"><label>Output folder:</label><div class="form-row"><input type="text" bind:value={musicOutput} readonly /><button class="btn" on:click={browseMusicOutput}>Select</button></div><div class="form-hint">Native Ogg files go to <code>ogg</code>; the generated list can be reused by PAK Replace.</div></div>
        <div class="form-group">
          <label>Conversion:</label>
          <div class="form-row checkbox-row">
            <label class="checkbox-label"><input type="checkbox" bind:checked={musicToMp3} /> Also create MP3 copies</label>
          </div>
        </div>
        <div class="form-actions">{#if running}<span class="running-indicator"></span> Running...{:else}<button class="btn btn-primary" on:click={startMusicExtract} disabled={!musicPak || !musicOutput}>Extract Music</button>{/if}</div>

      <!-- VOICE EXTRACT -->
      {:else if selectedOp === 'voice_extract'}
        <div class="form-title">VOICE / SYSVOICE.PAK — Audio Extract</div>
        <div class="form-group"><label>VOICE PAK file:</label><div class="form-row"><input type="text" bind:value={voicePak} readonly /><button class="btn" on:click={browseVoicePak}>Select</button></div></div>
        <div class="form-group"><label>Output folder:</label><div class="form-row"><input type="text" bind:value={voiceOutput} readonly /><button class="btn" on:click={browseVoiceOutput}>Select</button></div><div class="form-hint">Works with VOICE, VOICE0/1 and SYSVOICE PAKs. Native Ogg files go to <code>ogg</code>.</div></div>
        <div class="form-group">
          <label>Conversion:</label>
          <div class="form-row checkbox-row">
            <label class="checkbox-label"><input type="checkbox" bind:checked={voiceToMp3} /> Also create MP3 copies</label>
          </div>
        </div>
        <div class="form-actions">{#if running}<span class="running-indicator"></span> Running...{:else}<button class="btn btn-primary" on:click={startVoiceExtract} disabled={!voicePak || !voiceOutput}>Extract Voices</button>{/if}</div>

      <!-- AUDIO CONVERT -->
      {:else if selectedOp === 'audio_convert'}
        <div class="form-title">Audio Convert — Ogg / MP3</div>
        <div class="form-group">
          <label>Mode:</label>
          <div class="form-row">
            <select bind:value={audioConvDirection}>
              <option value="mp3">Native Ogg → MP3</option>
              <option value="native">MP3 → Native Ogg</option>
            </select>
          </div>
        </div>
        <div class="form-group"><label>{audioConvDirection === 'mp3' ? 'Native Ogg folder:' : 'MP3 folder:'}</label><div class="form-row"><input type="text" bind:value={audioConvInput} readonly /><button class="btn" on:click={browseAudioConvInput}>Select</button></div></div>
        <div class="form-group"><label>Output folder:</label><div class="form-row"><input type="text" bind:value={audioConvOutput} readonly /><button class="btn" on:click={browseAudioConvOutput}>Select</button></div></div>
        <div class="form-actions">{#if running}<span class="running-indicator"></span> Running...{:else}<button class="btn btn-primary" on:click={startAudioConvert} disabled={!audioConvInput || !audioConvOutput}>Convert Audio</button>{/if}</div>

      <!-- PAK CG REPLACE -->
      {:else if selectedOp === 'pak_cg_replace'}
        <div class="form-title">PAK (CG) — Replace</div>
        <div class="form-group"><label>Original PAK file:</label><div class="form-row"><input type="text" bind:value={pakRepSource} readonly /><button class="btn" on:click={browsePakRepSource}>Select</button></div></div>
        <div class="form-group">
          <label>{t("Mode d'entrée :", 'Input mode:', uiLanguage)}</label>
          <div class="form-row checkbox-row" style="margin-bottom:6px">
            <label class="checkbox-label"><input type="radio" bind:group={pakRepUseList} value={true} /> {t('Fichier liste', 'List file', uiLanguage)} (<code>*_list.txt</code>)</label>
            <label class="checkbox-label"><input type="radio" bind:group={pakRepUseList} value={false} /> {t('Dossier de fichiers', 'Files folder', uiLanguage)}</label>
          </div>
          {#if pakRepUseList}
            <div class="form-row"><input type="text" bind:value={pakRepListFile} placeholder="SYSCG_list.txt" readonly /><button class="btn" on:click={browsePakRepListFile}>Select</button></div>
            <div class="form-hint">{t("Fichier liste généré lors de l'extraction (ex : SYSCG_list.txt)", 'List file generated during extraction (e.g. SYSCG_list.txt)', uiLanguage)}</div>
          {:else}
            <div class="form-row"><input type="text" bind:value={pakRepInput} readonly /><button class="btn" on:click={browsePakRepInput}>Select</button></div>
            <div class="form-hint">{t('Dossier contenant les fichiers modifiés à réinjecter', 'Folder containing modified files to reinsert', uiLanguage)}</div>
          {/if}
        </div>
        <div class="form-group"><label>Output PAK:</label><div class="form-row"><input type="text" bind:value={pakRepOutput} readonly /><button class="btn" on:click={browsePakRepOutput}>Select</button></div></div>
        <div class="form-actions">
          {#if running}
            <span class="running-indicator"></span> Running...
          {:else}
            <button class="btn btn-primary" on:click={startPakReplace}
              disabled={!pakRepSource || !pakRepOutput || (pakRepUseList ? !pakRepListFile : !pakRepInput)}>
              Start Replace
            </button>
          {/if}
        </div>

      <!-- PAK FONT EXTRACT -->
      {:else if selectedOp === 'pak_font_extract'}
        <div class="form-title">PAK (Font) — Extract</div>
        <div class="form-group"><label>PAK file (Font) :</label><div class="form-row"><input type="text" bind:value={pakFontExtSource} readonly /><button class="btn" on:click={browsePakFontExtSource}>Select</button></div></div>
        <div class="form-group"><label>Charset :</label><div class="form-row"><select bind:value={pakFontExtCharset}><option value="UTF-8">UTF-8</option><option value="ShiftJIS">Shift-JIS</option><option value="GBK">GBK</option></select></div></div>
        <div class="form-group"><label>Output folder :</label><div class="form-row"><input type="text" bind:value={pakFontExtOutput} readonly /><button class="btn" on:click={browsePakFontExtOutput}>Select</button></div><div class="form-hint">{t('Tous les fichiers du PAK seront extraits ici', 'All files in the PAK will be extracted here', uiLanguage)}</div></div>
        <div class="form-actions">{#if running}<span class="running-indicator"></span> Running...{:else}<button class="btn btn-primary" on:click={startPakFontExtract} disabled={!pakFontExtSource || !pakFontExtOutput}>Start Extract</button>{/if}</div>

      <!-- PAK FONT REPLACE -->
      {:else if selectedOp === 'pak_font_replace'}
        <div class="form-title">PAK (Font) — Replace</div>
        <div class="form-group"><label>Original PAK file (Font) :</label><div class="form-row"><input type="text" bind:value={pakFontRepSource} readonly /><button class="btn" on:click={browsePakFontRepSource}>Select</button></div></div>
        <div class="form-group"><label>Charset :</label><div class="form-row"><select bind:value={pakFontRepCharset}><option value="UTF-8">UTF-8</option><option value="ShiftJIS">Shift-JIS</option><option value="GBK">GBK</option></select></div></div>
        <div class="form-group">
          <label>{t("Mode d'entrée :", 'Input mode:', uiLanguage)}</label>
          <div class="form-row checkbox-row" style="margin-bottom:6px">
            <label class="checkbox-label"><input type="radio" bind:group={pakFontRepMode} value="list" /> {t('Fichier liste', 'List file', uiLanguage)} (<code>*_list.txt</code>)</label>
            <label class="checkbox-label"><input type="radio" bind:group={pakFontRepMode} value="dir" /> {t('Dossier de fichiers', 'Files folder', uiLanguage)}</label>
            <label class="checkbox-label"><input type="radio" bind:group={pakFontRepMode} value="single" /> {t('Fichier unique par nom', 'Single file by name', uiLanguage)}</label>
            <label class="checkbox-label"><input type="radio" bind:group={pakFontRepMode} value="alias" /> {t('Alias de taille compatible', 'Compatible-size alias', uiLanguage)}</label>
          </div>
          {#if pakFontRepMode === 'list'}
            <div class="form-row"><input type="text" bind:value={pakFontRepListFile} placeholder="FONT__INFO_list.txt" readonly /><button class="btn" on:click={browsePakFontRepListFile}>Select</button></div>
            <div class="form-hint">{t("Fichier liste généré lors de l'extraction (ex : FONT__INFO_list.txt)", 'List file generated during extraction (e.g. FONT__INFO_list.txt)', uiLanguage)}</div>
          {:else if pakFontRepMode === 'dir'}
            <div class="form-row"><input type="text" bind:value={pakFontRepInput} readonly /><button class="btn" on:click={browsePakFontRepInput}>Select</button></div>
            <div class="form-hint">{t('Remplace uniquement les fichiers du dossier dont le nom existe dans le PAK.', 'Only replaces files whose names exist in the PAK.', uiLanguage)}</div>
          {:else if pakFontRepMode === 'single'}
            <div class="form-row"><input type="text" bind:value={pakFontRepSingleFile} readonly placeholder="ex : C:\dossier\info30" /><button class="btn" on:click={browsePakFontRepSingleFile}>Select</button></div>
            <div class="form-row" style="margin-top:6px"><input type="text" bind:value={pakFontRepSingleName} placeholder="Nom interne exact : info30 ou 明朝30" /></div>
            <div class="form-hint">{t('Recommandé pour Kanon : effectuez deux remplacements séparés,', 'Recommended for Kanon: make two separate replacements,', uiLanguage)} <code>info30</code> {t('dans', 'in', uiLanguage)} <code>FONT__INFO.PAK</code>, {t('puis', 'then', uiLanguage)} <code>明朝30</code> {t('dans', 'in', uiLanguage)} <code>FONT_MINCHO.PAK</code>.</div>
          {:else}
            <div class="form-row">
              <span style="min-width:110px;font-size:12px">{t('Copier depuis :', 'Copy from:', uiLanguage)}</span>
              <input type="text" bind:value={pakFontRepAliasFrom} placeholder="info30 ou 明朝30" />
            </div>
            <div class="form-row" style="margin-top:6px">
              <span style="min-width:110px;font-size:12px">{t('Vers :', 'To:', uiLanguage)}</span>
              <input type="text" bind:value={pakFontRepAliasTo} placeholder="info32 ou 明朝32" />
            </div>
            <div class="form-hint">{t("Adapte les données de la taille source à la structure de la taille cible. Test Kanon arabe : info30 → info32, puis 明朝30 → 明朝32 dans l'autre PAK. Le CZ2 conserve la largeur, les cellules et la longueur d'entrée attendues pour la taille 32. Cela affecte toute l'entrée 32, pas uniquement SELECT.", 'Adapts the source-size data to the target-size structure. Arabic Kanon test: info30 → info32, then 明朝30 → 明朝32 in the other PAK. CZ2 preserves the width, cells, and expected entry length for size 32. This affects the entire size-32 entry, not only SELECT.', uiLanguage)}</div>
          {/if}
        </div>
        <div class="form-group"><label>Output PAK :</label><div class="form-row"><input type="text" bind:value={pakFontRepOutput} readonly /><button class="btn" on:click={browsePakFontRepOutput}>Select</button></div></div>
        <div class="form-actions">
          {#if running}
            <span class="running-indicator"></span> Running...
          {:else}
            <button class="btn btn-primary" on:click={startPakFontReplace}
              disabled={!pakFontRepSource || !pakFontRepOutput || (pakFontRepMode === 'list' ? !pakFontRepListFile : pakFontRepMode === 'dir' ? !pakFontRepInput : pakFontRepMode === 'single' ? (!pakFontRepSingleFile || !pakFontRepSingleName) : (!pakFontRepAliasFrom || !pakFontRepAliasTo))}>
              Start Replace
            </button>
          {/if}
        </div>

      <!-- FONT EXTRACT -->
      {:else if selectedOp === 'font_extract'}
        <div class="form-title">Font Extract</div>
        <div class="form-group"><label>Font CZ file (e.g. 明朝32):</label><div class="form-row"><input type="text" bind:value={fontExtCz} readonly /><button class="btn" on:click={browseFontExtCz}>Select</button></div></div>
        <div class="form-group"><label>Info file (e.g. info32):</label><div class="form-row"><input type="text" bind:value={fontExtInfo} readonly /><button class="btn" on:click={browseFontExtInfo}>Select</button></div><div class="form-hint">Must match font size (info32 for 明朝32)</div></div>
        <div class="form-group"><label>Output PNG:</label><div class="form-row"><input type="text" bind:value={fontExtPng} readonly /><button class="btn" on:click={browseFontExtPng}>Select</button></div></div>
        <div class="form-group"><label>Output charset TXT (optional):</label><div class="form-row"><input type="text" bind:value={fontExtCharset} readonly /><button class="btn" on:click={browseFontExtCharset}>Select</button></div></div>
        <div class="form-actions">{#if running}<span class="running-indicator"></span> Running...{:else}<button class="btn btn-primary" on:click={startFontExtract} disabled={!fontExtCz || !fontExtInfo || !fontExtPng}>Start Extract</button>{/if}</div>

      <!-- FONT EDIT -->
      {:else if selectedOp === 'font_edit'}
        <div class="form-title">{t('Font Edit — Modification de glyphes', 'Font Edit — Glyph modification', uiLanguage)}</div>
        <div class="form-hint form-hint-warn">⚠ {t("Font Edit modifie les glyphes d'un fichier CZ avec un TTF. Pour simplement re-packer un PAK de police, utilisez", 'Font Edit modifies glyphs in a CZ file using a TTF. To simply repack a font PAK, use', uiLanguage)} <strong>PAK (Font) → Font Replace</strong>.</div>

        <div class="form-group"><label>Source CZ file:</label><div class="form-row"><input type="text" bind:value={fontEditCz} readonly /><button class="btn" on:click={browseFontEditCz}>Select</button></div></div>
        <div class="form-group"><label>Source info file:</label><div class="form-row"><input type="text" bind:value={fontEditInfo} readonly /><button class="btn" on:click={browseFontEditInfo}>Select</button></div></div>
        <div class="form-group"><label>TTF font file:</label><div class="form-row"><input type="text" bind:value={fontEditTtf} readonly /><button class="btn" on:click={browseFontEditTtf}>Select</button></div></div>

        <div class="form-group">
          <label>Mode :</label>
          <div class="form-row checkbox-row" style="margin-bottom:6px">
            <label class="checkbox-label"><input type="radio" bind:group={fontEditMode} value="redraw" /> Redraw all</label>
            <label class="checkbox-label"><input type="radio" bind:group={fontEditMode} value="append" /> Append to end</label>
            <label class="checkbox-label"><input type="radio" bind:group={fontEditMode} value="insert" /> Insert at index</label>
          </div>
          {#if fontEditMode === 'redraw'}
            <div class="form-hint">{t('Redessine TOUS les glyphes existants avec le TTF. Aucun charset requis.', 'Redraws ALL existing glyphs with the TTF. No charset is required.', uiLanguage)}</div>
          {:else if fontEditMode === 'append'}
            <div class="form-hint">{t('Ajoute les caractères du charset à la fin de la police.', 'Appends the charset characters to the end of the font.', uiLanguage)}</div>
          {:else if fontEditMode === 'insert'}
            <div class="form-row" style="margin-top:4px">
              <span style="min-width:90px;font-size:12px">Start index :</span>
              <input type="number" bind:value={fontEditIndex} min="0" style="width:80px;height:26px;padding:0 6px;border:1px solid #c0c0c0;border-radius:2px" />
            </div>
            <div class="form-hint">{t('Insère/remplace à partir de cette position (index commençant à 0).', 'Inserts/replaces from this position (zero-based index).', uiLanguage)}</div>
          {/if}
        </div>

        {#if fontEditMode !== 'redraw'}
          <div class="form-group"><label>Charset file <span class="required">*</span> :</label><div class="form-row"><input type="text" bind:value={fontEditCharsetFile} readonly /><button class="btn" on:click={browseFontEditCharset}>Select</button></div><div class="form-hint">{t('Fichier texte listant les caractères à ajouter/insérer (ex : accents_fr.txt)', 'Text file listing the characters to append/insert (e.g. accents_fr.txt)', uiLanguage)}</div></div>
        {/if}

        <div class="form-group">
          <label>Metrics adjustment :</label>
          <div class="form-row checkbox-row" style="margin-bottom:6px">
            <label class="checkbox-label"><input type="checkbox" checked={fontEditArabicMetrics} on:change={(e) => setFontEditArabicPreset(e.target.checked)} /> Arabic preset</label>
            <label class="checkbox-label"><input type="checkbox" bind:checked={fontEditMetricSetYEnabled} /> Set Y</label>
            {#if fontEditMetricSetYEnabled}
              <input type="number" bind:value={fontEditMetricSetY} style="width:70px;height:26px;padding:0 6px;border:1px solid #c0c0c0;border-radius:2px" />
            {/if}
          </div>
          <div class="form-row" style="gap:8px;flex-wrap:wrap">
            <span style="font-size:12px">Y offset</span>
            <input type="number" bind:value={fontEditMetricYOffset} style="width:70px;height:26px;padding:0 6px;border:1px solid #c0c0c0;border-radius:2px" />
            <span style="font-size:12px">X offset</span>
            <input type="number" bind:value={fontEditMetricXOffset} style="width:70px;height:26px;padding:0 6px;border:1px solid #c0c0c0;border-radius:2px" />
            <span style="font-size:12px">Advance offset</span>
            <input type="number" bind:value={fontEditMetricWOffset} style="width:70px;height:26px;padding:0 6px;border:1px solid #c0c0c0;border-radius:2px" />
            <span style="font-size:12px">Connector bleed</span>
            <input type="number" min="0" max="8" bind:value={fontEditArabicConnectorBleed} style="width:70px;height:26px;padding:0 6px;border:1px solid #c0c0c0;border-radius:2px" />
          </div>
          <div class="form-hint">Arabic preset shifts Arabic glyphs toward the Latin baseline. Connector bleed is experimental and stays manual.</div>
        </div>

        <div class="form-group"><label>Output CZ <span class="required">*</span> :</label><div class="form-row"><input type="text" bind:value={fontEditOutCz} placeholder="ex: C:\dossier\ゴシック26" /><button class="btn" on:click={browseFontEditOutCz}>📁</button></div><div class="form-hint">{t('Tapez le chemin complet sans extension — le bouton sélectionne le dossier', 'Enter the full path without an extension — the button selects the folder', uiLanguage)}</div></div>
        <div class="form-group"><label>Output info <span class="required">*</span> :</label><div class="form-row"><input type="text" bind:value={fontEditOutInfo} placeholder="ex: C:\dossier\info26" /><button class="btn" on:click={browseFontEditOutInfo}>📁</button></div><div class="form-hint">{t('Tapez le chemin complet sans extension — requis pour mettre à jour le compte de caractères', 'Enter the full path without an extension — required to update the character count', uiLanguage)}</div></div>

        <div class="form-actions">
          {#if running}
            <span class="running-indicator"></span> Running...
          {:else}
            <button class="btn btn-primary" on:click={startFontEdit}
              disabled={!fontEditCz || !fontEditInfo || !fontEditTtf || !fontEditOutCz || !fontEditOutInfo
                || (fontEditMode !== 'redraw' && !fontEditCharsetFile)}>
              Start Edit
            </button>
          {/if}
        </div>

      <!-- VIETNAMESE FONT PATCH -->
      {:else if selectedOp === 'viet_font_patch'}
        <div class="form-title">AIR / Planetarian SG — Vietnamese Font Patch</div>
        <div class="form-hint form-hint-warn">This dedicated workflow generates safe Vietnamese font PAKs directly from the original game font folder. Use the full 134-character Vietnamese charset file.</div>

        <div class="form-group"><label>Game files folder:</label><div class="form-row"><input type="text" bind:value={vietFontRoot} readonly placeholder="ex: C:\Games\AIR\files or C:\Games\Planetarian SG\files" /><button class="btn" on:click={browseVietFontRoot}>Select</button></div><div class="form-hint">Select the folder that contains <code>font_win32_1280</code>.</div></div>
        <div class="form-group"><label>Full Vietnamese charset file:</label><div class="form-row"><input type="text" bind:value={vietCharsetFile} readonly placeholder="examples\AIR_vietnamese_full_134.txt" /><button class="btn" on:click={browseVietCharset}>Select</button></div><div class="form-hint">Use the full 134-character charset. The tool keeps the 32 existing characters and injects only the missing 102.</div></div>
        <div class="form-group"><label>TTF / OTF font file:</label><div class="form-row"><input type="text" bind:value={vietTtfFile} readonly /><button class="btn" on:click={browseVietTtf}>Select</button></div></div>
        <div class="form-group"><label>Output folder:</label><div class="form-row"><input type="text" bind:value={vietOutputDir} readonly /><button class="btn" on:click={browseVietOutput}>Select</button></div><div class="form-hint">Each selected Y value creates a separate subfolder containing ready-to-test PAKs.</div></div>

        <div class="form-group">
          <label>Patch mode:</label>
          <div class="form-row checkbox-row">
            <label class="checkbox-label"><input type="checkbox" bind:checked={vietRedrawLatin} /> Experimental: redraw Latin alphabet from TTF</label>
          </div>
          <div class="form-hint">Default mode only injects missing Vietnamese glyphs. Experimental mode redraws A-Z/a-z and already-present Vietnamese glyphs at their original indexes, then injects the missing Vietnamese glyphs.</div>
        </div>

        <div class="form-group">
          <label>Target:</label>
          <div class="form-row">
            <select bind:value={vietSlot}>
              <option value="en">English slot (recommended)</option>
              <option value="zc">Chinese ZC slot</option>
              <option value="all">Both slots</option>
            </select>
            <select bind:value={vietFamily}>
              <option value="GOTHIC1">GOTHIC1 quick test</option>
              <option value="GOTHIC2">GOTHIC2</option>
              <option value="GOTHIC3">GOTHIC3</option>
              <option value="MINCHO">MINCHO</option>
              <option value="MODERN">MODERN</option>
              <option value="all">All English families</option>
            </select>
          </div>
          <div class="form-hint">For first tests, keep English slot + GOTHIC1. Generate all families only after a good Y value is found.</div>
        </div>

        <div class="form-group">
          <label>Y alignment values:</label>
          <div class="form-row checkbox-row">
            <label class="checkbox-label"><input type="checkbox" bind:checked={vietYMinus2} /> Y-2</label>
            <label class="checkbox-label"><input type="checkbox" bind:checked={vietYMinus1} /> Y-1</label>
            <label class="checkbox-label"><input type="checkbox" bind:checked={vietY0} /> Y+0</label>
            <label class="checkbox-label"><input type="checkbox" bind:checked={vietY1} /> Y+1</label>
            <label class="checkbox-label"><input type="checkbox" bind:checked={vietY2} /> Y+2</label>
            <label class="checkbox-label"><input type="checkbox" bind:checked={vietY3} /> Y+3</label>
          </div>
          <div class="form-hint">Y+2 is the validated AIR value. Select several values to create several test folders at once.</div>
        </div>

        <div class="form-actions">
          {#if running}
            <span class="running-indicator"></span> Running...
          {:else}
            <button class="btn btn-primary" on:click={startVietnameseFontPatch}
              disabled={!vietFontRoot || !vietCharsetFile || !vietTtfFile || !vietOutputDir || getVietYOffsets().length === 0}>
              Generate Vietnamese Font PAKs
            </button>
          {/if}
        </div>

      <!-- LUCA MENU DLL -->
      {:else if selectedOp === 'luca_menu_dll'}
        <div class="form-title">Luca Menu DLL</div>

        {#if !lucaInventory}
          <div class="form-actions" style="justify-content:flex-start">
            <button class="btn btn-primary" on:click={loadLucaInventory}>{t("Charger l'inventaire Luca", 'Load Luca inventory', uiLanguage)}</button>
          </div>
        {:else if !currentLucaProfile()}
          <div class="form-hint form-hint-warn">{t('Aucun profil Luca disponible dans le kit.', 'No Luca profile is available in the kit.', uiLanguage)}</div>
        {:else}
          <div class="form-group">
            <label>{t('Profil et slot :', 'Profile and slot:', uiLanguage)}</label>
            <div class="form-row">
              <select bind:value={lucaGame} on:change={() => setLucaGame(lucaGame)}>
                {#each lucaGameProfiles() as profile}
                  <option value={profile.id}>{profile.name} ({profile.id})</option>
                {/each}
              </select>
              <select value={lucaSlot} on:change={(e) => setLucaSlot(e.target.value)}>
                <option value="en">{t('Slot anglais', 'English slot', uiLanguage)}</option>
                <option value="jp">{t('Slot japonais', 'Japanese slot', uiLanguage)}</option>
                <option value="cn">{t('Slot chinois', 'Chinese slot', uiLanguage)}</option>
              </select>
              <button class="btn" on:click={loadLucaInventory}>Rescan</button>
            </div>
            <div class="form-hint">
              EN {lucaAvailableSlotCount('en')} · JP {lucaAvailableSlotCount('jp')} · CN {lucaAvailableSlotCount('cn')} · {t('FR sûr sur ce slot', 'safe FR strings in this slot', uiLanguage)} {lucaSafeFrenchCount()}
            </div>
            {#if lucaSlotSourceProfile() && lucaSlotSourceProfile().id !== currentLucaProfile().id}
              <div class="form-hint">{t('Inventaire du slot chargé depuis', 'Slot inventory loaded from', uiLanguage)} {lucaSlotSourceProfile().id}.</div>
            {:else if !lucaSlotSourceProfile()}
              <div class="form-hint form-hint-warn">{t('Aucune chaîne du slot', 'No strings from the', uiLanguage)} {slotLabel(lucaSlot)} {t("n'est encore inventoriée pour ce jeu.", 'slot have been inventoried for this game yet.', uiLanguage)}</div>
            {/if}
          </div>

          <div class="form-group">
            <label>{t('EXE du jeu', 'Game EXE', uiLanguage)} <span class="required">*</span> :</label>
            <div class="form-row"><input type="text" bind:value={lucaExe} placeholder={currentLucaProfile().gameExe || t('Sélectionnez le véritable EXE du jeu', 'Select the actual game EXE', uiLanguage)} /><button class="btn" on:click={browseLucaExe}>Select</button></div>
            <div class="form-hint">{t("Sélectionnez l'EXE présent dans le dossier du jeu afin de vérifier les offsets et la taille des chaînes.", 'Select the EXE located in the game folder so offsets and string sizes can be verified.', uiLanguage)}</div>
          </div>

          <div class="form-group">
            <label>{t('Dossier de sortie', 'Output folder', uiLanguage)} <span class="required">*</span> :</label>
            <div class="form-row"><input type="text" bind:value={lucaOutputDir} readonly /><button class="btn" on:click={browseLucaOutput}>Select</button></div>
            <div class="form-hint">{t('Le dossier recevra', 'The folder will receive', uiLanguage)} {lucaCustomPatch ? 'mixed_patches.py, custom_patches.py' : lucaFillMode === 'ru' ? 'mixed_patches.py, russian_preset.py' : 'patches.py'}, patches.h, patches.csv, version.c, {lucaProxyName()}.def {t('et', 'and', uiLanguage)} {lucaProxyName()}.dll {t('si la compilation réussit.', 'if compilation succeeds.', uiLanguage)}</div>
          </div>

          {#if (lucaGame || '').split('/')[0].toUpperCase() === 'LBEE'}
            <div class="form-group">
              <label>{t('Fichier PATCHES personnalisé (facultatif) :', 'Custom PATCHES file (optional):', uiLanguage)}</label>
              <div class="form-row">
                <input type="text" bind:value={lucaCustomPatch} readonly placeholder={t('russian_preset.py ou autre fichier contenant PATCHES', 'russian_preset.py or another file containing PATCHES', uiLanguage)} />
                <button class="btn" on:click={browseLucaCustomPatch}>Select</button>
                {#if lucaCustomPatch}<button class="btn" on:click={() => lucaCustomPatch = ''}>{t('Vider', 'Clear', uiLanguage)}</button>{/if}
              </div>
              <div class="form-hint">{t('Si renseigné, ce fichier remplace le preset russe interne. Le dossier de sortie reste uniquement la destination des fichiers générés.', 'When selected, this file replaces the built-in Russian preset. The output folder remains only the destination for generated files.', uiLanguage)}</div>
            </div>
          {/if}

          <div class="form-group">
            <label>{t('Identité du patch :', 'Patch identity:', uiLanguage)}</label>
            <div class="form-row">
              <input type="text" bind:value={lucaPatchName} placeholder={t('Nom affiché dans luckproxy.log', 'Name shown in luckproxy.log', uiLanguage)} />
              <input type="text" bind:value={lucaPatchVersion} placeholder="Version" style="max-width:140px" />
            </div>
          </div>

          <div class="form-group luca-proxy-group">
            <label>{t('DLL proxy à compiler :', 'Proxy DLL to compile:', uiLanguage)}</label>
            <div class="form-row checkbox-row luca-proxy-row">
              <label class:luca-proxy-selected={lucaBuildDll && lucaProxyChoice === 'version'} class="checkbox-label luca-proxy-option">
                <input type="checkbox" checked={lucaBuildDll && lucaProxyChoice === 'version'} on:change={(e) => selectLucaProxy('version', e.target.checked)} />
                <span><strong>version.dll</strong><small>Proxy Luca standard · x64</small></span>
              </label>
              <label class:luca-proxy-selected={lucaBuildDll && lucaProxyChoice === 'winmm'} class="checkbox-label luca-proxy-option">
                <input type="checkbox" checked={lucaBuildDll && lucaProxyChoice === 'winmm'} on:change={(e) => selectLucaProxy('winmm', e.target.checked)} />
                <span><strong>winmm.dll</strong><small>Little Busters! · PE32/x86</small></span>
              </label>
            </div>
            <div class="form-hint">{t("Choix exclusif : cocher une DLL décoche automatiquement l'autre. Décochez la sélection active pour générer le kit sans compiler.", 'Exclusive choice: selecting one DLL automatically clears the other. Clear the active selection to generate the kit without compiling.', uiLanguage)}</div>
            {#if lucaProxyChoice === 'winmm'}
              <div class="form-hint form-hint-warn"><strong>{t('LBEE sélectionné :', 'LBEE selected:', uiLanguage)}</strong> {t('la GUI compilera', 'the GUI will compile', uiLanguage)} <code>winmm.dll</code> {t('en 32 bits. Installez uniquement cette DLL à côté de', 'as 32-bit. Install only this DLL next to', uiLanguage)} <code>LITBUS_WIN32.exe</code> ; {t("n'utilisez pas", 'do not use', uiLanguage)} <code>version.dll</code>.</div>
            {/if}
          </div>

          <div class="form-group">
            <label>{t('Langue à injecter :', 'Language to inject:', uiLanguage)}</label>
            <div class="form-row checkbox-row luca-toolbar">
              <select value={lucaFillMode} on:change={(e) => setLucaFillMode(e.target.value)}>
                <option value="fr">FR</option>
                <option value="fr-safe">{t('FR (sûr)', 'FR (safe)', uiLanguage)}</option>
                <option value="en">ENG</option>
                <option value="en-safe">{t('ENG (sûr)', 'ENG (safe)', uiLanguage)}</option>
                <option value="ar">{t('Arabe', 'Arabic', uiLanguage)}</option>
                <option value="ru">{t('Russe (LBEE)', 'Russian (LBEE)', uiLanguage)}</option>
                <option value="jp">{t('Japonais', 'Japanese', uiLanguage)}</option>
                <option value="cn">{t('Chinois', 'Chinese', uiLanguage)}</option>
              </select>
              <button class="btn" on:click={clearLucaTargets}>{t('Vider', 'Clear', uiLanguage)}</button>
            </div>
            <div class="form-hint">{t('Les modes sûrs limitent la sélection aux chaînes communes aux quatre jeux et compatibles avec le budget du slot.', 'Safe modes limit selection to strings shared by all four games and compatible with the slot budget.', uiLanguage)}</div>
          </div>

          <div class="form-group">
            <label>{t('Filtre :', 'Filter:', uiLanguage)}</label>
            <div class="form-row">
              <input type="text" bind:value={lucaSearch} placeholder={t('source, cible, contexte...', 'source, target, context...', uiLanguage)} />
              <span class="luca-count">{lucaSelectedEntries().length} {t('sélectionnée(s)', 'selected', uiLanguage)} · {lucaVisibleEntries().length} {t('visible(s)', 'visible', uiLanguage)} · slot {slotLabel(lucaSlot)}</span>
            </div>
          </div>

          <div class="luca-table">
            <div class="luca-row luca-head">
              <div></div>
              <div>{t('Contexte', 'Context', uiLanguage)}</div>
              <div>Source</div>
              <div>{t('Cible', 'Target', uiLanguage)}</div>
              <div>Budget</div>
            </div>
            {#each lucaVisibleEntries() as entry (entry.rawOffset + entry.source)}
              <div class="luca-row" class:entry-warn={entryTooLong(entry)}>
                <div class="luca-check"><input type="checkbox" bind:checked={entry.include} /></div>
                <div>
                  <div class="luca-context">{entry.context}</div>
                  <div class="luca-meta">{entry.rawOffset} · {entry.encoding || 'utf-8'} · {entry.textKind}{entry.commonCount ? ` · ${entry.commonCount}/4` : ''}{entry.risk ? ` · ${entry.risk}` : ''}</div>
                </div>
                <div class="luca-source">{entry.source}</div>
                <div><input type="text" bind:value={entry.target} placeholder={entry.suggestedFr || t('Traduction', 'Translation', uiLanguage)} /></div>
                <div class="luca-budget">{lucaEncodedLen(entry)} / {entry.budget >= 0 ? entry.budget : '?'}</div>
              </div>
            {/each}
          </div>

          <div class="form-actions">
            {#if running}
              <span class="running-indicator"></span> Running...
            {:else}
              <button class="btn btn-primary" on:click={startLucaGenerate}
                disabled={!lucaExe || !lucaOutputDir || (lucaSelectedEntries().length === 0 && !lucaCustomPatch)}>
                {t('Générer le kit DLL', 'Generate DLL kit', uiLanguage)}
              </button>
            {/if}
          </div>
        {/if}

      <!-- IMAGE EXPORT -->
      {:else if selectedOp === 'image_export'}
        <div class="form-title">Image Export (CZ → PNG)</div>
        <div class="form-group">
          <div class="form-row checkbox-row">
            <label class="checkbox-label"><input type="checkbox" bind:checked={imgExpBatch} on:change={toggleExpBatch} /> Batch mode (entire folder)</label>
          </div>
        </div>
        <div class="form-group"><label>{imgExpBatch ? 'Input CZ folder:' : 'Input CZ file:'}</label><div class="form-row"><input type="text" bind:value={imgExpInput} readonly /><button class="btn" on:click={browseImgExpInput}>Select</button></div></div>
        <div class="form-group"><label>{imgExpBatch ? 'Output PNG folder:' : 'Output PNG file:'}</label><div class="form-row"><input type="text" bind:value={imgExpOutput} readonly /><button class="btn" on:click={browseImgExpOutput}>Select</button></div>
          {#if imgExpBatch}<div class="form-hint">CZ files are converted to PNG. MVT movies are exported as WebM.</div>{/if}
        </div>
        <div class="form-actions">{#if running}<span class="running-indicator"></span> Running...{:else}<button class="btn btn-primary" on:click={startImageExport} disabled={!imgExpInput || !imgExpOutput}>Start Export</button>{/if}</div>

      <!-- IMAGE IMPORT -->
      {:else if selectedOp === 'image_import'}
        <div class="form-title">Image Import (PNG → CZ)</div>
        <div class="form-group">
          <div class="form-row checkbox-row">
            <label class="checkbox-label"><input type="checkbox" bind:checked={imgImpBatch} on:change={toggleImpBatch} /> Batch mode (entire folder)</label>
            <label class="checkbox-label"><input type="checkbox" bind:checked={imgImpFill} /> Fill to original size (CZ1 only)</label>
          </div>
        </div>
        <div class="form-group"><label>{imgImpBatch ? 'Original CZ folder:' : 'Original CZ file:'}</label><div class="form-row"><input type="text" bind:value={imgImpSource} readonly /><button class="btn" on:click={browseImgImpSource}>Select</button></div><div class="form-hint">Original CZ file(s) for format reference</div></div>
        <div class="form-group"><label>{imgImpBatch ? 'Input PNG folder:' : 'Input PNG file:'}</label><div class="form-row"><input type="text" bind:value={imgImpInput} readonly /><button class="btn" on:click={browseImgImpInput}>Select</button></div></div>
        <div class="form-group"><label>{imgImpBatch ? 'Output CZ folder:' : 'Output CZ file:'}</label><div class="form-row"><input type="text" bind:value={imgImpOutput} readonly /><button class="btn" on:click={browseImgImpOutput}>Select</button></div>
          {#if imgImpBatch}<div class="form-hint">PNG files matching a CZ source will be converted</div>{/if}
        </div>
        <div class="form-actions">{#if running}<span class="running-indicator"></span> Running...{:else}<button class="btn btn-primary" on:click={startImageImport} disabled={!imgImpSource || !imgImpInput || !imgImpOutput}>Start Import</button>{/if}</div>

      <!-- DIALOGUE EXTRACT -->
      {:else if selectedOp === 'dlg_extract'}
        <div class="form-title">{t('Extraire les dialogues', 'Extract Dialogues', uiLanguage)}</div>
        <div class="form-hint" style="margin-bottom:10px">
          {t('Extrait les lignes', 'Extracts', uiLanguage)} <strong>MESSAGE</strong>, <strong>LOG_BEGIN</strong> {t('et', 'and', uiLanguage)} <strong>SELECT</strong> {t('des scripts décompilés (.txt) vers un fichier TSV éditable.', 'lines from decompiled scripts (.txt) to an editable TSV file.', uiLanguage)}<br>
          {t("Les colonnes correspondent aux chaînes entre guillemets dans l'ordre d'apparition. L'attribution des langues varie selon le jeu — vérifiez manuellement.", 'Columns match quoted strings in their order of appearance. Language assignment varies by game — check it manually.', uiLanguage)}
        </div>
        <div class="form-group">
          <div class="form-row checkbox-row">
            <label class="checkbox-label"><input type="checkbox" bind:checked={dlgExtBatch} on:change={toggleDlgExtBatch} /> {t('Mode par lot (dossier entier)', 'Batch mode (entire folder)', uiLanguage)}</label>
          </div>
        </div>
        <div class="form-group">
          <label>{t('Colonnes à extraire :', 'Columns to extract:', uiLanguage)}</label>
          <div class="form-row checkbox-row">
            <label class="checkbox-label"><input type="checkbox" bind:checked={dlgExtLang1} /> Lang 1</label>
            <label class="checkbox-label"><input type="checkbox" bind:checked={dlgExtLang2} /> Lang 2</label>
            <label class="checkbox-label"><input type="checkbox" bind:checked={dlgExtLang3} /> Lang 3</label>
            <label class="checkbox-label"><input type="checkbox" bind:checked={dlgExtLang4} /> Lang 4</label>
          </div>
          <div class="form-hint">{t('Chaque numéro correspond à la Nième chaîne entre guillemets dans le script. Ex. : pour AIR, Lang 1 = JAP, Lang 2 = ENG, Lang 3 = CN.', 'Each number matches the corresponding quoted string in the script. Example for AIR: Lang 1 = JAP, Lang 2 = ENG, Lang 3 = CN.', uiLanguage)}</div>
        </div>
        <div class="form-group"><label>{dlgExtBatch ? t('Dossier scripts (.txt) :', 'Scripts folder (.txt):', uiLanguage) : t('Fichier script (.txt) :', 'Script file (.txt):', uiLanguage)}</label><div class="form-row"><input type="text" bind:value={dlgExtInput} readonly /><button class="btn" on:click={browseDlgExtInput}>Select</button></div>
          {#if dlgExtDetectedFmt}<div class="form-hint">{t('Format détecté :', 'Detected format:', uiLanguage)} <strong>{dlgExtDetectedFmt}</strong></div>{/if}
        </div>
        <div class="form-group"><label>{dlgExtBatch ? t('Dossier de sortie :', 'Output folder:', uiLanguage) : t('Fichier TSV de sortie :', 'Output TSV file:', uiLanguage)}</label><div class="form-row"><input type="text" bind:value={dlgExtOutput} readonly /><button class="btn" on:click={browseDlgExtOutput}>Select</button></div>
          {#if dlgExtBatch}<div class="form-hint">{t('Un fichier', 'One', uiLanguage)} <code>*.ext.txt</code> {t('sera créé par script contenant des MESSAGE', 'file will be created for each script containing MESSAGE lines', uiLanguage)}</div>{/if}
        </div>
        <div class="form-actions">
          {#if running}<span class="running-indicator"></span> Running...
          {:else}<button class="btn btn-primary" on:click={startDlgExtract}
            disabled={!dlgExtInput || !dlgExtOutput || (!dlgExtLang1 && !dlgExtLang2 && !dlgExtLang3 && !dlgExtLang4)}>
            Start Extract
          </button>{/if}
        </div>

      <!-- DIALOGUE IMPORT -->
      {:else if selectedOp === 'dlg_import'}
        <div class="form-title">{t('Importer les dialogues', 'Import Dialogues', uiLanguage)}</div>
        <div class="form-hint" style="margin-bottom:10px">
          {t('Réinjecte les dialogues traduits (TSV) dans les fichiers scripts (.txt).', 'Reinserts translated dialogues (TSV) into script files (.txt).', uiLanguage)}<br>
          {t("Le TSV doit avoir été généré par l'extraction ci-dessus. Supporte MESSAGE, LOG_BEGIN et SELECT.", 'The TSV must have been generated by the extraction tool above. Supports MESSAGE, LOG_BEGIN, and SELECT.', uiLanguage)}
        </div>
        <div class="form-group">
          <div class="form-row checkbox-row">
            <label class="checkbox-label"><input type="checkbox" bind:checked={dlgImpBatch} on:change={toggleDlgImpBatch} /> {t('Mode par lot (dossier entier)', 'Batch mode (entire folder)', uiLanguage)}</label>
          </div>
        </div>
        <div class="form-group">
          <label>{t('Colonne cible à réinjecter :', 'Target column to reinsert:', uiLanguage)}</label>
          <div class="form-row">
            <select bind:value={dlgImpTargetCol}>
              <option value={1}>{t('Lang 1 (1re chaîne)', 'Lang 1 (1st string)', uiLanguage)}</option>
              <option value={2}>{t('Lang 2 (2e chaîne)', 'Lang 2 (2nd string)', uiLanguage)}</option>
              <option value={3}>{t('Lang 3 (3e chaîne)', 'Lang 3 (3rd string)', uiLanguage)}</option>
              <option value={4}>{t('Lang 4 (4e chaîne)', 'Lang 4 (4th string)', uiLanguage)}</option>
            </select>
          </div>
          <div class="form-hint">{t('La colonne sélectionnée sera lue dans le TSV et réinjectée dans la chaîne correspondante entre guillemets du script.', 'The selected column will be read from the TSV and reinserted into the corresponding quoted string in the script.', uiLanguage)}</div>
        </div>
        <div class="form-group"><label>{dlgImpBatch ? t('Dossier scripts originaux :', 'Original scripts folder:', uiLanguage) : t('Fichier script original :', 'Original script file:', uiLanguage)}</label><div class="form-row"><input type="text" bind:value={dlgImpScript} readonly /><button class="btn" on:click={browseDlgImpScript}>Select</button></div>
          <div class="form-hint">{t('Les fichiers .txt décompilés (originaux ou déjà traduits)', 'The decompiled .txt files (original or already translated)', uiLanguage)}</div>
        </div>
        <div class="form-group"><label>{dlgImpBatch ? t('Dossier TSV traduits :', 'Translated TSV folder:', uiLanguage) : t('Fichier TSV traduit :', 'Translated TSV file:', uiLanguage)}</label><div class="form-row"><input type="text" bind:value={dlgImpTsv} readonly /><button class="btn" on:click={browseDlgImpTsv}>Select</button></div>
          {#if dlgImpBatch}<div class="form-hint">{t('Fichiers', 'Files', uiLanguage)} <code>*.ext.txt</code> — {t('chaque TSV sera associé au script correspondant', 'each TSV will be matched with its corresponding script', uiLanguage)}</div>{/if}
        </div>
        <div class="form-group"><label>{dlgImpBatch ? t('Dossier de sortie :', 'Output folder:', uiLanguage) : t('Fichier de sortie :', 'Output file:', uiLanguage)}</label><div class="form-row"><input type="text" bind:value={dlgImpOutput} readonly /><button class="btn" on:click={browseDlgImpOutput}>Select</button></div></div>
        <div class="form-actions">
          {#if running}<span class="running-indicator"></span> Running...
          {:else}<button class="btn btn-primary" on:click={startDlgImport}
            disabled={!dlgImpScript || !dlgImpTsv || !dlgImpOutput}>
            Start Import
          </button>{/if}
        </div>

      <!-- ABOUT -->
      {:else if selectedOp === 'about'}
        <div class="form-title">{t('À propos', 'About', uiLanguage)}</div>
        <div class="about-panel">
          <div class="about-logo">LuckSystem</div>
          <div class="about-subtitle">Fork · Yoremi-v3.30</div>
          <div class="about-desc">
            {t("Interface graphique pour LuckSystem, l'outil de traduction de visual novels Visual Art's / Key.", "Graphical interface for LuckSystem, the Visual Art's / Key visual novel translation tool.", uiLanguage)}<br>
            {t('Inclut des correctifs CZ (CZ1, CZ4), script et PAK, ainsi que la gestion des processus.', 'Includes CZ (CZ1, CZ4), script, and PAK fixes, plus process management.', uiLanguage)}
          </div>
          <div class="about-links">
            <div class="about-link-row">
              <span class="about-link-label">{t('Projet source :', 'Upstream project:', uiLanguage)}</span>
              <span class="about-link-url">https://github.com/wetor/LuckSystem</span>
            </div>
            <div class="about-link-row">
              <span class="about-link-label">{t('Fork Yoremi :', 'Yoremi fork:', uiLanguage)}</span>
              <span class="about-link-url">https://github.com/yoremi-trad-fr/LuckSystem-2.3.2-Yoremi-Update</span>
            </div>
          </div>
          <div class="about-version">v3.30 GUI · Wails + Svelte</div>
        </div>
      {/if}
    </div>
  </div>

  <!-- CONSOLE -->
  <div class="console-wrapper">
    <div class="console-header">
      <span>{t('Sortie de la console', 'Console Output', uiLanguage)}</span>
      <div style="display:flex;gap:6px;align-items:center">
        {#if running}
          <button class="console-stop" on:click={stopProcess}>■ {t('Arrêter', 'Stop', uiLanguage)}</button>
        {/if}
        <button class="console-clear" on:click={clearConsole}>{t('Effacer', 'Clear', uiLanguage)}</button>
      </div>
    </div>
    <div
      class="console"
      bind:this={consoleEl}
      role="textbox"
      aria-label={t('Sortie de la console', 'Console Output', uiLanguage)}
      aria-readonly="true"
      tabindex="0"
      on:contextmenu={openConsoleMenu}
      on:keydown={handleConsoleKeydown}
    >
      {#each consoleLines as line}<div class={line.cls}>{line.text}</div>{/each}
    </div>
  </div>

  {#if consoleMenuVisible}
    <div
      class="console-menu"
      role="menu"
      style="left: {consoleMenuX}px; top: {consoleMenuY}px;"
    >
      <button type="button" on:click={copyConsoleSelection}>{t('Copier la sélection', 'Copy selection', uiLanguage)}</button>
      <button type="button" on:click={copyConsoleAll}>{t('Copier tout', 'Copy all', uiLanguage)}</button>
      <button type="button" on:click={pasteConsoleClipboard}>{t('Coller', 'Paste', uiLanguage)}</button>
    </div>
  {/if}
</div>
