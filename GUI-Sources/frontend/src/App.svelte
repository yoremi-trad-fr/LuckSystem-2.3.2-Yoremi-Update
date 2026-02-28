<script>
  import { onMount, onDestroy, tick } from 'svelte';
  import { EventsOn, EventsOff } from '../wailsjs/runtime/runtime.js';
  import {
    GetLuckSystemPath,
    SetLuckSystemPath,
    SelectPakFile,
    SelectFile,
    SelectDirectory,
    SelectSaveFile,
    StopProcess,
    ScriptDecompile,
    ScriptCompile,
    PakExtract,
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
  let lsPath = '';

  // --- Script fields ---
  let pakFile = '';
  let opcodeFile = '';
  let pluginFile = '';
  let charsetVal = 'UTF-8';
  let outputDir = '';
  let importDir = '';
  let outputPak = '';

  // --- PAK fields ---
  let pakExtSource = '';
  let pakExtOutput = '';
  let pakRepSource = '';
  let pakRepListFile = '';
  let pakRepInput = '';
  let pakRepOutput = '';
  let pakRepUseList = true; // mode par défaut : fichier liste

  // --- PAK Font fields ---
  let pakFontExtSource = '';
  let pakFontExtCharset = 'UTF-8';
  let pakFontExtOutput = '';
  let pakFontRepSource = '';
  let pakFontRepCharset = 'UTF-8';
  let pakFontRepListFile = '';
  let pakFontRepInput = '';
  let pakFontRepOutput = '';
  let pakFontRepUseList = true; // mode par défaut : fichier liste

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
    { id: 'decompile', label: 'Script Decompile' },
    { id: 'compile', label: 'Script Compile' },
    { id: '_s2', label: 'PAK (CG)', section: true },
    { id: 'pak_cg_extract', label: 'CG Extract' },
    { id: 'pak_cg_replace', label: 'CG Replace' },
    { id: '_s2b', label: 'PAK (Font)', section: true },
    { id: 'pak_font_extract', label: 'Font Extract' },
    { id: 'pak_font_replace', label: 'Font Replace' },
    { id: '_s3', label: 'FONT', section: true },
    { id: 'font_extract', label: 'Font Extract' },
    { id: 'font_edit', label: 'Font Edit' },
    { id: '_s4', label: 'IMAGE', section: true },
    { id: 'image_export', label: 'Image Export' },
    { id: 'image_import', label: 'Image Import' },
    { id: '_s6', label: 'DIALOGUE', section: true },
    { id: 'dlg_extract', label: 'Extract Dialogues' },
    { id: 'dlg_import', label: 'Import Dialogues' },
    { id: '_s5', label: '', section: true },
    { id: 'about', label: 'À propos' },
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

  onMount(async () => {
    EventsOn('log', (msg) => addLine(msg));
    lsPath = await GetLuckSystemPath();
    if (lsPath) {
      addLine('LuckSystem 2.3.2 - Yoremi fork v3');
      addLine('Executable: ' + lsPath);
    } else {
      addLine('[ERROR] lucksystem.exe not found!');
      addLine('Place lucksystem.exe next to the GUI, or click "Locate" below.');
    }
    addLine('Ready.');
  });
  onDestroy(() => { EventsOff('log'); });

  // ===== Browse helpers =====
  async function browsePak() { const f = await SelectPakFile(); if (f) pakFile = f; }
  async function browseOpcode() { const f = await SelectFile('Select Opcode (.txt)', '*.txt', 'Opcode files'); if (f) opcodeFile = f; }
  async function browsePlugin() { const f = await SelectFile('Select Plugin (.py)', '*.py', 'Python plugins'); if (f) pluginFile = f; }
  async function browseOutputDir() { const d = await SelectDirectory('Select output directory'); if (d) outputDir = d; }
  async function browseImportDir() { const d = await SelectDirectory('Select translated scripts directory'); if (d) importDir = d; }
  async function browseOutputPak() { const f = await SelectSaveFile('Save output PAK', 'SCRIPT_FR.PAK', '*.PAK;*.pak', 'PAK files'); if (f) outputPak = f; }

  async function browsePakExtSource() { const f = await SelectPakFile(); if (f) pakExtSource = f; }
  async function browsePakExtOutput() { const d = await SelectDirectory('Select extraction output'); if (d) pakExtOutput = d; }
  async function browsePakRepSource() { const f = await SelectPakFile(); if (f) pakRepSource = f; }
  async function browsePakRepListFile() { const f = await SelectFile('Sélectionner le fichier liste (_list.txt)', '*.txt', 'Fichiers liste'); if (f) pakRepListFile = f; }
  async function browsePakRepInput() { const d = await SelectDirectory('Select folder with modified files'); if (d) pakRepInput = d; }
  async function browsePakRepOutput() { const f = await SelectSaveFile('Save output PAK', 'FONT.out.PAK', '*.PAK;*.pak', 'PAK files'); if (f) pakRepOutput = f; }

  async function browsePakFontExtSource() { const f = await SelectPakFile(); if (f) pakFontExtSource = f; }
  async function browsePakFontExtOutput() { const d = await SelectDirectory('Dossier d\'extraction'); if (d) pakFontExtOutput = d; }
  async function browsePakFontRepSource() { const f = await SelectPakFile(); if (f) pakFontRepSource = f; }
  async function browsePakFontRepListFile() { const f = await SelectFile('Sélectionner le fichier liste (_list.txt)', '*.txt', 'Fichiers liste'); if (f) pakFontRepListFile = f; }
  async function browsePakFontRepInput() { const d = await SelectDirectory('Dossier des fichiers modifiés'); if (d) pakFontRepInput = d; }
  async function browsePakFontRepOutput() { const f = await SelectSaveFile('Save output PAK', 'FONT.out.PAK', '*.PAK;*.pak', 'PAK files'); if (f) pakFontRepOutput = f; }

  async function browseFontExtCz() { const f = await SelectFile('Select font CZ file', '*.*', 'Font CZ files'); if (f) fontExtCz = f; }
  async function browseFontExtInfo() { const f = await SelectFile('Select info file', '*.*', 'Info files'); if (f) fontExtInfo = f; }
  async function browseFontExtPng() { const f = await SelectSaveFile('Save font PNG', 'font.png', '*.png', 'PNG Images'); if (f) fontExtPng = f; }
  async function browseFontExtCharset() { const f = await SelectSaveFile('Save charset TXT', 'charset.txt', '*.txt', 'Text files'); if (f) fontExtCharset = f; }
  async function browseFontEditCz() { const f = await SelectFile('Select source CZ', '*.*', 'Font CZ files'); if (f) fontEditCz = f; }
  async function browseFontEditInfo() { const f = await SelectFile('Select source info', '*.*', 'Info files'); if (f) fontEditInfo = f; }
  async function browseFontEditTtf() { const f = await SelectFile('Select TTF font', '*.ttf;*.otf', 'Font files'); if (f) fontEditTtf = f; }
  async function browseFontEditOutCz() { const d = await SelectDirectory('Dossier de sortie pour le CZ modifié'); if (d) fontEditOutCz = d + '\\'; }
  async function browseFontEditOutInfo() { const d = await SelectDirectory('Dossier de sortie pour le fichier info'); if (d) fontEditOutInfo = d + '\\'; }
  async function browseFontEditCharset() { const f = await SelectFile('Select charset file', '*.txt', 'Text files'); if (f) fontEditCharsetFile = f; }

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
    if (lsPath) addLine('Executable set: ' + lsPath);
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

  function startDecompile() { run(() => ScriptDecompile(pakFile, opcodeFile, pluginFile, charsetVal, outputDir)); }
  function startCompile() { run(() => ScriptCompile(pakFile, opcodeFile, pluginFile, charsetVal, importDir, outputPak)); }
  function startPakExtract() { run(() => PakExtract(pakExtSource, pakExtOutput)); }
  function startPakReplace() {
    const listArg = pakRepUseList ? pakRepListFile : '';
    const dirArg  = pakRepUseList ? '' : pakRepInput;
    run(() => PakReplace(pakRepSource, dirArg, listArg, pakRepOutput));
  }
  function startPakFontExtract() { run(() => PakFontExtract(pakFontExtSource, pakFontExtCharset, pakFontExtOutput)); }
  function startPakFontReplace() {
    const listArg = pakFontRepUseList ? pakFontRepListFile : '';
    const dirArg  = pakFontRepUseList ? '' : pakFontRepInput;
    run(() => PakFontReplace(pakFontRepSource, pakFontRepCharset, dirArg, listArg, pakFontRepOutput));
  }
  function startFontExtract() { run(() => FontExtract(fontExtCz, fontExtInfo, fontExtPng, fontExtCharset)); }
  function startFontEdit() {
    const redraw  = fontEditMode === 'redraw';
    const append  = fontEditMode === 'append';
    const index   = (fontEditMode === 'insert') ? fontEditIndex : 0;
    run(() => FontEdit(fontEditCz, fontEditInfo, fontEditTtf, fontEditOutCz, fontEditOutInfo, fontEditCharsetFile, redraw, append, index));
  }

  function startImageExport() {
    if (imgExpBatch) run(() => ImageBatchExport(imgExpInput, imgExpOutput));
    else run(() => ImageExport(imgExpInput, imgExpOutput));
  }
  function startImageImport() {
    if (imgImpBatch) run(() => ImageBatchImport(imgImpSource, imgImpInput, imgImpOutput, imgImpFill));
    else run(() => ImageImport(imgImpSource, imgImpInput, imgImpOutput, imgImpFill));
  }

  function selectOp(op) { if (!op.disabled && !op.section) selectedOp = op.id; }

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
    <span>LuckSystem 2.3.2 - Yoremi fork v3</span>
    <span class="titlebar-path" on:click={locateLuckSystem} title="Click to change">
      {#if lsPath}📁 {lsPath}{:else}⚠ lucksystem.exe not found - Click to locate{/if}
    </span>
  </div>

  <div class="content">
    <!-- LEFT SIDEBAR -->
    <div class="sidebar">
      <div class="sidebar-title">Select option:</div>
      <div class="sidebar-list">
        {#each operations as op}
          {#if op.section}
            <div class="sidebar-section">{op.label}</div>
          {:else}
            <div class="sidebar-item" class:active={selectedOp === op.id} class:disabled={op.disabled} on:click={() => selectOp(op)}>
              {op.label}
            </div>
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
        <div class="form-group"><label>Opcode file (.txt):</label><div class="form-row"><input type="text" bind:value={opcodeFile} placeholder="e.g. data/AIR.txt" readonly /><button class="btn" on:click={browseOpcode}>Select</button></div><div class="form-hint">Game opcode definitions (AIR.txt, KANON.txt...)</div></div>
        <div class="form-group"><label>Plugin file (.py):</label><div class="form-row"><input type="text" bind:value={pluginFile} placeholder="e.g. data/AIR.py" readonly /><button class="btn" on:click={browsePlugin}>Select</button></div></div>
        <div class="form-group"><label>Charset:</label><div class="form-row"><select bind:value={charsetVal}><option value="UTF-8">UTF-8</option><option value="ShiftJIS">Shift-JIS</option><option value="GBK">GBK</option></select></div></div>
        <div class="form-group"><label>Output folder:</label><div class="form-row"><input type="text" bind:value={outputDir} readonly /><button class="btn" on:click={browseOutputDir}>Select</button></div><div class="form-hint">A SCRIPT.PAK subfolder will be created automatically inside</div></div>
        <div class="form-actions">{#if running}<span class="running-indicator"></span> Running...{:else}<button class="btn btn-primary" on:click={startDecompile} disabled={!pakFile || !outputDir}>Start Decompile</button>{/if}</div>

      <!-- SCRIPT COMPILE -->
      {:else if selectedOp === 'compile'}
        <div class="form-title">Script Compile</div>
        <div class="form-group"><label>Original SCRIPT.PAK:</label><div class="form-row"><input type="text" bind:value={pakFile} readonly /><button class="btn" on:click={browsePak}>Select</button></div><div class="form-hint">The original unmodified SCRIPT.PAK</div></div>
        <div class="form-group"><label>Opcode file (.txt):</label><div class="form-row"><input type="text" bind:value={opcodeFile} readonly /><button class="btn" on:click={browseOpcode}>Select</button></div></div>
        <div class="form-group"><label>Plugin file (.py):</label><div class="form-row"><input type="text" bind:value={pluginFile} readonly /><button class="btn" on:click={browsePlugin}>Select</button></div></div>
        <div class="form-group"><label>Charset:</label><div class="form-row"><select bind:value={charsetVal}><option value="UTF-8">UTF-8</option><option value="ShiftJIS">Shift-JIS</option><option value="GBK">GBK</option></select></div></div>
        <div class="form-group"><label>Translated scripts folder:</label><div class="form-row"><input type="text" bind:value={importDir} readonly /><button class="btn" on:click={browseImportDir}>Select</button></div><div class="form-hint">Select the parent folder containing SCRIPT.PAK (e.g. TRAD), not SCRIPT.PAK itself</div></div>
        <div class="form-group"><label>Output PAK file:</label><div class="form-row"><input type="text" bind:value={outputPak} readonly /><button class="btn" on:click={browseOutputPak}>Select</button></div></div>
        <div class="form-actions">{#if running}<span class="running-indicator"></span> Running...{:else}<button class="btn btn-primary" on:click={startCompile} disabled={!pakFile || !importDir || !outputPak}>Start Compile</button>{/if}</div>

      <!-- PAK CG EXTRACT -->
      {:else if selectedOp === 'pak_cg_extract'}
        <div class="form-title">PAK (CG) — Extract</div>
        <div class="form-group"><label>PAK file (CG) :</label><div class="form-row"><input type="text" bind:value={pakExtSource} readonly /><button class="btn" on:click={browsePakExtSource}>Select</button></div></div>
        <div class="form-group"><label>Output folder:</label><div class="form-row"><input type="text" bind:value={pakExtOutput} readonly /><button class="btn" on:click={browsePakExtOutput}>Select</button></div><div class="form-hint">Le fichier liste <code>&lt;NOM&gt;_list.txt</code> sera généré automatiquement dans ce dossier</div></div>
        <div class="form-actions">{#if running}<span class="running-indicator"></span> Running...{:else}<button class="btn btn-primary" on:click={startPakExtract} disabled={!pakExtSource || !pakExtOutput}>Start Extract</button>{/if}</div>

      <!-- PAK CG REPLACE -->
      {:else if selectedOp === 'pak_cg_replace'}
        <div class="form-title">PAK (CG) — Replace</div>
        <div class="form-group"><label>Original PAK file:</label><div class="form-row"><input type="text" bind:value={pakRepSource} readonly /><button class="btn" on:click={browsePakRepSource}>Select</button></div></div>
        <div class="form-group">
          <label>Mode d'entrée :</label>
          <div class="form-row checkbox-row" style="margin-bottom:6px">
            <label class="checkbox-label"><input type="radio" bind:group={pakRepUseList} value={true} /> Fichier liste (<code>*_list.txt</code>)</label>
            <label class="checkbox-label"><input type="radio" bind:group={pakRepUseList} value={false} /> Dossier de fichiers</label>
          </div>
          {#if pakRepUseList}
            <div class="form-row"><input type="text" bind:value={pakRepListFile} placeholder="SYSCG_list.txt" readonly /><button class="btn" on:click={browsePakRepListFile}>Select</button></div>
            <div class="form-hint">Fichier liste généré lors de l'extraction (ex : SYSCG_list.txt)</div>
          {:else}
            <div class="form-row"><input type="text" bind:value={pakRepInput} readonly /><button class="btn" on:click={browsePakRepInput}>Select</button></div>
            <div class="form-hint">Dossier contenant les fichiers modifiés à réinjecter</div>
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
        <div class="form-group"><label>Output folder :</label><div class="form-row"><input type="text" bind:value={pakFontExtOutput} readonly /><button class="btn" on:click={browsePakFontExtOutput}>Select</button></div><div class="form-hint">Tous les fichiers du PAK seront extraits ici</div></div>
        <div class="form-actions">{#if running}<span class="running-indicator"></span> Running...{:else}<button class="btn btn-primary" on:click={startPakFontExtract} disabled={!pakFontExtSource || !pakFontExtOutput}>Start Extract</button>{/if}</div>

      <!-- PAK FONT REPLACE -->
      {:else if selectedOp === 'pak_font_replace'}
        <div class="form-title">PAK (Font) — Replace</div>
        <div class="form-group"><label>Original PAK file (Font) :</label><div class="form-row"><input type="text" bind:value={pakFontRepSource} readonly /><button class="btn" on:click={browsePakFontRepSource}>Select</button></div></div>
        <div class="form-group"><label>Charset :</label><div class="form-row"><select bind:value={pakFontRepCharset}><option value="UTF-8">UTF-8</option><option value="ShiftJIS">Shift-JIS</option><option value="GBK">GBK</option></select></div></div>
        <div class="form-group">
          <label>Mode d'entrée :</label>
          <div class="form-row checkbox-row" style="margin-bottom:6px">
            <label class="checkbox-label"><input type="radio" bind:group={pakFontRepUseList} value={true} /> Fichier liste (<code>*_list.txt</code>)</label>
            <label class="checkbox-label"><input type="radio" bind:group={pakFontRepUseList} value={false} /> Dossier de fichiers</label>
          </div>
          {#if pakFontRepUseList}
            <div class="form-row"><input type="text" bind:value={pakFontRepListFile} placeholder="FONT__INFO_list.txt" readonly /><button class="btn" on:click={browsePakFontRepListFile}>Select</button></div>
            <div class="form-hint">Fichier liste généré lors de l'extraction (ex : FONT__INFO_list.txt)</div>
          {:else}
            <div class="form-row"><input type="text" bind:value={pakFontRepInput} readonly /><button class="btn" on:click={browsePakFontRepInput}>Select</button></div>
            <div class="form-hint">⚠ Le mode dossier peut échouer selon lucksystem — préférer le fichier liste</div>
          {/if}
        </div>
        <div class="form-group"><label>Output PAK :</label><div class="form-row"><input type="text" bind:value={pakFontRepOutput} readonly /><button class="btn" on:click={browsePakFontRepOutput}>Select</button></div></div>
        <div class="form-actions">
          {#if running}
            <span class="running-indicator"></span> Running...
          {:else}
            <button class="btn btn-primary" on:click={startPakFontReplace}
              disabled={!pakFontRepSource || !pakFontRepOutput || (pakFontRepUseList ? !pakFontRepListFile : !pakFontRepInput)}>
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
        <div class="form-title">Font Edit — Modification de glyphes</div>
        <div class="form-hint form-hint-warn">⚠ Font Edit modifie les glyphes d'un fichier CZ avec un TTF. Pour simplement re-packer un PAK de font, utilisez <strong>PAK (Font) → Font Replace</strong>.</div>

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
            <div class="form-hint">Redessine TOUS les glyphes existants avec le TTF. Aucun charset requis.</div>
          {:else if fontEditMode === 'append'}
            <div class="form-hint">Ajoute les caractères du charset à la fin de la police.</div>
          {:else if fontEditMode === 'insert'}
            <div class="form-row" style="margin-top:4px">
              <span style="min-width:90px;font-size:12px">Start index :</span>
              <input type="number" bind:value={fontEditIndex} min="0" style="width:80px;height:26px;padding:0 6px;border:1px solid #c0c0c0;border-radius:2px" />
            </div>
            <div class="form-hint">Insère/remplace à partir de cette position (0-indexé).</div>
          {/if}
        </div>

        {#if fontEditMode !== 'redraw'}
          <div class="form-group"><label>Charset file <span class="required">*</span> :</label><div class="form-row"><input type="text" bind:value={fontEditCharsetFile} readonly /><button class="btn" on:click={browseFontEditCharset}>Select</button></div><div class="form-hint">Fichier texte listant les caractères à ajouter/insérer (ex : accents_fr.txt)</div></div>
        {/if}

        <div class="form-group"><label>Output CZ <span class="required">*</span> :</label><div class="form-row"><input type="text" bind:value={fontEditOutCz} placeholder="ex: C:\dossier\ゴシック26" /><button class="btn" on:click={browseFontEditOutCz}>📁</button></div><div class="form-hint">Tapez le chemin complet sans extension — le bouton sélectionne le dossier</div></div>
        <div class="form-group"><label>Output info <span class="required">*</span> :</label><div class="form-row"><input type="text" bind:value={fontEditOutInfo} placeholder="ex: C:\dossier\info26" /><button class="btn" on:click={browseFontEditOutInfo}>📁</button></div><div class="form-hint">Tapez le chemin complet sans extension — requis pour mettre à jour le compte de caractères</div></div>

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
          {#if imgExpBatch}<div class="form-hint">All CZ files will be converted to PNG</div>{/if}
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
        <div class="form-title">Extract Dialogues</div>
        <div class="form-hint" style="margin-bottom:10px">
          Extrait les lignes <strong>MESSAGE</strong> et <strong>LOG_BEGIN</strong> des scripts décompilés (.txt) vers un fichier TSV éditable.<br>
          Les colonnes correspondent aux chaînes entre guillemets dans l'ordre d'apparition. L'attribution des langues varie selon le jeu — vérifiez manuellement.
        </div>
        <div class="form-group">
          <div class="form-row checkbox-row">
            <label class="checkbox-label"><input type="checkbox" bind:checked={dlgExtBatch} on:change={toggleDlgExtBatch} /> Batch mode (dossier entier)</label>
          </div>
        </div>
        <div class="form-group">
          <label>Colonnes à extraire :</label>
          <div class="form-row checkbox-row">
            <label class="checkbox-label"><input type="checkbox" bind:checked={dlgExtLang1} /> Lang 1</label>
            <label class="checkbox-label"><input type="checkbox" bind:checked={dlgExtLang2} /> Lang 2</label>
            <label class="checkbox-label"><input type="checkbox" bind:checked={dlgExtLang3} /> Lang 3</label>
            <label class="checkbox-label"><input type="checkbox" bind:checked={dlgExtLang4} /> Lang 4</label>
          </div>
          <div class="form-hint">Chaque numéro correspond à la Nième chaîne entre guillemets dans le script. Ex: pour AIR, Lang 1 = JAP, Lang 2 = ENG, Lang 3 = CN.</div>
        </div>
        <div class="form-group"><label>{dlgExtBatch ? 'Dossier scripts (.txt) :' : 'Fichier script (.txt) :'}</label><div class="form-row"><input type="text" bind:value={dlgExtInput} readonly /><button class="btn" on:click={browseDlgExtInput}>Select</button></div>
          {#if dlgExtDetectedFmt}<div class="form-hint">Format détecté : <strong>{dlgExtDetectedFmt}</strong></div>{/if}
        </div>
        <div class="form-group"><label>{dlgExtBatch ? 'Dossier de sortie :' : 'Fichier TSV de sortie :'}</label><div class="form-row"><input type="text" bind:value={dlgExtOutput} readonly /><button class="btn" on:click={browseDlgExtOutput}>Select</button></div>
          {#if dlgExtBatch}<div class="form-hint">Un fichier <code>*.ext.txt</code> sera créé par script contenant des MESSAGE</div>{/if}
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
        <div class="form-title">Import Dialogues</div>
        <div class="form-hint" style="margin-bottom:10px">
          Réinjecte les dialogues traduits (TSV) dans les fichiers scripts (.txt).<br>
          Le TSV doit avoir été généré par l'extraction ci-dessus. Supporte MESSAGE et LOG_BEGIN.
        </div>
        <div class="form-group">
          <div class="form-row checkbox-row">
            <label class="checkbox-label"><input type="checkbox" bind:checked={dlgImpBatch} on:change={toggleDlgImpBatch} /> Batch mode (dossier entier)</label>
          </div>
        </div>
        <div class="form-group">
          <label>Colonne cible à réinjecter :</label>
          <div class="form-row">
            <select bind:value={dlgImpTargetCol}>
              <option value={1}>Lang 1 (1ère chaîne)</option>
              <option value={2}>Lang 2 (2ème chaîne)</option>
              <option value={3}>Lang 3 (3ème chaîne)</option>
              <option value={4}>Lang 4 (4ème chaîne)</option>
            </select>
          </div>
          <div class="form-hint">La colonne sélectionnée sera lue dans le TSV et réinjectée dans la Nième chaîne entre guillemets du script.</div>
        </div>
        <div class="form-group"><label>{dlgImpBatch ? 'Dossier scripts originaux :' : 'Fichier script original :'}</label><div class="form-row"><input type="text" bind:value={dlgImpScript} readonly /><button class="btn" on:click={browseDlgImpScript}>Select</button></div>
          <div class="form-hint">Les fichiers .txt décompilés (originaux ou déjà traduits)</div>
        </div>
        <div class="form-group"><label>{dlgImpBatch ? 'Dossier TSV traduits :' : 'Fichier TSV traduit :'}</label><div class="form-row"><input type="text" bind:value={dlgImpTsv} readonly /><button class="btn" on:click={browseDlgImpTsv}>Select</button></div>
          {#if dlgImpBatch}<div class="form-hint">Fichiers <code>*.ext.txt</code> — chaque TSV sera associé au script correspondant</div>{/if}
        </div>
        <div class="form-group"><label>{dlgImpBatch ? 'Dossier de sortie :' : 'Fichier de sortie :'}</label><div class="form-row"><input type="text" bind:value={dlgImpOutput} readonly /><button class="btn" on:click={browseDlgImpOutput}>Select</button></div></div>
        <div class="form-actions">
          {#if running}<span class="running-indicator"></span> Running...
          {:else}<button class="btn btn-primary" on:click={startDlgImport}
            disabled={!dlgImpScript || !dlgImpTsv || !dlgImpOutput}>
            Start Import
          </button>{/if}
        </div>

      <!-- ABOUT -->
      {:else if selectedOp === 'about'}
        <div class="form-title">À propos</div>
        <div class="about-panel">
          <div class="about-logo">LuckSystem</div>
          <div class="about-subtitle">Fork · Yoremi-v3</div>
          <div class="about-desc">
            Interface graphique pour LuckSystem, l'outil de traduction de visual novels Visual Art's / Key.<br>
            Inclut des correctifs CZ (CZ1, CZ4), script, PAK, et une interface subprocess.
          </div>
          <div class="about-links">
            <div class="about-link-row">
              <span class="about-link-label">Projet source :</span>
              <span class="about-link-url">https://github.com/wetor/LuckSystem</span>
            </div>
            <div class="about-link-row">
              <span class="about-link-label">Fork Yoremi :</span>
              <span class="about-link-url">https://github.com/yoremi-trad-fr/LuckSystem-2.3.2-Yoremi-Update</span>
            </div>
          </div>
          <div class="about-version">v3 GUI · Wails + Svelte</div>
        </div>
      {/if}
    </div>
  </div>

  <!-- CONSOLE -->
  <div class="console-wrapper">
    <div class="console-header">
      <span>Console Output</span>
      <div style="display:flex;gap:6px;align-items:center">
        {#if running}
          <button class="console-stop" on:click={stopProcess}>■ Stop</button>
        {/if}
        <button class="console-clear" on:click={clearConsole}>Clear</button>
      </div>
    </div>
    <div class="console" bind:this={consoleEl}>
      {#each consoleLines as line}<div class={line.cls}>{line.text}</div>{/each}
    </div>
  </div>
</div>
