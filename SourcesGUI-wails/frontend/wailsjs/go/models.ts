export namespace main {
	
	export class DialogueFormatInfo {
	    format: string;
	    maxCols: number;
	
	    static createFrom(source: any = {}) {
	        return new DialogueFormatInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.format = source["format"];
	        this.maxCols = source["maxCols"];
	    }
	}
	export class GamePreset {
	    name: string;
	    opcodeFile: string;
	    pluginFile: string;
	    gameFlag: string;
	
	    static createFrom(source: any = {}) {
	        return new GamePreset(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.opcodeFile = source["opcodeFile"];
	        this.pluginFile = source["pluginFile"];
	        this.gameFlag = source["gameFlag"];
	    }
	}
	export class LucaMenuCommonRow {
	    source: string;
	    suggestedFr: string;
	    games: string[];
	    count: number;
	
	    static createFrom(source: any = {}) {
	        return new LucaMenuCommonRow(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.source = source["source"];
	        this.suggestedFr = source["suggestedFr"];
	        this.games = source["games"];
	        this.count = source["count"];
	    }
	}
	export class LucaMenuEntry {
	    rawOffset: string;
	    source: string;
	    target: string;
	    suggestedFr: string;
	    suggestedEn: string;
	    suggestedAr: string;
	    suggestedRu: string;
	    suggestedJp: string;
	    suggestedCn: string;
	    catalogId: string;
	    context: string;
	    note: string;
	    slot: string;
	    category: string;
	    textKind: string;
	    encoding: string;
	    sourceBytes: number;
	    targetBytes: number;
	    budget: number;
	    commonCount: number;
	    commonGames: string[];
	    safeAuto: boolean;
	    risk: string;
	    include: boolean;
	
	    static createFrom(source: any = {}) {
	        return new LucaMenuEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.rawOffset = source["rawOffset"];
	        this.source = source["source"];
	        this.target = source["target"];
	        this.suggestedFr = source["suggestedFr"];
	        this.suggestedEn = source["suggestedEn"];
	        this.suggestedAr = source["suggestedAr"];
	        this.suggestedRu = source["suggestedRu"];
	        this.suggestedJp = source["suggestedJp"];
	        this.suggestedCn = source["suggestedCn"];
	        this.catalogId = source["catalogId"];
	        this.context = source["context"];
	        this.note = source["note"];
	        this.slot = source["slot"];
	        this.category = source["category"];
	        this.textKind = source["textKind"];
	        this.encoding = source["encoding"];
	        this.sourceBytes = source["sourceBytes"];
	        this.targetBytes = source["targetBytes"];
	        this.budget = source["budget"];
	        this.commonCount = source["commonCount"];
	        this.commonGames = source["commonGames"];
	        this.safeAuto = source["safeAuto"];
	        this.risk = source["risk"];
	        this.include = source["include"];
	    }
	}
	export class LucaMenuPatchEdit {
	    rawOffset: string;
	    source: string;
	    target: string;
	    context: string;
	    note: string;
	    encoding: string;
	    include: boolean;
	    budget: number;
	
	    static createFrom(source: any = {}) {
	        return new LucaMenuPatchEdit(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.rawOffset = source["rawOffset"];
	        this.source = source["source"];
	        this.target = source["target"];
	        this.context = source["context"];
	        this.note = source["note"];
	        this.encoding = source["encoding"];
	        this.include = source["include"];
	        this.budget = source["budget"];
	    }
	}
	export class LucaMenuGenerateRequest {
	    profileId: string;
	    gameExe: string;
	    outputDir: string;
	    patchGameName: string;
	    patchVersion: string;
	    slot: string;
	    buildDll: boolean;
	    proxyDll: string;
	    preset: string;
	    customPatch: string;
	    entries: LucaMenuPatchEdit[];
	
	    static createFrom(source: any = {}) {
	        return new LucaMenuGenerateRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.profileId = source["profileId"];
	        this.gameExe = source["gameExe"];
	        this.outputDir = source["outputDir"];
	        this.patchGameName = source["patchGameName"];
	        this.patchVersion = source["patchVersion"];
	        this.slot = source["slot"];
	        this.buildDll = source["buildDll"];
	        this.proxyDll = source["proxyDll"];
	        this.preset = source["preset"];
	        this.customPatch = source["customPatch"];
	        this.entries = this.convertValues(source["entries"], LucaMenuPatchEdit);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class LucaMenuProfile {
	    id: string;
	    name: string;
	    folder: string;
	    patchFile: string;
	    gameExe: string;
	    rvaDelta: string;
	    rvaMode: string;
	    proxyDll: string;
	    architecture: string;
	    patchGameName: string;
	    patchVersion: string;
	    entries: LucaMenuEntry[];
	    totalCount: number;
	    englishCount: number;
	    japaneseCount: number;
	    chineseCount: number;
	    safeAutoCount: number;
	
	    static createFrom(source: any = {}) {
	        return new LucaMenuProfile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.folder = source["folder"];
	        this.patchFile = source["patchFile"];
	        this.gameExe = source["gameExe"];
	        this.rvaDelta = source["rvaDelta"];
	        this.rvaMode = source["rvaMode"];
	        this.proxyDll = source["proxyDll"];
	        this.architecture = source["architecture"];
	        this.patchGameName = source["patchGameName"];
	        this.patchVersion = source["patchVersion"];
	        this.entries = this.convertValues(source["entries"], LucaMenuEntry);
	        this.totalCount = source["totalCount"];
	        this.englishCount = source["englishCount"];
	        this.japaneseCount = source["japaneseCount"];
	        this.chineseCount = source["chineseCount"];
	        this.safeAutoCount = source["safeAutoCount"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class LucaMenuInventory {
	    kitDir: string;
	    profiles: LucaMenuProfile[];
	    common: LucaMenuCommonRow[];
	
	    static createFrom(source: any = {}) {
	        return new LucaMenuInventory(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.kitDir = source["kitDir"];
	        this.profiles = this.convertValues(source["profiles"], LucaMenuProfile);
	        this.common = this.convertValues(source["common"], LucaMenuCommonRow);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	

}

