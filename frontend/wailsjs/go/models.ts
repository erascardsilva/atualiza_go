export namespace main {
	
	export class DistroInfo {
	    name: string;
	    id: string;
	    version: string;
	    family: string;
	    packageManager: string;
	    hasFlatpak: boolean;
	    hasSnap: boolean;
	
	    static createFrom(source: any = {}) {
	        return new DistroInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.id = source["id"];
	        this.version = source["version"];
	        this.family = source["family"];
	        this.packageManager = source["packageManager"];
	        this.hasFlatpak = source["hasFlatpak"];
	        this.hasSnap = source["hasSnap"];
	    }
	}
	export class SystemStats {
	    memTotal: number;
	    memUsed: number;
	    memPercent: number;
	    diskTotal: string;
	    diskUsed: string;
	    diskPercent: number;
	    diskMessage: string;
	
	    static createFrom(source: any = {}) {
	        return new SystemStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.memTotal = source["memTotal"];
	        this.memUsed = source["memUsed"];
	        this.memPercent = source["memPercent"];
	        this.diskTotal = source["diskTotal"];
	        this.diskUsed = source["diskUsed"];
	        this.diskPercent = source["diskPercent"];
	        this.diskMessage = source["diskMessage"];
	    }
	}
	export class UpdateStep {
	    id: string;
	    label: string;
	    command: string;
	    needRoot: boolean;
	
	    static createFrom(source: any = {}) {
	        return new UpdateStep(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.label = source["label"];
	        this.command = source["command"];
	        this.needRoot = source["needRoot"];
	    }
	}

}

