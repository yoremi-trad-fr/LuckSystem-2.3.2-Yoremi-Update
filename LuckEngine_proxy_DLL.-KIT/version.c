/*
 * version.dll proxy for Luck Engine (Steam) — in-memory string patch toolkit
 *
 * How it works:
 *   1) The game exe imports VERSION.dll. Windows resolves imports from the exe's
 *      directory BEFORE System32 (VERSION.dll is not a Known DLL), so our
 *      proxy gets loaded instead of the real one.
 *   2) On DLL_PROCESS_ATTACH we spawn a worker thread that waits for the
 *      SteamStub DRM to finish decrypting / verifying .text + .rdata.
 *      Detection: poll a sentinel byte at a known .rdata offset until it
 *      matches the expected plaintext. SteamStub typically finishes within
 *      ~200 ms; we poll for up to ~30 s and give up silently on timeout.
 *   3) Apply the string patches via VirtualProtect + memcpy, restore page
 *      protection, done. The on-disk exe is never modified: SteamStub
 *      integrity checks pass, Steam is happy.
 *   4) Exports are runtime-forwarded to the real C:\Windows\System32\version.dll
 *      so the game's VERSION.dll imports resolve normally.
 *
 * Build:
 *   x86_64-w64-mingw32-gcc -O2 -s -shared -o version.dll \
 *       version.c version.def -static-libgcc -Wl,--subsystem,windows
 */

#include <windows.h>
#include <stdio.h>
#include <string.h>

/* ------------------------------------------------------------------ */
/*  Patch table                                                        */
/* ------------------------------------------------------------------ */

typedef struct {
    DWORD        rva;          /* RVA within the game module */
    SIZE_T       len;          /* byte count to check & write (src length + null) */
    const BYTE  *expected;     /* must match before we patch */
    const BYTE  *replacement;  /* translated bytes padded with \0 to src length */
    const char  *comment;      /* for logging                 */
} LuckPatch;

#include "patches.h"

/* Sentinel = the first patch: we use it to know when SteamStub is done */
#define SENTINEL_INDEX 0

/* ------------------------------------------------------------------ */
/*  Logging (optional, enabled if LUCKPROXY_LOG env var is set)        */
/* ------------------------------------------------------------------ */

static FILE *g_log = NULL;

static void log_open(void) {
    char buf[8];
    if (GetEnvironmentVariableA("LUCKPROXY_LOG", buf, sizeof(buf)) > 0) {
        char path[MAX_PATH];
        GetModuleFileNameA(NULL, path, MAX_PATH);
        char *slash = strrchr(path, '\\');
        if (slash) strcpy(slash + 1, "luckproxy.log");
        else strcpy(path, "luckproxy.log");
        g_log = fopen(path, "a");
    }
}

static void log_msg(const char *fmt, ...) {
    if (!g_log) return;
    SYSTEMTIME st; GetLocalTime(&st);
    fprintf(g_log, "[%02d:%02d:%02d.%03d] ",
            st.wHour, st.wMinute, st.wSecond, st.wMilliseconds);
    va_list ap; va_start(ap, fmt);
    vfprintf(g_log, fmt, ap);
    va_end(ap);
    fprintf(g_log, "\n");
    fflush(g_log);
}

/* ------------------------------------------------------------------ */
/*  Real VERSION.dll forwarding                                         */
/* ------------------------------------------------------------------ */

static HMODULE g_hRealDll = NULL;

/* Windows 10+ exports (superset; older Windows lack the *Ex variants, we
 * tolerate their absence by returning 0/FALSE) */
typedef DWORD (WINAPI *fn_GetFileVersionInfoSizeA)(LPCSTR, LPDWORD);
typedef DWORD (WINAPI *fn_GetFileVersionInfoSizeW)(LPCWSTR, LPDWORD);
typedef BOOL  (WINAPI *fn_GetFileVersionInfoA)(LPCSTR, DWORD, DWORD, LPVOID);
typedef BOOL  (WINAPI *fn_GetFileVersionInfoW)(LPCWSTR, DWORD, DWORD, LPVOID);
typedef BOOL  (WINAPI *fn_VerQueryValueA)(LPCVOID, LPCSTR, LPVOID*, PUINT);
typedef BOOL  (WINAPI *fn_VerQueryValueW)(LPCVOID, LPCWSTR, LPVOID*, PUINT);
typedef DWORD (WINAPI *fn_VerFindFileA)(DWORD, LPCSTR, LPCSTR, LPCSTR, LPSTR, PUINT, LPSTR, PUINT);
typedef DWORD (WINAPI *fn_VerFindFileW)(DWORD, LPCWSTR, LPCWSTR, LPCWSTR, LPWSTR, PUINT, LPWSTR, PUINT);
typedef DWORD (WINAPI *fn_VerInstallFileA)(DWORD, LPCSTR, LPCSTR, LPCSTR, LPCSTR, LPCSTR, LPSTR, PUINT);
typedef DWORD (WINAPI *fn_VerInstallFileW)(DWORD, LPCWSTR, LPCWSTR, LPCWSTR, LPCWSTR, LPCWSTR, LPWSTR, PUINT);
typedef DWORD (WINAPI *fn_VerLanguageNameA)(DWORD, LPSTR, DWORD);
typedef DWORD (WINAPI *fn_VerLanguageNameW)(DWORD, LPWSTR, DWORD);

