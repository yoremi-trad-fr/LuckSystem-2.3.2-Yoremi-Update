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

}

