export namespace persist {
	
	export class AttemptsSummary {
	    id: string;
	    templateId: string;
	    name: string;
	    categoryName: string;
	    attemptCount: number;
	    updatedAt: number;
	
	    static createFrom(source: any = {}) {
	        return new AttemptsSummary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.templateId = source["templateId"];
	        this.name = source["name"];
	        this.categoryName = source["categoryName"];
	        this.attemptCount = source["attemptCount"];
	        this.updatedAt = source["updatedAt"];
	    }
	}
	export class ColorSettings {
	    aheadGaining: string;
	    aheadLosing: string;
	    behindGaining: string;
	    behindLosing: string;
	    bestTime: string;
	
	    static createFrom(source: any = {}) {
	        return new ColorSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.aheadGaining = source["aheadGaining"];
	        this.aheadLosing = source["aheadLosing"];
	        this.behindGaining = source["behindGaining"];
	        this.behindLosing = source["behindLosing"];
	        this.bestTime = source["bestTime"];
	    }
	}
	export class HotkeyBindings {
	    startSplit: string;
	    pause: string;
	    reset: string;
	    undoSplit: string;
	    skipSplit: string;
	
	    static createFrom(source: any = {}) {
	        return new HotkeyBindings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.startSplit = source["startSplit"];
	        this.pause = source["pause"];
	        this.reset = source["reset"];
	        this.undoSplit = source["undoSplit"];
	        this.skipSplit = source["skipSplit"];
	    }
	}
	export class Settings {
	    alwaysOnTop: boolean;
	    hotkeys: HotkeyBindings;
	    comparison: string;
	    colors: ColorSettings;
	
	    static createFrom(source: any = {}) {
	        return new Settings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.alwaysOnTop = source["alwaysOnTop"];
	        this.hotkeys = this.convertValues(source["hotkeys"], HotkeyBindings);
	        this.comparison = source["comparison"];
	        this.colors = this.convertValues(source["colors"], ColorSettings);
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
	export class TemplateSummary {
	    id: string;
	    name: string;
	    segmentCount: number;
	    updatedAt: number;
	
	    static createFrom(source: any = {}) {
	        return new TemplateSummary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.segmentCount = source["segmentCount"];
	        this.updatedAt = source["updatedAt"];
	    }
	}

}

export namespace split {
	
	export class Attempt {
	    id: number;
	    // Go type: time
	    startedAt: any;
	    splitTimesMs: number[];
	    completed: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Attempt(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.startedAt = this.convertValues(source["startedAt"], null);
	        this.splitTimesMs = source["splitTimesMs"];
	        this.completed = source["completed"];
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
	export class Delta {
	    segmentIndex: number;
	    deltaMs: number;
	    isBestEver: boolean;
	    isAhead: boolean;
	    gainedTime: boolean;
	    skipped: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Delta(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.segmentIndex = source["segmentIndex"];
	        this.deltaMs = source["deltaMs"];
	        this.isBestEver = source["isBestEver"];
	        this.isAhead = source["isAhead"];
	        this.gainedTime = source["gainedTime"];
	        this.skipped = source["skipped"];
	    }
	}

}