static fn_GetFileVersionInfoSizeA p_GetFileVersionInfoSizeA;
static fn_GetFileVersionInfoSizeW p_GetFileVersionInfoSizeW;
static fn_GetFileVersionInfoA     p_GetFileVersionInfoA;
static fn_GetFileVersionInfoW     p_GetFileVersionInfoW;
static fn_VerQueryValueA          p_VerQueryValueA;
static fn_VerQueryValueW          p_VerQueryValueW;
static fn_VerFindFileA            p_VerFindFileA;
static fn_VerFindFileW            p_VerFindFileW;
static fn_VerInstallFileA         p_VerInstallFileA;
static fn_VerInstallFileW         p_VerInstallFileW;
static fn_VerLanguageNameA        p_VerLanguageNameA;
static fn_VerLanguageNameW        p_VerLanguageNameW;

static void load_real_version_dll(void) {
    if (g_hRealDll) return;
    char path[MAX_PATH];
    UINT n = GetSystemDirectoryA(path, MAX_PATH);
    if (!n || n >= MAX_PATH - 16) return;
    lstrcatA(path, "\\version.dll");
    g_hRealDll = LoadLibraryA(path);
    if (!g_hRealDll) { log_msg("ERROR: cannot load %s", path); return; }
    log_msg("Loaded real version.dll from %s", path);

#define LOAD(name) p_##name = (fn_##name)GetProcAddress(g_hRealDll, #name)
    LOAD(GetFileVersionInfoSizeA);
    LOAD(GetFileVersionInfoSizeW);
    LOAD(GetFileVersionInfoA);
    LOAD(GetFileVersionInfoW);
    LOAD(VerQueryValueA);
    LOAD(VerQueryValueW);
    LOAD(VerFindFileA);
    LOAD(VerFindFileW);
    LOAD(VerInstallFileA);
    LOAD(VerInstallFileW);
    LOAD(VerLanguageNameA);
    LOAD(VerLanguageNameW);
#undef LOAD
}

/* ------------------------------------------------------------------ */
/*  Exports (forward to real DLL)                                       */
/* ------------------------------------------------------------------ */

DWORD WINAPI LKPRX_GetFileVersionInfoSizeA(LPCSTR f, LPDWORD h) {
    load_real_version_dll();
    return p_GetFileVersionInfoSizeA ? p_GetFileVersionInfoSizeA(f, h) : 0;
}
DWORD WINAPI LKPRX_GetFileVersionInfoSizeW(LPCWSTR f, LPDWORD h) {
    load_real_version_dll();
    return p_GetFileVersionInfoSizeW ? p_GetFileVersionInfoSizeW(f, h) : 0;
}
BOOL WINAPI LKPRX_GetFileVersionInfoA(LPCSTR f, DWORD h, DWORD l, LPVOID d) {
    load_real_version_dll();
    return p_GetFileVersionInfoA ? p_GetFileVersionInfoA(f, h, l, d) : FALSE;
}
BOOL WINAPI LKPRX_GetFileVersionInfoW(LPCWSTR f, DWORD h, DWORD l, LPVOID d) {
    load_real_version_dll();
    return p_GetFileVersionInfoW ? p_GetFileVersionInfoW(f, h, l, d) : FALSE;
}
BOOL WINAPI LKPRX_VerQueryValueA(LPCVOID b, LPCSTR s, LPVOID *p, PUINT l) {
    load_real_version_dll();
    return p_VerQueryValueA ? p_VerQueryValueA(b, s, p, l) : FALSE;
}
BOOL WINAPI LKPRX_VerQueryValueW(LPCVOID b, LPCWSTR s, LPVOID *p, PUINT l) {
    load_real_version_dll();
    return p_VerQueryValueW ? p_VerQueryValueW(b, s, p, l) : FALSE;
}
DWORD WINAPI LKPRX_VerFindFileA(DWORD a, LPCSTR b, LPCSTR c, LPCSTR d, LPSTR e, PUINT f, LPSTR g, PUINT h) {
    load_real_version_dll();
    return p_VerFindFileA ? p_VerFindFileA(a,b,c,d,e,f,g,h) : 0;
}
DWORD WINAPI LKPRX_VerFindFileW(DWORD a, LPCWSTR b, LPCWSTR c, LPCWSTR d, LPWSTR e, PUINT f, LPWSTR g, PUINT h) {
    load_real_version_dll();
    return p_VerFindFileW ? p_VerFindFileW(a,b,c,d,e,f,g,h) : 0;
}
DWORD WINAPI LKPRX_VerInstallFileA(DWORD a, LPCSTR b, LPCSTR c, LPCSTR d, LPCSTR e, LPCSTR f, LPSTR g, PUINT h) {
    load_real_version_dll();
    return p_VerInstallFileA ? p_VerInstallFileA(a,b,c,d,e,f,g,h) : 0;
}
DWORD WINAPI LKPRX_VerInstallFileW(DWORD a, LPCWSTR b, LPCWSTR c, LPCWSTR d, LPCWSTR e, LPCWSTR f, LPWSTR g, PUINT h) {
    load_real_version_dll();
    return p_VerInstallFileW ? p_VerInstallFileW(a,b,c,d,e,f,g,h) : 0;
}
DWORD WINAPI LKPRX_VerLanguageNameA(DWORD l, LPSTR s, DWORD n) {
    load_real_version_dll();
    return p_VerLanguageNameA ? p_VerLanguageNameA(l, s, n) : 0;
}
DWORD WINAPI LKPRX_VerLanguageNameW(DWORD l, LPWSTR s, DWORD n) {
    load_real_version_dll();
    return p_VerLanguageNameW ? p_VerLanguageNameW(l, s, n) : 0;
}

