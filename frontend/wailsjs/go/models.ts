export namespace docker {
	
	export class DockerImage {
	    repository: string;
	    tag: string;
	    image_id: string;
	
	    static createFrom(source: any = {}) {
	        return new DockerImage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.repository = source["repository"];
	        this.tag = source["tag"];
	        this.image_id = source["image_id"];
	    }
	}

}

export namespace main {
	
	export class RepoDataPaginationFE {
	    data: models.Repository[];
	    total: number;
	
	    static createFrom(source: any = {}) {
	        return new RepoDataPaginationFE(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.data = this.convertValues(source["data"], models.Repository);
	        this.total = source["total"];
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
	export class ScanDataPaginationFE {
	    data: models.ScanFull[];
	    total: number;
	
	    static createFrom(source: any = {}) {
	        return new ScanDataPaginationFE(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.data = this.convertValues(source["data"], models.ScanFull);
	        this.total = source["total"];
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

export namespace models {
	
	export class Repository {
	    id: number;
	    user_id: number;
	    name: string;
	    description: string;
	    path: string;
	    create_at: string;
	    update_at: string;
	
	    static createFrom(source: any = {}) {
	        return new Repository(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.user_id = source["user_id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.path = source["path"];
	        this.create_at = source["create_at"];
	        this.update_at = source["update_at"];
	    }
	}
	export class RepositoryCreate {
	    user_id: number;
	    name: string;
	    description: string;
	    path: string;
	
	    static createFrom(source: any = {}) {
	        return new RepositoryCreate(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.user_id = source["user_id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.path = source["path"];
	    }
	}
	export class RepositoryUpdate {
	    id: number;
	    user_id: number;
	    name: string;
	    description: string;
	    path: string;
	
	    static createFrom(source: any = {}) {
	        return new RepositoryUpdate(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.user_id = source["user_id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.path = source["path"];
	    }
	}
	export class Scan {
	    id: number;
	    repository_id: number;
	    scan_time: string;
	    result: string;
	    vulnerabilities: number;
	    status: string;
	    created_at: string;
	    updated_at: string;
	
	    static createFrom(source: any = {}) {
	        return new Scan(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.repository_id = source["repository_id"];
	        this.scan_time = source["scan_time"];
	        this.result = source["result"];
	        this.vulnerabilities = source["vulnerabilities"];
	        this.status = source["status"];
	        this.created_at = source["created_at"];
	        this.updated_at = source["updated_at"];
	    }
	}
	export class ScanCreate {
	    repository_id: number;
	    user_id: number;
	    scan_time: string;
	    result: string;
	    vulnerabilities: number;
	    status: string;
	
	    static createFrom(source: any = {}) {
	        return new ScanCreate(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.repository_id = source["repository_id"];
	        this.user_id = source["user_id"];
	        this.scan_time = source["scan_time"];
	        this.result = source["result"];
	        this.vulnerabilities = source["vulnerabilities"];
	        this.status = source["status"];
	    }
	}
	export class ScanFull {
	    scan: Scan;
	    repository_name: string;
	    repository_path: string;
	    repository_description: string;
	
	    static createFrom(source: any = {}) {
	        return new ScanFull(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.scan = this.convertValues(source["scan"], Scan);
	        this.repository_name = source["repository_name"];
	        this.repository_path = source["repository_path"];
	        this.repository_description = source["repository_description"];
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
	export class User {
	    id: number;
	    username: string;
	    password: string;
	    is_completed_setup: boolean;
	    session_token: string;
	    session_expired: string;
	    created_at: string;
	    updated_at: string;
	
	    static createFrom(source: any = {}) {
	        return new User(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.username = source["username"];
	        this.password = source["password"];
	        this.is_completed_setup = source["is_completed_setup"];
	        this.session_token = source["session_token"];
	        this.session_expired = source["session_expired"];
	        this.created_at = source["created_at"];
	        this.updated_at = source["updated_at"];
	    }
	}
	export class UserLogin {
	    username: string;
	    password: string;
	
	    static createFrom(source: any = {}) {
	        return new UserLogin(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.username = source["username"];
	        this.password = source["password"];
	    }
	}
	export class UserUpdate {
	    id: number;
	    username?: string;
	    password?: string;
	    previous_password?: string;
	    is_completed_setup?: boolean;
	    session_token?: string;
	    session_expired?: string;
	
	    static createFrom(source: any = {}) {
	        return new UserUpdate(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.username = source["username"];
	        this.password = source["password"];
	        this.previous_password = source["previous_password"];
	        this.is_completed_setup = source["is_completed_setup"];
	        this.session_token = source["session_token"];
	        this.session_expired = source["session_expired"];
	    }
	}

}

export namespace parser {
	
	export class PerFileTime {
	    mean: number;
	    std_dev: number;
	
	    static createFrom(source: any = {}) {
	        return new PerFileTime(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.mean = source["mean"];
	        this.std_dev = source["std_dev"];
	    }
	}
	export class ParsingTime {
	    total_time: number;
	    per_file_time: PerFileTime;
	    very_slow_files: any[];
	
	    static createFrom(source: any = {}) {
	        return new ParsingTime(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.total_time = source["total_time"];
	        this.per_file_time = this.convertValues(source["per_file_time"], PerFileTime);
	        this.very_slow_files = source["very_slow_files"];
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
	
	export class Position {
	    line: number;
	    col: number;
	    offset: number;
	
	    static createFrom(source: any = {}) {
	        return new Position(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.line = source["line"];
	        this.col = source["col"];
	        this.offset = source["offset"];
	    }
	}
	export class ProfilingTimes {
	    config_time: number;
	    core_time: number;
	    ignores_time: number;
	    total_time: number;
	
	    static createFrom(source: any = {}) {
	        return new ProfilingTimes(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.config_time = source["config_time"];
	        this.core_time = source["core_time"];
	        this.ignores_time = source["ignores_time"];
	        this.total_time = source["total_time"];
	    }
	}
	export class ResultExtras {
	    message: string;
	    metadata: Record<string, any>;
	    severity: string;
	    fingerprint: string;
	    lines: string;
	    validation_state: string;
	    engine_kind: string;
	
	    static createFrom(source: any = {}) {
	        return new ResultExtras(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.message = source["message"];
	        this.metadata = source["metadata"];
	        this.severity = source["severity"];
	        this.fingerprint = source["fingerprint"];
	        this.lines = source["lines"];
	        this.validation_state = source["validation_state"];
	        this.engine_kind = source["engine_kind"];
	    }
	}
	export class Span {
	    file: string;
	    start: Position;
	    end: Position;
	
	    static createFrom(source: any = {}) {
	        return new Span(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.file = source["file"];
	        this.start = this.convertValues(source["start"], Position);
	        this.end = this.convertValues(source["end"], Position);
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
	export class SemgrepError {
	    code: number;
	    level: string;
	    type: any[];
	    message: string;
	    path: string;
	    spans: Span[];
	
	    static createFrom(source: any = {}) {
	        return new SemgrepError(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.level = source["level"];
	        this.type = source["type"];
	        this.message = source["message"];
	        this.path = source["path"];
	        this.spans = this.convertValues(source["spans"], Span);
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
	export class SemgrepPaths {
	    scanned: string[];
	
	    static createFrom(source: any = {}) {
	        return new SemgrepPaths(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.scanned = source["scanned"];
	    }
	}
	export class SemgrepTime {
	    rules: any[];
	    rules_parse_time: number;
	    profiling_times: ProfilingTimes;
	    parsing_time: ParsingTime;
	    targets: any[];
	    total_bytes: number;
	    max_memory_bytes: number;
	
	    static createFrom(source: any = {}) {
	        return new SemgrepTime(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.rules = source["rules"];
	        this.rules_parse_time = source["rules_parse_time"];
	        this.profiling_times = this.convertValues(source["profiling_times"], ProfilingTimes);
	        this.parsing_time = this.convertValues(source["parsing_time"], ParsingTime);
	        this.targets = source["targets"];
	        this.total_bytes = source["total_bytes"];
	        this.max_memory_bytes = source["max_memory_bytes"];
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
	export class SemgrepResult {
	    check_id: string;
	    path: string;
	    start: Position;
	    end: Position;
	    extra: ResultExtras;
	
	    static createFrom(source: any = {}) {
	        return new SemgrepResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.check_id = source["check_id"];
	        this.path = source["path"];
	        this.start = this.convertValues(source["start"], Position);
	        this.end = this.convertValues(source["end"], Position);
	        this.extra = this.convertValues(source["extra"], ResultExtras);
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
	export class SemgrepReport {
	    version: string;
	    results: SemgrepResult[];
	    errors: SemgrepError[];
	    paths: SemgrepPaths;
	    time: SemgrepTime;
	
	    static createFrom(source: any = {}) {
	        return new SemgrepReport(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.results = this.convertValues(source["results"], SemgrepResult);
	        this.errors = this.convertValues(source["errors"], SemgrepError);
	        this.paths = this.convertValues(source["paths"], SemgrepPaths);
	        this.time = this.convertValues(source["time"], SemgrepTime);
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

export namespace semgrep {
	
	export class ScanResult {
	    ExitCode: number;
	    Stdout: string;
	    Stderr: string;
	
	    static createFrom(source: any = {}) {
	        return new ScanResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ExitCode = source["ExitCode"];
	        this.Stdout = source["Stdout"];
	        this.Stderr = source["Stderr"];
	    }
	}

}

