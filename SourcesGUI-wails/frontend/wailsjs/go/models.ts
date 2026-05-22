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

}