/* ------------------------------------------------------------------ */
/*  Patch thread                                                        */
/* ------------------------------------------------------------------ */

static BOOL sentinel_ready(const BYTE *base) {
    const LuckPatch *s = &g_patches[SENTINEL_INDEX];
    return memcmp(base + s->rva, s->expected, s->len) == 0;
}

static void apply_patch(BYTE *base, const LuckPatch *p) {
    BYTE *addr = base + p->rva;
    DWORD old_protect = 0;
    if (!VirtualProtect(addr, p->len, PAGE_READWRITE, &old_protect)) {
        log_msg("VirtualProtect(RW) failed at RVA 0x%X: err=%lu",
                p->rva, GetLastError());
        return;
    }
    memcpy(addr, p->replacement, p->len);
    DWORD tmp;
    VirtualProtect(addr, p->len, old_protect, &tmp);

    /* Flush icache (belt and suspenders; data patches don't need this
     * but cheap insurance). */
    FlushInstructionCache(GetCurrentProcess(), addr, p->len);

    log_msg("Patched RVA 0x%X (%llu bytes): %s",
            p->rva, (unsigned long long)p->len, p->comment);
}

static DWORD WINAPI patch_thread(LPVOID unused) {
    (void)unused;
    HMODULE hExe = GetModuleHandleA(NULL);
    if (!hExe) { log_msg("ERROR: GetModuleHandleA(NULL) failed"); return 1; }
    BYTE *base = (BYTE *)hExe;
    log_msg("Module base = %p", (void*)base);

    /* Poll for SteamStub to finish. Typical: <1s. We wait up to 30s. */
    BOOL ready = FALSE;
    for (int tries = 0; tries < 300; ++tries) {
        if (sentinel_ready(base)) { ready = TRUE; break; }
        Sleep(100);
    }
    if (!ready) {
        const LuckPatch *s = &g_patches[SENTINEL_INDEX];
        BYTE got[8] = {0};
        memcpy(got, base + s->rva, s->len < 8 ? s->len : 8);
        log_msg("Sentinel never matched after 30s. Got: %02X %02X %02X %02X %02X",
                got[0], got[1], got[2], got[3], got[4]);
        return 2;
    }
    log_msg("Sentinel ready, applying %d patch(es)", (int)N_PATCHES);

    for (size_t i = 0; i < N_PATCHES; ++i)
        apply_patch(base, &g_patches[i]);

    log_msg("Patch thread done.");
    return 0;
}

/* ------------------------------------------------------------------ */
/*  DllMain                                                             */
/* ------------------------------------------------------------------ */

BOOL WINAPI DllMain(HINSTANCE hinst, DWORD reason, LPVOID reserved) {
    (void)reserved;
    if (reason == DLL_PROCESS_ATTACH) {
        DisableThreadLibraryCalls(hinst);
        log_open();
        log_msg("DLL_PROCESS_ATTACH (%s proxy v%s, %d patches)",
                PATCH_GAME_NAME, PATCH_VERSION, (int)N_PATCHES);
        HANDLE h = CreateThread(NULL, 0, patch_thread, NULL, 0, NULL);
        if (h) CloseHandle(h);
        else log_msg("ERROR: CreateThread failed: %lu", GetLastError());
    }
    return TRUE;
}
